# Review-Report — slice-233 (`links`: Referenz-Definitionen `[label]: ziel`), R2

**Review-Art:** Code (Korrektur-Verifikation gegen R1-H1, ADR-0095, Hard Rules).

**Gegenstand:** `slice-233`, Fix-Range `a83a97d0..022e79f1` (zwei Commits:
`e1f58b2d` Vertrag/ADR-0095 (supersedes ADR-0094), `022e79f1` Code-Fix); zur
ADR-Immutabilitäts- und Bestandsprüfung zusätzlich die volle Slice-Range
`6f8d6906..022e79f1`.

**Skill:** `.harness/skills/reviewer.md` @ Version 1.16.0, 2026-09-07.

**Modell-ID:** `claude-sonnet-5`.

**Datum:** 2026-09-27.

**Eingangs-Kontext:**

- R1-Report `docs/reviews/2026-09-27-slice-233-links-referenz-definitionen-review-r1.md`
  (vollständig gelesen) — Befund R1-H1 (HIGH): `definitionRe` validierte nicht,
  dass der Zeilenrest hinter dem Ziel ein delimitierter Titel ist; `[TERM]:
  First In, First Out` wurde fälschlich als Definition mit erfundenem Ziel
  „First" erkannt.
- `docs/plan/adr/0094-…backslash-grenze-korrigiert.md` (Status `Superseded by
  ADR-0095`) und `docs/plan/adr/0095-…titel-delimiter-pflicht.md` (Status
  `Accepted`), beide vollständig gelesen.
- `spec/lastenheft.md` (Version 0.93.2, Historie-Zeile 0.93.2) und
  `spec/spezifikation.md` (§`DC-FA-LINK-001.a` Schritt 3, präzisierte Fassung).
- `AGENTS.md` §3.5 (ADR-Immutabilität), §3.7 (Kommentar-Klassen), §3.4
  (Referenzrichtung).
- **Nicht** erhalten (laut Skill): die DoD-Abhakung selbst.

**Gates unabhängig nachgefahren** (Docker/`make`, kein Host-Go — AGENTS.md §3.1):

