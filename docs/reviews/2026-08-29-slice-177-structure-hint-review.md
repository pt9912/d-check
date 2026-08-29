# Review slice-177 — Die zentrale Zusage war über fremden Text brechbar

**Gegenstand:** [slice-177](../plan/planning/in-progress/slice-177-structure-hint.md), Stand `c894c7f`, `d8c89c8`, `e693df9`, `f107003`, `86d2985`.
**Datum:** 2026-08-29. **Reviewer:** unabhängiger Subagent, Skill `.harness/skills/reviewer.md` v1.13.0.
**Eigener Lauf:** `make gates` ⇒ `593 Datei(en) geprüft, 0 Befund(e)`, Coverage 94,70 %, semgrep 0 Findings, Exit 0. **Alle Befunde sind an grünen Gates vorbeigelaufen.**

---

## Urteil

**Blockiert.** HIGH 4 · MEDIUM 3 · LOW 5 · INFO 2. Die Substanz trägt — die
Byte-Identität ohne `message` ist gemessen und hält, der Vorrang ist vollständig
(auch für die Tabellen-Zellen), die Handbuch-Beispiele sind echt, die Abgrenzung
gegen `fixCandidate` ist sauber. Blockierend ist, dass die tragende Zusage über
einen Text gegeben wird, den das Werkzeug nicht kontrolliert.

## HIGH

### F-1 · `unsanitisierter-text-im-zeilenvertrag`

- **quelle:** `DC-FA-CLI-004`, ADR-0073
- **pfad:** `internal/adapter/driven/report/report.go:34-39`, `internal/adapter/driven/configyaml/configyaml.go:442-444`
- **befund:** Der Reporter hängt `f.Message` ungeprüft mit `\t` an; die
  Konfiguration weist nur leer/Whitespace ab. Gemessen gegen das frisch gebaute
  Image: ein `hint` mit `\n` erzeugt aus **einem** Befund **zwei** Zeilen, einer
  mit `\t` **fünf** Felder. Genau das, was der Kommentar im selben Commit
  zusagt und was ADR-0073 gegen die Fortsetzungszeile ins Feld führt. Zweite
  Instanz ohne `hint`: `commits.go:62` setzt `Message: subject` — ein
  Commit-Betreff kann einen Tabulator tragen. Die Änderung macht ein
  unvalidiertes Feld zum Bestandteil eines zeilenorientierten Vertrags; §3.8-Frage
  an den Report-Adapter.
- **verifizierbar:** ja — kein bestehender Test deckt es.

### F-2 · `benannte-grenze-untermengig`

- **quelle:** `DC-FA-STRUCT-001`, ADR-0073
- **pfad:** `internal/hexagon/core/rules/structure.go:100-102`
- **befund:** Lastenheft, Spezifikation, ADR und der `MessageFor`-Kommentar
  zählen **zwei** ausgenommene Befunde auf. Es gibt einen **dritten**: die
  unlesbare **Einzeldatei** meldet über `structureFinding`, also über
  `MessageFor`. Gemessen — mit `hint`: `… section-missing  <Hinweis>`, ohne:
  `… section-missing  Datei ist unlesbar (fail-closed)`. Die fail-closed-Ursache
  verschwindet damit aus Befund-Zeile, `--doctor` **und** `--json`, und an ihrer
  Stelle steht eine Anweisung auf eine nie gemessene Bedingung.
- **verifizierbar:** ja.

### F-3 · `test-ohne-assertion`

- **quelle:** `AGENTS.md` §4 (Harness-Lüge), `DC-FA-CLI-004`
- **pfad:** `internal/adapter/driving/cli/cli_hint_test.go:57-84`
- **befund:** `TestHint_OhneErlaeuterungDreiSpalten` prüft die Spaltenzahl mit
  `t.Logf` statt `t.Errorf`; einzige Assertion ist `Contains("section-missing")`.
  Seine Fixture erzeugt zudem eine **vierspaltige** Zeile. Die zweite Hälfte des
  neuen Akzeptanzkriteriums ist in `make test` ungedeckt.
- **verifizierbar:** ja.

### F-4 · `herkunfts-prosa-im-kommentar`

- **quelle:** `AGENTS.md` §3.7, Baseline §Was ein Kommentar trägt
- **pfad:** `report.go:32`, `:88-90`; `cli_hint_test.go:13-15`, `:57-58`, `:98`
- **befund:** „Sie stand bisher nur in `--json`/`--yaml`", „dort war es schon",
  „byte-identisch zu **vorher**" beschreiben den abwesenden Vorzustand und
  bestehen den Zeitform-Test nicht. Neuzugänge, also keine Bestandsgrenze.
