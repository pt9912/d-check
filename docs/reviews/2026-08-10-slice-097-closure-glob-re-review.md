# Re-Review-Report: slice-097 — 2026-08-10

**Review-Art:** Code — **Re-Review** (bestätigend). Geprüft wird nicht der Slice
noch einmal von vorn, sondern die **Heilung**: hält jeder Erst-Befund am Lauf,
und hat die Heilung Neues aufgemacht? Kein Verifikations-Lauf — die DoD-Abhakung
bleibt ausdrücklich außen vor.

**Gegenstand:** [slice-097](../plan/planning/in-progress/slice-097-closure-glob-entkopplung.md),
Heilungs-Commit `16def65` (Diff-Range `21de5a3..HEAD`); Gesamt-Slice
`39a3c6a..HEAD` (fünf Commits). Vorlauf-Beleg: der blockierende Erst-Report
[2026-08-10-slice-097-closure-glob-review.md](2026-08-10-slice-097-closure-glob-review.md).

**Skill:** `.harness/skills/reviewer.md` @ 1.3.0 ·
**Modell:** claude-opus-5[1m] · **Datum:** 2026-08-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne diese Liste
ist der Lauf nicht reproduzierbar):

- [slice-097](../plan/planning/in-progress/slice-097-closure-glob-entkopplung.md)
  §3 Abnahme-Punkt 2, §4 DoD, §5 Risiken
