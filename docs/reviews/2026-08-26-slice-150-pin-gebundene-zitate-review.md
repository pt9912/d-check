# Review-Report: slice-150 — Etappe C-1, „Ein Zitat der Baseline ist pin-gebunden wie ein Link"

**Datum:** 2026-08-26 · **Review-Art:** Plan-/Konventions-Review (geprüft gegen Slice-Plan slice-150 §1/§2/§3/§5/§7, Wellendokument welle-85 §4/§6, den **vendorten** Kanon `v5.12.0` — `grundlagen-harness-dateien.md` §Konventionsspeicher / §Template-Schichtung, `grundlagen-source-precedence.md` §Source Precedence, `modul-04-adrs.md` §Hard Rule für Accepted-ADRs, `modul-02-harness-bootstrap.md` §Freshness-Audit, `modul-03-spec.md` §Ziel-Form: Architektur-Sicht, `templates/AGENTS.template.md` §3.4, `templates/harness/conventions/MR-NNN-titel.template.md`, `templates/docs/plan/adr/NNNN-titel.template.md` —, `AGENTS.md` §3.5/§5, `DC-FA-CITE-001`/`DC-FA-CITE-001.a`, ADR-0045, Beobachtungs-Register `BEO-008`/`BEO-009`/`BEO-011`/`BEO-012`), unabhängiger Reviewer ohne Anteil an der Arbeit
**Gegenstand:** Commit-Kette `17a1eae..HEAD` — `c94dd31` (Feat: `MR-038`, `MR-033`-Ergänzung, slice-152-Schnitt), `0e1b776` (Berichtigung der Zitat-Zusage); **6 Dateien, 195 Einfügungen / 6 Löschungen** (`git diff --stat 17a1eae..HEAD`)
**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 (`9ee805b`) · **Modell-ID:** `claude-opus-5[1m]`

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/done/slice-150-pin-gebundene-zitate.md` (vollständig, inkl. §3 „Ausdrücklich NICHT", §5 Risiken, §7 Sichtung)
- Der neue Eintrag `harness/conventions/MR-038-zitate-pin-gebunden.md` und der ergänzte `harness/conventions/MR-033-sicht-ohne-modul-pfade.md`
- `harness/conventions/MR-021-vendored-verweise-pin-gebunden.md` (der laut Titel geschärfte Eintrag, im Diff unverändert) und `harness/conventions.md` (§Baseline, §Adaptions-Block, beide Index-Tabellen)
- Der vendorte Kanon `.harness/baseline/v5.12.0/` (Regelwerk **und** Templates), gelesen als Datei; zusätzlich der rekonstruierte `v5.11.0`-Baum (`git archive 9ee805b^`) als Vergleichsstand
- `docs/plan/planning/open/slice-152-citations-scharfschalten.md`, `docs/plan/adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md`, `spec/lastenheft.md` §`DC-FA-CITE-001`, `spec/spezifikation.md` §`DC-FA-CITE-001.a`, `internal/hexagon/core/rules/citations.go`
- `docs/plan/planning/observations.md` — `BEO-008`, `BEO-009`, `BEO-011`, `BEO-012`
- **Vorherige Findings am gleichen Vorgang:** `docs/reviews/2026-08-26-slice-149-delta-audit-review.md` (sieben MEDIUM, u. a. `zahl-ohne-messvorschrift`, `geltungsbereich-einer-eigenen-regel-gedehnt`, `zahl-aus-der-quelle-umgedeutet`) und `docs/reviews/2026-08-26-slice-148-baseline-v5120-review.md` (zwei MEDIUM zur Form des neuen Eintrags: fehlendes `Löst auf`, fehlendes `Begründung`)
- **Zusätzlich hinzugezogen** (nicht im Auftrag genannt, für D relevant): `docs/reviews/2026-07-18-slice-079-implementation-r1.md` F-3

**Nicht erhalten** (Skill §Eingangs-Kontext): die DoD-Abhakung als Prüf-Gegenstand — Plan-/DoD-Konformität prüft die Verifikation. Die vier offenen DoD-Haken in `slice-150` §4 werden hier **nicht** bewertet; der Slice liegt in `in-progress/` und ist nicht geschlossen.

**Vom Reviewer selbst gefahren** (Exit je Lauf in eine Datei umgeleitet und direkt gelesen):

- `make gates` Exit **0** — zehn Glieder (`baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check`), `targets`/`planning` je 512 Dateien / 0 Befunde.
- **Der Probelauf der Behauptung (b):** `docker run --rm --network none -v "$PWD":/repo:ro d-check:latest --enable citations` ⇒ Exit **2**, `d-check: error: CHANGELOG.md:592: malformte d-check:cite-Direktive — erwartet <!-- d-check:cite <pfad>:<von>-<bis> --> (DC-FA-CITE-001.a Schritt 1, fail-closed)`.
- **Eigener Zensus der Direktiven-Erwähnungen**: Nachbau von `proseLines`/`FenceToggle` (Python) über alle gescannten `.md` ⇒ **19 Treffer in 10 Dateien**, davon **0** innerhalb einer Fence. Anschließend **jede der zehn Dateien einzeln** in ein Scratchpad-Repo kopiert (`.d-check.yml` mit `modules: [citations]`) und gefahren — **alle zehn** brechen mit Exit 2.
- **Eigene Zitat-Messung**: alle `„…"`-Spannen der 17 aktiven `MR-`Einträge extrahiert (auf dem **ganzen** Dateitext, damit Umbrüche nicht zählen; Schluss-Zeichen `"` U+0022 **und** `“` U+201C berücksichtigt) und gegen den `v5.12.0`- **und** den `v5.11.0`-Baum gehalten.
- **Gegenprobe mit dem Produkt statt mit `grep`:** neun der Zitate mit einer echten `d-check:cite`-Direktive in einem Scratchpad-Repo gegen ihre Quell-Spanne gefahren (`modules: [citations]`) — fünf `citation-mismatch`, vier ohne Befund.
- **Keine Repo-Datei wurde verändert** — alle Sonden liegen unter `/tmp/claude-1000/.../scratchpad/`, kein `git checkout --` nötig, `git status --short` nach jedem Lauf leer. Die einzige geschriebene Datei ist dieser Report.

**Verdikt: blockierend** — ein HIGH, sieben MEDIUM, sechs LOW.

---

## Findings

### F-1

- **kategorie:** HIGH
- **quelle:** `MR-038` §Adaption (*„wird das Zitat **ergänzt, nicht ersetzt**"* / *„Es wird nichts entfernt"*) · `.harness/baseline/v5.12.0/regelwerk/grundlagen-harness-dateien.md` §Konventionsspeicher (*„Einträge werden nie überschrieben."*) · `BEO-008` (*„eines lebend (`MR-033`, **zwei Zitate**)"*)
- **pfad:** `harness/conventions/MR-033-sicht-ohne-modul-pfade.md:4-17` (Diff `c94dd31`, **12 Einfügungen / 4 Löschungen**)
- **befund:** Der Bestandsfall bestand aus **zwei** veralteten Zitaten; die neue Regel wurde nur auf eines angewandt. Das `v5.11.0`-Zitat aus `modul-03-spec.md` (*„Sprach- und meilensteinfrei — referenziert Modul-Pfade, aber …"*) ist ersatzlos **entfernt** und durch den neuen `v5.12.0`-Satz *„Die Erlaubnis ist keine Pflicht"* ersetzt worden — der alte Wortlaut steht nirgends mehr im Eintrag, während der Block *„Fassung, gegen die dieser Eintrag geschrieben wurde"* nur die `AGENTS.template.md`-Hälfte bewahrt.
- **verifizierbar:** nein — kein Gate liest Eintragskörper (`make gates` ist auf genau diesem Stand Exit 0, zehn Glieder). Belegt durch `git show c94dd31 -- harness/conventions/MR-033-sicht-ohne-modul-pfade.md` (die vier `-`-Zeilen) gegen `git show 9ee805b^:.harness/baseline/v5.11.0/regelwerk/modul-03-spec.md:120-121` und gegen `harness/conventions/MR-033-sicht-ohne-modul-pfade.md:12-17`.
- **klasse:** `neue-regel-im-eigenen-einfuehrungsfall-verletzt`

### F-2

- **kategorie:** MEDIUM
- **quelle:** `BEO-009` Richtung (a) — *„sie behauptet eine Probe oder Änderung, die **nicht stattfand**"* · `AGENTS.md` §5
- **pfad:** Commit-Botschaft `c94dd31`, Absatz 4 (*„MR-033 ist entsprechend ergänzt … **Nichts entfernt, nichts umgeschrieben** — dieselbe Bewegung wie ein Geschichte-Anhang an einer ADR."*)
- **befund:** Der Diff derselben Botschaft entfernt vier Zeilen aus `MR-033` (`git show --numstat c94dd31` ⇒ `12 4`), darunter das vollständige `modul-03-spec.md`-Zitat, und stellt die verbleibende `AGENTS.template.md`-Fassung von der ersten in die zweite Position um. Die Botschaft ist der einzige Ort, an dem ein späterer Leser die Bewegung nachvollzieht, und sie beschreibt das Gegenteil des Diffs.
- **verifizierbar:** nein — kein Gate vergleicht Botschaft und Diff (der genannte `BEO-009`-Hook deckt Datei-**Pfade**, nicht Aussagen über Zeilen). Belegt durch `git show --numstat c94dd31` und den Diff-Hunk selbst.
- **klasse:** `botschaft-behauptet-eine-nicht-stattgefundene-nicht-aenderung`

### F-3

- **kategorie:** MEDIUM
- **quelle:** `.harness/baseline/v5.12.0/regelwerk/grundlagen-harness-dateien.md` §Konventionsspeicher (*„schärft er ihn nur (der alte gilt weiter, die Regel wird **strenger**), steht das im Titel"*) · `MR-021` §Geltungsbereich · `BEO-008` (vierte Spiegel-Klasse)
- **pfad:** `harness/conventions/MR-038-zitate-pin-gebunden.md:1` (Titel) gegen `:14-16` (Geltungsbereich) und `harness/conventions/MR-021-vendored-verweise-pin-gebunden.md:6-14`
- **befund:** Der Titel sagt *„(schärft `MR-021`)"*, aber der Geltungsbereich der Schärfung (*„aktiven Einträgen unter `harness/conventions/`"*) ist ein echter Teil des geschärften Eintrags, dessen Geltungsbereich *„das Briefing, `harness/README.md`, dieser Konventionsspeicher …, **der Reviewer-Skill** und die **lebenden Planungs-Dokumente**"* ausdrücklich mitnennt. In genau diesen ausgeschlossenen Dateien stehen heute **neun** wörtliche Baseline-Zitate (`AGENTS.md`, `.harness/skills/reviewer.md` ×2, `.harness/skills/closure-note-reviewer.md` ×2, `docs/plan/planning/open/slice-151-…md`, `docs/plan/planning/welle-85-…md` ×3), für die beim nächsten Bump keine Ergänzungspflicht gilt und die kein Gate sieht.
- **verifizierbar:** nein — kein Gate kennt Geltungsbereiche. Belegt durch eine eigene Zitat-Extraktion über die lebenden Dokumente, jede Spanne gegen den `v5.12.0`-Baum aufgelöst (Fundstellen oben namentlich), gegen `MR-021:6-14` und `MR-038:14-16`.
- **klasse:** `schaerfung-reicht-weniger-weit-als-der-geschaerfte-eintrag`

### F-4

- **kategorie:** MEDIUM
- **quelle:** `BEO-011` (*„Eine Zahl ohne ihre Vorschrift ist keine Messung"*) · `DC-FA-CITE-001.a` Schritt 4 (whitespace-Normalisierung) · vorheriges Finding am gleichen Vorgang: slice-149-Review F-8 `zahl-ohne-messvorschrift`
- **pfad:** `harness/conventions/MR-038-zitate-pin-gebunden.md:38-44`
- **befund:** *„Von den 22 Zitaten … sind **neun** Zitate der Baseline, und die übrigen acht davon sind gegen `v5.12.0` wortgleich (normalisiert gemessen)"* nennt weder, was als Zitat zählt, noch welche Normalisierung gemeint ist; unter der Normalisierung, die der Eintrag selbst zum Gegenstand hat (whitespace-only, `citations.go:39/238`), liefern **fünf** Baseline-Zitate aktiver Einträge `citation-mismatch` — `MR-032` (zwei), `MR-033`s neue `v5.12.0`-Fassung, `MR-035` und `MR-038`s eigenes Kanon-Zitat. Meine Klassifikation zählt zudem **elf** Baseline-Zitate, nicht neun, sodass die Aufteilung 9/13 aus dem Bestand nicht reproduzierbar ist.
- **verifizierbar:** ja — die fünf Fehlschläge sind mit dem Produkt reproduziert: Zitattext plus echter `<!-- d-check:cite … -->`-Direktive in einem Scratchpad-Repo mit `modules: [citations]` ⇒ Exit 1, fünfmal `citation-mismatch`; die vier geprüften Treffer (`MR-005`, `MR-031` ×3) ⇒ Exit 0, 0 Befunde.
- **klasse:** `wortgleich-zusage-ohne-normalisierungs-vorschrift`

### F-5

- **kategorie:** MEDIUM
- **quelle:** `BEO-012` (Reichweite eines Zitats, Zähler 4) · `slice-150` §4 DoD (*„die Reparatur belegt (Zitat gegen das neue Ziel)"*)
- **pfad:** `harness/conventions/MR-033-sicht-ohne-modul-pfade.md:5-7` (neue Fassung) und `:13-15` (historische Fassung)
- **befund:** Beide Fassungen des reparierten Bestandsfalls tragen Auszeichnung, die ihre Quelle nicht trägt: die neue zitiert *„`spec/architecture.md` **darf** Pfade zu **Code-Modulen** referenzieren"*, während `AGENTS.template.md:130` *„`spec/architecture.md` darf Pfade zu **Code-Modulen** referenzieren"* schreibt — die Hervorhebung liegt ausgerechnet auf dem Wort, um dessen Erlaubnis-Charakter der ganze Eintrag streitet; die historische Fassung setzt *„**referenziert Modul-Pfade**"* fett und nimmt der Quelle zugleich das Fett auf *„**keine**"* (`v5.11.0` `AGENTS.template.md:130-131`).
- **verifizierbar:** ja — dieselbe Sonde wie F-4: die neue Fassung gegen `.harness/baseline/v5.12.0/templates/AGENTS.template.md:130-132` ⇒ `citation-mismatch`. Die historische Fassung ist gegen `git show 9ee805b^:.harness/baseline/v5.11.0/templates/AGENTS.template.md` Zeile für Zeile geprüft.
- **klasse:** `zitat-mit-hinzugefuegter-auszeichnung`

### F-6

- **kategorie:** MEDIUM
- **quelle:** `BEO-011` (Vollständigkeits-Aussagen) · `slice-152` §2 Schritt 2 (Wegwahl mit benanntem Preis) · `AGENTS.md` §3.5 (ADRs nach `Accepted` immutabel)
- **pfad:** Commit-Botschaft `c94dd31`, Absatz 7 · `docs/plan/planning/open/slice-152-citations-scharfschalten.md:43-45` · `docs/plan/planning/in-progress/roadmap.md:143`
- **befund:** Der Zensus nennt sechs Dokumente (CHANGELOG, beide READMEs, Lastenheft, Spezifikation, Handbuch); tatsächlich brechen **zehn** den Lauf — ungenannt bleiben `docs/plan/adr/0045-…md` und die drei Reporte `docs/reviews/2026-07-18-slice-079-{design-r1,implementation-r1,realdatenbeleg-ai-harness-init}.md`. Das sind genau die vier, für die der in `slice-152` §2 zuerst genannte Weg (Syntax nur noch in Fenced-Blöcken) nicht offensteht — eine `Accepted`-ADR ist immutabel, Review-Reporte sind Lauf-Belege, die `slice-152` §3 ausdrücklich nicht anfasst.
- **verifizierbar:** ja — jede der zehn Dateien einzeln in ein Scratchpad-Repo (`modules: [citations]`) kopiert und gefahren: **zehnmal Exit 2**. Der Marker-Zensus (Nachbau von `proseLines`) liefert 19 Treffer in denselben zehn Dateien, keiner in einer Fence.
- **klasse:** `zensus-nennt-nur-die-editierbaren-fundstellen`

### F-7

- **kategorie:** MEDIUM
- **quelle:** `BEO-012` (*„vor jedem Zitat einer Regel deren Geltungs-Feld lesen"*) · `.harness/baseline/v5.12.0/templates/docs/plan/adr/NNNN-titel.template.md:124-127` · `AGENTS.md` §3.5
- **pfad:** `harness/conventions/MR-038-zitate-pin-gebunden.md:34-36` (*„Dieselbe Bewegung wie der `## Geschichte`-Anhang einer ADR, die der Kanon ausdrücklich zulässt."*)
- **befund:** Der Kanon lässt den Anhang nicht ausdrücklich zu — seine ADR-Vorlage sagt unter `## Geschichte` selbst *„Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**"*, und die einzige andere Kanon-Aussage zu der Überschrift (`grundlagen-referenz-richtung.md:222`) nimmt sie vom **Token-Gate** aus, nicht von der Immutabilität; die Erlaubnis stammt aus `AGENTS.md` §3.5 und ihrer repo-lokalen Mechanik (`ADR-0016`, `.d-check.yml:258` `exclude-sections: [Geschichte]`). Der Vergleich trägt auch strukturell nicht: `## Geschichte` ist ein vom gehashten Kern **ausgenommener** Abschnitt, während die Änderung an `MR-033` im Pflichtfeld `Ersetzt-Baseline-Regel` liegt, für das keine Ausnahme deklariert ist.
- **verifizierbar:** nein — kein Gate prüft Kanon-Zuschreibungen. Belegt durch `grep -rn "Geschichte" .harness/baseline/v5.12.0/` (sechs Treffer, keiner mit einer Erlaubnis für Anhänge nach `Accepted`) gegen `AGENTS.md:195-200` und `.d-check.yml:248-258`.
- **klasse:** `repo-lokale-erlaubnis-dem-kanon-zugeschrieben`

