# Review-Report — slice-051 (matrix: token-basierte Referenz-Richtung) — R2

| Feld | Wert |
|---|---|
| Datum | 2026-06-28 |
| Gegenstand | Uncommittete Änderungen slice-051 — Modul `matrix`, `DC-FA-MTX-003` (Token-Referenz-Richtung, Provenance-Marker, `exempt-paths`-Grandfathering); R2 = Verifikation der R1-Auflösungen + kritische Delta-Prüfung |
| Reviewer | unabhängiger R2-Reviewer (Skill 1.2.0) |
| Eingang | R1-Report `2026-06-28-slice-051-matrix-token-richtung-r1.md`; aktueller Working-Tree (`git diff` + untracked ADR-0022/slice-051); Lastenheft `DC-FA-MTX-003`; ADR-0022; Spezifikation §DC-FA-MTX-001.a Schritt 6 + §2-Schema + §4; `.d-check.yml`; Handbuch §4.7; Reviewer-Skill |
| Gate-Belege (eigener Lauf) | `make test` grün; `make lint` 0 issues; `make doc-check` 142 Dateien / 0 Befunde; Negativ-Probe eigenhändig reproduziert (unmarkierter `slice-099` im neuen ADR-0099 ⇒ `matrix-forbidden`, Exit 1) |
| Code im Fokus | `internal/hexagon/core/rules/matrix.go`, `…/matrix_test.go`, `…/model/config.go`, `…/configyaml/configyaml.go`, `…/app/suggest.go`, `…/cli/config_template.go`, `.d-check.yml`, Spec/Lastenheft/Handbuch |

