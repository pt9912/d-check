# Verifikations-Report — slice-244 (DoD)

**Verifikations-Art:** DoD/Plan-Konformität (Modul 11 — „Bauen wir es richtig?"; *nicht* der Diff-Review — der steht im R1-Report)

**Gegenstand:** `slice-244` (welle-91, wellenlos) — `.claude`-Rollen und Commands aus ai-harness-init evaluieren

**Range:** `e0165e8c..HEAD` — `c7091b7b` (Beanspruchung: Ruhe-Marker verlässt die Roadmap), `cb918ee3` (Adoption: 5 Dateien neu), `a90b5d3a` (R1-Report + F-1/F-2/F-3-Korrekturen)

**Plan:** `docs/plan/planning/in-progress/slice-244-claude-rollen-commands-evaluieren.md`

**Modell-ID:** glm-5.3-flash · **Datum:** 2026-09-29

---

## DoD §2 — je Punkt geprüft

### DoD 1 — Je Rolle (6) und Command (3) eine belegte Entscheidung: **ERFÜLLT**

| Kandidat | Entscheidung | Beleg |
|---|---|---|
| architect | ADOPTIERT | `.claude/agents/architect.md` neu in `cb918ee3` |
| planner | ADOPTIERT | `.claude/agents/planner.md` neu in `cb918ee3` |
| validator | ADOPTIERT | `.claude/agents/validator.md` neu in `cb918ee3` |
| implementer | ABGELEHNT mit benanntem Grund | Plan §1 (`slice-244…md:32-34`): „läuft im Hauptlauf nach AGENTS.md §6 bzw. über den implement-slice-Command — ein Agent-Duplikat trüge dieselbe Anweisung an zweiter Stelle"; nachgereicht in `a90b5d3a` (R1-F-1) |
| reviewer | bereits als Agent vorhanden | `.claude/agents/reviewer.md` existiert an `e0165e8c` (Vorbestand slice-235, per `git log` geprüft); Entscheidung im Plan §1 (`:29-30`) |
| verifier | bereits als Agent vorhanden | `.claude/agents/verifier.md` existiert an `e0165e8c`; Entscheidung im Plan §1 (`:29-30`) |
| plan-welle | ADOPTIERT | `.claude/commands/plan-welle.md` neu in `cb918ee3` |
| close-welle | ADOPTIERT | `.claude/commands/close-welle.md` neu in `cb918ee3` |
| implement-slice | bereits als Command vorhanden | `.claude/commands/implement-slice.md` existiert an `e0165e8c`; Entscheidung im Plan §1 (`:30-31`) |

Gegenprobe zum Verzeichnis: `.claude/agents/` enthält **kein** `implementer.md` — die Ablehnung ist auch als Nicht-Akt belegt. R1-F-1 ist wirksam korrigiert: der Plan-Diff von `a90b5d3a` zeigt die Form-Korrektur („implement-slice als **Skill**" → „Command unter `.claude/commands/`"), das Entfernen der Dopplung „wellenlosen/wellenlosen" und die neu gesetzte Implementer-Ablehnung mit Grund.

*Form-Feinheit, kein Mangel:* die strikte Ablehnungs-Form gilt nur für implementer; reviewer/verifier/implement-slice sind als „bereits in eigener Form vorhanden" dokumentiert — beides ist je Kandidat eine belegte Entscheidung im Sinn des DoD („adoptiert … oder abgelehnt"), und die Vorbestands-Evidenz habe ich selbst gegen `e0165e8c` geprüft, nicht nur behauptet übernommen.

### DoD 2 — Kontext-Trennung, d-check-Pendants, Rollen-Achse: **ERFÜLLT**

