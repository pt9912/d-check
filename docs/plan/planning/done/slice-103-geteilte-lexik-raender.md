# Slice slice-103: Dieselbe Klasse, andere Lexiken — Absatzbildung, Anker-Auflösung, git-Revisionen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** welle-74.

**Bezug:** [`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in),
[`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in).
Vorgeschichte: [ADR-0050](../../adr/0050-fence-unclosed-in-spans.md) und die drei
Review-Runden an slice-101.

**Autor:** pt9912. **Datum:** 2026-08-10.

---

## 1. Ziel

Drei Befunde aus der dritten Review-Runde an slice-101 abarbeiten, die dort
ausdrücklich **nicht** hingehörten: sie betreffen andere Module mit eigenem
Vertrag, und sie sind älter als der Fence-Wächter. Gemeinsam ist ihnen die
Klasse, die slice-101 für die **Fence**-Lexik geschlossen hat.

## 2. Die Klasse

*Eine geteilte Lexik driftet an den Rändern, weil jeder Konsument sie selbst
vorbereitet.* Bei den Fences waren es fünf Stellen, die dieselbe Zeile
verschieden trimmten, und zwei Automaten, die dieselbe Frage verschieden
beantworteten. Jede für sich vertretbar — zusammen ein stiller Grün-Pfad.

slice-101 hat das für die Fence-Lexik geschlossen: ein geteiltes Trimm-Prädikat,
beide Schluss-Lesarten geteilt und ausgewertet, je Konsument eine Assertion
gegen Wieder-Divergenz. Die drei Befunde hier zeigen dieselbe Bauform in
**anderen** Lexiken.

## 3. Die drei Befunde

1. **Absatzbildung in `citations`** (Review R-2, HIGH). Das Modul gruppiert
   Absätze selbst, und ein Fence ist dort **keine** Grenze. Dieselbe Datei
   liefert Exit 0 statt des zugesagten fail-closed Exit 2, nur weil zwischen
   Direktive und Zitat ein Code-Block statt einer Leerzeile steht.
2. **Anker-Auflösung in `headingSection`** (Review R-3, MEDIUM). Sie beantwortet
   die Anker-Frage roh: ein HTML-Anker **innerhalb** eines Fence erfüllt die
   fail-closed-Bedingung von `versions.current-from`, während `anchors` im
   selben Lauf `anchor-missing` meldet. Zwei Module, dieselbe Frage, zwei
   Antworten.
3. **git-Revisionen als dritte unerreichbare Eingabe-Achse** (Review R-7, LOW).
   `vcs` rechnet die fence-empfindliche Section-Maske auf git-Blobs, die kein
   scannendes Modul je sieht — ein Fixture belegt ein falsches
   `core-drift-vcs`. Das ist die **dritte** Wiederholung von „Modul-Grenze nur
   auf der Quell-Achse gedacht" (slice-101 Review F-3, N-2, R-7).

## 4. Abnahme-Punkte

1. **Erst messen, dann entscheiden** — wie bei slice-101. Wie viele Dokumente im
   Ökosystem lösen die drei Fälle heute aus? Die Fence-Messung dort drehte die
   Entscheidung (776 Dateien, null Vorkommen ⇒ latent statt aktiv); ohne Zahl
   ist die Reichweite jeder Variante Spekulation.
   — **Ausgang:** gemessen (§4a). **Alle drei Fälle sind latent**, null
   Vorkommen über drei Repos. Die Reparatur ändert damit für keinen heutigen
   Konsumenten ein Ergebnis; die Richtung ist trotzdem zu nennen, weil sie
   **findet mehr** lautet.
2. **Ein Slice oder drei?** Die Klasse ist gemeinsam, die Verträge sind es
   nicht. Zu entscheiden nach der Messung.
   — **Ausgang: ein Slice.** Die Messung nimmt den Grund für einen Schnitt
   weg: kein Fall ist dringlicher als die anderen, keiner erreicht heute einen
   Konsumenten, und die beiden Reparaturen sind **je eine Zeile Lexik** plus
   ihre Assertionen. Fall 3 liefert **keinen** Code, sondern eine benannte
   Grenze (Abnahme-Punkt 3).
