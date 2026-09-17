# CO-001: `vcs` überspringt still, solange das publizierte Image den Fix nicht trägt

**Status:** Aufgelöst (2026-09-17).

**Datum angelegt:** 2026-09-08. **Letzte Prüfung:** 2026-09-17.

**Betroffenes Gate:** `adr-check` (Modul `vcs`) — in beiden Modi, und darunter der
Bindepunkt der PR-/Push-CI.

**Geltungsbereich — aufgelöst.** Alle vier Ausprägungen (Begründung unten) sind
seit `v0.76.1` im publizierten Image behoben, empirisch gegen das **gezogene**
Image geprüft (nicht nur den Quellstand) — siehe Auflösungs-Trigger. Ein
Konsument, der `ghcr.io/pt9912/d-check:v0.76.0` oder früher pinnt
(typischerweise über das `DCHECK_IMAGE` eines per `--print-mk` erzeugten
`d-check.mk`), fährt weiterhin den blinden bzw. fehldiagnostizierenden Pfad und
löst den Carveout für sich durch einen Pin-Wechsel auf `v0.76.1` oder neuer
auf. **Dieses Repo ist nicht betroffen** — `make adr-check` trägt die
Prerequisite `build` und fährt `$(IMAGE):latest` aus dem lokalen Quellstand
([`Makefile`](../../../Makefile), `DCHECK_RUN`), der den Fix seit dem
Feature-Commit von slice-220 trägt.

**Zwei frühere Fassungen des Geltungsbereichs lagen falsch**, in beide
Richtungen: die erste schrieb „dieses Repo eingeschlossen", ohne den
Image-Pin nachzuschlagen; die zweite nahm dieses Repo ganz aus, als alle
Ausprägungen behoben wären, während die dritte/vierte noch offen waren.

**Folge-Slice:** [`slice-220`](../planning/done/slice-220-vcs-pfadmenge-statt-diff.md)
— der **Klassen**-Fix für die offene dritte Ausprägung; danach
[`slice-219`](../planning/in-progress/slice-219-release-loest-co-001.md), das Release,
das beides ausliefert. **Nicht** slice-218: der liefert zwei der drei
Ausprägungen und schließt davor; ein Carveout, dessen Folge-Slice vor ihm
schließt, hat faktisch keinen.

Regeln: Baseline-Regelwerk `modul-07-carveouts.md` §Ziel-Form: Carveout — ein
Carveout braucht immer einen Auflösungs-Trigger **und** einen Folge-Slice.

---

## Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-07-carveouts.md`
§Ziel-Form: Carveout — technische Begründung, keine
„noch nicht geschafft"-Aussagen.

**Der Gate ist blind, wenn die Objektdatenbank nicht kanonisch
gepackt ist — in beiden Modi.** go-git meldet ein **unlesbares** Objekt mit demselben Fehler wie
eine **im Tree fehlende** Datei (`object.ErrFileNotFound`). Der zweite Fall ist
harmlos und muss befundfrei bleiben — eine später angelegte Datei —, der erste
ist ein Umgebungsfehler. Vor dem Fix behandelte der Adapter beide gleich und
übersprang die Datei ohne Befund.

**Gemessen, mit echter Verletzung, in allen vier Ausprägungen** — dieselbe
Range, eine geänderte `**Status:** Accepted`-Datei, kanonisch je
`1 Befund(e)`/Exit 1:

| Was der unsichtbar benannte Pack verschluckt | wie die Datei ankommt | vor slice-218 | vor slice-220 | **heute** |
| --- | --- | --- | --- | --- |
| das BASE-**Blob** | `M` | `0 Befund(e)`, Exit 0 | zu (Exit 2) | **zu** (Exit 2) |
| das BASE-**Tree**, mit Pendant auf der Gegenseite | `A` | `0 Befund(e)`, Exit 0 | zu (Exit 2) | **zu** (Exit 2) |
| das BASE-**Tree**, **ohne** Pendant (Verzeichnis gelöscht/umbenannt) | **gar nicht** | `0 Befund(e)`, Exit 0 | offen | **zu** (Exit 2) |
| der **HEAD**-Tree | — | fälschlich `core-drift-vcs`, Exit 1 | fälschlich `core-drift-vcs`, Exit 1 | **zu** (Exit 2) |

`git diff` zeigt die Änderung in allen vier Fällen unverändert an. **Jede
Zeile fand eine andere Review-Runde** — die erste Fassung dieses Carveouts
kannte nur die erste und sagte den Abschluss trotzdem über die ganze Klasse
zu.

**Die dritte Zeile war der eigentliche Grund, warum dieser Carveout bestand.**
Sie entstand **vor** jeder Stelle, an der das Modul prüfen konnte: Der
Tree-Walker der Bibliothek machte aus einem nicht ladbaren Unterbaum ein
`io.EOF`, die Löschung erreichte die Änderungsliste also nie. Gemessen war
auch, dass der naheliegende Wachposten nicht trug — `tree.Files()` benutzt
denselben Walker und schwieg ebenso. Der Fix war deshalb eine
**Entwurfsänderung** (die geschützte Pfad-Menge direkt gegen beide Trees
auflösen, statt dem Diff zu vertrauen) —
[slice-220](../planning/done/slice-220-vcs-pfadmenge-statt-diff.md).

**Der vierte Ausgang war kein stiller, aber eine Fehldiagnose:** Verschluckte
der Pack den **HEAD**-Tree, meldete der Lauf `core-drift-vcs` *„gelöscht oder
umbenannt"* mit Exit 1 — laut, aber falsch; das Umgebungsproblem erschien als
Inhalts-Befund. slice-220s Entwurfsänderung schließt beide (dritte und
vierte) auf demselben Codepfad, ohne sie einzeln zu unterscheiden.

