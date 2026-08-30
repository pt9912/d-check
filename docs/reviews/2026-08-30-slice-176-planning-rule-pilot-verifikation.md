# Verifikation slice-176 — unabhängiger Lauf

**Gegenstand:** [slice-176](../plan/planning/done/slice-176-planning-rule-pilot.md) §4, Commits `8879ff1 · 776c0dd · a8f2c6b · b0850a9 · 0cfe876`. HEAD = `0cfe876`.
**Arbeitsort:** alle Mutationsproben in `…/scratchpad/probe` (eigenes `git init`, Kopie von Skript + Baseline-Baum + Aliasen); Scan-Messungen gegen `git archive`-Abzüge in `/tmp`, gemountet als `/repo:ro`. Im Repo selbst nur zwei lesende Läufe (`make baseline-verify`, `git`-Abfragen).

## 1. DoD-Tabelle

| # | DoD-Punkt | Urteil | Messung |
|---|---|---|---|
| 1 | vier Symlinks, `readlink -e` alle OK, **805** von 4787 Zeilen | **erfüllt** | s. §2 |
| 2 | Zustellung belegt, Beleg nennt die Dateien (`/memory`) | **nicht selbst verifizierbar** — Fremdbeleg, Beleg-Grenze im Slice **nicht** benannt | s. §2 / A-9 |
| 3 | Preis gemessen: +29,4k Token, 59 868 Bytes, ~2,0 B/Token | **teilweise** — Bytes und Arithmetik erfüllt, Token-Zahlen Fremdbeleg | s. §2 |
| 4 | Bump-Träger **gewächtert**, gemessen, als `MR-055` geführt | **erfüllt**, Wächter beweisbar tragend — aber vier unbenannte blinde Flecken + vier nicht nachgezogene Spiegel | s. §3 / A-1…A-8 |
| 5 | Nicht-Zusagen geschrieben (§2, „fünf Punkte") | **erfüllt** (5 vorhanden); die DoD-Zeile selbst zählt nur vier auf | A-11 |
| 6 | Nachfolge-Entscheid benannt (§2) | **erfüllt** — Hook nachrangig mit Kriterium (29,4k), weitere Module, `AGENTS.md`-Schnitt |
| 7 | `make gates` grün, unabhängiger Review, Verifikation | **offen** (vom Slice selbst offen markiert). Von mir gefahren: die hermetische Hälfte (unten). Diese Verifikation ist der eine der beiden Läufe. |

**Punkt 7, was ich selbst gefahren habe** (gegen den `git archive HEAD`-Abzug, `/repo:ro`):

```
docker run --rm --network none -v <HEAD-Abzug>:/repo:ro d-check:latest
  → d-check: 600 Datei(en) geprüft, 0 Befund(e)   EXIT=0      (doc-check: links/anchors/ids/matrix/codepaths/spans/hostpaths/versions/structure/diagrams/citations)
--enable planning  → 600 Datei(en), 0 Befund(e)   EXIT=0
--enable targets   → 600 Datei(en), 0 Befund(e)   EXIT=0
--enable workflows → 600 Datei(en), 0 Befund(e)   EXIT=0
make baseline-verify (echtes Repo) → verify ok (51 Dateien, vollständig)  EXIT=0
```

Nicht von mir gefahren: `lint`, `test`, `arch-check`, `coverage-gate`, `semgrep` — sie hätten den geteilten Docker-Tag `d-check:latest` neu gebaut und damit in die Parallel-Sitzung hineingegriffen. Gemessen statt behauptet, dass sie unberührt sind: `git diff --name-only 8879ff1^..0cfe876` enthält **keine** Go-Datei, kein `Dockerfile`, kein `.golangci.yml` — nur 4 Symlinks, 10 Markdown-Dateien und `tools/harness/fetch-baseline-cache.sh`. Ergänzend, **nicht mein Beleg**: `.harness/state/gates-passed.diffsha` trägt `f37462f7…`, identisch mit dem aktuellen `working-tree-hash.sh`-Ergebnis — ein `make gates` ist auf genau diesem Baum gelaufen (07:07, Parallel-Sitzung).

## 2. Nachgeprüfte Zahlen

**Alle vier Aliase lösen auf** (`readlink -e`, im echten Repo):

```
modul-01-entwicklungszyklus.md → /Development/d-check/.harness/baseline/v5.12.0/regelwerk/modul-01-entwicklungszyklus.md
modul-05-planning-harness.md   → …/modul-05-planning-harness.md
modul-06-roadmap.md            → …/modul-06-roadmap.md
modul-08-agentenrollen.md      → …/modul-08-agentenrollen.md
```

Alle vier sind git-Modus `120000` (echte Symlinks im Index, keine Textdateien), und es sind die **einzigen** vier Symlinks im ganzen Repo (`git ls-files -s | awk '$1=="120000"'` und `find . -type l`). Der Geltungsbereich des Wächters deckt heute also 100 % der Aliase.

**805 Zeilen — bestätigt.** `wc -l`: 86 + 248 + 228 + 243 = **805**. Gesamtbaum `cat regelwerk/*.md | wc -l` = **4787**. Beides exakt wie im Plan.

**59 868 Bytes — bestätigt.** `wc -c`: 4676 + 14848 + 25979 + 14365 = **59 868**.

**Arithmetik der Ableitung — korrekt.** 58,3k − 28,9k = **29,4k**; 59 868 / 29 400 = **2,036** ⇒ „rund 2,0 Bytes je Token" trägt. Nebenrechnungen im Slice ebenfalls korrekt: 29,4/1000 = 2,94 % („2,9 % des Fensters"), 29,4/200 = 14,7 % („rund 15 %"), 805 + 527 = 1332 („über 1300 Zeilen"). `wc -l AGENTS.md` = **527**, `templates/AGENTS.template.md` = **236** — beide Zahlen aus §1 stimmen.

