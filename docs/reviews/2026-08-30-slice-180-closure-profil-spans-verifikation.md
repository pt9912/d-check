# Verifikation slice-180 — `spans` am Closure-Bindepunkt

**Gegenstand:** `085459e~1..6835944` (4 Commits), HEAD `6835944`, ungepusht.
**Rolle:** unabhängiger Verifier (Prüfung gegen DoD und Spec, nicht Code-Review).
**Gefahrene Sensors (alle lesend, keiner schreibend):** `make verify-closure-notes` · `make doc-check` · `make gate-consistency` · `make planning-check` · `make adr-check RANGE=…` · `make completeness-check` · `make test` · `make nightly-state` · direkte `docker run`-Läufe des Images gegen Scratch-Proben.
**Nicht gefahren:** `make gates`, `make fullbuild`, `make record-gates` (Auftrag — `gates` zieht `record-gates`).

## 1. DoD-Tabelle

| # | Behauptet (§4) | Gemessen | Verdikt |
|---|---|---|---|
| 1 | `verify-closure-notes` fährt `spans`; Befundzahl über den Bestand unverändert (0 aus `spans`) | Make echot das Rezept `--enable planning --enable structure --enable spans`. Vorher-Rezept über den echten Bestand: **546 Dateien, 0 Befunde, Exit 0**. Nachher-Rezept: **546/0, Exit 0**. Closure-Profil mit **nur** `spans`: 546/0. Dass `spans` nicht bloß schweigt, ist über die Probe (Punkt 2) belegt, nicht über den grünen Lauf | **erfüllt** |
| 2 | Neue Deckung: Probe mit vergessenem Schluss-Fence vorher grün, nachher rot (`fence-unclosed`) | Eigene Probe (§3): vorher **0 Befunde/Exit 0**, nachher **`fence-unclosed`/Exit 1**. Zusätzlich eine **Kontrolle**, die zeigt, dass das Vorher-Grün nicht leer ist | **erfüllt, stärker belegt als zugesagt** |
| 3 | Vier Deklarations-Flächen tragen die dritte Modul-Angabe; `gate-consistency` grün | `AGENTS.md:386` · `harness/README.md:98` · `Makefile:338` (`##`-Hilfetext) · `.d-check.closure.yml:15` + `:19`. `make gate-consistency`: **608/0, Exit 0** | **erfüllt** (Einschränkung: A-3) |
| 4 | ADR begründet Verortung, Deckung, Nicht-Änderung der Lexik; im Index | `docs/plan/adr/0076-spans-am-closure-bindepunkt.md`, Entscheidungen 1–5 decken alle drei; Index-Zeile `docs/plan/adr/README.md:86` | **erfüllt** (Einschränkung: A-4) |
| 5 | ADR-0042 trägt die Neu-Messung als `## Geschichte`; `adr-check` grün | Reiner Anhang einer Tabellenzeile (`0042-…:149`), Kern unberührt. `make adr-check RANGE=085459e~1..6835944`: **608/0, Exit 0** | **erfüllt** |
| 6 | slice-178 §1 nennt die richtige Zähl-Einheit; `BEO-020` ist fortgeschrieben | Zähl-Einheit: korrigiert (`slice-178:53-66`) — aber die neue Zahl ist nicht reproduzierbar (A-1) und der Schluss trägt nicht (A-2). **`docs/plan/planning/observations.md` ist im ganzen Range nicht angefasst** (`git log … -- observations.md` leer) | **teilweise** |
| 7 | `make gates` + `make fullbuild` grün, unabhängiger Review, Verifikation | Auftragsgemäß nicht gefahren. Gefahren: die vier `gates`-Glieder, die dieser Diff berühren kann (`doc-check`, `gate-consistency`, `planning-check`, `test`) — alle grün; kein `.go`- und kein `.github/`-File im Range, also sind `lint`/`arch-check`/`coverage-gate`/`semgrep`/`workflow-pins` unberührt | **offen (planmäßig)** |

