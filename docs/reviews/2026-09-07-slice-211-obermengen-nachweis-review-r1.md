# Review-Report — slice-211, Runde 1 (Obermengen-Nachweis gegen `markdownlint` MD013 und der CR-Entscheid)

**Review-Art:** Code/Design — geprüft werden der **Obermengen-Nachweis** und der **Entscheid** im eingehenden Zeilenlängen-CR gegen den Slice-Plan slice-211, gegen den Kanon `v6.5.0` · `regelwerk/modul-11-verification.md` §Fitness Function ohne Standard-Tool und gegen die Akzeptanzkriterien des CR selbst. **Nicht** geprüft: die DoD-Abhakung (Verifier-Rolle) und der §Zwischenstand des CR (liegt vor der Range).
**Gegenstand:** `a15f2dee` (Plan angelegt) · `f9161d96` (Beanspruchung) · `6088b616` (Nachweis) · `c63a5746` (Entscheid). Diff-Range `a15f2dee^..c63a5746`. **Achtung:** `c63a5746` ist während des Reviews von `HEAD` nach `HEAD~1` gerutscht — eine parallele Sitzung hat während des Reviews `8b10cbf4` und `ecbb5e5f` nachgeschoben (ein anderer, unbeteiligter CR zu links/anchors). Geprüft ist die Range oben, nicht `HEAD`.
**Skill:** `.harness/skills/reviewer.md` v1.15.0 @ `c63a5746` — einschließlich Anker 17, der mit slice-210 entstand und hier zum ersten Mal auf eine Messung trifft.
**Modell-ID:** claude-opus-5[1m] · **Datum:** 2026-09-07
**Eingangs-Kontext:** der Slice-Plan slice-211 (§1, §3, §6, §8, DoD 1–3); der CR `docs/plan/cr/2026-09-06-cr-eingehend-auftraggeber-zeilenlaenge.md` (§Akzeptanzkriterien, §Zwischenstand, §Obermengen-Nachweis, §Entscheid); `AGENTS.md` §3.1/§3.3/§3.4/§3.7/§3.8/§4/§5/§6; `harness/README.md`; `harness/conventions.md` samt `MR-013`, `MR-046`, `MR-053`, `MR-054`, `MR-070`; `DC-FA-STRUCT-001`, `DC-QA-02`, `DC-QA-03`; Baseline `v6.5.0` · `regelwerk/modul-05-planning-harness.md`, `regelwerk/modul-11-verification.md`, `regelwerk/modul-13-quality-gates.md` sowie `templates/docs/reviews/review-report.template.md`.

**Eigene Läufe (echte Ausgaben, keine behaupteten):**

- `make gates` — `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`, `d-check: 712 Datei(en) geprüft, 0 Befund(e)`, `coverage-gate: OK — Coverage 94.60% erfüllt Schwelle 93%`, Exit 0.
- `make doc-check` (nach den Fremd-Commits) — `d-check: 713 Datei(en) geprüft, 0 Befund(e)`, mit diesem Report im Baum `714 Datei(en) geprüft, 0 Befund(e)`. Die Differenz zu den 712 des `gates`-Laufs sind die beiden Fremd-Commits und dieser Report, nicht ein Artefakt dieses Slice; ein noch späterer Lauf meldete `715 Datei(en) … 0 Befund(e)`, weil die parallele Sitzung weiterschreibt. Jeder Lauf grün.
- `make nightly-state` — `gruen — juengster Lauf 2026-09-07T05:33:45Z` (`upstream-drift.yml`) und `gruen — juengster Lauf 2026-09-07T08:21:32Z` (`image-scan.yml`).
- **Beide Sensoren über die acht Proben** des Nachweises, mit dem im CR dokumentierten Digest `davidanson/markdownlint-cli2@sha256:173cb697…fed63a` (lokal vorhanden, kein Netz nötig) und `d-check:latest`.
- **Vier eigene Gegenproben** zur Proben-Isolation (K3 mit `code_blocks: true`, K4 mit `tables: true`, K7 mit auf den anderen Abschnitt umskopierter Regel, MD013 mit Schwelle 2046 über den ganzen Bestand).
- **Zwei eigene Bestands-Läufe** (MD013 und `forbid-pattern`, Schwelle 1000) über eine frisch aus `git ls-files` gebaute Kopie **und** über die Original-Kopie des Implementers.
- **Rohzählung** der Zeilenlängen über vier Skopen mit `awk`.
- Kein Host-Interpreter, keine Schreiboperation im Repo: alle Läufe über `make` bzw. `docker run --rm --network none` mit `:ro`-Mount; `git status` im Repo trug zu Beginn nichts und am Ende nur diesen Report.