| Gate | Ergebnis |
|---|---|
| `make gates` (voll, zehn Glieder) | grün — `baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; Coverage `94.30 %` (Schwelle `93 %`); `semgrep`: 55 Regeln/65 Dateien, 0 Befunde; `doc-check`: `849 Datei(en) geprüft, 0 Befund(e)` (zwei ADR-Dateien mehr als R1s 847) |
| `make adr-check RANGE=6f8d6906..022e79f1` (volle Slice-Range, nicht nur Fix-Range) | grün — `849 Datei(en) geprüft, 0 Befund(e)` (Modul `vcs`) |
| `go test ./...` mit temporärer, nicht committeter Adversarial-Testdatei gegen den fertigen Fix-Stand (siehe unten) | siehe Abschnitt „Adversarielle Regex-Prüfung" |
| Runtime-Image aus `git worktree` auf `6f8d6906` (vor slice-233) und Runtime-Image auf `022e79f1` (fertiger Fix), beide gegen **denselben** aktuellen Bestand, `--enable external --enable tracked` | je `849 Datei(en) geprüft, 107 Befund(e)` — identisch, 0 neue Befunde durch das Feature auf dem heutigen Bestand |
| Runtime-Image auf `82c941ab` (Feature vor R1-Fix) gegen Fixture `[TERM]: First In, First Out` | `1 Befund(e)` (`First target-missing`) — Rot-Beleg reproduziert |
| Runtime-Image auf `022e79f1` (nach Fix) gegen dasselbe Fixture | `0 Befund(e)` — Grün-Beleg reproduziert |
| Runtime-Image auf `022e79f1` gegen zwei neu konstruierte Fixtures `[TODO]: (spaeter)` / `[Status]: (offen)` | `2 Befund(e)` (`target-missing` je Datei) — siehe R2-M1 |

---

## Zentraler Prüfpunkt: Ist R1-H1 jetzt strukturell ausgeschlossen?

### Methodik

Die neue Regex wurde nicht nur gegen den mitgelieferten Test geprüft, sondern
gegen 22 selbst konstruierte Fälle — teils über eine temporäre, **nicht
committete** Testdatei `internal/hexagon/core/rules/zzz_r2adversarial_test.go`
(direkt gegen `definitionRe`/`ExtractLinks` per `go test` in Docker, danach
gelöscht — Arbeitsbaum vor und nach der Prüfung sauber, `git status` bestätigt
leer), teils direkt gegen den gebauten Container mit Fixture-Dateien.

### Ergebnis: Die konkrete R1-H1-Klasse ist ausgeschlossen — die Fehlerklasse als Ganzes nicht vollständig

**Ausgeschlossen (11 Fälle, alle korrekt `nil`/kein Fund):**

| Fixture | Ergebnis |
|---|---|
| `[TERM]: First In, First Out` (Original-Repro R1-H1) | kein `LinkRef` |
| `[Siehe]: docs/plan/adr/README.md fuer den Index` | kein `LinkRef` |
| `[Hinweis]: Diese Datei ist noch nicht fertig` | kein `LinkRef` |
| `[TODO]: spaeter nachtragen` | kein `LinkRef` |
| `[Beispiel]: siehe unten (Kapitel 3)` | kein `LinkRef` (Klammer-Titel **nicht** unmittelbar nach dem Ziel-Token — korrekt verworfen, keine naive „Klammer irgendwo in der Zeile"-Erkennung) |
| `[l]: ziel.md "eins" "zwei"` (zwei Titel-Kandidaten) | kein `LinkRef` (korrekt verworfen — nur ein Titel erlaubt) |
| `[x]: ziel.md "titel ohne schluss` (unterminierter Titel) | kein `LinkRef` |
| `[x]: ziel.md(kein leerzeichen)` | kein `LinkRef` |
| `[l]: ziel.md "Der Titel \"mit\" Anfuehrungszeichen"` (eingebettete Quotes) | kein `LinkRef` — siehe R2-L1 (Kehrseite: falsch-negativ statt falsch-positiv) |
| `[l]: ziel.md "Titel"` mit Trailing-Whitespace/-Tabs | korrekt erkannt, `Target=ziel.md` |
| `[l]: ziel.md ""` (leerer Titel) | korrekt erkannt, `Target=ziel.md` |

Die **Fünf Prosa-Gegenproben** (Punkt „Groß angelegte Prosa-Gegenprobe" im
Auftrag) sind darunter: „Siehe/Hinweis/TODO/Beispiel"-Zeilen — realistische,
selbst konstruierte Analoga zu Glossar-/Verweis-Prosa, da `AGENTS.md`,
`spec/lastenheft.md`, `spec/spezifikation.md` und `docs/user/*.md` **keine**
echte Zeile der Form `^\s{0,3}\[…\]:` enthalten (per `grep` bestätigt) — der
Bestand liefert selbst keine Gegenprobe, die Konstruktion ersetzt sie. Alle
vier bestehen.

**Nicht ausgeschlossen — zwei neue, unbenannte Restklassen (R2-M1, R2-L1
unten):** ein **einzelnes, Whitespace-freies** Token nach dem Ziel-Doppelpunkt,
das selbst wie ein Titel aussieht (`[TODO]: (spaeter)`, `[Status]: (offen)`),
wird **als Ziel** erkannt (nicht als Titel, weil kein Whitespace davor die
optionale Titel-Gruppe eröffnet) — verifiziert am laufenden Container:
`target-missing` auf `(spaeter)` bzw. `(offen)`. Ebenso `[l]:
ziel(unbalanced.md` und `[note]: draft(wip` (unbalancierte Klammer im bloßen
Ziel-Token, kein Whitespace).

