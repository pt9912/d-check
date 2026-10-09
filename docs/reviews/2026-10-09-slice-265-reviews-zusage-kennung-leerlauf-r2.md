# Review R2: slice-265, `reviews` (Zusage-Muster, benannte Kennungen, Leerlauf und Unterverzeichnisse)

- **Review-Art:** Code. Geprüft wird der Fix-Commit gegen den Slice-Plan (`slice-265` §1, §3 mit den
  beiden Plan-Änderungen nach R1, §6, §8), gegen die Befunde F-1 bis F-11 aus R1, gegen die
  Entscheidungen (`ADR-0081`, `MR-025`, `MR-035`/`MR-036`) und gegen die Hard Rules `AGENTS.md`
  §3.4/§3.7/§3.8 sowie §5 Regel 13/15. Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-265` · Fix-Commit `311fdc40`. Als Kontext gelesen: `7d538243` und `578231ee`
  (die zwei Plan-Änderungen) sowie `81faa360`/`c4072ef7`/`d52cac33` (ausgehender CR und Antwort).
- **Skill:** `reviewer.md` @ 1.19.0 (`980d6314`)
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** R1-Report `2026-10-09-slice-265-reviews-zusage-kennung-leerlauf-r1.md`;
  `DC-FA-RVW-001` im Lastenheft 0.101.0; Spezifikation `DC-FA-RVW-001.a` Schritte 2–5, Grenzen-Absatz,
  §2-Zeilen `reviews.*`, `SPEC-081`, §8-Historie; ausgehender CR
  `docs/plan/cr/2026-10-09-cr-ai-harness-course-reviews-zusage.md` und Antwort; Register-Zeilen
  `BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet` (2 Belege) und
  `BEO-ALL/review-fix-applied-only-at-cited-site` (2 Belege). Im Code gelesen:
  `internal/hexagon/core/rules/reviews.go` (vollständig), `reviews_muster_test.go`,
  `internal/hexagon/core/model/config.go`, `internal/adapter/driven/configyaml/configyaml.go`
  (`rawReviews`, Fehlermeldungen von `applyReviews`), `internal/adapter/driving/cli/config_template.go`,
  `harness/sensors/review-coverage.md` (vollständig), `Makefile` (`review-coverage`).
- **Proben** (Wegwerf-Worktree auf `311fdc40` im Scratchpad, `make test IMAGE=dcprobe-r2`; danach
  Worktree und Image entfernt, Arbeitsbaum unverändert):
  1. Probetest `match: name`, je ein Slice mit Zusage und ein Report:
     `slice-a-grö` / `…-slice-a-größe-r1.md` **gedeckt**; `slice-a` / `…-slice-a_b-r1.md` gedeckt;
     `slice-a` / `…-slice-a.b-r1.md` gedeckt; `slice-a-foo` / `…-slice-a-foo.md` gedeckt;
     `slice-a-foo` / `…-slice-a-foox-slice-a-foo-r1.md` gedeckt (zweites Vorkommen, richtig);
     `slice-a-foo` / `…-preslice-a-foo-r1.md` gedeckt (links keine Grenze).
  2. Probetest Default-Muster, Datei `slice-001-a.md`, `reviews-dir` ohne passenden Report:
     `- [x] Adaptions-Review durchgeführt (intern, ohne Report)` → 1 Befund;
     `- [x] kein Review durchgeführt (entfallen)` → 1; `- [ ] Review durchgeführt?` → 1;
     `- [x] Kein unabhängiger Review nötig` → 1; `- [x] Review **durchgeführt**` → 0;
     `- [x] review durchgeführt` → 0.
  3. Mutationen in `reviews.go`: `isSkipDir` ausgeschaltet und die Wortgrenze ausgeschaltet
     (`if true || after == "" …`). Ganze Suite: `TestReviewsRecursive_SkipDirsBleibenUnbetreten` und
     `TestReviewsMatch_NameWortgrenze` werden rot, die übrigen Pakete bleiben `ok`.
  4. Bestand: `grep` auf Checkbox-Zeilen mit „Review durchgeführt" unter `docs/plan/planning/done/`
     (rekursiv) ergibt 26 Treffer. Formen `X-Review durchgeführt`, `kein Review durchgeführt` und
     `Review durchgeführt?` kommen im eigenen Bestand und in den vendored Templates nicht vor.

## Stand der R1-Befunde

| R1 | Stand | Beleg |
|---|---|---|
| F-1 (Walk-Abbruch, Doppelmeldung) | **geschlossen** | `visit` ohne Abbruch, `badDirs` je Verzeichnis, `case len(badDirs) > 0:` unterdrückt beide Leerlauf-Befunde; zwei Tests, einer davon mit `require-promises`. Spezifikation Schritt 2/5 und Lastenheft-Kriterium nachgezogen |
| F-2 (Kommentar-Chronik) | **geschlossen** am gemeldeten Ort; ein neuer Grenzfall siehe F-6 unten | Mess-Label, `slice-138` und „Die Abnahme des CR" entfernt, „neuen" aus dem Testkommentar entfernt |
| F-3 (`SKIP_DIRS`-Test) | **geschlossen** | Probe 3: die Mutation macht den neuen Test rot |
| F-4 (Byte-Identität ohne `reviews`) | in der Sache **überholt**: der Default hat sich geändert, die Byte-Identität ist nicht mehr der Vertrag. Siehe F-2 unten | Botschaft `311fdc40` |
| F-5 (Muster trifft Fence, Abschnitt, Verneinung) | **teilweise**: Fence und Abschnitt sind benannt und getestet, die Verneinung weder benannt noch getestet. Mit dem neuen Default trifft die Lücke jetzt jeden Nutzer ohne `promise-pattern`. Siehe F-1 unten | Probe 2 |
| F-6 (`null` still auf Default) | **geschlossen durch Entscheidung**: als „gilt als abwesend" in Lastenheft und Spezifikation benannt, mit Test. Ein Restpunkt in der Vorlage siehe F-5 unten | `TestDecode_ReviewsPromisePatternNullIstAbwesend` |
| F-7 (Begründung „dieselbe Semantik") | **geschlossen** | Kommentar nennt jetzt die ABGRENZUNG, und sie trifft am Code zu |
| F-8 (`rawReviews`-Kommentar) | **geschlossen** | alle acht Schlüssel genannt |
| F-9 (eingefrorene Zahl) | **geschlossen** | Grenze 2 nennt das Kommando-Verfahren statt der Zahl |
| F-10 (Präfix-Deckung `slice-26`) | **geschlossen** für Ziffern und ASCII-Buchstaben; der neue Vertragssatz greift weiter als der Code. Siehe F-4 unten | Probe 1, Probe 3 |
| F-11 („fail-closed"-Begründung) | **geschlossen** | Rückfall entfernt, `MustCompile` |

## Findings

### F-1 (MEDIUM): Der neue Default erkennt die Vorlagen-Form ohne linke Wortgrenze und auch verneint, und der Negativtest liegt nicht nahe am Positivfall

- **quelle:** `DC-FA-RVW-001` (Lastenheft 0.101.0, „Warum diese beiden Formen als Default": „Bloßes
  „Review" wäre zu breit … („Adaptions-Review") ohne externen Report"; Kriterium „Negative
  (Default)"); Skill-Fragen 13, 18 und 20; R1 F-5; `BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`
- **pfad:** `internal/hexagon/core/model/config.go` · „``const DefaultPromisePattern = `[Uu]nabhängiger Review|Review durchgeführt` ``“;
  `reviews_muster_test.go` · „`- [x] Adaptions-Review notiert`“
- **befund:** Die zweite Alternative hat links keine Grenze. Deshalb ist
  „Adaptions-Review durchgeführt“ eine Zusage, also genau das Konzept ohne Report, mit dem die
  Default-Wahl begründet wird. Dasselbe gilt für „kein Review durchgeführt (entfallen)“ und für
  „Review durchgeführt?“ (Probe 2: je ein `review-missing`). Seit dem Fix trifft das jedes Repo mit
  aktivem `reviews` ohne `promise-pattern`, nicht nur ein opt-in-Muster. Die Grenzen in Lastenheft und
  Spezifikation nennen Fence und Abschnitt, aber weder Verneinung noch Komposita. Der
  Negativtest prüft „Adaptions-Review notiert“ und „Review-Report liegt vor“, und beide können die neue
  Alternative gar nicht treffen. R1 F-5 hatte die Verneinung ausdrücklich genannt, behoben wurde nur
  der Fence-Teil.
- **Failure-Szenario:** Ein Adopter führt `reviews` mit dem Default, seine DoD trägt
  „- [x] Adaptions-Review durchgeführt“ oder einen entfallenen Punkt „kein Review durchgeführt“. Nach
  dem Update auf das Release meldet das Gate `review-missing`, obwohl kein Report zugesagt ist. Der
  Lastenheft-Text versichert ihm das Gegenteil.
- **verifizierbar:** ja, Probe 2 als Test.
- **klasse:** `erkennungs-regex-nur-gegen-positivfaelle-getestet` · `review-fix-applied-only-at-cited-site`

### F-2 (MEDIUM): Der Plan verspricht weiter Byte-Identität ohne die neuen Schlüssel, aber die zweite Plan-Änderung hat den Default geändert

- **quelle:** Slice-Plan `slice-265` §1 Ziel („mit opt-in-Schlüsseln, ohne die der Befundsatz
  byte-identisch bleibt“), §2 („Ohne die neuen Schlüssel ist die Ausgabe unverändert
  (`make blackbox-probe`)“), Kopf `**Bezug:**` („Default der Phrase bleibt“), §3-Zeile
  `review-coverage.md` („die Default-Phrase trifft keinen DoD-Punkt dieses Repos“), §8 Entwurf („die
  einzige Änderung am Default-Verhalten“); `AGENTS.md` §6 Schritt 4; Skill-Fragen 8 und 17
- **pfad:** `docs/plan/planning/in-progress/slice-265-reviews-zusage-kennung-leerlauf.md` ·
  „ohne die der Befundsatz byte-identisch“
- **befund:** Die Plan-Änderung `578231ee` streicht nur den Abgrenzungs-Punkt in §1. Ziel, DoD,
  Kopf, eine §3-Zeile und der §8-Entwurf behaupten weiter, dass sich ohne die neuen Schlüssel nichts
  ändert. Der Code zeigt das Gegenteil: Vorher lautete der Default `[Uu]nabhängiger Review`, jetzt
  erkennt er zusätzlich „Review durchgeführt“. Der Test `TestReviewsDefault_ErkenntVorlagenForm` wird
  ohne jeden Schlüssel rot. Die F-4-Messung der Botschaft („altes gegen neues Image … alle gleich“)
  stimmt, ist aber so gebaut, dass sie die Änderung nicht sehen kann. Auf dem Repo haben alle Slices
  mit der Vorlagen-Form einen Report: das alte Image prüft dort 0 Zusagen, das neue mehr als 0, und
  beide geben dieselbe Ausgabe aus. Der Fremd-Bestand trug die Vorlagen-Form nicht.
- **Failure-Szenario:** Bei der Closure hakt jemand den DoD-Punkt „Ausgabe unverändert“ mit dieser
  Messung ab. In der Release-Prep übernimmt jemand die übliche Handbuch-Formel „Opt-in, keine Breaking
  Changes“. Ein Bestand mit „Review durchgeführt“ und aktivem `reviews`, darunter das Kurs-Beispiel
  laut Antwort auf den CR, wird nach dem Update rot, ohne dass die Release-Notiz es ankündigt.
- **verifizierbar:** ja, der Default-Test oben bzw. ein Lauf von altem gegen neues Image auf einem
  Bestand mit „- [x] Review durchgeführt“ und leerem `reviews-dir`.
- **klasse:** `plan-aenderung-nur-an-zitierter-stelle` · `botschaft-ueberdehnt-messung`

### F-3 (MEDIUM): Die Spiegel des Defaults sind nicht nachgezogen, darunter eine Nutzer-Meldung und der Vertrag der eigenen Sensor-Datei

- **quelle:** `MR-025` (Spiegel vor dem Editieren); Slice-Plan §8 (Spiegel-Grep mit dem Muster
  „unabhängiger Review“); `AGENTS.md` §5 Regel 13; `BEO-ALL/review-fix-applied-only-at-cited-site`
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go` · „(weglassen ⇒ Phrase „unabhängiger
  Review“)“; weitere Stellen: `internal/hexagon/core/rules/reviews.go` · „(abwesend: die Phrase
  "unabhängiger Review")“; `Makefile` · „jede DoD-Zusage "unabhängiger Review" braucht“;
  `harness/sensors/review-coverage.md` §Vertrag · „ein DoD-Haken, dessen Zeile „unabhängiger Review"
  nennt“