- [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md)
  (Entscheidungen 1–6, Verglichene Alternativen, Konsequenzen, Fitness Function,
  Re-Evaluierungs-Trigger, Geschichte) und
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (Lastenheft-Version 0.53.1, §Akzeptanzkriterien) und
  [§`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C1/C2, §2-Schema-Zeilen, §4-Grund-Code-Tabelle
- [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) (Byte-Identität)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules)
- der Erst-Report als Prüf-Checkliste (alle neun Befunde, auch LOW und INFO)

**Läufe dieses Reviews.** Alle Fixtures liegen in einem Temp-Verzeichnis
außerhalb des Repos, alle Läufe netzlos und read-only. Gebaut wurden **zwei**
Vergleichs-Images: HEAD (`make build`) und der Stand **vor** der Heilung
(`21de5a3`, über ein `git archive`-Auspack ins Temp-Verzeichnis,
`make build IMAGE=d-check-pre`). Gefahren: `make build` (dreimal), `make test`
(sechsmal, davon fünfmal mutiert), `make gates`, `make verify-closure-notes`,
`make adr-check` sowie rund 35 Fixture-Läufe gegen die beiden Images.
`make gates` grün (356 Dateien, 0 Befunde), `make verify-closure-notes` grün
(327 Dateien, 0 Befunde), `make adr-check RANGE=39a3c6a..HEAD` grün
(356 Dateien, 0 Befunde) — jeweils **vor** Ablage dieses Reports; mit ihm zählen
`make doc-check` 357 und `make verify-closure-notes` 328 Dateien, beide
weiterhin 0 Befunde. Der Arbeitsbaum ist am Ende unverändert
(`git status --short` leer, Prüfsummen der vier mutierten Dateien gegen die
Sicherungskopien gleich — wiederhergestellt per Rückkopie, **nicht** per
`git checkout`).

---

## Erst-Befund → Status → Beleg

| Erst-Befund | Kat. | Status | Beleg (Lauf dieses Reviews) |
|---|---|---|---|
| **F-1** Abschnitts-Grenze = Dokument-Grenze | MEDIUM | **geheilt — durch Rücknahme** | Der semantische Teil von `.d-check.closure.yml` ist gegen `39a3c6a` **identisch** (Diff nur Kommentar-Zeilen: `dir` und `boilerplate` unverändert, kein `glob`, kein `heading-pattern`). Kandidaten-Zählung am Lauf: 95 unter dem geerbten Filter, 110 unter `*.md` — die 15 Wellen-Dokumente sind wieder draußen, das Falsch-Negativ damit unerreichbar. `make verify-closure-notes` grün |
| **F-2** doppeltes Akzeptanzkriterium | MEDIUM | **geheilt** | Vier Fixtures gegen das HEAD-Image: (a) nur `welle-01-x.md` + `closure.glob: "*.md"` ⇒ Kandidat vorhanden, Exit 0; (b) dieselbe Datei **ohne** `closure.glob` ⇒ `closure-note-missing` auf dem Verzeichnis, Exit 1; (c) nur `slice-001-x.md` + `closure.glob: "*.txt"` ⇒ Exit 1; (d) `slice-glob: "*.md"` geerbt ⇒ kein Closure-Befund. Die zusammengezogene Fassung ist in allen vier Richtungen wahr; die alte, auf `planning.slice-glob` fixierte Fassung ist entfernt |
| **F-3** Bestandszahlen 96/111 | LOW | **teilweise** | 95/110 am Lauf bestätigt (Verzeichnis-Listing **und** Lauf mit nie passendem `heading-pattern`: 110 bzw. 95 Befunde). Slice-Plan, Fitness Function, Konsequenzen und Config-Kommentar sind korrigiert — die ADR-Tabellenzeile **nicht** (siehe N-2) |
| **F-4** YAML-`null` umgeht Exit 2 | LOW | **als Grenze dokumentiert, Text am Lauf wahr** | Schritt C2 sagt jetzt „YAML-`null` gilt als **abwesend**, nicht als leer — wie bei `min-sentences`". Gemessen: `glob:`, `glob: ~` und `glob: null` ⇒ Exit 0 mit geerbtem Filter; `glob: ''` ⇒ Exit 2; `min-sentences:` ⇒ Exit 0, `min-sentences: 0` ⇒ Exit 2. Die Analogie trägt, das Verhalten ist unverändert und jetzt benannt |
| **F-5** unquotierter Glob in der Nullmengen-Meldung | LOW | **geheilt — mit unbenannter Nebenwirkung** | Whitespace-Glob liefert jetzt „keine Datei nach `"   "`", `*.txt` liefert „nach `"*.txt"`". Kein Golden-Test und kein Fixture pinnt den Text (einziger Zitat-Test prüft `strings.Contains` und überlebt die Quotierung). **Aber:** derselbe Pfad läuft ohne den neuen Schlüssel — Byte-Identität gebrochen (N-1) |
| **F-6** Klartext nennt den Kandidaten „Slice" | LOW | **teilweise** | „Slice" ⇒ „Kandidat" ist in **beiden** Texten nachgezogen (§4-Zeile und `--doctor`). Der zweite Halbsatz des Befunds — der Leerlauf-Fall fehlt in beiden — ist **nur** im `--doctor`-Klartext geheilt; die §4-Zeile nennt ihn weiter nicht (siehe N-3). Der `AllReasons`-Lockstep ist unberührt (beide Deckungs-Tests prüfen Codes, keinen Text) |
| **F-7** Handbuch führt den Schlüssel nicht | LOW | **nicht geheilt** | `docs/user/benutzerhandbuch.md:1170` listet unverändert `dir`, `heading-pattern`, `min-sentences`, `boilerplate`; die Exit-2-Ränder daneben nennen weiter nur `min-sentences` und `boilerplate`; das Aufgaben-Beispiel §4.17 kennt den Schlüssel ebenfalls nicht. Die Versions-Tabelle endet bei 1.44 / v0.53.0. Der DoD-Punkt „Release" ist offen — die Lücke ist erwartbar, aber sie **steht noch** |
| **F-8** sechster Rückbau ungefangen | LOW | **geheilt — Test selbst nachgewiesen** | Mutation „Prüfung an `closure.dir` koppeln" (`c.Glob != nil && c.Dir != ""`, `internal/adapter/driven/configyaml/configyaml.go:969`) angewandt ⇒ `make test` **rot**, genau in `TestDecode_ClosureRandGiltAuchOhneDir`. Der Test ist nicht tautologisch: ohne die Mutation grün, mit ihr rot, und er deckt einen Fall ab, den keine bestehende Tabelle abdeckte |
| **F-9** Nachbar-Schlüssel, gegenteiliger Config-Rand | INFO | **unverändert offen** | `slice-glob: ''` ⇒ Exit 0 (stiller Rückfall), `closure.glob: ''` ⇒ Exit 2. Kein Dokument benennt die Asymmetrie; die Heilung fasst sie nicht an. Als INFO auch nicht zugesagt |

## Findings

### N-1 — Die Heilung bricht die Byte-Identitäts-Zusage, die der Slice gar nicht anfassen wollte

- `kategorie`: MEDIUM
- `quelle`: [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md)
  §Fitness Function („Ohne den Schlüssel byte-identisch zum Stand vor dieser
  Änderung") und Entscheidung 6 („ohne ihn ist der Befundsatz byte-identisch.
  Ein Adopter merkt beim Update nichts")
- `pfad`: `internal/hexagon/core/rules/planning.go:104` und
  `internal/hexagon/core/app/diagnose.go:119`
- `befund`: Beide Textänderungen liegen auf Pfaden, die **ohne** den neuen
  Schlüssel erreicht werden. Gemessen an einer identischen Fixture (Closure-
  Verzeichnis ohne Kandidaten unter dem **geerbten** Filter, kein `closure.glob`
  gesetzt) gegen zwei Images, `21de5a3` und HEAD: die `--json`-Meldung weicht ab
  („nach slice-*.md" gegen „nach `\"slice-*.md\"`"), und der `--doctor`-Klartext
  weicht ab („Abgeschlossener Slice … (oder Closure-Verzeichnis fehlt —
  fail-closed)" gegen „Abgeschlossener Kandidat … (oder Closure-Verzeichnis
  fehlt bzw. ist leer — fail-closed)"). Der Klartext trifft jeden Konsumenten mit
  **irgendeinem** `closure-note-missing`-Befund, nicht nur den Nullmengen-Fall.
  Keine der drei Historien führt die Änderung: die Lastenheft-Zeile 0.53.1 nennt
  nur das zusammengezogene Kriterium, die Spezifikations-Historie nur die
  §4-Zeile und `YAML-null`, `CHANGELOG.md` §Unreleased ist leer. Das Repo
  behandelt genau diese Klasse sonst als konsumenten-relevante Release-Notiz
  (Lastenheft-Zeile 0.52.3 markiert eine Klartext-Änderung ausdrücklich so).
