# Slice slice-029: Maschinenlesbare Diagnose (`--doctor --json`)

**Status:** open (geplant).

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

- [ ] **CLI:** `--doctor --json` ist erlaubt; der Kombinations-Check
  weist nur noch `--repair`+`--json` und `--doctor`+`--repair` als
  Nutzungsfehler (Exit 2) ab.
- [ ] **Renderer:** JSON-Diagnose im report-Adapter — `findings` mit
  `reasonText` und `fixCandidate` (explizit `null`, nicht weggelassen),
  dazu `summary`/`exitCode`; stdout enthält nur das JSON-Dokument.
- [ ] **Spezifikation:**
  [`DC-FA-CLI-007.a`](../../../../spec/spezifikation.md#dc-fa-cli-007a--diagnose-modus)
  + JSON-Ausgabe-Schema (§2) um `reasonText`/`fixCandidate` ergänzt.
- [ ] **Tests:** JSON-Variante (Felder vorhanden, stdout parsbar, Exit 1);
  `--doctor --repair` → Exit 2; Determinismus 10× ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus));
  read-only ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- [ ] **Doku-Drift schließen:** `docs/user/benutzerhandbuch.md` (§4.9
  Hinweise + FAQ §8) sagt derzeit „`--doctor` nicht mit `--json`
  kombinierbar" — auf den neuen, definierten Modus umstellen und eine
  Gegenüberstellung `--json` / `--doctor` / `--doctor --json` ergänzen.
  `operations.md` ebenfalls.
- [ ] `make gates` grün; Closure-Notiz.

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

*(wird bei Closure gefüllt — Umsetzung, Belege, Lerneintrag, Review-Runde.)*

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Spec-/Code-/Doku-Arbeit; Greenfield-Default).
