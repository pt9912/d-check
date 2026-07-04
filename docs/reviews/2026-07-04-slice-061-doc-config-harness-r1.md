# Review slice-061 — Doku-Config-Beispiel-Harness (Dimension A), R1

## Kopf-Metadaten

- **Gegenstand:** Feat-Commit `0e7d3b4` — `internal/adapter/driven/configyaml/docexamples_test.go`
  (`TestDocExamples_ConfigBeispieleValidieren` + Extraktor `extractYAMLBlocks`)
  und der not-config-Marker in `docs/user/benutzerhandbuch.md:549`.
- **Slice:** slice-061 (welle-50), Dimension A (Config-Fragment ↔ `configyaml.Decode`).
- **Anforderung:** [`DC-FA-CONF-001`](../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
  (Vollvalidierung, jeder Fehler → Exit 2). Kein CR, kein ADR (Test-Harness im bestehenden Schnitt).
- **Reviewer:** unabhängig, Code nicht mitverfasst. **Datum:** 2026-07-04. **Lauf:** R1.
- **Methode:** `git show 0e7d3b4`, Lesen von Plan/Spec/Validator; Gegenprobe des Extraktor-Randes
  gegen die realen Doku-Dateien (Fence-Dump, CRLF-Prüfung, alternative Fence-Varianten,
  Scope-Sweep über `docs/user/*.md` + READMEs). **Keine** `make`-Läufe (Verifier-Rolle).

---

## Findings

### MEDIUM-1 · DC-FA-CONF-001 / Slice-§1-Zielaussage · `internal/adapter/driven/configyaml/docexamples_test.go:16`

`docExampleFiles()` listet nur Handbuch + `README.md` + `README.de.md`.
`docs/user/operations.md` fehlt — obwohl es **zwei echte, kopierbare
`.d-check.yml`-Config-Beispiele** trägt: `ids:`-Muster (`operations.md:66–73`,
inkl. `link-policy`, `exempt-paths`) und `matrix.status:` (`operations.md:108–114`).
operations.md ist Live-Nutzer-Doku: vom Handbuch (`benutzerhandbuch.md:9`) und von
`README.md:169/200` („Invocation reference: options, exit codes, configuration")
verlinkt, nicht in `scan.ignore`, kein `done/`-Slice. Die Slice-§1-Aussage „damit
ist die Blindspot-Klasse für Config-Beispiele geschlossen, nicht nur der
hostpaths-Einzelfall" trifft für diese zwei Blöcke nicht zu; der Plan-§2-Scope
nennt operations.md weder als eingeschlossen noch in seiner „Nicht in
Scope"-Aufzählung (done/, Reviews, vendored Cache, `.d-check.yml`) — es ist
übersehen, nicht bewusst ausgeschlossen.
**Failure-Szenario:** Ein Autor tippt in `operations.md` künftig ein Config-
Beispiel mit unbekanntem Schlüssel oder führendem `/` (genau die
`hostpaths.prefixes`-Klasse aus `c8c33a0`). `configyaml.Decode` würde es beim Lauf
des kopierenden Nutzers mit Exit 2 ablehnen — der Harness bleibt aber grün, weil die
Datei nie extrahiert wird. Der Silent-Grün-Pfad, den der Slice schließen soll,
besteht für operations.md fort. (Kontext-Eskalation: Beobachtung liegt im
Gate-Pfad `make test`/`gates`.)
**Verifizierbar:** ja — ein bewusst kaputtes ` ```yaml `-Beispiel in
`operations.md` hinzufügen, `make test`: bleibt fälschlich grün.

### LOW-1 · Maintainability / Plan-§4 (Fence-Parser-Robustheit) · `internal/adapter/driven/configyaml/docexamples_test.go:79`

`extractYAMLBlocks` führt keinen Zustand „innerhalb eines Nicht-yaml-Fences". Es
öffnet einen Block bei jedem ` ```yaml ` und schließt am nächsten nackten ` ``` `,
ohne umgebende Fences (` ```markdown `, ` ```text `) zu berücksichtigen. Das
Handbuch trägt bereits einen ` ```markdown `-Template-Block (`benutzerhandbuch.md:869–885`);
aktuell enthält er kein inneres ` ```yaml `, daher heute keine Fehl-Erkennung.
Es existiert **kein** Unit-Test für den Extraktor selbst (nur der Lauf gegen die
echten Doku-Dateien), obwohl Plan §4 „Testfälle dafür" (eingerückte Fences,
Info-String-Zusatz, Balance) benennt.
**Failure-Szenario:** Ein Autor fügt zur Illustration ein ` ```yaml `-Beispiel
**innerhalb** eines ` ```markdown `/` ```text `-Blocks ein. Der Extraktor deutet
das innere ` ```yaml ` als Top-Level-Öffner, validiert dessen Inhalt und
überspringt den äußeren Schließer als nackten ` ``` `; ein absichtlich
illustrativ-kaputtes Fragment macht den Test fälschlich rot, ein valides wird
als bestandenes Beispiel fehlgezählt.
**Verifizierbar:** ja — solchen verschachtelten Block einfügen, `make test`:
unerwartete Extraktion/Rotfärbung.

### LOW-2 · Plan-§4 (Autoren-Hinweis) · `docs/user/benutzerhandbuch.md:549`

Der Commit setzt den Marker
`<!-- d-check-test:not-config: … -->`, ergänzt aber **keinen** autoren-sichtbaren
Hinweis in Handbuch/README, der die Marker-Konvention und ihre harte Bedingung
(„Marker muss in der **unmittelbaren** Vorzeile stehen", geprüft per
`strings.Contains(lines[i-1], …)`) erklärt. Plan §4 verlangt ausdrücklich, die
Wartungslast „im Handbuch-Autoren-Hinweis benennen". Die Konvention ist derzeit
nur im Slice-Plan dokumentiert.
**Failure-Szenario:** Ein Autor fügt ein neues ` ```yaml `-Ausgabe-Beispiel ein
oder formatiert das bestehende mit einer Leerzeile zwischen Marker und Fence um
(natürliche Markdown-Spationierung). Der Marker greift nicht mehr (Vorzeile ≠
Marker), `Decode` lehnt den `findings:`-Block ab → roter Test ohne jede
In-Doku-Anleitung, was der Marker ist. (Fail-laut, daher LOW.)
**Verifizierbar:** ja — Leerzeile zwischen `benutzerhandbuch.md:549` und `:550`
einfügen, `make test`: das `findings:`-Beispiel wird rot.

### INFO-1 · Maintainability · `internal/adapter/driven/configyaml/docexamples_test.go:74` / `:36`

Zwei dokumentationswürdige Annahmen ohne aktuellen Trigger:
(a) `skip` ist ein Substring-Test (`strings.Contains`) über die Vorzeile. Landet
künftig ein Autoren-Hinweis-Absatz, der den Token `d-check-test:not-config`
literal nennt, **unmittelbar** vor einem echten Config-` ```yaml `, wird dieses
echte Beispiel **still** übersprungen (Silent-Grün). (b) Der Leere-Menge-Guard
`total == 0` zählt `total++` **vor** der `skip`-Prüfung; ein hypothetischer
Zustand, in dem jeder ` ```yaml ` marker-übersprungen ist, liefe grün mit **null**
tatsächlichen Validierungen durch. Beide geringe Wahrscheinlichkeit, aber
undokumentiert.

### INFO-2 · Maintainability · `internal/adapter/driven/configyaml/docexamples_test.go:64`

Die Fence-Erkennung matcht ausschließlich `"```yaml"` und `"```yaml "` (Space).
Alternative/Attribut-Info-Strings (` ```yaml{.line-numbers} `, ` ```yml `,
` ```YAML `) werden still übersprungen. Heute keine im Repo (Sweep sauber).
**Failure-Szenario (latent):** Ein echtes Config-Beispiel in einer solchen Fence
liefe nie durch `Decode` — Silent-Grün. Undokumentierte Annahme, dass Doku-Autoren
nur die exakte ` ```yaml `-Form nutzen.

---

## Negativbefunde (geprüft, ohne Befund)

- **Extraktor-Korrektheit gegen die 3 In-Scope-Dateien:** alle ` ```yaml `-Fences
  sind Top-Level, sauber balanciert und alternieren; der ` ```markdown `-Block
  @869–885 enthält kein inneres ` ```yaml ` → **keine** Fehl-Erkennung in der
  aktuellen Doku. `openLine = i+1` korrekt (1-basiert, zeigt auf die Fence-Zeile).
- **Marker-Anwendung:** genau der eine erwartete Nicht-Config-Block (das
  `--yaml`-Ausgabe-Beispiel `findings:` @550) ist korrekt, unmittelbar vor dem
  Fence, mit Grund annotiert. Kein weiterer Nicht-Config- ` ```yaml ` in der
  In-Scope-Menge.
- **Fail-closed-Guards:** fehlende Datei (`os.ReadFile`→`Fatalf`), leere Menge
  (`total==0`→`Fatal`), unbalancierter Öffner (`!closed`→`Fatalf`) alle vorhanden
  und fail-closed (Rand-Nuance siehe INFO-1b).
- **Hexagon/Imports/Determinismus (DC-QA-02/03):** externes Testpaket
  `configyaml_test`, Imports nur `os`/`path/filepath`/`strings`/`testing` + der
  eigene Adapter; kein Netz; keine Inline-Suppression; relativer Repo-Pfad
  präzedenzgestützt (`diagnose_test.go:71` liest `../../../../spec/spezifikation.md`,
  slice-060). Tragfähig im Docker-`make test` (Linux, voller Repo-Kontext).
- **Zeilenenden:** Handbuch/README.md/README.de.md sämtlich LF (kein CRLF) — die
  `strings.Split(content, "\n")`-Zeilenzerlegung greift; kein `\r`-Rest an
  ` ```yaml `/` ``` `.
- **README.md/de.md ` ```yaml ` @177:** echtes Config-Beispiel
  (`scan`/`modules`/`ids`), korrekt in Scope, Struktur `Decode`-konform.
- **DC-FA-CONF-001-Bezug:** der Harness übt exakt `configyaml.Decode` (den
  CONF-001-Validator) aus; Mechanik deckt sich mit dem Vertrag, kein neuer Vertrag.

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 2 |

---

## Verdikt

**NACHBESSERN.** Kein HIGH; die Kernmechanik (Extraktion, Marker-Opt-out,
fail-closed-Guards, Determinismus, Hexagon) ist korrekt und für die drei
In-Scope-Dateien belegt sauber. Blockierend ist **MEDIUM-1**: `operations.md`
mit zwei echten, verlinkten Config-Beispielen liegt außerhalb der Doku-Menge,
womit die vom Slice §1 beanspruchte Schließung der Blindspot-Klasse für genau
diese Blöcke nicht erreicht ist und der Silent-Grün-Pfad dort fortbesteht.
Auflösung: entweder operations.md in die Menge aufnehmen **oder** den bewussten
Ausschluss im Plan-§2-Scope begründen (z. B. Verschiebung nach slice-062).
LOW-1/LOW-2 sind vor Closure zu klären (Extraktor-Robustheit ohne Test bzw.
fehlender Autoren-Hinweis, beide von Plan §4 vorgezeichnet); INFO-1/INFO-2 sind
dokumentationswürdige Annahmen (Won't-Fix zulässig, wenn notiert).
