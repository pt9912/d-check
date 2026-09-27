# Review-Report: slice-238 (Runde 2) — 2026-09-27

**Review-Art:** Code — Verifikation eines Fix-Commits gegen die Findings von
Runde 1 (`docs/reviews/2026-09-27-slice-238-modul-registrierungs-spiegel-review-r1.md`,
F-1 HIGH/F-2 MEDIUM), gegen `AGENTS.md` §3.7 (Kommentar-Disziplin), §3.8
(Modul-Scan-Grenze), §5 (Grenzen-Liste/Messmethode) und die
Reviewer-Skill-Prüffragen (Modul 10). Unabhängiger Kontext — kein Self-Review,
kein Rückgriff auf die Einschätzungen von Runde 1 ohne eigene Nachprüfung.

**Gegenstand:** slice-238 — Fix-Commit `ef7bb94c` (Vorgänger `688f2481` ist nur
der R1-Report-Commit; Basis `7702c593`, der von R1 geprüfte Implementierungs-
Commit). Geänderte Dateien: `internal/hexagon/core/app/registry_mirror_test.go`,
`internal/adapter/driven/configyaml/gate_consistency_test.go`. Arbeitsbaum
gleich HEAD (`ef7bb94c`), sauber vor und nach allen eigenen Proben.

**Skill:** `.harness/skills/reviewer.md` @ Version 1.16.0
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-27

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- R1-Report `docs/reviews/2026-09-27-slice-238-modul-registrierungs-spiegel-review-r1.md`
  (F-1 HIGH, F-2 MEDIUM — die beiden zu verifizierenden Behebungen)
- Slice-Plan `slice-238-modul-registrierungs-spiegel-checkliste.md` (§1
  Ziel/Abgrenzung, §2 DoD-Fundort-Liste)
- `AGENTS.md` §3.7 (fünf Kommentar-Klassen, Herkunft als **ein** auflösbares
  Feld, Bestandsgrenze/Neuzugangs-Pflicht), §3.8, §5
- `spec/lastenheft.md` `DC-FA-DIST-001` (Out-of-Scope der Docker-Hub-
  Beschreibungsseite)
- `internal/hexagon/core/model/config.go` (`validModules()`, 24 Namen) als
  Ground Truth
- `docs/reviews/2026-09-27-slice-236-file-modul-review-r1.md` (Vorgänger-
  Review, hält den `packaging/dockerhub/overview.md`-Rückstand bereits fest)
- Vollständiger Neu-Lesedurchgang beider Testdateien (nicht nur der
  geänderten Zeilen), `git blame` auf jede Fundstelle mit Review-/Slice-Marker
  zur Neuzugang-vs-Bestand-Unterscheidung
- Unabhängig nachgefahrene Gate-Läufe: `make test`, `make gates` (vollständig)
- Isolierte `go test -v -run` gegen alle neun betroffenen Testfunktionen
  (Docker/`golang:1.27.1`, kein Host-Go)
- Zwei selbst gebaute Mutationsproben (README.md, benutzerhandbuch.md),
  Arbeitsbaum nach jeder Probe zurückgesetzt und per `git status` verifiziert
