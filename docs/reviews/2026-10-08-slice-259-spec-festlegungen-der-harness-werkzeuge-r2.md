# Review R2 — slice-259: Festlegungen der Harness-Werkzeuge in der Spezifikation

- **Review-Art:** Code. Geprüft wurde der Fix-Commit gegen den Slice-Plan (`slice-259`, Plan-Änderung
  nach R1), gegen die eigene Festlegung `SPEC-089`..`SPEC-091` (Spezifikation §7), gegen die Hard
  Rules `AGENTS.md` §3.6/§3.7/§6 Schritt 4 und §5 Regel 13/15, und gegen die R1-Findings F-1..F-4.
  Gegenstand ist die Maintainability; die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-259` · Commit `7209e90f` (Range `26d63773..7209e90f`); Dateien
  `tools/coverage-gate.sh`, `Dockerfile`, `Makefile`, `spec/spezifikation.md`,
  `harness/sensors/semgrep.md`, Slice-Plan.
- **Skill:** `reviewer.md` @ 1.18.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-08
- **Eingangs-Kontext:** R1-Report zu `slice-259` (Findings F-1..F-6); Slice-Plan mit Plan-Änderung
  nach R1; Spezifikation §7; `.golangci.yml`; `tools/semgrep.sh`; `Dockerfile` Stage `coverage`;
  `Makefile` (`THRESHOLD`, `coverage-gate`, `semgrep`); `harness/rules/kommentare-fuenf-klassen.md`.
- **Proben (eigene Läufe, echte Ausgabe):**
  - `tools/coverage-gate.sh` direkt (Host-bash, präparierte Eingaben im Scratchpad), Eingabe 50,0 %:
    Schwelle `93` → 1 · `93.0` → 1 · `0` → 0 · `50` → 0 · `93.` → 2 · `.5` → 2 · `1e2` → 2 ·
    `"93 "` → 2 · `" 93"` → 2 · `"93\n"` → 2 · `+5` → 2 · `0x10` → 2 · `""` → 2 · `abc` → 2 ·
    `-5` → 2 · `093` → 1 · `99999999999999999999999` → 1 · `100.000001` → 1.
  - Eingabe-Formen bei Schwelle 93: `93.0%` → 0 · `100.0%` → 0 · `n/a` → 2 (mit Meldung) ·
    `93%` → 2 · `93,0%` → 2 · `93.0% ` (Leerzeichen am Ende) → 2 · `93.0%\r` → 2 · zwei
    `total:`-Zeilen (50,0 / 99,0) → 1 mit `printf: 50.0\n99.0: invalid number` auf stderr.
  - Locale: neues Skript unter `LC_ALL=de_DE.UTF-8`, 93,0 % / 93 → 0. Altes Skript (`56de3d24`)
    unter derselben Locale → 1 mit `printf: 93.0: Ungültige Zahl.` — die Begründung für
    `LC_ALL=C` ist bestätigt.
  - `make coverage-gate THRESHOLD=` → Build-Log `coverage-gate: Schwelle '' ist keine nicht
    negative Zahl`, docker `exit code: 2`, **make-Exit 2**. `make coverage-gate THRESHOLD=99` →
    `FAIL — Coverage 94.70% unter Schwelle 99%`, docker `exit code: 1`, **make-Exit 2**.
    `make coverage-gate` → `OK — Coverage 94.70% erfüllt Schwelle 93%`, Exit 0.
  - `tools/semgrep.sh` mit drei temporären, nicht getrackten Probe-Dateien (`crypto/md5`,
    Regel `use-of-md5`) unter `internal/zzrvwprobe/`, `internal/zzrvwprobe/test/`,
    `internal/zzrvwprobe/build/`: `Ran 55 rules on 66 files: 1 finding.` — gemeldet nur
    `internal/zzrvwprobe/a.go`; `Files matching .semgrepignore patterns: 81` (79 Testdateien + die
    beiden Proben unter `test/` und `build/`). Proben danach entfernt, `git status` sauber.
  - `git ls-files '*.go' | grep -E '(^|/)(test|tests|build|dist|vendor|node_modules)/'` → leer.

## Findings

### F-1 — HIGH: Neuer Skript-Kommentar erzählt das frühere Verhalten; die Herkunfts-Zeile bleibt im umgeschriebenen Kopf

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §3.7 (Kommentar-Klassen); Skill-Prüffrage 6; `harness/rules/kommentare-fuenf-klassen.md` §Bestandsgrenze („geräumt wird beim nächsten Anfassen der Zeile"); `AGENTS.md` §5 Regel 15
- **pfad:** `tools/coverage-gate.sh` · „verglich awk sie bisher als 0 — jede Coverage bestand still"
  (neu, Zeile 27); `tools/coverage-gate.sh` · „Muster: u-boot
  # scripts/coverage-gate.sh (gleiche Build-Familie)" (im Fix umgebrochen, Zeilen 4–5)
- **befund:** Der neue Kommentar über der Schwellen-Prüfung beschreibt, was das Skript **früher**
  tat („bisher … bestand still") — Chronik, keine der fünf Klassen; der Code darunter beantwortet
  die Frage nach dem Warum mit dem folgenden Satz bereits allein. Die Herkunfts-Angabe „Muster:
  u-boot …" steht in einer Zeile, die der Fix neu schreibt, und ist kein auflösbares Feld
  (`DC-*`/`ADR-*`/`MR-*`/`seit welle-`); die Commit-Botschaft sagt dagegen „dabei Herkunfts-Prosa
  entfernt". Failure-Szenario: Der nächste Leser des Gate-Skripts bekommt Review-Historie als
  Zusage zu lesen, und die Botschaft behauptet eine Bereinigung, die für diese Zeile nicht
  stattfand — genau die Klasse, die F-4 aus R1 schließen sollte.
- **verifizierbar:** nein — kein Gate liest Kommentar-Klassen.
- **klasse:** `kommentar-traegt-chronik`

### F-2 — MEDIUM: `SPEC-089` unterscheidet Exit 1 und Exit 2 für `make coverage-gate`; am Werkzeug sind beide Exit 2

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §5 Regel 13 und Regel 14; Skill-Prüffrage 10 (Messmethode gegen Spec-Stelle) und 17 (Proxy-Messung)
- **pfad:** `spec/spezifikation.md` §7 · „genau 93,0 besteht; darunter Exit 1" und „endet mit
  Exit 2, ebenso eine fehlende Coverage-Eingabe" (Zeile `SPEC-089`, im Fix neu geschrieben);
  Commit-Botschaft `7209e90f` · „make coverage-gate THRESHOLD=abc -> Exit 2 mit Meldung"
- **befund:** Die Zeile steht unter der Spalte *Werkzeug* = `make coverage-gate`. Dort endet jedes
  Scheitern mit Exit 2: gemessen `THRESHOLD=99` (Skript 1, docker `exit code: 1`) → make 2, und
  `THRESHOLD=` (Skript 2) → make 2 — `make` meldet ein gescheitertes Rezept immer mit 2. Die
  Unterscheidung „verfehlt" (1) gegen „gescheitert, nicht bestanden" (2) existiert nur am Skript
  und im Build-Log. Die Messung der Botschaft (make-Exit 2 bei `abc`) belegt deshalb nicht, was sie
  belegen soll: dasselbe Exit 2 käme auch bei einer gültigen Schwelle über der Coverage. Dieselbe
  Lesart trägt `SPEC-091` („endet er mit Exit 2 statt grün"; `tools/semgrep.sh` endet bei Befund
  mit 1, über `make` beides 2). Failure-Szenario: Wer nach `SPEC-089` am Exit-Code von
  `make coverage-gate` entscheidet, ob Carveout-Pflicht (Verfehlung) oder ein kaputter Lauf
  vorliegt, liest in beiden Fällen 2.
- **verifizierbar:** ja — `make coverage-gate THRESHOLD=99; echo $?` gegen `THRESHOLD=abc`.
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand-geprueft`

