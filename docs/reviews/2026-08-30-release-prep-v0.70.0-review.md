# Review Release-Prep v0.70.0 — `11d8a60`

**Review-Art:** Release-Prep (gegen [`docs/user/releasing.md`](../user/releasing.md) §Release-Prep) — **kein** Slice-Review
**Gegenstand:** `11d8a60` (fünf Dateien: [`version.md`](../../version.md), [`CHANGELOG.md`](../../CHANGELOG.md), [`README.de.md`](../../README.de.md), [`README.md`](../../README.md), [`docs/user/benutzerhandbuch.md`](../user/benutzerhandbuch.md)); der Inhalt von slice-178 und slice-172 ist **nicht** Gegenstand (je einmal reviewt und verifiziert)
**Modell:** claude-opus-5[1m] · **Datum:** 2026-08-30
**Eingangs-Kontext:** `docs/user/releasing.md` §Release-Prep samt seinen eigenen Warnungen, `AGENTS.md` §3.1/§4/§5, `spec/lastenheft.md` `DC-FA-STRUCT-001` (0.79.0), `spec/spezifikation.md` `SPEC-078`, ADR-0074, MR-056, die beiden Closure-Notizen in `done/`

---

## Eigener Lauf

| Lauf | Ausgabe |
|---|---|
| `make ci` | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` · `[ci] gates + image-test green` · Exit 0 |
| `make build` | Image-ID `sha256:fb5c3b9070…` — identisch mit dem in der slice-178-Closure genannten Hash |
| Anker-Probe gegen eine Kopie von `version.md` | s. B1 |
| Pin-Zählung (`grep -ro`) über beide READMEs + Handbuch | 26 `ghcr`-Pins, 2 nackte Tags — s. B2 |
| Bedingungs-Zählung Produkt/Lastenheft gegen die READMEs | 10 — s. B5 |
| Closure-Profil-Probe ohne/mit Bestands-Ausnahme | 144 Befunde in 37 Dateien / 0 — s. B6 |
| Byte-Identitäts-Probe `d-check:vor178` gegen `d-check:latest` | `diff` leer, 178 Befunde beidseits — s. B7 |
| 5 Probe-Bäume à 3–7 Dateien gegen `d-check:latest` (vier Grenzen, Lexik, Backtick-Ausfall, negativer Wert) | s. B8–B11 |
| Grund-Code-Mengenvergleich Produkt ↔ Handbuch §6 | 16 gegen 14 — s. M1 |
| Handbuch-`structure`-Beispiel gegen `--config` | `562 Datei(en) geprüft, 0 Befund(e)` · Exit 0 — s. B12 |

`git status --short` zeigt am Ende nur diesen Report. `make fullbuild`, `make record-gates` und `make image-scan` bewusst **nicht** gefahren (Nachweis-Datei bzw. Netz).

---

## Urteil

**TAGGEN — mit einer Empfehlung.** 0 HIGH · 4 MEDIUM · 3 LOW · 4 INFO.

Die mechanische Hälfte der Prep ist **vollständig und gemessen richtig**: der Anker ist
gewandert und bricht genau den Pin, den er brechen soll; alle 26 `ghcr`-Pins und **beide**
nackten Tags sind gezogen; alle sieben Datumsstempel stimmen und stammen vom Commit-Tag;
die §11-Zeile steht chronologisch unten und widerspricht dem CHANGELOG nirgends; alle vier
behaupteten Grenzen sind gegen das laufende Binary reproduziert und **keine** ist schärfer
formuliert als gemessen; die drei `operations.md`-Aussagen stimmen; `make ci` ist grün.
Der Fund, wegen dem die Prep mehr als Pins war — die fehlende zehnte Bedingung in beiden
READMEs — ist korrekt gefunden, korrekt gezählt und in beiden Sprachfassungen synchron
nachgetragen.

Die vier MEDIUM sind **alle** derselben Klasse: eine Aufzählung oder Zahl in einer
Deklarations-Fläche, die kein Gate hält. Drei davon hätte dieselbe Prep-Runde finden
können — zwei stehen in Dateien, die sie angefasst hat. Keine davon macht das Release rot
und keine ist nach dem Tag nur mit einem neuen Release korrigierbar (die Doku liegt im
Repo, nicht im Image; für den Handbuch-Nachtrag gibt es mit dem Digest-Pin-Folge-Commit
ohnehin einen Träger). **Empfehlung: M1 und M3 vor dem Tag ziehen** — beides sind
Ein-Zeilen-Änderungen, und M3 sitzt im eingefrorenen Release-Protokoll.

---

## MEDIUM

### M1 — Die Grund-Code-Spalte des Handbuchs kennt die zehnte Bedingung nicht (und die neunte auch nicht)

**Fundstelle:** `docs/user/benutzerhandbuch.md:2239` (§6-Modul-Tabelle, letzte Spalte der `structure`-Zeile)

**Behauptung:** Die Prep hat die neue Fähigkeit in beiden READMEs nachgetragen und das
Handbuch auf `1.64` gehoben. Die §6-Zeile beschreibt `max-open-tasks` in ihrer
Beschreibungsspalte ausdrücklich als *„die zehnte"* und nennt `section-tasks-open` dort im
Fließtext.

**Messung:**

```text
$ grep -o "section-[a-z-]*" internal/hexagon/core/model/finding.go | sort -u > produkt.txt
$ sed -n '2239p' docs/user/benutzerhandbuch.md | awk -F'|' '{print $(NF-1)}' \
    | grep -o "section-[a-z-]*" | sort -u > handbuch.txt
