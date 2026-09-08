# Review-Report: slice-218 — 2026-09-08

**Review-Art:** Code-Review einer Dokumentations-Änderung, deren Gegenstand eine
**gemessene Produkt-Eigenschaft** ist. Geprüft wird der Diff gegen Plan und Hard
Rules — und vor allem: **trägt die Messung, und trägt die entlastende Hälfte der
Grenze, die der Slice aufschreibt?** Die sechs Proben wurden nicht geglaubt,
sondern nachgefahren.

**Gegenstand:** slice-218 · Commit-Range `e20b7107..0d6dafd7` (`834b8ccd` Plan +
Vorprüfungen · `db998226` Beanspruchung · `0d6dafd7` die Änderung). Geändert:
`harness/sensors/adr-check.md` (+35), `harness/sensors/trace-check.md` (+14),
`docs/user/benutzerhandbuch.md` (+27) — 76 Zeilen, keine Löschung.

**Skill:** `.harness/skills/reviewer.md` @ 1.16.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-09-08

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als Tag plus Pfad in Inline-Code
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>) statt als Link.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan slice-218, §1 Ziel und Abgrenzung (die vier Ausschlüsse), §2 DoD,
  §3 die Vier-Teile-Form der Aussage, §6 Risiken, §8 Vorprüfungen
- `DC-FA-VCS-001` und `DC-FA-COMMITS-001` (die zwei Zusagen, deren Reichweite
  der Slice benennt), `DC-FA-TRK-001` (das Negativ-Ergebnis), `DC-QA-03`
- `AGENTS.md` §3.1 (Werkzeug-Klasse), §3.5, §3.7, §3.8, §4 (*„Halluzinierte
  Gates sind die häufigste Form von Harness-Lüge"*), §5 (Grenzen-Liste,
  Zählmethode, Reichweite einer Botschaft, Zitat-Geltungsbereich)