- **verifizierbar:** nein — kein Gate.

## MEDIUM

### F-5 · `mess-label-ungleich-gemessene-menge`

„22 von 31 Regel-Dateien" zählt **Dateien** mit dem Literal `Message:`; acht der
31 sind geteilte Helfer ohne Regel (`run.go`, `scan.go`, `markdown.go` …). Der
reichweiten-relevante Schnitt lautet **21 von 23** Regel-Einstiegs-Dateien;
ungedeckt sind genau `hostpaths` und `spans`. Der Schluss trägt in beiden
Zählungen; die Zahl führt einen späteren Leser in die Irre. Sie steht in einer
`Accepted`-ADR. **verifizierbar:** ja.

### F-6 · `falsche-quelle-fuer-den-entscheid`

`slice-177…md:67-70` begründet den Verzicht auf ein zweites Feld mit ADR-0069.
Zurückgebaut hat die Redundanz **ADR-0070**; ADR-0069 hat die Bedingung
eingeführt und trägt deshalb `Accepted (teil-superseded: ADR-0070)`. ADR-0073
und der Feat-Commit nennen korrekt ADR-0070. `BEO-012`-Klasse.
**verifizierbar:** nein — Link und Anker lösen auf.

### F-7 · `aufzaehlung-nicht-nachgezogen`

`spec/spezifikation.md:2102-2113` und `spec/lastenheft.md:2691-2705` führen die
fail-closed-Ränder als **geschlossene** Aufzählung. Der neue Rand — `hint`
gesetzt, aber leer ⇒ Exit 2 — steht in keiner der beiden.
**verifizierbar:** ja.

## LOW

- **F-8** `--print-config` führt `hint` nicht (`config_template.go:178-199`).
- **F-9** Glossar-Zeile `benutzerhandbuch.md:2221` bleibt dreispaltig.
- **F-10** Der Replay-Harness ankert die vierte Spalte und die `Hinweis:`-Zeile
  nicht (`handbook_examples_test.go:143`, `:157`) — sie können still driften.
- **F-11** „`--json`/`--yaml` bleiben unverändert" gilt nur schema-, nicht
  wertseitig; `TestHint_JSONUnveraendert` assertiert die Änderung, die sein Name
  bestreitet.
- **F-12** `derefString`-Kommentar spricht von „Bedingung aus"; `hint` ist keine
  Bedingung.

## INFO

- **F-13** Der `hint` **ersetzt**: mit ihm verlieren alle Befunde der Regel die
  quantitative Angabe der modul-eigenen Meldung — gemessen an
  `section-cell-undersized`. Weder ADR noch Handbuch nennen das.
- **F-14** `section-ambiguous` und der Spalten-Leerlauf laufen ebenfalls durch
  `MessageFor`; beide sind als Bedingungs-Verletzung verteidigbar, stehen aber
  in der Grenz-Diskussion nicht.

## Negativbefunde (geprüft, ohne Befund)

- **Byte-Identität ohne `message`** — `ghcr.io/pt9912/d-check:v0.66.1` gegen das
  neue Image auf demselben Probe-Repo, Befund `span-unclosed`: `sha256`
  beidseitig identisch.
- **Kein Befund-Pfad gewinnt oder verliert ein `message`** — vier
  `structure`-Befunde alt/neu verglichen: identische Menge, Grund-Codes, Zeilen.
- **Vorrang in `structure_tablecell.go`** vollständig; `structure_tableorder.go`
  geht durch `structureFinding`, ebenfalls gedeckt.
- **Exit 2 am Config-Rand** für `''` und `'   '` gemessen.
- **`--doctor`-Reihenfolge** `Stelle:` → `Hinweis:` → `Fix-Kandidat:`, konsistent
  mit `DC-FA-CLI-007.a`; die Schritt-Neunummerierung ist in allen Rückverweisen
  nachgezogen.
- **`DC-FA-CLI-004` Happy Path** — zwei kaputte Links liefern genau zwei Zeilen.
- **Handbuch-Beispiele gemessen**, nicht erfunden.
- **§3.4, §3.5, §3.2, §3.6, Hexagon-Richtung, `DC-QA-03`, §3.8 der Scan-Menge,
  Slice-Plan-Form, Zustandsfelder** — je geprüft, kein Befund.
