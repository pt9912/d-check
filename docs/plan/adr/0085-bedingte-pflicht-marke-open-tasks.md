# ADR-0085: Eine bedingte Pflicht-Marke koppelt `max-open-tasks` an eine Marke

**Status:** Accepted

**Datum:** 2026-09-17

**Autor:** pt9912

**Bezug:**
[der eingehende CR des Adopters `ai-harness-init`](../cr/2026-09-17-cr-eingehend-ai-harness-init-stilllegungs-form.md)
(Antrag und Beleg des Absenders),
[ADR-0074](0074-offene-tasks-auf-rohen-zeilen.md) (`max-open-tasks` selbst, die
Bedingung, an die diese koppelt),
[ADR-0075](0075-erklaerte-teilmenge-in-structure.md) (`hasMarker`-Form und
`require-all`, deren Erkennung diese ADR wiederverwendet),
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(die erweiterte Anforderung),
slice-225 <!-- d-check:status-provenance --> (der Träger, samt eigenem
Zweck: `CO-002` auflösen)

**Schärft:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v6.9.0/regelwerk/modul-04-adrs.md).

---

## Kontext

Baseline `v6.9.0`s `modul-05-planning-harness.md` führt einen vierten
Slice-Lifecycle-Zweig: Ein Slice, dessen Gegenstand ein anderer übernimmt oder
der ganz entfällt, geht **ohne Lieferung** nach `done/` — §7 trägt eine Zeile
`Gegenstand:` (`übernommen von slice-<Kennung>` | `entfallen: <Grund>`), und
die Liefer-Punkte der DoD bleiben **leer**. `max-open-tasks: 0` kennt diese
Form nicht: Sie verlangt ausnahmslos, dass jeder DoD-Haken gesetzt ist.

Dieses Repo hat den ersten Fall (slice-221 <!-- d-check:status-provenance -->) über einen blunten
`exempt-paths`-Eintrag der Regel überbrückt
([`CO-002`](../carveouts/CO-002-slice-221-gegenstand-entfallen.md)) — eine
namentliche Ausnahme pro Datei, die weder das Muster erkennt noch mit dem
Bestand mitwächst.

**Unabhängig davon traf ein eingehender CR eines zweiten Adopters
(`ai-harness-init`) ein**, der dieselbe Lücke von der anderen Seite misst:

> `planning.closure` liest den stillgelegten Slice wie jeden anderen in
> `done/`. Die Form der Stilllegung liest es nicht: Fehlt die Zeile
> `Gegenstand:`, bleibt der Lauf grün.

Der CR bittet um eine **bedingte Pflichtzeile**: Trägt ein Slice offene
Task-Items, muss sein Closure-Abschnitt eine Zeile `Gegenstand:` mit einem der
zwei Werte tragen — fehlt sie, ein **eigener** Befund. Das ist die
**Umkehrung** dessen, was dieses Repo für sich selbst braucht: Dort soll eine
vorhandene Zeile offene Punkte **erlauben** (Erlaubnis); der CR will, dass eine
**fehlende** Zeile bei offenen Punkten **auffällt** (Pflicht).

**Beide Hälften sind dieselbe Kopplung, nur mit vertauschtem Vorzeichen** — und
derselbe Code-Ort: `structureOpenTasks` zählt bereits die Überschuss-Items
gegen `max-open-tasks`; die Marke `**Gegenstand:**` ist exakt die Form, die
`hasMarker` (`require-all`, ADR-0075) bereits erkennt.

## Entscheidung

**1. Ein neuer, konditionierter Schlüssel: `open-tasks-require-marker`.**
Wirksam **nur**, wenn `max-open-tasks` für den geprüften Abschnitt bereits
mindestens einen Überschuss-Fund liefert — im Normalfall (alle Haken gesetzt)
bleibt die Kopplung wirkungslos, unabhängig davon, ob die Marke dasteht. Das
hält die Zusage klein: `max-open-tasks` bleibt die einzige Zählung, die Marke
entscheidet nur, **wie** ein bereits gefundener Überschuss gemeldet wird.

**2. Vorhandene Marke ⇒ alle Einzelbefunde entfallen (Erlaubnis).** Trägt der
bereinigte Abschnitts-Text die Marke, werden die sonst je Item gemeldeten
`section-tasks-open` vollständig unterdrückt. Das ist die Ablösung von
`CO-002`s `exempt-paths`-Eintrag: Die Erlaubnis hängt jetzt am **Inhalt** des
Abschnitts, nicht mehr an einer namentlich geführten Datei-Liste, die bei
jedem neuen Fall wieder anzufassen wäre.

