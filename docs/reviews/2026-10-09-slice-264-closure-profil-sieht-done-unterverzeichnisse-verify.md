# Verifikation — slice-264: Das Closure-Profil prüft die Slices in den Unterverzeichnissen von `done/`

- **Rolle:** Verifier (Modul 11). Geprüft werden DoD, Plan und Spec gegen den Stand. Den Diff gegen Plan und ADR prüft der Reviewer.
- **Gegenstand:** `docs/plan/planning/in-progress/slice-264-closure-profil-sieht-done-unterverzeichnisse.md`, DoD-Punkte 1–3 und die Plan-Änderungen in §3
- **Range:** `fe54da7e~1..HEAD` (11 Commits, Spitze `37bdb24e`)
- **Eingang:** die Commits; die Reports R1–R3 unter `docs/reviews/`;
  [`SPEC-095`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge);
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
  [`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in);
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md),
  [ADR-0081](../plan/adr/0081-reviews-modul.md),
  [ADR-0105](../plan/adr/0105-reviews-liest-done-unterverzeichnisse.md) (Proposed);
  [`MR-049`](../../harness/conventions.md#mr-049)
- **Datum:** 2026-10-09

## Messungen (Kommando und Ergebnis)

Alle Läufe habe ich selbst gefahren, nur über make/Docker und Lesewerkzeuge. Die
Brech-Proben liefen in einem frischen `git clone` von `HEAD` im Scratchpad.
Gemessen wurde jeweils mit
`docker run --rm --network none -v <klon>:/repo:ro d-check:latest --config <profil> …`.
Das Image war mit `make gates` aus `HEAD` gebaut. Der Slice ändert keinen Go-Code,
deshalb ist dasselbe Image für das alte und das neue Profil gültig. Als
**altes Profil** dienten `git show c55bae02~1:.d-check.closure.yml` und
`git show f851c64f~1:.d-check.yml`, beide als ungetrackte Kopien im Klon. Nach
jeder Probe wurde der Klon per `git checkout HEAD -- …` zurückgesetzt. Am
Arbeitsbaum des Repos habe ich nichts geändert außer diesem Report.

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `make gates` | Exit 0. `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; lint `0 issues.`; alle Testpakete `ok`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; semgrep `Ran 55 rules on 66 files: 0 findings.`; `fetch-baseline-cache: verify ok (54 Dateien, vollständig)` |
| M2 | `make verify-closure-notes` | Exit 0, `981 Datei(en) geprüft, 0 Befund(e)` |
| M3 | `make review-coverage` | Exit 0, `1117 Datei(en) geprüft, 0 Befund(e)` |
| M4 | `make adr-check RANGE=fe54da7e~1..HEAD`, `make trace-check RANGE=fe54da7e~1..HEAD` | beide Exit 0, je 11 Commits, 0 Befunde. Die Zeile, die ADR-0081 unter `## Geschichte` neu bekommt, besteht das Immutable-Gate |
| M5 | Klon, Inventur `done/`: Marker `^> \*\*ARCHIVIERT\*\* — Volltext:` | 254 Dateien tragen den Marker (alle in Unterverzeichnissen, keine flach), davon 229 `slice-*.md`. Unter `wellenlos/` liegen 80 Stubs und 24 Volltexte. Eine Datei nennt `ARCHIVIERT`, ohne die Marker-Form zu tragen: slice-263, ein Volltext, der den Marker inline zitiert. Sie bleibt Kandidatin. Keine Symlinks, keine `welle-*-results.md` in Unterverzeichnissen |
| M6 | Klon, unveränderter Stand, altes und neues Closure-Profil, `--enable planning --enable structure --enable spans --enable reviews` | beide Exit 0, je 0 Befunde |
| **M7** | Klon, **realer Vorzustand**: slice-240 bis slice-243 per `git checkout c55bae02~1 -- …` zurückgesetzt | **neu:** Exit 1, **4 Befunde**, je `section-forbidden` auf `done/**/slice-*.md :: ^## [0-9]+\. (?:Abnahme-Punkte / Risiken\|Risiken und offene Punkte)$`, `verbotenes Muster trifft: \*\*Ausgang:…`, für slice-240:109, -241:93, -242:90 und -243:95. **alt:** Exit 0, 0 Befunde |
| **M8** | Klon, `planning`-Hälfte: §7 des Volltexts `wellenlos/slice-270-…` auf `Fertig.` gekürzt | **neu:** Exit 1, `closure-note-thin` (`trägt 1 Satzende-Zeichen …, verlangt sind 4`) und `closure-note-boilerplate` (`Floskel „fertig“`), Ziel `docs/plan/planning/done`. **alt:** Exit 0, 0 Befunde |
| M9 | Klon, wie M8, aber die Überschrift `## 7. Closure-Notiz` in `## 7. Notizen` umbenannt | **neu:** Exit 1, `closure-note-missing` (planning) und `section-missing` auf `^#{2,3} .*Closure-Notiz` (structure). **alt:** Exit 0 |
| **M10** | Klon, Stub-Ausnahme: Marker in `wellenlos/slice-141-…` und `welle-90/slice-199-…` auf `> Archiv:` verändert | **neu:** Exit 1, 9 Befunde. slice-199: `closure-note-missing` und viermal `section-missing` (Closure, DoD, Risiko, DoD-Haken). slice-141: `closure-note-missing` und dreimal `section-missing`; die DoD-Haken-Regel fehlt dort, weil die Ausnahme `done/**/slice-1[0-6]?-*` in einem Unterverzeichnis greift. **alt:** Exit 0. Mit Marker bleiben die Stubs still (M6, 254 Stubs, 0 Befunde) |
| **M11** | Klon, `reviews`-Hälfte: alle sechs Reports `docs/reviews/*slice-270-*` gelöscht; slice-270 trägt `- [x] Review durchgeführt …` (Z. 58) | Closure-Profil **neu:** Exit 1, `review-missing … für slice-270` (`wellenlos/slice-270-…:58`); **alt:** Exit 0. Hauptprofil in der Form von `make review-coverage`: **neu** Exit 1, derselbe Befund; **alt** (`f851c64f~1`) Exit 0. Damit ist die Fitness Function von ADR-0105 erfüllt |
| M12 | `git diff fe54da7e~1 HEAD -- spec/spezifikation.md harness/sensors/{hooks,verify-closure-notes,review-coverage}.md .githooks/pre-commit` gelesen; `grep` nach Resten („nur die Slices direkt", „nicht rekursiv", „liest der Lauf nicht", „nur `done/` selbst") über `harness/`, `spec/`, `.githooks/`, `.github/`, beide Profile, `AGENTS.md`, `Makefile`, `docs/user/`, `docs/maintainer/` | Die Aussagen sind zurückgenommen (siehe DoD 3). Treffer gibt es nur noch an drei Stellen: in Historie-Zeilen, bei `reviews.reviews-dir` (das Report-Verzeichnis, weiterhin flach, korrekt) und in `docs/user/benutzerhandbuch.md` Z. 1432 (V-3) |
| M13 | `grep -n -i "rekursiv\|unmittelbar\|unterverzeichnis"` über ADR-0048, -0049, -0059, -0077, -0082 | keine Treffer. Außer ADR-0081 Entscheidung 4 trifft keine ADR eine Lage-Entscheidung, der die Rekursion widerspricht |

## DoD-Punkte

**DoD 1 — `.d-check.closure.yml` trifft die Volltexte in den Unterverzeichnissen und nimmt Stubs an ihrem Inhalt aus; bewusst gebrochen.** **Erfüllt.**
Das Profil setzt `planning.closure.recursive: true` und
`skip-pattern: '(?m)^> \*\*ARCHIVIERT\*\* — Volltext:'`. Alle sieben
Slice-Regeln von `structure` stehen auf `done/**/slice-*.md` und tragen dasselbe
`skip-pattern`. Für `reviews` gelten `recursive` und `skip-pattern`.

Die Brech-Proben decken alle drei Hälften aus dem richtigen Grund ab. In jedem
Fall bleibt die alte Fassung grün, und die neue wird mit dem Grund-Code der
zuständigen Regel rot:

| Hälfte | Probe | Grund-Code |
|---|---|---|
| `planning` | M8, M9 | `closure-note-thin` / `-boilerplate` / `-missing` |
| `structure` | M7, der reale Vorzustand | `section-forbidden` über den Ausgangs-Wortschatz von MR-049 |
| `reviews` | M11 | `review-missing` |

Zur Stub-Hälfte: Ein Stub mit Marker bleibt still (M6), ohne Marker wird er
rot (M10), und das gilt in `wellenlos/` ebenso wie in einem Wellen-Verzeichnis.
Das enge Muster nimmt den Volltext von slice-263 nicht aus, obwohl er
`ARCHIVIERT` zitiert (M5).

Die `**/`-Ausnahmen von `structure` greifen in der Tiefe. Das zeigt M10: Die
DoD-Haken-Regel feuert für slice-141 nicht. Dass das Matching dort `**` über
beliebig viele Segmente auflöst (`matchGlob`), habe ich am Code nachgelesen.
Die GRENZE bei `reviews.exempt-paths` (blankes `path.Match`) ist richtig
benannt. Im Profil ist sie heute ohne Gegenstand, weil der Block keine
`exempt-paths` trägt.

**DoD 2 — Altverstöße behoben; `make verify-closure-notes` und `make gates` grün.** **Erfüllt.**
M7 zeigt genau die vier Altverstöße, die der Plan in §8 gemessen hat. Auf `HEAD`
sind sie verschwunden (M2). `make gates` ist grün (M1).

Auch die `planning`-Hälfte meldet über die 24 Volltexte nichts (M6). Das ist
zugleich der Beleg für den Ausgang des Plan-Risikos in §6 („Weitere
Altverstöße"): Es ist **nicht eingetreten**, die `planning`-Hälfte ist auf dem
Bestand leer. Für die Closure-Notiz bietet sich deshalb *entfallen* mit
Verweis auf M6 an.

**DoD 3 — Zurücknahme der „nur `done/` direkt"-Aussagen.** **Erfüllt.**

- [`SPEC-095`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge): Der Satz zum ausgelösten Lauf sagt jetzt „prüft jeden Volltext unter `done/` samt Unterverzeichnissen …". Er nennt die beiden stillen Ausfälle (Marker-Zitat, Verzeichnis-Symlink) mit Verweis auf C2. Zwei Historie-Zeilen sind ergänzt.
- `harness/sensors/hooks.md`: Der Vertrag schließt die Unterverzeichnisse ein. Grenze 4 ist ersetzt; sie benennt jetzt nur noch, dass ein Stub zwar auslöst, aber nicht geprüft wird.
- `harness/sensors/verify-closure-notes.md`: Der Vertrag sagt „im Volltext in `done/` oder einem seiner Unterverzeichnisse". Grenze 8 ist ersetzt und trägt Marker-Zitat, Symlink und Wellen-Ergebnisnotizen nur flach; die letzte Grenze habe ich gegen den Bestand geprüft (M5). Der Bindungs-Satz „liest nur `done/`" ist gestrichen.
- `.githooks/pre-commit`: Der GRENZE-Kommentar ist durch eine Aussage über den Stub ersetzt.

Den Grep über die übrigen Träger zeigt M12. In `.github/workflows/ci.yml` und
`Makefile` steht keine Gegenaussage.

**Plan-Änderungen aus §3.** **Erfüllt.**

- Das Hauptprofil `.d-check.yml` (`reviews`) setzt `recursive` und `skip-pattern` in derselben Form wie das Closure-Profil. Beide Kommentare nennen die KOPPLUNG, und M11 zeigt, dass beide Läufe gleich antworten. Die Bestands-Ausnahme ist entfernt („keine"); die Grenze 4 alt in `review-coverage.md` ist gestrichen.
- [ADR-0105](../plan/adr/0105-reviews-liest-done-unterverzeichnisse.md) steht auf `Proposed`. Sie trägt `**Supersedes:**` auf ADR-0081 (nur Entscheidung 4), drei Alternativen plus die gewählte in der Form Option/Pro/Contra, Fitness Function (durch M11 bestätigt) und Re-Evaluierungs-Trigger. Die Index-Zeile ist ergänzt, ebenso die Geschichte-Zeile an ADR-0081; `adr-check` ist grün (M4).
- Den Register-Eintrag `BEO-ALL/aufnahme-kriterium-ohne-bestands-beleg` gibt es mit `observation.md`, `state.md` (`offen`) und `evidence/slice-242.md`. Der Dateiname der Evidence ist der Vorgang slice-242, wie der Kanon es verlangt (Form: Kennung eines abgeschlossenen Vorgangs).

## Nachgetragene Risiko-Ausgänge slice-240 bis slice-243

Geprüft habe ich je Risiko den Ausgang gegen §7 derselben Datei und gegen den
genannten Beleg.

| Slice | Risiko | Ausgang | Beleg gelesen | trägt? |
|---|---|---|---|---|
| 240 | Über-Hebung lebender Referenz in Frozen-Datei | entfallen | R1-Report slice-240, Negativbefund 3: „durchgehend zeilenweise getrennt, kein Über-Swap auf frozen Aussagen". §7 nennt das Risiko nicht | ja |
| 240 | `d-check:cite`-Spannen verschieben sich | eingetreten, im Slice behoben | §7: „7 statt 0 neu geankert, Wortlaut identisch" | ja (V-1) |
| 240 | frozen `target-missing` auf `v6.9.0` | eingetreten, im Slice behoben | §7: „4 Markdown-Links, Ventil gewachsen"; R1 Negativbefund 4 | ja (V-1) |
| 240 | Ersetzung trifft MR-Vorgänger-Zeile | eingetreten, revertiert | §7: „R1 … eingetreten, revertet" | ja (V-1) |
| 241 | Doppel-Dokumentation README/Skill | entfallen | §7: „vermieden, README gewinnt" | ja |
| 241 | Audit bleibt Prosa ohne Vollzugs-Beleg | entfallen | §7: Audit-Vollzug über vier Klassen | ja |
| 241 | Terminologie bootstrap-aware Gate | entfallen | §7: „n.a. begründet (… kalibrierte Konstanten ohne Hochschalt-Trigger)" | ja |
| 242 | versehentliche Lockerung des Senkungs-Verbots | entfallen | §7: Gegenprobe byte-identisch, R1 Negativbefund 1 | ja |
| 242 | „unabhängig lauffähig" unbelegt | weiter offen, mit Register-Eintrag | §7: „bewusst offener Punkt, Ausgang bei der nächsten Aufnahme notieren" | ja |
| 243 | vorzeitige Mechanisierung | entfallen | §7: „abgewehrt, der Prosa-Verbleib argumentiert gegen die Schwelle" | ja |
| 243 | Zustands- vs. Urteils-Prüfung | entfallen | §7: „Bewertung je Klasse", je Zeile getrennt | ja |

Alle elf Ausgänge sind durch ihre Belege getragen. Die Herkunfts-Zeile
„nachgetragen mit slice-264, aus §7" stimmt für 241 bis 243. Für 240 nennt sie
den R1-Report ausdrücklich mit.

In slice-240 nummeriert §7 die Risiken (R1–R4) anders als §6. Das §7-„R4
Zitat-Delta" steht nicht in §6, und die Über-Hebung aus §6 fehlt in §7. Die
Zuordnung im Nachtrag folgt deshalb dem Inhalt, nicht der Nummer, und sie ist
richtig.

## Befunde

**V-1 — INFO — „eingetreten" ohne Carveout und ohne Folge-Slice-Kennung.**
Für *eingetreten* nennt der Kanon zwei Ziele: Carveout oder Folge-Slice mit
Kennung (`modul-05-planning-harness.md` §Offene Risiken). Die drei Ausgänge in
slice-240 tragen keines davon, sondern „im selben Slice behoben/revertiert".

Die urteilsfreie Hälfte (eines der drei Wörter, MR-049) ist erfüllt, und das
Muster ist Hausbestand: slice-231 schreibt ausdrücklich „*geschlossen im Slice*,
nicht mit einer Folge-Slice-Kennung belegt", slice-236 ebenso. Gezählt habe ich
das Muster „eingetreten — im" dreimal in `done/`.

Eine DoD-Verletzung dieses Slice ist das nicht. Es ist aber eine
**undeklarierte Auslegung** des Kanons. Ein Ort dafür wäre `MR-049` oder ein
eigener Nachtrag. Ich melde den Punkt als Beobachtung und blockiere damit nicht.

**V-2 — LOW — Die Annahme von ADR-0105 steht in keinem DoD-Punkt.**
Plan §3 sagt „`Proposed` bis zur Closure". R3 F-14 (ADR-0081 nennt „abgelöst"
schon als Tatsache) und F-16 (Status und Index von ADR-0081 bei `Accepted`)
sind offen. Die DoD trägt den Übergang `Proposed → Accepted` samt Index-Zeile
nicht. DoD-Punkt 5 nennt nur Closure-Notiz, Register, Risiken und Paarungen.

Ohne einen ausdrücklichen Punkt kann der Slice nach `done/` wandern, während
die ADR, auf der die Konfiguration steht, noch `Proposed` ist. Empfehlung: den
Übergang in der Closure-Checkliste abhaken, also in DoD 5 oder einem eigenen
Punkt.

**V-3 — INFO — Benutzerhandbuch Z. 1432.**
Dort steht weiterhin „Beide Module listen ihr Verzeichnis **nicht rekursiv**"
für `workflows` und `reviews`. Für `reviews` mit `recursive: true` trifft das
nicht mehr zu, und gerade dann zählt die `path.Match`-Grenze bei `exempt-paths`.

Die Aussage gehört zum Produkt (slice-263, Release-Prep nach `AGENTS.md` §5
Regel 17). R1 F-8 hat sie als zurückgestellt akzeptiert. Sie liegt nicht im
DoD-Umfang dieses Slice; ich halte sie hier nur fest, damit die Release-Prep
sie nicht verliert.

## Negativbefunde

- Der Diff ändert keinen Produkt-Code; der Slice bleibt in seiner Schicht (Konfiguration, Doku, Bestand).
- Die Abgrenzung §1 ist eingehalten: Die 24 Volltexte unter `wellenlos/` sind nicht archiviert, und die Produkt-Schlüssel sind nicht berührt.
- Den Umfang erweitert nur, was der Plan in §3 vor dem Code festgehalten hat (R1 F-1, F-2; R2 F-9). Das ist nach `AGENTS.md` §6 Schritt 4 die zulässige Form.
- Keine `Accepted`-ADR ist inhaltlich überschrieben (M4).
- Ein Vergleich, der gegen den alten Profilstand nur zufällig grün bleibt, kommt nicht vor: In jeder Probe sind die Dateizahlen von altem und neuem Lauf gleich (981 bzw. 977/1113), nur die Befunde unterscheiden sich.

## Verdikt

**DoD 1–3 erfüllt, die Plan-Änderungen aus §3 erfüllt, die Risiko-Ausgänge in
slice-240 bis slice-243 getragen.** Das bewusste Brechen ist im frischen Klon
für die `planning`-, `structure`- und `reviews`-Hälfte belegt (M7–M11). Jede
Hälfte war vor `c55bae02` grün und ist mit dem Fix aus dem richtigen Grund rot.

Kein blockierender Befund. Vor dem `git mv` nach `done/` ist V-2 zu erledigen
(ADR-0105 annehmen samt R3 F-14/F-16). Der Ausgang des §6-Risikos kann sich
auf M6 stützen.
