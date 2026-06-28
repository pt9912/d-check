# Review-Report — slice-052 Modul `immutable` (R1)

## Kopf-Metadaten

- **Gegenstand:** Commit `b62a520` — slice-052: neues opt-in Modul
  `immutable` (12.), Immutabilitäts-Pin gegen Core-Drift.
- **Schwerpunkt (Auftrag):** Code-Korrektheit des Moduls `immutable` +
  Spec-Treue.
- **Reviewer:** unabhängiger Reviewer (Reviewer-Skill v1.2.0).
- **Datum:** 2026-06-28.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/in-progress/slice-052-immutable-modul.md`;
  Anforderung `DC-FA-IMM-001` (`spec/lastenheft.md` §DC-FA-IMM-001,
  Zeilen 1027–1089) + Algorithmus `DC-FA-IMM-001.a`
  (`spec/spezifikation.md` Zeilen 900–938) + Grund-Code `core-drift`
  (§4, Z. 1168) + Schema-Key `immutable.exclude-sections` (§2, Z. 1123);
  ADR-0023 (Proposed); Hard Rules `AGENTS.md` §3.
- **Rolle-Abgrenzung:** kein Verifier (DoD/Gate-Läufe sind nicht meine
  Rolle; `make gates` lief laut Auftrag bereits grün), kein Stil-Polizist.
- **Betrachtete Artefakte:** `internal/hexagon/core/rules/immutable.go`,
  `…/immutable_test.go`, `…/run.go` (Wiring),
  `internal/hexagon/core/model/config.go` (`ImmutableConfig`,
  `validModules`), `…/model/finding.go` (`ReasonCoreDrift`),
  `internal/adapter/driven/configyaml/configyaml.go`
  (`rawImmutable`/`applyImmutable`/`scopeOfImmutable`); Vergleichsbasis
  `rules/pins.go` (Normalisierung, `shortHash`) und `rules/matrix.go`
  (`excludedRanges`/`inRanges`) + `rules/markdown.go`
  (`PreprocessMarkdown`/`proseLines`/`extractHeadingLines`).

---

## Findings

### MEDIUM-1 — Kein `Modul-aus`-/Wiring-Test; gesamter `Run()`-Pfad des Moduls ungetestet

- **kategorie:** MEDIUM
- **quelle:** `DC-QA-02` / `DC-FA-IMM-001` (AK „Modul-aus") /
  Maintainability (Konsistenz-Lücke zwischen Modulen derselben
  Eingabe-Klasse)
- **pfad:** `internal/hexagon/core/rules/immutable_test.go` (Testdatei als
  Ganzes); fehlende Entsprechung zu
  `internal/hexagon/core/rules/pins_test.go:118` (`TestPinsNichtAktiv`)
  und `internal/hexagon/core/rules/versions_test.go:80` (Modul-aus);
  ungetesteter Wiring-Pfad `internal/hexagon/core/rules/run.go:147–149`.
- **befund:** Die sechs AKs aus `DC-FA-IMM-001` (Lastenheft Z. 1074)
  nennen „Modul-aus" (Befundsatz byte-identisch, `DC-QA-02`); die beiden
  jüngeren opt-in-Schwestermodule belegen das je mit einem Test, der über
  `Run()` mit **inaktivem** Modul 0 Befunde auf einer Fixture nachweist,
  die bei aktivem Modul einen Befund erzeugte (`pins_test.go:118`,
  `versions_test.go:80`). Für `immutable` fehlt dieser Test. Zusätzlich
  rufen alle 7 vorhandenen `immutable`-Tests `CheckImmutable` **direkt**
  auf (`immutable_test.go:41`), sodass der Gating-/Dispatch-Pfad in
  `run.go` (`st.applies("immutable", file)` + `CheckImmutable`-Aufruf,
  Z. 147–149) sowie die `discoverScopes`/`checkFile`-Anbindung von
  **keinem** Test durchlaufen wird.
- **verifizierbar:** ja — würde man `immutable` in `defaultModules()`
  (`config.go:17`) aufnehmen oder das `applies`-Gate falsch verdrahten,
  bliebe `make test` grün (kein `immutable`-Test failt), obwohl die
  Default-off-byte-Identität (`DC-QA-02`) gebrochen wäre. Ein
  `Run()`-basierter Modul-aus-Test (wie `TestPinsNichtAktiv`) würde den
  Regress fangen.

### LOW-1 — Inertheit eines Markers in **Inline**-Code nicht direkt getestet

- **kategorie:** LOW
- **quelle:** `DC-FA-IMM-001.a` Schritt 1 / Maintainability
  (Konsistenz-Lücke zur Schwestern-Testkonvention)
- **pfad:** `internal/hexagon/core/rules/immutable_test.go:90`
  (`TestImmutableMarkerInFenceInert` deckt nur den **Fenced**-Fall).
- **befund:** `DC-FA-IMM-001.a` Schritt 1 (Spezifikation Z. 913–914)
  erklärt einen Marker „in Fenced- **oder Inline**-Code" für inert;
  getestet ist nur der Fenced-Fall. Das Schwestermodul `pins` prüft beide
  Pfade (`pins_test.go` `TestPinsInlineCodeZwischenInert`). Bricht ein
  künftiger Eingriff in `stripInlineCodeByLine` (`markdown.go:135`) die
  positionserhaltende Inline-Code-Leerung für das Marker-Muster, würde ein
  Inline-`` `<!-- immutable: sha256:… -->` `` zum Live-Pin, ohne dass ein
  `immutable`-Test failt.
- **verifizierbar:** ja — ein Inline-Code-Fixture-Test (Marker zwischen
  Backticks, sonst leere Datei) würde die Inertheit absichern bzw. die
  Regression zeigen.

### INFO-1 — Commit-/DoD-Aussage „8 Tests" vs. 7 implementierte

- **kategorie:** INFO
- **quelle:** Maintainability (Commit-/DoD-Drift)
- **pfad:** Commit `b62a520` (Botschaft „8 Tests") +
  `slice-052-immutable-modul.md` §3 (8 Test-Szenarien gelistet) vs.
  `immutable_test.go` (7 `func Test…`).
- **befund:** Commit-Botschaft und Slice-DoD listen 8 Test-Szenarien
  (Happy/Reflow/ausgenommener-Abschnitt/Negative/kein-Marker/**Modul-aus**/
  Fence-inert/erster-Marker); implementiert sind 7 — das fehlende 8. ist
  genau „Modul-aus" (deckungsgleich mit MEDIUM-1).
- **verifizierbar:** ja (`grep -c '^func Test' immutable_test.go` → 7).

---

## Negativbefunde (geprüft, ohne Befund)

- **Core-Hash-Definition vs. `DC-FA-IMM-001.a`:** `immutableCoreHash`
  (`immutable.go:59`) rekonstruiert die rohe Datei ohne die **erste**
  Marker-Zeile (`no == markerLine`) und ohne die `exclude-sections`-Ranges
  (`inRanges(excluded, no)`), kollabiert dann jede Whitespace-Folge per
  `pinWhitespaceRE` zu einem Leerzeichen, trimmt den Rand und bildet
  SHA-256 — bytegleiche Normalisierungsfunktion wie `pins`
  (`pins.go:141`, `pinWhitespaceRE` aus `pins.go:20`). Deckt Schritt 2–4
  der Spezifikation. Konform.
- **Marker-Erkennung auf `Line.Text`:** `firstImmutablePin`
  (`immutable.go:44`) iteriert den `PreprocessMarkdown`-Output; Fence- und
  fence-interne Zeilen sind aus `lines` **entfernt** (`markdown.go`
  `proseLines`/`PreprocessMarkdown`), Inline-Code ist positionserhaltend
  geleert — der Regex `immutableRE` matcht dort nicht. Fenced-Inertheit
  per `TestImmutableMarkerInFenceInert` belegt. „Erster Marker je Datei"
  korrekt (Funktion kehrt beim ersten Treffer zurück; `TestImmutableFirstMarkerWins`
  verankert den Befund an Zeile 4). Konform.
- **Zweiter Marker im Body bleibt Core-Bestandteil:** `immutableCoreHash`
  schließt **nur** `markerLine` (die erste) aus; ein zweiter Marker bleibt
  im Core. Konsistent mit der Spezifikation („die **Marker-Zeile**",
  Singular, Z. 916–918) und selbst-konsistent (wer den Pin setzt, rechnet
  über denselben Core) — keine überraschenden `core-drift`-Befunde. Kein
  Befund.
- **`exclude-sections` am Dateiende (`to == 0`):** `excludedRanges`
  (`matrix.go:294`) lässt für den letzten ausgenommenen Abschnitt `r.to`
  bei 0; `inRanges` (`matrix.go:320`) schließt bei `to == 0` alle Zeilen
  `>= from` ein. `## Geschichte` am ADR-Ende ist damit vollständig
  abgedeckt (`TestImmutableExcludedSection` mit angehängter Zeile → 0
  Befunde). Heading-Erkennung läuft via `proseLines` (`markdown.go:279`),
  also fence-bewusst — kein falsch erkanntes Heading im Fence. Konform.
