# Review R1 — slice-259: Festlegungen der Harness-Werkzeuge in der Spezifikation

- **Review-Art:** Code (Doku-Diff). Geprüft wurde gegen den Slice-Plan (`slice-259`, §1
  Abgrenzung, §2, §3 Spiegel-Liste, §6 Risiken), gegen die Baseline `v6.17.0` ·
  `regelwerk/grundlagen-referenz-richtung.md` §Spec-Straten (Absatz „Auch die Werkzeuge des
  Harness treffen technische Festlegungen"), `v6.17.0` · `templates/spec/spezifikation.template.md`
  §7/§8 und `v6.17.0` · `templates/harness/sensors/gate.template.md` (Regeln der Datei,
  §Vertrag), gegen `MR-0098`, `MR-074` (Bewegung 2), den neuen Nachtrag `MR-075` und die
  Hard Rules `AGENTS.md` §3.4/§3.7/§3.8 sowie §5 Regel 13/16/17. Gegenstand ist die
  Maintainability; die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-259` · Range `d39c4b04..56de3d24` (zwei Commits: `efe8e683`
  Spezifikation §7/§8 + Kopplungen, `56de3d24` Sensor-Dateien).
- **Skill:** `reviewer.md` @ 1.18.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-08
- **Eingangs-Kontext:** Slice-Plan `slice-259`; Spezifikation §7 (`SPEC-089`..`SPEC-092`) und
  §8; die Träger der vier Festlegungen im Code: `Dockerfile` (Stage `coverage`),
  `tools/coverage-gate.sh`, `Makefile` (`THRESHOLD`, Targets `coverage-gate`, `semgrep`,
  `baseline-verify`), `.golangci.yml`, `tools/semgrep.sh`,
  `tools/harness/fetch-baseline-cache.sh` (`verify()`, `check_aliases()`); `.d-check.yml`
  (`matrix.exclude-sections`, `structure`). Vorherige Findings am selben Gegenstand: R1/Verify
  zu `slice-248` (Historie-Ausnahme, Klasse `grenzen-liste-wird-als-vollstaendig-gelesen`).
- **Proben (eigene Läufe, echte Ausgabe):**
  - `make doc-check` auf `56de3d24`: `d-check: 1007 Datei(en) geprüft, 0 Befund(e)`, Exit 0.
  - `make semgrep` auf `56de3d24`: `Scanning 1177 files tracked by git with 57 Code rules` ·
    `Scanning 65 files with 55 go rules` · `Files matching .semgrepignore patterns: 79` ·
    `Scan was limited to files tracked by git` · `Ran 55 rules on 65 files: 0 findings.`
    Gegenzählung: `git ls-files '*.go'` = 144, davon `*_test.go` = 79, übrige 65.
  - `tools/coverage-gate.sh` mit `LC_ALL=C` gegen synthetische Eingaben (Scratchpad, keine
    Repo-Änderung): Coverage 50,0 % mit Schwelle `93` → Exit 1; mit Schwelle `""` → Exit 0
    (`OK — Coverage 50.00% erfüllt Schwelle %`); mit `abc` → Exit 0; mit `-5` → Exit 0;
    `total:`-Zeile mit Wert `n/a` → Exit 1 **ohne Meldung**; leere Eingabe → Exit 2; Eingabe
    ohne `total:` → Exit 2.

## Findings

### F-1 — MEDIUM: `SPEC-089` legt eine Randform fest, die der Code anders entscheidet, und lässt den Schwellen-Parameter aus

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §5 Regel 13; Skill-Prüffrage 10 (Messmethode gegen Spec-Stelle) und 18
- **pfad:** `spec/spezifikation.md` §7 · „ist der Wert nicht lesbar, endet das Gate mit Exit 2:
  gescheitert, nicht bestanden" (Zeile 3634 als Lesehilfe); Gegenstand `tools/coverage-gate.sh` ·
  „total_pct=\"$(echo \"$total_line\" | grep -oE" und `Makefile` · „THRESHOLD ?= 93"
- **befund:** (a) Der Zweig „Prozent nicht parsbar → Exit 2" ist unerreichbar: Unter
  `set -euo pipefail` bricht die Zuweisung ab, wenn `grep -oE` nichts findet, und das Skript endet
  mit Exit 1 ohne Meldung — gemessen mit `n/a` als Wert; die Spezifikation legt jetzt Exit 2
  fest. (b) Die Schwelle ist ein `make`-Parameter (`THRESHOLD ?= 93`, durchgereicht als
  `COVERAGE_THRESHOLD`); `SPEC-089` nennt nur „Schwelle 93 %" und entscheidet die Randform einer
  leeren oder nicht-numerischen Schwelle nicht — der Code entscheidet sie als **bestanden**
  (gemessen: `""`, `abc`, `-5` → Exit 0 bei 50 % Coverage). (c) Ein roter Testlauf macht das Gate
  rot, bevor gemessen wird; auch diese Randform fehlt. Failure-Szenario: Wer `SPEC-089` als
  Autorität liest, hält einen leeren `THRESHOLD` (z. B. `THRESHOLD= make gates` — eine leer
  gesetzte Umgebungsvariable gilt für `?=` als definiert) für ein gescheitertes Gate; der Lauf
  ist grün bei beliebiger Coverage.
- **verifizierbar:** ja — die drei Proben oben; kein Gate fängt die Abweichung Spec ↔ Code.
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand-geprueft`
- **Hinweis zur Einordnung:** Der still-grüne Pfad selbst (b) liegt im **Bestand**
  (`tools/coverage-gate.sh`, nicht im Diff) und wäre in einem Diff, der das Skript berührt, ein
  HIGH nach Prüffrage 1. Der Plan schließt Code-Änderungen nicht ausdrücklich aus, führt aber nur
  Doku-Dateien (§3); die Behebung ist deshalb eine Plan-Frage, nicht still mitzunehmen
  (`AGENTS.md` §6 Schritt 4).

### F-2 — MEDIUM: `SPEC-090` zählt die Ausnahmen als vollständige Menge auf und lässt drei aus der Konfiguration weg

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §5 Regel 13; Skill-Prüffrage 18 (Grenzen-Liste ohne größte Lücke)
- **pfad:** `spec/spezifikation.md` §7 · „Ausnahmen gelten nur zentral in der Konfiguration,
  **sechs** Regeln"; `harness/sensors/lint.md` · „welche es sind, zählt"; Gegenstand
  `.golangci.yml` · „exclude-functions:", „generated: lax", „ireturn:"
- **befund:** Die sechs `exclusions.rules` stimmen (nachgezählt), aber die Konfiguration
  verkleinert den Prüfbereich an drei weiteren Stellen, die `SPEC-090` nicht nennt: `errcheck`
  ignoriert `fmt.Fprintln`/`fmt.Fprintf`/`fmt.Fprint` repo-weit (nicht nur in Tests),
  `exclusions.generated: lax` nimmt generierte Dateien heraus, und `ireturn` erlaubt eine
  Typ-Liste einschließlich Fremd-Paket `go-billy`. Die Formulierung „nur … sechs Regeln" liest
  sich als Menge, und die Sensor-Datei leitet ihre Grenze 3 jetzt ausdrücklich daraus ab
  („sauber außerhalb dieser Ausnahmen"). Failure-Szenario: Ein unbehandelter Schreibfehler über
  `fmt.Fprintf` in Produktions-Code (CLI-Reporter) bleibt grün, während Spezifikation und
  Sensor-Datei zusammen zusagen, `errcheck` gelte außerhalb der Testdateien ohne Ausnahme.
- **verifizierbar:** ja — Lesen von `.golangci.yml` gegen die §7-Zeile; kein Gate.
- **klasse:** `grenzen-liste-wird-als-vollstaendig-gelesen`

### F-3 — MEDIUM: `SPEC-091` legt „über das Repo" fest; gescannt werden 65 von 144 Go-Dateien

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §5 Regel 13; Skill-Prüffrage 18; `AGENTS.md` §3.8 (Zusage über die
  Scan-Menge)
- **pfad:** `spec/spezifikation.md` §7 · „netzlos über das Repo. Jeder Befund ist rot.";
  `harness/sensors/semgrep.md` §Grenze (Punkte 1–3)
- **befund:** Der Lauf selbst meldet `Files matching .semgrepignore patterns: 79` und
  `Scan was limited to files tracked by git`; das Repo führt keine eigene `.semgrepignore`, die
  Ausschlüsse kommen aus semgreps Voreinstellung und treffen genau die 79 `*_test.go`-Dateien.
  Weder `SPEC-091` (Prüfmenge) noch die Grenzen-Liste der Sensor-Datei nennt das — die größte
  Lücke (55 % der Go-Dateien) fehlt, während die Liste drei kleinere führt. Ob nicht getrackte,
  nicht ignorierte Dateien mitgescannt werden, ist nicht gemessen; die Lauf-Zeile sagt „tracked".
  Failure-Szenario: Eine Test-Hilfsfunktion mit einem `go/lang/security`-Muster (etwa
  `exec.Command` mit zusammengesetztem String) bleibt grün, obwohl die Festlegung „über das
  Repo" sie einschließt; ebenso eine neue, noch nicht `git add`-ete Datei im lokalen
  `make gates`-Lauf, falls die „tracked"-Zeile wörtlich gilt.
- **verifizierbar:** ja — `make semgrep` (Summary-Zeilen) gegen `git ls-files '*.go'`.
- **klasse:** `grenzen-liste-wird-als-vollstaendig-gelesen`

### F-4 — LOW: Rang-Zeiger im Code zeigen für die Schwelle weiter auf `harness/README.md`

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §3.7 (Kommentar-Klasse Rang-Zeiger); `MR-025` (Spiegel vor dem
  Editieren); Baseline `v6.17.0` · `regelwerk/grundlagen-referenz-richtung.md` §Spec-Straten
  („ebenso wenig der Kopf des Skripts — der trägt, womit das Werkzeug gedeckt ist, nicht, was es
  zusagt")
- **pfad:** `tools/coverage-gate.sh` · „aktuelle
  # Schwelle und Historie in harness/README.md §Sensors"; `Dockerfile` · „Kalibrierungs-Bindung
  (harness/README.md §Sensors): Schwelle 93 %"; `Makefile` · „Historie in harness/README §Sensors"
- **befund:** Seit `SPEC-089` ist die Spezifikation die Autorität der Schwelle; die drei
  Kommentare benennen weiter die Index-Zeile als Ort der „aktuellen Schwelle". Die Plan-Spiegelliste
  (§3) nennt die Träger, aber nicht ihre Zeiger-Kommentare. Die Kette löst noch auf (Index-Zelle
  → `SPEC-089`), aber über einen Hop, und der Kommentar sagt das Gegenteil von dem, was §7 regelt.
- **verifizierbar:** nein — kein Gate liest Kommentar-Zeiger.
- **klasse:** `rang-zeiger-nach-autoritaets-umzug-nicht-nachgezogen`

### F-5 — INFO: Zwei Aussagen stehen jetzt doppelt

- **kategorie:** INFO
- **quelle:** Baseline `v6.17.0` · `templates/harness/sensors/gate.template.md` §Vertrag („an zwei
  Orten liefen sie auseinander")
- **pfad:** `harness/sensors/semgrep.md` · „Lesen: das Holen des Regelsets"; `spec/spezifikation.md`
  §7 · „Das Holen des Regelsets am Pin ist Setup und braucht Netz"; ebenso
  `harness/sensors/baseline-verify.md` Grenze 0 · „Der Lauf beweist innere Konsistenz, nicht
  Echtheit" und `SPEC-092` · „Geprüft wird innere Konsistenz, nicht Echtheit"
- **befund:** Beide Aussagen stehen in Spezifikation und Sensor-Datei. Inhaltlich deckungsgleich,
  aber zwei Orte für eine Aussage; die Vorlage trennt Festlegung (Spec) und Lese-Hinweis/Grenze
  (Sensor-Datei). Designnotiz, kein Versagen.
- **verifizierbar:** nein.
- **klasse:** `aussage-an-zwei-orten`

### F-6 — INFO: Das Produkt schlägt Konsumenten weiter `"7. Historie"` vor

- **kategorie:** INFO
- **quelle:** Maintainability (kein Konventions-Anker — deshalb INFO, nicht LOW)
- **pfad:** `internal/hexagon/core/app/suggest.go` · „exclude-sections: [Historie, \"7. Historie\",
  Geschichte]"; `spec/spezifikation.md` §2-Beispiel · „exclude-sections: [Historie, \"7. Historie\",
  Geschichte]"
- **befund:** Die Baseline-Vorlage `v6.17.0` setzt die Historie der Spezifikation auf §8; ein
  Konsument, der ihr folgt und den `--suggest`-Vorschlag übernimmt, verliert die Ausnahme für seine
  Spezifikations-Historie. Außerhalb des Plans (Produkt-Code, Doku-Slice) — notiert, nicht
  gefordert.
- **verifizierbar:** nein.
- **klasse:** `produkt-vorschlag-hinkt-baseline-vorlage`

## Negativbefunde (geprüft, ohne Befund)

- **`SPEC-089` Messbasis:** `-coverpkg` aus `go list ./internal/...`, `-covermode=atomic`, `./...`
  als Testsuite, Vergleich `p+0 >= t+0` (genau 93,0 besteht), leere/fehlende Eingabe und fehlende
  `total:`-Zeile → Exit 2 — stimmt mit `Dockerfile` und `tools/coverage-gate.sh` (abgesehen von F-1).
- **`SPEC-090` Profil:** 29 eingeschaltete Linter nachgezählt (5 + 24), `default: none`; alle
  acht genannten Schwellen stimmen mit `settings` (cyclop/gocyclo 15, gocognit 20, funlen 100/60,
  nestif 5, dupl 150, maintidx 20, interfacebloat 10); sechs `exclusions.rules` nachgezählt und
  richtig beschrieben; `nolintlint` mit allen drei Schaltern scharf.
- **`SPEC-091` Pins und Abbruch:** Image digest-gepinnt, Commit-Pin, Umfang `go/lang/security`,
  `--error`, `--network none`, Abbruch mit Exit 2 ohne `Ran N rules` (N ≥ 1) — stimmt mit
  `tools/semgrep.sh` (abgesehen von der Prüfmenge, F-3).
- **`SPEC-092`:** drei Fragen gegen `verify()`/`check_aliases()` geprüft — `sha256sum -c`,
  Zählung `find -type f` ohne Manifest gegen `grep -c .` mit `> 0`, rekursives `find -type l`
  unter `.claude/rules/` mit Existenzprüfung; fehlende Host-Werkzeuge und fehlendes Manifest →
  Exit 1. Kein Befund.
- **Referenz-Richtung (§3.4):** §7 und die neue §8-Zeile verlinken weder ADR noch Slice noch eine
  Harness-Datei; die Werkzeuge stehen als `make`-Token, wie die Vorlage es zeigt. `make doc-check`
  grün (Matrix eingeschlossen).
- **Kopplungen der Umnummerierung:** `matrix.exclude-sections` trägt `"8. Historie"`, die
  `structure`-Regel der Spezifikation zeigt auf `## 8. Historie`, die Lastenheft-Regel bleibt auf
  `## 7. Historie`; der Kommentar über `exclude-sections` ist nachgezogen; kein Anker-Link auf
  `#7-historie` der Spezifikation im lebenden Bestand (grep). `MR-0098` bleibt unverändert und
  wird durch `MR-075` ergänzt — die Nachtrag-Form, die das Repo für akzeptierte Einträge führt.
- **`MR-075` Form und Index:** Pflichtfelder vorhanden, kein Status-Feld (Vorlage), Nummer dicht
  nach `MR-074`, Index-Zeile mit Kurz-Anker; Geltungsbereich und Auflösungs-Trigger tragen die
  Aussage, für die der Eintrag steht.
- **Kennungs-Vergabe und Adressierungs-Form (Prüffrage 16):** `SPEC-089`..`SPEC-092` dicht nach
  `SPEC-088`; die neuen Verweise (Index-Zelle `coverage-gate`, drei Sensor-Dateien) nennen die
  Kennung im Text und verlinken den Abschnitt.
- **Sensor-Dateien — verlorene Grenzen:** Gestrichen wurden nur Festlegungen, die jetzt in §7
  stehen (Linter-Zahl, `exclude-rules`-Inhalt, `--error`, Netz-Pfad, die drei Fragen samt
  „rekursiv und dotfile-bewusst"); alle Grenzen-Punkte (lint 1–3, semgrep 1–3, baseline-verify
  0–4) sind erhalten. Kein Lese-Hinweis verloren.
- **Kommentare (§3.7):** Der neue `.d-check.yml`-Kommentar trägt Abgrenzung und ein auflösbares
  Feld (`MR-075`, `ADR-0097`), keine Chronik.
- **Abgrenzung des Plans:** Der Diff berührt genau die Dateien aus Plan §3; keine Werkzeuge
  außerhalb von `make gates`, keine ADR, kein Lastenheft, kein Produkt-Code, kein `CHANGELOG.md`
  (§5 Regel 17).
- **Zustandsfelder (Prüffrage 7):** keine Roadmap-/Register-Zelle im Diff; die Bindung-Zelle von
  `coverage-gate` trägt keinen eingefrorenen Ist-Wert mehr.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 3 | F-1, F-2, F-3 |
| LOW | 1 | F-4 |
| INFO | 2 | F-5, F-6 |

Wiederkehrende Klasse: `grenzen-liste-wird-als-vollstaendig-gelesen` (F-2, F-3) — die Festlegung
wurde am Code geprüft, aber an den Stellen, die der Autor als Prüfgegenstand kannte
(`exclusions.rules`, die Skript-Zweige), nicht an der Konfiguration daneben (`exclude-functions`,
semgreps Voreinstellungen).

## Verdikt

**Nicht freigegeben.** Drei MEDIUM-Findings blockieren: Die Spezifikation ist mit diesem Slice die
Autorität über Schwelle und Randform, und an drei Stellen legt sie etwas fest, das der Gegenstand
anders tut. F-1(b) berührt zudem einen still-grünen Pfad im Bestand des Gate-Skripts — ob er im
Slice behoben (Plan-Änderung vor dem Code) oder als Folge-Slice geführt wird, ist eine
Plan-Entscheidung; die Festlegung in `SPEC-089` muss in beiden Fällen sagen, was der Code heute
tut. F-4 ist vor der Closure nachzuziehen oder zu begründen; F-5/F-6 sind Notizen.
