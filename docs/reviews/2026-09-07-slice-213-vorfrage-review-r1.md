# Review-Report — slice-213, Runde 1 (Vorfrage: trägt §6 den Ableiter?)

**Review-Art:** Design/Doc — geprüft werden die **Vorfrage-Antwort** gegen den **Wortlaut** von `AGENTS.md` §6, die **Zahl**, auf der die neue Regel steht, gegen die sieben Fundstellen aus slice-212, die **Selbstanwendung** der Regel auf ihre eigene Grenzen-Aussage, und die neue `state.md` gegen die Zustandsfeld-Form. **Nicht** geprüft: die DoD-Abhakung und die Gate-Lauf-Bestätigung (Verifier-Rolle).
**Gegenstand:** `8a3e823a` (Vorprüfungen) · `d9e05b84` (Beanspruchung) · `175413a1` (DoD 1+2+3). Diff-Range `HEAD~3..HEAD`, `HEAD == 175413a1` zu Beginn und am Ende des Reviews; Arbeitsbaum sauber.
**Skill:** `.harness/skills/reviewer.md` v1.15.0 @ `175413a1` — einschließlich Anker 8 (Botschaft über die Messung hinaus), Anker 9 (Quelle über den Geltungsbereich) und Anker 17 (Messung zählt einen Proxy).
**Modell-ID:** claude-opus-5[1m] · **Datum:** 2026-09-07
**Eingangs-Kontext:** der Slice-Plan slice-213 (§1, §2, §3, §4, §6, §8, DoD 1–3); der Vorgänger slice-212 samt Closure-Notiz §7 und seinem Review-Report (Runde 1, sieben Findings — **F-1**, **F-5**, **F-6** und **I-1** sind für diesen Lauf Vorbefunde am selben Gegenstand); der Registereintrag `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen` mit `observation.md`, `state.md` und drei Evidence-Dateien; die Nachbar-Einträge `citation-stretched-beyond-scope`, `rule-drawn-from-occasion-not-inventory`, `eigene-menge-gemessen-fremde-behauptet`, `zaehlmethode-misst-proxy-statt-gegenstand`, `module-promise-only-on-scan-axis`; `AGENTS.md` §3.1/§3.3/§3.6/§3.7/§3.8/§4/§5/§6; die sieben Sensor-Dateien `baseline-verify`, `adr-check`, `doc-check`, `lint`, `arch-check`, `semgrep`, `review-coverage` im Stand **vor** slice-212 (`9771a354`) und heute; `harness/conventions.md` samt MR-013, MR-053, MR-054; Baseline `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice und §Zwei Schritte vor der Modus-Begründung, `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register, `regelwerk/grundlagen-traceability.md` §Herkunfts-Anker.

**Eigene Läufe (echte Ausgaben, keine behaupteten):**

