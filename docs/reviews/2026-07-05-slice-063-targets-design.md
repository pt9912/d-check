# Review-Report — slice-063 (Modul `targets`) — Design-/Entwurf-Review

**Datum:** 2026-07-05
**Reviewer-Rolle:** unabhängig/adversarial, Fokus **Design-/Entwurf-Stand** vor
dem Doc-first-Fundament (die vom Slice erbetene „Bestätigung"). Der Slice ist
ENTWURF und **untracked** (`??`); Lastenheft-CR `DC-FA-TGT-001`, Spezifikation,
ADR-0031 und der Roadmap-Flip stehen laut Slice noch aus.
**Gegenstand:** `docs/plan/planning/in-progress/slice-063-targets-modul.md`
gegen das abzulösende `tools/gate-consistency.sh`, den `.d-check.yml`-Bestand,
`spec/lastenheft.md` (0.37.1) / `spec/spezifikation.md` §4, den ADR-Index und die
`--print-mk`-/DC-FA-CLI-010-Lage.
**Baseline:** `.harness/skills/reviewer.md` v1.2.0, `AGENTS.md` §3,
`harness/conventions.md` (MR-007, MR-013, MR-014).
**NICHT geprüft:** DoD-Abhakung (Verifikations-Rolle; es gibt noch keinen Impl-Diff);
a-check-Repo-interner Stand (eigener Harness drüben; die Drift-Behauptung ist von
hier aus nicht auditierbar — siehe INFO-3).

**Verifikations-Proben (Quellen-Lektüre statt Gate-Läufe — es gibt noch keinen Diff):**

- Abgelöstes Skript, Doku-Seite: `tools/gate-consistency.sh:26` — `doc_targets`
  scannt **nur Tabellenzeilen** (`grep -E '^\|'`) auf `` `make X` ``-Tokens;
  `make X` in Prosa (z. B. `AGENTS.md` §3.1 „**Richtig:** `make gates`", §3.3
  `make planning-check`, §3.5 `make adr-check`) zählt **nicht**.
- Abgelöstes Skript, Makefile-Seite: `tools/gate-consistency.sh:34-37` — Regel-
  zeilen `^[a-zA-Z][…]:([^=]|$)`; `.PHONY:`/Pattern-Rules (`%`) fallen über den
  Anfangsbuchstaben heraus, Zuweisungen über das `=`-Negat.
- Richtungen: Prüfung 1 (`:81`) läuft über `AGENTS.md` **und** `harness/README.md`;
  Prüfung 2 (`:84-91`) nur gegen `doc_targets AGENTS.md` — beide über
  **ganze Dateien**, nicht §4/§Sensors-scoped; **keine** exempt-Liste in Prüfung 2.
- Prüfung 3 (`:97-107`) = `.d-check.yml`-Modulliste (DC-QA-03) — der Slice
  belässt sie als repo-lokalen Rest (§2/§4), konsistent mit dem Skript.
- Reason-Code-Bestand (`spec/spezifikation.md` §4): `target-missing` (`links`,
  `:1153/1401`), `target-untracked` (`tracked`, `:1164/1423`) — „target" = **Link-/
  Referenz-Ziel**. Kein `target-phantom`/`target-undocumented`/`gate-*` vorhanden.
- `--print-mk`: Happy Path in `DC-FA-CLI-010` (`lastenheft.md:415`) listet 10
  Targets inkl. `doc-help` (`doc-tracked` mit 9→10 in `slice-059`); 10→11 stimmt.
- ADR-Index: `docs/plan/adr/0030-tracked-referenz-ziele.md` zuletzt — ADR-0031
  ist die nächste freie Nummer.
- `.d-check.yml` `ids` (`:24-37`): `link-policy: always` für `ADR-\d{4}`,
  `MR-\d{3}`, `DC-(FA-[A-Z]+|QA)-\d+`; exempt nur `CHANGELOG.md` + `docs/reviews/**`
  — Slices unter `docs/plan/planning/` sind **nicht** ausgenommen.
- Roadmap: §Aktuelle Welle trägt noch „Keine aktive Welle", während
  `slice-063-targets-modul.md` bereits in `in-progress/` liegt (untracked).

---

## Findings

### MEDIUM-1 — Tabellen-Scoping der `make X`-Erkennung ist fidelity-kritisch und unspezifiziert

- **kategorie:** MEDIUM
- **quelle:** `DC-QA-04` (Erkennungs-Differenz zur Alt-Mechanik)
- **pfad:** `docs/plan/planning/in-progress/slice-063-targets-modul.md:46-55`,
  `:56-63` (Config)
- **befund:** Das abgelöste Skript scannt Doku-Behauptungen **nur in
  Tabellenzeilen** (`gate-consistency.sh:26`, `^|`). `AGENTS.md`/`harness/README.md`
  tragen zahlreiche `make X` in **Prosa** (`AGENTS.md` §3.1/§3.3/§3.5). Der
  Config-Key `doc-tables` und die §2-Prosa sagen „Doku-**Tabelle**", benennen aber
  **keinen** Tabellen-Erkennungs-Mechanismus; wird Richtung 1 whole-file statt
  tabellen-zeilen-scoped umgesetzt, überdeckt sie Prosa-Erwähnungen.
- **Failure-Szenario:** Eine Prosa-Zeile nennt ein umbenanntes/entferntes Target
  (heute z. B. der `<!-- d-check:ignore -->`-annotierte historische `make doc-check`-
  Verweis-Stil) → spuriöser `target-phantom` auf dem sauberen Baum, obwohl das
  Skript grün ist. „Fidelity ≡ Skript" (§2 Z.71-72) bricht.
- **verifizierbar:** ja — Modul-Lauf vs. `make gate-consistency` auf demselben Baum;
  die Spec muss die `^|`-Scoping-Regel festschreiben.

### MEDIUM-2 — DoD ohne expliziten Paritäts-/Regressionsnachweis gegen das abgelöste Skript

- **kategorie:** MEDIUM
- **quelle:** `DC-QA-04`; Präzedenz slice-058 (20er-Proben-Matrix je
  Verbotszweig), `vcs`/`commits` (Paritäts-Belege)
- **pfad:** `docs/plan/planning/in-progress/slice-063-targets-modul.md:108-110`
  (DoD), `:71-72` (§2 „Fidelity ≡ Skript")
- **befund:** Der Slice **stilllegt ein aktives Gate** und stützt sich auf „keine
  Regression". Die DoD listet Happy/Boundary/Negative je Richtung + Negativ-
  Selbsttest, aber **keinen** „Modul-Befundsatz ≡ Skript-Befundsatz auf dem
  aktuellen Repo"-Beleg. Die Vorgänger-Migrationen trugen genau diesen Nachweis.
- **Failure-Szenario:** Das Modul fängt etwas nicht, das das Skript fing (oder
  umgekehrt); die Differenz reist unbemerkt mit dem Release, „Parität belegt" ist
  behauptet ohne belegt zu sein.
- **verifizierbar:** ja — Paritäts-Testlauf gegen die aktuelle Skript-Ausgabe in
  die DoD aufnehmen.

### LOW-1 — Reason-Code-Nomen „target" ist semantisch überladen

- **kategorie:** LOW · **quelle:** Maintainability (öffentlicher Vertrag)
- **pfad:** `docs/plan/planning/in-progress/slice-063-targets-modul.md:50/54`,
  `:136-138` (§4-Notiz)
- **befund:** `target-missing` (`links`) / `target-untracked` (`tracked`) meinen
  **Link-/Referenz-Ziel**; `target-phantom`/`target-undocumented` meinen
  **Makefile-/Build-Target**. Gleicher Reason-Code-Namensraum, zwei Bedeutungen von
  „target". Der Slice flaggt es („bewusst nur im Präfix; zu bestätigen"),
  unterschätzt es aber: es ist eine **Bedeutungs**-Kollision, keine bloße
  Präfix-Überlappung. Das Modul rahmt sich als Meta-Gate — `gate-*`/`make-target-*`
  wäre eindeutiger.
- **Failure-Szenario:** Ein Nutzer triagiert `--doctor`/`--json` und liest
  `target-phantom` als kaputtes Link-Ziel statt als halluziniertes Gate; nach
  Release ist der Code schwer umzubenennen.
- **verifizierbar:** ja — Spec §4-Grund-Code-Tabelle beim Doc-first festlegen.

### LOW-2 — „§4"/„§Sensors"-Section-Scoping fehlt im Config-Schema

- **kategorie:** LOW · **quelle:** Maintainability
- **pfad:** `docs/plan/planning/in-progress/slice-063-targets-modul.md:61`,
  `:131-135`
- **befund:** `authority: AGENTS.md` und `doc-tables: [AGENTS.md, …]` sind **ganze
  Dateien**; es gibt keinen Heading-Key (anders als `planning.heading`). Das Skript
  scannt ebenfalls whole-file-Tabellenzeilen — also fidelity-erhaltend — aber die
  wiederholte „§4"/„§Sensors"-Prosa verspricht ein Section-Scoping, das die Config
  nicht liefert.
- **Failure-Szenario:** Ein künftig **außerhalb** §4 eingefügter
  `| … make X … |`-Tabelleneintrag befriedigt Richtung 2 stillschweigend, obwohl
  §4 die deklarierte „nur hier gelistete"-Autorität ist.
- **verifizierbar:** ja — Section-Anker aufnehmen **oder** die „§4"-Präzision aus
  der Spec-Prosa streichen.

### LOW-3 — `exempt-targets: [help, clean]` ist für d-checks eigene Config irreführend

- **kategorie:** LOW · **quelle:** Maintainability
- **pfad:** `docs/plan/planning/in-progress/slice-063-targets-modul.md:62`
- **befund:** `help`/`clean` sind heute in `AGENTS.md` §4 dokumentiert
  (`:160`) — sie zu exemten ist unnötig und schwächt Richtung 2. Die DoD sollte
  festhalten, dass d-checks eigenes `targets.exempt-targets` **minimal/leer** bleibt
  (die heutige erschöpfende Makefile→§4-Prüfung zu erhalten); das `[help, clean]`
  im Schema ist Konsumenten-Illustration, nicht d-checks Config.
- **verifizierbar:** ja — `.d-check.yml`-Review nach Umsetzung.

### LOW-4 — DoD lässt die selbstbezüglichen Sensors-/§4-Doku-Updates für das geänderte `gate-consistency` aus

- **kategorie:** LOW · **quelle:** MR-007 / Harness-Ehrlichkeit
- **pfad:** `docs/plan/planning/in-progress/slice-063-targets-modul.md:111-115`
- **befund:** Beim Wechsel `make gate-consistency` (Skript → `targets`-Dogfood)
  ändern sich die Zeilen in `AGENTS.md` §4 **und** `harness/README.md` §Sensors
  (wie bei planning/commits/vcs). Pikant: `targets` **prüft genau diese
  Tabellen** — Doku-Edits und Modul-Lauf müssen wechselseitig konsistent sein. Die
  DoD nennt Handbuch §5/§6, aber nicht diese zwei Tabellenzeilen.
- **verifizierbar:** ja — `make gate-consistency` nach der Umstellung.

---

## Findings (INFO)

### INFO-1 — Forward-Referenzen + bare Tokens vs. `ids` beim Commit

- **kategorie:** INFO · **quelle:** `DC-FA-ID-001`
- **pfad:** `docs/plan/planning/in-progress/slice-063-targets-modul.md:13/15/107`
- **befund:** `DC-FA-TGT-001`/`ADR-0031` existieren noch nicht (bis Doc-first sie
  anlegt, nicht verlinkbar). Die **fett, nicht-backtickten** `**ADR-0031**`
  (Z.13/15/107) sind bare Prosa-Tokens; `ids` (`.d-check.yml:24-37`,
  `link-policy: always`, Slices **nicht** exempt) verlangt sie beim Commit
  backtickt-oder-verlinkt, sonst `id-unlinked` in `make doc-check`. Praktischer
  Heads-up, kein Design-Mangel.

### INFO-2 — Transienter `planning-drift`, solange der Entwurf in `in-progress/` liegt

- **kategorie:** INFO · **quelle:** `DC-FA-PLAN-001` / MR-013
- **pfad:** `docs/plan/planning/in-progress/slice-063-targets-modul.md:3-5`;
  `docs/plan/planning/in-progress/roadmap.md:13`
- **befund:** Die Roadmap sagt noch „Keine aktive Welle", der Slice liegt aber in
  `in-progress/`. Korrekt gemäß der bewussten „Roadmap-Flip nach Bestätigung"-
  Staffelung (Datei untracked ⇒ committed-tree-`planning-check` bleibt grün). Bei
  Annahme muss der Flip **atomar** im ersten Slice-Commit landen (MR-013), sonst
  ist der erste CI-Push rot.

### INFO-3 — Cross-repo-Drift-Beleg nicht aus diesem Repo verifizierbar

- **kategorie:** INFO · **quelle:** MR-007-Begründung
- **pfad:** `docs/plan/planning/in-progress/slice-063-targets-modul.md:24-26`,
  `:144-145`
- **befund:** Die motivierende Evidenz (a-checks `gate-consistency.sh` hat
  Pin-Block/Utility-Allowlist/`d-check.mk`-Quelle, d-checks nicht) ist von hier aus
  nicht prüfbar. Plausibel und konsistent mit dem a-check-Stand, aber tragend für
  „verteilen statt kopieren" — ein Datei-/Commit-Pointer auf die a-check-Fassung in
  der ADR-0031-Motivation macht den Drift auditierbar.

---

## Negativbefunde (geprüft, ohne Befund)

- **ADR-Nummer:** ADR-0031 ist die nächste freie (0030 zuletzt im Index). Ohne Befund.
- **Namensraum:** Modulname `targets` + Config-Key `targets.*` kollidieren mit
  keinem der 16 bestehenden Module/Keys; „17. Regelmodul" stimmt. Ohne Befund.
- **`--print-mk` 10→11:** arithmetisch korrekt (aktuell 10 inkl. `doc-help`,
  belegt in `DC-FA-CLI-010` Happy Path `lastenheft.md:415`). Ohne Befund.
- **Versions-Bump 0.37.1 → 0.38.0:** Minor passt für ein neues nutzersichtbares
  opt-in-Modul (Präzedenz slice-057/059). Ohne Befund.
- **Hexagon §6:** kein neuer Port — `planning` etablierte bereits den
  configured-path-`ReadFile` außerhalb des Markdown-Scans; die `Makefile` liegt im
  read-only-Mount. Ohne Befund.
- **`planning`↔`targets`-Analogie:** hermetisch, read-only, Doku-Behauptung ↔
  Repo-Struktur — trägt und begründet die „d-check-fähig"-Neubewertung korrekt.
  Ohne Befund.
- **Scope-Ausschluss Prüfung-3** (Modul-Listen-Selbstkonsistenz): sauber begründet
  (repo-lokal, nicht der cross-repo-Drift-Kern). Ohne Befund.
- **Relative Links + MR-Anker** im Slice lösen auf korrekter Tiefe auf
  (`../../../../` → Repo-Wurzel; MR-007/MR-013-Slugs korrekt). Ohne Befund.
- **Referenz-Richtung (SDP)/Marker-Ehrlichkeit:** kein Provenance-Marker im
  Gegenstand; keine getarnte Abwärts-Entscheidungsgrundlage. Ohne Befund.

---

## Kategorie-Summary

| Kategorie | Anzahl |
| --- | --- |
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 3 |

## Verdikt

**Design bestätigt — `targets` ist d-check-fit.** Die Kern-These (der cross-repo-
**gemeinsame** Kern von `gate-consistency.sh` ist „Doku-Behauptung über
Repo-Realität", also d-checks Mission) trägt, die `planning`-Analogie ist präzise,
Scope-Disziplin (Prüfung-3 draußen) und Hexagon-Einordnung stimmen. Keine HIGH,
keine Design-Blocker.

Vor dem Doc-first-Fundament zu klären: **MEDIUM-1** (die `^|`-Tabellen-Scoping-Regel
in die Spec schreiben — sonst bricht die Fidelity-These) und **MEDIUM-2**
(Paritätsnachweis in die DoD — ein stillgelegtes Gate braucht einen „≡ Skript"-Beleg
wie die Vorgänger-Migrationen). **LOW-1** (Reason-Code-Benennung) ist ein
öffentlicher Vertrag und billiger jetzt als nach Release. LOW-2 bis LOW-4 in
denselben Doc-first-Durchgang einarbeiten. Der Slice ist ENTWURF und darf
entsprechend nachgezogen werden.
