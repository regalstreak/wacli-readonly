# wacli-readonly specification

This document defines the specification for `wacli-readonly`: a **read-only** WhatsApp CLI that syncs messages locally and supports fast search. This is a security-hardened fork of [steipete/wacli](https://github.com/steipete/wacli) with all write operations (send messages, modify groups, etc.) removed by design. Implementation uses `whatsmeow` under the hood.

## Goals

- **Read-only by design**: no send-message, no group-write, no reactions, no presence-send. Write methods at the `wa.Client` layer return errors even if called directly (defense-in-depth).
- **Explicit authentication step**: `wacli auth` shows a QR code and completes login.
- **Auth starts syncing immediately**: after successful QR pairing, `wacli auth` begins initial sync (history + metadata).
- **Non-interactive sync**: `wacli sync` never displays a QR code; it fails with a clear error if not authenticated.
- **Fast offline message search**: local SQLite + FTS5 index.
- **Human-first output**: readable tables by default, `--json` opt-in for scripting.
- **Single-instance safety**: store locking to avoid multi-instance session conflicts.
- **Group management (read-only)**: list groups, inspect, refresh.

## Non-goals

- Sending messages, reactions, presence updates.
- Creating, renaming, joining, leaving groups.
- Managing group participants or invite links.
- End-to-end “contact creation” in WhatsApp.

## Terminology

- **JID**: WhatsApp Jabber ID, e.g. `1234567890@s.whatsapp.net` (user) or `123456789@g.us` (group).
- **Store directory**: directory containing all local state, default `~/.local/state/wacli-readonly` on Linux and `~/.wacli-readonly` elsewhere.

## Storage layout

Default store: `~/.local/state/wacli-readonly` on Linux and `~/.wacli-readonly` elsewhere (override with `--store DIR` or `WACLI_STORE_DIR`). Existing Linux `~/.wacli-readonly` stores are reused when the XDG state store does not exist.

Files:

- `<store>/session.db` — `whatsmeow` SQL store (device identity, keys, app-state).
- `<store>/wacli.db` — our SQLite DB (messages/chats, FTS, local metadata).
- `<store>/media/...` — downloaded media (optional, on-demand or background).
- `<store>/LOCK` — store lock to prevent concurrent access.

## Concurrency + locking

Every command that accesses the WhatsApp session must acquire an exclusive lock in the store dir.

## Authentication model

### Commands

- `wacli auth` (interactive)
  - If not authenticated: connect, show QR code, wait for success.
  - Optional `--qr-file PATH` saves the QR code as a PNG (requires `qrencode`).
  - After success: start initial sync (bootstrap) immediately.
  - Exits after initial sync “goes idle” (configurable), unless `--follow` is set.

- `wacli sync` (non-interactive)
  - Requires an existing authenticated session in `session.db`.
  - Never displays QR; if not authenticated, prints “run `wacli auth`”.
  - `--once` performs a bounded sync and exits.
  - Default (or `--follow`) stays connected and continues capturing messages.

### Device label

By default, `wacli-readonly` registers the linked device as "Chrome" to avoid detection. Override via `WACLI_DEVICE_LABEL` and `WACLI_DEVICE_PLATFORM` environment variables.

## CLI command surface

Global flags:

- `--store DIR` (default: XDG state dir on Linux, `~/.wacli-readonly` elsewhere; or `WACLI_STORE_DIR`)
- `--json` (default: human text)
- `--full` (disable table truncation; non-TTY output keeps full IDs)
- `--timeout DURATION` (non-sync commands; e.g. `5m`)
- `--lock-wait DURATION` (wait for the store lock before failing)
- `--version` (prints version and exits)

### Doctor

- `wacli doctor [--connect]`

### Auth

- `wacli auth [--follow] [--idle-exit 30s] [--qr-file PATH]`
- `wacli auth status`
- `wacli auth logout`

### Sync

- `wacli sync [--once] [--follow] [--download-media]`

### History backfill (best-effort)

- `wacli history backfill --chat JID [--count 50] [--requests N]`
- Backfill caps: `--count <= 500`, `--requests <= 100`.

### Messages (read-only)

- `wacli messages list [--chat JID] [--sender JID] [--from-me|--from-them] [--asc] [--limit N] [--before TS] [--after TS]`
- `wacli messages search <query> [--chat JID] [--from JID] [--limit N] [--before TS] [--after TS] [--type text|image|video|audio|document]`
- `wacli messages show --chat JID --id MSG_ID`
- `wacli messages context --chat JID --id MSG_ID [--before N] [--after N]`

### Contacts (read-only)

- `wacli contacts search <query>`
- `wacli contacts show --jid JID`

### Chats (read-only)

- `wacli chats list [--query TEXT]`
- `wacli chats show --jid JID`

### Groups (read-only)

- `wacli groups list [--query TEXT]`
- `wacli groups refresh`
- `wacli groups info --jid GROUP_JID`

### Media (read-only)

- `wacli media download ...`

## Output formats

Default: human-readable text (tables / aligned columns; TTY-aware wrapping).

Optional:

- `--json` prints `{"success":true,"data":...,"error":null}`-style responses.

## Security considerations

- Store contains encryption keys/session data; permissions are enforced:
  - store dir `0700`
  - DB files `0600`
- All write-capable methods in `internal/wa/Client` (SendText, SendProtoMessage, SendReaction, Upload, SetGroupName, UpdateGroupParticipants, GetGroupInviteLink, JoinGroupWithLink, LeaveGroup, SendChatPresence) return clear errors. Defense-in-depth: even if internal code paths call them, they fail safely.
- SQLite URI injection via StorePath is blocked.
- FTS5 queries are sanitized to prevent query-syntax injection.

## Prior art / credit

This fork is based on [steipete/wacli](https://github.com/steipete/wacli), which itself borrows ideas from `https://github.com/vicentereig/whatsapp-cli`.
