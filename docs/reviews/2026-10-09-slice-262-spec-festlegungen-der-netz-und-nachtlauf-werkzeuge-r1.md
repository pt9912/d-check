# Review R1 — slice-262: Festlegungen der Netz- und Nachtlauf-Werkzeuge in die Spezifikation

**Review-Art:** Code (Diff gegen Plan, ADRs, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-262, Commit `90344f3e` (`git show HEAD`)
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0 (`94f798ef`)
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-262 (§1 Ziel und Abgrenzung, §3 Plan, §8 Messung beim
Beanspruchen); `spec/spezifikation.md` §7 (`SPEC-092` bis `SPEC-103`);
[ADR-0011](../plan/adr/0011-digest-pins-build-gate-images.md),
[ADR-0066](../plan/adr/0066-cve-scan-gegen-das-publizierte-image.md);
[`MR-025`](../../harness/conventions.md#mr-025), [`MR-053`](../../harness/conventions.md#mr-053);
`AGENTS.md` §3.4, §3.7, §5 (Regeln 13, 15), §6 Schritt 4.
Gelesener Code: `tools/image-scan.sh`, `tools/harness/pin-freshness.sh`,
`tools/harness/fetch-baseline-cache.sh` (`--check-latest`, `--selftest`),
`tools/harness/nightly-state.sh`, `tools/harness/history-range-guard.sh`,
`tools/harness/selbstpruefung.sh`, die Makefile-Rezepte der betroffenen Targets
(samt `image-digest-axis`, `action-pin-axis`, `trace-check`, `adr-check`),
`.githooks/commit-msg`, `.github/workflows/upstream-drift.yml`, `image-scan.yml`.

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| A | Wegwerf-Repo mit einem Commit (Scratchpad): `bash tools/harness/history-range-guard.sh HEAD~1..HEAD` | Exit 2, „Range 'HEAD~1..HEAD' ist NICHT aufloesbar" |
| B | `make -n trace-check MSGFILE=/dev/null` | erste Rezeptzeile `bash tools/harness/history-range-guard.sh HEAD~1..HEAD` — der Wächter läuft auch im Message-Modus |

## Findings

### F-1 — MEDIUM — SPEC-102 verschweigt, dass der Wächter auch im Message-Modus des `commit-msg`-Hooks läuft

- **quelle:** `AGENTS.md` §5 Regel 13; Skill-Prüffrage 18
- **pfad:** `spec/spezifikation.md` · „Vorlauf der history-lesenden Targets (`trace-check`, `adr-check`)"; Gegenstand `Makefile` · „@bash tools/harness/history-range-guard.sh $(if $(RANGE),$(RANGE),HEAD~1..HEAD)" (Rezept von `trace-check`)
- **befund:** Das Rezept von `trace-check` ruft den Wächter unbedingt, auch mit `MSGFILE` — dem Modus des `commit-msg`-Hooks, der keine Historie liest. In einem Repo mit genau einem Commit oder in einem Klon der Tiefe 1 endet damit jeder Commit über den installierten Hook mit Exit 2, obwohl die Botschaft eine Kennung trägt (Messung A, B). SPEC-102 nennt die Unauflösbarkeit des Defaults, aber nicht, dass sie den Commit-Pfad sperrt; die Grenze steht nur im Code.
- **verifizierbar:** ja — Wegwerf-Repo mit einem Commit, `make hooks`, zweiter Commit mit Kennung
- **klasse:** grenze-nur-im-code

### F-2 — MEDIUM — Kopf von `upstream-drift.yml` behauptet „JEDER Pin" und zählt die Achsen falsch

- **quelle:** [`MR-025`](../../harness/conventions.md#mr-025) (Ableiter für eine Wortlaut-Präzisierung: `grep` nach dem alten Wortlaut); `AGENTS.md` §5 Regel 15; Beobachtung `BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen` (vom Plan selbst gesichtet)
- **pfad:** `.github/workflows/upstream-drift.yml` · „Gewacht wird JEDER Pin, den" sowie „semgrep, a-check) und die zwei Action-Pins."
- **befund:** SPEC-098/099 und die zwei Sensor-Dateien sagen jetzt richtig, dass `freshness-trivy` und `trivy-digest` nicht im Nachtlauf laufen; der Kopf des Nachtlaufs behauptet weiter, er wache **jeden** Fremd-Pin, nennt zwei statt drei Action-Pins und zwei Toolchain-Versionen ohne `semgrep`/`a-check`. Nicht gewacht sind außerdem die Builder-Images von `release.yml`. Failure-Szenario: wer den Nachtlauf-Kopf als Deckungsaussage liest, hält den Scanner-Pin des CVE-Nachtlaufs für überwacht, und er altert unbemerkt. Die Messung im Plan §8 führt diesen Spiegel nicht.
- **verifizierbar:** nein (Prosa gegen Workflow-Schritte; kein Gate)
- **klasse:** spiegel-nicht-nachgezogen

### F-3 — LOW — „fail-open" bei fehlendem `curl` steht an zwei Stellen weiter

- **quelle:** [`MR-025`](../../harness/conventions.md#mr-025)
- **pfad:** `harness/rules/docker-make-only.md` · „Die ersten drei sind fail-open, das vierte nicht"; `.github/workflows/upstream-drift.yml` · „(Netz-, Werkzeug- oder Manifest-Ausfall => SKIP je Teil)"
- **befund:** Beide Sätze stehen im Zusammenhang der Werkzeug-Erwartung (`curl`); `fetch-baseline-cache.sh --check-latest` endet ohne `curl` mit Exit 1, wie SPEC-100 und `baseline-freshness.md` jetzt sagen. Der Skript-Kopf von `fetch-baseline-cache.sh` ist im Teil (B) richtig gelesen, sein Abschnitt zu `check_latest` („KEIN fail-closed (Ausfall → SKIP je Teil)") ebenso pauschal.
- **verifizierbar:** nein
- **klasse:** spiegel-nicht-nachgezogen

### F-4 — INFO — SPEC-103 nennt zwei Grenzen des Laufs nicht

- **quelle:** Maintainability
- **pfad:** `spec/spezifikation.md` · „in einem Wegwerf-Klon (`file://`, ohne Netz, mit Docker)"
- **befund:** Der Klon trägt den Stand von `HEAD`, nicht den Arbeitsbaum (das Skript gibt es aus: „der Klon traegt den Stand von HEAD") — ein grüner Lauf sagt nichts über ungespeicherte Änderungen. „ohne Netz" gilt dem Klon-Transport; `make hooks`/`make gates` im Klon bauen Images und ziehen fehlende Basis-Images aus dem Netz.
- **verifizierbar:** nein
- **klasse:** grenze-nicht-benannt

### F-5 — INFO — Wer eine per `IMAGE_SCAN_PLATFORMS` genannte, fehlende Plattform fängt, ist nicht gemessen

- **quelle:** Maintainability
- **pfad:** `harness/sensors/image-scan.md` · „das fängt erst dieser Nachweis"; Commit-Botschaft · „prueft erst der Architektur-Nachweis"
- **befund:** SPEC-097 sagt vorsichtig „prüft dann erst der Scan". Ob Trivy bei einer Plattform, die der Index nicht führt, mit Fehler endet (→ „GESCHEITERT" im Nachweis-Lauf) oder still auf eine andere Variante ausweicht (→ Architektur-Vergleich), ist im Slice nicht gemessen; beide Wege enden mit Exit 2, die Aussage hält im Ausgang, nicht im Mechanismus.
- **verifizierbar:** ja, mit Netz — `IMAGE_SCAN_PLATFORMS=linux/s390x bash tools/image-scan.sh`
- **klasse:** mechanismus-unbelegt

### F-6 — INFO — Satzfehler in zwei neuen Grenzen-Punkten

- **quelle:** Maintainability
- **pfad:** `harness/sensors/freshness-go.md` und `harness/sensors/runtime-base-digest.md` · „wird nur gemeldet, wer die Achse"
- **befund:** „wer" statt „wenn man"; der Satz ist so nicht lesbar.
- **verifizierbar:** nein
- **klasse:** wortlaut

### F-7 — INFO — Plan-Tabelle §3 führt `harness/README.md` nicht

- **quelle:** `AGENTS.md` §6 Schritt 4
- **pfad:** `harness/README.md` · „kein Gate · adoptiert aus ai-harness-init · [`SPEC-102`]"
- **befund:** Der Diff ändert drei Zeilen des Gate-Index; §3 des Plans nennt nur Spezifikation und Sensor-Dateien. Die Abgrenzung schließt den Index nicht aus, und für SPEC-102/103 ohne Sensor-Datei ist er der einzige Verweisort — keine Ausweitung, aber eine Plan-Tabelle, die den Diff nicht ganz beschreibt.
- **verifizierbar:** nein
- **klasse:** plan-tabelle-unvollstaendig

## Negativbefunde

- **SPEC-097 gegen `tools/image-scan.sh` und `image-scan.yml`:** Default-Refs, Leerraum ⇒ Exit 2, Plattformen aus `imagetools inspect` mit `unknown/*`-Ausschluss und Alles-oder-nichts (`plattformen_aus_antwort`), drei Läufe je Plattform, zweites Segment als Soll-Architektur, Exit-Vorrang 2 vor 1, fail-closed ohne Docker/Netz, Nachtlauf liest den Log und wird bei Befund und Scheitern rot — geprüft, ohne Befund.
- **SPEC-098 gegen `pin-freshness.sh` und Makefile:** Gleich/Ungleich, `VERALTET` auf stderr, Exit 3; Redirect-Endstation `/releases/tag/`; symmetrisches `v`- bzw. `go`-Strippen; Versionsform bei go.dev; Pin-Quellen; erste `uses:`-Zeile über `head -1`; SKIP bei fehlendem `curl`, Netz, leerem Pin; Exit 1 nur bei Modus-Fehler; Nachtlauf ohne `freshness-trivy` — geprüft, ohne Befund.
- **SPEC-099:** Listen-Digest über `{{.Manifest.Digest}}`, Digest-Form beidseitig, `FROM`-Extraktion mit eingesetzten Build-Argumenten, SKIP bei fehlendem Docker, Nicht-Digest-Antwort, Pin ohne Digest und fehlender `FROM`-Zeile (leerer Pin), Nachtlauf ohne `trivy-digest` — geprüft, ohne Befund.
- **SPEC-100 gegen `check_latest`:** Pin aus der ersten `**Stand:**`-Zeile, sonst Exit 1; eine Seite à 100, `sort -V`; current/newer/ahead/skip; Content-Vergleich gegen `SHA256SUMS`, SKIP-Ursachen; Vorrang 4 vor 3; fehlendes `curl` ⇒ Exit 1 — geprüft, ohne Befund.
- **SPEC-101 gegen `nightly-state.sh`:** `per_page=1`, ohne Token/`gh`; `gruen` auf stdout, übrige Urteile auf stderr; `null` ⇒ SKIP; `failure` mit festem Hinweis; sonstige Ausgänge als Lauf-Störung; Form-Prüfung ⇒ SKIP; immer Exit 0 außer `--selftest`; Leerraum-`NIGHTLY_WORKFLOWS` still — geprüft, ohne Befund.
- **SPEC-102 Exit-Codes** (2 unauflösbar/keine Zahl, 1 leer, 0 sonst, `--staged` meldet und endet 0) und **SPEC-103 Glieder** (lokales `core.hooksPath` vorher leer, nachher gesetzt und `-ef` auf den Träger, Rot-Commit lässt `HEAD` stehen, Grün-Commit mit Betreff-Abgleich, `make gates` im Klon, jedes Scheitern Exit 1) — geprüft, ohne Befund jenseits F-1 und F-4.
- **Plan-Entscheidung `baseline-probe` ohne eigenen Eintrag:** `--selftest` fährt ausschließlich `check_aliases`, also Frage (3) von SPEC-092, und trifft keine eigene Festlegung; analog `guard-probe` unter SPEC-093. `history-range-guard` und `selbstpruefung` tragen eigene Exit-Semantik — eigener Eintrag richtig. Ohne Befund.
- **Referenz-Richtung (`AGENTS.md` §3.4):** SPEC-097 bis SPEC-103 tragen keine Slice-, ADR-, Wellen-, MR- oder Hash-Tokens; der Historie-Eintrag nennt keine. Ohne Befund.
- **Kommentar-Regel §3.7 / Plan-Abgrenzung:** der Diff berührt nur Markdown (`--stat`: Spezifikation, Gate-Index, sechs Sensor-Dateien); kein Skript, kein Makefile, kein Workflow. Ohne Befund.
- **Spiegel, ohne Befund:** `harness/README.md` (Trivy-Zeile korrigiert), `AGENTS.md` (nennt keine dieser Aussagen), [ADR-0066](../plan/adr/0066-cve-scan-gegen-das-publizierte-image.md) (behauptet die Trivy-Achsen nicht im Nachtlauf), [`MR-053`](../../harness/conventions.md#mr-053) („wird anders behandelt" meint dort den Menschen in der Vorprüfung, nicht das Werkzeug), Skript-Kopf `pin-freshness.sh` (fail-open bei fehlendem `curl` stimmt dort).
- **Zustandsfelder:** keine im Diff. Ohne Befund.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 2 | 1 | 4 |

Wiederkehrende Klasse: **spiegel-nicht-nachgezogen** (F-2, F-3) — Prosa außerhalb der
Spezifikation und der Sensor-Dateien trägt den alten Wortlaut weiter; dieselbe Klasse,
die der Plan als `BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen` gesichtet hat.

## Verdikt

**Nicht freigegeben, solange F-1 und F-2 offen sind.** Die sieben Einträge stimmen in jeder
geprüften Aussage mit dem Code überein; F-1 ist eine Randform, die der Code trifft und der
Eintrag nicht nennt, F-2 ein Spiegel, der die korrigierte Aussage weiter verneint. Beide
liegen in Dateien, die der Plan nicht als Gegenstand führt — ob ihre Behandlung den Plan
ändert (Workflow-Kommentar) oder in SPEC-102 bleibt, ist Sache der Implementation.
