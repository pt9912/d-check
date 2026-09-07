# Review-Report — slice-212, Runde 1 (Grenzen-Inventur über die 24 Sensor-Dateien)

**Review-Art:** Code/Design — geprüft werden die **vier ergänzten Grenzen** gegen das Skript bzw. Target, das sie beschreiben, die **Inventur-Behauptung „20 von 24"** gegen eine eigene Stichprobe, und der **neue Registereintrag** gegen seine beiden Nachbarn und seine drei Belege. **Nicht** geprüft: die DoD-Abhakung und die Gate-Lauf-Bestätigung (Verifier-Rolle).
**Gegenstand:** `c5fc6d63` (Vorprüfungen) · `af0c3fdd` (Beanspruchung) · `1f398f02` (DoD 1+2) · `03ca7b1c` (DoD 3). Diff-Range `c5fc6d63^..HEAD`, `HEAD == 03ca7b1c` zu Beginn und am Ende des Reviews.
**Skill:** `.harness/skills/reviewer.md` v1.15.0 @ `03ca7b1c` — einschließlich Anker 17 (Messung zählt einen Proxy statt des Gegenstands).
**Modell-ID:** claude-opus-5[1m] · **Datum:** 2026-09-07
**Eingangs-Kontext:** der Slice-Plan slice-212 (§1, §2, §3, §6, §8, DoD 1–3); der Folge-Slice slice-213 in `open/`; der neue Registereintrag `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen` samt `observation.md`, `state.md` und drei Evidence-Dateien; die beiden Nachbar-Einträge `BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt` und `BEO-ALL/citation-stretched-beyond-scope`; die Review-Reports von slice-209 und slice-210; `AGENTS.md` §3.1/§3.2/§3.3/§3.5/§3.6/§3.7/§3.8/§4/§5/§6; `harness/README.md` §Sensors; `harness/conventions.md` samt MR-011, MR-013, MR-045, MR-053, MR-054, MR-055, MR-069, MR-070; `DC-FA-RVW-001`, `DC-FA-VCS-001`, `DC-FA-CLI-011`, `DC-QA-03`; Baseline `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice und §Zwei Schritte vor der Modus-Begründung, `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register, `regelwerk/modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin).

**Eigene Läufe (echte Ausgaben, keine behaupteten):**

- `make gates` — `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`, `d-check: 724 Datei(en) geprüft, 0 Befund(e)`, `coverage-gate: OK — Coverage 94.60% erfüllt Schwelle 93%`, `fetch-baseline-cache: verify ok (54 Dateien, vollständig)`, Exit 0.
- `make doc-check` — `d-check: 724 Datei(en) geprüft, 0 Befund(e)`, Exit 0.
- **Eigener Bruch-Test zu DoD 1** in einer isolierten Kopie unter dem Scratchpad (eigenes `git init`, damit `git rev-parse --show-toplevel` im Skript greift): drei Läufe von `fetch-baseline-cache.sh --verify` und einer von `--check-latest`. Der Arbeitsbaum blieb unberührt (`git status --porcelain` leer vor und nach dem Test).
- **Code-Lektüre statt Vermutung** für die drei übrigen Ergänzungen: `internal/hexagon/core/rules/reviews.go`, `internal/hexagon/core/rules/vcs.go`, `internal/hexagon/core/rules/matrix.go` (`excludedRanges`), `tools/semgrep.sh`, `.golangci.yml`, `.a-check.yml`, `.d-check.yml` sowie die `edges`-Semantik im Schwester-Repo a-check (`Edges []Edge // allowed directed edges`, Grund-Code `wrong-direction`).
- **Eigene Stichprobe zur Inventur** über zwanzig als *vollständig* geführte Dateien, je gegen ihr Skript, ihr Make-Target und ihren Konfigurations-Block.
- **Zählungen** über das Beobachtungs-Register, `docs/reviews/`, `docs/plan/planning/**` und `docs/plan/adr/` mit `ls`/`find`/`grep`/`awk`.
- Kein Host-Go, kein Host-Skript-Interpreter, keine Schreiboperation im Repo außer diesem Report.

**Zitier-Form.** Dieser Report friert ein: Baseline-Stellen als Tag plus Pfad in Inline-Code, Slices als Kennung, Gates als `make <target>`-Token.

---

## Was gemessen reproduziert (vorweg, damit die Findings nicht als Gesamturteil gelesen werden)

