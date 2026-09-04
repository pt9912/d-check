# Review slice-171 — Der Slice erfüllt die von ihm eingeführte Regel in einem von zwei Fällen selbst nicht

**Gegenstand:** [slice-171](../plan/planning/done/slice-171-vorpruefungen-belegen.md), Stand `1fdc0ee` (Plan) + `4ef5025` (Implementierung).
**Datum:** 2026-08-29. **Reviewer:** unabhängiger Subagent, Skill `.harness/skills/reviewer.md` v1.13.0.
**Eigene Läufe:** `make doc-check` ⇒ 589 Dateien, 0 Befunde, Exit 0. Produkt-Gegenprobe (kein `grep`) über die drei neu ausgeschlossenen Verzeichnisse mit `modules: [citations]` und ohne `ignore` ⇒ 455 Dateien, 0 Befunde, Exit 0.

---

## Urteil

**Blockiert.** HIGH 1 · MEDIUM 4 · LOW 2 · INFO 1. Alle Befunde sind
Doku-/Konfigurationsform — keiner betrifft Produktcode, `doc-check` ist grün.
Am schwersten wiegen F-2 und F-3: F-2, weil der Slice die von ihm eingeführte
Regel in einem von zwei Fällen selbst nicht erfüllt und als Kopiervorlage jedes
künftigen Plans wirkt; F-3, weil eine falsche Prämisse im Plan ein drittes
Verzeichnis in den Ausschluss getragen hat, für das die gegebene Begründung
nicht gilt.

## HIGH

### F-1 · `kommentar-traegt-herkunfts-prosa`

- **quelle:** AGENTS.md §3.7 / Baseline §Was ein Kommentar trägt / Reviewer-Skill Frage 6
- **pfad:** `.d-check.yml:636`, `:647-654`
- **befund:** Der neue `citations`-Kommentarblock trägt (a) eine Slice-Nummer
  (`WARUM: Seit slice-171 …`), (b) Mess-Labels und die eigene
  Vorgangs-Geschichte (`GEMESSEN vor dem Setzen: …`, `587 -> 653 Dateien beim
  ersten Versuch` — also den **verworfenen ersten Versuch**), (c) über vier
  Absätze sechs Herkunfts-Kennungen. §3.7 nennt „keine Slice-Nummern und keine
  Mess-Labels"; die Baseline nennt Deliberation und Herkunfts-Prosa als zwei der
  drei Klassen, die keine der fünf tragen, und lässt Herkunft nur als **ein**
  auflösbares Feld zu, nie als Absatz. Neuzugang, also keine Bestandsgrenze; der
  direkt darüberstehende `workflows:`-Block zeigt die Form ohne diese Anteile.
  Die Kopplungs-Zusage im Block ist legitim — der Befund gilt den genannten
  Zeilen, nicht dem ganzen Kommentar.
- **verifizierbar:** nein — §3.7 sagt selbst „Kein Gate prüft das".

## MEDIUM

### F-2 · `cite-spanne-belegt-nicht-die-benannte-regel`

- **quelle:** MR-054 (Geltungsbereich), BEO-012, Reviewer-Skill Frage 9
- **pfad:** `docs/plan/planning/done/slice-171-vorpruefungen-belegen.md:165-169`
- **befund:** Der zweite Vorprüfungs-Block schreibt „Die Regel, die diesen
  Schritt vorschreibt:" und zitiert `modul-05-planning-harness.md:224-225` —
  eine **Notier-Pflicht**, nicht die Vorschrift des Schritts. Die vorschreibende
  Zeile ist **219**. MR-054 definiert die Form als Direktive auf die Zeilen, die
  den Schritt vorschreiben; der Slice, der die Regel einführt, erfüllt sie in
  einem von zwei Fällen nicht — und dieser Plan ist die Kopiervorlage jedes
  künftigen. Zwei Zeilen darüber benennt derselbe Block BEO-012 als die Klasse,
  die eine `cite`-Spanne angeblich eingrenzt. Folge: verschwindet beim Bump der
  „Keine Treffer"-Satz, entfällt die Direktive nach MR-039, obwohl der Schritt
  weiter vorgeschrieben ist; verschwindet umgekehrt Zeile 219, bleibt der Block
  grün.
