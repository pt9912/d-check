# Carveouts — <Projektname>

> **Template-Hinweis.** Vorlage für `docs/plan/carveouts/README.md`. Kopiere
> nach `docs/plan/carveouts/README.md`, ersetze `<Platzhalter>` und lösche
> diesen Block. **Derivativ:** Quelle der Wahrheit sind die Carveout-Dateien;
> bei jedem neuen/aufgelösten Carveout mitziehen.

Aktive Carveouts mit Auflösungs-Trigger. Aufgelöste Carveouts wandern nach
`done/` (reiner `git mv`).

## Aktive Carveouts

| ID | Titel | Gate | Trigger | Folge-Slice |
|---|---|---|---|---|
| [CO-<NNN>](CO-<NNN>-<titel>.md) | <Kurztitel> | `<make-target>` | <Trigger> | `slice-<NNN>` |

## Aufgelöste Carveouts

(noch keine)

## Konventionen

- Jeder aktive Carveout braucht: Trigger, Folge-Slice, letzten Prüf-Termin.
- Bei Welle-Closure: Carveout-Audit zwingend — welche gültig, welche aufgelöst?
- Siehe Baseline-Regelwerk `modul-07-carveouts.md`.
