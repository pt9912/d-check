# Slice slice-221: `AGENTS.md` §4 wird wieder ein Index — die Substanz zieht in die Sensor-Dateien

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`DC-FA-TGT-001`](../../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(das Modul `targets` liest §4 als Autoritäts-Tabelle — die Target-**Namen**
müssen alle bleiben) und
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(der Wächter, der den Zustand hält). Keine ADR: Der Slice trifft keine
Entscheidung, er stellt eine bereits geltende Regel wieder her.

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle).

**Verantwortlich:** — (wird bei der Beanspruchung gesetzt).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Ziel.** Die Gates-Tabelle in [`AGENTS.md`](../../../../AGENTS.md) §4 trägt
wieder **eine Zeile Zusage je Target** statt eines Absatzes; die Substanz zieht
dorthin, wo sie ohnehin hingehört — `harness/sensors/<target>.md`. Ein neuer
`structure`-Wächter hält den Zustand, damit er nicht zurückwächst.

**Der Anlass ist eine Auftraggeber-Beobachtung, und sie ist gemessen:**

| | `AGENTS.md` | `harness/README.md` |
|---|---|---|
| Größe | **69,6 KB**, 660 Zeilen | 23,7 KB, 233 Zeilen |
| davon Tabellen | **42 %** | 67 % |
| längste Zelle | **4090 Zeichen** | 817 Zeichen |
| §4-Tabelle | 43 Zeilen, **27 607 Zeichen**, Median 374 | — |
| Zellen > 1000 Zeichen | **8** | 0 |
| Zellen ≤ 200 Zeichen | **2** von 43 | — |

**Drei Gründe, warum das mehr kostet als Platz.** Erstens wird `AGENTS.md`
über `CLAUDE.md`s `@AGENTS.md` in **jeden** Lauf importiert — die 69,6 KB sind
Kontext-Grundlast, jedes Mal. Zweitens **verstößt es gegen die Regel, die die
Datei selbst in §1 aufstellt**: *„Diese Datei trägt Hard Rules und Pointer auf
die kanonischen Quellen und dupliziert deren Inhalt nicht — sonst entsteht
Drift."* Drittens ist die Drift real: **alle 54 `make`-Targets stehen in
beiden Dateien**, null nur in der einen, null nur in der anderen.

**Die Dopplung der Tabellen ist trotzdem kein Fehler, und das grenzt den
Slice ein.** Die beiden haben verschiedene deklarierte Rollen: §4 ist die
**Autoritäts-Liste**, gegen die `make gate-consistency` jede Makefile-Regel
prüft ([`DC-FA-TGT-001`](../../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in));
die Sensors-Tabelle trägt die **Bindung**-Spalte und eine eigene
`structure`-Regel. Gewuchert sind nicht die Tabellen, sondern die
**Beschreibungen in ihren Zellen**.

**Wohin die Substanz zieht, ist gemessen und nicht geraten:** **24** der 54
Targets haben bereits eine `harness/sensors/<target>.md`, und **13 der 14
größten Zellen** gehören dazu. Die eine Ausnahme unter den Großen ist
`archive-wave` (1894 Zeichen) — sie bekommt eine eigene Datei, wie
`guard-probe` und `hooks` sie als Werkzeuge längst haben.

**Abgrenzung — vier Punkte, jeder mit Grund:**

1. **§3 und §5 bleiben unberührt.** Sie sind mit 328 und 169 Zeilen die
   eigentliche Masse, aber sie tragen **Hard Rules**, keine Tabellen-Zellen.
   Sie zu kürzen ist eine andere Risiko-Klasse — jede Zeile dort ist eine
   Regel, deren Wegfall niemand meldet. *Es wäre ein anderer Vorgang*, und
   zwar einer, der jede einzelne Regel einzeln verantworten muss.
2. **Kein Inhalt wird gelöscht, nur verschoben.** Was heute in einer Zelle
   steht, steht danach in der Sensor-Datei — vollständig. Wer kürzen will,
   was inhaltlich überflüssig ist, führt ein Urteil je Satz; das ist
   *ein anderer Vorgang* — *Schicht-Abgrenzung*.