- `v6.5.0` · `regelwerk/modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin),
  Absatz *„Ein Gate ohne seine Grenze behauptet ebenfalls zu viel"*
- Reviewer-Anker 1 (stilles Grün), 2 (falscher Exit-Code), 8 (Reichweite),
  11 (zwei Module derselben Eingabe-Klasse), 15 (Scan-Achse), 17 (Proxy-Messung),
  18 (Grenzen-Liste)
- der Produkt-Code, den die Prosa beschreibt: `internal/adapter/driven/git/git.go`
  und `internal/hexagon/core/rules/vcs.go`
- go-git `v5.19.2` im Modul-Cache des Build-Images (`plumbing/object/tree.go`,
  `storage/filesystem/dotgit/dotgit.go`, `storage/filesystem/object.go`)
- die vier in §8 zitierten Register-Einträge, je `observation.md` und
  Evidence-Liste

**Eigene Läufe und Messungen** (echte Ausgaben, gekürzt auf die tragenden
Zeilen). Alle Probe-Repos im Scratchpad, alle Läufe read-only
(`--network none`, `:ro`-Mount) gegen `d-check:latest`. **Der Arbeitsbaum des
Repos war vor und nach diesem Review sauber** (`git status --short` leer).

- `make doc-check` — `d-check: 755 Datei(en) geprüft, 0 Befund(e)`, Exit 0
- sechs Pack-Namens-Proben, isoliert nachgebaut (Probe-Repo, `git repack -a -d`
  + `git prune-packed`, `.idx`/`.pack`/`.rev` gemeinsam umbenannt)
- vier weitere Proben, die der Slice **nicht** gefahren hat: partieller
  Pack-Zustand (`git maintenance run --task=loose-objects` mit vorhandenem
  kanonischen Pack), verstecktes Blob bei sichtbaren Trees, versteckter
  Zwischen-Commit, Wirksamkeit der Abhilfe
- `tracked`-Positiv-Kontrolle unter beiden Pack-Namen
- go-git-Quelle gelesen (im Build-Image, ohne Host-Go)

**Nachmessung der sechs Proben aus §1 des Plans** — Probe-Repo, alle Objekte in
*einem* Pack, `--enable vcs --disable links --staged`:

| Pack-Name | Plan behauptet | gemessen |
|---|---|---|
| `pack-<eigener Hash>` (Kontrolle) | `2 Datei(en) geprüft, 0 Befund(e)`, Exit 0 | wortgleich, Exit 0 |
| `loose-<eigener Hash>` | `staged-Basis "HEAD" nicht auflösbar`, Exit 2 | wortgleich (`… : reference not found`), Exit 2 |
| `pack-0123…4567` (gültige Form, fremder Hash) | Exit 2 | Exit 2 |
| `pack-zzzzzzzz` | Exit 2 | Exit 2 |
| `packXYZ` | Exit 2 | Exit 2 |
| `xpack-abc` | Exit 2 | Exit 2 |

**Alle sechs Zeilen stimmen.** Die Tabelle gibt die Proben korrekt wieder.

---

## Findings

### F-1 — Die entlastende Hälfte der Grenze ist falsch: `vcs` meldet im `--range`-Modus **still grün**, und zwar genau bei dem Auslöser, den der Text nennt

- `kategorie`: **HIGH** (blockierend)
- `quelle`: `AGENTS.md` §4 (*„Halluzinierte Gates sind die häufigste Form von
  Harness-Lüge"*) · `v6.5.0` · `regelwerk/modul-13-quality-gates.md` §Hard Rule
  (Doku-Disziplin), *„Ein Gate ohne seine Grenze behauptet ebenfalls zu viel"* ·
  `DC-FA-VCS-001` · Reviewer-Anker 1 und 2
- `pfad`: `harness/sensors/adr-check.md:53–55` (*„**Die entlastende Hälfte
  gehört dazu:** Exit **2** ist fail-closed — es gibt hier **kein stilles
  Grün**, der Lauf verweigert die Aussage, statt eine falsche zu treffen."*),
  `harness/sensors/adr-check.md:33–34` (*„der Lauf bricht ab, statt still grün
  zu melden"*), `harness/sensors/adr-check.md:59` (*„beide mit Exit 2"*),
  `docs/user/benutzerhandbuch.md:1593–1595` (*„**Es ist kein stiller Ausfall:**
  Der Lauf endet mit **Exit 2** und verweigert die Aussage, statt fälschlich
  grün zu melden."*)
- `befund`: Der Satz ist für `vcs` im `--range`-Modus **falsch**, und der
  Gegenbeweis läuft in genau dem Zustand ab, den der Text zwei Absätze weiter
  oben als Auslöser benennt. `git maintenance run --task=loose-objects`
  erzeugt **nicht** den gemessenen Zustand (ein einziger, unsichtbarer Pack),
  sondern einen **partiellen**: die vorhandenen kanonischen Packs bleiben
  liegen, nur die losen Objekte wandern in `loose-<Hash>.pack`. Liegt in diesem
  unsichtbaren Pack ein **Blob**, während Commit- und Tree-Objekte sichtbar
  bleiben, meldet `make adr-check RANGE=…`
  `d-check: 2 Datei(en) geprüft, 0 Befund(e)` mit **Exit 0** — und verschweigt
  dabei eine echte `core-drift-vcs`-Verletzung. Dasselbe Repo, dieselbe
  Kommandozeile, nur der Pack zurückbenannt:
  `docs/plan/adr/0001-test.md:3 … core-drift-vcs Core einer immutablen Datei hat
  sich über die Commit-Range geändert`, **Exit 1**. Zwei unabhängige Runden,
  identisches Ergebnis. Der Weg ins stille Grün ist gelesen, nicht vermutet:
  go-git `v5.19.2` `plumbing/object/tree.go:83–89` bildet
  `plumbing.ErrObjectNotFound` auf `ErrFileNotFound` ab; `git.go:239–241`
  übersetzt `ErrFileNotFound` zu `ok=false, err=nil` (*„ok=false, wenn an ref
  abwesend"*); `vcs.go:82–88` behandelt `!ok` als *nichts zu prüfen*
  (`return nil, nil`). Ein **fehlendes** Objekt wird damit ununterscheidbar von
  einer **abwesenden** Datei — die Klasse, gegen die die Zusage „fail-closed"
  gerade schützen soll. **Eskalation:** Der Modus mit dem stillen Grün ist der,
  den die klon-unabhängige CI fährt (`.github/workflows/ci.yml:77`,
  `make adr-check RANGE="$RANGE"`); der `--staged`-Hook, der fail-closed ist,
  ist der opt-in-Teil. Ein Adopter, der `git maintenance` laufen lässt, liest
  im Handbuch, sein grünes Ergebnis sei vertrauenswürdig, während sein
  ADR-Immutable-Gate nichts mehr prüft.
- `verifizierbar`: ja — Probe-Repo mit einem kanonischen Pack (C1) und einem
  zweiten Pack, der nur den geänderten ADR-Blob trägt; Trees und Commits lose.
  `docker run --rm --network none -v <repo>:/repo:ro d-check:latest --enable vcs
  --disable links --range HEAD~1..HEAD` liefert unter `loose-<Hash>` Exit 0 mit
  0 Befunden und unter `pack-<Hash>` Exit 1 mit `core-drift-vcs`. `git repack -A
  -d` im selben Repo stellt Exit 1 ebenfalls her.
- `klasse`: `fehlendes-objekt-als-abwesende-datei-gelesen` (nächste vorhandene
  Register-Klasse: `BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt` — dort für
  *Wächter, die nie fangen konnten*; der Eintrag sagt ausdrücklich *„sie trifft
  jede Prüfung, deren Erfolgsfall der Normalzustand ist"*, und das ist hier der
  Fall. Er ist damit einschlägig, deckt aber nicht den zweiten Teil: dass die
  **Doku** die fehlende Fähigkeit ausdrücklich zusagt.)

### F-2 — `vcs` und `commits` verhalten sich im selben Zustand **verschieden**; drei Stellen setzen sie gleich

- `kategorie`: MEDIUM
- `quelle`: Reviewer-Anker 11 (zwei Module derselben Eingabe-Klasse behandeln
  sie verschieden, ohne benannten Grund) · `AGENTS.md` §5 (Reichweite einer
  Aussage) · `DC-FA-COMMITS-001`
- `pfad`: `harness/sensors/trace-check.md:21` (*„Dieselbe Grenze wie bei
  [`make adr-check`], aus demselben Grund"*), `harness/sensors/adr-check.md:59`
  (*„beide mit Exit 2"*), `docs/user/benutzerhandbuch.md:1597` (*„Dasselbe gilt
  für `commits`"*)
- `befund`: Der Grund ist derselbe (beide lesen die Objektdatenbank über
  go-git), die **Grenze** ist es nicht. `commits` ist im partiellen Zustand
  fail-closed: bei einem unsichtbaren Zwischen-Commit meldet es
  `d-check: error: commit bb0e768… nicht lesbar: object not found`, **Exit 2**
  (`git.go:189`, und `ancestors()` gibt jeden `CommitObject`-Fehler weiter);
  unter kanonischem Namen findet dasselbe Repo **zwei**
  `commit-untraceable`-Befunde. `vcs` ist im selben Zustand still grün (F-1).
  Die Gleichsetzung ist damit in ihrer entlastenden Hälfte falsch — und sie
  bleibt es **auch nach** einer Korrektur von F-1: dann unterscheiden sich die
  beiden Module nachweislich, und *„dieselbe Grenze"* beschreibt sie nicht mehr.
  Wirkung: Wer `trace-check.md` liest, wird auf `adr-check.md` als vollständige
  Fundstelle verwiesen (`trace-check.md:29–31`) und übernimmt von dort eine
  Zusage, die für sein Modul zutrifft, für das andere aber nicht — die
  Ein-Ort-Pflege, die §6 Risiko 2 als Vorzug führt, transportiert hier den
  Fehler.
- `verifizierbar`: ja — Probe-Repo mit sichtbarem C1 und C3 und verstecktem C2:
  `--enable commits --range $C1..$C3` ⇒ Exit 2 mit `commit … nicht lesbar`;
  nach Rückbenennung Exit 1 mit zwei `commit-untraceable`-Befunden. Derselbe
  Zustand mit `--enable vcs` ⇒ Exit 0.
- `klasse`: `gleiche-eingabeklasse-ungleiches-versagen-als-gleich-dokumentiert`

### F-3 — Die sechs Proben messen einen Zustand, den der genannte Auslöser gar nicht erzeugen kann

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 (*„Vor einer Messung steht die Form ihres
  Gegenstands"* und *„wo der Gegenstand Code oder Konfiguration ist, gegen diese
  prüfen statt gegen die Prosa darüber"*) · Reviewer-Anker 17 und 18 ·
  `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand` ·
  `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`
- `pfad`: slice-218 §1, Absatz *„Sechs Proben, isoliert gefahren"* (Zeilen
  38–54 der Plandatei); wortgleich in `harness/sensors/adr-check.md:36–38`
  (*„**Gemessen, sechs Proben in einem isolierten Repo** (alle Objekte im
  Pack …)"*)
- `befund`: Alle sechs Proben laufen gegen ein Repo, in dem **jedes** Objekt in
  **einem** Pack liegt. Der Auslöser, den derselbe Abschnitt zwei Absätze
  später nennt (`adr-check.md:44–46`), erzeugt diesen Zustand nie: `git
  maintenance run --task=loose-objects` packt nur die **losen** Objekte und
  lässt die vorhandenen kanonischen Packs stehen — das Ergebnis ist immer ein
  **partieller** Pack. Der gemessene Zustand ist also ein Proxy für den
  Gegenstand, und er ist genau in der Eigenschaft unähnlich, um die es geht:
  im Ein-Pack-Fall wird schon der Range-Endpunkt unauflösbar (Exit 2), im
  partiellen Fall bleiben Endpunkte auflösbar und nur einzelne Objekte fehlen
  (F-1). Der zweite Griff derselben Regel hätte es unabhängig gefangen: der
  Gegenstand ist Code, und `vcs.go:82–88` (`if !ok { return nil, nil }`) sagt
  das Gegenteil der Prosa daneben. **Die Methode wird weiterverwendet:**
  `adr-check.md:64–66` erklärt einen go-git-Bump zum *„Anlass, hier
  nachzumessen"* — nachgemessen würde dann wieder der Ein-Pack-Fall, und die
  Lücke bliebe unentdeckt.
- `verifizierbar`: ja — in einem Repo mit vorhandenem kanonischen Pack
  hinterlässt `git maintenance run --task=loose-objects` nachweislich **zwei**
  `.idx`-Dateien (`ls .git/objects/pack`), nicht eine; die Proben-Beschreibung
  („alle Objekte im Pack") ist auf diesen Zustand nicht anwendbar.
- `klasse`: `probe-misst-anderen-zustand-als-den-genannten-ausloeser`

### F-4 — Die Symptom-Liste im Handbuch lässt bei `vcs` genau die Meldung aus, die der dort dokumentierte Aufruf erzeugt

- `kategorie`: LOW
- `quelle`: Maintainability (Doku-Drift) · der Zweck, den der Slice der Liste
  selbst gibt: `harness/sensors/adr-check.md:49–50` (*„Wer nur eine der beiden
  kennt, erkennt die andere nicht wieder."*)
- `pfad`: `docs/user/benutzerhandbuch.md:1589–1591` gegen den
  `make doc-immutable RANGE=…`-Aufruf in `docs/user/benutzerhandbuch.md:1579`
- `befund`: Der Absatz nennt `staged-Basis "HEAD" nicht auflösbar: reference
  not found` und `HEAD-Tree nicht lesbar: object not found` — beides Meldungen
  des `--staged`-Modus. Der Abschnitt dokumentiert unmittelbar darüber den
  **Range**-Aufruf, und der liefert gemessen eine dritte Form:
  `Range-Basis "HEAD~1" nicht auflösbar: reference not found`. Sie steht nur im
  `commits`-Abschnitt (`benutzerhandbuch.md:1621–1622`), also unter einem
  anderen Modul. Ein Adopter, der `make doc-immutable RANGE=…` fährt, sieht
  damit eine Meldung, die in *seinem* Abschnitt nicht aufgeführt ist — genau
  der Wiedererkennungs-Ausfall, den der Slice mit der Symptom-Variation
  verhindern will.
- `verifizierbar`: ja — `--enable vcs --disable links --range HEAD~1..HEAD`
  gegen ein Repo mit umbenanntem Pack liefert `d-check: error: Range-Basis
  "HEAD~1" nicht auflösbar: reference not found`, Exit 2.
- `klasse`: `symptomliste-deckt-den-eigenen-beispielaufruf-nicht`

---

## Negativbefunde (geprüft, ohne Befund)

- **Die sechs Probenergebnisse (§1 des Plans, `adr-check.md:36–42`).**
  Vollständig nachgefahren, nicht stichprobenweise: alle sechs Zeilen stimmen,
  inklusive des wörtlichen Kontroll-Strings `2 Datei(en) geprüft, 0 Befund(e)`
  und der wörtlichen Fehlermeldung. Die Tabelle gibt die Proben korrekt wieder.
- **„Beide Namensteile zählen" (`adr-check.md:40–41`).** Trifft zu — und der
  Mechanismus, den §6 Risiko 1 als *erschlossen, nicht gelesen* führt, ist jetzt
  **gelesen** und bestätigt die Aussage. Er hat zwei Stufen, nicht eine: go-git
  `v5.19.2` `storage/filesystem/dotgit/dotgit.go:298–310` verwirft beim
  Auflisten alles ohne Präfix `pack-` und ohne gültigen Hex-Hash (das erklärt
  `loose-…`, `pack-zzzzzzzz`, `packXYZ`, `xpack-abc`);
  `storage/filesystem/object.go:92–95` vergleicht danach die
  `PackfileChecksum` aus dem `.idx` mit dem Hash **aus dem Dateinamen** und
  bricht bei Abweichung ab (das erklärt die gültige Form mit fremdem Hash).
  Die Formulierung *„nur unter git's kanonischem Namen
  `pack-<Hash-des-Packs>.{idx,pack}`"* fasst beide Stufen korrekt zusammen.
  **Die zweite Erklärung ist also nicht wieder falsch** — anders als die
  erste, die §6 zu Recht verworfen hat.
- **„`git` selbst liest weiter" (`adr-check.md:41–42`).** Nachgemessen unter
  `xpack-abc`, `packXYZ` und `pack-zzzzzzzz`: `git log` und `git cat-file -t
  HEAD` arbeiten in allen drei Fällen unverändert.
- **Die go-git-Version.** `v5.19.2` in `adr-check.md:35` und
  `trace-check.md:23` stimmt mit `go.mod:6`
  (`github.com/go-git/go-git/v5 v5.19.2`) überein.
- **Die `tracked`-Aussage und ihre Positiv-Kontrolle (DoD 3).** Trägt. Probe
  mit einem gitignorierten Linkziel: unter `pack-<Hash>` **und** unter
  `loose-<Hash>` derselbe Befund `README.md:4 build/artefakt.md
  target-untracked`, Exit 1 in beiden Fällen. Das ist tatsächlich eine
  Positiv-Kontrolle und kein grüner Lauf — die Unterscheidung, die DoD (3)
  verlangt, ist eingelöst. Der Grund stimmt ebenfalls: `TrackedPaths()`
  (`git.go:258–267`) liest ausschließlich `Storer.Index()`.
- **Die Abhilfe.** `git repack -A -d` wurde end-to-end geprüft: im selben Repo
  vorher Exit 0 mit 0 Befunden, nachher Exit 1 mit dem echten
  `core-drift-vcs`-Befund. Die `-A`-Begründung (unerreichbare Objekte werden
  lose statt verworfen) ist die dokumentierte git-Semantik.
- **Der `--staged`-Pfad ist fail-closed.** Gegenprobe zu F-1: liegt der
  HEAD-Tree im unsichtbaren Pack, meldet `--staged`
  `HEAD-Tree nicht lesbar: object not found`, Exit 2. Ursache gelesen —
  `diffTreeIndex` (`git.go:124–131`) iteriert `headTree.Files()` **eager**, und
  ein fehlender Blob wird dort zum Fehler statt zu `ok=false`. Der
  `pre-commit`-Hook ist von F-1 also **nicht** betroffen; die Einschränkung
  gilt dem `--range`-Modus.
- **Reichweite auf der Modul-Achse (§3.8).** Kein Befund. Der Text sagt
  ausdrücklich *„drei Module sind geprüft, nicht alle"* und *„Über andere
  Stellen des Produkts sagt das nichts"* (`adr-check.md:62–64`); §6 Risiko 3
  führt dieselbe Frage. Die Formulierung spricht nirgends für „die
  git-basierten Module" pauschal. Diese Achse ist sauber — F-1 sitzt auf einer
  **anderen** Achse (Repo-Zustand desselben Moduls), die weder Plan noch Text
  stellt.
- **Abgrenzung §1 eingehalten.** Der Diff berührt genau drei
  Dokumentationsdateien, ausschließlich additiv (+76/−0): kein Produkt-Code
  (Punkt 1 und 2), kein neuer Sensor und keine Konfigurationsänderung
  (Punkt 3), keine Änderung an Lastenheft, Spezifikation, Sicht oder einer ADR
  (Punkt 4). Auch `CHANGELOG.md`, `README*.md` und der Handbuch-Kopf sind
  unberührt — konform zur Release-Prep-Regel aus `AGENTS.md` §5.
- **`make doc-check`** — `d-check: 755 Datei(en) geprüft, 0 Befund(e)`, Exit 0.
  Die beiden neuen Verweise lösen auf: der Anker
  `adr-check.md#grenze--was-das-grün-nicht-abdeckt` aus `trace-check.md:31` und
  der Handbuch-interne Anker auf den `vcs`-Abschnitt aus
  `benutzerhandbuch.md:1619`. Die neuen Absätze rendern als gewöhnliche
  Markdown-Absätze bzw. als Fortsetzung der nummerierten Liste; keine
  Fence-/Span-Artefakte (`spans` läuft im selben Lauf).
