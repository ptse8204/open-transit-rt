# Consumer Submission Artifacts

This directory is for redacted target-originated evidence after an operator
performs real-world adoption steps outside the repo.

The target directories intentionally contain README files only until real
redacted evidence exists. Do not add placeholder screenshots, fake receipts,
fake tickets, fake emails, example correspondence, or generated portal images.

## Target Directories

- `google-maps/`
- `apple-maps/`
- `transit-app/`
- `bing-maps/`
- `moovit/`
- `mobility-database/`
- `transit-land/`

## What May Be Committed

Only commit artifacts after redaction review. Public-safe examples include:

- redacted submission receipt;
- redacted ticket confirmation;
- redacted target email;
- redacted portal screenshot showing only the relevant status;
- redacted rejection or change-request reason;
- redacted acceptance confirmation;
- public official submission-path source note.

## What Must Stay Private

Do not commit:

- portal credentials;
- private ticket links;
- raw private correspondence;
- account names or personal contact details that are not public-safe;
- screenshots containing personal data or private portal navigation;
- tokens, DB URLs, private keys, admin URLs with secrets, or session IDs;
- raw private operator artifacts.

## Naming

Use:

```text
YYYY-MM-DD_<artifact-type>_<short-description>.<ext>
```

Examples of artifact types: `receipt`, `ticket`, `email`, `portal-screenshot`,
`rejection`, `blocker`, `acceptance`, `official-path`.

## Required Updates

After adding real evidence, update:

- the target current record under `../current/`;
- `../README.md`;
- `../status.json`;
- `docs/current-status.md` and `docs/handoffs/latest.md` if the next action
  changed.

The human-readable tracker and `status.json` must agree for target name,
status, packet path, prepared timestamp, and evidence references.
