# Review-Report — slice-244 (R1)

**Review-Art:** Code/Design (Adoption `.claude/`-Rollen und Commands — geprüft gegen Slice-Plan, Schwester-Snapshot `/tmp/aih-v6.13.0`, Hard Rules `AGENTS.md` §3; *nicht* gegen die DoD — das ist die Frage des Verifiers)

**Gegenstand:** `slice-244` (welle-91), Range `e0165e8c..HEAD` — Commits `c7091b7b` (Beanspruchung: Ruhe-Marker verlässt die Roadmap, reiner `git mv` des Plans) und `cb918ee3` (Adoption: 5 Dateien neu — `.claude/agents/{architect,planner,validator}.md`, `.claude/commands/{plan-welle,close-welle}.md` — plus Roadmap-Fortschreibung)

**Skill:** `reviewer.md` @ 1.16.0 · **Modell-ID:** glm-5.3-flash · **Datum:** 2026-09-29

**Eingangs-Kontext:** Slice-Plan `slice-244` (in-progress); welle-91-Plan; `harness/conventions.md` (MR-053, MR-054, MR-073); `AGENTS.md` §3; Beobachtungs-Register-README; Schwester-Originale aus `/tmp/aih-v6.13.0`; keine `DC-*` (Plan: „Rollen- und Command-Doku berührt kein Produkt-Anforderungs-Delta"). Vorherige Findings am selben Modul: keine übergeben (Erst-Review der Adoptions-Area; slice-235-Findings betrafen die Alt-Bestands-Agents, hier nicht relevant).

---

## Findings

### F1 — LOW: Die Ablehnungen der Schwester-Kandidaten sind nur teilweise belegt; die Implementer-Ablehnung trägt keinen benannten Grund

- **kategorie:** LOW
- **quelle:** Maintainability (slice-244 §1 — die Evaluierung ist der Zieldeliverable des Slices)
- **pfad:** `docs/plan/planning/in-progress/slice-244-claude-rollen-commands-evaluieren.md:26-31`
- **befund:** Der Plansatz nennt als Ablehnungs-Grund „d-check führt heute reviewer, verifier als Agents und implement-slice als Skill" — das belegt drei der vier Ablehnungen (reviewer, verifier, implement-slice existieren in eigener Form, geprüft: Dateien stammen aus slice-235/240, vor dem reviewed Commit). Die vierte Ablehnung (Implementer-Agent) ist nicht benannt: d-check hat **keinen** Implementer-Agent, das „existiert-bereits"-Argument trägt für sie nicht — ihr Grund (die Rolle läuft im Hauptlauf nach `AGENTS.md` §6 bzw. über das Command `implement-slice`) ist nur erschließbar, nicht geschrieben. Zwei Schwächen im selben Satz: „implement-slice als **Skill**" verweist in die falsche Form-Klasse (`implement-slice` ist ein Command unter `.claude/commands/`; `.harness/skills/` enthält kein implement-slice), und der Satz trägt die Dopplung „wellenlosen/wellenlosen-Betrieb" — die geforderte Grundlage (Modul 8 und d-checks wellenloser/wellengebundener Betrieb) ist in der Kopplung unlesbar. Folge-Szenario: liest die Welle-Closure (Trigger-Audit) oder ein späterer Lauf den Slice, ist die Implementer-Ablehnung von einer vergessenen Rolle nicht unterscheidbar — genau das Muster, das der adoptierte Validator-Text selbst benennt („ein ausgelassener Schritt und ein begründet übersprungener sehen im Nachhinein gleich aus").
- **verifizierbar:** nein — kein Gate; Lesen des Plans gegen den `.claude/`-Bestand
- **klasse:** `evaluierung-ablehnung-nur-teilweise-belegt`

### F2 — LOW: plan-welle verortet die drei Vorprüfungen im „Plan-Kopf", die Vorlage und die Haus-Form tragen sie in §8

- **kategorie:** LOW
- **quelle:** Maintainability (Adoption-Drift gegen `MR-053`/`MR-054` und die Baseline-Vorlage)
- **pfad:** `.claude/commands/plan-welle.md:29-30`
- **befund:** Der Bullet sagt: „§8 der Slice-Pläne trägt die beiden `d-check:cite`-Spannen (MR-054), der **Plan-Kopf** die drei Vorprüfungen samt Nachtlauf (MR-053)". Beide Hälften desselben Satzes widersprechen sich im Ort: die Baseline-Vorlage (`v6.13.0` · `templates/docs/plan/planning/slice.template.md`, §8) und die Haus-Form (slice-244s eigener Plan) tragen die drei Vorprüfungen in **§8**, nicht im Kopf — und MR-053s Geltungsbereich sagt „im Abschnitt der vorgelagerten Prüfungen". Folge-Szenario: ein Planner, der `/plan-welle` folgt, setzt die Vorprüfungen in den Plan-Kopf und divergiert von der Form, die Vorlage und die MR-054-Hälfte desselben Bullets verlangen; kein Sensor fängt das (Prosa).
- **verifizierbar:** nein — kein Gate (Prosa-Behauptung); Lesen gegen Vorlage und Bestands-Plan
- **klasse:** `adoption-verortet-inhalt-falsch`

### F3 — LOW: Der bei der Beanspruchung versprochene Nachtlauf-Stand ist nicht notiert

- **kategorie:** LOW
- **quelle:** Maintainability (slice-244 §8 gegen den Beanspruchungs-Commit; `MR-053`)
- **pfad:** `docs/plan/planning/in-progress/slice-244-claude-rollen-commands-evaluieren.md:114-116` (Versprechen) · `:89-95` (§7 leer)
- **befund:** §8 sagt: „`make nightly-state` wird bei der Beanspruchung gelesen und der Stand in §7 notiert"; §4 knüpft daran. Die Beanspruchung ist vollzogen (`c7091b7b`), §7 (Closure-Notiz) ist leer — der Stand des Nachtlaufs zum Beanspruchungszeitpunkt ist nicht belegt. Folge-Szenario: wird §7 erst bei der Closure gefüllt, trägt die Notiz den Stand eines späteren Tages; der Zustand „bei der Beanspruchung gelesen" ist dann nicht rekonstruierbar und die Vorprüfung war Form ohne Beleg.
- **verifizierbar:** nein — kein Gate; Lesen des Plans gegen den Beanspruchungs-Commit
- **klasse:** `vorpruefung-versprochen-ohne-vollzug`

---

## Negativbefunde (geprüft, ohne Befund)

- **Rollen-Achse der drei adoptierten Agents:** geprüft — der Schluss-Absatz („Warum dieser Typ existiert") trägt in allen drei Dateien unverändert die sechs kanonischen Rollen-Namen (`planner`, `architect`, `implementer`, `reviewer`, `verifier`, `validator`). Erhalten.
- **Inhaltstreue der Adoption (architect/planner/validator gegen die Schwester-Originale):** geprüft — die Abweichungen sind gezielte Adaptionen, keine Drift: der `ANPASSEN`-Platzhalter ist durch repo-eigene Quellen-Kommentare ersetzt, die Erfassungs-Tie („in der Erfassung", „Rollen-Feld" → „Rollen-Zuordnung") ist konsequent entfernt (die Erfassungsschicht ist Out-of-Scope von welle-91), und planner.md trägt die d-check-Schärfung der Move-Regel („100 %, Link-Tiefen-Korrektur als eigener Commit"). Die gezielte Kürzung in close-welle (Satz zum Carveout-Schluss) ist redundant gedeckt — die Regel steht im adoptierten planner.md.
- **close-welle zählt VIER Trigger-Audit-Klassen** (Carveout · bootstrap-aware Gate · ADR · Hard Rule) **— geprüft, kein Befund:** der Beleg „Adoption slice-241, dokumentiert im Beobachtungs-README" ist ehrlich; das Register-README dokumentiert die vier Klassen seit slice-241 wortgleich als Adoption des Baseline-Regelwerks (`v6.13.0` · `regelwerk/modul-06-roadmap.md` §Closure, Schritt 2), das Schwester-Original zählt drei.
- **Checkbare Infrastruktur-Behauptungen der Adopt-Dateien — geprüft, alle stimmen:** `doc-check` heißt in beiden Commands `doc-check` (der `docs-check`-Tippfehler der Schwester ist gefixt); „`welle-NN` ist keine Traceability-Klasse" steht in beiden Commands (`plan-welle.md:25-26`, `close-welle.md:23-24`); `make archive-wave WELLE=<welle-id>` existiert (Makefile-Zeile 421, Benutzungsform identisch); „`make gates` endet mit `record-gates`" stimmt; „link-policy: always" stimmt (`.d-check.yml`); `welle.template.md`/`welle-results.template.md` existieren in den vendored Templates.
- **Pfad-Tiefen und Anker — geprüft, kein Befund:** `plan-welle.md:29-30` nutzt `../../harness/conventions.md` aus `.claude/commands/` → Repo-Wurzel, korrekt; die Anker `#mr-053`/`#mr-054` existieren (Kurz-Anker im Adaptions-Block).
- **planning-check-Konsistenz (Prüfpunkt 4) — geprüft, kein Befund:** die welle-91-Zeile steht in der Offene-Wellen-Liste (angelegt in Basis-Commit `e0165e8c`), der Ruhe-Marker „Nichts in Arbeit." ist mit der Beanspruchung (`c7091b7b`) aus der Roadmap entfernt — genau der in der Roadmap deklarierte Normalfall (Marker steht, wenn `in-progress/` leer ist; „beides zugleich ist der Normalfall direkt nach der Wellen-Eröffnung"). Der `git mv` des Plans ist rein (Rename-Erkennung intakt).
- **Commit-Traceability der Range — geprüft, kein Befund:** beide Commits tragen gültige Kennungen (`slice-244`/`MR-053` bzw. `slice-244`/`MR-073`); welle-91 wird nur zusätzlich genannt, nicht als Klasse behauptet — konsistent mit der oben verifizierten Lehre.
- **Fünf-Klassen-Kommentare / Zustandsfelder in den Neuzugängen — geprüft, kein Befund:** die Markdown-Dateien tragen keine Review-Chronik und keine Befund-Marker; der Plan führt kein `**Status:**`-Feld (Lifecycle-Hinweis korrekt).
- **DoD-Abhakung — nicht geprüft** (Rolle des Verifiers); die DoD-Haken des Plans sind offen, der Slice ist in-progress.

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 0 |

## Verdikt

**Nicht blockiert.** Keine HIGH-/MEDIUM-Findings; die drei LOW-Findings sind Lesefracht der Evaluierungs-Belege und der Adoptions-Form, alle im laufenden Slice (vor Closure) behebbar. Die Adoption selbst ist inhaltstreu: die fünf Dateien folgen dem Schwester-Original mit sauber gezogenen d-check-Adaptionen, die Rollen-Achse ist erhalten, und jede behauptete d-check-Besonderheit (Gate-Namen, Target-Namen, MR-Anker, Vier-Klassen-Audit) hält der Überprüfung stand. Übergabe an den Verifier für die DoD-Abhakung (insbesondere DoD 1 — „je Rolle (6) und Command (3) eine belegte Entscheidung" — erfordert die Klärung von F1).
