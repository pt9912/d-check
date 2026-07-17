# Review — slice-073 Implementierung R2 (nach R1-REJECT, Stand v0.45.1)

**Datum:** 2026-07-17
**Review-Art:** unabhängiger, kontext-getrennter Implementierungs-Review
(adversarial, nicht bestätigend). **Nachgeholt** — der Gegenstand ist bereits
als **v0.45.1** veröffentlicht. Zweite Runde nach dem R1-REJECT (F-1…F-5).
**Gegenstand:**
[`slice-073`](../plan/planning/done/slice-073-link-transparente-range-fortsetzung.md) —
Range `2954e4d` (feat: Range-Fortsetzung link-transparent) + `6925987`
(fix: Link-Ziel klammer-balanciert; zugleich v0.45.1-Tag-Commit). HEAD `e4880ee`;
seit v0.45.1 hat nur die Rücknahme von slice-074 `trace_table.go` angefasst, **nicht**
die Range-Parser-Teile von slice-073.
**Reviewer:** Claude (kontext-unabhängiger Lauf, ohne Zugriff auf Vorsessions-Analysen)
**Skill:** `.harness/skills/reviewer.md` v1.2.0
**Modell:** Opus 4.8 (1M context)

## Verifikations-Basis

- Eigener Build `make build` ⇒ `d-check:latest`
  (Image-ID `sha256:7a414767c0cd51e840b00bcae3aa7694c4353f2a9cc60098fc98c2d5c0d689e5`,
  gegen HEAD `e4880ee`).
- Ausgeliefertes Vergleichs-Image `ghcr.io/pt9912/d-check:v0.45.1`
  (`sha256:99281f8417f380bc3bca795f6b7914dd03e7b0e6b89c41713c2679853a8af989`,
  RepoDigest `sha256:5c5cf2d669f9…4bf5ad`). **Beide Images liefern in allen
  Fixtures identische Ausgabe** — die Befunde gelten für den ausgelieferten Stand
  und für HEAD.
- Alle Läufe netzlos/read-only:
  `docker run --rm --network none -v <fixture>:/repo:ro -w /repo <img> --trace [--require-complete]`.
- Mutationstest in einem losgelösten Worktree (`git worktree add --detach … HEAD`),
  `make test`, danach `git worktree remove --force`; Produktivbaum unberührt.
- Baseline: `make test` gegen HEAD ⇒ **exit 0**
  (`ok core/app`, `ok core/rules`).

## Eingangs-Kontext

