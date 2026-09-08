# Review-Report: slice-218, Runde 3 — 2026-09-08

**Review-Art:** Code-Review über den Produkt-Fix, die drei Regressionstests, den
neuen Spezifikations-Absatz, `CO-001`, den Folge-Slice und die drei Doku-Orte.
Geprüft wird der Diff gegen Plan, Hard Rules und die Zusage, die der Fix
einlöst — mit dem vom Plan **vorab deklarierten Fokus dieser Runde**: der
*Vollständigkeit der Klasse*. Nichts wurde geglaubt; jede tragende Aussage ist
nachgefahren.

**Gegenstand:** slice-218 · Commit `7d25d2f1` (Vor-Stand `d308a882`), dazu der
Gesamtstand des Slice `e20b7107..7d25d2f1`. Geändert im jüngsten Commit:
`internal/adapter/driven/git/git.go`, `internal/adapter/driven/git/git_test.go`,
`internal/hexagon/core/rules/vcs.go`, `internal/hexagon/core/rules/vcs_test.go`,
`spec/spezifikation.md`, `docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md`,
`docs/user/benutzerhandbuch.md`, `harness/sensors/adr-check.md`, der Slice-Plan
und der neue Folge-Slice slice-219.

**Skill:** `.harness/skills/reviewer.md` @ 1.16.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-09-08

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als Tag plus Pfad in Inline-Code
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>) statt als Link.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan slice-218 nach der **zweiten** Plan-Änderung: §1 (vier Ausschlüsse,
  zwei deklarierte Öffnungen, beide `MR-066`-Pflichten), §2 DoD (1)–(5), §3–§6, §8
- `DC-FA-VCS-001` (die Zusage, die der Fix einlöst), `DC-FA-COMMITS-001`,
  `DC-FA-TRK-001`, `DC-QA-02`, `DC-QA-03`
- `spec/spezifikation.md` §`DC-FA-VCS-001.a` Schritte 1–5 samt neuer
  Historie-Zeile
- `AGENTS.md` §3.1, §3.2, §3.4, §3.6, §3.7, §3.8, §4, §5, §6 Schritt 4
- `MR-006` (Geltungsbereich **und** Scope-Grenze C-4 gelesen), `MR-066`
- `v6.5.0` · `regelwerk/grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP),
  8×8-Matrix · `regelwerk/modul-07-carveouts.md` §Ziel-Form: Carveout ·
  `regelwerk/modul-13-quality-gates.md` §Kernidee
- Vorherige Findings am gleichen Modul: Reports `…-review-r1.md` (F-1 … F-4) und
  `…-review-r2.md` (F-1 … F-12) — insbesondere R2/F-1
  `fix-schliesst-pfad-nicht-klasse`, R2/F-11
  `entstehungsgeschichte-im-lebenden-artefakt`, R2/F-12(b) („der befundfreie
  Zweig von `FileAt` ist über `CheckVCS` gar nicht erreichbar")
- go-git `v5.19.2` im Modul-Cache des `deps`-Images gelesen:
  `plumbing/filemode/filemode.go`, `plumbing/object/tree.go` (`File`,
  `FindEntry`, `dir`, `entry`, `Decode`, `canonicalTreeMode`, `TreeWalker.Next`),
  `plumbing/object/file.go` (`FileIter.Next`), `plumbing/object/treenoder.go`,
  `plumbing/object/difftree.go`, `utils/merkletrie/{difftree,iter,change}.go`,
  `utils/merkletrie/noder/path.go`

**Eigene Läufe und Messungen** (echte Ausgaben, gekürzt auf die tragenden
Zeilen; alle Probe-Repos außerhalb des Repos, alle Läufe `--network none`,
`:ro`-Mount, alle mit derselben `vcs`-Konfiguration wie das Repo):

```text
# 0) Vergleichs-Image
make build                        -> d-check:latest (sha256:49a816ad…, Stand 7d25d2f1)
Instrumentierte Kopie             -> dbg:latest (druckt die uebersetzte Aenderungsliste)
git archive 6e0b7960 / d308a882   -> drei Kopien mit je EINEM zurueckgenommenen Fix

# 1) Kontrollen (kanonischer Pack, echte Verletzung)
Core-Aenderung        : 1 Datei(en) geprueft, 1 Befund(e)  core-drift-vcs   exit=1
Loeschung             : 1 Datei(en) geprueft, 1 Befund(e)  "geloescht oder umbenannt" exit=1
Verzeichnis-Umbenennung: 2 Datei(en) geprueft, 1 Befund(e) "geloescht oder umbenannt" exit=1
--staged Loeschung    : 1 Datei(en) geprueft, 1 Befund(e)                   exit=1

