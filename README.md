# 🗃️ wacli-readonly — WhatsApp CLI (Read-Only Fork)

**Read-only, security-hardened fork of [wacli](https://github.com/steipete/wacli)**. All write/send capabilities are removed by design — not gated behind a flag. Even internal code paths that call write methods will fail with descriptive errors (defense-in-depth).

WhatsApp CLI built on top of `whatsmeow`, focused on:

- ✅ Best-effort local sync of message history + continuous capture
- ✅ Fast offline search (SQLite FTS5)
- ✅ Contact + group viewing
- ✅ Message context, quoted-reply inspection, media download
- ❌ **Sending messages — REMOVED**
- ❌ **Reactions, presence, uploads — REMOVED**
- ❌ **Group management/editing — REMOVED**

This is a third-party tool that uses the WhatsApp Web protocol via `whatsmeow` and is not affiliated with WhatsApp.

## What's Different?

This fork removes all write operations to provide a **read-only monitoring interface**:

### Removed Commands
- ❌ `send text` — Send text messages
- ❌ `send file` — Send files/media
- ❌ `send react` — React to messages
- ❌ `presence typing` / `presence paused` — Send presence indicators
- ❌ `groups rename` — Rename groups
- ❌ `groups participants` — Add/remove/promote/demote members
- ❌ `groups invite` — Manage invite links
- ❌ `groups join` — Join groups
- ❌ `groups leave` — Leave groups
- ❌ `contacts refresh` / alias / tags write commands

### Retained Commands
- ✅ `auth` — Authenticate (QR code), save QR as PNG via `--qr-file`
- ✅ `sync` — Sync message history
- ✅ `messages search` — Full-text search with filters
- ✅ `messages list` / `show` / `context` — Browse messages
- ✅ `chats list` / `show` — Browse chats
- ✅ `contacts search` / `show` — Browse contacts
- ✅ `groups list` / `info` / `refresh` — Browse groups
- ✅ `history backfill` — Backfill older messages
- ✅ `media download` — Download media files
- ✅ `doctor` — Diagnostics

## Use Case

Perfect for:
- **Monitoring WhatsApp without risk of sending messages**
- **Integration with automation tools that should only read**
- **Logging/archiving WhatsApp conversations**
- **Search and analysis without modification capabilities**

## Installation

### Homebrew (macOS & Linux)

```bash
brew tap regalstreak/tap
brew install wacli-readonly
```

### Direct Download

Download the latest release for your platform from the [Releases page](https://github.com/regalstreak/wacli-readonly/releases):

| Platform | File |
|----------|------|
| macOS (Intel & Apple Silicon) | `wacli-readonly-macos-universal.tar.gz` |
| Linux x64 | `wacli-readonly-linux-amd64.tar.gz` |
| Linux ARM64 | `wacli-readonly-linux-arm64.tar.gz` |
| Windows x64 | `wacli-readonly-windows-amd64.zip` |

### Build from Source

```bash
go build -tags sqlite_fts5 -o ./dist/wacli-readonly ./cmd/wacli
```

To save QR codes as PNG, install `qrencode`:

```bash
brew install qrencode         # macOS
sudo apt install qrencode     # Debian/Ubuntu
```

## Quick Start

Default store directory is `~/.local/state/wacli-readonly` on Linux and `~/.wacli-readonly` elsewhere. Override with `--store DIR` or `WACLI_STORE_DIR`.

```bash
# 1) Authenticate (shows QR), then bootstrap sync
wacli-readonly auth

# Save QR as PNG for easier scanning (requires qrencode)
wacli-readonly auth --qr-file /tmp/wa-qr.png

# 2) Keep syncing (never shows QR; requires prior auth)
wacli-readonly sync --follow

# Diagnostics
wacli-readonly doctor

# Search messages
wacli-readonly messages search "meeting"

# List chats
wacli-readonly chats list --limit 50

# List recent messages from a chat
wacli-readonly messages list --chat 1234567890@s.whatsapp.net --asc

# Show context around a message
wacli-readonly messages context --chat 1234567890@s.whatsapp.net --id <message-id>

# View group info (read-only)
wacli-readonly groups info --jid 1234567890@g.us

# Backfill older messages for a chat
wacli-readonly history backfill --chat 1234567890@s.whatsapp.net --requests 10 --count 50

# Download media for a message
wacli-readonly media download --chat 1234567890@s.whatsapp.net --id <message-id>
```

## High-level UX

- `wacli auth`: interactive login (shows QR code), then immediately performs initial data sync.
- `wacli sync`: non-interactive sync loop (never shows QR; errors if not authenticated).
- Output is human-readable by default; pass `--json` for machine-readable output.
- Pass `--full` to keep full IDs in table output; non-TTY output keeps full IDs automatically.

## Command surface

- `wacli auth [--follow] [--idle-exit 30s] [--download-media] [--qr-file PATH]`
- `wacli auth status`
- `wacli auth logout`
- `wacli sync [--once] [--follow] [--idle-exit 30s] [--max-reconnect 5m] [--download-media] [--refresh-contacts] [--refresh-groups]`
- `wacli messages list [--chat JID] [--sender JID] [--from-me|--from-them] [--asc] [--limit N] [--after DATE] [--before DATE]`
- `wacli messages search <query> [--chat JID] [--from JID] [--has-media] [--type text|image|video|audio|document]`
- `wacli messages show --chat JID --id MSG_ID`
- `wacli messages context --chat JID --id MSG_ID [--before N] [--after N]`
- `wacli media download --chat JID --id MSG_ID [--output PATH]`
- `wacli contacts search <query>`
- `wacli contacts show --jid JID`
- `wacli chats list [--query TEXT] [--limit N]`
- `wacli chats show --jid JID`
- `wacli groups list [--query TEXT] [--limit N]`
- `wacli groups refresh`
- `wacli groups info --jid GROUP_JID`
- `wacli history backfill --chat JID [--count 50] [--requests N]`
- `wacli doctor [--connect]`
- `wacli version`

## Storage

Defaults to `~/.local/state/wacli-readonly` on Linux and `~/.wacli-readonly` elsewhere. Existing Linux `~/.wacli-readonly` stores are reused when the XDG state store does not exist. Override with `--store DIR` or `WACLI_STORE_DIR`.

Global flags:

- `--store DIR`: store directory.
- `--json`: JSON output.
- `--full`: disable table truncation.
- `--timeout DURATION`: timeout for non-sync commands.
- `--lock-wait DURATION`: wait for the store lock before failing.

## Environment Overrides

- `WACLI_DEVICE_LABEL`: set the linked device label (shown in WhatsApp). Defaults to `Chrome` to blend in with normal Web sessions.
- `WACLI_DEVICE_PLATFORM`: override the linked device platform (defaults to `CHROME` if unset or invalid).
- `WACLI_STORE_DIR`: override the default store directory.

## Backfilling Older History

`wacli sync` stores whatever WhatsApp Web sends opportunistically. To try to fetch *older* messages, use on-demand history sync requests to your **primary device** (your phone).

Important notes:

- This is **best-effort**: WhatsApp may not return full history.
- Your **primary device must be online**.
- Requests are **per chat** (DM or group). Uses the *oldest locally stored message* in that chat as the anchor.
- Recommended `--count` is `50` per request; maximum is `500`.
- Maximum `--requests` per run is `100`.

### Backfill one chat

```bash
wacli-readonly history backfill --chat 1234567890@s.whatsapp.net --requests 10 --count 50
```

### Backfill all chats (script)

```bash
wacli-readonly --json chats list --limit 100000 \
  | jq -r '.[].JID' \
  | while read -r jid; do
      wacli-readonly history backfill --chat "$jid" --requests 3 --count 50
    done
```

## Security

- All write methods in `internal/wa/Client` return errors unconditionally (`SendText`, `SendProtoMessage`, `SendReaction`, `Upload`, `SetGroupName`, `UpdateGroupParticipants`, `GetGroupInviteLink`, `JoinGroupWithLink`, `LeaveGroup`, `SendChatPresence`).
- The CLI does not register any subcommand that attempts to write.
- Store contains encryption keys/session data; store dir `0700`, DB files `0600`.
- SQLite URI injection and FTS5 query injection are blocked.

## Prior Art / Credit

This is a read-only fork of the excellent `wacli` by Peter Steinberger:
- [`wacli`](https://github.com/steipete/wacli)

Which was inspired by:
- [`whatsapp-cli`](https://github.com/vicentereig/whatsapp-cli) by Vicente Reig

## License

See `LICENSE` (same as upstream wacli).
