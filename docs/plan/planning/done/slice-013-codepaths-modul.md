# Slice slice-013: Modul `codepaths` — Pfade in Inline-Code

**Status:** done.

**Welle:** welle-04-distribution-und-migration (Einschub vor
slice-012 — Nummerierung ist chronologisch, Reihenfolge per Trigger).

**Bezug:** [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(Change Request 0.3.0),
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
(Modul-Auswahl),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(Konfigurations-Vollvalidierung),
[ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md) (Layout).

**Autor:** pt9912. **Datum:** 2026-06-11.

---

## 1. Ziel

Das Modul `codepaths` ist normiert spezifiziert und implementiert —
die letzte generische Konsolidierungs-Lücke gegenüber der JS-Familie
(`docs-check.js`) ist geschlossen; slice-012 kann den Kurs danach
vollständig migrieren.

## 2. Definition of Done

- [x] Spezifikation fortgeschrieben (Doc führt, Code folgt):
  §[`DC-FA-CODE-001.a`](../../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) (Erkennungs-Algorithmus: Inline-Code-Spans,
  `./`/`../` immer + konfigurierte Präfixe, konservative Ausschlüsse,
  Marker-Semantik), `.d-check.yml`-Schema um `codepaths.roots`
  (string[], Constraint-Tabelle), §4 Grund-Codes (z. B.
  `codepath-missing`, Wiederverwendung von `repo-escape`), Historie.
- [x] Akzeptanzkriterien von [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) als Tests
  (Happy/Boundary/Negative); Boundary belegt zusätzlich, dass der
  Marker **nur** dieses Modul stilllegt (Link-Befund derselben Zeile
  bleibt bestehen).
- [x] Kern-Modul `internal/hexagon/core/codepaths.go` (nutzt die
  bestehende Inline-Code-Erkennung aus `markdown.go` — Spans werden
  hierfür *gelesen* statt gestrippt); Config-Durchreichung
  (`configyaml`), Modul in `validModules`.
- [x] Dogfooding: Selbstkonfiguration (`.d-check.yml`) aktiviert
  `codepaths` mit passenden Präfixen; eigene Doku ist befundfrei
  (Form-Fixes, wo nötig). Die QA-03-Modullisten-Prüfung in
  `tools/gate-consistency.sh` zieht das neue Modul in ihre
  „alle außer `external`"-Zusage ein.
- [x] Alle bestehenden Modul-Aufzählungen nachgezogen (Review-R1-
  Finding F-4): `spec/spezifikation.md` (JSON-Schema-Enum,
  `rule`-Tabellenzeile, `modules`-Constraint),
  [`docs/user/operations.md`](../../../../docs/user/operations.md)
  (Optionen-Tabelle), [`README.md`](../../../../README.md)
  (Status-Absatz), doc-check-Zeilen in
  [`harness/README.md`](../../../../harness/README.md) §Sensors und
  [`AGENTS.md`](../../../../AGENTS.md) §4 (Modul-Aufzählung +
  DC-Bindung um [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)).
- [x] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md);
  Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | update | §[`DC-FA-CODE-001.a`](../../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code), Schema, Grund-Codes (vor dem Code) |
| `internal/hexagon/core/codepaths.go` (+ Test) | neu | Kern-Modul gegen In-Memory-FS |
| [`internal/hexagon/core/markdown.go`](../../../../internal/hexagon/core/markdown.go) | update | Inline-Code-Spans als Werte zugänglich machen (heute nur positionserhaltendes Stripping) |
| [`internal/adapter/driven/configyaml/configyaml.go`](../../../../internal/adapter/driven/configyaml/configyaml.go) | update | `codepaths.roots` strikt validieren |
| [`.d-check.yml`](../../../../.d-check.yml), [`tools/gate-consistency.sh`](../../../../tools/gate-consistency.sh) | update | Dogfooding + QA-03-Modulliste |

## 4. Trigger

