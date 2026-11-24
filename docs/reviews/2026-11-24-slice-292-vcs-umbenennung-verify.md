# Verifikation — slice-292: Kern und git-Adapter, `vcs` folgt der Umbenennung

- **Rolle:** Verifier (Modul 11) — DoD, Plan und Spezifikation gegen den gelieferten Code. Die Reviews
  R1 und R2 dienten nur als Kontext; Zuordnungen, Zählungen, Mutationen und Messungen sind selbst
  gemacht, in einem frischen Kontext.
- **Gegenstand:** Plan
  [`slice-292`](../plan/planning/done/slice-292-vcs-umbenennung-modulkern.md) §1, §2, §3 (Tabelle und
  Testfälle) und §6;
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  (Lastenheft 0.114.3, 30 Kriterien),
  §[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  Schritt 2 und 4, [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 9 bis 13.
- **Range:** `6aca94b0..879e576d` auf `main`, lokal, fünf Commits. Code in `bea51566` (Stufe 1),
  `5999fb2b` (Stufe 2, Randfälle, Index-Angleichung, Handbuch, Sensor-Datei) und `95aef23b` (zwei Tests
  nach R1); dazu die Reports [R1](2026-11-24-slice-292-vcs-umbenennung-r1.md) und
  [R2](2026-11-24-slice-292-vcs-umbenennung-r2.md). 18 Dateien, 2730 Zeilen neu.
- **Datum:** 2026-11-24.

## Gates am Repo (selbst gefahren, Spitze `879e576d`, Arbeitsbaum sauber)

| Kommando | Ergebnis |
|---|---|
| `make gates` | Exit 0. Letzte Zeile: `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`. Davor `0 issues.` (lint), `gesamt: 0 Befund(e)` (arch-check), Semgrep `Ran 55 rules on 75 files: 0 findings.`, `targets` und `planning` je `1289 Datei(en) geprüft, 0 Befund(e)`; der Lauf lief als Hintergrundjob, ein zweiter `make gates` einer anderen Sitzung lief gleichzeitig und teilt die Image-Tags — beide Läufe blieben grün |
| `make coverage-gate` | Exit 0. `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%` |
| `make image-test` | Exit 0. Letzte Zeile: `image-test: OK — Akzeptanzkriterien erfüllt (DC-FA-DIST-001, DC-FA-CLI-007, DC-FA-CLI-008, DC-FA-CLI-013, DC-FA-MET-001, DC-FA-PSYNC-001, DC-FA-LLEN-001)` |
| `make build` | Exit 0; `d-check --print-config` im gebauten Image nennt `# follow-renames: true` im `vcs`-Block (Zeile 128) |
| `make blackbox-probe REF=6aca94b0` | Exit 0, `byte-identisch über 20 Vergleiche`. **Trägt für `vcs` nichts:** keine der vier Fixtures (`ids`, `links`, `sauber`, `targets`) fährt `--range`; den vcs-Vergleich ersetzt unten der Vorher/Nachher-Lauf |
| `make doc-check` nach dem Schreiben dieses Reports | siehe Schluss-Abschnitt |

## DoD-Punkte

### 1 — Port und Adapter: bestätigt

Gelesen: `Renames` am Port (`internal/hexagon/port/driven/vcs.go`) und der Adapter in
`internal/adapter/driven/git/git.go` (`rangeRenames`, `stagedRenames`, `diffFiles`, `finalRenames`).

- **Commit-Menge:** `rangeWalk.visit` läuft von HEAD, bricht an jedem Commit in `ancestors(BASE)` ab.
  Ein Schritt ist ein Commit mit genau einem Eltern-Commit; Merge und Wurzel tragen nichts bei
  (`statesOf`).
- **Gleichheit:** `fileEntry{hash, symlink}`; das Ausführbar-Bit steht nicht darin. Die Bibliotheks-Erkennung
  (`DiffTree`, `DetectRenames`) kommt im neuen Code nicht vor.
- **Kette, Mehrdeutigkeit, Ende:** `lineSet` je Quelle und Commit, über Eltern-Zustände fortgeführt, durch
  Mengen-Vereinigung am Merge; `finalRenames` verlangt Quelle in BASE, genau ein Ende, keine Linie ohne
  Ende, Ende in HEAD, nicht in BASE und kein Ende einer zweiten Quelle.
- **Basis kein Vorfahr:** `rangeRenames` liefert eine leere Zuordnung ohne Fehler; `--staged` ohne HEAD ebenso.
- **fail-closed:** Commit, Eltern-Commit, beide Bäume, jeder geänderte Teilbaum und der Index melden den
  Fehler weiter; Teilbäume gleicher Kennung bleiben ungelesen, kein Diff-Walker.
- **Index-Angleichung:** `AllPaths` nimmt im Index nur `isFileMode` (regulär, ausführbar, Symlink), derselbe
  Filter wie im Baum.

**Rot ohne die Regel — selbst gefahren.** Ein Klon im Scratchpad, je eine Mutation, `make test`, danach
`git checkout --`. Elf Mutationen, alle elf rot, mit der Ursache, die der Plan nennt:

| M | Mutation | rot in | Ursache gelesen |
|---|---|---|---|
| V1 | `AllPaths` im Index wieder jeder Eintrag | `TestAllPathsStagedGitlinkIstKeinePfad`, `TestVCS_ImmutableDateiDurchGitlinkErsetzt/--staged` | `Index-Menge = [keep.md mod], erwartet nur keep.md`; **am Vorzustand gemessen Exit 2** (`object not found`) statt 1 — der Rot-Beleg der Angleichung |
| V2 | nur der Schritt von HEAD gelesen (Eltern nicht besucht) | `TestRenamesLinear`, `TestRenamesMehrdeutigkeit`, `TestRenamesMerge`, `TestRenamesUnlesbaresObjekt`, `TestVCS_FollowRenames`, `TestVCS_FollowRenamesZehnLaeufe` | der Fall „Move und Nachzug in verschiedenen Commits“ |
| V3 | „Basis kein Vorfahr“ nicht geprüft | `TestRenamesMerge`, `TestVCS_FollowRenamesBasisKeinVorfahr` | Ausgabe `0 Befund(e)` statt `1 core-drift-vcs`, Meldung des Falles |
| V4 | Bedingung „Ende nicht in BASE“ entfernt | `TestRenamesMehrdeutigkeit` (Tausch, freigemachter Name, Rotation) | |
| V5 | zwei Quellen auf ein Ende zulässig | `TestRenamesMehrdeutigkeit` | |
| V6 | unlesbarer geänderter Teilbaum still übergangen (Diff-Walker-Verhalten) | `TestRenamesUnlesbaresObjekt/Teilbaum,_der_sich_im_Schritt_ändert` | `unlesbares Objekt still übergangen — erwartet war ein Fehler` |
| V7 | `FollowRenames` im Rand nicht verdrahtet | `TestDecode_VCSFollowRenames`, `TestVCS_FollowRenames`, `…FailClosed`, `…ZehnLaeufe`, `…Staged` | |
| V8 | `Renames` auch mit `ignore-link-targets` allein gerufen | `TestVCSFollowRenames`, `TestVCSFollowRenamesRandfaelle` | Zähler der Port-Aufrufe |
| V9 | Gitlink in der Baum-Enumeration | `TestAllPathsGitlinkWirdUebersprungen`, `TestVCS_ImmutableDateiDurchGitlinkErsetzt` | der `--range`-Zweig |
| V10 | das Ausführbar-Bit zählt | `TestRenamesArtDesEintrags` | |
| V11 | Gleichsetzen im Kern abgeschaltet | `TestVCSFollowRenames`, `TestVCSFollowRenamesRandfaelle`, `TestVCS_FollowRenames`, `…ZehnLaeufe`, `…Staged` | |

Eine zwölfte, selbst gesetzte Mutation blieb grün: die Zeile „Neuanlage lässt die Linie wegfallen“
(`if ev.created[p]` in `lineSet.land`). Das ist der Befund L-2 von R1, unverändert; siehe V-2 unten.

**Echte git-Fixtures:** die Adapter-Tests legen Objekte im Speicher eines echten `.git` an (Tabelle
`runRenameCases`, drei Kennungs-Salze), die CLI-Tests ein On-Disk-Repo mit go-git; kein Arbeitsbaum-Stat-Vergleich
in den Adapter-Fixtures, die CLI-Fixtures laufen über `writeAt` mit dem Wächter
`coretest.GitFixtureRewriteHazard` (siehe Risiko 2).

### 2 — Regel und Konfigurations-Rand: bestätigt

**Rand.** `rawVCS.FollowRenames bool`, gelesen wie `IgnoreLinkTargets`. Am gebauten Image selbst gefahren:
`'ja'`, `[true]` und `1` enden mit Exit 2 (`cannot unmarshal … into bool`), `true`, `false` und ein Schlüssel
ohne Wert mit Exit 0. `TestDecode_VCSFollowRenames` fährt beide Schlüssel nebeneinander und verlangt je
dasselbe Ergebnis.

**Kriterium → Test, selbst zugeordnet.** Alle 30 Kriterien tragen einen Test. Unter „Kern“ steht der Test, der
aus der Zuordnung ein Urteil macht, unter „Adapter“ der, der die Zuordnung liefert; wo beide Schichten
tragen, tragen sie zusammen (Port als Naht).

| # | Kriterium | Test |
|---|---|---|
| 1 | Happy Path | `TestVCS_RangeClean`, `TestVCSModified` (bestehend, unverändert) |
| 2 | Boundary (Modul-aus / git-frei) | `TestVCSDispatch` („Modul-aus“), `TestVCSInert` (bestehend) |
| 3 | Negative | `TestVCS_RangeDrift`, `TestVCSDeleteAddClass`, `TestVCSRenameOutOfClass` (bestehend) |
| 4 | fail-closed (git-Eingabe fehlt) | `TestVCS_FailClosed`, `TestVCSFailClosed`, `TestVCSAllPathsFehlerAnBeidenEnden` (bestehend) |
| 5 | Boundary (Pfad-Nachzug, ohne follow-renames) | `TestVCSIgnoreLinkTargets`, `TestVCSIgnoreLinkTargetsUmzug` — **Datei `vcs_pfad_nachzug_test.go` unverändert**; dazu die Spalte „ignore“ in `TestVCSFollowRenames` |
| 6 | Happy Path (Umbenennung folgen) | CLI `TestVCS_FollowRenames` („…in einem commit“, „reiner move und nachzug in getrennten commits“); Kern `TestVCSFollowRenames` (Nachzug, mit Anker, Bild und Referenz-Definition); Adapter `TestRenamesLinear` („ein Commit“, „Move und Inhaltsänderung in getrennten Commits“) |
| 7 | Negative (gelöscht und neu angelegt) | CLI `TestVCS_FollowRenames`; Kern `TestVCSFollowRenames`; Adapter `TestRenamesLinear` („gelöscht und anderes angelegt“) |
| 8 | Negative (Umbiegen) | CLI `TestVCS_FollowRenames`; Kern `TestVCSFollowRenames` („umbiegen auf eine bestehende, nicht umbenannte datei“) |
| 9 | Negative (Anker, Linktext) | CLI `TestVCS_FollowRenames` (zwei Fälle); Kern `TestVCSFollowRenames` (Anker, Linktext, übrige Zeile, Link entfernt) |
| 10 | Negative (umbenennen und ändern) | Adapter `TestRenamesLinear` („Umbenennen und im selben Commit ändern“ ⇒ keine Zuordnung; „Move und Inhaltsänderung in getrennten Commits“ und „Inhaltsänderung vor dem Move“ ⇒ Zuordnung) |
| 11 | Boundary (Kette) | `TestRenamesLinear` („Kette über drei Commits“); `TestRenamesMerge` („X → Y und Y → Z auf zwei Seiten setzen sich nicht fort“, „eine Seite benennt um, die andere ändert nur den Inhalt“); CLI `TestVCS_FollowRenamesZehnLaeufe`; die drei Salze in `runRenameCases` |
| 12 | Negative (Mehrdeutigkeit) | `TestRenamesMehrdeutigkeit` („zwei verschwundene …“, „zwei neue …“, „zwei Linien führen zu zwei Enden“, „Tausch über einen Zwischennamen“); `TestRenamesMerge` („die andere Seite löscht den Pfad …“); `TestRenamesLinear` („Kette: Ende in HEAD gelöscht“) |
| 13 | Negative (Merge-Seite löscht, legt später neu an) | `TestRenamesMerge` („Merge: die andere Seite löscht den Pfad und legt ihn später neu an“) |
| 14 | Negative (Ende gelöscht und neu angelegt) | `TestRenamesLinear` („Ende gelöscht und in einem späteren Commit als andere Datei neu angelegt“) |
| 15 | Negative (Zwischenpfad neu belegt) | `TestRenamesLinear` („Zwischenpfad einer Kette neu belegt“); Kern `TestVCSFollowRenamesRandfaelle` (gleichnamig) |
| 16 | Boundary (seiten-gleich) | `TestVCSFollowRenamesUnveraenderterLink` |
| 17 | Boundary (beide Schlüssel) | `TestVCSFollowRenamesRandfaelle` („beide Schlüssel: Umbiegen auf eine andere Datei mit dem Namen der Quelle“) |
| 18 | Boundary (Negativ-Kriterien mit beiden Schlüsseln) | die Spalte „beide“ in `TestVCSFollowRenames`, `TestVCSFollowRenamesRandfaelle`, `TestVCS_FollowRenames`, `TestVCS_FollowRenamesStaged`; Ausnahme gleichen Namens in „zwei nicht teilnehmende Ziele mit gleichem Namen“ |
| 19 | Negative (beide Schlüssel, Umbiegen auf den alten Namen) | `TestVCSFollowRenamesRandfaelle` („Umbiegen auf den alten Namen: d2/c.md wird zu d3/a.md“) |
| 20 | Negative (beide Schlüssel, Umbenennung ohne Namenswechsel) | `TestVCSFollowRenamesRandfaelle` („Umbenennung ohne Namenswechsel …“) |
| 21 | Negative (Name freigemacht und neu belegt) | `TestRenamesMehrdeutigkeit` („Name freigemacht und neu belegt“); Kern „freigemachter Name …“ (beide Hälften) |
| 22 | Negative (zwei Quellen auf ein Ende) | `TestRenamesMehrdeutigkeit` („zwei Quellen auf ein Ende (Merge)“); Kern „zwei Quellen auf ein Ende …“ (beide Hälften) |
| 23 | Boundary (Rotation) | `TestRenamesMehrdeutigkeit` („Rotation: A nimmt als Quelle teil, C nicht“); Kern „Rotation: Nachzug von A auf X folgt“, „… von C auf A bleibt Drift“ |
| 24 | Boundary (Löschen und Neuanlegen im selben Commit) | `TestRenamesLinear` („Ende in einem Commit gelöscht und angelegt: Inhaltsänderung an Ort und Stelle“); `TestRenamesMerge` (überschreibende Seite) |
| 25 | Boundary (Art des Eintrags, Verzeichnis) | `TestRenamesArtDesEintrags` (Ausführbar-Bit, Symlink gegen Datei, Symlink zu Symlink, Verzeichnis, Gitlink); Kern „Verzeichnis umbenannt: …“ (beide Hälften) |
| 26 | Negative (immutable Datei umbenannt) | `TestVCSFollowRenamesImmutableDateiSelbst`, `TestVCS_FollowRenamesImmutableDateiSelbst` |
| 27 | Boundary (Basis kein Vorfahr) | `TestRenamesMerge` („Basis kein Vorfahr von HEAD“), `TestRenamesFehlerFormen`, `TestVCS_FollowRenamesBasisKeinVorfahr` |
| 28 | Boundary (`--staged`) | CLI `TestVCS_FollowRenamesStaged` (beide Hälften, drei Konfigurationen); Adapter `TestRenamesStaged` |
| 29 | fail-closed (Zwischen-Commit) | `TestRenamesUnlesbaresObjekt` (alle vier Objekte: Commit eines Zwischen-Commits, Baum eines Schritts, Baum des Eltern-Commits vor BASE, Teilbaum); CLI `TestVCS_FollowRenamesFailClosed` (Exit 2 mit, Exit 0 ohne den Schlüssel); `TestVCSFollowRenamesPortMethode` |
| 30 | Boundary (ohne `vcs.follow-renames` byte-identisch) | `TestVCS_FollowRenamesAusByteIdentisch`, `TestVCS_FollowRenamesFailClosed` (ohne Schlüssel liest nichts), Zähler der Port-Aufrufe in `TestVCSFollowRenames`; die Vorfassung unten |

**„Ohne Schlüssel byte-identisch“.** Der Test belegt nur die Gleichheit zweier Konfigurationen im neuen Code
(Befund I-3 von R1); die Vorfassung trägt auch `make blackbox-probe` nicht (keine vcs-Fixture, siehe Tabelle
oben). Deshalb habe ich Vorher- und Nachher-Image (`d-check:probe-vorher` aus `6aca94b0`, `d-check:latest`)
auf identischen Eingaben verglichen, stdout, stderr und Exit getrennt per `cmp`:

| Eingabe | Ergebnis |
|---|---|
| Wegwerf-Klon, `--enable vcs --range HEAD~1..HEAD`, ohne Schlüssel (58 `core-drift-vcs`) | identisch (Exit 1, 536 Zeilen) |
| derselbe Klon, `ignore-link-targets` allein | identisch (Exit 1, 536 Zeilen) |
| derselbe Klon, `--json` | identisch (Exit 1, 4300 Zeilen) |
| derselbe Klon, `--staged` | identisch (Exit 1, 478 Zeilen) |
| dieses Repo, 100-Commit-Range mit Merge (siehe Risiko 3) | identisch (Exit 0) |
| dieses Repo, `--range HEAD~1000..HEAD --json` | identisch (Exit 0, 11 Zeilen) |

Die eine gewollte Ausnahme, Gitlink statt immutabler Datei mit `--staged`, ist durch V1 am Vorzustand gemessen
(Exit 2 `object not found`, nachher Exit 1 `core-drift-vcs`) und lockert nichts. Sie steht im Plan §3, im
Spezifikations-Satz und im Welle-Trigger
([`welle-96`](../plan/planning/welle-96-vcs-folgt-umbenennung.md), „gewollte Abweichung im gestagten Modus“).

**Zusätzlich, ohne die Test-Doppelgänger:** acht Szenarien mit echtem `git` und dem gebauten Image, Exit und
Zahl der `core-drift-vcs` je (ohne Schlüssel / `follow-renames` / beide):

| Szenario | erwartet nach Lastenheft | gemessen |
|---|---|---|
| Kette über zwei Merge-Seiten (X → Y, Y → Z), Nachzug X → Z | Drift | 1 / 1 / 1 |
| Rotation (A → X, dann C → A), Nachzug A → X | folgt | 1 / 0 / 0 |
| Rotation, Nachzug C → A | Drift | 1 / 1 / 1 |
| zwei Quellen auf ein Ende (Merge), Umbiegen A → B | Drift | 1 / 1 / 1 |
| Tausch über Zwischennamen, Nachzug A → B | Drift | 1 / 1 / 1 |
| Kette über drei Commits mit Inhaltsänderung, Nachzug A → C | folgt | 1 / 0 / 0 |
| Nachzug in einem **früheren** Commit als der Move | folgt | 1 / 0 / 0 |
| Umbenennen und im selben Commit ändern | Drift | 1 / 1 / 1 |
| immutable Datei selbst umbenannt | Drift | 1 / 1 / 1 |
| Basis kein Vorfahr (Seitenzweig als BASE) | wie ohne Schlüssel | 1 / 1 / 1, kein Exit 2 |

Jede Zeile stimmt mit dem Lastenheft überein.

### 3 — Spiegel: bestätigt

Die Suche aus Plan §3 nach `ignore-link-targets` und `IgnoreLinkTargets` ausgewertet: berührt sind
`spec/spezifikation.md`, `harness/sensors/adr-check.md`, `docs/user/benutzerhandbuch.md`, die Konfiguration
(`internal/hexagon/core/model/config.go`, `internal/adapter/driven/configyaml/configyaml.go`) und das
`--print-config`-Template. Nicht berührt, wie vorgesehen: README und README.de (Release-Prep), `.d-check.yml`,
ADR-0115 und der Index (Architect, Closure) sowie `AGENTS.md` und `harness/README.md`; die Spiegel-Dateien
von `ignore-link-targets` in `docs/plan/planning/done/` und `docs/reviews/` tragen Historie.

- **`--print-config`:** die Zeile steht im gebauten Image (Zeile 128) und gleicht dem Code.
- **Handbuch:** sieben Grenz-Aussagen am Code oder am Image geprüft: (1) ohne Schlüssel liest das Modul keinen
  Commit zwischen BASE und HEAD — `TestVCS_FollowRenamesFailClosed`; (2) ein Gitlink und ein Verzeichnis nehmen
  nicht teil — V9, V10; (3) Basis kein Vorfahr ohne Exit 2 — V3 und Szenario; (4) mit `--staged` einziger Schritt
  HEAD gegen Index — `stagedRenames`; (5) unlesbares Objekt ⇒ Exit 2 — V6; (6) „liest nur die Kennung“ — kein
  Blob wird im Adapter geöffnet; (7) Gitlink ersetzt immutable Datei ⇒ `core-drift-vcs` — V1, V9. Alle stimmen.
- **Sensor-Datei:** Grenze 5 und der geänderte Absatz zum Gitlink stimmen mit dem Code und
  `TestVCS_ImmutableDateiDurchGitlinkErsetzt` überein (vier Kombinationen, beide Modi).
- **Spezifikation:** der Nebensatz zur Index-Enumeration und der Satz vom „benannten Preis“ sind ersetzt, mit
  Historie-Zeile; die Bedingung des Plans (der Testfall bestätigt `core-drift-vcs` in beiden Modi) ist erfüllt.
  Kein neuer Grund-Code; §4 unberührt.
- **Abdeckungs-Datei:** `docs/user/abdeckung-tests.md` ist abgeleitet, `make test` hält sie; sie ist im Lauf grün.

**Nicht Teil dieses Punktes und offen:** `.d-check.yml` schaltet den Schlüssel noch nicht ein, ADR-0115 steht auf
`Proposed`, `AGENTS.md` §3.5 und `harness/README.md` nennen sie noch nicht. Der Plan bindet das an **einen**
Commit des Architect (Entscheid 14); es ist Closure-Arbeit, siehe Gesamturteil.

### 4 — `make gates`, `make coverage-gate`, `make image-test`: bestätigt

Siehe die Tabelle oben; alle drei selbst gefahren, Exit 0. Coverage 94,90 % gegen Schwelle 93 %.

### 5 — Reviews: bestätigt, mit Vorbehalt zum Register

[R1](2026-11-24-slice-292-vcs-umbenennung-r1.md): 0 HIGH, 2 MEDIUM (M-1, M-2), 2 LOW (L-1, L-2), 4 INFO.
[R2](2026-11-24-slice-292-vcs-umbenennung-r2.md): kein Befund, M-1 und M-2 geschlossen. Beide Mutationen
(`w.excluded[h] ||` entfernt, doppelte Marke einfach) sind in R2 gefahren und an den neuen Tests rot; M-1 und
M-2 sind damit **nicht** nur „behauptet“. Ich habe die Gegenrichtung selbst gefahren: die zwölf Mutationen
oben, darunter keine, die einen der beiden neuen Tests bräuchte. Es steht **kein offener HIGH oder MEDIUM**.

**Vorbehalt:** L-1, L-2 und I-1 bis I-4 „gehen ins Register“ — der Diff `6aca94b0..879e576d` berührt
`docs/plan/planning/observations/` nicht. Das gehört zur Closure (DoD-Punkt „Beobachtungs-Register
fortgeschrieben“), ist aber **noch nicht geschehen**.

### 6 — Das Mehr gegenüber der Welle, am gebauten Image: bestätigt

`make build`, dann ein Klon von `879e576d` im Scratchpad, ein Commit: `git mv spec/spezifikation.md
spec/specification.md` und in 59 `Accepted`-ADRs nur das Link-Ziel (Inline-Link und Referenz-Definition) auf
`spec/specification.md` gezogen, die Konfiguration des Klons auf die neue Datei umgestellt (sie nennt die Datei als
`ids`-Ziel). Aufruf je Lauf: `docker run --rm --network none -v <klon>:/repo:ro d-check:latest --enable vcs
--range HEAD~1..HEAD`; gezählt ist `core-drift-vcs` (die übrigen Module laufen mit und melden in allen vier Läufen
dieselben 478 Zeilen, bei denen die Klon-Konfiguration die Ursache ist).

| Lauf | Konfiguration im Klon | `core-drift-vcs` | Exit |
|---|---|---|---|
| A | die des Repos (`ignore-link-targets: true`), ohne den Schlüssel | **58** | 1 |
| B | `ignore-link-targets: true` **und** `follow-renames: true` | **0** | 1 (übrige Module) |
| C | `follow-renames: true` allein | **0** | 1 (übrige Module) |
| D | kein Schlüssel | **58** | 1 |
| E | Gegenprobe: Datei gelöscht und mit **anderem Inhalt** neu angelegt, derselbe Nachzug, beide Schlüssel | **58** | 1 |

Von den 59 geänderten ADRs meldet eine nichts, auch nicht ohne Schlüssel: bei ADR-0103 steht der Verweis nur im
ausgenommenen Abschnitt `Geschichte`. A und D sind außerhalb von `core-drift-vcs` Zeile für Zeile gleich, B und C
ebenso. Der Wegwerf-Schalter wurde nur im Klon gesetzt; im Repo steht `.d-check.yml` unverändert.

### 7 — Risiken aus §6

| # | Risiko | Vorgeschlagener Ausgang | Beleg |
|---|---|---|---|
| 1 | der Auflöser trägt eine zweite Gleichsetzung, Filter könnten verrutschen | **entfallen** | `internal/hexagon/core/rules/vcs_pfad_nachzug_test.go` im Diff nicht berührt und in `make test` grün (Exit 0); beide Schalter zusammen fahren `TestVCSFollowRenames` (Spalte „beide“), `TestVCSFollowRenamesRandfaelle`, `TestVCSFollowRenamesBeideSchluesselIgnoreBleibt`, `TestVCS_FollowRenames` |
| 2 | racily-clean-Fixtures bei gleichen Blobs | **entfallen** | die Adapter-Fixtures legen Objekte im Speicher an (kein Arbeitsbaum, kein Index); die CLI-Fixtures schreiben über `writeAt`, das `coretest.GitFixtureRewriteHazard` ruft, und kein neuer Fall schreibt dieselbe Datei gleich lang neu (jede Neufassung ändert den Link-Pfad und damit die Länge; `mvAt` legt einen neuen Pfad an) — [`BEO-ALL/racily-clean-git-fixture`](../plan/planning/observations/BEO-ALL/racily-clean-git-fixture/observation.md) |
| 3 | der Walk liest mehr, ein großer Push-Bereich wird langsamer | **entfallen** | gemessen, siehe unten: Median 1,14 s ohne, 1,56 s mit, **1,37×**, Schwelle 2× |
| 4 | die Angleichung der Index-Enumeration verschiebt ein Ergebnis mit `--staged` | **entfallen** | der Testfall „immutable Datei durch Gitlink ersetzt“ läuft in beiden Modi und beiden Konfigurationen grün; Rot am Vorzustand gemessen (V1: Exit 2 statt 1); der Closure-Trigger von [`welle-96`](../plan/planning/welle-96-vcs-folgt-umbenennung.md) benennt die Abweichung |

**Messung zu Risiko 3.** BASE `9e7b4017` (eine Merge-Seite zweigt am Eltern-Commit `7760865c` ab, der **vor** BASE
liegt), HEAD `b8546e6c`; `git rev-list --count BASE..HEAD` = 100, die Range enthält genau einen Merge, `25394098`.
Die Historie dieses Repos trägt drei Merges (`git log --merges`); nur dieser liegt in einer Range, in der eine
Seite vor einem wählbaren BASE abzweigt, ohne dass die Range sehr lang wird. Kommando je Lauf, Image `d-check:latest`,
Klon des Repos als Mount, Konfiguration des Repos (`ignore-link-targets: true`), mit Schlüssel eine Zeile
`follow-renames: true` mehr:

```bash
docker run --rm --network none -v <klon>:/repo:ro d-check:latest --enable vcs \
  --disable links --disable anchors --disable ids --disable matrix --disable codepaths \
  --disable spans --disable hostpaths --disable versions --disable structure \
  --disable diagrams --disable citations --disable file --disable pinsync \
  --range 9e7b4017daa0550bc3ca34e48672e1bb0da2e572..b8546e6c21a48ab4b93975ca54d7ff9f1679273b
```

Das ist der Aufruf, den `make adr-check RANGE=…` ausführt (`make -n` gelesen); ich habe ihn direkt gefahren,
weil `make adr-check` das Image zuerst baut und die Bauzeit (rund 10 s) die Messung überdeckt. Fünf Läufe je
Konfiguration, abwechselnd; beide Konfigurationen enden mit Exit 0 und `0 Befund(e)`:

| Konfiguration | Läufe (s) | Median (s) |
|---|---|---|
| ohne `follow-renames` | 1,163 · 1,128 · 1,139 · 1,177 · 1,272 | **1,139** |
| mit `follow-renames` | 1,510 · 1,300 · 1,562 · 1,598 · 1,799 | **1,562** |

Verhältnis 1,37. Als Gegenprobe `make adr-check RANGE=…` in beiden Konfigurationen: `real` 11,2 s ohne, 11,4 s
mit — von der Bauzeit dominiert, kein Gegenargument. R1 hatte als Einordnung höchstens das 1,35-Fache auf Ranges
bis 1000 Commits ohne die Merge-Bedingung gemessen; die hier gemessene Bedingung (Baum eines Eltern-Commits vor
BASE) erhöht die Zahl nicht über die Schwelle.

### 8 — Abgrenzung §1: bestätigt

`git diff 6aca94b0..HEAD --stat` über `docs/plan/adr/0115-vcs-folgt-umbenennung.md`, `docs/plan/adr/README.md`,
`.d-check.yml`, `AGENTS.md`, `harness/README.md`, `README.md`, `README.de.md` und `CHANGELOG.md`: leer. Der
Handbuch-Kopf ist nicht berührt (Diff nur im Abschnitt zum Modul `vcs`, in der Konfigurations-Skizze und in der
Zeile der Modul-Tabelle). Die Metriken, `tracked` und `commits` bleiben unberührt (`LineChanges` ohne Rename-Erkennung).
Die 18 Dateien der Stat-Liste stehen alle in Plan §3.

## Plan-vs-Code-Diff

| Plan §3 | Code | Abweichung |
|---|---|---|
| Port: Methode `Renames`, Zuordnung Quelle → Ende, leer bei „Basis kein Vorfahr“ und ohne HEAD | `Renames(base, head) (map[string]string, error)` | keine |
| Adapter: eigener Blob-Vergleich, Teilbäume gleicher Kennung ungelesen, kein Diff-Walker, Index-Angleichung | `diffFiles`/`diffEntry`/`isFileMode`; Teilbaum-Skip bei gleicher Kennung | keine; der Skip ist durch keinen Test gebrochen (I-1 von R1, im Ergebnis äquivalent) |
| `VCSConfig.FollowRenames`, `rawVCS` wie `ignore-link-targets` | beides, ein Test fährt beide Schlüssel nebeneinander | keine |
| Kern: Auflöser baut die Pfad-Vereinigung auch mit `follow-renames` allein | `newTargetPolicy` baut `tree` bei jedem der beiden Schlüssel | keine |
| `--print-config` | Zeile im `vcs`-Block | keine |
| Handbuch, Sensor-Datei | beide aktualisiert | keine |
| Spezifikation: zwei Sätze Schritt 2 + Historie-Zeile | beide, Historie 2026-11-24 | keine |
| Testdateien, Test-Doppelgänger des Ports | `git_renames_test.go`, `cli_vcs_renames_test.go`, `vcs_follow_renames_test.go`, `fakeVCS.Renames` | siehe V-1 |
| `.d-check.yml`, ADR-0115 `Accepted`, Index, `AGENTS.md` §3.5, `harness/README.md` | **nicht in der Range** | gewollt, Closure (Architect-Akt, Entscheid 14) |
| READMEs, CHANGELOG | unberührt | gewollt, Release-Prep |

## Findings der Verifikation

Kein HIGH, kein MEDIUM. Zwei Hinweise unter der Schwelle, beide ohne Folge für das Urteil:

- **V-1 (LOW):** Der Plan nennt für den Happy Path „gespiegelt der Nachzug in einem früheren Commit“. Kein Test
  fährt ihn. Es gibt auch keinen Code-Pfad, der ihn anders behandelt als den Fall „Move zuerst“ — der Kern
  vergleicht BASE- und HEAD-Fassung der immutablen Datei, die Zuordnung ist range-weit —, und mein Szenario
  (Nachzug in einem früheren Commit als der Move, drei Konfigurationen) liefert das Erwartete (0 Befunde mit dem
  Schlüssel). Ein Testfall ist billig und gehört ins Register.
- **V-2 (LOW, übernommen):** die zwei Lesepfade aus R1 L-1 (Baum einer Merge-Seite ohne Kind-Schritt, Art eines
  gestagten Eintrags) und L-2 (Neuanlage lässt die Linie wegfallen) tragen keinen Rot-Beleg; für L-2 habe ich
  die Mutation selbst gefahren: grün. Beide Punkte liegen **über** den 30 Kriterien (Sätze der Spezifikation, nicht
  eines der Akzeptanzkriterien) und sind ins Register zu tragen, wie R1 es vorsieht.

## Gesamturteil

**Der Code und die Spiegel dieses Slice sind verifiziert; der Slice ist closure-reif unter den folgenden
Closure-Handlungen, die in dieser Range nicht stehen und nicht Gegenstand der Verifikation waren.** Alle acht
geprüften Punkte sind bestätigt; die 30 Kriterien tragen je einen Test, jede der zwölf Mutationen, die ich gefahren
habe, außer der einen, die R1 bereits als L-2 führt, ist rot, und ohne den Schlüssel ist die Ausgabe gegen das
Vorher-Image byte-identisch. Bleibt für die Closure:

1. **Architect-Akt, ein Commit** (Entscheid 14): ADR-0115 auf `Accepted`, Index, `.d-check.yml` mit
   `follow-renames: true`, `AGENTS.md` §3.5 und `harness/README.md` §Traceability rules. Ohne ihn ist DoD-Punkt 3
   des Plans (Zeile zur Konfiguration, `AGENTS.md`, `harness/README.md`) nicht abgehakt.
2. **Register:** L-1, L-2, I-1 bis I-4 aus R1 und V-1 hier ins Beobachtungs-Register tragen (eine Antwort ist
   auch „keine Beobachtung angefallen“, aber hier sind welche angefallen).
3. **Closure-Notiz** mit Lerneintrag; die Messung zu Risiko 3 (Median 1,14 s gegen 1,56 s, Faktor 1,37) und die vier
   Ausgänge aus §7 oben in die Notiz übernehmen.
4. Der Welle-Trigger „Der Anlass trägt am gebauten Image“ ist hier schon gefahren (Tabelle zu Punkt 6); er gehört
   in `welle-96-results.md`, nicht in diesen Slice.

Ich habe im Repo nichts außer diesem Report geschrieben; Wegwerf-Klone, Mutationen und Logs liegen im Scratchpad.
