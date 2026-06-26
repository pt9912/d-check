## Modul 13 — Quality Gates

*Quelle: [04-qualitaet/modul-13-quality-gates.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-13-quality-gates.md)*

### Harness-Einordnung (Modul 13)

Gates = *computational feedback* (siehe
[`grundlagen/klassifikation.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/klassifikation.md)).
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
| Coverage / Critical Coverage | Coverage-Loch — gesamt bzw. auf dem kritischen Pfad | Qualität der Tests, Spec-Lücken ([Modul 11](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-11-verification.md)) |
| Replay-/Determinism-Gate | nicht-deterministischer Test oder Lauf | semantische Drift außerhalb des Golden Sets ([Modul 12](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-12-replay-evaluierung.md)) |
| Integrationstest | Verhalten im Zusammenspiel: Komponenten-Vertrag bricht erst in Kombination | lokale Muster und Typgrenzen — dafür zu teuer und zu spät |

Trennlinie ist die *Regel-Klasse*, nicht das Tool: Linter machen lokale
Mustererkennung, Security-Regeln verlangen Datenfluss-Analyse,
Architekturtests prüfen Struktur, Integrationstests Verhalten im
Zusammenspiel.

### Hard Rule (Doku-Disziplin)

In `harness/README.md` und in jeder Doku, die Gates aufzählt: keine
Befehle behaupten, die es nicht gibt. Wenn `make fullbuild` strukturell
rot ist, wird das als Carveout in `docs/plan/carveouts/CO-<NNN>-…`
dokumentiert ([Modul 7](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/02-planung/modul-07-carveouts.md)) und in
der Bindung-Spalte der Sensors-Tabelle per `CO-<NNN>`-ID verlinkt — nicht
ausgelassen, nicht geschönt, nicht in einer Status-Spalte versteckt
(die Sensors-Tabelle trägt keinen Lauf-Status; Lauf-Wahrheit pro Commit
liegt in CI, siehe
[`grundlagen/konventionen.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/konventionen.md#harnessreadmemd-als-einstiegspunkt)).
Halluzinierte Gates sind die häufigste Form von Harness-Lüge — und der
Implementation-Agent vertraut ihnen.

### Bootstrap-aware Gates

In der Frühphase eines Projekts ist eine harte Coverage-Schwelle Unsinn.
Statt sie zu verschweigen: bekenne den Reifegrad. Ein bootstrap-aware
Gate dokumentiert seine Stufe und seinen Hochschalt-Trigger im
Make-Target:

```
coverage-gate: ## Coverage threshold gate (bootstrap-aware, LH-FA-BUILD-008).
```

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
[Modul 2 §Kernidee](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/01-spec-und-architektur/modul-02-harness-bootstrap.md#kernidee)).
**Die BF-Sub-Area-Markierung ist nicht selbst ein Closure-Werkzeug**,
sondern der Sub-Area-Kontext, in dem Carveout und Bootstrap-aware
Gate als Closure-Antworten strukturell legitim werden —
Disambiguierung in
[Modul 7 §Worked Example A Schritt 6](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/02-planung/modul-07-carveouts.md#worked-example-a-einen-carveout-dokumentieren).

**Begriffsklärung:** *Bootstrap-aware Gate* (oben) ist nicht zu
verwechseln mit *Harness-Bootstrap* aus
[`grundlagen/konventionen.md` §Harness-Bootstrap](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/konventionen.md#harness-bootstrap).
Letzteres ist der **Repo-Einstiegsprozess** (Lebenszyklus eines Harness
im Repo); ersteres ist die **Reifestufe eines einzelnen Sensors**.
Beide Begriffe teilen das Wort, sind strukturell verschieden.

### Reichhaltige Gate-Landschaft als Inspiration

Ein reifes Repo (Beispiel `pt9912/grid-gym`, siehe
[`grundlagen/fallstudien.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/fallstudien.md)) hat
deutlich mehr als sechs Gates:

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

Ein zweites Beispiel in einer anderen Sprach-Welt: `pt9912/bess-ems`
(C#/.NET, Safety/Control) bringt Gate-Familien mit, die `grid-gym`
nicht hat — `solid-suppression-gate` (C#-Pendant zum noqa-gate),
`test-mpc-property` (Property-Based-Sensor für Regelungstechnik),
`native-sanitizer` (für C/C++-Interop-Anteile), `test-hil-*`
(Hardware-in-the-Loop). Voll ausgeschrieben in
[`grundlagen/fallstudien.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/fallstudien.md).

Pro Sprache wachsen also unterschiedliche Gate-Familien.

### Regeln gegen typische Fehlannahmen (Modul 13)

- Lint ist *ein* Gate-Typ. Architekturtests, Coverage-Gates, Security-Gates, Replay-Determinism-Gates sind weitere. Pro Repo entstehen sprachen- und domänenabhängige Gate-Familien.
- Dann ist es kein Gate, sondern ein Vorschlag. Pragmatik gehört in Carveouts oder bootstrap-aware Gates — mit Trigger und Folge-Slice.
- Es gibt keine universelle Schwelle. Critical Coverage (Security, Geld, Datenintegrität) ≠ Gesamt-Coverage. Schwellen sind ADR-pflichtig.
- Nur wenn lokal und CI dasselbe Image benutzen (Modul 14). Sonst debuggst du den Unterschied.
- Falsch in zwei Richtungen. Erstens: 80 % Gesamt-Coverage über *unkritischem* Code verbirgt 0 % Coverage auf dem Sicherheitspfad — Critical Coverage misst *gezielt*. Zweitens: Tests gegen Beispiele decken nur Realität ab, *wo das Golden Set repräsentativ ist* ([Modul 12](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-12-replay-evaluierung.md)); Tests gegen die *Spec* erschließt Verifikation ([Modul 11](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-11-verification.md)). Wer Test-Anzahl als Qualitätsmaß nimmt, baut Coverage-Anstiege, deren Wert auf 0 fällt, sobald die Realität die Coverage-Annahme bricht. Faustregel: *Verteilung vor Anzahl*. Ein zusätzlicher Test gegen einen bereits gut abgedeckten Pfad ist Boilerplate; ein zusätzlicher Test gegen einen *bisher unabgedeckten kritischen* Pfad ist Sensor.

### Worked Example: vom ADR-Satz zur Fitness Function

**Ausgangs-ADR:** ADR-0007 (siehe Worked Example in [Modul 4](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/01-spec-und-architektur/modul-04-architektur-adrs.md#worked-example-vom-diskussionsfaden-zum-prüfbaren-adr)) sagt:

> "Service-Layer importiert ausschließlich aus `adapter/`-Paket."

**Schritt 1 — Aussage maschinell formulieren.** Aus *"importiert
ausschließlich aus"* wird:
> Keine Datei unter `src/service/**` darf einen Import enthalten, dessen
> Modul nicht mit `adapter.` beginnt oder ein Standardbibliotheks-Modul ist.

**Schritt 2 — Werkzeug wählen.** Python → `import-linter` oder
`grimp`. Java → `ArchUnit`. Go → `depguard`. Allgemein:
`dep-cruiser` für Node, eigene AST-Scanner für Nischensprachen.

**Schritt 3 — Implementierung (Python-Beispiel mit `import-linter`):**

```ini
# .importlinter
[importlinter]
root_packages = service

[importlinter:contract:service-adapter-only]
name = service imports only from adapter or stdlib
type = forbidden
source_modules =
    service
forbidden_modules =
    requests
    urllib3
    httpx
```

**Schritt 4 — Als Gate verdrahten:**

```makefile
arch-check:  ## LH-QA-COUPLING-002 / ADR-0007 — Service-Adapter-Trennung
	lint-imports
```

**Schritt 5 — `make gates` lokal grün — und im CI mit gepinnter
Toolchain (Modul 14).**

**Schritt 6 — Bewusstes Brechen:** Implementer fügt zu Debug-Zwecken
`import requests` in `service/foo.py`. `make arch-check` läuft rot mit
`ADR-0007 violated`. Genau der Effekt, der eine ADR von einer
Absichtserklärung trennt.
