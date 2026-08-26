# Review-Report: slice-145 — Ein Sensor auf Pfad-Token in der Architektur-Sicht — 2026-08-26

**Review-Art:** Code-/Konfigurations-Review (Config- und Hard-Rule-Diff gegen
Slice-Plan, Kanon, Spezifikation und Modul-Code) · unabhängiger Reviewer ohne
Anteil an der Arbeit
**Gegenstand:** Commit-Kette `147476a..11054fa` von slice-145 — `147476a`
(Anspruchs-/Lifecycle-Move, nur `git mv` + Roadmap-Zeile), `701a567` (nur die
Schwelle, +17 Zeilen im Slice-Plan), `11054fa` (Feat: `.d-check.yml` +19,
`AGENTS.md` +17/−6)
**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 (`9ee805b`) · **Modell-ID:**
`claude-opus-5[1m]` · **Datum:** 2026-08-26

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-145-pfad-token-sicht.md`
  vollständig — §1 Ziel, §2 Schritt 0 (die vorab notierte Schwelle, vier
  Bedingungen), §3 „Ausdrücklich NICHT", §4 DoD, §5 Risiken, §7 Vorgelagert
- `.d-check.yml`, Block `structure` vollständig (der neue Eintrag samt
  Kommentar, Zeilen 339–356, und die sieben Bestandsregeln), dazu `matrix`
  (alle fünf Klassen, `rules`, `exempt-paths`), `codepaths`, `ids`, `scan`
- `AGENTS.md` §3.4 vorher (`git show 11054fa^:AGENTS.md`) und nachher, dazu
  §3.2, §3.6, §3.7, §3.8, §4 (Gate-Tabelle samt Kopfsatz), §5
- `spec/architecture.md` vollständig (158 Zeilen)
- `spec/lastenheft.md` → `DC-FA-STRUCT-001` (Kandidaten-Menge,
  Abschnitts-Bestimmung, `sections`-Kardinalität, Bedingungs-Tabelle mit
  `forbid-pattern` ⇒ `section-forbidden`), `DC-FA-MTX-001`, `DC-FA-MTX-003`,
  `DC-FA-CLI-003` (Exit-Codes), `DC-FA-CODE-001`;
  `spec/spezifikation.md` → `§DC-FA-STRUCT-001.a`
- Quellcode: `internal/hexagon/core/rules/structure.go`
  (`CheckStructure`, `checkStructureRule`, `checkStructureFile`,
  `structureConditions`, `structureExempt`),
  `internal/hexagon/core/rules/sections.go` (`FindSectionHeads`, `SectionEnd`,
  `SectionProse`), `internal/hexagon/core/rules/markdown.go`
  (`PreprocessMarkdown`, `proseLines`, `stripInlineCodeByLine`,
  `forEachInlineCodeSpan`), `internal/hexagon/core/rules/matrix.go`
  (`CheckMatrix`, `classOf`, `orderedClass`, `downwardFinding`,
  `tokenFindings`)
- `harness/conventions/MR-033-sicht-ohne-modul-pfade.md` vollständig
  (insbesondere die Felder `Geltungsbereich`, `Adaption`,
  `Ersetzt-Baseline-Regel`, `Auflösungs-Trigger`)
- Vendorte Baseline `v5.12.0`: `regelwerk/modul-03-spec.md`
  §Ziel-Form: Architektur-Sicht (Volltext),
  `regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt
- `docs/plan/planning/observations.md` → `BEO-009` (Zähler 6), `BEO-011`,
  `BEO-012`, `BEO-016`
- Vorherige Findings an derselben Mechanik:
  `docs/reviews/2026-08-26-slice-144-commit-hash-muster-review.md`
  (F-1 HIGH `konfig-kommentar-traegt-mess-historie-und-slice-nummer`, F-2/F-3/
  F-4/F-5 MEDIUM) und
  `docs/reviews/2026-08-26-slice-143-structure-abschnitts-skopus-review.md`
  (F-1/F-2 HIGH, F-3 `abschnitts-spanne-still-gekappt-durch-fremde-h1`)

**Nicht erhalten:** die DoD-Abhakung (Verifikations-Rolle, getrennter Kontext).
Die Haken in §4 des Slice-Plans sind in diesem Report weder gesetzt noch
bewertet.

**Vom Reviewer selbst gefahren:**

- `git show`/`git log --format=%ad` auf alle drei Commits; `git show
  147476a:<slice-datei>` für den Plan-Stand **vor** dem Schwellen-Commit.
- `make doc-check` (Exit 0, 517 Dateien / 0 Befunde) — das Image ist auf dem
  Stand von `HEAD` (`docker build` reexportiert `sha256:f41f99ad9a1e`).
- **Reproduktion der Vier-Kandidaten-Messung:** Sonden-Repo im Scratchpad mit
  einer Kopie von `spec/architecture.md` und einer Sonden-`.d-check.yml`
  (`modules: [structure]`, dieselbe Regel, nur `forbid-pattern` variiert), je
  Kandidat ein Lauf gegen das Runtime-Image — 0 / 0 / **1** / 0, Exit
  0 / 0 / **1** / 0.
