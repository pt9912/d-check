# Slice slice-061: Config-Beispiel-Verifikation (Fenced-YAML ↔ Validator)

**Status:** done (welle-50-handbuch-beispiel-verifikation, Closure 2026-07-04).

**Welle:** welle-50-handbuch-beispiel-verifikation (Trigger: Nutzer-Frage
2026-07-04 „Haben wir keine Tests dafür geschrieben?" nach dem
`hostpaths.prefixes`-Doku-Fix `c8c33a0`).

**Bezug:** Verifikations-Mechanik gegen den bestehenden Vertrag
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(Config-Vollvalidierung, jeder Fehler → Exit 2). **Kein Change Request** (kein
neuer/geänderter Vertrag), **kein ADR** (Test-Harness im bestehenden Schnitt).
Schwester-Muster zur slice-060-Verriegelung `AllReasons` ↔ Spezifikation §4
(Code ↔ Doc), hier **Doc-Config-Beispiel ↔ Validator-Code**.

**Stufung (Nutzer-Entscheid 2026-07-04):** Dieser Slice deckt **Dimension A**
(Config-Fragmente ↔ Parse). Die vom Auftraggeber angeregte **Dimension B**
(E2E-Verankerung der Kommando-/Ausgabe-Beispiele gegen das echte Binary) ist in
den Folge-Slice `slice-062` ausgegliedert (Backlog unter `next/`), weil sie
eigenständige Fixture-/Ausgabe-Matching-Entscheidungen trägt.

**Autor:** pt9912. **Datum:** 2026-07-04.

---

## 1. Ziel

Die gefenceten ` ```yaml `-Config-Beispiele der Live-Nutzer-Doku
(Benutzerhandbuch, README) laufen heute **nie** durch die eigene
`configyaml`-Validierung — d-check scannt Fenced-Code nicht (Ausnahme:
`versions`/`pins` gescopt). Ein Beispiel kann darum dauerhaft gegen den
Validator laufen, ohne dass ein Gate es fängt — real geschehen beim
`hostpaths.prefixes: ["/home"]`-Beispiel (führender `/` ⇒ Exit 2,
Fix `c8c33a0`). **Neu:** ein Test-Harness extrahiert die Config-Beispiele aus
der Live-Doku und validiert jedes gegen `configyaml.Decode`; ein Beispiel, das
der Validator ablehnt (und nicht bewusst als Nicht-Config markiert ist), macht
den Test rot. Damit ist die Blindspot-Klasse für Config-Beispiele geschlossen,
nicht nur der `hostpaths`-Einzelfall.

## 2. Entscheidungen

- **Test, kein Gate-Target.** Go-Test im bestehenden `make test`-Lauf (Teil von
  `gates`), der die Doku-Dateien über einen repo-relativen Pfad liest —
  Präzedenz: die slice-060-Verriegelung liest
  `../../../../spec/spezifikation.md`. Kein neues Make-Target.
- **Zentrale Entscheidung — Inklusionsregel (welcher Block wird validiert):**
  Die Sondierung zeigt, dass **nicht jeder ` ```yaml `-Block ein Config-Input**
  ist — Handbuch-Block @549 ist ein `--yaml`-**Ausgabe**-Beispiel (`findings:`),
  das `Decode` zu Recht ablehnt. Drei Optionen:

  | Ansatz | Pro | Contra |
  |---|---|---|
  | **(A) Opt-out-Marker** — alle ` ```yaml `-Blöcke der Doku-Menge validieren; ein Nicht-Config-Block trägt einen expliziten Skip-Marker (HTML-Kommentar, mit Grund) | fail-closed **zum Prüfen hin**: neues echtes Config-Beispiel automatisch abgedeckt; neuer Nicht-Config-Block bricht **laut** bis annotiert | Marker an den (wenigen) Ausgabe-Beispielen nötig |
  | (B) Opt-in-Marker — nur markierte Blöcke | kein False-Positive | neues echtes Beispiel **still** unabgedeckt (dieselbe Silent-Drift-Falle) |
  | (C) Schlüssel-Heuristik — nur bekannte Top-Level-Keys | markerlos | getippter Key (genau der `unbekannter Schlüssel`-Bug) still übersprungen |

  **Empfehlung (A):** deckungsgleich mit der slice-060-/`d-check:ignore`-
  Philosophie (fail-closed, laut). Marker-Vorschlag:
  `<!-- d-check-test:not-config: <Grund> -->` unmittelbar vor dem Fence. Genau
  **ein** Marker-Fall aus der Sondierung erwartet (das `--yaml`-Ausgabe-
  Beispiel @549); die `--json`-/Text-Ausgaben sind ` ```text `/` ```json `,
  keine ` ```yaml `.
