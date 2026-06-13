# Slice slice-007: Modul `matrix` + Dogfooding-Selbstkonfiguration

**Status:** done.

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
Selbstkonfiguration — die [`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)-Referenzrichtungs-Regel
(Spec-Straten verweisen nie abwärts) wird damit maschinell erzwungen.

## 2. Definition of Done

- [x] Akzeptanzkriterien von [`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix) als Tests: Slice →
  aktives ADR ok; Referenz auf `Superseded …` → `matrix-inactive`;
  Lastenheft → ADR → `matrix-forbidden` mit beiden Klassen.
- [x] Status-Extraktion in fester Reihenfolge (`**Status:**` vor
  `Status`-Heading), Präfix-Match case-insensitiv; ohne Status aktiv
  (Spezifikation §[`DC-FA-MTX-001.a`](../../../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung)).
- [x] Klassen-Glob-Präzedenz (Deklarationsreihenfolge) und
  `exclude-sections` (getrimmt, ohne Auszeichnung, case-sensitiv)
  getestet.
- [x] Selbstkonfiguration aktiv: `.d-check.yml` deklariert
  Dokumentklassen (Spec-Straten, ADR, Slice), Regeln
  `{spec-straten → adr/slice: verboten}`, `status.forbidden`,
  `exclude-sections` (Historie/Geschichte) sowie `ids`-Muster für
  `DC-*`, `MR-*`, `ADR-*`; `make doc-check` läuft mit
  `links, anchors, ids, matrix` und ist auf dem eigenen Repo grün.
- [x] `matrix` in `isImplemented`; `make gates` grün;
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

- Nackte Kennungen in **immutablen** Texten: Bereinigt wird nur
  lebende Doku — einmal `Accepted`e ADRs und akzeptierte MR-Einträge
  werden **nicht** editiert, auch nicht zum Nachverlinken
  (`AGENTS.md` §3.5). Inline-Code ist linkpflichtfrei und
  Historie-Sektionen sind per `exclude-sections` ausnehmbar; bleiben
  darüber hinaus Befunde in unveränderlichen Passagen, ist eine
  dokumentierte Entscheidung nötig (Spez-Fortschreibung einer
  Sektions-/Pfad-Ausnahme für `ids` ODER Form-Fix analog [`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)) —
  nie stillschweigend lockern.
- Klassen-Zuordnung über Globs muss mit der `"."`-Wurzel des
  Dogfoodings zusammenspielen (Pfade relativ zur Repo-Wurzel).
- Zeilenbasierte Link-Extraktion ist normative Grenze (Spezifikation
  §[`DC-FA-LINK-001.a`](../../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion) Schritt 3, Review R1 zu slice-006): Kennungen in
  *mehrzeiligem* Linktext gelten für `ids` als nackt. Bei der
  Selbstkonfiguration darauf achten, dass linkpflichtige Kennungen in
  der eigenen Doku in einzeiligen Links stehen; auftretende Fälle
  sind Form-Fixes der lebenden Doku, keine Modul-Bugs.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Commit `18a4976` (Modul `matrix`, Config-Durchreichung,
Selbstkonfiguration, ids-Fortschreibung, Doku-Bereinigung).

- **Was hat funktioniert:** Die normierte Status-Extraktion und die
  Sektions-Ausnahmen ließen sich direkt aus der Spezifikation in
  Tabellen-Tests übersetzen; der gemeinsame Heading-Scanner
  (`extractHeadingLines`) trägt jetzt `anchors` und `matrix` ohne
  Dopplung. Die [`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)-Referenzrichtung ist maschinell kodiert —
  der Kern-Vertrag dieses Slices.
- **Anders als geplant:** Der erste Selbstlauf lieferte 111 Befunde
  und legte eine konzeptuelle `ids`-Lücke offen: Definitions-Orte
  (Lastenheft-Headings, MR-Einträge, ADR-Korpus) und Heading-Zeilen
  wurden geflaggt, obwohl eine Definition nicht auf sich selbst
  verlinken kann. Statt der im Risiko-Block erwogenen Datei-Listen-
  Ausnahme wurde die Spezifikation prinzipiell fortgeschrieben
  (Headings und Muster-Target sind kein linkpflichtiger Fließtext) —
  das löste zugleich die Immutable-ADR-Frage ohne jede Editierung
  der `Accepted`-Texte. Restliche ~50 Fließtext-Befunde wurden als
  Form-Fixes bereinigt (lebende Docs: Links; historische
  Docs/CHANGELOG: Code-Spans, analog [`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)).
- **Steering-Loop-Lerneintrag:** (a) Der Dogfooding-Selbstlauf ist
  der schärfste Spec-Reviewer — die Definitions-Ort-Lücke war in
  Spec-Review und slice-006-Implementierung unsichtbar und wurde
  erst durch die Anwendung auf den eigenen Bestand messbar (111 → 0
  Befunde). Lehre für slice-008: das Netzlos-Gate ebenfalls zuerst
  gegen das eigene Repo kalibrieren. (b) `gocognit` schlug beim
  Erweitern von `checkIDs` an — Komplexitäts-Split statt Ausnahme
  bestätigt die slice-005-Politik.
- **Folge-Slices:** keine neuen; slice-008 (`external`) und slice-009
  (Gates) stehen wie geplant; nach slice-008 den
  Interim-Mechanismus `isImplemented`/`SkippedModules` entfernen
  (Hinweis aus slice-006-Closure bleibt gültig).

**Review R1 (nach Closure, Agent-Review mit getrenntem Kontext):**
6 Findings (2 MEDIUM, 3 LOW, 1 INFO), alle nachgeschärft in einem
Folge-Commit: Status-Extraktion fence-aware (Form 2 las Fence-Inhalt
als Statuswert — der dritte unabhängige Fence-Automat war die
Ursache; jetzt gemeinsamer Scanner `proseLines` für
PreprocessMarkdown/Headings/Status); `.md`-Guard vor dem Status-Read
(kein Voll-Read von Binärzielen); gemeinsame Ziel-Auflösung
`localTarget` für `links`+`matrix` (dritte Kopie beseitigt);
Doppel-Befund forbidden+inactive als unabhängige Verletzungen in der
Spezifikation expliziert (mit Test); `exclude-sections` der
Selbstkonfiguration um „7. Historie" ergänzt — der Eintrag „Historie"
war gegen die real nummerierten Spec-Headings wirkungslos (tote
Config als kleine Harness-Lüge, vom Review gefunden, nicht vom
grünen Dogfooding-Lauf: dort existieren noch keine Links in den
Historie-Sektionen). Steering-Loop-Beleg: ein grüner Selbstlauf
beweist nur die geprüften Fälle — tote Konfiguration findet erst der
Review mit Negativ-Blick.

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (spec-first; siehe Kurs Modul 5 §Worked
Mini-Example).
