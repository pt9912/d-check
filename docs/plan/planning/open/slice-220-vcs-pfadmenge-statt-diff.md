# Slice slice-220: `vcs` löst die geschützte Pfad-Menge auf, statt dem Diff zu vertrauen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
— und [`CO-001`](../../carveouts/CO-001-vcs-range-stiller-skip.md), dessen
offene dritte Ausprägung dieser Slice schließt.

**Berührte Spec-Stellen:** [`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 2 (die Kandidaten-Menge —
der Slice ändert, **woher** sie kommt).

**Verantwortlich:** — (wird bei der Beanspruchung gesetzt).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Ziel.** Das Modul `vcs` bestimmt seine Kandidaten **nicht mehr allein aus dem
Tree-Diff**, sondern löst die über `vcs.paths` geschützte Menge **direkt gegen
beide Trees** auf. Damit ist eine Änderung, die der Diff gar nicht liefert,
kein blinder Fleck mehr.

**Warum ein Entwurf und kein vierter Patch — das ist der Kern dieses Slice.**
In slice-218 wurde derselbe Defekt **dreimal** hintereinander gefunden, jedes
Mal in einer Form, die der vorige Fix nicht deckte, und jedes Mal vom
unabhängigen Review statt vom Autor:

| Form | wo der Fix ansetzte | Ergebnis |
|---|---|---|
| BASE-**Blob** unsichtbar | Adapter, `FileAt` | zu |
| BASE-**Tree** unsichtbar, mit Pendant | Regel, `A`-Zweig | zu |
| BASE-**Tree** unsichtbar, **ohne** Pendant | — | **offen** |

Die dritte entsteht **vor** jeder Stelle, an der das Modul prüfen kann: Der
Tree-Walker der Bibliothek macht aus einem nicht ladbaren Unterbaum ein
`io.EOF`. **Gemessen ist auch, dass der naheliegende Wachposten nicht trägt** —
`tree.Files()` benutzt denselben Walker und schweigt ebenso. Ein vierter Patch
auf derselben Schicht wäre also nicht nur wahrscheinlich unvollständig, er ist
an der entscheidenden Stelle **unmöglich**.

**Der Ansatz kehrt die Richtung um.** Statt zu fragen *„was hat sich laut Diff
geändert, und ist etwas davon geschützt?"* fragt das Modul *„was ist geschützt,
und wie steht es in beiden Ständen?"*. Die geschützte Menge ist konfiguriert
(`vcs.paths`) und damit bekannt, bevor irgendein Objekt gelesen wird.

**Abgrenzung — vier Punkte, jeder mit Grund:**

1. **Kein Repack durch das Produkt.** Wie in slice-218: das schriebe in das
   geprüfte Repository und bräche
   [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
   — *es wäre ein anderer Vorgang*.
2. **Kein anderes Modul.** `commits` liest Commit-Objekte und ist gemessen
   fail-closed; `tracked` liest den Index und ist gar nicht betroffen. Ob ein
   **anderes** Modul dieselbe Diff-Abhängigkeit trägt, ist **nicht gemessen**
   und wird hier nicht gemessen — *es wäre ein anderer Vorgang*.
3. **Keine Änderung an der Befund-Semantik.** Grund-Codes und Meldungen
   bleiben; es ändert sich, **welche** Kandidaten geprüft werden, nicht wie
   geurteilt wird — *Schicht-Abgrenzung*.
4. **Kein Performance-Ziel.** Die Auflösung der Pfad-Menge kann teurer sein als
   der Diff. Ob das messbar stört, ist eine eigene Frage; dieser Slice sichert
   die **Korrektheit** und misst die Laufzeit nur, um sie zu benennen —
   *Bestand bleibt bewusst stehen*.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Die Kandidaten-Menge kommt aus `vcs.paths`, gegen **beide**
      Stände aufgelöst. Ein Pfad, der in einem Stand nicht lesbar ist, ist ein
      **Umgebungsfehler** (Exit 2) — nicht „nicht vorhanden".
- [ ] **(2)** **Alle drei** Ausprägungen aus der Tabelle in §1 brechen mit
      Exit 2 ab, gemessen gegen dasselbe Probe-Repo; und die **vierte**, heute
      fehldiagnostizierte (unsichtbarer HEAD-Tree ⇒ `core-drift-vcs`
      *„gelöscht oder umbenannt"*, Exit 1) meldet danach den Umgebungsfehler
      statt eines Inhalts-Befunds.
- [ ] **(3)** **Die Gegenrichtung ist mitgemessen und als Testmenge benannt** —
      nicht als Einzelprobe: eine wirklich neu angelegte Datei, eine wirklich
      gelöschte, ein bewegter Gitlink, ein reiner Rename und eine unveränderte
      Datei bleiben, wie sie heute sind. Für jeden ein Test, der ohne den
      Umbau grün ist und mit ihm grün bleibt.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Die Menge steht vor dem Umbau fest.** Bevor eine Zeile geändert wird, wird
die Prüf-Matrix aufgeschrieben: **vier** Fehler-Achsen (BASE-Blob, BASE-Tree
mit/ohne Pendant, HEAD-Tree) mal **fünf** Gegenrichtungs-Fälle aus DoD (3).
Der Ableiter von
[`fix-schliesst-pfad-nicht-klasse`](../observations/BEO-ALL/fix-schliesst-pfad-nicht-klasse/observation.md)
gilt hier wörtlich: **erst die Schicht benennen, auf der ein Fix wirkt, dann
prüfen, was davor liegt.** Bei diesem Umbau ist die Schicht die
Kandidaten-Bestimmung, und davor liegt nur noch das Auflösen der Refs.

**Schritte:**

1. Prüf-Matrix schreiben, Probe-Repos bauen, **alle** Zellen gegen den
   heutigen Stand messen — die Vorher-Spalte ist der Beleg.
2. Kandidaten-Bestimmung umbauen.
3. Matrix erneut fahren; jede Zelle muss sich erklären lassen.
4. `make gates`, Handoff.

## 4. Trigger

**Beanspruchung:** slice-218 ist geschlossen, und der Auftraggeber gibt den
Umbau frei — er ändert das Verhalten eines Gates in der CI.

**Rückführung nach `next/`** (`in-progress→next`): wenn die Prüf-Matrix aus
Schritt 1 zeigt, dass mehr als die vier benannten Achsen betroffen sind —
dann ist der Zuschnitt zu klein und die Menge muss vorher geschnitten werden.

**Rückführung nach `open/`** (`in-progress→open`): wenn der Umbau die
Laufzeit auf großen Repos unbrauchbar macht (Abgrenzung 4 nennt das als
gemessene, nicht als optimierte Größe) — dann braucht es vorher eine
Entscheidung über den Preis.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `make gates` grün mit echter Ausgabe, unabhängiger
Review durchgeführt und eingearbeitet, Closure-Notiz geschrieben, Register
fortgeschrieben, jedes Risiko aus §6 mit einem der drei Ausgänge.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Umbau ist die vierte Antwort auf denselben Defekt — und die dritte
  war auch überzeugend.** Jeder der drei Vorgänger-Fixe sah beim Schreiben
  vollständig aus. Dass dieser die Schicht wechselt statt sie zu patchen, ist
  ein Argument, kein Beweis; der Beleg ist allein die Matrix aus §3, **vollständig
  gefahren**. — **Ausgang:** \<offen\>
- **Eine umgekehrte Kandidaten-Bestimmung kann Fälle verlieren, die der Diff
  heute liefert.** Ein Pfad, der in **keinem** der beiden Stände existiert,
  aber im Diff auftaucht (etwa über eine Zwischenrevision), fiele aus der
  Menge. Ob es solche Fälle gibt, ist **nicht gemessen** — DoD (3) fragt
  danach. — **Ausgang:** \<offen\>
- **`vcs.paths` ist ein Glob, kein Verzeichnis.** Die Menge über Globs
  aufzulösen heißt, **beide** Trees vollständig zu durchlaufen — genau das,
  was der Diff heute vermeidet. Der Preis steht in Abgrenzung 4 und ist damit
  benannt, nicht gelöst. — **Ausgang:** \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende.

**Die drei Vorprüfungen entstehen spätestens bei der Beanspruchung**
([`AGENTS.md`](../../../../AGENTS.md) §5) — dieser Plan liegt in `open/` und
trägt sie noch nicht. **Einer ist trotzdem schon absehbar:**
[`fix-schliesst-pfad-nicht-klasse`](../observations/BEO-ALL/fix-schliesst-pfad-nicht-klasse/observation.md)
steht nach slice-218 bei **1×** (drei Instanzen, aber ein Vorgang) und ist der Anlass dieses Slice; §3
trägt seinen Ableiter bereits.
