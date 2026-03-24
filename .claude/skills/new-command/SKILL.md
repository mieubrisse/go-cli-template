---
name: new-command
description: >-
  Scaffolds a new Cobra command following the directory-per-command convention.
  Creates the command file, test file, and wires it into the parent command.
  Invoke when adding a new CLI command or subcommand.
---

Scaffold a New Command
======================

Gather Input
------------

Ask the user for:

1. **Command path** — space-separated hierarchy, e.g., `status` or `config get` or `service logs tail`
2. **Short description** — one line for the leaf command's Cobra `Short` field

Parse the Path
--------------

Split the command path into segments. Each segment maps to a directory under `src/cmd/`:

- The **last segment** is the leaf command (has a `Run` function, gets a test file)
- **Intermediate segments** are parent commands (no `Run` function, wire children in `init()`)
- The **root** (`src/cmd/root.go`, package `cmd`) is the parent of all top-level commands

Read the module path from the `module` line in `src/go.mod` to construct import paths.

Before creating any files, check whether the target files already exist. If the leaf command file already exists, stop and tell the user — do not overwrite existing commands.

<examples>
<example>
**Input:** `status` with description "Show current status"

**Creates:**
- `src/cmd/status/status.go` — leaf command (package `status`)
- `src/cmd/status/status_test.go` — test file

**Wires into:** `src/cmd/root.go` — add import and `rootCmd.AddCommand(status.Cmd)`
</example>

<example>
**Input:** `config get` with description "Get a configuration value"

**Creates (if `config` parent does not exist yet):**
- `src/cmd/config/config.go` — parent command (package `config`), Short: `"config commands"`
- `src/cmd/config/get/get.go` — leaf command (package `get`)
- `src/cmd/config/get/get_test.go` — test file

**Wires:**
- `config.go` wires `get.Cmd` in its `init()`
- `root.go` wires `config.Cmd` in its `init()`

**If `config` parent already exists**, only create the `get/` leaf and add the wire to `config.go`.
</example>
</examples>

Reference Implementation
-------------------------

Read these files before scaffolding — they are the canonical examples of the command pattern in this repo:

- `src/cmd/version/version.go` — leaf command (exports `CmdStr` and `Cmd`, uses `Run`)
- `src/cmd/version/version_test.go` — command test with output capture
- `src/cmd/root.go` — root command (wires children in `init()`, uses unexported `cmdStr` because nothing imports it)

Match the style and conventions of these files exactly.

File Templates
--------------

### Leaf command

```go
package <segment>

import (
	"github.com/spf13/cobra"
)

const CmdStr = "<segment>"

var Cmd = &cobra.Command{
	Use:   CmdStr,
	Short: "<user-provided description>",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: implement
	},
}
```

### Parent command

```go
package <segment>

import (
	"<module-path>/cmd/<path>/<child-segment>"
	"github.com/spf13/cobra"
)

const CmdStr = "<segment>"

var Cmd = &cobra.Command{
	Use:   CmdStr,
	Short: "<segment> commands",
}

func init() {
	Cmd.AddCommand(<child-segment>.Cmd)
}
```

### Test file (leaf commands only)

Scaffold a smoke test with output capture. The user will add meaningful assertions after the command is implemented.

```go
package <segment>

import (
	"bytes"
	"testing"
)

func TestCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)

	Cmd.SetArgs([]string{})
	if err := Cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// TODO: assert on buf.String() once the command produces output
}
```

Wire Into the Parent
--------------------

Read the parent command file before modifying it. Then add the import and `AddCommand()` call:

- If the parent already has an `init()` function, add the new `AddCommand()` call to it
- If the parent has no `init()` function, create one
- The root command is special: it uses `rootCmd.AddCommand(...)` instead of `Cmd.AddCommand(...)`

Verify
------

Run `make check` to confirm the new command compiles, passes linting, and all tests pass. If failures are unrelated to the new command, report them to the user rather than attempting to fix pre-existing issues.
