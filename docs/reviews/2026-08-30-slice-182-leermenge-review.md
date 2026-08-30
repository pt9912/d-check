# Review slice-182 — Eine erklärte Teilmenge darf die Menge leeren, wenn sie sagt, wie viele


**Review-Art:** Code/Design (gegen Slice-Plan, ADR, Hard Rules)
**Gegenstand:** `a062fe8^..HEAD` = `a062fe8` (Claim-Move) · `9c383ba` (feat) · `307311c` (Handbuch) · `47b0ccb` (CR-Antwort); der eingehende CR selbst ist **nicht** Gegenstand
**Skill:** `.harness/skills/reviewer.md` @ v1.13.0 · **Modell:** claude-opus-5[1m] · **Datum:** 2026-08-30
**Eingangs-Kontext:** `AGENTS.md` §2/§3.1/§3.2/§3.4/§3.6/§3.7/§3.8/§4/§5, `harness/conventions.md` (MR-013, MR-053, MR-054), `.harness/baseline/v5.12.0/regelwerk/` (`grundlagen-harness-dateien.md` §Was ein Kommentar trägt, `modul-04-adrs.md`, `modul-05-planning-harness.md`, `modul-10-review-harness.md`), `DC-FA-STRUCT-001`, `DC-FA-CLI-003`, `DC-FA-CLI-007`, `DC-FA-CLI-010`, ADR-0073, ADR-0075, ADR-0078, Beobachtungs-Register (BEO-009, BEO-012, BEO-013, BEO-020, BEO-023), Nutzer-Vorgabe „CR-Dokumente tragen Bitte und Beleg"

### Eigener Lauf

| Lauf | Ausgabe |
|---|---|
| `make gates` | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` · Exit 0 |
| `make doc-check` | `d-check: 618 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `make test` | alle Pakete `ok` · Exit 0 |
| `make trace-check RANGE=a062fe8^..HEAD` | `d-check: 618 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `make adr-check RANGE=a062fe8^..HEAD` | `d-check: 618 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `d-check --print-mk` (Produkt, im Image) | drei Zeilen mit `doc-doctor`/`--doctor` (Zeilen 37–39 der Ausgabe) |
| 4 Mutationsproben an `structure.go` / `configyaml.go`, je `make test` | s. u. — zwei davon bleiben **grün** |
| 2 Produkt-Proben mit temporärem `--config`-Profil | s. u. (Arbeitsbaum danach sauber, `git status --short` leer) |

`make fullbuild`, `make record-gates` und `make image-scan` bewusst **nicht** gefahren (Nachweis-Datei bzw. Netz).

---

## Urteil

**BLOCKIERT** — 4 HIGH · 4 MEDIUM · 3 LOW · 2 INFO.

Der Entwurf selbst ist gut und die Zerlegung trägt: die Zahl statt der Erlaubnis ist die
stärkere Antwort auf denselben Antrag, die Trennung Konfigurationsdefekt ↔ Bestandszustand
ist sauber gezogen, und alle sechs Deklarations-Oberflächen sind bedient. Blockierend ist
etwas anderes: **die tragende Messung, mit der wir dem Absender seine Form absagen, ist
falsch** — `--print-mk` verteilt sehr wohl ein Doctor-Target, und dasselbe Lastenheft
schreibt es in einem Akzeptanzkriterium vor. Dazu ein Kollateralschaden am
Spezifikations-Stratum, den kein Gate sieht, zwei Zusagen des neuen Schlüssels ohne
fangende Probe (die ADR behauptet ausdrücklich das Gegenteil) und ein Kommentar mit
rückwärts zeigendem Konjunktiv.

---

## Findings

### H-1 · „`--print-mk` verteilt kein Doctor-Target" ist falsch — und trägt die Abweichung vom Antrag

