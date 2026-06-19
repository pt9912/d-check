# Review — slice-030 Implementierung (`--suggest-config ai-harness` / `ai-harness-init`)

## Kopf-Metadaten

- **Review-Art:** Code-Review (Diff gegen Plan/Spec/Anforderungen/Hard
  Rules — kein Verifier; DoD-Abhaken und Gate-Lauf-Bestätigung sind nicht
  Gegenstand, Gates werden NICHT als grün angenommen).
- **Datum:** 2026-06-19
- **Gegenstand:** Working-Tree-Diff (unstaged) der slice-030-Umsetzung
  plus die untrackte Slice-Datei
  `docs/plan/planning/done/slice-030-suggest-config-ai-harness.md` — die
  ai-harness-course-Vorlage für `--suggest-config`, aufgeteilt in **zwei
  explizite Modi**: `ai-harness-init` (Mode 1, Voll-Kanon — alle Blöcke
  aktiv ohne Existenzprüfung) und `ai-harness` (Mode 2, repo-bewusst —
  fehlende Pfade/Targets auskommentiert). Neuer Render-Pfad in
  `internal/hexagon/core/suggest.go`
  (`renderHarness`/`renderHarnessIDs`/`renderHarnessMatrix`,
  `pathExists`/`existingRoots`, `harnessIDPatterns`/`harnessClasses`),
  Dispatch in `SuggestConfig`, fünf neue Akzeptanztests, Lastenheft
  (`DC-FA-CLI-006`, v0.18.1), Spezifikation (`DC-FA-CLI-006.a`),
  Operations/Handbuch.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md`
  v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/done/slice-030-suggest-config-ai-harness.md`;
  Anforderungen `DC-FA-CLI-006` (Haupt — beide Modi + AK),
  `DC-FA-CONF-001` (Parser-Treue), `DC-QA-02` (Determinismus), `DC-QA-03`
  (read-only/netzlos); Spezifikation §`DC-FA-CLI-006.a`; Referenz-Konvention
  `.d-check.yml`; Hard Rules `AGENTS.md` §3; ADR-0005 (Hexagon-Import-
  Regeln); MR-006 (Referenzrichtung). **Die DoD-Abhakung des Slices lag
  bewusst nicht vor** — die Findings sind eigenständig gebildet.

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

Keine.

### INFO

#### INFO-1 — Vorrang bei kombinierten Tokens `ai-harness,ai-harness-init` undokumentiert

- **Kategorie:** INFO
- **Quelle:** Maintainability (Spec-/AK-Lücke einer Verhaltensvariante)
- **Pfad:** `internal/hexagon/core/suggest.go:73` (`return renderHarness(fsys, patterns, !initMode), nil`)
- **Befund:** Sind in einem Aufruf **beide** reservierten Tokens gesetzt
  (`--suggest-config ai-harness,ai-harness-init`), gewinnt deterministisch
  Mode 1: `initMode` ist gesetzt, `repoAware = !initMode = false` → Voll-Kanon,
  das repo-bewusste Auskommentieren von `ai-harness` entfällt stillschweigend.
  Das Verhalten ist im Code per Kommentar festgehalten, aber weder in
  `spec/spezifikation.md` §`DC-FA-CLI-006.a` noch in den Lastenheft-AK
  benannt; ein Nutzer, der beide Tokens kombiniert, erhält die Voll-Kanon-
  Ausgabe ohne dokumentierte Erwartung. Kein Korrektheitsdefekt (eindeutig
  und stabil), nur eine undokumentierte Annahme.
