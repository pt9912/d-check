# Slice slice-148: Baseline-Pin auf v5.12.0 — Etappe A der Migration

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-85-baseline-v5120-migration](../welle-85-baseline-v5120-migration.md).

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011) (Pin auf
Release-Tag, diese Hebung ist ihre Fortschreibung),
[`MR-021`](../../../../harness/conventions.md#mr-021) (in-Repo-Verweise sind
pin-gebunden), [`MR-023`](../../../../harness/conventions.md#mr-023)
(self-contained Bundle-Layout), [`MR-030`](../../../../harness/conventions.md#mr-030)
(der abzulösende Vorgänger-Pin), [`BEO-008`](../observations.md) (drei
Spiegel-Klassen).

**Berührte Spec-Stellen:** — (Harness-Bestand; das Produkt bleibt unberührt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

---

## 1. Ziel

Der vendorte Baum wandert von `v5.11.0` auf `v5.12.0`, und **alle** Verweise
wandern mit. Das Delta ist bereits gemessen (Wellendokument §1): 28 von 52
Dateien unterscheiden sich, davon **fünf** mit echtem Regel-Inhalt.

**Der Inhalt ist nicht Gegenstand dieses Slice** — nur der Stand. Was die fünf
Änderungen für dieses Repo bedeuten, beantwortet
[slice-149](../open/slice-149-baseline-v5120-delta-audit.md).

## 2. Vorgehen

1. Bundle vendoren und **verifizieren** (`--verify`, Manifest über beide Bäume).
2. **Je Spiegel-Klasse aus [`BEO-008`](../observations.md) eine Zählung vor und
   nach der Hebung:** Pfad-Verweise (gate-gedeckt), Release-/Tree-URLs,
   Prosa-/Ellipsen-Pins. Die Zahl kommt in die Closure-Notiz, nicht nur der
   grüne Lauf.
3. **Die vierte Klasse prüfen, die `BEO-008` nicht führt:** ein Verweis, der
   nicht nur auf eine Datei zeigt, sondern deren **Wortlaut zitiert**. Bei
   Punkt 2 des CR hat sich genau dieser Wortlaut geändert
   ([`MR-033`](../../../../harness/conventions.md#mr-033) zitiert die alte
   Fassung). Ob das eine eigene Klasse ist, gehört gemessen und benannt.
4. Ein neuer Adaptions-Eintrag trägt die Pin-Hebung; [`MR-030`](../../../../harness/conventions.md#mr-030)
   nach `conventions/done/` samt Link-Tiefen-Fix im Move-Commit
   ([`MR-013`](../../../../harness/conventions.md#mr-013)).
5. Alt-Baum `v5.11.0` entfernen; eingefrorene Zitate in `done/` bleiben, ihr
   Pfad geht bei Bedarf ins Quell-skopierte Ventil.
6. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine inhaltliche Folge.** Was die fünf geänderten Dateien für dieses Repo
  heißen, ist Etappe B.
- **Kein Freshness-Audit der Adaptionen.** Ebenfalls Etappe B.
- **Kein Retrofit eingefrorener Dokumente.** `done/` und Review-Reporte
  zitieren den Stand ihrer Zeit.

## 4. Definition of Done

- [ ] Bundle vendored und verifiziert; Datei-Zahl genannt.
- [ ] **Je Spiegel-Klasse eine Zahl** vor und nach — auch für die ungedeckten.
- [ ] Die Zitat-Klasse ist gemessen und beantwortet: eigene Klasse oder nicht.
- [ ] Der neue Adaptions-Eintrag ist geschrieben, [`MR-030`](../../../../harness/conventions.md#mr-030)
      aufgelöst und verschoben.
- [ ] Alt-Baum entfernt; kein hängender Verweis (`make doc-check` Exit 0).
- [ ] `make gates` grün (Exit explizit), `make fullbuild` grün; unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Pfad-`grep` kennt keine Zeitform.** [`BEO-008`](../observations.md)
  nennt die Über-Hebung ausdrücklich: eine Aussage über die **Vergangenheit**
  darf nicht mitgehoben werden. Jede Fundstelle ist auf Gegenwart oder
  Historie zu prüfen. — **Ausgang:** *(bei Closure)*
- **Ein gehobener Link kann auf einen geänderten Wortlaut zeigen.** Der Pfad
  löst auf, das Zitat daneben stimmt nicht mehr — und kein Gate sieht das. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): Welle eröffnet.

**Rückführungen:** `in-progress` → `next`, falls die Verifikation des Bundles
scheitert.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Bestand (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-008`](../observations.md) ist der Anlass;
  [`BEO-011`](../observations.md) für jede Aussage darüber, dass eine
  Spiegel-Klasse „vollständig" gehoben sei.

Slice-ID: slice-148. Betroffene IDs: — (kein `DC-`-Bezug). Module:
Harness-Bestand. Gates: `make doc-check`, `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Fortschreibung eines etablierten Vorgangs.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