**Zitier-Form.** Dieser Report friert ein: Baseline-Stellen als Tag + Pfad in Inline-Code, Slices als Kennung, Gates als `make <target>`-Token.

---

## Was gemessen reproduziert (vorweg, damit die Findings nicht als Gesamturteil gelesen werden)

| Behauptung | Quelle | Nachmessung |
|---|---|---|
| MD013 meldet K1, schweigt bei K2–K6a | §Obermengen-Nachweis Teil 1 | **reproduziert exakt** — `k1.md:7:301 … Actual: 343`, sonst nur `k6b` und `k7` |
| `forbid-pattern` meldet K1, K4, K5, K6a; schweigt bei K2, K3 | §Teil 1 | **reproduziert exakt** — fünf `section-forbidden`, Exit 1 |
| „Sechs von sechs gegen drei von sechs" | §Teil 1 | **reproduziert** — unter der Klassen-Form, die der Nachweis benutzt (siehe F-1) |
| Nebenfund K6b: Inline-Code **mit** Leerzeichen meldet bei beiden | §Teil 1 | **reproduziert** — `k6b.md:7:301 … Actual: 776` und `k6b.md:3 section-forbidden` |
| RE2-Decke: `.{1000,}` läuft, `.{1001,}` ist Konfigurationsfehler mit Exit 2 | §Teil 2 | **reproduziert wörtlich** — `forbid-pattern ".{1001,}" ist kein gültiges Regex: … invalid repeat count` |
| K7: MD013 meldet, `structure` schweigt (Abschnitts-Skopus) | §Teil 3 | **reproduziert** — und die Gegenprobe zeigt, dass das Schweigen **nur** aus der Skopierung kommt |
| Bestand bei Schwelle 1000: MD013 **9** Dateien | §Der unveränderte Bestand | **reproduziert exakt**, auf beiden Kopien; die Rohzählung „Datei mit Nicht-Tabellen-Zeile > 1000" ergibt ebenfalls 9 |
| „null Dateien meldet nur MD013" | §Der unveränderte Bestand | **reproduziert** — die 9er-Menge ist echte Teilmenge der `forbid-pattern`-Menge |
| „Gegenwert: neun Dateien oberhalb von 1000 Zeichen" | §Entscheid, Frage 3 | **reproduziert** |
| „vierte Toolchain … Pin, Target, zwei Deklarationen, Pin-Spiegel-Klassen" | `MR-046` | **wortgetreu übernommen**, Geltungsbereich trägt |
| §8: 37 Register-Verzeichnisse; fünf Einträge mit 3×/11×/8×/9×/2× | Slice-Plan §8 | **reproduziert exakt**, Zähler und Stände |
| §8: Nachtlauf-Stand beider Läufe grün, mit Zeitstempeln | Slice-Plan §8 | **reproduziert exakt**, beide Zeitstempel identisch |
| `d-check:cite`-Anker `268-269` und `274-274` | Slice-Plan §8 | **wortgleich**, vorschreibende Zeilen, identisch mit der Haus-Form von slice-207/208/209 |
| „zehn Gates, 712 Dateien, 0 Befunde" | Botschaften `6088b616`, `c63a5746` | **reproduziert exakt** |
| Beanspruchungs-Commit ist reiner Rename plus Ruhe-Marker | `f9161d96` | **reproduziert** — exakte Umkehr des slice-210-Closure-Commits, Blob-Hash `a7bfab5f` wiederhergestellt |
| „Zwei Proben mussten neu gebaut werden" | §Teil 1 | **plausibilisiert, nicht nachprüfbar** — die Vorfassungen existieren nicht mehr; die heutigen Proben tragen genau die beschriebenen Korrekturen (k1 endet auf Prosa mit Leerzeichen, k6a enthält **kein** Whitespace) |