- **Verifizierbar:** ja — ein Akzeptanztest mit beiden Tokens würde die
  Mode-1-Ausgabe (keine „fehlt im Repo"-Hinweise) zeigen; ein Spec-Diff
  gegen §`DC-FA-CLI-006.a` würde das Fehlen der Vorrang-Regel belegen.

#### INFO-2 — Kein dedizierter DC-QA-03-Test für den `ai-harness-init`-Pfad

- **Kategorie:** INFO
- **Quelle:** DC-QA-03 (Seiteneffektfreiheit) — Slice-DoD nennt „read-only"
  für beide Modi
- **Pfad:** `internal/adapter/driving/cli/cli_acceptance_test.go`
  (Test `TestCLI006_AiHarnessInit_VollKanon`)
- **Befund:** Der read-only-Nachweis (kein `.d-check.yml` im Repo nach dem
  Lauf) ist nur für Mode 2 explizit geführt
  (`TestCLI006_AiHarness_Happy`). Der Mode-1-Test
  `TestCLI006_AiHarnessInit_VollKanon` prüft Voll-Kanon-Inhalt und
  Decodierbarkeit, stellt aber nicht fest, dass der Lauf das Repo nicht
  schreibt. Die Read-only-Eigenschaft ist strukturell gegeben (beide Modi
  laufen durch dieselbe `renderHarness`-Funktion, die nur in den `strings.Builder`
  schreibt; `SuggestConfig` macht kein Datei-I/O), daher kein
  Korrektheitsdefekt, sondern eine für Mode 1 nicht abgesicherte DoD-Zusage.
- **Verifizierbar:** ja — ein `os.Stat`-Nachweis im Mode-1-Test (wie in
  `TestCLI006_AiHarness_Happy:138` vorhanden) würde die Zusage decken; sein
  Fehlen ist im Test-Bestand sichtbar.

## Negativbefunde (geprüft, ohne Befund)

- **Mode-1-Semantik (Voll-Kanon):** `renderHarness` mit `repoAware=false`
  überspringt jede Existenzfilterung — `scan.roots` bleibt
  `[spec, docs, harness]` (`suggest.go:290-296`); in `renderHarnessIDs`
  greift die Aktiv-Bedingung `!p.always && (!repoAware || …)` für alle vier
  Target-Muster (`suggest.go:321`), das Carveout `CO-\d{3}` (`always: true`)
  bleibt korrekt auskommentiert; in `renderHarnessMatrix` ist
  `active[c.name] = !repoAware || …` für alle Klassen true, also alle Klassen
  und beide Regeln aktiv (`suggest.go:359`). `TestCLI006_AiHarnessInit_VollKanon`
  belegt vier aktive Muster, aktives ADR-Muster trotz fehlendem Target, keine
  „fehlt im Repo"-Hinweise und volle `scan.roots`.
- **Mode-2-Semantik (repo-bewusst):** `existingRoots` filtert `scan.roots`
  und `ids.scope` auf vorhandene Pfade (`suggest.go:290-296`, `312-315`);
  Muster ohne Target werden auskommentiert mit Hinweis ausgegeben
  (`suggest.go:326-336`); Klassen ohne Probe-Pfad auskommentiert
  (`suggest.go:363-370`). `TestCLI006_AiHarness_Boundary` belegt
  ADR-Muster/-Klasse auskommentiert und `roots: [spec, docs]` (harness
  gefiltert).
- **Baumelnde matrix-Regel vermieden:** beide Regeln tragen `spec-straten`
  als `from`; eine Regel ist nur aktiv, wenn **beide** Endpunkt-Klassen
  aktiv sind (`!active[r[0]] || !active[r[1]]` → `pre = "# "`,
  `suggest.go:372-377`). Damit referenziert keine aktive Regel eine
  auskommentierte Klasse — `applyMatrix` (`configyaml.go:320-323`) würde
  eine Regel mit undeklarierter Klasse mit Fehler ablehnen; dieser Fall
  ist konstruktiv ausgeschlossen.
- **Determinismus (DC-QA-02):** Die einzige Map (`active` in
  `renderHarnessMatrix`, `suggest.go:357`) wird ausschließlich per Lookup
  gelesen; alle Ausgabe-Iterationen laufen über feste Slices
  (`harnessIDPatterns()`, `harnessClasses()`, das literale
  `[][2]string{…}`-Rules-Slice, `extra`). `existingRoots` erhält die
  Eingabe-Reihenfolge. Keine `range` über eine Map im Renderpfad.
  `TestCLI006_AiHarness_Determinismus` (10× byte-identisch) belegt Mode 2;
  für Mode 1 ist die Ableitung pfad-unabhängig fix.
- **Parser-Treue (DC-FA-CONF-001):** `Decode` (`configyaml.go:107`) prüft
  **keine** Target-Existenz — die ids-Target-Prüfung liegt in
  `ensureIDTargetsExist` (`run.go:183`), das erst beim Lauf
  (`run.go:31`) greift, nicht im Decode. Damit dekodiert der aktive Teil in
  Mode 1 trotz fehlender Targets (`TestCLI006_AiHarnessInit_VollKanon:221`)
  und in Mode 2 auch degeneriert/leer: leere `patterns:`/`classes:`/`rules:`-
  Blöcke binden an nil-Slices (`rawIDs.Patterns`, `rawMatrix.Classes/Rules`,
  `configyaml.go:32-68`), der weggelassene `ids.scope` bei leerem Scope
  (Guard `len(scope) > 0`, `suggest.go:316`) ist zulässig. Die
  Modul-Liste wird beim Decode nicht validiert (`configyaml.go:124`), die
  fünf gerenderten Module sind regulär.
- **Dispatch/Abgrenzung:** `SuggestConfig` zweigt `ai-harness`/`ai-harness-init`
  per `switch` ab, bevor `resolveConfigPath`/`fsys.Kind` läuft
  (`suggest.go:45-54`) — die reservierten Tokens werden nie als Pfad
  aufgelöst, also kein „Quelle existiert nicht"/Exit 2.
  `TestCLI006_AiHarness_KeinQuellenfehler` belegt Exit 0 ohne `ai-harness`-Datei.
  Echte Quellen landen in `realSrc` und durchlaufen den Ableitungs-Pfad; ihre
  Muster werden als `extra` an `renderHarnessIDs` angehängt (`suggest.go:338-346`).
- **read-only (DC-QA-03):** `SuggestConfig` und alle Render-Helfer schreiben
  nur in einen `strings.Builder`; die Existenz-Prüfung nutzt `fsys.Kind`
  (`suggest.go:221`) → `os.Lstat` (`fs`-Adapter), kein Schreibzugriff,
  kein git, kein Netz. Das CLI gibt das Ergebnis auf `stdout` aus
  (`cli.go:289`), der Aufrufer leitet um. `TestCLI006_AiHarness_Happy:138`
  prüft, dass kein `.d-check.yml` entsteht.
- **ADR-0005-Imports:** `suggest.go` importiert nur `fmt`, `regexp`, `sort`,
  `strings` (stdlib, keine I/O-APIs) und den Port `…/hexagon/port/driven`
  (`suggest.go:3-10`) — kein `os`/`io/fs`/`net`, kein `internal/adapter/*`,
  kein `gopkg.in/yaml.v3`. Der Kern bleibt I/O-frei; das Datei-I/O liegt
  hinter dem `driven.Filesystem`-Port.
- **Spec-/Doku-Treue (zwei Modi konsistent):** Lastenheft
  (`spec/lastenheft.md` §DC-FA-CLI-006, v0.18.1), Spezifikation
  (`spec/spezifikation.md` §DC-FA-CLI-006.a), `operations.md` (Options-Tabelle +
  Abschnitt) und Benutzerhandbuch (§4.4-Block) beschreiben durchgängig
  dieselben zwei Modi (Mode 1 Voll-Kanon/leeres Repo, Mode 2
  repo-bewusst), die Baseline `v1.3.0`, dieselbe kanonische Vorlage und die
  Kombinierbarkeit. Die neue AK „`ai-harness-init` Voll-Kanon"
  (`lastenheft.md:235`) deckt sich mit `TestCLI006_AiHarnessInit_VollKanon`.
  Das YAML-Beispiel der Spec (`spezifikation.md:145-180`) spiegelt die
  Render-Ausgabe und die `.d-check.yml`-Konvention (vier Muster mit
  `link-policy: always` + `exempt-paths`, gleiche matrix-Klassen/-Regeln).
- **Kanon-Spiegelung der Vorlage:** `harnessIDPatterns`/`harnessClasses` und
  `harnessExempt` (`suggest.go:248-272`, `:32`) entsprechen der Referenz
  `.d-check.yml` (ADR/MR/DC-Muster + Targets + `exempt-paths`, spec-straten/
  adr/slice-Klassen, beide Referenzrichtungs-Regeln, `status.forbidden`,
  `exclude-sections`); ergänzt um das in `.d-check.yml` nicht separat
  geführte `slice-\d{3}`-Muster (in der Spec-Vorlage und den AK gefordert)
  und das auskommentierte `CO-\d{3}`-Carveout.
- **MR-006 (Referenzrichtung):** Die neuen Spec-Stellen (`spezifikation.md`,
  ein Spec-Stratum) verlinken nur auf `lastenheft.md`-Abschnitte
  (auf-/seitwärts), nicht abwärts auf ADRs/Slices. Die Lastenheft-Prosa zu
  `DC-FA-CLI-006` enthält keinen Markdown-Link auf `docs/plan/*`; die
  `DC-FA-SCAN-001`/`DC-FA-ID-001`-Links zeigen auf dasselbe Lastenheft.
  Das `slice-030`-Token in der Änderungshistorie (`lastenheft.md:779`) ist
  ein unverlinktes Klartext-Label im Muster aller Vorzeilen — kein
  abwärts-Link. Keine neue MR-006-Bedingung.
- **codepaths/ids-Sauberkeit der neuen Doku/Spec:** Die kanonischen Regex-
  Muster und Pfade der Vorlage stehen in der Spezifikation in einem
  ```yaml-Fence (`spezifikation.md:145-180`) — vom `ids`/`codepaths`-Scan
  exempt. In der Prosa erscheinen IDs nur als verlinkte Backtick-Spans
  (`[`DC-FA-CONF-001`](…)`) oder als Regex in Inline-Code (`CO-\d{3}` —
  kein linkpflichtiges Präfix); keine bare `DC-`/`MR-`/`ADR-`/`slice-`-ID
  außerhalb Fence/Link in den neuen Prosa-Zeilen. Alle vom Slice/Spec
  referenzierten Lastenheft-Anker (`#dc-fa-cli-006…`, `#dc-fa-conf-001…`,
  `#dc-fa-scan-001…`, `#dc-fa-id-001…`, `#dc-qa-02…`, `#dc-qa-03…`) lösen
  auf vorhandene Headings auf.