$ echo "Produkt: $(wc -l < produkt.txt) · Handbuch-Spalte: $(wc -l < handbuch.txt)"
Produkt: 16 · Handbuch-Spalte: 14
$ comm -23 produkt.txt handbuch.txt
section-exempt-mismatch
section-tasks-open
```

Gegenprobe über das **ganze** Handbuch: von 48 Grund-Codes des Produkts fehlt keiner
irgendwo im Dokument — die Lücke ist also **genau** diese eine Spalte, nicht eine breite
Doku-Drift.

`section-exempt-mismatch` ist der Code aus `v0.69.0`; er wurde in der **vorigen** Prep
schon übersehen. Das ist die Wiederholung, nicht der Einzelfall — und die Spalte ist die
Stelle, an der ein Konsument nachschlägt, welche Befunde ein Modul werfen kann.

**Vorschlag:** die beiden Codes an die Spalte anhängen (`…, section-column-missing,
section-exempt-mismatch, section-tasks-open`), im selben Commit wie der Tag oder im
Digest-Pin-Folge-Commit. Ohne Handbuch-Versions-Bump — die Spalte holt nach, was die
Zeile `1.63`/`1.64` bereits versprochen hat.

### M2 — Beide READMEs behaupten seit ~zwei Monaten „acht Module" im Dogfooding-Vollausbau

**Fundstelle:** `README.de.md:251` · `README.md:247`

**Behauptung:** *„Dogfooding: d-check validiert die eigene Doku bei jedem Gate-Lauf — mit der
Selbstkonfiguration im **Vollausbau (acht Module** inkl. Referenzmatrix, Span-Artefakten,
Host-Pfad-Hygiene und Versions-Pin-Konsistenz)."* — englisch wortgleich („eight modules").

**Messung:**

```text
$ grep "^modules:" .d-check.yml
modules: [links, anchors, ids, matrix, codepaths, spans, hostpaths, versions, structure, diagrams, citations]
$ grep "^modules:" .d-check.yml | tr ',' '\n' | wc -l
11

$ git log -1 --format="%h %ad" --date=short 720b675
720b675 2026-07-01
$ git show 720b675:.d-check.yml | grep "^modules:"
modules: [links, anchors, ids, matrix, codepaths, spans, hostpaths, versions]
```

Die Zahl war am 2026-07-01 richtig und ist es seit `structure`, `diagrams` und `citations`
nicht mehr. Auch die Klammer-Enumeration nennt keines der drei. Beide Sprachfassungen sind
**synchron falsch** — die DE↔EN-Disziplin hat gehalten, die Currency nicht.

Das ist genau die Klasse, die [`releasing.md`](../user/releasing.md) selbst benennt
(*„die Modul-Liste blieb so von v0.25 bis v0.37 still bei acht Modulen stehen"*) — nur an
einer anderen Datei als der dort genannten.

**Vorschlag:** Zahl auf `elf` / `eleven` und die Enumeration um Struktur-Invarianten,
Diagramm-Kennungen und Zitat-Verifikation ergänzen; DE zuerst.

### M3 — Der CHANGELOG publiziert eine Zahl, die die eigene Verifikation als nicht reproduzierbar protokolliert hat

**Fundstelle:** `CHANGELOG.md:45` — *„gemessen gegen ein aus dem Vor-Commit gebautes Image: **276 Befunde beidseits**, `diff` leer"*

**Behauptung:** Die tragende Aussage ist *„ohne den Schlüssel ist der Befundsatz
byte-identisch"*. Die Zahl `276` soll sie belegen.

**Messung:** Die tragende Aussage **hält** — mit eigenem Vorgänger-Image gegengeprüft:

```text
$ docker run --rm --network none -v "$PWD:/repo:ro" -v "<probe>.yml:/repo/.d-check.closure.yml:ro" \
    d-check:vor178  --config .d-check.closure.yml --enable structure > alt.txt
