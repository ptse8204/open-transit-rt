# Public Product Reality Audit

Date: 2026-05-18

This audit checks the actual `main/site` source, the checked-out `gh-pages`
worktree, and the private Operations Console code. It does not rely on status
summaries.

## Result

The public site and private admin UI are not yet at the product bar for small
agency staff. The current work is real, but the experience still reads like
maintainer output in several places.

## Live Website

The `gh-pages` worktree contains the same public pages as `site/`:

- `index.html`
- `ui-tour.html`
- `connectors.html`
- `readiness.html`
- `video.html`
- `assets/site.css`
- `assets/open-transit-rt-browser-first-tutorial.mp4`
- `assets/open-transit-rt-browser-first-tutorial.vtt`

So the live branch has the current page set. The problem is quality, not
publication drift.

## Favicon

The public pages still use `href="data:,"` favicon placeholders. The previous
favicon is not restored in the current `site/` tree or the checked-out
`gh-pages` branch.

## UI Tour

The UI tour does not contain actual browser screenshots. It uses generated
browser skeletons and still points readers to screenshot-capture rules. That
makes the page look like an internal artifact instead of a product tour.

Product-critical fix: capture real browser screenshots from the local app and
replace the skeleton UI.

## Video Page

The video page embeds a small MP4, but the page still says the video is a
generated walkthrough and includes recording checklist and storyboard material.
That is maintainer content, not the public tutorial experience.

Product-critical fix: publish an understandable interface walkthrough built
from real UI captures or a real browser recording, keep captions/transcript,
and move recording instructions back to maintainer docs only.

## Readiness Page

The readiness page is still framed as "readiness without overclaiming." It has
useful checks, but the first user question is simpler: "Is the system running,
and what should I do next?"

Product-critical fix: reframe the public page around GTFS import, five feed
URLs, realtime publishing, fresh vehicle data, validators, issues, and next
actions. Keep CAL-ITP-style readiness as a secondary public-data preparation
section.

## Connector Page

The connector page lists useful connector categories, but it still mixes user
actions, local commands, caveats, and roadmap candidates. A prospective user
should immediately see:

- what works today in local/self-hosted evaluation;
- what requires deployment-owned setup;
- what is planned or candidate-only.

Product-critical fix: rewrite the page and connector docs with clear "Works
today" and "Planned / candidate" sections.

## Admin Console

`/admin/local-login` exists and is guarded for local/demo use, but the console
after login is still overloaded. The dashboard includes an action queue,
progress tiles, more action cards, readiness summaries, detailed tables, and
repeated boundary text. The primary flow is present in code, but it is not yet
a straight-line workflow.

Product-critical fix: make the dashboard follow one visible path:

1. Start setup
2. Import GTFS
3. Check feeds
4. Connect vehicles
5. Review realtime
6. Fix issues
7. Share public URLs
8. Maintain system

Move secondary status, technical help, JSON exports, and claim details behind
details panels.

## Maintainer / AI-Agent Writing Still Visible

These phrases and patterns still make public/user surfaces feel internal:

- generated browser capture;
- screenshot capture rules;
- recording checklist and storyboard language on the public video page;
- long repeated "does not prove" paragraphs;
- claim-flag and phase-history framing;
- raw route paths as primary navigation;
- command-heavy public connector copy.

## Docs To Move Out Of The Main User Path

The stable-branch filter already excludes many agent-only paths:

- `docs/agent/`
- `docs/handoffs/`
- `docs/prompts/`
- `docs/roadmaps/`
- `docs/phase-*.md`
- `docs/*phase-*.md`
- `docs/**/*phase-*.md`
- Codex task and conversation summary files.

The product problem is that user pages still link through GitHub markdown for
normal tasks. Keep maintainer and phase-history docs available for development,
but make the website the primary path for:

- trying locally;
- importing GTFS;
- checking feeds;
- connecting vehicles;
- troubleshooting;
- reviewing connector support.

## Product-Critical Versus Documentation-Only

Product-critical:

- restore the favicon on public and private browser pages;
- rewrite public pages in product language;
- replace fake UI tour captures with real rendered screenshots;
- replace public video slop with a real interface walkthrough;
- redesign the admin dashboard around a straight-line workflow;
- make normal user docs available on the website;
- clarify connector support as current versus planned.

Documentation-only:

- archive or filter AI-agent phase history;
- keep maintainer recording rules in `docs/tutorials/video-recording-guide.md`;
- keep detailed connector/conformance references in technical docs;
- update closeout/status pages after the product behavior is verified.

## Claim Boundary

This audit does not claim CAL-ITP/Caltrans compliance, production readiness,
agency adoption or approval, consumer submission/review/acceptance/ingestion,
hosted service availability, vendor compatibility, SLA/uptime, production AVL
reliability, ETA quality, or final-root readiness.
