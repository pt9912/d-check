# Verifikation slice-253 — Baseline-Pin-Hebung `v6.13.0` → `v6.17.0` (DoD)

- **Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?"; geprüft gegen §2 DoD und §1 Abgrenzung des Slice-Plans, nicht gegen Review-Entscheidungen.
- **Gegenstand:** `0c0d46dc..93288754` — 5 Commits: `f8382b60` (reiner Move MR-073 → `conventions/done/`), `c7ace79b` (MR-073 Link-Tiefen, ±5), `957aeedc` (Baum-Swap, MR-074, Retargets), `2201353c` (R1-Report), `93288754` (R1-Einarbeitung in MR-074).
- **Sensor-Evidence:** jede Zahl unten selbst gemessen: Vorher-Baum per `git archive 0c0d46dc .harness/baseline/v6.13.0`, `diff -rq -I '<!-- Quelle:'`, `git grep` gegen `0c0d46dc`/`HEAD`, `make gates`, `make baseline-freshness`, drei Bruch-Proben. Aus dem Implementer-Bericht ist nichts übernommen.
- **Modell-ID:** claude-opus-5-5 (Claude Agent SDK)
- **Datum:** 2026-10-07

---

## DoD-Prüfung je Punkt (§2)

### 1. Baum auf `v6.17.0`, `SHA256SUMS` verifiziert, `v6.13.0` entfernt, §Baseline-Pin, Vorgänger nach `done/`, neuer MR im Index: **ERFÜLLT**

- `ls .harness/baseline/` → nur `v6.17.0`. `make baseline-verify` (im `gates`-Lauf): `fetch-baseline-cache: verify ok (54 Dateien, vollständig)`.
- Bytes gegen das Release: `make baseline-freshness` → `check-latest OK (Currency) — Pin v6.17.0 ist der neueste Release-Tag` und `check-latest OK (Content) — … upstream unverändert (Bytes == vendored SHA256SUMS)`, Exit 0.
- **Bruch-Probe `baseline-verify`:** ein angehängtes Byte in `regelwerk/modul-15-observability.md` → `sha256sum: WARNUNG: 1 berechnete Prüfsumme passte NICHT`, `make: *** [Makefile:209: baseline-verify] Fehler 1`. Wird aus dem richtigen Grund rot. Danach zurückgesetzt.
- `harness/conventions.md:27-29` Stand `v6.17.0`, gepinnt mit `MR-074`; `:41,46,47,63` retargetet; `:136` MR-074 aktiv; `:179` MR-073 unter §Aufgelöste Adaptionen mit beiden Alt-Ankern und Nachfolger MR-074.
- `f8382b60` ist laut `git diff -M --summary` ein reiner Rename; `c7ace79b` hat 5 Zeilen +/- nur an `done/MR-073`. Damit ist `AGENTS.md` §3.3 Fall 1 eingehalten.

### 2. Lebende pin-gebundene Referenzen retargetet, Cite-Spannen neu geankert, Zitat-Delta und Frozen-Liste gemessen: **ERFÜLLT** (ein LOW zum Bewegungs-Inventar, V-1)

