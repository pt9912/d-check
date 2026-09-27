# Neue ADRs tragen einen Re-Evaluierungs-Trigger

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

Neue ADRs tragen die Sektion `## Re-Evaluierungs-Trigger` (oder
„permanent"); die vor Einführung `Accepted`-ADRs sind immutable und
**grandfathered** (das Trigger-Feld liegt im ADR-Core, nachträgliches
Ergänzen bräche `make adr-check`). Der Welle-Closure-Trigger-Audit
(Baseline-Regelwerk Modul 6) bestätigt oder revidiert sie (Folge-ADR mit
`supersedes`).
