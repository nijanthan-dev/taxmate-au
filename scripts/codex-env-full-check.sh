#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

bash scripts/codex-env-setup.sh
bash scripts/check-publication-ready.sh
gitleaks detect --source . --redact --no-banner
