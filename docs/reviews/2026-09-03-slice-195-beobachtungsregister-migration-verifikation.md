# Verifikations-Report: slice-195 — 2026-09-03

**Verifikations-Art:** DoD-/Spec-Konformität gegen den Slice-Plan (Modul 8
§Verification vs. Validation), unabhängig von Implementer und Reviewer.

**Gegenstand:** Commit `4d3a386` (feat: Beobachtungs-Register-Datenmigration
Tabelle → Verzeichnis).

**Nachtrag:** `4d3a386` ist nach diesem Report in zwei Commits zerlegt
worden (`94b19bd` Beanspruchung, `b1b960b` Migrationsinhalt) und trägt
die Behebung der beiden HIGH-Befunde — siehe slice-195 §9.

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-03

**Eingangs-Kontext:**
[slice-195](../plan/planning/in-progress/slice-195-beobachtungsregister-migration.md)
§4 Definition of Done, §5 Risiken, §9 Closure-Notiz; `AGENTS.md` (Hard
Rules); [`harness/conventions.md`](../../harness/conventions.md) MR-013,
MR-049, MR-059.

---

## DoD-Punkte

- [x] **Alle Bestandseinträge migriert — 28.** Nachgezählt: das Verzeichnis
  führt **29** `BEO-<KUERZEL>/<slug>/`-Ordner (27× `ALL`, 1× `HARN`, plus
  `BEO-HARN/check-latest-blind-before-pin`), nicht 28. Die Differenz ist
  aufgeklärt und kein Fehlbestand: der 29. Ordner
  (`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`) ist die **eigene**
  Beobachtung dieses Slice aus §9 „Lerneintrag" — eine legitime dritte Quelle
  des Closure-Eintrags nach Modul 5, nicht Teil der 28 migrierten
  Alt-Zeilen. `evidence/`-Dateien: **89** gesamt (88 migriert + 1 neu für
  dieselbe Beobachtung) — deckt sich exakt. Fünf Stichproben (`BEO-023/012/013/027/028`
  aus der alten Tabelle) gegen die neuen `evidence/`-Verzeichnisse geprüft:
  Zähler 6/12/1/1/1 stimmen in allen fünf Fällen exakt, ebenso die
  `Stand`-Texte (wortgleich übernommen, z. B. BEO-012 → `citation-stretched-beyond-scope`).
  **Bestätigt.**
- [x] `docs/plan/planning/observations.md` entfernt (Datei existiert nicht
  mehr), `docs/plan/planning/observations/README.md` neu und **substantiell**
  (30 Zeilen: Verzeichnisform, drei Dateien/drei Lebensdauern, Schreiber/Leser,
  „gestrichen ≠ gelöscht", leere-Ablage-Konvention, Review-Report-Ausnahme —
  kein Platzhalter). **Bestätigt.**