- `make gates` → `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; darin `coverage-gate: OK — Coverage 94.60% erfüllt Schwelle 93%`, `Ran 55 rules on 63 files: 0 findings.`, `d-check: 726 Datei(en) geprüft, 0 Befund(e)` (targets- und planning-Lauf je einmal).
- `make doc-check` → `d-check: 726 Datei(en) geprüft, 0 Befund(e)`.
- Auszählungen (kein Gate zählt das): `ls .../observations/*/*/ | wc -l` ⇒ **38** Verzeichnisse; `evidence/`-Dateien: `grenzen-liste-wird-als-vollstaendig-gelesen` **3**, `citation-stretched-beyond-scope` **18**, `eigene-menge-gemessen-fremde-behauptet` **13**, `rule-drawn-from-occasion-not-inventory` **9**.
- `git show --stat --find-renames d9e05b84` ⇒ die Slice-Datei mit `| 0` — reiner Move.

## Was gemessen reproduziert (vorweg, damit die Findings nicht als Gesamturteil gelesen werden)

| Behauptung | Ort | Ergebnis |
|---|---|---|
| Register **38** Verzeichnisse, vier einschlägig | slice-213 §8, Botschaft `8a3e823a` | **exakt** — 37 unter `BEO-ALL/`, 1 unter `BEO-HARN/` |
| `citation-stretched-beyond-scope` 18×, `eigene-menge-gemessen-fremde-behauptet` 13× | slice-213 §8 | **beide exakt** — Evidence-Dateien nachgezählt |
| §6 trägt den **zweiten** Ableiter-Teil wörtlich | slice-213 §3 | **trägt** — §6 richtet den fremden Leser für *jedes* Artefakt eines Slice ein; eine zweite Regel daneben wäre eine zweite Quelle. Das **Urteil** hält; seine **Begründung** nicht (F-1) |
| §6 trägt den **ersten** Teil nicht | slice-213 §3 | **trägt er nicht** — die acht Schritte kennen keinen Vertrags-Umkehr-Schritt, Schritt 5 meint einen Sensor. Die *Reichweite* der Aussage ist zu groß (F-4) |
| „zwei Grenzen … beschrieben Mechanismen, die es nicht gibt" | `AGENTS.md:562` | **exakt** — `review-coverage` (Substring-Match statt Gleichheit) und `arch-check` (`edges` ist Erlaubnisliste) |
| „in einem Fall widersprach der Code-Kommentar dem Verhalten" | `AGENTS.md:563` | **exakt** — `review-coverage`, in slice-212 §7 belegt |
| `make gates` grün, zehn Gates, 726 Dateien | alle drei Botschaften | **exakt**, eigener Lauf oben |
| Reiner Move-Commit, Roadmap-Flip und drei Pfad-Verweise mitgezogen | `d9e05b84` | **exakt** — `| 0`, Ruhe-Marker entfernt, `planning-check` grün |

## Findings

### F-1 · MEDIUM · Die Grenze der neuen Regel schreibt dem Leser, den §6 einrichtet, sieben Funde zu — im Anlass-Fall sagt der Beleg wörtlich das Gegenteil

- **kategorie:** MEDIUM (Anker 8 + Anker 9; die Aussage ist die **tragende Begründung** dafür, Teil 2 des Ableiters nicht zu schreiben)
- **quelle:** `AGENTS.md` §5 (*„behauptet nicht mehr, als die Arbeit trägt"*) · `BEO-ALL/citation-stretched-beyond-scope` (18×)
- **pfad:** `AGENTS.md:564-567` · dieselbe Aussage in slice-213:104-106 und in der Botschaft `175413a1`
- **befund:** Die Grenze lautet: *„den [fremden Leser] richtet §6 ein, und **er** hat in allen sieben Fällen gefunden, was hier vermieden werden soll."* §6 richtet **den Reviewer und den Verifier** ein — sonst niemanden. Für den ersten der sieben Fälle sagt der eigene Beleg das Gegenteil: `evidence/slice-212.md` schreibt *„**Gefunden hat es der Auftraggeber** beim Lesen eines Zwischenbescheids, **nicht ein Review und kein Gate**"*, und slice-212 §7 zieht daraus ausdrücklich den Schluss *„Der fremdeste Leser war diesmal **kein Reviewer**."* Der Satz macht aus dem, was slice-212 als *„in keinem der sieben Fälle war es der Autor"* geschrieben hat, ein *„der von §6 eingerichtete Leser hat es gefunden"* — zwei verschiedene Mengen. **Versagensszenario:** Ein späterer Lauf stellt dieselbe Vorfrage für eine benachbarte Klasse, liest in §5 die eingefrorene Bilanz *„§6s Leser fängt das zuverlässig"* und verzichtet aus demselben Grund auf eine Regel — auf einer Bilanz, die genau den Fall verschweigt, in dem der Mechanismus **nicht** gegriffen hat und ein Zufallsleser einsprang. Verschärfend: Es ist die Klasse, gegen die der Slice sich in §3 und §8 selbst gewappnet hat (*„der Beleg muss aus dem Wortlaut kommen"*), angewandt auf §6s Wortlaut in der einen Richtung und unterlassen in der anderen.
- **verifizierbar:** ja — `grep -n "Auftraggeber" docs/plan/planning/observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/evidence/slice-212.md` und slice-212 §7. Kein Gate liest Prosa gegen Prosa.
- **klasse:** `fund-dem-falschen-mechanismus-zugeschrieben`

### F-2 · MEDIUM · „In sieben von sieben Fundstellen" steht ohne die Form, die drei Absätze weiter unten in derselben Sektion verlangt wird — und einer der sieben trägt sie nicht

- **kategorie:** MEDIUM (Anker 17; die Zahl trägt hier den **ganzen** Nachweis der Regel, nicht nur eine Illustration)
- **quelle:** `AGENTS.md` §5 (*„Vor einer Messung steht die Form ihres Gegenstands … **Sie steht dort, wo die Zahl steht**"*, `seit slice-210`)
- **pfad:** `AGENTS.md:560-562`
- **befund:** Der Absatz sagt *„Gemessen: In sieben von sieben Fundstellen stand die fehlende Grenze bereits im Vertrags-Teil oder in der Konfiguration"* und definiert nirgends, was einen Treffer zu einem Treffer macht — der Nachbarabsatz elf Zeilen später verlangt genau das, **im selben Artefakt, das die Aussage trägt**. Die Form fehlt nicht folgenlos: Bei `review-coverage`, dem siebten Fall, stand die fehlende Grenze **weder** im Vertrags-Teil **noch** in der Konfiguration. Der Vertrags-Teil sagte *„Substring-Match, 1:N zulässig"* — und das ist falsch; die echte Grenze (*„der Abgleich vergleicht die erste Kennung im Dateinamen, und nur sie"*) ist am **Code** abgelesen, was slice-212 §7 wörtlich festhält: *„Meine Grenze erbte den Fehler."* Die vier übrigen Klassen halten (`baseline-verify`: *„gegen `SHA256SUMS`"* im Vertrag; `semgrep`: *„gepinntes Regelset"* im Vertrag; `adr-check`/`doc-check`/`lint`/`arch-check`: in `.d-check.yml`, `.golangci.yml`, `.a-check.yml`). **Versagensszenario:** Ein Autor liest *„sieben von sieben"* als Zusage, dass der Vertrags-Umkehr-Schritt die Lücke immer zutage fördert, dreht den Vertrag um und schreibt — wie slice-212 es tat — eine Grenze ab, die es nicht gibt; die Zahl, die ihn dazu bringt, enthält den einen Fall, der das widerlegt. Ohne ausgeschriebene Form ist die Differenz nicht einmal benennbar.
- **verifizierbar:** ja — `git show 9771a354:harness/sensors/review-coverage.md` gegen `harness/sensors/review-coverage.md` §Grenze Punkt 5. Kein Gate prüft eine Zahl gegen ihre Menge.
- **klasse:** `zahl-ohne-form-des-treffers`

### F-3 · MEDIUM · Die einzige Urteils-Regel in §5 ohne gepaarten Reviewer-Anker — während ihre eigene Grenze das Fangen an den Reviewer delegiert

- **kategorie:** MEDIUM (Anker 15-nah: eine Zusage wird über eine Achse gegeben, auf der niemand misst)
- **quelle:** `AGENTS.md` §3.7, §3.8, §5 (drei Vorbilder) · Baseline `v6.5.0` · `regelwerk/modul-11-verification.md` §Fitness Function ohne Standard-Tool (*„Hard Rule in zwei Quadranten"* — inferential feedforward **und** eine prüfende Gegenseite)
- **pfad:** `AGENTS.md:567` (Ende des neuen Absatzes) gegen `AGENTS.md:304`, `AGENTS.md:332`, `AGENTS.md:357`, `AGENTS.md:584`, `AGENTS.md:591`, `AGENTS.md:602`
- **befund:** Jede andere Urteils-Regel dieser Datei schließt mit *„der Reviewer-Skill trägt den Anker dazu"* — §3.7 (zwei), §3.8, und in §5 die drei Nachbarn `seit slice-210`, `seit welle-82`, `seit slice-147`. Der neue Absatz endet mit *„Urteil, kein `grep`."* und bricht ab; `grep -n "Grenzen-Liste" .harness/skills/reviewer.md` ⇒ **kein Treffer**, die sechzehn Prüffragen kennen die Klasse nicht. Zugleich stützt sich die Grenze desselben Absatzes ausdrücklich auf den fremden Leser (*„ersetzt den fremden Leser nicht — den richtet §6 ein"*). Damit ist die Feedforward-Hälfte geschrieben und die Feedback-Hälfte an eine Rolle abgegeben, der niemand gesagt hat, wonach sie sehen soll. Der Slice schließt in §1 *„ein Sensor darauf"* aus — ein Skill-Anker ist keiner, und die Abgrenzung deckt seine Abwesenheit deshalb nicht; der Plan trifft die Entscheidung schlicht nicht. **Versagensszenario:** Der nächste Slice ergänzt einen `## Grenze`-Abschnitt, der Autor lässt die Vertrags-Umkehrung aus (kein Sensor, keine Spur), der Reviewer arbeitet seine sechzehn Fragen ab, unter denen die Klasse nicht vorkommt, und die Lücke geht durch — der Ablauf, den der Eintrag dreimal belegt hat. Der Registereintrag steht dann auf *verkörpert*, ohne dass sich am Erkennungspfad etwas geändert hätte.
- **verifizierbar:** ja — `grep -n "Reviewer-Skill trägt" AGENTS.md` (sechs Treffer, keiner im neuen Absatz) und `grep -n "Grenzen-Liste" .harness/skills/reviewer.md` (leer). Kein Gate paart Regel und Anker.
- **klasse:** `regel-ohne-feedback-haelfte`