- **befund:** Der Default hat sich geändert, aber vier Stellen nennen weiter nur die Phrase. Eine
  davon ist die Exit-2-Meldung, die der Nutzer bei leerem `promise-pattern` sieht. Eine andere ist der
  §Vertrag derselben Sensor-Datei, deren Grenze 2 im selben Commit sagt, der Default erkenne
  „Review durchgeführt“. Die Datei widerspricht sich damit selbst. Der Spiegel-Grep aus §8 trifft alle
  vier Stellen, er ist für die Default-Änderung offenbar nicht neu gelaufen. README, README.de und
  Handbuch sind nach §1 Release-Prep und nicht gemeint.
- **Failure-Szenario:** Ein Nutzer bekommt Exit 2 für `promise-pattern: ''`, folgt der Meldung und
  lässt den Schlüssel weg. Danach meldet das Gate DoD-Punkte, die laut Meldung nicht erkannt würden.
  Wer die Sensor-Datei liest, findet im Vertrag eine andere Erkennung als in der Grenze darunter.
- **verifizierbar:** ja, der Grep aus §8 des Plans.
- **klasse:** `spiegel-unvollstaendig`

### F-4 (LOW): „weder Buchstabe noch Ziffer“ gilt im Code nur für ASCII, und die Grenze nennt nur den Bindestrich

