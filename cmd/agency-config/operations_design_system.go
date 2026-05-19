package main

import "html/template"

const operationsDesignCSS = `
:root{
--color-page:#eef3f6;
--color-surface:#ffffff;
--color-surface-muted:#f6f8fa;
--color-surface-raised:#fbfcfd;
--color-text:#17202a;
--color-muted:#53616f;
--color-border:#cfd8e3;
--color-border-strong:#9aa8b7;
--color-action:#0f5f6d;
--color-action-strong:#0b4b55;
--color-focus:#0b76d1;
--color-info:#dbeafe;
--color-info-text:#1e3a5f;
--color-success:#dff6e8;
--color-success-text:#14532d;
--color-warning:#fff3c4;
--color-warning-text:#713f12;
--color-danger:#ffe1df;
--color-danger-text:#7f1d1d;
--space-1:.25rem;
--space-2:.5rem;
--space-3:.75rem;
--space-4:1rem;
--space-5:1.5rem;
--space-6:2rem;
--radius-1:4px;
--radius-2:6px;
--radius-3:8px;
--shadow-1:0 1px 2px rgba(15,23,42,.08);
--font-small:.875rem;
--font-base:1rem;
--font-large:1.2rem;
--font-title:1.75rem;
}
*{box-sizing:border-box}
html{background:var(--color-page)}
body{font-family:system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;margin:0;line-height:1.5;color:var(--color-text);background:var(--color-page);overflow-wrap:anywhere}
a{color:var(--color-action)}
a:hover{color:var(--color-action-strong)}
.sr-only{position:absolute;width:1px;height:1px;padding:0;margin:-1px;overflow:hidden;clip:rect(0,0,0,0);white-space:nowrap;border:0}
.skip-link{position:absolute;left:-999px;top:var(--space-2);background:var(--color-text);color:#fff;padding:var(--space-2) var(--space-3);border-radius:var(--radius-1);z-index:10}
.skip-link+.skip-link:focus,.skip-link+.skip-link:focus-visible{top:calc(var(--space-2) + 2.8rem)}
.skip-link:focus,.skip-link:focus-visible{left:var(--space-2)}
a:focus-visible,button:focus-visible,input:focus-visible,select:focus-visible,textarea:focus-visible,summary:focus-visible,main:focus-visible{outline:3px solid var(--color-focus);outline-offset:3px}
.operations-header{max-width:92rem;margin:0 auto var(--space-3);padding:var(--space-3) var(--space-4);background:var(--color-surface);border-bottom:1px solid var(--color-border);display:flex;align-items:end;justify-content:space-between;gap:var(--space-4)}
.operations-header h1{font-size:var(--font-title);line-height:1.2;margin:0;font-weight:750;letter-spacing:0}
.header-title{min-width:0}
.app-kicker{font-size:var(--font-small);font-weight:700;text-transform:uppercase;letter-spacing:0;color:var(--color-muted);margin:0 0 var(--space-1)}
.app-breadcrumb{font-size:var(--font-small);margin:0;color:var(--color-muted)}
.app-meta{display:flex;flex-wrap:wrap;justify-content:flex-end;gap:var(--space-2);margin:0;color:var(--color-muted);font-size:var(--font-small)}
.app-meta span{display:inline-flex;align-items:center;gap:var(--space-1)}
.operations-frame{max-width:92rem;margin:0 auto var(--space-6);display:grid;grid-template-columns:minmax(12rem,16rem) minmax(0,1fr);gap:var(--space-4);align-items:start;padding:0 var(--space-4)}
.page-next-action{border:1px solid #94b8c2;background:#f0f8f8;border-radius:var(--radius-2);padding:var(--space-3);margin:var(--space-4) 0 0}
.page-next-action:first-child{margin-top:0}
.page-next-action h2{font-size:1rem;margin:0 0 var(--space-1)}
.page-next-action p{margin:.25rem 0}
.scope-banner{border:1px solid var(--color-border);border-radius:var(--radius-2);background:var(--color-surface-muted);padding:var(--space-3);margin:var(--space-4) 0 0}
.scope-banner h2{font-size:1rem;margin:0 0 var(--space-2);letter-spacing:0}
.scope-banner p{margin:.35rem 0}
.scope-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(14rem,1fr));gap:var(--space-3);margin:0}
.scope-grid div{border:1px solid var(--color-border);border-radius:var(--radius-1);background:var(--color-surface);padding:var(--space-2)}
.scope-grid dt{font-size:var(--font-small);font-weight:750;color:var(--color-muted);margin:0 0 var(--space-1)}
.scope-grid dd{margin:0}
.operations-nav{position:sticky;top:var(--space-3);align-self:start;border:1px solid var(--color-border);border-radius:var(--radius-2);background:var(--color-surface);padding:var(--space-3);max-height:calc(100vh - 2rem);overflow:auto}
.nav-group{border:0;border-top:1px solid var(--color-border);border-radius:0;padding:var(--space-2) 0;background:transparent;box-shadow:none}
.nav-group:first-child{border-top:0;padding-top:0}
.nav-group:last-child{padding-bottom:0}
.nav-group-label{font-weight:750;margin:0 0 var(--space-1);font-size:.78rem;color:var(--color-muted);letter-spacing:0;text-transform:uppercase}
.nav-links{display:grid;gap:.2rem}
.nav-link{border:1px solid transparent;border-radius:var(--radius-1);padding:.32rem .45rem;min-height:2rem;text-decoration:none;color:var(--color-text);background:transparent;display:flex;align-items:center;justify-content:space-between;gap:var(--space-1);font-weight:650;font-size:var(--font-small)}
.nav-link:focus,.nav-link:hover{border-color:var(--color-border-strong);background:#eef6f7;color:var(--color-action-strong)}
.nav-link.current{border-color:var(--color-action-strong);background:var(--color-action-strong);color:#fff}
.nav-link.current:focus-visible{outline-color:#111827;box-shadow:0 0 0 3px #fff}
.nav-surface{font-size:.72rem;font-weight:700;color:inherit;border:1px solid currentColor;border-radius:var(--radius-1);padding:0 var(--space-1);opacity:.85}
main{min-width:0;margin:0;padding:var(--space-5);background:var(--color-surface);border:1px solid var(--color-border);border-radius:var(--radius-3);box-shadow:var(--shadow-1)}
main h2{font-size:var(--font-large);line-height:1.25;margin:var(--space-5) 0 var(--space-3);letter-spacing:0}
main h2:first-child{margin-top:0}
main h3{font-size:1.05rem;line-height:1.25;margin:var(--space-4) 0 var(--space-2);letter-spacing:0}
table{border-collapse:collapse;width:100%;margin:var(--space-4) 0;background:var(--color-surface);font-size:.94rem}
caption{text-align:left;font-weight:750;color:var(--color-muted);padding:0 0 var(--space-2)}
th,td{border:1px solid var(--color-border);padding:var(--space-2);text-align:left;vertical-align:top}
th{background:var(--color-surface-muted);font-weight:750}
tbody tr:nth-child(even){background:#fbfdfe}
tbody tr:focus-within{outline:2px solid var(--color-focus);outline-offset:-2px}
.pill,.status-chip{display:inline-block;border:1px solid var(--color-border);border-radius:var(--radius-1);padding:.12rem .4rem;background:var(--color-surface-muted);font-size:var(--font-small)}
.status-chip{font-weight:750}
.status-ready,.status-ready-for-local-review,.status-ok,.status-configured,.status-recorded{border-color:#80b88d;background:var(--color-success);color:var(--color-success-text)}
.status-needs-review,.status-warning,.status-yellow{border-color:#d5a54f;background:var(--color-warning);color:var(--color-warning-text)}
.status-missing,.status-blocked,.status-failed,.status-red{border-color:#df9188;background:var(--color-danger);color:var(--color-danger-text)}
.status-unknown,.status-diagnostic-only,.status-not-run,.status-not-available{border-color:var(--color-border-strong);background:var(--color-surface-muted);color:#364250}
.hero{border:1px solid #a9bfce;background:#f5f9fb;padding:var(--space-4);border-radius:var(--radius-3);margin:var(--space-4) 0}
.hero.start-here{border-color:#7fb0b8;background:#f0f8f8}
.workflow-hero{display:grid;grid-template-columns:minmax(0,1fr) minmax(18rem,.45fr);gap:var(--space-4);align-items:start;border:1px solid #83adb7;background:#f0f8f8;padding:var(--space-4);border-radius:var(--radius-3);margin:0 0 var(--space-4)}
.workflow-hero h2{margin-top:0}
.workflow-summary{border:1px solid #adc7d4;border-radius:var(--radius-2);background:var(--color-surface);padding:var(--space-3)}
.workflow-summary h3{margin-top:0}
.workflow-summary dl{display:grid;gap:var(--space-2);margin:0}
.workflow-summary div{border-top:1px solid var(--color-border);padding-top:var(--space-2)}
.workflow-summary div:first-child{border-top:0;padding-top:0}
.workflow-summary dt{font-size:var(--font-small);font-weight:750;color:var(--color-muted)}
.workflow-summary dd{margin:0}
.workflow-steps{list-style:none;counter-reset:workflow;display:grid;grid-template-columns:repeat(auto-fit,minmax(22rem,1fr));gap:var(--space-4);padding:0;margin:var(--space-4) 0}
.workflow-step{counter-increment:workflow;border:1px solid #96b8c0;border-radius:var(--radius-3);padding:var(--space-4);background:#f7fbfb;box-shadow:var(--shadow-1)}
.workflow-step-header{display:grid;grid-template-columns:auto 1fr;gap:var(--space-3);align-items:start}
.workflow-step-number::before{content:counter(workflow);display:inline-grid;place-items:center;width:2rem;height:2rem;border-radius:999px;background:var(--color-action);color:#fff;font-weight:800}
.workflow-step h3{margin:0 0 var(--space-1)}
.dashboard-tools,.dashboard-details{margin:var(--space-4) 0}
.action-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(17rem,1fr));gap:var(--space-4);margin:var(--space-4) 0}
.action-card{border:1px solid #96b8c0;border-radius:var(--radius-3);padding:var(--space-4);background:#f7fbfb;box-shadow:var(--shadow-1)}
.action-card h3{margin-top:0}
.action-card p{margin:.45rem 0}
.action-link{display:inline-flex;align-items:center;justify-content:center;min-height:2.25rem;border:1px solid var(--color-action-strong);border-radius:var(--radius-1);background:var(--color-action);color:#fff;text-decoration:none;font-weight:750;padding:.5rem .75rem;margin-top:var(--space-2)}
.action-link:hover,.action-link:focus{background:var(--color-action-strong);color:#fff}
.secondary-action{background:var(--color-surface);color:var(--color-action-strong)}
.secondary-action:hover,.secondary-action:focus{background:#eef6f7;color:var(--color-action-strong)}
.status-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(15rem,1fr));gap:var(--space-3);margin:var(--space-4) 0}
.status-tile{border:1px solid var(--color-border);border-radius:var(--radius-2);padding:var(--space-3);background:var(--color-surface-raised)}
.status-tile h3{font-size:1rem;margin:0 0 var(--space-2)}
.compact-actions{display:flex;flex-wrap:wrap;gap:var(--space-2);margin:var(--space-3) 0}
.compact-actions a:not(.action-link){display:inline-flex;align-items:center;min-height:2.1rem;border:1px solid var(--color-border);border-radius:var(--radius-1);padding:.4rem .55rem;background:var(--color-surface-muted);text-decoration:none;font-weight:650}
.card-grid,.feed-copy-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(16rem,1fr));gap:var(--space-4);margin:var(--space-4) 0}
.feed-copy-grid{grid-template-columns:repeat(auto-fit,minmax(18rem,1fr))}
.card,.feed-url-card{border:1px solid var(--color-border);border-radius:var(--radius-3);padding:var(--space-4);background:var(--color-surface-raised);box-shadow:var(--shadow-1)}
.card h3,.feed-url-card h3{margin-top:0}
.card p,.feed-url-card p{margin:.45rem 0}
.empty-state{border-color:#9dbfd6;background:#f5fbff}
.path-card{border-color:#8cb8c2}
.path-developer{border-color:#aaa0c6}
.status{font-weight:650}
.copy-value{display:block;border:1px solid var(--color-border);background:var(--color-surface-muted);border-radius:var(--radius-1);padding:var(--space-2);white-space:pre-wrap;overflow-wrap:anywhere}
.review-tools{border:1px solid var(--color-border);background:var(--color-surface-raised);border-radius:var(--radius-3);padding:var(--space-4);margin:var(--space-4) 0}
.review-tools h3{margin-top:0}
.review-controls{display:flex;flex-wrap:wrap;gap:var(--space-3);align-items:end}
.review-controls label{margin:0}
.review-controls input,.review-controls select{min-width:12rem}
.review-status{display:block;margin:.35rem 0;color:var(--color-muted);font-size:var(--font-small)}
.copy-action{margin:.35rem 0;background:var(--color-surface-muted);color:var(--color-action-strong);border-color:var(--color-border-strong)}
[hidden]{display:none!important}
.section-note,.boundary-notice{border:1px solid var(--color-border);background:var(--color-surface-raised);padding:var(--space-3);border-radius:var(--radius-2);margin:var(--space-3) 0}
.support-panels{margin-top:var(--space-5)}
.support-details{background:#f8fbfc}
.support-details>.scope-banner,.support-details>.context-help{margin:var(--space-3) 0 0}
.context-help{border:1px solid #adc7d4;background:#f4f9fb;border-radius:var(--radius-3);padding:var(--space-4);margin:var(--space-4) 0;box-shadow:none}
.context-help h2{font-size:1.05rem;margin:0 0 var(--space-3);letter-spacing:0}
.context-help-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(14rem,1fr));gap:var(--space-3)}
.context-help-topic{border:1px solid #b8d0db;border-radius:var(--radius-2);background:var(--color-surface);padding:var(--space-3)}
.context-help-topic h3{font-size:1rem;margin:.1rem 0;letter-spacing:0}
.context-help-topic p{margin:.3rem 0}
.warning{background:var(--color-warning);color:var(--color-warning-text)}
.ok{background:var(--color-success);color:var(--color-success-text)}
.bad{background:var(--color-danger);color:var(--color-danger-text)}
.muted{color:var(--color-muted)}
.token{border:1px solid #d6a44f;background:var(--color-warning);padding:var(--space-4);border-radius:var(--radius-2)}
form{margin:var(--space-4) 0;padding:var(--space-4);border:1px solid var(--color-border);border-radius:var(--radius-3);background:var(--color-surface-raised)}
fieldset{border:1px solid var(--color-border);border-radius:var(--radius-2);padding:var(--space-3);margin:var(--space-3) 0}
legend{font-weight:750;padding:0 var(--space-1)}
label{display:block;margin:.35rem 0;font-weight:650}
input,select,textarea{min-width:22rem;max-width:100%;padding:.5rem;border:1px solid var(--color-border-strong);border-radius:var(--radius-1);background:#fff;color:var(--color-text)}
button{padding:.55rem .85rem;min-height:2.25rem;border:1px solid var(--color-action-strong);border-radius:var(--radius-1);background:var(--color-action);color:#fff;font-weight:750;cursor:pointer}
button:hover{background:var(--color-action-strong)}
details{border:1px solid var(--color-border);border-radius:var(--radius-2);padding:var(--space-3);margin:var(--space-3) 0;background:var(--color-surface-raised)}
summary{font-weight:750;cursor:pointer}
@media (prefers-reduced-motion:reduce){*,*::before,*::after{scroll-behavior:auto!important;transition:none!important;animation:none!important}}
@media (prefers-contrast:more){:root{--color-page:#fff;--color-surface:#fff;--color-surface-muted:#f2f2f2;--color-text:#000;--color-muted:#1f2937;--color-border:#111827;--color-border-strong:#000;--color-action:#003f8c;--color-action-strong:#001f4d;--color-focus:#ffbf00}.status-ready,.status-ready-for-local-review,.status-ok,.status-configured,.status-recorded,.status-needs-review,.status-warning,.status-yellow,.status-missing,.status-blocked,.status-failed,.status-red,.status-unknown,.status-diagnostic-only,.status-not-run,.status-not-available{border-color:#000}.nav-link.current{background:#000;color:#fff}}
@media (max-width:900px){.operations-header{display:grid;align-items:start}.app-meta{justify-content:flex-start}.operations-frame{grid-template-columns:1fr}.operations-nav{position:static;max-height:none}.nav-links{grid-template-columns:repeat(auto-fit,minmax(10rem,1fr))}}
@media (max-width:700px){body{margin:0}.operations-header,.operations-frame{padding-left:1rem;padding-right:1rem}.operations-header,main,.context-help{margin-left:0;margin-right:0;padding:1rem}.operations-nav,.card-grid,.feed-copy-grid,.context-help-grid,.action-grid,.status-grid,.workflow-hero,.workflow-steps{grid-template-columns:1fr}.nav-links,.review-controls{display:grid;grid-template-columns:1fr}.nav-link,button{width:100%}.action-link{width:100%}table{display:block;max-width:100%;overflow-x:auto;-webkit-overflow-scrolling:touch}input,select,textarea{min-width:0;width:100%}.review-controls input,.review-controls select{min-width:0;width:100%}.app-meta{display:grid;grid-template-columns:1fr}}
`

func operationsCSS() template.CSS {
	return template.CSS(operationsDesignCSS)
}
