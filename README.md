> # This is a template for Go CLIs
> Clone it and get a working project with:
>
> - **Claude-guided onboarding** — run `/initialize-repo` and an agent sets up the repo, creates your Homebrew tap, and configures GitHub
> - **CI and Homebrew publishing** — push a version tag and GitHub Actions builds binaries and publishes to your Homebrew tap automatically
> - **Issue tracking with Beads** — `bd` integration baked into git hooks and Claude config, backed by a shared Dolt server
> - **Git hooks enforcing quality gates** — formatting, vet, and tests run on every commit structurally, not by convention
> - **Make automation** — build, check, test, run, and compile targets ready to go
> - **Automatic version injection** — the binary reports its real version from git tags via ldflags
> - **Claude skills for the dev lifecycle** — `/releasing` for semver tagging, `/brainstorm` before new features, `/technical-writer` for docs
>
> To get started:
> 1. Click **"Use this template"** on GitHub to create your repo
> 2. Clone your new repo and open it in [Claude Code](https://claude.com/claude-code)
> 3. Run `/initialize-repo` to configure the project

project-replaceme
=================

A command-line tool built with Go and [Cobra](https://github.com/spf13/cobra).

Installation
------------

Install via [Homebrew](https://brew.sh):

```bash
brew install owner-replaceme/project-replaceme/project-replaceme
```

This installs the `project-replaceme` binary and makes it available on your PATH.

Verify the installation:

```bash
project-replaceme version
```

Upgrading
---------

```bash
brew upgrade project-replaceme
```

Development
-----------

### Prerequisites

- [Go](https://go.dev/dl/) (version specified in `src/go.mod`)
- [GNU Make](https://www.gnu.org/software/make/)

### Getting started

Clone the repository and run the full build pipeline:

```bash
git clone https://github.com/owner-replaceme/project-replaceme.git
cd project-replaceme
make build
```

`make build` configures Git hooks, runs formatting and linting checks, runs the test suite, and compiles the binary. The output binary is at `_build/project-replaceme`.

### Common tasks

| Command | What it does |
|---|---|
| `make build` | Full pipeline: configure hooks, check, compile |
| `make compile` | Build the binary only (skip checks) |
| `make check` | Run formatting, static analysis, and tests |
| `make test` | Run the test suite only |
| `make run` | Build and run the binary (pass args via `ARGS="..."`) |
| `make clean` | Remove build artifacts |

### Running the binary

Use `make run` to build and run in one step:

```bash
make run
make run ARGS="version"
make run ARGS="--help"
```

Or run the compiled binary directly:

```bash
_build/project-replaceme
_build/project-replaceme version
```

### Project layout

```
src/                    Go source code
  cmd/                  CLI command definitions (Cobra)
  internal/             Internal packages
  main.go              Entry point
_build/                 Compiled binary (gitignored)
.githooks/              Git hooks (activated by make setup)
```
