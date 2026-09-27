# Review-Report — slice-232 (`links`: Ziel hinter dem Zeilenumbruch nach `](`), R1

**Review-Art:** Code (gegen Slice-Plan, betroffene `DC-*`-Anforderung, ADR-0091, Hard Rules).

**Gegenstand:** `slice-232`, Commit-Range `c97fa95a..377311d8` (drei Commits:
`45d55827` Lifecycle-Move, `4e07cf4c` Vertrag/ADR, `377311d8` Feature-Code).

**Skill:** `.harness/skills/reviewer.md` @ `8dcc0452` (Version 1.16.0, 2026-09-07).

**Modell-ID:** `claude-sonnet-5`.

**Datum:** 2026-09-27.

**Eingangs-Kontext:**

- Slice-Plan `docs/plan/planning/in-progress/slice-232-links-ziel-hinter-zeilenumbruch.md`
  (§1–§8, vollständig gelesen).
- `DC-FA-LINK-001` (`spec/lastenheft.md`), Spec-Stelle `DC-FA-LINK-001.a`
  Schritt 3 (`spec/spezifikation.md`).
- `docs/plan/adr/0091-links-absatzweise-zeilenumbruch-extraktion.md` (vollständig
  gelesen: Form und Inhalt geprüft).
- `AGENTS.md` §3 (§3.7 Kommentar-Klassen, §3.8 Modul-Scan-Grenze), §5
  (Grenzen-Regel, MR-025, Zitat-Geltungsbereich, Commit-Overclaim-Regel).
- `docs/reviews/` nach früheren Reviews zu `markdown.go`/`ExtractLinks`/
  `DC-FA-LINK-001` durchsucht: keine themengleiche Vorgänger-Review gefunden;
  kein Wiederholungs-Finding aus einem früheren Lauf zu übernehmen.
- **Nicht** erhalten (laut Skill): die DoD-Abhakung selbst.

**Gates unabhängig nachgefahren** (Docker/`make`, kein Host-Go — AGENTS.md §3.1):

