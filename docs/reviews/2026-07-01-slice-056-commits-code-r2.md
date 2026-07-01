# Review R2 — slice-056 Modul `commits` (Go-Code / Gate-Umbau / Tests)

## Kopf-Metadaten

- **Datum:** 2026-07-01
- **Reviewer:** unabhängig/adversarial (Code-Fokus; Doc-Straten prüft R1)
- **Gegenstand:** slice-056 — opt-in Modul `commits` (Commit-Message-Traceability),
  Portierung des entfernten `tools/trace-check.sh` über den reine-Go-VCS-Port.
- **Betroffene IDs:** `DC-FA-COMMITS-001`, `DC-FA-CLI-002/003/010`, `DC-QA-02/03/04`,
  ADR-0027 (löst die Skript-Mechanik von ADR-0013 ab), MR-007.
- **Diff:** `git diff HEAD` + neue Dateien (`commits.go`, `commits_test.go`,
  `cli_commits_test.go`, ADR-0027, slice-056-Plan); `tools/trace-check.sh` gelöscht.
- **Prüf-Fokus:** clean_message-Parität, ID-/Exempt-Logik, fail-closed/Silent-Grün,
  git-Adapter `CommitMessages`, arch-check R2, CLI-Verdrahtung, Gate-/Hook-Umbau, Tests.
- **Nicht geprüft (Rolle):** DoD-Abhakung, Plan-/Doc-Konsistenz, Prosa-Straten.
- **Methode:** Code-Lesung + reale Proben gegen `d-check:latest` (Image 41 min alt,
  kennt `commits` — spiegelt den neuen Code); Paritäts-Batterie gegen das aus git
  wiederhergestellte `tools/trace-check.sh`. Kein neues Image gebaut.

## Verifikations-Belege (reproduzierbar)

- **Paritäts-Batterie (14 Fälle)** altes Skript `--message` vs. `--commit-msg -`
  über identische Message-Dateien: **alle 14 OK, keine Divergenz** — u. a.
  ID-nur-in-#-Kommentar, ID-nach-scissors, scissors ohne Leerzeichen (`#x>8y`),
  CRLF-Kommentar-ID, CRLF-Inhalts-ID, führender #-Kommentar als Betreff, leere
  Message, `merge` kleingeschrieben (nicht ausgenommen).
- **Range-Modus real** (`--enable commits --range`): rev-list-Parität (Set == `git
  rev-list --no-merges base..head`), Basis ausgeschlossen, umgekehrte Range leer
  (Exit 0), zwei Läufe byte-identisch (DC-QA-02), Ausgabe SHA-aufsteigend sortiert.
- **fail-closed:** commits aktiv ohne Range → Exit 2; `--staged` → Exit 2
  (IndexRef nicht unterstützt); unauflösbare Nullbasis → Exit 2; `--commit-msg`
  ohne id-patterns → Exit 2; `--commit-msg + --range` → Exit 2 (comboError).
- **Config-Guards** (Probe): `id-patterns: ['.*']` → Exit 2 „matcht den Leerstring";
  leerer Eintrag → Exit 2; `[unclosed` → Exit 2.

---

## Findings

### MEDIUM-1 — `applyCommits` Config-Ablehnungspfade ohne Negativtest (Konventionsbruch)