### F-4 · MEDIUM · „Der erste Teil steht nirgends" ist über drei Sektionen gemessen und über die ganze Datei behauptet — §3.8 verlangt die Zusagen-Umkehrung bereits

- **kategorie:** MEDIUM (Anker 8; die Negativ-Aussage trägt DoD (1) und damit die Existenz der neuen Regel)
- **quelle:** `AGENTS.md` §5 (*„ihr Schluss reicht nicht weiter als die gemessene Menge"*) · `BEO-ALL/eigene-menge-gemessen-fremde-behauptet` (13×, von slice-213 §8 selbst als einschlägig geführt)
- **pfad:** slice-213:33 (*„**Der erste Teil steht nirgends**"*) und slice-213:110-113 (*„§5 regelt, was eine Aussage behaupten darf; §3.7, was ein Kommentar trägt. **Keine der drei sagt**, was man tut, bevor man eine Grenzen-Liste aus der Hand gibt."*)
- **befund:** Geprüft sind ausweislich des Plans **§6, §5 und §3.7**; behauptet ist *nirgends*. Ungeprüft blieb der nächste Nachbar: §3.8 sagt *„Wer ein Modul anlegt oder ändert, beantwortet deshalb: welche Eingaben liest es, die es nicht scannt — und gilt dort dieselbe Zusage? **Wo sie nicht gilt, gehört die Grenze in die Anforderung**"* — das ist eine Zusage umdrehen und die Kehrseite aufschreiben, verpflichtend beim Schreiben, für Module auf einer Achse. Die neue Regel ist echt breiter (jedes Artefakt mit Vertrags-Teil, plus die Code-Hälfte), also fällt die Entscheidung *Nein* nicht — aber ihr Verhältnis zu §3.8 ist ungeklärt, und dieselbe Datei trägt jetzt eine allgemeine Regel in §5 und ihren Spezialfall als Hard Rule in §3.8, ohne dass eine die andere nennt. **Versagensszenario:** Genau das Argument, mit dem der Slice Teil 2 verwirft — *„eine zweite Regel daneben wäre eine zweite Quelle für dieselbe Aussage"* —, trifft nun die Paarung §5/§3.8: Wer §3.8 später anfasst (sein Registereintrag steht nach der `v6.0.0`-Korrektur **unter** der Schwelle, also ist Anfassen absehbar), liest die Achsen-Frage als die ganze Pflicht und weiß nicht, dass §5 sie inzwischen verallgemeinert; wer §5 anfasst, sieht §3.8 nicht.
- **verifizierbar:** ja — `sed -n '/^### 3.8/,/^### 3.9/p' AGENTS.md` gegen slice-213:110-113. Kein Gate prüft Negativ-Aussagen über eine Datei.
- **klasse:** `negativ-behauptung-ueber-ungemessene-menge`

### F-5 · MEDIUM · Die neue `state.md` holt „zwei der drei Belege" zurück — genau die Hälfte, die der slice-212-Review (F-6) entfernen ließ

- **kategorie:** MEDIUM (LOW-Klasse aus dem Vorgänger-Report, **eine Stufe höher wegen Wiedereinführung nach dokumentierter Entfernung**)
- **quelle:** `AGENTS.md` §3.7 · Baseline `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register (*„Der Zähler wird abgeleitet, nicht geführt … Ein gespeicherter Zähler neben einer Belegliste sind zwei Quellen für denselben Zustand"*)
- **pfad:** `docs/plan/planning/observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md:1` (*„**Bleibt wach:** **zwei der drei** Belege sind nachgetragen"*)
- **befund:** Der slice-212-Report meldete als F-6, dass die `state.md` *„Die Schwelle ist mit dem dritten Beleg erreicht; zwei der drei sind nachgetragen"* einen abgeleiteten Stand speichert; `f36c6a1e` entfernte den Satz, der Stand bei `HEAD~3` trug ihn nicht mehr. `175413a1` schreibt *„zwei der drei Belege sind nachgetragen"* zurück. **Versagensszenario:** wortgleich das aus F-6 und von slice-213 §6 selbst vorhergesagt — *„Eine vierte Instanz wäre der bessere Beleg als eine dritte nachgetragene"*: Sobald `evidence/` vier Dateien führt, sagt das Zustandsfeld *„zwei der drei"*, während der abgeleitete Zähler bei vier steht; das Feld behauptet einen Stand, den es nicht besitzt, und kein Sensor meldet es. Die tragende Information (*welche* Belege nachgetragen sind) steht bereits in den beiden Evidence-Dateien selbst.
- **verifizierbar:** ja — `git show f36c6a1e:…/state.md` gegen den heutigen Stand; reproduzierbar durch Anlegen einer vierten Evidence-Datei. Kein Gate prüft Zustandsfeld-Form.
- **klasse:** `state-md-speichert-abgeleiteten-stand`

### F-6 · MEDIUM · Die Grenze der Regel nennt ihre Belegmenge nicht — `state.md` und Commit-Botschaft tun es, der Nachbarabsatz führt sie in seiner Grenze

- **kategorie:** MEDIUM (die Regel verfehlt an ihrer eigenen Grenzen-Aussage genau das, was sie vorschreibt — Anker 8)
- **quelle:** `AGENTS.md` §5 · `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen` (der Gegenstand selbst)
- **pfad:** `AGENTS.md:564-567` gegen `state.md:1` und die Botschaft `175413a1`
- **befund:** `state.md` und Botschaft halten fest: *„die Regel ist an **Sensor-Beschreibungen** belegt, nicht an Grenzen-Listen überhaupt."* Die Grenzen-Aussage in §5 sagt das nicht; sie nennt nur die Abgrenzung zum fremden Leser. Der Nachbarabsatz (`seit slice-210`) zeigt die Form, die hier fehlt: *„belegt ist sie nur an Doku-/Planning-Messungen."* Zweitens fehlt die Vorbedingung, die slice-212 teuer gelernt und der dortige Report als I-1 notiert hat: Die erste Hälfte der Regel setzt voraus, dass der Vertrags-Teil **stimmt**; die zweite Hälfte fängt das nur ab, **wo der Gegenstand Code oder Konfiguration ist** — für einen Gegenstand, der selbst Prosa ist (eine Regel, eine Kanon-Stelle, ein Slice-Plan), steht die erste Hälfte allein und erbt Vertragsfehler ungebremst. **Versagensszenario:** Ein Autor wendet die Regel auf die Grenzen-Liste eines Prosa-Artefakts an, dreht dessen Vertrags-Teil um, erbt einen dort stehenden Fehler und hat, anders als bei Code, keine zweite Instanz zum Gegenprüfen — die Regel sagt ihm nicht, dass sie für diesen Fall unbelegt ist. Das ist die Klasse des Eintrags, angewandt auf den Absatz, der ihn verkörpert: Die Aufzählung liest sich als Menge und ist eine Auswahl.
- **verifizierbar:** nein — Urteil; die Belege sind `state.md`, die Botschaft und der Nachbarabsatz, alle drei zitierbar.
- **klasse:** `grenzen-liste-laesst-die-eigene-einschraenkung-weg`

### F-7 · LOW · Zwei Zähler-Stände für denselben Registereintrag im selben Plan — zum dritten Mal in drei aufeinanderfolgenden Slices

- **kategorie:** LOW (die Zahl trägt hier keine Schlussfolgerung; beide Werte liegen weit über der Schwelle) — **mit Steering-Loop-Vermerk**, siehe unten
- **quelle:** `AGENTS.md` §5 · Baseline `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register (*„Der Zähler wird abgeleitet, nicht geführt"*)
- **pfad:** slice-213:13-14 (§Bezug, *„(10×)"*) gegen slice-213:216-217 (§8, *„(9×)"*)
- **befund:** `BEO-ALL/rule-drawn-from-occasion-not-inventory` trägt **neun** Evidence-Dateien; der Plan nennt den Stand zweimal und einmal falsch. Die Kopfzeile stammt aus `03ca7b1c` (slice-212) und liegt damit vor dieser Range — aber `f36c6a1e` korrigierte alle vier Stellen in slice-212 auf 9× und ließ die Schwesterdatei stehen, und `8a3e823a` schrieb den richtigen Wert in §8 **desselben** Dokuments, im Commit, dessen Gegenstand die Register-Sichtung ist. **Versagensszenario:** wie im Vorgänger-Report — der nächste Lauf sichtet das Register, findet 9, hält den Plan für den frischeren Stand und schreibt einen zehnten Beleg, den es nicht gibt. **Steering-Loop-Vermerk:** Dieselbe Klasse ist jetzt dreimal berichtet (slice-210 F-5, slice-212 F-5, hier) und zweimal an derselben Ursache — ein abgeleiteter Zähler, der in Prosa kopiert wird. Zusammen mit F-5 dieses Reports sind das in **einem** Slice zwei Instanzen derselben Wurzel; nach der Kontext-Eskalation des Skills ist das ein Register-Kandidat, kein Einzelbefund mehr.
- **verifizierbar:** ja — `ls docs/plan/planning/observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/evidence/ | wc -l` ⇒ 9. Kein Gate zählt Belege gegen Prosa.
- **klasse:** `zwei-zaehler-staende-fuer-denselben-eintrag`

## Negativbefunde (geprüft, ohne Befund)

- **Die Vorfrage-Antwort selbst (DoD 1), Ja-Hälfte.** §6 richtet den fremden Leser für **jedes** Artefakt eines Slice ein und verbietet den Self-Review; eine zweite Regel *„die Grenzen-Liste braucht einen fremden Leser"* wäre eine zweite Quelle für dieselbe Aussage. Das ist **keine** Dehnung: Der Schluss geht vom allgemeinen Satz auf einen Fall, den sein Wortlaut umfasst, nicht von einem Fall auf einen allgemeinen Satz. Geprüft habe ich die Gegenrichtung mit: Nichts in §6 nimmt eine Artefaktklasse aus. Beanstandet ist nur die **Begründung** (F-1) und die fehlende Feedback-Hälfte (F-3), nicht das Urteil.
- **Die Vorfrage-Antwort, Nein-Hälfte.** §6s acht Schritte enthalten keinen Vertrags-Umkehr-Schritt; Schritt 5 (*„Engsten nützlichen Sensor laufen lassen"*) meint einen ausführbaren Sensor, Schritt 6 den Gate-Lauf, Schritt 7 die Doku-Nachführung. Auch der Reviewer-Skill enthält nichts, was den **Autor vor der Übergabe** bindet — er adressiert den Reviewer. Die Verengung ist nicht künstlich; ihre **Reichweite** ist es (F-4).
- **§1-Abgrenzung eingehalten.** Alle drei Ausschlüsse sind unberührt: keine erneute Inventur über die Sensor-Dateien (die Range fasst `harness/sensors/` nicht an), kein Sensor angelegt, kein `adr-check`-Bruch-Test. Die Range ändert fünf Dateien, alle im Plan angekündigt oder von MR-013 gefordert.
- **§3.3 / MR-013.** `d9e05b84` ist ein reiner Move (`| 0`, `--find-renames` bestätigt) und bündelt Roadmap-Flip und die drei Pfad-Verweise; der Review-Report von slice-212 bleibt zu Recht unangetastet (Inline-Code, kein Link). `make planning-check` grün.
- **§8, alle drei Vorprüfungen.** Beide `d-check:cite`-Direktiven zeigen auf die **vorschreibenden** Zeilen (`modul-05-planning-harness.md:268-269` und `:274`) und tragen ihre Zitate wortgleich — das Modul `citations` läuft fail-closed im inneren Loop und ist grün. Die Form entspricht slice-210 und slice-212 byteweise. Der dritte Block (Nachtlauf) trägt bewusst keine Direktive, wie MR-053/MR-054 es vorsehen. Register-Zahl **38** nachgezählt und korrekt.
- **Anker-Paarung.** `state.md` nennt als Zielort `AGENTS.md` §5, und §5 trägt `seit slice-213` im neuen Absatz; der Zielort löst ab Repo-Wurzel auf, der §6-Verweis darin (`#6-minimal-agent-workflow`) löst auf (`anchors` grün). Dass der Anker im Absatz statt in der §5-Überschrift steht, ist die etablierte Haus-Form dieser Datei — `seit welle-82`, `seit slice-147` und `seit slice-210` stehen genauso; kein Befund gegen diesen Slice.
- **Zustandsfeld-Form der `state.md` im Übrigen (§3.7).** *„verkörpert"* ist der Zustand, `AGENTS.md` §5 mit `seit slice-213` der auflösbare Beleg. Die Einschränkung *„nur der erste Teil des Ableiters"* ist **Zustand**, keine Chronik: Sie sagt, was heute verkörpert ist und wo der Rest liegt, nicht, wie es dazu kam. Form und Register wie beim Nachbarn `zaehlmethode-misst-proxy-statt-gegenstand`. Beanstandet ist allein der gespeicherte Stand (F-5).
- **Commit-Botschaften gegen die Arbeit (§5).** Die drei Gate-Angaben stimmen mit meinen eigenen Läufen überein (zehn Gates, 726 Dateien, 0 Befunde; `make doc-check` 726/0). Die Botschaften behaupten keine Probe, die nicht lief, und benennen die Einschränkungen (*„nur der erste Teil"*, *„zwei der drei Belege sind nachgetragen"*, *„an Sensor-Beschreibungen belegt"*) ausdrücklich mit. Die einzige überdehnte Aussage ist die aus F-1, und sie steht in Regel und Botschaft gleichlautend — sie ist dort kein zusätzlicher Befund.
- **Bestands-Effekt in §5.** Der neue Absatz überlappt **nicht** mit seinen beiden Nachbarn: *„Vor einer Messung steht die Form ihres Gegenstands"* prüft die Messung, *„behauptet nicht mehr, als die Arbeit trägt"* den Schluss, der neue Absatz die **Vollständigkeit einer Aufzählung**. Drei verschiedene Gegenstände; die Abgrenzung zum Mess-Absatz ist dort sogar explizit ausgeschrieben. Die einzige ungeklärte Nachbarschaft ist §3.8 (F-4).
- **Hard-Rule-Form.** Der Absatz trägt Herkunfts-Anker, Registerzeiger und den Auflösungs-Trigger *permanent* — dieselbe Form wie die drei Vorbilder in §5.
- **`make gates` und `make doc-check`.** Beide grün, Ausgaben oben. Keiner der sieben Befunde ist gate-sichtbar, und keiner ist es der Natur nach: Es geht durchweg um Aussagen über Aussagen.

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 6 | F-1, F-2, F-3, F-4, F-5, F-6 |
| LOW | 1 | F-7 |
| INFO | 0 | — |

## Verdikt

**Blockierend.** Das Urteil des Slice — §6 trägt Teil 2, Teil 1 wird geschrieben — hält der Prüfung am Wortlaut stand und ist die richtige Antwort; blockierend sind nicht die Entscheidung, sondern die drei Stellen, an denen die neue Regel selbst tut, wogegen sie geschrieben ist: Sie schreibt einen Fund dem falschen Mechanismus zu (F-1), stützt sich auf eine Zahl ohne Form, deren siebter Fall sie widerlegt (F-2), und lässt in ihrer eigenen Grenzen-Liste die Einschränkung weg, die ihr Autor zweimal daneben aufgeschrieben hat (F-6) — dazu die fehlende Feedback-Hälfte (F-3), die überdehnte Negativ-Aussage (F-4) und die zurückgeholte Zähler-Kopie (F-5).
