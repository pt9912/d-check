# Review-Report: slice-225 — 2026-09-17 (R2)

**Review-Art:** Code-/Plan-Review — geprüft wird die Behebung der Findings aus
Review-Runde 1 (`docs/reviews/2026-09-17-slice-225-open-tasks-marker-review-r1.md`):
die neue `structure`-Bedingung `open-tasks-require-marker-section`
(`internal/hexagon/core/rules/structure.go`, `internal/hexagon/core/model/config.go`,
`internal/adapter/driven/configyaml/configyaml.go`), die `ADR-0085`-Geschichte-
Ergänzung, die Spec-Bumps (`spec/lastenheft.md` 0.87.1, `spec/spezifikation.md`),
die neue Beobachtungs-Registerzeile
(`BEO-ALL/messmethode-scope-enger-als-dokumentierte-ziel-form`) und die
Fortschreibung von `slice-225`s §6/§7/§8 — gegen den R1-Report, `AGENTS.md`
§3.5/§3.7/§5, `harness/conventions.md` und Baseline `v6.9.0`
`modul-05-planning-harness.md`/`modul-06-roadmap.md`.

**Gegenstand:** `docs/plan/planning/in-progress/slice-225-gegenstand-entfallen-uebernommen.md`,
Commit-Kette `35c783f5..8036f9c5` (`35c783f5` R1-Report, `c9384a27`
R1-F-1-Fix, `8036f9c5` Risiko-Ausgänge/Closure-Notiz/Register).

**Skill:** `.harness/skills/reviewer.md` @ v1.16.0 / 2026-09-07 ·
**Modell:** `claude-sonnet-5` · **Datum:** 2026-09-17