- **Scope der Doku-Menge:** Benutzerhandbuch (§4/§5/§6 tragen die
  Config-Beispiele) **und** die je eine ` ```yaml `-Fence in `README.md`/
  `README.de.md` (vor Implementierung Natur sichten: Config vs. Illustration).
  Nicht in Scope: `done/`-Slices, Review-Reports, der vendored Regelwerk-/
  Templates-Cache (Fremdinhalt, ohnehin `scan.ignore`), `.d-check.yml` selbst
  (läuft bereits jeden Lauf durch den Validator).
- **Partielle Configs sind gültig.** `configyaml.Decode` akzeptiert
  Teil-Configs (alle Felder optional — belegt durch
  `TestDecode_LeereUndKommentarDatei`); ein Einzel-Modul-Fragment wie
  `hostpaths:\n  prefixes: [home]` validiert sauber. Der Harness erwartet
  **kein** Vollausbau-Beispiel.
- **Fail-closed.** Fehlt eine Doku-Datei, eine erwartete ` ```yaml `-Fence oder
  ist ein Fence unbalanciert ⇒ Test rot (kein stilles Grün mit leerer Menge —
  Guard-Klasse aus slice-060/057).
- **Determinismus/Read-only** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  rein Go, ohne Netz; identischer Baum ⇒ identisches Ergebnis.

## 3. Definition of Done

- [x] **Harness:** Go-Test, der die Doku-Menge liest, je ` ```yaml `-Block
  extrahiert (Fence-Zustandsverfolgung), die per Marker ausgenommenen
  überspringt und den Rest gegen `configyaml.Decode` validiert; jede Ablehnung
  ⇒ Testfehler mit Datei:Zeile + Validator-Meldung.
- [x] **Marker-Anwendung:** der eine Nicht-Config-Block (das
  `--yaml`-Ausgabe-Beispiel `findings:`) trägt den Skip-Marker mit Grund.
- [x] **Fail-closed-Guards** (fehlende Datei / unbalancierter Fence / leere
  Menge) mutations-verifiziert; zusätzlich `TestExtractYAMLBlocks` (Extraktion,
  yaml-in-markdown, Marker-Vorzeilen-Regel, unbalanciert, Info-Varianten).
- [x] **Regressions-Beleg:** kaputtes Config-Beispiel ⇒ Harness rot an der
  Zeile (Handbuch-`hostpaths` **und** operations.md-`ids`), Revert ⇒ grün.
- [x] **Belege/Prozess:** `make gates`/`make ci` grün; unabhängiges Review R1
  (NACHBESSERN → alle Befunde eingearbeitet); Closure-Move nach `done/` +
  Roadmap-Flip
  ([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)),
  Closure-Body als Folge-Commit. **Kein Produkt-Code, Image byte-identisch,
  kein Release** (Test-/Harness-Infra, analog slice-055).

## 4. Risiken / offene Punkte

- **Marker-Grundsatzentscheid (§2):** Opt-out (A) empfohlen; der HTML-Kommentar-
  Marker (wie `d-check:ignore`) wird vom `hostpaths`-/`ids`-Scan nicht als
  Befund gewertet — prüfen.
- **README-YAML-Natur:** die je eine ` ```yaml `-Fence in README.md/de.md vor
  Implementierung sichten (Config-Beispiel validieren oder Illustration
  markieren).
- **Fence-Parser-Robustheit:** eingerückte Fences, ` ```yaml `-Info-Strings mit
  Zusatz; der Extraktor muss die Balance korrekt führen (Testfälle dafür).
- **Wartungslast:** ein neues Ausgabe-Beispiel in ` ```yaml ` (selten) bricht
  den Test bis annotiert — bewusst (laut > still); im Handbuch-Autoren-Hinweis
  benennen.

## 5. Trigger

Nutzer-Frage 2026-07-04 nach dem `hostpaths.prefixes`-Doku-Fix (`c8c33a0`) +
Validator-Unit-Test (`22300a3`): die grundsätzliche Schließung der Blindspot-
Klasse für Config-Beispiele. Die E2E-Ausbaustufe (Kommando-/Ausgabe-Beispiele)
ist als `slice-062` ausgegliedert. Notiert in
[[release-prep-doc-currency-blindspots]] Punkt 5.

## 6. Sub-Area-Modus-Begründung

