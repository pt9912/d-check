# Review-Report: slice-139 — closure-outcomes im fullbuild — 2026-08-25

**Review-Art:** Code-Review (Skript-/Doku-Diff gegen Kanon und Slice-Plan,
Modul 10 §Drei Review-Arten) · **Gegenstand:** Commit `612a619` (Diff
`HEAD~1..HEAD`) — der Feature-Commit von slice-139 (`feat(harness):
closure-outcomes im fullbuild — die Drei-Ausgänge-Regel bekommt ihre
Feedback-Hälfte (BEO-015, slice-139)`). 4 Dateien laut `git show --stat`
(`AGENTS.md` +2/-1, `Makefile` +9/-2, `harness/README.md` +1, neu
`tools/harness/closure-outcomes.sh` +68), 80 Einfügungen/3 Löschungen gesamt.

**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 · **Modell-ID:**
`claude-sonnet-5` · **Datum:** 2026-08-25

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-139-closure-ausgang-waechter.md`
  (§1–§9, insbesondere §3 NICHT, §4 DoD, §5 Risiken)
- `tools/harness/closure-outcomes.sh` vollständig (68 Zeilen)
- `Makefile` (neues Target `closure-outcomes`, `.PHONY`-Zeile, `fullbuild`-Kette)
- `AGENTS.md` §4 (Diff der `closure-outcomes`-/`fullbuild`-Zeilen)
- `harness/README.md` §Sensors (Diff, plus die vorbestehenden
  `fullbuild`-/Meta-Governance-Zeilen)
- `.d-check.closure.yml` vollständig (das bestehende Closure-Profil samt
  seiner `planning.closure.placeholder`-Begründung, ADR-0052)
- Baseline `modul-05-planning-harness.md` §Offene Risiken werden bei Closure
  aufgelöst (Volltext) und §Trigger je Lifecycle-Übergang
- Vendorte Vorlage
  `.harness/baseline/v5.11.0/templates/docs/plan/planning/slice.template.md`
  vollständig
- `docs/plan/planning/observations.md` → `BEO-007`, `BEO-010`, `BEO-011`,
  `BEO-012`, `BEO-015`
- Quellcode: `internal/hexagon/core/rules/planning.go`
  (`placeholderRE`, `checkClosurePlaceholder`, `placeholderRejected`),
  `internal/hexagon/core/rules/markdown.go` (`PreprocessMarkdown`,
  `proseLines`, `stripInlineCodeByLine`) — zur Prüfung, was die bestehende
  Platzhalter-Erkennung des Produkts mechanisch leistet, im Vergleich zum
  neuen Skript

**Nicht erhalten:** die DoD-Abhakung (Verifikations-Rolle, getrennter
Kontext).

**Vom Reviewer selbst gefahren** (nur Lesekommandos plus eine isolierte
Bash-Semantik-Probe in `/tmp`, kein `make`, keine Dateiänderung im Repo):
`git show`/`git log` auf den Commit; `find`/`grep` über
`docs/plan/planning/done/**` (Bestandszahl, manueller Nachvollzug der drei
Muster, Suche nach weiteren Klammer-Platzhaltern); Lesen von
`planning.go`/`markdown.go`, um die produkteigene Platzhalter-Erkennung
gegen das neue Skript zu vergleichen; eine Kopie von
`closure-outcomes.sh`s Kernschleife in
`/tmp/.../scratchpad/co-test/probe.sh` gegen eine unlesbar gemachte
(`chmod 000`) Testdatei gefahren, um das Fail-Closed-Verhalten bei
Leseversagen zu verifizieren (F-2) — die eigentliche Skript-Datei wurde
dabei **nicht** im Repo ausgeführt.

**Verdikt: blockierend** — zwei HIGH, ein MEDIUM, ein LOW.

---

## Findings

### F-1 — Die Platzhalter-Liste deckt nur eine von vielen Klammerformen der Vorlage; der Wächter startet grün, obwohl sein Kernfall unentdeckt bliebe

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §4 („Halluzinierte Gates sind die häufigste Form
  von Harness-Lüge") · Baseline `modul-05-planning-harness.md` §Offene
  Risiken werden bei Closure aufgelöst · slice-139 §5 Risiko 1 („Ein Wächter
  auf eine Zeichenkette ist so gut wie seine Liste")
- **pfad:** `tools/harness/closure-outcomes.sh:40-44` (`PATTERNS`) vs.
  `.harness/baseline/v5.11.0/templates/docs/plan/planning/slice.template.md:118`
  (`- <Risiko> — **Ausgang:** <eingetreten: CO-NNN / slice-NNN | entfallen:
  Grund | weiter offen: → BEO-NNN im Register>`) und dieselbe Vorlage an
  weiteren Stellen (`:1` `<Titel>`, `:14` `<welle-id>`, `:28`/`:35` `<Name>`,
  `:94-95` `<Bedingung>`, `:146` `<slice-NNN (<Titel>)`, `:175` `<Name>`)
- **befund:** `PATTERNS` enthält genau **eine** Vorlagen-Form, den
  wortwörtlichen Ellipsen-Platzhalter `<…>` (Byte-für-Byte identisch mit der
  Vorlage, geprüft). Die Vorlage selbst benutzt aber durchgehend **benannte**
  Klammer-Platzhalter (`<Name>`, `<Titel>`, `<Bedingung>`, `<welle-id>` …),
  und ausgerechnet ihr **§6-Ausgang-Feld** — das Feld, das die
  Drei-Ausgänge-Regel dieses Slice überhaupt erst erzwingen soll — trägt die
  Form `<eingetreten: … | entfallen: … | weiter offen: …>`, die keines der
  drei `PATTERNS` enthält. Ein Slice, der die Vorlage wörtlich kopiert und
  ein Risiko ohne Ausgang belässt (statt die repo-übliche Kurzform
  `*(bei Closure)*` zu verwenden), landet unentdeckt in `done/` — exakt der
  stille Grün-Pfad, den `BEO-015`/der Slice zu schließen vorgibt.
- **verifizierbar:** ja — eine Kopie eines `done/`-Slice mit einer
  Risikozeile `- <Risiko> — **Ausgang:** <eingetreten: CO-NNN / slice-NNN |
  entfallen: Grund | weiter offen: → BEO-NNN im Register>` (ohne Backticks)
  lässt `bash tools/harness/closure-outcomes.sh` mit Exit 0 durchlaufen.
- **klasse:** platzhalter-liste-deckt-nur-eine-vorlagenform

### F-2 — Fail-Closed gilt nur für die leere Prüfmenge, nicht für Leseversagen einer einzelnen Datei — der Wächter kann grün enden, ohne den Dateiinhalt gelesen zu haben

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §4 Harness-Lüge-Kriterium ·
  [`BEO-007`](../plan/planning/observations.md) („Ein Gate-Aufruf hinter
  einer Pipe meldet den Exit des letzten Pipe-Glieds … ein roter Lauf wird
  als grün behandelt") — dieselbe Fehlerklasse, hier innerhalb des
  Skripts selbst statt in seinem Aufruf
- **pfad:** `tools/harness/closure-outcomes.sh:60`
  (`` done < <(sed 's/`[^`]*`//g' "$f" | grep -n . || true) ``)
- **befund:** Schlägt `sed` beim Lesen von `$f` fehl (z. B. Rechte-Fehler,
  I/O-Fehler), liefert die Pipeline unter `pipefail` den Exit-Code von
  `grep` (1, da leere Eingabe → keine Treffer) statt den von `sed`; `|| true`
  neutralisiert selbst diesen. Die `while read`-Schleife erhält keine
  Zeilen, `findings` bleibt für diese Datei bei 0, und das Skript meldet am
  Ende „ok" (Exit 0) — obwohl der Dateiinhalt nie gelesen wurde. `set -euo
  pipefail` fängt das nicht, weil der fehlschlagende Befehl in einer
  `< <(...)`-Prozess-Substitution steckt, deren Exit-Status vom äußeren
  Skript nicht geprüft wird. **Experimentell reproduziert** (isolierte
  Kopie der Kernschleife in `/tmp`, nicht im Repo ausgeführt): eine Datei
  mit `chmod 000` und dem Inhalt `- Risiko — **Ausgang:** *(bei Closure)*`
  ergab `findings: 0` und Skript-Exit `0`, obwohl `sed` sichtbar
  „kann nicht gelesen werden: Keine Berechtigung" meldete.
- **verifizierbar:** ja — s. o., reproduziert außerhalb des Repos; dieselbe
  Mechanik greift für jede Form eines Leseversagens auf einer
  `done/slice-*.md`-Datei im echten Lauf.
- **klasse:** pipe-exit-durch-fallback-in-prozess-substitution-maskiert

### F-3 — Dieselbe Aufgabe (Vorlagen-Platzhalter erkennen, Inline-/Fence-Code ausklammern) ist im Produkt bereits robuster gelöst; das neue Skript baut eine engere Nachbildung ohne Fence-Behandlung

- **kategorie:** MEDIUM
- **quelle:** ADR-0052 / `DC-FA-PLAN-001` · Reviewer-Skill-Kategorie
  „Konsistenz-Lücke zwischen Modulen derselben Eingabe-Klasse"
- **pfad:** `tools/harness/closure-outcomes.sh:60` vs.
  `internal/hexagon/core/rules/planning.go:290-327`
  (`checkClosurePlaceholder`, allgemeine `<...>`-Erkennung mit Nachfiltern)
  und `internal/hexagon/core/rules/markdown.go:159-171`
  (`PreprocessMarkdown`: entfernt **Fences und** Inline-Code-Spans
  positionserhaltend)
- **befund:** Das Produkt besitzt bereits eine für exakt diesen Zweck
  gebaute, am eigenen Bestand kalibrierte Platzhalter-Erkennung
  (`planning.closure.placeholder`, aktiv in `.d-check.closure.yml`), die
  sowohl Fenced-Code als auch Inline-Code ausklammert und eine
  **generische** `<...>`-Form erkennt statt einer festen Liste. Das neue
  Skript reimplementiert denselben Grundgedanken (Platzhalter finden,
  Code-Kontext ausklammern) mit `` sed 's/`[^`]*`//g' `` — das behandelt nur
  paarige Inline-Backticks pro Zeile und kennt Fenced-Code-Blöcke
  (Fenced-Code-Bloecke) ueberhaupt nicht. Ein Beispiel-Fence, das eine der drei
  Zeichenketten ohne Inline-Backticks zeigt (z. B. eine künftige
  Dokumentations-Passage, die `- <Risiko> — **Ausgang:** *(bei Closure)*`
  als Zitat in einem Fenced-Block wiedergibt statt es einzeln zu
  escapen), würde als `closure-outcome-open` gemeldet, obwohl es kein
  offenes Risiko ist. Der Skript-Kopf benennt nur die
  Inline-Code-Ausnahme und ihre Begründung, nicht die fehlende
  Fence-Ausnahme — die Lücke ist nirgends dokumentiert.
- **verifizierbar:** ja — ein `done/`-Slice mit einem Fenced-Block, der die
  Zeile `- <Risiko> — **Ausgang:** *(bei Closure)*` ohne Inline-Backticks
  enthält, würde vom neuen Skript, nicht aber vom bestehenden
  `checkClosurePlaceholder`-Mechanismus fälschlich gemeldet.
- **klasse:** platzhalter-erkennung-doppelt-implementiert-ohne-fence-schutz

### F-4 — Zwei weitere Prosa-Fassungen der `fullbuild`-Kette in derselben Datei blieben unaktualisiert, obwohl die Commit-Botschaft Vollständigkeit der Doku-Nachzucht behauptet

- **kategorie:** LOW
- **quelle:** [`BEO-010`](../plan/planning/observations.md) (dieselbe
  Klasse: „Modulliste … hat Spiegel außerhalb der Config — und keiner
  davon ist gate-gedeckt", hier innerhalb derselben Zieldatei)
- **pfad:** `harness/README.md:93` (`make fullbuild`-Zeile: „gates +
  image-test + bench + completeness-check + verify-closure-notes" —
  `closure-outcomes` fehlt) und `harness/README.md:111`
  (Meta-/Governance-Gates-Zusammenfassung: nennt `completeness-check`,
  `verify-closure-notes` als Closure-Bindepunkte, `closure-outcomes`
  fehlt)
- **befund:** Der Commit fügt in `harness/README.md` die neue
  `closure-outcomes`-Sensors-Zeile (Zeile 90) hinzu und zieht die
  `fullbuild`-Kettenprosa in `AGENTS.md`/`Makefile` nach — lässt aber
  zwei weitere Prosa-Nennungen derselben Kette **in derselben Datei**
  unverändert. Die Commit-Botschaft benennt „Drei Doku-Flächen
  nachgezogen (Makefile-Hilfe, AGENTS.md §4 samt der fullbuild-Kette,
  Sensors-Tabelle)" — trifft für die genannten drei Flächen zu, deckt
  aber nicht, dass dieselbe Kette an zwei weiteren Stellen derselben
  dritten Fläche weiterhin den alten Stand zeigt. `gate-consistency`
  prüft nur Target-Namen/-Existenz, nicht Ketten-Prosa, fängt dies also
  nicht.
- **verifizierbar:** ja — `grep -n "verify-closure-notes)" harness/README.md`
  zeigt die veraltete Zeile 93 neben der aktuellen Zeile 90.
- **klasse:** doku-flaeche-teilweise-nachgezogen

## Negativbefunde

- geprüft, ohne Befund: **Scope `done/slice-*.md`, `-maxdepth 1`** — passt
  zum Wortlaut von `modul-05` („Ein Risiko **aus dem Slice-Plan**"); `done/`
  hat keine Unterverzeichnisse (`find … -maxdepth 1 -type d` leer); die
  ca. 40 `welle-*.md`/`welle-*-results.md`-Dateien in `done/` liegen
  außerhalb der Regel — dieselbe Artefakt-Unterscheidung, die
  `.d-check.closure.yml` für die analoge Frage bei `verify-closure-notes`
  bereits explizit dokumentiert (Wellen-Ergebnisnotiz *ist* eine
  Closure-Notiz, *enthält* keine) — kein zu enger Scope gegenüber der
  Regel, die tatsächlich aufgestellt wird.
- geprüft, ohne Befund: **Einordnung `fullbuild` statt `gates`** — die
  Begründung (Regel gilt dem Übergang nach `done/`, derselbe
  Closure-Bindepunkt wie `completeness-check`/`verify-closure-notes`) ist
  intern konsistent und widerspricht der genannten Präzedenz nicht; die
  Alternative „Aufnahme in `.d-check.closure.yml`/Modul `structure`" wird
  in §3 des Slice-Plans **explizit** erwogen und mit nachvollziehbarem
  Grund (fehlende Modul-Fähigkeit `forbid-match`, benötigt ADR + Zensus)
  vertagt — keine stillschweigend übergangene Option (s. aber F-3 zur
  Wiederverwendung der bereits vorhandenen Bausteine).
- geprüft, ohne Befund: **Bestandszahl 137 und „null offene Platzhalter"**
  — eigener Nachvollzug (`find … | wc -l` = 137; manueller Re-Scan mit
  denselben drei Mustern über alle 137 Dateien = 0 Treffer) bestätigt die
  Commit-Angabe exakt.
- geprüft, ohne Befund: **`modul-05`-Zitat** — der in `AGENTS.md`/im
  Skript-Kopf/im Slice-Plan wiederholte Satz „Ein Slice geht nicht nach
  `done/`, während ein Risiko ohne Ausgang dasteht" ist wortgleich mit der
  Baseline-Quelle, im richtigen Geltungsbereich (Slice-Closure, nicht
  Wellen-Closure) zitiert.
- geprüft, ohne Befund: **Makefile-Mechanik** — `.PHONY` um
  `closure-outcomes` ergänzt, Target trägt die geforderte
  „fullbuild-nicht-gates"-Begründung im Kommentar, `fullbuild`-Prerequisite-
  Liste korrekt erweitert.
- geprüft, ohne Befund: **`AGENTS.md` §4-Zeile für `closure-outcomes`** —
  behauptet nicht mehr, als das Skript leistet (nennt nur die
  Inline-Code-Ausnahme, nicht Fence-Schutz; nennt „Ob ein Ausgang
  inhaltlich trägt, bleibt Urteil" korrekt als Nicht-Zusage) — kein
  Über- oder Unter-Claim gegenüber dem tatsächlichen Mechanismus.
- geprüft, ohne Befund: **Odd-Backtick-Verhalten** (ungerade
  Backtick-Zahl je Zeile lässt Rest sichtbar stehen) — entspricht der im
  Skript-Kopf selbst benannten Fail-Loud-Entscheidung; bevorzugt
  Falsch-Positiv über Falsch-Negativ, kein unbenannter Seiteneffekt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** platzhalter-liste-deckt-nur-eine-vorlagenform
· pipe-exit-durch-fallback-in-prozess-substitution-maskiert ·
platzhalter-erkennung-doppelt-implementiert-ohne-fence-schutz ·
doku-flaeche-teilweise-nachgezogen

## Verdikt

**Merge-blockierend:** ja — zwei HIGH. F-1 ist ein stiller Grün-Pfad in
genau dem Feld, das die Drei-Ausgänge-Regel erzwingen soll (die
Vorlagen-eigene `<eingetreten: … | entfallen: … | weiter offen: …>`-Form
bleibt unentdeckt), F-2 ist dieselbe Fehlerklasse wie das bereits
registrierte `BEO-007` („Pipe meldet den Exit des letzten Glieds, ein
roter Zustand erscheint grün"), jetzt innerhalb des neuen Wächter-Skripts
selbst statt in seinem Aufruf, und experimentell reproduziert. F-3 ist ein
Konsistenz-/Redundanz-Mangel gegenüber der bereits vorhandenen,
robusteren Platzhalter-Erkennung des Moduls `planning` (Fence-Schutz).
F-4 ist reine Doku-Drift ohne Gate-Deckung.

**Übergabe:** Findings gehen an den Implementer. F-1/F-2 sind
Kandidaten für einen neuen `BEO-<NNN>`-Eintrag bzw. — für F-2 — für eine
weitere Instanz der bereits registrierten `BEO-007`-Klasse; die
Einordnung obliegt dem Maintainer bei der Slice-Closure (§7), nicht
diesem Report. Dieser Report ist ein Lauf-Beleg und ersetzt keine
Verifikation (DoD-/Spec-Konformität prüft der Verifier separat).
