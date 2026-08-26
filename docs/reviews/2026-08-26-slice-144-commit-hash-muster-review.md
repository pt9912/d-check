# Review-Report: slice-144 — Commit-Hashes in den Spec-Straten, ein Muster mit vertretbarer Falsch-Positiv-Last

**Datum:** 2026-08-26 · **Review-Art:** Code-/Konfigurations-Review (Config- und
Hard-Rule-Diff gegen Slice-Plan, Kanon, Spezifikation und Modul-Code) ·
unabhängiger Reviewer ohne Anteil an der Arbeit
**Gegenstand:** Commit-Kette `7dc8738..8e0a6f0` von slice-144 — `7dc8738`
(Anspruchs-/Lifecycle-Move), `89363e3` (nur die Schwelle, +15 Zeilen im
Slice-Plan), `8e0a6f0` (Feat: `.d-check.yml` +27, `AGENTS.md` +18/−13)
**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 (`9ee805b`) · **Modell-ID:**
`claude-opus-5[1m]`

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-144-commit-hash-muster.md`
  vollständig — insbesondere §2 Schritt 0 (die vorab notierte Schwelle), §3
  „Ausdrücklich NICHT", §5 Risiken, §7 Vorgelagert
- `.d-check.yml`, Block `matrix` vollständig (alle fünf Klassen, `rules`,
  `exclude-sections`, `exempt-paths`) sowie `ids`/`codepaths`/`versions`/
  `structure` zur Abgrenzung
- `AGENTS.md` §3.4 vorher (`git show 8e0a6f0^:AGENTS.md`) und nachher, dazu
  §3.3, §3.7, §3.8, §4, §5
- `spec/lastenheft.md` §`DC-FA-MTX-001`, §`DC-FA-MTX-003`, §`DC-FA-CLI-003`;
  `spec/spezifikation.md` §`DC-FA-MTX-001.a` Schritt 6 und die
  Schema-Tabelle §2 (`matrix.classes[].name`/`.paths`/`.order`/`.direction`/
  `.token`)
- Quellcode: `internal/hexagon/core/rules/matrix.go` (`CheckMatrix`,
  `tokenFindings`, `linkSpanRe`, `provenanceMarker`),
  `internal/hexagon/core/rules/markdown.go` (`proseLines`),
  `internal/hexagon/core/model/finding.go` (`SortFindings`, Dedup),
  `internal/adapter/driven/configyaml/configyaml.go` (`applyMatrix`),
  `internal/hexagon/core/model/config.go` (`MatrixClass`)
- `harness/conventions/MR-006-referenzrichtung-matrix.md` und
  `harness/conventions/MR-034-matrix-scope-welle.md`
- Baseline `v5.12.0`: `grundlagen-referenz-richtung.md` §Referenz-Richtung
  (SDP) und `grundlagen-harness-dateien.md` §Was ein Kommentar trägt
- `docs/plan/planning/observations.md` — `BEO-009`, `BEO-011`, `BEO-012`,
  `BEO-017`
- **Vorherige Findings am gleichen Modul (Pflicht):**
  `docs/reviews/2026-08-25-slice-138-matrix-wellen-klasse-review.md` (F-1 HIGH
  „gedeckt trägt weniger als der Wortlaut", F-3 MEDIUM „`paths` ist für eine
  Klasse nicht erforderlich") und
  `docs/reviews/2026-08-26-slice-150-pin-gebundene-zitate-review.md`
  (Finding-Klassen `zahl-ohne-messvorschrift`,
  `botschaft-behauptet-eine-nicht-stattgefundene-nicht-aenderung`,
  `probe-ohne-exit-und-ohne-gegenstand`)

**Nicht erhalten** (Skill §Eingangs-Kontext): die DoD-Abhakung. Die vier
offenen Haken in slice-144 §4 werden hier **nicht** bewertet; der Slice liegt in
`in-progress/`.

**Vom Reviewer selbst gefahren** (nur lesend; jede Sonde in einem eigenen
Scratchpad-Repo bzw. als schreibgeschützter Config-Overlay über den
read-only-Mount — **keine Repo-Datei wurde verändert**, `git status --short`
vor und nach jedem Lauf leer, kein `git checkout --` nötig):

- `make doc-check` ⇒ **Exit 0**, „515 Datei(en) geprüft, 0 Befund(e)"; das
  Image wurde dabei aus HEAD neu gebaut (`sha256:f41f99ad9a1e…`) und ist die
  Basis aller folgenden Läufe.
- **Messung mit dem Produkt** (`docker run --network none`, Repo read-only,
  ein zweiter Mount legt eine Sonden-Config über `.d-check.yml`;
  `modules: [matrix]`, `exclude-sections: [Geschichte]`, `exempt-paths` wie im
  Repo, eine Quellklasse plus eine Token-Klasse mit verbotener Kante):
  - Quellklasse = die drei Straten: `\b[0-9a-f]{7,40}\b` ⇒ 0 · `\b[0-9a-f]{40}\b`
    ⇒ 0 · `\b[0-9a-f]{7,12}\b` ⇒ 0 · `(?i)commit\s+[0-9a-f]` ⇒ **2**
    (`spec/lastenheft.md:2032`, `:2063`) · `(?i)commit\s+[0-9a-f]{7,40}` ⇒ 0.
  - Quellklasse = alle Markdown-Dateien (`["**/*.md","*.md"]`):
    `\b[a-f]{7,40}\b` ⇒ **9** · `\b[0-9a-f]{7,40}\b` ⇒ **1146** ·
    `\b[0-9]{7,40}\b` ⇒ **104**. Variante ohne `exclude-sections`/`exempt-paths`:
    1151 bzw. 9.
  - Quellklasse = ADR-Körper: `\b[0-9a-f]{7,40}\b` ⇒ **0** (ohne `exempt-paths`
    genau **1**: `docs/plan/adr/0010-semgrep-hermetisches-gate.md:34`, ein
    Regelset-Pin).
- **Messung über den Rohtext** (`python3`/`re` über dieselben 515 gescannten
  `.md`, ohne Fence-, Link- oder Sektions-Filter): `\b[a-f]{7,40}\b` ⇒ **10**
  (ausnahmslos `deadbeef`) · `\b[0-9a-f]{7,40}\b` ⇒ **1191** (in den drei
  Straten: **0**) · `\b[0-9]{7,40}\b` ⇒ **107**. Dieselben drei Zahlen an den
  Bäumen `89363e3` und `5272d5c` (`git archive`) — unverändert.
- **Falsch-Negativ-Sonde:** ein Spec-Stratum mit 25 Schreibweisen desselben
  Hashes (bare kurz/lang, Inline-Code, Markdown-Link mit Hash im Text, Markdown-Link
  mit Hash in der URL, bare URL, Autolink, HTML-Kommentar, Tabellenzelle,
  Satzzeichen, Klammern, Großschrift, Mischschrift, Zeilenumbruch, voller
  `sha256`-Digest, gekürzter Digest mit und ohne Endstück, `@sha256:`-Pin,
  Provenance-Marker-Zeile, Code-Fence, eingerückter Code, Sechssteller,
  `git show`-Kommando, Referenz-Link-Definition, Bindestrich-Trennung) gegen die
  reale Token-Form ⇒ 18 Befunde, sieben Formen stumm.
- **Falsch-Positiv-Sonde:** UUID, Lauf-Nummer, Kompaktdatum, Byte-Grenze,
  Hex-Farbwert, Docker-Kurz-ID, `defaced`, `deadbeef` ⇒ **9 von 9** gemeldet.
- **Positivkontroll-Gegenprobe:** kurz/lang/Inline-Code gegen alle vier
  Kandidaten (bestätigt die Trennung, die die Botschaft behauptet).
- `git log --format='%h %ad %s' --date=iso`, `git show`/`--numstat` auf alle drei
  Commits, `git diff 89363e3 HEAD` auf die Slice-Datei (leer),
  `git log --follow` auf die Slice-Datei.

**Verdikt: blockierend** — ein HIGH, fünf MEDIUM, vier LOW, ein INFO.

---

## Findings

### F-1

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §3.7 (*„Keine Review-Historie …, keine Deliberation
  über Verworfenes, keine Herkunfts-Prosa, keine Slice-Nummern und keine
  Mess-Labels; Herkunft nur als **ein** auflösbares Feld nach dem
  Baseline-Schema (`DC-*` …, `ADR-*`, `MR-*`, `seit welle-<NN>`)"*, dazu
  *„Neuzugänge fallen überall unter den Anker"*) ·
  `.harness/baseline/v5.12.0/regelwerk/grundlagen-harness-dateien.md` §Was ein
  Kommentar trägt (*„Wer Herkunft nennt, nennt sie als **ein** auflösbares Feld
  … und nie als Absatz"*; Adressaten-Test: *„Der Entscheider hat ADR und
  Slice"*; *„Unzulässig ist er über die **verworfene Alternative**"*) ·
  Reviewer-Skill §HIGH-Anker (*Kommentar trägt keine der fünf Klassen*, seit
  1.5.0, Auflösungs-Trigger permanent)
- **pfad:** `.d-check.yml:246–251`
- **befund:** Der neue Klassen-Kommentar trägt eine Slice-Nummer als
  Herkunfts-Zeiger (*„DIE SCHWELLE STAND VOR DER MESSUNG (slice-144 §2 Schritt
  0)"*) und zwei Zeilen reine Lauf-Historie — *„Bestand heute 0 Befunde (vier
  Kandidaten, alle 0)"* nennt die drei verworfenen Kandidaten, *„Positivkontrolle
  Kurz- UND Langform gemeldet, auch in Inline-Code"* protokolliert das Ergebnis
  eines einmaligen Probelaufs; keine der fünf Klassen trägt das, und
  `slice-` steht nicht unter den zugelassenen Herkunfts-Formen. Die beiden
  Geschwister-Kommentare derselben Klassen-Liste lösen ihre Herkunft
  ausschließlich über `DC-FA-MTX-003` bzw. `MR-006`/`MR-034` auf
  (`.d-check.yml:214–219`, `:232–238`).
- **verifizierbar:** nein — kein Gate liest Kommentar-Inhalte (`make doc-check`
  ist auf genau diesem Stand Exit 0, 515 Dateien, 0 Befunde); die Prüfung ist
  das im Skill benannte Urteil. Belegt durch `git show 8e0a6f0 -- .d-check.yml`
  gegen `AGENTS.md` §3.7 und die Baseline-Stelle.
- **klasse:** `konfig-kommentar-traegt-mess-historie-und-slice-nummer`

### F-2

- **kategorie:** MEDIUM
- **quelle:** `BEO-009` Richtung (b) (*„die genannten Proben liefen und
  stimmen, aber der Schluss daraus gilt weiter als sie reichen"*) ·
  Reviewer-Skill §MEDIUM-Anker *Botschaft verallgemeinert über die Messung
  hinaus* · Vorbefund `zahl-ohne-messvorschrift`
  (`docs/reviews/2026-08-26-slice-150-pin-gebundene-zitate-review.md`
  Eingangs-Kontext) · `spec/spezifikation.md` §`DC-FA-MTX-001.a` Schritt 6
  (`matrix` liest den Prosa-Körper *„außerhalb Fenced-Code, außerhalb
  `exclude-sections` und ohne Markdown-Link-Spans"*)
- **pfad:** `.d-check.yml:248–255` (Kopfzeile *„Gemessen mit dem Produkt, nicht
  per grep:"* über der Zeile *„zehn Treffer"*); Commit-Botschaft `8e0a6f0`,
  Absatz *„DIE RISIKO-KLASSE IST GEMESSEN"* (*„1181 hash-artige Token"*)
- **befund:** Die beiden Risiko-Zahlen sind über den **Rohtext** gemessen, nicht
  über den Text, den `matrix` liest: über das Produkt gefahren liefert
  `\b[a-f]{7,40}\b` **9** statt zehn Befunde (die zehnte Fundstelle ist der
  zweite `deadbeef` derselben Zeile
  `docs/plan/planning/done/slice-047-print-mk-doctor-repair-help-digest.md:102`,
  den `SortFindings` als identisches Tupel dedupliziert) und
  `\b[0-9a-f]{7,40}\b` **1146** statt 1181; die 1181 ist rechnerisch der
  Rohtext-Wert 1191 minus die zehn `deadbeef`. Die Kommentar-Zeile steht damit
  unter einer Kopfzeile, die für sie nicht gilt, und die Botschaft, die
  „gemessen mit dem Produkt, nicht per grep" zu ihrem Kernanspruch macht, trägt
  eine Zahl, die genau so nicht entstanden sein kann.
- **verifizierbar:** ja — Sonden-Config mit `modules: [matrix]`, Quellklasse
  `["**/*.md","*.md"]`, Token `\b[a-f]{7,40}\b` bzw. `\b[0-9a-f]{7,40}\b` über
  den read-only-Mount; Gegenrechnung per `python3`/`re` über dieselben 515
  Dateien (oben unter *Vom Reviewer selbst gefahren* mit Vorschrift und Zahl).
- **klasse:** `zahl-ohne-messvorschrift`

### F-3

- **kategorie:** MEDIUM
- **quelle:** `BEO-011` Ausprägung (a) (*„Wer eine Eigenschaft am gerade
  bearbeiteten Fall beobachtet, erklärt sie für exklusiv, allgemein oder
  abschließend"*) · slice-144 §2 Schritt 0 (*„Null Falsch-Positive im Bestand
  … Ein Gate in `gates`, das legitime Sätze meldet, erzwingt Umformulierung"*)
- **pfad:** `AGENTS.md:171–176` (*„die Risiko-Klasse sind reine `a`–`f`-Wörter
  ab sieben Zeichen"*); `.d-check.yml:252–255` (dieselbe Aussage)
- **befund:** Die benannte Risiko-Klasse ist enger als das Muster: `\b[0-9a-f]{7,40}\b`
  trifft jede Ziffern-/Hex-Kette der Länge 7–40 zwischen Wortgrenzen, nicht nur
  reine `a`–`f`-Wörter — eine Sonde mit UUID, Lauf-Nummer `28679242142`,
  Kompaktdatum `20260826`, Byte-Grenze `16777216`, Hex-Farbwert und Docker-Kurz-ID
  wird in **allen** Fällen als `commit-hash` gemeldet, und der Bestand dieses
  Repos führt heute 107 rein dezimale Token dieser Länge in Markdown (Produkt:
  104). Nur die `a`–`f`-Hälfte wurde gemessen; die Zusage *„Ihre
  Falsch-Positiv-Last ist gemessen, nicht geschätzt"* steht über der ganzen
  Klasse.
- **verifizierbar:** ja — Sonden-Repo mit der realen Token-Form und den acht
  konstruierten Formen (9 von 9 Befunde); Zählung `\b[0-9]{7,40}\b` über
  Produkt und Rohtext wie oben.
- **klasse:** `risiko-klasse-schmaler-benannt-als-das-muster`

### F-4

- **kategorie:** MEDIUM
- **quelle:** slice-144 §2 Schritt 0 (*„Die Falsch-Negativ-Klasse ist zu
  **benennen**, nicht zu minimieren: welche echten Hashes das Muster nicht
  fängt, gehört ins Ergebnis"*) · `internal/hexagon/core/rules/matrix.go`
  `linkSpanRe` = `\[[^\]]*\]\([^)]*\)`, angewandt als
  `linkSpanRe.ReplaceAllString(pl.raw, " ")` · `BEO-009` Richtung (b)
- **pfad:** `.d-check.yml:257–261` (*„BENANNTE GRENZEN … Eine URL mit Hash wird
  erfasst."*); `AGENTS.md:176–178` (*„Sie hat drei **benannte Grenzen**"*);
  Commit-Botschaft `8e0a6f0`, Absatz *„DREI GRENZEN, BENANNT STATT
  WEGKONFIGURIERT"*
- **befund:** Die Liste nennt die größte Falsch-Negativ-Klasse nicht und sagt für
  sie das Gegenteil: `tokenFindings` entfernt die **ganze** Markdown-Link-Spanne
  vor der Token-Suche, also Linktext *und* Ziel — ein Hash in
  `[8e0a6f0](https://…/commit/8e0a6f0…)` bleibt ebenso unentdeckt wie einer in
  `[Beleg](https://…/commit/8e0a6f0)`, während nur die **bare** URL ohne
  Markdown-Klammern erfasst wird. In einem Repo, dessen Prosa Verweise
  durchgängig als Markdown-Links schreibt, ist das die dominante Schreibweise
  eines Commit-Verweises.
- **verifizierbar:** ja — die Falsch-Negativ-Sonde meldet die bare URL und den
  Autolink, nicht aber die beiden Markdown-Link-Formen (Zeilen 13 und 16 der
  Sonde ohne Befund, Zeile 19/22 mit Befund).
- **klasse:** `falsch-negativ-liste-nennt-den-groessten-fall-als-erfasst`

### F-5

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §3.4 (*„ein solches Token im Spec-Körper ist ein
  `matrix-forbidden`-Befund"*, ohne Vorbehalt) · `spec/lastenheft.md`
  §`DC-FA-MTX-003` §Provenance-Marker (der Marker deklariert *„Provenance/
  Verifikations-Zeiger, keine Entscheidungsgrundlage"*) ·
  `harness/conventions/MR-034-matrix-scope-welle.md` §Adaption (*„ein flaches
  Verbot — **ohne** den Provenance-Marker-Ausweg, den `ADR→Slice`
  offenlässt"*) · `internal/hexagon/core/rules/matrix.go` `tokenFindings`
  (`strings.Contains(pl.raw, provenanceMarker)` ⇒ `continue`, klassenübergreifend)
- **pfad:** `.d-check.yml:257–261` und `AGENTS.md:176–178` (die Aufzählung der
  drei Grenzen); Befund-Meldung in
  `internal/hexagon/core/rules/matrix.go:110–115`
- **befund:** Eine vierte, unbenannte Grenze ist der Zeilen-Marker: eine Zeile
  mit `<!-- d-check:status-provenance -->` nimmt **alle** Token-Klassen aus,
  auch `commit-hash` — die Sonde meldet den Hash auf der markerlosen Zeile und
  schweigt auf der markierten. Die Meldung des Gates empfiehlt diesen Ausweg
  sogar wörtlich (*„Token-Referenz spec-straten → commit-hash (8e0a6f0) ist
  nicht erlaubt — Provenance via `<!-- d-check:status-provenance -->`
  deklarieren"*), während §3.4 für die Spec-Straten keine Provenance-Ausnahme
  kennt und `MR-034` für die Nachbarklasse ausdrücklich festhält, dass der
  Marker-Ausweg dort nicht offensteht.
- **verifizierbar:** ja — Sonden-Repo mit zwei identischen Zeilen, eine davon
  mit Marker: ein Befund statt zwei, Exit 1 statt 1 mit zwei Zeilen; die
  Meldungs-Zeichenkette per `--json` gelesen.
- **klasse:** `gate-meldung-empfiehlt-einen-ausweg-den-die-regel-nicht-kennt`

### F-6

- **kategorie:** MEDIUM
- **quelle:** `spec/spezifikation.md:2577` (`matrix.classes[].paths` · Typ
  `string[]` · Default `—` · Constraint *„Glob (Mitgliedschaft)"* — dieselbe
  Default-Zelle wie das verpflichtende `matrix.classes[].name`, im Unterschied
  zu `.order`/`.token` mit „leer" und zu `ignore-refs[].in` mit „— (offen)") ·
  `spec/lastenheft.md` §`DC-FA-MTX-001` (*„Die Konfiguration deklariert
  Dokumentklassen über Pfad-Muster"*) ·
  `internal/hexagon/core/model/config.go:144` (*„MatrixClass ist eine über
  Pfad-Globs deklarierte Dokumentklasse"*) · Vorbefund F-3 in
  `docs/reviews/2026-08-25-slice-138-matrix-wellen-klasse-review.md` ·
  slice-144 §Berührte Spec-Stellen (*„—"*)
- **pfad:** `.d-check.yml:241–244` (*„matrix verlangt nur `name`, und genau das
  macht die Kategorie ausdrueckbar"*); `AGENTS.md:171–172` (*„eine Token-Klasse
  **ohne Zieldateien**"*)
- **befund:** Eine Hard Rule stützt sich jetzt auf eine Konfigurations-Freiheit,
  die nur in `applyMatrix` existiert (validiert wird ausschließlich
  `c.Name != ""` und die Eindeutigkeit) und der die Spezifikation
  widerspricht — deren Schema-Tabelle führt `paths` ohne Default, und
  `DC-FA-MTX-001` beschreibt eine Klasse als über Pfad-Muster deklariert. Der
  Slice weist ausdrücklich keine berührte Spec-Stelle aus, obwohl die neue
  Klasse die erste ohne `paths` ist.
- **verifizierbar:** ja — ein fail-closed-Constraint auf nicht-leere `paths`
  (den die Default-Zelle `—` beschreibt) würde `make doc-check` repo-weit auf
  Exit 2 setzen; heute belegt durch Lesen von
  `internal/adapter/driven/configyaml/configyaml.go:1944–1958` gegen
  `spec/spezifikation.md:2576–2580`.
- **klasse:** `hard-rule-stuetzt-sich-auf-undokumentierte-config-freiheit`

### F-7

- **kategorie:** LOW
- **quelle:** `BEO-017` (*„wer einen Wächter probt, liest die Meldung, nicht den
  Exit"* — hier die Gegenrichtung: das Notat muss die Zahl tragen, die es
  behauptet) · slice-144 §2 Schritt 1 (*„je Kandidat am Bestand messen"*)
- **pfad:** Commit-Botschaft `8e0a6f0`, Kandidaten-Tabelle, vierte Zeile
  (`(?i)commit\s+[0-9a-f]` ⇒ *„0 Befunde"*)
- **befund:** Der vierte Kandidat reproduziert die notierte Null nicht: genau so
  gefahren meldet er **2** Befunde in den Straten
  (`spec/lastenheft.md:2032` „Commit e", `:2063` „Commit b"); erst die
  vervollständigte Form `(?i)commit\s+[0-9a-f]{7,40}` liefert 0. Die
  Commit-Botschaft ist der einzige Ort, an dem der Kandidatensatz überhaupt
  festgehalten ist — die Konfigurations-Zeile nennt nur *„vier Kandidaten, alle
  0"*.
- **verifizierbar:** ja — Sonden-Config mit Quellklasse = drei Straten und
  Token `(?i)commit\s+[0-9a-f]`; die beiden Fundstellen sind
  `--commit-msg`-Prosa bzw. „neuer Commit bzw. ein menschliches".
- **klasse:** `mess-notat-reproduziert-die-notierte-zahl-nicht`

### F-8

- **kategorie:** LOW
- **quelle:** `BEO-011` Ausprägung (b) (*„ein Kriterium, das eine bereits
  korrigierte Menge nachträglich begründet"*) · `spec/spezifikation.md:1647`
  (`pins`: `target` = `sha256:<gekürzter-Pin>`)
- **pfad:** `.d-check.yml:263–264`; Commit-Botschaft `8e0a6f0` (*„Die Obergrenze
  40 ist tragend: ein 64-stelliger sha256-Digest hat keine Wortgrenze in der
  Mitte"*)
- **befund:** Die Begründung trägt nur für die Vollform. Die Schreibweise, die
  dieses Repo tatsächlich verwendet, ist der **gekürzte** Digest
  (`sha256:0cbe2d54…`, 50 Vorkommen in Markdown) — dort steht hinter dem
  Hex-Rest eine Wortgrenze, und die Sonde meldet ihn als `commit-hash`, während
  der volle 64-Steller und die `@sha256:`-Pin-Form korrekt stumm bleiben. Ein
  Image-Digest ist keine der fünf Kategorien des §3.4.
- **verifizierbar:** ja — Falsch-Negativ-Sonde, Zeilen mit vollem Digest (kein
  Befund) gegen die Zeilen mit `sha256:4980715a…` und `sha256:4980715a…b5c6`
  (je ein Befund).
- **klasse:** `begruendung-gilt-nur-fuer-die-vollform-der-schreibweise`

### F-9

- **kategorie:** LOW
- **quelle:** `.d-check.yml:266–270` (`{from: spec-straten, to: commit-hash,
  allow: false}` — die einzige Kante der Klasse) ·
  `internal/hexagon/core/rules/matrix.go` `tokenFindings` (`ruleFor(cfg.Rules,
  srcClass, c.Name)`; ohne Regel wird die Klasse übersprungen) · `BEO-009`
  Richtung (b)
- **pfad:** Commit-Botschaft `8e0a6f0` (*„Daneben stehen 1181 hash-artige Token
  ausserhalb der Straten — die Menge, gegen die die Regel überhaupt schützt."*)
- **befund:** Die Regel wirkt ausschließlich **innerhalb** der drei Straten;
  Token außerhalb werden von keiner Kante berührt und sind damit gerade nicht
  die Menge, gegen die sie schützt. Der Zusatz *„ausserhalb der Straten"* grenzt
  zudem nichts ein — die Straten tragen null solcher Token, die Zahl ist die
  Gesamtzahl des Repos.
- **verifizierbar:** ja — Rohtext-Zählung getrennt nach Straten (0) und Rest
  (1191); Produktlauf mit ADR-Quellklasse und derselben Token-Form ⇒ 0 Befunde,
  weil keine `{from: adr, to: commit-hash}`-Regel existiert.
- **klasse:** `zahl-als-schutzumfang-gerahmt-den-die-regel-nicht-hat`

### F-10

- **kategorie:** LOW
- **quelle:** `.d-check.yml:210–213` (Klasse `adr`: `paths`, **kein** `token`) ·
  `internal/hexagon/core/rules/matrix.go` `tokenFindings` (*„if c.Token == nil …
  continue"*) · Vorbefunde F-1 in
  `docs/reviews/2026-08-25-slice-138-matrix-wellen-klasse-review.md` und die
  dort genannte Kette slice-132/slice-136 (dreimal dasselbe Muster)
- **pfad:** `AGENTS.md:166–169`
- **befund:** Der neu gefasste Satz *„`make doc-check` (Modul `matrix`) hält
  **ADRs**, **Slices**, **Wellen** und **Commit-Hashes** — ein solches Token im
  Spec-Körper ist ein `matrix-forbidden`-Befund"* gilt in seiner Token-Hälfte
  für drei der vier Aufzählungsglieder: die Klasse `adr` trägt kein `token`,
  ein bares `ADR-0022` im Spec-Körper erzeugt keinen `matrix`-Befund (die
  Linkpflicht dafür trägt `ids`, ein anderes Modul). Der Wortlaut wurde in
  diesem Commit umgeschrieben und die Aussage damit erneuert.
- **verifizierbar:** ja — Sonden-Repo mit der vollständigen realen Klassen- und
  Regel-Liste und einer Stratum-Zeile, die `ADR-0022`, `welle-84`, `slice-144`
  und einen Hash nebeneinander nennt: drei Befunde, keiner für die ADR-Kennung.
- **klasse:** `token-zusage-fuer-eine-klasse-ohne-token`

### F-11

- **kategorie:** INFO
- **quelle:** `89363e3` (Botschaft: *„Vorab gemessen, weil AGENTS.md §3.4 es
  behauptet: eine Token-Klasse OHNE Zieldateien ist ausdrückbar"*) gegen den im
  selben Commit eingefügten Text
- **pfad:** `docs/plan/planning/in-progress/slice-144-commit-hash-muster.md:31–32`
- **befund:** Der Schwellen-Text sagt *„festgehalten am 2026-08-26, vor dem
  ersten Lauf"*, während dieselbe Commit-Botschaft einen bereits gefahrenen Lauf
  berichtet. Beides ist nur vereinbar, wenn „Lauf" die Falsch-Positiv-Messung
  meint und nicht den Ausdrückbarkeits-Nachweis — der Satz sagt das nicht, und
  die Reihenfolge-Zusage dieses Commits ist gerade sein Kernanspruch.
- **verifizierbar:** nein — Lesart-Frage, kein Gate. Die Reihenfolge selbst ist
  belegt (siehe Negativbefund 1).
- **klasse:** `reihenfolge-zusage-mit-eigener-ausnahme-im-selben-commit`

## Negativbefunde

1. **Geprüft, ohne Befund — Kernfrage A, die Schwelle stand vorher.**
   `git log --format='%h %ad %s' --date=iso` zeigt `7dc8738` 16:39:24 →
   `89363e3` 16:41:47 → `8e0a6f0` 16:48:21; `89363e3` fasst **nur** die
   Slice-Datei an (+15/−0, `git show --numstat`) und enthält die drei
   Bedingungen im Wortlaut. `git diff 89363e3 HEAD` auf die Slice-Datei ist
   leer — der Schwellentext ist seit seiner Niederschrift unverändert. Die drei
   Bedingungen decken sich mit denen, gegen die der Feat-Commit entscheidet
   (null Falsch-Positive · Positivkontrolle mit gelesener Ursache und grünem
   Rückbau · Falsch-Negativ-Klasse benannt); der vierte Satz (*„Ein Kandidat,
   der die erste Bedingung nur mit Ausnahmen erfüllt, ist durchgefallen"*) ist
   eingehalten — der Diff fügt weder `exempt-paths` noch `ignore-refs` noch eine
   `exclude-sections`-Zeile für die neue Klasse hinzu. Keine nachträgliche
   Anpassung im Sinne von `BEO-011`.
2. **Geprüft, ohne Befund — die Null im Bestand der drei Straten.**
   Alle vier Kandidaten in der Form, in der sie 0 tragen sollen, über das
   Produkt gefahren: 0 Befunde; Rohtext-Gegenprobe über die drei Straten
   ebenfalls 0 für `\b[0-9a-f]{7,40}\b`. Die Straten sprechen über Hashes
   ausschließlich in Platzhalter-Form (`<hex>`, `<64-hex>`,
   `spec/spezifikation.md:2549`, `:2590`) — kein literaler Hex-String.
3. **Geprüft, ohne Befund — die Positivkontrolle trennt die Kandidaten wie
   behauptet.** Kurzform, Langform und Inline-Code gegen alle vier Muster:
   `\b[0-9a-f]{7,40}\b` meldet alle drei, `\b[0-9a-f]{40}\b` nur die Langform,
   `\b[0-9a-f]{7,12}\b` nur die Kurzformen, `(?i)commit\s+[0-9a-f]{7,40}`
   Kurz- und Langform, aber nicht die Inline-Code-Form. Die Aussage *„Ohne diese
   Probe wäre 0 Befunde von die Klasse tut nichts nicht zu unterscheiden
   gewesen"* trägt und ist die von `BEO-017` verlangte Gegenrichtung.
4. **Geprüft, ohne Befund — Kernfrage D, Reichweite auf ADRs.** Die Klasse
   **muss** nicht für ADRs gelten: die kanonische SDP-Matrix
   (`grundlagen-referenz-richtung.md`) kennt keine Kategorie „Commit-Hash", der
   §3.4-Satz, aus dem die fünf Kategorien stammen, ist ausdrücklich auf
   Spec-Straten skopiert, und ein Commit-SHA in einem ADR ist regelmäßig ein
   **Pin**, kein Abwärtsverweis — der einzige heutige Fund in ADR-Körpern ist
   der Regelset-Pin in `docs/plan/adr/0010-semgrep-hermetisches-gate.md:34`, und
   §3.9 verlangt für Workflows sogar ausdrücklich volle 40-Steller. Die
   Auslassung ist damit anders gelagert als die `adr→welle`-Kante, die
   `MR-034` nachziehen musste (dort führte die Baseline eine ❌-Kante, hier gibt
   es keine). Kein Finding; die Asymmetrie ist aber der Grund für F-10s
   Wortlaut-Befund.
5. **Geprüft, ohne Befund — Kernfrage D, Namenskollision `commit-hash` ↔ Modul
   `commits`.** Klassennamen und Modulnamen liegen in getrennten Namensräumen:
   `--enable`/`--disable` und `modules:` prüfen gegen die Modul-Liste,
   `matrix.rules` gegen `matrix.classes[].name`; kein Codepfad vergleicht die
   beiden. Der Befund-Grund bleibt `matrix-forbidden` (nicht
   `commit-untraceable`), und `--doctor` gruppiert nach Regel, nicht nach
   Klasse. Verwechslungsgefahr besteht nur in der Prosa.
6. **Geprüft, ohne Befund — Kernfrage E, die Obergrenze 40 für die Vollform.**
   Ein echter 64-Steller aus diesem Repo wird nicht gemeldet, auch nicht mit
   `sha256:`-Präfix und nicht in der `@sha256:`-Pin-Form — die Begründung „keine
   Wortgrenze in der Mitte" trägt für diese Schreibweisen. Der Vorbehalt steht in
   F-8.
7. **Geprüft, ohne Befund — Kernfrage F, der verbleibende Satz zum
   Closure-Datum.** *„Die Spec-Straten führen ihre eigenen Historie-Zeilen
   voller Daten"* stimmt: `spec/lastenheft.md:2970` und
   `spec/spezifikation.md:2790` tragen je einen `## 7. Historie`-Abschnitt mit
   Datumsspalte, und `exclude-sections: [Geschichte]` nimmt ihn **nicht** aus
   (exakter Heading-Vergleich). Der frühere Auflösungs-Trigger für Commit-Hashes
   ist eingelöst und ersatzlos entfernt; für das Closure-Datum steht *„keiner"*
   samt Begründung, wie die Botschaft es sagt.
8. **Geprüft, ohne Befund — Modul-Grenze auf der Ziel-Achse (Skill §3.8).** Die
   neue Klasse liest keine Eingabe, die sie nicht scannt: ohne `paths` ist sie
   nie Quell-Klasse (`classOf` matcht sie nie), nie Link-Ziel und tritt
   ausschließlich in `tokenFindings` über den bereits gescannten Prosa-Körper
   auf. Kein Status-Read, kein `exempt-paths`-Nebeneffekt.
9. **Geprüft, ohne Befund — `make doc-check` ist grün und die Zahl stimmt.**
   Exit 0, 515 Dateien, 0 Befunde, mit einem aus HEAD frisch gebauten Image.
   Die Botschaft *„Kein natürliches Wort dieser Form kommt in über 500
   Dokumenten vor"* trägt: 515 gescannte Markdown-Dateien, und die einzige
   Fundform ist der Platzhalter `deadbeef`.
10. **Geprüft, mit bekannter, geteilter Grenze — kein neuer Befund dieses
    Commits.** Mischschrift (`8e0A6f0`) bleibt wie Großschrift stumm (dieselbe
    Ursache, nur die Vollform ist benannt); ein über einen Zeilenumbruch
    getrennter Hash wird als **zwei** Token gemeldet, nicht als einer; zwei
    identische Hashes auf derselben Zeile ergeben wegen `SortFindings` **einen**
    Befund. Eingerückter Code, HTML-Kommentare, Tabellenzellen,
    Referenz-Link-Definitionen und Satzzeichen-/Klammer-Nachbarschaft werden
    korrekt erfasst. Alles davon ist geteilte, schon für `slice`/`welle`
    wirksame Mechanik.
11. **Geprüft, ohne Befund — §3-Verbot „keine Ausnahme-Liste als Ersatz für
    Präzision".** Der Diff fügt der neuen Klasse keine einzige Ausnahme hinzu;
    `exempt-paths` und `exclude-sections` sind unverändert, und die drei
    grandfatherten ADR-Globs stammen aus früheren Slices.
12. **Nicht gefahren:** `make gates` (zehn Glieder) und `make fullbuild` in
    voller Länge — die Botschafts-Zeilen *„make gates Exit 0"* und
    *„make fullbuild Exit 0"* sind weder bestätigt noch widerlegt; bestätigt ist
    nur `make doc-check` Exit 0. Die Zeile *„make doc-check Exit 2"* für den
    Rückbau-Beleg ist der Exit von `make` (GNU make meldet jeden
    Recipe-Fehler als 2), nicht der des Werkzeugs — `d-check` liefert bei
    Befunden Exit 1 (`DC-FA-CLI-003`); dieselbe Schreibweise steht in älteren
    Botschaften dieses Repos und wird darum nicht als Finding geführt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 5 |
| LOW | 4 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:**
`konfig-kommentar-traegt-mess-historie-und-slice-nummer` ·
`zahl-ohne-messvorschrift` · `risiko-klasse-schmaler-benannt-als-das-muster` ·
`falsch-negativ-liste-nennt-den-groessten-fall-als-erfasst` ·
`gate-meldung-empfiehlt-einen-ausweg-den-die-regel-nicht-kennt` ·
`hard-rule-stuetzt-sich-auf-undokumentierte-config-freiheit` ·
`mess-notat-reproduziert-die-notierte-zahl-nicht` ·
`begruendung-gilt-nur-fuer-die-vollform-der-schreibweise` ·
`zahl-als-schutzumfang-gerahmt-den-die-regel-nicht-hat` ·
`token-zusage-fuer-eine-klasse-ohne-token` ·
`reihenfolge-zusage-mit-eigener-ausnahme-im-selben-commit`

## Verdikt

**Merge-blockierend:** ja — ein HIGH und fünf MEDIUM.

Die **Entscheidung** dieses Slice hält: die Schwelle stand nachweislich vorher
und wurde nicht angepasst, die Null im Bestand ist reproduzierbar, die
Positivkontrolle trennt die Kandidaten, und die Ausklammerung der ADRs ist
sachlich richtig. Was nicht hält, ist die **Buchführung über die Messung**. Vier
der elf Findings betreffen Zahlen und Zusagen, die weiter reichen als das, was
gemessen wurde: die beiden Risiko-Zahlen sind über den Rohtext entstanden und
stehen unter einer Kopfzeile, die „mit dem Produkt, nicht per grep" verspricht
(F-2); die benannte Risiko-Klasse deckt nur die Buchstaben-Hälfte des Musters
(F-3); die Liste der Grenzen nennt die größte Falsch-Negativ-Klasse nicht und
behauptet für sie das Gegenteil (F-4); und der vierte Kandidat reproduziert
seine notierte Null nicht (F-7). Das ist `BEO-009` Richtung (b) und `BEO-011`
Ausprägung (a) in einem Slice, der beide Beobachtungen ausdrücklich als seinen
Anlass führt.

F-1 ist davon unabhängig und die härteste Einzelstelle: der neue
Konfigurations-Kommentar hält die Entstehungsgeschichte fest, die §3.7 dort
verbietet — dieselbe Bewegung, die der Nutzer-Vorgabe zu CR-Dokumenten
zugrunde liegt (Beleg gehört in den Kommentar, Entstehung in die
Commit-Botschaft), nur in die andere Richtung. F-5 ist der Kandidat für eine
Eskalation, wenn er wiederkehrt: die Meldung des Gates empfiehlt einen Ausweg,
den die Hard Rule nicht kennt, und das ist ein Gate-Pfad.

**Übergabe:** Findings gehen an den Implementer; F-2/F-3/F-4/F-7 sind
gemeinsam Kandidaten für den `BEO-009`-Zähler (Richtung b) und F-3 zusätzlich
für `BEO-011` (a) — die Einordnung trifft der Maintainer, nicht dieser Report.
Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation; DoD- und
Plan-Konformität prüft der Verifier separat.