# 2) BASE-/HEAD-Objekte im unsichtbar benannten Pack — fail-closed (gedeckt)
M, BASE-Blob unsichtbar   : nicht lesbares Objekt zu "<pfad>" in "<ref>": file not found  exit=2
M, HEAD-Blob unsichtbar   : nicht lesbares Objekt zu "<pfad>" in "HEAD": file not found   exit=2
D, BASE-Blob unsichtbar   : nicht lesbares Objekt zu "<pfad>" in "HEAD~1": file not found exit=2
BASE-Tree docs/plan/adr   : nicht lesbarer Tree-Eintrag zu "<pfad>" in "HEAD~1": object not found exit=2
BASE-Tree docs/plan       : dito                                                          exit=2
BASE-Tree docs            : dito                                                          exit=2
HEAD-Root-Tree unsichtbar : Range-Spitze "HEAD" nicht aufloesbar: object not found         exit=2
BASE-Commit unsichtbar    : Range-Basis "<sha>" nicht aufloesbar: reference not found      exit=2
--staged, HEAD-Blob       : HEAD-Tree nicht lesbar: object not found                       exit=2
alle Objekte im Pack      : staged-Basis "HEAD" nicht aufloesbar: reference not found      exit=2

# 3) >>> F-1: dieselbe Ursache, STILL GRUEN — in beiden Modi
Verzeichnis der ADRs umbenannt, BASE-Tree des Verzeichnisses unsichtbar:
                          : 2 Datei(en) geprueft, 0 Befund(e)                              exit=0
   (git diff dazu: "D docs/plan/adr/0001-a.md" / "A docs/beschluesse/0001-a.md")
ADR-Verzeichnis geloescht, BASE-Tree unsichtbar:
                          : 1 Datei(en) geprueft, 0 Befund(e)                              exit=0
--staged, dieselbe Loeschung, HEAD-Tree unsichtbar:
                          : 1 Datei(en) geprueft, 0 Befund(e)                              exit=0
instrumentiert            : DBG base.plan.Tree(adr) err=directory not found
                            DBG diffTrees: n=1 err=<nil>
                            DBG change A README.md        <- die Loeschung fehlt ganz

# 4) >>> F-4: unlesbarer HEAD-Stand kommt als Loeschung an
HEAD-Tree docs/plan/adr unsichtbar (Datei unveraendert):
                          : 1 Befund core-drift-vcs "immutable Datei geloescht oder
                            umbenannt"                                                     exit=1

# 5) >>> F-5: Gitlink ersetzt die immutable ADR am geschuetzten Pfad
                          : 0 Datei(en) geprueft, 0 Befund(e)                              exit=0
   Gegenprobe Symlink     : 1 Befund core-drift-vcs                                        exit=1

# 6) >>> F-11: --staged mit unlesbarem INDEX-Blob
                          : error: object not found        (ohne Pfad, ohne Ref)           exit=2

# 7) Regressionstests, je EINER zurueckgenommen (docker build --target test)
Fix 1 zurueck (6e0b7960-Adapter) : --- FAIL: TestFileAtUnlesbaresObjekt      (nur dieser)
Fix 2 zurueck (Modus-Pruefung)   : --- FAIL: TestFileAtEintragOhneBlob       (nur dieser)
Fix 3 zurueck (Added-Zweig)      : --- FAIL: TestVCSAddedMeldetUnlesbareBasis (nur dieser)

