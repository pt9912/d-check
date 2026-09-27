# CHANGELOG.md wird in der Release-Prep gepflegt, nicht im Feature-Commit

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

`CHANGELOG.md` wird bei nutzersichtbaren Änderungen gepflegt — **in der
Release-Prep, nicht im Feature-Commit.** Die Datei führt **keinen**
`[Unreleased]`-Abschnitt: jeder Eintrag steht unter seiner Versions-Nummer,
und die steht erst fest, wenn das Release geschnitten wird. Ein Slice, der
seine Zeile vorzieht, muss sie beim Bump wieder anfassen. Dieselbe Grenze
gilt den beiden `README*.md` und dem Handbuch-Kopf. Ein fehlender Eintrag
dort im Feature-Commit ist deshalb kein Rückstand.
