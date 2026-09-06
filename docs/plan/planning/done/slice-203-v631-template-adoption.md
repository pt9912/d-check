# Slice slice-203: Sensors-Sektion auf die v6.3.1-Form — Träger, Nicht-Gates, Rollenwechsel

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`MR-065`](../../../../harness/conventions.md#mr-065) (der Bump, der
die Vorlagen brachte). Ob die Adoption eigene `MR`-Einträge braucht,
entscheidet der Slice — dieser Plan setzt es nicht voraus.

**Berührte Spec-Stellen:** — (Harness-Form, keine Spec-Stelle).

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-06.

---

## 1. Ziel

Die Sensors-Sektion in [`harness/README.md`](../../../../harness/README.md)
auf die mit `v6.3.1` adoptierte Form bringen: überladene Verträge wandern in
eigene Träger-Dateien, Nicht-Gates werden als solche **in der Zeile**
gekennzeichnet, und der Minimal Agent Workflow benennt Schritt 8 als
Rollenwechsel statt als Abschluss.

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/sensors/` <!-- d-check:ignore (wird mit diesem Slice angelegt) --> | neu | Träger-Verzeichnis; je Datei ein Target, Ziel-Form ist die vendorte `gate.template.md` |
| `harness/README.md` §Sensors | update | die drei ausgelagerten Zellen werden Index-Zeile mit Link; Nicht-Gates in eine zweite Tabelle |
| [`AGENTS.md`](../../../../AGENTS.md) §6 | update | Schritt 8 ist Rollenwechsel, kein Abschluss; kein Self-Review |

**Gemessen, nicht geschätzt — und die Messgröße ist der ZELLINHALT, nicht die
Zeilenlänge.** Markdown-Tabellen polstern jede Zeile auf die Breite der
Trennzeile; wer die Zeile misst, misst die Formatierung mit und bekommt eine
Rangfolge, die es nicht gibt. Nach Zellinhalt lag der Bestand vor diesem Slice
bei `verify-closure-notes` 4997 · Digest-Achsen 1420 · `nightly-state` 1347 ·
`baseline-verify` 1253 · … · `record-gates` 44 Zeichen; **24 von 31** Zeilen
lagen über 200.

## 3. Ausdrücklich NICHT in diesem Slice

- **Jede weitere Tabelle.** Die Zellenlängen-Regel gilt der Sensors-Sektion
  und dem Adaptions-Block; andere Tabellen des Repos (Roadmap, Meilensteine,
  Beobachtungs-Register) sind nicht geprüft und nicht gekürzt.
- **Die `Bindung`-Spalte als Längen-Bedingung.** Sie ist beim Aufräumen
  vollständig **gefüllt** worden — zwei Zeilen fehlte sie ganz —, trägt aber
  keine Schwelle: Eine Kennungs-Liste ist keine Prosa, und ihre Länge wächst
  mit der Zahl der Bindungen, nicht mit Geschwätzigkeit.
- **Jede Änderung an Produkt-Code.** Der Slice berührt Harness-Form, Config
  und Briefing, nicht `internal/` oder `cmd/`.
## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; Gate-Läufe und die
Closure-Pflichten darunter zählen nicht mit. **Dieser Slice trägt fünf**, und
das ist eine bewusste, vom Auftraggeber getriebene Erweiterung: (4) und (5) kamen
hinzu, nachdem sich zeigte, dass das Produkt die Regel selbst halten kann.
Ohne (4) wäre (1) eine einmalige Aufräumaktion statt einer gehaltenen Form;
(5) wendet dieselbe Regel auf den Ort an, den der Kanon als ihr Vorbild nennt.

- [x] **(1)** `harness/sensors/` <!-- d-check:ignore (entsteht mit diesem Slice) --> trägt jeden Vertrag, der mehr als einen
      Satz braucht, nach der vendorten `gate.template.md`; die Index-Zellen
      sind auf **einen Satz** gekürzt und verlinken ihre Datei. Der Link ist
      die einzige geprüfte Fassung der Zuordnung — eine bloße
      Namenskonvention bliebe still grün, wenn die Datei verschwände.
- [x] **(2)** Nicht-Gates stehen in einer **zweiten Tabelle** und tragen
      `kein Gate` **in der Bindung-Spalte**, nicht in Prosa daneben. Kriterium
      ist, **worüber** ein Target urteilt — Zustand des Repos gegen
      Vorbedingungen des eigenen Laufs. Kein Target, das ein Lauf braucht,
      fehlt in einer der beiden Tabellen.
- [x] **(3)** [`AGENTS.md`](../../../../AGENTS.md) §6 Schritt 8 benennt den
      Rollenwechsel (Handoff an Reviewer, kein Self-Review).
- [x] **(4)** Eine `structure`-Regel in [`.d-check.yml`](../../../../.d-check.yml)
      hält die Zellenlänge beider Tabellen (`table.column`, `cell-max-chars`)
      und läuft in `make gates`. Die Schwelle bildet die **Ein-Satz-Regel** ab,
      nicht den Bestand — eine auf den Bestand kalibrierte Schwelle hätte ihn
      eingefroren, statt die Regel zu verkörpern.
- [x] **(5)** Dieselbe Regel gilt dem **Adaptions-Block** in
      [`harness/conventions.md`](../../../../harness/conventions.md): Titel,
      Geltungsbereich und Ersetzt-Baseline-Regel unter der Schwelle, Prosa in
      der Eintrags-Datei. Beide Tabellen tragen zusätzlich eine **Untergrenze**
      — eine Obergrenze allein lässt die leere Zelle passieren.
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §5 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang.

- **Die Auslagerung wird zur Umschichtung ohne Gewinn.** Der Kanon warnt
  ausdrücklich davor, das neue Verzeichnis zur Halde zu machen, die die Sektion
  vorher war: hinein gehört **Grenze und Bedienvertrag**, nicht der
  Deckungsnachweis des Werkzeugs (der lebt in ADR, Spec-Zeile, Skriptkopf).
  Prüfbar am Ergebnis: steht in einer Träger-Datei ein Satz über *welcher Test
  welche Hälfte deckt*, ist er am falschen Ort. — **Ausgang:** eingetreten — der
  zweite Review fand ihn als M-9 in **fünf** Träger-Dateien, drei davon unter
  der Überschrift *Grenze*, wo ein bestandener Selbsttest keine Lücke ist. Alle
  fünf Stellen entfernt; der Deckungsnachweis lebt in ADR und Modul-Test. Der
  Prüfsatz aus diesem Risiko hat gehalten — er hat den Befund benannt, bevor
  der Reviewer ihn fand, nur hat ihn beim Schreiben niemand angewandt.
- **Der Rest-Bestand bleibt unsichtbar.** — **Ausgang:** entfallen — es gibt
  keinen: Der ursprüngliche Zuschnitt (drei von 30) beruhte auf einer falschen
  Messgröße; nach der Korrektur ist der Bestand vollständig geräumt, und die
  Schwelle in [`.d-check.yml`](../../../../.d-check.yml) hält ihn. Ein Rest,
  den eine Prosa-Zeile ausweisen müsste, ist nicht übrig.
- **Die Nicht-Gate-Zuordnung ist ein Urteil, kein `grep`.** Ob `make bench`
  über den Zustand des Repos urteilt oder über seinen eigenen Lauf, entscheidet
  kein Muster; eine falsche Einordnung macht aus einem Gate ein Werkzeug oder
  umgekehrt — und das erste ist die gefährlichere Richtung. — **Ausgang:**
  eingetreten, und zwar in der benannten gefährlicheren Richtung:
  `doc-complete` und `archive-wave-test` standen als Werkzeuge da, obwohl beide
  über den Repo-Zustand urteilen — `doc-complete` sogar mit einem Recipe, das
  mit `completeness-check` **byte-identisch** ist. Beide in die Gate-Tabelle
  verschoben; das Kriterium im Einleitungstext nennt jetzt den **Gegenstand**
  statt fail-open, weil `guard-probe` und `baseline-probe` fail-closed urteilen
  und trotzdem keine Gates sind.

## 6. Trigger

**Start** (`open` → `in-progress`): [`MR-065`](../../../../harness/conventions.md#mr-065)
ist geschlossen und die Vorlagen liegen unter `.harness/baseline/v6.3.1/templates/`.
WIP-Limit frei (alle Lifecycle-Verzeichnisse leer).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt sich beim Schreiben, dass schon die
  **drei** Verträge nicht in einer Review-Sitzung prüfbar sind, wird die
  Auslagerung auf den einen Ausreißer reduziert und der Rest zieht mit dem
  Bestand mit.
- `in-progress` → `open` (blockiert): Erweist sich die Nicht-Gate-Zuordnung als
  Entscheid, der eine eigene Grundlage braucht (welche Targets urteilen worüber),
  ruht der Slice bis zu diesem Entscheid.

**Closure-Trigger.** Zwei beobachtbare Kriterien **und** ein Lerneintrag:
(a) `make gates` grün, und jede angelegte Träger-Datei ist aus ihrer
Index-Zeile heraus verlinkt — der Link-Sensor hält die Zuordnung, eine bloße
Namenskonvention täte das nicht; (b) keine der drei ausgelagerten Index-Zellen
ist länger als ein Satz — gemessen, nicht geschätzt.

## 7. Vorgelagert (vor der Modus-Begründung)

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md:223-224 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area berührt: `*` (Repo-Default). Sie erfüllt die Schwelle über
eigene Konventionsregeln (der Adaptions-Block), eigene Sensoren (`make gates`)
und eigene Artefaktklassen. `tools/harness/` ist **nicht** berührt — der Slice
ändert kein Skript, nur Doku-Form; die Pfad-Nähe von `harness/sensors/` <!-- d-check:ignore (entsteht mit diesem Slice) --> zu
`tools/harness/` ist Namensähnlichkeit, keine Berührung. Eine
Ausdifferenzierung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md:229-229 -->

> **Offene Beobachtungen sichten.**

Register durchgegangen (gemergter Stand, 34 Verzeichnisse). Vier Einträge
betreffen `*` und diesen Slice:

- [`large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)
  (2×, offen) — **gesichtet, und der Schutz hat nicht getragen.** Beim
  Schreiben dieses Blocks lautete die Antwort: „der Schnitt auf drei Träger
  ist die Antwort, ein dritter Beleg entsteht nicht". Der Schnitt beruhte auf
  einer falschen Messgröße und wurde korrigiert; der Slice liefert 23 Träger
  plus zwei zusätzliche Liefer-Punkte. **Die Beobachtung ist damit
  eingetreten** — nicht durch das Wachstum allein, sondern durch seine Folge:
  der zweite Review fand 19 Befunde, darunter mehrere, die direkt aus dem
  Wachstum stammen (nicht nachgezogene Zahlen, ein veralteter Vorprüfungs-Block,
  Deckungsnachweis in fünf neuen Dateien). Der Beleg wird bei der Closure
  eingetragen; der Zähler steht damit bei **3×** und die Schwelle ist erreicht.
