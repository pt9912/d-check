# Review-Report R2 — slice-206 (Modul `mentions`, Korrekturen nach R1)

- **Review-Art:** Code — geprüft werden die **Korrekturen** gegen den R1-Report, gegen Anforderung, Entscheid und Hard Rules; Maßstab sind [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) (zehn Akzeptanzkriterien, Lastenheft jetzt 0.86.3), [ADR-0084](../plan/adr/0084-mentions-eigenes-modul.md), [`spec/spezifikation.md`](../../spec/spezifikation.md) §`DC-FA-MENT-001.a`/§2/§4, [`AGENTS.md`](../../AGENTS.md) §3.1/§3.4/§3.7/§3.8/§4/§5 und der Baseline-Kanon (`modul-11-verification.md` §Fitness Function ohne Standard-Tool, `modul-13-quality-gates.md` §Hard Rule)
- **Gegenstand:** `e5cf62d` (`fix(mentions)`, drei HIGH + sechs MEDIUM, 13 Dateien, +672/−71) und `3a31c8f` (`fix(mentions)`, sechs LOW, 6 Dateien, +41/−38). Range `ac12993..3a31c8f` — der während des Laufs darüber entstandene Commit `af29684` (Planung für den Baseline-Sprung) gehört nicht zum Gegenstand
- **Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) @ 1.13.0 · **Modell-ID:** claude-opus-5[1m] · **Datum:** 2026-09-06
- **Eingangs-Kontext:** [R1-Report](2026-09-06-slice-206-mentions-modul-review.md) (3 HIGH / 8 MEDIUM / 6 LOW / 4 INFO, blockierend), Slice-Plan [`slice-206`](../plan/planning/in-progress/slice-206-mentions-modul.md), `internal/hexagon/core/rules/mentions.go` + `mentions_test.go`, `internal/adapter/driven/report/report.go` + `report_test.go`, `internal/hexagon/core/rules/{run.go,workflows.go,paths.go,scan.go}`, `.d-check.yml`, `Makefile`, `AGENTS.md`, `harness/README.md`, `harness/sensors/mention-coverage.md`, `spec/lastenheft.md`, `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`
- **Vorherige Findings am gleichen Gegenstand:** R1 zu slice-206 (dieser Lauf), davor [R1](2026-09-06-slice-205-mentions-anforderung-review.md) und [R2](2026-09-06-slice-205-mentions-anforderung-review-r2.md) zu slice-205. Dominante Klassen: `wortlaut-behauptet-pruefung-die-fehlt`, `selbstauskunft-zahl-reproduziert-nicht`, `semantic-change-body-only-edges-stale` — **alle drei erscheinen in diesem Lauf wieder.**