- **Trefferzählung** für den durchgefallenen Kandidaten: eine Python-Nachbildung
  von `proseLines` → `proseParagraphs` → `stripInlineCodeByLine` →
  `SectionProse` (dieselben Regeln, inklusive der CommonMark-Infozeilen-Regel
  in `FenceToggle`), gegengeprüft am Produkt über einzelne Sonden-Muster.
- **Falsch-Negativ-Sonde:** eine `spec/architecture.md` aus 23 H1-Abschnitten,
  je Abschnitt genau eine Schreibform; ein Lauf, 12 Befunde — die Zuordnung
  Form → Befund ergibt sich aus der Befund-Zeile (= H1-Zeile des Abschnitts).
- **Mehrzeilen-Spannen-Sonde** (`BEO-016`): zwei identische Abschnitte, einer
  mit einer absatzweiten Inline-Code-Spanne über dem Pfad, `modules:
  [structure, spans]`.
- **Abschnitts-Selektor-Sonden:** H1 verloren · zweite H1 mit Pfad dahinter ·
  Datei fehlt ganz · H1 nur innerhalb einer Fence · `exempt-paths` auf die
  einzige Kandidaten-Datei · `scan.ignore` auf `spec/**`.
- **matrix-Alternativ-Sonden:** drei Läufe über ein Sonden-Repo mit den drei
  Straten-Dateien und einem Abwärts-Link Lastenheft → Sicht — (A) Token-Klasse
  ohne `paths` an `spec-straten`, (B) eigene Klasse `sicht` **vor**
  `spec-straten`, (C) dieselbe **nach** `spec-straten`.
- **Bestandszählung** `\b(internal|cmd)/[a-z]` über alle Markdown-Dateien des
  Repos (ohne `.git`, `.harness`): roh **901**, außerhalb Fences 899, im
  bereinigten Text **30**.
- **Temporäre Repo-Änderungen** (jeweils unmittelbar nach dem Lauf mit
  `git checkout -- spec/architecture.md` zurückgebaut, Byte-Gleichheit gegen
  eine vorher gezogene Kopie geprüft): an `spec/architecture.md` wurde
  nacheinander angehängt — (1) `Der Kern liegt unter internal/hexagon/core/rules.`
  (Prosa) für die Positivkontrolle, (2) dieselbe Zeile mit dem Pfad in
  Inline-Code. Nach jedem Rückbau wurde der Lauf wiederholt (wieder 0 Befunde);
  **`git status --short` ist am Ende dieses Reviews leer.**

**Verdikt: blockierend** — ein HIGH, fünf MEDIUM, zwei LOW, ein INFO.

---

## Findings

