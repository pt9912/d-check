# Slice slice-021: Reichhaltige `--help` (Synopsis, `[pfad]`, Config-Pointer)

**Status:** in-progress.

**Welle:** welle-11-help (Trigger: Priorisierung durch den Auftraggeber
— die Default-`flag`-Usage verschweigt das `[pfad]`-Argument).

**Bezug:** [`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
(Schärfung — die Hilfe nennt Synopsis + Pfad-Argument + Config-Pointer),
[`DC-FA-CLI-005`](../../../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
/ [`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(die Hilfe verweist darauf für das Config-Format).

**Autor:** pt9912. **Datum:** 2026-06-13.

---

## 1. Ziel

`d-check --help` zeigt heute nur die nackte `flag`-Paket-Default-Liste —
ohne Synopsis, ohne das Positions-Argument `[pfad]` (die Scan-Wurzel),
ohne Kurzbeschreibung. Ein Erstnutzer erfährt nicht, dass er einen Pfad
übergeben kann. Eine eigene `Usage`-Funktion behebt das. Das
Config-**Format** wird **nicht** dupliziert — die Hilfe verweist auf
`--print-config` (Gerüst) und `--suggest-config` (Vorschlag); dort lebt
das Format kommentiert (kein Drift).

## 2. Definition of Done

- [ ] **Lastenheft-Schärfung** [`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
  (Version 0.11.0): neues Akzeptanzkriterium — `--help` nennt die
  Synopsis `d-check [optionen] [pfad]`, beschreibt das Pfad-Argument
  (Scan-Wurzel, Default cwd) und verweist für die Konfiguration auf
  `--print-config`.
- [ ] **Spezifikation** §[`DC-FA-CLI-001.a`](../../../../spec/spezifikation.md#dc-fa-cli-001a--ablauf-eines-prüflaufs): Inhalt der `Usage`-Ausgabe
  präzisiert (Synopsis, Pfad-Zeile, Config-Pointer, Flag-Liste auf stderr).
- [ ] **Implementierung**: eigene `flags.Usage`-Funktion; bestehender
  Exit-0-/stderr-Pfad unverändert; Test prüft Synopsis, `[pfad]` und den
  `--print-config`-Verweis.
- [ ] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md);
  Closure-Notiz.

## 3. Plan (vor Code)

| Datei | Art | Begründung |
|---|---|---|
| [`spec/lastenheft.md`](../../../../spec/lastenheft.md) | update | [`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)-Schärfung (0.11.0) |
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | update | §[`DC-FA-CLI-001.a`](../../../../spec/spezifikation.md#dc-fa-cli-001a--ablauf-eines-prüflaufs) Usage-Inhalt |
| `internal/adapter/driving/cli/cli.go` | update | eigene `Usage`-Funktion |
| `internal/adapter/driving/cli/cli_acceptance_test.go` | update | Hilfe-Inhalt geprüft |

## 4. Trigger

Priorisierung durch den Auftraggeber: die Hilfe verschweigt das
`[pfad]`-Argument — eine Erstkontakt-Lücke.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün.

## 6. Risiken und offene Punkte

- **Drift Hilfe ↔ Config-Format:** bewusst vermieden — die Hilfe
  dupliziert das Format nicht, sondern verweist auf `--print-config`.

## 7. Closure-Notiz (nach `done/`)

*(folgt mit dem Lifecycle-Übergang nach `done/`.)*

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Spec-/Code-/Doku-Arbeit; Greenfield-Default).
