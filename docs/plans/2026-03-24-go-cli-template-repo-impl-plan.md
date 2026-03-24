Go CLI Template Repo Implementation Plan
=========================================

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a GitHub template repo that produces a working Go CLI with Cobra, GoReleaser, Homebrew publishing, Claude Code config, Beads integration, and Git hooks.

**Architecture:** All Go code lives under `src/` with module path `github.com/owner-replaceme/project-replaceme`. Build output goes to `_build/` (gitignored). Placeholder names `owner-replaceme` and `project-replaceme` are replaced by the `/initialize-repo` skill post-cloning. The template is a valid, compilable Go project as-is.

**Tech Stack:** Go, Cobra CLI framework, GoReleaser, GitHub Actions, Beads (bd), Make

---

### Task 1: Initialize the Git repo and create .gitignore

**Files:**
- Create: `.gitignore`

**Step 1: Initialize the git repo**

Run:
```bash
git init
```

**Step 2: Create .gitignore**

```gitignore
# Build output
_build/

# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Environment
.env
.env.*
```

**Step 3: Commit**

```bash
git add .gitignore
git commit -m "Add .gitignore"
```

---

### Task 2: Create the Go module and version package

**Files:**
- Create: `src/go.mod`
- Create: `src/internal/version/version.go`

**Step 1: Create go.mod**

```bash
cd src
go mod init github.com/owner-replaceme/project-replaceme
cd ..
```

**Step 2: Create the version package**

Write `src/internal/version/version.go`:

```go
package version

// Version is set at build time via ldflags.
var Version = "dev"
```

**Step 3: Commit**

```bash
git add src/go.mod src/internal/version/version.go
git commit -m "Add Go module and version package"
```

---

### Task 3: Create Cobra CLI scaffolding

**Files:**
- Create: `src/cmd/command_str_consts.go`
- Create: `src/cmd/root.go`
- Create: `src/cmd/version.go`
- Create: `src/main.go`

**Step 1: Create command string constants**

Write `src/cmd/command_str_consts.go`:

```go
package cmd

// Centralized command name strings. Use these constants in Cobra Use fields
// and user-facing messages so that command names are defined in one place.
const (
	rootCmdStr   = "project-replaceme"
	versionCmdStr = "version"
)
```

**Step 2: Create root command**

Write `src/cmd/root.go`:

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   rootCmdStr,
	Short: "project-replaceme CLI",
	Long:  "project-replaceme is a command-line tool.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

**Step 3: Create version command**

Write `src/cmd/version.go`:

```go
package cmd

import (
	"fmt"

	"github.com/owner-replaceme/project-replaceme/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   versionCmdStr,
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
```

**Step 4: Create main.go**

Write `src/main.go`:

```go
package main

import "github.com/owner-replaceme/project-replaceme/cmd"

func main() {
	cmd.Execute()
}
```

**Step 5: Install Cobra dependency**

Run from `src/`:
```bash
cd src
go get github.com/spf13/cobra@latest
go mod tidy
cd ..
```

**Step 6: Verify it compiles**

Run:
```bash
cd src && go build -o /dev/null . && cd ..
```
Expected: success, no output

**Step 7: Commit**

```bash
git add src/
git commit -m "Add Cobra CLI scaffolding with version subcommand"
```

---

### Task 4: Create the Makefile

**Files:**
- Create: `Makefile`

**Step 1: Write the Makefile**

The Makefile must:
- Detect version from git state (tag > hash > dirty)
- Inject version via ldflags into `internal/version.Version`
- Run all Go commands from the `src/` directory
- Output binary to `_build/project-replaceme`

```makefile
VERSION_PKG := github.com/owner-replaceme/project-replaceme/internal/version

GIT_DIRTY := $(shell git diff --quiet 2>/dev/null && echo clean || echo dirty)
GIT_HASH  := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
GIT_TAG   := $(shell git describe --tags --exact-match HEAD 2>/dev/null)

ifeq ($(GIT_DIRTY),clean)
  ifneq ($(GIT_TAG),)
    VERSION := $(GIT_TAG)
  else
    VERSION := $(GIT_HASH)
  endif
else
  VERSION := $(GIT_HASH)-dirty
endif

LDFLAGS := -X $(VERSION_PKG).Version=$(VERSION)
BINARY  := project-replaceme
BUILD_DIR := _build

.PHONY: build check clean compile setup test

setup:
	@if git rev-parse --git-dir >/dev/null 2>&1; then \
		current=$$(git config core.hooksPath 2>/dev/null); \
		if [ "$$current" != ".githooks" ]; then \
			git config core.hooksPath .githooks; \
			echo "Git hooks configured (.githooks/)"; \
		fi; \
	fi

check:
	@echo "Checking code formatting..."
	@cd src && unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "❌ Files need formatting:"; \
		echo "$$unformatted"; \
		echo ""; \
		echo "Run: cd src && gofmt -w ."; \
		exit 1; \
	fi
	@echo "✓ Formatting OK"
	@echo "Running go vet..."
	@cd src && go vet ./...
	@echo "✓ Static analysis OK"
	@echo "Running tests..."
	@cd src && go test ./...
	@echo "✓ Tests passed"

compile:
	@echo "Building $(BINARY)..."
	@mkdir -p $(BUILD_DIR)
	@cd src && go build -ldflags "$(LDFLAGS)" -o ../$(BUILD_DIR)/$(BINARY) .
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY)"

build: setup check compile

test:
	@echo "Running tests..."
	@cd src && go test ./...
	@echo "✓ Tests passed"

clean:
	rm -rf $(BUILD_DIR)
```

