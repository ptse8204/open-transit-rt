# Product Language And Layout Guide

Open Transit RT should read like self-hosted transit operations software for
small agencies, operator staff, deployment owners, and civic technologists.
Public pages and primary private console pages should help people decide what
to do next. They should not read like an audit log, release ledger, or internal
implementation handoff.

## Voice

Use direct product language:

- Import a schedule.
- Check feed URLs.
- Connect vehicle data.
- Fix feed issues.
- Maintain the deployment.
- Ask an administrator or deployment owner.

Avoid internal shorthand on public pages and primary console HTML:

- phase-ledger wording;
- maintainer-only status language;
- raw internal safety flag names;
- repeated caveats on every card;
- route lists before the user sees the next action.

## Terms To Prefer

| Use | Avoid on primary pages |
| --- | --- |
| Fix feed issues | Common Next Actions |
| Tour the operator workflow | Follow the same path staff use in the console |
| When an administrator is needed | What stays technical? |
| Administrator or deployment owner | technical helper when repeated as a catch-all |
| Limits | repeated does-not-prove paragraphs |
| Feed URLs | configured public route URLs |

## Boundaries

Keep the product truthful. Local and reference deployment review still does not
prove compliance, production readiness, consumer acceptance, agency adoption,
hosted service availability, vendor compatibility, SLA coverage, production AVL
reliability, or ETA quality.

Place those boundaries where they help:

- Public pages should have one concise Limits section unless a risky action
  needs an inline warning.
- Primary console pages may show one compact limits or permission note.
- Detailed safety flags belong in JSON, tests, or collapsed advanced sections.
- Evidence, compliance, and release documents may stay more formal, but should
  not be linked as the first path for normal users.

## Layout

Use layouts that make the next action visible quickly:

- Start with an action list, short current-state summary, or focused form.
- Use cards only for distinct repeated items, not for every paragraph.
- Prefer ordered task flows, tables, split panes, accordions, and compact action
  lists over large grids.
- Keep diagnostic tables below the main action path or behind disclosure.
- Keep navigation compact enough that the first useful action is visible at
  1366 by 768.
- Preserve skip links, landmarks, keyboard focus, labels, and readable mobile
  layout.

## Audit Scope

`make audit-product-language` checks primary user-facing source for banned
phrases and raw internal flag names where they would appear as copy. It
intentionally excludes evidence archives, release ledgers, tests, JSON field
definitions, and the prompt file used to request this work.

`make audit-ui-layout` checks stable static layout rules. It is not a visual
review replacement; it exists to catch obvious regressions before review.
