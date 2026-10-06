# Verifikation — slice-251: `targets.authority` nimmt eine Liste an

- **Rolle:** Verifier (Modul 8/11). Die Frage lautet „Bauen wir es richtig?“. Geprüft
  wird gegen die Plan-DoD (§2), die Abgrenzung (§1), die CR-Akzeptanzkriterien und
  `DC-FA-TGT-001`. Die Maintainability ist Sache von R1 und wird hier nicht geprüft.
- **Gegenstand:** `slice-251` @ `6c108f44`. Kette: `92dde5e8` feat → `34d8df6f` R1-Report →
  `6c108f44` R1-Einarbeitung.
- **Eingang:** Slice-Plan `docs/plan/planning/in-progress/slice-251-targets-authority-liste.md`
  (DoD-Punkt 2 mit vermerkter Plan-Änderung nach R1); eingehender CR
  `docs/plan/cr/2026-10-06-cr-eingehend-ai-harness-init-targets-authority-liste.md`;
  `spec/lastenheft.md` §`DC-FA-TGT-001` (0.96.1); `spec/spezifikation.md`
  §`DC-FA-TGT-001.a` Schritt 5, die Schema-Zeile `targets.authority` und `SPEC-061`;
  `ADR-0100` (Proposed); R1-Report `docs/reviews/2026-10-06-slice-251-targets-authority-liste-r1.md`.
- **Modell-ID:** claude-opus-5-5 · **Datum:** 2026-10-06
- **Arbeitsbaum:** Jede Mutation wurde per `git checkout -- internal` zurückgenommen, die
  Zusatz-Images `dcverify-before:251` und `dcverify-after:251` sind entfernt. Bis auf diesen
  Report ist `git status` leer.

## Selbst gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make gates` (HEAD, sauberer Baum) | Exit 0. Ausgabe `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`, `coverage-gate: OK — Coverage 94.60% erfüllt Schwelle 93%`. Der `targets`- und der `planning`-Dogfooding-Lauf melden je 0 Befunde |
| Vorher-Image | `git archive 92dde5e8^ \| docker build --target runtime -t dcverify-before:251 -` |
| Nachher-Image | `git archive HEAD \| docker build --target runtime -t dcverify-after:251 -` |
| `make test` × 4 Mutations-Sätze | siehe §Bewusstes Brechen |
| Black-Box, rund 60 Läufe je Image | siehe §Black-Box |

## Bewusstes Brechen (Modul 11)

Die Mutations-Sätze kombinieren Eingriffe in **verschiedenen** Paketen (`core/rules` und
`configyaml`), damit sich jeder Rot-Grund seiner Mutation zuordnen lässt. Zu jedem Lauf wurde
die Fehlerursache gelesen.

| Satz | Mutation | Rot in | gelesene Ursache | richtiger Grund? |
|---|---|---|---|---|
| A | `targets.go`: nur `authority[:1]` gelesen (keine Vereinigung) | `TestCheckTargetsAuthorityListe`, `…ListeFehlend` | `tool` fälschlich `gate-undocumented` (`in der Autoritäts-Doku harness/README.md`); fehlende zweite Datei: `bekam <nil>` | ja |
| A | `configyaml.go` ← `92dde5e8` (vor R1) | `TestDecode_TargetsAuthority`, 10 Fälle | `null`/`~` ⇒ `Authority = ["null"]`/`["~"]`; `harness/[b].md` ⇒ „ist ein Muster“; Alias ⇒ „muss ein Pfad …“; `-`, `[~]`, `[…, null]` ⇒ `bekam <nil>`; `""`/`[[a.md]]`/`{a: b}` ⇒ andere Meldung | ja (deckt H1 und M1 ab) |
| B | `targets.go`: path.Clean-Dublette ausgeschaltet (`if false && seenDoc[key]`) | **nichts** — `ok github.com/pt9912/d-check/internal/hexagon/core/rules` | — | **ungeschützt (V2)** |
| B | `configyaml.go`: `el.Tag == "!!null"` entfernt | `TestDecode_TargetsAuthority` | `[harness/README.md, null]` und `[~]` ⇒ `bekam <nil>` | ja |
| C | `targets.go`: Plural-Zweig der Meldung aus | `TestCheckTargetsAuthorityListe` | `Message = "… in der Autoritäts-Doku harness/README.md"` | ja |
| C | `configyaml.go`: String-Form über `n.Value` statt `Decode`, Alias-Auflösung oben aus, `authority` aus der Pfad-Regel in `applyTargets` genommen | `TestDecode_TargetsAuthority` | `~`/`null` ⇒ `["~"]`/`["null"]`; Alias ⇒ `["a"]`; `{a: b}` ⇒ `bekam <nil>`; `[../x.md]` ⇒ `bekam <nil>` (Pfad-Regel) | ja |
| D | `configyaml.go`: **nur** die Alias-Auflösung auf oberster Ebene (`:1696-1698`) aus | **nichts** — beide Pakete `ok` | — | **ungeschützt (V3)** |

