# Review-Report — slice-079 Zitat-Verifikation (Implementierungs-Review, R1)

## Kopf-Metadaten

- **Gegenstand:** Implementierung von `slice-079` (v0.50.0-Kandidat) — die zwei
  Fähigkeiten `codepaths.check-lines` (`DC-FA-CODE-001`-Erweiterung) und das neue
  18. Modul `citations` (`DC-FA-CITE-001`). Prüftiefe: Code-Korrektheit,
  Vertragstreue, Mutations-Härte, Release-Prep-Vollständigkeit.
- **Geprüftes Artefakt:** HEAD `a387dab` (== `origin/main`), Working-Tree clean.
  Feat-Commit `bd257d1`, Realdatenbeleg `124cd52`, Release-Prep `a387dab`.
- **Betroffene Artefakte:** `internal/hexagon/core/rules/citations.go` (neu),
  `codepaths.go`, `run.go`, `model/config.go`, `configyaml/configyaml.go`,
  `app/diagnose.go`; Tests `citations_test.go`, `codepaths_checklines_test.go`;
  Spec `DC-FA-CODE-001.a`/`DC-FA-CITE-001.a`/§4/§2; Lastenheft
  `DC-FA-CODE-001`/`DC-FA-CITE-001`; ADR-0045; Release-Prep (CHANGELOG,
  version.md, READMEs, benutzerhandbuch, operations.md).
- **Reviewer-Rolle:** unabhängig, kontext-getrennt; Prozess `.harness/skills/reviewer.md`
  (Kategorien HIGH/MEDIUM/LOW/INFO, Befund-Schema, Negativbefund-Pflicht).
  Keine Behauptung der Implementierungs-Session übernommen — alles empirisch belegt.
- **Empirie-Baseline:** `make build` grün (Exit 0); `make gates` grün — 255 Dateien,
  0 Befunde (`doc-check + lint + test + arch-check + coverage-gate + semgrep +
  gate-consistency + planning-check green`). Alle folgenden Befunde sind über das
  gebaute Image `d-check:latest` bzw. `make test`-Mutation reproduziert.

---

## Findings (priorisiert)

### F-1 — `citations` **stürzt ab** (panic) bei `<von> = 0`; `codepaths` akzeptiert dieselbe Eingabe still — untere Bereichsgrenze wird nirgends validiert

- **kategorie:** MEDIUM
- **quelle:** `DC-FA-CITE-001.a` Schritt 3/4 (1-basierte Zeilennummern) ·
  Konsistenz zweier Module derselben Eingabe-Klasse (`datei:<von>-<bis>`)
- **pfad:** `internal/hexagon/core/rules/citations.go:137`
  (`strings.Join(lines[from-1:to], "\n")`); Gegenprobe
  `internal/hexagon/core/rules/codepaths.go:271-289` (`checkCodepathLineRange`)
- **befund:** Weder `citations` noch `codepaths` erzwingen `<von> ≥ 1`, obwohl die
  Spezifikation die Zeilennummern als **1-basiert** deklariert. In `citations`
  indiziert `citationSpan` die Quell-Zeilen mit `lines[from-1:to]`; bei `from = 0`
  wird das zu `lines[-1:…]` und Go **paniced** (`slice bounds out of range [-1:]`).
  Der Vergleich `from > to` (Zeile 103) fängt das nicht (`0 > to` ist falsch), und
  der Existenz-/Bereichs-Guard `to > n` (Zeile 134) greift erst hinter der
  Unterschranke. `codepaths` teilt dieselbe fehlende Validierung, indiziert aber
  nicht — dort läuft `datei:0` still als „in-range" durch (kein Befund). Dieselbe
  out-of-contract-Eingabe erzeugt so in zwei Schwester-Prüfungen zwei verschiedene
  Fehlmodi: **Absturz** (citations) vs. **stilles Grün** (codepaths).
