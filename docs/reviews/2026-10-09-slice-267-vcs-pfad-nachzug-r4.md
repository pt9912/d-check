# Review R4: slice-267, `vcs` lässt einen reinen Pfad-Nachzug durch

- **Review-Art:** Code. Geprüft wird `45176f55..ab64a1e9` — `2115d0a2` (ADR-0103 §Geschichte) und
  `ab64a1e9` (`fix(vcs)`, Antwort auf R3 M-1, L-1, L-2, I-1) — gegen den Slice-Plan `slice-267`,
  gegen die Findings aus R1 bis R3 (`2026-10-09-slice-267-vcs-pfad-nachzug-r1.md` bis `-r3.md`),
  gegen `ADR-0103`, `DC-FA-VCS-001` und die Hard Rules `AGENTS.md` §3.5/§3.7 sowie §5 Regel 13 und
  15. Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-267` · `45176f55..ab64a1e9`; dazu `git diff 273bebea --
  docs/plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md`.
- **Skill:** `reviewer.md` @ 1.19.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** `DC-FA-VCS-001` im Lastenheft 0.102.2 (Absatz „Pfad-Nachzug (opt-in)",
  Kriterium „Boundary (Pfad-Nachzug)", Historie 0.102.0 bis 0.102.2); Spezifikation
  `DC-FA-VCS-001.a` Schritt 4 und §8-Historie; `ADR-0103` samt Geschichte; Baseline `v6.17.0` ·
  `regelwerk/modul-04-adrs.md` §Nachzug ist keine Überschreibung; `harness/sensors/adr-check.md`
  Grenze 4. Im Code: `vcsCore`, `blankLinkTargets`, `overlapsSpan`, `indentedLine` in
  `internal/hexagon/core/rules/vcs.go`; `proseLineSet`, `forEachInlineCodeSpan` in `markdown.go`;
  `vcs_pfad_nachzug_test.go`. Vorherige Findings am Modul: R1 bis R3 zu `slice-267`.
- **Proben** (Image `d-check:latest` per `make build` vom Stand `ab64a1e9`; je Probe ein
  Wegwerf-git-Repo im Scratchpad mit einer `Accepted`-ADR unter `## Kontext`, zwei Commits, Lauf
  `--enable vcs --range HEAD~1..HEAD` mit `vcs.ignore-link-targets: true` und der vcs-Konfiguration
  aus `.d-check.yml`; danach Repos entfernt, Arbeitsbaum unverändert). Zusätzlich Probe 11 in einem
  Wegwerf-Klon dieses Repos mit `FOCUS_DISABLE` aus dem `Makefile`.

  | # | BASE → HEAD | Ergebnis | Markdown-Lesart |
  |---|---|---|---|
  | 01 | ``Siehe [`releasing.md`](../../docs/user/releasing.md).`` → `…/maintainer/releasing.md` | **`core-drift-vcs`, Exit 1** | Link; reiner Nachzug |
  | 02 | `[releasing](../../docs/user/releasing.md)` → `…/maintainer/…` | Exit 0 | Link; reiner Nachzug |
  | 03 | `Text` / Leerzeile / `␠⇥[a](Nicht)` → `…(Immer)` | **Exit 0, still** | eingerückter Code (Leerzeichen + Tab = Spalte 4) |
  | 04 | `<div>` / `[a](Nicht)` / `</div>` → `…(Immer)` | **Exit 0, still** | HTML-Block, sichtbarer Text, kein Link |
  | 05 | `<pre>` / `[a](Nicht)` / `</pre>` → `…(Immer)` | **Exit 0, still** | HTML-Block, sichtbarer Text, kein Link |
  | 06 | `Text` / Leerzeile / `␠␠␠␠[a](Nicht)` → `…(Immer)` | Drift, Exit 1 | eingerückter Code |
  | 07 | ``[a `b](Nicht)` c`` → ``…(Immer)` c`` | Drift, Exit 1 | kein Link |
  | 08 | `<!-- [a](Nicht) -->` → `…(Immer)` | Exit 0 | HTML-Kommentar, unsichtbar |
  | 09 | `[x](<a b.md>)` → `[x](<a c.md>)` | Exit 0 | Link; reiner Nachzug |
  | 10 | ``[`fns[k]`](a.md) Rest`` → `…(b.md) Rest` | Drift, Exit 1 | Link; reiner Nachzug |
  | 11 | `docs/plan/adr/0014-…` mit allen vier Zielen `../../../docs/user/releasing.md` → `…/docs/user/maintainer/releasing.md` (`sed`, 4 Zeilen) | **`core-drift-vcs`, Exit 1** | vier Links, reiner Nachzug — der Anlass aus dem Slice-Plan |

## Findings

### H-1 — Ein Link mit Code-Span im Linktext gilt nicht mehr als Link: der Nachzug des Anlasses ist wieder Drift

- **Kategorie:** HIGH
- **Quelle:** Prüffrage 2 und 13 · `DC-FA-VCS-001` Kriterium „Boundary (Pfad-Nachzug)" ·
  `ADR-0103` §Entscheidung 1 · Slice-Plan §1 *Ziel*
- **Pfad:** `internal/hexagon/core/rules/vcs.go` · „`if escapedAt(line, open) || overlapsSpan(code,
  m[0], m[1]) {`"; `spec/spezifikation.md` · „Kein Link
  ist, was einen Code-Span der Zeile berührt"
- **Befund:** `overlapsSpan` verwirft jeden Treffer, der mit einem Code-Span überlappt — auch den,
  dessen Code-Span ganz **im** Linktext liegt (``[`releasing.md`](…)``), was nach CommonMark ein
  Link ist. Vor `ab64a1e9` prüfte `inSpans` nur den Trefferanfang `[` und leerte diese Form. Jetzt
  meldet `vcs` für einen reinen Ziel-Nachzug `core-drift-vcs` (Proben 01, 10), und zwar für genau
  den Anlass aus Slice-Plan und `ADR-0103` §Kontext: Probe 11 zieht die vier `releasing.md`-Ziele
  von `ADR-0014` nach, Ergebnis `core-drift-vcs`, Exit 1. Im ADR-Bestand haben 687 der 1344
  Zeilen mit `](` diese Form (`grep -c '\[`[^]]*`\](' docs/plan/adr/*.md`). Kein Testfall trägt
  einen Link mit Code-Span im Linktext; die Fehlschlag-Listen (Spezifikation Schritt 4,
  `adr-check.md` Grenze 4, Code-Kommentar) nennen die Form nicht. `adr-check.md` sagt sogar
  „außerhalb von Code (… Code-Span)", was liest, als bleibe ein Link mit Code-Span im Text erhalten.
- **Failure-Szenario:** slice-268 verschiebt `docs/user/releasing.md` und zieht die Links in
  `ADR-0014`/`ADR-0067` nach; `pre-commit` (`make adr-check STAGED=1`) und die CI-Range melden
  `core-drift-vcs`, obwohl Lastenheft (Rang 1) zusagt: „Eine Änderung, die nur Ziele ändert, ist
  damit keine Drift". Zurück bleiben die zwei Wege, die `ADR-0103` §Kontext abschaffen sollte: tote
  Links per `ignore-refs` oder das Gate umgehen.
- **Warum HIGH, obwohl fail-safe:** Prüffrage 2 nennt den falschen Befund ausdrücklich; der Fehler
  trifft die Hälfte des Bestands und den einen Fall, für den der Schlüssel eingeführt wurde. Der
  Stand liegt auf `origin/main`.
- **Verifizierbar:** ja — Probe 11 im Wegwerf-Klon; ein Kern-Test
  ``"[`x`](a.md)"`` → ``"[`x`](x/a.md)"`` mit erwarteten 0 Befunden liefe rot.
- **Klasse:** `fix-ueberkorrigiert-ohne-bestandsprobe`

### M-1 — „Code bleibt Teil des Vergleichs" gilt weiter nicht für alle Code- und Nicht-Link-Formen

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 1 und 18 · `DC-FA-VCS-001` · `AGENTS.md` §5 Regel 13
- **Pfad:** `internal/hexagon/core/rules/vcs.go` · „`return strings.HasPrefix(line, "\t") ||
  strings.HasPrefix(line, "    ")`"; `spec/lastenheft.md` · „ihre benannten Grenzen
  sind ein Code-Span über mehrere Zeilen und eine Absatz-Folgezeile"; `harness/sensors/adr-check.md`
  · „außerhalb von Code (Fenced, eingerückt, Code-Span)"
- **Befund:** `indentedLine` erkennt Einzug nur als führenden Tab oder vier Leerzeichen; eine Zeile
  mit 1–3 Leerzeichen vor einem Tab ist nach CommonMark ebenfalls eingerückter Code und wird
  geleert (Probe 03). Ebenso werden Zeilen in einem HTML-Block (`<div>`, `<pre>`) geleert, in dem
  `[a](Nicht)` sichtbarer Text und kein Link ist (Proben 04, 05). Eine inhaltliche Umkehr passiert
  dort mit dem Schlüssel still. Das Lastenheft nennt als geleerte Grenzen nur „zwei", `adr-check.md`
  sagt „außerhalb von Code (… eingerückt …)", die Spezifikation „ohne Einzug von Tab oder vier
  Leerzeichen" (das stimmt mit dem Code, aber nicht mit der Code-Zusage des Lastenhefts).
- **Warum nicht HIGH:** im ADR-Bestand kein Treffer (keine Zeile `^ {1,3}\t`, kein HTML-Block-Anfang).
  Dieselbe Klasse wie R2 M-1 und R3 M-1, eine Form weiter.
- **Verifizierbar:** ja — Proben 03–05; Kern-Test `"Text\n\n \t[a](Nicht)"` → `(Immer)` liefe ohne
  Befund.
- **Klasse:** `muster-leert-mehr-als-gegenstand`

### I-1 — Die Berichtigungs-Zeile nennt nicht, was an der vorigen Zeile falsch war

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Pfad:** `docs/plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md` · „Berichtigt zur vorigen
  Zeile:"
- **Befund:** Die neue Zeile wiederholt die Spec-Verweise und gibt die Baseline-Lesart; welche
  Aussage der wiederhergestellten Zeile sie zurücknimmt („Alternative nachgetragen" — die ADR trägt
  im Körper keine nachgetragene Alternative), muss der Leser aus dem Vergleich der beiden Zeilen
  ableiten. Kein Failure-Szenario: beide Zeilen widersprechen dem Lastenheft nicht.
- **Verifizierbar:** nein.
- **Klasse:** `berichtigung-ohne-gegenstand`

### I-2 — Die R2-Zeile der Spezifikations-Historie wurde neu gefasst statt um eine R3-Zeile ergänzt

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Pfad:** `spec/spezifikation.md` · „nicht in einer eingerückten Zeile, nicht, wo er einen
  Code-Span berührt"
- **Befund:** Die gepushte Zeile aus `575e8d08` („nicht in einem Code-Span") wurde an die Spitze
  verschoben und umformuliert; das Lastenheft trägt für denselben Fix eine eigene Zeile 0.102.2.
  Die Reihenfolge ist jetzt richtig (R3 L-2 aufgelöst); dass die Spezifikation zwischen `575e8d08`
  und `ab64a1e9` „Code-Span" als Ausschluss nannte, zeigt nur noch `git log`. Keine Regel verlangt
  Append-only für diese Historie, deshalb ohne Kategorie darüber.
- **Verifizierbar:** nein.
- **Klasse:** `historie-zeile-umgeschrieben`

### I-3 — Steering-Loop-Signal: dieselbe Klasse in jeder Runde

- **Kategorie:** INFO
- **Quelle:** dieser Skill, §Kontext-Eskalation („dritte Wiederholung derselben Klasse in einer
  Sitzung")
- **Pfad:** `internal/hexagon/core/rules/vcs.go` · „`func blankLinkTargets(line string) string {`"
- **Befund:** Die Abgrenzung „was ist ein Link" am zeilenweisen Muster lieferte in R1, R2, R3 und
  R4 je ein MEDIUM oder HIGH in einer der zwei Richtungen (leert zu viel: R2 M-1, R3 M-1, R4 M-1;
  leert zu wenig: R4 H-1). Jeder Fix wurde gegen die gemeldeten Formen gebrochen, nicht gegen den
  unveränderten Bestand.
- **Verifizierbar:** nein.
- **Klasse:** `muster-leert-mehr-als-gegenstand`

## Status der R3-Findings

- **R3 M-1:** für die gemeldeten Formen eingelöst (Proben 06, 07; Tests vorhanden, Bewusstes
  Brechen in der Botschaft). Weitere Formen offen (M-1); die Lösung für den Code-Span bricht die
  Gegenrichtung (H-1).
- **R3 M-2:** aufgelöst. `git diff 273bebea -- docs/plan/adr/0103-…` zeigt eine einzige
  `+`-Zeile; die Zeile aus `273bebea` steht wortgleich wieder da. Gegen den zuletzt veröffentlichten
  Stand vor `cd7ebe04` ist der Unterschied ein Anhang (`AGENTS.md` §3.5).
- **R3 M-3:** aufgelöst. Die neue Zeile sagt, das Regelwerk nehme den Referenz-/Pfad-Nachzug „von
  der Überschreibung aus" — gedeckt durch `modul-04-adrs.md` „Zwei Fälle bleiben davon ausdrücklich
  getrennt" und „Beide Fälle ändern die Entscheidung nicht — nur ihre Adresse". Sie behauptet keine
  repo-eigene Form mehr und deckt sich mit Lastenheft und ADR-Kontext; ein `MR` ist nicht nötig.
  „empfiehlt, ihn … zu vermeiden" ist eine zulässige Lesart der Kennungs-Doktrin desselben Absatzes.
- **R3 M-4:** aufgelöst, soweit möglich. Die Botschaft von `575e8d08` ist gepusht und unveränderlich;
  `2115d0a2` benennt die Fehlzuschreibung und den tatsächlichen Commit `cd7ebe04`, trägt `ADR-0103`
  im Betreff und erscheint in `git log -- docs/plan/adr/0103-…` direkt neben `cd7ebe04`. Wer die
  ADR-Historie liest, findet die Richtigstellung.
- **R3 L-1:** eingelöst (Code-Kommentar, Spezifikation, Lastenheft, `adr-check.md`; Test
  „absatz-folgezeile in referenz-form").
- **R3 L-2:** eingelöst (siehe I-2).
- **R3 I-1:** eingelöst (Test „bild hinter escaptem ausrufezeichen").

## Negativbefunde (geprüft, ohne Befund)

- **`overlapsSpan`-Arithmetik:** halboffene Intervalle, `start < sp[1] && sp[0] < end` ist die
  korrekte Überlappung; Code-Span neben dem Link (Test „link neben inline-code") bleibt geleert.
  Kein Befund zur Arithmetik (zur Semantik: H-1).
- **`indentedLine` gegen Listen:** ein Listen-Folgeabsatz mit vier Leerzeichen wird nicht geleert
  — fail-safe, im Code-Kommentar, in der Spezifikation und in `adr-check.md` benannt. Im Bestand
  vier solche Zeilen (`ADR-0016`, `ADR-0040`, `ADR-0058`). Kein Befund.
- **Kommentare §3.7:** `vcsCore`, `blankLinkTargets`, `overlapsSpan`, `indentedLine` tragen Zusage
  bzw. `GRENZE`, keine Review-Historie, keine Befund-Nummern. Kein Befund zur Klasse (Wahrheit:
  H-1, M-1).
- **Testfall „absatz-folgezeile in referenz-form":** der Kommentar `// die benannte Grenze: …` ist
  Klasse Grenze. Kein Befund.
- **Spezifikation Schritt 4 gegen Code:** Einzug „Tab oder vier Leerzeichen" und Fenced-Code
  stimmen mit `indentedLine`/`proseLineSet`; Fußnote, Escapes, Referenz-Ziel mit Pfadzeichen
  unverändert. Kein Befund über H-1/M-1 hinaus.
- **Lastenheft 0.102.2:** Versionskopf und Historie-Zeile zueinander passend, neueste oben. Kein
  Befund (Inhalt: M-1).
- **Opt-in/byte-identisch:** der Schlüssel-Pfad ist an `ignoreLinkTargets` gebunden, der Zweig ohne
  ihn unverändert. Kein Befund.
- **`ADR-0005`/Hexagon, `DC-QA-03`, Suppressions:** nur `strings`/`regexp`, kein neuer Import, kein
  Netz, kein `//nolint`. Kein Befund.
- **Traceability der Commits:** beide Betreffe tragen `slice-267` und `DC-FA-VCS-001` bzw.
  `ADR-0103`. Kein Befund.
- **Botschaft `ab64a1e9` (§5 Regel 15):** nennt die gebrochenen Fälle und die Gates, behauptet
  keinen Erhalt des Nachzugs über die Proben hinaus. Kein Befund (dass die Bestandsprobe fehlt,
  trägt H-1).
- **Plan-Abgrenzung:** kein Umzug von `releasing.md`, kein Template-Feld-Nachzug. Kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 (H-1) |
| MEDIUM | 1 (M-1) |
| LOW | 0 |
| INFO | 3 (I-1, I-2, I-3) |

## Verdikt

**Blockiert durch H-1.** Die ADR-Hälfte ist erledigt: `2115d0a2` stellt die veröffentlichte Zeile
wortgleich wieder her, die Berichtigung ist ein Anhang, und die Baseline-Lesart deckt sich jetzt
mit Lastenheft und ADR-Kontext (R3 M-2 bis M-4 aufgelöst). Der Code-Fix für R3 M-1 verwirft
dagegen jeden Link, dessen Linktext einen Code-Span enthält — die Form der Hälfte des ADR-Bestands
und des Anlasses selbst. Ein reiner Nachzug der `releasing.md`-Links in `ADR-0014` ist am gebauten
Image `core-drift-vcs`, entgegen dem Abnahmekriterium „Boundary (Pfad-Nachzug)". Der Stand ist
gepusht; slice-268 sollte erst danach anlaufen. Daneben passieren eingerückter Code mit
Leerzeichen-Tab-Einzug und HTML-Blöcke weiter still (M-1). Da dieselbe Klasse in jeder Runde
wiederkehrt (I-3), ist eine Probe gegen den unveränderten Bestand vor dem nächsten Fix angezeigt.
