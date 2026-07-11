# ADR-0034 — RTM-Quellen und Kennungs-Muster konfigurierbar (opt-in `trace`-Block)

**Status:** Accepted
**Datum:** 2026-07-11
**Autor:** pt9912
**Bezug:** [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(die geänderte Anforderung; Spezifikation
[§`DC-FA-CLI-009.a`](../../../spec/spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix));
[`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(erbt die konfigurierten Quellen — gleicher RTM-Lauf);
[ADR-0026](0026-completeness-in-product-gate.md) (das `completeness-check`-Gate
dogfoodet denselben `--trace`-Pfad);
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).

## Kontext

Der `--trace`-Modus
([`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
leitet die Requirements Traceability Matrix aus drei Quell-Achsen ab, die **hart
an d-checks eigene Konvention gebunden** sind — keine davon ist heute
konfigurierbar:

1. **Kennungs-Gestalt** der Anforderungen: der Regex
   `[A-Z][A-Z0-9]*-(?:FA-[A-Z]+|QA)-\d+[A-Za-z]?` matcht nur die
   `-FA-<BEREICH>-`/`-QA-`-Mittelsegmente von d-checks Schema.
2. **Slice-Dateiname:** `^(slice-\d+)-.*\.md$` — nur `slice-NNN-…md`.
3. **ADR-Dateiname:** `^(\d{4})-.*\.md$` → `ADR-NNNN` (der einzige, den
   Fremd-Repos häufig teilen).

Die Quell-Pfade (`spec/lastenheft.md`, `docs/plan/adr`, `docs/plan/planning`)
sind ebenfalls Konstanten. Die Erst-Fassung
([`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix),
Lastenheft 0.21.0) schloss **frei konfigurierbare Quell-Pfade jenseits der
adoptierten Harness-Konvention** ausdrücklich als Out-of-Scope aus — eine bewusst
enge Erst-Fassung.

Ein Konsumenten-Befund kehrt diese Einschätzung um: das Schwester-Repo
**grid-gym** bindet d-check via `d-check.mk` ein und ruft `make doc-trace`. Sein
Lastenheft definiert **243** Anforderungen in ~40 Familien
(`GG-ARCH-*`, `GG-PRINC-*`, `GG-CC-*`, `GG-TEST-*`, `GG-QA-*` …) nach dem Schema
`GG-<FAMILIE>-NNN`; seine Slices heißen `NNN-titel.md`. Die RTM sah davon **6**
Anforderungen — ausschließlich die `GG-QA-*`-Familie, die als einzige zufällig
den `-QA-`-Zweig des Default-Regex trifft — und markierte alle als Waisen, weil
kein `slice-NNN-…md`-Dateiname (und damit kein Slice-Owner) erkannt wird. Ein via
`--print-mk` an **jeden** Konsumenten verteiltes Feature
([`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
Zieltargets `doc-trace`/`doc-complete`) funktioniert also nur für Repos, die
d-checks internes Kennungs-/Datei-Schema nachahmen. Für alle anderen ist die
Ausgabe stilles Rauschen — und `doc-complete` (`--require-complete`) als CI-Gate
wäre still fehlerhaft grün (prüft nur die zufällig sichtbare Teilmenge).

## Entscheidung

### 1. Ein opt-in `trace`-Config-Block, der die vier Achsen überschreibt

`.d-check.yml` erhält einen optionalen `trace`-Block mit genau den Feldern, die
heute Konstanten sind:

```
trace:
  requirements: { source, id-pattern }
  adrs:         { dir, file-pattern, id-prefix }
  slices:       { dir, file-pattern, id-prefix }
```

- `requirements.source` — Quell-Datei der Anforderungs-Headings.
- `requirements.id-pattern` — Regex, der eine Anforderungs-Kennung erkennt:
  als **Ganz-Token** im Heading (der Treffer muss das ganze erste Token sein)
  **und** als Vorkommen in ADR-/Slice-Dateien (Referenz-Zählung). Ein Regex für
  beide Rollen — konsistent mit dem heutigen `reqIDFull`.
- je Referenzklasse (`adrs`, `slices`): `dir` (Scan-Verzeichnis),
  `file-pattern` (Regex auf den Basisnamen; **Capture-Gruppe 1** = Owner-Kennung)
  und `id-prefix` (der Capture-Gruppe vorangestellt).

### 2. Alle Felder optional; Default = heutige Konstante ⇒ byte-identisch

Jedes abwesende/leere Feld fällt auf den heutigen hart kodierten Wert zurück
(`Effective*()`-Methoden im Modell, wie bei den 15 Modul-Configs). Eine Config
ohne `trace`-Block liefert eine **byte-identische** RTM
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)); d-checks
eigenes `make doc-trace`/`make completeness-check` bleibt unverändert. Der
`trace`-Block ist reine Erweiterung, kein Verhaltens-Change im Default.

### 3. Design spiegelt `ids.patterns`, nicht dessen Ableitung

Der Block ist ein eigenständiges, **explizit gesetztes** Schema (Regex + Ziel,
wie `ids.patterns`). Er wird **nicht** aus vorhandenen `ids.patterns`-Einträgen
abgeleitet: `ids` trägt mehrere Muster mit `link-policy`/`exempt-paths`-Semantik;
welches davon „die Anforderungs-Kennung" ist, wäre eine implizite, brüchige
Heuristik. Ein dediziertes, selbst-beschreibendes `trace`-Schema ist ehrlicher
und stabiler.

### 4. Fail-closed zur Config-Zeit

Jede gesetzte `id-pattern`/`file-pattern` muss ein gültiges Regex sein, und jede
gesetzte `file-pattern` muss **mindestens eine Capture-Gruppe** tragen (sonst
wäre die Owner-Kennung undefiniert — heute ein `m[1]`-Panic-Risiko). Verstoß ⇒
Exit 2 mit erklärender Meldung, vor jedem Scan — dieselbe fail-closed-Politik wie
`versions.pin-pattern`/`matrix.classes[].token`.

### Verglichene Alternativen

| Option | Pro | Contra |
| --- | --- | --- |
| **(A) Status quo — Konvention hart kodiert** | null Code; d-checks Dogfood unberührt | von grid-gym widerlegt: `--print-mk` verteilt ein Feature, das nur d-checks Schema-Familie bedient; `doc-complete` still falsch-grün |
| **(B) `trace`-Muster aus `ids.patterns` ableiten** | keine doppelte Config | implizit/brüchig (welches der N `ids`-Muster ist die Anforderung?); koppelt zwei unabhängige Vertragsflächen |
| **(C, gewählt) dedizierter opt-in `trace`-Block, Default = Konvention** | explizit, selbst-beschreibend, byte-identisch im Default, spiegelt die 15 Modul-Configs | eine neue Config-Fläche (Schema, Decode, Validierung, print-config) |

**Fitness-Funktion:**

- grid-gym mit `trace.requirements.id-pattern` + `trace.slices.file-pattern` sieht
  **alle 243** Anforderungen mit ihren Slices statt 6 (Beleg im umsetzenden Slice).
- d-checks eigener `make doc-trace`/`make doc-check` bleibt **byte-identisch**
  (kein `trace`-Block; Dogfood-Regressionstest).
- `trace`-Config mit ungültiger Regex bzw. `file-pattern` ohne Capture-Gruppe ⇒
  Exit 2 (Negativ-Test).

## Konsequenzen

- **Geänderte Anforderung (Lastenheft-CR, Versions-Bump + Historie).** Der
  `trace`-Block ist eine nutzersichtbare Config-Fläche und **kehrt die
  0.21.0-Out-of-Scope-Zeile um** → geänderte Anforderung im Lastenheft
  ([`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)),
  **kein** neues Modul, **keine** neue `DC-*`-ID, **kein** neuer Grund-Code
  (`--trace` bleibt advisory).
- **Modell** (`internal/hexagon/core/model/config.go`): neues `TraceConfig`
  (kompilierte Regex + Strings) mit `Effective*()`-Defaults, analog zu
  `PlanningConfig`/`VersionsConfig`. `Config` trägt ein `Trace`-Feld.
- **Config-Decode** (`internal/adapter/driven/configyaml/configyaml.go`): neuer
  `trace`-Top-Level-Key; Regex-Kompilierung + Capture-Gruppen-Guard fail-closed.
- **RTM** (`internal/hexagon/core/app/trace.go`): `BuildTraceMatrix` nimmt die
  aufgelöste `TraceConfig` statt der Konstanten; Quell-Pfade, Kennungs-Regex und
  Datei-Gestalten kommen aus der Config.
- **CLI** (`internal/adapter/driving/cli/cli.go`): der `--trace`-Zweig lädt die
  Config (`loadConfig`) **vor** `runTrace` (heute danach) und reicht `cfg.Trace`
  durch. Ein Config-Fehler ⇒ Exit 2 (fail-closed), sonst unverändert.
- **`--print-config`** (`config_template.go`): ein kommentierter `trace`-Block
  wird ergänzt (Vollständigkeit der Verfügbar-Liste — dieselbe Harness-
  Ehrlichkeit wie die Modul-Nachträge früherer Slices).
- **Determinismus/Read-only unberührt**
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)/[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)):
  kein neuer Eingabe-Scope (weiterhin nur der read-only Markdown-Baum), keine
  Toolchain, kein Netz. VCS-/git-historienbasierte Erkennung bleibt Out-of-Scope.
- **Reversibel** im Verhalten (Default byte-identisch), aber Vertrags-Änderung —
  daher Lastenheft-CR statt stiller Code-Änderung. **Release** (Config-Fläche
  nutzersichtbar).

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-07-11 | Entwurf (slice-066, welle-55; Konsumenten-Befund grid-gym: `make doc-trace` sah 6 von 243 Anforderungen, weil nur die `GG-QA-*`-Familie zufällig d-checks `-QA-`-Default-Gestalt trifft und `NNN-…md`-Slices unerkannt bleiben). Opt-in `trace`-Block überschreibt die vier Konventions-Achsen (Anforderungs-Quelle + Kennungs-Regex; je Referenzklasse Verzeichnis + Basisnamen-Gestalt + Owner-Präfix); Default = heutige Konstante ⇒ byte-identisch; Design spiegelt `ids.patterns` (dediziert, nicht abgeleitet); fail-closed (Regex/Capture-Gruppe). Kehrt die Lastenheft-0.21.0-Out-of-Scope-Zeile um. Lastenheft-CR (v0.40.0), Release geplant. Status Proposed. |
| 2026-07-11 | Angenommen mit der slice-066-Closure: `model.TraceConfig` + `configyaml.applyTrace` (Regex-/Capture-/Pfad-Guards, Exit 2) + `app.resolveTrace`/`BuildTraceMatrix(fsys, tc)`; `--print-config` führt einen kommentierten `trace`-Block. End-to-End gegen grid-gyms Realdaten verifiziert (6 → 243 Anforderungen). Impl-Review R1 ACCEPT-WITH-NITS (Voll-Custom-Test für alle 8 Achsen + `TestCLI053`-Verankerung eingearbeitet). `make gates` grün, Release v0.40.0. Status Accepted. |
