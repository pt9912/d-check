# Verifikation — slice-280: `hostpaths` erkennt `$HOME` und `${HOME}` als Home-Verweis

- **Rolle:** Verifier (Modul 11) — DoD, Plan und Spiegel gegen den Stand
- **Gegenstand:** `docs/plan/planning/done/welle-92/slice-280-hostpaths-home-variablen.md` §1, §2, §3 (Negativliste, Plan-Änderungen, Spiegel), §5, §6; `ADR-0111`
- **Range:** `ab7de630..HEAD` (Spitze `57ba771b`); der Review-Report `2026-10-29-slice-280-hostpaths-home-variablen.md` nur als Kontext
- **Datum:** 2026-10-29

## Messungen (selbst gefahren)

Alle Mutationen am Arbeitsbaum, jeweils `make test`, danach die Datei aus einer vorab gezogenen Kopie zurückgespielt; am Ende `git diff` 0 Zeilen, `git status --short` leer, SHA-256 beider Dateien gleich der Kopie.

| Nr | Prüfung | Behauptet | Gemessen | |
|---|---|---|---|---|
| 1 | `make gates` auf sauberem Arbeitsbaum | grün | Exit 0, Schlusszeile `[gates] … green`, coverage-gate „94.70% erfüllt Schwelle 93%", Semgrep/doc-check/planning-check 0 Befunde | ok |
| 2 | Mutation A: `variableHostpathRE` liefert nie (`len(names) >= 0` → `nil`) | `$HOME`-Test rot aus dem richtigen Grund | rot, acht Tests, durchweg „kein Fund statt `hostpath-forbidden`": `…Happy` (`Targets = []`, want `$HOME/Dev/x/y.log` …), `…Liste`, `…ExemptTargets`, `…ZielWechselt`; die übrigen (`…OhneUnixVerlust`, `…VerdraengtUnix`, `…PraefixAlsErstesSegment`, `…UeberlappungGenauEinmal`, `…EmphasisLikeTilde`) verlieren den Variablen-Fund. Boundary/Negatives bleiben grün (erwartet, sie prüfen Stille) | ok |
| 3 | Mutation B: `dropInside` gibt `ms` unverändert zurück | Verdrängungs-Test rot wegen zweitem Fund | rot: `…VerdraengtUnix` (der Unix-Teil des Fund-Pfads zusätzlich), `…PraefixAlsErstesSegment` (Unix-Treffer mit Präfix-Name als erstem Segment zusätzlich), `…UeberlappungGenauEinmal` (Doppelfunde für Tilde- und Variablen-Anteil), `…ZielWechselt` („das Variablen-Glob nimmt den Fund aus“ mit einem Unix-Fund in der Liste), dazu der Bestandstest `TestHostpathsTilde` (Unix-Treffer innerhalb der Tilde nicht mehr verdrängt) | ok |
| 4 | Mutation C: Punkt-Segment-Regel fehlt (`.` aus der Segment-Klasse der Variablen-Regex) | Boundary rot | genau `TestHostpathsVariablesBoundary`: „Punkt-Segment und nackte Variable bleiben still, got [$HOME/.config/x]" | ok |
| 5 | Mutation D: `startsWithListedVariable` immer wahr | Config-Test rot | `TestDecode_HostpathsVariables`: „$-Glob ohne passende Variable muss Exit 2 sein" (ohne `variables`) | ok |
| 6 | Mutation E: `shellIdentRE` akzeptiert alles (`.*`) | Config-Test rot | `TestDecode_HostpathsVariables`: `hostpaths.variables "": err = <nil> (Ablehnung erwartet)` | ok |
| 7 | Zusatz-Mutation F: `withoutClosingQuote` gibt den Pfad unverändert zurück | (nicht benannt) | rot: `…Happy` (`$HOME"/Dev/w`), `…ExemptTargets`, `…PraefixAlsErstesSegment` | ok |
| 8 | Zusatz-Mutation G: Wortgrenze ohne `$` und Backslash | (nicht benannt) | rot: `…Negatives` (`$$HOME/x/y`, `\$HOME/x/y` melden) | ok |
| 9 | Vorher-Image aus `git archive v0.87.1`, `--target runtime`; Nachher-Image = `d-check:latest` aus HEAD (von `make gates` gebaut) | — | beide gebaut; das Vorher-Image kennt `variables` nicht (Config mit Schlüssel ⇒ Exit 2 „field variables not found") | — |
| 10 | (d)(i) Repo-Kopie (1611 getrackte Dateien ohne `.harness/baseline`), `modules: [hostpaths]`, ohne Schlüssel, `--enable hostpaths`, `--network none`, `:ro` | byte-identisch | Standard-Präfixliste: stdout (0 Byte), stderr (46 Byte), Exit (0) identisch. Weil das leer ist, zusätzlich mit erweiterter Präfixliste (home, mnt, Users, tmp, Development, Dev, docs, internal, repo, var, etc): 16 Zeilen stdout (1681 B) und JSON (3534 B), stderr, Exit 1 — jeweils `cmp` identisch | ok |
| 11 | (d)(ii) Fixture Issue 12 (`$HOME/Dev/x/y.log`, `${HOME}/Dev/z`, `"$HOME"/Dev/w`, Inline-Code `$HOME/Dev/code`, ein Tilde-Pfad; Präfix `Dev` gesetzt), ohne Schlüssel | byte-identisch | Text, `--json`, `--yaml`: stdout, stderr, Exit je `cmp` identisch (3 Funde: zwei Unix-Reste der Formen mit `}` und Anführungszeichen, ein Tilde-Pfad, Exit 1) | ok |
| 12 | (e) dasselbe Fixture, `hostpaths: {variables: [HOME]}`, Nachher-Image | alle drei Formen + Inline-Code + Tilde, sonst nichts | genau fünf Funde: `$HOME/Dev/x/y.log` (Z.3), `$HOME/Dev/w` (Z.4, ohne schließendes Anführungszeichen), `${HOME}/Dev/z` (Z.4), `$HOME/Dev/code` (Z.6), der Tilde-Pfad (Z.6); keine Unix-Doppelfunde `/Dev/w`, `/Dev/z`; Exit 1 | ok |
| 13 | (f) Randfälle am Image mit `variables: [HOME]` | still bzw. dokumentiert | still: `$HOMEDIR/x/y`, `$HOME/.config/x`, `${HOME:-/tmp}/x/y`, `\$HOME/x/y`, nacktes `$HOME`, `$HOME//x`, `${HOME}x/y`, `x$HOME/y`, `$$HOME/x`, URL, `$OTHER/x`. Gemeldet: `**$HOME/x**` → `$HOME/x**` (F-4, dokumentierte Grenze), `$HOME/ok/pfad`, und die drei Formen mit erstem Segment `tmp` je einmal in voller Form | ok |
| 14 | (f) Config-Rand am Image | Exit 2 | `exempt-targets: ["$HOME/**"]` ohne `variables`: Exit 2; `["$OTHER/**"]` bei `variables: [HOME]`: Exit 2; `variables: ["$HOME"]`, `[""]`, `["1A"]`: Exit 2, Meldung nennt `hostpaths.variables`. Gegenprobe `["$HOME/Dev/x/**"]` mit `variables: [HOME]`: Exit 1, genau dieser Fund (`$HOME/Dev/x/y.log`) entfällt, vier Funde bleiben | ok |
| 15 | `--print-config` am Image | Gerüst führt den Schlüssel | Zeile 73 `variables: [HOME]` (aus, solange abwesend), Zeile 74 `exempt-targets` nennt „einer in variables genannten Variablen" | ok |

## DoD-Punkte

| Punkt | Befund |
|---|---|
| Test referenziert (Happy/Boundary/Negative) + Rot-Beleg | Erfüllt. Happy `TestHostpathsVariablesHappy`, Boundary `TestHostpathsVariablesBoundary`, Negative `TestHostpathsVariablesNegatives` / `…ExemptTargets` / `TestDecode_HostpathsVariables`. Rot-Beleg nachgefahren (Messung 2–8): Der `$HOME`-Test wird ohne die Änderung mit „kein Befund statt `hostpath-forbidden`" rot, nicht aus einem Folgefehler. Akzeptanz in `docs/user/abdeckung-tests.md` eingetragen (`make test` hält sie gegen die Ableitung, grün). |
| Fassung nachgezogen | Erfüllt, mit einem Befund (V-1 unten). Lastenheft 0.108.0 in Kopf und Historie, Beschreibung Punkt 3, Wortgrenze, Ventil, Kriterien, Out-of-Scope; Spezifikation Schritt 2, Schritt 4, beide Schema-Zeilen, Historie; ADR-0111 samt Index-Zeile. |
| Spiegel (Gerüst) | Erfüllt (Messung 15, am Image, nicht nur am Quelltext). |
| `make gates` grün | Erfüllt (Messung 1). |
| Review-Report liegt vor | Erfüllt. F-1…F-5, F-7 im Nachzug-Commit bzw. in den Plan-Änderungen; die Wirkung zeigen Messung 3 (allgemeine Verdrängung), 13 (F-4) und die Plan-Änderungen vor §4. |
| Doku-Update Schnittstelle | Erfüllt (Spec, Gerüst; Handbuch laut Plan Release-Prep). |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | Offen — Slice liegt in `in-progress/`, §7 „Wird bei der Closure gefüllt", Haken leer. Zustand vor Closure, keine Verletzung. |
| Closure-Trigger §5 | Beide Hälften gemessen: Fixture Issue 12 meldet alle drei Formen mit `[HOME]` (Messung 12), ohne Schlüssel byte-gleich gegen `v0.87.1` (Messung 10, 11). |

## Prüfung (b): Fitness-Funktion in ADR-0111

Bestätigt. „`$HOME`-Test ohne Variablen-Erkennung rot (leere Funde)" = Messung 2; „Verdrängungs-Test ohne die Verdrängung rot (zweiter Unix-Fund innerhalb des Variablen-Pfads)" = Messung 3 (Unix-Funde zusätzlich). „Jede Stelle einzeln gebrochen und gesehen" trägt für die in der ADR genannten Stellen und darüber hinaus: Punkt-Segment (4), beide Config-Prüfungen (5, 6), schließendes Anführungszeichen (7), Wortgrenze (8) — jede wird von mindestens einem Test aus eigener Meldung rot. Nicht gebrochen habe ich die Reihenfolge der Verdrängung zwischen Tilde und Variable im Einzelnen (nur gemeinsam über Mutation B).

## Prüfung (c): Spiegel

Code, Spezifikation, Gerüst und ADR tragen dieselbe Form: Variablen `NAME` aus der Liste, `$NAME`/`${NAME}`, optional schließendes Anführungszeichen, Punkt-Segment-Regel der Tilde, Wortgrenze der Tilde plus `$` und Backslash, Meldung ohne schließendes Anführungszeichen, früher beginnender Treffer gewinnt, `$`-Glob nur mit gelisteter Variable. Version 0.108.0 in Kopf und Historie. Die generische „Bekannte Grenze" zur Hervorhebung am Pfadrand (Spezifikation, Abschnitt zu `DC-FA-HOST-001.a`) deckt `**$HOME/x**` (Messung 13) mit, ohne eigene Variablen-Zeile.

## Befunde / DoD-Verletzungen

Keine DoD-Verletzung. Ein Spiegel-Befund und Hinweise:

- **V-1 (LOW, Fassung):** Das Lastenheft-Kriterium *Variablen-Home (Boundary)* sagt „Given dieselben Zeilen **ohne** `hostpaths.variables` … then kein Befund". Das trifft nur zu, solange das erste Segment kein konfigurierter Präfix-Name ist. Gemessen (Messung 11): im Fixture mit Präfix `Dev` meldet `${HOME}/Dev/z` ohne den Schlüssel nur den Unix-Rest hinter der Klammer und `"$HOME"/Dev/w` den hinter dem Anführungszeichen — der alte Unix-Treffer hinter `}` bzw. `"`, in v0.87.1 genauso. Das ist mit dem zweiten Halbsatz („Befundsatz byte-identisch") und dem Code vereinbar, aber der erste Halbsatz überspricht für diese Zeilen; der Boundary-Test (`TestHostpathsVariablesBoundary`) nutzt `Dev` bei Standard-Präfixen und sieht das nicht. Gleiches gilt für „melden nichts" im Kontext von Slice §1 und ADR-0111 (dort gilt es nur bei Standard-Präfixen, F-7 benennt das Präfix `Dev` nur im Plan). Kein Code-Fehler.
- ADR-0111 steht auf `Proposed`; der Übergang nach `Accepted` ist Closure-Arbeit und gehört vor den `git mv`.
- Die Boundary-Aussage „byte-identisch" ist am Eigenbestand mit leerem Befundsatz (Standard-Präfixe) und mit erweiterter Präfixliste (16 Funde) sowie am Issue-Fixture belegt, nicht an Fremd-Repos.
- Plan §3 erwartet für `variables: [HOME]` über den Eigenbestand „die zwei Dateien oben". Gemessen an der Repo-Kopie sind es fünf Dateien (`ADR-0111`, `slice-280`, Welle-92-Plan, der Review-Report, `spec/lastenheft.md`); die Mehrfunde sind die neuen eigenen Dokumente dieses Slice. Das gibt der später möglichen Selbst-Aktivierung (ADR-0111, Re-Evaluierungs-Trigger) ein Gewicht, die Aktivierung selbst ist laut Plan ausgeschlossen und nicht erfolgt (`.d-check.yml` unverändert).
- Die Risiken §6 tragen noch keinen Ausgang (Closure). Datenlage: Risiko 3 (Meldung nennt den Schlüssel nicht) ist an Messung 12 sichtbar — die Meldung ist `hostpath-forbidden` mit Pfad, ohne Hinweis auf `variables`.

## Nicht geprüft

- `make verify-closure-notes`, Closure-Notiz, Register-Eintrag, drei Paarungen (Slice nicht geschlossen).
- `make image-test`, `make fullbuild`, `make completeness-check` (nicht Teil der DoD dieses Slice).
- Verhalten gegen reale Fremd-Repos mit `$HOME` in Prosa (kein Fremdkorpus verfügbar); `suggest.go`/`run.go` auf Aussagen über Präfixliste und Ventil (Plan §3) habe ich nicht gesondert gelesen — `suggest.go` liegt unter `internal/hexagon/core/app/` und wurde nicht gelesen.
- Das Vorher-Image wurde aus dem Tag-Archiv gebaut, nicht gegen die veröffentlichte GHCR-Fassung verglichen.