- **quelle:** `AGENTS.md` §5 Regel 13 und 15; `DC-FA-RVW-001` (Beschreibung und Kriterium „Boundary
  (Slug-Kennung)“); Spezifikation Schritt 4 und §2-Zeile `reviews.match`
- **pfad:** `internal/hexagon/core/rules/reviews.go` · „`func isNameRune(b byte) bool {`“
- **befund:** Lastenheft, Spezifikation, Codekommentar und Commit-Botschaft sagen, nach dem Basisnamen
  müsse ein Zeichen stehen, das „weder Buchstabe noch Ziffer“ ist. Der Code prüft ein einzelnes Byte
  gegen `[A-Za-z0-9]`. Ein Umlaut nach dem Basisnamen gilt deshalb als Grenze, und `slice-a-grö` wird
  vom Report zu `slice-a-größe` gedeckt (Probe 1). Die benannte Grenze nennt nur den Bindestrich.
  Unterstrich und Punkt wirken genauso (`slice-a` wird von `slice-a_b`/`slice-a.b` gedeckt), und links
  vom Basisnamen gibt es überhaupt keine Grenze (`preslice-a-foo`). Für Slug-Kennungen sind diese
  Fälle selten, aber der Vertragssatz greift weiter als der Code.
- **verifizierbar:** ja, Probe 1.
- **klasse:** `vertragssatz-weiter-als-code`