- [x] **0 lebende `BEO-<NNN>`-Zitate außerhalb der Frozen-Bestände.** Eigener
  `grep -rn "BEO-[0-9]\+"` abzüglich `done/`, `docs/reviews/`,
  `harness/conventions/done/`, `docs/plan/cr/` liefert 18 Dateien. Davon:
  **13 ADRs** (0051, 0053, 0058, 0068–0074, 0079–0081) — jeder Treffer liegt
  vor der jeweiligen `## Geschichte`-Sektion bzw. die ADR führt gar keine,
  liegt also im **Core** und ist durch §3.5 korrekt eingefroren (die vierte,
  im Plan nachträglich benannte Frozen-Klasse). Die restlichen 5 sind
  entweder der Slice-Plan selbst (historische Erwähnung „alte BEO-024" bzw.
  Aufzählung der migrierten Alt-IDs), die neue `registerzeile-ohne-ausgang-nach-schwelle`-Beobachtung
  (referenziert „die alten BEO-008/BEO-015/BEO-020" explizit als Historie) oder
  reine Syntax-Beispiele ohne Registerzeile (`BEO-999` in `spec/lastenheft.md`
  Zeile 2478, `BEO-024`/`BEO-999` in `docs/user/benutzerhandbuch.md`
  §4.20 — laut ADR-0079 bewusst nie registrierte Illustrationen). Kein echtes
  Nacharbeits-Delta gefunden. **Bestätigt.**
- [x] Templates/Boilerplate-Träger nachgezogen: `AGENTS.md` §3.8/§5 (drei
  Fundstellen, per Diff geprüft), `MR-054` (Zeile 37, neuer Pfad),
  `.harness/skills/reviewer.md` und `closure-note-reviewer.md` (alle
  BEO-Zitate umgehängt). **Eine Unschärfe:** die DoD-Zeile nennt „ADR-0082"
  als mit-nachgezogenen Träger; `git diff HEAD~1 HEAD -- docs/plan/adr/0082-*.md`
  ist **leer** — ADR-0082 enthielt nie einen `observations.md`- oder
  `BEO-<NNN>`-Verweis und musste demnach nicht angefasst werden. Die
  Formulierung suggeriert eine Änderung, die nicht stattfand; inhaltlich
  falsch ist nichts (ADR-0082 ist tatsächlich konform), aber die DoD-Zeile
  überzeichnet ihren eigenen Beleg um eine Datei. **Bestätigt mit Anmerkung
  (LOW).**
- [x] `make gates` grün (zehn Gates) — **selbst nachgefahren**, aktueller
  Working Tree (HEAD `4d3a386`): `baseline-verify + workflow-pins +
  doc-check + lint + test + arch-check + coverage-gate + semgrep +
  gate-consistency + planning-check green`, Coverage 94,70 % ≥ 93 %,
  Semgrep 0 Findings. **Bestätigt.**
- [x] `make fullbuild` grün — **selbst nachgefahren**: `image-test` (4/4
  Fälle), `bench` (Median 706 ms < 5000 ms), `--trace --require-complete`
  (51 Anforderungen, 0 Waisen), `verify-closure-notes` (`--config
  .d-check.closure.yml --enable planning --enable structure --enable spans
  --enable reviews`) meldet **581 Dateien, 0 Befunde** — deckt sich exakt mit
  der Commit-Botschaft. Schluss: `[fullbuild] green — image-hash
  sha256:60c3fc51781e4769a42cc92502336e1e39bbfa6b09b5adf3e12b58dfa4529815`.
  **Bestätigt.**
- [ ] Unabhängiger Review — außerhalb des Verifikations-Scopes, nicht geprüft
  (Unabhängigkeits-Gebot dieses Reports).
- [x] Unabhängige Verifikation durchgeführt (dieser Report).
- [ ] **Closure-Notiz (§9): Substanz vorhanden, aber ein Risiko-Ausgang trägt
  nicht.** Siehe §Hard-Rule-Konformität unten — Risiko 3 aus §5 hat keinen der
  drei kanonischen Ausgänge. Solange das offen steht, ist dieser Haken zu
  Recht ungesetzt.

## Hard-Rule-Konformität

- **§3.5 (ADR-Immutabilität):** `make adr-check --range HEAD~1..HEAD`
  **selbst gefahren** → `649 Datei(en) geprüft, 0 Befund(e)`. Manuell
  nachvollzogen: alle vier in diesem Commit inhaltlich geänderten Accepted-ADRs
  (0055, 0074, 0075, 0078) ändern ausschließlich Zeilen **innerhalb der
  `## Geschichte`-Tabelle** (Umhängen von `BEO-<NNN>` auf den neuen Pfad in
  bereits bestehenden Chronik-Einträgen), keine Kern-Zeile ist betroffen.
  Konform.

- **§3.3/MR-013 (git mv + Inhaltsänderung = zwei Commits) — Justifikation
  nicht tragfähig, HIGH.** Die Commit-Botschaft begründet den **einzigen**
  Commit ausdrücklich mit „dieselbe Begründung wie bei MR-059 für den
  Wellen-Archiv-Stub-Move". MR-059 selbst grenzt seinen Geltungsbereich aber
  explizit ein: *„Geltungsbereich: `tools/archive-wave` (`Apply()`), jeder
  Wellen-Archivierungs-Commit"* und nennt als **Grenze** wörtlich: *„Diese
  Regel deckt nur die Wellen-Archivierung selbst (Modul 6 Schritt 4). Sie ist
  keine Blankovollmacht für beliebige Content-Move-Commits — jeder andere
  Fall prüft weiterhin gegen die Zwei-Commit-Grundregel und die drei
  MR-013-Ausnahmen."* Die Beobachtungs-Register-Migration ist keiner der vier
  in `AGENTS.md` §3.3 namentlich aufgeführten Fälle (Slice-Lifecycle-Move,
  Beanspruchung, MR-/Wellen-Lifecycle-Move, Wellen-Archiv-Stub-Move) — sie ist
  weder ein Slice- noch ein Wellen-Lifecycle-Übergang, und
  `tools/archive-wave` ist nicht beteiligt. Die praktische Beobachtung des
  Commits (keine Zwischenphase mit grünem `make gates` bei getrennten
  Commits, `git diff-tree -r --name-status -M` zeigt 0 Renames, 2 D / 149 A /
  34 M — **selbst nachgefahren**) mag in der Sache zutreffen, aber genau
  dafür sieht dieser Kanon den vorgesehenen Weg vor: eine **neue**,
  eigenständige `MR-<NNN>`-Adaption (Nachtrag zu MR-013, analog zu MR-059
  selbst) statt einer Zitat-Analogie in der Commit-Botschaft. Eine zitierte
  Quelle trägt nur, was in ihrem Geltungsbereich steht (`AGENTS.md` §5,
  Hard Rule seit slice-147) — MR-059 wörtlich gelesen schließt diese Analogie
  eher aus, als sie zu stützen. Das ist keine Formalie: `AGENTS.md` §3.3 ist
  eine Hard Rule mit geschlossener Ausnahmeliste, und dieser Commit fügt ihr
  faktisch eine fünfte, unbeschlossene Ausnahme hinzu.

- **Modul 5 §Offene Risiken werden bei Closure aufgelöst / MR-049 — Risiko 3
  hat keinen kanonischen Ausgang, HIGH.** §5 dritter Punkt: *„Der Umfang
  könnte die Ein-Sitzungs-Review-Grenze sprengen — **eingetreten**, aber
  nicht rückgeführt."* Der Kanon definiert `eingetreten` als *„→ Carveout
  (Modul 7) oder Folge-Slice mit ID"*. Beides fehlt: `docs/plan/carveouts/`
  enthält außer der `README.md` **keinen** Eintrag, kein Folge-Slice mit ID
  ist genannt. Der Text beschreibt stattdessen eine dritte, selbst gewählte
  Kompensation („gezielter Stichproben-Fokus" statt Rückführung) — das ist
  exakt die in diesem Repo bereits als Anti-Muster registrierte Klasse
  [`BEO-ALL/invented-fourth-closure-outcome`](../plan/planning/observations/BEO-ALL/invented-fourth-closure-outcome/observation.md)
  (*„Der Kanon nennt drei Ausgänge … Der vierte, erfundene Ausgang ist ein
  Absatz in der Closure-Notiz … er sieht aus wie Buchführung, ist aber
  Ablage"*). Verschärfend: §6 desselben Slice-Plans hatte den Trigger für
  genau diesen Fall selbst vorab festgelegt — *„Rückführungen: `in-progress`
  → `next`, falls §5s dritter Punkt eintritt … dann Aufteilung nach
  Sub-Area"* — und die Closure-Notiz bestätigt, dass die Bedingung eintrat,
  aber bewusst **nicht** dem eigenen Trigger gefolgt wurde. `MR-049`s
  `forbid-pattern`-Sensor greift hier **nicht**, weil er nur das
  Komplement der drei Ausgangswörter verbietet und `eingetreten` als
  Präfix formal genügt (bis Tiefe 5) — der Sensor prüft die Form, nicht ob
  „eingetreten" tatsächlich einen Carveout oder Folge-Slice trägt; das ist
  laut Kanon ausdrücklich die verbleibende **Urteils**-Hälfte, die dieser
  Report zieht. Ein Übergang nach `done/` mit diesem Stand wäre nach Modul 5
  unzulässig (*„Ein Slice geht nicht nach `done/`, während ein Risiko ohne
  Ausgang dasteht"*), solange hier nicht entweder ein Carveout eröffnet, ein
  Folge-Slice mit ID benannt oder der Punkt ehrlich als `weiter offen` ins
  Beobachtungs-Register überführt wird.

- **§3.7 (Kommentarklassen):** die neuen `<!-- d-check:ignore … -->`-Marker
  im Slice-Plan (Quell-/Ziel-Form-Begründungen) und die neuen
  Config-Kommentare in `.d-check.yml` (Tombstone-Begründung, Kopplung zu
  `FOCUS_DISABLE`) tragen jeweils Zusage/Abgrenzung/Kopplung, keine
  Review-Historie oder Deliberation. Konform. Die `state.md`-Zustandsfelder
  der migrierten Beobachtungen (Stichprobe geprüft) tragen Zustand + Anker im
  Indikativ, keine Chronik. Konform.

- **§3.8 (Modul-Scan-Grenze):** nicht berührt von diesem Slice.

## Verdikt

Die **mechanischen** DoD-Behauptungen (Verzeichnis-/Beleg-Zahlen,
Zitations-Vollständigkeit, Template-Nachzug, `make gates`/`make fullbuild`)
sind sämtlich unabhängig nachgeprüft und **korrekt** — inklusive der auf den
ersten Blick abweichenden Verzeichniszahl (29 statt 28), die sich sauber
durch die eigene Lerneintrags-Beobachtung erklärt.

**Nicht bereit für `done/`** wegen zweier substantieller Befunde:

1. Die Ein-Commit-Entscheidung (Ausnahme von §3.3/MR-013) stützt sich auf ein
   Zitat (MR-059), dessen eigener Geltungsbereich und explizite Grenzklausel
   diese Anwendung ausschließen. Die Sache mag richtig entschieden sein —
   aber nicht mit dieser Begründung, und nicht ohne die vom Kanon dafür
   vorgesehene Form (neue `MR-<NNN>`-Adaption).
2. Risiko 3 aus §5 trägt keinen der drei kanonischen Ausgänge (Carveout,
   Folge-Slice-ID oder Beobachtungs-Register-Eintrag), obwohl der Slice-Plan
   selbst für genau diesen Fall einen Trigger vorab festgelegt hatte, dem
   nicht gefolgt wurde. Das ist dieselbe Anti-Musterklasse, die dieses Repo
   bereits als `BEO-ALL/invented-fourth-closure-outcome` führt.

Beide Punkte sind vor dem Übergang nach `done/` zu klären — durch Eröffnen
eines Carveouts, Benennen eines Folge-Slice mit ID, ehrliches `weiter offen`
im Register, oder (für Punkt 1) eine eigene MR-Adaption bzw. eine
Architect-Entscheidung nach Modul 8 §Konflikt-Pfad statt einer
Zitat-Analogie in der Commit-Botschaft.
