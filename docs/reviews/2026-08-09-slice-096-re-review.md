# Review-Report: slice-096 — bestätigende Re-Review der Heilung — 2026-08-09

**Review-Art:** **Bestätigende Re-Review** (Design/Plan, eng geführt) — geprüft
werden genau zwei Fragen: (1) ist jeder der 32 Vor-Befunde geheilt, (2) hat die
Überarbeitung Neues kaputt gemacht. Keine Vollprüfung; alles, was keiner der
beiden Fragen zugeordnet werden kann, bleibt außerhalb.

**Gegenstand:** Heilungs-Commit `5b6c7f4` (Diff-Range `cb5c3ee..HEAD`) zu
[slice-096](../plan/planning/in-progress/slice-096-structure-modul-analyse.md);
berührt [`spec/lastenheft.md`](../../spec/lastenheft.md),
[`spec/spezifikation.md`](../../spec/spezifikation.md),
[ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md), den
Slice-Plan, das
[Wellendokument](../plan/planning/welle-69-structure-schnitt.md) und die Slices
unter `docs/plan/planning/open/`.

**Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) @ 1.3.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-08-09

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne diese Liste
ist der Lauf nicht reproduzierbar):

- Die beiden Vor-Reviews als Checkliste:
  [Vertrags-Review](2026-08-09-slice-096-vertrag-review.md) (17 Befunde) und
  [Umsetzbarkeits-Review](2026-08-09-slice-096-umsetzbarkeit-review.md)
  (15 Befunde)
- [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
  [`DC-FA-CONF-002`](../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope),
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl),
  [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben),
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)
- §[`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure),
  §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning),
  §2-Schema und §7-Historie derselben Datei
- [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) (Proposed)
  samt [ADR-Index](../plan/adr/README.md); Bezugs-ADRs
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md),
  [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md),
  [ADR-0042](../plan/adr/0042-markdown-lexik-folgt-commonmark.md)
- [slice-099](../plan/planning/open/slice-099-structure-modul.md) (neu gefasst),
  [slice-101](../plan/planning/open/slice-101-fence-unbalanciert.md) (neu),
  [slice-094](../plan/planning/open/slice-094-closure-zaehl-paritaet.md),
  [slice-097](../plan/planning/open/slice-097-closure-glob-entkopplung.md),
  [slice-098](../plan/planning/open/slice-098-closure-note-placeholder.md)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules), Dogfooding-Konfiguration
  [`.d-check.yml`](../../.d-check.yml), Befund-Dedup in
  `internal/hexagon/core/model/finding.go`
- **Nur lesend beigezogen:** die drei vermessenen Prüfskripte des Schwester-Repos
  `a-check` samt seinem Lastenheft und seinem Closure-Bestand

**Belegläufe:** das ausgelieferte Image gegen das Repo
(`docker run --rm --network none -v "$PWD":/repo:ro d-check:latest` ⇒
`343 Datei(en) geprüft, 0 Befund(e)`) und `make trace` (⇒ `48 Anforderung(en),
0 Waise(n)`,
[`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
ohne Waise auf den neuen Folge-Slice abgebildet). Marken- und
Zahlen-Gegenproben per `grep` gegen beide Repos; keine Fixture im Repo angelegt.

---

## Frage 1 — Heilungs-Tabelle

Urteile: **geheilt** · **teilweise** · **nicht geheilt** · **bewusst nicht**.
Geprüft wurde inhaltlich (trägt die neue Fassung den Fall?), nicht per
Stichwort-Suche.

### Vertrags-Review — 17 Befunde

