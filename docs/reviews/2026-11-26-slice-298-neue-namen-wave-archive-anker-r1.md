# Review slice-298 — Neue Namen: `wave-<NN>`, `archive.zip`, `seit wave-<NN>`; der Bestand bleibt

- **Review-Art:** Code/Doku (gegen Plan samt den zwei Plan-Änderungen, Architect-Verdikt, `ADR-0117`/`ADR-0055`, `MR-050`/`MR-069`/`MR-070`/`MR-081`/`MR-082`, Hard Rules)
- **Gegenstand:** `slice-298`, Diff `db329ff2..HEAD` (Plan-Änderung `73f52639`; Werkzeug `928cb73f`; reiner Move `9d295355`; Plan-Änderung `21d6f711`; Verweis-Nachzug und Konfiguration `2afe2d47`; `MR-082`, Commands, Briefing, Sensor-Dateien `423bc2d5`)
- **Skill:** `reviewer.md` @ 1.20.0
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-11-26
- **Eingangs-Kontext:** `slice-298` (§1 Abgrenzung, §3 Plan samt beiden Plan-Änderungen, §6 Risiken), `wave-98` (§3, §5), das Architect-Verdikt zu `wave-98` (Fragen 1 bis 3, 8, 12, Nachträge V1 und V2), `ADR-0117` (`Proposed`), `ADR-0055`, `MR-050`, `MR-069`, `MR-070`, `MR-081`, `MR-082`, die Vorlage `MR-NNN-title.template.md` des vendorten `v7.0.0`, `AGENTS.md` §3 und §5, die Vorbild-Reports zu `slice-297` und `slice-295`. Nicht erhalten: die DoD-Abhakung (Sache des Verifiers).
- **Zusätzlich gefahren (Wegwerf-Klone im Scratchpad, nichts im Repo und nichts im Arbeitsbaum):** (1) Bestands-Dry-Run mit dem Image des Vorzustands `db329ff2` und dem Image von `HEAD` auf demselben Baum, je `welle-97`, `welle-96`, `welle-95`; (2) neun Mutationen am Werkzeug (Tabelle unten) und eine Gesamt-Rücknahme der Änderung; (3) eine Regex-Probe mit vierzehn Formen; (4) Planning-Proben (`wave-drift` mit und ohne Zeiger, unregistrierte Bestands-Welle, entfernte Bestands-Registerzeile, entfernte Bestands-Ergebnisnotiz, `wave-results-missing`); (5) `matrix`-Proben je Stratum und `ids`-Proben für `done/wave-*/`; (6) `verify-closure-notes`-Proben für die Closure-Regel; (7) `doc-check` an den Ständen `9d295355` und `21d6f711`.

## Findings

Keine HIGH, kein MEDIUM. Sechs LOW, vier INFO.

### F-1 — LOW — Der Kommentar zur Closure-Regel trägt eine „heute“-Zahl, die mit der Closure dieser Welle falsch wird

- **quelle:** `AGENTS.md` §3.7 (ein Kommentar beschreibt, was da ist), dieser Skill (LOW: latente Wartungsfalle, hart verdrahteter Wert, der erst beim künftigen Edit zündet)
- **pfad:** `.d-check.closure.yml` · „heute sind es genau die welle-<NN>-Notizen“
- **befund:** Der GRENZE-Kommentar der neuen Regel `docs/plan/planning/done/*-results.md` sagt, heute seien es genau die Bestandsnotizen, und nennt das Kommando mit dem Ergebnis 0. Die Closure von `wave-98` legt die erste `wave-98-results.md` an; das Kommando ergibt dann 1, und die Zeile beschreibt nicht mehr, was da ist. Die Grenze selbst („jede `*-results.md` unter `done/`“) bleibt wahr.
- **verifizierbar:** ja — `git ls-files | grep -E 'done/[^/]*-results.md$' | grep -vc 'done/welle-'` ergibt am Stand `HEAD` 0 und nach der Closure 1
- **klasse:** kommentar-zaehlt-heute-statt-grenze

### F-2 — LOW — Derselbe Kommentar nennt die verworfene Alternative („ein Glob je Schreibweise“)

