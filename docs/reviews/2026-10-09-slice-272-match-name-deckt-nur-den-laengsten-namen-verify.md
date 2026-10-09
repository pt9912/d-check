# Verifikation — slice-272: Unter `match: name` deckt ein Report nur den längsten passenden Slice

**Rolle:** Verifier (Baseline-Regelwerk `modul-11-verification.md`): DoD und Plan gegen den tatsächlichen Stand, nicht Diff gegen Plan
**Gegenstand:** slice-272 auf `b369f472` (Code- und Spec-Stand `e5f56d02`, R1-Korrekturen `b369f472`); DoD-Punkte 1 und 2
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan §1, §2, §6, §8;
[Befund von `ai-harness-course`](../plan/cr/2026-10-09-befund-ai-harness-course-reviews-gleichgewicht.md), Punkt 2 samt Entscheidung („längster Name gewinnt");
Review R1 zu slice-272 (nur als Hinweis auf Messbäume, nicht als Beleg übernommen);
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
samt `.a`-Algorithmus und §2-Schema `reviews.match`.

Alle Läufe in frischen Klonen bzw. Wegwerf-Bäumen unter dem Session-Scratchpad, nur über
`make`/Docker. Am Arbeitsbaum ist nichts geändert außer dieser Datei.

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| V1 | `git clone` des Arbeits-Repos (HEAD `b369f472`), `make gates IMAGE=dcheck-v272` | Exit 0, `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; lint `0 issues.`; `ok …/internal/hexagon/core/rules`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; semgrep `Ran 55 rules on 66 files: 0 findings.`; doc-check `0 Befund(e)`; `reviewCandidates`, `hasReviewContaining`, `nameInReport`, `longerNames` je 100,0 % |
| V2 | zweiter Klon, Mutationen (Tabelle *Bewusstes Brechen*), je `make test IMAGE=dcheck-v272-mut` | siehe unten |
| V3 | Wegwerf-Bäume, je `docker run --rm --network none -v <baum>:/repo:ro <image> [--json]`, `ghcr.io/pt9912/d-check:v0.85.0` gegen `dcheck-v272:latest` (V1), stdout, stderr und Exit-Code getrennt per `cmp` | Tabellen *Black-box* und *`match: id`* |
| V4 | `grep` nach der alten Grenzformulierung („Bindestrich, Unterstrich oder Punkt", „Präfix eines anderen", `slice-a-foo-bar`) über `spec/`, `internal/`, `docs/user/`, `docs/plan/cr/` | Treffer nur noch in eingefrorenen Artefakten (Review-Reports zu slice-265, die `observation.md` der Beobachtung) und in der sf-connector-CR-Entscheidung, die der Nachtrag aus `b369f472` ausdrücklich zurücknimmt; Lastenheft, Spezifikation, Code-Kommentar, `--print-config`-Vorlage tragen die neue Zuordnung |

Die Testbilder wurden danach entfernt.

## DoD-Punkt 1 — Lastenheft und Spezifikation

**Erfüllt.**

- Lastenheft 0.104.0, `DC-FA-RVW-001`: der Satz „deckt ein Report auch einen Slice, dessen
  Basisname über einen Bindestrich, Unterstrich oder Punkt Präfix des eigenen ist" ist ersetzt
  durch „deckt ein Report einen Slice **nicht**, wenn sein Name auch einen **längeren**
  Slice-Basisnamen trägt, der den kürzeren enthält"; gezählt über alle Slice-Dateien, die das
  Modul sieht, auch die von `exempt-paths`/`skip-pattern` ausgenommenen. Neue Grenze: ein
  Report mit beiden Namen deckt nur den längeren. Kriterium „Boundary (Slug-Kennung)" nennt
  `slice-cache`/`slice-cache-warmup`, den archivierten Stub und den Fall ohne Slice-Datei des
  längeren. Historie-Zeile 0.104.0 vorhanden.
- Spezifikation §`DC-FA-RVW-001.a` Schritt 4 (Ausnahme längerer Basisname, Menge vor beiden
  Abzügen), Grenze des Moduls (Präfix-Satz zurückgenommen; Restgrenze: längeres Gegenstück
  außerhalb der gelesenen Menge, etwa ohne `recursive`), §2-Schema `reviews.match`,
  Historie-Zeile.
- Die Zusage stimmt mit der Entscheidung zu Punkt 2 des Befunds überein (längster Name gewinnt,
  übersprungene Stubs zählen mit).

## DoD-Punkt 2 — Modul, Tests, `make gates`

**Erfüllt.** `make gates` grün (V1). Die drei verlangten Tests existieren
(`reviews_laengster_name_test.go`): `TestReviewsMatchName_LaengsterNameGewinnt` (Fall des
Befunds), `TestReviewsMatchName_StubDesLaengerenZaehlt` (archivierter längerer Stub),
`TestReviewsMatchName_KeinePraefixKollisionOhneGrenze` (Gegenprobe `slice-cachex`); dazu aus
R1 `…_ExemptDesLaengerenZaehlt` und `…_StubImUnterverzeichnis`.

### Bewusstes Brechen (V2)

| Mutation | Ergebnis von `make test` |
|---|---|
| **Vorzustand:** `internal/hexagon/core/rules/reviews.go` auf `dbcec84d` zurückgesetzt, neue Tests behalten | Exit 2; FAIL `LaengsterNameGewinnt` („slice-cache ohne eigenen Report meldet, slice-cache-warmup nicht, got []"), `StubDesLaengerenZaehlt` („…zählt als längerer Name, got []"), `ExemptDesLaengerenZaehlt`, `StubImUnterverzeichnis` (je `got []`). Richtiger Grund: `slice-cache` bleibt durch den Report von `slice-cache-warmup` gedeckt — genau das Verhalten des Befunds. Sonst kein Test rot |
| `w.names = append(…)` hinter den `exempt-paths`-Abzug | Exit 2; FAIL nur `ExemptDesLaengerenZaehlt` (`got []`) — die R1-Lücke M1 ist geschlossen |
| `w.names = append(…)` hinter beide Abzüge (vor `w.out = append`) | Exit 2; FAIL `StubDesLaengerenZaehlt`, `ExemptDesLaengerenZaehlt`, `StubImUnterverzeichnis` (je `got []`) |

`KeinePraefixKollisionOhneGrenze` und `TestReviewsMatch_NameWortgrenze` bleiben im Vorzustand
grün. Das ist richtig so: Beide halten ein Verhalten fest, das sich nicht ändern darf
(`slice-cachex` deckt `slice-cache` nicht; ohne Slice-Datei des längeren Namens deckt dessen
Report den kürzeren weiter). Sie sind Regressionswächter, keine Testbehauptung zum Fix
(siehe V-2).

### Black-box gegen `v0.85.0` (V3)

DoD-Punkt jeder Slice-Datei: `- [x] Review durchgeführt, Report unter docs/reviews/ liegt vor`;
`.d-check.yml` mit `modules: [reviews]`, `done-dir: done`, `reviews-dir: reviews`.

| Baum | `v0.85.0` | HEAD |
|---|---|---|
| **f1** `match: name`; `slice-cache`, `slice-cache-warmup`; nur Report `2026-10-09-slice-cache-warmup.md` | Exit 0, `0 Befund(e)` | Exit 1, `done/slice-cache.md:3 reviews review-missing Review-Zusage ohne Report unter reviews fuer slice-cache` |
| **f2** wie f1 plus `2026-10-09-slice-cache.md` | Exit 0 | Exit 0 — alle Kanäle identisch |
| **f3** `match: name`, `recursive`, `skip-pattern`; Stub `done/welle-1/slice-cache-warmup.md`, Report warmup | Exit 0 | Exit 1, `review-missing` für `slice-cache` |
| **f4** `slice-cache`, `slice-cachex`, Report nur `slice-cachex` | Exit 1, `review-missing` `slice-cache` | identisch |

### `match: id` byte-identisch (V3)

Drei Bäume, je Text und `--json`, stdout, stderr und Exit-Code per `cmp`:
**i1** ohne `match` (Default `id`), Slices `slice-001-a`, `slice-002-b`, `slice-003-c`,
`slice-0010-d`, `slice-cache`, `slice-cache-warmup`, Reports zu 001, 0010 und
`slice-cache-warmup`; **i2** `match: id` explizit, plus Report zu 002;
**i3** `match: id`, `recursive`, `require-promises`, `skip-pattern`, plus Stub
`done/welle-1/slice-004-x.md`. Alle sechs Läufe **IDENTISCH**: Exit 1, 4/3/4 Befunde,
darunter die Slug-Kennungen mit dem Hinweis auf `match: name`. Der `id`-Zweig von
`reviewFinding` ist im Diff unverändert; die Namenssammlung liest nur Namen, die der Walk
ohnehin listet.

## Befunde

### V-1 — Die Grenze „Report mit beiden Namen deckt nur den längeren" hat keinen Unit-Test

- **kategorie:** INFO
- **pfad:** `internal/hexagon/core/rules/reviews_laengster_name_test.go`
- **befund:** Lastenheft und Spezifikation sagen die Grenze zu; R1 hat sie black-box gemessen
  (Fixture A), kein Test hält sie fest. Nicht Teil der DoD, die drei verlangten Fälle sind
  gedeckt. Ebenso ohne Test: `longerNames` prüft `strings.Contains`, nicht das Präfix — der
  kürzere Name in der Mitte des längeren (R1 Fixture H) ist zugesagt („enthält"), aber
  ungetestet; alle Test-Fixtures tragen den kürzeren Namen als Präfix, eine Mutation auf
  `strings.HasPrefix` bliebe danach grün (aus den Fixtures gelesen, nicht gefahren).

### V-2 — Die Gegenprobe `slice-cachex` ist kein Beleg für den Fix

- **kategorie:** INFO
- **pfad:** `reviews_laengster_name_test.go` · `TestReviewsMatchName_KeinePraefixKollisionOhneGrenze`
- **befund:** Der Test läuft im Vorzustand grün; er hält die Wortgrenze von `nameInReport` fest,
  nicht die Längste-Name-Regel. Der DoD-Wortlaut („deckt weiter nicht") verlangt genau das, ein
  Regressionswächter. Für die Closure-Notiz: diesen Test nicht als Rot-Beleg des Fixes zitieren.

### V-3 — Closure-Zustände noch offen

- **kategorie:** INFO
- **befund:** Der Kopf des Befunds sagt zu Punkt 2 „geplant"; das Risiko aus §6
  (bestehende Repos werden rot) wird durch f1/f3 bestätigt: genau diese Verhaltensänderung muss
  die Release-Notiz nennen. Die Beobachtung `BEO-ALL/name-zuordnung-deckt-praefix-slices` trägt
  eine Restgrenze (längeres Gegenstück außerhalb der gelesenen Menge, R1 Fixture E). Alles
  Gegenstand von DoD-Punkt 4, nicht dieser Verifikation.

## Negativbefunde

- **Behauptung gegen Beleg:** die Commit-Botschaft von `b369f472` behauptet, das Verschieben
  hinter `exempt-paths` mache den neuen Test rot — selbst nachgefahren, bestätigt.
- **Abgrenzung §1:** `match: id` unverändert (byte-identisch gemessen); Handbuch und
  Release-Notiz nicht berührt. Die Historie-Zeile zu slice-271 in der Spezifikation hat R1 als
  F-5 bereits gemeldet; hier nicht doppelt.
- **Spiegel nach MR-025:** siehe V4; keine lebende Stelle trägt die alte Präfix-Grenze noch als
  aktuell.

## Verdikt

**DoD-Punkte 1 und 2 erfüllt.** Die Tests werden ohne den Fix aus dem richtigen Grund rot,
`make gates` ist im frischen Klon grün, black-box meldet HEAD `slice-cache` gegenüber
`v0.85.0` und ist mit eigenem Report grün, `match: id` bleibt byte-identisch. Drei
INFO-Befunde, kein Nachbesserungsbedarf vor der Closure.
