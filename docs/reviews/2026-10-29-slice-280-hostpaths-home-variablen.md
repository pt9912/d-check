# Review slice-280: `hostpaths` erkennt `$HOME` und `${HOME}` als Home-Verweis

**Review-Art:** Code (Diff gegen Plan, ADR und Hard Rules; nicht gegen die DoD)
**Gegenstand:** `slice-280`, Range `ab7de630..HEAD` (f9a1b893 feat, 9e1f9fe7 docs/spec)
**Skill:** `reviewer.md` @ 1.20.0
**Modell-ID:** claude-sonnet-5-5
**Datum:** 2026-10-29
**Eingangs-Kontext:** `slice-280` (§1 Ziel und Abgrenzung, §3 Plan inkl. Negativliste und Plan-Änderung, §6 Risiken), `ADR-0111` (Proposed), `welle-92`, `DC-FA-HOST-001`, `DC-QA-02`, `AGENTS.md` §3.2/§3.4/§3.7/§3.8. Vorheriger Review am Zwilling: `slice-279`.

## Findings

### F-1 — MEDIUM — Die Plan-Änderung, die ADR-Grenze und ein Test-Kommentar begründen die Verdrängung mit einer Prämisse, die für zwei der drei Schreibweisen falsch ist; der Test fährt nur die Form, in der die Verdrängung nichts tut

- **quelle:** Reviewer-Skill Prüffrage 20 (Begründung trifft am Gegenstand nicht zu) und 13 (Negativtest); `slice-280` §3 Negativliste „ein Variablen-Pfad mit Host-Präfix-Namen als erstem Segment (genau ein Befund, in voller Form)"
- **pfad:** `docs/plan/planning/in-progress/slice-280-hostpaths-home-variablen.md` · „weil vor dem Präfix das letzte Zeichen des Variablen-Namens steht (Wortgrenze)"; `docs/plan/adr/0111-hostpaths-variablen-home-opt-in.md` · „trifft das Unix-Muster ohnehin nicht, weil davor das letzte Zeichen des Variablen-Namens steht"; `internal/hexagon/core/rules/hostpaths_variables_test.go` · „träfe das Unix-Muster ohnehin nicht — vor `/<Präfix>/` steht das `E` von `HOME`"
- **befund:** Das gilt nur für `$HOME/<Präfix>/x`. Bei `${HOME}/<Präfix>/x` steht `}` und bei `"$HOME"/<Präfix>/x` steht `"` vor dem Schrägstrich; beide Zeichen liegen nicht in der Ausschlussklasse `[^A-Za-z0-9_.:/-]` von `unixHostpathRE`, das Unix-Muster trifft dort also sehr wohl. Der Code (`insideSpan(m[4], homes)`) verdrängt korrekt, aber `TestHostpathsVariablesVerdraengtUnix` prüft die Präfix-Namen-Form nur mit `$HOME/<Präfix>/d-check/README.md`, wo die Verdrängung ohne Wirkung ist; die `=`-Form deckt den Mechanismus, nicht die beiden Schreibweisen, in denen der Negativlisten-Fall real auftritt.
- **Failure-Szenario:** Jemand zieht die Verdrängung oder die Vorbedingung des Variablen-Musters enger (die Plan-Änderung und die ADR laden dazu ein: „ohnehin nicht"). Alle Tests bleiben grün, und `${HOME}/<Präfix>/x` meldet danach zweimal, `/<Präfix>/x` und `${HOME}/<Präfix>/x`. Die ADR-Grenze („greift nur, wenn … `=`") und der Plan-Satz lesen sich dann als Beleg, dass die Form nicht vorkommen kann.
- **verifizierbar:** ja. Ein Test mit `${HOME}/<Präfix>/x/y` und `"$HOME"/<Präfix>/x/y` (erwartet je genau ein Ziel) rot ohne die Verdrängung.
- **klasse:** `begruendung-gilt-nur-fuer-eine-schreibweise`

### F-2 — MEDIUM — Mit gesetztem `variables` ändert sich das Ziel eines bisherigen Unix-Funds; ein Bestands-Glob in `exempt-targets` trifft ihn dann nicht mehr, und kein Kriterium nennt das