## 2. Nachgeprüfte Zahlen

| Zahl | Fundstelle | Behauptet | Gemessen |
|---|---|---|---|
| Markdown-Dateien | slice §1, ADR-0076 §Kontext, ADR-0042 §Geschichte, Commit-Botschaft | **676** | **676 rekonstruiert** = 659 getrackte `.md` bei `9d22a44` + 17 ungetrackte `.harness/cache`-Templates. Enthält die **4 Symlinks** unter `.claude/rules/`, die auf Baseline-Dateien zeigen → **672 verschiedene Dateien**, 4 doppelt gezählt. Heute: 678 Pfade / 674 reguläre Dateien |
| Abweichungen der beiden Fence-Lesarten | ebd. | **0** | **0** — eigene Nachbildung (`FenceToggle`-Toggle vs. `FenceRun`+`FenceCloses`, Öffner in **beiden** über `FenceToggle`, wie `spans.go:89`/`trace_table.go:337` es tun), gefahren über **678 Pfade** (Superset der 676). Kein Zeilen-Unterschied, kein End-Zustands-Unterschied |
| „Die Messung *kann* finden" — 3 konstruierte Fälle | ebd. | 3 gefunden, saubere Datei nicht | reproduziert: Infostring hinter dem Schluss-Fence, `` ``` ``→`~~~`, zu kurzer Schluss-Run — je 2 Divergenz-Meldungen; `clean.md` still |
| **Modellierungs-Gegenprobe** | — | — | Unter der **falschen** Modellierung (strikt öffnet bei jedem Run ≥ 3, ohne `FenceToggle`-Gate) meldet derselbe Baum **1 Datei / 20 Zeilen** (`docs/plan/adr/0042-…md`). Die veröffentlichte **0 ist genau die der richtigen Modellierung** — die Korrektur des Autors war notwendig und ist wirksam |
| `spans`-Befunde über den ganzen Baum | ADR-0076 E3, `harness/README.md:98` | **0** | **0 über 674 Dateien** (`roots: ["."]`, kein `ignore`, Config container-seitig übergelagert). Die Belegläufe des Autors decken 546 bzw. 608 — die Aussage hält also, auf einer **größeren** Menge als ihr Beleg (siehe A-5) |
| Bindepunkt | Commit-Botschaft, DoD 1 | **546 Dateien / 0** | **546/0, Exit 0**. 546 = `docs` (537) + `spec` (3) + Wurzel-`*.md` (6) — das Closure-Profil hat **keinen `scan:`-Block**, es fährt die **Default-Wurzeln** (`rules/scan.go:25`), nicht den Baum |
| `doc-check` / `gate-consistency` | Commit-Botschaft | **608 / 0** | beide **608/0, Exit 0** |
| slice-178: Backticks in slice-061/slice-076 | `slice-178:59` (korrigiert) | **21 bzw. 4** | **nicht reproduzierbar** (A-1). DoD-Abschnitt roh: **25 / 45** — das sind die *alten* Zahlen, und als Abschnitts-Summe **korrekt**. Ohne Fenced-Code: 25/45. Backtick-**Läufe**: 23/37. Zeilen mit Backtick: 8/13. Absatzweise: der DoD-Abschnitt ist in beiden Dateien **genau ein** Absatz → wieder 25/45 |
| „`spans` meldet dort nichts; kein Absatz ist unbalanciert" | `slice-178:63` | — | **für diese beiden Dateien zutreffend** (eigene Nachbildung von `unclosedRuns`/`sticksToText`/`stripInlineCodeByLine`, gegen das Produkt an 4 Proben kalibriert): kein ungeschlossener Lauf, kein verschlucktes Task-Item |
| Nachtlauf-Stand (`MR-053`) | `slice-180:191-193` | beide grün, `2026-08-30T06:08:17Z` / `2026-08-29T10:07:43Z` | `make nightly-state`: **byte-genau dieselben Zeitstempel**, beide `gruen` |
| Register-Stand | `slice-180:173` | höchste Kennung `BEO-024`, `BEO-020` Zähler 3, `BEO-023` Zähler 3 | alle drei bestätigt (`observations.md:11,12,15`) |

