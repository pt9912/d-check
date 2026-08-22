# Review — slice-029 Implementierung (`--doctor --json`)

## Kopf-Metadaten

- **Review-Art:** Code-Review (Diff gegen Plan/Spec/Anforderungen/Hard
  Rules — kein Verifier; DoD-Abhaken und Gate-Lauf-Bestätigung sind
  nicht Gegenstand).
- **Datum:** 2026-06-19
- **Gegenstand:** Working-Tree-Diff (unstaged) der slice-029-Umsetzung
  — maschinenlesbare Diagnose `--doctor --json`: neuer
  `report.DoctorJSON`, umgestellter CLI-Kombi-Check + render-Switch,
  ersetzte/ergänzte Akzeptanztests, Spec/Handbuch/Operations/Changelog.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md`
  v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/done/slice-029-doctor-json.md` (zum Review-Zeitpunkt
  unter `open/`); Anforderungen
  `DC-FA-CLI-007` (Haupt — JSON-Rendering der Diagnose),
  `DC-FA-CLI-004` (Ausgabeformate/Kombinierbarkeit), `DC-FA-CLI-003`
  (Exit-Codes), `DC-QA-02` (Determinismus), `DC-QA-03`
  (read-only/netzlos); Spezifikation §`DC-FA-CLI-007.a` Schritt 6 + §2
  „JSON-Diagnose"; Hard Rules `AGENTS.md` §3; ADR-0005 (Hexagon-Import-
  Regeln); MR-006 (Referenzrichtung). **Die DoD-Abhakung des Slices lag
  bewusst nicht vor** — die Findings sind eigenständig gebildet.

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

Keine.

### INFO

#### INFO-1 — `message`-Feld im JSON-Diagnose-Schema undokumentiert

- **Kategorie:** INFO
- **Quelle:** Maintainability (Spec-Treue der Feldliste)
- **Pfad:** `internal/adapter/driven/report/report.go` Z. 165–169
  (eingebettetes `core.Finding`) · `spec/spezifikation.md` §2 „JSON-Diagnose"
- **Befund:** Das eingebettete `core.Finding` promotet auch sein
  `Message`-Feld (`json:"message,omitempty"`,
  `internal/hexagon/core/model/finding.go:27`) flach in jeden
  `findings`-Eintrag; das §2-Delta-Schema der JSON-Diagnose listet
  `message` weder im `required`-Block noch in den `properties`. Ist
  ein Befund mit nicht-leerer `Message` vorhanden, erscheint im
  stdout-Dokument ein Feld, das das dokumentierte Diagnose-Schema nicht
  nennt — ein Schema-Konsument, der strikt gegen die §2-Liste
  validiert, sähe ein unerwartetes Feld. Wegen `omitempty` bleibt es im
  Normalfall (leere Message) aus, daher INFO statt MEDIUM.
- **Verifizierbar:** ja — ein Akzeptanztest mit einem Befund mit
  gesetzter `Message` würde das Feld im JSON sichtbar machen; ein
  Schema-Diff-Test gegen die §2-Properties würde die Abweichung melden.

#### INFO-2 — Kein dedizierter DC-QA-03-Test für den doctor-json-Pfad

- **Kategorie:** INFO
- **Quelle:** DC-QA-03 (Seiteneffektfreiheit) — Slice-DoD nennt
  „read-only (DC-QA-03)" für die JSON-Variante
- **Pfad:** `internal/adapter/driving/cli/cli_acceptance_test.go`
  (neue Tests `TestCLI007_DoctorJSON_Happy`,
  `TestCLI007_DoctorJSON_Determinismus`)
- **Befund:** Die beiden neuen Tests prüfen Happy-Path und Determinismus
  der JSON-Diagnose, aber keiner stellt explizit fest, dass der
  `--doctor --json`-Lauf das Repo nicht schreibt (z. B. read-only-Mount
  oder mtime-/Inhalts-Vergleich vor/nach). Die Read-only-Eigenschaft ist
  strukturell gegeben (derselbe `core.Run`/`render`-Pfad wie die anderen
  Modi; `DoctorJSON` schreibt nur in den injizierten `stdout`), daher
  kein Korrektheitsdefekt, sondern eine nicht abgesicherte DoD-Zusage.
