# Review slice-172 — Der offene DoD-Haken eines geschlossenen Slice wird gewächtert


**Review-Art:** Code/Design (gegen Slice-Plan, ADR, Hard Rules)
**Gegenstand:** `fd4cbb4^..HEAD` = `fd4cbb4` (Prep) · `e8401c7` (Claim-Move) · `79bf375` (feat)
**Skill:** `.harness/skills/reviewer.md` @ v1.13.0 · **Modell:** claude-opus-5[1m] · **Datum:** 2026-08-30
**Eingangs-Kontext:** `AGENTS.md` §2/§3.1/§3.2/§3.4/§3.6/§3.7/§3.8/§4/§5, `harness/conventions.md` (MR-013, MR-049, MR-053, MR-054, MR-056), `.harness/baseline/v5.12.0/regelwerk/` (`grundlagen-harness-dateien.md` §Was ein Kommentar trägt, `modul-05-planning-harness.md`, `modul-06-roadmap.md`, `modul-10-review-harness.md`), DC-FA-STRUCT-001, DC-FA-PLAN-001, DC-FA-SPAN-001, ADR-0059, ADR-0073, ADR-0074, ADR-0077, Beobachtungs-Register (BEO-011, BEO-012, BEO-013, BEO-015, BEO-020, BEO-022, BEO-023, BEO-024), `docs/plan/planning/welle-86-closure-uebergang-durchsetzen.md`, Slice-Plan §2 (die Vorarbeit aus drei Anläufen)

### Eigener Lauf

| Lauf | Ausgabe |
|---|---|
| `make verify-closure-notes` (HEAD, unverändert) | `d-check: 560 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `make doc-check` (HEAD, unverändert) | `d-check: 623 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| Gegenprobe **ohne** Bestands-Ausnahme, echtes Closure-Profil | `d-check: 560 Datei(en) geprüft, 144 Befund(e)` · Exit 1 |
| Positiv-Probe an `slice-181` (Haken geöffnet, echtes Profil) | 1 Befund, `:92`, `section-tasks-open`, Exit 1; zurückgesetzt wieder 0 |
| Fence-Probe isoliert (offener Haken hinter vergessenem Schluss-Fence) | ohne `spans` `0 Befund(e)` / Exit 0 · mit `spans` `fence-unclosed` / Exit 1 |
| Fence-Probe **im echten Repo** (Fence in §4 von `slice-181`) | `section-empty` + `fence-unclosed`; die neue Regel selbst **stumm** |
| 4 Gegenproben an synthetischen `done/`-Bäumen (`--network none`, read-only) | s. M-1, M-2 |
| Glob-Probe über 16 synthetische Slice-Nummern | s. Negativbefunde |
| Mutations-Probe an der `cite`-Direktive in MR-056 | `citation-mismatch` · der Beleg ist gate-gebunden |
| Bestands-Zensus `done/` (`ls`, `grep`) | 177 Slice-Dateien · 7 DoD-Überschriften-Formen · **0** H3-Formen |

`make gates`, `make fullbuild`, `make record-gates`, `make image-scan` bewusst **nicht**
gefahren — Gate-Lauf-Bestätigung ist Verifier-Rolle, nicht Reviewer-Rolle. Arbeitsbaum
nach allen Proben sauber (`git status --short` leer bis auf diesen Report).

---

## Urteil

**BLOCKIERT** — 2 HIGH · 3 MEDIUM · 3 LOW · 3 INFO.

Der dritte Bau ist der beste der drei. **Die tragenden Zahlen stimmen alle**: 144 Befunde
in 37 Dateien ohne Ausnahme, null mit ihr, dieselben Nummern (025–104, 160, 168, 169,
170), die Positiv-Probe fängt an der genannten Stelle, und die Fence-Abhängigkeit ist
nicht behauptet, sondern reproduzierbar — beide Hälften auf die Zeile genau. Die
Ausnahme-Globs halten den Maßstab von MR-049 nachweislich ein: vier- und zweistellige
Nummern laufen unter die Regel, die lautlose Klasse ist vermieden.

Blockierend sind zwei Dinge:

1. **Die eine Zahl, die niemand nachgezählt hat, ist die Belegzahl.** *„Elf Closures ohne
   einen einzigen offenen Haken"* stützt die Aussage, dass die korrigierte Praxis hält.
   Gemessen liegen in `done/` **acht** Slices oberhalb der Ausnahme; 173, 174 und 175
   existieren nirgends im Repo, und das steht in der Welle-Datei, auf die derselbe Plan
   verlinkt. Die Zahl ist aus der Nummernspanne gerechnet, nicht aus dem Verzeichnis
   gelesen — die BEO-020-Gestalt (Register-Zähler 5, jüngster Beleg slice-182).
