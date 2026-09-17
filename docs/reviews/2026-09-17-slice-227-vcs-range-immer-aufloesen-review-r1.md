# Review-Report: slice-227 — 2026-09-17 (R1)

**Review-Art:** Code-Review — geprüft wird die Implementierung gegen den
Slice-Plan, `DC-FA-VCS-001`/`DC-FA-COMMITS-001` und `AGENTS.md` §3.

**Gegenstand:** Commit-Range `a4e3bb96^..1a380bb4` (drei Commits:
`a4e3bb96` Beanspruchung, `7aa51f74` Fix + Tests, `1a380bb4`
DoD-Häkchen).

**Skill:** `.harness/skills/reviewer.md` @ v1.16.0 / 2026-09-07 ·
**Modell:** `claude-sonnet-5` · **Datum:** 2026-09-17

**Eingangs-Kontext:**

- Slice-Plan `docs/plan/planning/in-progress/slice-227-vcs-range-immer-aufloesen.md`
- [GitHub Issue #4](https://github.com/pt9912/d-check/issues/4) Punkt 1 (Anlass)
- `DC-FA-VCS-001`/`DC-FA-VCS-001.a`, `DC-FA-COMMITS-001`/`DC-FA-COMMITS-001.a`
  (`spec/lastenheft.md`, `spec/spezifikation.md`)
- `AGENTS.md` §3 (Hard Rules), §6 (Workflow)
- Vorherige Findings am selben Modul: `docs/reviews/2026-09-17-slice-226-pack-alias-fremdpraefix-review-r1.md`,
  `docs/reviews/2026-09-17-slice-220-vcs-pfadmenge-review-r1.md` — keine
  davon betrifft die hier geänderte Aufrufreihenfolge oder Kommentarklasse.

---

## Verifikationen dieses Laufs

Alle folgenden Läufe wurden **selbst ausgeführt** (Docker/`make`-only, kein
Host-Go; `git` als normales Host-Tool nur für die Repro-Repos), nicht aus
der Commit-Botschaft übernommen:

- **Bewusstes Brechen (Modul 11) für DoD (1)/(2):** `vcs.go`/`commits.go`
  im Arbeitsbaum auf den Elternstand (`7aa51f74^`) zurückgesetzt, die
  **neuen** Testdateien unverändert gelassen, `make test` gefahren.
  Ergebnis: `TestCheckCommitsInert` und `TestVCSInert` schlagen **beide**
  mit exakt der erwarteten Meldung fehl (`"unauflösbare Range mit leerer
  Klassen-Config still passiert — erwartet war ein Fehler"`), alle
  übrigen Pakete bleiben grün. Die beiden neuen Testfälle diskriminieren
  also echt zwischen Vor- und Nachzustand, keiner ist tautologisch grün.
  Danach `git checkout --` — Arbeitsbaum wieder sauber.
- **`make gates`** (voller Lauf, kein Cache-Vorbehalt genannt): grün —
  `baseline-verify + workflow-pins + doc-check + lint + test + arch-check
  + coverage-gate + semgrep + gate-consistency + planning-check green`.
  `coverage-gate`: **94,70 %** (Schwelle 93 %, keine Regression ggü.
  vorherigen Läufen).
- **Eigene End-to-End-Probe** gegen das aus `HEAD` gebaute
  `d-check:latest` (Scratch-Git-Repo, zwei Commits, **kein** `.d-check.yml`):
  - `--enable vcs --range deadbeef..cafebabe` → `d-check: error: Range-Basis
    "deadbeef" nicht auflösbar: reference not found`, **Exit 2**.
  - `--enable commits --range deadbeef..cafebabe` → dieselbe Meldung,
    **Exit 2**.
  - `--enable vcs --range <auflösbar>` bzw. `--enable commits --range
    <auflösbar>`, weiterhin ohne jeden `vcs:`/`commits:`-Block: beide
    `d-check: 1 Datei(en) geprüft, 0 Befund(e)`, **Exit 0** — die
    dokumentierte Opt-in-Trägheit bleibt für den Happy Path erhalten.
  - Zusätzlich (nicht in DoD/Plan benannt, siehe R1-F-2): `--enable
    commits --staged` ohne `commits:`-Block gegen dasselbe Repo →
    `d-check: error: commits: --staged wird nicht unterstützt — nutze
    --range oder --commit-msg …`, **Exit 2**. Dieselbe Kombination gegen
    ein aus dem Elternstand gebautes Image (`commits.go` vor dem Fix)
    ergab **Exit 0, 0 Befunde** — die Kombination hat also tatsächlich
    ihr Verhalten gewechselt.
- **Statischer Check** auf verbleibende Kommentar-Verstöße
  (`AGENTS.md` §3.7): `grep -rn "slice-227" internal/*.go` — keine
  Treffer in Produktionscode; `grep -rn "Issue #4" internal/*.go` — zwei
  Treffer, beide in `_test.go`-Dateien (siehe R1-F-1).
- **Hexagon-Import-Grenze** (`ADR-0005`): Diff fügt keine neuen Imports
  hinzu, beide Funktionen bleiben innerhalb `internal/hexagon/core/rules`.
  Kein Befund.
- **Edge Case „vcs/commits aktiv, aber base/head leer"** (Aufgabenpunkt 7,
  erster Spiegelstrich): `internal/adapter/driving/cli/cli.go` `vcsRefs`
  (Zeilen 685–699) liefert bei fehlendem `--range`/`--staged` eine
  `msg != ""`, und `resolveVCS` (Zeilen 705–733) gibt dann **vor** dem
  Öffnen des Adapters `nil, "", "", 2` zurück — der Adapter (und damit ein
  von `nil` verschiedener `vcs`-Port) existiert in diesem Fall nie
  gleichzeitig mit leeren `base`/`head`. Die Reihenfolge-Änderung in
  `CheckVCS`/`CheckCommits` kann diesen Pfad also nicht erreichen. Kein
  Befund.

