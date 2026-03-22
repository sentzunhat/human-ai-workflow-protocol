# TASK-065 Evidence: GitHub-hosted distribution auto-sync proof

Date: 2026-05-15

## Source Query

Command:

```bash
curl -s https://api.github.com/repos/sentzunhat/human-ai-workflow-protocol/actions/workflows/sync-distribution-generated.yml/runs?per_page=5
```

## Verified Run

- Workflow: Validate Distribution Generated
- Workflow file: .github/workflows/sync-distribution-generated.yml
- Run ID: 25897172755
- Run URL: https://github.com/sentzunhat/human-ai-workflow-protocol/actions/runs/25897172755
- Branch: dev
- Event: push
- Status: completed
- Conclusion: success
- Head SHA: 0da2ffd7e7fc65b461ad273c4173c7a10edd36c1
- Created At (UTC): 2026-05-15T02:35:42Z
- Updated At (UTC): 2026-05-15T02:35:57Z

## Repo-root Proof (redacted)

```bash
pwd
<repo-root-abs>

git rev-parse --show-toplevel
<repo-root-abs>

git rev-parse --show-prefix


git status --short
 M .hawp/bin/hawp
 M .hawp/kit/guidance/da-schema-planning.md
 M .hawp/work/BACKLOG.md
RM .hawp/work/active/TASK-013.md -> .hawp/work/closed/2026/05/15/TASK-013.md
 M README.md
 M core/.hawp/kit/standards/README.md
 M librarian/package.json
?? .awp/
?? .hawp/kit/guidance/shared-standards-review-rubric.md
?? .hawp/work/active/TASK-065.md
?? .hawp/work/closed/2026/05/15/TASK-050.md
?? .hawp/work/closed/2026/05/15/TASK-051.md
?? .hawp/work/closed/2026/05/15/TASK-052.md
?? .hawp/work/closed/2026/05/15/TASK-053.md
?? .hawp/work/closed/2026/05/15/TASK-054.md
?? .hawp/work/closed/2026/05/15/TASK-055.md
?? .hawp/work/closed/2026/05/15/TASK-056.md
?? .hawp/work/closed/2026/05/15/TASK-057.md
?? .hawp/work/closed/2026/05/15/TASK-058.md
?? .hawp/work/closed/2026/05/15/TASK-059.md
?? .hawp/work/closed/2026/05/15/TASK-060.md
?? .hawp/work/closed/2026/05/15/TASK-061.md
?? .hawp/work/closed/2026/05/15/TASK-062.md
?? .hawp/work/closed/2026/05/15/TASK-063.md
?? .hawp/work/closed/2026/05/15/TASK-064.md
?? .hawp/work/evidence/2026/05/15/
?? .hawp/work/evidence/db-decision-template.md
?? core/.hawp/kit/standards/database/
?? core/.hawp/kit/standards/patterns/
?? core/.hawp/kit/standards/service-design/
?? librarian/scripts/backlog-validate/
?? shared_standards/
```