- **quelle:** `DC-FA-HOST-001` (Kriterien „Ziel-Ventil", „kein Verlust"), `DC-QA-02`; Prüffrage 11/18 des Skills (Grenzen-Liste ohne ihre größte Lücke)
- **pfad:** `internal/hexagon/core/rules/hostpaths.go` · „if !insideSpan(m[4], homes) {"
- **befund:** Ohne den Schlüssel meldet `${HOME}/<Präfix>/x` (und `"$HOME"/<Präfix>/x`, `$HOME/a=/<Präfix>/x`) bisher das Ziel `/<Präfix>/x` bzw. `/<Präfix>/x`; mit `variables: [HOME]` entfällt dieser Unix-Fund zugunsten des Variablen-Funds `${HOME}/<Präfix>/x`. Ein vorhandenes Ventil-Glob `/<Präfix>/**` exempt den Fund dann nicht mehr. Lastenheft-Kriterium „Variablen-Home (Boundary)" sichert Byte-Identität nur *ohne* den Schlüssel; zum Wechsel des Ziels *mit* ihm schweigen Lastenheft, Spezifikation (Schritt 4), ADR-Konsequenzen und das Handbuch-Gerüst. Ein Kriterium „kein Verlust" nach dem Muster der Tilde fehlt für die Variable.
- **Failure-Szenario:** Ein Adopter hält `exempt-targets: ["/<Präfix>/shared/**"]`, schaltet `variables: [HOME]` ein, und sein grüner Lauf wird rot, ohne dass ein neuer Pfad hinzukam. Die Ursache steht in keiner Zusage; die Abhilfe (zusätzlich ein `$`-/`${…}`-Glob) ist nirgends genannt.
- **verifizierbar:** ja. Test: `${HOME}/<Präfix>/shared/x` mit Glob `/<Präfix>/shared/**` und `variables: [HOME]` meldet; ohne Schlüssel meldet nichts.
- **klasse:** `schluessel-aendert-ziel-bestehender-funde`

### F-3 — LOW — Zwei Home-Funde in einer Zeile können überlappend doppelt melden (Variable und Tilde bzw. Unix)

- **quelle:** `DC-FA-HOST-001` (Beschreibung: „in voller Form und genau einmal gemeldet"); Spezifikation `DC-FA-HOST-001.a` Schritt 2 („Unix-Treffer … entfällt")
- **pfad:** `internal/hexagon/core/rules/hostpaths.go` · „if !insideSpan(m[4], homes) {"
- **befund:** Verdrängt wird nur der Unix-Treffer. Ein Tilde-Treffer innerhalb eines Variablen-Treffers (`$HOME/a,` plus Tilde-Pfad), ein Variablen-Treffer innerhalb eines Tilde-Treffers (Tilde-Pfad plus `,$HOME/b/c`) und ein Variablen-Treffer innerhalb eines Unix-Treffers (`/<Präfix>/a,$HOME/b/c`) liefern je zwei Befunde mit überlappendem Text. Die zweite Richtung der Tilde-/Unix-Überlappung war schon vorher so; neu ist, dass die Variable hinzukommt und die Spezifikation „genau einmal" nur für das Verhältnis zum Unix-Muster ausschreibt. Der Komma-Fall ist erreichbar, weil `,` keine Vorbedingungs-Sperre ist (der Doppelpunkt schon, deshalb sind PATH-Listen sicher).
- **verifizierbar:** ja, per Unit-Test mit diesen drei Formen.
- **klasse:** `ueberlappende-treffer-verschiedener-muster`

### F-4 — LOW — Negativlistenpunkt `**$HOME/x**` ist im Code das Gegenteil, ohne dass die Plan-Änderung ihn nennt

- **quelle:** `AGENTS.md` §6 Schritt 4 (Plan-Änderung gehört vor den Code); `slice-280` §3 Negativliste
- **pfad:** `docs/plan/planning/in-progress/slice-280-hostpaths-home-variablen.md` · „`**$HOME/x**` und `_$HOME/x_` (Markdown-Hervorhebung am Rand)"; `internal/hexagon/core/rules/hostpaths_variables_test.go` · „func TestHostpathsVariablesEmphasisLikeTilde"
- **befund:** Der Plan führt `**$HOME/x**` unter „kein Fund"; der Test erwartet `$HOME/x/y**` als Fund (Stern-Zeichen im Ziel). Nur `_$HOME/x_` ist wie geplant still. Die Plan-Änderung vor §4 behandelt ausschließlich den Präfix-Namen-Fall. Die Wirkung ist durch die bekannte Grenze der Spezifikation („Markdown-Hervorhebung am Pfadrand gehört zum gemeldeten Pfad … ein Glob mit `**` am Ende schon") und durch die ADR-Grenze gedeckt, die Abweichung vom Plan aber nicht ausgewiesen.
- **verifizierbar:** nein (Doku gegen Doku).
- **klasse:** `negativliste-umgekehrt-ohne-plan-aenderung`

### F-5 — LOW — Lastenheft-Kriterium „Ziel-Ventil" und die Wortgrenzen-Aussage der Beschreibung sind nicht auf die Variable nachgezogen

- **quelle:** `DC-FA-HOST-001`; Skill-Anker „Doku-Drift"
- **pfad:** `spec/lastenheft.md` · „then entfällt genau dieser Befund; ein Windows- oder UNC-Fund bleibt bestehen" und „vor einem Home-relativen Pfad auch keine Tilde"
- **befund:** Die Beschreibung und die Config-Rand-Zeile nennen den Variablen-Fund im Ventil, das Akzeptanzkriterium „Ziel-Ventil" spricht weiter nur von „Unix- oder Home-relativen Pfad". Der Wortgrenzen-Satz nennt für die Variable die zusätzlichen Sperren `$` und Backslash nicht (nur die Spezifikation tut es). Das Ventil für die Variable ist über einen Test belegt, aber kein Kriterium trägt es.
- **verifizierbar:** nein.
- **klasse:** `spiegel-teilweise-nachgezogen`

### F-6 — INFO — Die Fehlermeldung für ein Ventil-Glob ohne `/`/`~` ist umformuliert

- **quelle:** `DC-QA-02` (nur Befundsatz zugesagt)
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go` · „beginnt weder mit '/' noch mit '~' (noch mit einer in hostpaths.variables genannten Variablen)"
- **befund:** Der Exit-2-Text für den Altfall ändert sich auch ohne den neuen Schlüssel. Die Zusage der Byte-Identität gilt dem Befundsatz, nicht der Konfigurationsfehler-Meldung; der bestehende Test prüft nur den Schlüsselnamen. Dokumentationswürdig, kein Verstoß.
- **verifizierbar:** ja, `make test` ist grün.
- **klasse:** `fehlertext-aenderung-ohne-schluessel`

### F-7 — INFO — Die Ursache im Anlass-Absatz setzt ein konfiguriertes Präfix `Dev` voraus

- **quelle:** Skill Prüffrage 20
- **pfad:** `docs/plan/adr/0111-hostpaths-variablen-home-opt-in.md` · „sieht vor `/Dev…` das `E` von `HOME`"
- **befund:** Im Default heißt das Präfix `Development`, `/Dev/` ist keines. Die Erklärung trifft nur zu, wenn das Feld-Repo `Dev` als Präfix führt. Die Entscheidung (ein eigener Schlüssel) bliebe bei falscher Ursache dieselbe, die Variable ist unabhängig vom Präfix nie erkannt worden.
- **verifizierbar:** nein.
- **klasse:** `ursache-setzt-konfiguration-voraus`

## Beantwortung der Schwerpunkte

1. **Negativliste vollständig umgesetzt und getestet?** Bis auf den Präfix-Namen-Fall (F-1) und `**$HOME/x**` (F-4) ja:
   - `$HOMEDIR/x`, `${HOME:-/tmp}/x`, `${HOME}x/y`, `x$HOME/y`, `$$HOME/x`, `\$HOME/x`, URL, `_$HOME/x_`, `$OTHER/x`, `$HOME//x`, Zeilenumbruch, Fence: `TestHostpathsVariablesNegatives`.
   - Nacktes `$HOME` und `$HOME/.config/x`: `TestHostpathsVariablesBoundary`.
   - Ohne den Schlüssel kein Fund: ebenda (leere Konfiguration).
   - `$`-Glob ohne Variable, anderer Name, Namens-Präfix, ohne Schrägstrich, leerer/ungültiger Name: `TestDecode_HostpathsVariables`.
   - Der Lauf über den Eigenbestand und der Vorher/Nachher-Vergleich gegen `v0.87.1` liegen nicht im Diff; sie sind Messung und Verifier-Sache.
2. **Nimmt die neue Erkennung dem Bisherigen etwas weg?**
   - Ohne Schlüssel: nein. `varRE` ist `nil`, `variableMatches` liefert nichts, `homes` ist die Tilde-Liste unverändert; die Windows-Schleife in `windowsPaths` liefert dieselbe Reihenfolge (Laufwerk, dann UNC) wie vorher.
   - Unix- und Tilde-Funde mit Schlüssel: ein Fund geht nicht verloren (nur ein Unix-Treffer *innerhalb* eines Variablen-Treffers entfällt), aber sein Ziel ändert sich; Folge für Bestands-Globs in F-2. `TestHostpathsVariablesOhneUnixVerlust` deckt nur den Nachbar-Fall.
3. **Regex-Grenzen gegen Lastenheft, Spezifikation, ADR:** Wortgrenze (`[^A-Za-z0-9_.:/~$\\-]`), Punkt-Segment (erste Zeichenklasse wie `tildeRE`), `"`-Form (`"?` hinter der Alternation), `${NAME}` und längere Namen (`HOME|HOMEDIR` fällt über die Alternation korrekt auf `HOMEDIR` zurück) stimmen mit Spezifikation Schritt 2 überein. `regexp.QuoteMeta` schützt Namen mit Metazeichen; die Config-Validierung lässt solche ohnehin nicht zu. `"$HOME"/x`: das gemeldete Ziel entfernt genau das eine `"`, weil die Restpfad-Klasse `"` ausschließt.
4. **Spiegel untereinander:** Lastenheft 0.108.0 (Beschreibung, Kriterien, Out-of-Scope, Historie), Spezifikation (Schritt 2 und 4, zwei Schema-Zeilen, Historie), Gerüst-Kommentar in `config_template.go`, `HostpathsConfig`-Kommentar und ADR-Index sagen dasselbe wie der Code, außer F-5 (Kriterium/Wortgrenze) und F-1 (ADR-Grenze).
   - `internal/hexagon/core/app/suggest.go` liest `hostpaths` nur als Modulname (Zeilen mit `optIn` und `modules`), `internal/hexagon/core/rules/run.go` ruft nur `CheckHostpaths` auf: beide korrekt unberührt.
   - `harness/sensors/doc-check.md` nennt `hostpaths` nur in der Modulliste, ohne Aussage über Präfixe oder Ventil: korrekt unberührt.
   - `docs/user/benutzerhandbuch.md` (Beispiel, Modultabelle) und `docs/user/operations.md` sind Release-Prep (`AGENTS.md` §5 Zeile 17): kein Finding. `docs/user/abdeckung-tests.md` trägt die neue Zeile.
5. **Zerlegung `CheckHostpaths` / `applyHostpaths`:** verhaltensgleich ohne den Schlüssel (siehe 2). In `applyHostpaths` ändert sich die Fehlerreihenfolge nur mit gesetztem `variables` (Variablen-Fehler vor Ventil-Fehler); ein einziger Fehlertext wurde umformuliert (F-6). `make lint` meldet 0 Issues.
6. **Falsche oder doppelte Funde:** doppelt melden kann F-3. Falsche Funde: `$HOME/…`-Platzhalter wie `$HOME/…` (Auslassungszeichen als erstes Zeichen) melden, wie bei der Tilde; das ist die zugesagte Erkennung, keine neue Klasse. Eine zweite Variable in derselben Zeile hinter dem ersten Treffer (`$HOME/a,$WORK/b/c`) ist durch den ersten Treffer mitgedeckt; ein Ventil-Glob auf den ersten Fund verschluckt sie, wie bei der Tilde.

## Negativbefunde (geprüft, ohne Befund)

- **Suppression (§3.2):** keine `//nolint` im Diff. Kein Befund.
- **Hexagon-Richtung (`ADR-0005`):** `hostpaths.go` importiert nur `model`; `configyaml.go` bleibt im Adapter. Kein Befund.
- **Spec-Straten nie abwärts (§3.4):** Lastenheft und Spezifikation nennen weder ADR, Slice noch Welle; die Historie-Zeilen tragen „Anlass: Feldmeldung eines Adopters" ohne Kennung. Kein Befund.
- **Kommentar-Klassen (§3.7):** neue Kommentare in `hostpaths.go`, `configyaml.go`, `config_template.go`, `model/config.go` tragen Zusage, Kopplung oder Grenze mit `DC-FA-HOST-001` als einzigem Herkunftsfeld; die Test-Kommentare tragen keine Slice-Nummern oder Mess-Labels (die Sachaussage im Verdrängungs-Test ist F-1). Kein weiterer Befund.
- **Modul liest Eingaben, die es nicht scannt (§3.8, Prüffrage 15):** `variables` und `exempt-targets` sind Konfiguration, keine Zieldateien; die Grenze für `${HOME:-…}`, `%NAME%` und die nackte Variable steht im Lastenheft. Kein Befund.
- **Netz/Schreibzugriff (`DC-QA-03`):** nur reine Funktionen. Kein Befund.
- **Eigene Erkennung statt Produktes (Prüffrage 21):** Pfad-Erkennung innerhalb des Moduls selbst, im Muster der Tilde. Kein Befund.
- **ADR-0111:** Re-Evaluierungs-Trigger, Alternativen-Tabelle, `Schärft:` mit Kennung (`DC-FA-HOST-001.a`), Provenance-Marker zeigt statt zu begründen, Index-Zeile Proposed mit passendem Datum. Kein Befund außer F-1/F-7.
- **Gates als Beleg (kein Verifier-Urteil):** `make test` Exit 0, `make lint` 0 Issues, `make doc-check` ohne Befund.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 2 | 3 | 2 |

## Verdikt

Nicht blockierend außer F-1 und F-2. F-1 verlangt vor dem Merge einen Test mit `${HOME}/<Präfix>/…` und `"$HOME"/<Präfix>/…` und eine Korrektur der falschen Prämisse in Plan-Änderung, ADR-Grenze und Test-Kommentar. F-2 verlangt, dass Lastenheft oder Spezifikation und ADR benennen, dass mit gesetztem `variables` ein bisheriger Unix-Fund sein Ziel wechselt, samt Test. F-3 bis F-7 sind optional beziehungsweise Doku-Nachzug.