**Gefahrene Läufe.**
`make gates` → Exit 0; `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `d-check: 677 Datei(en) geprüft, 0 Befund(e)`; `coverage-gate: OK — Coverage 94.60% erfüllt Schwelle 93%`; `Ran 55 rules on 63 files: 0 findings.` — Wort für Wort wie in der Botschaft.
`make mention-coverage` → Exit 0; `d-check: mentions: 84 von 84 Artefakt(en) erwähnt, über 1 Dokument(e)`, `0 Befund(e)`.

**Acht eigene Bruch-Tests am Produktivcode**, je mit Wiederherstellung. Sie beantworten die Frage, ob die neuen Tests fallen **können**:

| # | Eingriff | Erwartung | Ergebnis |
|---|---|---|---|
| A | ADR-0084-Zeile aus `docs/plan/adr/README.md` entfernt (`grep -c` danach 0) | Befund | **Exit 1**, `83 von 84`, `artifact-unmentioned` auf dem Artefakt-Pfad — H-1 geschlossen |
| B | `sort.Strings(out)` aus `mentionsFilter` entfernt | `make test` rot | **rot**: `TestMentionsSortierungTrenntDFSVonLexikografisch` — M-3 geschlossen |
| C | fail-closed bei unlesbarem Verzeichnis auf `return nil` zurückgedreht | rot | **rot**: `TestMentionsUnlesbaresVerzeichnisFailClosed` (prüft auch den Meldungstext) — M-8 geschlossen |
| D | **`dirIgnored` UND `ignored` vollständig aus `mentionsWalk` entfernt** | rot | **grün** — kein einziger Test fällt. Siehe N-3 |
| E | `mentionsOccurs(...)` → `strings.Contains(corpus, needle)` | rot | **rot**: `TestMentionsBasenameKeineTeilzeichenkette`, `TestMentionsPathKeinPraefixTreffer`, `TestMentionsRechteGrenze` — H-2 geschlossen |
| F | `matchGlob` → `path.Match` in `mentionsMatchAny` | rot | **rot**: `TestMentionsGlobDoppelsternUeberMehrereSegmente` — H-3 geschlossen |
| G | `omitempty` am `Notes`-Tag entfernt **und** `oneLine(n)` umgangen | rot | **rot**: `TestJSONYAML_OhneNoteKeinNotesSchluessel`, `TestText_NoteBleibtEineZeile` |
| H | `Notes` auf `json:"-"`/`yaml:"-"` **und** die Notes-Schleife in `Text` totgelegt | rot | **rot**: `TestText_NotesAufStderrVorDerZaehlzeile`, `TestJSONYAML_MitNoteImSummary` |

**Elf Mess-Läufe gegen eigene Fixtures und Fremd-Bestände** (Container gegen einen eigenen Mount, `--network none`, read-only) — die Zahlen stehen bei den Findings. Alle Kalibrierungs-Zahlen der Botschaft reproduzieren; die Grenz-Prüfung tut es nicht in allen Formen.

**Arbeitsbaum.** Am Ende ist die einzige Änderung dieser Report. Zwei Dateien (`docs/plan/planning/open/slice-207-…md`, `…/slice-208-…md`) sind während des Laufs von außerhalb dieses Reviews entstanden und inzwischen als `af29684` committet; sie wurden von diesem Lauf nicht angefasst.

---

## R1-Finding → Status

| R1 | Kurzfassung | Status | Beleg |
|---|---|---|---|
| **H-1** | Union-Korpus deckt zwei getrennt gemeinte Invarianten | **behoben** | `.d-check.yml` trägt ein Paar; Bruch-Test A ⇒ Exit 1. Die Sensors-Invariante ist in Sensor-Datei und Anforderung ausdrücklich als **nicht gehalten** benannt (Rest: N-9) |
| **H-2** | Teilzeichenketten-Kollision `test.md` ⊂ `image-test.md` | **teilweise** | Der gemeldete Fall ist zu (Fixture: genau `s/test.md` wird gemeldet); Bruch-Test E belegt die drei Regressionstests. **Aber:** die Korrektur erzeugt vier gemessene Falsch-Alarme (N-1) und lässt die Klasse bei Nicht-ASCII-Nachbarn offen (N-2) |
| **H-3** | `**` trug nicht — `path.Match` statt `scan.ignore`-Semantik | **behoben** | `tools/**/*.sh` ⇒ **11** Mitglieder, 11 Funde (war 6); Bruch-Test F. Am ausgelieferten Glob `docs/plan/adr/[0-9]*.md` ist `matchGlob` ≡ `path.Match` (84 unverändert) — keine Nebenwirkung |
| **M-1** | Kalibrierungs-Zahlen reproduzieren nicht | **behoben** | Neu gemessen und reproduziert: `path` ⇒ **84 von 84** ADRs gegen den Index, **24 von 24** Sensor-Dateien (nicht 23); „21 der 84" ⇒ 21 gemessen. „ausnahmslos falsch" trägt jetzt, weil je Paar gemessen wird |
| **M-2** | DoD (3) verlangt Fremd-Bestand, geliefert war der eigene | **behoben, Begründung falsch** | Am Fremd-Bestand nachgefahren und exakt reproduziert (5 Mitglieder, 4 Funde, beide Formen identisch). Die **Begründung** der Historie-Zeile ist unwahr (N-5) |
| **M-3** | Determinismus-Test misst die Sortierung nicht | **behoben** | Neue Fixture trennt DFS- von lexikografischer Ordnung; Bruch-Test B |
| **M-4** | `Notes`-Vertrag ohne Test | **behoben** | Vier Reporter-Tests, alle können fallen (Bruch-Tests G/H) — darunter die Rückwärtskompatibilität am `omitempty` |
| **M-5** | `scan.ignore`/`scan.roots` ungenutzt und unbenannt | **teilweise** | Das **Produkt** honoriert `scan.ignore` (gemessen) und die `scan.roots`-Wahl ist an fünf Stellen benannt. Der **Test dazu kann nicht fallen** (N-3) |
| **M-6** | §2-Schema-Tabelle ohne die drei Schlüssel | **behoben** | `mentions.artifacts`/`.documents`/`.match` stehen in §2; Typ, Default und Exit-2-Bedingungen decken sich mit `applyMentions` |
| **M-7** | `.d-check.yml` nennt `make mentions-check` | **behoben** | Nennt jetzt `make mention-coverage`; `grep -c "^mention-coverage:" Makefile` ⇒ 1 |
| **M-8** | Unlesbares Verzeichnis verkleinert die Soll-Menge still | **behoben** | Exit 2 mit Ursachen-Nennung, gemessen; Bruch-Test C. Preis: N-10 |
| **L-1** | Plan zählt neun Akzeptanzkriterien, es sind zehn | **behoben** | Plan sagt „zehn"; nachgezählt ⇒ 10 |
| **L-2** | Make-Eintrag trägt den Kommentar des Nachbarn | **behoben** | Beide Targets tragen jetzt einen eigenen Block, im je eigenen Wortschatz |
| **L-3** | Gate-Taxonomie führt `mention-coverage` nicht | **behoben** | Zeile *Meta-/Governance-Gates* führt es samt Bindepunkt „Erwähnungs-Deckung" |
| **L-4** | Baum wird zweimal gelaufen | **behoben** | Ein Durchlauf, zwei Filter; Dedup, Sortierung und Glob-Reihenfolge verhalten sich identisch zur Zwei-Pass-Fassung (Code-Vergleich gegen `ac12993`) |
| **L-5** | `matchAnyGlob`-Kommentar beschreibt weder Konsument noch Wahrheit | **behoben** | Nennt seine zwei Konsumenten (`grep` bestätigt: `reviews.go:145`, `workflows.go:249`) und die Grenze. Folge: N-8 |
| **L-6** | Keine §7-Historie-Zeile für §`DC-FA-MENT-001.a` | **behoben** | Zeile 2026-09-06 in der Spezifikations-Historie |

**Bilanz: 13 von 17 wirklich behoben**, **2 teilweise** (H-2, M-5), **1 behoben mit unwahrer Begründung** (M-2), **1 behoben mit neuem Preis** (M-8 ⇒ N-10). **Nicht behoben: keines.** Kein Finding ist *falsch* behoben im Sinne von „an der falschen Stelle repariert"; der schwerere Fall ist H-2, dessen Korrektur zwei neue Befunde erzeugt.

---

## Akzeptanzkriterium → misst ein Test es wirklich?

Gegenüber R1 („3 vollständig, 4 teilweise, 3 gar nicht") verschieben die elf neuen Tests die Bilanz auf **3 vollständig, 7 teilweise, 0 gar nicht**. Alle Tests liegen weiterhin auf Kern- bzw. Reporter-Ebene; ein CLI-/Akzeptanztest für `mentions` existiert nach wie vor **nicht**, deshalb bleibt jedes Kriterium mit zugesagtem Exit-Code „teilweise".

| # | Kriterium | Test | Misst es wirklich? |
|---|---|---|---|
| 1 | Happy Path — Exit 0, kein Befund | `TestMentionsHappyPath` | **teilweise**, unverändert — Exit 0 kennt der Test nicht |
| 2 | Bezugsmenge, in **beiden** Ausgabe-Formen | `TestText_NotesAufStderrVorDerZaehlzeile`, `TestJSONYAML_MitNoteImSummary` (neu) | **teilweise, deutlich verbessert.** Beide Formen sind jetzt gemessen und beide Tests fallen (Bruch-Test H). Ungemessen bleibt die **Naht**: dass `CheckMentions` → `runPostPasses` → `report.Summary` die Note wirklich durchreicht, prüft kein Test — nur manuell |
| 3 | Boundary (Ist-Seite ist eine Menge) | `TestMentionsIstMengeIstVereinigung` | **ja**, unverändert |
| 4 | Negative (Fund) — Exit 1, `file`/`line` | `TestMentionsFindet` | **teilweise**, unverändert — Exit 1 nicht |
| 5 | Erkennungsform `basename` vs. `path` | `TestMentionsErkennungsform` | **ja**, unverändert |
| 6 | Leere Soll-Menge, fail-closed | `TestMentionsLeereSollMengeFailClosed` | **teilweise**, unverändert |
| 7 | Leere Ist-Menge, fail-closed | `TestMentionsLeereIstMengeFailClosed` | **teilweise**, unverändert |
| 8 | Negative (Config) | fünf Tests | **ja**, unverändert |
| 9 | Determinismus — byte-identisch, stabile Sortierung | `TestMentionsDeterministisch` + `TestMentionsSortierungTrenntDFSVonLexikografisch` (neu) | **teilweise, Kern geschlossen.** Die Sortierung ist jetzt lasttragend geprüft (Bruch-Test B). Die zweite Hälfte („nichts wird geschrieben", [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)) misst weiterhin kein `mentions`-Test |
| 10 | Default byte-identisch | `TestJSONYAML_OhneNoteKeinNotesSchluessel` (neu) | **teilweise.** Der einzige neue Serialisierungs-Schlüssel ist gewächtert (Bruch-Test G); „byte-identisch zur Fassung davor" über die **gesamte** Ausgabe prüft nichts |

**Was die Liste nicht mehr abbildet.** Die Anforderung hat mit 0.86.3 **drei** neue Vertrags-Aussagen bekommen; die Kriterien-Liste steht unverändert bei zehn — siehe N-7.

---

## Findings

### HIGH

**N-1 · Die neue Grenz-Prüfung meldet erwähnte Artefakte als unerwähnt — vier gemessene Formen, keine davon benannt**
`quelle`: [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) · [`AGENTS.md`](../../AGENTS.md) §3.8 · `klasse`: `haertung-kippt-fehlerpolitik-ungeprueft`
`pfad`: `internal/hexagon/core/rules/mentions.go:211-241`, `spec/lastenheft.md` §Was „kommt vor" heißt, `spec/spezifikation.md` §`DC-FA-MENT-001.a` Schritt 5
`verifizierbar`: **ja** — zwei Fixture-Läufe, unten wiedergegeben.

Der H-2-Fix prüft links und rechts auf ein „Namens-Byte" (`[A-Za-z0-9_.-]` und `/`). Damit fallen Nennungen aus der Wertung, die vorher zählten und die in gewöhnlicher Prosa der Regelfall sind. Gemessen an einer Fixture mit **18** Artefakten, je einer Nennungs-Form (Container-Lauf, `match: basename`): 15 zählen, **drei nicht** —

```
a/b10.md   in "die b10.md-Datei als kompositum"   ⇒ artifact-unmentioned
a/b11.md   in "siehe b11.md."                     ⇒ artifact-unmentioned
a/b13.md   in "_b13.md_ kursiv"                   ⇒ artifact-unmentioned
```

Der Satz-Schlusspunkt ist die schwerste der drei: `.` ist Namens-Byte **und** Satzzeichen, und *„siehe `docs/user/operations.md`."* ist die häufigste Art, eine Datei in Prosa zu nennen. Der Bindestrich trifft jedes deutsche Kompositum (`…md-Zeile`, `…md-Regeln` — eine Form, die dieses Repo selbst schreibt), der Unterstrich die Markdown-Kursivierung.

**Die vierte Form trifft die Default-Erkennungsform.** Unter `match: path` ist ein `/` links **nicht** ausgenommen; damit zählt keine `../`-relative Verlinkung mehr, obwohl sie exakt auf das Mitglied zeigt:

```
docs/plan/a.md, verlinkt als "[ADR-a](../docs/plan/a.md)"  ⇒ artifact-unmentioned
docs/plan/b.md, verlinkt als "docs/plan/b.md"              ⇒ erwähnt
```

Am eigenen Bestand gemessen (ADRs gegen `harness/README.md`, `match: path`): **0 von 84** erwähnt — vor dem Fix waren es 21, und alle 21 sind korrekte `../docs/plan/adr/…`-Links. Der Code-Kommentar begründet die Wahl mit *„`x/docs/a.md` ist eine andere Datei als `docs/a.md`"* — das trifft ein **fremdes** Präfix, nicht `../`, und die Regel kann beide nicht unterscheiden. Versagensbild: Ein Adopter mit dem üblichen Layout (Doku im Unterverzeichnis, relative Links nach oben) aktiviert `mentions` im Default und bekommt seine gesamte Soll-Menge als Befund. Keine der vier Formen steht als Grenze in Anforderung, Spezifikation, Sensor-Datei oder Modul-Kommentar. **Einstufung:** HIGH nach Prüffrage 2 (falscher Befund eines Kern-Moduls). Gegen-Argument, ausdrücklich notiert: der Fehlschlag ist **laut** — nach derselben Logik, mit der R1 M-7 auf MEDIUM setzte, wäre MEDIUM vertretbar. Ausschlaggebend ist, dass es sich um eine **Regression** der Erkennungs-Menge handelt, die der Fix ohne Preis-Angabe eingeführt hat.

**N-2 · Zwei Review-Befund-Marker stehen in `.d-check.yml` — die Klasse, die der Aufräum-Commit für beseitigt erklärt**
`quelle`: [`AGENTS.md`](../../AGENTS.md) §3.7 · Baseline `grundlagen-harness-dateien.md` §Was ein Kommentar trägt · `klasse`: `kommentar-traegt-review-historie`
`pfad`: `.d-check.yml:892` (`… ungewächtert (unabhängiger Review, H-1). …`), `.d-check.yml:905` (`… zu (Review H-2, im Modul behoben).`)
`verifizierbar`: **ja** — `grep -nE "Review|H-[0-9]" .d-check.yml`.

§3.7 gilt für **Code, Konfiguration und Skripte**; Review-Befund-Marker sind dort namentlich verboten. Die Botschaft von `3a31c8f` führt unter „ZWEI EIGENE VERSTOESSE" genau diese Klasse auf und erklärt sie für vollständig beseitigt (*„Alle entfernt"*) — für `mentions.go`, `mentions_test.go` und `workflows.go` stimmt das (nachgeprüft, dort steht keiner mehr). `.d-check.yml` wurde im **vorherigen** Commit geschrieben und im Aufräum-Commit nicht mehr angefasst; die beiden Marker sind seit `e5cf62d` neu im Baum. Derselbe Block trägt zusätzlich die Chronik der verworfenen Fassung (*„die Grenze, an der die erste Fassung dieser Konfiguration scheiterte: Sie führte ZWEI Invarianten in einem Block"*) — Deliberation über Verworfenes, dieselbe verbotene Klasse; die **Grenze** selbst (*„Die Ist-Menge wird als VEREINIGUNG gelesen"*) gehört dagegen dorthin. Versagensbild: der nächste Leser der Konfiguration sucht einen Report, der außerhalb dieses Repos nicht mitreist, und die Regel, die diese Marker verbietet, ist im selben Slice zweimal zitiert worden.

### MEDIUM

**N-3 · Der einzige Test des `scan.ignore`-Fixes kann nicht fallen — die Fixture benutzt `vendor/`, und das prunt schon die feste Skip-Liste**
`quelle`: [`BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt`](../plan/planning/observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md) · Reviewer-Anker „fehlende Negativtests bei neuem öffentlichen Vertrag" · `klasse`: `wortlaut-behauptet-pruefung-die-fehlt`
`pfad`: `internal/hexagon/core/rules/mentions_test.go` (`TestMentionsHonoriertScanIgnore`) gegen `internal/hexagon/core/rules/scan.go:13-21`
`verifizierbar`: **ja** — Bruch-Test D.

`TestMentionsHonoriertScanIgnore` legt `vendor/fremd.sh` an und erwartet, dass `scan.ignore: ["vendor/**"]` es aus der Soll-Menge nimmt. `isSkipDir` führt `vendor` aber bereits in seiner **festen** Liste — das Verzeichnis wird geprunt, ob `ignore` ausgewertet wird oder nicht. Gemessen: Nach vollständigem Entfernen **beider** `ignore`-Auswertungen aus `mentionsWalk` (`dirIgnored` bei Verzeichnissen, `ignored` bei Dateien) bleibt `make test` **grün**; kein Test fällt. Das ist genau die Fixture-Klasse, die der Slice-Plan §7 mit 7× Belegen als „unmittelbar einschlägig für die Tests" gesichtet hat, und die der M-3-Fix im selben Commit an anderer Stelle aufgelöst hat. Das **Produkt** ist in Ordnung — mit einem nicht geskippten Verzeichnis gegengeprüft: `scan.ignore: ["fremd/**"]` nimmt `fremd/fremd.md` wirklich aus der Menge. Ungewächtert ist die Zusage, nicht das Verhalten.

**N-4 · Der Lastenheft-Wortlaut der Grenz-Prüfung ist weiter als Spezifikation und Code — bei einem Nicht-ASCII-Nachbarn besteht die H-2-Klasse fort**
`quelle`: [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) (Rang 1) gegen `spec/spezifikation.md` §`DC-FA-MENT-001.a` Schritt 5 (Rang 2) · `klasse`: `zaehlmethode-misst-proxy-statt-gegenstand`
`pfad`: `spec/lastenheft.md` §Was „kommt vor" heißt, `internal/hexagon/core/rules/mentions.go:232-242`
`verifizierbar`: **ja** — Fixture-Lauf unten.

Das Lastenheft sagt zu: *„Unmittelbar davor und dahinter darf kein Zeichen stehen, das selbst Teil eines Datei- oder Pfadnamens sein kann."* Die Spezifikation nennt dafür die **Aufzählung** `[A-Za-z0-9_.-]` und `/`, und `mentionsNameByte` arbeitet **byte-weise** über genau diese ASCII-Bereiche. Ein Umlaut ist ein Zeichen, das Teil eines Dateinamens sein kann — seine UTF-8-Bytes liegen aber außerhalb der Bereiche und blockieren die Grenze nicht. Gemessen (`match: basename`):

```
Soll-Menge: s/test.md, s/haupt.sh
Ist-Dokument: "hier steht Ätest.md …" / "und hier haupt.shö als endung"
⇒ d-check: mentions: 2 von 2 Artefakt(en) erwähnt, 0 Befund(e), Exit 0
```

Beide Mitglieder gelten als erwähnt, obwohl keines genannt ist — dasselbe **stille Grün**, für das R1 H-2 als HIGH führte, nur mit einem Nicht-ASCII-Nachbarn statt eines Bindestrichs. Die Klasse ist also verkleinert, nicht geschlossen. Die drei Ebenen widersprechen sich dabei absteigend (Lastenheft ⊃ Spezifikation ⊃ Code); nach Source Precedence müsste die höherrangige Aussage stehen bleiben und der Code ihr folgen, oder das Lastenheft die Aufzählung übernehmen — die Entscheidung fehlt. **Dieselbe Byte-Basis ist zugleich der Grund, warum die legitimen Formen `„a.sh"` und `»a.sh«` funktionieren** (gemessen, beide zählen): die Regel ist an dieser Stelle absichtlich grob und muss es benennen.