| Behauptung | Quelle | Nachmessung |
|---|---|---|
| Änderung ohne nachgezogenes Manifest ⇒ `GESCHEITERT`, Exit 1 | `harness/sensors/baseline-verify.md` §Grenze (0) | **reproduziert exakt** — `regelwerk/modul-05-planning-harness.md: GESCHEITERT`, Exit 1 |
| Dieselbe Änderung **mit** nachgezogenem Manifest ⇒ `verify ok (54 Dateien, vollständig)`, Exit 0 | ebd. | **reproduziert exakt**, Exit 0 — die manipulierte Zeile stand dabei nachweislich in der Datei |
| `--check-latest` fängt es, Exit 4, `UPSTREAM-CONTENT-DRIFT` | ebd. | **reproduziert exakt** — `check-latest OK (Currency)` plus `UPSTREAM-CONTENT-DRIFT`, Exit 4 |
| Das semgrep-Regelset ist gepinnt und lokal gecacht, altert also | `harness/sensors/semgrep.md` §Grenze (2) | **trägt** — `RULES_COMMIT` ist ein fester Commit-Pin, der Cache liegt außerhalb des Repos, der Scan läuft `--network none` |
| Der Verweis `make freshness-semgrep` zeigt auf die Datei einer anderen Achse | ebd. | **kein Befund** — jene Datei führt im Titel *„u. a."* und im Vertrag ausdrücklich alle vier Achsen inklusive `SEMGREP_VERSION`; `harness/README.md` verlinkt die Familie ebenso |
| Die erste Fehlmessung: Listenmarker gezählt, `mention-coverage` als null gemeldet | slice-212 §3 | **trägt** — jene Datei führt ihre Grenzen als Fettabsätze, kein einziger Listenmarker |
| Register: 37 Verzeichnisse zum Lesezeitpunkt | slice-212 §8 | **trägt** — heute 38, das 38. ist der von diesem Slice angelegte |
| `citation-stretched-beyond-scope` 18×, `wortlaut-behauptet-pruefung-die-fehlt` 8×, `eigene-menge-gemessen-fremde-behauptet` 12×, `zaehlmethode-misst-proxy-statt-gegenstand` 4× | slice-212 §8 | **alle vier exakt** — Evidence-Dateien nachgezählt |
| Die beiden nachgetragenen Belege geben her, was ihnen zugeschrieben wird | `evidence/slice-209.md`, `evidence/slice-210.md` | **beide belegt** — F-2 des slice-209-Reports nennt die sechste Grenze und *„`MR-070` fünf andere Grenzen ausschreibt"*; F-4 des slice-210-Reports trägt Titel und Begründung wortgleich zur Zuschreibung |
| Die Abgrenzung gegen die beiden Nachbarn | `observation.md` | **trägt** — der eine Nachbar handelt vom Wächter, der zu fangen behauptet, der andere von der überdehnten Quelle; keiner deckt „Geltungsbereich richtig gelesen, Kehrseite nicht aufgeschrieben" |
| §1-Abgrenzung *„keine Verhaltens-Änderung"* | slice-212 §1 | **eingehalten** — der Diff fasst ausschließlich `.md`-Dateien an |
| §1-Abgrenzung *„die Grenzen-Abschnitte nicht inhaltlich neu schreiben"* | ebd. | **eingehalten** — alle vier Eingriffe sind Nachträge; keine bestehende Nummer wurde umgeschrieben oder verschoben (`baseline-verify` fügt bewusst eine `0.` vor der bestehenden `1.` ein, statt 1–4 zu renumerieren) |
| Die beiden `d-check:cite`-Anker in §8 | slice-212 §8 | **wortgleich und vorschreibend** — Zeilen 268–269 und 274 der Kanon-Datei sind die beiden nummerierten Vorschriften aus §Zwei Schritte vor der Modus-Begründung, nicht eine Erläuterung daneben |

---

## Findings

### F-1 · MEDIUM · Die ergänzte Grenze in `review-coverage` beschreibt einen Mechanismus, den das Modul nicht hat — und misst ihn dann am Bestand

