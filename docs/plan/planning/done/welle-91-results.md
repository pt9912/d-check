# Welle welle-91 — Adoption aus ai-harness-init — Closure-Notiz

**Welle:** welle-91
**Abschluss:** 2026-09-29
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — Ergebnis, nicht Tätigkeit.

- **slice-244** (`.claude`-Rollen und Commands evaluieren): 5 adoptiert
  (architect, planner, validator als Agents; plan-welle, close-welle als
  Commands), 4 abgelehnt — d-check führt die Rolle in eigener Form.
- **slice-245** (`tools/harness`-Werkzeuge evaluieren): 3 adoptiert
  (slice-mv, history-range-guard als Vorlauf-Wächter vor trace-check/adr-check
  und als Einzellauf-Target, selbstpruefung), 2 abgelehnt (e2e-abdeckung,
  traeger-fetch — je belegt); slice-mv erhielt Identity-Fallback und die
  DONE-Unterordner-Anpassung mit Messung.
- **slice-246** (`.d-check.yml`-Positionen evaluieren): Position welle-Klasse
  bereits Bestand in schärferer Form
  ([MR-034](../../../../harness/conventions.md#mr-034--die-referenzmatrix-bewacht-auch-die-kante-adr--welle),
  Token `welle-\d{2,}`, beide Regeln); Positionen `aussen`/`adaptionsblock`
  Grundsatz bejaht, Umfang delegiert an slice-248 (Probe: 40 Befunde);
  Positionen ids-Weite und exclude-sections-Scoping abgelehnt, je gemessen.

## Was hat funktioniert?

- **Je Werkzeug/Position eine Einzelentscheidung statt
  Pauschalübernahme** — jede adoptierte oder abgelehnte Position trägt ihre
  Messung (scratch-Repo, shallow-Clone, grep-Zählung), nicht eine Behauptung.
- **Die Empirie-Route vor der Doku**: Verhaltens-Behauptungen
  (Identity-Fallback, Unterordner-Mapping, Wächter-Exit-Codes, leere Range)
  wurden gemessen, bevor sie in Doku oder Folge-Plan standen — die
  shallow-Clone-Gegenprobe machte die stille-Grün-Lücke am eigenen Adapter
  sichtbar und schnitt slice-247.
- **Die zwei-Slice-Schneidung** (Werkzeug adoptieren, dann operieren) hielt:
  die Probe von slice-246 lief, bevor eine Entscheidung fiel; der Ersteinsatz
  von slice-mv (Flachfall-Claim von slice-246) ging fehlerfrei durch.

## Was ging anders als geplant?

- **Das adoptierte Werkzeug war unvollständig angepasst**: slice-mvs
  ausgehende Verweis-Richtung zog die Pfadtiefe des Unterordners nicht mit
  (R1-F-1, Repro des Reviewers) — behoben in 09aeb4d5, bevor der erste
  produktive Einsatz; der Closure-Unterordner-Fall blieb manuell (Grenze 5
  im Skriptkopf).
- **Die README-Werkzeug-Zeilen behaupteten ein Target ohne Makefile-Regel**
  (gate-phantom) — gefunden im repo-weiten Handoff-Lauf, behoben in
  8637f08a.
- **Die Zählung aus dem Gedächtnis war falsch**: die Probe lief 40 Befunde,
  notiert waren 29 — ein `tail`-Artefakt, das der Reviewer durch
  Probe-Wiederholung aufdeckte. Konsequenz: die Verifier-Rolle wiederholt
  Messungen statt sie zu übernehmen (Modul 11, dieselbe Linie wie
  „Bewusstes Brechen").

## Steering-Loop-Einträge

Kein Register-Eintrag erreicht in dieser Welle die 3×-Schwelle — es gibt
keine Verkörperung. Gezählt, nicht verkörpert: die stille-Grün-Klasse
(`BEO-ALL/stilles-gruen-ueber-leerer-range`, neu angelegt, 1×) und die
Kommentar-Klasse (`BEO-ALL/kommentar-traegt-herkunfts-prosa-statt-fuenf-
klassen`, auf 2× gestiegen). Beide tragen ihren Ausgang in den Folge-Slices
bzw. der individuellen Behebung, nicht in einer Regel.

## Beobachtungs-Register (Zeiger)

Der Zähler steht im stehenden Register `../observations/` — diese Notiz
trägt keine Daten. Was in dieser Welle 3× erreicht hätte, steht oben unter
*Steering-Loop-Einträge*: nichts.

## Folge-Slices

- slice-247 — vcs-Modul meldet stilles Grün über leerer, auflösbare Range
  (wellenlos, in `open/`).
- slice-248 — matrix-Klassen `aussen` + `adaptionsblock` adoptieren, samt
  Bestands-Nachzug (wellenlos, in `open/`).

## Verifikation

`make gates` Exit 0 an 93a3c216 (alle zehn Glieder),
`make verify-closure-notes` Exit 0, `make fullbuild` Exit 0 an demselben
Stand — Image-Hash sha256:4ac70ae439573dabc8a2d50a271e16de75ff083816e0574e3fd54f21658d99ac.
Der Trigger-Audit (Schritt 2) ist je Klasse geprüft: Carveout — keine aktiven;
bootstrap-aware Gate — coverage-gate kalibriert, kein Hochschalt-Trigger
gefeuert; ADR — 76 Trigger-Sektionen, keine Bedingung feuerte in der Welle;
Hard Rule — permanente Trigger bzw. benannte ungedeckte Kategorien, nichts
feuerte.
