# Review-Report: slice-217 — 2026-09-08

**Review-Art:** Code-Review einer Whitespace-Operation — geprüft wird der Diff
gegen Plan und Hard Rules, und vor allem: **trägt die Messung, auf der der
Anlass ruht, und ist die Zeichen-Klasse, die der Slice sich selbst gibt,
vollständig?**

**Gegenstand:** slice-217 · Commit-Range `HEAD~2..HEAD` (`c2cdb5e9` Plan +
Vorprüfungen · `6fd053c4` Beanspruchung · `34812641` die Änderung). Geändert:
`harness/README.md` (22 Zeilen), der Ruhe-Marker der Roadmap, der Slice-Plan.

**Skill:** `.harness/skills/reviewer.md` @ 1.16.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-09-08

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als Tag plus Pfad in Inline-Code
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>) statt als Link.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan slice-217, §1 Ziel und Abgrenzung, §2 DoD, §3 Plan (die
  Zeichen-Klasse), §6 Risiken, §8 Vorprüfungen
- `DC-FA-STRUCT-001` und `DC-FA-TGT-001` in voller Länge (Kandidaten-Menge,
  Bedingungen, Tabellen-Scoping, Akzeptanzkriterien)
- `SPEC-068`/`SPEC-070` in `spec/spezifikation.md` (Zell-Messung)
- `.d-check.yml` — die beiden Blöcke, die `harness/README.md` binden
  (`structure`-Regel auf §Sensors, `targets.doc-tables`)
- `MR-070` (Geltungsbereich-Feld, nicht die Index-Zeile), `MR-054`, `MR-013`
- `AGENTS.md` §3.1, §3.4, §3.7, §3.8, §5, §6
- Reviewer-Anker 8 (Reichweite einer Botschaft), 9 (Geltungsbereich eines
  Zitats), 15 (Scan-Achse), 17 (Proxy-Messung), 18 (Grenzen-Liste)
- die vier im Plan zitierten Register-Einträge, je `observation.md`, `state.md`
  und die Evidence-Liste
- die Baseline-Vorlage `.harness/baseline/v6.5.0/templates/harness/README.template.md`

**Eigene Läufe und Messungen** (echte Ausgaben, gekürzt auf die tragenden
Zeilen). Alle Messungen mit `awk` über die aus `git show` extrahierten Stände;
alle Läufe read-only gegen den Arbeitsbaum bzw. gegen eine Kopie im
Scratchpad. Der Arbeitsbaum war vor und nach diesem Review sauber.

- `make doc-check` — `d-check: 748 Datei(en) geprüft, 0 Befund(e)`, Exit 0
- `make gate-consistency` — `d-check: 748 Datei(en) geprüft, 0 Befund(e)`, Exit 0
- `make planning-check` — `d-check: 748 Datei(en) geprüft, 0 Befund(e)`, Exit 0
- `make baseline-verify` — `fetch-baseline-cache: verify ok (54 Dateien, vollständig)`, Exit 0
- `make mention-coverage` — `84 von 84 Artefakt(en) erwähnt`, Exit 0
  (berührt `harness/README.md` nicht — Ist-Menge ist der ADR-Index)
- **Positiv-Kontrolle reproduziert**, auf einer Kopie des Repos im Scratchpad,
  zweimal: ein Phantom-Token in die **entpaddete** §Guides-Zeile 62 eingesetzt
  ⇒ `harness/README.md:62  phantom-a  gate-phantom  dokumentiertes Target
  \`make phantom-a\` ohne Makefile-Regel`, Exit 1. Eine zusätzlich eingefügte
  Phantom-**Zeile** ⇒ derselbe Befund auf Z. 63. Die Kopie war danach
  byte-identisch zum Arbeitsbaum. Das Modul `targets` parst die entpaddete
  Tabelle also weiterhin — die Behauptung der Commit-Botschaft trifft zu.

**Eigene Nachmessung des Gegenstands** (gegen die in §3 deklarierte Form,
Vorstand `34812641^`):