3. **Reparieren oder melden?** slice-101 hat den Zustand gemeldet statt die
   Paarung zu reparieren. Ob das hier trägt, ist offen: bei `citations` und
   `headingSection` geht es um eine **falsche Antwort**, nicht um einen
   unentscheidbaren Zustand.
   — **Ausgang: reparieren (Fall 1 und 2), benennen (Fall 3).** Die
   Unterscheidung ist nicht Aufwand, sondern **Erreichbarkeit**: Fall 1 und 2
   sind Fragen, die das Produkt an einer anderen Stelle bereits **richtig**
   beantwortet — dort ist die Reparatur die Übernahme der vorhandenen Antwort.
   Fall 3 dagegen betrifft eine Eingabe, die kein scannendes Modul je sieht
   (git-Blob); dort gibt es keine vorhandene Antwort zu übernehmen, sondern nur
   einen neuen Mechanismus zu bauen — für einen Fall, der in **152**
   Revisions-Blobs nie eingetreten ist.
4. **Die dritte Wiederholung.** „Modul-Grenze nur auf der Quell-Achse“ hat mit
   R-7 die Schwelle des Beobachtungs-Registers erreicht. Zu entscheiden ist,
   welche Form sie verkörpert — die Register-Regel sagt: verkörpern statt
   weiterzählen.
   — **Ausgang: erledigt, außerhalb dieses Slice.** Die Klasse steht als
   **BEO-004** im Register und ist in der welle-73-Closure verkörpert worden —
   als MEDIUM-Anker im Reviewer-Skill („welche Eingaben liest dieses Modul, die
   es nicht scannt?“). Damit **gehört Fall 3 gar nicht zu BEO-003:** er ist
   keine geteilte Lexik, die driftet, sondern eine unerreichbare Eingabe-Achse.
   Die Vermischung stand im Slice-Text und löst sich hier auf.

## 4a. Messung: wie oft treten die drei Fälle im Bestand auf?

Methode wie in slice-101: der **reale** Bestand dreier Repos (`d-check`,
`a-check`, `ai-harness-course`), und wo eine Aussage vom Fence-Verhalten
abhängt, entscheidet **das Produkt**, nicht eine Nachrechnung.

| Fall | Grundmenge | Vorkommen |
|---|---|---|
| 1 — Fence zwischen Direktive und Zitat | 18 Marker-Vorkommen (alle in `d-check`) | **0** — und die Grundmenge ist schwächer, als die Zahl aussieht: **keines** der 18 ist eine produktive Direktive, alle stehen in Inline-Code oder Prosa (Handbuch, READMEs, Spec, Reports). Das Ergebnis „null betroffene Fälle“ steht, die Zahl 18 belegt es nicht |
| 2 — HTML-Anker **innerhalb** eines Fence | 94 HTML-Anker, davon 93 außerhalb | **1** (`vX.Y.Z`, ein Formbeispiel im Handbuch) — und **kein** Konsument fragt ihn: `versions.current-from` zeigt auf einen Überschriften-Anker, die drei `dpin`-Fragmente auf Slugs |
| 3 — unbalancierter Fence in einer immutablen Revision | 53 ADR-Pfade, **152** deduplizierte Revisions-Blobs | **0** |

**Fall 3 ist der Beleg für die Methode.** Eine naive Paritätszählung meldete
**zwei** verdächtige Blobs; beide liegen in
[ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md), der die
Fence-Lexik selbst dokumentiert und dafür **geschachtelte** Beispiel-Fences
zeigt. Das Produkt über dieselben Blobs gefahren (`--enable spans`): **0
Befunde** — eine Öffnung wird nur von einer mindestens gleich langen Folge
geschlossen. Die Nachrechnung irrte, nicht der Bestand.

### Ist die Klasse damit geschlossen?

Die Leitfrage aus slice-101 lautet nicht „sind die Befunde abgehakt“, sondern
„ist die Klasse geschlossen“. Dafür sind **alle** Stellen aufgezählt worden, die
eine Lexik-Frage beantworten:

- **Fence-Erkennung:** **vier** Dateien führen eine eigene Fence-Zustands-Schleife
  — `markdown.go` (die Lexik selbst), `trace_table.go`, `sections.go` (zweimal,
  für Abschnitts-Kopf und -Ende, konsumiert von `planning` und `structure`) und
  `spans.go` (der Wächter). Alle vier speisen sich aus denselben geteilten
  Prädikaten `TrimFenceIndent` / `FenceRun` / `FenceCloses`, die **Antwort**
  stimmt also überein. Die erste Fassung dieser Aufzählung nannte zwei von vier
  — sie war als Beleg für „die Klasse ist geschlossen" gedacht und hätte einem
  Leser `sections.go` verschwiegen.
