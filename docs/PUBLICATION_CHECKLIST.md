# Publication Checklist

Use this before publishing TaxMate Australia outside a local install.

- Confirm legal license and repository URL.
- Confirm `DISCLAIMER.md` is present and linked from README.
- Confirm no wording implies ATO, Commonwealth, state revenue office, insurer, super fund, or financial-institution endorsement.
- Run `go test ./...`.
- Run `bin/taxmate-au-validate` and require `100/100`.
- Refresh or recrawl ATO sources near release date.
- Confirm source pack contains only official ATO pages and expected state revenue routing notes.
- Confirm no private paths in publication-facing docs or plugin skills.
- Confirm compatibility wrappers are clearly marked as local install helpers.
- Confirm no legacy file-sanitisation code or binaries.
- Confirm workbook and taxpack skills cannot make independent tax treatment calls.
- Confirm every high-risk area defaults to `Accountant review` when facts are incomplete.
