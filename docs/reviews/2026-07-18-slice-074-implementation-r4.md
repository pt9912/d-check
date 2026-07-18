# Review slice-074 (R4) — Direktiven-Zelle in Tabellenzeilen (neue, sichere Fassung)

**Datum:** 2026-07-18 · **Rolle:** unabhängiger Reviewer (kontext-getrennt,
nicht der Autor) · **Lauf:** R4 — die neue Fassung, nachdem R1/R2/R3 der alten
Fassung 2026-07-17 mit BLOCK endeten und der Slice zurückgenommen wurde. Diese
Fassung setzt als **Aufsatz** auf slice-077/[ADR-0043](../plan/adr/0043-tabellengrenze-am-relevanten-header.md)
(Grenze am relevanten Header, released v0.48.0) auf. **Vor** Release v0.48.1.

**Gegenstand:** `b0d1a4f` (feat), Range `5993ea3..HEAD` (origin/main = `5993ea3`;
davor die doc-first-/Aktivierungs-Commits `e944be0`/`f6e5189`/`1572831`). Kern:
`internal/hexagon/core/app/trace_table.go` (`htmlCommentCell` + Toleranz im
`len(cells) != len(t.header)`-Zweig), Tests in
`internal/adapter/driving/cli/cli_acceptance_test.go` (`TestCLI074_*`).

