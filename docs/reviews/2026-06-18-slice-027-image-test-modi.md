# Review — slice-027 (`image-test` deckt `--doctor`/`--repair`)

## Kopf-Metadaten

- **Datum:** 2026-06-18
- **Gegenstand:** `tools/image-test.sh` — neue Stufe (4): `--doctor` und
  `--repair` nativ vs. Container byte-identisch.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Eingangs-Kontext:** Slice-Plan `slice-027-image-test-modi.md`;
  DC-FA-DIST-001 (Image: identisch zur nativen Ausführung, CLI-Optionen
  als Container-Argumente), DC-QA-02 (byte-identisch), DC-FA-CLI-007/008
  (Modi). `make image-test` + `make gates` grün.

## Findings

### INFO-1 — nur konservatives `--repair` in der Image-Stufe

- **Quelle:** DC-FA-DIST-001 · **Pfad:** `tools/image-test.sh` Stufe (4)
- **Befund:** Stufe (4) prüft `--repair` (konservativ), nicht
  `--repair-broad`. Bewusst (im Slice-Plan §6 begründet): die breite Stufe
  bräuchte ein Basisnamen-Treffer-Fixture; der nativ==Container-Byte-
  Vergleich ist von der Stufe unabhängig und mit dem konservativen Lauf
  belegt. Akzeptiert.

## Negativbefunde (geprüft, ohne Befund)

- **Vertrags-Parität (DC-FA-DIST-001):** Stufe (4) belegt am **realen
  Image** „Ergebnis identisch zur nativen Ausführung" jetzt auch für die
  Modi `--doctor`/`--repair` (Exit-Code gleich, stdout **und** stderr
  byte-identisch — DC-QA-02).
- **Aussagekraft:** Fixture trägt einen `id-unlinked`-Befund → `--doctor`
  liefert eine Diagnose, `--repair` einen nicht-leeren Patch (gegen den
  Hunk geprüft); keine Leer-Ausgabe-Tautologie.
- **read-only:** Container-Lauf mit `:ro`-Mount; das Werkzeug schreibt
  nichts (DC-QA-03).
- **Kein Vertrags-Change:** reine Gate-Härtung unter DC-FA-DIST-001; kein
  Lastenheft-/Spezifikations-Eingriff (korrekt — der Vertrag fordert die
  Parität bereits generisch).

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 0 | 1 |

## Verdikt

**Freigegeben.** Reine, am realen Image grüne Gate-Härtung; keine
HIGH/MEDIUM/LOW. Closure kann erfolgen; welle-16 ist damit vollständig und
die E2E-Lücke vor dem Release v0.16.0 geschlossen.