3. **`harness/README.md` wird nicht angefasst.** Seine Sensors-Tabelle hat
   eine eigene `structure`-Regel mit `cell-min-chars` und eine andere Rolle.
   Ob **sie** zu lang ist, ist **nicht gemessen** und wird hier nicht
   gemessen — *ein Folge-Slice übernähme es*.
4. **Keine Änderung an `gate-consistency` oder der Autoritäts-Rolle.** Alle
   54 Target-**Namen** bleiben in §4; nur ihre Beschreibungen schrumpfen.
   *Bestand bleibt bewusst stehen.*

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Jede Zelle der §4-Tabelle trägt **eine Zusage** und, wo es eine
      gibt, den Zeiger auf `harness/sensors/<target>.md`. Ein neuer
      `structure`-Eintrag mit `cell-max-chars` hält den Zustand — Präzedenz im
      Bestand: dieselbe Bedingung steht bereits dreimal in
      [`.d-check.yml`](../../../../.d-check.yml). **Die Schwelle wird aus dem
      Ergebnis abgeleitet, nicht vorher gesetzt**, und ihre Wahl steht im
      Konfigurations-Kommentar.
- [ ] **(2)** **Die Substanz ist vollständig erhalten**, in
      `harness/sensors/<target>.md`: 24 vorhandene Dateien ergänzt, **eine
      neue** für `archive-wave`. Für die 30 Targets ohne Sensor-Datei gilt:
      Wo die Zelle kurz genug ist, bleibt sie; wo nicht, entsteht eine Datei —
      welche das sind, sagt die Inventur aus DoD (3), nicht eine Schätzung.
