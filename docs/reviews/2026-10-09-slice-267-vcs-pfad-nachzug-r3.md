# Review R3: slice-267, `vcs` lässt einen reinen Pfad-Nachzug durch

- **Review-Art:** Code. Geprüft wird der Fix-Diff `cd7ebe04..575e8d08` (`fix(vcs)` `575e8d08`) gegen
  den Slice-Plan `slice-267`, gegen die Findings aus R1 und R2
  (`2026-10-09-slice-267-vcs-pfad-nachzug-r1.md`, `-r2.md`), gegen `ADR-0103`, `DC-FA-VCS-001`,
  `MR-032` und die Hard Rules `AGENTS.md` §3.5/§3.7 sowie §5 Regel 13 und 15. Mitgeprüft ist die
  Änderung an `ADR-0103` §Geschichte, die R2 M-2 einlösen soll: Sie liegt nicht im Fix-Commit,
  sondern in `cd7ebe04` (dem Commit des R2-Reports). Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-267` · `cd7ebe04..575e8d08`, dazu die ADR-Hälfte von `cd7ebe04`.
- **Skill:** `reviewer.md` @ 1.19.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** `DC-FA-VCS-001` im Lastenheft 0.102.1 (Absatz „Pfad-Nachzug (opt-in)",
  Historie 0.102.0/0.102.1); Spezifikation `DC-FA-VCS-001.a` Schritt 4, §2-Zeile
  `vcs.ignore-link-targets`, §8-Historie; `ADR-0103` samt Geschichte; Baseline `v6.17.0` ·
  `regelwerk/modul-04-adrs.md` §Nachzug ist keine Überschreibung; `harness/conventions.md`
  §Purpose; `MR-032`. Im Code: `vcsCore`, `inlineLinkTargetRE`, `refDefTargetRE`,
  `blankLinkTargets`, `escapedAt`, `inSpans` in `internal/hexagon/core/rules/vcs.go`;
  `proseLineSet`, `proseLines`, `forEachInlineCodeSpan` in `markdown.go`; Kommentar von
  `VCSConfig.IgnoreLinkTargets`; `vcs_pfad_nachzug_test.go`; `harness/sensors/adr-check.md`
  Grenzen 2 und 4. Vorherige Findings am Modul: R1 und R2 zu `slice-267`.
- **Proben** (Image `dcr267r3` per `make build IMAGE=dcr267r3` vom Stand `575e8d08`; je Probe ein
  Wegwerf-git-Repo im Scratchpad mit einer `Accepted`-ADR, zwei Commits, Lauf
  `--enable vcs --disable links --disable anchors --range HEAD~1..HEAD`, Schlüssel an und aus;
  danach Repos und Image entfernt, Arbeitsbaum unverändert). Ohne Schlüssel meldet jede Probe
  `core-drift-vcs`, Exit 1.

  | # | BASE → HEAD | mit Schlüssel | Markdown-Lesart |
  |---|---|---|---|
  | 01 | ``Regel [a `b](Nicht)` gilt.`` → ``…(Immer)`…`` | **still** | kein Link: der Code-Span beginnt vor `]` und hat Vorrang |
  | 02 | Indented Code (4 Leerzeichen) `fns[k](true)` → `fns[k](false)` | **still** | Code-Block, kein Link |
  | 03 | Code-Span über zwei Zeilen, Zeile 2 `fns[k](true)` → `(false)` | still | benannte Grenze |
  | 04 | Absatz-Folgezeile `[Status]: Abgelehnt.` → `Angenommen.` | **still** | gerenderter Fließtext, keine Referenz-Definition |
  | 05 | ``\`[A](a.md)\` `` → Ziel nachgezogen | Befund | echter Link — fail-safe |
  | 06 | `![A](a.png)` → Ziel nachgezogen | still | Nachzug |
  | 07 | `\![A](a.png)` → Ziel nachgezogen | still | Nachzug (`!` escaped, Link bleibt) |
  | 08 | ``[A](a.md) `c[k](1)` [B](b.md)`` → Ziele nachgezogen **und** `(1)` → `(2)` | Befund | wie zugesagt |
  | 09 | wie 08, nur Ziele nachgezogen | still | Nachzug |
  | 10/11 | ungeschlossener Backtick vor `[A](…)` | still | echter Link; 11 ist Wort-Ziel, gedeckt durch „das Gate sieht nur die Form" |
  | 12 | drei Links, ein Linktext geändert | Befund | wie zugesagt |
  | 13 | `[a\]b](a.md)` → nachgezogen | still | Nachzug |
  | 14 | `\\[R](Nicht)` → `(Immer)` | still | echter Link, Wort-Ziel |
  | 15 | `<!-- [A](Nicht) -->` → `(Immer)` | still | unsichtbar; nicht gemeldet |
  | 16 | `~~~`-Fence `fns[k](true)` → `(false)` | Befund | wie zugesagt |
  | 17 | ``[A](`x`)`` → ``(`y`)`` | still | echter Link, Wort-Ziel |
  | 18 | `<https://x/[a](Nicht)>` → `(Immer)` | still | Autolink; Grenzfall, nicht gemeldet |

