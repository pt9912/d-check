# Review — slice-055 (Completeness-Rückbau) · R1 (Doc-first-Kohärenz / Harness-Konformität / Ehrlichkeit)

## Kopf-Metadaten

- **Rolle:** unabhängiger Reviewer R1 — Doc-first-Kohärenz, Harness-Konformität, Ehrlichkeit.
  Kein Verifier (Gate-Läufe nicht meine Rolle: `make ci`/`completeness-check`/`gate-consistency`
  grün; adversariale Waisen-Probe trieb `completeness-check` rot mit `WAISE`-Zeile + Anzahl,
  dann revertiert).
- **Datum:** 2026-06-29. **Reviewer-Skill:** `.harness/skills/reviewer.md` v1.2.0.
- **Gegenstand:** Working-Tree-Änderungen slice-055 (uncommitted) + lokaler Commit `7ac74da`
  (Backlog welle-44/45/46).
- **Eingangs-Kontext:** ADR-0026 (Proposed), ADR-0017 (Accepted/immutable), DC-FA-CLI-011/009,
  AGENTS §3.5.

## Findings

### F-1 — MEDIUM — §3d annotiert die ADR-0017-Index-Zeile nicht (Teil-Supersede)

- **Quelle:** ADR-Index-Konvention (README) + Präzedenz ADR-0016/0024.
- **Pfad:** `slice-055-completeness-rueckbau.md` §3d.
- **Befund:** Die Closure-Checkliste listet „neue ADR → Accepted + Geschichte-Append an ADR-0017
  + welle-44 entfernen", nennt aber **nicht** das Annotieren der ADR-0017-**Index**-Zeile mit der
  Teil-Supersede-Notiz. Die Konvention (Index: „Status der alten wird Superseded by ADR-NNNN")
  und beide Präzedenzfälle führen die abgelöste ADR im Index an: ADR-0016 = „Accepted
  (Skript-Mechanik teil-superseded durch ADR-0024)", ADR-0024 = „Accepted (…superseded durch
  ADR-0025)". **Failure-Szenario:** §3d wörtlich abgearbeitet → ADR-0017 bleibt blankes
  „Accepted", inkonsistent zur eigenen Konvention und reproduziert genau die Index-Status-Lücke,
  die slice-055 für ADR-0025 gerade aufräumt. Kein Gate erzwingt es (gate-consistency misst nur
  Target-Existenz).
- **Verifizierbar:** ja (nach Closure: `grep ADR-0017 README.md` blankes „Accepted", während
  0016/0024 die Notiz tragen).

### F-2 — LOW — ADR-0025-Index-Flip ohne Scope-Spur; wiederkehrende Klasse (Steering)

- **Quelle:** Maintainability/Provenance.
- **Pfad:** `docs/plan/adr/README.md` (ADR-0025 Proposed→Accepted).
- **Befund:** slice-055 flippt nebenbei die ADR-0025-Index-Zeile auf „Accepted" (Korrektur eines
  slice-054-Closure-Miss; die ADR-0025-**Datei** ist Accepted). Sachlich korrekt, aber weder im
  Slice-Plan noch in ADR-0026 als Scope erwähnt. **Failure-Szenario:** der Flip landet im
  Feat-Commit ohne Begründung; ein Audit „warum wechselt ADR-0025 in einem Completeness-Slice auf
  Accepted?" findet keine Spur. F-1+F-2 sind dieselbe Index-Status-Honesty-Klasse zum zweiten Mal
  → **Steering-Loop-Signal:** die Closure-Checkliste sollte „abgelöste/abhängige ADR-Index-Zeile
  mitziehen" mechanisieren statt es pro Slice manuell zu treffen.
- **Verifizierbar:** ja (`git diff HEAD README.md` zeigt den Flip; keine Erwähnung im Slice/ADR).

## Negativbefunde (geprüft, ohne Befund)