- **Determinismus/Read-only (`DC-QA-02`/`DC-QA-03`):** `CheckImmutable`
  (`immutable.go:24`) nimmt nur `file`/`lines`/`content`/`cfg` —
  **kein** `Filesystem`, kein git, kein Netz; gelesen wird ausschließlich
  die gescannte Datei selbst. Default-off ist strukturell garantiert
  (`immutable` ∉ `defaultModules()`, `config.go:17`; Dispatch hinter
  `st.applies("immutable", …)`, `run.go:147`). Diagnose-only (kein
  `--repair`-Pfad). Konform (Test-Lücke separat unter MEDIUM-1).
- **Arch-/Import-Regeln (ADR-0005/0012):** `immutable.go` importiert nur
  `core/model` (+ stdlib `crypto/sha256`, `encoding/hex`, `regexp`,
  `strings`); die Kern-Richtung `model ← rules` ist gewahrt, kein
  `app`-/`adapter`-Import, kein `port/driven` nötig (strenger als `pins`).
  Konform.
- **Config-Wiring:** `rawImmutable` (`configyaml.go:85`, `scope` +
  `exclude-sections`), `applyImmutable` (`configyaml.go:289`) überträgt
  `ExcludeSections` ohne Pflichtfeld (opt-in pro Datei), `scopeOfImmutable`
  (`configyaml.go:393`) ist in `applyScopes` eingehängt (Roots-Pflicht bei
  gesetztem Scope). `validModules`/`EffectiveModules` (`config.go:13,232`)
  kennen `immutable`; Glossar + `DC-FA-CLI-002` führen es (Lastenheft
  Z. 1201). Konform.

