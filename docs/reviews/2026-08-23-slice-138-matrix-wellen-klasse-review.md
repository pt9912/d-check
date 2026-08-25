# Review-Report: slice-138 — matrix bekommt die welle-Klasse — 2026-08-23

**Review-Art:** Code-Review (Config-/Doku-Diff gegen Kanon und Slice-Plan,
Modul 10 §Drei Review-Arten) · **Gegenstand:** Commit `a375222` (Diff
`HEAD~1..HEAD`) — der Feature-Commit von slice-138 (`feat(config): matrix
bekommt die welle-Klasse — §3.4 von zwei auf drei von fünf (DC-FA-MTX-001,
DC-FA-MTX-003, slice-138)`). 3 Dateien laut `git show --stat`
(`.d-check.yml` +17, `AGENTS.md` +18/-6 [netto +12], `spec/lastenheft.md`
+5/-2), 32 Einfügungen/8 Löschungen gesamt.

**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 (`9a7654a`) · **Modell-ID:**
`claude-sonnet-5` · **Datum:** 2026-08-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-138-matrix-wellen-klasse.md`
  (§1–§9, insbesondere §3 NICHT, §4 DoD, §5 Risiken)
- `.d-check.yml` Block `matrix` vollständig (Klassen, Regeln,
  `exclude-sections`, `exempt-paths`)
- `AGENTS.md` §3.4 (Diff) und §1 (Source Precedence/Konflikt-Meldepflicht)
- `spec/lastenheft.md` §7 (die geänderte 0.60.0-Zeile **und** die neue
  0.65.4-Zeile)
- `docs/plan/adr/0047-matrix-spec-historie-nicht-provenance-exempt.md`
  (Volltext, inkl. `## Kontext` und `## Geschichte`)
- `spec/lastenheft.md` → `DC-FA-MTX-001`, `DC-FA-MTX-002`, `DC-FA-MTX-003`
  (Volltext)
- `harness/conventions.md` → `MR-006` (`harness/conventions/MR-006-referenzrichtung-matrix.md`),
  `MR-032` (`harness/conventions/MR-032-historie-vor-accepted.md`)
