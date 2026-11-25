project-replaceme
=================

A Go CLI built with Cobra. This file contains the non-obvious context an agent needs to work effectively in this repo — conventions, build commands, and workflow rules that cannot be inferred from reading the code alone.

Project Structure
-----------------

```
src/                    Go source code (module: github.com/owner-replaceme/project-replaceme)
  cmd/                  Cobra command definitions (directory-per-command)
    version/            "version" subcommand
  internal/             Internal packages (not importable by other modules)
    buildinfo/          Build-time version injection via ldflags
  main.go              Entry point — calls cmd.Execute()
_build/                 Build output (gitignored)
.githooks/              Git hooks (configured via make setup)
.claude/                Claude Code settings and skills
.github/workflows/      CI and release automation
```

All Go code lives under `src/`. The module path is `github.com/owner-replaceme/project-replaceme`. The `go.mod` file is in `src/`, so run Go commands from that directory (or use the Makefile, which handles this).

Command Structure
-----------------

Each command lives in its own package under `src/cmd/`, one directory per level of the command hierarchy:

```
src/cmd/
  root.go                  Root command (package cmd)
  version/version.go       "version" subcommand (package version)
  config/                  "config" subcommand group (package config)
    config.go
    get/get.go             "config get" (package get)
    set/set.go             "config set" (package set)
```

Each command package exports:
- `CmdStr` — the command name as a constant
- `Cmd` — the `*cobra.Command` instance

Parent commands wire children via `init()` using `AddCommand()`. Use `/new-command` to scaffold new commands following this pattern.

Building and Running
--------------------

Always build via the Makefile — it injects the version string via ldflags from git state. Running `go build` directly produces a binary that reports version "dev".

```bash
make build      # Full pipeline: configure hooks → check → compile
make compile    # Build the binary only (skip checks)
make run        # Build and run (pass args via ARGS="...")
make check      # Modules → formatting → vet → lint → govulncheck → deadcode → tests (with -race)
make test       # Tests only (with -race)
make clean      # Remove _build/
```

Use `make run` to build and run the binary in one step:

```bash
make run
make run ARGS="version"
```

The sandbox is configured to allow writes to the Go build cache, so `make build`, `make check`, and `make test` work without disabling the sandbox.

Dependencies
------------

To add a Go library dependency, run from the `src/` directory:

```bash
cd src && go get <package>@latest && go mod tidy
```

Do not use `go install` for library dependencies — that installs binaries, not libraries. Use `go get`.

Required Workflows
------------------

These are interactive skills that guide you through a multi-step workflow. Invoke them and follow their instructions — do not attempt to replicate their behavior manually.

### Adding commands

Invoke `/new-command` when adding a CLI command or subcommand. The skill scaffolds the command file, test file, and wires it into the parent following the directory-per-command convention.

### New features

Invoke `/brainstorm` before starting any new feature work. This ensures the design is explored and validated before code is written, preventing wasted effort on approaches that miss requirements or have better alternatives.

### README changes

Invoke `/technical-writer` when modifying the README. The skill produces documentation that is clear, well-structured, and consistent with the project's voice.

### Releasing

Invoke `/releasing` to tag and publish a new version. The skill walks through pre-flight checks, changelog generation, semver version selection (with mandatory user confirmation), and release verification. Do not create version tags manually — the skill ensures all release steps are completed correctly.

Issue Tracking
--------------

Where the `/linear-system` skill is available, work on this repo is tracked in Linear, each issue carrying this repo's label from the `repo` group. Invoke `/linear-system` before reading or writing there; it covers the labels, including a repo that has none yet.