---

## Kategorie-Summary

| Kategorie | Anzahl |
| --- | --- |
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

---

## Verdikt: **NACHBESSERN** (ein MEDIUM)

Die **Code-Korrektheit des Moduls selbst ist sauber**: Core-Definition,
Normalisierung (pins-identisch), Marker-/Fence-Inertheit,
Erster-Marker-Vertrag, `exclude-sections` bis EOF, Read-only/Determinismus
und die Kern-Import-Richtung entsprechen `DC-FA-IMM-001` /
`DC-FA-IMM-001.a` / ADR-0023. Es wurde **kein** Korrektheits- oder
Spec-Treue-Fehler im Produktionspfad gefunden (kein HIGH).

Das einzige blockierende Finding (MEDIUM-1) ist eine **Test-Lücke**, kein
Laufzeitfehler: Der von den Schwestermodulen `pins`/`versions` etablierte
`Modul-aus`-Test (über `Run()`, inaktives Modul → 0 Befunde) fehlt, und
damit ist der gesamte Wiring-/Gating-Pfad des Moduls ungetestet — der
`DC-QA-02`-Default-off-Vertrag (zugleich AK 6 der Anforderung) hat keinen
ausführbaren Wächter. Per Reviewer-Skill blockiert ein MEDIUM typischer-
weise; der Befund ist eng umrissen und mit einem einzelnen `Run()`-Test
(Vorbild `TestPinsNichtAktiv`) zu schließen. LOW-1 (Inline-Code-Inertheit
ohne direkten Test) und INFO-1 (Commit-/DoD-Zähler) sind nicht blockierend
und können im selben Nachbesserungs-Schritt mitlaufen.