- **Slice-Plan-Treue:** Die im Plan (§3) benannten Dateien sind genau die
  geänderten; die DoD-Punkte (reservierte Schlüsselwörter ohne Exit 2,
  zwei-Modi-Vorlage mit Carveout auskommentiert, Spec-Erweiterung,
  Tests inkl. Determinismus/read-only, Doku) sind im Diff abgebildet.
  (Die DoD-*Abhakung* ist nicht Reviewer-Sache.)

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 0 | 2 |

## Verdikt

**Freigegeben.** Keine HIGH/MEDIUM/LOW offen. Der Kern des Slices — saubere
Token-Abzweigung vor der Pfad-Auflösung (kein Exit 2 für die reservierten
Wörter, Kombi mit echten Quellen), korrekte Mode-1-Semantik (Voll-Kanon
ohne Existenzprüfung, Carveout dennoch auskommentiert), korrekte
Mode-2-Semantik (repo-bewusstes Auskommentieren, baumelnde matrix-Regel
konstruktiv vermieden), deterministischer Renderpfad (Map nur als Lookup,
Iteration über feste Slices), Parser-Treue in beiden Modi (Decode prüft
keine Target-Existenz; aktive Teile dekodieren — auch degeneriert/leer),
Layer-Konformität (Kern I/O-frei, ADR-0005), und durchgängig konsistente
Spec/Doku zu den zwei Modi — ist tragfähig umgesetzt. Die zwei INFO-Punkte
(undokumentierter Vorrang bei kombinierten Tokens; für Mode 1 nicht
abgesicherte read-only-Zusage) sind dokumentationswürdige Annahmen ohne
Korrektheits- oder Sicherheitswirkung und blockieren die Freigabe nicht.
Die Gate-Bestätigung obliegt der getrennten Verifikation (hier nicht als
grün angenommen).
