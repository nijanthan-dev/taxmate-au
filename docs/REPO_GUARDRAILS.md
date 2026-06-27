# Repository Guardrails

This repo should be easy to contribute to, but hard to accidentally weaken. These guardrails are based on a live benchmark of mature OSS repos and Claude/Codex skill repos on 2026-06-27.

## Benchmark

| Repo | Relevant practice | TaxMate adoption |
| --- | --- | --- |
| [alirezarezvani/claude-skills](https://github.com/alirezarezvani/claude-skills) | Contributor docs, skill authoring standard, PR template, branch-protection config snapshot. | Keep a source-backed guardrail doc, PR template, and local process check. |
| [Jeffallan/claude-skills](https://github.com/Jeffallan/claude-skills) | Skill-specific issue templates for new skills and Claude issues. | Keep issue forms that ask for scope, source basis, and safety handling. |
| [SawyerHood/dev-browser](https://github.com/SawyerHood/dev-browser) | Contributing and release docs for a focused Claude skill. | Keep release and validation expectations explicit before merge. |
| [posit-dev/skills](https://github.com/posit-dev/skills) | Contributor guide for a curated skill collection. | Separate contribution workflow from runtime docs. |
| [Dimillian/Skills](https://github.com/Dimillian/Skills) | Codex skill collection with examples organized by skill. | Keep public skills discoverable and generated artifacts checked. |
| [ComposioHQ/awesome-codex-skills](https://github.com/ComposioHQ/awesome-codex-skills) | Large Codex skill catalog with simple per-skill layout. | Keep portable skill surfaces simple, but add stricter TaxMate safety gates. |
| [kubernetes/kubernetes](https://github.com/kubernetes/kubernetes) | PR template and ownership metadata. | Keep CODEOWNERS and required PR checklist. |
| [rust-lang/rust](https://github.com/rust-lang/rust) | Issue templates, PR template, workflows, and dependency automation. | Keep templates/workflows; skip dependency automation until manifests exist. |
| [microsoft/vscode](https://github.com/microsoft/vscode) | CODEOWNERS, issue templates, workflow checks, Dependabot, and agent guidance. | Keep CODEOWNERS and agent instructions. Add dependency automation only with committed manifests. |
| [Homebrew/brew](https://github.com/Homebrew/brew) | PR template, issue templates, CodeQL, action linting, dependency/license checks. | Keep security scan and publication checks; add action linting later if workflows grow. |
| [nodejs/node](https://github.com/nodejs/node) | CODEOWNERS, support docs, PR template, labels, CodeQL, and many workflows. | Keep support/security boundaries and required status checks clear. |

## Current GitHub Protection Target

Verified on 2026-06-27:

- Default branch: `main`
- Merge policy: squash only; merge commits and rebase merges disabled
- Delete branch on merge: enabled
- Required checks: `Validate`, `Plugin Package`, `Gitleaks`
- Required status checks strict: enabled
- Conversation resolution: required
- Linear history: required
- Admin enforcement: enabled
- Force pushes: disabled
- Branch deletion: disabled

Known gap:

- Required approving reviews are not enabled. Do not treat `@Codex` as a privileged approver. Treat it as a review signal, then merge only when the owner is satisfied and all required checks/threads are clean.

Check live state before merge:

```bash
gh repo view nijanthan-dev/taxmate-australia --json defaultBranchRef,mergeCommitAllowed,rebaseMergeAllowed,squashMergeAllowed,deleteBranchOnMerge
gh api repos/nijanthan-dev/taxmate-australia/branches/main/protection
```

## Required Repo Files

- `CONTRIBUTING.md`
- `SECURITY.md`
- `CODE_OF_CONDUCT.md`
- `.github/CODEOWNERS`
- `.github/pull_request_template.md`
- `.github/ISSUE_TEMPLATE/bug_report.yml`
- `.github/ISSUE_TEMPLATE/feature_request.yml`
- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `scripts/check-repo-guardrails.sh`
- `scripts/check-publication-ready.sh`

## PR Gate

Local:

```bash
PYTHONPYCACHEPREFIX=/tmp/taxmate-pycache python3 -m py_compile scripts/*.py
./scripts/taxmate validate
./scripts/taxmate skills generate --check
./scripts/taxmate skills audit --check
scripts/check-repo-guardrails.sh
scripts/check-publication-ready.sh
gitleaks detect --source . --redact
gitleaks dir . --redact
```

Remote:

- `Validate` green
- `Plugin Package` green
- `Gitleaks` green
- `mergeStateStatus` is `CLEAN`
- No unresolved review threads
- Latest-head `@Codex` review has no blocking findings

## Maintainer Rules

- Never merge with unresolved review threads.
- Never merge without a current local or CI secret scan.
- If a secret reaches remote Git, stop and remove it from branch history before merge.
- Reply to fixed review threads and resolve them after verification.
- Use squash merge only.