| Nr. | Kat. | Urteil | Beleg |
|---|---|---|---|
| F-1 Kardinalität nicht ausdrückbar | HIGH | **geheilt** | `spec/lastenheft.md:1985–1997` führt `sections` mit `one`/`each` ein; `spec/spezifikation.md:1752–1758` ist die Kardinalität ein eigener Algorithmus-Schritt; `spec/spezifikation.md:2217` trägt den Schema-Schlüssel; [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) Entscheidung 7 begründet ihn. **Leerlauf unter `each` scharf geprüft:** alle drei Stellen sagen übereinstimmend „0 Treffer ⇒ `section-missing` (`line` = 1)", Mehrfachtreffer sind unter `each` ausdrücklich kein Befund. Die Motiv-Regel ist damit schreibbar (`files` auf das Lastenheft, `section-pattern` auf die Anforderungs-Überschrift, `sections: each`, `require-all` mit den vier Marken) |
| F-2 `require-any` deckt Messzeile 8 nicht | MEDIUM | **geheilt** | `require-any` ist ersatzlos entfallen (`spec/lastenheft.md:2027–2032`); an seine Stelle tritt `require-pattern` als Spiegelbild von `forbid-pattern`. Die Marke `Lerneintrag` trifft die im Adopter-Bestand gelebte Form `**Lerneintrag — Form: …**` unter der neuen Semantik (Inhalt beginnt mit der Marke, läuft nicht-alphanumerisch weiter), die Formen-Alternation liegt jetzt im Muster statt in der Marke |
| F-3 Summenzeile geht nicht auf | MEDIUM | **geheilt** | `docs/plan/planning/in-progress/slice-096-structure-modul-analyse.md:76` nennt „5 nicht gedeckt … (Summe 11)"; nachgezählt trifft das (nicht gedeckt: Zeilen 2, 4, 7, 8, 10). Beide Folgestellen nachgezogen: `spec/lastenheft.md:2449` und `docs/plan/adr/0049-structure-modul-schnitt-und-preset.md:39` |
| F-4 Messzeile 4 mit Fremd-Beleg | MEDIUM | **nicht geheilt** | Zeile 4 der Messtabelle steht unverändert (`…-analyse.md:67`): Einstufung „nicht gedeckt", Beleg weiterhin „vier Platzhalter-Sätze passieren alle drei Codes" — also der Beleg aus der Messung eines anderen Antrags, nicht aus dem hier vermessenen Skript (dessen fünf literale Phrasen `planning.closure.boilerplate` heute schon entgegennimmt). Kein Grund für das Stehenlassen genannt |
| F-5 `min-sentences`-Default vs. Exit-2-Grenze | MEDIUM | **geheilt** | `spec/lastenheft.md:2049–2050` und `spec/spezifikation.md:1730–1734` sagen beide „**explizit** gesetztes `min-sentences` < 1 oder `max-tasks` < 0"; die Spezifikation ergänzt den Satz, dass ein abwesender Zahlen-Schlüssel „Bedingung aus" heißt. `spec/spezifikation.md:2219–2220` führt beide Defaults als „abwesend (aus)" statt `0` |
| F-6 Ventil hebelt Nullmengen-Guard aus | MEDIUM | **geheilt** | Beide Straten sagen jetzt dasselbe: `spec/lastenheft.md:2051–2053` („`exempt-paths` hebelt den Leerlauf-Befund nicht aus") und `spec/spezifikation.md:1742–1745` („auch dann, wenn erst `exempt-paths` die Menge geleert hat"); das Akzeptanzkriterium `spec/lastenheft.md:2068` nennt den Fall ausdrücklich |
| F-7 Preset-Kopplung ohne Akzeptanzkriterium | MEDIUM | **teilweise** | Ein Kriterium existiert jetzt — aber nur in **einer** der beiden Anforderungen (`spec/lastenheft.md:2067`). Die Kriterienliste von [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (`spec/lastenheft.md:1928–1942`) trägt keines. Beide Straten behaupten dagegen „ein Akzeptanzkriterium **beider** Anforderungen" (`spec/lastenheft.md:2449`, `spec/spezifikation.md:1801`) ⇒ siehe F-N5 |
| F-8 „erstes Modul ohne Referenz-Invariante" | MEDIUM | **teilweise** | Korrigiert an drei Stellen (`spec/lastenheft.md:17–24`, `spec/lastenheft.md:1960–1963`, `docs/plan/adr/0049-…md:26–30`). **Unverändert** an vier weiteren: `docs/plan/adr/0049-…md:174` (§Konsequenzen), `docs/plan/adr/README.md:64`, `docs/plan/planning/welle-69-structure-schnitt.md:24`, `…-analyse.md:36` und `…-analyse.md:112` ⇒ siehe F-N2/F-N3 |
| F-9 Sammel-Code auf instabilem Feld | MEDIUM | **geheilt** | Sechs Bedingungs-Codes statt eines (`spec/lastenheft.md:2003–2010`, `spec/spezifikation.md:1767–1774`), begründet in [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) Entscheidung 8 samt neuer Alternativen-Zeile. Die Begründung ist am Dedup-Tupel in `internal/hexagon/core/model/finding.go:50` belegt und trifft zu |
| F-10 Abwärts-Referenz als Prosa-Zeiger | MEDIUM | **nicht geheilt** | Die gerügte Stelle wurde entschärft („Grund-Codes sind stabil zugesagt", `spec/lastenheft.md:2037–2038`), aber im selben Diff ist ein neuer Prosa-Abwärtszeiger derselben Klasse entstanden: `spec/lastenheft.md:2015` („die laut §4 der Spezifikation gerade **nicht** stabil zugesagt ist"). Wieder ohne Link, ohne Inline-Code-Pfad, ohne Token — von `matrix`/`links`/`codepaths` unsichtbar, und wieder **normativ** (er trägt die Begründung der Code-Aufteilung) |
| F-11 „der erste Treffer" an drei Stellen | MEDIUM | **teilweise** | Nachgezogen ist genau **eine** der drei: die §2-Schema-Zeile (`spec/spezifikation.md:2211` sagt jetzt „**genau ein** Treffer … mehrere ⇒ `closure-note-ambiguous`"). Unverändert: `spec/lastenheft.md:1877` („in ihr der **erste** Abschnitt") und der Eingangssatz von Schritt C3 (`spec/spezifikation.md:1678`, „Je Kandidat wird die **erste** Zeile gesucht") — letzterer widerspricht drei Zeilen später seiner eigenen Mehrdeutigkeits-Regel |
| F-12 Ventil-Granularität Messzeile 11 | MEDIUM | **teilweise** | Der Vertrag benennt den Fall jetzt: namentliche Ausnahmen innerhalb einer Datei sind ausdrückliches Nicht-Ziel (`spec/lastenheft.md:2077–2079`, [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) Entscheidung 9). Damit ist die Einstufung derselben Messzeile als „**Ventil** — über die bestehende Ausnahme-Mechanik ausdrückbar" (`…-analyse.md:74`) jetzt **positiv widerlegt** statt nur unbelegt, und steht unverändert da ⇒ siehe F-N4 |
| F-13 Zeitangabe Contract-Churn | LOW | **teilweise** | Die ADR ist korrigiert (`docs/plan/adr/0049-…md:41`, „**am selben Tag**"). Der Slice-Plan trägt die widerlegte Angabe weiter: `…-analyse.md:248` („ein Bruch nach wenigen Tagen") |
| F-14 Glossar bleibt bei 19 Modulen | LOW | **nicht geheilt** | `spec/lastenheft.md:2435` führt weiterhin 19 Namen mit `sources` als letztem; die neue Historie-Zeile nennt weiterhin nur `DC-FA-CLI-002` als nachgezogene Stelle (`spec/lastenheft.md:2449`). `DC-FA-CLI-002` selbst trägt `structure` (20 Namen) |
| F-15 Leerlauf-Befund ohne `line`-Wert | LOW | **geheilt** | `line` = 1 ist an beiden Stellen gepinnt (`spec/lastenheft.md:1973`, `spec/spezifikation.md:1742`); zusätzlich pinnt Schritt 7 (`spec/spezifikation.md:1787`) `line` generell auf die Überschriftszeile bzw. 1 |
| F-16 Zählung übergeht die Guard-Klasse | LOW | **nicht geheilt** | Die Messtabelle hat unverändert elf Zeilen; die drei Grundgesamtheits-Guards der vermessenen Skripte bekommen weiterhin keine Zeile, obwohl der Vertrag sie als Leerlauf-Befund neu aufnimmt. Kein Grund genannt |
| F-17 Kennzahl im Dogfooding-Kommentar | INFO | **bewusst nicht** | `.d-check.yml:184` nennt weiterhin „42 Anforderungen"; Ist ist 48 (`make trace`). Der Befund war ausdrücklich als „nicht Gegenstand dieses Diffs" gemeldet — das Stehenlassen ist konsistent damit und tragfähig |

