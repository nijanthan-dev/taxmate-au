# Full Plugin Runtime

Use this advanced path only when you need live ATO source refresh, CSV finance review, calculator commands, skill regeneration, source coverage, and Codex plugin orchestration.

Prerequisites:

- Node.js 18+ if also testing portable install.
- Bash 5+.
- Python 3.9+.
- Git.

## Clean checkout setup

```bash
git clone https://github.com/nijanthan-dev/taxmate-australia.git
cd taxmate-australia
```

Run full-runtime commands through the launcher (bash entrypoint + python runtime):

```bash
./scripts/taxmate refresh --help
```

The same commands also work by calling the Python module directly:

```bash
./scripts/taxmate.py refresh --help
```

Validate:

```bash
./scripts/taxmate validate
./scripts/taxmate skills validate
./scripts/taxmate skills audit --check
./scripts/taxmate skills audit --format markdown --output /tmp/source-coverage.md
scripts/check-publication-ready.sh
```

If you need local-speed, rebuild native binaries in `bin/` once and keep those ahead of the launcher.

## Local plugin setup

This repo includes `.codex-plugin/plugin.json` for advanced local plugin testing. Local marketplace configuration is development-only.

If you create a user-global marketplace file, its path is:

```text
~/.agents/plugins/marketplace.json
```

For a cloned repo at `/absolute/path/taxmate-australia`, the local plugin entry path should be that absolute repo path. The path is interpreted relative to the marketplace file only when it is relative.

Example:

```json
{
  "name": "taxmate-local",
  "interface": { "displayName": "Local plugins" },
  "plugins": [
    {
      "name": "taxmate-australia",
      "source": {
        "source": "local",
        "path": "/absolute/path/taxmate-australia"
      },
      "policy": {
        "installation": "AVAILABLE",
        "authentication": "ON_INSTALL"
      },
      "category": "Productivity"
    }
  ]
}
```

Do not claim official plugin discovery unless a published listing has been verified.
