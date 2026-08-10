# Welle welle-69-structure-schnitt: Struktur-Invarianten — Schnitt und Ablöse-Pfad

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-69-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Werkzeug-Fähigkeit auf
Konsumenten-Antrag).

**Verantwortlich:** pt9912. **Datum:** 2026-08-09.

---

## 1. Welle-Ziel

Entscheiden **und vertraglich festhalten**, wie d-check **Struktur-Invarianten
innerhalb eines Dokuments** bekommt — und wie die in v0.52.0 ausgelieferte,
planning-lokale Closure-Fähigkeit darin aufgeht, ohne Adopter-Konfigurationen zu
brechen.

Der Modulsatz deckt heute **Referenz**-Invarianten lückenlos ab (Ziel existiert,
Kennung verlinkt, Richtung erlaubt, Target deklariert, Core unverändert), aber
keine einzige Aussage über die **Form eines Dokuments selbst**. Das ist keine
Nachlässigkeit, sondern eine nie ausgesprochene Grenze: die Module sind entlang
„zeigt dieses Dokument korrekt auf andere?" gewachsen, nie entlang „ist dieses
Dokument selbst richtig gebaut?".

## 2. Trigger (Welle startet)

- Der Change Request „Struktur-Invarianten über Dokumentklassen" aus dem
  Schwester-Repo a-check liegt vor und ist freigegeben (2026-08-09). Der Trigger
  liegt **vor** der Welle und ist kein Ergebnis von ihr.
- Kein WIP-Konflikt: `docs/plan/planning/in-progress/` trug beim Start keinen
  Slice.

## 3. Closure-Trigger (Welle schließt)

- Alle Slices dieser Welle liegen in `done/`.
- **Das Mehr gegenüber der Slice-DoD:** die drei bereits liegenden
  Closure-Slices ([slice-094](../done/slice-094-closure-zaehl-paritaet.md),
  [slice-097](../done/slice-097-closure-glob-entkopplung.md),
  [slice-098](../done/slice-098-closure-note-placeholder.md)) haben eine
  **Entscheidung** — eigenständig, aufgegangen oder neu zugeschnitten — und, wo
  sie bestehen bleiben, eine Wellen-Zuordnung. Solange das offen ist, ist der
  Schnitt nicht fertig, auch wenn der Slice-Plan abgehakt ist.
- Die Folge-Slices der Umsetzung existieren als Dateien im Planning-Lifecycle —
  genannt ohne angelegt wäre dieselbe Klasse wie ein halluziniertes Gate.
- `make gates` grün; Trigger-Audit durchlaufen (`modul-06` Closure-Schritt 2).
- Closure-Notiz `done/welle-69-results.md` geschrieben.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [slice-096](slice-096-structure-modul-analyse.md) | Modul `structure` — Analyse, Modul-Schnitt und Ablöse-Pfad | [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in), [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) |

**Genau ein Slice** — die Implementierung ist nach §6 ausdrücklich
**Out-of-Scope** dieser Welle. Die aus slice-096 geschnittenen Folge-Slices
stehen deshalb **nicht** hier, sondern derivativ in der Ergebnis-Notiz: sie in
§4 zu führen hieße, den eigenen Closure-Trigger („alle Slices dieser Welle in
`done/`") unerfüllbar zu machen.

## 5. Abhängigkeiten

- **Blockiert:** die Folge-Welle mit den Closure-Slices — deren Zuschnitt ist ein
  Ergebnis dieser Welle, nicht ihre Voraussetzung.
- **Wird blockiert von:** nichts. [slice-095](../open/slice-095-links-resolve-from.md)
  (`links.resolve-from`) läuft unabhängig daneben und gehört nicht hierher.

## 6. Out-of-Scope für diese Welle

- **Die Implementierung des Moduls.** Diese Welle liefert Schnitt, Vertrag und
  Ablöse-Pfad; Go-Code, Paritäts-Beleg und Release sind Folge-Slices.
- **Die Ablösung fremder Skripte.** Sie ist der *Anlass*, nicht das Ergebnis —
  d-check liefert eine Fähigkeit, das Zurückziehen entscheidet der Adopter.
- **Nicht-Markdown-Quellen.** Der Antragsteller hat selbst gemessen, dass zwei
  seiner fünf Pin-Stellen außerhalb von Markdown liegen; das bleibt außerhalb
  von d-checks Gegenstand.
- **`links.resolve-from`** — eigener Antrag, eigener Slice, keine gemeinsame
  Closure-Bedingung.

## 7. Closure-Notiz

Geschlossen am 2026-08-09. Die Ergebnis-Notiz der Welle steht — der Baseline-Form
folgend — in einer **eigenen** Datei neben dieser: [`welle-69-results.md`](welle-69-results.md). Sie
trägt, was geliefert wurde, was funktioniert hat, was anders lief, den
Lese-Schritt am Beobachtungs-Register und die Verifikation.

Der Trigger-Audit dieser Welle hat einen eingetretenen
Re-Evaluierungs-Trigger von [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md)
festgestellt; er ist ohne Supersede als Verfeinerung in
[ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md) beantwortet.

Diese Plan-Datei hält nur noch fest, **dass** die Welle geschlossen ist; ihr
Zustand ist die Verzeichnis-Position.