Hinweis: das `befund`-Feld bleibt beobachtend (Skill-Regel „kein Lösungsvorschlag im
Befund"). Die R1-Verifikation ist der Pflichtteil dieses Laufs; neue Befunde folgen
darunter.

---

## R1-Auflösungen verifiziert

### F-4 (MEDIUM) — Carve-out-Regressionstests Fence/`exclude-sections` — AUFGELÖST
- **Beleg:** `internal/hexagon/core/rules/matrix_test.go` `TestMatrixTokenReferenz`
  trägt jetzt zwei zusätzliche Fixtures. `docs/plan/adr/0053-f.md` =
  `"# A\n```\nslice-077\n```\n"` (Token `slice-077` in einer Fence) und
  `docs/plan/adr/0054-h.md` = `"# A\n## Geschichte\nWegen slice-088 (Provenance).\n"`
  (Token `slice-088` unter `## Geschichte`). Beide sind ADR-klassifiziert
  (`docs/plan/adr/[0-9]*.md`) und **nicht** grandfathered (exempt-paths deckt nur
  `0001-*`), die Regel `adr→slice` ist verboten, die Klasse `slice` trägt
  `token: slice-\d{3}` und `cfg.ExcludeSections = ["Geschichte"]`. Das `want` ist
  unverändert (`0050-x.md:2`, `0052-l.md:2`) — beide neuen Fixtures müssen
  **befundfrei** bleiben.
- **Branch-Treffer geprüft:** 0053-f isoliert den Fence-Zweig — `proseLines`
  (markdown.go:24) verwirft Zeile 3 im Fence, ohne `exclude-sections` und ohne
  Marker; eine Mutation des Token-Scans auf rohe Zeilen ließe `0053-f.md:3`
  auftauchen → Test rot. 0054-h isoliert den `inRanges(excluded,…)`-Zweig — das
  Heading `## Geschichte` (parseATXHeading liefert „Geschichte", deckungsgleich mit
  `ExcludeSections`) erzeugt `excludedRanges` `{from:2,to:0}`; die Prosa-Zeile trägt
  **nicht** den Marker-Substring `d-check:status-provenance`, also greift hier
  ausschließlich der Sektions-Carve-out. Wegfall des `inRanges`-Checks ließe
  `0054-h.md:3` auftauchen → Test rot.
- **Lauf:** `make test` grün (Paket `…/rules` ok).

### F-1 (LOW) — alle Token je Zeile — AUFGELÖST
- **Beleg:** `matrix.go:112` lautet jetzt
  `for _, loc := range c.Token.FindAllStringIndex(stripped, -1)` und emittiert je
  Treffer ein Finding mit `Target = stripped[loc[0]:loc[1]]`. `FindAllStringIndex`
  liefert die nicht-überlappenden Treffer in Dokument-Reihenfolge (links→rechts),
  also deterministisch; `SortFindings` (File/Line/Rule/Target/Reason) re-sortiert
  und dedupliziert ohnehin total. Determinismus (DC-QA-02) gewahrt.
- **Einschränkung:** kein dediziertes Regressions-Fixture (siehe N-1).

### F-2 / F-3 / F-5 / F-6 — als bekannte Grenzen dokumentiert — ERFÜLLT
- **Beleg:** `slice-051-matrix-token-richtung.md` §4 listet sie unter
  „Bekannte LOW-Grenzen (R1, Folge-CR-tauglich)": Inline-Code wird nicht geleert
  (F-5), Badge-/verschachtelte Links nur partiell gestrippt (F-3), der Marker ist
  nackter Substring und nimmt die ganze Zeile aus (F-2); `exempt-paths`
  überspringt die Datei komplett (F-6, in ADR-0022 begründet). Alle vier sind als
  bewusst-minimal/ohne aktuellen Repo-Trigger benannt — Folge-CR-tauglich, kein
  Eingriff in diesem Slice gefordert.

---

## Findings (neuer Stoff seit R1)

### N-1 — F-1-Fix („alle Token je Zeile") ist ungetestet
- **Kategorie:** LOW
- **Quelle:** Maintainability / DC-FA-MTX-003
- **Pfad:** `internal/hexagon/core/rules/matrix_test.go:356` (`TestMatrixTokenReferenz`)
- **Befund:** Die Umstellung auf `FindAllStringIndex` (matrix.go:112) ist korrekt,
  aber kein Fixture trägt **zwei** verbotene Token in einer Prosa-Zeile: 0050-x.md
  hat einen Token, 0052-l.md einen Link, alle übrigen Fixtures höchstens einen
  Treffer je Zeile. Eine Rück-Mutation auf `FindStringIndex` (nur Erst-Treffer)
  ließe keinen Test rot werden — die in slice-051 §4 als „behoben" geführte
  Eigenschaft ist nicht regressionsgesichert.
- **Verifizierbar:** ja — Mutationsprobe (`FindAllStringIndex` → `FindStringIndex`)
  bleibt grün; ein Fixture mit zwei verbotenen Token je Zeile (erwartet zwei
  Befunde) würde die Lücke schließen.

---

## Negativbefunde (geprüft, ohne Befund)

- **FindAllStringIndex-Verhalten im Dogfood:** `make doc-check` meldet 142 Dateien
  / 0 Befunde — die Umstellung erzeugt keine zusätzlichen Dogfood-Befunde (es
  existiert keine unmarkierte Mehr-Token-Zeile im nicht-grandfatherten Körper).
  Ohne Befund.
- **Komplexität/Lint:** `make lint` 0 issues; gocognit-Schwelle 20 / gocyclo 15
  (`.golangci.yml`) — die neue innerste `FindAllStringIndex`-Schleife in
  `tokenFindings` und `CheckMatrix` bleiben darunter. Ohne Befund.
- **lineageValues-Refactor:** `lineageValues(cfg, content)` (matrix.go:128) =
  `if cfg.AllowSupersedeLineage { return supersedeFieldValues(content,
  cfg.SupersedeFields) } return nil` ist semantisch identisch zur Inline-Form
  (Flag aus ⇒ nil; Flag an ⇒ Feldwerte); `supersedeFieldValues` liefert bei
  leeren Feldern ohnehin nil. Default-aus byte-identisch. Ohne Befund.
- **F-4-Fixtures korrekt konstruiert:** 0053-f/0054-h sind ADR-klassifiziert,
  nicht exempt, treffen genau je einen Carve-out (Fence bzw. `## Geschichte`),
  tragen keinen Marker und ändern das `want` nicht. `ExcludeSections:
  ["Geschichte"]` deckt das Fixture-Heading `## Geschichte` exakt. Ohne Befund.
- **Negativ-Probe eigenhändig:** ein temporärer `docs/plan/adr/0099-probe.md`
  (Status Proposed, „gründet auf slice-099 als Entscheidungsgrundlage", unmarkiert,
  nicht von exempt-paths gedeckt) erzeugt `0099-probe.md:5 slice-099
  matrix-forbidden`, Exit 1; nach Entfernen Working-Tree sauber. Das Token-Gate
  feuert. Ohne Befund.
- **Config-Surface/Fail-closed:** `compileMatrixToken` (configyaml.go) lehnt
  nicht kompilierende Regexe und Leerstring-Matcher ab (Exit 2),
  `TestDecode_MatrixTokenFailClosed` (`[unclosed`, `x*`) + Happy-Pfad
  `TestDecode_MatrixTokenUndExemptPaths` decken es ab. `model.MatrixClass.Token`
  / `MatrixConfig.ExemptPaths` korrekt durchgereicht. Ohne Befund.
- **suggest.go / config_template.go:** der generierte `token:`-Block, die
  `{from: adr, to: slice}`-Regel und der `# exempt-paths`-Kommentar sind
  konsistent mit `.d-check.yml`; die Regel wird auskommentiert, wenn adr oder
  slice inaktiv ist; `make test` (cli-Akzeptanz, Round-Trip-Decode) grün. Ohne
  Befund.
- **Dogfood-Disziplin Spec:** `spec/spezifikation.md` entfernt den bare
  `slice-015`-Token aus dem §-Körper (Kalibrierungs-Befund-Prosa); alle übrigen
  bare `slice-NNN` in lastenheft.md/spezifikation.md liegen in den
  Historie-/Änderungstabellen, die `exclude-sections: [Historie, "7. Historie",
  Geschichte]` deckt. doc-check bestätigt 0 Befunde. Ohne Befund.
- **ADR-0022-Verweisrichtung:** `Bezug:` zeigt aufwärts (DC-FA-MTX-003/Lastenheft)
  bzw. seitlich (Regelwerk §Referenz-Richtung, ADR-0021, MR-006); `Schärft:` zeigt
  aufwärts auf die Spezifikation. `slice-051` erscheint nur Z.77 (Marker auf der
  rohen Zeile, Inline-Code `'slice-\d{3}'` matcht das Regex nicht, Link auf
  `.d-check.yml` link-gestrippt) und Z.83 unter `## Geschichte` (excluded) — kein
  unmarkierter Abwärts-Token. Ohne Befund.
- **Kohärenz Lastenheft↔Spec↔ADR↔Code↔Dogfood↔Handbuch:** `DC-FA-MTX-003`
  (Lastenheft 0.31.0) ↔ §DC-FA-MTX-001.a Schritt 6 (Marker auf **roher** Zeile,
  Match comment-whitespace-unabhängig; Links/Fences/exclude-sections aus;
  exempt-paths Ganz-Datei; fail-closed) ↔ ADR-0022 (Accepted, drei Mechaniken) ↔
  `matrix.go`/`configyaml.go`/`model` ↔ `.d-check.yml` (token auf slice, adr→slice
  verboten, exempt-paths 0001–0021) ↔ Handbuch §4.7 ↔ Reviewer-Skill 1.2.0
  (Marker-Ehrlichkeit) sind deckungsgleich; keine Quelle behauptet Verhalten, das
  der Code nicht erfüllt. Spec-Schema führt `matrix.classes[].token` und
  `matrix.exempt-paths`; §4 `matrix-forbidden` um die Token-Form erweitert (kein
  neuer Grund-Code). Ohne Befund.
- **Index/Version/Roadmap:** ADR-README listet ADR-0022 (Accepted, 2026-06-28);
  CHANGELOG 0.31.0; Roadmap-Flip auf aktive welle-40; Lastenheft-Version 0.31.0 —
  konsistent. Ohne Befund.

---

## Kategorie-Summary

| Severity | Anzahl |
|---|---|
| BLOCKER | 0 |
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 (N-1) |
| INFO | 0 |

---

## Verdikt

**Mergebar: ja.**

Die R1-Befunde sind belegt aufgelöst: F-4 (MEDIUM) durch die zwei isolierten,
branch-treffenden Fixtures 0053-f/0054-h (Carve-outs Fence und `## Geschichte`,
`want` unverändert); F-1 (LOW) durch die `FindAllStringIndex`-Umstellung (korrekt,
links→rechts deterministisch); F-2/F-3/F-5/F-6 sind als Folge-CR-taugliche Grenzen
in slice-051 §4 dokumentiert. Die Deltas sind sauber: kein neuer Dogfood-Befund
(doc-check 142/0), Lint 0 issues (Komplexität unter Schwelle), `lineageValues`
semantisch identisch, Negativ-Probe eigenhändig reproduziert. Es gibt keinen
BLOCKER/HIGH/MEDIUM.

Das einzige neue Finding N-1 (LOW) ist die fehlende Regressionssicherung der
F-1-Eigenschaft („alle Token je Zeile") — eine Test-Lücke ohne Verhaltensfehler,
analog zur R1-Logik bei F-4 als **nicht-blockierend** eingestuft. Empfehlung: ein
Zwei-Token-Zeilen-Fixture vor oder nach dem Release nachziehen.
