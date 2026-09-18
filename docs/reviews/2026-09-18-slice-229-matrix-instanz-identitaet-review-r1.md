# Review-Report: slice-229 — 2026-09-18

**Review-Art:** Code — geprüft gegen Slice-Plan, [ADR-0087](../plan/adr/0087-matrix-instanz-identitaets-ausnahme.md),
`spec/lastenheft.md`/`spec/spezifikation.md` (DC-FA-MTX-003) und
`AGENTS.md` §3 (insbesondere §3.2, §3.6, §3.7, §3.8).

**Gegenstand:** `git log 498bfe82~1..HEAD` (drei Commits: `498bfe82`
Spec-Erweiterung, `32f0358e` ADR+Slice-Claim, `d4656a63` Go-Implementierung).

**Skill:** `.harness/skills/reviewer.md` @ `1b689b1f` (Version 1.16.0)

**Modell:** claude-sonnet-5 · **Datum:** 2026-09-18

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-229-matrix-instanz-identitaets-ausnahme.md`
  (vollständig gelesen)
- [ADR-0087](../plan/adr/0087-matrix-instanz-identitaets-ausnahme.md) (vollständig gelesen)
- `docs/plan/cr/2026-09-18-cr-eingehend-pg-change-feed-matrix-instanz-identitaet.md`
  (Anlass, vollständig gelesen — nicht Teil des Diffs)
- `spec/lastenheft.md` DC-FA-MTX-003 (Diff + neue Akzeptanzkriterien)
- `spec/spezifikation.md` §DC-FA-MTX-001.a Schritt 7, §2-Schema, Grund-Code-Zeile
- `AGENTS.md` §3.1/§3.2/§3.6/§3.7/§3.8/§5
- Präzedenzfall Supersede-Lineage-Ausnahme (Stil-Vergleich, kein bekannter
  offener Finding dazu)
- Gegenprobe: `make test` (grün, inkl. unverändertem `TestMatrixTokenReferenz`),
  `make lint` (0 issues), `make doc-check` (795 Dateien, 0 Befunde)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Fitness-Function-Tabelle von ADR-0087 behauptet, `make doc-check` verifiziere, dass das erweiterte YAML-Config-Beispiel in `spec/spezifikation.md` syntaktisch gültig bleibt. `make doc-check` fährt laut `.d-check.yml` nur die Module `links, anchors, ids, matrix, codepaths, spans, hostpaths, versions, structure, diagrams, citations` — keines davon dekodiert gefencte YAML-Blöcke. Der einzige Mechanismus, der Doku-YAML gegen `configyaml.Decode` prüft (`TestDocExamples_ConfigBeispieleValidieren`, läuft unter `make test`), ist per eigenem Kommentar bewusst auf vier Nutzer-Doku-Dateien beschränkt (`docs/user/benutzerhandbuch.md`, `docs/user/operations.md`, `README.md`, `README.de.md`) und nennt `spec/spezifikation.md` nicht. Die „matrix (Selbstanwendung)"-Zeile prüft Klassen-/Token-Referenzen zwischen Dokumentklassen, keine YAML-Schema-Gültigkeit — beide Teile der Zeile zusammen verweisen auf ein Gate, das diesen Vertrag nicht hält. | ADR-0087 §Fitness Function | `docs/plan/adr/0087-matrix-instanz-identitaets-ausnahme.md` (Tabellenzeile „matrix (Selbstanwendung)"); Gegenbeleg `internal/adapter/driven/configyaml/docexamples_test.go:16-24` (`docExampleFiles()`) | ja — `make doc-check` gegen eine absichtlich mit ungültigem YAML verfälschte Kopie des Beispiels bliebe grün; `go test ./internal/adapter/driven/configyaml/... -run TestDocExamples` deckt die Datei nicht ab | adr-fitness-function-names-wrong-gate |
| F-2 | MEDIUM | `TestDecode_MatrixAllowIfSameIDFailClosed` variiert in allen drei Unterfällen ausschließlich das `token` der `from`-Klasse (`slice`); die `to`-Klasse (`review`) trägt in jedem Unterfall ein gültiges Token mit genau einer Capture-Gruppe. Da `validateMatrixAllowIfSameID` beide Seiten in einer Schleife prüft und beim ersten Fehler zurückkehrt, bricht jeder der drei Testfälle bereits auf der `from`-Seite ab — der symmetrische Code-Pfad für die `to`-Seite (`side.token == nil` bzw. `side.token.NumSubexp() != 1` für die Ziel-Klasse) wird von keinem Test isoliert ausgelöst. Die im Lastenheft wörtlich formulierte Fehlkonfigurations-AK („deren `to`-Klasse kein `token` trägt [...]") bleibt damit ohne dedizierten Beleg, obwohl der Code sie korrekt behandelt. | DC-FA-MTX-003 (Akzeptanzkriterium „Instanz-Identität Fehlkonfiguration") | `internal/adapter/driven/configyaml/configyaml_test.go:698-719` (`TestDecode_MatrixAllowIfSameIDFailClosed`) | ja — ein vierter Unterfall, der nur das `token` der `review`-Klasse bricht (`from` bleibt gültig), würde die Lücke schließen oder eine echte Asymmetrie aufdecken | negativtest-deckt-nur-eine-regelseite |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/core/rules/matrix.go` (`tokenFindings`, `tokenFindingsOnLine`, `sourceInstanceID`, `sameInstance`) | geprüft, ohne Befund — `sameInstance` behandelt `loc[2]/loc[3] < 0` (nicht teilgenommene Capture-Gruppe) und Regexe ohne Capture-Gruppe (`len(loc) < 4`) korrekt; `FindAllStringIndex` → `FindAllStringSubmatchIndex` ändert `loc[0]:loc[1]` (Gesamt-Match, für `Target`/Message verwendet) nicht, verifiziert durch grünes `TestMatrixTokenReferenz` (Klassen ohne Capture-Gruppe, unverändert im Diff) |
| `internal/hexagon/core/rules/matrix.go` (Link-Prüfschleife in `CheckMatrix`, `lineageValues`/`lineageExempt`, `matrix-downward`) | geprüft, ohne Befund — Diff berührt ausschließlich die Token-Form-Funktionen; Link-Form, Supersede-Lineage und `matrix-downward` sind textlich unverändert (§1-Abgrenzung des Slice-Plans eingehalten) |
| `internal/adapter/driven/configyaml/configyaml.go` (`applyMatrix`-Zerlegung, `validateMatrixAllowIfSameID`) | geprüft, ohne Befund — `tokenOf[rule.From]`/`tokenOf[rule.To]` sind zum Zeitpunkt des Aufrufs immer gesetzt (Existenz bereits durch die vorgelagerte `classes[...]`-Prüfung erzwungen); Validierung läuft strikt vor `cfg.Matrix.Rules = append(...)`, kein Pfad, auf dem eine `allow-if-same-id: true`-Regel ungeprüft durchrutscht |
| `internal/hexagon/core/model/config.go` (`MatrixRule.AllowIfSameID`) | geprüft, ohne Befund |
| §3.2 Suppression-Verbot | geprüft, ohne Befund — kein `//nolint` im Diff |
| §3.6 Gate-Lockerung nur per ADR | geprüft, ohne Befund — keine Schwellen-Senkung; neue Fähigkeit mit eigener ADR |
| §3.7 Kommentar-Klassen | geprüft, ohne Befund — alle neuen Kommentare (matrix.go, configyaml.go, config.go) tragen Zusage/Kopplung/Grenze mit auflösbarem `DC-FA-MTX-003`-Anker, keine Review-Historie oder Herkunfts-Prosa |
| §3.8 Modul-Scan-Grenze | geprüft, ohne Befund — `sourceInstanceID` liest nur `file` (die bereits im Scan befindliche Quelldatei selbst), keine Eingabe außerhalb der Scan-Menge |
| ADR-Treue (Rule-Feld statt Klassen-Feld, Token-Wiederverwendung, kein neues Klassen-Feld) | geprüft, ohne Befund — Implementierung entspricht Alternative D aus ADR-0087 |
| Byte-Identität ohne `allow-if-same-id` (DC-QA-02) | geprüft, ohne Befund — `sourceInstanceID` wird zwar unabhängig vom Flag berechnet, ihr Ergebnis fließt aber nur über `sameInstance` ein, das ausschließlich bei `rule.AllowIfSameID == true` aufgerufen wird; `TestMatrixAllowIfSameID` bestätigt Byte-Identität zwischen Flag-aus-Konfiguration und einer Baseline-Konfiguration ohne das Feld |
| Akzeptanzkriterien-Deckung „Instanz-Identität Happy/Boundary/Negative (keine Korrelation)/Default" | geprüft, ohne Befund — `TestMatrixAllowIfSameID` und `TestDecode_MatrixAllowIfSameID{Happy,DefaultAus}` decken diese vier wortgetreu (siehe F-2 für die fünfte, Fehlkonfiguration) |
| ADR-0087 Kontext-Abschnitt auf Restbestände der CR-Erzählung (bare `slice-<NNN>`-Fehltreffer beim Absender) | geprüft, ohne Befund — der Kontext-Abschnitt übernimmt die „11 Selbst-Zitate"-Erzählung des CR nicht wörtlich, sondern paraphrasiert nur den strukturellen Befund; kein Leftover-Text aus einem Entwurfsstand |
| Scope-Treue §1 (Link-Form von `matrix-forbidden`, Supersede-Lineage-Ausnahme, Antwort an CR-Absender) | geprüft, ohne Befund — keiner der drei ausgeschlossenen Punkte wurde berührt (`git diff --stat` zeigt keine Änderung an `docs/plan/cr/`, keine neue Datei für eine CR-Antwort) |
| `harness/README.md` §Sensors, `CHANGELOG.md`, `README*.md`, `docs/user/` | geprüft, ohne Befund — keine dieser Dateien im Diff, DoD-Aussage „kein weiterer öffentlicher Vertrag berührt" korrekt |
| Referenz-Richtung (SDP) der neuen `Schärft:`/`Bezug:`-Felder in ADR-0087 | geprüft, ohne Befund — zeigen aufwärts auf `DC-FA-MTX-003`/`DC-FA-MTX-001`/`DC-FA-MTX-001.a`, keine Abwärts-Referenz von Spec auf ADR/Slice |
| Provenance-Marker-Ehrlichkeit (`<!-- d-check:status-provenance -->` im `Schärft:`-Feld von ADR-0087) | geprüft, ohne Befund — zeigt, wo der CR/die Messungen liegen (Provenance), begründet keine Entscheidung |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** adr-fitness-function-names-wrong-gate ·
negativtest-deckt-nur-eine-regelseite

## Verdikt

**Merge-blockierend:** nein — beide Findings sind MEDIUM und dokumentations-
bzw. testabdeckungsbezogen, keine Korrektheitsverstöße im Laufzeitverhalten.
Die Implementierung selbst (Instanz-ID-Extraktion, Capture-Gruppen-Vergleich,
Fail-closed-Validierung, Byte-Identität ohne das Feature) ist gegen ADR-0087
und DC-FA-MTX-003 korrekt und durch grüne Gegenprobe (`make test`, `make
lint`, `make doc-check`) bestätigt. F-1 und F-2 sollten vor der finalen
Closure-Notiz nachgezogen werden (Fitness-Function-Zeile korrigieren bzw.
löschen; vierten Fehlkonfigurations-Unterfall ergänzen), blockieren aber
nicht den Review-Schritt selbst.

**Übergabe:** Findings gehen an den Implementer (Rückkante Review → Plan
nicht nötig — kein Plan-Defekt, reine Nacharbeit am bestehenden Diff). Die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in
den Zähler. Dieser Report selbst ist ein **Lauf-Beleg** und wird über Läufe
hinweg nicht wieder gelesen. Der Report ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat.
