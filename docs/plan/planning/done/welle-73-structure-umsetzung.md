# Welle welle-73-structure-umsetzung: Das 20. Regelmodul und die verkörperte Spiegel-Regel

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-73-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Modul-Umsetzung nach vorliegendem
Vertrag).

**Verantwortlich:** pt9912. **Datum:** 2026-08-15.

---

## 1. Welle-Ziel

Das Modul `structure` liefern — Anforderung, Algorithmus und ADR liegen seit
welle-69 vollständig vor, dieser Zug setzt sie um.

**Das Mehr gegenüber der Slice-DoD ist eine Verkörperung.** Der Slice ändert die
Grund-Code-Menge um **neun** Einträge und berührt damit mehr Spiegel als jeder
Slice davor. Genau diese Klasse führt das Beobachtungs-Register als **BEO-002**,
sie ist **dreimal** eingetreten und kein einziges Mal von einem Gate gefunden
worden — jedes Mal von einem Review. Die Schwelle war erreicht; sie hier ein
viertes Mal zu zählen wäre die falsche Antwort gewesen.

**Verkörpert ist sie mit [`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)**
— vor dem Editieren einer zugesagten Semantik wird die Liste ihrer Spiegel
aufgeschrieben. Die Regel steht **vor** dem Slice, nicht danach: sie soll an ihm
zum ersten Mal wirken, nicht an ihm zum vierten Mal fehlen.

## 2. Trigger (Welle startet)

Freigabe des Auftraggebers (2026-08-15), WIP-Slot frei, und die **bindende**
Start-Bedingung des Slice ist erfüllt: [slice-096](slice-096-structure-modul-analyse.md)
**und** [slice-101](slice-101-fence-unbalanciert.md) liegen in `done/`. Die
zweite war nicht Vorliebe, sondern Bedingung — sonst erbte das neue Modul über
die geteilte Mechanik einen **bekannten** stillen Grün-Pfad.

## 3. Closure-Trigger (Welle schließt)

- [slice-099](slice-099-structure-modul.md) liegt in `done/`.
- **[`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten) hat an diesem Slice gewirkt:** die Spiegel-Liste steht im Slice,
  **bevor** die Grund-Codes geändert wurden, und ist am Ende abgehakt.
- Release als **Minor**, Digest-Backfill.
- `make fullbuild` grün.

## 4. Slices in dieser Welle

| Slice | Rolle |
|---|---|
| [slice-099](slice-099-structure-modul.md) | Modul `structure` vollständig + Preset-Kopplung der Closure-Fähigkeit samt `closure-note-ambiguous` |

Ein Slice, und das ist Absicht: der Umsetzbarkeits-Review hat den früheren
Zwei-Slice-Schnitt an der **Release-Grenze** verworfen — ein Modul, das die
Hälfte seines eigenen veröffentlichten Schemas mit Exit 2 ablehnt, ist kein
lieferbarer Zwischenstand.

## 5. Abhängigkeiten

- [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  und [ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md) liegen
  vollständig vor — dieser Zug erfindet nichts, er liefert.
- **Paritäts-Fixtures** aus dem Schwester-Repo: beizuziehen, nicht nachzubauen.
  Beim Wellen-Start geprüft, das Repo liegt vor.

## 6. Out-of-Scope für diese Welle

- **Die Marken-Syntax konfigurierbar machen.** `**M:**` ist die Form der beiden
  vermessenen Repos; ob ein drittes anders schreibt, ist eine Frage für den
  Re-Evaluierungs-Trigger, nicht für diese Umsetzung.
- **Ein Gate für die Spiegel-Konsistenz.** Die Regel ist eine für Menschen;
  sie durch ein Modul zu ersetzen wäre ein eigener Slice — und der
  Auflösungs-Trigger der Regel sagt genau das.
- [slice-103](slice-103-geteilte-lexik-raender.md) und
  [slice-095](slice-095-links-resolve-from.md).

## 7. Closure-Notiz

Geschlossen am 2026-08-15 mit **v0.57.0**. Alle vier Closure-Trigger sind
erfüllt: [slice-099](slice-099-structure-modul.md) liegt in `done/`, die
Spiegel-Regel hat an ihm gewirkt und ist mit Bilanz abgehakt, das Release samt
Digest-Backfill ist draußen, und `make fullbuild` ist grün.

Die vollständige Notiz — geliefert, gelernt, Register-Lese-Schritt und
Trigger-Audit — steht in [`welle-73-results.md`](welle-73-results.md).
