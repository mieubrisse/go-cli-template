Go CLI Template Repo Design
============================

Date: 2026-03-24

Goal
----

Create a GitHub template repo for Go CLI projects using Cobra. The template produces a working,
compilable CLI that an agent transforms into a real project via the `/initialize-repo` skill.
The template includes CI/CD, Homebrew publishing via GoReleaser, Claude Code configuration,
Beads integration, and Git hooks.

Placeholder Strategy
--------------------

The template ships as a compilable Go CLI using placeholder names:

- `owner-replaceme` — GitHub username/org
- `project-replaceme` — repo/binary/project name

These are real, valid identifiers so the template builds and tests pass as-is. The
`/initialize-repo` skill performs a global find-and-replace as its first step.

Repo Structure
--------------

```
project-replaceme/
├── .claude/
│   ├── settings.json
│   └── skills/
│       ├── initialize-repo/
│       │   └── SKILL.md
│       └── releasing/
│           └── SKILL.md
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── release.yml
├── .githooks/
│   ├── pre-commit
│   ├── prepare-commit-msg
│   ├── post-checkout
│   ├── post-merge
│   └── pre-push
├── .goreleaser.yaml
├── .beads/
│   └── config.yaml              # issue-prefix deliberately broken
├── src/
│   ├── cmd/
│   │   ├── command_str_consts.go
│   │   ├── root.go
│   │   └── version.go
│   ├── internal/
│   │   └── version/
│   │       └── version.go
│   ├── main.go
│   ├── go.mod
│   └── go.sum
├── _build/                       # gitignored
├── .gitignore
├── Makefile
├── CLAUDE.md
├── SETUP_CHECKLIST.md
└── README.md
```

Components
----------

### Go CLI (src/)

- `main.go` — entry point, calls `cmd.Execute()`
- `cmd/root.go` — Cobra root command with placeholder description
- `cmd/version.go` — version subcommand, prints build-time injected version
- `cmd/command_str_consts.go` — centralized string constants for command/flag names
- `internal/version/version.go` — holds `Version` variable set via ldflags
- `go.mod` — module `github.com/owner-replaceme/project-replaceme`

### Makefile

Targets (all Go commands run from `src/` directory):

- `setup` — configures git hooks path to `.githooks/`
- `check` — gofmt + go vet + go test
- `compile` — go build with ldflags → `_build/project-replaceme`
- `build` — setup + check + compile
- `test` — go test ./...
- `clean` — rm -rf _build/

Version detection from git state (tag > commit hash > dirty), injected via ldflags:
`-X github.com/owner-replaceme/project-replaceme/internal/version.Version=$(VERSION)`

### GoReleaser (.goreleaser.yaml)

- Version 2 config
- Builds for darwin/linux × amd64/arm64, CGO_ENABLED=0
- Binary output from `src/` directory
- Publishes to `owner-replaceme/homebrew-project-replaceme` tap
- Uses `HOMEBREW_TAP_TOKEN` env var
- Changelog filtering (excludes docs, test, ci, chore)
- Checksum generation

### GitHub Actions

**ci.yml** — triggers on push to main + PRs:
- Setup Go from go.mod
- Run `make check`

**release.yml** — triggers on `v*` tags:
- Checkout with full history (fetch-depth: 0)
- Setup Go from go.mod
- GoReleaser with GITHUB_TOKEN and HOMEBREW_TAP_TOKEN

### Claude Code Configuration (.claude/settings.json)

Permissions:
- Allow: Read/Edit/Write/Glob/Grep on `./**`, Bash for `go`, `make`, `gofmt`,
  `golangci-lint`, `bd`, `goreleaser`
- Deny: `.env`, `.env.*`, `secrets/**`, `rm -rf`, `sudo`, `git commit --no-verify`

Sandbox:
- Filesystem: allow write to `/tmp`, `/private/tmp/claude`
- Network: allow `127.0.0.1`, `localhost` (for Beads Dolt shared server)
- Unix sockets: `/tmp/claude`, `/private/tmp/claude`
- Excluded commands: `bd`, `dolt`

