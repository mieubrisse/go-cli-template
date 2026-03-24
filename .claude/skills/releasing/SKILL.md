---
name: releasing
description: >-
  Guides the semantic versioning release workflow. Invoke when the user asks to
  tag, release, cut a release, bump the version, or publish a new version.
  Determines the next semver version, presents it for user approval, then tags
  and pushes. Monitors the release CI workflow to completion.
---

Releasing a New Version
=======================

This project uses **semantic versioning** and **GoReleaser**. Releases are triggered by pushing a Git tag — CI builds binaries and updates the Homebrew tap automatically. Do not run GoReleaser locally.

Workflow
--------

### 1. Pre-flight checks

Verify the environment is ready. Do not proceed until all checks pass.

```bash
git branch --show-current
git status --porcelain
git fetch origin
git status -uno
```

- Must be on the default branch
- Working tree must be clean (no uncommitted changes)
- Local branch must be up to date with the remote

If dirty, commit or stash first. If behind, pull before proceeding.

### 2. Determine the current version

```bash
git fetch --tags origin
git tag --sort=-v:refname | head -5
```

Identify the latest `vX.Y.Z` tag. This is the baseline for the next version.

### 3. Analyze changes and build a changelog

```bash
git log <latest-tag>..HEAD --oneline
```

Group commits into categories:

- **Features:** new capabilities, commands, integrations
- **Fixes:** bug fixes, crash fixes, correctness improvements
- **Internal:** refactoring, docs, CI, dependency updates

Classify the overall change scope:

| Change type | Version bump | Examples |
|---|---|---|
| Breaking or incompatible changes | **Major** (X.0.0), or **Minor** while pre-1.0 | Removed a command, changed config schema incompatibly |
| New features or capabilities | **Minor** (0.X.0) | Added a subcommand, new config option |
| Bug fixes, docs, internal | **Patch** (0.0.X) | Fixed a crash, corrected a typo, refactored internals |

While the project is pre-1.0, breaking changes bump minor (not major).

### 4. Present the proposed version — then stop

Present the user with:

1. The changelog summary
2. The suggested version with reasoning
3. The commit that will be tagged (short SHA + first line)

<confirmation-gate>
**STOP.** Wait for the user to explicitly confirm or override the version number. Do not tag until the user states the version they want.

The user decides the version — your suggestion is advisory. If the user picks a different version, use theirs.
</confirmation-gate>

### 5. Tag and push

Verify the tag does not already exist:

```bash
git tag -l v<confirmed-version>
```

If it exists, ask the user for a different version. Otherwise:

```bash
git tag v<confirmed-version>
```

```bash
git push origin v<confirmed-version>
```

### 6. Verify the release

The tag push triggers `.github/workflows/release.yml`, which runs GoReleaser to build binaries and update the Homebrew tap.

Monitor until complete:

```bash
gh run list --workflow=release.yml --limit=1
```

If still in progress, stream logs:

```bash
gh run watch <run-id>
```

The release is not done until the workflow reports success. If it fails, investigate and resolve before telling the user the release is complete.

When successful, confirm with:
- The version released
- A link to the GitHub release page
- A reminder that `brew upgrade project-replaceme` will pick up the new version