$ docker run --rm --network none -v "$PWD:/repo:ro" -v "<probe>.yml:/repo/.d-check.closure.yml:ro" \
    d-check:latest  --config .d-check.closure.yml --enable structure > neu.txt
$ wc -l alt.txt neu.txt
178 alt.txt
178 neu.txt
$ diff alt.txt neu.txt && echo "DIFF LEER"
DIFF LEER
```

Die **Zahl** dagegen ist aus dem publizierten Text nicht nachvollziehbar: der Eintrag nennt
weder Korpus noch Glob noch Selektor. Drei plausible Varianten ergeben `223` / `227` / `364`,
keine `276`. Und die Closure-Notiz von slice-178 protokolliert denselben Befund bereits als
Korrektur an sich selbst:

> *„«276 Befunde» nannte weder Glob noch Selektor; die Verifikation maß mit drei plausiblen
> Varianten 164/331/161 und konnte sie ohne den Korpus nicht reproduzieren"*
> ([slice-178](../plan/planning/done/slice-178-offene-tasks-roh.md), §Lerneintrag)

Die Prep hat die Zahl aus der Closure-Notiz ins **eingefrorene** Release-Protokoll
übernommen — samt der Eigenschaft, wegen der sie dort schon beanstandet war. `AGENTS.md` §5
macht daraus eine Hard Rule (*„ihr Schluss reicht nicht weiter als die gemessene Menge"*).

**Vorschlag:** die Zahl streichen (*„gemessen gegen ein aus dem Vor-Commit gebautes Image:
`diff` leer"*) oder den Korpus dazuschreiben. Vor dem Tag — danach ist es eine Korrektur am
veröffentlichten Protokoll.

### M4 — Die Prep-Checkliste zeigt auf eine README-Zeile, die es nicht gibt

**Fundstelle:** `docs/user/releasing.md:69-70`

**Behauptung:** *„(a) die **Status-Zeile** („alle N Regelmodule (…)" — Zahl *und* Enumeration
*und* das „zuletzt das Modul X"-Fragment)"*.

**Messung:**

```text
$ grep -rn "Regelmodule (" README.md README.de.md
$ grep -rn "zuletzt das Modul" README.md README.de.md
$ grep -rn "zuletzt das Modul" --include="*.md" . | grep -v done/ | grep -v reviews/
docs/user/releasing.md:70:     „zuletzt das Modul X"-Fragment) und (b) die **Modul-Liste** unter
```

Beide Muster kommen in **keinem** der READMEs vor; der einzige Treffer im Repo ist die
Checkliste, die sie beschreibt. Wer Punkt (a) buchstäblich abarbeitet, findet nichts und
hakt ihn ab — während die tatsächlich zahl-tragende Zeile (M2) ungenannt bleibt. Das ist
kein Schreibfehler, sondern der **Grund**, warum M2 zwei Monate überlebt hat: ein Wächter,
der nie fangen konnte.

**Vorschlag:** Punkt (a) auf die real existierenden Zahl-Träger umschreiben — die
Dogfooding-Zeile („N Module") und die Modul-Bullet-Liste unter §Was ist d-check — und den
Prüfsatz danebenstellen: `grep "^modules:" .d-check.yml | tr ',' '\n' | wc -l` gegen die
Zahl im Text. Eigener Folge-Slice, nicht Teil dieses Tags.

---

## LOW

### L1 — `dafür fence-treu` behauptet einen Gegensatz zu `max-tasks`, den es nicht gibt (nur DE)

**Fundstelle:** `README.de.md:136-137` — *„anders als `max-tasks` blind gegen die absatzweise
Inline-Code-Paarung, **dafür fence-treu**"*. Englisch: *„unlike `max-tasks` it is immune to
the paragraph-wide inline-code pairing, **while staying fence-true**"*.

Das DE-*„dafür"* liest sich als Kompensation („dafür ist sie wenigstens fence-treu") und
legt nahe, `max-tasks` sei es nicht. Gemessen ist `max-tasks` ebenfalls fence-treu:

```text
Probe a-fence.md, Abschnitt mit einem freien und einem Task-Item IM Fence:
  max-open-tasks: 0  ⇒  a-fence.md:5  section-tasks-open      (nur das freie)
  max-tasks: 0       ⇒  a-fence.md:3  "Abschnitt trägt 1 Task-Items, erlaubt sind 0"
