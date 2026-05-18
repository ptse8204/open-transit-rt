# Agent History Archive

Long phase ledgers, historical handoffs, old autonomous roadmap packs, and
prompt archives were preserved on the local branch:

```text
archive/agent-history
```

That branch points at commit `ad6cc9c`, immediately before the main-branch
history trim. Use it when you need old implementation archaeology.

Main keeps only the active maintenance context:

- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/roadmaps/external-connector-runtime-integration/`
- `docs/roadmaps/post-rc2-browser-first-product/`
- `docs/agent/`
- `docs/phase-60-final-claim-review-and-public-closeout.md`
- evidence and consumer-tracker boundary docs

Do not make agency users read archived history to run the product. Update the
public site, README, user docs, and current status first.