- **Rohe Zeilen-Spaltung** an 22 Stellen — die überwiegende Mehrheit schneidet
  nur nach Zeilennummer und holt die **Antwort** aus der geteilten Lexik
  (`matrix.SelectSections`/`SectionMask`/`HeadingSections` etwa aus
  `extractHeadingLines`). Vertraglich roh lesen `versions` (Pins auch in Fences,
  [ADR-0019](../../adr/0019-versions-pin-fence-ausnahme.md)), `pins`
  ([ADR-0020](../../adr/0020-content-pin-fence-ausnahme.md)) und `immutable` —
  das ist keine Drift, sondern die zugesagte Antwort.
- **Übrig bleiben genau zwei** Stellen, die eine Lexik-Frage **roh**
  beantworten: die Absatzbildung in `citations` und der HTML-Anker-Zweig in
  `headingSection`. Beide sind Fall 1 und 2.
- **Der erste Anlauf hat die zweite nur halb repariert** (unabhängiger Review):
  übernommen war die **Zeilen-Menge**, nicht die **Erkennung** — vier
  Zeichenfolgen galten weiter für `versions`/`pins` als Anker und für `anchors`
  nicht (Anker in Inline-Code, `data-id`, `name` an beliebigem Element,
  anker-förmige Prosa ohne Tag). Gemessen trug die reparierte Achse **1**
  Vorkommen, die stehengebliebene **40**. Die Lehre gehört zur Klasse: eine
  geteilte Antwort ist nicht halb übernehmbar, und „dieselbe Grundmenge“ ist
  noch keine „dieselbe Antwort“.

**Die Aufzählung hat ein zweites Mal nicht getragen.** Die bestätigende
Re-Review fand zwei weitere Stellen, beide im Modul `planning`: der
Aktiv-Status-Guard zählt die kanonische Überschrift **roh** und liest den
Ruhe-Marker **roh**. Die Folge war zweimal ein **Falsch-Rot** — eine Roadmap,
die ihren eigenen Abschnitt in einem Beispiel-Block zeigt, galt als mehrdeutig,
und eine Raute-Zeile in einem Beispiel-Block beendete den Aktiv-Block vorzeitig.
Beide sind repariert.

**Eine dritte Runde hat drei weitere Module gefunden**, die in keinem
Vorgänger-Report vorkamen: `vcs` (Status-Zeile und Immutabilitäts-Entscheidung
roh — **stilles Grün im Immutabilitäts-Gate**), `planning` ein zweites Mal (die
H2-Grenze war ein roher Präfix-Vergleich, obwohl derselbe Commit die
Fence-Hälfte gerade geheilt hatte) und `targets` (Tabellenzeilen roh). Damit ist
auch die Out-of-Scope-Begründung dieser Welle widerlegt: die Tabellen-Lexik ist
**kein neuer** Rand, sie driftet bereits.

**Die Aufzählung ist als Methode endgültig widerlegt.** Dreimal hintereinander
hat eine Liste, die „vollständig“ hieß, Stellen nicht gekannt — und jedes Mal
fand ein Review sie in Minuten. Die
Antwort ist deshalb kein dritter Listen-Anlauf, sondern ein **Kopplungs-Test**:
`TestAnkerFrageHatEineAntwort` fährt zehn Anker-Schreibweisen durch **beide**
Konsumenten und schlägt fehl, sobald sie verschieden antworten. Er prüft nicht,
ob eine bestimmte Funktion aufgerufen wird, sondern ob die Antworten
übereinstimmen — eine Liste kann eine Stelle vergessen, eine Kopplung nicht.
Das ist die Verkörperung von **BEO-003** (Zähler 3).

## 5. Definition of Done

- [x] Bestandsmessung für alle drei Fälle über die drei Repos (§4a): 0 · 1 · 0.
- [x] Abnahme-Punkte 1–4 entschieden (§4), Vertragsanpassung geliefert
      (Lastenheft 0.58.0, Spezifikation, [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)).
- [x] Je Fall ein mutations-echter Test; die Gegenprobe über eine **Dateikopie**,
      nicht über `git checkout`. **Fünf** Assertionen, jede fängt ihren Rückbau:
      zwei an `citations`, je eine an `versions` und `pins`, und der **bestehende**
      Absatz-Test hängt am selben Prädikat — das ist der Beleg, dass die Lexik
      wirklich geteilt ist und nicht nur gleich aussieht.