**N-5 · Die Historie-Zeile 0.86.3 behauptet einen Beleg als „nachgeholt", den die Anforderung seit 0.86.1 führt**
`quelle`: [`AGENTS.md`](../../AGENTS.md) §5 (nicht mehr behaupten, als die Arbeit trägt) · `klasse`: `commit-message-overclaims-work`
`pfad`: `spec/lastenheft.md` §Historie, Zeile 0.86.3, Schluss von Punkt (5); Gegenbeleg: Zeile 0.86.1 derselben Tabelle und der Absatz §Die Mengen-Wahl ist das Urteil in `ac12993`
`verifizierbar`: **ja** — `git show ac12993:spec/lastenheft.md` trägt den Fremd-Bestands-Satz bereits.

Die neue Zeile schließt: *„Damit ist zugleich der Kalibrierungs-Beleg an einem **fremden** Bestand nachgeholt, den die Anforderung bisher nur am eigenen führte."* Die Anforderung führte ihn bereits: Zeile 0.86.1 sagt wörtlich *„an ihre Stelle tritt eine zweite Pfad-Menge an einem Fremd-Repo (`scripts/*.sh` + `tools/*.py`: vier Funde bei fünf Mitgliedern)"*, und derselbe Satz stand vor diesem Commit im Anforderungstext. **Nachgeholt ist etwas anderes und es ist echt** — der Form-Vergleich (beide Erkennungsformen liefern dasselbe) und das Urteil über die vier Funde; beides habe ich gegen den Fremd-Bestand reproduziert (5 Mitglieder, 4 Funde, `path` und `basename` identisch, einer der vier ist eine Test-Datei). Beanstandet ist der Schluss, nicht die Messung: er macht aus einer **Ergänzung** eine **Erst-Erfüllung** und entwertet damit die eigene Historie-Zeile 0.86.1.