- [`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
  (Lastenheft, unqualifizierte Range-Zusage),
  [`DC-FA-COV-001.a`](../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
  Schritt 3 (§Link-Transparenz),
  [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)/[`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  (mit-betroffen über den geteilten Parser),
  [`DC-FA-ID-001`](../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
- [ADR-0039](../plan/adr/0039-link-transparente-range-fortsetzung.md) (Proposed;
  §Fitness-Funktion, §Konsequenzen, §Geschichte),
  [`AGENTS.md`](../../AGENTS.md) §3,
  Vorlauf-Review [R1](2026-07-17-slice-073-implementation-r1.md) (REJECT)
- Kern-Code: `internal/hexagon/core/app/trace.go` (`skipLinkSuffix`, `expandRange`,
  `rangeAwareIDs`, `coverageRefs`), `internal/hexagon/core/rules/markdown.go`
  (`LinkSuffixEnd`, `matchBracket`), `internal/hexagon/core/app/trace_cross.go`
  (`forwardEdges`/`backwardEdges`), Tests in `trace_coverage_test.go`/`trace_cross_test.go`

---

## Findings

### R2-F-1 — LOW · Die Link-Transparenz ist auf der `trace.cross-consistency`-Seite von keinem Test verriegelt

**kategorie:** LOW
**quelle:** [ADR-0039](../plan/adr/0039-link-transparente-range-fortsetzung.md)
§Fitness-Funktion / Entscheidung 4 („Ein Fix, zwei Konsumenten"),
[`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency),
fehlende Negativtests bei einem öffentlichen Vertrag (Rest-Befund aus R1-F-4)
**pfad:** `internal/hexagon/core/app/trace_cross_test.go` (kein link-transparenter
Fall); `internal/hexagon/core/app/trace_cross.go:252,278` (`rangeAwareIDs` über
`row.cells[…]`)

**befund:** Kein Test in `trace_cross_test.go` fährt eine **verlinkte** Range/Enum
durch den Kreuzverweis-Abgleich; `TestCrossConsistencyRangeAware` nutzt die
**unverlinkte** Form `GG-SIM-001..004`. Die Vorwärts-/Rück-Sicht speist ihre
Zellen (`row.cells[bt.primary/secondary]`) über den **Tabellen-Zellen-Leser**
(`markdownTables`) in `rangeAwareIDs` — ein anderer Vorlauf als der Prosa-Text der
Coverage-Quelle. Genau diese Kombination (Tabellenzelle **mit** Link-Suffix) hat
keinen Testfall. **Failure-Szenario:** ein späterer Umbau der Zellen-Extraktion der
Kreuzverweis-Sichten (z. B. Vor-Splitten der Zelle oder Code-Span-Stripping vor
`rangeAwareIDs`) bräche die Link-Transparenz **nur** auf der Cross-Seite; die
Coverage-Tests blieben grün (sie durchlaufen die Zellen-Extraktion nicht), und die
Regression fiele unbemerkt aus. Der Slice-DoD hakt „Tests (positiv) … für
`trace.cross-consistency` (2 Differenzen ⇒ 0)" ab — dieser Test existiert nicht.

**verifizierbar:** ja — die R2-Mutation (siehe Negativbefund zu F-4) kippt
**keinen** `TestCrossConsistency*`-Fall; die Cross-Seite ist laufzeit-korrekt (eigene
Fixture unten), aber test-blind für diese Achse.

---

## Negativbefunde (geprüft, ohne Befund)

### Stand der fünf R1-Befunde (je mit eigenem Beleg)

- **R1-F-1 (HIGH) — GESCHLOSSEN.** Eigene Fixture, `trace.coverage`,
  `ranges: true`, Zelle **ohne jede** Range-Notation:
  | Zelle | v0.45.1 | HEAD |
  |---|---|---|
  | `` [`GG-QA-001`](../specs/Rev(2)/002/003.md) `` | **2 Waisen, Exit 1** | 2 Waisen, Exit 1 |
  | `` [`GG-QA-001`](../a(1)..003.md) `` | **2 Waisen, Exit 1** | 2 Waisen, Exit 1 |
  `GG-QA-002`/`GG-QA-003` erscheinen als `WAISE`, `--require-complete` bleibt **rot**.
  Die URL-Interna (`/002/003`, `..003`) werden **nicht** mehr als Enum/Range
  expandiert. Beleg: `LinkSuffixEnd` (`markdown.go:318`) grenzt das Ziel
  **klammer-balanciert** via `matchBracket` ab (statt der verworfenen Regex `[^)]*`)
  und konsumiert `` [`GG-QA-001`](…/Rev(2)/002/003.md) `` vollständig; hinter `)`
  steht dann nichts mehr. Der stille-Grün-Pfad aus R1 ist weg.
- **R1-F-2 (MEDIUM) — GESCHLOSSEN.** Fixture `` [`GG-QA-001`](https://x.org/A_(b))..003 ``
  ⇒ v0.45.1 **und** HEAD expandieren auf `001/002/003`, **0 Waisen, Exit 0** (unter
  R1: 2 Waisen). Der Range-Parser liest das Ziel jetzt mit **demselben**
  `matchBracket` wie `parseLinkAt`/`ids` — eine Link-Definition im Repo. Die
  Ziel-Klasse „geklammertes Ziel" ist im Spec-Text
  ([`DC-FA-COV-001.a`](../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
  §Link-Transparenz, „geklammertes Ziel") und im ADR (§Alternativen, „Ziel
  klammer-balanciert über den kanonischen Reader") ausgesprochen.
- **R1-F-3 (MEDIUM) — GESCHLOSSEN.** Die drei durch F-1 widerlegten Zusagen sind
  korrigiert: CHANGELOG §0.45.1 dokumentiert jetzt „**Keine Falsch-Expansion mehr bei
  Klammern im Linkziel**", „**Betroffen: v0.44.1 und v0.45.0**", „**falsche Deckung,
  die Waisen versteckt** (Exit 1 → Exit 0)"; das Handbuch §5 sagt „Klammern **im
  Linkziel** sind dagegen unproblematisch: das Ziel wird klammer-balanciert
  abgegrenzt … `` [`GG-QA-001`](https://x.org/A_(b))..003 `` **expandiert**"; der
  falsche Kommentar `trace.go:369–371` existiert nicht mehr (Funktion neu geschrieben,
  Kommentar `trace.go:370–382` beschreibt die balancierte Abgrenzung korrekt).
- **R1-F-4 (MEDIUM) — WEITGEHEND GESCHLOSSEN, LOW-Rest → R2-F-1.** Die
  Fitness-Funktion ist jetzt **mechanisiert**: `TestCoverageRefsLinkTransparent`
  (`trace_coverage_test.go:191`) läuft über `coverageRefs` (den Waisen-Pfad) und
  `TestExpandRangeBalancedLinkTarget` (`:159`) pinnt die F-1-Zeichenklasse.
  **Mutationsbeleg (eigener Lauf):** `LinkSuffixEnd` auf die naive „bis zur ersten
  `)`"-Form zurückgedreht ⇒ `make test` **kippt** (Exit 2), und zwar
  `TestExpandRangeBalancedLinkTarget/*` (u. a. `Pfadsegmente_hinter_Klammer-URL_…` =
  der F-1-Fall) **und** `TestCoverageRefsLinkTransparent/verlinkt_mit_Klammer-URL,_ohne_Range`
  (Konsumenten-Ebene: `abgedeckte Anforderungen = 3, want 1`). Die R1-Behauptung
  „`make test` bleibt grün, wenn der Fix zurückgedreht wird" ist damit für die
  Coverage-Seite **widerlegt**. **Rest:** kein `TestCrossConsistency*` kippt bei
  derselben Mutation ⇒ die Cross-Seite ist test-blind (siehe R2-F-1). Die
  Fail-closed-Tests (`TestExpandRangeLinkTransparentFailClosed`) prüfen weiterhin nur
  `err != nil`, nicht die Fehlerklasse — das war in R1 eine Randnotiz und bleibt eine.
- **R1-F-5 (LOW) — GESCHLOSSEN.** Die ADR-0039-Indexzeile (`docs/plan/adr/README.md:54`)
  trägt jetzt **6** Pipe-Zeichen, gleich wie ADR-0038 (`awk`-Zählung: `0038: 6`,
  `0039: 6`); Titel, Status (`Proposed`), Datum (`2026-07-17`) und Bezug
  (`DC-FA-COV-001`, `DC-FA-XREF-001`) sind vollständig.

### Weitere geprüfte Bereiche

- **`trace.cross-consistency` laufzeit-korrekt (Gegenprobe zu R2-F-1):** eigene
  Fixture, Rück-`Bezug`-Zelle `` [`GG-SIM-001`](../docs/x.md)..009 `` gegen
  Vorwärts-`GG-SIM-001..004` ⇒ v0.45.1 meldet genau **5** Rück-Kanten-Differenzen
  (`GG-SIM-005..009`), Exit 1; identische Mengen (`..004` verlinkt vs. unverlinkt) ⇒
  **0 Differenzen**, Exit 0. Der geteilte Fix wirkt am realen Tabellen-Zellen-Pfad —
  die Lücke ist die **Test**-Verriegelung, nicht die Funktion. Ohne Befund (Laufzeit).
- **Fail-closed hinter Link:** `skipLinkSuffix` sitzt nach der Familien-/Breiten-
  Ableitung (`trace.go:513–519`), vor beiden Suffix-Zweigen; `AAA>BBB`/Breiten-Mismatch
  bleiben hinter einem Link Fehler (`TestExpandRangeLinkTransparentFailClosed`, Lauf
  grün). Ohne Befund.
- **`matchBracket` nicht gierig:** liefert am **ersten** ausbalancierten `)` zurück
  (`markdown.go:291–305`); ein `](a.md)` gefolgt von späterem `(x)..003` konsumiert nur
  `(a.md)`, `..003` bleibt hinter Text/Leerzeichen und expandiert **nicht** — keine neue
  Über-Konsumption in die Falsch-Deckungs-Richtung. Ein **unbalanciertes** Ziel
  (`](a(b.md)..003`) ⇒ `matchBracket` false ⇒ `LinkSuffixEnd` −1 ⇒ keine Expansion
  (sichere, laute Richtung). Ohne Befund.
- **No-op für unverlinkte/`ranges:false`-Formen:** `skipLinkSuffix("..003")` gibt
  `..003` unverändert zurück (`LinkSuffixEnd` −1); `rangeAwareIDs` überspringt bei
  `ranges:false` die Expansion (`trace.go:494`). `TestExpandRange`/`TestRangeAwareIDs`
  grün. Byte-Identität der Alt-Formen gewahrt. Ohne Befund.
- **„Genau ein Suffix" / Mehrfach-Peeling:** `skipLinkSuffix` ist nicht geschleift;
  zwei Suffixe, ein Zeichen/Whitespace zwischen `)` und `..` ⇒ keine Expansion
  (`TestExpandRangeLinkTransparent` Negativfälle, grün). Deckt sich mit ADR-0039
  Entscheidung 2. Ohne Befund.
- **Code-Span-ID ohne Link + Direkt-Range (`` `GG-QA-001`..003 ``):** expandiert
  **nicht** (eigene Fixture: 2 Waisen). `skipLinkSuffix` gibt bei Backtick ohne
  folgendes Link-Suffix das **Original**-`rest` zurück (der Backtick blockt
  `rangeSuffix`). Verhalten **unverändert** gegenüber vor slice-073 (der alte Code
  übersprang gar nichts) — keine Regression, außerhalb der Link-Transparenz-Zusage,
  sichere Richtung. Ohne Befund.
- **`<>`-quotiertes Linkziel mit Klammer — INFO, nicht gemeldet als Finding:**
  `LinkSuffixEnd` balanciert nur `(`/`)` und honoriert das `<…>`-Quoting **nicht**,
  das `parseLinkAt` via `NormalizeTarget` kennt. Ein Ziel `` [`ID`](<a)b>)..003 ``
  würde am `)` innerhalb `<…>` frühzeitig abgegrenzt ⇒ keine Expansion. Das
  qualifiziert die CHANGELOG-/ADR-Aussage „**eine** Link-Definition im Repo" minimal,
  hat aber nur ein **sicheres** (lautes, false-orphan) Failure-Szenario und eine
  exotische Eingabe (angle-bracket-Ziele für ID-Links kommen im belegten Bestand nicht
  vor, ADR §Kontext: 40/40). Dokumentationswürdige Annahme, kein blockierender Befund.
- **Referenz-Richtung (SDP):** `slice-073` erscheint in ADR-0039 nur in §Geschichte
  (Provenienz: „Umsetzender Slice", „Anlass … R1-F-1") und in der Spec-Historie in der
  **Slice-Spalte** — nirgends als Entscheidungsgrundlage im ADR-/Spec-Körper. Kein
  Provenance-Marker `<!-- d-check:status-provenance -->` im Diff zu prüfen. Ohne Befund.
- **Hard Rules ([`AGENTS.md`](../../AGENTS.md) §3):** kein `//nolint` in `trace.go`/
  `markdown.go` (§3.2); Verifikation ausschließlich über `make build`/`make test` +
  Image, kein Host-`go` (§3.1); der Fix ist rein lexikalisch — kein neuer Import,
  ADR-0005-Schnitt unberührt; ADR-0039 `Proposed` (§3.5 unberührt); keine
  Gate-Schwelle gesenkt (§3.6). Ohne Befund.
- **Netzlosigkeit ([`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)):**
  alle Fixtures liefen mit `--network none`; der Fix fügt keinen I/O-/Netz-Pfad hinzu.
  Ohne Befund.
- **Determinismus ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)):**
  `matchBracket`/`LinkSuffixEnd` sind rein positional-deterministisch; wiederholte
  Fixture-Läufe (beide Images) byte-identisch. Ohne Befund.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 1 | R2-F-1 |
| INFO | 0 | (eine INFO-Notiz als Negativbefund geführt) |

---

## Verdikt

**ACCEPT-WITH-NITS.**

Alle blockierenden R1-Befunde sind **durch eigene Messung** geschlossen:
R1-F-1 (HIGH) und R1-F-2 (MEDIUM) haben **dieselbe Wurzel** — die zweite
Link-Definition (`[^)]*`) — und sind durch die klammer-balancierte Abgrenzung über
`rules.LinkSuffixEnd` (Commit `6925987`) behoben; die kritische, stille Richtung
(Falsch-Deckung, die Waisen versteckt) liefert jetzt wieder **laut rot**
(2 Waisen / Exit 1) und die legitime verlinkte Range **grün** (0 Waisen / Exit 0),
belegt gegen `ghcr.io/pt9912/d-check:v0.45.1` **und** HEAD. R1-F-3 (falsche Zusagen)
und R1-F-5 (Index-Zeile) sind mit demselben Zug korrigiert.

R1-F-4 — der Grund, warum F-1 durchkam — ist der wichtigste Fortschritt: die
Fitness-Funktion ist jetzt an ihrer echten Wirkungsstelle mechanisiert
(`TestCoverageRefsLinkTransparent` + `TestExpandRangeBalancedLinkTarget`), und meine
Mutation der balancierten Abgrenzung **kippt** diese Tests — der Sensor fängt den
F-1-Rückdreh. Es bleibt **eine** Nit: die zweite geteilte Konsumenten-Seite
`trace.cross-consistency` ist laufzeit-korrekt (eigene Fixture: 5/0 Differenzen), aber
von **keinem** Test verriegelt (R2-F-1, LOW) — die Mutation kippt dort nichts. Das ist
kein stiller-Grün-Pfad im ausgelieferten Code (der geteilte Kern ist über die
Coverage-Seite gesensort), sondern eine enge Test-Vollständigkeitslücke gegen einen
künftigen, cross-spezifischen Umbau. Kein HIGH/MEDIUM offen ⇒ nicht blockierend.
