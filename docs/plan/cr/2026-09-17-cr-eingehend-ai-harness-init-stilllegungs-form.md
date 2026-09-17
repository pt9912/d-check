# Eingehender Change Request — `planning.closure` liest die Stilllegungs-Form nicht

**Absender:** Adopter `ai-harness-init` · **Eingegangen:** 2026-09-17
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (Modul `planning`, Fähigkeit `planning.closure`)
**Stand:** **entschieden und umgesetzt am 2026-09-17** — Bitte angenommen,
mit einer Abweichung von der Form (Belege unten); Träger
slice-225 <!-- d-check:status-provenance -->.

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

## Antwort

**Angenommen: eine bedingte Pflichtzeile, gekoppelt an die bestehende
Überschuss-Zählung.** Umgesetzt als neue, opt-in `structure`-Bedingung
`open-tasks-require-marker`
([ADR-0085](../adr/0085-bedingte-pflicht-marke-open-tasks.md),
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
0.87.0): trägt der geprüfte Abschnitt Überschuss-Task-Items gegen
`max-open-tasks`, ersetzt eine vorhandene Marke alle Einzelbefunde
(Erlaubnis — der eigene Bedarf dieses Repos), eine fehlende ersetzt sie
durch **einen** neuen Grund-Code `section-open-tasks-marker-missing`
(Pflicht — Ihre Bitte).

**Eine Abweichung von der beantragten Form, mit Grund.** Sie fragten nach
einer Bedingung, die eine **konkrete** Zeile `Gegenstand:` erkennt.
Umgesetzt ist eine **generische** Marken-Kopplung (`hasMarker`-Form,
dieselbe Erkennung wie `require-all`) — der Marken-**Name** ist
Konfiguration, nicht Produktkonstante. Grund: `Gegenstand` ist eine
Konvention der Baseline `v6.9.0`-Slice-Vorlage, keine Eigenschaft des
Werkzeugs, und eine hartkodierte Erkennung hätte jedem Adopter mit
abweichender Feld-Benennung nicht gedient.

**Die Abgrenzung der gezählten Task-Items bleibt unkonfiguriert.** Sie
merken an, dass Ihre Slice-Vorlage Liefer-Punkte von übrigen DoD-Zeilen
trennt, und lassen die Abgrenzung offen. Für den ersten Anwendungsfall
(dieses Repo, `slice-221`) genügt „alle offenen Task-Items" — die
Boilerplate-Haken sind zum Zeitpunkt der Stilllegung bereits gesetzt.
Eine engere Abgrenzung (z. B. eine zweite, unter einer benannten
Überschrift verankerte Zählung) ist **nicht** umgesetzt; folgt sie einem
zweiten gemessenen Bedarf, ist das ein eigener Schlüssel, keine Erweiterung
dieses.

**Nachtrag 2026-09-17 (nach unabhängigem Review, vor der ersten Closure
dieser Fähigkeit): die Marke lebt in einem eigenen Abschnitt — Ihre Bitte
zitiert das korrekt, der Erstentwurf oben prüfte es falsch.** Baseline
`v6.9.0` verortet `Gegenstand:` in „§7 Closure-Notiz", einem **anderen**
Abschnitt als dem, den `max-open-tasks` zählt (bei Ihnen vermutlich
ebenfalls die DoD-Sektion). Der Erstentwurf suchte die Marke ausschließlich
im gezählten Abschnitt selbst — ein Dokument, das exakt Ihrer eigenen
Baseline-Zitierung folgt (Marke nur in Closure-Notiz), hätte fälschlich
`section-open-tasks-marker-missing` erhalten. Unabhängiger Review fand das
empirisch, **bevor** dieser Slice geschlossen wurde. Behoben durch
`open-tasks-require-marker-section` (RE2, dieselbe Zeile wie
`section-pattern`): durchsucht **jeden** Abschnitt der Datei, dessen
Überschrift trifft, statt nur den gezählten. Für Ihre eigene Umsetzung
heißt das: **setzen Sie diesen zweiten Schlüssel** auf ein Muster, das
Ihren Closure-Abschnitt trifft (z. B. `^#{1,3} [0-9]+\. Closure-Notiz` oder
Ihre eigene Überschriftenform) — ohne ihn sucht die Bedingung weiterhin nur
im gezählten Abschnitt und trifft Ihre eigene Ziel-Form nicht.

**Die genannte DC-ID war nicht zutreffend.** Sie schreiben
[`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
— in diesem Repos Schnitt lebt `max-open-tasks` und die
neue Bedingung in
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(Modul `structure`, nicht `planning`) — beide Module laufen unter
`.d-check.closure.yml` nebeneinander, was die Verwechslung erklärt.

**Beleg.** `internal/hexagon/core/rules/structure_offene_tasks_test.go`:
vier Tests decken Erlaubnis, Pflicht (genau **ein** Befund statt
mehrerer), Normalfall unberührt (keine Überschuss-Items ⇒ Kopplung
wirkungslos, unabhängig von der Marke) und den abwesenden Schlüssel
(byte-identisches Verhalten); drei weitere Tests
(`TestOpenTasksRequireMarkerSection_*`) belegen die Abschnitts-Verlegung
aus dem Nachtrag oben — Erkennung in einem anderen Abschnitt, fehlender
benannter Abschnitt gilt als fehlende Marke, abweichende
Überschriften-Nummerierung. `slice-221` selbst ist der lebende Beleg:
`make verify-closure-notes` läuft grün gegen seine `**Gegenstand:**`-Zeile
in „## 7. Closure-Notiz", ohne einen namentlichen `exempt-paths`-Eintrag.