**Beleg-Grenze — nicht benannt (A-9).** Die Größen `58.3k`/`28.9k` und die `/memory`-Auflistung stammen aus einer fremden Sitzung des Auftraggebers; sie sind in keiner Datei des Repos abgelegt und von hier aus nicht reproduzierbar. Der Slice sagt das an keiner Stelle — weder in §4 noch in §5 noch in den fünf Nicht-Zusagen. Nur die Commit-Botschaft `b0850a9` nennt die Herkunft („Der Auftraggeber hat eine frische Sitzung gestartet und zwei Anzeigen geliefert"). **Verschärfend:** im Scratchpad liegt der verworfene Entwurf `beleg.txt`, der die Grenze noch ausdrücklich trug — *„Die Grenze des Belegs, benannt: `/context` führt keine Zeile je Datei; die Zuordnung ruht auf dem Delta plus der Größenordnung, nicht auf einer Namensnennung."* Beim Umschreiben auf den stärkeren `/memory`-Beleg (`beleg2.txt` → Slice) ist der Grenz-Satz ersatzlos entfallen. Die `/memory`-Namensnennung entkräftet zwar *diese* Grenze — die verbleibende, andere Grenze (Fremdherkunft, nicht wiederholbar) tritt aber nicht an ihre Stelle.

**Gegenbeobachtung zur Zustellung (A-10).** In *diesem* Kontext — einem Subagenten-Thread desselben Werkzeugs — sind die vier Module **nicht** präsent: der `claudeMd`-Block führt `CLAUDE.md`, `AGENTS.md` und die Nutzer-`MEMORY.md`, keinen der vier Regelwerk-Pfade. Das widerspricht dem Beleg nicht (der galt einer frischen Haupt-Sitzung), aber es ist eine Reichweiten-Grenze, die keine der fünf Nicht-Zusagen nennt — und sie trifft genau die Rollen, deren Ausfall den Anlassfall bildete: ein Review- oder Verifikations-Subagent bekommt `modul-01` (Zyklus) und `modul-06` (Register) nach dieser Messung **nicht** mit. Das ist die Gestalt von `BEO-024` eine Ebene tiefer: der Kanal hängt nicht nur am Werkzeug, sondern am **Sitzungstyp** innerhalb des Werkzeugs.

## 3. Mutationsproben am Sensor (`BEO-023`)

Sandbox: Kopie von `fetch-baseline-cache.sh`, `.harness/baseline/v5.12.0/`, `.claude/rules/`, `harness/conventions.md`-Standzeile; eigenes `git init`. Kommando je Probe: `cd <probe> && bash tools/harness/fetch-baseline-cache.sh --verify`.

| Probe | Erwartung | Ergebnis | Urteil |
|---|---|---|---|
| **P0** unverändert | Exit 0 | `verify ok (51 Dateien, vollständig)` · **Exit 0** | Kontrolle |
| **P1** ein toter Symlink | Exit ≠ 0 + Meldung | `toter Symlink .claude/rules/probe-tot.md — Ziel fehlt (Baseline-Bump?)` · **Exit 1** | **beißt** |
| **P1c** getreue Bump-Simulation (Tag-Verzeichnis `v5.12.0`→`v5.13.0` **und** `**Stand:**`-Zeile gehoben) | alle vier gemeldet | **alle vier** einzeln gemeldet · **Exit 1** | **der eigentliche Fall — gedeckt** |
| **P2** Prüfung entfernt (Zeilen 135–144 gelöscht), toter Symlink bleibt | Exit 0 | `verify ok` · **Exit 0** | **Wächter ist tragend** |
| **P3** Bedingung invertiert (`&& continue` → `\\|\\| continue`), alle Links gesund | Exit ≠ 0 | alle vier gesunden Aliase gemeldet · **Exit 1** | **Bedingungsrichtung tragend** |
| **P4** `.claude/rules/` fehlt ganz | ? | `verify ok` · **Exit 0**, keine Meldung | **stiller Pfad** (A-5) |
| **P4b** `.claude/rules/` leer | Exit 0, kein Falsch-Positiv | `verify ok` · **Exit 0** | korrekt (kein unexpandiertes Glob) |
| **P5** Symlink auf ein **Verzeichnis** | ? | `verify ok` · **Exit 0** | **passiert** (A-6) |
| **P6** Symlink-**Schleife** (a→b→a) | ? | beide gemeldet · **Exit 1** | erkannt, **Ursache falsch benannt** (A-8) |
| **P7** **echte** Datei statt Symlink (mit totem Link darin) | ? | `verify ok` · **Exit 0** | erwartet (`[ -L ]`-Filter) |
| **P8** Symlink **außerhalb** des Pins (löst auf) | Exit 0 | `verify ok` · **Exit 0** | **benannte** Grenze — korrekt |
| **P9** toter Symlink als **Dotfile** (`.probe-tot.md`) | ? | `verify ok` · **Exit 0** | **blind** (A-3) |
| **P10** toter Symlink in `.claude/rules/**sub/**` | ? | `verify ok` · **Exit 0** | **blind** (A-4) |
| **P11** toter Symlink **und** SHA-Drift | ? | `regelwerk/modul-01…: GESCHEITERT` · **Exit 1**, Symlink **nicht** gemeldet | fail-closed, aber dritte Frage übersprungen (A-7) |

**Modus-Abdeckung (statisch, `sed`):** die Prüfung sitzt in `verify()`. Erreicht wird sie von `--verify` **und** vom re-vendor-Pfad (der `verify` am Ende ruft) — **nicht** von `--check-latest`, also nicht von `make baseline-freshness`. Das deckt sich mit der Zusage („`--verify` prüft als dritte Frage") und ist keine Abweichung.

## 4. `make gates` erreicht den Sensor — gemessen

`baseline-verify` ist die **erste** Voraussetzung von `gates` (`Makefile:299`), und `.NOTPARALLEL:` (`Makefile:52`) erzwingt sequenzielle Abarbeitung. Probe mit vollständig kopiertem `Makefile`+`a-check.mk` in der Sandbox, ein toter Alias eingelegt:

```
$ make gates
fetch-baseline-cache: verify .harness/baseline/v5.12.0 gegen SHA256SUMS
fetch-baseline-cache: toter Symlink .claude/rules/probe-tot.md — Ziel fehlt (Baseline-Bump?)
make: *** [Makefile:190: baseline-verify] Fehler 1
make gates EXIT=2
```

Kein weiteres Gate ist gelaufen (54 Zeilen Gesamtausgabe, alle aus `sha256sum -c`). Die Kette wird also **rot und zwar sofort** — Nebeneffekt, nicht im Plan benannt: ein toter Alias **verdeckt** jeden anderen Gate-Befund desselben Laufs.

## 5. Wirkung auf die Scan-Menge — gemessen

`git archive` von `a8f2c6b^`, `a8f2c6b` und `HEAD` nach `/tmp`, je `docker run --rm --network none -v …:/repo:ro d-check:latest`:

```
a8f2c6b^ : 597 Datei(en) geprüft, 0 Befund(e)   EXIT=0
a8f2c6b  : 598 Datei(en) geprüft, 0 Befund(e)   EXIT=0
HEAD     : 600 Datei(en) geprüft, 0 Befund(e)   EXIT=0
```

`a8f2c6b` legt **fünf** Dateien an: vier Symlinks + `harness/conventions/MR-055-symlink-als-pin-traeger.md`. Der Zuwachs ist **+1**. **Die vier Symlinks zählen nicht als Dateien** — genau wie zugesagt, keine Überraschung. (HEAD +2 gegenüber `a8f2c6b` stammt aus `ce984e6`/slice-179, nicht aus diesem Slice.)

Direkte Gegenprobe zur Asymmetrie-Behauptung aus §2 Punkt 3, an einer HEAD-Kopie:

```
A) echte Datei .claude/rules/probe-echt.md mit totem Link
   → .claude/rules/probe-echt.md:3  ./gibt-es-nicht.md  target-missing
   → 601 Datei(en), 1 Befund   EXIT=1        (echte Datei WIRD gescannt)
B) Wurzel-Datei probe-src.md (toter Link) + Symlink-Alias darauf unter .claude/rules/
   → probe-src.md:3  ./gibt-es-nicht.md  target-missing
   → 601 Datei(en), 1 Befund   EXIT=1        (der Alias erhöht weder Zahl noch Befund)
```

Die Zusage stimmt **exakt**. Ergänzend: `.harness/baseline/**` steht in `scan.ignore` (`.d-check.yml:15`) — die Ziele werden also auch **direkt** nicht geprüft; „ihr Inhalt ist nicht gate-geprüft" ist damit sogar breiter wahr als die Formulierung „über den Alias" nahelegt.

## 6. `MR-055` gegen `MR-021` / `MR-011` und die Index-Zeile

**Gegen `MR-021` — trägt, Feld gelesen.** `MR-021` §Geltungsbereich: *„**alle** Markdown-Links auf `.harness/baseline/<tag>/…` in der Live-Doku … die Menge bestimmt der Zensus der Bump-Prozedur, nicht diese Zeile."* Ein Symlink ist kein Markdown-Link; `MR-055`s Abgrenzung liest also das **Feld**, nicht den Titel, und die Behauptung „von `MR-021`s Zensus nicht erfasst" ist gegen dieses Feld korrekt. `MR-021` §Ersetzt-Baseline-Regel nennt `grundlagen-harness-dateien.md §Verzeichniskonvention`; `MR-055` nennt dort korrekt „keine — Nachtrag zu MR-021" und verweist auf dieselbe Kanon-Stelle als *bekannt, aber nicht ersetzt*. Sauber.

**Gegen `MR-011` — kein Konflikt, aber unerwähnt.** `MR-011` (`harness/conventions/done/`, `Aufgelöst durch: MR-012`) §Geltungsbereich sind `conventions.md` §Baseline und §Adoptierte Konventions-Quellen — die **Pin-Mechanik**, nicht der Sensor. `MR-055` fügt dem Träger der „MR-011-Kette" (`make baseline-verify`) eine dritte Frage hinzu, ohne `MR-011` zu nennen. Inhaltlich unbedenklich; die Folge ist aber, dass die vier Stellen, die `baseline-verify` als „MR-011-Kette" mit **zwei** Hälften beschreiben, nicht nachgezogen wurden (A-1…A-4).

**Index-Zeile — formal korrekt.** `harness/conventions.md:135`, in aufsteigender Reihenfolge nach `MR-052/053/054`, mit beiden Ankern (`#mr-055--…` und `#mr-055`) nach Hausform, Spalte 3 = Geltungsbereich, Spalte 4 = Ersetzt-Baseline-Regel. Alle vier Slice-Verweise auf `#mr-055` lösen auf (`doc-check` grün, `anchors` aktiv).

## 7. Abweichungen Zusage ↔ Zustand

**Spiegel-Klasse — der Sensor hat drei Fragen, vier Deklarationen sagen zwei** (`MR-025` „Spiegel vor dem Editieren auflisten"; kein Gate fängt es, `targets` prüft **Vorhandensein** von Targets, nicht Beschreibungstext):

- **A-1** `harness/README.md:90` — *„**Zwei Hälften, beide nötig:** `sha256sum -c` … die Manifest-Deckung …"*. Das ist die Tabelle, auf die `AGENTS.md` §4 ausdrücklich für „Details und Bindungen" verweist; sie steht jetzt im Widerspruch zum Code. Die Bindungs-Spalte nennt `MR-011`-Kette und `MR-021`, **nicht** `MR-055`. *(schwerwiegendste der vier — direkte Falschaussage, nicht nur Lücke)*
- **A-2** `AGENTS.md:373` — Zeile `make baseline-verify` nennt nur `sha256sum -c` + Manifest-Deckung.
- **A-3** `Makefile:167–176` (Kommentar: *„Beide Hälften nötig"*) und `Makefile:189` (Help-Text: *„Integrität + Manifest-Deckung"*).
- **A-4** `tools/harness/fetch-baseline-cache.sh:22–23` — der eigene Kopfkommentar des Skripts beschreibt `--verify` als *„gegen SHA256SUMS + Manifest-Deckung"*, während der Rumpf 20 Zeilen tiefer die dritte Frage stellt.

**Grenzen des Wächters, die `MR-055` §Grenze nicht nennt** (§Grenze nennt heute nur „Auflösung, nicht Ziel"):

- **A-5** `.claude/rules/` **ganz gelöscht** ⇒ still Exit 0 (P4). Der Verlust der gesamten Zustellung ist für den Sensor nicht von „gibt es hier nicht" unterscheidbar.
- **A-6** Alias auf ein **Verzeichnis** passiert (P5) — `readlink -e` gilt Verzeichnissen ebenso.
- **A-7** Alias in `.claude/rules/**sub/**` wird nicht gesehen (P10) — `for l in "$rules"/*` ist einstufig. Das kollidiert mit dem Wortlaut des Geltungsbereichs: *„Symlinks **unterhalb** von `.claude/rules/`"* liest sich rekursiv, die Implementierung ist es nicht.
- **A-8** **Dotfile**-Alias wird nicht gesehen (P9) — `*` fasst keine Punktdateien.
- **A-7b** SHA-Drift bricht vor der dritten Frage ab (P11); die Befunde akkumulieren nicht.
- **A-8b** Eine Symlink-**Schleife** wird zwar gemeldet, aber mit der Ursache *„Ziel fehlt (Baseline-Bump?)"* — die Meldung benennt eine Ursache, die im Schleifenfall falsch ist.
- **Umgekehrte Richtung:** `MR-055` §Geltungsbereich engt auf Aliase ein, *„deren Ziel in `.harness/baseline/<tag>/` liegt"*; geprüft wird tatsächlich **jeder** Symlink dort (P8 zeigt: ein Fremd-Alias wird geprüft, besteht aber). Implementierung breiter als Deklaration — harmlose Richtung, aber eine Ungenauigkeit.

**Beleg und Text:**

- **A-9** Die Beleg-Grenze zu DoD 2/3 ist im Slice **nicht** benannt (nur in der Commit-Botschaft); der Vorentwurf trug sie und hat sie beim Umschreiben verloren. Siehe §2.
- **A-10** Reichweiten-Grenze „Subagenten-Kontext" ist in keiner der fünf Nicht-Zusagen benannt; gemessen an diesem Lauf trifft sie genau die Rollen des Anlassfalls. Siehe §2.
- **A-11** DoD-Punkt 5 sagt *„(§2, fünf Punkte)"* und zählt dann **vier** auf — *„Anwesenheit ist nicht Wirkung"* fehlt in der Aufzählung (im Text §2 ist er vorhanden).
- **A-12** Wortwahl „Module": `regelwerk/` führt 26 **Dateien** = 17 `modul-*` + 8 `grundlagen-*` + `README.md`. „die übrigen **22 Module**" und „wer alle **26** einhängt" rechnen mit Dateien; nach `modul-*` wären es 17 gesamt und 13 übrig. Die Zeilenzahl 4787 ist korrekt (Dateisumme).
- **A-13** Kein wiederholbares Probe-Ziel für den neuen Sensor. Der Beleg lebt allein in der Commit-Botschaft `a8f2c6b`. Das Repo kennt für genau diese Lage die Gegenform — `make guard-probe`, dessen Makefile-Kommentar sagt: *„Ohne wiederholbare Proben wäre seine Zusage eine Erinnerung."* Die dritte Frage von `baseline-verify` hat weder Probe-Target noch Go-Test.
- **A-14** Formal: unnummerierte `## Zwei Träger wurden geprüft, einer davon verworfen` zwischen §1 und §2 (`structure` beanstandet es nicht, aber die Nummerierung des Slice-Templates ist damit unterbrochen).

**Was gegenüber der Zusage *besser* ist:** `MR-055` behauptet den Sensor; die Mutationsproben zeigen darüber hinaus, dass er **tragend** ist (P2: ohne ihn grün) und die **Richtung** seiner Bedingung tragend ist (P3), und P1c belegt den *tatsächlichen* Bump-Fall, nicht nur einen künstlichen toten Link. `working-tree-hash.sh` behandelt Symlinks bereits seit `MR-005` gesondert (`LINK %s -> %s`), der Gate-Nachweis bricht durch die Aliase also nicht — geprüft, dass das kein Nachtrag dieses Slice war.

## 8. Repo unverändert

```
$ diff repo-before.txt repo-after.txt
UNVERAENDERT: repo-before == repo-after
$ git status --porcelain
(leer)
$ sha256sum tools/harness/fetch-baseline-cache.sh
58e98015c02193b0401da08730a9333d47098e31f697c3b6f7a5d1eae0998948   (vorher == nachher)
$ sha256sum Makefile
98527cd6388b2a88c23ff2feb8cecc81f113f38669c2a9eb226c343ca7db82d2   (vorher == nachher)
alle vier Aliase: OK
```

Im Repo gefahren wurden ausschließlich lesende Kommandos (`git`, `grep`, `sed`, `wc`, `find`, `readlink`, `make baseline-verify`). `.harness/state/gates-passed.diffsha` wurde von mir **nicht** geschrieben — ich habe weder `make gates` noch `make record-gates` im Repo aufgerufen; die Datei trägt den Hash des aktuellen Baums und stammt aus der Parallel-Sitzung. Die temporären Abzüge unter `/tmp` sind entfernt, die Probe-Sandbox (`…/scratchpad/probe`) steht auf ihrem Ausgangszustand (`verify ok (51 Dateien, vollständig)`).

## 9. Empfehlung

Die vier inhaltlichen Zusagen des Slice halten der Messung stand, und der Sensor ist stärker belegt, als der Plan behauptet. **Vor der Closure zu schließen sind A-1 bis A-4** (die vier Deklarationen sagen „zwei Hälften", der Code hat drei Fragen — A-1 ist eine Falschaussage in der Tabelle, auf die `AGENTS.md` selbst verweist) und **A-9/A-10** (die zwei unbenannten Beleg- bzw. Reichweiten-Grenzen; A-10 trifft ausgerechnet die Rollen des Anlassfalls). **A-5 bis A-8** gehören als benannte Grenzen in `MR-055` §Grenze — oder, wo billig, in den Wächter (`shopt -s nullglob dotglob` plus `find`-statt-Glob deckt A-7/A-8 in einer Zeile). **A-11 bis A-14** sind Kleinigkeiten, aber A-13 ist die, die dieses Repo sonst konsequent anders entscheidet.
