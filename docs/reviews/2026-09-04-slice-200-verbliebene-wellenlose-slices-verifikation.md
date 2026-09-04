# Verifikation: slice-200 — Verbliebene wellenlose Slices archivieren

**Art:** Verifikation
**Gegenstand:** Commit `ce5fb50ed7207c39d86dcbd8725ea9c7d611e07e` — „feat(planning):
slice-200 -- sieben verbliebene wellenlose Slices archiviert (slice-200)“
**Modell-ID:** claude-sonnet-5
**Datum:** 2026-09-04
**Rolle:** unabhängiger Verifier (gegen DoD/Spec, nicht gegen Plan/ADR)

## Vorgehen

Slice-Plan vollständig gelesen
(`docs/plan/planning/in-progress/slice-200-verbliebene-wellenlose-slices-archivieren.md`),
jeder DoD-Punkt aus §4 unabhängig gegen den tatsächlichen Repo-Zustand
geprüft (eigene Kommandos, keine übernommenen Behauptungen), §5-Risiken
gegen den Ausgang bewertet, Commit-Struktur und Stub-Inhalt stichprobenartig
gegengelesen.

## DoD-Punkt für Punkt

### 1. „Bestand exakt erhoben (Zahl im DoD eingetragen)“

`git show ce5fb50 --stat` zählt genau **7** gelöschte
`done/slice-*.md`-Dateien: `slice-141-fixture-racily-clean.md`,
`slice-168-adr-index-titelspalte.md`, `slice-169-workflow-rechte-pruefung.md`,
`slice-170-workflows-modul.md`, `slice-183-baseline-v5150.md`,
`slice-184-uses-tag-kohaerenz.md`, `slice-188-register-gegen-neuen-kanon.md`.
Namen und Anzahl stimmen exakt mit der Planungs-Liste (§2 Punkt 1) überein.
**Bestätigt.** — Anmerkung: Der DoD-Haken selbst ist im Plan-Kopf noch
**nicht** gesetzt (Datei liegt noch in `in-progress/`); das ist zu diesem
Zeitpunkt der Lifecycle-Sequenz korrekt (Body-Fill folgt vor dem `git mv`
nach `done/`, siehe Präzedenz unten).

### 2. „Alle erhobenen Slices archiviert (Stub + Zip je Slice)“

Für alle 7 Kennungen geprüft: `docs/plan/planning/done/wellenlos/<basisname>.md`
(Stub) **und** `docs/plan/planning/done/wellenlos/slice-<NNN>-archiv.zip`
existieren beide. **Bestätigt**, keine Lücke.

Stub-Form stichprobenartig gegen ältere Stubs (slice-083, slice-095,
slice-102 aus slice-197) verglichen: identische Form (`ARCHIVIERT`-Kopf,
`**Welle:**`, `**Archiviert:** <manuell auszufuellen>` als
werkzeug-generierter Platzhalter — kein Defekt, sondern durchgängiges
Tool-Verhalten, siehe Vorlage `archiv-stub-slice.template.md`).

### 3. „`done/` enthält außerhalb `done/wellenlos/` keinen flachen wellenlosen Slice mehr“

`ls docs/plan/planning/done/*.md | grep -v 'results\.md'` liefert **keine**
Zeile (Exit 1 / leer). Der flache `done/`-Bestand besteht ausschließlich noch
aus den 31 `welle-*-results.md`. **Bestätigt.**

### 4. „`make gates` grün (zehn Gates) auf dem Endstand“

Selbst ausgeführt: `make gates` läuft vollständig durch, Schlusszeile:

```
[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green
```

Alle zehn Gates grün, Coverage 94,70 % über Schwelle 93 %, Semgrep 0 Findings.
**Bestätigt.**

### 5. „`make fullbuild` grün (`--require-complete`: 0 neue Trace-Waisen)“

**Teilweise rot, mit klarem, selbstauflösendem Grund.** `make fullbuild`
läuft bis `image-test`, `bench` und der RTM (`--trace --require-complete`,
„51 Anforderung(en), 0 Waise(n)“) grün durch — die
`--require-complete`-Teilaussage des DoD-Punkts ist damit erfüllt. Danach
bricht `verify-closure-notes` mit 7 Befunden ab, alle der Form
`closure-note-missing` / `review-missing` / `section-missing`, jeweils mit
der Meldung „das Gate liefe leer“ bzw. „leere Pruefmenge“.