- **Verifizierbar:** ja — ein read-only-Lauf-Test (wie für die anderen
  Modi vorhanden) würde die Zusage decken; sein Fehlen ist im
  Test-Bestand sichtbar.

## Negativbefunde (geprüft, ohne Befund)

- **Determinismus der JSON-Diagnose (DC-QA-02):** `jsonDiagFinding`
  bettet `core.Finding` ein und ergänzt zwei Felder; alle Felder sind
  Skalare in fester Struct-Deklarationsreihenfolge, `Summary` zwei
  Ints — keine Map im gesamten Renderpfad
  (`report.go:165-197`). Iteriert wird über die bereits stabil
  sortierte `res.Findings` (`SortFindings`,
  `internal/hexagon/core/model/finding.go:33`). Go-`encoding/json` gibt
  promotete Felder eines eingebetteten Structs in Deklarationsreihenfolge
  aus — keine nicht-deterministische Quelle. Der Test
  `TestCLI007_DoctorJSON_Determinismus` (10× byte-identisch) belegt es
  empirisch.
- **`fixCandidate: null` explizit (kein Weglassen):**
  `FixCandidate *jsonFixCandidate` trägt KEIN `omitempty`
  (`report.go:168`); fehlt der Kandidat, marshalt ein nil-Pointer als
  `null`. Der Test prüft `strings.Contains(stdout, "\"fixCandidate\":
  null")` (`cli_acceptance_test.go:457`). Die Aussage „kein eindeutiger
  Fix" geht nicht verloren.