- **quelle:** `AGENTS.md` §5 (Proben-Disziplin: *„eine genannte Probe muss gelaufen sein"*), BEO-020, Reviewer-Skill HIGH-nahe Hälfte von BEO-009
- **pfad:** `spec/lastenheft.md:3364` (Historie 0.78.0) · `docs/plan/adr/0078-erklaerte-leermenge-mit-zahl.md:62-66` (Entscheidung 2) · `docs/plan/planning/in-progress/slice-182-erklaerte-leermenge.md:60-63` · `docs/plan/cr/2026-08-30-antwort-a-check-leermenge.md:61` (Tabellenzeile) · Commit-Botschaft `9c383ba`
- **zugesagt:** wörtlich, an drei der vier Stellen gleichlautend: *„`--doctor` läuft in **keinem** Gate dieses Repos — nicht im `Makefile`, nicht in `.github/workflows/`, nicht im `pre-commit`-Hook —, und `--print-mk` verteilt kein Doctor-Target an Konsumenten."* Daraus folgt der ganze Rest: *„Gemessen ist er null"*, und deshalb wird die beantragte Form abgelehnt.
- **gemessen:** die erste Hälfte stimmt (`grep -n doctor Makefile` → 0, `.github/workflows/` → 0, `.githooks/` → 0). **Die zweite ist falsch.** Produkt-Lauf `d-check --print-mk` im Image:

  ```
  37:.PHONY: doc-doctor
  38:doc-doctor: ## erklärende Diagnose mit Fix-Kandidaten (DC-FA-CLI-007)
  39:	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --doctor
  ```

  Quelle: `internal/adapter/driving/cli/print_mk.go:87-89`; der Doc-Kommentar `:22-27`
  führt `doc-doctor` als eines der **zwölf** `##`-annotierten Targets. Und die Verteilung
  ist nicht Zufall, sondern **Vertrag**: `spec/lastenheft.md:516` (`DC-FA-CLI-010`, Happy
  Path) verlangt genau dieses Target, `:517` (Boundary) verlangt `doc-doctor` → `--doctor`,
  und die Lastenheft-Historie `:3439` (0.27.0) protokolliert seine Aufnahme. Damit
  widerspricht die 0.78.0-Historie einem Akzeptanzkriterium **derselben Datei**.
- **die Adressaten-Hälfte ist repo-lokal belegt:** `docs/plan/adr/0029-arch-check-via-a-check.md:56`
  — *„konsumiert a-check bereits d-check über ein eingebundenes `d-check.mk`"*. Der
  Absender, dem wir schreiben, *hat* also ein `make doc-doctor` im eigenen Makefile
  stehen. Genau die Reichweite, die das Argument bestreitet.
- **die Fassung in der CR-Antwort ist enger, aber nicht ehrlicher:** die Zelle lautet
  *„kein `doc-doctor` im Gate-Pfad"* und steht in einer Tabelle, deren Spaltenkopf
  `--doctor` heißt und deren drei andere Zellen `0`/`0 Fundstellen` sagen. Sie liest sich
  als Abwesenheit, während die wahre Zelle *„ein verteiltes Target `doc-doctor`, das
  `--doctor` fährt"* lautete.
- **warum es zählt:** es ist nicht irgendeine Zeile, sondern *die* Begründung, mit der wir
  einem Adopter seine selbst gemessene Form absagen — und sie steht im **ranghöchsten**
  Stratum (Lastenheft, Rang 1) sowie in einer `Accepted`-ADR, die nach `AGENTS.md` §3.5
  nicht mehr inhaltlich korrigiert werden darf. Die ADR trägt zwar den passenden
  Re-Evaluierungs-Trigger (*„Wenn ein Adopter `--doctor` in einem Gate fährt"*), aber der
  ist auf eine Auskunft des Absenders gestellt, nicht auf die Tatsache, dass wir ihm das
  Target selbst ausliefern.
- **verifizierbar:** ja, ohne Gate-Lauf — `d-check --print-mk` gegen `spec/lastenheft.md:516`.
- **klasse:** `messung-widerlegt-durch-eigenen-vertrag`

### H-2 · Der Feat-Commit hat die Kopfzeile der Tabelle §6 *Externe Verträge* überschrieben

- **quelle:** `AGENTS.md` §2 (Source Precedence, Rang 3 der Spec-Straten), Reviewer-Skill HIGH *Korrektheits-Verstoß*
- **pfad:** `spec/spezifikation.md:3050-3056`; eingeführt von `9c383ba`
- **zugesagt:** §6 ist eine vierspaltige Vertragstabelle mit dem Kopf
  `| Kennung | System | Version/Stand | Vertrag |`.
- **gemessen:** `git show 9c383ba -- spec/spezifikation.md` enthält
  `-| Kennung | System | Version/Stand | Vertrag |` — die Kopfzeile ist **gelöscht** und
  durch den 0.78.0-Historie-Eintrag ersetzt worden, der zusätzlich korrekt in §7 steht
  (`:3062`). Aktueller Stand:

  ```
  3050:## 6. Externe Verträge
  3052:| 2026-08-30 | §DC-FA-STRUCT-001.a Schritt 3 und das §2-Schema … |
  3053:|---|---|---|---|
  3054:| SPEC-064 | gopkg.in/yaml.v3 | gepinnt via go.sum | … |
  ```

  Die drei externen Verträge `SPEC-064`/`SPEC-065`/`SPEC-066` stehen damit unter einer
  Kopfzeile, die ein Datum und einen Änderungstext trägt; der Eintrag steht doppelt im
  Dokument.
- **kein Gate sieht das, gemessen:** `make gates` → Exit 0, `make doc-check` →
  `618 Datei(en) geprüft, 0 Befund(e)`. Das `structure`-Profil dieses Repos führt für
  §6 weder eine Spalten- noch eine Überschriften-Bedingung; `links`/`anchors` stören sich
  nicht an einer Kopfzelle.
- **warum es zählt:** die Spezifikation ist Rang 2 der kanonischen Quellen und §6 ist die
  Stelle, an der ein Implementer nachschlägt, welche fremden Systeme mit welchem Stand
  gebunden sind. Der Schaden ist still — und er ist in demselben Commit entstanden, der
  die Historie-Zeile schreiben wollte, also durch ein Einfügen an der falschen
  Tabellen-Position.
- **verifizierbar:** nein — kein Sensor; Beleg ist der Diff.
- **klasse:** `tabellenkopf-ueberschrieben`

### H-3 · Zwei Zusagen des neuen Schlüssels haben keine Probe, die sie fängt — die ADR behauptet das Gegenteil

- **quelle:** BEO-023 (Zähler **5**, Schwelle erreicht), Slice-DoD *„Umkehr-Proben je Zusage, jede von den Tests gefangen, die dagegen stehen"*, Reviewer-Skill HIGH #1 (stiller Grün-Pfad)
- **pfad:** `docs/plan/adr/0078-erklaerte-leermenge-mit-zahl.md:154-161` (Fitness Function) · `internal/hexagon/core/rules/structure_leermenge_test.go` (sieben Tests) · `internal/adapter/driven/configyaml/configyaml_test.go:264-266` (zwei Ränder)
- **zugesagt:** ADR-0078, wörtlich: *„**Jede Zusage ist von der Mutation gefangen, gegen die sie steht**, gemessen: Zählung deaktiviert ⇒ 5 rot · Drift nur einseitig ⇒ 3 rot · `case`-Guard entfernt … ⇒ 6 rot · Mismatch als `raw` ⇒ 1 rot · beide Config-Ränder entfernt ⇒ 1 rot."*
- **gemessen, Probe (a) — die Verdrahtung Config → Modell ist ungetestet:** in
  `internal/adapter/driven/configyaml/configyaml.go:497` die Zeile
  `ExemptExpectCount:    r.ExemptExpectCount,` aus `applyStructureRule` entfernt. Danach
  wird der Schlüssel weiterhin geparst und validiert, erreicht aber **nie** die Regel: die
  Nullmengen-Härte feuert wieder, `exempt-expect-count: 19` liefert `section-missing`, das
  Feature ist für jeden Anwender tot. `make test` → alle Pakete `ok`, **Exit 0**, kein
  einziger Fehlschlag. Grund: alle sieben Verhaltens-Tests bauen `model.StructureRule`
  direkt; die beiden `configyaml`-Tests sind **reine Negativ-Tests** (`ohne Muster`,
  `negativ`) und prüfen nur, dass `Decode` einen Fehler liefert. Ein Positiv-Test, der
  einen gültigen Wert bis ins Modell verfolgt, existiert nicht
  (`grep -rn exempt-expect-count --include=*_test.go internal/` → genau die zwei
  Negativ-Zeilen).
- **gemessen, Probe (b) — die Zeiger-Semantik der deklarierten Null ist ungetestet:**
  `internal/hexagon/core/rules/structure.go:137` von
  `case r.ExemptExpectCount != nil:` auf
  `case r.ExemptExpectCount != nil && *r.ExemptExpectCount > 0:` geändert — das ist genau
  die Mutation, die eine deklarierte `0` wieder von *„nicht deklariert"* ununterscheidbar
  macht und damit die `*int`-Wahl entwertet. `make test` → **Exit 0**, kein Fehlschlag.
  Auch `TestExemptExpectCount_NullIstEineAussage` bleibt grün: sein Fixture nimmt **null**
  Abschnitte aus, und dieselben drei `section-marker-missing` entstehen mit und ohne
  Feature. Der Test misst die Zusage nicht, gegen die sein Kommentar steht.
- **die aufgezählten Mutationen reproduzieren dagegen:** *Zählung deaktiviert*
  (`case … && n < 0`) ⇒ 5 rote Blatt-Tests (`DeklarierteLeermengeIstStumm`,
  `DriftIstBeidseitig/zu_wenig`, `…/zu_viel`, `PrueftAuchBeiRestmenge`,
  `HintGiltFuerDenMismatch`) — exakt die behaupteten 5. *Drift nur einseitig*
  (`!=` → `<`) ⇒ 2 Blatt-Tests / 3 FAIL-Zeilen; die ADR sagt 3, zählt hier also den
  Eltern-Test mit, bei der ersten Mutation nicht. Die Zählweise ist uneinheitlich, die
  Größenordnung stimmt.
- **warum es zählt:** die ADR ist `Accepted` und ihre Fitness Function ist eine
  Allaussage (*„jede Zusage"*), die zwei Gegenbeispiele hat — davon eines, das das
  gesamte Feature stumm sterben ließe, ohne dass ein Gate mucks. Das ist die Klasse, die
  BEO-023 bei Zähler 5 führt, und die zweite Richtung von BEO-020 in einem Zug.
- **verifizierbar:** ja — die beiden Mutationen sind reproduzierbar und in ~40 s je Lauf
  messbar.
- **klasse:** `zusage-ohne-fangende-probe`

### H-4 · Konjunktiv über die verworfene Alternative plus Mess-Label im Kommentar

- **quelle:** `AGENTS.md` §3.7, Baseline `grundlagen-harness-dateien.md` §Was ein Kommentar trägt (Zeitform-Test; *„Drei Klassen fallen heraus … **Deliberation** (der Konjunktiv über die verworfene Alternative …; ihr Ort ist die ADR oder §3/§6 des Slice)"*)
- **pfad:** `internal/hexagon/core/rules/structure.go:157-158`
- **zugesagt:** derselbe Kanon verbietet zusätzlich Mess-Labels und Review-Historie im
  Kommentar; Herkunft nur als **ein** auflösbares Feld.
- **gemessen:** der Satz lautet *„Ein `return nil` hier waere toter Code (gemessen: keine
  Mutation daran macht einen Test rot)."* Das ist beides zugleich: ein **rückwärts**
  zeigender Konjunktiv über eine Variante, die es einmal gab und die entfernt wurde
  (die ADR `:163-167` und die Commit-Botschaft `9c383ba` erzählen dieselbe Geschichte —
  dort gehört sie hin), und ein **Mess-Label** über einen Mutationslauf. Der Rest des
  Blocks (`:153-157`, *„KEIN Zweig für die geleerte Menge … weil dieser case die
  Nullmengen-Haerte darunter UEBERSPRINGT"*) ist eine saubere Kopplung/Abgrenzung; der
  Befund gilt nur dem letzten Satz.
- **warum es zählt:** die Zeile reist mit dem Code weiter und läuft bei jeder künftigen
  Änderung erneut gegen den Zeitform-Test; das Mess-Label altert zusätzlich — es
  behauptet ein Mutationsergebnis, das nach der nächsten Test-Ergänzung nicht mehr
  stimmen muss. Dieselbe Klasse war in slice-181 bereits ein HIGH.
- **verifizierbar:** nein — kein Gate; `make lint` meldet `0 issues` (gemessen).
- **klasse:** `kommentar-konjunktiv-rueckwaerts`

### M-1 · Die erwartete Anzahl gilt **je Datei**, jede Deklarations-Oberfläche beschreibt sie je Regel

- **quelle:** `AGENTS.md` §3.8 (Zusagen nur über die gemessene Menge), Reviewer-Skill MEDIUM *Spec-Treue-Lücke einer Messmethode*
- **pfad:** `internal/hexagon/core/rules/structure.go:112-165` (die Zählung liegt in `checkStructureFile`, nicht in `checkStructureRule`) · `spec/lastenheft.md:2532` · `spec/spezifikation.md:2166-2172`, `:2902` · `docs/plan/adr/0078-erklaerte-leermenge-mit-zahl.md:44-47` · `internal/adapter/driving/cli/config_template.go:207-210` · `docs/user/benutzerhandbuch.md:2016-2018`, `:2138-2152`
- **zugesagt:** durchgängig auf Regel-Ebene formuliert — Spezifikation §2: *„**erwartete Anzahl** der von `exempt-section-pattern` abgezogenen Abschnitte"*; Handbuch: *„`exempt-expect-count: 19` heißt: „neunzehn sollen ausgenommen sein""*. Kein Satz sagt „je Datei".
- **gemessen (Produkt, temporäres `--config`-Profil, danach entfernt):** eine Regel
  `files: "docs/plan/cr/*.md"` · `section-pattern: "^## "` · `sections: each` ·
  `exempt-section-pattern: "^## Vorab"` · `exempt-expect-count: 1` liefert

  ```
  d-check: 556 Datei(en) geprüft, 3 Befund(e)
  docs/plan/cr/2026-08-25-cr-regelwerk-v5110.md:1 … section-exempt-mismatch  … nimmt 0 von 4 Abschnitten aus, deklariert sind 1
  docs/plan/cr/2026-08-30-cr-a-check-leermenge.md:1 … section-exempt-mismatch  … nimmt 0 von 3 Abschnitten aus, deklariert sind 1
  docs/plan/cr/2026-08-30-cr-a-check-structure-teilmenge.md:1 … section-exempt-mismatch … nimmt 0 von 4 Abschnitten aus, deklariert sind 1
  ```

  Die Zahl muss also in **jeder** getroffenen Datei stimmen. Für den Fall des Absenders
  (`files: "spec/lastenheft.md"`, eine Datei) fällt das nicht auf; für jedes Glob mit
  mehr als einem Treffer ist es eine andere Zusage als die dokumentierte.
- **ungeprüft:** alle sieben neuen Tests laufen über `laufe(…)`, und der Helfer
  (`internal/hexagon/core/rules/structure_teilmenge_test.go:28-31`) legt ein MemFS mit
  **genau einer** Datei `docs/a.md` an. Die Mehrdatei-Achse ist nirgends getestet.
- **warum es zählt:** ein Adopter, der die Fähigkeit über ein Verzeichnis-Glob einsetzt —
  der naheliegende zweite Anwendungsfall, den die ADR unter §Re-Evaluierungs-Trigger
  selbst erwartet —, bekommt Befunde, die keine Deklarations-Oberfläche ankündigt.
- **verifizierbar:** ja — die Probe oben.
- **klasse:** `zusage-skopus-datei-vs-regel`

### M-2 · Ein `section-exempt-mismatch` bricht die restliche Prüfung dieser Datei ab — undeklariert

- **quelle:** `AGENTS.md` §3.6 (die geprüfte Menge zu verkleinern ist ein bewusster Akt), Reviewer-Skill MEDIUM *Spec-Treue-Lücke*
- **pfad:** `internal/hexagon/core/rules/structure.go:147-152` (`return []model.Finding{…}`)
- **zugesagt:** Spezifikation `:2166-2172` und Lastenheft `:2571-2576` sagen, dass bei
  Abweichung *„`section-exempt-mismatch` entsteht"* — nicht, dass danach nichts mehr
  gemessen wird. Der Nachbar-Abbruch bei `sections: one` ist im Code ausdrücklich
  begründet (*„ohne eindeutigen Abschnitt sagt eine Messung nichts"*); hier fehlt eine
  solche Begründung, denn die verbleibenden Abschnitte sind sehr wohl messbar.
- **gemessen:** der eigene Test belegt es.
  `internal/hexagon/core/rules/structure_leermenge_test.go:75-88`: bei
  `ExemptExpectCount = ptr(3)` und einem vierten, nicht ausgenommenen Abschnitt entsteht
  genau **ein** Befund `section-marker-missing`; bei `ptr(4)` entsteht genau **ein**
  Befund `section-exempt-mismatch` — der Marken-Befund des vierten Abschnitts ist weg. Die
  Assertion `len(f) != 1` fixiert diese Maskierung als gewolltes Verhalten.
- **warum es zählt:** wer die Zahl korrigiert, sieht schlagartig Befunde auftauchen, die
  vorher verdeckt waren, und liest das als Folge seiner Korrektur. Der Exit-Code bleibt in
  beiden Fällen 1, ein stilles Grün entsteht also nicht — die Diagnose wird aber unvollständig,
  und das steht in keiner der sechs Oberflächen.
- **verifizierbar:** ja — Testquelle und Code.
- **klasse:** `frueher-abbruch-undeklariert`

### M-3 · Die neue Zeile steht in der Tabelle *Bedingungen im Abschnitt* und erbt deren Zusagen

- **quelle:** Reviewer-Skill MEDIUM *Spec-Treue-Lücke einer Messmethode*, `AGENTS.md` §3.8
- **pfad:** `spec/lastenheft.md:2510-2532`
- **zugesagt:** der Tabellen-Vorspann `:2510-2515` lautet: *„**Bedingungen im Abschnitt**,
  **je Abschnitt geprüft**, jede optional, jede **fence-treu** (der Abschnitts-Text wird
  zuvor um Fenced-Code bereinigt; **zwei** Bedingungen lesen einen anderen Text und sind
  unten je benannt …) — und jede mit **eigenem** Grund-Code"*. Die neue Zeile `:2532`
  hängt darunter.
- **gemessen:** `exempt-expect-count` ist keine Bedingung im Abschnitt. Es wird **einmal
  je Datei** ausgewertet (`structure.go:132-152`), **bevor** irgendein Abschnitts-Text
  gelesen wird; es läuft nicht durch `structureConditions`, es liest weder den
  bereinigten noch den rohen Abschnitts-Text, und sein Befund steht per Zusage auf
  `line = 1` statt auf der Abschnitts-Zeile — im Widerspruch zur Zusage der
  0.72.0-Historie, dass zeilen-gebundene Bedingungen *„auf der Zeile melden, an der
  repariert wird"*. Nebenbei bricht die Zeile die Enumeration im Vorspann: sie führt
  „zwei Bedingungen, die einen anderen Text lesen" abschließend auf; die neue liest gar
  keinen.
- **warum es zählt:** dieselbe Enumerations-Drift, die die Lastenheft-Historie für andere
  Stellen schon zweimal protokolliert hat. Die richtige Nachbarschaft wäre der
  darunterliegende Absatz *„Eine Regel darf ihre Grundmenge erklären"* (`:2535-2537`),
  in dem der Schlüssel bereits inhaltlich behandelt wird (`:2569-2583`).
- **verifizierbar:** nein — kein Sensor über Prosa-Enumerationen.
- **klasse:** `spec-zeile-in-falscher-klasse`

### M-4 · „169 Befunde, `diff` leer" ist aus den Artefakten nicht nachvollziehbar

- **quelle:** `AGENTS.md` §5 (*„eine genannte Probe muss gelaufen sein … und ihr Schluss reicht nicht weiter als die gemessene Menge"*), BEO-020
- **pfad:** `docs/plan/adr/0078-erklaerte-leermenge-mit-zahl.md:136-137` · `docs/plan/cr/2026-08-30-antwort-a-check-leermenge.md:110-112`
- **zugesagt:** *„Ohne den Schlüssel ist der Befundsatz byte-identisch — gemessen gegen
  das Vorgänger-Image: **169 Befunde, `diff` leer**."* In der CR-Antwort steht dieselbe
  Zahl unter *Angenommen, unverändert* — also in einem Dokument, das an einen Dritten geht.
- **gemessen:** weder die ADR noch die CR-Antwort nennt Korpus, Konfiguration oder Befehl.
  Erst die Commit-Botschaft `9c383ba` präzisiert *„derselbe max-tasks-Lauf"* — ohne Profil.
  Die naheliegenden Kandidaten liefern die Zahl nicht: `make doc-check` → **0**, der
  Closure-Lauf `--config .d-check.closure.yml --enable planning --enable structure
  --enable spans` → **0**, ein Lauf mit sieben zusätzlich aktivierten Modulen → **0**, und
  eine Rekonstruktion des `max-tasks`-Laufs (`files: "docs/plan/planning/**/*.md"` ·
  `section: "## 4. Definition of Done"` · `max-tasks: 3`) liefert **218**, nicht 169.
- **warum es zählt:** *byte-identisch* ist die zentrale Rückwärtskompatibilitäts-Zusage
  dieses Slice, und der Absender bekommt sie mit einer Zahl belegt, die er nicht
  nachrechnen kann. Ich behaupte **nicht**, dass die Messung nicht stattfand — nur, dass
  sie in dieser Form kein Beleg ist. Der billigste Fix ist, das Profil zu nennen.
- **verifizierbar:** ja, sobald der Korpus benannt ist.
- **klasse:** `messung-ohne-nennbaren-korpus`

### L-1 · Der Identity-Kommentar wiederholt die ADR-Begründung statt die Kopplung

- **quelle:** Baseline `grundlagen-harness-dateien.md` §Was ein Kommentar trägt, Leser-Modell (*„nicht an jemanden, der die Entscheidung noch einmal treffen will. Der zweite Leser hat die ADR."*)
- **pfad:** `internal/hexagon/core/model/config.go:589-593`
- **gemessen:** Satz 1 (*„exempt-expect-count geht BEWUSST NICHT ein (ADR-0078)"*) ist eine
  saubere Abgrenzung samt auflösbarem Herkunfts-Feld. Die restlichen zwei Sätze
  (*„… sind kein Paar verschiedener Zusagen, sondern ein Widerspruch — eine davon muss
  falsch sein. Sie als Duplikat abzuweisen ist die richtige Antwort, nicht sie zu
  trennen."*) sind wortgleich die Entscheidung 8 der ADR `:130-136`; der Schlusssatz
  bewertet zusätzlich die verworfene Option. Adressat ist der Entscheider, nicht der
  Ändernde.
- **verifizierbar:** nein.
- **klasse:** `kommentar-dupliziert-adr-begruendung`

### L-2 · „wortgleich" ist für eines der beiden zertifizierten Zitate zu stark

- **quelle:** BEO-012, `AGENTS.md` §5 (Zitat-Disziplin)
- **pfad:** `docs/plan/cr/2026-08-30-antwort-a-check-leermenge.md:32-33`
- **zugesagt:** *„Die beiden Spec-Zitate der zweiten Fassung sind nachgelesen — **wortgleich und im Geltungsbereich**."*
- **gemessen:** das zweite Zitat (*„auf stdout unabhängig vom Code"*) ist wortgleich
  (`spec/lastenheft.md:346-347`). Das erste nicht ganz: der CR zitiert
  *„… und eine Diagnose, die ‚0 Befunde' ausweist."*, das Boundary-Kriterium
  (`spec/lastenheft.md:360`) endet *„… ausweist, **ohne Fix-Kandidaten**."* — die
  Fortsetzung ist abgeschnitten und durch einen Punkt ersetzt, dazu sind die inneren
  Anführungszeichen getauscht und *ohne Befunde* ist gefettet. Am Geltungsbereich ändert
  das nichts, an *wortgleich* schon.
- **warum es zählt:** wir stellen dem Absender ein Prüf-Zeugnis aus; ein Zeugnis über
  Wortgleichheit sollte die Wortgleichheit tragen.
- **verifizierbar:** ja — Zeichenvergleich gegen `spec/lastenheft.md:360`.
- **klasse:** `zitat-zertifikat-zu-stark`

### L-3 · Das Handbuch nennt die zwei Exit-2-Ränder des neuen Schlüssels nicht

- **quelle:** Reviewer-Skill LOW *Doku-Drift*
- **pfad:** `docs/user/benutzerhandbuch.md:2016-2018` (Beispiel-Kommentar) und `:2138-2152` (Prosa)
- **gemessen:** beide Stellen erklären Bedeutung, Beidseitigkeit, die Null und den Preis —
  aber nicht, dass der Schlüssel **ohne** `exempt-section-pattern` und mit einem Wert
  `< 0` zu Exit 2 führt. Das `--print-config`-Gerüst nennt beide
  (`internal/adapter/driving/cli/config_template.go:210`), die Spezifikation ebenfalls
  (`:2902`). Wer das Handbuch als Einstieg nimmt — genau die Zielgruppe, die der
  Handbuch-Standard nennt —, läuft in einen Exit 2 ohne Vorwarnung.
- **verifizierbar:** nein.
- **klasse:** `handbuch-config-rand-fehlt`

### I-1 · Release-Prep-Schuld: beide READMEs sagen die Nullmengen-Härte unbedingt

- **pfad:** `README.md:141-146` · `README.de.md:143-149`
- **gemessen:** *„Both only **shrink** the checked set, carry no reason code of their own,
  and without them the finding set is byte-identical; **if the section exemption empties
  the set, `section-missing` is reported** rather than silent green."* Das gilt seit
  diesem Slice nur noch ohne `exempt-expect-count`, und der neue Grund-Code fehlt in der
  Aufzählung. `CHANGELOG.md` ist ebenfalls unberührt.
- **kein Befund gegen diesen Slice:** beide Dateien sind nach der gelebten
  Commit-Grenze Release-Prep-Fläche (`git log -- CHANGELOG.md` zeigt ausschließlich
  `docs(release)`-Commits). Hier notiert, damit die Zeile beim nächsten Release-Prep nicht
  durchrutscht — die Modul-Listen der READMEs sind genau die Fläche, die kein Gate
  erreicht.

### I-2 · Die CR-Antwort protokolliert den zurückgezogenen Erstentwurf des Absenders zweimal

- **pfad:** `docs/plan/cr/2026-08-30-antwort-a-check-leermenge.md:34-35`, `:104-106`
- **gemessen:** *„dass ihr euren ersten Entwurf selbst korrigiert habt, statt ihn zu
  verteidigen, hat diese Runde kurz gemacht"* und *„Ihr habt das selbst zurückgezogen; wir
  bestätigen es"*. Beides ist im Ton fair bis freundlich, und die zweite Stelle trägt eine
  echte Bestätigung. Gemessen an der Nutzer-Vorgabe (*CR-Dokumente tragen Bitte und
  Beleg*) ist eine zurückgezogene Fassung trotzdem Entstehungsgeschichte; der Ort dafür
  ist die Closure-Notiz. Kein Blocker — die Stelle steht hier, weil sie beim nächsten
  ausgehenden Dokument ohne Ton-Puffer dieselbe Klasse wäre.

---

## Negativbefunde (geprüft, ohne Befund)

- **Ton und Fairness der CR-Antwort:** angemessen. Der Antrag wird in der Sache
  ausdrücklich als tragend anerkannt, die Zerlegung des Absenders wird als *„die richtige
  Zerlegung"* übernommen, jede der drei Abweichungen wird benannt statt untergeschoben,
  und die Grenze der eigenen Messung steht als Re-Evaluierungs-Trigger daneben. Es gibt
  keine d-check-interne Forensik (kein Slice-Verlauf, keine Review-Historie, keine
  Deliberation über den eigenen toten Code). Der einzige Einwand ist H-1: die Messung, mit
  der wir absagen, ist falsch — das ist ein Sach-, kein Ton-Befund.
- **`AGENTS.md` §3.1 (Docker/make-only):** kein Host-Werkzeug im Diff; alle neuen Läufe
  gehen über `make`-Targets bzw. das Image.
- **§3.2 (Suppression-Verbot):** keine `nolint`-Direktive im Diff; `make lint` → `0 issues`.
- **§3.4 (Spec-Straten nie abwärts):** die neuen Lastenheft-/Spezifikations-Absätze nennen
  *„Begründung in begleitender ADR"* ohne ADR-Kennung, keinen Slice, keine Welle, keinen
  Commit-Hash. Das `matrix`-Modul bestätigt es mechanisch (`make doc-check` → 0 Befunde).
- **§3.5 (ADR-Immutabilität):** ADR-0078 ist neu; `make adr-check RANGE=a062fe8^..HEAD`
  → 0 Befunde. Der ADR-Index ist um genau eine Zeile ergänzt
  (`docs/plan/adr/README.md:88`), Zell-Längen halten die eigenen `structure`-Regeln.
- **§3.6 (Gate-Lockerung nur mit ADR):** die Lockerung ist ADR-0078, sie ist als solche
  benannt (*„Die geprüfte Menge zu verkleinern bleibt eine Lockerung"*) und steht mit
  gleichem Wortlaut in §5 des Slice-Plans, in der ADR und in der CR-Antwort.
- **§3.8 (Modul verspricht nur über die Scan-Menge):** der neue Schlüssel liest keine
  Eingabe außerhalb der Scan-Menge — er zählt Überschriften, die das Modul ohnehin
  gefunden hat. Keine neue Ziel-Achse. (Die Skopus-Frage, die tatsächlich offen ist, ist
  Datei vs. Regel — M-1, andere Achse.)
- **Hexagon-Richtung (ADR-0005/ADR-0012):** `model` bekommt ein Feld, `rules` liest es,
  `configyaml` befüllt es, `app` trägt den Klartext — keine neue Kante.
  `make arch-check` grün innerhalb von `make gates`.
- **Netz (`DC-QA-03`):** kein Zugriff außerhalb `external`; alle neuen Läufe hermetisch.
- **Vollständigkeit der Deklarations-Oberflächen:** `--print-config`-Gerüst ✔,
  `AllReasons()` ✔, `--doctor`-Klartext (`reasonTexts`) ✔, Spezifikation §Schritt 3 ✔,
  §2-Schema ✔, §4 `SPEC-077` ✔ samt Historie, Lastenheft-Tabelle ✔ (Platzierung siehe
  M-3), Lastenheft-Akzeptanzkriterien ✔ (zwei neue), Lastenheft-Historie ✔, Handbuch ✔.
  `FixCandidateFor` braucht nichts (nur eindeutig ableitbare Codes). Keine fehlende
  Oberfläche gefunden.
- **Zitate im Geltungsbereich (BEO-012), außer L-2:** ADR-0075 wird korrekt für die
  Nullmengen-Härte zitiert (Wortlaut gegen `docs/plan/adr/0075-…:113-114` geprüft, bis auf
  entfernte Fett-Marker identisch); `DC-FA-CLI-003` trägt die Aussage
  *„Differenzierte Exit-Codes pro Befund-Kategorie"* als Out-of-Scope tatsächlich
  (`spec/lastenheft.md:146`); ADR-0073 trägt die `hint`-Ausnahme für die Klasse
  *„die Regel hat nicht gemessen"*; ADR-0070 hat die Verdopplung der Tabellen-Bedingungen
  tatsächlich zurückgebaut. Der Verweis auf `DC-FA-STRUCT-001` als Vorschrift für einen
  eigenen Grund-Code trägt: `spec/lastenheft.md:2510-2515` schreibt genau das vor
  (*„jede mit eigenem Grund-Code, weil jede eine andere Reparatur verlangt"*).
- **„Neue Bauform ohne Präzedenz":** nachgezählt gegen das §2-Schema — `min-sentences`,
  `max-tasks`, `cell-min-chars`, `cell-max-chars` sind Schwellen, keine erwartete Anzahl.
  Die Aussage hält.
- **Die vier Ergebnis-Zeilen der CR-Antwort-Tabelle (19/19/19 · 20/19/19 · 19/18/19 ·
  Selektor leer):** alle vier gegen den Code und die Tests geprüft, alle vier korrekt.
- **Slice-Kopf und Vorprüfungen:** `**Lifecycle:**` statt `**Status:**` ✔,
  `**Verantwortlich:**` bei der Beanspruchung gesetzt ✔, `**Berührte Spec-Stellen:**` mit
  Kennung ✔, drei Vorprüfungen vorhanden ✔, die beiden kanonischen mit `cite`-Direktive
  und die Nachtlauf-Vorprüfung bewusst ohne — konform zu MR-053/MR-054;
  `make doc-check` (Modul `citations`) bestätigt die Anker.
- **MR-013 (Lifecycle-Move):** `a062fe8` ist ein reiner `git mv` plus der gekoppelte
  Roadmap-Flip (zwei gelöschte Zeilen), die Slice-Datei selbst unverändert —
  `make planning-check` grün.
- **Handbuch verkauft keinen Preis als Vorteil:** der neue Absatz endet ausdrücklich mit
  *„**Der Preis gehört dazu:** die Zahl altert wie jeder andere Autoren-Text, und wer sie
  blind mitzieht, hat einen Wächter, der nur noch sich selbst bestätigt."* Das ist die
  Korrektur des slice-181-Befunds, und sie sitzt. Das Beispiel steht in einem
  auskommentierten Block und ist deshalb kein Validator-Kandidat; der
  `TestDocExamples_ConfigBeispieleValidieren` (Scope: Handbuch, `operations.md`, beide
  READMEs) bleibt grün.
- **`exempt-paths` bekommt nichts:** die Abgrenzung ist in ADR-0078, im Slice-Plan §3 und
  in der CR-Antwort benannt und übernimmt die Begründung des Absenders — kein
  unbenannter Konsistenz-Spalt zwischen zwei Ventilen derselben Klasse.
- **Arbeitsbaum:** nach allen Proben `git status --short` leer bis auf diesen Report.

---

## Kategorie-Summary

| Kategorie | Anzahl | Klassen |
|---|---|---|
| HIGH | 4 | `messung-widerlegt-durch-eigenen-vertrag` · `tabellenkopf-ueberschrieben` · `zusage-ohne-fangende-probe` · `kommentar-konjunktiv-rueckwaerts` |
| MEDIUM | 4 | `zusage-skopus-datei-vs-regel` · `frueher-abbruch-undeklariert` · `spec-zeile-in-falscher-klasse` · `messung-ohne-nennbaren-korpus` |
| LOW | 3 | `kommentar-dupliziert-adr-begruendung` · `zitat-zertifikat-zu-stark` · `handbuch-config-rand-fehlt` |
| INFO | 2 | Release-Prep-Schuld READMEs/CHANGELOG · zurückgezogener Erstentwurf in der CR-Antwort |

**Wiederkehrende Klassen dieser Sitzung:** BEO-020 (Behauptung sieht aus wie Beleg) trifft
H-1 und M-4 — zweimal in einem Slice. BEO-023 (Wächter wird still schwächer) trifft H-3
mit zwei unabhängigen Mutationen. BEO-012 trifft nur L-2, und zwar in der milden Form —
das ist gegenüber slice-181 eine deutliche Verbesserung.
