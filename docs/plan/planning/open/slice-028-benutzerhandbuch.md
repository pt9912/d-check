# Slice slice-028: Benutzerhandbuch (`docs/user`)

**Status:** open (geplant).

**Welle:** welle-17-benutzerhandbuch (Trigger: Auftraggeber — ein
nutzbares Benutzerhandbuch nach dem mit-getrackten Standard; baut auf dem
v0.12.0-Funktionsumfang).

**Bezug:**
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(Docker-Aufruf) und
[`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
(CLI-Aufruf) als Einstiegs-Verträge. Das Handbuch ist **abgeleitete
Nutzer-Doku** über die nutzersichtbaren CLI-, Modul- und Distributions-
Verträge (kein neuer Vertrag); `docs/user` ist das Rang-6-Doku-Stratum.

**Autor:** pt9912. **Datum:** 2026-06-18.

---

## 1. Ziel

Ein **aufgabenbasiertes** Benutzerhandbuch unter
docs/user/benutzerhandbuch.md, das **alle Use Cases** von d-check abdeckt
— geschrieben nach dem mit-getrackten
`docs/user/benutzerhandbuch-standard.md` (zielgruppenorientiert,
aufgabenbasiert statt funktionsbasiert, eindeutige Sprache,
Versionsbezug, wartbar). CLI-Adaption des GUI-orientierten Standards:
**Befehle statt Klicks**, **Terminal-Ausgaben statt Screenshots**, keine
Rollen/Rechte (zustandsloses Lese-Tool).

## 2. Definition of Done

- [ ] docs/user/benutzerhandbuch.md nach Standard-Struktur: Einleitung ·
  Erste Schritte · Aufgaben (Schritt-für-Schritt) · Konfiguration ·
  Modul-Referenz · Fehlerbehebung · FAQ · Glossar · Support ·
  Änderungshistorie. Versionsbezug **v0.12.0**.
- [ ] **Alle Use Cases:** Doku prüfen (Default), Docker- und CI-Einbindung,
  `--print-config`/`--suggest-config`, `--enable`/`--disable`, alle acht
  Regelmodule, `--doctor`, `--repair` (konservativ/breit), `--json`,
  Exit-Codes, Fehlerfälle.
- [ ] `benutzerhandbuch-standard.md` als adoptierte Methodik
  **mit-getrackt** (Auftraggeber-Vorlage); das Handbuch verweist darauf.
- [ ] **d-check-doc-check-rein** (das Handbuch wird selbst gescannt):
  keine `DC-*`-IDs im Fließtext, keine Host-Pfade, valide Links/Anker,
  saubere Spans; Beispiel-Befehle/-Ausgaben in Fenced-Blöcken
  (codepaths-frei).
- [ ] `make gates` grün; Closure-Notiz.

## 3. Plan (vor Code)

| Datei | Art | Begründung |
|---|---|---|
| docs/user/benutzerhandbuch.md | neu | das Handbuch (aufgabenbasiert, alle Use Cases) |
| `docs/user/benutzerhandbuch-standard.md` | track | adoptierte Methodik (Auftraggeber-Vorlage), bislang untracked |

Lastenheft/Spezifikation unverändert (abgeleitete Nutzer-Doku, kein
Vertrag).

## 4. Trigger

Auftraggeber. Voraussetzung erfüllt: v0.12.0 released — der nutzersichtbare
Funktionsumfang steht.

## 5. Closure-Trigger

DoD vollständig inkl. `make gates` grün.

## 6. Risiken und offene Punkte

- **Doc-check-Selbstprüfung:** Beispiel-Pfade und -Ausgaben gehören in
  **Fenced-Blöcke** (sonst `codepaths`/`links`-Befunde auf fiktive
  Pfade); interne Anker (Inhaltsverzeichnis) müssen auflösen.
- **Wartbarkeit** (Standard §12): Versionsbezug an v0.12.0 koppeln, bei
  Feature-Änderungen nachziehen — Review-Prozess mit Releases verbinden.
- **GUI-Standard, CLI-Produkt:** „Klicks/Screenshots/Rollen" werden
  bewusst auf Befehle/Terminal-Ausgaben/„kein Auth" übertragen.

## 7. Closure-Notiz (nach `done/`)

*(wird bei Closure gefüllt — Umsetzung, Belege, Lerneintrag, Review-Runde.)*

## 8. Sub-Area-Modus-Begründung

Sub-Area `docs/user` (GF, Greenfield-Default): neue Nutzer-Doku, kein
Bestandscode, Doc führt.