**N-6 · Die Erkennungsform-Messung im Anforderungstext ist unter der neuen Semantik nicht nachgezogen — und ihre Stichprobe ist nicht benannt**
`quelle`: [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) · [`AGENTS.md`](../../AGENTS.md) §5 · `klasse`: `selbstauskunft-zahl-reproduziert-nicht`
`pfad`: `spec/lastenheft.md` §Was „kommt vor" heißt, Satz *„über zwölf Artefakte aus drei Repos weichen die Formen **einmal** ab"*
`verifizierbar`: **ja** — die Messung (a) unten; die Stichprobe selbst ist nicht reproduzierbar, und das ist der zweite Teil des Befundes.

Der Satz stützt die Default-Wahl `path` mit einer Messung, die unter `strings.Contains` entstand. Genau diese Vergleichs-Grundlage ändert 0.86.3. Am eigenen Bestand ist die Divergenz zwischen den Formen durch den Fix **gewachsen**: ADRs gegen `harness/README.md` liefern unter `basename` 21 Treffer und unter `path` jetzt **0** (vor dem Fix ebenfalls 21) — 21 neue Abweichungen an einem einzigen Paar. Ob „einmal" über die zwölf Artefakte noch stimmt, ist damit offen; die Zeile 0.86.3 zählt unter Punkt (5) auf, welche Zahlen neu gemessen wurden, und diese ist nicht dabei. Erschwerend: **welche zwölf Artefakte aus welchen drei Repos** steht nirgends, also lässt sich der Wert weder bestätigen noch widerlegen — bei den beiden anderen Fremd-Messungen (fünf Mitglieder, acht Dokumente) war das Nachfahren möglich und hat funktioniert.

