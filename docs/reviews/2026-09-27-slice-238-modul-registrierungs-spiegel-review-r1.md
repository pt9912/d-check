# Review-Report: slice-238 — 2026-09-27

**Review-Art:** Plan/Code — geprüft gegen den Slice-Plan
(`slice-238-modul-registrierungs-spiegel-checkliste.md`: §1 Ziel und
Abgrenzung/Option A, §2 DoD, §6 Risiken), `AGENTS.md` §3 (Hard Rules,
insbes. §3.7/§3.8) und §5, die Beobachtung
`BEO-ALL/modulliste-spiegel-ungegated` (observation.md, state.md,
evidence/slice-115.md, evidence/slice-152.md, evidence/slice-236.md) sowie
die Reviewer-Skill-Prüffragen (Modul 10).

**Gegenstand:** slice-238 — Commit-Range `bc54d396..7702c593` (drei Commits:
`0394acab` reiner `git mv` open→in-progress, `8d02b273` Beanspruchung/dritter
Vorprüfungsblock, `7702c593` Implementierung: `internal/hexagon/core/app/registry_mirror_test.go`
neu, `internal/adapter/driven/configyaml/gate_consistency_test.go` erweitert).
Arbeitsbaum gleich HEAD (`7702c593`).

**Skill:** `.harness/skills/reviewer.md` @ Version 1.16.0
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-27

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-238-modul-registrierungs-spiegel-checkliste.md` (§1
  Ziel/Abgrenzung/Option A vs. B, §2 DoD, §6 Risiken)
- `BEO-ALL/modulliste-spiegel-ungegated` (observation.md, state.md,
  evidence/slice-115.md, evidence/slice-152.md, evidence/slice-236.md) —
  der Auslöser und seine drei Vorgänge
- `AGENTS.md` §3.7 (Kommentar-Disziplin, fünf Klassen, Verbot von
  Review-Historie/Review-Befund-Markern), §3.8 (Modul verspricht nur über
  Scan-Menge), §5 (Grenzen-Regel — Grenzen-Liste ohne größte Lücke)
- Vorheriger Review-Report `docs/reviews/2026-09-27-slice-236-file-modul-review-r1.md`
  (F-2/F-3/F-4 — die drei Fundorte, die dieser Slice mechanisch decken soll)
- Unabhängig nachgefahrene Gate-Läufe: `make test` (isoliert und via
  `make gates`), `make gates` vollständig, sowie gezielte `go test -v -run`
  gegen die sechs neuen/geänderten Testfunktionen
- Drei selbst gebaute adversarielle Gegenproben (mutierte Arbeitskopie,
  siehe Findings/Negativbefunde)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Drei neu geschriebene Kommentare zitieren Review-Runden-/Befund-Marker einer **anderen** Rolle bzw. eines **anderen** Slice als Design-Begründung: „Deckt genau die drei Orte, die slice-236s Review nachtragen musste (F-2/F-3/F-4)" (Review-Historie + Befund-Marker aus dem Review von slice-236) und zweimal „(slice-057-R3-Lehre: nur der Guard löst den Befund aus)" (Review-Runde 3 eines fremden Slice als Lehre zitiert). Beides sind Review-Historie/Review-Befund-Marker — genau die Klasse, die `AGENTS.md` §3.7 und der HIGH-Anker des Reviewer-Skills ausdrücklich verbieten: „Keine Review-Historie und keine Review-Befund-Marker … Herkunft ist nur als **ein** auflösbares Feld zulässig". Der Verweis auf `slice-236`/`BEO-ALL/…` allein (Herkunfts-Anker) wäre zulässig; die zusätzliche Nennung von `F-2/F-3/F-4` und `R3` überschreitet das erlaubte eine Feld. Das Muster ist im Bestand verbreitet (33 Fundstellen repo-weit), aber diese drei Stellen sind **Neuzugänge** in brandneuen Dateien/Funktionen dieses Slice — §3.7 nimmt nur Alt-Code vor der Schärfung aus, nicht neue Kopien des Musters: „Neuzugänge fallen überall unter den Anker." | `AGENTS.md` §3.7; Reviewer-Skill Frage 6 | `internal/hexagon/core/app/registry_mirror_test.go:17` (F-2/F-3/F-4), `:39` (slice-057-R3-Lehre, neue Funktion `moduleMirrorIssues`); `internal/adapter/driven/configyaml/gate_consistency_test.go:171` (slice-057-R3-Lehre, neue Funktion `focusDisableIssues`) | ja — Zeilen liegen vollständig im Diff `bc54d396..7702c593`, Wortlaut per `grep -n "F-2/F-3/F-4\|slice-057-R3-Lehre" …` reproduzierbar | Review-Historie/Review-Befund-Marker im Code-Kommentar (Neuzugang) |
| F-2 | MEDIUM | Zwei der drei bekannten, wörtlich-vollständigen `model.ValidModules()`-Spiegel außerhalb der drei jetzt gedeckten Fundorte bleiben ungedeckt: (a) `README.md` und `README.de.md` führen je eine vollständige Bullet-Liste, ein Eintrag pro Modul (24/24 Module, exakte `` `modul` `` -Backtick-Form, geprüft durch Auszählen), (b) `docs/user/benutzerhandbuch.md` §6 „Regelmodule" führt dieselben 24 Module als Tabellenzeilen (Spalte 1). Der Slice-Plan §1 nennt „README-Bullets" als eine der Kategorien, für die Option A ausdrücklich trägt („Trägt für Prosa-Listen (Glossar, Bereichskürzel, README-Bullets) nur, wenn die Namen dort wörtlich und vollständig stehen"), und §2 DoD verlangt, `docs/user/*.md` und `README*.md` „durch tatsächliches Lesen" als Fundort-Kandidaten zu prüfen. Weder die vier neuen Tests noch ein bestehender Test hält diese drei Dateien gegen `model.ValidModules()`. Damit bleibt der wahrscheinlichste **nächste** Ort für die vierte Evidenz von `BEO-ALL/modulliste-spiegel-ungegated` exakt der, den dieser Slice laut eigenem Kriterium hätte decken sollen. | Plan §1 (Option-A-Kriterium „README-Bullets"), §2 DoD (Fundort-Liste); `AGENTS.md` §5 (Grenzen-Liste ohne größte Lücke); Reviewer-Skill Frage 2/18 | `README.md:23-215` (24 Bullets), `README.de.md:22-217` (24 Bullets), `docs/user/benutzerhandbuch.md:2616-2639` (24 Tabellenzeilen `## 6. Regelmodule`) | ja — Auszählung `grep -c "^- \`[a-z]*\` " README.md README.de.md` liefert je 24, `grep -c "^| \`[a-z]*\` " docs/user/benutzerhandbuch.md` liefert 24; `model.ValidModules()` (`internal/hexagon/core/model/config.go:25`) liefert ebenfalls 24 Namen; kein neuer Test in diesem Diff referenziert eine der drei Dateien | Spiegel unvollständig — plan-eigenes Kriterium nicht auf alle benannten Fundorte angewendet |

## Negativbefunde (die sechs Prüfpunkte des Auftrags)

| Prüfpunkt | Ergebnis |
|---|---|
| 1. Robustheit der vier Regex-Extraktionen (Einfügen mitten in die Kette, Verschwinden der Ankerphrase) | geprüft, ohne Befund über F-1/F-2 hinaus. Drei adversarielle Gegenproben in einer mutierten Arbeitskopie: (a) ein unbekanntes Backtick-Token mitten in die `DC-FA-CLI-002`-Kette eingefügt (`` `matrix`, `zzznew`, `external` ``) → `TestModulRegistrySpiegel_LastenheftCLI002` schlägt korrekt fehl (`Eintrag "zzznew" ist kein bekanntes Modul`); (b) ein echtes neues Modul nur in `model.ValidModules()` ergänzt, Prosa unverändert gelassen → alle drei `TestModulRegistrySpiegel_*`-Tests schlagen korrekt mit `Modul "zzznewmod" fehlt` fehl; (c) die Ankerphrase in `spec/lastenheft.md` redaktionell umformuliert (`gegliedert:` → `besteht aus den Bausteinen:`) → `TestModulRegistrySpiegel_LastenheftCLI002` bricht fail-closed mit `t.Fatal("… nicht gefunden (fail-closed)")`, nicht still grün. `FindStringSubmatch`/`FindSubmatch` werden an allen vier Fundstellen (drei in `registry_mirror_test.go`, eine in `gate_consistency_test.go`) explizit auf `nil` geprüft, bevor auf das Ergebnis zugegriffen wird — kein ungeprüfter Zugriffspfad gefunden. |
| 2. Fehlen weitere echte, mechanisch tragfähige Fundorte in der vom Plan genannten längeren Liste? | F-2 beantwortet dies mit ja (README.md/README.de.md/benutzerhandbuch.md §6). Die übrigen vom Plan genannten Orte wurden einzeln geprüft, ohne weiteren Befund: `harness/README.md` und `AGENTS.md` führen keine vollständige Backtick-Kette (`grep` liefert keinen Treffer); `spec/spezifikation.md` §2/§4 führt nur das **fixe 7-Modul-Standardset** (`links, anchors, ids, matrix, codepaths, spans, hostpaths`, Zeile 136), identisch zur Konstante in `renderHarness` (`internal/hexagon/core/app/suggest.go:416`) — eine bewusste Teilmenge, keine Vollständigkeits-Behauptung; `internal/adapter/driving/cli/config_template.go` und `internal/hexagon/core/app/diagnose.go` sind bereits durch die bestehenden Tests `TestPrintConfigVerfuegbarDecktRegistry`/`TestAllReasonsDeckungGegenSpezifikationGrundCodes` gewächtert (unverändert in diesem Diff, per Plan §1 explizit als bereits gedeckt benannt). |
| 3. Ist der Ausschluss von Bereichskürzel-Liste und „fixem Standard-Modulset" sachlich begründet? | geprüft, ohne Befund — die Entscheidung trägt. `spec/lastenheft.md` §3 Zeile 59f. führt **Abkürzungen** (`CLI`, `SCAN`, `LINK`, `ANCH`, `ID`, `MTX`, …), keine wörtlichen Modulnamen — eine 1:1-Prüfung gegen `model.ValidModules()` wäre kategorisch falsch (`LINK` ≠ `links`). Das „fixe Standard-Modulset" (`internal/hexagon/core/app/suggest.go:416`, sieben Module, für den `ai-harness`-Vorlagenmodus) ist im Code nachweislich eine Konstante unabhängig von `model.ValidModules()`, und die zugehörigen Doku-Stellen (`spec/spezifikation.md:136`, `docs/user/operations.md:165`, `docs/user/benutzerhandbuch.md:269`) benennen es ausdrücklich als „fixes Standard-Modulset" — keine Vollständigkeits-Zusage, also zu Recht außerhalb einer Vollständigkeitsprüfung. |
| 4. `assertNetlessModules`s dritte Prüfrichtung — schließt sie einen echten, vorher unentdeckten Fehlermodus? | geprüft, ohne Befund. Vor diesem Diff prüften die beiden bestehenden Schleifen nur `netlessDocModules() ⊆ modules` und `forbiddenInNetless() ∩ modules = ∅`; ein drittes, weder gelistetes noch verbotenes Modul (z. B. das reale `mentions`, das in keiner der beiden Listen steht) fiel nachweislich durch beide Maschen — verifiziert durch Nachvollzug der Vor-Diff-Logik gegen den neuen Testfall `"unbekanntes Modul gesetzt (slice-238)"` (`append(full, "mentions")`), der isoliert gegen die alte Zwei-Schleifen-Logik grün geblieben wäre. Der neue dritte Block (`gate_consistency_test.go:56-70`) und der Guard-Fall (`gate_consistency_test.go:227`) treffen exakt diesen Fall; `TestQA03_NetlessModuleList_Guards` lief isoliert gegen den Docker-Build grün mit allen Subtests inkl. des neuen. |
| 5. Deckt `TestFocusDisable_DecktDCheckYmlModules` einen vorher komplett ungeprüften Fundort? | geprüft, ohne Befund. Der Vorgänger-Review-Report `docs/reviews/2026-09-27-slice-236-file-modul-review-r1.md` (Negativbefund #8) hält ausdrücklich fest: „`.d-check.yml`/`FOCUS_DISABLE`-Paar (bewusst unberührt, ADR-Konsequenz korrekt)" — also zu diesem Zeitpunkt bekannt ungedeckt. `git grep FOCUS_DISABLE` gegen den Stand vor diesem Slice (`bc54d396`) findet keinen Test, der `FOCUS_DISABLE` gegen `.d-check.yml` hält. `TestFocusDisable_DecktDCheckYmlModules` ist damit eine echte Neu-Deckung, kein Doppel-Test. |
| 6. `make gates`/`make test` unabhängig nachgefahren; laufen die neuen/geänderten Tests wirklich? | geprüft, ohne Befund. `make test` (Docker-Build, `CGO_ENABLED=0 go test ./...`) lief grün für alle Pakete. `make gates` lief vollständig grün (baseline-verify, workflow-pins, doc-check `863 Datei(en), 0 Befund(e)` ×2, lint, test, arch-check, `coverage-gate: 94.30% erfüllt Schwelle 93%`, semgrep `55 rules, 0 findings`, gate-consistency, planning-check). Zusätzlich gezielt `go test -v -run 'TestModulRegistrySpiegel|TestModuleMirrorIssues_Guards|TestFocusDisable'` gegen beide betroffenen Pakete gefahren: alle sechs Testfunktionen (`TestModuleMirrorIssues_Guards`, `TestModulRegistrySpiegel_LastenheftCLI002`, `TestModulRegistrySpiegel_LastenheftGlossar`, `TestModulRegistrySpiegel_OperationsMdOptionen`, `TestFocusDisable_DecktDCheckYmlModules`, `TestFocusDisableIssues_Guards`) laufen tatsächlich (`=== RUN` + `--- PASS` je Test) — kein Build-Tag, kein `t.Skip`, keine falsche Paket-Zuordnung; beide Dateien tragen kein `//go:build`. |
| 7. Kommentar-Disziplin (`AGENTS.md` §3.7) | F-1 ist ein Verstoß. Die übrigen neuen Kommentare wurden einzeln geprüft: der Datei-Kopf-Kommentar in `registry_mirror_test.go:1-9` (Provenienz `slice-238 (BEO-ALL/…, 3. Evidenz)` als einzelnes Herkunftsfeld + Abgrenzungs-Aussage „nicht die Bereichskürzel-Liste … und nicht die 'fixe Standard-Modulset'-Stellen") ist für sich genommen zulässig (Abgrenzungs-Klasse + ein Herkunftsfeld) — der Verstoß liegt ausschließlich in der zusätzlichen Nennung von `F-2/F-3/F-4` zwei Sätze weiter (F-1). Die Doc-Kommentare zu `moduleTokensIn`, `assertModuleMirrorComplete`, `readRepoFile`, `focusDisableIssues` (ohne die F-1-Stelle) und die drei `TestModulRegistrySpiegel_*`-Funktionskommentare tragen jeweils Zusage/Kopplung-Klasse (was die Funktion tut, warum sie eine reine Funktion ist) ohne Review-Historie. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Review-Historie/Review-Befund-Marker im
Code-Kommentar (Neuzugang) · Spiegel unvollständig — plan-eigenes Kriterium
nicht auf alle benannten Fundorte angewendet.

## Verdikt

**Merge-blockierend: ja.**

F-1 (HIGH) ist ein direkter Verstoß gegen eine Hard Rule (`AGENTS.md` §3.7)
in neu geschriebenem Code dieses Slice — unabhängig davon, dass dasselbe
Muster im Bestand verbreitet ist, gilt für Neuzugänge keine Ausnahme
(„Neuzugänge fallen überall unter den Anker"). F-2 (MEDIUM) ist inhaltlich
das Gegenteil dessen, was dieser Slice beheben soll: Er wurde ausgelöst,
weil ein Modul-Spiegel nach dem anderen unbemerkt veraltet, und lässt nun
selbst die beiden vollständigsten, plan-eigens benannten Spiegel
(README.md/README.de.md, benutzerhandbuch.md §6) ungedeckt zurück — ein
plausibler Kandidat für die vierte Evidenz derselben Beobachtung. Beide
Findings sind vor Closure zu beheben, nicht nur im Beobachtungs-Register zu
vermerken: F-1, weil die Hard Rule keine Abwägung vorsieht, F-2, weil der
Slice sonst mit einer bekannten, benennbaren Lücke genau der Klasse
schließt, die er selbst schließen soll.
