# Changes from Original wacli

This document lists all modifications made to create the read-only fork.

## Files Modified

### 1. `cmd/wacli/root.go`
**Line 44:** Commented out send command registration
```go
// Send command removed - read-only version
// rootCmd.AddCommand(newSendCmd(&flags))
```

### 2. `cmd/wacli/groups.go`
**Lines 18-24:** Commented out write operations
```go
// Write operations removed - read-only version
// cmd.AddCommand(newGroupsRenameCmd(flags))
// cmd.AddCommand(newGroupsParticipantsCmd(flags))
// cmd.AddCommand(newGroupsInviteCmd(flags))
// cmd.AddCommand(newGroupsJoinCmd(flags))
// cmd.AddCommand(newGroupsLeaveCmd(flags))
```

**Line 11:** Updated description
```go
Short: "Group management (read-only)",
```

### 3. `README.md`
Completely rewritten to document read-only nature and removed features.

## Removed Capabilities

### Send Commands (Completely Disabled)
- `wacli send text` - Send text messages
- `wacli send file` - Send files/attachments

### Group Write Operations (Disabled)
- `wacli groups rename` - Rename group
- `wacli groups participants add|remove|promote|demote` - Manage participants
- `wacli groups invite link revoke` - Revoke invite links
- `wacli groups join` - Join groups via invite
- `wacli groups leave` - Leave groups

## Retained Capabilities

All read-only operations work normally:
- ✅ Authentication (`auth`)
- ✅ Message syncing (`sync`)
- ✅ Message search (`messages search`)
- ✅ Message listing (`messages list`)
- ✅ Chat listing (`chats list`)
- ✅ Contact management (`contacts list`)
- ✅ Group viewing (`groups list`, `groups info`, `groups refresh`)
- ✅ History backfilling (`history backfill`)
- ✅ Media downloads (`media download`)
- ✅ Diagnostics (`doctor`)

## Build Instructions

```bash
cd wacli-readonly
go build -tags sqlite_fts5 -o ./dist/wacli-readonly ./cmd/wacli
```

## Testing

After building, verify removed commands are unavailable:

```bash
# These should fail with "unknown command"
./dist/wacli-readonly send
./dist/wacli-readonly groups rename
./dist/wacli-readonly groups join

# These should work
./dist/wacli-readonly --help
./dist/wacli-readonly messages --help
./dist/wacli-readonly groups list --help
```

## Integration with OpenClaw

To use with OpenClaw for WhatsApp monitoring:

1. Build this read-only version
2. Authenticate: `./dist/wacli-readonly auth`
3. Set up sync in background: `./dist/wacli-readonly sync --follow`
4. Configure OpenClaw to read from `~/.wacli` database
5. Use `messages search` and `chats list` for queries

The read-only nature ensures OpenClaw agents cannot accidentally send messages to WhatsApp contacts.

## Future Maintenance

When syncing with upstream wacli:

1. Check for new write commands in `root.go` and `groups.go`
2. Comment out any new write operations
3. Update `README.md` and this file with changes
4. Test that build still works

## Why Read-Only?

This fork was created to provide a **safe monitoring interface** for automation systems where:
- Agents should observe conversations but never send
- Risk of accidental messages must be eliminated
- Logging and analysis are the primary use cases
- Integration with tools like OpenClaw requires guaranteed read-only access
