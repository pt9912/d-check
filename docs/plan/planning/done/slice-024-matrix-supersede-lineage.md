# Slice slice-024: `matrix` Supersede-Lineage-Carve-out

**Status:** done.

**Welle:** welle-14-matrix-lineage (Trigger: Change Request des
Auftraggebers aus dem Fremd-Repo `grid-gym`, Dialog 2026-06-17).

**Bezug:** [`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
(Change Request 0.14.0 — opt-in Supersede-Lineage-Ausnahme von der
Status-Prüfung),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Abwärtskompatibilität: Default aus ⇒ Befundsatz byte-identisch),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(netzloser Dogfooding-Lauf bleibt grün).

**Autor:** pt9912. **Datum:** 2026-06-17.

---

## 1. Ziel

Die Status-Regel des Moduls `matrix`
(`status.forbidden: [superseded, deprecated]`) flaggt **jede** Referenz
auf ein inaktives Dokument als `matrix-inactive` — auch die legitime
Supersede-Lineage: die ablösende ADR verweist per Definition auf die
ADR, die sie ablöst. Im Melde-Repo `grid-gym` (`.d-check.yml`, v0.10.0)
bricht das das fail-closed `make docs-check`, obwohl der Verweis
normative ADR→ADR-Lineage ist. Der dortige Workaround (Lineage-Link auf
Inline-Code umgestellt) ist informationsärmer als nötig.

Ziel: Eine **opt-in** Konfiguration nimmt genau die deklarierte
Lineage-Kante X → Y von der Status-Prüfung aus, wenn die Quelle X über
ein Feld die Ablösung von Y benennt. Strukturelle Ausnahme (wie
`exclude-sections`), kein zeilenweiser Opt-out-Marker — `matrix` bleibt
marker-frei (Entscheidung des Auftraggebers, „dokumentieren statt
erweitern"; konsistent mit „deterministische Befunde werden behoben,
nicht stummgeschaltet").

## 2. Definition of Done

- [x] **Lastenheft-Change-Request**
  [`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
  (Version-Bump 0.14.0, Historie-Zeile): Beschreibung um die opt-in
  Supersede-Lineage-Ausnahme; drei neue AKs (Lineage Happy/Boundary/
  Default-Negative); Out-of-Scope ergänzt (kein semantischer Schluss,
  nur Status-Prüfung, kein `matrix`-Marker).
- [x] **Spezifikation**
  §[`DC-FA-MTX-001.a`](../../../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung):
  Schritt 4 (Lineage-Ausnahme: Feld-Extraktion `**Feld:** Wert`/
  `Feld: Wert`, Match über normalisierten Linktext bzw. Zielpfad, nur
  `matrix-inactive`); Schema-Tabelle + Beispiel-YAML um
  `allow-supersede-lineage` und `supersede-fields`. **Kein** neuer
  Grund-Code (weiter `matrix-inactive`).
- [x] **Implementierung** im Modul `matrix`: `MatrixConfig` um
  `AllowSupersedeLineage`/`SupersedeFields`; `rawMatrix.Status` +
  Validierung (nicht-leere Feldnamen, Exit 2); `LinkRef.Text` getragen;
  `checkMatrix` ehrt die Lineage-Ausnahme vor dem `matrix-inactive`.
  Die drei AKs als Tests.
- [x] **Abwärtskompatibilitäts-Beleg**
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  Ohne `allow-supersede-lineage` wird `supersede-fields` nie konsultiert;
  Befundsatz byte-identisch (Default-Test).
- [x] **Doku-Nachzug** (`README.md`, `docs/user/operations.md`):
  Lineage-Carve-out dokumentiert; Marker-Abdeckung (`ids`/`codepaths`,
  `matrix` bewusst ausgenommen) explizit.
- [x] `make gates` grün (echte Ausgabe);
  [`CHANGELOG.md`](../../../../CHANGELOG.md); Closure-Notiz §8.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`spec/lastenheft.md`](../../../../spec/lastenheft.md) | update | CR [`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix) 0.14.0: Beschreibung, AKs, Out-of-Scope, Historie |
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | update | §[`DC-FA-MTX-001.a`](../../../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung) Schritt 4, Schema-Tabelle, Beispiel |
| `internal/hexagon/core/config.go` | update | `MatrixConfig` +`AllowSupersedeLineage`/`SupersedeFields` |
| `internal/adapter/driven/configyaml/configyaml.go` | update | `rawMatrix.Status` + Mapping + Validierung |
| `internal/hexagon/core/markdown.go` | update | `LinkRef.Text` (Linktext für den Lineage-Match) |
| `internal/hexagon/core/matrix.go` | update | Lineage-Ausnahme vor `matrix-inactive`; Feld-/Match-Helper |
| `internal/hexagon/core/matrix_test.go` | update | AKs Lineage Happy/Boundary/Default |
| `internal/adapter/driven/configyaml/configyaml_test.go` | update | Parsing der zwei neuen Status-Schlüssel + Validierung |
| `internal/adapter/driving/cli/config_template.go` | update | `--print-config`-Gerüst um die opt-in Keys (Kommentar) |
| [`README.md`](../../../../README.md) | update | Lineage-Carve-out + Marker-Abdeckung |
| [`docs/user/operations.md`](../../../../docs/user/operations.md) | update | Lineage-Carve-out + Marker-Abdeckung |
| [`CHANGELOG.md`](../../../../CHANGELOG.md) | update | 0.x-Eintrag (nutzersichtbar) |

## 4. Trigger

Change Request des Auftraggebers (verbatim) aus `grid-gym`: normative
ADR→ADR-Lineage (`Aenderungstyp: Supersedes …` plus `Bezug:`-Link auf
das abgelöste ADR) wird als `matrix-inactive` gemeldet und bricht das
fail-closed `make docs-check`. Zweiter Befund im CR: der zeilen-scoped
`d-check:ignore`-Marker greift nicht für `matrix`. Entscheidung des
Auftraggebers im Dialog 2026-06-17: **Teil (A)** (Lineage-Carve-out)
umsetzen; **Teil (B)** als „dokumentieren, nicht erweitern" — `matrix`
bleibt marker-frei, die Modul-Abdeckung wird in README/operations
explizit gemacht.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün (echte Ausgabe), Default-Lauf zeigt
unveränderten Befundsatz (Abwärtskompatibilität), Closure-Notiz §8 mit
Lerneintrag.

## 6. Risiken und offene Punkte

- **Match-Breite:** Der Lineage-Match ist ein normalisierter
  Teilzeichenketten-Vergleich (Linktext bzw. Zielpfad gegen den
  Feldwert). Zu breite Feldwerte könnten mehr Kanten ausnehmen als
  beabsichtigt; die Ausnahme bleibt aber auf `matrix-inactive` *einer
  Quelldatei mit passendem Feld* beschränkt und ist opt-in. Im Review
  bestätigen.
- **Abwärtskompatibilität**
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  Default `allow-supersede-lineage: false` ⇒ `supersede-fields` wird nie
  gelesen, kein Pfad-Unterschied. Durch Default-Test belegt.
- **`LinkRef.Text`-Erweiterung:** zusätzliches Feld an einem geteilten
  Typ; rein additiv (bestehende Nutzer ignorieren es). Befüllung in
  `ExtractLinks` über die bereits geparste Linktext-Spanne.
- **Marker-Politik:** `matrix` bleibt bewusst ohne `d-check:ignore`
  (Teil B des CR). Wird in README/operations dokumentiert, damit die
  Modul-Abdeckung des Markers nicht implizit bleibt.

## 7. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas Greenfield (Spec-/Code-/Doku-Arbeit;
Greenfield-Default der Modus-Tabelle in
[`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)).

