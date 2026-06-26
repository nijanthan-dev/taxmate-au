#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="${TAXMATE_AUSTRALIA_ROOT:-}"
if [[ -z "$ROOT" ]]; then
  ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
fi

case "$ROOT" in
  /*) ;;
  *) echo "error: root must be absolute" >&2; exit 1 ;;
esac

cd "$ROOT"

REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null)" || {
  echo "error: refusing cleanup outside git worktree" >&2
  exit 1
}
cd "$REPO_ROOT"

rm -rf \
  .cache/ato \
  .tmp \
  bin \
  coverage \
  dist \
  htmlcov \
  .coverage \
  .mypy_cache \
  .pytest_cache \
  .ruff_cache

find . \
  -path './.git' -prune -o \
  -path './.venv' -prune -o \
  -type d -name '__pycache__' -exec rm -rf {} +

find . \
  -path './.git' -prune -o \
  -path './.venv' -prune -o \
  -type f \( -name '*.pyc' -o -name '*.pyo' \) -exec rm -f {} +

echo "post-task cleanup complete"
