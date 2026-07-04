# Slice slice-061: Config-Beispiel-Verifikation (Fenced-YAML ↔ Validator)

**Status:** in-progress (welle-50-handbuch-beispiel-verifikation).

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

- [ ] **Harness:** Go-Test, der die Doku-Menge liest, je ` ```yaml `-Block
  extrahiert (Fence-Balance geprüft), die per Marker ausgenommenen überspringt
  und den Rest gegen `configyaml.Decode` validiert; jede Ablehnung ⇒ Testfehler
  mit Datei:Zeile + Validator-Meldung.
- [ ] **Marker-Anwendung:** die identifizierten Nicht-Config-Blöcke (Stand
  Sondierung: das `--yaml`-Ausgabe-Beispiel) tragen den beschlossenen
  Skip-Marker mit Grund.
- [ ] **Fail-closed-Guards** (fehlende Datei / keine Fence / unbalancierter
  Fence / leere Menge) mutations-verifiziert (slice-057-R3-Lehre: jede Probe
  verriegelt genau ihren Guard).
- [ ] **Regressions-Beleg:** ein bewusst kaputtes Config-Beispiel (z. B.
  `hostpaths.prefixes: ["/x"]` temporär) ⇒ Harness rot; Revert ⇒ grün.
- [ ] **Belege/Prozess:** `make gates`/`make ci` grün; unabhängiges Review vor
  Closure; Closure-Move nach `done/` + Roadmap-Flip
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

*(bei Closure: Umsetzung, Belege, Lerneintrag, Steering-Loop-Eintrag.)*