- **Delta-Zählung (MR-074):** 55/55 Dateien, gleiche Dateimenge, keine neu, keine entfallen; 22 mit Inhaltsdelta (8 Regelwerk, 14 Templates) plus `SHA256SUMS`. Diff-Zeilen: `grundlagen-harness-dateien.md` 65, `grundlagen-referenz-richtung.md`/`spezifikation.template.md`/`harness/README.template.md` je 17, `modul-13` 15. Alles wie behauptet.
- **Spiegel-Messung:** `git grep -o 'v6\.13\.0' 0c0d46dc -- . ':!.harness'` → **73 Dateien / 184 Vorkommen**. Ohne den Slice-Plan sind es **72 / 180**. `.harness/skills/` → `reviewer.md` 6 + `closure-note-reviewer.md` 2 = **8**. Alles wie behauptet, der Plan-Ausschluss ist jetzt deklariert.
- **Nur Tag-Swaps:** Im Range-Diff `0c0d46dc..957aeedc` außerhalb des Baums, ohne MR-073/074, wurde jede Minus-Zeile mit `v6.13.0→v6.17.0` normalisiert. Sie ist dann mit ihrer Plus-Zeile identisch. Ausnahmen: genau die Cite-Nummern und die Index-Umstellung in `conventions.md`. Nichts wurde inhaltlich adoptiert.
- **Verbliebene `v6.13.0`-Nennungen in lebenden Dateien (HEAD):** `.d-check.yml:369` (Tombstone der Vorstufe), `CHANGELOG.md:103` (Historie), `observations/README.md:21` (Herkunft der Trigger-Audit-Adoption), `tools/harness/selbstpruefung.sh:3` (Stand von ai-harness-init, anderer Gegenstand), sowie Plan und MR-074 selbst als Gegenstand. Jede ist eine Vergangenheits- oder Gegenstands-Aussage, kein Pin-Verweis. Alles Übrige sind Frozen-Klassen.
- **Frozen byte-stabil:** `git diff --stat 0c0d46dc -- docs/plan/planning/done docs/reviews docs/plan/adr docs/plan/cr CHANGELOG.md .d-check.yml spec/lastenheft.md` zeigt nur zwei Dateien. Der R1-Report ist neu. `conventions/done/MR-073` ist der Zug samt Tiefen-Commit. Sonst ist kein eingefrorenes Byte geändert.
- **Symlinks (`MR-055`):** 11 unter `.claude/rules/`. 8 davon zeigen auf `.harness/baseline/v6.17.0/regelwerk/…`, alle lösen auf (`test -e`). Die 3 übrigen sind keine Pin-Träger.
- **`d-check:cite` (`MR-051`):** Bei `0c0d46dc` gibt es 14 lebende Direktiven außerhalb der `citations.scope.ignore`-Pfade. 4 haben neue Zeilennummern: MR-031 194→196, MR-035 49→50, `reviewer.md` 82→87, `closure-note-reviewer.md` 83→84. 10 sind tag-only. Alles wie behauptet. **Bruch-Probe:** MR-031 auf `:194-194` zurückgesetzt ergibt `citation-mismatch … Zitat-Fäule`, `doc-check` Fehler 1. Danach zurückgesetzt.
- **`ignore-refs` wächst nicht (`MR-069`):** `doc-check` ist auf HEAD grün. **Bruch-Probe:** ein Markdown-Link auf `.harness/baseline/v6.13.0/regelwerk/README.md` in einer lebenden Datei ergibt `target-missing`, `doc-check` Fehler 1. Der Sensor würde einen toten Link auf den entfernten Baum also sehen, und die 0 Befunde tragen die Behauptung. Danach zurückgesetzt.
- **Anker:** Lebende Links in den `v6.17.0`-Baum haben 15 eindeutige `#…`-Anker. Alle 15 lösen gegen die Heading-Slugs bzw. `id=` der Zieldateien auf. Abgleich per Shell, das Verfahren ist an einem Kontroll-Anker (`bogus-anker` → miss) geprüft. Das ist eine eigene Messung, weil die Zieldateien außerhalb der Scan-Wurzel liegen (`AGENTS.md` §3.8).
- **Zitat-Delta (`MR-039`):** `AGENTS.md:286` „Halluzinierte Gates sind die häufigste Form von Harness-Lüge" steht in `v6.17.0` `modul-13:99`. Das MR-056-Zitat „DoD-Häkchen und Closure-Notiz sind die Bedingung dafür, dass die Datei überhaupt nach `done/` darf" steht in `v6.17.0` `modul-05`. Der Delta im elidierten `.d-check.yml`-Zitat (`:932-934` gegen `AGENTS.template.md:206-207`) ist bestätigt und in MR-074 vermerkt. Zusätzlich gegengeprüft: Kein Präfix (40 Zeichen) einer im Delta entfernten Zeile kommt in einer lebenden Datei vor. Grenze: Das ist ein Präfix-Abgleich, er ersetzt keinen vollständigen Zitat-Index.

### 3. `make gates` grün, `make baseline-freshness` meldet den Pin aktuell: **ERFÜLLT** (beides selbst gefahren auf HEAD `93288754`, sauberer Baum)