- [ ] **(3)** **Nichts ist verlorengegangen, und das ist gemessen** — eine
      Vorher/Nachher-Inventur je Target über die Zeichenmenge, plus
      `make gate-consistency` grün (alle 54 Namen weiterhin in der
      Autoritäts-Doku) und `make doc-check` grün.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Die Form des Gegenstands, bevor gezählt wird** (Ableiter von
[`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)):
Gemessen wird die **Beschreibungs-Zelle** einer Zeile der §4-Tabelle — das
zweite Feld zwischen zwei `|`, ohne die Target-Zelle davor und ohne die
Trennzeichen. Nicht gemessen: Zeilen, die keine Target-Zeile sind (Kopf,
Trenner), und die Target-Zelle selbst.

**Schritte:**

1. **Inventur schreiben**: je Target die heutige Zellen-Zeichenmenge, ob eine
   Sensor-Datei existiert, und wohin die Substanz zieht. Die Inventur ist der
   Beleg für DoD (3) und entsteht **vor** der ersten Änderung.
2. Substanz je Target verschieben — **eine Datei nach der anderen**, nicht
   mechanisch über den Baum
   ([`MR-070`](../../../../harness/conventions.md#mr-070) gilt der mechanischen
   Ersetzung; hier ist jede Zelle ein eigenes Urteil darüber, was Zusage und
   was Detail ist).
3. Schwelle aus dem Ergebnis ableiten, `structure`-Regel einziehen, gegen den
   neuen Stand fahren.
4. `make gates`, Handoff.

## 4. Trigger

**Beanspruchung:** WIP-Limit frei und der Auftraggeber hat den Zuschnitt am
2026-09-08 beauftragt, nachdem er die Größe selbst beobachtet hatte.

**Rückführung nach `next/`** (`in-progress→next`): wenn die Inventur aus
Schritt 1 zeigt, dass mehr als etwa ein Dutzend Targets eine **neue**
Sensor-Datei brauchen — dann ist der Slice keine Umschichtung mehr, sondern
das Anlegen einer Doku-Ebene, und das ist ein anderer Zuschnitt.

**Rückführung nach `open/`** (`in-progress→open`): wenn sich beim Verschieben
zeigt, dass eine Zelle Substanz trägt, die **nirgendwo sonst** hingehört —
etwa eine Hard-Rule-Aussage, die sich als Gate-Beschreibung tarnt. Dann ist
vor dem Verschieben eine Entscheidung fällig, wo sie hingehört.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `make gates` grün mit echter Ausgabe, unabhängiger
Review durchgeführt und eingearbeitet, Closure-Notiz geschrieben, Register
fortgeschrieben, jedes Risiko aus §6 mit einem der drei Ausgänge.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Eine Zelle kann eine Hard Rule enthalten, die sich als Gate-Beschreibung
  tarnt.** Die längste (`verify-closure-notes`, 4090 Zeichen) beschreibt nicht
  nur, was das Target tut, sondern **welche Regeln es durchsetzt** und wo ihre
  Grenzen liegen. Wer sie schematisch in die Sensor-Datei schiebt, verschiebt
  womöglich eine Regel aus dem Dokument, das **jeder Lauf lädt**, in eines,
  das nur bei Bedarf gelesen wird. **Das ist der teuerste Fehler, den dieser
  Slice machen kann**, und §4 nennt ihn als Rückführungs-Bedingung. —
  **Ausgang:** \<offen\>
- **Der Slice ist groß, und seine Größe ist nicht die Zahl der Liefer-Punkte,
  sondern die Zahl der Urteile:** 43 Zellen, jede einzeln. Die
  Ein-Sitzungs-Review-Grenze
  ([`MR-066`](../../../../harness/conventions.md#mr-066)) ist damit absehbar
  überschritten, und die Ersatz-Form der Prüfung gehört **vorab** in den Plan,
  nicht in den Bericht. **Vorschlag, hier schon benannt:** eine mechanische
  Vollprüfung entlang der Zeichenmenge (Inventur vorher/nachher, kein Zeichen
  darf verschwinden) plus eine Stichprobe von fünf Zellen Wort für Wort.
  Bestätigt oder ersetzt wird sie bei der Beanspruchung. — **Ausgang:** \<offen\>
- **Ein `cell-max-chars` auf §4 ist eine neue Schwelle, und Schwellen altern.**
  Eine zu enge zwingt künftige Autoren, Substanz in die Sensor-Datei zu
  schreiben — das ist der Zweck. Eine zu weite hält nichts. Die Wahl ist ein
  Urteil, das kein Sensor prüft; sie gehört in den
  Konfigurations-Kommentar, samt der Zahl, aus der sie abgeleitet wurde. —
  **Ausgang:** \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende.

**Die drei Vorprüfungen entstehen spätestens bei der Beanspruchung**
([`AGENTS.md`](../../../../AGENTS.md) §5) — dieser Plan liegt in `open/`.
**Zwei sind trotzdem schon gelaufen**, weil sie den Zuschnitt tragen, und
werden bei der Beanspruchung gegen den dann gültigen Stand wiederholt:

- **Nachtlauf-Stand** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
  am 2026-09-08 beide grün (`upstream-drift.yml` 05:31:31Z, `image-scan.yml`
  08:07:20Z). Ohne Bezug zu diesem Slice.
- **Register** (40 Verzeichnisse): **zwei** Einträge sind einschlägig.
  [`registry-vs-authority-table-drift`](../observations/BEO-ALL/registry-vs-authority-table-drift/observation.md)
  (2×) trifft den Gegenstand unmittelbar — ein Register und seine
  Autoritäts-Tabelle driften; hier sind es zwei Tabellen über dieselbe Menge.
  **Erreicht der Eintrag mit diesem Slice 3×, ist er keine Notiz mehr,
  sondern eine Lücke** und braucht einen eigenen Folge-Slice; das ist bei der
  Beanspruchung zu entscheiden.
  [`large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)
  (3×, verkörpert als [`MR-066`](../../../../harness/conventions.md#mr-066))
  ist der Grund, warum §6 die Ersatz-Form der Prüfung schon jetzt benennt.
