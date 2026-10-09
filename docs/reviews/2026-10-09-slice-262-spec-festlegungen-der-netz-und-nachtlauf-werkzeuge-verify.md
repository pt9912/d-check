# Verifikation — slice-262: Festlegungen der Netz- und Nachtlauf-Werkzeuge in die Spezifikation

**Rolle:** Verifier (Baseline-Regelwerk `modul-11-verification.md`). Geprüft wird gegen Plan und DoD, nicht gegen Diff und Hard Rules (das tut der Reviewer).
**Gegenstand:** slice-262, committeter Stand `4f6803ea` (Commits `65a63366`, `5231f227`, `90344f3e`, `bc660ea4`, `40f9c434`, `4f6803ea`)
**Geprüfte DoD-Punkte:** §2 Punkt 1 und 2. Punkt 3 (Review/Verifikation) und 4 (Closure) folgen.
**Lauf-Ort:** frischer Klon von `4f6803ea` im Scratchpad. Nur bash, make und Docker; kein Host-Interpreter.
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09

**Abgrenzung des Stands:** Im Arbeitsbaum von `dem Arbeitsbaum` lag während des Laufs eine
**nicht committete** Änderung an `spec/spezifikation.md` (SPEC-102 „ohne oder mit nur einem Commit“,
SPEC-103 „schon die beiden Commit-Versuche bauen das Image“), vermutlich die Antwort auf R2. Sie ist
**nicht** Gegenstand dieser Verifikation. Beide Aussagen stimmen aber mit den Messungen unten überein
(V-9: ohne Commit gibt es kein `HEAD~1`; V-11: die Commit-Versuche bauen `d-check:latest`).

## Sensor-Belege (selbst gefahren)

| # | Kommando | Ergebnis |
|---|---|---|
| G | `make gates` im frischen Klon | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`, **Exit 0**. semgrep 0 Befunde, d-check 0 Befunde. |
| V-1 | `pin-freshness.sh --compare x 1.2.3 1.2.3` / `x 1.2.3 1.2.4` / `x 1.2.3 ''` | ok, Exit 0 / `VERALTET` auf stderr, **Exit 3** / SKIP, Exit 0 |
| V-2 | `pin-freshness.sh --compare x v1.2.3 1.2.3` | **`VERALTET`, Exit 3**, nicht ok (siehe B-1) |
| V-3 | `pin-freshness.sh` ohne Modus mit `PINNED=1` / ohne `PINNED` · `--github` und `--digest` mit `PATH` ohne curl/docker · `--digest` mit Pin `sha256:abc` | Exit 1 / SKIP, Exit 0 (siehe B-2) · SKIP, Exit 0 · SKIP „Pin ist kein Digest“, Exit 0 |
| V-4 | `make runtime-base-digest` in einer Kopie, deren `Dockerfile` keine `distroless`-`FROM`-Zeile hat | „PINNED ist leer — SKIP“, Exit 0 (SPEC-099: keine passende `FROM`-Zeile ⇒ SKIP) |
| V-5 | **Netz:** `make freshness-golangci`, `freshness-go`, `checkout-pin-freshness`, `runtime-base-digest`, `trivy-digest` | alle ok, Exit 0. `golangci-lint ok — Pin 2.14.0`: der Pin `v2.14.0` erscheint ohne `v`, das Präfix fällt also im `--github`-Zweig auf beiden Seiten weg (SPEC-098) |
| V-6 | `tools/image-scan.sh --selftest` · `IMAGE_SCAN_REFS="   "` · `PATH` ohne docker (mit und ohne `IMAGE_SCAN_PLATFORMS`) | 0 Fehlschläge, Exit 0 · Exit 2 · beide Exit 2 „GESCHEITERT … UNBEKANNT, nicht gruen“ (fail-closed, SPEC-097) |
| V-7 | `nightly-state.sh --selftest` · `--parse` mit `success` (stderr verworfen) · `--parse` mit `timed_out` · `PATH` ohne curl · `NIGHTLY_WORKFLOWS="   "` · `--parse /nonexist` | 6/6 ok, Exit 0 · `gruen` auf stdout · `ROT (timed_out)` + Lauf-Störung auf stderr, Exit 0 · SKIP, Exit 0 · keine Ausgabe, Exit 0 · Exit 0 |
| V-8 | **Netz:** `make baseline-freshness`, `make nightly-state` | Currency OK + Content OK, Exit 0 · upstream-drift `gruen`, image-scan `ROT (failure)` mit festem Hinweis, Exit 0 |
| V-9 | `history-range-guard.sh` in einem Wegwerf-Repo: ein Commit, `HEAD~1..HEAD` · zwei Commits, `HEAD..HEAD` · `HEAD~1..HEAD` · `--staged` mit leerem Index · `nope..HEAD` | **Exit 2** · **Exit 1** („aufloesbar, aber LEER“) · Exit 0 · Meldung, Exit 0 · Exit 2 |
| V-10 | `fetch-baseline-cache.sh --check-latest` mit `PATH` ohne `curl` (alle anderen Programme aus `/usr/bin` per Symlink) · `make baseline-freshness` ebenso · Tag `vX` | „'curl' nicht gefunden“, **Exit 1** · make Exit 2 · Exit 1 |
| V-11 | Currency offline mit einem Stub-`curl` im Scratchpad (bash-Skript, liefert eine feste Release-Liste) | Pin neuester ⇒ Exit 0; neuerer Tag ⇒ Exit 3 „NEUER RELEASE“; Pin fehlt in Liste ⇒ Exit 3 „unbestimmt“; leere Antwort ⇒ SKIP, Exit 0 |
| V-12 | Klon der Tiefe 1 (`git clone --depth 1 file://…`), `make hooks`, `git commit --allow-empty -m "probe slice-262"` | Hook baut `d-check:latest`, dann `history-range-guard: Range 'HEAD~1..HEAD' ist NICHT aufloesbar`, `make: *** [Makefile:409: trace-check] Fehler 2`, Commit abgebrochen (git Exit 1), 1 Commit im Log. Bestätigt den Satz, den R1 in SPEC-102 nachgetragen hat (F-1). |
| V-13 | `selbstpruefung.sh` mit `SELBSTPRUEFUNG_MSG_ROT="falsch rot slice-262"` · mit `SELBSTPRUEFUNG_AKTIVIERUNG=true` | FEHLER „geht durch“, **Exit 1** · FEHLER „core.hooksPath im Klon weiter leer“, **Exit 1** |
| V-14 | `make selbstpruefung` (positiver Gesamtlauf) | ROT: der Commit ohne Kennung fällt am Träger (Exit 1), `HEAD` bleibt stehen. GRÜN: der Commit mit Kennung geht durch, `HEAD` trägt seinen Betreff. GATE: `make gates` im Klon green (coverage-gate 94.80 % ≥ 93 %). `selbstpruefung: OK`, **Exit 0** |

