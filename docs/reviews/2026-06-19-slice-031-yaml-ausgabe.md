# Review — slice-031 Implementierung (YAML-Ausgabeformat `--yaml`)

## Kopf-Metadaten

- **Review-Art:** Code-Review (fertiger Diff gegen Plan/Spec/Anforderungen/
  Hard Rules — **kein Verifier**: DoD-Abhaken und Gate-Lauf-Bestätigung sind
  nicht Gegenstand; Gates werden NICHT als grün angenommen).
- **Datum:** 2026-06-19
- **Gegenstand:** Working-Tree-Diff (unstaged) der slice-031-Umsetzung plus
  die untrackte Slice-Datei
  `docs/plan/planning/done/slice-031-yaml-ausgabe.md`. `--yaml` gibt die
  Befunde als YAML auf stdout aus — strukturgleich zu `--json`;
  `--doctor --yaml` analog `--doctor --json`. Neu: yaml.v3 auch im
  report-Adapter (formatneutrale Output-Structs `outDoc`/`outDiagDoc`/
  `outDiagFinding`/`outFixCandidate` mit json+yaml-Tags; `encodeJSON`/
  `encodeYAML`; `buildDiagDoc`; `JSON`/`YAML`/`DoctorJSON`/`DoctorYAML`),
  yaml-Tags an `core.Finding`, `--yaml`-Flag + Kombi-Checks + render→
  `render`/`renderStdout`-Refactor in der CLI, arch-check R3-Lockerung,
  vier neue Akzeptanztests, Spezifikation (§2 „YAML-Ausgabe"),
  Operations/Handbuch. **Bezug (bereits committet, nicht im Diff):**
  Lastenheft `DC-FA-CLI-004` 0.19.0 (beide neuen AK), ADR-0009 (Accepted
  via diesen Diff).
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/done/slice-031-yaml-ausgabe.md`; Anforderungen
  `DC-FA-CLI-004` (Haupt — YAML als Ausgabeformat, beide AK),
  `DC-FA-CLI-007` (`--doctor --yaml`), `DC-FA-CLI-003` (Exit-Codes),
  `DC-QA-02` (Determinismus), `DC-QA-03` (read-only/netzlos); Spezifikation
  §„YAML-Ausgabe"; ADR-0009 (yaml.v3 im report-Adapter, erweitert ADR-0005
  Regel 3), ADR-0005 (Hexagon-Import-Regeln 1/3/5); Hard Rules `AGENTS.md`
  §3. **Die DoD-Abhakung des Slices lag bewusst nicht vor** — die Findings
  sind eigenständig gebildet. Tests/Gates wurden NICHT ausgeführt
  (make/Docker-only; Reviewer ist kein Verifier).

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

#### LOW-1 — Vergleichstabelle „Die drei Ausgaben im Vergleich" ohne YAML-Pendants

- **Kategorie:** LOW
- **Quelle:** Maintainability (Doku-Drift zwischen zwei Stellen derselben
  Feature-Oberfläche)
- **Pfad:** `docs/user/benutzerhandbuch.md:390`
- **Befund:** Die §4.9-Tabelle listet die maschinenlesbare Oberfläche als
  `--json`/`--doctor`/`--doctor --json`, wurde aber nicht um die
  YAML-Pendants (`--yaml`/`--doctor --yaml`) ergänzt. operations.md
  (Options-Tabelle) und §4.11 nennen `--yaml` hingegen; ein Leser, der die
  Vergleichstabelle als vollständige Modus-Übersicht nimmt, sieht das
  YAML-Format dort nicht. Kein Korrektheitsdefekt — eine Konsistenz-Lücke
  zwischen zwei Doku-Stellen, die dieselbe Ausgabe-Oberfläche beschreiben.
- **Verifizierbar:** ja — ein Diff der Tabelle gegen die Options-Tabelle in
  `operations.md` zeigt die fehlenden YAML-Zeilen.

#### LOW-2 — §4.11-Überschrift trägt nur `(--json)`, deckt nun beide Formate

- **Kategorie:** LOW
- **Quelle:** Maintainability (Doku-Drift Heading↔Inhalt)
- **Pfad:** `docs/user/benutzerhandbuch.md:445`
- **Befund:** Der ergänzte §4.11-Abschnitt beschreibt nun JSON **und** YAML
  (inkl. `--doctor --yaml`), die Überschrift lautet aber weiterhin
  „4.11 Maschinenlesbare Ausgabe (`--json`)". Der Anker bleibt intakt (kein
  gebrochener §4.11-Verweis), die Suffix-Auszeichnung im Titel nennt das
  zweite Format jedoch nicht. Reine Heading-Drift, kein Korrektheits- oder
  Verweis-Defekt.
- **Verifizierbar:** ja — sichtbar im Diff (Heading unverändert, Body um
  YAML erweitert).

### INFO

#### INFO-1 — Text-Render-Fehler erzeugt jetzt eine stderr-Meldung (vorher still)

- **Kategorie:** INFO
- **Quelle:** Maintainability (beobachtbare Verhaltensänderung des
  render-Refactors)
- **Pfad:** `internal/adapter/driving/cli/cli.go:260`
- **Befund:** Der frühere `default`-Zweig behandelte einen Fehler von
  `report.Text` mit `if err != nil { return 2 }` — **ohne** stderr-Meldung.
  Nach dem Refactor liefert `renderStdout` `(exit, err)`, und `render`
  schreibt für jeden Render-Fehler `d-check: error: %v` auf stderr und gibt
  2 zurück. Damit erscheint im (praktisch nur bei stdout-Schreibfehler
  erreichbaren) Text-Fehlerpfad neu eine stderr-Zeile; Exit-Code (2) bleibt
  gleich. Verhalten vereinheitlicht sich mit den übrigen Modi; kein Test
  fordert Stille auf diesem Pfad. Dokumentationswürdige Verhaltensangleichung
  ohne Korrektheitswirkung.
- **Verifizierbar:** ja — ein Test, der `report.Text` einen
  stdout-Schreibfehler unterschiebt, würde die neue stderr-Zeile zeigen
  (vorher leer).

## Negativbefunde (geprüft, ohne Befund)

- **Struktur-Parität YAML↔JSON (knappe Befunde):** `outDoc` trägt für
  jedes Feld json- **und** yaml-Tag (`report.go` `outDoc`-Struct,
  `findings`/`summary`/`exitCode`); `Summary` ebenso
  (`filesChecked`/`findingCount`, `report.go:14`); `core.Finding` ist mit
  yaml-Tags pro Feld versehen
  (`internal/hexagon/core/finding.go:23` ff., `file`/`line`/`rule`/`target`/
  `reason`/`message,omitempty`). yaml.v3 v3.0.1 (`go.mod:5`) würde untaggte
  Feldnamen kleinschreiben (`filesChecked`→`fileschecked`); die expliziten
  Tags verhindern das. `TestCLI004_YAML`
  (`internal/adapter/driving/cli/cli_acceptance_test.go:232`) prüft
  `findingCount:` vorhanden und `fileschecked:` abwesend sowie gleiche
  `findingCount`/`exitCode`/Listenlänge wie JSON.
- **Embedded-Flattening der Diagnose:** `outDiagFinding` bettet
  `core.Finding` anonym ein und trägt `yaml:",inline"`
  (`report.go` `outDiagFinding`-Struct) — yaml.v3 flacht anonyme Structs nur
  mit `,inline` ab, sonst verschachtelt. Die JSON-Seite flacht über die
  anonyme Promotion (kein Tag). `TestCLI007_DoctorYAML`
  (`cli_acceptance_test.go:285`) liest `reason`/`reasonText` flach auf einer
  Ebene; fehlte das `reasonText` (Verschachtelung), schlüge der Test fehl
  (`cli_acceptance_test.go:292`).
- **`fixCandidate: null` explizit in YAML:** `outDiagFinding.FixCandidate`
  ist `*outFixCandidate` **ohne** omitempty (json+yaml)
  (`report.go` `outDiagFinding`-Struct); bleibt der Zeiger nil, serialisieren
  beide Formate explizit `null`. `TestCLI007_DoctorYAML`
  (`cli_acceptance_test.go:308`) prüft `strings.Contains(stdout,
  "fixCandidate: null")` und zählt genau einen null- und einen
  Kandidaten-Eintrag.
- **Determinismus (DC-QA-02):** Render-Pfad serialisiert ausschließlich
  Structs in fester Feld-Reihenfolge (`outDoc`/`outDiagDoc`), keine Map;
  `buildDiagDoc` (`report.go:201`) iteriert über das stabil sortierte
  `findings`-Slice und füllt `make([]outDiagFinding, 0, len)`.
  yaml.v3-Struct-Marshal ist feldreihenfolge-stabil.
  `TestCLI004_YAML_Determinismus` (`cli_acceptance_test.go:328`) belegt 10
  byte-identische Läufe.
- **encodeYAML-Korrektheit (Flush/Fehler):** `encodeYAML` (`report.go:170`)
  setzt 2-Space-Einrückung, encodet und ruft in beiden Pfaden `Close()`
  (Fehlerfall: `_ = enc.Close()` + return err; Erfolgsfall: `return
  enc.Close()`); yaml.v3 schreibt erst beim `Close()` das
  Dokument-Terminierende/Final-Flush, sodass kein Output verloren geht und
  der Close-Fehler propagiert. Der Aufrufer `render` mappt einen Rückgabe-
  Fehler auf Exit 2 (`cli.go:226`).
- **leere Befundliste = `[]` nicht `null`:** `nonNil` (`report.go:160`)
  ersetzt eine nil-Befundliste durch `[]core.Finding{}` für `JSON`/`YAML`;
  `buildDiagDoc` nutzt `make(…, 0, len)` (`report.go:202`). Beide
  Serialisierungen geben damit eine leere Liste statt `null` aus — Parität
  zur bisherigen JSON-Ausgabe.
- **Kombi-Vertrag (DC-FA-CLI-004/003):** `--json`+`--yaml` → Exit 2
  (`cli.go:161`), `--repair`+`--json/--yaml` → Exit 2 (`cli.go:165`),
  `--doctor`+`--repair` → Exit 2 (`cli.go:169`); `--doctor`+`--yaml`
  **erlaubt** (Switch-Zweig `cli.go:246`). `TestCLI004_YAML_Negative`
  (`cli_acceptance_test.go:315`) prüft beide Exit-2-Fälle samt leerem stdout
  und Meldungstext. Deckt sich mit den AK „YAML"/„YAML Negative" des
  committeten `DC-FA-CLI-004` (`spec/lastenheft.md`).
- **Bestands-Test `TestCLI008_Repair_Inkompatibel` nicht gebrochen:** der
  Test prüft `strings.Contains(stderr, "nicht mit --json")`
  (`cli_acceptance_test.go:689`); die geänderte Meldung lautet „--repair ist
  nicht mit --json/--yaml kombinierbar" (`cli.go:166`) und **enthält** den
  Teilstring „nicht mit --json". Der zweite Teil (`--doctor --repair`,
  „nicht kombinierbar", `cli_acceptance_test.go:693`) trifft weiterhin
  `cli.go:170`.
- **render-Refactor verhaltensgleich:** `renderStdout` (`cli.go:237`)
  setzt `exit` weiterhin 0/1 nach Befund-Stand und gibt für jeden Modus
  `(exit, modusErr)` zurück; `render` (`cli.go:224`) bildet jeden Fehler auf
  Exit 2 ab. Die Repair-stderr-Ausgabe bleibt in `report.Repair`
  (`cli.go:255`, Best-Guess-Marker auf stderr unverändert); Doctor-/Text-
  Zusammenfassung weiterhin auf stderr (`report.go:31`, `:66`). Exit-Codes
  0/1 (Befund) bzw. 2 (Render-/Repair-Fehler) unverändert (Ausnahme: die
  unter INFO-1 notierte zusätzliche stderr-Zeile im Text-Fehlerpfad).
- **arch-check R3 (ADR-0009) exakt:** die yaml.v3-Regel lässt per
  `case` genau `internal/adapter/driven/configyaml` **und**
  `internal/adapter/driven/report` zu, sonst `violation R3`
  (`tools/arch-check.sh:48`). Kein weiteres Paket; Regel-/Kommentartext
  nachgezogen (`tools/arch-check.sh:11`). R1 blockt yaml in
  `internal/hexagon/*` unverändert (`tools/arch-check.sh:39`) — der Kern
  bleibt yaml-frei. R5 (`tools/arch-check.sh:59`) prüft nur
  Modul-interne driven↔driven-Importe; der yaml.v3-Import von `report` ist
  eine externe Lib, kein anderer driven-Adapter, also keine R5-Verletzung;
  `report` importiert weiterhin nur `core` + stdlib + yaml.v3
  (`report.go:5`-Importblock). Deckt sich mit ADR-0009 Entscheidung/Fitness
  Function („genau zwei Adapter").
- **read-only / netzlos (DC-QA-03):** der YAML-Pfad schreibt nur in den
  übergebenen stdout-Writer (`encodeYAML`, `report.go:170`); kein Datei-,
  git- oder Netzzugriff. `YAML`/`DoctorYAML` nehmen dieselben
  Eingaben wie die JSON-Varianten und führen kein zusätzliches I/O ein.
- **ADR-0009/ADR-0005-Doku-Konsistenz:** ADR-0009 auf `Accepted`
  (`docs/plan/adr/0009-yaml-im-report-adapter.md:3`) samt
  Geschichts-Zeile (`:72`) und README-Tabelle
  (`docs/plan/adr/README.md` ADR-0009-Zeile); der ADR-Text erweitert
  ADR-0005 Regel 3 (configyaml **und** report) und nennt R5-Begründung —
  deckt sich mit der arch-check-Lockerung.
- **Spec/Doku-Treue (zwei Formate konsistent):** Spezifikation §„YAML-
  Ausgabe" (`spec/spezifikation.md:732`) beschreibt „dasselbe Dokument wie
  `--json`", camelCase-Feldnamen, gegenseitigen Ausschluss,
  `fixCandidate: null` auch in YAML und Determinismus; operations.md
  (`--yaml`-Zeile + aktualisierte `--doctor`/`--repair`-Zeilen,
  `docs/user/operations.md:25`) und Handbuch §4.11
  (`docs/user/benutzerhandbuch.md:469` ff.) nennen durchgängig „gleiche
  Struktur, nur Serialisierung" und den `--json`/`--yaml`-Ausschluss.
- **codepaths/ids-Sauberkeit der neuen Doku/Spec:** der neue Spec-Block
  (`spec/spezifikation.md:732`) führt IDs nur als verlinkte Backtick-Spans
  (`[`DC-FA-CLI-004`](…)`); die operations.md-/Handbuch-Ergänzungen nennen
  `DC-FA-CLI-004` ausschließlich als verlinkten Backtick-Span. Der
  Handbuch-YAML-Beispielblock steht in einem YAML-Code-Fence
  (`docs/user/benutzerhandbuch.md:476`) — vom ids/codepaths-Scan exempt;
  keine bare `DC-`/`ADR-`/`slice-`-ID außerhalb Link/Fence in den neuen
  Prosa-Zeilen.
- **Slice-Plan-Treue:** die in §3 des Slice-Plans benannten Dateien
  (`finding.go`, `report.go`, `cli.go`, `arch-check.sh`, `spezifikation.md`,
  `operations.md`/`benutzerhandbuch.md`) sind genau die geänderten
  Quell-/Doku-Dateien (plus die im Plan vorgesehenen Tests und die
  ADR-Status-/README-Pflege). Die DoD-*Abhakung* ist nicht Reviewer-Sache.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 2 | 1 |

## Verdikt

**Freigegeben.** Keine HIGH/MEDIUM offen. Der Kern des Slices ist tragfähig
umgesetzt: vollständige Struktur-Parität YAML↔JSON über formatneutrale
Output-Structs mit doppelten json+yaml-Tags (camelCase-Schlüssel gegen das
yaml.v3-Kleinschreiben gesichert), korrektes Embedded-Flattening der
Diagnose via `yaml:",inline"`, explizites `fixCandidate: null` in beiden
Formaten, deterministischer Struct-Marshal (keine Map), korrekt geflushter
YAML-Encoder mit propagiertem Close-Fehler, vollständiger Kombi-Vertrag
(`--json`+`--yaml`/`--repair`+`--yaml` → Exit 2; `--doctor --yaml` erlaubt)
ohne Bruch des Bestands-Tests `TestCLI008_Repair_Inkompatibel`,
verhaltensgleicher render/renderStdout-Refactor und eine exakt auf
`configyaml`+`report` begrenzte arch-check-R3-Lockerung bei yaml-freiem Kern
(R1) und unverletzter R5. Die zwei LOW-Punkte (unvollständige
Vergleichstabelle in §4.9; §4.11-Heading nennt nur `--json`) sind
Doku-Drift ohne Korrektheits-/Verweis-Defekt; INFO-1 ist eine bewusste
Verhaltensangleichung im praktisch unerreichbaren Text-Fehlerpfad. Keiner
blockiert die Freigabe. Die Gate-Bestätigung obliegt der getrennten
Verifikation (hier NICHT als grün angenommen; Tests wurden nicht
ausgeführt).
