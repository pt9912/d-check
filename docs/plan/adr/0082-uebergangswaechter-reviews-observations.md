# ADR-0082: Der Übergangs-Wächter bindet jetzt Review-Report-Deckung und Register-Deckung

**Status:** Accepted

**Datum:** 2026-09-03

**Autor:** pt9912

**Bezug:** [ADR-0081](0081-reviews-modul.md) §Re-Evaluierungs-Trigger („Zweiter
Trigger: Aufnahme in `gates`"), [ADR-0048](0048-closure-note-struktur-im-planning-modul.md)/[ADR-0077](0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md)
(`.d-check.closure.yml`, der Bindepunkt, der hier wächst), welle-86 <!-- d-check:status-provenance --> §1 (die
vier Vorbedingungen) und §3 (der Closure-Trigger, dessen Proben diesen Fund
auslösten).

**Schärft:** [`DC-FA-RVW-001`](../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
(keine inhaltliche Änderung — nur eine neue Bindung).

**Regeln:** Baseline-Regelwerk
[`modul-06-roadmap.md` §Wellen-Closure-Prozedur](../../../.harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md#wellen-closure-prozedur-modul-6),
Schritt 2 (Trigger-Audit).

---

## Kontext

Die Trigger-Prüfung von welle-86 <!-- d-check:status-provenance --> (Modul 6, Schritt 1) verlangt vier reale
Proben — je eine pro Vorbedingung, real als `git commit`-Versuch abgewiesen.
Zwei davon zeigten ein gemischtes Bild:

| Vorbedingung | Lokaler Hook (`verify-closure-notes`) | `make gates`/CI |
|---|---|---|
| DoD-Häkchen gesetzt | abgewiesen | abgewiesen |
| Beobachtungs-Register fortgeschrieben | **angenommen** | abgewiesen (`planning-check`, Hauptprofil) |
| Review-Report-Deckung | **angenommen** | **angenommen** |
| Bindung an den Übergang | abgewiesen (slice-175 <!-- d-check:status-provenance -->) | abgewiesen |

**Ursache, gemessen:** `.d-check.closure.yml` — das Profil, das der lokale
Hook über `verify-closure-notes` fährt — deklariert keinen
`planning.observations`-Block; die Register-Deckung läuft dort nicht mit,
nur im Hauptprofil. Und `reviews` ist laut ADR-0081 bewusst aus **jedem**
erzwungenen Lauf herausgehalten — weder Hook noch `make gates` aktivieren
es.

**Das trifft genau ADR-0081s eigenen Re-Evaluierungs-Trigger:**

> Zweiter Trigger: Aufnahme in `gates`. Sobald ein voller Wellen-Zyklus mit
> `review-coverage` grün gelaufen ist, ist die Aufnahme in `gates` fällig zu
> prüfen — dieselbe Reifung wie bei `workflows`.

welle-86 <!-- d-check:status-provenance --> ist dieser volle Wellen-Zyklus. `make review-coverage` lief während
seiner gesamten Laufzeit grün (mit der gemessenen Bestands-Ausnahme).

## Entscheidung

1. **`.d-check.closure.yml` bekommt einen `planning.observations`-Block**
   (Register + Verzeichnis, identisch zum Hauptprofil — dieselbe bewusste
   Duplikation, die der Datei-Kopf für `planning.heading`/`marker` bereits
   begründet: zwei Profile, eine Quelle wäre ohne zweite Herkunft nicht
   machbar).
2. **`.d-check.closure.yml` bekommt einen `reviews`-Block** (`done-dir`,
   `reviews-dir`, dieselbe Bestands-Ausnahme wie im Hauptprofil).
3. **`make verify-closure-notes` aktiviert zusätzlich `--enable reviews`.**
   Kein neues Modul, keine neue Prüf-Logik — dieselbe Wiederverwendung wie
   bei slice-175 <!-- d-check:status-provenance -->.
4. **`make review-coverage` bleibt als eigenständiges Ziel bestehen** — ein
   fokussierter Lauf ohne den restlichen Closure-Bindepunkt bleibt nützlich
   zum Debuggen, genau wie `make planning-check` neben `verify-closure-notes`
   bestehen bleibt, obwohl beide dieselbe Fähigkeit teilweise berühren.
5. **Kein Folge-ADR mit `supersedes` auf ADR-0081.** Dessen Entscheidung
   („noch nicht Teil von `gates`") war zum damaligen Zeitpunkt richtig und
   benannte selbst die Bedingung, unter der sie sich ändert. Diese ADR **löst
   die benannte Bedingung ein**, sie widerruft sie nicht.
6. **Keine Aufnahme von `reviews` in die unconditional `gates:`-Zusammensetzung
   selbst** (die zehn Glieder, die auf **jedem** Commit laufen, unabhängig von
   einer `done/`-Transition). Der Bindepunkt bleibt gezielt: `reviews` prüft
   eine Eigenschaft von `done/`-Slices, und genau dort — beim Übergang — soll
   sie greifen, nicht bei jeder beliebigen Änderung anderswo im Repo.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, Lücke offen lassen (welle-86 <!-- d-check:status-provenance --> mit Carveout schließen) | kein Eingriff | die zwei Lücken sind real und gemessen; ein Carveout ohne Not, wo eine Lösung eine Konfigurationszeile kostet |
| B — `reviews` in die unconditional `gates`-Zusammensetzung aufnehmen | maximale Abdeckung, jede Änderung geprüft | prüft `docs/reviews/`-Deckung bei jedem Commit, auch wenn kein `done/`-Slice betroffen ist — unnötige Laufzeit und eine Zusage, die weiter reicht als der Bindepunkt, der sie auslöste |
| **C — beide Lücken im Closure-Bindepunkt schließen (gewählt)** | passt genau auf den Bedarf, den die Proben gemessen haben; folgt dem bei slice-175 <!-- d-check:status-provenance --> etablierten Muster (Wiederverwendung, kein neues Modul) | zwei Config-Blöcke leben jetzt doppelt (Haupt- und Closure-Profil) — dieselbe, bereits akzeptierte Kosten-Klasse wie bei `planning.heading`/`marker` |

## Konsequenzen

- **Positiv, gemessen:** beide Proben (Register-Zitat ohne Zeile,
  fehlender Review-Report), die zuvor den lokalen Hook passierten, werden
  nach diesem Slice abgewiesen — nachgemessen, nicht nur behauptet.
- **Positiv:** ADR-0081s eigener Re-Evaluierungs-Trigger ist eingelöst,
  nicht stillschweigend liegen gelassen.
- **Negativ, benannt:** `verify-closure-notes` baut jetzt effektiv **vier**
  Fähigkeiten statt drei zusammen; ein Fund in einer von ihnen blockiert
  weiterhin den ganzen Übergang — dieselbe Eigenschaft, die slice-175 <!-- d-check:status-provenance -->
  schon hatte, jetzt mit einer Fähigkeit mehr.
- **Negativ, benannt:** die `reviews`-Bestands-Ausnahme aus ADR-0081 lebt
  jetzt an zwei Config-Stellen (Haupt- und Closure-Profil) — wächst sie,
  müssen beide nachgezogen werden. Kein Sensor hält das (dieselbe Grenze wie
  bei `planning.heading`/`marker`, benannt im Kopf von
  `.d-check.closure.yml`).

## Fitness Function (falls maschinell prüfbar)

Zwei reale Proben, real als `git commit`-Versuch gefahren (nicht simuliert),
nach der Umsetzung erneut gefahren:

- Ein Test-Slice mit zitierter, nicht registrierter `BEO-<NNN>` wird
  abgewiesen.
- Ein Test-Slice mit Review-Zusage ohne passenden Report wird abgewiesen.

Beide Proben-Ausgaben stehen im umsetzenden Slice, nicht hier (ADRs sind nach
`Accepted` immutabel).

## Re-Evaluierungs-Trigger

**Ein dritter Bindepunkt, an dem dieselbe Lücken-Klasse auftritt.** Sollte
ein künftiges `planning`- oder `reviews`-Feature erneut nur im Hauptprofil,
nicht im Closure-Profil landen, ist das ein drittes Vorkommen derselben
Klasse „zwei Config-Quellen, eine vergessen" — dann lohnt sich die Frage, ob
die Duplikation selbst (statt nur ihrer Instanzen) eine strukturelle Antwort
braucht.
