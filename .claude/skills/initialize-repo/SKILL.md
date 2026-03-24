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

This skill transforms the template into a working project.

Read `SETUP_CHECKLIST.md` and execute each item in order, checking them off as you go. The checklist is the single source of truth for all initialization steps.

For items that require user input (GitHub owner, project name, license choice, HOMEBREW_TAP_TOKEN), ask the user and wait for their response before proceeding.

For the `🚨 USER ACTION` item, surface it clearly and wait for the user to confirm completion before continuing.