Die Rot-Belege beider Commit-Botschaften habe ich nachgemessen. `92dde5e8`: Eine Fassung, die
nur die erste Datei liest, färbt `…Liste` und `…ListeFehlend` rot (Satz A). `6c108f44`: Gegen
die `configyaml.go` aus `92dde5e8` sind null/~/`[`/Alias rot, außerdem drei leere Elemente mit
`bekam <nil>` sowie `""` und `[[a.md]]` mit abweichender Meldung (Satz A). Zusätzlich ist
`{a: b}` rot, das die Botschaft nicht nennt. Beide Botschaften behaupten also nicht mehr, als
gemessen wurde.

## Black-Box (echtes Image, `--network none`, Fixture `chmod -R a+rX`)

Die Fixture enthält ein `Makefile` mit `own`, `tool`, `both` und `ghost`. Dazu kommen
`harness/README.md` (`own`, `both`), `harness/targets.md` (`tool`, `both`) und
`harness/[b].md` (`own`). Die Konfiguration ist `targets: {makefiles: [Makefile], authority: <Wert>}`.
stdout, stderr und Exit-Code wurden getrennt erfasst und byte-weise verglichen.

**String-Form, Vorher gegen Nachher:** Alle folgenden Werte sind **byte-identisch**:
`harness/README.md`, `"harness/README.md"`, `./harness/README.md`, `'harness/[b].md'`
(rc 1, 2 Befunde), `null`, `~`, `Null`, `""`, `''`, leerer Wert, fehlender Schlüssel (rc 0),
`"  "`, `123`, `true`, `!!str null`, `nope.md`, `/abs.md`, `../x.md`, ein unvollständiges `[…`
(rc 2), ein Alias `*a` auf einen Pfad-Scalar (rc 1). Für den String-Pfad und `[b]` gilt das auch
mit `--json`, `--doctor` und `--doctor --json`, für `null` mit `--json` und `--doctor`. Einzige
Abweichung ist `{a: b}`: Beide Images liefern rc 2, nachher steht jedoch der Präfix
`line 3: targets.authority muss ein Pfad oder eine Liste von Pfaden sein:` vor der
yaml-Meldung. Das ist ein Fehlerfall, kein gültiger String, und entspricht der R1-L2-Einarbeitung.

**Einelementige Liste gegen die alte String-Form:** `[README]`, `[README, ./README]` und
`[README, README]` sind im Nachher-Image byte-identisch zum String-Lauf des Vorher-Images, in
allen vier Ausgabeformen. Auch Dubletten in zwei Schreibweisen ergeben also den Singular-Wortlaut.

**Listenform (Nachher):**

| Wert | rc | Ergebnis |
|---|---|---|
| `[README, targets]` | 1 | genau `Makefile:7 ghost`, Meldung `in einer der Autoritäts-Dokus harness/README.md, harness/targets.md`; `both` (in beiden) ohne Befund |
| `[README, targets, README]` | 1 | wie oben, die Dublette wird still verworfen |
| `[README, nope.md]` / `[nope.md, README]` | 2 | `kann die Doku-Datei "nope.md" nicht lesen (DC-FA-TGT-001, fail-closed)`, unabhängig von der Position |
| `[]` | 0 | Richtung 2 entfällt |
| `[~]`, `[README, null]`, `[README, NULL]`, `[README, ""]`, `[README, "  "]`, `[[README]]`, `[{a: b}]`, Block-`-` allein, Block-`- README` + `-`, Alias-Element auf `&n ~` | 2 | `line N: targets.authority enthält einen leeren oder ungültigen Eintrag (erwartet: Pfad)` |
| `[README, /abs.md]`, `[../x.md]` | 2 | `targets-Pfad … muss relativ zur Repo-Wurzel liegen` |
| `[README, 'harness/[b].md']` | 1 | wörtlich gelesen, beide Dateien in der Meldung |
| `['harness/*.md']` | 2 | Laufzeit-Fehler „nicht lesen“, das Muster wird nicht expandiert |
| `[123]`, `[!!str null]` | 2 | als Pfad `123`/`null` gelesen, Laufzeit-Fehler (siehe V5) |
| Alias auf eine Sequenz, Alias-Element | 1 | aufgelöst, Vereinigung korrekt |
| `--doctor` | 1 | Klartext im Singular, `Hinweis:` nennt beide Dateien (ADR-0100 Entscheidung 5) |
| mit `exempt-targets: [ghost]` | 0 | Ausnahme greift weiter |
| `makefiles: [Makefile2, "harness/mk/*.mk"]` + Liste | 1 | `harness/mk/a.mk:5 frag`, Fundstelle im Fragment (CR-AK 6) |

