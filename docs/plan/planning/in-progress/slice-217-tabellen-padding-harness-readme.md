# Slice slice-217: Tabellen-Padding in `harness/README.md` auf die Vorlagen-Form

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle. Es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden wäre — kein repo-weiter Beleg über die DoD hinaus, kein
zweiter Slice, der mitschließt (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
und [`DC-FA-TGT-001`](../../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
— **nicht als geänderte Anforderungen, sondern als die zwei Zusagen, die auf
die angefasste Datei zeigen** und deshalb nach der Änderung noch gelten
müssen. Keine ADR: Der Slice ändert keine Entscheidung, sondern gleicht eine
Darstellungsform an die bereits adoptierte Vorlage an.

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle).

**Verantwortlich:** pt9912 (Implementer-Rolle).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Ziel.** Die beiden ausgerichteten Tabellen in
[`harness/README.md`](../../../../harness/README.md) — §Source precedence und
§Guides — verlieren ihr Spalten-Padding und tragen danach dieselbe schlanke
Form wie die drei übrigen Tabellen derselben Datei (§Sensors, §Werkzeuge,
§Gate-Taxonomie) **und wie die Baseline-Vorlage**
[`harness/README.template.md`](../../../../.harness/baseline/v6.5.0/templates/harness/README.template.md).

**Der Anlass ist gemessen, nicht ästhetisch.** Von **4737** Padding-Zeichen
stehen **3394 in der jeweils letzten Spalte** — dort, wo hinter dem Pipe
nichts mehr folgt, das ausgerichtet werden könnte. Die restlichen 1297 richten
aus, aber unvollständig: Z. 43, 44, 60 und 62 brechen die Ausrichtung bereits
heute, Z. 60 sogar am **inneren** Pipe. Eine Ausrichtung, die vier von 23
Zeilen nicht einhalten, ist keine mehr; sie kostet nur noch.

**Die Form-Frage ist an der Vorlage geklärt, nicht am Nachbarn.** Die
Baseline-Vorlage führt **39** Tabellenzeilen mit **null** Padding (Trenner
`|---|---|---|`). Das Entpadden bewegt diese Datei also **zur** Vorlage hin.
Ohne diese Prüfung wäre der Slice eine Angleichung an den lokalen Bestand
gewesen — genau der Fehler, den
[`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
führt.

**Abgrenzung — vier Punkte, jeder mit Grund:**

1. **Kein Inhalt ändert sich, kein Zeichen außerhalb der Padding-Klasse.**
   Der Slice ist eine Whitespace-Operation; jede inhaltliche Mitnahme wäre
   eine andere Arbeit und im Diff nicht mehr von ihr zu trennen — *Schicht-
   Abgrenzung*, beim Review sofort prüfbar.
2. **Keine andere Datei.** `AGENTS.md` führt ebenfalls große Tabellen und
   steht in derselben `doc-tables`-Liste; ob sie dasselbe Bild zeigt, ist
   **nicht gemessen** und wird hier auch nicht gemessen. Ein Befund dort wäre
   ein eigener Vorgang — *es wäre ein anderer Vorgang*.
3. **Kein Sensor auf die Tabellen-Form.** Ein Wächter, der Padding meldet,
   wäre eine neue Zusage mit eigener Grenzen-Frage und eigener ADR-Last. Der
   Slice stellt einen Zustand her, er verankert ihn nicht — *Bestand bleibt
   bewusst stehen*, und zwar der Zustand *ohne* Wächter.
4. **Die Baseline-Vorlagen unter `.harness/baseline/` bleiben unberührt.**
   Sie sind vendored und pin-gebunden; eine Änderung dort bräche
   `make baseline-verify` gegen `SHA256SUMS` — *Schicht-Abgrenzung*.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** **Die Form des Gegenstands steht in §3 ausgeschrieben, bevor
      gezählt wird** — was genau ein überflüssiges Leerzeichen ist und was
      nicht. Die Trefferliste wird stichprobenweise dagegen gehalten, nicht
      nur ihre Zahl gelesen.
- [ ] **(2)** §Source precedence und §Guides tragen die schlanke Form:
      genau ein Leerzeichen beidseits jedes Zellinhalts, Trenner `| --- |`.
      Die drei übrigen Tabellen der Datei sind **unverändert**.
- [ ] **(3)** **Die zwei gemessenen Konfigurations-Kanten sind nach der
      Änderung nachweislich intakt** — mit echter Ausgabe, nicht mit einem
      behaupteten Exit-Code: `make gate-consistency` (das Modul `targets`
      liest `harness/README.md` als `doc-tables` und findet in §Guides
      **ein** `make`-Target) und `make doc-check` (die `structure`-Regel auf
      §Sensors mit `cell-min-chars`).
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Die Form des Gegenstands — ausgeschrieben, bevor gezählt wird.** Das
verlangt
[`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md),
und die Regel hat sich in der Voruntersuchung dieses Slice bereits **zweimal
bezahlt gemacht** (§6). Als *überflüssiges Leerzeichen* zählt:

- ein Lauf von **zwei oder mehr** Leerzeichen unmittelbar vor einem `|`
  innerhalb einer Tabellenzeile;
- in der Trennzeile die Bindestrich-Strecke über `---` hinaus.

**Nicht** dazu zählt das **eine** Leerzeichen nach und vor einem `|` — das ist
die übliche Zell-Begrenzung und steht auch in der Vorlage. Die Unterscheidung
ist der ganze Punkt: Ein Muster, das *ein* Leerzeichen mitnimmt, trifft
**jede** Tabellenzeile und liefert eine plausible, falsche Zahl.

**Gemessener Ausgangsstand** (gegen diese Form, Stichprobe geprüft):

| Tabelle | Zeilen | Padding-Zeilen | Trennzeile |
|---|---|---|---|
| §Source precedence | 11 | 10 | 226 Zeichen |
| §Guides | 12 | 10 | 455 Zeichen |
| §Sensors | 25 | **0** | `\| --- \|` |
| §Werkzeuge | 22 | **0** | `\| --- \|` |
| §Gate-Taxonomie | 5 | **0** | `\| --- \|` |

**Schritte:**

1. Padding und Trennzeilen der zwei Tabellen normalisieren — **marker-**,
   nicht zeilennummern-basiert (§6 nennt den Grund).
2. Gegenprobe, dass **nur** Whitespace fiel: ein Vergleich, der Leerzeichen
   ignoriert, muss gegen den Vorstand **leer** sein.
3. Die zwei Kanten fahren: `make gate-consistency`, `make doc-check`.
4. `make gates`, dann Handoff an den unabhängigen Review.

**Was der Plan ausdrücklich nicht mitnimmt**, steht in §1 — die vier Punkte
binden den Lauf, auch wenn unterwegs etwas Naheliegendes auffällt.

## 4. Trigger

**Beanspruchung:** WIP-Limit frei (`in-progress/` ist leer, gemessen nach der
Closure von slice-216), und der Auftraggeber hat den Zuschnitt am 2026-09-08
ausdrücklich beauftragt.

**Rückführung nach `next/`** (`in-progress→next`): wenn sich zeigt, dass die
Normalisierung mehr als Whitespace bewegt — etwa weil eine Zelle ein
bedeutungstragendes Doppel-Leerzeichen führt (ein Zeilenumbruch-Marker in
manchen Markdown-Dialekten). Dann ist der Gegenstand nicht die
Whitespace-Klasse, sondern eine Inhalts-Frage, und der Slice ist falsch
geschnitten.

**Rückführung nach `open/`** (`in-progress→open`): wenn eine der zwei Kanten
nach der Änderung meldet und die Ursache nicht in der Änderung liegt, sondern
in einer Zusage des Produkts über Tabellen-Zellen. Das wäre ein Produkt-Befund
und blockiert diesen Slice.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `make gates` grün mit echter Ausgabe, unabhängiger
Review durchgeführt und seine Befunde eingearbeitet, Closure-Notiz geschrieben,
Register fortgeschrieben, jedes Risiko aus §6 mit einem der drei Ausgänge.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Messung des eigenen Gegenstands ist in diesem Slice schon zweimal
  gescheitert — beide Male in der Voruntersuchung.** Erstens traf das Muster
  `[^ ]  *|` auch **ein** Leerzeichen (`*` ist null oder mehr) und meldete
  75 von 75 Zeilen statt 20. Zweitens vergaß das Trennzeilen-Muster
  `^\|[ -]+\|$` die **inneren** Pipes im Zeichensatz und fand null Zeilen, wo
  zwei stehen. Beide Male lieferte das falsche Muster eine **plausible** Zahl,
  und beide Male fiel es nur auf, weil das Ergebnis gegen den Augenschein
  gehalten wurde. **Ein drittes Mal wäre kein Zufall mehr**, und der
  Bruch-Test dagegen ist Schritt 2: Ein whitespace-ignorierender Vergleich
  muss leer sein. — **Ausgang:** \<offen\>
- **Zeilennummern-basiertes Patchen hat in diesem Repo wiederholt Dateien
  beschädigt** (duplizierte Zeilen, abgeschnittene Sätze) — zuletzt in der
  Closure von slice-216, wo eine `awk`-Ersetzung eine Zeile doppelt schrieb,
  die schon dastand. Der Gegenstand hier ist eine Datei mit **75**
  Tabellenzeilen, und die Operation läuft über die meisten davon. — **Ausgang:**
  \<offen\>
- **Die zwei Konfigurations-Kanten sind gemessen, aber ihre Wirkung ist
  nicht.** Dass `targets` die Datei als `doc-tables` liest und `structure` eine
  Zell-Untergrenze auf §Sensors führt, steht in
  [`.d-check.yml`](../../../../.d-check.yml). **Ob** eine dieser Prüfungen
  Padding **mitzählt** — ob also `cell-min-chars` gegen den gepaddeten oder den
  getrimmten Zellinhalt misst —, ist **nicht** gemessen. Für §Sensors ist es
  gleichgültig, weil dort kein Padding steht; für §Guides mit seinem einen
  `make`-Target ist es die offene Frage. **Der Slice beantwortet sie durch
  Fahren, nicht durch Lesen** — und wenn eine Kante meldet, greift §4. —
  **Ausgang:** \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende. Dieses Repo führt **drei** Prüfungen — die
zwei kanonischen und, als Adaption, den Nachtlauf-Stand
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:268-269 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice ändert **ein** getracktes
Artefakt. `harness/README.md` liegt zwar im Harness-Baum, aber die deklarierte
Sub-Area `tools/harness/` ist ein **Code**-Pfad und trägt diese Datei nicht;
sie hier heranzuziehen wäre ein Geltungsbereich, den die Modus-Deklaration
nicht hergibt.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **38** Verzeichnisse über beide
Kürzel; `BEO-HARN` einzeln geöffnet — der eine offene Eintrag dort betrifft
`--check-latest` und berührt diesen Slice nicht). **Drei** Einträge sind
einschlägig, einer davon **nicht**, und das ist die interessantere Feststellung:

- [`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (4×, verkörpert) — **der tragende.** Sein Ableiter ist DoD (1), und er hat
  in der Voruntersuchung **zweimal** gegriffen (§6). Ein Slice, dessen ganzer
  Gegenstand eine Zeichen-Klasse ist, steht und fällt mit ihrer Definition.
- [`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
  (3×, verkörpert als [`MR-070`](../../../../harness/conventions.md#mr-070))
  — **einschlägig dem Thema nach, aber nicht dem Geltungsbereich nach.**
  [`MR-070`](../../../../harness/conventions.md#mr-070) gilt *„jeder mechanischen Ersetzung über **mehr als eine**
  Datei"*; dieser Slice fasst **eine** an. Die Regel greift also nicht, und
  das Feld ist gelesen statt der Titel — genau, was
  [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  (18×) verlangt. **Was trotzdem bleibt**, ist die Frage hinter der Regel, und
  §1 beantwortet sie: die eingefrorene Klasse hier sind die vendored
  Baseline-Vorlagen, und Abgrenzung (4) nimmt sie aus.
- [`semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md)
  (15×) — die Kanten sind **gesucht und gefunden**: keine `d-check:cite`-Spanne
  ankert zeilengenau in die Datei (gemessen: kein Treffer), aber zwei
  Konfigurations-Einträge lesen sie. §6 führt sie als drittes Risiko.

**Keiner der drei erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-08 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (jüngster Lauf 2026-09-08T05:31:31Z) und `image-scan.yml`
(2026-09-07T08:21:32Z). Für diesen Slice ohne Bezug: Er ändert kein gepinntes
Artefakt und keinen Sensor. **Notiert, weil die Prüfung unbedingt ist.**

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** **hoch, und zwar an der Vorlage statt am Bestand.**
  Die Tabellen-Form ist in
  [`harness/README.template.md`](../../../../.harness/baseline/v6.5.0/templates/harness/README.template.md)
  vorgegeben (39 Zeilen, null Padding) — die Datei ist eine ausgefüllte Kopie
  davon. Damit ist die Zielform nicht Geschmack, sondern die adoptierte Form.
- **Phase-Reife:** Phase 5. Die Datei ist seit dem Bootstrap in Gebrauch (117
  Commits), ihre Struktur ist gewachsen und gewächtert.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig für den Bestand, mittel für die
  Operation.** Zu inventarisieren ist nichts — Doc und Code können hier nicht
  divergieren, weil es keinen Code gibt. Das Risiko sitzt vollständig in der
  Ausführung: in der Definition der Zeichen-Klasse (zweimal gescheitert) und
  im Patch-Verfahren (in slice-216 gescheitert).
- **Reconciliation-Aufwand:** keiner (GF).
