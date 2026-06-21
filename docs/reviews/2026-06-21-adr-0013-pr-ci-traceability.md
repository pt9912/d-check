# Review — ADR-0013 (PR-/Push-CI und Traceability-Gate), Proposed

## Kopf-Metadaten

- **Review-Art:** ADR-/Design-Review (Proposed-Entwurf gegen Source
  Precedence, Hard Rules, MR-Adaptionen, doc-check-Selbstkonformität und
  faktische Korrektheit der zitierten Mechanik — **kein Verifier**:
  DoD-Abhaken und Gate-Lauf-Bestätigung sind nicht Gegenstand; kein Gate
  wurde ausgeführt).
- **Datum:** 2026-06-21
- **Gegenstand:** `docs/plan/adr/0013-pr-ci-und-traceability-gate.md`
  (Status Proposed, NEU, untrackt) und die laut Auftrag erwartete neue
  Index-Zeile in `docs/plan/adr/README.md`. Inhalt des ADR: (1) ein
  PR-/Push-CI-Workflow neben dem Tag-Release-Workflow, der `make ci`
  ruft; (2) ein Traceability-Gate in zwei Quadranten (lokaler
  `commit-msg`-Hook via `core.hooksPath` + CI-Range-Check); (3) eine
  gemeinsame Skript-Quelle `tools/trace-check.sh` mit `make`-Wrappern <!-- d-check:ignore (geplanter Pfad, ADR-0013-Folge-Slice) -->
  `trace-check`/`hooks`; (4) Kennungs-Muster deckungsgleich mit
  `.d-check.yml` (`ids`); (5) `make trace-check` bewusst nicht in
  `make gates`. Status Proposed — der ADR baut nichts, die Umsetzung ist
  einem Folge-Slice zugewiesen.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext (selbst gelesen):** Reviewer-Skill
  `.harness/skills/reviewer.md` v1.0.0; Format-Vorbild
  `docs/reviews/2026-06-20-slice-033-digest-pins.md`; das Review-Objekt
  und `docs/plan/adr/README.md`; kanonische Quellen `harness/README.md`
  (Source precedence, §Traceability rules, §Sensors, Safety/scope),
  `AGENTS.md` (§3 Hard Rules, §4 Gates, §5 Doku-Regeln),
  `harness/conventions.md` (MR-006 Referenzrichtung, MR-004/005
  Hook-Mechanik, Modus-Tabelle, Zusatzklassen), ADR-Index-README-Format
  (Schärft-Feld, Status-Lifecycle, Immutabilität); die zitierten
  Anforderungen `DC-QA-02`, `DC-FA-DIST-001` sowie `DC-FA-ID-001`,
  `DC-FA-MTX-001`, `DC-FA-CODE-001` (für die doc-check-Selbstkonformität);
  `.d-check.yml` (ids link-policy, matrix-Regeln/Status, codepaths-Roots);
  Präzedenz-ADRs `0011-digest-pins-build-gate-images.md` und
  `0007-repository-lizenz-mit.md`; Mechanik-Quellen `Makefile`,
  `.github/workflows/release.yml`, `.claude/hooks/stop-require-gates.sh`,
  `tools/gate-consistency.sh`, `internal/hexagon/core/rules/codepaths.go`.
  **Netzzugriff war nicht erforderlich** (reines Doku-/Mechanik-Review,
  alle Belege lokal). Kein Gate ausgeführt; Lauf-Wahrheit obliegt der
  getrennten Verifikation.

## Findings

### HIGH

Keine.

### MEDIUM

#### MEDIUM-1 — Neuer ADR nicht im Index und untrackt; verletzt die Hard Rule „Neue ADRs müssen den ADR-Index aktualisieren"

