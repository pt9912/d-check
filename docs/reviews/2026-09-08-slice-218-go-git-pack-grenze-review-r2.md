# Review-Report: slice-218, Runde 2 — 2026-09-08

**Review-Art:** Code-Review über Produkt-Code, Regressionstest, Carveout und die
umgeschriebene Doku. Geprüft wird der Diff gegen Plan, Hard Rules und die
Zusage, die der Fix einlöst — mit dem vom Plan **vorab deklarierten Fokus**:
die Fallunterscheidung im git-Adapter gegen die **Gegenrichtung**. Nichts wurde
geglaubt; jede tragende Aussage ist nachgefahren.

**Gegenstand:** slice-218 · Commit `d308a882` (Vor-Stand `6e0b7960`). Geändert:
`internal/adapter/driven/git/git.go` (+29), `internal/adapter/driven/git/git_test.go`
(+58), `docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md` (neu, +106),
`docs/plan/carveouts/README.md`, `harness/README.md`,
`harness/sensors/adr-check.md`, `harness/sensors/trace-check.md`,
`docs/user/benutzerhandbuch.md`, der Slice-Plan selbst.

**Skill:** `.harness/skills/reviewer.md` @ 1.16.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-09-08

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als Tag plus Pfad in Inline-Code
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>) statt als Link.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan slice-218 in der Fassung nach der Plan-Änderung: §1 (vier
  Ausschlüsse + die deklarierte Öffnung), §2 DoD (1)–(5), §4 Trigger, §5
  Closure-Trigger, §6 Risiken
- `DC-FA-VCS-001` (die Zusage, die der Fix einlöst), `DC-FA-COMMITS-001`,
  `DC-FA-TRK-001`, `DC-QA-02`, `DC-QA-03`
- `spec/spezifikation.md` §`DC-FA-VCS-001.a` Schritte 1–5 samt der
  Historie-Zeile vom 2026-08-31 (die Rename-Erkennungs-Vorgeschichte)
- `AGENTS.md` §3.1 (Werkzeug-Klasse), §3.2, §3.4, §3.6, §3.7 (Kommentar-Klassen),
  §3.8 (Scan-Achse), §4, §5, §6 Schritt 4 (Abgrenzung bindet den Lauf)
- `MR-066` (Geltungsbereich **und** beide Pflichten gelesen, nicht nur der Titel)
- `v6.5.0` · `regelwerk/modul-07-carveouts.md` §Ziel-Form: Carveout und
  §Werkzeug-Wahl bei Diskrepanz · `regelwerk/modul-13-quality-gates.md`
  §Kernidee und §Hard Rule (Doku-Disziplin)
- Vorherige Findings am gleichen Modul: Report `…-review-r1.md`,
  F-1 `fehlendes-objekt-als-abwesende-datei-gelesen` (HIGH),
  F-3 `probe-misst-anderen-zustand-als-den-genannten-ausloeser` (MEDIUM)
- go-git `v5.19.2` im Modul-Cache des `deps`-Images gelesen:
  `plumbing/object/tree.go` (`File`, `FindEntry`, `dir`, `entry`),
  `plumbing/object/treenoder.go`, `internal/pathutil/tree.go`

**Eigene Läufe und Messungen** (echte Ausgaben, gekürzt auf die tragenden
Zeilen; alle Probe-Repos außerhalb des Repos, alle Läufe `--network none`):

