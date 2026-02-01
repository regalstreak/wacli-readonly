# 🗃️ wacli-readonly — WhatsApp CLI (Read-Only Fork)

**Read-only fork of [wacli](https://github.com/steipete/wacli)** with all write/send capabilities removed.

WhatsApp CLI built on top of `whatsmeow`, focused on:

- ✅ Best-effort local sync of message history + continuous capture
- ✅ Fast offline search
- ✅ Contact + group viewing
- ❌ **Sending messages (REMOVED)**
- ❌ **Group management/editing (REMOVED)**

This is a third-party tool that uses the WhatsApp Web protocol via `whatsmeow` and is not affiliated with WhatsApp.

## What's Different?

This fork removes all write operations to provide a **read-only monitoring interface**:

### Removed Commands:
- ❌ `send text` - Send text messages
- ❌ `send file` - Send files/media
- ❌ `groups rename` - Rename groups
- ❌ `groups participants` - Add/remove/promote/demote members
- ❌ `groups invite` - Manage invite links
- ❌ `groups join` - Join groups
- ❌ `groups leave` - Leave groups

### Retained Commands:
- ✅ `auth` - Authenticate (QR code)
- ✅ `sync` - Sync message history
- ✅ `messages search` - Search messages
- ✅ `messages list` - List messages
- ✅ `chats list` - List chats
- ✅ `contacts list` - List contacts
- ✅ `groups list` - View groups
- ✅ `groups info` - View group details
- ✅ `groups refresh` - Refresh group list
- ✅ `history backfill` - Backfill older messages
- ✅ `media download` - Download media files
- ✅ `doctor` - Diagnostics

## Use Case

Perfect for:
- **Monitoring WhatsApp without risk of sending messages**
- **Integration with automation tools that should only read**
- **Logging/archiving WhatsApp conversations**
- **Search and analysis without modification capabilities**

## Build

```bash
go build -tags sqlite_fts5 -o ./dist/wacli-readonly ./cmd/wacli
```

## Quick Start

Default store directory is `~/.wacli` (override with `--store DIR`).

```bash
# 1) Authenticate (shows QR), then bootstrap sync
./dist/wacli-readonly auth

# 2) Keep syncing (never shows QR; requires prior auth)
./dist/wacli-readonly sync --follow

# Diagnostics
./dist/wacli-readonly doctor

# Search messages
./dist/wacli-readonly messages search "meeting"

# List chats
./dist/wacli-readonly chats list --limit 50

# View group info (read-only)
./dist/wacli-readonly groups info --jid 1234567890@g.us

# Backfill older messages for a chat
./dist/wacli-readonly history backfill --chat 1234567890@s.whatsapp.net --requests 10 --count 50

# Download media for a message
./dist/wacli-readonly media download --chat 1234567890@s.whatsapp.net --id <message-id>
```

## Storage

Defaults to `~/.wacli` (override with `--store DIR`).

## Environment Overrides

- `WACLI_DEVICE_LABEL`: set the linked device label (shown in WhatsApp).
- `WACLI_DEVICE_PLATFORM`: override the linked device platform (defaults to `CHROME` if unset or invalid).

## Backfilling Older History

`wacli sync` stores whatever WhatsApp Web sends opportunistically. To try to fetch *older* messages, use on-demand history sync requests to your **primary device** (your phone).

Important notes:

- This is **best-effort**: WhatsApp may not return full history.
- Your **primary device must be online**.
- Requests are **per chat** (DM or group). Uses the *oldest locally stored message* in that chat as the anchor.
- Recommended `--count` is `50` per request.

### Backfill one chat

```bash
./dist/wacli-readonly history backfill --chat 1234567890@s.whatsapp.net --requests 10 --count 50
```

### Backfill all chats (script)

```bash
./dist/wacli-readonly --json chats list --limit 100000 \
  | jq -r '.[].JID' \
  | while read -r jid; do
      ./dist/wacli-readonly history backfill --chat "$jid" --requests 3 --count 50
    done
```

## Prior Art / Credit

This is a fork of the excellent `wacli` by Peter Steinberger:
- [`wacli`](https://github.com/steipete/wacli)

Which was inspired by:
- [`whatsapp-cli`](https://github.com/vicentereig/whatsapp-cli) by Vicente Reig

## License

See `LICENSE` (same as upstream wacli).
