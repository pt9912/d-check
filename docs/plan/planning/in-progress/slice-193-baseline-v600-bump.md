# Slice slice-193: Die Baseline steht auf v6.0.0

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-88](../welle-88-baseline-v600-migration.md) — erster von
drei Slices: reiner Pin-Bump, die Register-Neugestaltung ist ausdrücklich
nicht Gegenstand (§3).

**Bezug:**
[`MR-021`](../../../../harness/conventions.md#mr-021) (In-Repo-Verweise sind
pin-gebunden — die Pfad-Hälfte des Bumps);
[`MR-051`](../../../../harness/conventions.md#mr-051) (die `cite`-Spannen,
die zweite pin-gebundene Größe);
[`MR-055`](../../../../harness/conventions.md#mr-055) (die Symlink-Aliase,
die denselben Pin binden);
[`MR-013`](../../../../harness/conventions.md#mr-013) (Lifecycle-Move-Commit
bündelt gekoppelte Verweise — Gegenstand der [`MR-058`](../../../../harness/conventions.md#mr-058)→`done/`-Migration);
[`MR-058`](../../../../harness/conventions.md#mr-058) (dokumentiert die
vorige Hebung, v5.15.0→v5.18.0, wird mit diesem Slice selbst superseded und
wandert nach `done/`).

**Berührte Spec-Stellen:** — (Adoptions-Stand einer externen Konvention;
keine Produkt-Anforderung berührt).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-03.

---

## 1. Ziel

**Der vendorte Baseline-Baum steht auf `v5.15.0`↔`v5.18.0` … Zielstand
`v6.0.0`.** Gemessen (netzloser Diff der vendorten Bäume `v5.18.0` gegen
`v6.0.0`, alle 53 Bundle-Dateien klassifiziert): **14** Dateien mit echtem
Inhalt, Rest reiner Versions-Stempel. Zwei inhaltlich unabhängige Stränge:

1. **Wellenlose Zeitdokumente-Archivierung** (Kurs-Wellen 112–114, bereits
   über die eigene CR-Antwort
   [`2026-09-03-antwort-ai-harness-course-wellenlose-review-archivierung.md`](../../cr/2026-09-03-antwort-ai-harness-course-wellenlose-review-archivierung.md)
   bekannt): `modul-05`, `modul-06`, `modul-10` bekommen eine sechste
   Trägerzeile — jeder wellenlos geschlossene Slice archiviert sich **selbst**,
   nach den drei Paarungen, Schlüssel `done/slice-<NNN>-archiv.zip` flach
   neben dem Stub. **Wird mit diesem Slice adoptiert** (Regelwerk-Ebene);
   die Umsetzung in `tools/archive-wave` (ein neuer Einzel-Slice-Modus) ist
   ein eigener Folge-Slice, kein Teil dieses Bumps (§3).
2. **Beobachtungs-Register neu gestaltet** (Kurs-Wellen 115–116):
   `grundlagen-begriffe.md`, `grundlagen-harness-dateien.md`,
   `grundlagen-traceability.md`, `modul-05`, `modul-06` — Tabelle wird
   Verzeichnis, `BEO-<NNN>` wird Pfad `BEO-<KUERZEL>/<slug>`, Zähler wird
   abgeleitet statt gepflegt. **Wird mit diesem Slice NICHT adoptiert** —
   das ist slice-194 (Fähigkeit) und slice-195 (Migration), siehe §3.

Bis slice-194/195 schließen, gilt für den eigenen Bestand weiter die alte
Form (`observations.md`, `BEO-<NNN>`) — dieser Slice hebt nur den Pin und
alle davon unabhängigen Pfad-/Zitat-Verweise; er ändert an `observations.md`
selbst nichts.

## 2. Vorgehen

1. **Vendoring** (bereits erledigt, netzlos verifiziert):
   `.harness/baseline/v6.0.0/{regelwerk,templates}` via
   `tools/harness/fetch-baseline-cache.sh v6.0.0`, 53/53 Dateien OK.
2. **[`MR-058`](../../../../harness/conventions.md#mr-058) → `done/` als eigener, vorgezogener Commit** — Lehre aus
   slice-189s eigener Closure-Notiz (F-3 des unabhängigen Reviews: das
   Ein-Commit-Muster aus MR-Move + Pfad-Bump ist seit mehreren Bumps
   ungewächtert). Diesmal: `git mv
   harness/conventions/MR-058-baseline-v5180.md
   harness/conventions/done/MR-058-baseline-v5180.md` samt eigener
   Link-Tiefen-Fixes, **vor** dem Haupt-Bump-Commit, nicht danach.
3. **Cite-Re-Anker ([`MR-051`](../../../../harness/conventions.md#mr-051)):** alle lebenden `d-check:cite`-Direktiven auf
   `.harness/baseline/v5.18.0/…` — 12 Dateien, Pfad-Bump auf `v6.0.0`; davon
   4 Zitate in `modul-05-planning-harness.md` mit Zeilenspann-Verschiebung
   (+4 ab der Einfügestelle, nachgerechnet gegen den Volltext, nicht
   geschätzt). Frozen-Bestand (`done/`, `docs/reviews/`,
   `harness/conventions/done/`) bleibt unangetastet (`citations.scope`).
4. **Pfad-Bump aller übrigen lebenden `.harness/baseline/v5.18.0/…`-Referenzen**
   auf `v6.0.0` (Markdown-Links ohne Cite-Partner, [`MR-021`](../../../../harness/conventions.md#mr-021)s eigenes
   Ist-Beispiel „aktuell `…/v5.18.0/…`").
5. **Adaptions-Review** über alle lebenden `MR`-Einträge (Stand: Anzahl beim
   Ausführen erheben, nicht schätzen) gegen ihr `Ersetzt-Baseline-Regel`-Feld.
6. **`.harness/baseline/v5.18.0/` entfernen**, `v6.0.0` bleibt als einziger
   vendorter Baum.
7. `harness/conventions.md` §Baseline `**Stand:**`-Zeile auf `v6.0.0`.

## 3. Ausdrücklich NICHT in diesem Slice

- **Die Beobachtungs-Register-Neugestaltung** — Architektur, Produktcode und
  Datenmigration sind slice-194/195. `observations.md` bleibt in diesem
  Slice byte-identisch (bis auf ggf. den einen Cite-Pfad-Bump ihrer eigenen
  Baseline-Verweise).
- **Die Umsetzung der wellenlosen Archivierung in `tools/archive-wave`** —
  die Regelwerk-Adoption ist Teil dieses Slice, ein neuer
  Einzel-Slice-Archivierungsmodus im Werkzeug ist ein eigener,
  unverbindlicher Folge-Slice (wie schon bei welle-86s eigenem
  Stop-Hook-Kandidaten).
- Jede Prosa-Currency-Änderung an `docs/user/*` — dieser Bump berührt kein
  Nutzer-sichtbares CLI-Verhalten.

## 4. Definition of Done

- [ ] [`MR-058`](../../../../harness/conventions.md#mr-058) nach `harness/conventions/done/` migriert, eigener Commit,
      `make doc-check` grün auf dem Zwischenstand.
- [ ] Alle 15 lebenden `d-check:cite`-Direktiven (11 Dateien) auf `v6.0.0` re-geankert
      (Pfad + ggf. Zeilenspanne), `make doc-check --enable citations` grün.
- [ ] Alle übrigen lebenden `.harness/baseline/v5.18.0/`-Pfad-Referenzen auf
      `v6.0.0` gezogen (0 Treffer bei `grep -rn "\.harness/baseline/v5\.18\.0"`
      außerhalb von `done/`, `docs/reviews/`, `harness/conventions/done/`,
      `docs/plan/adr/` [immutable ADRs], `docs/plan/cr/`, `CHANGELOG.md`,
      `spec/lastenheft.md`). Bloße Bestands-Prosa, die einen vergangenen
      Kanon-Stand als Herkunfts-Anker nennt (z. B. „seit `v5.18.0`" in
      `observations.md`), ist davon nicht betroffen — sie trägt kein
      `.harness/baseline/`-Pfadsegment und bleibt unverändert.
- [ ] Adaptions-Review aller lebenden `MR`-Einträge durchgeführt, Ergebnis
      benannt (wie viele bleiben gültig, wie viele angepasst).
- [ ] `.harness/baseline/v5.18.0/` entfernt, `v6.0.0` einziger vendorter Baum.
- [ ] `harness/conventions.md` §Baseline `**Stand:**` auf `v6.0.0`.
- [ ] `make baseline-verify` und `make baseline-freshness` (Content-Teil)
      grün auf `v6.0.0`.
- [ ] `make gates` grün (zehn Gates).
- [ ] `make fullbuild` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`.
- [ ] Unabhängige Verifikation durchgeführt.
- [ ] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Das Ein-Commit-Muster aus §2 Schritt 2 könnte wieder zusammenfallen**,
  wenn der Bump unter Zeitdruck in einem Rutsch committet wird — dieselbe
  Beobachtung wie bei slice-189 (F-3), diesmal bewusst vorab benannt statt
  erst im Review gefunden.
- **Die Zeilenspann-Verschiebung in `modul-05` ist von Hand nachgerechnet**,
  nicht durch ein Werkzeug — ein Rechenfehler bliebe unentdeckt, bis
  `citations` ihn meldet (fail-closed, würde also auffallen, aber ggf. erst
  spät im Lauf).
- **Der Adaptions-Review-Umfang** (Zahl lebender `MR`-Einträge) ist zum
  Planungszeitpunkt nicht neu gezählt — Risiko einer Unterschätzung wie bei
  slice-189 selbst benannt.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt
keinen Slice; `v6.0.0` ist upstream verfügbar (gemessen per direkter
GitHub-Release-Abfrage, `make baseline-freshness`s `check-latest` zeigt es
aus ungeklärtem Grund nicht — separat zu beobachten, kein Blocker für
diesen Slice).

**Rückführungen:** `in-progress` → `open`, falls der Adaptions-Review eine
Kollision wie bei [`MR-057`](../../../../harness/conventions.md#mr-057)/[`MR-058`](../../../../harness/conventions.md#mr-058) aufdeckt, deren Auflösung selbst ein
Slice ist.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `.harness/baseline/` (vendorte Fremd-Konvention) und
  `harness/` (der Konventionsspeicher). Beide fallen unter den Default `*` =
  **Greenfield** ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration).

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:223-224 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Stand 2026-09-03, höchste Kennung
  `BEO-027`): [`BEO-008`](../observations.md) (**5**, verkörpert, aber
  weiter offen — die Pin-Spiegel-Klassen eines Bumps, hier durch den
  expliziten Schritt 4 im Vorgehen adressiert); [`BEO-009`](../observations.md)
  (**9**, verkörpert, weiter offen — Commit-Botschaften nicht über die
  gemessene Menge hinaus verallgemeinern, hier besonders beim
  Adaptions-Review-Ergebnis); [`BEO-013`](../observations.md) (**1**) — die
  Bestands-Stichprobe gehört gefahren, auch wenn der Pin danach aktuell ist;
  [`BEO-027`](../observations.md) (**1**, `docs/plan/planning/observations.md`
  als Sub-Area) — betrifft das Register selbst, nicht diesen Slice direkt,
  aber relevanter Kontext für slice-194/195.

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:229-229 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)) — beide Läufe
  grün (`upstream-drift.yml` 2026-09-03T14:47:17Z, `image-scan.yml`
  2026-09-03T08:05:17Z). **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-193. Betroffene IDs:
[`MR-013`](../../../../harness/conventions.md#mr-013),
[`MR-021`](../../../../harness/conventions.md#mr-021),
[`MR-051`](../../../../harness/conventions.md#mr-051),
[`MR-055`](../../../../harness/conventions.md#mr-055),
[`MR-058`](../../../../harness/conventions.md#mr-058). Module: `links`,
`anchors`, `citations`. Gates: `make baseline-verify`,
`make baseline-freshness`, `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter
den Default: Doc führt, Code folgt. Ein Adoptions-Stand und sein
Konventionsspeicher; kein Produkt-Code, keine Reconciliation. Das
**Evidenz-Risiko** ist die einzige Achse mit Substanz: der vendorte Baum ist
Fremd-Inhalt, und ob eine Adaption noch trägt, entscheidet sein Delta — nicht
unser Bestand.

## 9. Closure-Notiz (nach `done/`)

<!-- wird erst bei Closure gefüllt -->