```text
# 0) Vergleichs-Images
make build                        -> d-check:latest  (sha256:145e2a2b…, identisch zum Stand)
git archive 6e0b7960 | tar -x     -> d-check:prefix  (Vor-Fix)
git archive d308a882 | tar -x     -> d-check:mutant  (entryUnreadable liefert immer true)
                                  -> d-check:dbg     (druckt die uebersetzte Aenderungsliste)
                                  -> d-check:revert-test (Fallunterscheidung zurueckgenommen)

# 1) Kern-Fall: unlesbares BASE-BLOB (die Probe des Slice, nachgefahren)
kanonischer Pack : 1 Datei(en) geprueft, 1 Befund(e)  core-drift-vcs   exit=1
loose-<hash>     : d-check:prefix -> 1 Datei(en) geprueft, 0 Befund(e) exit=0
loose-<hash>     : d-check:latest -> error: nicht lesbares Objekt zu
                   "docs/plan/adr/0001-erste.md" in "d04206aa…": file not found  exit=2

# 2) Gegenrichtung, drei Formen (Neuzugang, Rename, neues Unterverzeichnis)
d-check:latest und d-check:prefix melden identisch 1 Befund (der Rename) —
kein Fehlalarm auf 0002-neu.md und unter/0004-tief.md.

# 3) Regressionstest ohne den Fix
docker build --target test  (revert-src)
--- FAIL: TestFileAtUnlesbaresObjekt (0.01s)
    git_test.go:345: unlesbares Objekt still uebersprungen: ok=false len=0 — erwartet war ein Fehler

# 4) unlesbares BASE-TREE-Objekt statt eines Blobs   >>> F-1
kanonischer Pack : 1 Datei(en) geprueft, 1 Befund(e)  core-drift-vcs   exit=1
loose-<hash>     : d-check:latest -> 1 Datei(en) geprueft, 0 Befund(e) exit=0
loose-<hash>     : d-check:prefix -> 1 Datei(en) geprueft, 0 Befund(e) exit=0
instrumentiert   : DBG change status=65 path=docs/plan/adr/0001-erste.md   (65 = 'A')
kanonisch        : DBG change status=77 path=docs/plan/adr/0001-erste.md   (77 = 'M')

# 5) dasselbe im --staged-Modus (der pre-commit-Hook)                 >>> F-1
kanonisch        : 1 Datei(en) geprueft, 1 Befund(e)  core-drift-vcs   exit=1
loose-<hash>     : 1 Datei(en) geprueft, 0 Befund(e)                   exit=0

# 6) --staged mit unlesbarem BLOB (Regressions-Kontrolle)
d-check:latest und d-check:prefix -> error: HEAD-Tree nicht lesbar: object not found  exit=2

# 7) Tree-Eintrag ohne Blob (gitlink) in der geschuetzten Klasse       >>> F-2
d-check:prefix -> 1 Datei(en) geprueft, 0 Befund(e)                    exit=0
d-check:latest -> error: nicht lesbares Objekt zu "docs/plan/adr/sub"
                  in "f97a1216…": file not found                       exit=2

# 8) blinde Gegenrichtungs-Probe                                       >>> F-4
Probe "nur neue ADRs" gegen d-check:mutant (Gegenrichtung absichtlich kaputt):
                  3 Datei(en) geprueft, 0 Befund(e)                    exit=0

# 9) Repo-Gates nachgefahren
make doc-check   -> 757 Datei(en) geprueft, 0 Befund(e)  exit=0
make adr-check   -> 757 Datei(en) geprueft, 0 Befund(e)  exit=0
make trace-check -> 757 Datei(en) geprueft, 0 Befund(e)
```

Reproduktion von (4), vollständig — dieselbe Mechanik, die der Slice für sein
Blob benutzt, nur mit einem **Tree**-Objekt:

```bash
TREE=$(git rev-parse "$BASE:docs/plan/adr")
echo "$TREE" | git pack-objects --quiet .git/objects/pack/pack
git prune-packed
for f in .git/objects/pack/pack-*; do mv "$f" "${f/\/pack-/\/loose-}"; done
docker run --rm --network none -v "$PWD":/repo:ro d-check:latest \
  --enable vcs --disable links --range "$BASE"..HEAD
```

---

## Findings

### F-1 — Der stille Pfad ist nur zur Hälfte zu: ein unsichtbares **Tree**-Objekt liefert weiterhin `0 Befunde`/Exit 0 — in **beiden** Modi