- **`AGENTS.md` §3.7.** Kein Befund. Es wird kein Zustandsfeld eingeführt und
  keines geändert; die neuen Texte sind Prosa in Dokumentationsdateien, keine
  Code-/Konfigurations-Kommentare. Der Herkunfts-Anker `*(seit slice-218)*`
  (`adr-check.md:66`, `trace-check.md:32`) folgt der Baseline-Form
  `seit slice-<NNN>` und hat in `harness/sensors/` bereits Präzedenz
  (`baseline-freshness.md` mit `seit slice-215`, `hooks.md` mit
  `seit welle-79`). Keine Review-Historie, keine Deliberation über Verworfenes,
  keine Chronik.
- **Die Zählungen in §8 des Plans.** Nachgemessen: 39 Beobachtungs-Verzeichnisse
  über beide Kürzel (Plan: 39); `module-promise-only-on-scan-axis` 2 Evidence-
  Dateien (2×), `wortlaut-behauptet-pruefung-die-fehlt` 9 (9×),
  `grenzen-liste-wird-als-vollstaendig-gelesen` 8 (8×),
  `form-vom-nachbarn-statt-von-der-vorlage` 1 (1×). Alle vier stimmen, und
  keiner erreicht mit diesem Slice erstmalig die Schwelle — wie behauptet.
