# Welle welle-68-planning-roadmap-harness: Planning-/Roadmap-Harness vollenden

**Lifecycle:** Diese Datei liegt **flach** unter `docs/plan/planning/` solange die Welle
läuft; bei Closure wandert sie per `git mv` nach `done/` (neben ihre
`welle-68-results.md`). Der Zustand ist die **Verzeichnis-Position** — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (interne Harness-Vertiefung; das
`verify-closure-notes`-Gate wird released).

**Verantwortlich:** pt9912. **Datum:** 2026-08-03.

---

## 1. Welle-Ziel

Die Planning-/Roadmap-Harness-Schicht nach der v5.0.0-Migration **vollenden**:
(a) `## Aktuelle Welle` trägt — während eine Welle läuft — die
**Baseline-Struktur-Felder** (Welle-ID · Start · Geplantes Ende · Closure-Trigger), womit
der Abschnitt die Template-Form erreicht **ohne** `planning`-Modul-Umbau (der Ruhe-Marker
gilt nur noch im wellenlosen Zustand); und (b) der **Closure-Note-Qualitäts-Nachlauf** ist
mechanisiert — ein `verify-closure-notes`-Gate (Struktur) plus der inferentielle
`closure-note-reviewer.md`-Skill (Inhalt vs. Floskel).

## 2. Trigger (Welle startet)

- welle-67-baseline-v500-migration liegt in `done/` (die v5.0.0-Migration ist komplett) —
  der Start-Trigger liegt **vor** der Welle, ist kein Ergebnis dieser Welle.

## 3. Closure-Trigger (Welle schließt)

- Alle Slices dieser Welle (092–093) liegen in `done/`.
- `make gates` grün **und** das **neue** `verify-closure-notes`-Gate grün über den
  `done/`-Bestand (repo-weit — steht in keiner einzelnen Slice-DoD; das ist das *Mehr*
  gegenüber den Slice-DoDs).
- Trigger-Audit durchlaufen (`modul-06` Closure-Schritt 2).
- Closure-Notiz `done/welle-68-results.md` geschrieben.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [slice-092](done/slice-092-roadmap-aktuelle-welle-template-form.md) | `## Aktuelle Welle` auf die Template-Struktur-Felder (aktive-Welle-Form) + Adaption verfeinern | Roadmap-Template-Konformität |
| [slice-093](done/slice-093-closure-note-gate.md) | Closure-Note-Reviewer-Skill + `verify-closure-notes`-Gate | Etappe-D-Finding D-7 |

## 5. Abhängigkeiten

- **slice-092 zuerst** (macht `## Aktuelle Welle` template-konform, sobald welle-68
  aktiv ist); **slice-093** danach (die Gate-Adoption, Produkt-Code + eigene ADR).
- **Wird blockiert von:** nichts Externem.

## 6. Out-of-Scope für diese Welle

- **RTM-Generator** und **`--print-version-md`-Scaffold** — separate CLI-Kandidaten
  (Roadmap §Nächste Wellen), nicht Teil dieser Planning-/Harness-Welle.
- **Kein `planning`-Modul-Code-Umbau:** die aktive Welle trägt die Struktur-Felder ohne
  Ruhe-Marker (`hasActive == hasSlices` grün); der Ruhe-Marker bleibt die wellenlose Form.

## 7. Closure-Notiz

_Ausstehend (wird bei der Welle-Closure gefüllt; die Pointer lösen vom Ruheort `done/` auf)._