### Bewusstes Brechen (Mutation in einer Kopie im Scratchpad, nicht im Repo)

| Mutation | erwartet | gemessen |
|---|---|---|
| `image-scan.sh`: `grep -v '^unknown/'` aus `plattformen_aus_antwort` entfernt | Selbsttest rot an der Attestations-Regel aus SPEC-097 | `FAIL Attestation ausgenommen … war: [linux/amd64 unknown/unknown ]`, `FAIL nur Attestations`, 2 Fehlschläge, Exit 1 |
| `nightly-state.sh`: `null\|"")` → `"")` | Selbsttest rot bei „noch kein Ausgang ⇒ SKIP“ (SPEC-101) | `FAIL laeuft gerade (null) … war: ROT (null)`, Exit 1 |
| `history-range-guard.sh`: Leer-Prüfung `"$count" -eq 0` → `false` | leere Range wird still grün (der Fall, den SPEC-102 ausschließt) | `HEAD..HEAD aufgeloest, 0 Commit(s) — OK`, Exit 0 statt 1 |
| `selbstpruefung.sh` über seine Marker (V-13) | ein Hook, der nicht aufhält, bzw. eine Aktivierung ohne Wirkung ⇒ Exit 1 | beide Exit 1 aus dem richtigen Grund |

Jede Mutation ergibt die Fehlermeldung, die zur jeweiligen SPEC-Aussage passt, nicht bloß irgendeinen Fehler.

## DoD-Punkt 1 — je Werkzeug mit eigener Festlegung ein §7-Eintrag, am Code geprüft

**Bestätigt.** `spec/spezifikation.md` §7 trägt SPEC-097 bis SPEC-103 (Zeilen 3813–3819), dazu zwei
Historie-Zeilen. Die Abdeckung entspricht §1 des Plans: `image-scan`, Versions-Achsen samt der drei
Action-Pins, Digest-Achsen, `baseline-freshness`, `nightly-state`; außerdem die drei beim Beanspruchen
zu prüfenden Werkzeuge: `history-range-guard` (SPEC-102) und `selbstpruefung` (SPEC-103) haben einen
Eintrag, `baseline-probe` keinen. Begründet ist das in §8 („Selbsttest der Alias-Frage aus SPEC-092“).
Die Gegenprobe bestätigt es: `--selftest` fährt nur `check_aliases`.