- **Symptom-Variation.** Die Aussage, dass die Meldung mit dem Repo-Zustand
  variiert, trifft zu und ist an beiden genannten Formen nachgemessen. Sie ist
  sogar unterschätzt: es gibt mindestens drei Formen (siehe F-4).
- **Arbeitsbaum.** Vor und nach diesem Review sauber; alle Proben liefen in
  Kopien im Scratchpad, alle d-check-Läufe read-only und netzlos.

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 0 |

## Verdikt

**Blockierend.** F-1 ist ein HIGH im Gate-Pfad: Die Änderung schreibt in die
Sensor-Datei des ADR-Immutable-Gates **und** ins adopter-gerichtete Handbuch,
dass es hier kein stilles Grün gebe — während genau dieses Gate in dem Modus,
den die CI fährt, unter dem im selben Absatz genannten Auslöser still grün
melden kann und dabei eine echte Immutabilitäts-Verletzung verschweigt. Das ist
die Klasse, die `AGENTS.md` §4 als häufigste Harness-Lüge führt, und sie wiegt
hier schwerer als eine fehlende Grenze: Eine unvollständige Grenzen-Liste lässt
den Leser im Unklaren, eine falsche entlastende Zusage macht ihn sicher.

Bemerkenswert ist, **wo** der Slice trägt und wo nicht. Die sechs Proben sind
korrekt und vollständig wiedergegeben, die `tracked`-Positiv-Kontrolle ist
methodisch die stärkste Einzelleistung, die Reichweite auf der Modul-Achse ist
ausdrücklich begrenzt, und die Erklärung, die §6 zu Recht misstraute, hält der
Quelle stand. Was nicht trägt, ist die eine Achse, die §6 **nicht** gestellt
hat: nicht *welche Module*, sondern *welche Repo-Zustände* die Messung deckt.
Der Slice hat den Auslöser richtig benannt und dann einen Zustand gemessen, den
dieser Auslöser nicht erzeugt (F-3) — und aus dem Ausbleiben eines stillen
Grüns in diesem Zustand die entlastende Zusage gezogen (F-1). Der zweite Griff
aus `AGENTS.md` §5, den §8 selbst als doppelt einschlägig notiert — *wo der
Gegenstand Code ist, gegen den Code prüfen* —, hätte `vcs.go:82–88` gezeigt und
beides gefangen.

