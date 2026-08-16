# Re-Review-Report: slice-102 — Wellen-Lifecycle-Invariante — 2026-08-16

**Review-Art:** Code — **bestätigende Re-Review**: geprüft wird, ob die vierzehn
Befunde des Erst-Reports
([2026-08-16-slice-102-wellen-invariante-review.md](2026-08-16-slice-102-wellen-invariante-review.md))
wirklich geheilt sind — und ob die Heilung neue Defekte eingeführt hat. Jedes
Urteil stützt sich auf eigene Läufe oder Code-/Vertragszitate, nicht auf den
Commit-Text. Nicht geprüft wird die DoD-Abhakung (getrennter Kontext,
Verifikation).

**Gegenstand:** Commit-Range `9c7506d..d3e3501` (sieben Commits), im Besonderen
der Heilungs-Commit `d3e3501`; Arbeitsbaum-Stand `d3e3501` (= HEAD, clean).

**Skill:** `.harness/skills/reviewer.md` @ 1.4.0 ·
**Modell:** claude-fable-5 · **Datum:** 2026-08-16

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  §Wellen-Invariante (Lastenheft-Historie bis 0.59.1)
- §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritte W1–W5, das §2-Schema (`planning.waves.*`), die vier §4-Zeilen und
  §[`DC-FA-TGT-001.a`](../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
  Schritt 3 (Tabellenzeilen-Lexik, jetzt geteiltes Prädikat `tableRowLine`)
- [ADR-0055](../plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
  (Proposed), [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md),
  [ADR-0028](../plan/adr/0028-planning-lifecycle-modul.md),
  [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md),
  [`AGENTS.md`](../../AGENTS.md) §3, [`CLAUDE.md`](../../CLAUDE.md)
- Der Slice-Plan
  [slice-102](../plan/planning/in-progress/slice-102-wellen-lifecycle-invariante.md)
  (§3, §3a samt F-8-Korrektur) und das Wellendokument
  [welle-75](../plan/planning/welle-75-wellen-register.md)

**Läufe dieses Re-Reviews.** Das Image wurde aus dem HEAD-Stand frisch gebaut;
Baseline `make test` grün (Exit 0). **Sechs Mutationsläufe**, jeder über eine
Dateikopie (Sicherung außerhalb des Repos, mutieren → `make test` → Kopie
byte-identisch zurückschreiben, `diff` leer), Ergebnis am **Exit-Code**
abgelesen, nie per `git checkout`: M1 (dirErrs-Zweig aus), M2
(next-heading-Guard aus), M3n (naive Aktiv-Status-Eigenbestimmung ohne
Fence-Behandlung), M3s (Fence-bewusste Eigenbestimmung mit rohem
Abschnitts-Vergleich), M4 (flach-Hälfte von W5a entfernt), M5
(results-glob-Präfix-Validierung entfernt). Rund **zwanzig Fixture-Läufe**
gegen das HEAD-Image (netzlos, read-only) in einem Scratch-Verzeichnis
außerhalb des Repos; keine Probe-Config wurde im Repo abgelegt. Nachmessung:
`make gates` grün, der `planning-check`-Lauf mit aktivem `waves` meldet **379
Dateien, 0 Befunde**; `make verify-closure-notes` grün (349/0). Die
Schwester-Repo-Messung lief über eine **lesende Kopie** des Planungs-Baums
(am Original nichts geändert; die zunächst versuchte Overlay-Einblendung einer
Probe-Config in den read-only-Mount scheitert am Mountpoint — deshalb Kopie).
`git status --short` ist am Ende leer bis auf diesen Report.

---

## Urteils-Tabelle F-1 … F-14

| Befund | Urteil | Beleg (eigener Lauf / Zitat) |
|---|---|---|
| F-1 (HIGH) unlesbares `waves.dir` still | **geheilt** | Fixture Ruhe-Marker + flache Datei + dir-Tippfehler ⇒ zwei `wave-drift` (dir und abgeleitetes done-dir), Exit 1 — in beiden Zuständen; Mutation M1 (Zweig aus) wird von genau `TestWavesUnlesbaresVerzeichnisFailClosed` rot (Exit 2, FAIL-Zeile nennt den Test); W2 trägt die Zusage jetzt wörtlich. Rest-Beobachtungen: N-1, N-5 |
| F-2 (MEDIUM) fehlende Register-Überschrift still | **geheilt** | Fixture mit umbenanntem Vorschau-Register ⇒ `wave-drift` an Zeile 1, Exit 1; der Erst-Report-Fall (umbenanntes Abschluss-Register) ist getestet (`TestWavesFehlendeRegisterUeberschriftFailClosed`); W4 trägt den Guard. Rest: die next-Hälfte ist unbewacht — Mutation M2 überlebt (N-6) |
| F-3 (MEDIUM) Kennungs-Lexik unterspezifiziert, stiller `welle-`-Rückfall | **geheilt** | Glob mit führendem Platzhalter ⇒ Exit 2 mit benanntem Schlüssel; fremdes results-Präfix ⇒ Exit 2; konsistenter `wave-*`-Konsument läuft mit 0 Befunden durch (Fixture + `TestWavesFremdesPraefix`); W1/W2 sagen die Ableitung jetzt („literales Glob-Präfix plus Ziffernfolge"). Rest: der Lastenheft-Körpersatz nennt weiter wörtlich `welle-<n>` (N-8), Meldungsqualität am Default-Rand (N-7) |
| F-4 (MEDIUM) `wave-unregistered`-target = Kennung | **geheilt** | Fixture „Notiz ohne Registerzeile": target-Spalte zeigt den **Notiz-Pfad** im Ruheort; `TestWavesUnregisteredTargetIstDieNotiz` nagelt es fest; W5c um die Ventil-Paritäts-Formulierung ergänzt |
| F-5 (MEDIUM) geteilte Aktiv-Status-Antwort nicht festgenagelt | **geheilt** | Der im Erst-Report beschriebene Mutant (roher `##`-Vergleich, keine Fence-Behandlung, kein Mehrfach-Guard) nachgebaut (M3n) ⇒ `make test` Exit 2, rot durch genau `TestWavesAktivStatusIstDieGeteilteAntwort`. Rest: die Assertion deckt nur die Fence-Achse (N-5) |
| F-6 (MEDIUM) flach-Hälfte von W5a untestet | **geheilt** | Mutation M4 (Bedingung auf Ruheort verkürzt) ⇒ Exit 2, rot durch genau `TestWavesVorschauMitFlachemDokument` |
| F-7 (MEDIUM) explizit leer fällt still auf Default | **geheilt** | `rawWaves` führt die vier optionalen Schlüssel als Zeiger; Fixtures `next-heading: ""` und `glob: ""` ⇒ Exit 2 mit Klartext („weglassen ⇒ Default"); W1 unterscheidet explizit/abwesend ausdrücklich. Rest: der ganze W1-Config-Rand ist untestet (N-2) |
| F-8 (MEDIUM) zwölfter Schwester-Befund als Bestands-Faktum | **teilweise** | Eigene Nachmessung über die lesende Kopie: konsument-gerechter Marker ⇒ **11** `wave-results-missing`, **kein** `wave-drift`; Default-Marker ⇒ 13 (elf + Artefakt-`wave-drift` + Artefakt-`planning-drift`). Korrigiert: Slice-§3a-Schlussabsatz, ADR-Konsequenzen, CHANGELOG, Lastenheft-Historie. **Nicht korrigiert:** die §3a-Messtabelle behauptet die widerlegte Aussage-1/2-Verletzung weiter, und der Config-Kommentar trägt die un-annotierte Zwölf (N-3) |
| F-9 (MEDIUM) `--print-config` ohne `waves` | **geheilt** | Die emittierte Vorlage des Images führt den `waves`-Block mit allen sechs Schlüsseln, Defaults und fail-closed-Hinweisen (`config_template.go`, live verifiziert) |
| F-10 (MEDIUM) Tabellenzeilen-Erkennung zweifach inline | **geheilt** | `tableRowLine` in `internal/hexagon/core/rules/markdown.go` ist das eine Prädikat; `targets.go` und `planning_waves.go` rufen es beide (strukturell verifiziert); der Docstring erklärt es zur einen Antwort des Produkts |
| F-11 (LOW) fehlende Roadmap meldet doppelt | **geheilt** | Fixture ohne Roadmap ⇒ genau **ein** `planning-drift`; W2 dokumentiert die Schweigen-Disziplin („kein Doppelbefund") — der Randfall ist jetzt vertraglich entschieden |
| F-12 (LOW) „Grund-Codes folgen"-Satz im W5-Körper | **geheilt** | Der Satz ist aus dem Algorithmus-Körper entfernt; er steht nur noch in Historien-Zeilen (dort korrekt, als Zeitpunkt-Aussage — Haus-Konvention der Vorgänger-Fähigkeiten) |
| F-13 (LOW) Roadmap wird zweimal gelesen | **nicht geheilt** (offen) | Beide Fähigkeiten lesen die Datei weiterhin je einmal über den Port und berechnen den Status getrennt (`planning.go` und `planning_waves.go`); der Commit-Text beansprucht die Heilung nicht. Unter [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) weiterhin keine beobachtbare Divergenz — als LOW tragbar |
| F-14 (INFO) Kopf-/Trennzeile nur über den Kennungs-Filter | **nicht geheilt** (offen) | Fixture: eine Abschluss-Register-**Kopfzeile**, deren Spaltentitel eine Kennung trägt, erzeugt weiterhin ein falsches `wave-results-missing`; W4 sagt weiter „ohne Kopf- und Trennzeile" ohne strukturelle Entsprechung. Der Commit-Text beansprucht die Heilung nicht — als INFO tragbar, aber das Vertragswort steht |

## Neue Findings

### N-1 — Der `wave-drift`-Sammelcode kollidiert in der Deduplikation: zwei verschiedene Verletzungen fallen zu einem Befund zusammen

- `kategorie`: MEDIUM
- `quelle`: [ADR-0055](../plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
  Entscheidung 4 (die eigene Begründung: getrennte Codes, **weil** die
  Deduplikation über (Datei, Zeile, Regel, Ziel, Grund) zwei Verletzungen
  derselben Zeile sonst zusammenfallen ließe) ·
  §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W3/W4
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:110` (W3-Drift), `:129`
  und `:143` (Register-Guards) gegen die Deduplikation in
  `internal/hexagon/core/model/finding.go:101`
- `befund`: Die Heilung hat `wave-drift` zum Sammelcode für **drei** Bedeutungen
  gemacht (Status-Widerspruch, unlesbares Verzeichnis, fehlende
  Register-Überschrift) — und zwei davon können auf dasselbe Tupel fallen.
  Konstruiert und differentiell belegt: (a) Roadmap, deren Aktiv-Status-Zeile
  Zeile 1 ist, W3-Drift **und** fehlende Vorschau-Überschrift ⇒ beide Befunde
  tragen (Roadmap, 1, planning, dir, wave-drift) mit verschiedenen Meldungen —
  gemeldet wird **einer** (Kontroll-Fixtures: jede Verletzung allein meldet);
  (b) `done-dir` gleich `dir` konfiguriert und **beide** Register-Überschriften
  fehlen ⇒ ebenfalls ein Befund statt zwei. Es bleibt nie still (Exit 1), aber
  eine der beiden Reparaturen wird verschluckt — exakt der Effekt, den die
  Vier-Codes-Entscheidung ausschließen sollte, jetzt innerhalb des Codes. Die
  Negativ-Zeile des Erst-Reports („die einzige konstruierbare
  Selbe-Zeile-Kollision wird durch die getrennten Codes aufgelöst") gilt seit
  der Heilung nicht mehr.
- `verifizierbar`: ja — die drei Scratch-Fixtures (Kombination plus zwei
  Kontrollen) gegen das Image: 1 Befund, wo einzeln je 1 gemeldet wird.
- `klasse`: „Sammelcode kollidiert in der Befund-Deduplikation"

### N-2 — Der gesamte W1-Config-Rand ist untestet: die Präfix-Validierung lässt sich entfernen, ohne dass ein Test rot wird

- `kategorie`: MEDIUM
- `quelle`: Reviewer-Anker „fehlende Negativtests bei neuem öffentlichen
  Vertrag" ·
  §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W1 (fünf Exit-2-Zusagen)
- `pfad`: `internal/adapter/driven/configyaml/configyaml.go:1079`–`:1153`
  (`applyWaves`/`wavesGlob`/`wavesHeading`) gegen
  `internal/adapter/driven/configyaml/configyaml_test.go` (kein Vorkommen von
  `waves`)
- `befund`: Mutation M5 (die results-glob-Präfix-Prüfung per Kurzschluss
  deaktiviert) überlebt `make test` mit Exit 0; eine Suche über alle Tests
  findet **keinen**, der einen `planning.waves`-Konfigurationsfehler stellt.
  Alle in W1 zugesagten Exit-2-Ränder (Pfad-Flucht, explizit leere/ungültige
  Globs, leere Überschriften, Platzhalter-Präfix, fremdes results-Präfix) —
  darunter die Heilungen von F-3 und F-7 — sind heute nur Behauptungen des
  Codes. Der Haus-Standard steht in derselben Testdatei direkt daneben: die
  Closure-Fähigkeit hat eine tabellengetriebene Negativ-Reihe für jeden ihrer
  Exit-2-Ränder.
- `verifizierbar`: ja — M5 gegen `make test`: Exit 0.
- `klasse`: „Config-Rand-Zusage ohne Negativtest"

### N-3 — F-8-Rest: die §3a-Messtabelle behauptet die widerlegte Aussage-1/2-Verletzung weiter, und der Config-Kommentar trägt die un-annotierte Zwölf

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §3a (Messung) · eigene Nachmessung (lesende Kopie des
  Schwester-Baums: 11 Befunde, **kein** `wave-drift` mit konsument-gerechtem
  Marker)
- `pfad`: [slice-102 §3a](../plan/planning/in-progress/slice-102-wellen-lifecycle-invariante.md)
  Zeile 123 (Tabellenzeile zu Aussage 1+2) · `.d-check.yml:188`
- `befund`: Die Heilung hat den §3a-**Schlussabsatz**, die ADR-Konsequenzen,
  den CHANGELOG und die Lastenheft-Historie korrigiert — die **Messtabelle**
  fünfundzwanzig Zeilen darüber sagt aber weiterhin: das Schwester-Repo habe „0
  bei aktiver Welle — letzteres ist eine Verletzung, die Roadmap nennt eine
  Welle ohne Wellendokument". Das ist exakt die von F-8 widerlegte Aussage
  (die Roadmap nennt keine aktive Welle; sie ruht unter dem Marker
  „**Keine.**"), und sie widerspricht jetzt dem korrigierten Absatz im selben
  Abschnitt. Zusätzlich transportiert der Kommentar der eigenen Konfiguration
  die „12" weiter als Messwert des Schwester-Repos, ohne den Artefakt-Anteil zu
  nennen. Der Commit-Text beansprucht „Slice §3a … korrigiert" — für die
  Tabelle trifft das nicht zu.
- `verifizierbar`: ja — Textbefund plus die beiden Nachmess-Läufe (11 mit
  Marker; 13 inkl. zweier Artefakte ohne).
- `klasse`: „Messwert trägt Konfigurations-Artefakt als Bestands-Befund"
  (Fortschreibung aus F-8)

### N-4 — Der `--doctor`-Klartext von `wave-drift` benennt nur eine der drei Bedeutungen — in den fail-closed-Fällen ist die Diagnose faktisch falsch

- `kategorie`: MEDIUM
- `quelle`: §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  §4-Zeile `wave-drift` (nennt die fail-closed-Fälle) ·
  [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel einer Semantik-Änderung)
- `pfad`: `internal/hexagon/core/app/diagnose.go:131`
- `befund`: Die Heilung hat die §4-Zeile des Codes um die zwei
  fail-closed-Bedeutungen erweitert, den Diagnose-Klartext aber nicht: er
  lautet weiterhin nur „Roadmap-Aktiv-Status und Präsenz eines flachen
  Wellendokuments widersprechen sich". Am Tippfehler-Fixture zeigt `--doctor`
  für das unlesbare Verzeichnis genau diesen Satz — dort widerspricht sich
  nichts, das Verzeichnis ist unlesbar; die Detail-Meldung mit der echten
  Ursache erscheint in der Diagnose-Ansicht nicht. Die Schwester-Codes zeigen
  die verlangte Form (`planning-drift` und `closure-note-missing` tragen ihre
  fail-closed-Klausel im Klartext); die im Erst-Report bestätigte
  Deckungsgleichheit Klartext ↔ §4 ist durch die Heilung gebrochen.
- `verifizierbar`: ja — `--doctor` über das Tippfehler-Fixture: der Klartext
  behauptet einen Widerspruch, die Ursache ist ein unlesbares Verzeichnis.
- `klasse`: „Semantik-Spiegel bei der Heilung ausgelassen"

### N-5 — Die Kopplungs-Assertion deckt nur die Fence-Achse: eine Fence-bewusste Eigen-Bestimmung überlebt

- `kategorie`: LOW
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 4 · Lastenheft-Historie 0.58.3 (die dort geheilten Achsen)
- `pfad`: `internal/hexagon/core/rules/planning_waves_test.go:219`
  (`TestWavesAktivStatusIstDieGeteilteAntwort`)
- `befund`: Mutation M3s — eine lokale Aktiv-Status-Eigenbestimmung, die
  Fences korrekt behandelt, das Abschnitts-Ende aber am rohen `##`-Präfix
  bestimmt und den Mehrfach-Überschrift-Guard weglässt — überlebt `make test`
  (Exit 0). Die geforderte eine Assertion je Konsument existiert und tötet den
  naiven Mutanten; sie legt den Status aber in keine Lage, in der die übrigen
  Divergenz-Achsen tragen (eingerückte/tab-getrennte H2, H1 beendet H2,
  doppelte Überschrift) — exakt die Achsen, die in der 0.58.3-Runde als reale
  Defekte geheilt wurden. Kein aktueller Fehler; eine Wartungsfalle für den
  Tag, an dem jemand die Kopplung „vereinfacht".
- `verifizierbar`: ja — M3s gegen `make test`: Exit 0 (M3n dagegen Exit 2).
- `klasse`: „geteilte Antwort mit Assertion auf nur einer Achse"

### N-6 — Die next-Hälfte des Register-Guards ist unbewacht: der Guard lässt sich entfernen, ohne dass ein Test rot wird

- `kategorie`: LOW
- `quelle`: Reviewer-Anker „fehlende Negativtests bei neuem öffentlichen
  Vertrag" ·
  §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W4 („Fehlt **eine der beiden** Register-Überschriften")
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:128` gegen
  `internal/hexagon/core/rules/planning_waves_test.go:198` (Test nur für die
  closed-Hälfte)
- `befund`: Mutation M2 (den `nextNo == 0`-Zweig per Kurzschluss deaktiviert)
  überlebt `make test` mit Exit 0 — der einzige Guard-Test stellt das
  Abschluss-Register um, nie die Vorschau. Dieselbe Halbform, die F-6 bei W5a
  war („Vertrags-Oder ohne Test für eine Hälfte"), reproduziert von der
  Heilung, die F-6 geschlossen hat. Verhalten und Vertrag stimmen heute
  (Fixture belegt den Guard); nur die Bewachung fehlt.
- `verifizierbar`: ja — M2 gegen `make test`: Exit 0.
- `klasse`: „Vertrags-Und ohne Test für eine Hälfte"

### N-7 — Die Präfix-Meldung beschuldigt einen nie gesetzten Schlüssel und zitiert das Präfix statt des Werts

- `kategorie`: LOW
- `quelle`: Maintainability (Meldungsqualität am Config-Rand); Vergleichsmaßstab
  die übrigen `waves`-Meldungen („weglassen ⇒ Default", benannter Schlüssel,
  zitierter Wert)
- `pfad`: `internal/adapter/driven/configyaml/configyaml.go:1108`–`:1110`
- `befund`: Setzt ein Konsument nur `glob` (etwa auf ein `wave-*`-Schema) und
  lässt `results-glob` weg, bricht der Lauf korrekt fail-closed ab — die
  Meldung lautet aber, `results-glob` „welle-" trage nicht das Präfix: sie
  beschuldigt einen Schlüssel, den der Konsument nie geschrieben hat, ohne zu
  sagen, dass dessen **Default** gemeint ist, und zitiert in beiden
  Präfix-Meldungen `rp` — das abgeleitete Präfix — an der Stelle, an der die
  Meldungs-Grammatik den Glob-Wert verspricht. Der nächste Schritt des
  Konsumenten (welchen Wert er wo korrigieren soll) ist aus der Meldung nicht
  ablesbar. Die Abbruch-Richtung selbst ist richtig und vertragsgedeckt (W1,
  effektive Werte).
- `verifizierbar`: ja — Fixture „glob gesetzt, results-glob abwesend": Exit 2
  mit der beschriebenen Meldung.
- `klasse`: „Config-Rand-Meldung nennt Ableitung statt Eingabe"

### N-8 — Der Lastenheft-Körpersatz nennt das Kennungs-Schema weiter wörtlich `welle-<n>`

- `kategorie`: LOW
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  §Wellen-Invariante gegen
  §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W2 (generalisiert) und die eigene 0.59.1-Historien-Zeile
- `pfad`: `spec/lastenheft.md:2135`–`2136`
- `befund`: Die Heilung hat die Kennung in der Spezifikation ans literale
  Glob-Präfix gebunden und sagt das auch in der Lastenheft-**Historie**; der
  **Körpersatz** der Anforderung sagt aber weiterhin „Verglichen wird über das
  Zahlen-Präfix `welle-<n>`" — ohne Default-Qualifikation. Für einen
  `wave-*`-Konsumenten (seit der Heilung ausdrücklich getragen und getestet)
  beschreibt der abnahmebindende Text die falsche Kennung; da das Lastenheft
  in der Source Precedence über der Spezifikation steht, gewinnt formal der
  veraltete Satz. Der Rest der F-3-Vertragshälfte.
- `verifizierbar`: ja — Textbefund; W2 und die 0.59.1-Zeile sagen das
  Gegenteil desselben Dokuments.
- `klasse`: „Rand auf der alten Fassung stehengeblieben"

## Negativbefunde

- geprüft, ohne Befund: **F-1-Heilung in beiden Zuständen und die
  dirErrs-Deduplizierung.** Ruhe- und Aktiv-Zustand melden beim Tippfehler-dir
  fail-closed (Fixtures, Exit 1); `dir` gleich `done-dir` und unlesbar ⇒ genau
  **ein** Befund (der Gleichheits-Guard in `waveSets` dedupliziert korrekt);
  die Meldung nennt Verzeichnis und Ursache statt der falschen
  Ursachen-Nennung aus dem Erst-Report.
- geprüft, ohne Befund: **Config-Rand, Verhaltens-Seite — alle vier
  Glob-Kombinationen.** Beide gesetzt/konsistent ⇒ läuft (0 Befunde am
  `wave-*`-Baum); beide abwesend ⇒ Defaults laufen; je eines gesetzt mit
  Präfix-Konflikt ⇒ Exit 2 in **beiden** Richtungen — es rutscht kein
  inkonsistentes Paar durch, und der Abbruch des glob-only-Falls ist die
  sichere Richtung (mit dem Default-results-glob liefe die Zuordnung real
  auseinander). Explizit leere Schlüssel: Exit 2 für `glob` und
  `next-heading`; `done-dir` leer fällt dokumentiert auf `<dir>/done` (W1
  sagt es so — keine stille Abweichung).
- geprüft, ohne Befund: **Geteiltes Prädikat und targets-Parität.**
  `tableRowLine` ist das einzige Tabellenzeilen-Prädikat (ein gemeinsames
  Symbol, beide Aufrufer verifiziert); die `targets`-Tests der Baseline
  bleiben grün; das Fence-Verhalten des neuen Konsumenten ist weiter getestet.
- geprüft, ohne Befund: **W5-Befundformen nach der Heilung.**
  `wave-preview-exists` an der Zeile mit der Kennung, `wave-results-missing`
  an der Zeile, `wave-unregistered` an der `closed-heading`-Zeile mit dem
  Notiz-Pfad — alle drei vertragskonform (Fixtures).
- geprüft, ohne Befund: **Eigene Nachmessung.** `make gates` grün; der
  `planning-check`-Lauf mit aktivem `waves` meldet 379 Dateien, 0 Befunde;
  `make verify-closure-notes` grün (349/0) — beide Zahlen decken sich mit dem
  Commit-Text.
- geprüft, ohne Befund: **Schwester-Repo-Zahlen.** Elf `wave-results-missing`
  mit konsument-gerechtem Marker („**Keine.**"), kein `wave-drift`; mit
  Default-Marker reproduzierbar 13 (elf plus zwei Marker-Artefakte) — die
  korrigierten Aussagen in ADR-Konsequenzen, CHANGELOG und Lastenheft-Historie
  stimmen mit der eigenen Messung überein (Rest: N-3).
- geprüft, ohne Befund: **ADR-Immutabilität der Heilung.**
  [ADR-0055](../plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
  ist Proposed — die inhaltliche Korrektur der Konsequenzen war zulässig, kein
  Supersede nötig; Entscheidungen 1–5 sind unverändert und decken die
  Implementierung weiterhin.
- geprüft, ohne Befund: **Hexagon-Schnitt und Hermetik nach dem Umbau.**
  `planning_waves.go` importiert weiter nur Kern-Model und den
  Filesystem-Port; die neuen Zweige lesen keine neuen Eingaben (konform
  [ADR-0005](../plan/adr/0005-modul-layout-hexagon-ordner.md) /
  [ADR-0028](../plan/adr/0028-planning-lifecycle-modul.md)); Determinismus
  ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)) über die
  bestehende Sortierung unberührt, Wiederholungsläufe identisch.
- geprüft, ohne Befund: **`--print-config`-Spiegel**
  ([`DC-FA-CLI-005`](../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)):
  die Vorlage nennt alle sechs Schlüssel mit den §2-Defaults;
  `--suggest-config` bleibt konsistent ohne `waves` (wie ohne `closure` —
  unverändert, dort kein Spiegel-Anspruch).
- geprüft, ohne Befund: **F-11-Gegenprobe.** Der Doppelbefund ist weg, und die
  Schweigen-Richtung ist jetzt vertraglich gedeckt (W2) — die im Erst-Report
  monierte Uneinheitlichkeit ist als Entscheidung dokumentiert statt offen.

## Summary

| Kategorie | Anzahl (neue Findings) |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 4 |
| INFO | 0 |

**Urteils-Bilanz der vierzehn Erst-Befunde:** elf **geheilt** (F-1 bis F-7 und
F-9 bis F-12), einer **teilweise** (F-8), zwei **nicht geheilt und vom
Commit nicht beansprucht** (F-13 LOW, F-14 INFO — beide tragbar, F-14 lässt
das W4-Vertragswort ohne Entsprechung stehen).

**Wiederholungs-Signal.** Zwei der neuen Findings sind Wiedergänger von
Klassen, die die Heilung selbst geschlossen hat: N-6 ist F-6 in der
next-Hälfte („Hälfte ohne Test"), N-3 ist F-8 auf der übersehenen Fläche. Und
N-1 ist die Vier-Codes-Begründung von
[ADR-0055](../plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
gegen die eigene Heilung gewendet: wer fail-closed-Fälle in einen Sammelcode
legt, erzeugt die Kollision, die die Entscheidung 4 ausschließen wollte.

## Verdikt

**Merge-blockierend:** ja — vier MEDIUM (kein HIGH).

**Was die Heilung trägt:** Der blockierende Kern des Erst-Reports ist real
geschlossen und eigenständig belegt: alle drei Still-Abschalt-Pfade melden
fail-closed (Fixtures, Exit 1), die vier Config-Rand-Härtungen greifen (Exit 2
mit benanntem Schlüssel), die beiden überlebenden Mutanten des Erst-Reports
sind tot (M1/M3n/M4 werden von genau ihrem Test gefangen), das geteilte
Prädikat existiert als ein Symbol, die Befund-Form der Gegenrichtung ist
Ventil-paritätisch, und die korrigierte Messaussage (elf robust, die Zwölf ein
Marker-Artefakt) hält der unabhängigen Replikation stand. Eigener Baum 0,
Gates grün.

**Was bleibt:** N-1 ist eine Produkt-Semantik-Frage und gehört **vor das
Release** entschieden — die Grund-Codes werden mit v0.59 öffentlich, danach
ist die Korrektur des Sammelcodes ein Vertragsbruch statt ein Nachzug (eigene
Codes für die fail-closed-Fälle oder kollisionsfreie Tupel; die Entscheidung
gehört in den Slice, nicht hierher). N-3 und N-4 sind kleine Züge (zwei
Textstellen, ein Klartext), N-2 eine Negativ-Testtabelle nach vorhandenem
Haus-Muster ohne Produktänderung. N-5 bis N-8 sind tragbar und ohne
Release-Risiko nachziehbar; F-13/F-14 können erklärt offen bleiben, F-14
saubererweise mit einer Anpassung des W4-Wortlauts oder des Codes.

**Release-Empfehlung: noch nicht releasen** — ein Nachzug im selben Slice
(N-1 bis N-4), dann steht dem Minor nichts entgegen. Die Minor-Einordnung
selbst (opt-in, ohne die neuen Schlüssel byte-identisch) ist unverändert
korrekt und durch den Modul-aus-Lauf gedeckt.

**Übergabe:** Die Findings gehen an die Implementation; N-1/N-6/N-3 gehören
als Klassen-Fortschreibung in Slice-Closure und Beobachtungs-Register. Dieser
Report ist ein Lauf-Beleg (dieser Diff, dieser Skill, dieses Modell, dieses
Verdikt) und ersetzt keine Verifikation.
