# Welle wave-98: Baseline-Pin auf `v7.0.0` heben, die Namensregel für Neues einführen und die übrigen Neuerungen vorlegen

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-11-25.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Der vendorte Baseline-Baum steht auf
[`v7.0.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v7.0.0) (Kurs-Welle 168,
2026-11-23), jeder lebende pin-gebundene Verweis zeigt auf ihn, und die **Namensregel** des Regelwerks gilt
in diesem Repo **ab jetzt für Neues**: die Spezifikation heißt `specification.md`, neue Wellen heißen
`wave-<NN>`, neue Archive `archive.zip`, neue Herkunfts-Anker `seit wave-<NN>` — **der Bestand bleibt**, und
beide Schreibweisen nebeneinander sind eine benannte Adaption. Das Produkt ändert sich an den Stellen, die §3 (d)
einzeln nennt, an keiner anderen; zu jeder **übrigen** Neuerung des Bundles hat der Auftraggeber eine
Entscheidungsvorlage. **Übernommen wird außer der Namensregel nichts.**

Zwei Akzeptanzkriterien, an denen sich das Ziel spiegelt:

- *Given der Pin steht auf `v7.0.0`, when `make baseline-freshness` läuft, then meldet es den Pin als aktuell
  (heute Exit 3, „NEUER RELEASE verfügbar (Pin v6.18.1)“, gemessen 2026-11-25), und `make gates` ist grün.*
- *Given eine neue flache Datei `wave-<NN>-….md` unter `docs/plan/planning/` ohne Zeiger in der Roadmap, when
  `make planning-check` läuft, then meldet es `wave-drift` (heute: Exit 0, 0 Befunde — die Datei ist für das Gate
  unsichtbar, gemessen im Probelauf unten).*

**Der Zuschnitt ist die kleinste tragende Fassung je Frage** — die Lehre aus
[`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md):
Der Bump ist mechanisch und gut bekannt (die Pin-Kette von [`MR-011`](../../../harness/conventions.md#mr-011)
bis [`MR-080`](../../../harness/conventions.md#mr-080); zuletzt
[`slice-277`](done/wellenlos/slice-277-baseline-v6181-bump.md)); die Namensregel hat eine **kleinste** Antwort
(Konfiguration und Werkzeug, **kein** neues Produkt-Merkmal) und eine größere (das Modul `planning` liest mehrere
Präfixe), die in §5 Frage 3 mit ihrem Preis daneben steht; die Sichtung ist eine **Entscheidungsvorlage, kein
Umbau**. Ein Release ist nicht Teil der Welle (§6).

**Entscheide des Auftraggebers (2026-11-25, bindend; der zweite ersetzt den ersten Entscheid der Anfrage, der die
Spezifikation deutsch halten wollte):**

1. **Namensregel, Reichweite „ab jetzt, Bestand bleibt“.** Das „Deutsch bleiben“ galt dem **Inhalt**: Die Prosa
   bleibt deutsch (*Welle*, *Spezifikation*). Die **Dateinamen** folgen der Regel aus `v7.0.0`
   (`grundlagen-harness-dateien.md`, Absatz *Dateinamen sind englisch*): englisch, wo es ein gängiges englisches
   Wort gibt; `lastenheft.md` und fremde Namen (`AGENTS.md`, `README.md`, `Makefile`, `.d-check.yml`) bleiben; das
   Präfix der Welle gilt mit; der Titel-Teil bleibt frei und darf deutsch sein. **Jetzt umbenannt:** die
   Spezifikation (`spec/spezifikation.md` → `specification.md` (unter `spec/`)) und die sieben Vorlagen-Namen (sie kommen mit
   dem vendorten Baum). **Ab jetzt für Neues:** `wave-<NN>-<titel>.md` mit `wave-<NN>-results.md`, `archive.zip`,
   `seit wave-<NN>`. **Bestand bleibt:** abgeschlossene `welle-<NN>`-Dateien, `archiv.zip`, vorhandene `seit welle-`-Anker
   werden nicht umbenannt. Die Koexistenz ist eine Adaption (Form: Architect, §5 Frage 1).
2. **Umfang „Bump plus Neuerungen“:** (a) der Bump selbst, (b) eine Sichtung der Neuerungen, an deren Ende der
   Auftraggeber je Neuerung entscheidet, was übernommen wird. Die Sichtung übernimmt nichts; Übernahmen sind eigene
   Folge-Slices nach seiner Entscheidung. Die Namensregel ist die **eine** Neuerung, die der Auftraggeber schon
   entschieden hat; sie steht deshalb in dieser Welle und nicht in der Vorlage.

**Gemessener Ausgangsstand des Bundles** (Planner, 2026-11-25, am entpackten Release-Asset `lab-regelwerk.zip`
von `v7.0.0` — SHA256 `34a7866d…e0e`, gleich dem `SHA256SUMS` des Release — gegen den vendorten
Baum `.harness/baseline/v6.18.1/`; der Implementer **wiederholt** die Messung am vendorten Baum
von `v7.0.0` und am alten Baum aus `git archive`, die Zahlen hier sind Erwartung, nicht Beleg).
Gegenstand der Zählung ist die **Datei im Bundle**, nicht der Verweis auf sie:

| Gemessen | Zahl | Form der Zählung |
|---|---|---|
| Bundle-Dateien alt / neu | 54 / 55 | `find` über `regelwerk/` und `templates/` ohne `SHA256SUMS`; vendort je +1 für `SHA256SUMS` (55 / 56) |
| gleichnamig in beiden | 47 | Rest: 7 umbenannt (alt) ↔ 7 (neu), 1 neu |
| davon byte-verschieden | 38 | `cmp` je Datei; **862** Diff-Zeilen (`diff`, Zeilen mit `<`/`>`) |
| davon **inhaltlich** verschieden | 31 (19 Regelwerk, 12 Vorlagen) | ohne die `<!-- Quelle:`-Zeilen, Versionsstempel auf einen Platzhalter normalisiert; **804** Diff-Zeilen. Die übrigen 7 ändern nur Stempel und Quell-Verweis |
| umbenannt | 7 | `welle.template.md`→`wave.template.md`, `welle-results.template.md`→`wave-results.template.md`, `archiv-stub-slice`/`-welle`→`archive-stub-slice`/`-wave`, `spezifikation.template.md`→`specification.template.md`, `NNNN-titel`→`NNNN-title`, `MR-NNN-titel`→`MR-NNN-title` (je 2–28 inhaltliche Diff-Zeilen) |
| neu | 1 | `templates/harness/sensors/environment-check.template.md` (74 Zeilen) |

**Größte Bewegungen** (inhaltliche Diff-Zeilen): `modul-14-docker-harness` 258,
`templates/.d-check.yml` 74, `modul-04-adrs` 57, `modul-05-planning-harness` 49,
`templates/README.md` 35.

**Gemessene Reichweite der Namensregel** (Planner, 2026-11-25, Stand `590dde3d`; gezählt wird die git-getrackte
Datei bzw. Zeile mit der Zeichenkette, nicht der Verweis auf sie; der Implementer zählt neu). Die Klassen sind
die Frozen-Klassen der Eigenschaft nach ([`MR-070`](../../../harness/conventions.md#mr-070)), nicht nach
Verzeichnis:

| Klasse | Dateien | Zeilen | Kommando |
|---|---|---|---|
| `Accepted`/`Proposed`/`Superseded`-ADRs, Verweis als Link-Ziel | 64 (von 72, die den Namen nennen) | 119 | `git grep -c '\](\([^)]*/\)\?spezifikation\.md' -- docs/plan/adr` |
| eingefrorene Nicht-ADR-Dateien mit Link-Ziel (`done/` 23, `docs/reviews` 13, `docs/plan/cr` 2, `conventions/done` 1) | 39 | 80 | dasselbe, über `docs/plan/planning/done docs/reviews harness/conventions/done docs/plan/cr` |
| lebende Markdown-Dateien ohne die Spezifikation | 29 | 53 (davon 42 mit Link-Ziel) | `git grep -c 'spezifikation\.md' -- '*.md'` ohne Historie, ADRs, Register-Belege und die Pläne dieser Welle |
| die Spezifikation selbst | 1 | 161 (158 als Link-Ziel auf sich selbst) | `grep -c 'spezifikation\.md' spec/spezifikation.md` |
| Konfiguration und Skripte | `.d-check.yml` 6 Zeilen; `Dockerfile`, `Makefile` je 1 Kommentar; `tools/` 3 Zeilen (davon `tools/image-test.sh` Z. 163 **funktional**, 2 Kommentare) | | `git grep -n 'spezifikation\.md' -- .d-check.yml Dockerfile Makefile tools` |
| Go, Produkt | `manual.go` (Pfad-Konstante und `go:embed`); 2 Tests lesen oder erwarten den Namen; 90 Kommentarzeilen in 33 Dateien; 4 Zeilen sind **Konsumenten-Beispiele** (`config_template.go`, `suggest.go`) | | `git grep -n 'spezifikation\.md' -- '*.go'` |
| **Dateinamen**, die auf `spezifikation.md` enden (außer der Spezifikation) | 2 (`done/…-architektur-und-spezifikation.md`, `done/…-spec-vergabe-spezifikation.md`) | | `git ls-files \| grep -c 'spezifikation\.md$'` ergibt 3 — ein Token-Tausch ohne Wortgrenze trifft sie |
| Wellen: Pfade mit `welle-<NN>` im Namen | 344 (14 außerhalb von `done/`) | | `git ls-files \| grep -c 'welle-[0-9]'` |
| Wellen: Abschluss-Register-Zeilen / Ergebnisnotizen | 38 / 38 | | `grep -c '^\| welle-' docs/plan/planning/in-progress/roadmap.md` |
| Archive `archiv.zip` (Welle und wellenlos) | 188 | | `git ls-files \| grep -c 'archiv\.zip$'` |
| Herkunfts-Anker `seit welle-<NN>` (Instanzen) | | 24 | `git grep -h -o 'seit welle-[0-9]*' -- . ':!.harness/baseline' ':!docs/plan/planning/done' ':!docs/reviews'` |

**Probelauf in einem Wegwerf-Klon** (Planner, 2026-11-25; nicht im Repo, nichts davon ist umgesetzt) — was die
Namensregel nach sich zieht, **gemessen** statt geschätzt:

1. **Der Build bricht am Umbenennen.** Ohne Nachzug in `manual.go` scheitert `go:embed` (`pattern spec/spezifikation.md:
   no matching files found`); `make doc-check`, `make adr-check` und jedes Target, das das Image baut, brechen
   mit. Der **reine Move-Commit allein ist deshalb nicht baubar** — die Zusage von
   [`AGENTS.md`](../../../AGENTS.md) §3.3 (Zwischenstand nicht die Spitze eines Push) trägt hier, und der Hook prüft
   den Arbeitsbaum.
2. **`make test` wird an genau zwei Stellen rot:** `TestManual_GibtAbschnitteAus` (erwartet die Kopfzeile
   `==> spec/spezifikation.md:`) und `TestAllReasonsDeckungGegenSpezifikationGrundCodes` (liest die Spezifikation
   am alten Pfad).
3. **Der ADR-Gate trägt den Nachzug — in einem Commit, nicht in zweien.** Ein Commit mit dem reinen Rename **und** den
   Link-Ziel-Nachzügen der 64 ADRs: `make adr-check STAGED=1` Exit 0, 0 Befunde
   ([`ADR-0103`](../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md),
   [`ADR-0115`](../adr/0115-vcs-folgt-umbenennung.md), `vcs.follow-renames: true` steht in der
   [`.d-check.yml`](../../../.d-check.yml)). **Zwei Commits** (erst der Move, dann die ADR-Links): der zweite meldet
   `core-drift-vcs` für **60** der ADRs — die Rename-Erkennung sieht nur die Spanne des Commits. Über die ganze Spanne
   beider Commits (`RANGE=`, wie die CI sie fährt) ist es wieder Exit 0.
4. **`make doc-check` ist erreichbar, mit einem Ventil.** Nach Rename, Link-Nachzug in allen Markdown-Dateien,
   `.d-check.yml` und `manual.go`: 123 Befunde `codepath-missing` (`spec/spezifikation.md` im Inline-Code) in 83
   Dateien — **104** davon in eingefrorenen Klassen (`done/` 47, ADRs 44, Register-Belege 5, `conventions/done` 5,
   `CHANGELOG.md` 2, `docs/plan/cr` 1), **0** in lebenden Dokumenten nach deren Nachzug (Probelauf, Stand `590dde3d`;
   überholt durch die Endstand-Messung im nächsten Absatz). Ein einziger Eintrag in
   `codepaths.ignore-refs` bringt `make doc-check` auf **0 Befunde**; er wirkt **repo-weit** (die Grenze steht an
   der Liste in der `.d-check.yml`, [`ADR-0025`](../adr/0025-codepaths-ignore-refs.md)).
   **Endstand-Messung** (Planner, 2026-11-26, Stand `00d4f352` nach [`slice-297`](done/slice-297-spezifikation-heisst-specification.md);
   Wegwerf-Klon ohne den Ventil-Eintrag, `docker run --rm --network none -v <klon>:/repo:ro d-check:latest`): **120**
   Befunde `codepath-missing` in 80 Dateien (`cut -f1 <ausgabe> | sed 's/:.*//' | sort -u | wc -l`), **109** davon in
   eingefrorenen Klassen (`done/` 51, ADRs 44, Register-Belege 5, `conventions/done` 6, `CHANGELOG.md` 2,
   `docs/plan/cr` 1), **11** in lebenden Dateien, die den alten Pfad als **Gegenstand** nennen (diese Welle 4,
   `slice-297` 4, `slice-296` 2, [`MR-081`](../../../harness/conventions.md#mr-081) 1; je Klasse
   `cut -f1 <ausgabe> | grep -c '^<pfad>'`). Die Zahl zählt ihre eigene Nennung mit: jede Zeile, die den Pfad in
   Inline-Code trägt, ist ein Befund, und das Ventil deckt sie alle. Die Zahl 119 aus
   [`slice-297`](done/slice-297-spezifikation-heisst-specification.md) (Zweite Plan-Änderung) gilt für ihren
   Messstand und stimmt am Endstand nicht mehr (120).
5. **Das Modul `planning` liest genau ein Präfix** (`wavePrefix` in `planning_waves.go`, acht Stellen; `glob` und
   `results-glob` müssen dasselbe literale Präfix tragen, `configyaml.go` Z. 1567–1572). Folge, gemessen am Klon: mit der
   Konfiguration von heute (`welle-*.md`) ist eine neue `wave-99-….md` **unsichtbar** — `make planning-check` Exit 0, 0
   Befunde, auch ohne Zeiger; mit `glob: "wave-*.md"` und `results-glob: "wave-*-results.md"` meldet es `wave-drift`
   für die Datei ohne Zeiger, und die 38 Abschluss-Register-Zeilen `welle-NN` samt ihren 38 Ergebnisnotizen sind
   **ungewächtert** (kein Befund, auch wo einer stünde). **Beide Bevölkerungen zugleich kann nur eine Änderung am
   Produkt** (mehrere Präfixe) — das ist die größere Antwort in Frage 3.
6. **`tools/archive-wave` bricht fail-closed, nicht still:** `make archive-wave WELLE=wave-99` endet mit Exit 2
   („weder Welle-Plan noch Slices für wave-99 gefunden — Abbruch, fail-closed statt leeres Archiv“), weil
   `collect.go` (Z. 16, 124) nur `welle-(\d+)` im `**Welle:**`-Feld liest; dasselbe Fixture mit `welle-99` findet
   den Slice.

**Vier Befunde prägen den Zuschnitt:**

1. **Die sieben Vorlagen-Umbenennungen treffen das Repo mechanisch.** `plan-welle.md` und
   `close-welle.md` nennen die alten Namen (`git grep -n "welle\.template\|welle-results\.template\|archiv-stub"
   -- .claude`), ebenso `harness/conventions.md` (`MR-<NNN>-titel.template.md`) und zwei Kommentarzeilen in
   `tools/archive-wave/stub.go`. Das gehört zum Bump ([`slice-295`](done/slice-295-baseline-v700-bump.md)).
2. **Die Namensregel hat zwei Hälften mit zwei Slices.** Die **Spezifikation** umzubenennen ist ein Rename mit
   gemessener Reichweite und eigenem Commit-Zwang ([`slice-297`](done/slice-297-spezifikation-heisst-specification.md));
   die **neuen Namen** der Wellen, Archive und Anker sind ein Werkzeug-/Sensor-/Command-Thema
   ([`slice-298`](done/slice-298-neue-namen-wave-archive-anker.md)). Beide in einem Slice wären vier Liefer-Punkte.
3. **Eingebettete Dokumente und `--manual`.** Das Produkt bettet `spec/spezifikation.md` und das Handbuch ein
   ([`ADR-0107`](../adr/0107-mitgelieferte-dokumente-im-modul-root.md), `manual.go`); die Kopfzeile der Ausgabe
   trägt den Pfad (`==> <pfad>:<zeile>`). **Vertrag:** [`DC-FA-CLI-013`](../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen)
   verlangt eine Kopfzeile, „die Dokument und Zeile nennt“, die Spezifikation §[`DC-FA-CLI-013.a`](../../../spec/specification.md#dc-fa-cli-013a--handbuch-und-spezifikation-ausgeben---manual) Schritt 5 `==> <pfad>:<zeile>`,
   Handbuch und `operations.md` `==> <datei>:<zeile>` — **keine** dieser Stellen nennt den Dateinamen wörtlich
   (`git grep -n '==> ' -- spec docs/user`). Der Wortlaut des Vertrags bleibt; **beobachtbar** ändern sich die Bytes
   der Ausgabe (Kopfzeile, die Link-Ziele der Spezifikation, die eine Handbuch-Zeile, der Pin), und ein Konsument, der
   `==> spec/spezifikation.md:` greppt, bricht. Das gehört in den CHANGELOG der Release-Prep
   ([`AGENTS.md`](../../../AGENTS.md) §5 Zeile 17), nicht in einen Feature-Commit.
4. **Die Spezifikation trägt ihren Pin, und `spec/lastenheft.md` Zeile 4762 nennt ihren Pfad** (Messmethode zu
   [`DC-QA-01`](../../../spec/lastenheft.md#dc-qa-01--performance)). Eine Änderung am Lastenheft ist ein
   Change-Request-Fußabdruck (Version, Historie; [`MR-032`](../../../harness/conventions.md#mr-032)) — Frage 6.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis**
dieser Welle: Steht er in der Slice-Liste unten, ist er falsch platziert.

- **Ein neuer Baseline-Release liegt vor, und der gepinnte Stand ist unverändert.** `make baseline-freshness`
  endet mit Exit 3 und der Zeile „NEUER RELEASE verfügbar (Pin v6.18.1): v7.0.0“, zugleich „gepinnter Tag
  v6.18.1 upstream unverändert (Bytes == vendored SHA256SUMS)“ — gelesen am 2026-11-25. Der Nachtlauf
  `upstream-drift` ist seit demselben Tag rot (`make nightly-state`; die Meldung ist **planmäßig** gelesen, nicht
  weggeklickt: ein Fremd-Release, [`MR-051`](../../../harness/conventions.md#mr-051)). Es ist eine reine
  Fortschreibung, kein Nachholen stiller Drift.
- **Der Auftraggeber hat die Welle angefordert** (2026-11-25, in der Sitzung) und am selben Tag den Entscheid zur
  Namensregel **korrigiert** (§1, Entscheid 1).
- **welle-97 ist geschlossen:** ihre Welle-Datei liegt unter `done/welle-97/`
  (`ls docs/plan/planning/done/welle-97/`), ihre Closure-Notiz als `done/welle-97-results.md`.
- **`git fetch`:** `origin/main` liegt nicht vor dem lokalen Hauptzweig (`git rev-list --left-right --count
  origin/main...HEAD`: links `0`; rechts stehen die ungepushten Commits dieser Welle). `in-progress/` trägt nur die
  Roadmap, `next/` ist leer, `open/` trägt nur die Slices dieser Welle.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle vier Slices der Welle sind `done/`.
- `make fullbuild` grün auf dem Endstand.
- **Das Mehr gegenüber den Slice-DoDs.** (d) bis (g) sind Aussagen über das Repo als Ganzes, die kein
  einzelner Slice abhaken kann — sie stehen deshalb **nicht** in einer Slice-DoD
  ([`AGENTS.md`](../../../AGENTS.md) §5 Regel 10). (a) bis (c) liefern
  [`slice-295`](done/slice-295-baseline-v700-bump.md) und [`slice-297`](done/slice-297-spezifikation-heisst-specification.md);
  die Welle **misst sie am Endstand gegen**, nachdem alle Slices geschrieben haben — ein Nachzug
  dort darf kein Rot erzeugen. Zu jeder Aussage das Rot **ohne** die Arbeit der Welle:
  - **(a) Pin und Bestand gehören zusammen.** `make baseline-verify` Exit 0 **und** `make baseline-freshness`
    Exit 0 mit „Pin v7.0.0 ist der neueste Release-Tag“ und „upstream unverändert“. *Rot ohne die Welle:*
    `make baseline-freshness` endet heute mit Exit 3. Erscheint bis zur Closure ein neuerer Release, ist das
    ein neuer Anlass und nicht Teil dieser Welle; die Closure-Notiz nennt ihn.
  - **(b) Kein lebender Verweis nennt mehr den alten Pin.** Die Treffermenge von
    `git grep -l 'v6\.18\.1' -- . ':!docs/plan/planning/done' ':!docs/reviews' ':!CHANGELOG.md' ':!harness/conventions/done' ':!.harness/baseline'`
    ist **genau** die Menge, die der Bump-Eintrag als Vergangenheits-Aussage oder eingefrorene Nennung
    mit benanntem Grund führt (erwartet: der Bump-Eintrag selbst und die eingefrorenen Dateien hinter dem
    Ventil, §5 Frage 5). **Und** die Symlinks lösen auf den neuen Baum auf
    (`for l in .claude/rules/*; do readlink "$l"; done | grep -c 'v6\.18\.1'` ergibt 0 — `git grep` liest keinen
    Symlink-Inhalt). *Rot ohne die Welle:* 55 Dateien mit 103 Zeilen (Stand `76cf7d51`, vor Anlage dieser Pläne)
    und 7 Symlinks. Die Zählung nennt ihre Form: gezählt wird die Zeichenkette in Dateien, nicht der Verweis.
  - **(c) Die gehaltenen Zitate und Verweise stehen auf dem neuen Stand.** `make doc-check` ohne Befund, darunter
    `citations` für die **15** lebenden `d-check:cite`-Spannen auf den Baum (davon verschieben sich **7**, am
    Vergleich derselben Zeilenbereiche beider Bäume gemessen) und die Spannen der Pläne dieser Welle, **und** `links`
    und `codepaths` nach dem Umbenennen der Spezifikation; `ignore-refs` in der
    [`.d-check.yml`](../../../.d-check.yml) wächst nur um die fünf Einträge des Bump-Eintrags
    ([`MR-069`](../../../harness/conventions.md#mr-069)) und den einen `codepaths`-Eintrag der Spezifikation
    ([`ADR-0025`](../adr/0025-codepaths-ignore-refs.md)). `make adr-check RANGE=<Vor-Bump>..HEAD` — die
    Spanne, die die CI fährt — endet mit Exit 0. *Rot ohne die Welle:* nach dem Entfernen des alten Baums
    melden die Pfad-Links und die Spannen Befunde (Pfad-Links `target-missing`, Spannen das Modul `citations`); nach
    dem bloßen Tag-Tausch bleiben die 7 verschobenen Spannen rot; nach dem bloßen Umbenennen der Spezifikation
    melden `make doc-check` Befunde (123 im Probelauf, Stand `590dde3d`) und der Build bricht (Probelauf 1 und 4 in §1).
  - **(d) Das Produkt ändert sich an genau den benannten Stellen.** (d1) **Quelltext:** `git diff --stat
    <Vor-Bump>..HEAD -- '*.go' ':!tools' Dockerfile go.mod` nennt **nur** `manual.go` (die Pfad-Konstante und das
    `go:embed`), `internal/adapter/driving/cli/cli_manual_test.go` und `internal/hexagon/core/app/diagnose_test.go`
    — keine Kommentarzeile (§5 Frage 7). Kein Standardwert, keine Logik, kein Grund-Code. (d2)
    `make blackbox-probe REF=<Vor-Bump>` —
    der Commit **vor** dem ersten Bump-Commit — endet mit „byte-identisch“ über seine Vergleiche (stdout, stderr,
    Exit getrennt; Formen: Standard, `--json`, `--yaml`, `--doctor`). **Nicht `REF=v0.93.0`:** auch dieses
    Vorher-Image kennt das Profil (`.d-check.yml` ist seit `v0.93.0` unverändert, `git diff --stat v0.93.0..HEAD --
    .d-check.yml` ist leer), aber `v0.93.0` trägt gegenüber dem Hauptzweig zwei Änderungen an Bauteilen des
    Images (`git diff --stat v0.93.0..HEAD -- '*.go' Dockerfile go.mod manual.go spec docs/user`: `Dockerfile` und das
    Handbuch), die nicht aus dieser Welle stammen. (d3) **`--manual`** — hier ist die Aussage **enger**, nicht
    gelockert: Für die Begriffe `planning`, `Spezifikation` und `manual` wird die Ausgabe von Vorher- und
    Nachher-Image verglichen (`docker run … --manual <begriff>`, `diff`). **Jede** abweichende Zeile ist eine von
    fünf Klassen: (i) die Kopfzeile `==> spec/spezifikation.md:<n>` → `==> spec/specification.md:<n>` (die Zeilennummer
    bleibt, weil kein Zeilenumbruch entsteht); (ii) eine Zeile der Spezifikation, deren Link-Ziel
    `spezifikation.md#…` → `specification.md#…` wechselt (erwartet bis zu 158 Zeilen, nicht alle im Ausschnitt);
    (iii) die Pin-Nennung `v6.18.1` → `v7.0.0`; (iv) die eine Handbuch-Zeile, die die Spezifikation beim Pfad nennt
    (`docs/user/benutzerhandbuch.md`, Abschnitt der Konfigurations-Defaults). (v) die Prosa-Zeile 920 der Spezifikation, die den Pfad zu `--manual` nennt und in der Ausgabe zu `manual` und `Spezifikation` steht (Fundstelle: [`slice-297`](done/slice-297-spezifikation-heisst-specification.md) §3, Zweite Plan-Änderung, dritter Punkt; die Messung bleibt bei dieser einen Zeile). Eine sechste Klasse ist ein Befund.
    (d4) Der Vertrag [`DC-FA-CLI-013`](../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen)
    bleibt im Wortlaut; die Closure-Notiz nennt den Hinweis für die Release-Prep (CHANGELOG-Eintrag zu (d3)). *Rot
    ohne die Welle:* (d) ist **nicht** rot, wenn nichts geschieht — es hält eine Zusage fest, die der Rename
    verletzen könnte (ein mitgeänderter Quelltext, eine sechste Abweichungsklasse in `--manual`).
  - **(e) Das Trigger-Audit** (seit [`slice-241`](done/wellenlos/slice-241-trigger-audit-adoption.md)) ist über
    die Closure vollzogen und in `wave-98-results.md` belegt. Fällig ist mindestens der Auflösungs-Trigger von
    [`MR-080`](../../../harness/conventions.md#mr-080) („die nächste Pin-Hebung“); die übrigen Trigger-Klassen
    liest die Closure, sie sind hier nicht vorweggenommen.
  - **(f) Die Entscheidungsvorlage liegt dem Auftraggeber vor.** Jede der **39** Dateien mit Bewegung (31
    inhaltlich verschieden, 7 umbenannt, 1 neu — auch die umbenannten tragen inhaltliche Diff-Zeilen, die der
    Entscheid zur Namensregel nicht beantwortet) hat eine Zeile mit Klassifikation und Empfehlung; die 7 Dateien, die nur ihren
    Stempel ändern, stehen in **einer** Sammelzeile mit ihrem Beleg. Gezählt wird mit dem Kommando aus
    [`slice-296`](next/slice-296-baseline-v700-neuerungen-sichtung.md) §3. **Die Entscheidungen des Auftraggebers
    sind nicht Closure-Bedingung** — die Welle schließt, wenn die Vorlage ihm vorliegt; jede Übernahme ist ein
    eigener Slice nach seiner Antwort (§6). *Rot ohne die Welle:* die Datei existiert nicht.
  - **(g) Die neuen Namen werden gesehen, der Bestand bleibt.** Am Endstand, in einer Wegwerf-Kopie des Repos: eine
    flache `wave-99-probe.md` ohne Zeiger in der Roadmap ⇒ `make planning-check` meldet `wave-drift`, mit Zeiger Exit 0;
    `make archive-wave WELLE=wave-99` findet den Slice mit `**Welle:** wave-99`, und `make archive-wave WELLE=welle-97`
    findet weiter die Slices der Bestands-Welle; eine `wave-99-results.md` ohne den Abschnitt `## Was wurde geliefert`
    meldet `make verify-closure-notes`; ein `wave-99`-Token im Körper von Spec, Sicht oder ADR meldet `matrix`. **Und**
    `git diff --stat -M <Vor-Bump>..HEAD -- docs/plan/planning/done` nennt außer den Nachzügen aus Frage 5 keine Datei
    (Bestand bleibt). *Rot ohne die Welle:* gemessen im Probelauf (Punkt 5 und 6 in §1): `wave-99` unsichtbar für
    `planning-check`, `archive-wave` bricht mit Exit 2.
- Closure-Notiz in `wave-98-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-295`](done/slice-295-baseline-v700-bump.md) | **Bump:** Baseline-Pin auf `v7.0.0` — Re-Vendor, Spiegel, `d-check:cite`-Spannen, Vorlagen-Namen in den Verweisen, Bump-Eintrag (Delta vollständig nach Datei, keine Neuerung übernommen) | [`MR-011`](../../../harness/conventions.md#mr-011) (Pin-Kette), [`MR-021`](../../../harness/conventions.md#mr-021) |
| [`slice-297`](done/slice-297-spezifikation-heisst-specification.md) | **Namensregel, Hälfte 1:** `spec/spezifikation.md` → `specification.md` (unter `spec/`) — reiner Rename samt Link-Nachzug, Inhalt, `manual.go` und die zwei Tests, Ventil für die eingefrorenen Nennungen | [`DC-FA-CLI-013`](../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen) (Vertrag bleibt), [`ADR-0103`](../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md), [`ADR-0115`](../adr/0115-vcs-folgt-umbenennung.md) |
| [`slice-298`](done/slice-298-neue-namen-wave-archive-anker.md) | **Namensregel, Hälfte 2:** neue Wellen `wave-<NN>`, Archive `archive.zip`, Anker `seit wave-<NN>` — `tools/archive-wave`, die Sensoren der `.d-check.yml` und des Closure-Profils, die Command-Dateien; der Adaptions-Eintrag zur Koexistenz | [`MR-011`](../../../harness/conventions.md#mr-011) |
| [`slice-296`](next/slice-296-baseline-v700-neuerungen-sichtung.md) | **Sichtung:** je Bewegung der übrigen Dateien eine Klassifikation (bereits befolgt · neu und anwendbar · nicht anwendbar) mit Empfehlung, dazu die Frage zu `environment-check` und zur „Umgebungs-Zusage“ gegen die vorhandenen Gates — **eine** Entscheidungsvorlage, keine Übernahme | [`MR-011`](../../../harness/conventions.md#mr-011) |

**Die Reihenfolge ist bindend: 295 → 297 → 298 → 296** — und die Nummern folgen ihr **nicht**. Eine Slice-Nummer ist
die Reihenfolge der Anlage, nicht der Arbeit; `slice-297` und `slice-298` entstanden, als der Entscheid zur
Namensregel kam, und liegen vor dem schon angelegten `slice-296`, weil die Sichtung das Repo **nach** der
Namensregel liest. Die Vorläufer-Form steht im Drift-Log der Roadmap (ein „davorgeschnittener“ Slice trägt eine
höhere Nummer als der, vor den er tritt). Die Nummern bleiben, weil die Pläne ungepusht und unbeansprucht sind
und eine Umnummerierung nur Verweise bewegte, ohne etwas zu klären.

**Vier Slices statt eines Namens-Slice.** Der Namensregel-Slice wäre der dritte, hätte aber **vier**
Liefer-Punkte gehabt (Rename samt Verweisen · Produkt-Code und Tests · neue Namen im Werkzeug · neue Namen in den
Sensoren und Commands, dazu die Adaption) und zwei Schichten (Dokumentation/Konfiguration, Code in Produkt **und**
Werkzeug) — die Grenze von `modul-05-planning-harness.md` §Ziel-Form: Slice ist ≤ 3 Liefer-Punkte. Die zwei Hälften
sind unabhängig: Der Rename kommt ohne neue Namen aus, die neuen Namen ohne Rename. Der Adaptions-Eintrag zur
Koexistenz gehört in `slice-298` (er beschreibt die beiden Schreibweisen der Wellen, Archive und Anker); die
Spezifikation selbst steht nach dem Rename **konform** da und braucht keine Adaption.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine Welle. Jede Übernahme einer **weiteren** Neuerung ist ein eigener Slice **nach** der Antwort des
  Auftraggebers auf die Vorlage (§6) und hängt an dieser Welle nicht.
- Wird blockiert von: keiner.
- **Die Reihenfolge ist fachlich bindend: erst pinnen, dann dem Kanon folgen.** Erst
  [`slice-295`](done/slice-295-baseline-v700-bump.md) (der neue Baum steht; die Vorlagen heißen jetzt `wave.template.md`
  und so weiter), dann [`slice-297`](done/slice-297-spezifikation-heisst-specification.md) (der Rename berührt die Zeile der Spezifikation,
  die der Bump ändert, und die `d-check:cite`-Spannen, die der Bump neu ankert) und
  [`slice-298`](done/slice-298-neue-namen-wave-archive-anker.md) (die Commands nennen die neuen Vorlagen-Namen), zuletzt
  [`slice-296`](next/slice-296-baseline-v700-neuerungen-sichtung.md): Die Sichtung liest den **vendorten** Baum von
  `v7.0.0` und das Repo **nach** der Namensregel, nicht das Scratchpad, und den alten Baum aus `git archive <Vor-Bump>`;
  ein Audit gegen zwei Stände, von denen einer nirgends im Repo liegt, ist kein Audit (dieselbe Begründung wie in
  der Vorgänger-Migration [`welle-85`](done/welle-85/welle-85-baseline-v5120-migration.md) §5). Das WIP-Limit 1 hält
  sie ohnehin nacheinander. **Rollenfolge:** vor `in-progress` steht das Architect-Verdikt zu den Fragen unten (liegt vor), je Slice
  ein frischer Kontext für Implementer, Reviewer und Verifier.
- **Die zwölf Fragen an den Architect — beantwortet** im
  [Architect-Verdikt vom 2026-11-26](../../reviews/2026-11-26-welle-98-architect-verdikt.md) (dort je Frage Begründung, kleinste
  Antwort und die größere Antwort mit ihrem Mehraufwand). Die Nummern sind die der Fragen, wie sie hier standen; die Antworten:
  1. *Koexistenz-Adaption:* **ein eigener Eintrag**, geschrieben im Lauf von [`slice-298`](done/slice-298-neue-namen-wave-archive-anker.md),
     Auflösungs-Trigger **der Bestand ist weg** (Kommando im Verdikt), nicht *permanent*; die Senkung steht in
     [`ADR-0117`](../adr/0117-wellen-namen-ueber-die-konfiguration.md), der Eintrag zeigt dorthin.
  2. *`wave-98`:* **ja, aber erst in `slice-298`**, zusammen mit der Konfiguration (Auftraggeber, 2026-11-26): reiner Move, dann der
     Nachzug aller Verweise; keine Lücke für `planning-check`. Bis dahin heißt die Welle `welle-98`.
  3. *Sensor-Reichweite:* **die Konfiguration**; 38 Register-Zeilen und 38 Ergebnisnotizen sind danach ungewächtert — benannt in
     [`ADR-0117`](../adr/0117-wellen-namen-ueber-die-konfiguration.md) (Senkung nach [`AGENTS.md`](../../../AGENTS.md) §3.6).
     Größere Antwort: das Modul liest mehrere Präfixe — Produkt-Merkmal des Auftraggebers, nicht in dieser Welle.
  4. *Commit-Zerlegung:* **ein Commit** für Rename und Link-Nachzüge, die Spezifikation selbst im zweiten; **zulässig ohne
     Ausnahme** — [`AGENTS.md`](../../../AGENTS.md) §3.3 schützt die bewegte Datei, sie bleibt `R100`; beide Commits in **einen** Push.
  5. *Eingefrorene Link-Ziele:* Spezifikation → **Nachzug** der Link-Ziele (64 ADRs, 39 Dateien), dazu **ein**
     `codepaths.ignore-refs`-Eintrag ([`ADR-0025`](../adr/0025-codepaths-ignore-refs.md)); alter Baum (7 Dateien; Nachtrag V1 des Verdikts: fünf `Accepted`-ADRs → **Ventil**, `slice-274` und das Verdikt zu welle-93 → Retarget)
     (fünf datei-skopierte `links.ignore-refs`-Einträge, [`MR-069`](../../../harness/conventions.md#mr-069)).
  6. *Lastenheft Z. 4762:* **nachziehen** (Auftraggeber, 2026-11-26) — Version, eine Historie-Zeile
     ([`MR-032`](../../../harness/conventions.md#mr-032)); Plan-Auftrag in [`slice-297`](done/slice-297-spezifikation-heisst-specification.md).
  7. *Kommentarzeilen:* **nur was funktioniert**; die 94 Kommentarzeilen bleiben als benanntes Negativ
     ([`AGENTS.md`](../../../AGENTS.md) §3.7, Grund im Verdikt); der Nachzug ist die größere Antwort (ein Commit, 36 Dateien).
  8. *`archive.zip`:* **alle drei Modi** (Welle folgt dem Präfix der Kennung; Slice und Review schreiben `-archive.zip`).
  9. *Alter Baum fällt im Bump:* **Praxis bleibt**; die Abweichung vom Kanon trägt [`MR-021`](../../../harness/conventions.md#mr-021) (1)
     schon, der Bump-Eintrag zeigt dorthin und nennt den Vor-Bump-Commit.
  10. *Freshness-Audit der Adaptionen:* **nein**; die Obergrenze von 17 Treffern bleibt Messung.
  11. *Ort der Entscheidungsvorlage:* **`baseline-v700-neuerungen.md`** unter `docs/plan/planning/`, flach; kein Eintrag für die Klasse.
  12. *`stub.go`, zwei Kommentarzeilen:* **im Bump nachziehen**, in [`slice-295`](done/slice-295-baseline-v700-bump.md).
- **Review-Vorgabe der Welle** (Lehre aus welle-96 und welle-97): je Slice **höchstens zwei Runden**,
  eingearbeitet werden nur HIGH und MEDIUM; LOW und INFO bleiben im Report. Jede Nachbesserung zitiert die
  Stelle, die sie trägt, statt sie nachzuerzählen
  ([`AGENTS.md`](../../../AGENTS.md) §5 Regel 18).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Die Umbenennung des Bestands.** Abgeschlossene `welle-<NN>`-Dateien (344 Pfade), `archiv.zip` (188), die 24
  `seit welle-`-Anker, die 38 Register-Zeilen: nichts davon wird umbenannt (Entscheid 1, „Bestand bleibt“). Eine spätere
  Migration ist ein eigener Beschluss des Auftraggebers.
- **Ein Produkt-Merkmal „mehrere Wellen-Präfixe“** (Frage 3, größere Antwort) und jede Änderung der **Defaults** des
  Werkzeugs (`welle-*.md`, `welle-*-results.md`): Konsumenten verlassen sich auf sie. Die Welle ändert die
  Konfiguration **dieses** Repos.
- **Die Übernahme jeder weiteren Neuerung** (Entscheid 2) — auch der kleinen und der mechanischen: das Vorlagen-Paar
  `environment-check`, die Worktree-Regel samt `.claude/worktrees/**` im Scan-Ignore, die Bedrohungs-Ausgänge, die
  Umgebungs-Zusage, der `waves`-Block der Vorlagen-`.d-check.yml` über die zwei Globs hinaus. Die Sichtung klassifiziert
  und empfiehlt; ein Haken „übernommen“ gehört in einen Folge-Slice nach der Antwort.
- **Eine Änderung am Produkt über §3 (d) hinaus** — Standardwerte, Grund-Codes, Spezifikation (außer dem Selbst-Link-
  Nachzug und der Pin-Nennung), Handbuch (außer der einen Zeile) und der Wortlaut von [`DC-FA-CLI-013`](../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen).
- **Ein Durchgang durch die Adaptions-Liste** (Freshness-Audit) — Frage 10 in §5; ohne Antwort bleibt es bei der
  Obergrenze von 17 Treffern als Messung.
- **Ein CR an den Kurs.** Findet die Sichtung einen Widerspruch im Regelwerk selbst, ist ein ausgehender
  Change Request ([`MR-035`](../../../harness/conventions.md#mr-035)) ein eigener Vorgang.
- **Die Tool-Pins** (Go, golangci-lint, semgrep, a-check) und alle Dependabot-Hebungen — eigene Vorgänge mit
  eigener Messung, wie in [`slice-277`](done/wellenlos/slice-277-baseline-v6181-bump.md).
- **Nacharbeit am eingefrorenen Bestand über die Link-Ziele hinaus** — `done/`, `docs/reviews/`, `CHANGELOG.md`,
  `conventions/done/` und der Kern der `Accepted`-ADRs nennen den alten Pin und den alten Dateinamen als
  Vergangenheits-Aussage und bleiben byte-stabil; Ausnahmen entscheidet allein die Antwort auf Frage 5 in §5.
- **Release, Tag, Release-Prep** (CHANGELOG, READMEs, Handbuch-Kopf, Digest-Pins) — eigener Vorgang nach der
  Welle ([`AGENTS.md`](../../../AGENTS.md) §5 Zeile 17); der CHANGELOG-Hinweis zu (d3) gehört dorthin.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `wave-98-results.md` (Geschwister im Ruheort `done/`, entsteht bei der Closure)
Zähler: [Beobachtungs-Register](observations/README.md)
