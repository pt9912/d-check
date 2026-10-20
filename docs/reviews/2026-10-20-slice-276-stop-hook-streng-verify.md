# Verifikation — slice-276: Stop-Hook streng, HEAD-Stempel, Selbstprüfung

- **Rolle:** Verifier (Modul 11) — DoD, Plan und Spec gegen den Stand
- **Gegenstand:** `docs/plan/planning/done/wellenlos/slice-276-stop-hook-haertung-aus-ai-harness-init-v0100.md` §1, §2, §3 (samt der Plan-Änderung nach Runde 1), §4, §6; [MR-079](../../harness/conventions/MR-079-stop-hook-streng-mit-head-stempel.md); `SPEC-094`, `SPEC-103`
- **Range:** `2032fe1d..HEAD` (Spitze `565d13d9`); Review-Reports nur als Kontext
- **Datum:** 2026-10-20

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `make selbstpruefung` (selbst gefahren, Docker) | Exit 0; die vier Zeilen `STOP STRENG`, `STOP LOCKER`, `STOP VARIABLE`, `STOP BLOCK` stehen in der Ausgabe, danach `OK —`; der rote Commit-Versuch (ohne Kennung) fällt, der grüne geht durch |
| M2 | Gegenprobe Plan-Risiko 1: Hook von `2032fe1d` gegen aktuellen Hook im Wegwerf-Repo, Default (ohne Datei, ohne Variable), 17 Lagen | siehe Tabelle G; 16 von 17 deckungsgleich, eine bewusste Abweichung (kein Repo) |
| M3 | Aufhebungs-/Löschungs-Proben A–K gegen beide Hooks | siehe Tabelle L |
| M4 | Bewusstes Brechen an vier Kopien des Hooks im Scratchpad (Repo unberührt, `git status` leer) | siehe Tabelle B |
| M5 | `make gates` | Exit 0; `[gates] baseline-verify + … + planning-check green`, doc-check `1178 Datei(en) geprüft, 0 Befund(e)`, Coverage 94.70 % |
| M6 | `grep -rn "stop-gate-streng"` außerhalb `docs/reviews`, `planning/done` | nur Plan (Plan-Änderung, Risiko-Ausgang, Abgrenzung: historisch) und MR-079 (Upstream-Vergleich, Begründung: historisch); kein Hook, Skript, Spec oder README hält die Datei als Halter; `.harness/stop-gate-streng` existiert nicht, `.harness/stop-gate-locker` existiert nicht (`ls .harness`: nur `baseline`, `skills`, `state`) |

### Tabelle G — Gegenprobe im Default (alt = `2032fe1d`, neu = HEAD)

Neun geforderte Lagen und acht weitere; Ausgabe `rc` und `decision`.

| Lage | alt | neu |
|---|---|---|
| frischer Klon, sauber, kein Nachweis | 0 approve | 0 approve |
| Änderung ohne Nachweis | 0 block | 0 block |
| passender Nachweis, sauber | 0 approve | 0 approve |
| Änderung nach Nachweis, ohne Commit | 0 block | 0 block |
| Änderung nach Nachweis, mit Commit | 0 block | 0 block |
| dito, mit Commit und weiterer Änderung | 0 block | 0 block |
| Nachweisdatei unlesbar (Verzeichnis), sauberer Baum | 0 approve | 0 approve |
| Hash-Skript fehlt (gelöscht) | 0 block | 0 block |
| Hash-Skript fehlt (aus Index entfernt) | 0 block | 0 block |
| kaputter Index | 0 block | 0 block |
| `stop_hook_active: true` mit Änderung | 0 approve | 0 approve |
| `stop_hook_active: false` mit Änderung | 0 block | 0 block |
| Nachweis-Stempel gleich `HEAD`, Änderung (Kern der Lockerung) | 0 block | 0 block |
| Stempel gelöscht, Änderung | 0 block | 0 block |
| `stop-gate-streng` (alte Datei), Änderung | 0 block | 0 block (ohne Wirkung) |
| `diffsha` als Verzeichnis, Änderung | 0 block | 0 block |
| `diffsha` `chmod 000`, Änderung | 0 block | 0 block |
| **kein Repository** | **rc 1** (nicht blockierend) | **rc 2** (blockierend) |