## 8. Closure-Notiz (nach `done/`)

**Umsetzung:** Vertrag + Spezifikation (Commit `1bd860e`), `matrix`-Code
+ Tests (`9e10e05`), Doku + Slice + Roadmap (`c0228c3`), Lifecycle-Move
(`d7b4c16`). `make gates` grün (echte Ausgabe: doc-check + lint + test +
arch-check + coverage-gate **95,00 %** ≥ 93 % + gate-consistency).

**Vorher/Nachher-Beleg am gebauten Image** (`d-check:latest`,
`--network none`, read-only-Mount): ein ablösendes ADR mit
`**Aenderungstyp:** Supersedes ADR 0003` und Link auf das mit
`Superseded by …` markierte Ziel erzeugte ohne das Flag **1**
`matrix-inactive` (Exit 1); mit `allow-supersede-lineage: true` +
`supersede-fields: [Supersedes, Aenderungstyp]` **0** Befunde (Exit 0).
Boundary bestätigt: eine fremde Quelle ohne Supersede-Feld auf dasselbe
Ziel blieb `matrix-inactive` — der Carve-out ist eng auf die deklarierte
Lineage-Kante beschränkt.

- **Was hat funktioniert:** Die Ausnahme als **strukturelle**
  Konfiguration zu modellieren (wie `exclude-sections`) statt als
  Zeilen-Marker. Der Match lokal in der Quelldatei (Linktext bzw.
  Zielpfad gegen den Feldwert) vermeidet das ID-vs-Pfad-Formatproblem
  über Dateigrenzen — beides ist in X' eigener Schreibweise notiert.
- **Anders als geplant:** Das Tragen des Linktexts erforderte ein neues
  Feld `LinkRef.Text`; zwei bestehende `ExtractLinks`-Tests prüften
  implizit, dass `Text` leer ist, und mussten nachgezogen werden (kein
  Verhaltens-, nur ein Erwartungs-Update).
- **Steering-Loop-Lerneintrag (geschärfte Regel):** Legitime Ausnahmen
  eines deterministischen Gates gehören in **deklarierte Konfiguration**,
  nicht in verstreute Opt-out-Marker. Der Auftraggeber-Vorschlag bot
  beides an (Lineage-Carve-out **oder** `d-check:ignore` für `matrix`);
  gewählt wurde der Carve-out plus die explizite Doku, dass `matrix`
  marker-frei bleibt — so bleibt „deterministische Befunde werden
  behoben oder strukturell ausgenommen, nicht stummgeschaltet" intakt
  ([`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
  0.14.0). Abwärtskompatibel
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  Default aus ⇒ byte-identisch.
- **Folge-Slices:** keine. Das Melde-Repo `grid-gym` (und andere)
  profitieren, sobald sie die Release-Version pinnen und
  `allow-supersede-lineage` aktivieren; der Inline-Code-Workaround dort
  kann auf den klickbaren Lineage-Link zurückgebaut werden.
- **Offen:** Release `v0.11.0` (GHCR-Digest-Pin) folgt nach dem Review
  des Auftraggebers („Review vor Release").
