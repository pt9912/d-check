# Slice slice-087: Spec-§7-Referenzrichtung konform (C-3-Nachzug)

**Status:** In Arbeit (welle-67).

**Welle:** welle-67-baseline-v500-migration (C-3-Nachzug, nach
[slice-086](../done/slice-086-etappe-c-mr-bereinigung.md)).

**Bezug:** Umsetzung der in [slice-085](../done/slice-085-etappe-b-modul-delta.md)
§3 als **C-3/B-5/F-8** doppelt belegten **Historie-Provenance-Revocation** (v5.0.0
`grundlagen-referenz-richtung`: Spec-Straten verweisen nie abwärts auf ADR/Slice —
**auch nicht in der Historie**), die in
[slice-086](../done/slice-086-etappe-c-mr-bereinigung.md) §3 (Abnahme-Punkt 3) aus
der Doc-Migration **herausgeschnitten** wurde. **Korrigiert** dessen Annahme „braucht
ein d-check-Code-Feature": die Messung zeigte, die Heading-Namen `## Geschichte` (46
ADRs) und `## 7. Historie` (2 Spec-Straten) trennen ADR- und Spec-Provenance
**bereits sauber** → ein chirurgischer Konfig-Schnitt genügt. Der Entscheidungs-Record
ist [ADR-0047](../../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md)
(Supersede-Verfeinerung zu
[ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md)).
**Kein Release, keine Go-Code-Änderung.**

**Autor:** pt9912. **Datum:** 2026-08-02.

---

## 1. Ziel

`spec/spezifikation.md` §7 und `spec/lastenheft.md` §7 v5.0.0-konform machen:
**keine Abwärtsverweise auf `slice-NNN`/ADR** — auch nicht in der Historie. Die
immutable ADR-`Geschichte` (marker-lose Provenance) bleibt **unberührt exempt**.

## 2. Vorgehen

1. **Entscheidungs-Record** — [ADR-0047](../../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md)
   (Proposed): Historie-Provenance-Ausnahme für Spec-Straten widerrufen;
   `matrix.exclude-sections` behält nur `Geschichte`; §7-Chronik ent-tokenisiert
   (gewählte Form: §7 bleibt `Datum | Änderung`-Chronik). Bezug
   [`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)/[`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix).
2. **Spec §7 putzen** (beide Straten, ~55/54 Zeilen): die `Verweis`-Spalte streichen;
   die `slice-NNN`-Token und echte `docs/plan/adr/`-Abwärtslinks aus der
   `Änderung`-Prosa entfernen. §7 bleibt lesbare `Datum | Änderung`-Chronik. Die
   fiktiven Beispiel-Kennungen im Lastenheft (redaktionelle Meta-Kommentare, **keine**
   Links) bleiben — `matrix` flaggt bare ADR-Token nicht (kein `adr`-Token, nur
   Links + `slice`-Token).
3. **Konfig verengen** — `Historie` und `7. Historie` aus `matrix.exclude-sections`
   in `.d-check.yml` streichen (`Geschichte` bleibt). **Atomar mit Schritt 2**
   committen — sonst flaggt `matrix` §7 zwischen den Commits.
4. **Gate** — `make gates` (matrix/doc-check/planning-check) + `make adr-check` grün;
   **keine** Accepted-ADR berührt; unabhängiger Frischkontext-Review.

## 3. Abnahme-Punkte

1. **§7-Konformitätsform** — entkoppeln (§7 bleibt Chronik) vs. §7 auslagern.
   → **Entschieden 2026-08-02: entkoppeln** — `Verweis`-Spalte + Prosa-Token raus,
   §7 bleibt `Datum | Änderung`-Chronik.

## 4. Definition of Done

- [ ] `.d-check.yml`: `matrix.exclude-sections` = `[Geschichte]` (Historie/7. Historie
  entfernt).
- [ ] `spec/spezifikation.md` §7 + `spec/lastenheft.md` §7 ohne `slice-NNN`-Token und
  ohne `docs/plan/adr/`-Abwärtslink; §7 bleibt lesbare Chronik.
- [ ] [ADR-0047](../../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md)
  Accepted; Supersede-Verfeinerung zu
  [ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md) dokumentiert;
  das immutable Original bleibt unberührt.
- [ ] `make gates` + `make adr-check` grün; unabhängiger Frischkontext-Review.

## 5. Risiken / offene Punkte

- **Kanonische Spec** (§7 beider Straten): review-pflichtig; nur die §7-Chronik, keine
  Anforderungs-Semantik.
- **Vollständigkeit des Putzes**: der `matrix`-Gate nach der Konfig-Verengung ist die
  Verifikation — ein übersehener Token bricht den Guard (fail-closed, gewollt).
- **Zukunft**: neue §7-Einträge dürfen keine Abwärts-Kennungen tragen — als
  Prozess-Notiz (AGENTS/Handbuch) prüfen, ggf. Etappe D.

## 6. Trigger

[slice-086](../done/slice-086-etappe-c-mr-bereinigung.md) §3 Abnahme-Punkt 3 (C-3
herausgeschnitten) + die Korrektur „kein Code-Feature".

## 7. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc/Prozess führt. Berührt die *Spec*-Sub-Area (§7 beider Straten)
und *Harness/Konfig* (`.d-check.yml`); greenfield-Konformität an die adoptierte
Baseline, ohne Brownfield-Spec.

## 8. Closure-Notiz (nach `done/`)

_Ausstehend._
