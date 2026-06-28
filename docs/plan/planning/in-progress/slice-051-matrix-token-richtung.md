# Slice slice-051: Modul `matrix` — Token-basierte Referenz-Richtung + Provenance-Marker

**Status:** in-progress (welle-40-matrix-token-richtung).

**Welle:** welle-40-matrix-token-richtung (Trigger: Auftraggeber — die
Referenz-Richtung mechanisieren, die das adoptierte Regelwerk bewusst dem
Reviewer überließ; der Provenance-Marker macht die semantische Unterscheidung
grep-bar).

**Bezug:** Neue Anforderung
[`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
(Lastenheft 0.31.0) +
[ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md). Setzt die
`matrix`-Mechanik aus
[ADR-0021](../../adr/0021-matrix-klasseninterne-verweisrichtung.md) fort; dieselbe
Linie wie
[`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs).

**Autor:** pt9912. **Datum:** 2026-06-28.

---

## 1. Ziel

`matrix` prüft Referenz-Richtungen nur über **Links**. Eine Referenz kann aber
als **bare ID-Token** im Fließtext stehen (eine Slice-Kennung in einem
ADR-Körper), die der Link-Scan nicht sieht — und für Slice-Kennungen gibt es kein
`ids`-Linkpflicht-Muster. Das Regelwerk benennt die Lücke und überlässt sie dem
Reviewer (ein nackter Token-Grep flaggt legitime Verifikations-Zeiger
falsch-positiv). Neu: eine Klasse kann ein `token`-Regex tragen; `matrix` fängt
verbotene Token-Referenzen, ein **Provenance-Marker**
`<!-- d-check:status-provenance -->` deklariert legitime aus → die Semantik wird
grep-bar. Immutable `Accepted`-ADRs werden per `exempt-paths` grandfathered.

## 2. Entscheidungen

- **Token-Erkennung in `matrix`** (nicht eigenes Modul): `matrix` *ist* das
  Referenz-Richtungs-Modul; Link- und Token-Form sind dieselbe Frage. Klasse-Feld
  `token: <regex>`; Scan über den Prosa-Körper außerhalb Fences/`exclude-sections`/
  Links.
- **Provenance-Marker als Ausnahme** (nicht Token-Detektion als Primär-Signal):
  der Marker nimmt eine bereits geflaggte verbotene Token-Referenz aus —
  `d-check:ignore`-Mechanik, aber **benannt** (näher an `allow-supersede-lineage`).
  Erster `matrix`-Zeilen-Marker; die „nur strukturell"-Haltung wird eng begrenzt
  umgekehrt ([ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md)). Links bleiben markerlos.
- **Grandfathering per `exempt-paths`** statt Immutable-Edit/`adr-check`-Eingriff:
  21 `Accepted`-ADRs sind eingefroren und nennen Slices als legitime
  Verifikations-Zeiger im Körper — sie werden ausgenommen (Regelwerk:
  „Gate prüft nur ab Einführung neu"), neue ADRs ab 0022 tragen die Deklaration.
- **`matrix-forbidden` wiederverwendet** (kein neuer Grund-Code): dieselbe
  Verletzung, andere Form (Link/Token).
- **Ehrlichkeit bleibt Reviewer:** der Marker macht „deklariert?" grep-bar, nicht
  „ehrlich?" — der Reviewer-Anker schrumpft auf Marker-Missbrauch-Audit.

## 3. Definition of Done

- [x] **Spec (doc-first):** [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix) (Lastenheft 0.31.0 + §7) +
  [ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md) + ADR-Index
  + [§DC-FA-MTX-001.a](../../../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung) Schritt 6 + Schema-Keys (`matrix.classes[].token`,
  `matrix.exempt-paths`) + §4-`matrix-forbidden`-Token-Form.
- [ ] **Code:** `model.MatrixClass.Token` + `MatrixConfig.ExemptPaths`; Config-Adapter
  parst/validiert fail-closed; `CheckMatrix` — `exempt-paths`-Guard + Token-Pass
  (Prosa, Link-/Fence-frei, Marker-Ausnahme). Tests: Happy (markiert)/Boundary
  (unmarkiert → `matrix-forbidden`)/Negative (Config-Fehler)/grandfathered/
  Default-aus byte-identisch/Token-in-Link-und-Fence-zählt-nicht.
- [ ] **Dogfood:** `.d-check.yml` — `token: 'slice-\d{3}'` auf `slice`-Klasse,
  Regel `{from: adr, to: slice, allow: false}`, `exempt-paths` für die Alt-ADRs 0001–0021.
  Spec-Fix `spezifikation.md` (Slice-Token aus dem Körper). [ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md) dogfoodet den
  Marker auf seinem eigenen `slice-051`-Beleg.
- [ ] **Reviewer-Anker** in `.harness/skills/reviewer.md` von „Referenz-Richtung
  beurteilen" auf „Provenance-Marker-Ehrlichkeit prüfen" abspecken.
- [ ] `make gates` grün; zwei unabhängige Reviews; CHANGELOG; Release v0.31.0;
  Closure (Move + Roadmap-Flip,
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

## 4. Risiken / offene Punkte

- **Negativ-Probe als Beleg:** ein temporärer unmarkierter Slice-Token in einem
  *neuen* (nicht-grandfatherten) ADR-Körper muss `make doc-check` röten; mit
  Marker grün — beweist, dass das Token-Gate feuert und der Marker greift.
- **Dogfood-Disziplin:** Spec-Straten und neue ADRs dürfen Slice-Token im Körper
  nur markiert oder in der Historie tragen — daher `slice-NNN`-Platzhalter in der
  Spec-Prosa und der markierte Beleg in [ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md).
- **Token-in-Link-Doppelzählung** vermieden, indem der Token-Scan Markdown-Links
  vorher entfernt (Links deckt die Link-Prüfung ab).
- **Bekannte LOW-Grenzen (R1, Folge-CR-tauglich):** der Token-Scan leert keinen
  Inline-Code (ein backtick-Slice-Pfad im Körper triggert) und strippt
  Badge-/verschachtelte Links nur partiell; der Marker ist nackter Substring
  (nimmt die ganze Zeile aus). Kein aktueller Repo-Trigger; bewusst minimal
  gehalten. `exempt-paths` überspringt die Datei komplett (breiter als die
  Token-Motivation, in [ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md) begründet). R1-MEDIUM (Fence-/Section-Test-Lücke)
  + F-1 (alle Token je Zeile) sind behoben.

## 5. Trigger

Auftraggeber (2026-06-28): nach der Referenz-Richtungs-Diskussion (Regelwerk
§Referenz-Richtung) — „setzen wir das regelkonform um (Release v0.31.0)", Weg B
(adr→slice auch mechanisch, via Token + Marker + Grandfathering).

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code + Spec; „Doc führt, Code folgt"). Keine BF-Sub-Area.