- **Kombinations-Vertrag (DC-FA-CLI-004):** `parseOptions`
  (`cli.go:159-166`) weist `opts.json && opts.repair` und
  `opts.doctor && opts.repair` je mit Exit 2 ab — jede Kombination mit
  `--repair` wird gefangen, auch das Tripel `--doctor --json --repair`.
  Neu erlaubt ist ausschließlich `--doctor && --json && !repair`.
  `TestCLI008_Repair_Inkompatibel` (`cli_acceptance_test.go:570-580`)
  deckt `--repair --json` (Exit 2, „nicht mit --json") und
  `--doctor --repair` (Exit 2, „nicht kombinierbar") weiter ab; die in
  cli.go geänderte Fehlermeldung („--repair ist nicht mit --json
  kombinierbar") erfüllt die Test-Assertion `Contains(…, "nicht mit
  --json")` weiterhin.
- **stdout-Reinheit:** `DoctorJSON` schreibt allein das JSON-Dokument
  auf `stdout`; die Zusammenfassung steckt im `summary`-Feld, nicht auf
  stderr (`report.go:185-197`). Der Test prüft Parsbarkeit per
  `json.Unmarshal([]byte(stdout))` und schließt Prosa-Marker
  („Diagnose", „Fix-Kandidat:") aus (`cli_acceptance_test.go:459-464`).
- **render-Switch-Präzedenz:** `case opts.doctor && opts.json` steht VOR
  `case opts.doctor` und `case opts.json` (`cli.go:225-249`); die
  JSON-Diagnose verdrängt korrekt sowohl die Prosa- als auch die
  knappe-JSON-Variante. Der Render-Fehler wird auf Exit 2 abgebildet
  wie die übrigen Zweige.
- **Hexagon (ADR-0005, Import-Regeln):** `DoctorJSON` liegt im
  driven-Adapter `report`; das Paket importiert nur `encoding/json`,
  `fmt`, `io` (stdlib) und `internal/hexagon/core`
  (`report.go:6-12`) — kein Adapter-Import, kein I/O-API außer dem
  injizierten Writer. Der Renderer liegt im richtigen Layer; der Kern
  (`core.ReasonText`, `core.FixCandidateFor`) wird nur konsumiert.
- **Spec-Treue & Referenzrichtung (MR-006):** Die neuen Spec-Stellen
  (§`DC-FA-CLI-007.a` Schritt 6, §2 „JSON-Diagnose") verweisen nur
  aufwärts/seitwärts auf `lastenheft.md`-Anforderungen und auf andere
  Spec-Abschnitte (`#spec-002--json-ausgabe---json`,
  `#dc-fa-cli-007a--diagnose-modus`), nicht abwärts auf ADRs/Slices im
  bindenden Text. Die Slice-Spalte der Änderungshistorie (`slice-029`,
  Zeile 811) ist ein **unverlinkter Klartext-Label** im selben Muster
  wie die zwölf Vorzeilen — MR-006 verbietet abwärts-*Links* auf
  Planning-Artefakte; das ist keiner. Keine neue MR-006-Bedingung.
- **Anker-/Link-Auflösung:** Alle neuen Querverweise lösen auf —
  Spec-Heading `### JSON-Diagnose (\`--doctor --json\`)` →
  `#spec-003--json-diagnose---doctor---json` (referenziert in der Historie);
  `### JSON-Ausgabe (\`--json\`)` → `#spec-002--json-ausgabe---json`
  (Schritt 6 + §2); Handbuch `### 4.11 …` →
  `#411-maschinenlesbare-ausgabe---json` (§4.9-Beispiel);
  `#dc-fa-cli-007a--diagnose-modus`. Kein toter Anker.
- **Doku-Konsistenz der Kombi-Aussagen:** Die „nicht mit
  --json"-Aussage zu `--doctor` ist in Handbuch §4.9
  (`benutzerhandbuch.md:316-318`) und in der Optionstabelle von
  `operations.md` (Zeile mit `--doctor`) auf den neuen Vertrag
  umgestellt; die FAQ §8 enthält keine `--doctor`+`--json`-Verbots-
  Aussage (entgegen der Slice-Plan-Sorge — keine Drift offen). Die
  weiterhin gültigen `--repair`-Verbote bleiben unverfälscht: Handbuch
  §4.10 Zeile 392 („Nicht mit `--json` oder `--doctor` kombinierbar")
  gehört zu `--repair` und wurde korrekt nicht angefasst; `operations.md`
  `--repair`-Zeile unverändert; Spec §`DC-FA-CLI-008.a` `--repair`-Verbot
  unverändert.
- **Negativtest-Bestand:** Der alte
  `TestCLI007_Doctor_JSONInkompatibel` (der `--doctor --json` als
  Nutzungsfehler prüfte) wurde durch die JSON-Variante-Tests ersetzt —
  korrekt, da die Aussage sich umgekehrt hat. Die beiden weiter
  gültigen Negativfälle bleiben in `TestCLI008_Repair_Inkompatibel`
  gedeckt.
- **Changelog-Eintrag:** Added (neuer Modus) + Changed (`--doctor`
  jetzt mit `--json` kombinierbar) im `[Unreleased]`-Block, mit
  korrekten Anforderungs-Links; konsistent zur Spec-/Code-Aussage.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 0 | 2 |

## Verdikt

**Freigegeben.** Keine HIGH/MEDIUM/LOW offen. Der Kern des Slices —
deterministischer JSON-Renderer mit explizitem `fixCandidate: null`,
korrekt verengter Kombinations-Vertrag (nur `--doctor --json` neu
erlaubt, alle `--repair`-Kombinationen weiter Exit 2),
Layer-konformer Renderer, konsistente Spec/Doku — ist tragfähig
umgesetzt. Die zwei INFO-Punkte (undokumentiertes promotetes
`message`-Feld im §2-Schema; fehlender dedizierter
DC-QA-03-doctor-json-Test) sind dokumentationswürdige Annahmen ohne
Korrektheits- oder Sicherheitswirkung und blockieren die Freigabe
nicht. Die Gate-Bestätigung obliegt der getrennten Verifikation.