- **kategorie:** MEDIUM
- **quelle:** DC-QA-04 / Reviewer-Skill „fehlende Negativtests bei neuem
  öffentlichen Vertrag"
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:281-305` (`applyCommits`);
  Lücke in `internal/adapter/driven/configyaml/configyaml_test.go`
- **befund:** `applyCommits` trägt drei Ablehnungs-Guards — leerer id-patterns-Eintrag,
  **Leerstring-matchendes** Muster (`re.MatchString("")` → Exit 2, der
  Silent-Grün-Schutz: sonst gälte jede Message als getraced) und nicht kompilierbares
  Regex. Kein einziger dieser Pfade ist im configyaml-Test abgedeckt; `cli_commits_test.go`
  fährt ausschließlich gültige Configs, und `TestCommits_FailClosed` testet nur den
  Laufzeit-Guard in `runCommitMsg`, nicht die Decode-Guards. Das bricht eine im Repo
  durchgängige Konvention: `TestDecode_VersionsFehler` (pin-pattern leer/`x*`/`[unclosed`),
  `TestDecode_LeerstringRegex` (ids) und der matrix-token-Negativtest prüfen exakt diese
  Guard-Klasse für jedes andere Modul. Failure-Szenario: entfernt ein künftiger Refactor
  den `re.MatchString("")`-Guard (oder invertiert die Bedingung), fällt kein Test — ein
  `id-patterns: ['.*']` machte das Gate danach still grün (jede Message „getraced"), ohne
  dass die Test-Suite es fängt.
- **verifizierbar:** ja — `make test` (bzw. `go test ./internal/adapter/driven/configyaml/`)
  bliebe grün, wenn man den `MatchString("")`-Guard in `applyCommits` entfernte; ein
  analoger Testfall wie `TestDecode_VersionsFehler` würde ihn binden. Die Guards
  funktionieren aktuell (Probe: `.*`/leer/`[unclosed` je Exit 2) — die Lücke ist der
  fehlende Regressionsschutz, kein Laufzeitfehler.

---

### INFO-1 — Range-Modus inert bei leeren `id-patterns` (Asymmetrie zum Message-Modus)

- **kategorie:** INFO
- **quelle:** DC-QA-04 / Maintainability (Gate-Exposition)
- **pfad:** `internal/hexagon/core/rules/commits.go:24` (`CheckCommits` early-return)
  vs. `internal/adapter/driving/cli/cli.go:210` (`runCommitMsg`-Guard Exit 2)
- **befund:** Bei aktivem `commits`, aber leeren `id-patterns` (z. B. ganzer
  `commits:`-Block aus `.d-check.yml` entfernt) ist der Range-/Scan-Modus **inert**
  (Exit 0) — verifiziert: `--enable commits --range base..head` über einen
  kennungslosen Commit meldet „0 Befund(e)". Der Message-Modus (`--commit-msg`)
  bricht bei derselben Bedingung fail-closed ab (Exit 2). Diese Asymmetrie ist
  **kein Defekt**: der Range-Inert ist das akzeptierte, familienweite opt-in-Vertrag
  — `CheckVCS` (`internal/hexagon/core/rules/vcs.go:20`) hat das byte-identische
  Muster `if len(cfg.Paths)==0 || vcs==nil { return nil,nil }` und wurde in slice-053
  (v0.33.0) so abgenommen; der Message-Modus-Guard ist eine sinnvolle Sonderbehandlung
  (ein dedizierter Ein-Message-Modus wäre inert nutzlos). Tippfehler im
  Config-Schlüssel fängt der strikte YAML-Decoder (KnownFields → Exit 2). Erwähnenswert
  bleibt die Gate-Exposition für Konsumenten des verteilten `doc-commits`-Targets, die
  es ohne konfigurierte `id-patterns` verdrahten: dieselbe Klasse wie das bereits
  ausgelieferte `doc-immutable` (vcs).
- **verifizierbar:** ja — Probe A oben (Exit 0 trotz kennungslosem Commit bei fehlendem
  `commits:`-Block).

### INFO-2 — Über-permissives `id-patterns`-Muster passiert still

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:296` (`re.MatchString("")`-Guard)
- **befund:** Der Guard lehnt nur Leerstring-**matchende** Muster ab. Ein Muster wie
  `.` (matcht jede nicht-leere Zeile) passiert die Validierung und lässt jeden Commit
  als „getraced" durchgehen — verifiziert: `id-patterns: ['.']` → Exit 0 über den
  kennungslosen Commit. Inhärent jeder Regex-Konfiguration (auch das Alt-Skript war
  über eine editierte `ID_RE` brechbar); dokumentationswürdige Annahme, kein Codefehler.
- **verifizierbar:** ja — Probe D oben.

### INFO-3 — `--commit-msg` mit früheren Kurzschluss-Flags still ignoriert

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `internal/adapter/driving/cli/cli.go:402/411/443` (Reihenfolge der Run-Kurzschlüsse)
- **befund:** `comboError` sperrt `--commit-msg` gegen `--range/--staged/--trace/
  --doctor/--repair`, nicht gegen `--print-config/--print-mk/--suggest-config`. Wird
  `--commit-msg` mit einem dieser Generatoren kombiniert, gewinnt der frühere
  Kurzschluss und `--commit-msg` wird stillschweigend ignoriert. Harmlos (reine
  Print-and-Exit-Modi, kein Repo-Schreiben, kein verfälschtes Gate-Ergebnis).
- **verifizierbar:** ja (Code-Reihenfolge; kein Gate-Effekt).

### LOW-1 — `FOCUS_DISABLE` handgepflegter Spiegel der `.d-check.yml`-Module

- **kategorie:** LOW
- **quelle:** Maintainability (latente Wartungsfalle)
- **pfad:** `Makefile:196-197` (`FOCUS_DISABLE`)
- **befund:** `FOCUS_DISABLE` listet genau die acht Default-Module aus
  `.d-check.yml` `modules:` (aktuell deckungsgleich — verifiziert). Wird dort ein
  neuntes Default-Modul ergänzt und hier nicht nachgezogen, liefe dieses in
  `make trace-check`/`make adr-check` auf den Arbeitsbaum-Inhalt mit — **lautes**
  Über-Feuern (Falsch-Positiv), kein Silent-Grün. Der Kommentar warnt explizit; die
  Umbenennung `VCS_DISABLE→FOCUS_DISABLE` (geteilt von beiden Fokus-Targets) ist eine
  DRY-Verbesserung. `--print-mk`/`disableAllExcept` ist demgegenüber selbst-ableitend
  (aus `ValidModules`) und braucht keine Pflege.
- **verifizierbar:** ja — bei künftigem `modules:`-Zuwachs würde `make trace-check`
  über das nicht-abgewählte Modul im Arbeitsbaum feuern.

---

## Negativbefunde (geprüft, ohne Befund)

- **clean_message-Parität:** `cleanCommitMessage` (commits.go:71) bildet
  `sed -e '/^#.*>8/,$d' -e '/^#/d'` exakt ab — scissors-Break dropt den Rest
  inkl. der Zeile, `#`-Präfix-Skip danach. 14/14 Grenzfälle byte-gleich zum
  wiederhergestellten Skript; keine divergierende Message-Klasse gefunden.
- **ID-/Exempt-Logik:** `IDPatterns` (.d-check.yml `ADR-\d{4}`/`MR-\d{3}`/
  `DC-(FA-[A-Z]+|QA)-\d+`/`slice-\d+`) deckungsgleich mit dem Skript-`ID_RE`;
  ID-Suche über die **bereinigte** Message ⇒ eine ID nur auf einer `#`-Zeile zählt
  nicht (Probe E/L). Exempt betreff-basiert auf der bereinigten ersten Zeile wie
  `is_exempt`.
- **git-Adapter `CommitMessages`:** rev-list-Parität, `len(ParentHashes)<=1`-Merge-
  Filter (TestCommitMessagesSkipsMerges + Code), SHA-Sortierung deterministisch,
  `ancestors`-Ausschluss korrekt für vorwärts/umgekehrt/leer, `.git` rein lesend
  (ResolveRevision/CommitObject; `:ro`-Mount ohne Schreibfehler), IndexRef-Guard.
- **arch-check R2:** `go-git/go-git/v5*` ausschließlich in
  `internal/adapter/driven/git/git.go` importiert; `commits.go` importiert nur
  `regexp`/`strings`/`model`/Port — Kern rein/fakebar (`fakeVCS`).
- **CLI-Verdrahtung:** `--commit-msg`-Kurzschluss (Run:443) VOR
  `EffectiveModules`/`resolveVCS`/Scan; kein Repo-Scan, kein Schreiben; stdin bei
  `-` via `io.ReadAll` (TestCommits_CommitMsgStdin); `gitModuleActive` erfasst
  `commits`; comboError-Kombinationen fail-closed.
- **Gate-/Hook-Umbau:** `make trace-check` fährt beide Modi (MSGFILE→`DCHECK_RUN_I`
  `-i` stdin; sonst Range, Default `HEAD~1..HEAD`); `FOCUS_DISABLE` == 8 Default-Module
  (nur `commits` läuft); Hook `exec make trace-check MSGFILE="$1"` propagiert non-zero;
  CI-Workflow ruft `make trace-check RANGE=...` (kein Skript-Bezug); Skript-Löschung
  vollständig (keine ausführenden Referenzen mehr; Rest tombstoned via
  `codepaths.ignore-refs`).
- **Selbsttest-Klassen:** commits_test.go deckt ID erkannt (adr/mr/dc/qa/slice),
  fehlend, Merge/Revert exempt, ID-nur-in-#, ID-nach-scissors, Body-Zeile, leere
  Message, no-exempt; git_test.go Range/Merges/fail-closed; cli_commits_test.go
  Range/CommitMsg/stdin/exempt/fail-closed-Combos. **Einzige Lücke:** die
  `applyCommits`-Decode-Guards (MEDIUM-1).

## Kategorie-Summary

| Kategorie | Anzahl |
|-----------|--------|
| HIGH      | 0      |
| MEDIUM    | 1      |
| LOW       | 1      |
| INFO      | 3      |

## Verdikt

**NACHBESSERN** (ein MEDIUM). Der Kern ist korrekt und in echter Parität zum
abgelösten `tools/trace-check.sh` — keine clean_message-/ID-/Range-Divergenz, kein
arch-check-Verstoß, fail-closed durchgängig, DC-QA-02/03 gewahrt. Das MEDIUM ist
kein Korrektheits-/Parität-Defekt, sondern eine **Testlücke**: die
Config-Ablehnungspfade von `applyCommits` (insbesondere der Leerstring-Match-Guard
als Silent-Grün-Schutz) sind — entgegen der repo-weiten Konvention
(`TestDecode_VersionsFehler` u. a.) — ungetestet. Empfehlung an die Implementation:
Regressionstest analog `TestDecode_VersionsFehler` ergänzen; die INFO/LOW-Notizen
sind optional.