**3. Fehlende Marke ⇒ EIN eigener Befund statt der vielen (Pflicht).**
`section-open-tasks-marker-missing`, `line` = Überschriftszeile. Kein
wiederverwendeter Code: Die Reparatur ist eine andere als bei
`section-tasks-open` (dort Haken setzen oder Punkt auflösen, hier die Marke
ergänzen), und dieselbe Regel — "jede Bedingung ihr eigener Grund-Code, damit
zwei Verletzungen nicht unter der Deduplikation zusammenfallen" — gilt hier
genauso. **Und es ersetzt, statt zu ergänzen:** ein Abschnitt mit drei offenen
Items und fehlender Marke trägt **einen** Befund, nicht vier. Der CR bittet
ausdrücklich um „einen eigenen Befund", nicht um einen zusätzlichen neben den
bestehenden — und vier Meldungen für dieselbe Ursache wären Rauschen, keine
zusätzliche Diagnose.

**4. `hasMarker`, keine neue Erkennungsform.** Dieselbe Funktion, die
`require-all` bereits nutzt: ein hervorgehobener Textlauf (`**M:**`,
`- **M:**`, `**M (Zusatz):**`) am Zeilen-Anfang des **bereinigten** Textes.
Eine zweite, ähnlich geschriebene Erkennung wäre die Falle, die ADR-0075 schon
einmal benannt hat (zwei RE2 einer Regel mit zwei verschiedenen Zielen) — hier
gibt es nicht einmal ein zweites RE2, weil `hasMarker` bereits die richtige
Form prüft.

**5. Halbe Aktivierung ⇒ Exit 2, wie bei `tasks-ignore-pattern`.** Der
Schlüssel **ohne** `max-open-tasks` ist eine Zusage, die nie greift — derselbe
Config-Rand wie bei jeder anderen Kopplung dieses Moduls
(`table.order-column` ohne `table.order`, `tasks-ignore-pattern` ohne
`max-tasks`, `exempt-expect-count` ohne `exempt-section-pattern`).