- **verifizierbar:** teilweise — `doc-check` bleibt grün (das Zitat *ist*
  wortgleicher Teilstring der Spanne); die Diskrepanz Etikett↔Spanne ist Urteil.

### F-3 · `ausschluss-begruendung-deckt-nur-teil-der-menge`

- **quelle:** MR-054, Reviewer-Skill Frage 1/9
- **pfad:** `.d-check.yml:662`, `harness/conventions/MR-054-vorpruefungen-belegen-ihre-regel.md:56-62`, Slice §2 Punkt 3
- **befund:** Die Begründung lautet „Der Beleg zählt zum Zeitpunkt seiner
  Prüfung: vor dem `git mv` läuft die Direktive im inneren Loop; danach ist sie
  Lauf-Beleg." Das trifft `done/` (Live-Phase in `in-progress/`) und
  `conventions/done/` (Live-Phase in `conventions/`). **`docs/reviews/` hat
  keine solche Phase** — ein Report entsteht direkt dort. Eine `cite`-Direktive
  in einem Review-Report wird damit **nie** geprüft, und eine *malformte* nimmt
  den Lauf nicht mehr fail-closed mit; vor diesem Commit tat sie das. Wurzel ist
  eine falsche Prämisse im Plan: „nimmt `done/` aus, **wie es `docs/reviews/`
  faktisch schon ist**" — das gilt für `codepaths`/`ids`/`versions`
  (`exempt-paths`), für `citations` galt es nicht.
- **verifizierbar:** ja — `git show 4ef5025^:.d-check.yml | grep citations:` ⇒
  leer. Die Produkt-Gegenprobe zeigt zugleich, dass **rückwirkend nichts**
  stillgelegt wurde (0 Befunde), der Ausschluss also rein vorwärts wirkt.

### F-4 · `praezedenz-ohne-achsen-differenz`

- **quelle:** BEO-012 / Reviewer-Skill Frage 9
- **pfad:** `harness/conventions/MR-054-…md:60-62`, `.d-check.yml:643-645`, dieselbe Aussage in der Commit-Botschaft
- **befund:** „Die **Link**-Achse löst dasselbe längst so — `ignore-refs` nimmt
  die eingefrorenen Verweise … quell-skopiert aus." `ignore-refs` ist
  quell-skopiert **und ziel-benannt**; die Datei bleibt im Scan. Der
  repo-eigene Kommentar sagt das ausdrücklich (`.d-check.yml:142`).
  `citations.scope.ignore` nimmt dagegen die **ganze Datei** aus dem Modul. Die
  Präzedenz stützt die Richtung, nicht die Reichweite; „dasselbe" überdehnt sie.
- **verifizierbar:** ja — `.d-check.yml:139-150`.

### F-5 · `regel-in-kraft-ohne-zustellung`

- **quelle:** MR-054 §Geltungsbereich; AGENTS.md §5
- **pfad:** `harness/conventions/MR-054-…md:9-13`, `:21-27`
- **befund:** MR-054 gilt ab sofort für jeden neu angelegten Slice-Plan, und
  kein Sensor prüft es — ein Plan ohne Direktiven ist grün. Einziger Trägerort
  ist der Index in `harness/conventions.md`; `AGENTS.md` §5, das die drei
  Vorprüfungen aufzählt und für die dritte MR-053 nennt, sagt zur
  Direktiven-Form nichts, und die Baseline-Vorlage kann sie nicht tragen
  (SHA-gepinnt, `make baseline-verify` in `gates`). Das ist dieselbe Klasse —
  Regel ohne Träger —, gegen die der Slice gebaut ist. **Mildernd:** slice-176
  ist genau diese Arbeit, liegt aber in `open/` und ist unbeansprucht, während
  MR-054 bereits gilt.
