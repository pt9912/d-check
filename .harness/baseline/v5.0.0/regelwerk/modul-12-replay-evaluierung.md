## Modul 12 — Replay und Evaluierung

<!-- Quelle: [04-qualitaet/modul-12-replay-evaluierung.md](https://github.com/pt9912/ai-harness-course/blob/v5.0.0/kurs/de/04-qualitaet/modul-12-replay-evaluierung.md) -->

### Kernidee (Modul 12)

Ohne Replay ist jeder Agenten-Lauf ein einmaliges Experiment. Mit Replay
wird er zur Messung.

### Regeln gegen typische Fehlannahmen (Modul 12)

- Replay grün heißt: das Modell hat das wiederholt, was *im Golden Set steht*. Ob das Golden Set noch die Realität abbildet, ist eine andere Frage.
- Statische Golden Sets überfitten. Rotation und neues Sampling sind Pflicht, nicht Kür.
- Determinismus erfordert: Modellversion + Seed + Inputs *und* Tool-Versionen, Wetter im Container, Zeitstempel-Maskierung. Wer nur den Seed pinnt, pinnt eine *einzige* von mehreren Drift-Quellen — Modellversion, Sampling-Parameter, Tool-Umgebung und Prompt-Kontext driften unabhängig davon weiter.

### Replay-Manifest (Modul 12)

Ein Baseline-Replay hält einen Agentenlauf als Messung fest, gegen die
Modellwechsel verglichen werden. Layout `evals/golden/welle-NN-baseline/`
mit `manifest.yaml`, `inputs/`, `expectations/`. Regeln:

- **Mindestens drei Fälle — Happy · Boundary · Negative** (dieselbe
  Spec-Disziplin wie Akzeptanzkriterien, [Modul 3](modul-03-spec.md)).
  Ein Replay mit einem Fall ist eine Demo.
- **Manifest-Pflichtfelder:** `model.version`, `model.seed`, `inputs_ref`
  (Selbstcheck-Pflicht); dazu `runtime.image_hash` (Toolchain-Drift
  abgrenzen) und `recorded_at` (späteren Diff datieren) — sie trennen
  ernsthaftes von symbolischem Replay.
- **Erwartungen als Verhalten, nicht als Wortlaut:** `must_include` ·
  `must_not_include` · `tool_calls`-Zähler statt wörtlichem Vergleich
  (der bricht bei Modellwechsel sofort). Exact-Match nur für
  strukturierte Schnittstellen (JSON-Felder), nie für Fließtext.
- **Baseline einfrieren:** wird der erste Lauf nicht grün, *erst* das
  Manifest schärfen (meist Erwartung zu eng), nicht das Modell tauschen.
- **Drift als Zahl:** **Drift-Rate** = rote Fälle ÷ Gesamt-Fälle. Die
  Zahl macht *Trend* über Modellversionen und eine *Schwelle* für den
  Steering Loop („ab Drift-Rate > X Carveout-Pflicht") prüfbar — eine
  „zwei rot"-Notiz lässt sich zwischen Läufen nicht vergleichen.
- **Drift-Diagnose in fester Reihenfolge** (wer zuerst „echte Regression"
  tippt, baut den Carveout an der falschen Stelle ein):

| Reihenfolge | Verdächtiger | Belegquelle |
|---|---|---|
| 1 | Toolchain-Drift | `runtime.image_hash` verglichen |
| 2 | Modell-Routing | `model.version` plus Provider-Status |
| 3 | Erwartungs-Drift | Eingaben vs. Spec (Modul 3) |
| 4 | echte Regression | alles oben ausgeschlossen |

- **Rotation:** Replay-Sets verrotten — Fälle aus Steering-Loop-Einträgen
  ergänzen, giftig gewordene (Schnittstelle real geändert) entfernen,
  datiert im Set-eigenen `CHANGELOG.md`.