Die einzige Abweichung ist die in MR-079 §Begründung benannte, und sie ist eine Verschärfung (Exit 1 ließ den Stop durch). Ohne `git` im `PATH`: alt Exit 127, neu `decision: block` mit Grund — wie in MR-079 behauptet. Eine stille Lockerung im Default liegt nicht vor.

### Tabelle L — Löschung, toter Symlink, fehlende Datei (Stempel gleich `HEAD`, ungedeckte Änderung)

| Lage | alt | neu |
|---|---|---|
| keine Datei | block | block |
| toter Symlink `stop-gate-locker` | block | block |
| Datei angelegt, dann gelöscht | block | block |
| alte Datei `stop-gate-streng` | block | block |
| Datei + `STOP_GATE_STRENG=0` / `=true` | block | approve (Variable sagt nichts, wie dokumentiert) |
| Verzeichnis `stop-gate-locker` | block | approve (`-e`, wie in MR-079 §Grenze dokumentiert) |
| Datei, Stempel gleich `HEAD`, Änderung | block | approve (der lockere Zweig, nur mit Datei) |
| Datei, neuer `HEAD` plus Änderung | block | block |

Die Zusage der Runde 2 trägt: Löschung, fehlende Datei und toter Symlink lockern nichts; nur eine vorhandene Datei lockert.

### Tabelle B — Bewusstes Brechen (Kopien im Scratchpad; Repo danach unverändert, `git status` leer)

| Zurückgenommen | Wegwerf-Lauf | Rot durch |
|---|---|---|
| Default-Zweig (`locker=0` → `locker=1`) und Datei-Prüfung (`-e …` → `locker=1`) | Lage „Änderung nach Nachweis, ohne Commit": block → **approve**; Lagen A–D (keine Datei, toter Symlink, gelöscht, alte Datei): block → **approve** | genau die stille Lockerung; entspricht dem Rot-Fall der Selbstprüfungs-Zeile `STOP STRENG` (Ausgang „blockiert" erwartet) und der Prüfung nach `rm -f …locker` am Ende |
| `trap` | Lage „kein Repository": rc 2 → **rc 1** | Exit 1 ließe den Stop durch; die Hook-Entscheidung der übrigen Lagen bleibt gleich |
| Aufhebung durch `STOP_GATE_STRENG=1` | Lage „Datei + `STOP_GATE_STRENG=1`": block → **approve** | entspricht der Selbstprüfungs-Zeile `STOP VARIABLE` |

Die Selbstprüfung selbst habe ich nicht gegen eine mutierte Hook-Kopie gefahren (sie klont `HEAD`; eine Mutation bräuchte einen Commit). Die Rot-Aussage für ihre Zeilen ist daher eine Entsprechung der Wegwerf-Lage, kein wörtlicher Lauf.

## DoD-Punkte

