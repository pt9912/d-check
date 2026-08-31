# Review-Report — slice-174 (Register-Deckung: die vierte `planning`-Fähigkeit)

- **Review-Art:** Code (gegen Plan, ADR-0079, `DC-FA-PLAN-001`/`.a` §Register-Deckung, Hard Rules)
- **Gegenstand:** Commits `f6bdb7e..d3b0dac` (Diff `git diff f6bdb7e~1..d3b0dac`) — sechs Commits (Plan angelegt · beansprucht · Verantwortlicher gesetzt · Entscheidungen §2a · `feat` Implementierung · `docs(adr)` ADR-0079). Arbeitsbaum-Stand = HEAD (`a4bc209`); der jüngste Commit schneidet einen unabhängigen Slice (slice-187) und ist nicht Gegenstand.
- **Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) v1.13.0
- **Modell-ID:** claude-sonnet-5
- **Datum:** 2026-08-31
- **Eingangs-Kontext:** Slice-Plan [`slice-174`](../plan/planning/in-progress/slice-174-register-deckung.md); [ADR-0079](../plan/adr/0079-register-deckung-zaehlt-linktext.md); [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (Fassung 0.81.0, §Register-Deckung) und `SPEC-079` in [`spec/spezifikation.md`](../../spec/spezifikation.md); Wellen-Kontext [welle-86](../plan/planning/welle-86-closure-uebergang-durchsetzen.md); [`AGENTS.md`](../../AGENTS.md) §3 (bes. §3.7/§3.8), §5; vorherige Findings am selben Modul: [`2026-07-01-slice-057-planning-code-r2.md`](2026-07-01-slice-057-planning-code-r2.md), [`2026-07-01-slice-057-planning-doc-r1.md`](2026-07-01-slice-057-planning-doc-r1.md), [`2026-08-16-slice-102-wellen-invariante-review.md`](2026-08-16-slice-102-wellen-invariante-review.md), [`2026-08-16-slice-102-wellen-invariante-re-review.md`](2026-08-16-slice-102-wellen-invariante-re-review.md), [`2026-08-02-slice-088-planning-layer-review.md`](2026-08-02-slice-088-planning-layer-review.md); Beobachtungs-Register [`BEO-011`](../plan/planning/observations.md), [`BEO-013`](../plan/planning/observations.md).

Repo unverändert vor/nach dem Review (`git status --porcelain` leer, HEAD `a4bc209`). Alle Proben liefen gegen das lokal gebaute `d-check:latest`-Image (`make build`) oder über `make`-Gates; eine temporäre Fixture-Kopie lag außerhalb des Repos unter der Scratchpad-Wurzel. Eine einzige Mess-Mutation am Repo selbst (`.d-check.yml`: `# observations:` auskommentiert → aktiv, um die Fähigkeit gegen den echten Bestand zu fahren) wurde per `git checkout -- .d-check.yml` exakt zurückgestellt; `git status --porcelain .d-check.yml` danach leer.

---

## Eigener Lauf (Ausgabe, nicht behauptet)

| Lauf | Ausgabe |
|---|---|
| `make test` | `ok internal/hexagon/core/rules 0.032s`, `ok internal/adapter/driven/configyaml 0.014s`, alle elf Pakete grün |
| `make lint` | `golangci-lint run ./...` → `0 issues.` |
| `make arch-check` | `docker run … ghcr.io/pt9912/a-check:v0.19.0…` → `gesamt: 0 Befund(e)` |
| `make coverage-gate` | `coverage-gate: OK — Coverage 94.60% erfüllt Schwelle 93%`; je Funktion (Auszug): `CheckPlanningObservations 90.9%`, `declaredObservationIDs 100.0%`, `citedWithoutRow 93.8%`, `inCodeOnly 100.0%`, `markdownFilesUnder 64.7%`, `observationFinding 100.0%`, `applyObservations 70.0%` (Konfig-Layer) gegen `applyWaves 100.0%` (Schwester-Fähigkeit) |
| `make planning-check` | `d-check: 643 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-check` | `d-check: 643 Datei(en) geprüft, 0 Befund(e)` |
| Fähigkeit scharf gegen den echten Bestand (`.d-check.yml` temporär mit `planning.observations` aktiviert, `--enable planning`, alle anderen Module ab) | `d-check: 643 Datei(en) geprüft, 0 Befund(e)`, Exit 0 — reproduziert die im Slice/ADR behauptete Messung |
| Vier differentielle Fixture-Läufe gegen `d-check:latest` (eigenes `.d-check.yml`, `docs/plan/planning/{observations.md,done/slice-1.md}`, außerhalb des Repos) | s. F-1 |
| Drei fail-closed-Proben am Produkt (Register fehlt · `dirs`-Eintrag ist eine Datei/existiert nicht · `pattern: "BEO-["`) | s. Negativbefunde |

---

## Findings

### F-1 · HIGH · `DC-FA-PLAN-001.a` / Prüffrage 2 (Kern-Modul meldet falsche Menge) · `internal/hexagon/core/rules/planning_observations.go:162-171` (`observationFinding`) gegen `internal/hexagon/core/model/finding.go:121-150` (`SortFindings`)

**Befund:** `observationFinding` setzt `Target: file` — für **jeden** Befund derselben Datei denselben Wert, unabhängig von der zitierten Kennung. `SortFindings` dedupliziert stabil über das Tupel `(Datei, Zeile, Regel, Ziel, Grund)`, **ohne** `Message`. Stehen zwei **verschiedene**, beide unregistrierte Kennungen in derselben Zeile derselben Datei, tragen ihre zwei Befunde identisches `(File, Line, Rule="planning", Target=file, Reason="observation-unregistered")` und unterscheiden sich nur in `Message` — der zweite fällt der Deduplikation lautlos zum Opfer. Differentiell am gebauten Image belegt (Fixture außerhalb des Repos, `planning.observations.register`/`dirs` auf die Fixture gesetzt):

| Zeile enthält | Befunde |
|---|---|
| nur `` [`BEO-998`](…) `` | 1 (`BEO-998`) |
| nur `` [`BEO-999`](…) `` | 1 (`BEO-999`) |
| **beide** `` [`BEO-998`](…) `` und `` [`BEO-999`](…) `` in **einer** Zeile | **1** (nur `BEO-998`; `BEO-999` fehlt spurlos) |
| dieselben zwei Kennungen auf **zwei** Zeilen (Kontrolle) | 2 |

`ids.go` löst genau dieselbe Kollisionsgefahr für seine eigenen Befunde bereits, indem `Target` die **Kennung** trägt statt eines Datei-weiten Werts (`internal/hexagon/core/rules/ids.go:120` `Target: val[m[0]:m[1]]`, `:146` `Target: ln.Text[m[0]:m[1]]`) — die neue Fähigkeit übernimmt dieses Muster nicht. Der Exit-Code bleibt in jedem Fall 1 (kein stilles Grün), aber die **Menge** der gemeldeten Verstöße ist falsch — exakt die Situation, die die Fähigkeit laut Ziel (§1 des Slice-Plans, ADR-Kontext) verhindern soll: *„ein erfundenes `BEO-999` … fiele … niemandem auf"*. Dichte Mehrfach-Zitate in einer Zeile/Zelle sind im Beobachtungs-Register selbst der **Regelfall** (mehrere Zeilen der Haupttabelle nennen zwei oder mehr `BEO-`-Kennungen im selben Fließtext-Absatz), das Kollisionsrisiko ist also nicht hypothetisch.

**Verifizierbar:** ja — die vier Fixture-Läufe oben sind reproduzierbar; kein Gate fängt es (`planning.observations` ist im eigenen Profil noch nicht scharf, s. Negativbefunde).

**Klasse:** `dedup-kollision-generisches-target` — dieselbe Wurzel wie `N-1` im slice-102-Re-Review (dort: drei Bedeutungen von `wave-drift` kollidieren auf einem Datei-Ziel, dort MEDIUM eingestuft, weil die Kollision zwischen **verwandten** Symptomen **desselben** Konfigurationsobjekts lag und iterative Läufe konvergieren). Hier kollidieren **beliebige, unabhängige** Kennungen aus einem 780+ Vorkommen großen Baum — die höhere Trefferwahrscheinlichkeit bei der beobachteten Schreibweise des Registers trägt die Einstufung HIGH statt MEDIUM.

---

### F-2 · MEDIUM · Prüffrage 13 (fehlender Negativtest zu neuem öffentlichen Vertrag) · `internal/adapter/driven/configyaml/configyaml.go:2236-2253` (`applyObservations`)

**Befund:** Die dritte Fähigkeit (`waves`) hat ein vollständiges Config-Rand-Testpaar (`TestDecode_WavesFehler` — 13 Fälle inkl. `dir absolut`, `dir mit ..`, `glob ungültig` —, `TestDecode_WavesHappy`, `configyaml_test.go:386/411`). Für die vierte Fähigkeit (`observations`) existiert **kein** Pendant: `grep -rn "observations" internal/adapter/driven/configyaml/*_test.go` liefert null Treffer. Gemessen über `make coverage-gate`: `applyObservations` **70.0 %** und `markdownFilesUnder` **64.7 %**, gegen `applyWaves` **100.0 %**. Die im Lastenheft selbst formulierte Grenzbedingung — [„ein ungültiges `pattern` ist ein Konfigurationsfehler (Exit 2)“](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (Zeile 2462, Akzeptanzkriterium „Negative (fail-closed)“) — ist damit durch **keinen** Test gedeckt, weder am Config-Layer noch am Kern (`CheckPlanningObservations`s eigener `regexp.Compile`-Zweig ist durch die 70/90,9 %-Lücken ebenfalls unbewacht und praktisch unerreichbar, weil `applyObservations` jedes ungültige `pattern` schon vor dem Kern abweist). Am Produkt geprüft funktioniert die Zusage heute korrekt (s. Negativbefunde) — ungewächtert bleibt sie trotzdem.

**Verifizierbar:** ja — `make coverage-gate`-Ausgabe (oben) und der leere `grep`-Treffer sind reproduzierbar.

**Klasse:** `neuer-config-rand-ohne-negativtest` — dieselbe Klasse wie `N-2` im slice-102-Re-Review (dort: der W1-Config-Rand von `waves` blieb nach der ersten Heilung unbewacht). Zweites Auftreten in derselben Modul-Familie.

---

### F-3 · MEDIUM · `AGENTS.md` §3.7 (Baseline-Merksatz „ein Kommentar beschreibt, was da ist“) / Maintainability · `internal/hexagon/core/rules/planning_observations.go:108-111` (`inCodeOnly`)

**Befund:** Der Kommentar behauptet: *„Die Linktext-Ausnahme ist dieselbe, die `ids` für seine Linkpflicht trifft (`IDOccurrenceExempt`) — eine zweite Antwort auf dieselbe Frage wäre ein Defekt.“* Der Code ruft `IDOccurrenceExempt` (`internal/hexagon/core/rules/ids.go:216-226`) jedoch **nicht** auf — `inCodeOnly` reimplementiert die Linktext-Überlappungsprüfung (`start >= sp.TextStart && end <= sp.TextEnd && !sp.IsImage`) eigenständig inline. Damit entsteht real die im selben Satz benannte Gefahr: **zwei unabhängig gepflegte Antworten auf dieselbe Frage.** Ändert sich künftig, was `IDOccurrenceExempt` als Linktext zählt (z. B. eine dritte Ausnahme), folgt `inCodeOnly` dem nicht automatisch — der Kommentar verspricht eine Kopplung, die im Code nicht besteht.

**Verifizierbar:** ja — Codevergleich der beiden Funktionen; kein Gate prüft Kommentar-Wahrheitsgehalt (`dupl` im Lint-Profil greift hier nicht: `make lint` bleibt bei „0 issues“, weil die beiden Funktionen strukturell zu unähnlich für den Klon-Detektor sind).

**Klasse:** `kopplungs-kommentar-ohne-delegation`

---

### F-4 · MEDIUM · Baseline-Regelwerk Modul 2 (Bootstrap-Modus) / §8-Sub-Area-Begründung des Slice-Plans · `docs/plan/planning/in-progress/slice-174-register-deckung.md:239-244` (§8) gegen [ADR-0079](../plan/adr/0079-register-deckung-zaehlt-linktext.md) §Kontext und die Commit-Reihenfolge

**Befund:** §8 des Slice-Plans erklärt alle drei berührten Sub-Areas als reines GF (*„die Regel steht im Kanon, der Bestand ist gemessen und konform (24/24), es gibt keine Reconciliation“*, *„Evidenz-Risiko … gemessen, nicht geschätzt“*). ADR-0079 §Kontext sagt selbst: *„Beim Bau stellte sich heraus, dass die naheliegende Erkennungs-Regel dieses Repos hier das Gegenteil des Gewollten tut“* — die tragende technische Entscheidung (Prosa+Linktext zählen, reines Inline-Code nicht) wurde also **während** der Implementierung empirisch entdeckt, nicht vorab aus der Doku abgeleitet. Konsistent damit: der `feat`-Commit `c5a7569` (Code **und** Lastenheft-Update in einem Commit) liegt **vor** dem `docs(adr)`-Commit `d3b0dac`, der die Entscheidung erst formal verankert — genau umgekehrt zur in Modul 1 gezeichneten Kette Spec → ADR → Plan → Code. §8 benennt diese Korrektur-Schleife nicht; sie ist dieselbe Klasse Ereignis, die das Evidenz-/Diskrepanz-Risiko-Kriterium eigentlich einfangen soll.

**Verifizierbar:** ja — `git log --oneline f6bdb7e~1..d3b0dac` zeigt die Reihenfolge; der zitierte ADR-Satz steht wörtlich in §Kontext.

**Klasse:** `gf-deklaration-ohne-reconciliation-bei-tatsaechlicher-bf-entdeckung`

---

### F-5 · LOW · Maintainability (latente Wartungsfalle) · `internal/hexagon/core/rules/planning_observations.go:63-77` (`declaredObservationIDs`)

**Befund:** Eine Zeile deklariert nur, wenn ihre erste Zelle **exakt** und undekoriert aus der Kennung besteht (`m == cell`). Im übrigen Register ist die Backtick-Schreibweise für Kennungen die verbreitete Konvention (z. B. `` `BEO-022` ``, `` `BEO-023` `` in Prosa-Zellen); käme dieselbe Konvention künftig auf die **erste** Zelle einer neuen Zeile (`` | `BEO-025` | … ``), würde diese Zeile nicht mehr deklarieren, und jedes reguläre Zitat von `BEO-025` anderswo meldete fälschlich `observation-unregistered`. Bestand geprüft: 24 von 24 Zeilen der Haupttabelle und der Gestrichenen-Sektion führen die Kennung heute undekoriert (`grep -cP '^\|\s*BEO-\d{3}\s*\|'` → 24, `'^\|\s*`BEO-\d{3}`\s*\|'` → 0) — die Falle ist heute nicht gezündet.

**Verifizierbar:** ja — Bestandsgrep oben, Verhalten aus dem Code (`m == cell`) hergeleitet.

**Klasse:** `stille-form-abhaengigkeit-der-erst-zell-deklaration`

---

## Negativbefunde (geprüft, ohne Befund)

- **Hexagon-Importrichtung** ([ADR-0005](../plan/adr/0005-modul-layout-hexagon-ordner.md)): `make arch-check` → `gesamt: 0 Befund(e)`. Der neue Code liegt vollständig in `internal/hexagon/core/{model,rules}` und `internal/adapter/driven/configyaml`; keine neue Abhängigkeitsrichtung.
- **Netzzugriff außerhalb `external`** ([`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)): kein Netz-Code im Diff; die Fähigkeit liest ausschließlich über `driven.Filesystem`.
- **Inline-Suppression/Schwellen-Senkung ohne ADR** (`AGENTS.md` §3.2/§3.6): kein neues `//nolint`, `THRESHOLD` unverändert; `make lint` 0 issues.
- **§3.4 Abwärtsreferenzen in Spec/ADR** (Slice-/Wellen-/Commit-Bezug in `spec/lastenheft.md`, `spec/spezifikation.md` oder ADR-0079-Körper): keine gefunden; `grep -n "slice-\|welle-" docs/plan/adr/0079-*.md` liefert nichts.
- **Referenz-Richtung/Marker-Ehrlichkeit** (`DC-FA-MTX-003`): kein `d-check:status-provenance`-Marker im Diff — nicht anwendbar.
- **Fail-closed-Rand 1 (unlesbares Register)**: am Produkt geprüft — `register: docs/plan/planning/nicht-vorhanden.md` → Exit 1, Befund `Register … fehlt oder ist unlesbar (fail-closed)`.
- **Fail-closed-Rand 2 (`dirs`-Eintrag ist keine Datei/kein Verzeichnis)**: am Produkt geprüft, beide Varianten (Datei statt Verzeichnis; Pfad existiert nicht) → je Exit 1, Befund `Zitier-Verzeichnis … fehlt oder ist unlesbar (fail-closed)`. Die Grenzfall-Rechtfertigung im Kommentar (*„manche Adapter liefern eine leere Liste statt eines Fehlers“*) ist gegen beide Adapter geprüft: der echte `fs.Adapter` (`internal/adapter/driven/fs/fs.go:54-71`) lässt `os.ReadDir` selbst fehlschlagen, `coretest.MemFS.List` (`internal/hexagon/core/coretest/memfs.go:73-94`) dagegen liefert für einen fehlenden Pfad nie einen Fehler — die Kommentar-Aussage trifft exakt zu.
- **Fail-closed-Rand 3 (ungültiges `pattern`)**: am Produkt geprüft — `pattern: "BEO-["` → `d-check: error: … kein gueltiger Ausdruck …`, Exit 2. Funktioniert korrekt (Regressionslücke separat als F-2).
- **DC-QA-02 (byte-identisch ohne den Schlüssel)**: früher Return bei `o.Register == ""`; `TestObservationsFailClosedAndInert` deckt den Fall, Verhalten stimmt.
- **Deklarations-Regel, Grundfall** (Quer-Referenz im Fließtext einer Nicht-Erst-Zelle deklariert nicht): `TestObservationsOnlyFirstCellDeclares` und eigene Produktprobe stimmen überein; die Backtick-Randlage ist gesondert als F-5 (LOW) geführt.
- **Erkennungs-Form-Messung** (366/293/112/5 laut Slice/ADR): näherungsweise unabhängig nachvollzogen — grep/awk-Klassifikation über den Planning-Baum am Commit `d3b0dac` (fenced Code ausgeschlossen) ergibt 367/5/116/296 = 784 gegen die behaupteten 366+5+293+112 = 776 (≈ 1 % Abweichung, plausibel durch die markdown-bewusste Link-Span-Erkennung des Produkts gegenüber einer grep-Näherung). Die tragende **qualitative** Aussage — Linktext+Prosa dominieren gegenüber reinem Inline-Code um Größenordnungen — bleibt bei beiden Zählweisen unverändert; kein Finding.
- **`ObservationsConfig`/`rawObservations`-Kommentare** (`internal/hexagon/core/model/config.go:370-396`, `internal/adapter/driven/configyaml/configyaml.go` `rawObservations`) und der neue Grund-Code-Kommentar (`internal/hexagon/core/model/finding.go:55-60`): tragen Zusage/Abgrenzung, keine Review-Historie, keine Slice-Nummern, keine Mess-Labels.
- **`SPEC-079`-Zeile** (`spec/spezifikation.md`): Form konsistent mit den Nachbarzeilen, kein Abwärtsverweis.
- **Coverage-Schwelle**: `make coverage-gate` 94,60 % ≥ 93 % — die neuen Dateien senken den Gesamtwert nicht unter die Schwelle (die Lücken aus F-2 sind Randfälle, keine Breitenlücke).

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 1 | F-1 |
| MEDIUM | 3 | F-2, F-3, F-4 |
| LOW | 1 | F-5 |
| INFO | 0 | — |

**Finding-Klassen dieses Laufs:** `dedup-kollision-generisches-target` · `neuer-config-rand-ohne-negativtest` · `kopplungs-kommentar-ohne-delegation` · `gf-deklaration-ohne-reconciliation-bei-tatsaechlicher-bf-entdeckung` · `stille-form-abhaengigkeit-der-erst-zell-deklaration`

---

## Verdikt

**Mergeblockierend: ja — F-1 (HIGH).**

Der Kern der Fähigkeit — die Erkennungs-Form (Prosa+Linktext vs. reines Inline-Code), die Deklarations-Regel, die Scan-Mengen-Begrenzung nach `AGENTS.md` §3.8 und zwei der drei fail-closed-Ränder — trägt und ist am Produkt bestätigt; die im Slice/ADR behauptete Erkennungs-Form-Messung ist der Größenordnung nach unabhängig nachvollziehbar. Blockierend ist **F-1**: die generische `Target`-Belegung lässt den repo-weiten Dedup-Mechanismus (`SortFindings`) zwei unabhängige, in derselben Zeile zitierte unregistrierte Kennungen zu einem Befund zusammenfallen — genau die Fehlerklasse, die die Fähigkeit laut ihrem eigenen Ziel verhindern soll, entsteht dadurch innerhalb ihres eigenen Codes. Der Gate-Ausgang bleibt zwar rot (kein Silent-Grün), aber die gemeldete Menge ist falsch; da `planning.observations` im eigenen Profil noch nicht scharf ist, hat der Bestand dieses Repos die Lücke bislang nicht ausgelöst.

Die drei MEDIUM sind unabhängig von F-1 zu klären: **F-2** ist eine Regressionslücke am neu eingeführten Config-Vertrag (zweites Auftreten derselben Klasse wie `N-2` im slice-102-Re-Review — ein wiederkehrendes Muster in der `planning`-Modul-Familie, das über diesen Review-Report hinaus für die Closure-Sichtung relevant sein dürfte). **F-3** ist eine Kommentar-Zusage, die im Code nicht eingelöst wird. **F-4** ist eine Diskrepanz zwischen der behaupteten reinen GF-Konformität in §8 und der im ADR selbst dokumentierten Bau-zeitlichen Entdeckung der tragenden Entscheidung. F-5 (LOW) ist eine heute nicht gezündete, aber real angelegte Formfalle.

Dieser Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (anderer Kontext, anderes Prüf-Artefakt).