- **Kategorie:** MEDIUM
- **Quelle:** Hard Rule `AGENTS.md` §5 („Neue ADRs müssen den ADR-Index
  aktualisieren"); `harness/README.md` §Traceability rules („Neue ADRs
  müssen im ADR-Index ergänzt werden"); ADR-Index-Konvention
  (`docs/plan/adr/README.md`: „Neue ADRs werden in der Tabelle unten
  ergänzt")
- **Pfad:** `docs/plan/adr/README.md:26` (letzte Tabellenzeile = ADR-0012)
- **Befund:** Die ADR-Index-Tabelle endet bei ADR-0012; eine Zeile für
  ADR-0013 fehlt vollständig (`grep '0013'` über `README.md` liefert
  nichts), und die ADR-Datei selbst ist untrackt (`git status`:
  `?? docs/plan/adr/0013-pr-ci-und-traceability-gate.md`). Der Auftrag
  benennt eine „neue Index-Zeile", die im Repo nicht existiert. Damit ist
  der ADR nicht aus dem Index erreichbar: ein Konsument, der dem
  ADR-Index als kanonischer Source-Precedence-Quelle (`AGENTS.md` §2
  Rang 4) folgt, sieht ADR-0013 nicht und referenziert eine in der Liste
  fehlende Entscheidung — genau die Traceability-Lücke, die der ADR
  selbst zu schließen beansprucht. Die Hard Rule ist verletzt, kein Gate
  fängt das ab (`doc-check` prüft Links/Anker, nicht die Vollständigkeit
  der Index-Tabelle).
- **Verifizierbar:** ja — `grep -n '0013' docs/plan/adr/README.md`
  (leer) und `git status --porcelain docs/plan/adr/` (zeigt die untrackte
  ADR-Datei, keine README-Änderung).

#### MEDIUM-2 — Innerer Widerspruch zur Bindung von `make trace-check`: Punkt 5 sagt „`make ci` ruft beide", Punkt 3 lässt die CI `make ci` + `make trace-check` getrennt rufen

- **Kategorie:** MEDIUM
- **Quelle:** ADR-0013 §Entscheidung Punkt 3 vs. Punkt 5 (Selbst-Konsistenz
  des „zwei Bindepunkte, eine Wahrheit"-Designs); `DC-FA-DIST-001`
  (`make ci`-Vertrag) als bezogener Kontext
- **Pfad:** `docs/plan/adr/0013-pr-ci-und-traceability-gate.md:80` ↔ `:99`
- **Befund:** Punkt 5 entscheidet wörtlich „Zwei Bindepunkte, zwei Targets
  — `make ci` ruft beide" (Zeile 99), d. h. `make ci` umfasst Gates **und**
  `make trace-check`. Der Skizzen-Codeblock in Punkt 3 schreibt dagegen
  `.github/workflows/ci.yml # ruft `make ci` + `make trace-check` (Range)`
  (Zeile 80), d. h. der Workflow ruft `make ci` **und zusätzlich** noch
  `make trace-check` — was nur stimmig ist, wenn `make ci` trace-check
  gerade **nicht** enthält. Beide Aussagen können nicht gleichzeitig
  gelten: entweder ist trace-check in `make ci` (dann ist der separate
  Aufruf im Workflow eine Doppel-Ausführung des Range-Checks), oder es ist
  nicht enthalten (dann ist Punkt 5 falsch). Der Folge-Slice erbt eine
  widersprüchliche Entscheidung über den zentralen Bindepunkt; eine
  Implementierung, die Punkt 5 folgt, lässt den Range-Check zweimal laufen,
  eine, die Punkt 3 folgt, widerlegt die Kern-These „`make ci` ruft beide".
  Der heutige `make ci` (`Makefile:128`: `ci: gates image-test`) enthält
  trace-check nicht — die Entscheidung, ob das so bleibt, lässt der ADR
  offen statt sie zu treffen (Aufgabe eines Proposed-ADR).
- **Verifizierbar:** ja — Lektüre der beiden ADR-Zeilen genügt; gegen
  `Makefile:128` (`ci: gates image-test`) zeigt sich, dass beide
  Lesarten heute nicht erfüllt sind und der ADR die Wahl nicht festlegt.

### LOW

#### LOW-1 — Bewachungs-Zusage benennt nur `harness/README.md §Sensors`; `make gate-consistency` verlangt neue Targets zusätzlich in `AGENTS.md` §4, sonst rotes `make gates`

- **Kategorie:** LOW
- **Quelle:** Maintainability (Reviewer-Skill LOW-Anker „latente
  Wartungsfalle, die erst bei künftigem Edit zündet"); ADR-0013
  §Konsequenzen; `tools/gate-consistency.sh` (Richtung 2)
- **Pfad:** `docs/plan/adr/0013-pr-ci-und-traceability-gate.md:116`
- **Befund:** Der ADR sagt zu, `make trace-check`/`make hooks` würden
  „in `harness/README.md` §Sensors dokumentiert und von
  `make gate-consistency` bewacht (sonst Harness-Lüge)". `gate-consistency`
  prüft jedoch **zwei** Richtungen: Richtung 1 (Doku→Makefile) liest
  `AGENTS.md` UND `harness/README.md` (`gate-consistency.sh:81`), Richtung 2
  (Makefile→Doku) liest **ausschließlich** `AGENTS.md` §4
  (`gate-consistency.sh:84`) und feuert für jedes Makefile-Target, das dort
  fehlt (`:88`). Ein Folge-Slice, der die zwei neuen Targets nur in
  `harness/README.md` §Sensors einträgt (wie der ADR sie anleitet), läuft
  in einen roten `make gate-consistency` und damit rotes `make gates` —
  der ADR-Leitfaden ist unvollständig und führt zur entgegengesetzten
  Wirkung der „sonst Harness-Lüge"-Absicht.
- **Verifizierbar:** ja — sobald ein Target im Makefile steht, aber in
  `AGENTS.md` §4 fehlt, gibt `bash tools/gate-consistency.sh` „FAIL —
  Makefile-Target '<name>' fehlt in AGENTS.md §4" mit Exit 1 aus.

#### LOW-2 — `commit-msg`-Hook als Host-Check außerhalb `make` und am `make trace-check`-Wrapper vorbei steht in Spannung zu §3.1 „Alle Checks laufen über `make`" und zur eigenen „eine Wahrheit"-Klausel

- **Kategorie:** LOW
- **Quelle:** Hard Rule `AGENTS.md` §3.1 („Alle Checks laufen über `make`";
  Host braucht nur `git`, GNU `make`, `bash`, Docker); ADR-0013
  §Entscheidung Punkt 3 (`make`-Target „umschließt" das Skript)
- **Pfad:** `docs/plan/adr/0013-pr-ci-und-traceability-gate.md:79`
- **Befund:** Der ADR hängt einen `commit-msg`-Hook über
  `git config core.hooksPath` ein, der laut Skizze
  `.githooks/commit-msg # ruft tools/trace-check.sh <message-file>`
  (Zeile 79) das Prüf-Skript **direkt** aufruft — nicht über
  `make trace-check`. Damit läuft erstens ein Check zur Commit-Zeit
  außerhalb von `make` (Spannung zu §3.1 „Alle Checks laufen über `make`";
  bash/git sind zwar erlaubte Host-Werkzeuge, der Pfad geht aber bewusst
  am `make`-Vertrag vorbei), und zweitens umgeht der Hook gerade den
  `make`-Wrapper, den Punkt 3 als Auditierbarkeits-Argument („Ein
  `make`-Target umschließt es") einführt. Das geteilte Skript hält die
  „eine Skript-Wahrheit", aber die „eine `make`-Wahrheit" gilt für den
  Hook-Pfad nicht; der ADR benennt diese Asymmetrie nicht. Failure-Szenario:
  ein künftiger Edit ändert das `make trace-check`-Target (z. B. zusätzliche
  Flags/Normalisierung), der direkt aufgerufene Hook zieht die Änderung
  nicht mit — Hook- und CI-/lokales-make-Verhalten divergieren still.
- **Verifizierbar:** ja — sobald `make trace-check` Logik trägt, die das
  Skript nicht selbst kapselt, divergiert ein `git commit` (Hook ruft
  `tools/trace-check.sh` direkt) von `make trace-check` reproduzierbar. <!-- d-check:ignore (geplanter Pfad, ADR-0013-Folge-Slice) -->

### INFO

#### INFO-1 — Faktenangabe „12/15 Commits mit ID" weicht vom Repo-Stand ab; `slice-NNN` ist kein gültiges Muster, mehrere reale Commits würden geblockt

- **Kategorie:** INFO
- **Quelle:** Faktische Korrektheit (ADR-0013 §Kontext Punkt 1);
  ADR-0013 §Entscheidung Punkt 4 (Muster `ADR-\d{4}` / `DC-(FA-[A-Z]+|QA)-\d+`)
- **Pfad:** `docs/plan/adr/0013-pr-ci-und-traceability-gate.md:30`
- **Befund:** Der ADR behauptet „zuletzt 12/15 Commits mit ID-Token". Im
  aktuellen 15-Commit-Fenster tragen nur **6/15** ein ID-Token in der
  Subject-Zeile und **11/15** im vollen Message-Body
  (`ADR-NNNN`/`DC-…`-Token); die Differenz 11 vs. 12 ist klein und für
  ein anderes Audit-Fenster (Stand 2026-06-21) plausibel, daher nicht
  REFUTED, aber unpräzise. Materiell relevant: vier reale jüngere Commits
  tragen **kein** gültiges Token (z. B. `docs(handbuch): …`,
  `docs(plan): slice-036 …`), und `slice-NNN` ist im Muster aus Punkt 4
  bewusst nicht enthalten — diese Commit-Klasse würde der vorgeschlagene
  Gate blockieren. Der ADR-Kontext stellt die Disziplin als „gelebt" dar,
  ohne diese ID-lose Klasse zu beziffern; die Ausnahmeliste (Punkt 4)
  deckt nur Merge-/Revert-Commits, nicht reine `slice-NNN`-/Doku-Commits.
- **Verifizierbar:** ja —
  `git log -15 --format='%B'` je Commit gegen
  `grep -E 'ADR-[0-9]{4}|DC-(FA-[A-Z]+|QA)-[0-9]+'` zählt 11/15 (Body)
  bzw. 6/15 (Subject); die NO-ID-Commits sind im Log direkt sichtbar.

#### INFO-2 — Branch-Protection-Restlücke und shallow-clone-Range ehrlich als offene Punkte benannt (bewusste Won't-Fix-/Folge-Designnotiz)

- **Kategorie:** INFO
- **Quelle:** ADR-0013 §Konsequenzen / §Re-Evaluierungs-Trigger;
  Reviewer-Skill INFO-Anker (bewusste, dokumentierte Designnotiz)
- **Pfad:** `docs/plan/adr/0013-pr-ci-und-traceability-gate.md:119`
- **Befund:** Der ADR benennt die zentrale Restlücke selbst: ohne extern
  konfigurierte Branch Protection ist die PR-CI nur *advisory* (Zeilen
  119–123), und der shallow-clone-Range steht als Re-Evaluierungs-Trigger
  (Zeile 145–146). Das ist die korrekte Behandlung einer nicht aus dem
  Klon auditierbaren Grenze (analog der Stop-Hook-Restlücke in
  `.claude/hooks/stop-require-gates.sh:9-11`, auf die der ADR-Kontext sich
  beruft) — dokumentationswürdig als bewusst offener Folgepunkt, kein
  Mangel des Proposed-Entwurfs.
- **Verifizierbar:** ja — die Restlücken-Absätze stehen explizit im ADR;
  ein Failure-Szenario tritt nur ein, wenn der Folge-Slice die externe
  Konfiguration nicht als Betriebshinweis dokumentiert (dort zu prüfen).

## Negativbefunde (geprüft, ohne Befund)

- **doc-check / `ids` (`DC-FA-ID-001`, link-policy `always`):** REFUTED,
  dass der ADR `ids` reißt. `docs/plan/adr/` ist **nicht** in den
  `exempt-paths` (`.d-check.yml`: nur `CHANGELOG.md`, `docs/reviews/**`),
  also greift `link-policy: always` voll. Alle realen ID-Token außerhalb
  von Fenced-Blöcken (`grep` auf Nicht-Fence-Zeilen: Zeilen 6/8/10/11/15/56/126/129)
  stehen als Markdown-Link (`[ADR-0011](…)` bzw. `[`DC-QA-02`](…)`) und
  erfüllen damit `always`; die Heading-Zeile (`ADR-0013`, Zeile 1) ist
  per Spec von der Linkpflicht ausgenommen. Die ID-Regexe in den
  Fenced-Blöcken (Punkt 4) sind per `DC-FA-ID-001` („außerhalb … von
  Fenced-Code-Blöcken"; `always` „Fenced-Code-Blöcke … bleiben frei")
  ausgenommen.
- **doc-check / `codepaths` (`DC-FA-CODE-001`):** REFUTED, dass die noch
  nicht existierenden Pfade ein `codepath-missing` erzeugen. Die drei
  Pfade `tools/trace-check.sh`, `.githooks/commit-msg`, <!-- d-check:ignore (geplanter Pfad, ADR-0013-Folge-Slice) -->
  `.github/workflows/ci.yml` stehen im **Fenced-Block** (Punkt 3); die
  Anforderung schließt das aus („Vorkommen in Fenced-Code-Blöcken werden
  nicht geprüft"; Out-of-Scope „Pfade in Fenced-Code-Blöcken"), bestätigt
  durch die Implementierung `internal/hexagon/core/rules/codepaths.go`
  (arbeitet auf `proseLines`/`inlineSpansByLine`, nicht auf
  Fence-Inhalten). `tools/trace-check.sh` läge zwar unter dem <!-- d-check:ignore (geplanter Pfad, ADR-0013-Folge-Slice) -->
  codepaths-Root `tools` — wäre es in einem Inline-Span, gäbe es einen
  Befund; im Fence ist es frei. Kein Inline-Code-Pfad außerhalb der
  Fences zeigt auf ein fehlendes Ziel.
- **doc-check / Links + Anker:** Alle drei Anker des ADR resolven gegen
  reale Heading-Slugs: `#dc-qa-02--determinismus` ←
  `### DC-QA-02 — Determinismus` (`spec/lastenheft.md:741`),
  `#dc-fa-dist-001--docker-image` ← `### DC-FA-DIST-001 — Docker-Image`
  (`:715`), `#traceability-rules` ← `## Traceability rules`
  (`harness/README.md:73`). Alle Link-Ziele existieren
  (`spec/lastenheft.md`, `harness/README.md`, die vier referenzierten
  ADR-Dateien, `.github/workflows/release.yml`,
  `.claude/hooks/stop-require-gates.sh`, `tools`). Das Prosa-„§Sensors"
  (Zeile 117) trägt **keinen** `#…`-Anker, also kein Bruch trotz
  abweichendem Slug (`sensors-feedback-gates`).
- **doc-check / `matrix` (`DC-FA-MTX-001`, status-forbidden):** REFUTED,
  dass der ADR `matrix-inactive`/`matrix-forbidden` reißt. Der ADR ist
  Klasse `adr` (`.d-check.yml`: `docs/plan/adr/[0-9]*.md`); die
  `from`-Regeln verbieten nur `spec-straten → adr/slice`, nicht
  `adr → *`. Die Status-Prüfung (`forbidden: [superseded, deprecated]`)
  betrifft alle Referenzen — der ADR verlinkt aber ausschließlich
  ADR-0011/0002/0007, alle mit `**Status:** Accepted`; ADR-0004 (im Index
  als „superseded"-tragende Statuszelle) wird **nicht** verlinkt. Eine
  ADR→Slice-Pflicht existiert nicht (`matrix` verbietet Kanten, fordert
  keine).
- **Referenzrichtung / MR-006 (`Bezug:`/`Schärft:`):** Der ADR verweist im
  bindenden Text **aufwärts** (auf `DC-QA-02`, `DC-FA-DIST-001`, aktive
  ADRs, `harness/README.md`); kein Spec-Stratum verweist hier abwärts —
  MR-006 betrifft `spec/*.md`, nicht ADRs, und ist nicht berührt. Der ADR
  greift mit Dateinamen/Skriptpfaden/Regex in Slice-Detail vor, kennzeichnet
  das aber explizit als Skizze (Punkt 3 Codeblock) und delegiert Bau + DoD
  an den Folge-Slice (§Konsequenzen Zeile 113–115) — Scope-Disziplin
  gewahrt, kein Vorgriff im entscheidenden Text.
- **`Schärft: keine Spec-Stelle` (Ehrlichkeit):** Geprüft, ob der ADR
  eigentlich `spec/spezifikation.md` schärfen müsste. Die durchgesetzte
  Regel („Commits nennen eine ID") lebt in `harness/README.md`
  §Traceability rules, nicht in einem Spec-Stratum; `DC-QA-02`
  (Determinismus) ist ein Produktvertrag, nicht die Traceability-Regel.
  „keine Spec-Stelle — Prozess-/Durchsetzungs-ADR" ist damit ehrlich und
  durch ADR-0011 (Zeile 10: „keine Spec-Stelle — Prozess-…-ADR") und
  ADR-0007 (Zeile 7: „keine Spec-Stelle — Projektmetadaten-ADR")
  präzedenziert; der nicht-DC-`Bezug:` auf `harness/README.md` ist durch
  dieselbe ADR-0007-Linie (nicht-DC-`Bezug:` „Repository-Veröffentlichung
  …") gedeckt.
- **Faktencheck zur zitierten Mechanik:** `make ci` (gates + image-test)
  **existiert** (`Makefile:128`) und ist die Release-Engine
  (`release.yml:70` ruft `make ci`); `release.yml` ist **tag-only**
  (`release.yml:20-23`: `on: push: tags: ['v*']`, kein `pull_request`/
  `push: branches`); der Stop-Hook **gibt cleane Klone frei**
  (`stop-require-gates.sh:28-34`: kein State + leerer Tree ⇒
  `{"decision":"approve"}`); `make gate-consistency`/`make image-test`
  existieren (`Makefile:73`/`:79`). Die behaupteten Dateipfade
  (`release.yml`, Stop-Hook) sind real. Einzige Abweichung: die
  ID-Quote-Angabe (INFO-1).
- **Status-Lifecycle / Immutabilität (`AGENTS.md` §3.5, Index-Konvention):**
  Status `Proposed`, korrekt für einen Entwurf; keine bestehende
  Accepted-ADR wird editiert (`git status` zeigt nur die untrackte
  ADR-0013-Datei, keine Änderung an `0001`–`0012`). Kein Immutabilitäts-Bruch.
- **CI-Trigger / Doppelläufe:** Punkt 1 schließt Tag-Läufe explizit aus
  (`push` auf Branches **ohne** Tags; Tag-Läufe bleiben `release.yml`
  vorbehalten) — die offensichtliche Doppellauf-Falle (Tag-Push triggert
  beide Workflows) ist im Entscheidungstext adressiert. (Die getrennte
  trace-check-Doppelung im PR-Pfad ist MEDIUM-2, ein anderer Mechanismus.)

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 2 | 2 | 2 |

## Verdikt

**Vor Annahme zu klären** (MEDIUM-1, MEDIUM-2) — als Proposed
**noch nicht merge-/accept-fähig**. Der Entwurf ist inhaltlich tragfähig:
Referenzrichtung, `Schärft: keine Spec-Stelle` und der nicht-DC-`Bezug:`
sind ehrlich und durch ADR-0011/ADR-0007 präzedenziert; die
doc-check-Selbstkonformität (`ids` `always`, `codepaths`, `matrix`,
Anker/Links) ist sauber (Fenced-Blöcke und Heading nehmen die ID-Regexe
und die Beispiel-Pfade korrekt aus); die zitierte Mechanik stimmt
faktisch (`make ci` existiert, `release.yml` tag-only, Stop-Hook gibt
cleane Klone frei, gate-consistency/image-test existieren); die
Doppellauf-Falle (Tag) und die Branch-Protection-/shallow-clone-Lücken
sind adressiert (INFO-2). Es blockieren zwei Punkte: **MEDIUM-1** — der
ADR ist weder im Index ergänzt noch überhaupt getrackt, was die Hard
Rule „Neue ADRs müssen den ADR-Index aktualisieren" (`AGENTS.md` §5,
`harness/README.md` §Traceability rules) verletzt und den ADR aus dem
kanonischen Source-Precedence-Pfad ausblendet (kein Gate fängt das ab —
genau die aspirativ-statt-bindend-Lücke, die der ADR adressieren will);
und **MEDIUM-2** — der zentrale „zwei Bindepunkte, eine Wahrheit"-Entwurf
widerspricht sich selbst über die Bindung von `make trace-check`
(Punkt 5 „`make ci` ruft beide" vs. Punkt 3 „CI ruft `make ci` +
`make trace-check`"), sodass der Folge-Slice eine unentschiedene
Kernfrage erbt. LOW-1 (gate-consistency verlangt neue Targets in
`AGENTS.md` §4, der ADR-Leitfaden nennt nur `harness/README.md` §Sensors)
und LOW-2 (`commit-msg`-Hook umgeht `make` und den eigenen
`make`-Wrapper) sind latente Wartungsfallen für die Umsetzung; INFO-1
(ID-Quote unpräzise; `slice-NNN`-Commits würden geblockt) und INFO-2
(bewusst delegierte Restlücken) sind dokumentationswürdig. Gate-/CI-Lauf
und DoD obliegen der getrennten Verifikation (hier NICHT als grün
angenommen; kein Gate ausgeführt).
