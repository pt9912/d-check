# Review slice-074 (R2) — Direktiven-Zelle in Tabellenzeilen

**Datum:** 2026-07-17 · **Rolle:** unabhängiger Reviewer (kontext-getrennt,
nicht der Autor) · **Lauf:** R2 nach BLOCK/HIGH in
[R1](2026-07-17-slice-074-implementation-r1.md), **vor** Release v0.45.2.

**Gegenstand:** `1210842` (Ansatz-Umkehr: Fix + Tests + ADR-/Spec-Neufassung),
`990dc68` (Handbuch). Diff `6944147..HEAD`.

**Quellen:** [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritte 3/5, [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency),
[ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md) (neu gefasst),
[ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md),
[slice-074](../plan/planning/open/slice-074-kommentar-suffix-tabellenzeilen.md),
[`AGENTS.md`](../../AGENTS.md) §3.

**Verifikations-Basis:** `make build` (HEAD, Image-ID `425a8f1c9ba6`) gegen
`make build IMAGE=d-check-v0451` in einem Worktree auf `44a5201~1` (= der
ausgelieferte Stand v0.45.1). Läufe `docker run --rm --network none -v
<fixture>:/repo:ro -w /repo <image> --trace [--require-complete]`. Mutationen via
`make test` in einem Scratch-Worktree. Realdaten: Kopie von grid-gyms
`spec/architecture.md`/`docs/plan/traceability.md` (grid-gym selbst nur gelesen).
Produktivbaum unverändert, beide Worktrees entfernt.

**Vorbemerkung — die Ansatz-Umkehr trägt.** Der Diff gegen den **ausgelieferten**
Stand ist minimal und strukturell aussagekräftig:

```text
$ git diff 44a5201~1..HEAD -- internal/hexagon/core/app/trace_table.go
+ var htmlCommentCell = regexp.MustCompile(...)
- if len(cells) != len(t.header) {
+ if !cellCountOK(cells, t.header) {
+ func cellCountOK(cells, header []string) bool { ... }
```

`splitPipeTableLine` ist gegenüber v0.45.1 **byte-identisch**; `tableHeaderAt` ist
unberührt. Damit ist die R1-F-1-Mechanik (Zellen verschwinden vor dem Zählen)
tatsächlich beseitigt — kein Codepfad entfernt mehr eine Zelle. Die drei
R1-Befunde sind geschlossen (Belege in den Negativbefunden). Die Findings unten
betreffen eine **andere** Mechanik derselben Klasse sowie die Vertragskette.

---

## Findings

### F-1 · MEDIUM · Der stille Übersprung ist **nicht** strukturell unmöglich: die tolerierte Zeile verlängert die Tabelle über eine Zeile hinweg, die sie vorher beendete — ein nachfolgender Tabellen-Header wird verschluckt, Anforderungen und Rück-Kanten verschwinden lautlos, Exit 1 wird Exit 0

**quelle:** [ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md)
§Entscheidung 3 („Damit ist ein stiller Übersprung **strukturell unmöglich**,
nicht bloß behoben") und §Entscheidung 4 („gegenüber v0.45.1 wird **keine** Zeile
anders gelesen, die heute gelesen wird — es kommt nur die N+1-Datenzeile mit
Direktive hinzu"); Reviewer-Anker HIGH („Stilles-Grün-Pfad in einem Gate").

**pfad:** `internal/hexagon/core/app/trace_table.go:186-191` (`cellCountOK`) im
Zusammenspiel mit `internal/hexagon/core/app/trace_table.go:150-166`
(`consumeTableRows` liefert den Fortsetzungs-Index) und
`internal/hexagon/core/app/trace_table.go:140-144` (`i = next - 1`).

**befund:** Eine Zeile, die v0.45.1 als `badLine` **beendete**, wird jetzt als
Datenzeile **konsumiert**; der Scan setzt dadurch eine Zeile später wieder auf.
Trägt der Header einer unmittelbar folgenden Tabelle die Direktive und hat er
gegenüber der laufenden Tabelle exakt N+1 Zellen, wird er als tolerierte
Datenzeile verschluckt, seine Trennzeile beendet die laufende Tabelle, und die
**folgende Tabelle wird nie erkannt**. Ist die laufende Tabelle nicht relevant,
bleibt ihr `badLine` folgenlos — der Lauf schweigt.

Belegt gegen beide Images, beide Konsumenten:

```text
=== fx-l · trace.requirements.format: table · --trace --require-complete ===
v0.45.1  EXIT=1   | F-1 | Alpha | ADR-0001 | slice-001 | ok |
                  | F-2 | Beta  | — | — | WAISE |
                  | F-3 | Gamma | — | — | WAISE |
                  3 Anforderung(en), 2 Waise(n).
HEAD     EXIT=0   | F-1 | Alpha | ADR-0001 | slice-001 | ok |
                  1 Anforderung(en), 0 Waise(n).

=== fx-x3 · trace.cross-consistency · --trace --require-complete ===
v0.45.1  EXIT=1   | F-2 | C-2 | Rück-Kante, ohne RTM-Eintrag | spec/architecture.md:14 |
                  1 Differenz(en).
HEAD              0 Differenz(en).
```

Die Richtung ist damit erneut laut → still, und erneut über den geteilten Reader
in **beide** Konsumenten. Die Vorbedingungen sind eng und wurden abgegrenzt: keine
Leerzeile zwischen den Tabellen (`fx-l2` mit Leerzeile: beide Images identisch,
3/2), die vorangehende Tabelle nicht relevant (`fx-l3` mit relevanter
Vortabelle: beide Images Exit 2, laut), exakt N+1, letzte Zelle ganzzellig
Kommentar (`fx-i`, `fx-g`: unverändert fail-closed). Die auslösende Form ist
zudem in GFM **keine** zwei Tabellen — ohne Leerzeile ist der zweite Header dort
ebenfalls eine Body-Zeile; HEAD liest insoweit GFM-treuer als v0.45.1, und der
verschluckte Header trägt genau die Direktiven-Form, die das Handbuch seit
`990dc68` empfiehlt. Der Befund ist deshalb nicht als Praxisrisiko HIGH, sondern
als **falsifizierte Universal-Zusage** in einem normativen Dokument MEDIUM: „kein
stiller Übersprung, strukturell unmöglich" und „keine Zeile wird anders gelesen"
sind als Beweis formuliert und durch ein 15-Zeilen-Fixture widerlegt. Genau die
Erschließung der Klasse aus einem geprüften Zweig ist die Klasse, die in dieser
Code-Region zweimal ausgelieferte bzw. abgefangene stille Grüns erzeugt hat.

**verifizierbar:** ja — `docker run --rm --network none -v fx-l:/repo:ro -w /repo
d-check:latest --trace --require-complete` liefert Exit 0/„0 Waise(n)", das Image
auf `44a5201~1` Exit 1/„2 Waise(n)". Als Regressionstest bände eine
Konsumenten-Fixture mit zwei leerzeilenlos benachbarten Tabellen den Befund.

---

### F-2 · MEDIUM · Der Slice-Plan schreibt weiterhin genau die Platzierung vor, die ADR-0040 als Defekt verworfen hat — DoD und ADR widersprechen sich frontal

**quelle:** [slice-074](../plan/planning/open/slice-074-kommentar-suffix-tabellenzeilen.md)
§3 DoD gegen [ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md)
§Entscheidung 2; Source Precedence ([`AGENTS.md`](../../AGENTS.md) §3).

**pfad:** `docs/plan/planning/open/slice-074-kommentar-suffix-tabellenzeilen.md:52-53`
(auch `:30-31`, `:35-41`, `:60-61`, `:72-74`).

**befund:** Die DoD verlangt unverändert „**Implementierung:** im geteilten
Zeilen-**Splitter**, vor dem Zellen-Zählen — Header-, Trenn- und Datenzeilen
gleichermaßen" und „**Mutations-Härte:** Suffix-Abstreifung entfernt kippt genau
einen Test"; §1/§2/§4 tragen durchgehend die Suffix-Prämisse („macht den
abschließenden Kommentar zu einem Zeilen-**Suffix** statt zu einer Zelle", „Der
Splitter trägt ein weiteres Stück Markdown-Wissen"). ADR-0040 verwirft in der
Neufassung exakt diese Platzierung als Ursache des stillen Übersprungs, der Code
tut das Gegenteil, und es gibt weder eine `dropCommentSuffix`-Funktion noch eine
Suffix-Abstreifung mehr. R1 hat denselben DoD-Satz als Quelle des Defekts zitiert
(R1-F-1 §befund); ADR, Spec, Code und Handbuch wurden nachgezogen, der Slice
nicht. Eine DoD-Abhakung gegen diesen Wortlaut ist entweder unmöglich oder
verlangt die Rückkehr zum Defekt.

**verifizierbar:** ja — `grep -n "Splitter" docs/plan/planning/open/slice-074-kommentar-suffix-tabellenzeilen.md`
gegen `grep -n "nicht im Splitter" docs/plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md`;
`grep -rn "dropCommentSuffix" internal/` ist leer.

---

### F-3 · LOW · Die neu gefasste Spec verweist auf einen Schritt, den es nicht gibt

**quelle:** [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 3 (Rang-2-Vertrag); Maintainability (Doku-Drift).

**pfad:** `spec/spezifikation.md:380`.

**befund:** Der Satz „Die Nachsicht ist nach **Zeilenart** getrennt (Schritte
**3a**/5)" verweist auf einen Schritt `3a`; `DC-FA-REQ-001.a` hat die Schritte
1–6 und keinen Schritt 3a. Der Verweis ist im selben Commit entstanden, der die
Regel normativ festschreibt. Kein Gate greift: es ist Prosa, kein Link/Anker.

**verifizierbar:** ja — `grep -n "3a" spec/spezifikation.md` liefert genau diesen
einen Treffer; die Schritt-Nummern der Sektion sind
`grep -E "^[0-9]+\. \*\*"`-sichtbar.

---

### F-4 · LOW · Handbuch: „Alles andere bleibt Exit 2" ist unwahr für den fehlgeformten Direktiven-Header — der ist still

**quelle:** [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 3; Handbuch-Wahrheit (Doku-Drift, latente Falle).

**pfad:** `docs/user/benutzerhandbuch.md:1240-1241` (Aussage) gegen `:1236-1238`
(Header-Anweisung).

**befund:** Das Handbuch führt Nutzer erstmals aktiv an den Direktiven-**Header**
(„die Trennzeile braucht dann eine Spalte mehr") und schließt mit „Alles andere
bleibt Exit 2: zwei überzählige Zellen, eine überzählige Nicht-Kommentar-Zelle,
eine fehlende Zelle". Wer der Anweisung **falsch** folgt — Marker im Header, aber
N-Trennzeile — bekommt kein Exit 2, sondern Stille: die Tabelle wird nicht
erkannt, ihre Anforderungen verschwinden, Exit 0.

```text
=== fx-d · Header mit Marker + N-Trennzeile · --trace --require-complete ===
v0.45.1  EXIT=0   1 Anforderung(en), 0 Waise(n).   (F-2/F-3 unsichtbar)
HEAD     EXIT=0   1 Anforderung(en), 0 Waise(n).   (F-2/F-3 unsichtbar)
```

Das Verhalten ist **vorbestehend** (beide Images identisch, GFM-konform: die Form
rendert auch dort nicht als Tabelle) und keine slice-074-Regression — die
Enumeration, die es als Exit 2 mitbehauptet, ist neu.

**verifizierbar:** ja — Fixture `fx-d` gegen beide Images.

---

### F-5 · LOW · Für `trace.requirements.format: table` pinnt kein Konsumententest die tolerierte **Datenzeile** — also gerade den motivierenden Realfall

**quelle:** [slice-074](../plan/planning/open/slice-074-kommentar-suffix-tabellenzeilen.md)
§3 DoD („Tests (positiv): … **je einmal** für `trace.requirements.format: table`
und `trace.cross-consistency`, auf **Konsumenten-Ebene**"); Reviewer-Anker MEDIUM
(Negativtests am neuen Vertrag) — hier positiv-seitig, daher LOW.

**pfad:** `internal/hexagon/core/app/trace_cross_test.go:623`
(`TestTraceTableHeaderMitDirektiveWirdErkannt`, Header-Fall),
`internal/hexagon/core/app/trace_coverage_test.go:250` (`TestCellCountOK`,
Unit-Ebene).

**befund:** Nachgemessen per Mutation gegen `make test` — Toleranz in
`cellCountOK` entfernt (M-B):

| Mutation | kippt |
|---|---|
| M-A Regel **zurück in den Splitter** (R1-F-3-Gegenprobe) | `TestSplitPipeTableLineLaesstKommentarZelleStehen/*`, `TestCellCountOK/eine_Kommentar-Zelle…`, **`TestTraceTableHeaderMitDirektiveWirdErkannt`** |
| M-B Toleranz ganz entfernt | `TestCellCountOK/eine_Kommentar…`, `TestCrossConsistencyKommentarSuffix` |
| M-C „genau eine" aufgehoben (`>= N+1`) | `TestCellCountOK/zwei_zu_viel…` |
| M-D Kommentar-Bedingung entfernt | `TestCellCountOK/*`, `TestTraceTableEchterZellenbruchBleibtLaut`, `TestCrossConsistencyZellenzahlGuardBleibtScharf`, `TestCLI070_TraceTable_Negative/row_width` |

M-B kippt **keinen** `format: table`-Konsumententest: der Konsumenten-Fall dieses
Readers ist dort der **Header** (`TestTraceTableHeaderMitDirektiveWirdErkannt`),
die tolerierte Datenzeile — grid-gym `architecture.md:913`, der Anlass des Slices
— nur unit-seitig (`TestCellCountOK`) und über den **anderen** Konsumenten
(`TestCrossConsistencyKommentarSuffix`, Datenzeile). Für `format: table`, den
älteren der beiden Konsumenten, bleibt der positive Datenzeilen-Pfad
konsumentenseitig ungepinnt.

**verifizierbar:** ja — die Mutationsläufe sind mit `make test` in einem Worktree
reproduzierbar.

---

### F-6 · LOW · R1-F-4 (doppelte Backlog-Zeile in der Roadmap) ist unverändert offen

**quelle:** [R1](2026-07-17-slice-074-implementation-r1.md) F-4;
Maintainability (Doku-Drift).

**pfad:** `docs/plan/planning/in-progress/roadmap.md:60` und `:62`.

**befund:** Beide wortgleichen Vorkommen von „**Im Backlog (`next/`), auf Aufnahme
in eine Welle wartend:** derzeit keiner." stehen weiterhin da; der R2-Diff fasst
die Roadmap nicht an.

**verifizierbar:** ja — `grep -n "Im Backlog" docs/plan/planning/in-progress/roadmap.md`
liefert zwei Treffer.

---

### F-7 · INFO · Die Fehler-Zeilennummer eines echten Zellenbruchs kann sich um eine Zeile verschieben

**quelle:** dokumentationswürdige Folge der Toleranz; [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)-nah
(Meldungsstabilität), kein Korrektheitsbefund.

**pfad:** `internal/hexagon/core/app/trace_table.go:159-162`.

**befund:** Wird eine Direktiven-Zeile toleriert, wandert ein nachfolgender
echter Bruch als `badLine` eine Zeile weiter. `fx-l3` (relevante Vortabelle,
Direktiven-Header dahinter): v0.45.1 meldet „Tabellenzeile 6 hat 3 statt 2
Zellen", HEAD „Tabellenzeile **7** hat 3 statt 2 Zellen". Beide Exit 2 — der Pfad
bleibt laut, nur die genannte Zeile ist eine andere.

**verifizierbar:** ja — Fixture `fx-l3` gegen beide Images.

---

## Negativbefunde (geprüft, ohne Befund)

- **R1-F-1 (Header-Suffix ⇒ stilles Grün): geschlossen, gegen das echte Image
  nachgefahren.** Fixture `fx-c` (Header mit Direktive + N+1-Trennzeile, zweite
  Tabelle relevant): v0.45.1 **und** HEAD liefern identisch `3 Anforderung(en), 2
  Waise(n)`, Exit 1. Das Image auf `44a5201` lieferte hier Exit 0/„0 Waise(n)".
  XREF-Gegenstück `fx-x2`: beide Images `1 Differenz(en)`. Ursächlich strukturell
  gedeckt: `splitPipeTableLine` ist gegenüber v0.45.1 byte-identisch, kein Pfad
  entfernt Zellen.
- **R1-F-2 (legitime Kommentar-**Spalte** ⇒ Exit 2): geschlossen.** Fixture `fx-e`
  (`| F-2 | <!-- offen --> |`, N Zellen): beide Images `2 Anforderung(en), 1
  Waise(n)`. `cellCountOK` greift erst bei N+1; die N-Zeile läuft in den
  `==`-Zweig und bleibt Zellinhalt.
- **R1-F-3 (Suite pinnt die defekte Platzierung): geschlossen und invertiert.**
  Mutation M-A (Regel zurück in den Splitter) kippt jetzt drei Tests, darunter den
  Konsumententest `TestTraceTableHeaderMitDirektiveWirdErkannt` — die von R1
  geforderte Richtung. Die sichere Platzierung wird nicht mehr abgelehnt.
- **Regression gegen v0.45.1 über die Fixture-Matrix:** geprüft. Byte-gleiches
  Verhalten in `fx-a` (sauber), `fx-c` (Header-Direktive), `fx-d` (Header +
  N-Trennzeile), `fx-e` (Kommentar-Spalte), `fx-f` (Trennzeile mit Marker), `fx-g`
  (N+2), `fx-h` (N−1), `fx-i` (Marker in erster Spalte), `fx-k` (nur
  Header+Trennzeile), `fx-l2` (Tabellen mit Leerzeile). Abweichungen **nur** in
  `fx-b`/`fx-j` (N+1-Datenzeile mit Direktive — die beabsichtigte Erweiterung,
  Exit 2 ⇒ gelesen) und in `fx-l`/`fx-x3` (F-1). Die Behauptung der
  Byte-Identität für Header-Fall und legitime Kommentar-Spalte ist damit
  **bestätigt**, als Universalaussage aber falsch (F-1).
- **XREF-Indizierung (`forwardEdges`/`backwardEdges`, out-of-range / falsche
  Spalte):** geprüft, ohne Befund. Beide lesen ausschließlich `row.cells[bt.primary]`
  / `[bt.secondary]`; beide Indizes stammen aus `bindCrossColumns` über
  `t.header` und sind damit `< len(header) <= len(cells)`. Kein Konsument iteriert
  über alle Zellen (`grep` über `.rows`: nur `trace_table.go:85`,
  `trace_cross.go:251/273`), die tolerierte Kommentar-Zelle wird also nie gelesen
  und kann keine Kante injizieren. Kein Panic-Pfad: `t.rows` wird nur bei
  bestandenem `cellCountOK` gefüllt. Empirisch: `fx-x1` (tolerierte Datenzeile in
  der Rück-Tabelle) hält die Kante `F-2 → C-2` sichtbar.
- **`bindCrossColumns`/`bindTableColumns` mit Kommentar-Spalte im Header:**
  geprüft, ohne Befund. Beide binden über exakte Namensgleichheit; die Zelle
  `<!-- d-check:ignore (…) -->` trifft keinen konfigurierten Rollen-Namen, erhöht
  nur `counts[…]` ihres eigenen Literals und bindet an keine Rolle — wie in ADR
  und Spec zugesagt (`fx-c`, `fx-x2` belegen es am Image).
- **`htmlCommentCell`-Regex:** geprüft, ohne Befund — unverändert gegenüber
  `44a5201`, R1 hat ihn in der gefährlichen Richtung durchgespielt (kein `-->` im
  Body, keine echte Zelle matcht fälschlich). Neu ist nur der **Aufrufort**;
  `MatchString(cells[len(cells)-1])` ist durch den `len(cells) == len(header)+1`-Zweig
  gegen `len(cells) == 0` geschützt (Header ist nie leer).
- **Marker in der ersten Spalte / leere Zellen / Tabelle nur aus Header+Trennzeile:**
  geprüft, ohne Befund. `fx-i` (`| <!-- x --> | F-2 | Beta |`): fail-closed, die
  Toleranz ist strikt end-verankert. `fx-j` (`| F-2 |  | <!-- x -->`): toleriert,
  Titel leer — dieselbe Semantik wie `| F-2 |  |` unter v0.45.1. `fx-k`: Tabelle
  ohne Datenzeilen, beide Images identisch.
- **Trennzeile mit Marker:** geprüft, ohne Befund — `fx-f`: `isTableDelimiter`
  verwirft die Kommentar-Zelle, die Tabelle wird nicht erkannt; identisch in
  v0.45.1 (der Splitter ist unverändert) und GFM-konform.
- **Realdatenbeleg grid-gym:** geprüft, ohne Befund — `spec/architecture.md:913`
  ist eine Datenzeile mit vier Zellen bei 3-Spalten-Header. Gegen die Kopie der
  echten Dateien: v0.45.1 `error: trace.cross-consistency.backward: Tabellenzeile
  913 hat 4 statt 3 Zellen` (Exit 2), HEAD `174 Differenz(en)` (Exit 0) — der
  motivierende Fall ist gelöst, der Reader liest die Datei durch. (grid-gym nur
  gelesen, nichts geändert.)
- **„Defekt-Fix, kein CR":** geprüft, weiterhin **haltbar**. Das Lastenheft
  ([`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen))
  sagt keine Zellenzahl zu; die Regel lebt in Rang 2 (Spec) und ist dort
  fortgeschrieben. Unverändert gegenüber R1.
- **SemVer-Patch:** geprüft, haltbar **mit** Einschränkung. Die Bestandsregression
  aus R1-F-2 ist weg; gegenüber v0.45.1 ist die Änderung für alle
  GFM-renderbaren Eingaben eine reine Erweiterung (Exit 2 ⇒ lesbar). Die einzige
  Verengung ist der degenerierte Pfad aus F-1. Patch bleibt vertretbar, sobald die
  Universal-Zusage in ADR/Spec die tatsächliche Reichweite nennt.
- **ADR-Immutabilität ([`AGENTS.md`](../../AGENTS.md) §3.5):** geprüft, ohne
  Befund — ADR-0040 steht auf `Proposed`; die Regel bindet erst `Accepted`. Die
  Neufassung ist zudem in `## Geschichte` mit Anlass, verworfener Prämisse und
  Konsequenz protokolliert — vorbildlich nachvollziehbar. Bestehende
  `Accepted`-ADRs unangetastet.
- **ADR-0040 ↔ Spec ↔ Code ↔ Handbuch (Wortlaut-Konsistenz):** geprüft, ohne
  Widerspruch **zwischen diesen vieren** — alle vier sagen „Header GFM-streng,
  Datenzeile genau eine ganzzellige Kommentar-Zelle". Der Widerspruch liegt zum
  **Slice-Plan** (F-2) und in der Reichweite der Zusage (F-1).
- **Handbuch §5 gegen das Binary (positive Richtung):** geprüft, ohne Befund —
  „Datenzeile toleriert genau eine" ist mit `fx-b`/`fx-g`/`fx-i` belegt, „Header
  braucht N+1-Trennzeile" mit `fx-c`. Die Einschränkung betrifft nur die
  Exit-2-Enumeration (F-4). Der §4.12-nahe Verweis auf §5 (`:628-631`) und die
  §7-Marker-Konvention (`:1229`) sind konsistent; das Handbuch sagt ehrlich, dass
  die Zeile „gelesen wird wie ohne Marker" — der Marker nimmt sie nicht aus der RTM.
- **ADR-0005-Import-Regeln / Modul-Layout:** geprüft, ohne Befund — keine neuen
  Imports, die Änderung bleibt in `core/app`; `dropCommentSuffix` ersatzlos
  entfernt.
- **[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  (Netz/Seiteneffekte):** geprüft, ohne Befund — nur `regexp`/`strings`, alle
  Läufe mit `--network none`.
- **Gate-Suppression / Schwellen-Senkung ([`AGENTS.md`](../../AGENTS.md) §3.6):**
  geprüft, ohne Befund — keine `//nolint`, keine Gate-/Schwellen-Änderung.
- **Referenz-Richtung (SDP), Marker-Ehrlichkeit:** geprüft, ohne Befund —
  ADR-0040 nennt slice-074 weiterhin nur in `## Geschichte` als Provenance.
- **CHANGELOG / Release-Prep fehlen im Diff:** geprüft, **kein Finding** — die
  Commit-Grenzen-Konvention legt sie in den Release-Prep-Commit; der Review läuft
  vertragsgemäß davor. (Das Handbuch kam mit `990dc68` bereits vorab.)
- **DoD-Abhakung / Gate-Lauf-Bestätigung (`make gates`):** nicht geprüft — Rolle
  der Verifikation, nicht des Reviews (Reviewer-Skill §Anti-Pattern). F-2 betrifft
  den **Wortlaut** der DoD, nicht ihre Abhakung.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 2 | F-1, F-2 |
| LOW | 4 | F-3, F-4, F-5, F-6 |
| INFO | 1 | F-7 |

---

## Verdikt

**BLOCK — v0.45.2 darf noch nicht getaggt werden.** Der Abstand zu Grün ist
allerdings klein und liegt überwiegend **außerhalb** des Produktivcodes.

Die Ansatz-Umkehr ist richtig und gut ausgeführt. R1-F-1, R1-F-2 und R1-F-3 sind
geschlossen — nachgefahren gegen das echte Image, nicht bloß gelesen: `fx-c` und
`fx-x2` liefern wieder die Befunde, die `44a5201` verschluckte; `fx-e` liest die
Kommentar-Spalte wie v0.45.1; die Mutation „Regel zurück in den Splitter" kippt
jetzt drei Tests statt keinen. Dass `splitPipeTableLine` gegenüber dem
ausgelieferten Stand byte-identisch ist, macht die R1-F-1-Mechanik strukturell
unmöglich — diese Zusage der ADR trägt.

**Was blockiert:** ADR-0040 verallgemeinert diese echte Teil-Aussage zu „ein
stiller Übersprung ist **strukturell unmöglich**" und „gegenüber v0.45.1 wird
**keine** Zeile anders gelesen". Beides ist mit `fx-l`/`fx-x3` widerlegt: die
Toleranz verlängert die Tabelle über eine vormals beendende Zeile hinweg, ein
nachfolgender Direktiven-Header wird verschluckt, Waisen und Rück-Kanten
verschwinden, Exit 1 wird Exit 0 — wieder in beiden Konsumenten (F-1). Die
Vorbedingungen sind eng und die auslösende Form ist nicht GFM-renderbar; das
**Verhalten** ist vertretbar und teils sogar GFM-treuer als v0.45.1. Nicht
vertretbar ist die als Beweis formulierte Universal-Zusage in einem normativen
Dokument — in genau der Code-Region, in der dieselbe Erschließung („einen Zweig
geprüft, auf die Klasse geschlossen") schon zweimal ein stilles Grün erzeugt hat.
F-1 ist damit **entweder** durch eine Verengung im Code **oder** durch eine
ehrliche Reichweiten-Angabe zu entlasten; beides ist eine Entscheidung des Autors,
nicht des Reviews.

F-2 blockiert unabhängig davon: der Slice-Plan verlangt in §3 DoD weiterhin die
Splitter-Platzierung, die ADR-0040 als Defekt verwirft und die R1 als Ursache des
stillen Grüns zitiert hat. ADR, Spec, Code und Handbuch wurden nachgezogen, der
Slice nicht — die DoD ist gegen den ausgelieferten Code nicht abhakbar, und ein
Verifier, der ihrem Wortlaut folgt, forderte den Defekt zurück. F-3 bis F-7 sind
nicht release-blockierend.

**Steering-Loop-Hinweis (Reviewer-Skill §Kontext-Eskalation).** Das ist die
**vierte** Runde derselben Klasse in dieser Code-Region. Sie ist diesmal nicht im
Code gelandet, sondern in der Zusage **über** den Code: nachdem die konkrete
Mechanik geschlossen war, wurde die Klasse als geschlossen deklariert, statt sie
zu prüfen — die Fitness-Funktion der ADR listet den Header-Fall, den
Datenzeilen-Fall und die Kommentar-Spalte, aber keinen Fall, in dem die
**Toleranz selbst** eine Zeile umdeutet. Der von slice-074 §4 selbst benannte
Ansatzpunkt (Dogfood-Lücke, Konsumenten-naher Fixture-Anker) und die
konsumentenseitige Lücke aus F-5 zeigen auf dieselbe Stelle. Der Wert des
Verfahrens hat sich hier bestätigt: der Review **vor** dem Release hat den ersten
Defekt gefangen, und die Umkehr hat ihn strukturell beseitigt.

Der motivierende Realfall (grid-gym `architecture.md:913`) wird korrekt gelöst —
v0.45.1 bricht dort ab, HEAD liest die echte Datei durch.
</content>
</invoke>
