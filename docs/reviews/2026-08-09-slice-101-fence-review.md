# Review-Report: slice-101 — 2026-08-09

**Review-Art:** Code — geprüft wird der Diff gegen Slice-Plan, Anforderung,
Spezifikation und ADR (Modul 10 §Drei Review-Arten).

**Gegenstand:** [slice-101](../plan/planning/in-progress/slice-101-fence-unbalanciert.md),
Diff-Range `38c36b2..aba7792` (drei Commits: Bestandsmessung, CR + ADR,
Implementierung).

**Skill:** `.harness/skills/reviewer.md` @ 1.3.0 ·
**Modell:** claude-opus-5[1m] · **Datum:** 2026-08-09

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [slice-101](../plan/planning/in-progress/slice-101-fence-unbalanciert.md)
  (Ziel, Bestandsmessung §3, Abnahme-Punkte §4, Risiken §6)
- [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md) (Entscheidungen 1–5,
  Fitness Function, Konsequenzen) und
  [ADR-0042](../plan/adr/0042-markdown-lexik-folgt-commonmark.md) (die bewusst
  offen gelassene Fence-Grenze)
- [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
  (Klasse 3 + zwei neue Akzeptanzkriterien) und
  [§`DC-FA-SPAN-001.a`](../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
  Schritt 3 sowie die §4-Grund-Code-Zeile
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (der Konsument, dessen Silent-Grün den Anlass gab)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules)

**Läufe dieses Reviews** (alle gegen das aus dem Diff-HEAD gebaute Image bzw.
über `make`-Targets; Fixtures in einem Temp-Verzeichnis außerhalb des Repos):
`make build`, `make test`, `make lint`, `make doc-check`,
`docker run --network none … d-check:latest` gegen 16 Fixtures.

---

## Findings

### F-1 — Der Wächter wertet einen anderen Fence-Automaten aus als der Tabellen-Leser; das Vollständigkeits-Gate bleibt still grün

- `kategorie`: HIGH
- `quelle`: [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
  (Klasse 3: „Alles dahinter gilt für **jede** Vorverarbeitung als Code und wird
  von **allen** Modulen übersprungen"),
  [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md) §Konsequenzen
  („Jedes Modul profitiert, weil die Lexik geteilt ist")
- `pfad`: `internal/hexagon/core/rules/spans.go:39` gegen
  `internal/hexagon/core/app/trace_table.go:340`
- `befund`: Das Produkt trägt **zwei** Fence-Automaten mit unterschiedlicher
  Schluss-Regel — `checkUnclosedFence` wertet den naiven Toggle aus,
  `markdownTableLines` schließt zeichen- **und** längenabgeglichen. Eine Datei,
  in der ein Backtick-Öffner von einer Tilde-Zeile „geschlossen" wird, ist für
  den Wächter balanciert (0 Befunde, Exit 0), für den Tabellen-Leser bis zum
  Dateiende offen: `--trace --require-complete` meldet dort „2 Anforderung(en),
  0 Waise(n)" und **Exit 0**, obwohl eine ungedeckte Anforderung hinter dem
  Fence steht; dieselbe Datei mit passendem Schluss liefert 3 Anforderungen,
  1 Waise und Exit 1.
- `verifizierbar`: ja — Fixture mit `requirements.format: table`; die Zeilen
  1–6 tragen zwei gedeckte Anforderungen, Z. 10 einen Backtick-Öffner mit
  Infozeile `yaml`, Z. 12 eine Zeile aus drei Tilden, danach eine zweite
  Anforderungstabelle mit der ungedeckten Kennung. Lauf A (Tilden-Zeile):
  `2 Anforderung(en), 0 Waise(n)`, Exit 0, und der `spans`-Lauf derselben Datei
  `0 Befund(e)`, Exit 0. Lauf B (dieselbe Datei, Schluss-Zeile aus drei
  Backticks): `3 Anforderung(en), 1 Waise(n)`, Exit 1. Der Effekt trägt bis in
  `make completeness-check`.
- `klasse`: „Wächter misst einen anderen Automaten als der bewachte Verbraucher"

### F-2 — Wächter und Vorverarbeitung trimmen verschieden; eine einzige Zeile schaltet den Befund ab

- `kategorie`: HIGH
- `quelle`: [§`DC-FA-SPAN-001.a`](../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
  Schritt 3 („derselbe Fence-Automat wie in der Vorverarbeitung")
- `pfad`: `internal/hexagon/core/rules/spans.go:39` gegen
  `internal/hexagon/core/rules/markdown.go:54`
- `befund`: Der Wächter trimmt mit `strings.TrimSpace` (unicode-weit), die
  Vorverarbeitung `proseLines` mit `strings.TrimLeft(raw, " \t")`. Eine mit
  Unicode-Whitespace eingerückte Fence-Zeile (belegt mit U+00A0) ist deshalb
  **nur** für den Wächter ein Toggle: in einer Datei mit einer solchen Zeile und
  einem echten, nie geschlossenen Fence meldet der Lauf 0 Befunde und Exit 0,
  obwohl die Vorverarbeitung ab dem echten Öffner alles überspringt.
- `verifizierbar`: ja — Fixture: Z. 3 besteht aus U+00A0 gefolgt von drei
  Backticks, Z. 7 ist ein echter Öffner, Z. 9 ein kaputter Link. Ergebnis
  `0 Befund(e)`, Exit 0; dieselbe Datei ohne die Zeile 3 liefert
  `a.md:5 fence-unclosed`, Exit 1. Der kaputte Link wird in beiden Fällen
  verschluckt.
- `klasse`: „Wächter und Bewachtes trimmen verschieden"

### F-3 — Die dokumentierte Grenze deckt nur die Quell-Achse; die Ziel-Achse regulärer Module hat dieselbe Asymmetrie

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
  Absatz „Grenze"; [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md)
  §Konsequenzen
- `pfad`: `internal/hexagon/core/rules/matrix.go:436` (`statusOf`), Vertragstext
  in `spec/lastenheft.md`
- `befund`: Als Grenze benannt sind nur „Module, die ihre Eingabe selbst
  benennen (Post-Pässe über deklarierte Verzeichnisse)". Dieselbe Asymmetrie
  besteht auf der **Ziel-Achse** gewöhnlicher In-Scope-Module: `matrix` liest den
  Status der **Ziel**-Datei über `statusOf` → `proseLines`; liegt die Zieldatei
  außerhalb `scan.roots` und steht ihre Status-Zeile hinter einem offenen Fence,
  meldet der Lauf weder `matrix-inactive` noch `fence-unclosed`.
- `verifizierbar`: ja — Fixture mit `scan.roots: ["docs"]`, `docs/a.md`
  verlinkt `../b.md`, `b.md` trägt vor der Status-Zeile einen offenen Fence:
  `0 Befund(e)`, Exit 0. Gegenprobe ohne den Fence:
  `docs/a.md:3 ../b.md matrix-inactive`. Dieselbe Bauform tragen `anchors`
  (Slug-Cache über Zieldateien), `diagrams` (`defined-in`) und `versions`
  (`current-from`).
- `klasse`: „Modul-Grenze nur auf der Quell-Achse gedacht"

### F-4 — Gezählt wird Parität, nicht Paarung: legale geschlossene Verschachtelung meldet, fremde Toggle-Quellen verschlucken

- `kategorie`: MEDIUM
- `quelle`: [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md)
  §Fitness Function („Balanciert bleibt grün: eine Datei mit vielen, jeweils
  geschlossenen Fences — auch verschiedener Länge — erzeugt keinen Befund")
- `pfad`: `internal/hexagon/core/rules/spans.go:35`
- `befund`: Die Fitness-Zusage hält nicht. Eine legale, vollständig geschlossene
  Verschachtelung — ein vier-Backtick-Block, der einen drei-Backtick-Block zeigt
  — erzeugt einen `fence-unclosed`-Befund an der **Schluss**-Zeile des äußeren
  Blocks, und der Rest der Datei bleibt trotzdem übersprungen (ein kaputter Link
  dahinter wird nicht gemeldet). Dasselbe gilt ohne jede Längen-Differenz, sobald
  Fence-Zeichen gemischt werden (Tilden-Block mit einer Backtick-Zeile darin).
  In die andere Richtung verschluckt eine **gerade** Zahl von Toggles aus
  Nicht-Fence-Quellen (zwei Backtick-Zeilen innerhalb eines 4-Space-eingerückten
  Codeblocks, zwischen denen Prosa steht) diese Prosa ohne jeden Befund.
- `verifizierbar`: ja — drei Fixtures: (i) vier-Backtick-Öffner mit Infozeile
  `markdown`, innen ein drei-Backtick-Block, Schluss aus vier Backticks →
  `a.md:6 fence-unclosed`, der Link in Z. 8 fehlt im Befundsatz; (ii) Tilden-
  Öffner, Backtick-Zeile, Tilden-Schluss → `a.md:5 fence-unclosed`, Link in Z. 7
  verschluckt; (iii) Z. 3 und Z. 7 je vier Leerzeichen plus drei Backticks, Z. 5
  kaputter Link → `0 Befund(e)`, Exit 0, Gegenprobe ohne die beiden Zeilen
  meldet `target-missing`.
- `klasse`: „Paritätszählung statt Paarung — legale Verschachtelung kippt beide
  Richtungen"

### F-5 — Zwei weitere Rückbauten bleiben grün; das zugesagte Befund-Ziel hat keine Assertion

- `kategorie`: MEDIUM
- `quelle`: [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md)
  §Fitness Function („der Test ist mutations-echt"); Reviewer-Skill §MEDIUM
  (fehlende Negativtests bei neuem öffentlichen Vertrag)
- `pfad`: `internal/hexagon/core/rules/spans_test.go:97`
- `befund`: Über die im Slice geprüften zwei Mutationen hinaus bleiben zwei
  weitere unbemerkt: die in [§`DC-FA-SPAN-001.a`](../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
  zugesagte Kappung auf 30 Zeichen von 30 auf 300 geändert (`make test` grün),
  und die Trimmung der Zeile entfernt, sodass eingerückte Fences nicht mehr
  zählen (`make test` grün). Keiner der drei neuen Tests prüft das `Target` des
  Befundes, und keiner enthält eine eingerückte Fence-Zeile; `clipRunes` hat
  keinen eigenen Test.
- `verifizierbar`: ja — beide Mutationen einzeln angewendet, `make test` je
  grün (Exit 0), Arbeitsbaum danach wiederhergestellt (`git diff` leer).
- `klasse`: „zugesagtes Befund-Feld ohne Assertion"

### F-6 — Die gemeldete Zeile ist bei fehlendem Schluss systematisch nicht die Reparaturstelle

- `kategorie`: LOW
- `quelle`: [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md)
  Entscheidung 3 („Die Öffnungszeile ist der Ort der Reparatur, also der Ort des
  Befundes")
- `pfad`: `internal/hexagon/core/rules/spans.go:47`
- `befund`: Gemeldet wird die zuletzt **öffnend gewertete** Toggle-Zeile, nicht
  die Öffnung ohne Schluss. Fehlt in einem Dokument mit mehreren Blöcken ein
  einzelner Schluss, kippt die Parität aller folgenden Fence-Zeilen, und der
  Befund steht auf der **schließenden** Zeile eines intakten späteren Blocks.
- `verifizierbar`: ja — Fixture mit Öffner Z. 3 (Schluss fehlt), Block B
  Z. 8/10, Block C Z. 12/14: Befund `a.md:14`, Ziel die Schluss-Zeile; die
  Reparatur gehört hinter Z. 5.
- `klasse`: „gemeldete Zeile ist nicht die Reparaturstelle"

### F-7 — Der Sammelsatz der Anforderung ist mit der dritten Klasse nicht mitgewachsen

- `kategorie`: LOW
- `quelle`: [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
- `pfad`: `spec/lastenheft.md` (der unveränderte Satz unmittelbar hinter der
  Klassen-Aufzählung)
- `befund`: Der Satz „Der Befund nennt Datei, Zeile (Opener- bzw. Muster-Zeile),
  die betroffene Backtick-Folge bzw. das Muster und den Grund" deckt Klasse 3
  nicht: dort ist das Ziel die gesamte getrimmte Fence-Zeile samt Infozeile, also
  weder eine Backtick-Folge noch ein Muster.
- `verifizierbar`: ja — Lauf gegen eine Datei mit langer Infozeile liefert als
  Ziel `jsonc-mit-sehr-langer-Infoz` samt der drei führenden Backticks.
- `klasse`: „Sammelsatz einer Anforderung nicht mit der neuen Klasse mitgewachsen"

## Negativbefunde

- geprüft, ohne Befund: **Bestandsmessung §3 des Slice-Plans.** Methode
  unabhängig nachgebaut (Nachbau von `FenceToggle` inklusive der
  CommonMark-Infozeilen-Regel und der `TrimLeft`-Behandlung) und über alle drei
  Repos gefahren. Die Denominatoren sind exakt reproduzierbar als „alle
  `*.md` ohne die beiden vendorten Baseline-/Cache-Bäume": 347 (eigenes Repo,
  heute 348 wegen des neuen ADR) + 184 + 245 = **776**. Null ungerade
  Toggle-Zahlen, null Tilden-Fences, null Fence-Längen ungleich drei — auch über
  die Obermenge inklusive der vendorten Bäume (882 Dateien). Die Messung misst
  tatsächlich d-checks eigene Lexik und ist korrekt.
- geprüft, ohne Befund: **Behauptung „(b) löst den Fall nicht".** Sie hält: ein
  Fence ohne jede weitere Fence-Zeile bleibt unter jeder Paarungsregel offen
  (Fixture belegt). Einschränkend gehört zum Bild, dass der längenabgeglichene
  Schluss im Produkt bereits existiert (`internal/hexagon/core/app/trace_table.go`)
  — die Aussage „wirkungslos" gilt für den heutigen Dokumentbestand, nicht für
  die Produkt-Interna (siehe F-1).
- geprüft, ohne Befund: **„genau ein Befund je Datei".** Fixture mit drei
  Öffnern in einer Datei und einem zweiten Dokument: je genau ein Befund, zwei
  insgesamt.
- geprüft, ohne Befund: **„dateiweit statt absatzweise".** Der Automat läuft über
  `splitLines(content)` ohne Absatzgruppierung; Befund entsteht am Dateiende.
- geprüft, ohne Befund: **„Ziel auf 30 gekappt".** Die Kappung ist runen-basiert
  (`clipRunes`) und schneidet Umlaute nicht mitten im Zeichen; gemessen exakt
  30 Runen.
- geprüft, ohne Befund: **Zeilenenden und Datei-Ränder.** CRLF (Befund korrekt,
  Ziel ohne `\r`), Datei ohne abschließenden Zeilenumbruch, Fence in der
  allerletzten Zeile, Tilden-Öffner, Infozeile mit Backtick (kein Öffner, kein
  Befund) — alle wie zugesagt.
- geprüft, ohne Befund: **Befund-Deduplikation.** `model.SortFindings`
  dedupliziert über das Tupel (Datei, Zeile, Regel, Ziel, Grund) inklusive
  Grund; eine Kollision von `fence-unclosed` mit `span-unclosed` auf derselben
  Zeile ist ausgeschlossen.
- geprüft, ohne Befund: **`--doctor` und `--repair`.** Der Klartext ist
  vorhanden und wird gruppiert ausgegeben; `--repair-broad` liefert
  erwartungsgemäß keinen Hunk (`FixCandidateFor` bedient nur `id-unlinked`) und
  bricht nicht.
- geprüft, ohne Befund: **Grund-Code-Lockstep.** Neuer Code in `AllReasons`,
  in `reasonTexts` und in der §4-Tabelle der Spezifikation; `make test` grün.
- geprüft, ohne Befund: **Blockquote-Fences.** Eine Fence-Zeile hinter `>` ist
  in **beiden** Automaten kein Toggle; kein Befund, keine Divergenz.
- geprüft, ohne Befund: **Import-Regeln (ADR-0005).**
  `internal/hexagon/core/rules/spans.go` importiert nur `strings` und das
  `model`-Paket; `make lint` grün.
- geprüft, ohne Befund: **eigener Bestand.** `make doc-check` über das Repo:
  348 Dateien, 0 Befunde, Exit 0 — die Zusage „der eigene Bestand bleibt bei
  null" hält.
- geprüft, ohne Befund: **Referenz-Richtung (SDP).** Der Diff trägt keinen
  Provenance-Marker und keinen neuen Abwärts-Verweis aus ADR oder Spec auf
  Planning-Artefakte.
- geprüft, ohne Befund: **Lifecycle-/Ablage-Konventionen des Diffs.**
  ADR im Index ergänzt, Lastenheft-Version und beide Historien-Zeilen gesetzt,
  Slice-Datei im `in-progress`-Verzeichnis.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Wächter misst einen anderen Automaten als der
bewachte Verbraucher · Wächter und Bewachtes trimmen verschieden · Modul-Grenze
nur auf der Quell-Achse gedacht · Paritätszählung statt Paarung — legale
Verschachtelung kippt beide Richtungen · zugesagtes Befund-Feld ohne Assertion ·
gemeldete Zeile ist nicht die Reparaturstelle · Sammelsatz einer Anforderung
nicht mit der neuen Klasse mitgewachsen

## Verdikt

**Merge-blockierend:** ja — zwei HIGH und drei MEDIUM.

Die Richtung des Change ist richtig und der belegte Reproduktionsfall meldet
zuverlässig; der Wächter schließt den Fall, für den er gebaut wurde. Blockierend
ist nicht dieser Kern, sondern die **Reichweite**, die Anforderung und ADR ihm
zuschreiben: „jede Vorverarbeitung", „alle Module", „jedes Modul profitiert".
F-1 zeigt einen weiterhin stillen Grün-Pfad in einem Closure-Gate, F-2 einen im
neuen Code selbst — beide mit Lauf belegt, beide innerhalb der zugesagten
Klasse. Damit steht die Zusage höher als das Verhalten, und das ist bei einem
Wächter gegen Silent-Grün genau die Eigenschaft, die er beseitigen soll.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist
ein Lauf-Beleg (dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und
ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
