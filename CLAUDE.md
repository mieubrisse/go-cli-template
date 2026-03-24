project-replaceme
=================

A Go CLI built with Cobra. This file contains the non-obvious context an agent needs to work effectively in this repo — conventions, build commands, and workflow rules that cannot be inferred from reading the code alone.

Project Structure
-----------------

```
src/                    Go source code (module: github.com/owner-replaceme/project-replaceme)
  cmd/                  Cobra command definitions
  internal/             Internal packages (not importable by other modules)
    version/            Build-time version injection
  main.go              Entry point — calls cmd.Execute()
_build/                 Build output (gitignored)
.githooks/              Git hooks (configured via make setup)
.beads/                 Beads issue tracking config
.claude/                Claude Code settings and skills
.github/workflows/      CI and release automation
```

All Go code lives under `src/`. The module path is `github.com/owner-replaceme/project-replaceme`. The `go.mod` file is in `src/`, so run Go commands from that directory (or use the Makefile, which handles this).

Building and Running
--------------------

Always build via the Makefile — it injects the version string via ldflags from git state. Running `go build` directly produces a binary that reports version "dev".

```bash
make build      # Full pipeline: configure hooks → check → compile
make compile    # Build the binary only (skip checks)
make check      # Formatting (gofmt) → static analysis (go vet) → tests
make test       # Tests only
make clean      # Remove _build/
```

Run the binary from the repo root:

```bash
_build/project-replaceme
_build/project-replaceme version
```

The sandbox is configured to allow writes to the Go build cache, so `make build`, `make check`, and `make test` work without disabling the sandbox.

Dependencies
------------

To add a Go library dependency, run from the `src/` directory:

```bash
cd src && go get <package>@latest && go mod tidy
```

Do not use `go install` for library dependencies — that installs binaries, not libraries. Use `go get`.

Command String Constants
------------------------

All Cobra command and flag names are defined as constants in `src/cmd/command_str_consts.go`. When adding a new command or flag, add its string constant there first, then reference the constant in the command definition. This keeps command names defined in exactly one place.

Required Workflows
------------------

These are interactive skills that guide you through a multi-step workflow. Invoke them and follow their instructions — do not attempt to replicate their behavior manually.

### New features

Invoke `/brainstorm` before starting any new feature work. This ensures the design is explored and validated before code is written, preventing wasted effort on approaches that miss requirements or have better alternatives.

### README changes

Invoke `/technical-writer` when modifying the README. The skill produces documentation that is clear, well-structured, and consistent with the project's voice.

### Releasing

Invoke `/releasing` to tag and publish a new version. The skill walks through pre-flight checks, changelog generation, semver version selection (with mandatory user confirmation), and release verification. Do not create version tags manually — the skill ensures all release steps are completed correctly.

Issue Tracking
--------------

This project uses **bd** (Beads) for issue tracking, configured in shared Dolt server mode. The server runs at `127.0.0.1` and auto-starts when needed.

```bash
bd prime          # Load workflow context
bd ready          # Find unblocked work
bd create "Title" # Create an issue
bd close <id>     # Complete an issue
```

Run `bd prime` for full workflow details.

If `bd` is not installed, issue tracking commands will fail. Install it before using any `bd` commands. The git hooks handle a missing `bd` gracefully — they skip beads integration and continue — so a missing `bd` will not block commits.
