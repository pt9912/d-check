# Impl-Review R1 — slice-066 (RTM-Quellen konfigurierbar, `trace`-Block)

**Datum:** 2026-07-11 · **Reviewer:** unabhängiger Subagent · **Verdikt:
ACCEPT-WITH-NITS** (kein Korrektheits-Defekt; 1 MEDIUM Test-Deckungslücke,
2 LOW/INFO).

**Gegenstand:** slice-066 / [ADR-0034](../plan/adr/0034-trace-konfigurierbare-quellen.md)
— die Requirements Traceability Matrix (`--trace`, `DC-FA-CLI-009`) wird über
einen opt-in `trace`-Block quell-/kennungs-konfigurierbar (vier Achsen:
Anforderungs-Quelldatei + Kennungs-Regex; je Referenzklasse ADR/Slice
Verzeichnis + Basisnamen-Regex mit Capture-Gruppe + Owner-Präfix). Default =
Konvention ⇒ byte-identisch. Lastenheft-CR (v0.40.0), Release geplant, **kein**
neues Modul, **kein** neuer Grund-Code.

**Geprüfte Artefakte:** `internal/hexagon/core/model/config.go`,
`internal/adapter/driven/configyaml/configyaml.go`,
`internal/hexagon/core/app/trace.go`,
`internal/adapter/driving/cli/cli.go`,
`internal/adapter/driving/cli/config_template.go`,
`internal/adapter/driving/cli/cli_acceptance_test.go`,
`spec/lastenheft.md`, `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`,
`docs/user/operations.md`, `CHANGELOG.md`, `docs/plan/adr/README.md`,
`docs/plan/planning/in-progress/roadmap.md`.

---

## Findings

### MEDIUM — Maintainability / DC-FA-CLI-009: 6 der 8 `trace`-Override-Achsen nicht differenziell getestet

**Pfad:** `internal/adapter/driving/cli/cli_acceptance_test.go:334` (`TestCLI066_*`)
· `internal/hexagon/core/app/trace.go:72` (`resolveTrace`).

**Befund:** Von den acht Config-Achsen exerzieren die Tests nur zwei so, dass ein
kaputter/entfernter Override-Zweig einen Test rot macht: `requirements.id-pattern`
und `slices.file-pattern` (beide über `TestCLI066_Trace_FremdKonvention`). Die
übrigen sechs Achsen — `requirements.source`, `adrs.dir`, `adrs.file-pattern`,
`adrs.id-prefix`, `slices.dir`, `slices.id-prefix` — hat kein Test mit einem
**Nicht-Default-Wert**, dessen Wirkung geprüft würde. `TestCLI066_Trace_DefaultByteIdentisch`
setzt alle acht Felder auf **ihre Defaults**; die Override-Zweige
(`if tc.Source != "" { … }` usw.) laufen zwar durch, ändern das Ergebnis aber
nicht — die Mutation ist inert und überlebt. Der Slice fordert im DoD/§5
ausdrücklich „mutations-hart".

**Failure-Szenario (Beispiel):** Ein Refactor lässt versehentlich
`if tc.ADRPrefix != "" { rt.adrPrefix = tc.ADRPrefix }`
([trace.go:90](../../internal/hexagon/core/app/trace.go)) weg (oder eine
Copy-Paste-Vertauschung `adrDir`↔`sliceDir`). Ein Konsument mit
`trace.adrs.id-prefix: 'RFC-'` bekäme still weiter `ADR-`-Owner-Kennungen; kein
Test wird rot. Ebenso: `trace.requirements.source` oder `trace.adrs.dir` auf ein
Nicht-Default-Verzeichnis zeigend würde bei defektem Override still die
Default-Quelle lesen — genau der „andere Quell-Pfad"-Fall, den die
Out-of-Scope-Umkehr adressiert.

**verifizierbar:** ja — ein Test mit non-default `source`/`adrs.dir`/`adrs.id-prefix`/
`slices.id-prefix`, der die veränderte Zeile/den veränderten Owner assertet,
würde die jeweilige Mutation töten; heute überleben sie `go test ./...`.

### LOW — Maintainability: keine Zusicherung, dass `--print-config` den `trace`-Block emittiert

**Pfad:** `internal/adapter/driving/cli/config_template.go:144` ·
`internal/adapter/driving/cli/cli_acceptance_test.go:834`
(`TestCLI053_PrintConfig_VollesModulset`).

