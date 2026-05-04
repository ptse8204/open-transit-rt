# Device And AVL Evidence

This folder is a template-only scaffold for future device, GPS emitter, or vendor AVL integration review.

It does not contain real telemetry evidence, vendor approvals, hardware certifications, production AVL proof, consumer acceptance, CAL-ITP/Caltrans compliance evidence, or agency endorsement.

## Contents

- `templates/integration-review-template.md`: reusable review template for future public-safe integration evidence.
- Phase 29B synthetic adapter fixtures live under `testdata/avl-vendor/`; they are test fixtures, not real device or vendor evidence.

## Evidence Rules

Do not add artifacts here until they have passed source, permission, and redaction review.

Do not commit:

- device tokens;
- admin tokens;
- JWT or CSRF secrets;
- DB URLs with passwords;
- vendor credentials;
- private AVL payloads;
- raw private telemetry;
- private device, vehicle, or vendor identifiers;
- private operator notes;
- private logs with credentials;
- `.cache` files.

Future evidence must clearly label whether it comes from:

- simulator or no-hardware testing;
- a pilot device;
- a real production-directed device;
- vendor-owned middleware;
- agency-owned adapter code.

Every evidence record must state what it proves and what it does not prove. Simulator evidence does not prove production AVL reliability. Validator success does not prove consumer acceptance.

Phase 29B dry-run adapter output is transform-review evidence only when recorded in a future reviewed packet. It is not submitted telemetry, production integration evidence, successful vendor compatibility proof, certified hardware support, or production AVL reliability proof.
