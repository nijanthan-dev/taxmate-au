#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
export GOCACHE="${GOCACHE:-$ROOT/.cache/go-build}"

fail() {
  echo "error: $*" >&2
  exit 1
}

[[ -f .codex-plugin/plugin.json ]] || fail "missing plugin manifest"
[[ -f README.md ]] || fail "missing README"
[[ -f DISCLAIMER.md ]] || fail "missing DISCLAIMER.md"
[[ -f LICENSE ]] || fail "missing LICENSE"
[[ -f SECURITY.md ]] || fail "missing SECURITY.md"
[[ -f CONTRIBUTING.md ]] || fail "missing CONTRIBUTING.md"
[[ -f docs/PUBLICATION_CHECKLIST.md ]] || fail "missing publication checklist"

if git ls-files 'bin/*' | grep -q .; then
  fail "built binaries are tracked"
fi

if git grep -nE 'public[-]work|taxmate-au-public[-]work' -- . ':!data/ato_knowledge_base/raw/**' ':!data/ato_knowledge_base/text/**'; then
  fail "temporary staging name leaked"
fi

if git grep -nE '/Users/[[:alnum:]_.-]+|custom[_]apps/skills[_]and[_]plugins|Developer/custom[_]apps' -- README.md .codex-plugin agents skills docs; then
  fail "private machine path leaked into public docs"
fi

git grep -Eq 'not (professional )?tax, legal, accounting, financial' -- README.md DISCLAIMER.md .codex-plugin skills || fail "missing professional-advice disclaimer"
git grep -q 'not affiliated with' -- README.md DISCLAIMER.md .codex-plugin skills || fail "missing affiliation disclaimer"
git grep -q 'does not lodge' -- DISCLAIMER.md || fail "missing lodgment disclaimer"
git grep -q 'Accountant review' -- DISCLAIMER.md skills || fail "missing accountant-review boundary"

if git grep -nE 'taxmate-au-re[d]act|internal/pri[v]acy|cmd/taxmate-au-re[d]act|RE[D]ACTED' -- . ':!data/ato_knowledge_base/raw/**' ':!data/ato_knowledge_base/text/**'; then
  fail "legacy file-sanitisation artifact found"
fi

go test ./...
mkdir -p bin
go build -o bin/taxmate-au-refresh ./cmd/taxmate-au-refresh
go build -o bin/taxmate-au-validate ./cmd/taxmate-au-validate
go build -o bin/taxmate-au-finance ./cmd/taxmate-au-finance
go build -o bin/taxmate-au-calc ./cmd/taxmate-au-calc
bin/taxmate-au-validate >/tmp/taxmate-au-validate.json

echo "publication checks passed"