### F-5 (LOW): Die `--print-config`-Vorlage sagt „leer ⇒ Exit 2“, ohne den Wert ist der Schlüssel aber abwesend

- **quelle:** Spezifikation §2-Zeile `reviews.promise-pattern` („**Explizit** leer ⇒ Exit 2 … ohne Wert
  (YAML-`null`) ⇒ abwesend“); R1 F-6
- **pfad:** `internal/adapter/driving/cli/config_template.go` · „RE2 gegen den Punkt-Text ab dem Bullet;
  leer ⇒ Exit 2“
- **befund:** Die Spezifikation unterscheidet zwischen `''` und `null`. Die Vorlage sagt nur „leer“, und
  das liest man naheliegend als `promise-pattern:` ohne Wert. Genau diese Form läuft aber still mit dem
  Default.
- **verifizierbar:** nein (Doku).
- **klasse:** `null-wert-faellt-still-auf-default`

### F-6 (LOW): Neue Kommentare tragen die Slice-Nummer als Beispiel und ein Chronik-Wort

- **quelle:** `AGENTS.md` §3.7 („keine Slice-Nummern“); `harness/rules/kommentare-fuenf-klassen.md`
- **pfad:** `internal/hexagon/core/rules/reviews.go` · „slice-26 wird nicht vom Report zu slice-265-x
  gedeckt“; zweite Stelle `reviews.go` · „MustCompile trifft hier nur noch gueltige Muster“