## Findings

### R1-F-1 (HIGH) — Neue Testkommentare tragen Herkunfts-Prosa statt einer der fünf zulässigen Klassen

- **Kategorie:** HIGH
- **Quelle:** `AGENTS.md` §3.7 (Kommentare tragen eine der fünf Klassen —
  Zusage · Kopplung · Abgrenzung · Rang-Zeiger · Grenze; Herkunft nur als
  **ein** auflösbares Feld nach dem Baseline-Schema `DC-*`/`ADR-*`/`MR-*`/
  `seit welle-<NN>`)
- **Pfad:** `internal/hexagon/core/rules/commits_test.go:95-96`,
  `internal/hexagon/core/rules/vcs_test.go:233-234`
- **Befund:** Beide Doc-Kommentare nennen in Klammern
  `„(GitHub Issue #4 Punkt 1)"` als Begründung, warum der Test existiert —
  eine externe Vorgangs-Referenz außerhalb des im Regelwerk zugelassenen
  Anker-Schemas. Das ist Herkunfts-Prosa (Chronik, wie der Test entstand),
  keine der fünf zulässigen Klassen; die Bestandsgrenze in §3.7 gilt
  ausdrücklich nicht rückwirkend für **Neuzugänge** — beide Zeilen wurden
  mit `7aa51f74` neu geschrieben, nicht geerbt. Dieselbe Repo-Konvention
  wurde im aktuellen Bestand bereits einmal für Slice-/Befund-Bezüge in
  Kommentaren gefordert (Nutzer-Feedback zu Code-/Test-Kommentaren: solche
  Bezüge gehören in Commit-Message/Review-Report, nicht in den Kommentar);
  eine GitHub-Issue-Nummer ist dieselbe Klasse externer Prosa wie eine
  Slice- oder Befund-Nummer.
- **Verifizierbar:** ja — `grep -n "Issue #4" internal/hexagon/core/rules/*_test.go`.
- **Klasse:** `kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen`