**Quellen:** [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 5, [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency),
[ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md),
[ADR-0043](../plan/adr/0043-tabellengrenze-am-relevanten-header.md),
[ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md),
[slice-074](../plan/planning/in-progress/slice-074-kommentar-suffix-tabellenzeilen.md),
[R3](2026-07-17-slice-074-implementation-r3.md), [`AGENTS.md`](../../AGENTS.md) §3.

**Verifikations-Basis:** `make test` (Baseline **grün**) und `make build`
(Image `sha256:78944c5d9114…`) in einem **Scratch-Worktree** (nach dem Lauf
entfernt; Produktivbaum `git status` leer, HEAD `b0d1a4f`). Läufe
`docker run --rm --network none -v <fixture>:/repo:ro -w /repo <image> --trace
[--require-complete] [--json]` gegen die **ausgelieferten** Images
`ghcr.io/pt9912/d-check:v0.47.0` (Vor-Grenze) und `:v0.48.0` (Grenze released,
korrekte Patch-Basis) sowie HEAD. Mutation via `if false && …` in der
Worktree-Kopie, `make test`. Realdaten: **Kopie** von grid-gyms
`spec/architecture.md` (grid-gym selbst nur gelesen, `git status` danach leer).

**Vorbemerkung.** Der strukturelle Kern der fünf Fehlschläge — das Tolerieren
entfernt den `badLine`-Wiederaufsetz-Punkt und verschluckt die Folgetabelle
(R3-F-1) — ist in dieser Fassung **nicht mehr Sache der Toleranz**: die
Tabellengrenze zieht seit slice-077 der geteilte Reader am relevanten Header
(`consumeTableRows:180-184`, **vor** dem Zellenzahl-Zweig). Ich habe den
R3-F-1-Fall gegen HEAD **direkt nachgebaut** und ihn geschlossen gefunden (siehe
Negativbefunde). Die beiden Findings unten sind Nits, kein BLOCK.

---

## Findings

### F-1 · LOW · Der historisch fünfmal gescheiterte Kombinationsfall (tolerierte Direktiven-Datenzeile *unmittelbar vor* einer relevanten Folgetabelle) ist nur **transitiv** gepinnt — kein lokaler Regressionstest in der `TestCLI074_*`-Suite

**quelle:** [slice-074](../plan/planning/in-progress/slice-074-kommentar-suffix-tabellenzeilen.md)
§3 DoD („Mutations-Härte … gemessen, nicht zugesagt" — die R3-F-2-Lehre);
[ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md) Fitness-Funktion
(„Folgt der laufenden Tabelle **ohne Leerzeile** eine **relevante** Folgetabelle,
wird deren Header **nicht** als tolerierte Datenzeile verschluckt … Gemessen
(`fx-m`: `total 3`)"); Reviewer-Anker MEDIUM/LOW (Sensor-Lücke an der Stelle, an
der die Klasse zuletzt fünfmal zuschlug).

**pfad:** `internal/adapter/driving/cli/cli_acceptance_test.go:2409-2495`
(die vier `TestCLI074_*`).

**befund:** Die vier 074-Tests decken die Toleranz in **einer** Tabelle
(`DirektivenDatenzeileToleriert`: R-1/R-2-Direktive/R-3 in derselben Tabelle),
den Guard (`ZweiExtraZellen`), den Panic-Pfad (`KommentarLetzteZeile`) und den
Cross-Pfad ab — **nicht** den exakten R3-F-1-Fall: eine **tolerierte**
Direktiven-Datenzeile, der **ohne Leerzeile** eine **relevante Folgetabelle**
folgt. Genau dieser Fall (`fx-m` der ADR) ließ den Slice fünfmal fallen. Die
ADR-Fitness-Funktion behauptet ihn als „gemessen (`fx-m`: `total 3`)", aber im
Repo existiert dafür kein committeter Test — `grep -n "fx-m\|Folgetabelle"
cli_acceptance_test.go` trifft nur die slice-077-Tests. Der Fall ist **transitiv**
gepinnt (die Toleranz durch die drei 074-Tests, die Grenze durch
`TestCLI077_StillerUebersprungGrenze`), und ich habe verifiziert, dass **keine
einzelne** Mutation den stillen Übersprung geräuschlos zurückbringt: das Entfernen
der Grenze kippt den 077-Test, das Deaktivieren der Toleranz kippt drei 074-Tests.
Ein Versagen entstünde erst, wenn ein künftiger Refactor **Grenze und** die
Reihenfolge (Grenze vor Zellenzahl-Zweig) zugleich anfasst; dann fehlt der lokale
Wächter genau für die 5×-Klasse. Kein Korrektheitsbefund am Code — eine
Härtungslücke am Sensor.

**verifizierbar:** ja — `docker run … -v fx-m-rel:/repo:ro … --trace
--require-complete` liefert unter HEAD `total 3`/Exit 1 (kein Verlust); ein
committeter Test dieser Form bände die Klasse lokal.

---

### F-2 · INFO · `htmlCommentCell`-Regex `^<!--.*-->$` ist greedy — eine überzählige Zelle mit *sichtbarem* Text zwischen zwei Kommentar-Markern gilt als tolerierbar und wird stumm verworfen

**quelle:** [ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md)
Entscheidung 1 („nur, wenn sie **ganzzellig** ein HTML-Kommentar ist");
dokumentationswürdige Annahme.

**pfad:** `internal/hexagon/core/app/trace_table.go:20-26`
(`htmlCommentCellRe`/`htmlCommentCell`).

**befund:** `^<!--.*-->$` mit greedy `.*` matcht auch Zellen, die **nicht** ganz
ein Kommentar sind: `<!-- x --> VISIBLE <!-- y -->` und `<!-- a --><!-- b -->`
(mehrere Kommentare) werden beide als tolerierbar gewertet (gemessen: `total 2`,
Exit 0). Da die tolerierte Zelle immer die **überzählige** (N+1.) ist und per
`cells[:len(t.header)]` verworfen wird — und GFM überzählige Body-Zellen ohnehin
ignoriert — geht **keine** gelesene Anforderungs-/Titel-/Kanten-Zelle verloren;
das Ergebnis ist byte-identisch zur reinen Kommentar-Form. Die Härte des Guards
bleibt in der gefährlichen Richtung erhalten: eine überzählige Zelle, die **nicht**
mit `-->` schließt (`foo -->`, `<!-- x --> tail`), bleibt Exit 2 (gemessen). Es
ist also kein Korrektheits-, sondern ein Ehrlichkeits-Detail: „ganzzellig
Kommentar" ist in der Umsetzung „beginnt mit `<!--`, endet mit `-->`", nicht
„enthält ausschließlich Kommentar-Syntax". Weil die Zelle verworfen wird, ist das
folgenlos — aber die Annahme steht weder im Code-Kommentar noch in der Spec.

**verifizierbar:** ja — `-v fx-sandwich:/repo:ro … --trace --json` liefert
`total 2`, Exit 0; die überzählige Zelle taucht in keiner Ausgabe auf.

---

## Negativbefunde (geprüft, ohne Befund)

- **R3-F-1 (der 5×-Fehlschlag: tolerierte Datenzeile verschluckt die Folgetabelle)
  — GESCHLOSSEN, gegen HEAD direkt nachgebaut.** Fixture `fx-m-rel` (tolerierte
  Direktiven-Datenzeile R-1, **ohne Leerzeile** gefolgt von einer relevanten
  Tabelle mit R-2/R-3): HEAD `total 3`/Exit 1 — nichts verschluckt. Die Grenze
  (`trace_table.go:180-184`) feuert am relevanten Header **vor** dem
  Zellenzahl-Zweig; die tolerierte Zeile entfernt den Wiederaufsetz-Punkt **nicht**
  mehr, weil die Grenze ihn ersetzt. `fx-m-irrel` (tolerierte Zeile + irrelevante
  Tabelle + relevante Tabelle): `total 2`/Exit 1 — die relevante Folgetabelle wird
  trotz dazwischenliegender irrelevanter erkannt. Fail-closed intakt.
- **Toleranz-Korrektheit (genau N+1, ganzzellig Kommentar, Truncation):**
  geprüft, ohne Befund — die Bedingung `len(cells) == len(t.header)+1 &&
  htmlCommentCell(cells[len(cells)-1])` (`trace_table.go:194`) greift nur bei genau
  einer überzähligen Zelle; `cells[:len(t.header)]` behält ID- und Textspalte
  korrekt (`fx-mreq`: die tolerierte Zeile **ist** die Anforderung, Titel „only
  one", der Kommentar landet nicht im Titel). Stimmt mit Spec Schritt 5 und
  ADR-0040 Entscheidung 1/3 überein.
- **Guard scharf (echte Verrutschung bleibt Exit 2):** geprüft, ohne Befund —
  zwei überzählige Zellen (`| … | extra | <!-- x -->`), überzählige
  Nicht-Kommentar-Zelle (`foo -->`, `<!-- x --> tail`) und ein Kommentar in der
  **Mitte** mit Datenzelle dahinter bleiben Exit 2. Der ADR-0037-Nullmengen-/
  Guard-Vertrag ist nicht aufgeweicht.
- **Panic/Bounds (R3-M-4):** geprüft, ohne Befund — `cells[len(cells)-1]` ist
  durch `len(cells) == len(t.header)+1` und `len(t.header) >= 1` (`tableHeaderAt`
  ⇒ `isTableDelimiter` verlangt `len > 0`) auf `>= 2` untergrenzt; kein leerer
  Slice. Kommentar als **letzte** Zeile ohne Trailing-Newline: HEAD `total 2`,
  Exit 0, kein Panic (`TestCLI074_KommentarLetzteZeileKeinPanic`, Mutation kippt
  ihn). Die Grenze (`trace_table.go:180`) ist mit `j+1 < len(lines)` gegen den
  Zugriff auf `lines[j+1]` geschützt.
- **Beide Konsumenten über den geteilten Reader:** geprüft, ohne Befund — die
  Toleranz wirkt auf `format: table` (`fx-m-rel`), auf `cross.backward`
  (`TestCLI074_Cross`, grid-gym-Realfall) **und** auf `cross.forward`
  (eigene Fixture: v0.47.0 `forward: Tabellenzeile 5 hat 3 statt 2 Zellen`/Exit 2
  → HEAD `0 Differenz(en)`/Exit 0). `consumeTableRows` ist der einzige Lese-Kern
  für alle drei.
- **Mutations-Härte (Toleranz gepinnt — R3-F-2 adressiert):** geprüft, ohne
  Befund — `if false && …` vor der Toleranz kippt in `make test` **genau**
  `TestCLI074_DirektivenDatenzeileToleriert`, `…KommentarLetzteZeileKeinPanic`,
  `…Cross_DirektivenZeileToleriert`; `…ZweiExtraZellenBleibtExit2` bleibt grün;
  die slice-070/077-Suite bleibt grün. Die R3-F-2-Lücke (ungepinnte Grenzen der
  alten Fassung) besteht in dieser Fassung nicht — die Toleranz selbst ist
  gepinnt, die Grenze durch die 077-Suite.
- **SemVer-Patch (rot→grün, kein grüner Lauf wird rot):** geprüft, ohne Befund —
  gegen die **korrekte** Patch-Basis v0.48.0 (Grenze bereits enthalten): der
  volle `--trace`-Output auf grid-gyms committeter Config ist zwischen v0.48.0 und
  HEAD **byte-identisch** (`diff` leer). Die einzige Verhaltensänderung ist die
  N+1-Direktiven-Datenzeile, die von Exit 2 nach lesbar wechselt; keine heute
  gelesene Zeile wird anders gelesen.
- **Realdatenbeleg grid-gym `architecture.md:913`:** geprüft, ohne Befund —
  gegen eine **Kopie** der echten Datei (Header `| Testart | Verortung | Bezug |`,
  Zeile 913 mit `<!-- d-check:ignore (geplant: …) -->` als 4. Zelle), Rück-Sicht
  über `cross-consistency.backward`: v0.47.0 `error: … backward: Tabellenzeile 913
  hat 4 statt 3 Zellen`/Exit 2 → HEAD Exit 0. Der motivierende Realfall läuft
  durch. grid-gym `git status` danach leer.
- **`htmlCommentCell`-Kanten (leer, Whitespace, Nicht-Kommentar):** geprüft, ohne
  Befund — `<!---->` (leer) toleriert; Whitespace um den Kommentar via
  `strings.TrimSpace` toleriert; `foo -->`/`<!-- x --> tail` (kein ganzzelliger
  Kommentar) Exit 2. Die greedy-`.*`-Beobachtung ist F-2 (folgenlos).
- **Legitime N-spaltige Zeile fälschlich toleriert:** geprüft, ohne Befund — die
  Toleranz betritt nur den `len(cells) != len(t.header)`-Zweig; eine N-spaltige
  Zeile (Zellenzahl == Header) läuft nie hinein. Ein Header, der die Direktive als
  echte Spalte trägt (N+1 gegen N+1-Trennzeile), wird regulär gebunden — kein
  Toleranz-Pfad.
- **Reihenfolge Grenze vs. Toleranz:** geprüft, ohne Befund — die Grenze
  (`:180-184`) steht **vor** dem Split und dem Zellenzahl-Zweig (`:185-200`). Für
  den Header einer relevanten Folgetabelle feuert die Grenze, bevor die Toleranz
  ihn je als Datenzeile sehen könnte; die Toleranz kann keinen relevanten Header
  verschlucken.
- **Masken-Kopplung (R3-F-5):** geprüft, ohne Befund — die Grenze prüft jetzt
  `maskAllows(mask, lines[j+1].no)` (`:180`), bevor sie `tableHeaderAt` aufruft;
  die in R3-F-5 benannte Masken-Blindheit von `isNewTableHeader` ist mit dem
  neuen Grenz-Prädikat nicht mehr gegeben.
- **XREF-Kanten-Injektion über die tolerierte Zelle:** geprüft, ohne Befund —
  `forwardEdges`/`backwardEdges` lesen nur die gebundenen Spaltenindizes
  (`< len(header) <= len(cells)`); die verworfene Kommentar-Zelle wird nie
  gelesen und kann keine Kante erzeugen (`fx`-Forward-/Backward-Läufe:
  `0 Differenz(en)`).
- **`splitPipeTableLine` byte-identisch:** geprüft, ohne Befund — der Diff fügt
  nur `htmlCommentCell` und die Toleranz-Verzweigung hinzu; der Splitter entfernt
  nichts (keine `dropCommentSuffix`-Wiederkehr, `grep` leer). R1-F-1/F-2/F-3
  bleiben geschlossen.
- **ADR-Immutabilität ([`AGENTS.md`](../../AGENTS.md) §3.5):** geprüft, ohne
  Befund — ADR-0040 steht auf `Proposed` (Wiederaufnahme erlaubt), Historie
  vollständig; ADR-0043 (`Accepted`) ist im Range **unangetastet** (`git diff …
  -- 0043….md` leer). Die Änderungen an den R2/R3-Reviews sind **reine
  Link-Pfad-Aktualisierungen** (`open/` → `in-progress/`), kein Eingriff in
  Befunde.
- **ADR-0005-Import-Regeln / Modul-Layout:** geprüft, ohne Befund — keine neuen
  Imports, alles in `core/app`; `regexp`/`strings` bereits vorhanden.
- **[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit):**
  geprüft, ohne Befund — alle Läufe `--network none`; kein Netz-, kein
  Schreibzugriff.
- **Gate-Suppression / Schwellen ([`AGENTS.md`](../../AGENTS.md) §3.6):** geprüft,
  ohne Befund — keine `//nolint`, keine Gate-/Schwellen-Änderung.
- **Referenz-Richtung (SDP), Marker-Ehrlichkeit:** geprüft, ohne Befund — ADR-0040
  nennt slice-074 nur in `## Geschichte`/Kopf als Provenance, nicht als
  Entscheidungsgrundlage; kein `status-provenance`-Marker missbraucht.
- **Spec-Treue Schritt 5:** geprüft, ohne Befund — der ergänzte Absatz
  (`spezifikation.md:398-405`) beschreibt genau „genau eine überzählige,
  ganzzellige HTML-Kommentar-Zelle … auf Header-Breite gelesen; zwei überzählige
  Zellen oder eine Nicht-Kommentar-Zelle bleiben Exit 2" — deckt sich mit dem Code.
- **CHANGELOG / Release-Prep / DoD-Abhakung / `make gates`:** nicht geprüft —
  Release-Prep und Verifikation sind nicht Reviewer-Rolle (Reviewer-Skill
  §Anti-Pattern). Der Feat-Commit trägt sie konventionsgemäß nicht.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 1 | F-1 |
| INFO | 1 | F-2 |

---

## Verdikt

**ACCEPT-WITH-NITS — v0.48.1 kann getaggt werden.**

Die Klasse, die diesen Slice fünfmal fallen ließ, ist geschlossen — und zwar an
der **Wurzel**, nicht am Nachbarfall: die Tabellengrenze zieht seit slice-077 der
geteilte Reader am relevanten Header (`consumeTableRows:180-184`), **vor** dem
Zellenzahl-Zweig, in dem die Toleranz sitzt. Ich habe den exakten R3-F-1-Fall
gegen HEAD nachgebaut (`fx-m-rel`: tolerierte Direktiven-Datenzeile, ohne
Leerzeile gefolgt von einer relevanten Tabelle) und ihn **geschlossen** gefunden
(`total 3`, kein Verlust), ebenso den Cross-Spiegel (forward **und** backward).
Der Guard bleibt scharf (zwei Zellen, Nicht-Kommentar, mittiger Kommentar → Exit
2), der Panic-Pfad ist geschützt und gepinnt, und die R3-F-2-Sensor-Lücke ist
weg: das Deaktivieren der Toleranz kippt **gemessen** drei 074-Tests. Der
SemVer-Patch trägt gegen die korrekte Basis v0.48.0 (byte-identischer
grid-gym-Output; nur die N+1-Direktiven-Zeile wechselt rot→grün), und der
motivierende Realfall `architecture.md:913` läuft unter HEAD durch, wo v0.47.0
mit „913 hat 4 statt 3 Zellen" abbrach.

Es bleiben zwei Nits, beide nicht release-blockierend. **F-1 (LOW):** der
historisch fragile Kombinationsfall ist nur **transitiv** gepinnt (Toleranz durch
074-Tests, Grenze durch `TestCLI077_StillerUebersprungGrenze`) — korrekt und
gegen jede Einzel-Mutation abgesichert, aber ohne lokalen Regressionstest für die
5×-Klasse; ein `fx-m-rel`-Test würde sie hart binden. **F-2 (INFO):** die
greedy-`^<!--.*-->$`-Regex akzeptiert überzählige Zellen mit sichtbarem Text
zwischen zwei Kommentar-Markern — folgenlos, weil die überzählige Zelle ohnehin
verworfen wird, aber „ganzzellig Kommentar" ist in der Umsetzung schwächer als in
der ADR-Prosa. Beide gehören in die Übergabe, nicht in den Weg zum Tag.