### F-1

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §3.4 Kopfsatz (*„benennt Schichten und Rollen statt
  Technologie — keine Sprach-/Modul-Pfade"*, unverändert) · `AGENTS.md` §3.6
  (*„Jede Schwellen-Senkung (Coverage, Linter-Strenge, **Prüfregel**) ist ein
  ADR, kein PR-Kommentar"*) · `harness/conventions/MR-033-sicht-ohne-modul-pfade.md`
  §Adaption (*„Nennt die Sicht Pfade, wandert sie bei jedem Verzeichnis-Umbau
  mit — und das ist genau die Kopplung, die das Stratum vermeiden soll"*) ·
  `BEO-012` (eine Regel wird über ihren Geltungsbereich hinaus gedehnt) ·
  Kontext-Eskalation des Reviewer-Skills (Gate-Pfad)
- **pfad:** `AGENTS.md:201-202` (*„Inline-Code und Fenced Blöcke zählen nicht
  (die Sicht darf einen Pfad *zitieren*, ohne ihn zu *führen*)"*) und
  `.d-check.yml:349-350` (derselbe Satz) gegen
  `internal/hexagon/core/rules/sections.go:65` (`SectionProse` →
  `PreprocessMarkdown`) und `internal/hexagon/core/rules/markdown.go:276`
  (`stripInlineCodeByLine`)
- **befund:** Die lexikalische Blindstelle des Sensors wird in derselben
  Änderung zur **Erlaubnis** erklärt: ein Modul-Pfad in Inline-Code sei ein
  *Zitat* und kein *Führen*. Diese Unterscheidung kennt weder der Kopfsatz von
  §3.4 (unverändert *„keine Sprach-/Modul-Pfade"*) noch `MR-033`
  (`grep -n "zitier"` über den Eintrag und über die vendorte
  `modul-03-spec.md` §Ziel-Form: Architektur-Sicht: kein Treffer), und die in
  `MR-033` genannte **Begründung der Strenge** — der Pfad wandert bei jedem
  Verzeichnis-Umbau mit — trifft einen Pfad in Backticks genauso wie einen in
  Prosa; der Kanon schreibt sein eigenes Beispiel eines erlaubten Modul-Pfads
  (`src/service/`) selbst in Inline-Code. Gemessen ist das die
  **Haus-Schreibweise**: von 901 rohen `internal/`-/`cmd/`-Pfad-Token in den
  Markdown-Dateien dieses Repos stehen 869 in Inline-Code und nur 30 im Text,
  den die Regel liest. Am Bestand vorgeführt: die Zeile ``Der Kern liegt in
  `internal/hexagon/core/rules`.`` in `spec/architecture.md` passiert den
  vollen `doc-check` mit **Exit 0, 517 Dateien, 0 Befunde** — `structure` sieht
  sie nicht, und `codepaths` bestätigt sogar die Existenz des Pfads. Damit
  senkt der Commit die geschriebene Reichweite einer Hard Rule auf die
  Lexik-Grenze ihres neuen Sensors, ohne `MR-033` anzufassen: Regeltext und
  Abweichungs-Register sagen seither Verschiedenes.
- **verifizierbar:** ja — die Inline-Code-Zeile an `spec/architecture.md`
  anhängen und `make doc-check` fahren (Exit 0); dieselbe Zeile ohne Backticks
  ergibt `section-forbidden`, Exit 1. Die Zählung 869/901 ist über
  `PreprocessMarkdown`-Nachbildung und Rohtext reproduzierbar.
- **klasse:** `sensor-luecke-als-erlaubnis-in-die-regel-geschrieben`

### F-2

- **kategorie:** MEDIUM
- **quelle:** `BEO-009` Richtung (b) (*„die genannten Proben liefen und
  stimmen, aber der Schluss daraus gilt weiter als sie reichen"*) ·
  Reviewer-Skill §MEDIUM-Anker *Botschaft verallgemeinert über die Messung
  hinaus* · `spec/lastenheft.md` §`DC-FA-STRUCT-001` Bedingungs-Tabelle
  (`forbid-pattern` ⇒ **ein** `section-forbidden` je Abschnitt) ·
  Vorbefund `zahl-ohne-messvorschrift` (slice-144-Report F-2)
- **pfad:** Commit-Botschaft `11054fa`, Absatz *„Der dritte fängt einen
  DOKUMENT-Pfad in der Sicht und fällt damit durch die Schwelle — die Regel
  gilt Modul-Pfaden, und ein Dokument ist kein Modul."* gegen
  `internal/hexagon/core/rules/structure.go:175`
  (`regexp.MustCompile(r.ForbidPattern).MatchString(body)` — ein boolescher
  Test, ein Befund)
- **befund:** „1 Befund" ist bei dieser Bedingung eine **Abschnitts**-Aussage,
  keine Trefferzahl: `structureConditions` prüft `MatchString` und meldet
  höchstens einmal je Abschnitt; `spec/architecture.md` hat genau eine H1, also
  ist 1 die Obergrenze. Auf dem bereinigten Text trifft
  `\b[a-z][a-z0-9-]*/[a-z]` tatsächlich **fünfmal an drei Stellen**: einmal in
  `spec/architecture.md:7` (das Markdown-Link-Ziel
  `../.harness/baseline/…/regelwerk/modul-03-spec.md`, zwei Treffer), zweimal
  auf Zeile 86 in `stdout/stderr` und einmal auf Zeile 154 in
  `fehlt/unlesbar`. Zwei der drei Stellen sind **überhaupt keine Pfade**,
  sondern Schrägstrich-Wortpaare — der Kandidat scheitert also nicht nur an der
  Dokument-Pfad-Grenze, sondern an gewöhnlicher Prosa, und die Botschaft nennt
  die Ursache, die sie nicht gemessen hat.
- **verifizierbar:** ja — Sonden-Muster einzeln gegen das Produkt fahren
  (`stdout/s` und `fehlt/u` liefern je einen `section-forbidden`), und die
  Trefferzählung über die `SectionProse`-Nachbildung.
- **klasse:** `befundzahl-als-trefferzahl-gelesen`

### F-3

- **kategorie:** MEDIUM
- **quelle:** slice-145 §2 Schritt 0, vierte Bedingung (*„Keine dritte Mechanik
  … die Wahl … wird **begründet**, nicht bequem getroffen"*) ·
  `AGENTS.md` §3.4 (*„Die Commit-Hash-Klasse ist eine Token-Klasse **ohne
  Zieldateien**: ihr Gegenstand ist eine Zeichenkette, kein Dokument"*) ·
  `.d-check.yml:239-243` (Klasse `commit-hash`, ohne `paths`) ·
  `BEO-011` Ausprägung (b)
- **pfad:** Commit-Botschaft `11054fa`, Absatz *„matrix regelt Referenzen
  ZWISCHEN Dokumentklassen, und ein Pfad-Token in EINER Datei ist keine solche
  Referenz"* gegen `internal/hexagon/core/rules/matrix.go:219` (`classOf` —
  First-Match) und `:92` (`tokenFindings`)
- **befund:** Der tragende Satz der Wegwahl ist durch die Klasse widerlegt, die
  dieses Repo im Commit davor gebaut hat: `commit-hash` ist genau „ein Token in
  **einer** Datei", und die Botschaft nennt sie selbst als verworfene
  Alternative. Gemessen: eine Token-Klasse `modul-pfad` ohne `paths` plus
  `{from: spec-straten, to: modul-pfad, allow: false}` meldet den Pfad in der
  Sicht als `matrix-forbidden` — der Weg war mechanisch offen. Was ihn
  tatsächlich versperrt, steht in der Botschaft nur als Nebensatz (*„stört die
  Klassen-Zugehörigkeit von matrix nicht"*) und ist ungemessen: `classOf`
  liefert die **erste** passende Klasse, also kippt eine eigene Klasse `sicht`
  mit `paths: [spec/architecture.md]` in beide Richtungen still — steht sie
  **vor** `spec-straten`, verliert `spec/architecture.md` seine
  Straten-Mitgliedschaft und der `matrix-downward`-Befund für den
  Abwärts-Link Lastenheft → Sicht verschwindet (Sonde: 2 Befunde → 1); steht
  sie **danach**, feuert die Regel `sicht → modul-pfad` nie (Sonde: der
  Pfad-Befund fehlt, Exit 1 nur wegen `matrix-downward`). Wer diese Begründung
  später wiederverwendet, schließt `matrix` für Token-in-einer-Datei aus einem
  Grund aus, den die Nachbarklasse widerlegt.
- **verifizierbar:** ja — die drei Sonden-Konfigurationen (A)/(B)/(C) über ein
  Repo mit den drei Straten-Dateien und einem Abwärts-Link; Befundsätze wie
  oben.
- **klasse:** `wegwahl-begruendet-mit-einer-widerlegten-modul-eigenschaft`

### F-4

- **kategorie:** MEDIUM
- **quelle:** slice-145 §2 Schritt 0, dritte Bedingung (*„Die
  Falsch-Negativ-Klasse ist zu **benennen**, nicht zu minimieren"*) ·
  `internal/hexagon/core/rules/sections.go:69` (`ln.No <= headingNo` — die
  Überschriften-Zeile fällt aus dem Text) · Vorbefund
  `falsch-negativ-liste-nennt-den-groessten-fall-als-erfasst`
  (slice-144-Report F-4)
- **pfad:** `.d-check.yml:349-352` und `AGENTS.md:200-203` (jeweils die
  abschließend gemeinte Aufzählung: Inline-Code, Fenced Blöcke,
  Großschreibung, `tools/`)
- **befund:** Drei weitere Formen fallen gemessen durch und stehen in keiner
  der drei Fassungen: (a) ein Pfad, der über einen **Zeilenumbruch** geht —
  `internal/` am Zeilenende, `hexagon` in der nächsten Zeile — bleibt stumm,
  weil `[a-z]` kein `\n` matcht; (b) die **Überschriften-Zeile des Abschnitts
  selbst** ist nicht Teil des gemessenen Textes (`SectionProse` beginnt hinter
  ihr), heute also `spec/architecture.md:1`; (c) jedes Zeichen außer einem
  ASCII-Kleinbuchstaben **direkt hinter dem Schrägstrich** —
  `internal/_helper` und `internal/2fa` sind stumm, was „Großschreibung fällt
  heraus" nicht abdeckt. Alle drei sind an der 23-Abschnitte-Sonde einzeln
  gemessen (Abschnitte S07, S17, S21, S22 ohne Befund, während die
  Kontroll-Abschnitte melden). Die Zusage der Schwelle war, die Klasse zu
  benennen; benannt sind vier von sieben.
- **verifizierbar:** ja — die Sonden-Datei mit je einem Abschnitt pro Form; ein
  Lauf liefert 12 Befunde, und die vier genannten Abschnitte sind nicht
  darunter.
- **klasse:** `falsch-negativ-liste-unvollstaendig`

### F-5

- **kategorie:** MEDIUM
- **quelle:** `BEO-016` §Prozedur bis dahin (*„wer einen Wächter auf
  bereinigten Text stellt, prüft eine Probe **innerhalb** einer mehrzeiligen
  Spanne mit"*) · slice-145 §7 Vorgelagert (sichtet als offene Beobachtung nur
  `BEO-012`) · Vorbefund
  `mehrzeilige-inline-code-spanne-verschluckt-platzhalter`
  (slice-143-Report F-1, HIGH, dieselbe Bedingung desselben Moduls)
- **pfad:** `.d-check.yml:353-356` (die neue Regel) in Verbindung mit
  `internal/hexagon/core/rules/markdown.go:276` (`stripInlineCodeByLine` —
  Spannen werden **absatzweise** über Zeilengrenzen hinweg geleert);
  Commit-Botschaft `11054fa`, Absatz *„Alle fünf Formen einzeln gefahren"*
- **befund:** Die vom Beobachtungs-Register für genau diesen Akt
  vorgeschriebene Probe fehlt unter den fünf gefahrenen Formen, und sie fällt
  durch: ein Modul-Pfad, der zwischen zwei Backticks **verschiedener Zeilen
  desselben Absatzes** liegt, wird positionserhaltend geleert und ist für die
  Regel unsichtbar — Sonde Exit 0, der identische Kontroll-Abschnitt ohne
  Spanne Exit 1. `spans` deckt es nicht ab (die Parität ist gerade); der Absatz
  rendert zugleich falsch, sodass die Fundstelle auch dem Leser entgeht. In
  einem Dokument, das wie die Sicht durchgängig Backticks für Rollen- und
  Komponentennamen führt, ist eine unbeabsichtigte Spanne der Normalfall dieses
  Fehlers, nicht der konstruierte.
- **verifizierbar:** ja — zwei identische Abschnitte, einer mit der Spanne über
  dem Pfad; ein Lauf mit `modules: [structure, spans]` meldet nur den
  Kontroll-Abschnitt.
- **klasse:** `bereinigter-text-ohne-mehrzeilen-spannen-probe`

### F-6

- **kategorie:** MEDIUM
- **quelle:** `BEO-009` Richtung (a) (*„sie behauptet eine Probe … die
  Botschaft wird vor oder während der Arbeit formuliert"*) · `BEO-011`
  (Aussage aus dem Anlass statt aus dem Bestand) · slice-145 §2 Schritt 0
  (*„festgehalten am 2026-08-26, vor dem ersten Lauf und **vor jedem Blick in
  die Sicht**"*)
- **pfad:** Commit-Botschaft `701a567` (*„Eigener Commit vor jedem Lauf und vor
  jedem Blick in die Sicht … Eine Schwelle, die den Bestand schon kennt, ist
  keine Schwelle mehr."*) gegen `git show
  147476a:docs/plan/planning/in-progress/slice-145-pfad-token-sicht.md`, §1
- **befund:** Derselbe Slice-Plan trägt drei Absätze über der Schwelle bereits
  das Messergebnis: *„Gemessen: die Sicht trägt heute **null** solcher Token."*
  Der Satz stand dort schon zum Anspruchs-Commit `147476a` (18:21:46, drei
  Minuten vor `701a567`) und stammt aus dem Plan-Datum 2026-08-25; die Zahl
  selbst geht auf slice-136 zurück (*„heute gemessen: `spec/architecture.md`
  null"*). Die erste Schwellen-Bedingung — null Falsch-Positive in der Sicht —
  wurde damit gegen einen bereits bekannten und im selben Dokument notierten
  Bestand gesetzt, und die Begründung, die den eigenen Commit rechtfertigt,
  gilt für ihn nicht. (Die **Form** — eigener Commit, nur die Slice-Datei,
  zeitlich vor dem Feat-Commit — ist eingehalten; der Befund gilt der
  Behauptung, nicht der Reihenfolge.)
- **verifizierbar:** ja — `git show 147476a:<slice-datei> | sed -n '19,26p'`
  gegen die Botschaft von `701a567`; `git log --format=%ad --date=iso` für die
  drei Zeitstempel.
- **klasse:** `schwelle-gegen-bereits-bekannten-bestand`

### F-7

- **kategorie:** LOW
- **quelle:** `spec/lastenheft.md` §`DC-FA-CLI-003` (*„`1` = Prüfung gelaufen,
  mindestens ein Befund; `2` = Nutzungs- oder Umgebungsfehler … die Prüfung hat
  dann **keine verlässliche Aussage** geliefert"*) · slice-145 §2 Schritt 0,
  zweite Bedingung (Positivkontrolle mit **gelesener Ursache**)
- **pfad:** Commit-Botschaft `11054fa`, Gate-Zeile (*„Bewusstes Brechen am
  scharfen Gate: make doc-check Exit 2, Rückbau Exit 0."*)
- **befund:** Der protokollierte Code ist der des `make`-Rahmens, nicht der des
  Produkts: gemessen liefert der Container-Lauf **Exit 1** mit einem
  `section-forbidden`-Befund, `make` meldet dazu `Fehler 1` und beendet sich
  selbst mit 2. Wer die Positivkontrolle später gegen `DC-FA-CLI-003` prüft,
  liest die notierte 2 als Umgebungs-/Konfigurationsfehler — also als Beleg
  dafür, dass der Sensor gerade **nicht** ausgelöst hat. Die Nachbar-Botschaft
  `8e0a6f0` nennt daneben wenigstens die gelesene Ursache (*„Befund nennt den
  Hash selbst"*); hier steht neben dem Code keine.
- **verifizierbar:** ja — Pfad-Zeile anhängen, einmal `docker run … d-check`
  (Exit 1, ein Befund) und einmal `make doc-check` (Exit 2, Ausgabe
  `make: *** [Makefile:135: doc-check] Fehler 1`).
- **klasse:** `make-wrapper-exit-als-produkt-exit-protokolliert`

### F-8

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §4 Kopfsatz (*„Halluzinierte Gates sind die häufigste
  Form von Harness-Lüge"*) · Präzedenz derselben Tabelle: die Zeile
  `make verify-closure-notes` weist die blinden Bereiche der
  Schwester-Regel ausdrücklich aus (*„Vier Bereiche sieht sie nicht"*) ·
  Reviewer-Skill LOW *Doku-Drift*
- **pfad:** `AGENTS.md:196-197` (*„meldet einen Pfad unter den **Code-Wurzeln**
  dieses Repos in der Sicht als `section-forbidden`"*) und `.d-check.yml:344-352`
  gegen `internal/hexagon/core/rules/structure.go:158` (`structureFinding(r,
  file, line, …)` mit `line` = Abschnitts-Zeile) und `:176` (`"verbotenes
  Muster trifft: " + r.ForbidPattern`)
- **befund:** Die Eigenschaft, die die Commit-Botschaft ausdrücklich als
  benennenswert führt — der Befund sitzt auf der Überschriftszeile und die
  Meldung nennt das **Muster** statt des Treffers —, steht weder in `AGENTS.md`
  §3.4 noch im Konfigurations-Kommentar; sie überlebt damit nur in der
  Commit-Historie. Gemessen lautet der Befund `spec/architecture.md:1
  section-forbidden` mit der Meldung `verbotenes Muster trifft:
  \b(internal|cmd)/[a-z]`, während der Verstoß auf Zeile 160 lag. Wer die Regel
  reparieren soll, liest in §3.4 „meldet einen Pfad" und bekommt eine Zeile,
  die keinen Pfad nennt und auf den Datei-Anfang zeigt.
- **verifizierbar:** ja — die Positivkontrolle mit `--json` fahren; das Feld
  `message` trägt das Muster, `line` ist 1.
- **klasse:** `befund-ort-nur-in-der-commit-botschaft`

### F-9

- **kategorie:** INFO
- **quelle:** `AGENTS.md` §3.8 (*„welche Eingaben liest es, die es nicht scannt
  — und gilt dort dieselbe Zusage?"*) · `DC-FA-STRUCT-001` §Bedingungen im
  Abschnitt (fence-treu)
- **pfad:** `spec/architecture.md:38-56` und `:107-143` (die beiden
  mermaid-Fences) gegen `.d-check.yml:349` (*„Fenced Bloecke zaehlen nicht"*)
- **befund:** Die benannte Fence-Grenze hat in **dieser** Datei ein
  ungewöhnliches Gewicht: 56 der 158 Zeilen (35 %) liegen in den beiden
  mermaid-Blöcken, und dort leben die Komponenten-Labels, in denen eine
  sprachkonkrete Benennung am ehesten entstünde. Das Modul `diagrams` liest
  diese Fences, prüft dort aber ausschließlich `ARC-\d{3}`. Repo-weit ist die
  Fence-Grenze dagegen belanglos (2 von 901 Token) — die Grenze ist also
  richtig benannt, ihr Anteil an der einen bewachten Datei aber nirgends.
- **verifizierbar:** nein (Aufteilung der Zeilen, kein Gate-Lauf); der
  beobachtbare Zustand sind die Fence-Grenzen der Datei.
- **klasse:** `fence-grenze-ohne-anteils-angabe-an-der-bewachten-datei`

---

## Negativbefunde

- geprüft, ohne Befund: **Reproduktion der Vier-Kandidaten-Messung.** Alle vier
  Muster einzeln gegen das Produkt gefahren, dieselbe Regel-Form, nur
  `forbid-pattern` variiert: 0 / 0 / 1 / 0 Befunde bei Exit 0 / 0 / 1 / 0 —
  exakt die Zahlen der Botschaft. Auch die Aussage *„DER ZWEITE HÄTTE AUCH NULL
  GEMELDET"* stimmt (`\b(internal|cmd|tools)/[a-z]` ⇒ 0).
- geprüft, ohne Befund: **Nullbestand der Sicht.** `spec/architecture.md`
  trägt heute kein einziges rohes `internal/`- oder `cmd/`-Token — auch nicht
  in Inline-Code oder in den Fences. Der gewählte Kandidat ist also nicht
  „grün, weil blind": er ist grün, weil nichts da ist. Dasselbe gilt für
  `spec/lastenheft.md` und `spec/spezifikation.md` (je 0 rohe Token).
- geprüft, ohne Befund: **Positivkontrolle in fünf Formen.** Prosa und
  Markdown-Link-Ziel melden (`section-forbidden`, Exit 1); Inline-Code, Fenced
  Block und Großschreibung melden nicht — genau wie die Botschaft sagt. Der
  Rückbau ist in beiden Fällen Exit 0, 517 Dateien, 0 Befunde.
- geprüft, ohne Befund: **Abschnitts-Selektor und Kardinalität (§E).**
  `spec/architecture.md` hat **genau eine** `^# `-Überschrift, auf Zeile 1;
  `SectionEnd` liefert für sie 0, der Abschnitt reicht also bis zum Dateiende.
  Vor der H1 liegt nichts. Verliert die Sicht ihre H1 (`##` statt `#`) ⇒
  `section-missing`, Exit 1. Bekommt sie eine **zweite** H1, misst `sections:
  each` beide Abschnitte — der Pfad hinter der zweiten H1 wird gemeldet;
  damit greift die Kappungs-Falle aus slice-143 F-3
  (`abschnitts-spanne-still-gekappt-durch-fremde-h1`) hier **nicht**. Steht die
  H1 nur innerhalb einer Fence ⇒ `section-missing`.
- geprüft, ohne Befund: **Ventile hebeln die Regel nicht aus.**
  `exempt-paths: [spec/architecture.md]` leert die Kandidatenmenge und löst die
  Nullmengen-Härte aus (`section-missing`, Exit 1, „Regel trifft keine Datei"),
  ist also kein stiller Ausweg; `scan.ignore: ["spec/**"]` wirkt gar nicht
  (`structure` ist ein Post-Pass mit eigenem Glob, `DC-FA-STRUCT-001`
  §Kandidaten-Menge) — der Befund erscheint bei „0 Dateien geprüft". Auch die
  beiden Zeilen-Marker des Produkts ziehen nicht: `d-check:ignore` und
  `d-check:status-provenance` auf derselben Zeile lassen den Befund stehen
  (Sonden-Abschnitte S15/S16 melden beide).
- geprüft, ohne Befund: **Weitere Falsch-Negativ-Kandidaten, die halten.**
  Gemessen **gemeldet** werden: HTML-Kommentar, Tabellenzelle, `./internal/…`,
  `github.com/pt9912/d-check/internal/…`, Bild-Ziel `![…](internal/…)`,
  Autolink `<https://…/internal/…>`, `cmd/d-check` und eine
  Unterabschnitts-Überschrift `## internal/…`. Nur die in F-4 genannten
  Formen fallen durch.
- geprüft, ohne Befund: **Die Zusage über die Code-Wurzeln stimmt.**
  `find . -name '*.go'` außerhalb von `internal/` und `cmd/` ist leer; die
  Formulierung *„die Code-Wurzeln dieses Repos (`internal/`, `cmd/`)"* ist
  keine unbelegte Exklusivitäts-Aussage im Sinne von `BEO-011` (a).
  `tools/` enthält ausschließlich Shell-Skripte.
- geprüft, ohne Befund: **Dokument-Pfade sind wirklich außen vor.**
  `\b(internal|cmd)/[a-z]` trifft weder `docs/`, `spec/`, `harness/` noch
  `tools/`; die Aussage des Kommentars ist über das Muster selbst prüfbar und
  hält.
- geprüft, ohne Befund: **Kommentar-Klassen (§3.7).** Der neue Block trägt
  Zusage (*„ZUSAGE, GENAU"*), Abgrenzung und Grenze (*„GRENZEN, benannt"*)
  sowie einen Rang-Zeiger auf die Regel, die er mechanisiert. Er enthält
  **keine** Slice-Nummer, **keine** Lauf-/Mess-Historie und keine
  Review-Marker — der HIGH-Anker, den slice-144 hier kassiert hat
  (`konfig-kommentar-traegt-mess-historie-und-slice-nummer`), ist diesmal
  gemieden. Die Nennung von `AGENTS.md §3.4` **und** `MR-033` liest sich als
  ein Rang-Zeiger plus ein Herkunfts-Feld, nicht als Herkunfts-Prosa; die
  Geschwister-Kommentare derselben Datei tragen dieselbe Form. Der `tools/`-Satz
  ist als **Grenze** formuliert (was die Regel nicht abdeckt), nicht als
  Deliberation über eine verworfene Alternative — knapp, aber innerhalb der
  fünf Klassen. Inhaltlich gilt der Befund F-1, nicht die Klassen-Frage.
- geprüft, ohne Befund: **Keine dritte Mechanik.** Der Slice übernimmt die
  Fähigkeit aus slice-143 (`structure.forbid-pattern`) unverändert; es entsteht
  kein neues Modul, kein neuer Grund-Code, keine neue Konfigurations-Fläche.
  §3 des Plans ist in diesem Punkt eingehalten — der Befund F-3 gilt der
  **Begründung** der Wahl, nicht der Wahl.
- geprüft, ohne Befund: **Keine Ausweitung auf die anderen Straten.** Die Regel
  bindet ausschließlich `spec/architecture.md`; `spec/lastenheft.md` und
  `spec/spezifikation.md` bleiben unberührt (§3 des Plans).
- geprüft, ohne Befund: **Modul-Grenze auf der Ziel-Achse (§3.8).** Die Regel
  liest keine Eingabe außerhalb der Datei, die sie scannt — kein Zielpfad, kein
  git, kein Post-Pass-Verzeichnis. `structureTree` läuft über den Baum, aber
  ausschließlich zur Kandidaten-Findung.
- geprüft, ohne Befund: **Referenz-Richtung (SDP).** Der Diff setzt keinen
  Provenance-Marker und keinen neuen Abwärts-Verweis; `spec/architecture.md`
  ist unverändert.
- geprüft, ohne Befund: **Gate-Stand.** `make doc-check` auf dem Stand von
  `HEAD`: Exit 0, 517 Dateien, 0 Befunde. Die Gate-Zeile *„make gates Exit 0
  (zehn Glieder)"* ist mit der Echo-Zeile des `gates`-Targets konsistent
  (zehn Gates plus `record-gates` als Nachweis-Schritt); `make gates` und
  `make fullbuild` selbst sind Verifikations-Sache und hier nicht gefahren.
- geprüft, ohne Befund: **Doku-Spiegel.** `AGENTS.md` §4 und die
  Sensors-Tabelle in `harness/README.md` beschreiben `make doc-check` generisch
  („Abschnitts-Invarianten (Modul `structure`)") und werden von der neuen Regel
  nicht drift-pflichtig; `MR-033` behält seinen Auflösungs-Trigger zu Recht
  (der Kanon erlaubt der Sicht Modul-Pfade weiterhin, die Adaption bleibt).
  Ein zentrales Trigger-Register außerhalb von `AGENTS.md` existiert nicht.

---

## Kein Finding, aber gemessen und benannt

- Der `structure`-Eintrag steht **vor** den sieben Bestandsregeln, obwohl der
  Block-Kommentar darüber („Zuerst die Chronologie-Monotonie … darunter, als
  siebte Regel, die Kennungs-Pflicht") eine Reihenfolge beschreibt, in der die
  neue Regel nicht vorkommt. Verhaltensrelevant ist das nicht (Regeln sind
  unabhängig, die Befund-Sortierung ist `DC-QA-02`-stabil), und der
  Block-Kommentar ist Bestand; es ist der nächste Kandidat für Doku-Drift, wenn
  eine achte Regel dazukommt.
- Der Selektor `^# ` trifft die Überschriftszeile über einen RE2-Match auf die
  **getrimmte** Zeile (`structureMatcher`), zusätzlich gefiltert durch
  `parseATXHeading`. Eine eingerückte H1 wird damit erfasst — die Lexik-Falle
  aus `DC-FA-STRUCT-001` §Grenzen (eine Bedingung, die ihre eigene
  Heading-Lexik nachbaut) besteht hier nicht.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 5 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:**
`sensor-luecke-als-erlaubnis-in-die-regel-geschrieben` ·
`befundzahl-als-trefferzahl-gelesen` ·
`wegwahl-begruendet-mit-einer-widerlegten-modul-eigenschaft` ·
`falsch-negativ-liste-unvollstaendig` ·
`bereinigter-text-ohne-mehrzeilen-spannen-probe` ·
`schwelle-gegen-bereits-bekannten-bestand` ·
`make-wrapper-exit-als-produkt-exit-protokolliert` ·
`befund-ort-nur-in-der-commit-botschaft` ·
`fence-grenze-ohne-anteils-angabe-an-der-bewachten-datei`

---

## Verdikt

**Merge-blockierend:** ja — ein HIGH und fünf MEDIUM.

Die Wegwahl selbst trägt: `structure` ist dateigenau, braucht keine neue
Klasse, erfindet keine dritte Mechanik, und `sections: each` schließt genau die
Kappungs-Falle, an der die Schwester-Regel aus slice-143 ein MEDIUM kassiert
hat. Der Selektor deckt die Datei ab Zeile 2 vollständig ab, die Nullmengen-Härte
ist echt, keines der drei Ventile hebelt die Regel aus, und die vier
Kandidaten-Zahlen der Botschaft reproduzieren exakt.

Blockierend ist, was **um** den Sensor herum geschrieben wurde. F-1 macht aus
seiner Lexik-Grenze eine Erlaubnis und senkt damit die geschriebene Reichweite
einer Hard Rule auf das, was der Sensor zufällig sieht — gemessen 30 von 901
Vorkommen; die Form, in der dieses Repo Pfade tatsächlich schreibt, ist die
freigegebene. F-2, F-3 und F-6 sind dieselbe Familie wie `BEO-009`/`BEO-011`:
die Messungen stimmen, ihre Deutungen nicht — eine Befundzahl wird als
Trefferzahl gelesen, eine Modul-Eigenschaft aus dem Anlass statt aus der
Nachbarklasse gebildet, und eine Schwelle gegen einen Bestand gesetzt, der drei
Absätze weiter oben schon notiert war. F-4 und F-5 betreffen die dritte
Schwellen-Bedingung: die Falsch-Negativ-Klasse ist zu vier Siebteln benannt,
und die Probe, die das Beobachtungs-Register für genau diesen Akt vorschreibt,
fehlt.
