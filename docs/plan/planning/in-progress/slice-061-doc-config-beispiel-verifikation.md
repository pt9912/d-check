# Slice slice-061: Handbuch-Beispiel-Verifikation (Config-Parse + E2E-Verhalten)

**Status:** in-progress (welle-50-handbuch-beispiel-verifikation).

**Welle:** welle-50-handbuch-beispiel-verifikation (Trigger: Nutzer-Frage
2026-07-04 „Haben wir keine Tests dafür geschrieben?" nach dem
`hostpaths.prefixes`-Doku-Fix `c8c33a0` + Nutzer-Ergänzung „von E2E-Tests
könnte ein Handbuch profitieren, da ja im Handbuch Beispiele aufgeführt sind").

**Bezug:** Verifikations-Mechanik gegen bereits bestehende Verträge
([`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
Config-Vollvalidierung;
[`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)
Exit-Codes;
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
Image-Verhalten). **Kein Change Request** (kein neuer/geänderter Vertrag),
**kein ADR** (Test-Harness im bestehenden Schnitt). Schwester-Muster zur
slice-060-Verriegelung `AllReasons` ↔ Spezifikation §4 (Code ↔ Doc), hier
**Handbuch-Beispiel ↔ Realität**.

**Autor:** pt9912. **Datum:** 2026-07-04.

---

## 1. Ziel

Das Benutzerhandbuch führt Beispiele in drei Klassen, von denen **keine**
heute an die Realität gekoppelt ist:

1. **Config-Fragmente** (` ```yaml `, `.d-check.yml`) — laufen nie durch die
   eigene `configyaml`-Validierung (Fenced-Code wird nicht gescannt); real
   gegen den Validator gelaufen beim `hostpaths.prefixes: ["/home"]`-Beispiel
   (führender `/` ⇒ Exit 2, Fix `c8c33a0`).
2. **Kommando-Aufrufe mit dokumentierter Ausgabe/Exit-Code** (` ```bash ` +
   ` ```text `) — behaupten ein Verhalten (sauberes Repo ⇒ Exit 0 +
   „0 Befund(e)"; kaputter Link ⇒ Befund-Zeile `Datei:Zeile  Ziel  Grund` +
   Exit 1; `--doctor` ⇒ gruppierte Diagnose), das nie gegen das echte Binary
   geprüft wird.
3. **Reine Ausgabe-Beispiele** (` ```text `/` ```json `/` ```yaml findings `) —
   Struktur-Behauptungen ohne Kopplung an die reale Ausgabe.

**Ziel:** die überprüfbaren Behauptungen der Beispiele an die Realität binden —
Config-Fragmente gegen `configyaml.Decode` (Parse-Ebene) und
Kommando-/Ausgabe-Beispiele gegen das **echte Binary** über Fixtures
(E2E-/Verhaltens-Ebene). Damit ist die Doku-Drift-Blindspot-Klasse geschlossen,
nicht nur der `hostpaths`-Einzelfall.

## 2. Entscheidungen

Zwei Dimensionen, ein Ziel. **Stufung ist der zentrale offene Entscheid**
(siehe §4): Dimension A ist der enge, direkte Nachzug zum `hostpaths`-Bug;
Dimension B (E2E) hat spürbar mehr Design-Fläche und wird ggf. als **slice-062**
abgespalten.

### Dimension A — Config-Fragmente ↔ Validator (Parse)

- **Test, kein Gate-Target.** Go-Test im bestehenden `make test`-Lauf, der die
  Doku-Menge über einen repo-relativen Pfad liest (Präzedenz: slice-060 liest
  `../../../../spec/spezifikation.md`).
- **Inklusionsregel (welcher Block wird validiert):** nicht jeder
  ` ```yaml `-Block ist Config-Input — Handbuch-Block @549 ist ein
  `--yaml`-**Ausgabe**-Beispiel (`findings:`), das `Decode` zu Recht ablehnt.
  Drei Optionen:

  | Ansatz | Pro | Contra |
  |---|---|---|
  | **(A) Opt-out-Marker** — alle ` ```yaml `-Blöcke validieren; Nicht-Config-Block trägt Skip-Marker (HTML-Kommentar, mit Grund) | fail-closed **zum Prüfen hin**: neues echtes Beispiel automatisch abgedeckt; neuer Nicht-Config-Block bricht **laut** bis annotiert | Marker an den (wenigen) Ausgabe-Beispielen nötig |
  | (B) Opt-in-Marker — nur markierte Blöcke | kein False-Positive | neues echtes Beispiel **still** unabgedeckt (dieselbe Silent-Drift-Falle) |
  | (C) Schlüssel-Heuristik — nur bekannte Top-Level-Keys | markerlos | getippter Key (genau der `unbekannter Schlüssel`-Bug) still übersprungen |

  **Empfehlung (A)** — deckungsgleich mit der slice-060-/`d-check:ignore`-
  Philosophie (fail-closed, laut). Marker-Vorschlag:
  `<!-- d-check-test:not-config: <Grund> -->` vor dem Fence. Genau **ein**
  erwarteter Fall (das `--yaml`-Ausgabe-Beispiel @549).
- **Partielle Configs sind gültig** (`Decode` akzeptiert Teil-Configs, alle
  Felder optional — `TestDecode_LeereUndKommentarDatei`).

### Dimension B — Kommando-/Ausgabe-Beispiele ↔ echtes Binary (E2E)

- **Vorhandene E2E-Infrastruktur nutzen:** `tools/image-test.sh`
  ([`make image-test`](../../../../harness/README.md#sensors-feedback-gates))
  fährt schon nativ == Container mit Fixtures und prüft Exit-Codes/Ausgabe;
  `cli_acceptance_test.go` fährt CLI-E2E gegen `MemFS`/Temp-Repos. Dimension B
  reiht sich dort ein statt einen neuen Runner zu bauen.
- **Nur Beispiele mit prüfbarer Verhaltensbehauptung** werden E2E-verankert:
  je ein Fixture, das die Prämisse des Beispiels herstellt (sauberes Repo /
  kaputter Link / `--doctor`-Fall), der dokumentierte Aufruf (die **Flags**,
  nicht der wörtliche `docker run …@sha256`-Pull), Assertion auf Exit-Code +
  Ausgabe-**Form** (Befund-Zeilen-Schema, „N Befund(e)"-Zeile, Diagnose-Kopf).
- **Nicht replaybare Beispiele markieren:** Aufrufe, die externen Zustand
  referenzieren (fremde Repos, konkrete Digests, Netz) sind nicht E2E-fähig —
  gleicher Opt-out-Marker wie Dimension A, mit Grund.
- **Version-/Digest-Zeilen sind kein E2E-Ziel:** die `:vX.Y.Z`-/`@sha256:`-
  Pins in den Beispielen wandern legitim pro Release (der `versions`-Gate
  deckt die ghcr-Pins bereits) — E2E prüft **Verhalten**, nicht die
  Pin-Strings.

### Gemeinsam

- **Scope der Doku-Menge:** Benutzerhandbuch (Config in §4/§5/§6, Kommandos in
  §3/§4) + die je eine ` ```yaml `-Fence in `README.md`/`README.de.md`. Nicht
  in Scope: `done/`-Slices, Review-Reports, der vendored Regelwerk-/Templates-
  Cache, `.d-check.yml` selbst (läuft bereits jeden Lauf durch den Validator).
- **Fail-closed:** fehlende Doku-Datei, keine erwartete Fence oder
  unbalancierter Fence ⇒ Test rot (kein stilles Grün mit leerer Menge —
  Guard-Klasse aus slice-060/057).
- **Determinismus/Read-only** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  rein Go, ohne Netz (Dimension B nutzt das lokal gebaute Image wie
  `image-test`, kein Pull), identischer Baum ⇒ identisches Ergebnis.

## 3. Definition of Done

**Dimension A (in jedem Fall):**

- [ ] Go-Harness liest die Doku-Menge, extrahiert je ` ```yaml `-Block
  (Fence-Balance geprüft), überspringt Marker-ausgenommene, validiert den Rest
  gegen `configyaml.Decode`; jede Ablehnung ⇒ Testfehler mit Datei:Zeile +
  Validator-Meldung.
- [ ] identifizierte Nicht-Config-Blöcke (Sondierung: `--yaml`-Ausgabe @549)
  tragen den beschlossenen Skip-Marker mit Grund.
- [ ] Fail-closed-Guards mutations-verifiziert (jede Probe verriegelt genau
  ihren Guard); adversariale Probe (temporär kaputtes Config-Beispiel ⇒ rot,
  Revert ⇒ grün).

**Dimension B (dieser Slice oder slice-062 — s. §4):**

- [ ] repräsentative Handbuch-Kommando-Beispiele mit Verhaltensbehauptung als
  E2E-Fälle (Fixture + Flags + Exit-Code-/Ausgabe-Form-Assertion), eingereiht
  in `cli_acceptance_test.go` bzw. `tools/image-test.sh`.
- [ ] nicht-replaybare Beispiele markiert; Auswahl-Begründung dokumentiert
  (kein stiller Ausschluss — welche Beispiele sind E2E-verankert, welche nicht
  und warum).

**Gemeinsam:**

- [ ] `make gates`/`make ci` grün; unabhängiges Review vor Closure;
  Closure-Move nach `done/` + Roadmap-Flip
  ([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)),
  Closure-Body als Folge-Commit. **Kein Produkt-Code, Image byte-identisch,
  kein Release** (Test-/Harness-Infra, analog slice-055).

## 4. Risiken / offene Punkte

- **Stufung A vs. A+B (zentraler Entscheid):** Dimension A ist der enge
  Nachzug zum `hostpaths`-Bug und in sich abschließbar; Dimension B (E2E) hat
  mehr Design-Fläche (Fixture-Entwurf, Beispiel-Auswahl, Ausgabe-Form-Matching
  ohne die wandernden Pin-Zeilen). Optionen: **beide in slice-061**, oder **A
  jetzt / B als slice-062**. Vor Implementierung mit dem Auftraggeber
  festlegen.
- **Marker-Grundsatzentscheid (§2 A):** Opt-out (A) empfohlen; HTML-Kommentar-
  Marker (wie `d-check:ignore`) wird vom `hostpaths`-/`ids`-Scan nicht als
  Befund gewertet — prüfen.
- **E2E-Ausgabe-Matching:** auf **Form** prüfen (Regex/Schema der Befund-Zeile,
  „N Befund(e)"), nicht auf wörtliche Zeilen mit Datei-Zahlen/Versions-Pins —
  sonst bricht der Test bei jeder harmlosen Doku-/Release-Änderung (Wartungs-
  falle statt Wert).
- **README-YAML-Natur:** die je eine ` ```yaml `-Fence in README.md/de.md vor
  Implementierung sichten (Config vs. Illustration).
- **Wartungslast:** ein neues Ausgabe-Beispiel bricht den Test bis annotiert —
  bewusst (laut > still); im Handbuch-Autoren-Hinweis benennen.

## 5. Trigger

Nutzer-Frage 2026-07-04 nach dem `hostpaths.prefixes`-Doku-Fix (`c8c33a0`) +
Validator-Unit-Test (`22300a3`) — die grundsätzliche zweite Stufe („ein Test,
der die Doku-Beispiele gegen die Realität prüft") — plus die Nutzer-Ergänzung,
dass das Handbuch mit seinen **Beispielen** von **E2E-Tests** profitiert
(Dimension B). Notiert in [[release-prep-doc-currency-blindspots]] Punkt 5.

## 6. Sub-Area-Modus-Begründung

GF (Test-/Harness-Erweiterung im bestehenden Schnitt; „Doc führt, Code folgt"
trägt über die Verifikations-Mechanik — der Test koppelt die Doku-Beispiele an
die bestehenden Verträge/das reale Binary). Kein neuer Adapter, keine
BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

*(bei Closure: Umsetzung, Belege, Lerneintrag, Steering-Loop-Eintrag.)*