- Baseline-Kanon `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)
  (Volltext, inkl. Matrix-Tabelle, Planungs-Ebene-Diagramm, Regel 5)
- Quellcode: `internal/hexagon/core/rules/matrix.go` (`CheckMatrix`,
  `tokenFindings`, `proseLines`), `internal/adapter/driven/configyaml/configyaml.go`
  (`applyMatrix`, `compileMatrixToken`) — zur Prüfung, was `matrix` mechanisch
  tatsächlich leistet, unabhängig von der Commit-Prosa
- `docs/plan/planning/observations.md` (`BEO-011`, `BEO-012`)
- Vorherige Review-Reports am selben Abschnitt (Kalibrierungs-Anker):
  `docs/reviews/2026-08-23-slice-132-hard-rule-zensus-review.md` und
  `docs/reviews/2026-08-23-slice-136-agents-34-klaerung-review.md` (beide F-1:
  ein „gedeckt"-Verdikt deckt nur einen Teil der genannten Kategorien)

**Nicht erhalten:** die DoD-Abhakung (Verifikations-Rolle, getrennter Kontext).

**Vom Reviewer selbst gefahren** (nur Lesekommandos, kein `make`, keine
Dateiänderung): `git show`/`git log -p`/`git diff` auf mehrere Commits
(`a375222`, `100ec3f`, `HEAD~1`); `find`/`grep` über `docs/plan/planning/**`
und alle drei Spec-Straten (roh, case-insensitiv, raumtolerant, gegen
`welle-[0-9]+` und Varianten); `awk`-Sektionierung aller nicht-grandfatherten
ADRs mit `welle-NN`-Vorkommen gegen ihre `##`-Überschrift; Lesen von
`matrix.go`/`configyaml.go`, um `token`/`paths`-Unabhängigkeit zu verifizieren;
manuelle `grep -oP`-Proben des Token-Regex gegen Suffix-, Groß-/Klein- und
Inline-Code-Fälle.

**Verdikt: blockierend** — ein HIGH, zwei MEDIUM.

---

## Findings

### F-1 — Die fallengelassene `adr → welle`-Regel widerspricht der zitierten SDP-Matrix und lässt MR-006 unaufgelöst stehen

- **kategorie:** HIGH
- **quelle:** `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP),
  Matrix-Zeile **ADR** → Spalte **Welle** = `❌` (strenger als ADR→Slice, das
  dort als `Kontext … der zulässige Zeiger wird in seiner Zeile markiert`
  geführt wird) sowie der Satz „Dazu die einzige Kontext-Kante, die den
  kanonischen Block verlässt: **ADR → Slice**. Sie ist deshalb die einzige,
  die eine Markierung in ihrer Zeile verlangt" · `MR-006`
  (`harness/conventions/MR-006-referenzrichtung-matrix.md`), Feld
  „Adaption": „die weiteren ❌-Kanten der 8×8-Matrix — **ADR→Carveout/Welle/Roadmap**
  … sind **bewusst unbewacht** (d-check modelliert Carveout/Welle/Roadmap
  **nicht** als `matrix`-Klassen; eine Erweiterung wäre ein eigener Change)" ·
  `AGENTS.md` §1 („Melde den Widerspruch, statt ihn stillschweigend nach einer
  Seite aufzulösen") · slice-138 §5 Risiko 1 · Kalibrierung: F-1 in
  `docs/reviews/2026-08-23-slice-132-hard-rule-zensus-review.md` und F-1 in
  `docs/reviews/2026-08-23-slice-136-agents-34-klaerung-review.md` — exakt
  dasselbe Fehlermuster („gedeckt" trägt weniger, als der Satz suggeriert)
  ein drittes Mal, jetzt eine Ebene tiefer innerhalb der `welle`-Kategorie
  selbst.
- **pfad:** `.d-check.yml:208–223` (Klasse `welle`, fehlende
  `{from: adr, to: welle}`-Regel); `docs/plan/adr/0047-matrix-spec-historie-nicht-provenance-exempt.md:21`
  (lebender, ungeprüfter Fund: „Der adoptierte Baseline-Pin ist auf **v5.0.0**
  gehoben (welle-67)." — `## Kontext`, **nicht** `## Geschichte`, also nicht
  über `exclude-sections` gedeckt)
- **befund:** Der Commit fügt `{from: spec-straten, to: welle, allow: false}`
  hinzu, lässt `{from: adr, to: welle}` aber bewusst weg — mit der
  Begründung, `§3.4 verlangt sie ohnehin nicht: sie gilt den SPEC-STRATEN`,
  und der Alternative „Falsch-Verdikt vs. pauschale Ausnahme". Diese
  Begründung trifft `§3.4` korrekt (die Hard Rule spricht ausschließlich von
  Spec-Straten), zieht daraus aber einen Schluss, den dieselbe Logik für die
  bereits **bestehende** `{from: adr, to: slice}`-Regel ebenso träfe — auch
  diese ist aus `§3.4`s Wortlaut nicht ableitbar, sondern aus `MR-006`/der
  SDP-Matrix. Genau diese Quelle wird für die `welle`-Entscheidung nicht
  zitiert und nicht aufgelöst, obwohl sie (a) `ADR→Welle` als striktes `❌`
  führt — strenger als das markerfähige `ADR→Slice` — und (b) in `MR-006`
  eine **explizite, Accepted** Scope-Grenze trägt, deren Vorbedingung
  („d-check modelliert Carveout/Welle/Roadmap nicht als `matrix`-Klassen")
  dieser Commit selbst für `Welle` aufhebt, ohne die Grenze nachzuziehen oder
  zu kommentieren. Die genannte Zwei-Wege-Alternative ist zudem eine falsche
  Dichotomie: `DC-FA-MTX-003`s Grandfathering (`matrix.exempt-paths`) ist
  bereits als **nicht-pauschale**, dateispezifische Ausnahme etabliert (aktuell
  `docs/plan/adr/00[01][0-9]-*.md`/`002[01]-*.md` für exakt denselben Fall bei
  `slice`) — ein analoger Eintrag für `docs/plan/adr/0047-*.md` hätte die Regel
  für **alle künftigen** ADRs scharf gehalten und nur das eine unmarkierbare
  Bestandsdokument ausgenommen, ohne „pauschal" zu sein. Der tatsächliche
  Effekt: jedes künftige ADR kann jede `welle-NN`-Kennung unbegrenzt und
  ungeprüft im Körper nennen; `AGENTS.md:166–168`s „gedeckt sind drei" trägt
  für „Wellen" damit eine andere (schwächere, nur einseitige) Bedeutung als
  für „ADRs" und „Slices" (dort beidseitig: spec-straten- **und**
  adr-seitig), ohne dass der Unterschied benannt wird.
- **verifizierbar:** ja — eine probeweise ergänzte Regel `{from: adr, to:
  welle, allow: false}` plus `exempt-paths`-Eintrag für `docs/plan/adr/0047-*.md`
  ließe `make doc-check` exakt einen historischen Fund (den bereits
  grandfatherten) unterdrücken und jeden **neuen** unmarkierten `welle-NN`-Token
  in einem beliebigen anderen ADR-Körper als `matrix-forbidden` melden; nicht
  selbst ausgeführt (Auftrag untersagt `make`-Läufe).
- **klasse:** sdp-kante-durch-lokale-regeltext-lesart-uebergangen

### F-2 — Die de-tokenisierte 0.60.0-Zeile ersetzt die Kennung durch einen unauflösbaren Verweis statt einer eigenständigen Formulierung

- **kategorie:** MEDIUM
- **quelle:** `grundlagen-referenz-richtung.md` Regel 5 (Grund:
  „Unreparierbarkeit … Eine Historie-Zeile ist ein Protokoll; sie wird nicht
  rückwirkend geändert") · slice-138 §5 Risiko 3 („Die Unterscheidung [Form
  vs. Aussage] ist zu belegen, nicht zu behaupten") · Präzedenz-Commit
  `100ec3f` (ADR-0047 Entscheidung 2, slice-087): Zeile 0.47.0 ersetzte
  „gemeinsames Kriterium **mit `slice-079`**" durch das **eigenständige**
  „gemeinsames Kriterium" (ersatzlose Streichung) und „folgen **in
  `slice-078`**" durch das **generische** „folgen **separat**" (kein
  Rückverweis nötig) — dieselbe Zeile verwendet an anderer Stelle „jene"
  korrekt **innerhalb desselben Satzes** mit gegenwärtigem Antezedens
  („`scan.ignore` … — jene entfernt Dateien").
- **pfad:** `spec/lastenheft.md:2988` (0.60.0-Zeile)
- **befund:** Die Zeile lautete vor dem Commit „… reproduziert den
  historischen 19-Link-Bruch **der welle-69-Eröffnung** zeichengenau" und
  lautet danach „… reproduziert den historischen 19-Link-Bruch **jener
  Wellen-Eröffnung** zeichengenau". „welle-69" war die **einzige** Stelle in
  der gesamten Zeile, die identifizierte, welche Welle gemeint ist; nach der
  Streichung verweist „jener" auf nichts mehr innerhalb der Zeile — ein Leser
  ohne Git-Historie kann nicht mehr feststellen, welche Welleneröffnung
  gemeint ist. Das unterscheidet sich vom Präzedenzfall, der entweder
  ersatzlos strich oder einen generischen, referenzlosen Ersatz wählte. Die
  neue 0.65.4-Zeile behauptet „die Aussage bleibt unverändert, der Verweis
  ist aufgelöst" — für diese spezifische Zeile ist die Aussage aber nicht
  unverändert, sondern um ihre einzige Identifizierungs-Information ärmer;
  das ist eher eine Grenz- als eine Formänderung, und Risiko 3 verlangt genau
  hierfür einen Beleg, den der Commit nicht liefert (er erwähnt diese
  Zeilenänderung nicht einmal explizit).
- **verifizierbar:** nein — Lesbarkeits-/Vollständigkeitsurteil, kein Gate
  (`matrix` prüft nur Token-Abwesenheit, nicht Referenz-Auflösbarkeit der
  Ersatzformulierung).
- **klasse:** detokenisierung-hinterlaesst-unauflösbaren-verweis

### F-3 — „Kein Modul kennt ein Muster-Verbot" für Commit-Hashes ist unbelegt und vom eigenen Code widersprochen

- **kategorie:** MEDIUM
- **quelle:** `internal/adapter/driven/configyaml/configyaml.go:1938–1957`
  (`applyMatrix`: `c.Paths` wird nirgends als nicht-leer validiert;
  `compileMatrixToken` kompiliert `token` unabhängig von `Paths`) ·
  `internal/hexagon/core/rules/matrix.go` `tokenFindings` (iteriert
  `cfg.Classes` nach `c.Token`, referenziert `c.Paths` nicht) · derselbe
  Commit, Belegprobe „welle-99 im Spec-Körper => matrix-forbidden" (ein Token
  ohne zugehörige reale Datei löst bereits heute einen Fund aus — derselbe
  Mechanismus, der für eine Commit-Hash-Klasse fehlen soll) · `BEO-011`/`BEO-012`
  (Bestandsaussage „kein Modul kennt X" ungeprüft gegen den eigenen
  Mechanismus).
- **pfad:** `AGENTS.md:170–171` („Ein **Commit-Hash** ist kein Dokument;
  `matrix` verbietet Verweise auf eine Klasse von *Dateien*, und dafür gibt es
  hier keine — es bräuchte ein Muster-Verbot, das kein Modul kennt.")
- **befund:** `matrix.classes[].token` ist bereits genau ein solches
  Muster-Verbot, und `paths` ist für eine Klasse **nicht** erforderlich (der
  Config-Loader validiert nur `Name`; `Paths`/`Order`/`Direction`/`Token`
  bleiben optional). Eine Klasse mit leeren/nicht-matchenden `paths` und
  `token: '[0-9a-f]{7,40}'` würde über denselben Pfad wie `welle-99`
  (Token-Fund ohne reale Zieldatei) jeden Hex-String im Spec-Körper als
  `matrix-forbidden` melden — mechanisch also **kein** fehlendes Modul,
  sondern ein bereits vorhandener, nur ungenutzter Mechanismus. Der reale
  Hinderungsgrund ist die **Präzision** eines bloßen Hex-Musters (hohe
  Falsch-Positiv-Last: jede zufällige Hex-ähnliche Zeichenkette träfe), nicht
  das Fehlen einer Modul-Fähigkeit; der Satz „das kein Modul kennt"
  formuliert eine Kategorie-Aussage, die die eigene, im selben Commit
  demonstrierte Mechanik übersieht.
- **verifizierbar:** ja, durch Code-Lesen bestätigt (nicht durch einen
  tatsächlichen Lauf, da `make`-Läufe untersagt sind) — ein Testfall mit
  leerem `paths` + `token`-Klasse wäre der Gate-Beleg.
- **klasse:** bestandsaussage-ueber-modul-faehigkeit-unbelegt

## Negativbefunde

- geprüft, ohne Befund: **Klassen-Pfade** `["docs/plan/planning/**/welle-*.md"]`
  matchen exakt alle 43 Wellendokumente (25 Wellen: 7×nur `-results.md`,
  18×Plan- **und** `-results.md`-Datei) über alle vier Lifecycle-Verzeichnisse
  hinweg via `find`; kein Wellendokument bleibt unerfasst, kein Fremd-Treffer.
  `**`-Segment-Matching in `paths.go:98–116` ist geteilte, bereits an `slice`
  erprobte Maschinerie (rekursive Segment-Suche inkl. null Zwischensegmente),
  kein neues Risiko.
