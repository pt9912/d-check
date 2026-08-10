# Review-Report: slice-101 (dritte Runde) — 2026-08-09

**Review-Art:** Code — **dritte Runde**, klassen-orientiert. Erst-Review und
Re-Review waren beide merge-blockierend und fanden **dieselbe** Klasse an einer
neuen Stelle. Geprüft wird deshalb nicht, ob N-1 bis N-6 abgehakt sind, sondern
ob die Klasse **geschlossen** ist (Modul 10 §Drei Review-Arten). Der Dateiname
trägt den Datums-Stamm des Erst-Reports; der Lauf fand am 2026-08-10 statt.

**Gegenstand:** [slice-101](../plan/planning/in-progress/slice-101-fence-unbalanciert.md),
zweiter Heilungs-Commit `fc94503` (Diff-Range `fbb2d73..HEAD`); Gesamt-Slice
`38c36b2..HEAD`.

**Skill:** `.harness/skills/reviewer.md` @ 1.3.0 ·
**Modell:** claude-opus-5[1m] · **Datum:** 2026-08-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [Erst-Report](2026-08-09-slice-101-fence-review.md) (F-1 bis F-7) und
  [Re-Review-Report](2026-08-09-slice-101-fence-re-review.md) (N-1 bis N-6) als
  Checkliste
- [slice-101](../plan/planning/in-progress/slice-101-fence-unbalanciert.md)
  (§4 Abnahme-Punkte, §5 DoD, §6 Risiken, §8 Beobachtungs-Kandidat)
- [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md) (Entscheidungen 1–5,
  Alternativen, Konsequenzen, Fitness Function) und
  [ADR-0042](../plan/adr/0042-markdown-lexik-folgt-commonmark.md)