Jeder Eintrag wurde gegen den Code gelesen, und mindestens eine Aussage je Eintrag lief ausführend:
SPEC-097 V-6 + Mutation, SPEC-098 V-1/V-3/V-5, SPEC-099 V-3/V-4/V-5, SPEC-100 V-8/V-10/V-11,
SPEC-101 V-7/V-8 + Mutation, SPEC-102 V-9/V-12 + Mutation, SPEC-103 V-13/V-14. Die Nachtlauf-Aussagen
(„alle außer `freshness-trivy`“ bzw. „außer `trivy-digest`“) stimmen mit den 13 Schritten von
`upstream-drift.yml`. Die Aussage „Nachtlauf liest die Ausgabe, wird bei Befund wie bei Scheitern rot“
stimmt mit `image-scan.yml` (grep auf `GESCHEITERT`, dann auf `behebbare CRITICAL/HIGH-Befunde`, dann `rc`).
Referenz-Richtung (`AGENTS.md` §3.4): die sieben Zeilen tragen kein ADR-, Slice-, Wellen-, MR- oder
Hash-Token (grep leer). Den `matrix`-Teil von `doc-check` deckt G.

Kein Eintrag widerspricht dem Code. Zwei Präzisierungen ohne DoD-Verletzung stehen unter B-1 und B-2.

## DoD-Punkt 2 — Sensor-Dateien verlinken die Kennung statt Schwelle und Randform zu führen; `make gates` grün

**Bestätigt, mit einem LOW-Befund (B-3).** `make gates` ist im frischen Klon grün (G). Die sechs
geänderten Sensor-Dateien verlinken ihre Kennung jeweils in Vertrag oder Ausgänge **und** in der
Bindung. `trace-check.md` verlinkt SPEC-102, und der Gate-Index nennt SPEC-098/099/102/103 an den
Zeilen ohne eigene Sensor-Datei. Die Exit-Tabelle von `image-scan.md` ist entfernt, die übrigen Dateien
führen weder eine Exit-Tabelle noch einen Vorrang. Die gemessenen Abweichungen aus §8 des Plans sind
nachgezogen: fünf Versions-Achsen, sechs Digest-Achsen, keine Trivy-Achse im Nachtlauf (Sensoren,
Gate-Index und Kopf von `upstream-drift.yml`), `curl` fehlt ⇒ Exit 1 (`baseline-freshness.md`,
`docker-make-only.md`, Workflow-Kopf), `nightly-state` entscheidet nicht über „planmäßig“. Die
Selbsttest-Zahlen in `image-scan.md` (7/4/6 Proben) stimmen mit dem Skript.

## Befunde

### B-1 — INFO — `--compare` normalisiert nicht; der Präfix-Abgleich aus SPEC-098 ist netzlos nicht prüfbar

- **pfad:** `tools/harness/pin-freshness.sh` (Zweig `--compare`); `harness/sensors/freshness-go.md` · „Netzlos prüfbar über `--compare` …; ohne diesen Einstieg wäre die Semantik nur mit Netz zu prüfen“
- **befund:** `--compare x v1.2.3 1.2.3` ergibt `VERALTET`, Exit 3 (V-2). Das `v`/`go`-Strippen steht im `--github`- bzw. `--godev`-Zweig, nicht in `compare()`. SPEC-098 ist dabei **richtig**, denn es ordnet das Strippen den Achsen zu und sagt über `--compare` nichts. Die Sensor-Datei verspricht aber „die Semantik“ netzlos, und die Präfix-Hälfte gehört nicht dazu: Sie ist nur mit Netz belegbar (V-5). Auch die im Auftrag genannte Erwartung „`v1.2.3` gegen `1.2.3` ⇒ ok“ trifft auf den Testeinstieg nicht zu.
- **wirkung:** keine DoD-Verletzung. Eine Grenze der Netzlos-Zusage, die die Sensor-Datei nicht nennt.

### B-2 — INFO — SPEC-098 „nur ein fehlender oder unbekannter Modus ist Exit 1“ gilt nur bei gesetztem Pin

