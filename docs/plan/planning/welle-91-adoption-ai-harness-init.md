# Welle welle-91: Adoption aus ai-harness-init (Stand v6.13.0)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-91-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-29.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die Adoptions-Kandidaten aus dem Frisch-Bootstrap der Schwester
(ai-harness-init, Stand Kurs v6.13.0; Snapshot `/tmp/aih-v6.13.0`) sind je
Area bewertet und entschieden: die `.claude`-Rollen und Commands, die
Werkzeuge unter `tools/harness/` und die fünf über die Baseline-Vorlage
hinausgehenden `.d-check.yml`-Positionen. Je Area steht eine belegte
Entscheidung — Übernahme (an d-check angepasst) oder begründete Ablehnung.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf
erwähnt werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein
Ergebnis dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch
platziert.

- Der Auftraggeber hat die Wellen-Wiederaufnahme und die Evaluierung
  angefragt; der Snapshot `/tmp/aih-v6.13.0` liegt vor (2026-09-29).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle sind `done/`.
- Das Trigger-Audit (seit slice-241) ist über die Welle-Closure vollzogen
  und in `welle-91-results.md` belegt.
- `make fullbuild` grün.

## 4. Slices in dieser Welle

<!-- BEDIENHINWEIS: keine Status-Spalte ergaenzen. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-244 | `.claude`-Rollen und Commands evaluieren | Snapshot `/tmp/aih-v6.13.0` |
| slice-245 | `tools/harness`-Werkzeuge evaluieren | Snapshot `/tmp/aih-v6.13.0` |
| slice-246 | `.d-check.yml`-Positionen evaluieren | Snapshot `/tmp/aih-v6.13.0` |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine Welle.
- Wird blockiert von: keiner.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Die **Erfassungsschicht** der Schwester (`span-emit.sh`, `erfassung.mk`
  und die werkzeug-erzeugte Feldliste — LH-FA-10/Festlegung 5 des
  Schwester-Bootstrap): eigene Telemetrie-Entscheidung, wird hier nicht
  evaluiert.
- **Produkt-Release/Tag**: die Welle berührt keine Distributions-Fläche;
  Release-Prep bleibt eigener Vorgang.
- **Inhaltliche Baseline-Deltas**: die v6.13.0-Adoption ist abgeschlossen
  (slice-240/241/242/243); diese Welle bewertet ausschließlich Schwester-Artefakte.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: <Zeiger auf `welle-91-results.md`, Geschwister im Ruheort `done/`>