- **befund:** Das Beispiel ist zwar eine Grenze, trägt aber die Kennung dieses Slice und den
  Probenfall aus R1. Ein neutrales Paar wie `slice-12`/`slice-123` sagte dasselbe ohne Herkunft.
  „nur noch“ beschreibt eine Änderung statt eines Zustands. Beides ist eine Grenzlage: der Kommentar
  hat eine der fünf Klassen, deshalb nicht HIGH.
- **verifizierbar:** nein (kein Gate prüft §3.7).
- **klasse:** `kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen`

### F-7 (LOW): Die Historie-Zeile 0.101.0 zählt die Kriterien ungenau

- **quelle:** `AGENTS.md` §5 Regel 15
- **pfad:** `spec/lastenheft.md` · „ersetzen die Vorgänger, die die alte Default-Annahme trugen“
- **befund:** Die Zeile sagt, drei Kriterien ersetzen ihre Vorgänger. Tatsächlich ist „Happy Path
  (Vorlagen-Form)“ neu hinzugekommen. Ersetzt wurden nur „Negative (Muster)“ und „Boundary (benannte
  Kennung)“. Die Zeile verschweigt außerdem, dass „Happy Path (eigenes Muster)“ inhaltlich umgeschrieben
  wurde (anderes Muster) und „Boundary (Unterverzeichnisse und Stubs)“ erweitert wurde.
- **verifizierbar:** nein.
- **klasse:** `historie-zeile-zaehlt-ungenau`

### F-8 (INFO): Neben einem unlesbaren Unterverzeichnis entfällt auch die Diagnose eines unlesbaren `reviews-dir`

- **quelle:** Spezifikation Schritt 5 („entfallen beide Leerlauf-Befunde“); Maintainability
- **pfad:** `internal/hexagon/core/rules/reviews.go` · „`case len(badDirs) > 0:`“
- **befund:** Sind ein Unterverzeichnis und `reviews-dir` gleichzeitig unlesbar und ist keine Zusage
  da, meldet der Lauf nur das Unterverzeichnis. Er bleibt rot, die zweite Ursache sieht man aber erst
  nach der ersten Reparatur. Das stimmt mit der Spezifikation überein und ist kein stilles Grün.
- **verifizierbar:** ja (MemFS mit beiden Fehlern).
- **klasse:** `diagnose-verdeckt-zweite-ursache`

### F-9 (INFO): `MustCompile` im Kern bricht bei einem ungültigen Muster ab, `planning` meldet im selben Fall

- **quelle:** Skill-Frage 11; Maintainability
- **pfad:** `internal/hexagon/core/rules/reviews.go` · „`promiseRE := regexp.MustCompile(cfg.EffectivePromisePattern())`“
- **befund:** Ein Aufrufer, der den Config-Rand umgeht, bekommt einen Panic statt eines Befunds. Über
  den Adapter ist der Pfad nicht erreichbar, der Kommentar nennt die Kopplung. Der Unterschied zu
  `planning` ist bewusst, aber nicht benannt.
- **verifizierbar:** nein.
- **klasse:** `module-gleicher-eingabe-verschieden`

## Negativbefunde

- **F-1-Fix (Walk):** `visit` läuft über alle Einträge weiter, `badDirs` ist sortiert, je Verzeichnis
  gibt es einen Befund mit `Target` = Pfad (`recursive` ist unreleased, also keine Vertragsänderung
  ggü. einem Release). Ohne `recursive` entsteht kein `badDir`. Ohne Befund.
- **Wortgrenze, mehrere Vorkommen und Namensende:** `rest = rest[i+1:]` findet auch ein späteres und
  überlappendes Vorkommen (Probe 1, `slice-a-foox-slice-a-foo`). Ein Basisname direkt vor `.md` deckt,
  weil `.` die Grenze ist. Ein leerer Rest deckt. Ohne Befund (Unicode siehe F-4).