**N-7 · Drei neue Vertrags-Aussagen der Anforderung haben kein Akzeptanzkriterium**
`quelle`: [`BEO-ALL/spec-randbedingung-ohne-test`](../plan/planning/observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md) · `klasse`: `spec-randbedingung-ohne-test`
`pfad`: `spec/lastenheft.md` §Akzeptanzkriterien von `DC-FA-MENT-001` (unverändert zehn Punkte) gegen §Historie 0.86.3 Punkte (1), (2), (3)
`verifizierbar`: **ja** — der Diff fasst den Kriterien-Block nicht an; nachgezählt bleiben es zehn.

0.86.3 legt drei prüfbare Dinge neu fest: *eigenständige Nennung statt Teilzeichenkette* (mit der ausdrücklichen Ausnahme des `/` links unter `basename`), *ein Block ist ein Paar* und *unlesbares Verzeichnis ⇒ Exit 2*. Für alle drei gibt es Unit-Tests — aber die Kriterien-Liste der Anforderung, die die Verifikation abhakt und die R1 als L-1 gerade erst nachgezählt hat, kennt sie nicht. Damit ist der zentrale Vertragswechsel dieses Bumps auf der Ebene, die ihn festhält, unsichtbar: Wer die zehn Kriterien abhakt, hat die Grenz-Prüfung nicht abgehakt. Der neue Exit-2-Pfad steht zusätzlich neben zwei Kriterien, die die *anderen* Exit-2-Pfade einzeln aufführen — die Auslassung fällt genau dort auf, wo die Systematik sonst vollständig ist.