## CR-Akzeptanzkriterien und Abgrenzung

| CR-AK | Verhalten | Test-Schutz |
|---|---|---|
| 1 String byte-identisch | erfüllt, einschließlich `null`/`~`/Alias/`[`/`--json`/`--doctor` (Black-Box) | `TestDecode_TargetsAuthority`, `…EinzelnWortlaut` (Satz A, C rot) |
| 2 Vereinigung | erfüllt | `…AuthorityListe` (Satz A rot) |
| 3 fehlende Datei Exit 2 | erfüllt, beide Positionen, mit Dateinamen | `…ListeFehlend` (Satz A rot) |
| 4 Doppelnennung kein Befund | erfüllt (`both`) | `…AuthorityListe` |
| 5 Pfad-Regel je Eintrag | erfüllt (`/abs.md`, `../x.md`) | `TestDecode_TargetsAuthority` (Satz C rot) |
| 6 rotes Gegenbeispiel mit Fundstelle im Fragment | erfüllt (`harness/mk/a.mk:5`) | `…AuthorityListe` (Datei und Zeile geprüft) |
| 7 Glob | entschieden: kein Glob, Einträge wörtlich auch mit Glob-Zeichen (ADR-0100 Entscheidung 3); `*` scheitert laut zur Laufzeit | `TestDecode_TargetsAuthority` (`[b]` wörtlich) |

Die offenen CR-Fragen sind beantwortet: Frage 1 durch das Listen-Beispiel im
`--print-config`-Gerüst (im Image gesehen, Zeilen 248–249), Frage 2 durch ADR-0100
Entscheidung 2. Zur CR-Abgrenzung und zu Plan §1: `doc-tables`, `exempt-targets` und die
Regel-Extraktion sind unverändert (der Diff berührt sie nicht, `exempt-targets` greift in der
Black-Box). Es gibt keine Disjunktheits-Prüfung, keinen Glob für `authority` und keine
Handbuch-, README- oder CHANGELOG-Änderung. Die Rückführungs-Bedingung aus §4 ist nicht
eingetreten, denn `decodeStrict` ist unverändert und die String-oder-Liste-Form lebt allein in
`decodeTargetsAuthority`.

## R1-Einarbeitung gegen den R1-Report

| R1 | Einarbeitung | Beleg |
|---|---|---|
| H1 (Null-Element still verworfen) | behoben: Prüfung je Knoten | Black-Box `-`, `[~]`, `[README, null]` ⇒ rc 2 mit Zeile; Satz A/B rot |
| M1 (String-Form nicht byte-identisch) | behoben: Decode als string, Ablehnung von Glob-Zeichen zurückgenommen; die Plan-Änderung ist in DoD-Punkt 2 vermerkt | Black-Box byte-identisch für `null`/`~`/`Null`/Alias/`[b]` |
| L1 (Kopfkommentar Singular) | behoben (`targets.go:45` Plural) | Diff |
| L2 (Zeilenangabe) | behoben: `{a: b}` meldet `line 3`. Beim Alias-Element auf ein Null-Anker zeigt die Zeile auf den Anker, nicht auf `authority` (V4) | Black-Box |
| I1 (`--doctor`-Singular undokumentiert) | als ADR-0100 Entscheidung 5 festgehalten | ADR |
| I2 (Null-Elemente in `makefiles`/`doc-tables`) | nicht angefasst, Bestand außerhalb von §1 | — |
| I3 (Release-Prep) | offen, zu Recht (AGENTS §5 Regel 17) | — |

Nach der Einarbeitung gab es keine zweite Review-Runde. DoD-Punkt 5 ist mit R1 formal erfüllt,
und die Einarbeitung ist oben verhaltensseitig nachgemessen.

## Verdikt je DoD-Punkt

