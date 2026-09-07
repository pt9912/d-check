# ADR-NNNN: <Titel der Entscheidung>

> **Template-Hinweis.** Diese Datei ist eine Vorlage im MADR-/Nygard-
> Stil. Kopiere sie nach `docs/plan/adr/<NNNN>-<kurzer-titel-kebab>.md`
> — bzw. `<BEREICH>-<NNNN>-…`, wenn dein Repo den Zählraum je Sub-Area
> führt (Deklaration in `harness/conventions.md`) —
> und ersetze alle `<Platzhalter>`. Lösche diesen Block nach dem
> Ausfüllen. Vergiss nicht, den ADR-Index in
> `docs/plan/adr/README.md` zu aktualisieren.

**Status:** Proposed | Accepted | Deprecated | Superseded by ADR-NNNN

**Datum:** YYYY-MM-DD

**Autor:** <Name>

**Bezug:** [`<LH-FA-NN>`](../../../spec/lastenheft.md#<anker>), [`<LH-QA-NN>`](../../../spec/lastenheft.md#<anker>), [ADR-<NNNN>](<NNNN>-<titel>.md) (optional)

**Schärft:** [`<SPEC-NNN>`](../../../spec/spezifikation.md#<anker>) / [`<PREFIX>-FA-<NN>.<Buchstabe>`](../../../spec/spezifikation.md#<anker>) / [`<ARC-NNN>`](../../../spec/architecture.md#<anker>) — welche
Spec-Stelle diese ADR verbindlich macht. **Die Kennung nennen, wo das
Zielelement eine trägt** — `SPEC-*`, die Verfeinerung `<PREFIX>-FA-<NN>.<Buchstabe>`
oder `ARC-*`; erst wo eine Sektion keine vergibt, den Abschnitt
selbst (`<architecture.md §N>`). Steht die Kennung in einer **Überschrift**,
zeigt der Link auf sie; steht sie in einer **Tabellenzelle**, auf die Sektion —
dort hat sie keinen eigenen Anker
(Baseline-Regelwerk `grundlagen-source-precedence.md` §ID-Schema als Klammer).
Aufwärts-Deklaration der
Änderungskopplung: wer diese ADR ändert, zieht von hier die betroffenen
Spec-Stellen nach. `—` eintragen, wenn Prozess-ADR ohne Spec-Stratum.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

> **IDs als Markdown-Link** (klickbar zur Quelle, Baseline-Regelwerk `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)).
> Der `<anker>` ist der GitHub-Heading-Slug der Ziel-Überschrift. Der
> Der Referenz-Richtungs-Gate prüft heute nur Token-Richtung, **nicht** die
> Anker-Auflösung — ein umbenannter Abschnitt rottet den Link still. Die
> anker-validierende Reifestufe löst Links auf und prüft die Anker-Existenz
> am Zielknoten.

---

## Kontext

<!--
Was ist die Ausgangslage? Welche Anforderung oder welcher Druck führt
zu dieser Entscheidung? Welche Annahmen gelten? Wenn diese Annahmen
kippen, kippt die Entscheidung.
-->

<…>

## Entscheidung

<!--
Die Wahl, in einem Satz oder einem kurzen Absatz. Eindeutig, ohne
"vielleicht", "könnte".
-->

Wir wählen **<Variante X>**.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

<!--
Mindestens drei Optionen mit Pro/Contra. Alternativ "nichts tun" ist
auch eine Option.
-->

| Option | Pro | Contra |
|---|---|---|
| A — <Bezeichnung> | <…> | <…> |
| B — <Bezeichnung> | <…> | <…> |
| **C — <gewählt>** | <…> | <…> |

## Konsequenzen

<!--
Was folgt aus der Entscheidung? Sowohl Positives als auch Schmerzen.
Was wird leichter, was schwerer.
-->

- Positiv: <…>
- Negativ: <…>
- Folgepflicht: <Fitness Function, Doku-Update, Folge-Slice>

## Fitness Function (falls maschinell prüfbar)

<!--
Wenn die Entscheidung sich in einer prüfbaren Eigenschaft des Codes
niederschlägt: hier die konkrete Regel benennen. Beispiel:
"depguard verbietet Import von internal/runtime aus internal/service."
-->

| Tooling | Regel | Make-Target |
|---|---|---|
| <z.B. depguard> | <…> | `make arch-check` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

<!--
Wann sollte diese Entscheidung erneut geprüft werden?
"Wenn Bibliothek X v2 verfügbar ist." "Wenn Kostenbudget Y überschritten."
"Bei Meilenstein M3."
-->

<…>

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| YYYY-MM-DD | Proposed | <Slice-Datei> |
| YYYY-MM-DD | Accepted | <PR-Link> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