- **Supersede-Korrektheit:** ADR-0026 scopt sauber **nur** die „Skript-als-Gate-Quelle"-Mechanik
  (ADR-0017 §Entscheidung Pkt. 2); Policy „Waise→Exit 1" + Bindepunkt „fullbuild, nicht gates/ci"
  explizit unverändert. ADR-0017 §Entscheidung Pkt. 1/3/4/5 gewahrt. **AGENTS §3.5 gewahrt:**
  ADR-0017-Datei unangetastet (nicht in `git status`); nur Geschichte-Append zur Closure geplant.
- **Negativ-Selbsttest-Äquivalenz:** die ADR-0017-Garantie „fail-closed" bleibt strukturell — die
  fragile Bash-Parsing-Schicht entfällt ersatzlos, `matrix.Orphans`→Exit kommt aus dem
  akzeptanz-getesteten Produkt, ein gescheiterter `docker run` lässt das Target rot.
- **ADR-0025-Index ↔ Datei:** nach dem Flip konsistent (Datei Accepted ↔ Index Accepted).
- **Gate-Doku-Ehrlichkeit:** AGENTS §4 + harness/README §Sensors beschreiben die neue Mechanik
  faktisch korrekt (`--trace --require-complete`, `WAISE`-Zeilen via `report.go` `status="WAISE"`
  im Default-Markdown, stderr-Anzahl, netzlos/read-only). Der tote `tools/completeness-check.sh`-
  Pfad ist aus den lebenden Gate-Doku-Zeilen entfernt; Restnennungen sind immutable ADR-0017
  (Tombstone-gedeckt), ADR-0026/slice-055 (Beschreibung), done/-Slices + Reviews (Historie).
- **„Kein Release":** `git diff HEAD -- internal/ cmd/` leer → kein Produkt-Code; Image
  byte-identisch zu v0.34.0, kein Versions-Bump/GHCR korrekt.
- **Zweiter Tombstone:** `.d-check.yml` `ignore-refs += tools/completeness-check.sh` mit
  attribuierendem Kommentar; referenz-weit, deckt die immutable ADR-0017-Inline-Referenz + die
  done/slice-042-Nennung. Ehrlicher slice-054-Dogfood.
- **Referenz-Richtung (SDP)/Doc-Rollen:** ADR-0026 trägt keinen Provenance-Marker; slice-055-
  Nennungen sind Provenance, keine getarnte Entscheidungsgrundlage. AGENTS §4 knapp-operativ,
  harness/README Sensor-Tiefe, ADR Rationale, Slice Plan — Rollen gewahrt.
- **Backlog 7ac74da:** welle-45/46 nennen generisch „neue Anforderung + neue ADR" — keine
  konkreten Zukunfts-DC-IDs (keine Waisen); Nicht-Kandidaten-Notiz plausibel; welle-44 korrekt aus
  §Nächste Wellen entfernt (jetzt aktiv).
- **Lastenheft-Zitat-Treue:** ADR-0026 zitiert „…ohne die `completeness-check.sh`-Parsing-Logik zu
  kopieren" deckungsgleich mit der Lastenheft-0.24.0-Historie.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 1 (F-1) | 1 (F-2) | 0 |

## Verdikt: NACHBESSERN (ein MEDIUM vor Closure zu klären)

Kein HIGH; Doc-first-/Supersede-/Ehrlichkeits-Kern sauber und gegen die Quelle belegt. F-1
(MEDIUM): §3d muss die ADR-0017-Index-Annotation aufnehmen. F-1+F-2 = wiederkehrende
Index-Status-Honesty-Klasse → Steering an die Closure-Checkliste.

## Einarbeitung (Implementation, 2026-06-29)

- **F-1 — behoben:** slice-055 §3d ergänzt um „die ADR-0017-Index-Zeile mit der Teil-Supersede-
  Notiz auf die neue ADR annotieren (Form wie ADR-0016/ADR-0024)".
- **F-2 — behoben + Steering aufgenommen:** §3a-Artefaktliste nennt jetzt die ADR-0025-Index-
  Korrektur als Scope; §4 hält die wiederkehrende Index-Status-Honesty-Klasse fest und notiert die
  **Mechanisierungs-Idee** (ein Check „ADR-Datei-Status ↔ Index-Zeilen-Status konsistent", selbst
  d-check-nah) als Folge-Idee — nicht in diesem Slice gebaut.