**N-8 · Handbuch §5 und der neue `matchAnyGlob`-Kommentar sagen jetzt das Gegenteil voneinander**
`quelle`: [`DC-FA-CONF-001`](../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei) · `klasse`: `semantic-change-body-only-edges-stale`
`pfad`: `docs/user/benutzerhandbuch.md:1211-1216` gegen `internal/hexagon/core/rules/workflows.go:258-263`
`verifizierbar`: **ja** — beide Stellen nebeneinander lesen.

Der L-5-Fix schreibt in den Kommentar: *„GRENZE: `**` traegt hier nicht ueber mehrere Segmente … Wer die Semantik von scan.ignore braucht, nimmt matchGlob"* — für `workflows.exempt-paths` und `reviews.exempt-paths`. Das Handbuch sagt im selben Baum: *„`**` steht für beliebig viele Segmente … Das gilt für alle Glob-Felder — `scan.ignore`, `matrix.classes[].paths`/`.order` und die `exempt-paths` der Module."* R1 hatte den Bestandsbefund als I-3 notiert; seit diesem Commit ist er kein stiller Bestand mehr, sondern ein **ausgeschriebener Widerspruch** zwischen einer Nutzer-Zusage und einer Code-Grenze. Dieselbe Aufzählung ist doppelt veraltet: `mentions.artifacts`/`mentions.documents` sind neue Glob-Felder, für die `**` trägt, und sie stehen dort nicht. Versagensbild: Ein Adopter schreibt `exempt-paths: ["docs/**/alt-*.yml"]`, das Handbuch verspricht Wirkung, und die Ausnahme greift still nicht.

### LOW

**N-9 · Die zuvor als gehalten deklarierte Sensors-Tabellen-Invariante entfällt ersatzlos — ohne Träger**
`pfad`: `.d-check.yml` (Soll-/Ist-Menge um je einen Eintrag verkürzt), `harness/sensors/mention-coverage.md` §Grenze · `quelle`: Baseline `modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst · `klasse`: `registerzeile-ohne-ausgang-nach-schwelle`
Der H-1-Fix ist der richtige, und die Sensor-Datei benennt die Folge ehrlich (*„Die Sensors-Tabellen-Invariante ist deshalb hier ausdrücklich NICHT gehalten; sie bräuchte einen zweiten Lauf"*). Was fehlt, ist der Träger dieses „bräuchte": kein Folge-Slice, kein Eintrag im Beobachtungs-Register, kein Risiko in §5 des Plans (`git log` zeigt: `docs/plan/planning/observations/` ist seit slice-205 unberührt). Eine benannte, nicht gehaltene Invariante ohne Vorgang, der sie einlöst, ist genau die Form, die beim nächsten Aufräumen verschwindet. `verifizierbar`: ja.

**N-10 · Das fail-closed bei unlesbarem Verzeichnis gilt dem ganzen Baum — auch dort, wo kein Glob je treffen könnte**
`pfad`: `internal/hexagon/core/rules/mentions.go:138-142` in Verbindung mit dem `scan.roots`-Verzicht · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.8 · `klasse`: `haertung-kippt-fehlerpolitik-ungeprueft`
Gemessen: ein Verzeichnis `geheim/` mit Modus `000`, das ausschließlich eine `.txt` enthält, bringt einen Lauf mit `artifacts: ["**/*.md"]` auf `d-check: error: mentions: Verzeichnis "geheim" nicht lesbar … Exit 2`. Weil das Modul bewusst **nicht** auf `scan.roots` eingeschränkt ist, ist die Fläche dieses fail-closed die gesamte Repository — größer als bei jedem Geschwister-Modul und unabhängig davon, ob das Verzeichnis für die konfigurierten Mengen überhaupt relevant wäre. Im Container läuft d-check als `nonroot`; ein Verzeichnis, das dieser UID nicht offensteht, genügt. Das Verhalten ist zugesagt (Spezifikation Schritt 2), die Fläche nicht. `verifizierbar`: ja.

**N-11 · Ein `scan.ignore`-Muster kann die Ist-Menge leeren; die Fehlermeldung zeigt dann auf den falschen Schalter**
`pfad`: `internal/hexagon/core/rules/mentions.go:74-77` · `klasse`: `gate-rot-aus-falschem-grund`
Neue Kopplung durch den M-5-Fix: `documents` wird aus derselben, jetzt `scan.ignore`-geprunten Liste gefiltert. Gemessen mit `scan.ignore: ["docs/**"]` bei `documents: [docs/doc.md]` — die Datei existiert und das Glob trifft sie: `d-check: error: mentions.documents [docs/doc.md] trifft kein Dokument …`, Exit 2. Die Meldung nennt das Glob als Ursache; die Ursache ist `scan.ignore`. Laut, also kein stilles Grün — aber der Leser dreht am falschen Schalter. `verifizierbar`: ja.

**N-12 · Die Test-Zahlen der Botschaft reproduzieren nicht**
`pfad`: Commit-Botschaft `e5cf62d`, letzter Absatz · `quelle`: [`AGENTS.md`](../../AGENTS.md) §5 · `klasse`: `selbstauskunft-zahl-reproduziert-nicht`
*„Sieben neue Tests, davon drei als Regression gegen H-2 und H-3."* Gezählt sind es **elf** neue Tests (sieben in `mentions_test.go`, vier in `report_test.go` — die Botschaft nennt die vier zwei Absätze vorher selbst) und **vier** Regressionen (drei gegen H-2, eine gegen H-3). Die Richtung stimmt, die Ziffern nicht — dieselbe Klasse, die R1 als M-1 und L-1 führte und die R1 als „siebten Lauf in Folge" zählte. `verifizierbar`: ja.

### INFO

- **I-1 · Die `seen`-Map in `mentionsFilter` ist strukturell tot.** Seit dem L-4-Fix stammt `all` aus **einem** Walk und ist damit duplikatfrei; `seen[rel]` kann nie `true` werden. Der Kommentar sagt weiterhin „dedupliziert (DC-QA-02)" — die Zusage stimmt, aber sie wird jetzt vom Walk getragen, nicht von dieser Zeile. Kein Defekt, nur eine Zeile, die verspricht, was sie nicht mehr tut.
- **I-2 · `Config.Mentions` trägt weiter keinen Feld-Kommentar** (`internal/hexagon/core/model/config.go`), während `Tracked`, `Targets`, `Workflows` und `Reviews` je eine `// <Name>: Parameter des Moduls … (DC-…)`-Zeile führen. R1-I-1 unverändert; §3.7 regelt Kommentare, die es gibt, nicht ihre Abwesenheit.
- **I-3 · Der Konfigurations-Referenzblock des Benutzerhandbuchs dokumentiert `mentions:` weiterhin nicht** (R1-I-2 unverändert) — die Modul-Tabelle ist nachgezogen, der Referenzblock nicht. Ein Nutzer, der das Modul dort aktiviert, findet die Pflicht-Schlüssel nirgends und bekommt Exit 2.
- **I-4 · Paralleler Schreibzugriff während des Laufs.** Die Slice-Pläne `slice-207` und `slice-208` unter `docs/plan/planning/open/` sind während dieses Reviews von außerhalb entstanden und als `af29684` committet worden. Sie gehören nicht zum Gegenstand und wurden nicht angefasst; sie stehen hier, damit der nächste Leser die Commit-Reihenfolge richtig liest.

