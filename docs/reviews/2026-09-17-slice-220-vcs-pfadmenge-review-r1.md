# Review-Report: slice-220 — 2026-09-17 (R1)

**Review-Art:** Code-Review — geprüft wird die Implementierung gegen den
Slice-Plan, `CO-001`, `DC-FA-VCS-001`/`DC-FA-VCS-001.a` und `AGENTS.md` §3.

**Gegenstand:** Commit `3100e2f8` (feat(vcs): geschützte Pfad-Menge direkt
gegen beide Trees auflösen, statt dem Diff zu vertrauen).

**Skill:** `.harness/skills/reviewer.md` @ v1.16.0 / 2026-09-07 ·
**Modell:** `claude-sonnet-5` · **Datum:** 2026-09-17

**Eingangs-Kontext:**

- Slice-Plan `docs/plan/planning/in-progress/slice-220-vcs-pfadmenge-statt-diff.md`
- `CO-001` (`docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md`) — dessen
  dritte und vierte Ausprägung dieser Commit adressiert
- Beobachtungs-Register: `BEO-ALL/fix-schliesst-pfad-nicht-klasse` (Anlass
  des Entwurfswechsels)
- `DC-FA-VCS-001`/`DC-FA-VCS-001.a` (`spec/lastenheft.md`,
  `spec/spezifikation.md`)
- `ADR-0024` (Accepted; Schärft-Ziel `DC-FA-VCS-001.a`)
- `AGENTS.md` §3 (Hard Rules), §6 (Workflow)

---

## Verifikationen dieses Laufs

Alle folgenden Läufe wurden **selbst ausgeführt** (Docker/`make`-only, kein
Host-Go), nicht nur aus der Commit-Botschaft übernommen:

- `make test`: grün, alle Pakete — `internal/adapter/driven/git` und
  `internal/hexagon/core/rules` eingeschlossen.
- `make lint`: `0 issues`.
- `make arch-check`: `gesamt: 0 Befund(e)` (Hexagon-Grenze über a-check).
- `make coverage-gate`: **94.70 %** (Schwelle 93 %) — deckt sich exakt mit
  der Commit-Botschaft; `AllPaths` 90.0 %, `pathsAt` 100.0 %, `walkTree`
  91.7 %, `CheckVCS` 95.8 %, `protectedSet` 100.0 %.
