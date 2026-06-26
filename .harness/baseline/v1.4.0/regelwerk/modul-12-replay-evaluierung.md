## Modul 12 — Replay und Evaluierung

*Quelle: [04-qualitaet/modul-12-replay-evaluierung.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-12-replay-evaluierung.md)*

### Kernidee (Modul 12)

Ohne Replay ist jeder Agenten-Lauf ein einmaliges Experiment. Mit Replay
wird er zur Messung.

### Regeln gegen typische Fehlannahmen (Modul 12)

- Replay grün heißt: das Modell hat das wiederholt, was *im Golden Set steht*. Ob das Golden Set noch die Realität abbildet, ist eine andere Frage.
- Statische Golden Sets überfitten. Rotation und neues Sampling sind Pflicht, nicht Kür.
- Determinismus erfordert: Modellversion + Seed + Inputs *und* Tool-Versionen, Wetter im Container, Zeitstempel-Maskierung. Wer nur den Seed pinnt, pinnt eine *einzige* von mehreren Drift-Quellen — Modellversion, Sampling-Parameter, Tool-Umgebung und Prompt-Kontext driften unabhängig davon weiter.

### Worked Example: ein Replay-Manifest aufbauen

**Ausgangssituation:** Du hast einen Agentenlauf gemacht, der den Slice
`SL-024` (kleiner Replay-Erweiterung) abgeschlossen hat. Du willst ihn
als *Baseline-Replay* festhalten, gegen den Modellwechsel verglichen
wird.

**Schritt 1 — Pfad und Skelett anlegen.**

```
evals/golden/welle-1-baseline/
├── manifest.yaml
├── inputs/
│   ├── case-001.json
│   ├── case-002.json
│   └── case-003.json
└── expectations/
    ├── case-001.json
    └── ...
```

Drei Fälle ist das Minimum: Happy / Boundary / Negative — dieselbe
Spec-Disziplin wie bei Akzeptanzkriterien
([Modul 3](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/01-spec-und-architektur/modul-03-lastenheft.md)). Ein
Replay mit einem Fall ist eine Demo, kein Replay.

**Schritt 2 — Pflichtfelder im Manifest fixieren.**

```yaml
# evals/golden/welle-1-baseline/manifest.yaml
slice: SL-024
recorded_at: 2026-06-15T10:31:00Z
model:
  name: claude-opus-4-7
  version: "20260301"
  seed: 42
runtime:
  image_hash: sha256:9c7f4a...   # siehe Vorgriff-Block oben
  toolchain:
    python: 3.12.4
    ruff: 0.4.10
inputs_ref: inputs/
expectations_ref: expectations/
```

Drei Felder sind im Selbstcheck Pflicht: `model.version`, `model.seed`,
`inputs_ref`. Zwei weitere unterscheiden ernsthaftes von symbolischem
Replay: `runtime.image_hash` (Toolchain-Drift abgrenzen) und
`recorded_at` (späteren Diff datieren).

**Schritt 3 — Erwartungen *als Verhalten*, nicht als Wortlaut.**
Schlecht: *"Agent antwortet exakt 'index updated'"* — bricht bei
Modellwechsel sofort. Gut:

```yaml
# expectations/case-001.json
{
  "must_include": ["index", "updated"],
  "must_not_include": ["error", "traceback"],
  "tool_calls": {
    "writer.write_index": {"min": 1, "max": 1}
  }
}
```

Drei semantische Aussagen statt eines wörtlichen Vergleichs. Exact-Match
bewahre für strukturierte Schnittstellen (JSON-Felder), nie für
Fließtext.

**Schritt 4 — Erster Lauf, Baseline einfrieren.**

```bash
make replay RUN=welle-1-baseline
```

Erwartet: drei grüne Fälle. Wenn nicht: *erst* das Manifest schärfen
(meist Schritt 3 zu eng), nicht das Modell tauschen.

**Schritt 5 — Modellwechsel-Drift messen.**

```bash
make replay RUN=welle-1-baseline MODEL=claude-sonnet-4-6
```

Drei mögliche Ergebnisse:
* alle grün → kein Drift in dieser Klasse.
* einer rot → erste Drift-Diagnose: ist die Erwartung zu eng (Schritt 3 nachschärfen)
  oder hat das neue Modell ein neues Verhalten?
* zwei rot → Modellwechsel nicht ohne Anpassung möglich; Carveout +
  Folge-Slice für Erwartungs-Update.

*Quantifizieren statt nur einordnen.* Halte den Drift als **Zahl** fest,
nicht nur als "ein paar rot": die **Drift-Rate** = rote Fälle ÷
Gesamt-Fälle des Golden Sets (zwei von zwanzig → 10 %). Die Zahl macht
zweierlei prüfbar, was die ordinale Einordnung verbirgt: (1) den *Trend*
über mehrere Modellversionen (steigt die Rate von 5 % auf 15 %, ist der
Modellpfad selbst der Verdächtige, nicht ein Einzelfall), und (2) eine
*Schwelle* für den Steering Loop ("ab Drift-Rate > X Carveout-Pflicht").
Eine reine "zwei rot"-Notiz lässt sich zwischen Läufen nicht vergleichen
— ein Prozentwert schon.

**Schritt 6 — Drift-Diagnose-Reihenfolge.** Wenn ein Lauf rot wird, ist
die Reihenfolge der Verdächtigen *nicht beliebig*:

| Reihenfolge | Verdächtiger | Belegquelle |
|---|---|---|
| 1 | Toolchain-Drift | `runtime.image_hash` verglichen |
| 2 | Modell-Routing | `model.version` plus Provider-Status |
| 3 | Erwartungs-Drift | Eingaben vs. Spec (Modul 3) |
| 4 | echte Regression | alles oben ausgeschlossen |

Wer zuerst auf "echte Regression" tippt, baut den Carveout an der
falschen Stelle ein.

**Schritt 7 — Lerneintrag und Rotation.**
Replay-Sets verrotten. In
`evals/golden/welle-1-baseline/CHANGELOG.md`:

```markdown
2026-06-15 — Baseline mit drei Fällen aufgesetzt.
2026-08-02 — Fall hinzugefügt aus Steering-Loop-Eintrag #4
             (vorher unerkanntes Negativ-Muster, siehe
             reflexion-vorlage.md).
2026-09-10 — Fall-001 entfernt — Realität hat
             Schnittstelle geändert, Fall war giftig.
```

Sieben Schritte, ein reproduzierbares Manifest. Vergleich im Lab:
[`../../lab/example/evals/golden/welle-1-baseline/`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/evals/golden/welle-1-baseline/)
mit `manifest.yaml`, `inputs/case-{001,002,003}.json`,
`expectations/case-{001,002,003}.json` und `CHANGELOG.md` in derselben
Verzeichnis-Struktur. Das Lab zeigt ein *Retrieval*-Replay (Embedding-
Modell `local-embed-v3`, drei Search-Cases gegen LH-FA-02); das Worked
Example oben demonstriert dasselbe Schema für einen *LLM-Agentenlauf*
— die Struktur trägt beides.

