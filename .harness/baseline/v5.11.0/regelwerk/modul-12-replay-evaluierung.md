## Modul 12 — Replay und Evaluierung

<!-- Quelle: [04-qualitaet/modul-12-replay-evaluierung.md](https://github.com/pt9912/ai-harness-course/blob/v5.11.0/kurs/de/04-qualitaet/modul-12-replay-evaluierung.md) -->

### Kernidee (Modul 12)

Ohne Replay ist jeder Lauf ein einmaliges Experiment. Mit Replay wird er
zur Messung — und damit zum Sensor, der merkt, wenn eine schnell
gemachte Änderung das *Verhalten* gedreht hat, nicht nur die Signatur.

**Wovon diese Regeln sprechen, wenn sie „Modell" sagen:** vom
*nicht-deterministischen Kern* deines Produkts — der Komponente, die auf
dieselbe Eingabe nicht garantiert dieselbe Ausgabe liefert. Im Regelfall
ist das dein **Domänen-Modell**: Simulation, Optimierer, Ranking,
Scoring — samt seiner Zufallsquelle und seiner Entscheidungsregeln. Ein
Inferenz-Modell (Embedding-Modell, LLM) ist derselbe Fall mit anderen
Drift-Quellen. Hat dein Repo beides, gilt das Verfahren für beide
getrennt.

### Regeln gegen typische Fehlannahmen (Modul 12)

- Replay grün heißt: das Modell hat das wiederholt, was *im Golden Set steht*. Ob das Golden Set noch die Realität abbildet, ist eine andere Frage.
- Statische Golden Sets überfitten. Rotation und neues Sampling sind Pflicht, nicht Kür.
- Determinismus erfordert, dass *jede* Quelle gepinnt ist: Inputs, die Zufallsquelle (Seed **und** ihre Ableitungsregel), die Modellversion, die Tool-Versionen, der Umgebungszustand des Containers (Zeit, Locale, Env-Variablen, Netz, sichtbare CPU-Zahl) und die Zeitstempel-Maskierung. Wer nur den Seed pinnt, pinnt eine *einzige* davon — die übrigen driften unabhängig weiter.

### Replay-Manifest (Modul 12)

Ein Baseline-Replay hält das Verhalten eines Laufs als Messung fest,
gegen die spätere Änderungen verglichen werden — gleich ob sie von
Hand, vom Agenten oder aus einer Abhängigkeit kommen.

#### Ziel-Form: Golden Set

**Ein Verzeichnis je Set**, üblicherweise unter `evals/golden/`. Der
Set-Name ist frei — er darf an der Closure hängen, die das Set erzeugt
hat, oder an dem, was es misst.

```
evals/golden/<set-name>/
├── manifest.yaml         # womit gelaufen wurde — Felder unten
├── inputs/
│   ├── case-001.json     # Happy
│   ├── case-002.json     # Boundary
│   └── case-003.json     # Negative
├── expectations/
│   ├── case-001.json     # gleichnamig — der Dateiname ist die Kopplung
│   ├── case-002.json
│   └── case-003.json
└── CHANGELOG.md          # datierte Veränderung des Sets
```

`inputs/<fall>` trägt die Eingabe des Laufs, `expectations/<fall>` die
Zusagen dazu (Form: siehe *Erwartungen als Verhalten* unten). Die
Zuordnung läuft über den gleichen Dateinamen.

Beide Seiten tragen über die Nutzlast hinaus drei Felder: eine `id`, die
Fehlerklasse (`kind`: happy · boundary · negative) und den `bezug` auf
das Akzeptanzkriterium, das dieser Fall bewacht. Der `bezug` ist die
Traceability-Zeile des Sets — ohne ihn ist nach ein paar Closures nicht
mehr erkennbar, welche Zusage an einem Fall hängt.

```json
{
  "id": "<fall-kennung>",
  "kind": "<happy | boundary | negative>",
  "bezug": "<Akzeptanzkriterium>",
  "input": { }
}
```

**Mindestens drei Fälle — Happy · Boundary · Negative** (dieselbe
Spec-Disziplin wie Akzeptanzkriterien, [Modul 3](modul-03-spec.md)).
Ein Replay mit einem Fall ist eine Demo.

Zu jedem Fall gehört **ein Auswahlkriterium in einem Satz** — welche
Fehlerklasse fängt genau dieser Fall, die die anderen zwei nicht fangen?
Hättest du nur die Fälle im Set, die eine Fehlerklasse nicht sehen, wäre der
Replay grün und die Realität rot. Genau das ist der Unterschied zwischen einem
Demo-Set und einem Golden Set.

#### Ziel-Form: manifest.yaml

Welche Felder ein Manifest braucht, richtet sich danach, **welche
Quellen im Lauf noch zufällig sind.** Die Art des Kerns sagt das im
Regelfall voraus — darum zwei Formen. Mischformen kommen vor (ein
Inferenz-Modell mit Seed, eine nachgelagerte Regel-Stufe); dann
entscheidet die Quellen-Frage, nicht das Etikett.

Domänen-Modell (Simulation, Optimierer, Ranking, Scoring):

```yaml
slice: <erzeugende Closure>    # Traceability zurück auf den Plan
recorded_at: <ISO-8601-Zeitstempel mit Zeitzone>
model:                         # die gemessene Stufe
  name: <Bezeichnung der Stufe>
  version: <Kennung der Stufe>
  seed: <Startwert der Zufallsquelle>
determinism:                   # Regeln ohne eigenes Versionsfeld
  seed_derivation: <wie aus dem Seed die Werte je Fall entstehen>
  <regel>: <gewählte Strategie>   # ein Schlüssel je weiterer Regel
runtime:
  image_hash: <sha256:...>
  env:                         # was der Image-Hash nicht pinnt
    tz: <Zeitzone>
    locale: <Locale>
    network: <none | erlaubte Ziele>
    cpu_count: <sichtbare CPU-Zahl>
  timestamp_masking: <Strategie>
inputs_ref: <Pfad relativ zu diesem Manifest>
expectations_ref: <Pfad relativ zu diesem Manifest>
```

Inferenz-Modell (Embedding-Modell, LLM): `seed:` entfällt, wenn die
Inferenz-API keinen anbietet, und `prompt_context:` kommt hinzu. Die
übrigen Felder stehen wie oben.

```yaml
slice: <auslösende Closure>
recorded_at: <ISO-8601-Zeitstempel mit Zeitzone>
model:
  name: <Modell-Kennung>       # ohne gleitenden Alias
  version: <Release-Snapshot>
prompt_context:                # die Regeln dieses Laufs
  system_prompt_hash: <sha256:...>
  tool_definitions_hash: <sha256:...>
  tool_order: [<werkzeug>, ...]
runtime:
  image_hash: <sha256:...>
  env: { tz: ..., locale: ..., network: ..., cpu_count: ... }
  timestamp_masking: <Strategie>
inputs_ref: <Pfad relativ zu diesem Manifest>
expectations_ref: <Pfad relativ zu diesem Manifest>
```

Im **reinen** Inferenz-Lauf steht `prompt_context:` an Stelle von
`determinism:` — dort gibt es keine Entscheidungsregeln ohne
Versionsfeld. Hat der Lauf beides, ein Modell *und* eine nachgelagerte
Regel-Stufe, stehen **beide Blöcke**: der Prompt-Kontext für das Modell,
`determinism:` für die Regeln der Stufe.

Was die einzelnen Felder festhalten:

| Feld | Was hineingehört | Wozu |
|---|---|---|
| `slice` | die Closure, die das Set **erzeugt** hat — spätere Änderungen stehen im `CHANGELOG.md`, nicht hier | Traceability zurück auf den Plan; wer das Feld fortschreibt, überschreibt den Anker bei jeder Rotation |
| `recorded_at` | Zeitpunkt der Aufnahme | damit spätere Läufe ihren Diff datieren können |
| `model.name` | Bezeichnung des Messgegenstands; beim Inferenz-Modell die Modell-Kennung **ohne gleitenden Alias** | ein Alias zeigt morgen auf etwas anderes |
| `model.version` | die konkrete Kennung, nicht die Familie. Beim Inferenz-Modell trägt sie die Hauptlast — aus demselben Grund, aus dem ein Image-Hash kein Tag ist | grenzt Modell-Drift ab (Rang 2 der Diagnose unten) |
| `model.seed` | Startwert der Zufallsquelle | pinnt den stochastischen Anteil des Laufs |
| `determinism:` | Entscheidungsregeln ohne eigenes Versionsfeld — Tie-Break, Sortierstabilität, Grenzwerte, dazu `seed_derivation` | pinnt, was sonst still driftet: Niemand bemerkt einen Tie-Break-Wechsel an einer Versionsnummer |
| `determinism.seed_derivation` | wie aus dem Seed die Werte je Fall entstehen — Substream je Fall, fortlaufender Strom, Neuinitialisierung | der Seed allein reicht nicht: zwei Läufe mit gleichem Seed und anderer Ableitung sind verschiedene Läufe |
| `prompt_context:` | System-Prompt und Werkzeug-Definitionen **als Hash**, dazu deren Reihenfolge | dasselbe beim Inferenz-Modell — als Hash fällt jede Änderung auf, ohne dass du den ganzen Prompt einfrieren musst |
| `runtime.image_hash` | byte-genaue Adresse des Container-Images, kein Tag | grenzt Toolchain-Drift ab (Rang 1) |
| `runtime.env` | Zeitzone, Locale, Netz-Zugang, sichtbare CPU-Zahl | der Image-Hash pinnt die Toolchain, nicht den Laufzeit-Zustand — der entsteht erst beim Start des Containers |
| `runtime.timestamp_masking` | wie Zeitstempel in der Ausgabe ersetzt werden | ohne sie wird ein Lauf schon an seinem eigenen Datum rot |
| `inputs_ref` | Verweis auf das Eingabe-Verzeichnis — referenziert, nicht als Inline-Text ins Manifest kopiert | sagt dem Replay, gegen welchen Datensatz er läuft |
| `expectations_ref` | Verweis auf das Erwartungs-Verzeichnis, ein Gegenstück je Eingabe-Fall | ohne Gegenstück gibt es nichts zu vergleichen |

**Drei Positionen** sind die Merkformel dafür, was den *Lauf*
rekonstruierbar macht — zwei feste Felder und eine Familie:
`inputs_ref`, `recorded_at` und **je ein Feld pro Zufallsquelle des
Laufs**. Wie viele Felder die dritte umfasst, entscheidet der Lauf: beim
Domänen-Modell `model.seed` *und* `determinism:`, beim Inferenz-Modell
`model.version` *und* `prompt_context:`. `model.seed` **entfällt**, wenn
die Inferenz-API keinen Seed-Parameter anbietet — ein Feld, dessen Wert
auf nichts zeigt, ist keine Pflicht, sondern Dekoration.

Die Merkformel zählt nicht die ganze Datei. `expectations_ref` gehört
ebenso hinein, steht aber auf der Vergleichs- statt auf der
Rekonstruktions-Seite; `slice` und `model.name` benennen Set und
Messgegenstand. `runtime.image_hash` und `model.version` liegen beim
Domänen-Modell außerhalb der Merkformel, unterscheiden aber ernsthaftes
von symbolischem Replay.

**Ins Manifest gehört nur, was der Lauf selbst noch tut.** Was
vorgelagert und eingefroren ist — Daten, die fertig in `inputs/` liegen
— ist keine Drift-Quelle mehr und bekommt kein Feld. Sonst behauptest du
eine Abhängigkeit, die im Replay gar nicht mehr wirkt.

#### Regeln

- **Erwartungen als Verhalten, nicht als Wortlaut:** semantische Anker
  statt wörtlichem Vergleich — die Feldnamen kommen aus deiner Domäne,
  die drei Anker-Arten nicht. Bei einem Domänen-Modell etwa der erwartete
  Spitzenwert (*welches Element steht oben*), eine Untergrenze (*eine
  Schwelle statt eines exakten Werts*) und eine Invariante (*was über
  zwei Läufe gleich bleibt*, etwa die Stabilität der Reihenfolge
  bei Gleichstand); bei einem Inferenz-Modell `must_include` ·
  `must_not_include` · `tool_calls`-Zähler. Eine wörtlich kopierte
  Ausgabe bricht, sobald sich ein Wert in der letzten Nachkommastelle
  ändert. Exact-Match nur für strukturierte Schnittstellen
  (JSON-Felder), nie für Fließtext.
- **Baseline einfrieren:** wird der erste Lauf nicht grün, *erst* das
  Manifest schärfen (meist Erwartung zu eng), nicht die Implementierung
  anfassen.
- **Drei Ergebnisse, drei Reaktionen** — nach einer Änderung am Kern:
  alle grün → kein Drift in dieser Klasse. Einer rot → erste
  Drift-Diagnose: ist die Erwartung zu eng (nachschärfen) oder hat der
  Kern ein neues Verhalten? Zwei rot → die Änderung ist nicht ohne
  Anpassung möglich; **Carveout plus Folge-Slice** für das
  Erwartungs-Update.
- **Drift als Zahl:** **Drift-Rate** = rote Fälle ÷ Gesamt-Fälle. Die
  Zahl macht *Trend* über mehrere Änderungen und eine *Schwelle* für den
  Steering Loop („ab Drift-Rate > X Carveout-Pflicht") prüfbar — eine
  ordinale „einer rot"-Notiz lässt sich zwischen Läufen nicht
  vergleichen. **Als Schwelle trägt die Rate aber erst, wenn das Set
  über Rotation gewachsen ist:** Bei den drei Pflicht-Fällen springt
  jeder einzelne rote Fall um 33 Punkte, und eine Schwelle wie „ab 10 %
  Carveout-Pflicht" löst dann bei jedem Rot aus. Bis dahin ist sie eine
  Notiz, keine Schwelle.
- **Drift-Diagnose in fester Reihenfolge** (wer zuerst „echte Regression"
  tippt, baut den Carveout an der falschen Stelle ein):

| Reihenfolge | Verdächtiger | Belegquelle |
|---|---|---|
| 1 | Toolchain-Drift | `runtime.image_hash` verglichen |
| 2 | Modell-Drift | `model.version` · `model.seed` · `determinism:` bzw. `prompt_context:` verglichen |
| 3 | Erwartungs-Drift | Eingaben vs. Spec (Modul 3) |
| 4 | echte Regression | alles oben ausgeschlossen |

Beim Inferenz-Modell kommt auf Rang 2 der **Provider-Status** hinzu:
gleiche Version, anderes Subroute ist eine reale Drift-Quelle.

- **Rotation:** Replay-Sets verrotten — Fälle aus Steering-Loop-Einträgen
  ergänzen, giftig gewordene (Schnittstelle real geändert) entfernen,
  datiert im Set-eigenen `CHANGELOG.md`. Je Eintrag: Datum, was sich am
  Set geändert hat, welche Closure es ausgelöst hat und **warum** ein
  Fall ergänzt oder entfernt wurde.

