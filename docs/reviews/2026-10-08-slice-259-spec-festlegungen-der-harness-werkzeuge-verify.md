# Verifikation slice-259 — Festlegungen der Harness-Werkzeuge in der Spezifikation

**Rolle:** Verifier („Bauen wir es richtig?" — gegen Plan und DoD, nicht gegen den Diff-Stil).
**Gegenstand:** `docs/plan/planning/in-progress/slice-259-spec-festlegungen-der-harness-werkzeuge.md`, Stand `94157a96` (sauberer Baum).
**Datum:** 2026-10-08.

## Summary

Verdikt: **konform** für die DoD-Punkte 1 bis 3. Punkt 4 ist mit diesem Bericht erfüllt (die Reviews R1 und R2 liegen vor). Punkt 5 (Closure) ist erwartungsgemäß offen. Ich habe jede der vier Festlegungen `SPEC-089` bis `SPEC-092` gegen Code und Konfiguration gehalten und ihre korrektheitskritischen Randformen bewusst gebrochen: Jede wurde aus dem behaupteten Grund rot. Es gibt keine DoD-Verletzung, nur einen LOW-Prozessbefund und drei INFO.

## Sensor-Belege (selbst gefahren)

| Sensor | Ergebnis |
|---|---|
| `make gates` auf `94157a96` | `EXIT=0`, `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`. Einzelbelege: `verify ok (54 Dateien, vollständig)`; `d-check: 1004`/`1010 Datei(en) geprüft, 0 Befund(e)`; `0 issues.`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; `Ran 55 rules on 65 files: 0 findings.` |

## DoD-Abgleich

| DoD-Punkt | Befund | Beleg |
|---|---|---|
| 1. §7 nach der Vorlage, Historie §8, Kopplungen nachgezogen, Nachtrag zu MR-0098 | erfüllt | §7 steht zwischen §6 und §8 und trägt die Spaltenform der Vorlage `templates/spec/spezifikation.template.md` §7 (die Kopfzeile heißt im Haus-Stil `Kennung`, wie in §6). `.d-check.yml`: `exclude-sections: [Geschichte, "7. Historie", "8. Historie"]`. Die `structure`-Regel der Spezifikation zeigt auf `## 8. Historie`; die Lastenheft-Regel bleibt bei `## 7. Historie`, wie es sein muss. `harness/conventions/MR-075-…md` ist nach der Vorlage `MR-NNN-titel.template.md` gebaut, die Index-Zeile steht in `harness/conventions.md`. `MR-0098` ist nicht überschrieben. |
| 2. Vier Einträge, am Code geprüft | erfüllt | siehe §Festlegungen gegen Code und §Bewusstes Brechen |
| 3. Sensor-Dateien bzw. Index-Zeile verlinken die Kennung; `make gates` grün | erfüllt | `lint.md`, `semgrep.md` und `baseline-verify.md` verlinken `SPEC-090`, `SPEC-091` und `SPEC-092`. In `harness/README.md` nennt die Zeile von `coverage-gate` `SPEC-089`, die Zahl 93 steht dort nicht mehr. Die Schwelle als Zahl steht außerhalb der Spezifikation nur noch in ihren Trägern (`Makefile` `THRESHOLD ?= 93`, `Dockerfile` `ARG COVERAGE_THRESHOLD=93`), wie Plan §3 es vorsieht. Gates grün (oben). |
| 4. Review + Verifikation | erfüllt | `docs/reviews/2026-10-08-slice-259-…-r1.md`, `…-r2.md`, dazu dieser Bericht |
| 5. Closure | offen (erwartet) | §7 ist leer, beide Risiken in §6 stehen auf `(offen)`. Das ist vor dem `git mv` nach `done/` zu leisten. |

## Festlegungen gegen Code

- **SPEC-089** (`Dockerfile` Stage `coverage`, `tools/coverage-gate.sh`, `Makefile`): `-coverpkg` über `go list ./internal/...`, `-covermode=atomic`, `./...` des Hauptmoduls, Prüfgröße `^total:`. Die `&&`-Kette sorgt dafür, dass ein roter Test vor der Messung abbricht. Der Vergleich ist `p+0 >= t+0`. Die Prüfung auf eine Zahl (Regex `^[0-9]+(\.[0-9]+)?$`) endet mit Exit 2. Gesetzt sind `LC_ALL=C`, `THRESHOLD ?= 93` und `ARG COVERAGE_THRESHOLD=93`. Alles stimmt mit der Festlegung überein.
- **SPEC-090** (`.golangci.yml`): Gezählt ergeben sich 29 `enable`-Einträge, `default: none`, und die Schwellen cyclop 15, gocyclo 15, gocognit 20, funlen 100/60, nestif 5, dupl 150, maintidx 20, interfacebloat 10. Es gibt sechs `exclusions.rules`, die genau so geschnitten sind wie beschrieben. Dazu kommen errcheck `fmt.Fprint*`, die ireturn-Allowlist inklusive `port/driven` und go-billy sowie `generated: lax`. Für nolintlint sind die drei Schalter scharf gesetzt. Abweichung: siehe INFO-2.
- **SPEC-091** (`tools/semgrep.sh`): Image mit `@sha256:`-Pin, `RULES_COMMIT`-Pin, `RULES_SUBSET="go/lang/security"`, `--network none`, `--error`, Abbruch mit Exit 2 ohne Zeile `Ran [1-9]… rules`. Stimmt überein.
- **SPEC-092** (`tools/harness/fetch-baseline-cache.sh --verify`): Prüft auf `sha256sum`, `find` und `readlink` und darauf, dass das Manifest vorhanden ist. Danach folgen `sha256sum -c`, die Zählung über das ganze Tag-Verzeichnis ohne Manifest mit `-gt 0` und `-eq` sowie `check_aliases`, rekursiv über `find -type l`. Stimmt überein.

## Bewusstes Brechen

Jede Probe wurde zurückgesetzt. `git status --short` ist danach leer, abgesehen von diesem Bericht.

1. **Matrix-Kopplung (Risiko 2, DoD 1):** Ich habe `"8. Historie"` aus `exclude-sections` entfernt und `make doc-check` gefahren. Ergebnis: rc 2, `4 Befund(e)`, alle `matrix-forbidden` in `spec/spezifikation.md:3663/3690/3768`, also in §8 (spec-straten → aussen/adaptionsblock). Nach dem Zurücksetzen ist der Lauf grün.
2. **structure-Kopplung:** Den Selektor der Spezifikations-Regel habe ich auf `## 7. Historie` zurückgesetzt. Ergebnis: rc 2, `section-missing` („kein Abschnitt passt auf den Selektor"). Als Gegenprobe mit korrektem Selektor habe ich zwei Zeilen der §8-Historie vertauscht. Ergebnis: `section-unordered` bei Zeile 3645. Die Regel greift also wirklich auf §8.
3. **SPEC-089, alte gegen neue Skript-Fassung** (`d39c4b04` gegen HEAD, im Image `d-check:coverage`):

   | Eingabe | alt | neu |
   |---|---|---|
   | total 93.0, Schwelle 93 / 93.0 | 0 / 0 | 0 / 0 |
   | total 92.9, Schwelle 93 | 1 | 1 |
   | Schwelle `""` / `abc` / `-1` | **0 / 0 / 0** | 2 / 2 / 2 |
   | Eingabe fehlt / keine `total:` | 2 / 2 | 2 / 2 |
   | `total:` nicht parsbar | **1** (toter Exit-2-Zweig) | 2 |
   | Host, `LC_ALL=de_DE.UTF-8`, total 93.0 | **1** (`printf: 93.0: Ungültige Zahl`) | 0 |

   Die alte Fassung ist also genau aus den Gründen rot bzw. still grün, die der Fix behauptet. Über make habe ich zusätzlich `make coverage-gate THRESHOLD=abc` gefahren: rc 2, Meldung `Schwelle 'abc' ist keine nicht negative Zahl`.
4. **SPEC-091:**
   - Eine *nicht getrackte* Datei `internal/zzprobe/probe.go` mit `md5.Sum` ergibt `use-of-md5`, `Ran 55 rules on 66 files: 1 finding`. Das Skript endet mit rc 1, `make semgrep` mit rc 2.
   - Dieselbe Datei als `probe_test.go` ergibt 0 findings bei 65 files, die `.semgrepignore`-Zählung steigt von 79 auf 80.
   - Dieselbe Datei im `.gitignore`-Pfad `.harness/cache/` ergibt 0 findings bei 65 files, sie wurde also nicht gescannt.
   - Ein leerer Regel-Umfang (`SEMGREP_RULES_CACHE` auf ein leeres `go/lang/security`) ergibt rc 2 mit `0 Regeln geladen … breche ab statt still grün`.
   - Fazit: Der Umfang „nicht ignoriert, getrackt oder nicht, ohne Voreinstellung" und beide Exit-Codes sind bestätigt.
5. **SPEC-092:**
   - Ein Geschwister `zz-sibling.md` im Tag-Verzeichnis ergibt rc 1 mit `Manifest (54 Zeilen) != Dateien auf Platte (55)`.
   - Ein toter Symlink `.claude/rules/.zz/.dead.md` (Unterverzeichnis und Punkt-Name) ergibt rc 1 mit `toter Symlink …`.
   - Eine geänderte Datei im Baum ergibt rot. Ein fehlendes `SHA256SUMS` ergibt rot.
   - Über make endet jeder dieser Fälle mit rc 2.

`SPEC-090` habe ich nicht gebrochen (der Guard verweigert den direkten Aufruf von `golangci-lint` im Container). Die Zählungen habe ich mit `awk` am Konfig-Text bestätigt: 29 Linter, 6 Ausschluss-Regeln. Der Punkt ist nicht korrektheitskritisch im Sinne von Modul 11, denn er behauptet keinen Test.

## Befunde

- **LOW-1 — Mitnahme nicht vor dem Code im Plan.** Der Absatz „Plan-Änderung nach R1, vor dem Code" stand in `26d63773` ohne den Satz zu `LC_ALL=C`. Dieser Satz („Mitgenommen: `LC_ALL=C` …") kam erst mit dem Code-Commit `7209e90f` in den Plan. AGENTS.md §6 Schritt 4 verlangt, dass eine Mitnahme vor dem Code im Plan steht. Inhaltlich ist die Änderung richtig und bruchgetestet (Probe 3). Der Befund gilt nur der Reihenfolge.
- **INFO-2 — SPEC-090 ordnet `generated: lax` als „Linter-Einstellung" ein.** In der Konfiguration steht der Schalter unter `linters.exclusions`, neben den sechs Pfad-Regeln, nicht unter `settings`. Anzahl und Wirkung stimmen, nur die Zuordnung zur zweiten Form ist ungenau.
- **INFO-3 — Risiko 1 nennt einen Träger zu wenig.** Die Schwelle hat außerhalb der Spezifikation zwei Träger: `Makefile` `THRESHOLD ?= 93`, der über make wirkt, und `Dockerfile` `ARG COVERAGE_THRESHOLD=93`. Das Risiko nennt nur `COVERAGE_THRESHOLD`. Das ist beim Risiko-Ausgang zur Closure mitzunehmen.
- **INFO-4 — Der Kopfkommentar der `.golangci.yml` kann der Zahl 29 widersprechen.** Er lautet „5 Default- + 23 SOLID-nahe … dazu nolintlint …, zusammen 24" und ist liest sich gegen die 29 aus SPEC-090 mehrdeutig. Er stammt aus dem Bestand und ist nicht Teil des Slice-Diffs. Er wäre ein Kandidat für den Folge-Slice.

## Abgrenzung eingehalten

Ich habe die Abgrenzung des Plans gegen den Diff gehalten:
- Es gibt keinen §7-Eintrag für Werkzeuge außerhalb von `gates` und keinen für Gates mit eigener DC-Anforderung.
- Es wurde keine ADR geändert und keine Lastenheft-Anforderung angelegt.
- Den Produkt-Code (`internal/hexagon/core/app/suggest.go`, der noch `"7. Historie"` vorschlägt) hat der Slice nicht angefasst. Dafür ist slice-261 in `open/` geschnitten.
- Der Move-Commit `d39c4b04` ist rein (0 Zeilen).
- Die Plan-Änderungen nach R1 und R2 stehen jeweils vor ihrem Code-Commit, mit Ausnahme von LOW-1.