Hooks:
- SessionStart: `bd prime`
- PreCompact: `bd prime`

### Git Hooks (.githooks/)

**pre-commit:**
1. Beads prefix guard — fails if `.beads/config.yaml` still contains
   `REPLACE-ME-WITH-YOUR-PROJECT-PREFIX`. Removed during `/initialize-repo`.
2. Run `make check` (formatting, vet, tests)
3. Run beads pre-commit hook (30s timeout, graceful fallback)

**Other hooks** (prepare-commit-msg, post-checkout, post-merge, pre-push):
- Beads hooks only, 30s timeout, graceful fallback if bd unavailable

### Beads (.beads/)

Initialized via `bd init --shared-server` during template creation. The prefix in
`.beads/config.yaml` is deliberately set to `REPLACE-ME-WITH-YOUR-PROJECT-PREFIX`
which is invalid and will cause bd commands to fail, forcing users to set the correct
prefix during initialization.

### Skills

**`/initialize-repo`** — runs the SETUP_CHECKLIST.md step by step:

1. Ask for GitHub owner and project name
2. Find-and-replace `owner-replaceme` and `project-replaceme` across all files
3. Set beads prefix in `.beads/config.yaml` to the project name
4. Remove the beads prefix guard from `.githooks/pre-commit`
5. Choose and add a LICENSE file
6. Create Homebrew tap repo via `gh repo create`
7. Enable tag protection rules on the code repo (`v*` pattern)
8. **User action:** create HOMEBREW_TAP_TOKEN PAT and add as repo secret
9. Verify: `make build` passes, `goreleaser check` validates
10. Delete `SETUP_CHECKLIST.md` and `.claude/skills/initialize-repo/`
11. Commit cleanup

**`/releasing`** — semver tagging and release workflow:

1. Pre-flight checks (clean tree, on default branch, synced with remote)
2. Determine current version from git tags
3. Analyze commits since last tag, build changelog
4. Propose semver bump (major/minor/patch)
5. **User must explicitly confirm the version** — no auto-tagging
6. Create tag and push
7. Wait for release CI to complete successfully

### CLAUDE.md

Contains:
- Project structure overview (src/ layout, build output in _build/)
- How to build: `make build`
- How to run: `_build/project-replaceme`
- "Use `/technical-writer` for any README changes"
- "Use `/brainstorm` before any new feature work"
- "Use `/releasing` for version tagging and releases"

### README.md

Opens with a template usage callout:

> **This is a template repo.** To create a new CLI project from it:
> 1. Click "Use this template" on GitHub to create your repo
> 2. Clone your new repo and run Claude Code in it
> 3. Tell the agent: `/initialize-repo`

Followed by (using placeholder names, updated during initialization):
- Installation via Homebrew
- Upgrading
- Development setup (clone, make build, make check)

Written via `/technical-writer` skill.

Build & Release Flow
--------------------

**Local development:**
```
make build        → full pipeline (setup hooks, check, compile)
make compile      → just build the binary
make check        → lint + vet + test
_build/project-replaceme   → run the CLI
```

**Release:**
```
/releasing        → agent guides through semver tag + push
                  → GitHub Actions release.yml triggers
                  → GoReleaser builds + publishes to Homebrew tap
```

**CI:**
```
Push to main/PR   → ci.yml → make check
Push v* tag        → release.yml → GoReleaser
```

Error Handling
--------------

- All hook scripts use `set -euo pipefail` with graceful fallback if bd unavailable
- Makefile `check` fails fast on first issue
- Beads prefix guard blocks commits until initialization is complete
- `/initialize-repo` clearly separates agent-executable steps from user-required actions
- GoReleaser config validated via `goreleaser check`
- `/releasing` skill requires explicit user confirmation before any tag is created

Testing the Template
--------------------

The template itself is a valid, compilable Go project:
- `make build` succeeds on the template as-is
- CI can run against the template repo to catch breakage
- `goreleaser check` validates the release config
