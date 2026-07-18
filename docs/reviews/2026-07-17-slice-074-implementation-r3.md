# Review slice-074 (R3) — Direktiven-Zelle in Tabellenzeilen

**Datum:** 2026-07-17 · **Rolle:** unabhängiger Reviewer (kontext-getrennt,
nicht der Autor) · **Lauf:** R3 nach BLOCK in
[R1](2026-07-17-slice-074-implementation-r1.md) und
[R2](2026-07-17-slice-074-implementation-r2.md), **vor** Release v0.45.2.

**Gegenstand:** `e8b66ec` (R2-Befunde: Toleranz an der Tabellengrenze verengt),
`670ebaf` (Selbstbefund: Lookahead auf die Nachsichts-Entscheidung begrenzt).
Diff `990dc68..HEAD`.

**Quellen:** [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritte 3/5, [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency),
[ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md),
[ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md),
[slice-074](../plan/planning/in-progress/slice-074-kommentar-suffix-tabellenzeilen.md),
[`AGENTS.md`](../../AGENTS.md) §3.

**Verifikations-Basis:** `make build` (HEAD, Image `sha256:2f77484826…`) gegen das
**ausgelieferte** `ghcr.io/pt9912/d-check:v0.45.1` (nicht gegen einen
nachgebauten Stand). Läufe `docker run --rm --network none -v <fixture>:/repo:ro
-w /repo <image> --trace [--require-complete]`. Mutationen via `make test` in
einem Scratch-Worktree (entfernt; Produktivbaum unverändert, `git status` leer).
Realdaten: **Kopie** von grid-gyms `spec/`+`docs/` (grid-gym selbst nur gelesen).

**Vorbemerkung.** `e8b66ec` schließt R2-F-1 in der geprüften Form: der
Direktiven-**Header** einer Folgetabelle wird nicht mehr gefressen
(`TestTraceTableFolgeHeaderWirdNichtGefressen` pinnt es, Mutation bestätigt). Die
R2-Doku-Befunde F-3 (Schritt `3a`), F-5 (Konsumententest Datenzeile) und
F-6/R1-F-4 (Roadmap) sind geschlossen. **Der Befund unten ist keine Wiederholung
von R2-F-1, sondern dessen Spiegelbild** — und er ist nicht abgefangen.

---

## Findings

### F-1 · HIGH · Die Toleranz greift doch über die Tabellengrenze: eine tolerierte Direktiven-**Datenzeile** verschluckt die **gesamte** folgende Tabelle — Anforderungen und Rück-Kanten verschwinden lautlos, Exit 1 wird Exit 0

**quelle:** [ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md)
§Entscheidung 3 („Die Zusage lautet jetzt enger und prüfbar: keine Zeilenart
verliert Zellen, und **die Toleranz greift nie über eine Tabellengrenze**; beide
Grenzen sind per Mutation gepinnt") und §Entscheidung 4 („SemVer-**Patch**:
gegenüber v0.45.1 wird **keine** Zeile anders gelesen, die heute gelesen wird");
[`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 5 („**Die Nachsicht endet an der Tabellengrenze:** eine Zeile, der eine
passende Trennzeile folgt, ist der Header einer **neuen** Tabelle und wird nie als
Datenzeile der laufenden toleriert — sonst verschwänden die Anforderungen der
Folgetabelle **lautlos**"); Reviewer-Anker HIGH („Stilles-Grün-Pfad in einem Gate
— Harness-Lüge").

**pfad:** `internal/hexagon/core/app/trace_table.go:166-170` (die
Nachsichts-Bedingung) im Zusammenspiel mit
`internal/hexagon/core/app/trace_table.go:180-186` (`isNewTableHeader`) und
`internal/hexagon/core/app/trace_table.go:132-144` (`i = next - 1`).

**befund:** `isNewTableHeader` fragt ausschließlich: *„ist die N+1-Zeile
**selbst** der Header einer neuen Tabelle?"* Sie fragt nicht: *„führt das
Tolerieren dieser Zeile dazu, dass ich einen **nachfolgenden** Header
verschlucke?"* Genau das passiert.

In v0.45.1 setzte **jede** `badLine` den Header-Scan neu auf: `consumeTableRows`
gab `j` zurück, `markdownTables` setzte `i = j` und bekam an `j`, `j+1`, `j+2` …
wieder die Chance, einen Header zu erkennen. Die Toleranz **entfernt diesen
Wiederaufsetz-Punkt**. Die unmittelbar folgende Zeile — der Header der neuen
Tabelle — hat exakt `N` Zellen und läuft damit in den `==`-Zweig, der
**bewusst keinen Lookahead** hat (`670ebaf`). Sie wird als Datenzeile konsumiert;
Trennzeile und alle Datenzeilen der Folgetabelle folgen ihr. Die Folgetabelle
wird **nie erkannt**, bindet keine Rolle, und ihre Anforderungen existieren nicht.
Ist die laufende Tabelle nicht relevant, bleibt ihre `badLine` folgenlos — der
Lauf schweigt. Bleibt **eine andere** Tabelle relevant, greift auch der
`foundTable`-Guard aus [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md)
nicht.

Belegt gegen das **ausgelieferte** Image, beide Konsumenten. Fixture `fx-m`
(`format: table`; Tabelle 1 relevant/gedeckt, Tabelle 2 **irrelevant** mit
Direktiven-Datenzeile, Tabelle 3 ohne Leerzeile dahinter, relevant, `F-2`/`F-3`
echte Waisen):

```markdown
| ID | Titel | Notiz |
|---|---|---|
| X | Y | Z | <!-- d-check:ignore (Grund) -->     ← toleriert (N+1)
| Anforderung | Beschreibung | Status |            ← N Zellen ⇒ als Datenzeile gefressen
|---|---|---|
| F-2 | Beta | ok |
| F-3 | Gamma | ok |
```

```text
=== fx-m · trace.requirements.format: table · --trace --require-complete ===
v0.45.1  EXIT=1   | F-1 | Alpha | — | slice-001 | ok    |
                  | F-2 | Beta  | — | —         | WAISE |
                  | F-3 | Gamma | — | —         | WAISE |
                  3 Anforderung(en), 2 Waise(n).
HEAD     EXIT=0   | F-1 | Alpha | — | slice-001 | ok |
                  1 Anforderung(en), 0 Waise(n).

=== fx-x4 · trace.cross-consistency · --trace --require-complete ===
v0.45.1  EXIT=1   1 Differenz(en).   (Rück-Kante C-2 → F-2)
HEAD     EXIT=0   0 Differenz(en).
```

Die Kausalität ist durch zwei Kontrollen **isoliert** — die Toleranz ist die
alleinige Ursache, nicht die Nachbarschaft der Tabellen:

| Kontrolle | v0.45.1 | HEAD |
|---|---|---|
| `fx-m2` — identisch, **mit** Leerzeile vor Tabelle 3 | 3/2, Exit 1 | 3/2, Exit 1 |
| `fx-m3` — identisch, 4. Zelle **kein** Kommentar (`\| X \| Y \| Z \| W \|`) | 3/2, Exit 1 | 3/2, Exit 1 |
| `fx-m` — 4. Zelle **ist** die Direktive | 3/2, **Exit 1** | 1/0, **Exit 0** |

Damit sind **beide** Zusagen falsifiziert:

- „Die Toleranz greift **nie** über eine Tabellengrenze" — sie greift hier über
  eine ganze Tabelle hinweg, nicht nur über deren Grenze.
- „Gegenüber v0.45.1 wird **keine** Zeile anders gelesen, die heute gelesen wird"
  — die Zeile `| Anforderung | Beschreibung | Status |` wird von v0.45.1 als
  **Header** gelesen und von HEAD als **Datenzeile**. Sie wird heute gelesen, und
  sie wird anders gelesen.

Ebenso die normative Rang-2-Zusage in Spec Schritt 5: die dort geforderte Zeile
(„eine Zeile, der eine passende Trennzeile folgt, … wird **nie** als Datenzeile
der laufenden toleriert") **ist** in `fx-m` Zeile 9 mit passender Trennzeile in
Zeile 10 — und sie wird konsumiert. Der Einwand, ein exakt breiter Treffer werde
nicht „toleriert", sondern gehöre nach Satz 1 schlicht zur Tabelle, trägt die
Regel nicht: die Spec **begründet** sie mit „sonst verschwänden die Anforderungen
der Folgetabelle lautlos", und genau das tritt ein. Eine Regel, die ihren eigenen
angegebenen Zweck im Nachbarfall verfehlt, ist als Vertrag nicht erfüllt.

Die ADR-Fitness-Funktion listet den Fall „Folgetabelle mit **Direktiven**-Header"
(gedeckt, `e8b66ec`), aber keinen Fall „Folgetabelle mit **normalem** Header nach
einer tolerierten Direktiven-**Datenzeile**". Das ist dieselbe Bewegung wie in
R1, R2 und im Selbstbefund `670ebaf`: **ein Zweig geprüft, auf die Klasse
geschlossen** — hier der Fall „die tolerierte Zeile *ist* der Header" geprüft, der
Fall „die tolerierte Zeile steht *vor* dem Header" nicht.

Der motivierende Realfall bleibt gelöst: grid-gym `architecture.md:913` ist eine
Datenzeile inmitten ihrer Tabelle; der Realdatenlauf ist gegen beide Images
identisch (`151 Anforderung(en), 0 Waise(n)`) — der Befund ist **kein**
Praxis-Regress an den Realdaten, sondern ein stiller Gate-Pfad und eine
falsifizierte normative Zusage.

**verifizierbar:** ja — `docker run --rm --network none -v fx-m:/repo:ro -w /repo
d-check:latest --trace --require-complete` liefert Exit 0/„0 Waise(n)",
`ghcr.io/pt9912/d-check:v0.45.1` Exit 1/„2 Waise(n)". Als Regressionstest bände
eine Konsumenten-Fixture „irrelevante Tabelle mit Direktiven-Datenzeile, ohne
Leerzeile gefolgt von einer relevanten Tabelle mit N-Header" den Befund.

---

### F-2 · MEDIUM · Beide Grenzen, die `670ebaf` einzieht, sind **ungepinnt** — die Suite ist mutations-blind an genau der Stelle, an der der Autor seinen Selbstbefund verortet hat

**quelle:** [ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md)
§Entscheidung 3 („beide Grenzen sind **per Mutation gepinnt**");
`internal/hexagon/core/app/trace_table.go:159-165`/`176-179` (die Kommentare
behaupten die Invariante ausdrücklich); Reviewer-Anker MEDIUM (fehlende
Negativtests am neuen Vertrag).

**pfad:** `internal/hexagon/core/app/trace_table.go:166-167`
(`len(cells) != len(t.header) &&`) und
`internal/hexagon/core/app/trace_table.go:181-183` (`j+1 >= len(lines)`).

**befund:** `670ebaf` ist der Selbstbefund „Lookahead war zu breit platziert" und
zieht zwei Grenzen ein. **Keine** von beiden kippt einen Test — nachgemessen
gegen `make test` im Scratch-Worktree:

| Mutation | Suite | tatsächliches Verhalten |
|---|---|---|
| **M-1** `\|\| isNewTableHeader(lines, j)` entfernt | **kippt 1** (`TestTraceTableFolgeHeaderWirdNichtGefressen`) | ✔ wie zugesagt |
| **M-3** `len(cells) != len(t.header) &&` entfernt (= `670ebaf` zurückgedreht) | **GRÜN** | Verhalten ändert sich |
| **M-4** `j+1 >= len(lines)` entfernt | **GRÜN** | **Panic** |

Beide Mutationen sind **keine** No-ops — belegt am echten Image:

```text
=== fx-n · exakt breite Zeile, die zugleich ein Tabellen-Header ist · KEINE Direktive ===
v0.45.1        2 Anforderung(en), 1 Waise(n).            EXIT=0
HEAD           2 Anforderung(en), 1 Waise(n).            EXIT=0
M-3-Image      error: Tabellenzeile 6 hat 3 statt 3 Zellen   EXIT=2

=== fx-o · tolerierte Direktiven-Zeile = LETZTE Zeile, ohne Trailing-Newline ===
v0.45.1        error: Tabellenzeile 6 hat 3 statt 2 Zellen   EXIT=2
HEAD           2 Anforderung(en), 1 Waise(n).            EXIT=0
M-4-Image      panic: runtime error: index out of range [6] with length 6
```

M-3 bricht die Byte-Identität zu v0.45.1 auf einer Datei **ohne jede Direktive**
und erzeugt die absurde Meldung „3 statt 3 Zellen" — die gesamte Suite bleibt
grün. M-4 ist ein realer Panic-Pfad: eine Datei ohne abschließenden Zeilenumbruch,
deren letzte Zeile die Direktive trägt, ist keine Konstruktion, sondern die
`fx-o`-Form (HEAD behandelt sie korrekt) — die Suite bleibt grün.

Die Zusage „beide Grenzen sind per Mutation gepinnt" ist damit für die **eine**
Grenze wahr (M-1) und für die **beiden** Grenzen aus `670ebaf` falsch. Praktische
Folge: ein späterer Refactor kann den Selbstbefund geräuschlos rückgängig machen
oder den Panic-Guard entfernen, und `make gates` bleibt grün — dieselbe
Sensor-Lücke, die R1-F-3 („die Suite pinnte die defekte Platzierung") schon
einmal als Klasse benannt hat.

**verifizierbar:** ja — die drei Mutationsläufe sind mit `make test` in einem
Worktree reproduzierbar; `fx-n`/`fx-o` gegen die Mutations-Images belegen, dass
die Mutationen verhaltensändernd sind.

---

### F-3 · LOW · Die Mutations-Zusage „Splitter-Platzierung wiederhergestellt kippt 7" ist nicht reproduzierbar

**quelle:** Zusage des Autors zum R3-Lauf; slice-074 §3 DoD („**Mutations-Härte**",
Checkbox offen); Maintainability (prüfbare Zusage).

**pfad:** `internal/hexagon/core/app/trace_table.go:364-386` (Splitter),
`internal/hexagon/core/app/trace_coverage_test.go:250` (`TestCellCountOK`).

**befund:** Mutation M-2 (`dropCommentSuffix` aus `44a5201` zurück in
`splitPipeTableLine`, `cellCountOK` unverändert) kippt **5** FAIL-Zeilen und
**3** distinkte Testfunktionen, nicht 7:

```text
--- FAIL: TestSplitPipeTableLineLaesstKommentarZelleStehen
--- FAIL: TestSplitPipeTableLineLaesstKommentarZelleStehen/Kommentar_am_Ende_ist_eine_Zelle
--- FAIL: TestSplitPipeTableLineLaesstKommentarZelleStehen/auch_mit_Schluss-Pipe
--- FAIL: TestTraceTableFolgeHeaderWirdNichtGefressen
--- FAIL: TestTraceTableHeaderMitDirektiveWirdErkannt
```

Die **Richtung** der Zusage trägt (die defekte Platzierung wird abgelehnt, und
zwar auch auf Konsumenten-Ebene — die von R1-F-3 geforderte Eigenschaft); nur die
Zahl stimmt nicht. Kein Korrektheitsbefund, aber die Zahl steht als Beleg in der
DoD-Diskussion und ist so nicht haltbar.

**verifizierbar:** ja — M-2 mit `make test` in einem Worktree.

---

### F-4 · LOW · R2-F-2 ist nur zur Hälfte geschlossen: §1/§2 und die Mutations-Härte-Zeile des Slices tragen weiter die **verworfene** Suffix-Prämisse

**quelle:** [R2](2026-07-17-slice-074-implementation-r2.md) F-2;
[ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md) §Entscheidung
(„Die erste Fassung dieser ADR nannte den abschließenden Kommentar ein ‚Suffix,
keine Zelle'. Das ist **falsch**"); Source Precedence
([`AGENTS.md`](../../AGENTS.md) §3).

**pfad:** `docs/plan/planning/in-progress/slice-074-kommentar-suffix-tabellenzeilen.md:30`,
`:38`, `:64` (auch der Titel `:1`).

**befund:** Die DoD-Zeile „**Implementierung:** **nicht** im Splitter … Grenze am
nächsten Tabellen-Header (R2-F-1)" ist korrekt nachgezogen — R2-F-2 ist insoweit
geschlossen. Unverändert stehen geblieben sind:

- `:30` „Der Slice macht den abschließenden Kommentar zu einem Zeilen-**Suffix**
  statt zu einer Zelle" — die ADR sagt explizit das Gegenteil.
- `:38` „Der Guard bleibt scharf: **nach dem Entfernen des Suffixes** ist eine
  abweichende Zellenzahl unverändert Exit 2" — es wird nichts mehr entfernt.
- `:64` „**Mutations-Härte:** **Suffix-Abstreifung** entfernt kippt genau einen
  Test" — die einzige offene DoD-Checkbox neben Release/Qualität benennt eine
  Mutation an Code, den es nicht gibt (`grep -rn "dropCommentSuffix" internal/`
  ist leer). Abhakbar ist sie nur, wenn man sie als „Toleranz entfernt" liest —
  das ist M-1 und trifft zu, steht aber so nicht da.

Kein Gate greift (Prosa). Die Verwechslungsgefahr ist real: die DoD ist der
Wortlaut, gegen den die Verifikation abhakt.

**verifizierbar:** ja — `grep -n "Suffix" docs/plan/planning/in-progress/slice-074-kommentar-suffix-tabellenzeilen.md`
gegen `grep -n "ist eine Zelle\|nicht entfernt" spec/spezifikation.md`.

---

### F-5 · INFO · `isNewTableHeader` ist masken-blind — geprüft, in der Praxis unerreichbar

**quelle:** [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
Schritte 2/3 (`sections`/`exclude-sections`); dokumentationswürdige Asymmetrie.

**pfad:** `internal/hexagon/core/app/trace_table.go:180-186`.

**befund:** `consumeTableRows` prüft `maskAllows(mask, lines[j].no)`, reicht die
Maske aber **nicht** an `isNewTableHeader` weiter; `tableHeaderAt` prüft nur
`prose`. Die Funktion kann daher einen „neuen Header" an einer Position sehen,
an der `markdownTables` (das `lines[i]` **und** `lines[i+1]` gegen die Maske
prüft) nie eine Tabelle erkennen würde. Die Richtung ist **fail-closed**: das
Ergebnis wäre eine `badLine` statt einer Toleranz — also exakt v0.45.1-Verhalten,
kein stiller Pfad.

Erreichbar ist der Fall zudem nicht: damit `lines[j]` in der Maske und
`lines[j+1]` außerhalb liegt, muss `j` die letzte Zeile eines Abschnitts sein —
die Folgezeile ist dann die **Heading-Zeile** des nächsten Abschnitts und kann
nie eine Trennzeile sein, `isNewTableHeader` also nie `true`. Empirisch gegen
`fx-s` (`backward.sections`, Direktiven-Zeile am Abschnittsrand) bestätigt: HEAD
toleriert wie vorgesehen (`1 Differenz(en)`), v0.45.1 Exit 2. Kein Finding, aber
die Kopplung ist unausgesprochen und wäre bei einer künftigen Masken-Semantik
(z. B. zeilen- statt abschnittsbasiert) sofort scharf.

**verifizierbar:** ja — Code-Lektüre + `fx-s` gegen beide Images.

---

## Negativbefunde (geprüft, ohne Befund)

- **R2-F-1 (Direktiven-Header der Folgetabelle wird gefressen): geschlossen,
  gegen das ausgelieferte Image nachgefahren.** Die R2-Fixture-Form (Folgetabelle
  mit **Direktiven**-Header, N+1 gegen N+1) läuft unter HEAD laut statt still;
  `TestTraceTableFolgeHeaderWirdNichtGefressen` pinnt es, und Mutation M-1 kippt
  genau diesen Test. Der Fix ist an der geprüften Achse korrekt — F-1 ist der
  **spiegelverkehrte** Fall, nicht derselbe.
- **R2-F-3 (Verweis auf „Schritt 3a"): geschlossen** — `grep -n "3a"
  spec/spezifikation.md` ist leer; der Satz nennt jetzt „Schritte 4/5", was den
  tatsächlichen Schritt-Nummern entspricht.
- **R2-F-5 (kein `format: table`-Konsumententest für die tolerierte Datenzeile):
  geschlossen** — `TestTraceTableDatenzeileMitDirektive` (`trace_cross_test.go:687`)
  prüft Konsumenten-seitig `m.Total == 2` **und** die Titel-Bindung
  (`Requirements[1].Title == "Mit Marker."`), also nicht nur „kein Fehler".
- **R2-F-6 / R1-F-4 (doppelte Backlog-Zeile in der Roadmap): geschlossen** —
  `grep -c "Im Backlog" docs/plan/planning/in-progress/roadmap.md` liefert 1.
- **R1-F-1 / R1-F-2 / R1-F-3: unverändert geschlossen** — `splitPipeTableLine`
  ist gegenüber v0.45.1 weiterhin **byte-identisch** (`git diff
  44a5201~1..HEAD -- internal/hexagon/core/app/trace_table.go` zeigt nur
  `htmlCommentCell`, den `cellCountOK`-Aufruf, `cellCountOK`, `isNewTableHeader`).
  Kein Codepfad entfernt eine Zelle.
- **Byte-Identität zu v0.45.1, sechste Achse gesucht — GEFUNDEN, siehe F-1.** Die
  fünf vom Autor behaupteten Achsen (sauberer Lauf, Header-Direktive, legitime
  Kommentar-Spalte, N+2/N−1, Trennzeile mit Marker) sind bestätigt; die
  **sechste** — eine Zeile, die v0.45.1 als *Header* und HEAD als *Datenzeile*
  liest — bricht sie (`fx-m` Zeile 9). Zusätzlich geprüft und **ohne** Bruch:
  Direktive in der ersten Spalte (fail-closed), Zeile nur aus Kommentar,
  Tabelle nur aus Header+Trennzeile, Datei-Ende ohne Trailing-Newline (`fx-o`,
  HEAD liest korrekt), Marker in der Trennzeile (beide Images überspringen),
  Fences (`markdownTableLines` unverändert).
- **`||`-Kurzschluss-Reihenfolge in `consumeTableRows`:** geprüft, ohne Befund —
  `cellCountOK` und `isNewTableHeader` sind beide seiteneffektfrei; ein Tausch der
  Operanden ist semantisch identisch (nur Kosten, keine Wirkung). Load-bearing ist
  allein der `len(cells) != len(t.header) &&`-Präfix — der ist ungepinnt (F-2).
- **`htmlCommentCell`-Regex:** geprüft, ohne Befund — unverändert seit `44a5201`;
  R1 hat ihn in der gefährlichen Richtung durchgespielt (kein `-->` im Body, keine
  echte Zelle matcht fälschlich). `cells[len(cells)-1]` ist durch den
  `len(cells) == len(header)+1`-Zweig gegen einen leeren Slice geschützt
  (`len(header) >= 1`, da `isTableDelimiter` `len(cells) > 0` verlangt).
- **XREF-Regression über den geteilten Reader:** geprüft — der Reader ist
  derselbe, und F-1 trifft `trace.cross-consistency` **ebenso** (`fx-x4`:
  Exit 1 ⇒ Exit 0). Darüber hinaus keine eigene XREF-Mechanik betroffen:
  `forwardEdges`/`backwardEdges` lesen ausschließlich `row.cells[bt.primary]` /
  `[bt.secondary]`, beide Indizes stammen aus `bindCrossColumns` über `t.header`
  und sind `< len(header) <= len(cells)`; die tolerierte Kommentar-Zelle wird nie
  gelesen und kann keine Kante injizieren.
- **`bindTableColumns` mit Kommentar-Spalte im Header:** geprüft, ohne Befund —
  Bindung über exakte Namensgleichheit; die Direktiven-Zelle trifft keinen
  konfigurierten Rollen-Namen (`TestTraceTableHeaderMitDirektiveWirdErkannt`).
- **Realdatenbeleg grid-gym:** geprüft, ohne Befund — `--trace` gegen eine
  **Kopie** der echten `spec/`+`docs/`: v0.45.1 und HEAD liefern identisch
  `151 Anforderung(en), 0 Waise(n)`. `architecture.md:913` ist eine Datenzeile
  inmitten ihrer Tabelle; die F-1-Vorbedingung (Tabelle ohne Leerzeile direkt
  hinter einer Direktiven-Datenzeile) kommt dort nicht vor. Der motivierende Fall
  ist gelöst. (grid-gym selbst nur gelesen, nichts geändert.)
- **Handbuch §5 / §4.12:** geprüft, ohne **neuen** Befund — die Zeilenart-Regel
  ist korrekt beschrieben, und der neue Absatz „Eine Grenze, die Sie kennen
  sollten" nennt die Mitigation („Lassen Sie zwischen zwei Tabellen eine
  Leerzeile"), die F-1 tatsächlich verhindert (`fx-m2` belegt es). Die
  R2-F-4-Einschränkung besteht fort: „Alles andere bleibt Exit 2" deckt weder den
  fehlgeformten Direktiven-Header noch F-1 — beide sind still. Da R2 das bereits
  als LOW führt und der Absatz die Warnung teilweise nachzieht, kein eigenes
  Finding; mit F-1 fällt die Aussage ohnehin neu an.
- **ADR-Immutabilität ([`AGENTS.md`](../../AGENTS.md) §3.5):** geprüft, ohne
  Befund — ADR-0040 steht auf `Proposed`; die Regel bindet erst `Accepted`. Die
  zweite Neufassung ist in `## Geschichte` mit Anlass, verworfenem Ausweg und
  Konsequenz protokolliert. Bemerkenswert und ausdrücklich positiv: der von R2
  angebotene Ausweg „ehrliche Reichweiten-Angabe statt Code-Verengung" wurde mit
  der Begründung **verworfen**, „eine Zusage zu entschärfen, damit ein stiller
  Verlust hineinpasst, ist die Bewegung, gegen die dieses Werkzeug gebaut ist" —
  das ist die richtige Wahl.
- **„Defekt-Fix, kein CR":** geprüft, weiterhin **haltbar** — das Lastenheft
  ([`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen))
  sagt keine Zellenzahl zu; die Regel lebt in Rang 2 und ist dort fortgeschrieben.
- **SemVer-Patch:** geprüft, **nicht haltbar in der aktuellen Formulierung** —
  siehe F-1: die Zusage „keine Zeile wird anders gelesen" ist falsifiziert. Die
  Einstufung als Patch bleibt vertretbar (die betroffene Form ist degeneriert und
  in GFM ohnehin eine einzige Tabelle), die **Begründung** in der ADR nicht.
- **ADR-0005-Import-Regeln / Modul-Layout:** geprüft, ohne Befund — keine neuen
  Imports, alles in `core/app`.
- **[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit):**
  geprüft, ohne Befund — nur `regexp`/`strings`, alle Läufe mit `--network none`.
- **Gate-Suppression / Schwellen-Senkung ([`AGENTS.md`](../../AGENTS.md) §3.6):**
  geprüft, ohne Befund — keine `//nolint`, keine Gate-/Schwellen-Änderung.
- **Referenz-Richtung (SDP), Marker-Ehrlichkeit:** geprüft, ohne Befund —
  ADR-0040 nennt slice-074 weiterhin nur in `## Geschichte` als Provenance.
- **CHANGELOG / Release-Prep fehlen im Diff:** geprüft, **kein Finding** — die
  Commit-Grenzen-Konvention legt sie in den Release-Prep-Commit.
- **DoD-Abhakung / `make gates`:** nicht geprüft — Rolle der Verifikation, nicht
  des Reviews (Reviewer-Skill §Anti-Pattern). F-4 betrifft den **Wortlaut**.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 1 | F-1 |
| MEDIUM | 1 | F-2 |
| LOW | 2 | F-3, F-4 |
| INFO | 1 | F-5 |

---

## Verdikt

**BLOCK — v0.45.2 darf nicht getaggt werden.**

F-1 ist ein Stilles-Grün-Pfad in einem Gate: ein Lauf, der unter dem
**ausgelieferten** v0.45.1 zwei echte Waisen bzw. eine echte Kreuzverweis-Differenz
meldete (Exit 1), endet unter HEAD Exit 0 und verschweigt sie — in **beiden**
Konsumenten, über den geteilten Reader. Zwei Kontroll-Fixtures isolieren die
Direktiven-Toleranz als alleinige Ursache. Damit ist die zentrale, in `e8b66ec`
**neu** gegebene Zusage („die Toleranz greift **nie** über eine Tabellengrenze")
ebenso falsifiziert wie die SemVer-Begründung („keine Zeile wird anders gelesen")
und der normative Rang-2-Satz in
[`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 5.

**Das Muster ist unverändert, und es ist jetzt das fünfte Mal in dieser
Code-Region.** R1: Body-Regel am Header. R2: Toleranz frisst den
Direktiven-Header. `670ebaf` (Selbstbefund): Lookahead zu breit. Und jetzt: die
**Grenze** gegen R2-F-1 ist zu schmal gedacht — sie fragt „ist diese Zeile ein
Header?", nicht „verschlucke ich durch das Tolerieren einen Header?". Jedes Mal
wurde ein Zweig geprüft und auf die Klasse geschlossen; dreimal war das Ergebnis
ein stilles Grün. Der Reviewer-Skill nennt das ein **Steering-Loop-Signal**: die
Konsequenz gehört über den Einzelfix hinaus. Der strukturelle Kern ist benennbar
und in F-1 §befund ausgeschrieben — **die Toleranz entfernt den
Wiederaufsetz-Punkt des Header-Scans**, den `badLine` in v0.45.1 an jeder
abweichenden Zeile garantierte. Wer diese Invariante direkt adressiert (statt
Nachbarfall um Nachbarfall zu verengen), schließt die Klasse; wer weiter
Einzelfälle verengt, wird R4 bekommen.

F-2 blockiert unabhängig davon **und ist die eigentlich beunruhigende Hälfte**:
die Zusage „beide Grenzen sind per Mutation gepinnt" ist für die beiden Grenzen
aus `670ebaf` **falsch**. Der Selbstbefund lässt sich geräuschlos zurückdrehen
(M-3: Byte-Identität zu v0.45.1 bricht auf einer Datei ohne jede Direktive, Suite
grün), und der Panic-Guard lässt sich entfernen (M-4: realer Panic auf einer
Datei ohne Trailing-Newline, Suite grün). Der Sensor, der die Wiederholung dieser
Klasse verhindern soll, ist an der Stelle blind, an der die Klasse zuletzt
zuschlug — genau der Befund, den R1-F-3 schon einmal als Klasse benannt hat.
F-3/F-4 sind nicht release-blockierend, F-5 ist ohne Befund.

**Was trägt:** Die Ansatz-Umkehr aus R2 bleibt richtig, `splitPipeTableLine` ist
gegenüber v0.45.1 byte-identisch, R1-F-1/F-2/F-3 und R2-F-1/F-3/F-5/F-6 sind
geschlossen, der motivierende Realfall ist gelöst (grid-gym identisch zu v0.45.1,
`architecture.md:913` läuft durch), und die Entscheidung, R2s angebotenen Ausweg
„Zusage entschärfen statt Code verengen" **auszuschlagen**, war die richtige. Der
Abstand zu Grün ist nicht groß — aber er liegt diesmal wieder im Produktivcode,
nicht nur in der Zusage darüber.
