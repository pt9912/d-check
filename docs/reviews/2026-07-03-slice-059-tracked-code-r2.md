# Review-Report — slice-059 (Modul `tracked`) — Code-Review (R2)

**Datum:** 2026-07-03
**Reviewer-Rolle:** unabhängig/adversarial, Fokus **Code** (Commit `d39dd3f`).
**Gegenstand:** neues opt-in Modul `tracked` (`DC-FA-TRK-001`,
[ADR-0030](../plan/adr/0030-tracked-referenz-ziele.md)) — Port-Erweiterung
`TrackedPaths()` (`internal/hexagon/port/driven/vcs.go`,
`internal/adapter/driven/git/git.go`), Regelmodul
`internal/hexagon/core/rules/tracked.go`, `run.go`-Verdrahtung, `configyaml`
(`rawTracked`/`applyTracked`), CLI (`resolveVCS` ohne Range-Pflicht,
`--print-mk doc-tracked`, `config_template`, `suggest`), Tests.
**Bindend:** `spec/lastenheft.md` §DC-FA-TRK-001 (7 AK),
`spec/spezifikation.md` §DC-FA-TRK-001.a (Schritte 1–6), ADR-0030,
`AGENTS.md` §3, Slice-Plan §2/§4.
**Baseline:** `.harness/skills/reviewer.md` v1.2.0.
**NICHT geprüft:** DoD-Abhakung (Verifikations-Rolle), Doku-Treue der
Spec-/Handbuch-Prosa (R1), CI-Pipeline-Lauf.

## Verifikations-Proben (eigene, adversarial — Baum nach jeder Probe restauriert)

Alle Läufe gegen das gebaute Image (`docker run --rm --network none
-v …:/repo:ro d-check:latest`) auf Wegwerf-Fixture-Repos bzw. via
`make test`/`make arch-check`; Abschluss-Zustand: `git status` sauber
(nur die Review-Reports neu), `make test` grün, `make arch-check`
`gesamt: 0 Befund(e)`.

1. **CLEAN:** unveränderter Baum → `make test` grün (alle Pakete `ok`),
   `make arch-check` 0 Befunde.
2. **Symlink-Durchgang:** committeter Verzeichnis-Symlink `docs-link → docs`,
   committetes `docs/page.md`, Link `[p](docs-link/page.md)` → Voll-Lauf
   liefert **zwei** Befunde auf derselben Referenz (`symlink` + `target-untracked`);
   frischer `git clone` des Fixtures: Symlink + Datei werden hergestellt, die
   Referenz löst auf → der `target-untracked`-Befund ist falsch (MEDIUM-2).
3. **Untracktes Symlink-Ziel:** `link.md → real.md` untracked, Link auf
   `link.md`, fokussierter Lauf (`--enable tracked --disable links --disable
   anchors` = `doc-tracked`-Gestalt) → **Exit 0, 0 Befunde** — auf einem
   frischen Klon fehlt `link.md` (MEDIUM-3).
4. **intent-to-add:** `git add -N u.md`, Link auf `u.md` → Exit 0 (grün);
   anschließendes `git commit` (ohne weiteres add) nimmt `u.md` **nicht** mit
   (`git ls-tree HEAD` ohne `u.md`) — frischer Klon wäre rot (INFO-1).
5. **Prozent-Dekodierungs-Parität:** `[a](my%20file.md)`/`[b](%C3%BCber.md)`
   auf getrackte Dateien → grün; `my file.md` entrackt → Befund. Parität
   links↔tracked↔Index bestätigt; die target-Spalte zeigt dabei
   `my%20file.md` (roh — MEDIUM-1).
6. **Roh- vs. aufgelöstes Ziel:** `[u](./docs/u.md#top)` auf untracktes Ziel →
   Befund-target `./docs/u.md#top` statt `docs/u.md` (MEDIUM-1).
7. **Ventil-Glob-Semantik:** `exempt-targets: ["u.*"]` matcht `docs/u.md`
   **nicht** (segmentweise, wie dokumentiert — Negativbefund);
   `["docs/[a-u].md"]` nimmt aus (grün); `["[a/b].md"]` **passiert die
   config-zeitige Validierung** (kein Exit 2) und ist zur Laufzeit still
   wirkungslos — Befund bleibt (LOW-1).