2. **Der Kommentar erzählt die Geschichte der Datei.** Der neue Block referiert, was
   *früher hier notiert* war, und widerlegt es. Das ist Review-/Edit-Historie im
   Kommentar (§3.7), und derselbe Absatz steht bereits — richtig — in der
   Commit-Botschaft.

Die drei MEDIUM sind: ein `hint`, der auf einem Befund erscheint, den er nicht meint
(gemessen, im echten Profil); eine dreifach zugestellte **Zwei-Grenzen**-Zusage, die eine
dritte, gemessene Blindstelle nicht nennt; und eine Vorprüfung, deren Register-Stand bei
der dritten Beanspruchung nicht nachgezogen wurde, obwohl der Nachtlauf-Block am selben
Tag nachgezogen wurde.

---

## Nachtrag — Stand des Arbeitsbaums bei der Übergabe

**Dieser Report prüft die drei Commits, nicht den Arbeitsbaum.** Während er entstand,
hat ein paralleler Kontext den Slice-Plan, `harness/conventions.md` und MR-056
**uncommittet** geändert. Gemessen am Übergabe-Zeitpunkt:

| Befund | Stand im Arbeitsbaum |
|---|---|
| H-1 | **adressiert** — `grep -rn` auf die Wendung über `docs/plan` und `harness/` liefert keine Fundstelle mehr; Plan und MR-056 nennen jetzt **acht** Dateien und benennen die Lücken 172–175 |
| H-2 | **offen** — `.d-check.closure.yml` ist unverändert (`git status` führt sie nicht) |
| M-1 | **offen** — der `hint` ist unverändert |
| M-2 | **offen — und die laufende Bearbeitung verschärft ihn.** `AGENTS.md` sagt weiterhin zweimal *„Zwei Grenzen"*. `harness/README.md` trägt neu: *„**Die vier Grenzen unten gelten ihr NICHT:** sie liest die **rohen** Zeilen … ihre eigene Grenze ist der **vergessene Schluss-Fence**"*. Gemessen ist das für eine der vier falsch: **Fenced Blocks** bleiben auch für diese Bedingung außen vor (`prose` überspringt sie); ein Haken in einem **wohlgeformten** Fence im DoD-Abschnitt liefert `0 Befunde, Exit 0` und **kein** `fence-unclosed`. Die Nicht-Vererbung gilt der Inline-Code-Paarung, nicht dem Fence |
| M-3 | **adressiert** — der Sichtungs-Block trägt jetzt Register-Stand 2026-08-30, höchste Kennung `BEO-024`, und benennt das Versäumnis |
| L-1 | **adressiert** — `harness/conventions.md:136` hat wieder fünf Pipes wie der Tabellenkopf |
| L-2 | **offen** |
| L-3 | **adressiert** — die Teilsumme ist auf 33 korrigiert |
| I-1 · I-2 · I-3 | unverändert |

Die Findings bleiben stehen: sie beschreiben den Stand der Commits, gegen die geprüft
wurde. Die Spalte oben ist eine **Messung des Arbeitsbaums**, keine Nachprüfung der
Korrekturen — ob sie tragen, gehört in den nächsten Lauf.


## Findings

### H-1 · „Elf Closures" ist aus dem Verzeichnis nicht belegbar — gemessen sind es acht

- **quelle:** `AGENTS.md` §5 (*ein Schluss reicht nicht weiter als die gemessene Menge*), BEO-020, Reviewer-Skill Prüffrage 8 / HIGH-nahe Hälfte von BEO-009
- **pfad:** `docs/plan/planning/in-progress/slice-172-closure-uebergang-waechtern.md:118-127` (§2); `harness/conventions/MR-056-dod-haken-waechter.md:44-45`; Commit-Botschaft `fd4cbb4`
- **zugesagt:** wörtlich in MR-056: *„Gemessen halten die Slices 171–182 ihre Haken
  gesetzt — elf Closures ohne einen einzigen offenen."* Im Plan §2 mit dem Zusatz:
  *„**Das ist zugleich der Beleg, dass die korrigierte Praxis hält** — elf Closures ohne
  einen einzigen offenen Haken."*
- **gemessen:** `ls docs/plan/planning/done/slice-17*.md docs/plan/planning/done/slice-18*.md`

  ```
  done/slice-170-workflows-modul.md
  done/slice-171-vorpruefungen-belegen.md
  done/slice-176-planning-rule-pilot.md
  done/slice-177-structure-hint.md
  done/slice-178-offene-tasks-roh.md
  done/slice-179-strukture-teilmenge.md
  done/slice-180-closure-profil-spans.md
  done/slice-181-handbuch-link-in-cli.md
  done/slice-182-erklaerte-leermenge.md
  ```

  Oberhalb der Ausnahme (also ab 171) liegen **acht** Dateien, nicht elf.
  `find . -name "slice-17[345]*" -not -path "./.git/*"` liefert **nichts** — 173, 174 und
  175 existieren in keinem Lifecycle-Verzeichnis. Gegenzählung über die Globs:
  `slice-0??-*` 99 · `slice-1[0-6]?-*` 69 · `slice-170-*` 1 = **169** ausgenommen von
  **177**, gemessen also **8**.

  Die Zahl 11 entsteht nur, wenn man die Spanne 171–182 als lückenlos annimmt und
  slice-172 abzieht. Genau diese Annahme widerlegt
  `docs/plan/planning/welle-86-closure-uebergang-durchsetzen.md:92`, auf das derselbe
  Plan-Absatz verlinkt: *„Die drei letzten sind **noch nicht angelegt**."*
