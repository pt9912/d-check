## Modul 13 — Quality Gates

<!-- Quelle: [04-qualitaet/modul-13-quality-gates.md](https://github.com/pt9912/ai-harness-course/blob/v5.9.0/kurs/de/04-qualitaet/modul-13-quality-gates.md) -->

### Harness-Einordnung (Modul 13)

Gates = *computational feedback* (siehe
[`grundlagen/klassifikation.md`](grundlagen-klassifikation.md)).
Schnellste und billigste Sensoren des Harness. Was hier prüfbar wird,
muss nicht mehr im Review-Agent landen — das ist die wichtigste
Einsparung im gesamten System.

### Kernidee (Modul 13)

Gates sind Aussagen, die *immer* gelten müssen. Wenn ein Gate "manchmal"
rot sein darf, ist es kein Gate, sondern ein Vorschlag.

### Gate-Typ ↔ Fehlerbild

Wer einen neuen Sensor in den Steering Loop einzieht, muss wissen,
*welche Sensor-Klasse welche Fehlerklasse fängt* — sonst reagiert er
auf einen wiederkehrenden Fehler mit dem falschen Sensor, und der
Steering Loop läuft leer. Die Zuordnung in Kurzform:

| Gate-Typ | typisches Fehlerbild | was er NICHT fängt |
|---|---|---|
| Linter | lokale Muster: toter Import, verbotenes Idiom, Suppression-Marker | Datenfluss über Funktionsgrenzen, Struktur-Regeln |
| Typecheck | Typgrenzen-Verstoß: falsche Signatur, `None` am falschen Ort | Vertrauensgrenzen — `str` bleibt `str`, ob nutzerkontrolliert oder nicht |
| Architekturtest | Struktur-/Import-Regel: Layer-Bruch, Domäne importiert Infrastruktur | Verhalten zur Laufzeit, lokale Muster |
| Security-Gate | Datenfluss-Befund: SQL-Injection, Secret-/Entropie-Treffer | Architektur-Schnitt, Coverage-Lücken |
| Coverage / Critical Coverage | Coverage-Loch — gesamt bzw. auf dem kritischen Pfad | Qualität der Tests, Spec-Lücken ([Modul 11](modul-11-verification.md)) |
| Replay-/Determinism-Gate | nicht-deterministischer Test oder Lauf | semantische Drift außerhalb des Golden Sets ([Modul 12](modul-12-replay-evaluierung.md)) |
| Integrationstest | Verhalten im Zusammenspiel: Komponenten-Vertrag bricht erst in Kombination | lokale Muster und Typgrenzen — dafür zu teuer und zu spät |

Trennlinie ist die *Regel-Klasse*, nicht das Tool: Linter machen lokale
Mustererkennung, Security-Regeln verlangen Datenfluss-Analyse,
Architekturtests prüfen Struktur, Integrationstests Verhalten im
Zusammenspiel.

### Hard Rule (Doku-Disziplin)

In `harness/README.md` und in jeder Doku, die Gates aufzählt: keine
Befehle behaupten, die es nicht gibt. Wenn `make fullbuild` strukturell
rot ist, wird das als Carveout in `docs/plan/carveouts/CO-<NNN>-…`
dokumentiert ([Modul 7](modul-07-carveouts.md)) und in
der Bindung-Spalte der Sensors-Tabelle per `CO-<NNN>`-ID verlinkt — nicht
ausgelassen, nicht geschönt, nicht in einer Status-Spalte versteckt
(die Sensors-Tabelle trägt keinen Lauf-Status; Lauf-Wahrheit pro Commit
liegt in CI, siehe
[`grundlagen-harness-dateien.md`](grundlagen-harness-dateien.md#harnessreadmemd-als-einstiegspunkt)).
**Eine neue Hard Rule trägt ab ihrer Einführung einen Auflösungs-Trigger oder
die Kennzeichnung *permanent*** — dieselbe Disziplin, die ADR
(Re-Evaluierungs-Trigger, [Modul 4](modul-04-adrs.md)) und Carveout
([Modul 7](modul-07-carveouts.md)) längst tragen. Für den **Altbestand** gilt:
kein Nachrüsten — ein nachgetragener Trigger wäre erfunden, nicht
rekonstruiert; der leere Zustand ist die ehrliche Information. Deklarierter
Backfill bleibt möglich, wo sich der Trigger wirklich herleiten lässt.

Halluzinierte Gates sind die häufigste Form von Harness-Lüge — und der
Implementer-Agent vertraut ihnen.

**Vorhanden ≠ behauptet.** Verboten ist ein *behauptetes* Gate ohne Deckung,
nicht ein *vorhandenes* Target ohne Anspruch. Ein tool-generiertes Gate-Fragment
(`d-check.mk` aus `d-check --print-mk`, per `-include` eingebunden statt
handgeschrieben) bringt mehr Targets mit, als du als Gate führst: nur das genutzte
(`docs-check`) steht in `harness/README.md`/`AGENTS.md` und `make gates`; die
advisory-Targets (`doc-trace`, `doc-doctor`, …) sind **verfügbar, aber nicht als
Gate behauptet** — wie ein Maintenance-Target (`regelwerk-check`), das nicht in
`gates` läuft. Die Lüge wäre, ein Gate zu versprechen, das nicht läuft.

### Bootstrap-aware Gates

In der Frühphase eines Projekts ist eine harte Coverage-Schwelle Unsinn.
Statt sie zu verschweigen: bekenne den Reifegrad. Ein bootstrap-aware
Gate dokumentiert seine Stufe und seinen Hochschalt-Trigger im
Make-Target:

```
coverage-gate: ## Coverage threshold gate (bootstrap-aware, LH-FA-BUILD-008).
```

Kam das Gate aus dem Steering Loop statt aus einer Anforderung, trägt der
Kommentar zusätzlich den Herkunfts-Anker `· seit welle-<NN>` — ohne Welle `· seit slice-<NNN>`
([`grundlagen-traceability.md` §Herkunfts-Anker](grundlagen-traceability.md#herkunfts-anker)).

Das Gate prüft heute z. B. 40 %, schaltet bei Meilenstein M2 auf 70 %
hoch. Das macht "bootstrap-aware" nicht zum Schlupfloch, sondern zum
**explizit terminierten Reifestufen-Gate** — ein Werkzeug eigener
Klasse, kein Subtyp von Carveout (die Werkzeug-Triade-Einordnung
steht direkt unter diesem Absatz).

**Werkzeug-Triade-Einordnung.** Bootstrap-aware Gate ist eine der
drei legitimen Antworten auf gelockerte Gate-Disziplin neben
*Carveout* (punktuelle Ausnahme mit Folge-Slice) und
*BF-Sub-Area-Markierung* (Sub-Area-weiter Übergangs-Modus mit
Graduation-Plan, Konzept in
[Modul 2 §Kernidee](modul-02-harness-bootstrap.md#kernidee-modul-2)).
**Die BF-Sub-Area-Markierung ist nicht selbst ein Closure-Werkzeug**,
sondern der Sub-Area-Kontext, in dem Carveout und Bootstrap-aware
Gate als Closure-Antworten strukturell legitim werden —
Disambiguierung in
[Modul 7 §Werkzeug-Wahl bei Diskrepanz](modul-07-carveouts.md#werkzeug-wahl).

**Begriffsklärung:** *Bootstrap-aware Gate* (oben) ist nicht zu
verwechseln mit *Harness-Bootstrap* aus
[`grundlagen-bootstrap.md` §Harness-Bootstrap](grundlagen-bootstrap.md#harness-bootstrap).
Letzteres ist der **Repo-Einstiegsprozess** (Lebenszyklus eines Harness
im Repo); ersteres ist die **Reifestufe eines einzelnen Sensors**.
Beide Begriffe teilen das Wort, sind strukturell verschieden.

### Reichhaltige Gate-Landschaft als Inspiration

Ein reifes Repo hat deutlich mehr als sechs Gates. Eine real gewachsene
Landschaft (Python/Docker, Simulations-Domäne) sieht etwa so aus:

```
lint · format-check · typecheck
arch-check · arch-check-imports · arch-check-custom
docs-check · spdx-check · noqa-check · noqa-gate
test-unit · test-determinism · test-replay · test-fault
test-integration
coverage-gate · coverage-gate-critical
dep-audit · image-audit · openapi-validate
```

Pointe: Domänenspezifische Gates (`test-determinism`, `test-replay`,
`noqa-gate`) entstehen aus dem Steering Loop — nicht aus einem
Standard-Setup. Wenn dein Repo nur die generischen sechs hat, weißt du
nur, dass du noch keine Schmerzen hattest.

`test-replay` ist dabei das **Gate**, nicht die Praxis: Das Replay-Set
selbst baut [Modul 12](modul-12-replay-evaluierung.md), dieses Target
setzt es durch. Gleiches Wort, zwei Ebenen.

In einer anderen Sprach-Welt wächst eine andere Landschaft: Ein Repo aus
C#/.NET mit Safety-/Control-Anteil bringt Gate-Familien mit, die im
Python-Beispiel oben gar nicht vorkommen — `solid-suppression-gate`
(C#-Pendant zum noqa-gate),
`test-mpc-property` (Property-Based-Sensor für Regelungstechnik),
`native-sanitizer` (für C/C++-Interop-Anteile), `test-hil-*`
(Hardware-in-the-Loop).

Pro Sprache wachsen also unterschiedliche Gate-Familien.

### Regeln gegen typische Fehlannahmen (Modul 13)

- Lint ist *ein* Gate-Typ. Architekturtests, Coverage-Gates, Security-Gates, Replay-Determinism-Gates sind weitere. Pro Repo entstehen sprachen- und domänenabhängige Gate-Familien.
- Dann ist es kein Gate, sondern ein Vorschlag. Pragmatik gehört in Carveouts oder bootstrap-aware Gates — mit Trigger und Folge-Slice.
- Es gibt keine universelle Schwelle. Critical Coverage (Security, Geld, Datenintegrität) ≠ Gesamt-Coverage. Schwellen sind ADR-pflichtig.
- Nur wenn lokal und CI dasselbe Image benutzen (Modul 14). Sonst debuggst du den Unterschied.
- Falsch in zwei Richtungen. Erstens: 80 % Gesamt-Coverage über *unkritischem* Code verbirgt 0 % Coverage auf dem Sicherheitspfad — Critical Coverage misst *gezielt*. Zweitens: Tests gegen Beispiele decken nur Realität ab, *wo das Golden Set repräsentativ ist* ([Modul 12](modul-12-replay-evaluierung.md)); Tests gegen die *Spec* erschließt Verifikation ([Modul 11](modul-11-verification.md)). Wer Test-Anzahl als Qualitätsmaß nimmt, baut Coverage-Anstiege, deren Wert auf 0 fällt, sobald die Realität die Coverage-Annahme bricht. Faustregel: *Verteilung vor Anzahl*. Ein zusätzlicher Test gegen einen bereits gut abgedeckten Pfad ist Boilerplate; ein zusätzlicher Test gegen einen *bisher unabgedeckten kritischen* Pfad ist Sensor.

<a id="adr-zur-fitness-function"></a>

### Fitness Function aus einem ADR-Satz (Modul 13)

Eine ADR *mit* Fitness Function ist ein Constraint statt einer
Absichtserklärung. Die Übersetzung in sechs Schritten (die sprachkonkrete
Implementierung zwischen Werkzeugwahl und Verdrahtung ist Illustration und
entfällt hier): **Aussage maschinell formulieren** → **Werkzeug wählen** →
**als Gate verdrahten** → `make gates` lokal grün — und **im CI mit gepinnter
Toolchain** (Modul 14) → **Bewusstes Brechen**: das Gate läuft rot mit
`ADR-<NNNN> violated`. Genau der Effekt, der eine ADR von einer
Absichtserklärung trennt.

| ADR-Satz (Beispiel) | Werkzeug | Make-Target | Failure-Beispiel |
|---|---|---|---|
| „Service importiert nur aus `adapter/`" | `import-linter`/`grimp` (Py) · `ArchUnit` (Java) · `depguard` (Go) · `dep-cruiser` (Node) | `arch-check:` ## LH-QA-COUPLING-002 / ADR-0007 | `import requests` in `service/foo.py` → `make arch-check` rot mit `ADR-0007 violated` |

Die **maschinelle Formulierung** ist die eigentliche Arbeit: aus
„importiert ausschließlich aus `adapter/`" wird „keine Datei unter
`src/service/**` enthält einen Import, dessen Modul nicht mit `adapter.`
beginnt oder Standardbibliothek ist" — erst diese Präzision ist als
`forbidden`-Contract eines Import-Linters prüfbar.

<a id="guard-haertung"></a>

### Guard-Härtung: Wächter reifen in Wellen (Modul 13)

Ein Wächter der [Durchsetzungsschicht](grundlagen-durchsetzungsschicht.md)
(Tool-Call-Gate, Handoff-Gate) ist Code *im* Harness und unterliegt
demselben Steering-Loop. Regeln für seine Härtung:

- **Auslöser ist Beobachtung, nicht Bedrohungsmodell.** Gehärtet wird
  gegen eine **dreimal beobachtete** Umgehung (Vorfall → Symptom →
  Lücke), belegt über Lerneinträge. Eine Wächter-Regel ohne
  Sensor-Evidenz ist Aufwand ohne Begründung und fällt beim ersten
  Fehlalarm.
- **Gehärtet wird die Zerlegung, nicht die Denylist.** Umgeht ein Aufruf
  den Guard über eine Sub-Shell (`bash -c "…"`, kombinierte Flags `-lc`,
  `-ec`), wird der Payload **rekursiv** derselben Prüfung unterworfen —
  mit Tiefenlimit, darüber fail-closed blockiert. Die Denylist um den
  Interpreter zu erweitern ist die falsche Reaktion: sie blockiert
  legitime Shell-Arbeit inklusive `make`. Ein **abgeschalteter Wächter
  ist schlechter als ein löchriger**, weil die Doku ihn weiter behauptet.
- **Jede Härtung landet als neuer `MR-<NNN>`**, der den vorherigen
  *schärft* — nie als inhaltliche Änderung eines akzeptierten Eintrags
  (Adaptions-Block-Disziplin,
  [`grundlagen-harness-dateien.md`](grundlagen-harness-dateien.md#harnessconventionsmd-als-konventionsspeicher)).
  Ein überschriebener Eintrag löscht, *welche* Umgehung die Härtung
  ausgelöst hat; die Regel wirkt später wie Overengineering.
- **Die Grenz-Zeile wird mitgezogen.** Jeder Wächter-`MR` trägt, was der
  Wächter *nicht* kann (`python -c "…"`, `env`-Umwege, Wrapper-Skripte —
  Netz dafür ist CI). Bleibt sie nach einer Härtung stehen, verspricht
  die Doku zu wenig oder zu viel; letzteres ist eine Harness-Lüge.
- **Wächter gehören nicht in die Gate-Typ-↔-Fehlerbild-Tabelle.** Ein
  Gate prüft ein *Ergebnis* (computational feedback), ein Wächter
  verhindert eine *Handlung* (computational feedforward). Er fängt kein
  Fehlerbild, er nimmt einen Weg weg.

Ziel-Form des Eintrags: `MR-<NNN>`-Schema in
[`../templates/harness/conventions.template.md`](../templates/harness/conventions.template.md)
§Adaptions-Block. Wo die Grenz-Aussage steht — als eigene `Grenze:`-Zeile
oder innerhalb von *Adaption* — ist Wahl; *dass* sie dasteht und bei jeder
Härtung mitwandert, ist Pflicht.
