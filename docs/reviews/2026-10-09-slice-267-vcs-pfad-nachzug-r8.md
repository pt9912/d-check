# Review R8: slice-267, `vcs` lässt einen reinen Pfad-Nachzug durch

- **Review-Art:** Code. Geprüft wird `b60bf321..dc1fa440`: `056dc4cf` (`fix(vcs)`, Auflösung gegen
  die Vereinigung der Pfad-Bäume, Marke vor der normierten Form) und `dc1fa440`
  (`ADR-0103`-Geschichte). Geprüft gegen den Slice-Plan `slice-267`, die Findings aus R1 bis R7
  (`2026-10-09-slice-267-vcs-pfad-nachzug-r1.md` bis `-r7.md`), `ADR-0103`, `DC-FA-VCS-001`,
  `MR-025` und die Hard Rules `AGENTS.md` §3.5, §3.7, §3.8 sowie §5 Regel 13. Die DoD-Abhakung
  prüft dieses Review nicht.
- **Gegenstand:** `slice-267` · `b60bf321..dc1fa440`
- **Skill:** `reviewer.md` @ 1.19.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** Lastenheft 0.102.6 (`DC-FA-VCS-001`, Absatz „Pfad-Nachzug (opt-in)",
  Kriterium „Boundary (Pfad-Nachzug)", Historie 0.102.6), Spezifikation `DC-FA-VCS-001.a`
  Schritt 4, §2-Zeile `vcs.ignore-link-targets` und Historie, `harness/sensors/adr-check.md`
  Grenze 4, `ADR-0103` (Geschichte, Diff gegen `273bebea`). Im Code: `vcs.go` (`CheckVCS`,
  `vcsModified`, `vcsCore`), `vcs_link_targets.go` (`normalizedTargetMark`, `linkTargetResolver`,
  `pathTree`, `normalizedLinkTargetLines`), `pins.go` (`pinWhitespaceRE`), `model/config.go`,
  `cli/config_template.go`, `vcs_pfad_nachzug_test.go`. Vorherige Findings am Modul: R1 bis R7.
- **Proben:** Image per `make build` vom Stand `dc1fa440` (`d-check:latest`). Alle Proben in einem
  Wegwerf-Klon im Scratchpad: Änderung, Commit, Lauf mit dem Aufruf von `make adr-check`
  (`--enable vcs` plus `FOCUS_DISABLE`, `--range HEAD~1..HEAD`) gegen die `.d-check.yml` des Repos
  (`ignore-link-targets: true`), danach `git reset --hard origin/main`. Klon danach entfernt.

## Proben am gebauten Image

| # | Änderung im Klon | Ergebnis |
|---|---|---|
| P1 | `git mv docs/user/releasing.md docs/user/maintainer/releasing.md`, die vier Links in `ADR-0014` nachgezogen | 0 Befunde, Exit 0 |
| P2 | `git rm harness/README.md`, keine ADR geändert (6 `Accepted`-ADRs verlinken sie) | 0 Befunde, Exit 0 (mit `ecb8443e` laut Commit-Botschaft 6 Befunde) |
| P3 | ohne Umzug: Ziel in `ADR-0014` Zeile 21 von `../../../docs/user/releasing.md` auf `../../../docs/user/operations.md` | `core-drift-vcs`, Exit 1 |
| P3b | wie P1, ein Link zeigt stattdessen auf `operations.md` | `core-drift-vcs`, Exit 1 |
| P4 | ohne Umzug: Ziel auf bloßes `releasing.md` (R7 M-1) | `core-drift-vcs`, Exit 1 — R7 M-1 für Text ohne Steuerzeichen eingelöst |
| P4n | ohne Umzug: Ziel auf `\x00releasing.md` (NUL-Byte vor dem Dateinamen) | **0 Befunde, Exit 0**; Default-Lauf (Modul `links`) meldet `target-missing` Zeile 21 |
| P5 | Commit 1 legt `docs/tmp/releasing.md` an; Commit 2 löscht sie und zieht `ADR-0014` Zeile 21 dorthin | 0 Befunde, Exit 0 (benannte Grenze, wie dokumentiert) |

## Findings

### M-1 — Die Marke ist ein Zeichen, das roher Text tragen kann; „gleicht ihr nie" stimmt nicht

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 20 und 18 · `AGENTS.md` §5 Regel 13 · `DC-FA-VCS-001.a` Schritt 4
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`roher Text trägt das`" /
  „`const normalizedTargetMark = "\x00"`"; `spec/spezifikation.md` · „die Marke trägt kein roher
  Text, ein Ziel, das nur so aussieht, gleicht ihr nie"; `ADR-0103` Geschichte · „ein Ziel aus
  bloßem Dateinamen gleicht ihr nicht"
- **Befund:** Die Marke ist das NUL-Byte. Ein Markdown-Blob kann es tragen: git speichert die Datei,
  die Link-Erkennung nimmt `\x00` als Zielzeichen an (`linkDestRE` `[^\s<>]+`), der Resolver löst
  `\x00releasing.md` nicht auf und lässt es roh, und `pinWhitespaceRE` (`\s+`) berührt es nicht.
  Damit gleicht der rohe HEAD-Text der normierten BASE-Form byte-genau (P4n). Die Entscheidung für
  eine Marke trägt (P4 ist jetzt Drift); die Begründung „roher Text trägt das Zeichen nicht" und
  das „nie" in der Spezifikation treffen am Gegenstand nicht zu, und keine Grenzen-Liste nennt den
  Restpfad. Die Spezifikation illustriert die Marke außerdem als `‹…›` um das Ziel, der Code setzt
  ein unsichtbares Präfix — ein Leser kann aus dem Beispiel nicht ablesen, welche Zeichen die
  Garantie tragen.
- **Failure-Szenario:** Ein Commit setzt in einer `Accepted`-ADR das Ziel
  `../../../docs/user/releasing.md` auf `\x00releasing.md` (etwa durch ein Werkzeug, das
  Steuerzeichen einschleust). `make adr-check` und der `pre-commit`-Hook bleiben grün (P4n). In
  diesem Repo fängt `make doc-check` den toten Link; ein Konsument, der `vcs` ohne `links` fährt,
  bekommt das Grün, das ihm Schritt 4 mit „gleicht ihr nie" ausdrücklich nicht verspricht.
- **Warum nicht HIGH:** Die Form verlangt ein Steuerzeichen im Zielslot, kein realistischer
  Tippfehler; sichtbarer Text ändert sich nicht, und das Partner-Gate `links` meldet sie — dieselbe
  Sicherung wie bei der benannten BASE-only-Grenze.
- **Verifizierbar:** ja. P4n; ein Kern-Fall mit `"[S](\x00s.md#a)"` als HEAD-Ziel und erwartetem
  Befund läuft heute rot.
- **Klasse:** `grenze-gegen-gegenstand`

### L-1 — Der Funktionskommentar von `normalizedLinkTargetLines` nennt die Marke nicht

- **Kategorie:** LOW
- **Quelle:** `MR-025` (Spiegel einer Semantik-Änderung) · `AGENTS.md` §3.7 („Ein Kommentar
  beschreibt, was da ist")
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`es wird durch
  // Dateiname und Anker ersetzt, sodass ein reiner Pfad-Nachzug denselben Text`" (Kopf von
  `normalizedLinkTargetLines`) und „`im Stand der Datei auflöst`"
- **Befund:** Der Kopfkommentar sagt, ein Ziel werde „im Stand der Datei" aufgelöst und durch
  Dateiname und Anker ersetzt. Seit `056dc4cf` löst es gegen die Vereinigung beider Stände auf und
  wird zu Marke, Dateiname und Anker; der Kommentar von `linkTargetResolver` sagt das richtig, der
  der aufrufenden Funktion nicht.
- **Verifizierbar:** nein (Lesen).
- **Klasse:** `spiegel-nicht-nachgezogen`

### I-1 — Die benannte BASE-only-Grenze hat keinen Fall in der Suite

- **Kategorie:** INFO
- **Quelle:** Prüffrage 13 (unter MEDIUM-Schwelle: das Verhalten ist eine benannte Grenze, kein
  Vertrag)
- **Pfad:** `internal/hexagon/core/rules/vcs_pfad_nachzug_test.go` · „`func
  TestVCSIgnoreLinkTargetsUmzug(t *testing.T) {`"
- **Befund:** Die vier Fälle decken Umzug, unveränderten Link auf gelöschte Datei, bloßen
  Dateinamen und ein in beiden Ständen fehlendes Ziel. Der Fall „Nachzug auf eine Datei, die nur in
  BASE existiert ⇒ 0 Befunde" (P5) steht nur in Prosa. Ändert eine spätere Härtung ihn zu Drift,
  bleibt die Suite grün und die Grenz-Texte in Spezifikation, Lastenheft, Sensor-Datei und
  Kommentar sind still falsch. Richtung fail-safe.
- **Verifizierbar:** ja. P5.
- **Klasse:** `grenze-ohne-probe`

## Status der R7-Findings

- **R7 M-1:** für Text ohne Steuerzeichen aufgelöst (P4 Drift; Test „neues ziel nur dateiname").
  Restpfad über das Markenzeichen selbst: M-1 oben. Die Grenz-Texte nennen jetzt „in BASE oder HEAD
  existiert", passend zum Code.
- **R7 M-2:** aufgelöst. `TestVCSIgnoreLinkTargetsUmzug` legt getrennte Bäume an. Ich habe die
  Mutationen gedanklich gegen die Fälle gelegt: Baum nur aus BASE ⇒ „nachzug nach umzug" rot (HEAD
  löst `maintainer/releasing.md` nicht auf); Baum nur aus HEAD ⇒ dieselbe Probe rot (BASE löst
  nicht auf); getrennte Bäume (Vorzustand) ⇒ „unveraenderter link auf geloeschte datei" rot; Marke
  leer ⇒ „neues ziel nur dateiname" rot. Das deckt sich mit den Mutationsangaben in `056dc4cf`;
  selbst gefahren habe ich die Mutationen nicht.
- **R7 L-1:** aufgelöst (Modell-Kommentar, `opaqueLines`, `--print-config`-Vorlage). Ein weiterer
  Spiegel: L-1 oben.
- **R7 L-2:** aufgelöst. Spezifikation und Sensor-Datei nennen den zeilenweise falsch gepaarten
  Code-Span.
- **R7 I-1, I-2:** unverändert, keine Änderung erwartet.

## Bewertung der Vereinigung (Frage 1)

- **Korrektheit:** `pathTree(baseAll ∪ headAll)`; beide Seiten bekommen denselben Resolver. Ein
  unveränderter Link ergibt damit auf beiden Seiten denselben Text, unabhängig davon, ob sein Ziel
  gelöscht, angelegt oder verschoben wurde (P2). Der Fehler aus `ecb8443e` ist strukturell weg.
- **Neue Durchlass-Pfade:** Gegenüber getrennten Bäumen löst zusätzlich auf (a) ein HEAD-Ziel, das
  nur in BASE existiert — benannt (P5) — und (b) ein BASE-Ziel, das nur in HEAD existiert (Link in
  BASE war tot, Datei entsteht in der Range). (b) geht nur durch, wenn HEAD auf ein auflösendes Ziel
  **gleichen Namens** zeigt; das ist die bereits benannte Klasse „gleicher Name, nicht dieselbe
  Datei" (R7 I-2), kein neuer Pfad mit sichtbarer Inhaltsänderung. Ein Inhaltswort, ein fehlendes
  Ziel oder ein anderer Name bleibt Drift (P3, P3b, P4).
- **Marke und Whitespace:** `pinWhitespaceRE` ist `\s+` (Go/RE2: `[\t\n\f\r ]`), NUL fällt nicht
  darunter; die Normalisierung wirkt auf die Marke nicht. Ob roher Text sie tragen kann: ja, M-1.

## Negativbefunde (geprüft, ohne Befund)

- **Hexagon-Richtung (ADR-0005):** keine neuen Imports.
- **Kommentar-Klassen (§3.7):** neue Kommentare (`normalizedTargetMark`, `linkTargetResolver` samt
  `GRENZE`, Testkommentar) tragen Zusage bzw. Grenze; keine Review-Historie, keine Befund-Nummern.
  Die Zusage der Marke ist inhaltlich falsch (M-1), nicht klassenlos.
- **ADR-Immutabilität (§3.5):** `git diff 273bebea..HEAD` an `ADR-0103` enthält nur hinzugefügte
  Zeilen in der Geschichte-Tabelle; die Status-Zeile ist `Accepted`, der Kern unverändert.
- **Lastenheft:** Version 0.102.6 im Kopf und als Historie-Zeile; Absatz „Pfad-Nachzug" und
  Kriterium „Boundary (Pfad-Nachzug)" („auf keine in BASE oder HEAD existierende") decken sich mit
  dem Code und mit P3–P5.