### F-8

- **kategorie:** MEDIUM
- **quelle:** `slice-150` §5 Risiko 1 (*„Die bequeme Antwort ist die Pflege-Antwort … die Begründung muss aus dem Kanon kommen, nicht aus dem Ergebnis"*) · `.harness/baseline/v5.12.0/regelwerk/modul-02-harness-bootstrap.md` §Freshness-Audit
- **pfad:** `harness/conventions/MR-038-zitate-pin-gebunden.md:4-12` · gleichlautend `harness/conventions.md:120` (Index-Zeile *„keine — der Kanon regelt den Fall nicht"*) und Commit-Botschaft `c94dd31`, Absatz 2 (*„sagt keine Stelle"*)
- **befund:** Die Prämisse ist an genau zwei Stellen geprüft; drei weitere Kanon-Sätze sprechen zur Frage, ob und wie ein bestehender `MR-`Eintrag geändert wird, und keiner ist genannt oder gewogen — `modul-02` §Freshness-Audit *„**Rückbau ist ein neuer Eintrag, kein Edit** … Die alte Zeile ist die historisch korrekte Aussage über den damaligen Zustand"*, ebd. für wiederkehrende Vorlagen *„bestehende werden nicht rückwirkend umgeschrieben"*, und `grundlagen-harness-dateien.md:49` *„**Der Zeiger ist kein Zitat.** Ein Template, das den Normtext ausschreibt, führt ihn ein zweites Mal — und zwei Fassungen driften"* (der Zustand, den `MR-038` dauerhaft herstellt). Alle drei weisen in die Gegenrichtung der gewählten Antwort; die Feststellung *„der Kanon regelt den Fall nicht"* steht jetzt als stehende Index-Zeile, die der nächste Freshness-Audit liest, statt sie zu wiederholen.
- **verifizierbar:** nein — kein Gate liest Begründungen. Belegt durch `grep -rn "Zitat\|zitier\|Wortlaut\|wörtlich\|verbatim" .harness/baseline/v5.12.0/` über **Regelwerk und Templates** (nicht nur die drei im Auftrag genannten Dateien) und `sed -n '285,300p' .harness/baseline/v5.12.0/regelwerk/modul-02-harness-bootstrap.md`.
- **klasse:** `kanon-suche-auf-zwei-stellen-verengt`

### F-9

- **kategorie:** LOW
- **quelle:** `DC-FA-CITE-001.a` Schritt 1 vs. Schritt 2 (`citations.go:77`, `:94`) · `slice-152` §2 Schritt 1
- **pfad:** Commit-Botschaft `c94dd31`, Absatz 7 (*„Der Scan ist fence-bewusst, aber nicht inline-code-bewusst; die Syntax steht dort in Backticks und ist damit von einer echten Direktive nicht zu unterscheiden"*) · `docs/plan/planning/open/slice-152-citations-scharfschalten.md:36-46`
- **befund:** Der Lauf bricht an **zwei** verschiedenen Stellen des Algorithmus: fünfzehn Erwähnungen scheitern an Schritt 1 (*malformte Direktive*), zwei Dateien aber an Schritt 2 (*`d-check:cite` ohne folgendes Zitat*, `docs/reviews/2026-07-18-slice-079-implementation-r1.md:48` und `…realdatenbeleg-ai-harness-init.md:36`) — dort ist die dokumentierte Beispiel-Direktive vollständig und **parst**. Die Diagnose nennt nur die erste Form, und `slice-152` §2 Schritt 1 zählt entsprechend nur „Inline-Code gegen Fenced-Block".
- **verifizierbar:** ja — die Einzel-Läufe der zehn Dateien liefern die zwei abweichenden Fehlermeldungen wörtlich.
- **klasse:** `diagnose-nennt-einen-von-zwei-bruchpfaden`

### F-10

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §5 (Gate-Aussagen mit Exit) · `BEO-009` · `slice-152` §4 DoD (*„Exit und Befundzahl genannt"*)
- **pfad:** Commit-Botschaft `c94dd31`, Absatz 8 (*„Die Probe lief mit einer echten Direktive auf ein bestätigtes Zitat und wurde sofort zurückgebaut; der Arbeitsbaum war danach sauber."*)
- **befund:** Die Probe ist der einzige Beleg dafür, dass die Mechanik am eigenen Bestand *positiv* trägt, nennt aber weder Exit noch Umfang noch das geprüfte Zitat und ist damit nicht nachfahrbar; ein Lauf über den Bestand kann sie nicht gewesen sein, weil er fail-closed an `CHANGELOG.md:592` endet, bevor ein Eintrag erreicht wird. Von den elf Baseline-Zitaten aktiver Einträge scheitern fünf an genau dieser Prüfung (F-4), sodass die Auswahl des „bestätigten" Zitats das Ergebnis bestimmt.
- **verifizierbar:** ja, in der Gegenrichtung — mein Lauf `--enable citations` über die Repo-Wurzel endet Exit 2 an `CHANGELOG.md:592`.
- **klasse:** `probe-ohne-exit-und-ohne-gegenstand`

### F-11

- **kategorie:** LOW
- **quelle:** `.harness/baseline/v5.12.0/templates/harness/conventions/MR-NNN-titel.template.md:13-15` (*„`Löst auf` und `Ausgelöst durch Baseline-Stand` nur, wenn dieser Eintrag einen früheren ablöst"*) und `:29` (*„Pflicht zusammen mit „Löst auf""*) · vorheriges Finding am gleichen Vorgang: slice-148-Review F-2 (dieselbe Feld-Paarung, andere Hälfte)
- **pfad:** `harness/conventions/MR-038-zitate-pin-gebunden.md:50`
- **befund:** Der Eintrag trägt `Ausgelöst durch Baseline-Stand: v5.12.0`, ist aber ausweislich seines Titels und seiner Index-Zeile **keine** Ablösung, sondern eine Schärfung; das Feld ist nach der Vorlage nur die eine Hälfte eines Ablöse-Paares. Weil der Kanon an die zwei Klassen verschiedene Folgen knüpft (*der alte gilt weiter* vs. *engerer Nachfolger*), ist das genau das Feld, an dem ein späterer Freshness-Audit `MR-021` fälschlich als abgelöst lesen kann.
- **verifizierbar:** nein — kein Gate prüft die Feld-Kombination. Belegt durch den Feld-Zensus über alle 17 aktiven Einträge (`MR-031`–`MR-036` tragen dieselbe Kombination, `MR-037` als einziger Ablöser trägt sie ohne `Löst auf`) gegen die Vorlage.
- **klasse:** `abloese-feld-in-einer-schaerfung`

### F-12

- **kategorie:** LOW
- **quelle:** `BEO-012` · `slice-150` §1 (der Gegenstand des Slice ist der Wortlaut eines Zitats)
- **pfad:** `docs/plan/planning/open/slice-152-citations-scharfschalten.md:36-38`
- **befund:** Der als Lauf-Ausgabe gesetzte `` ```text ``-Block gibt die Meldung gekürzt und ohne Auslassungszeichen wieder; das Werkzeug schreibt `… malformte d-check:cite-Direktive — erwartet <!-- d-check:cite <pfad>:<von>-<bis> --> (DC-FA-CITE-001.a Schritt 1, fail-closed)`. Der Block ist der einzige Reproduktions-Anker des Slice, und die Herkunftsangabe (`Schritt 1`), die den zweiten Bruchpfad aus F-9 unterscheidbar machte, fällt dabei weg.
- **verifizierbar:** ja — mein eigener Lauf liefert die volle Zeile.
- **klasse:** `gekuerztes-lauf-zitat-ohne-auslassungszeichen`

### F-13

- **kategorie:** LOW
- **quelle:** `slice-152` §7 (Sichtung offener Beobachtungen) · `docs/reviews/2026-07-18-slice-079-implementation-r1.md` F-3
- **pfad:** Commit-Botschaft `c94dd31`, Absätze 6–7 (*„mit einem Ergebnis, das ich nicht erwartet hatte"* / *„auch das ist gemessen statt vermutet"*) · `docs/plan/planning/open/slice-152-citations-scharfschalten.md:1-46`
- **befund:** Derselbe Blocker steht seit dem 18.07.2026 als INFO-Finding im Repo, mit derselben Datei-Aufzählung (*„wie d-check es in README/Handbuch/CHANGELOG tut"*), derselben Ursache und einer ausdrücklichen Einordnung (*„Bewusste, dokumentierte Fail-closed-Semantik"*) samt der damals festgehaltenen Entschärfung (*„d-checks **eigene** `.d-check.yml` aktiviert `citations` **nicht**"*). Weder die Botschaft noch `slice-152` §7 nennt ihn, sodass die Wegwahl in §2 ohne die bereits protokollierte Abwägung getroffen würde.
- **verifizierbar:** nein — kein Gate verknüpft Slices mit alten Reports. Belegt durch `sed -n '85,104p' docs/reviews/2026-07-18-slice-079-implementation-r1.md`.
- **klasse:** `bekannter-befund-als-neuentdeckung-gerahmt`

### F-14

- **kategorie:** LOW
- **quelle:** `.harness/baseline/v5.12.0/regelwerk/grundlagen-harness-dateien.md` §Konventionsspeicher (*„Einträge werden nie überschrieben."*) · `slice-150` §3 (*„Einträge werden nicht überschrieben."*)
- **pfad:** `harness/conventions/MR-038-zitate-pin-gebunden.md:38-44`, Diff `0e1b776` (**5 Einfügungen / 2 Löschungen**)
- **befund:** Die Berichtigung ersetzt zwei Zeilen im Körper eines Eintrags, der seit seinem ersten Commit `Status: Accepted` trägt; die abgelöste Zusage (*„Alle übrigen Zitate aktiver Einträge sind gegen `v5.12.0` wortgleich"*) steht im Eintrag nicht mehr und lebt nur noch in der Commit-Botschaft. Die Kette setzt damit im selben Zug, in dem sie die Frage „Überschreiben oder Pflege?" entscheidet, einen unmarkierten dritten Fall — die Korrektur einer eigenen Fehlmessung —, auf den sich der nächste Bearbeiter berufen kann.
- **verifizierbar:** nein — kein Gate prüft Eintrags-Immutabilität (`adr-check`/`vcs` decken `docs/plan/adr/`, nicht `harness/conventions/`; `.d-check.yml:352-357`). Belegt durch `git show 0e1b776` gegen `harness/conventions/MR-038-zitate-pin-gebunden.md:3`.
- **klasse:** `accepted-eintrag-im-koerper-korrigiert`

---

## Negativbefunde

1. **`make gates` ist grün, und die Zahl stimmt.** Exit **0**, Abschlusszeile nennt genau zehn Glieder (`baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check`); `targets` und `planning` je 512 Dateien / 0 Befunde. Die Gate-Zusage beider Botschaften (*„make gates Exit 0 (zehn Glieder)"*) ist eigenständig reproduziert.
2. **Die Kern-Prämisse hält für die zwei genannten Stellen — gezielt und breiter gesucht.** Ich habe den **gesamten** vendorten Baum (Regelwerk **und** Templates, 49 Dateien) nach `Zitat`/`zitier`/`Wortlaut`/`wörtlich`/`verbatim`/`Fassung` durchsucht. Weder `modul-04-adrs.md` §Hard Rule (*„Eine ADR mit Status `Accepted` wird nicht **inhaltlich** überschrieben"*) noch `grundlagen-source-precedence.md` §Source Precedence sagt, ob ein Zitat der Baseline zum geschützten Kern oder zur pin-gebundenen Referenz gehört. Die Formulierung des Eintrags (*„sagt keine der **beiden** Stellen"*) ist damit richtig; F-8 betrifft die **Verallgemeinerung** in der Botschaft und die drei ungewogenen Nachbarsätze, nicht diese Aussage.
3. **Das Modul `citations` hat die Form, die der Slice ihm zuschreibt — am Code gelesen, nicht aus dem CHANGELOG.** `citations.go:31-40` erkennt ausschließlich die Kommentar-Form `<!-- d-check:cite` (keine bloße String-Erwähnung reicht als *Direktive* im Sinne des Parsers), `:107-118` normalisiert Quell-Spanne und Zitattext whitespace-weit und prüft **Teilstring**, `:117` erzwingt eine Mindestlänge von 16 Runen, `:77`/`:94` sind fail-closed. Das ist genau die vom Slice geforderte Prüfung *Zitat gegen sein eigenes Link-Ziel* und nicht der in `slice-148` gescheiterte Korpus-Test. Die Aussage *„es gibt es her"* trägt.
4. **Der fail-closed-Bruch ist echt, an der behaupteten Zeile und aus dem behaupteten Grund.** Eigener Lauf: Exit 2, `CHANGELOG.md:592`; `sed -n '590,594p' CHANGELOG.md` zeigt die Direktiv-Syntax in **Inline-Code** innerhalb eines Aufzählungspunkts, nicht in einer Fence. Die Erklärung *„fence-bewusst, aber nicht inline-code-bewusst"* stimmt mit `proseLines` (`markdown.go:85-99`) überein: Fence-Zeilen entfallen, Inline-Code bleibt roh stehen. `citations` ist in `.d-check.yml:16` **nicht** aktiviert — die Selbst-Betroffenheit ist keine Regression von `make gates`.
5. **Die Form-Lücke, die der Vorgänger-Review an `MR-037` fand, ist nicht wiederholt.** `MR-038` trägt alle Pflichtfelder der Pflichtgliederung (`grundlagen-harness-dateien.md:234`): `Datum` (`:13`), `Geltungsbereich` (`:14`), `Ersetzt-Baseline-Regel` (`:4`), `Adaption` (`:17`), **`Begründung`** (`:45`) und `Auflösungs-Trigger` (`:51`) — das Feld `Begründung`, dessen Fehlen die slice-148-Review-F-4 für `MR-032`–`MR-037` beanstandete, ist wieder da. Beanstandet bleibt in F-11 nur die überzählige Ablöse-Hälfte.
6. **Titel-Form und Index-Zeile folgen dem Kanon.** Die Kanon-Vorgabe verlangt die Schärfung **im Titel** (`(schärft MR-<NNN>)`) — `MR-038:1` trägt sie dort und wiederholt sie in der Titel-Spalte der Index-Zeile (`harness/conventions.md:120`), nicht in der Spalte `Ersetzt-Baseline-Regel` (wo `MR-034` sie als abweichende Bestandsform führt). Die Zeile steht nur in §Aktive Adaptionen, in aufsteigender Ordnung am Ende, mit beiden `<a id>`-Ankern (Voll-Slug und `#mr-038`); alle Verweise `conventions.md#mr-038` aus `MR-033`, `slice-152` und `MR-038` selbst lösen auf (`make gates` mit `links`+`anchors` grün). §Baseline und §Adoptierte Konventions-Quellen sind unberührt und weiterhin auf `v5.12.0` konsistent. *Benannt, nicht gemeldet:* die zwei Kanon-Links im Feld `Ersetzt-Baseline-Regel` (`:7`, `:9`) tragen keinen Anker, obwohl die Vorlage *„als Link mit Anker"* schreibt — das Feld führt hier den Wert `keine`, die Links sind erläuternd, und `MR-035`/`MR-036` setzen dieselbe Form.
7. **Der `MR-033`-Bestandsfall ist inhaltlich richtig entschieden.** Die neue `v5.12.0`-Fassung gibt die Kanon-Aussage sachlich korrekt wieder (`AGENTS.template.md:130-132` erlaubt Code-Modul-Pfade), der Zusatz *„Die Erlaubnis ist keine Pflicht"* steht wörtlich in `modul-03-spec.md:127`, und die Schlussfolgerung *„damit bleibt dieses Repo bei einer Verschärfung"* trägt: `AGENTS.md` §3.4 formuliert weiterhin ein Verbot. Die Beanstandungen F-1 und F-5 betreffen die **Bewegung** und den **Wortlaut**, nicht das Ergebnis.
8. **Die Berichtigung `0e1b776` ist in ihrem Gegenstand ehrlich.** Sie schneidet eine zu breite Zusage zurück, benennt die Richtung (`BEO-009` (b)) und die Wiederholung („heute zum dritten Mal") selbst und verschweigt nicht, dass sie vor dem Review gefunden wurde. Ein Fund, den sie hätte mitnehmen können, ist die eigene Zahl (F-4); F-14 betrifft die **Ablage-Form** der Korrektur, nicht ihren Inhalt.
9. **Der Schnitt von `slice-152` ist vollständig gebucht und formgerecht.** Der Slice steht in `welle-85` §4 mit Rolle (`Etappe C-3`), im Drift-Log der Roadmap mit Datum und Begründung (`roadmap.md:143`), trägt `Verantwortlich:`/`Autor:`/`Datum:`, ein `Berührte Spec-Stellen`-Feld mit `DC-FA-CITE-001`-Kennung (nicht nur `§N`), §3 „Ausdrücklich NICHT", §7-Sichtung, §8-Modus-Begründung und zwei §5-Ausgänge in der Repo-üblichen Platzhalter-Form. `slice-150` liegt in `in-progress/`, die Roadmap führt keinen Ruhe-Marker — `planning-check` grün.
10. **Die Zählbasis „22 Zitate" ist mit einer Vorschrift rekonstruierbar, wenn auch nicht mit ihrer eigenen.** Die 17 aktiven Einträge tragen **23** `„`-Öffner; zieht man `MR-025`s `„slice-099“` ab (eine Kennung in Anführungszeichen, kein Zitat, und der einzige `„…“`-Fall des Bestands), bleiben genau 22. Die zwei reinen `"…"`-Paare in `MR-005` und `MR-007` fallen dabei heraus. Umbrüche habe ich ausgeschlossen, indem ich auf dem gesamten Dateitext statt zeilenweise extrahiert habe — die Falle aus `slice-148` ist hier nicht eingetreten. Die Beanstandung in F-4 betrifft die **Aufteilung** 9/8/13, nicht die 22.
11. **Kein Zitat außerhalb des Konventionsspeichers ist durch den `v5.12.0`-Bump veraltet.** Eigene Gegenprobe über `AGENTS.md`, `harness/README.md`, beide Skills, `harness/conventions.md` und alle lebenden Planungs-Dokumente: jede `„…"`-Spanne ab 20 Zeichen gegen **beide** gepinnten Bäume gehalten — **null** Treffer, die nur in `v5.11.0` stehen. Die Aussage *„Anwendungsfall bei Einführung: genau einer"* ist damit für den heutigen Bestand richtig; F-3 betrifft den **nächsten** Bump, nicht diesen.
12. **§3 des Slice-Plans ist eingehalten.** Der Diff `17a1eae..HEAD` fasst `MR-021` nicht an, editiert kein `done/`-Dokument und keinen Review-Report, baut kein Gate und aktiviert `citations` nicht — die drei Selbstbeschränkungen halten am Diff. Auch der vendorte Baum ist unberührt (`git diff --stat 17a1eae..HEAD` nennt sechs Dateien, keine unter `.harness/`).
13. **Referenz-Richtung und Marker-Ehrlichkeit ohne Befund.** Keines der neuen Dokumente trägt einen `d-check:status-provenance`-Marker; die Abwärts-Verweise in `slice-152` (auf ADR-0045, `BEO-*`, `MR-*`) laufen von einem Planungs- in Entscheidungs-/Planungs-Artefakte und verletzen die SDP-Matrix nicht. `make gates` mit aktivem `matrix` ist grün.
14. **Kein Zustandsfeld trägt Chronik.** `MR-038` und der ergänzte Teil von `MR-033` führen kein `Status:`-Feld mit Entstehungsgeschichte; die Drift-Log-Zeile der Roadmap protokolliert einen Schnitt (Kanon-konform), keine Schließung. Der HIGH-Anker des Skills ist nicht berührt.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 1 | F-1 |
| MEDIUM | 7 | F-2, F-3, F-4, F-5, F-6, F-7, F-8 |
| LOW | 6 | F-9, F-10, F-11, F-12, F-13, F-14 |
| INFO | 0 | — |

## Verdikt

**Blockierend.** Ein HIGH, sieben MEDIUM.

Die **Sacharbeit** ist an drei Stellen ohne Befund, und das sind keine kleinen. Die Regelfrage ist an den zwei genannten Kanon-Stellen sauber geprüft — ich habe den ganzen vendorten Baum durchsucht und finde dort keine Antwort, die der Eintrag übersehen hätte (N-2). Die Messung zum Produkt ist am **Code** richtig, nicht nur am CHANGELOG: `citations` prüft genau *Zitat gegen sein eigenes Link-Ziel*, whitespace-normalisiert, fail-closed (N-3). Und der Bruch ist echt, an der behaupteten Zeile, aus der behaupteten Ursache (N-4). Der Slice hat gemessen, wo er messen wollte.

Blockierend ist, was mit dem Ergebnis geschah. **Der Eintrag wird in seinem eigenen Einführungsfall verletzt** (F-1): `MR-033` trug **zwei** veraltete Zitate — `BEO-008` sagt es, `slice-150` §1 sagt es —, und die Reparatur hat eines ergänzt und das andere entfernt. Der `v5.11.0`-Wortlaut aus `modul-03-spec.md` steht nirgends mehr; genau die Bewegung, gegen die `MR-038` drei Absätze lang argumentiert, ist in derselben Datei ausgeführt worden. Die Botschaft desselben Commits sagt dazu *„Nichts entfernt, nichts umgeschrieben"*, während der Diff vier Zeilen löscht (F-2). Kein Gate sieht das — `make gates` ist auf genau diesem Stand grün.

Die zweite Gruppe trifft die **Reichweite**. Der Titel sagt *„schärft `MR-021`"*, aber der Geltungsbereich der Schärfung ist ein echter Teil des geschärften Eintrags: `MR-021` nennt den Reviewer-Skill und die lebenden Planungs-Dokumente ausdrücklich, `MR-038` deckt nur `harness/conventions/` — und in den ausgeschlossenen Dateien stehen heute neun wörtliche Baseline-Zitate (F-3). Der `Geschichte`-Vergleich, der die ganze Pflege-Antwort trägt, schreibt dem Kanon eine Erlaubnis zu, die repo-lokal ist und strukturell einen **ausgenommenen** Abschnitt betrifft, nicht ein Pflichtfeld (F-7). Und die Kanon-Suche, die der Slice sich selbst als Risiko aufgeschrieben hat, ist auf zwei Stellen verengt; drei Nachbarsätze weisen in die Gegenrichtung und sind nicht gewogen (F-8). Das ist `BEO-012` in einem Commit-Strang, der `BEO-012` führt.

Die dritte Gruppe sind **Zahlen und Mengen**, und hier hat sich das Produkt selbst als Prüfer angeboten. Die Zusage *„acht wortgleich (normalisiert gemessen)"* nennt ihre Normalisierung nicht; unter derjenigen, die der Eintrag zum Gegenstand hat, liefern fünf Baseline-Zitate `citation-mismatch` — darunter `MR-038`s eigenes Kanon-Zitat und die gerade reparierte Fassung in `MR-033` (F-4, F-5). Ich habe das nicht behauptet, sondern gefahren: neun Zitate mit echter Direktive gegen ihre Quell-Spanne, fünf rot, vier grün. Und der Zensus zum Blocker nennt sechs Dokumente, während zehn den Lauf brechen; die vier ungenannten sind eine immutable ADR und drei eingefrorene Review-Reporte — genau die, für die der billigere der zwei Wege in `slice-152` gar nicht offensteht (F-6).

Die sechs LOW reisen mit: ein zweiter Bruchpfad, den die Diagnose nicht nennt (F-9), eine Probe ohne Exit und ohne Gegenstand (F-10), eine Ablöse-Feldhälfte in einer Schärfung (F-11), ein gekürztes Lauf-Zitat im Slice über wörtliche Zitate (F-12), ein Blocker, der seit Juli protokolliert ist (F-13), und eine Korrektur im Körper eines `Accepted`-Eintrags, die im selben Zug einen unmarkierten dritten Fall setzt (F-14).