- geprüft, ohne Befund (Frage aus Auftrag 1): **`welle-*-results.md` als Teil
  der Klasse** ist korrekt, nicht fehlerhaft — beide Dokumentformen
  repräsentieren dieselbe zeitliche Schicht (`§3.4`), keine bestehende oder
  angelegte Regel unterscheidet sie, und kein Link von Spec-Straten auf eine
  Wellendatei existiert im Bestand (weder Plan- noch Ergebnisnotiz-Form).
- geprüft, ohne Befund (Frage aus Auftrag 2, Suffix-Fall): `welle-\d{2,}`
  matcht korrekt als Teilstring in „welle-84-durchsetzung" (`grep -oP`-Probe);
  kein Falsch-Negativ bei Kennung-plus-Suffix.
- geprüft, mit bekannter, **geteilter** Grenze — kein neuer Befund dieses
  Commits: das Token ist case-sensitiv (verpasst „Welle-84") und
  `tokenFindings`/`proseLines` entfernen nur **Fenced**-Code, nicht
  Inline-Backticks — ein `` `welle-84` `` in Prosa würde mitgezählt. Beides ist
  dieselbe, bereits für `slice` etablierte Mechanik (`matrix.go` `proseLines`,
  geteilt), nicht neu durch diesen Commit eingeführt, und im Bestand aktuell
  folgenlos (kein solcher Fall vorhanden, siehe nächster Punkt).
- geprüft, ohne Befund (Frage aus Auftrag 2/3, Klassen-`paths` vs. Token
  verwechselt?): die beiden Mechanismen sind im Code sauber getrennt (`paths`
  bestimmt Klassen-Mitgliedschaft für Link-Ziele, `token` matcht Text
  unabhängig von realen Zieldateien) — der Commit selbst demonstriert das
  bewusst korrekt mit den nicht-existenten Wellen `welle-99`/`welle-123`.
  Keine Vermischung, aber auch keine erläuternde Zeile dazu im
  `.d-check.yml`-Kommentar (nicht separat als Finding gewertet — reine
  Doku-Ergänzung, kein Verhaltensrisiko).
- geprüft, ohne Befund (Frage aus Auftrag 3): weitere, dem Wächter verborgene
  `welle-NN`-Vorkommen in den drei Spec-Straten-Dateien — roher,
  case-insensitiver, raumtoleranter Scan (`[Ww]elle[n]?[^a-zA-Z]{0,15}[0-9]{1,3}`)
  liefert außer dem bereits behandelten Fund nur „Kurs-Welle 81"
  (`spec/lastenheft.md:2983`, 0.62.0-Zeile) und „welle 4" (Teilstring von
  „**Sch**welle 4", `spec/lastenheft.md:3000` — Fehlalarm des eigenen groben
  Musters). „Kurs-Welle 81" referenziert die Wellenzählung des adoptierten
  **Kurs**-Repos (`ai-harness-course`), nicht d-checks eigene Planungsebene —
  kein Wächter-Blindfleck, sondern ein legitim anderer Referent.