- [x] `make gates` grün (375/0); **SemVer: Minor** — die Änderung **findet
      mehr** (fail-closed statt Zufalls-Paarung, Anker im Fence löst nicht mehr
      auf) **und weniger** an mehreren — die Aufzählung der Richtungen steht in
      Release-Notiz und ADR und ist ausdrücklich **offen** formuliert, weil sie
      in drei Review-Runden dreimal unvollständig war. Am eigenen Bestand ist
      **kein** Fall betroffen; ohne die reparierten Konsumenten-Fälle ist der
      Befundsatz byte-identisch.

## 6. Risiken / offene Punkte

- **Drei Verträge in einem Slice** könnten den Schnitt sprengen. Abnahme-Punkt 2
  entscheidet das nach der Messung, nicht vorher.
- **Der `vcs`-Fall ist LOW und teuer.** Die git-Achse für ein scannendes Modul
  erreichbar zu machen, ist ein anderer Umbau als ein geteiltes Prädikat.
  Möglich, dass hier nur die Grenze benannt gehört.

## 7. Trigger

**Start** (`open` → `next`): nach der Closure von slice-101 und wenn ein
WIP-Slot frei ist. Keine Kopplung an einen Release.

## 8. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** bei der Planung erneut lesen — die
  Ziel-Achsen-Klasse aus Abnahme-Punkt 4 steht dann voraussichtlich im Register.

## 9. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — pro Fall wird zuerst die Zusage formuliert
(welche Antwort gilt?), dann geliefert.

## 10. Closure-Notiz (nach `done/`)

**Geliefert:** sechs Module beantworten ihre Lexik-Fragen jetzt gleich —
`citations` (Absatzgrenze), `versions` und `pins` (Anker-Auflösung vollständig:
Fence, Inline-Code, Tag-Kontext, Duplikat-Slug, Prozent-Dekodierung,
Groß-/Kleinschreibung), `planning` (Überschrift, Marker, Block-Grenze), `vcs`
(Status-Zeile, `immutable-when`) und `targets` (Tabellenzeilen). Ausgeliefert als
**v0.58.0**. Für die git-Revisions-Achse liefert der Slice bewusst **keinen**
Code, sondern eine benannte Grenze mit beobachtbarem Trigger.

**Der Slice hat dreimal geglaubt, am Ende zu sein.** Nach der Messung standen
zwei betroffene Stellen fest, nach dem ersten Review vier, nach dem zweiten
sechs — und jede Zwischenbilanz war als vollständig formuliert. Zwei der zuletzt
gefundenen trugen ein **stilles Grün in einem Gate**: eine reale Core-Änderung
einer `Accepted`-ADR passierte das Immutabilitäts-Gate ohne Ausgabe, und
`planning-drift` entfiel, wenn der Aktiv-Block an einer eingerückten H2 endete.

**Die teuerste Lehre ist methodisch.** Eine Aufzählung belegt keinen
Klassen-Abschluss: dreimal hat eine Liste, die „vollständig“ hieß, Stellen nicht
gekannt, und jedes Mal fand ein Review sie in Minuten. Was trägt, ist eine
**Kopplung** — ein Test, der alle Konsumenten derselben Frage dieselbe Eingabe
beantworten lässt — oder eine **erschöpfende Prüfung mit Negativbefund je
Kandidat**. Beides liegt jetzt vor: `TestAnkerFrageHatEineAntwort` und die
Modul-für-Modul-Liste der dritten Review-Runde.

**Auch die Reparatur war zweimal halb.** Die Anker-Frage wurde erst zur
Fence-Hälfte vereinheitlicht, dann zur HTML-Hälfte, und erst im dritten Anlauf
auch in der Slug-Hälfte. Eine geteilte Antwort ist nicht halb übernehmbar, und
„dieselbe Grundmenge“ ist noch keine „dieselbe Antwort“.

**Offen bleibt benannt:** die Klasse ist dort nicht mechanisch geschlossen, wo
noch kein Kopplungs-Test existiert — die Absatz-, Überschriften- und
Tabellen-Achse tragen bis heute nur Einzel-Assertionen. Die Register-Bewegungen
stehen in der [Wellen-Ergebnisnotiz](welle-74-results.md).
