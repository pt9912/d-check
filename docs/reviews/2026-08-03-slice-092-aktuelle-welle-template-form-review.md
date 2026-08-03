# Review — slice-092: `## Aktuelle Welle` auf die Template-Struktur-Felder

- **Review-Art:** Plan/Design-Review (Doku-/Harness-Slice, kein Produkt-Code) —
  geprüft gegen Baseline-Template, `planning`-Gate-Invariant und
  Konventionsspeicher-Regeln.
- **Gegenstand:** slice-092 (welle-68), Commit-Range `e4161d2..HEAD`
  (`8a3aeee` welle-68-Eröffnung + slice-092-Open + Roadmap-Struktur-Feld-Form;
  `4a12846` `MR-024`-Verfeinerung).
- **Skill:** `reviewer.md` @ Version 1.3.0.
- **Modell-ID:** claude-opus-4-8.
- **Datum:** 2026-08-03.
- **Eingangs-Kontext (Verträge, gegen die geprüft wurde):**
  - Ziel-Form `.harness/baseline/v5.0.0/templates/docs/plan/planning/roadmap.template.md`
    (§Aktuelle Welle, §Nächste Wellen) + `welle.template.md`.
  - Slice-Plan `docs/plan/planning/in-progress/slice-092-roadmap-aktuelle-welle-template-form.md`.
  - `planning`-Modul-Invariant (`DC-FA-PLAN-001` / `ADR-0028`): `hasActive == hasSlices`.
  - Konventionsspeicher `harness/conventions.md` + `harness/conventions/MR-024-…md`.
  - Hard Rules `AGENTS.md` §3.

## Verifikation (selbst ausgeführt)

- `make planning-check` → `324 Datei(en) geprüft, 0 Befund(e)` (grün).
- `make adr-check` (`vcs`, `HEAD~1..HEAD`) → `0 Befund(e)` (grün).
- `make gates` → alle Module `0 Befund(e)`; `coverage-gate: OK — 94.20% ≥ 93%`.
- Einzelläufe `links` / `anchors` / `ids` über das Image → `0 Befund(e)`.
- `git diff --stat e4161d2..HEAD`: nur `roadmap.md`, `slice-092-…md` (neu),
  `welle-68-…md` (neu), `conventions.md`, `MR-024-…md` — **keine** `.go`-,
  `spec/`- oder `adr/`-Datei berührt.
- `awk`-Extraktion des `## Aktuelle Welle`-Blocks (H2 bis nächste H2);
  `grep`-Sweeps auf Ruhe-Marker-String, alten/neuen Anker-Slug, `#NAME?`.

## Findings

### INFO-1 — `MR-024`-`Datum`-Feld nicht mit der Verfeinerung fortgeschrieben

- **kategorie:** INFO
- **quelle:** `MR-024` (Konventionsspeicher-Eintrag)
- **pfad:** `harness/conventions/MR-024-aktuelle-welle-ruhe-marker-form.md:5`
- **befund:** Das Kopf-`Datum` steht auf `2026-08-02` (Erst-Annahme in slice-088,
  Commit `7d8d085`), während dieser Commit (`4a12846`, 2026-08-03) die Adaptions-
  Semantik materiell umgekehrt hat (H1-Titel + „aktiv = Template-Felder, wellenlos =
  Ruhe-Marker"). Der Rumpf trägt die richtige Provenienz („seit slice-092",
  „Nachtrag slice-092"); nur die Kopf-Datumszeile bleibt auf dem Erst-Annahme-Tag.
- **failure-scenario:** Ein Auditor quervergleicht das `Datum` `2026-08-02` mit dem
  Meilenstein-Log (`roadmap.md`, Zeile 163: welle-68/slice-092 = `2026-08-03`) und
  datiert die neue „aktiv = template-konform"-Regel einen Tag vor den Slice, der sie
  eingeführt hat. Kein Gate erzwingt das Feld; Auswirkung rein dokumentarisch.
- **verifizierbar:** nein (kein Gate prüft `Datum`-Aktualität) — belegbar per
  `git log --follow` auf die Datei.
- **klasse:** `mr-datum-nicht-fortgeschrieben`
- **Hinweis:** Falls `Datum` bewusst als „Erst-Annahme-Tag" geführt wird (die
  Rumpf-Nachträge tragen dann die Änderungshistorie), ist der Befund Won't-Fix.

