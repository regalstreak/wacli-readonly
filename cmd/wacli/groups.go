package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/steipete/wacli/internal/out"
	"github.com/steipete/wacli/internal/store"
	"github.com/steipete/wacli/internal/wa"
	"go.mau.fi/whatsmeow/types"
)

func newGroupsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "Group management (read-only)",
	}
	cmd.AddCommand(newGroupsListCmd(flags))
	cmd.AddCommand(newGroupsRefreshCmd(flags))
	cmd.AddCommand(newGroupsInfoCmd(flags))
	return cmd
}

func newGroupsRefreshCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Fetch joined groups (live) and update local DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withTimeout(context.Background(), flags)
			defer cancel()

			a, lk, err := newApp(ctx, flags, true, false)
			if err != nil {
				return err
			}
			defer closeApp(a, lk)

			if err := a.EnsureAuthed(); err != nil {
				return err
			}
			if err := a.Connect(ctx, false, nil); err != nil {
				return err
			}

			gs, err := a.WA().GetJoinedGroups(ctx)
			if err != nil {
				return err
			}
			for _, g := range gs {
				if g == nil {
					continue
				}
				_ = persistGroupInfo(a.DB(), g)
				_ = a.DB().UpsertChat(g.JID.String(), "group", g.GroupName.Name, time.Now())
			}

			if flags.asJSON {
				return out.WriteJSON(os.Stdout, map[string]any{"groups": len(gs)})
			}
			fmt.Fprintf(os.Stdout, "Imported %d groups.\n", len(gs))
			return nil
		},
	}
	return cmd
}

func newGroupsListCmd(flags *rootFlags) *cobra.Command {
	var query string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List known groups (from local DB; run sync to populate)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withTimeout(context.Background(), flags)
			defer cancel()

			a, lk, err := newApp(ctx, flags, false, false)
			if err != nil {
				return err
			}
			defer closeApp(a, lk)

			gs, err := a.DB().ListGroups(query, limit)
			if err != nil {
				return err
			}
			if flags.asJSON {
				return out.WriteJSON(os.Stdout, gs)
			}

			w := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tJID\tCREATED")
			for _, g := range gs {
				name := g.Name
				if name == "" {
					name = g.JID
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", truncate(name, 40), g.JID, g.CreatedAt.Local().Format("2006-01-02"))
			}
			_ = w.Flush()
			return nil
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "search query")
	cmd.Flags().IntVar(&limit, "limit", 50, "limit")
	return cmd
}

func newGroupsInfoCmd(flags *rootFlags) *cobra.Command {
	var jidStr string
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Fetch group info (live) and update local DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(jidStr) == "" {
				return fmt.Errorf("--jid is required")
			}
			ctx, cancel := withTimeout(context.Background(), flags)
			defer cancel()

			a, lk, err := newApp(ctx, flags, true, false)
			if err != nil {
				return err
			}
			defer closeApp(a, lk)

			if err := a.EnsureAuthed(); err != nil {
				return err
			}
			if err := a.Connect(ctx, false, nil); err != nil {
				return err
			}

			gjid, err := types.ParseJID(jidStr)
			if err != nil {
				return err
			}
			info, err := a.WA().GetGroupInfo(ctx, gjid)
			if err != nil {
				return err
			}
			if info != nil {
				_ = persistGroupInfo(a.DB(), info)
			}

			if flags.asJSON {
				return out.WriteJSON(os.Stdout, info)
			}

			fmt.Fprintf(os.Stdout, "JID: %s\nName: %s\nOwner: %s\nCreated: %s\nParticipants: %d\n",
				info.JID.String(),
				info.GroupName.Name,
				info.OwnerJID.String(),
				info.GroupCreated.Local().Format(time.RFC3339),
				len(info.Participants),
			)
			return nil
		},
	}
	cmd.Flags().StringVar(&jidStr, "jid", "", "group JID (…@g.us)")
	return cmd
}

func persistGroupInfo(db *store.DB, info *types.GroupInfo) error {
	if info == nil {
		return nil
	}
	if err := db.UpsertGroup(info.JID.String(), info.GroupName.Name, info.OwnerJID.String(), info.GroupCreated); err != nil {
		return err
	}
	var ps []store.GroupParticipant
	for _, p := range info.Participants {
		role := "member"
		if p.IsSuperAdmin {
			role = "superadmin"
		} else if p.IsAdmin {
			role = "admin"
		}
		ps = append(ps, store.GroupParticipant{
			GroupJID: info.JID.String(),
			UserJID:  p.JID.String(),
			Role:     role,
		})
	}
	return db.ReplaceGroupParticipants(info.JID.String(), ps)
}
