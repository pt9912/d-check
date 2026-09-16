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
(dessen Auflösung dieser Slice trägt).

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle; die
Änderung liegt in `.d-check.closure.yml`, einer Repo-Konfiguration, nicht im
Produkt).

**Verantwortlich:** — (wird bei der Beanspruchung gesetzt).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Ziel.** Baseline `v6.9.0`s `modul-05-planning-harness.md` führt einen
vierten Slice-Lifecycle-Zweig: Ein Slice, dessen Gegenstand ein anderer
Slice übernimmt oder der ganz entfällt, geht **ohne Lieferung** nach
`done/` — §7 trägt dann eine Zeile `Gegenstand:` (`übernommen von
slice-<Kennung>` | `entfallen: <Grund>`), und die Liefer-Punkte der DoD
bleiben **leer**. Dieser Slice macht `make verify-closure-notes` (Modul
`structure`, Regel `max-open-tasks: 0`) mit dieser Form kompatibel — heute
verlangt sie ausnahmslos, dass jeder DoD-Haken gesetzt ist, ohne Rücksicht
auf einen `Gegenstand:`-Ausgang.

**Der Anlass ist [slice-221](../next/slice-221-agents-md-tabellenzellen.md)**,
das erste Repo-Beispiel dieses Zwecks (Ausgang „entfallen" — slice-222 hat
seinen Gegenstand durch eine andere Lösung erledigt) und deshalb per
[`CO-002`](../../carveouts/CO-002-slice-221-gegenstand-entfallen.md) statt
regulärer Closure nach `done/` gewandert.

**Abgrenzung — drei Punkte, jeder mit Grund:**

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

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** `.d-check.closure.yml`s `structure`-Regel für die DoD-Sektion
      erkennt eine `**Gegenstand:**`-Zeile im Closure-Notiz-Abschnitt (§7)
      und lässt für **diesen** Slice offene Liefer-Punkte zu — die
      Boilerplate-Haken (`make gates`, Review, Closure-Notiz, Register,
      Risiken, Paarungen) bleiben **weiterhin** Pflicht.
- [ ] **(2)** [`CO-002`](../../carveouts/CO-002-slice-221-gegenstand-entfallen.md)
      ist aufgelöst: `make verify-closure-notes` läuft gegen
      [slice-221](../next/slice-221-agents-md-tabellenzellen.md) grün ohne
      dessen `exempt-paths`-Eintrag; der Eintrag ist entfernt, der Carveout
      liegt in `docs/plan/carveouts/done/`. <!-- d-check:ignore (done/ entsteht erst bei erster Carveout-Auflösung) -->
- [ ] **(3)** **Ein Bruch-Test bestätigt beide Richtungen**: ein `done/`-Slice
      mit `Gegenstand:`-Zeile und offenem Liefer-Haken bleibt grün; ein
      `done/`-Slice **ohne** `Gegenstand:`-Zeile und offenem Liefer-Haken
      bleibt wie bisher rot (`section-tasks-open`).
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

1. Erkennungsform für `**Gegenstand:**` in `internal/hexagon/core/rules/structure.go`
   (oder deklarativ über die YAML-Config, falls die bestehende
   `conditions`-artige Mechanik das trägt — **erst prüfen, dann
   entscheiden**, ob Code- oder Config-Änderung).
2. `.d-check.closure.yml` anpassen, Bruch-Test-Fixtures ergänzen.
3. `exempt-paths`-Eintrag für `slice-221` aus der DoD-Regel entfernen,
   `CO-002` auflösen (`git mv` nach `done/`).
4. `make gates`, Handoff.

## 4. Trigger

**Beanspruchung:** WIP-Limit frei.

**Rückführung nach `next/`** (`in-progress→next`): wenn Schritt 1 zeigt,
dass die bestehende `structure`-Config-Sprache die Bedingung nicht tragen
kann und eine Go-Code-Änderung an einem öffentlichen Modul-Vertrag nötig
wird, die eine Spec-Änderung nach sich zieht — dann ist der Zuschnitt zu
klein für einen Config-only-Slice.

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

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende.

**Die drei Vorprüfungen entstehen spätestens bei der Beanspruchung**
([`AGENTS.md`](../../../../AGENTS.md) §5) — dieser Plan liegt in `open/` und
trägt sie noch nicht.
