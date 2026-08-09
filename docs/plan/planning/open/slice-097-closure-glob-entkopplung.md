# Slice slice-097: `planning.closure.glob` — eigener Kandidaten-Filter für die Closure-Fähigkeit

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** Welle folgt, gemeinsam mit
[slice-098](slice-098-closure-note-placeholder.md) — der Konsument kann sein
letztes handgeschriebenes Prüfskript **erst dann** zurückziehen, wenn **beide**
Slices liegen. Das ist eine Closure-Bedingung jenseits der beiden Slice-DoDs und
damit genau der Fall, für den es eine Welle gibt.

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Closure-Fähigkeit, Spezifikation Schritt C2),
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md).
**Change Request** des Konsumenten `ai-harness-course` (CR 1).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Die Closure-Fähigkeit bekommt einen **eigenen** Kandidaten-Filter
(`planning.closure.glob`, Basisname-Glob, Default = `planning.slice-glob`).
Heute teilen sich zwei Fähigkeiten mit **verschiedenen Grundmengen** einen Knopf;
er ist nicht einstellbar, ohne eine von beiden zu verbiegen.

## 2. Der gemessene Befund

Die Aussagen sind verschieden: die Roadmap-Invariante fragt „liegt hier noch
Arbeit?", die Closure-Fähigkeit „ist jedes abgeschlossene Paket dokumentiert?".
Beide lesen heute `planning.slice-glob`.

Nachgemessen am 2026-08-09 gegen v0.52.0, Endzustand *kein Slice mehr in
Arbeit, Ruhe-Marker gesetzt* (muss grün sein):

```text
slice-glob "slice-*.md"  → 0 Befunde                    ✓
slice-glob "*.md"        → roadmap.md  planning-drift    ✗ falsch-rot
```

Ursache: die Roadmap-Datei liegt selbst im gezählten Verzeichnis und matcht
`*.md`. Die Invariante sieht dauerhaft „eine Datei liegt da" und verlangt, dass
der Ruhe-Marker **nie** gesetzt wird — falsch genau dann, wenn die Welle
wirklich ruht.

**Das ist ein Entwurfsfehler aus
[slice-093](../done/slice-093-closure-note-gate.md)**, nicht eine fehlende
Fähigkeit: die zweite Fähigkeit hat den Glob der ersten mitbenutzt, weil beide
zufällig Slice-Dateien meinten. Sobald eine Grundmenge weiter ist, bricht die
Kopplung.

## 3. Abnahme-Punkte

1. **Default = `planning.slice-glob`**, nicht ein eigener Literal-Default. Nur so
   ist der Befundsatz ohne den Schlüssel byte-identisch zu v0.52.0.
2. **Der eigene Bestand ist mitbetroffen — was tun wir damit?** In `done/` liegen
   **9** `welle-*-results.md`, jede mit einer Closure-Notiz; sie sind heute
   unsichtbar. Weitet man den Glob, greift sofort der nächste Unterschied: ihre
   Notiz-Überschrift ist eine **H1** (`# Welle NN — … — Closure-Notiz`), das
   Default-`heading-pattern` verlangt H2/H3 ⇒ neun `closure-note-missing`.
   Zu entscheiden: Muster auf `^#{1,3}` weiten, die Wellen-Notizen umbauen, oder
   den eigenen Glob bewusst eng lassen. **Der Slice liefert den Schlüssel; die
   eigene Konfiguration ist eine getrennte Entscheidung.**

## 4. Definition of Done

- [ ] `planning.closure.glob` in Lastenheft, Spezifikation (Schritt C2) und
      Config-Schema; Default = `planning.slice-glob`.
- [ ] Akzeptanzkriterien des CR als Tests: (1) ohne Schlüssel byte-identisch;
      (2) `closure.glob: "*.md"` bei unverändertem `slice-glob` prüft alle
      Markdown-Dateien im Closure-Verzeichnis, **ohne** die Roadmap-Invariante zu
      berühren; (3) fail-closed unverändert — leerer Glob ⇒ Exit 2, null
      Kandidaten unter dem gesetzten Glob ⇒ Befund.
- [ ] `make gates` + `make verify-closure-notes` grün; Abnahme-Punkt 2 beantwortet
      und die eigene Konfiguration entsprechend gesetzt.

## 5. Risiken / offene Punkte

- **Neun eigene Befunde bei geweitetem Glob** (siehe Abnahme-Punkt 2).
  — **Ausgang:** offen; die Entscheidung gehört in diesen Slice, die etwaige
  Sanierung der Wellen-Notizen nicht.
- **Zwei Globs, die fast immer gleich sind,** laden zu Drift ein (einer wird
  gepflegt, der andere nicht). — **Ausgang:** offen; der Default-Verweis statt
  eines Literals ist die Gegenmaßnahme, die geprüft werden muss.

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei. Unabhängig von
[slice-094](slice-094-closure-zaehl-paritaet.md) umsetzbar, sinnvoll aber
**danach** — sonst ändert sich die Zähl-Semantik unter den neuen Tests weg.

**Rückführungen:** `in-progress` → `open`, falls Abnahme-Punkt 2 eine
Bestands-Sanierung nach sich zieht, die eigenständig geschnitten gehört.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001**
  (Datei-Register driften unbemerkt gegen ihre Autoritäts-Tabelle). Andere
  Klasse — hier geht es um die Kandidaten-Menge **einer** Prüfung, dort um
  Referenzen zwischen Dokumenten. Nichts zu berücksichtigen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Change Request und Spezifikations-Schärfung
schreiben die Zusage, der Go-Code liefert sie.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