### F-3 — MEDIUM: `SPEC-091` legt die Prüfmenge als „ohne Testdateien" fest; semgrep lässt zusätzlich ganze Verzeichnisse aus

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §5 Regel 13; `AGENTS.md` §3.8; Skill-Prüffrage 18
- **pfad:** `spec/spezifikation.md` §7 · „Gescannt werden die Go-Dateien des Arbeitsbaums, die
  `.gitignore` nicht ausnimmt — getrackt oder nicht —, **ohne** Testdateien (`*_test.go`)";
  `harness/sensors/semgrep.md` §Grenze · „**Testdateien werden nicht gescannt**"
- **befund:** Die Voreinstellung von semgrep (`.semgrepignore`-Default) nimmt nicht nur
  `*_test.go`, sondern auch Verzeichnisse wie `test/` und `build/` aus — gemessen: zwei nicht
  getrackte Probe-Dateien mit Befund unter `internal/zzrvwprobe/test/` und `…/build/` werden
  übersprungen (Skip-Zahl 79 → 81), nur die Probe direkt in `internal/zzrvwprobe/` wird gemeldet.
  Heute liegt keine Go-Datei in einem solchen Verzeichnis, die Festlegung liest sich aber als
  vollständige Prüfmenge und die Sensor-Grenze nennt nur die Testdateien. Failure-Szenario: Ein
  künftiges Paket `internal/…/build/` mit einem `go/lang/security`-Muster bleibt grün, während
  `SPEC-091` es als gescannt festlegt.