Der praktische Auslöser ist `git maintenance run --task=loose-objects`: Es
schreibt `loose-<Hash>.pack` und packt in Stapeln, hinterlässt also genau die
**partiellen** Packs, in denen der stille Pfad greift.

**Warum das ein Carveout ist und nicht nur ein behobener Fehler — zwei
Gründe, und der zweite ist der schwerere.** Erstens fährt ein Konsument ein
**gepinntes Image**: Bis zum nächsten Release zieht jeder, der `v0.75.0` pinnt,
auch die beiden behobenen Pfade weiterhin blind. Zweitens ist die dritte
Ausprägung **überhaupt nicht behoben** — sie liegt unterhalb der Stelle, an der
das Modul eingreifen kann, und niemand sieht sie aus seinem Repo heraus, weil
der Lauf grün meldet. Die Lücke ist real und datiert, nicht hypothetisch, und
sie steht in der Sensors-Tabelle, statt nur in einer Grenzen-Prosa zu wohnen.

**Der Rest-Risiko-Umfang ist klein und benennbar:** Es braucht (a) eine
git-Wartungsaufgabe, die unkanonisch benannte Packs schreibt, **und** (b) eine
echte Kern-Änderung an einer `Accepted`-ADR im selben Lauf — für die offene
dritte Ausprägung zusätzlich (c), dass das geschützte **Verzeichnis** gelöscht
oder umbenannt wird. Trifft nur (a) zu, ist der Lauf grün und auch korrekt
grün.

## Auflösungs-Trigger (erfüllt am 2026-09-17)

Regeln dieser Sektion: Baseline-Regelwerk `modul-07-carveouts.md`
§Ziel-Form: Carveout — konkret und prüfbar. „Wenn Zeit ist" ist kein Trigger.

**Der Klassen-Fix aus slice-220 ist drin UND das nächste Release ist
veröffentlicht** — `v0.76.1` (slice-219), getaggt und über GHCR publiziert
(Digest `sha256:1470ecdcaa686a5ef4513dee9b0ae522586f54b87d568b06fc6b5b2741b633b3`).
Ein Konsument löst den Carveout für sich auf, indem er sein `DCHECK_IMAGE`
darauf zieht. Geprüft mit:

```bash
docker run --rm --network none -v "$PWD:/repo:ro" \
  ghcr.io/pt9912/d-check@sha256:1470ecdcaa686a5ef4513dee9b0ae522586f54b87d568b06fc6b5b2741b633b3 \
  --enable vcs --disable links --range <base>..<head>
```

gegen ein Probe-Repo mit **partiell** unsichtbarem Pack, in **allen vier**
Ausprägungen der Tabelle oben — unsichtbares BASE-**Blob**, unsichtbares
BASE-**Tree** mit Pendant, unsichtbares BASE-**Tree ohne** Pendant und
unsichtbarer **HEAD**-Tree. **Gemessen** (2026-09-17, gegen das gezogene
Image, nicht nur den Quellstand): alle vier brechen mit **Exit 2** ab —
`nicht lesbares Objekt zu "adr-x.md" in "<ref>": file not found` (Blob) bzw.
`Range-Basis/-Spitze "<ref>" nicht auflösbar: nicht vollständig lesbarer Tree
zu "<ref>": nicht lesbarer Unterbaum "sub": object not found` (alle drei
Tree-Fälle).

**Alle vier wurden gefahren, und der Grund steht in der Geschichte dieses
Carveouts:** Die erste Fassung kannte nur die Blob-Ausprägung, die zweite
zwei — und hätte den Trigger jedes Mal für erfüllt gehalten, während die
nächste noch offen war. **Denselben Fehler zweimal zu machen war der Grund
für die Zahl in dieser Zeile; die vierte kam erst durch slice-220 selbst
dazu, nicht durch eine weitere Wiederholung desselben Fehlers.**

## Geltungs-Konfiguration

Der Carveout ist **nicht** im Werkzeug konfiguriert — es gibt keinen Schalter,
der ihn ausdrückt. Er beschreibt einen Zustand des gepinnten Images.

| Datei | Zeile/Section | Wert |
|---|---|---|
| Repo des Konsumenten: `d-check.mk` | `DCHECK_IMAGE` | der Pin, der den blinden Stand trägt |
| [`harness/README.md`](../../../harness/README.md) | §Sensors, Zeile `make adr-check` | `CO-001` in der Bindung-Spalte |

## Verifikation (nach Auflösung)

- [x] Gate ist für den Geltungsbereich aktiviert (Gate-Konfiguration aktualisiert).
- [x] `make gates` grün ohne Ausnahme.
- [x] Datei wird nach `docs/plan/carveouts/done/` bewegt (reiner `git mv`). <!-- d-check:ignore (done/ entsteht erst bei erster Carveout-Auflösung) -->
- [x] Folge-Slice geschlossen oder explizit dokumentiert.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-08 | Angelegt | [slice-218](../planning/done/slice-218-go-git-pack-namenskonvention.md) |
| 2026-09-17 | Aufgelöst — Klassen-Fix (slice-220) im Release `v0.76.1` (`slice-219`, zum Schreibzeitpunkt dieser Zeile noch nicht selbst geschlossen) publiziert, alle vier Ausprägungen gegen das gezogene Image gemessen | [slice-220](../planning/done/slice-220-vcs-pfadmenge-statt-diff.md) |
