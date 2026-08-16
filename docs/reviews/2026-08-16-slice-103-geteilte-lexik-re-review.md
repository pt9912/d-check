# Re-Review-Report: slice-103 — Geteilte Lexik, bestätigende Re-Review — 2026-08-16

**Review-Art:** bestätigende Re-Review — geprüft wird, ob die dreizehn Befunde
des Erst-Reports geheilt sind **und** ob die Heilung selbst neue Defekte
eingeführt hat. Jedes Urteil steht auf einem eigenen Lauf oder einem Code-/
Vertragszitat, nicht auf dem Commit-Text. Nicht geprüft wird die DoD-Abhakung
(getrennter Kontext, Verifikation).

**Gegenstand:** Commit-Range `d2aaf90..346a223` (sechs Commits: Wellen-Eröffnung ·
Messung · Vertrag · Implementierung · Erst-Report · Heilung); im Besonderen der
Heilungs-Commit `346a223`, der alle dreizehn Befunde zu schließen beansprucht.
Arbeitsbaum-Stand `346a223`.

**Skill:** `.harness/skills/reviewer.md` @ 1.4.0 ·
**Modell:** claude-opus-5[1m] · **Datum:** 2026-08-16

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Der [Erst-Report](2026-08-16-slice-103-geteilte-lexik-review.md) mit den
  Befunden L-1 bis L-13 (Prüfliste dieser Runde)