- **quelle:** `AGENTS.md` §3.7 (keine Deliberation über Verworfenes; Auflösung nur Kopplung, Abgrenzung, Grenze); dieser Skill, Anker zur Kommentar-Klasse (HIGH, hier auf LOW gestellt, weil der Satz zugleich die tragende Kopplung an `SPEC-049` nennt)
- **pfad:** `.d-check.closure.yml` · „Ein Glob je Schreibweise meldet section-missing, sobald einer keine Datei trifft“
- **befund:** Der Kommentar erklärt, warum die Regel ein Glob ist, indem er beschreibt, was zwei Regeln täten. Die Kopplung (eine Regel ohne Kandidaten ist fail-closed, `SPEC-049`) ist richtig und gemessen (`section-missing`, `Regel trifft keine Datei`); sie steht dort im Gewand der Alternativen-Erzählung.
- **verifizierbar:** nein — Urteil, kein Gate fängt Kommentar-Klassen
- **klasse:** kommentar-erzaehlt-verworfene-alternative

### F-3 — LOW — Zwei Test-Kommentare nennen eine Rot-Bedingung, die der Test nicht trägt

- **quelle:** `AGENTS.md` §5 Regel 13 (eine Grenze gegen den Gegenstand prüfen, nicht gegen ihre Beschreibung), `slice-298` §6 Risiko 2 (jeder neue Test rot aus dem richtigen Grund)
- **pfad:** `tools/archive-wave/wave_names_test.go` · „Rot, wenn runWelle `wave-` nicht als Kennung findet (Fail-Closed“ und im Dateikopf „Jeder Test nennt die Aenderung, die ihn rot faerbt.“
- **befund:** Mutation M1 (`collect.go` auf den Stand `db329ff2` zurück): `TestRun_PlanOhneEigeneSlices_Wave` bleibt grün, weil `FindWellePlan` den Plan über das Präfix findet und der fail-closed-Zweig („weder Welle-Plan noch Slices“) nie erreicht wird; rot wird er erst über `archive.zip` (M3). Der Dateikopf sagt zudem „jeder Test“, und `TestArchiveName` sowie `TestRunSlice_AcceptsContextualWaveMention` tragen keine Rot-Bedingung.
- **verifizierbar:** ja — Wegwerf-Klon, `git show db329ff2:tools/archive-wave/collect.go > collect.go`, `make test`: `TestRun_PlanOhneEigeneSlices_Wave` fehlt in der Liste der roten Tests (rot sind `TestCollectSlices_WaveUndWelleGleicheZiffernGetrennt`, `TestFixture_EndToEnd_Wave`, `TestFixture_Gemischt_WelleUndWaveImSelbenRepo`, `TestRunSlice_RejectsWaveSlice`)
- **klasse:** test-kommentar-rot-bedingung-nicht-getragen

### F-4 — LOW — Zwischenstand `21d6f711` ist ebenfalls rot und im Plan nicht als solcher benannt

- **quelle:** `AGENTS.md` §3.3 (ein roter Zwischenstand ist zulässig, solange er nicht die Spitze eines Push wird), `slice-298` §3 Plan-Änderung 1 Punkt 4 („Zwischenstand, der nicht die Spitze eines Push werden darf: (2a) allein“)
- **pfad:** `docs/plan/planning/in-progress/slice-298-neue-namen-wave-archive-anker.md` · „Zwischenstand, der nicht die Spitze eines Push werden darf“
- **befund:** Der Plan nennt nur den reinen Move (2a) als roten Zwischenstand; die zweite Plan-Änderung schiebt die Plan-Datei (2b) davor und sagt nichts dazu. Gemessen sind beide rot: `9d295355` meldet `1328 Datei(en) geprüft, 28 Befund(e)` (`target-missing`) und `wave-drift` für `welle-98`; `21d6f711` meldet dieselben Klassen. Die Botschaft von `9d295355` benennt es, die von `21d6f711` nicht.
- **verifizierbar:** ja — Wegwerf-Klon, `git checkout 21d6f711`, `make planning-check`: `docs/plan/planning/in-progress/roadmap.md:10 welle-98 wave-drift … kein flaches Wellendokument dieser Kennung`
- **klasse:** roter-zwischenstand-nur-teilweise-benannt

### F-5 — LOW — Der Nachzug hat im Verifikations-Report zu `slice-297` eine Mess-Aussage selbstwidersprüchlich gemacht

- **quelle:** `MR-070` (Eigenschafts-Test: würde ein korrigierter Wert verfälschen, was das Artefakt zu seinem Datum festhielt?), `slice-298` §3 Plan-Änderung 1 Punkt 3 (Klasse „Aussagen, die den alten Namen als Gegenstand tragen“)
- **pfad:** `docs/reviews/2026-11-26-slice-297-spezifikation-heisst-specification-verify.md` · „Keine Datei zu `slice-296`/`slice-298`, `archive`, `stub`, `wave-` oder einer neuen ADR im Diff“
- **befund:** Der Satz hielt fest, dass der Filter `wave-` über den Diff `f44c6e65..HEAD` leer blieb und die Welle (damals `welle-98`) nur als Link-Ziel-Nachzug vorkam. Der Nachzug setzte in der Klammer `wave-98` ein; der Satz sagt nun im selben Atemzug „keine Datei zu `wave-`“ und „`wave-98` kommt vor“. Die Bedeutung der Messung ist nicht mehr lesbar.
- **verifizierbar:** ja — `git grep -n 'Keine Datei zu' -- docs/reviews`; `git show 2afe2d47 -- docs/reviews/2026-11-26-slice-297-spezifikation-heisst-specification-verify.md`
- **klasse:** nachzug-macht-filter-aussage-selbstwidersprüchlich

### F-6 — LOW — `ADR-0117` trägt zwei Zahlen, die kein Kommando reproduziert, und wird mit der Closure unveränderlich

- **quelle:** `AGENTS.md` §5 Regel 14 (vor einer Messung steht die ausgeschriebene Form ihres Gegenstands), `AGENTS.md` §3.5 (`Accepted` ist immutabel); Architect-Verdikt Frage 3 (die ADR geht mit der Closure nach `Accepted`)
- **pfad:** `docs/plan/adr/0117-wellen-namen-ueber-die-konfiguration.md` · „344 Pfade mit `welle-<NN>`, 188 `archiv.zip`, 24 Anker und 38 Register-Zeilen sind Belege“
- **befund:** Die Alternative B nennt 344 Pfade und 24 Anker ohne Kommando. Am Arbeitsbaum ergeben die Kommandos aus `MR-082` 345 Pfade und 19 Anker (der Plan selbst hält fest, dass die 24 „kein Kommando reproduziert“, gemessen 36 Zeilen, davon 17 Platzhalter und 19 mit Ziffern). Die ADR ist bis zur Closure `Proposed` und änderbar; danach steht die Abweichung zwischen ADR und `MR-082` unveränderlich. Nicht Teil des Diffs; Gegenstand hier, weil sie die Senkung trägt, die der Diff konfiguriert.
- **verifizierbar:** ja — `git ls-files | grep -c 'welle-[0-9]'` (345), `git grep -h -o 'seit welle-[0-9][0-9]*' -- . ':!.harness/baseline' ':!docs/plan/planning/done' ':!docs/reviews' | wc -l` (19)
- **klasse:** zahl-in-immutabler-adr-ohne-kommando

### F-7 — INFO — Fünf der dreizehn neuen Tests bleiben ohne die Änderung grün

- **quelle:** `slice-298` §6 Risiko 2 (Ausgang „entfallen“ nur, wenn jeder neue Test rot war), `BEO-ALL/testfall-faehrt-seine-regel-nicht`
- **pfad:** `tools/archive-wave/wave_names_test.go` · „func TestExtractTitle_Wave“
- **befund:** Gesamt-Rücknahme (M7: `collect.go` auf `db329ff2`, `ArchiveName` immer `archiv.zip`, Suffix `-archiv.zip`): rot sind elf Tests (sieben neue, vier bestehende), **grün bleiben** `TestExtractTitle_Wave`, `TestFindWellePlan_Wave`, `TestRun_DryRun_NoWrites_Wave` (das Werkzeug ist dort präfix-generisch, sie sind Regressionswächter, keine Prüfung der Änderung) und `TestFixture_Bestand_SchreibtWeiterArchiv`, `TestRunSlice_AcceptsContextualWaveMention` (Gegenseiten, gewollt grün). Der Ausgang von Risiko 2 darf „jeder rot“ so nicht behaupten; die Botschaft von `928cb73f` behauptet es nicht (neun Mutationen, gemessen reproduziert).
- **verifizierbar:** ja — Wegwerf-Klon, Rücknahme wie beschrieben, `make test`, Liste der `--- FAIL` gegen die Funktionsnamen der Datei
- **klasse:** neuer-test-bestaetigt-nur-bestand-oder-generisches

### F-8 — INFO — Die Wortgrenze der Kennungs-Regex ist nicht getestet

- **quelle:** `slice-298` §6 Risiko 2; Prüffrage 13 (Negativtest zu einem neuen Vertrag) — hier nicht öffentlich und im Bestand vorher ebenso ungetestet, daher INFO
- **pfad:** `tools/archive-wave/collect.go` · „\b(welle|wave)-(\d+)\b“
- **befund:** M6 (beide `\b` entfernt): alle Tests grün. Die Regex-Probe selbst verhält sich richtig: `wave-1` und `wave-100` treffen, `wave-x`, `awave-98`, `xwelle-98`, `wave-98x`, `wave-98_`, `Wave-98`, `wave--98` nicht; `welle-wave-98` liefert `wave-98` und `wave-welle-97` liefert `welle-97` (erstes Vorkommen, wie bisher bei `welle-`). Kein Test hält eine dieser Formen.
- **verifizierbar:** ja — Wegwerf-Klon, `\b` aus `welleIDInFieldRE` entfernen, `make test`
- **klasse:** wortgrenze-ohne-negativtest

### F-9 — INFO — Der bestehende Test zum Slice-Modus bleibt aus dem falschen Grund grün

- **quelle:** `BEO-ALL/testfall-faehrt-seine-regel-nicht` (2×, offen)
- **pfad:** `tools/archive-wave/slice_mode_test.go` · „func TestRunSlice_RejectsWelleSlice“
- **befund:** M8 (die Wellen-Zugehörigkeits-Prüfung in `runSlice` abgeschaltet): nur der neue `TestRunSlice_RejectsWaveSlice` wird rot; der Bestands-Test bleibt grün, weil sein Fixture kein `docs/reviews` anlegt und der Lauf an anderer Stelle mit einem Fehler endet. Der neue Test schließt das mit `.gitkeep` und der Meldungs-Prüfung für die `wave-`-Seite; die `welle-`-Seite bleibt schwach. Eine dritte Instanz der beobachteten Klasse, nicht Gegenstand dieses Diffs (unveränderte Zeilen).
- **verifizierbar:** ja — Wegwerf-Klon, `if welleIDInFieldRE.MatchString(firstLine)` durch `if false && …` ersetzen, `make test`
- **klasse:** testfall-faehrt-seine-regel-nicht

### F-10 — INFO — Kommentare nennen `welle-<NN>` als Wert für das Wellen-Verzeichnis von `slice-mv`

- **quelle:** `MR-082` (Geltungsbereich nennt weder `Makefile` noch `tools/harness/`)
- **pfad:** `Makefile` · „(Wellen-Slice: welle-<NN>, wellenlos: ungesetzt)“
- **befund:** Das Target `slice-mv` reicht das Verzeichnis je Aufruf durch; für eine neue Welle wäre es `wave-<NN>`. Gemessen trägt `tools/harness/slice-mv.sh` kein Muster im Code (nur die Kommentarzeilen 103, 170, 295); das Verhalten stimmt. Der Plan (§8) hat das als „nur Kommentarzeilen“ erwartet; die Kommentare nennen nur die Bestandsform.
- **verifizierbar:** ja — `git grep -n 'welle-' -- Makefile tools/harness/slice-mv.sh`
- **klasse:** kommentar-nennt-nur-bestandsform

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Kennungs-Erkennung `collect.go` | geprüft, ohne Befund — Regex `(welle\|wave)-(\d+)` mit Wortgrenzen, Vergleich auf Präfix **und** Ziffern (M2, nur Ziffern: `TestCollectSlices_WaveUndWelleGleicheZiffernGetrennt` rot mit der Meldung `wave-98 sammelt [slice-701-a.md slice-702-b.md slice-703-c.md slice-704-d.md]`); Regex-Probe siehe F-8 |
| `ArchiveName` und Stubs | geprüft, ohne Befund — M3 (immer `archiv.zip`): `TestArchiveName` (`ArchiveName("wave-98") = "archiv.zip", erwartet "archive.zip"`), `TestFixture_EndToEnd_Wave`, `TestFixture_Gemischt_…`, `TestStubs_Wave` rot; M4 (immer `archive.zip`): sieben Bestands-Tests rot (u. a. `TestFixture_EndToEnd`, `TestSliceStub`, `TestArchiveName` für `welle-97`); je aus dem benannten Grund |
| Slice- und Review-Modus (`-archive.zip`, Verdikt Frage 8) | geprüft, ohne Befund — M5 (Suffix `-archiv.zip`): vier Tests rot (`TestRunSlice_SchreibtNurDieNeueForm`: `erwartet genau slice-601-archive.zip, got [slice-601-archiv.zip]`); Konstante `SliceArchiveSuffix` an beiden Stellen, beide Stub-Texte mitgezogen |
| Bestand byte-gleich (Dry-Run) | geprüft, ohne Befund — je `welle-97`, `welle-96`, `welle-95` mit dem Image von `db329ff2` und dem von `HEAD` auf demselben Baum: stdout je `cmp` gleich, Exit je 0 gleich; stderr unterscheidet sich nur im Docker-Build-Log (Reihenfolge der Schritte, `transferring context: 98.94kB` gegen `115.19kB`, Image-Namen), nicht in der Ausgabe des Werkzeugs |
| Slice-Modus lehnt wellengebundene Slices ab (`**Welle:** wave-98`) | geprüft, ohne Befund — am Baum von `HEAD` lehnt das neue Image `slice-295` ab (`gehoert zu einer Welle … -- das gehoert in den -welle-Modus, nicht -slice`), das alte Image nimmt ihn wegen des geänderten Feldes als wellenlos an. Das ist keine Verhaltensänderung für den Bestand, sondern die Folge des umbenannten Feldes; M8/M1 belegen den Test dazu |
| `-welle=wave-98` am Baum von `HEAD` | geprüft, ohne Befund — Dry-Run findet `slice-295`, `slice-297` und vier Reports, zeigt `Welle-Plan: (keiner …)` (die Welle liegt noch flach, nicht in `done/`), Exit 0 |
| `planning.waves` (Konfiguration) | geprüft, ohne Befund — Probe in Wegwerf-Klon: unregistrierte flache `wave-99`-Datei ⇒ `wave-drift` (`roadmap.md:10 wave-99 … nennt es nicht`, Exit 2), mit Zeiger `0 Befund(e)`; `wave-results-missing` für eine Registerzeile `wave-97` ohne Notiz; **Gegenseite** (Bestand ungewächtert, wie im Kommentar benannt): unregistrierte `welle-77`-Datei, entfernte Registerzeile `welle-97`, entfernte `welle-97-results.md` je `0 Befund(e)` |
| Wortlaut der Senkung (Regel 13) | geprüft, ohne Befund — Konfig-Kommentar, `planning-check.md`, `MR-082`, die zwei Commands und `ADR-0117` nennen dieselbe Grenze (Wellen-Invariante und Register-Deckung der Bestands-Wellen laufen nicht mehr; `matrix`, Closure-Profil, Werkzeug lesen beide); die Kommandos zu 38 Registerzeilen und 38 Ergebnisnotizen ergeben 38 und 38 |
| Matrix-Klasse `wave` und Verbote (V2) | geprüft, ohne Befund — `wave-98` im Fließtext von `spec/specification.md` (Z. 30), `spec/lastenheft.md` (Z. 29), `spec/architecture.md` (Z. 159), `ADR-0117` (Z. 30) ⇒ je ein `matrix-forbidden` (`spec-straten → wave`, `sicht → wave`, `adr → wave`); `wave-drift`, `wave-<NN>`, `wave-9` ⇒ `0 Befund(e)`; `ADR-0117` selbst trägt 0 Befunde. (Zeilen am Dateiende hinter der Historie-Tabelle sind nicht scanbar — eine erste Probe dort war ein Fehler meiner Probe, kein Befund; die Probe wurde mitten im Text wiederholt.) |
| `ids`-`exempt-paths` `done/wave-*/*.md` | geprüft, ohne Befund — Stub mit nackten `ADR-0001`/`DC-FA-PLAN-001` unter `done/wave-99/` und `done/welle-99/` ⇒ kein Befund; Kontrolle `done/xwave-99/` ⇒ zwei `id-unlinked` |
| Closure-Regel `*-results.md` (Plan-Änderung `21d6f711`) | geprüft, ohne Befund bis auf F-1/F-2 — trifft genau die 38 Bestandsnotizen (`git ls-files | grep -cE 'done/[^/]*-results.md$'` = 38, davon ohne `done/welle-` 0); Probe: `wave-99-results.md` ohne den Abschnitt ⇒ `section-missing`, mit leerem Abschnitt ⇒ `section-empty`; nichts Fremdes erfasst (`*` überschreitet das Verzeichnis nicht) |
| Reiner Move `9d295355` | geprüft, ohne Befund — `R100`, ein Pfad |
| Nachzug von `welle-98` (MR-070) | geprüft, mit F-5 — am Vorzustand 71 Zeilen in 13 Dateien, am Stand `HEAD` 32 Zeilen in 7 Dateien; die verbleibenden Klassen sind Dateiname des Verdikts (lebende Links auf eine vorhandene Datei), Aussagen über den Umzug (Plan §1/§3, Welle §5 Frage 2, Verdikt, `slice-297` Z. 384) und Fixture-Kennungen der Testdatei; keine ist ein lebender Verweis, der bräche; alle Links lösen auf (`make doc-check` 0 Befunde) |
| `MR-082` gegen die Vorlage und `MR-081` | geprüft, ohne Befund — Pflichtfelder Datum, Geltungsbereich, Ersetzt-Baseline-Regel (genau eine, Link mit Anker), Adaption, Begründung, Auflösungs-Trigger; Template-Hinweis und Regel-Absatz entfernt; Index-Zeile in `harness/conventions.md`; Bezug auf `ADR-0117` ohne Nacherzählung (Regel 18) |
| `MR-082` Zahlen und Trigger | geprüft, ohne Befund — `git ls-files \| grep -c 'welle-[0-9]'` = **345**, `grep -c 'archiv\.zip$'` = **188**, `grep -c '^\| welle-'` Roadmap = **38**, Trigger-Kommando 1 = **435**, Trigger-Kommando 2 = **19**; der Trigger ist nicht „permanent“ und nennt beide Zähler (Verdikt Frage 1) |
| `AGENTS.md` §3.7 | geprüft, ohne Befund — eine Zeile ersetzt (Zeilenzahl gleich, `d-check:cite` in `reviewer.md` bleibt gültig, `AGENTS.md:287` unverändert); `seit wave-<NN>` ist die Form des vendorten `v7.0.0` (`modul-09-implementierung.md`, `modul-13-quality-gates.md`), die Instanzen `seit welle-73`/`-80` bleiben |
| Commands, Sensor-Dateien, Werkzeug-README | geprüft, ohne Befund — nur Namen und die Grenze geändert, Platzhalter und Prosa deutsch unberührt; die Aussagen zu Archiv-Namen je Modus stimmen mit dem Code (`ArchiveName`, `SliceArchiveSuffix`) |
| Abgrenzung | geprüft, ohne Befund — `git diff --name-only db329ff2..HEAD \| grep -c '^internal/'` = 0; `done/` nur `slice-295`/`slice-297` (Link-Nachzug), `observations/` unberührt, `slice-296` nur Link-Ziele; keine Default-Änderung (`internal/…/config.go` unberührt, `spec/` unberührt); keine Neuerung aus `slice-296`; `ADR-0117` unverändert `Proposed` (letzter Commit an der ADR: `db329ff2`) |
| Hooks, Skripte, CI auf feste `welle-`-Muster | geprüft, ohne Befund — `.githooks/pre-commit` und CI erkennen `(.+/)?slice-[0-9]+`, kein Präfix; `.github/workflows`, `tools/*.sh`, `tools/harness/*.sh`: keine `welle-`-Muster im Code |
| Zählungen des Plans (Regel 14/18) | geprüft, ohne Befund — selbst nachgezählt: `archiv.zip` im Werkzeug 12 (`archive.go` 7, `stub.go` 4, `main.go` 1), `seit welle-` 36 Zeilen (17 Platzhalter, 19 mit Ziffern) am Vorzustand, `welle-98` 71 Zeilen/13 Dateien am Vorzustand, 345 (Vorzustand 346) Pfade, 188 Archive, 38 Registerzeilen |
| Botschaften (Regel 15) | geprüft, ohne Befund — `928cb73f`: „Neun Mutationen“ = Regex, Ziffern, `ArchiveName` je Richtung (2), Suffix, vier Stub-Texte; die fünf Mutationsklassen, die ich fuhr (M1 bis M5), sind rot aus dem benannten Grund. `2afe2d47`: „elf Dateien“ stimmt (Roadmap, `slice-295`, `slice-297`, `slice-296`, Welle, vier Reports, Verdikt, dazu der Anker in `.d-check.yml`; die Aufzählung in der Botschaft nennt das Verdikt nicht, die Zahl trägt) — **kein** Overclaim; `.d-check.closure.yml` kommt als Konfiguration hinzu (zwölf Dateien im Commit) |
| Kommentare in Go und Konfiguration (§3.7) | geprüft, mit F-1, F-2, F-3 — `ArchiveName`, `SliceArchiveSuffix`, `welleIDInFieldRE`, die Kommentare an `planning.waves`, der Klasse `wave` und der `ids`-Pfade tragen Zusage, Kopplung oder Grenze, Herkunft ein Feld (`MR-082`, `ADR-0117`); keine Review-Historie, keine Slice-Nummern |
| Gates, die ich gefahren habe | geprüft, ohne Befund — `make archive-wave-test` Exit 0; `make doc-check` `d-check: 1329 Datei(en) geprüft, 0 Befund(e)`; `make planning-check` `1323 Datei(en) geprüft, 0 Befund(e)`; `make verify-closure-notes` `1182 Datei(en) geprüft, 0 Befund(e)`; `make gate-consistency` `1323 Datei(en) geprüft, 0 Befund(e)`; `make adr-check RANGE=db329ff2..HEAD` (6 Commits) `1323 Datei(en) geprüft, 0 Befund(e)`; `make trace-check RANGE=db329ff2..HEAD` `1323 Datei(en) geprüft, 0 Befund(e)` |

## Mutationsproben am Werkzeug (Wegwerf-Klon, nie im Arbeitsbaum)

| # | Änderung | rot (Auswahl) | Meldung gelesen |
|---|---|---|---|
| M1 | `collect.go` auf `db329ff2` | vier Tests, nicht `TestRun_PlanOhneEigeneSlices_Wave` | `wave-98 sammelt []`; `erwartet 3 Moves …, got 1`; `erwartet die Ablehnung wegen der Wellen-Zugehoerigkeit, got <nil>` |
| M2 | Vergleich nur über die Ziffern | `TestCollectSlices_WaveUndWelleGleicheZiffernGetrennt` | beide Kennungen sammeln alle vier Slices |
| M3 | `ArchiveName` immer `archiv.zip` | fünf Tests | `ArchiveName("wave-98") = "archiv.zip", erwartet "archive.zip"` |
| M4 | `ArchiveName` immer `archive.zip` | sieben Tests, darunter Bestands-Tests | `ArchiveName("welle-97") = "archive.zip", erwartet "archiv.zip"` |
| M5 | `SliceArchiveSuffix` `-archiv.zip` | vier Tests | `erwartet genau slice-601-archive.zip, got [slice-601-archiv.zip]` |
| M6 | `\b` aus der Regex | keiner (F-8) | — |
| M7 | Gesamt-Rücknahme (M1+M3+M5) | elf Tests; fünf neue bleiben grün (F-7) | — |
| M8 | Wellen-Zugehörigkeits-Prüfung in `runSlice` aus | nur `TestRunSlice_RejectsWaveSlice` (F-9) | `erwartet die Ablehnung wegen der Wellen-Zugehoerigkeit, got <nil>` |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 6 (F-1 bis F-6) |
| INFO | 4 (F-7, F-8, F-9, F-10) |


## Verdikt

Der Diff trägt den Plan: Das Werkzeug erkennt beide Schreibweisen, vergleicht Präfix und Ziffern, schreibt das Archiv nach dem Präfix der Kennung und lässt den Bestands-Lauf byte-gleich; die Konfiguration bewacht die neue Form (beide Seiten gefahren, die Senkung ist im Konfig-Kommentar, in `planning-check.md`, `MR-082` und `ADR-0117` mit demselben Wortlaut benannt und ihre Zahlen stimmen); der Move ist `R100`; kein Bestand wurde umbenannt, `internal/` und die Defaults sind unberührt, `ADR-0117` ist `Proposed`. Kein HIGH und kein MEDIUM blockiert. Die LOW sind Wartungs- und Beleg-Fragen: zwei Kommentare der Closure-Regel (F-1, F-2), zwei Test-Kommentare (F-3), die Benennung des zweiten roten Zwischenstands (F-4), ein Nachzugs-Nebeneffekt im Verifikations-Report (F-5) und zwei Zahlen in einer ADR, die vor `Accepted` noch änderbar ist (F-6). Zur Übergabe: F-1, F-2, F-3, F-5 an den Implementer; F-4 und die Hinweise F-7, F-9 für die Closure-Notiz (Risiko 2, Register) an den Planner; F-6 an den Architect, **bevor** `ADR-0117` nach `Accepted` geht.
