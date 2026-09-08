# Review-Report: slice-222 — 2026-09-08 (R1)

**Review-Art:** Code-Review — geprüft wird der Diff gegen Slice-Plan,
Hard Rules und die zitierten `MR`-Einträge.

**Gegenstand:** Commit-Kette `c0730c13~1..HEAD` (`c0730c13`, `9c5f4d63`,
`ad47dc22`) plus der zum Prüfzeitpunkt **nicht committete** Arbeitsbaum-Stand
(Closure-Notiz in §6/§7 des Slice-Plans, drei `evidence/slice-222.md`).

**Skill:** `.harness/skills/reviewer.md` @ v1.16.0 / `ad47dc2` ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-09-08

> **Zitier-Form.** Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb Kennung statt Adresse: `slice-NNN` statt Lifecycle-Pfad,
> `make <target>` als Token statt Link auf die Sensor-Datei, eine
> Baseline-Stelle als Tag plus Pfad in Inline-Code
> (`v6.6.0` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext:**

- Slice-Plan `slice-222` (Fassung nach der Plan-Änderung vom 2026-09-08)
- `MR-011`-Pin-Serie · `MR-021` · `MR-025` · `MR-051` · `MR-055` · `MR-066` ·
  `MR-069` · `MR-070` · neu `MR-071`
- `DC-FA-TGT-001` (Modul `targets`), `DC-FA-PLAN-001`, `DC-QA-03`
- `AGENTS.md` §3.1 · §3.3 · §3.4 · §3.7 · §3.8 · §4 · §5 · §6
- Beobachtungs-Register: `BEO-ALL/pin-bump-mirrors-ungated`,
  `BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`,
  `BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage`
- Baseline `v6.6.0` · `templates/AGENTS.template.md` §4 ·
  `templates/.d-check.yml` · `templates/harness/README.template.md`

---

## Messformen dieses Laufs (vor den Zahlen, `AGENTS.md` §5)

Damit die Zahlen unten nachvollziehbar und nicht selbst ein Proxy sind:

- **Vorkommen** = jede Fundstelle der Zeichenkette `v6.5.0`
  (`grep -rIo "v6\.5\.0" | wc -l`). Zwei Fundstellen in einer Zeile zählen
  **zwei**.
- **Zeile** = jede Zeile mit mindestens einer Fundstelle
  (`grep -rIc "v6\.5\.0"`, Summe). Zwei Fundstellen in einer Zeile zählen
  **eins**.
- **Pfad-Verweis** = eine Zeile, die die Zeichenkette `.harness/baseline/v6.5.0`
  trägt (die Form, unter der die Zahl `77` des Plans reproduzierbar ist —
  siehe F-3). Die relative Form `../baseline/v6.5.0` wird **getrennt**
  ausgewiesen.
- **Makefile-Regel** / **Doku-Target**: die Formen des Moduls `targets` selbst
  (`internal/hexagon/core/rules/targets.go`) — Regelzeile
  `^([A-Za-z][A-Za-z0-9 _-]*):([^=]|$)`, je Feld ein Name; Doku-Target
  `` `make X` `` mit `X = [a-z][a-z0-9_-]*` **nur** in Zeilen, deren erstes
  Zeichen `|` ist.
- Vergleichsstand „vor der Ersetzung" ist `c0730c13`, ausgecheckt über
  `git archive`.

---

## Findings

### F-1 — Die Autoritäts-Umschaltung ist in den zwei Schlüsseln vollzogen und in beiden Beschreibungen des Gates nicht

- `kategorie`: MEDIUM
- `quelle`: `MR-025` §Geltungsbereich (*„jede Änderung an einer zugesagten
  Semantik — Grund-Code, Algorithmus-Schritt, **Config-Schlüssel**,
  Schwellenwert, Erkennungs-Form"*) und dessen Spiegel-Tabelle, die als
  Autoritäts-Doku ausdrücklich `AGENTS.md` **und** `harness/README.md` führt ·
  `AGENTS.md` §3.7 (Kommentar-Klasse *Kopplung*)
- `pfad`: `.d-check.yml:814` und `.d-check.yml:815` ·
  `harness/sensors/gate-consistency.md:5-7`
- `befund`: Der Blockkommentar **drei Zeilen über** den geänderten Schlüsseln
  sagt weiter *„jede Makefile-Regel ohne Eintrag in der Autoritäts-Doku
  (AGENTS.md §4) ⇒ gate-undocumented"* und *„exempt-targets bleibt leer —
  d-check dokumentiert alle Targets in AGENTS.md §4"*; darunter stehen
  `doc-tables: [harness/README.md]` und `authority: harness/README.md`. Die
  Vertrags-Sektion der Sensor-Datei sagt dasselbe: der geprüfte Kern sei
  *„in `AGENTS.md` §4 und in der Sensors-Tabelle dokumentierte `make X` ↔
  Makefile-Regeln, beide Richtungen"*. Damit nennt **keine** Beschreibung des
  Gates im Repo die tatsächliche Autorität, während `AGENTS.md` §4 zugleich
  maschinelle Deckung behauptet. Eine Spiegel-Liste nach `MR-025` liegt weder
  im Slice-Plan noch in einer Commit-Botschaft vor; die Datei
  `harness/sensors/gate-consistency.md` ist von der Kette nicht angefasst.
- `verifizierbar`: nein — kein Gate liest Kommentare oder Sensor-Prosa.
  Nachvollziehbar über `sed -n '810,822p' .d-check.yml` neben
  `sed -n '1,10p' harness/sensors/gate-consistency.md`.
- `klasse`: `gate-kommentar-beschreibt-den-alten-wert-der-zeile-darunter`
- `Kategorie-Begründung`: Basis MEDIUM nach dem Präzedenzfall
  `2026-09-07-slice-207-baseline-v650-review.md` F-2
  (`gate-kommentar-nennt-falschen-grund`, ebenfalls `.d-check.yml`). Nicht
  eskaliert, obwohl im Gate-Pfad — der Präzedenzfall lag im selben Pfad und
  blieb MEDIUM; Konsistenz der Skala sticht hier.

### F-2 — Eine Release-URL der Klasse „Spiegel 2" ist stehen geblieben, im Harness-Einstieg

- `kategorie`: MEDIUM
- `quelle`: DoD (2) des Slice-Plans · `BEO-ALL/pin-bump-mirrors-ungated`
  (Klasse *Release-/Tree-URLs*) · `MR-021`
- `pfad`: `harness/README.md:60`
- `befund`: Die Zeile verlinkt das vendorte Regelwerk auf
  `.harness/baseline/v6.6.0/regelwerk/` und sagt im selben Satz, es sei
  *„vendored aus dem self-contained `lab-regelwerk.zip`"* — mit der URL
  `…/releases/download/**v6.5.0**/lab-regelwerk.zip`. Die beiden Geschwister
  derselben Aussage sind gehoben (`AGENTS.md:36`, `harness/conventions.md:46`);
  diese dritte nicht. Die Zeile behauptet damit, der `v6.6.0`-Baum stamme aus
  dem `v6.5.0`-Bundle. Commit-Botschaft und `MR-071` führen für diese Klasse
  *„fünf Release-/Tree-URLs"* als nachgezogen und nennen die Fundstelle nicht
  unter den fünf bewusst stehen gelassenen Vergangenheits-Aussagen.
- `verifizierbar`: nein — kein Gate deckt diese Klasse (der Registereintrag
  sagt das ausdrücklich). Nachvollziehbar über
  `grep -rn "releases/download/v6" AGENTS.md harness/conventions.md harness/README.md`
  — zwei Treffer auf `v6.6.0`, einer auf `v6.5.0`.
- `klasse`: `release-url-spiegel-beim-pin-bump-uebersehen`

### F-3 — Die Frozen-Klassen-Tabelle zählt Zeilen und nennt sie Vorkommen; drei ihrer Spaltenwerte sind nicht reproduzierbar

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 (*„Vor einer Messung steht die Form ihres
  Gegenstands … Sie steht dort, wo die Zahl steht"*, seit `slice-210`) ·
  `MR-070` (die Tabelle **ist** das Artefakt, das `MR-070` verlangt) ·
  `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`
- `pfad`: Slice-Plan `slice-222` §3, Tabelle *Frozen-Klassen* ·
  `harness/conventions/MR-071-baseline-v660.md:36-42` · Commit-Botschaften
  `c0730c13` und `9c5f4d63`
- `befund`: Die Tabelle nennt keine Zählform. Gemessen am Stand `c0730c13`:
  108 Dateien (stimmt), aber **315 Vorkommen** und **302 Zeilen** — nicht 299.
  Die sechs Zeilensummen der Tabelle stimmen exakt gegen die **Zeilen**-Zählung
  (29 / 55 / 109 / 1 / 1 / 105); die Überschrift nennt sie *Vorkommen*, und
  die Tabelle summiert sich auf 300 gegen die Kopfzeile 299. Der Split
  *Pfad-Verweise* trägt in drei von fünf Zeilen einen Wert, der unter der Form,
  die die vierte reproduziert (Zeile enthält `.harness/baseline/v6.5.0`), nicht
  herauskommt: `done/`-Slices **27** statt 16, Review-Reports **6** statt 16,
  vendorter Baum **0** statt 29 (dessen 29 sind `blob/v6.5.0`-URLs in den
  Quelle-Kopfzeilen, also eine dritte Form). Fünf Zeilen der lebenden Klasse
  tragen die relative Form `../baseline/v6.5.0/` und stehen dadurch in der
  Spalte *nur Versionsnennung* — genau die Fehlklassifikation, die der Slice
  später als „übersehene Form" berichtet. Aus derselben Zahl folgt eine falsche
  Vorhersage: *„Die 33 eingefrorenen Pfad-Verweise werden beim Entfernen des
  Baums zu `target-missing`"* — gemessen sind es **4** in **2** Dateien.
- `verifizierbar`: ja, netzlos und ohne Werkzeug —
  `git archive c0730c13 | tar -x -C <tmp>` und darin
  `grep -rIo "v6\.5\.0" . | wc -l` (315) gegen
  `grep -rIc "v6\.5\.0" . | awk -F: '{s+=$NF} END {print s}'` (302). Die
  Vorhersage-Hälfte über den elften `ignore-refs`-Eintrag: in einer Kopie
  entfernt, `d-check` meldet **4** `target-missing`, Exit 1.
- `klasse`: `zaehlform-nicht-ausgeschrieben-zeilen-als-vorkommen`

### F-4 — Der elfte `ignore-refs`-Eintrag begründet sich mit einer Änderung, die seine Artefakte nicht betrifft

- `kategorie`: MEDIUM
- `quelle`: `MR-069` (das Ventil ist eine **deklarierte** Gate-Senkung — die
  Deklaration ist die Begründung) · `AGENTS.md` §5 (*„Eine zitierte Quelle
  trägt nur, was in ihrem Geltungsbereich steht"*) ·
  `BEO-ALL/citation-stretched-beyond-scope`
- `pfad`: `.d-check.yml:229-238` · wortgleich weitergetragen in
  `harness/conventions/MR-071-baseline-v660.md:49-53` und in der
  Commit-Botschaft `9c5f4d63`
- `befund`: Die Begründung lautet *„Ein Lift machte die Aussage still falsch
  (die v6.6.0-Vorlage führt gar keine §4-Tabelle mehr)"*. Die beiden
  ausgenommenen Artefakte verweisen ausschließlich auf
  `templates/harness/README.template.md` (vier Links, gemessen); die
  §4-Tabelle steht in `templates/AGENTS.template.md`, das keines der beiden
  Artefakte nennt. Der Grund, der hier tatsächlich trägt, ist ein anderer:
  `harness/README.template.md` führt in `v6.5.0` **39** Tabellenzeilen und in
  `v6.6.0` **40** — genau die Zahl, die `slice-217` als Messung festhält, würde
  durch den Lift falsch. Zweiter, kleinerer Punkt derselben Zeile: `refs`
  nimmt den ganzen Glob `.harness/baseline/v6.5.0/**` aus, gemessen nötig ist
  eine einzige Datei.
- `verifizierbar`: teilweise — die Trefferliste über
  `diff -u <alt>/templates/harness/README.template.md .harness/baseline/v6.6.0/templates/harness/README.template.md`
  (eine zusätzliche Tabellenzeile `make docs-check`) und
  `grep -c '^|'` auf beiden Ständen (39 / 40). Kein Gate liest Kommentare.
- `klasse`: `gate-kommentar-nennt-falschen-grund` — **zweite** Instanz nach
  `2026-09-07-slice-207-baseline-v650-review.md` F-2, an derselben Datei und
  in derselben Bump-Prozedur.

### F-5 — Der Adoptions-Teil liegt in einem Commit, dessen Botschaft und Kennungen ihn nicht nennen

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 (Traceability-Kennung; *„Eine Commit-Botschaft …
  behauptet nicht mehr, als die Arbeit trägt"*) · `DC-FA-COMMITS-001` ·
  Baseline `v6.6.0` · `regelwerk/modul-09-implementierung.md` §Minimal Agent
  Workflow (Schritt 8: berichten, was gelaufen ist)
- `pfad`: Commit `ad47dc22`
- `befund`: Der Commit trägt die zentrale Lieferung dieses Slice — Streichung
  der §4-Tabelle in `AGENTS.md` (−62 Zeilen), Umschaltung von `doc-tables` und
  `authority` in `.d-check.yml`, das Neu-Ankern der `reviewer.md`-Spanne und
  die Plan-Änderung in `slice-222` §1 — zusammen mit einem neuen, sachlich
  unabhängigen Plan-Dokument. Betreff, Rumpf und Kennungsliste nennen
  ausschließlich `slice-223` (dazu `MR-013`, `MR-059`, `MR-061`–`MR-064`,
  `DC-FA-PLAN-001`); weder `slice-222` noch `DC-FA-TGT-001` kommen vor, und
  der Rumpf beschreibt die Adoption mit keinem Wort. `make trace-check` bleibt
  grün, weil Kennungen **vorhanden** sind. `git log --grep=slice-222` findet
  damit den Commit nicht, der `AGENTS.md` §4 entfernt und die Gate-Autorität
  umgehängt hat. Der Bump-Commit `9c5f4d63` sagt an derselben Stelle noch
  *„Die Übernahme ist Sache des Adoptions-Teils, nicht dieses Commits"* — der
  angekündigte Teil hat keine eigene Botschaft bekommen.
- `verifizierbar`: ja —
  `git show --stat ad47dc22` neben `git log --grep=slice-222 --oneline`.
- `klasse`: `traegerkommit-nennt-den-falschen-vorgang`

### F-6 — Sechs lebende Stellen nennen `AGENTS.md` §4 als Deklarationsort, der dort nicht mehr existiert

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §4 (der Abschnitt sagt jetzt selbst *„Diese Datei führt
  die Liste nicht"*) · Maintainability
- `pfad`: `harness/README.md:115`, `:117`, `:119` (Spalte *Bindung*:
  „kein Gate · `AGENTS.md` §4") · `harness/conventions/MR-042-…:41` ·
  `harness/conventions/MR-046-…:62` · `tools/archive-wave/README.md:49`
- `befund`: Drei Zeilen der Werkzeug-Tabelle führen als **Bindung** einen
  Zeiger auf `AGENTS.md` §4; §4 verweist seinerseits auf ebendiese Tabelle
  zurück. Die Bindungs-Spalte ist für diese drei Targets damit inhaltsleer.
  Dieselbe Adresse benutzen `MR-042` (*„Die Proben sind ein `make`-Target
  (`AGENTS.md` §4, `make guard-probe`)"*), `MR-046` (*„ein `make`-Target, die
  Deklaration in `AGENTS.md` §4"*) und die Werkzeug-README als Ort einer
  Target-Deklaration.
- `verifizierbar`: nein — der Link auf `AGENTS.md` selbst löst auf, nur seine
  Aussage ist leer; `make doc-check` kann das nicht sehen.
- `klasse`: `zeiger-auf-einen-entleerten-abschnitt`

### F-7 — `AGENTS.md` §4 sagt maschinelle Deckung für eine Regelhälfte zu, die kein Mechanismus trägt

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.8 (*„Ein Modul verspricht nur über das, was es
  scannt"*) · `AGENTS.md` §5 (Grenzen-Absatz) ·
  `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`
- `pfad`: `AGENTS.md:410-418`
- `befund`: Der Abschnitt stellt die Regel auf — *„Kein Target nennen, das im
  Makefile nicht existiert — **auch nicht in Prosa**"* — und fährt fort
  *„Maschinell gehalten wird das von `make gate-consistency` in beide
  Richtungen"*. Das Modul `targets` liest ausschließlich Zeilen, deren erstes
  Zeichen `|` ist (`extractDocTargets`), und seit diesem Slice ausschließlich
  in `harness/README.md`. Prosa ist damit von keiner Richtung gedeckt, und
  `AGENTS.md` selbst ist seit der Umschaltung weder `doc-tables`-Eintrag noch
  `authority`. Die Vorlage `v6.6.0` · `templates/AGENTS.template.md` §4 stellt
  die Regel ohne diesen Zusatz auf. **Heute folgenlos** (gemessen: `AGENTS.md`
  trägt null Tabellenzeilen, und alle elf in Prosa genannten Targets
  existieren); die Zusage gilt aber jedem künftigen Zugang.
- `verifizierbar`: teilweise — die Scan-Menge steht in
  `internal/hexagon/core/rules/targets.go` (`tableRowLine`) und in
  `.d-check.yml:820-821`; ein Prosa-Phantom in `AGENTS.md` erzeugt
  reproduzierbar **null** Befunde.
- `klasse`: `deckungs-zusage-groesser-als-die-scan-menge`

### F-8 — Der neue Inline-Kommentar trägt einen Prosa-Pin statt eines auflösbaren Herkunfts-Feldes

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (*„Herkunft nur als **ein** auflösbares Feld nach
  dem Baseline-Schema (`DC-*`, `ADR-*`, `MR-*`, `seit welle-<NN>`)"*) ·
  `BEO-ALL/pin-bump-mirrors-ungated` (Klasse *Prosa-/Ellipsen-Pins*)
- `pfad`: `.d-check.yml:820`
- `befund`: Der Kommentar lautet `# der EINE Gate-Index (Baseline v6.6.0)`. Die
  Klammer ist eine Herkunfts-Angabe in einer Form, die das Baseline-Schema
  nicht führt und die nicht auflöst; das auflösbare Feld für genau diese
  Änderung wäre `MR-071`. Zugleich ist es ein **Versions-Pin in Prosa** —
  dieselbe gate-blinde Klasse, deren Nachziehen dieser Slice zum Gegenstand
  hat; beim nächsten Bump muss die Zeile von Hand mitwandern, und nichts
  meldet es.
- `verifizierbar`: nein — kein Gate prüft Kommentar-Klassen.
- `klasse`: `herkunfts-prosa-statt-aufloesbarem-feld`

### F-9 — Aus einer Zusage ist beim Chronik-Streichen eine Vollzugsmeldung geworden

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §5 (*„behauptet nicht mehr, als die Arbeit trägt"*)
- `pfad`: `AGENTS.md` §3.7, Absatz *Bestandsgrenze*
- `befund`: Aus *„der vorhandene Bestand **wird** mit dem v5.9.0-Bump
  umgestellt, nicht grandfathered"* wurde *„der vorhandene Bestand **ist**
  umgestellt, nicht grandfathered"*. Die Streichung der Bump-Chronik ist
  richtig und §3.7-konform; mitgegangen ist dabei die Zeitform. Der Satz
  behauptet jetzt einen abgeschlossenen Zustand, für den in der Kette keine
  Messung liegt. Alles Übrige der Chronik-Entfernung ist prüfbar folgenlos:
  der Auflösungs-Trigger *permanent* steht weiter, und die entfernten Angaben
  („seit dem v5.9.0-Bump", „deklariert mit dem v5.6.0-Bump") waren keine
  Herkunfts-Anker der Form `seit welle-<NN>` / `seit slice-<NNN>`.
- `verifizierbar`: nein.
- `klasse`: `zeitform-wechsel-beim-kuerzen`

### F-10 — Die produkt-seitige Beispiel-Konfiguration nennt weiter `authority: AGENTS.md`

- `kategorie`: INFO
- `quelle`: `MR-025` §Spiegel (*Emittierte Vorlage* · *Nutzer-Doku*)
- `pfad`: `docs/user/benutzerhandbuch.md:1307` ·
  `internal/adapter/driving/cli/config_template.go:243`
- `befund`: Beide Stellen zeigen `authority: AGENTS.md` als Beispiel; die
  adoptierte Vorlage `v6.6.0` · `templates/.d-check.yml` zeigt an derselben
  Stelle `authority: harness/README.md`. Das ist eine Produkt-Oberfläche und
  fällt unter Abgrenzung 1 des Plans (kein Urteil über den Regel-Delta) — es
  ist deshalb **kein** Verstoß, sondern der Spiegel, den die Umschaltung
  sichtbar macht und dessen Adresse (Folge-Slice) noch fehlt.
- `verifizierbar`: nein.
- `klasse`: `produkt-beispiel-hinter-der-adoptierten-vorlage`

### F-11 — Die Grenzen-Liste von `make baseline-verify` nennt den Fall nicht, der beim Bump zählt

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §5 (Grenzen-Absatz) · `MR-055`
- `pfad`: `harness/sensors/baseline-verify.md` §Grenze, Punkt 1
- `befund`: Punkt 1 lautet *„Geprüft wird die Auflösung, nicht das Ziel — ein
  Alias auf ein Verzeichnis passiert"*. Der Fall, der bei einer Pin-Hebung
  wirklich trägt, steht dort nicht: Ein vergessener Alias auf den **alten,
  noch vorhandenen** Baum löst auf und meldet grün, während der Lauf die alten
  Regeln liest. In dieser Kette wäre er aufgefallen, weil der `v6.5.0`-Baum im
  selben Commit verschwand und jeder stehengebliebene Alias tot gewesen wäre;
  in jedem Fenster, in dem beide Bäume existieren, meldet er nichts. Der
  Slice **benennt diese Grenze bereits** in seiner Closure-Notiz („Was offen
  bleibt (1)"); die dauerhafte Stelle — die Sensor-Datei — trägt sie nicht,
  und die Sensor-Datei ist von diesem Slice nicht angefasst.
- `verifizierbar`: ja — `readlink` auf einen Alias mit Ziel im alten Baum
  liefert, solange dieser existiert, ein gültiges Ziel;
  `bash tools/harness/fetch-baseline-cache.sh --verify` bleibt grün.
- `klasse`: `grenzen-liste-nennt-den-bump-fall-nicht`

---

## Negativbefunde

- **geprüft, ohne Befund: der Bruch-Test am Autoritäts-Wechsel, beide
  Richtungen, selbst gefahren** (Repo-Kopie unter dem Scratchpad,
  `d-check:latest`, `--enable targets` plus `FOCUS_DISABLE`). Ausgangslage:
  `774 Datei(en) geprüft, 0 Befund(e)`, Exit 0. **Richtung 1:** eine
  eingefügte Tabellenzeile `` | `make erfundenes-target` | … | `` in
  `harness/README.md` ⇒ `harness/README.md:122 erfundenes-target gate-phantom
  „dokumentiertes Target `make erfundenes-target` ohne Makefile-Regel"`,
  Exit 1. **Richtung 2:** eine angehängte Makefile-Regel
  `undokumentiertes-target:` ⇒ `Makefile:425 undokumentiertes-target
  gate-undocumented „Makefile-Regel `undokumentiertes-target` ohne Deklaration
  in der Autoritäts-Doku harness/README.md"`, Exit 1. Beide Befunde nennen
  `harness/README.md` als Autorität. Die vom Plan nach `MR-066` **vorab
  benannte Ersatz-Form ist damit eingelöst**, unabhängig nachgefahren und
  bestätigt; ihr Beleg im Repo nennt das Ergebnis, nicht die Ausgabe
  (`AGENTS.md` §6 Schritt 8) — was der Hausform der übrigen Commit-Botschaften
  entspricht und deshalb kein Befund ist. (Zum Zeitpunkt des Prüflaufs stand
  der Beleg nur in der noch nicht committeten Closure-Notiz; siehe Nachtrag.)
- **geprüft, ohne Befund: die Bijektion 54 ↔ 54.** Mit den Formen des Moduls
  gemessen: 54 Makefile-Regeln (54 eindeutig), 54 `make X`-Tokens in
  Tabellenzeilen von `harness/README.md`; `comm` in **beide** Richtungen leer.
  `exempt-targets: []` ist leer. Die Machbarkeits-Aussage des Plans trägt.
- **geprüft, ohne Befund: der gemessene Baseline-Delta.** Unabhängig
  reproduziert mit `diff -rq -I '<!-- Quelle:'` gegen den `c0730c13`-Baum:
  **sechs** von 26 Regelwerk-Dateien (`README.md`,
  `grundlagen-harness-dateien.md`, `modul-02-harness-bootstrap.md`,
  `modul-09-implementierung.md`, `modul-13-quality-gates.md`,
  `modul-15-observability.md`) und **fünf** Templates (`.d-check.yml`,
  `AGENTS.template.md`, `Makefile`, `harness/README.template.md`,
  `harness/conventions.template.md`) — deckungsgleich mit `MR-071`, Datei für
  Datei. 55 Dateien je Baum, keine neu, keine entfallen.
- **geprüft, ohne Befund: die `MR`-Kette.** `MR-067` liegt in
  `conventions/done/` und trägt in der Tabelle *Aufgelöste Adaptionen* seinen
  Nachfolger `MR-071`; `MR-071` steht in *Aktive Adaptionen*; §Baseline zeigt
  auf `#mr-071`. Beide Voll-Slug-Anker und beide Kurz-Anker existieren genau
  einmal. Die Ordinalzahl *dreizehnter Nachtrag* setzt die Reihe von
  *zwölfter* (`MR-067`) korrekt fort.
- **geprüft, ohne Befund: der `git mv` von `MR-067`.** Rename-Score `R086` —
  weit über der 50-%-Schwelle aus `AGENTS.md` §3.3; die mitgeführten
  Änderungen sind ausschließlich Link-Tiefen-Fixes (sieben Zeilen, `../..` →
  `../../..` bzw. `done/` entfernt), genau der Umfang, den die
  MR-Lifecycle-Ausnahme erlaubt. Der `Geltungsbereich` steht korrekt weiter
  auf `.harness/baseline/v6.5.0/` — die Über-Hebung ist zurückgenommen.
- **geprüft, ohne Befund: die acht Baseline-Symlinks** unter `.claude/rules/`.
  Alle acht zeigen auf `.harness/baseline/v6.6.0/regelwerk/…` und lösen auf;
  die drei repo-internen Aliase (`AGENTS.md`, `conventions.md`,
  `harness-README.md`) sind unberührt. `bash tools/harness/fetch-baseline-cache.sh
  --verify` meldet `verify ok (54 Dateien, vollständig)`, Exit 0. Die
  Blindstelle dieses Sensors ist als F-11 geführt, nicht als Fehler dieser
  Kette.
- **geprüft, ohne Befund: die Substanz der entfernten §4-Tabelle**, Wort für
  Wort an sechs Stichproben. `make test` / `MR-048`: *„jeder verdrahtete
  Hook-Pfad existiert … Permission-Sperrliste deckt jeden Namen … als **ganze**
  Befehlsklasse"* samt Grenze *„Kein Beleg, dass die Durchsetzung läuft"* →
  `harness/sensors/test.md:15-31`. Digest-Achsen: *„Digests haben keine
  Ordnung"* und *„`nonroot` führt keine Version"* →
  `harness/sensors/runtime-base-digest.md:13,26-27`. Action-Pins:
  *„Verglichen wird der Kommentar gegen upstream, nicht der SHA"* →
  `harness/sensors/checkout-pin-freshness.md:18`. Proben-Targets: *„Ohne
  wiederholbare Proben wäre seine Zusage eine Erinnerung"* → `Makefile:161`,
  `tools/harness/fetch-baseline-cache.sh:143`, `MR-042`. `make coverage-gate`:
  Kalibrierungs-Bindung → `Makefile:35,73` und `harness/README.md:75`.
  `make baseline-freshness`: *„Prereleases überspringt"* und der Herkunfts-Anker
  `seit slice-215` → `harness/sensors/baseline-freshness.md:32,54`,
  `Makefile:313`. **Eine** Einbuße, bewusst notiert statt gemeldet: die
  Aufzählung der neun `baseline-probe`-Fälle stand nur in §4 und steht jetzt
  nur noch im Skript; Zahl und Zeiger tragen `harness/README.md:124` und
  `harness/sensors/baseline-verify.md:48`.
- **geprüft, ohne Befund: kein Herkunfts-Anker ist mit der Tabelle
  verschwunden.** Die entfernten 55 Zeilen tragen genau einen (`seit
  slice-215`), und der steht weiter in `Makefile:313` und in der Sensor-Datei.
  Kein `liegt in`-Pflichtfeld im ganzen Repo zeigt auf `AGENTS.md §4`; die
  Anker-Paarung ist unberührt. **Nebenbei behoben:** §4 trug mit *„die drei
  Umkehr-Proben aus slice-190"* einen Slice-Verweis, den `MR-045` verbietet.
- **geprüft, ohne Befund: kein Phantom-Target in `AGENTS.md`.** Alle elf in
  Prosa genannten `make X` existieren als Makefile-Regel (`comm` leer). Die
  §4-Regel selbst ist heute eingehalten; ihre Deckungs-Zusage ist als F-7
  geführt.
- **geprüft, ohne Befund: die sieben Prosa-Nennungen.** Genau sieben lebende
  Stellen tragen den Tag ohne Pfad und stehen auf `v6.6.0`:
  `.d-check.closure.yml:185`, `MR-021:20`, `MR-053:9`, `MR-049:14`,
  `roadmap.md:13`, `observations/README.md:18`, `AGENTS.md:604`.
- **geprüft, ohne Befund: die verbliebenen `v6.5.0`-Fundstellen außerhalb der
  Frozen-Klassen.** 35 Dateien nennen den Tag noch; 28 davon sind die
  deklarierte Frozen-Menge, sieben sind neu oder legitim: `MR-071` und
  `conventions/done/MR-067` (beschreiben die Hebung), der Slice-Plan und die
  neue `evidence/slice-222.md`, die vier `.d-check.yml`-Tombstone-Zeilen und
  zwei wörtliche Fremd-Zitate im eingehenden CR. **Einzige Ausnahme:**
  `harness/README.md:60` (F-2).
- **geprüft, ohne Befund: die `d-check:cite`-Spannen.** `make doc-check` läuft
  auf dem Arbeitsbaum grün (`775 Datei(en) geprüft, 0 Befund(e)`, Exit 0),
  und `citations` ist Teil davon; die neu geankerte Spanne
  `reviewer.md:59 → AGENTS.md:412-413` deckt das Zitat *„Halluzinierte Gates
  sind die häufigste Form von Harness-Lüge"* korrekt. Die vielen
  `v6.5.0`-Direktiven in `done/` und `docs/reviews/` sind über
  `citations.scope.ignore` ausgenommen — deklariert, mit benannter Grenze, und
  nicht durch diesen Slice entstanden.
- **geprüft, ohne Befund: die Adoption gegen die Vorlage.** Der übernommene
  Satz ist wortgleich mit `v6.6.0` · `templates/AGENTS.template.md` §4
  (*„Der Gate-Index steht **einmal**, in `harness/README.md` §Sensors — dort
  steht auch die *Bindung* jedes Targets. Diese Datei führt die Liste
  nicht."*), die repo-eigene Ergänzung (Weg zur `DC-*`-ID, zur ADR, zum
  Carveout) ist als solche erkennbar. Die Template-`.d-check.yml` setzt beide
  Schlüssel auf `harness/README.md` — übernommen.
- **geprüft, ohne Befund: `slice-221`.** Angefasst wurden nur seine zwei
  `d-check:cite`-Spannen (Tag-Nachzug, `MR-021`/`MR-051`); die Zeilennummern
  `268-269` / `274-274` stimmen weiter, weil `modul-05-planning-harness.md`
  nicht im Delta liegt. Die Abgrenzungs-Formulierung *„weder beansprucht noch
  nachgezogen"* ist gegen diesen Pin-Nachzug ungenau, die Absicht aber
  gewahrt: an der Substanz von `slice-221` ist nichts geändert.
- **geprüft, ohne Befund: Referenz-Richtung und Spec-Straten.** Die einzigen
  Änderungen an `spec/architecture.md` und `spec/spezifikation.md` sind zwei
  Baseline-Pin-Nachzüge; kein Spec-Stratum hat einen Abwärts-Verweis
  bekommen. Das Modul `matrix` läuft im grünen `make doc-check` mit.
- **geprüft, ohne Befund: `AGENTS.md` §3.1.** Für diesen Review kamen
  ausschließlich `git`, POSIX-Werkzeuge (`grep`, `awk`, `sed`, `comm`,
  `diff`, `find`, `tar`), `bash` und `docker` zum Einsatz; kein Host-Go, kein
  Host-Paketmanager, kein Host-Skript-Interpreter. Der Arbeitsbaum wurde
  nicht verändert — die Bruch-Tests liefen in einer Kopie unter dem
  Scratchpad.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 5 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:**
`gate-kommentar-beschreibt-den-alten-wert-der-zeile-darunter` ·
`release-url-spiegel-beim-pin-bump-uebersehen` ·
`zaehlform-nicht-ausgeschrieben-zeilen-als-vorkommen` ·
`gate-kommentar-nennt-falschen-grund` ·
`traegerkommit-nennt-den-falschen-vorgang` ·
`zeiger-auf-einen-entleerten-abschnitt` ·
`deckungs-zusage-groesser-als-die-scan-menge` ·
`herkunfts-prosa-statt-aufloesbarem-feld` ·
`zeitform-wechsel-beim-kuerzen` ·
`produkt-beispiel-hinter-der-adoptierten-vorlage` ·
`grenzen-liste-nennt-den-bump-fall-nicht`

**Zwei Klassen sind Wiederholungen und gehören in den Zähler:**
`gate-kommentar-nennt-falschen-grund` (F-4) ist die zweite Instanz nach dem
`slice-207`-Review, an derselben Datei und in derselben Prozedur.
`release-url-spiegel-beim-pin-bump-uebersehen` (F-2) und
`zaehlform-nicht-ausgeschrieben-zeilen-als-vorkommen` (F-3) sind je eine
weitere Instanz von `BEO-ALL/pin-bump-mirrors-ungated` bzw.
`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand` — beide zählen bereits
über der Schwelle, was die Kategorie nicht senkt, aber die Frage schärft, ob
die verkörperte Antwort (`MR-070`, die Suchform-Regel) die Klasse trifft.

**Kein HIGH — und warum das eine Aussage ist:** Die vier
HIGH-Anker mit Mechanik-Bezug wurden einzeln geprüft. Kein Gate hat einen
still-grünen Pfad **bekommen** (beide Richtungen des umgeschalteten Moduls
feuern, selbst gemessen), kein Kern-Modul ist berührt, keine Schwelle gesenkt,
keine Suppression, kein Netzzugriff. Was bleibt, sind Beschreibungen, die
hinter der Konfiguration zurückstehen — schmerzhaft in einem Repo, dessen
Gegenstand genau das ist, aber ohne Mechanik-Folge.

## Verdikt

**Merge-blockierend:** ja — fünf MEDIUM.

F-1, F-2 und F-4 sind punktuell und ohne Risiko zu beheben (drei Textstellen
und eine URL). F-3 verlangt mehr als eine Korrektur: Die Zahlen stehen in
einem `MR`-Eintrag, der **lebend** ist und den nächsten Bump anleitet — dort
gehört die Zählform neben die Zahl, nicht nur die korrigierte Zahl. F-5 ist
nicht mehr reparabel, ohne Historie umzuschreiben; hier trägt die
Closure-Notiz die Aufgabe, den Träger-Commit der Adoption zu benennen, damit
der Vorgang auffindbar bleibt.

**Nicht blockierend, aber vor der Closure zu entscheiden:** F-6 bis F-8 sind
LOW; F-10 nennt einen Spiegel ohne Adresse — entweder ein Folge-Slice oder
eine ausdrückliche Ablage unter Abgrenzung 1.

---

## Nachtrag (noch im selben Lauf, vor der Übergabe)

Während dieses Reviews hat der Autor **F-5 selbst gefunden und behoben** —
nicht auf diesen Report hin, sondern beim Nachmessen einer Größenangabe für
die Closure-Notiz. Der Befund bleibt oben stehen, wie er gegen den geprüften
Stand galt; hier steht, was seither gilt:

- `ad47dc22` ist in `cd0b6940` (slice-222: `AGENTS.md` §4, `.d-check.yml`,
  `reviewer.md`, Plan-Änderung) und `a7c352f9` (slice-223: nur die neue
  Plan-Datei) zerlegt. **Der Baum ist bit-identisch** — `git rev-parse
  ad47dc22^{tree}` und `HEAD^{tree}` liefern beide `b2d02104…`; alle übrigen
  Findings gelten unverändert.
- Die neue Botschaft `cd0b6940` nennt `slice-222`, beschreibt die Adoption und
  **trägt den Bruch-Test in beide Richtungen** als `MR-066`-Ersatzform.
- Neu im Register: `BEO-ALL/liefer-punkt-in-fremdem-commit` mit
  `evidence/slice-222.md`. Die Abgrenzung gegen
  `BEO-ALL/commit-message-overclaims-work` ist dort sauber gezogen: dort
  behauptet die Botschaft zu viel, hier ist sie richtig und am falschen
  Vorgang befestigt.
- **Nachgeprüft und korrekt:** `AGENTS.md` fällt von **59 826** auf **40 320**
  Bytes (−32,6 %), §4 verliert **43** Tabellenzeilen. Die Closure-Notiz führt
  diese Zahlen und weist ausdrücklich aus, dass die 9 785 Bytes, die
  `slice-221` an Zell-Kürzung beigetragen hatte, **nicht** mitgezählt werden —
  die frühere Angabe „69 611 auf 40 320 (−42 %)" hätte sie mitgenommen.

**F-5 gilt damit als erledigt.** Die vier übrigen MEDIUM (F-1 bis F-4) stehen
unverändert; F-3 ist der einzige, dessen Zahlen bereits in einem **lebenden**
`MR`-Eintrag stehen.

**Übergabe:** Findings gehen an den Implementer. Dieser Report ist Lauf-Beleg
und wird über Läufe hinweg nicht wieder gelesen; die Finding-Klassen gehen in
§7 des Slice-Plans und von dort in den Zähler. Er ersetzt keine
Verifikation — DoD-Abhakung und Spec-Konformität prüft der Verifier separat.