**Step 2: Verify make compile works**

Run:
```bash
make compile
```
Expected: Binary at `_build/project-replaceme`

**Step 3: Verify make check works**

Run:
```bash
make check
```
Expected: formatting OK, vet OK, tests passed

**Step 4: Commit**

```bash
git add Makefile
git commit -m "Add Makefile with build, check, and compile targets"
```

---

### Task 5: Create Git hooks

**Files:**
- Create: `.githooks/pre-commit`
- Create: `.githooks/prepare-commit-msg`
- Create: `.githooks/post-checkout`
- Create: `.githooks/post-merge`
- Create: `.githooks/pre-push`

**Step 1: Create pre-commit hook**

This hook has three sections: beads prefix guard, make check, beads integration.

Write `.githooks/pre-commit`:

```bash
#!/usr/bin/env bash
set -euo pipefail

# === BEADS PREFIX GUARD ===
# Blocks commits until the beads prefix has been set to a real value.
# Remove this section after running /initialize-repo.
if grep -q "REPLACE-ME-WITH-YOUR-PROJECT-PREFIX" .beads/config.yaml 2>/dev/null; then
    echo "❌ Beads prefix has not been configured."
    echo ""
    echo "Run /initialize-repo to set up this project, or manually edit"
    echo ".beads/config.yaml and set issue-prefix to your project name."
    exit 1
fi
# === END BEADS PREFIX GUARD ===

# Clear git environment variables that leak into child processes.
unset GIT_DIR GIT_INDEX_FILE GIT_WORK_TREE GIT_OBJECT_DIRECTORY GIT_ALTERNATE_OBJECT_DIRECTORIES

# Run quality checks (formatting, vet, tests)
make check

# --- BEGIN BEADS INTEGRATION v0.61.0 ---
# This section is managed by beads. Do not remove these markers.
if command -v bd >/dev/null 2>&1; then
  export BD_GIT_HOOK=1
  _bd_timeout=${BEADS_HOOK_TIMEOUT:-30}
  if command -v timeout >/dev/null 2>&1; then
    timeout "$_bd_timeout" bd hooks run pre-commit "$@"
    _bd_exit=$?
    if [ $_bd_exit -eq 124 ]; then
      echo >&2 "beads: hook 'pre-commit' timed out after ${_bd_timeout}s — continuing without beads"
      _bd_exit=0
    fi
  else
    bd hooks run pre-commit "$@"
    _bd_exit=$?
  fi
  if [ $_bd_exit -eq 3 ]; then
    echo >&2 "beads: database not initialized — skipping hook 'pre-commit'"
    _bd_exit=0
  fi
  if [ $_bd_exit -ne 0 ]; then exit $_bd_exit; fi
fi
# --- END BEADS INTEGRATION v0.61.0 ---
```

**Step 2: Create the four beads-only hooks**

Each of prepare-commit-msg, post-checkout, post-merge, pre-push follows the identical beads integration pattern from the agenc repo. Only the hook name changes in each script.

Write `.githooks/prepare-commit-msg`, `.githooks/post-checkout`, `.githooks/post-merge`, `.githooks/pre-push` using the beads integration template — substituting the hook name in each.

**Step 3: Make all hooks executable**

Run:
```bash
chmod +x .githooks/*
```

**Step 4: Commit**

```bash
git add .githooks/
git commit -m "Add git hooks with beads integration and prefix guard"
```

---

### Task 6: Create GoReleaser config

**Files:**
- Create: `.goreleaser.yaml`

**Step 1: Write .goreleaser.yaml**

Pattern after agenc's config, adjusted for `src/` directory:

