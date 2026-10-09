# Review R1 — slice-272: Unter `match: name` deckt ein Report nur den längsten passenden Slice

**Review-Art:** Code (Diff gegen Plan, ADRs, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-272, Commit `e5f56d02` (`git show HEAD`)
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-272 (§1 Ziel und Abgrenzung, §6 Risiko, §8 Spiegel-Liste);
eingehender Befund `docs/plan/cr/2026-10-09-befund-ai-harness-course-reviews-gleichgewicht.md`,
Punkt 2 samt Entscheidung („längster Name gewinnt");
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
samt `.a`-Algorithmus, §2-Schema `reviews.match`, SPEC-081;
[ADR-0081](../plan/adr/0081-reviews-modul.md) (Entscheidung 3),
[ADR-0105](../plan/adr/0105-reviews-liest-done-unterverzeichnisse.md),
[ADR-0106](../plan/adr/0106-skip-allows-empty-erklaert-den-ruhezustand.md);
[`MR-025`](../../harness/conventions.md#mr-025); `AGENTS.md` §3.6, §3.7, §3.8, §6 Schritt 4.
Vorherige Findings am selben Modul: slice-265 R1 F-10 (Präfix-Deckung, damals
als Grenze benannt), die Reviews zu slice-271.

## Messungen (Kommando und Ergebnis)

Alle im Wegwerf-Baum bzw. -Klon unter dem Scratchpad, nicht im Arbeitsbaum.
Fixtures: `done/` + `reviews/`, `.d-check.yml` mit `modules: [reviews]`,
`match: name`, DoD-Punkt „- [x] Review durchgeführt";
`docker run --rm --network none -v <fx>:/repo:ro -w /repo d-check:latest` (HEAD).

| # | Fall | Ergebnis |
|---|---|---|
| A | `slice-cache`, `slice-cachex`; Reports `…-slice-cachex.md` und `…-slice-cachex-und-slice-cache.md` | 1 Befund: `slice-cache` ungedeckt — `slice-cachex` zählt als längerer Name |
| B | wie A, aber eigener Report `…-slice-cache.md` statt des gemeinsamen | 0 Befunde |
| C | `slice-cache`, `slice-cache-warmup` in `exempt-paths`; Report nur für warmup | 1 Befund (`slice-cache`) — wie zugesagt |
| D | `recursive`, `skip-pattern`; Stub `done/welle-1/slice-cache-warmup.md`; Report warmup | 1 Befund (`slice-cache`) — wie zugesagt |
| E | wie D ohne `recursive` | 0 Befunde — die benannte Grenze |
| F | `slice-é`, `slice-é-ü`; Report `…-slice-é-ü.md` | 1 Befund (`slice-é`) |
| G | `slice-cache`, `slice-Cache-Warmup`; Reports für beide in ihrer Schreibung | 0 Befunde — groß/klein getrennt, wie vorher |
| H | `slice-a`, `slice-b-slice-a` (kürzerer in der Mitte); Report `…-slice-b-slice-a.md` | 1 Befund (`slice-a`) |
| I | `slice-c`, `slice-c-x`, `slice-c-x-y`; Report nur für `slice-c-x-y` | 2 Befunde (`slice-c`, `slice-c-x`) |
| M1 | Klon, Mutation: `w.names = append(…)` **hinter** den `exempt-paths`-Abzug verschoben; `make test IMAGE=dcheck-rvw272` | **grün** — kein Test bemerkt es |
| M2 | Klon, Mutation: `longerNames` durch `nil` ersetzt (Vorzustand) | rot: `LaengsterNameGewinnt`, `StubDesLaengerenZaehlt`, je mit `got []` — richtiger Grund |
| M3 | Klon, Mutation: `append` hinter den `skip-pattern`-Abzug verschoben | rot: `StubDesLaengerenZaehlt` — richtiger Grund |

Das Mutations-Image wurde danach entfernt.

## Findings

### F-1 — Die Hälfte „vor dem Abzug von `exempt-paths`" ist ungetestet

- **kategorie:** MEDIUM
- **quelle:** Skill-Prüffrage 13 (Negativtest zu neuem öffentlichem Vertrag); `DC-FA-RVW-001`
- **pfad:** `internal/hexagon/core/rules/reviews_laengster_name_test.go` · „Ein archivierter Stub des längeren Namens zählt mit"
- **befund:** Lastenheft und Spezifikation sagen zu, dass die Basisnamen **vor beiden** Abzügen gezählt werden; die Tests decken nur `skip-pattern` (M3 rot) und nur den flachen Stub. Wird das Sammeln hinter `exempt-paths` verschoben (M1), bleibt `make test` grün, und ein ausgenommener längerer Slice deckt mit seinem Report den kürzeren wieder still (Fixture C liefe dann auf 0 Befunde). Ebenso ungetestet: der Stub unter einem Unterverzeichnis mit `recursive: true` (Fixture D) — der reale Archiv-Fall.
- **verifizierbar:** ja — Mutation M1 plus `make test`.
- **klasse:** vertrag-teilweise-ungetestet

### F-2 — Der Kommentar der Gegenprobe behauptet etwas, das der Code nicht tut

- **kategorie:** MEDIUM
- **quelle:** Skill-Prüffrage 20 (Begründung trifft am Gegenstand nicht zu); `AGENTS.md` §5 Regel 13
- **pfad:** `internal/hexagon/core/rules/reviews_laengster_name_test.go` · „slice-cachex ist kein längerer Name im Sinn der Zuordnung"
- **befund:** `longerNames` prüft `strings.Contains` ohne Wortgrenze; `slice-cachex` **ist** damit ein längerer Name von `slice-cache`. Ein gemeinsamer Report `…-slice-cachex-und-slice-cache.md` deckt `slice-cache` deshalb nicht (Fixture A, 1 Befund). Der Test besteht nur, weil sein Report `slice-cache` gar nicht trägt — er läuft auch im Vorzustand grün und prüft die behauptete Eigenschaft nicht. Das Verhalten selbst ist von der Grenze in Lastenheft („deckt nur den längeren") und Spezifikation gedeckt; falsch ist die Begründung daneben, und wer ihr glaubt, nimmt eine Wortgrenze in der Längenprüfung an, die es nicht gibt.
- **verifizierbar:** ja — Fixture A.
- **klasse:** begruendung-trifft-gegenstand-nicht

### F-3 — Herkunfts-Prosa in einem Test-Kommentar

- **kategorie:** HIGH (Skill-Anker §3.7; Fix eine Zeile)
- **quelle:** `AGENTS.md` §3.7; Skill-Prüffrage 6
- **pfad:** `internal/hexagon/core/rules/reviews_laengster_name_test.go` · „// Der Fall des Befunds: zwei Slices, ein Report für den längeren."
- **befund:** „Der Fall des Befunds" verweist auf die Herkunft des Tests (den eingehenden Befund) als Prosa, nicht als auflösbares Feld; ein Leser in einem Jahr kann nicht sagen, welcher Befund gemeint ist. Der Rest der Zeile trägt die Zusage und genügt allein. Der Datei-Kopfkommentar trägt `DC-FA-RVW-001` als Feld und ist in Ordnung.
- **verifizierbar:** nein — kein Gate prüft die Kommentar-Klasse.
- **klasse:** kommentar-herkunfts-prosa

### F-4 — Spiegel aus der eigenen §8-Liste nicht nachgezogen: sf-connector-CR

- **kategorie:** LOW
- **quelle:** `MR-025`; Plan slice-272 §8 („Die Präfix-Grenze steht in … Befund-Datei und im sf-connector-CR")
- **pfad:** `docs/plan/cr/2026-10-09-cr-eingehend-sf-connector-reviews-zusage.md` · „**Grenze:** ein\n   Basisname, der über Bindestrich, Unterstrich oder Punkt Präfix eines"
- **befund:** Die Entscheidung zu Punkt 2 nennt die Präfix-Deckung weiter als Grenze, und der Abschnitt „Für die Abnahme in sf-connector" liest sich als aktuelle Anleitung an den Adopter; der Kopf `**Stand:**` trägt keinen Hinweis, dass die Grenze entfallen ist. Der Plan führt die Datei selbst als Spiegel.
- **verifizierbar:** nein.
- **klasse:** spiegel-nicht-nachgezogen

### F-5 — Mitnahme außerhalb der Abgrenzung: Historie-Zeile zu slice-271

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §6 Schritt 4; Plan slice-272 §1 („Die Leere-Regel nach `skip-pattern` — slice-271")
- **pfad:** `spec/spezifikation.md` · „Nachzug nach Review an §[`DC-FA-PLAN-001.a`]"
- **befund:** Der Commit trägt eine Historie-Zeile zum `skip-allows-empty`-Nachzug aus slice-271 nach — Gegenstand, den §1 ausdrücklich ausschließt. Die Botschaft nennt es, der Plan nicht; nach §6 Schritt 4 ist das eine Plan-Änderung, die vor den Code gehört. Inhaltlich harmlos (fehlende Zeile).
- **verifizierbar:** nein.
- **klasse:** abgrenzung-still-ausgeweitet

### F-6 — Zustände, die die Closure noch nachziehen muss

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `docs/plan/cr/2026-10-09-befund-ai-harness-course-reviews-gleichgewicht.md` · „Punkt 2 geplant"; `docs/plan/planning/observations/BEO-ALL/name-zuordnung-deckt-praefix-slices/state.md` · „offen — 1×"
- **befund:** Der Befund-Kopf sagt „geplant", die Beobachtung beschreibt ein Verhalten, das es mit diesem Commit nur noch für längere Namen **außerhalb** von `done-dir` gibt. Beides ist DoD-/Closure-Gegenstand, kein Code-Mangel; notiert, damit die Ausgang-Zuweisung der Beobachtung die verbleibende Restgrenze (Fixture E, längerer Slice nicht in `done/`) trägt statt sie als erledigt zu streichen.
- **verifizierbar:** nein.
- **klasse:** closure-zustand-offen

## Negativbefunde

- **Zuordnung, Fälle aus dem Auftrag:** mehrere längere Namen (I), kürzerer in der Mitte (H), Report mit beiden Namen (A, benannte Grenze), Unicode-Grenze (F), längerer Name nur im nicht gelesenen Unterverzeichnis (E, benannte Grenze), Groß-/Kleinschreibung (G, unverändert case-sensitiv auf beiden Seiten) — geprüft, Verhalten entspricht Lastenheft und Spezifikation; kein Fall, in dem der Code etwas anderes tut als die Spec sagt.
- **Richtung der Änderung / §3.6:** geprüft. Die neue Deckung ist eine Teilmenge der alten (`nameInReport(key)` **und** kein längerer Name) — keine Lockerung, kein Fall wird neu gedeckt. ADR-0081 Entscheidung 3 regelt nur den `slice-<NNN>`-Abgleich unter `match: id`, der unverändert ist; `match: name` hat keine eigene ADR, sein Vertrag lebt im Lastenheft. Kein neuer Grund-Code, kein neuer Scope — keine ADR-Pflicht. ADR-0105/0106: unberührt (`recursive`, `skip-pattern`, `skip-allows-empty` behalten ihre Semantik; die Namenssammlung liest nur, was der Walk ohnehin listet).
- **§3.8 (Eingaben, die nicht gescannt werden):** geprüft. Neu gelesen werden die **Namen** ausgenommener und übersprungener Dateien, nicht ihr Inhalt; die Spezifikation sagt das zu und benennt die Restgrenze (längeres Gegenstück außerhalb `done-dir`).
- **Spiegel nach MR-025:** Lastenheft (Text, Kriterium, Historie), Spezifikation (Schritt 4, Grenze, §2-Schema, Historie), SPEC-081 (sagt nur „nach `reviews.match`", nichts Veraltetes), Modul- und Funktionskommentare in `reviews.go`, `--print-config`-Vorlage, `harness/sensors/review-coverage.md` (beschreibt `match: id` dieses Repos, Grenze 5 betrifft nur `id` — nicht veraltet) — ohne Befund außer F-4/F-6.
- **Tests, richtiger Grund:** M2 und M3 rot mit der behaupteten Meldung; nur M1 überlebt (F-1).
- **Hexagon-Richtung, Netz, Suppressions:** keine neuen Imports, kein `//nolint`, kein Netz.
- **Übrige Kommentare §3.7:** `reviews.go` (`GRENZE:`-Absatz, `nameInReport`, `longerNames`, `reviewCandidates`), `config_template.go`, `reviews_muster_test.go` — tragen Zusage/Grenze, ohne Befund.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 1 | 2 | 2 | 1 |

Wiederkehrende Klasse: `begruendung-trifft-gegenstand-nicht` (Prüffrage 20).

## Verdikt

**Nachbessern vor Closure.** F-3 ist eine Zeile, F-1 und F-2 sind Test-Arbeit
ohne Änderung am Produkt-Code; das Verhalten selbst entspricht in allen
gemessenen Fällen der nachgezogenen Spec.
