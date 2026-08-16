# Welle welle-75-wellen-register: Die Roadmap-Aussage gegen das Verzeichnis

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-75-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (dritte Fähigkeit eines
bestehenden Moduls).

**Verantwortlich:** pt9912. **Datum:** 2026-08-16.

---

## 1. Welle-Ziel

Die Lifecycle-Invariante, die das Modul `planning` heute auf der **Slice**-Ebene
prüft, eine Ebene höher ziehen: die **Wellen**-Abschnitte der Roadmap gegen die
Wellen-Dateien. [slice-102](slice-102-wellen-lifecycle-invariante.md)
trägt die vier Aussagen.

**Das Mehr gegenüber der Slice-DoD ist ein Beleg, den es vorher nicht gab.**
Als der Slice geschnitten wurde, war seine Motivation eine Beobachtung: zwei
Wellen-Closures hatten den Ruhe-Marker gesetzt, während die flache Wellen-Datei
noch danebenlag. Seither ist die **Gegenrichtung** eingetreten und gemessen:
in der Closure von welle-73 fehlten **drei** Zeilen in §Abgeschlossene Wellen,
obwohl alle drei Ergebnisnotizen im Ruheort lagen. Das ist die Richtung
„Artefakt ⇒ Register“ — im Beobachtungs-Register als **BEO-001** geführt, dort
seit Langem als „heute konsistent, aber ebenso ungeprüft“ benannt, und
inzwischen auf Zähler 2. Diese Welle liefert damit nicht eine vermutete, sondern
eine **belegte** Invariante.

**Der Register-Schritt ist ausdrücklich Teil des Wellen-Ziels:** BEO-001 wird
hier entweder geschlossen oder mit benannter Form weitergeführt. Die
Vorgänger-Welle hat gezeigt, was passiert, wenn dieser Schritt der nächsten
Welle überlassen wird.

## 2. Trigger (Welle startet)

Freigabe des Auftraggebers (2026-08-16), WIP-Slot frei (welle-74 geschlossen,
`in-progress/` trägt nur die Roadmap). Der Slice hat keine Vorgänger-Bedingung —
er ist von den `structure`- und Closure-Strängen unabhängig.

## 3. Closure-Trigger (Welle schließt)

- [slice-102](slice-102-wellen-lifecycle-invariante.md) liegt in
  `done/`.
- **Die Bestandsmessung liegt vor, bevor eine Aussage scharfgeschaltet wird** —
  je Aussage die Zahl der heutigen Verletzungen über die drei Repos. Eine
  Aussage, die den eigenen Bestand rot färbt, wird begründet oder fällt.
- **BEO-001 ist entschieden** — geschlossen oder mit benannter Form.
- Release als **Minor** (neue Fähigkeit, opt-in über neue Schlüssel; ohne sie
  byte-identisch), die Richtung in der Notiz **offen** formuliert.
- `make fullbuild` grün.

## 4. Slices in dieser Welle

| Slice | Rolle |
|---|---|
| [slice-102](slice-102-wellen-lifecycle-invariante.md) | Die vier Aussagen der Wellen-Invariante plus die gemessene Gegenrichtung von Aussage 4 |

**Ein Slice, alle vier Aussagen** — entschieden beim Wellen-Schnitt statt in
Abnahme-Punkt 1. Der Grund ist eine Änderung der Lage: die Aussagen 3 und 4
brauchen Tabellenzeilen, und die Tabellenzeilen-Lexik ist in welle-74 gerade
**entdriftet** worden (sie zählt jetzt nur außerhalb von Fences). Sie existiert
also, und sie braucht für diesen Slice nur eine **Spalten-Adresse** — nicht mehr
den Neubau, gegen den Abnahme-Punkt 1 ursprünglich abgewogen hat.

## 5. Abhängigkeiten

- [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  und [ADR-0028](../../adr/0028-planning-lifecycle-modul.md) liegen vor; diese
  Welle erweitert eine Anforderung, sie erfindet keine.
- Die **Tabellenzeilen-Lexik** aus
  [`DC-FA-TGT-001`](../../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
  ist der einzige bestehende Leser von Tabellenzeilen. Was hier entsteht, ist
  ihr **zweiter Konsument** — und nach
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 4 heißt das: geteilte Antwort, und je Konsument eine Assertion.
- Der eigene Bestand ist die Messgrundlage; die Roadmap dieses Repos ist
  zugleich Prüfling und Beispiel.

## 6. Out-of-Scope für diese Welle

- **Die Ordnungs-Bedingung aus BEO-005** (chronologische Tabellen kippen still
  ihre Richtung). Sie braucht dieselbe Tabellen-Lexik, ist aber eine Bedingung
  in einem **anderen** Modul und hängt nicht am Wellen-Register. Sie bekommt
  ihre eigene Welle — und dann hat die Lexik ihren dritten Konsumenten und
  verdient einen **Kopplungs-Test** statt Einzel-Assertionen.
- **Der Ausbau der Tabellen-Lexik zu einem adressierbaren Modell** (Spalten nach
  Kopfzeile, Typen, Ordnung). Hier entsteht nur, was die vier Aussagen brauchen.
- [slice-095](slice-095-links-resolve-from.md).

## 7. Closure-Notiz

Geschlossen am 2026-08-16 mit **v0.59.0**. Alle fünf Closure-Trigger sind
erfüllt: [slice-102](slice-102-wellen-lifecycle-invariante.md) liegt in `done/`,
die Bestandsmessung lag vor jeder Scharfschaltung, **BEO-001 ist entschieden**
(gestrichen — die eingetretene Instanz ist mechanisiert, der Rest benannt), das
Release samt Digest-Backfill ist draußen, und `make fullbuild` ist grün.

Die vollständige Notiz steht in [`welle-75-results.md`](welle-75-results.md).
