Setup Checklist
===============

This checklist is consumed by the `/initialize-repo` skill. Run that skill instead of following this manually.

- [ ] Get GitHub owner and project name from the user
- [ ] Find-and-replace `owner-replaceme` with actual GitHub owner across all files
- [ ] Find-and-replace `project-replaceme` with actual project name across all files
- [ ] Verify Go code compiles after replacement: `cd src && go build ./...`
- [ ] Set beads prefix in `.beads/config.yaml` to the project name
- [ ] Remove the beads prefix guard section from `.githooks/pre-commit`
- [ ] Choose a license and create the LICENSE file
- [ ] Update the `license` field in `.goreleaser.yaml`
- [ ] Create Homebrew tap repo: `gh repo create <owner>/homebrew-<project> --public`
- [ ] Enable tag protection on the code repo for `v*` tags
- [ ] 🚨 USER ACTION: Create a fine-grained PAT with write access to the tap repo, add as `HOMEBREW_TAP_TOKEN` secret
- [ ] Verify: `make build` passes
- [ ] Verify: `goreleaser check` passes
- [ ] Remove the template usage callout from README.md
- [ ] Delete this file (SETUP_CHECKLIST.md)
- [ ] Delete `.claude/skills/initialize-repo/`
- [ ] Commit and push all changes