| Größe | gemessen | Plan / Botschaft |
|---|---|---|
| Tabellenzeilen gesamt | 75 | 75 |
| Padding-Zeilen | 20 | 20 |
| Padding-Zeichen, **ganzer Lauf** | 4737 | 4737 |
| Padding-Zeichen, **Lauf minus ein Trennzeichen** | 4691 | — |
| davon letzte Spalte, ganzer Lauf / minus eins | 3412 / **3394** | 3394 |
| Rest, ganzer Lauf / minus eins | 1325 / **1297** | 1297 |
| überlange Trennzeilen | 2 | 2 |
| Zeilen mit gebrochener Ausrichtung | Z. 43, 44, 60, 62 (60 am inneren Pipe) | dieselben vier |
| Zeilen/Bytes vorher → nachher | 233/28873 → 233/23533 (−5340) | 233, −5340 |
| Vorlage: Tabellenzeilen / Padding | 39 / 0 | 39 / 0 |

---

## Findings

### F-1 — Die Zerlegung der tragenden Kennzahl geht nicht auf: 4737 ≠ 3394 + 1297, und die Gesamtzahl widerspricht der eigenen Klasse

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 (*„Vor einer Messung steht die Form ihres
  Gegenstands"*) · Reviewer-Anker 17 ·
  `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`
- `pfad`: slice-217 §1, Absatz *„Der Anlass ist gemessen, nicht ästhetisch"*
  (Zeilen 35–40 der Plandatei); wortgleich in Commit `c2cdb5e9`, erster Absatz
  der Botschaft
- `befund`: Der Satz *„Von **4737** Padding-Zeichen stehen **3394** in der
  jeweils letzten Spalte … Die restlichen **1297**"* mischt zwei
  Zählmethoden. **4737** ist die Summe der **ganzen** Läufe von ≥ 2 Leerzeichen
  vor einem Pipe (nachgemessen: 4737). **3394** und **1297** sind dieselben
  Läufe **je um ein Zeichen gekürzt** (nachgemessen: letzte Spalte 3394, Rest
  1297, Summe **4691**). Damit gilt 4737 − 3394 = 1343 ≠ 1297; die Differenz
  **46** ist exakt die Zahl der Läufe. Verschärfend: §3 nimmt das eine
  Leerzeichen ausdrücklich aus (*„**Nicht** dazu zählt das **eine** Leerzeichen
  nach und vor einem `|`"*) — nach der **eigenen** Klasse ist die richtige
  Gesamtzahl 4691, nicht 4737, und die beiden Teilzahlen sind die einzigen
  klassenkonformen Zahlen des Absatzes. Der Slice führt DoD (1) und §6 genau
  gegen diese Klasse ins Feld; es ist das dritte Mal in demselben Slice, dass
  die Messung des eigenen Gegenstands rutscht, und §6 hält fest: *„Ein drittes
  Mal wäre kein Zufall mehr"*.
- `verifizierbar`: ja — `awk` über `git show 34812641^:harness/README.md`,
  Läufe von ≥ 2 Leerzeichen vor einem `|` einmal voll und einmal um eins
  gekürzt summiert; die vier Zahlen 4737 / 4691 / 3394 / 1297 stehen in der
  Ausgabe.
- `klasse`: `zwei-zaehlmethoden-in-einem-satz`

### F-2 — §1 zitiert `rule-drawn-from-occasion-not-inventory` für eine Fehlerklasse, die der Eintrag nicht führt — mit umgekehrter Polarität

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 (*„Eine zitierte Quelle trägt nur, was in ihrem
  Geltungsbereich steht"*) · Reviewer-Anker 9 ·
  `BEO-ALL/citation-stretched-beyond-scope`
- `pfad`: slice-217 §1, Absatz *„Die Form-Frage ist an der Vorlage geklärt,
  nicht am Nachbarn"* (Zeilen 42–48): *„Ohne diese Prüfung wäre der Slice eine
  Angleichung an den lokalen **Bestand** gewesen — genau der Fehler, den
  `rule-drawn-from-occasion-not-inventory` führt."*
- `befund`: Der Register-Eintrag heißt *„Eine Aussage oder Regel wird aus dem
  **Anlass** gezogen statt aus dem **Bestand**"* und benennt in
  `docs/plan/planning/observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md`
  **drei** Ausprägungen: (a) eine Exklusivitäts-Aussage, (b) ein Kriterium, das
  eine bereits korrigierte Menge nachträglich begründet, (c) eine Regel, die
  nur für den Ort geschrieben wird, an dem sie wehtat. Alle **neun**
  Evidence-Dateien sind Fälle derselben Bauart — *„aus drei Anlässen gezogen,
  die Inventur nicht gemacht"* (slice-204, slice-209, slice-210 wörtlich
  nachgelesen). Der im Slice beschriebene Fehler — die **Ziel-Form** vom
  lokalen Nachbarn statt von der adoptierten Vorlage nehmen — ist keine der
  drei; er benutzt zudem *Bestand* in der **entgegengesetzten** Rolle: Beim
  Register ist der Bestand das, was man hätte zählen sollen, im Slice das, was
  man nicht hätte abschreiben sollen. Ein Eintrag für *„Form am Nachbarn
  abgelesen statt im Regelwerk nachgeschlagen"* existiert im Register nicht
  (37 Slugs unter `BEO-ALL` durchgesehen).
- `verifizierbar`: nein (Urteil) — Beleg ist der Wortlaut von `observation.md`
  und der neun Evidence-Dateien.
- `klasse`: `beobachtung-mit-umgekehrter-polaritaet-zitiert`

### F-3 — Die Commit-Botschaft schreibt der Baseline-Vorlage eine Trenner-Form zu, die sie nicht führt

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 (*„behauptet nicht mehr, als die Arbeit trägt"*) ·
  Reviewer-Anker 8 und 9
- `pfad`: Commit `34812641`, erster Absatz: *„§Source precedence und §Guides
  tragen jetzt dieselbe schlanke Form wie die drei uebrigen Tabellen der Datei
  **und wie die Baseline-Vorlage**: ein Leerzeichen beidseits jedes
  Zellinhalts, **Trenner `| --- |`**."* Korrespondierend slice-217 §1, Zeile 32.
- `befund`: Die Vorlage
  `.harness/baseline/v6.5.0/templates/harness/README.template.md` führt ihre
  **vier** Trennzeilen als `|---|---|---|` **ohne** Leerzeichen (gemessen); die
  gelandete Form `| --- | --- | --- |` stammt von den drei Nachbar-Tabellen
  derselben Datei, nicht von der Vorlage. Für die Datenzeilen trifft die
  Aussage zu, für den Trenner nicht — und der Trenner ist der eine Teil, den die
  Botschaft ausdrücklich nennt. Der Plan selbst hat die Vorlagen-Form korrekt
  notiert (§1: *„Trenner `|---|---|---|`"*); die Botschaft komprimiert das zu
  einer falschen Zuschreibung. Das trifft gerade die Stelle, die der Slice
  tragend macht: *„Die Form-Frage ist an der Vorlage geklärt, nicht am
  Nachbarn."* — beim Trenner war sie am Nachbarn geklärt.
- `verifizierbar`: ja — `grep -n '^|[ |:-]*$'` über die Vorlage liefert
  viermal `|---|---|---|`, über `harness/README.md` fünfmal `| --- | … |`.
- `klasse`: `form-der-falschen-autoritaet-zugeschrieben`

### F-4 — Die deklarierte Zeichen-Klasse ist einseitig, DoD (2) behauptet beidseitig; der Absatz „Nicht dazu zählt …" liest sich als vollständige Abgrenzung und nennt nur eine von fünf Auslassungen

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §5 (Grenzen-Liste / Zählmethode) · Reviewer-Anker 17
  und 18
- `pfad`: slice-217 §3, Zeilen 99–108 (die Klasse und ihr
  *„**Nicht** dazu zählt …"*-Absatz) gegen §2 DoD (2), Zeile 78
- `befund`: §3 definiert als Padding ausschließlich *„einen Lauf von zwei oder
  mehr Leerzeichen **unmittelbar vor** einem `|`"* plus die Bindestrich-Strecke
  der Trennzeile. DoD (2) verlangt dagegen *„genau ein Leerzeichen
  **beidseits** jedes Zellinhalts"*. Die Klasse misst nur die eine Seite; **vier
  weitere Formen überflüssigen Whitespace in Markdown-Tabellen fallen nicht
  darunter** und werden im Abgrenzungs-Absatz auch nicht als Auslassung
  benannt: Links-Padding hinter dem Pipe (`|   Zelle`, die Regelform bei
  rechtsbündigen Spalten `|---:|`), Whitespace nach dem letzten Pipe,
  Tabulatoren (die Klasse sagt *Leerzeichen*) und das geschützte Leerzeichen
  U+00A0. Ungeklärt bleibt außerdem, was eine *Tabellenzeile* ist — eine um bis
  zu drei Stellen eingerückte Zeile rendert weiterhin als Tabelle. **Das
  Ergebnis ist trotzdem richtig:** Ich habe alle fünf Formen selbst gemessen,
  in beiden Ständen und in beiden Tabellen null Treffer (Links-Padding
  zeichengenau geprüft, nicht per `grep '|  '`). Belegt ist der Zustand damit
  durch diese Messung, nicht durch die des Slice — DoD (2) behauptet eine
  Eigenschaft, die die eigene Klasse nicht prüfen kann. Zur Einordnung: unter
  derselben Klasse trägt `AGENTS.md` 38 Padding-Zeilen / 3179 Zeichen und
  ebenfalls kein Links-Padding, die Lücke bisse dort also auch nicht.
- `verifizierbar`: ja — `awk` über beide Stände, Läufe von ≥ 2 Leerzeichen
  **nach** einem Pipe mit folgendem Nicht-Leerzeichen; dazu
  `grep -c '[[:space:]]$'`, ein Tab-`grep` und ein `grep` auf `\302\240`.
- `klasse`: `deklarierte-klasse-deckt-nur-eine-seite`

### F-5 — DoD (3) führt zwei „nachweislich intakte" Kanten; die `structure`-Kante zeigt auf einen Abschnitt, den der Slice nicht anfasst

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.8 (*„Ein Modul verspricht nur über das, was es
  scannt"*) · Reviewer-Anker 15 · `DC-FA-STRUCT-001`
- `pfad`: slice-217 §2 DoD (3), Zeilen 80–85; §Bezug, Zeilen 13–18; §6 drittes
  Risiko, Zeilen 177–186. Gegenstand: `.d-check.yml` Zeilen 650–671
- `befund`: Die `cell-min-chars`-Regel bindet `harness/README.md` an
  `section: "## Sensors (Feedback-Gates)"` — eine der **drei unveränderten**
  Tabellen. Für die zwei geänderten Tabellen ist die Regel nicht zuständig; ein
  grüner `make doc-check` ist auf dieser Kante deshalb keine Messung, sondern
  eine Tautologie. Bewiesen wurde die eine Kante, die den Diff berührt
  (`targets`), und zwar sauber per Positiv-Kontrolle — die andere nicht, weil
  es dort nichts zu beweisen gab. DoD (3) und §Bezug führen beide gleichrangig.
  Hinzu kommt: §6 erklärt die Frage *„misst `cell-min-chars` gegen den
  gepaddeten oder den getrimmten Zellinhalt"* für **nicht gemessen** und den
  Slice für zuständig, sie *„durch Fahren, nicht durch Lesen"* zu beantworten —
  `spec/spezifikation.md` beantwortet sie bei `SPEC-068` ausdrücklich
  (*„gemessen wird die getrimmte Zelle, wie sie dasteht"*). Fahren ist der
  stärkere Beleg, aber nicht der einzige verfügbare, und gefahren wurde die
  Kante ohnehin nicht dort, wo sich etwas ändern konnte. Der zweite Halbsatz
  desselben Risikos (*„für §Guides mit seinem einen `make`-Target ist es die
  offene Frage"*) hängt die `cell-min-chars`-Frage an einen Abschnitt, auf dem
  keine `cell-min-chars`-Regel liegt.
- `verifizierbar`: ja — `grep -n 'harness/README' .d-check.yml` liefert genau
  zwei Bindungen (Zeile 650 mit `section: "## Sensors …"`, Zeile 810
  `doc-tables`); die geänderten Zeilen 37–62 liegen außerhalb von §Sensors.
- `klasse`: `kante-ausserhalb-der-aenderung-als-beleg-gefuehrt`

### F-6 — §8 nennt drei einschlägige Register-Einträge; §1 stützt die zentrale Entwurfsentscheidung auf einen vierten

- `kategorie`: LOW
- `quelle`: `v6.5.0` · `regelwerk/modul-05-planning-harness.md`
  §Zwei Schritte vor der Modus-Begründung (der Sichtungs-Schritt) · `MR-054`
- `pfad`: slice-217 §8, Zeile 222 (*„**Drei** Einträge sind einschlägig"*)
  gegen §1, Zeilen 45–48
- `befund`: Der Sichtungs-Block zählt `zaehlmethode-misst-proxy-statt-gegenstand`,
  `mechanical-id-rewrite-misses-frozen-classes` und
  `semantic-change-body-only-edges-stale` auf und schließt die Menge mit
  *„Drei"*. §1 begründet die Prüfung gegen die Baseline-Vorlage — die
  Entscheidung, die den Slice vom bloßen Angleichen an den Nachbarn trennt —
  mit einem vierten Eintrag, `rule-drawn-from-occasion-not-inventory`, der in
  §8 nicht vorkommt. Ein Eintrag, der den Entwurf nachweislich gesteuert hat,
  fehlt damit in genau dem Block, der die Steuerung dokumentieren soll. (Ob er
  dort **hingehört**, ist wegen F-2 offen — beides gehört zusammen entschieden.)
- `verifizierbar`: nein (Urteil) — beide Stellen stehen in derselben Datei.
- `klasse`: `sichtungs-block-nennt-weniger-als-der-plan-benutzt`

### F-7 — `AGENTS.md` trägt unter derselben Klasse 38 Padding-Zeilen; Abgrenzung (2) ist eingehalten, ihre Begründung ist ab jetzt überholt

- `kategorie`: INFO
- `quelle`: Maintainability · slice-217 §1 Abgrenzung (2)
- `pfad`: slice-217 §1, Zeilen 56–59
- `befund`: Abgrenzung (2) nimmt `AGENTS.md` mit dem Grund heraus, *„ob sie
  dasselbe Bild zeigt, ist **nicht gemessen** und wird hier auch nicht
  gemessen"*. Die Abgrenzung ist eingehalten — der Diff fasst `AGENTS.md` nicht
  an. Für den Fall, dass jemand den Folge-Vorgang schneidet, ist die Zahl
  jetzt da: 43 Tabellenzeilen, davon **38** mit Padding, **3179** Padding-Zeichen
  (ganzer Lauf), **eine** überlange Trennzeile, kein Links-Padding, kein
  Zeilenende-Whitespace. Damit ist die Begründung *„nicht gemessen"* nicht mehr
  die Lage, sondern nur noch die Lage zum Zeitpunkt des Plans.
- `verifizierbar`: ja — dieselbe `awk`-Messung über `AGENTS.md`.
- `klasse`: `abgrenzungs-begruendung-durch-den-review-ueberholt`

---

## Negativbefunde

Eine Zeile je betrachtetem Bereich — sonst ist „keine Findings" nicht von
„nicht geprüft" unterscheidbar.

- **geprüft, ohne Befund: der Diff fällt vollständig in die deklarierte
  Klasse.** `diff -w` nach Normalisierung der Bindestrich-Läufe ist gegen den
  Vorstand **leer** (Exit 0). Ohne die Bindestrich-Normalisierung bleiben genau
  die zwei Trennzeilen stehen — die Commit-Botschaft benennt diesen Unterschied
  von sich aus und weist die schwächere Probe zurück; das trifft zu und ist die
  seltenere Richtung von Anker 8.
- **geprüft, ohne Befund: die drei übrigen Tabellen sind unverändert.**
  Geänderte Zeilen sind exakt 37–47 und 51–59, 61, 62 (22 Stück, wie die
  Statistik sagt). Zeile 60 (die lange Baseline-Zeile) ist **nicht** angefasst,
  weil sie kein Padding trug — korrekt. `cmp` über Z. 1–36, Z. 48–50, Z. 60 und
  Z. 63–233 (171 Zeilen) meldet jeweils Identität. Die Botschaft sagt
  *„cmp ueber 170 Zeilen"*; die Aussage *„byte-identisch"* trifft unabhängig
  von der Grenzwahl zu.
- **geprüft, ohne Befund: die Datei rendert weiterhin als Tabelle.** Alle elf
  Zeilen von §Source precedence führen 4 Pipes (3 Zellen), alle zwölf von
  §Guides 3 Pipes (2 Zellen); die drei übrigen Tabellen unverändert 4/4/5
  Pipes. Keine Zeile hat eine abweichende Spaltenzahl, kein Pipe steht in einer
  Inline-Code-Spanne (zeichengenau geprüft, nicht per `grep`), keine escapte
  Pipe `\|` kommt vor. Ein HTML-Marker ging nicht verloren: in beiden Tabellen
  stand vor **und** nach der Änderung keiner.
- **geprüft, ohne Befund: das Modul `targets` liest die entpaddete Tabelle
  weiterhin.** Die Positiv-Kontrolle der Botschaft ist auf einer Kopie
  reproduziert, mit demselben Grund-Code und derselben Zeile (`62`). Der
  Erkennungs-Vertrag von `DC-FA-TGT-001` hängt am `|`-**Präfix** der Zeile —
  das hat die Operation nicht angetastet, und genau deshalb konnte sie hier
  nichts still abschalten. §Guides trägt danach weiterhin **ein**
  `make`-Target (`make verify-closure-notes`), §Source precedence keines; beide
  Zahlen decken sich mit dem Plan.
- **geprüft, ohne Befund: kein Stilles-Grün-Pfad (Reviewer-Frage 1).** Keine
  der fünf `harness/README.md`-berührenden Prüfungen wird durch die Änderung
  blind: `structure` scannt einen anderen Abschnitt (F-5), `targets` ist per
  Positiv-Kontrolle nachweislich wach, `links`/`anchors`/`ids`/`matrix`/
  `codepaths`/`versions`/`citations`/`spans`/`diagrams` laufen in `doc-check`
  über dieselbe Datei mit 0 Befunden. Repo-weit binden nur **zwei**
  Konfigurations-Einträge diese Datei — die Kanten-Suche des Plans ist
  vollständig, das habe ich selbst nachgesucht.
- **geprüft, ohne Befund: alle Mengenangaben des Plans.** Die §3-Tabelle des
  Ausgangsstands (11/10, 12/10, 25/0, 22/0, 5/0; Trennzeilen 226 und 455
  Zeichen) reproduziert exakt, ebenso 75 Tabellenzeilen, 233 Zeilen vorher wie
  nachher, 5340 Bytes Differenz, die vier ausrichtungsbrechenden Zeilen
  einschließlich des inneren Pipes in Z. 60 (Pipe-Position 83 statt 122),
  die 39 padding-freien Tabellenzeilen der Baseline-Vorlage, die 38
  Register-Verzeichnisse über beide Kürzel, die Zähler 4× / 3× / 18× / 15×
  (als Zahl der Evidence-Dateien), der eine offene `BEO-HARN`-Eintrag zu
  `--check-latest` und die 117 Commits auf `harness/README.md` zum Zeitpunkt
  des Plan-Commits (heute 118). **Einzige Ausnahme ist F-1.**
- **geprüft, ohne Befund: `MR-070` ist korrekt als nicht einschlägig
  eingeordnet.** Das Feld `Geltungsbereich` in
  `harness/conventions/MR-070-frozen-klassen-vor-mechanischer-ersetzung.md`
  lautet *„jede mechanische Ersetzung über mehr als eine Datei"* und nimmt
  *„die Änderung einer einzelnen Datei"* ausdrücklich aus. Der Slice fasst eine
  Datei an; §8 hat das Feld gelesen, nicht den Titel. Die Wiedergabe im Plan
  ist grammatisch angepasst, inhaltlich exakt. Auch die dort angeschlossene
  Überlegung trägt: Die eingefrorene Klasse ist hier die vendorte Vorlage, und
  Abgrenzung (4) nimmt sie aus — nachgewiesen durch grünes
  `make baseline-verify`.
- **geprüft, ohne Befund: die zwei `d-check:cite`-Direktiven nach `MR-054`.**
  Beide zeigen auf die **vorschreibende** Zeile (`modul-05` Z. 268–269 bzw.
  Z. 274), nicht auf eine Nebenregel; die zitierten Wortlaute sind identisch mit
  der Quelle, `citations` läuft im inneren Loop grün. Die Form ist byte-gleich
  mit der von slice-216 — Haus-Form, keine Abweichung. Der dritte Block
  (Nachtlauf) trägt korrekt keine Direktive.
- **geprüft, ohne Befund: die Abgrenzung wurde nicht ausgeweitet.** Alle vier
  Punkte halten: kein Zeichen außerhalb der Whitespace-Klasse (Gegenprobe
  oben), keine andere Datei (die drei Commits fassen zusammen nur
  `harness/README.md`, den Slice-Plan und den Roadmap-Ruhe-Marker an), kein
  neuer Sensor (`.d-check.yml` und `Makefile` unverändert), keine Baseline
  berührt.
- **geprüft, ohne Befund: `MR-013` beim Beanspruchungs-Move.** `6fd053c4` ist
  ein Rename mit Score 100 % und trägt zusätzlich den Roadmap-Flip; kein
  Dokument nannte den `open/`-Pfad. `make planning-check` grün.
- **geprüft, ohne Befund: `AGENTS.md` §3.7 (Kommentare und Zustandsfelder).**
  Der Diff enthält keinen Kommentar und kein Zustandsfeld; die einzige
  Zustands-Aussage der drei Commits ist der entfernte Ruhe-Marker, und der ist
  Zustand, nicht Chronik.
- **geprüft, ohne Befund: `AGENTS.md` §3.4 und Referenz-Richtung.** Der Slice
  berührt keine Spec-Stelle (Kopf-Feld `—`, zutreffend); `harness/README.md`
  ist kein Spec-Stratum; kein `d-check:status-provenance`-Marker im Diff;
  `matrix` meldet nichts.
- **geprüft, ohne Befund: `AGENTS.md` §3.1 in diesem Review-Lauf.** Kein
  Host-Go, kein Host-Paketmanager, kein Host-Skript-Interpreter — nur `bash`,
  `awk`, `sed`, `grep`, `find`, `cmp`, `diff`, `git`, `docker` und
  `make baseline-verify`. Der Arbeitsbaum war vor und nach dem Review sauber;
  alle Schreibversuche liefen in den Scratchpad.
- **geprüft, ohne Befund: die Slice-Größe.** Ein Liefer-Punkt in einer Datei,
  22 Zeilen Diff, in einer Sitzung prüfbar — weit unter der
  Ein-Sitzungs-Grenze; `MR-066` ist nicht berührt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:**
`zwei-zaehlmethoden-in-einem-satz` ·
`beobachtung-mit-umgekehrter-polaritaet-zitiert` ·
`form-der-falschen-autoritaet-zugeschrieben` ·
`deklarierte-klasse-deckt-nur-eine-seite` ·
`kante-ausserhalb-der-aenderung-als-beleg-gefuehrt` ·
`sichtungs-block-nennt-weniger-als-der-plan-benutzt` ·
`abgrenzungs-begruendung-durch-den-review-ueberholt`

## Verdikt

**Merge-blockierend:** ja — drei MEDIUM.

**Das Arbeitsergebnis selbst ist sauber, und das ist der eigentliche Befund
dieses Laufs.** Der Diff fällt restlos in die deklarierte Whitespace-Klasse,
die drei übrigen Tabellen sind byte-identisch, die Datei rendert weiterhin
korrekt, und die einzige Konfigurations-Kante, die die Änderung berühren
konnte, ist per Positiv-Kontrolle nachweislich wach — reproduziert, nicht
geglaubt. Diese Positiv-Kontrolle ist die stärkste Einzelleistung des Slice:
Sie beantwortet die Frage, die ein grüner Lauf offenlässt.

**Was nicht trägt, ist durchweg die Prosa um die Messung herum.** Ein Slice,
dessen ganzer Gegenstand eine Zeichen-Klasse ist und der sich DoD (1) genau
dafür gibt, hat seine tragende Kennzahl in zwei Methoden zerlegt, die nicht
zusammenpassen (F-1), hat die Klasse einseitig definiert und zweiseitig
behauptet (F-4), hat die Autorität, an der er die Form geklärt haben will, in
der Commit-Botschaft für eine Trenner-Form in Anspruch genommen, die sie nicht
führt (F-3), und hat die Register-Beobachtung, die ihm den Weg zur Vorlage
wies, mit umgekehrter Polarität zitiert (F-2). Keiner dieser Punkte ändert eine
Datei; alle vier ändern, was der Slice über sich selbst aussagt — und das ist
das Artefakt, das nach `done/` wandert und dort gelesen wird.

**Übergabe:** Findings an den Implementer (Rückkante Review → Plan; F-1, F-2
und F-4 betreffen die Messung von DoD (1) und den Text des Plans, F-3 die
Commit-Botschaft, die nur über einen Folge-Commit oder die Closure-Notiz
korrigierbar ist). F-6 gehört zusammen mit F-2 entschieden — beide betreffen
denselben Register-Eintrag. Die **Finding-Klassen** gehen zusätzlich in die
Slice-Closure §7 und von dort in den Zähler;
`zaehlmethode-misst-proxy-statt-gegenstand`,
`citation-stretched-beyond-scope`, `commit-message-overclaims-work` und
`grenzen-liste-wird-als-vollstaendig-gelesen` sind in diesem Lauf jeweils
erneut belegt. Dieser Report ersetzt keine Verifikation — DoD-Konformität und
der volle Gate-Stand sind Sache des Verifiers.
