# Slice slice-027: image-test deckt `--doctor`/`--repair` ab

**Status:** done

**Welle:** welle-16-image-test-modi (Trigger: Entscheidung Auftraggeber —
E2E-Lücke vor Release v0.12.0 schließen; baut auf slice-025/026).

**Bezug:**
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(Image: „Ergebnis und Exit-Code identisch zur nativen Ausführung; CLI-
Optionen als Container-Argumente" — Hauptvertrag),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(byte-identisch),
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
und
[`DC-FA-CLI-008`](../../../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch)
(die getesteten Ausgabe-Modi).

**Autor:** pt9912. **Datum:** 2026-06-18.

---

## 1. Ziel

`make image-test` ([`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image))
prüft bisher nur den Default-Modus nativ == Container. Diese Slice dehnt
den Distributions-Vertrag **explizit auf die Ausgabe-Modi `--doctor` und
`--repair` aus**: je eine Stufe, die den Container-Lauf gegen den nativen
byte-identisch vergleicht (stdout + stderr + Exit-Code). **Reine Gate-
Härtung/Beleg** — kein Vertrags-Change, da [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image) „identisch zur
nativen Ausführung" bereits generisch fordert (CLI-Optionen als
Container-Argumente).

## 2. Definition of Done

- [ ] `tools/image-test.sh` um zwei Stufen erweitert — `--doctor` und
  `--repair` (konservativ) —, je nativ vs. Container byte-identisch
  (stdout **und** stderr) bei gleichem Exit-Code, auf einem Fixture mit
  reparierbarem Befund (`id-unlinked` über eine `ids`-Config).
- [ ] `make image-test` grün (lokal gebautes Image).
- [ ] Kein Lastenheft-/Spezifikations-Change (Gate-Härtung unter
  [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image));
  fällt beim Bauen doch eine Schärfung auf, separat.
- [ ] `make gates` grün (Doku); Closure-Notiz.

## 3. Plan (vor Code)

| Datei | Art | Begründung |
|---|---|---|
| `tools/image-test.sh` | update | Stufen (4) `--doctor` und (5) `--repair`: nativ vs. Container byte-identisch ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)-Parität für die Modi) |

Lastenheft/Spezifikation unverändert (kein neuer Vertrag).

## 4. Trigger

Entscheidung Auftraggeber (E2E-Lücke vor Release v0.12.0). Voraussetzung
erfüllt: slice-025/026 done (`--doctor`/`--repair` existieren).

## 5. Closure-Trigger

DoD vollständig inkl. `make image-test` grün und `make gates` grün.

## 6. Risiken und offene Punkte

- **Breite Stufe im Container:** `--repair-broad` bräuchte ein Fixture mit
  eindeutigem Basisnamen-Treffer; die **konservative** Stufe (`--repair`)
  genügt für den Byte-Vergleich und hält das Fixture minimal.
- **Fixture muss einen reparierbaren Befund tragen** (`id-unlinked` mit
  `ids`-Pattern + existierendem Target-Verzeichnis), sonst liefern
  `--doctor`/`--repair` leere, wenig aussagekräftige Ausgabe.
- `make image-test` baut das Image und läuft nicht in `make gates` —
  Verifikation explizit über `make image-test`/`make ci`.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** `tools/image-test.sh` Stufe (4): auf einem `id-unlinked`-
Fixture werden `--doctor` und `--repair` (konservativ) nativ vs. Container
byte-identisch verglichen (stdout + stderr + Exit-Code), mit Inhalts-
Asserts (Diagnose-Ausgabe bzw. nicht-leerer Patch-Hunk). Kein Lastenheft-/
Spezifikations-Change — reine Gate-Härtung unter
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
/ [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus).

**Belege:** `make image-test` grün (Stufe (4) am realen Image);
`make gates` grün. Verifikation explizit über `make image-test` (läuft
nicht in `make gates`).

**Lerneintrag:** Der Distributions-Vertrag „identisch zur nativen
Ausführung" galt für die neuen Modi nur transitiv (gleiches Binary); die
Stufe macht ihn am ausgelieferten Image *explizit* prüfbar. Das Fixture
trägt bewusst einen reparierbaren Befund — sonst wäre der Byte-Vergleich
eine Leer-Ausgabe-Tautologie.

**Review R1** (Self-Review,
[Report](../../../reviews/2026-06-18-slice-027-image-test-modi.md)):
HIGH 0 / MEDIUM 0 / LOW 0 / INFO 1 (konservatives `--repair` genügt für
die Stufe) — freigegeben.

**Welle:** welle-16-image-test-modi damit vollständig; die E2E-Lücke vor
dem Release v0.12.0 ist geschlossen.

## 8. Sub-Area-Modus-Begründung

Sub-Area `Test-Infrastruktur` (GF, Greenfield-Default): Erweiterung des
bestehenden E2E-Gates ohne Bestandscode-Inventur.
