# Slice slice-023: `ids`-Ventile gelten auch für nackte Vorkommen

**Status:** in-progress.

**Welle:** welle-13-ventil-prosa (Trigger: reproduzierter Fremd-Repo-Befund,
Priorisierung durch den Auftraggeber im Dialog 2026-06-16).

**Bezug:** [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(Change Request 0.13.0 — Geltungsbereich der beiden Ventile auf nackte
Fließtext-Vorkommen erweitert),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Abwärtskompatibilität: Configs ohne gesetzte Ventile byte-identisch),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(netzloser Dogfooding-Lauf bleibt grün).

**Autor:** pt9912. **Datum:** 2026-06-16.

---

## 1. Ziel

Die beiden mit slice-018 eingeführten Ventile des Moduls `ids` —
`exempt-paths` (Datei-Glob) und der Zeilen-Marker `d-check:ignore` —
greifen heute **nur** auf Kennungs-Vorkommen *innerhalb von
Inline-Code-Spans* (der von `link-policy: always` zusätzlich geprüfte
Befundsatz). Ein **nacktes** Vorkommen im Fließtext läuft durch den
Prosa-Basispfad, der keines der beiden Ventile konsultiert. Folge: eine
Kennung als Klartext in einer `exempt-paths`-Datei (z. B. einem
Review-Report) wird weiterhin als `id-unlinked` gemeldet, und ein
`d-check:ignore` auf einer Prosa-Zeile bleibt für `ids` wirkungslos.

Reproduktion (Auftraggeber, verbatim): Fremd-Repo `ai-harness-init`,
`d-check` v0.8.0 und v0.9.0 — eine nackte Kennung in einem Review-Report
wurde trotz greifendem Review-Verzeichnis-Exempt als `id-unlinked`
gemeldet (Exit 1). Lokaler Gegenbeweis an der Quelle: dieselbe exempte
Datei mit beiden Formen ergab genau **einen** Befund (die nackte Form),
die Inline-Code-Form blieb frei.

Ziel: Beide Ventile werden ein **Ganzdatei- bzw. Ganzzeilen-Carve-out**
für *alle* Vorkommen des Musters — nackt wie in Inline-Code — und
unabhängig von der `link-policy`. „Datei ausnehmen" heißt dann die ganze
Datei, nicht nur ihre Backtick-Vorkommen.

## 2. Definition of Done

- [x] **Lastenheft-Change-Request**
  [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
  (Version-Bump 0.13.0, Historie-Zeile): Ventil-Beschreibung neugefasst
  (Ganzdatei-/Ganzzeilen-Carve-out, alle Vorkommen, politik-unabhängig);
  Negative-AK präzisiert (außerhalb `exempt-paths`, ohne
  `d-check:ignore`); neues AK „Ventile für nackte Vorkommen".
- [x] **Spezifikation**
  §[`DC-FA-ID-001.a`](../../../../spec/spezifikation.md#dc-fa-id-001a--kennungs-prüfung):
  der Prosa-Basisalgorithmus nennt beide Ventile als linkpflichtfrei;
  der `always`-Absatz stellt klar, dass (3)/(4) dieselben Ventile wie im
  Prosa-Pfad sind; Schema-Tabelle und Beispiel-Kommentar entkoppelt von
  „nur `always`". **Kein** neuer Grund-Code (weiter `id-unlinked`).
- [x] **Implementierung** im Modul `ids`: `checkIDLine` ehrt
  `ignored(file, p.ExemptPaths)`; `checkIDs` überspringt Prosa-Zeilen mit
  `d-check:ignore` (neuer Helper `ignoreMarkerLines`, roh-zeilenbasiert
  wie der Inline-Code-Pfad und `codepaths`). Die zwei AKs als Tests.
- [x] **Abwärtskompatibilitäts-Beleg**
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  Configs ohne gesetzte Ventile sind unberührt (`ignored(file, nil)`
  ⇒ false; leere Marker-Menge); Wirkung nur in Richtung *weniger*
  Befunde in explizit ausgenommenen Dateien/Zeilen.
- [x] **Doku-Nachzug** (`config.go`-Kommentar, `docs/user/operations.md`)
  von der „nur `always`"-Kopplung gelöst.
- [ ] `make gates` grün;
  [`CHANGELOG.md`](../../../../CHANGELOG.md); Closure-Notiz mit
  Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`spec/lastenheft.md`](../../../../spec/lastenheft.md) | update | CR [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids) 0.13.0: Ventil-Neufassung, AKs, Historie |
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | update | §[`DC-FA-ID-001.a`](../../../../spec/spezifikation.md#dc-fa-id-001a--kennungs-prüfung): Prosa-Ventile, Schema-Tabelle |
| `internal/hexagon/core/ids.go` | update | `checkIDLine` + `exempt-paths`; `checkIDs` + `ignoreMarkerLines` |
| `internal/hexagon/core/ids_test.go` | update | AK „Ventile für nackte Vorkommen" + politik-unabhängig |
| `internal/hexagon/core/config.go` | update | Kommentar `ExemptPaths` entkoppelt |
| [`docs/user/operations.md`](../../../../docs/user/operations.md) | update | Ventil-Beschreibung entkoppelt |
| [`CHANGELOG.md`](../../../../CHANGELOG.md) | update | 0.x-Eintrag (nutzersichtbar) |

## 4. Trigger

Reproduzierter Befund aus einem Fremd-Repo plus Priorisierung durch den
Auftraggeber (Dialog 2026-06-16). Die Diagnose grenzte stufenweise ein:
nicht Glob-Syntax (jede Schreibweise scheiterte), nicht Pfadform (der
emittierte Pfad ist repo-relativ ohne `./`), nicht Versions-Regression
(v0.8.0 ≡ v0.9.0 im betroffenen Pfad) — sondern die strukturelle
Asymmetrie zweier Prüfpfade.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün (echte Ausgabe), Dogfooding zeigt die
exempten Review-Reports weiter bei 0 Befunden, Closure-Notiz mit
Lerneintrag.

## 6. Risiken und offene Punkte

- **Semantik-Verbreiterung:** `exempt-paths` wirkt nun auch unter
  Default-Politik `prose` (vorher dort ein No-op). Das ändert Verhalten
  nur für Configs, die `exempt-paths` gesetzt haben, und ausschließlich
  in Richtung weniger Befunde in den ausgenommenen Dateien. Im Review
  bestätigen.
- **Marker-Erkennung auf der Roh-Zeile:** `ignoreMarkerLines` prüft die
  rohe Prosa-Zeile (`proseLines`), konsistent mit `alwaysLineFindings`
  und `codepaths`; eine in Backticks gesetzte Marker-Zeichenkette ist
  pathologisch und nicht die dokumentierte Nutzung.
- **Determinismus**
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  `ignored`/`matchGlob` ist bereits reihenfolge-/plattformstabil
  (unverändert wiederverwendet).

## 7. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas Greenfield (Spec-/Code-/Doku-Arbeit;
Greenfield-Default der Modus-Tabelle).

## 8. Closure-Notiz (nach `done/`)

*(folgt nach `make gates` grün und Move nach `done/`.)*