### R1-F-2 (MEDIUM) — `--enable commits --staged` mit leerer Klassen-Config wechselt undokumentiert und ungetestet das Verhalten

- **Kategorie:** MEDIUM
- **Quelle:** `DC-FA-COMMITS-001` (Lastenheft, Akzeptanzkriterium
  *„fail-closed (git-/Range-Eingabe fehlt)"*: eine fehlende/unauflösbare
  Range bricht mit Exit 2 ab — unbedingt, keine Ausnahme für leere
  `commits.id-patterns` genannt); Prüffrage 13 des Reviewer-Skills
  (fehlender Negativtest zu einem bestehenden öffentlichen Vertrag, der
  sich hier ändert)
- **Pfad:** `internal/hexagon/core/rules/commits.go:25-38`
- **Befund:** `CommitMessages` wird jetzt **unconditional** aufgerufen,
  sobald `vcs != nil` — auch im `--staged`-Modus (`head ==
  driven.IndexRef`), den der Port laut eigenem Vertrag
  (`internal/hexagon/port/driven/vcs.go:41`) **nie** auflöst
  („head == IndexRef wird nicht unterstützt … ⇒ Fehler"). Vor diesem Fix
  überlebte `--enable commits --staged` mit leerem `commits:`-Block genau
  deshalb unbemerkt: der alte Frühausstieg (`len(cfg.IDPatterns) == 0`)
  griff, bevor `CommitMessages` je aufgerufen wurde, und lieferte still
  `0 Befunde`/Exit 0 — dieselbe Fehlerklasse wie Issue #4, nur über
  `--staged` statt `--range` ausgelöst, und **nicht** Teil der
  Root-Cause-Analyse des Slice-Plans. Empirisch bestätigt (Alt-Image
  gegen den Elternstand `7aa51f74^` gebaut, s. o.): Alt → Exit 0/0
  Befunde, Neu → Exit 2 (Fehlermeldung „--staged wird nicht
  unterstützt"). Der Nettoeffekt ist in sich stimmig mit dem
  Lastenheft-Vertrag (jetzt korrekt fail-closed) — aber weder DoD, Plan
  noch die beiden neuen Testfälle erwähnen diese Kombination, und keiner
  der beiden neuen Tests deckt sie ab. Ein künftiger Umbau, der die
  Reihenfolge in `CheckCommits` erneut vertauscht, würde genau diesen
  Pfad lautlos zurückdrehen, ohne dass ein Test es fängt.
- **Verifizierbar:** ja — `docker run … d-check:latest --enable commits
  --staged` (kein `commits:`-Block) gegen ein Scratch-Repo vor und nach
  `7aa51f74`.
- **Klasse:** `unbelegte-verhaltensaenderung-staged-commits`

## Negativbefunde

- **Root-Cause-Reihenfolge:** in beiden Funktionen geprüft — `vcs == nil`
  zuerst (unverändert), dann `AllPaths`/`CommitMessages` unconditional,
  erst danach die Klassen-Config-Prüfung. Kein Befund.
- **Bewusstes Brechen (DoD 1/2):** beide neuen Testfälle schlagen gegen
  den Elternstand aus dem behaupteten Grund fehl (s. o.). Kein Befund.
- **DoD (3) — Rückwärtskompatibilität des Opt-in-Vertrags über `--range`:**
  empirisch gegen das gebaute Image bestätigt für **beide** Module. Kein
  Befund.
- **Edge Case „aktives Modul, aber leere base/head"** (Aufgabenpunkt 7):
  strukturell unerreichbar über die CLI (`vcsRefs`/`resolveVCS` lehnen
  vorher mit Exit 2 ab). Kein Befund.
- **Hexagon-Import-Grenze (ADR-0005):** keine neuen Imports, keine
  Schichtverletzung. Kein Befund.
- **Inline-Suppression / Gate-Lockerung ohne ADR** (`AGENTS.md`
  §3.2/§3.6): geprüft — kein `//nolint`, keine Schwellen-Senkung. Kein
  Befund.
- **Netzzugriff außerhalb `external`** (`DC-QA-03`): geprüft — reine
  Lese-Aufrufe über den vorhandenen VCS-Port, kein neuer Netzpfad. Kein
  Befund.
- **Kommentar-Klassen der Produktionscode-Kommentare** (`vcs.go`/
  `commits.go`, außerhalb der unter R1-F-1 genannten Testdateien): die
  vier geänderten/neuen Kommentare (`CheckVCS`/`CheckCommits`-Docstring,
  die beiden Inline-Kommentare vor dem `AllPaths`/`CommitMessages`-Aufruf)
  beschreiben ausschließlich die aktuelle Aufrufreihenfolge und ihre
  Begründung (Zusage-Klasse), tragen keine Slice-ID, kein Review-Vokabular
  und keine Herkunfts-Prosa. Kein Befund.
- **Spec-Treue der Abgrenzung „kein Spec-Nachtrag"** (§1 des Slice-Plans):
  geprüft gegen `spec/lastenheft.md` — beide betroffenen
  Akzeptanzkriterien (*„fail-closed (git-/Range-Eingabe fehlt)"*) sagen
  das jetzige Verhalten bereits **unbedingt** zu, ohne eine Ausnahme für
  leere Klassen-Config zu benennen; der Alt-Zustand war also bereits vor
  diesem Slice ein Spec-Verstoß, kein neues Verhalten. Auch
  `docs/user/benutzerhandbuch.md` und `CHANGELOG.md` wurden nach einer
  Beschreibung des alten (fehlerhaften) Verhaltens durchsucht — keine
  Fundstelle verspricht die alte, stille Passierbarkeit als Vertrag. Kein
  Befund.
- **Abgrenzungspunkt 2 des Slice-Plans** („kein pauschales Audit anderer
  opt-in-Module — nur `vcs`/`commits` tragen eine zweite, unabhängig
  gültige Range-Eingabe"): stichprobenartig gegen `tracked` geprüft
  (`internal/hexagon/core/rules/run.go:71-79`) — `tracked` ruft
  `vcs.TrackedPaths()` bereits **vor** jeder Config-Prüfung unconditional
  auf und trägt denselben Fehler von vornherein nicht. Die Abgrenzung
  hält für die geprüfte Stichprobe. Kein Befund.
- **`make gates` grün, Coverage ohne Regression:** s. Verifikationen
  oben. Kein Befund.
- **Commit-Botschaft überdehnt die Messung nicht** (§5 `AGENTS.md`,
  Reviewer-Anker MEDIUM): „beide neuen Testfälle schlagen gegen den
  Vorzustand aus dem erwarteten Grund fehl" und „End-to-End … reproduziert"
  sind beide durch diesen Review-Lauf unabhängig bestätigt, keine der
  beiden Aussagen reicht weiter als das tatsächlich Gemessene. Kein
  Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:**
`kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen` ·
`unbelegte-verhaltensaenderung-staged-commits`

## Verdikt

**Blockiert.** R1-F-1 (HIGH) ist ein unmittelbarer `AGENTS.md`
§3.7-Verstoß in **neu geschriebenem** Code (keine Bestandsgrenze greift)
und mechanisch trivial zu beheben (die Klammer-Referenz entfernen oder
durch eine reine Verhaltens-Beschreibung ersetzen; der Issue-Bezug gehört
in Commit-Message/Slice-Plan, wo er bereits steht). R1-F-2 (MEDIUM) hält
die Closure nicht zwingend auf, sollte aber vor `done/` entweder durch
einen dritten Regressionstest (`--staged` + leere `commits`-Config)
geschlossen oder als bewusst akzeptierte Fußnote in DoD/Closure-Notiz
festgehalten werden — sonst verschwindet die einzige Stelle, an der dieses
Verhalten je beobachtet wurde, mit diesem Review-Report.