**Fazit zur zentralen Prüffrage:** Die **konkrete** R1-H1-Reproduktion (mehrwortige
Prosa ohne jede Klammer-/Anführungszeichen-Struktur direkt nach dem
Ziel-Token) ist strukturell ausgeschlossen und durch elf unabhängige
Gegenproben bestätigt — deutlich robuster, als es der einzelne mitgelieferte
Test allein zeigen könnte. Die **Fehlerklasse als Ganzes** — eine
Annotations-Zeile ohne Link-Absicht wird zu einem erfundenen `LinkRef` — ist
jedoch **nicht** vollständig eliminiert: Ein einzelnes klammerartiges
Whitespace-freies Token danach reproduziert dasselbe Grundmuster in
schmalerer Form (siehe R2-M1). ADR-0095s Aussage „eliminiert die Fehlerklasse
durch Konstruktion" trifft für die **behandelte** Klasse (mehrwortiger
Titel-Rest) zu, ist als Aussage über *die* Fehlerklasse aber zu weit gefasst.

---

## Weitere Pflicht-Checks

### 1. ADR-Form/Prozess

- `git show e1f58b2d -- docs/plan/adr/0094-…md` ändert **ausschließlich** die
  `**Status:**`-Zeile (`Accepted` → `Superseded by ADR-0095`) und hängt eine
  Zeile an `## Geschichte` an — Kern (Kontext, Entscheidung, Alternativen,
  Konsequenzen, Fitness Function, Re-Evaluierungs-Trigger) byte-identisch.
  `make adr-check RANGE=6f8d6906..022e79f1` bestätigt unabhängig, über die
  **volle** Slice-Range (0 Befunde).