- geprüft, ohne Befund (Frage aus Auftrag 3, `exclude-sections`/`exempt-paths`
  als blinder Fleck): acht weitere `welle-NN`-Vorkommen in nicht-grandfatherten
  ADRs (0029, 0030, 0031, 0033, 0034, 0035, 0036, 0056) liegen ausnahmslos
  unter `## Geschichte` und sind seit ADR-0047 Entscheidung 3 legitim
  `exclude-sections`-exempt; nur `ADR-0047:21` selbst (`## Kontext`) liegt
  außerhalb — das ist der in F-1 behandelte, einzige lebende Fund. Kein
  weiterer verborgener Fund in den grandfatherten ADRs 0001–0021 gesucht
  (außerhalb des Auftrags-Skopus „Spec-Straten").
- geprüft, ohne Befund: **Risiko 2 aus §5** (Token-Weite `{2,}`) — die
  Rückbau-Probe mit `welle-99`/`welle-123` demonstriert korrekt, dass zwei-
  **und** dreistellige Kennungen greifen, und der Rückbau meldet 0 Befunde;
  keine Diskrepanz zur Behauptung gefunden.
- geprüft, ohne Befund: **Lastenheft-Bump/`MR-032`-Form** — Version 0.65.4
  gesetzt, Historie-Zeile vorhanden, `Verweis`-Spalte korrekt `—` (keine
  begonnene CR-Pflicht, konsistent mit `MR-032`).
- geprüft, ohne Befund: **Datums-Aussage** zu Closure-Daten (§3.4, „von einem
  legitimen Datum nicht unterscheidbar") — als Urteil (nicht als Gate-Zusage)
  formuliert und nicht mechanisch widerlegbar; kein Gegenbeleg gefunden.
- nicht nachvollzogen (Auftrag untersagt eigene `make`-Läufe): die
  Gate-Aussage „`make gates` Exit 0 (zehn Glieder, 483 Dateien, 0 Befunde)"
  in der Commit-Botschaft — weder bestätigt noch widerlegt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** sdp-kante-durch-lokale-regeltext-lesart-uebergangen
· detokenisierung-hinterlaesst-unauflösbaren-verweis ·
bestandsaussage-ueber-modul-faehigkeit-unbelegt

## Verdikt

**Merge-blockierend:** ja — ein HIGH (F-1) und zwei MEDIUM. F-1 ist die
dritte Wiederholung desselben Fehlermusters am selben Abschnitt
(„gedeckt"-Verdikt trägt weniger, als sein Wortlaut suggeriert — nach
`docs/reviews/2026-08-23-slice-132-hard-rule-zensus-review.md` F-1 und
`docs/reviews/2026-08-23-slice-136-agents-34-klaerung-review.md` F-1), diesmal
**innerhalb** der neu gedeckten `welle`-Kategorie: „gedeckt" bedeutet für
Wellen nur „von Spec-Straten aus", nicht symmetrisch wie bei ADRs/Slices —
und die einzige einschlägige, Accepted Konvention (`MR-006`), die genau diese
Asymmetrie schon vorwegnahm und an eine jetzt gefallene Vorbedingung band,
wird im Commit weder zitiert noch nachgezogen. F-2 und F-3 sind unabhängig
davon reale, aber kleinere Qualitätsmängel der begleitenden Doku-Änderungen.
Die Kernentscheidungen des Slice — Klassenzuschnitt (Auftrag 1), Token-Form
inkl. Suffix-/Weiten-Verhalten (Auftrag 2), Vollständigkeit der
Bestandsmessung (Auftrag 3) und die §3.4-Risikoproben (Auftrag 6) — halten der
Prüfung stand.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den Beobachtungs-Zähler —
F-1 ist ein Kandidat für `BEO-012` (Quelle mit einschlägigem Geltungsbereich,
hier `MR-006`, wird nicht konsultiert) bzw. für eine neue Beobachtung
„gedeckt-Verdikt hält nur einseitig", deren Einordnung dem Maintainer
obliegt, nicht diesem Report. Dieser Report ist ein Lauf-Beleg und ersetzt
keine Verifikation (DoD-/Spec-Konformität prüft der Verifier separat).
