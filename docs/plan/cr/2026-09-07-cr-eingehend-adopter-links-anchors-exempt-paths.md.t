# Eingehender Change Request — `links` und `anchors` brauchen `exempt-paths`

**Absender:** Adopter · **Eingegangen:** 2026-09-07
*(Der Wortlaut nennt den Absender nicht; er ist über den Vorgang bekannt und
wird hier nicht erfunden.)*
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links),
[`DC-FA-ANCH-001`](../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors),
[`DC-FA-REF-001`](../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
**Stand:** **entschieden am 2026-09-08 — nicht umgesetzt**; alle drei Argumente
beantwortet (§Entscheid), die tragende Messung steht unter §Messung. Die
Rückfrage an den Absender ist offen und im Entscheid benannt.

**Ablage-Hinweis.** Ein **eingehender** CR ist die dritte Klasse neben
[`MR-035`](../../../harness/conventions.md#mr-035) (ausgehend) und
[`MR-036`](../../../harness/conventions.md#mr-036) (Antworten darauf); der
Kanon führt sie als *externen Vorgang* und ausdrücklich als „bewusst kein
Harness-Konstrukt". Die Datei liegt hier aus demselben Grund wie die beiden
Vorgänger: Der erste Konsumenten-CR dieses Repos ging verloren, und mit ihm die
Frage, was genau gebeten und mit welcher Begründung entschieden wurde.

---

## Wortlaut (unverändert übernommen)

> # Change Request an d-check
>
> `links` und `anchors` brauchen `exempt-paths` — wie sechs andere Module es führen
>
> ## Der Anlass
>
> Ein Repo, das eine externe Baseline committet vendored hält, legt sie unter einem
> tag-gescopten Pfad ab (`.harness/baseline/<tag>/…`). Beim Versions-Sprung fällt <!-- d-check:ignore (zitierter Fremd-Wortlaut, kein Verweis) -->
> das alte Verzeichnis, und jede Adresse darauf stirbt.
>
> Das ist beherrschbar, solange die verweisenden Artefakte lebend sind — der Sprung
> zieht sie nach. Es ist nicht beherrschbar, wenn sie eingefroren sind:
> Review-Reports, geschlossene Slice-Pläne, `Accepted`-ADRs. Diese Artefakte halten
> eine Messung zu ihrem Datum fest und dürfen per Regel nicht mehr angefasst werden.
>
> Nach unserem Sprung `v6.0.0` → `v6.5.0`, gemessen am Ist-Stand:
>
> ```
> make docs-check | grep 'target-missing' | grep -cE '^docs/(reviews|plan/planning/done)/'    # 35
> make docs-check | grep 'target-missing' | awk -F'\t' '{print $1}' | sed 's/:[0-9]*$//' \
>   | grep -E '^docs/(reviews|plan/planning/done)/' | sort -u | wc -l                          # 15
> ```
>
> 35 Befunde über 15 eingefrorene Dateien. Sie sind weder reparierbar
> (Immutabilität) noch ignorierbar (kein Knopf) — das Gate bleibt rot, ohne dass ein
> Defekt vorliegt.
>
> ## Die Bitte
>
> `links` und `anchors` bekommen `exempt-paths` mit derselben Semantik wie anderswo:
> datei-weit, Glob-Liste, die genannten Dateien fallen aus dem Prüfbereich dieses
> Moduls.
>
> ```yaml
> links:
>   exempt-paths: ["docs/reviews/**", "docs/plan/planning/done/**"]
> anchors:
>   exempt-paths: ["docs/reviews/**"]
> ```
>
> ## Warum das kein Sonderwunsch ist
>
> Der Knopf ist im Werkzeug etabliert. Sechs Module führen ihn:
>
> ```
> grep -rhoE '\b[a-z-]+\.exempt-paths' spec/*.md *.md | sed 's/\.exempt-paths//' | sort -u
> # codepaths  diagrams  matrix  reviews  versions  workflows
> ```
>
> `links` und `anchors` haben gar keine Options-Sektion — bei uns gemessen:
> `grep -cE '^(links|anchors):' .d-check.yml` → 0.
>
> Und dieselbe Frage ist auf der anderen Achse längst beantwortet — durch
> `codepaths.exempt-paths`, mit exakt dieser Begründung in unserer Config:
>
> > Zeitdokumente `docs/reviews/**` frieren den Stand ihres Review-Laufs ein;
> > Lifecycle-Pfade (`next/`→`in-progress/`→`done/`) darin veralten per Definition. <!-- d-check:ignore (zitierter Fremd-Wortlaut, kein Verweis) -->
>
> Für einen Inline-Code-Pfad in einem eingefrorenen Report ist die Antwort also
> „darf veralten". Für einen Markdown-Link auf dasselbe Ziel, in derselben Datei,
> gibt es keine. Der Unterschied ist heute kein Beschluss, sondern die Abwesenheit
> eines Knopfes.
>
> ## Warum die vorhandenen Wege nicht tragen
>
> Drei Messungen aus unserer Entscheidung (ADR-0039), die genau diese Frage <!-- d-check:ignore (fremde ADR-Kennung des Absenders, kein Verweis in dieses Repo) -->
> abgewogen hat:
>
> | Weg | Warum er scheitert |
> |---|---|
> | `ignore-refs`-Paar je Ziel | Am eigenen Breiten-Wächter versperrt — 6 von 24 Paaren überschreiten die Kappung |
> | `scan.ignore` für `docs/reviews/**` | Kostet 3887 Prüfungen über 299 Dateien und macht das Gate trotzdem nicht grün |
> | Alten Baum koexistieren lassen | Nur ein Aufschub; bricht später in einem Artefakt, das niemand mehr anfasst |
>
> `scan.ignore` ist zudem das falsche Instrument: Es nimmt die Datei aus **allen**
> Modulen. Wir wollen sie in `ids`, `matrix` und `spans` weiter geprüft haben — nur
> ihre Links auf einen bewegten Baum nicht.
>
> ## Abgrenzung
>
> Kein neues Modul, keine neue Prüfung, keine geänderte Semantik. Es ist der
> vorhandene Mechanismus, auf zwei Module ausgedehnt, die ihn als einzige der
> Referenz-prüfenden nicht haben.
>
> Wenn die Voreinstellung leer bleibt, ändert sich für bestehende Configs nichts.

---

## Erste Messung am eigenen Werkzeug (2026-09-07, **kein Entscheid**)

Festgehalten, weil sie den Entscheid vorbereitet und weil sie eine Prämisse des
CR berührt. **Sie entscheidet nichts** — der Entscheid ist ein eigener Vorgang.

**Die sechs Module stimmen.** Nachgezählt über
[`spec/lastenheft.md`](../../../spec/lastenheft.md) und
[`spec/spezifikation.md`](../../../spec/spezifikation.md): `codepaths`,
`diagrams`, `matrix`, `reviews`, `versions`, `workflows`. Auch dass `links` und
`anchors` **keine** modul-lokale Options-Sektion für Ausnahmen führen, trifft
zu.

**Die Prämisse *„kein Knopf"* trifft dagegen nicht.** Seit
[`DC-FA-REF-001`](../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
gibt es ein **geteiltes** Ventil, und es nennt beide Module ausdrücklich:

> Ein optionales, geteiltes Ventil nimmt **bestimmte Referenz-Ziele** von der
> Existenz- und Anker-Prüfung der Module `links` …, `anchors` … und `codepaths`
> … aus

Es unterdrückt genau die Klassen, um die der CR bittet
(`target-missing` **und** `anchor-missing`), und es trägt mit `in:` einen
**Quell-Skopus** — also die Einschränkung auf `docs/reviews/**`, die der CR
über `exempt-paths` erreichen will.

**Dieses Repo fährt genau diesen Fall.** 25 Ventil-Einträge über zehn entfernte
Baseline-Bäume, dahinter 28 eingefrorene Dateien; deklariert als Gate-Senkung
in [`MR-069`](../../../harness/conventions.md#mr-069).

**Der Unterschied der beiden Instrumente ist eine Eigenschaft, kein Zufall:**
`exempt-paths` ist **datei**-weit, `ignore-refs` **ziel**-weit. Gemessen in
slice-208: Ein toter Link **ohne** Baseline-Bezug in einer vom Ventil gedeckten
Datei meldet weiterhin `target-missing`. Ein `exempt-paths` auf dieselbe Datei
würde auch ihn verschlucken.

**Was damit offen bleibt und den Entscheid trägt:** Der CR nennt als Grund
gegen `ignore-refs` einen **Breiten-Wächter** — *„6 von 24 Paaren überschreiten
die Kappung"*. Das ist eine Regel im Repo des Absenders, keine Eigenschaft
dieses Werkzeugs. Ob daraus eine Produkt-Änderung folgt, ist die eigentliche
Frage.

### Nachtrag vom Auftraggeber (2026-09-07): die Ursache liegt woanders

**Mitgeteilt, nicht von hier gemessen.** Der Anlass im Repo des Absenders
(`ai-harness-init`) war, dass ein Quellen-Kommentar seine **Form verlor**: aus

```
<!-- Quelle: [02-planung/modul-05-planning-harness.md](https://github.com/pt9912/ai-harness-course/blob/v6.5.0/kurs/de/02-planung/modul-05-planning-harness.md) -->
```

wurde wieder ein **relativer** Pfad (`../../kurs/de/…`). Ein solcher Pfad löst <!-- d-check:ignore (zitierte Fremd-Form, kein Verweis) -->
im Adopter-Repo nicht auf — dort gibt es kein `kurs/de/`; er zeigt in das
**Kurs**-Repo. Daher die Befunde.

**Das ist ein reparierbarer Defekt, keine Folge der Immutabilität** — und damit
eine andere Lage als die, die der CR beschreibt. Die absolute Form ist am Ziel
gültig und überlebt jeden Pin-Sprung, weil sie den Tag **im Link** trägt statt
im Verzeichnis.

**Hier gemessen, zur Einordnung:** Das vendorte Bundle liefert die **absolute**
Form — alle **25** `Quelle:`-Kommentare unter
`.harness/baseline/v6.5.0/regelwerk/`, und ebenso in den vier vorherigen
Pin-Generationen (`v6.3.1`, `v6.0.0`, `v5.18.0`, `v5.15.0`, je Stichprobe an
`modul-05`). Der Rückfall ist also **nicht** im ausgelieferten Bundle
entstanden; wo er entstand, ist von hier aus nicht feststellbar und wird nicht
vermutet.

**Was das für den Entscheid ändert:** Die Bitte bleibt für sich prüfbar — ein
Werkzeug-Knopf steht oder fällt nicht mit dem Anlass. Aber die Begründung
*„weder reparierbar noch ignorierbar"* trägt für diesen Fall nicht mehr: Er ist
reparierbar, und für den Rest gibt es
[`DC-FA-REF-001`](../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus).
Der Entscheid hat damit zu prüfen, ob nach Abzug dieses Anlasses ein Fall
übrig bleibt, den das vorhandene Ventil nicht trägt.

### Der Wächter dagegen existiert — gemessen, nicht behauptet

**Die Forderung *„das Original-Bundle muss unverändert verwendet werden"* ist in
diesem Repo ein Gate**, kein Vorsatz:
[`make baseline-verify`](../../../harness/sensors/baseline-verify.md) prüft den
vendorten Bestand gegen `SHA256SUMS` und läuft als erstes Glied in
`make gates`.

**Bruch-Test, in einer isolierten Kopie gefahren** (der Arbeitsbaum wurde nicht
angefasst), mit **genau** dem Rückfall aus dem Nachtrag oben — absolute URL
zurück auf `../../kurs/de/…` in `modul-05-planning-harness.md`: <!-- d-check:ignore (zitierte Fremd-Form, kein Verweis) -->

```
unverändert:  Exit 0  ·  "verify ok (54 Dateien, vollständig)"
nach dem Rückfall:  Exit 1  ·  "regelwerk/modul-05-planning-harness.md: GESCHEITERT"
```

**Ein einziges geändertes Zeichen im vendorten Baum macht `make gates` rot.**
Damit ist die Klasse, die den CR ausgelöst hat, an ihrer Wurzel gefangen — und
zwar **bevor** sie als toter Link in einem eingefrorenen Artefakt auftaucht.

**Was das für den Entscheid heißt.** Die Kette des CR lautet: Bundle-Inhalt
ändert sich → Links sterben → eingefrorene Artefakte werden rot → es braucht
`exempt-paths`. Das erste Glied ist gewächtert. Ob der Absender diesen Wächter
fährt, ist von hier aus **nicht** feststellbar und wird nicht vermutet; dass es
ihn gibt, gehört in die Antwort.

**Und die Grenze dieser Messung:** Sie gilt dem Fall *„Bundle verändert"*. Der
**andere** Fall des CR — der Pin-Sprung, bei dem das alte Verzeichnis
planmäßig fällt und Adressen darauf sterben — bleibt davon unberührt. Für ihn
gilt weiterhin die Antwort oben:
[`DC-FA-REF-001`](../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
mit `in:`-Skopus, in diesem Repo 25-fach gelebt.

---