- **verifizierbar:** ja — Probe-Datei unter einem `build/`-Verzeichnis, `make semgrep`.
- **klasse:** `grenzen-liste-wird-als-vollstaendig-gelesen`

### F-4 — LOW: Die Mitnahme `LC_ALL=C` steht im Plan erst im selben Commit wie der Code

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §6 Schritt 4 („gehört vor den Code, nicht in den Bericht danach")
- **pfad:** Slice-Plan `slice-259` · „Mitgenommen: `LC_ALL=C` im Skript"
- **befund:** Die Plan-Änderung nach R1 landete als eigener Commit (`26d63773`) vor dem Code; der
  Satz zur Mitnahme `LC_ALL=C` kam erst im Fix-Commit `7209e90f` dazu, zusammen mit der
  Skript-Änderung. Die Mitnahme ist begründet und gemessen; nur ihre Reihenfolge weicht ab.
- **verifizierbar:** nein — `git show 26d63773` vs. `git show 7209e90f -- docs/plan/`.
- **klasse:** `mitnahme-nach-dem-code-geplant`

### F-5 — LOW: Die Hilfe-Zeile des Targets verdrahtet die Schwelle erneut als Zahl

- **kategorie:** LOW
- **quelle:** dieser Skill, LOW-Anker „latente Wartungsfalle (hart verdrahteter Wert)"
- **pfad:** `Makefile` · „coverage-gate: ## Coverage-Schwelle 93 % (SPEC-089;"
- **befund:** Nach dem Umzug der Autorität auf `SPEC-089` steht 93 jetzt an drei Orten im Code
  (`THRESHOLD ?= 93`, `ARG COVERAGE_THRESHOLD=93`, Hilfe-Text) und einmal in der Spezifikation; die
  Hilfe-Zeile ist die einzige der drei ohne Wirkung und zündet erst beim nächsten Hub der Schwelle.
- **verifizierbar:** nein.
- **klasse:** `wert-an-mehreren-orten-verdrahtet`

### F-6 — INFO: Randformen, die `SPEC-089` nicht entscheidet

- **kategorie:** INFO
- **quelle:** Maintainability (kein Konventions-Anker — deshalb INFO)
- **pfad:** `tools/coverage-gate.sh` · „grep -E '^total:'"; Regex
  „`^[0-9]+(\.[0-9]+)?$`"
- **befund:** (a) Eine Schwelle über 100 (auch `99999999999999999999999`) ist gültig und endet
  immer mit 1; (b) `93.` und `+5`, die awk früher als Zahl las, sind jetzt Exit 2 — kein Aufrufer
  im Repo nutzt sie (`grep` über `THRESHOLD`/`coverage-gate.sh`); (c) zwei `total:`-Zeilen
  enden mit 1 und einer `printf`-Fehlermeldung statt mit 2 — Bestand, `go tool cover -func`
  schreibt genau eine. Keine dieser Formen ist ein bisher grüner, legitimer Aufruf.
- **verifizierbar:** ja — die Proben oben.
- **klasse:** `randform-nicht-festgelegt`

## Negativbefunde (geprüft, ohne Befund)

- **Skript-Ausgänge gegen `SPEC-089` (Skript-Ebene):** leere, Text-, negative, wissenschaftliche,
  Hex-, Vorzeichen- und Leerzeichen-Schwelle → 2 mit Meldung; fehlende/leere Eingabe → 2; fehlende
  `total:`-Zeile → 2; nicht lesbarer Wert (`n/a`, `93%`, `93,0%`, CRLF, Leerzeichen am Ende) → 2
  mit Meldung; darunter → 1; genau 93,0 → 0. Kein Stilles-Grün-Pfad mehr (Prüffrage 1).
- **Prüffrage 19 (`|| true` am Parse):** wirkt auf die ganze Pipeline, maskiert nur den
  Nicht-Treffer von `grep`; jeder maskierte Fall landet im Exit-2-Zweig (gemessen mit sechs
  Formen). Kein bisher grüner legitimer Aufruf fällt rot (`make coverage-gate` Default → 0).
- **`LC_ALL=C`:** `export` vor jedem Zahl-Werkzeug; awk, printf und grep laufen unter C; Begründung
  am alten Skript nachgemessen.
- **`Dockerfile`/`Makefile`-Durchreichung:** `THRESHOLD=` erreicht das Skript als leerer String
  und endet rot (gemessen); Default unverändert 93.
- **`SPEC-090` gegen `.golangci.yml`:** `errcheck.exclude-functions` = genau die drei `fmt.Fprint*`;
  `ireturn.allow` = error/empty/anon/stdlib/generic + `port/driven` + `go-billy`;
  `exclusions.generated: lax`; keine `presets`, keine `exclusions.paths`. Die drei Ausnahmen sind
  vollständig und richtig beschrieben; R1 F-2 geschlossen.
- **`SPEC-091` „getrackt oder nicht":** bestätigt — eine nicht getrackte Datei wird trotz der Zeile
  „Scan was limited to files tracked by git" gemeldet.
- **Kommentare `Dockerfile`/`Makefile` (§3.7):** Rang-Zeiger auf `SPEC-089` plus Zusage
  („Verfehlung ⇒ Carveout-Pflicht, Senkung nur per ADR"); keine Chronik mehr. Die `u-boot-Muster`-
  Erwähnung im `Dockerfile` steht in einer nicht angefassten Zeile (Bestand, grandfathered).
  `LC_ALL`- und `|| true`-Kommentar im Skript: Kopplung bzw. Grenze, Gegenwart.
- **R1-Schließung:** F-1 (a) toter Exit-2-Zweig, (b) still-grüne Schwelle, (c) roter Testlauf —
  sachlich geschlossen. F-2 geschlossen. F-3: Testdateien und nicht getrackte Dateien jetzt
  benannt — weitere Lücke siehe F-3 dieser Runde. F-4: Rang-Zeiger auf `SPEC-089` nachgezogen —
  neue Chronik im Skript siehe F-1 dieser Runde.
- **Referenz-Richtung §3.4:** neue §8-Zeile und §7-Zeilen ohne ADR-/Slice-/Harness-Verweis.
- **§3.6:** Default 93 unverändert, keine Schwellen-Senkung; die Validierung verschärft nur.
- **`harness/sensors/semgrep.md`:** Nummerierung 1–4 konsistent, kein Grenzen-Punkt verloren.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 1 | F-1 |
| MEDIUM | 2 | F-2, F-3 |
| LOW | 2 | F-4, F-5 |
| INFO | 1 | F-6 |

Wiederkehrende Klasse: `grenzen-liste-wird-als-vollstaendig-gelesen` — in diesem Slice zum
**dritten** Mal (R1 F-2, R1 F-3, R2 F-3): die Festlegung wird gegen die Konfiguration geprüft, die
der Autor kennt, nicht gegen die Voreinstellungen des Werkzeugs. Steering-Loop-Signal.

## Verdikt

**Nicht freigegeben.** F-1 (HIGH) ist ein Kommentar im Gate-Skript, der Chronik trägt, und eine
Botschaft, die für die Herkunfts-Zeile eine Bereinigung behauptet, die nicht stattfand. F-2 und F-3
sind Festlegungen in `SPEC-089`/`SPEC-091`, die der Gegenstand anders tut. Der Code-Fix selbst —
Schwellen-Validierung, `|| true`, `LC_ALL=C` — ist korrekt und gegen seine Fehlerformen
nachgemessen; die Befunde betreffen Kommentar und Spezifikationstext. F-4/F-5 vor der Closure
annehmen oder begründen; F-6 ist eine Notiz.
