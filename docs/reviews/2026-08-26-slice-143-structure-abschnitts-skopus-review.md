# Review-Report: slice-143 — structure-Abschnitts-Skopus / Ablösung `closure-outcomes.sh` — 2026-08-26

**Review-Art:** Code-Review (Config-/Doku-/Makefile-Diff gegen Kanon, Spec und
Slice-Plan, Modul 10 §Drei Review-Arten) · **Gegenstand:** Commit `f49a6c8`
(Range `HEAD~1..HEAD`) — der Feature-Commit von slice-143 (`feat(harness): Die
Drei-Ausgänge-Regel zieht ins Produkt, das Bash-Skript fällt (DC-FA-STRUCT-001,
slice-143)`). 6 Dateien laut `git show --stat` (`.d-check.closure.yml` +19,
`.d-check.yml` +2/-1, `AGENTS.md` +2/-3, `Makefile` +2/-9, `harness/README.md`
+3/-4, gelöscht `tools/harness/closure-outcomes.sh` -92), 28 Einfügungen /
109 Löschungen gesamt.

**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 · **Modell-ID:**
`claude-opus-5[1m]` · **Datum:** 2026-08-26

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-143-structure-abschnitts-skopus.md`
  (§1 Ziel, §2 Vorgehen, §3 „Ausdrücklich NICHT", §4 DoD, §5 Risiken, §7)
- `AGENTS.md` §3 (Hard Rules, insbesondere §3.6 und §3.7), §4 (Gate-Tabelle
  samt Kopfsatz), §5 (Doku-Regeln, insbesondere der `BEO-009`-Absatz)
- [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
  vollständig (Entscheidung 4 „Struktur, nicht Bedeutung" und Entscheidung 8
  „Fail-closed + geteilte Heading-Lexik"),
  [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  (Entscheidungen 1–4), [ADR-0025](../plan/adr/0025-codepaths-ignore-refs.md)
  (Tombstone-Register, Präzedenz der Eintrags-Form)
- `spec/lastenheft.md` → `DC-FA-STRUCT-001` (Kandidaten-Menge, Abschnitts-
  Bestimmung, Bedingungs-Tabelle mit `forbid-pattern` ⇒ `section-forbidden`),
  `DC-FA-PLAN-001`
- `spec/spezifikation.md` → `§DC-FA-STRUCT-001.a` Schritte 2–6 (insbesondere
  Schritt 5 „Bereinigen") und `§DC-FA-PLAN-001.a` Schritt C4b
- Quellcode: `internal/hexagon/core/rules/structure.go`
  (`CheckStructure`, `structureTree`, `checkStructureRule`,
  `checkStructureFile`, `structureConditions`),
  `internal/hexagon/core/rules/sections.go` (`FindSectionHeads`, `SectionEnd`,
  `SectionProse`, `SectionHeadings`),
  `internal/hexagon/core/rules/markdown.go` (`PreprocessMarkdown`,
  `proseLines`, `stripInlineCodeByLine`)
- Vorherige Findings am gleichen Gegenstand:
  [`2026-08-25-slice-139-closure-ausgang-waechter-review.md`](2026-08-25-slice-139-closure-ausgang-waechter-review.md)
  (F-1/F-2 HIGH, F-3 MEDIUM, F-4 LOW) und
  [`2026-08-25-slice-140-konsumenten-cr-review.md`](2026-08-25-slice-140-konsumenten-cr-review.md)
  (drei MEDIUM, drei LOW)
- `docs/plan/planning/observations.md` → `BEO-007`, `BEO-010`, `BEO-011`,
  `BEO-012`, `BEO-015`
- Vendorte Quellen: `.harness/baseline/v5.11.0/regelwerk/modul-05-planning-harness.md`
  §Offene Risiken werden bei Closure aufgelöst (Volltext),
  `.harness/baseline/v5.11.0/templates/docs/plan/planning/slice.template.md`
  (Platzhalter-Formen)

**Nicht erhalten:** die DoD-Abhakung (Verifikations-Rolle, getrennter Kontext).

**Vom Reviewer selbst gefahren:**

- `git show`/`git log --diff-filter=A` auf den Commit und auf das entfernte
  Skript; `grep`/`find` über `docs/plan/planning/done/**`, `docs/plan/adr/**`,
  die vendorte Baseline und alle lebenden Doku-Flächen.
- `make build` (Exit 0), `make gate-consistency` (Exit 0, 500 Dateien /
  0 Befunde), `make doc-check` (Exit 0, 500 Dateien / 0 Befunde).
- Das Prüf-Kommando des Gates direkt:
  das `DCHECK_RUN`-Kommando des Makefile mit
  `--config .d-check.closure.yml --enable planning --enable structure` — Ausgangslage 459 Dateien / 0 Befunde, Exit 0.
- Das **entfernte** Skript als Vergleichsmaßstab: `git show
  f49a6c8^:tools/harness/closure-outcomes.sh` in den Scratchpad extrahiert und
  von dort gefahren (es `cd`t selbst auf die Repo-Wurzel); Ausgangslage
  „ok (139 Slices in done/, kein offener Platzhalter)", Exit 0. **Keine
  Skript-Datei wurde im Repo angelegt.**
- **Temporäre Repo-Änderungen** (jeweils unmittelbar nach dem Lauf mit
  `git checkout -- <pfad>` zurückgebaut): an
  `docs/plan/planning/done/slice-001-adr-fundament.md` wurden nacheinander
  angehängt bzw. eingefügt — die vier verbotenen Formen einzeln, eine
  groß geschriebene Variante `(Bei Closure)`, eine ASCII-Ellipse `<...>`,
  ein Platzhalter in einem Fenced Block, ein Platzhalter innerhalb einer
  mehrzeiligen Inline-Code-Spanne, ein Platzhalter hinter einer zweiten H1,
  ein generisches `<Risiko>` am Dateiende und dasselbe innerhalb des
  Closure-Abschnitts, sowie ein `<…>` an die H1-Zeile selbst. Nach jedem
  Rückbau wurde der Lauf wiederholt (wieder 0 Befunde); **`git status
  --short` ist am Ende dieses Reviews leer.**

**Verdikt: blockierend** — zwei HIGH, drei MEDIUM, kein LOW.

---

## Findings

### F-1 — Eine mehrzeilige Inline-Code-Spanne verschluckt einen echten Platzhalter: das Gate bleibt grün, wo das abgelöste Skript rot war

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §4 („Halluzinierte Gates sind die häufigste Form von
  Harness-Lüge") · [`BEO-003`](../plan/planning/observations.md) (geteilte
  Lexik, „konfigurationsseitig — ein Muster in der Prüf-Config kann dieselbe
  Lexik anders sprechen als das Modul, das es benutzt") · `DC-FA-STRUCT-001`
- **pfad:** `.d-check.closure.yml:135` (`forbid-pattern`) in Verbindung mit
  `internal/hexagon/core/rules/markdown.go:276`
  (`stripInlineCodeByLine` — Spannen werden **absatzweise** über
  Zeilengrenzen hinweg geleert) vs. dem entfernten
  `tools/harness/closure-outcomes.sh:57` (ein `sed`-Ausdruck, der
  Backtick-Paare **je Zeile** entfernt)
- **befund:** Das entfernte Skript paarte Backticks nur innerhalb einer Zeile;
  die neue Regel liest den absatzweise bereinigten Text, in dem ein einzelner
  Backtick eine Spanne über mehrere Zeilen aufziehen kann. Ein echter, offener
  `*(bei Closure)*`-Ausgang, der zwischen zwei solchen Backticks eines Absatzes
  liegt, wird positionserhaltend geleert und ist für die Regel unsichtbar —
  gemessen: derselbe Bestand mit derselben eingefügten Risikozeile ergibt beim
  Skript Exit 1 (`closure-outcome-open`) und bei der neuen Regel
  „459 Datei(en) geprüft, 0 Befund(e)", Exit 0.
- **verifizierbar:** ja — an eine `done/slice-*.md` drei aufeinander
  folgende Zeilen anhängen — die erste endet mit einem einzelnen Backtick, die
  zweite ist eine Risikozeile mit dem offenen Ausgang in der Form aus der
  Alternation, die dritte trägt den zweiten Backtick — und `make verify-closure-notes` fahren: Exit 0. Dieselbe Einfügung
  gegen `git show f49a6c8^:tools/harness/closure-outcomes.sh`: Exit 1.
- **klasse:** mehrzeilige-inline-code-spanne-verschluckt-platzhalter

### F-2 — `AGENTS.md` §4 und die Sensors-Tabelle sagen „über die **ganze** Slice-Datei" zu; die Regel sieht vier Bereiche der Datei nicht, und der Ehrlichkeits-Vorbehalt der abgelösten Zeile ist ersatzlos entfallen

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §4 Kopfsatz („Nur hier gelistete Targets existieren im
  Makefile. Halluzinierte Gates sind die häufigste Form von Harness-Lüge") ·
  `spec/spezifikation.md` §DC-FA-STRUCT-001.a Schritt 5 (Bereinigung, „ein
  `forbid-pattern`, das auf ein Wort in Backticks zielt, trifft **nicht**") ·
  Kontext-Eskalation des Reviewer-Skills (Gate-Pfad)
- **pfad:** `AGENTS.md:321` („als `forbid-pattern` über die **ganze**
  Slice-Datei") und `harness/README.md:91` (derselbe Wortlaut) gegen
  `internal/hexagon/core/rules/sections.go:65` (`SectionProse`: `ln.No <=
  headingNo` schneidet die Überschriften-Zeile ab, `SectionEnd` endet an der
  nächsten H1, `PreprocessMarkdown` entfernt Fences und leert Inline-Code)
- **befund:** Beide Gate-Beschreibungen behaupten Reichweite über die **ganze**
  Datei. Gemessen sind vier Bereiche ausgenommen: (a) die H1-Zeile selbst
  (`# Slice slice-001: … <…>` ⇒ neue Regel Exit 0, Skript Exit 1), (b) Fenced
  Blocks, (c) Inline-Code-Spannen, (d) alles hinter einer zweiten H1 (F-3).
  Die abgelöste `AGENTS.md`-Zeile trug den Vorbehalt noch selbst („Inline-Code
  wird übersprungen, damit ein Slice über die Platzhalter schreiben kann");
  die neue Zeile nennt ihn nicht mehr und dehnt die Aussage zugleich auf die
  ganze Datei aus. Damit steht in der Tabelle, die laut ihrem eigenen Kopfsatz
  gegen Harness-Lüge steht, eine größere Zusage als der Mechanismus einlöst.
  (Die Ausnahmen (b)/(c) sind spec-gedeckt und lösen F-3 des
  slice-139-Reviews auf — der Befund gilt der **Beschreibung**, nicht der
  Bereinigung.)
- **verifizierbar:** ja — `<…>` an die H1-Zeile einer `done/slice-*.md`
  anhängen bzw. einen Fenced Block mit `<eingetreten: …>` einfügen und
  `make verify-closure-notes` fahren: jeweils Exit 0, während das Skript
  Exit 1 lieferte.
- **klasse:** gate-zusage-reicht-weiter-als-der-mechanismus

### F-3 — Eine zweite H1 kappt die Spanne still; die Kardinalitäts-Wache greift nicht, weil sie nur Treffer des Selektors zählt

- **kategorie:** MEDIUM
- **quelle:** `DC-FA-STRUCT-001` §Abschnitte bestimmen („ein Abschnitt reicht
  bis zur nächsten Überschrift gleicher oder höherer Ebene") ·
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
  Entscheidung 8 (b) (dieselbe Klasse: eine als H1 gelesene Zeile „hätte den
  Abschnitt dort abgeschnitten und alles dahinter unsichtbar gemacht") ·
  Reviewer-Skill LOW „latente Wartungsfalle", eine Stufe hoch wegen Gate-Pfad
- **pfad:** `.d-check.closure.yml:133-135` (kein `sections`-Schlüssel, Default
  `one`) gegen `internal/hexagon/core/rules/sections.go:40` (`SectionEnd`
  endet an **jeder** H1) und `internal/hexagon/core/rules/structure.go:104`
  (`section-ambiguous` nur bei mehr als einem **Selektor**-Treffer)
- **befund:** Trägt eine `done/`-Slice-Datei eine zweite H1, die **nicht** auf
  `^# Slice slice-` passt (etwa `# Anhang`), endet die geprüfte Spanne dort;
  alles dahinter ist ungeprüft und meldet weder `section-ambiguous` noch
  `section-missing`. Gemessen: `# Anhang` plus eine Zeile
  `- Risiko R — **Ausgang:** *(bei Closure)*` ergibt bei der neuen Regel Exit 0,
  beim entfernten Skript Exit 1. Die Commit-Botschaft führt als tragende
  Eigenschaft „ein Selektor, der nicht still leerlaufen kann — greift er nicht,
  meldet die Regel `section-missing`" — das deckt den Null-Treffer-Fall, nicht
  die Kappung. Der heutige Bestand ist unauffällig (139/139 Dateien: genau
  **eine** H1, auf Zeile 1), die Falle zündet beim nächsten Slice, der eine
  zweite H1 setzt.
- **verifizierbar:** ja — `# Anhang` und danach eine der vier verbotenen Formen
  an eine `done/slice-*.md` anhängen; `make verify-closure-notes` Exit 0.
- **klasse:** abschnitts-spanne-still-gekappt-durch-fremde-h1

### F-4 — Das Tombstone-Register schreibt die Skript-Entfernung zwei ADRs zu, die sie nicht entschieden haben — und die es zeitlich nicht konnten

- **kategorie:** MEDIUM
- **quelle:** [`BEO-012`](../plan/planning/observations.md) (Zähler 4, „eine
  ADR-Entscheidung wird als stehendes Verbot gelesen, obwohl sie einen
  einmaligen Akt beschreibt") · `AGENTS.md` §3.6 („Jede Schwellen-Senkung
  (Coverage, Linter-Strenge, **Prüfregel**) ist ein ADR, kein PR-Kommentar") ·
  Präzedenz der Eintrags-Form in
  [ADR-0025](../plan/adr/0025-codepaths-ignore-refs.md)
- **pfad:** `.d-check.yml:267` (`# - harness/closure-outcomes.sh: abgelöst
  durch structure.forbid-pattern im Closure-Profil (ADR-0048/ADR-0049,
  slice-143)`)
- **befund:** Die fünf Schwester-Einträge des Registers nennen jeweils die ADR,
  die die Ablösung **entschieden** hat (ADR-0025, ADR-0026, ADR-0027, ADR-0028,
  ADR-0029, ADR-0032). ADR-0048 und ADR-0049 tragen beide `**Datum:**
  2026-08-09`; das Skript entstand am 2026-08-25 (`612a619`,
  `git log --diff-filter=A`), und `grep -rn "closure-outcomes" docs/plan/adr/`
  liefert keinen Treffer. Beide ADRs entscheiden die Closure-Struktur bzw. den
  Modul-Schnitt, keine von ihnen die Entfernung dieses Gates. Wer dem Zeiger
  folgt, um zu erfahren, warum ein Closure-Gate verschwunden ist, findet dort
  keine Entscheidung — die Entfernung steht damit ohne den Beleg da, den §3.6
  für eine gelockerte Prüfregel verlangt (F-1/F-2/F-3 zeigen, dass die Ablösung
  in vier messbaren Punkten weniger deckt als das Entfernte).
- **verifizierbar:** ja — `grep -rn "closure-outcomes" docs/plan/adr/` (leer)
  und `grep -m1 "^\*\*Datum" docs/plan/adr/0048-*.md docs/plan/adr/0049-*.md`
  gegen `git log --diff-filter=A -- tools/harness/closure-outcomes.sh`.
- **klasse:** tombstone-zuschreibung-an-nicht-entscheidende-adr

### F-5 — Der neue Register-Kommentar trägt eine Slice-Nummer und drei Herkunfts-Token, wo §3.7 genau ein Feld zulässt

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §3.7 („keine Slice-Nummern und keine Mess-Labels;
  Herkunft nur als **ein** auflösbares Feld nach dem Baseline-Schema (`DC-*`
  …, `ADR-*`, `MR-*`, `seit welle-<NN>`)"; „Neuzugänge fallen überall unter
  den Anker") ·
  [Baseline §Was ein Kommentar trägt](../../.harness/baseline/v5.11.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte)
- **pfad:** `.d-check.yml:267`
- **befund:** Der Kommentar schließt mit `(ADR-0048/ADR-0049, slice-143)` — drei
  Herkunfts-Token statt eines, darunter eine Slice-Nummer, die §3.7 in Code-,
  Konfigurations- und Skript-Kommentaren ausdrücklich ausschließt. Die fünf
  darüberstehenden Zeilen tragen dieselbe Form, sind aber nach der
  Bestandsgrenze grandfathered; diese Zeile ist ein Neuzugang. (Die
  Kommentar-**Klasse** selbst ist getroffen — Abgrenzung/Grenze —, deshalb
  MEDIUM statt des HIGH-Ankers, der dem klassenlosen Kommentar gilt.) Kein Gate
  prüft das: `make doc-check` und `make gate-consistency` laufen beide grün.
- **verifizierbar:** nein (Urteil, kein `grep` — die Prüfung ist per §3.7
  ausdrücklich Reviewer-Sache); der beobachtbare Zustand ist die Zeile selbst.
- **klasse:** herkunftsfeld-ueberfuellt-mit-slice-nummer

---

## Negativbefunde

- geprüft, ohne Befund: **Form-für-Form-Parität der vier Muster** — jede der
  vier Alternativen wurde einzeln in eine `done/`-Slice eingefügt und gefahren;
  alle vier melden `section-forbidden`, Exit 1 (`(bei Closure)`, `wird mit dem
  Closure-Body gefüllt`, `<…>`, `<eingetreten:`), Rückbau je Exit 0. Auf dem
  reinen Prosa-Pfad trifft die Alternation dieselbe Menge wie die
  `PATTERNS`-Liste des Skripts.
- geprüft, ohne Befund: **YAML-Quoting und RE2-Escaping** — die einfach
  gequotete Skalare gibt `\(bei Closure\)` unverändert an RE2 weiter; die
  Klammern treffen literal (gemessen), die Ellipse `…` ist als Mehrbyte-Literal
  wirksam (gemessen), `|` trennt korrekt vier Alternativen, keine
  Kompilier-Fehler (kein Exit 2 in allen Läufen).
- geprüft, ohne Befund: **Groß-/Kleinschreibung und ASCII-Ellipse** —
  `(Bei Closure)` und `<...>` bleiben grün. Das ist **Parität** zum entfernten
  Skript (`case`-Vergleich, ebenfalls literal) und von der Config-Zeile
  `.d-check.closure.yml:130-132` als „die Alternation ist eine LISTE"
  ausdrücklich benannt — keine Verschlechterung, keine unbenannte Grenze.
- geprüft, ohne Befund: **Selektor-Deckung des Bestands** — alle 139
  `done/slice-*.md` tragen genau **eine** H1, und zwar auf Zeile 1, und alle
  passen auf `^# Slice slice-`. Es gibt heute keine Datei mit null Treffern
  (`section-missing`) und keine mit zwei Selektor-Treffern
  (`section-ambiguous`); die H1-Spanne reicht damit im heutigen Bestand
  tatsächlich vom Kopfteil bis zum Dateiende (Einschränkungen: F-2/F-3).
- geprüft, ohne Befund: **Nullmengen-Härte** — `checkStructureRule`
  (`structure.go:68-76`) meldet `section-missing` auf dem Glob, wenn die Regel
  keine Datei trifft; `structureTree` läuft über den ganzen Baum unabhängig von
  `scan.roots`/`scan.ignore`. Ein Umzug des `done/`-Bestands ließe die Regel
  also nicht still leerlaufen — die Eigenschaft, die die Commit-Botschaft für
  sich reklamiert, ist an dieser Stelle korrekt.
- geprüft, ohne Befund: **Inline-Code-Negativkontrolle** — dieselbe Form in
  einfachen Backticks bleibt grün; ohne sie meldete jeder Slice, der über die
  Platzhalter schreibt, seine eigene Dokumentation. Spec-gedeckt
  (§DC-FA-STRUCT-001.a Schritt 5, letzter Satz).
- geprüft, ohne Befund: **Fähigkeit ist zugesagt, „kein Produkt-Delta" trägt** —
  `DC-FA-STRUCT-001` führt `forbid-pattern` (RE2) ⇒ `section-forbidden` in
  seiner Bedingungs-Tabelle, ebenso `§DC-FA-STRUCT-001.a` Schritt 6; `sections`
  und `exempt-paths` stehen dort ebenfalls. Für die **benutzte Fähigkeit**
  braucht es weder Lastenheft-Bump noch neue ADR; die ADR-Frage stellt sich am
  **Entfernen** des Gates (F-4), nicht am Hinzufügen der Regel.
- geprüft, ohne Befund: **ADR-0048 Entscheidung 4 (Struktur statt Bedeutung)** —
  die Regel prüft die Abwesenheit vier benannter Zeichenketten, nicht ob ein
  eingetragener Ausgang trägt; beide Doku-Flächen sagen das ausdrücklich dazu.
  Die Vertragsgrenze bleibt gewahrt.
- geprüft, ohne Befund: **ADR-0048 Entscheidung 8 (b) geteilte Heading-Lexik** —
  die Regel benutzt `FindSectionHeads`/`parseATXHeading`, denselben Parser wie
  `anchors`/`matrix`/die Closure-Fähigkeit; eine eigene `#`-Zählung existiert
  nicht (`sections.go:19-36`).
- geprüft, ohne Befund: **Spiegel des entfernten Targets** — `Makefile`
  (Target-Rumpf, `.PHONY:47`, `fullbuild:209`), `AGENTS.md:321-322`,
  `harness/README.md:91`, `:92` (fullbuild-Kette) und `:110`
  (Meta-Gates-Klassifikation) sind alle nachgezogen; `grep -rn
  "closure-outcomes"` findet außerhalb von `.d-check.yml` nur eingefrorene
  Fundstellen (`done/slice-139`, zwei Reviews) und den lebenden slice-143-Plan,
  der den Pfad als Gegenstand nennt. Die Prosa-Kette „gates + image-test +
  bench + completeness-check + verify-closure-notes" steht an allen drei
  Stellen identisch — die `BEO-010`/F-4-Klasse des Vorgänger-Reviews ist diesmal
  vollständig bedient.
- geprüft, ohne Befund: **Tombstone wirkt und wird gebraucht** — der lebende
  slice-143-Plan führt `tools/harness/closure-outcomes.sh` in Inline-Code
  (`Bezug:`-Feld); ohne `codepaths.ignore-refs` wäre das `codepath-missing`.
  `make doc-check` läuft grün (500 Dateien, 0 Befunde), der Eintrag greift also
  (Mangel liegt in seiner **Zuschreibung**, F-4, nicht in seiner Wirkung).
- geprüft, ohne Befund: **`make gate-consistency`** — Exit 0 (500 Dateien,
  0 Befunde). Es deckt diesen Fall aber auch nur zur Hälfte ab: das Modul
  `targets` prüft Deklarations-Konsistenz Doku↔Makefile über Target-**Namen**;
  ein Target, das nirgends mehr genannt wird, ist für es unauffällig. Die
  Vollständigkeit der Spiegel oben ist von Hand belegt, nicht vom Gate.
- geprüft, ohne Befund: **Bindepunkt-Einordnung** — die Regel wohnt im
  `--config`-Profil und läuft damit an `fullbuild`, nicht an `gates`/`ci`;
  dieselbe Trennung, die ADR-0048 Entscheidung 7 und ADR-0026 tragen. Kein
  Closure-Gate ist in den Inner-Loop gerutscht, und `modules: []` ist unberührt
  (keine zweite Netz-Tür, `DC-QA-03`).
- geprüft, ohne Befund: **Kanon-Zitate (BEO-012-Gegenprobe)** — `#### Offene
  Risiken werden bei Closure aufgelöst` existiert wortgleich in
  `.harness/baseline/v5.11.0/regelwerk/modul-05-planning-harness.md:126`; der in
  `AGENTS.md:321` zitierte Satz *„Ein Slice geht nicht nach `done/`, während ein
  Risiko ohne Ausgang dasteht"* steht dort byte-gleich als Schluss-Satz des
  Absatzes und im richtigen Geltungsbereich (Slice-Closure, nicht
  Wellen-Closure). Der Name **Drei-Ausgänge-Regel** steht so **nicht** im
  vendorten Baum (`grep` leer) — er ist ein Repo-Label für eine Kanon-Regel,
  deren Inhalt („genau **einen** Ausgang: *eingetreten* · *entfallen* · *weiter
  offen*") dort steht; die Zuschreibung dehnt die Quelle nicht, sie benennt sie
  nur mit einem eigenen Wort, und sie ist gegenüber slice-139 unverändert.
- geprüft, ohne Befund: **`.d-check.closure.yml`-Kommentar gegen §3.7** — der
  neue Block (`:118-132`) trägt Zusage, Abgrenzung und eine ausdrücklich
  benannte Grenze, keine Slice-Nummer, kein Mess-Label und keine
  Review-Historie. Er ist die formkonforme Gegenprobe zu F-5.
- geprüft, ohne Befund: **Modul-Grenze auf der Ziel-Achse (§3.8/`BEO-004`)** —
  die Regel liest ausschließlich Dateien, die ihr eigener `files`-Glob
  bestimmt; sie folgt keinem Link, keinem Anker und keiner git-Revision. Es
  gibt keine zweite, ungeprüfte Eingabe-Achse.

## Kein Finding, aber gemessen und benannt

Die generischen Vorlagen-Platzhalter der Kanon-Vorlage (`<Titel>`, `<Name>`,
`<Risiko>`, `<Bedingung>`, `<welle-id>`, `<slice-NNN>` … — 21 whitespace-freie
Formen) sieht die neue Regel nicht (sie nennt zwei davon) und `planning` nur
**innerhalb** des Closure-Abschnitts. Gemessen: `<Risiko>` am Ende von
`slice-001` (dessen Closure-Notiz §7 ist, gefolgt von §8) ⇒ **beide** stumm,
Exit 0; dasselbe `<Risiko>` innerhalb des Closure-Abschnitts ⇒
`closure-note-placeholder`, Exit 1. In **31 von 139** `done/`-Slices steht die
Closure-Notiz nicht am Dateiende, dort existiert diese Zone real. Das ist die
Restmenge von F-1 des slice-139-Reviews und **keine Verschlechterung** durch
diesen Commit — das entfernte Skript deckte dieselben Formen ebenso wenig ab.
Notiert, weil die Gegenprobe der Commit-Botschaft („Gegenprobe mit einem
generischen `<Name>`: nur `planning` meldet") an einer Datei lief, deren
Closure-Notiz der **letzte** Abschnitt ist, und der daraus gezogene
Komplement-Schluss für die 31 anderen nicht gilt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 3 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:**
mehrzeilige-inline-code-spanne-verschluckt-platzhalter ·
gate-zusage-reicht-weiter-als-der-mechanismus ·
abschnitts-spanne-still-gekappt-durch-fremde-h1 ·
tombstone-zuschreibung-an-nicht-entscheidende-adr ·
herkunftsfeld-ueberfuellt-mit-slice-nummer

## Verdikt

**Merge-blockierend:** ja — zwei HIGH. Die Wegwahl selbst ist tragfähig: die
Fähigkeit ist zugesagt, der Selektor deckt den heutigen Bestand vollständig, die
vier Formen treffen, die Nullmengen-Härte ist echt, alle Spiegel sind gezogen,
und der Fence-Schutz löst F-3 des slice-139-Reviews auf. Blockierend ist die
**Reichweite**: F-1 ist ein stiller Grün-Pfad, den das abgelöste Skript rot
meldete (gemessen, beide Richtungen), und F-2 ist die Differenz zwischen dem,
was `AGENTS.md` §4 zusagt, und dem, was läuft — in genau der Tabelle, deren
Kopfsatz Halluzination als häufigste Harness-Lüge benennt, und unter Wegfall
des Vorbehalts, den die abgelöste Zeile noch trug. F-3 ist dieselbe Klasse
latent. F-4 und F-5 betreffen dieselbe neue Zeile im Tombstone-Register: eine
Zuschreibung an ADRs, die den Akt nicht entschieden haben (`BEO-012`, Zähler
steht bei 4), und ein überfülltes Herkunftsfeld.

Die Aussage der Commit-Botschaft „die Zusage wird nicht schwächer" hält damit in
drei von vier gemessenen Achsen **nicht**: H1-Zeile, mehrzeilige Inline-Spanne
und Zweit-H1 sind neu stumm; die vierte (Fenced Blocks) ist eine gewollte,
spec-gedeckte Verengung, die nur in der Doku fehlt.

**Übergabe:** Findings gehen an den Implementer. F-1 ist ein Kandidat für eine
weitere Instanz der `BEO-003`-Klasse (Config-Muster spricht eine andere Lexik
als das Modul, das es ausführt), F-4 für die fünfte Instanz von `BEO-012`; die
Einordnung obliegt dem Maintainer bei der Slice-Closure, nicht diesem Report.
Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation (DoD-/
Spec-Konformität prüft der Verifier separat).