**Übergabe:** Findings an den Implementer (Rückkante Review → Implementation).
F-1 und F-2 betreffen den Text an allen drei Orten und sind vor der Closure zu
entscheiden; F-3 betrifft die Proben-Beschreibung im Plan **und** die
Nachmess-Anweisung in `adr-check.md:64–66`, die die Methode fortschreibt; F-4
ist eine Ergänzung im Handbuch. Ob F-1 als **Doku**-Korrektur (die Grenze
richtig aufschreiben) oder als **Produkt**-Frage (`FileAt` unterscheidet
„abwesend" nicht von „nicht lesbar") beantwortet wird, ist eine
Architekt-Entscheidung, keine Reviewer-Entscheidung — §1 Punkt 2 schließt
Produkt-Code für *diesen* Slice aus, und §4 des Plans nennt für genau diesen
Fall die Rückführung nach `next/`. Die **Finding-Klassen** gehen zusätzlich in
die Slice-Closure §7 und von dort in den Zähler;
`wortlaut-behauptet-pruefung-die-fehlt`,
`zaehlmethode-misst-proxy-statt-gegenstand` und
`grenzen-liste-wird-als-vollstaendig-gelesen` sind in diesem Lauf je erneut
belegt. Dieser Report ersetzt keine Verifikation — DoD-Konformität und der
volle Gate-Stand sind Sache des Verifiers.