- **verifizierbar:** nein (kein Gate).

## LOW

### F-6 · `paraphrase-in-anfuehrungszeichen`

`harness/conventions.md:134`: Die MR-054-Indexzeile schreibt *„in `done/` steht
keine"* in Zitatform und schreibt es MR-051 zu. MR-051 sagt: *„Nicht `done/`,
nicht `docs/reviews/`, nicht `conventions/done/` — dort steht keine."*
`grep -cF` gegen MR-051 ⇒ **0**. Der MR-Körper markiert seine (korrekte)
Wiedergabe an derselben Stelle ausdrücklich als „wörtlich" — die Indexzeile
benutzt dieselbe Typografie für eine Paraphrase. In einem Slice über wortgleiche
Belege ist das die eigene Klasse. **verifizierbar:** ja.

### F-7 · `mr-ohne-vorwaerts-zeiger`

`harness/conventions/MR-051-cite-spannen-beim-bump.md:9-12` und die Indexzeile
bleiben unverändert, obwohl MR-054 ihre Feld-Aussage für absehbar unwahr
erklärt. AGENTS.md §5/BEO-012 verlangt, vor jedem Verweis genau dieses
`Geltungsbereich`-Feld zu lesen. Szenario: eine spätere Aufräumung liest MR-051,
hält den `done/`-Ausschluss für unnötig, entfernt ihn — und der innere Loop wird
beim nächsten Bump fail-closed rot. **Nicht höher eingestuft**, weil MR-021
gegenüber MR-051 denselben fehlenden Rückzeiger hat: der Index ist der Lesepfad
des Hauses. **verifizierbar:** ja.

## INFO

### F-8

MR-054 zitiert MR-051 wörtlich in Prosa **ohne** `d-check:cite`-Direktive,
obwohl das Ziel repo-lokal und auflösbar ist und der MR selbst den verifizierten
Beleg zum Thema hat. Von MR-054s eigenem Geltungsbereich nicht gefordert —
deshalb INFO, kein Verstoß.

## Negativbefunde (geprüft, ohne Befund)

- **Messung „keine wirksame Direktive in den drei Verzeichnissen"** — mit dem
  **Produkt** gegengeprüft: 455 Dateien, `modules: [citations]`, kein `ignore`
  ⇒ 0 Befunde, Exit 0. Die 13 rohen Marker-Treffer sind vollzählig (12
  Inline-Code, 1 im Fence `done/slice-152:41`). MR-051s „dort steht keine" war
  richtig.
- **Modul-lokaler Scope ersetzt den globalen** — `scan.ignore` führt genau die
  zwei `.harness/`-Einträge; beide sind wiederholt. Es fehlt kein dritter.
- **Ziel-Seite unberührt** — der `.harness/baseline/**`-Eintrag im Modul-Scope
  bricht die Skill-Direktiven **in** die Baseline nicht; `reviewer.md:191`
  bleibt in der Scan-Menge.
- **Beide Zitate wortgleich** — whitespace-normalisierte Teilstrings ihrer
  Spannen, über der 16-Zeichen-Schwelle.
- **Kein Widerspruch zu ADR-0045/ADR-0060** — ADR-0060 Entscheidung 6 lässt die
  zwei groben Achsen ausdrücklich stehen; ihre Ablehnung der groben Achse ist
  auf **lebende** Dokumente skopiert.
- **AGENTS.md §3.6** — kein Befund: die gelebte Präzedenz legt
  Dogfood-Scope-Ausnahmen ohne ADR in die Config; ADRs decken den
  Produkt-Mechanismus, nicht die Repo-Scope-Wahl.
- **`.d-check.closure.yml`** unberührt — der Closure-Bindepunkt fuhr `citations`
  nie.
- **Kein Go-Code, keine Imports, kein Netz, keine Inline-Suppression, keine
  Schwellen-Senkung** im Diff.
