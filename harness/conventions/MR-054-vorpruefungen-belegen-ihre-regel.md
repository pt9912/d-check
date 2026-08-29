# MR-054 — Slice-Vorprüfungen belegen ihre Regel mit `d-check:cite` (Nachtrag zu MR-051)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon verlangt die beiden
  Vorprüfungs-Schritte
  ([`modul-05-planning-harness.md` §Zwei Schritte vor der Modus-Begründung](../../.harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md)),
  sagt aber nichts darüber, **wie** ihre Ausführung belegt wird — er setzt kein
  Werkzeug voraus, das einen Beleg prüfen könnte. Diese Adaption ergänzt die
  Form, nicht die Pflicht.
- **Datum:** 2026-08-29
- **Geltungsbereich:** die beiden **kanonischen** Vorprüfungs-Blöcke im
  `Vorgelagert`-Abschnitt jedes neu angelegten Slice-Plans — Sub-Area-Wahl und
  Beobachtungs-Sichtung. **Nicht** der dritte Block (Nachtlauf-Stand,
  [`MR-053`](../conventions.md#mr-053)); **nicht** rückwirkend für Slices in
  `done/`.

## Adaption

**Ein Vorprüfungs-Block war bis hierher eine Selbstauskunft.** Der Autor schrieb
hin, dass er geprüft hat; ob er die Regel, die den Schritt vorschreibt, gelesen
hat, stand nirgends und war an nichts gebunden.

**Ab jetzt trägt jeder der beiden Blöcke eine `d-check:cite`-Direktive** auf die
Zeilen des Regelwerks, die ihn vorschreiben, mit dem wörtlichen Zitat darunter.
`citations` prüft es wortgleich gegen die Quell-Spanne — fail-closed, im inneren
Loop ([ADR-0045](../../docs/plan/adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)).
Wer die Spanne korrekt setzt, hat die Datei geöffnet.

**Der Anlass ist gemessen und liegt in diesem Repo.** slice-168, -169 und -170
liefen mit vollständig ausgefüllten Vorprüfungs-Blöcken durch, während der
zuständige Zyklus-Abschnitt des Regelwerks (`modul-01`, `-05`, `-06`, `-08`)
nicht gelesen war. Folge: Review und Verifikation fielen aus, das
Beobachtungs-Register wurde nie fortgeschrieben, drei Slices gingen mit offenem
DoD-Haken nach `done/`. **Kein Block war falsch ausgefüllt** — sie waren alle
nur Deklaration.

## Warum der dritte Block keine Direktive trägt

Der Nachtlauf-Schritt ist selbst eine Repo-Adaption
([`MR-053`](../conventions.md#mr-053)); sein Ziel liegt **nicht** unter
`.harness/baseline/`. Eine Direktive auf ein repo-eigenes Ziel meldet bei
**jeder** Änderung des Ziels, nicht nur bei einer inhaltlichen — der Preis
stünde in keinem Verhältnis zum Beleg. Das ist ein Entscheid, keine Auslassung.

## Der Preis, und wie er getragen wird

[`MR-051`](../conventions.md#mr-051) legt das Neu-Ankern der Spannen in die
Bump-Prozedur und grenzt seinen Geltungsbereich auf die **lebenden** Dokumente
ein — wörtlich *„Nicht `done/`, nicht `docs/reviews/`, nicht
`conventions/done/` — dort steht keine."* Diese Adaption macht die Aussage
absehbar unwahr: Ein Slice mit Direktiven **schließt** irgendwann, und dann
steht dort eine.

**Die Antwort ist nicht, den Beleg zu entfernen, sondern seinen Geltungszeitraum
zu benennen.** `citations.scope` nimmt die drei eingefrorenen Verzeichnisse aus.
Der Beleg zählt **zum Zeitpunkt seiner Prüfung**: vor dem `git mv` läuft die
Direktive im inneren Loop; danach ist sie Lauf-Beleg wie ein Review-Report. Die
**Link**-Achse löst dasselbe längst so — `ignore-refs` nimmt die eingefrorenen
Verweise auf entfernte Baseline-Bäume quell-skopiert aus.

**Zwei Messungen dazu, beide vor dem Setzen gefahren:**

1. **In den drei Verzeichnissen steht heute keine wirksame Direktive.** Ein
   rohes `grep` findet 13 Treffer; alle sind **Erwähnungen** der Syntax in
   Inline-Code, die [ADR-0060](../../docs/plan/adr/0060-citations-marker-scan-geteilte-prosa-antwort.md)
   ausdrücklich nicht als Direktive zählt — und genau deshalb sind die Gates
   grün. Der Ausschluss ist **Vorsorge, keine Reparatur**; MR-051s Aussage war
   zum Zeitpunkt ihres Schreibens richtig.
2. **Ein modul-lokaler Scope ersetzt den globalen** ([`DC-FA-CONF-002`](../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope)). Der
   erste Versuch ohne die beiden `.harness/`-Ausschlüsse hob die Scan-Menge von
   **587 auf 653** Dateien — `citations` hätte die vendored Baseline
   mitgescannt. Die Wiederholung der globalen Ausschlüsse im Modul-Scope ist
   deshalb Pflicht, nicht Redundanz.

## Grenze

**Ein Zitat belegt Zugriff, nicht Verständnis.** Es lässt sich aus einem anderen
Slice kopieren, und dann hat niemand etwas gelesen. Der Beleg ist damit
schwächer, als er aussieht — und ein Beleg, der stärker aussieht, als er ist,
ist die Klasse, gegen die diese Adaption gebaut ist. Sie hätte den Anlassfall
trotzdem verhindert: Wer die Spanne korrekt setzt, hat die Datei offen, in der
zwei Abschnitte weiter der Zyklus steht.

## Auflösungs-Trigger

**Permanent, solange die Vorprüfungs-Blöcke bestehen.** Fällt einer von ihnen
aus dem Kanon, entfällt seine Direktive mit ihm.

**Ein Beobachtungs-Trigger daneben:** meldet `citations` nach einem
Baseline-Bump mehr Spannen-Brüche in Slice-Plänen als in den Skills, ist die
Direktive im Planungs-Pfad teurer als gedacht — dann ist zu entscheiden, ob der
Beleg an eine **stabilere Adresse** gebunden wird (Abschnitts-Anker statt
Zeilenspanne) oder ob er den Preis wert bleibt.
