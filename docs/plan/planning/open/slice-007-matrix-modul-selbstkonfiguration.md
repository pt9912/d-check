# Slice slice-007: Modul `matrix` + Dogfooding-Selbstkonfiguration

**Status:** open.

**Welle:** welle-03-regelmodule.

**Bezug:** [`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix),
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(Selbstkonfiguration),
[`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)
(maschinelle Kodierung);
[ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md).

**Autor:** pt9912. **Datum:** 2026-06-10.

---

## 1. Ziel

Das Regelmodul `matrix` ist implementiert, und die eigene
`.d-check.yml` aktiviert `ids` + `matrix` mit der vollständigen
Selbstkonfiguration — die MR-006-Referenzrichtungs-Regel
(Spec-Straten verweisen nie abwärts) wird damit maschinell erzwungen.

## 2. Definition of Done

- [ ] Akzeptanzkriterien von `DC-FA-MTX-001` als Tests: Slice →
  aktives ADR ok; Referenz auf `Superseded …` → `matrix-inactive`;
  Lastenheft → ADR → `matrix-forbidden` mit beiden Klassen.
- [ ] Status-Extraktion in fester Reihenfolge (`**Status:**` vor
  `Status`-Heading), Präfix-Match case-insensitiv; ohne Status aktiv
  (Spezifikation §DC-FA-MTX-001.a).
- [ ] Klassen-Glob-Präzedenz (Deklarationsreihenfolge) und
  `exclude-sections` (getrimmt, ohne Auszeichnung, case-sensitiv)
  getestet.
- [ ] Selbstkonfiguration aktiv: `.d-check.yml` deklariert
  Dokumentklassen (Spec-Straten, ADR, Slice), Regeln
  `{spec-straten → adr/slice: verboten}`, `status.forbidden`,
  `exclude-sections` (Historie/Geschichte) sowie `ids`-Muster für
  `DC-*`, `MR-*`, `ADR-*`; `make doc-check` läuft mit
  `links, anchors, ids, matrix` und ist auf dem eigenen Repo grün.
- [ ] `matrix` in `isImplemented`; `make gates` grün;
  [`CHANGELOG.md`](../../../../CHANGELOG.md); Closure-Notiz.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/core/matrix.go` (+ Tests) | neu | Klassen-/Regel-/Status-Logik |
| `internal/hexagon/core/config.go`, `configyaml` | update | Matrix-Konfig in `core.Config` durchreichen |
| [`.d-check.yml`](../../../../.d-check.yml) | update | Selbstkonfiguration (ids + matrix) |
| Eigene Doku (lebende Dateien) | update | etwaige Befunde der Selbstkonfiguration bereinigen |

## 4. Trigger

slice-006 done (Selbstkonfiguration braucht das `ids`-Modul).

## 5. Closure-Trigger

DoD vollständig + Commit(s) auf `main` + Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Nackte Kennungen in **immutablen** Texten (ADR-Geschichte,
  akzeptierte MR-Einträge): Inline-Code ist linkpflichtfrei und
  Historie-Sektionen sind per `exclude-sections` ausnehmbar — falls
  darüber hinaus Befunde in unveränderlichen Passagen bleiben, ist
  eine dokumentierte Entscheidung nötig (Spez-Fortschreibung einer
  Sektions-Ausnahme für `ids` ODER Form-Fix analog MR-007); nicht
  stillschweigend lockern.
- Klassen-Zuordnung über Globs muss mit der `"."`-Wurzel des
  Dogfoodings zusammenspielen (Pfade relativ zur Repo-Wurzel).

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (spec-first; siehe Kurs Modul 5 §Worked
Mini-Example).