| # | DoD-Punkt | Verdikt |
|---|---|---|
| 1 | Kern: Vereinigung, Meldung (eine Datei: heutiger Wortlaut, mehrere: alle), Dokument-Dubletten einmal gelesen, Tests nach CR-AK 1–6 mit rotem Gegenbeispiel, Rot-Beleg | **bestätigt mit Vorbehalt.** Das Verhalten ist vollständig erfüllt (Black-Box), die Tests für AK 1–6 laufen ohne die Änderung aus dem richtigen Grund rot. Die zugesagte Dubletten-Behandlung schützt aber kein Test, die Mutation überlebt (V2) |
| 2 | Config-Rand: String oder Liste; leerer/Null-/Nicht-Pfad-Eintrag Exit 2; Pfad-Regel je Eintrag; String-Form wie zuvor; `--print-config` mit Listen-Beispiel; Plan-Änderung vermerkt | **bestätigt.** Alle Fälle sind in der Black-Box und per Mutation belegt. Ungeschützt bleibt nur die Alias-Auflösung auf oberster Ebene (V3), funktional ist sie korrekt |
| 3 | Lastenheft (Bump + Historie), Spezifikation, Schema, ADR, CR-Antwort | **teilweise.** Lastenheft 0.96.1 mit Historie (`MR-032`), Spezifikation Schritt 5, Schema-Zeile, `SPEC-061`, ADR-0100 und Index-Zeile sind vorhanden. Die Beschreibung von `gate-undocumented` im Lastenheft ist jedoch durch eine doppelte Verneinung invertiert (V1). Die CR-Antwort steht noch aus, das ist erwartet und kein Befund |
| 4 | `make gates` grün | **bestätigt** (selbst gefahren, Exit 0, Coverage 94,60 %) |
| 5 | Review durchgeführt, Report liegt vor | **bestätigt** (R1-Report `34d8df6f`, anderer Kontext) |
| 6 | Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | **offen**, wie erwartet. §7 und §6 sind ungefüllt, ADR-0100 bleibt bis zur Closure `Proposed`. Kein Befund |

## Neue Befunde

### V1 — MEDIUM — Die Lastenheft-Definition von `gate-undocumented` ist durch eine doppelte Verneinung invertiert

- **pfad:** `spec/lastenheft.md:3261-3263`
- **befund:** Der Text lautet „eine Makefile-Regel (minus `targets.exempt-targets`), die in
  **keiner** Autoritäts-Doku (`targets.authority`) **nicht** als `make X` steht“. Eingeführt
  wurde das in `92dde5e8`, als „in der“ durch „in keiner“ ersetzt wurde, ohne das fett gesetzte
  „**nicht**“ zu streichen. Wörtlich bedeutet der Satz eine Regel, die in jeder Autoritäts-Doku
  steht, also das Gegenteil der Zusage. Der Absatz *Mehrere Autoritäts-Dateien*, Spezifikation
  Schritt 5 und der Code sind korrekt. Fehlerhaft ist nur die Definitionszeile des Grund-Codes
  im abnahmebindenden Rang-1-Dokument. R1 hat die Stelle als nachgezogenen Spiegel gelistet
  („Beschreibung ‚in keiner Autoritäts-Doku‘“), ohne den Restsatz zu lesen.
- **Behebung:** Entweder „**nicht**“ streichen oder „die in keiner Autoritäts-Doku … steht“ schreiben.
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand-geprueft`, hier am Spec-Text selbst

### V2 — LOW — Die path.Clean-Dublette in `undocumentedFindings` schützt kein Test

- **pfad:** `internal/hexagon/core/rules/targets.go:131-135`
- **befund:** DoD-Punkt 1 („Dokument-Dubletten einmal gelesen“) und Spezifikation Schritt 5
  („eine doppelt genannte (gleicher bereinigter Pfad) einmal“) sagen das Verhalten zu. Die
  Black-Box bestätigt es: `[README, ./README]` ergibt den Singular-Wortlaut, byte-identisch zur
  String-Form. Mit ausgeschalteter Prüfung bleibt `make test` aber grün (Satz B,
  `ok …/core/rules`). Ohne die Prüfung würde `[README, ./README]` den Plural mit beiden
  Schreibweisen melden. Das ist ein Meldungsbruch, den kein Test fängt. Es ist dieselbe Klasse
  wie `grenze-ohne-negativtest` aus R1 zu slice-250.
- **Behebung:** einen Fall `Authority: {"harness/README.md", "./harness/README.md"}` mit
  erwartetem Singular-Wortlaut in `TestCheckTargetsAuthorityEinzelnWortlaut` oder einem
  Geschwister-Test ergänzen.

### V3 — LOW — Die Alias-Auflösung auf oberster Ebene schützt kein Test

- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:1696-1698`
- **befund:** Der Alias-Testfall (`authority: *a` auf einen Scalar) läuft auch ohne diese drei
  Zeilen grün, weil `n.Decode(&s)` einen Scalar-Alias selbst auflöst (Satz D, beide Pakete
  `ok`). Wirksam ist der Zweig nur für einen Alias auf eine **Sequenz**. Den deckt kein Test,
  funktional ist er korrekt (Black-Box „alias seq“ rc 1). Ohne den Zweig wäre ein solcher Alias
  rc 2 („cannot unmarshal !!seq into string“), also laut und nicht still.