**6. Der CR nennt die falsche DC-ID, übernommen wird sie nicht.** Er schreibt
[`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
— in diesem Repos Schnitt lebt `max-open-tasks` und diese
Bedingung in [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(Modul `structure`). Beide Module laufen unter
`.d-check.closure.yml` nebeneinander, was die Verwechslung erklärt, sie aber
nicht rechtfertigt: Der zutreffende Bezug steht in `spec/lastenheft.md`s
eigenem Schnitt, nicht im Titel des CR.

## Verglichene Alternativen

**Ein neuer Grund-Code je fehlendem Item statt einem je Abschnitt.** Verworfen:
Der CR bittet ausdrücklich um „einen eigenen Befund" — Mehrzahl wäre eine
andere Zusage, und sie widerspräche der Absicht, `section-tasks-open`
(mehrere) durch eine klarere Diagnose zu **ersetzen**, nicht zu vermehren.

**`require-all` direkt verwenden, ohne Kopplung.** Verworfen: `require-all`
ist **unbedingt** — die Marke müsste in jedem Abschnitt stehen, auch wenn alle
Haken gesetzt sind. Das bräche jeden bestehenden, regulär gelieferten Slice
dieses Repos, der die Marke nie trägt.

**Die Marke geht in die Regel-Identität ein.** Verworfen, aus demselben Grund
wie bei `tasks-ignore-pattern`/`exempt-section-pattern`: Sie **verkleinert**
keine Menge und **grenzt** keine Regel gegenüber einer Schwester-Regel ab — sie
ändert nur, **wie** ein bereits gefundener Zustand gemeldet wird.

**Zwei getrennte Schlüssel** (einer für die Erlaubnis, einer für die Pflicht).
Verworfen: Beide Hälften prüfen dieselbe Marke gegen denselben
Überschuss-Zustand — ein zweiter Schlüssel für dieselbe Frage ist die
Verdopplung, die ADR-0070 für die Tabellen-Bedingungen bereits zurückgebaut
hat.

## Konsequenzen

**Positiv.** `CO-002`s `exempt-paths`-Eintrag wird auflösbar (Folge-Slice
slice-225 <!-- d-check:status-provenance -->):
Die Erlaubnis hängt jetzt am Inhalt, nicht an einer Datei-Liste. Der
eingehende CR eines zweiten, unabhängigen Adopters ist mit derselben
Mechanik beantwortet — zwei Bedürfnisse, eine Kopplung. Ohne den Schlüssel
byte-identisches Verhalten.

**Negativ, und das ist die ehrliche Seite.** Die Erkennung könnte zu weit
greifen — ein `Gegenstand:`-Feld, das aus anderem Anlass in einem Fließtext
auftaucht (nicht als Closure-Feld gemeint), könnte die Prüfung fälschlich
entschärfen. `hasMarker` prüft nur die **Form** (hervorgehobener Textlauf am
Zeilen-Anfang), nicht den **Ort** innerhalb des Abschnitts — ein Adopter, der
dieselbe Marke anderswo im Closure-Abschnitt verwendet, bekäme dieselbe
Erlaubnis unbeabsichtigt. Benannt in slice-225 <!-- d-check:status-provenance --> §6, nicht in dieser ADR
aufgelöst.

**Zwei Grund-Codes für denselben Rohzustand.** `section-tasks-open` und
`section-open-tasks-marker-missing` beschreiben beide „zu viele offene
Task-Items", nur mit unterschiedlicher Reparatur-Erwartung. Ein Leser, der nur
den Grund-Code kennt und nicht die Kopplung, könnte beide für unabhängige
Bedingungen halten.

## Fitness Function (falls maschinell prüfbar)

`make test` — `internal/hexagon/core/rules/structure_offene_tasks_test.go`
(bestehende Tests, unverändert) und ein neuer Test für die Kopplung: Erlaubnis
(Marke vorhanden ⇒ befundfrei), Pflicht (Marke fehlt ⇒ genau ein
`section-open-tasks-marker-missing`, nicht die Einzelbefunde), Normalfall
unberührt (keine Überschuss-Items ⇒ befundfrei unabhängig von der Marke),
abwesender Schlüssel (byte-identisch zum Vorzustand). Plus die zwei neuen
Config-Ränder in `configyaml_test.go`.

**Nicht maschinell geprüft** ist, ob die Marke am **richtigen Ort** im
Abschnitt steht (§Konsequenzen, negativ) — das bleibt Urteil.

## Re-Evaluierungs-Trigger

**Wenn die Marke wiederholt am falschen Ort trifft.** Dann ist die reine
Form-Prüfung (§Konsequenzen, negativ) zu weit, und eine Orts-Einschränkung
(z. B. nur als erste Zeile nach der Überschrift) wird fällig.

**Wenn ein zweiter Anwendungsfall dieselbe Kopplung an einer anderen
Bedingung braucht** (z. B. `max-tasks` statt `max-open-tasks`). Dann ist die
Bauform keine Einzelfall-Antwort mehr und verdient eine gemeinsame,
generischere Form statt eines zweiten, ähnlich benannten Schlüssels.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-17 | **Zählfehler in der eigenen Korrektur — dieselbe Klasse, die sie behebt.** Der vorige Eintrag behauptet „Vier neue Tests … (`TestOpenTasksRequireMarkerSection_*`)"; `grep -c '^func TestOpenTasksRequireMarkerSection_'` liefert **drei**, nicht vier (unabhängiger Review, Runde 2, R2-F-1). Register: [`commit-message-overclaims-work`](../planning/observations/BEO-ALL/commit-message-overclaims-work/observation.md) (13×, bereits verkörpert als Hard Rule — dieser Fund zeigt, dass die geschriebene Regel weiter verfehlt werden kann). Kein Code betroffen, nur die Zählangabe in diesem und im Slice-Plan |
| 2026-09-17 | **§Entscheidung 4 traf nur den Gleichschnitt-Fall — unabhängiger Review fand es empirisch.** `hasMarker(body, …)` in Entscheidung 1/4 sucht `body`, den bereinigten Text **desselben** Abschnitts, den `max-open-tasks` bereits zählt. §Kontext zitiert die Baseline korrekt (*„§7 trägt dann eine Zeile `Gegenstand:`"*) — aber §7 „Closure-Notiz" ist in jedem gemessenen Slice-Template ein **anderer** Abschnitt als die DoD-Sektion (§2), gegen die `.d-check.closure.yml`s Regel bindet. Reproduziert mit einer isolierten Fixture gegen das gebaute Image: eine Marke ausschließlich in „## 7. Closure-Notiz" bei offenen DoD-Haken löst `section-open-tasks-marker-missing` aus — **die Baseline-Ziel-Form selbst wird als Verstoß gegen sie gemeldet**. Der einzige Grund, warum `slice-221` <!-- d-check:status-provenance --> zuvor grün lief, war eine zweite, an keiner Stelle vorgeschriebene Kopie der Marke direkt unter der DoD-Checkliste. **Fix: `OpenTasksRequireMarkerSection`** (`open-tasks-require-marker-section`, RE2, dieselbe rohe Überschriften-Zeile wie `SectionPattern`) — abwesend sucht `hasMarker` weiterhin im gezählten Abschnitt (byte-identisch zu diesem Eintrag), gesetzt durchsucht `markerBody` die **Vereinigung** aller Abschnitte, deren Überschrift trifft (`FindSectionHeads`, dieselbe Erkennung wie jede andere Abschnitts-Suche des Moduls). Kein Treffer ⇒ Marke gilt als fehlend — dieselbe Lesart wie eine Datei ohne den Abschnitt. Vier neue Tests reproduzieren den Vorzustand und belegen den Fix (`TestOpenTasksRequireMarkerSection_*`), zusätzlich zwei neue Config-Rand-Tests. **Der Fitness-Function-Absatz unten galt zum Anlegen dieser ADR nur für den Gleichschnitt-Fall** — die Lücke war für Coverage-Metriken unsichtbar, weil „Marke fehlt im Abschnitt" und „Marke steht nur in einem anderen Abschnitt" denselben Codepfad durchliefen; **beide** Zustände sind jetzt einzeln getestet. `spec/lastenheft.md` 0.87.1, `spec/spezifikation.md` (neue Historie-Zeile, §2-Schema-Zeile, §4-Bedingungstabellen-Zeile aktualisiert) |
| 2026-09-17 | Angelegt, zusammen mit `open-tasks-require-marker` |