- **warum das trägt:** Die Aussage ist nicht Beiwerk, sondern der einzige vorgelegte
  Beleg dafür, dass die Grenze bei 171 richtig liegt und nicht bloß bequem ist. Aus acht
  Closures statt elf folgt nichts anderes — aber die Belegdichte ist um 27 % kleiner als
  behauptet, und sie ist behauptet, nicht gezählt. Dieselbe Gestalt trägt BEO-020 zuletzt
  bei slice-182; der Register-Zähler steht bei 5.
- **verifizierbar:** ja — `ls docs/plan/planning/done/slice-1[78]*.md | wc -l`
- **klasse:** `belegzahl-aus-der-spanne-gerechnet`

### H-2 · Der neue Konfigurations-Kommentar referiert seine eigene Vorgänger-Fassung

- **quelle:** `AGENTS.md` §3.7 (Kommentar trägt eine der fünf Klassen; *keine Review-Historie, keine Deliberation über Verworfenes*), Baseline `grundlagen-harness-dateien.md` §Was ein Kommentar trägt; Reviewer-Skill Prüffrage 6 (HIGH)
- **pfad:** `.d-check.closure.yml:197-201`
- **befund:** Der Block enthält:

  ```
  # Die frueher hier notierte Ablehnung ("32 Befunde, und die sind berechtigt,
  # weil die WELLE die Box einloest") ist ueberholt: was eine Welle einloest,
  # gehoert in ihren Closure-Trigger, nicht in die DoD eines Slice
  # (AGENTS.md §5). Ein DoD-Punkt, den der Slice selbst nicht abhaken kann,
  # macht den Haken als Zustandsfeld unbrauchbar.
  ```

  Das Subjekt der ersten anderthalb Zeilen ist ein Text, den derselbe Commit **entfernt**
  hat (`.d-check.closure.yml:104-108` im Vor-Stand). Ein Kommentar beschreibt, was da ist; hier
  beschreibt er, was nicht mehr da ist, und dass es widerlegt wurde. Das ist die
  Edit-Historie der Datei — die Klasse, die §3.7 ausdrücklich ausschließt.