GF (Test-/Harness-Erweiterung im bestehenden Schnitt; „Doc führt, Code folgt"
trägt über die Verifikations-Mechanik — der Test koppelt die Doku-Beispiele an
den bestehenden Validator-Vertrag). Kein neuer Adapter, keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** Dimension A als Go-Test
`TestDocExamples_ConfigBeispieleValidieren`
(`internal/adapter/driven/configyaml/docexamples_test.go`): er liest die
wartungsaktive Nutzer-Doku (Benutzerhandbuch, `operations.md`, `README.md`/
`README.de.md`) über repo-relative Pfade (Präzedenz slice-060), extrahiert jeden
` ```yaml `-Block und validiert ihn gegen `configyaml.Decode`. Der Extraktor
`extractYAMLBlocks` führt echte **Fence-Zustandsverfolgung** (ein Fence öffnet
auf ` ``` `, schließt erst auf ein nacktes ` ``` `) — ein ` ```yaml ` innerhalb
eines ` ```markdown `-Blocks ist damit Body, kein eigener Block; als reine
Funktion (Fehler statt `t.Fatalf`) ist er über `TestExtractYAMLBlocks`
unit-testbar. Nicht-Config-Blöcke tragen den Opt-out-Marker
`<!-- d-check-test:not-config: <Grund> -->` in der Zeile unmittelbar vor dem
Fence (fail-closed zum Prüfen hin: neues echtes Beispiel automatisch abgedeckt,
neuer Nicht-Config-Block bricht laut). Der **Audit fand nur einen** Nicht-Config-
Block (das `--yaml`-Ausgabe-Beispiel `findings:`) und **keinen neuen echten
Config-Bug** (der `hostpaths`-Fall war mit `c8c33a0` bereits behoben). Dimension
B (E2E-Verankerung der Kommando-/Ausgabe-Beispiele) ist als
[`slice-062`](slice-062-handbuch-e2e-beispiele.md) ausgegliedert.
Commit-Kette: doc-first `df07a03` (+ Verengung `46a0127`) → feat `0e7d3b4` → R1
`<review>` → closure-move → closure-body.

**Belege.**
- `make gates`/`make ci` **grün** (doc-check 183/0, lint, test, arch-check,
  coverage, semgrep, gate-consistency, planning-check; image-test).
- **Vier Fail-closed-Mutations-Belege** (jede verriegelt genau ihren Guard):
  kaputtes Config-Beispiel ⇒ rot an der Zeile (Handbuch-`hostpaths` **und**,
  nach R1, operations.md-`ids` via Unbekannter-Key-Injektion an
  `operations.md:66`); unbalancierter Fence / fehlende Datei / leere Menge ⇒
  fail-closed Fatal. Plus `TestExtractYAMLBlocks` (Extraktion, yaml-in-markdown,
  Marker-Vorzeilen-Regel, unbalanciert, `yml`/`YAML`/Info-Suffix).
- **Unabhängiges Review R1**
  ([Report](../../../reviews/2026-07-04-slice-061-doc-config-harness-r1.md)):
  **NACHBESSERN**, 0 HIGH/1 MEDIUM/2 LOW/2 INFO — alle eingearbeitet. MEDIUM-1
  (operations.md trug zwei ungeprüfte Config-Beispiele) aufgenommen; LOW-1
  (Extraktor-Robustheit) per Fence-Zustandsverfolgung + Unit-Tests behoben;
  LOW-2/INFO-1a (Vorzeilen-Regel) in Fehlermeldung/Doc; INFO-1b (`validated>0`
  statt `total>0`) und INFO-2 (Info-String-Varianten) geschlossen.
- **Kein Release** (nur Test-Code + ein HTML-Kommentar-Marker; Image
  byte-identisch).

**Lerneintrag.** (1) Der R1-MEDIUM bestätigt: der **Scope einer Doku-Menge ist
selbst ein Silent-Green-Risiko** — die Blindspot-Klasse ist erst geschlossen,
wenn *alle* config-tragenden Live-Docs erfasst sind (operations.md fehlte
zuerst). Ein Scope-Entscheid gehört explizit begründet, nicht implizit gelassen.
(2) Ein naiver Fence-Extraktor („öffne auf ` ```yaml `") ist bei gemischten
Fences (` ```markdown `-Beispiele) falsch; korrekt ist **Fence-Zustand führen**
und am Info-String entscheiden. (3) Die Schwester-Beziehung zu slice-060 hält:
beide koppeln Doku an Code über eine Mengen-/Format-Verriegelung mit
Fail-closed-Guards — hier Doku-Beispiel ↔ Validator statt Grund-Code-Liste ↔
Spec-§4. Steering-Loop: Feedback (Doku-Beispiel bricht Validator) als
Sensor-Verankerung verkörpert, nicht nur einmalig behoben. Dimension B (E2E)
bleibt der offene Folge-Schritt ([`slice-062`](slice-062-handbuch-e2e-beispiele.md)).