# 8) Repo-Gate nachgefahren
make doc-check            : 759 Datei(en) geprueft, 0 Befund(e)                            exit=0
```

Reproduktion von (3), vollständig — der realistischste Fall, eine
Verzeichnis-Umbenennung:

```bash
git init -q .; mkdir -p docs/plan/adr
printf '# ADR-0001\n\n**Status:** Accepted\n\nKern.\n' > docs/plan/adr/0001-a.md
printf 'x\n' > README.md; git add -A; git commit -qm base
TA=$(git rev-parse HEAD:docs/plan/adr)
mkdir -p docs/beschluesse
git mv docs/plan/adr/0001-a.md docs/beschluesse/0001-a.md
git commit -qm "Verzeichnis der immutablen ADRs umbenannt"
echo "$TA" | git pack-objects --quiet .git/objects/pack/pack
git prune-packed
for f in .git/objects/pack/pack-*; do mv "$f" "${f/\/pack-/\/loose-}"; done
docker run --rm --network none -v "$PWD":/repo:ro d-check:latest --range HEAD~1..HEAD
```

---

## Findings

### F-1 — Der stille Pfad ist weiterhin offen: eine **umbenannte oder gelöschte** immutable ADR verschwindet befundfrei, in **beiden** Modi

- `kategorie`: HIGH
- `quelle`: `DC-FA-VCS-001` (*„**gelöschte** oder **umbenannte** immutable Datei
  → Befund `core-drift-vcs`"* — ohne Vorbehalt gegen einen Repo-Zustand);
  `AGENTS.md` §4 (Harness-Lüge)
- `pfad`: `internal/hexagon/core/rules/vcs.go:32-55` (die Fallunterscheidung
  greift erst auf der Änderungsliste); go-git `v5.19.2`
  `plumbing/object/tree.go`, `TreeWalker.Next()`
- `befund`: Liegt das **Tree-Objekt des geschützten Verzeichnisses** im
  unsichtbar benannten Pack und hat es auf der Gegenseite **kein Pendant**
  (Verzeichnis umbenannt oder gelöscht), erscheint die Löschung in der
  Änderungsliste **gar nicht** — instrumentiert gemessen: `diffTrees: n=1
  err=<nil>`, und der einzige Eintrag ist die unbeteiligte neue Datei. Der Lauf
  meldet `0 Befund(e)`, Exit 0, wo derselbe Stand mit kanonischem Pack-Namen
  `1 Befund(e)`, Exit 1 liefert; `git diff` zeigt `D
  docs/plan/adr/0001-a.md`. Dasselbe im `--staged`-Modus, also im
  `pre-commit`-Hook. Die Ursache ist eine **zweite, nicht dokumentierte**
  Zusammenfassung in go-git: `TreeWalker.Next()` macht aus einem nicht ladbaren
  Unterbaum `io.EOF` (`if entry.Mode == filemode.Dir { obj, err = GetTree(…) }`
  … `if err != nil { err = io.EOF; return }`), womit ein unlesbarer Unterbaum
  als *zu Ende iteriert* gilt statt als Fehler. Der Fix aus Runde 2 setzt
  **hinter** dieser Stelle an (er prüft die drei Diff-Zustände `A`/`M`/`D`) und
  kann eine Änderung nicht sehen, die nie entsteht. Damit gilt die in diesem
  Diff geschriebene Zusage *„**Ab der Version nach `v0.75.0` ist das behoben**,
  in beiden Ausprägungen und in beiden Modi"* nicht.
- `verifizierbar`: ja — Messung (3) oben, Reproduktions-Rezept vollständig;
  Kontrolle mit kanonischem Pack-Namen liefert Exit 1
- `klasse`: `fix-schliesst-pfad-nicht-klasse` (**dritte** Instanz: R1/F-1 Blob,
  R2/F-1 Tree-mit-Pendant, hier Tree-ohne-Pendant)

### F-2 — Der genannte Mechanismus erklärt die Tree-Ausprägung nicht; es ist der **zweite** falsche Mechanismus, und §6 hatte davor gewarnt

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-VCS-001.a` (Spec-Treue); Slice-Plan §6 Risiko 1 (*„Der
  Mechanismus ist erschlossen, nicht gelesen … Die erste Erklärung war bereits
  falsch"*)
- `pfad`: `spec/spezifikation.md:3318` (Historie-Zeile),
  `spec/spezifikation.md:1782-1791`, `harness/sensors/adr-check.md:60-76`,
  `docs/user/benutzerhandbuch.md:1601-1615`
- `befund`: Alle vier Stellen führen **eine** Ursache: *„ein unlesbares Objekt
  wird mit demselben Fehler gemeldet wie eine im Tree fehlende Datei"* — und
  ordnen ihr **beide** Ausprägungen zu, auch die, bei der die Datei als `A`
  ankommt. Für die Tree-Ausprägung stimmt das nicht: dort entscheidet nicht die
  `ErrFileNotFound`-Zusammenfassung in `Tree.File`, sondern die
  `io.EOF`-Zusammenfassung in `TreeWalker.Next()` (gelesen, `v5.19.2`), über die
  ein unlesbarer Unterbaum als *leer* gelesen wird — genau die Stelle, die F-1
  offen lässt. Wer die Doku-Erklärung für den Mechanismus hält, sucht die
  Restlücke an der falschen Stelle.
- `verifizierbar`: ja — go-git-Quelle im `deps`-Image plus Messung (3): die
  Änderungsliste ist *fehlerfrei und unvollständig*, nicht *fehlerhaft gefüllt*
- `klasse`: `mechanismus-erschlossen-statt-gelesen`

### F-3 — `spec/spezifikation.md` verweist abwärts auf einen Carveout; die kanonische Matrix führt diese Kante als ❌, und `MR-006` verbietet sie ausdrücklich

- `kategorie`: MEDIUM
- `quelle`: `MR-006` (*Geltungsbereich* `spec/*.md`; *Adaption*: „**kein
  Spec-Stratum (Rang 1–3) verweist abwärts auf ADRs oder Planning-Artefakte**");
  `v6.5.0` · `regelwerk/grundlagen-referenz-richtung.md`
  §Referenz-Richtung (SDP), Zeile *Technik* × Spalte *Carveout* = ❌ nebst
  Decken-Regel
- `pfad`: `spec/spezifikation.md:3318` (`[CO-001](../docs/plan/carveouts/…)`)
- `befund`: Die neue Historie-Zeile verlinkt aus dem Technik-Stratum in die
  Planungs-Ebene. `AGENTS.md` §3.4 zählt fünf Kategorien auf, unter denen
  *Carveout* nicht steht — die höherrangigen Quellen (Baseline-Matrix, `MR-006`)
  sind aber strenger, und `MR-006` benennt zugleich, warum das niemandem
  auffällt: die Kante ist *„bewusst unbewacht"*, weil `matrix` Carveouts nicht
  als Klasse führt (Scope-Grenze C-4). Der Slice brauchte den Verweis nicht: §1
  Punkt 3 begründet nur den **Spec-Absatz**, nicht den Carveout-Zeiger.
  Einordnung, damit der Befund nicht zu breit gelesen wird: dieselbe Tabelle
  nennt in der Zeile 2026-09-02 bereits einen CR-Pfad — aber als **Inline-Code**
  ohne Link, und *CR* ist in der 8×8-Matrix überhaupt keine Zeile, *Carveout*
  dagegen eine mit ❌.
- `verifizierbar`: nein für den Gate-Lauf (kein Sensor deckt die Kante — genau
  das sagt `MR-006`); ja gegen die zitierte Matrix-Zelle
- `klasse`: `spec-stratum-verweist-abwaerts-auf-planungsartefakt`

### F-4 — Ein unlesbarer **HEAD**-Stand kommt als Löschung an: `core-drift-vcs`/Exit 1 statt Umgebungsfehler/Exit 2

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-VCS-001.a` Schritt 2 in der neuen Fassung (*„trennt nicht
  vorhanden von nicht lesbar: … das Zweite ist ein Umgebungsfehler (Exit 2,
  **fail-closed**)"*); `AGENTS.md` §5 (Grenzen-Liste)
- `pfad`: `internal/hexagon/core/rules/vcs.go:64-77` (`vcsDeleted` fasst den
  HEAD-Stand nie an); Zusage in `harness/sensors/adr-check.md:44-48` und
  `docs/user/benutzerhandbuch.md:1593-1597` (*„— **alle mit Exit 2**"*)
- `befund`: Liegt das **HEAD**-Tree des geschützten Verzeichnisses im
  unsichtbaren Pack, während die Datei unverändert ist, meldet der Lauf
  `1 Befund core-drift-vcs` *„immutable Datei gelöscht oder umbenannt"* mit
  **Exit 1** (gemessen). Die Datei wurde nicht angefasst; die Diagnose schickt
  den Leser auf die Suche nach einer Umbenennung, die es nicht gibt, und die
  dokumentierte Abhilfe `git repack -A -d` ist von dieser Meldung aus nicht
  erreichbar. Die in diesem Diff geschriebene Symptom-Liste schließt den Fall
  ausdrücklich aus (*alle mit Exit 2*). Das Verhalten ist **nicht** neu — es
  stammt aus der Zeit vor slice-218 —; neu ist die Zusage, die es widerlegt.
- `verifizierbar`: ja — Messung (4)
- `klasse`: `umgebungsfehler-als-inhaltsbefund`

### F-5 — Ein Gitlink, der eine `Accepted`-ADR am geschützten Pfad ersetzt, bleibt befundfrei; die neue Spec-Zusage sanktioniert das, ohne die Folge zu nennen

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-VCS-001` (Kern-Änderung an einer immutablen Datei ⇒ Befund);
  `AGENTS.md` §5 (*wer eine Grenze aufschreibt, prüft sie gegen den Gegenstand*)
- `pfad`: `internal/adapter/driven/git/git.go:271-286`;
  `spec/spezifikation.md:1790-1791` (*„ein Eintrag **ohne** Datei-Inhalt
  (Verzeichnis, Gitlink) ist dabei kein unlesbares Objekt, sondern bleibt
  befundfrei"*)
- `befund`: Wird `docs/plan/adr/0001-a.md` im HEAD durch einen Gitlink
  **desselben Pfades** ersetzt, meldet der Lauf `0 Befund(e)`, Exit 0
  (gemessen) — der Core der immutablen Datei ist an diesem Pfad verschwunden.
  Die Gegenprobe mit einem Symlink liefert korrekt einen Befund, die Klasse ist
  also asymmetrisch. Der Zustand ist der des Vor-Fix-Stands und insoweit keine
  Regression; die Spec-Zeile erhebt ihn aber jetzt zur **Zusage** und nennt nur
  den entlastenden Fall (der wandernde Submodul-Zeiger), nicht den belastenden
  (der Submodul-Zeiger, der eine geschützte Datei verdrängt). Zugleich
  überholt der Fall R2/F-12(b): der befundfreie Zweig von `FileAt` ist über
  `CheckVCS` sehr wohl erreichbar.
- `verifizierbar`: ja — Messung (5)
- `klasse`: `grenzen-liste-ohne-ihre-belastende-haelfte`

### F-6 — Zwei Index-Stellen beschränken `CO-001` auf den `RANGE=`-Modus, während der Carveout selbst „in beiden Modi" führt

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §4 (Sensors-Tabelle behauptet nicht weniger, als gilt);
  `v6.5.0` · `regelwerk/modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin)
- `pfad`: `harness/README.md:88` (Bindung-Spalte `make adr-check`:
  *„RANGE=-Modus im publizierten Image blind"*),
  `docs/plan/carveouts/README.md:12` (Titel *„`vcs --range` überspringt still"*),
  `harness/sensors/trace-check.md:32` (*„Bis slice-218 meldete es im
  `RANGE=`-Modus still grün"*)
- `befund`: `CO-001` führt als betroffenes Gate ausdrücklich *„`adr-check`
  (Modul `vcs`) — **in beiden Modi**"*, und Runde 2 hat den `--staged`-Fall
  gemessen. Die drei Stellen, über die ein Leser den Carveout überhaupt findet
  (Sensors-Tabelle, Carveout-Index, Nachbar-Sensor), nennen weiterhin nur den
  Range-Modus. Wer die Sensors-Tabelle liest, schließt daraus, der lokale
  `pre-commit`-Hook sei nicht betroffen — er ist es.
- `verifizierbar`: ja — Textvergleich der vier Stellen gegen `CO-001` §Betroffenes Gate
- `klasse`: `index-zeile-enger-als-das-indizierte`

### F-7 — Der Folge-Slice zieht beim Auflösungs-`git mv` nur zwei von sieben Verweis-Trägern nach; einer davon ist ein Spec-Stratum

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-LINK-001` (Link-Ziel muss existieren); `AGENTS.md` §3.3
  (Lifecycle-Move bündelt die gekoppelten Verweise), `MR-013`
- `pfad`: `docs/plan/planning/open/slice-219-release-loest-co-001.md` DoD (3)
- `befund`: **Zählform vorab:** gezählt sind Markdown-Links, deren Ziel auf die
  Datei `docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md` auflöst, in
  lebenden Dokumenten (Review-Reports ausgenommen, weil sie einfrieren). Das
  sind **acht Vorkommen in sieben Dateien**: `harness/README.md` (1),
  `docs/plan/carveouts/README.md` (1), `spec/spezifikation.md` (1),
  `harness/sensors/adr-check.md` (1), `harness/sensors/trace-check.md` (1),
  slice-218 (1), slice-219 (2). DoD (3) nennt für den `git mv` nach
  `carveouts/done/` genau zwei davon (Carveout-Index und Sensors-Zeile). Die
  übrigen fünf werden mit dem Move zu toten Links; `make doc-check` läuft rot,
  darunter auf einem Spec-Stratum — und das ist derselbe Commit, den `MR-013`
  als gebündelt verlangt.
- `verifizierbar`: ja — `make doc-check` nach dem in DoD (3) beschriebenen `git mv`
- `klasse`: `folge-slice-kennt-nur-einen-teil-der-kopplung`

### F-8 — `entryUnreadable` zählt die Datei-Modi von Hand statt über go-gits `IsFile()`; die Deckung hängt an einem undokumentierten Verhalten des gepinnten go-git

- `kategorie`: LOW
- `quelle`: Maintainability (latente Wartungsfalle); Slice-Plan §8, dritter
  Vorprüfungs-Block (*„Der Gegenstand ist an `go-git v5.19.2` gebunden, und ein
  Dependabot-Bump dieser Abhängigkeit könnte das gemessene Verhalten ändern"*)
- `pfad`: `internal/adapter/driven/git/git.go:279-280`
- `befund`: Der Ausdruck `e.Mode == filemode.Regular || e.Mode ==
  filemode.Executable || e.Mode == filemode.Symlink` lässt `filemode.Deprecated`
  (`0100664`) aus, obwohl das ein Modus **mit** Datei-Inhalt ist; go-git führt
  dafür das Prädikat `FileMode.IsFile()`, das ihn einschließt. Heute ist die
  Lücke **nicht erreichbar** — gemessen an einem Probe-Repo mit einem von Hand
  auf `100664` gesetzten Tree-Eintrag (Exit 2, korrekt), und gelesen in
  `Tree.Decode`, das seit `v5.19.2` jeden Modus über `canonicalTreeMode` auf
  {Dir, Regular, Executable, Symlink, Submodule} normalisiert. Genau diese
  Normalisierung ist aber nirgends als tragende Annahme benannt: fällt sie mit
  einem Bump weg, liefert dieser Zweig für alte `0100664`-Bestände wieder ein
  stilles Grün, und zwar ohne dass ein Test darauf zeigt.
- `verifizierbar`: ja — go-git-Quelle im `deps`-Image; Probe mit `git mktree
  '100664 blob …'`
- `klasse`: `handzaehlung-statt-bibliotheks-praedikat`

### F-9 — `--staged` mit unlesbarem **Index**-Blob meldet ein kontextloses `error: object not found`

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §5 (Grenzen-Liste); Symptom-Liste in
  `harness/sensors/adr-check.md:40-48` und `docs/user/benutzerhandbuch.md:1593-1597`
- `pfad`: `internal/adapter/driven/git/git.go:315-317` (`fileFromIndex` reicht
  den Storer-Fehler ungewrappt durch)
- `befund`: Wandert der **gestagte** Blob in den unsichtbar benannten Pack,
  bricht der Lauf mit `error: object not found` ab — ohne Pfad, ohne Ref, ohne
  Bezug zum Pack-Namen (gemessen). Die drei anderen Meldeformen desselben
  Auslösers sind ausführlich dokumentiert und tragen den Satz *„Wer nur eine der
  Formen kennt, erkennt die anderen nicht wieder"*; diese vierte steht in keiner
  der Listen und ist die einzige, aus der sich der Auslöser gar nicht erschließen
  lässt.
- `verifizierbar`: ja — Messung (6)
- `klasse`: `meldeform-ohne-kontext`

### F-10 — `harness/sensors/adr-check.md` sagt „Zwei Regressionstests" und nennt drei

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §5 (Zählmethode/Aussage gegen den Gegenstand)
- `pfad`: `harness/sensors/adr-check.md:72-75`
- `befund`: Der Satz *„**Zwei** Regressionstests halten beide Richtungen und
  fallen ohne ihren Fix"* wird von einer Klammer mit **drei** Testnamen
  fortgesetzt. Alle drei existieren und fallen tatsächlich je ohne ihren
  eigenen Fix (unten als Negativbefund gemessen) — falsch ist nur die Zahl.
- `verifizierbar`: ja — Textvergleich; Messung (7)
- `klasse`: `zahl-widerspricht-der-liste-daneben`

### F-11 — Vier Plan-Sektionen beschreiben weiterhin den Doku-Slice von vor den zwei Plan-Änderungen

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §6 Schritt 4 (der Plan bindet den Lauf); `v6.5.0` ·
  `regelwerk/modul-05-planning-harness.md` §Zwei Schritte vor der
  Modus-Begründung
- `pfad`: Slice-Plan §3 (Zeile 195 ff.), §4 (Zeile 216), §5 (Zeile 224),
  §8 (Zeile 271)
- `befund`: §3 *Plan (vor Code)* führt weiterhin nur drei Doku-Schritte und
  keinen Code-Schritt. §8 sagt *„Der Slice ändert drei Dokumente"* und *„der
  Gegenstand ist zwar Produkt-Verhalten, aber es wird **gemessen, nicht
  geändert**"* — der Slice ändert inzwischen zwei Produkt-Dateien in zwei
  Schichten (driven-Adapter **und** Kern-Regel). §5 nennt als Closure-Trigger
  *„DoD (1) bis (3)"*, obwohl es (4) und (5) gibt; §4 stützt seine
  Rückführungs-Bedingung auf *„§1 Punkt 2 schließt Code aus"*, was Punkt 2 nicht
  sagt. Die letzten beiden standen als R2/F-7 und R2/F-8 schon im Report der
  Vorrunde und sind unverändert.
- `verifizierbar`: nein (Urteil), aber jede Einzelaussage ist am Diff prüfbar
- `klasse`: `plan-sektionen-nicht-mit-der-plan-aenderung-nachgezogen`
  (**dritte bis sechste** Instanz in dieser Sitzung ⇒ Steering-Loop-Signal,
  gehört in die Closure-Notiz)

### F-12 — Die Entstehungsgeschichte im Carveout ist nicht zurückgenommen, sondern gewachsen

- `kategorie`: INFO
- `quelle`: Maintainability; R2/F-11 (unverändert offen). *Ausdrücklich **kein**
  §3.7-Verstoß* — dessen Geltungsbereich sind Kommentare in
  Code/Konfiguration/Skripten und Zustandsfelder, nicht Fließtext in einem
  Planungs-Dokument.
- `pfad`: `docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md:18-20`, `:22-24`,
  `:52-54`, `:92-95`
- `befund`: Zu der von Runde 2 benannten Stelle (*„Eine frühere Fassung dieses
  Carveouts behauptete das Gegenteil"*) sind in diesem Commit zwei weitere
  gekommen (*„Die zweite Zeile fand erst die zweite Review-Runde …"*, *„Die
  erste Fassung dieses Carveouts kannte nur die Blob-Hälfte und hätte den
  Trigger für erfüllt gehalten"*) und eine dritte im Feld `Folge-Slice`. Ein
  Konsument, der den Carveout liest, um zu handeln, braucht keine davon; die
  Vorlage hat für Vorgangs-Historie den Abschnitt `## Geschichte`, und die
  Commit-Botschaft trägt dieselbe Information bereits.
- `verifizierbar`: nein (Urteil)
- `klasse`: `entstehungsgeschichte-im-lebenden-artefakt` (**zweite** Instanz)

---

## Negativbefunde (geprüft, ohne Befund)

- **Der Fix wirkt auf allen Achsen, die er beansprucht.** Zehn Proben mit
  unsichtbar benanntem Pack (Messung 2): BASE-Blob bei `M` und `D`, HEAD-Blob
  bei `M`, BASE-Tree auf drei Verzeichnis-Ebenen, HEAD-Root-Tree,
  BASE-Commit-Objekt, `--staged` mit unlesbarem HEAD-Blob, alle Objekte im Pack
  — **jede** liefert Exit 2 mit einer benannten Meldung. Der `D`-Zweig ist auf
  der BASE-Seite gedeckt; der `M`-Zweig auch auf der HEAD-Seite, weil die
  Fallunterscheidung in `FileAt` sitzt und nicht in einem ref-spezifischen
  Zweig.
- **Die drei Regressionstests fallen je ohne ihren eigenen Fix — und nur dann.**
  Drei Kopien, in jeder genau ein Fix originalgetreu zurückgenommen (der
  Adapter-Stand aus `6e0b7960` bzw. `d308a882`, die Kern-Regel aus `d308a882`);
  `docker build --target test` je einmal: `TestFileAtUnlesbaresObjekt`,
  `TestFileAtEintragOhneBlob`, `TestVCSAddedMeldetUnlesbareBasis` fallen
  einzeln, alle übrigen Pakete bleiben grün. Alle drei Kopien kompilieren — die
  Rücknahme ist echt, nicht nur ein Build-Fehler.
- **Die wörtlich zitierten Meldungen stimmen byte-genau.** Vier Formen gegen die
  tatsächliche Ausgabe gehalten: `staged-Basis "HEAD" nicht auflösbar:
  reference not found`, `HEAD-Tree nicht lesbar: object not found`, `nicht
  lesbares Objekt zu "<pfad>" in "<ref>": file not found`, `nicht lesbarer
  Tree-Eintrag zu "<pfad>" in "<ref>": object not found`. (Dass die Liste
  unvollständig ist, steht als F-9; die zitierten Formen selbst sind korrekt.)
- **Die Modus-Menge in `entryUnreadable` ist für `v5.19.2` vollständig.**
  `Tree.Decode` normalisiert über `canonicalTreeMode` auf genau fünf Modi;
  `Dir` und `Submodule` sind die beiden ohne Datei-Inhalt, die übrigen drei sind
  genau die aufgezählten. Die Handzählung ist heute korrekt — die Kopplung an
  diese Normalisierung ist als F-8 notiert, nicht die Deckung.
- **Typänderung Datei → Symlink.** Der Symlink-Blob wird gelesen, der
  Core-Vergleich schlägt an: `1 Befund core-drift-vcs`, Exit 1. Kein stiller
  Pfad. (Die Gegenrichtung Datei → Gitlink ist F-5.)
- **Die Gegenrichtung des Fixes hält.** Eine wirklich neu angelegte ADR bleibt
  befundfrei (Exit 0), auch mit dem neuen `FileAt`-Aufruf im `A`-Zweig; der
  Unit-Test deckt beide Richtungen in einem Testfall.
- **Hexagon-Richtung (ADR-0005/ADR-0012).** Der Kern ruft `vcs.FileAt` über den
  bereits vorhandenen `driven.VCS`-Port; `internal/hexagon/core/rules/vcs.go`
  bekommt **keinen** neuen Import. Der Adapter-Teil bleibt im driven-Adapter,
  sein einziger neuer Import ist `plumbing/filemode` — go-git ist dort bereits
  die deklarierte einzige Tür.
- **`AGENTS.md` §3.2/§3.6.** Keine `//nolint`, keine gesenkte Schwelle, keine
  Gate-Lockerung; `.golangci.yml` und die Schwellen sind unberührt.
- **`AGENTS.md` §3.7, Kommentar-Klassen.** Die geänderten Kommentare (`FileAt`,
  `entryUnreadable`, der `A`-Zweig, `fileErr`, die zwei Test-Doc-Blöcke) tragen
  *Zusage*, *Abgrenzung* und *Grenze*; keine Review-Historie, keine
  Slice-Nummer, kein Mess-Label. Herkunft steht je einmal und auflösbar
  (`DC-FA-VCS-001`).
- **`DC-QA-03`.** Kein Netzzugriff, kein Schreiben ins geprüfte Repository; alle
  Proben liefen `--network none` mit `:ro`-Mount, und der neue `FileAt`-Aufruf
  ist rein lesend.
- **`AGENTS.md` §3.1.** Kein Host-Go, kein Host-Paketmanager, kein
  Host-Skript-Interpreter im Diff; die Prüfung dieses Reports lief über `make
  build`, `docker build` und POSIX-Werkzeuge.
- **Abgrenzung eingehalten.** §1 Punkt 1 (kein Repack durch das Produkt), Punkt 2
  (kein Sensor auf Pack-Namen) und Punkt 4 (keine Behebung im publizierten Bild)
  sind über den ganzen Slice-Bereich `e20b7107..7d25d2f1` gehalten; Punkt 3 ist
  wie deklariert nur zur **Spezifikations**-Hälfte geöffnet —
  `spec/lastenheft.md` ist in keinem der sechs Commits angefasst. Die
  Produkt-Code-Öffnung steht vor dem Code im Plan.
- **`MR-066` korrekt zitiert und eingelöst.** Geltungsbereich trifft zu; für die
  **zweite** Plan-Änderung stehen erneut Grund und vorab benannte Ersatz-Form
  (diese dritte Runde mit deklariertem Fokus). Der Fokus hat getragen: F-1 und
  F-5 stammen beide aus ihm.
- **Carveout-Form und Folge-Slice.** Alle Pflichtfelder der Vorlage vorhanden
  (`Status`, `Datum angelegt`/`Letzte Prüfung`, `Betroffenes Gate`,
  `Geltungsbereich`, `Folge-Slice`, `Auflösungs-Trigger`); der
  Auflösungs-Trigger ist beobachtbar, mit Prüf-Kommando unterlegt und verlangt
  ausdrücklich **beide** Ausprägungen. slice-219 existiert in `open/`, trägt
  `Verantwortlich: —` (korrekt für `open/`), eine Abgrenzung mit drei begründeten
  Punkten und eine DoD, die den Carveout überlebt — die R2/F-5-Kritik ist
  eingelöst. (Die Verweis-Kopplung ist F-7, nicht die Form.)
- **Kein Abwärts-Token im Spec-Körper.** Der neue Absatz und die Historie-Zeile
  nennen weder Slice- noch Wellen-Kennung noch Commit-Hash noch einen
  Modul-Pfad; `make doc-check` (Modul `matrix`) ist grün. Der Carveout-**Link**
  ist F-3 — ihn deckt keine `matrix`-Klasse.
- **Repo-Gate nachgefahren.** `make doc-check`: `759 Datei(en) geprüft, 0
  Befund(e)`, Exit 0 — deckungsgleich mit der Zahl in der Commit-Botschaft. Die
  Packs dieses Repos tragen den kanonischen Namen; F-1 wirkt hier heute nicht.
- **`AGENTS.md` §5, Reichweite der Botschaft.** Die Kontroll-Matrix der
  Commit-Botschaft (fünf Fälle) ist nachgefahren und stimmt in allen fünf
  Zeilen. Überdehnt ist genau die Zusage über die *Klasse* (F-1/F-2), nicht die
  Messung.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 1 | F-1 |
| MEDIUM | 6 | F-2, F-3, F-4, F-5, F-6, F-7 |
| LOW | 4 | F-8, F-9, F-10, F-11 |
| INFO | 1 | F-12 |

**Wiederholungs-Signal.** F-1 ist die **dritte** Instanz derselben Klasse
(R1/F-1 Blob → R2/F-1 Tree mit Pendant → hier Tree ohne Pendant): *ein Fix
schließt den gemessenen Pfad, nicht die Klasse*. Damit ist die
Steering-Loop-Schwelle erreicht — der Eintrag gehört bei der Closure ins
Beobachtungs-Register und braucht einen Ausgang, nicht nur eine Notiz. F-11
zählt vier weitere Instanzen von *Plan-Sektion nicht mit der Plan-Änderung
nachgezogen* (nach R2/F-7, R2/F-8) und ist damit ebenfalls über der Schwelle.
Beide Klassen haben denselben Kern: Was aus einem einzelnen gemessenen Fall
folgt, wird als Aussage über die Menge geschrieben.

## Verdikt

**Blockiert.** F-1 steht der Closure entgegen, und zwar aus dem Grund, aus dem
der Slice überhaupt existiert: `vcs` meldet weiterhin `0 Befund(e)`/Exit 0,
während eine `Accepted`-ADR gelöscht oder ihr Verzeichnis umbenannt wird — in
**beiden** Modi, also auch im lokalen `pre-commit`-Hook. Solange das steht,
sagen Handbuch, Sensor-Datei, Spezifikation und Carveout einem Adopter mehr zu,
als das Produkt hält, und der Auflösungs-Trigger von `CO-001` würde nach zwei
grünen Proben für erfüllt gehalten, während die dritte Ausprägung offen ist.

Der Befund verlangt **keinen** größeren Fix als die beiden bisherigen — er
verlangt, ihn eine Ebene tiefer anzusetzen: Die Änderungsliste selbst ist die
Stelle, an der ein unlesbarer Unterbaum heute als *leer* durchgeht (F-2 nennt
den gelesenen Mechanismus). Wer nur den `D`-Zweig härtet, schließt wieder einen
Pfad statt der Klasse — das ist genau die Klasse, die mit diesem Befund ihre
dritte Instanz erreicht.

Was der Diff richtig macht, ist belegt und bleibt: die Blob- und die
Tree-mit-Pendant-Ausprägung sind gemessen geschlossen, alle drei
Regressionstests fallen einzeln ohne ihren eigenen Fix, die Gegenrichtung
erzeugt in den geprüften Formen keinen Fehlalarm, und die zitierten Meldungen
stimmen wörtlich. Die sechs MEDIUM sind vor der Closure zu klären; F-3 und F-7
berühren dabei ein Spec-Stratum und sollten nicht in den Folge-Slice wandern.