| Gate | Ergebnis |
|---|---|
| `make test` (voller `go test ./...` in Docker) | grün — alle Pakete `ok` |
| `make gates` (voll, alle zehn Glieder) | grün — `baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; Coverage `94.30%` (Schwelle `93%`); `arch-check`/`semgrep`: 0 Befunde; `doc-check`: `838 Datei(en) geprüft, 0 Befund(e)` |
| Runtime-Image gegen dieses Repo (`docker run --rm --network none -v $(pwd):/repo:ro d-check:latest`) | `838 Datei(en) geprüft, 0 Befund(e)` — reproduziert die Commit-Botschafts-Zahl unabhängig |

---

## Zentraler Prüfpunkt: Blast-Radius und Korrektheit der Absatz-Faltung

**Ergebnis: Der Bracket-Matching-Blast-Radius ist real, nicht mathematisch
ausgeschlossen — und praktisch nachgewiesen (nicht nur konstruiert).** Die
Bestandsmessung, mit der die Commit-Botschaft und die ADR die Sicherheit der
Änderung belegen, deckt zwei der sechs Konsumenten von `ExtractLinks` nicht ab
und hätte eine reale Verhaltensänderung im Bestand dieses Repos übersehen, wäre
sie nicht schädlich gewesen.

### 1. Bracket-Matching über Zeilengrenzen — nachgewiesenes Risiko

`matchBracket` (Zeile ~433) zählt `(`/`)` bzw. `[`/`]` unabhängig von
Zeilengrenzen, sobald `ExtractLinks` (Zeile 541) die Zeilen eines Absatzes zu
einem String zusammenfügt. Per temporärem Testfall (`go test` in Docker,
danach wieder entfernt — Arbeitsbaum ist clean, `git status` bestätigt) wurden
drei Fälle verifiziert:

1. **Literales unbalanciertes `(` ohne `[`** (z. B. „Der Satz (mit Klammer
   ohne Schluss" gefolgt von einer Zeile mit echtem Link) — **folgenlos,
   bestätigt**: `matchBracket` für `(`/`)` wird nur ab der Position direkt
   nach einem geschlossenen `](` aufgerufen (`parseLinkAt`, Zeile ~583); ein
   literales `(` ohne vorausgehendes `[…]` triggert `parseLinkAt` nie.
2. **Unbalanciertes `[` ohne `]` auf seiner Zeile, gefolgt von einer Zeile mit
   `x](target.md)`** — **löst tatsächlich einen Fehlfund aus**: das
   `ExtractLinks`-Ergebnis enthält `{Line: 1, Target: "target.md", Text:
   "bracket without a closer\nx"}` — der Fund wird der **Zeile des
   unbalancierten `[`** zugeschrieben, nicht der Zeile, auf der `](` tatsächlich
   steht. Für einen realen `target-missing`-Befund hieße das: der Befund zeigt
   auf die falsche Zeile.
3. **Ein echter, sauberer Link gefolgt von einer weiteren, für sich
   unbalancierten `[baz](Zeile 1) / qux](real.md)`-Konstruktion** — der
   fehlerfreie Link (`a → b.md`) bleibt unberührt (Isolations-Eigenschaft
   bestätigt), **aber** die zweite, in der Quelle gar nicht als
   Link gemeinte Zeichenfolge wird zu einem **erfundenen** Link fusioniert
   (`Text: "baz\nqux"`, `Target: "real.md"`, `Line: 1`) — eine Zeichenfolge,
   die unter der bisherigen (zeilenweisen) Logik **kein** Link gewesen wäre,
   weil weder Zeile 1 noch Zeile 2 für sich eine balancierte `[...](...)` -Form
   trägt.

Das ist eine andere, breitere Fehlerklasse als das in ADR-0091 und im
Slice-Plan §1 benannte „Linktext über Zeilenumbruch" (dort: ein Mensch würde
`[lang⏎text](ziel)` noch als *einen* beabsichtigten Link lesen, nur mit
ungewöhnlichem Soft-Break). Hier wird **Prosa, die kein Mensch als Link läse**
(ein einzelnes, für sich stehendes `[` in einem Satz, gefolgt von unbezogenem
Text mit einer späteren `]…(` -Sequenz im selben Absatz), durch die
Absatz-Faltung zu einem Fantom-Link oder einer Falsch-Zuordnung. Das ist weder
in §1 „Ausdrücklich NICHT" noch in ADR-0091 als Grenze benannt — die ADR
benennt nur die eng verwandte, aber andere Grenze „Zeilenumbruch im
Linktext/vor Titel". `TestExtractLinks_Kanten` (Bestandstest,
`internal/hexagon/core/rules/markdown_test.go:139`) enthält zufällig zwei
Zeilen mit unbalancierten Klammern gefolgt von echten Links im selben Absatz
und besteht weiterhin — aber dort ist die spätere `]`-Position, die den
unbalancierten Bracket schließen könnte, nie unmittelbar von `(` gefolgt,
weshalb der Fall dort nicht auftritt. Es ist Zufall, kein Test, der diese
Fehlerklasse gezielt abdeckt (siehe R1-H1 unten).

### 2. Bestandsmessung nachgefahren — bestätigt die Zahl, aber deckt nicht alle Konsumenten

- Neues Image (`377311d8`) gegen dieses Repo: `838 Datei(en) geprüft, 0
  Befund(e)` — reproduziert.
- Altes Image (Stand `c97fa95a`, per `git worktree`) gegen denselben
  (aktuellen) 838-Datei-Bestand: ebenfalls `838 Datei(en) geprüft, 0
  Befund(e)` — für die **in `.d-check.yml` aktiven** Module (`links, anchors,
  ids, matrix, codepaths, spans, hostpaths, versions, structure, diagrams,
  citations`) tatsächlich verhaltensgleich.
- **Aber:** `.d-check.yml`s `modules:`-Liste enthält **weder `external` noch
  `tracked`** — zwei der sechs `ExtractLinks`-Konsumenten laut ADR-0091. Mit
  `--enable external --enable tracked` (Netz aus, daher nur zur
  Fund-**Zahl**-Differenz aussagekräftig) ergibt der Vergleich Alt- vs.
  Neu-Binary auf demselben Bestand: **107 vs. 108 Befunde** — eine
  tatsächliche neue Zeile: `tools/archive-wave/README.md:3
  https://github.com/pt9912/ai-harness-course external-status`. Prüfung der
  Datei bestätigt: Zeile 3/4 tragen exakt die neue Form (`[Text,⏎ Fortsetzung]
  (https://…)`), vorher vom `external`-Modul gar nicht als Link erkannt, jetzt
  korrekt erkannt — in diesem Einzelfall ein **echter, harmloser** Treffer
  der neuen Fähigkeit, kein Fehlfund. Das ändert aber nichts daran, dass die
  Commit-Botschaft „838 Dateien, 0 Befunde vor und nach der Änderung" und die
  DoD-Formulierung „Befundzahl der Konsumenten-**Module**" (Plural) einen
  Vollständigkeitsanspruch über alle sechs `ExtractLinks`-Konsumenten
  suggerieren, den die tatsächlich gelaufene Messung nicht einlöst: zwei der
  sechs blieben ungemessen, und mindestens einer davon zeigt nachweislich
  *irgendeine* Verhaltensänderung auf dem eigenen Bestand. Für Konsumenten wie
  `external`/`tracked`, die keine `target-missing`-Zahl, sondern andere
  Aussagen (Erreichbarkeit, Tracked-Status) aus demselben `LinkRef` ziehen,
  wäre eine durch Bracket-Fusion **falsch attribuierte** Zeile durch eine reine
  Zählmessung ohnehin nicht sichtbar geworden (Befund-**Zahl** bliebe gleich,
  Befund-**Ort** könnte trotzdem falsch sein) — siehe Punkt 1.

### 3. Vier „unveränderte" String-Konsumenten — Code-Check

- `ids.go:92,132` und `internal/hexagon/core/app/repair.go:113` rufen
  tatsächlich `ExtractLinkSpans` (unverändert) auf — bestätigt.
- `internal/hexagon/core/rules/planning_observations.go:180` ruft
  `ExtractLinkSpans` auf — bestätigt.
- `internal/hexagon/core/rules/pins.go:70` ruft **nicht** `ExtractLinkSpans`,
  sondern **direkt** `forEachLink` (zeilenweise, je Zeile ein Aufruf in
  `CheckPins`, Zeile 29 ff.) — die ADR-Kontext-Sektion ordnet `pins` fälschlich
  der `ExtractLinkSpans`-Signatur zu. Die **Schlussfolgerung** (Verhalten für
  `pins` bleibt zeilenbasiert/unverändert) ist trotzdem richtig, weil
  `forEachLink` selbst nicht angerührt wurde und `pins` es weiterhin nur mit
  einer einzelnen Zeile aufruft — siehe R1-L1.

### 4. Line-Zuordnung, Randfall `<=` statt `<` — kein Bug, verifiziert

`linkParagraphs` (Zeile 513) hängt eine Zeile nur an `cur` an, wenn
`!blank` (Zeile 523) — eine Leerzeile wird nie Mitglied einer `grp`.
Damit ist `len(raws[i])` für jedes Gruppenmitglied ≥ 1, `lineEnd` liegt
exakt auf der Position des Trenn-`\n` (bzw. bei der letzten Zeile der Gruppe
außerhalb des gültigen Indexbereichs). `parseLinkAt` startet nur an
Positionen mit `s[i] == '['` oder `'!'` — niemals an einer `\n`-Position und
niemals außerhalb `len(text)`. Der Vergleich `span.Start <= lineEnd`
ist damit zwar technisch weiter gefasst als `<`, aber inert: `span.Start`
kann den Wert `lineEnd` nie annehmen. Verifiziert am Code, kein Finding.

### 5. Fund-Zeilen-Konsistenz — Vertrag/ADR/Code deckungsgleich

Lastenheft („Der Fund wird der öffnenden Zeile des Links zugeschrieben"),
Spezifikation Schritt 3 (dieselbe Formulierung), ADR-0091 Abschnitt (b) und
der Code (`ref.Line = grp[i].No` an der Position von `span.Start`, die stets
mit der Position von `[`/`![` übereinstimmt) sind konsistent. Bestätigt durch
`TestCLI232_ZeilenumbruchHinterKlammer_TotesZiel` (erwartet Zeile 1, nicht
Zeile 2).

---

## Weitere Prüfpunkte

- **`TestCLI232_UnveraenderteKontrollformen`:** Fixture platziert alle fünf
  Formen (Inline, Titel, Spitzklammer, Klammer-im-Ziel, Zeilenumbruch)
  tatsächlich zeilenweise ohne Leerzeile dazwischen (`cli_acceptance_test.go`,
  Zeilen 2-6 der Fixture) — der Test sagt also wirklich etwas über
  Absatz-Interferenz zwischen **wohlgeformten** Linkformen aus. Er deckt aber
  **nicht** den unter „Zentraler Prüfpunkt" Punkt 1 gezeigten Fall
  unbalancierter/gar nicht als Link gemeinter Klammer-Zeichen im selben
  Absatz (siehe R1-H1).
- **Abgrenzung §1 „Ausdrücklich NICHT":** Referenz-Definitionen,
  Linktext-Zeilenumbruch, Titel-hinter-Zeilenumbruch, `--repair` für die neue
  Form — im Diff nicht angefasst (`repair.go` unverändert, keine
  `[name]:`-Logik berührt). Handbuch/README/CHANGELOG/Release: keine dieser
  Dateien im Diff. Eingehalten.
- **Commit-Zerlegung:** `377311d8` (Feature) ändert ausschließlich
  `internal/hexagon/core/rules/markdown.go` und
  `internal/adapter/driving/cli/cli_acceptance_test.go` — rein `internal/`.
  `4e07cf4c` (Vertrag) ändert ausschließlich `docs/plan/adr/`, `spec/*.md` —
  fasst keinen Code an. `45d55827` ist ein reiner Lifecycle-Move (git mv +
  zwei erlaubte Pfad-Nachzüge in fremden Slice-Dateien, AGENTS.md §3.3).
  Sauber.
- **ADR-0091 Form:** drei verglichene Alternativen mit Trade-offs, Fitness
  Function (Tabelle), `Re-Evaluierungs-Trigger` (zwei benannte Bedingungen),
  `Schärft:` zeigt aufwärts auf `DC-FA-LINK-001` — vollständig.
- **§3.4 Referenzrichtung:** `grep -n "ADR-0091\|slice-232"
  spec/lastenheft.md spec/spezifikation.md` liefert keinen Treffer — beide
  Spec-Straten nennen ADR/Slice nicht im Körper, `matrix` bestätigt das über
  `make doc-check`.
- **MR-032/Decken-Regel:** Historie-Zeile 0.92.0 nennt weder ADR-0091 noch
  slice-232, `Verweis`-Spalte `—`. Eingehalten.

---

## Findings

| # | Kategorie | Quelle | Pfad | Befund | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| R1-H1 | HIGH | `DC-FA-LINK-001.a` Schritt 3 / Prüffrage 2 (Kern-Modul meldet falsch) | `internal/hexagon/core/rules/markdown.go:541` (`ExtractLinks`), Fehlender Test in `internal/adapter/driving/cli/cli_acceptance_test.go` | Die absatzweise Faltung lässt ein für sich unbalanciertes `[` (kein Link, gewöhnliche Prosa, kein `]` auf seiner eigenen Zeile) mit einer späteren, unabhängigen `]…(...)`-Sequenz im selben Absatz zu einem Fund verschmelzen — nachgewiesen per temporärem Testfall: `ExtractLinks` liefert für zwei Prosa-Fragmente, von denen keines für sich ein Link ist, einen erfundenen `LinkRef` mit Text über die Zeilengrenze hinweg **und/oder** ordnet einen tatsächlich existierenden Link der falschen (der vorausgehenden) Zeile zu. Das ist eine andere, breitere Fehlerklasse als die in ADR-0091/Slice-Plan §1 benannte Grenze „Zeilenumbruch im Linktext" (dort ist der Umbruch Teil eines vom Menschen als Link gelesenen Ausdrucks) und ist in keinem der beiden Verträge oder in ADR-0091 als Grenze benannt. Betroffen sind potenziell alle sechs `ExtractLinks`-Konsumenten (`links`, `anchors`, `matrix`, `external`, `tracked`, `resolve-from`); für `target-missing`-artige Befunde bedeutet das eine mögliche Falsch-Zuordnung der gemeldeten Zeile, für andere Konsumenten (z. B. `matrix`-Lineage-Match über `Text`) eine potenziell falsche Zuordnung des Linktexts. Auf dem aktuellen 838-Datei-Bestand tritt der Fall nicht auf (Bestandsmessung 0 Befunde), das schließt ihn aber nicht für künftige oder fremde Dokumente aus. | ja — Fixture mit einem freistehenden `[` in Prosatext (keine eigene schließende `]` in derselben Zeile), gefolgt im selben Absatz von einer Zeile mit `]`+`(...)` unmittelbar danach; `ExtractLinks` darauf ausführen und Line-/Text-Feld des Ergebnisses prüfen | link-paragraph-bracket-fusion |
| R1-M1 | MEDIUM | `AGENTS.md` §5 (Commit-Botschaft-Overclaim) / Prüffrage 8 | Commit `377311d8` (Botschaft), DoD-Punkt „Bestandsmessung" im Slice-Plan §2, ADR-0091 Abschnitt „Konsequenzen" | Die Botschaft „Bestandsmessung gegen dieses Repo: 838 Dateien, 0 Befunde vor und nach der Änderung" und die DoD-Formulierung „die Befundzahl der Konsumenten-Module gegen den Stand vor der Änderung" lesen sich als Aussage über **alle** Konsumenten-Module der geänderten Funktion. Tatsächlich lief die Messung nur über die in `.d-check.yml` aktiven Module (`links, anchors, ids, matrix, codepaths, spans, hostpaths, versions, structure, diagrams, citations`) — zwei der sechs `ExtractLinks`-Konsumenten (`external`, `tracked`) sind darin nicht enthalten. Ein unabhängiger Nachlauf mit `--enable external --enable tracked` zeigt auf demselben Bestand tatsächlich eine Verhaltensänderung (eine zusätzliche `external-status`-Zeile durch einen jetzt erkannten, zuvor unerkannten Link in `tools/archive-wave/README.md:3`) — im konkreten Fall harmlos (echter, jetzt korrekt erkannter Link), aber ein Beleg, dass „0 Befunde vor und nach" für die tatsächlich gelaufene Messmenge gilt, nicht für die im Text suggerierte volle Konsumentenmenge. | ja — `docker run --rm --network none -v $(pwd):/repo:ro <image> --enable external --enable tracked` mit Alt- und Neu-Image vergleichen | messmethode-deckt-nicht-alle-konsumenten |
| R1-L1 | LOW | Doku-Präzision, ADR-0091 Kontext | `docs/plan/adr/0091-links-absatzweise-zeilenumbruch-extraktion.md` (Abschnitt „Kontext", Aufzählung der `ExtractLinkSpans`-Konsumenten) | Die ADR ordnet `pins` der Signatur `ExtractLinkSpans(text string) []LinkSpan` zu („vier Konsumenten: … `pins` (`dpin`-Marker-Bindung …)"). Tatsächlich ruft `internal/hexagon/core/rules/pins.go:70` `forEachLink` direkt auf, nicht `ExtractLinkSpans`. Die Schlussfolgerung (verhaltensunverändert, weiterhin zeilenbasiert) bleibt richtig, weil `forEachLink` selbst nicht geändert wurde und `pins` es weiterhin je Zeile aufruft — die technische Zuordnung in der ADR ist trotzdem ungenau. | ja — `grep -n "ExtractLinkSpans\|forEachLink" internal/hexagon/core/rules/pins.go` | adr-attribution-ungenau |
| R1-H2 | HIGH | `AGENTS.md` §3.7 (Kommentar-Klassen, „keine Slice-Nummern") | `internal/adapter/driving/cli/cli_acceptance_test.go` (vier neue Kommentare vor `TestCLI232_ZeilenumbruchHinterKlammer_TotesZiel`, `_LebendesZiel`, `_Bild`, `TestCLI232_UnveraenderteKontrollformen`) | Alle vier neuen Test-Kommentare beginnen mit `// slice-232 …`. §3.7 verbietet explizit „Slice-Nummern" im Kommentar und lässt als Herkunfts-Feld nur `DC-*`, `ADR-*`, `MR-*` oder `seit welle-<NN>`/`seit slice-<Kennung>` (für verkörperte Steering-Loop-Regeln) zu — eine bloße `slice-232`-Kennzeichnung eines Testzwecks fällt unter keine der fünf erlaubten Klassen. Die Bestandsgrenze („Test-Kommentare … sind grandfathered") deckt nur **vor** der Einführung der Regel geschriebene Kommentare; „Neuzugänge fallen überall unter den Anker" — diese vier sind Neuzugänge dieses Diffs. Dasselbe Muster existiert bereits an anderer Stelle der Datei (z. B. `// slice-076 …`), was die Verbreitung erklärt, aber die Neuzugänge nicht grandfathered macht. | ja — Diff von `cli_acceptance_test.go` zeigt die vier neuen `// slice-232`-Kommentare als reine Neuzugänge | slice-nummer-in-kommentar |

## Negativbefunde (geprüft, ohne Befund)

- **Gates:** `make test` und der volle `make gates`-Lauf (zehn Glieder)
  unabhängig nachgefahren — alle grün, echte Ausgabe oben zitiert
  (Coverage 94,30 %, `arch-check`/`semgrep` je 0 Befunde).
- **Hexagon-Import-Richtung (ADR-0005):** keine neuen Imports in
  `markdown.go`; Diff berührt keine Adapter-Schicht.
- **Gate-Suppression / Schwellen-Senkung ohne ADR (§3.2/§3.6):** keine
  `//nolint`, keine Schwellen-Änderung im Diff.
- **Netzzugriff außerhalb `external`:** keiner im Feature-Commit.
- **Zustandsfelder:** keine Roadmap-/Register-Zustandszeile inhaltlich
  verändert außer dem bereits etablierten Lifecycle-Move-Begleiteffekt
  (`Nichts in Arbeit` entfernt).
- **Fund-Zeilen-Konsistenz (Vertrag/ADR/Code):** deckungsgleich, siehe
  „Zentraler Prüfpunkt" Punkt 5.
- **Line-Zuordnungs-Randfall `<=`:** verifiziert kein Bug, siehe „Zentraler
  Prüfpunkt" Punkt 4.
- **Drei String-Konsumenten (`ids`, `--repair`, `planning`):** rufen
  tatsächlich `ExtractLinkSpans` unverändert auf, wie in ADR-0091 behauptet.
- **Abgrenzung §1 „Ausdrücklich NICHT":** alle fünf Punkte eingehalten
  (Referenz-Definitionen, Linktext-Umbruch, Titel-Umbruch, `--repair`,
  Handbuch/README/CHANGELOG/Release — keiner im Diff berührt).
- **Commit-Zerlegung:** rein — Feature-Commit nur `internal/`, Vertrags-Commit
  ohne Code, Lifecycle-Move ohne Inhaltsänderung an der bewegten Datei.
- **ADR-0091 Form (Modul 4):** drei Alternativen, Fitness Function,
  Re-Evaluierungs-Trigger, `Schärft:` aufwärts — vollständig.
- **§3.4/MR-006 Referenzrichtung:** kein ADR-/Slice-Token im Körper von
  Lastenheft oder Spezifikation.
- **MR-032 Decken-Regel:** Historie-Zeile ohne ADR-/Slice-Nennung.
- **Bestandsmessung (Zahl):** „838 Dateien, 0 Befunde vor und nach" für die
  tatsächlich gelaufene Modul-Menge unabhängig reproduziert — korrekt für die
  gemessene Menge (siehe R1-M1 für die Reichweiten-Einschränkung).
- **Rot-Beleg gegen den alten Stand:** unabhängig nachgefahren (alte
  `markdown.go` per `git show c97fa95a:…` eingesetzt, `make test` in Docker
  gebaut) — exakt drei Tests schlagen fehl
  (`TestCLI232_ZeilenumbruchHinterKlammer_TotesZiel`,
  `_Bild`, `TestCLI232_UnveraenderteKontrollformen`), wie die Commit-Botschaft
  behauptet; `_LebendesZiel` besteht in beiden Ständen (erwartungsgemäß, kein
  Diskriminator).
- **Referenz-Richtung/Provenance-Marker:** zwei `<!-- d-check:status-provenance
  -->`-Marker in ADR-0091 zeigen jeweils nur, *wo* die Entscheidung entstand
  (Auftraggeber-Entscheid vom Datum, historischer Vorfall slice-073) — keine
  getarnte Entscheidungsgrundlage.
- **`DC-QA-04` Alt-Tool-Migrationsabdeckung:** nicht berührt.
- **Modul-Scan-Grenze (§3.8):** nicht einschlägig — `ExtractLinks` ist eine
  interne Parser-Funktion, kein scannendes Audit-Modul mit eigener
  Scan-Wurzel.

## Kategorie-Summary

- HIGH: 2 (R1-H1, R1-H2)
- MEDIUM: 1 (R1-M1)
- LOW: 1 (R1-L1)
- INFO: 0

## Verdikt

**Merge-blockierend: ja, wegen R1-H1 und R1-H2.**

Begründung: R1-H1 ist ein nachgewiesener (nicht nur hypothetischer)
Korrektheits-Blast-Radius in einer von sechs Modulen konsumierten
Kern-Funktion, der über die in ADR-0091/Slice-Plan §1 explizit benannten
Grenzen hinausgeht und durch keinen bestehenden oder neuen Test abgedeckt ist
— die Bestandsmessung mit 0 Befunden ist dafür kein hinreichender Beleg
(Fusion zu einem Fantom-Link ändert nicht notwendig eine Fund-**Zahl**, kann
aber eine Fund-**Zeile**/einen Fund-**Text** verfälschen, was ausgerechnet in
den von der Standardmessung nicht erfassten Modulen `external`/`tracked`
unbeobachtet bliebe, siehe R1-M1). Vor Merge sollte mindestens eine
Absatz-Boundary-Regel ergänzt werden, die verhindert, dass ein Bracket ohne
eigene schließende Klammer über eine Zeilengrenze hinweg matcht, **oder** die
Grenze muss explizit in ADR-0091/Lastenheft benannt und mit einem
Negativtest belegt werden. R1-H2 ist mit vier Ein-Wort-Präfix-Entfernungen
trivial zu beheben (die Testabsicht bleibt vollständig erhalten, wenn
`slice-232 ` am Zeilenanfang entfällt), blockiert aber laut Reviewer-Skill
formal als HIGH, weil §3.7 die Klasse „Slice-Nummer im Kommentar" explizit
und ohne Ausnahme für Neuzugänge nennt. R1-M1 sollte vor Closure durch eine
präzisere Formulierung der Botschaft/DoD (Nennung der tatsächlich
gemessenen Modul-Menge) aufgelöst werden, blockiert für sich allein nicht
zwingend. R1-L1 ist nicht merge-blockierend (Doku-Präzision, keine
Verhaltensfrage).
