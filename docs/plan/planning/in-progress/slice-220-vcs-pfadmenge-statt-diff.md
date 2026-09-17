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

**Verantwortlich:** pt9912 (Implementer-Rolle).

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

- [x] **(1)** Die Kandidaten-Menge kommt aus `vcs.paths`, gegen **beide**
      Stände aufgelöst. Ein Pfad, der in einem Stand nicht lesbar ist, ist ein
      **Umgebungsfehler** (Exit 2) — nicht „nicht vorhanden". Getragen von
      `driven.VCS.AllPaths(base, head)` + `walkTree` (fail-closed über
      `repo.TreeObject`) + `protectedSet`/`ignored()` im Kern.
- [x] **(2)** **Alle drei** Ausprägungen aus der Tabelle in §1 brechen mit
      Exit 2 ab, gemessen gegen dasselbe Probe-Repo; und die **vierte**, heute
      fehldiagnostizierte (unsichtbarer HEAD-Tree ⇒ `core-drift-vcs`
      *„gelöscht oder umbenannt"*, Exit 1) meldet danach den Umgebungsfehler
      statt eines Inhalts-Befunds. Getragen von
      `TestVCS_UnlesbareObjekte` (vier Unterfälle, End-to-End gegen ein
      echtes Repo mit real entferntem losen Objekt) und **empirisch
      unabhängig** vom Review R1 (eigenes Probe-Repo, Parent- vs.
      Feature-Commit-Binary) reproduziert.
- [x] **(3)** **Die Gegenrichtung ist mitgemessen und als Testmenge benannt** —
      nicht als Einzelprobe: eine wirklich neu angelegte Datei, eine wirklich
      gelöschte, ein bewegter Gitlink, ein reiner Rename und eine unveränderte
      Datei bleiben, wie sie heute sind. Für jeden ein Test, der ohne den
      Umbau grün ist und mit ihm grün bleibt. **Präzisierung nach R1s
      Negativbefund:** der Gitlink-Test deckt den **hinzugefügten** Gitlink;
      ein geänderter/bewegter wird nicht separat getestet, weil derselbe
      unbedingte `switch`-Zweig in `walkTree` jeden `filemode.Submodule`-
      Eintrag identisch behandelt, unabhängig davon, ob er neu, geändert oder
      unverändert ist (R1/R2 haben den Codepfad geprüft, kein Finding).
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (zwei Runden: R1 blockiert mit 1 MEDIUM/2 LOW, eingearbeitet; R2
      freigegeben, 0 Befunde).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

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
  gefahren**. — **Ausgang:** entfallen. Die Matrix (DoD 2/3) wurde vollständig
  gefahren und **zweifach unabhängig verifiziert** — Review R1 hat eigene
  Probe-Repos gebaut und Parent- gegen Feature-Commit-Binary gegenübergestellt
  (empirisch, nicht nur Code-Lesung), R2 hat zusätzlich eine Codepfad-Analyse
  über alle 18 Prüffragen gefahren. Beide Runden fanden keine fünfte, noch
  offene Form.
- **Eine umgekehrte Kandidaten-Bestimmung kann Fälle verlieren, die der Diff
  heute liefert.** Ein Pfad, der in **keinem** der beiden Stände existiert,
  aber im Diff auftaucht (etwa über eine Zwischenrevision), fiele aus der
  Menge. Ob es solche Fälle gibt, ist **nicht gemessen** — DoD (3) fragt
  danach. — **Ausgang:** entfallen. Das Szenario ist mit einem Zwei-Baum-
  Vergleich strukturell unmöglich, in der alten Fassung genauso wie in der
  neuen: `DiffTreeWithOptions(base, head)` (alt) wie die Mengen-Differenz aus
  `AllPaths(base, head)` (neu) können nur Pfade nennen, die in **mindestens
  einem** der beiden übergebenen Tree-Stände tatsächlich existieren — eine
  „Zwischenrevision" wird bei keinem der beiden Ansätze je gelesen. Die Sorge
  bezog sich auf einen Fall, den auch der abgelöste Mechanismus nie hatte.
- **`vcs.paths` ist ein Glob, kein Verzeichnis.** Die Menge über Globs
  aufzulösen heißt, **beide** Trees vollständig zu durchlaufen — genau das,
  was der Diff heute vermeidet. Der Preis steht in Abgrenzung 4 und ist damit
  benannt, nicht gelöst. — **Ausgang:** entfallen (aus dem Slice
  ausgekoppelt). Abgrenzung 4 (§1) benennt den Preis bereits als bewusst
  nicht angegangene Größe; kein Performance-Vorfall ist gemessen, kein
  Folge-Slice nötig. Wird die Laufzeit auf großen Repos je zum Problem, ist
  das ein neuer, eigener Anlass — kein Rest dieses Slice.

## 7. Closure-Notiz

**Geliefert.** Das Modul `vcs`
([`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)/[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs))
bestimmt seine Kandidaten jetzt über `driven.VCS.AllPaths(base, head)` —
jeder Datei-Pfad wird direkt aus dem jeweiligen Tree-Stand gelesen
(fail-closed über einen selbstgeschriebenen `walkTree`, statt über den
geteilten `Tree.Files()`-Walker, der einen nicht ladbaren Unterbaum
stillschweigend abschneidet), der Kern bildet die Mengen-Differenz selbst.
`ChangedPaths`/`VCSChange`/`VCSStatus`/`diffTrees`/`diffTreeIndex` sind
ersatzlos entfallen. Das schließt
[`CO-001`](../../carveouts/CO-001-vcs-range-stiller-skip.md)s dritte
Ausprägung (BASE-Tree ohne Pendant gelöscht — bisher `0 Befund(e)`/Exit 0)
und deckt eine vierte, bisher fehldiagnostizierte auf (unlesbarer HEAD-Tree
meldete fälschlich `core-drift-vcs`, Exit 1, statt eines Umgebungsfehlers).
Beide fallen jetzt auf denselben Codepfad wie die beiden bereits in
slice-218 behobenen Ausprägungen — der Ableiter aus
[`fix-schliesst-pfad-nicht-klasse`](../observations/BEO-ALL/fix-schliesst-pfad-nicht-klasse/observation.md)
(„erst die Schicht benennen, auf der ein Fix wirkt, dann prüfen, was davor
liegt") trägt.

**Review-Runde 1** ([`docs/reviews/2026-09-17-slice-220-vcs-pfadmenge-review-r1.md`](../../../reviews/2026-09-17-slice-220-vcs-pfadmenge-review-r1.md)):
0 HIGH · 1 MEDIUM · 2 LOW, Verdikt „Blockiert". **R1-F-1 (MEDIUM):**
`spec/spezifikation.md` §[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 2 — genau die Spec-Stelle,
die der Slice-Kopf selbst unter „Berührte Spec-Stellen" als geändert nennt —
beschrieb nach dem Feature-Commit weiterhin den abgelösten Diff-Mechanismus
(Diff-Status-Klassen `M`/`T`/`D`/`R`/`A`, die (ab-/an-)geschaltete
Rename-Erkennung, den Sonderzweig „`A` liest BASE nach"). Umgeschrieben auf
die tatsächliche Mengen-Differenz, mit neuem Historie-Eintrag. **R1-F-2/F-3
(LOW):** zwei Kommentare (der `vcsDeleted`-Doc-Kommentar, ein
Test-Doc-Kommentar) nannten noch die entfernte „Diff-Übersetzung im
VCS-Adapter" bzw. `ChangedPaths` — nachgezogen.

**Review-Runde 2** ([`docs/reviews/2026-09-17-slice-220-vcs-pfadmenge-review-r2.md`](../../../reviews/2026-09-17-slice-220-vcs-pfadmenge-review-r2.md)),
Satz-für-Satz-Vergleich der neuen Schritt-2-Fassung gegen den tatsächlichen
Code plus vollständiger Neu-Durchlauf aller 18 Prüffragen über den
gesamten Beanspruchungs-Diff: **0 HIGH · 0 MEDIUM · 0 LOW · 0 INFO**,
Verdikt „Freigegeben". R1-F-1 bis F-3 bestätigt behoben, keine neuen
Stale-Referenzen, keine neuen Findings.

**Beobachtungs-Register fortgeschrieben:**
[`fix-schliesst-pfad-nicht-klasse`](../observations/BEO-ALL/fix-schliesst-pfad-nicht-klasse/observation.md)
bleibt bei **1×**/`offen` (der Zähler misst Wiederholung über Vorgänge
hinweg, nicht die Bestätigung einer Diagnose) — `state.md` trägt jetzt, dass
der empfohlene Schichtwechsel umgesetzt und zweifach unabhängig verifiziert
wurde, statt nur „trägt den Ableiter bereits" zu sagen.

**Ein Risiko wurde beim Schreiben unterschätzt und beim Review korrigiert
festgehalten, nicht verschwiegen:** Die Sorge um eine „Zwischenrevision", die
der alten Diff-Fassung Kandidaten liefern könnte, die keinem der beiden
Enden angehören, erwies sich bei genauerer Prüfung als **strukturell
unmöglich** für einen Zwei-Baum-Vergleich — weder der alte noch der neue
Mechanismus liest je eine dritte Revision. Siehe §6 für die vollständige
Begründung.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende.

Dieses Repo führt **drei** Prüfungen — die zwei kanonischen und, als
Adaption, den Nachtlauf-Stand ([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Die drei Vorprüfungen sind bei der Beanspruchung am 2026-09-17 gelesen.**

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Zwei** Sub-Areas: `*` (Repo-Default, Konfig-Schema in `internal/hexagon/core/model/`
und Reason-Codes) und `internal/adapter/driven/git/` (der Tree-Walker,
den dieser Slice neu baut) — letztere trägt keine eigene Modus-Deklaration in
`harness/conventions.md` und fällt damit unter den Default `*`. Beide
Greenfield, wie der Rest des Produkts.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, 41 Verzeichnisse). **Zwei Einträge
sind einschlägig:**

- [`fix-schliesst-pfad-nicht-klasse`](../observations/BEO-ALL/fix-schliesst-pfad-nicht-klasse/observation.md)
  (Sub-Area `*`, 1×) — der Anlass dieses Slice selbst; §3 trägt seinen
  Ableiter bereits.
- [`racily-clean-git-fixture`](../observations/BEO-ALL/racily-clean-git-fixture/observation.md)
  (Sub-Area `internal/adapter/driven/git`, unter der Schwelle) — jedes neue
  git-Fixture in diesem Slice schreibt Dateien über den bestehenden
  `put()`-Helfer, der `coretest.GitFixtureRewriteHazard` bereits trägt; kein
  neuer Fundort.

**Keine** der übrigen Einträge trifft `internal/adapter/driven/git/`,
`internal/hexagon/core/rules/` oder `internal/hexagon/port/driven/` als
Sub-Area.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-17 gelesen: `image-scan.yml` grün
(2026-09-16T08:37:54Z). `upstream-drift.yml` **rot** (2026-09-17T05:38:51Z)
— `make freshness-semgrep` und `make go-base-digest` melden neuere
Fremd-Releases, laut eigener Meldung eine **planmäßige** Benachrichtigung
([`MR-051`](../../../../harness/conventions.md#mr-051)), keine unerwartete;
betrifft gepinnte Fremd-Bestände, nicht diesen Slice
(`internal/adapter/driven/git/`, `internal/hexagon/core/rules/`).

**Modus-Begründung:** alle berührten Sub-Areas GF — kein Begründungsblock
nötig.