- **Spezifikation:** Schritt 4, §2-Zeile und Historie nennen die Vereinigung und die BASE-only-Grenze;
  §2-Zeile sagt jetzt „Datei mit anderem Namen" (R7 M-1-Teil zur §2-Zeile eingelöst).
- **Sensor-Datei `adr-check.md` Grenze 4:** nennt Vereinigung, BASE-only-Grenze und den falsch
  gepaarten Code-Span; deckt sich mit dem Code.
- **`--print-config`-Vorlage / `model/config.go`:** beschreiben Normierung statt Leerung, passend.
- **Determinismus (`DC-QA-02`):** die Vereinigung ist eine Menge, Kandidaten bleiben sortiert.
- **Modul-Grenze (§3.8):** der Baum wird aus `AllPaths` (fail-closed) gebildet, keine neue
  ungescannte Eingabe.
- **Bestand am Image:** echter Umzug von `releasing.md` samt Nachzug ⇒ 0 (P1); Löschen einer
  verlinkten Datei ohne ADR-Änderung ⇒ 0 (P2); Nachzug auf ein anderes Dokument ⇒ Befund (P3, P3b).

## Kategorie-Summary

HIGH 0 · MEDIUM 1 · LOW 1 · INFO 1

## Verdikt

Die Vereinigung ist korrekt und schließt den in `ecb8443e` gefundenen Fehler; der geforderte Bestand
verhält sich am Image wie vertraglich zugesagt. M-1 blockiert nach Skill-Regel: nicht wegen der
Wahrscheinlichkeit der NUL-Form, sondern weil Kommentar, Spezifikation und ADR-Geschichte eine
absolute Garantie („nie") aussprechen, die am Code nicht gilt, und keine Grenzen-Liste den Rest
nennt.
