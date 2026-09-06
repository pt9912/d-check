# MR-054 — Slice-Vorprüfungen belegen ihre Regel mit `d-check:cite` (Nachtrag zu MR-051)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon verlangt die beiden
  Vorprüfungs-Schritte
  ([`modul-05-planning-harness.md` §Zwei Schritte vor der Modus-Begründung](../../.harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md#zwei-schritte-vor-der-modus-begründung)),
  sagt aber nichts darüber, **wie** ihre Ausführung belegt wird — er setzt kein
  Werkzeug voraus, das einen Beleg prüfen könnte. Diese Adaption ergänzt die
  Form, nicht die Pflicht.
- **Datum:** 2026-08-29
- **Geltungsbereich:** die beiden **kanonischen** Vorprüfungs-Blöcke im
  `Vorgelagert`-Abschnitt jedes neu angelegten Slice-Plans — Sub-Area-Wahl und
  Beobachtungs-Sichtung. **Nicht** der dritte Block (Nachtlauf-Stand,
  [`MR-053`](../conventions.md#mr-053)); **nicht** rückwirkend für Slices in
  `done/`. Eine `cite`-Direktive **anderswo** in einem Planungs-Dokument — in
  der Zielbegründung eines Slice, im Wellen-Ziel — ist davon unberührt: erlaubt,
  von dieser Adaption aber nicht gefordert. Sie regiert **zwei Blöcke**, nicht
  die Direktive überhaupt.

## Adaption

**Ein Vorprüfungs-Block war bis hierher eine Selbstauskunft.** Der Autor schrieb
hin, dass er geprüft hat; ob er die Regel, die den Schritt vorschreibt, gelesen
hat, stand nirgends und war an nichts gebunden.

**Ab jetzt trägt jeder der beiden Blöcke eine `d-check:cite`-Direktive** auf die
Zeilen des Regelwerks, die ihn vorschreiben, mit dem wörtlichen Zitat darunter.
`citations` prüft es wortgleich gegen die Quell-Spanne — fail-closed, im inneren
Loop ([ADR-0045](../../docs/plan/adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)).
Wer die Spanne korrekt setzt, hat die Datei geöffnet.

**Die Spanne trägt die *vorschreibende* Zeile, nicht irgendeine aus demselben
Absatz.** Ein Schritt-Punkt des Regelwerks führt neben der Vorschrift auch
Nebenregeln — Notier-Pflichten, Lesehinweise. Eine Direktive darauf ist grün und
belegt trotzdem nicht, was ihre Anrede behauptet; verschwindet die Nebenregel
beim Bump, entfällt der Beleg, während der Schritt weiter vorgeschrieben ist.
Das ist die [`BEO-ALL/citation-stretched-beyond-scope`](../../docs/plan/planning/observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)-Klasse in
Direktiven-Form.

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
Direktive im inneren Loop; danach ist sie Lauf-Beleg wie ein Review-Report.

**Das dritte Verzeichnis liegt aus einem anderen Grund in derselben Liste, und
der gehört benannt.** Ein Review-Report unter `docs/reviews/` hat **keine
Live-Phase** — er entsteht direkt an seinem Endort, das Argument „vor dem
`git mv` geprüft, danach Beleg" greift bei ihm nicht. Er ist **von Geburt an**
Lauf-Beleg. Eine Direktive dort einmal beim Schreiben zu prüfen, kaufte einen
schmalen Beleg gegen eine **dauerhafte** Neu-Anker-Pflicht — genau die, die
[`MR-051`](../conventions.md#mr-051) für diese Klasse ausdrücklich nicht will.

**Die Link-Achse steht Pate für die Richtung, nicht für die Reichweite.**
`ignore-refs` nimmt eingefrorene Verweise auf entfernte Baseline-Bäume aus, ist
dabei aber quell-skopiert **und ziel-benannt**: die Datei bleibt in der
Scan-Menge, nur ein benanntes Ziel zählt dort nicht mehr.
`citations.scope.ignore` nimmt die **ganze Datei** aus dem Modul. Das ist
gröber — und es ist gröber, weil `citations` kein Ziel-Muster kennt; die feinere
Form gibt es hier schlicht nicht.

**Der Preis ist in allen drei Verzeichnissen derselbe und wird hier genannt:**
mit dem Ausschluss entfällt dort auch der **fail-closed**-Pfad. Eine strukturell
malformte Direktive nahm vorher den ganzen Lauf mit; in den drei Verzeichnissen
tut sie es nicht mehr, sondern fällt erst einem Leser auf.

**Zwei Messungen dazu, beide selbst gefahren:**

1. **In den drei Verzeichnissen steht heute keine wirksame Direktive.** Gemessen
   mit dem **Produkt**, nicht mit `grep`: ein Lauf über die drei Verzeichnisse
   mit `modules: [citations]` und ohne jeden Ausschluss meldet **455 Dateien, 0
   Befunde**, Exit 0. Ein rohes `grep` findet dort 13 Marker-Treffer — **zwölf**
   in Inline-Code, **einer** in einem Fenced-Block; beide Formen zählt
   [ADR-0060](../../docs/plan/adr/0060-citations-marker-scan-geteilte-prosa-antwort.md)
   ausdrücklich nicht als Direktive. Der Ausschluss ist **Vorsorge, keine
   Reparatur**; MR-051s Aussage war zum Zeitpunkt ihres Schreibens richtig, und
   rückwirkend legt der Ausschluss nichts still.
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

**Kein Sensor hält sie.** Ein Slice-Plan ohne Direktiven ist grün; geprüft wird
nur, was dasteht. Die Adaption wirkt über die Zustellung — Index,
[`AGENTS.md`](../../AGENTS.md) §5 — und über das Review, nicht über ein Gate.

## Auflösungs-Trigger

**Permanent, solange die Vorprüfungs-Blöcke bestehen.** Fällt einer von ihnen
aus dem Kanon, entfällt seine Direktive mit ihm.

**Ein Beobachtungs-Trigger daneben:** meldet `citations` nach einem
Baseline-Bump mehr Spannen-Brüche in Slice-Plänen als in den Skills, ist die
Direktive im Planungs-Pfad teurer als gedacht — dann ist zu entscheiden, ob der
Beleg an eine **stabilere Adresse** gebunden wird (Abschnitts-Anker statt
Zeilenspanne) oder ob er den Preis wert bleibt.