- **failure-szenario / repro (empirisch, gebautes Image `d-check:latest`):**
  Fixture `docs/citing.md` = `<!-- d-check:cite docs/src.md:0-3 -->` gefolgt von
  einem Zitat, `docs/src.md` mit ≥ 3 Zeilen. `--enable citations` ⇒
  ```
  panic: runtime error: slice bounds out of range [-1:]
  … rules.citationSpan(…) /src/internal/hexagon/core/rules/citations.go:137
  EXIT: 2
  ```
  Ebenso mit `docs/src.md:0` (einzeln, `from=0,to=0`) und mit leerer Zieldatei.
  Gegenprobe `codepaths.check-lines: true` + `` `docs/target.md:0` `` ⇒
  `0 Befund(e)`, Exit 0 (still akzeptiert). Der Adopter tippt `:0`/`:0-N` als
  plausiblen Vertipper; statt des versprochenen kontrollierten Fail-closed
  (Exit 2 mit Klartext) bzw. eines `citation-out-of-range`-Befunds (Exit 1) liefert
  das Tool einen Go-Stacktrace, der interne `/src/…`-Pfade leakt. Der Netto-Exit
  ist zwar 2 (Gate rot, **kein** stilles Grün in citations), aber ein Absturz statt
  kontrollierter Diagnose untergräbt den Kern-Anspruch des Werkzeugs (kontrollierte,
  deterministische Befunde), und das Schwester-Modul beweist mit seiner sanften
  Behandlung die eigentlich intendierte Semantik.
- **verifizierbar:** ja — Exit-/Panic-Assertion „`d-check:cite …:0-N`" gegen
  `citations` und „`` `datei:0` ``" gegen `codepaths.check-lines`.

### F-2 — Toter Zweig `resolvable == false` in `resolveCitePath` (nie erreichbar)

- **kategorie:** LOW
- **quelle:** Maintainability
- **pfad:** `internal/hexagon/core/rules/citations.go:96-99` und `143-149`
  (`resolveCitePath` / `if !resolvable`)
- **befund:** `resolveCitePath` liefert `ok=false` nie: der `./`/`../`-Zweig gibt
  `ResolveTarget(...)` zurück, dessen dritter Rückgabewert im aktuellen Code
  **immer** `true` ist (`paths.go:16-38`), der Wurzel-relative Zweig setzt `ok`
  hart auf `true`. Der Fail-closed-Zweig `if !resolvable { return …fmt.Errorf(…) }`
  in `citationForDirective` ist damit unerreichbar. Kein Fehlverhalten heute, aber
  eine latente Falle: eine spätere Änderung, die `ResolveTarget` wirklich
  `ok=false` liefern lässt, verließe sich auf einen ungetesteten Pfad.
- **failure-szenario:** wartungslatent — der Zweig kann nie kippen, also deckt ihn
  kein Test ab; ein künftiger `ok=false`-Rückgabewert liefe ungeprüft.
- **verifizierbar:** nein (Erreichbarkeits-/Wartungs-Beobachtung, kein Gate-Lauf).

### F-3 — Fail-closed feuert auf **jede** `<!-- d-check:cite …`-Erwähnung in Prosa (Adopter-Doku, die die Direktive dokumentiert)

- **kategorie:** INFO
- **quelle:** `DC-FA-CITE-001.a` Schritt 1/2 (fail-closed) · R1-Design-Review F-6(e)
- **pfad:** `internal/hexagon/core/rules/citations.go:32,54` (`citeMarkerRe`,
  Marker-Schleife)
- **befund:** `citeMarkerRe` behandelt jede Nicht-Fence-Zeile mit `<!-- d-check:cite`
  als Direktive; eine unbrauchbare (kein Zielformat, kein folgendes Zitat) ist
  fail-closed (Exit 2). Ein Adopter, der die Direktive **in seiner eigenen Prosa
  dokumentiert** (wie d-check es in README/Handbuch/CHANGELOG tut) und `citations`
  über diese Datei aktiviert, bekäme dadurch Exit 2. Entschärft durch die
  Fence-Awareness (empirisch bestätigt: eine Direktive in einem Fence-Block wird
  **nicht** gescannt) und dadurch, dass d-checks **eigene** `.d-check.yml`
  `citations` **nicht** aktiviert (keine Selbst-Betroffenheit; `make gates` grün).
  Bewusste, dokumentierte Fail-closed-Semantik — hier nur als Betriebs-Hinweis
  festgehalten.
- **failure-szenario:** Adopter aktiviert `citations` auf einem Doku-Verzeichnis,
  das die Direktiven-Syntax in Prosa (außerhalb einer Fence) erklärt ⇒ Exit 2.