- **Tests aus dem richtigen Grund rot:** `SKIP_DIRS` und Wortgrenze per Mutation (Probe 3). Der
  Vorlagen-Default-Test ist rot, wenn die zweite Alternative fehlt (am Code nachvollzogen: das alte
  Muster trifft „Review durchgeführt“ nicht). Ohne Befund.
- **Exit-2-Liste (Regel 13, am Code gezählt):** sechs Fälle in Code, Lastenheft und Spezifikation. Die
  `null`-Regel ist jetzt beidseitig benannt. Ohne Befund.
- **Lastenheft und Spezifikation gegen Code:** Default (beide Alternativen), Leerlauf-Ersatz neben
  `badDirs` (beide Leerlauf-Befunde), „ab dem Bullet“, Fence und Abschnitt, Schritt 2 „übrige Einträge
  werden weiter gelesen“, `SPEC-081`, §2-Zeilen: Sie stimmen mit dem Code überein. Ausnahmen sind F-1
  (Grenze unvollständig) und F-4.
- **Referenz-Richtung (`AGENTS.md` §3.4):** Der Spec-Diff enthält kein Slice-, Wellen-, ADR- oder
  Commit-Token. „Slice-Vorlage“ benennt ein Konzept, keine Kennung. Ohne Befund.
- **Hexagon, Netz, Suppression:** Es kommt kein Import hinzu, kein Netz, kein `//nolint`, keine
  gesenkte Schwelle. Ohne Befund.
- **Lese-Achse (`AGENTS.md` §3.8):** `reviews-dir` wird gelistet und nicht gescannt. Die Grenze ist
  unverändert benannt. Ohne Befund.
- **Ausgehender CR und Antwort (`MR-035`/`MR-036`):** Die Ablage liegt unter `docs/plan/cr/`. Die
  Antwort liegt neben dem CR, der Stand im CR-Kopf verweist auf sie. Die Ergänzung (Slug-Kennungen →
  `match: name`) ist in der `--print-config`-Vorlage aufgenommen. Ohne Befund.
- **Abgrenzung §1:** Handbuch und README sind nicht angefasst, `.d-check.yml` ebenfalls nicht. Die
  Default-Erweiterung ist vor dem Code als Plan-Änderung eingetragen (`578231ee`). Ohne Befund, bis
  auf die nicht nachgezogenen Plan-Stellen (F-2).
- **Commit-Botschaft (Regel 15):** Die Zuordnung F-1 bis F-11 stimmt mit dem Diff überein, und
  „Aendert den Befundsatz nur, wo ein Checkbox-Punkt diese Form traegt“ ist richtig. Überdehnt sind
  nur „weder Buchstabe noch Ziffer“ (F-4) und die Reichweite der F-4-Messung (F-2).

## Kategorie-Summary

HIGH 0 · MEDIUM 3 · LOW 4 · INFO 2. Wiederkehrend: `review-fix-applied-only-at-cited-site` (F-1: der
Fence-Teil von R1 F-5 ist behoben, die Verneinung nicht; F-2: die Plan-Änderung streicht nur eine
Stelle; F-3: die Spiegel der Default-Änderung fehlen) und `erkennungs-regex-nur-gegen-positivfaelle-getestet`
(F-1). Beide Register-Zeilen stehen bei 2 Belegen; mit diesem Slice erreichen sie die Schwelle.

## Verdikt

Nicht freigegeben. Der HIGH-Befund aus R1 (F-1 Walk) und der Kommentar-Befund F-2 sind in der Sache
geschlossen, ebenso F-3 und F-6 bis F-11. Vor dem Merge zu klären sind F-1 (der neue Default trifft
Komposita und Verneinung ohne Grenze und ohne nahen Negativtest), F-2 (der Plan verspricht noch
Byte-Identität, die Release-Notiz braucht die Verhaltensänderung) und F-3 (Spiegel des Defaults,
darunter eine Nutzer-Meldung). F-4 bis F-7 sind nicht blockierend.
