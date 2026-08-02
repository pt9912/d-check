# Slice slice-087: Spec-§7-Referenzrichtung konform (C-3-Nachzug)

**Status:** Done (welle-67, C-3-Nachzug abgeschlossen 2026-08-02).

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

- [x] `.d-check.yml`: `matrix.exclude-sections` = `[Geschichte]` (Historie/7. Historie
  entfernt).
- [x] `spec/spezifikation.md` §7 + `spec/lastenheft.md` §7 ohne `slice-NNN`-Token und
  ohne `docs/plan/adr/`-Abwärtslink; §7 bleibt lesbare Chronik.
- [x] [ADR-0047](../../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md)
  Accepted; Supersede-Verfeinerung zu
  [ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md) dokumentiert;
  das immutable Original bleibt unberührt.
- [x] `make gates` + `make adr-check` grün; unabhängiger Frischkontext-Review.

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

Umgesetzt: die d-check-eigenen Spec-Straten sind in §7 („## 7. Historie")
v5.0.0-konform — **keine Abwärtsverweise auf slice-/ADR-Kennungen, auch nicht in der
Historie**. Konkret: in `spec/spezifikation.md` §7 und `spec/lastenheft.md` §7 die
`Verweis`-Spalte entfernt (trug je eine `slice-NNN`) und aus der Änderungs-Prosa die
**8** (spez) bzw. **17** (last) `slice-NNN`-Token sowie die **9** echten
adr-Abwärtslinks („Begründung in [ADR-XXXX]"-Klauseln) getilgt; §7 bleibt eine
lesbare `Datum | Änderung`-Chronik. In `.d-check.yml` `matrix.exclude-sections` von
`[Historie, "7. Historie", Geschichte]` auf `[Geschichte]` verengt — nur die
immutable ADR-`## Geschichte` bleibt provenance-exempt, §7 wird ab jetzt geprüft
(fail-closed gegen künftige Abwärts-Kennungen). Entscheidungs-Record:
[ADR-0047](../../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md) (Accepted).

**Korrektur der slice-086-Annahme (Kern-Erkenntnis).** slice-086 §3 schnitt C-3 als
**Folge-Produkt-Slice** heraus mit der Begründung, eine konforme Durchsetzung brauche
ein d-check-**Code-Feature** (per-Klasse/-Pfad `exclude-sections`), weil die Messung
zeigte: `matrix.exclude-sections` ist **global**, und global entfernen bräche **65**
Befunde in den immutablen ADR-`## Geschichte`-Abschnitten. Die Messung entfernte
jedoch die **ganze** Liste. Tatsächlich trennen die **Heading-Namen** sauber: die 46
Accepted-ADRs führen durchweg `## Geschichte`, die zwei Spec-Straten `## 7. Historie`
(keine ADR nutzt „Historie", keine Spec-Datei „Geschichte"); der Ausschluss greift per
**exaktem Heading-Klartext**. Also genügte der **chirurgische** Schnitt — nur `Historie`
/`7. Historie` streichen, `Geschichte` behalten — **ohne** Go-Code, ADR-Feature oder
Release. **Lehre: vor „Ausnahme entfernen" auf der richtigen Granularität messen** — die
slice-086-Messung war zu grob (ganze Liste statt namens-selektiv).

**Behalten (korrekt kein Abwärtsverweis):** die DC-Selbstlinks (`spezifikation.md#…`/
`lastenheft.md#…`, in-file bzw. aufwärts), die `conventions.md#mr-…`-Links (Ziel ist
weder ADR noch Slice → matrix-neutral) und die **fiktiven** Beispiel-Kennungen im
Lastenheft (redaktionell, `d-check:ignore`-markiert, kein Link).

**Enforcement live verifiziert:** ein testweise in §7 injiziertes `slice-999` erzeugte
`spec/spezifikation.md:2152 slice-999 matrix-forbidden`; nach Revert wieder grün — die
§7-Prüfung ist real, nicht still übersprungen.

**Review:** unabhängiger Frischkontext-Review
(`docs/reviews/2026-08-02-slice-087-spec-historie-review.md`): **abnahmereif**, HIGH 0 /
MEDIUM 0 / LOW 1 / INFO 1. Beide eingearbeitet — F-1 (LOW): ein durch die
ADR-Klausel-Entfernung hängendes Semikolon in zwei Zellen geheilt; F-2 (INFO): ein
Rest-Abwärtsbezug in Nicht-`slice-NNN`-Form (`(049/052/053)`) entkoppelt. `make gates`
+ `make adr-check` grün; die verfeinerte immutable
[ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md) byte-identisch.

**Anschluss:** **Etappe D (Form-Konformität)** ist der **letzte** offene
Migrations-Slice der welle-67 (die 11 D-Findings aus slice-085 §3).