- **Hook und Nachweis — erfüllt.** `stop-require-gates.sh` trägt `trap` auf Exit 2, `head_wert`, den lockeren Zweig nur mit `-e .harness/stop-gate-locker`, `STOP_GATE_STRENG` genau `1` hebt auf; `record-gates.sh` schreibt `.harness/state/gates-passed.head`; `head_wert` ist in beiden Fassungen gleich (Kopplung im Kopf genannt). Die Datei gibt es nicht. Default streng belegt (Tabellen G, L; M1).
- **Selbstprüfung — erfüllt.** Stop-Hook-Zweige und Marker, belegt für d-check; `make selbstpruefung` grün, beide Ausgänge (blockiert/frei) beobachtet (M1). Die Selbstprüfung bricht, wenn der Klon die Locker-Datei trägt.
- **Entscheidung und Spiegel — erfüllt.** MR-079 (Index-Zeile in `conventions.md`, schärft MR-076) hält Übernahme, Abweichung vom Upstream und Nicht-Übernahme des lockeren Defaults fest. SPEC-094 (Hook-Modi, Stempel, Exit 2), SPEC-103 (Stop-Hook-Zweige), Historie-Zeile und die zwei Absätze/Zeilen in `harness/README.md` stimmen mit dem gemessenen Verhalten überein (Tabellen G/L, M1).
- **`make gates` grün — erfüllt** (M5).
- **Review durchgeführt — erfüllt** (zwei Runden, Reports liegen vor); Verifikation mit diesem Bericht.
- **Offen, keine Verletzung:** Beobachtungs-Register, Closure-Notiz §7, Risiko-Ausgänge §6 (Risiko 1 und 3 stehen noch „offen bis Closure"; Risiko 1 ist mit M2 belegbar, Risiko 3: die Selbstprüfung lief ohne neue Netz-/Docker-Annahmen über die bekannten hinaus grün), drei Paarungen, `git mv` nach `done/`.

## Plan-vs-Code (Abgrenzung §1)

Dateien im Diff: Hook, `record-gates.sh`, `selbstpruefung.sh`, MR-079 samt Index, SPEC-094/103 samt Historie, `harness/README.md`, Plan, zwei Review-Reports. Alle durch §3 oder die Plan-Änderung nach Runde 1 gedeckt. `slice-mv.sh`, `extract-command.awk`, `pretooluse-command-guard.sh`, Baseline-Pin und -Baum: im Diff unverändert. Kein lockerer Zustand im Repo: die Locker-Datei ist nicht angelegt, und die Selbstprüfung hält das fest. Abgrenzung eingehalten.

## Befunde

- **HIGH / MEDIUM:** keine.
- **LOW 1:** Die Selbstprüfung misst das `trap` nicht. Ohne `trap` bleibt `make selbstpruefung` grün; rot wird erst die Wegwerf-Lage „kein Repository" (rc 1). Die DoD behauptet das `trap` im Hook, nicht in der Selbstprüfung, und der Plan nennt als Rot-Beleg „oder Wegwerf-Repo-Lauf" — keine Verletzung, aber das `trap` ist nur durch diesen Bericht belegt (`tools/harness/selbstpruefung.sh`).
- **LOW 2:** `record-gates.sh` endet jetzt mit Exit 1, wenn `HEAD` nicht auflösbar und das Repo nicht eng als commitlos erkannt ist; das lässt `make gates` im letzten Schritt rot werden. Es steht im Skriptkopf, aber weder in SPEC-094 noch in MR-079 (`tools/harness/record-gates.sh`).
- **INFO:** Der lockere Zweig ist als Code im Hook vorhanden (Upstream-Mechanik, in MR-079 und Plan §3 so gewollt). Er greift nur mit einer Datei, die dieses Repo nicht führt; ein Anlegen wäre im Diff sichtbar, gehalten wird es von `make selbstpruefung` (kein Teil von `gates`) und vom Review — in MR-079 §Grenze benannt.
- **INFO:** `-e` lässt auch ein Verzeichnis dieses Namens lockern; dokumentiert, Probe G bestätigt.
- **INFO:** Plan §6 Risiko 2 führt „eingetreten und im Design behoben" als Ausgang; die Form passt zu den drei Wörtern.

## Urteil

DoD-Liefer-Punkte erfüllt, Plan, Spec und MR konsistent mit dem Code; der strenge Default ist deckungsgleich mit dem Hook von `2032fe1d` bis auf die dokumentierte Verschärfung „kein Repository". Keine DoD-Verletzung. Offen sind nur die Closure-Punkte.
