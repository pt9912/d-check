# Review slice-178 — Eine Bedingung zählt offene Task-Items auf dem rohen Abschnitt


**Review-Art:** Code/Design (gegen Slice-Plan, ADR, Hard Rules)
**Gegenstand:** `3f8049e^..HEAD` = `3f8049e` (Prep) · `b2530f7` (Claim-Move) · `23b0f56` (feat)
**Skill:** `.harness/skills/reviewer.md` @ v1.13.0 · **Modell:** claude-opus-5[1m] · **Datum:** 2026-08-30
**Eingangs-Kontext:** `AGENTS.md` §2/§3.1/§3.2/§3.3/§3.4/§3.5/§3.6/§3.7/§3.8/§4/§5, `harness/conventions.md` (MR-013, MR-053, MR-054), `.harness/baseline/v5.12.0/regelwerk/` (`grundlagen-harness-dateien.md` §Was ein Kommentar trägt, `modul-04-adrs.md`, `modul-05-planning-harness.md`, `modul-06-roadmap.md`, `modul-10-review-harness.md`), `DC-FA-STRUCT-001`, `DC-FA-SPAN-001`, ADR-0042, ADR-0057, ADR-0059, ADR-0069, ADR-0073, ADR-0074, ADR-0077, Beobachtungs-Register (BEO-003, BEO-011, BEO-012, BEO-013, BEO-016, BEO-020, BEO-023), Slice-Plan §2 (die zehn Befunde des Vorgängers)

### Eigener Lauf

| Lauf | Ausgabe |
|---|---|
| `make gates` | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` · Exit 0 |
| `make doc-check` (in `gates`) | `d-check: 620 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `make test` (in `gates` und 10× einzeln) | alle Pakete `ok` · Exit 0 |
| 8 Mutationsproben (die behaupteten), je `make test` | alle acht reproduzieren, Zahlen exakt — s. Tabelle unter H-2 |
| 2 **eigene** Gegenproben, je `make test` | **beide bleiben grün** — s. H-2 und M-2 |
| 7 Produkt-Proben (`d-check:latest` und ein aus `b2530f7` gebautes Vorgänger-Image, je `--config <profil> --enable structure`, read-only, `--network none`) | s. u. |
| Zwei-Automaten-Scan (`awk`) über alle 690 `.md`-Dateien | **0** Fence-Lexik-Divergenzen; Selbsttest des Detektors an zwei synthetischen Dateien positiv |

`make fullbuild`, `make record-gates`, `make image-scan` bewusst **nicht** gefahren (Nachweis-Datei bzw. Netz). Arbeitsbaum nach allen Proben sauber (`git status --short` leer, Worktree entfernt, Vorgänger-Image gelöscht).

---

## Urteil

**BLOCKIERT** — 3 HIGH · 3 MEDIUM · 3 LOW · 3 INFO.

Der Neubau ist deutlich besser als sein Vorgänger. Sieben der zehn Befunde aus §2 sind
messbar geschlossen, die acht behaupteten Mutationsproben reproduzieren **alle** mit
exakt den genannten Zahlen, und die drei benannten Grenzen stimmen am Produkt gemessen
auf die Zeile genau. Byte-Identität und die `done/`-Korpus-Zahlen (144 Befunde,
37 Dateien, dieselben 37 wie die abgelöste Form) reproduzieren ebenfalls.

Blockierend sind drei Dinge, und alle drei sind Wiederholungen aus §2:

1. **Die Spezifikation widerspricht sich weiterhin** — §DC-FA-STRUCT-001.a Schritt 5
   deklariert nach wie vor **zwei** benannte Ausnahmen von der Bereinigung, während
   dasselbe Dokument drei roh lesende Bedingungen beschreibt. Die Commit-Botschaft
   behauptet, Lastenheft **und Spezifikation** seien korrigiert; korrigiert ist nur das
   Lastenheft.
2. **Der Kopplungs-Test kann nicht fangen.** `TestOffenerHaken_TeiltDieLexikMitTaskItemRE`
   prüft eine Implikation, die unter der gelieferten Implementierung eine Tautologie ist.
   Gemessen: genau die BEO-003-Drift, gegen die er namentlich steht, lässt die gesamte
   Suite grün. Das ist BEO-023 im Wächter gegen BEO-003 — dieselbe Doppelung wie beim
   Vorgänger, nur an einer anderen Stelle.
