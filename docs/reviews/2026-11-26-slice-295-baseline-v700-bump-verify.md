# Verifikation — slice-295: Baseline-Pin auf v7.0.0, mechanisch

- **Rolle:** Verifier (Modul 11) — DoD, Plan und MR gegen den Stand
- **Gegenstand:** `docs/plan/planning/in-progress/slice-295-baseline-v700-bump.md` §1, §2, §3 (samt beiden Plan-Änderungen), §4, §6; [MR-081](../../harness/conventions/MR-081-baseline-v700.md); Architect-Verdikt zu welle-98 samt Nachtrag V1; Review-Report [r1](2026-11-26-slice-295-baseline-v700-bump-r1.md) nur als Kontext
- **Range:** `738d2bca..HEAD` (Spitze `0985933f`, 12 Commits); Arbeitsbaum sauber
- **Datum:** 2026-11-26
- **Wegwerf:** alles im Scratchpad (`…/scratchpad/verify295/`), nichts im Arbeitsbaum; Mutationen in einem Klon

## Messungen (Kommando und Ergebnis, selbst gefahren, ohne Pipe auf den Exit-Code)

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `make baseline-verify` | Exit 0; `fetch-baseline-cache: verify ok (55 Dateien, vollständig)` |
| M2 | `make baseline-freshness` | Exit 0; `check-latest OK (Currency) — Pin v7.0.0 ist der neueste Release-Tag.` / `OK (Content) — gepinnter Tag v7.0.0 upstream unverändert (Bytes == vendored SHA256SUMS).` |
| M3 | unabhängig: `gh release download v7.0.0 -p lab-regelwerk.zip`, `unzip`, `diff -rq` gegen `.harness/baseline/v7.0.0` | einziger Unterschied `Nur in .harness/baseline/v7.0.0: SHA256SUMS`; zusätzlich Scratchpad-Bundle des Implementers `diff -rq` gleich |
| M4 | `make doc-check` | Exit 0; `d-check: 1319 Datei(en) geprüft, 0 Befund(e)` |
| M5 | `make adr-check RANGE=738d2bca..HEAD` | Exit 0; `1313 Datei(en) geprüft, 0 Befund(e)` |
| M6 | `make trace-check RANGE=738d2bca..HEAD` | Exit 0; `Range … aufgeloest, 12 Commit(s) — OK.` / `0 Befund(e)` |
| M7 | `make planning-check` | Exit 0; `0 Befund(e)` |
| M8 | `make archive-wave-test` | Exit 0 |
| M9 | `make gates` (Ausgabe in Datei, Exit-Code erfasst) | `EXIT=0`; letzte Zeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%`; semgrep `0 findings`; Gate `verify ok (55 Dateien, vollständig)` |
| M10 | `make verify-closure-notes` | `0 Befund(e)` (Slice liegt noch in `in-progress/`; kein Aussagewert für §7) |

## Nachzählung MR-081 (jede Zahl mit ihrem Kommando oder dem Verfahren des Eintrags)

Alter Baum aus `git archive 738d2bca .harness/baseline/v6.18.1`, neuer Baum der vendorte.

| Behauptung | Nachzählung |
|---|---|
| 54 / 55 Bundle-Dateien, vendort 55 / 56; 47 gleichnamig, 7 umbenannt, 1 neu | 54 / 55 (`find` über `regelwerk/`, `templates/`), `git ls-tree` alt 55, `git ls-files` neu 56, Manifest 55 Zeilen; `comm`: 47 / 7 / 1 (die sieben Namenspaare und `environment-check.template.md` wie genannt) |
| 38 byte-verschieden, 31 inhaltlich (19 Regelwerk, 12 Vorlagen), 7 nur Stempel, 862 / 804 Diff-Zeilen | `cmp` 38; normalisiert (ohne `<!-- Quelle:`, Stempel auf Platzhalter) 31 mit Diff, 19 + 12, 7 mit 0 (die im Eintrag benannten sieben); roh 862, normalisiert 804; Spitzen 258 / 74 / 57 / 49 / 35 stimmen |
| 77 Dateien / 190 Vorkommen / 175 Zeilen am Vorzustand | `git grep -l` / `-o` auf `738d2bca` ohne `.harness/baseline`: 77 / 190 / 175 |
| Frozen 17 / 51 (8/29 `done/`, 8/21 Reviews, 1/1 CHANGELOG); `conventions/done/` und CR 0 | gleich; `conventions/done` und `docs/plan/cr` 0 Dateien |
| 14 Dateien / 40 Vorkommen byte-stabil (7/28, 6/11, 1/1), drei mit je einer geänderten Zeile | Schleife des Eintrags: 14 / 40; `done/` 28, Reviews 11; die drei: `slice-274`, `slice-277`-Verifikation, Verdikt welle-93, je 1 Zeile, nur Link-Ziel (`git diff -U0`) |
| ADRs 5 / 5, Pläne der Welle 5 / 25, `MR-080` 1 / 9, lebend retargetet 49 / 100 | 5; 5 / 25 (`-o`; `-c` zählt 24 Zeilen — der Eintrag zählt Vorkommen, stimmt); 9; 49 / 100 |
| Treffermenge nach dem Bump: 10 Dateien | Kommando des Eintrags: `.d-check.yml`, `MR-081`, ADR 0112–0116, `slice-295`, `slice-296`, `welle-98` = 10; 0 Treffer in der Form `6.18.1` ohne `v` außerhalb |
| 7 Symlinks lösen gegen `v7.0.0` auf | `readlink`-Schleife `grep -c 'v6\.18\.1'` = 0; alle 9 Einträge unter `.claude/rules/` lösen auf (`[ -e ]`), Modus 120000 |
| Ventil 35 statt 30 Einträge, 13 statt 12 Tags | `grep -c` 35 / 30; Tags 13 / 12 (am Vorzustand mit `git show 738d2bca:.d-check.yml`); Diff `.d-check.yml` = ein Kommentarblock plus fünf Einträge (ADR 0112–0116) |
| „alt gegen neu: eine abweichende Zeile" (`modul-04` §Ziel-Form: ADR (MADR)) | Abschnitt je 14 Zeilen, `diff`: genau Zeile 3, `NNNN-titel.template.md` → `NNNN-title.template.md` |
| `slice-274` / Verdikt welle-93 wortgleich (`modul-11` §Bewusstes Brechen 18, `modul-05` §Ziel-Form: Slice 58 Zeilen) | `diff` je Abschnitt leer; 18 bzw. 58 Zeilen |
| 15 lebende Spannen, 7 verschoben, 8 unverschoben; plus 8 in den Plänen | 15 + 8 = 23 Direktiven auf Baum, alle `v7.0.0`; die 7 (+ 2 Plan-Bereiche) alt gegen neu Zeile für Zeile wortgleich; die 8 unverschobenen Bereiche in beiden Bäumen gleich |
| 9 Zeilen in 4 Dateien (Vorlagen-Namen) | Kommando des Eintrags am Vorzustand 9; Diff `stub.go` nur zwei Kommentarzeilen; **nach** dem Nachzug nicht 0, siehe V-1 |
| `slice.template.md` §1 bis §8 unverändert | `grep '^## '` beider Bäume gleich; `modul-06` §Roadmap-Struktur unterscheidet sich nur im Namen `wave-<Kennung>-results.md`, die Aussage zu *Offene Wellen* im alten Wortlaut |
| Zitat-Delta 36 / 35 / 1, weitere 381 | siehe unten — reproduziert |
| Release-Liste (`v6.17.0`, `v6.18.0`, `v6.18.1`, `v7.0.0`) | `gh release list`: dieselbe Folge, keine Zwischen-Tags |

**Retarget nur Tag-Tausch:** die 51 geänderten lebenden Textdateien (ohne Baum, Reviews, `conventions/done`) mit zurückgetauschtem Tag gegen `738d2bca` verglichen: Rest nur Pläne dieser Welle (Spannen, Plan-Änderungen), `.d-check.yml` (Ventil), die vier Dateien mit Vorlagen-Namen (`close-welle.md` 8, `plan-welle.md` 4, `stub.go` 4, `conventions.md` Z. 59), sechs Spannen-Zeilen in `reviewer.md`, je zwei in `MR-031/035/049` (Spannen), `baseline-verify.md` (54 → 55), Index-Zeilen in `harness/conventions.md`. `AGENTS.md`, `harness/README.md`, `harness/rules/`, `roadmap.md`, `observations/README.md`, `spec/*`: reiner Tag-Tausch.

**Zitat-Delta (besonderer Punkt).** Verfahren des Eintrags nachgefahren: `git ls-files` mit den genannten Ausschlüssen, Symlinks ausgenommen (630 Dateien), `tr`-Normalisierung, `grep -o '„[^„"“]*["“]'` (990 Treffer), je Zitat `grep -qF` in beiden normalisierten Bäumen. Mit der **Länge des Inneren in Zeichen ≥ 25** (`LC_ALL=C.UTF-8`, `wc -m`): 417 Zitate, **36** in mindestens einem Baum, **35** in beiden, **1** nur im alten (`„Zwei Formen des Reproduzierbarkeits-Ankers“`, `slice-295`), **381** in keinem — exakt die Zahlen des Eintrags. Die Differenz zum Reviewer (31 / 30 / 1) ist **fünf** Fundstellen; fünf Treffer liegen in `spec/` (3) und `docs/user/benutzerhandbuch.md` (2) — genau diese Menge erklärt die Differenz rechnerisch, ich kann den Scope des Reviewers aber nicht bestätigen; sie ist als erklärt-nicht-bestätigt zu lesen (V-5). Der Reviewer-Befund F-2 galt der früheren Zahl 51 / 50 ohne Kommando; die jetzige Fassung trägt das Verfahren.

## Bewusstes Brechen (DoD-Behauptungen, die sich auf Gate-Wirkung stützen)

Klon `…/verify295/clone` (HEAD `0985933f`), `d-check:latest` mit `--network none`; sauberer Klon: `1319 Datei(en) geprüft, 0 Befund(e)`.

| Mutation | Erwartung | Beobachtung |
|---|---|---|
| `d-check:cite`-Spanne in `reviewer.md` von `modul-10…:115-115` auf `114-114` | Zitat-Fäule | 1 Befund `citation-mismatch … Zitattext ist kein zusammenhängender Teilstring der Quell-Spanne` — aus dem richtigen Grund |
| dieselbe Direktive auf `…/v6.18.1/…` | Ziel fehlt | 1 Befund `citation-out-of-range … d-check:cite-Ziel fehlt` |
| die fünf Ventil-Einträge entfernt | Link-Bruch | 5 Befunde `target-missing` auf `…v6.18.1/regelwerk/modul-04-adrs.md` in ADR 0112 bis 0116, sonst nichts — das Ventil trägt genau diese fünf |
| ein Symlink unter `.claude/rules/` auf den alten Baum | `doc-check` blind, `baseline-verify` rot | `make doc-check`: 0 Befunde (erwartet, `git grep`/Link-Module lesen keinen Symlink-Inhalt); `make baseline-verify`: `toter Symlink .claude/rules/modul-11-verification.md — Ziel fehlt (Baseline-Bump?)`, Exit 2 — der Wächter ist `baseline-verify` ([MR-055](../../harness/conventions.md#mr-055)) |
| Inline-Code `harness/nicht/vorhanden/….md` einzeilig / mehrzeilig in `MR-081` | `codepath-missing` | einzeilig: Befund; **mehrzeilig (Pfad auf Öffnungs- oder Folgezeile des Spans): 0 Befunde** — siehe V-3 |

## DoD-Punkte

- **Vendoring — erfüllt.** M1–M3: 56 Dateien (`git ls-files`), 55 im Manifest, byte-gleich zum Upstream-Asset (unabhängig geladen); `ls .harness/baseline` nennt nur `v7.0.0`; der Baum `v6.18.1` (55 Dateien in `738d2bca`) ist entfernt; Pin in `harness/conventions.md` Z. 27, 41, 46, 47 (§Baseline, §Adoptierte Konventions-Quellen), `Index`-Zeile `MR-081` unter *Aktive*. Dass der Baum über `fetch-baseline-cache.sh v7.0.0` entstand, ist am Ergebnis nicht unterscheidbar von einem Kopieren; belegt ist der Bytebestand, nicht der Weg.
- **Spiegel — erfüllt.** 10 verbleibende Treffer außerhalb der Frozen-Klassen und des Baums (alle benannt, Stand `git grep -l`); 7 Symlinks lösen auf; 15 + 8 Spannen auf `v7.0.0`, die 7 verschobenen gelesen wortgleich; Vorlagen-Namen: die 9 Zeilen tragen die neuen Namen; `make doc-check` 0 Befunde (M4). Mutationsprobe: Spannen- und Ventil-Wirkung aus dem richtigen Grund rot.
- **Einträge — erfüllt, mit INFO.** `MR-081` trägt Delta nach Datei (Zahlen oben nachgezählt), Frozen-Klassen am Vorzustand, Zitat-Delta mit Verfahren; Pflichtfelder wie `MR-074`/`MR-080` und `MR-NNN-title.template.md` des neuen Baums (Datum, Geltungsbereich, Ersetzt-Baseline-Regel `—`, Adaption, Begründung, Auflösungs-Trigger, Löst auf, Ausgelöst durch). `MR-080` per `d67b6221` `R100` (`rename harness/conventions/{ => done}/MR-080-baseline-v6181.md (100%)`), Link-Tiefen in `3d8f3555`; Index: `MR-081` aktiv, `MR-080` aufgelöst. „Jede Bewegung als *nicht übernommen* markiert": Sammelüberschrift (Z. 30) plus „Nicht übernommen" bei Bewegung 2, 3, 4; Bewegung 1 nennt die Zuständigkeit (`slice-297`/`slice-298`), 5 und 6 nur die Sammelüberschrift — erfüllt in der Sammelform, Reviewer F-4 steht unverändert (V-6).
- **`make gates` grün; `make archive-wave-test` grün — erfüllt** (M9, M8). Kommentar-Nachzug in `stub.go`: zwei Zeilen, nur Dateinamen der Vorlagen, kein Verhalten.
- **Review durchgeführt — erfüllt, mit Vorbehalt.** Report `2026-11-26-slice-295-baseline-v700-bump-r1.md` liegt vor (1 MEDIUM, 2 LOW, 3 INFO); F-1, F-2, F-3 sind in den Fixrunden-Commits `5cda1d2d`…`0985933f` eingearbeitet (Ventil fünf statt sieben, Kommando zum Zitat-Delta, `55 Dateien`). Die Fixrunde war **nicht rein dokumentarisch** (V-2) und hat keine zweite Review-Runde.
- **Closure-Notiz mit Lerneintrag — offen.** §7 trägt `*(Wird bei der Closure gefüllt.)*`; Sache des Planners nach dieser Verifikation.
- **Beobachtungs-Register fortgeschrieben — offen.** Keine Evidenz-Datei für `slice-295` unter `docs/plan/planning/observations/` (`git ls-files … | grep 295` leer); `observations/README.md` nur Pin-Retarget. Kandidaten für Evidenz: V-3 (siehe dort) und V-1; „keine Beobachtung" wäre ebenfalls eine in §7 zu notierende Antwort.
- **Jedes Risiko aus §6 trägt einen Ausgang — offen (formal), Inhalt gemessen.** Noch keine Ausgänge eingetragen. Nach den Messungen oben: Risiko 1 *entfallen* (doc-check meldet nur die benannten Folgen, 0 Befunde); Risiko 2 *entfallen* mit der benannten Ausnahme (`git diff --stat` über die Frozen-Verzeichnisse nennt `slice-274`, Verdikt welle-93, `slice-277`-Verifikation je 1 Zeile, `MR-080` unter `conventions/done/` und die in dieser Welle neu angelegten Reviews); Risiko 3 *entfallen* (alte und neue Zeile je Spanne genannt, Wortlaut gleich — Sammelform „je reine Zeilenverschiebung bei gelesenem, wortgleichem Text"); Risiko 4 *entfallen unter Vorbehalt V-1* (alle Zahlen reproduzieren, eine Kommando-Behauptung „danach 0" nicht wörtlich); Risiko 5 (kopierte Welle-Vorlage trägt `wave-<Kennung>-results.md`) *weiter offen* → `slice-298`, belegt: `wave-results.template.md` Z. 4 nennt `done/<wave-id>-results.md`, `plan-welle.md` Z. 73–77 verlangt die Lifecycle-Note.
- **Drei Paarungen (Anker · Folge-Slice · Register) — offen, nicht Gegenstand dieser Verifikation** (trägt die nächste Welle-Closure).

## Abgrenzung (nichts aus slice-296/297/298)

Diff `738d2bca..HEAD` außerhalb Baum, Plänen, Reviews, `harness/`, `.claude/`, `.harness/skills/`: nur `spec/architecture.md` und `spec/spezifikation.md` (je die **Rolle**-Zeile, Pin), `tools/archive-wave/stub.go` (zwei Kommentarzeilen), `.d-check.yml`. `internal/`, `Makefile`, `Dockerfile`, `docs/user`, READMEs, `CHANGELOG.md`, `spec/lastenheft.md`, `tools/harness/` unberührt (`git diff --stat` leer). `spec/spezifikation.md` weder umbenannt noch sonst inhaltlich berührt; keine `wave-`-Konfiguration; `MR-081` übernimmt keine Neuerung. Die vier Pläne tragen nur Spannen-Neuankerung und die zwei Plan-Änderungs-Texte (`welle-98` Z. 213/324–325 auf fünf).

## Befunde

Keine HIGH, kein MEDIUM. Zwei LOW, vier INFO.

### V-1 — LOW — „danach 0" ist als Kommando nicht wörtlich reproduzierbar

- **Stelle:** `harness/conventions/MR-081-baseline-v700.md` Z. 139–141, „Befund `git grep -n "welle\.template|…" -- . ':!.harness/baseline' ':!docs/plan/planning' ':!docs/reviews' ':!harness/conventions/done'` war vor dem Nachzug 9 Zeilen, danach 0".
- **Messung:** `738d2bca`: 9 Zeilen (stimmt). Heute: 7 Zeilen, alle in `MR-081` selbst (Z. 128–133 und 215, die die Namen als Vergangenheits-Aussage nennen); außerhalb von `MR-081` 0. Die Ausschlussliste des Kommandos nimmt `MR-081` nicht aus; wer es wörtlich fährt, sieht 7, nicht 0 (§6 Risiko 4: Nachzählen „mit dem im Eintrag genannten Kommando").
- **Folge:** Zahl stimmt für den Gegenstand (Verweise des Repos auf die alten Namen), das genannte Kommando trägt sie nicht. Fix: `':!harness/conventions/MR-081-baseline-v700.md'` ins Kommando oder der Satz „außerhalb dieses Eintrags".

### V-2 — LOW — Die Fixrunde nach r1 änderte Gate-Konfiguration und wurde nicht zum zweiten Mal reviewt

- **Stelle:** `git diff --stat ef442d8c..HEAD`: `.d-check.yml` (Ventil sieben → fünf, `17 ++---`), zwei eingefrorene Dateien (je 1 Zeile), `MR-081` (86), Pläne, `baseline-verify.md`.
- **Befund:** Nicht „rein dokumentarisch": `.d-check.yml` ist die Konfiguration des `links`-Gates und die Fixrunde ändert, was es still durchlässt (zwei Dateien weniger im Ventil). Der Plan sieht höchstens zwei Runden vor; es gab eine plus Fixrunde ohne zweiten Report. Ich habe das Delta statt dessen selbst gemessen: `make doc-check` 0 Befunde; die fünf Einträge sind tragend (Mutation oben: ohne sie 5 `target-missing`); die beiden Retargets verändern je eine Zeile, die zitierten Abschnitte sind gelesen wortgleich (`diff` leer, 18 und 58 Zeilen); `adr-check` 0 Befunde. Kein Befund am Inhalt — nur die Lücke im Verfahren ist zu benennen.

### V-3 — INFO — `codepaths` sieht Pfade in mehrzeiligen Inline-Code-Spans nicht (Antwort auf die Frage zu `MR-081`)

- **Beobachtung:** `MR-081` trägt acht mehrzeilig geöffnete Inline-Code-Spans (Z. 150–153, 164–167, 198–200, 265–266). `make doc-check` ist trotzdem grün. Zwei verschiedene Module, zwei verschiedene Antworten:
  - `spans`: **Form-Eigenheit, spec-gedeckt** — `DC-FA-SPAN-001` Happy Path: „auch solchen, die einen Zeilenumbruch enthalten". Balancierte mehrzeilige Spans sind kein Befund.
  - `codepaths`: **Prüflücke** — empirisch im Klon: ein nicht existierender Pfad (`harness/nicht/vorhanden/….md`) in einem einzeiligen Span meldet `codepath-missing`; derselbe Pfad in einem mehrzeiligen Span (auf der Öffnungs- wie auf der Folgezeile) meldet **nichts**. Die Prüfung ist zeilenlokal. Ein Satz dazu steht weder in `DC-FA-CODE-001` noch in `spec/spezifikation.md` (kein Treffer für mehrzeilig bei Inline-Code-Pfaden).
- **Wirkung auf diesen Slice:** keine. Alle Pfade in den mehrzeiligen Spans von `MR-081` (`.harness/baseline`, `docs/plan/planning/done`, `docs/reviews`, `CHANGELOG.md`, `harness/conventions/done`) existieren (mit M4 und von Hand gelesen). Es ist ein Beobachtungs-Kandidat für das Register (Modul verspricht nur über das, was es scannt — [`AGENTS.md`](../../AGENTS.md) §3.8; nicht in diesem Slice zu beheben).

### V-4 — INFO — Zwei Zwischenstände sind einzeln rot (Push-Pflicht, Risiko der Zerlegung)

- **Messung** (`d-check:latest` auf je ausgechecktem Commit im Klon): `79a14c01` grün (0 Befunde); `04dca313` **1 Befund** (`harness/conventions.md:186 conventions/done/MR-080-baseline-v6181.md target-missing`); `d67b6221` rot (`MR-080` unter `done/`: `../conventions.md#mr-069`, `#mr-039`, `#mr-074`, … `target-missing`); `3d8f3555` und alle späteren grün.
- **Folge:** erwartet nach [`AGENTS.md`](../../AGENTS.md) §3.3 (Index-Zeile und Link-Tiefen im Arbeitsbaum vorab; der `pre-commit`-Hook prüft den Arbeitsbaum). `main` steht 18 Commits vor `origin/main`, nichts gepusht; `04dca313`…`d67b6221` dürfen nicht die Spitze eines Push werden — alle Commits bis mindestens `3d8f3555`, sinnvoll bis `0985933f`, in **einem** Push. Die Botschaft von `04dca313` nennt noch „sieben ignore-refs"; sie ist Historie und durch Plan-Änderung (`c7713f02`), `MR-081` und den Verdikt-Nachtrag überholt, nicht umzuschreiben.

### V-5 — INFO — Zitat-Delta: Verfahren nennt den Bezug der 25 Zeichen nicht

- Der Eintrag sagt „ab 25 Zeichen", nicht ob das Innere oder das Zitat samt Anführungszeichen, und nicht in Zeichen oder Bytes. Mit dem Inneren in Zeichen reproduziert sich 36 / 35 / 1 / 381 exakt. Die abweichende Zählung des Reviewers (31 / 30 / 1) ist um fünf Fundstellen kleiner; rechnerisch passen die fünf Treffer aus `spec/` und `docs/user/`. Ob der Reviewer diese Dateien ausnahm, ist offen und nicht entschieden; die Schlussfolgerung (kein lebendes Zitat verliert seinen Wortlaut außer der einen im Plan genannten Überschrift) trägt in beiden Zählungen.

### V-6 — INFO — Reviewer F-4, F-5, F-6 unverändert

- F-4 (kein eigenes „nicht übernommen" bei Bewegung 1, 5, 6): offen, DoD in Sammelform erfüllt. F-5 (Baum-Zuwachs in eigenem Commit statt „ein Hauptcommit"): `79a14c01` ist einzeln grün (oben gemessen), der Zwischenstand schadet nicht. F-6 (Inline-Code-Nennungen des alten `MR-080`-Pfads in drei eingefrorenen Reports): `codepaths` nimmt `docs/reviews/**` aus; bewusst, im Plan benannt.

## Offene Risiken am Handoff

1. **Push als eine Einheit** (V-4); nach Push `make ci`/CI prüfen, nicht aus dem lokalen Stand schließen.
2. **Zwei offene DoD-Punkte, ein formaler** (Closure-Notiz, Register, Risiko-Ausgänge) — vom Planner zu füllen; V-1 sollte vorher ein Satz im Eintrag bekommen, weil Risiko 4 sonst nicht ohne Vorbehalt *entfallen* lauten kann. Eintrags-Änderung wäre ein weiterer Commit nach der Verifikation.
3. **Welle-Closure-Trigger** (Gegenmessung am Endstand, `blackbox-probe`, `fullbuild`) liegen bewusst außerhalb dieses Slice; die einzige Berührung des Produkts ist die Pin-Zeile in `spec/spezifikation.md` (eingebettet), `make test` ist grün.
4. **Gates decken den Spiegel nur teilweise** ([`BEO-ALL/pin-bump-mirrors-ungated`](../plan/planning/observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)): Prosa-Pins in `roadmap.md` Z. 13 und `observations/README.md` Z. 18 hält kein Gate; sie sind von Hand gelesen (Aussage im neuen Baum gleich).

## Urteil

Plan und Diff stimmen; jede nachgezählte Zahl des Eintrags reproduziert sich bis auf die Wörtlichkeit eines Kommandos (V-1). `make gates` Exit 0, die fünf Einzel-Gates Exit 0, der Bytebestand des Baums gegen Upstream gleich. DoD: Vendoring, Spiegel, Einträge, Gates, Review **erfüllt**; Closure-Notiz, Register, Risiko-Ausgänge, Paarungen **offen** und Sache des Planners bzw. der Welle-Closure. Keine Plan-Abweichung, die nicht im Plan als Plan-Änderung vor dem Code steht.
