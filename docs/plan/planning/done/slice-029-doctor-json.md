# Slice slice-029: Maschinenlesbare Diagnose (`--doctor --json`)

**Status:** done (abgeschlossen 2026-06-19).

**Welle:** welle-18-doctor-json (Trigger: Change Request 0.17.0 akzeptiert;
baut auf slice-025 (Fix-Kandidaten-Modell) und slice-026).

**Bezug:**
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
(Hauptanforderung — JSON-Rendering der Diagnose),
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
(Kombinierbarkeit + JSON-Schema),
[`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)
(Exit-Codes),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(byte-identisch),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(read-only).

**Autor:** pt9912. **Datum:** 2026-06-18.

---

## 1. Ziel

`d-check --doctor --json` gibt die Diagnose **maschinenlesbar** aus: ein
JSON-Dokument, dessen `findings` je Eintrag zusätzlich `reasonText`
(Grund-Klartext) und `fixCandidate` (`{original, replacement, note}` oder
`null`) tragen; die Gruppierung nach Datei trägt das `file`-Feld. Es ist
**dasselbe Modell wie die Prosa-Diagnose**, ein drittes Rendering neben
Prosa (slice-025) und Patch (slice-026) — über `core.ReasonText` und
`core.FixCandidateFor`.

## 2. Definition of Done

- [x] **CLI:** `--doctor --json` ist erlaubt; der Kombinations-Check
  weist nur noch `--repair`+`--json` und `--doctor`+`--repair` als
  Nutzungsfehler (Exit 2) ab.
- [x] **Renderer:** JSON-Diagnose im report-Adapter — `findings` mit
  `reasonText` und `fixCandidate` (explizit `null`, nicht weggelassen),
  dazu `summary`/`exitCode`; stdout enthält nur das JSON-Dokument.
- [x] **Spezifikation:**
  [`DC-FA-CLI-007.a`](../../../../spec/spezifikation.md#dc-fa-cli-007a--diagnose-modus)
  + JSON-Ausgabe-Schema (§2) um `reasonText`/`fixCandidate` ergänzt.
- [x] **Tests:** JSON-Variante (Felder vorhanden, stdout parsbar, Exit 1);
  `--doctor --repair` → Exit 2; Determinismus 10× ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus));
  read-only ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- [x] **Doku-Drift schließen:** `docs/user/benutzerhandbuch.md` (§4.9
  Hinweise + FAQ §8) sagt derzeit „`--doctor` nicht mit `--json`
  kombinierbar" — auf den neuen, definierten Modus umstellen und eine
  Gegenüberstellung `--json` / `--doctor` / `--doctor --json` ergänzen.
  `operations.md` ebenfalls.
- [x] `make gates` grün; Closure-Notiz.

## 3. Plan (vor Code)

| Datei | Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | JSON-Diagnose-Schema (`reasonText`, `fixCandidate`) |
| `internal/…` (CLI-Kombi-Check + report-JSON-Renderer) | update | `--doctor --json` erlauben; JSON-Rendering der Diagnose (Wiederverwendung `ReasonText`/`FixCandidateFor`) |
| `docs/user/benutzerhandbuch.md` | update | `--doctor --json` dokumentieren, „nicht kombinierbar"-Stellen korrigieren |
| `docs/user/operations.md` | update | Options-Tabelle nachziehen |

Lastenheft ist mit Change Request 0.17.0 bereits gesetzt.

## 4. Trigger

Change Request 0.17.0 akzeptiert. Voraussetzung erfüllt: slice-025 liefert
`FixCandidateFor`/`ReasonText`.

## 5. Closure-Trigger

DoD vollständig inkl. geschlossener Doku-Drift und grüner Gates.

## 6. Risiken und offene Punkte

- **Schema-Stabilität** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  feste Struct-Reihenfolge, kein Map-Iter; `fixCandidate: null` explizit
  (nicht `omitempty`) — sonst verschwindet die Aussage „kein eindeutiger
  Fix".
- **Doku-Drift ist die eigentliche Falle:** die `--doctor`+`--json`-
  Verbots-Aussagen in Handbuch und `operations.md` müssen mit umgestellt
  werden, sonst widerspricht die Doku dem neuen Vertrag.
- **`--repair --json` bleibt Fehler** — der Kombinations-Check darf nicht
  versehentlich auch `--repair`+`--json` freigeben.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Vertrag
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
(Lastenheft 0.17.0, JSON-Variante) + Spezifikation
[`DC-FA-CLI-007.a`](../../../../spec/spezifikation.md#dc-fa-cli-007a--diagnose-modus)
Schritt 6 und §2 „JSON-Diagnose". `report.DoctorJSON` ist das dritte
Rendering desselben Fix-Kandidaten-Modells (`core.ReasonText` +
`core.FixCandidateFor`) — eingebettetes `core.Finding` plus `reasonText`
und `fixCandidate`; der Pointer trägt **kein** `omitempty`, fehlt der
Kandidat steht explizit `null`. CLI-Kombi-Check umgestellt: nur noch
`--repair`+`--json` und `--doctor`+`--repair` sind Nutzungsfehler (Exit 2),
neuer render-Switch-Zweig `--doctor --json`. Doku-Drift in Handbuch §4.9 +
`operations.md` geschlossen, Drei-Wege-Gegenüberstellung
(`--json`/`--doctor`/`--doctor --json`) ergänzt. `make gates` grün
(Coverage 93,80 %).

**Belege:**

- **JSON-Variante** (Happy): `findings` mit `reasonText` und
  `fixCandidate`, `"fixCandidate": null` explizit, stdout reines JSON,
  Exit 1 (`TestCLI007_DoctorJSON_Happy`).
- `--doctor --repair` → Exit 2 (`TestCLI008_Repair_Inkompatibel`);
  `--repair --json` → Exit 2 unverändert.
- Determinismus 10× byte-identisch
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
  `TestCLI007_DoctorJSON_Determinismus`); read-only
  ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit))
  strukturell + global per `make arch-check` / netzlosem `make doc-check`.

**Lerneintrag:** „Eine Quelle, drei Ausgaben" ist eingelöst — dieselbe
Ableitung speist Prosa (slice-025), Patch (slice-026) und JSON
(slice-029). Das `fixCandidate`-Pointer-Feld **ohne** `omitempty` ist der
Trick für explizites `null`: die Aussage „kein eindeutiger Fix" überlebt
die Serialisierung. Nebenbefund aus dem Review: ein Review-Report, der
`pfad.go:ZEILE-BEREICH` als Inline-Code schreibt, trippt das eigene
`codepaths`-Gate (der Range-Suffix ist kein existierender Pfad) — Vollpfad
nur mit Einzelzeile oder den Bereich aus dem Code-Span heraushalten.

**Review R1** (unabhängiger Reviewer-Subagent, eigener Kontext ohne
DoD-Wissen,
[Report](../../../reviews/2026-06-19-slice-029-doctor-json.md)): HIGH 0 /
MEDIUM 0 / LOW 0 / INFO 2 — freigegeben. INFO-1 (promotetes
`message`-Feld im §2-Schema undokumentiert) im selben Stand geschlossen
(Spec-Klarstellung). INFO-2 (kein dedizierter read-only-Test des
`--doctor --json`-Pfads) akzeptiert: read-only ist strukturell geerbt
(gleicher `render`-Pfad, Schreiben nur in den injizierten stdout) und wird
wie bei slice-025 global gemessen, kein per-Modus-Test.

**Welle:** welle-18-doctor-json ist damit vollständig (slice-029).

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Spec-/Code-/Doku-Arbeit; Greenfield-Default).