- [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
  (6×, geplant) — unmittelbar im Gegenstand: Die auszulagernden Verträge sind
  genau die Stellen, an denen die Sektion Prüfungen beschreibt. Wer Prosa
  verschiebt, verschiebt auch ihre Behauptungen; jede übernommene Grenz-Aussage
  ist beim Umzug gegen das Target zu halten, nicht nur zu kopieren.
- [`begruendung-traegt-entscheidung-nicht`](../observations/BEO-ALL/begruendung-traegt-entscheidung-nicht/observation.md)
  (1×, offen, aus slice-202) — dieselbe Gefahr eine Ebene tiefer: Eine
  Träger-Datei besteht überwiegend aus Begründungen, und keine davon prüft ein
  Sensor.
- [`module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
  (1×, offen) — die Träger-Dateien tragen je einen Abschnitt *Grenze*; genau
  die Frage, die dieser Eintrag stellt.

Keiner der vier erreicht mit diesem Slice die Schwelle von 3× erstmalig.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-06 gelesen. `image-scan.yml` **grün**
(07:56:19Z). `upstream-drift.yml` **ROT** (05:23:44Z) — und hier greift die
benannte Grenze des Targets: es liest den **jüngsten** Lauf, nicht sein
**Alter**. Der Lauf datiert **vor** der Pin-Hebung aus slice-202; seine
Ursache — der Baseline-Rückstand — ist seither behoben, `make baseline-freshness`
meldet lokal beide Teile grün. Kein offener Befund, aber auch kein grüner
Nachtlauf, bis der nächste läuft.

## 8. Sub-Area-Modus-Begründung

**Modus:** die eine berührte Sub-Area `*` ist **GF** (Greenfield,
Repo-Default) — Doc führt, Code folgt. Kein Produkt-Code, keine
Reconciliation. **Konventionen-Dichte** hoch: die Form ist vollständig als
vendorte Vorlage verankert. **Phase-Reife** hoch: die Sektion existiert, ist
gewachsen und wird nur umgeformt. **Evidenz-/Diskrepanz-Risiko** liegt nicht
im Code-vs-Doku-Abstand, sondern in der **Aussagen-Treue beim Umzug** — die
vier gesichteten Beobachtungen oben benennen genau das.

## 9. Closure-Notiz (nach `done/`)

**Was hat funktioniert.** Die Mechanisierung. Die Regel *„eine Zelle, die zum
Absatz geworden ist, ist der Fund"* stand vorher nur im Kanon; jetzt hält sie
ein Sensor in `make gates`. Entscheidend war die Wahl der **Schwelle**: 200
Zeichen bilden die Ein-Satz-Regel ab, nicht den Bestand. Eine auf den Bestand
kalibrierte Schwelle (1420/1253) hätte ihn eingefroren und die Regel nie
verkörpert — der Auftraggeber hat genau das zurückgewiesen, und erst dadurch
wurde aus einer Aufräumaktion eine gehaltene Form. Der Sensor fand dabei drei
vorbestehende Defekte, die kein anderer sah: drei Zeilen, denen eine Spalte
**ganz** fehlte.

**Was ging anders als geplant.** Der Slice wuchs von drei auf 23 Träger plus
zwei Liefer-Punkte. Ursache war ein Messfehler: Der ursprüngliche Zuschnitt maß
die **Zeilen**länge inklusive Markdown-Spaltenpolsterung statt des Zellinhalts.
Nach echter Messung waren die drei ausgewählten Kandidaten Rang 1, 9 und 11 —
die Ränge 2 bis 8 wären stehen geblieben. Der erste Review fand das als F-1.

**Zwei Review-Runden, 16 MEDIUM.** Bemerkenswert ist ihre Verteilung: Der
**Kern** trug in beiden Runden — die Regel greift in beiden Tabellen und beiden
Richtungen, ist fail-closed gegen ihre eigene Abschaltung, und der Umzug der
größten Zelle (4997 Zeichen) ist Satz für Satz aussagen-treu. Blockiert hat
beide Male die **Schicht darüber**: Zahlen und Zuschreibungen über die eigene
Arbeit.

**Der schwerste Einzelbefund war eine Behauptung über einen Sensor.** Die
Commit-Botschaft von `aadb07e` schrieb der neuen Regel zwei Funde zu, die sie
nicht machen konnte — sie band nur die Vertrags-Spalte, eine Zeile ohne
Bindung-Zelle passierte still. Der Reviewer hat das per **Break-Test**
widerlegt, nicht per Argument. Behoben wurde es durch Bindung der Spalte, nicht
durch Rücknahme des Satzes: Die Regel findet die Klasse jetzt wirklich, belegt
durch einen eigenen Break-Test. **Die Botschaft von `aadb07e` bleibt trotzdem
falsch** — sie ist Lauf-Beleg und wird nicht umgeschrieben; sie stimmt erst ab
`c7da2c8`.

**Steering-Loop-Eintrag.** Der Schwellen-Übertritt ist
[`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)
(3×). Die drei Instanzen teilen die **Folge**, nicht den **Grund**: bei
slice-195 hätte eine Teilung den Zähler-Beleg zerrissen, bei slice-197 die
Werkzeug-Korrektur vervielfacht, hier sprengte **kein einzelner
Nachsteuerungs-Schritt** die Grenze — nur ihre Summe. Der Ausgang ist
**geplant**, nicht *verkörpert*, und das ist die inhaltliche Entscheidung
dieser Closure: Die Auflösung in genau dem Slice zu schreiben, der die
Schwelle ausgelöst hat, würde den Fehler wiederholen, den sie benennt. Sie ist
[slice-204](../in-progress/slice-204-slice-wachstum-aufloesung.md).
**Verkörpert wurde mit diesem Slice nichts** — der Eintrag ist gezählt und
geplant, nicht verkörpert.

**Beobachtungs-Register (`../observations/`):** vier Belege
`evidence/slice-203.md` — bei
[`large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)
(3×, Schwelle erreicht, Ausgang *geplant*),
[`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
(8×),
[`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
(7×) und
[`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(14×).

**Eine Beobachtung außerhalb des Registers, weil sie kein Repo-Defekt ist:**
Ein Review-Agent, der Break-Tests fährt, und ein paralleler `make gates`-Lauf
sehen **verschiedene Repos**. Ein Gate-Lauf meldete hier vier Befunde, die
zehn Minuten später verschwunden waren — echt zum Zeitpunkt der Messung,
Artefakt des Nebeneinanders. Wer beides parallel fährt, muss den roten Lauf
gegen `git status` halten, bevor er ihn glaubt.

**Folge-Slices:** [slice-204](../in-progress/slice-204-slice-wachstum-aufloesung.md)
(Auflösung für den wachsenden Slice) — ist eine Datei in `open/`.

**Risiken aus §5:** drei, jedes mit genau einem Ausgang — **zwei eingetreten**
(Deckungsnachweis in fünf Dateien; Nicht-Gate-Zuordnung in der benannten
gefährlicheren Richtung falsch), eines entfallen. Beide eingetretenen hatte der
Plan wörtlich vorhergesagt, und beide fand nicht der Autor, sondern der Review.

**Verifikation.** `make gates` Exit 0 (zehn Gates, 656 Dateien, 0 Befunde) ·
`make doc-check` grün auf jedem Zwischenstand · zwei eigene Break-Tests der
neuen Regel gefahren und zurückgerollt, dazu sieben des Reviewers · zwei
unabhängige Review-Runden (7+9 MEDIUM, alle behoben; 8 LOW behoben, 2 INFO
benannt), Reports unter
[`docs/reviews/`](../../../reviews/2026-09-06-slice-203-sensors-form-review-r2.md).

**Drei Paarungen:** nach dem `git mv` geprüft — Ergebnis im letzten
DoD-Häkchen.
