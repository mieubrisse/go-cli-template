---
name: initialize-repo
description: >-
  Transforms this template repo into a real project. Invoke after cloning a repo
  created from the Go CLI template. Replaces placeholder names, configures Beads,
  creates the Homebrew tap repo, and sets up GitHub release infrastructure.
  Self-destructs after completion.
---

Initialize Repository
=====================

This skill transforms the template into a working project. Read `SETUP_CHECKLIST.md` and execute each step in order, checking items off as you go.

Workflow
--------

### 1. Gather project details

Ask the user for:
- **GitHub owner** (username or org) — replaces `owner-replaceme`
- **Project name** — replaces `project-replaceme` (used as binary name, repo name, and Homebrew formula name)

### 2. Replace placeholders

Perform a global find-and-replace across all files in the repo:

1. Replace `owner-replaceme` with the actual GitHub owner
2. Replace `project-replaceme` with the actual project name

Files affected include: `go.mod`, `go.sum`, all `.go` files, `Makefile`, `.goreleaser.yaml`, `CLAUDE.md`, `README.md`, GitHub Actions workflows, and the releasing skill.

After replacement, verify imports resolve:

```bash
cd src && go build ./... && cd ..
```

### 3. Configure Beads

Edit `.beads/config.yaml`: set `issue-prefix` to the project name (replacing `REPLACE-ME-WITH-YOUR-PROJECT-PREFIX`).

Remove the beads prefix guard from `.githooks/pre-commit` — delete the entire block between `# === BEADS PREFIX GUARD ===` and `# === END BEADS PREFIX GUARD ===` (inclusive).

### 4. Choose a license

Ask the user which license they want (MIT, Apache-2.0, AGPL-3.0, etc.). Create the `LICENSE` file with the full license text, filling in the current year and the user's name/org.

Update the `license` field in `.goreleaser.yaml` to match.

### 5. Create the Homebrew tap repo

```bash
gh repo create <owner>/homebrew-<project> --public --description "Homebrew tap for <project>"
```

### 6. Enable tag protection

Enable tag protection for `v*` tags on the code repo:

```bash
gh api repos/<owner>/<project>/rulesets --method POST --input - <<'EOF'
{
  "name": "Protect release tags",
  "target": "tag",
  "enforcement": "active",
  "conditions": {
    "ref_name": {
      "include": ["refs/tags/v*"],
      "exclude": []
    }
  },
  "rules": [
    {
      "type": "creation"
    },
    {
      "type": "deletion"
    }
  ]
}
EOF
```

### 7. User action — HOMEBREW_TAP_TOKEN

Surface this clearly:

🚨 **ACTION REQUIRED** 🚨

Create a fine-grained Personal Access Token with **write access to the `<owner>/homebrew-<project>` repository**, then add it as a secret named `HOMEBREW_TAP_TOKEN` in your code repo's GitHub settings:

**Settings → Secrets and variables → Actions → New repository secret**

Wait for the user to confirm they have completed this step before proceeding.

### 8. Verify

```bash
make build
goreleaser check
```

Both must pass. Fix any issues before proceeding.

### 9. Clean up

Delete the following files and directories:
- `SETUP_CHECKLIST.md`
- `.claude/skills/initialize-repo/` (this skill)

Remove the template usage callout from the top of `README.md` (the blockquote section about "This is a template repo").

### 10. Commit

Stage and commit all changes with a message like:

```
Initialize project as <project-name>
```

Push to the remote.
