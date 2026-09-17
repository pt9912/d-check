# Eingehender Change Request — `planning.closure` liest die Stilllegungs-Form nicht

**Absender:** Adopter `ai-harness-init` · **Eingegangen:** 2026-09-17
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (Modul `planning`, Fähigkeit `planning.closure`)
**Stand:** **offen** — eingegangen, nicht entschieden.

**Ablage-Hinweis.** Ein **eingehender** CR ist die dritte Klasse neben
[`MR-035`](../../../harness/conventions.md#mr-035) (ausgehend) und
[`MR-036`](../../../harness/conventions.md#mr-036) (Antworten darauf). Die Datei
liegt hier aus demselben Grund wie ihre Vorgänger in diesem Verzeichnis: Was
gebeten wurde und wie entschieden wird, soll im Repo stehen und nicht nur im
Vorgang.

---

## Wortlaut

> **Betreff:** Ein stillgelegter Slice in `done/` trägt eine Pflichtzeile, die
> kein Modul liest.
>
> **Einreicher:** ai-harness-init (Adopter), 2026-09-17, gegen d-check
> `@sha256:e31a372b66dbde26305982424854cfce7c9ab7ce555a94debeee7ee26e6d4641`
> (`v0.74.1`).
>
> ### Anlass
>
> Die Baseline `v6.9.0` des AI-Harness-Kurses führt eine neue Form ein:
> `modul-05-planning-harness.md`, Abschnitt *Ein Slice, dessen Gegenstand ein
> anderer übernimmt*. Ein solcher Slice geht ohne Lieferung über die Kanten
> `open → done` oder `next → done`. Die Liefer-Punkte seiner DoD bleiben
> offen, und §7 trägt die Zeile `Gegenstand:`, entweder *übernommen von*
> mit der Kennung des Nehmers oder *entfallen* mit einem Grund. Die Ziel-Fassung
> nennt es **urteilsfrei**, dass diese Zeile dasteht und einen der zwei Werte
> trägt.
>
> ### Messung
>
> An einer Kopie des Adopter-Repos, netzlos, über allen Modulen seiner
> `.d-check.yml`, je ein Slice über `open → done` und über `next → done`
> stillgelegt (Inhalt vor dem Move committet):
>
> | Lage im stillgelegten Slice | Meldung zu diesem Slice |
> |---|---|
> | Ziel-Form vollständig | keine |
> | §7 auf einen Satz gekürzt | `closure-note-thin` |
> | die Zeile `Gegenstand:` fehlt | **keine** |
> | `Gegenstand:` als unausgefüllter Vorlagen-Platzhalter, `placeholder: true` | `closure-note-placeholder` |
>
> `planning.closure` liest den stillgelegten Slice also wie jeden anderen in
> `done/`. Die Form der Stilllegung liest es nicht: Fehlt die Zeile
> `Gegenstand:`, bleibt der Lauf grün. Die Konfigurations-Vorlage
> (`--print-config`) führt keine Regel, die eine Zeile an eine Bedingung im
> selben Dokument knüpft. `structure` kennt `max-open-tasks` und
> `forbid-pattern`, aber keine Verknüpfung der beiden.
>
> ### Bitte
>
> Eine **bedingte Pflichtzeile** in `planning.closure`: Trägt ein Slice in
> `closure.dir` offene Task-Items in seiner DoD, muss sein Closure-Abschnitt
> eine Zeile `Gegenstand:` mit einem der zwei Werte tragen. Fehlt sie, ergibt
> das einen eigenen Befund.
>
> Wie die offenen Task-Items abgegrenzt werden (alle, oder nur die unter einer
> konfigurierten Überschrift wie *Liefer-Punkte*), ist Ihre Entscheidung. Der
> Adopter braucht nur, dass die Abgrenzung konfigurierbar ist, denn seine
> Slice-Vorlage trennt Liefer-Punkte von den übrigen DoD-Zeilen.
>
> ### Was ausdrücklich **nicht** gebeten ist
>
> - **Die Auflösung der genannten Kennung.** Ob der Nehmer existiert, lässt
>   die Ziel-Fassung selbst als Urteil oder eigenen Sensor offen.
> - **Die Frage, ob ein abgehakter Punkt ein Liefer-Punkt ist.** Das bleibt
>   Urteil.
>
> ### Wo der Adopter die Lücke führt
>
> Im Adopter-Repo steht die Lücke in der Sensor-Datei seines Doku-Gates
> (Ziel `docs-check`), Abschnitt *Ein stillgelegter Slice in `done/`*, mit
> dieser Datei als Adresse.
