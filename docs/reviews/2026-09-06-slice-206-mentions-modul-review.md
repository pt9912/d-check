# Review-Report — slice-206 (Modul `mentions`, Implementierung)

- **Review-Art:** Code — geprüft wird die Implementierung gegen Plan, Anforderung, Entscheid und Hard Rules; Maßstab sind [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) (zehn Akzeptanzkriterien), [ADR-0084](../plan/adr/0084-mentions-eigenes-modul.md) samt `## Geschichte`, [`spec/spezifikation.md`](../../spec/spezifikation.md) §`DC-FA-MENT-001.a`/§4 `SPEC-082`, [`AGENTS.md`](../../AGENTS.md) §3.1/§3.2/§3.4/§3.5/§3.7/§3.8/§4/§5/§6 und der Baseline-Kanon (`modul-11-verification.md` §Fitness Function ohne Standard-Tool, `modul-13-quality-gates.md` §Fitness Function / §Hard Rule)
- **Gegenstand:** `ac12993` (`feat(mentions)`, 18 Dateien, +734/−10) — der Hauptgegenstand; Kontext `d33c961` (Vorprüfungen) und `401ec1d` (Beanspruchung). Range `1f6933f..HEAD`
- **Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) @ 1.13.0 · **Modell-ID:** claude-opus-5[1m] · **Datum:** 2026-09-06
- **Eingangs-Kontext:** Slice-Plan [`slice-206`](../plan/planning/in-progress/slice-206-mentions-modul.md) (§3 Out-of-Scope, §4 DoD, §5 Risiken, §6 Rückführungen, §7 Vorprüfungen), `internal/hexagon/core/rules/mentions.go` + `mentions_test.go`, `internal/adapter/driven/configyaml/{configyaml.go,mentions_test.go}`, `internal/adapter/driven/report/report.go`, `internal/hexagon/core/{model/config.go,rules/run.go,app/diagnose.go}`, `.d-check.yml`, `Makefile`, `harness/README.md`, `harness/sensors/mention-coverage.md`, `docs/user/{benutzerhandbuch,operations}.md`, das Beobachtungs-Register (35 Verzeichnisse), der ADR-Index und die Sensors-Tabelle als Ist-Bestände
- **Vorherige Findings am gleichen Gegenstand:** [R1](2026-09-06-slice-205-mentions-anforderung-review.md) (9 MEDIUM/5 LOW/4 INFO) und [R2](2026-09-06-slice-205-mentions-anforderung-review-r2.md) (7 MEDIUM/4 LOW/4 INFO) zu slice-205. Dominante Klassen dort: `selbstauskunft-zahl-reproduziert-nicht` (sechs Läufe in Folge), `semantic-change-body-only-edges-stale`, `praemisse-gegen-ein-dokument-statt-gegen-die-menge`. **Alle drei erscheinen in diesem Lauf wieder** — die dritte als HIGH.

**Gefahrene Läufe.**
`make gates` → Exit 0; `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `d-check: 676 Datei(en) geprüft, 0 Befund(e)`; `coverage-gate: OK — Coverage 94.50% erfüllt Schwelle 93%`; `Ran 55 rules on 63 files: 0 findings.`
`make mention-coverage` → Exit 0; `d-check: mentions: 108 von 108 Artefakt(en) erwähnt, über 2 Dokument(e)` / `d-check: 676 Datei(en) geprüft, 0 Befund(e)`.
`--json` / `--yaml` mit `--enable mentions` → `summary.notes` trägt die Bezugsmenge (beide Formen einzeln gelesen). Default-Lauf ohne `mentions` → `summary` **ohne** `notes`-Schlüssel (byte-identisch).
Exit-2-Pfade einzeln gemessen: leere Soll-Menge, leere Ist-Menge, aktives Modul ohne Block — je `d-check: error: …`, Exit 2.
`make doc-check` mit diesem Report im Baum → `d-check: 677 Datei(en) geprüft, 0 Befund(e)`, Exit 0.
**Eigene Bruch-Tests, vier**, jeder mit anschließender Wiederherstellung; Arbeitsbaum am Ende sauber (`git status --short` leer), einzige Änderung ist diese Report-Datei:

| # | Eingriff | Erwartung | Ergebnis |
|---|---|---|---|
| A | ADR-0084-Zeile aus `docs/plan/adr/README.md` entfernt (`grep -c` danach 0) | Befund | **Exit 0, „108 von 108"** — H-1 |
| B | ADR-0001-Zeile entfernt (Kontrolle: nirgends sonst genannt) | Befund | Exit 1, `docs/plan/adr/0001-…md:1 … artifact-unmentioned`, „107 von 108" |
| C | Link auf `sensors/test.md` aus der Sensors-Tabelle in `harness/README.md` entfernt | Befund | **Exit 0, „108 von 108"** — H-2 |
| D | `sort.Strings(out)` in `mentionsResolve` entfernt | `make test` rot | **`make test` Exit 0** — M-3 |

B ist das *bewusste Brechen* aus `modul-13` §Fitness Function; A, C und D sind die Gegenrichtung aus `modul-11` §Fitness Function ohne Standard-Tool — der Bestand, auf dem der Sensor schweigen **müsste** und schweigt, obwohl der Defekt da ist.

---

## Akzeptanzkriterium → misst ein Test es wirklich?

[`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) führt **zehn** Akzeptanzkriterien (nicht neun, siehe L-1). Alle Tests liegen auf der Kern-Ebene (`CheckMentions`); es gibt **keinen** CLI-/Akzeptanztest für `mentions`.