### Umsetzbarkeits-Review — 15 Befunde

| Nr. | Kat. | Urteil | Beleg |
|---|---|---|---|
| F-1 Unbalancierter Fence macht still grün | HIGH | **teilweise** | Der Defekt ist als eigener Slice registriert und mit dem Reproduktionsfall dokumentiert ([slice-101](../plan/planning/open/slice-101-fence-unbalanciert.md)); [slice-099](../plan/planning/open/slice-099-structure-modul.md) §5 nennt die Vererbung über die geteilte Mechanik. **Offen bleibt der Vertrag selbst:** `spec/spezifikation.md:1761` schreibt die neue Anforderung weiterhin auf die defekte `FenceToggle`-Lexik fest, kein Akzeptanzkriterium deckt den unbalancierten Fence, und die Reihenfolge ist nur weich formuliert („vorzugsweise nach", „sinnvoll vor") und in keiner Wellen-Abhängigkeit verankert |
| F-2 Sammel-Code kollidiert mit Befund-Dedup | HIGH | **teilweise** | Der Hauptfall ist geheilt: sechs verschiedene Grund-Codes plus `line` = Zeile der Abschnitts-Überschrift (`spec/spezifikation.md:1787`) machen zwei verletzte Bedingungen und zwei verletzende Abschnitte unterscheidbar; das Akzeptanzkriterium `spec/lastenheft.md:2062` fordert es. **Der zweite, ausdrücklich genannte Fall ist ungelöst:** zwei Regeln mit **demselben** `files`-Glob, die beide `section-missing` melden, haben identisches (Datei, Zeile 1, Regel, Ziel, Grund) und fallen unter `SortFindings` (`internal/hexagon/core/model/finding.go:50`) zu einem Befund zusammen — die zweite fehlende Abschnitts-Zusage verschwindet |
| F-3 Zeilenanker ignoriert Listen-Marker | HIGH | **geheilt** | `spec/spezifikation.md:1776–1781` und `spec/lastenheft.md:2017–2025` verankern die Marke „nach optionalem **Listen-Marker** und Whitespace" mit einem hervorgehobenen Textlauf, dessen Inhalt mit der Marke beginnt und dort endet oder nicht-alphanumerisch weitergeht. **Gegenprobe gefahren:** die drei zugesagten Formen decken die gemessenen Realformen ab — im Lastenheft dieses Repos 108 Listen-Item-Marken, 44 bare, 64 qualifizierte; im Lastenheft des Adopters 33/16/32. Neue Falsch-Positive: nur die gewollte Klasse (eine Marke `Boundary` trifft auch `- **Boundary-Fall …**`); Fließtext-Vorkommen bleiben ausgeschlossen. Zwei im Vor-Befund genannte Unterfälle bleiben unabgedeckt (Tabellenzelle bewusst, eingerückte Unterliste unbenannt ⇒ F-N12), beide im gemessenen Bestand mit **null** Vorkommen |
| F-4 Post-Pass ohne definierte Eingabemenge | MEDIUM | **teilweise** | Zwei der drei offenen Punkte sind entschieden: `spec/lastenheft.md:1968–1977` und `spec/spezifikation.md:1735–1741` legen fest, dass die Globs Wurzel-relativ über den gesamten Baum laufen, **unabhängig** von `scan.roots`/`scan.ignore`, dass die `SKIP_DIRS` weiter gelten und dass `structure` deshalb kein `<modul>.scope` kennt. **Unentschieden bleibt der dritte:** ob Nicht-Markdown-Dateien Kandidaten sind. Ein Glob, der nur Nicht-Markdown trifft, ist unter der einen Auslegung Leerlauf (`section-missing`), unter der anderen eine Reihe zeilenweise gelesener Fremdformate |
| F-5 Schema-Zeile nicht mitgeführt | MEDIUM | **teilweise** | Identisch zu Vertrags-F-11: die Schema-Zeile ist nachgezogen, `spec/lastenheft.md:1877` nicht. Die Anforderungs-Beschreibung und ihre eigene Mehrdeutigkeits-Regel stehen damit weiterhin gegensätzlich im selben Absatz-Block |
| F-6 Vergleichsgegenstand von `section` | MEDIUM | **geheilt** | Entschieden und an beiden Technik-Stellen identisch: „getrimmte Überschriften-Zeile **einschließlich der `#`-Folge**, exakt" (`spec/spezifikation.md:1749–1751`, `spec/spezifikation.md:2215`). Der Happy-Path-Test ist damit ohne Raten schreibbar; die Überschriftsebene ist bewusst mitgenagelt |
| F-7 Zwei Anker-Regime für RE2-Schlüssel | MEDIUM | **geheilt** | `spec/spezifikation.md:1782–1784`: beide Muster matchen gegen den gesamten bereinigten Abschnitts-Text, `^`/`$` binden an Text- statt Zeilengrenzen, wer Zeilen meint, schreibt `(?m)` — RE2-ausdrückbar |
| F-8 Schema größer als der Kern-Slice | MEDIUM | **geheilt** | Die beiden Folge-Slices sind zu einem zusammengeführt; [slice-099](../plan/planning/open/slice-099-structure-modul.md) §2 nennt genau die Release-Grenze als Grund und die strikte Dekodierung als Mechanik |
| F-9 Nullwert vs. „nicht gesetzt" | MEDIUM | **geheilt** | `spec/spezifikation.md:1732–1734` formuliert die Unterscheidung ausdrücklich als Teil der Zusage; `max-tasks: 0` ist damit die schärfste Setzung und vom Default unterscheidbar |
| F-10 Universelle Config-Zusage mit stiller Ausnahme | MEDIUM | **teilweise** | Die Ausnahme ist nicht mehr still (`spec/lastenheft.md:1974–1977`). Aber [`DC-FA-CONF-002`](../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope) selbst ist unverändert und sagt weiterhin „**Jedes** Regelmodul akzeptiert … `<modul>.scope`" (`spec/lastenheft.md:2352`). Zwei Anforderungen desselben Stratums widersprechen sich jetzt ausdrücklich statt stillschweigend |
| F-11 CLI-Mit-Modifikation ohne Slice | MEDIUM | **teilweise** | [slice-099](../plan/planning/open/slice-099-structure-modul.md) §3 Punkt 3 und die DoD tragen sie jetzt. Der Vertrag zieht nicht mit: [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) nennt weiterhin „elf Targets" und erklärt weitere ausdrücklich zum Out-of-Scope ⇒ siehe F-N8 |
| F-12 Task-Item-Lexik per Beispiel | LOW | **geheilt** | `spec/spezifikation.md:1775–1777` definiert das Task-Item als Regel: optionaler Whitespace, Listen-Marker (`-`, `*`, `+` oder Ziffernfolge mit Punkt), Whitespace, `[ ]` bzw. `[x]`/`[X]` |
| F-13 Abbruch-Skopus Regel vs. Datei | LOW | **geheilt** | `spec/spezifikation.md:1755–1756`: „Der Abbruch gilt der **Datei innerhalb dieser Regel** — andere Regeln und andere Dateien bleiben unberührt" |
| F-14 Config-Rand in drei Fassungen | LOW | **teilweise** | Die beiden Listen sind jetzt wortgleich und tragen die fehlende Glob-Validierung (`spec/lastenheft.md:2046–2051`, `spec/spezifikation.md:1728–1731`). Dabei ist der zuvor vorhandene Exit-2-Guard für einen **leeren `exempt-paths`-Eintrag** ersatzlos herausgefallen ⇒ siehe F-N11 |
| F-15 Nullmengen-Guard endet an der Modulgrenze | INFO | **nicht geheilt** | `spec/lastenheft.md:2042–2043` hält unverändert fest, dass ein aktiviertes Modul ohne Regeln inert und grün ist; die im Befund benannte Asymmetrie zur Regel-Ebene wird nirgends aufgelöst oder als bewusst erklärt |