```yaml
version: 2

project_name: project-replaceme

before:
  hooks:
    - go mod tidy

builds:
  - id: project-replaceme
    dir: src
    main: ./main.go
    binary: project-replaceme
    env:
      - CGO_ENABLED=0
    goos:
      - darwin
      - linux
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X github.com/owner-replaceme/project-replaceme/internal/version.Version={{.Version}}

archives:
  - id: project-replaceme
    formats:
      - tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: "checksums.txt"
  algorithm: sha256

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^ci:"
      - "^chore:"

brews:
  - name: project-replaceme
    repository:
      owner: owner-replaceme
      name: homebrew-project-replaceme
      token: "{{ .Env.HOMEBREW_TAP_TOKEN }}"
    directory: Formula
    homepage: "https://github.com/owner-replaceme/project-replaceme"
    description: "project-replaceme CLI"
    license: ""
    install: |
      bin.install "project-replaceme"
    test: |
      system "#{bin}/project-replaceme", "version"
```

**Step 2: Validate the config**

Run:
```bash
goreleaser check
```
Expected: no errors (warnings about missing token are OK)

**Step 3: Commit**

```bash
git add .goreleaser.yaml
git commit -m "Add GoReleaser config for Homebrew tap publishing"
```

---

### Task 7: Create GitHub Actions workflows

**Files:**
- Create: `.github/workflows/ci.yml`
- Create: `.github/workflows/release.yml`

**Step 1: Write ci.yml**

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version-file: src/go.mod

      - name: Run checks
        run: make check
```

**Step 2: Write release.yml**

```yaml
name: Release

on:
  push:
    tags:
      - "v*"

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version-file: src/go.mod

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
```

**Step 3: Commit**

```bash
git add .github/
git commit -m "Add CI and release GitHub Actions workflows"
```

---

### Task 8: Create Claude Code settings.json

**Files:**
- Create: `.claude/settings.json`

**Step 1: Write settings.json**

```json
{
    "permissions": {
        "allow": [
            "Read(./**)",
            "Edit(./**)",
            "Write(./**)",
            "Glob(./**)",
            "Grep(./**)",
            "Bash(go:*)",
            "Bash(make:*)",
            "Bash(gofmt:*)",
            "Bash(golangci-lint:*)",
            "Bash(bd:*)",
            "Bash(goreleaser:*)"
        ],
        "deny": [
            "Read(./.env)",
            "Read(./.env.*)",
            "Read(./secrets/**)",
            "Bash(rm -rf:*)",
            "Bash(sudo:*)",
            "Bash(git commit --no-verify:*)"
        ]
    },
    "sandbox": {
        "filesystem": {
            "allowWrite": ["//tmp", "//private/tmp/claude"]
        },
        "excludedCommands": ["bd", "dolt"],
        "network": {
            "allowLocalBinding": true,
            "allowedDomains": ["127.0.0.1", "localhost"],
            "allowUnixSockets": ["//tmp/claude", "//private/tmp/claude"]
        }
    },
    "hooks": {
        "SessionStart": [
            {
                "matcher": "",
                "hooks": [
                    {
                        "type": "command",
                        "command": "bd prime"
                    }
                ]
            }
        ],
        "PreCompact": [
            {
                "matcher": "",
                "hooks": [
                    {
                        "type": "command",
                        "command": "bd prime"
                    }
                ]
            }
        ]
    }
}
```

**Step 2: Commit**

```bash
git add .claude/settings.json
git commit -m "Add Claude Code settings with Beads and Go permissions"
```

---

### Task 9: Initialize Beads with broken prefix

**Files:**
- Create: `.beads/config.yaml` (via bd init, then modify)

**Step 1: Initialize beads**

Run:
```bash
bd init --shared-server
```

**Step 2: Replace the prefix with the broken placeholder**

Edit `.beads/config.yaml` and set:
```yaml
issue-prefix: "REPLACE-ME-WITH-YOUR-PROJECT-PREFIX"
```

**Step 3: Commit**

```bash
git add .beads/
git commit -m "Add beads config with deliberately broken prefix"
```

---

### Task 10: Create CLAUDE.md

**Files:**
- Create: `CLAUDE.md`

**Step 1: Write CLAUDE.md using /prompt-engineer skill**

Invoke `/prompt-engineer` to write the CLAUDE.md. It must contain:

- Project structure: Go code in `src/`, build output in `_build/`, module path
- Building: `make build` (full pipeline), `make compile` (just binary), `make check` (lint/vet/test)
- Running: `_build/project-replaceme`
- Skill guidance: use `/technical-writer` for README changes, `/brainstorm` for new features, `/releasing` for version tagging
- Beads: project uses `bd` for issue tracking, shared server mode

**Step 2: Commit**

```bash
git add CLAUDE.md
git commit -m "Add CLAUDE.md with project instructions"
```

---

### Task 11: Create the /releasing skill

**Files:**
- Create: `.claude/skills/releasing/SKILL.md`

**Step 1: Write the releasing skill using /prompt-engineer**

Invoke `/prompt-engineer`. Pattern after agenc's `tagging-a-new-release` skill with these key elements:

- Frontmatter: name `releasing`, description about tagging releases
- Pre-flight checks (clean tree, default branch, synced)
- Current version from git tags
- Changelog from commits since last tag (features/fixes/internal)
- Semver bump rules (major/minor/patch, pre-1.0 rules)
- **Confirmation gate**: present changelog + proposed version, STOP and wait for user to explicitly confirm. User decides the version.
- Tag and push
- Monitor release workflow via `gh run list` / `gh run watch`
- Confirm success with version, GitHub release link, brew upgrade reminder

**Step 2: Commit**

```bash
git add .claude/skills/releasing/
git commit -m "Add /releasing skill for semver tagging workflow"
```

---

### Task 12: Create the SETUP_CHECKLIST.md

**Files:**
- Create: `SETUP_CHECKLIST.md`

**Step 1: Write the checklist**

This is consumed by the `/initialize-repo` skill, not read directly by humans. Markdown checklist format:

```markdown
Setup Checklist
===============

