# Welle welle-72-closure-semantik: Die Semantik eines ausgelieferten Gates nachziehen

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-72-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Semantik-Angleichung an
ausgeliefertem Verhalten).

**Verantwortlich:** pt9912. **Datum:** 2026-08-10.

---

## 1. Welle-Ziel

Beide Slices ändern die **Semantik des ausgelieferten** Closure-Gates, statt
etwas hinzuzufügen. Das ist die Klammer dieser Welle — und der Grund, warum sie
in welle-71 bewusst **nicht** mitliefen: dort war alles additiv (Default-Verweis,
Default `false`), ein Adopter merkte beim Update nichts.

Hier merkt er etwas. Ziel der Welle ist deshalb mehr als die Summe der beiden
DoDs: **eine** Release-Notiz, die einem Konsumenten in **einem** Zug sagt, was
sich an seinem Gate ändert — statt zweier aufeinanderfolgender Releases, die
jeweils ein Stück verschieben.

## 2. Trigger (Welle startet)

Freigabe des Auftraggebers (2026-08-10) und ein freier WIP-Slot. Beide
Slice-Trigger sind erfüllt: [slice-096](slice-096-structure-modul-analyse.md)
liegt in `done/` und hat den Zuschnitt entschieden — beide bleiben eigenständig.

## 3. Closure-Trigger (Welle schließt)

- Beide Slices dieser Welle liegen in `done/`.
- **Ein** Release trägt beide Änderungen, und seine Notiz nennt **jede** Richtung
  einzeln (siehe §5).
- Der eigene Bestand ist **vor** jeder Änderung gemessen und danach grün.
- `make fullbuild` grün.

## 4. Slices in dieser Welle

| Slice | Rolle | Richtung |
|---|---|---|
| [slice-094](slice-094-closure-zaehl-paritaet.md) | Zähl-Parität: Inline-Code zählt nicht, Satzende nur vor Whitespace | **beide** |
| [slice-104](slice-104-floskel-wortgrenze.md) | Floskel-Vergleich an der Wortgrenze statt als Teilstring | findet **weniger** |

**Reihenfolge: slice-094 zuerst.** Es trägt die ADR, die den Rahmen setzt
(eine Prüfregel wird bewusst gelockert,
[`AGENTS.md` §3.6](../../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)).
slice-104 ist danach die zweite, kleinere Lockerung **derselben** Prüfung und
kann sich auf den gesetzten Rahmen beziehen.

## 5. Die Kopplung, die keiner der beiden Slices nennt

**Beide lockern die Floskel-Prüfung — aus verschiedenen Richtungen.**

- slice-094 entfernt Inline-Code aus dem bereinigten Abschnittstext. Die
  Floskel-Prüfung liest **denselben** Text; eine Phrase in Backticks wird danach
  nicht mehr gefunden.
- slice-104 vergleicht an Wortgrenzen. Eine Phrase, die heute **innerhalb** eines
  Wortes trifft, wird danach nicht mehr gefunden.

Einzeln betrachtet ist jede Lockerung begründet und sachlich richtig. Zusammen
sind es **zwei** Verengungen derselben Prüfung in **einem** Release, und genau
deshalb gehören sie in dieselbe Welle: ein Konsument soll einmal lesen, was seine
Floskel-Liste danach noch findet — nicht zweimal ein Stück davon.

Die Gegenrichtung von slice-094 (`closure-note-thin` wird **schärfer**) steht
daneben und muss in derselben Notiz erscheinen.

## 6. Abhängigkeiten

- Beide hängen an [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  Schritt C4 und an [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md).
- **Paritäts-Fixtures:** slice-094 will sie aus dem Schwester-Repo `a-check`
  **beiziehen**, nicht nachbauen. Beim Wellen-Start geprüft: das Repo liegt vor
  (87 abgeschlossene Slices). Ist es beim Slice-Start nicht nutzbar, verengt sich
  die Zusage auf eigene Fixtures — so steht es im Slice.

## 7. Out-of-Scope für diese Welle

- **Die Abschnitts-Anzahl.** Das Adopter-Skript prüft zusätzlich, wie viele
  Closure-Abschnitte ein Dokument trägt; d-check liest laut Spezifikation nur den
  ersten. Diese Lücke ist ausdrücklich **nicht** hier zu schließen.
- **Eine Sanierung des eigenen Bestands**, falls die Vorab-Messung Notizen rot
  macht. Das wäre ein Befund über die Notizen und ein eigener Slice — kein Grund,
  die Zählung aufzuweichen.
- **Das Modul `structure`** ([slice-099](../open/slice-099-structure-modul.md)) —
  eigene Welle.

## 8. Closure-Notiz

Geschlossen am 2026-08-10. Die Ergebnis-Notiz steht — der Baseline-Form folgend —
in einer **eigenen** Datei daneben: [`welle-72-results.md`](welle-72-results.md).

Das Wellen-Ziel jenseits beider Slice-DoDs ist eingelöst: **ein** Release
(**v0.56.0**) mit **einer** Notiz, die alle drei Richtungen einzeln nennt.

Diese Plan-Datei hält nur noch fest, **dass** die Welle geschlossen ist; ihr
Zustand ist die Verzeichnis-Position.
