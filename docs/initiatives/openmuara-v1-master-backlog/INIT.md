> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**
>
> ⚠️ **This file is TEMPORARY. After initialization, delete it.**
>
> Run this checklist ONCE when creating this initiative from the template.

---

## Step 0: Human Pre-Flight
- [x] Read `PREREQUISITES.md` and ensure all required skills and MCPs are installed.
- [x] Complete the Environment & Access Checklist in `PREREQUISITES.md`.
- [x] Define Human Approval Gates and notify approvers.
- [x] Sign off `PREREQUISITES.md` before allowing AI execution.

## Step 1: Identity
- [x] In `README.md`, replace `Started`, `Target End`, `Human Reviewer`.
- [x] In `README.md`, confirm `Target Repo` and `Product Branch`.

## Step 2: Configure Tracking
- [x] In `TRACKING.md`, fill the `Scope` line.
- [x] In `TRACKING.md`, review the Prompt Inventory rows and adjust as needed.

## Step 3: Configure Metadata
- [x] In `DECISIONS.md`, replace `<Initiative Name>` with `OpenMuara v1 Master Backlog`.
- [x] In `RISKS.md`, replace `<Initiative Name>` with `OpenMuara v1 Master Backlog`.
- [x] In `HANDOFF.md`, replace `<Initiative Name>` with `OpenMuara v1 Master Backlog`.

## Step 4: Prepare Prompts
- [x] Delete `prompts/_example.md` (instructional only).
- [x] Keep `prompts/_template.md` as the authoring standard.
- [x] Review/create real prompt files.
- [x] Create matching `tasks/##-*.md` files for non-trivial prompts.

## Step 5: Configure Initiative Metadata
- [x] In `KNOWN_ISSUES.md`, replace placeholders and fill pre-existing bugs / out-of-scope areas.
- [x] In `REFERENCES.md`, replace placeholders and add initiative-specific links.
- [x] In `PREREQUISITES.md`, replace placeholders and fill checklists.

## Step 6: Environment Check
- [x] From the repo root, run quality gates.
- [x] All gates passed before starting execution.

## Step 7: Commit Planning Docs
- [x] `git add docs/initiatives/openmuara-v1-master-backlog/`
- [x] `git commit -m "docs(initiatives): add openmuara-v1-master-backlog priority tracker"`
- [x] Commit made on root repo `dev` branch (planning docs only).

## Step 8: Cleanup
- [ ] Delete this `INIT.md` file at project closeout.
- [x] Update `TRACKING.md` status to 🟡 Active.
