# Review-Report: slice-220 — 2026-09-17 (R2)

**Review-Art:** Code-Review — Runde 2. Geprüft wird (a) die inhaltliche
Vollständigkeit der R1-Fixes und (b) ein eigenständiger Durchlauf gegen den
gesamten Diff seit der Beanspruchung von slice-220.

**Gegenstand:** Diff-Range `769a2bf0..800a3166` (Claim-Commit bis Fix-Commit),
mit Fokus auf `800a3166` (fix(vcs): Review-Runde 1 einarbeiten). Der
Feature-Commit `3100e2f8` wurde bereits in R1
(`docs/reviews/2026-09-17-slice-220-vcs-pfadmenge-review-r1.md`) geprüft und
ist hier Teil des eigenständigen Zweitdurchlaufs, nicht erneut von Grund auf
neu bewertet.

**Skill:** `.harness/skills/reviewer.md` @ v1.16.0 / 2026-09-07 ·
**Modell:** `claude-sonnet-5` · **Datum:** 2026-09-17

**Eingangs-Kontext:**

- R1-Report `docs/reviews/2026-09-17-slice-220-vcs-pfadmenge-review-r1.md`
  (vollständig gelesen — Findings R1-F-1..F-3, alle Negativbefunde, alle
  Verifikationsläufe)
- Slice-Plan `docs/plan/planning/in-progress/slice-220-vcs-pfadmenge-statt-diff.md`
- `CO-001` (`docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md`)
- Beobachtungs-Register: `BEO-ALL/fix-schliesst-pfad-nicht-klasse`
- `DC-FA-VCS-001`/`DC-FA-VCS-001.a` (`spec/lastenheft.md`,
  `spec/spezifikation.md`)
- `ADR-0024` (Accepted; Schärft-Ziel `DC-FA-VCS-001.a`)
- `AGENTS.md` §3 (Hard Rules), §6 (Workflow)

---

## Verifikationen dieses Laufs

Selbst ausgeführt (Docker/`make`-only, kein Host-Go), gegen den Stand
`800a3166`:

- `make test`: grün, alle elf Pakete — deckt sich mit der Fix-Commit-Botschaft.
- `make gates`: grün — `baseline-verify + workflow-pins + doc-check + lint +
  test + arch-check + coverage-gate + semgrep + gate-consistency +
  planning-check`; `coverage-gate` meldet **94.70 %** (Schwelle 93 %),
  identisch zu R1/zum Feature-Commit; `targets`/`planning`-Läufe je
  **798 Datei(en) geprüft, 0 Befund(e)** — deckt sich exakt mit der
  Fix-Commit-Botschaft („zehn Gates, 798 Dateien, 0 Befunde").
- `make doc-check` **separat und mit dem vollen Default-Modulsatz** (nicht nur
  den in `gates` isoliert gefahrenen `targets`/`planning`-Läufen) ausgeführt,
  um `matrix`/`ids`/`links` explizit gegen die beiden neuen
  `<!-- d-check:status-provenance -->`-Marker in `spec/spezifikation.md`
  laufen zu lassen: **798 Datei(en) geprüft, 0 Befund(e)** — kein
  `matrix-forbidden`, keine Anker-/Link-Lücke.
- Code-gegen-Spec-Textvergleich Satz für Satz: neue Fassung von
  `spec/spezifikation.md` §`DC-FA-VCS-001.a` Schritt 2 gegen
  `internal/hexagon/core/rules/vcs.go` (`CheckVCS`, `protectedSet`,
  `vcsDeleted`) und `internal/adapter/driven/git/git.go` (`AllPaths`,
  `pathsAt`, `walkTree`) am Stand `800a3166` (s. u., Abschnitt R1-F-1).
- Repo-weiter Grep auf verbliebene Referenzen der abgelösten Diff-Typen
  (`ChangedPaths`, `VCSChange`, `VCSStatus`, `diffTrees`, `diffTreeIndex`,
  „Rename-Erkennung", „Diff-Übersetzung", „Diff-Status") außerhalb von
  `docs/reviews/`, historischen Slice-/Historie-Einträgen und CR-Dokumenten
  (s. u., Abschnitt R1-F-2/F-3).

## Prüfung R1-F-1 (MEDIUM) — Spec-Nachzug

**Bestätigt behoben.** Die neue Fassung von Schritt 2
(`spec/spezifikation.md:1762-1799`) wurde Satz für Satz gegen den Code am
Stand `800a3166` geprüft:

| Spec-Aussage (neu) | Code-Beleg |
|---|---|
| „Die Klasse `vcs.paths` … wird direkt gegen beide Tree-Stände aufgelöst, statt einem Diff zu vertrauen" | `CheckVCS` ruft `vcs.AllPaths(base, head)`, kein `ChangedPaths`/Diff mehr (`vcs.go:34`) |
| „der VCS-Port enumeriert jeden Datei-Pfad (regulär/ausführbar/symlink) … über einen vollständigen Durchlauf beider Bäume" | `walkTree` (`git.go:117-137`), `switch e.Mode { case filemode.Regular, filemode.Executable, filemode.Symlink: … }` |
| „Ein Unterbaum, dessen Objekt nicht lesbar ist, bricht diesen Schritt fail-closed ab" | `walkTree`, `case filemode.Dir: sub, err := repo.TreeObject(e.Hash); if err != nil { return fmt.Errorf(…) }` (`git.go:124-128`), propagiert über `pathsAt`/`AllPaths` als `error` |
| „Ein Tree-Eintrag ohne Datei-Inhalt (Verzeichnis, Gitlink) … wird beim Durchlaufen übersprungen" | `walkTree`-`switch` hat für `filemode.Submodule` keinen `case` — fällt durch, ohne Eintrag in `out` |
| „ein Pfad in beiden gefilterten Mengen → Core-Vergleich; nur in BASE → Pfad-Stabilitäts-Prüfung; nur in HEAD → frei" | `CheckVCS`: iteriert `paths` (aus `baseSet`), `if headSet[p] { vcsModified(...) } else { vcsDeleted(...) }`; ein Pfad nur in `headSet` wird nie besucht (`vcs.go:40-60`) |
| „Eine Umbenennung ist … strukturell Löschung(alt) + Hinzufügung(neu)" | `baseAll`/`headAll` unabhängig aus je einem eigenen Tree-Walk, keine Rename-Erkennung im Adapter (`git_test.go::TestAllPathsPureRenameYieldsDelete` bestätigt: alter Pfad nur in `baseAll`, neuer nur in `headAll`) |
| „Eine unlesbare BASE lässt bereits die Enumeration selbst fehlschlagen, bevor irgendein Pfad klassifiziert wird" | `AllPaths`: `baseAll, err := a.pathsAt(base)` mit sofortigem `return` bei Fehler, **vor** `headAll`/Filterung (`git.go:84-92`) |

Der neue Historie-Eintrag (`spec/spezifikation.md:3334`) wurde ebenfalls
geprüft: Er benennt korrekt, was ersatzlos entfällt (Diff-Status-Klassen
`M`/`T`/`D`/`R`/`A`, Rename-Erkennungs-Absatz, nachgelesene BASE-Lesbarkeit
für `A`), was unverändert bleibt (Gitlink-Zusage, fehlende
Inhalts-Ähnlichkeits-Messung) und den `CO-001`-Anlass korrekt (dritte
Ausprägung *ohne* Pendant, plus die vierte, zuvor fehldiagnostizierte
Ausprägung mit `core-drift-vcs`/Exit 1 statt Exit 2) — deckt sich mit der
empirischen R1-Gegenprobe und mit den neuen End-to-End-Tests
`TestVCS_UnlesbareObjekte` (vier Unterfälle, `cli_vcs_test.go:242-333`). Kein
Lastenheft-Bump nötig — bestätigt: `spec/lastenheft.md` ist im gesamten Diff
`769a2bf0..800a3166` unverändert (letzter berührender Commit ein unabhängiger
STRUCT-001-Fix).

## Prüfung R1-F-2/R1-F-3 (LOW) — stale Kommentare

**Beide bestätigt behoben, keine neuen Stale-Referenzen hinterlassen.**

- `internal/hexagon/core/rules/vcs.go:76-79` (`vcsDeleted`-Doc-Kommentar):
  nennt jetzt korrekt „die Mengen-Differenz in `CheckVCS` (der Pfad fehlt in
  `headSet`)" statt der nicht mehr existierenden „Diff-Übersetzung im
  VCS-Adapter". `headSet` ist tatsächlich die dafür zuständige Variable
  (`vcs.go:39`).
- `internal/adapter/driving/cli/cli_vcs_test.go:95`: `ChangedPaths` durch
  `AllPaths` ersetzt — die tatsächlich aufgerufene Methode.
- Repo-weiter Grep (`ChangedPaths`, `VCSChange`, `VCSStatus`, `diffTrees`,
  `diffTreeIndex`) findet außerhalb von Review-Reports und historischen
  Historie-/Slice-Einträgen nur noch `internal/hexagon/port/driven/vcs.go:32`
  (`ChangedPaths` als bewusster Herkunfts-Zeiger im `AllPaths`-Doc-Kommentar,
  bereits von R1 geprüft und für korrekt befunden — unverändert seit
  `3100e2f8`) sowie `internal/hexagon/core/rules/vcs_test.go:281`
  (`TestVCSStatusImFenceZaehltNicht` — der Name bezieht sich auf die
  fachliche **Status-Zeile** einer Datei, nicht auf den entfernten Typ
  `VCSStatus`; keine Verwechslungsgefahr, kein stale Kommentar).

## Eigenständiger Zweitdurchlauf gegen den Gesamt-Diff (`769a2bf0..800a3166`)

Alle 18 Prüffragen erneut gegen `vcs.go`, `git.go`, `git_test.go`,
`cli_vcs_test.go`, `vcs_test.go`, `port/driven/vcs.go`,
`port/driven/workflow.go` und `spec/spezifikation.md` durchgegangen. Keine
neuen HIGH/MEDIUM/LOW-Findings über die bereits von R1 gefundenen und jetzt
behobenen hinaus. Details in den Negativbefunden unten; hervorzuheben:

- **Frage 2 (Kern-Modul meldet falsch)** erneut mit eigener Herleitung
  geprüft (nicht nur R1s Empirie übernommen): Ein Pfad mit unlesbarem
  **Blob** (nicht Unterbaum) an BASE wird von `walkTree` dennoch korrekt in
  `baseAll` gelistet (die Tree-**Struktur**-Enumeration braucht das Blob
  nicht), landet also in `baseSet`/`headSet` wie ein normaler Pfad und wird
  anschließend über `vcsModified`/`vcsDeleted` → `FileAt` behandelt — dort
  greift die seit slice-218 bestehende `entryUnreadable`-Unterscheidung
  weiterhin fail-closed. Es gibt also **keinen** Pfad, auf dem ein
  unlesbares Blob fälschlich als „Added" (frei) durchrutschen könnte, obwohl
  der neue Entwurf die alte, alleinige BASE-Nachlese für den `A`-Fall
  ersatzlos strich — die Eigenschaft ist strukturell erhalten, nicht nur
  zufällig getestet.
- **Frage 11 (zwei Module derselben Eingabe-Klasse verschieden)**: `commits`
  (`CheckCommits`/`CommitMessages`) und `tracked` (`CheckTracked`/
  `TrackedPaths`) wurden nicht angefasst und bleiben von diesem Umbau
  unberührt; keine neue Inkonsistenz zwischen den drei VCS-Port-Konsumenten.
- **Frage 14 (Provenance-Marker begründet statt zeigt)**: beide neuen
  `<!-- d-check:status-provenance -->`-Marker (Fließtext und Historie-Zeile)
  zeigen den Ursprungs-Slice, ohne eine Entscheidung zu begründen — dieselbe,
  bereits etablierte Verwendungsform wie beim `slice-225`-Marker im
  benachbarten Historie-Eintrag.

## Negativbefunde

- **Hexagon-Grenze** (ADR-0005/ADR-0012): geprüft — der Fix-Commit ändert
  keine Import-Struktur; `make arch-check` (Teil von `make gates`): 0
  Befunde. Kein Befund.
- **Inline-Suppression / Gate-Lockerung ohne ADR** (`AGENTS.md` §3.2/§3.6):
  geprüft — kein `//nolint`, keine Schwellen-Senkung im Fix-Commit. Kein
  Befund.
- **Netzzugriff außerhalb `external`** (`DC-QA-03`): geprüft — Fix-Commit
  ändert nur Doc-Kommentare, einen Test-Kommentar und Spec-Prosa; kein
  neuer Code-Pfad. Kein Befund.
- **Kommentar-Klassen** (`AGENTS.md` §3.7): der geänderte `vcsDeleted`-Kommentar
  und der geänderte Test-Kommentar beschreiben weiterhin nur, was da ist
  (Zusage/Kopplung), keine Review-Historie, keine Slice-Nummern im
  Kommentar-Text selbst (der Provenance-Marker im Spec-Text ist eine
  eigene, dafür vorgesehene Form, keine Kommentar-Klasse nach §3.7). Kein
  Befund.
- **Referenz-Richtung (SDP) / Marker-Ehrlichkeit**: geprüft — beide neuen
  `d-check:status-provenance`-Marker zeigen Herkunft, begründen keine
  Entscheidung; `make doc-check` mit vollem Default-Modulsatz (798 Dateien)
  bestätigt zusätzlich maschinell 0 `matrix-forbidden`-Befunde. Kein Befund.
- **Botschaft/Historie-Eintrag verallgemeinert über die Messung hinaus**
  (Anker 8): der Fix-Commit-Botschaft und der Historie-Eintrag nennen
  „zehn Gates, 798 Dateien, 0 Befunde" — selbst nachgefahren, deckungsgleich.
  Keine weiterreichende Schlussfolgerung (z. B. „damit ist jede
  Diff-Abhängigkeit im Repo behoben") wird behauptet; Abgrenzung 2 des
  Slice-Plans benennt explizit, dass andere Module nicht geprüft wurden.
  Kein Befund.
- **Grenzen-Liste ohne größte Lücke** (Anker 18): Abgrenzung 4 des Slice-Plans
  benennt bewusst offen, dass die Pfad-Auflösung teurer sein kann als der
  Diff, ohne ein Performance-Ziel zu setzen — das ist die im Vertragstext
  selbst genannte Grenze, keine verschwiegene. Kein Befund.
- **`spec/lastenheft.md` unverändert, obwohl Verhalten sich ändert**: geprüft
  — bewusst, weil dieselbe Zusage (Löschung/Umbenennung einer immutablen
  Datei wird erkannt) über einen anderen Mechanismus eingelöst wird, keine
  neue oder geänderte Anforderung. Konsistent mit dem Vorgehen der
  2026-09-08/2026-08-31-Einträge zu derselben Anforderung. Kein Befund.
- **DoD-Checkboxen des Slice-Plans**: bleiben nach dem Fix-Commit
  unverändert unabgehakt — korrekt, da der Slice noch nicht geschlossen ist;
  DoD-/Closure-Konformität ist Aufgabe der Verifikation, nicht dieses
  Reviews (Anti-Pattern „Kein Verifier"). Kein Befund an dieser Stelle.
- **`CO-001` nicht aktualisiert**: die Carveout-Datei nennt weiterhin die
  dritte Ausprägung als „offen" — korrekt für den aktuellen Stand, da die
  Carveout-Auflösung an das **Release** gebunden ist (Auflösungs-Trigger:
  „Der Klassen-Fix aus slice-220 ist drin UND das nächste Release ist
  veröffentlicht"), nicht an den Merge des Fixes. Aktualisierung gehört zur
  Closure/zum Release-Slice `slice-219`, nicht zu diesem Fix-Commit. Kein
  Befund.
- **`CHANGELOG.md`**: unverändert — korrekt nach `AGENTS.md` §5
  (Release-Prep, nicht Feature-/Fix-Commit). Kein Befund.
- **Historie-Tabellen-Reihenfolge** (`spec/spezifikation.md` §7): der neue
  Eintrag steht als oberste (jüngste) Zeile unter den drei 2026-09-17-Einträgen;
  reverse-chronologisch konsistent zum Bestand. Kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** keine neuen.

## Verdikt

**Freigegeben.** R1-F-1 (MEDIUM) ist inhaltlich vollständig und korrekt
behoben — Satz-für-Satz-Vergleich der neuen Spec-Fassung gegen den
tatsächlichen Code am Stand `800a3166` deckt sich lückenlos, der neue
Historie-Eintrag ist präzise. R1-F-2/R1-F-3 (LOW) sind behoben, ohne neue
Stale-Referenzen zu hinterlassen. Der eigenständige Zweitdurchlauf gegen den
gesamten Diff seit der Beanspruchung (`769a2bf0..800a3166`), einschließlich
aller 18 Prüffragen des Reviewer-Skills, findet keine neuen Findings.
`make test` und `make gates` wurden selbst ausgeführt und bestätigen die in
der Fix-Commit-Botschaft behaupteten Werte (798 Dateien, 0 Befunde, Coverage
94.70 %). Aus Review-Sicht steht der Closure von slice-220 nichts entgegen —
verbleibende Schritte (DoD-Abhaken, Closure-Notiz, Register-Fortschreibung,
Risiko-Ausgänge, `CO-001`/`slice-219`-Nachzug bei Release) sind Sache der
Verifikation und der Closure, nicht dieses Reviews.