- `make gates` → `EXIT=0`, Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`. Einzelbelege: `verify ok (54 Dateien, vollständig)`; `d-check: 952 Datei(en) geprüft, 0 Befund(e)` (doc-check, workflows, targets, planning); `golangci-lint … 0 issues.`; Tests ohne `FAIL`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; semgrep `Ran 55 rules on 65 files: 0 findings.`
- `make baseline-freshness` → Currency OK, Content OK (siehe Punkt 1).

### 4. Review durchgeführt, Report liegt vor: **ERFÜLLT**

- `docs/reviews/2026-10-07-slice-253-baseline-v6170-bump-r1.md` (`2201353c`), eigener Kontext.
- **R1-Einarbeitung (`93288754`) gegen den Report:**
  - **F-1** eingearbeitet. MR-074:28-35 nennt `a-check.mk` und den Folge-Slice-Kandidaten. Gemessen: `Makefile:33 include a-check.mk`; einziges Target `a-check:` (`a-check.mk:54`); `make a-check` fehlt in `harness/README.md`, dort steht nur `a-check-digest`; `harness/mk/` existiert nicht. Die neue Aussage stimmt.
  - **F-2** eingearbeitet, aber an zwei Stellen falsch zugeordnet, siehe V-1. Die Dateimenge der sechs Bewegungen deckt alle 22 Delta-Dateien ab.
  - **F-3** eingearbeitet (Zitat-Delta-Absatz).
  - **F-4** eingearbeitet. Kein `Status`-Feld mehr, `Löst auf` und `Ausgelöst durch Baseline-Stand` vorhanden. Das entspricht der Vorlage `v6.17.0` `templates/harness/conventions/MR-NNN-titel.template.md`.
  - **F-5** in Bewegung 2 benannt (Kopplung an MR-0098).
  - **F-6** brauchte keine Aktion.
  - **F-7** eingearbeitet (72/180 bzw. 73/184).

### 5. Closure-Notiz, Register, Risiko-Ausgänge, Paarungen: **OFFEN, wie erwartet**, kein Befund

§6 und §7 stehen auf `(offen)` bzw. `—`. Das ist der Closure-Schritt nach der Verifikation. `make verify-closure-notes` ist nicht anwendbar, solange der Slice nicht in `done/` liegt.

---

## Abgrenzung (§1): gehalten

Laut normalisiertem Range-Diff wurde keine inhaltliche Bewegung übernommen: kein `harness/mk/`, kein `make a-check` im Index, kein Spec-§7. `spec/` hat nur je eine Rolle-Zeile mit Tag-Swap. `tools/harness/selbstpruefung.sh` ist unberührt. Kein Dependabot-Bezug, kein Release.

---

## Verifier-Befunde

### V-1: LOW

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §5 Regel 15 (nicht mehr behaupten, als gemessen); Einarbeitung von R1 F-2
- **Pfad:** `harness/conventions/MR-074-baseline-v6170.md:32` · „`templates/harness/sensors/gate.template.md`, `templates/Makefile` und" sowie `:58` · „`templates/.harness/skills/closure-note-reviewer.template.md`. Dieses Repo"
- **Befund:** Die Liste heißt jetzt „vollständig nach Datei". Ihre Dateimenge ist vollständig (22/22), die Zuordnung zu den Bewegungen stimmt aber an zwei Stellen nicht. Gemessen mit `diff -I '<!-- Quelle:'` gegen den Vorher-Baum:
  - (a) `gate.template.md` steht unter Bewegung 1 (Werkzeug-Teil). Sein Delta ist aber ganz Bewegung 2: Spec-Kennung, und was das Werkzeug prüft, steht „in der Spezifikation". Dazu fehlen unter Bewegung 2 die Spec-Anteile von `harness/README.template.md` (`:92`, `:132-133`), von `grundlagen-harness-dateien.md` (`:344-348`) und die `sensors`-Zeile in `grundlagen-begriffe.md:38`.
  - (b) `closure-note-reviewer.template.md` steht unter Bewegung 4 (Register-Kennung). Sein einziges Delta ist aber `pfad` als wörtliches Kurzzitat, also Bewegung 3. Die tatsächlichen `BEO-<NNN>`-Nachzüge in `review-report.template.md:93` und `modul-10-review-harness.md:140` fehlen unter 4.
  - (c) Bewegung 3 nennt die neue Regel „Kein Stil-Polizist" nicht, auch nicht LOW nur mit Konventions-Anker (`modul-10:77-81`, `reviewer.template.md:66-68,85-86`).
- **Failure-Szenario:** Ein Adoptions-Slice für Bewegung 3 (Reviewer-Skill) stützt sich auf MR-074 als Inventar, denn der Vorher-Baum liegt nur noch in der Historie. Er übernimmt dann die Failure-Szenario-Regel ohne die Stil-Polizist-/LOW-Anker-Regel. Und er sucht die `pfad`-Änderung des Closure-Note-Skills unter der falschen Bewegung.
- **Verifizierbar:** ja. `git archive 0c0d46dc .harness/baseline/v6.13.0 | tar -x -C <tmp>`, dann `diff -I '<!-- Quelle:' <tmp>/…/<datei> .harness/baseline/v6.17.0/<datei>` für die vier genannten Dateien.
- **Klasse:** record-claim-vs-diff (Restbestand der R1-F-2-Einarbeitung)

Sonst keine Befunde. Negativ geprüft ohne Befund: Delta-Zahlen, Spiegel-Zahlen, Skills, Cite 14/4/10, 8 Symlinks, `ignore-refs`, Frozen-Byte-Stabilität, Anker, Zitate, Index-Kette, MR-Feldmenge.

---

## Verdict

**DoD 1–4 erfüllt, DoD 5 erwartungsgemäß offen.** Gates und Freshness wurden selbst gefahren und sind grün. Alle Messzahlen in MR-074 sind nachgezählt und stimmen. Drei Bruch-Proben (`baseline-verify`, `doc-check`/links, `doc-check`/citations) werden aus dem richtigen Grund rot. Ein LOW (V-1) betrifft die Bewegungs-Zuordnung in MR-074 und blockiert die Closure nicht: annehmen oder begründen.
