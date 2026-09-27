# Slice-Kopf-Feld Berührte Spec-Stellen nennt die Kennung

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

Das Slice-Kopf-Feld `**Berührte Spec-Stellen:**` nennt die **Kennung**, wo
das Zielelement eine trägt (`SPEC-<NNN>`, `ARC-<NNN>`,
`<DC-ID>.<Buchstabe>`), sonst den Abschnitt; `—`, wenn der Slice keine
Spec-Stelle berührt. Der Verweis zeigt **aufwärts** — die Spec nennt den
Slice nie (`AGENTS.md` §3.4). Feld-Form aus der Baseline-`slice.template.md`,
die Kennungs-Regel aus
[`MR-000`](../../conventions.md#mr-000--baseline-aussage).
