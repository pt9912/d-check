# Review-Report: slice-099 (Modul `structure`, das 20. Regelmodul) — 2026-08-15

**Review-Art:** Code — geprüft wird der Diff gegen Lastenheft, Spezifikation,
[ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md), den
Slice-Plan und die Konventionsregel
[`MR-025`](../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten).
Die DoD-Abhakung ist ausdrücklich **nicht** Gegenstand (das ist Verifikation,
anderer Kontext, anderes Prüf-Artefakt).

**Gegenstand:** `64c62cb..HEAD` (`59a73a2`) — vier Commits: Refactor der
Abschnitts-Mechanik (`0960185`), Modul-Kern samt neun Grund-Codes (`3ea969e`),
CLI-Spiegel und Selbst-Aktivierung (`e93d6a9`), Release-Prep v0.57.0
(`59a73a2`).

**Skill:** `.harness/skills/reviewer.md` @ 1.3.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-08-15

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [slice-099](../plan/planning/in-progress/slice-099-structure-modul.md),
  besonders §3a (Spiegel-Liste) und §4a (Messung)
- [welle-73](../plan/planning/welle-73-structure-umsetzung.md) und
  [`MR-025`](../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
- [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  und [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  (acht Schritte), §2-Schema, §4-Grund-Code-Tabelle
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  und [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritte C3–C5 (Preset-Kopplung)
- [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben),
  [`DC-FA-CLI-006`](../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten),
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
- [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) (Accepted),
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
- [`AGENTS.md`](../../AGENTS.md) Hard Rules

**Eigene Läufe dieses Reviews:** `make build`, `make lint`, `make test`
(Baseline plus zwölf Mutationsläufe), `make gates`, `make adr-check`,
`make verify-closure-notes`; dazu ein aus dem Vor-Commit `64c62cb` gebautes
Vergleichs-Image (`make build IMAGE=d-check-old` in einem `git archive`-Auszug)
und Image-Läufe gegen acht Fixture-Repos in einem Temp-Verzeichnis. Der
Arbeitsbaum ist unverändert (Beleg am Ende).

---

## Findings

### F-1 — `closure-note-ambiguous` existiert als Grund-Code, aber kein Codepfad erzeugt ihn

- `kategorie`: HIGH
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  Akzeptanzkriterium „Negative (mehrdeutig)",
  [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C3, [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  §Fitness Function („Mehrdeutigkeit schlägt Messung")
- `pfad`: `internal/hexagon/core/rules/planning.go:171-179` und
  `internal/hexagon/core/rules/planning_closure_test.go:757-773`
- `befund`: `closureHeadingLine` sammelt über `FindSectionHeads` alle Treffer
  und gibt `heads[0]` zurück; die übrigen werden verworfen, ohne dass irgendwo
  ein Befund entsteht. Der Lauf gegen ein `done/`-Fixture mit zwei auf
  `heading-pattern` passenden Überschriften (Zeile 3 leer, Zeile 7 mit Substanz)
  meldet `closure-note-thin` in Zeile 3 statt `closure-note-ambiguous` in
  Zeile 7 — Exit 1, aber mit dem falschen Code und ohne den zugesagten Abbruch.
  Die Konstante steht in `internal/hexagon/core/model/finding.go:46`, der
  Klartext in `internal/hexagon/core/app/diagnose.go:128`, die Zeile in
  Spezifikation §4 und im Handbuch — der einzige fehlende Teil ist die
  Erzeugung. Der im selben Diff hinzugefügte Test
  `TestClosureErsterTrefferBeiMehrerenUeberschriften` schreibt das alte
  Verhalten ausdrücklich fest („der ERSTE Abschnitt wird gemessen ⇒ thin
  erwartet") und trägt den Kommentar „bis der Grund-Code existiert" — der Code
  existiert seit demselben Diff. Die Lockstep-Verriegelung
  (`internal/hexagon/core/app/diagnose_test.go:41`) greift nicht: sie prüft
  Katalog-Deckung zwischen `AllReasons()`, `reasonTexts()` und §4, nicht
  Erreichbarkeit.
- `verifizierbar`: ja — Fixture mit zwei Closure-Überschriften,
  `--enable planning --json`: Ist `closure-note-thin` Zeile 3, Soll
  `closure-note-ambiguous` Zeile 7 und **kein** `-thin`.
- `klasse`: „Grund-Code deklariert und dokumentiert, aber von keinem Codepfad
  erzeugt"

### F-2 — Release-Notiz, Lastenheft-Historie und Handbuch sagen eine Verhaltensänderung zu, die nicht eintritt

- `kategorie`: MEDIUM
- `quelle`: [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  §Konsequenzen („Ein Repo mit zwei Closure-Abschnitten wird danach rot … Das
  gehört in die Release-Notiz"), [slice-099](../plan/planning/in-progress/slice-099-structure-modul.md)
  §4 DoD
- `pfad`: `CHANGELOG.md:25-28`, `spec/lastenheft.md:2609`,
  `docs/user/benutzerhandbuch.md:1841`
- `befund`: Der Changelog-Eintrag sagt „ohne eindeutigen Abschnitt wird jetzt
  **nicht** gemessen, statt still den ersten zu nehmen", die Lastenheft-Zeile
  0.57.0 nennt die „seit 0.51.0 zugesagte, aber nicht implementierte
  Mehrdeutigkeits-Härte" als geliefert, die Handbuch-§11-Zeile „Dazu
  `closure-note-ambiguous` (mehrere Closure-Überschriften)". Ein Konsument, der
  nach diesem Text seinen Bestand prüft, erwartet einen neuen roten Befund und
  bekommt denselben Befundsatz wie unter v0.56.0 — belegt durch den Vergleich
  gegen das aus `64c62cb` gebaute Image: über ein Fixture mit vierzehn
  Closure-Randfällen ist die `--json`-Ausgabe beider Images byte-identisch,
  einschließlich des Falls mit zwei Überschriften. Die Zusage ist damit in drei
  ausgelieferten Flächen falsifizierbar.
- `verifizierbar`: ja — derselbe Fixture-Lauf gegen `v0.56.0`-Image und
  Arbeitsbaum-Image liefert identische Befunde.
- `klasse`: „Release-Notiz sagt Verhaltensänderung zu, die der Lauf nicht zeigt"

### F-3 — Marken-Grenze prüft ASCII-Bytes; ein Umlaut hinter der Marke gilt als Grenze

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  §„Marken sind Auszeichnungs-Marken" und
  [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 6 („endet oder mit einem nicht-alphanumerischen Zeichen weitergeht")
- `pfad`: `internal/hexagon/core/rules/structure.go:197` und
  `internal/hexagon/core/rules/structure.go:204-206`
- `befund`: `hasMarker` entscheidet die Fortsetzungs-Grenze über
  `isAlnumByte(rest[0])`, also über ein einzelnes Byte im ASCII-Bereich. Bei
  einem mehrbyte-kodierten Buchstaben ist das erste Fortsetzungs-Byte kein
  ASCII-Alphanumerikum, und die Marke gilt als vorhanden. Gemessen im Fixture:
  `require-all: ["Ergebnis"]` gilt gegen `- **Ergebnisüberblick:**` als erfüllt
  (kein Befund), gegen `- **Ergebnisliste:**` als verletzt
  (`section-marker-missing`). Die zugesagte Unterscheidung trennt damit nach
  ASCII-Zugehörigkeit statt nach Alphanumerik, und zwar in der Richtung, die
  falsch grün wird. In einem deutschsprachigen Bestand ist die Klasse
  einschlägig (Überblick, Änderung, Begründung, Lösung als Marken-Fortsetzung).
  Für die verwandte Floskel-Bedingung desselben Moduls hat die
  Lastenheft-Historie 0.56.1 dieselbe ASCII-Grenze ausdrücklich in beiden
  Richtungen in den Vertrag geschrieben; für `require-all` steht sie nirgends.
- `verifizierbar`: ja — Fixture mit `- **Ergebnisüberblick:**` und
  `require-all: ["Ergebnis"]`, `--enable structure --json`: Ist kein Befund,
  Soll `section-marker-missing`.
- `klasse`: „ASCII-Bytegrenze als Alphanumerik-Grenze (Nicht-ASCII gilt als
  Trennzeichen)"

### F-4 — Schritt 5 und 6 beschreiben die geteilte Bereinigung und Zählung anders als der Preset-Partner und als der Code

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritte 5–6 und §2-Schema-Zeile zu `structure[].min-sentences`;
  [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  §Konsequenzen („eine Änderung an einer der beiden Stellen ohne die andere ist
  ein Spezifikations-Bug")
- `pfad`: `spec/spezifikation.md:1865-1868`, `spec/spezifikation.md:1876`,
  `spec/spezifikation.md:2330` gegen
  `internal/hexagon/core/rules/sections.go:63-74` und
  `internal/hexagon/core/rules/planning.go:308-322`
- `befund`: Schritt 5 sagt, es würden „die Fenced-Code-Blöcke entfernt";
  `SectionProse` ruft `PreprocessMarkdown` und leert zusätzlich die
  Inline-Code-Spans. Beobachtbar: `require-pattern: 'Beleg'` meldet
  `section-pattern-missing`, obwohl der Abschnitt „Der `Beleg` steht hier"
  trägt, und `forbid-pattern: 'TODO'` meldet **nicht**, obwohl „`TODO`" im
  Abschnitt steht — ein stiller Falsch-Negativ am Verbots-Muster. Schritt 6
  sagt „die Zahl der Satzende-Zeichen (`.`, `!`, `?`)"; `countSentenceEnds`
  zählt nur, wenn Whitespace, Zeilenende oder Wagenrücklauf folgt: ein
  Abschnitt mit sechs Punkten wird mit 2 gemessen und meldet bei
  `min-sentences: 4` `section-thin`. Beide Einschränkungen stehen für den
  Preset-Partner seit 0.55.0 ausdrücklich in
  [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C4 — für `structure` in keiner der drei Flächen (Algorithmus,
  §2-Schema, Handbuch-§5). Die Mechanik ist eine, die Beschreibung nicht.
- `verifizierbar`: ja — Fixture mit Inline-Code-Marke und Punkt-Kette,
  `--enable structure --json`.
- `klasse`: „Preset-Partner-Spezifikationen driften an derselben geteilten
  Mechanik"

### F-5 — `DC-FA-CLI-010`: die Akzeptanzkriterien führen `doc-structure` nicht

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  Akzeptanzkriterien Happy Path und Boundary;
  [`MR-025`](../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
  Spiegel-Tabelle („Anforderung: Beschreibung **und** Akzeptanzkriterien")
- `pfad`: `spec/lastenheft.md:490` und `spec/lastenheft.md:491`
- `befund`: Beschreibung („zwölf") und Out-of-Scope-Aufzählung sind auf
  `doc-structure` nachgezogen, die beiden Akzeptanzkriterien nicht: sie
  enumerieren weiterhin `doc-check`, `doc-trace`, `doc-complete`, `doc-doctor`,
  `doc-repair`, `doc-immutable`, `doc-commits`, `doc-planning`, `doc-tracked`,
  `doc-targets`, `doc-help` — elf Targets. `--print-mk` liefert zwölf. Eine
  AK-getriebene Verifikation prüfte das zwölfte Target nicht, und die
  Anforderung widerspricht sich innerhalb einer Bildschirmseite. Dieselbe
  Stelle wurde in der Lastenheft-Historie 0.37.1 schon einmal als
  „Selbstwiderspruch behoben" saniert (neun auf zehn Targets, `doc-tracked`).
- `verifizierbar`: ja — `--print-mk` zählt zwölf `doc-*`-Targets, das
  Akzeptanzkriterium verlangt elf.
- `klasse`: „Zahl im Beschreibungstext nachgezogen, Enumeration in den
  Akzeptanzkriterien nicht"

### F-6 — `--suggest-config`: das 20. Modul fehlt in der Nicht-aktiviert-Enumeration, im Vertrag wie in der Ausgabe

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-CLI-006`](../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
  §Aufnahme ins Modulset („**Alle übrigen** Module bleiben inaktiv: …"; „Die
  Vorlage nennt die nicht aktivierten Module in einem Kommentar … Auffindbarkeit
  ohne stilles Aktivieren eines inerten Moduls");
  [slice-099](../plan/planning/in-progress/slice-099-structure-modul.md) §3a
  Zeile `--suggest-config`
- `pfad`: `spec/lastenheft.md` §`DC-FA-CLI-006` (Aufnahme-Absatz) und
  `internal/adapter/driving/cli` (erzeugte Vorlage, Kommentarzeile
  „Weitere opt-in-Module sind situativ …")
- `befund`: Die Vertragsklausel klassifiziert jedes Modul entweder als
  aktiv-in-der-Vorlage oder als eines der „alle übrigen" mit Grund; `structure`
  ist in beiden Aufzählungen nicht enthalten. Die erzeugte Ausgabe von
  `--suggest-config ai-harness` nennt „external, diagrams, versions, pins,
  immutable, tracked, targets" — ohne `structure`. Die zugesagte
  Auffindbarkeit gilt für das neueste Modul damit nicht. Die Spiegel-Liste in
  §3a führt genau diese Zeile mit „**ja** — Aufnahme prüfen"; die Prüfung ist
  entweder nicht erfolgt oder ihr Ergebnis nirgends festgehalten. (Dieselbe
  Stelle führt auch `citations` und `sources` nicht — ältere Herkunft, gleiche
  Klasse.)
- `verifizierbar`: ja — `--suggest-config ai-harness` gegen ein
  Minimal-Fixture, Kommentarzeile ohne `structure`.
- `klasse`: „neues Modul nicht in der Nicht-aktiviert-Enumeration von
  `--suggest-config`"

### F-7 — Vier Rückbauten am neuen Modul bleiben grün, darunter ein fail-closed-Rand

- `kategorie`: MEDIUM
- `quelle`: Reviewer-Anker „fehlende Negativtests bei neuem öffentlichen
  Vertrag";
  [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  §fail-closed
- `pfad`: `internal/hexagon/core/rules/structure_test.go` (gesamt), gegen
  `internal/hexagon/core/rules/structure.go:86-87`,
  `internal/hexagon/core/rules/structure.go:135`,
  `internal/hexagon/core/rules/structure.go:120`,
  `internal/hexagon/core/rules/structure.go:205`
- `befund`: Vier gezielte Rückbauten am ausgelieferten Kern lassen `make test`
  grün. (1) Die Schwellen-Bedingung `got < *r.MinSentences` auf `<=` gedreht —
  ein Abschnitt, der die Schwelle **genau** erreicht, meldete fortan
  `section-thin`; kein Test bindet die Schwelle an ihrem Rand. (2) Der
  Ziffernbereich in `isAlnumByte` entfernt — `**Beleg2:**` erfüllte fortan die
  Marke `Beleg`; kein Fixture trägt eine Ziffern-Fortsetzung. (3) Das
  `TrimSpace` im `section-pattern`-Prädikat entfernt — eine Überschrift mit
  nachlaufendem Whitespace träfe ein auf `$` verankertes Muster nicht mehr.
  (4) Der fail-closed-Zweig für eine unlesbare Kandidaten-Datei durch
  `return nil` ersetzt — ein stiller Grün-Pfad, den die Anforderung
  ausdrücklich ausschließt und den kein Test bewacht. Der Commit-Text nennt
  neun Rückbauten mit acht roten; die hier gefundenen liegen außerhalb dieser
  neun.
- `verifizierbar`: ja — jede Mutation einzeln, `make test` grün (Läufe
  durchgeführt, Datei danach aus einer Kopie wiederhergestellt).
- `klasse`: „Schwellen- und fail-closed-Ränder eines neuen Moduls ohne Test"

### F-8 — Die Spiegel-Liste erklärt `.d-check.closure.yml` für unberührt und lässt drei weitere Spiegel aus

- `kategorie`: MEDIUM
- `quelle`: [`MR-025`](../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
  §Adaption; [slice-099](../plan/planning/in-progress/slice-099-structure-modul.md)
  §3a
- `pfad`: `docs/plan/planning/in-progress/slice-099-structure-modul.md:78-80`
  gegen `.d-check.closure.yml` und `Makefile:186-187`
- `befund`: §3a schließt mit „Nicht auf der Liste, weil unberührt:
  `.d-check.yml` …, `.d-check.closure.yml`, die Referenzmatrix." Der Slice fügt
  `.d-check.closure.yml` einen siebenunddreißigzeiligen `structure`-Block hinzu
  und ändert das `verify-closure-notes`-Recipe im `Makefile` — beide Dateien
  fehlen auf der Liste, eine davon steht sogar ausdrücklich als unberührt
  darauf. Ebenfalls nicht gelistet, aber berührt: die Modul-Registry
  `validModules()` in `internal/hexagon/core/model/config.go:13`
  (Vertragsfläche von
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl))
  und die Akzeptanzkriterien von
  [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben),
  die §3a auf „die **Zahl** in der Out-of-Scope-Zeile" verengt, obwohl die
  MR-Tabelle die Akzeptanzkriterien als eigenen Spiegel führt (Folge: F-5).
  Die Zeile zu `AGENTS.md` und `harness/README.md` ist auf „**nur falls** ein
  Gate-Target dazukommt" konditioniert; es kam keines dazu, beide Dateien
  wurden dennoch geändert — die Bedingung beschreibt den eingetretenen Fall
  nicht. Die Liste hat drei Auslassungen produziert, von denen zwei zu
  Findings dieses Laufs geworden sind.
- `verifizierbar`: ja — `git diff 64c62cb..HEAD --stat` gegen die Tabelle in
  §3a.
- `klasse`: „Spiegel-Liste nennt einen Spiegel als unberührt, der berührt wird"

### F-9 — Der Wächter-Kommentar über dem `--print-mk`-Template zählt weiter elf Targets und sechs Format-Verben

- `kategorie`: LOW
- `quelle`: Maintainability (latente Wartungsfalle; der Kommentar ist
  ausdrücklich als Zähl-Wächter formuliert)
- `pfad`: `internal/adapter/driving/cli/print_mk.go:22-30`
- `befund`: Der Kommentar sagt „elf `##`-annotierte Targets" und „Das Template
  hat genau SECHS fmt-Verben … sonst bräche `fmt.Sprintf`". Das Template führt
  jetzt zwölf Targets, und `makefileFragment` übergibt sieben Argumente an
  ebenso viele Verben. Der Kommentar wurde um die Zeile `doc-structure/`
  ergänzt, die beiden Zahlen nicht. Wer sich beim nächsten Target auf den
  Wächter verlässt, zählt gegen einen falschen Sollwert.
- `verifizierbar`: nein — kein Gate deckt Prosa-Zahlen in Kommentaren;
  durch Lesen von `print_mk.go` belegbar.
- `klasse`: „Zähl-Kommentar als Wächter formuliert, aber nicht mitgezogen"

### F-10 — Die `planning`-Zeile der Handbuch-Modultabelle führt den neuen Grund-Code nicht

- `kategorie`: LOW
- `quelle`: [`MR-025`](../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
  Spiegel „Nutzer-Doku: Modul-Tabelle";
  [slice-099](../plan/planning/in-progress/slice-099-structure-modul.md) §3a
  Zeile Benutzerhandbuch
- `pfad`: `docs/user/benutzerhandbuch.md:1660`
- `befund`: Die Zeile behauptet in der Prosa „fail-closed bei
  fehlender/mehrdeutiger Überschrift", listet in der Grund-Code-Spalte aber nur
  `planning-drift`, `closure-note-missing`, `closure-note-thin`,
  `closure-note-boilerplate` und `closure-note-placeholder`. Die §4-Tabelle
  (`docs/user/benutzerhandbuch.md:508`) und die §11-Zeile
  (`docs/user/benutzerhandbuch.md:1841`) führen `closure-note-ambiguous`, die
  Modultabelle nicht — die drei Handbuch-Enumerationen desselben Codes stehen
  auf zwei verschiedenen Ständen. Der Defekt bleibt bestehen, auch wenn F-1
  behoben wird.
- `verifizierbar`: nein — Doku-Vergleich innerhalb einer Datei.
- `klasse`: „Grund-Code in nur einer von mehreren Handbuch-Enumerationen
  nachgezogen"

### F-11 — Der Kategorien-Satz der deutschen README trägt eine doppelte Aufzählungs-Konjunktion

- `kategorie`: LOW
- `quelle`: Maintainability (Doku-Drift in einer ausgelieferten
  Nutzer-Oberfläche, von diesem Diff eingeführt)
- `pfad`: `README.de.md:15-19`
- `befund`: Der Satz lautet nach der Änderung „… bis zu Versions-Pin-,
  Commit-Traceability-, Planning-Lifecycle- und Getrackt-Status-Konsistenz bis
  zu Struktur-Invarianten **innerhalb** eines Dokuments:". Die Konstruktion
  „von … bis zu … bis zu" liest sich als zwei konkurrierende Aufzählungen; der
  Satz ist die erste inhaltliche Aussage der README über den Modul-Umfang. Die
  englische Fassung (`README.md:16-20`) hat das Problem nicht, sie hängt
  „up to structure invariants" an ein Komma.
- `verifizierbar`: nein — kein Gate prüft Satzbau.
- `klasse`: „Doku-Drift in der READMEs-Statuszeile"

### F-12 — Symlink-Kandidaten werden still übergangen, ohne dass eine Vertragsfläche das sagt

- `kategorie`: INFO
- `quelle`: [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 2 (Kandidaten-Menge)
- `pfad`: `internal/hexagon/core/rules/scan.go:101-126`
- `befund`: `walkMarkdown` behandelt nur `KindDir` und `KindFile`; eine
  Markdown-Datei, die über einen Symlink erreichbar ist, und ein
  Verzeichnis-Symlink erscheinen nicht in der Kandidaten-Menge. Im Fixture
  bleibt `docs/link.md` (Symlink auf `docs/a.md`) ungeprüft, ohne Befund und
  ohne Hinweis. Ein stilles Grün entsteht daraus nicht — trifft eine Regel
  ausschließlich Symlinks, greift die Nullmengen-Härte —, aber eine gemischte
  Kandidaten-Menge prüft den Symlink still nicht. Für das Modul `tracked` hat
  [`DC-FA-TRK-001`](../../spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)
  genau diese Frage ausdrücklich beantwortet; für `structure` steht sie nirgends.
- `verifizierbar`: ja — Fixture mit Datei- und Verzeichnis-Symlink unter einem
  `files`-Glob.
- `klasse`: „undokumentierte Kandidaten-Ausnahme im Baum-Walk"

### F-13 — `CheckStructure` verwirft einen Walk-Fehler und liefert einen leeren Befundsatz

- `kategorie`: INFO
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  §fail-closed
- `pfad`: `internal/hexagon/core/rules/structure.go:21-24`
- `befund`: Schlägt `structureTree` fehl, kehrt die Funktion mit `nil` zurück —
  das Modul meldet dann nichts, obwohl jede andere Randbedingung desselben
  Moduls (unlesbare Datei, Nullmenge, Config-Fehler) fail-closed meldet oder
  abbricht. Der Zweig ist heute praktisch unerreichbar, weil der Scanner die
  Repo-Wurzel zuvor öffnet; er wird erreichbar, sobald ein Filesystem-Port
  Teil-Lesefehler durchreicht. Damit ist der einzige stille Pfad des Moduls
  genau der, den kein Test und kein Gate berührt.
- `verifizierbar`: nein — nicht über die CLI auslösbar, nur durch Lesen bzw.
  einen Port-Doppelgänger im Test.
- `klasse`: „Fehler-Rückgabe im Post-Pass verworfen statt gemeldet"

## Negativbefunde

- geprüft, ohne Befund: **Refactor-Äquivalenz der ausgelieferten
  Closure-Fähigkeit.** Ein aus `64c62cb` gebautes Image und der Arbeitsbaum-Bau
  liefern über vierzehn Randfälle (normal, thin, Abschnitt fehlt, CRLF,
  Setext-Überschrift, Inhalt nur im Fence, ungeschlossener Fence, Tilde- und
  eingerückter Fence, Abschnitt am Dateiende, Platzhalter, Floskel,
  verschachtelte Ebene, Überschrift mit Inline-Code, zwei Überschriften)
  byte-identische `--json`- und `--doctor`-Ausgaben; ebenso für fehlendes
  Closure-Verzeichnis, null Kandidaten und den vollständigen Default-Lauf über
  denselben Baum.
- geprüft, ohne Befund: **Schritt 2 (Kandidaten-Menge).** `**`-Globs greifen;
  Dateien außerhalb `scan.roots` werden geprüft; Dateien unter `scan.ignore`
  werden geprüft (wie zugesagt: unabhängig vom globalen Scan-Scope); die
  `SKIP_DIRS` gelten weiter (Regeln auf `.git/` und `node_modules/` laufen in
  `section-missing`); nur `.md` ist Kandidat; Sortierung stabil;
  `exempt-paths` nimmt einzelne Dateien aus, ohne die Nullmengen-Härte
  auszuhebeln.
- geprüft, ohne Befund: **Schritt 4 (Kardinalität).** `one` meldet
  `section-ambiguous` mit der Zeile des **zweiten** Treffers und bricht für
  diese Datei ab; `each` prüft jeden Treffer und meldet bei drei Abschnitten
  genau den verletzenden — beides in Fixture und Unit-Test.
- geprüft, ohne Befund: **Schritt 6 (Bedingungen).** Alle sechs mit eigenem
  Grund-Code; zwei Verletzungen desselben Abschnitts stehen nebeneinander
  (Dedup-Schlüssel in `internal/hexagon/core/model/finding.go:97-100` enthält
  die Meldung nicht, das frühe `break` in der Marken-Schleife ist damit
  verhaltensgleich); Task-Items brauchen den Listen-Marker und zählen im Fence
  nicht; `^` und `$` binden an Text-Grenzen, `(?m)` schaltet auf Zeilen, ein
  Muster über die Zeilengrenze trifft.
- geprüft, ohne Befund: **Schritt 7 (Befund-Form).** `file`, `line`, `rule` und
  `target` entsprechen der Zusage; die Regel-Identität aus Glob und Selektor
  trennt zwei Regeln über derselben Datei; identische Identität bricht mit
  Exit 2 ab.
- geprüft, ohne Befund: **Schritt 1 (Config-Rand).** Alle acht aufgezählten
  Fehlerfälle liefern Exit 2 vor dem Lauf, einschließlich eines ungültigen
  Globs in einem **späteren** Pfadsegment; leere Regel-Liste ist inert.
- geprüft, ohne Befund: **Schritt 8 (Determinismus, read-only).** Wiederholte
  Läufe liefern identische, sortierte Befundsätze; alle Läufe dieses Reviews
  liefen mit `--network none` und `:ro`-Mount.
- geprüft, ohne Befund: **Preset-Kopplungs-Test nicht tautologisch.** Driftet
  man **nur** den Closure-Konsumenten (Abschnitts-Ebene `level` auf `level+1`
  in `internal/hexagon/core/rules/planning.go:134`), wird genau
  `TestStructurePresetKopplungMitClosure` rot — der Test bindet also die
  gemeinsame Semantik, nicht nur den gemeinsamen Funktionsnamen.
- geprüft, ohne Befund: **Mutations-Echtheit der behaupteten Rückbauten.**
  Sieben selbst nachgestellt, alle rot: Task-Item ohne Listen-Marker,
  Marken-Grenze aufgehoben, `section-ambiguous` meldet den ersten statt den
  zweiten Treffer, `exempt-paths` wirkungslos, Regel-Identität ohne Selektor,
  Default-Kardinalität `each` statt `one`, `SectionEnd` nur bei höherer Ebene.
- geprüft, ohne Befund: **Messung der Selbst-Aktivierung (§4a).**
  `make verify-closure-notes` läuft mit den drei aufgenommenen Regeln über 337
  Dateien mit null Befunden. Die verworfene Regel „abgeschlossener Slice ohne
  offene Task-Boxen" habe ich unabhängig nachgemessen und komme auf **exakt
  32** — die Zahl ist reproduzierbar und die Begründung trägt. Die drei
  ADR-Zahlen liegen bei mir je um eins höher (16/21/15 statt 15/20/14); die
  Differenz erklärt sich vollständig durch `docs/plan/adr/README.md`, das mein
  weiterer Glob mitnimmt.
- geprüft, ohne Befund: **SemVer-Einstufung.** Minor ist korrekt: additives
  Modul, additive Grund-Codes, nichts entfernt, Default-Befundsatz
  byte-identisch. Der einzige als schärfend angekündigte Punkt tritt nicht ein
  (F-2) — das macht den Minor nicht falsch, sondern die Notiz.
- geprüft, ohne Befund: **Vertrags-Flächen außer den in F-5/F-6/F-10 genannten.**
  `DC-FA-CLI-002`-Modul-Liste, `validModules()`, §2-Schema,
  §4-Grund-Code-Tabelle (neun Zeilen), `AllReasons()`, `reasonTexts()`,
  `--print-config`-Gerüst, `--print-mk` (zwölf Targets, `doc-structure` mit
  `--enable structure` und fokussierter `--disable`-Liste, ohne Range),
  `docs/user/operations.md`-Optionstabelle, Handbuch-Modultabelle und §5, beide
  README-Modullisten, `AGENTS.md`- und `harness/README.md`-Gate-Beschreibungen,
  `version.md`, `CHANGELOG.md`, Handbuch-Kopfstempel und ghcr-Pins auf
  `v0.57.0`.
- geprüft, ohne Befund: **`--doctor`, `--json`, `--yaml`.** Alle neun neuen
  Codes haben Klartexte; die Diagnose gruppiert nach Datei und zeigt die
  Regel-Identität als Stelle; `--yaml` und `--json` führen dieselben Schlüssel.
- geprüft, ohne Befund: **ADR-0049.** Status bereits `Accepted`, Bezug,
  Konsequenzen und Re-Evaluierungs-Trigger vollständig; die Fitness Function
  ist bis auf den Mehrdeutigkeits-Punkt (F-1) erfüllt.
- geprüft, ohne Befund: **Architektur und Import-Regeln.** `structure.go` und
  `sections.go` liegen im Kern `internal/hexagon/core/rules` und importieren nur
  `model` und `port/driven`; `make arch-check` innerhalb `make gates` grün.
- geprüft, ohne Befund: **`FOCUS_DISABLE`.** Muss `structure` nicht führen — das
  Modul steht nicht in der `modules`-Liste der `.d-check.yml`, die
  Fokus-Liste spiegelt genau diese acht.
- geprüft, ohne Befund: **Gate-Läufe.** `make lint` (0 issues), `make test`,
  `make gates`, `make adr-check`, `make verify-closure-notes` — alle grün auf
  dem unveränderten Arbeitsbaum.
- geprüft, ohne Befund: **Sehr große Dateien.** Das Modul liest Kandidaten über
  denselben Filesystem-Port wie der Scanner und kennt keine eigene Grenze; das
  entspricht allen bestehenden Modulen, und der `DC-QA-01`-Benchmark deckt den
  Scan-Pfad ab. Kein eigener Befund.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 7 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Grund-Code deklariert und dokumentiert, aber
von keinem Codepfad erzeugt · Release-Notiz sagt Verhaltensänderung zu, die der
Lauf nicht zeigt · ASCII-Bytegrenze als Alphanumerik-Grenze · Preset-Partner-
Spezifikationen driften an derselben geteilten Mechanik · Zahl im
Beschreibungstext nachgezogen, Enumeration in den Akzeptanzkriterien nicht ·
neues Modul nicht in der Nicht-aktiviert-Enumeration von `--suggest-config` ·
Schwellen- und fail-closed-Ränder eines neuen Moduls ohne Test · Spiegel-Liste
nennt einen Spiegel als unberührt, der berührt wird · Zähl-Kommentar als
Wächter formuliert, aber nicht mitgezogen · Grund-Code in nur einer von
mehreren Handbuch-Enumerationen nachgezogen · Doku-Drift in der
READMEs-Statuszeile · undokumentierte Kandidaten-Ausnahme im Baum-Walk ·
Fehler-Rückgabe im Post-Pass verworfen statt gemeldet

## Verdikt

**Merge-blockierend:** ja — ein HIGH und sieben MEDIUM.

**Release-Empfehlung: v0.57.0 nicht taggen.** Der Grund ist nicht das Modul
`structure`: dessen Kern hält die acht Algorithmus-Schritte in allen von mir
geprüften Punkten bis auf F-3 und die Beschreibungslücke F-4, und der Refactor
an ausgeliefertem Code ist am Lauf gegen ein Vor-Stand-Image als
verhaltens-erhaltend belegt. Blockierend ist die zweite Hälfte des Slice: die
Mehrdeutigkeits-Härte ist auf **jeder** Vertrags-, Doku- und Release-Fläche
angekündigt und in **keiner** Codezeile umgesetzt (F-1), und drei ausgelieferte
Flächen behaupten sie gegenüber Konsumenten als vollzogen (F-2). Ein Release in
diesem Zustand verkauft eine Härte, die der Lauf nicht zeigt — dieselbe Klasse
von Harness-Lüge, gegen die die Grund-Code-Verriegelung angetreten ist, nur
eine Ebene tiefer: die Verriegelung prüft Katalog-Deckung, nicht Erreichbarkeit.

Die drei MEDIUM an den Vertragsflächen (F-5, F-6) und die vier überlebenden
Rückbauten (F-7) sind vor dem Tag zu klären, blockieren aber für sich genommen
nicht das Modul, sondern die Vollständigkeit seiner Zusage. F-8 ist der
lehrreiche Befund für die Regel selbst: `MR-025` hat bei ihrer ersten Anwendung
neun Spiegel korrekt getroffen und drei verfehlt — zwei davon sind hier zu
Findings geworden, und einen hat die Liste ausdrücklich als „unberührt"
deklariert, bevor er berührt wurde. Die Regel wirkt; ihre Liste ist noch nicht
vollständig.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in
die Slice-Closure §7 und von dort in den Steering-Loop-Zähler. Dieser Report
ist ein Lauf-Beleg (dieser Diff, dieser Skill, dieses Modell, dieses Verdikt)
und ersetzt keine Verifikation.