```

`1`, nicht `2` — der Fence bleibt für **beide** außen vor. Zusätzlich kollidiert *„blind
gegen"* mit dem Wortgebrauch derselben Release-Notiz, in der *„blind"* zwei Absätze weiter
den **Defekt** bezeichnet (*„ein vergessener Schluss-Fence macht die Bedingung blind"*).
Die EN-Fassung hat beide Probleme nicht — die kanonische Fassung ist hier die schwächere.

**Vorschlag:** *„…, die ein überzähliger Backtick im Absatz nicht abschaltet; der Fence
bleibt wie bei `max-tasks` außen vor"*.

### L2 — Der CHANGELOG nennt den Abschnitts-Selektor enger, als er ist

**Fundstelle:** `CHANGELOG.md:52` — *„hält jetzt `max-open-tasks: 0` über den
`## N. Definition of Done`-Abschnitt"*.

Der Selektor im Profil ist `^#{2,3} [0-9]+\. Definition of Done`
([`.d-check.closure.yml`](../../.d-check.closure.yml):216) und trifft H2 **und** H3. Die
Notiz ist damit nicht falsch, aber enger als die Regel — und wer sie als Vorlage kopiert,
schreibt eine Regel, die einen Teil des Bestands nicht sieht.

**Vorschlag:** `## N.`/`### N.` oder schlicht *„den Definition-of-Done-Abschnitt"*.

### L3 — Der `Changed`-Eintrag trägt als einziger keinen auflösbaren Anker

**Fundstelle:** `CHANGELOG.md:48-55`

Jeder andere Eintrag der Datei nennt seine Quelle (ADR bzw. `DC-*`-Anforderung) als Link.
Der slice-172-Eintrag beschreibt eine Regel, deren Begründung in MR-056 und ADR-0073/ADR-0074
liegt, und nennt keine davon. Für eine Änderung, die ausdrücklich *„die Regel bisher als
verworfen führte"* sagt, ist der Zeiger auf die ablösende MR die halbe Aussage.

**Vorschlag:** `([MR-056](harness/conventions/MR-056-dod-haken-waechter.md))` anhängen.

---

## INFO

### I1 — `releasing.md` beschreibt einen `[Unreleased]`-Schnitt, den dieses Repo nicht fährt

`docs/user/releasing.md:14-15` und `:36`: *„Vor dem Tag wird dort der `[Unreleased]`-Stand
unter die neue Version geschnitten."* Gemessen am Stand vor der Prep:

```text
$ git show 11d8a60^:CHANGELOG.md | head -7 | tail -1
## [0.69.0] — 2026-08-30
$ grep -c "Unreleased" CHANGELOG.md
0
```

Es gab keinen `[Unreleased]`-Abschnitt und gibt keinen. Die Prep hat richtig gehandelt
(Überschrift direkt angelegt) — die Prozessbeschreibung hinkt der Praxis nach. Kein Defekt
dieses Commits.

### I2 — Der Digest-Pin des Handbuchs zeigt weiter auf `v0.69.0`