- **verifizierbar:** ja — `--enable citations` gegen eine Datei mit Prosa-`<!-- d-check:cite`
  ohne folgendes Zitat.

---

## Negativbefunde (geprüft, ohne Befund)

- **N-1 Mutations-Härte (Kern-Kanten alle gepinnt):** je EINE Zeile neutralisiert,
  `make test`, GENAU der zuständige Test kippt, zurückgedreht (Tree danach clean):
  - `strings.Contains(…)` (Teilstring-Vergleich) → `TestCitationsInlineMismatch`
  - `if checkLines` (check-lines-Gate) → `TestCodepathsCheckLines`
  - `< citationMinLen` (Mindestlänge) → `TestCitationsMindestlaenge`
  - `to > countLines` (codepaths-Zeilenzahl) → `TestCodepathsCheckLines`
  - `to > n` (citations-Span-Zeilenzahl) → `TestCitationsZitatFaeule`
  - fail-closed „kein Zitat" (`if !ok`) → `TestCitationsFailClosed`
  Kein grüner Test blieb bei einer Mutation stehen — kein ungepinnter Kern-Pfad.
- **N-2 Byte-Identität (Default aus):** `codepaths` ohne `check-lines` gibt für
  `datei:zeile`-Refs 0 Befunde (Code-Pfad `codepathValueAndRange(raw, false)` →
  unveränderte `normalizeCodepath`-Ausgabe; `TestCodepathsCheckLinesDefaultAus`;
  Dogfood-`.d-check.yml` setzt `check-lines` nicht, `make gates` grün). `citations`
  ohne Direktive 0 Befunde (`TestCitationsOhneDirektiveStill`; Realdaten-Baseline
  0/61). `DC-QA-02` gewahrt.
- **N-3 Spec/Code-Lockstep (§4 ↔ AllReasons ↔ reasonTexts):** exakt drei
  Grund-Codes (`citation-out-of-range`, `citation-inverted-range`,
  `citation-mismatch`) in `spezifikation.md` §4, `AllReasons()` und `reasonTexts()`;
  `codepaths` teilt die ersten zwei, `citation-mismatch` ist citations-eigen. Der
  gate-consistency-Test (Teil des grünen `make gates`) verriegelt den Lockstep.
- **N-4 Realdatenbeleg unabhängig nachgestellt:** Mini-Fixture, korrektes,
  re-wrapptes inline-`„…"`-Zitat einer Quell-Zeile ⇒ **0 Befunde, Exit 0**; ein
  einziges gedriftetes Wort (`vendored`→`gebuendelt`) ⇒
  `docs/src.md:2-2 citation-mismatch`, **Exit 1**. Deckt sich mit dem Beleg in
  `docs/reviews/2026-07-18-slice-079-realdatenbeleg-ai-harness-init.md`.
- **N-5 Vertragstreue der Fäule-/Fehler-Kanten:** `citations` reiht Schritt 1
  (malformter Span → fail-closed) · Schritt 2 (kein Zitat → fail-closed) · Schritt 3
  (Repo-Escape → Befund Exit 1 · `von>bis` → `citation-inverted-range` ·
  Span > Datei-Ende / Datei fehlt → `citation-out-of-range`) · Schritt 4/5
  (Normalisierung, Mindestlänge 16, Teilstring) in genau der Spec-Reihenfolge;
  Zitat-Fäule ist Befund (Exit 1), nicht fail-closed — kohärent zu
  `codepaths.check-lines` (`TestCitationsZitatFaeule`, `codepaths_checklines_test.go`).
