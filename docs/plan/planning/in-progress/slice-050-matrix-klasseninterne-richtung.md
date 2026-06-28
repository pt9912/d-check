# Slice slice-050: Modul `matrix` — Verweisrichtung innerhalb einer geordneten Klasse

**Status:** in-progress (welle-39-matrix-richtung).

**Welle:** welle-39-matrix-richtung (Trigger: Auftraggeber — der Wunsch, die
Source-Precedence-Richtung `architecture → spezifikation → lastenheft` **innerhalb**
des Spec-Stratums maschinell zu prüfen; die naive Einzelklassen-Variante war als
Richtung nicht erkennbar und feuerte wegen First-Match nicht).

**Bezug:** Neue Anforderung
[`DC-FA-MTX-002`](../../../../spec/lastenheft.md#dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix)
(Lastenheft 0.30.0) + begleitender
[ADR-0021](../../adr/0021-matrix-klasseninterne-verweisrichtung.md) (Design:
`order`/`direction`, Glob-Rang, fail-closed). Additiv zu
[`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix);
kodiert intern dieselbe Richtung wie
[`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs).

**Autor:** pt9912. **Datum:** 2026-06-28.

---

## 1. Ziel

`matrix` prüft Referenzrichtungen **zwischen** Klassen (`{from,to,allow}`). Die
Schichtung **innerhalb** eines Stratums — `lastenheft` (autoritativ) über
`spezifikation` über `architecture`, Verweise nur aufwärts — war nicht
ausdrückbar, weil alle drei Spec-Dateien in **einer** Klasse `spec-straten`
liegen und eine klasseninterne Referenz keine Regel hat. Neu: eine Klasse trägt
optional `order` (Glob-Rangliste, autoritativste Schicht zuerst) + `direction:
no-downward`; ein klasseninterner Verweis von höher- auf niederrangig erzeugt
`matrix-downward`. So wird die Precedence-Richtung im eigenen Repo und bei
Konsumenten (d-migrate: 23 Spec-Dateien ⇒ Glob-Schichten) prüfbar — an **einer**
lesbaren Stelle statt über *n·(n−1)/2* flache Paare.

## 2. Entscheidungen

- **Klassen-Feld statt Einzelklassen.** Die Gruppe `spec-straten` bleibt intakt
  (Lesbarkeit + bestehende adr/slice-Regeln) und gewinnt `order`/`direction`. Die
  Alternative (Auflösen in Einzelklassen) bricht zweifach: First-Match
  verschattet die Klassen → tote Regeln (stilles Grün), und die Richtung ist über
  Paare nicht erkennbar ([ADR-0021](../../adr/0021-matrix-klasseninterne-verweisrichtung.md)
  §Alternativen).
- **Rang über Globs, First-Match.** `order`-Einträge sind Pfad-Globs; Rang = Index
  des ersten Treffers (wie die Klassenzuordnung). Generalisiert auf viele Dateien
  je Schicht ohne Einzeldatei-Listing. Rangfreie Mitglieder nehmen nicht teil.
- **Transitiv automatisch.** Da der Rang verglichen wird (nicht Paare), ist der
  direkte Sprung `lastenheft → architecture` ohne Extra-Regel erfasst.
- **Fail-closed-Config.** Unbekannter `direction`-Wert, `direction` ohne `order`,
  `order` ohne `direction` ⇒ Konfigurationsfehler (Exit 2). Keine still
  wirkungslose Richtungs-Deklaration — die direkte Lehre aus dem toten-Regel-Bruch.
- **Default-aus byte-identisch.** Ohne `order`/`direction` ist der Befundsatz
  unverändert ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- **Nur `no-downward` in v1.** Andere Politiken (`no-upward`, `strict-adjacent`)
  sind spätere CRs.

## 3. Definition of Done

- [x] **Spec (doc-first):** neue Anforderung [`DC-FA-MTX-002`](../../../../spec/lastenheft.md#dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix) im Lastenheft
  (§3-MTX existiert, Versions-Bump 0.30.0 + §7-Historie) + begleitender
  [ADR-0021](../../adr/0021-matrix-klasseninterne-verweisrichtung.md) + ADR-Index
  + spezifikation [§DC-FA-MTX-001.a](../../../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung) Schritt 5 + Schema-Keys
  (`matrix.classes[].order`/`.direction`) + Grund-Code `matrix-downward` (§4) +
  Config-Beispiel.
- [x] **Code:** `model.MatrixClass` um `Order`/`Direction`; Config-Adapter
  (`applyMatrix`) parst + validiert fail-closed; `CheckMatrix` prüft die
  klasseninterne Rangrichtung (`matrix-downward`); Doctor-Klartext-Mapping ergänzt.
  Tests: Happy (aufwärts ok) / Boundary (abwärts, auch transitiv → `matrix-downward`) /
  Negative (Config-Fehler) / rangfrei / Gleichrang / Selbstverweis / Default-aus byte-identisch.
- [x] **Nutzersichtbare Config-Ausgaben:** `--suggest-config` (`suggest.go`) und
  `--print-config` (`config_template.go`) zeigen `order`/`direction` an der
  geschichteten Klasse; Benutzerhandbuch §4.7 dokumentiert `matrix-downward`.
- [x] **Dogfood:** `spec-straten` in [`.d-check.yml`](../../../../.d-check.yml)
  trägt `order`/`direction`; die toten Einzelklassen/Regeln (Sackgasse der
  additiven Variante) raus.
- [ ] `make gates` grün; CHANGELOG; Closure (Move nach `done/` + Roadmap-Flip,
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

## 4. Risiken / offene Punkte

- **Rang-Semantik ist ein Vertrag** (First-Match-Glob, rangfrei = ausgenommen) —
  exakt in der `.a` dokumentiert; konservativ (nur klasseninterne, beidseitig
  rangbehaftete Kanten).
- **`matrix-downward` vs. `matrix-forbidden`** sind unabhängige Prüfungen
  (Spezifikation [§DC-FA-MTX-001.a](../../../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung): „unabhängig von matrix-forbidden/-inactive").
  Solange die Klassen-Paar-`rules` klassenübergreifend bleiben, kann eine Kante
  nicht beide auslösen; deklarierte man zusätzlich eine klasseninterne Selbstregel
  `{from: X, to: X, allow: false}`, feuerten an einer Abwärtskante beide — kein
  Defekt, sondern zwei verschiedene Verletzungen (Annahme aus dem Impl-Review
  expliziert).
- **Negativ-Probe als Beleg:** ein temporär eingebauter Abwärtsverweis in einem
  Spec-Dokument muss `make doc-check` röten (verifiziert, dass die Regel feuert —
  genau der Test, den die tote Variante nicht bestanden hätte).

## 5. Trigger

Auftraggeber (2026-06-28): die `spec-straten`-Richtung intern prüfbar machen;
nach Design-Diskussion (Transparenz: kein stilles Grün; Erkennbarkeit: kein
Paar-Listing; Generalität: Glob-Schichten für vielfiles-`spec/`).

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code + Spec; „Doc führt, Code folgt"). Keine BF-Sub-Area.