## Findings

### M-1 — „Code bleibt unverändert" gilt nicht für Indented Code und nicht für einen Code-Span, der im Linktext beginnt

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 1 und 18 · `DC-FA-VCS-001` · `AGENTS.md` §5 Regel 13
- **Pfad:** `spec/lastenheft.md` · „eine Fußnote, Code, ein escapter Link oder eine Zeile ohne
  vollständigen Link"; `internal/hexagon/core/rules/vcs.go` · „`if escapedAt(line, open) || inSpans(code, m[0]) {`"
  und „`prose := proseLineSet(content)`"; `spec/spezifikation.md` · „in jeder Prosa-Zeile (außerhalb
  von Fenced-Code)"
- **Befund:** Die Prosa-Menge schließt nur Fenced-Code aus, eingerückter Code gilt als Prosa (Probe
  02); und `inSpans` prüft nur den **Anfang** des Treffers, so dass ein Code-Span, der im Linktext
  beginnt und hinter der Zielklammer schließt, den Treffer nicht schützt (Probe 01 — nach CommonMark
  kein Link). Beide Male passiert eine inhaltliche Umkehr im Code mit dem Schlüssel still, während das
  Lastenheft „Code … bleibt unverändert Teil des Vergleichs" zusagt. Die Grenzen-Listen
  (Spezifikation Schritt 4, `adr-check.md` Grenze 4, Code-Kommentar) nennen nur den mehrzeiligen
  Code-Span; das Lastenheft nennt gar keine Grenze in dieser Richtung, nur die fail-safe-Richtung.
  Kein Test hält eine der beiden Formen.
- **Warum nicht HIGH:** Beide Formen kommen im ADR-Bestand heute nicht vor (eingerückte Zeilen mit
  `](` außerhalb von Listen: keine; Code-Span, der im Linktext öffnet: keiner). Das Modul ist aber
  ein Produkt-Vertrag (`DC-FA-VCS-001`), nicht nur dieses Repos Gate; derselbe Klassen-Pfad wie R2
  M-1, eine Form weiter.
- **Verifizierbar:** ja — Proben 01, 02 am Image; ein Kern-Test mit
  `"Text.\n\n    fns[k](true)"` liefe ohne Befund.
- **Klasse:** `muster-leert-mehr-als-gegenstand`

### M-2 — Eine veröffentlichte Geschichte-Zeile einer `Accepted`-ADR wurde ersetzt statt ergänzt

- **Kategorie:** MEDIUM
- **Quelle:** `AGENTS.md` §3.5 · `harness/sensors/adr-check.md` Grenze 2
- **Pfad:** `docs/plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md` · „Präzisiert: Die Grenzen der
  Erkennung legt"