---

## Negativbefunde

- **Alle Gate-Läufe der beiden Botschaften sind echt und reproduzieren.** `make gates` Exit 0 über zehn Gates, `677 Datei(en) geprüft, 0 Befund(e)`, `Coverage 94.60%`, semgrep `0 findings`; `make mention-coverage` Exit 0 mit `84 von 84 … über 1 Dokument(e)`. Die Zahlen 94,60 % und 677 stehen wörtlich so in `e5cf62d`.
- **Alle sechs Kalibrierungs-Zahlen reproduzieren einzeln.** `docs/plan/adr/[0-9]*.md` ⇒ **84**; `harness/sensors/*.md` ⇒ **24**; `match: path` gegen den ADR-Index ⇒ **84 Befunde**; `match: path` für die Sensor-Dateien gegen `harness/README.md` ⇒ **24 Befunde**; `match: basename` für die ADRs gegen `harness/README.md` ⇒ **21 von 84**; `tools/**/*.sh` ⇒ **11 Mitglieder, 11 Funde** (5 flach, 6 eine Ebene tiefer). Die R1-Beanstandung „23 statt 24" ist berichtigt, „ausnahmslos falsch" trägt jetzt, weil je Paar statt gegen die Vereinigung gemessen wird.
- **Der Fremd-Bestands-Beleg ist echt und exakt nachgefahren.** `scripts/*.sh` + `tools/*.py` gegen `docs/user/*.md` eines Fremd-Repos: **fünf** Mitglieder, **acht** Dokumente, **vier** Funde, `path` und `basename` **identisch** — und der einzige Treffer wird dort tatsächlich mit vollem Pfad genannt, wie die Begründung behauptet. Einer der vier ist eine Test-Datei. Nur der Rahmen-Satz um diese Messung stimmt nicht (N-5).
- **Der H-2-Fix schließt den gemeldeten Fall wirklich.** Fixture mit `s/test.md` und `s/image-test.md`, nur letzteres verlinkt: genau **ein** Befund, auf `s/test.md`, und die Gegenprobe zählt: der Pfad-Präfix `s/` vor dem Basisnamen blockiert nicht.
- **Der H-3-Wechsel auf `matchGlob` hat am ausgelieferten Bestand keine Nebenwirkung.** `docs/plan/adr/[0-9]*.md` und `docs/plan/adr/README.md` enthalten kein `**`; `matchGlob` zerlegt segmentweise und ruft je Segment dasselbe `path.Match` — die Mengen bleiben 84 bzw. 1, gemessen vor und nach der Umstellung. Auch die Ungültigkeits-Politik ist gleich geblieben (ein fehlerhaftes Muster matcht nicht, statt zu werfen).
- **`mentionsFilter` ist gegenüber der Zwei-Pass-Fassung verhaltensgleich.** Dieselbe Iteration über dieselbe Walk-Reihenfolge, dieselbe `seen`-Prüfung, dasselbe abschließende `sort.Strings`; die **Reihenfolge der Globs** ist ohne Wirkung, weil je Datei über alle Globs disjunktiv geprüft wird. Der einzige Unterschied ist, dass der Walk einmal statt zweimal läuft — das war der Zweck.
- **Die Grenz-Prüfung hält an den geprüften Rändern.** Aus 18 Nennungs-Formen zählen 15: fett (`**x.md**`), eckige Klammer, Backticks, `x.md:12`, runde Klammer, Zeilenanfang, Zeilenende, deutsche und französische Anführungszeichen, Pfad-Präfix unter `basename`, Komma, Sternchen-Kursiv, spitze Klammern, alleinstehend. Ein Treffer am Korpus-Ende ist abgedeckt (`end == len(corpus)`), Groß-/Kleinschreibung bleibt wie zuvor signifikant, und `a.sh` in `a.sh.bak` wird korrekt **nicht** gezählt. Die drei Ausnahmen stehen in N-1.
- **Der Determinismus-Fix ist keine Kosmetik.** Die neue Fixture trennt die Ordnungen wirklich: `tools.md` steht lexikografisch vor `tools/a.sh` (`.` < `/`), die Tiefensuche liefert es danach — ohne `sort.Strings` fällt der Test (gemessen).
- **Die vier Notes-Tests hängen an vier verschiedenen Eigenschaften.** Reihenfolge auf stderr, Einzeiligkeit über `oneLine`, Abwesenheit des Schlüssels ohne Note (`omitempty`) und Anwesenheit in **beiden** maschinenlesbaren Formen — jede einzeln durch einen Eingriff zum Fallen gebracht.
- **§3.5 ist gewahrt** — [ADR-0084](../plan/adr/0084-mentions-eigenes-modul.md) ist in beiden Commits unberührt; der Diff zeigt die Datei nicht.
- **§3.1/§3.2 sind gewahrt** — kein Host-Werkzeug in den Commits, keine `//nolint`-Direktive; alle eigenen Läufe dieses Reviews liefen über `make` bzw. `docker run` gegen das gebaute Image.
- **§3.4 ist gewahrt** — die neuen Spec-Absätze nennen keine ADR, keinen Slice, keine Welle, keinen Commit-Hash; `matrix` ist grün.
- **§3.7 im Go-Code ist sauber.** In `mentions.go`, `mentions_test.go`, `report_test.go` und `workflows.go` steht kein Review-Befund-Marker mehr (`grep` über alle vier). Die verbliebenen Kommentare tragen Zusage, Abgrenzung und Grenze; Herkunft steht je als **ein** auflösbares Feld. Die Ausnahme ist die Konfiguration, siehe N-2.
- **Kein Netzzugriff, kein Schreibpfad** — die Änderungen benutzen ausschließlich `driven.Filesystem.List`/`ReadFile`; `make arch-check` ist grün, alle Läufe fuhren `--network none` und `-v …:ro`.
- **Die Doku-Ränder des H-1-Fixes widersprechen einander nicht.** Die Makefile-Hilfe nennt die gehaltene Invariante (*„jede ADR steht im ADR-Index"*), die Sensor-Datei nennt sie und dazu die **nicht** gehaltene, die Sensors-Tabelle in `harness/README.md` beschreibt sie generisch. `AGENTS.md` beschreibt weiterhin nur die Mechanik und verweist für die Mengen-Wahl auf die Sensor-Datei — unvollständig, aber kein Widerspruch und dieselbe Arbeitsteilung wie vor dem Fix.
- **Die `basename`-Begründung ist jetzt am richtigen Gegenstand gemessen.** Über die 84 ADR-Basisnamen gibt es 0 Duplikate **und** 0 Fälle, in denen einer Endstück eines anderen ist (beides nachgezählt) — die zweite Prüfung, deren Fehlen R1 als H-2 meldete, ist damit belegt und nicht nur gefordert.
- **Die §2-Schema-Zeilen decken sich mit dem Code.** `applyMentions` erzwingt: beide Schlüssel oder keiner, kein Weißraum-Glob, gültiges Glob beidseitig, `match` aus zwei Werten — je Exit 2, genau wie die drei neuen Tabellenzeilen sagen.

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 6 |
| LOW | 4 |
| INFO | 4 |

**Die Klassen dieses Laufs.** Zwei der drei R1-Wiederholungsklassen kehren wieder, aber verschoben: `selbstauskunft-zahl-reproduziert-nicht` betrifft jetzt nicht mehr die Kalibrierung (die reproduziert vollständig), sondern die **Test-Zahl** und eine **nicht nachgezogene** Messung (N-12, N-6); `wortlaut-behauptet-pruefung-die-fehlt` trifft genau den Test, der einen der Fixes belegen sollte (N-3). Neu und die wichtigste Beobachtung: **die Härtung selbst wird nicht gegen ihre Fehlerpolitik geprüft** — der H-2-Fix ist ausschließlich mit Fällen belegt, die er fangen **soll** (drei Regressionstests), und mit keinem einzigen, den er weiterhin **durchlassen** muss. Das ist die Gegenrichtung aus `modul-11` §Fitness Function ohne Standard-Tool, und sie fehlt: der unveränderte Bestand, auf dem der Sensor schweigen muss, ist für die Grenz-Prüfung nie erhoben worden. Kandidat für einen Registereintrag oder eine Evidenz-Datei an [`BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft`](../plan/planning/observations/BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft/observation.md).

## Verdikt

**Blockierend — N-1 und N-2.**

Die Korrekturen sind in der Sache gut: **13 der 17** R1-Findings sind wirklich behoben, jeder der drei HIGH-Bruch-Tests kehrt sich um, und zehn der elf neuen Tests können nachweislich fallen. Blockierend bleiben zwei Punkte. **N-1** ist die Kehrseite des H-2-Fixes: Vier gemessene, in Prosa und Markdown alltägliche Nennungs-Formen zählen nicht mehr, darunter der Satz-Schlusspunkt und — unter dem **Default** `path` — jede `../`-relative Verlinkung; keine davon ist als Grenze benannt. **N-2** ist ein §3.7-Verstoß derselben Klasse, die der zweite Commit ausdrücklich für beseitigt erklärt, an einer Datei, die er nicht mehr angefasst hat.

**Nicht blockierend, aber vor der Closure zu entscheiden:** N-3 (der Test des `scan.ignore`-Fixes kann nicht fallen — die Zusage ist ungewächtert, das Verhalten nicht), N-4 (Lastenheft, Spezifikation und Code sagen drei verschieden weite Grenzen zu; die H-2-Klasse besteht bei Nicht-ASCII-Nachbarn fort) und N-9 (eine benannte, nicht gehaltene Invariante ohne Träger).