- `verifizierbar`: ja — Vor-Image aus `21de5a3` bauen
  (`make build IMAGE=d-check-pre`) und beide Images auf derselben Fixture mit
  `--json` und `--doctor` fahren; die Ausgaben differieren zeilenweise.
- `klasse`: Byte-Identitäts-Zusage von der Heilung mitgerissen, ohne Historie

### N-2 — Die ADR-Tabellenzeile führt weiter 96, fünf Zeilen unter der auf 95 korrigierten Messung

- `kategorie`: LOW
- `quelle`: [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md)
  §Kontext und §Geschichte („die Bestandszahlen sind auf die gemessenen 95/110
  korrigiert")
- `pfad`: `docs/plan/adr/0051-eigener-kandidaten-filter-closure.md:48`
- `befund`: Zeile 43 derselben Passage sagt „110 Markdown-Dateien, von denen der
  heutige Filter **95** sieht"; die Tabelle darunter sagt in derselben Variante
  „96". Die identische Tabelle im Slice-Plan ist auf 95 korrigiert, die ADR-
  Fassung nicht. Gemessen: 95 (Verzeichnis-Listing **und** Lauf mit nie
  passendem `heading-pattern`). Die §Geschichte behauptet die Korrektur als
  abgeschlossen; der Rückstand ist damit nicht nur eine Zahl, sondern eine
  Selbstaussage, die der eigene Text widerlegt. Mit der Closure geht die ADR auf
  `Accepted` und ist ab dann immutabel; die Zahl wandert außerdem: sobald
  `slice-097` selbst in `done/` liegt, sind es 96 Slices — die falsche Zeile
  würde später zufällig richtig aussehen und die Messung unauffindbar machen.
- `verifizierbar`: ja — Lauf mit `closure.glob: "*.md"` bzw. ohne den Schlüssel,
  je mit einem `heading-pattern`, das nichts trifft: die Zahl der
  `closure-note-missing`-Befunde ist die Kandidatenzahl (110 bzw. 95).
- `klasse`: Bestandszahl in der Fitness Function nicht reproduzierbar

### N-3 — §4-Zeile und `--doctor`-Klartext sind nach der Heilung in der anderen Richtung auseinander

- `kategorie`: LOW
- `quelle`: [§`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C2 („Ebenso fail-closed: null Kandidaten") und §4-Grund-Code-Tabelle
- `pfad`: `spec/spezifikation.md:2376` gegen
  `internal/hexagon/core/app/diagnose.go:119`
- `befund`: Der `--doctor`-Klartext nennt jetzt drei Auslöser („ohne
  Closure-Notiz-Abschnitt", „Verzeichnis fehlt", „bzw. ist leer"), die §4-Zeile
  nur zwei — die Wendung „ist leer" kommt in der gesamten Spezifikation nicht
  vor. Vor der Heilung fehlte der Leerlauf-Fall in **beiden** Texten (Erst-Befund
  F-6, zweiter Halbsatz); jetzt fehlt er in einem. Ein Konsument, der den
  Grund-Code in §4 nachschlägt, findet die Ursache seines Befunds dort nicht,
  während `--doctor` sie nennt. Die beiden Lockstep-Tests decken das nicht ab:
  sie prüfen die Deckung der **Codes**, nicht die der Texte.
- `verifizierbar`: ja — Fixture mit gesetztem `closure.dir` ohne passenden
  Kandidaten gegen das Image, `--doctor` gelesen, daneben die §4-Zeile.
- `klasse`: Grund-Code-Klartext und §4-Zeile ohne Text-Lockstep

### N-4 — Eine verworfene Alternative ist mit einer falschen Aussage über `path.Match` begründet

- `kategorie`: LOW
- `quelle`: [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md)
  §Verglichene Alternativen
- `pfad`: `docs/plan/adr/0051-eigener-kandidaten-filter-closure.md:124`
- `befund`: Die Zeile begründet die Ablehnung mit „`path.Match` kennt keine
  Suffix-Negation; jedes Muster, das `welle-70-fence-lexik.md` trifft, trifft
  auch `welle-70-results.md`". Das Benutzerhandbuch dokumentiert für genau diese
  Glob-Familie die negierte Zeichenklasse `[^…]` (§5, Glob-Syntax). Gemessen am
  eigenen Bestand mit `closure.glob: '*-*-[^r]*.md'`: 97 Kandidaten, darunter
  **alle vier** Wellen-Plan-Dokumente und **keine** der elf
  `welle-*-results.md` — die behauptete Implikation ist am Lauf widerlegt (das
  Muster verfehlt zwei Slice-Dateien, die Aussage der Zeile betrifft aber die
  Wellen-Trennung). Was `path.Match` fehlt, ist Alternation, nicht Negation. Die
  ADR friert die Begründung mit dem Statuswechsel ein; ein späterer
  Re-Evaluierungs-Trigger liest sie als Fakt.
- `verifizierbar`: ja — derselbe Lauf mit dem genannten Glob und einem nie
  passenden `heading-pattern`; die Befund-Dateiliste zeigt die Trennung.
- `klasse`: verworfene Alternative mit widerlegbarer Begründung eingefroren

### N-5 — Der neu geschriebene Config-Kommentar nennt ein Minimum, das die Nachmessung nicht bestätigt

- `kategorie`: LOW
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  §Closure-Struktur (Substanz-Untergrenze) · Maintainability
- `pfad`: `.d-check.closure.yml:47`
- `befund`: Der Kommentar sagt „Gemessen am eigenen Bestand: 95/95 Slices tragen
  einen Abschnitt, das Minimum liegt bei 5 Satzende-Zeichen außerhalb Code — die
  Schwelle ist ein Boden, keine Decke". Die erste Zahl ist am Lauf richtig, die
  zweite nicht: mit `min-sentences: 6` und `min-sentences: 7` bleibt der Bestand
  befundfrei, erst `min-sentences: 8` liefert einen `closure-note-thin` mit der
  Meldung „trägt 7 Satzende-Zeichen". Das gemessene Minimum ist **7**, nicht 5.
  Der Satz ist im Heilungs-Commit angefasst worden (92/92 ⇒ 95/95); die dritte
  Zahl derselben Messung ist ungeprüft mitgewandert. Wer den Abstand zur Schwelle
  aus diesem Kommentar liest, unterschätzt ihn um zwei.
- `verifizierbar`: ja — dasselbe Profil mit `min-sentences` von 4 bis 8 gegen das
  Image; der erste Wert mit einem Befund benennt das Minimum.
- `klasse`: Bestandszahl in der Fitness Function nicht reproduzierbar

### N-6 — Siebter Rückbau: der Config-Rand der Nachbar-Schlüssel lässt sich weiter an `closure.dir` koppeln, ohne dass ein Test rot wird

- `kategorie`: LOW
- `quelle`: [§`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C1 (Exit-2-Liste am Config-Rand)
- `pfad`: `internal/adapter/driven/configyaml/configyaml.go:959` und
  `internal/adapter/driven/configyaml/configyaml.go:991`
- `befund`: Der neue Test pinnt die Unabhängigkeit vom Aktivierungs-Schalter für
  `glob` (zweimal) und `min-sentences` — nicht für `heading-pattern` und
  `boilerplate`. Selbst ausgeführt: beide Prüfungen auf `c.Dir != ""` verengt ⇒
  `make test` **grün** (Exit 0). Gegen die daraus gebauten Images kippt derselbe
  Rand messbar: `heading-pattern: '^(['` **ohne** `closure.dir` meldet am
  HEAD-Image Exit 2 mit Nennung des Schlüssels, am mutierten Image Exit 1 mit
  einem unbezogenen `planning-drift`; für `boilerplate: ['']` dasselbe Bild. Die
  Erst-Befund-Klasse F-8 ist damit für zwei von vier Schlüsseln geschlossen und
  für zwei offen — die Testtabelle setzt dort weiterhin durchgängig `dir: done`.
- `verifizierbar`: ja — die genannte Verengung anwenden, `make test` fahren
  (bleibt grün), dann `make build` und dieselbe Config gegen beide Images.
- `klasse`: Config-Rand-Verhalten bei inerter Fähigkeit ungetestet

### N-7 — Die `--print-config`-Vorlage stellt den Verweis-Default als Literal dar

- `kategorie`: INFO
- `quelle`: [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md)
  Entscheidung 2 (Default als Verweis, nicht als kopiertes Literal)
- `pfad`: `internal/adapter/driving/cli/config_template.go:128`
- `befund`: Die Zeile zeigt `glob: 'slice-*.md'` mit dem Zusatz „weglassen ⇒
  slice-glob". Wer den auskommentierten `closure`-Block als Ganzes übernimmt —
  die Form, die die Vorlage anbietet —, pinnt damit genau das Literal, dessen
  Vermeidung die ADR als tragendes Argument führt; ändert er später
  `slice-glob`, wandert die Closure-Menge nicht mit und der Lauf läuft
  fail-closed in die Nullmenge. Die Spezifikation verlangt für `--print-config`
  nur „kommentierte Beispiele", nicht die Wiedergabe des Defaults — der Rand ist
  daher zulässig, aber im Vertrag nicht als bewusst benannt.
- `verifizierbar`: nein (kein Gate) — Sichtprüfung gegen die §2-Schema-Zeile
  (Default „Wert von `planning.slice-glob`").
- `klasse`: Vorlage bildet einen Verweis-Default als Literal ab

## Negativbefunde

- geprüft, ohne Befund: **Rücknahme der eigenen Weitung.** Der wirksame Teil von
  `.d-check.closure.yml` ist gegen den Vor-Slice-Stand identisch — der Diff
  gegen `39a3c6a` berührt ausschließlich Kommentarzeilen; `dir` und die
  fünf `boilerplate`-Einträge stehen unverändert, `glob` und `heading-pattern`
  sind abwesend. Der eigene Lauf prüft damit wieder genau dieselben 95
  `slice-*.md` wie vor dem Slice, mit demselben Default-Muster.
- geprüft, ohne Befund: **Die beiden nachgetragenen Wellen-Closure-Notizen
  tragen.** Beide Notizen (Wellen-Plan-Dokumente 69 und 70 in `done/`) enthalten
  je zwei datei-spezifische, am Repo nachprüfbare Aussagen: der Zeiger auf die
  jeweilige Ergebnis-Notiz (beide Dateien existieren), der eingetretene
  Re-Evaluierungs-Trigger, der ohne Supersede als Verfeinerung beantwortet ist
  (die genannte ADR existiert und trägt `Accepted`), und die beiden neuen
  Register-Einträge samt Zählerstand 3 an der Verkörperungs-Schwelle (in
  [observations.md](../plan/planning/observations.md) so belegt). Die übrigen
  vier Sätze sind zwischen beiden Dateien wortgleich (Schluss-Datum,
  Baseline-Form, Verzeichnis-Position) — die Satz-Schwelle würden sie allein
  nehmen, die **Aussage** kommt aus den datei-spezifischen Sätzen. Kein Befund:
  die Notizen sind nicht auf die Schwelle geschrieben. Nebenbefund: beide
  Dateien sind unter dem zurückgenommenen Filter ohnehin keine Kandidaten mehr —
  die Füllung ist freiwillig geblieben und nicht wieder zurückgenommen worden.
- geprüft, ohne Befund: **Bestandszahlen 95/110 zweifach reproduziert.**
  Verzeichnis-Listing (110 Markdown-Dateien, 95 `slice-*.md`, 15
  Wellen-Dokumente: 11 Ergebnis- + 4 Plan-Notizen) und Lauf mit absichtlich nie
  passendem `heading-pattern` (110 bzw. 95 `closure-note-missing`). Slice-Plan
  §3-Tabelle, ADR-§Konsequenzen, ADR-§Fitness Function und Config-Kommentar
  nennen 95; einzige Abweichung ist die ADR-Tabellenzeile (N-2).
- geprüft, ohne Befund: **Das zusammengezogene Akzeptanzkriterium ist gegen die
  Implementierung wahr.** Vier Fixtures decken beide Seiten des „effektiven"
  Filters ab (gesetzter Glob trifft / trifft nicht, geerbter Glob trifft /
  trifft nicht); in jedem Fall stimmen Befund und Exit-Code mit dem Kriterium
  überein. Die widerrufene Fassung ist rückstandslos entfernt — die Zeichenkette
  der alten Formulierung kommt in der Anforderung nicht mehr vor.
- geprüft, ohne Befund: **Kein Golden-Test, kein Fixture an der quotierten
  Meldung.** Die einzige Test-Zusicherung auf diesen Text prüft per
  `strings.Contains` auf das gesetzte Muster und überlebt die Quotierung; im
  ganzen Repo gibt es außerhalb der Produktionsstelle keine weitere Fundstelle
  des Meldungstextes. Auch der geänderte `--doctor`-Klartext ist nirgends
  gepinnt — beide Änderungen sind ungetestet (ihre Konsumenten-Wirkung ist N-1).
- geprüft, ohne Befund: **`AllReasons`-Lockstep.** Die beiden Deckungs-Tests
  (`reasonTexts` gegen `AllReasons`, `AllReasons` gegen die §4-Code-Spalte)
  greifen unverändert; die Heilung ändert nur einen Klartext-Wert, keinen Code
  und keine Tabellenzeilen-Kennung. Der §4-Parser bleibt an der um einen
  Klammerzusatz verlängerten Zeile stabil (`make test` grün).
- geprüft, ohne Befund: **Mutations-Echtheit, vier Rückbauten selbst
  nachvollzogen** (Sicherungskopie im Temp-Verzeichnis, Mutation, `make test`,
  Rückkopie — kein `git checkout`): Filter zurück auf `slice-glob` ⇒ zwei rote
  Tests; Default als Literal statt Verweis ⇒ ein roter Test; leerer Glob fällt
  still zurück ⇒ zwei rote Tests; Config-Rand an `closure.dir` gekoppelt ⇒ ein
  roter Test (der neue). Ein **siebter**, nicht gefangener Rückbau ist als N-6
  gemeldet.
- geprüft, ohne Befund: **Der neue Test ist nicht tautologisch.** Er lässt in
  allen drei Fällen `closure.dir` weg — genau die Konstellation, die keine
  bestehende Tabelle abdeckt — und wird von der Mutation, gegen die er gebaut
  ist, rot. Ohne Mutation grün.
- geprüft, ohne Befund: **YAML-`null` in drei Schreibweisen.** `glob:`,
  `glob: ~` und `glob: null` verhalten sich gleich (abwesend ⇒ geerbter Filter,
  Exit 0); die im Spec-Text gezogene Analogie zu `min-sentences` hält am Lauf
  (`min-sentences:` ⇒ Default, `min-sentences: 0` ⇒ Exit 2).
- geprüft, ohne Befund: **Nichts hing von der Weitung ab.** Die einzigen
  Fundstellen des zurückgenommenen `^#{1,3}`-Musters außerhalb der drei
  Erzähl-Dokumente sind Kommentare; `harness/README.md` beschreibt das Gate
  weiterhin als „je Slice im `done/`-Bestand" und wird durch die Rücknahme
  wieder wahr statt falsch. Der Go-Test, der `.d-check.closure.yml` dekodiert,
  läuft grün; der ADR-Index nennt keine Bestands-Zahl.
- geprüft, ohne Befund: **ADR-Immutabilität und Referenz-Richtung.**
  `make adr-check RANGE=39a3c6a..HEAD` grün — die Umkehrung von Entscheidung 5
  passiert in einer ADR mit Status `Proposed`, der Verfeinerungs-Zeiger in
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
  §Geschichte bleibt unberührt. Die neue §Geschichte-Zeile deklariert die
  Umkehrung ausdrücklich, statt sie still zu ersetzen.
- geprüft, ohne Befund: **Gates.** `make gates` grün (356 Dateien, 0 Befunde,
  alle acht Gates), `make verify-closure-notes` grün (327 Dateien, 0 Befunde).
- geprüft, ohne Befund: **SemVer-Einordnung Minor bleibt richtig.** Der Schlüssel
  ist additiv, die Kandidaten-Semantik ohne ihn unverändert, ältere Fassungen
  brechen am unbekannten Feld fail-closed ab. Die beiden Textänderungen sind
  kein Vertragsbruch (Meldungstexte sind nicht stabilitätsgarantiert) — sie
  gehören in die Release-Notiz, nicht in einen Major (N-1).
- geprüft, ohne Befund: **Formatierung.** Das Lint-Profil führt keinen
  Formatier-Linter (`.golangci.yml` enable-Liste ohne `gofmt`/`gci`); die
  Spalten-Ausrichtung der neuen Test-Map und die Leerzeile am Dateiende sind
  ohne Konventions-Anker kein Finding (Skill §Anti-Pattern).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 5 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Byte-Identitäts-Zusage von der Heilung
mitgerissen, ohne Historie · Bestandszahl in der Fitness Function nicht
reproduzierbar (**zweimal**: N-2, N-5) · Grund-Code-Klartext und §4-Zeile ohne
Text-Lockstep · verworfene Alternative mit widerlegbarer Begründung eingefroren ·
Config-Rand-Verhalten bei inerter Fähigkeit ungetestet · Vorlage bildet einen
Verweis-Default als Literal ab

## Verdikt

**Merge-blockierend:** ja — N-1. Die beiden Erst-MEDIUM sind **echt** geheilt
und am Lauf belegt: die eigene Weitung ist zurückgenommen (Konfiguration
semantisch identisch zum Vor-Slice-Stand), und das zusammengezogene
Akzeptanzkriterium hält gegen vier Fixtures in beiden Richtungen des effektiven
Filters. Auch der sechste Rückbau ist mit einem Test geschlossen, den ich selbst
rot gemacht habe. Blockierend ist, was die Heilung **neu** aufgemacht hat: sie
ändert zwei konsumenten-sichtbare Ausgabetexte auf Pfaden, die ohne den neuen
Schlüssel erreicht werden, und widerlegt damit die Byte-Identitäts-Zusage, die
in der Fitness Function und in Entscheidung 6 derselben ADR steht — ohne dass
eine der drei Historien es führt.

**Release-Empfehlung:** **noch nicht freigeben; die Reihenfolge zählt.** Vor dem
Statuswechsel der [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md)
auf `Accepted` — ab dann immutabel — gehören N-1 (Fitness Function und
Entscheidung 6 sagen etwas, das der Lauf widerlegt), N-2 (die Tabellenzeile
widerspricht der eigenen Messung zwei Absätze weiter oben, und die §Geschichte
behauptet die Korrektur als erledigt) und N-4 (eine eingefrorene Begründung, die
das eigene Handbuch widerlegt) geklärt. In dieselbe Runde gehört N-1s zweite
Hälfte: die Release-Notiz für den Minor muss die geänderte `--doctor`-Zeile und
die quotierte Nullmengen-Meldung nennen — `CHANGELOG.md` §Unreleased ist derzeit
leer, und das Repo hat für genau diese Klasse eine eigene Konsumenten-Zeile
etabliert. In die Release-Prep-Runde gehören der offene Erst-Befund F-7
(Handbuch führt `glob` weder im Schlüssel-Block noch im Aufgaben-Beispiel; die
Versions-Tabelle endet bei v0.53.0) und N-3 (§4-Zeile). N-5 gehört in dieselbe
Runde wie N-2, weil beide dieselbe angefasste Messung betreffen. N-6 und N-7
sind release-verträglich und passen in einen Folge-Schnitt — N-6 naheliegend
zusammen mit [slice-098](../plan/planning/open/slice-098-closure-note-placeholder.md),
N-7 an der nächsten Berührung der Vorlage.

**Minor (v0.54.0) ist die richtige Stufe** — additiver Schlüssel, unveränderte
Semantik ohne ihn, fail-closed bei Versions-Rückschritt. Die zwei Textänderungen
sind keine Major-Frage, sondern eine Release-Notiz-Frage.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die
Slice-Closure §7 und von dort in den Zähler des Beobachtungs-Registers — die
Klasse „Bestandszahl in der Fitness Function nicht reproduzierbar" fällt in
diesem Lauf zum zweiten und dritten Mal und stand schon im Erst-Report, ist
damit an der Verkörperungs-Schwelle. Dieser Report ist ein **Lauf-Beleg** (dieser
Diff, dieser Skill, dieses Modell, dieses Verdikt) und ersetzt keine Verifikation
— die DoD-/Spec-Konformität prüft der Verifier separat.
