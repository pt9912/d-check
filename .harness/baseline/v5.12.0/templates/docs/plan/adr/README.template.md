# ADR-Index — <Projektname>

> **Template-Hinweis.** Vorlage für `docs/plan/adr/README.md`. Kopiere nach
> `docs/plan/adr/README.md`, ersetze `<Platzhalter>` und lösche diesen Block.
> **Derivativ:** Quelle der Wahrheit sind die ADR-Dateien; dieser Index ist
> eine Bequemlichkeits-Sicht — bei jedem neuen/akzeptierten ADR mitziehen.

| ID | Titel | Status | Bezug |
|---|---|---|---|
| [<NNNN>](<NNNN>-<titel>.md) | <Titel der Entscheidung> | Proposed \| Accepted | `<LH-FA-NN>` |

## Konventionen

- ADRs sind nach `Accepted` **immutable** (siehe Baseline-Regelwerk `modul-04-adrs.md`).
- Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN`.
- Bei `Accepted`: diesen Index aktualisieren (Status, Datum).
- Jede ADR deklariert im `**Schärft:**`-Feld *aufwärts*, welche Spec-Stelle
  sie verbindlich macht (Baseline-Regelwerk `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)) —
  als Kennung (`SPEC-*`, `ARC-*`, `<PREFIX>-FA-*.<Buchstabe>`), ersatzweise als
  Abschnitt, wo die Sektion keine Kennungen vergibt.
  Prozess-ADRs ohne Spec-Stratum tragen `—`.