- `kategorie`: **HIGH** (blockierend)
- `quelle`: `DC-FA-VCS-001`; `AGENTS.md` §4 (*„Halluzinierte Gates sind die
  häufigste Form von Harness-Lüge"*); `v6.5.0` · `regelwerk/modul-13-quality-gates.md`
  §Kernidee (*ein Gate, das manchmal rot sein darf, ist keins*)
- `pfad`: `internal/adapter/driven/git/git.go:238-253` (der Fix greift erst,
  wenn `tree.File` überhaupt aufgerufen wird) zusammen mit
  `internal/adapter/driven/git/git.go:98-116` (`diffTrees`) und
  `internal/hexagon/core/rules/vcs.go:33` (`case driven.VCSAdded:` ohne Rumpf)
- `befund`: Liegt statt des BASE-**Blobs** das BASE-**Tree-Objekt** des
  geschützten Verzeichnisses in einem unsichtbar benannten Pack, meldet
  `d-check --enable vcs --range` bei einer echten, mit kanonischem Pack-Namen
  gemeldeten `core-drift-vcs`-Verletzung weiterhin `1 Datei(en) geprüft,
  0 Befund(e)`, Exit 0 — vor **und** nach dem Fix, mit Positiv-Kontrolle
  gemessen; der `--staged`-Pfad, also der lokale `pre-commit`-Hook, verhält
  sich genauso. Ursache ist eine Ebene über der behobenen: go-gits Tree-Noder
  bildet ein unlesbares Unterverzeichnis auf *„Verzeichnis nicht vorhanden"*
  ab, der Tree-Diff liefert die Datei deshalb als `A` statt `M`
  (instrumentiert gemessen: `status=65` statt `status=77`), und `CheckVCS`
  überspringt `A` planmäßig — `FileAt` und damit die neue Fallunterscheidung
  werden nie erreicht.
- `verifizierbar`: ja — Probe-Repo, `git pack-objects` über
  `$BASE:docs/plan/adr`, `git prune-packed`, Pack auf `loose-*` umbenennen;
  danach `--range`-Lauf (Exit 0, 0 Befunde) gegen den Lauf mit
  zurückbenanntem Pack (Exit 1, `core-drift-vcs`). Ein Regressionstest auf
  Adapter-Ebene fängt diesen Fall nicht, weil er nicht im Adapter entsteht.
- `klasse`: `vcs-stilles-gruen-bei-unsichtbarem-tree-objekt` — **dieselbe
  Klasse wie R1/F-1** (`fehlendes-objekt-als-abwesende-datei-gelesen`), eine
  Ebene höher im Objektgraphen
- **Zur Erreichbarkeit, ehrlich:** Die Probe ist konstruiert — genauso
  konstruiert wie die Blob-Probe, mit der der Slice seinen eigenen Fall belegt.
  Ein natürlicher Zwei-Lauf-Zyklus von `git maintenance run --task=loose-objects`
  lieferte in meiner Messung Exit 2 (`Range-Basis … nicht auflösbar`). Was den
  Befund trägt, ist nicht die Häufigkeit, sondern die **Reichweite der Zusage**:
  Der Slice begründet die Erreichbarkeit seines eigenen Falls damit, dass
  `git maintenance` *„in Stapeln packt"* — ein Stapel enthält Trees so gut wie
  Blobs, die Begründung trägt für beide gleich weit.

### F-2 — Neuer Fehlalarm: ein Tree-Eintrag **ohne** Blob (Submodul/gitlink) wird als „nicht lesbares Objekt" gemeldet, Exit 2

- `kategorie`: **HIGH** (blockierend)
- `quelle`: `DC-FA-VCS-001`; Reviewer-Anker 2 (Kern-Modul meldet falschen
  Befund/Exit-Code)
- `pfad`: `internal/adapter/driven/git/git.go:265-278` (`entryUnreadable`:
  `case err == nil: return true, nil`) und der Kommentar darüber
  (*„existiert der Tree-Eintrag zu path, war das zugehörige Blob unlesbar"*)
- `befund`: `tree.File` liefert `object.ErrFileNotFound` auch dann, wenn der
  Tree-Eintrag existiert, aber **kein Blob** dahinter steht — bei einem
  gitlink-Eintrag (Submodul, Modus `160000`) ist das der Normalfall, weil
  `GetBlob` auf einen Commit-Hash zwangsläufig scheitert. `entryUnreadable`
  liest die Existenz des Eintrags als *„Blob war unlesbar"* und bricht den
  ganzen Lauf mit Exit 2 und der Diagnose `nicht lesbares Objekt zu
  "docs/plan/adr/sub" …` ab, wo der Vor-Fix-Stand korrekt befundfrei blieb;
  die mitgelieferte Abhilfe (`git repack -A -d`) hilft in diesem Zustand nicht.
  Ausgelöst wird es, sobald ein Submodul-Pfad die konfigurierte Klasse
  `vcs.paths` trifft und sein Zeiger sich über die Range bewegt.
- `verifizierbar`: ja — Probe-Repo mit
  `git update-index --add --cacheinfo 160000,<sha>,docs/plan/adr/sub`,
  `vcs.paths: ["docs/plan/adr/*"]`, Zeiger in einem zweiten Commit bewegt:
  `d-check:prefix` ⇒ `0 Befund(e)`, Exit 0; `d-check:latest` ⇒ Exit 2.
- `klasse`: `tree-eintrag-ohne-blob-als-unlesbar-gelesen`
- **Zur Erreichbarkeit, ehrlich:** Mit der Konfiguration *dieses* Repos
  (`vcs.paths: ["docs/plan/adr/[0-9]*.md"]`) ist der Fall praktisch
  ausgeschlossen. `vcs.paths` ist Adopter-Konfiguration, und ein Glob wie
  `docs/**` über einem Doku-Submodul ist keine erfundene Form. Der Fall ist
  fail-closed, also laut — aber er ist genau der Tausch, den der Plan als
  Fokus dieser Runde benennt: *„Ein Fix, der beide Fälle gleich behandelt,
  tauscht ein stilles Übersehen gegen einen Fehlalarm."*

### F-3 — Drei Doku-Orte und der Carveout sagen den Abschluss über die **ganze** Klasse zu; gemessen ist die Blob-Hälfte

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 (Reichweite einer Botschaft; Grenzen-Liste);
  Reviewer-Anker 8 und 18; `v6.5.0` · `regelwerk/modul-13-quality-gates.md`
  §Hard Rule, Absatz *„Ein Gate ohne seine Grenze behauptet ebenfalls zu viel"*
- `pfad`: `docs/user/benutzerhandbuch.md:1596-1597` (*„— alle mit Exit 2. Der
  Lauf verweigert die Aussage, statt fälschlich grün zu melden."*) und
  `:1610-1611` (*„Ab der Version nach `v0.75.0` ist das behoben."*);
  `harness/sensors/adr-check.md:33` (*„der Lauf bricht dann ab"*) und `:47`;
  `docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md:61-64` (*„Der
  Rest-Risiko-Umfang ist klein und benennbar"*)
- `befund`: Alle vier Stellen formulieren über die Klasse *„unsichtbar
  benannter Pack"* und zählen drei Symptome auf, die **alle** Exit 2 tragen;
  der gemessene vierte Ausgang derselben Klasse — `0 Befunde`, Exit 0 bei
  einem verschluckten Tree-Objekt (F-1) — fehlt, und die Handbuch-Zeile
  *„Ab der Version nach `v0.75.0` ist das behoben"* sagt einem Adopter das
  Gegenteil des gemessenen Verhaltens. Die Grenzen-Liste liest sich durch ihre
  Form als Menge und ist eine Auswahl.
- `verifizierbar`: ja — dieselbe Probe wie F-1 gegen die zitierten Sätze halten.
- `klasse`: `grenzenliste-nennt-nur-die-gemessene-haelfte-der-klasse`

### F-4 — Die Gegenrichtungs-Probe der Botschaft ist blind: eine neu angelegte ADR erreicht `FileAt` nie

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 (Vor einer Messung steht die Form ihres
  Gegenstands); Reviewer-Anker 17 (Proxy-Messung)
- `pfad`: Commit-Botschaft `d308a882`, dritter Spiegelstrich des Bruch-Tests
  (*„Neu angelegte ADR, im BASE-Tree wirklich nicht vorhanden: weiterhin 0
  Befunde, Exit 0 — KEIN Fehlalarm"*), gestützt auf slice-218 §2 DoD (4);
  Gegenstelle: `internal/hexagon/core/rules/vcs.go:33`
- `befund`: `CheckVCS` behandelt den Diff-Status `A` mit einem leeren
  `case`-Zweig — eine hinzugefügte Datei erreicht `FileAt` und damit die neue
  Fallunterscheidung überhaupt nicht. Die als Beleg für die Gegenrichtung
  genannte Probe kann den behaupteten Unterschied deshalb nicht sehen:
  gegen ein Image, in dem `entryUnreadable` absichtlich **immer** „unlesbar"
  meldet, liefert dieselbe Probe unverändert `3 Datei(en) geprüft, 0
  Befund(e)`, Exit 0. Die Eigenschaft selbst *ist* geprüft — aber vom
  Unit-Test `TestFileAtUnlesbaresObjekt`, nicht von der zitierten End-zu-End-Probe.
- `verifizierbar`: ja — Mutanten-Image (`entryUnreadable` ⇒ `return true, nil`)
  gegen ein Probe-Repo, dessen HEAD nur neue ADRs hinzufügt.
- `klasse`: `probe-kann-den-behaupteten-unterschied-nicht-sehen` — Verwandte
  von R1/F-3 `probe-misst-anderen-zustand-als-den-genannten-ausloeser`; zweite
  Instanz derselben Proben-Qualitäts-Klasse in diesem Slice

### F-5 — `CO-001` nennt als Folge-Slice den Slice, der ihn anlegt — und der schließt vor der Auflösung

- `kategorie`: MEDIUM
- `quelle`: `v6.5.0` · `regelwerk/modul-07-carveouts.md` §Ziel-Form: Carveout
  (*„Fehlt der Folge-Slice, ist der Carveout de facto permanent"*) und
  §Regeln gegen typische Fehlannahmen (*„Slice schlägt Memo"*);
  `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice
  (*„Ein Folge-Slice, der … vor dem verweisenden schließt, ist keine
  [Adresse]"*)
- `pfad`: `docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md:22-24`;
  slice-218 §1 Abgrenzung Punkt 4 (*„ein Folge-Slice übernimmt es, nämlich das
  Release"*)
- `befund`: Das Pflichtfeld *Folge-Slice* zeigt auf slice-218 selbst; der
  Auflösungs-Trigger ist *„das nächste Release ist veröffentlicht"*, und dieses
  Release ist in slice-218 weder DoD-Punkt noch Closure-Trigger. Nach der
  Closure trägt damit kein benannter Vorgang die Auflösung — und der
  Carveout-Audit des Kanons hängt an der Wellen-Closure, die dieses Repo
  wellenlos nicht fährt. Der Punkt in §1 nennt den Träger als *„das Release"*
  ohne Kennung, also ohne Adresse.
- `verifizierbar`: ja — nach der Closure von slice-218 zeigt
  `docs/plan/carveouts/` einen aktiven Carveout, dessen Folge-Slice in `done/`
  liegt und dessen Trigger in keiner offenen Datei genannt ist.
- `klasse`: `carveout-folge-slice-zeigt-auf-den-eigenen-slice`

### F-6 — Der Verzicht auf eine Spezifikations-Zeile steht auf einem Kriterium, das der eigene Präzedenzfall widerlegt

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 (Zitat-/Kriterien-Geltungsbereich; Reviewer-Anker 9);
  `spec/spezifikation.md` §`DC-FA-VCS-001.a` samt Historie-Zeile 2026-08-31
- `pfad`: slice-218 §1 Abgrenzung Punkt 3 (*„Ein Spec-Eintrag käme erst in
  Frage, wenn ein **neuer** Grund-Code entstünde"*)
- `befund`: Die Historie derselben Spec-Sektion führt einen strukturgleichen
  Vorgang vom 2026-08-31 — ein `vcs`-Defekt, der im CI-Modus still grün
  meldete, ohne Lastenheft-Bump und **ausdrücklich ohne neuen Grund-Code**
  behoben — und dieser Vorgang bekam trotzdem einen Spezifikations-Eintrag
  (der Mechanismus-Absatz in Schritt 2 plus Historie-Zeile). Das hier
  aufgestellte Kriterium ist damit enger als die belegte Praxis des Repos, und
  die Folge ist beobachtbar: Schritt 1 sagt fail-closed nur für ein fehlendes
  oder unlesbares `.git` zu; dass ein **einzelnes** unlesbares Objekt eines
  Kandidaten den Lauf jetzt abbricht, steht in keinem Spec-Stratum.
- `verifizierbar`: ja — `spec/spezifikation.md` §`DC-FA-VCS-001.a` Schritte 1–5
  gegen das gemessene neue Verhalten (Exit 2 mit `nicht lesbares Objekt zu …`).
- `klasse`: `kriterium-gesetzt-statt-praezedenz-nachgeschlagen`

### F-7 — §5 Closure-Trigger nennt weiterhin nur DoD (1) bis (3)

- `kategorie`: LOW
- `quelle`: `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Roadmap-Regeln
  (Closure-Kriterien als beobachtbare Bedingung)
- `pfad`: slice-218 §5, Zeile 197
- `befund`: Die Plan-Änderung hat DoD (4) (der Fix samt Regressionstest) und
  DoD (5) (`CO-001`) ergänzt; der Closure-Trigger blieb auf *„DoD (1) bis (3)
  abgehakt"* stehen. Wer §5 als Schließbedingung liest, kann den Slice
  formal schließen, ohne dass Fix und Carveout in der Bedingung stehen.
- `verifizierbar`: ja — §2 gegen §5 desselben Dokuments halten.
- `klasse`: `closure-trigger-nach-plan-aenderung-nicht-nachgezogen`

### F-8 — §4 zitiert „§1 Punkt 2 schließt Code aus"; Punkt 2 sagt das nicht mehr, und die Bedingung ist eingetreten

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §5 (die direkteste Quelle, und das Feld statt des
  Titels lesen); `v6.5.0` · `regelwerk/modul-05-planning-harness.md`
  §Trigger je Lifecycle-Übergang
- `pfad`: slice-218 §4, Zeile 189
- `befund`: Der Rückführungs-Trigger nach `next/` begründet sich mit
  *„(§1 Punkt 2 schließt Code aus)"*. Nach der Renummerierung lautet Punkt 2
  *„Kein Sensor auf Pack-Namen"*, und Produkt-Code ist überhaupt nicht mehr
  ausgeschlossen. Zugleich ist die beschriebene Bedingung — die Grenze
  verlangt eine Verhaltens-Antwort — tatsächlich eingetreten; §4 steht
  unverändert daneben und widerspricht §1, ohne dass eine Zeile die beiden
  in Beziehung setzt.
- `verifizierbar`: ja — §1 Punkt 2 gegen den Klammerverweis in §4 halten.
- `klasse`: `stehengebliebener-trigger-widerspricht-der-plan-aenderung`

### F-9 — Das Handbuch nennt dieselben zwei Meldungen in zwei aufeinanderfolgenden Absätzen

- `kategorie`: LOW
- `quelle`: Maintainability (Doku-Drift)
- `pfad`: `docs/user/benutzerhandbuch.md:1588-1590` gegen `:1593-1597`
- `befund`: Der Absatz vor der Umschreibung endet mit *„Sie sehen dann je nach
  Repo-Zustand `staged-Basis …` oder `HEAD-Tree …`"*; der neue Absatz beginnt
  mit *„Je nachdem, was der unsichtbare Pack verschluckt, sehen Sie
  `staged-Basis …`, `HEAD-Tree …` oder …"*. Die zwei Meldungen stehen zweimal
  unmittelbar hintereinander — ein Rest der Umschreibung.
- `verifizierbar`: ja — beide Absätze lesen.
- `klasse`: `doppelte-symptomliste-nach-umschreibung`

### F-10 — Die Plan-Änderung sagt „das steht hier vor dem Code" und reist mit dem Code in einem Commit

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §6 Schritt 4 (*„das gehört vor den Code, nicht in den
  Bericht danach"*)
- `pfad`: slice-218 §1, Absatz *Plan-Änderung*; Commit `d308a882`
- `befund`: Plan-Änderung, Produkt-Fix, Regressionstest, Carveout und die
  umgeschriebene Doku liegen in **einem** Commit. Die Aussage *„nicht still
  ausgeweitet, und das steht hier vor dem Code"* ist inhaltlich plausibel, aus
  der Historie aber nicht belegbar: Es gibt keinen Stand, in dem die geöffnete
  Abgrenzung ohne den Code dasteht.
- `verifizierbar`: ja — `git show --stat d308a882` zeigt Plan und
  `internal/adapter/driven/git/` im selben Commit; `6e0b7960` fasst den Plan
  nicht an.
- `klasse`: `plan-aenderung-und-code-im-selben-commit`

### F-11 — `CO-001` trägt Deliberation über eine verworfene Fassung

- `kategorie`: INFO
- `quelle`: Maintainability. *Ausdrücklich **kein** §3.7-Verstoß* — dessen
  Geltungsbereich sind Kommentare in Code/Konfiguration/Skripten und
  Zustandsfelder, nicht Fließtext in einem Planungs-Dokument; die Regel wird
  hier bewusst **nicht** über ihren Bereich hinaus gezogen.
- `pfad`: `docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md:18-20`
- `befund`: Der Kopfblock erklärt, eine *„frühere Fassung dieses Carveouts"*
  habe das Gegenteil behauptet. Diese Fassung war nie committet
  (`git log --diff-filter=A` nennt nur `d308a882`), der Satz erzählt also die
  Entstehung eines Laufs in einem lebenden Artefakt; die Vorlage hat für
  Vorgangs-Historie den Abschnitt `## Geschichte`, und die Commit-Botschaft
  trägt dieselbe Information bereits.
- `verifizierbar`: nein (Urteil)
- `klasse`: `entstehungsgeschichte-im-lebenden-artefakt`

### F-12 — Drei Randbefunde zur Fehler-Menge, die den Fix nicht belasten, aber zum Bild gehören

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.8 (Modul-Zusage nur über die Scan-Menge)
- `pfad`: `internal/adapter/driven/git/git.go:270-278`; go-git `v5.19.2`
  `plumbing/object/tree.go`, `internal/pathutil/tree.go`
- `befund`: (a) `FindEntry` in `v5.19.2` prüft **vorab** `pathutil.ValidTreePath`
  und liefert dafür eine dritte Fehlerklasse (`ErrInvalidPath` bei
  Steuerzeichen, `.`/`..`-Komponenten, `.git`-Tarnungen). Sie fällt in den
  `default`-Zweig und macht aus einem vormals befundfreien Pfad Exit 2 — für
  reguläre Repos nicht erreichbar, aber es ist kein „gibt es nicht". (b) Der
  befundfreie Zweig von `FileAt` ist über `CheckVCS` gar nicht erreichbar:
  `A` ruft `FileAt` nie, `M`/`D` liefern nur Pfade, die im jeweiligen Stand
  existieren — die Gegenrichtung ist damit eine reine Adapter-Eigenschaft.
  (c) Im `--staged`-Modus greift die Fallunterscheidung für unlesbare Blobs
  nie: `diffTreeIndex` liest über `headTree.Files()` bereits jedes Blob und
  bricht vorher mit `HEAD-Tree nicht lesbar` ab (gemessen, vor und nach dem
  Fix identisch).
- `verifizierbar`: ja — go-git-Quelle im `deps`-Image; Messungen (6) und (8).
- `klasse`: `fehler-menge-und-erreichbarkeit-des-neuen-zweigs`

---

## Negativbefunde (geprüft, ohne Befund)

- **Kern-Fall des Slice, nachgefahren statt geglaubt.** Probe-Repo mit
  `**Status:** Accepted`-ADR, Core-Änderung, nur das BASE-Blob im unsichtbar
  benannten Pack: Vor-Fix `1 Datei(en) geprüft, 0 Befund(e)` / Exit 0, nach dem
  Fix `error: nicht lesbares Objekt zu … : file not found` / Exit 2, mit
  kanonischem Pack-Namen `1 Befund(e)` / Exit 1. Der Fix wirkt genau dort, wo
  er es zusagt.
- **Gegenrichtung, drei Formen, kein Fehlalarm.** Neu angelegte ADR, `git mv`
  einer Accepted-ADR und eine ADR in einem **neuen** Unterverzeichnis: `latest`
  und `prefix` melden byte-gleich denselben einen Rename-Befund; auf die
  Neuzugänge feuert nichts. (Dass diese Läufe die Fallunterscheidung nicht
  berühren, steht als F-4 — der *Befund* „kein Fehlalarm" bleibt richtig.)
- **Der Regressionstest fällt ohne den Fix.** Fallunterscheidung in einer
  Kopie zurückgenommen, `docker build --target test`:
  `--- FAIL: TestFileAtUnlesbaresObjekt … unlesbares Objekt still übersprungen:
  ok=false len=0 — erwartet war ein Fehler`, wörtlich die in der Botschaft
  zitierte Ausgabe. Alle übrigen Pakete grün.
- **Vollständigkeit der „gibt es nicht"-Fehlermenge.** go-git `v5.19.2`-Quelle
  gelesen: `Tree.entry` liefert `ErrEntryNotFound`, `Tree.dir` liefert
  `ErrDirectoryNotFound` für einen fehlenden oder nicht-Baum-Eintrag. Die
  beiden im `switch` genannten Werte decken die Abwesenheit vollständig ab.
- **Unlesbarer Zwischen-Tree wird im `entryUnreadable`-Pfad nicht
  verschluckt.** `Tree.dir` reicht den Storer-Fehler (`ErrObjectNotFound`,
  Decode-Fehler) durch; er landet im `default`-Zweig und über den
  `ferr`-Zweig als `nicht lesbarer Tree-Eintrag zu … ` im Ergebnis. Der
  gebaute Fall zeigt die Lücke eine Ebene früher (F-1), nicht hier.
- **`--staged` bringt keine Regression.** Unlesbares BASE-Blob: vor und nach
  dem Fix identisch `HEAD-Tree nicht lesbar: object not found`, Exit 2.
- **Wörtliche Fehlermeldungen.** Die dritte, neu dokumentierte Form stimmt
  byte-genau mit der gemessenen Ausgabe (`nicht lesbares Objekt zu "<pfad>" in
  "<ref>": file not found`); die zwei älteren Formen stehen so im Adapter
  (`staged-Basis %q nicht auflösbar`, `HEAD-Tree nicht lesbar`).
- **`CO-001`s Kern-Aussage „dieses Repo ist nicht betroffen" trägt.**
  `Makefile:380` führt `adr-check: build`, das Recipe fährt `$(IMAGE):latest`
  aus dem Quellstand (`DCHECK_RUN`, `Makefile:122`), und `.github/workflows/ci.yml`
  ruft im Schritt *Traceability + ADR-Immutable + Slice-Closure-Übergang*
  `make adr-check RANGE="$RANGE"` — kein gepinntes Release-Image im Spiel. Für
  die Blob-Hälfte stimmt die Aussage; für die Tree-Hälfte gilt sie nicht, das
  ist F-1/F-3 und kein Fehler dieses Satzes.
- **Carveout-Form.** Alle sechs Pflicht-Header-Felder der Vorlage vorhanden;
  Auflösungs-Trigger ist beobachtbar und mit einem Prüf-Kommando unterlegt;
  Index in `docs/plan/carveouts/README.md` und die Bindung-Spalte in
  `harness/README.md` §Sensors tragen `CO-001` (beides fordert Modul 7). Dass
  keine Gate-Konfiguration den Carveout ausdrückt, ist im Abschnitt
  *Geltungs-Konfiguration* offen benannt und für einen Carveout ohne
  Werkzeug-Schalter die ehrliche Antwort.
- **`AGENTS.md` §3.7, Kommentar-Klassen.** Die vier neuen Kommentarblöcke
  (`FileAt`, `entryUnreadable`, Test-Doc, vier Inline-Zeilen) tragen *Grenze*,
  *Zusage* und *Abgrenzung*; keine Slice-Nummer, kein Review-Bezug, kein
  Mess-Label, keine Deliberation über Verworfenes. Herkunft steht genau einmal
  und auflösbar (`DC-FA-VCS-001` im Test-Kommentar).
- **`AGENTS.md` §1/§3.4/§3.6.** Keine Suppression, keine gesenkte Schwelle,
  kein `//nolint`; keine Spec-Stratum-Änderung, kein Abwärts-Verweis.
- **Hexagon-Richtung (ADR-0005/ADR-0012).** Der Fix bleibt vollständig im
  driven-Adapter, führt keine neuen Importe ein und lässt Port und Kern
  unberührt; `make arch-check` läuft in `gates`.
- **`DC-QA-03`.** Kein Netzzugriff, kein Schreibzugriff auf das geprüfte Repo;
  alle Proben liefen mit `--network none` und `:ro`-Mount.
- **`MR-066` korrekt zitiert und eingelöst.** Geltungsbereich (*Slice über der
  Ein-Sitzungs-Review-Grenze, der nicht zurückgeführt wird*) trifft zu; beide
  Pflichten stehen im Plan — Grund und **vorab benannte** Ersatz-Form
  (*zwei Runden gegen je einen abgeschlossenen Stand, Runde 2 mit deklariertem
  Fokus*). Der von `MR-066` verlangte **Vollzug im Review-Report** ist dieser
  Report: Runde 2 lief gegen den abgeschlossenen Stand `d308a882` und gegen
  den deklarierten Fokus — und der Fokus hat getragen, F-1 und F-2 stammen
  beide aus ihm.
- **Repo-Gates nachgefahren.** `make doc-check` und `make adr-check` je
  `757 Datei(en) geprüft, 0 Befund(e)`, Exit 0; `make trace-check` ebenso. Die
  Packs dieses Repos tragen den kanonischen Namen, F-1 wirkt hier heute nicht.
- **`AGENTS.md` §5, Reichweite der Botschaft.** Die übrigen Aussagen der
  Commit-Botschaft halten der Messung stand: der partielle Pack, der
  kanonische Kontrolllauf, das Zurücknehmen des Fixes, die Prerequisite
  `build`. Überdehnt ist genau eine (F-4), unvollständig genau eine (F-3).

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 2 | F-1, F-2 |
| MEDIUM | 4 | F-3, F-4, F-5, F-6 |
| LOW | 4 | F-7, F-8, F-9, F-10 |
| INFO | 2 | F-11, F-12 |

**Wiederholungs-Signal.** F-1 ist dieselbe Fehlerklasse wie R1/F-1, eine Ebene
höher; F-4 ist die zweite Instanz der Proben-Qualitäts-Klasse aus R1/F-3
(*eine Probe, die den behaupteten Unterschied nicht sehen kann*). Beide sind
noch keine dritte Wiederholung, aber beide gehören in die Closure-Notiz und
ins Beobachtungs-Register — die zweite ist der Kandidat für einen eigenen
Register-Eintrag: *ein Bruch-Test, der die Code-Verzweigung nicht erreicht,
die er belegen soll.*

## Verdikt

**Blockiert.** Zwei HIGH stehen der Closure entgegen, und beide liegen im
deklarierten Fokus dieser Runde: Der stille Pfad ist für Blobs zu und für
Tree-Objekte offen — in beiden Modi, auch im lokalen `pre-commit`-Hook (F-1) —,
und die Fallunterscheidung erzeugt für einen Tree-Eintrag ohne Blob einen
Fehlalarm mit irreführender Diagnose (F-2). Solange F-1 steht, sagen Handbuch,
Sensor-Datei und Carveout einem Adopter mehr zu, als das Produkt hält (F-3);
das ist die Klasse, gegen die dieser Slice überhaupt angetreten ist.

Der Fix selbst ist an seiner Stelle richtig und belegt: die Blob-Hälfte ist
gemessen geschlossen, der Regressionstest fällt ohne ihn, die Gegenrichtung
erzeugt in den drei geprüften Alltagsformen keinen Fehlalarm. Was fehlt, ist
die Reichweite — und, für die Zukunft, ein Bruch-Test, der die Verzweigung
erreicht, die er belegen soll.
