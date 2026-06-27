#!/usr/bin/env bash
set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"
cd "$ROOT"

fail() {
  echo "error: $*" >&2
  exit 1
}

require_file() {
  [[ -f "$1" ]] || fail "missing $1"
}

require_grep() {
  local pattern="$1"
  local file="$2"
  grep -Eq "$pattern" "$file" || fail "$file missing pattern: $pattern"
}

require_file CONTRIBUTING.md
require_file SECURITY.md
require_file CODE_OF_CONDUCT.md
require_file docs/REPO_GUARDRAILS.md
require_file .github/CODEOWNERS
require_file .github/pull_request_template.md
require_file .github/ISSUE_TEMPLATE/bug_report.yml
require_file .github/ISSUE_TEMPLATE/feature_request.yml
require_file .github/ISSUE_TEMPLATE/config.yml
require_file .github/workflows/ci.yml
require_file .github/workflows/release.yml
require_file docs/PUBLICATION_CHECKLIST.md
require_file scripts/check-publication-ready.sh

require_grep 'docs/REPO_GUARDRAILS.md' CONTRIBUTING.md
require_grep 'gitleaks detect --source \. --redact' CONTRIBUTING.md
require_grep 'gitleaks dir \. --redact' CONTRIBUTING.md
require_grep 'mergeStateStatus' CONTRIBUTING.md
require_grep '@Codex' CONTRIBUTING.md

require_grep 'alirezarezvani/claude-skills' docs/REPO_GUARDRAILS.md
require_grep 'ComposioHQ/awesome-codex-skills' docs/REPO_GUARDRAILS.md
require_grep 'kubernetes/kubernetes' docs/REPO_GUARDRAILS.md
require_grep 'Required checks: `Validate`, `Plugin Package`, `Gitleaks`' docs/REPO_GUARDRAILS.md
require_grep 'Required approving reviews are not enabled' docs/REPO_GUARDRAILS.md

require_grep 'scripts/check-repo-guardrails.sh' .github/pull_request_template.md
require_grep 'secret scan' .github/pull_request_template.md
require_grep 'gitleaks dir \. --redact' .github/pull_request_template.md
require_grep 'Generated Artifacts' .github/pull_request_template.md
require_grep '@Codex' .github/pull_request_template.md
require_grep 'mergeStateStatus' .github/pull_request_template.md

require_grep 'Use synthetic data only' .github/ISSUE_TEMPLATE/bug_report.yml
require_grep 'CI or repo guardrail' .github/ISSUE_TEMPLATE/bug_report.yml
require_grep 'Source basis' .github/ISSUE_TEMPLATE/feature_request.yml
require_grep 'required: true' .github/ISSUE_TEMPLATE/feature_request.yml
require_grep 'blank_issues_enabled: false' .github/ISSUE_TEMPLATE/config.yml
require_grep 'scripts/check-repo-guardrails.sh' docs/PUBLICATION_CHECKLIST.md

require_grep 'scripts/check-repo-guardrails.sh' .github/workflows/ci.yml
require_grep 'scripts/check-repo-guardrails.sh' scripts/check-publication-ready.sh

echo "repo guardrails passed"
