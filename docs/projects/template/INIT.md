> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This template is subordinate to it.**
>
> ⚠️ **This file is TEMPORARY. After you complete this checklist, you MUST delete it. Do not leave it in the project folder.**
>
> Run this checklist ONCE when creating a new project from the template.
> After initialization, delete this file so it does not clutter the workspace.

---

## Step 0: Human Pre-Flight (do not skip)
- [ ] Read `PREREQUISITES.md` and ensure all Required skills and MCPs are installed.
- [ ] Complete the Environment & Access Checklist in `PREREQUISITES.md`.
- [ ] Define Human Approval Gates and notify approvers.
- [ ] Sign off `PREREQUISITES.md` before allowing AI execution.

## Step 1: Identity
- [ ] Rename this folder from `template` to `<project-kebab-case-name>`.
- [ ] In `README.md`, replace `<PROJECT_NAME>` with the real project name.
- [ ] In `README.md`, fill `Started`, `Target End`, `Human Reviewer`.
- [ ] In `README.md`, fill `Target Repo` with `<repo-root>`.
- [ ] In `README.md`, confirm `Product Branch: dev`.

## Step 2: Configure Tracking
- [ ] In `TRACKING.md`, replace `<Project Name>` with the real name.
- [ ] In `TRACKING.md`, fill the `Scope` line.
- [ ] In `TRACKING.md`, delete example rows and create real prompt inventory rows.

## Step 3: Configure Metadata
- [ ] In `DECISIONS.md`, replace `<Project Name>` with the real name.
- [ ] In `RISKS.md`, replace `<Project Name>` with the real name.
- [ ] In `HANDOFF.md`, replace `<Project Name>` with the real name.

## Step 4: Prepare Prompts
- [ ] Delete `prompts/_example.md` (it is instructional only).
- [ ] Keep `prompts/_template.md` as the authoring standard.
- [ ] Create real prompt files: `01-*.md`, `02-*.md`, etc.
- [ ] (Optional) Create real task files in `tasks/` if using dual-layer.

## Step 5: Configure Project Metadata
- [ ] In `KNOWN_ISSUES.md`, replace `<Project Name>` and fill the pre-existing bugs / out-of-scope areas.
- [ ] In `REFERENCES.md`, replace `<Project Name>` and add project-specific links.
- [ ] In `PREREQUISITES.md`, replace `<Project Name>` and fill all checklists.

## Step 6: Verify Baseline
- [ ] Run `git branch --show-current` and confirm `dev`.
- [ ] Run `task check` and confirm it passes.
- [ ] Run `task smoke` and confirm it passes.

## Step 7: Commit Planning Docs
- [ ] `git add docs/projects/<project-name>/`
- [ ] `git commit -m "docs(repo): init <project-name> planning workspace"`
- [ ] Push to `origin/dev` (planning docs only — no product code).

## Step 8: Cleanup
- [ ] Delete this `INIT.md` file.
- [ ] Update `TRACKING.md` status to 🟡 In Progress.
- [ ] Start execution with Step 01 on `dev`.