Der **Entscheid selbst — dreimal Nein — wird von keinem Finding berührt.** Alle Findings treffen die **Belege**, nicht die Schlüsse; zwei von ihnen machen den Schluss sogar stärker, als der Nachweis ihn führt.

---

## Findings

### F-1 · MEDIUM · Die Klasse K6 wurde nach §3 neu definiert; die Zählung „sechs von sechs" zählt eine Klasse, die der Plan nicht führt — und die Klasse, die er führt, fällt bei MD013 durch

- **quelle:** `AGENTS.md` §5 (*Vor einer Messung steht die Form ihres Gegenstands*, `seit slice-210`); Reviewer-Anker 17
- **pfad:** `docs/plan/cr/2026-09-06-cr-eingehend-auftraggeber-zeilenlaenge.md:210` und `:212` gegen `docs/plan/planning/in-progress/slice-211-obermengen-nachweis-md013.md:95`
- **befund:** §3 des Plans schreibt die Form **vor** der Messung aus und führt `| K6 | Inline-Code-Spanne | kein Befund |`. Gemessen und gezählt wurde stattdessen `K6a unteilbares Inline-Code-Token`; die Probe, die der Plan-Klasse entspricht — eine Zeile, die **nur** aus einer Inline-Code-Spanne besteht —, ist als K6b aus der Zählung herausgenommen und meldet bei MD013 (reproduziert: `k6b.md:7:301 MD013/line-length … Actual: 776`). Dieselbe Klasse steht im §Zwischenstand desselben Dokuments als CR-Forderung (`:139` — *„Inline-Code-Spanne | meldet | soll schweigen"*). Die Spalte *„vom CR gefordert"* schreibt dem CR damit in Zeile 210 eine Forderung zu, die er in dieser Verengung nicht stellt. Die Folge steht im Entscheid: Umkehr-Bedingung 2 (`:328`) sagt, die fehlende Abschnitts-Skopierung werde dann *„zur einzigen offenen Frage"* — gemessen bliebe mindestens eine zweite offen.
- **verifizierbar:** ja — der dokumentierte `docker run` über die Proben zeigt `k6b` in der MD013-Ausgabe; §3 des Plans und §Zwischenstand `:139` zeigen die abweichende Klassendefinition.
- **klasse:** `klassen-definition-nach-der-messung-verengt`

### F-2 · MEDIUM · Teil 2 und Umkehr-Bedingung 3 belegen sich mit Zahlen aus einem anderen Bestand als dem gemessenen — und mit der Zeilen-Klasse, die der geprüfte Sensor gerade **nicht** ausnimmt

- **quelle:** `AGENTS.md` §5; Reviewer-Anker 17
- **pfad:** `docs/plan/cr/2026-09-06-cr-eingehend-auftraggeber-zeilenlaenge.md:226`, `:294`, `:330`
- **befund:** *„Die längste Zeile des Bestands misst 2045 Zeichen"* trägt Teil 2 und den Entscheid zu Frage 3. Im Bestand, den der Nachweis zwei Absätze später misst (`docs/`, `spec/`, `harness/`), ist die längste **Nicht-Tabellen**-Zeile 1478 (`spec/lastenheft.md:397`) und die längste Zeile überhaupt 3083 (`spec/lastenheft.md:3847`); die 2045 steht in `.harness/baseline/v6.5.0/regelwerk/modul-06-roadmap.md:62`, also im vendorten Fremdtext, den derselbe CR im §Anlass ausdrücklich als Fremdtext ausweist. Die gewählte Population zählt Tabellenzeilen heraus — genau die Klasse, die `forbid-pattern` gemessen **meldet** (K4). Umkehr-Bedingung 3 (`:330`) mischt zusätzlich zwei Mengen und zwei Einheiten: die *neun Dateien* stammen aus dem Bestandslauf und reproduzieren, die *„eine über 2000"* ist dort **null** (gemessen: 0 Dateien mit Nicht-Tabellen-Zeile über 2000) und stammt aus der §Anlass-Tabelle, wo sie eine **Zeile** zählt, keine Datei. Die Bedingung ist damit nicht nachmessbar. Der Schluss *„konstruktiv unerreichbar"* bleibt richtig — mit 3083 statt 2045 sogar deutlicher —, aber er ruht auf einer Zahl, die den Gegenstand nicht misst.
- **verifizierbar:** ja — `awk`-Rohzählung über `git ls-files 'docs/*.md' 'spec/*.md' 'harness/*.md'` ohne `/archiv/`; zusätzlich zeigt MD013 mit `line_length: 2046` über denselben Bestand `Summary: 0 error(s)`, die Kanon-Hälfte *„der unveränderte Bestand, auf dem beide schweigen"* ist also für MD013 erreichbar und nur für `forbid-pattern` nicht.
- **klasse:** `zahl-aus-anderer-menge-als-die-aussage`

### F-3 · MEDIUM · „13 nur `forbid-pattern` — und das sind genau die drei Klassen aus Teil 1" — gemessen ist es **eine** Klasse

- **quelle:** `AGENTS.md` §5; Reviewer-Anker 17
- **pfad:** `docs/plan/cr/2026-09-06-cr-eingehend-auftraggeber-zeilenlaenge.md:247`; wiederholt in `:310` als *„seine drei Falsch-Positiv-Klassen"*
- **befund:** Alle 13 Dateien der Differenzmenge sind nachgezählt: **jede** Zeile ab 1000 Zeichen in ihnen ist eine **Tabellenzeile** (K4). K5 (unteilbares Token) und K6a (Inline-Code-Token) kommen in der Differenzmenge bei Schwelle 1000 **nicht vor** — von 8 Zeilen in `AGENTS.md` über 23 in `spec/spezifikation.md` bis 16 in `docs/user/benutzerhandbuch.md` beginnt jede mit `|`. Die Trefferliste wird damit als Beleg für drei Ursachen gelesen, obwohl eine sie vollständig erklärt. Wer den `forbid-pattern`-Weg auf dieser Grundlage bewertet — der Entscheid empfiehlt ihn als Ad-hoc-Messung mit dem Zusatz *„wer ihn benutzt, weiß, was er bekommt"* —, rechnet mit drei zu behandelnden Klassen, wo gemessen eine trägt.
- **verifizierbar:** ja — `awk 'length($0)>=1000'` über die 13 Dateien; die drei Klassen selbst über die Proben k4/k5/k6a.
- **klasse:** `ursachen-attribution-ohne-nachzaehlung`

### F-4 · MEDIUM · Die Bestands-Messung ist aus ihrer eigenen Beschreibung nicht reproduzierbar: mit „`docs/`, `spec/`, `harness/`" ergibt der `forbid-pattern`-Lauf **21** Dateien, nicht 22

- **quelle:** `DC-FA-STRUCT-001` (§Kandidaten-Menge); DoD 2 des Slice-Plans (*„reproduzierbar dokumentiert"*)
- **pfad:** `docs/plan/cr/2026-09-06-cr-eingehend-auftraggeber-zeilenlaenge.md:244`, dazu `:191-192`
- **befund:** Aus einer Kopie, die exakt die genannten Wurzeln enthält, meldet `forbid-pattern` **21** Dateien; die 22. ist `AGENTS.md`, die außerhalb der genannten Wurzeln liegt und nur deshalb im Lauf des Implementers steht, weil seine Kopie sie zusätzlich enthielt. Die Ursache ist keine Nachlässigkeit im Kopieren, sondern eine Eigenschaft, die genau die Anforderung beschreibt, um die es hier geht: `DC-FA-STRUCT-001` sagt, der `files`-Glob werde *„über den gesamten Baum — **unabhängig** von `scan.roots`/`scan.ignore`"* ausgewertet. Die verwendete Konfiguration setzte aber `scan.roots: ["docs","spec","harness"]` und `ignore: ["**/archiv/**"]`; beides ist für `structure` wirkungslos. Der Kandidaten-Satz der `forbid-pattern`-Seite war der **ganze** kopierte Baum (702 Markdown-Dateien einschließlich `AGENTS.md`, `README.md` und 11 Archiv-Dateien), nicht der beschriebene — nachgestellt an einem Minimalbaum mit `scan.roots: ["docs"]`, in dem `structure` einen Befund in `root.md` meldet, während der Zähler `1 Datei(en) geprüft` sagt. Dazu fehlen zwei Angaben, ohne die der Lauf nicht wiederholbar ist: der Abschnitts-Selektor (`section-pattern: '^## '`, `sections: each`), von dem der CR selbst sagt, die Zahl hänge an ihm (§Zwischenstand: *„342 gegen 645"*), und die tatsächlich gefahrene Schwelle `.{1000,}` — dokumentiert ist `.{N+1,}` (`:192`), was bei N = 1000 nicht läuft.
- **verifizierbar:** ja — zwei Läufe mit identischer Konfiguration über zwei Kopien (21 gegen 22) und die Gegenprobe am Minimalbaum.
- **klasse:** `kandidaten-menge-nicht-die-dokumentierte`

### F-5 · MEDIUM · Der Beanspruchungs-Commit reklamiert `MR-070` für einen Vorgang, den `MR-070`s Geltungsbereich ausdrücklich ausnimmt — zum dritten Mal

- **quelle:** `MR-070` (Feld *Geltungsbereich*); Reviewer-Anker 9
- **pfad:** Commit `f9161d96`, Botschaft, Absatz 2
- **befund:** Die Botschaft schreibt *„Dritter Lauf unter MR-070: Ausschluss-Menge vor dem Lauf gebildet, sie ist leer"*. `MR-070` nimmt genau diesen Fall heraus: *„**Nicht** erfasst: … und, ausdrücklich, der **Pfad-Nachzug eines Lifecycle-Moves** nach `MR-013`"* — mit eigener Begründung (der Eigenschafts-Test), die der Eintrag über einen halben Abschnitt führt. Dieselbe Etikettierung tragen bereits `971aaa3d` und `98e6c04c`; es ist die dritte Instanz. `MR-070` schreibt über seine erste Pflicht selbst: *„Gefangen wurde er dreimal von Pflicht 2 oder vom unabhängigen Review — **nie** von einer Vorab-Liste"*. Drei als *„Lauf unter MR-070"* etikettierte Lifecycle-Moves mit jeweils leerer Ausschluss-Menge lesen sich beim nächsten Trigger-Audit als Anwendungs-Belege für genau die Pflicht, die der Eintrag als unbelegt ausweist. Für die Traceability war der Verweis nicht nötig — `MR-013` steht ohnehin in derselben Botschaft.
- **verifizierbar:** nein (Commit-Botschaft, kein Gate-Gegenstand) — belegbar über `git log --grep=MR-070` und das Geltungsbereich-Feld des Eintrags.
- **klasse:** `mr-zitat-ausserhalb-seines-geltungsbereichs`

### F-6 · LOW · „kein Pin in einer getrackten Datei" — der Digest steht jetzt in einer getrackten Datei

- **quelle:** `MR-046` (Feld *Auflösungs-Trigger*)
- **pfad:** `docs/plan/planning/in-progress/slice-211-obermengen-nachweis-md013.md:146-147`; Commit `6088b616`, Botschaft, letzter Absatz
- **befund:** Die Gültigkeit von `MR-046` ruht auf drei Nein: kein Target, keine Deklaration, *„kein Pin in einer getrackten Datei"*. Das dritte ist wörtlich unzutreffend — `docs/plan/cr/…-zeilenlaenge.md:184` trägt `davidanson/markdownlint-cli2@sha256:173cb697…fed63a`, und DoD 2 verlangt ihn dort. Der Schluss trägt trotzdem, aber aus einem anderen Grund als dem genannten: der Digest hat keinen Konsumenten, kein Target zieht ihn, keine Nachtlauf-Achse wacht über ihn. In der wörtlichen Fassung ist die Prüfliste beim nächsten `MR-046`-Anlass unbrauchbar — wer sie abhakt, findet den Digest und muss entweder die Liste aufweichen oder den Eintrag ohne Sachgrund ablösen.
- **verifizierbar:** ja — `grep -n 'sha256:' docs/plan/cr/2026-09-06-cr-eingehend-auftraggeber-zeilenlaenge.md`.
- **klasse:** `nein-aussage-woertlich-unzutreffend`

### F-7 · INFO · Der herangezogene Kanon-Absatz regelt die **Ablösung eines bestehenden** Skripts; hier existiert keines

- **quelle:** `v6.5.0` · `regelwerk/modul-11-verification.md` §Fitness Function ohne Standard-Tool
- **pfad:** `docs/plan/cr/2026-09-06-cr-eingehend-auftraggeber-zeilenlaenge.md:254-256` (§Antwort auf Frage 2)
- **befund:** Der Absatz beginnt mit *„Ein selbstgebautes Gate ist auf Zeit gebaut … Erscheint später eines"* und endet mit *„Ist das Werkzeug Obermenge, wird das Skript retired — sonst benennt man die fehlende Klasse und behält es."* Er adressiert ein **vorhandenes** Gate und ein **später** erscheinendes Werkzeug. In slice-211 ist die Lage umgekehrt: das Werkzeug existiert, das Gate nicht — §1 des Plans schließt jede `.d-check.yml`-Änderung ausdrücklich aus, der `forbid-pattern`-Weg ist nie scharf geschaltet worden. Die **Drei-Teile-Form** und die **Break-Test-Form** übertragen sich unverändert; die Ablösungs-Klausel, mit der der Nachweis schließt, hat hier kein Objekt. Die Übertragung stammt aus der CR-Fassung vom 2026-09-06 (§Abgrenzung) und ist nicht in dieser Range entstanden — sie ist aber auch nirgends als Übertragung benannt.
- **verifizierbar:** nein — Urteil über den Geltungsbereich einer Kanon-Stelle.
- **klasse:** `kanon-absatz-auf-umgekehrte-lage-uebertragen`

### F-8 · INFO · K5 und K6a prüfen dieselbe MD013-Ausnahme; als zwei Klassen gezählt lassen sie die Zählung breiter aussehen, als sie ist

- **quelle:** Reviewer-Anker 17
- **pfad:** `docs/plan/cr/2026-09-06-cr-eingehend-auftraggeber-zeilenlaenge.md:209-210`
- **befund:** Gemessen enthält weder `k5.md:7` (lange URL) noch `k6a.md:7` (Inline-Code-Token) **irgendein** Whitespace; beide werden von MD013 aus demselben einen Grund übergangen — der Ausnahme für Zeilen ohne Leerzeichen jenseits der Schwelle. Zwei Zeilen der Ergebnistabelle, ein Mechanismus. Für `forbid-pattern` gilt dasselbe (beide melden). Die Zählung *„sechs von sechs gegen drei von sechs"* ist dadurch nicht falsch — beide Seiten werden gleich gezählt —, aber sie suggeriert sechs unabhängige Prüfachsen, wo fünf gemessen sind. Der CR selbst führt die beiden in **einem** Akzeptanzkriterium (*„eine lange URL oder ein langer Inline-Code"*), und der Nebenfund im selben Abschnitt sagt es sogar: *„«Inline-Code» ist keine Klasse, «unteilbares Token» ist eine."*
- **verifizierbar:** ja — `grep -c '[[:space:]]'` auf beiden Probenzeilen ergibt 0.
- **klasse:** `zwei-klassen-ein-mechanismus`

---

## Negativbefunde (geprüft, ohne Befund)

- **Proben-Isolation K1:** Die Zeile trägt Whitespace jenseits von Spalte 300 (`… und noch mehr text hier`); MD013 meldet aus dem beabsichtigten Grund, nicht zufällig. Die im Nachweis beschriebene Neuanlage ist an der Probe ablesbar.
- **Proben-Isolation K3 (Fenced Block):** Gegenprobe mit `code_blocks: true` meldet `k3.md:8:301` — das Schweigen kommt nachweislich aus der Fenced-Ausnahme, nicht daher, dass die Zeile die Schwelle verfehlt.
- **Proben-Isolation K4 (Tabellenzeile):** korrekt gebaut — Kopfzeile, Trennzeile, Datenzeile, führendes und schließendes `|`, Whitespace jenseits der Schwelle. Gegenprobe mit `tables: true` meldet `k4.md:9:301 … Actual: 347`; das Schweigen kommt aus der Tabellen-Ausnahme.
- **Proben-Isolation K7 (Skopus):** Die lange Zeile ist inhaltsgleich mit K1, liegt aber im zweiten Abschnitt. Gegenprobe mit auf `## Anderer Abschnitt` umskopierter Regel liefert `k7.md:7 section-forbidden` — `structure`s Schweigen kommt **nur** aus der Abschnitts-Skopierung. Teil 3 des Nachweises trägt damit sauber.
- **K2 (Summe über der Schwelle):** kein Befund. Die Probe ist per Konstruktion für beide Sensoren unbrechbar (RE2 `.` trifft keinen Zeilenumbruch; MD013 ist eine Zeilen-Regel) — ein Schweige-Test statt eines Break-Tests, aber ein ausdrückliches CR-Kriterium, und beide Seiten werden gleich behandelt. Nicht gemeldet.
- **Teil 2, RE2-Decke:** wörtlich reproduziert, inklusive der Aussage, dass der Fehler den **ganzen** Lauf mitnimmt (Exit 2, nicht Schweigen).
- **Off-by-one `.{1000,}` gegen MD013s „> 1000":** geprüft — im Bestand existiert keine Zeile von exakt 1000 Zeichen, die Differenz hat keine Wirkung auf 9 gegen 22.
- **„null Dateien meldet nur MD013":** reproduziert; die MD013-Menge ist echte Teilmenge. Auch über die 702-Datei-Kopie (mit `AGENTS.md`, `README.md`, Archiv) meldet MD013 dieselben 9.
- **Slice-Plan §8, drei Vorprüfungen:** vollständig, in der geforderten Reihenfolge; beide kanonischen Blöcke tragen `d-check:cite` auf die vorschreibende Zeile, wortgleich, mit denselben Spannen wie slice-207/208/209; der dritte Block (Nachtlauf, `MR-053`) trägt bewusst keine.
- **Slice-Plan §1, Abgrenzung:** alle vier Ausschlüsse eingehalten. `.d-check.yml`, `.d-check.closure.yml`, `Makefile`, `AGENTS.md` und `harness/` sind über die gesamte Range unberührt; kein Retrofit am Bestand; kein markdownlint-Artefakt im Baum (`git status` zu Beginn des Reviews leer).
- **`AGENTS.md` §3.1:** keine Spur eines Host-Skript-Interpreters; der Vergleich läuft über ein digest-gepinntes Image, denselben Weg wie `semgrep` und `trivy`.
- **`AGENTS.md` §3.3 / `MR-013`:** `f9161d96` ist ein reiner Rename plus der gekoppelte Roadmap-Marker; der Marker-Block ist die byte-exakte Umkehr des slice-210-Closure-Commits (Blob `b547cb7c` → `a7bfab5f`). Keine Inhaltsänderung an der Slice-Datei im Move-Commit.
- **`AGENTS.md` §3.7, Zustandsfeld:** die `Stand:`-Zeile des CR nennt Zustand (*„entschieden am 2026-09-07 — nicht umgesetzt"*) und zwei Belege (§Entscheid, §Obermengen-Nachweis). Der Satz über den §Zwischenstand kennzeichnet einen Abschnitt als überholt und erzählt nicht, wie der Zustand entstand — an der Grenze, aber kein Chronik-Befund.
- **`AGENTS.md` §3.4/§3.5:** kein Spec-Stratum und keine ADR berührt; keine Abwärts-Referenz entstanden.
- **`AGENTS.md` §3.2/§3.6/§3.9:** keine Suppression, keine Schwellen-Senkung, keine `uses:`-Referenz im Diff.
- **Reviewer-Anker 1–5, 10–16:** ohne Gegenstand — der Diff enthält keinen Produkt-Code, kein Modul, keine Konfiguration, keinen Gate-Pfad und keinen Kommentar in Code/Konfiguration/Skript.
- **DoD-Punkt 3 („alle drei Fragen beantwortet, `Stand:`-Zeile, Begründung je Frage"):** die drei Antworten sind getrennt geführt, jede mit eigener Begründung; die Warnung des Plans (§6), aus einem Ja auf Frage 2 kein Ja auf Frage 3 zu machen, ist auch in der Gegenrichtung eingehalten — Frage 3 wird **nicht** aus dem Nein zu Frage 2 abgeleitet, sondern für beide Wege getrennt begründet (Decke bzw. Preis).
- **Umkehr-Bedingungen 1 und 2:** beobachtbar formuliert (eine Bedingung im Modul `structure`; eine Node-Toolchain im Repo) — beide ohne Rückfrage entscheidbar. Nur Bedingung 3 ist es nicht (F-2).
- **Parallele Sitzung:** die während des Reviews aufgelaufene Fremd-Änderung betrifft einen anderen CR und keinen Gegenstand dieser Range; sie ist hier nur als Grund für die Zähl-Differenz 712/713/714 vermerkt.

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 5 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** `klassen-definition-nach-der-messung-verengt` · `zahl-aus-anderer-menge-als-die-aussage` · `ursachen-attribution-ohne-nachzaehlung` · `kandidaten-menge-nicht-die-dokumentierte` · `mr-zitat-ausserhalb-seines-geltungsbereichs` · `nein-aussage-woertlich-unzutreffend` · `kanon-absatz-auf-umgekehrte-lage-uebertragen` · `zwei-klassen-ein-mechanismus`

Vier der acht Klassen sind Spielarten **einer** Familie: eine Zahl oder eine Zusage steht neben einer Menge, die nicht die ist, über die geredet wird (F-1, F-2, F-3, F-4, F-8). Das ist `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand` bzw. `BEO-ALL/eigene-menge-gemessen-fremde-behauptet`. Bemerkenswert dabei: der Slice hat diese Klasse an **zwei** Proben selbst gefunden und behoben — die Regel aus §3 hat gewirkt. Sie hat nur nicht bis in die **Zahlen** durchgeschlagen, die aus den Proben gezogen wurden.

---

## Verdikt

**Merge-blockierend: ja** — fünf MEDIUM. Die Blockade betrifft ausdrücklich **nicht** den Entscheid: alle drei Antworten des CR halten der Nachmessung stand, und zwei Findings (F-2, F-3) machen die tragende Begründung sogar robuster, als sie geschrieben ist. Zu klären ist die **Beleglage**: die Klassen-Definition, die die Zählung trägt (F-1), die Herkunft der Zahlen 2045 / „eine über 2000" (F-2), die Ursachen-Zuschreibung der 13 (F-3) und die Reproduzierbarkeit der 22 (F-4). Solange sie stehen, behauptet der Nachweis an vier Stellen mehr Deckung, als er gemessen hat — bei einem Dokument, dessen einziger Zweck es ist, eine Entscheidung **auf Messung statt auf Datenblatt** zu stellen, ist genau das der blockierende Punkt.

**Übergabe:** Findings an den Implementer; die Finding-Klassen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist Lauf-Beleg und ersetzt die Verifikation nicht — DoD- und Spec-Konformität prüft der Verifier separat, insbesondere DoD 1 (der Kanon-Test *„der unveränderte Bestand, auf dem beide schweigen"* ist für `forbid-pattern` nicht gefahren worden, sondern durch einen Differenz-Lauf ersetzt) und DoD 2 (F-4).