- **abgegrenzt, damit der Befund nicht zu weit greift:** Die **zweite** Hälfte desselben
  Absatzes (*„was eine Welle einlöst, gehört in ihren Closure-Trigger … (AGENTS.md §5)"*)
  ist ein sauberer Rang-Zeiger plus Abgrenzung und trägt für sich. Ebenso der Absatz
  `:187-191` (`max-open-tasks` statt `forbid-pattern`): er warnt einen künftigen Editor
  vor einer Vereinfachung und ist damit **Abgrenzung**, keine Deliberation. Und die
  Mess-Zahlen `:203-210` folgen der etablierten Form derselben Datei (MR-049-Block
  `:167-173`) — kein Befund.
- **verifizierbar:** nein (Urteil; kein Gate fängt §3.7). Beobachtbar ist, dass derselbe
  Absatz bereits in der Commit-Botschaft `79bf375` steht — dem Ort, an den ihn die
  Hausregel verweist.
- **klasse:** `kommentar-erzaehlt-die-datei-historie`

### M-1 · Der `hint` erscheint auf `section-missing` — dort meint er den falschen Defekt

- **quelle:** ADR-0073 (*„ein Befund ohne Bedingungs-Verletzung behält seine eigene Meldung"*), `AGENTS.md` §3.8, Reviewer-Skill Prüffrage 2/15
- **pfad:** `.d-check.closure.yml:219` (`hint`); `internal/hexagon/core/rules/structure.go:126-128` (`section-missing` läuft durch `structureFinding`, also durch `MessageFor`)
- **gemessen, im echten Closure-Profil.** Probe: ein vergessener Schluss-Fence in Zeile 1
  von `slice-181` (danach zurückgesetzt) — die fence-bewusste Überschriften-Suche findet
  keine Überschrift mehr, und vier Regeln melden auf derselben Datei:

  ```
  slice-181…:1  :: ## 5. Abnahme-Punkte / Risiken            section-missing  kein Abschnitt passt auf den Selektor
  slice-181…:1  :: ^## [0-9]+\. Definition of Done           section-missing  kein Abschnitt passt auf den Selektor
  slice-181…:1  :: ^#{2,3} .*Closure-Notiz                   section-missing  kein Abschnitt passt auf den Selektor
  slice-181…:1  :: ^#{2,3} [0-9]+\. Definition of Done       section-missing  offener DoD-Haken in einem geschlossenen Slice: Haken setzen, wenn der Punkt erledigt ist — sonst gehört der Slice zurück nach in-progress/
  slice-181…:3  ```text                                      fence-unclosed
  ```

  Drei der vier sagen die Wahrheit. Die neue Regel sagt *„offener DoD-Haken … Haken
  setzen"* — und es gibt keinen offenen Haken; es gibt einen kaputten Fence. Zweite Probe
  an einem synthetischen Baum (DoD-Abschnitt unter anderem Titel, sonst intakt): derselbe
  Text auf `section-missing`, Exit 1.
- **gegen die Quelle gehalten:** ADR-0073 nimmt **zwei** Befunde vom `hint` aus — den
  unlesbaren Dateibaum und die leer laufende Regel — mit der Begründung, die Regel habe
  *„dort gar nicht gemessen"*. `section-missing` (Selektor trifft nichts) und
  `section-ambiguous` sind derselbe Fall und sind **nicht** ausgenommen. Die Ausnahme ist
  im Produkt also schmaler als in der ADR, die sie deklariert.
- **wer es trägt:** Die Ausnahme-Liste ist Produkt (slice-177). Sichtbar wird sie erst
  hier, weil slice-172 der **erste** Konsument von `hint` ist und einen Hinweis wählt,
  der ausschließlich die Task-Bedingung beschreibt. Der Befund gehört deshalb an diesen
  Diff, nicht in die Vergangenheit.
- **verifizierbar:** ja — `make verify-closure-notes` gegen eine `done/`-Datei mit
  umbenanntem DoD-Abschnitt
- **klasse:** `hint-auf-befund-ohne-bedingungs-verletzung`

### M-2 · „Meldet **jeden** offenen Haken" mit „**Zwei** Grenzen" — gemessen gibt es eine dritte, stille

- **quelle:** `AGENTS.md` §3.8 (*ein Modul verspricht nur über das, was es scannt*), BEO-023, Reviewer-Skill Prüffrage 15
- **pfad:** `AGENTS.md:386` (§4-Zeile) · `AGENTS.md:471-482` (§5) · `harness/README.md:98` (Sensors) · `harness/conventions/MR-056-dod-haken-waechter.md:54-70` (Grenzen-Liste)
- **zugesagt:** §4: *„meldet jeden offenen Haken eines `done/`-Slice … **Zwei Grenzen:**
  ein Haken ist eine Selbstauskunft, und ein vergessener Schluss-Fence schaltet die
  Bedingung ab"*. §5 wortgleich in der Sache (*„**Zwei Grenzen gehören dazu**"*). MR-056
  führt vier Grenzen — Selbstauskunft, vergessener Fence, Zustand-statt-Übergang,
  alternde Ausnahme —, aber keine davon deckt den gemessenen Fall.
- **gemessen:** synthetischer `done/`-Baum, identische Regel, ein Haken in einem
  **wohlgeformten** Fenced Block innerhalb des DoD-Abschnitts:

  ```
  ## 4. Definition of Done

  <fence markdown>
  - [ ] offener Haken IM geschlossenen Fence
  <fence>

  - [x] echter Punkt.
  ```

  Ergebnis: `d-check: 2 Datei(en) geprüft, 1 Befund(e)` — und der eine Befund stammt aus
  der Nachbardatei. Für diese Datei: **0 Befunde, Exit 0, kein `fence-unclosed`.** `spans`
  fängt hier nichts, weil nichts kaputt ist: der Fence ist korrekt.
- **die Grenze ist nicht neu, ihre Zustellung fehlt.** ADR-0074 §Entscheidung 5 (`docs/plan/adr/0074-offene-tasks-auf-rohen-zeilen.md:91-95`) nennt sie
  (*„Der Fence bleibt außen vor"*), und der Code-Kommentar
  `internal/hexagon/core/rules/structure.go:349-352` ebenfalls. Keiner der **vier**
  Zustell-Orte, die dieser Slice bespielt, nennt sie — und drei von ihnen behaupten
  daneben eine geschlossene Zahl (*zwei*) und eine Allaussage (*jeden*). Wer die Grenze
  sucht, findet sie nur im Produkt-Quelltext, den §3.1 dem Leser gar nicht zumutet.
- **Failure-Szenario:** ein `done/`-Slice ab 171, dessen DoD-Punkt als Beispiel-Block
  formatiert ist — der Wächter meldet grün, die Zusage *„kein abgeschlossener Slice trägt
  einen offenen DoD-Haken"* ist trotzdem verletzt. Das ist BEO-023 in seiner
  Nicht-Test-Ausprägung, die die Registerzeile ausdrücklich mitführt.
- **verifizierbar:** ja — die Probe oben gegen `--enable structure --enable spans`
- **klasse:** `zwei-grenzen-zugesagt-drei-gemessen`

### M-3 · Die Beobachtungs-Sichtung ist bei der dritten Beanspruchung nicht nachgezogen worden

- **quelle:** Baseline `modul-05-planning-harness.md` §Zwei Schritte vor der Modus-Begründung (*„Gelesen wird der **gemergte** Stand"*), `AGENTS.md` §5 (drei Vorprüfungen), MR-054, BEO-024
- **pfad:** `docs/plan/planning/in-progress/slice-172-closure-uebergang-waechtern.md:342-346` (§7, zweiter Block)
- **befund:** Der Block deklariert *„Register-Stand 2026-08-29, höchste Kennung
  `BEO-023`"* und sichtet fünf Einträge. Der Prep-Commit `fd4cbb4` vom **2026-08-30**
  hat den **dritten** Block (Nachtlauf) ausdrücklich nachgezogen — *„bei der dritten
  Beanspruchung neu gelesen"*, neue Zeitstempel — und den zweiten unverändert gelassen.
- **gemessen:** `git show fd4cbb4^:docs/plan/planning/observations.md | grep -c "BEO-024"`
  → `1`. BEO-024 stand also **vor** dem Prep-Commit im Register (eingetragen
  2026-08-29 durch `b35662e`, slice-176). Höchste Kennung zum Zeitpunkt der dritten
  Beanspruchung ist BEO-024, nicht BEO-023.
- **und der ausgelassene Eintrag ist einschlägig:** BEO-024 lautet *„Ein Zustell-Kanal
  hängt an der ARBEITSWEISE, die Regel aber am Inhalt"*, Sub-Area `*`, und die
  Registerzeile weitet die Klasse selbst aus: *„jede Bedingung, die nach dem WAS klingt
  und am WIE hängt … ein Gate, das nur bei bestimmten Commit-Formen läuft"*. Genau das
  ist dieser Wächter: er hängt nicht am `mv`-Commit, sondern daran, dass jemand
  `make verify-closure-notes` fährt. MR-056 nennt die Sache (*„Sie prüft den Zustand,
  nicht den Übergang"*) — aber die Vorprüfung, die die Klasse benannt hätte, hat nicht
  stattgefunden, und der Block behauptet, sie habe.
- **Failure-Szenario:** die Vorprüfung ist der einzige Leser für alles unter der
  Schwelle. Ein Block, der einen überholten Register-Stand deklariert, ist genau die
  Selbstauskunft, gegen die MR-054 gebaut ist — nur eine Ebene höher: die Direktive
  belegt, dass die *Regelwerk*-Zeile gelesen wurde, nicht dass das *Register* gelesen
  wurde.
- **verifizierbar:** ja — `git show fd4cbb4^:docs/plan/planning/observations.md | grep BEO-024`
- **klasse:** `vorpruefung-mit-veraltetem-register-stand`

### L-1 · Die MR-056-Zeile in `harness/conventions.md` trägt eine Spalte zu viel

- **quelle:** Maintainability (Tabellenform des MR-Index)
- **pfad:** `harness/conventions.md:136`
- **gemessen:** Pipe-Zählung über den Tabellenblock:

  ```
  102 pipes=5 :: | MR | Titel | Geltungsbereich | Ersetzt-Baseline-Regel |
  ...
  135 pipes=5 :: | [MR-055](...)
  136 pipes=6 :: | [MR-056](...)
  ```

  Die Zeile endet auf `… keine Abweichung | |` — eine fünfte, leere Zelle über einem
  vierspaltigen Kopf. Kein Gate fängt es (`make doc-check` ist grün); GFM verwirft die
  Überzahl beim Rendern, `structure.table` läuft über diese Datei nicht.
- **verifizierbar:** ja — `awk '{n=gsub(/\|/,"|"); print NR, n}' harness/conventions.md`
- **klasse:** `tabellenzeile-mit-ueberzaehliger-zelle`

### L-2 · `^#{2,3}` deckt eine Überschriften-Ebene, die es im Bestand nicht gibt — und weicht von der Nachbarregel ab

- **quelle:** Reviewer-Skill Prüffrage 11 (dieselbe Eingabe-Klasse, verschieden behandelt, ohne benannten Grund)
- **pfad:** `.d-check.closure.yml:216` gegen `.d-check.closure.yml:120`
- **befund:** Acht Zeilen über der neuen Regel selektiert dieselbe Eingabe-Klasse
  (`done/slice-*.md`, DoD-Abschnitt) mit `^## [0-9]+\. Definition of Done` und der
  Kardinalität `one` (Default). Die neue Regel nimmt `^#{2,3} …` und `sections: each`.
  Beide Abweichungen sind unbegründet gelassen.
- **gemessen:** über alle 177 `done/`-Slices:

  ```
  H3-DoD-Ueberschriften: (keine)
  H2-DoD gesamt: 177   Dateien gesamt: 177
       33 ## 2. Definition of Done
       45 ## 3. Definition of Done
        1 ## 3. Definition of Done (R1 eingearbeitet)
        4 ## 3. Definition of Done (vorläufig)
       88 ## 4. Definition of Done
        3 ## 4. Definition of Done (vorläufig)
        3 ## 5. Definition of Done
  ```

  Die vom Plan genannten **sieben Formen** stimmen exakt (33+45+1+4+88+3+3 = 177). Die
  H3-Öffnung deckt **null** Realfälle, und `each` unterscheidet sich heute von `one`
  nicht, weil jede Datei genau eine DoD-Überschrift trägt. Der Preis ist, dass eine
  künftige H3-Wiederholung (etwa eine zitierte Vorlage) von der einen Regel gemessen und
  von der anderen als `section-ambiguous` gemeldet würde — zwei Antworten auf dieselbe
  Datei.
- **verifizierbar:** ja — `grep -c "^### [0-9]\+\. Definition of Done" docs/plan/planning/done/slice-*.md`
- **klasse:** `selektor-weiter-als-der-nachbar-ohne-grund`

### L-3 · „025–104 (34 Stück)" — gemessen sind es 33, und 34+4 widerspricht der 37 im selben Absatz

- **quelle:** `AGENTS.md` §5, BEO-020
- **pfad:** `docs/plan/planning/in-progress/slice-172-closure-uebergang-waechtern.md:133-138`
- **gemessen:** die 37 Treffer-Dateien, nach Nummern sortiert:

  ```
  025 026 027 028 040 041 042 043 044 045 046 048 049 050 051 052 053 054 055 057
  064 065 066 067 068 076 082 094 097 098 101 102 104 | 160 168 169 170
  ```

  Im Bereich 025–104 liegen **33** Dateien, nicht 34; 33 + 4 = 37, wie überall sonst
  behauptet. Mit der 34 ergäbe derselbe Satz 38.
- **abgegrenzt:** Die Zeile ist Bestand aus dem zweiten Anlauf und vom Prep-Commit nicht
  angefasst worden. Gemeldet wird sie trotzdem, weil sie mit der dritten Beanspruchung
  ausdrücklich als weiter gültig erklärt wird (*„Die drei Ausnahme-Muster aus der vorigen
  Zeile tragen unverändert"*) und weil sie sich gegen die tragende 37 stellt.
- **verifizierbar:** ja — Datei-Liste des Laufs ohne Ausnahme
- **klasse:** `teilsumme-widerspricht-der-gesamtsumme`

### I-1 · Die Ausnahme ist 4½-mal so groß wie ihr Anlass — und das ist keine der vier genannten Grenzen

- **quelle:** MR-056 §Grenzen, BEO-013
- **pfad:** `.d-check.closure.yml:211-215`; `harness/conventions/MR-056-dod-haken-waechter.md:68-70`
- **gemessen:** 169 der 177 `done/`-Slices sind ausgenommen; **37** von ihnen tragen
  einen offenen Haken. Die übrigen **132** sind heute sauber und trotzdem dauerhaft
  stumm: wird eine dieser Dateien später bearbeitet und ein Haken geöffnet, meldet
  nichts.
- **warum nur INFO:** Das ist der Preis der Datums-Grenze und in MR-049 identisch
  vorhanden; eine treffer-genaue Ausnahme wäre eine Liste, die schlechter altert. Die
  vierte MR-056-Grenze (*„Die Ausnahme altert"*) nennt allerdings nur das **Wachsen** der
  Ausnahme als Befund, nicht diese 132-Datei-Zone.

### I-2 · `?` im Ausnahme-Glob ist „ein Zeichen", nicht „eine Ziffer"

- **pfad:** `.d-check.closure.yml:212-215`
- **befund:** `slice-0??-*` und `slice-1[0-6]?-*` nähmen auch `slice-0ab-…` oder
  `slice-16z-…` heraus. Bei rein numerischer Slice-Benennung folgenlos; MR-049 trägt
  dieselbe Eigenschaft. Notiert, damit die nächste Instanz der Form sie nicht neu
  entdeckt.

### I-3 · §1 rechnet noch mit 169 `done/`-Slices

- **pfad:** `docs/plan/planning/in-progress/slice-172-closure-uebergang-waechtern.md:45-49`
- **befund:** *„von 169 `done/`-Slices tragen **37** mindestens einen offenen Haken"* —
  heute sind es 177 Dateien. Die 37 stimmt unverändert; nur der Nenner ist von zwei
  Rückführungen überholt. Der Satz ist als Bestandsaufnahme des ersten Anlaufs lesbar,
  trägt aber kein Datum.

---

## Negativbefunde (geprüft, ohne Befund)

- **Die beidseitige Messung reproduziert exakt.** Ohne Bestands-Ausnahme, im echten
  Closure-Profil (Zeilen 212–215 entfernt, sonst byte-gleich):
  `d-check: 560 Datei(en) geprüft, 144 Befund(e)`, Exit 1 — **alle 144** tragen den Code
  `section-tasks-open`, verteilt auf **37** Dateien; höchste Dichte `slice-054` mit 12.
  Mit Ausnahme: `make verify-closure-notes` → `560 Datei(en), 0 Befund(e)`, Exit 0. Die
  Nummern-Aufzählung im Kommentar (*„025-104, dann 160, 168, 169, 170"*) stimmt Zeichen
  für Zeichen.
- **Die Positiv-Probe fängt, wie beschrieben.** Erster Haken im §4 von `slice-181`
  geöffnet, echtes Profil: `560 Datei(en) geprüft, 1 Befund(e)`, `slice-181…:92`,
  `section-tasks-open`, Hinweistext vollständig, Exit 1. Nach `git checkout --` wieder 0.
  Datei, Zeile und Hinweis — alle drei wie zugesagt.
- **Die `spans`-Abhängigkeit ist echt, in beiden Richtungen.** Isolierter Baum, ein
  offener Haken hinter einem vergessenen Schluss-Fence im DoD-Abschnitt:
  `--enable structure` → `1 Datei(en) geprüft, 0 Befund(e)`, Exit 0 — vollständig blind.
  `--enable structure --enable spans` → `1 Befund(e)`, `fence-unclosed`, Exit 1. Im
  echten Repo (Fence in §4 von `slice-181`) meldet zusätzlich die Nachbarregel
  `section-empty`; die neue Bedingung selbst bleibt stumm. Die Botschaft sagt genau das
  und nennt das Nachbar-Signal ausdrücklich *„Zufall, nicht Konstruktion"* — kein
  überdehnter Schluss. Und `spans` deckt `done/` im Closure-Profil tatsächlich ab; die
  Abhängigkeit ist also nicht nur behauptet, sondern wirksam.
- **Die Ausnahme-Globs halten den MR-049-Maßstab — gemessen, nicht abgeleitet.** 16
  synthetische `done/`-Dateien, je ein offener Haken, dieselbe Regel:

  | Nummer | Ergebnis | Nummer | Ergebnis |
  |---|---|---|---|
  | 025 · 099 · 100 · 104 · 160 · 169 · 170 | ausgenommen (stumm) | 171 · 182 | gemeldet |
  | 999 · 1000 · 1099 · 1600 · 1700 | gemeldet | 17 · 1a5 | gemeldet |

  Die lautlose Klasse, die MR-049 benennt (`slice-1[0-3]*` verschluckt 1000–1399), ist
  vermieden: alle vierstelligen Nummern laufen unter die Regel. Der abschließende
  Bindestrich im Muster ist es, der die Ziffernzahl fixiert.
- **Die Nullmengen-Härte ist intakt.** `checkStructureRule` meldet
  `„Regel trifft keine Datei (auch nach Abzug von exempt-paths) — das Gate liefe leer"`;
  eine Ausnahme, die eines Tages alles herausnimmt, schaltet die Regel nicht still ab.
- **Die `cite`-Direktive in MR-056 ist ein gate-gebundener Beleg, kein Zierrat.** Probe:
  zwei Umlaute im Zitat transliteriert → `MR-056…:5 … citation-mismatch  Zitattext ist
  kein zusammenhängender Teilstring der Quell-Spanne (Zitat-Fäule)`. Und die Spanne
  `modul-05-planning-harness.md:33-34` ist die **vorschreibende** Zeile, nicht eine
  Nebenregel. MR-054 regiert diese Direktive ohnehin nicht — sein Geltungsbereich sind
  die zwei Vorprüfungs-Blöcke des Slice-Plans, und er erklärt eine Direktive *„anderswo
  in einem Planungs-Dokument"* ausdrücklich für erlaubt-aber-nicht-gefordert. Auch die
  zwei Direktiven im Plan selbst (`:213-214`, `:219`) ankern korrekt auf
  *„Sub-Area-Wahl prüfen"* bzw. *„Offene Beobachtungen sichten."*
- **Die übrigen Berufungen tragen.** MR-049 wird als **Form**-Präzedenz zitiert und nennt
  seinen eigenen Gegenstand dabei mit (*„wie MR-049 es für die Drei-Ausgänge-Regel
  führt"*) — keine Reichweiten-Übertragung auf §5. ADR-0074 trägt das rohe Zählen,
  ADR-0077 das `spans` im Profil, ADR-0073 die Arbeitsteilung Code/Hinweis — alle drei
  sagen an der zitierten Stelle, was von ihnen behauptet wird. Die Baseline-Zeile
  33–34 trägt die Übergangs-Bedingung wörtlich. BEO-013 wird so verwendet, wie das
  Register es selbst glossiert (BEO-022-Zeile: *„dort fängt ein vorhandener Wächter
  nichts mehr"*), BEO-023 innerhalb seiner ausdrücklich mitgeführten Nicht-Test-Hälfte.
- **§3.6 ist nicht berührt.** Eine Bestands-Ausnahme an einer **neu eingeführten** Regel
  senkt keine Schwelle: vorher gab es hier gar keine Prüfung, nachher eine über 8 von 177
  Dateien. Der Diff enthält keine gelockerte Schwelle, keine entfernte Regel und keine
  Inline-Suppression (`git diff fd4cbb4^..HEAD` berührt weder `.golangci.yml` noch ein
  `nolint`). Die Rücknahme der früheren *„bewusst verworfen"*-Notiz ist keine Lockerung,
  sondern eine Verschärfung, und sie beruft sich auf `AGENTS.md` §5 — eine höherrangige
  Quelle als ein Konfigurations-Kommentar. Dass die Regel dort seit `01db8a5` steht, ist
  nachgeprüft.
- **Der genannte Anlassfall wird tatsächlich anders behandelt, und das ist benannt.**
  `slice-094` und `slice-104` — die beiden Fälle, die die alte Ablehnung trugen — liegen
  in der Ausnahme (`slice-0??-*` bzw. `slice-1[0-6]?-*`) und melden nicht; der Plan sagt
  das ausdrücklich (*„Der Altbestand bleibt unberührt"*). `slice-168/169/170` sind in
  MR-056 als bewusst ausgenommen begründet. Kein stiller Wechsel.
- **§3.1 ist eingehalten.** Der Diff enthält keine neue Toolchain, kein Skript, keine
  Host-Abhängigkeit; die Regel läuft über das bestehende Image im bestehenden Target.
  Kein neues Make-Target, also auch keine `targets`-Deklarationslücke.
- **§3.3 ist eingehalten.** `e8401c7` ist ein reiner `git mv` plus die nach MR-013
  gekoppelten Verweise: Roadmap-Ruhe-Marker entfernt, Pfad-Verweise in `welle-86`, in den
  Closure-Notizen von `slice-177`/`slice-178` und in der Roadmap-Drift-Zeile umgehängt.
  Kein Inhalt am bewegten Slice selbst.
- **Die Zustellung widerspricht sich nicht.** `AGENTS.md` §4, `AGENTS.md` §5,
  `harness/README.md` und MR-056 nennen übereinstimmend: `max-open-tasks: 0`, den
  Grund-Code `section-tasks-open`, den Befund je Haken auf seiner Zeile, den Bindepunkt
  `verify-closure-notes` statt `gates`, und die Ausnahme bis `slice-170`. Die einzige
  Abweichung zwischen ihnen ist der Grenzen-**Umfang**, und der ist als M-2 gemeldet.
- **Die §4-Zeile bleibt trotz der Einfügung lesbar.** Die vorhandene Aussage
  *„**Drei** Grund-Codes können den Bindepunkt rot machen"* steht weiterhin im
  `spans`-Satzblock, dessen Subjekt zwei Sätze vorher eingeführt wird (*„Es deckt
  nicht …"*); sie behauptet keine geschlossene Menge über den ganzen Bindepunkt. Die
  neue DoD-Passage ist davor eingefügt, nicht hinein.
- **Keine Referenz-Richtungs-Verletzung.** Der Diff berührt kein Spec-Stratum; die
  Spec-Stellen-Zeile des Plans steht korrekt auf `—`. `make doc-check` (Modul `matrix`)
  ist grün.

---

## Kategorie-Summary

| Kategorie | Anzahl | Klassen |
|---|---|---|
| HIGH | 2 | `belegzahl-aus-der-spanne-gerechnet` · `kommentar-erzaehlt-die-datei-historie` |
| MEDIUM | 3 | `hint-auf-befund-ohne-bedingungs-verletzung` · `zwei-grenzen-zugesagt-drei-gemessen` · `vorpruefung-mit-veraltetem-register-stand` |
| LOW | 3 | `tabellenzeile-mit-ueberzaehliger-zelle` · `selektor-weiter-als-der-nachbar-ohne-grund` · `teilsumme-widerspricht-der-gesamtsumme` |
| INFO | 3 | Ausnahme-Zone · Glob-`?`-Semantik · veralteter Nenner in §1 |

**Wiederkehrende Klassen für den Steering Loop:** BEO-020 (H-1, L-3 — eine Zahl aus der
Spanne gerechnet statt aus dem Verzeichnis gelesen; wäre die sechste Instanz) · BEO-023
(M-2 — ein Wächter mit einer stillen Blindstelle, die seine Zustellung nicht nennt) ·
BEO-022 (M-2 — die Zustellung existiert, sagt aber weniger, als sie zusagt).
