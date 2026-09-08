# CO-001: `vcs` überspringt still, solange das publizierte Image den Fix nicht trägt

**Status:** Aktiv.

**Datum angelegt:** 2026-09-08. **Letzte Prüfung:** 2026-09-08.

**Betroffenes Gate:** `adr-check` (Modul `vcs`) — in beiden Modi, und darunter der
Bindepunkt der PR-/Push-CI.

**Geltungsbereich:** **Konsumenten**, die ein veröffentlichtes Image pinnen —
`ghcr.io/pt9912/d-check:v0.75.0` oder früher, typischerweise über das
`DCHECK_IMAGE` eines per `--print-mk` erzeugten `d-check.mk`.

**Dieses Repo ist ausdrücklich nicht betroffen** — und das ist gemessen, nicht
angenommen: `make adr-check` trägt die Prerequisite `build` und fährt
`$(IMAGE):latest` aus dem lokalen Quellstand
([`Makefile`](../../../Makefile), `DCHECK_RUN`). Der Fix wirkt hier also ab
dem Commit, auch in der PR-/Push-CI. **Eine frühere Fassung dieses Carveouts
behauptete das Gegenteil** („dieses Repo eingeschlossen") — sie hatte den
Image-Pin nicht nachgeschlagen.

**Folge-Slice:** [`slice-219`](../planning/open/slice-219-release-loest-co-001.md) — das Release, das den Fix ausliefert. **Nicht** slice-218: der liefert den Fix und schließt davor; ein Carveout, dessen Folge-Slice vor ihm schließt, hat faktisch keinen.
Der Carveout deckt die Strecke zwischen Fix und Auslieferung.

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

| Was der unsichtbar benannte Pack verschluckt | wie die Datei ankommt | vor dem Fix |
| --- | --- | --- |
| das BASE-**Blob** | `M` — der Core-Vergleich überspringt sie | `0 Befund(e)`, Exit 0 |
| das BASE-**Tree** des Verzeichnisses | `A` — als wäre sie neu | `0 Befund(e)`, Exit 0 |

`git diff` zeigt die Änderung in beiden Fällen unverändert an. **Die zweite
Zeile fand erst die zweite Review-Runde** — die erste Fassung dieses Carveouts
kannte nur die Blob-Hälfte und sagte den Abschluss trotzdem über die ganze
Klasse zu.

Der praktische Auslöser ist `git maintenance run --task=loose-objects`: Es
schreibt `loose-<Hash>.pack` und packt in Stapeln, hinterlässt also genau die
**partiellen** Packs, in denen der stille Pfad greift.

**Warum das ein Carveout ist und nicht nur ein behobener Fehler:** Der Fix
liegt im Quellstand, aber ein Konsument fährt ein **gepinntes Image**. Bis
zum nächsten Release zieht jeder, der `v0.75.0` pinnt, weiterhin den blinden
Pfad — und er kann es aus seinem Repo heraus nicht sehen, weil der Lauf grün
meldet. Die Lücke ist also real und datiert, nicht hypothetisch, und sie steht
in der Sensors-Tabelle, statt nur in einer Grenzen-Prosa zu wohnen.

**Der Rest-Risiko-Umfang ist klein und benennbar:** Es braucht (a) eine
git-Wartungsaufgabe, die unkanonisch benannte Packs schreibt, **und** (b) eine
echte Core-Änderung an einer `Accepted`-ADR im selben Lauf. Trifft nur (a) zu,
ist der Lauf grün und auch korrekt grün.

## Auflösungs-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-07-carveouts.md`
§Ziel-Form: Carveout — konkret und prüfbar. „Wenn Zeit ist" ist kein Trigger.

**Das nächste Release ist veröffentlicht** — also der Tag nach `v0.75.0`, mit
dem der Fix aus slice-218 das publizierte Image erreicht; ein Konsument löst
den Carveout für sich auf, indem er sein `DCHECK_IMAGE` darauf zieht. Prüfbar
mit:

```bash
docker run --rm --network none -v "$PWD:/repo:ro" <neues-image> \
  --enable vcs --disable links --range <base>..<head>
```

gegen ein Probe-Repo mit **partiell** unsichtbarem Pack, und zwar in **beiden**
Ausprägungen — unsichtbares BASE-**Blob** (die Datei kommt als `M` an) und
unsichtbares BASE-**Tree** (sie kommt als `A` an). Vor der Auflösung meldet der
Lauf je `0 Befund(e)` mit Exit 0, danach `nicht lesbares Objekt zu …` bzw.
`nicht lesbarer Tree-Eintrag zu …` mit Exit 2. **Beide sind zu fahren:** Die
erste Fassung dieses Carveouts kannte nur die Blob-Hälfte und hätte den
Trigger für erfüllt gehalten, während die andere Hälfte noch offen war.

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
