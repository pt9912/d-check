# Verifikations-Report: slice-193 — 2026-09-03

**Verifikations-Art:** DoD-/Spec-Konformität gegen den Slice-Plan (Modul 8
§Verification vs. Validation).

**Gegenstand:** Commits `faac05c` (MR-058 → `done/`), `8bd9ab7`
(feat: Baseline-Pin auf v6.0.0)

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-03

**Eingangs-Kontext:**
[slice-193](../plan/planning/done/slice-193-baseline-v600-bump.md)
§4 Definition of Done, `AGENTS.md` (Hard Rules).

---

## DoD-Punkte

- [x] MR-058 nach `harness/conventions/done/` migriert, eigener Commit —
  ursprünglich nicht sauber trennbar (`harness/conventions.md`s
  Aktive-Adaptionen-Zeile hing an derselben Datei wie der Haupt-Bump);
  behoben durch einen Zwischenzustand (MR-058 bleibt vorübergehend in
  „Aktive Adaptionen", nur ihr Pfad zeigt schon auf `done/`) — `faac05c`
  ist dadurch ein echter, `make gates`-grüner Zwischenstand.
- [x] Alle 15 lebenden `d-check:cite`-Direktiven (11 Dateien, nicht 12 wie
  ursprünglich geschätzt) re-geankert — `make doc-check --enable
  citations` grün.
- [x] Alle übrigen lebenden `.harness/baseline/v5.18.0/`-Pfad-Referenzen
  gezogen — ein übersehener Treffer (`roadmap.md:13`, „§Roadmap-Struktur
  (v5.18.0)") beim ersten Durchlauf gefunden und auf `v6.0.0` korrigiert,
  weil die zitierte Sektion inhaltlich unverändert ist und der Tag damit
  ein Live-Zeiger ist, keine Herkunfts-Prosa.
- [x] `.harness/baseline/v5.18.0/` entfernt, `v6.0.0` einziger vendorter
  Baum.
- [x] `harness/conventions.md` §Baseline `**Stand:**` → v6.0.0 + MR-060.
- [x] `make baseline-verify` und `make baseline-freshness` (Content) grün.
- [x] `make gates` grün (zehn Gates).
- [x] `make fullbuild` grün.
- [x] Unabhängiger Review durchgeführt
  ([Report](2026-09-03-slice-193-baseline-v600-code-r1.md), 1 MEDIUM,
  behoben).
- [x] Unabhängige Verifikation durchgeführt (dieser Report).

## Hard-Rule-Konformität

- §3.3/MR-013: Move-Commit trägt nur eigene Link-Tiefen-Fixes plus den
  einen externen Pointer, der sonst sofort `target-missing` wäre —
  konform mit der MR-/Wellen-Lifecycle-Move-Ausnahme.
- §3.5: keine der drei immutable ADRs (0080/0081/0082) wurde editiert.
- §3.7: der neue `.d-check.yml`-Tombstone-Kommentar wurde vor der Closure
  von Slice-Nummern-Aufzählung und vergleichender Bump-übergreifender
  Prosa bereinigt (Herkunft bleibt auf `MR-*`/`ADR-*`-Kennungen
  beschränkt).

## Verdikt

Alle DoD-Punkte erfüllt, zwei während der Verifikation gefundene Lücken
(Commit-Trennung, `roadmap.md`-Treffer) vor der Closure behoben. Bereit
für `done/`.