- [`DC-FA-CITE-001`](../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in),
  [`DC-FA-VER-001`](../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
  [`DC-FA-PIN-001`](../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in),
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
  [`DC-FA-ANCH-001`](../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
  und [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (Lastenheft-Fassung 0.58.1)
- §[`DC-FA-CITE-001.a`](../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
  Schritt 2, §[`DC-FA-VER-001.a`](../../spec/spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions)
  Schritt 1, §[`DC-FA-PIN-001.a`](../../spec/spezifikation.md#dc-fa-pin-001a--content-pin-gegen-inhaltlichen-drift-pins)
  Schritt 2, §[`DC-FA-ANCH-001.a`](../../spec/spezifikation.md#dc-fa-anch-001a--github-slug-algorithmus)
  Schritt 5, §[`DC-FA-ANCH-001.b`](../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker),
  §[`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  Schritt 2, §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritte 3/4 und C3
- [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  (Entscheidungen 1–4, Konsequenzen, Re-Evaluierungs-Trigger),
  [ADR-0019](../plan/adr/0019-versions-pin-fence-ausnahme.md),
  [ADR-0020](../plan/adr/0020-content-pin-fence-ausnahme.md),
  [ADR-0024](../plan/adr/0024-vcs-immutable-gate.md),
  [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md)
- [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel-Tabelle), [`AGENTS.md`](../../AGENTS.md) §3 Hard Rules,
  [`CLAUDE.md`](../../CLAUDE.md)
- Der [Slice-Plan](../plan/planning/in-progress/slice-103-geteilte-lexik-raender.md)
  (§4 Abnahme-Punkte, §4a Messung und Klassen-Abschluss) und das
  [Wellendokument](../plan/planning/welle-74-geteilte-lexik-raender.md) (§3
  Closure-Trigger)

**Läufe dieses Reviews.** Alle Fixtures liegen in einem Scratch-Verzeichnis
außerhalb des Repos; alle Läufe sind netzlos und read-only. Gefahren wurden:
`make build`, `make gates` (grün, 373 Dateien, 0 Befunde), `make verify-closure-notes`
(grün, 343/0), `make test` in drei Mutationsläufen sowie rund 60 Fixture-Läufe
gegen drei Images — den HEAD-Build, einen aus Dateikopien rekonstruierten
**Vor-Heilungs**-Build (Stand `6461bd6`, nur die vier berührten Regel-Dateien
zurückgeschrieben) und das veröffentlichte Bild `v0.52.0` als Bestands-Gegenprobe.
Für jede Mutation wurde die Repo-Datei **vor** dem Eingriff in das
Scratch-Verzeichnis kopiert und danach aus der Kopie zurückgeschrieben — **kein**
`git checkout`. **Der Arbeitsbaum ist wiederhergestellt:** `git status --short`
ist leer (einzige neue Datei ist dieser Report), und der Neubau nach der letzten
Mutation liefert dieselbe Image-ID wie der erste Build vor dem ersten Eingriff.

---

## Urteil je Erst-Befund

| Befund | Urteil | Beleg dieser Runde |
|---|---|---|
| L-1 (HIGH) Anker-Frage nur halb vereinheitlicht | **teilweise geheilt** | Die vier benannten Zeichenfolgen sind geheilt: vier Fixtures (Anker in Inline-Code, `data-id`, `name` an `<area>`, anker-förmige Prosa ohne Tag) liefern jetzt bei **allen drei** Modulen dieselbe Antwort — `anchors` meldet `anchor-missing`, `versions` bricht mit Exit 2 ab, `pins` schweigt. Zwei **andere** Achsen derselben Frage divergieren weiter: ein Duplikat-Slug-Suffix und ein prozent-kodiertes Fragment ⇒ N-1 |
| L-2 (HIGH) falsche Richtungs-Zusage | **teilweise geheilt** | Beide Richtungen stehen jetzt in [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) §Konsequenzen, in der (umgeschriebenen) Historien-Zeile 0.58.0, in der neuen Zeile 0.58.1 und in beiden `CHANGELOG.md`-Einträgen. Die neue Aufzählung ist aber selbst **geschlossen** formuliert („an zwei Stellen“, „zwei Fälle“) und der Heilungs-Commit fügt eine **dritte** Verlust-Stelle hinzu, die er nicht nennt ⇒ N-2 |
| L-3 (MEDIUM) zwei unassertierte Wirkstellen | **geheilt** | Beide Rückbauten einzeln über eine Dateikopie angewendet: Fence-Grenze im Blockquote-Sammler entfernt ⇒ `TestCitationsFenceImBlockquoteBeendetDenBlock` rot; im Absatz-Sammler entfernt ⇒ `TestCitationsFenceImAbsatzTrennt` rot. Je genau ein Test, kein anderer. Arbeitsbaum danach wiederhergestellt |
| L-4 (MEDIUM) geschlossene Terminator-Aufzählung | **geheilt** | §[`DC-FA-CITE-001.a`](../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations) Schritt 2 nennt den Fenced-Block jetzt ausdrücklich als dritten Blockquote-Terminator („eine Leer-, eine Nicht-Größer-Zeile **oder ein Fenced-Block** beendet ihn“); das Handbuch sagt denselben Satz. Die Akzeptanzkriterien folgen nicht ⇒ N-6 |
| L-5 (MEDIUM) Anforderungen auf der alten Fassung | **geheilt** (Beschreibungs-Hälfte) | Die **Beschreibungen** aller drei Anforderungen tragen die neue Semantik: `spec/lastenheft.md:1351` (Absatz-Grenze), `:1563` („Sie gilt dem Pin, nicht dem Anker“ — genau an der Fence-Ausnahme, wo der Widerspruch entstand), `:1629` (Anker-Erkennung für `pins`). Die **Akzeptanzkriterien**, die [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md) in derselben Zeile mitnennt, blieben unverändert ⇒ N-6. Der neue Text der Beschreibung ist zudem gegen einen Lauf falsch ⇒ N-4 |
| L-6 (MEDIUM) vcs-Grenze ohne stille Richtung | **geheilt** | `spec/lastenheft.md` §Grenze und [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 3 beschreiben jetzt „Öffnung im ausgenommenen Abschnitt ⇒ Ausnahme bis Dateiende ⇒ Exit 0 ohne Ausgabe“. Der Trigger hängt jetzt an etwas Beobachtbarem, und die Beobachtung ist im eigenen Gate erreichbar: `spans` steht in der `modules`-Liste von `.d-check.yml`, und die stille Richtung setzt die offene Fence in der **HEAD**-Fassung voraus, also in der Fassung, die `spans` scannt |
| L-7 (LOW) Fence-Aufzählung nennt zwei von vier | **geheilt** (Fence-Achse) | Eigene Vollzählung der Aufrufstellen von `FenceToggle`, `FenceRun`, `FenceCloses` und `TrimFenceIndent` im Nicht-Test-Code: genau vier Dateien führen eine eigene Fence-Zustands-Schleife — `markdown.go`, `trace_table.go`, `sections.go`, `spans.go`. Die Aufzählung in §4a nennt jetzt alle vier. Auf der **Überschriften**-Achse ist dieselbe Aufzählung weiter unvollständig ⇒ N-3 |
| L-8 (LOW) Nutzer-Doku auf der alten Paarung | **geheilt, mit neuem Defekt** | Alle vier gemeldeten Stellen (`docs/user/benutzerhandbuch.md:1142`, die §6-Modul-Tabelle, `README.de.md:104`, `README.md:103`) tragen die Fence-Grenze, und `docs/user/releasing.md:67` bekommt die fehlende Checklisten-Zeile. Die neue Formulierung ist aber an beiden Stellen sachlich falsch ⇒ N-4, und dieselbe neue Checklisten-Regel ist für `versions`/`pins` nicht angewandt ⇒ N-8 |
| L-9 (LOW) doctor-Klartext halbiert | **geheilt** | `--doctor` liefert jetzt „Anker entspricht keinem Heading-Slug und keinem HTML-Anker der Zieldatei“, wortgleich mit der Befund-Meldung beider Emitter (`anchors.go:241`, `codepaths.go:264`) und mit der §4-Grund-Code-Tabelle. Vier Fundstellen, ein Wortlaut |
| L-10 (LOW) verwaistes Pronomen | **geheilt** | Der Einschub steht jetzt vor beiden Zweigen; das „er“ folgt unmittelbar auf „inline-Zitat-Span“. Volltext von Schritt 2 gegen die Vorfassung gelesen |
| L-11 (LOW) Link/Prosa-Widerspruch in der BEO-003-Zelle | **geheilt** | `docs/plan/planning/observations.md` sagt in derselben Zelle jetzt „liegt … in Arbeit“ statt „in `open/`“; Link und Prosa stimmen überein. Der Zähler steht unverändert auf 2 — das ist Wellen-Closure, nicht Slice-DoD, und N-3 gibt ihm jetzt einen Anlass |
| L-12 (INFO) Grundmenge zählt Erwähnungen | **teilweise geheilt** | §4a des Slice, die Historien-Zeile 0.58.0 und der `CHANGELOG.md`-Eintrag sind korrigiert. [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) §Kontext trägt weiterhin „null von 18 `d-check:cite`-Direktiven“ ⇒ N-9 |
| L-13 (INFO) äquivalenter Mutant | **geheilt** | Die Wache ist in beiden Sammlern entfernt. Randfall geprüft statt behauptet: Direktive in **Zeile 1** (also der kleinstmögliche Index, `j == 1`) im Blockquote- **und** im Absatz-Zweig — kein Panik-Pfad, und das Ergebnis ist identisch zum Vor-Heilungs-Bild. Der Zugriff bleibt sicher, weil `j` immer mindestens 1 ist und die Eintritts-Prüfung den Schritt bis `j` bereits ausgewertet hat |

---

## Neue Findings

### N-1 — „Dieselbe Anker-Erkennung wie in `DC-FA-ANCH-001`" gilt weiter nur für die HTML-Hälfte; zwei Achsen geben in **einem** Lauf zwei Antworten

- `kategorie`: HIGH
- `quelle`: [`DC-FA-VER-001`](../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
  (neuer Absatz: „welchen Span `versions.current-from` adressiert, entscheidet
  **dieselbe Anker-Erkennung wie in** [`DC-FA-ANCH-001`](../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
  — ein Heading-Slug oder ein Inline-HTML-Anker") und
  [`DC-FA-PIN-001`](../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in)
  (gleichlautender Einschub) gegen
  §[`DC-FA-ANCH-001.a`](../../spec/spezifikation.md#dc-fa-anch-001a--github-slug-algorithmus)
  Schritt 5 (Duplikat-Suffixe `-1`, `-2`, …) und
  §[`DC-FA-ANCH-001.b`](../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker)
  („Der Fragment-Vergleich **nach RFC-3986-Dekodierung**");
  [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 1 und 2; Reviewer-Skill §HIGH und §MEDIUM
- `pfad`: `internal/hexagon/core/rules/versions.go:121`–`136` (`slugSection`)
  gegen `internal/hexagon/core/rules/anchors.go:84`–`98` (`HeadingSlugs`); und
  `internal/hexagon/core/rules/anchors.go:201` (`url.PathUnescape`) ohne Pendant
  in `internal/hexagon/core/rules/versions.go:71`–`74` bzw.
  `internal/hexagon/core/rules/pins.go:101`–`104`
- `befund`: Der Commit vereinheitlicht die **HTML**-Hälfte vollständig — das ist
  nachvollzogen und hält. Er lässt aber die **Slug**-Hälfte unangetastet und
  erweitert zugleich die Vertragszusage von „erkannt wie in `DC-FA-ANCH-001.b`"
  auf „dieselbe Anker-Erkennung wie in `DC-FA-ANCH-001`", also gerade auf die
  Hälfte, die weiter divergiert. Zwei Eingaben genügen: (a) eine Datei mit zwei
  gleichnamigen Überschriften, adressiert über den Duplikat-Slug `#alt-1` — für
  `anchors` ein gültiger Anker (0 Befunde), für `slugSection` keiner, weil dort
  ein nackter `Slugify`-Vergleich ohne den Duplikat-Zähler steht; (b) ein
  prozent-kodiertes Fragment `#a%20b` auf `id="a b"` — `anchors` dekodiert und
  trifft, `versions`/`pins` vergleichen die kodierte Zeichenkette. In beiden
  Fällen bricht **derselbe** Lauf mit Exit 2 („Anker nicht auflösbar") ab,
  während `anchors` denselben Anker für gültig hält, und `pins` verliert seinen
  Drift-Schutz kommentarlos. Das ist wörtlich die Bauform, die L-1 gemeldet hat,
  eine Ebene weiter — und die neue Zusage schließt die Suche erneut.
- `verifizierbar`: ja — zwei Fixtures in einem Scratch-Verzeichnis außerhalb des
  Repos, je ein Lauf mit `--enable anchors --enable versions --enable pins`:
  Exit 2 mit „versions.current-from: Anker #alt-1 nicht auflösbar" bzw.
  „… Anker #a%20b nicht auflösbar", während der Lauf mit nur `anchors` 0 Befunde
  meldet. Gegenprobe mit dem Basis-Slug `#alt` derselben Datei: Exit 1 mit
  `link-stale` **und** `version-stale`, also volle Auflösung. Beide Divergenzen
  liegen auch im Vor-Heilungs-Bild vor; neu ist die Zusage, nicht das Verhalten.
  Reichweite gemessen: über die drei Repos existiert heute **eine** produktive
  Konsumenten-Referenz (`version.md#aktuell`, ein Heading-Slug ohne Duplikat) und
  ein Doku-Beispiel — die Divergenz ist latent, wie der Rest des Slice.
- `klasse`: „geteilte Lexik, halb übernommen"

### N-2 — Die Vereinheitlichung macht die Anker-Erkennung case-sensitiv; das ist eine **dritte**, ungenannte „findet weniger"-Stelle, und bei `pins` ist sie still

- `kategorie`: HIGH
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  §Konsequenzen („**Weniger:** … und ein `dpin` …" — zwei Stellen, geschlossen
  aufgezählt), Historien-Zeile 0.58.0 („findet **weniger** an **zwei**
  konstruierbaren Stellen") und 0.58.1 („**zwei** Fälle verlieren einen Befund,
  einer davon still"), `CHANGELOG.md` („Sie **findet weniger** an **zwei**
  Stellen"); [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus);
  Reviewer-Skill §HIGH (Stilles-Grün-Pfad) mit Kontext-Eskalation (dritte
  Wiederholung derselben Klasse in derselben Sitzung)
- `pfad`: `internal/hexagon/core/rules/versions.go:144` (`htmlAnchorLines(content)[anchor]`,
  wörtlicher Map-Zugriff) gegen die von diesem Commit entfernte Regex derselben
  Funktion, die den Ankernamen unter einem `(?i)`-Flag über dem **gesamten**
  Muster verglich
- `befund`: Vor der Heilung löste `versions.current-from` mit dem Fragment
  `#aktuell` einen Anker `id="Aktuell"` auf; danach nicht mehr. Die neue Antwort
  ist die richtige und von
  §[`DC-FA-ANCH-001.b`](../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker)
  ausdrücklich gedeckt („wörtlich … HTML-Fragment-Auflösung ist case-sensitiv") —
  gemeldet wird nicht das Verhalten, sondern dass **kein** konsumentensichtbarer
  Rand diese Änderung nennt, obwohl alle vier Ränder ihre Richtungs-Aufzählung in
  diesem Commit neu geschrieben und dabei geschlossen formuliert haben. Die
  Änderung wirkt in beide Richtungen: `versions` kippt von Exit 0 auf Exit 2
  (findet mehr), und `pins` verliert seinen `link-stale`-Befund **ohne jede
  Ausgabe** (findet weniger, still) — genau die Ausprägung, die L-2 als die
  unangenehmere benannt hat. Sie fällt außerdem nicht unter die vier Zeichenfolgen,
  die der Commit als „sieht nur wie ein Anker aus" aufzählt: hier ist ein echter
  Anker im Spiel, nur anders geschrieben.
- `verifizierbar`: ja — ein Fixture (`id="Aktuell"`, Referenz und `current-from`
  auf `#aktuell`) gegen beide Images: Vor-Heilungs-Build `versions` Exit 0 /
  0 Befunde und `pins` Exit 1 `link-stale`; HEAD-Build `versions` Exit 2 („Anker
  #aktuell nicht auflösbar") und `pins` Exit 0 / 0 Befunde. `anchors` meldet in
  beiden Fassungen `anchor-missing`.
- `klasse`: „Verhaltensänderung eines ausgelieferten Moduls ohne gedeckte Richtung"

### N-3 — Das Modul `planning` beantwortet „ist das eine Überschrift" zweimal **roh**; die Aufzählung, die den Klassen-Abschluss trägt, nennt beide Stellen nicht

- `kategorie`: MEDIUM
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 2 („Wer … fragt ‚ist das eine **Überschrift**‘ … bekommt die
  geteilte Antwort") und §Re-Evaluierungs-Trigger („Eine **vierte** Stelle
  beantwortet eine Lexik-Frage selbst. Dann genügt die Aufzählung als Beleg nicht
  mehr"); §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C3 („**außerhalb** eines Fenced-Code-Blocks — eine Raute-Zeile in einem
  Beispielblock ist keine Überschrift") gegen Schritt 3/4 derselben Datei;
  §4a des [Slice-Plans](../plan/planning/in-progress/slice-103-geteilte-lexik-raender.md)
  („**Übrig bleiben genau zwei** Stellen, die eine Lexik-Frage **roh**
  beantworten"); Reviewer-Skill §MEDIUM (Konsistenz-Lücke zwischen Modulen
  derselben Eingabe-Klasse)
- `pfad`: `internal/hexagon/core/rules/planning.go:357`–`367`
  (`planningHeadingLine`: `strings.TrimRight`-Gleichheit über **alle** Rohzeilen)
  und `internal/hexagon/core/rules/planning.go:373`–`385`
  (`planningBlockHasMarker`: `strings.HasPrefix` auf `## ` über Rohzeilen) gegen
  `internal/hexagon/core/rules/planning.go:179`–`191` (`closureHeadingLine`) und
  `internal/hexagon/core/rules/sections.go:19`–`36` (`FindSectionHeads`)
- `befund`: Dieselbe Datei, dieselbe Frage, zwei Antworten — und diesmal
  innerhalb **eines** Moduls. Die Closure-Fähigkeit nimmt die geteilte,
  fence-bewusste Antwort und dokumentiert das ausdrücklich; der Aktiv-Status-Guard
  daneben zählt Rohzeilen. Die Folge ist an zwei Stellen ein **Falsch-Rot**: eine
  Roadmap, die ihren eigenen Abschnitt in einem Beispiel-Fence zeigt, gilt als
  „Überschrift mehrfach vorhanden" und meldet `planning-drift`; eine Roadmap, in
  deren Aktiv-Block ein Beispiel-Fence eine Raute-Zeile trägt, verliert den
  Ruhe-Marker hinter dem Fence und meldet ebenfalls `planning-drift`. Beides ist
  älter als dieser Diff — gemeldet wird, dass die Aufzählung, mit der §4a den
  Klassen-Abschluss **belegt**, diese beiden Stellen nicht kennt und die Zusage
  „die Klasse ist geschlossen" damit nicht trägt. Nach dem Wortlaut von
  [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) ist das
  der erste Re-Evaluierungs-Trigger, und für den Closure-Trigger des
  [Wellendokuments](../plan/planning/welle-74-geteilte-lexik-raender.md) („BEO-003
  ist entschieden — geschlossen oder auf 3") ist es der Anlass, den der Zähler
  bisher nicht hatte.
- `verifizierbar`: ja — zwei Fixtures außerhalb des Repos, je mit
  `--enable planning`. (a) Roadmap mit kanonischer Überschrift, Ruhe-Marker und
  einem Markdown-Fence, der dieselbe Überschrift als Beispiel zeigt, ohne
  Slice-Datei im Verzeichnis: `planning-drift`, Exit 1; ohne den Beispiel-Block
  0 Befunde, Exit 0. (b) Roadmap mit Ruhe-Marker **hinter** einem Fence, der eine
  Raute-Zeile enthält: `planning-drift`, Exit 1; derselbe Marker **vor** dem
  Fence: 0 Befunde, Exit 0. Das veröffentlichte Bild `v0.52.0` liefert dieselben
  Zahlen, der Defekt ist also älter als die Range.
- `klasse`: „Aufzählung als Beleg, aus dem Gedächtnis abgeleitet"

### N-4 — Die neue Absatz-Zusage in Anforderung, Handbuch und beiden READMEs ist gegen einen Lauf falsch: eine Leerzeile trennt gerade **nicht**

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-CITE-001`](../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in)
  §Beschreibung (neu); [`AGENTS.md`](../../AGENTS.md) §2 (das Lastenheft ist
  vertraglich abnahmebindend); §[`DC-FA-CITE-001.a`](../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
  Schritt 2 („Ist die **nächste nicht-leere Zeile** ein Größer-Blockquote");
  Reviewer-Skill §MEDIUM (Spec-Treue-Lücke einer Messmethode)
- `pfad`: `spec/lastenheft.md:1351`–`1353`, `docs/user/benutzerhandbuch.md:1142`–`1145`,
  `README.de.md:104`–`105`, `README.md:103`–`104` gegen
  `internal/hexagon/core/rules/citations.go:158`–`165` (`citationQuote`
  überspringt Leerzeilen, bevor der Kandidat bestimmt wird)
- `befund`: Der neue Satz lautet „**Unmittelbar heißt: im selben Absatz.** Ein
  Fenced-Block zwischen Direktive und Zitat trennt **genauso wie eine
  Leerzeile**". Beide Hälften sind falsch. Die Direktive paart **über**
  Leerzeilen hinweg, also über Absatzgrenzen hinweg — das ist keine Randlage,
  sondern der dokumentierte Normalfall des Blockquote-Zweigs („die nächste
  nicht-leere Zeile"), und der Erst-Report hat ihn als Negativbefund
  ausdrücklich festgehalten (die naive Ein-Schritt-Form würde „jede Direktive mit
  mehr als einer Leerzeile Abstand fälschlich trennen"). Eine Leerzeile trennt
  in dieser Position also überhaupt nicht, und der Vergleich, mit dem die neue
  Fence-Regel erklärt wird, sagt dem Leser das Gegenteil. Dieselbe Konstruktion
  steht im Handbuch („er ist eine Absatzgrenze wie eine Leerzeile, und dann folgt
  der Direktive kein Zitat") — dort im Widerspruch zum Satz zwei Zeilen darüber —
  und in beiden READMEs, die die Fehlaussage sogar in ihre Kurzform heben
  („markiert das folgende Zitat **im selben Absatz**"). Das Failure-Szenario ist
  der Adopter, der nach dem Lesen einen Zeilenabstand einbaut, um die Paarung zu
  lösen, und einen `citation-mismatch` bekommt.
- `verifizierbar`: ja — ein Fixture außerhalb des Repos mit `--enable citations`:
  Direktive, **zwei** Leerzeilen, dann ein Absatz mit inline-Zitat. Passender
  Zitattext ⇒ 0 Befunde, Exit 0; abweichender Zitattext ⇒ `citation-mismatch`,
  Exit 1 — die Paarung findet also über die Absatzgrenze hinweg statt. Ebenso im
  Blockquote-Zweig (Direktive, Leerzeile, abweichendes Größer-Zitat ⇒
  `citation-mismatch`). Die Anforderung beschreibt beide Läufe falsch.
- `klasse`: „Rand behauptet mehr als der Lauf"

### N-5 — Die geteilte Anker-Menge ist mit „nie ein Falsch-Befund" begründet; bei ihrem neuen Konsumenten `versions` erzeugt sie genau einen

- `kategorie`: MEDIUM
- `quelle`: §[`DC-FA-ANCH-001.b`](../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker)
  (Schluss-Absatz: „Eine zu groß geratene Anker-Menge führt nur dazu, dass das
  Modul **schweigt** (nie ein Falsch-Befund); die permissive Richtung ist daher
  bewusst und mit dem Schweige-Charakter des Moduls konsistent") gegen den neuen
  Einschub in [`DC-FA-VER-001`](../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
  der `versions` an dieselbe Menge bindet; Reviewer-Skill §MEDIUM (Modul-Grenze
  auf der Ziel-Achse: welche Eingaben liest dieses Modul, die es nicht scannt —
  und gilt dort dieselbe Zusage?)
- `pfad`: `internal/hexagon/core/rules/anchors.go:123`–`142` (`htmlAnchorLines`,
  ohne Kommentar- und Indented-Code-Behandlung) konsumiert von
  `internal/hexagon/core/rules/versions.go:143`–`156` (`htmlAnchorSection`) und
  `internal/hexagon/core/rules/pins.go:128`–`148` (`spanHash`)
- `befund`: Die Erkennung liest HTML-Tags auf allen Prosa-Zeilen. Zwei Formen,
  die GitHub nicht als Sprungziel rendert, zählen darum mit: ein Anker in einer
  **HTML-Kommentar**-Zeile und ein Anker in einem **eingerückten** Code-Block
  (die Vorverarbeitung kennt nur Fenced- und Inline-Code). Für `anchors` ist das
  folgenlos — die Menge ist zu groß, das Modul schweigt. Für `versions` ist es
  das nicht: `htmlAnchorLines` liefert das **erste** Vorkommen, der Span beginnt
  dort, und `versions` zieht seine „aktuelle Version" aus einem auskommentierten
  bzw. eingerückten Beispiel. Ergebnis ist ein `version-stale`-Befund auf einem
  Pin, der die tatsächlich aktuelle Version trägt — ein Falsch-Rot in genau der
  Richtung, die der Begründungstext ausschließt. Der Diff hat das Verhalten nicht
  geändert; er hat die **Begründung** importiert, die dafür nicht gilt.
- `verifizierbar`: ja — zwei Fixtures außerhalb des Repos, je eine Zieldatei mit
  einem echten Anker und einem gleichnamigen Anker davor (einmal in einem
  HTML-Kommentar, einmal in einem mit vier Leerzeichen eingerückten Block), beide
  mit einer älteren Version im Kommentar-/Code-Beispiel: `--enable versions`
  meldet `version-stale`, Exit 1, obwohl der Pin die Version des echten
  Anker-Spans trägt; `--enable anchors` meldet 0 Befunde. Ohne die vorgelagerte
  Zeile: 0 Befunde, Exit 0. Vor-Heilungs-Bild identisch.
- `klasse`: „Zusage des Erzeugers gilt beim neuen Konsumenten nicht"

### N-6 — Die Akzeptanzkriterien der drei geänderten Anforderungen stehen unverändert auf der alten Fassung

- `kategorie`: LOW
- `quelle`: [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel-Zeile „Anforderung — `spec/lastenheft.md` (Beschreibung **und**
  Akzeptanzkriterien)"); Lastenheft-Historie 0.57.1 und 0.57.2, die dieselbe
  Klasse an derselben Stelle schon zweimal saniert haben; Reviewer-Skill §LOW
- `pfad`: die Abschnitte **Akzeptanzkriterien** von
  [`DC-FA-CITE-001`](../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in),
  [`DC-FA-VER-001`](../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
  und [`DC-FA-PIN-001`](../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in)
- `befund`: Der Commit heilt L-5 in den **Beschreibungen** und lässt die
  Akzeptanzkriterien aller drei Anforderungen unberührt. Bei
  `DC-FA-VER-001` nennt das Kriterium „Negative" ausdrücklich den Pin „auch
  innerhalb eines Fenced-Code-Blocks" — also genau die Roh-Lesung, deren
  Abgrenzung zum Anker der neue Beschreibungstext einführt —, ohne die
  Anker-Seite zu erwähnen. Bei `DC-FA-CITE-001` deckt „Fail-closed … dem
  fehlenden folgenden Zitat vorbehalten" den neuen Fall nur generisch ab; das
  Trio benennt ihn nicht. Das Failure-Szenario ist die Abnahme: wer gegen die
  Kriterien abnimmt statt gegen die Beschreibung, prüft die geänderte Semantik
  nicht.
- `verifizierbar`: ja — Diff der Range auf `spec/lastenheft.md` zeigt in den drei
  Anforderungen nur Beschreibungs-Absätze; Volltext der drei Kriterien-Listen
  gelesen, keine Fence-, Absatz- oder Anker-Bedingung darin.
- `klasse`: „Semantik im Körper geändert, der Rand referiert die andere Fassung"

### N-7 — Die Spezifikation ändert einen normativen Schritt und eine Grund-Code-Zeile, ohne dass ihre §7-Historie eine Zeile bekommt

- `kategorie`: LOW
- `quelle`: `spec/spezifikation.md` §7 Historie (eine Zeile je Änderung, auch für
  Review-Nachzüge — die Einträge vom 2026-08-15 und die vier vom 2026-08-10
  belegen die Granularität); [`AGENTS.md`](../../AGENTS.md) §2; Reviewer-Skill §LOW
- `pfad`: `spec/spezifikation.md` §7 (Zeile 2508 trägt weiter nur den Eintrag des
  Vorgänger-Commits) gegen die in dieser Range geänderten Stellen
  §[`DC-FA-CITE-001.a`](../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
  Schritt 2 und die §4-Zeile zu `anchor-missing`
- `befund`: Der Heilungs-Commit fügt der Spezifikation einen dritten
  Blockquote-Terminator hinzu und ändert den Klartext einer Grund-Code-Zeile —
  beides normativ. Das Lastenheft bekommt für denselben Vorgang eine eigene
  Versionszeile 0.58.1; die Spezifikations-Historie bekommt nichts, und ihre
  vorhandene Zeile vom selben Tag beschreibt ausschließlich den Stand **vor** der
  Heilung. Wer die §7-Chronik als Delta-Register liest, findet die
  Terminator-Änderung nirgends.
- `verifizierbar`: ja — Diff der Range auf `spec/spezifikation.md`: zwei
  Hunks, beide außerhalb von §7.
- `klasse`: „Semantik im Körper geändert, der Rand referiert die andere Fassung"

### N-8 — Die neue Release-Prep-Regel („eine geänderte Zusage zählt wie ein neues Feature") wird im selben Commit für `versions` und `pins` nicht angewandt

- `kategorie`: LOW
- `quelle`: `docs/user/releasing.md:67`–`73` (neu, in diesem Commit);
  [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel „Nutzer-Doku"); Reviewer-Skill §LOW (Doku-Drift)
- `pfad`: `docs/user/benutzerhandbuch.md` §5 (der Abschnitt zum Versions-Register
  und der Abschnitt zum Content-Pin), die §6-Modul-Tabellenzeilen zu `versions`
  und `pins`, `README.de.md:60`–`67` und `README.md:61`–`68`
- `befund`: Der Commit schreibt die Regel auf, dass eine **geänderte** Zusage
  mitten in einem bestehenden Abschnitt nachzuziehen ist — und zieht sie
  anschließend nur für `citations` nach. Die Anker-Semantik von `versions` und
  `pins` hat sich in derselben Range zweimal geändert (kein Anker in Fences, dann
  kein Anker in Inline-Code / ohne Tag / an `data-id`), und keine der genannten
  Stellen erwähnt eine Grenze. Der Handbuch-Abschnitt zum Versions-Register
  empfiehlt Nutzern ausdrücklich den expliziten HTML-Anker als Verweispunkt und
  zeigt ihn in Beispielform — er ist damit der Ort, an dem die neue Grenze
  gebraucht wird. Bewusst LOW: der Stand ist unreleased, und das Repo zieht
  Handbuch und READMEs planmäßig im Release-Prep nach.
- `verifizierbar`: ja — Volltext der fünf Stellen gelesen; die Zeichenfolgen
  „Fence", „Code-Block" und „Inline-Code" kommen dort nicht in Bezug auf die
  Anker-Erkennung vor. Diff der Range berührt in beiden READMEs nur die
  `citations`-Aufzählungspunkte.
- `klasse`: „Rand auf der alten Fassung stehengeblieben"

### N-9 — Der ADR-Kontext trägt weiter die Grundmenge, die der Commit an drei anderen Stellen korrigiert

- `kategorie`: LOW
- `quelle`: L-12 des [Erst-Reports](2026-08-16-slice-103-geteilte-lexik-review.md);
  Reviewer-Skill §LOW
- `pfad`: `docs/plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md:22`
- `befund`: Der Commit korrigiert die Grundmenge im Slice (§4a), in der
  Historien-Zeile 0.58.0 und im `CHANGELOG.md`-Eintrag; der §Kontext-Absatz von
  [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  behauptet weiter „null von 18 `d-check:cite`-**Direktiven**" — die Formulierung,
  die L-12 als unbelegt gemeldet hat. Damit trägt ausgerechnet das
  Entscheidungs-Dokument die schwächste Fassung der Messung, auf der es beruht.
- `verifizierbar`: ja — Diff der Range auf der ADR-Datei zeigt drei Hunks
  (Entscheidung 3, Konsequenzen, Re-Evaluierungs-Trigger), keinen im Kontext.
- `klasse`: „Rand auf der alten Fassung stehengeblieben"

### N-10 — Die Historien-Zeile 0.58.0 ist nachträglich umgeschrieben, nicht nur ergänzt worden

- `kategorie`: INFO
- `quelle`: `spec/lastenheft.md` §Historie; Reviewer-Skill §INFO
  (dokumentationswürdige Annahme)
- `pfad`: `spec/lastenheft.md:2643`–`2644`
- `befund`: Der Commit fügt die Zeile 0.58.1 hinzu **und** ersetzt zugleich den
  Richtungs-Satz der Zeile 0.58.0. Die Historie ist damit an dieser Stelle kein
  Append-Log mehr: 0.58.0 sagt heute etwas anderes, als die Fassung 0.58.0 gesagt
  hat. Sachlich ist die Korrektur richtig und folgenlos, weil 0.58.0 nie
  released wurde und beide Zeilen denselben Tag tragen. Festgehalten, damit ein
  späterer Leser die Doppelung („beide Richtungen" in 0.58.0 **und** „die Zusage
  war falsch" in 0.58.1) nicht für einen Widerspruch hält.
- `verifizierbar`: ja — Diff der Range auf `spec/lastenheft.md` zeigt die
  0.58.0-Zeile als geändert, nicht als unverändert.
- `klasse`: „Append-Log nachträglich editiert"

---

## Negativbefunde

- geprüft, ohne Befund: **Die HTML-Hälfte der Anker-Frage ist wirklich
  vereinheitlicht.** Sieben eigene Fixtures, die der Erst-Report nicht hatte —
  Groß-/Kleinschreibung, zwei gleichnamige Anker, Anker in einem eingerückten
  Code-Block, Anker in einer HTML-Kommentar-Zeile, `id` und `name` am selben
  Element, Anker hinter einem unbalancierten Fence, Anker in einem geschlossenen
  Fence — geben bei `anchors`, `versions` und `pins` in **jedem** Fall dieselbe
  Antwort. Die verbleibenden Unterschiede liegen auf der Slug- und der
  Fragment-Dekodierungs-Achse (N-1), nicht im HTML-Zweig.
- geprüft, ohne Befund: **Die Reihenfolge „erstes Vorkommen gewinnt" ist
  unverändert.** Fixture mit zwei gleichnamigen HTML-Ankern, deren Spans
  verschiedene Versionen tragen: Vor-Heilungs-Bild und HEAD melden beide
  `version-stale` mit demselben Ziel-Wert, wählen also denselben (den ersten)
  Anker. Das gilt auch innerhalb einer Zeile, weil die Sammelfunktion je Name
  nur beim **ersten** Treffer schreibt.
- geprüft, ohne Befund: **Das Entfernen der Wache in beiden Sammlern zerstört
  keinen Randfall.** Der kleinste erreichbare Index ist `j == 1` (die Direktive
  steht in Zeile 1, das Zitat in Zeile 2); ein negativer Index ist konstruktiv
  ausgeschlossen. Beide Zweige mit dieser Eingabe gefahren: kein Abbruch,
  identisches Ergebnis zum Vor-Heilungs-Bild. Die Eintritts-Prüfung wertet den
  Schritt bis `j` bereits aus, der zusätzliche Test am Schleifenanfang ist damit
  konstant falsch.
- geprüft, ohne Befund: **Kein Bestands-Delta durch die Heilung.** Vor-Heilungs-
  und HEAD-Build über die drei Repos (`d-check`, `a-check`, `ai-harness-course`)
  mit `anchors`, `pins`, `versions`, `spans` und mit `citations`: byte-identische
  Ausgaben und Exit-Codes an jeder Stelle. Die Latenz-Aussage des Slice hält auch
  für die zweite Reparatur-Runde.
- geprüft, ohne Befund: **Die gescopten Roh-Lesungen sind weiterhin unberührt.**
  Der Pin-Scan von `versions` läuft über alle Rohzeilen einschließlich Fences
  ([ADR-0019](../plan/adr/0019-versions-pin-fence-ausnahme.md)), der gehashte
  Ziel-Span von `pins` bleibt roh
  ([ADR-0020](../plan/adr/0020-content-pin-fence-ausnahme.md)); der Diff berührt
  in beiden Modulen nur die Anker-**Erkennung**, nicht den Span. Fixture-Beleg:
  ein Pin innerhalb eines Fence wird weiterhin gefunden.
- geprüft, ohne Befund: **Kein dritter Anker-Automat.** Nach der Heilung gibt es
  genau **eine** HTML-Anker-Erkennung; `htmlAnchors` ist ihre Namens-Sicht,
  `htmlAnchorSection` ihre Zeilen-Sicht. Vollzählung der Aufrufstellen bestätigt
  das; eine zweite Regex auf Ankernamen existiert im Nicht-Test-Code nicht mehr.
- geprüft, ohne Befund: **Die Fence-Achse der Klassen-Aufzählung ist vollständig.**
  Eigene Vollzählung über `FenceToggle`, `FenceRun`, `FenceCloses` und
  `TrimFenceIndent`: vier Dateien, alle vier in §4a genannt, alle vier über die
  geteilten Prädikate. Die Lücke liegt auf der Überschriften-Achse (N-3).
- geprüft, ohne Befund: **Die zwei neuen `citations`-Tests sind mutations-echt.**
  Je Rückbau genau ein roter Test, kein Kollateral-Rot; beide über Dateikopien,
  Arbeitsbaum danach wiederhergestellt.
- geprüft, ohne Befund: **Der `--doctor`-Spiegel ist vollständig geschlossen.**
  Vier Fundstellen des `anchor-missing`-Klartexts (beide Emitter, `reasonTexts`,
  §4-Tabelle) tragen denselben Wortlaut.
- geprüft, ohne Befund: **Der neue `vcs`-Trigger ist im eigenen Repo wirklich
  beobachtbar.** `spans` steht in der `modules`-Liste der `.d-check.yml`, und die
  stille Richtung setzt die offene Fence in der HEAD-Fassung voraus — also in
  genau der Fassung, die der Arbeitsbaum-Scan sieht. Der Trigger ist damit kein
  Platzhalter.
- geprüft, ohne Befund: **Hard Rules.** Kein Spec-Stratum nennt eine ADR-Kennung;
  [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) steht auf
  `Proposed`, ist im Index eingetragen und trägt den geforderten
  Re-Evaluierungs-Trigger-Abschnitt; keine Inline-Suppression, keine
  Gate-Lockerung, kein Netzzugriff außerhalb der Netz-Module; die
  Import-Richtung des Kerns ist unverändert (nur `model` und Ports).
- geprüft, ohne Befund: **Referenz-Richtung (SDP) und Marker-Ehrlichkeit.** Die
  Range führt keinen neuen Provenance-Marker ein; die Slice-Kennung steht in der
  ADR ausschließlich im ausgenommenen `## Geschichte`-Abschnitt. `matrix` ist im
  grünen Gate-Lauf aktiv.
- geprüft, ohne Befund: **Keine neue Vertragsfläche unterschlagen, wo eine
  entsteht.** Der Diff führt keinen neuen Grund-Code, keinen Config-Schlüssel und
  kein Modul ein; §2-Schema und die `--print-config`-Vorlage bleiben zu Recht
  unverändert. Die berührten Spiegel sind Anforderung (N-6), Algorithmus-Historie
  (N-7), Nutzer-Doku (N-4/N-8) und ADR (N-9).
- geprüft, ohne Befund: **Gate-Stand.** `make gates` Exit 0 (373 Dateien, 0
  Befunde: `doc-check`, `lint`, `test`, `arch-check`, `coverage-gate`, `semgrep`,
  `gate-consistency`, `planning-check`) und `make verify-closure-notes` Exit 0
  (343/0), beide nach der Wiederherstellung des Arbeitsbaums.
- geprüft, ohne Befund: **Arbeitsbaum.** Jede Mutation aus der Scratch-Kopie
  zurückgeschrieben, kein `git checkout`; `git status --short` ist am Ende leer
  (dieser Report ist die einzige neue Datei), und der Neubau liefert dieselbe
  Image-ID wie vor dem ersten Eingriff.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 3 |
| LOW | 4 |
| INFO | 1 |

**Erst-Befunde:** 9 geheilt · 3 teilweise geheilt (L-1, L-2, L-12) · 0 nicht
geheilt · 0 gegenstandslos. Ein weiterer (L-8) ist geheilt, hat aber einen neuen
Defekt eingeführt.

**Finding-Klassen dieses Laufs:** geteilte Lexik, halb übernommen ·
Verhaltensänderung eines ausgelieferten Moduls ohne gedeckte Richtung ·
Aufzählung als Beleg, aus dem Gedächtnis abgeleitet · Rand behauptet mehr als der
Lauf · Zusage des Erzeugers gilt beim neuen Konsumenten nicht · Semantik im Körper
geändert, der Rand referiert die andere Fassung · Rand auf der alten Fassung
stehengeblieben · Append-Log nachträglich editiert.

**Wiederholungs-Signal.** Beide HIGH sind Wiedergänger der beiden HIGH, die sie
schließen sollen: N-1 ist L-1 eine Ebene weiter (die geteilte Antwort wurde
wieder nur zur Hälfte übernommen, diesmal HTML statt Fence), N-2 ist L-2 in
derselben Bauform (eine geschlossen formulierte Richtungs-Aufzählung, der genau
die stille Stelle fehlt). Das ist die **dritte** Runde derselben zwei Klassen an
demselben Vorhaben — nach dem Reviewer-Skill §Kontext-Eskalation ein
Steering-Loop-Signal: die Frage ist nicht mehr, ob dieser Commit sie schließt,
sondern was sie **mechanisch** schließt. N-3 liefert dafür den Anlass, den das
Beobachtungs-Register braucht.

## Verdikt

**Merge-blockierend:** ja — zwei HIGH und drei MEDIUM.

**Was die Heilung geliefert hat, ist echt.** Neun der dreizehn Befunde sind ohne
Einschränkung geheilt, und die Belege halten unabhängiger Nachprüfung stand: die
vier gemeldeten Anker-Zeichenfolgen antworten jetzt in allen drei Modulen gleich,
die zwei neuen `citations`-Tests fangen ihren jeweiligen Rückbau und nur ihn, der
äquivalente Mutant ist ohne Randfall-Schaden entfernt, der `--doctor`-Klartext ist
an allen vier Fundstellen wortgleich, und die `vcs`-Grenze beschreibt jetzt die
stille Richtung mit einem Trigger, der im eigenen Gate tatsächlich feuern kann.
Der Bestand ist über drei Repos byte-identisch geblieben — die Latenz-Aussage des
Slice trägt auch nach der zweiten Runde.

**Die Klasse ist trotzdem nicht geschlossen, und zwar an drei Stellen innerhalb
dessen, was der Commit zu schließen behauptet.** Die Anker-Frage ist jetzt zur
HTML-Hälfte vereinheitlicht statt zur Fence-Hälfte; die Slug-Hälfte und die
Fragment-Dekodierung geben weiter zwei Antworten, und der Vertragstext hat seine
Zusage in derselben Bewegung von `DC-FA-ANCH-001.b` auf die **ganze**
`DC-FA-ANCH-001` erweitert — also gerade auf die Hälfte, die divergiert (N-1).
Die Vereinheitlichung selbst hat eine dritte, ungenannte Richtungsänderung
eingeführt, die bei `pins` still ist (N-2) — in genau dem Commit, dessen
Kernanspruch „beide Richtungen benannt" lautet. Und die Aufzählung, mit der §4a
den Klassen-Abschluss belegt, kennt zwei rohe Überschriften-Antworten im Modul
`planning` nicht, die beide als Falsch-Rot beobachtbar sind (N-3).

**Die neu geschriebenen Ränder sind ungleich verlässlich.** ADR, Historie und
`CHANGELOG` tragen die Richtungsfrage jetzt ehrlich; die **Anforderung** und die
Nutzer-Doku haben sie mit einem Satz eingeführt, der gegen einen Lauf falsch ist
(N-4), und die Akzeptanzkriterien derselben drei Anforderungen stehen unverändert
da (N-6) — dieselbe Halbierung wie in L-5, eine Ebene tiefer.

**Release-Empfehlung: noch nicht releasen.** Die Einordnung als **Minor** bleibt
richtig; ihre Begründung ist erneut unvollständig, weil eine dritte Stelle
weniger findet und dabei schweigt. Vor den Tag gehören N-1, N-2 und N-4: das
erste, weil der Vertrag sonst eine Parität zusagt, die es nicht gibt, und die
Suche damit ein drittes Mal schließt; das zweite, weil ein Konsument nach dem
Update einen `pins`-Drift-Befund verliert, ohne dass ihm irgendetwas davon sagt;
das vierte, weil eine falsche normative Aussage im abnahmebindenden Lastenheft
und in beiden READMEs steht. N-3 und N-5 sind Vertrags- und Entscheidungs-Arbeit
auf demselben Diff — N-3 zusätzlich die Entscheidung, die der Closure-Trigger der
Welle für **BEO-003** ohnehin verlangt. N-6, N-7, N-8, N-9 sind klein und billig;
N-8 ist Release-Prep-Material.

**Übergabe:** Die Findings gehen an den Implementer. Für die Slice-Closure und
das Beobachtungs-Register ist die tragende Beobachtung dieser Runde nicht ein
einzelner Befund, sondern dass **beide** HIGH ihre eigenen Vorgänger wiederholen:
eine geteilte Antwort ist zweimal hintereinander halb übernommen worden, und eine
Richtungs-Aufzählung ist zweimal hintereinander geschlossen formuliert und
unvollständig gewesen. N-3 macht zugleich den ersten Re-Evaluierungs-Trigger von
[ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) fällig,
noch bevor die ADR `Accepted` wird. Dieser Report ist ein Lauf-Beleg (diese
Range, dieser Skill, dieses Modell, dieses Verdikt) und ersetzt keine
Verifikation — DoD- und Plan-Konformität prüft der Verifier separat.