- **kategorie:** MEDIUM (Anker 9 *Quelle über ihren Geltungsbereich hinaus zitiert*, hier auf einen **Code-Kommentar** angewandt, plus Anker 17 *Messung zählt einen Proxy*. Kontext-Eskalation erwogen und **verworfen**: kein Code-Pfad ändert sich, kein Gate wird schwächer — der Schaden liegt im Modell des Lesers)
- **quelle:** `DC-FA-RVW-001` · ADR-0081 · `AGENTS.md` §5 (*„Eine zitierte Quelle trägt nur, was in ihrem Geltungsbereich steht"*)
- **pfad:** `harness/sensors/review-coverage.md:24-29` (§Grenze, neuer Eintrag 5)
- **befund:** Der Eintrag sagt: *„Wäre eine Kennung Präfix einer anderen, deckte der Report der längeren die kürzere mit ab."* Das kann nicht eintreten. `hasMatchingReview` in `internal/hexagon/core/rules/reviews.go:163` vergleicht nicht auf Enthaltensein, sondern **extrahiert und vergleicht auf Gleichheit**: `if m := sliceIDRE.FindString(e.Name); m == id`. Für den Slice `slice-21` liefert ein Report `…-slice-212-…` den Treffer `"slice-212"`, und `"slice-212" != "slice-21"` — kein Match. Dasselbe im Schwester-Werkzeug (`tools/archive-wave/collect.go:197`, `FindStringSubmatch` plus Set-Lookup). Der Satz übernimmt das Wort *Substring-Match* aus dem Code-Kommentar und aus dem eigenen Vertrags-Teil derselben Datei (Zeile 9) und liest es als Mechanismus, statt die Stelle aufzuschlagen, die es benennt. Der angehängte Beleg *„Im Bestand nicht eingetreten (alle Kennungen dreistellig, keine ist Präfix einer anderen)"* misst danach ein Risiko, das es nicht gibt — die Kennungs-Länge stimmt zwar (344 Slice-Dateien, alle dreistellig; 16 Kennungen in `docs/reviews/`, alle dreistellig), sagt aber über diesen Gate nichts. **Versagensszenario:** Ein Leser leitet aus der Zeile den Handlungsbedarf für vierstellige Kennungen ab und baut einen Wächter gegen eine Nicht-Lücke; gleichzeitig bleibt die **tatsächliche** Grenze derselben Zeile unbenannt — `FindString` nimmt den **ersten** Treffer im Report-Dateinamen, ein Report, der zwei Slices im Namen führt (`…-slice-209-und-slice-210-…`), deckt also nur den erstgenannten, und der zweite bleibt `review-missing`, obwohl sein Report existiert.
- **verifizierbar:** ja — `internal/hexagon/core/rules/reviews.go:163` gelesen; reproduzierbar durch einen Report mit zwei Kennungen im Namen und `make review-coverage`. Kein Gate hält die Aussage einer Sensor-Datei über sich selbst.
- **klasse:** `grenze-beschreibt-mechanismus-der-nicht-existiert`

### F-2 · MEDIUM · Die ergänzte Grenze in `arch-check` sagt das Gegenteil der `edges`-Semantik: eine undeklarierte Kante ist ein Befund, nicht unsichtbar

- **kategorie:** MEDIUM (Anker 10 *Messmethode klafft gegen die Stelle, die sie erfüllen soll*; hier: Beschreibung gegen Werkzeug-Vertrag)
- **quelle:** ADR-0005 · ADR-0012 · ADR-0029
- **pfad:** `harness/sensors/arch-check.md:19-23` (§Grenze, neuer Eintrag 2)
- **befund:** Der Eintrag sagt: *„eine Kante, für die keine Regel existiert, ist kein Befund, sondern unsichtbar."* Im Werkzeug ist `edges` eine **Erlaubnisliste**, keine Verbotsliste: das Feld heißt dort `Edges []Edge // allowed directed edges`, und der Grund-Code `wrong-direction` ist definiert als *„ein Import quert eine Schicht-Kante entgegen `edges`/`allow`"* — mit dem Akzeptanzkriterium *„Given eine Kante gegen die deklarierte Richtung … then ein Befund (`wrong-direction`) und Exit-Code 1"*. Genau darauf beruht auch die Absicht der eigenen Konfiguration: `.a-check.yml` listet neun Kanten und **keine** aus `model` heraus, weil `model` keine Ausgänge haben soll (R6/ADR-0012). Wäre der Satz wahr, wäre diese Liste Dokumentation statt Durchsetzung. **Versagensszenario:** Jemand liest die Zeile, hält den `edges`-Block für unverbindlich und streicht einen Eintrag als „ohnehin nicht geprüft" — womit er den Gate **verschärft** und einen bislang erlaubten Import rot macht; oder er fügt eine Kante ohne Entscheid hinzu, weil sie „nichts erzwingt" (`AGENTS.md` §3.6-nah). Umgekehrt bleibt das, was wirklich unsichtbar ist, unbenannt: `exclude: ["**/*_test.go", "tools/archive-wave/**"]` — der zweite Eintrag, ein vollständiger Werkzeug-Baum, steht in dieser Datei weder im Vertrag noch in der Grenze (siehe F-3).
- **verifizierbar:** ja — Modell-Kommentar und Grund-Code-Tabelle im a-check-Repo gelesen; reproduzierbar durch einen Import `internal/hexagon/core/model` → `internal/hexagon/core/rules` und `make arch-check` (erwartet: `wrong-direction`, nicht Schweigen).
- **klasse:** `grenze-kehrt-werkzeug-semantik-um`

### F-3 · MEDIUM · „20 von 24 vollständig" hält der Stichprobe nicht stand — drei als *vollständig* geführte Dateien lassen eine Lücke aus, die ihr eigener Vertrag oder ihre Konfiguration bereits benennt

- **kategorie:** MEDIUM (Anker 8 *Schluss reicht weiter als die gemessene Menge* — die Zahl ist die tragende Aussage von DoD (2) und wird in der Commit-Botschaft wiederholt)
- **quelle:** `AGENTS.md` §5 · Baseline `v6.5.0` · `regelwerk/modul-13-quality-gates.md` §Hard Rule (*„Ein Gate ohne seine Grenze behauptet ebenfalls zu viel"*)
- **pfad:** slice-212 §3 (Inventur-Tabelle) · `harness/sensors/lint.md:12-20` · `harness/sensors/doc-check.md:16-33` · `harness/sensors/arch-check.md:14-23`
- **befund:** Der Slice definiert *vollständig* in §3 selbst: *„jede Lücke, die das Skript oder das Target selbst kennt — als Kommentar, als Exit-Code, als benannter Ausgang —, steht auch im Abschnitt."* Nach dieser Form fallen mindestens drei der zwanzig durch:
  **(a) `lint`.** Der Vertrag der Datei nennt die Tatsache (*„Ausnahmen leben zentral in `.golangci.yml` … mit Begründung"*), die Grenze nennt ihre Kehrseite nicht. `.golangci.yml` schaltet unter `exclusions` für jede `_test.go`-Datei fünf Komplexitäts-Linter (`cyclop`, `gocognit`, `gocyclo`, `nestif`, `funlen`) sowie `noctx`, `unparam` und zwei `revive`-Regeln ab, dazu `testpackage` für den ganzen Kern-Baum — jede Ausnahme mit eigenem `Why:`-Kommentar, also benannt. Ein grüner `make lint` sagt über Komplexität in Testcode nichts; die Grenze behauptet nur, `nolintlint` prüfe Form statt Berechtigung.
  **(b) `doc-check`.** Die Grenze nennt `d-check:ignore` ausschließlich in der Merge-Richtung (*„greift dort nicht"*) und `scan.ignore`/`citations.scope` als die groben Ventile. Ungenannt bleibt, was das Grün am stärksten formt: 45 `ignore-refs`-Paare in `.d-check.yml` und 246 `d-check:ignore`-Marker im Baum. Für das erste Ventil steht die Bewertung seit demselben Tag im eigenen Konventionsspeicher — MR-069 führt es ausdrücklich als *„deklarierte Gate-Senkung"*, misst 25 Baseline-Einträge über zehn entfernte Tags und benennt, dass davon sieben nichts mehr skopieren. Eine Grenzen-Liste, die eine als Gate-Senkung deklarierte Konfiguration nicht führt, ist die Klasse, gegen die dieser Slice geschrieben ist.
  **(c) `arch-check`.** `.a-check.yml` nimmt neben Testdateien `tools/archive-wave/**` vollständig aus, mit Begründung im Kommentar. Der Vertrag der Sensor-Datei nennt nur die Testdateien, die Grenze keines von beiden.
  **Versagensszenario:** Der Slice führt genau diese drei Dateien als Beleg dafür, dass der Bestand *„besser als erwartet"* ist, und slice-213 baut darauf auf (§1: *„Eine erneute Inventur … liegt in slice-212 und ist dort abgeschlossen"*). Wer die Zahl übernimmt, hält eine Menge für geprüft, die es an drei Stellen nicht ist.
- **verifizierbar:** ja — `.golangci.yml` §`exclusions`, `grep -c "^  - in:" .d-check.yml` (45), `grep -rn "d-check:ignore" --include=*.md .` ohne die vendorte Baseline (246), `.a-check.yml` §`exclude`. Kein Gate hält die Vollständigkeit einer Grenzen-Liste; das sagt der Slice selbst.
- **klasse:** `inventur-behauptet-vollstaendigkeit-die-die-stichprobe-bricht`

### F-4 · MEDIUM · Die dritte Antwort zu `adr-check` ist am Code entscheidbar — der Bruch-Test war nicht nötig, und die größte Grenze der Datei bleibt deshalb unbenannt

- **kategorie:** MEDIUM (Anker 10; zusätzlich Anker 15 — die Frage betrifft eine Eingabe, die das Modul liest, ohne sie zu scannen: zwei git-Stände)
- **quelle:** `DC-FA-VCS-001` · ADR-0024 · `AGENTS.md` §3.5
- **pfad:** slice-212 §3 (Inventur-Tabelle, Zeile `adr-check`, und der Absatz *„`adr-check` ist die dritte Antwort"*) · `harness/sensors/adr-check.md:16-19`
- **befund:** Der Slice erklärt die Frage — *„ob eine Kern-Änderung, die innerhalb eines `## Geschichte`-Abschnitts abgelegt wird, den Vergleich passiert"* — für *„nur durch einen Bruch-Test zu beantworten"* und begründet den Verzicht damit, der Test bräuchte eine manipulierte `Accepted`-ADR im Arbeitsbaum. Die Antwort steht in drei lesbaren Stellen, die der Slice für die übrigen 23 Dateien auch gelesen hat: `.d-check.yml:721` setzt `exclude-sections: [Geschichte]`; `vcsCore` (`internal/hexagon/core/rules/vcs.go:134-147`) verwirft vor dem Hashen jede Zeile aus `excludedRanges`; und `excludedRanges` (`internal/hexagon/core/rules/matrix.go:294-318`) schneidet vom passenden Heading bis zum nächsten Heading gleicher oder höherer Ebene — **existiert keines, bleibt `r.to == 0`, und `inRanges` liest das als „bis Dateiende"**. Im Bestand ist genau das der Regelfall: **79 von 79** ADRs mit einem `## Geschichte`-Abschnitt führen ihn als **letzte** Sektion. Das Grün von `make adr-check` sagt damit über alles, was unterhalb dieser Überschrift steht, nichts — beliebig viel, beliebigen Inhalts, bis zum Dateiende. Die Tatsache steht bereits im Vertrag der Sensor-Datei (*„Erlaubt bleiben zwei Dinge: `## Geschichte`-Anhänge und der `**Status:**`-Übergang"*); ihre Kehrseite — die Reichweite dieses Ausschlusses — fehlt in der Grenze. **Versagensszenario:** Eine `Accepted`-ADR bekommt unter `## Geschichte` eine neu formulierte Entscheidung angehängt, die der Kopf-Entscheidung widerspricht; `make adr-check` bleibt grün, und die Sensor-Datei hat den Leser nicht gewarnt. Der Verzicht ist damit nicht *„keine Ausrede"*, sondern eine vermeidbare dritte Antwort: die Methode des Slice (beide Seiten lesen) reichte aus.
- **verifizierbar:** ja — die drei Code-/Config-Stellen sind gelesen; `for f in docs/plan/adr/[0-9]*.md` mit `grep -n "^## "` zeigt 79 von 79. Reproduzierbar durch einen Anhang unter `## Geschichte` einer `Accepted`-ADR und `make adr-check STAGED=1` (erwartet: grün).
- **klasse:** `dritte-antwort-statt-lesbarer-quelle`

### F-5 · LOW · Zwei verschiedene Zähler-Stände für denselben Registereintrag im selben Plan — und der falsche steht auch in der Commit-Botschaft

- **kategorie:** LOW (Anker 17; die Zahl trägt hier keine Schlussfolgerung, beide Werte liegen weit über der Schwelle)
- **quelle:** `AGENTS.md` §5 · Baseline `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register (*„Der Zähler wird abgeleitet, nicht geführt"*)
- **pfad:** slice-212:233 (§8, *„10×, zuletzt slice-210"*) gegen slice-212:16 (§Bezug, *„9×"*) und slice-212:102 (§3, *„steht bei 9×"*) · Commit `c5fc6d63` (Botschaft, *„(10x)"*)
- **befund:** `BEO-ALL/rule-drawn-from-occasion-not-inventory` trägt **neun** Evidence-Dateien, die jüngste `slice-210.md`. Der Plan nennt den Stand dreimal und zweimal richtig; §8 und die Botschaft des Vorprüfungs-Commits führen 10×. **Versagensszenario:** Der nächste Lauf sichtet das Register, findet 9, hält den Plan für den frischeren Stand und schreibt einen zehnten Beleg, den es nicht gibt — oder umgekehrt, er glaubt eine Instanz sei verlorengegangen. Dieselbe Klasse hat der slice-210-Report als F-5 bereits gemeldet (*„Zwei verschiedene Zähler-Stände für denselben Eintrag im selben Plan"*); dies ist ihre Wiederholung zwei Slices später.
- **verifizierbar:** ja — `ls docs/plan/planning/observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/evidence/ | wc -l` ⇒ 9. Kein Gate zählt Belege gegen Prosa.
- **klasse:** `zwei-zaehler-staende-fuer-denselben-eintrag`

### F-6 · LOW · Die neue `state.md` speichert einen abgeleiteten Zähler und wiederholt die Provenance, die schon in jeder Evidence-Datei steht

- **kategorie:** LOW (§3.7-nah; das Feld nennt Zustand *und* auflösbaren Beleg, deshalb keine HIGH-Chronik — der Zusatz ist eine latente Wartungsfalle)
- **quelle:** `AGENTS.md` §3.7 · Baseline `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register (*„Ein gespeicherter Zähler neben einer Belegliste sind zwei Quellen für denselben Zustand"*)
- **pfad:** `docs/plan/planning/observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md:1`
- **befund:** Der zweite Satz lautet *„Die Schwelle ist mit dem dritten Beleg erreicht; zwei der drei sind nachgetragen und als solche gekennzeichnet."* Beide Hälften sind aus `evidence/` ableitbar, und die zweite steht wörtlich in den beiden betroffenen Dateien selbst. **Versagensszenario:** Sobald ein vierter Beleg landet — und slice-213 §6 nennt genau das als besseren Beleg —, sagt die `state.md` „mit dem dritten" und „zwei der drei", während `evidence/` vier führt; das Feld behauptet dann einen Stand, den es nicht besitzt, und niemand meldet es. Die beiden Vorgänger-Reports haben je eine Hälfte dieser Klasse gemeldet: slice-210 F-3 (gespeicherter Zähler neben der Belegliste) und slice-209 F-7 (`state.md` kopiert, statt zu zeigen).
- **verifizierbar:** ja — durch Anlegen einer vierten Evidence-Datei; der Text der `state.md` ändert sich nicht mit. Kein Gate prüft Zustandsfeld-Form.
- **klasse:** `state-md-speichert-abgeleiteten-stand`

### F-7 · LOW · Eine lebende Sensor-Datei verweist für ihren Beleg in ein Zeitdokument, das planmäßig archiviert wird

- **kategorie:** LOW (latente Wartungsfalle; der Verweis ist heute korrekt)
- **quelle:** `AGENTS.md` §3.3 (Ausnahme Slice-Lifecycle-Move) · MR-062 · Baseline `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Wann Arbeit eine Welle braucht (Zeitdokumente archivieren)
- **pfad:** `harness/sensors/baseline-verify.md:49-52`
- **befund:** Der Schlusssatz des Grenzen-Abschnitts lautet: *„die Echtheits-Grenze (0) hat keine [Probe] … ihr Bruch-Test steht in [slice-212]"*, mit einem Link auf den Lifecycle-Pfad in `in-progress/`. Das ist der **einzige** Link aus einer der 24 Sensor-Dateien in den Planungsbaum; die drei anderen Nennungen von `docs/plan/planning` dort sind Verzeichnis-Muster in Inline-Code, keine Verweise. Der Pfad wandert zweimal: beim `git mv` nach `done/` (dort zieht ihn die Kopplung aus MR-013 mit) und danach bei `make archive-wave SLICE=212 APPLY=1` nach `done/wellenlos/`. **Nach der Archivierung liegt dort ein Stub**; der Volltext mit der Messung liegt in `archiv.zip`. Der Link bleibt grün, weil der Stub existiert — die Zusage *„ihr Bruch-Test steht in slice-212"* wird still falsch. Hinzu kommt, dass am Zielort kein Test steht, sondern die Prosa-Aufzeichnung einer einmaligen Messung: `make baseline-probe`, `make guard-probe` und die `--selftest`-Modi sind in diesem Repo das, was *Probe* heißt, und die Zeile stellt sich mit dem Wort *Bruch-Test* daneben. **Versagensszenario:** Ein späterer Leser folgt dem Verweis, findet einen Stub ohne Messung und hat für die Grenze, die *„alle folgenden überwiegt"*, keinen Beleg mehr.
- **verifizierbar:** ja — `grep -rn "docs/plan/planning" harness/sensors/*.md` zeigt den Einzelfall; die Staleness ist reproduzierbar durch den Archiv-Lauf. `make doc-check` bleibt in beiden Zuständen grün.
- **klasse:** `lebende-datei-verweist-in-zeitdokument`

### I-1 · INFO · Der Ableiter des neuen Eintrags setzt voraus, dass der Vertrags-Teil richtig ist — die dritte seiner vier Instanzen widerlegt das an sich selbst

- **kategorie:** INFO (dokumentationswürdige, aber undokumentierte Annahme)
- **quelle:** Maintainability · Baseline `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register
- **pfad:** `docs/plan/planning/observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md` (Absatz *Ableiter*, Schritt 1)
- **befund:** Schritt 1 lautet *„den Vertrags-Teil desselben Artefakts durchgehen und jede Zusage einmal umdrehen"*. Bei `review-coverage` ist die umgedrehte Zusage (*„Substring-Match, 1:N zulässig"*) selbst ungenau, und das Umdrehen erzeugt deshalb eine falsche Grenze statt einer fehlenden (F-1). Der Eintrag benennt diese Vorbedingung nicht; sie ist der Unterschied zwischen „die Kehrseite aufschreiben" und „die Kehrseite **am Code** aufschreiben".
- **verifizierbar:** nein — Urteil; die Instanz dazu ist F-1.
- **klasse:** `ableiter-ohne-benannte-vorbedingung`

### I-2 · INFO · Der semgrep-Regel-Cache trägt keine Integritätsprüfung — dieselbe Frage, die dieser Slice für die Baseline beantwortet

- **kategorie:** INFO (bewusste Designnotiz, keine benannte Lücke des Skripts — deshalb **kein** DoD-2-Befund)
- **quelle:** ADR-0010 · `DC-QA-02`
- **pfad:** `tools/semgrep.sh` (Setup-Zweig, `RULES_DIR` unter `$XDG_CACHE_HOME`)
- **befund:** Der Regel-Cache wird einmalig am Commit-Pin geholt und danach nur noch über die Existenz von `$RULES_DIR/$RULES_SUBSET` erkannt; es gibt kein `SHA256SUMS`-Gegenstück wie beim vendorten Baseline-Baum. Der Schutz gegen stilles Grün ist die `Ran N rules`-Zeile mit N ≥ 1 — sie fängt ein leeres Subset, nicht ein verändertes. Damit gilt für das Regelset dieselbe Unterscheidung, die dieser Slice für `baseline-verify` aufgeschrieben hat (innere Konsistenz gegen Echtheit), nur ohne den zweiten Träger. Kein Versagen im Bestand gemessen; als Beobachtung notiert, weil die Symmetrie zum Slice-Gegenstand auffällt.
- **verifizierbar:** nein — kein Befund im Bestand; die Asymmetrie ist am Skript ablesbar.
- **klasse:** `pin-ohne-integritaets-gegenprobe`

---

## Negativbefunde (geprüft, ohne Befund)

- **DoD (1), Bruch-Test.** Beide Hälften der gemessenen Aussage reproduzieren exakt, einschließlich der Ausgabetexte und der drei Exit-Codes (1 / 0 / 4). Der Arbeitsbaum blieb unberührt.
- **DoD (1), Geltungsbereich.** Der Satz *„Integrität ist nicht Aktualität"* steht im Skript-Kopf und wird korrekt wiedergegeben; die Zuordnung zu `--check-latest` als **Netz, fail-open, kein Gate** stimmt mit `harness/README.md` §Sensors und der Skript-Dispatch-Zeile überein.
- **`semgrep`-Ergänzung.** Pin und lokaler Cache belegt; der Verweis auf die Frische-Achse zeigt auf die richtige Datei (siehe Tabelle oben). Kein Befund.
- **§1-Abgrenzung.** Vier Ausschlüsse, je mit Begründung, alle vier eingehalten: keine Verhaltens-Änderung (nur `.md` im Diff), keine Änderung an `baseline-verify` selbst, kein Sensor auf die Vollständigkeit, keine Neuformulierung bestehender Grenzen-Einträge.
- **Nummerierungs-Disziplin.** Die `0.` in `baseline-verify` vermeidet eine Renumerierung von 1–4 und hält damit alle bestehenden Bezüge auf „Grenze 1–4" gültig. Bewusst und richtig.
- **§8, drei Vorprüfungen.** Alle drei vorhanden und in der vorgeschriebenen Reihenfolge; die beiden kanonischen Blöcke tragen ihre `d-check:cite`-Direktive auf die **vorschreibenden** Zeilen (268–269 und 274), wortgleich — `citations` bestätigt das im grünen Lauf.
- **§8, Sub-Area-Wahl.** Eine Sub-Area (`*`), und die Aussage über `tools/harness/` (gelesen, nicht angefasst) trägt: kein Skript ist im Diff.
- **§8, Register-Sichtung.** Fünf einschlägige Einträge, vier Zähler exakt (der fünfte ist F-5); die Aussage *„Keiner der fünf erreicht mit diesem Slice die Schwelle erstmalig"* stimmt — der neue Eintrag ist keiner der fünf.
- **Anker-Paarung.** Der Ausgang ist *geplant*, also einer der drei kanonischen; er nennt eine Kennung, kein `liegt in`-Feld — richtig, denn *geplant* verlangt die Kennung des schreibenden Vorgangs, nicht den Zielort.
- **Folge-Slice-Paarung.** `docs/plan/planning/open/slice-213-grenzen-liste-braucht-fremden-leser.md` existiert im Lifecycle. Seine DoD (1)+(3) lösen den Ausgang ein: (1) beantwortet die Vorfrage aus dem **Wortlaut** von `AGENTS.md` §6, (3) setzt den Registereintrag auf einen Ausgang mit auflösbarem Zielort. Die Adresse nimmt die Sendung an — sie schließt den verwiesenen Punkt nicht aus und schließt nicht vor diesem Slice.
- **Register-Paarung.** Das neue Verzeichnis trägt `observation.md`, `state.md` und ein nicht leeres `evidence/`; die Belege sind formgebunden (der Dateiname **ist** die Vorgangs-Kennung), und alle drei Vorgänge existieren im Repo.
- **Die zwei nachgetragenen Belege.** Beide sind als nachgetragen gekennzeichnet, beide sind aus den zitierten Review-Reports belegbar (F-2 des slice-209-Reports, F-4 des slice-210-Reports), und keiner ist rekonstruiert. Der Einwand, dass rückwärts gesammelte Belege schwächer sind, steht als Risiko in slice-213 §6 — als Aussage über die **Gültigkeit** der Schwelle ist er offen und richtig offen gehalten.
- **Klassen-Abgrenzung.** Die Unterscheidung gegen `wortlaut-behauptet-pruefung-die-fehlt` und `citation-stretched-beyond-scope` ist nicht konstruiert: der erste Nachbar handelt vom Wächter, der zweite von der Quelle, der neue von der **Beschreibung** eines ehrlichen Wächters. Auch die übrigen 36 Verzeichnisse führen keine Zeile, die den Gegenstand deckt.
- **`AGENTS.md`/`harness/README.md`.** Beide unberührt — richtig: MR-045 verbietet ihnen Slice-Verweise, und keine der vier Ergänzungen ändert eine dort geführte Zusage. **Ausdrücklich nicht gemeldet:** MR-045 deckt den Verweis aus F-7 **nicht** — sein Geltungsbereich nennt genau diese zwei Dateien, nicht `harness/sensors/`. Ihn dorthin zu ziehen wäre `citation-stretched-beyond-scope`; F-7 steht deshalb auf der Archiv-Kopplung, nicht auf MR-045.
- **`AGENTS.md` §3.1.** Keine Host-Toolchain im Diff und keiner im Review: alle Läufe über `make`, der Bruch-Test mit `bash`/`sha256sum` in einer Kopie außerhalb des Repos.
- **`AGENTS.md` §3.3.** Der Beanspruchungs-Commit `af0c3fdd` ist ein reiner Move plus Roadmap-Flip; `make planning-check` ist grün, und die Botschaft begründet zutreffend, warum kein Pfad-Nachzug nötig war (auf den Plan zeigte noch nichts) und warum MR-070 nicht im Betreff steht.
- **`AGENTS.md` §3.7, Kommentare.** Keine Code-, Konfigurations- oder Skript-Kommentare im Diff.
- **Kommandos statt eingefrorener Zahlen.** Die Ergänzungen führen keine Bestands-Zahl ein, die beim nächsten Commit falsch wäre — die eine Zahl (`54 Dateien`) ist Teil einer zitierten Ausgabe, nicht eine Behauptung über den Bestand.
- **Alle 24 Dateien tragen einen `## Grenze`-Abschnitt.** Nachgezählt; die Formulierung des Slice stimmt.
- **`make gates` und `make doc-check`.** Beide grün, Ausgaben oben; keiner der sieben Befunde ist gate-sichtbar, und keiner ist es der Natur nach.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 4 | F-1, F-2, F-3, F-4 |
| LOW | 3 | F-5, F-6, F-7 |
| INFO | 2 | I-1, I-2 |

---

## Verdikt

**Blockiert** — vier MEDIUM. Der Anlass (DoD 1) ist die stärkste Arbeit in diesem Slice: der Bruch-Test reproduziert in allen drei Ausgängen exakt, die Formulierung gibt den Geltungsbereich des Skripts wieder statt seines Titels, und die `0.` vor der bestehenden `1.` ist eine bewusst rückwärtskompatible Setzung. Die **Inventur** dagegen trägt ihre eigene These nicht: zwei der vier Ergänzungen beschreiben einen Mechanismus, den das Werkzeug nicht hat (F-1, F-2), die Zahl „20 von 24" bricht an drei Stichproben (F-3), und die als *nicht entscheidbar* abgelegte vierundzwanzigste Datei war mit der Methode dieses Slice — beide Seiten lesen — entscheidbar (F-4). Bemerkenswert ist die Symmetrie: dreimal steht die Tatsache im Vertrags-Teil und wird beim Umdrehen **falsch** statt gar nicht aufgeschrieben. Das ist eine Aussage über den Ableiter des neuen Registereintrags (I-1) und gehört in seine Bewertung, bevor slice-213 aus ihm eine Regel macht. Der Registereintrag selbst trägt: die Abgrenzung gegen beide Nachbarn ist echt, beide nachgetragenen Belege sind aus den zitierten Reports gedeckt, und die Unsicherheit über die rückwärts gesammelte Schwelle steht als offenes Risiko dort, wo sie entschieden wird.
