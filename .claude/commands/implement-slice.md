# Slice implementieren (Harness)

Argument: $ARGUMENTS

Dieser Command führt die **Implementer**-Rolle (Modul 9) für *einen* Slice — innerhalb der
Rollen-Sequenz Planner → Architect → Implementer → Reviewer → Verifier → Validator →
Planner-Closure (Modul 8). **Rollen-Trennung ist Kontext-Trennung:** die nachgelagerten Rollen
(Review, Verifikation, Validation, Closure) laufen in **frischem Kontext** (Subagent / geleerter
Kontext), nie im Kontext, der den Code schrieb — sonst wiederholt sich derselbe blinde Fleck.
Keine Rolle springt rückwärts ohne Übergabe-Artefakt (Findings · Folge-ADR · Carveout, Modul 8).

Kanonische Quellen (vendored Regelwerk, `.harness/baseline/<tag>/regelwerk/`): Modul 9
(Implementierung), Modul 5 (Lifecycle), Modul 8 (Rollen), Modul 10 (Review), Modul 11
(Verifikation).

## Repo-lokale Adaptionen, die du beachten MUSST

Über das Regelwerk hinaus trägt dieses Repo lokale Adaptionen gegenüber der Baseline. Lies den
Adaptions-Block („MR-Block") in `harness/conventions.md`; die workflow-relevanten:

- **Docker-only, kein Host-Toolchain.** Jeder Gate und jedes Tool läuft in einem gepinnten
  Docker-Image; der PreToolUse-Guard (`.claude/hooks/pretooluse-command-guard.sh`) blockt
  Host-Go, Host-Paketmanager und Host-Skript-Interpreter (und prüft Sub-Shell-Strings,
  `AGENTS.md` §3.1). Rufe nie einen Host-Toolchain auf — nur die `make`-Targets.
- **Gate-Nachweis + Stop-Hook.** `make gates` endet mit `record-gates`, das einen Content-Hash des
  Working Tree und die Commit-SHA von HEAD stempelt. Der Stop-Hook ist **streng, ohne Schalter**:
  jedes Turn-Ende mit Inhalt, den kein grüner Lauf deckt, blockiert — auch ohne neuen Commit, und
  **ein Commit oder `git mv`-Commit ohne frischen `make gates`-Lauf lässt ihn rot**. Locker wäre er
  nur mit der Datei `.harness/stop-gate-locker`; sie fehlt, und sie einzuführen ist eine Lockerung
  des Handoff-Gates, die eine ADR braucht ([`MR-079`](../../harness/conventions.md#mr-079)).
- **Strenges Doc-Gate (d-check).** Jede `DC-`/`ADR-`/`MR-`-Kennung in einer gescannten `.md` muss
  ein klickbarer Anker-Link sein (link-policy: always) — ein bares Kennungs-Token bricht
  `make doc-check` (`id-unlinked`). `codepaths` verlangt, dass Pfade in Inline-Code existieren: eine
  *geplante* Datei braucht einen Inline-`d-check:ignore`-Marker, eine *bewusst entfernte* gehört in
  `ignore-refs`. Spec verweist nie abwärts auf ADR/Slice; ein Verweis auf eine superseded ADR nur
  via Inline-Code + `d-check:ignore`. `docs/reviews/**` (Zeitdokumente) ist von der
  Referenz-Richtung (`matrix`) und von der Pfad-Existenz (`codepaths`) ausgenommen; Kennungen
  prüft das Gate dort wie überall.
- **Neue Artefakte per `cp` aus den vendored Templates** (`.harness/baseline/<tag>/templates/…`),
  dann ausfüllen — keine handgeschriebenen oder repo-gepflegten Template-Kopien.
- **Commit via Message-Datei** (`git commit -F <datei>`): der Guard scannt den Command-String,
  also nie eine Commit-Message inline, die ein geblocktes Tool-Token enthält.
- **Commit-Kennung.** Eine Commit-Message ohne Kennung (`DC-*`, `ADR-NNNN`, `MR-NNN`,
  `slice-<Kennung>`; `AGENTS.md` §5 Zeile 1) weist der versionierte Hook `.githooks/commit-msg` ab
  (`make trace-check`, Modul `commits`), sobald er aktiviert ist. Er **reist mit dem Klon, seine
  Aktivierung nicht**: `make hooks` setzt `core.hooksPath`; `git commit --no-verify` umgeht ihn, und
  das Netz dafür ist die PR-/Push-CI. Geprüft wird die **Anwesenheit** einer Kennung, nicht ihre
  Wahrheit. Was er **nicht** erreicht: die zweite Hälfte der Traceability-Zusage — ein
  Doku-Update bei berührtem öffentlichem Vertrag — ist von einem Commit-Wächter nicht mechanisch
  prüfbar und bleibt deine Arbeit. Er **erreicht** dagegen jede Commit-Klasse, die `git` erzeugt —
  auch die der Repo-Werkzeuge (`make slice-mv`, `make archive-wave`).

## Kontext lesen (Modul 9, Schritte 1–3)

1. `CLAUDE.md` lesen (der Werkzeug-Einstieg; er importiert `AGENTS.md`).
2. `harness/README.md` lesen.
3. `AGENTS.md` lesen.
4. `harness/conventions.md` lesen.
5. Den Regelwerk-Index (`.harness/baseline/<tag>/regelwerk/README.md`) und das aufgabenrelevante
   Modul **on-demand** lesen (Source Precedence, committet vendored Baseline). Nicht den ganzen
   Baum laden.
6. Die als Argument übergebene Slice-Datei lesen.
7. Alle referenzierten ADRs und Anforderungen lesen.
8. Berichten: Slice-ID · DC-IDs · ADR-IDs · betroffene Komponenten · zu laufende Gates.

## Nach in-progress eintreten (Modul 5 Lifecycle + Modul 8 Übergabe)

9. Der Implementer erhält den Slice **in `in-progress/`** (Planner→Implementer-Übergabe,
   Modul 8; `next → in-progress` = „Implementer beginnt", Modul 5). Liegt er noch in `open/`,
   zuerst dorthin verschieben (`open → next → in-progress`); `open → next` setzt dabei das
   Kopf-Feld `Verantwortlich:`. **Jeder dieser Übergänge läuft über
   `make slice-mv SLICE=slice-<Kennung> TO=<next|in-progress>`** — das Werkzeug bewegt die Datei
   per `git mv` und committet den **reinen Move** sofort als eigenen Commit (Hard Rule 3.3), und
   es zieht danach die **Verweise** nach: **eingehend** jede Präfix-Form auf die bewegte Datei,
   repo-weit, dazu präfixlose Links aus den Geschwistern im Ausgangsverzeichnis, und **ausgehend**
   die präfixlosen Ziele innerhalb der bewegten Datei selbst. Fielen
   Verweise an, committet es sie als **zweiten**, vom Move getrennten Commit; sonst bleibt es beim
   einen Move-Commit. Der Move landet damit **auf dem Hauptzweig, vor der Arbeit** — der Branch
   entsteht danach. Reist er erst im PR mit, ist der Zustand zweigelokal, und `in-progress/`
   bleibt für alle anderen leer, bis die Arbeit fertig ist.
   **Was das Werkzeug nicht kann, steht in seinem Kopf (`tools/harness/slice-mv.sh`), und der Rest
   ist deine Handarbeit:** es verlangt einen sauberen Arbeitsbaum, zieht Pfade nach, keine
   Zustandsätze (den Ruhe-Marker der Roadmap hält `make planning-check`), und einen präfixlosen
   **eingehenden** Verweis erkennt es nur als Markdown-Link in den Geschwistern, die flach im
   Ausgangsverzeichnis liegen. `make doc-check` nach dem Move zeigt, was stehen blieb — zieh es
   nach, bevor der nächste Schritt startet.
   **Vor der Übergabe an den Implementer liegt das Architect-Verdikt vor** (Modul 8: Planner →
   Architect → Planner, dann Planner → Implementer): bestätigte Entscheidungs-Bezüge oder ein
   Folge-Entscheidungs-Vorschlag, als Report unter `docs/reviews/`. **Eine ADR schreibt der
   Architect, nicht der Implementer** (Modul 8 §Rollen-Regeln: „Architect schreibt; Reviewer
   prüft auf Konsistenz; Implementer liest als Constraint"). Liefert der Slice eine ADR, liest der
   Implementer sie als Constraint; eine nötige Änderung meldet er als Befund an den Architect,
   statt sie selbst zu schreiben.
10. WIP-Limit = 1 pro Implementer (Modul 5): kein paralleles `in-progress/`.
11. Lifecycle-Rücksprungkanten (Modul 5), falls sich der Slice als falsch erweist: zu groß →
    `in-progress → next` (zurück zur Zerlegung); blockiert → `in-progress → open` (Carveout,
    Modul 7). Zurückführen ist Disziplin, kein Scheitern.

## Plan vor Code (Modul 9, Schritt 4 — nicht optional)

12. **Den Ist-Zustand gegen den Slice-Plan messen, bevor du editierst** (`grep`/`diff`, nicht
    `edit`) — Geschwister-Slices lassen Pläne altern (gelöschte Pfade, verschobene
    Lifecycle-Dateien). Drift zuerst abgleichen; keinen veralteten Plan blind abarbeiten.
13. Die kleinste sinnvolle Änderung gegen die DoD planen. Erst planen, dann coden.

## Implementieren und gaten (Modul 9, Schritte 5–6)

14. Die kleinste sinnvolle Änderung implementieren.
15. Zuerst den engsten nützlichen Gate laufen lassen (z. B. eine Testdatei / ein Gate).
16. `make gates` laufen lassen.

**Vor jedem Commit — repo-eigene Regeln aus dem Steering Loop** (je mit Herkunfts-Anker; alle
*permanent*, denn jede ist ein Urteil, das kein Gate sieht):

- Vor jedem Code-Commit `git diff --cached --name-only` gegen die Plan-Tabelle in §3 und ihre
  Plan-Änderungs-Notizen halten. Jede dort nicht genannte Datei bekommt zuerst eine
  Plan-Änderungs-Notiz, als eigener Commit vor dem Code (`AGENTS.md` §6 Schritt 4).
  (seit slice-259; permanent — ob eine Notiz eine Datei deckt, ist ein Urteil)
- Vor jedem Code-Commit jede Kommentarzeile, die der Diff berührt — auch eine nur umbrochene oder
  verschobene —, gegen die fünf Klassen aus `AGENTS.md` §3.7 lesen: keine Geschichte früheren
  Verhaltens, keine Herkunfts-Prosa, kein externer Verweis. (seit slice-259; permanent — Chronik
  ist ein Urteil, kein Gate sieht sie)
- Ersetzt eine Änderung einen Mechanismus (einen Pfad, ein Muster, einen Ort, an dem etwas läuft),
  vor dem Code-Commit den alten Begriff oder das alte Muster über `harness/`, `docs/user/`,
  `.claude/` und die Kommentare der berührten Skripte greppen und jeden Treffer im Plan aufführen
  — die Führungs-Dokumente außerhalb der Spec sind ebenfalls Spiegel, wie in
  [MR-025](../../harness/conventions.md#mr-025). (seit slice-260; permanent — ob ein Treffer noch
  den alten Mechanismus beschreibt, ist ein Urteil)
- Jede Liste, die in einen Spec-Eintrag, die Grenzen einer Sensor-Datei oder einen Kommentar
  geschrieben wird — was blockiert, was fällt, was ausgenommen ist —, wird am Code gezählt, und
  das Zähl-Kommando steht im Plan neben der Zahl. Eine aus einer anderen Beschreibung gelesene
  Liste ist nicht gezählt. (seit slice-260; permanent — ob eine Liste vollständig ist, ist ein
  Urteil über den Gegenstand, kein Sensor hält es)
- Wer ein Erkennungs-Muster schreibt oder ändert (einen regulären Ausdruck, der entscheidet, was als
  Befund, Zusage oder Kandidat zählt), listet und testet vor dem Code-Commit Negativfälle —
  Verneinung, Kompositum, Satzmitte, Zeilenumbruch, Code-Block, Zitat — und fährt das Muster über
  den eigenen Korpus des Repos, das Zähl-Kommando im Plan. Beantwortet das Produkt die Frage schon
  (was ist ein Link, eine Überschrift, ein Code-Span), dient dessen Erkennung statt eines neuen
  Musters. (seit slice-265, slice-267; permanent — welcher Freitext einer Form ähnelt, ist ein
  Urteil)
- Vor jedem Commit `git diff --cached --stat` lesen und gegen die Botschaft halten: eine früher
  gestagte Datei (für einen Gate-Lauf, einen Move) reist in den nächsten Commit mit. (seit
  slice-267; permanent — was in einen Commit gehört, ist Absicht, kein Hook sieht sie)

**Plan-Defekt-Rücksprungkanten (Modul 9):** ein roter Sensor (15) oder rotes Gate (16) führt
zurück zum **Plan** (13) — den Plan verfeinern, nicht den Kontext neu lesen. Ein Rücksprung zu
Schritt 1 signalisiert einen Kontext-Defekt. Ein struktureller Fehlschnitt (zu groß / blockiert)
ist eine Lifecycle-Rücksprungkante (11).

## Pre-completion-Checkliste (Modul 9, Schritt 8 — letzte Handlung der Implementer-Rolle)

17. Doku aktualisieren, falls ein öffentlicher Vertrag berührt ist; ADR und ADR-Index trägt der
    Architect (Schritt 9). CHANGELOG, READMEs und Handbuch-Kopf sind Release-Prep, nicht der
    Feature-Commit (`AGENTS.md` §5 Zeile 17).
18. Die Pre-completion-Checkliste laufen: die DoD Punkt für Punkt **behaupten** und die
    **Sensor-Belege** anhängen — `make gates` **und die Nicht-Gate-Sensoren, die den Slice
    betreffen** (die dieses Repo führt — z. B. ein Mutations-Sensor, wenn Wächter neu/geändert sind;
    ein Emit-/Integrations-Smoke, wenn der betroffene Pfad berührt ist). Modul 11 verlangt genau
    hier den Lauf: *„der Implementer-Agent läuft `make verify-*` **selbst** vor der
    ‚fertig'-Meldung"* — ein Sensor, der erst zur Wellen-Closure feuert, ist pro Slice keiner.
    **Ein nicht gelaufener Sensor ist ein Befund, kein Formfehler:** ihn wegzulassen ist eine
    Aussage („betrifft diesen Slice nicht"), die begründet werden muss. Kein Gate erzwingt das —
    der Stop-Hook deckt nur `make gates`. Das ist die *Behauptung* der Implementer-Rolle und die
    *Eingabe* des Verifiers — **nicht** das finale DoD-Urteil (Modul 11: „Behauptung ohne
    Bestätigung ist die häufigste Verifier-Lücke"; eine DoD-Verletzung ist eine Verifier-only-Klasse,
    unsichtbar für Review und Tests). Ausgeführte Sensors + Restrisiken berichten.
19. **Zu jedem neuen oder geänderten Wächter die rot färbende Mutation benennen**
    (`AGENTS.md` §3.6). Ein grüner Gate-Lauf belegt nur, dass nichts *bricht* — nicht, dass
    der Wächter greift. Pro Zusage also: *welche Änderung am geprüften Code müsste diesen
    Test rot machen, und wurde sie einmal gesehen?* Wo die Antwort dauerhaft interessant
    ist, gehört sie in den Mutations-Sensor (falls vorhanden); wo sie einmalig ist, in
    den Bericht. **Keine Antwort ist ein Befund**, kein Formfehler — die Klasse „Zusage greift
    weiter als Abdeckung" ist in der Praxis teuer erkauft.
20. **Jeden in diesem Lauf neu geschriebenen oder geänderten Kommentar gegen `AGENTS.md` §3.7
    prüfen** (Code, Konfiguration, Skripte). Die Probe: beschreibt der Satz den **Ist-Zustand**
    (indikativ, auflösbar), oder trägt er eine Slice-Nummer als Begründung, ein „(… , entschieden)"
    ohne Anker-Form, oder einen Konjunktiv über eine verworfene Alternative bzw. eine noch nicht
    existierende künftige Änderung (**„sobald Slice X das tut …"**)? Herkunft steht nur als **ein**
    auflösbares Feld in den dort genannten Formen (`DC-*`, `ADR-*`, `MR-*`, `· seit welle-<Kennung>`, wellenlos
    `· seit slice-<Kennung>`) — alles andere ist Zustand, keine Chronik, und wird vor der Übergabe
    umformuliert statt mitgeschleift.

Hier endet die Implementation. Die übrigen Rollen laufen in **getrennten Kontexten** (Modul 8).

## Übergaben an nachgelagerte Rollen (Modul 8 → 10 → 11)

21. **→ Reviewer (Code-Review, Modul 10):** den Diff + Plan-Verweis an einen **unabhängigen**
    Reviewer übergeben (`.harness/skills/reviewer.md`, frischer Kontext — kein Selbst-Review). Er
    kategorisiert Findings (HIGH/MEDIUM/LOW/INFO) in einen Report unter `docs/reviews/` und prüft
    den Diff gegen **Plan + ADR + Hard Rules** (nicht die DoD). HIGH/MEDIUM auflösen; ein HIGH mit
    Rollen-Konflikt folgt Modul 8 §Konflikt-Pfad (Sequenz mit Übergabe-Artefakten, nie
    „herabstufen, weil der Implementer widerspricht"). Befunde an einer ADR gehen an den
    Architect, nicht an den Implementer (Schritt 9).
    - Vor einem Fix-Commit nach dem Review jeden Befund über seine **Aussage** suchen, nicht über
      die zitierte Fundstelle: den alten Begriff oder die alte Behauptung über Code, Doku und Plan
      greppen und die Treffer in der Plan-Änderungs-Notiz aufführen. (seit slice-265; permanent —
      ob ein Treffer dieselbe Aussage trägt, ist ein Urteil)
    - Schreibt ein Fix für Review- oder Verifikations-Befunde mehr, als die Befunde verlangen —
      einen neuen Mechanismus, eine Überarbeitung, einen neuen Lesepfad —, läuft eine weitere
      Reviewer-Runde über genau die Fix-Commits, bevor der Verifier sie sieht; sonst liest der
      größte Code-Anteil des Slice kein Review. (seit slice-258; permanent — die Größe eines Fixes
      ist ein Urteil, kein Gate sieht sie)
22. **→ Verifier (Modul 11):** in getrenntem Kontext die DoD-/Spec-Behauptung und den
    Plan-vs-Code-Diff **bestätigen**, dazu ADR-Konformität. Das fängt, was Tests übersehen und der
    Reviewer nicht sieht (DoD-Verletzung).
23. **→ Validator (Modul 8):** falls der Slice End-Nutzer-Wert liefert, gegen den realen Bedarf
    validieren („das Richtige bauen"). Meist n/a bei interner Wartung — dann explizit sagen statt
    still überspringen.

## Closure — Planner-Rolle (Modul 8 + Modul 5)

24. Erst wenn der Review konform **und** die Verifikation die DoD bestätigt hat, schließt der
    **Planner**: die Closure-Notiz mit einem **Steering-Loop-Eintrag** schreiben (geschärfte Regel ·
    neuer Sensor · benannte Spec-Lücke — Modul 5: der `→ done`-Übergang verlangt einen Lerneintrag,
    nicht nur grüne Gates), dann den Slice `in-progress → done` verschieben —
    **`make slice-mv SLICE=slice-<Kennung> TO=done`**, derselbe Aufruf wie in Schritt 9: reiner
    Move als eigener Commit, Verweis-Nachzug getrennt davon (Hard Rule 3.3). Ein rotes Gate erreicht `done/` **nur** mit
    dokumentiertem Carveout (Modul 7), nie als stilles Rot. **Jedes offene Risiko aus dem Slice-Plan
    bekommt dabei genau einen von drei Ausgängen** (Modul 5): *eingetreten* → Carveout oder
    Folge-Slice mit ID · *entfallen* → gestrichen **mit Begründung** · *weiter offen* → wandert ins
    Beobachtungs-Register (Schritt 25). Ein Slice geht nicht nach `done/`, während ein Risiko ohne
    Ausgang dasteht.
25. **Das Beobachtungs-Register fortschreiben** (`docs/plan/planning/observations/`, Modul 6) —
    der **Schreib**-Schritt, und er hängt an der Closure, nicht an der Implementation. Für jede
    Beobachtung aus der Closure-Notiz: führt das Register die Klasse schon, dann die vorhandene
    Kennung `BEO-<KUERZEL>/<slug>` **zitieren** und eine weitere Datei in ihrem `evidence/` anlegen
    — wer neu formuliert, spaltet eine Klasse in zwei Pfade, und keiner der beiden erreicht je 3×.
    Sonst ein neues Verzeichnis `BEO-<KUERZEL>/<slug>/` mit `observation.md` und `state.md` anlegen
    — Kürzel aus der Modus-Deklaration nachschlagen, nicht erfinden; das Register ist zugleich die
    Vergabestelle für den `<slug>`-Teil. Der Beleg ist **formgebunden**: `evidence/slice-<Kennung>.md`,
    kein Freitext, eine Datei je Auftreten. Geschrieben wird er **vor** dem Move aus Schritt 24 — die
    Slice-Datei liegt dann noch nicht in `done/`, und das ist richtig so, weil Move und Inhalt
    getrennt committen (Hard Rule 3.3). Der Zähler wird **nicht gesetzt**, er ist die Zahl der
    Evidence-Dateien und **folgt** aus ihnen. **Bei null Beobachtungen** bleibt die Ablage
    unverändert und die Closure-Notiz trägt den Satz *keine Beobachtung angefallen*: das Auslassen
    ist keine Antwort. Erreicht ein Eintrag **mit diesem Slice** 3× (die Zahl seiner
    Evidence-Dateien), wandert er in die Steering-Loop-Einträge der laufenden Welle-Closure
    (`/close-welle`); läuft keine Welle, löst die Slice-Closure den Lese-Schritt selbst aus, und der
    Herkunfts-Anker lautet dann `seit slice-<Kennung>` statt `seit welle-<Kennung>`.

Gates nicht überspringen. Keine Erfolgsmeldung ohne Command-Ausgabe.