- **Kontext-Trennung:** `architect.md:23` und `planner.md:22` tragen „Was du NICHT bist" wortgleich aus der Vorlage (Schwester-Snapshot `/tmp/aih-v6.13.0` per diff gegen die Adopt-Dateien geprüft). `validator.md` trägt sie in der Form **der Vorlage**: deren validator hat keine „Was du NICHT bist"-Sektion — die Kontext-Trennung liegt dort in „Dein Kontext-Zuschnitt … Du liest die Anforderung und das Ergebnis, nicht den Diff" (`validator.md:21-22`), byte-identisch adoptiert. Die Adoptions-Abweichungen sind gezielte Schnitte (Erfassungs-Referenzen entfernt, ANPASSEN-Platzhalter ersetzt, planner-Move-Regel auf d-checks 100 %-Form geschärft) — Vorlage-Treue, keine Drift.
- **Schluss-Absatz Rollen-Achse:** „Warum dieser Typ existiert" mit den sechs kanonischen Namen (`planner`, `architect`, `implementer`, `reviewer`, `verifier`, `validator`) in allen drei Dateien — `architect.md:44-48`, `planner.md:28-32`, `validator.md:33-37`; durch eigenen diff gegen die Originale bestätigt.
- **d-check-Pendants:** `planner.md:13-15` verlinkt die eigenen Commands `plan-welle`/`close-welle`; `architect.md:50-52` nennt `docs/plan/adr/` samt Index, `spec/`, `harness/conventions.md`; `validator.md:39-40` nennt Closure-Notiz/Welle-results-Notiz und Bedarfsträger; `validator.md:10-11` verweist auf den d-check-Verifier. Die ANPASSEN-Platzhalter der Vorlage sind durchgehend ersetzt.

### DoD 3 — `make gates` grün: **ERFÜLLT** (selbst gefahren)

`make gates` in frischer Ausführung, Exit-Code 0 — die Ausgabe liegt vor, der Exit-Code ist nicht behauptet:

- Zusammenfassung: `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`
- `doc-check` und `planning-check`: **906 Datei(en) geprüft, 0 Befund(e)** je zweifach in der Ausgabe sichtbar — der behauptete Stand stimmt.
- `coverage-gate`: „OK — Coverage 94.60 % erfüllt Schwelle 93 %"
- `semgrep`: „Ran 55 rules on 65 files: 0 findings."
- `baseline-verify`: „fetch-baseline-cache: verify ok (54 Dateien, vollständig)"
- `test`: alle Pakete `ok`; `lint`/`arch-check` (a-check v0.20.0, digest-gepinnt) grün; `workflow-pins`/`gate-consistency` im grünen Gesamtdurchlauf.

---

## Abgrenzung des Slices gegen den tatsächlichen Diff

- **Erfassungsschicht out:** korrekt eingehalten — `span-emit.sh`/`erfassung.mk` sind nicht im Range-Diff; die Out-of-Scope-Verzahnung zur Welle steht in `welle-91-adoption-ai-harness-init.md:80` (§6).
- **Keine Umbauten an reviewer/verifier:** beide Dateien tauchen im Range-Diff nicht auf (Vorbestand unverändert).
- **Keine Besetzungs-Pflicht:** der Plan §1 schließt sie aus; die Adopt-Dateien beschreiben Rollen, ohne sie zu erzwingen.
- Der Range fasst genau: 5 Adopt-Dateien + Plan + Roadmap-Ruhe-Marker (`c7091b7b`, reiner Move-Commit des Plans von `open/` nach `in-progress/`) + R1-Report. Keine Ausweitung über die deklarierte Abgrenzung hinaus.

## Beobachtungen (nicht blockierend, closure-seitig)

1. Die DoD-Haken in §2 sind offen — korrekt für den Zustand `in-progress` (der Haken-Wächter greift laut MR-056 erst unter `done/`); beim Closure-Zug mit abhaken.
2. §7 Closure-Notiz trägt den Nachtlauf-Stand (R1-F-3-Korrektur: „beide Nachtläufe grün (upstream-drift 2026-09-29, image-scan 2026-09-28)"); die übrigen Notiz-Felder stehen auf `—` — §5 verlangt das Füllen vor dem `git mv` nach `done/`.

## Verdikt

**DoD vollständig erfüllt.** Alle drei DoD-Punkte sind gegen eigenständig erhobene Belege bestätigt (Adopt-Dateien im Verzeichnis, Plan-Text nach R1-Korrektur, Vorlage-Diffs, frischer `make gates`-Lauf mit Exit 0). Der Slice ist DoD-seitig closure-fähig; offen sind nur die closure-seitigen Schritte (DoD-Haken, Volltext der §7-Notiz, `verify-closure-notes` beim `fullbuild`).