- **Befund:** `273bebea` hängte die Zeile „Grenzen präzisiert in … trägt für künftige ADRs, nicht für
  den Bestand, dessen Verweise schon Adressen sind" an; `cd7ebe04` ersetzte sie durch eine andere
  Zeile (`git show cd7ebe04 -- docs/plan/adr/`: `-`/`+` derselben Tabellenzeile). §3.5 erlaubt
  **Anhänge** an `## Geschichte`, keine Neufassung; das Gate sieht es nicht, weil der Abschnitt
  vollständig ausgenommen ist (Grenze 2: „die Unterscheidung *Anhang gegen verkleidete
  Kern-Änderung* ist ein Urteil"). Die Zeile war **nicht** unveröffentlicht: `origin/main` stand ab
  10:28:19 auf `273bebea` (Reflog `refs/remotes/origin/main`), die Neufassung kam 10:38. Folge: der
  Pfad-Anker von R2 M-2 („trägt für künftige ADRs, nicht für den Bestand …") löst in der ADR nicht
  mehr auf, und die Geschichte zeigt nicht mehr, was die ADR zehn Minuten lang auf dem Hauptzweig
  sagte.
- **Verifizierbar:** nein — kein Gate vergleicht Geschichte-Zeilen; `git log -p` zeigt es.
- **Klasse:** `geschichte-zeile-ueberschrieben`

### M-3 — Die neue Geschichte-Zeile sagt über die Baseline das Gegenteil von Lastenheft und Slice-Plan, und die so benannte Abweichung hat keinen `MR`

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 9 und 20 · `harness/conventions.md` §Purpose (Adaptionen ggü. Baseline) ·
  Baseline `modul-04-adrs.md` §Nachzug ist keine Überschreibung
- **Pfad:** `docs/plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md` · „der Pfad-Nachzug ist eine
  dritte Form, die dieses Repo hier deklariert"; `spec/lastenheft.md` · „(die Baseline nimmt ihn von
  der Immutabilität einer angenommenen Entscheidung aus)"
- **Befund:** Die Geschichte-Zeile erklärt, die Baseline nenne als Formen des Nachzugs nur Kennung
  und Verfall-Vermerk, der In-Place-Pfad-Nachzug sei eine **repo-eigene dritte Form**. Lastenheft
  0.102.1 (Anforderungstext und Historie 0.102.0), der Slice-Plan („auch wenn die Baseline den
  Nachzug erlaubt") und der Kontext der ADR sagen weiter, die Baseline nehme ihn aus. Eine der
  beiden Lesarten ist falsch: Stimmt die Geschichte-Zeile, ist der Nachzug eine Abweichung von der
  Baseline, die `harness/conventions.md` nicht führt (kein `MR`, kein Auflösungs-Trigger dort; der
  Re-Evaluierungs-Trigger der ADR „die Baseline ändert die Ausnahme" setzt die Ausnahme voraus);
  stimmt das Lastenheft, behauptet die immutable Zeile eine Abweichung, die es nicht gibt.
  Failure-Szenario: der Freshness-Audit gegen die Baseline sucht Abweichungen im Adaptions-Block und
  findet diese nicht; ein Leser des Lastenhefts (Rang 1) hält das Verhalten für baseline-konform,
  während die ADR es als Eigenbau führt.
- **Verifizierbar:** nein — Urteil über die Lesart einer Kanon-Stelle; der Widerspruch selbst ist
  per `grep -n "Baseline nimmt" spec/lastenheft.md` und der Geschichte-Zeile sichtbar.
- **Klasse:** `quellen-widersprechen-sich-ueber-baseline`

### M-4 — Der Fix-Commit behauptet die ADR-Änderung, die in einem anderen Commit liegt

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 8 · `AGENTS.md` §5 Regel 15 · Baseline `modul-08-agentenrollen.md`
  §Rollen-Regeln
- **Pfad:** Commit `575e8d08` · „R2 M-2: die eigene Geschichte-Zeile von ADR-0103 neu gefasst";
  Commit `cd7ebe04` · „docs(reviews): slice-267 R2 — link-foermiger Text ohne Link, ADR-Geschichte"
- **Befund:** `575e8d08` berührt `docs/plan/adr/` nicht (`git show --stat`: sieben Dateien, keine
  ADR). Die Neufassung liegt in `cd7ebe04`, dem Commit des R2-Reports, dessen Botschaft nur die
  Kategorie-Zahlen nennt. Damit steht die Implementer-Antwort auf ein Finding im selben Commit wie
  das Finding, und der Report wurde mit einem Pfad-Anker committet, der in diesem Commit schon nicht
  mehr auflöste. Wer M-2 am Fix-Commit prüft, findet nichts; wer die ADR-Historie liest, findet die
  Änderung unter einem Review-Commit ohne `ADR-0103` in der Botschaft.
- **Verifizierbar:** ja — `git show --stat 575e8d08`, `git show cd7ebe04 -- docs/plan/adr/`.
- **Klasse:** `botschaft-behauptet-fremden-diff`

### L-1 — „wie Markdown sie liest" stimmt für eine Absatz-Folgezeile nicht

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §5 Regel 13
- **Pfad:** `internal/hexagon/core/rules/vcs.go` · „gilt als
  Referenz-Definition, wie Markdown sie liest"; `spec/spezifikation.md` · „(auch wenn
  sie nicht gerendert wird)"; `harness/sensors/adr-check.md` · „auch wenn Markdown sie nicht rendert"
- **Befund:** Eine Referenz-Definition kann keinen Absatz unterbrechen; als Folgezeile ist
  `[Status]: Abgelehnt.` sichtbarer Fließtext und wird trotzdem geleert (Probe 04). Die Grenzen
  beschreiben nur den Fall „nicht gerendert" und behaupten Markdown-Lesart, wo das Muster sie nicht
  hat. Kein Bestandsfall (keine Zeile im ADR-Baum trifft `refDefTargetRE`), deshalb LOW.
- **Verifizierbar:** ja — Probe 04.
- **Klasse:** `grenzen-beispiel-enger-als-regel`

### L-2 — Neue Historie-Zeile der Spezifikation steht unter der älteren

- **Kategorie:** LOW
- **Quelle:** `.d-check.yml` `structure`-Regel für `spec/spezifikation.md` (`order: desc`, neue
  Zeile oben)
- **Pfad:** `spec/spezifikation.md` · „Nachzug nach Review an §[`DC-FA-VCS-001.a`]"
- **Befund:** Die Nachzug-Zeile ist neuer als die Zeile „`vcs.ignore-link-targets` leert beim
  Bilden des Core jedes Link-Ziel", steht aber darunter; die übrigen gleichtägigen Zeilen
  (`DC-FA-RVW-001.a`) stehen neueste zuerst. Das Gate sieht es nicht, weil beide dasselbe Datum
  tragen. Das Lastenheft hat 0.102.1 richtig oben.
- **Verifizierbar:** nein — gleiche Daten.
- **Klasse:** `historie-reihenfolge-bei-gleichem-datum`

### I-1 — Die `!`-Verschiebung in `blankLinkTargets` hat keinen tötenden Test

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Pfad:** `internal/hexagon/core/rules/vcs.go` · „`if line[open] == '!' {`"
- **Befund:** Entfiele `open++`, bliebe nur `\![A](a.png)` ungeleert (fail-safe); kein Testfall
  trägt ein escaptes `!`. Probe 07 zeigt das gewollte Verhalten am Image.
- **Verifizierbar:** nein.
- **Klasse:** `zweig-ohne-test`

## Einlösung der R2-Findings

| R2 | Stand | Beleg |
|---|---|---|
| M-1 | weitgehend eingelöst | Proben R2-1/R2-2/R1-esc entsprechen jetzt den Testfällen „escapte klammern", „funktionsaufruf in inline-code", „escapte oeffnende klammer" (Befund); Restpfad Indented Code und Code-Span im Linktext: R3 M-1 |
| M-2 | inhaltlich adressiert, Form verletzt | neue Zeile gibt für Kennung und Vermerk je einen Grund und korrigiert die Kontext-Lesart; aber per Ersetzung (R3 M-2), im Review-Commit (R3 M-4), und die neue Lesart widerspricht dem Lastenheft (R3 M-3) |
| L-1 | eingelöst | Lastenheft 0.102.1 mit Historie, Spezifikation mit Historie-Zeile (Reihenfolge: R3 L-2) |
| L-2 | eingelöst | Kommentar zeigt auf Schritt 4 statt die Zusage zu wiederholen |
| L-3 | eingelöst | Beispiel jetzt `[Label]: wort.` mit „ein Pfadzeichen" (Folgezeilen-Fall: R3 L-1) |
| I-1 | bleibt INFO | bewusste Entscheidung, unverändert |
| I-2 | eingelöst | `rawVCS` und Test-Struct ausgerichtet |

## Negativbefunde (geprüft, ohne Befund)

- **Index-Arithmetik der Submatches:** `m[2]:m[3]` ist Gruppe 1 (Linktext samt `!`), `m[4]:m[5]`
  Gruppe 2 (Titel und schließende Klammer); `last = m[1]` übernimmt den Rest. Mehrere Treffer pro
  Zeile, angrenzende Links und Code-Span zwischen zwei Links (Proben 08, 09, 12) verhalten sich wie
  zugesagt. Kein Befund.
- **`!`-Bilder:** `open++` prüft den Escape vor `[`, nicht vor `!` — `\![A](…)` bleibt ein Link
  und wird richtig geleert (Probe 07). Kein Befund (Test: I-1).
- **Ungeschlossene Backticks:** `forEachInlineCodeSpan` behandelt den öffnenden Lauf als literal,
  der folgende echte Link wird geleert (Proben 10/11) — CommonMark-konform. Kein Befund.
- **Escapter Backtick:** ``\`[A](a.md)\` `` gilt dem Scanner als Code, ist aber ein Link — fail-safe
  Drift (Probe 05). Kein Befund.
- **Escape-Zählung:** `escapedAt` zählt Backslash-Läufe paritätisch; `\\[R](…)` bleibt Link
  (Test, Probe 14). Kein Befund.
- **Referenz-Definition nach der Inline-Leerung:** wird auf den bereits bereinigten Text
  angewandt, setzt `^` mit 0–3 Leerzeichen voraus und kann deshalb weder in Code-Spans noch in
  Indented Code greifen. Kein Befund.
- **Fenced-Code:** ``` ``` ``` und `~~~` aus `proseLineSet` ausgenommen (Test, Probe 16). Kein Befund.
- **Opt-in/byte-identisch:** ohne Schlüssel meldet jede der 19 Proben. Kein Befund.
- **Kommentare §3.7:** `escapedAt`, `inSpans`, `blankLinkTargets`, `vcsCore`, `VCSConfig` tragen
  Zusage bzw. Grenze, keine Review-Historie, keine Befund-Nummern. Kein Befund (Wahrheitsfrage:
  M-1, L-1).
- **Spezifikation §2-Zeile:** verweist für die Grenzen auf Schritt 4, wiederholt sie nicht. Kein
  Befund.
- **`ADR-0005`/Hexagon, `DC-QA-03`, Suppressions:** nur `regexp`/`strings` und der Kern-eigene
  Markdown-Scanner, kein Netz, kein `//nolint`. Kein Befund.
- **Plan-Abgrenzung:** kein Umzug von `releasing.md`, kein Template-Feld-Nachzug, kein Eingriff in
  den Linktext-Vergleich. Kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 (M-1, M-2, M-3, M-4) |
| LOW | 2 (L-1, L-2) |
| INFO | 1 (I-1) |

## Verdikt

**Nicht blockiert durch HIGH; MEDIUM vor Merge zu klären.** Der Code-Fix trägt R2 M-1 für die
gemeldeten Formen; zwei weitere Code-Formen (Indented Code, Code-Span ab dem Linktext) passieren
weiter still, und das Lastenheft sagt für Code das Gegenteil (M-1). Die übrigen drei MEDIUM betreffen
die ADR-Hälfte: eine veröffentlichte Geschichte-Zeile wurde ersetzt statt ergänzt (M-2), die neue
Zeile widerspricht dem Lastenheft über die Baseline (M-3), und die Änderung liegt im Review-Commit,
nicht im Fix-Commit, der sie behauptet (M-4). M-2 und M-3 berühren eine `Accepted`-ADR; ihre
Auflösung ist eine Architect-Frage (Anhang mit Korrektur, Folge-ADR oder Lastenheft-Nachzug samt
`MR`), kein Implementer-Edit.
