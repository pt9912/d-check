# CO-001: `vcs` überspringt still, solange das publizierte Image den Fix nicht trägt

**Status:** Aktiv.

**Datum angelegt:** 2026-09-08. **Letzte Prüfung:** 2026-09-08.

**Betroffenes Gate:** `adr-check` (Modul `vcs`) — in beiden Modi, und darunter der
Bindepunkt der PR-/Push-CI.

**Geltungsbereich — zwei ungleiche Hälften, und die zweite trifft auch uns:**

- **Zwei der drei Ausprägungen sind im Quellstand behoben** (slice-218). Für
  sie gilt der Carveout nur **Konsumenten**, die ein veröffentlichtes Image
  pinnen (`ghcr.io/pt9912/d-check:v0.75.0` oder früher, typischerweise über das
  `DCHECK_IMAGE` eines per `--print-mk` erzeugten `d-check.mk`). **Dieses Repo
  ist da nicht betroffen** — gemessen, nicht angenommen: `make adr-check` trägt
  die Prerequisite `build` und fährt `$(IMAGE):latest` aus dem lokalen
  Quellstand ([`Makefile`](../../../Makefile), `DCHECK_RUN`).
- **Die dritte Ausprägung ist offen** und trifft damit **jeden**, dieses Repo
  eingeschlossen — Quellstand wie gepinntes Image.

**Zwei frühere Fassungen dieses Feldes lagen falsch**, in beide Richtungen: die
erste schrieb „dieses Repo eingeschlossen", ohne den Image-Pin nachzuschlagen;
die zweite nahm dieses Repo ganz aus, als alle Ausprägungen behoben wären. Die
Unterscheidung oben ist die gemessene.

**Folge-Slice:** [`slice-220`](../planning/open/slice-220-vcs-pfadmenge-statt-diff.md)
— der **Klassen**-Fix für die offene dritte Ausprägung; danach
[`slice-219`](../planning/open/slice-219-release-loest-co-001.md), das Release,
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

**Gemessen, mit echter Verletzung, in zwei Ausprägungen** — dieselbe Range,
eine geänderte `**Status:** Accepted`-Datei, kanonisch je
`1 Befund(e)`/Exit 1:

| Was der unsichtbar benannte Pack verschluckt | wie die Datei ankommt | vor slice-218 | **heute** |
| --- | --- | --- | --- |
| das BASE-**Blob** | `M` | `0 Befund(e)`, Exit 0 | **zu** (Exit 2) |
| das BASE-**Tree**, mit Pendant auf der Gegenseite | `A` | `0 Befund(e)`, Exit 0 | **zu** (Exit 2) |
| das BASE-**Tree**, **ohne** Pendant (Verzeichnis gelöscht/umbenannt) | **gar nicht** | `0 Befund(e)`, Exit 0 | **weiterhin offen** |

`git diff` zeigt die Änderung in allen drei Fällen unverändert an. **Jede Zeile
fand eine andere Review-Runde** — die erste Fassung dieses Carveouts kannte nur
die erste und sagte den Abschluss trotzdem über die ganze Klasse zu.

**Die dritte Zeile ist der eigentliche Grund, warum dieser Carveout bleibt.**
Sie entsteht **vor** jeder Stelle, an der das Modul prüfen könnte: Der
Tree-Walker der Bibliothek macht aus einem nicht ladbaren Unterbaum ein
`io.EOF`, die Löschung erreicht die Änderungsliste also nie. Gemessen ist auch,
dass der naheliegende Wachposten nicht trägt — `tree.Files()` benutzt denselben
Walker und schweigt ebenso. Ein Fix ist deshalb eine **Entwurfsänderung** (die
geschützte Pfad-Menge direkt gegen beide Trees auflösen, statt dem Diff zu
vertrauen) und liegt bei
[slice-220](../planning/open/slice-220-vcs-pfadmenge-statt-diff.md).

**Ein vierter Ausgang gehört daneben, obwohl er kein stiller ist:** Verschluckt
der Pack den **HEAD**-Tree, meldet der Lauf `core-drift-vcs` *„gelöscht oder
umbenannt"* mit Exit 1 — laut, aber eine **Fehldiagnose**; das Umgebungsproblem
erscheint als Inhalts-Befund.

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

## Auflösungs-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-07-carveouts.md`
§Ziel-Form: Carveout — konkret und prüfbar. „Wenn Zeit ist" ist kein Trigger.

**Der Klassen-Fix aus slice-220 ist drin UND das nächste Release ist veröffentlicht** — also der Tag nach `v0.75.0`, mit
dem der Fix aus slice-218 das publizierte Image erreicht; ein Konsument löst
den Carveout für sich auf, indem er sein `DCHECK_IMAGE` darauf zieht. Prüfbar
mit:

```bash
docker run --rm --network none -v "$PWD:/repo:ro" <neues-image> \
  --enable vcs --disable links --range <base>..<head>
```

gegen ein Probe-Repo mit **partiell** unsichtbarem Pack, in **allen drei**
Ausprägungen der Tabelle oben — unsichtbares BASE-**Blob** (`M`), unsichtbares
BASE-**Tree** mit Pendant (`A`) und unsichtbares BASE-**Tree ohne** Pendant
(Verzeichnis gelöscht, erscheint gar nicht). Vor der Auflösung meldet der Lauf
je `0 Befund(e)` mit Exit 0; danach müssen alle drei mit **Exit 2** abbrechen.

**Alle drei sind zu fahren, und der Grund steht in der Geschichte dieses
Carveouts:** Die erste Fassung kannte nur die Blob-Ausprägung, die zweite zwei
— und hätte den Trigger jedes Mal für erfüllt gehalten, während die nächste
noch offen war. **Dreimal denselben Fehler zu machen ist der Grund für die
Zahl in dieser Zeile.**

## Geltungs-Konfiguration

Der Carveout ist **nicht** im Werkzeug konfiguriert — es gibt keinen Schalter,
der ihn ausdrückt. Er beschreibt einen Zustand des gepinnten Images.

| Datei | Zeile/Section | Wert |
|---|---|---|
| Repo des Konsumenten: `d-check.mk` | `DCHECK_IMAGE` | der Pin, der den blinden Stand trägt |
| [`harness/README.md`](../../../harness/README.md) | §Sensors, Zeile `make adr-check` | `CO-001` in der Bindung-Spalte |

## Verifikation (nach Auflösung)

- [ ] Gate ist für den Geltungsbereich aktiviert (Gate-Konfiguration aktualisiert).
- [ ] `make gates` grün ohne Ausnahme.
- [ ] Datei wird nach `docs/plan/carveouts/done/` bewegt (reiner `git mv`). <!-- d-check:ignore (done/ entsteht erst bei erster Carveout-Auflösung) -->
- [ ] Folge-Slice geschlossen oder explizit dokumentiert.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-08 | Angelegt | [slice-218](../planning/in-progress/slice-218-go-git-pack-namenskonvention.md) |
