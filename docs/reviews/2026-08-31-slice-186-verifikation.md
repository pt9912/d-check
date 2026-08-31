# Verifikation slice-186 — `vcs` meldet den reinen Rename auch über `--range`

**Art:** Verifikation.
**Gegenstand:** [slice-186](../plan/planning/in-progress/slice-186-vcs-rename-im-range.md), Range `5a6a75c..HEAD` — `f01b4d9` (Plan angelegt), `8117758` (Claim-Move `open/` → `in-progress/`), `84c0ca8` (Implementierung: Rename-Erkennung im Range-Diff aus), `4e9cbe8` (Test für den Rename **aus** der geschützten Klasse heraus), `b377c27`/HEAD (Review-Report + fünf Kommentar-Korrekturen).
**Modell-ID:** `claude-sonnet-5`.
**Datum:** 2026-08-31.
**Rolle:** unabhängiger Verifier — geprüft wird gegen **DoD und Spec** (Baseline-Regelwerk [`modul-08-agentenrollen.md` §Welche Rolle braucht welche Artefaktklasse](../../.harness/baseline/v5.12.0/regelwerk/modul-08-agentenrollen.md#welche-rolle-braucht-welche-artefaktklasse-modul-8)), nicht gegen Plan/ADR-Maintainability. Der Review-Report `docs/reviews/2026-08-31-slice-186-vcs-rename-code-r1.md` wurde bewusst **nicht** gelesen, um einen von der Review-Rolle unabhängigen Blick zu halten.
**Prüfgrundlage:** Slice-Plan `docs/plan/planning/in-progress/slice-186-vcs-rename-im-range.md` §4 (DoD), §1–§3; [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in) im Lastenheft (unverändert); §[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) in der Spezifikation; der eingehende Befund [`2026-08-31-befund-ai-harness-course-vcs-rename.md`](../plan/cr/2026-08-31-befund-ai-harness-course-vcs-rename.md); das tatsächlich gebaute Produkt (`make build`, `d-check:latest`, Image-ID `sha256:446ab542dac4…`), nicht der Quelltext allein.
**Nicht verändert:** kein Repo-Artefakt außer dieser Datei. Für die Umkehr-Probe wurde `internal/adapter/driven/git/git.go` temporär auf den Vor-Fix-Stand (`git show 84c0ca8~1:…`) zurückgesetzt und danach exakt auf den Post-Fix-Stand zurückkopiert (Beleg unten); `git status --short` war danach und ist bei Abgabe leer bis auf diese Report-Datei. Alle Mess-Läufe gegen ein eigenes Wegwerf-Repo unter `/tmp/.../scratchpad/vcs-verify` (außerhalb dieses Repos, keine Spur hier).

---

## 1. DoD-Tabelle (§4 des Slice-Plans)

| # | Behauptet | Gemessen | Verdikt |
|---|---|---|---|
| 1a | Der reine Rename einer immutablen Datei meldet `core-drift-vcs` **auch über `--range`** | Wegwerf-Repo: `git mv adr/0001-kern.md adr/0002-kern.md`, committet, `d-check:latest --enable vcs --range CFG..HEAD2` → **1 Befund**, `core-drift-vcs` auf `adr/0001-kern.md` (§2 Fall 1 unten). Vor dem Fix (Umkehr-Probe, s. §2 Fall B) lieferte derselbe Lauf **0 Befunde** | **erfüllt** |
| 1b | … mit Test, der **vor** dem Fix rot war (Ausgabe **in der Closure-Notiz**) | `TestRangePureRenameYieldsDelete` existiert (`git_test.go`); selbst reproduziert: mit dem Vor-Fix-`git.go` schlägt genau dieser Test fehl (`--- FAIL: TestRangePureRenameYieldsDelete`, alle anderen Pakete `ok`), mit dem Fix ist er grün. Die **Substanz** ist damit belegt. Die **Ablage** dieses Belegs in der Closure-Notiz ist der Closure-Schritt selbst — §9 des Slice-Plans ist eine leere Überschrift (Slice liegt in `in-progress/`) | Substanz **erfüllt**; Beleg-Ablage **noch nicht fällig** |
| 2 | Die drei Kontrollfälle unverändert: `--staged`-Rename, `--range`-Löschen, `--range`-Rename-mit-Umformulierung | Alle drei selbst nachgefahren (§2 Fälle 1, 3, 4) — jeweils genau **1 Befund** `core-drift-vcs`, Exit 1, identisch zur Absender-Tabelle | **erfüllt** |
| 3 | **Die zwei falschen Zusagen sind weg** (`git.go`-Kommentar, `.githooks/pre-commit`), und `AGENTS.md` §3.5 sagt wieder die Wahrheit — oder trägt die verbleibende Grenze benannt | `git.go`: der Kommentar über `diffTrees` beschreibt jetzt den tatsächlichen Mechanismus (Rename-Erkennung aus) statt die vorherige, falsche Behauptung. `.githooks/pre-commit`: die Zeile „EINE Wahrheit: ruft dasselbe Gate wie CI" steht unverändert, trägt aber jetzt vier Zeilen darunter, die die Modus-Abhängigkeit einräumen. `AGENTS.md` §3.5 selbst ist **unverändert** (kein Diff) — zu Recht: seine Aussage („`make adr-check` erzwungen über `pre-commit`-Hook + PR-/Push-CI") ist mit dem Fix wieder wahr, selbst nachgemessen (§2 Fälle 1+2: beide Modi fangen den reinen Rename). Zusätzlich, über die zwei im Plan benannten Stellen hinaus, korrigiert HEAD zwei **weitere** Kommentare mit derselben falschen Behauptung: `internal/hexagon/core/rules/vcs.go` (`vcsDeleted`) und `internal/hexagon/port/driven/vcs.go` (`VCSStatus`) — beide zeigen jetzt auf den Adapter als Halter der Eigenschaft, statt sie selbst festzustellen. Das ist eine **positive** Erweiterung über den Plan hinaus, siehe §4 | **erfüllt** |
| 4 | Spezifikation §[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) nennt den Mechanismus; Historie-Zeile gesetzt; **kein** Lastenheft-Bump | Schritt 2 trägt einen neuen Absatz: Range-Pfad difft ohne Rename-Erkennung, Preis benannt (keine Inhalts-Ähnlichkeit), `--staged` trägt dieselbe Eigenschaft über `diffTreeIndex`. Historie-Tabelle: neue Zeile `2026-08-31` mit vollständiger Begründung. `git diff 5a6a75c..HEAD -- spec/lastenheft.md` ist **leer** — kein Bump | **erfüllt** |
| 5 | Der Befund ist als eingehendes Dokument abgelegt, und die zwei CR-Dokumente vom 2026-08-31 sind um den zurückgezogenen Abschnitt korrigiert | Neue Datei `docs/plan/cr/2026-08-31-befund-ai-harness-course-vcs-rename.md`. `2026-08-31-cr-ai-harness-course-observations-relational.md` trägt einen Rückzugs-Vermerk am betroffenen Absatz; `2026-08-31-antwort-ai-harness-course-observations-relational.md` trägt einen Nachtrag-Abschnitt mit Bestätigung, Fix-Beschreibung und den beiden mitgenommenen Lehren | **erfüllt** |
| 6a | `make gates` grün (Exit explizit) | Selbst gefahren: **Exit 0**, Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; semgrep „Ran 55 rules on 55 files: 0 findings." | **erfüllt** |
| 6b | **unabhängiger Review**; **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten | Review-Commit `b377c27` trägt den Report `docs/reviews/2026-08-31-slice-186-vcs-rename-code-r1.md` in eigenem Kontext (Inhalt bewusst ungelesen, s. Kopf). Diese Datei ist die Verifikation, in ihrem eigenen, vom Review getrennten Kontext | **erfüllt** (durch diesen Report selbst eingelöst) |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register fortgeschrieben; jedes Risiko aus §5 mit Ausgang; die drei Paarungen geprüft | §9 leer; alle vier §5-Risiken tragen `*(bei Closure)*`; `git diff 5a6a75c..HEAD -- docs/plan/planning/observations.md` ist leer (Register unangetastet) | **noch nicht fällig** |

**Zusammenfassung:** 7 von 8 Zeilen **erfüllt** (1a, 2, 3, 4, 5, 6a, 6b), 1 Zeile mit **Substanz erfüllt / Beleg-Ablage noch nicht fällig** (1b), 1 Zeile **noch nicht fällig** (7 — Closure-Schritt, korrekt: Slice liegt in `in-progress/`). **0 Zeilen nicht erfüllt.**

## 2. Die fünf nachgefahrenen Fälle (eigenes Wegwerf-Repo, Produkt-Binary)

Aufbau: ein Repo mit `docs/plan/adr/0001-kern.md` (`**Status:** Accepted`) und `.d-check.yml` mit `vcs: {paths: ["docs/plan/adr/*.md"], immutable-when: '^\*\*Status:\*\* Accepted'}`, Basis-Commit `CFG = d44adf2`. Jeder Fall startet frisch von `CFG` auf einem eigenen Branch, wird danach verworfen. Alle Läufe: `docker run --rm --network none -v <repo>:/repo:ro d-check:latest …` gegen das lokal aus HEAD gebaute `d-check:latest`.

| Fall | Aufbau | Aufruf | Ausgabe | Erwartung |
|---|---|---|---|---|
| 1 | `git mv 0001-kern.md 0002-kern.md`, **staged**, nicht committet | `--enable vcs --staged` | `d-check: 1 Datei(en) geprüft, 1 Befund(e)` / `docs/plan/adr/0001-kern.md:1 … core-drift-vcs … gelöscht oder umbenannt …`, **Exit 1** | `core-drift-vcs` ✓ |
| 2 (der gemeldete Fall) | derselbe Rename, **committet** | `--enable vcs --range CFG..HEAD2` | `docs/plan/adr/0001-kern.md:1 … core-drift-vcs …`, `d-check: 1 Datei(en) geprüft, 1 Befund(e)`, **Exit 1** | `core-drift-vcs` ✓ (vor dem Fix: 0 Befunde, s. Fall B) |
| 3 | Rename **mit** vollständiger Umformulierung (Titel, Fließtext, Länge), committet | `--enable vcs --range CFG..HEAD3` | `docs/plan/adr/0001-kern.md:1 … core-drift-vcs …`, **Exit 1** | `core-drift-vcs` ✓ (Kontrollfall) |
| 4 | reines `git rm 0001-kern.md`, committet | `--enable vcs --range CFG..HEAD4` | `d-check: 0 Datei(en) geprüft, 1 Befund(e)` / `docs/plan/adr/0001-kern.md:1 … core-drift-vcs …`, **Exit 1** | `core-drift-vcs` ✓ (Kontrollfall) |
| 5 (vom Slice ergänzt) | Rename **aus** der Klasse heraus: `docs/plan/adr/0001-kern.md` → `docs/notes/kern.md` (`docs/notes/**` ist nicht in `vcs.paths`), committet | `--enable vcs --range CFG..HEAD5` | genau **1** Befund, auf dem **alten** Pfad `docs/plan/adr/0001-kern.md`, **Exit 1** | genau ein Befund auf dem alten Pfad — der neue Pfad ist frei, weil außerhalb der Klasse; deckungsgleich mit dem neuen Unit-Test `TestVCSRenameOutOfClass` |

**Fall B — Umkehr-Probe (vor dem Fix):** `internal/adapter/driven/git/git.go` auf den Stand `84c0ca8~1` zurückgesetzt (`base.Diff(head)` statt `object.DiffTreeWithOptions(…, DetectRenames: false)`), `docker build --target test .` (volles `go test ./...` im Container) → **genau ein** Fehlschlag:
```
--- FAIL: TestRangePureRenameYieldsDelete (0.01s)
    git_test.go:294: reiner Rename, adr/0001-kern.md: D erwartet, got   ok=false (alle [{77 adr/0002-kern.md}])
FAIL	github.com/pt9912/d-check/internal/adapter/driven/git	0.034s
```
alle übrigen Pakete `ok`. Danach `git.go` exakt auf den Post-Fix-Stand zurückkopiert; `git diff` gegen den ursprünglichen Arbeitsbaum-Stand ist leer.

**Die vier Zeilen der Absender-Tabelle stimmen zellenweise** mit den Fällen 1–4; Fall 2 ist der einzige, der sich durch den Fix bewegt (0 → 1 Befund), die drei anderen bleiben unverändert — genau das, was das Vorgehen (§2 Schritt 3) und Risiko 1 (§5) verlangen.

## 3. Spec-vs-Produkt

- **Mechanismus.** §[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 2 behauptet: der Range-Pfad difft ohne Rename-Erkennung, sonst kommt eine Umbenennung als eine Änderung auf dem neuen Pfad an und die Pfad-Stabilitäts-Prüfung läuft still ins Leere. Der Code (`diffTrees`, `object.DiffTreeWithOptions(…, &object.DiffTreeOptions{DetectRenames: false})`) tut genau das; Fall B oben belegt am Produkt, dass mit eingeschalteter Erkennung (dem Vor-Fix-Stand) exakt dieses stille Verschwinden eintritt.
- **Preis.** Dieselbe Textstelle benennt die Grenze: der Range-Pfad misst keine Inhalts-Ähnlichkeit — Fall 5 belegt das direkt: ein Rename mit identischem Inhalt, aber außerhalb der Klasse, wird nicht als „derselbe Inhalt, anderer Pfad" erkannt, sondern der alte Pfad-Befund und der neue freie Pfad stehen unabhängig nebeneinander (kein Ähnlichkeits-Matching, wie zugesagt).
- **Kein Lastenheft-Bump — geprüft, nicht nur behauptet.** [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in) sagt bereits zu: „**gelöschte** oder **umbenannte** immutable Datei → Befund `core-drift-vcs` (der Pfad einer immutablen Datei ist stabil)" — ohne Modus-Einschränkung. Stärker noch: die **Negative**-Akzeptanzkriterium-Zeile benennt bereits ausdrücklich den `--range`-Aufruf für genau diesen Fall („… oder die Datei gelöscht/umbenannt wird, when `d-check --enable vcs --range <base>..<head>` läuft, then ein Befund `core-drift-vcs`"). Der Vor-Fix-Code verletzte damit sein **eigenes**, unveränderte Akzeptanzkriterium — die Behauptung „die Anforderung sagte den Befund bereits zu, sie wurde nur nicht eingelöst" trägt, es ist kein Modus-Vorbehalt aufzufinden, der einen Bump nötig machen würde.

## 4. Plan-vs-Code-Diff, beide Richtungen

**Richtung 1 — im Diff, aber vom Plan nicht wörtlich angekündigt:**

- Zwei zusätzliche Kommentar-Korrekturen (`internal/hexagon/core/rules/vcs.go`, `internal/hexagon/port/driven/vcs.go`) über die im Plan (§2 Schritt 4) namentlich genannten zwei Stellen (`git.go`, `.githooks/pre-commit`) hinaus — die Commit-Botschaft von HEAD legt das offen („fünf Befund-Korrekturen: zwei weitere Kommentare, die dieselbe Eigenschaft behaupteten, zeigen jetzt auf ihren Halter statt sie festzustellen"). Beide sind reine Kommentar-Korrekturen ohne Verhaltensänderung, tragen dieselbe Korrektur-Absicht wie DoD-Punkt 3 und verletzen keine der fünf §3-Grenzen. Bewertung: **positive, offengelegte Erweiterung**, kein Plan-Verstoß.
- `docs/plan/planning/in-progress/roadmap.md` (−2 Zeilen „Nichts in Arbeit."): die deklarierte Ausnahme aus `AGENTS.md` §3.3/[`MR-013`](../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise), die den Lifecycle-Move-Commit begleitet — Konvention, kein unangekündigter Fund.
- Der Review-Report selbst (`docs/reviews/2026-08-31-slice-186-vcs-rename-code-r1.md`) — Artefakt der Review-Rolle, kein Plan-Liefer-Punkt, erwartungsgemäß im selben Commit wie die Kommentar-Korrekturen, die er ausgelöst hat.

**Richtung 2 — vom Plan zugesagt, aber im Diff nicht wiederzufinden:** keine Lücke gefunden. Beide Liefer-Punkte aus §2 (Fix + Umkehr-Probe-Test; Spezifikation + Dokumente) sind im Diff vollständig vertreten.

**§3 „Ausdrücklich NICHT" — alle fünf Grenzen gehalten:**

1. **Keine Ähnlichkeits-Erkennung** — `DetectRenames: false`, kein Schwellenwert, kein Ähnlichkeits-Code irgendwo im Diff; Fall 5 bestätigt das am Produkt.
2. **Kein neuer Grund-Code** — `core-drift-vcs` unverändert (kein Diff im Grund-Code-Abschnitt der Spezifikation), alle fünf Fälle oben melden denselben Code.
3. **Keine Vertrags-Änderung** — `spec/lastenheft.md` unverändert (`git diff 5a6a75c..HEAD -- spec/lastenheft.md` leer).
4. **Keine Änderung am `--staged`-Pfad** — `diffTreeIndex` textuell unverändert (nur in einem neuen Kommentar an anderer Stelle *erwähnt*); Fall 1 bestätigt unverändertes Verhalten.
5. **Kein CHANGELOG-Eintrag im Feature-Commit** — `git diff 5a6a75c..HEAD -- CHANGELOG.md` leer.

## 5. Gate-/Mess-Läufe (echte Ausgabe)

| Lauf | Ergebnis |
|---|---|
| `make build` | Neubau aus HEAD, `d-check:latest`, Image-ID `sha256:446ab542dac4…` |
| `make test` (Post-Fix) | alle Pakete `ok`, u. a. `internal/adapter/driven/git` und `internal/hexagon/core/rules` |
| Fall 1 — `--staged`, reiner Rename | `core-drift-vcs`, **Exit 1** |
| Fall 2 — `--range`, reiner Rename committet | `core-drift-vcs`, **Exit 1** (vor dem Fix: 0 Befunde) |
| Fall 3 — `--range`, Rename mit Umformulierung | `core-drift-vcs`, **Exit 1** |
| Fall 4 — `--range`, reines Löschen | `core-drift-vcs`, **Exit 1** |
| Fall 5 — `--range`, Rename aus der Klasse heraus | 1 Befund auf dem alten Pfad, **Exit 1** |
| Umkehr-Probe: `docker build --target test .` mit Vor-Fix-`git.go` | **Exit 1** (Build scheitert am `go test`), genau `--- FAIL: TestRangePureRenameYieldsDelete`, sonst alles `ok` |
| Rückstellung `git.go` auf Post-Fix-Stand | `git status --short` leer |
| `make adr-check RANGE="5a6a75c..HEAD"` (Sanity, eigener Bestand) | `d-check: 639 Datei(en) geprüft, 0 Befund(e)`, **Exit 0** |
| `make adr-check` (Default-Range `HEAD~1..HEAD`, Sanity) | `d-check: 639 Datei(en) geprüft, 0 Befund(e)`, **Exit 0** |
| `make gates` | **Exit 0** — `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; Coverage 94,70 % ≥ 93 %; semgrep 55 Regeln, 0 Befunde |

## 6. Verdikt

**Bestanden.** Alle fälligen DoD-Punkte (1a, 2, 3, 4, 5, 6a, 6b) sind gegen das laufende Binary erfüllt — inklusive einer selbst reproduzierten Umkehr-Probe, die exakt den behaupteten einen Test rot schaltet, und einer eigenständigen Nachfahrt aller vier gemeldeten Fälle plus des vom Slice ergänzten fünften. Die Substanz von 1b ist erfüllt, ihre Beleg-Ablage in der Closure-Notiz ist ein Closure-Schritt und planmäßig offen; Punkt 7 ist vollständig Closure. Spezifikation und Produkt stimmen überein, einschließlich des benannten Preises (keine Inhalts-Ähnlichkeit). Die Behauptung „kein Lastenheft-Bump nötig" trägt — das unveränderte Akzeptanzkriterium sagt den `--range`-Fall bereits namentlich zu. Die einzige Abweichung zwischen Plan und Diff sind zwei zusätzliche, im Commit offengelegte Kommentar-Korrekturen ohne Verhaltensänderung — keine der fünf `§3`-Ausschluss-Grenzen ist verletzt. `make gates` ist mit echtem Exit 0 gegengeprüft.