- **Behebung:** einen Fall `doc-tables: &l [a.md, b.md]` / `authority: *l` mit erwarteter
  Zwei-Element-Liste ergänzen.

### V4 — INFO — Bei einem Alias-Element zeigt die Fehlerzeile auf den Anker, nicht auf `authority`

- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:1705-1709`
- **befund:** `doc-tables: [&n ~]` (Zeile 3) und `authority: [harness/README.md, *n]` (Zeile 4)
  liefern `line 3: targets.authority enthält einen leeren …`. Nach `el = el.Alias` trägt
  `el.Line` die Zeile des Ankers. Das Verhalten ist laut und korrekt, nur der Ort ist ungenau.
  Ein Randfall.

### V5 — INFO — Lastenheft „Nicht-Pfad-Listeneintrag“ gegen Spezifikation „Nicht-Skalar-Element“

- **pfad:** `spec/lastenheft.md` (Absatz *Mehrere Autoritäts-Dateien*) und `spec/spezifikation.md` Schritt 5
- **befund:** Die Spezifikation und der Code lehnen beim Laden Nicht-**Skalare** ab. Ein
  Skalar wie `123`, `true` oder `!!str null` gilt als Pfad und scheitert erst beim Lesen, ebenfalls
  mit Exit 2. Das ist konsistent zur String-Form, deren Verhalten Vorher und Nachher
  byte-identisch ist. Der Exit-Code stimmt also mit der Lastenheft-Zusage überein, die
  Lade-Zeitpunkt-Lesart von „Nicht-Pfad“ dagegen nicht. Es geht nur um den Wortlaut.

## Negativbefunde

- **Byte-Identität der String-Form (CR-AK 1, `DC-QA-02`):** In der Black-Box vollständig
  bestätigt, einschließlich der R1-M1-Fälle. Geprüft, ohne Befund.
- **Stilles Grün:** Kein konfigurierter, aber leerer Eintrag legt Richtung 2 still. Geprüft
  wurden alle Null-, Leer- und Nicht-Skalar-Varianten in Flow- und Block-Form sowie per Alias.
  Ohne Befund.
- **Spezifikation Schritt 1:** Der Plan nennt Schritt 1 als berührt, der Diff ändert ihn nicht.
  Sein Wortlaut („nur bei nicht-leerem `targets.authority`“, „jede konfigurierte Datei … fehlend
  ⇒ Exit 2“) trägt die Listenform jedoch unverändert. Ohne Befund.
- **Hexagon-Richtung:** Der Kern nutzt nur `path`/`strings`, `yaml.Node` bleibt im Adapter;
  `make arch-check` ist grün. Ohne Befund.
- **ADR-0100 Grenzen gegen das Verhalten:** Entscheidung 3 („Muster scheitert laut zur
  Laufzeit“) ist bestätigt (`['harness/*.md']` ⇒ rc 2). Entscheidung 4 („leer heißt entfällt —
  nur auf oberster Ebene“) ist bestätigt: `[]`/`null`/`""` rc 0, Elemente rc 2. Entscheidung 5
  (`--doctor` Singular, Hinweis nennt alle) ist bestätigt. Die Fitness-Function-Tabelle nennt
  vier Tests, alle existieren und sind per Mutation rot-fähig. Ohne Befund.
- **Erwartet offen, kein Befund:** Closure-Notiz §7, Risiko-Ausgänge §6, ADR-Status `Proposed`,
  CR-Antwort und Release-Prep-Flächen (Handbuch §5/§6, README, CHANGELOG; AGENTS §5 Regel 17).

## Verdikt

**DoD-konform bis auf V1.** Das Verhalten erfüllt alle sieben CR-Akzeptanzkriterien und die
Abgrenzung, die R1-Einarbeitung ist nachgemessen, und `make gates` ist grün. V1 ist eine
Ein-Wort-Korrektur im Lastenheft und sollte vor der Closure erfolgen. V2 und V3 sind fehlende
Negativtests für korrekt arbeitenden Code. V4 und V5 dienen der Information.
