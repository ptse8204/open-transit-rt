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
.skip-link{position:absolute;left:-999px;top:var(--space-2);background:var(--color-text);color:#fff;padding:var(--space-2) var(--space-3);border-radius:var(--radius-1);z-index:10}
.skip-link:focus,.skip-link:focus-visible{left:var(--space-2)}
a:focus-visible,button:focus-visible,input:focus-visible,select:focus-visible,textarea:focus-visible,main:focus-visible{outline:3px solid var(--color-focus);outline-offset:2px}
.operations-header{max-width:82rem;margin:var(--space-5) auto var(--space-3);padding:var(--space-5);background:var(--color-surface);border:1px solid var(--color-border);border-radius:var(--radius-3);box-shadow:var(--shadow-1)}
.operations-header h1{font-size:var(--font-title);line-height:1.2;margin:0 0 var(--space-2);font-weight:750;letter-spacing:0}
.shell-kicker{font-size:var(--font-small);font-weight:700;text-transform:uppercase;letter-spacing:0;color:var(--color-muted);margin:0 0 var(--space-1)}
.shell-breadcrumb{font-size:var(--font-small);margin:0 0 var(--space-3);color:var(--color-muted)}
.shell-meta{display:flex;flex-wrap:wrap;gap:var(--space-2);margin:0;color:var(--color-muted)}
.shell-meta span{display:inline-flex;align-items:center;gap:var(--space-1)}
.operations-nav{max-width:82rem;display:grid;grid-template-columns:repeat(auto-fit,minmax(13rem,1fr));gap:var(--space-3);margin:0 auto var(--space-4)}
.nav-group{border:1px solid var(--color-border);border-radius:var(--radius-3);padding:var(--space-3);background:var(--color-surface);box-shadow:var(--shadow-1)}
.nav-group-label{font-weight:750;margin:0 0 var(--space-2);font-size:var(--font-small);color:var(--color-muted);letter-spacing:0}
.nav-links{display:flex;flex-wrap:wrap;gap:var(--space-2)}
.nav-link{border:1px solid var(--color-border);border-radius:var(--radius-1);padding:.45rem .6rem;min-height:2.25rem;text-decoration:none;color:var(--color-text);background:var(--color-surface-muted);display:inline-flex;align-items:center;gap:var(--space-1);font-weight:650}
.nav-link:focus,.nav-link:hover{border-color:var(--color-border-strong);background:#eef6f7;color:var(--color-action-strong)}
.nav-link.current{border-color:var(--color-action-strong);background:var(--color-action-strong);color:#fff}
.nav-surface{font-size:.72rem;font-weight:700;color:inherit;border:1px solid currentColor;border-radius:var(--radius-1);padding:0 var(--space-1);opacity:.85}
main{max-width:82rem;margin:0 auto var(--space-6);padding:var(--space-5);background:var(--color-surface);border:1px solid var(--color-border);border-radius:var(--radius-3);box-shadow:var(--shadow-1)}
main h2{font-size:var(--font-large);line-height:1.25;margin:var(--space-5) 0 var(--space-3);letter-spacing:0}
main h2:first-child{margin-top:0}
main h3{font-size:1.05rem;line-height:1.25;margin:var(--space-4) 0 var(--space-2);letter-spacing:0}
table{border-collapse:collapse;width:100%;margin:var(--space-4) 0;background:var(--color-surface);font-size:.94rem}
th,td{border:1px solid var(--color-border);padding:var(--space-2);text-align:left;vertical-align:top}
th{background:var(--color-surface-muted);font-weight:750}
tbody tr:nth-child(even){background:#fbfdfe}
.pill,.status-chip{display:inline-block;border:1px solid var(--color-border);border-radius:var(--radius-1);padding:.12rem .4rem;background:var(--color-surface-muted);font-size:var(--font-small)}
.status-chip{font-weight:750}
.status-ready,.status-ready-for-local-review,.status-ok,.status-configured,.status-recorded{border-color:#80b88d;background:var(--color-success);color:var(--color-success-text)}
.status-needs-review,.status-warning,.status-yellow{border-color:#d5a54f;background:var(--color-warning);color:var(--color-warning-text)}
.status-missing,.status-blocked,.status-failed,.status-red{border-color:#df9188;background:var(--color-danger);color:var(--color-danger-text)}
.status-unknown,.status-diagnostic-only,.status-not-run,.status-not-available{border-color:var(--color-border-strong);background:var(--color-surface-muted);color:#364250}
.hero{border:1px solid #a9bfce;background:#f5f9fb;padding:var(--space-4);border-radius:var(--radius-3);margin:var(--space-4) 0}
.hero.start-here{border-color:#7fb0b8;background:#f0f8f8}
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
.section-note,.boundary-notice{border:1px solid var(--color-border);background:var(--color-surface-raised);padding:var(--space-3);border-radius:var(--radius-2);margin:var(--space-3) 0}
.context-help{max-width:82rem;border:1px solid #adc7d4;background:#f4f9fb;border-radius:var(--radius-3);padding:var(--space-4);margin:0 auto var(--space-4);box-shadow:var(--shadow-1)}
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
@media (max-width:700px){body{margin:0;padding:1rem}.operations-header,main,.context-help{margin-left:0;margin-right:0;padding:1rem}.operations-nav,.card-grid,.feed-copy-grid,.context-help-grid{grid-template-columns:1fr}.nav-links{display:grid;grid-template-columns:1fr}.nav-link,button{width:100%}table{display:block;max-width:100%;overflow-x:auto;-webkit-overflow-scrolling:touch}input,select,textarea{min-width:0;width:100%}.shell-meta{display:grid;grid-template-columns:1fr}}
`

func operationsCSS() template.CSS {
	return template.CSS(operationsDesignCSS)
}
