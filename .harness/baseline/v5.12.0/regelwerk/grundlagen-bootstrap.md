## Harness-Bootstrap
<!-- Quelle: [grundlagen/bootstrap.md](https://github.com/pt9912/ai-harness-course/blob/v5.12.0/kurs/de/grundlagen/bootstrap.md) -->

### Harness-Bootstrap

*Harness-Bootstrap* bezeichnet den **Einstiegsprozess** in den
Harness-Lebenszyklus eines Repos — der Weg von "leeres Repo" oder
"Repo ohne Harness" bis zur Stelle, an der inhaltliche Arbeit (Slices,
Code) auf einem etablierten Harness aufsetzt. Es ist eine *Trajektorie
durch Dokument-Zustände*, kein *Ereignis*. Konkreter Walkthrough mit
Schritten in [Modul 1](modul-01-entwicklungszyklus.md#source-precedence-block).

> **Begriffsklärung:** "Harness-Bootstrap" meint hier den
> Einstiegsprozess in den Harness. Nicht zu verwechseln mit
> *Bootstrap-aware Gate* ([Modul 13](modul-13-quality-gates.md)) — das ist ein
> einzelnes Gate mit Reifestufe und Hochschalt-Trigger (Coverage 0 →
> 70 %). Beide Begriffe teilen das Wort, sind strukturell verschieden:
> *Harness-Bootstrap* betrifft den **Repo-Lebenszyklus**,
> *Bootstrap-aware Gate* die **Reifestufe eines Sensors**.

#### Was ist eine Sub-Area?

Eine *Sub-Area* ist eine **Doku-/Code-Sektion, die als Träger einer
Modus-Entscheidung dient** — mit eigener Konventions-Härte (eigene
`MR-NNN` möglich), eigener Inventur-Linie und eigener Pfad-/Datei-Familie
im Repo. Sie ist nicht das Repo (zu grob) und nicht der Slice (ein Slice
*berührt* Sub-Areas, *trägt* aber keinen Modus).

*Modul, Verzeichnis, Komponente* (siehe §Modus pro Sub-Area unten) sind
die **typischen Träger** — sie nennen, *welche Strukturen* eine Sub-Area
sein können. Ob eine konkrete Struktur als Sub-Area **qualifiziert**,
entscheiden drei Inklusions-Achsen (bottom-up):

| Achse                         | Test                                                      | erfüllt, wenn …                                                                                                                   |
| ----------------------------- | --------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| **1 — Konventions-Härte**     | Ist eine eigene `MR-NNN`-Adaption plausibel formulierbar? | … die Sektion eine eigene Strukturregel tragen *könnte* (nicht: schon trägt).                                                     |
| **2 — Inventur-Linie**        | Ist eine eigene Diskrepanz-Bericht-Zeile sinnvoll?        | … Code-Bestand und Doku-Aussage dieser Sektion als Paar abgleichbar sind, ohne dass eine Nachbar-Sub-Area mitgezogen werden muss. |
| **3 — Struktureller Cluster** | Gibt es eine eigene Pfad-/Datei-Familie?                  | … ein eigenes Verzeichnis, Dateimuster oder Konventions-Präfix die Sektion trägt.                                                 |

**Schwelle: mindestens zwei der drei Achsen.** Eine Achse allein ist zu
schwach — der typische Fall ist *Struktur ohne Substanz*: ein Verzeichnis
existiert (Achse 3), hat aber keine eigene Konvention (Achse 1) und keine
eigenständig abgleichbare Inventur-Linie (Achse 2). Das ist noch keine
Sub-Area, sondern eine **Sub-Area-Aspirantin** — in winzigen Repos
normal, mit wachsender Struktur wird daraus eine Sub-Area.

**Positiv-Beispiele:**

- *Audit-Logging* — eigene MR-Adaption denkbar (Format-Standard für
  Log-Einträge, Achse 1), eigene Inventur-Linie (entstehen alle
  Audit-Events wie spezifiziert?, Achse 2), eigener `services/audit/`-
  Pfad-Cluster (Achse 3). Alle drei → klar Sub-Area.
- *Test-Infrastruktur* — eigenes Pfadnaming-Schema (Achse 3) und eine
  eigene Inventur-Linie (Tests ohne `LH-*`-ID als Diskrepanz, Achse 2).
  Zwei von drei → Sub-Area.

**Negativ-Beispiele:**

- *"Backend"* ist zu grob — verletzt Achse 1 (keine *einzelne*
  `MR-NNN`-Adaption denkbar; API-Pattern, Persistence-Layout und
  Hintergrund-Jobs bräuchten je eigene) und Achse 3 (mehrere
  Pfad-Familien). *"Backend"* bündelt typischerweise *drei* Sub-Areas.
- *"Frontend"* — analog: eigene Konventionen pro Schicht (Komponenten,
  State, Styling), keine gemeinsame Inventur-Linie. Auch hier:
  ausdifferenzieren, nicht als *eine* Sub-Area führen.

> **Abgrenzung zu den vier Modus-Pflichtkriterien.** Die drei Achsen
> hier beantworten *ob eine Struktur eine Sub-Area ist* (Granularitäts-
> Gate). Sie sind **nicht** zu verwechseln mit den vier Pflichtkriterien,
> mit denen [Modul 5](modul-05-planning-harness.md#ziel-form-sub-area-modus-begründung)
> begründet, *welcher Modus* (GF/BF/Hybrid) für eine bereits erkannte
> Sub-Area gilt (Konventionen-Dichte · Phase-Reife · Evidenz-/Diskrepanz-
> Risiko · Reconciliation-Aufwand). Erst Inklusion (hier), dann
> Modus-Wahl (Modul 5).

**Was heißt „berührt"?** Ein Slice *berührt* eine Sub-Area, wenn er ihren
**Doku-/Code-Abgleich bewegt** — wenn er also ihre Konventions-Härte oder ihre
Inventur-Linie verändert. Das ist die Bedingung, die entscheidet, für welche
Sub-Areas ein Slice einen Begründungsblock schreibt ([Modul 5 §8](modul-05-planning-harness.md#ziel-form-sub-area-modus-begründung)) und unter welche
Sub-Area eine Beobachtung ins Register geht ([Modul 6 §Das Beobachtungs-Register](modul-06-roadmap.md#das-beobachtungs-register-modul-6)). Zwei Wege führen dorthin,
und nur einer steht im Diff:

- **Pfad-Berührung** — der Slice ändert eine Datei aus dem Pfad-Cluster der
  Sub-Area. Mechanisch ablesbar, aber **nicht hinreichend**: Additive Arbeit
  *innerhalb* einer bereits deklarierten Konvention bewegt den Abgleich nicht.
  Ein Slice, der eine Testdatei nach dem geltenden Schema ergänzt, berührt
  *Test-Infrastruktur* nicht — sonst trüge jeder Slice diesen Block, und der
  Block verlöre genau die Aussage, für die es ihn gibt.
- **Aussagen-Berührung** — der Slice ändert eine Aussage, gegen die der Cluster
  abgeglichen wird, ohne eine seiner Dateien anzufassen. Eine ADR, die die
  Schreib-Semantik des Index festlegt, berührt die Implementierungs-Sub-Area
  auch dann, wenn in dieser Welle noch kein Index-Code entsteht.

> **Grenze — ehrlich benannt:** Keine der beiden Hälften ist maschinell
> entscheidbar. Was der Diff liefert, ist eine **Kandidatenliste**: Ein Gate
> kann verlangen, dass jeder Pfad-Kandidat in §8 entweder einen Block oder eine
> Abweisung mit Grund bekommt. Ob die Abweisung trägt — und ob ein Kandidat
> fehlt, den nur die Aussagen-Berührung findet —, bleibt Urteil.

**Aggregation — die Kehrseite der Inklusion.** Wie die Schwelle ein
*Zuviel an Struktur* abweist (die Aspirantin oben), weist dieselbe Logik
rückwärts gelesen ein *Zuwenig an Trennung* ab: Zwei Sub-Areas, die
**permanent dieselben Trigger** erzeugen *und* **dieselbe Modus-Aussage**
tragen, sind in Wahrheit *eine* — sie getrennt zu führen erzeugt zwei
Inventur-Linien ohne eigene Diskrepanz (Anti-Refactoring). Die
Diagnose-Frage ist die Achsen-Frage rückwärts: *„Feuern die beiden je
**unabhängig** — eigener Trigger, eigene `MR-NNN`?"* Über mehrere Wellen
nein → zusammenführen; sobald eine Hälfte eine eigene Adaption oder
Inventur-Linie bekommt (Achse 1/2 divergiert) → trennen. Aggregation ist
damit keine Einmal-Entscheidung, sondern eine wiederkehrende
Wartungs-Praxis. Faustregel: *was nie getrennt feuert, ist
eine Sub-Area; eine Sub-Area, deren Hälften auseinanderdriften, sind
zwei.* Beispiel: die sechs Sprach-Skelette (`go/`, `python/`,
…) werden *nicht* als sechs `Implementierung`-Sub-Areas geführt, sondern
als *eine* — sie teilen Spec und Modus (alle GF) und tragen nie eine
*unabhängige* Modus- oder Trigger-Entscheidung; die per-Sprache-Stilunterschiede
(`gofmt` vs. `black`) sind Sub-Sub-Area-Nuancen, keine eigenen
Inventur-Linien. Split-Trigger: kippte ein Skelett nach BF (etwa ein
Alt-Port mit Bestandscode), bekäme es eine eigene Modus-Aussage — und
*dann* wäre es eine eigene Sub-Area. Die Gegenrichtung zeigt
`harness/conventions.md`: `Test-Infrastruktur`, `Verifikation` und
`Replay-/Eval-Infrastruktur` sehen ähnlich aus („Korrektheits-Sensoren"),
sind aber *drei* Sub-Areas, weil Achse 1 divergiert — sie zu mergen wäre
der „zu grob"-Fehler.

#### Modus pro Sub-Area: Greenfield vs Brownfield

Pro Sub-Area eines Repos (Modul, Verzeichnis, Komponente) wird ein
**Modus** deklariert (im Adaptions-Block von
`harness/conventions.md`). Die Modus-Wahl bestimmt die
*Trigger-Richtung* — wer wem folgt:

| Modus               | Trigger-Richtung          | Bild im Kopf                                                                                         |
| ------------------- | ------------------------- | ---------------------------------------------------------------------------------------------------- |
| **Greenfield** (GF) | Doc → Code                | Spec führt, Code folgt. "Wir versprechen X, dann liefern wir X." Steady-State.                       |
| **Brownfield** (BF) | Code → Doc                | Code existiert, Doku folgt. Inventur des Bestands. **Übergangs-Modus mit Konvergenz-Auftrag** zu GF. |
| **Hybrid**          | gemischt pro Sub-Sub-Area | Realistisch: alte Komponenten BF, neue GF.                                                           |

**Konvergenz-Auftrag.** BF ist *keine Daueroption*. Jede BF-Sub-Area
trägt eine **Graduation-Bedingung** (im Adaptions-Block dokumentiert):
*was muss erfüllt sein, damit die Sub-Area in GF-Modus wechselt?*
Typisch: alle entdeckten Diskrepanzen aufgelöst (als Carveouts oder
Reconciliation-Slices); Spec/ADR/Sensors decken Code-Stand ab;
ID-Schema retrofitted. Eine BF-Sub-Area ohne Graduation-Plan ist eine
*permanente Ausnahme als temporär getarnt* — analog zur
Carveout-Disziplin in [Modul 7](modul-07-carveouts.md).

Permanente BF-Erklärung (für Code, der absehbar entfernt wird —
Legacy, Drittsystem-Adapter) ist möglich, mit Begründung und
Folge-Slice.

**Beide Modi regeln dieselbe Achse: Doc ↔ Code.** Die Spalte
*Trigger-Richtung* sagt es wörtlich — GF ist `Doc → Code`, BF ist `Code → Doc`.
Eine dritte Beziehung fällt deshalb durch beide hindurch: **adoptierte Norm ↔
ausgefülltes Artefakt**. Wer eine Regelwerks-Migration für einen
Brownfield-Fall hält, weil dort *„Inventur des Bestands"* steht, greift zum
falschen Werkzeug: BF regelt, ob **Code oder Doku führt** — bei einer Migration
sind beide längst da und stimmen miteinander überein; abweichen kann das
Artefakt von der *adoptierten Norm*. Für diese Achse ist der Freshness-Audit
zuständig
([`modul-02-harness-bootstrap.md` §Freshness-Audit](modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)),
nicht die Modus-Wahl.

#### Sektionsweise Reife: Phasen pro Dokument

Ein Harness-Dokument ist während Bootstrap nicht "entweder leer oder
fertig". Sektionen reifen mit unterschiedlichem Tempo durch fünf
Phasen:

| Phase        | Beschreibung                                                                     |
| ------------ | -------------------------------------------------------------------------------- |
| 0 — leer     | Datei existiert nicht                                                            |
| 1 — Skelett  | Template kopiert, Pflichtgliederung mit Platzhaltern                             |
| 2 — Outline  | Top-Level ausformuliert, Details `<…>`                                           |
| 3 — partiell | einige Sektionen voll, andere noch `<…>`                                         |
| 4 — kohärent | alle Sektionen gefüllt, intern konsistent — *freigegeben* für Verweise von außen |
| 5 — stabil   | Änderungen nur über Change-Process                                               |

*Sektionen* eines Dokuments können in unterschiedlichen Phasen sein.
Beispiel: §Source precedence von `harness/README.md` kann durch
Template-Adoption früh auf Phase 2 sein, während §Sensors auf Phase 1
verharrt, bis das Makefile existiert. **Sektionsweise Reife ist Regel,
nicht Ausnahme** — Schreibreife wird sektionsweise beurteilt, nicht
dateiweise.

#### Vier Trigger-Klassen

Während Bootstrap (und auch danach im Steering-Loop) lösen Änderungen
in einem Dokument *Folgeaktionen* in anderen aus. Vier Klassen:

| Klasse                      | Wirkung                                                                                                                        | Beispiel                                                                                                                                                                                                                                                                                                                                     |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Sync-Trigger**            | Pointer in einem Dokument muss in einem anderen ergänzt werden                                                                 | Neuer Eintrag in `conventions.md` → Pointer in `harness/README.md`                                                                                                                                                                                                                                                                           |
| **Promotion-Trigger**       | Eintrag wandert aus "Nicht behauptet"-Block in Haupt-Tabelle                                                                   | Make-Target real im Makefile entstanden → Sensor-Zeile gepromoted                                                                                                                                                                                                                                                                            |
| **Cross-Reference-Trigger** | Verlinkung zwischen Dokumenten, normativ **nur volatil→stabil** ([`grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)](grundlagen-referenz-richtung.md#referenz-richtung-sdp-wer-darf-wen-referenzieren)) | Neue ADR *deklariert aufwärts, was sie schärft* (ADR → Spec-§) und referenziert die Anforderung; der Acceptance-Trigger zieht die Spec nach. Ein Spec→ADR-Rückzeiger im bindenden Text existiert nicht (auch nicht als Quellen-Spalte) — Provenance nur in der Historie-Tabelle (Regel 5); der Referenz-Richtungs-Gate erzwingt das über alle Straten |
| **Acceptance-Trigger**      | Phase-Übergang via Sign-off (z. B. ADR Proposed → Accepted)                                                                    | ADR-Review-Runde abgeschlossen → bindend                                                                                                                                                                                                                                                                                                     |

Trigger werden zwischen Bootstrap-Schritten ausgewertet — sie sind die
"Inbox" der nicht-Vorderscene-Arbeit. Eine zwischen Schritten
übersehene Trigger-Pflicht ist ein häufiges Drift-Symptom.

#### Harness-Bootstrap-Ende vs Workflow-Beginn

Harness-Bootstrap ist *abgeschlossen*, wenn der Repo bereit ist für
inhaltliche Slices. In **Greenfield**: erster ADR akzeptiert,
Roadmap-Outline mit Welle-Sequenz, Sensors-Roster als "Nicht
behauptet"-Block. In **Brownfield**: Reconciliation-Backlog steht,
Konvergenzpfad zu GF ist sichtbar (mit ersten Reconciliation-Slices in
`open/`). Ab dann übernimmt der **Workflow** (Slice-Lebenszyklus,
Modul 5–9). Bootstrap und Workflow sind getrennte Lebenszyklen — kein
Übergang ohne Sichtbarkeit.

**Der Reconciliation-Backlog hat einen Ort:** `docs/plan/planning/reconciliation.md`,
flach neben dem Beobachtungs-Register. Eine Zeile je Fund des Rückbaus — Fund,
Sub-Area, Klasse, auflösendes Artefakt, Stand. Ohne diese Datei wäre die Menge
nicht bestimmbar: Weder die Kennung eines Carveouts noch seine sechs Kopffelder
sagen, dass er aus der Inventur stammt, und in `open/` sieht man einem Slice
nicht an, ob er eine Zusage einlöst oder eine Lücke schließt.

Damit ist das Bootstrap-Ende-Kriterium nachzählbar: *„Der Backlog steht"* heißt
**jeder Fund hat eine Zeile, und jede Zeile trägt ihre Auflösung** — ein
Carveout mit Folge-Slice, einen Reconciliation-Slice oder eine retroaktive ADR.
Nicht „das Register ist leer": Beim Bootstrap-Ende ist es im Gegenteil voll.
Leer wird es erst bei der Graduation, und zwar je Sub-Area — deshalb führt es
die offenen Zeilen pro Sub-Area, statt sie über Carveout- und
Slice-Verzeichnisse verstreut zu lassen.

Die dritte Klasse ist mit dem Eintrag erledigt: Eine retroaktive ADR entsteht in
Schritt 7 mit Status *Accepted* (oder *Superseded*) und wartet auf nichts; ihre Zeile geht sofort
nach *aufgelöst*.

#### Einführungs-Reihenfolge über mehrere Repos

Bootstrap gilt pro Repo — in einer Mehrfach-Repo-Landschaft stellt sich
zusätzlich die Frage, *welches Repo zuerst*. Die Antwort folgt der
Repo-Klasse (§Kernbegriffe):

**Beginne immer beim Referenz-Repo**, portiere erst nach erfolgreicher
Steering-Loop-Iteration auf die Flagships (Safety/Control,
Policy/Compliance). Alle Repos parallel mit demselben Master-Prompt zu
treiben skaliert nicht — der Agent verteilt dann halbgare
Standardtexte über alle.

Begründung: das Referenz-Repo ist der *Demonstrator*, in dem
experimentiert werden darf; ein Flagship trägt nicht verhandelbare Hard
Rules und ist der falsche Ort, um eine Konvention zum ersten Mal
auszuprobieren. Was sich im Referenz-Repo über eine Steering-Loop-Runde
bewährt hat, wandert in die Flagships — nicht umgekehrt.

#### Verbindung zum Steering-Loop

Harness-Bootstrap ist im Grunde der **Steering-Loop ([Modul 11](modul-11-verification.md)),
einmal in Folge angewendet, bis Graduation erreicht ist**. Das
Werkzeug ist identisch (Beobachtung → Guide/Sensor); was sich
unterscheidet, ist die Anwendungsphase: Bootstrap = initial bis
Steady-State; Steering-Loop = laufend im Steady-State. Wer den
Steering-Loop versteht, versteht Bootstrap — und umgekehrt.

#### Querverweise

- **[Modul 2 — Harness-Bootstrap](modul-02-harness-bootstrap.md)**: ausgearbeiteter Lehrtext mit GF/BF-Walkthroughs, Trigger-Klassen-Inline-Ankern und Phasen-Karten-Übung — Vollform des Bootstrap-Konzepts.
- **Modul 1 §Schritt 0** ([§Source-Precedence-Block](modul-01-entwicklungszyklus.md#source-precedence-block)): kompakter Vorgriff auf das Modus-Konzept als Eingang in den Lebenszyklus (Baseline und Modus festlegen plus den sechs Folge-Schritten); Vollform in Modul 2.
- **[§harness/conventions.md als Konventionsspeicher](grundlagen-harness-dateien.md#harnessconventionsmd-als-konventionsspeicher)**: Adaptions-Block
trägt Modus-Deklaration pro Sub-Area; Graduation-Bedingung wird dort
dokumentiert.