**Bilanz Frage 1:** 14 geheilt · 12 teilweise · 5 nicht geheilt · 1 bewusst nicht.

---

## Frage 2 — neue Findings

### F-N1 — Die §7-Historie der Spezifikation beschreibt die verworfene erste Fassung, im selben Dokument, dessen Algorithmus-Sektion das Gegenteil sagt

- `kategorie`: MEDIUM
- `quelle`: §[`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure) /
  [`AGENTS.md`](../../AGENTS.md) §2 (Source Precedence)
- `pfad`: `spec/spezifikation.md:2352`
- `befund`: Die Historie-Zeile der Spezifikation ist vom Heilungs-Commit nicht
  angefasst worden und trägt weiterhin die Schlüsselliste mit `require-any`, den
  Sammel-Code `section-constraint` („je verletzter Bedingung einer, mit Ist/Soll
  in der Meldung"), „**genau ein** Abschnitt" und „Marken sind **zeilenverankert
  und hervorgehoben** (`**M:**`) … `require-all` fordert jede, `require-any`
  mindestens eine". `sections`, `require-pattern` und die sechs
  Bedingungs-Codes fehlen. Die Historie-Zeile des Lastenhefts wurde vollständig
  neu geschrieben, die der Spezifikation nicht — das Änderungsprotokoll des
  Dokuments beschreibt damit die Fassung, die dasselbe Dokument 600 Zeilen
  weiter oben widerruft.
- `verifizierbar`: nein maschinell (kein Gate vergleicht §7 gegen den Körper) —
  ja durch Textvergleich innerhalb einer Datei.
- `klasse`: „Änderungsprotokoll nicht mit der Änderung mitgeführt"

### F-N2 — Der ADR-Index trägt zwei zurückgezogene Entscheidungen als Zusammenfassung

- `kategorie`: MEDIUM
- `quelle`: [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  Entscheidung 5 + §Kontext
- `pfad`: `docs/plan/adr/README.md:64`
- `befund`: Die Index-Zeile ist unverändert und fasst die ADR mit zwei
  Aussagen zusammen, die die ADR selbst zurückgenommen hat: „das erste Modul
  ohne Referenz-Invariante" (in §Kontext als Einzelfall-Formulierung korrigiert,
  weil `spans`/`hostpaths` längst die Form prüfen) und „zwei Marken-Formen statt
  einer (`require-all`/`require-any`, zeilenverankert statt Teilstring)"
  (`require-any` ist laut Entscheidung 5 ersatzlos entfallen, die Verankerung neu
  gefasst). Der Index ist die Einstiegsfläche für Leser, die die ADR **nicht**
  öffnen; er referiert einen Stand, den es nicht mehr gibt.
- `verifizierbar`: nein maschinell — ja durch Vergleich der Index-Zeile mit
  Entscheidung 5 und §Kontext derselben ADR.
- `klasse`: „Index-Zusammenfassung nicht mit der Revision mitgeführt"

### F-N3 — Innerhalb der ADR widerspricht §Konsequenzen dem neu gefassten §Kontext

- `kategorie`: MEDIUM
- `quelle`: [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
- `pfad`: `docs/plan/adr/0049-structure-modul-schnitt-und-preset.md:173–176`
- `befund`: §Kontext sagt nach der Korrektur, Aussagen über die Form eines
  Dokuments habe es als Einzelfälle längst gegeben (`spans`, `hostpaths`);
  §Konsequenzen sagt unverändert, `structure` sei „das erste Modul ohne
  Referenz-Invariante". Dieselbe ADR behauptet damit beides. Zwei kleinere
  Nachzüge derselben Klasse in derselben Datei: die Alternativen-Zeile
  „Mehrdeutigkeit als `section-constraint` mitführen" (Zeile 162) verwirft eine
  Bauform anhand eines Grund-Codes, den der Vertrag nicht mehr kennt, und die
  Fitness-Function-Zeile „Mehrdeutigkeit schlägt Messung" (Zeile 193) ist
  unqualifiziert formuliert, obwohl Entscheidung 7 sie ausdrücklich zur
  Eigenschaft des Modus `one` erklärt. Für die Entscheidungen 7 und 8 ist keine
  Fitness Function nachgetragen.
- `verifizierbar`: nein maschinell — ja durch Vergleich der Abschnitte
  derselben Datei.
- `klasse`: „Nachträglich erweiterte ADR nur im vorderen Teil nachgezogen"

### F-N4 — Der Abnahme-Punkt, der die Marken-Entscheidung festhält, widerspricht der Entscheidung

- `kategorie`: MEDIUM
- `quelle`: [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  Entscheidung 5 /
  [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
- `pfad`: `docs/plan/planning/in-progress/slice-096-structure-modul-analyse.md:178–187`,
  `…-analyse.md:71`, `…-analyse.md:74`
- `befund`: Abnahme-Punkt 4 des Slice-Plans ist das Entscheidungs-Protokoll
  dieses Slice und steht unverändert: „der Vertrag braucht **zwei** Formen
  („alle von" und „mindestens eine von")" und „**Die Marke ist zeilenverankert
  und ausgezeichnet** … das Skript prüft `**Name:**` **am Zeilenanfang**".
  Beides ist vom Vertrag und von der ADR widerrufen. Ebenso stehen geblieben:
  Messzeile 8 mit dem Beleg „Alternation über benannte Marken" (die neue Fassung
  sagt ausdrücklich, das sei keine Marken-Frage) und Messzeile 11 mit der
  Einstufung „Ventil", die Entscheidung 9 zum Nicht-Ziel erklärt hat. Wer den
  Slice-Plan als Entscheidungsgrundlage liest — die Rolle, die er im Prozess hat
  —, bekommt drei widerrufene Entscheidungen.
- `verifizierbar`: nein maschinell — ja durch Vergleich von Abnahme-Punkt 4 und
  Messtabelle gegen Entscheidung 5/9 der ADR.
- `klasse`: „Entscheidungs-Protokoll nicht mit der Revision mitgeführt"

### F-N5 — „ein Akzeptanzkriterium beider Anforderungen" ist in einer der beiden nicht vorhanden

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
  [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  Entscheidung 3
- `pfad`: `spec/lastenheft.md:2449`, `spec/spezifikation.md:1801`,
  `spec/lastenheft.md:1926–1942`
- `befund`: Beide Straten sagen zu, die Preset-Kopplung werde durch „ein
  Akzeptanzkriterium **beider** Anforderungen" zusammengehalten. Vorhanden ist
  es nur in
  [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  (`spec/lastenheft.md:2067`); die Kriterienliste von
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  ist unverändert und enthält keinen Satz über zwei Oberflächen. Die Aussage,
  mit der der Vor-Befund als geheilt gilt, ist damit zur Hälfte unbelegt: ändert
  jemand die Closure-Fähigkeit, fällt in ihrer eigenen Anforderung kein
  Kriterium.
- `verifizierbar`: ja — Auszählen der Given/When/Then-Sätze beider
  Anforderungen.
- `klasse`: „Zugesagte Invariante ohne Akzeptanzkriterium"

### F-N6 — Die Grund-Code-Zählung des Umsetzungs-Slice untererfasst um zwei Codes

- `kategorie`: MEDIUM
- `quelle`: [slice-099](../plan/planning/open/slice-099-structure-modul.md) §3/§4
- `pfad`: `docs/plan/planning/open/slice-099-structure-modul.md:41–43`,
  `docs/plan/planning/open/slice-099-structure-modul.md:55`
- `befund`: Der Vertrag führt **acht** `section-*`-Codes (`section-missing`,
  `section-ambiguous`, `section-empty`, `section-thin`, `section-oversized`,
  `section-forbidden`, `section-pattern-missing`, `section-marker-missing`) plus
  den additiven `closure-note-ambiguous` — neun neue Codes. Der Slice nennt „die
  sechs `section-*` plus `closure-note-ambiguous`" und die DoD „Sieben neue
  Grund-Codes im Lockstep mit `AllReasons()` und Spezifikation §4". Gezählt sind
  offenbar nur die sechs Bedingungs-Codes; die beiden Abschnitts-Codes fallen
  aus der Lockstep-Liste. Der Lockstep ist genau die Stelle, an der ein Code
  ohne §4-Zeile und ohne Doctor-Klartext ausgeliefert würde.
- `verifizierbar`: ja — Auszählen der Codes in `spec/lastenheft.md:2003–2010`
  plus `spec/lastenheft.md:1973`/`:1988` gegen die DoD-Zeile.
- `klasse`: „Aggregat widerspricht der Einzelaufstellung"

### F-N7 — Zwei abgeschlossene Review-Reports wurden nachträglich editiert; ein Link trägt jetzt ein Label, das nicht zu seinem Ziel passt

- `kategorie`: MEDIUM
- `quelle`: [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md)
  §Ablage („Nie überschreiben — Folgeläufe bekommen eine neue Datei")
- `pfad`: `docs/reviews/2026-08-09-slice-096-umsetzbarkeit-review.md:26`,
  `…-umsetzbarkeit-review.md:110`, `…-umsetzbarkeit-review.md:207`,
  `…-umsetzbarkeit-review.md:265–266`,
  `docs/reviews/2026-08-09-slice-096-vertrag-review.md:24`
- `befund`: Der Heilungs-Commit hat beide Lauf-Belege mitgeändert, um die Links
  auf den entfallenen zweiten Folge-Slice auflösbar zu halten. Ergebnis sind
  sechs Links, deren sichtbares Label eine Slice-Nummer nennt und deren Ziel
  eine **andere** Slice-Datei ist; ein Leser, der dem Zeiger folgt, landet
  nachweisbar woanders, als das Label sagt. Zugleich ist derselbe Report jetzt
  in sich widersprüchlich: das `pfad`-Feld seines Befundes F-11 nennt weiterhin
  zwei Dateien, die es nicht mehr gibt. Kein Gate fängt beides —
  [`.d-check.yml`](../../.d-check.yml) nimmt `docs/reviews/**` von `codepaths`
  aus, und die umgebogenen Links lösen auf.
- `verifizierbar`: nein am Gate (bewusstes Ventil) — ja durch `git show` des
  Heilungs-Commits gegen die beiden Report-Dateien.
- `klasse`: „Lauf-Beleg nachträglich verändert"

### F-N8 — Der Umsetzungs-Slice sagt eine `--print-mk`-Erweiterung zu, die der Vertrag ausdrücklich ausschließt

- `kategorie`: LOW
- `quelle`: [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
- `pfad`: `docs/plan/planning/open/slice-099-structure-modul.md:44–46`,
  `spec/lastenheft.md:433–493`
- `befund`: Der Slice zieht die „`--print-mk`-Target-Liste" als
  CLI-Mit-Modifikation ein.
  [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  nennt aber „**elf** `##`-annotierte Targets", zählt sie im Happy-Path-Kriterium
  einzeln auf und führt unter Out-of-Scope „weitere Targets jenseits der
  gelisteten elf". Eine Mit-Modifikation dieser Anforderung — die jede
  Modul-Einführung mit `--print-mk`-Target bisher mitgebracht hat — fehlt im CR.
  Die Umsetzung kann die DoD entweder erfüllen und gegen den Vertrag laufen oder
  den Vertrag halten und die DoD offen lassen. Genau diesen Zustand hat die
  Historie schon einmal als „Selbstwiderspruch" verbucht (Eintrag 0.37.1).
- `verifizierbar`: ja — `--print-mk` nach dem Slice gegen die Aufzählung in
  [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben).
- `klasse`: „DoD-Zusage ohne Vertrags-Deckung"

### F-N9 — Verwaiste Wellen-Zuordnung nach der Slice-Zusammenführung

- `kategorie`: LOW
- `quelle`: [welle-69](../plan/planning/welle-69-structure-schnitt.md) §4 /
  [slice-096](../plan/planning/in-progress/slice-096-structure-modul-analyse.md) §4
- `pfad`: `docs/plan/planning/in-progress/slice-096-structure-modul-analyse.md:242`
- `befund`: Die Zeile „**Wellen-Zuordnung:** 099 und 100 gehören zu dieser
  Welle" nennt eine Slice-Nummer, deren Datei im selben Commit gelöscht wurde;
  die Wellen-Tabelle führt folgerichtig nur noch einen Folge-Slice. Umgekehrt
  taucht der im selben Commit neu angelegte Fence-Slice in keiner
  Wellen-Zuordnung auf — er trägt „ohne Welle" (eine im Repo etablierte Form),
  aber die DoD desselben Slice-Plans führt ihn als Ergebnis dieser Welle. Die
  Zuordnungs-Aussage ist damit an beiden Enden nicht mehr wahr.
- `verifizierbar`: nein am Gate (die Nummer steht als Prosa, nicht als Link) —
  ja durch Verzeichnis-Listing von `docs/plan/planning/open/`.
- `klasse`: „Aufzählung nach Datei-Entfall nicht nachgezogen"

### F-N10 — Die Marken-Messzahlen sind einem Repo entnommen, aber beiden zugeschrieben; die zweite Messzahl ist heute nicht reproduzierbar

- `kategorie`: LOW
- `quelle`: [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  Entscheidung 5 /
  [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
- `pfad`: `docs/plan/adr/0049-structure-modul-schnitt-und-preset.md:99–101`,
  `spec/lastenheft.md:2021–2022`, `spec/lastenheft.md:2029–2030`
- `befund`: „Gemessen an den **beiden** Repos, die den Antrag tragen … zu 108 als
  Listen-Item, zu 44 bare". Nachgezählt sind 108 und 44 exakt die Zahlen des
  **eigenen** Lastenhefts; das Lastenheft des Adopters trägt zusätzlich 33
  Listen-Item- und 16 bare Marken, die Vereinigung wäre 141/60. Die zweite Zahl
  („37 von 61 Lerneintrag-Formen tragen die Form-Angabe im Inneren des
  Textlaufs") ist heute nicht reproduzierbar: im Closure-Bestand des Adopters
  finden sich 68 hervorgehobene Lerneintrag-Läufe, davon 38 mit der Form-Angabe
  im Inneren. Keine der Zahlen nennt ihre Messmethode oder ihren Messstand,
  obwohl sie die tragende Begründung der neu gefassten Entscheidung 5 sind.
- `verifizierbar`: ja — Auszählung der Marken-Formen in beiden Lastenheften und
  der Lerneintrag-Läufe im Closure-Bestand des Adopters.
- `klasse`: „Messzahl ohne Methoden- und Standangabe"

### F-N11 — Der Exit-2-Guard für einen leeren `exempt-paths`-Eintrag ist bei der Neufassung entfallen

- `kategorie`: LOW
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  §fail-closed
- `pfad`: `spec/spezifikation.md:1728–1731`, `spec/spezifikation.md:2224`
- `befund`: Die abgelöste Fassung führte „leerer Eintrag in
  `require-all`/`require-any`/`exempt-paths`" als Exit-2-Fall. Die neue Liste
  nennt nur noch `require-all`; die neu aufgenommene Glob-Validierung fängt den
  Fall nicht, weil ein leeres Muster unter der Glob-Semantik dieses Repos
  (segmentweiser Vergleich, `internal/hexagon/core/rules/paths.go:98`) ein
  **gültiges** Muster ohne Treffer ist. Ein leerer `exempt-paths`-Eintrag ist
  damit ein still wirkungsloses Ventil — genau die Fehlerform, die der Vor-Befund
  benannt hatte und die die Schwester-Schlüssel (`tracked.exempt-targets`,
  `planning.closure.boilerplate`) ausdrücklich abfangen.
- `verifizierbar`: ja — Konfiguration mit einem leeren `exempt-paths`-Eintrag:
  Exit 2 oder stilles Weiterlaufen entscheidet.
- `klasse`: „Guard bei der Neufassung verloren"

### F-N12 — Die Marken-Verankerung nennt keinen führenden Whitespace, die Task-Item-Lexik zwei Zeilen darüber schon

- `kategorie`: INFO
- `quelle`: §[`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 6
- `pfad`: `spec/spezifikation.md:1775–1781`
- `befund`: Das Task-Item ist als „Zeile, die — nach optionalem Whitespace — mit
  einem Listen-Marker … beginnt" definiert; die Marke als „eine Zeile beginnt —
  nach optionalem Listen-Marker und Whitespace — mit einem hervorgehobenen
  Textlauf". Wörtlich gelesen kennt die Marken-Regel keine Einrückung **vor** dem
  Listen-Marker, die Task-Item-Regel schon. Eine Marke in einer eingerückten
  Unterliste — im Vor-Befund ausdrücklich als Lücke genannt — fiele damit
  weiterhin durch. Im gemessenen Bestand beider Repos gibt es **null** solche
  Zeilen, der Fall ist heute folgenlos; dokumentationswürdig ist die
  unbegründete Asymmetrie zweier Lexiken im selben Absatz.
- `verifizierbar`: ja — Fixture mit einer Marke in einer eingerückten
  Unterliste.
- `klasse`: „Zwei Lexiken derselben Sektion mit ungleicher Whitespace-Regel"

---

## Negativbefunde

- geprüft, ohne Befund: **RE2-Ausdrückbarkeit des neuen Vertrags.** Verlangt wird
  an keiner Stelle Lookaround, Rückwärtsreferenz oder rekursives Matching.
  `section-pattern`, `forbid-pattern` und `require-pattern` sind reine, vom
  Nutzer gelieferte Treffer-Prädikate; die neu aufgenommene Mehrzeilen-Klärung
  benennt mit `(?m)` genau ein RE2-Konstrukt; die Marken- und Task-Item-Semantik
  ist als Lexik-Regel formuliert, nicht als Regex. Auch der Motiv-Fall
  („eine von drei Formen") ist als Alternation ohne Lookahead schreibbar.
- geprüft, ohne Befund: **`sections: each` gegen den Leerlauf-Guard.** Lastenheft,
  Algorithmus-Schritt 4 und Schema-Zeile sagen übereinstimmend „0 Treffer ⇒
  `section-missing` (`line` = 1)"; `each` schaltet den Nullmengen-Guard also
  nicht ab, und Mehrfachtreffer sind dort ausdrücklich kein Befund. Die
  Kandidaten-Nullmenge (Schritt 2) und die Abschnitts-Nullmenge (Schritt 4)
  teilen sich denselben Code ohne Kollision, weil `file` verschieden ist (Glob
  vs. Datei).
- geprüft, ohne Befund: **Benennung der sieben Konfigurations- und der acht
  Befund-Bezeichner über alle Fundstellen.** `sections`, `one`, `each`,
  `require-pattern`, `require-all`, `section-empty`, `section-thin`,
  `section-oversized`, `section-forbidden`, `section-pattern-missing`,
  `section-marker-missing` sind in Lastenheft-Beschreibung, Lastenheft-Tabelle,
  Akzeptanzkriterien, Algorithmus-Tabelle, Schema-Tabelle und Historie-Zeile des
  Lastenhefts **zeichengleich** geschrieben. Die einzigen Abweichungen liegen in
  der Historie-Zeile der Spezifikation (F-N1) und in der Zählung des Slice
  (F-N6), nicht in der Schreibweise.
- geprüft, ohne Befund: **`require-any` als Restbestand.** Außerhalb der beiden
  in F-N1/F-N2 genannten Stellen kommt der entfallene Schlüssel nirgends mehr
  vor; in der ADR steht er ausschließlich als benannter Entfall (Entscheidung 5)
  und als zwei Verwerfungs-Zeilen der Alternativen-Tabelle.
- geprüft, ohne Befund: **§4-Grund-Code-Tabelle der Spezifikation.** Die neuen
  Codes fehlen dort bewusst; beide Historie-Zeilen führen den etablierten
  Lockstep („Grund-Codes (§4) folgen mit der Implementierung"), wie zuvor bei
  `targets`. Kein Widerspruch.
- geprüft, ohne Befund: **Auszählung der Messtabelle.** Elf Zeilen, Verteilung
  2/2/5/1/1, Summe 11 — die korrigierte Summenzeile trifft, und beide
  Folgestellen (Lastenheft-Historie, ADR-Kontext) tragen jetzt dieselbe Zahl.
- geprüft, ohne Befund: **Verwaiste Datei-Verweise nach dem Slice-Umbau.** Kein
  Link im Repo zeigt auf eine der beiden gelöschten Slice-Dateien; die
  Wellen-Tabelle, die DoD des Analyse-Slice und die Trigger-Abschnitte beider
  neuer Slices verweisen konsistent auf die zusammengeführte Datei. Die
  Ausnahmen sind die Prosa-Nummer in F-N9 und die Inline-Code-Pfade in F-N7.
- geprüft, ohne Befund: **RTM und Dogfooding.** `make trace` liefert 48
  Anforderungen und 0 Waisen; die neue Anforderung ist auf den
  zusammengeführten Folge-Slice abgebildet. Der Lauf des ausgelieferten Image
  gegen das Repo meldet 343 geprüfte Dateien und 0 Befunde — die Anker- und
  Linkpflicht der neuen Passagen hält.
- geprüft, ohne Befund: **Befund-Dedup unter `sections: each`.** Weil Schritt 7
  `line` auf die Zeile der Abschnitts-Überschrift pinnt, tragen zwei verletzende
  Abschnitte derselben Datei verschiedene Tupel; die Zusage „genau ein Befund für
  den verletzenden Abschnitt" ist unter dem vorhandenen `SortFindings` erreichbar.
  Die verbleibende Kollision betrifft ausschließlich `section-missing` zweier
  Regeln mit gleichem Glob (Umsetzbarkeit F-2, „teilweise").
- geprüft, ohne Befund: **Status der ADR.** Sie ist unverändert `Proposed`; die
  nachträgliche Erweiterung verletzt damit keine Immutabilitäts-Regel, und der
  Geschichts-Eintrag hält Anlass und Umfang der Revision fest. Der `adr`-Gate-Lauf
  bleibt davon unberührt.
- geprüft, ohne Befund: **Hermetik-, Opt-in- und Determinismus-Zusagen.** Die
  Neufassung ändert an „byte-identisch ohne aktives Modul", „nur lesend,
  netzlos", „diagnose-only" nichts; die erweiterte Kandidaten-Menge (Baumlauf
  unabhängig vom Scan-Scope) bleibt lesend und deterministisch und berührt die
  Messmethode von
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) nicht.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 7 |
| LOW | 4 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Änderungsprotokoll nicht mit der Änderung
mitgeführt · Index-Zusammenfassung nicht mit der Revision mitgeführt ·
Nachträglich erweiterte ADR nur im vorderen Teil nachgezogen ·
Entscheidungs-Protokoll nicht mit der Revision mitgeführt · Zugesagte Invariante
ohne Akzeptanzkriterium · Aggregat widerspricht der Einzelaufstellung ·
Lauf-Beleg nachträglich verändert · DoD-Zusage ohne Vertrags-Deckung ·
Aufzählung nach Datei-Entfall nicht nachgezogen · Messzahl ohne Methoden- und
Standangabe · Guard bei der Neufassung verloren · Zwei Lexiken derselben Sektion
mit ungleicher Whitespace-Regel

Wiederholung aus den Vor-Läufen (Steering-Loop-Signal): „Zugesagte Invariante
ohne Akzeptanzkriterium" und „Aggregat widerspricht der Einzelaufstellung"
treten hier zum **zweiten** Mal auf, „Teil-Nachzug einer Semantik-Änderung"
sinngemäß zum **dritten** (Vertrags-F-11, Umsetzbarkeits-F-5 und jetzt F-N1 bis
F-N4). Die dritte Wiederholung derselben Klasse ist laut Skill ein Signal, den
Sensor nachzuziehen statt nur zu melden — hier: eine Nachzugs-Liste für die
Stellen, die eine Vertrags-Revision immer berührt (beide §7-Historien,
ADR-Index, ADR-Konsequenzen, Slice-Abnahme-Punkte, Wellendokument, Glossar).

## Verdikt

**Merge-blockierend: ja** — sieben MEDIUM, kein HIGH.

Die inhaltliche Arbeit trägt. Die beiden strukturellen Befunde, die scharf zu
prüfen waren, sind **echt geheilt**, nicht bloß beschriftet: `sections` mit
`one`/`each` macht die Dokumentklasse prüfbar, die den Modul-Schnitt begründet
hat, und der Modus ist an allen vier Vertragsstellen widerspruchsfrei
beschrieben, einschließlich des Leerlauf-Verhaltens bei null Treffern. Die
Marken-Verankerung ist gemessen statt angenommen; die Gegenprobe gegen beide
Repos bestätigt, dass die drei zugesagten Formen den realen Bestand decken und
keine neue Falsch-Positiv-Klasse öffnen. Auch die schwereren
Umsetzbarkeits-Befunde sind sauber adressiert: je Bedingung ein Grund-Code plus
gepinnte Befund-Zeile lösen die Dedup-Kollision im Hauptfall, die
Kandidaten-Menge ist definiert, die Anker-Semantik entschieden, und die
Zusammenführung der beiden Folge-Slices ist die richtige Antwort auf die
Release-Grenze. Von 32 Vor-Befunden sind 14 vollständig erledigt.

Blockierend ist etwas anderes, und es ist ein Muster, kein Einzelfall: **die
Revision hat den Körper der Dokumente korrigiert und ihre Ränder stehen
lassen.** Die Historie-Zeile der Spezifikation, der ADR-Index, der
Konsequenzen-Abschnitt der ADR und der Abnahme-Punkt des Slice-Plans beschreiben
alle vier weiterhin die Fassung, die dieser Commit gerade verworfen hat (F-N1
bis F-N4). Jede dieser Stellen ist die Fläche, die ein Leser **statt** des
Körpers liest — ein Index, ein Änderungsprotokoll, ein Entscheidungs-Protokoll.
Dazu kommen zwei Zusagen, die im Diff behauptet, aber nicht eingelöst sind (F-N5
Akzeptanzkriterium „beider Anforderungen", F-N6 Grund-Code-Zählung), und ein
Eingriff in zwei abgeschlossene Lauf-Belege, der einen Link mit falschem Label
hinterlässt (F-N7).

Von den zwölf „teilweise"-Urteilen sind drei substanziell und nicht nur
Nachzugs-Arbeit: die verbleibende Dedup-Kollision zweier Regeln mit gleichem
Glob (Umsetzbarkeit F-2), die offene Frage nach Nicht-Markdown-Kandidaten
(Umsetzbarkeit F-4) und die Widerspruchslage zu
[`DC-FA-CONF-002`](../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope)
(Umsetzbarkeit F-10), die vom stillen zum ausgesprochenen Widerspruch zweier
Anforderungen desselben Stratums geworden ist. Der ausgelieferte
Silent-Grün-Pfad im Fence-Automaten ist korrekt als eigener Slice ausgelagert,
aber die Reihenfolge, die verhindert, dass das neue Modul ihn erbt, ist nirgends
bindend formuliert.

**Zeitkritisch:** F-N2 und F-N3 wandern mit dem `Accepted`-Übergang der ADR in
ein immutables Dokument bzw. in eine Index-Zeile, die dann eine widerlegte
Zusammenfassung festschreibt ([`AGENTS.md`](../../AGENTS.md) §3.5) — sie sind
**vor** der Statusänderung zu klären, nicht danach. Dass die ADR noch `Proposed`
war, hat diese Runde überhaupt erst möglich gemacht; dasselbe Fenster steht für
diese Findings noch offen.

**Übergabe:** Die Findings gehen an den Implementer; die Rückkante Review → Plan
betrifft F-N4 (Abnahme-Punkt 4 und zwei Messzeilen sind Entscheidungs-Protokoll,
nicht Prosa) und F-N6/F-N8 (DoD und Vertrag von
[slice-099](../plan/planning/open/slice-099-structure-modul.md)). Die
Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den
Zähler. Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses
Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat.