- `make semgrep`: `0 findings`.
- `make gates`: grün, `baseline-verify + workflow-pins + doc-check + lint +
  test + arch-check + coverage-gate + semgrep + gate-consistency +
  planning-check`; `targets`/`planning`-Läufe je **797 Datei(en) geprüft, 0
  Befund(e)** — deckt sich exakt mit der Commit-Botschaft
  („zehn Gates, 797 Dateien, 0 Befunde").
- **Empirische Gegenprobe von CO-001s dritter Ausprägung**, gegen ein selbst
  gebautes Probe-Repo (nicht die mitgelieferte Test-Suite): Verzeichnis `sub/`
  mit `sub/adr-x.md` (`**Status:** Accepted`) committet, dann gelöscht; das
  lose Tree-Objekt zu `sub` am BASE-Commit entfernt. Binary aus `3100e2f8^`
  (Parent-Commit, eigens per `git worktree` + `make build` gebaut) meldet
  `d-check: 0 Datei(en) geprüft, 0 Befund(e)` (Exit 0) — der stille Fehlschlag,
  den CO-001 beschreibt. Binary aus `3100e2f8` (`make build`) meldet gegen
  **dasselbe** Probe-Repo:
  `error: Range-Basis "…" nicht auflösbar: nicht vollständig lesbarer Tree zu
  "…": nicht lesbarer Unterbaum "sub": object not found` (Exit 2).
- **Empirische Gegenprobe von CO-001s vierter Ausprägung** (HEAD-Tree
  unlesbar), ebenfalls gegen ein selbst gebautes Probe-Repo: Parent-Binary
  meldet fälschlich `core-drift-vcs` („immutable Datei gelöscht oder
  umbenannt") mit Exit 1; das neue Binary meldet stattdessen
  `error: Range-Spitze "…" nicht auflösbar: … nicht lesbarer Unterbaum "sub":
  object not found` (Exit 2) — die Fehldiagnose ist behoben.
- **Determinismus (DC-QA-02):** dreimaliger Lauf gegen ein Probe-Repo mit
  fünf gleichzeitig geänderten Klassen-Dateien (gelöscht, modifiziert,
  reiner Rename, neu, neu) lieferte dreimal byte-identische, alphabetisch
  sortierte Ausgabe (`adr-2.md`, `adr-3.md`, `adr-5.md`).

## Findings

### R1-F-1 (MEDIUM) — `DC-FA-VCS-001.a` Schritt 2 beschreibt weiterhin den abgelösten Diff-Mechanismus

- **Kategorie:** MEDIUM
- **Quelle:** `AGENTS.md` §6 Schritt 7 (*„Doku/Indizes aktualisieren, falls
  ein öffentlicher Vertrag berührt"*); Slice-Plan-Kopf selbst, Feld
  **Berührte Spec-Stellen** (*„`DC-FA-VCS-001.a` Schritt 2 — der Slice
  ändert, **woher** sie kommt"*)
- **Pfad:** `spec/spezifikation.md:1762-1799` (§`DC-FA-VCS-001.a` Schritt 2)
- **Befund:** Schritt 2 sagt weiterhin *„Aus dem Diff `BASE..HEAD` (bzw.
  staged) werden die Pfade gewählt … Pro Eintrag zählt der Diff-Status:
  Modifikation/Typänderung (M/T) → Core-Vergleich; Löschung (D) /
  Umbenennung (R) → Pfad-Stabilitäts-Prüfung; Hinzufügung (A) → frei"`,
  gefolgt von zwei Absätzen über *„der Range-Pfad difft ohne
  Rename-Erkennung"* und *„A ist auch die Antwort, die ein Tree-Diff gibt,
  wenn er den BASE-Stand gar nicht lesen konnte"*. Keiner dieser Mechanismen
  existiert nach diesem Commit noch: Es gibt keinen Diff mehr
  (`ChangedPaths`/`diffTrees`/`diffTreeIndex` sind ersatzlos entfernt), keine
  Diff-Status-Klassen M/T/D/R/A, keine (ab-/an-)geschaltete
  Rename-Erkennung und keinen dedizierten „A liest den BASE-Stand nach"-Zweig
  (dieser ist mit dem alten `VCSAdded`-Fall komplett entfallen). Die
  Spec-Stelle ist genau die, die der Slice-Kopf selbst als geändert benennt
  — sie beschreibt nach diesem Commit einen Mechanismus, den der Code nicht
  mehr hat, statt der tatsächlichen Mengen-Differenz
  (`baseSet`/`headSet`/`protectedSet`, `internal/hexagon/core/rules/vcs.go:30-52`).
- **Verifizierbar:** ja — Textvergleich `spec/spezifikation.md` Schritt 2
  gegen `internal/hexagon/core/rules/vcs.go`/`internal/adapter/driven/git/git.go`;
  kein Gate prüft die semantische Deckung von Spec-Prosa gegen Code
  (`doc-check` prüft Struktur/Referenzen, nicht Algorithmus-Treue).
- **Klasse:** `messmethode-spec-text-nicht-nachgezogen`

## Negativbefunde

- **Gitlink-Behandlung** (`walkTree`, `internal/adapter/driven/git/git.go:117-135`):
  geprüft gegen `spec/spezifikation.md` Historie-Eintrag 2026-09-08 (*„ein
  Tree-Eintrag ohne Datei-Inhalt (Verzeichnis, Gitlink) … bleibt
  befundfrei"*) — das unbedingte Überspringen von `filemode.Submodule` (kein
  `case`, kein Fehler) ist die zugesagte, spec-dokumentierte Zusage, keine
  neue stille Lücke. `TestAllPathsGitlinkWirdUebersprungen` deckt den
  Codepfad ab. Kein Befund.
- **DoD (3)/§3 „ein bewegter Gitlink"** — die vorhandene Testfunktion deckt
  nur den *hinzugefügten* Gitlink (nur in HEAD), keinen an derselben Stelle
  *geänderten* (in BASE und HEAD, unterschiedlicher Ziel-Hash). Geprüft: der
  `switch`-Zweig in `walkTree` behandelt jeden `filemode.Submodule`-Eintrag
  identisch, unabhängig davon, ob er neu, geändert oder unverändert ist — der
  Pfad landet so oder so nie in `baseAll`/`headAll`. Kein plausibles
  Versagens-Szenario, das der vorhandene Test nicht bereits abdeckt; deshalb
  kein Finding (Anti-Pattern „kein Finding ohne Failure-Szenario"), nur zur
  Kenntnis für die Verifikation der DoD-Wortlaut-Deckung.
- **--staged-Sonderfall** (`AllPaths`, `head == driven.IndexRef` ohne
  `hasHead()`): geprüft — Rückgabe `(nil, nil, nil)` ist im Doc-Kommentar
  des Ports (`internal/hexagon/port/driven/vcs.go:24-32`) als bewusste
  Ausnahme mit Begründung *und* Herkunfts-Zeiger („Parität zum abgelösten
  `ChangedPaths`") benannt. Kein Befund.
- **Determinismus** (`sort.Strings(paths)` in `CheckVCS`,
  `internal/hexagon/core/rules/vcs.go:37`): geprüft, mit Kommentar-Verweis
  auf `DC-QA-02`, und empirisch dreifach reproduziert (s. o.). Kein Befund.
- **Hexagon-Grenze** (ADR-0005/ADR-0012): geprüft — `ignored()`/`cfg.Paths`
  bleiben in `internal/hexagon/core/rules/paths.go`; der Adapter
  (`internal/adapter/driven/git/git.go`) liefert ausschließlich rohe
  Pfad-Listen, kennt `vcs.paths` nicht. `make arch-check`: 0 Befunde. Kein
  Befund.
- **Testabdeckung der vier `CO-001`-Ausprägungen**
  (`internal/adapter/driving/cli/cli_vcs_test.go::TestVCS_UnlesbareObjekte`):
  geprüft — alle vier Unterfälle (BASE-Blob, BASE-Tree mit Pendant, BASE-Tree
  ohne Pendant, HEAD-Tree) laufen End-to-End gegen ein echtes, per
  `initVCSRepo`/`commitAll` erzeugtes git-Repo mit real entferntem losen
  Objekt (`removeLooseObject`), nicht gegen `fakeVCS`. Zwei der vier
  (BASE-Tree ohne Pendant, HEAD-Tree) zusätzlich unabhängig gegen eigene,
  außerhalb der Test-Suite gebaute Probe-Repos verifiziert (s. o.). Kein
  Befund.
- **Inline-Suppression / Gate-Lockerung ohne ADR** (`AGENTS.md` §3.2/§3.6):
  geprüft — kein `//nolint`, keine Schwellen-Senkung. Kein Befund.
- **Netzzugriff außerhalb `external`** (`DC-QA-03`): geprüft — reine
  git-Objekt-Lektüre über go-git, kein Netz, kein Schreiben. Kein Befund.
- **Verbleibende `ChangedPaths`/`VCSChange`/`VCSStatus`/`diffTrees`-Referenzen**
  außerhalb der beiden unter R1-F-1 verwandten Kommentar-Fundstellen (siehe
  unten): repo-weit gegrept — keine weiteren Referenzen in `*.go` (der
  Port-Kommentar in `internal/hexagon/port/driven/vcs.go:32` nennt
  `ChangedPaths` bewusst als Herkunfts-Zeiger, nicht als aktive Referenz).
  Kein Befund an dieser Stelle (die zwei stale Kommentare stehen separat
  unten).

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** `messmethode-spec-text-nicht-nachgezogen` ·
`stale-kommentar-nach-mechanismus-entfernung` (2×)

## LOW-Findings (nice-to-fix)

- **R1-F-2** — `internal/hexagon/core/rules/vcs.go:76-79`
  (`vcsDeleted`-Doc-Kommentar): *„Dass ein Rename überhaupt als
  Delete-Hälfte hier ankommt, hält die Diff-Übersetzung im VCS-Adapter,
  nicht diese Funktion."* Es gibt nach diesem Commit keine
  „Diff-Übersetzung im VCS-Adapter" mehr — der Adapter liefert nur rohe
  Pfad-Listen; dass ein Pfad als „Delete-Hälfte" ankommt, entscheidet jetzt
  die Mengen-Differenz in `CheckVCS` selbst (`headSet[p]`), nicht der
  Adapter. Ein Engineer, der wegen eines Rename-Bugs hierher navigiert,
  wird auf eine nicht mehr existierende Komponente verwiesen. Quelle:
  `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist). Verifizierbar: ja,
  Code-Lesung. Klasse: `stale-kommentar-nach-mechanismus-entfernung`.
- **R1-F-3** — `internal/adapter/driving/cli/cli_vcs_test.go:95`
  (`TestVCS_RangeDrift`-Doc-Kommentar): *„… deckt setupVCS/Open/
  ChangedPaths/FileAt/CheckVCS ab"* — `ChangedPaths` existiert nicht mehr;
  der Kommentar wurde beim Umbau nicht nachgezogen (die Zeile lag außerhalb
  der geänderten Hunks). Quelle: `AGENTS.md` §3.7. Verifizierbar: ja,
  `grep -n ChangedPaths internal/adapter/driving/cli/cli_vcs_test.go`.
  Klasse: `stale-kommentar-nach-mechanismus-entfernung`.

## Verdikt

**Blockiert.** R1-F-1 (MEDIUM) fällt unter die Reviewer-Skill-Kategorie
*„Spec-Treue-Lücke einer Messmethode … vor Merge zu klären"* — und trifft
genau die Spec-Stelle, die der Slice-Plan selbst unter „Berührte
Spec-Stellen" als geändert benennt. `spec/spezifikation.md` ist Rang 2 der
Source Precedence (technisch verbindlich, fortschreibbar); dass Schritt 2
weiterhin einen entfernten Diff-Mechanismus beschreibt, ist kein
kosmetischer Rückstand, sondern eine normative Aussage, die dem
tatsächlichen Verhalten widerspricht. Vor der nächsten Closure-Runde ist
Schritt 2 (und die beiden Rename-/Added-Absätze darunter) auf die
Mengen-Differenz umzuschreiben, mit einem neuen Historie-Eintrag nach dem
etablierten Muster (vgl. den 2026-09-08-Eintrag zur selben Anforderung).
R1-F-2/R1-F-3 (LOW) sind unabhängig davon zu beheben, halten die Closure
für sich genommen aber nicht auf.