- [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
  (Klasse 3, Lastenheft-Version 0.52.2) und
  [§`DC-FA-SPAN-001.a`](../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
  Schritt 3 samt der §4-Grund-Code-Zeile
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (der Konsument, dessen Silent-Grün den Anlass gab — und der in diesem Commit
  sein Verhalten ändert)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules)

**Läufe dieses Reviews.** Alle Fixtures in einem Temp-Verzeichnis außerhalb des
Repos, alle Läufe netzlos und read-only. Gebaut wurden drei Images: HEAD, der
Stand **vor** dem zweiten Heilungs-Commit (`fbb2d73`, über einen git-Worktree im
Temp-Verzeichnis) und ein bewusst mutiertes für eine gezielte Sonde. Gefahren:
`make build`, `make test` (14 Mutationsläufe), `make gates`,
`make completeness-check`, `make verify-closure-notes` sowie rund 30
Fixture-Läufe gegen die Images. `make gates` grün (350 Dateien, 0 Befunde,
Coverage 94,30 %), `make completeness-check` 48 Anforderungen / 0 Waisen,
`make verify-closure-notes` 321 Dateien / 0 Befunde. Der Arbeitsbaum ist am Ende
unverändert (`git status --short` leer, Worktree entfernt).

---

## Vor-Befunde: Status

| Vor-Befund | Kat. | Status | Beleg (Lauf) |
|---|---|---|---|
| F-1 — Wächter wertet einen anderen Fence-Automaten aus als der Tabellen-Leser | HIGH | **geheilt** | Beide Lesarten sind ausgewertet und die strenge ist real geteilt: der Rückbau „strenge Lesart entfernt" ist rot (4 Fälle), „naive entfernt" rot, „Längen-Abgleich in `FenceCloses` entfernt" rot (3 Fälle). Lauf: ein Backtick-Block, den eine Tilden-Zeile „schließt", meldet `fence-unclosed` |
| F-2 — Wächter und Vorverarbeitung trimmen verschieden | HIGH | **geheilt im Code, ungesichert im Test** | Alle sieben Aufrufstellen von `FenceToggle`/`FenceRun`/`FenceCloses` werden aus `TrimFenceIndent` gespeist (Vollzählung unten). Aber drei der fünf Konsumenten haben keine Assertion — R-1 |
| F-3 — Grenze deckt nur die Quell-Achse | MEDIUM | **erweitert, weiterhin unvollständig** | Die sechs genannten Ziel-Achsen-Leser stimmen (je im Code belegt). Nicht genannt ist eine dritte Klasse: git-**Revisionen**, die das Modul `vcs` mit derselben fence-empfindlichen Section-Maske liest — R-7 |
| F-4 — Parität statt Paarung | MEDIUM | **geheilt bzw. benannt** | Legale Verschachtelung bleibt grün (Testfall vorhanden, Rückbau des Längen-Abgleichs macht ihn rot); die zweite Hälfte (eingerückter Code-Block) steht jetzt als eigene Grenze in der Anforderung |
| F-5 — zugesagtes Befund-Ziel ohne Assertion | MEDIUM | **für die sechs zugesagten Rückbauten bestätigt** | Selbst nachvollzogen über Dateikopien: sechs von sechs rot (Tabelle unten). Drei **weitere** Rückbauten bleiben grün — R-1 |
| F-6 — gemeldete Zeile ist nicht die Reparaturstelle | LOW | **im Vertrag ehrlich, im Entscheidungs-Protokoll nicht** | Anforderung, Spezifikation und ADR-Entscheidung 3 nennen sie **Fundstelle**. Slice §4 Punkt 3, die ADR-Alternativen-Tabelle und die ADR-Index-Zeile begründen weiter mit „Ort der Reparatur" — R-6 |
| F-7 — Sammelsatz nicht mit der dritten Klasse mitgewachsen | LOW | **geheilt** | Der Sammelsatz nennt das Ziel je Klasse; Lauf gegen eine lange Infozeile liefert die getrimmte Fence-Zeile auf 30 Runen |
| N-1 — `planning` trimmt unicode-weit, Anlassfall bleibt still | HIGH | **geheilt** | Alt-Image gegen HEAD-Image auf demselben Fixture (Floskel hinter einer mit U+00A0 eingerückten Fence-Zeile): alt 0 Befunde, HEAD `closure-note-boilerplate`. Beide `planning`-Stellen haben eine eigene Assertion (Rückbauten rot). Die Umstellung ändert zugleich das Verhalten eines **ausgelieferten** Gates — R-5 |
| N-2 — Grenzen-Aufzählung nennt `pins`/`codepaths` nicht | MEDIUM | **geheilt** | Anforderung und ADR nennen beide, samt der stillen Richtung bei `pins` |
| N-3 — CR einer CRLF-Zeile im Befund-Ziel | MEDIUM | **geheilt** | `--json` auf einer CRLF-Datei liefert das Ziel ohne Steuerzeichen; der Rückbau der CR-Trimmung ist rot |
| N-4 — Vorrang der strengen Lesart ohne Assertion | MEDIUM | **geheilt** | Der Rückbau „Vorrang getauscht" ist jetzt rot (`strenge Lesart hat Vorrang`) |
| N-5 — Ehrlichkeits-Klausel am seltenen Fall aufgehängt | MEDIUM | **geheilt in den drei Verträgen** | Anforderung, Spezifikation und ADR nennen jetzt **beide** Lesarten als potenziell danebenzeigend. Die widerrufene Begründung überlebt im Slice und in der Alternativen-Tabelle — R-6 |
| N-6 — eingerückter Code-Block in keiner Lesart modelliert | INFO | **benannt** | Als eigener Absatz „Ebenfalls außerhalb" in der Anforderung |

### Mutations-Gegenprobe: eigener Nachvollzug

Verfahren: Dateikopie im Temp-Verzeichnis sichern, mutieren, `make test`, aus der
Kopie zurückschreiben (**nicht** `git checkout`); jede Ersetzung bricht ab, wenn
sie nicht **genau einmal** greift. Der Arbeitsbaum war nach jedem Durchlauf leer.

| Rückbau | Ergebnis | roter Testfall |
|---|---|---|
| `TrimFenceIndent` wieder unicode-weit | rot | vier Fälle, davon zwei aus dem `planning`-Modul |
| strenge Lesart ersatzlos entfernt | rot | Tilde/Backtick-Kreuzfälle, kürzerer Schluss, Vorrang |
| naive Lesart ersatzlos entfernt | rot | „nur die naive Lesart kippt" |
| Kappung von 30 auf 300 | rot | „auf 30 Runen gekappt" |
| CR nicht mehr aus dem Ziel getrimmt | rot | Ziel-ohne-CR-Test |
| Vorrang der Lesarten getauscht | rot | „strenge Lesart hat Vorrang" |
| Längen-Abgleich in `FenceCloses` entfernt | rot | drei Fälle |
| `closureHeadingLine` wieder unicode-weit | rot | U+00A0 vor der Überschrift |
| `closureSectionProse` wieder unicode-weit | rot | U+00A0 in der Notiz |
| Aufruf aus der Modul-Oberfläche entfernt | rot | Verdrahtungs-Test |
| **Tabellen-Leser trimmt wieder unicode-weit** | **grün** | keiner — R-1 |
| **`proseLines` trimmt wieder unicode-weit** | **grün** | keiner — R-1 |
| **`diagramFenceLines` trimmt wieder unicode-weit** | **grün** | keiner — R-1 |

Die im Slice behaupteten „sechs Rückbauten, alle rot" sind damit bestätigt und um
vier weitere rote ergänzt; **drei** Rückbauten bleiben grün.

---

## Findings

### R-1 — Die Trimm-Invariante trägt die ganze Heilung und ist an drei von fünf Konsumenten nicht assertiert; ein Einzeiler stellt den stillen Grün-Pfad wieder her

- `kategorie`: HIGH
- `quelle`: [§`DC-FA-SPAN-001.a`](../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
  Schritt 3 („mit derselben Trimmung wie die Vorverarbeitung"),
  [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md) §Fitness Function
  („der Test ist mutations-echt"), Reviewer-Skill §HIGH (Stilles-Grün in einem Gate)
- `pfad`: `internal/hexagon/core/rules/markdown.go:45` (die Zusage im Kommentar)
  gegen `internal/hexagon/core/app/trace_table.go:328`,
  `internal/hexagon/core/rules/markdown.go:89` und
  `internal/hexagon/core/rules/markdown.go:128`
- `befund`: Der Kommentar über `TrimFenceIndent` erklärt die Funktion zum
  „einzigen Trimmer vor `FenceToggle`, `FenceRun` und `FenceCloses`" — genau die
  Invariante, auf der die Heilung von F-2 und N-1 ruht. Assertiert ist sie nur an
  den beiden `planning`-Stellen und im Wächter selbst. Wird sie in einem der drei
  übrigen Konsumenten auf `strings.TrimSpace` zurückgedreht, bleibt `make test`
  grün, obwohl der Silent-Grün-Pfad wieder offen ist: mit dem Rückbau im
  Tabellen-Leser meldet `--trace --require-complete` auf einem Dokument mit einer
  mit U+00A0 eingerückten Fence-Zeile „1 Anforderung(en), 0 Waise(n)" und Exit 0,
  wo HEAD „2 Anforderung(en), 1 Waise(n)" und Exit 1 liefert — und der Wächter
  meldet in **beiden** Fällen 0 Befunde.
- `verifizierbar`: ja — die drei Rückbauten einzeln über eine Dateikopie
  angewendet, `make test` je Exit 0; aus dem Rückbau im Tabellen-Leser zusätzlich
  ein Image gebaut und gegen ein Fixture mit tabellarischer Anforderungsquelle
  gefahren (Zahlen oben). Arbeitsbaum danach wiederhergestellt.
- `klasse`: „geteiltes Prädikat ohne Assertion gegen Wieder-Divergenz"

### R-2 — Das Modul `citations` gruppiert Absätze selbst; ein Fence ist dort keine Absatzgrenze, und die fail-closed-Zusage kippt auf Exit 0

- `kategorie`: HIGH
- `quelle`: [§`DC-FA-CITE-001.a`](../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
  Schritt 2 („der nächste inline-Zitat-Span im **selben Absatz**" … „fail-closed
  — die Direktive ist unbrauchbar, kein Schweigen") und der Kopfsatz derselben
  Sektion („fence-aware wie die übrigen Module");
  [`DC-FA-CITE-001`](../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in);
  Absatz-Semantik aus [§`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  Schritt 2
- `pfad`: `internal/hexagon/core/rules/citations.go:176` gegen
  `internal/hexagon/core/rules/markdown.go:167` (`proseParagraphs`)
- `befund`: Die Absatzgrenze ist im Produkt eine geteilte Lexik — `proseParagraphs`
  beendet einen Absatz an einer Leerzeile **und** an einer Lücke in der
  Zeilennummerierung, also an einem übersprungenen Fence. `citations` baut seinen
  Absatz selbst und bricht ausschließlich an Leerzeilen; ein Fenced-Code-Block
  zwischen Direktive und Zitat trennt dort nicht. Dieselbe Datei liefert deshalb
  je nach Trennzeichen zwei entgegengesetzte Ergebnisse, und die Richtung, in der
  sie grün wird, ist die, in der das Modul eigentlich fail-closed abbrechen müsste.
- `verifizierbar`: ja — Fixture mit einer `d-check:cite`-Direktive, einer
  Prosa-Zeile ohne Zitat, einem Fenced-Code-Block und dahinter einem Zitat, das
  zur Quell-Spanne passt: **0 Befunde, Exit 0**. Dieselbe Datei mit einer
  Leerzeile statt des Fence-Blocks: `d-check:cite ohne folgendes Zitat …
  fail-closed`, **Exit 2**. Passt das Zitat hinter dem Fence **nicht**, liefert
  die Fence-Variante ein falsches `citation-mismatch` (Exit 1) statt des
  fail-closed. Gegenprobe zur Absatz-Semantik des restlichen Produkts: dieselbe
  Konstellation im Modul `spans` liefert mit Fence zwei `span-unclosed`-Befunde
  (zwei Absätze), ohne Fence null (ein Absatz).
- `klasse`: „geteilte Lexik, vom Konsumenten selbst vorbereitet"

### R-3 — `headingSection` beantwortet die Anker-Frage roh; ein HTML-Anker in einem Fence löst `versions.current-from` auf, während `anchors` denselben Anker im selben Lauf für nicht existent hält

- `kategorie`: MEDIUM
- `quelle`: [§`DC-FA-ANCH-001.b`](../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker)
  („zeilenbasiert, **außerhalb Fences**"),
  [`DC-FA-VER-001`](../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
  (fail-closed bei unauflösbarer Quelle),
  [`DC-FA-PIN-001`](../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in)
- `pfad`: `internal/hexagon/core/rules/versions.go:121` (der HTML-Anker-Zweig
  läuft über die rohen Zeilen) gegen `internal/hexagon/core/rules/anchors.go:120`
  (`htmlAnchors` liest die vorverarbeiteten, fence-freien Zeilen)
- `befund`: Ob ein Inline-HTML-Anker existiert, beantworten zwei Stellen
  verschieden. `htmlAnchors` fragt die Vorverarbeitung und übersieht Anker in
  Fenced-Code (so zugesagt — GitHub rendert sie nicht als Sprungziel);
  `headingSection`, der Span-Lieferant von `versions` und `pins`, scannt die
  rohen Zeilen und findet sie. Ein Anker, der ausschließlich in einem
  Beispielblock steht, erfüllt damit die fail-closed-Bedingung von
  `versions.current-from` und liefert die „aktuelle Version" aus dem Beispiel.
- `verifizierbar`: ja — Fixture mit `versions.current-from` auf einen Anker, der
  nur in einem Fenced-Code-Block steht: **Exit 0**, die Versionsprüfung läuft
  gegen den Wert aus dem Beispielblock. Gegenprobe mit einem gar nicht
  existierenden Anker: `versions.current-from: Anker … nicht auflösbar`,
  **Exit 2**. Im selben Lauf meldet `anchors` auf einen Link, der genau diesen
  Anker adressiert, `anchor-missing` — ein Lauf, zwei Antworten auf dieselbe Frage.
- `klasse`: „geteilte Lexik, vom Konsumenten selbst vorbereitet"

### R-4 — Grund-Code-Tabelle und `--doctor`-Klartext tragen die mit 0.52.1 widerrufene Reichweiten-Zusage; ein Lauf widerlegt sie in derselben Ausgabe

- `kategorie`: MEDIUM
- `quelle`: [§4 der Spezifikation](../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
  (Grund-Codes, „stabil, maschinenlesbar") gegen
  [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
  Klasse 3 in der Fassung 0.52.1/0.52.2
- `pfad`: `spec/spezifikation.md:2354` und
  `internal/hexagon/core/app/diagnose.go:109` (Klartext des Grund-Codes)
- `befund`: Die Zeile beschreibt den Befund weiter als „alles dahinter gilt
  **jeder** Vorverarbeitung als Code und wird von **allen** Modulen übersprungen.
  Befund an der Öffnungszeile" — beides ist genau die Formulierung, die der
  Erst-Review als überzogene Reichweite beanstandet hat und die im Lastenheft
  seither korrigiert ist (Fundstelle statt Reparaturstelle; „mindestens ein
  Modul" statt „alle"). Seit 0.52.1 kann ein `fence-unclosed` entstehen, ohne dass
  irgendein naiver Leser etwas überspringt, und die Aussage ist dann falsch. Der
  `--doctor`-Klartext trägt dieselbe Zusage und ist die Fassung, die der Nutzer
  liest.
- `verifizierbar`: ja — Datei mit einem Backtick-Block, den eine Tilden-Zeile
  „schließt", und einem kaputten Link dahinter, Module `spans` und `links`: ein
  Lauf, zwei Befunde — `fence-unclosed` mit dem Klartext „alles dahinter wird von
  jedem Modul übersprungen" **und** ein `target-missing` auf genau der Zeile
  dahinter, die laut Klartext übersprungen worden sein müsste.
- `klasse`: „Rand auf widerrufener Fassung stehengeblieben"

### R-5 — Die Verhaltensänderung am ausgelieferten `planning`-Gate steht in keiner konsumentensichtbaren Zeile, und die SemVer-Begründung deckt ihre Richtung nicht

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (unverändert), Lastenheft-Historie 0.52.2 (nennt nur `DC-FA-SPAN-001`),
  [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md) Entscheidung 5
  („d-check findet danach mehr")
- `pfad`: `internal/hexagon/core/rules/planning.go:168` und
  `internal/hexagon/core/rules/planning.go:194`; Vertragstext in
  `spec/lastenheft.md`
- `befund`: Der Commit ändert die Fence-Erkennung der Closure-Fähigkeit, die mit
  Release 0.52.0 ausgeliefert ist. Die Richtung ist **nicht** nur „findet mehr":
  zwei der drei messbaren Wirkungen sind „findet weniger". Die Änderung ist
  fachlich richtig und von Schritt C4 der Spezifikation gedeckt („Fence-Lexik wie
  im übrigen Scanner"), aber weder das Lastenheft-Delta 0.52.2 noch die
  Anforderung selbst noch die SemVer-Begründung der ADR erwähnen, dass ein
  Konsument nach dem Update an einem Closure-Gate **Befunde verliert**.
- `verifizierbar`: ja — dasselbe Fixture gegen das Image vor der Heilung und
  gegen HEAD, Profil `--enable planning` mit gesetztem Closure-Verzeichnis: (a)
  Sätze hinter einer mit U+00A0 eingerückten Fence-Zeile — alt
  `closure-note-thin`, HEAD still; (b) Closure-Überschrift hinter derselben
  Zeile — alt `closure-note-missing`, HEAD still; (c) Floskel dahinter — alt
  still, HEAD `closure-note-boilerplate`.
- `klasse`: „Verhaltensänderung eines ausgelieferten Gates ohne Vertrags-Delta"

### R-6 — Entscheidungs-Protokoll, Alternativen-Tabelle und ADR-Index begründen weiter mit einer Aussage, die Entscheidung 3 widerrufen hat

- `kategorie`: LOW
- `quelle`: [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md) Entscheidung 3
  („eine **Fundstelle**, keine Reparaturstelle") gegen dieselbe ADR
  §Verglichene Alternativen und
  [slice-101](../plan/planning/in-progress/slice-101-fence-unbalanciert.md) §4
  Punkt 3
- `pfad`: `docs/plan/adr/0050-fence-unclosed-in-spans.md` (Alternativen-Zeile
  „Befund an der letzten Zeile der Datei"), `docs/plan/adr/README.md` (Index-Zeile
  zu ADR-0050), `docs/plan/planning/in-progress/slice-101-fence-unbalanciert.md`
  §4 Punkt 3
- `befund`: Drei Ränder stehen auf der widerrufenen Fassung. Der Slice hält als
  **Entscheidung** fest: „Die Öffnungszeile ist der Ort der Reparatur, also der
  Ort des Befundes." Die Alternativen-Tabelle verwirft „Befund an der letzten
  Zeile der Datei" mit der Begründung, ein Befund am Dateiende „zeigt auf eine
  Zeile, die niemand ändern muss" — genau die Eigenschaft, die Entscheidung 3
  inzwischen für den gewählten Ort einräumt. Die Index-Zeile nennt weiterhin nur
  „Befund an der **Öffnungszeile**".
- `verifizierbar`: ja — Fixture mit einem Öffner ohne Schluss in Zeile 3 und zwei
  intakten Blöcken dahinter: der Befund steht auf Zeile 14, der Schlusszeile eines
  intakten Blocks; die Reparatur gehört hinter Zeile 5. Der ausgeschlossene
  Kandidat „letzte Zeile der Datei" wäre in diesem Dokument Zeile 14 gewesen —
  dieselbe Zeile.
- `klasse`: „Rand auf widerrufener Fassung stehengeblieben"

### R-7 — Die Grenze zählt Scan-Scope und Ziel-Achse auf; die dritte Klasse unerreichbarer Eingaben sind git-Revisionen

- `kategorie`: LOW
- `quelle`: [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
  Absatz „Grenze"; [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md)
  §Konsequenzen; [§`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  Schritt 4
- `pfad`: `internal/hexagon/core/rules/vcs.go:129` (`vcsCore` ruft
  `excludedRanges`) und `internal/hexagon/core/rules/matrix.go:302`
  (`excludedRanges` läuft über `extractHeadingLines`, also über den
  fence-empfindlichen Prosa-Automaten)
- `befund`: Die Section-Maske des Immutabilitäts-Gates wird auf **git-Blobs**
  gerechnet, die der Wächter niemals sieht: er scannt den Arbeitsbaum. Ein
  unbalancierter Fence in einer Revision verschiebt dort die Abgrenzung von
  `vcs.exclude-sections` und lässt das Gate einen anderen Core vergleichen als
  gemeint. Die Grenze nennt zwei Klassen (Post-Pässe über selbst benannte
  Verzeichnisse, Zieldateien außerhalb der Scan-Wurzeln) — Revisionen fehlen, und
  für jeden bereits existierenden Commit ist der Wächter rückwirkend blind.
- `verifizierbar`: ja — git-Fixture mit einer Accepted-ADR und
  `exclude-sections` auf den Geschichts-Abschnitt: ein reiner Geschichts-Eintrag
  über zwei Commits ist befundfrei (Exit 0); dieselben zwei Commits mit einem
  unbalancierten Fence vor dem Abschnitt melden `core-drift-vcs` (Exit 1),
  obwohl außerhalb des ausgenommenen Abschnitts nichts geändert wurde.
- `klasse`: „Modul-Grenze nur auf der Quell-Achse gedacht"

### R-8 — Die drei Historien-Zeilen des Heilungs-Commits datieren einen Tag vor die Re-Review, die sie ausgelöst hat

- `kategorie`: LOW
- `quelle`: `harness/conventions.md` (Historien-Pflege der Spec-Straten);
  [Re-Review-Report](2026-08-09-slice-101-fence-re-review.md) (Kopf: Datum
  2026-08-10)
- `pfad`: `spec/lastenheft.md:2508`, `spec/spezifikation.md:2393` und
  `docs/plan/adr/0050-fence-unclosed-in-spans.md` §Geschichte (dritter Eintrag)
- `befund`: Alle drei Zeilen tragen das Datum 2026-08-09 und beschreiben den
  Nachzug „nach bestätigender Re-Review". Der Commit `fc94503` und der Commit
  des Re-Review-Reports sind auf 2026-08-10 datiert, und der Report selbst nennt
  im Kopf ausdrücklich den 2026-08-10 als Lauf-Datum. Die Historie ist das
  Audit-Artefakt der Spec-Straten; sie behauptet hier eine Änderung am Tag vor
  ihrem Anlass.
- `verifizierbar`: ja — `git log --date=short` für `fbb2d73` und `fc94503`
  gegen die drei Zeilen.
- `klasse`: „Historien-Zeile trägt nicht das Datum ihrer Änderung"

## Negativbefunde

- geprüft, ohne Befund: **Vollzähligkeit der Trimmung (Kernfrage A).** Alle
  Aufrufstellen der Fence-Lexik im Repo aufgezählt und einzeln nachverfolgt:
  `FenceToggle` an sieben Stellen (Wächter zweimal, `planning` zweimal,
  `proseLines`, `diagramFenceLines`, Tabellen-Leser), `FenceRun` an vier,
  `FenceCloses` an zwei — **jede** wird aus `TrimFenceIndent` gespeist bzw.
  bekommt eine bereits so getrimmte Zeile. Es gibt keine Aufrufstelle mehr, die
  selbst trimmt, auch nicht in Tests oder außerhalb des Regel-Pakets.
- geprüft, ohne Befund: **kein dritter Fence-Automat.** Im Produkt existieren
  genau zwei Schluss-Lesarten (naiver Toggle, längenabgeglichener Schluss), beide
  geteilt. Das Modul `versions` liest Pins bewusst **auch** in Fences (gescopte,
  spezifizierte Ausnahme, kein eigener Automat); `targets` liest selbst benannte
  Doku-Tabellen roh und fällt damit in die bereits dokumentierte Post-Pass-Klasse.
- geprüft, ohne Befund: **Heading-Lexik.** `ExtractHeadings`, `HeadingSlugs`,
  `extractHeadingLines`, `excludedRanges`, `headingSection` und der Slug-Cache
  von `anchors`/`codepaths` hängen alle am selben `proseLines`-Automaten; eine
  zweite Heading-Erkennung existiert nicht. Auch `closureHeadingLine` nutzt den
  geteilten ATX-Parser.
- geprüft, ohne Befund: **Inline-Code-Lexik.** Der Wächter (`unclosedRuns`) und
  die Vorverarbeitung (`forEachInlineCodeSpan`) paaren über dieselbe Hilfsfunktion
  `findClosingRun` mit derselben Regel (die öffnende Folge bestimmt die
  schließende, ungeschlossene Folgen sind literal). Keine Divergenz.
- geprüft, ohne Befund: **`immutable` (Arbeitsbaum-Variante).** Es rechnet den
  Core über dieselbe fence-empfindliche Section-Maske, aber auf der **gescannten**
  Datei — dort greift der Wächter. Nur die git-Revision liegt daneben (R-7).
- geprüft, ohne Befund: **Zitat-Quellspanne.** `citations` liest die zitierte
  Quelle rein zeilennummern-basiert; dort wird keine Fence-Frage gestellt, also
  auch keine verschieden beantwortet. Die Divergenz von R-2 liegt ausschließlich
  in der Absatz-Gruppierung der **zitierenden** Datei.
- geprüft, ohne Befund: **Der Eingriff in den Tabellen-Leser ist ein No-op.**
  `TrimFenceIndent` hat exakt den Rumpf, der dort ersetzt wurde; der RTM-Pfad kann
  sich durch diesen Commit nicht geändert haben. `make completeness-check` über
  das eigene Repo: 48 Anforderungen, 0 Waisen.
- geprüft, ohne Befund: **Kein Falsch-Positiv im Ökosystem.** Reines
  `spans`-Profil: eigenes Repo 350 Dateien, die beiden Schwester-Repos 224 bzw.
  222 Dateien — **796 Dokumente, 0 Befunde**. Die Zusage „der Bestand bleibt bei
  null" hält auch nach dem zweiten Heilungs-Commit.
- geprüft, ohne Befund: **Gate-Stand.** `make gates` Exit 0 (350 Dateien / 0
  Befunde, Coverage 94,30 % gegen Schwelle 93 %, alle acht Gates),
  `make verify-closure-notes` 321 Dateien / 0 Befunde.
- geprüft, ohne Befund: **Befund-Ziel auf CRLF.** Das JSON-Feld trägt die reine
  Fence-Zeile ohne Steuerzeichen; der `--doctor`-Klartext zeigt sie an der
  Fundstelle.
- geprüft, ohne Befund: **Ziel-Achsen-Aufzählung.** Die sechs in der Anforderung
  genannten Module lesen Zieldateien tatsächlich über die fence-empfindliche
  Vorverarbeitung (je im Code belegt); es fehlt kein siebtes Modul dieser Bauform.
- geprüft, ohne Befund: **Referenz-Richtung (SDP).** Der Heilungs-Commit trägt
  keinen Provenance-Marker und keinen neuen Abwärtsverweis aus ADR oder
  Spec-Straten auf Planning-Artefakte; das Referenzmatrix-Modul ist im grünen
  `make gates`-Lauf aktiv.
- geprüft, ohne Befund: **Import-Regeln und Ablage.** Das neue geteilte Prädikat
  liegt im Regel-Paket und wird vom App-Paket konsumiert — Richtung unverändert;
  `make lint`, `make arch-check` und `make semgrep` im Gate-Lauf grün. ADR im
  Index, Lastenheft-Version auf 0.52.2 gehoben, Historien-Zeile in beiden
  Spec-Straten (Datum siehe R-8), Slice-Datei in `in-progress`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** geteiltes Prädikat ohne Assertion gegen
Wieder-Divergenz · geteilte Lexik, vom Konsumenten selbst vorbereitet · Rand auf
widerrufener Fassung stehengeblieben · Verhaltensänderung eines ausgelieferten
Gates ohne Vertrags-Delta · Modul-Grenze nur auf der Quell-Achse gedacht ·
Historien-Zeile trägt nicht das Datum ihrer Änderung

**Wiederholungs-Signal — die Schwelle ist überschritten.** „Modul-Grenze nur auf
der Quell-Achse gedacht" ist mit R-7 zum **dritten** Mal dran (F-3, N-2, R-7);
nach dem Reviewer-Skill ist das der Punkt, an dem nicht mehr gemeldet, sondern
der Sensor nachgezogen wird. Die im Slice §8 selbst formulierte Oberklasse —
„eine geteilte Lexik driftet an den Rändern, weil jeder Konsument sie selbst
vorbereitet" — steht mit R-2 und R-3 bei vier Vorkommen über drei Reviews. Beides
gehört bei der Closure ins Beobachtungs-Register.

## Verdikt

**Merge-blockierend:** ja — zwei HIGH und drei MEDIUM.

**Ist die Klasse geschlossen? Nein — aber die Fence-Lexik ist es.** Für Fences
gibt es kein Leck mehr: alle sieben Aufrufstellen der Fence-Lexik werden aus
einem Trimmer gespeist, es gibt genau zwei Schluss-Lesarten und beide sind
geteilt, und weder Heading-Erkennung noch Slug-Cache noch Section-Masken noch
Diagramm-Fences noch der Zitat-Quellpfad noch die Inline-Code-Paarung tragen eine
zweite Antwort. Eine dritte Instanz **im Fence-Bereich** existiert nicht; das ist
in diesem Report der einzige belastbare Negativbefund, der die Leitfrage direkt
beantwortet.

Die Klasse ist trotzdem offen, und zwar auf beiden Achsen, nach denen zu suchen
war. **Zeitlich:** die Zusammenführung ist nicht gegen ihr Wiederauseinanderdriften
gesichert — drei Einzeiler stellen den stillen Grün-Pfad wieder her, ohne dass ein
Test rot wird, und für einen davon ist der Weg von der Zeile bis zum Exit 0 über
einer ungedeckten Anforderung hier durchgemessen (R-1). **Sachlich:** dieselbe
Bauform steckt in zwei weiteren geteilten Lexiken. Die Absatzgrenze beantwortet
`citations` selbst, und die Richtung, in der es dadurch grün wird, ist genau die,
in der das Modul fail-closed abbrechen müsste — Exit 2 wird zu Exit 0, allein weil
zwischen Direktive und Zitat ein Code-Block statt einer Leerzeile steht (R-2). Die
Anker-Frage beantwortet `headingSection` roh, sodass ein Anker aus einem
Beispielblock die fail-closed-Bedingung von `versions.current-from` erfüllt,
während `anchors` im selben Lauf sagt, dass es ihn nicht gibt (R-3). Beide sind
älter als dieser Slice, aber beide sind exakt der Befund, den der Slice zu
schließen behauptet.

Der Kern der Arbeit ist gut. Sechs zugesagte Rückbauten sind unabhängig
nachvollzogen und rot, vier weitere kommen dazu; N-1 bis N-6 sind sachlich
abgearbeitet; 796 reale Dokumente in drei Repos erzeugen kein einziges
Falsch-Positiv; das größte Regressionsrisiko — der Tabellen-Leser — ist durch
einen beweisbaren No-op geschützt und durch die Vollständigkeits-Zahlen bestätigt.

**Release-Empfehlung: noch nicht releasen.** Die Einordnung als **Minor** ist
korrekt — kein Schema-Delta, keine neue Option, keine geänderte Exit-Code-Semantik,
und die neue Befundfläche macht grüne Konsumentenläufe rot, nicht umgekehrt. Zwei
Dinge würden einen Konsumenten beim Update aber unangenehm überraschen und stehen
in keiner Release-Notiz: er verliert an einem ausgelieferten Closure-Gate Befunde
(R-5, Richtung „findet weniger", von der SemVer-Begründung nicht gedeckt), und
der Klartext, den er zum neuen Befund liest, sagt ihm eine Reichweite zu, die
derselbe Lauf widerlegt (R-4). Vor dem Tag gehören R-1, R-4, R-5 und die drei
Ränder aus R-6/R-8 in diesen Slice — sie sind klein und liegen alle auf dem
Gesamtdiff. R-2, R-3 und R-7 gehören **nicht** hierher: sie sind ältere Defekte in
anderen Modulen mit eigenem Vertrag. Sie gehören in einen Folge-Slice und ins
Beobachtungs-Register, und die Grenze in
[`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
sollte vorher noch die dritte unerreichbare Eingabe-Klasse nennen.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler — mit der
ausdrücklichen Notiz, dass eine Klasse die Drei-Schwelle erreicht hat. Dieser
Report ist ein Lauf-Beleg (dieser Diff, dieser Skill, dieses Modell, dieses
Verdikt) und ersetzt keine Verifikation — DoD- und Spec-Konformität prüft der
Verifier separat.