## Negativbefunde (geprüft, ohne Befund)

- **planning-Gate-Korrektheit (Kern):** Der `## Aktuelle Welle`-Block (H2→H2) in
  `roadmap.md` enthält die Ruhe-Marker-Zeichenfolge **nicht** (`grep` → not found);
  `planning-check` grün, welle-68 aktiv + `slice-092` in `in-progress/`
  (`hasActive == hasSlices`, beide true). Kein `planning-drift`. **Ohne Befund.**
- **Vier Struktur-Felder in Template-Form:** Der Block trägt `**Welle-ID:**`,
  `**Start:**`, `**Geplantes Ende:**` (`… (Schätzung, korrigierbar)`) und
  `**Closure-Trigger:**` plus die Slice-Liste (Pflicht-Bestandteil Slice-IDs aus
  `modul-06`). Feld-Namen und Reihenfolge decken sich mit `roadmap.template.md`.
  **Ohne Befund.**
- **Wellen-/Slice-Konsistenz:** Das Wellendokument `welle-68-…md` §4 nennt
  `slice-092` (verlinkt, `in-progress/…` löst auf: Datei existiert) und `slice-093`
  (Klartext, noch keine Datei — korrekt für einen geplanten Folge-Slice, kein toter
  Link). `slice-092` führt **kein** `**Status:**`-Feld und trägt den `**Lifecycle:**`-
  Hinweis (D-5-Konvention). **Ohne Befund.**
- **`MR-024`-Faktentreue:** Text deckt sich mit der Realität — aktiv = Template-Felder
  seit slice-092, wellenlos = Ruhe-Marker (gate-erzwungen), ausdrücklich **kein**
  `planning`-Modul-Umbau. **Ohne Befund.**
- **Index-Anker-Konsistenz:** Der `<a id>`-Slug in `conventions.md:100` stimmt mit dem
  neuen H1-Titel überein; der **alte** Slug wird repo-weit **nirgends** referenziert
  (`grep` → 0 Treffer), also bricht die Ankeränderung keinen Verweis. Kein `#NAME?`/
  `#REF`/`#VALUE` in der Index-Tabelle. Kein anderer `MR-*`-Eintrag zeigt auf `MR-024`;
  das Entfernen des `MR-014`-Tie-ins hinterlässt keine tote Rückreferenz. `anchors`/
  `links`/`ids` grün. **Ohne Befund.**
- **Template-Treue der „Regeln dieser Sektion"-Zeiger:** §Aktuelle Welle und
  §Nächste Wellen decken sich **wörtlich** mit `roadmap.template.md`. **Ohne Befund.**
- **Berührte `Accepted`-ADR/Spec/Code:** Keine — `git diff --stat` zeigt ausschließlich
  Planning-/Konventions-Doku. **Ohne Befund.**
- **`observations.md`:** trägt `— keine —`; die slice-092-§7-Aussage („nichts zu
  berücksichtigen") ist faktentreu. **Ohne Befund.**
- **Forward-Referenz Closure-Trigger:** `roadmap.md` und `welle-68-…md` nennen als
  Closure-Kriterium das noch nicht existierende `verify-closure-notes`-Gate — korrekt,
  da es slice-093-Ergebnis ist und den Closure-Zustand (nicht den Ist-Zustand)
  beschreibt (das „Mehr" gegenüber den Slice-DoDs, `modul-06`). **Ohne Befund.**

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

## Verdikt

**abnahmereif.** Der Kern-Prüfpunkt hält: der `## Aktuelle Welle`-Block trägt die vier
Template-Struktur-Felder ohne Ruhe-Marker-Zeichenfolge, das `planning`-Modul ist grün
(aktive Welle + Slice in `in-progress/`), und alle Gates (`make gates`, `adr-check`)
laufen sauber. Wellen-/Slice-Dokumente, `MR-024`-Verfeinerung, Index-Anker und
Template-Zeiger sind konsistent und faktentreu; kein `Accepted`-Artefakt berührt. Der
einzige Befund (INFO-1) ist rein dokumentarisch und blockiert nicht.
