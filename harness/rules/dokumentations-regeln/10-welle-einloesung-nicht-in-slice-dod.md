# Was eine Welle einlöst, gehört in ihren Closure-Trigger — nicht in die Slice-DoD

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

Ein DoD-Punkt, den der Slice **selbst nicht abhaken kann** (das Release,
das erst mit der Welle fällt), zwingt ihn, mit offenem Haken zu schließen.
Damit ist der Haken als Zustandsfeld unbrauchbar: er sagt nicht mehr *„hier
fehlt etwas"*, sondern *„hier fehlt vielleicht etwas"*. Der
Wellen-Closure-Trigger trägt ihn; die Slice-DoD nennt ihn gar nicht.
**Bestands-Grenze:** vor dieser Regel geschriebene `done/`-Slices behalten
ihre Form — ein nachträglich umgeschriebener DoD-Punkt fälschte einen
Lauf-Beleg.

## Ein Sensor hält das

Und zwar am **Ruheort**: eine `structure`-Regel im Closure-Profil meldet
jeden offenen DoD-Haken eines `done/`-Slice (`max-open-tasks: 0` ⇒
`section-tasks-open`, je Haken auf seiner Zeile, mit verfasstem
Reparatur-Hinweis). Sie läuft in `make verify-closure-notes`, **nicht** in
`gates` — sonst meldete sie beim Arbeiten an einem laufenden Slice. Der
Altbestand ist mit fester Ziffernzahl ausgenommen
([`MR-056`](../../conventions.md#mr-056)).

**Drei Grenzen gehören dazu:** ein Haken ist eine **Selbstauskunft** — die
Regel verschiebt die Lücke von *unsichtbar* nach *behauptet* und prüft
keinen Review; und ein **vergessener Schluss-Fence** macht die
**Bedingung** blind (isoliert gemessen: 0 Befunde, Exit 0), weshalb
dasselbe Profil `spans` fährt — `fence-unclosed` meldet den Fall. **Der
Bindepunkt als Ganzes wird davon nicht grün:** im heutigen Profil melden
Nachbarregeln, und `spans` nennt die Ursache. Und **ein Haken INNERHALB
eines wohlgeformten Fenced-Blocks ist unsichtbar** — dort meldet auch
`fence-unclosed` nichts; dieselbe Fence-Treue, die eine Illustration
schützt, ist der Weg, einen Haken zu verstecken.