## 3. Die diskriminierende Probe

Gebaut aus einem echten `done/`-Slice (`slice-179`) in einem Scratch-Repo mit vier Dateien, das unter **beiden** Rezepten grün startet (die erste Fassung war es nicht — sie meldete `section-missing` und `planning-drift`, also musste eine `welle-*-results.md` und ein `in-progress/`-Slice dazu; genau die Falle, die die Commit-Botschaft beschreibt).

| Probe | Altes Rezept (`planning`+`structure`) | Neues Rezept (`+spans`) |
|---|---|---|
| **Basis** (unverändert) | 4 Dateien, 0 Befunde, **Exit 0** | 4/0, **Exit 0** |
| **P1** — nur der vergessene Schluss-Fence am Dateiende | **0 Befunde, Exit 0** | **`fence-unclosed` :379, Exit 1** |
| **P2** — Fence **plus** ein Platzhalter `(bei Closure)` dahinter | **0 Befunde, Exit 0** — die `section-forbidden`-Zusage wird **still wahr** | **`fence-unclosed`, Exit 1** |
| **P3 — Kontrolle:** derselbe Platzhalter **ohne** Fence | **`section-forbidden`, Exit 1** | ebenso |

P3 ist der Punkt: es belegt, dass das Vorher-Grün in P2 **nicht** daher kommt, dass die Regel nicht gilt — sie gilt, sie feuert, und der Fence hat sie stummgeschaltet. Damit ist die tragende Aussage von ADR-0076 §Kontext („der Preis hat eine zweite Hälfte") **gemessen**, nicht nur behauptet. Die DoD verlangt nur P1; P2+P3 sind stärker.

**Zwei Präzisierungen aus derselben Probe:**
- In P2 meldet das neue Rezept **nur** `fence-unclosed`, nicht den verschluckten Platzhalter. `spans` macht den **Defekt am Text** laut; es stellt die verschluckte Zusage nicht wieder her. ADR-0076 formuliert genau das korrekt („nicht mehr möglich, **ohne dass etwas meldet**") — keine Abweichung, aber der Unterschied gehört benannt.
- Ein vergessener Schluss-Fence ist **nicht durchweg** still gewesen. Platziert **direkt hinter** der Closure-Notiz-Überschrift war das alte Rezept schon rot (`closure-note-thin` + `section-empty`). Still ist er nur, wo er **nach** allen bewachten Abschnitten steht — genau die Stelle, an der die Slice-Vorlage die Closure-Notiz führt. Der Anlass ist damit real und schmaler, als §1 nahelegt.

## 4. Spec-Konformität

**„Berührte Spec-Stellen: —" trägt.** Gemessen:
- Der Diff fasst **keine Datei unter `spec/`** an (`git diff --name-only`).
- `DC-FA-SPAN-001`, `DC-FA-PLAN-001`, `DC-FA-CLI-012`: keine Aussage-Änderung — das Modul, das Profil-Flag und die Closure-Fähigkeit bleiben, wie sie sind; geändert wird, an welchem **zweiten Ort dieses Repo** sie fährt.
- `DC-QA-03`: Anforderung („schreibt nie ins Repo, öffnet außer `external`/`sources` keine Verbindung") und Messmethode („alle Module außer `external`/`sources`/`vcs` aktiv" = der `doc-check`-Lauf, 608/0) unberührt.
- `spec/` nennt weder `verify-closure-notes` noch `.d-check.closure.yml` — die Abwärts-Sperre aus §3.4 ist eingehalten.
- `make completeness-check`: **50 Anforderungen, 0 Waisen, Exit 0**.
- Kein `CHANGELOG`-Eintrag — korrekt, es gibt keine nutzersichtbare Änderung.

**Der Go-Test verliert keine Aussagekraft.** `TestQA03_ClosureProfil_KeineZweiteNetzTuer` prüft die `modules:`-Liste der `.d-check.closure.yml` (weiterhin `[]`) plus `planning.closure.dir ≠ ""`. Beides unverändert; `make test` grün, inkl. `configyaml`. Er war schon vor slice-180 gegenüber dem **Rezept** blind (`structure` stand dort bereits) — die Zahl der außerhalb seiner Reichweite scharfgeschalteten Module geht 2 → 3. Praktisch gedeckt bleibt die Netz-Zusage durch `--network none` in `DCHECK_RUN` (`Makefile:122`), nicht durch diesen Test. Siehe A-7.

## 5. Abweichungen Zusage ↔ Zustand

**A-1 — `docs/plan/planning/open/slice-178-offene-tasks-roh.md:59` (MEDIUM).**
*„Nachgemessen: 21 bzw. 4 Backticks"* ist unter keiner geprüften Lesart reproduzierbar: DoD-Abschnitt roh 25/45, ohne Fenced-Code 25/45, Läufe 23/37, Zeilen-mit-Backtick 8/13, absatzweise (der Abschnitt **ist** ein Absatz) 25/45. Schärfer noch: die **alten** Zahlen 25/45 waren als Abschnitts-Summe **richtig** — falsch war allein die *Relevanz* der Einheit. Die Korrektur ersetzt damit eine korrekte Zahl durch eine unbelegte und begeht dabei die Klasse, die sie korrigiert (`BEO-020`: gemessen wird die eine Menge, ausgesagt über die andere).

**A-2 — `slice-178:63-65` (MEDIUM).**
Die Kette *„`spans` meldet nichts → kein Absatz ist unbalanciert → **Die Exposition dieses Repos ist damit heute null**"* trägt nicht. Gemessen mit einer gegen das Produkt kalibrierten Nachbildung über 678 Pfade:
- **7 unbalancierte Absätze** repo-weit (u. a. `spec/lastenheft.md:1046`, `docs/user/benutzerhandbuch.md:2204`, `docs/plan/adr/0042-…:120`) — `spans` schweigt dort **bewusst** (`sticksToText`), Schweigen belegt also keine Balance.
- **6 Stellen in 5 Dateien**, an denen ein `- [ ]` durch eine **wohlgeformte** Spanne aus dem bereinigten Text verschwindet — darunter zwei `done/`-Slices (`slice-096:176`, `slice-179:355`) und slice-178 selbst (`:121`, `:218`). Das ist genau der Mechanismus, um den es slice-178 geht.
- Das **Fazit** hält trotzdem: alle 6 sind bewusste Prosa über den Marker, kein echter unquittierter Haken ist verdeckt. Die Aussage stimmt — der angegebene Grund belegt sie nicht.

**A-3 — `AGENTS.md:386`, `harness/README.md:98`, ADR-0076 §Kontext, `slice-180:41-46` (MEDIUM).**
Alle Flächen führen die neue Deckung als **`fence-unclosed`, `span-unclosed`**. Das Modul kennt einen **dritten** Grund-Code, `span-nested-link` (`model/finding.go:21`). Gemessen: er macht den Bindepunkt rot (`docs/plan/planning/done/slice-179-…:379 ](ziel.md)]( span-nested-link`, Exit 1). Das Risiko in §5 („ein drittes Modul ist ein drittes, das rot werden kann") ist damit um eine unbenannte Klasse breiter, als jede Deklarations-Fläche sagt.

**A-4 — `docs/plan/adr/0076-…md` §Verglichene Alternativen (LOW).**
*„`non-empty: true` an die bestehenden Closure-Regeln hängen. Macht denselben Fall laut (gemessen: `section-empty`)."* Gemessen, beide Platzierungen:
- Fence am Dateiende (der Fall, auf dem der Entscheid steht): mit `non-empty` an der `forbid-pattern`-Regel **0 Befunde, Exit 0** — es wird *nicht* laut.
- Fence direkt hinter der H1-Titelzeile: `section-empty` **erscheint** — aber dort ist das heutige Profil **ohnehin schon rot** (4 weitere Befunde ohne jede Änderung).
Die verworfene Alternative ist also **schwächer**, als die ADR ihr zugesteht; der Verwerfungs-Entscheid wird dadurch richtiger, nicht falscher. Die ADR ist `Accepted` → Korrektur nur als `## Geschichte`-Anhang oder Folge-ADR.

**A-5 — ADR-0076 E3, `harness/README.md:98` (LOW).**
*„`spans` meldet über den ganzen Baum 0 Befunde."* Die Aussage ist **wahr** (nachgemessen: 674/0), aber keiner der genannten Belegläufe deckt den ganzen Baum: der Bindepunkt sieht **546** Dateien (Default-Wurzeln `docs`/`spec`/Wurzel-`*.md` — das Closure-Profil hat keinen `scan:`-Block), `doc-check` **608** (Baum minus `.harness/baseline`/`cache`). 66 `.md`-Dateien liegen in keinem der beiden. Gleiche Klasse wie A-1/A-2, nur mit richtigem Ergebnis.

**A-6 — DoD 6, zweite Hälfte (BLOCKIEREND für die Closure).**
`docs/plan/planning/observations.md` ist im Range nicht angefasst; `BEO-020` steht unverändert bei Zähler 3 mit den Belegen slice-142/160/179. Kanon (`modul-06` §Beobachtungs-Register): der Beleg wird **vor** dem `git mv` nach `done/` geschrieben. Offene Closure-Pflicht, kein jetziger Defekt — der Haken ist korrekt ungesetzt.

**A-7 — `internal/adapter/driven/configyaml/gate_consistency_test.go:75-76` (LOW).**
Der Kommentar sagt, das Closure-Profil sei „ein fokussiertes Profil, das **nur `planning`** per Kommandozeile dazuschaltet" — es sind drei. Der Kommentar war schon vor slice-180 veraltet; der Slice vergrößert die Abweichung, ohne sie anzufassen. Dazu ADR-0076 E2: „ein Go-Test hält bereits, dass **das Profil** keine Netz-/Range-Tür öffnet" — er hält das für die `modules:`-Liste, also für die Stelle, die der Entscheid **leer lässt**, nicht für das Rezept, in das er die Module legt. §Fitness Function benennt die Rezept-Lücke für `spans`, nicht für ein Netz-Modul.

**A-8 — Beobachtung, keine Abweichung (INFO).**
`slice-180` führt sich als **wellenlos**, während `welle-86` („der Closure-Übergang trägt seine Vorbedingungen") offen ist und `slice-178` — den dieser Slice entsperrt — seinen Anlass ausdrücklich dort verortet. Der Baseline-Kanon lässt wellenlose Slices neben offenen Wellen zu; die Zuordnung ist Planner-Urteil und nicht meine Rolle. Notiert, weil die Sub-Area dieselbe ist.

## 6. Lifecycle (E)

- `6b485a6` ist ein **reiner Move**: `R100`, Slice-Datei byte-identisch (`git diff` zwischen beiden Blobs: 0 Zeilen).
- Der **gekoppelte Roadmap-Flip** ist im selben Commit: `Nichts in Arbeit.` verlässt §Offene Wellen, der `welle-86`-Zeiger bleibt (`MR-013`). `make planning-check`: **608/0, Exit 0**.
- Die zwei `d-check:cite`-Direktiven treffen die **vorschreibenden** Zeilen exakt: `modul-05:213-214` = *„**Sub-Area-Wahl prüfen.** … Schwelle ≥ 2"*, `modul-05:219` = *„**Offene Beobachtungen sichten.**"* (`MR-054`); `citations` verifiziert sie wortgleich im grünen `doc-check`. Der dritte Block trägt korrekt keine.
- `make adr-check` über den ganzen Range: grün — ADR-0042 ist reiner `## Geschichte`-Anhang.

## 7. Was besser ist als zugesagt (F)

1. Die **0-Abweichungen-Messung ist modellierungs-kritisch und hält**: unter der falschen Öffner-Modellierung ergäbe derselbe Baum 1 divergente Datei / 20 Zeilen. Die veröffentlichte Null ist die der richtigen Modellierung — die Korrektur war nötig und ist wirksam.
2. Sie hält auf **678 Pfaden** statt 676.
3. `spans`-Null hält auf **674 Dateien** statt auf den 546/608 der Belegläufe.
4. Die neue Deckung ist mit **Kontrolle** (P3) belegbar, nicht nur vorher/nachher — das ist genau die Antwort auf `BEO-023`, die §7 verspricht, und sie ist stärker, als die DoD verlangt.
5. Der Lifecycle-Move ist formal sauber (R100 + gekoppelter Flip), die Nachtlauf-Zeitstempel reproduzieren byte-genau, die Zitat-Anker sitzen auf der Zeile.

## 8. Repo unverändert

`git status --porcelain` leer, `git diff --stat` leer, `HEAD = 6835944` unverändert. `.harness/state/gates-passed.diffsha` trägt weiterhin `mtime 2026-08-30 10:32:14` (der Lauf des Autors) — `record-gates` lief nicht. Alle Proben liefen unter `…/scratchpad/verify/` bzw. als container-seitige Bind-Mount-Überlagerung; das Scratch-Verzeichnis ist entfernt.

## 9. Empfehlung

**Vor der Closure zu schließen:**
- **A-6** — `BEO-020` fortschreiben (Kanon: vor dem `git mv`). Dabei ist zu entscheiden, ob slice-180 eine vierte Instanz ist oder nur eine Korrektur der dritten; ebenso `BEO-023` (die Probe P1/P3 ist die Antwort, die §7 ankündigt).
- **A-1** — die Zahl „21 bzw. 4" belegen oder ersetzen. Belastbar wäre: *„25 bzw. 45 Backticks im DoD-Abschnitt — als Abschnitts-Summe richtig, für den absatzweise paarenden Mechanismus aber ohne Aussage; beide Abschnitte sind je **ein** Absatz, und in keinem kippt die Paarung."*
- **A-2** — den Schluss an das binden, was ihn trägt. Gemessen liegt vor: 0 `spans`-Befunde, 7 unbalancierte Absätze ohne Befund, 6 durch wohlgeformte Spannen verschluckte Task-Items, alle sechs bewusste Prosa.
- **A-3** — `span-nested-link` in die vier Deklarations-Flächen und in §5 aufnehmen (eine Zeile je Fläche); es macht den Bindepunkt rot und steht nirgends.

**Danach möglich (kein Closure-Blocker):**
- **A-4** und **A-5** als `## Geschichte`-Anhang an ADR-0076 bzw. als Präzisierung in `harness/README.md` („über 674 Dateien gemessen"; „der Bindepunkt-Lauf sieht die Default-Wurzeln, nicht den Baum").
- **A-7** — Kommentar in `gate_consistency_test.go` nachziehen; optional die Feststellung, dass die Netz-Zusage am Bindepunkt von `--network none` getragen wird, nicht von diesem Test.

**Fachlich trägt der Slice.** Die Kern-Zusage — der Closure-Bindepunkt sieht eine Defekt-Klasse, die er vorher übersah, ohne Bestands-Rauschen — ist von mir unabhängig reproduziert und mit einer Kontrolle abgesichert. Die Abweichungen betreffen durchweg **Beleg-Präzision**, nicht die Mechanik; A-3 ist die einzige mit funktionaler Kante.
