# Neue oder geänderte DC-*-Anforderungen entstehen nur im Lastenheft

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

Neue oder geänderte `DC-*`-Anforderungen entstehen nur in
[`spec/lastenheft.md`](../../../spec/lastenheft.md) — nie per ADR (ADRs
schärfen die Spezifikation, nicht das Lastenheft). Der Anlege-Prozess
(Akzeptanzkriterien-Trio, Versions-Bump + Historie, Beleg-Pflicht) folgt
dem Baseline-Regelwerk
([`modul-03-spec`](../../../.harness/baseline/v6.13.0/regelwerk/modul-03-spec.md));
das repo-spezifische ID-Schema steht in `spec/lastenheft.md` §3.
