# Review — slice-055 (Completeness-Rückbau) · R2 (Mechanik- / Gate-Korrektheit / fail-closed)

## Kopf-Metadaten

- **Rolle:** unabhängiger Reviewer R2 — Makefile-/Verdrahtungs-/fail-closed-Korrektheit. NICHT
  Doc-first/Prozess (separater Reviewer R1). Kein Verifier (Gate-Status als gegeben).
- **Datum:** 2026-06-29. **Reviewer-Skill:** `.harness/skills/reviewer.md` v1.2.0.
- **Gegenstand:** Working-Tree-Änderungen slice-055 (uncommitted): `Makefile`-Recipe
  `completeness-check` (Skript → `$(DCHECK_RUN) --trace --require-complete`), `git rm
  tools/completeness-check.sh`, `.d-check.yml` `codepaths.ignore-refs +=`, Gate-Tabellen in
  AGENTS/harness, ADR-0026 + Slice.

## Findings

### F-1 — LOW — `completeness-check` und `doc-complete`: identische Recipe ohne Identitäts-Wächter

- **Quelle:** Maintainability.
- **Pfad:** `Makefile` (`doc-complete` + `completeness-check`).
- **Befund:** Beide Targets tragen die byte-identische Recipe `$(DCHECK_RUN) --trace
  --require-complete` mit identischem `build`-Prerequisite. **Failure-Szenario:** eine künftige
  divergierende Edit (ein Flag nur an einem) bleibt unbemerkt — kein Gate erzwingt Recipe-Identität,
  und der Kommentar „Dieselbe Mechanik wie doc-complete" würde still zur Unwahrheit, während
  Closure-Gate und Konsumenten-Dogfood auseinanderlaufen.
- **Verifizierbar:** nein (kein Gate prüft Recipe-Gleichheit; per Diff sichtbar).

### F-2 — INFO — Leerer Scan-Scope ⇒ Exit 0 (silent-green), vorbestehend

- **Quelle:** DC-FA-CLI-011 / ADR-0026.
- **Pfad:** `internal/adapter/driving/cli/cli.go` (`requireComplete && matrix.Orphans > 0`) +
  `core/app/trace.go`.
- **Befund:** `--require-complete` bindet Exit 1 allein an `matrix.Orphans > 0`. Fehlt
  `spec/lastenheft.md` oder verlieren alle Anforderungs-Headings ihre Kennungs-Form, liefert die
  RTM `Total=0/Orphans=0` ⇒ Exit 0 — das Gate meldet „vollständig", ohne etwas bewiesen zu haben.
  **Pre-existing:** das gelöschte Skript bestand `{"total":0,"orphans":0}` ebenso grün (kein
  `total>0`-Boden). slice-055 ist darauf **neutral**.
- **Verifizierbar:** ja (`make completeness-check` gegen einen Baum ohne `spec/lastenheft.md` → Exit 0).

## Negativbefunde (geprüft, ohne Blocker)

- **Recipe-Korrektheit:** `make completeness-check` ruft exakt `$(DCHECK_RUN) --trace
  --require-complete`; `DCHECK_RUN` netzlos (`--network none`), read-only (`:ro`), gegen
  `$(IMAGE):latest`. `build`-Prerequisite vorhanden. Bindepunkt intakt (`fullbuild: ci bench
  completeness-check`); NICHT in `gates`/`ci`.
- **Fail-closed-Parität — REFUTED als Regression:** die Hypothese „der Flag verliert eine
  fail-closed-Garantie des Skripts" ist falsch. (a) Image fehlt/Lauf bricht: das Skript fing das per
  `[ -z "$json" ]` → Exit 1; die neue Recipe propagiert den realen Nicht-Null-Exit von
  `docker run`/d-check (kein `|| true`, kein `2>/dev/null`) → make failt — Parität, transparenter.
  (b) JSON-Schema-Drift-Absicherung war gegen die grep/awk-Fragilität des Skripts selbst; in-Produkt
  ist `matrix.Orphans` ein Prozess-Integer — der Versagensmodus existiert nicht mehr. (c) Der
  „0 Waisen bei leerem Scope ⇒ grün"-Pfad ist **kein neues** Silent-Green (Skript identisch) → daher
  F-2 nur INFO.
- **Ketten-Duplikat (Rollentrennung):** `doc-complete` (Konsumenten-/print-mk-Name) vs.
  `completeness-check` (Producer-Closure-Gate) ist eine dokumentierte, vertretbare Rollentrennung;
  das Finding betrifft nur die unverdrahtete Recipe-Identität (F-1).
- **gate-consistency intakt:** Target-Name `completeness-check` unverändert, weiter in `.PHONY` und
  in AGENTS §4 **und** harness/README §Sensors gelistet. `tools/gate-consistency.sh` parst nur
  Target-Namen + die `modules:`-Zeile — Recipe/Skript-Pfad gleichgültig; Modulliste unangetastet.
- **Tombstone-Glob:** `ignored()`→`matchGlob()`→`path.Match`; `tools/completeness-check.sh` matcht
  den aufgelösten Pfad exakt, referenz-weit (Skip vor allen Prüfungen) — deckt alle Live-Doc-
  Vorkommen (immutable ADR-0017, ADR-0026, slice-055, roadmap).
- **Kein Produkt-Code/Hexagon:** `git status` zeigt keine `internal/`-/`cmd/`-Änderung; Image
  byte-identisch; kein neuer Import, kein arch-check-Belang.
- **Selbsttest-Kompensation:** der Skript-Selbsttest ist durch in-Produkt-Akzeptanztests ersetzt,
  die in `make test` (Teil von `gates`, läuft immer — strenger als der alte fullbuild-only-Selbsttest)
  feuern: `TestCLI044_RequireComplete_KeineWaisen/_Waise/_JSON/_OhneTrace` (beide Richtungen +
  Usage-Negative). Einzig der Leer-Scope-Boundary bleibt ungetestet (= F-2).

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO | REFUTED |
|---|---|---|---|---|
| 0 | 0 | 1 (F-1) | 1 (F-2) | 1 (fail-closed-Regression) |

## Verdikt: MERGEBAR (aus R2-Sicht)

Keine HIGH/MEDIUM: Recipe, Bindepunkt, Netzlosigkeit, `build`-Prerequisite, gate-consistency und
Tombstone-Glob korrekt verdrahtet; keine fail-closed-Regression durch Skript→Flag (REFUTED mit
Skript-Zitat). F-1 (Recipe-Duplikat ohne Wächter) und F-2 (vorbestehendes Leer-Scope-Silent-Green)
nicht blockierend.

## Einarbeitung (Implementation, 2026-06-29)

- **F-1 — behoben:** der Flag-Satz lebt jetzt in **einer** geteilten Makefile-Variable
  `COMPLETE_FLAGS = --trace --require-complete`; `doc-complete` und `completeness-check` rufen beide
  `$(DCHECK_RUN) $(COMPLETE_FLAGS)` → keine stille Divergenz der Flags mehr möglich.
- **F-2 — dokumentiert:** slice-055 §4 hält das (vorbestehende, von slice-055 nicht eingeführte)
  Leer-Scope-Silent-Green fest; ADR-0026 Re-Evaluierungs-Trigger nennt einen `total>0`-Boden als
  DC-FA-CLI-011-Erweiterung (eigener Trigger, Produkt-Change, bewusst nicht hier).