`docs/user/benutzerhandbuch.md:99` trägt `sha256:dca9b350…`, den Digest von `v0.69.0`,
während Kopf und alle Pull-Beispiele `v0.70.0` sagen. Das ist der dokumentierte Ablauf
(`releasing.md`: *„Der Digest-Pin in Handbuch §2 entsteht **nach** dem Tag"*), und für
`v0.69.0` hat ihn `91df58c` nachgezogen. **Erinnerung, kein Befund:** ohne diesen
Folge-Commit ist das Handbuch zwischen Tag und Nachzug in sich widersprüchlich.

### I3 — „eine einzeilige Inline-Code-Spanne meldet nicht" trägt ohne das Handbuch-Beispiel weiter, als gemeint

Der CHANGELOG (`:38`) und die §11-Zeile lassen das Beispiel weg, das im Handbuch danebensteht
(``` `- [ ] Beispiel` ```). Gemeint ist eine Spanne, die den **Listen-Marker einschließt**;
ein Task-Item, das inline-Code nur *enthält*, meldet sehr wohl:

```text
f-item-mit-code.md, Zeile 5:  - [ ] Punkt mit `inline-code` im Text
  ⇒ f-item-mit-code.md:5   section-tasks-open
```

Die Handbuch-Fassung ist präzise; die beiden Kurzfassungen sind es knapp nicht. Kein
Widerspruch, nur eine Lesart, die ein Konsument falsch nehmen kann.

### I4 — `--id-prefix` hat keine eigene Zeile in der Optionen-Tabelle

`docs/user/operations.md`: 15 Tabellenzeilen gegen 17 Flags des Binaries. Aufgelöst:
`--enable`/`--disable` teilen sich eine Zeile, `--id-prefix` steht nur im Fließtext der
`--suggest-config`-Zeile (`:42`). Vorbestehend, **nicht** durch dieses Release entstanden —
die Prep-Behauptung *„keine neue CLI-Option"* ist korrekt.

---

## Bestätigte Zusagen (Belege)

### B1 — Der wandernde Anker

```text
$ grep -c "<a id=" version.md
1
$ grep -n "<a id=" version.md
35:| `v0.70.0` <a id="v0.70.0"></a> | 2026-08-30 | [Tag v0.70.0](…/releases/tag/v0.70.0) |
```

Probe-Baum mit einer Kopie von `version.md` und vier festen Pins:

```text
pin.md:5   version.md#v0.69.0   anchor-missing
pin.md:6   version.md#v0.68.0   anchor-missing
```

`#aktuell` und `#v0.70.0` melden **nicht** — sie lösen auf. `#v0.69.0` bricht, und genau das
ist der Zweck des Registers.

### B2 — Alle Versions-Pins, `ghcr` und nackt

```text
$ grep -ro "ghcr.io/pt9912/d-check:v0.70.0" README.md README.de.md docs/user/benutzerhandbuch.md | wc -l
26
$ grep -rn ":v0.70.0" README.md README.de.md docs/user/benutzerhandbuch.md | grep -v "ghcr.io"
docs/user/benutzerhandbuch.md:82:docker pull pt9912/d-check:v0.70.0
docs/user/benutzerhandbuch.md:87:- `:v0.70.0` — eine feste Version (empfohlen für reproduzierbare Läufe; die jeweils
```

26 + 2 — exakt die Zahlen der Commit-Botschaft, und die Zwei aus `releasing.md` stimmt.
Repo-weit bleibt `v0.69.0` an genau zwei Stellen stehen, beide Historie: die
`version.md`-Verlaufszeile und die Handbuch-§11-Zeile `1.63`.

Zusätzlich geprüft: `packaging/dockerhub/overview.md` nennt weiter *„21 rule modules"*
(richtig, kein neues Modul) und pinnt per Digest-Platzhalter — kein Nachzug nötig.

### B3 — Die sieben Datumsstempel

```text
$ git log -1 --date=short --format=%ad
2026-08-30
```

`version.md` §Aktuell · `version.md` §Verlauf (neue Zeile) · `CHANGELOG`-Überschrift ·
Handbuch-Kopf (`**Stand:**`) · Handbuch-§11-Zeile `1.64` · Lastenheft-Historie `0.79.0` ·
Spezifikations-Historie — **alle sieben** tragen `2026-08-30`, und das ist der Kalendertag
des Commits, nicht die Vorgängerzeile.

### B4 — Die §11-Zeile

`docs/user/benutzerhandbuch.md:2438` ist die **letzte** Zeile der Datei; darüber steht
`1.63 | v0.69.0`. Die Tabelle ist die einzige aufsteigende des Repos, und `structure` hält
ihre Monotonie im inneren Loop — `make doc-check` grün. Inhaltlich deckt sie sich mit dem
CHANGELOG-Eintrag (vier Grenzen, Lexik, ein Befund je Haken, negativer Wert ⇒ Exit 2,
byte-identisch ohne den Schlüssel) und widerspricht ihm an keiner Stelle. Sie nennt
slice-172 **nicht** — richtig: das Handbuch beschreibt das Werkzeug, und slice-172 hat es
nicht angefasst.

### B5 — Die Zahl `zehn`

Bedingungen der `structure`-Regel laut `spec/lastenheft.md` `DC-FA-STRUCT-001`
(Tabelle *Bedingungen im Abschnitt*, Zeilen mit eigenem Schlüssel, „dieselbe Bedingung"
nicht doppelt gezählt): `non-empty` · `min-sentences` · `max-tasks` · `max-open-tasks` ·
`forbid-pattern` · `require-pattern` · `require-all` · `table.order` · `headings-match` ·
`table.column` = **10**. Vorher neun. Beide READMEs sagen `zehn`/`ten`, und beide tragen
denselben Vier-Zeilen-Einschub an derselben Stelle mit derselben Aussage — die
DE↔EN-Synchronität der neuen Passage hält (die Bullet-Modul-Liste zählt in beiden Fassungen
21 Einträge in identischer Reihenfolge).

### B6 — `144 Befunde in 37 Dateien`, mit Ausnahme `null`

```text
$ docker run --rm --network none -v "$PWD:/repo:ro" -v "<ohne-ausnahme>.yml:/repo/.d-check.closure.yml:ro" \
    d-check:latest --config .d-check.closure.yml --enable structure
d-check: 223 Datei(en) geprüft, 144 Befund(e)
$ … | cut -f1 | sed 's/:[0-9]*$//' | sort -u | wc -l
37
$ … mit den drei exempt-paths der Regel …
d-check: 223 Datei(en) geprüft, 0 Befund(e)
```

Exakt reproduziert, beide Zahlen.

### B7 — Byte-Identität ohne den Schlüssel

s. M3 — `diff` leer, 178 Befunde beidseits. Zusätzlich das Hauptprofil über den ganzen Baum:
`d-check:vor178` und `d-check:latest` melden beide Exit 0 und null Befunde, `diff` leer.
Gegenprobe zur Zuschreibung: `git show --stat` für `79bf375` und `a362692` (slice-172) zeigt
`.d-check.closure.yml`, `AGENTS.md`, `harness/**` und den Slice-Plan — **kein**
`internal/`, `cmd/` oder `Dockerfile`. `Changed` statt `Added` ist damit belegt richtig, und
die Zusage *„Modulsatz, Grund-Codes und Konfigurations-Fläche des Produkts sind unberührt"*
hält.

### B8 — Grenze 1 und 4: Fence außen vor, vergessener Schluss-Fence blind

```text
a-fence.md   (Haken frei + Haken im wohlgeformten Fence)  ⇒ 1 Befund, Zeile 5
e-fence-offen.md (Haken hinter nie geschlossenem Fence)   ⇒ 0 Befunde
e-fence-offen.md mit --enable spans                       ⇒ e-fence-offen.md:5  ```text  fence-unclosed
```

Beide Grenzen exakt wie beschrieben, samt der Konsequenz „`spans` im selben Profil".

### B9 — Grenze 2: einzeilig schweigt, mehrzeilig zählt mit

```text
b-inline-einzeilig.md   Text `- [ ] offen-inline` Text            ⇒ 0 Befunde
c-inline-mehrzeilig.md  Spanne über drei Zeilen mit Haken darin   ⇒ c-inline-mehrzeilig.md:6  section-tasks-open
```

Dieselbe mehrzeilige Spanne meldet unter `max-tasks: 0` **nicht** — die Aussage „das ist der
ausgewiesene Preis der rohen Lesung" trägt.

### B10 — Grenze 3: Blockquote und Tabulator zählen für keine der beiden Bedingungen

```text
d-blockquote-tab.md   > - [ ] blockquote   und   - [<TAB>] tabbox
  max-open-tasks: 0   ⇒ 0 Befunde
  max-tasks: 0        ⇒ 0 Befunde
```

Beide Bedingungen, wie behauptet.

### B11 — Lexik, Reihenfolge, negativer Wert, Backtick-Ausfall

```text
lexik.md, max-open-tasks: 0  ⇒ Befunde auf 5,6,7,8,9,10
  (- · * · + · geordnete Liste · eingerückt · Tab-Trenner)
  [x] und [X] auf 11/12 melden NICHT
erste.md, max-open-tasks: 1  ⇒ Befunde auf 6 und 7, nicht auf 5
  (die ersten N in Dokument-Reihenfolge sind erlaubt; ein Befund je Item auf SEINER Zeile)
max-open-tasks: -1           ⇒ "max-open-tasks -1 muss >= 0 sein", Exit 2
```

Der Backtick-Ausfall, auf den die ganze Begründung steht:

```text
ohne-backtick.md / mit-backtick.md, derselbe offene Haken
  forbid-pattern '- \[ \]'  ⇒ 1 Befund  (nur ohne-backtick.md)
  max-open-tasks: 0         ⇒ 2 Befunde (beide)
```

*„fällt von einem Befund auf null"* — genau so gemessen.

### B12 — `operations.md` und die gefenceten YAML-Beispiele

```text
$ grep -o "`[a-z]*`" der --enable/--disable-Zeile  ⇒ 21 Modulnamen
$ grep "^modules:" … internal/hexagon/core/model/config.go
   21 Einträge, gleiche Reihenfolge
$ grep -n "max-tasks\|section-\|min-sentences" docs/user/operations.md
   (keine Ausgabe)
```

Alle drei Prep-Behauptungen zu `operations.md` bestätigt. Das `structure`-Beispiel des
Handbuchs (`:1983-2021`) gegen den echten Validator:

```text
$ docker run --rm --network none -v "$PWD:/repo:ro" -v "<hb-beispiel>.yml:/repo/.d-check.closure.yml:ro" \
    d-check:latest --config .d-check.closure.yml
d-check: 562 Datei(en) geprüft, 0 Befund(e)   Exit 0
```

Es lädt, validiert und läuft. **Anmerkung zur Deckung:** die gefenceten YAML-Beispiele sind
inzwischen **nicht mehr** ungewächtert — `internal/adapter/driven/configyaml/docexamples_test.go`
zieht jeden gefenceten YAML-Block aus Handbuch, `operations.md` und beiden READMEs durch
`configyaml.Decode` und ist fail-closed zum Prüfen hin. Der Blind-Spot, den
`releasing.md` §4 noch als offen beschreibt, ist an dieser Achse geschlossen; offen bleibt,
ob ein Beispiel das **tut**, was der Text daneben verspricht — das prüft kein Gate.

### B13 — SemVer

`0.70.0` (MINOR) ist richtig. Zugewachsen sind ein Konfigurations-Schlüssel und ein
Grund-Code, beide opt-in; nichts entfernt, nichts umbenannt, keine Default-Änderung —
gemessen als byte-identischer Befundsatz ohne den Schlüssel (B7). Eine Konfiguration, die
`v0.69.0` grün fuhr, fährt `v0.70.0` unverändert grün; eine, die `max-open-tasks` nennt,
brach vorher mit Exit 2 und wird jetzt akzeptiert — rein additiv. `--print-mk` verteilt den
Tag aus der Build-Version, die Handbuch-Wiedergabe (`:850`) ist entsprechend auf `v0.70.0`
gezogen und nennt die zwölf Targets, die das Binary erzeugt (nachgezählt: 12, gleiche Namen).

---

## Was ich nicht geprüft habe

- Den **Inhalt** von slice-178 und slice-172 (je einmal reviewt und verifiziert; Auftrag
  ausdrücklich ausgenommen).
- Die **Gültigkeit** des Digests in Handbuch §2 gegen die Registry (Netz).
- Ob die Zahl `276` mit *irgendeinem* Korpus reproduzierbar ist — geprüft wurden drei
  plausible Varianten; die Aussage in M3 ist „aus dem publizierten Text nicht
  nachvollziehbar", nicht „falsch".
- Die Übersetzungsqualität der EN-Fassung jenseits der von diesem Commit berührten Passage.