3. **Der stille Fence-Pfad ist entschieden, aber nirgends in der Anforderung benannt.**
   Gemessen bleibt ein echt offenes Task-Item unsichtbar, sobald ein Schluss-Fence fehlt;
   sichtbar wird es nur durch ein **anderes**, opt-in-Modul. Kein Stratum, kein Handbuch
   und kein Kommentar sagt das — und das Lastenheft sagt an derselben Stelle das
   Gegenteil (*„ein Wächter, den ein Tippfehler abschaltet, meldet grün"*).

---

## Findings

### H-1 · Die Spezifikation führt weiter **zwei** benannte Ausnahmen — und die Commit-Botschaft behauptet die Korrektur

- **quelle:** `AGENTS.md` §2 (Source Precedence, Rang 2), §5 (Proben-/Aussagen-Disziplin), BEO-020, Reviewer-Skill HIGH-nahe Hälfte von BEO-009; Slice-Plan §2 Befund 5 (*„Die Spezifikation widersprach sich intern"*)
- **pfad:** `spec/spezifikation.md:2193-2198` (§DC-FA-STRUCT-001.a Schritt 5) gegen `spec/spezifikation.md:2349` (neuer Block) und `spec/spezifikation.md:2302` (Zellenlängen-Block); Commit-Botschaft `23b0f56`
- **zugesagt:** die Botschaft von `23b0f56`, wörtlich: *„Dabei korrigiert: Lastenheft und
  Spezifikation fuehrten ZWEI Bedingungen, die einen anderen Text lesen als den
  bereinigten. Es sind VIER."*
- **gemessen:** `grep -n "benannten Ausnahmen" spec/spezifikation.md` → **eine** Fundstelle,
  `:2194`, unverändert:

  ```
  2193:   Schritt C4). Alle Bedingungen arbeiten auf **diesem** Text — mit **zwei**
  2194:   benannten Ausnahmen: die Chronologie-Bedingung (`table.order`, Schritt 6)
  2195:   liest die **rohen** Abschnitts-Zeilen, weil die Bereinigung Inline-Code
  2196:   leert und reale Schlüsselspalten genau dort stehen; die
  2197:   Überschriften-Bedingung (`headings-match`, Schritt 6) liest die
  2198:   Überschriften selbst.
  ```

  `git diff 3f8049e^..HEAD -- spec/spezifikation.md` berührt diesen Absatz nicht (Hunks
  liegen bei `:2345`, `:2913`, `:3061`, `:3090`). Das Wort *vier* kommt im gesamten
  `structure`-Bereich der Spezifikation genau einmal vor, und zwar in *„alle vier
  Listen-Marker"* (`:2357`). Korrigiert wurde ausschließlich der Lastenheft-Tabellenkopf
  (`spec/lastenheft.md:2511-2517`).
- **warum es zählt:** Schritt 5 ist eine **abschließende** Aufzählung („mit zwei benannten
  Ausnahmen"), auf die Schritt 6 und die Zellenlängen-Prosa (`:2274`, *„benannte Ausnahme
  aus Schritt 5"*) sich beziehen. Wer sie befolgt, entfernt beim nächsten Umbau die
  Roh-Lesung der dritten Bedingung regelkonform wieder — genau die Klasse, die die
  Spezifikations-Historie für §DC-FA-CLI-001.a schon einmal protokolliert hat
  (*„wer eine geschlossene Enumeration in der Spezifikation befolgt, entfernt den Zeiger
  beim nächsten Umbau regelkonform wieder"*). Dazu die zweite Hälfte: die Botschaft nennt
  eine Korrektur an einem Dokument, die dort nicht stattgefunden hat — eine Behauptung,
  die wie ein Beleg aussieht.
- **verifizierbar:** nein, kein Sensor — `make doc-check` ist grün (620 Dateien,
  0 Befunde); Beleg sind `grep` und der Diff.
- **klasse:** `enumerations-drift-im-stratum-plus-behauptete-korrektur`

### H-2 · Der Kopplungs-Test kann nicht fangen — die BEO-003-Drift lässt die ganze Suite grün

- **quelle:** BEO-023 (Zähler **6**), BEO-003 (verkörpert in welle-74; Register-Stand:
  *„Ein geteiltes Prädikat allein genügt nicht; es braucht je Konsument eine Assertion"*),
  Slice-DoD *„Umkehr-Probe: je Zusage eine Mutation, die genau einen Test rot macht"*,
  Reviewer-Skill HIGH #1
- **pfad:** `internal/hexagon/core/rules/structure_offene_tasks_test.go:166-183`
  (`TestOffenerHaken_TeiltDieLexikMitTaskItemRE`) gegen
  `internal/hexagon/core/rules/structure.go:336-338` (`offenerHaken`)
- **zugesagt:** der Kommentar über dem Test, wörtlich: *„DIE LEXIK IST GETEILT, NICHT
  KOPIERT: offenerHaken erkennt genau die Zeilen, die taskItemRE erkennt, verengt auf die
  leere Box. Ein zweites RE2 waere ein woertliches Praefix des ersten und driftete beim
  ersten Zusatz still auseinander (BEO-003)."* ADR-0074 `## Geschichte` (2026-08-30) und
  die Commit-Botschaft führen dieselbe Zusage: *„die Lexik kann nicht driften, statt dass
  ein Test ihre Kopplung ueberwacht."*
- **gemessen, Struktur:** die einzige Assertion des Tests lautet `if offen && !erkannt`.
  `offenerHaken` ist definiert als
  `strings.HasSuffix(taskItemRE.FindString(line), "[ ]")`; `FindString` liefert bei
  Nicht-Treffer den Leerstring, und `HasSuffix("", "[ ]")` ist `false`. `offen ⇒ erkannt`
  ist damit **wahr per Konstruktion** — die Bedingung ist unerfüllbar, der Test kann
  keinen Zustand melden. Die im Kommentar genannte Richtung (*taskItemRE wird erweitert,
  die Box-Erkennung nicht*) verletzt diese Implikation ebenfalls nicht.
- **gemessen, Gegenprobe B (die benannte Drift, hergestellt):** `taskItemRE` um die
  CommonMark-Form `1)` erweitert und daneben ein zweites RE2 als **wörtliches Präfix der
  alten Form** eingeführt, auf das `offenerHaken` umgestellt wird:

  ```
  327: var taskItemRE     = regexp.MustCompile(`^[ \t]*(?:[-*+]|[0-9]+[.)])[ \t]+\[[ xX]\]`)
  328: var openTaskItemRE = regexp.MustCompile(`^[ \t]*(?:[-*+]|[0-9]+\.)[ \t]+\[ \]`)
  338:     return openTaskItemRE.MatchString(line)
  ```

  `make test` → **Exit 0**, `--- FAIL`-Zeilen: **0**. Genau der Zustand, gegen den der
  Test namentlich steht, passiert die gesamte Suite.
- **gemessen, Gegenprobe (Kollaps):** `offenerHaken` → `taskItemRE.MatchString(line)`
  (Verengung auf die leere Box vollständig entfernt) macht drei Tests rot
  (`AlleMarkerFormen`, `EinBefundJeItem`, `NurImAbschnitt`) — der Kopplungs-Test ist
  **keiner** davon.
- **warum es zählt:** BEO-023 ist die Klasse *„ein Wächter, der nie fangen konnte, liest
  sich wie einer, der fängt"*, und ihre zweite Instanz **ist** slice-178. Der Neubau
  wiederholt sie in dem Test, der die Lehre einlösen soll. Zusätzlich: die verkörperte
  Form von BEO-003 verlangt ausdrücklich *„je Konsument eine Assertion"* — hier gibt es
  eine Assertion, die keinen Konsumenten prüft. Und der Test steht in der DoD-Zeile
  *„Umkehr-Probe … die Probe kostet einen Lauf und ist der einzige Beleg, dass der
  Wächter beißt"*: die acht gefahrenen Mutationen belegen für die **anderen** acht Tests,
  dass sie beißen — für diesen ist keine gefahren worden, und es gibt keine, die es
  könnte.
- **verifizierbar:** ja, ohne Gate — Gegenprobe B ist reproduzierbar und endet grün.
- **klasse:** `tautologischer-kopplungstest`

### H-3 · Der stille Fence-Pfad ist entschieden, aber in keiner Anforderung benannt — und das Lastenheft sagt das Gegenteil

- **quelle:** `AGENTS.md` §3.8 (*„Wo sie nicht gilt, gehört die Grenze in die
  Anforderung"*), Reviewer-Skill HIGH #1 (stiller Grün-Pfad) und MEDIUM *Modul-Grenze auf
  der Ziel-Achse* (BEO-004), Slice-Plan §2 Befund 1, BEO-016 (*„die Klasse ist breiter als
  gedacht: derselbe stille Ausfall entsteht über einen vergessenen Schluss-Fence"*)
- **pfad:** `spec/lastenheft.md:2535-2557` (Prosa-Block und *„Drei Grenzen"*) ·
  `spec/spezifikation.md:2349-2373` (Ablauf-Block, *„Grenze, benannt statt zugesagt"*) ·
  `docs/user/benutzerhandbuch.md:2178-2202` · `internal/hexagon/core/rules/structure.go:340-356`
  (Kommentar-Blöcke ZUSAGE/GRENZE)
- **zugesagt:** `spec/lastenheft.md:2539-2540`, wörtlich: *„für eine Zusage, die eine
  **Closure-Vorbedingung** trägt, kippt er: ein Wächter, den ein Tippfehler abschaltet,
  meldet grün, wo er nichts gesehen hat."* Darunter *„Drei Grenzen, und sie sind gemessen,
  nicht geschätzt"* — Fence, einzeilige/mehrzeilige Inline-Spanne, Blockquote/Tabulator.
  Der Kommentar in `structure.go:349-353` nennt dieselben drei.
- **gemessen:** zwei Fixtures, identisch bis auf einen fehlenden Schluss-Fence, Profil
  `max-open-tasks: 0` auf dem DoD-Abschnitt:

  ```
  printf '# F\n\n## 4. Definition of Done\n\n```md\nBeispiel\n\n- [ ] wirklich offener Punkt\n' > docs/fixture/blind.md
  printf '# F\n\n## 4. Definition of Done\n\n- [ ] wirklich offener Punkt\n'                    > docs/fixture/laut.md
  ```

  ```
  --- --enable structure ---
  docs/fixture/laut.md:5   … section-tasks-open  offenes Task-Item über der Grenze von 0
  d-check: 560 Datei(en) geprüft, 1 Befund(e)        (blind.md: KEIN Befund)

  --- --enable structure --enable spans ---
  docs/fixture/blind.md:5  ```md   fence-unclosed
  docs/fixture/laut.md:5   … section-tasks-open  offenes Task-Item über der Grenze von 0
  ```

  Ein **echt** unquittierter Haken bleibt für die Bedingung unsichtbar; sichtbar wird nur
  die Fence-Ursache, und die meldet ein **anderes**, opt-in-Modul mit einem anderen
  Grund-Code.
- **die Entscheidung selbst ist nicht der Befund.** Sie ist getroffen und begründet
  (Slice-Plan §6, ADR-0074 `## Geschichte` 2026-08-30: Fall 1 fängt `fence-unclosed`,
  Fall 2 hat null Realfälle). Beides habe ich nachgemessen und beides trägt — der
  Zwei-Automaten-Scan über alle 690 `.md`-Dateien findet **0** Divergenzen (s. I-2).
  Der Befund ist, dass **keine** der vier Deklarations-Oberflächen die Grenze führt:
  `grep -n "fence-unclosed\|Schluss-Fence\|unbeendet" spec/lastenheft.md spec/spezifikation.md docs/user/benutzerhandbuch.md`
  liefert im gesamten `structure`-Bereich **null** Treffer; alle Treffer liegen im
  `spans`-Bereich (`DC-FA-SPAN-001`). Damit steht die **Abhängigkeit** von einem zweiten
  Modul nirgends: dass die Zusage nur so weit trägt, wie ein Konsument `spans` mitfährt,
  ist eine Eigenschaft des Makefile-Rezepts dieses Repos
  (`verify-closure-notes: … --enable structure --enable spans`) und keine des
  ausgelieferten Schlüssels.
- **warum es zählt:** §3.8 stellt genau diese Frage — welche Eingabe liest das Modul, für
  die seine Zusage nicht gilt. Hier ist es dieselbe Datei, nur hinter einem Fence-Zeichen.
  Der Satz *„ein Wächter, den ein Tippfehler abschaltet"* steht als **Begründung der
  Anforderung** im ranghöchsten Stratum und beschreibt einen Zustand, der weiter besteht;
  ein Adopter, der `max-open-tasks` ohne `spans` fährt, liest eine Zusage, die er nicht
  bekommt. Das ist derselbe Defekt, der zur zweiten Rückführung geführt hat — diesmal
  nicht im Code, sondern in der Deklaration.
- **verifizierbar:** ja am Produkt (Fixture oben); kein Gate.
- **klasse:** `zusage-ohne-benannte-modul-grenze`

### M-1 · ADR-0074 verspricht die Abgrenzung in **beiden** Schema-Zeilen — geliefert ist eine

- **quelle:** ADR-0074 §Konsequenzen (`Accepted`, Kern), Slice-Plan §5 Risiko 1, Reviewer-Skill MEDIUM *Konsistenz-Lücke zwischen Modulen derselben Eingabe-Klasse*
- **pfad:** `docs/plan/adr/0074-offene-tasks-auf-rohen-zeilen.md:123-125` gegen
  `spec/spezifikation.md:2915` (`structure[].max-tasks`), `spec/lastenheft.md:2523`
  (Tabellenzeile `max-tasks`) und `internal/hexagon/core/model/config.go:455`
  (`MaxTasks *int`)
- **zugesagt:** wörtlich: *„`max-tasks` (bereinigt, alle Items) und `max-open-tasks` (roh,
  offene Items). Wer den falschen greift, bekommt stillschweigend die schwächere Zusage.
  Dagegen steht die Abgrenzung in beiden Schema-Zeilen, nicht ein Sensor."*
- **gemessen:** die Abgrenzung existiert genau einmal, auf der **neuen** Seite —
  `spec/spezifikation.md:2916` (*„**Abgrenzung zu `max-tasks`:** jenes zählt **alle** Items
  auf dem **bereinigten** Text"*) und `internal/hexagon/core/model/config.go:462-464`
  (`ABGRENZUNG zu MaxTasks`). Die Gegenrichtung fehlt vollständig:

  ```
  spec/spezifikation.md:2915: | `structure[].max-tasks` | int | abwesend (aus) | Obergrenze der Task-Items **im Abschnitt**, nicht dateiweit ⇒ sonst `section-oversized`; **explizit** < 0 ⇒ Exit 2 |
  spec/lastenheft.md:2523:    | `max-tasks` (int ≥ 0; erklärte Teilmenge über `tasks-ignore-pattern`) | `section-oversized` | zerlegen statt dehnen |
  internal/…/model/config.go:455:  MaxTasks       *int
  ```

  Auch `internal/hexagon/core/model/finding.go` trägt beim älteren
  `ReasonSectionOversized` (`:60`) keinen Zeiger auf den neuen Code, während der neue
  (`:92-98`) einen auf den alten trägt.
- **warum es zählt:** die Richtung, die fehlt, ist die einzige, die wirkt. Wer bereits
  `max-tasks` benutzt oder es in der Schema-Tabelle findet, ist der Leser, der die
  schwächere Zusage bekommt; wer `max-open-tasks` liest, hat die stärkere schon gewählt.
  Die ADR benennt die Mitigation als das, was an die Stelle eines Sensors tritt — sie ist
  halb geliefert und damit an der Stelle unwirksam, für die sie gedacht war.
- **verifizierbar:** nein, kein Sensor; Beleg ist `grep`.
- **klasse:** `einseitige-abgrenzung-gegen-adr-zusage`

### M-2 · Die Abschnittsgrenze hat weiterhin keine beißende Probe

- **quelle:** BEO-023, Slice-Plan §2 Befund 3 (*„Die Abschnittsgrenze hielt kein Test"*),
  ADR-0074 `## Geschichte` 2026-08-30 (*„über die ganze Datei statt im Abschnitt gezählt ⇒
  1 rot"*), Spezifikation `:2367-2369` (*„Abschnittsgrenze. Gezählt wird bis zur nächsten
  Überschrift derselben oder höherer Ebene, nicht dateiweit."*)
- **pfad:** `internal/hexagon/core/rules/structure.go:366` ·
  `internal/hexagon/core/rules/structure_offene_tasks_test.go:125-138`
  (`TestMaxOpenTasks_NurImAbschnitt`)
- **zugesagt:** die Zusage ist die Abschnittsgrenze; die gefahrene Mutation ist
  *„über die ganze Datei statt im Abschnitt"* (`for i := 0; i < len(lines)`), und sie
  macht korrekt genau einen Test rot — das habe ich reproduziert.
- **gemessen, eigene Gegenprobe A:** die Grenze um **eins** verschoben statt aufgehoben —
  `for i := headingNo; i < end-1` → `for i := headingNo; i < end-2`. Damit fällt die
  **letzte** Zeile jedes Abschnitts aus der Zählung; ein offenes Item unmittelbar vor der
  nächsten Überschrift wird still übersehen. `make test` → **Exit 0**, `--- FAIL`-Zeilen:
  **0**. Kein Test deckt diese Lage: in allen sechs Fixtures steht zwischen dem letzten
  Item und der nächsten Überschrift (bzw. dem Dateiende) mindestens eine Leerzeile.
- **zweite, nicht mutations-gemessene Lücke derselben Zusage:** kein Test kombiniert
  mehrere passende Abschnitte (`sections: each`) mit einer Schwelle **> 0**. Die
  Rückstellung des Budgets je Abschnitt (`frei := *r.MaxOpenTasks` in
  `structure.go:364`, lokal je Aufruf) ist damit unbelegt — `TestMaxOpenTasks_NurImAbschnitt`
  fährt zwei Abschnitte, aber mit Schwelle 0, wo ein geteiltes und ein
  zurückgesetztes Budget dasselbe Ergebnis liefern.
- **warum es zählt:** die DoD verlangt *„je Zusage eine Mutation, die genau einen Test rot
  macht"*, und ADR-0074 führt die Abschnittsgrenze als eine der acht. Die gewählte
  Mutation ist die grobe; die naheliegende Regressions-Form (Off-by-one an derselben
  Zeile) läuft grün durch. Nach dem Vorgänger, bei dem dieselbe Zusage **gar** keinen Test
  hatte, ist sie damit halb geschlossen, nicht ganz.
- **verifizierbar:** ja — Gegenprobe A ist reproduzierbar und endet grün.
- **klasse:** `zusage-nur-grob-gewaechtert`

### M-3 · Zwei Kern-Kommentare tragen weiter die alte Zahl — einer davon in der Funktion, die die neue Bedingung umgeht

- **quelle:** `AGENTS.md` §3.7 (*„Ein Kommentar beschreibt, was da ist"*), Slice-Plan §2
  Befund 4 (*„Die Aufzählung stand an vier Stellen falsch"*) und §2 *„die Zahl ‚drei' ist
  überall ‚vier'"*, BEO-011 (*„wer ‚einzig/nur/alle' schreibt, zählt vorher die
  Nachbarn"*)
- **pfad:** `internal/hexagon/core/rules/sections.go:62-64` ·
  `internal/hexagon/core/rules/structure_tableorder.go:4-5` ·
  (Kollision) `docs/user/benutzerhandbuch.md:2072` gegen `spec/spezifikation.md:2349`
- **zugesagt:** `internal/hexagon/core/rules/structure.go:239-243` ist korrekt auf **vier**
  gezogen; die Commit-Botschaft und der Plan sagen, die Zahl sei „überall" nachgezogen.
- **gemessen:**

  ```
  internal/hexagon/core/rules/sections.go:62-64
    // Die Bedingungen beider Konsumenten arbeiten auf diesem Text —
    // mit zwei benannten Ausnahmen: die Chronologie-Monotonie liest die rohen
    // Abschnitts-Zeilen, die Überschriften-Bedingung die Überschriften selbst.

  internal/hexagon/core/rules/structure_tableorder.go:4-5
    // §DC-FA-STRUCT-001.a Schritt 6, ADR-0057): die eine structure-Bedingung, die
    // die ROHEN Abschnitts-Zeilen liest …
  ```

  Dazu eine **Ordinal-Kollision**, die dieser Slice neu erzeugt:
  `docs/user/benutzerhandbuch.md:2072` nennt die **Zellenlängen**-Bedingung *„die dritte
  Bedingung auf den rohen Abschnitts-Zeilen"*, `spec/spezifikation.md:2349` nennt
  `max-open-tasks` *„die **dritte** Bedingung auf den rohen Abschnitts-Zeilen"*, und
  `spec/spezifikation.md:2302` führt die Zellenlänge als *„die **zweite**"*. Drei
  Dokumente, zwei Gegenstände, dieselbe Ordnungszahl.
- **warum es zählt:** `sections.go:62-64` ist der Doc-Kommentar von `SectionProse` — genau
  der Funktion, an der die neue Bedingung **vorbeigeht**. Wer sie ändert, liest dort eine
  abschließende Liste der Ausnahmen und findet die neue nicht; die Kopplung, die der
  Kommentar tragen soll, zeigt ins Leere. Der Kommentar in `structure_tableorder.go` war
  schon vor diesem Slice falsch (ADR-0069 machte die Zellenlänge zur zweiten), aber der
  Slice hat sich das Nachziehen der Zahl ausdrücklich vorgenommen und die beiden Stellen
  nicht gemessen.
- **verifizierbar:** nein, kein Sensor; Beleg ist `grep -rn "zwei benannten Ausnahmen\|die eine structure-Bedingung" --include=*.go internal/`.
- **klasse:** `enumerations-drift-im-kommentar`

### L-1 · Die Modul-Tabelle des Handbuchs führt weiter neun Bedingungen und kennt den neuen Grund-Code nicht

- **quelle:** Reviewer-Skill LOW *Doku-Drift*; Slice-DoD *„Das Benutzerhandbuch führt die Bedingung dort, wo es die übrigen `structure`-Schlüssel führt"*
- **pfad:** `docs/user/benutzerhandbuch.md:2230` (Modul-Tabelle) · nachrangig `README.md:130-141`
- **gemessen:** die Zeile sagt *„bis zu **neun** Bedingungen mit je eigenem Grund-Code"*
  und listet in der Grund-Code-Spalte 14 Codes — ohne `section-tasks-open`. Es fehlt
  dort auch `section-exempt-mismatch` (aus slice-182, also Bestand). `README.md:130`
  sagt englisch ebenfalls *„nine conditions"*.
- **warum es zählt:** der Slice hat das Handbuch **angefasst** (Konfigurations-Beispiel
  `:1993-1995` und Prosa-Block `:2178-2202`), also greift hier nicht die
  Release-Prep-Konvention, die README und `CHANGELOG.md` abdeckt (belegt: `9c383ba` hat
  das Handbuch gar nicht berührt, `bd57f07` hat es nachgezogen). Die Grund-Code-Spalte ist
  die Stelle, an der ein Anwender nachschlägt, welche Codes ein Modul überhaupt erzeugen
  kann.
- **verifizierbar:** nein; kein Sensor hält Prosa-Modullisten (benannte Grenze in
  `AGENTS.md` §4 zu `targets`).
- **klasse:** `modultabellen-drift`

### L-2 · Überlappende `sections: each` melden auf Zeilen, deren eigener Abschnitt die Schwelle einhält

- **quelle:** Reviewer-Skill MEDIUM/LOW *dokumentationswürdige, aber undokumentierte Annahme*; ADR-0074 §Entscheidung 4
- **pfad:** `internal/hexagon/core/rules/structure.go:357-378` · `spec/spezifikation.md:2362-2366`
- **gemessen:** Fixture mit `# Titel` (3 offene Items) und darin `## Unter` (2 offene
  Items), Regel `section-pattern: '^#'`, `sections: each`:

  ```
  --- max-open-tasks: 2 ---
  docs/fixture/nest.md:5   … section-tasks-open  offenes Task-Item über der Grenze von 2
  docs/fixture/nest.md:9   … section-tasks-open  offenes Task-Item über der Grenze von 2
  docs/fixture/nest.md:10  … section-tasks-open  offenes Task-Item über der Grenze von 2

  --- max-tasks: 2 (dieselbe Datei, dieselbe Regel-Form) ---
  docs/fixture/nest.md:1   … section-oversized   Abschnitt trägt 5 Task-Items, erlaubt sind 2
  ```

  Die Befunde auf `:9` und `:10` stammen aus dem **umschließenden** H1-Abschnitt; der
  Abschnitt `## Unter`, in dem sie stehen, trägt genau zwei offene Items und erfüllt die
  Bedingung.
- **warum es zählt:** die neue Kombination *ein Befund je Item* + *Schwelle > 0* macht die
  Überlappung erstmals sichtbar und irreführend — bei `max-tasks` liegt der Befund auf der
  Abschnitts-Überschrift und ist eindeutig zuzuordnen, hier liegt er auf einer Zeile,
  deren nächstliegende Regel-Instanz grün ist. Das Closure-Profil dieses Repos benutzt
  genau die Form `section-pattern: '^# '` mit `sections: each`
  (`.d-check.closure.yml`), der Fall ist also erreichbar. Die Spezifikation sagt
  *„die ersten `max-open-tasks` offenen Items in Dokument-Reihenfolge"* ohne den Zusatz,
  dass bei überlappenden Treffern jeder Abschnitt sein eigenes Kontingent führt.
- **verifizierbar:** ja am Produkt (Fixture oben).
- **klasse:** `ueberlappende-abschnitte-mit-eigenem-kontingent`

### L-3 · Ein Zitat in Anführungszeichen, das die Quelle so nicht schreibt

- **quelle:** `AGENTS.md` §5 (BEO-012, *„Eine zitierte Quelle trägt nur, was in ihrem Geltungsbereich steht"*), Reviewer-Skill MEDIUM *Geltungsbereich*
- **pfad:** `docs/plan/adr/0074-offene-tasks-auf-rohen-zeilen.md:172` (`## Geschichte`,
  2026-08-30) · `docs/plan/planning/in-progress/slice-178-offene-tasks-roh.md:225` ·
  Commit-Botschaft `3f8049e`
- **zugesagt:** dreimal gleichlautend, in Anführungszeichen: ADR-0042 *„hat die Frage
  ausdrücklich offen gelassen und ihre Bedingung benannt — sie bekommt ‚erst eine Regel,
  wenn ein Realfall existiert'."*
- **gemessen:** `docs/plan/adr/0042-markdown-lexik-folgt-commonmark.md:136-137` lautet
  *„Beide sind **unbelegt** — kein Realfall in den 522 Dateien — und bekommen erst eine
  Regel, wenn einer existiert."* Die zitierte Zeichenkette steht dort nicht; sie ist eine
  Paraphrase in Zitat-Zeichen.
- **die Geltungs-Prüfung selbst besteht:** der Satz in ADR-0042 gilt beiden dort offen
  gelassenen Punkten, und der hier gemeinte (der naive Toggle) ist Punkt (a) — die
  **Aussage** trägt also. Beanstandet ist nur die Form: ein Zitat sieht aus wie ein Beleg.
  `citations` fängt es nicht, weil keine `cite`-Direktive daranhängt.
- **verifizierbar:** nein; Beleg ist der Vergleich mit `:136-137`.
- **klasse:** `paraphrase-in-zitatzeichen`

### I-1 · „276 Befunde" ist aus den Artefakten nicht nachvollziehbar — die Aussage dahinter trägt

- **pfad:** Commit-Botschaft `23b0f56` (*„derselbe Korpus (Slice-Plaene, max-tasks 3)
  liefert vor und nach der Aenderung 276 Befunde, `diff` leer"*)
- **gemessen:** die Byte-Identität ist unabhängig bestätigt — Vorgänger-Image aus
  `b2530f7` gebaut, beide Images gegen **denselben** Baum:

  | Korpus | Regel | vorher | nachher | `diff` |
  |---|---|---|---|---|
  | A | `done/slice-*.md`, `section-pattern: '^## '`, `sections: each`, `max-tasks: 3` | 162 | 162 | leer |
  | B | dieselbe Datei-Menge, `section-pattern: '^#{1,3} '` | 327 | 327 | leer |

  Die Zahl **276** trifft keiner der naheliegenden Zuschnitte; das Profil, das sie
  erzeugt, ist in keinem Artefakt genannt. Nach BEO-020 (*„die gezählte Menge benennen,
  bevor die Zahl fällt"*) fehlt der Nenner. **Die beiden anderen Produkt-Zahlen der
  Botschaft reproduzieren dagegen exakt** (s. Negativbefunde).
- **klasse:** `zahl-ohne-benannten-korpus`

### I-2 · Der „einzige Treffer" der Fence-Divergenz ist mit einem Zwei-Automaten-Scan nicht auffindbar — das Ergebnis ist trotzdem stärker als behauptet

- **pfad:** `docs/plan/planning/in-progress/slice-178-offene-tasks-roh.md:220-226` ·
  ADR-0074 `## Geschichte` 2026-08-30
- **zugesagt:** *„Beide Automaten über alle 620 Markdown-Dateien gegeneinander gefahren:
  der **einzige** Treffer ist die Prosa von ADR-0042, die den Unterschied beschreibt."*
- **gemessen:** eigener `awk`-Detektor, der beide Lesarten Zeile für Zeile mitführt
  (Toggle: jede links-getrimmte Fence-Zeile kippt, Backtick-Infozeilen-Regel wie
  `FenceToggle`; CommonMark: Zeichen-, Längen- und Whitespace-Abgleich wie `FenceCloses`)
  und die erste Zustands-Divergenz je Datei meldet. Über alle **690** `.md`-Dateien:
  **0** Treffer, ADR-0042 eingeschlossen. Selbsttest des Detektors an zwei synthetischen
  Dateien (`~~~` schließt einen Backtick-Fence; ```` ```txt ```` als Schließer) meldet
  beide korrekt.
- **warum das notiert gehört:** die Schlussfolgerung (*null Realfälle*) hält — sie hält
  sogar strenger als behauptet. Aber der genannte „einzige Treffer" ist mit einem
  Zustands-Vergleich nicht reproduzierbar; die zugrundeliegende Messung ist im Artefakt
  nicht so beschrieben, dass ein zweiter Leser dieselbe Zahl bekommt.
- **klasse:** `messverfahren-nicht-beschrieben`

### I-3 · ADR-0074 §Fitness Function nennt sieben Kern-Tests; gebaut sind neun

- **pfad:** `docs/plan/adr/0074-offene-tasks-auf-rohen-zeilen.md:143-149` ·
  `internal/hexagon/core/rules/structure_offene_tasks_test.go`
- **gemessen:** der Kern zählt sieben Zusagen auf; die Datei trägt neun `Test…`-Funktionen
  (zusätzlich `NurImAbschnitt` und `InlineCodeGrenze`, dazu der unter H-2 beanstandete
  Kopplungs-Test). Die `## Geschichte` hat die **Mutations**-Zahl (drei → acht) nachgeführt,
  die Test-Zahl nicht.
- **warum das kein MEDIUM ist:** der Kern einer `Accepted`-ADR darf nach `AGENTS.md` §3.5
  nicht überschrieben werden, und der 2026-08-29-Eintrag der `## Geschichte` sagt bereits,
  dass die alte Fitness-Function-Liste Tests beschreibt, die es nicht mehr gibt. Die
  Lücke ist damit deklariert, nur nicht beziffert.
- **klasse:** `fitness-function-zahl-nicht-nachgefuehrt`

---

## Die zehn Vorgänger-Befunde aus §2 — Stand je Befund

| # | Befund des Vorgängers | Stand | Beleg |
|---|---|---|---|
| 1 | Ein vergessener Schluss-Fence schaltet die Bedingung ab | **offen als Deklaration** (Verhalten bewusst so entschieden) | H-3: Fixture-Messung; keine der vier Oberflächen nennt die Grenze |
| 2 | Regressions-Test traf die Blindstelle nicht | **geschlossen** | `structure_offene_tasks_test.go:20-40` trägt den Vorzustand als eigene Assertion (`forbid-pattern` muss blind sein) und fällt bei Mutation 6 rot |
| 3 | Abschnittsgrenze hielt kein Test | **halb geschlossen** | M-2: grobe Mutation rot, Off-by-one grün (Exit 0) |
| 4 | „drei statt vier" Roh-Bedingungen | **halb geschlossen** | Lastenheft korrigiert; H-1 (Spezifikation Schritt 5) und M-3 (zwei Kern-Kommentare) offen |
| 5 | Spezifikation widersprach sich intern | **halb geschlossen** | Schritt 6, Schritt 7, §2-Schema, §4 und Historie sind gezogen; Schritt 5 nicht — H-1 |
| 6 | Lexik war kopiert, nicht geteilt | **anders gelöst — und der Beleg trägt nicht** | Es gibt kein zweites RE2 mehr (`structure.go:336-338`), das ist die stärkere Form; der Test, der sie sichert, kann nicht fangen — H-2 |
| 7 | Bei Schwelle > 0 meldete sie alle Items | **geschlossen** | `SchwelleMeldetNurDenUeberhang` prüft Zeilen 7/8 bei Grenze 2; Mutation 3 macht ihn rot |
| 8 | Inline-Code-Grenze überzeichnet | **geschlossen** | Produkt-Fixture: einzeilige Spanne stumm, mehrzeilige meldet auf ihrer Zeile; `InlineCodeGrenze` deckt beides |
| 9 | Blockquote und Tab-in-der-Box ungenannt | **geschlossen** | in Lastenheft, Spezifikation, Handbuch und `AlleMarkerFormen` als Grenze geführt; Produkt-Fixture bestätigt beide |
| 10 | BEO-016 in der Sichtung übersehen | **geschlossen** | Slice-Plan §7 führt BEO-016 mit der Begründung *„das ist die Klasse dieses Defekts"* |

---

## Negativbefunde (geprüft, ohne Befund)

- **Die acht behaupteten Mutationsproben reproduzieren alle, mit exakt den genannten
  Zahlen.** Je Mutation ein `make test`; gezählt sind rote Test-**Funktionen**:

  | # | Mutation | behauptet | gemessen |
  |---|---|---|---|
  | 1 | Fence-Gate entfernt (`!prose[i+1]` gestrichen) | 1 | 1 — `FenceBleibtAussen` |
  | 2 | Box nicht auf leer verengt | 3 | 3 — `AlleMarkerFormen`, `EinBefundJeItem`, `NurImAbschnitt` |
  | 3 | Schwelle ignoriert | 1 | 1 — `SchwelleMeldetNurDenUeberhang` |
  | 4 | über die ganze Datei statt im Abschnitt | 1 | 1 — `NurImAbschnitt` |
  | 5 | ein Befund je Datei (`break`) | 2 | 2 — `EinBefundJeItem`, `SchwelleMeldetNurDenUeberhang` |
  | 6 | bereinigt statt roh gelesen | 2 | 2 — `BacktickSchaltetNichtAb`, `InlineCodeGrenze` |
  | 7 | negativer Wert geschluckt | 1 | 1 — `TestDecode_StructureFehler` |
  | 8 | Verdrahtung Config → Modell entfernt | 1 | 1 — `TestDecode_StructureFehler` |

- **Die zwei Produkt-Zahlen über den `done/`-Bestand reproduzieren exakt.** Profil
  `done/slice-*.md`, `section-pattern: '^## [0-9]+\. Definition of Done'`,
  `max-open-tasks: 0` → `144 Befund(e)` in **37** Dateien. Dieselbe Datei-Menge mit
  `forbid-pattern: '- \[ \]'` statt der neuen Bedingung → ebenfalls **37** Dateien, und
  die beiden sortierten Datei-Listen sind `diff`-gleich. Die Botschaft sagt genau das und
  nicht mehr (*„der Gewinn hier ist die Granularität … NICHT mehr Treffer"*) — kein
  überdehnter Schluss.
- **Die drei benannten Grenzen stimmen am Produkt, Zeile für Zeile.** Ein Fixture mit
  einzeiliger Inline-Spanne, mehrzeiliger Inline-Spanne, Blockquote, Tab-Box,
  `*   [ ]` und `2. [ ]` liefert genau drei Befunde, auf den Zeilen der mehrzeiligen
  Spanne, des Sternchen-Items und des geordneten Items. Keine der vier stummen Formen
  meldet, keine der zwei erwarteten fehlt.
- **`structureAmAbschnitt` ist verhaltenserhaltend.** Die Reihenfolge der Bedingungen ist
  unverändert (`structureConditions` → `table.order` → `table.column` → **neu** →
  `headings-match`), und die Byte-Identitäts-Läufe über zwei Korpora (162 und 327 Befunde)
  sind `diff`-gleich. Der Schnitt selbst ist ein echter — Datei-Ebene trägt die
  Kandidaten-Auswahl (Ausnahme-Ventile, Kardinalität, `prose`-Aufbau), Abschnitts-Ebene
  die Bedingungen —, keine Verlegenheits-Extraktion.
- **Keine Suppression, keine gesenkte Schwelle** (§3.2/§3.6): `git diff 3f8049e^..HEAD`
  enthält kein `nolint`; `THRESHOLD ?= 93` im `Makefile` unverändert; `make lint` und
  `make coverage-gate` in `gates` grün.
- **Hexagon-Richtung** (§3.4, ADR-0005/ADR-0012): keine neuen Imports; `structure.go`
  bleibt bei `regexp`/`strconv`/`strings` und `core/model`. `make arch-check` (a-check,
  netzlos) in `gates` grün.
- **Kein Netz außerhalb `external`** (`DC-QA-03`): die Bedingung ist rein rechnend;
  `make doc-check` läuft mit `--network none` und meldet 620 Dateien, 0 Befunde.
- **Commit-Zerlegung** (§3.3, MR-013): `b2530f7` ist ein reiner `git mv` (die Slice-Datei
  selbst mit 0 geänderten Zeilen) plus die gekoppelten Verweise — Roadmap-Ruhe-Marker
  entfernt, `observations.md` und die Closure-Notiz von slice-180 auf `in-progress/`
  gezogen. `make planning-check` in `gates` grün. Die zwei Review-Reports zu slice-180
  sind bewusst **nicht** umgehängt, mit Begründung in der Botschaft — Lauf-Belege bleiben
  stehen, und `links` ist grün.
- **Kommentar-Klassen** (§3.7): die neuen Blöcke in `model/config.go:456-465`,
  `model/finding.go:92-98` und `rules/structure.go:329-356` tragen Zusage, Kopplung,
  Abgrenzung und Grenze; keine Review-Historie, keine Befund-Marker, keine Slice-Nummern,
  keine Mess-Labels. Der eine Konjunktiv (*„ein eigenes RE2 waere ein woertliches Praefix
  … und driftete beim ersten Zusatz still auseinander"*) zeigt nach **vorn** — er
  adressiert den nächsten Ändernden und beschreibt den Bruch der Zusage; das ist die
  ausdrücklich zulässige Form (Baseline §Zeitform-Test), nicht die verworfene Alternative.
- **Zustandsfelder** (§3.7): `observations.md` und `roadmap.md` sind ausschließlich um
  Pfad-Verweise bzw. den Ruhe-Marker geändert; keine Chronik in einer `Stand`-Zelle, kein
  Eintrag im Drift-Log.
- **ADR-Immutabilität** (§3.5, ADR-0016): `23b0f56` fügt ADR-0074 ausschließlich zwei
  Zeilen unter `## Geschichte` an, der Kern ist unberührt. Die inhaltliche Abweichung von
  §Entscheidung 3 (`openTaskItemRE`) ist dort begründet und ist eine **Übererfüllung** der
  Entscheidung (*„Die Lexik kommt aus dem Modul, nicht aus der Konfiguration"*), keine
  Umkehr — eine Folge-ADR ist dafür nach `modul-04-adrs.md` nicht verlangt. Beanstandet
  ist nur, dass der Beleg dafür nicht trägt (H-2).
- **Deklarations-Oberflächen:** `--print-config` führt den Schlüssel (Produkt-Lauf im
  Image, Zeile 180 der Ausgabe); `--doctor`-Klartext und `AllReasons()` sind ergänzt und
  durch den bestehenden Lockstep-Test gegen die Spezifikations-§4-Tabelle verriegelt
  (`internal/hexagon/core/app/diagnose_test.go:17,41`); Spezifikation §2-Schema, §4
  (`SPEC-078`), Ablauf-Block und Schritt 7 sind gezogen; Lastenheft-Tabelle, drei
  Akzeptanzkriterien, Versions-Bump 0.79.0 und Historie ebenso. `--suggest-config` gibt
  keine `structure`-Bedingungen aus (`grep` in `internal/hexagon/core/app/suggest.go` →
  0 Treffer), hat also keine Fläche.
- **§3.8, Ziel-Achse:** außer dem unter H-3 benannten Punkt liest die neue Bedingung keine
  Eingabe außerhalb der Scan-Menge des Moduls — sie arbeitet auf demselben `content`, das
  `structure` ohnehin über seine eigenen Globs liest.
- **Config-Rand:** `max-open-tasks: -1` ⇒ Exit 2 ist implementiert
  (`configyaml.go:328-330`) **und** durch einen Test gedeckt, der bei Entfernen der Prüfung
  rot wird (Mutation 7); die explizite Null wird bis ins Modell durchgereicht und ist von
  einem abwesenden Schlüssel unterscheidbar (Mutation 8 plus
  `TestMaxOpenTasks_AbwesendIstAus`).
- **`hint`-Interaktion:** `structureOpenTasks` meldet über `structureFinding`, also mit
  `MessageFor` — dieselbe Behandlung wie alle messenden Bedingungen; die `raw`-Form bleibt
  den Befunden vorbehalten, bei denen die Regel gar nicht gemessen hat (ADR-0073).
- **§3.9 Workflow-Pins:** unberührt; `make workflow-pins` in `gates` grün.

---

## Kategorie-Summary

| Kategorie | Anzahl | Klassen |
|---|---|---|
| HIGH | 3 | `enumerations-drift-im-stratum-plus-behauptete-korrektur` · `tautologischer-kopplungstest` · `zusage-ohne-benannte-modul-grenze` |
| MEDIUM | 3 | `einseitige-abgrenzung-gegen-adr-zusage` · `zusage-nur-grob-gewaechtert` · `enumerations-drift-im-kommentar` |
| LOW | 3 | `modultabellen-drift` · `ueberlappende-abschnitte-mit-eigenem-kontingent` · `paraphrase-in-zitatzeichen` |
| INFO | 3 | `zahl-ohne-benannten-korpus` · `messverfahren-nicht-beschrieben` · `fitness-function-zahl-nicht-nachgefuehrt` |

**Wiederholte Klassen in dieser Sitzung** (Steering-Loop-Signal): *Enumerations-Drift*
tritt zweimal auf (H-1 im Stratum, M-3 im Kommentar) und ist in der
Spezifikations-Historie bereits zweimal für andere Stellen protokolliert;
*Wächter ohne Fangvermögen* tritt zweimal auf (H-2, M-2) und ist BEO-023 in seiner
sechsten Instanz.
