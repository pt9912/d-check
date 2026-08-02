# ADR-0047 — matrix: Spec-Straten-Historie nicht mehr provenance-exempt (v5.0.0-Referenz-Richtung)

**Status:** Proposed
**Datum:** 2026-08-02
**Autor:** pt9912
**Bezug:** [`DC-FA-MTX-001`](../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
(Modul `matrix`),
[`DC-FA-MTX-003`](../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix);
mechanisiert die §Referenz-Richtung des **v5.0.0**-Regelwerks
([`grundlagen-referenz-richtung.md`](../../../.harness/baseline/v5.0.0/regelwerk/grundlagen-referenz-richtung.md#referenz-richtung-sdp-wer-darf-wen-referenzieren)),
die die Historie-Ausnahme streicht. **Verfeinert** die Entscheidung aus
[ADR-0022](0022-matrix-token-richtung-provenance-marker.md) (dort:
„`exclude-sections` deckt Provenance unter `## Geschichte`/Historie ab") — ADR-0022
bleibt immutabel unberührt, nur der Geltungsbereich seiner
`exclude-sections`-Aussage wird für die Spec-Straten verengt.
**Schärft:** — (kein Spec-Schema-Delta; ein Config-Wert in `.d-check.yml` +
§7-Content zweier Spec-Straten).

## Kontext

Der adoptierte Baseline-Pin ist auf **v5.0.0** gehoben (welle-67). Deren
Referenz-Richtung verbietet Abwärtsverweise aus den Spec-Straten auf ADR/Slice
**auch in der Historie** — die frühere Provenance-Ausnahme für Historie-Abschnitte
ist widerrufen (Baseline-Migrations-Befund C-3/B-5/F-8). d-checks eigene Spec
verletzt das: `spec/spezifikation.md` §7 und `spec/lastenheft.md` §7 tragen eine
`Datum | Änderung | Verweis`-Chronik, deren `Verweis`-Spalte und Prosa `slice-NNN`
(und einige `docs/plan/adr/`-Links) nennen.

[ADR-0022](0022-matrix-token-richtung-provenance-marker.md) stützte einen Teil
seiner Begründung darauf, dass „`exclude-sections` Provenance unter
`## Geschichte`/Historie abdeckt". Eine frühere Messung schloss daraus, eine
konforme Durchsetzung brauche ein **per-Klasse**-`exclude-sections`-Code-Feature
(global entfernen bräche 65 Befunde in den immutablen ADR-`Geschichte`-Abschnitten).
Diese Annahme ist **falsch**: die Abschnitts-Namen trennen bereits sauber — die 46
Accepted-ADRs führen durchweg `## Geschichte`, die zwei Spec-Straten `## 7.
Historie` (keine ADR nutzt „Historie", keine Spec-Datei „Geschichte"). Der
Ausschluss greift per **exaktem Heading-Klartext**, ist also schon heute
namens-selektiv.

## Entscheidung

1. **`matrix.exclude-sections` auf `[Geschichte]` verengen** — `Historie` und
   `7. Historie` entfallen. Das prüft ab sofort die Spec-§7-Historie und lässt die
   immutable ADR-`Geschichte` unberührt exempt.
2. **Spec §7 ent-tokenisieren** (beide Straten): die `Verweis`-Spalte streichen und
   die `slice-NNN`-Token sowie echte `docs/plan/adr/`-Abwärtslinks aus der
   `Änderung`-Prosa entfernen. §7 bleibt eine lesbare `Datum | Änderung`-Chronik;
   die Rückwärts-Traceability liegt ohnehin in der RTM (`trace`) und der
   git-Historie, nicht im Top-Stratum.
3. **`## Geschichte` bleibt exempt** — die marker-lose Slice-Provenance in den
   immutablen ADRs ist per [ADR-0016](0016-adr-immutable-gate.md) unveränderbar; die
   vor Einführung Accepted-ADRs bleiben zusätzlich per `matrix.exempt-paths`
   grandfathered.

Die drei Schritte sind eine reine **Baseline-Konformität** — kein Go-Code, kein
neues Config-Schema, kein Release. Der Baseline-Default sticht die repo-lokale
Historie-Ausnahme.

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **Heading-Namens-Verengung (gewählt)** | kein Code; trifft exakt Spec §7; immutable ADR-`Geschichte` unberührt; nutzt die vorhandene exakte Heading-Semantik | §7 verliert die Slice-Chronik-Kopplung (liegt aber in RTM/git) |
| Per-Klasse `exclude-sections` als Code-Feature (die zunächst gemessene Annahme) | granular pro Klasse | Go-Code + eigene ADR + Release für einen Effekt, den die Namens-Trennung gratis liefert |
| §7 ganz auslagern (nach `CHANGELOG.md`/git) | strukturell sauberste Trennung | größerer Eingriff; die Chronik verlässt die Spec |
| Historie-Ausnahme behalten | kein Aufwand | verstößt gegen die adoptierte v5.0.0-Baseline |

## Fitness Function

- Nach Verengung + §7-Putz: `make doc-check` (Modul `matrix`) **grün** — kein
  `slice-NNN`-Token und kein `docs/plan/adr/`-Abwärtslink mehr in den beiden
  Spec-§7-Abschnitten.
- Die **46** ADR-`## Geschichte`-Abschnitte bleiben befundfrei (`Geschichte` weiter
  in `exclude-sections`).
- `make adr-check` **grün**: keine Accepted-ADR inhaltlich berührt; ADR-0022 bleibt
  byte-identisch.
- Dogfood: `matrix.exclude-sections: [Geschichte]` in [`.d-check.yml`](../../../.d-check.yml).

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-08-02 | Entwurf (slice-087, welle-67 „C-3-Nachzug"): Korrektur der slice-086-Annahme (kein per-Klasse-Code-Feature nötig — die Heading-Namen `Geschichte`/`7. Historie` trennen ADR- und Spec-Provenance schon selektiv). Status Proposed. |
