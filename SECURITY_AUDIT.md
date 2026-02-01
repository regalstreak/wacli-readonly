# Security Audit - wacli-readonly

## Audit Date
2026-02-01 21:44 IST

## Objective
Ensure no write operations are available and that an agent cannot discover or use write functionality.

## Findings

### ❌ CRITICAL: Write Functions Still Present in Codebase

Despite removing CLI commands, the underlying write functions remain in internal packages:

#### 1. **internal/wa/client.go** - Send Functions
- `SendText()` (line 175)
- `SendProtoMessage()` (line 184)
- `Upload()` (line 194)

#### 2. **internal/wa/groups.go** - Group Write Operations
- `SetGroupName()`
- `UpdateGroupParticipants()`
- `JoinGroupWithLink()`
- `LeaveGroup()`

#### 3. **internal/app/app.go** - WAClient Interface
Interface still declares all write methods:
- `SetGroupName`
- `UpdateGroupParticipants`
- `GetGroupInviteLink`
- `JoinGroupWithLink`
- `LeaveGroup`
- `SendText`
- `SendProtoMessage`
- `Upload`

### Current Protection Level: **WEAK** ⚠️

**Why it's a problem:**
- An agent with code access could import `internal/wa` package
- Functions can be called directly without CLI commands
- No runtime protection, only CLI layer removed

### ✅ What's Actually Removed
- ✅ CLI commands: `send`, `groups rename/participants/join/leave`
- ✅ Command files: `send.go`, `send_file.go`, `send_file_cmd.go` deleted
- ✅ No calls to write functions from `cmd/` directory

## Remediation Required

### Option 1: Stub Out Write Functions (Recommended)
Replace write function bodies with error returns:

```go
func (c *Client) SendText(ctx context.Context, to types.JID, text string) (types.MessageID, error) {
    return "", fmt.Errorf("send operations disabled in read-only build")
}
```

**Pros:**
- Maintains interface compatibility
- Clear error messages
- Build still compiles
- No risk of accidental sends

### Option 2: Remove Functions Entirely
Delete all write functions from codebase.

**Cons:**
- Breaks interface
- Requires updating fake_wa_test.go
- More invasive changes

### Option 3: Build Tag Protection
Add build tags to conditionally compile write functions only when explicitly enabled.

**Cons:**
- Complex to implement
- Easy to bypass with wrong build flags

## Recommended Actions

1. **Stub out all write functions** in:
   - `internal/wa/client.go`: SendText, SendProtoMessage, Upload
   - `internal/wa/groups.go`: SetGroupName, UpdateGroupParticipants, JoinGroupWithLink, LeaveGroup
   
2. **Update interface** in `internal/app/app.go`:
   - Keep interface methods but document as disabled
   - OR remove from interface (requires test updates)

3. **Add clear documentation** in all stubbed functions

4. **Update tests** in `internal/app/fake_wa_test.go`

## Verification Steps

After remediation:
1. ✅ Grep for function calls: no calls to write functions from cmd/
2. ✅ Attempt to use write functions: should return clear errors
3. ✅ Build succeeds without warnings
4. ✅ Read operations still work correctly

## Current Safety Assessment

**Before remediation:** ⚠️ **UNSAFE** - Write functions accessible via code
**After stubbing:** ✅ **SAFE** - Write functions return errors at runtime

---

## Update Log

### 2026-02-01 21:44 - Initial Audit
- Identified write functions in internal packages
- Documented risk level
- Proposed remediation strategy