**Ursache verifiziert:** `.d-check.closure.yml` bindet `closure.dir` auf
`docs/plan/planning/done` **nicht-rekursiv** mit implizitem
`slice-*.md`-Glob. Nach der Archivierung aller sieben verbliebenen Slices
(DoD-Punkt 3) und da slice-200 selbst noch nicht als flache Datei in `done/`
liegt (es liegt noch in `in-progress/`), trifft dieser Glob auf **null**
Dateien — das Modul ist laut eigener Design-Zusage fail-closed bei leerer
Kandidatenmenge (siehe `make review-coverage`-Beschreibung in AGENTS.md §4:
„Fail-closed bei leerer Kandidatenmenge […], nicht bei null gefundenen
Zusagen“ — dieselbe Bindung gilt strukturanalog für `structure`/`planning`
im Closure-Profil).

**Reproduktion der Selbstauflösung, nicht-destruktiv getestet:** die
aktuelle (unvollständige) `slice-200`-Plandatei probeweise nach
`docs/plan/planning/done/slice-200-….md` kopiert (nicht committet) und
`make verify-closure-notes` erneut gelaufen. Ergebnis: die
leere-Kandidatenmenge-Befunde (`closure-note-missing`,
`review-missing` auf Verzeichnisebene, `section-missing`) **verschwinden
vollständig** — stattdessen meldet das Werkzeug jetzt die *echten*,
inhaltlichen Befunde des unfertigen Standes (`section-tasks-open` für die
noch offenen DoD-Haken, `closure-note-thin` für die noch leere
Closure-Notiz, `review-missing` für den fehlenden Report). Das bestätigt
präzise den erwarteten Mechanismus: Sobald **irgendeine** flache
`slice-*.md`-Datei in `done/` liegt, ist die Kandidatenmenge nicht mehr leer
und das Gate bewertet echten Inhalt statt eines leeren Scans. Testdatei
danach entfernt, Arbeitsbaum wieder sauber (`git status --short` leer
bestätigt).

**Lifecycle-Präzedenz geprüft:** anhand von slice-199 (`417d3ab` →
`620d00a`) zeigt sich die tatsächliche Reihenfolge: **erst** wird der
Slice-Body vollständig ausgefüllt (DoD-Haken gesetzt, Closure-Notiz
geschrieben, Risiko-Ausgänge eingetragen) — **während die Datei noch in
`in-progress/` liegt** —, **danach** erst der reine `git mv` nach `done/`
(Diff zeigt „| 0" — die Datei selbst bleibt im Move-Commit inhaltlich
unverändert, Rename-Detection hält). Übertragen auf slice-200: Wenn dieselbe
Reihenfolge gefahren wird (Body zuerst, dann `git mv`), landet die Datei
bereits **vollständig ausgefüllt** in `done/` — genau der Zustand, den mein
Test simuliert hat, nur mit erfüllten statt offenen Haken. Ein vollständig
ausgefüllter slice-200 mit gesetzten DoD-Haken, ausgefüllter Closure-Notiz
und einem Review-Report unter `docs/reviews/` sollte die verbleibenden
inhaltlichen Prüfungen (`section-tasks-open`, `closure-note-thin`,
`review-missing`) bestehen, da die Struktur (§4 „Definition of Done“, §9
„Closure-Notiz“, §5 „Abnahme-Punkte / Risiken“) exakt dem Template folgt,
das die Muster erwarten.

**Einschätzung: kein Blocker, sondern eine Momentaufnahme.** Der aktuelle
rote `make fullbuild`-Lauf ist eine Eigenschaft des *jetzigen*
Zwischenzustands (7 archiviert, slice-200 selbst noch nicht verschoben) —
nicht eine strukturelle Verletzung, die die Schließung von slice-200
verhindert. Im Gegenteil: **die Schließung von slice-200 selbst ist der Akt,
der die Kandidatenmenge wieder befüllt**, weil der etablierte
Lifecycle-Ablauf verlangt, den Body auszufüllen, **bevor** die Datei nach
`done/` wandert. Sobald das geschieht, sollte ein erneuter
`make fullbuild`-Lauf (auf dem dann wahren Endstand, mit slice-200 selbst
flach in `done/`) wieder grün sein. Empfehlung: den DoD-Haken für diesen
Punkt erst nach diesem finalen Lauf setzen und in der Closure-Notiz (§9)
ausdrücklich vermerken, dass der Zwischenstand (7 archiviert, slice-200 noch
`in-progress/`) einen `closure-note-missing`/`section-missing`-Fail-closed-
Befund durch eine leere Kandidatenmenge zeigte — als Beleg, keine
Behauptung eines durchgehend grünen Laufs auf einem Stand, auf dem er es
nicht war.

### 6. „Unabhängiger Review durchgeführt, Report unter `docs/reviews/`“

Außerhalb meines Auftrags (Review läuft parallel); nicht durch mich
bestätigt.

### 7. „Unabhängige Verifikation durchgeführt“

Dieser Report ist der Beleg.

### 8. „Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang“