- Eigene repo-weite Grep-/Fenster-Suche nach weiteren, noch ungedeckten
  Modul-Spiegeln (unabhängig von R1s eigener Suche)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Zwei **neue** Kommentare (aus Commit `7702c593`, von `ef7bb94c` nicht angefasst) bündeln drei Angaben in einer Klammer — `slice-238`, den Beobachtungs-Pfad `BEO-ALL/modulliste-spiegel-ungegated` und die Ordinalzahl „3. Evidenz". Eine bloße `BEO-ALL/<slug>`-Nennung ist in diesem Repo breit precedented und als **ein** auflösbares Feld etabliert (54 Fundstellen repo-weit, u. a. `internal/adapter/driving/cli/cli.go:137`, `internal/hexagon/core/rules/structure.go:359`); die zusätzliche Bündelung mit einer Slice-ID **und** einem Evidenz-Zähler geht über „ein Feld" hinaus und trägt ein Element, das veralten kann (steigt die Beobachtung auf eine vierte/fünfte Evidenz, bleibt „3." stehen). Es ist **keine** Review-Runden-/Befund-Marker-Nennung (kein `F-<n>`, kein `R<n>`) und damit nicht dieselbe Klasse wie die drei von F-1/R1 behobenen Stellen — R1 hat exakt denselben Wortlaut an der (inzwischen umformulierten) Kopf-Stelle von `registry_mirror_test.go` bereits geprüft und als „für sich genommen zulässig (Abgrenzungs-Klasse + ein Herkunftsfeld)" bewertet. Der Fix-Commit hat diese Formulierung an der Kopf-Stelle vollständig entfernt, an den zwei baugleichen Stellen in `gate_consistency_test.go` aber unverändert stehen lassen — eine unbeabsichtigte Inkonsistenz innerhalb desselben Fix-Commits, kein neuer Verstoß der F-1-Klasse. | `AGENTS.md` §3.7 (Herkunft als ein Feld); Reviewer-Skill Frage 6 | `internal/adapter/driven/configyaml/gate_consistency_test.go:56` (Kommentar in `assertNetlessModules`), `:129-131` (Kommentar über `TestFocusDisable_DecktDCheckYmlModules`) | ja — `git blame` zeigt beide Zeilen aus Commit `7702c593`, unverändert durch `ef7bb94c`; Wortlaut-Vergleich mit der entfernten Kopf-Stelle in `registry_mirror_test.go` (Commit-Diff `ef7bb94c`) | Herkunfts-Feld mit zusätzlichem, potenziell veraltendem Ordinal-Element (Bündelung über „ein Feld" hinaus, keine Review-Befund-Marker-Klasse) |

## Negativbefunde

| Prüfpunkt | Ergebnis |
|---|---|
| 1. Sind F-1 (HIGH, drei Review-Befund-Marker) und F-2 (MEDIUM, README/README.de/Handbuch ungedeckt) wirklich behoben? | **Ja, beide.** F-1: alle drei von R1 zeilengenau benannten Stellen wurden geprüft — `registry_mirror_test.go:17` (`F-2/F-3/F-4`-Referenz) ist aus dem komplett umgeschriebenen Kopf-Kommentar entfernt (heute Zeilen 12-24, keine Review-Runden-/Befund-Referenz mehr); `registry_mirror_test.go:39` (`slice-057-R3-Lehre`) ist zu „…treffen denselben Code." gekürzt (heute Zeile 43); `gate_consistency_test.go:171` (`slice-057-R3-Lehre` über `focusDisableIssues`) ist identisch gekürzt (heute Zeile 170). Kein Rest der drei markierten Wortlaute im aktuellen Stand (`grep -n "F-2/F-3/F-4\|slice-057-R3-Lehre"` liefert für diese drei Stellen keinen Treffer mehr). F-2: drei neue Testfunktionen `TestModulRegistrySpiegel_ReadmeEn`, `TestModulRegistrySpiegel_ReadmeDe`, `TestModulRegistrySpiegel_BenutzerhandbuchTabelle` existieren, referenzieren die drei von F-2 benannten Dateien über `readRepoFile`, und laufen tatsächlich (isoliert per `go test -v -run` bestätigt, alle PASS). |
| 2. Sind die drei neuen Regex-Muster robust — adversarielle Gegenproben? | Geprüft, ohne Befund über F-1 hinaus. `bulletModuleLineRE` (`^- `wort` — `) liefert in README.md **und** README.de.md je exakt 24 Treffer (`grep -noE "^- \`[a-z]+\` — "`), alle 24 sind reale Modulnamen aus `validModules()`; eine Gegenprobe `grep -c "^- \`" README.md` (ohne die Anforderung des Gedankenstrichs) liefert ebenfalls 24 — es gibt in beiden Dateien **keine** andere Backtick-Bullet-Zeile, die die engere Form verfehlen und damit einen stillen Blindspot erzeugen könnte. `tableModuleRowRE` (`^\| `wort``) liefert in `docs/user/benutzerhandbuch.md` exakt 24 Treffer, alle in §6 (Zeilen 2616-2639), keine einzige Fehlplatzierung außerhalb dieses Abschnitts im gesamten Dokument. Die von den Testkommentaren behauptete Aussage „Trefferzahl entspricht genau der Modul-Anzahl" wurde damit **nicht** aus dem Kommentar übernommen, sondern durch eigenes `grep`-Nachzählen bestätigt. |
| 3. Sind alle drei von F-1 geflaggten Stellen bereinigt, und wurde keine vierte, vorher unbemerkte Neuzugangs-Stelle übersehen? | Ja, bereinigt (siehe Prüfpunkt 1) und keine vierte übersehen. Vollständiger `grep -n "slice-\|R[0-9]-[A-Z]\|F-[0-9]\|H[0-9]-"` über beide Dateien, jeder Treffer einzeln per `git blame` auf Neuzugang vs. Bestand geprüft: **zwei** Treffer sind waschechte, aber **vorbestehende** (Commit `ea857a687`, 2026-07-06, von slice-238 nicht angefasst) Review-Befund-Marker derselben verbotenen Klasse — `gate_consistency_test.go:20` „ADR-0032 R1-F-6" und `:39` „slice-057-R3-Lehre" (im **unveränderten** `assertNetlessModules`-Kommentar, zu unterscheiden vom oben bereits bereinigten `focusDisableIssues`-Kommentar). Beide sind nach der Bestandsgrenze aus `AGENTS.md` §3.7 grandfathered („geräumt wird beim nächsten Anfassen der Zeile") — sie wurden von slice-238 nicht editiert und sind daher nicht Gegenstand dieses Slice; siehe Negativbefund unten. F-1 selbst (dieses Reports) ist der einzige neue Fund, kategorisch verschieden (kein Review-Runden-Bezug). |
| 4. Bleibt nach den drei neuen Tests noch ein prominenter, wörtlich-vollständiger Modul-Spiegel ungedeckt? | Eigene, von R1 unabhängige Fenster-Suche über den gesamten Tracked-Bestand (Backtick-Token-Dichte in gleitenden 40- und 250-Zeilen-Fenstern) findet **keinen neuen** ungedeckten Kandidaten. Alle Dateien mit hoher Trefferdichte sind entweder bereits gedeckt (`spec/lastenheft.md`, `README.md`, `README.de.md`, `docs/user/operations.md`, `docs/user/benutzerhandbuch.md`), Streufunde ohne kohärente Liste (`spec/spezifikation.md` — Modulnamen als Config-Schlüssel-Präfixe verstreut über 4000 Zeilen, kein einziger zusammenhängender 24er-Block; `docs/plan/adr/README.md` — ein Modulname pro ADR-Zeile, kein Registrierungs-Anspruch; `CHANGELOG.md` — chronologisch über die gesamte Historie verstreut) oder eine historisch eingefrorene ADR-Aussage (`docs/plan/adr/0033…`, Stand „17 Module" zum Entscheidungszeitpunkt, per §3.5 immutabel und nicht nachzuziehen). Der einzige echte Treffer, `packaging/dockerhub/overview.md` (aktuell „21 rule modules", tatsächlich 24 — `reviews`/`mentions`/`file` fehlen), ist **kein neuer Fund**: siehe eigenen Negativbefund unten. |
| 5. `make gates`/`make test` unabhängig grün, neue Tests laufen wirklich, werden bei Mutation rot? | Ja. `make test` (Docker-Build) grün für alle Pakete. `make gates` vollständig grün (baseline-verify, workflow-pins, doc-check `864 Datei(en), 0 Befund(e)` ×2, lint, test, `coverage-gate: 94.30% erfüllt Schwelle 93%`, semgrep `55 rules, 0 findings`, gate-consistency, planning-check). Isoliert per `go test -v -run` (Docker/`golang:1.27.1`) alle neun Testfunktionen einzeln bestätigt als tatsächlich laufend: `TestModulRegistrySpiegel_LastenheftCLI002`, `_LastenheftGlossar`, `_OperationsMdOptionen`, `_ReadmeEn`, `_ReadmeDe`, `_BenutzerhandbuchTabelle`, `TestModuleMirrorIssues_Guards`, `TestFocusDisable_DecktDCheckYmlModules`, `TestFocusDisableIssues_Guards`, `TestQA03_NetlessModuleList_Guards` (inkl. des Subtests `unbekanntes_Modul_gesetzt_(slice-238)`) — alle PASS, keine `t.Skip`, kein Build-Tag. Zwei eigene Mutationsproben: (a) `` `file` `` → `` `filex` `` in README.md → `TestModulRegistrySpiegel_ReadmeEn` bricht korrekt mit „Modul \"file\" fehlt" + „Eintrag \"filex\" ist kein bekanntes Modul"; (b) `` `links` `` → `` `linksx` `` in der ersten Handbuch-Tabellenzeile → `TestModulRegistrySpiegel_BenutzerhandbuchTabelle` bricht korrekt mit derselben Fehlermodus-Paarung. Arbeitsbaum nach jeder Probe per `cp`-Backup zurückgesetzt und `git status --short` als sauber verifiziert. |
| 6. Kommentar-Disziplin — sonstige neue Kommentare im Fix-Diff (`ef7bb94c` selbst) | Geprüft, ohne Befund über F-1 hinaus. Der komplett umgeschriebene Kopf-Kommentar in `registry_mirror_test.go:12-24` beschreibt jetzt zwei Erkennungsformen und ihre Abgrenzung (Zusage-/Grenze-Klasse) ohne jede Review-/Slice-Referenz. Die beiden neuen Doc-Kommentare zu `bulletModuleLineRE`/`tableModuleRowRE`/`moduleLinesIn` tragen Zusage-/Grenze-Klasse (was erkannt wird, warum keine Kollision entsteht) ohne Chronik. |
| 7. Commit-Botschaft von `ef7bb94c` — Übertreibung ggü. gemessener Arbeit (`AGENTS.md` §5)? | Geprüft, ohne Befund. Die Botschaft benennt die drei behobenen Fundorte konkret und verallgemeinert nicht („…waren…ungeprüft — der wahrscheinlichste nächste Fundort…"), ohne einen Vollständigkeitsanspruch über die drei hinaus zu erheben. |

**Zusätzlicher, eigenständig recherchierter Negativbefund (nicht Teil der
sechs Auftragspunkte, aber relevant für die Vollständigkeits-Frage):**
`packaging/dockerhub/overview.md` §Modules zählt aktuell wörtlich „21 rule
modules" und listet 21 Namen (`reviews`/`mentions`/`file` fehlen gegenüber den
24 aus `model.ValidModules()`). Das ist ein reales, gegenwärtig falsches
Zahlen-Versprechen — aber **kein neuer Fund**: Der Vorgänger-Review
`docs/reviews/2026-09-27-slice-236-file-modul-review-r1.md` (Negativbefund 8)
hat diesen exakten Rückstand bereits vor slice-238 benannt und ausdrücklich
als „außerhalb dessen Abgrenzung" eingeordnet — mit Beleg: `spec/lastenheft.md`
`DC-FA-DIST-001` Out-of-Scope sagt wörtlich, der **Inhalt** der Hub-
Beschreibungsseite „ist aber nicht Teil der Distributions-Zusage: sein
Fehlschlag lässt das Release grün". Das ist eine Entscheidung auf
Lastenheft-Rang (Rang 1 der Source Precedence), nicht ein Versäumnis dieses
oder eines vorherigen Slice — eine mechanische Deckung dieser Datei wäre eine
Anforderungsänderung, kein Bugfix. Aufgeführt hier als Negativbefund, nicht
als Finding.

**Zusätzlicher Negativbefund — pre-existing Review-Befund-Marker (Bestand,
außerhalb des Diffs):** `gate_consistency_test.go:20` („ADR-0032 R1-F-6") und
`:39` („slice-057-R3-Lehre", im unveränderten `assertNetlessModules`-
Kommentar) sind echte Instanzen derselben verbotenen Klasse wie das
ursprüngliche F-1 — beide aus Commit `ea857a687` (2026-07-06), von keinem der
beiden slice-238-Commits berührt. Nach `AGENTS.md` §3.7 Bestandsgrenze
grandfathered; nicht Gegenstand dieses Reviews, aber notiert, falls das
Steering-Register künftig eine eigene Beobachtung für Review-Befund-Marker-
Bestand führen will (andere Beobachtung als `BEO-ALL/modulliste-spiegel-
ungegated`).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Herkunfts-Feld mit zusätzlichem,
potenziell veraltendem Ordinal-Element (Bündelung über „ein Feld" hinaus).

## Verdikt

**Merge-blockierend: nein.**

F-1 (HIGH) und F-2 (MEDIUM) aus Runde 1 sind beide vollständig und korrekt
behoben — durch eigenständiges Neu-Lesen beider Dateien, `git blame`-Herkunfts-
Nachweis, adversarielle Regex-Proben (24/24-Treffer-Auszählung ohne
Fehltreffer in beiden README-Varianten und im Handbuch) sowie zwei
Mutationsproben (Arbeitsbaum sauber restauriert) unabhängig bestätigt. `make
gates` und `make test` laufen grün, alle neun betroffenen Testfunktionen
laufen nachweislich (nicht nur behauptet). Der einzige neue Fund (F-1 dieses
Reports) ist LOW: eine Inkonsistenz innerhalb des Fix-Commits selbst (dieselbe
Formulierung an einer Stelle entfernt, an zwei baugleichen stehen gelassen),
kategorisch verschieden von der behobenen HIGH-Klasse, da kein Review-Runden-
/Befund-Bezug vorliegt und ein bloßer `BEO-ALL/<slug>`-Verweis in diesem Repo
breit etabliert ist. Die beiden entdeckten, aber vorbestehenden
Review-Befund-Marker sind nach der Bestandsgrenze grandfathered und nicht
diesem Slice zuzurechnen; der entdeckte Docker-Hub-Zahlen-Rückstand ist ein
bereits dokumentierter, per Lastenheft ausdrücklich ausgenommener
Alt-Befund. Nichts davon erfordert eine weitere Runde vor Closure.