Sofort — der Change Request (Lastenheft 0.3.0) ist vom Auftraggeber
freigegeben (2026-06-11).

## 5. Closure-Trigger

DoD vollständig + Commit(s) auf `main` + Closure-Notiz; entriegelt
gemeinsam mit slice-011 den slice-012.

## 6. Risiken und offene Punkte

- Falsch-Positiv-Gefahr der Pfad-Heuristik (Prosa-Backticks, die wie
  Pfade aussehen): konservative Erkennung übernehmen
  (Whitespace/Glob-Zeichen/Ellipsen ausschließen — erprobtes
  docs-check-Verhalten) und am eigenen Repo plus Kurs-Repo
  gegentesten, **bevor** der Slice schließt.
- Der Marker ist ein neuer CLI-Vertrag (`d-check:ignore`) —
  Begründungs-Klammer empfohlen, nicht erzwungen (Erzwingung wäre
  Stil-Polizei; das Review-Argument lebt in der Doku).
- `markdown.go`-Umbau berührt gemeinsame Infrastruktur aller Module —
  Regressionsschutz über die bestehende Testbasis (95 % Coverage).

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Commits `a6fcec4` (Spec + Modul + Dogfooding) und
`169f971` (Linktext-Ausnahme aus dem Gegentest). **Gegentest am
Kurs-Repo (Risiko-Auflage):** Clean-Lauf 116 Dateien, **0 = 0**
Befunde gegenüber `docs-check.js` (nach Marker-Umbenennung);
Negativ-Probe fängt alle vier injizierten Fehlerklassen — inklusive
der Inline-Pfad-Klasse, deren Fehlen den Change Request auslöste.

- **Was hat funktioniert:** Spec-first hat sich doppelt ausgezahlt —
  und das Dogfooding war erneut der schärfste Reviewer: Der Erstlauf
  am eigenen Repo (28 Befunde) erzwang zwei prinzipielle
  Fortschreibungen statt Marker-Spam: die `Datei:Zeile`-Konvention
  gehört in die Normalisierung (acht Review-Report-Felder), und
  **Headings sind ausgenommen** — entdeckt, weil ein Marker im
  [`MR-003`](../../../../harness/conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh)-Heading
  dessen Anker-Slug brach und sechs Links rot wurden
  (ein Drei-Module-Zusammenspiel, das nur der Selbstlauf zeigt).
- **Anders als geplant:** (a) Die Linktext-Span-Ausnahme
  (`` [`pfad`](ziel) ``) fehlte in Anforderung und Spez — der
  Kurs-Gegentest fand sie als einzige Paritäts-Differenz; docs-check
  hatte die Regel, unser Faktencheck hatte sie übersehen. (b) 17
  begründete `d-check:ignore`-Marker im eigenen Repo (historische
  Pfade in ADRs/MRs/Closures, AK-Beispiele im Lastenheft, das eigene
  CHANGELOG-Syntax-Beispiel) — der Marker-mit-Begründung-Mechanismus
  trägt auch als Form-Annotation an immutablen Artefakten.
- **Steering-Loop-Lerneintrag:** Erkennungs-Heuristiken über fremden
  Bestand kalibriert man nicht am Schreibtisch — erst der eigene
  Selbstlauf, dann der Gegentest am Ziel-Repo fanden je eine
  Regel-Lücke, die weder Review noch Faktencheck des Quell-Tools
  sahen. Für slice-012 heißt das: Vergleichslauf **vor** der
  CI-Umstellung ist Pflicht, nicht Kür.
- **Folge-Slices:** keine neuen; slice-012 ist entriegelt
  (011 ✓ + 013 ✓). Hinweis dorthin: `codepaths` ist erst nach dem
  nächsten Release im GHCR-Image — der Kurs-Vergleichslauf dort
  braucht v0.2.x oder einen lokalen Build.

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (spec-first; siehe Kurs Modul 5 §Worked
Mini-Example).