Noch nicht geschrieben (Plan liegt noch in `in-progress/`, §9 trägt noch den
Platzhalter-Kommentar). Nicht Gegenstand dieser Verifikation vor der
Body-Füllung.

## §5-Risiken gegen den Ausgang

- **„Stille Bedeutungsverschiebung“** — Voll-Dry-Run wurde laut Commit-Botschaft
  vor der Anwendung gefahren; die 7 Ergebnisse stimmen exakt mit der
  Vorab-Erhebung überein. Kein Hinweis auf Verschiebung. **Ausgang: entfallen.**
- **„Ein Slice könnte einzige Zitierstelle sein“** — `--require-complete`
  meldet 0 Waisen (51/51 Anforderungen `ok`). **Ausgang: entfallen.**
- **„Umfang sprengt die Ein-Sitzungs-Grenze“** — sieben mechanische,
  werkzeug-verifizierte Einzel-Läufe plus ein Bündel-Commit; keine
  Rückführung nötig. **Ausgang: entfallen.**

Alle drei Risiken sind, wie bei slice-197/199, als „entfallen“ auflösbar.

## Commit-Struktur (kurzer Eindruck, Review prüft das parallel vertieft)

`git diff-tree --name-status -M` zeigt für die 7 Slice-Archivierungen reine
`D`/`A`-Paare (keine erkannten Renames) plus drei reine Referenz-Fixes
(`AGENTS.md`, `welle-87-results.md`, `MR-057-baseline-v5150.md` — alle drei
korrekt auf `wellenlos/`-Pfade nachgezogen, von `doc-check`/`links` bereits
grün bestätigt). Das entspricht MR-062 für **jeden einzelnen**
`-slice=<id>`-Move. **Offener Punkt, nicht durch MR-062 gedeckt:** MR-062
benennt in seinem eigenen §Grenze ausdrücklich, dass es **nur** begründet,
dass je Slice ein Commit die Untergrenze ist — **nicht**, ob mehrere
unabhängige Einzel-Slice-Archive in einem gemeinsamen Commit gebündelt
werden dürfen. Diese Bündelung (7 in einem Commit) ist dieselbe Praxis wie
bei slice-197/199, aber laut MR-062-Wortlaut eine offene, nicht
abschließend beantwortete Frage — kein klarer Regelverstoß, aber auch keine
ausdrückliche Deckung. Vertiefte Bewertung obliegt dem parallelen Review.

## `**Welle:**`-Feld / Pfad-Tiefen-Risiko (Punkt 5 des Auftrags)

Alle 7 neuen Stubs geprüft: Jedes `**Welle:**`-Feld trägt reine Prosa
(„— **wellenlos**, …“), **keinen** relativen Link auf eine Wellen-Datei.
Die bekannte Fehlerklasse (`../../` statt Geschwister-Pfad bei zwei Dateien
im selben `wellenlos/`-Verzeichnis) ist hier **strukturell nicht
einschlägig**, da kein Link existiert, der eine falsche Tiefe tragen
könnte. **Bestätigt: kein Befund.**

## Fazit

- DoD erfüllt: **Noch nicht vollständig** — Punkte 1–3 belegt korrekt,
  Punkt 4 (`make gates`) grün bestätigt, Punkt 5 (`make fullbuild`) zeigt
  einen roten Zwischenstand, der jedoch nachvollziehbar und reproduzierbar
  als **selbstauflösende Momentaufnahme** identifiziert wurde (leere
  Kandidatenmenge des nicht-rekursiven `done/slice-*.md`-Globs, solange
  slice-200 selbst noch nicht flach in `done/` liegt), **kein struktureller
  Blocker**. Punkte 6–8 sind zum jetzigen Zeitpunkt (Datei noch in
  `in-progress/`) plangemäß offen.
- Schließbar: **Ja, mit einer Bedingung.** Slice-200 ist schließbar, sobald
  der etablierte Lifecycle-Ablauf eingehalten wird: Body füllen (DoD-Haken,
  Closure-Notiz, Review-Report) **während** die Datei noch in
  `in-progress/` liegt, **dann** `make fullbuild` auf dem dann wahren
  Endstand (mit slice-200 selbst als vollständig ausgefüllte, aber noch
  nicht verschobene Datei) erneut prüfen — dieser Lauf sollte grün sein, da
  mein nicht-destruktiver Test bestätigt hat, dass eine anwesende
  `done/slice-*.md`-Datei die Fail-closed-Befunde durch echte
  Struktur-Prüfung ersetzt. Erst danach der reine `git mv` nach `done/`.
  Die Commit-Bündelung (7 Archivierungen + 3 Referenz-Fixes in einem
  Commit) ist Praxis-konform zu slice-197/199, aber von MR-062 nicht
  ausdrücklich gedeckt — dies ist eine offene Frage, kein Blocker im engeren
  DoD-Sinn.