**Befund:** Der neue `trace`-Block im `--print-config`-Gerüst ist vollständig
auskommentiert und damit für `configyaml.Decode` inert — `TestCLI005`/`TestCLI053`
dekodieren nur (⇒ trivially grün) bzw. prüfen eine Präsenz-Stringliste, in der
`trace` **nicht** steht. Entfernte man den Block aus dem Template (oder tippt sich
in seinem Regex-Kommentar), fällt kein Test.

**Failure-Szenario:** Ein Merge-Konflikt tilgt den `# --- trace:`-Block aus
`config_template.go`; die „Verfügbar-Vollständigkeit" (ADR-0034 §Konsequenzen)
regrediert unbemerkt. **verifizierbar:** ja — eine `want`-Zeile `"# --- trace:"`
(bzw. `"trace:"`) in `TestCLI053` fixiert es.

### INFO — DC-FA-CLI-009 / ADR-0034: `--trace` verlangt nun eine gültige `.d-check.yml`

**Pfad:** `internal/adapter/driving/cli/cli.go:431` (`loadConfig` vor dem
`--trace`-Zweig).

**Befund:** Vor slice-066 schloss `--trace` **vor** `loadConfig` kurz und ignorierte
`.d-check.yml` vollständig; jetzt lädt/validiert `Run` die Config zuerst. Ein Repo
mit einer defekten, `trace`-**fremden** `.d-check.yml` (z. B. ungültiges
`versions.pin-pattern`) liefert bei `--trace` nun Exit 2 statt einer RTM. Das ist
die **beabsichtigte** fail-closed-Semantik (ADR-0034 §Konsequenzen; DoD „Config-Fehler
⇒ Exit 2") und zugleich ein Verhaltens-Delta ohne eigenen Test. Nicht blockierend,
hier nur dokumentiert.

---

## Verifiziert sauber (Negativbefunde)

- **Byte-Identität im Default (Kernrisiko):** `resolveTrace` mit Nullwert-
  `TraceConfig` ([trace.go:72](../../internal/hexagon/core/app/trace.go)) ergibt
  exakt die alten Konstanten/Regexe — `source=spec/lastenheft.md`, `reqPat=reqIDFull`,
  `adrDir/adrFile/"ADR-"`, `sliceDir/sliceFile/""` — deckungsgleich mit den
  vorherigen hart kodierten `traceRefs(…, "ADR-")` / `traceRefs(…, "")`-Aufrufen.
  Zeile-für-Zeile geprüft; `traceRequirements`/`traceRefs` nutzen den durchgereichten
  `reqPat` statt der Globalen, im Nullwert identisch. `TestCLI066_Trace_DefaultByteIdentisch`
  vergleicht „ohne Block" vs. „Default-Block 1:1" byte-genau. Kein Verhaltens-Change
  im Default. Sauber.
- **Fail-closed-Vollständigkeit (Capture-Guard):** `compileTracePattern`
  ([configyaml.go:311](../../internal/adapter/driven/configyaml/configyaml.go))
  fordert `requireCapture=true` für **beide** `file-pattern` (adrs **und** slices),
  `false` für `id-pattern` (das keine Gruppe braucht). Der `m[1]`-Zugriff in
  `traceRefs` ([trace.go:209](../../internal/hexagon/core/app/trace.go)) läuft nur
  über `fileShape`, das im Default eine Gruppe trägt und im Custom-Fall garantiert
  ≥1 hat ⇒ **kein** Panic-Pfad. Ein nicht-teilnehmender Alternations-Zweig gäbe
  `m[1]==""` (leerer Owner), keinen Absturz. Guard ist load-bearing:
  `TestCLI066_Trace_FilePatternOhneCapture` würde ohne Guard das reale `m[1]`-Panic
  (traceRepo-Slice `slice-099-x.md` matcht das gruppenlose Muster) statt Exit 2
  auslösen ⇒ Mutation gefangen. Sauber.
- **`--require-complete` erbt die Config:** `runTrace` reicht `cfg.Trace` durch
  denselben `BuildTraceMatrix`-Lauf ([cli.go:181](../../internal/adapter/driving/cli/cli.go));
  `TestCLI066_Trace_RequireCompleteErbtConfig` ist ein echter Vorher/Nachher-
  Differenzialtest (ohne Config Exit 0, mit `id-pattern` Exit 1). Sauber.
- **Prefix-Default-Semantik konsistent:** ADR-Präfix leer ⇒ Default `ADR-` (nicht
  ausdrückbar), Slice-Präfix Default leer und direkt übernommen — identisch in Code
  (`resolveTrace`), Spezifikation §2-Tabelle und Lastenheft-Out-of-Scope. Sauber.
- **Norm↔Code-Parität:** die acht `rawTrace`-YAML-Keys entsprechen 1:1 den acht
  §2-Schema-Zeilen; alle Default-Werte der Tabelle stimmen mit den Code-Konstanten
  (`traceReqSource`, `reqIDFull`, `traceADRDir`, `adrFileShape`, `"ADR-"`,
  `traceSliceDir`, `sliceFileShape`, `""`). §DC-FA-CLI-009.a Schritte 1/2 +
  Config-Auflösungs-Absatz decken den Algorithmus. Sauber.
- **Fail-closed zur Config-Zeit (Exit 2, vor Scan):** ungültige Regex
  (`TestCLI066_Trace_UngueltigeRegex`), `file-pattern` ohne Gruppe
  (`…_FilePatternOhneCapture`), Pfad-Escape `..` (`…_PfadEscape`,
  `validateTracePath` analog `planning.roadmap`) — alle drei Negativ-Tests grün-
  fordernd Exit 2 mit Feld-Nennung im stderr. Sauber.
- **Usage-Fehler unberührt vom `loadConfig`-Reorder:** `--trace`+`--doctor`/`--repair`
  und `--require-complete` ohne `--trace` werden in `comboError`
  ([cli.go:141](../../internal/adapter/driving/cli/cli.go)) zur **Parse-Zeit** (vor
  `loadConfig`) auf Exit 2 abgebildet; `--commit-msg`+`--trace` in `gitComboError`.
  Der Reorder verschiebt nur den Scan-Pfad, nicht die Modus-Validierung. Sauber.
- **Fremd-Konvention echt differenziell:** `TestCLI066_Trace_FremdKonvention`
  prüft Vorher (Default sieht `GG-ARCH-001` nicht, `GG-QA-001` Waise) → Nachher
  (Total 3 / Orphans 1, `GG-ARCH-001` mit `ADR-0007`+Slice `063`). Belegt die
  Fitness-Funktion des ADR (grid-gym-Gestalt vollständig statt 6/243). Sauber.
- **Harness-Regeln:** ADR-0034 nennt `slice-066` **nur** in `## Geschichte`
  (Zeile 153; Body frei von Abwärts-Token — matrix-konform). Die zwei neuen
  Handbuch-```yaml```-Blöcke (§4.12, §5) sind gültige `trace`-Configs und laufen
  über die slice-061-`docexamples`-Decode-Harness (kein `not-config`-Marker nötig);
  der ADR-``` ```-Block ist ein nackter Fence (kein `yaml`-Tag) und wird nicht
  dekodiert. Lastenheft-Versions-Header `0.40.0` == oberste §7-Historien-Zeile
  `0.40.0` == CHANGELOG-Top. ADR-README-Indexzeile ergänzt. Sauber.
- **Determinismus/Read-only (DC-QA-02/03):** `BuildTraceMatrix` nutzt nur
  `ReadFile`/`DiscoverFiles`/`pathExists` + `sort.Strings`; kein Schreibpfad, kein
  git/Netz. `TestCLI036_Trace_Markdown` bekräftigt „nichts geschrieben". Sauber.
- **Lint/Komplexität:** `applyTrace` (~11) und `resolveTrace` (~8) liegen unter der
  gocyclo-Schwelle 15 (`.golangci.yml`); die pre-existierenden Regex-Globalen
  (`reqIDFull`/`adrFileShape`/`sliceFileShape`) sind nur read, unverändert.
  Keine tote/duplizierte Logik. Sauber.

---

## Kategorie-Summary

| Kategorie | Zahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 (Test-Deckung: 6/8 Override-Achsen nicht mutationsgetestet) |
| LOW | 1 (keine `--print-config`-`trace`-Präsenz-Assertion) |
| INFO | 1 (`--trace` verlangt nun gültige Config — beabsichtigt) |

## Verdikt

**ACCEPT-WITH-NITS.** Kern korrekt und byte-identisch im Default; Fail-closed,
Capture-Guard, `--require-complete`-Vererbung und Norm↔Code-Parität sauber
verifiziert. Empfehlung vor Closure (gemäß DoD „mutations-hart"): je einen
differenziellen Testfall für `trace.requirements.source`, ein Referenz-`dir` und
einen Nicht-Default-`id-prefix` ergänzen (MEDIUM), sowie eine `trace`-Präsenz-
Zeile in `TestCLI053` (LOW). Kein Produkt-Code-Blocker.
