# Welle welle-71-closure-konsumenten-paritaet: Die Closure-Fähigkeit wird Obermenge des Konsumenten-Skripts

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-71-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Konsumenten-Change-Requests).

**Verantwortlich:** pt9912. **Datum:** 2026-08-10.

---

## 1. Welle-Ziel

Der Konsument `ai-harness-course` prüft seine Closure-Notizen heute noch mit
einem **handgeschriebenen Skript**, weil d-checks Closure-Fähigkeit zwei Lücken
lässt. Ziel dieser Welle ist, dass er es **zurückziehen** kann.

Das ist ausdrücklich **mehr als die Summe der beiden Slice-DoDs**: jeder für sich
bringt Deckung, aber erst beide zusammen machen die Fähigkeit zur Obermenge des
Skripts. Genau dafür gibt es eine Welle — die Bedingung liegt jenseits dessen,
was ein einzelner Slice zusagen kann.

## 2. Trigger (Welle startet)

Freigabe des Auftraggebers (2026-08-10) und ein freier WIP-Slot. Beide
Slice-Trigger sind erfüllt: [slice-096](done/slice-096-structure-modul-analyse.md)
liegt in `done/`, und dessen Schnitt hat entschieden, dass beide Slices
**eigenständig** bleiben (nicht im `structure`-Modul aufgehen) — der Vorbehalt in
den Slice-Köpfen ist damit aufgelöst.

## 3. Closure-Trigger (Welle schließt)

- Beide Slices dieser Welle liegen in `done/`.
- **Die Konsumenten-Bedingung ist belegt:** die Closure-Fähigkeit deckt das
  handgeschriebene Prüfskript des Konsumenten vollständig ab — nachgewiesen an
  seinem realen Bestand, nicht an einem Fixture.
- `make fullbuild` grün; Release erfolgt, sofern die Slices Produkt-Code liefern.

## 4. Slices in dieser Welle

| Slice | Rolle |
|---|---|
| [slice-097](done/slice-097-closure-glob-entkopplung.md) | Eigener Kandidaten-Filter `planning.closure.glob` — die Fähigkeit bekommt ihre **Grundmenge** zurück, die sie heute mit `planning.slice-glob` teilt |
| [slice-098](in-progress/slice-098-closure-note-placeholder.md) | Neuer opt-in-Grund-Code `closure-note-placeholder` — unausgefüllte Template-Rümpfe |

**Reihenfolge:** zuerst slice-097. Beide Slices halten fest, dass keine
Reihenfolge nötig ist; der Kandidaten-Filter bestimmt aber die Menge, auf der die
Platzhalter-Erkennung dann misst — in dieser Richtung ist die zweite Messung
aussagekräftiger. WIP-Limit 1, also ohnehin nacheinander.

## 5. Abhängigkeiten

- Beide hängen an [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  und [ADR-0048](../adr/0048-closure-note-struktur-im-planning-modul.md); ob [ADR-0048](../adr/0048-closure-note-struktur-im-planning-modul.md)
  verfeinert oder eine Folge-ADR nötig wird, entscheidet sich je Slice.
- Keine Abhängigkeit zu [slice-094](open/slice-094-closure-zaehl-paritaet.md) —
  siehe Out-of-Scope.

## 6. Out-of-Scope für diese Welle

- **[slice-094](open/slice-094-closure-zaehl-paritaet.md) (Zähl-Parität) bleibt
  draußen**, obwohl die Roadmap alle drei zusammen geplant hatte. Der Grund ist
  die **Risiko-Klasse**, nicht der Umfang: 097 und 098 sind rein **additiv**
  (neuer Config-Schlüssel mit rückwärtskompatiblem Default; opt-in-Grund-Code mit
  Default `false`) — ein Adopter merkt beim Update nichts. slice-094 ändert
  dagegen die **Zähl-Semantik eines ausgelieferten Gates**: `closure-note-thin`
  meldet danach anders, und zwar in beide Richtungen. Das ist dieselbe Klasse,
  die in welle-70 eine eigene konsumentensichtbare Zeile gekostet hat, und sie
  verdient eine eigene Welle mit eigener Bestandsmessung — nicht ein Mitfahren
  im Windschatten zweier additiver Slices.
- Die `structure`-Umsetzung ([slice-099](open/slice-099-structure-modul.md)) —
  eigene Welle, Trigger eingetreten, Freigabe offen.

## 7. Closure-Notiz

_Ausstehend._
