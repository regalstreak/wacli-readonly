package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/steipete/wacli/internal/out"
	"go.mau.fi/whatsmeow/types"
)

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