- **pfad:** `spec/spezifikation.md` SPEC-098; `pin-freshness.sh` (Leer-Pin-Prüfung steht vor dem `case "$mode"`)
- **befund:** Ohne `PINNED` endet auch ein fehlender oder unbekannter Modus mit SKIP, Exit 0 (V-3). Über `make` tritt das nicht auf, weil jedes Target einen Modus setzt. Die Reihenfolge der beiden Regeln nennt der Eintrag nicht.
- **wirkung:** keine. Nur ein direkter Skript-Aufruf ist betroffen.

### B-3 — LOW — `image-scan.md` nennt den Vollbericht „nie fällt“; Code und SPEC-097 lassen ihn bei einem Trivy-Fehler scheitern

- **pfad:** `harness/sensors/image-scan.md` · „ein Vollbericht über alle Schweregrade, der nie fällt“; Gegenstand `tools/image-scan.sh` (Vollbericht-Lauf: `if ! trivy … --format table` ⇒ `GESCHEITERT`, `errored=1`)
- **befund:** SPEC-097 sagt richtig „der nur bei einem Trivy-Fehler scheitert“. Die Sensor-Datei, die dieser Slice angefasst hat, behauptet weiter „nie“. Der Satz stammt nicht aus diesem Slice, steht aber in einer Datei, für die DoD-Punkt 2 verlangt, dass sie die Randform der Kennung überlässt, statt eine eigene, abweichende zu führen. Gemeint ist vermutlich „fällt nie wegen Befunden“, doch so steht es nicht da.
- **verifizierbar:** ja. Lesbar am Code. Ausführend nur mit Netz und einem gezielten Trivy-Fehler, nicht gefahren.
- **vorschlag:** „…, der nicht wegen Befunden fällt“ oder der Verweis auf SPEC-097.

### B-4 — INFO — der Skript-Kopf von `check_latest` sagt weiter pauschal „Ausfall → SKIP je Teil“

- **pfad:** `tools/harness/fetch-baseline-cache.sh` · „KEIN fail-closed (Ausfall → SKIP je Teil)“, zwei Zeilen über dem `curl`-Riegel mit Exit 1
- **befund:** Diesen Rest benennt schon R1 F-3. Er liegt außerhalb der Plan-Abgrenzung (§3: nur Text in Doku, keine Werkzeug-Datei) und ist darum keine DoD-Verletzung. Er bleibt aber der letzte Ort, der der korrigierten Aussage widerspricht.

### Nicht gemessen

- Ob Trivy bei einer per `IMAGE_SCAN_PLATFORMS` genannten, im Index fehlenden Plattform scheitert oder ausweicht (R1 F-5). `image-scan.md` sagt seit `4f6803ea` ausdrücklich, dass das nicht gemessen ist, und der Ausgang Exit 2 gilt in beiden Fällen. Ein voller Trivy-Lauf mit Netz wurde nicht gefahren.

## Negativbefunde

- **Plan-Treue:** Der Diff berührt nur, was §3 des Plans in seiner Fassung nach R1 nennt: Spezifikation, Sensor-Dateien, Gate-Index, Kopf von `upstream-drift.yml` (nur Kommentar, die 13 Schritte sind unverändert), `docker-make-only.md`, außerdem `trace-check.md`. Kein Skript, kein Makefile, keine Workflow-Logik.
- **Rückführungs-Bedingung (§4):** Kein Werkzeug brauchte mehr als eine §7-Zeile. Die Bedingung ist nicht eingetreten.
- **Spiegel-Suche** nach den alten Formulierungen („vier Versions-Achsen“, „fünf Achsen“, „Werkzeug-Ausfall“, „JEDER Pin“) in lebender Doku, Workflows, Skripten und Makefile: keine Treffer außer den korrigierten Stellen, der eingefrorenen ADR-0067 („fünf Digest-Achsen“ zum Zeitpunkt ihrer Annahme) und der Ausgabezeile von `nightly-state.sh` (wörtlicher Hinweistext, den SPEC-101 als „festen Hinweis“ beschreibt).

## Verdikt

**DoD-Punkte 1 und 2: erfüllt.** `make gates` ist im frischen Klon grün. Alle sieben Einträge stimmen in
jeder ausführend geprüften Aussage mit dem Code überein. Bei den drei Selbsttests und beim Range-Wächter
hat das bewusste Brechen gezeigt, dass sie aus dem richtigen Grund rot werden. Es gibt einen LOW-Befund (B-3)
und drei INFO-Befunde. Keiner blockiert die Closure. B-3 ist ein Einzeiler und sollte vor der Closure mitgehen.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 1 | 3 |