8. **Stray-Range:** `--enable tracked --staged` bzw. `--range HEAD~1..HEAD`
   ohne vcs/commits → still ignoriert, Exit 0; Gegenprobe links-only +
   `--range` → identisch still ignoriert (Vorbestandskontrakt, kein
   Regressions-Delta — Negativbefund).
9. **Linked worktree:** `git worktree add`, Lauf im Worktree-Mount →
   Exit 2 `kein lesbares git-Repository … repository does not exist`
   (gitdir-Redirect zeigt auf nicht gemounteten Host-Pfad) — laut,
   kein Silent-Grün (INFO-2).
10. **Gelöschter Index (Repo MIT Commits):** `rm .git/index` → leere
    Index-Menge, **alle** existierenden Ziele werden gemeldet (Exit 1) —
    Fehlrichtung laut/rot, kein Silent-Grün (Negativbefund zur
    adversarialen Frage „leerer Index in einem Repo mit Commits").
11. **Bild-Ziel:** `![i](img.png)` untracked → Befund (Verhalten korrekt);
    `--doctor` zeigt `target-untracked [tracked]` ohne Klartext (LOW-2).
12. **Mutations-Probe Bild-Einschluss:** `if ref.IsImage { continue }` in
    `CheckTracked` injiziert → **`make test` komplett grün** (`go test ./...`
    alle Pakete `ok`) — der Bild-Einschluss (Spec Schritt 3) ist durch keinen
    Test verriegelt (MEDIUM-4). Revert, Baseline wieder grün.
13. **Mutations-Nachweis nil-Port-Guard:** Guard in `run.go` entfernt →
    genau `TestRunWithVCS_TrackedFailClosedOhnePort` rot (nil-Panik), alle
    übrigen Pakete grün — Commit-Behauptung unabhängig reproduziert.
14. **Mutations-Nachweis Glob-Validierung:** Validierungs-Schleife in
    `applyTracked` entfernt → genau `TestTracked_UngueltigesGlobExit2` rot —
    Commit-Behauptung unabhängig reproduziert (das Test-Fixture enthält
    keine Links, nur der Guard kann Exit 2 liefern — slice-057-R3-Lehre
    eingehalten).

---

## Findings

### MEDIUM-1 — Befund-`target` trägt das rohe Linkziel statt des aufgelösten Zielpfads (Spec Schritt 5); die target-Spalte ist damit nicht ventil-fähig

- **kategorie:** MEDIUM
- **quelle:** `spec/spezifikation.md` §DC-FA-TRK-001.a Schritt 5 („`target` =
  der **aufgelöste** Zielpfad"); Lastenheft-AK Negative („Datei, Zeile,
  aufgelöstes Ziel")
- **pfad:** `internal/hexagon/core/rules/tracked.go:40`
- **befund:** Der Befund setzt `Target: ref.Target` (rohes Ziel wie
  geschrieben). Probe 6: `[u](./docs/u.md#top)` meldet `./docs/u.md#top`
  statt `docs/u.md`; Probe 5: `my%20file.md` statt `my file.md`. Das Ventil
  `exempt-targets` matcht dagegen den **aufgelösten** Pfad — wer den
  angezeigten target-Wert in die Config kopiert, erzeugt ein Glob, das nie
  matcht (`./docs/u.md#top` ↛ `docs/u.md`). Der Unit-Test verriegelt nur den
  Fall roh==aufgelöst (`got.Target != "u.md"` bei Link-Text `u.md`), die
  Spec-Abweichung bleibt testfrei. Bei `links` ist die Roh-Anzeige
  spec-konform (§2: „geprüftes Ziel"), für `tracked` fordert Schritt 5
  explizit den aufgelösten Pfad — doc-first: die Spec ist die höherrangige
  Quelle.
- **verifizierbar:** ja — Proben 5/6 (Fixture-Repo, ein Link mit
  `./`-Präfix/Fragment/Prozent-Encoding auf untracktes Ziel).

### MEDIUM-2 — Pfad durch getrackten Verzeichnis-Symlink ⇒ falscher `target-untracked` zusätzlich zum `links`-`symlink`-Befund (Doppelbefund entgegen dem deklarierten Prinzip)

- **kategorie:** MEDIUM
- **quelle:** ADR-0030 (§Entscheidung „Kein Doppelbefund … der strukturelle
  Befund bleibt beim Struktur-Modul"); `DC-FA-TRK-001` (Befund-Semantik
  „wäre auf jedem frischen Klon ein `target-missing`"); Slice-Plan §4
  (Symlink-Ränder)
- **pfad:** `internal/hexagon/core/rules/tracked.go:28-29`
- **befund:** `fsys.Kind(rel)` lstatet nur die **letzte** Pfad-Komponente;
  Zwischen-Symlinks werden vom OS aufgelöst. Probe 2: Link
  `docs-link/page.md` über den committeten Symlink `docs-link → docs` auf
  das committete `docs/page.md` ⇒ `Kind`=Datei, aber der Index führt nur
  `docs-link` und `docs/page.md`, nicht den durchgereichten Pfad ⇒
  `target-untracked` — **falsch**: der frische Klon stellt Symlink und Datei
  her, die Referenz löst auf (per Clone-Gegenprobe belegt); die
  Befund-Message behauptet das Gegenteil. Zugleich meldet `links` dieselbe
  Referenz bereits als `symlink` (`symlinkInPath` prüft jede Komponente) —
  zwei Befunde für eine Referenz, deren struktureller Grund beim
  Struktur-Modul liegt.
- **verifizierbar:** ja — Probe 2 inkl. Clone-Gegenprobe; nach einer
  Klärung muss dasselbe Fixture ohne `tracked`-Befund laufen.

### MEDIUM-3 — Existierendes, untracktes Symlink-Ziel wird still übersprungen: die beworbene Frischer-Klon-Drift entgeht genau dem fokussierten `doc-tracked`-Gate

- **kategorie:** MEDIUM (Stilles-Grün-Beobachtung im verteilten Gate-Pfad;
  Begründung gegen HIGH unten)
- **quelle:** `DC-FA-TRK-001` (Zweck: „Existiert das Ziel nur im
  Arbeitsbaum … `target-untracked`"); §DC-FA-TRK-001.a Schritt 3
  (Kandidaten-Definition entscheidet den Symlink-Fall nicht);
  `internal/adapter/driving/cli/print_mk.go:603-605` (Verteilungs-Vehikel)
- **pfad:** `internal/hexagon/core/rules/tracked.go:29`
- **befund:** `kind != driven.KindFile` überspringt Symlink-Ziele
  (`KindSymlink`, Lstat). Probe 3: untracktes `link.md → real.md`, Link auf
  `link.md`, Lauf in `doc-tracked`-Gestalt (links/anchors disabled) ⇒
  Exit 0 — auf jedem frischen Klon fehlt `link.md`, exakt die Drift-Klasse,
  für die das Modul existiert. Kein Test dokumentiert die Skip-Entscheidung;
  Spec und ADR schweigen zu Symlink-Zielen (Slice-Plan §4 nannte die Ränder
  ausdrücklich). **Nicht HIGH**, weil die Repo-Doktrin (`DC-FA-LINK-002`)
  Symlink-Referenzen ohnehin verbietet — der Voll-Lauf (`doc-check`) meldet
  die Referenz als `symlink`; still grün bleibt nur ein Konsument, der das
  fokussierte Gate ohne die Default-Module fährt und bereits links-rot wäre.
  Zusammen mit MEDIUM-2 ist die Symlink-Semantik in **beide** Richtungen
  inkonsistent: der reale Drift-Fall (untrackter Symlink) ist still grün,
  der Nicht-Drift-Fall (Pfad durch getrackten Symlink) ist rot.
- **verifizierbar:** ja — Probe 3; ein AK-artiger Test (untracktes
  Symlink-Ziel im fokussierten Lauf) würde heute grün statt rot laufen.

### MEDIUM-4 — Bild-Einschluss (Spec Schritt 3) und Normalisierungs-Parität sind unverriegelt: die Mutation „Bilder überspringen" überlebt die gesamte Test-Suite

- **kategorie:** MEDIUM
- **quelle:** `spec/spezifikation.md` §DC-FA-TRK-001.a Schritt 3
  („Bild-Ziele eingeschlossen"); Slice-Plan §4 („Normalisierungs-Parität
  testen"); Reviewer-Skill-Anker „fehlende Negativtests bei neuem
  öffentlichen Vertrag"; slice-057-R3-Lehre (Test verriegelt genau seinen
  Zweig)
- **pfad:** `internal/hexagon/core/rules/tracked_test.go:28-35` (Fixture
  ohne Bild-, Symlink- und Encoding-Fall; dito `cli_tracked_test.go`,
  `git_tracked_test.go`)
- **befund:** Probe 12: `if ref.IsImage { continue }` in `CheckTracked`
  injiziert ⇒ `make test` bleibt **komplett grün** — ein Refactoring, das
  Bild-Referenzen aus der Prüfung filtert, würde von keinem Test gefangen,
  und ein untracktes Bild wäre im `doc-tracked`-Gate still grün. Dieselbe
  Lücke gilt für die vom Slice-Plan §4 geforderte Pfad-Normalisierungs-
  Parität (Prozent-Dekodierung, `./`-Präfix, Fragment-Abtrennung): heute
  korrekt (Proben 5/6/11 grün bzw. rot wie erwartet), aber ausschließlich
  durch geteilte Helfer garantiert, nicht durch einen Test des neuen
  Vertrags.
- **verifizierbar:** ja — Probe 12 (Mutation + `make test`); nach
  Test-Nachrüstung muss dieselbe Mutation rot laufen.

### LOW-1 — Ventil-Validierung prüft das ungesplittete Muster, die Laufzeit matcht segmentweise: ein Glob mit `/` in einer Zeichenklasse passiert die Validierung und ist still wirkungslos

- **kategorie:** LOW (Fehlrichtung laut — Befund bleibt trotz Ventil;
  Muster exotisch)
- **quelle:** `DC-FA-CONF-001` („vollständig validiert … keine Prüfung mit
  stillschweigenden Defaults"); Kommentar-Anspruch in `applyTracked`
  („sonst … das Ventil wäre still wirkungslos")
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:384`
- **befund:** `path.Match(g, "probe")` validiert das **ganze** Muster;
  `matchGlob` splittet zur Laufzeit an `/` und matcht segmentweise
  (`ok, _ := path.Match` schluckt `ErrBadPattern`). Probe 7c:
  `exempt-targets: ["[a/b].md"]` ⇒ kein Exit 2, zur Laufzeit matcht das
  Segment-Fragment `"[a"` nie — der Befund bleibt trotz deklariertem Ventil.
  Genau die Falle, die der Kommentar ausschließen will, nur eine
  Dialekt-Ebene tiefer. Zum Vergleich: `planning.slice-glob` validiert und
  matcht **beide** Male ganzheitlich (`rules/planning.go:106`) — dort ist
  die Probe-Validierung exakt; `tracked` ist das erste Modul, dessen
  Validierungs-Dialekt vom Laufzeit-Dialekt abweicht.
- **verifizierbar:** ja — Probe 7c (Config-Fixture, kein Code-Edit).

### LOW-2 — `target-untracked` fehlt in der „kanonischen Liste aller Grund-Codes" (`AllReasons`/`reasonTexts`): `--doctor` zeigt den nackten Code; die Lücken-Klasse wiederholt sich seit mehreren Slices

- **kategorie:** LOW
- **quelle:** Maintainability; Kommentar-Vertrag in
  `internal/hexagon/core/app/diagnose.go` („Ein neuer Grund-Code wird hier
  ergänzt; fehlt der Klartext, bricht der Test")
- **pfad:** `internal/hexagon/core/app/diagnose.go:67`
- **befund:** `AllReasons()`/`reasonTexts()` enden bei `hostpath-forbidden`;
  `target-untracked` (und Vorbestand: die Grund-Codes von
  diagrams/versions/pins/immutable/vcs/commits/planning) fehlen. Der
  Vollständigkeits-Test prüft nur „jeder gelistete Code hat Klartext" —
  er kann das Fehlen neuer Codes nicht fangen; `--doctor` fällt auf den
  nackten Code zurück (Probe 11: `target-untracked [tracked]` ohne
  Klartext-Zeile, während Alt-Codes Klartext tragen). Kein
  Korrektheits-Bruch (fail-safe-Fallback), aber der Kommentar-Vertrag ist
  seit mindestens sechs Modulen tote Doku — dritte+ Wiederholung derselben
  Klasse: Steering-Loop-Signal laut Skill (Sensor/Kommentar nachziehen
  statt je Slice neu vergessen).
- **verifizierbar:** ja — Probe 11 (`--doctor` auf Fixture);
  `grep -c Reason internal/hexagon/core/app/diagnose.go` vs.
  `model`/`rules`-Konstanten.

### INFO-1 — intent-to-add-Einträge (`git add -N`) gelten als getrackt, werden von `git commit` aber nicht mitgenommen: die Frischer-Klon-Garantie hat eine undokumentierte Index-Ausnahme

- **kategorie:** INFO (dokumentationswürdige, undokumentierte Annahme)
- **quelle:** `DC-FA-TRK-001` („eine … gestagte, noch nie committete Datei
  gilt als getrackt"); ADR-0030 (Index = eine Wahrheit)
- **pfad:** `internal/adapter/driven/git/git.go:240-250`
- **befund:** `TrackedPaths` nimmt jeden Index-Eintrag (Flags ignoriert).
  Probe 4: `git add -N u.md` ⇒ Modul grün, aber plain `git commit` schließt
  `u.md` nicht ein — der frische Klon ist rot, obwohl das Modul grün war.
  „Index-Wahrheit" deckt den Fall formal (der Eintrag existiert), die
  beworbene Garantie („beim Erzeuger grün ⇒ Klon grün") nicht. Randfall
  ohne Handlungszwang; als Annahme notierenswert (Spec-Out-of-Scope-Liste
  oder Handbuch).
- **verifizierbar:** ja — Probe 4.

### INFO-2 — Linked worktrees: `.git`-Datei mit gitdir-Redirect ⇒ Exit 2 mit go-git-Roh-Meldung; das Handbuch nennt nur „kein lesbares `.git`"

- **kategorie:** INFO
- **quelle:** Maintainability; `DC-FA-TRK-001` fail-closed-AK (Verhalten
  korrekt); Slice-Plan §4 (Konsumenten-Ränder)
- **pfad:** `docs/user/benutzerhandbuch.md:813-819`;
  `internal/adapter/driven/git/git.go:31-37`
- **befund:** Probe 9: im per `git worktree add` erzeugten Checkout zeigt
  die `.git`-**Datei** auf einen Host-Pfad außerhalb des `:ro`-Mounts ⇒
  `kein lesbares git-Repository unter /repo: repository does not exist`,
  Exit 2 — fail-closed und laut (korrekt), aber für Konsumenten mit
  Worktree-basierten CI-Checkouts wenig diagnostisch; Spec listet nur
  „verschachtelte Arbeitsbäume" als Out-of-Scope, das Handbuch nennt den
  Randfall nicht.
- **verifizierbar:** ja — Probe 9.

---

## Negativbefunde (geprüft, ohne Befund)

- **Adapter-/Index-Semantik (`TrackedPaths`):** Index-Wahrheit korrekt —
  committete UND frisch gestagte Dateien enthalten, reine
  Arbeitsbaum-Dateien nicht (Adapter-Fixture-Test + Proben 5/10);
  leere Menge nur bei leerem/fehlendem Index-File, und diese Richtung
  **über-feuert laut** statt still grün (Probe 10: gelöschter Index ⇒ alle
  Ziele gemeldet, Exit 1 — deckungsgleich mit `git ls-files` auf demselben
  Zustand); unlesbarer Index ⇒ Fehler ⇒ Exit 2. Kein Silent-Grün-Pfad im
  Adapter gefunden.
- **Fail-closed-Kette Kern+CLI:** beide Guards (nil-Port in `run.go`,
  Glob-Validierung in `applyTracked`) von mir **unabhängig
  mutations-reproduziert** — je Mutation genau der benannte Test rot
  (Proben 13/14), Fixtures so gebaut, dass nur der Guard den Ausgang trägt
  (slice-057-R3-Lehre eingehalten); zusätzlich CLI-seitig `Open`-Fehler ⇒
  Exit 2 mit stderr-Hinweis (E2E-Test + Probe 9). Redundante Absicherung
  Kern/CLI unabhängig voneinander wirksam.
- **`resolveVCS`-Ränder:** `tracked` ohne Range gültig, `vcs`/`commits`
  behalten die Range-Pflicht auch in Kombination (E2E-Test verriegelt beide
  Richtungen); stray `--range`/`--staged` bei tracked-only wird still
  ignoriert — identisch zum Vorbestand bei links-only (Probe 8), kein neues
  Delta; `--commit-msg` bleibt eigener Kurzschluss-Modus (überspringt ALLE
  Scan-Module, nicht tracked-spezifisch; `gitComboError` unverändert);
  `--trace` läuft vor der Modul-Auflösung und nutzt keine Module — von
  `tracked` unberührt.
- **Kein-Doppelbefund (missing/dir/escape):** nicht existierende Ziele,
  Verzeichnis-Ziele und Repo-Escapes sind keine Kandidaten (Unit- +
  Integrations-Test; `links` meldet `target-missing`/`repo-escape` allein) —
  die Ausnahme ist die Symlink-Klasse (MEDIUM-2/-3).
- **Ventil-Glob-Doku-Parität:** segmentweise Semantik „wie `scan.ignore`"
  stimmt mit dem dokumentierten Dialekt überein — `u.*` matcht `docs/u.md`
  **nicht** (Probe 7a), `docs/[a-u].md` und `build/**` matchen wie erwartet;
  leeres Glob ⇒ Exit 2 (Test verriegelt).
- **Default-aus byte-identisch/opt-in-Kette:** Modul-aus-Test + nil-Port
  strukturell ungenutzt (Probe 13 belegt: ohne aktives `tracked` keine
  Port-Berührung, sonst wäre die Panik-Mutation breiter rot);
  `validModules`/`EffectiveModules`/`applyScopes` konsistent erweitert;
  `.d-check.yml`-`modules:`-Basis unverändert.
- **`--print-mk doc-tracked`:** ohne Range, fokussierte `--disable`-Liste aus
  `ValidModules` abgeleitet (kein `--disable tracked`), FÜNF-Verben-Zählung
  des Templates stimmt mit dem Kommentar überein; Test verriegelt Target,
  Recipe-Flags und Range-Freiheit.
- **`FOCUS_DISABLE` (eigenes Makefile):** unverändert korrekt — die Variable
  spiegelt die `.d-check.yml`-`modules:`-Liste, `tracked` ist kein
  Default-Modul (ADR-0030-Konsequenz geprüft, kein Nachzieh-Bedarf).
- **Architektur/Hard Rules:** `rules/tracked.go` importiert nur
  `core/model` + `port/driven` (ADR-0005-konform), go-git bleibt allein im
  git-Adapter — eigener `make arch-check`-Lauf: 0 Befunde; keine
  Inline-Suppressions im Diff.
- **DC-QA-03:** keine neue Netz-Tür; sämtliche Proben liefen mit
  `--network none` + `:ro`-Mount grün/rot wie erwartet (Index-Lesen im
  read-only-Mount funktioniert).
- **Referenz-Richtung (SDP)/Marker-Ehrlichkeit:** keine Provenance-Marker im
  Diff; ADR-/Slice-Nennungen in Code-Kommentaren sind Provenance. Ohne
  Befund.

---

## Kategorie-Summary

| Kategorie | Anzahl |
| --- | --- |
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 2 |
| INFO | 2 |

## Verdikt

**NACHBESSERN (vor Closure).** Kern-Mechanik, fail-closed-Kette,
Range-Freiheit und Config-Surface sind sauber gebaut; beide
Guard-Mutations-Behauptungen des Commits habe ich unabhängig reproduziert,
und die adversarialen Silent-Grün-Kandidaten (leerer Index, gelöschter
Index, Worktree, stray Range) lösen alle laut aus. Kein HIGH. Die vier
MEDIUMs teilen sich in zwei direkt behebbare Code-/Test-Punkte (MEDIUM-1:
rohes statt aufgelöstes `target` entgegen Spec Schritt 5; MEDIUM-4:
Bild-Einschluss und Normalisierungs-Parität unverriegelt — Mutation
überlebt die Suite) und eine Design-Klärung (MEDIUM-2/-3: die
Symlink-Semantik ist unspezifiziert und in beide Richtungen inkonsistent —
falscher Befund beim Durchgangs-Symlink, stilles Grün beim untrackten
Symlink-Ziel; gehört als Entscheidung in Spec/ADR, nicht nur in den Code).
LOW-1/LOW-2 sind nicht blockierend; LOW-2 ist als wiederholte Klasse ein
Steering-Kandidat, INFO-1/-2 sind Doku-Kandidaten ohne Handlungszwang.