> **Zitier-Form.** Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb Kennung statt Adresse: `slice-NNN` statt Lifecycle-Pfad,
> `make <target>` als Token statt Link auf die Sensor-Datei, eine
> Baseline-Stelle als Tag plus Pfad in Inline-Code
> (`v6.9.0` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext:**

- R1-Report `docs/reviews/2026-09-17-slice-225-open-tasks-marker-review-r1.md`
  (1 HIGH, 2 MEDIUM, 1 LOW; blockiert)
- Slice-Plan `slice-225` (Fassung `in-progress/`, Stand nach `8036f9c5`)
- `ADR-0085` (Accepted, mit neuem `## Geschichte`-Eintrag 2026-09-17)
- `spec/lastenheft.md` 0.87.1, `spec/spezifikation.md` (Schritt-6-/Schema-
  Tabellen-Nachzug)
- eingehender CR `docs/plan/cr/2026-09-17-cr-eingehend-ai-harness-init-stilllegungs-form.md`
  samt Antwort (unverändert seit `da7205bc`, **nicht** Teil der Fix-Commits)
- `AGENTS.md` §3.5 (ADR-Immutabilität), §3.7 (Kommentar-Klassen), §5
  (Zähl-/Zitat-Disziplin, Botschafts-Reichweite)
- Baseline `v6.9.0` · `modul-05-planning-harness.md` §Offene Risiken werden
  bei Closure aufgelöst; `modul-06-roadmap.md` §Das Beobachtungs-Register
  (*„Drei Quellen speisen den Closure-Eintrag"*)
- Beobachtungs-Register: `BEO-ALL/commit-message-overclaims-work` (Bestand
  vor diesem Slice, zwölf Einträge), `BEO-ALL/messmethode-scope-enger-als-dokumentierte-ziel-form`
  (neu, dieser Slice), `BEO-ALL/citation-stretched-beyond-scope`

---

## Verifikationen dieses Laufs

- `make test`: grün, alle Pakete inkl. `internal/hexagon/core/rules` und
  `internal/adapter/driven/configyaml`.
- `make lint`: **0 issues.**
- `make gates`: grün — `baseline-verify + workflow-pins + doc-check + lint +
  test + arch-check + coverage-gate + semgrep + gate-consistency +
  planning-check`; Coverage 94,60 % (Schwelle 93 %, `structureOpenTasks`/
  `markerBody`/`hasMarker` je 100,0 %); Semgrep 0 Findings (55 Regeln, 63
  Dateien).
- `make verify-closure-notes`: **691 Datei(en) geprüft, 0 Befund(e)** —
  bestätigt, dass `slice-221`s Doppel-Marke weiterhin grün durchläuft
  (Bestandsschutz) und kein anderer `done/`-Slice durch die neue
  `open-tasks-require-marker-section`-Zeile rot wird.
- **Eigene empirische Gegenprobe gegen das lokal gebaute Image**
  (`d-check:latest`, `sha256:b78c960b…`), unabhängig von R1s Fixture neu
  aufgesetzt, mit der exakten R1-F-1-Konstellation (Marke nur in
  „## 7. Closure-Notiz", offene Haken in „## 2. Definition of Done"):

  | Config | Ergebnis |
  |---|---|
  | **ohne** `open-tasks-require-marker-section` (Negativ-Kontrolle) | `section-open-tasks-marker-missing` auf der Überschriftszeile — **derselbe Fehlschlag wie in R1** |
  | **mit** `open-tasks-require-marker-section: '^#{1,3} [0-9]+\. Closure-Notiz'` | `1 Datei(en) geprüft, 0 Befund(e)`, Exit 0 |

  **R1-F-1 ist damit doppelt bestätigt behoben**: einmal über die eigene
  Fixture gegen das Image (diese Probe), einmal über die neue Unit-Test-Suite
  (`TestOpenTasksRequireMarkerSection_BaselineZielFormWirdErkannt` reproduziert
  denselben Vorzustand/Fix-Kontrast inline).
- Code-Lesung `structureOpenTasks`/`markerBody`
  (`internal/hexagon/core/rules/structure.go:380-420`): Reihenfolge korrekt
  (Zählung zuerst, Marken-Prüfung nur bei Überschuss, `markerBody` nur bei
  gesetztem `OpenTasksRequireMarker`), `FindSectionHeads`/`SectionProse`
  korrekt wiederverwendet (dieselbe Erkennung wie jede andere
  Abschnitts-Suche des Moduls), Regex-Validität bereits am Config-Rand
  geprüft (`structureBedingungsFehler`) — `regexp.MustCompile` in `markerBody`
  kann laufzeitseitig nicht auf einen ungültigen Pattern treffen. Mehrere
  Treffer werden korrekt vereinigt (`strings.Builder`), kein Treffer liefert
  korrekt den leeren String (⇒ Marke gilt als fehlend).

## Findings

### R2-F-1 (MEDIUM) — Der Fix-Commit für R1-F-2 wiederholt R1-F-2s Fehlerklasse

- **Kategorie:** MEDIUM
- **Quelle:** `AGENTS.md` §5 (*„eine genannte Probe muss gelaufen sein … wer
  N Formen geprüft hat, berichtet N"*); `BEO-ALL/commit-message-overclaims-work`
- **Pfad:** Commit `c9384a27` (Botschaft, dritter Absatz); `docs/plan/adr/0085-bedingte-pflicht-marke-open-tasks.md`
  (Geschichte-Zeile 2026-09-17); `docs/plan/planning/in-progress/slice-225-gegenstand-entfallen-uebernommen.md`
  §7 (*„durch vier neue Tests (`TestOpenTasksRequireMarkerSection_*`) …
  belegt"*)
- **Befund:** Alle drei Stellen behaupten „vier neue Tests"
  (`TestOpenTasksRequireMarkerSection_*`), die den Vorzustand reproduzieren
  und den Fix belegen. `internal/hexagon/core/rules/structure_offene_tasks_test.go`
  trägt nach diesem Commit genau **drei** neue Funktionen mit diesem Präfix
  (`_BaselineZielFormWirdErkannt`, `_FehlenderAbschnittGiltAlsFehlend`,
  `_AndereNummerierungTrifftTrotzdem` — Zeilen 274, 299, 309), nicht vier.
  Das ist exakt dieselbe Fehlerklasse wie R1-F-2 (Commit `b2874d18` zählte
  „fünf" statt vier `TestOpenTasksRequireMarker_*`-Funktionen) — hier tritt
  sie ein zweites Mal auf, **in dem Commit, der R1-F-1 behebt**, und wird
  von dort unverändert in die ADR-Geschichte und die Closure-Notiz
  kopiert, statt beim zweiten Schreiben aufzufallen. Die begleitende
  `spec/lastenheft.md`-Historie-Zeile (0.87.1, „Vier neue
  **Akzeptanzkriterien**") zählt dagegen korrekt — sie zählt eine andere
  Menge (Given/Then-Kriterien, von denen tatsächlich vier existieren,
  eines davon mit zwei Given/Then-Sätzen in einem Bullet), was die
  Verwechslung erklärt, sie aber nicht auf die Testzahl überträgt.
- **Verifizierbar:** ja — `grep -c '^func TestOpenTasksRequireMarkerSection_' internal/hexagon/core/rules/structure_offene_tasks_test.go` liefert 3.
- **Klasse:** `commit-message-overclaims-work` (Zähl-Variante) — zweites
  Auftreten in diesem Slice.

### R2-F-2 (MEDIUM) — Die wiederkehrende Finding-Klasse speist das bestehende Register nicht

- **Kategorie:** MEDIUM
- **Quelle:** Baseline `v6.9.0` · `regelwerk/modul-06-roadmap.md`
  §Das Beobachtungs-Register (*„Drei Quellen speisen den Closure-Eintrag:
  eigene Beobachtung · offenes Risiko aus dem Slice-Plan · wiederkehrende
  Finding-Klasse aus dem Review dieses Slice"*)
- **Pfad:** `docs/plan/planning/observations/BEO-ALL/commit-message-overclaims-work/evidence/`
  (kein `slice-225.md`); `docs/plan/planning/in-progress/slice-225-gegenstand-entfallen-uebernommen.md`
  §7 (nennt R1-F-2 nur in Prosa, ohne Register-Bezug)
- **Befund:** R1-F-2 ist wörtlich eine „wiederkehrende Finding-Klasse aus
  dem Review dieses Slice" — `BEO-ALL/commit-message-overclaims-work`
  existiert bereits (zwölf Evidence-Dateien, Stand „verkörpert", aber
  laut eigenem `state.md` weiterhin wach: *„Bleibt trotz Verkörperung
  offen, weil eine geschriebene Regel weiter verfehlt werden kann"*).
  §7 der Closure-Notiz hält R1-F-2 korrekt in Prosa fest (siehe R1-F-2s
  eigene Bewertung unten), trägt aber **keine** neue
  `evidence/slice-225.md` in diesem Register nach — anders als für R1-F-1s
  Klasse, für die korrekt ein **neues** Register-Verzeichnis
  (`messmethode-scope-enger-als-dokumentierte-ziel-form`) angelegt wurde.
  Mit R2-F-1 (oben) liegen jetzt **zwei** unabhängige Instanzen derselben
  Klasse in diesem einen Slice vor, von denen keine im bestehenden
  Register verzeichnet ist.
- **Verifizierbar:** nein (Register-Paarung ist ein Urteil über
  Vollständigkeit, kein Gate — `make verify-closure-notes` prüft nur, dass
  zitierte Pfade auflösen, nicht dass jede einschlägige Review-Finding-
  Klasse verzeichnet wurde).
- **Klasse:** `register-quelle-review-finding-uebergangen`

### R2-F-3 (MEDIUM) — Die CR-Antwort an den Adopter bleibt hinter dem R1-Fix zurück

- **Kategorie:** MEDIUM
- **Quelle:** R1-Report, Verdikt (*„… und, falls zutreffend, in der
  CR-Antwort an `ai-harness-init` (die Baseline-Wortlaut zitiert, den der
  Code so nicht einlöst)"*); `AGENTS.md` §5 (Zitier-/Botschafts-Disziplin)
- **Pfad:** `docs/plan/cr/2026-09-17-cr-eingehend-ai-harness-init-stilllegungs-form.md`
  §Antwort, Absatz „Beleg" (unverändert seit Commit `da7205bc`)
- **Befund:** `git log` bestätigt: Die CR-Datei wurde von keinem der beiden
  Fix-Commits (`c9384a27`, `8036f9c5`) berührt. Ihr „Beleg"-Absatz nennt
  weiterhin nur `open-tasks-require-marker` (0.87.0-Stand) und schweigt zu
  `open-tasks-require-marker-section`. Der CR selbst zitiert exakt die
  Baseline-Stelle, die den Fehler verursachte (*„§7 trägt … eine Zeile
  `Gegenstand:`"*, ein **eigener** Abschnitt) — und die Baseline-Ziel-Form,
  die der CR damit referenziert, trennt DoD- und Closure-Notiz-Abschnitt
  strukturell (Modul 5, Ziel-Form: Slice). Ein Adopter, der
  `open-tasks-require-marker` exakt nach dem dokumentierten Stand der
  Antwort einführt (ohne die neue `-section`-Variante zu kennen, weil sie
  dort nicht erwähnt wird), reproduziert `R1-F-1`s Fehler in seinem
  eigenen Repo, sobald sein Slice-Template — wie das hiesige und wie die
  von ihm selbst zitierte Baseline — DoD und Closure-Notiz in getrennten
  Abschnitten führt. R1s Verdikt hatte die CR-Antwort ausdrücklich als
  einen der zu prüfenden Orte benannt („falls zutreffend"); dieser Fall
  ist zutreffend, wurde aber nicht abgearbeitet.
- **Verifizierbar:** ja — `git log --oneline -- docs/plan/cr/2026-09-17-cr-eingehend-ai-harness-init-stilllegungs-form.md`
  zeigt nur `d8e30b7d`/`da7205bc`, keinen der beiden Fix-Commits.
- **Klasse:** `citation-stretched-beyond-scope` (Geschwister-Variante: die
  zitierte Antwort bleibt hinter ihrem eigenen späteren Korrektur-Stand
  zurück, nicht hinter einer fremden Quelle)

## Negativbefunde

- **R1-F-1 (HIGH) behoben:** geprüft — empirisch gegen das gebaute Image
  (eigene, von R1 unabhängig aufgesetzte Fixture) **und** über die neue
  Test-Suite; Negativ-Kontrolle (ohne den neuen Schlüssel) reproduziert den
  R1-Fehlschlag unverändert. Kein Befund.
- **R1-F-3 (MEDIUM) korrekt nachgeschärft:** die in §6 dritte Risikozeile
  genannte Paketliste (`hexagon/core/rules`, `hexagon/core/model`,
  `hexagon/core/app`, `adapter/driven/configyaml`, `adapter/driving/cli`)
  stimmt exakt mit `git show b2874d18 --stat` überein (drei Core-Pakete,
  ein Driven-, ein Driving-Adapter). Kein Befund.
- **R1-F-4 (LOW) korrekt berichtigt:** §8 nennt jetzt „40 Verzeichnisse" für
  den **gemergten Stand zum Sichtungszeitpunkt** (vor diesem Slices eigener
  Registerzeile); `git ls-tree -r c9384a27~1` (Vorzustand vor der neuen
  Registerzeile) bestätigt exakt 40 Beobachtungs-Verzeichnisse. Kein Befund.
- **R1-F-2 (MEDIUM) ehrlich behandelt:** §7 hält fest, dass die
  Commit-Botschaft von `b2874d18` „unkorrigierbar am Commit selbst" ist und
  die Diskrepanz „hier festgehalten, damit sie nicht verschwindet" —
  beschönigt nichts, verschweigt nichts. (Dass dieselbe Klasse in der
  Fix-Runde selbst wiederkehrte, ist R2-F-1, kein Vorwurf gegen diese
  Behandlung von R1-F-2 selbst.) Kein Befund an dieser Stelle.
- **ADR-Immutabilität** (`AGENTS.md` §3.5): geprüft — der Diff an
  `docs/plan/adr/0085-bedingte-pflicht-marke-open-tasks.md` in `c9384a27`
  besteht aus **genau einer** neuen Zeile in der `## Geschichte`-Tabelle;
  `Kontext`, `Entscheidung`, `Konsequenzen`, `Fitness Function` und
  `Re-Evaluierungs-Trigger` sind byte-identisch zum Vorzustand, der
  `**Status:**`-Wert bleibt `Accepted`. Kein Befund.
- **Kommentar-Klassen** (`AGENTS.md` §3.7): geprüft in `structure.go`,
  `model/config.go`, `configyaml.go`. Die neuen Kommentare (`markerBody`,
  `structureMarkenFehler`, das erweiterte `structureOpenTasks`-Doc, das
  `OpenTasksRequireMarkerSection`-Feld) tragen Zusage/Kopplung/Grenze-
  Anker. Eigens geprüft: die Wendung „der urspruengliche Entwurf traf nur
  den Gleichschnitt-Fall" (`model/config.go:588`) — sie beschreibt eine
  **Grenze** des Abwesend-Zustands (nicht Review-Prozess-Historie) und
  entspricht demselben, bereits im Bestand etablierten Stil wie
  „Genau daran scheiterte die Vorgaenger-Form als Closure-Vorbedingung"
  (`structure.go`, unveränderter Bestand, Zeile ~347). Kein Befund.
- **Hexagon-Import-Richtung** (`ADR-0005`): geprüft — keine neuen Imports,
  `regexp`/`strings` bereits vorhandene Abhängigkeiten des Pakets. Kein
  Befund.
- **Halbe Aktivierung / Config-Ränder** (ADR-0085 Entscheidung 5,
  Erweiterung): geprüft — `structureMarkenFehler` deckt beide Stufen
  (`open-tasks-require-marker` ohne `max-open-tasks`,
  `open-tasks-require-marker-section` ohne `open-tasks-require-marker`),
  beide mit Bruch-Test in `configyaml_test.go`. Kein Befund.
- **Referenzrichtung/Provenance-Marker** (`AGENTS.md` §3.4): geprüft — alle
  neuen `slice-221`/`slice-225`-Nennungen in ADR-Geschichte und
  Spec-Historie tragen `<!-- d-check:status-provenance -->`; `make
  doc-check` (Teil von `make gates`, 0 Befunde) bestätigt zusätzlich
  keine unmarkierte Abwärts-Referenz. Kein Befund.
- **Gocyclo/Lint der neuen Verzweigung** (`structureUeberschriftFehler` →
  `structureMarkenFehler`-Auslagerung): `make lint` meldet 0 issues; die
  Auslagerung ist explizit mit dem Zweck begründet („damit
  `structureUeberschriftFehler` unter der gocyclo-Schwelle bleibt"). Kein
  Befund.
- **`Sub-Area`/Register-Form der neuen Beobachtung**
  (`BEO-ALL/messmethode-scope-enger-als-dokumentierte-ziel-form`): geprüft
  gegen Baseline `v6.9.0` · `modul-06-roadmap.md` §Das Beobachtungs-Register
  — Pfadschema (`BEO-<KUERZEL>/<slug>`, Kürzel `ALL` aus der
  Modus-Deklaration für Sub-Area `*`), Drei-Datei-Form
  (`observation.md`/`state.md`/`evidence/<vorgang>.md`), `state.md` trägt
  `offen` (korrekt für ein Erstauftreten unter der 3×-Schwelle), Titel-Form
  und `**Sub-Area:**`-Feld entsprechen exakt bestehenden Nachbar-Einträgen
  (`citation-stretched-beyond-scope`). Kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 0 |
| INFO | 0 |

## Verdikt

**R1-F-1 (HIGH) ist behoben** — doppelt empirisch bestätigt (eigene
Image-Gegenprobe plus Test-Suite), keine neue Ausweich-Lücke gefunden. R1-F-2
bis R1-F-4 sind sachlich korrekt und ehrlich eingearbeitet.

**Drei neue MEDIUM-Befunde**, die vor der finalen Closure zu klären sind
(MEDIUM blockiert nach Skill-Konvention typischerweise, hier mit
geringerem Gewicht als R1s HIGH — keiner erfordert eine weitere
Code-Änderung an `structure.go` selbst):

- R2-F-1: Zähl-Wiederholung in Commit-Botschaft/ADR-Geschichte/
  Closure-Notiz (vier statt drei Tests) — Korrektur ist eine Text-Präzisierung,
  keine Code-Änderung.
- R2-F-2: die zweite Instanz derselben Zähl-Fehlerklasse in diesem Slice
  sollte `BEO-ALL/commit-message-overclaims-work` als Evidence zugeführt
  werden (Modul 6: „wiederkehrende Finding-Klasse aus dem Review" ist eine
  der drei Pflichtquellen für den Closure-Eintrag).
- R2-F-3: die CR-Antwort an `ai-harness-init` sollte vor Closure um einen
  Hinweis auf `open-tasks-require-marker-section` ergänzt werden, sonst
  trägt dieses Repo die Korrektur, während der Adopter, der die Lücke
  ursprünglich meldete, weiter auf dem Stand steht, der sie hatte.

**Kein neues HIGH.** Keiner der drei Befunde erzwingt eine dritte
Review-Runde am Code; sie sind Dokumentations-/Register-Nachträge, die der
Implementer vor dem `git mv` nach `done/` einarbeiten sollte.
