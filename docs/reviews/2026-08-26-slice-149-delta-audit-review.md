# Review-Report: slice-149 — Etappe B, Delta-Audit `v5.12.0` und Freshness-Audit aller Adaptionen

**Datum:** 2026-08-26 · **Review-Art:** Plan-/Doku-Review (geprüft gegen Slice-Plan slice-149 §1/§2/§3/§4/§5/§7/§9, Wellendokument welle-85 §1/§3/§4/§6, den **vendorten** Kanon `v5.12.0` — `modul-02` §Freshness-Audit, `modul-05` §Lifecycle als State Machine / §Closure- und Lerneintrag-Regeln / §Offene Risiken / §Ziel-Form: Slice, `modul-06` §Das Beobachtungs-Register, `modul-13` §Fitness Function, `modul-03` §Ziel-Form: Architektur-Sicht, `grundlagen-source-precedence.md`, `grundlagen-harness-dateien.md` §Konventionsspeicher, `templates/AGENTS.template.md` —, `AGENTS.md` §1/§3/§5, **alle 16** aktiven `MR-`Einträge, Beobachtungs-Register `BEO-008`/`BEO-009`/`BEO-011`/`BEO-012`/`BEO-015`/`BEO-016`/`BEO-017`), unabhängiger Reviewer ohne Anteil an der Arbeit
**Gegenstand:** Commit-Kette `297a58a..HEAD` — `ac1a23b` (Move + Etappe-C-Schnitt), `3bdc94f` (Closure-Body), `cdc54c5` (Tatsachenberichtigung); **8 Dateien, 418 Einfügungen / 113 Löschungen** (`git diff --stat 297a58a..HEAD`)
**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 (`9ee805b`) · **Modell-ID:** `claude-opus-5[1m]`

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/done/slice-149-baseline-v5120-delta-audit.md` (vollständig, inkl. §3 „Ausdrücklich NICHT", §4 DoD, §5 Risiko-Ausgänge, §9)
- Wellendokument `docs/plan/planning/welle-85-baseline-v5120-migration.md`
- Die zwei geschnittenen Slices `docs/plan/planning/open/slice-150-pin-gebundene-zitate.md`, `docs/plan/planning/open/slice-151-urteilsfreie-haelfte-voll.md`
- Der vendorte Kanon `.harness/baseline/v5.12.0/` (Regelwerk + Templates), gelesen als **Datei**, nicht als Antwort des Kurses
- `docs/plan/cr/2026-08-26-antwort-regelwerk-v5110.md` (die Erwartung, gegen die das Delta gehalten wurde)
- `harness/conventions.md` und **alle 16** aktiven Einträge unter `harness/conventions/`
- `docs/plan/planning/observations.md` — `BEO-008`, `BEO-009`, `BEO-011`, `BEO-012`, `BEO-015`, `BEO-016`, `BEO-017` (neu)
- **Vorherige Findings am gleichen Vorgang:** `docs/reviews/2026-08-26-slice-148-baseline-v5120-review.md` (Etappe A, zwei MEDIUM zu falsch gerahmten Zahlen) und `docs/reviews/2026-08-23-slice-129-delta-audit-review.md` (das vorige Delta-Audit, dieselbe Slice-Art)

**Nicht erhalten** (Skill §Eingangs-Kontext): die DoD-Abhakung als Prüf-Gegenstand — Plan-/DoD-Konformität prüft die Verifikation. Der offene DoD-Haken wird hier **nur** dort berührt, wo er eine Aussage des Kanons über die Verzeichnis-Position trägt (F-7).

**Vom Reviewer selbst gefahren** (Exit je Lauf in eine Datei umgeleitet und direkt gelesen):

- `make gates` Exit **0** — zehn Glieder (`baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check`), `targets`/`planning` je 509 Dateien / 0 Befunde, Coverage 94,80 %.
- `make fullbuild` Exit **0** — Abschluss `image-hash sha256:f41f99ad9a1e6d256423de4cf3aec77586e16fcedc3aac7d2d1cee408dfafc72`, Closure-Profil (`planning` + `structure`) 467 Dateien / 0 Befunde, 48 Anforderungen / 0 Waisen.
- **Beide gepinnten Bäume rekonstruiert und die fünf Inhalts-Dateien Zeile für Zeile verglichen** (`git show 9ee805b^:.harness/baseline/v5.11.0/<f>` gegen `git show HEAD:.harness/baseline/v5.12.0/<f>`, `diff -u6`) — der vendorte Text ist die Grundlage jeder Zuschreibung unten, nicht die Kurs-Antwort.
- **Eigene Link-/Anker-Probe** über alle 38 Dateien des Konventionsspeichers: jeder Markdown-Verweis in den vendorten Baum, Ziel-Existenz und Anker gegen Überschriften **und** `<a id>` aufgelöst (Python, Scratchpad).
- **Eigene Datums-Probe** über `git log` je Beleg-Slice von `BEO-017` (Anlage-, Feat-, Review-, Move- und Closure-Commit).
- Zählungen: `done/`-Slice-Dateien am Schnitt-Commit `ac1a23b`, `Ersetzt-Baseline-Regel`-Ziele aller 16 aktiven Einträge gegen die fünf geänderten Dateien.
- **Keine Repo-Datei wurde verändert** — kein `git checkout --` nötig, `git status --short` nach jedem Lauf leer. Die einzige geschriebene Datei ist dieser Report.

**Verdikt: blockierend** — kein HIGH, sieben MEDIUM, vier LOW.

---

## Findings

### F-1

- **kategorie:** MEDIUM
- **quelle:** `.harness/baseline/v5.12.0/regelwerk/modul-02-harness-bootstrap.md` §Freshness-Audit (*„Der Review geht durch die **Adaptions-Liste**, nicht nur durch den Diff — Frage pro Eintrag: Regelt die neue Fassung das, wofür diese Adaption angelegt wurde?"*, dann **fünf Ausgänge**) · welle-85 §3 Trigger 2 · `BEO-011`
- **pfad:** `docs/plan/planning/done/slice-149-baseline-v5120-delta-audit.md:68-70` (DoD) und `:138-146` (§9)
- **befund:** Die „gemessene" Verengung fragt nicht die Frage des Kanons (*Gegenstand* des Eintrags), sondern eine Eigenschaft seines **Zeigers** — ob das `Ersetzt-Baseline-Regel`-Ziel unter den fünf geänderten Dateien liegt. Für **vier** der 16 Einträge nennt dieses Feld gar keine Baseline-Datei (`MR-034` verweist auf `MR-006`, `MR-035`/`MR-036` tragen „keine", `MR-037` trägt „—"), sodass Schritt 1 für sie unter **keinem** Delta „betroffen" liefern kann; `MR-035`/`MR-036` haben als Auflösungs-Trigger ausgerechnet *„der Kanon benennt selbst einen Ruheort für ausgehende Change Requests"* — genau die Änderung, die der Filter nicht sähe. Das Dokument nennt zudem nur 4 der 16 Einträge namentlich; welche 12 gefragt wurden und mit welchem der **fünf** kanonischen Ausgänge, ist aus der Akte nicht nachvollziehbar.
- **verifizierbar:** nein — kein Gate prüft eine Audit-Aussage (`planning` prüft Lifecycle-Invarianten, `structure` Überschriften/Abschnitts-Bedingungen; beide liefen hier grün). Belegt durch den Feld-Zensus über alle 16 aktiven `MR-`Dateien und die Index-Tabelle `harness/conventions.md:127-142` gegen `modul-02` §Freshness-Audit.
- **klasse:** `freshness-filter-misst-den-zeiger-statt-den-gegenstand`

### F-2

- **kategorie:** MEDIUM
- **quelle:** Slice-Plan §4 DoD (*„die Reduktion **gemessen** statt angenommen … ob die neuen Zeilen in **seinem** Abschnitt landeten (1 von 4)"*) · `.harness/baseline/v5.12.0/regelwerk/grundlagen-source-precedence.md:142-147`
- **pfad:** `docs/plan/planning/done/slice-149-baseline-v5120-delta-audit.md:143-145`
- **befund:** Schritt 2 wird für die vier Kandidaten in **zwei verschiedenen Granularitäten** ausgeführt: `MR-007` und `MR-013` werden auf **Abschnitts**-Ebene entschieden (§Hard Rule vs. §Fitness Function; §Lifecycle als State Machine vs. §Offene Risiken), `MR-032` dagegen auf **Absatz**-Ebene (*Wann die CR-Pflicht beginnt*), obwohl sein Feld zuerst **§Source Precedence** nennt — und genau dort, am Ende von §Source Precedence unmittelbar vor §Spec-Stratifizierung, landen die sechs neuen Zeilen. Auf der Granularität, mit der die anderen drei entschieden wurden, ergibt der beschriebene Schritt „2 von 4", nicht „1 von 4"; die Zahl ist aus ihrer eigenen Vorschrift nicht reproduzierbar. Hinzu kommt, dass der genannte Absatz *Wann die CR-Pflicht beginnt* gar nicht in §Source Precedence steht (`grundlagen-source-precedence.md:215`, §Spec-Stratifizierung ab `:151`) — die Prüfung des Feldes an der Datei hätte diese Unstimmigkeit gezeigt.
- **verifizierbar:** nein — Prosa-Messung ohne Gate. Belegt durch `grep -n '^###' .harness/baseline/v5.12.0/regelwerk/grundlagen-source-precedence.md` (§Source Precedence `:4`, §Spec-Stratifizierung `:151`) gegen die Einfügeposition `:142-147` und gegen `harness/conventions/MR-032-historie-vor-accepted.md:4-7`.
- **klasse:** `messvorschrift-reproduziert-eigenes-ergebnis-nicht`

### F-3

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §1 · `BEO-012` (Reichweite eines Zitats, Zähler 4)
- **pfad:** `docs/plan/planning/observations.md:11` (`BEO-017`, *„Darum hier keine eigene Regelzeile: sie zu duplizieren verstieße gegen `AGENTS.md` §1"*) · gleichlautend `slice-149:133` und Commit-Botschaft `3bdc94f`
- **befund:** Der Satz, auf den sich die Begründung stützt, gilt ausweislich seines eigenen Wortlauts nur für `AGENTS.md`: *„**Diese Datei** trägt Hard Rules und Pointer auf die kanonischen Quellen und dupliziert deren Inhalt nicht"*. Auf das Beobachtungs-Register angewandt reicht er weiter, als er trägt — und in seiner gedehnten Lesart trifft er den Bestand: `MR-031` und `MR-033` zitieren beide den Kanon-Satz, von dem sie abweichen, als Pflichtfeld `Ersetzt-Baseline-Regel`. Dieselbe §1 sagt zudem, die **verkörperte Form** führe und das Regelwerk sei die nachschlagbare Vertiefung — womit sie für einen Eintrag an der 3×-Schwelle, dessen kanonische Folge die Verkörperung ist, gerade nicht in die gezogene Richtung argumentiert.
- **verifizierbar:** nein — kein Gate liest Begründungen. Belegt durch `sed -n '1,20p' AGENTS.md` gegen `docs/plan/planning/observations.md:11` und gegen die `Ersetzt-Baseline-Regel`-Felder von `MR-031:5-7` und `MR-033:5-11`.
- **klasse:** `geltungsbereich-einer-eigenen-regel-gedehnt`

### F-4

- **kategorie:** MEDIUM
- **quelle:** `BEO-009` Richtung (a) — *„sie behauptet eine Probe oder Änderung, die nicht stattfand"* · `modul-06` §Das Beobachtungs-Register (Spalte `Stand` trägt Zustand und Beleg)
- **pfad:** `docs/plan/planning/observations.md:11` (`BEO-017`, Spalte `Stand`: *„alle drei an einem einzigen Arbeitstag"*) · gleichlautend Commit-Botschaft `3bdc94f` (*„Zaehler 3, alle drei Belege von einem einzigen Arbeitstag"*)
- **befund:** Die drei Belege liegen nicht auf einem Arbeitstag: sämtliche Commits zu `slice-132` und `slice-134` (Anlage, Feat, Review, Move, Closure-Body) tragen das Datum **2026-08-23**, sämtliche Commits zu `slice-138` das Datum **2026-08-25**; auch die `Datum:`-Felder der drei Pläne (`2026-08-23`, `2026-08-25`, `2026-08-23`) ergeben unter keiner Lesart einen gemeinsamen Tag. Die Aussage steht in einer **stehenden** Registerzeile, die jede Wellen-Closure am Lese-Schritt liest, und sie ist dort das Intensitäts-Argument für die erreichte Schwelle. Sie steht im selben Commit, dessen Nachfolger eine andere `BEO-009`(a)-Instanz berichtigt.
- **verifizierbar:** nein — kein Gate liest Datumsangaben in Prosa. Belegt durch `git log --format='%ad %s' --date=short --all --grep='slice-132|slice-134|slice-138'` und `grep '^\*\*Verantwortlich' docs/plan/planning/done/slice-13{2,4,8}-*.md`.
- **klasse:** `registerzeile-behauptet-ungeprueftes-datum`

### F-5

- **kategorie:** MEDIUM
- **quelle:** `modul-02` §Freshness-Audit (fünf benannte Ausgänge; *„Widerspricht sie der neuen Fassung … der Widerspruch gehört benannt"*) · `AGENTS.md` §1 Source Precedence (*„Bei Konflikt gilt die höherrangige Quelle"*) · `BEO-012`
- **pfad:** `docs/plan/planning/done/slice-149-baseline-v5120-delta-audit.md:148-156`
- **befund:** Die Frage nach `MR-033` wird von einem **Kanon**-Satz aufgeworfen (*„Die Erlaubnis ist keine Pflicht: Eine Sicht, die ihre Komponenten nur über Rollen und `ARC-*` führt, ist ebenso konform"*, `modul-03-spec.md:126-129`) und ausschließlich mit **repo-lokalen** Einträgen beantwortet (`MR-032` → `MR-031`); keiner der fünf kanonischen Ausgänge wird benannt, und der neue Satz wird gegen die Beibehaltung nirgends gewogen. Der als Beleg zitierte Satz ist zudem eine Verkürzung, die genau seinen Geltungsbereich tilgt: `MR-031:22-24` lautet *„Wer **ihn** verschärft, weicht von einem **Baseline-Default** ab — auch wenn er nur mehr verlangt"* und ist auf `AGENTS.md` §6 Schritt 3 skopiert; zitiert wird *„wer verschärft, weicht ab, auch wenn er nur mehr verlangt"*. Damit steht in der Akte ein Ergebnis („bleibt"), aber nicht, welcher Ausgang gewählt wurde und woraus — die Angabe, gegen die der nächste Freshness-Audit fragt.
- **verifizierbar:** nein — kein Gate prüft Audit-Begründungen. Belegt durch `sed -n '20,30p' harness/conventions/MR-031-schritt-3-benennen.md`, `sed -n '18,30p' harness/conventions/MR-032-historie-vor-accepted.md` und `sed -n '118,150p' .harness/baseline/v5.12.0/regelwerk/modul-02-harness-bootstrap.md`.
- **klasse:** `kanon-frage-mit-repo-praezedenz-beantwortet`

### F-6

- **kategorie:** MEDIUM
- **quelle:** `BEO-012` (*„vor jedem Zitat einer Regel deren Geltungs-Feld lesen … Wer die Reichweite nicht nachgelesen hat, zitiert nicht, sondern erinnert sich"*) · slice-148 §4/§9
- **pfad:** `docs/plan/planning/open/slice-150-pin-gebundene-zitate.md:85-87`
- **befund:** Das zweite §5-Risiko beziffert die Falsch-Positiv-Last mit *„in slice-148 hat dieselbe Prüfung **sechs Anläufe** gebraucht"*. `slice-148` sagt an beiden einschlägigen Stellen (`:69`, `:131`) **„sechs Paare"** — Dokument/Quell-Paare der vierten Spiegel-Klasse, also das *Ergebnis* der Messung, nicht die Zahl ihrer Versuche; die Wendung „Anläufe" kommt in `slice-148`, im zugehörigen Review-Report und in keiner Commit-Botschaft des Repos vor (`git log --all --format=%B | grep -c 'sechs Anläufe'` ⇒ 0). Die Größenordnung, mit der slice-150 seine Vorsicht gegen einen Zitat-Wächter begründet, hat damit keine Fundstelle.
- **verifizierbar:** nein — kein Gate vergleicht Prosa-Zahlen mit ihrer Quelle. Belegt durch `grep -n 'sechs' docs/plan/planning/done/slice-148-baseline-v5120-vendoring.md` (zwei Treffer, beide „sechs Paare") und die genannte `git log`-Zählung.
- **klasse:** `zahl-aus-der-quelle-umgedeutet`

### F-7

- **kategorie:** MEDIUM
- **quelle:** `.harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md` §Lifecycle als State Machine — *„DoD-Häkchen und Closure-Notiz sind die Bedingung dafür, dass die Datei überhaupt nach `done/` darf"* (State-Machine-Kante `in_progress --> done: DoD + Lerneintrag + Risiko-Ausgänge`) · `BEO-015`
- **pfad:** `docs/plan/planning/done/slice-149-baseline-v5120-delta-audit.md:77-79`
- **befund:** Die Tatsachenberichtigung stellt einen **unangekreuzten** DoD-Punkt her und lässt die Datei zugleich in `done/` liegen; der Kanon macht die DoD-Häkchen zur Bedingung für genau diese Verzeichnis-Position, und die Botschaft von `cdc54c5` begründet die Weiterführung in `done/` nicht (sie begründet nur, dass die Zusage falsch war). Der Zustand ist die Gestalt, die `BEO-015` benennt: `ls open/` und `ls in-progress/` melden nichts Offenes, während das Archiv-Dokument eine ausstehende Pflicht trägt. **Kein Gate fängt es** — das Closure-Profil aus `make fullbuild` (`planning` + `structure`, 467 Dateien) läuft auf genau diesem Stand grün, weil die Regel aus slice-143 den Vorlagen-Platzhalter im §5-Ausgang prüft, nicht den DoD-Haken.
- **verifizierbar:** ja, in der Gegenrichtung — `make fullbuild` Exit 0 auf `HEAD` belegt, dass kein vorhandener Sensor den Zustand meldet; die Aussage des Kanons selbst ist per `grep -n 'Bedingung dafür' .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md` nachlesbar.
- **klasse:** `done-mit-offenem-dod-haken`

### F-8

- **kategorie:** LOW
- **quelle:** `BEO-011` · vorheriges Finding am gleichen Vorgang: slice-148-Review F-5/F-7 (*„Eine Zahl ohne ihre Vorschrift ist keine Messung, sondern eine Behauptung mit Ziffern"*)
- **pfad:** `docs/plan/planning/open/slice-151-urteilsfreie-haelfte-voll.md:49` und `:86`
- **befund:** Die Prüfmenge wird zweimal mit *„139 Slices"* beziffert; `docs/plan/planning/done/` führt am Schnitt-Commit `ac1a23b` **142** `slice-*`-Dateien (dieselbe Zahl an `HEAD`), und „139" ist im Repo an keiner anderen Stelle als gemessene Zahl belegt. Eine Messvorschrift steht nicht dabei. Der Slice leitet aus der Zahl seine Retrofit-Sorge ab (*„könnte breit rot laufen"*), und die nächste Bearbeitung würde sie als Vergleichspunkt übernehmen.
- **verifizierbar:** nein — Prosa-Zahl. Belegt durch `git ls-tree -r --name-only ac1a23b docs/plan/planning/done/ | grep -c '/slice-'` ⇒ 142 und `grep -rn '\b139\b' docs/plan/planning/ | grep -v slice-139` (nur die zwei beanstandeten Stellen).
- **klasse:** `zahl-ohne-messvorschrift`

### F-9

- **kategorie:** LOW
- **quelle:** welle-85 §6 (*„Out-of-Scope für diese Welle … Keine Ausweitung der Closure-Regel auf die volle urteilsfreie Hälfte"*) gegen welle-85 §4 (Slice-Tabelle) · Präzedenz welle-83 §4/§6
- **pfad:** `docs/plan/planning/welle-85-baseline-v5120-migration.md:89-90` gegen `:107-113`
- **befund:** §4 führt `slice-151` jetzt als Slice **dieser** Welle, während §6 unter der Überschrift *Out-of-Scope für diese Welle* dieselbe Arbeit ausschließt. Bei der Wellen-Closure geben die zwei Abschnitte gegensätzliche Antworten auf dieselbe Frage: nach §4 liegen zwei Wellen-Slices offen, nach §6 ist ihr Gegenstand nicht Sache dieser Welle. Die Präzedenz spricht für §4 — in welle-83 wurden die beiden Etappe-C-Slices (`slice-130`, `slice-131`) innerhalb der Welle gebaut und die Welle schloss erst danach.
- **verifizierbar:** nein — kein Gate liest Wellen-Scope-Aussagen (`planning` prüft nur Roadmap-Aktiv-Status gegen `in-progress/`). Belegt durch den Vergleich beider Abschnitte und durch `docs/plan/planning/done/welle-83-baseline-v5110-migration.md` §4/§6.
- **klasse:** `wellen-scope-widerspricht-slice-tabelle`

### F-10

- **kategorie:** LOW
- **quelle:** `BEO-012` (die Reichweite eines Zitats ist genau das, was die Auszeichnung markiert) · `BEO-016` (unbeabsichtigte Spanne, kein Gate sieht sie)
- **pfad:** `docs/plan/planning/done/slice-149-baseline-v5120-delta-audit.md:131`
- **befund:** Die Kanon-Zuschreibung in §9, Zeile 1 der Delta-Tabelle trägt eine unbalancierte Auszeichnung: `**und welcher der drei*` — zwei öffnende, ein schließender Stern. Dieselbe Wendung steht an drei anderen Stellen korrekt als `**und welcher der drei**` (`slice-149:46`, `welle-85:51`, `slice-151:27`). Wo das Zitat aus dem Kanon endet, ist damit im gerenderten Dokument nicht mehr eindeutig ablesbar; `spans` deckt es nicht — das Modul prüft Backtick-Spannen, Fences und verschachtelte Links (`DC-FA-SPAN-001`), keine Emphasis-Parität.
- **verifizierbar:** nein — belegt durch `sed -n '131p' … | cat -A` gegen die drei korrekten Vorkommen und gegen `spec/spezifikation.md:2736-2738` (die drei `spans`-Grund-Codes).
- **klasse:** `unbalancierte-auszeichnung-in-kanon-zitat`

### F-11

- **kategorie:** LOW
- **quelle:** welle-85 §3 Trigger 3 (*„Handlung mit Slice-ID oder belegt folgenlos"*) · welle-85 §4 (*„Zwei Slices, aus den zwei Kurs-Wellen mit Handlungs-Antwort; die übrigen zwei sind belegt folgenlos"*)
- **pfad:** `docs/plan/planning/done/slice-149-baseline-v5120-delta-audit.md:134`
- **befund:** Die Antwort auf Punkt 4 lautet „folgenlos", nennt als zweite Hälfte ihrer Begründung aber einen **Träger**: *„und `slice-147` liegt für die Feedforward-Hälfte bereits in `open/`"*. `slice-147` ist ausweislich seines Kopfes ausdrücklich **wellenlos** und steht folglich nicht in welle-85 §4. Damit verzeichnet die Welle für CR-Punkt 4 „folgenlos", während ein Teil der Antwort an einem Slice hängt, den die Welle nicht führt — fällt oder verschiebt sich `slice-147`, weist kein Wellen-Artefakt darauf hin.
- **verifizierbar:** nein — kein Gate prüft Wellen-Buchführung dieser Art. Belegt durch `sed -n '1,10p' docs/plan/planning/open/slice-147-reviewer-anker-reichweite.md` (Feld `**Welle:** — wellenlos`) gegen welle-85 §4.
- **klasse:** `folgenlos-mit-traeger-ausserhalb-der-welle`

---

## Negativbefunde

1. **Das Delta selbst ist exakt und deckungsgleich mit der Ankündigung — eigene Rekonstruktion beider Bäume.** Fünf Dateien tragen Regel-Inhalt, keine sechste: `modul-05` **+11** (§Offene Risiken, unmittelbar nach *„Ein Slice geht nicht nach `done/`, während ein Risiko ohne Ausgang dasteht"*), `modul-13` **+8** (§Fitness Function, hinter dem Absatz *Die maschinelle Formulierung*), `grundlagen-source-precedence` **+6** (Ende §Source Precedence), `modul-03-spec` **5 für 2** (§Ziel-Form: Architektur-Sicht), `templates/AGENTS.template.md` **3 für 2** (§3.4) — jeweils nach Abzug der Versions-Stempelzeile. Der §5-Ausgang *„geprüft und deckungsgleich … fünf Dateien, keine sechste"* hält wörtlich.
2. **Die vier Kanon-Zuschreibungen in §9, Spalte *Was landete*, sind wörtlich richtig.** *„Die Erlaubnis ist keine Pflicht"* steht so in `modul-03-spec.md:127`; *„das Rot muss von dieser Regel kommen … die Meldung nennt die gebrochene Regel und die Fundstelle, und beides gehört angesehen"* steht (umbrochen über vier Zeilen, Auslassung korrekt markiert) in `modul-13:176-178`; *„Welches Werkzeug die urteilsfreie Hälfte prüft, ist Repo-Entscheidung; dass sie eine hat, ist es nicht"* in `modul-05:152-153`; die Reichweitenfrage in `grundlagen-source-precedence:142-147`. Alle vier gegen den **umbrochenen** Text geprüft, nicht zeilenweise — die Falle, die der Slice selbst beschreibt, ist hier nicht eingetreten.
3. **Der Fehlalarm ist echt und richtig eingeordnet.** `MR-032`s Zitat *„Vor `Accepted` ist das Lastenheft ein Entwurf — frei änderbar, ohne Change Request, ohne Historie-Zeile"* steht in `grundlagen-source-precedence.md:217-218` — über zwei Zeilen umbrochen, ein zeilenweises `grep` verfehlt es. Die Beschreibung in §9 stimmt Wort für Wort mit dem Befund überein.
4. **Schritt 1 der Verengung ist exakt reproduzierbar.** Eigene Auswertung der Spalte `Ersetzt-Baseline-Regel` über alle 16 Index-Zeilen ergibt genau die vier genannten Kandidaten (`MR-007` → `modul-13`, `MR-013` → `modul-05`, `MR-032` → `grundlagen-source-precedence`, `MR-033` → `AGENTS.template.md` **und** `modul-03-spec.md`); die übrigen zwölf zeigen auf `grundlagen-durchsetzungsschicht`, `grundlagen-referenz-richtung`, `grundlagen-harness-dateien`, `modul-10`, `modul-09` oder auf keine Baseline-Datei. Die Beanstandung in F-1 betrifft die **Frage**, nicht diese Zählung. Nebenbei: `MR-025` führt das Feld als `Ergänzt-Baseline-Regel` — der Wert steht in der Index-Tabelle, die Auswertung bleibt richtig.
5. **Schritt 2 hält für `MR-007` und `MR-013`.** `MR-007` zeigt auf `modul-13` §Hard Rule (Doku-Disziplin) (`:40`), die Änderung liegt in §Fitness Function (`:155`); `MR-013` zeigt auf `modul-05` §Lifecycle als State Machine (`:11-78`), die Änderung in `#### Offene Risiken werden bei Closure aufgelöst` (`:126`, unter §Closure- und Lerneintrag-Regeln). Beide Zuordnungen sind an den Überschriften nachgeprüft. Die Beanstandung in F-2 betrifft ausschließlich `MR-032`.
6. **`MR-033` ist tatsächlich der eine betroffene Eintrag, und seine Prämisse ist bestätigt.** Sein Feld `Ersetzt-Baseline-Regel` zitiert zweimal den Wortlaut, den `v5.12.0` ersetzt hat (`AGENTS.template.md` §3.4 und `modul-03-spec.md` §Ziel-Form), beide Pfade sind gehoben und lösen auf, beide Zitate existieren am neuen Ziel nicht mehr. Der Kanon liest die Stelle jetzt ausdrücklich als **Erlaubnis** für **Code**-Modul-Pfade — die Lesart, auf der `MR-033` steht. Dass der Eintrag stehen bleibt, ist inhaltlich tragfähig: `AGENTS.md` §3.4 formuliert ein **Verbot** und weicht damit vom Vorlagen-Text ab, unabhängig davon, dass eine rollen-geführte Sicht „ebenso konform" ist. Die Beanstandung in F-5 betrifft die **Begründung**, nicht das Ergebnis.
7. **Kein toter Anker im Konventionsspeicher — die Formcheck-Voraussetzung des Kanons ist am Bestand erfüllt.** Eigene Probe über alle 38 Dateien unter `harness/conventions/`: jeder Verweis in den vendorten Baum löst auf, jeder Anker findet seine Überschrift oder ihr `<a id>` (inkl. der Doppel-Bindestrich-Slugs `#template-schichtung--was-der-rumpf-trägt-und-was-der-kommentar`). Die einzigen nicht auflösenden Ziele liegen in `conventions/done/` und zeigen auf entfernte Alt-Bäume (`v5.6.0`, `v5.7.0`) — genau die eingefrorenen Fälle, die `ignore-refs` deckt. `make gates` ist grün und die Aussage damit nicht still.
8. **`BEO-017` trägt seine drei Gestalten, und es sind drei verschiedene.** (1) *falscher Linter, vom Nachbarn rot* — `slice-132:98-101` und `:167`; (2) *Direktiv-Form gilt als „wohlgeformt", war nie gefahren* — `slice-134:145-146`; (3) *Probe ans Dateiende angehängt, landet in `## Geschichte`, dem einzigen von `matrix` ausgenommenen Abschnitt* — `slice-138:178-181`. Drei Ereignisse, drei Slices, keine Doppelzählung; dass `slice-134` die erste Gestalt rückblickend miterzählt, ändert ihre Zurechnung zu `slice-132` nicht. Auch die Gegenrichtung im Text (*der unveränderte Bestand, auf dem der Sensor schweigen muss*) hat ihre Kanon-Deckung: `modul-13:179-181` verweist dafür auf `modul-11` §Fitness Function ohne Standard-Tool.
9. **Die Beleg-Form von `BEO-017` erfüllt `modul-06` §Das Beobachtungs-Register.** **Form:** drei `slice-<NNN>`-Kennungen, kein Freitext · **Anzahl:** drei Belege bei Zähler 3 · **Lage:** alle drei Ziele liegen in `docs/plan/planning/done/` und die Prüfung erfolgt nach dem `git mv`. Die Spalte `Sub-Area` trägt `*`, das in der Modus-Deklaration von `harness/conventions.md` geführte Repo-Default. Die Registerzeile steht in der bestehenden Zeilen-Ordnung oben; die Beanstandung in F-4 betrifft eine Aussage **in** der Zelle, nicht die Form.
10. **Die Prämisse von `slice-150` ist kanon-gedeckt, und seine Vorfrage ist im Kanon nicht bereits beantwortet.** *„Einträge werden nie überschrieben"* steht wörtlich in `grundlagen-harness-dateien.md` §Konventionsspeicher (Pflichtgliederung, Zeile Adaptions-Block); `grundlagen-source-precedence.md:134-140` ergänzt für den Baseline-Update-Fall die Nachfolger-Disziplin. Beide sprechen über **Gegenstandslosigkeit und Widerspruch**, keiner über die Pflege eines veralteten **Zitats** in einem geltenden Eintrag — und `MR-021` bindet ausdrücklich **Links**, nicht Zitate. Die Frage ist damit offen und als Regelfrage richtig geschnitten. Gezielt gesucht wurde nach „Zitat"/„zitiert" im gesamten vendorten Regelwerk.
11. **Beide neuen Slices halten `modul-05` §Ziel-Form: Slice.** `slice-150`: drei Liefer-Punkte (Regelfrage · Eintrag/Bestandsfall · gemessene mechanische Form), Sub-Areas Konventionsspeicher + Produkt-Module, einzeln lieferbar, in einer Sitzung prüfbar. `slice-151`: drei Liefer-Punkte (Ausdrückbarkeits-Messung · Bestandsmessung · Umsetzung oder begründeter Verzicht), Module `structure` + `planning`, das Produkt-Delta ausdrücklich als eigene Anforderung ausgelagert. Beide tragen `Verantwortlich:`/`Autor:`, ein `Berührte Spec-Stellen`-Feld (bei `slice-151` mit der Kennung `DC-FA-STRUCT-001`, nicht nur mit `§N`), §7-Sichtung, §8-Modus-Begründung und §5-Ausgänge in der Repo-üblichen Platzhalter-Form. Ohne Befund.
12. **Der Etappe-C-Schnitt selbst ist vollständig gebucht.** Beide Slices stehen in welle-85 §4 mit Rolle, im Drift-Log der Roadmap mit Datum und Begründung, und der Roadmap-Ruhe-Marker (*„Nichts in Arbeit."*) ist zurückgesetzt; `in-progress/` enthält nur noch `roadmap.md`. `make gates` (mit `planning-check`) bestätigt die Roadmap-↔-`in-progress`-Invariante. Der Move nach `MR-013` ist gebündelt und die vier gekoppelten Verweise (`slice-148:30`, `:93`, welle-85 §4, Drift-Log) sind mitgezogen.
13. **§3 ist eingehalten: das Audit hat geschnitten, nicht gebaut.** Der Diff `297a58a..HEAD` fasst weder den vendorten Baum noch `harness/conventions/` noch `spec/` an; das veraltete Zitat in `MR-033` steht unverändert, kein eingefrorenes Dokument wurde editiert, und keine Regel des neuen Stands wurde angewandt. Die Selbst-Auskunft *„Was das Audit NICHT getan hat"* ist am Diff überprüft und stimmt.
14. **Die Tatsachenberichtigung ist in ihrem Gegenstand vollständig und ehrlich.** Der falsche Haken ist nicht stillschweigend entfernt, sondern als verfrüht ausgewiesen; die Gate-Zusage ist von der Review-Zusage getrennt; die zusätzlich aufgenommene Angabe `make fullbuild` Exit 0 ist durch die Botschaft des Closure-Commits `3bdc94f` gedeckt und von mir eigenständig reproduziert (Exit 0, identischer `image-hash`). Weitere DoD-Haken oder §5-Ausgänge, die eine nicht stattgefundene Handlung behaupten, habe ich nicht gefunden — insbesondere sind die zwei §5-Ausgänge Aussagen über eigene Messungen, nicht über fremde Läufe. Die Beanstandung in F-7 betrifft die **Verzeichnis-Position**, nicht die Berichtigung.
15. **Gate-Belege der Botschaften exakt reproduziert.** *„`make gates` Exit 0 (zehn Glieder)"* ✓ — die Abschlusszeile nennt genau zehn; *„`make fullbuild` Exit 0"* ✓ mit `image-hash sha256:f41f99ad…`, Closure-Profil 467 Dateien / 0 Befunde, 48 Anforderungen / 0 Waisen. Der in `3bdc94f` berichtete Gate-Fang (*„sechs nackte Kennungen im Closure-Body"*, Linkpflicht) ist am Ergebnis sichtbar: alle `MR-`/`BEO-`/`slice-`-Nennungen in §9 tragen Links, und die `MR-`Verweise adressieren die Index-Zeile (`harness/conventions.md#mr-<NNN>`) statt der Eintrags-Datei — die vom Kanon verlangte bruchsichere Form.
16. **Referenz-Richtung und Marker-Ehrlichkeit ohne Befund.** Die neuen Dokumente enthalten keinen Provenance-Marker; die Abwärts-Verweise in `slice-150`/`slice-151` (auf ADRs, `done/`-Slices, Register) laufen von Planungs- in Planungs-/Entscheidungs-Artefakte und verletzen die SDP-Matrix nicht. `make gates` mit aktivem `matrix` ist grün.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 7 | F-1, F-2, F-3, F-4, F-5, F-6, F-7 |
| LOW | 4 | F-8, F-9, F-10, F-11 |
| INFO | 0 | — |

## Verdikt

**Blockierend.** Kein HIGH; sieben MEDIUM.

Die **Sacharbeit am Delta** ist ohne Befund: Ich habe beide gepinnten Bäume rekonstruiert, die fünf Inhalts-Dateien Zeile für Zeile gelesen und jede Kanon-Zuschreibung in §9 gegen den **umbrochenen** Text gehalten — die vier Antworten treffen den vendorten Wortlaut, der Umfang stimmt bis auf die Datei genau, und der beschriebene Fehlalarm ist echt. Auch der Schnitt selbst trägt: zwei Slices, beide baseline-konform dimensioniert, beide vollständig gebucht, und das veraltete Zitat in `MR-033` steht bewusst unrepariert. Der Slice hat getan, was er zu tun ankündigte.

Blockierend ist, **womit** die zwei Kernaussagen belegt sind. Der Freshness-Audit stellt nicht die Frage des Kanons: er misst, wohin der `Ersetzt-Baseline-Regel`-Zeiger eines Eintrags zeigt, während der Kanon nach dem **Gegenstand** fragt und ausdrücklich verlangt, *„durch die Adaptions-Liste, nicht nur durch den Diff"* zu gehen (F-1) — für vier der 16 Einträge kann der Filter „betroffen" gar nicht liefern, und die Akte nennt nur vier Einträge namentlich. Der zweite Schritt reproduziert dazu sein eigenes Ergebnis nicht: auf der Granularität, mit der `MR-007` und `MR-013` entschieden wurden, wäre `MR-032` der zweite Kandidat (F-2) — und ausgerechnet der neue Absatz dort spricht über die Reichweite einer `MR-<NNN>`. Die `MR-033`-Entscheidung ist im Ergebnis richtig, aber mit repo-lokaler Präzedenz statt mit dem höherrangigen Kanon begründet, und der als Beleg zitierte Satz verliert in der Verkürzung genau seinen Geltungsbereich (F-5). Dieselbe Bewegung trägt F-3: `AGENTS.md` §1 sagt „diese Datei", nicht „dieses Repo".

Zwei Befunde sind reine Tatsachen. `BEO-017` behauptet in einer **stehenden** Registerzeile — der Zeile, die jede Wellen-Closure am Lese-Schritt liest — „alle drei an einem einzigen Arbeitstag"; die drei Belege liegen auf zwei Tagen (F-4). Und `slice-150` begründet seine Vorsicht mit „sechs Anläufe" aus `slice-148`, wo „sechs Paare" steht (F-6). Beide sind `BEO-012`/`BEO-009`-Klasse in einem Commit-Strang, der genau diese Klassen registriert.

F-7 ist die unbequemste Stelle: Die Berichtigung ist in ihrem Gegenstand vorbildlich — sie schreibt nichts still um —, lässt aber einen Zustand stehen, den der Kanon in einem Satz ausschließt (*„DoD-Häkchen … sind die Bedingung dafür, dass die Datei überhaupt nach `done/` darf"*). Dass `make fullbuild` auf genau diesem Stand grün läuft, ist der Beleg, nicht die Entlastung: es ist die Gestalt aus `BEO-015`, diesmal mit einer offenen Zusage statt eines erfundenen vierten Ausgangs.

Die vier LOW reisen mit: eine Zahl, die die nächste Bearbeitung als Prüfmenge übernehmen würde (F-8), zwei Abschnitte desselben Wellendokuments mit gegensätzlicher Antwort auf die Closure-Frage (F-9), ein unbalanciertes Zitat-Ende in einer Kanon-Zuschreibung (F-10) und ein „folgenlos", dessen Träger außerhalb der Welle liegt (F-11).
