# Verifikation welle-97 — Closure-Trigger, Schritt 1

**Rolle:** Verifier (Modul 8, §Rollen-Sequenz für eine Welle: Verifier → Planner) ·
**Gegenstand:** Closure-Trigger aus
[welle-97](../plan/planning/done/welle-97-linelength-eingerueckte-tabellen.md) §3, geprüft in einem frischen Kontext ·
**Bezug:**
[`DC-FA-LLEN-001`](../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in),
[`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus),
[ADR-0114](../plan/adr/0114-linelength-modul.md),
[ADR-0116](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) · Datum 2026-11-24 · Modell-ID claude-sonnet-5-5.

Jeder Beleg stammt aus einem Lauf, den ich selbst gefahren habe. Das Trigger-Audit (Schritt 2) ist hier nicht
geprüft. Der Beleg geht über die Slice-DoDs hinaus; deren Verifikation steht in
[slice-293](2026-11-24-slice-293-linelength-eingerueckte-tabellen-verify.md) und
[slice-294](2026-11-24-slice-294-linelength-eingerueckte-tabellen-verify.md).

**Stand.** `make fullbuild` lief auf `5c8ded99` mit sauberem Arbeitsbaum. Während der Läufe kam `192cf345` auf
`main` (Planner-Kontext, parallel zu diesem Lauf: `AGENTS.md`, ein Beobachtungs-Zustand, ein Regel-Volltext und der
Trigger-Audit-Report — vier Dateien, nur Dokumentation). `git diff --name-only 5c8ded99 192cf345` nennt weder
`.go`, `Makefile`, `Dockerfile`, `go.mod`, `spec/` noch `docs/user/`; kein Image-Bestandteil hat sich bewegt. Probe,
Vergleiche, Probe-Fixtures und die Gates des Punktes 5 liefen auf `192cf345`, Arbeitsbaum sauber bis auf diesen
Report. Ich habe `make fullbuild` nicht wiederholt.

## Trigger-Punkte

| Punkt | Urteil | Beleg |
|---|---|---|
| 1 — alle Slices in `done/` | bestätigt | [slice-293](../plan/planning/done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md) und [slice-294](../plan/planning/done/slice-294-linelength-eingerueckte-tabellen-kern.md) liegen in `done/`; nur diese beiden tragen `**Welle:**` auf welle-97 (`grep -rl` über `docs/plan/planning`). `ls docs/plan/planning/{open,next,in-progress}`: `open/` und `next/` nur `.gitkeep`, `in-progress/` nur `roadmap.md`. |
| 2 — `make fullbuild` grün | bestätigt | Exit 0, Abschluss `[fullbuild] green — image-hash sha256:dee681826eea27a2163f4e4b3416e65ceef5140a5c740f7739fec9067fb9aac9`. Zeilen unten. |
| (a1) — Probe `REF=69c53e0f` | bestätigt | `blackbox-probe: byte-identisch über 20 Vergleiche (Vorher 69c53e0f, stdout/stderr/Exit getrennt, Kanarienlauf vorher und nachher).` Exit 0. |
| (a2) — `--enable linelength` Vorher/Nachher | bestätigt | 1176 Vergleiche; 0 Abweichungen ohne eingerückte Tabellenzeile; die eine gewollte Abweichung ist die einzige. |
| (b) — der Anlass am gebauten Image | bestätigt | Tabelle Probe → Ergebnis unten. |
| (c) — Trigger-Audit | nicht Gegenstand dieses Schritts | Der Report `2026-11-24-welle-97-trigger-audit.md` liegt seit `192cf345` vor; `welle-97-results.md` gibt es noch nicht (Planner, Schritt 3c). |

## Beleg 1 — `make fullbuild`

Exit 0 auf `5c8ded99`. Auszug aus dem Lauf (Wortlaut der Ausgabe):

- `fetch-baseline-cache: verify ok (54 Dateien, vollständig)`; `0 issues.`; `gesamt: 0 Befund(e)` (arch-check)
- `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%`
- `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`
- `image-test: OK — Akzeptanzkriterien erfüllt (DC-FA-DIST-001, DC-FA-CLI-007, DC-FA-CLI-008, DC-FA-CLI-013, DC-FA-MET-001, DC-FA-PSYNC-001, DC-FA-LLEN-001)`
- `bench: Median 1000 ms (Pass-Kriterium < 5000 ms — DC-QA-01)` (Läufe 982, 1030, 1000 ms)
- `58 Anforderung(en), 0 Waise(n).` (completeness-check)
- verify-closure-notes: `d-check: 1166 Datei(en) geprüft, 0 Befund(e)`
- Abschluss: `[fullbuild] green — image-hash sha256:dee681826eea27a2163f4e4b3416e65ceef5140a5c740f7739fec9067fb9aac9`

Anmerkung zum Hash: Direkt nach dem Lauf habe ich `docker image inspect d-check:latest --format '{{.Id}}'` nicht
mehr auf `dee68…` gesehen, sondern erst `e6fe9…`, später `4a7cc…`. `make blackbox-probe` und ein paralleles
`make gates` des Planner-Kontexts bauen das Image unter demselben Tag neu; das Daemon-Backend (containerd-Snapshotter)
führt für ein neu gebautes Image mit Attestation eine neue Kennung. Der Build dort zeigte die Runtime-Schicht
(`COPY --from=build /out/d-check`) als `CACHED`, das Binary hat sich also nicht bewegt. Als Hash des Laufs gilt der
von `make fullbuild` gedruckte; die Vergleiche unten liefen gegen das jeweils gerade getaggte `d-check:latest` desselben Quellstands.

## Beleg 2 — (a1) Probe

`make blackbox-probe REF=69c53e0f` auf `192cf345`, Exit 0. Gleich sind die Fixtures `ids`, `links`, `sauber`,
`targets` und der Eigenlauf `repo`, je Standard, `--json`, `--yaml`, `--doctor` (20 Vergleiche). Der Eigenlauf endet
vorher wie nachher mit Exit 0 — anders als gegen `v0.92.0` (Exit 2 vorher wegen `follow-renames`, siehe
[Verifikation zu welle-96](2026-11-24-welle-96-verifikation.md)); das Vorher-Image `69c53e0f` kennt den Schlüssel.

Was die Probe trägt, liest sich an `.d-check.yml`, wie der Plan es verlangt:

- Sie fährt `linelength` nicht: `grep -rl linelength tools/blackbox-probe/ tools/blackbox-probe.sh` findet nichts
  (Exit 1); die Fixtures sind `ids`, `links`, `sauber`, `targets`.
- Das Profil des Repos (`modules:` in `.d-check.yml`, Zeile 32) schaltet ein: `links, anchors, ids, matrix, codepaths,
  spans, hostpaths, versions, structure, diagrams, citations, file, pinsync`. **`structure` ist also im Eigenlauf
  `repo` enthalten; `planning` (und damit `planning.waves`) und `targets` sind es nicht** — `planning` läuft nur mit
  `--enable planning` (Kommentar zu Zeile 892 der Datei), `targets` hat das Fixture `targets`.
- Für `planning.waves` und `structure` im Lesen von Tabellenzeilen trägt deshalb nicht die Probe, sondern:
  `make planning-check` und `make verify-closure-notes` (beide mit `--enable planning --enable structure`) liefen
  auf dem neuen Image grün (Beleg 5), und der Kopplungs-Test aus [slice-294](../plan/planning/done/slice-294-linelength-eingerueckte-tabellen-kern.md).

## Beleg 3 — (a2) Vorher/Nachher mit `--enable linelength`

**Aufbau.** Vorher: `d-check:w97-vorher`, gebaut aus `git archive 69c53e0f` mit
`docker build --build-arg VERSION=0.0.0-dev --target runtime`, wie `tools/blackbox-probe.sh`. Nachher: `d-check:latest`.
Je Lauf `docker run --rm --network none -v <Fixture>:/repo:ro <Image> --enable linelength [--json|--yaml]`; stdout,
stderr und Exit getrennt gespeichert und mit `cmp` verglichen. Je Zelle ein eigenes Fixture-Verzeichnis mit
`.d-check.yml` (`modules: [links]`, `linelength:` mit Schalter und `max-chars`). Matrix je Fixture:
**7 Schalter-Stellungen** (Default · `exempt-tables: true` · `false` · `exempt-code-blocks: false` ·
`exempt-unbreakable: false` · alle drei `false` · alle drei `true`) × **`max-chars` 20, 40, 80, 120** × **3 Formen**
(Text, `--json`, `--yaml`) = 84 Vergleiche je Fixture; 14 Fixtures.

| Gruppe | Fixtures | Inhalt | Vergleiche | Abweichend |
|---|---|---|---|---|
| ohne eingerückte `\|`-Zeile (1 bis 3 Leerzeichen) | `user`, `maint`, `rules`, `adr` | Kopien von `docs/user/`, `docs/maintainer/`, `harness/rules/`, `docs/plan/adr/` (aus dem letzten die eine Datei mit einer solchen Zeile entfernt, 116 von 117); `grep -rEc '^ {1,3}\|'` über die Kopien: 0 | 336 | 0 |
| ohne eingerückte `\|`-Zeile, synthetisch | `s1` bis `s4` | Tabelle Spalte 0, vier Leerzeichen, Tab, `> \|`, `a \| b`, `- \|`, `10. \|`, U+00A0, Umlaute/Emoji, Fenced-Code, lange URL (`s1`); CRLF, `~~~`, unabgeschlossener Zaun, ohne Schluss-Newline, leere Datei (`s2`); nur Tabelle, nur Code, eingerückter Code mit `\|` (`s3`); Unterverzeichnisse und `CHANGELOG.md` (`s4`) | 336 | 0 |
| Zaun-Rand | `z1` | wie `s1`/`s2` mit einer `   \|`-Zeile (drei Leerzeichen) **im Zaun** — enthält eine eingerückte `\|`-Zeile, aber nicht außerhalb des Zauns | 84 | 0 |
| Spalte 0 / vier Leerzeichen | `d0`, `d4` | `1. Schritt:`, dann eine Tabelle mit null bzw. vier Leerzeichen und einer 162 Zeichen langen Datenzeile | 168 | 0 |
| **gewollt: ein bis drei Leerzeichen** | `d1`, `d2`, `d3` | dieselbe Tabelle unter `1. Schritt:` mit ein, zwei, drei Leerzeichen | 252 | **180** |
| | | | **1176** | **180** |

(1176 Vergleiche zu je stdout, stderr und Exit, also 3528 Einzelvergleiche.) Die 180 Abweichungen sind **genau** die
Zellen von `d1`, `d2`, `d3`, in denen `exempt-tables` eingeschaltet ist (Default, explizit `true`, `exempt-code-blocks:
false`, `exempt-unbreakable: false`, alle drei `true`: 5 Stellungen × 4 Grenzen × 3 Formen × 3 Fixtures). In den 72
Zellen mit `exempt-tables: false` (und alle drei `false`) sind Vorher und Nachher byte-gleich — „Nachher mit
`exempt-tables: false` wieder derselbe Befund wie Vorher“ ist belegt. Die Art der Abweichung, am Text nachgemessen:
stdout des Nachher-Laufs ist stdout des Vorher-Laufs ohne die Befunde der Zeilen 3 bis 5 (60 von 60 Text-Zellen), stderr
unterscheidet sich in der Bezugsmengen-Zeile (`7 Zeile(n) geprüft` gegen `4`) und der Befundzahl (je zwei Zeilen
vorher und nachher), der Exit wechselt von 1 auf 0; `--json` und `--yaml` tragen dieselbe Änderung in `findings`,
`summary.findingCount`, `summary.notes` und `exitCode`. Keine andere Art von Abweichung kommt vor.

**Die Vergleiche sind nicht leer.** Exit des Vorher-Images: 1116 mal 1, 60 mal 0 (die 60 sind `d0`, dessen Tabelle in
Spalte 0 steht und deshalb schon vorher ausgenommen ist). Länge des Vorher-stdout im Text-Lauf mit allen Ausnahmen aus
und `max-chars: 20`: `user` 4199 Zeilen, `adr` 12832, `maint` 328, `rules` 150, `s1` 12, `s2` 6, `s3` 5, `s4` 4, `z1` 18.
Kein Lauf endete mit einem Exit außerhalb von 0, 1, 2.

## Beleg 4 — (b) Der Anlass am gebauten Image

Wegwerf-Fixtures im Scratchpad, je eine Datei `liste.md` (sieben Zeilen: `1. Schritt:`, leer, Kopfzeile, Trennzeile,
**Datenzeile mit 56 Zeichen** — Zeile 5 —, leer, `2. Weiter`), `max-chars: 30`. Läufe mit
`docker run --rm --network none -v <Fixture>:/repo:ro … --enable linelength`. Die Bezugsmenge ist die Zeile
`linelength: 1 Datei(en), N Zeile(n) geprüft, …` auf stderr.

| Probe | Einrückung der Tabelle | Schalter | Image | erwartet | gemessen |
|---|---|---|---|---|---|
| 1 | 3 Leerzeichen | Default | Nachher | kein Befund | Exit 0, 0 Befunde, 4 Zeilen geprüft |
| 2 | 3 Leerzeichen | `exempt-tables: false` | Nachher | ein Befund | Exit 1, 1 Befund (`liste.md:5`), 7 Zeilen |
| 3 | 4 Leerzeichen | Default | Nachher | ein Befund | Exit 1, 1 Befund (Zeile 5), 7 |
| 4 | 4 Leerzeichen | `exempt-tables: false` | Nachher | ein Befund | Exit 1, 1 Befund (Zeile 5), 7 |
| 5 | ein Tab | Default | Nachher | ein Befund | Exit 1, 1 Befund (Zeile 5), 7 (Fixture mit `cat -A` als `^I\|` bestätigt) |
| 6 | `> \|` | Default | Nachher | ein Befund | Exit 1, 1 Befund (Zeile 5), 7 |
| 7 | `>   \|` (hinter `>` drei Leerzeichen) | Default | Nachher | ein Befund | Exit 1, 1 Befund (Zeile 5), 7 |
| 8 | 0 Leerzeichen | Default | Nachher | kein Befund (wie bisher) | Exit 0, 0 Befunde, 4 |
| 9 | 1 Leerzeichen | Default | Nachher | kein Befund | Exit 0, 0 Befunde, 4 |
| 10 | 3 Leerzeichen | `exempt-code-blocks: false` | Nachher | kein Befund (Zeile liegt nicht im Zaun) | Exit 0, 0 Befunde, 4 |
| 11 | 3 Leerzeichen, Zeile mit **einem** 60 Zeichen langen Wort | Default | Nachher | kein Befund | Exit 0, 0 Befunde, 4 (Tabelle ausgenommen, Wort-Zähler 0) |
| 12 | wie 11 | `exempt-tables: false` | Nachher | bestehende Regel: Wort ausgenommen | Exit 0, 0 Befunde, 7 Zeilen, `1 davon mit einem Wort über der Grenze ausgenommen` |
| 13 | wie 11 | `exempt-tables: false`, `exempt-unbreakable: false` | Nachher | ein Befund | Exit 1, 1 Befund (Zeile 5) |
| 14 | wie 11 | `exempt-unbreakable: false` (Tabellen-Ausnahme an) | Nachher | kein Befund | Exit 0, 0 Befunde, 4 |
| 15 | 4 Leerzeichen, Wortzeile | Default | Nachher | bestehende Regel: Wort ausgenommen | Exit 0, 0 Befunde, 7 Zeilen, 1 mit Wort ausgenommen |
| 16 | 4 Leerzeichen, Wortzeile | `exempt-unbreakable: false` | Nachher | ein Befund | Exit 1, 1 Befund (Zeile 5) |
| 17 | 3 Leerzeichen | Default | **Vorher** (`69c53e0f`) | ein Befund | Exit 1, 1 Befund (Zeile 5), 7 Zeilen |
| 18 | 1 Leerzeichen | Default | Vorher | ein Befund | Exit 1, 1 Befund (Zeile 5), 7 |
| 19 | 3 Leerzeichen | `exempt-tables: false` | Vorher | ein Befund | Exit 1, 1 Befund — gleich Probe 2 |
| 20 | 4 Leerzeichen · Tab · `> \|` · 0 Leerzeichen | Default | Vorher | wie Nachher | Exit 1, 1, 1 und 0 — je gleich Nachher (Proben 3, 5, 6, 8) |
| 21 | 3 Leerzeichen, Wortzeile | Default | Vorher | bestehende Regel | Exit 0, 7 Zeilen, 1 mit Wort ausgenommen (Nachher 4 Zeilen: die Zeile wird gar nicht mehr gemessen) |
| 22 | 3 Leerzeichen, Wortzeile | `exempt-tables: false`, auch mit `exempt-unbreakable: false` | Vorher | wie Nachher | Exit 0 bzw. Exit 1 — gleich Proben 12 und 13 |
| 23 | 3 Leerzeichen | Default, `--json` | Nachher / Vorher | `findings: []` / ein Befund | Nachher: `findings: []`, `findingCount: 0`, `exitCode: 0`, Note `4 Zeile(n)`; Vorher: ein Befund `line-chars-exceeded`, Exit 1 |
| 24 | 3 Leerzeichen | `exempt-tables: false`, `--json` | Nachher | ein Befund | `reason: line-chars-exceeded`, Zeile 5, Exit 1 |

Das Rot der Proben 2 bis 7 kommt aus der Einrückung und dem Schalter und nicht aus einer Zufallsregel: Probe 1 und
Probe 2 haben dieselbe Datei und unterscheiden sich nur im Schalter; Probe 17 zeigt, dass ohne die Änderung schon der
Default meldet; die Proben 3 bis 7 sind byte-gleich zum Vorher-Verhalten. Die Grenze, die der Wunsch **enger** zieht
als der Wortlaut des Issues („2 bis 5 Leerzeichen“), steht damit in der Probe: vier Leerzeichen bleibt ein Befund.

## Beleg 5 — Gates auf `192cf345`

| Kommando | Exit | letzte Zeile |
|---|---|---|
| `make gates` | 0 | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` (davor `d-check: 1306 Datei(en) geprüft, 0 Befund(e)`) |
| `make image-test` | 0 | `image-test: OK — Akzeptanzkriterien erfüllt (DC-FA-DIST-001, DC-FA-CLI-007, DC-FA-CLI-008, DC-FA-CLI-013, DC-FA-MET-001, DC-FA-PSYNC-001, DC-FA-LLEN-001)` |
| `make coverage-gate` | 0 | `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%` (die Zeile steht im Docker-Build-Log der Stage; das Target endet danach mit Exit 0) |
| `make completeness-check` | 0 | `58 Anforderung(en), 0 Waise(n).` |
| `make verify-closure-notes` | 0 | `d-check: 1167 Datei(en) geprüft, 0 Befund(e)` |

Zum Wortlaut von `image-test`: Fall (9) lautet dort `linelength — 40 grün, 41 rot, Codepunkte statt Bytes,
Fenced-Code/Tabelle/Wort über der Grenze ausgenommen und je einzeln schaltbar (:ro, fremder Benutzer, ohne Netz)`.

## Beleg 6 — Lifecycle

`ls docs/plan/planning/{open,next,in-progress}`: `open/` und `next/` je nur `.gitkeep`, `in-progress/` nur `roadmap.md`.
`ls docs/plan/planning/done | grep -E "29[34]|welle-97"` nennt `slice-293-…` und `slice-294-…`. Das flache
Wellendokument `welle-97-linelength-eingerueckte-tabellen.md` liegt noch unter `docs/plan/planning/` und wandert mit der
Closure (Planner). Die Roadmap führt den Zeiger auf welle-97 unter *Offene Wellen* (Zeile 35).

## Abweichungen vom Welle-Plan

1. **Welcher Teil des Profils in der Probe läuft (zu (a1)).** Der Plan sagt, die Closure lese an `.d-check.yml`, ob das
   Profil `structure` und `planning.waves` einschaltet. Gelesen: `structure` ja (Zeile 32); `planning` nicht in
   `modules:`, nur per `--enable planning`; `targets` nicht im Profil des Repos, aber im Fixture `targets`. Der
   Eigenlauf `repo` der Probe fährt also `structure`, nicht `planning.waves`. Das ist keine Abweichung vom Plan, sondern
   die Antwort auf seine Frage; die Ergebnis-Notiz nennt es. Für `planning.waves` trägt `planning-check` (grün).
2. **Zwei Stände.** `make fullbuild` auf `5c8ded99`, der Rest auf `192cf345`; der Unterschied ist ein reiner
   Dokumentations-Commit des Planner-Kontexts (siehe Stand).
3. **Fixtures.** Der Plan nennt für (a2) ein Fixture, eine Kopie von `docs/user/`. Ich habe 14 gefahren; `z1` trägt
   bewusst eine `   |`-Zeile **im Zaun** und erfüllt das Merkmal „ohne eingerückte `|`-Zeile“ nicht streng — er steht
   als eigene Gruppe, mit 0 Abweichungen, und fließt nicht in die 672 Vergleiche der ersten beiden Gruppen ein.
4. **Image-Kennung.** Der gedruckte Hash `dee68…` ist nicht mehr per `docker image inspect` abrufbar, weil
   `d-check:latest` im Verlauf neu getaggt wurde (siehe Beleg 1). Das trifft die Reproduzierbarkeits-Aussage des
   Plans nicht (Binary-Schicht `CACHED`), aber ein späterer Leser findet den Hash nicht mehr im Daemon.

Keine weitere Abweichung: die Slice-Zahl, die Reihenfolge und die gewollte Abweichung (eine eingerückte Tabellenzeile mit
ein bis drei Leerzeichen; null Leerzeichen war schon vorher ausgenommen und weicht erwartungsgemäß nicht ab) entsprechen dem Plan.

## Grenzen dieser Verifikation

- Die 14 Fixtures sind eine Auswahl, kein Beweis; die Gleichheit im ganzen Raum tragen die Tests und die Mutationen
  aus der Slice-Verifikation von [slice-294](2026-11-24-slice-294-linelength-eingerueckte-tabellen-verify.md).
- Der Bestand des Melders von [Issue #20](https://github.com/pt9912/d-check/issues/20) (41 von 85 Befunden) liegt
  außerhalb dieses Repos und ist von hier nicht gemessen; Probe 3 und 4 sind der Beleg, dass vier Leerzeichen und mehr
  Befunde behalten.
- Die Probe `blackbox-probe` fährt `linelength` nicht; (a2) ist der Beleg für das Modul, die Probe der für alle anderen.
- Trigger-Audit, die Paarungen und `welle-97-results.md` sind nicht Gegenstand dieses Schritts.

## Ergebnis

Die Punkte 1, 2, (a1), (a2) und (b) sind bestätigt; (c) liegt beim Planner. Die eine gewollte Abweichung — eine
`|`-Zeile mit ein bis drei Leerzeichen außerhalb von Fenced-Code wird bei eingeschaltetem `exempt-tables` nicht mehr
gemessen — ist die einzige: 180 von 1176 Vergleichen, alle in `d1` bis `d3` mit eingeschaltetem Schalter. Image-Hash des
Laufs `make fullbuild`: `sha256:dee681826eea27a2163f4e4b3416e65ceef5140a5c740f7739fec9067fb9aac9`. Offen für die
folgenden Closure-Schritte (Planner): `welle-97-results.md` mit der Tabelle Probe → Ergebnis (Beleg 4) und dem Hinweis
zu Abweichung 1 und 4.
