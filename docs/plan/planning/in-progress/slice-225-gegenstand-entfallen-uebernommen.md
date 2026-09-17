# Slice slice-225: Der Lifecycle-Zweig „Gegenstand entfallen/übernommen" wird gate-tragfähig

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`MR-072`](../../../../harness/conventions.md#mr-072) (Delta-Messung
v6.6.0→v6.9.0, Punkt 1: neuer vierter Slice-Lifecycle-Zweig — dieser Slice
ist die dort angekündigte Adoption), [`CO-002`](../../carveouts/CO-002-slice-221-gegenstand-entfallen.md)
(dessen Auflösung dieser Slice trägt), [eingehender CR des Adopters
`ai-harness-init`](../../cr/2026-09-17-cr-eingehend-ai-harness-init-stilllegungs-form.md)
(seit 2026-09-17 in diesen Plan aufgenommen — Auftraggeber-Entscheid, siehe
§1 unten; der CR nennt selbst [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in), gemeint ist nach diesem
Repos Schnitt [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) —
`max-open-tasks` und die neue Bedingung leben im Modul `structure`, nicht
`planning`; `.d-check.closure.yml` aktiviert beide Module nebeneinander,
was die Verwechslung erklärt, sie aber nicht auflöst).

**Berührte Spec-Stellen:** [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(neue Bedingung, siehe §1) — **erweitert um den CR, ursprünglich `—`**: der
Ursprungs-Plan (2026-09-08) sah nur eine `.d-check.closure.yml`-Anpassung
ohne Produkt-Berührung vor; das ist mit der CR-Aufnahme nicht mehr richtig.

**Verantwortlich:** pt9912 (Implementer-Rolle).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Plan-Änderung (2026-09-17, bei der Beanspruchung).** Der ursprüngliche
Plan (2026-09-08, unten in den Abschnitten erhalten) deckte nur die
**Erlaubnis-Hälfte**: eine `Gegenstand:`-Zeile lässt offene Liefer-Punkte
zu. Am selben Tag traf ein eingehender CR des Adopters `ai-harness-init`
ein
([`docs/plan/cr/2026-09-17-cr-eingehend-ai-harness-init-stilllegungs-form.md`](../../cr/2026-09-17-cr-eingehend-ai-harness-init-stilllegungs-form.md)),
der die **Pflicht-Hälfte** bittet: Trägt ein Slice offene Liefer-Punkte,
muss die `Gegenstand:`-Zeile mit einem der zwei Werte **vorhanden** sein,
sonst ein **eigener** Befund — nicht nur das generische `section-tasks-open`,
das heute schon feuert, wenn die Zeile fehlt und Punkte offen sind. Beide
Hälften sind derselbe Code-Ort (dieselbe `structure`-Regel derselben Sektion)
und derselbe Bruch-Test-Rahmen; getrennt zu implementieren hieße, dieselbe
Erkennungsform zweimal zu schreiben. Auftraggeber-Entscheid: **aufnehmen**,
nicht als eigenen Folge-Slice abtrennen.

**Ziel.** Baseline `v6.9.0`s `modul-05-planning-harness.md` führt einen
vierten Slice-Lifecycle-Zweig: Ein Slice, dessen Gegenstand ein anderer
Slice übernimmt oder der ganz entfällt, geht **ohne Lieferung** nach
`done/` — §7 trägt dann eine Zeile `Gegenstand:` (`übernommen von
slice-<Kennung>` | `entfallen: <Grund>`), und die Liefer-Punkte der DoD
bleiben **leer**. Dieser Slice macht `make verify-closure-notes` (Modul
`structure`) mit dieser Form **beidseitig** kompatibel: **(a)** eine
vorhandene `Gegenstand:`-Zeile lässt offene Liefer-Punkte zu (Erlaubnis;
der ursprüngliche Zweck), **und (b)** offene Liefer-Punkte **ohne** eine
solche Zeile ergeben einen eigenen, benannten Befund statt des heutigen
pauschalen `section-tasks-open` (Pflicht; der CR-Zweck) — heute kennt
`max-open-tasks: 0` nur die ausnahmslose Zahl, ohne Rücksicht auf einen
`Gegenstand:`-Ausgang in beide Richtungen.

**Der Anlass ist [slice-221](../done/slice-221-agents-md-tabellenzellen.md)**,
das erste Repo-Beispiel dieses Zwecks (Ausgang „entfallen" — slice-222 hat
seinen Gegenstand durch eine andere Lösung erledigt) und deshalb per
[`CO-002`](../../carveouts/CO-002-slice-221-gegenstand-entfallen.md) statt
regulärer Closure nach `done/` gewandert; der CR bestätigt unabhängig
(anderer Adopter, andere Messung), dass dieselbe Lücke auch dort auffällt.

**Abgrenzung — fünf Punkte, jeder mit Grund:**

1. **Kein Retrofit auf `done/`-Altbestand.** Kein bisheriger Slice trägt
   die neue `Gegenstand:`-Form; sie gilt ab Einführung, wie jede
   Struktur-Neuerung in diesem Repo (`AGENTS.md` §3.7 Bestandsgrenze,
   analog). *Bestand bleibt bewusst stehen.*
2. **Keine Änderung an `slice.template.md` als Repo-Artefakt.** Dieses Repo
   führt keine eigene Kopie der Slice-Vorlage (`AGENTS.md` §5: „Baseline
   v5.5.0, template-forward, kein Retrofit") — die Feld-Form kommt direkt
   aus der vendorten `v6.9.0`-Vorlage. *Es wäre ein anderer Vorgang*, eine
   lokale Kopie einzuführen, die es bisher nicht gibt.
3. **Keine weiteren `v6.9.0`-Delta-Punkte.** [`MR-072`](../../../../harness/conventions.md#mr-072)
   nennt einen zweiten offenen Punkt (Review-Report-Tabellenformat) — der
   ist unabhängig von diesem und *ein Folge-Slice übernähme ihn*, keiner,
   den dieser Slice mitzieht.
4. **Die Auflösung der genannten `Gegenstand:`-Kennung bleibt Urteil.** Der
   CR nennt das selbst ausdrücklich als nicht gebeten — ob ein `übernommen
   von slice-<Kennung>` genannter Nehmer existiert oder den Gegenstand
   wirklich führt, prüft kein Sensor dieses Slice (Baseline-Regelwerk
   `modul-05-planning-harness.md` §Was Maschine hier kann: dieselbe Grenze
   gilt bereits für die Register- und Folge-Slice-Paarungen).
5. **Ob ein abgehakter Punkt ein Liefer-Punkt ist, bleibt Urteil.** Ebenfalls
   ausdrücklich vom CR ausgenommen. Die Abgrenzung, **welche** offenen
   Punkte die neue Bedingung überhaupt zählt (alle DoD-Punkte, oder nur die
   unter einer konfigurierten Überschrift), macht dieser Slice
   **konfigurierbar** (Bitte des CR) — er entscheidet nicht inhaltlich,
   welche Zeile in einem fremden Slice-Template ein Liefer-Punkt ist.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** `structure`-Modul um eine **bedingte Pflicht-Zeile** erweitert
      ([`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
      neue Bedingung, begleitende ADR nach dem etablierten Muster dieses
      Moduls): Trägt der geprüfte Abschnitt offene Task-Items (Menge
      konfigurierbar, Default alle — CR-Bitte, Abgrenzung 5), lässt eine
      vorhandene `**Gegenstand:**`-Zeile sie zu (Erlaubnis, unverändert
      gegenüber dem Ursprungs-Plan); fehlt die Zeile dabei, meldet ein
      **eigener, benannter Grund-Code** (nicht das generische
      `section-tasks-open`) genau diese Lücke (Pflicht, CR-Bitte). Ohne
      die neue Bedingung byte-identisches Verhalten (Kalibrierungs-Disziplin
      dieses Moduls). `.d-check.closure.yml`s DoD-Regel nutzt sie für
      **diesen** Slice; die Boilerplate-Haken (`make gates`, Review,
      Closure-Notiz, Register, Risiken, Paarungen) bleiben **weiterhin**
      Pflicht.
- [ ] **(2)** [`CO-002`](../../carveouts/CO-002-slice-221-gegenstand-entfallen.md)
      ist aufgelöst: `make verify-closure-notes` läuft gegen
      [slice-221](../done/slice-221-agents-md-tabellenzellen.md) grün ohne
      dessen `exempt-paths`-Eintrag; der Eintrag ist entfernt, der Carveout
      liegt in `docs/plan/carveouts/done/`. <!-- d-check:ignore (done/ entsteht erst bei erster Carveout-Auflösung) -->
- [ ] **(3)** **Ein Bruch-Test bestätigt alle vier Zustände**: `Gegenstand:`
      + offene Punkte → grün (Erlaubnis) · **ohne** `Gegenstand:` + offene
      Punkte → der neue eigene Befund, **nicht** `section-tasks-open`
      (Pflicht) · ohne `Gegenstand:` + alle Punkte gesetzt → grün
      (Normalfall unverändert) · mit `Gegenstand:` + alle Punkte gesetzt →
      grün (Randfall). Die eingehende CR-Datei bekommt ihre Antwort: `Stand:`
      von „offen" auf „beantwortet, `slice-225`" mit Verweis auf den
      Bruch-Test als Beleg — Kennung ohne Link, da die Zeile erst nach dem
      `git mv` nach `done/` geschrieben wird und der Zielort bis dahin nicht
      existiert.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

1. `internal/hexagon/core/rules/structure.go` und
   `internal/hexagon/core/model/config.go` lesen (bereits als Teil dieses
   Plans gesichtet): `StructureRule` trägt bereits `MaxOpenTasks`,
   `RequirePattern`/`RequireAll` (unbedingt) und `TasksIgnorePattern`/
   `ExemptSectionPattern` (Mengen-Abgrenzung) — kein bestehendes Feld
   koppelt eine Bedingung an eine andere. Die neue Bedingung ist eine
   **Kopplung**, kein neues Modul; Namens- und Schnitt-Entscheidung
   (eigener Schlüssel vs. Erweiterung von `RequirePattern` um eine
   Bedingung) fällt hier, **vor** dem Code, mit derselben Sorgfalt wie die
   übrigen `structure`-Erweiterungen (jede trägt eine eigene ADR).
2. Konfigurierbare Mengen-Abgrenzung entscheiden (Abgrenzung 5): Default
   „alle offenen Task-Items des Abschnitts" (deckt diesen Repos eigenen
   Bedarf ohne neuen Schlüssel), opt-in-Verengung nach demselben Muster wie
   `TasksIgnorePattern`, falls ein Adopter — wie der CR es tut — nur eine
   benannte Unter-Überschrift zählen will.
3. Neuen Grund-Code für die Pflicht-Hälfte festlegen (nicht
   `section-tasks-open` wiederverwenden — der CR bittet ausdrücklich um
   einen **eigenen** Befund, unterscheidbar von „nur vergessen abzuhaken").
4. `spec/lastenheft.md`s [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) um die neue Bedingung
   erweitern (Akzeptanzkriterium, Versions-Bump + Historie-Zeile, wie jede
   frühere Erweiterung dieses Moduls) und begleitende ADR schreiben.
5. `.d-check.closure.yml` anpassen, Bruch-Test-Fixtures für alle vier
   Zustände aus DoD (3) ergänzen.
6. `exempt-paths`-Eintrag für `slice-221` aus der DoD-Regel entfernen,
   `CO-002` auflösen (`git mv` nach `done/`).
7. Antwort in den eingehenden CR eintragen (DoD (3)).
8. `make gates`, Handoff.

## 4. Trigger

**Beanspruchung:** WIP-Limit frei.

**Plan-Änderung überholt diesen Trigger.** Der Ursprungs-Plan sah die
Rückführung vor, falls Schritt 1 eine Go-Code-Änderung samt Spec-Bezug
nötig macht — genau das ist mit der CR-Aufnahme jetzt der **erwartete**
Weg, nicht mehr ein Anzeichen für Fehlschnitt (§1 Plan-Änderung). Der
Trigger bleibt aus einem engeren Grund bestehen:

**Rückführung nach `next/`** (`in-progress→next`): wenn Schritt 1 zeigt,
dass die Kopplung zweier `structure`-Bedingungen (eine schaltet die andere
scharf) eine Erweiterung des **Bedingungs-Modells** selbst braucht, die
über eine neue Bedingung hinausgeht (z. B. eine allgemeine
Abhängigkeits-Syntax zwischen beliebigen Bedingungen) — dann ist das ein
eigener, größerer Vorgang, kein einzelnes Feld nach dem etablierten
Muster dieses Moduls.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `make gates` grün mit echter Ausgabe, unabhängiger
Review durchgeführt und eingearbeitet, Closure-Notiz geschrieben, Register
fortgeschrieben, jedes Risiko aus §6 mit einem der drei Ausgänge.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Erkennung könnte zu weit greifen** — ein `Gegenstand:`-Feld, das
  aus anderem Anlass in einem Fließtext auftaucht (nicht als Closure-Feld
  gemeint), könnte die Prüfung fälschlich entschärfen. Die Form muss eng
  genug sein (z. B. nur als Zeilenanfang direkt unter `## 7.
  Closure-Notiz`), um das auszuschließen. — **Ausgang:** \<offen\>
- **Zwei parallele Wächter-Sprachen** — heute `max-open-tasks: 0` als
  Zahl, künftig zusätzlich eine Bedingung. Ob sich das sauber in die
  bestehende `structure`-Modul-Konfiguration einfügt oder eine neue
  Regel-Klasse braucht, ist vor Schritt 1 nicht abschließend geklärt. —
  **Ausgang:** \<offen\>
- **Die CR-Aufnahme könnte den Zuschnitt sprengen** — `modul-05-planning-harness.md`
  §Ziel-Form: Slice zählt Liefer-Punkte, nicht berührte Artefakte; DoD (1)
  bündelt jetzt Go-Code, Config, Spec-Erweiterung und eine begleitende ADR
  in einem Punkt. Trägt der Bruch-Test in DoD (3) am Ende **eine**
  Review-Sitzung, war das Bündeln richtig; sonst zieht §4s Trigger. —
  **Ausgang:** \<offen\>
- **Die CR-eigene DC-Zuordnung ist falsch übernehmbar** — der CR nennt
  [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
  der Kanon dieses Repos trägt die Bedingung in
  [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  (Bezug oben). Wer nur den Titel des CR zitiert statt
  das Feld nachzuschlagen, überträgt die falsche Kennung in Commit oder
  Spec — genau die Klasse, die
  [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  benennt. — **Ausgang:** \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende. Dieses Repo führt **drei** Prüfungen — die
zwei kanonischen und, als Adaption, den Nachtlauf-Stand
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Die drei Vorprüfungen sind bei der Beanspruchung am 2026-09-17 gelesen.**

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Zwei** Sub-Areas: `*` (Repo-Default, die Config-Hälfte in
`.d-check.closure.yml`) und `internal/hexagon/core/rules/` (die neue
`structure`-Bedingung in Go) — letztere trägt keine eigene
Modus-Deklaration in `harness/conventions.md` und fällt damit unter den
Default `*`. Beide Greenfield, wie der Rest des Produkts.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **39** Verzeichnisse). **Ein
Eintrag ist einschlägig:**

- [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  (19×, verkörpert als Hard Rule seit slice-147) — genau die Klasse, die
  §6 als Risiko gegen die CR-eigene [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)-Zuordnung führt: das
  Feld (`Geltungsbereich`/`Ersetzt-Baseline-Regel`, hier die konkrete
  Modul-Zuordnung), nicht der Titel des CR, entscheidet.

**Keine** der übrigen Einträge trifft `structure`, `spec/lastenheft.md`
oder die CR-Ablage `docs/plan/cr/` als Sub-Area.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-17 gelesen: `image-scan.yml` grün
(2026-09-16T08:37:54Z). `upstream-drift.yml` **rot** (2026-09-16T05:33:02Z)
— laut eigener Meldung eine **planmäßige** Fremd-Release-Benachrichtigung
([`MR-051`](../../../../harness/conventions.md#mr-051)), keine
unerwartete; betrifft die Baseline-Currency, nicht diesen Slice
(`internal/hexagon/core/rules/`, `spec/`, `.d-check.closure.yml`).

**Modus-Begründung:** alle berührten Sub-Areas GF — kein Begründungsblock
nötig.
