package wa

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/types"
)

func (c *Client) GetJoinedGroups(ctx context.Context) ([]*types.GroupInfo, error) {
	c.mu.Lock()
	cli := c.client
	c.mu.Unlock()
	if cli == nil || !cli.IsConnected() {
		return nil, fmt.Errorf("not connected")
	}
	return cli.GetJoinedGroups(ctx)
}

// SetGroupName is disabled in the read-only build.
func (c *Client) SetGroupName(ctx context.Context, jid types.JID, name string) error {
	return fmt.Errorf("group write operations disabled in read-only build")
}

type GroupParticipantAction string

const (
	GroupParticipantAdd     GroupParticipantAction = "add"
	GroupParticipantRemove  GroupParticipantAction = "remove"
	GroupParticipantPromote GroupParticipantAction = "promote"
	GroupParticipantDemote  GroupParticipantAction = "demote"
)

// UpdateGroupParticipants is disabled in the read-only build.
func (c *Client) UpdateGroupParticipants(ctx context.Context, group types.JID, users []types.JID, action GroupParticipantAction) ([]types.GroupParticipant, error) {
	return nil, fmt.Errorf("group write operations disabled in read-only build")
}

// GetGroupInviteLink is disabled in the read-only build.
func (c *Client) GetGroupInviteLink(ctx context.Context, group types.JID, reset bool) (string, error) {
	return "", fmt.Errorf("group invite operations disabled in read-only build")
}

// JoinGroupWithLink is disabled in the read-only build.
func (c *Client) JoinGroupWithLink(ctx context.Context, code string) (types.JID, error) {
	return types.JID{}, fmt.Errorf("group write operations disabled in read-only build")
}

// LeaveGroup is disabled in the read-only build.
func (c *Client) LeaveGroup(ctx context.Context, group types.JID) error {
	return fmt.Errorf("group write operations disabled in read-only build")
}