| # | Kriterium | Test | Misst es wirklich? |
|---|---|---|---|
| 1 | **Happy Path** — Exit 0, kein Befund | `TestMentionsHappyPath` | **teilweise.** Befundfreiheit und Bezugsmenge gemessen; **Exit 0 nicht** — der Test kennt keinen Exit-Code. Manuell reproduziert (`make mention-coverage`) |
| 2 | **Bezugsmenge im sauberen Lauf**, in **beiden** Ausgabe-Formen | `TestMentionsNoteNenntDieBezugsmenge` | **nein.** Der Test ruft `MentionsResult.Note()` — den *Erzeuger* des Strings. Weder die stderr-Zeile (`report.Text`) noch `summary.notes` unter `--json`/`--yaml` wird von irgendeinem Test berührt (`grep -rn "Notes" --include=*_test.go internal/` liefert **null** Treffer). Genau die Hälfte, die das Kriterium ausdrücklich hinzufügt („Eine Zusammenfassung, die es nur auf stderr gäbe, wäre unter `--json` nicht vorhanden"), ist die ungemessene |
| 3 | **Boundary (Ist-Seite ist eine Menge)** | `TestMentionsIstMengeIstVereinigung` | **ja.** Vier statt acht Dokumente, aber die Eigenschaft ist gemessen: ein Vorkommen in einem von vieren genügt, `Documents == 4` geprüft |
| 4 | **Negative (Fund)** — Exit 1, `file`/`line` | `TestMentionsFindet` | **teilweise.** `File`, `Line: 1`, `Reason`, `Target` und die Bezugsmenge gemessen; **Exit 1 nicht**. Manuell reproduziert (Bruch-Test B) |
| 5 | **Erkennungsform** — `basename` vs. `path` | `TestMentionsErkennungsform` | **ja.** Dieselbe Eingabe, beide Richtungen, ein Befund gegen null |
| 6 | **Leere Soll-Menge, fail-closed** | `TestMentionsLeereSollMengeFailClosed` | **teilweise.** `error != nil` **und** die Nennung des Globs gemessen; **Exit 2 nicht**. Manuell reproduziert |
| 7 | **Leere Ist-Menge, fail-closed** | `TestMentionsLeereIstMengeFailClosed` | **teilweise, aber der beste Test der Datei.** Er schließt ausdrücklich aus, dass die *erste* Bedingung zuschlägt (`if strings.Contains(err.Error(), "artifacts")`) — genau der Fixture-Fehler, den `wortlaut-behauptet-pruefung-die-fehlt` beschreibt. Exit 2 wieder nicht gemessen |
| 8 | **Negative (Config)** — fehlender/leerer Schlüssel bricht ab | `TestMentionsFehlenderSchluesselBrichtAb` + vier `configyaml`-Tests | **ja.** Drei Kern-Fälle plus Config-Rand (nur eine Seite, leeres Glob, ungültiges Glob beidseitig, unbekanntes `match`) |
| 9 | **Determinismus** — byte-identisch, stabile Sortierung | `TestMentionsDeterministisch` | **nein.** Bruch-Test D: `sort.Strings(out)` entfernt ⇒ `make test` bleibt **Exit 0**. `coretest.MemFS.List` sortiert selbst (`memfs.go:105`), also kann die Fixture die Reihenfolge gar nicht verletzen, die der Test prüft. Die zweite Hälfte („nichts wird geschrieben", [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)) misst kein `mentions`-Test |
| 10 | **Default byte-identisch** | — | **nein, kein Test.** Manuell gemessen (Default-`--json` trägt kein `notes`); die Zusage hängt allein an `omitempty` in `report.go:30` und ist ungewächtert |

**Bilanz: 3 von 10 vollständig gemessen** (3, 5, 8), **4 teilweise** (1, 4, 6, 7 — jeweils fehlt der zugesagte Exit-Code), **3 gar nicht** (2, 9, 10).

---

## Findings

### HIGH

**H-1 · Die zugesagte Invariante „jede ADR steht im ADR-Index" hält nicht — die Ist-Menge ist EIN Korpus**
`quelle`: [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) · [`AGENTS.md`](../../AGENTS.md) §4 · `klasse`: `union-korpus-deckt-getrennt-gemeinte-invariante`
`pfad`: `.d-check.yml:886-906`, `harness/sensors/mention-coverage.md:16-20`, `AGENTS.md:422`, `Makefile:197`, `internal/hexagon/core/rules/mentions.go:73,82,144-160`
`verifizierbar`: **ja** — Bruch-Test A oben, wiederholbar.

`.d-check.yml`, die Sensor-Datei, die Makefile-Hilfe und die Commit-Botschaft sagen übereinstimmend, die Mengen-Wahl trage **zwei** Invarianten: *jede ADR steht im ADR-Index* **und** *jede Datei unter `harness/sensors/` ist aus der Sensors-Tabelle verlinkt*. Das Modul kann diese zwei Aussagen nicht trennen: `mentionsCorpus` verkettet die Ist-Menge zu **einem** String (`mentions.go:149-159`), und `strings.Contains` (`:82`) fragt nur, ob die Zeichenkette **irgendwo** darin steht. Gemessen: **21 der 84** ADR-Basisnamen kommen auch in `harness/README.md` vor. Bruch-Test A entfernt die ADR-0084-Zeile vollständig aus `docs/plan/adr/README.md` — `make mention-coverage` meldet danach `108 von 108 Artefakt(en) erwähnt`, `0 Befund(e)`, **Exit 0**. Gehalten ist damit nur die schwächere Aussage *„jedes Mitglied kommt in mindestens einem der beiden Dokumente vor"*; für 21 der 84 ADRs ist der Index-Eintrag ungewächtert. Die Kontroll-Probe B zeigt, dass der Sensor an sich funktioniert — der Defekt sitzt in der Mengen-Wahl, und die Mengen-Wahl ist nach [ADR-0084](../plan/adr/0084-mentions-eigenes-modul.md) §Was aus diesen beiden folgt „das einzige verbliebene Urteil".

**H-2 · Unter `basename` deckt eine Teilzeichenketten-Kollision ein Mitglied ab — `test.md` ⊂ `image-test.md`**
`quelle`: [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) („ein bloßer Dateiname kollidiert") · `klasse`: `zaehlmethode-misst-proxy-statt-gegenstand`
`pfad`: `.d-check.yml:892-898` („beide Soll-Mengen haben eindeutige Basisnamen (gemessen, 0 Duplikate)"), `internal/hexagon/core/rules/mentions.go:80-84`
`verifizierbar`: **ja** — Bruch-Test C oben.

Die Begründung für `match: basename` stützt sich auf eine Messung: *0 Duplikate*. Gemessen ist damit **Gleichheit** von Basisnamen; die Erkennung arbeitet aber über **Enthaltensein** (`strings.Contains`). Über die 108 Mitglieder gibt es genau ein Paar, bei dem der eine Basisname Endstück des anderen ist: `test.md` ⊂ `image-test.md` — und **beide** liegen in der Soll-Menge. Bruch-Test C entfernt den Link `[`make test`](sensors/test.md)` aus der Sensors-Tabelle; `harness/sensors/test.md` ist danach nirgends mehr verlinkt, und `make mention-coverage` meldet `108 von 108`, `0 Befund(e)`, **Exit 0** — gedeckt allein durch `sensors/image-test.md`. Das ist genau die Kollision, gegen die der Default `path` laut Anforderung steht; die Messung, die sie ausschließen sollte, misst den falschen Gegenstand. (Der Register-Eintrag [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../plan/planning/observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md) ist im Plan §7 als „trifft diesen Slice unmittelbar" gesichtet und sein Ableiter ausdrücklich auf den Kalibrierungs-Beleg bezogen worden.)

**H-3 · Die Glob-Semantik der beiden neuen Mengen ist eine andere, als Spezifikation und Handbuch zusagen — `**` trägt nicht**
`quelle`: `spec/spezifikation.md` §`DC-FA-MENT-001.a` Schritt 2 · [`DC-FA-CONF-001`](../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei) · `klasse`: `glob-semantik-zusage-gegen-implementierung`
`pfad`: `spec/spezifikation.md:2801-2806`, `docs/user/benutzerhandbuch.md:1211-1216`, `internal/hexagon/core/rules/mentions.go:109`, `internal/hexagon/core/rules/workflows.go:258-267`
`verifizierbar`: **ja** — Lauf mit `mentions.artifacts: ["tools/**/*.sh"]` unten.

Die neue Spec-Stelle sagt zu: *„Mitglied ist jede Datei, deren '/'-relativer Pfad mindestens eines der Globs aus `mentions.artifacts` trifft (**dieselbe Glob-Semantik wie `scan.ignore`**, `**` als Segment eingeschlossen)"*. `mentionsResolve` benutzt dafür `matchAnyGlob` (`mentions.go:109`), und das ist blankes `path.Match` (`workflows.go:260-267`). `scan.ignore` benutzt dagegen `matchGlob` → `matchSegs` (`internal/hexagon/core/rules/paths.go:98-121`), das `**` segmentweise über **beliebig viele** Segmente auflöst. Gemessen an diesem Repo (11 Dateien unter `tools/`, 5 flach, 6 eine Ebene tiefer):

```
mentions.artifacts: ["tools/**/*.sh"]  gegen docs/user/*.md
⇒ d-check: mentions: 0 von 6 Artefakt(en) erwähnt, über 4 Dokument(e)
mentions.artifacts: ["tools/*.sh", "tools/*/*.sh"]  (dieselbe Menge, explizit)
⇒ 11 Befunde
```

Fünf Mitglieder fallen **still** aus der Soll-Menge; fail-closed greift nicht, weil die Menge nicht leer ist, und die Bezugsmenge „0 von 6" nennt die Sechs, nicht die fehlenden Fünf. Dieselbe Zusage steht im Benutzerhandbuch §5: *„`**` steht für beliebig viele Segmente … Das gilt für **alle** Glob-Felder"*. Und sie trifft die Anforderung selbst: `DC-FA-MENT-001` führt als eigenen Beleg *„die Menge `tools/**/*.sh` gegen `docs/user/` liefert **elf** Funde"* — mit der ausgelieferten Semantik sind es **sechs**.

### MEDIUM

**M-1 · Die zentrale Kalibrierungs-Zahl reproduziert nicht — „84 von 84 und 23 von 23, ausnahmslos falsch"**
`quelle`: [`AGENTS.md`](../../AGENTS.md) §5 · `klasse`: `selbstauskunft-zahl-reproduziert-nicht`
`pfad`: `.d-check.yml:894-898`, `harness/sensors/mention-coverage.md:33-40`, Commit-Botschaft `ac12993` Absatz „DER KALIBRIERUNGS-BEFUND"
`verifizierbar`: **ja** — Läufe unten.

Behauptet ist: *„Mit der Default-Form `match: path` meldet dasselbe Repo 84 von 84 ADRs und 23 von 23 Sensor-Dateien — ausnahmslos falsch."* Gemessen mit **der konfigurierten Mengen-Form** (beide Soll-Globs gegen beide Ist-Dokumente, nur `match: path` gewechselt):

```
d-check: mentions: 21 von 108 Artefakt(en) erwähnt, über 2 Dokument(e)
d-check: 676 Datei(en) geprüft, 87 Befund(e)     (63 ADRs + 24 Sensor-Dateien)
```

Die genannten Zahlen entstehen nur, wenn jede Soll-Menge gegen **ihr eigenes** Dokument gehalten wird (`84 von 84` bzw. `24 von 24`) — also in einer Mengen-Form, die die Konfiguration nicht fährt. Das ist derselbe Bruch wie H-1, hier auf der Mess-Seite: die per-Paar-Zahl wird als Ergebnis des Union-Laufs berichtet. Zwei weitere Ungenauigkeiten in derselben Zeile: **„23" ist 24** (`harness/sensors/*.md` trifft heute 24 Dateien, seit dieser Slice `mention-coverage.md` hinzugefügt hat — die Zahl 108 im selben Commit rechnet bereits mit 24), und **„ausnahmslos"** ist auf der ADR-Seite falsch, weil `harness/README.md` 21 ADRs mit vollem `docs/plan/adr/…`-Pfad verlinkt. Das Argument trägt trotzdem (87 Fehlalarme sind erdrückend), die Ziffern nicht — und sie stehen in **lebenden** Dokumenten, nicht nur in der Botschaft.

**M-2 · DoD (3) verlangt einen Fremd-Bestand; geliefert ist eine Messung am eigenen**
`quelle`: Slice-Plan §2/§4 · [ADR-0084](../plan/adr/0084-mentions-eigenes-modul.md) §Konsequenzen · `klasse`: `eigene-menge-gemessen-fremde-behauptet`
`pfad`: `docs/plan/planning/in-progress/slice-206-mentions-modul.md:60-64`, Commit-Botschaft `ac12993`
`verifizierbar`: **ja** (Diff-Lektüre: keine Fremd-Repo-Messung im Commit).

Der Plan sagt in §2: *„Dieser Slice muss deshalb an einem **Fremd-Bestand** messen und das Ergebnis als solches ausweisen"*, DoD (3) fordert *„Ein **Kalibrierungs-Beleg** an einem Fremd-Bestand"*, und ADR-0084 §Konsequenzen nennt das *„der ehrlichste Preis dieser Entscheidung"*. Der gelieferte Beleg — die einzige Messung im Commit — ist die `path`-vs-`basename`-Gegenüberstellung **an d-check selbst**; kein Fremd-Repo kommt vor. Die Botschaft führt sie unter der Überschrift „DER KALIBRIERUNGS-BEFUND aus DoD (3)". Versagensbild: die Closure hakt DoD (3) ab, und die Ausnahme-Klasse ist an dem Bestand justiert, von dem slice-205 ausdrücklich festgehalten hat, dass er **kein eigenes Rauschen** trägt. (Der Befund selbst ist wertvoll und gehört berichtet — beanstandet ist, dass er als der geforderte durchgeht.)

**M-3 · Der Determinismus-Test misst die stabile Sortierung nicht**
`quelle`: [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) · `klasse`: `wortlaut-behauptet-pruefung-die-fehlt`
`pfad`: `internal/hexagon/core/rules/mentions_test.go:174-198`, `internal/hexagon/core/coretest/memfs.go:101-106`
`verifizierbar`: **ja** — Bruch-Test D.

`TestMentionsDeterministisch` ruft fünfmal dieselbe Funktion auf derselben `MemFS` und vergleicht die Befund-Reihenfolge, zuletzt gegen die feste Erwartung `tools/a.sh|tools/m.sh|tools/z.sh`. `MemFS.List` sortiert seine Einträge selbst (`memfs.go:105`), und der Walk ist tiefen-erst über sortierte Namen — die Fixture kann die Reihenfolge, gegen die geprüft wird, gar nicht verletzen. Bruch-Test D entfernt `sort.Strings(out)` aus `mentionsResolve`; `make test` bleibt **Exit 0**. Die Sortierung ist real lasttragend (DFS-Reihenfolge und lexikografische Pfad-Reihenfolge fallen auseinander, sobald ein Dateiname und ein Verzeichnisname sich am Trennzeichen unterscheiden — `tools.md` < `tools/a.sh`), aber ungewächtert. Genau der Register-Eintrag, den Plan §7 mit 7× Belegen als „unmittelbar einschlägig für die Tests" gesichtet hat.

**M-4 · Der neue, generische `Notes`-Vertrag hat keinen einzigen Test**
`quelle`: [`DC-FA-CLI-004`](../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate) · Reviewer-Anker „fehlende Negativtests bei neuem öffentlichen Vertrag" · `klasse`: `neuer-vertrag-ohne-negativtest`
`pfad`: `internal/adapter/driven/report/report.go:24-31,58-62,136-140`, `internal/hexagon/core/rules/run.go:13-19`, `internal/adapter/driving/cli/cli.go:548`
`verifizierbar`: **ja** — `grep -rn "Notes" --include=*_test.go internal/` liefert null Treffer.

`Summary` bekommt ein drittes, serialisiertes Feld; `Text`, `doctorSummary`, `JSON`, `YAML`, `DoctorJSON` und `DoctorYAML` geben es aus. Getestet ist davon nichts: nicht die stderr-Zeile, nicht `summary.notes` unter `--json`/`--yaml` (Akzeptanzkriterium 2, das die *beiden* Formen ausdrücklich zur Festlegung erhebt), und vor allem nicht die Rückwärtskompatibilität — dass ohne aktives `mentions` **kein** `notes`-Schlüssel erscheint (Akzeptanzkriterium 10), hängt allein am `omitempty` in `report.go:30`. Fiele es weg, trüge jede JSON-Ausgabe des Produkts `"notes": null`, und kein Test meldete es. Ich habe beide Formen manuell gemessen (siehe §Gefahrene Läufe); der Punkt ist, dass die Messung nicht wiederholbar verankert ist. Das Slice-Risiko §5 („ob das rückwärtskompatibel bleibt, entscheidet sich an den vorhandenen Konsumenten") ist damit unbeantwortet geblieben, nicht beantwortet.

**M-5 · Das Modul liest den ganzen Baum ab der Repo-Wurzel und honoriert weder `scan.roots` noch `scan.ignore` — unbenannte Grenze, und die Spec sagt etwas anderes**
`quelle`: [`AGENTS.md`](../../AGENTS.md) §3.8 · [`BEO-ALL/module-promise-only-on-scan-axis`](../plan/planning/observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md) · `klasse`: `module-promise-only-on-scan-axis`
`pfad`: `internal/hexagon/core/rules/mentions.go:106,119-142`, `spec/spezifikation.md:2799-2806`, `internal/hexagon/core/rules/scan.go:113`
`verifizierbar`: **ja** — Läufe unten.

`mentionsWalk(fsys, "", …)` startet immer an der Repo-Wurzel und prunt ausschließlich über `isSkipDir`. Die Spec sagt: *„Der Dateibaum wird ab der Scan-Wurzel rekursiv gelaufen, unter **derselben Skip-Liste wie die Markdown-Discovery**"* — die Discovery prunt zusätzlich über `scan.ignore` (`scan.go:113`: `isSkipDir(e.Name) || dirIgnored(rel, ignore)`). Gemessen:

```
scan.ignore: [".harness/baseline/**"] ; artifacts: [".harness/baseline/*/regelwerk/modul-0*.md"]
⇒ mentions: 0 von 10 Artefakt(en) erwähnt      (10 Befunde aus dem ausgenommenen Baum)
scan.roots: ["docs"] ; artifacts: ["harness/sensors/*.md"]
⇒ mentions: 0 von 24 Artefakt(en) erwähnt      (Mitglieder außerhalb der deklarierten Wurzel)
```

Das Verhalten ist mit den Geschwister-Modulen (`reviews`, `workflows`, `targets`) konsistent — beanstandet ist nicht die Wahl, sondern dass sie **nirgends benannt** ist: Der Modul-Kommentar spricht die Zwei-Eingaben-Asymmetrie aus (`mentions.go:15-20`, wie Plan §7 verlangt), aber nicht diese zweite Achse; die Spec behauptet stattdessen Deckung mit der Discovery, und das Lastenheft sagt „relativ zur **Scan-Wurzel**", was ein Leser auf `scan.roots` bezieht. Versagensbild: ein Adopter, der einen vendorten Fremdbaum per `scan.ignore` ausnimmt, bekommt ihn über `mentions.artifacts: "**/*.md"` als Soll-Mitglieder zurück.

**M-6 · Die Konfigurations-Schema-Tabelle der Spezifikation führt die drei neuen Schlüssel nicht**
`quelle`: [`DC-FA-CONF-001`](../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei) · `klasse`: `semantic-change-body-only-edges-stale`
`pfad`: `spec/spezifikation.md` §2 (Konfigurations-Tabelle; `workflows.dir`, `reviews.done-dir`, `targets.*`, `tracked.exempt-targets` stehen dort)
`verifizierbar`: **ja** — `grep -n "mentions\." spec/spezifikation.md` trifft nur §`DC-FA-MENT-001.a` und §4.

Jedes Vorgänger-Modul hat mit seiner Verfeinerungs-Sektion auch seine §2-Schema-Zeilen bekommen; die §7-Historie protokolliert das wörtlich (*„§2-Schema (`targets.makefiles`/…) ergänzt"*, *„§2-Schema (`tracked.exempt-targets`)"*). `mentions.artifacts`, `mentions.documents` und `mentions.match` fehlen dort — Typ, Default und die Exit-2-Bedingungen stehen ausschließlich in der Prosa von §`DC-FA-MENT-001.a` Schritt 1. Versagensbild: die Tabelle, die ein Konsument als Schema-Nachschlagewerk liest, ist für das neue Modul stumm, und der nächste Slice, der sie fortschreibt, hält sie für vollständig.

**M-7 · `.d-check.yml` nennt ein Make-Target, das es nicht gibt (`make mentions-check`)**
`quelle`: [`AGENTS.md`](../../AGENTS.md) §4 („Nur hier gelistete Targets existieren im Makefile") · Baseline `modul-13` §Hard Rule · `klasse`: `phantom-target-in-config-kommentar`
`pfad`: `.d-check.yml:884` gegen `Makefile:197` (`mention-coverage`)
`verifizierbar`: **ja** — `grep -c "^mentions-check:" Makefile` ⇒ 0.

Der Begründungs-Block der Mengen-Wahl sagt: *„Opt-in, läuft nur über `make mentions-check`, nicht in `gates`"*. Das Target heißt `mention-coverage`. Alle übrigen `make X`-Nennungen in `.d-check.yml` (`adr-check`, `trace-check`, `planning-check`, `gate-consistency`, `trace`, `completeness-check`) treffen reale Regeln — diese als einzige nicht. `make gate-consistency` fängt es nicht: das Modul `targets` liest `make X` nur aus **Doku-Tabellenzeilen**, nicht aus YAML-Kommentaren (gemessen: der grüne `gates`-Lauf). Kein stilles Grün — der Aufruf scheitert laut —, aber die Datei, die die Mengen-Wahl begründet, nennt den falschen Träger. **Einstufung MEDIUM statt HIGH**, weil der Fehlschlag laut ist; wer die Kanon-Regel („keine Befehle behaupten, die es nicht gibt") streng liest, kann eine Stufe höher gehen.

**M-8 · Ein unlesbares Verzeichnis verkleinert die Soll-Menge still**
`quelle`: Reviewer-Anker „Konsistenz-Lücke zwischen Modulen derselben Eingabe-Klasse" · `klasse`: `stiller-skip-bei-unlesbarer-eingabe`
`pfad`: `internal/hexagon/core/rules/mentions.go:123-126`
`verifizierbar`: **ja** (Code-Lektüre; Reproduktion braucht ein nicht lesbares Verzeichnis im Mount).

`mentionsWalk` verwirft den `List`-Fehler (`if err != nil { return }`). Ein unlesbarer Teilbaum entfernt seine Dateien aus **beiden** Mengen, ohne dass irgendetwas gemeldet wird; fail-closed greift nur, wenn danach *gar nichts* übrig ist. Die Nachbarn derselben Eingabe-Klasse machen es anders und begründen es: `reviews` meldet ein unlesbares `reviews-dir` als denselben fail-closed-Befund wie null Kandidaten, `workflows` meldet unlesbares YAML als `workflow-unparsable` statt es zu überspringen. Für die **Ist**-Seite ist der Fall ausgesprochen (`mentions.go:146-148`, Spec Schritt 4: „Ein unlesbares Dokument trägt nichts bei, zählt aber als Mitglied") — für die **Verzeichnis**-Ebene nirgends.

### LOW

**L-1 · Der Plan zählt neun Akzeptanzkriterien; es sind zehn — und die zehnte ist die einzige ohne Test**
`pfad`: `docs/plan/planning/in-progress/slice-206-mentions-modul.md:184-187` · `klasse`: `selbstauskunft-zahl-reproduziert-nicht`
Der Sichtungs-Block zu [`BEO-ALL/spec-randbedingung-ohne-test`](../plan/planning/observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md) sagt: *„die Anforderung trägt **neun** Akzeptanzkriterien … Jedes einzelne braucht seinen Test."* Gezählt sind es zehn (`sed -n '3530,3548p' spec/lastenheft.md | grep -c '^- \*\*'` ⇒ 10). Die aus der Zählung fallende letzte — *Default byte-identisch* — ist genau die ohne Test. `verifizierbar`: ja.

**L-2 · Der neue Make-Eintrag übernimmt den Kommentar seines Nachbarn; `review-coverage` bleibt unkommentiert**
`pfad`: `Makefile:192-200` gegen `git show 401ec1d:Makefile` · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.7 · `klasse`: `kommentar-bindet-an-falsches-ziel`
Der Block *„Netzlos, fail-closed auch bei leerer **Kandidatenmenge** …"* stand vor dem Commit über `review-coverage` (dessen Wortschatz er trägt) und steht jetzt über `mention-coverage`; `review-coverage` hat keinen mehr. Ein Kommentar beschreibt, was da ist — hier beschreibt er den Nachbarn. `verifizierbar`: ja (Diff).

**L-3 · Die Gate-Taxonomie-Tabelle führt `mention-coverage` nicht**
`pfad`: `harness/README.md:148-152` · `klasse`: `semantic-change-body-only-edges-stale`
Die Sensors-Tabelle hat ihre Zeile bekommen, die Taxonomie darunter nicht — obwohl sie `review-coverage` als *„eigenständiger Fokus-Lauf, **nicht** in `gates`/`ci`"* ausdrücklich führt und in der Bindepunkt-Spalte jede Gate-Klasse aufzählt. `make gate-consistency` deckt das nicht (es prüft `make X` gegen Makefile-Regeln, nicht die Klassen-Zuordnung). `verifizierbar`: ja.

**L-4 · Der Baum wird zweimal vollständig gelaufen**
`pfad`: `internal/hexagon/core/rules/mentions.go:103-117` (`mentionsResolve` ruft `mentionsWalk` je Aufruf; `CheckMentions:64,68` ruft es zweimal) · `quelle`: [`DC-QA-01`](../../spec/lastenheft.md#dc-qa-01--performance) · `klasse`: `wiederholter-vollscan-statt-einmal-sammeln`
Die gesammelte Datei-Liste ist für beide Mengen dieselbe und wird zweimal erzeugt. Bei einem großen Bestand ist das der doppelte Verzeichnis-Durchlauf zusätzlich zum Markdown-Scan; das Modul ist opt-in, deshalb LOW. `verifizierbar`: ja (Code).

**L-5 · Der Kommentar des wiederverwendeten Glob-Helfers beschreibt weder seinen neuen Konsumenten noch die Wahrheit**
`pfad`: `internal/hexagon/core/rules/workflows.go:258-259` · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.7 · `klasse`: `kommentar-bindet-an-falsches-ziel`
*„matchAnyGlob prueft die **Ventil**-Globs (Go path.Match, **wie jedes andere d-check-Glob**)."* Seit diesem Commit prüft er zusätzlich die **mengendefinierenden** Globs von `mentions` — kein Ventil —, und die Klammer-Aussage ist gegen `matchGlob`/`ignored` (`paths.go:98-131`) falsch. Die Zeile ist der Ursprung von H-3. `verifizierbar`: ja.

**L-6 · Keine §7-Historie-Zeile für §`DC-FA-MENT-001.a`/`SPEC-082`**
`pfad`: `spec/spezifikation.md:3272-3276` (jüngste Zeile 2026-09-02) · `klasse`: `semantic-change-body-only-edges-stale`
Jede frühere Modul-Verfeinerung hat ihre Zeile. **Nicht entlastend, aber ehrlich zu nennen:** das unmittelbar vorangegangene Modul `reviews` (§`DC-FA-RVW-001.a`) hat ebenfalls keine — der Bestand hat den Bruch schon einmal gemacht. `verifizierbar`: ja.

### INFO

- **I-1 · `Config.Mentions` ohne Feld-Kommentar** (`internal/hexagon/core/model/config.go:89`), während `Planning`, `Tracked`, `Targets`, `Workflows` und `Reviews` je eine `// <Name>: Parameter des Moduls … (DC-…)`-Zeile tragen. `Structure` trägt ebenfalls keine — es gibt also Präzedenz; §3.7 regelt Kommentare, die es gibt, nicht ihre Abwesenheit.
- **I-2 · Der Konfigurations-Referenzblock des Benutzerhandbuchs (`docs/user/benutzerhandbuch.md:1200-1310`) dokumentiert `mentions:` nicht** — genau wie er `reviews:` nicht dokumentiert. Ein Nutzer, der das Modul aus der Modul-Tabelle (`:2335`) aktiviert, findet die Pflicht-Schlüssel im Handbuch nirgends und bekommt Exit 2. Bestands-Muster, deshalb INFO statt Finding gegen diesen Slice.
- **I-3 · Bestandsbefund neben H-3:** `reviews.exempt-paths` und `workflows.exempt-paths` benutzen denselben `path.Match`-Helfer und brechen damit die Handbuch-Zusage *„Das gilt für alle Glob-Felder … und die `exempt-paths` der Module"* bereits vor diesem Slice. Nicht diesem Commit anzulasten, aber dieselbe Wurzel.
- **I-4 · `target` ist die komma-verkettete Glob-Liste** (`mentions.go:74`), während das Lastenheft *„das Glob der Ist-Menge"* (Singular) sagt — die in R2-L-2 benannte Unbestimmtheit. Sie ist jetzt **entschieden und dokumentiert** (Spec Schritt 5: „komma-getrennt"), also aufgelöst; festgehalten nur, damit die Kette lesbar bleibt.

---

## Negativbefunde

- **Alle Gate-Läufe der Botschaft sind echt und reproduzieren.** `make gates` Exit 0 über **zehn** Gates, `676 Datei(en) geprüft, 0 Befund(e)`, `Coverage 94.50%`, semgrep `0 findings` — Wort für Wort wie in der Botschaft. `make mention-coverage` Exit 0 mit `108 von 108 … über 2 Dokument(e)`.
- **Die Mengen-Größe 108 reproduziert exakt.** `docs/plan/adr/[0-9]*.md` ⇒ 84, `harness/sensors/*.md` ⇒ 24, Summe 108. Auch „0 Duplikate" stimmt: über beide Mengen gibt es keinen doppelten Basisnamen (beanstandet ist in H-2 nur, dass Gleichheit den falschen Gegenstand misst).
- **Der Bruch-Test der Botschaft reproduziert.** Eine ADR ohne Index-Eintrag wird mit `artifact-unmentioned`, `file` = Artefakt-Pfad, `line` = 1 gemeldet, die Bezugsmenge fällt auf `107 von 108`, Exit 1 — Bruch-Test B, exakt wie beschrieben.
- **Alle drei Exit-2-Pfade erreichen wirklich Exit 2**, nicht Exit 1: leere Soll-Menge, leere Ist-Menge und aktives Modul ohne Block enden je als `d-check: error:` mit Exit 2. Die Fehlermeldungen nennen jeweils das verantwortliche Glob bzw. beide Pflicht-Schlüssel. Der Fehler wandert über `CheckMentions` → `runPostPasses` → `RunWithVCS` → CLI durch, ohne unterwegs Befund zu werden.
- **Kein Vorkommen kann über eine Dokumentgrenze hinweg entstehen.** `mentionsCorpus` schreibt nach jedem Dokument ein `'\n'` (`mentions.go:157`) — auch nach leerem Inhalt und auch, wenn das Dokument nicht mit Zeilenumbruch endet. Ein Suchbegriff ist ein Pfad oder Basisname und enthält kein `\n`, kann die Fuge also nicht überspannen. Die Hypothese ist **widerlegt**.
- **Die drei angefassten Modul-Mengen-Spiegel sind wirklich gegatet, und die Behauptung ist vollständig für das, was sie behauptet.** `validModules()` ↔ Verfügbar-Zeile des Config-Templates hält `cli_acceptance_test.go:3031-3038`; `AllReasons()` ↔ `reasonTexts()` ↔ Spezifikation §4 halten `diagnose_test.go:17-56`; die Doku-Tabellen hält `make gate-consistency`. **Eigene Suche nach weiteren Spiegeln:** `FOCUS_DISABLE` (`Makefile:374-376`) und `netlessDocModules()` (`gate_consistency_test.go:17-52`) spiegeln die `.d-check.yml`-`modules:`-Zeile, in der `mentions` als opt-in-Modul zu Recht **nicht** steht — der Verzicht auf einen Nachzug ist korrekt, nicht vergessen. `--print-config` trägt nur die Verfügbar-Zeile (kein Modul-Block), `--print-mk` keine Modul-Aufzählung. `DC-FA-CLI-002` (`lastenheft.md:111`) und `operations.md:33` führen `mentions` bereits. **Ein vierter Spiegel ist ungepflegt geblieben** — die Gate-Taxonomie (L-3) —, aber er spiegelt Targets, nicht die Modul-Menge.
- **`mentions.scope` wird abgelehnt, nicht stillschweigend geschluckt.** `rawMentions` trägt kein `Scope`-Feld, und der Decoder läuft mit `KnownFields(true)`: `field scope not found`, Exit 2. Das entspricht [`DC-FA-CONF-002`](../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope) („ein wirkungsloser Schlüssel, den der Decoder annimmt, ist ein stilles Grün"). Dass die Aufzählung der nicht-scannenden Module dort `mentions` nicht nennt, teilt das Modul mit `workflows` und `reviews` — Bestands-Drift, kein Neuzugang.
- **Der Befund-Ort ist der in ADR-0084 festgelegte, und die Bestands-Behauptung stimmt.** `file` = Artefakt-Pfad, `line` = 1, `target` = Ist-Globs; **sechs** Produkt-Module führen `Line: 1` bereits als Platzhalter (`commits.go`, `links_resolvefrom.go`, `reviews.go`, `structure.go`, `vcs.go`, `workflows.go`) — nachgezählt, die Zahl stimmt. Der `--doctor`-Klartext benennt den Platzhalter ausdrücklich als solchen.
- **Plan §3 (Out-of-Scope) ist eingehalten, alle drei Punkte.** Keine zweite Quell-Form (nur Pfad-Globs, `rawMentions` kennt keine Literale); `mentions` ist **nicht** in `make gates` (der Gates-Lauf nennt zehn Gates ohne es, und `.d-check.yml`-`modules:` führt es nicht); weder `spec/lastenheft.md` noch `docs/plan/adr/0084-…md` sind im Diff von `ac12993`.
- **§3.5 ist gewahrt** — ADR-0084 ist unberührt; `make adr-check` lief im Rahmen der Hook-Kette nicht rot, und der Diff zeigt die Datei nicht.
- **§3.1/§3.2 sind gewahrt** — kein Host-Werkzeug im Diff, keine `//nolint`-Direktive; der Botschafts-Satz zum exportierten `MentionsResult` stimmt (`golangci`s `revive`-Familie verlangt bei exportierter Funktion einen exportierten Rückgabetyp; der Typ ist exportiert statt unterdrückt, und `make lint` ist grün).
- **§3.4 ist gewahrt** — die neue Spec-Sektion nennt keine ADR, keinen Slice, keine Welle, keinen Commit-Hash; sie verweist ausschließlich auf `lastenheft.md`-Anker und Modulnamen. `matrix` ist grün.
- **§3.7 hält im neuen Code** — der Modul-Kopf trägt Zusage, Abgrenzung und Grenze; keine Review-Historie, keine Slice-Nummer, keine Mess-Label. Herkunft steht je als **ein** auflösbares Feld (`DC-FA-MENT-001`, `ADR-0084`, `AGENTS.md §3.8`). Die Ausnahmen sind L-2 und L-5, und beide betreffen Kommentare an falschem Ziel, nicht falsche Klassen.
- **§5, Release-Prep-Grenze:** dass die beiden `README*.md`, der Handbuch-Kopf und `CHANGELOG.md` **nicht** angefasst sind, ist regelkonform und in der Botschaft korrekt begründet — die zitierte Regel deckt genau diese drei Ziele. Die Handbuch-**Modul-Tabelle** ist dagegen Körper und richtigerweise mitgezogen.
- **Die Erkennungsform-Wahl ist im Ergebnis richtig, auch wenn die Zahl in M-1 nicht stimmt.** `match: basename` ist für dieses Repo die tragfähige Form (87 Fehlalarme unter `path`, gemessen), und die Schlussfolgerung der Botschaft — der Befund widerlege die **Default**-Wahl von ADR-0084 nicht, weil dort an Mengen gemessen wurde, deren Dokumente Vollpfade schreiben — **trägt**: die ADR-Messung ist in `## Geschichte` (M-2-Berichtigung) über drei benannte Fremd-Mengen dokumentiert, und der einzige dort gemessene Formunterschied ist ein Fehlalarm der laxen Form. **Rückführung §6 „blockiert" also zu Recht nicht gezogen** — keine Festlegung aus ADR-0084 ist widerlegt. (Was der Befund *doch* berührt, ist nicht die Default-Wahl, sondern die Kollisions-Begründung — siehe H-2.)
- **Rückführung §6 „zu groß" ebenfalls zu Recht nicht gezogen** — der Diff ist in einer Sitzung prüfbar (18 Dateien, drei Liefer-Punkte, zwei Schichten: Kern + Adapter/Doku).
- **Der Ablauf des Moduls im Lauf ist korrekt eingehängt** — `mentions` ist Post-Pass wie `reviews`/`targets`, `Run` delegiert an `RunWithVCS` (also verlieren Alt-Aufrufer die Notes nicht), Befunde gehen durch `model.SortFindings`, und die Reihenfolge der Post-Pässe ist für das Ergebnis irrelevant.
- **Die Config-Validierung am Rand ist vollständig gegen die Spec Schritt 1** — nur eine Seite gesetzt, Weißraum-Glob, ungültiges Glob **auf beiden Seiten**, unbekanntes `match`: je Exit 2, je mit Test. `path.Match(g, "probe")` ist als Gültigkeitsprobe ausreichend, weil `ErrBadPattern` allein vom Muster abhängt.
- **Kein Netzzugriff, kein Schreibpfad** — das Modul benutzt ausschließlich `driven.Filesystem.List`/`ReadFile`; `make arch-check` ist grün, der Lauf fährt `--network none` und `-v …:ro`.

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 3 |
| MEDIUM | 8 |
| LOW | 6 |
| INFO | 4 |

**Wiederkehrende Klassen — und die Steering-Loop-Beobachtung dieses Laufs.** H-1, M-1 und (auf der Mess-Seite) H-2 teilen **eine** Wurzel: *eine Aussage über die Vereinigung wird als Aussage über ihre Teile geführt*. Das Modul hält eine Menge gegen **einen** Korpus; die Doku, die Sensor-Datei und die Botschaft sprechen von **zwei** Invarianten, und die Kalibrierungs-Zahlen entstehen aus **zwei** getrennten Läufen. Das ist dieselbe Klasse, die R1 als `praemisse-gegen-ein-dokument-statt-gegen-die-menge` gemeldet hat — hier mit umgekehrtem Vorzeichen: nicht *gegen eines statt gegen die Menge*, sondern *gegen die Menge, berichtet als wären es einzelne*. Kandidat für einen Registereintrag; der bestehende trägt die Richtung nicht.

Zweite Klasse, dreifach belegt: **die Zusage steht, die Probe misst sie nicht** (Akzeptanzkriterien 2, 9, 10; M-3 mit Bruch-Test bewiesen) — [`BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt`](../plan/planning/observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md), 7× vorbelegt und im Plan §7 ausdrücklich als „unmittelbar einschlägig für die Tests" gesichtet. Dass die Sichtung stattfand und die Lücke trotzdem entstand, ist selbst das Signal.

Dritte, unverändert: `selbstauskunft-zahl-reproduziert-nicht` (M-1, L-1) — der siebte Lauf in Folge.

## Verdikt

**Blockierend — H-1, H-2, H-3.**

Die drei HIGH sind je durch einen wiederholbaren Bruch-Test bzw. einen Lauf belegt und betreffen das, was der Sensor über sich behauptet: Zwei zeigen, dass `make mention-coverage` grün bleibt, während die Invariante gebrochen ist, die er zu halten verspricht; die dritte, dass die zugesagte Glob-Semantik nicht geliefert wird und Mitglieder still aus der Soll-Menge fallen. Von den zehn Akzeptanzkriterien sind **drei vollständig** gemessen, vier nur bis zum Exit-Code, drei gar nicht.

**Nicht blockierend, aber vor der Closure zu entscheiden:** M-2 — DoD (3) verlangt einen Fremd-Bestand, geliefert ist eine Messung am eigenen. Das ist eine Plan-/DoD-Frage und gehört der Verifikation; hier steht sie, weil die Commit-Botschaft die eigene Messung unter der DoD-(3)-Überschrift führt.