- **N-6 R1/R2-NITs im Code eingelöst:** N-1 (Mindestlänge 16, rune-basiert:
  `len([]rune(nq)) < citationMinLen` — Unicode-Falle vermieden), N-2 (inline-Suche
  im Absatz bis Leerzeile, Zeilenüberspannung, „frühester Öffner"-Paarung `„`+`"` /
  `"`+`"`), N-3 (Leerzeile vor `>`-Block toleriert: `citationQuote` überspringt
  Leerzeilen), N-4 (Repo-Escape auf **Befund Exit 1** angeglichen, nicht Exit 2).
  Die §-2-Delimiter-Paarung ist einheitengetestet (`TestInlineQuoteSpanEinheit`).
- **N-7 Fence-/Determinismus-Eigenschaften:** `citations` scannt nur
  `proseLines` (fence-aware — Direktive in einer Code-Fence wird nicht gescannt,
  empirisch bestätigt); Befunde deterministisch (`model.SortFindings`, feste
  File/Line-Herkunft). Hermetisch — nur Filesystem-Port, kein git/Netz (`DC-QA-03`).
- **N-8 Config-Oberfläche:** `codepaths.check-lines` (bool, Default false) und
  `citations` (parameterlos, nur `scope`) sauber im strikten Decoder
  (`KnownFields(true)`); `validModules()` führt genau 18 Namen inkl. `citations`;
  `EffectiveModules` weist Unbekanntes zurück (Realdaten: `v0.49.0 --enable
  citations` bricht mit „unbekanntes Modul" ab — Feature nachweislich neu).
- **N-9 Release-Prep vollständig:** `version.md` #aktuell + `<a id="v0.50.0">`-Anker
  auf v0.50.0 verschoben; alle ghcr-`:tag`-Pins in Live-Docs = v0.50.0; README.md +
  README.de.md je 18 Modul-Bullets inkl. `citations`; `operations.md` `--enable`-Liste
  + benutzerhandbuch §6-Modultabelle enthalten `citations`; benutzerhandbuch-Header
  v0.50.0/1.39; §5 + §11 (Grund-Codes) + Versions-Historie gepflegt; CHANGELOG-Block
  `[0.50.0]`. Handbuch-YAML-Beispiel (`codepaths.check-lines: true`) ist gegen den
  configyaml-Validator zulässig (bool-Schlüssel, strikt geparst). `make gates` grün
  fängt Doku-Link-/Anker-Drift.
- **N-10 ADR-0045:** korrekt `Proposed` (nicht überschrieben; R1-Einarbeitung als
  `## Geschichte`-Anhang). Der Slice liegt noch in `in-progress/` — der ADR-Flip auf
  `Accepted` gehört an die Closure, nicht in dieses Review. Keine
  Immutabilitäts-Verletzung.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 |
| LOW | 1 | F-2 |
| INFO | 1 | F-3 |
| Negativbefund | 10 | N-1 … N-10 |

---

## Verdikt: BLOCK (eng — nur F-1)

Das Design (whitespace-normalisierter Teilstring, opt-in, fail-closed, byte-identisch)
ist am realen Modell **und** an der Mutations-Härte solide belegt: alle sechs
geprüften Kern-Kanten sind an genau ihren Test gepinnt, die Byte-Identität hält, der
Lockstep §4↔Code stimmt, die Release-Prep ist vollständig, und `make gates`/`make build`
sind grün. Die R1/R2-Design-NITs sind im ausgelieferten Code eingelöst.

Blockierend ist **ein** Punkt, eng umrissen:

- **F-1 (MEDIUM):** `citations` **paniced** (`lines[-1:]`) bei `<von> = 0` —
  empirisch mit Stacktrace + Exit 2 reproduziert; `codepaths` akzeptiert dieselbe
  out-of-contract-Eingabe still. Die fehlende Untergrenzen-Validierung eines
  1-basierten Bereichs ist ein echter Absturz-/Konsistenz-Defekt in einem gerade
  frisch ausgelieferten Kern-Modul; die Behebung ist klein (eine Untergrenzen-Prüfung
  vor der Indizierung, einheitlich für beide Zeilen-Checks) und JETZT billig,
  vor breiter Adopter-Adoption.

**Einordnung / mögliche begründete Abweichung:** Der Blast-Radius ist gering —
`citations` ist opt-in, hat heute **null** produktive Direktiven (auch beim Adopter),
und der Netto-Exit ist rot (kein stilles Grün in citations). Entscheidet der
Auftraggeber, dass ein Absturz-auf-Vertipper in einem nutzerlosen opt-in-Modul den
Release nicht aufhält, ist ein Ship-with-Fast-Follow eine **vertretbare, aber
begründungspflichtige** Abweichung, die als Risiko in den Slice gehört. F-2/F-3 sind
nicht blockierend und mit dem Fix bzw. nach Wahl nachziehbar.