- ADR-0095 selbst: **drei** verglichene Alternativen mit Trade-offs (Tabelle),
  Fitness Function (zwei `go test`-Zeilen mit Make-Target), Re-Evaluierungs-Trigger
  (eine benannte Bedingung, sonst „permanent"), `Supersedes: ADR-0094` gesetzt,
  `Schärft:` zeigt aufwärts auf `DC-FA-LINK-001`. Zwei
  `<!-- d-check:status-provenance -->`-Marker im `Bezug:`-Feld zeigen korrekt
  *wo* verifiziert/entstanden (Review-Report, Slice) — keine getarnte
  Entscheidungsgrundlage.
- ADR-Index (`docs/plan/adr/README.md`): beide Zeilen im selben Commit
  aktualisiert (ADR-0094 → `Superseded by ADR-0095`, neue ADR-0095-Zeile).
- **Vollständig, keine Findings.**

### 2. Bestandsmessung — nachgefahren gegen den aktuellen Endstand

Wie im Auftrag verlangt: Runtime-Image aus `git worktree` auf `6f8d6906`
(Vor-Slice-Stand) **und** Runtime-Image auf `022e79f1` (jetziger Endstand),
beide gegen **denselben** aktuellen Arbeitsbaum, `--enable external --enable
tracked`:

- Alt (`6f8d6906`): `849 Datei(en) geprüft, 107 Befund(e)`.
- Neu (`022e79f1`): `849 Datei(en) geprüft, 107 Befund(e)`.
- **Identisch — 0 neue Befunde**, bestätigt auf dem finalen (nicht nur dem
  R1-Zwischen-)Stand. Die Dateizahl liegt bei 849 statt R1s 847, weil in der
  Zwischenzeit zwei ADR-Dateien (ADR-0094/0095 zusammen mit den
  Vertrags-Änderungen) hinzukamen — beide Läufe zählen dieselben 849, die
  Differenz ist kein Messfehler.
- Ergänzend geprüft: keine der beiden im Repo tatsächlich vorkommenden
  Zeilen entspricht den zwei neuen Restklassen aus R2-M1 (`grep -rnE
  '^\s{0,3}\[[^]]+\]: *\('` über alle `*.md` außerhalb der vendorten Baseline:
  0 Treffer) — die Restklassen sind für **dieses** Repo folgenlos, wie schon
  R1-H1 es vor seiner Korrektur war.

### 3. Kommentar-Disziplin (`AGENTS.md` §3.7)

**Neuer Befund — siehe R2-H1 unten.** Zwei der neuen Kommentare
(`markdown.go:578`, `markdown_test.go:243-245`) tragen den Text „R1-H1" bzw.
„unabhaengiger Review, R1-H1" — ein Review-Befund-Marker im Code-Kommentar,
zusätzlich zur (für sich zulässigen) ADR-Referenz. Das ist eine
Verschlechterung gegenüber R1s eigenem Befund zu diesem Slice: R1 bescheinigte
den damaligen Kommentaren (vor dem Fix) saubere Disziplin — der Fix, der
R1-H1 behebt, führt selbst einen neuen §3.7-Verstoß ein. Der neue Kommentar in
`cli_acceptance_test.go` (nur „ADR-0093/0094/0095") ist dagegen sauber.

### 4. `make gates` und `make test` unabhängig grün

Siehe Gates-Tabelle oben — beide grün, echte Ausgabe zitiert (Coverage
94,30 %, `doc-check` 849/0, `semgrep` 0/55).

### 5. R1-L1 (matrix-Lineage über Definitions-Label) — weiterhin zutreffend

`matrix.go:77` (`lineageExempt(supersedeValues, ref.Text, rel)`) ist von
diesem Fix-Commit **nicht** berührt (`git show 022e79f1 --stat` listet nur
`markdown.go`/`markdown_test.go`/`cli_acceptance_test.go`). Die in R1-L1
beschriebene Konsequenz — eine Definition mit Label „ADR-0001" erhält
dieselbe Supersede-Lineage-Ausnahme wie ein Inline-Link mit demselben
Linktext — besteht unverändert und ist weiterhin in keiner ADR/Spec
ausbuchstabiert. Nicht merge-blockierend, wie in R1 bewertet.

### Weitere geprüfte Punkte ohne Befund

- **Hexagon-Import-Richtung (ADR-0005):** keine neuen Imports im Diff.
- **Gate-Suppression/Schwellen-Senkung (§3.2/§3.6):** kein `//nolint`, keine
  Schwellen-Änderung.
- **Netzzugriff außerhalb `external`:** keiner.
- **Commit-Zerlegung:** `e1f58b2d` (Vertrag) ändert ausschließlich
  `spec/`/`docs/plan/adr/`, kein Code; `022e79f1` (Fix) ändert ausschließlich
  `internal/` — sauber getrennt, wie im etablierten Muster dieses Slice.
- **§3.4/MR-006 Referenzrichtung (ADR/Slice/Welle/Commit-Hash im Spec-Körper):**
  `grep -n "ADR-009[345]\|slice-233" spec/lastenheft.md spec/spezifikation.md`
  liefert keinen Treffer im Fließtext. **Beobachtung, nicht meldepflichtig:**
  die neue Historie-Zeile 0.93.2 nennt „Der Review (R1-H1, HIGH)" — eine
  Review-Report-Kennung ist keine der fünf in §3.4 benannten Kategorien
  (ADR/Slice/Welle/Commit-Hash/Closure-Datum) und wird vom `matrix`-Modul
  nicht als Klasse geführt; kein Konventions-Anker, der das verböte. Als
  INFO vermerkt, nicht als Finding gewertet.
- **MR-032/Decken-Regel:** Historie-Zeile 0.93.2 nennt weder ADR- noch
  Slice-Nummer im Fließtext, `Verweis`-Spalte `—` (implizit über die
  Tabellenstruktur, wie bei 0.93.0/0.93.1).
- **Zustandsfelder:** keine inhaltlich veränderte Roadmap-/Register-Zeile
  außer dem Status-Übergang der ADR (zulässig, §3.5/§3.7-Ausnahme).
- **`DC-QA-04` Alt-Tool-Migrationsabdeckung:** nicht berührt.
- **Modul-Scan-Grenze (§3.8):** nicht einschlägig, wie in R1.
- **Fitness Function aus ADR-0095 nachvollzogen:** beide `go test`-Zeilen
  (Rot-Beleg gegen die Vorfassung, Bestandstests bleiben grün) durch `make
  test`/eigene Läufe bestätigt.

---

## Findings

| # | Kategorie | Quelle | Pfad | Befund | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| R2-H1 | HIGH | `AGENTS.md` §3.7 (Kommentar-Klassen) / Prüffrage 6 | `internal/hexagon/core/rules/markdown.go:578`, `internal/hexagon/core/rules/markdown_test.go:243,245` | Zwei neue Kommentare tragen den Text „R1-H1" bzw. „unabhaengiger Review, R1-H1" — ein Review-Befund-Marker im Code-Kommentar. `AGENTS.md` §3.7 lässt Herkunft nur als **ein** auflösbares Feld nach dem Baseline-Schema (`DC-*`/`ADR-*`/`MR-*`/`seit welle-<NN>`) zu und verbietet ausdrücklich „Review-Historie" und „Review-Befund-Marker" im Kommentar; die parallel vorhandene, für sich zulässige ADR-0095-Referenz macht die zusätzliche Review-ID-Nennung nicht zulässig. Derselbe Anti-Pattern findet sich bereits (grandfathered) an vier weiteren Stellen dieses Repos (`cli_acceptance_test.go:1353`, `tracked_test.go:36,57,77`) — dieser Fix fügt zwei **neue** Vorkommen hinzu, statt die Klasse zu meiden. | ja — `grep -n "R1-H1" internal/hexagon/core/rules/markdown.go internal/hexagon/core/rules/markdown_test.go` zeigt beide Fundstellen; kein Gate fängt das (Urteil, kein `grep`, siehe Skill) | kommentar-traegt-review-befund-marker |
| R2-M1 | MEDIUM | Prüffrage 2 (Kern-Modul meldet falsch) / ADR-0095 §Verglichene Alternativen | `internal/hexagon/core/rules/markdown.go:579` (`definitionRe`) | Ein einzelnes, whitespace-freies Token direkt nach dem Ziel-Doppelpunkt, das selbst klammerartig aussieht (`[TODO]: (spaeter)`, `[Status]: (offen)`), wird als **Ziel** erkannt statt verworfen, weil die optionale Titel-Gruppe nur nach vorangehendem Whitespace greift — ohne dieses Whitespace absorbiert `\S+` das gesamte Token als Ziel. Verifiziert am gebauten `022e79f1`-Container: beide Fixtures erzeugen `target-missing`. Dieselbe Grundfigur wie R1-H1 (Annotations-Zeile ohne Link-Absicht → erfundenes Ziel), nur auf ein schmaleres Muster (ein Token, kein Whitespace) verengt. Zusätzlich verifiziert: ein Ziel-Token mit unbalancierter Klammer ohne Whitespace (`[l]: ziel(unbalanced.md`, `[note]: draft(wip`) wird ebenfalls unverändert als Ziel akzeptiert, obwohl echtes CommonMark unbalancierte, unescapte Klammern in einem undelimitierten Ziel nicht erlaubt — hier ist d-check nachweisbar **großzügiger** als die eigene „nach CommonMark"-Zusage im Kommentar über `definitionRe`. ADR-0095 benennt in §Verglichene Alternativen nur die Grenze „nicht delimitierter Titel bleibt unerkannt" — diese beiden Restklassen (whitespace-freies Klammer-Token als Ziel; unbalancierte Klammern im Ziel-Token) sind dort nicht aufgeführt, ADR-0095s Aussage „eliminiert die Fehlerklasse durch Konstruktion" ist insofern eine Übergeneralisierung über die tatsächlich geschlossene (schmalere) Lücke. Für den aktuellen Bestand dieses Repos folgenlos (0 Treffer per `grep`), betrifft aber wie schon bei R1-H1 fremde/künftige Dokumente. | ja — Fixtures `todo.md`/`status.md` mit Inhalt `[TODO]: (spaeter)` bzw. `[Status]: (offen)`, `d-check` ohne Flags darauf ausführen: je ein `target-missing`-Befund | definitionRe-bare-token-nicht-klammer-validiert |
| R2-L1 | LOW | Konsistenz mit CommonMark-Zusage im Kommentar | `internal/hexagon/core/rules/markdown.go:579` (`definitionRe`) | Ein Titel mit einem via Backslash „escapten" eingebetteten Anführungszeichen (`[l]: ziel.md "Der Titel \"mit\" Anfuehrungszeichen"`) lässt die **ganze** Zeile unerkannt bleiben (verifiziert: `nil`), obwohl echtes CommonMark Backslash-Escapes innerhalb von Link-Titeln unterstützt und diese Zeile dort als gültige Definition mit einem toten `ziel.md`-Ziel gelesen würde. `[^"\n]*` in der Titel-Alternative kennt keine Escape-Semantik und bricht am ersten (auch escapten) `"`. Das ist eine falsch-**negative** Abweichung (ein echtes totes Ziel bliebe hier unentdeckt) statt einer falsch-positiven wie bei R1-H1 — geringere Schwere, da ein Ausbleiben eines Befundes stiller, aber nicht irreführend ist. Nicht in ADR-0095 oder der Spezifikation als Grenze benannt. | ja — Fixture mit obigem Inhalt, `d-check`: `0 Befund(e)`, obwohl `ziel.md` nicht existiert | definitionRe-keine-escape-semantik-im-titel |
| R2-L2 (= R1-L1, unverändert) | LOW | Doku-Präzision, `matrix.go:77` | `docs/plan/adr/0093-links-referenz-definitionen-gemeinsame-extraktion.md` | Wie in R1 berichtet — `matrix` nimmt für eine Definition auch deren Label in die Supersede-Lineage-Ausnahme auf, was ADR-0093 nicht ausbuchstabiert. Von diesem Fix-Commit nicht berührt, weiterhin zutreffend. | ja — `grep -n "ref.Text" internal/hexagon/core/rules/matrix.go` | adr-konsequenz-nicht-ausbuchstabiert |

## Negativbefunde (geprüft, ohne Befund)

- **Gates:** `make gates` (zehn Glieder, echte Ausgabe oben) und unabhängig
  `make adr-check RANGE=6f8d6906..022e79f1` (volle Slice-Range) grün.
- **Zentrale R1-H1-Reproduktion:** 11 unabhängige Gegenproben (davon 4
  realistische Prosa-Konstrukte, da der Bestand selbst keine liefert) bestehen
  alle — die spezifische, im R1-Report demonstrierte Fehlerform ist
  strukturell ausgeschlossen.
- **Rot-/Grün-Beleg der Korrektur, unabhängig nachgefahren:** Alt-Image
  (`82c941ab`, vor dem Fix) reproduziert `First target-missing`; Neu-Image
  (`022e79f1`) liefert `0 Befunde` auf demselben Fixture.
- **Bestandsmessung, inklusive `external`/`tracked`, gegen den Endstand:**
  849/107 auf Alt- (`6f8d6906`) vs. Neu-Binary (`022e79f1`), identisch — 0
  neue Befunde.
- **ADR-Immutabilität (§3.5):** `1fad1b64`-artiger Diff auf ADR-0094 betrifft
  ausschließlich Status-Feld und Geschichte-Anhang; `make adr-check` über die
  volle Range bestätigt (0 Befunde).
- **ADR-0095 Form (Modul 4):** drei Alternativen, Fitness Function,
  Re-Evaluierungs-Trigger, `Schärft:` aufwärts, `Supersedes`-Feld,
  Provenance-Marker korrekt (zeigt wo verifiziert, keine getarnte
  Entscheidungsgrundlage) — vollständig; ADR-Index aktualisiert.
- **Hexagon-Import-Richtung (ADR-0005):** keine neuen Imports.
- **Gate-Suppression/Schwellen-Senkung (§3.2/§3.6):** keine `//nolint`, keine
  Schwellen-Änderung.
- **Netzzugriff außerhalb `external`:** keiner.
- **Commit-Zerlegung:** Vertrags-Commit ohne Code, Fix-Commit ausschließlich
  `internal/`.
- **§3.4/MR-006 Referenzrichtung (mechanisierte Kategorien):** kein
  ADR-/Slice-/Welle-/Commit-Hash-Token im Körper von Lastenheft oder
  Spezifikation (Review-Report-Erwähnung siehe INFO oben, nicht
  meldepflichtig — keine der fünf benannten Kategorien).
- **MR-032 Decken-Regel:** Historie-Zeile 0.93.2 ohne ADR-/Slice-Nennung im
  Fließtext.
- **Zustandsfelder:** keine inhaltlich veränderte Roadmap-/Register-Zeile
  außer dem zulässigen ADR-Status-Übergang.
- **`DC-QA-04` Alt-Tool-Migrationsabdeckung:** nicht berührt.
- **Modul-Scan-Grenze (§3.8):** nicht einschlägig.
- **Provenance-Marker:** beide `<!-- d-check:status-provenance -->`-Marker in
  ADR-0095 zeigen nur, *wo* die Entscheidung/Verifikation entstand.
- **Kommentar-Disziplin, Rest der Datei:** `cli_acceptance_test.go`s neuer
  Kommentar trägt ausschließlich `ADR-0093/0094/0095` — sauber (im
  Unterschied zu `markdown.go`/`markdown_test.go`, siehe R2-H1).

## Kategorie-Summary

- HIGH: 1 (R2-H1)
- MEDIUM: 1 (R2-M1)
- LOW: 2 (R2-L1, R2-L2 = R1-L1 fortgeführt)
- INFO: 1 (Review-Report-Kennung in Spec-Historie, unbenannte Kategorie, nicht meldepflichtig)

## Verdikt

**Merge-blockierend: ja, wegen R2-H1.**

Begründung: R2-H1 ist ein klarer, mechanisch nachweisbarer Verstoß gegen
`AGENTS.md` §3.7 (zwei neue Code-Kommentare tragen einen Review-Befund-Marker,
„R1-H1") — dieselbe Prüffrage, die R1 selbst an diesem Slice noch sauber
vorfand. Er ist trivial zu beheben (Text entfernen, ADR-0095-Referenz allein
genügt als Herkunfts-Feld) und stellt for sich genommen kein
Korrektheitsrisiko dar, ist aber ein Hard-Rule-Verstoß und damit nach dem im
Reviewer-Skill festgelegten Kategorien-Schema blockierend.

**Zur zentralen Frage des Auftrags:** R1-H1 selbst — die konkrete,
demonstrierte Fehlform (mehrwortige Prosa als Definition mit erfundenem Ziel)
— ist jetzt **strukturell ausgeschlossen**, bestätigt durch elf unabhängige
Gegenproben zusätzlich zum mitgelieferten Test. Die **Fehlerklasse im
weiteren Sinne** (Annotations-Zeile ohne Link-Absicht → erfundenes Ziel) ist
jedoch **nicht vollständig** geschlossen: R2-M1 zeigt eine schmalere, aber
reale Restform (ein einzelnes klammerartiges Token ohne Whitespace). R2-M1 ist
nicht für sich merge-blockierend (schmalerer Trigger, auf dem aktuellen
Bestand folgenlos, teilweise sogar CommonMark-konform statt eindeutig
falsch), gehört aber vor der nächsten Closure entweder in ADR-0095/Spec als
benannte Grenze nachgetragen oder durch eine weitere Regex-Schärfung
geschlossen — dieselbe Wahl, die ADR-0095 für R1-H1 bereits getroffen hat.
R2-L1 (Escape-Semantik im Titel) und R2-L2 (fortgeführtes R1-L1) sind
nice-to-fix, nicht blockierend.