- [ ] Get GitHub owner and project name from the user
- [ ] Find-and-replace `owner-replaceme` with actual GitHub owner across all files
- [ ] Find-and-replace `project-replaceme` with actual project name across all files
- [ ] Set beads prefix in `.beads/config.yaml` to the project name
- [ ] Remove the beads prefix guard section from `.githooks/pre-commit`
- [ ] Choose and add a LICENSE file
- [ ] Create Homebrew tap repo: `gh repo create <owner>/homebrew-<project> --public`
- [ ] Enable tag protection on code repo for `v*` pattern
- [ ] 🚨 USER ACTION: Create a fine-grained PAT with write access to the tap repo, add as `HOMEBREW_TAP_TOKEN` secret in the code repo's GitHub settings
- [ ] Verify: `make build` passes
- [ ] Verify: `goreleaser check` passes
- [ ] Delete this file (SETUP_CHECKLIST.md)
- [ ] Delete `.claude/skills/initialize-repo/`
- [ ] Commit all changes
```

**Step 2: Commit**

```bash
git add SETUP_CHECKLIST.md
git commit -m "Add setup checklist for post-template initialization"
```

---

### Task 13: Create the /initialize-repo skill

**Files:**
- Create: `.claude/skills/initialize-repo/SKILL.md`

**Step 1: Write the initialize-repo skill using /prompt-engineer**

Invoke `/prompt-engineer`. The skill must:

- Frontmatter: name `initialize-repo`, description about transforming the template into a real project
- Read SETUP_CHECKLIST.md and execute each item in order
- Ask the user for GitHub owner and project name upfront
- Perform global find-and-replace for `owner-replaceme` and `project-replaceme`
- Fix the beads prefix and remove the guard from pre-commit
- Ask user to choose a license
- Create the Homebrew tap repo via `gh`
- Enable tag protection via GitHub API
- Clearly surface the HOMEBREW_TAP_TOKEN user action with `🚨 ACTION REQUIRED 🚨` formatting
- Verify build and goreleaser
- Self-destruct: delete SETUP_CHECKLIST.md and the initialize-repo skill directory
- Commit everything

**Step 2: Commit**

```bash
git add .claude/skills/initialize-repo/
git commit -m "Add /initialize-repo skill for template transformation"
```

---

### Task 14: Create the README.md

**Files:**
- Create: `README.md`

**Step 1: Write README using /technical-writer skill**

Invoke `/technical-writer`. The README must contain:

- Template usage callout at the top:
  > **This is a template repo.** To create a new CLI project from it:
  > 1. Click "Use this template" on GitHub
  > 2. Clone your new repo and run Claude Code in it
  > 3. Tell the agent: `/initialize-repo`

- Installation section (Homebrew):
  ```
  brew install owner-replaceme/project-replaceme/project-replaceme
  ```

- Upgrading section:
  ```
  brew upgrade project-replaceme
  ```

- Development section:
  - Clone the repo
  - `make build` — full build pipeline
  - `make check` — run linting and tests
  - `_build/project-replaceme` — run the binary

**Step 2: Commit**

```bash
git add README.md
git commit -m "Add README with template usage, install, and dev instructions"
```

---

### Task 15: Final verification

**Step 1: Verify the full build works**

Run:
```bash
make build
```
Expected: hooks configured, formatting OK, vet OK, tests passed, binary built

**Step 2: Verify the binary runs**

Run:
```bash
_build/project-replaceme version
```
Expected: prints a version string (git hash or "dev")

**Step 3: Verify goreleaser config**

Run:
```bash
goreleaser check
```
Expected: no errors

**Step 4: Verify the repo is clean**

Run:
```bash
git status
```
Expected: clean working tree

**Step 5: Commit any remaining changes if needed**

If any files were modified during verification, commit them.
