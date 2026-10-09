# Review R1 — slice-265: `reviews` — Zusage-Muster, benannte Kennungen, Leerlauf und Unterverzeichnisse (Produkt)

- **Review-Art:** Code. Geprüft wird der Produkt-Diff gegen den Slice-Plan (`slice-265`: §1 Ziel und
  Abgrenzung, §3, §6, §8 mit Entwurf und Spiegel-Liste), die Entscheidungen (`ADR-0081`, `MR-025`) und
  die Hard Rules `AGENTS.md` §3.4/§3.7/§3.8 sowie §5 Regel 13/15. Die DoD-Abhakung gehört nicht zu
  diesem Review.
- **Gegenstand:** `slice-265` · Produkt-Commit `1837d8d2` (Kontext: `86a0d887..97be2feb`, Plan,
  Übergang, Plan-Änderung).
- **Skill:** `reviewer.md` @ 1.19.0 (`980d6314`)
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** Slice-Plan `slice-265`; eingehender CR
  `docs/plan/cr/2026-10-09-cr-eingehend-sf-connector-reviews-zusage.md`; `DC-FA-RVW-001` im
  Lastenheft (0.100.0); Spezifikation §`DC-FA-RVW-001.a` Schritte 1–5, §2-Zeilen `reviews.*`, §4-Zeile
  `SPEC-081`, §8; vorherige Befunde am Modul bzw. an der parallelen Achse: Review R1 von `slice-263`
  (F-5 `SKIP_DIRS`-Negativtest, F-6 Byte-Identitäts-Behauptung, F-10 verdeckte Restmenge) und der
  H2-Regressionstest in `reviews_test.go` (widersprüchliche „leere Prüfmenge"-Meldung). Im Code
  gelesen: `internal/hexagon/core/rules/{reviews,planning,scan}.go`, `reviews_test.go`,
  `reviews_muster_test.go`, `internal/adapter/driven/fs/fs.go` (`List` sortiert),
  `internal/adapter/driven/configyaml/configyaml.go` (`rawReviews`, `applyReviews`),
  `internal/hexagon/core/model/config.go` (`ReviewsConfig`),
  `internal/adapter/driving/cli/config_template.go`, `harness/sensors/review-coverage.md`,
  `tools/blackbox-probe.sh`, `Makefile` (`blackbox-probe`, `review-coverage`), `.d-check.yml`.
- **Proben** (Wegwerf-Kopie im Scratchpad, `make test IMAGE=dcprobe` bzw. `make build IMAGE=dcprobe`
  und `docker run --network none` gegen die Kopie; Arbeitsbaum unverändert, Probe-Images danach
  entfernt):
  1. Test, `recursive: true`, unlesbares `done/archiv/` (sortiert **vor** `slice-*`), daneben flach
     `slice-001-a.md` mit Zusage und ohne Report. Ausgabe: `Verzeichnis …/archiv unlesbar — fail-closed`
     **und** `leere Pruefmenge: 0 Kandidat(en), 0 Review-Zusage(n), reviews-dir lesbar: true —
     fail-closed`; der fehlende Report von `slice-001` wird **nicht** gemeldet. Mit
     `require-promises` und einem weiteren Slice ohne Zusage: dieselben zwei Befunde.
  2. Test, `promise-pattern: 'Review durchgeführt'`: ein Checkbox-Punkt innerhalb eines
     Markdown-Fence in §9 und ein Punkt „- [x] kein Review durchgeführt (entfallen)"
     in §7 melden beide `review-missing`.
  3. Test, `match: name`: Slice `slice-26.md` mit Zusage, einziger Report
     `2026-slice-265-foo-r1.md` — kein Befund.
  4. Mutation `if !w.cfg.Recursive || (false && isSkipDir(name))` in `reviewWalk.descend`, ganze
     Suite: `ok github.com/pt9912/d-check/internal/hexagon/core/rules` — kein Test wird rot.
  5. Image gegen den Repo-Bestand: `promise-pattern:` ohne Wert und `promise-pattern: ~` laufen mit
     `0 Befund(e)` (die Phrase greift still), `promise-pattern: ''` bricht mit der Exit-2-Meldung ab.
     Die Messung aus `harness/sensors/review-coverage.md` reproduziert: Default-Phrase 0 Befunde,
     `require-promises` 1 Befund („keine Review-Zusage unter 31 Kandidat(en)"), CR-Muster mit leerem
     `reviews-dir` 26 Befunde, mit dem echten 0.

## Findings

### F-1 — HIGH — Ein unlesbares Unterverzeichnis bricht den ganzen Walk ab und erzeugt eine falsche Leermeldung

- **quelle:** `DC-FA-RVW-001` (Spezifikation Schritt 2: „ein **unlesbares** Unterverzeichnis ist ein
  Befund `review-missing` auf `reviews.done-dir` mit seinem Pfad in der Meldung"); Skill-Frage 2;
  vorheriger Befund am Modul: H2-Regressionstest `TestReviewsUnreadableReviewsDirWithPromisesNoRedundantFinding`
- **pfad:** `internal/hexagon/core/rules/reviews.go` · „`if w.badDir != "" {`" in `reviewWalk.visit`
  und der Kommentar „der erste unlesbare Unterordner beendet den Abstieg und steht in badDir"
- **befund:** `visit` kehrt nach dem ersten unlesbaren Unterverzeichnis auf **jeder** Ebene zurück;
  da `List` nach Namen sortiert, fallen alle später sortierten Einträge aus — auch flache
  `slice-*.md` neben einem Unterverzeichnis wie `archiv/` (Probe 1). Die Zusagen dort werden nicht
  geprüft, und weil die Menge dann leer ist, tritt neben den `badDir`-Befund ein zweiter mit „leere
  Pruefmenge: 0 Kandidat(en)", obwohl die Menge nicht leer, sondern ungelesen ist — genau die
  widersprüchliche Doppelmeldung, die der H2-Regressionstest für das unlesbare `reviews-dir`
  ausschließt. Der Kommentar sagt „beendet den Abstieg", der Code beendet die ganze Kandidatenwahl;
  die Spezifikation nennt den Abbruch nicht. `planning` bricht am selben Ereignis mit **einem**
  Befund ab und misst nichts Halbes.
- **verifizierbar:** ja — Probe 1 als Test (Unterverzeichnis-Name vor `slice-` sortiert).
- **klasse:** `fail-closed-verdeckt-restmenge` · `widerspruechliche-leermeldung`

### F-2 — HIGH — Angefasste Kommentare tragen Slice-Nummer, Mess-Label und Herkunfts-Prosa

- **quelle:** `AGENTS.md` §3.7; `harness/rules/kommentare-fuenf-klassen.md` §Bestandsgrenze („geräumt
  wird beim nächsten Anfassen der Zeile"); Skill-Frage 6
- **pfad:** `internal/hexagon/core/rules/reviews.go` · „(gemessen: mindestens sechs Faelle im Bestand,
  u. a." / „slice-138)" in `reviewPromise`; zweite Stelle
  `internal/hexagon/core/rules/reviews_muster_test.go` · „Die Abnahme des CR: die DoD-Zeile aus dem
  Bestand des Absenders"
- **befund:** Die beiden `reviewPromise`-Zeilen mit Mess-Label und Slice-Nummer sind im Diff neu
  umbrochen, also angefasst; die Bestandsgrenze endet damit, und §3.7 schließt beides ausdrücklich aus.
  Der neue Testkommentar nennt die Herkunft des Falls (CR, Bestand des Absenders) statt nur die Zusage
  des Tests; dazu „der neuen reviews-Schlüssel" in `configyaml_test.go` als Chronik-Wort. Die
  Vorprüfung in §8 des Plans (Chronik-Grep vor dem Commit) hat Text, der beim Umbruch mitwanderte,
  nicht erfasst.
- **verifizierbar:** nein (kein Gate prüft §3.7).
- **klasse:** `kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen`

### F-3 — MEDIUM — Der `SKIP_DIRS`-Negativtest fehlt, obwohl er an der Parallel-Achse gerade nachgezogen wurde

- **quelle:** Skill-Frage 13; Spezifikation Schritt 2 („die `SKIP_DIRS` ausgenommen"); Review R1 von
  `slice-263` F-5; Beobachtungs-Register `BEO-ALL/review-fix-applied-only-at-cited-site`
- **pfad:** `internal/hexagon/core/rules/reviews.go` · „`if !w.cfg.Recursive || isSkipDir(name) {`";
  `internal/hexagon/core/rules/reviews_muster_test.go`
- **befund:** Mit ausgeschaltetem `isSkipDir` bleibt die ganze Suite grün (Probe 4). `slice-263` R1 F-5
  meldete dieselbe Lücke für `planning`, dort ist sie seither durch
  `TestClosureRecursive_SkipDirsBleibenUnbetreten` gehalten; die Übertragung derselben Semantik auf
  `reviews` hat den Test nicht mitgenommen. Ein Abstieg in `node_modules/`, `vendor/` o. ä. unter
  `done-dir` liefe unbemerkt durch. Der Plan nennt die Register-Zeile mit Stand 2× — dies wäre der
  dritte Treffer.
- **verifizierbar:** ja — Test analog zum `planning`-Test, gegengeprüft mit der Mutation aus Probe 4.
- **klasse:** `review-fix-applied-only-at-cited-site`

### F-4 — MEDIUM — Die Byte-Identität ist mit einer Probe belegt, die `reviews` nicht einschaltet

- **quelle:** `AGENTS.md` §5 Regel 15; Skill-Fragen 8 und 17; Review R1 von `slice-263` F-6
- **pfad:** Commit `1837d8d2` · „make blackbox-probe REF=HEAD: byte-identisch ueber 20 Vergleiche."
- **befund:** `tools/blackbox-probe.sh` fährt das Image ohne Argumente über die Fixtures und das Repo;
  die Fixtures tragen keinen `reviews`-Block, und das Repo führt `reviews` nicht in `modules:`
  (`.d-check.yml` Zeile 32) — das geänderte Modul läuft in keinem der 20 Vergleiche. Die Botschaft stellt
  die Zahl als Beleg in den Kontext der `reviews`-Änderung; gemessen ist die Unverändertheit der
  **anderen** Module. Nach Lesen des Codes hält die Byte-Identität ohne die neuen Schlüssel (siehe
  Negativbefunde), getragen aber von den unveränderten Tests in `reviews_test.go`, nicht von der Probe.
  Dieselbe Klasse wie `slice-263` F-6. Die Botschaft ist eingefroren; die Korrektur gehört in die
  Closure-Notiz.
- **verifizierbar:** ja — `grep -rn reviews tools/blackbox-probe/fixtures` leer; `modules:`-Zeile.
- **klasse:** `botschaft-ueberdehnt-messung`

### F-5 — MEDIUM — Das Zusage-Muster trifft Checkbox-Punkte in Fences, in anderen Abschnitten und in Verneinung; getestet sind nur triviale Negativfälle

- **quelle:** Skill-Fragen 13 und 18; Lastenheft `DC-FA-RVW-001` Kriterium „Negative (Muster)" („…
  oder die Wortfolge außerhalb eines DoD-Punkts steht, then keine Zusage");
  `BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`
- **pfad:** `internal/hexagon/core/rules/reviews.go` · „`if !checkboxLineRE.MatchString(lines[i]) {`"
  in `reviewPromise`; `reviews_muster_test.go` · „`- [x] Adaptions-Review notiert`"
- **befund:** Ein „DoD-Punkt" ist im Code jede Checkbox-Zeile der Datei, gleich in welchem Abschnitt
  und auch innerhalb eines Fence; mit dem CR-Muster melden ein zitiertes DoD-Beispiel im Fence und ein
  Closure-Haken „kein Review durchgeführt (entfallen)" beide `review-missing` (Probe 2). Schritt 3
  „roh, zeilenweise" deckt das lexikalisch, die Grenzen-Liste der Spezifikation nennt es nicht, und
  das Lastenheft-Kriterium liest sich als Abschnitts-Aussage. Die Negativtests prüfen „Adaptions-Review
  notiert" und „Review-Report liegt vor" — beide können `Review durchgeführt` gar nicht treffen; kein
  Fall liegt nahe am Positiv. Ebenso unbenannt: das Muster läuft gegen den Punkt **samt** Bullet und
  Task-Box und ohne `(?m)`, ein Anker `^Review` trifft nie. Mit dem Default war die Reichweite durch die
  enge Phrase begrenzt; das opt-in-Muster ist gerade für breitere Formulierungen gedacht.
- **verifizierbar:** ja — Probe 2 als Test.
- **klasse:** `erkennungs-regex-nur-gegen-positivfaelle-getestet`

### F-6 — MEDIUM — `promise-pattern:` ohne Wert läuft still mit der Phrase

- **quelle:** Skill-Fragen 1 und 18; Spezifikation §2-Zeile `reviews.promise-pattern` („**Explizit**
  leer ⇒ Exit 2 (es träfe jeden Punkt)")
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go` · „`if r.PromisePattern != nil {`"
- **befund:** YAML-`null` (`promise-pattern:` ohne Wert, `~`) dekodiert in einen `nil`-Zeiger und
  gilt als abwesend; der Lauf nimmt dann still die Default-Phrase (Probe 5: `0 Befund(e)`), nur `''`
  bricht ab. Wer den Schlüssel hinschreibt und den Wert vergisst, bekommt genau den Zustand, den der CR
  meldet — grün über null Zusagen, solange `require-promises` fehlt. Die Spezifikation unterscheidet
  „explizit leer" von „abwesend", ordnet `null` aber keiner Seite zu. Dass `null` im Adapter allgemein
  als abwesend gilt, ist Bestand; an diesem Schlüssel fällt die Folge in die stille Richtung.
- **verifizierbar:** ja — Probe 5 als Decode-Test.
- **klasse:** `null-wert-faellt-still-auf-default`

### F-7 — LOW — „Dieselbe Semantik wie die Closure-Prüfung" trifft am Gegenstand nicht zu

- **quelle:** Skill-Frage 20; `AGENTS.md` §5 Regel 13
- **pfad:** `internal/hexagon/core/model/config.go` · „Dieselbe Semantik wie die
  Closure-Pruefung von planning."
- **befund:** Drei Unterschiede am Code: unlesbares Unterverzeichnis (`planning` ein Befund und
  Abbruch, `reviews` Teilmenge plus Leermeldung, F-1); unlesbare Datei (`planning` „bleibt Kandidatin",
  weil C3 sie meldet — in `reviews` bleibt sie Kandidatin und fällt danach still aus, zählt aber für
  die Nullmengen-Regel mit); nicht kompilierendes `skip-pattern` im Kern (`planning` nimmt nichts aus,
  `reviews` `regexp.MustCompile`). Die Grenze der unlesbaren Datei steht in der Spezifikation; der
  Kommentar verspricht trotzdem Gleichheit.
- **verifizierbar:** nein.
- **klasse:** `begruendung-trifft-am-gegenstand-nicht-zu`

### F-8 — LOW — Kopfkommentar von `rawReviews` zählt drei von acht Schlüsseln

- **quelle:** `MR-025` (Spiegel vor dem Editieren)
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go` · „done-dir (Aktivierungs-Schalter),
  reviews-dir und exempt-paths."
- **befund:** Der Kommentar über dem erweiterten Struct nennt die Parameter des Moduls als
  abgeschlossene Aufzählung und ist nicht nachgezogen; die Spiegel-Liste in §8 des Plans führt den
  Adapter, nicht diesen Kommentar.
- **verifizierbar:** nein (Doku-Drift).
- **klasse:** `spiegel-unvollstaendig`

### F-9 — LOW — Sensor-Datei friert Zahlen ein, statt das Kommando zu nennen

- **quelle:** `v6.17.0` · `regelwerk/modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin) · „mit dem
  **Kommando**, das den Ausschnitt zeigt, nicht mit einer eingefrorenen Zahl"
- **pfad:** `harness/sensors/review-coverage.md` · „26 (11 unter `done/`, 15 unter `done/wellenlos/`)"
- **befund:** Die Grenze „das Gate prüft heute keinen Slice" ist richtig gemessen (Probe 5), trägt
  aber Datum, Zählung und den Vorsatz „setzt die Schlüssel, sobald sie released sind"; mit der nächsten
  Closure stimmt die Zahl nicht mehr, und das Kommando, das sie reproduziert, steht nicht da.
- **verifizierbar:** nein.
- **klasse:** `grenze-mit-eingefrorener-zahl`

### F-10 — INFO — Das Risiko aus §6 ist eingetreten, nicht gebannt

- **quelle:** Slice-Plan §6 („die Zuordnung Report → Slice muss eindeutig bleiben"); Maintainability
- **pfad:** `internal/hexagon/core/rules/reviews.go` · „`strings.Contains(e.Name, key)`"
- **befund:** `match: name` deckt per Teilstring; das ist der Entwurf aus §8 und als Grenze in
  Lastenheft, Spezifikation und Code benannt. Die Präfix-Deckung reicht über das Beispiel
  `slice-a-foo`/`slice-a-foo-bar` hinaus: `slice-26` wird vom Report zu `slice-265-…` gedeckt (Probe 3).
  Bei der Closure braucht das Risiko einen Ausgang; eingetreten ist es.
- **verifizierbar:** ja — Probe 3.
- **klasse:** `teilstring-zuordnung-mehrdeutig`

### F-11 — INFO — Rückfall auf die Phrase heißt „fail-closed", ist es aber nicht

- **quelle:** Maintainability
- **pfad:** `internal/hexagon/core/rules/reviews.go` · „haelt aber den Lauf fail-closed statt leer"
- **befund:** Ein nicht kompilierendes Muster würde im Kern still gegen die Default-Phrase laufen —
  ohne Befund; `planning` meldet im selben Fall. Der Pfad ist vom Adapter aus unerreichbar, die
  Begründung im Kommentar trifft trotzdem nicht zu.
- **verifizierbar:** nein.
- **klasse:** `begruendung-trifft-am-gegenstand-nicht-zu`

## Negativbefunde

- **Hexagon-Richtung (`ADR-0005`):** `rules` importiert weiter nur `model` und `port/driven`; kein
  Adapter-Import. Ohne Befund.
- **Netz, Suppression:** kein Netzzugriff, kein `//nolint`, keine gesenkte Schwelle. Ohne Befund.
- **Byte-Identität ohne die neuen Schlüssel (Code gelesen):** Default-Muster ist wörtlich das alte
  `[Uu]nabhängiger Review`; ohne `recursive` geht ein Verzeichnis in `descend` und kehrt sofort zurück
  (alt: `Kind != KindFile` übersprungen), Symlinks bleiben übersprungen; `badDir` entsteht nur unter
  `recursive`, `skipRE` nur mit Muster, der `require-promises`-Zweig nur mit Schlüssel; Meldungstexte
  von `review-missing` und Leerlauf unverändert. Das vorgezogene `ReadFile` vor dem Kennungs-Check
  ändert nur, dass Dateien ohne Kennung gelesen werden; eine unlesbare Datei fällt wie zuvor still aus.
  Einzige Abweichung ist die gewollte: Zusage ohne `slice-<NNN>` meldet. Ohne Befund (Beleg-Lücke
  siehe F-4).
- **Exit-2-Liste (Regel 13, am Code gezählt):** Weißraum-`done-dir`, `done-dir` ohne `reviews-dir`,
  ungültiges Glob, explizit leeres oder nicht kompilierendes `promise-pattern`, `match` außer
  `id`/`name`/leer, nicht kompilierendes `skip-pattern` — sechs im Code, sechs in Lastenheft und
  Spezifikation, jede Meldung nennt den Schlüssel (Test `TestDecode_ReviewsMusterFehler`). Ohne Befund
  (bis auf `null`, F-6).
- **Grenzen-Liste der Spezifikation:** Präfix-Deckung, still ausfallende unlesbare Datei, zitierter
  Stub-Marker, nicht rekursives `reviews-dir`, Symlink-Verzeichnis benannt; fehlend F-1, F-5, F-6.
- **Referenz-Richtung (`AGENTS.md` §3.4):** Lastenheft- und Spezifikations-Diff ohne ADR-, Slice-,
  Wellen- oder Commit-Token; der Historien-Eintrag nennt den Anlass ohne Slice. Ohne Befund.
- **Tests aus dem richtigen Grund rot:** Kennungs-Befund (ohne den neuen `case id == ""` liefe der
  Kandidat leer durch, Test erwartet einen Befund mit „match: name"); Abnahme des CR (ohne Muster kein
  Befund, Test erwartet Zeile 3); `require-promises`, Abstieg, `skip-pattern`,
  Unterverzeichnis-Befund — je der zugehörige Test fällt bei Wegnahme, am Code nachvollzogen. Nicht
  gehalten: `SKIP_DIRS` (F-3) und der Abbruch der Restmenge (F-1).
- **`--print-config`-Vorlage:** die falsche Aussage „eine Zeile, die Review nennt" ist ersetzt; die
  fünf Schlüssel stehen mit Wirkung; das `skip-pattern`-Beispiel ist in YAML-Einzelanführung gültiges
  RE2. Ohne Befund.
- **Abgrenzung §1:** kein Handbuch-/README-Edit, keine Änderung an `.d-check.yml`, keine Modul-Doku im
  Image; die Sensor-Datei ist per Plan-Änderung vor dem Code aufgenommen. Ohne Befund.
- **Lese-Achse (`AGENTS.md` §3.8):** `reviews-dir` wird gelistet, nicht gescannt; die Grenze steht in
  Code und Spezifikation. Ohne Befund.

## Kategorie-Summary

HIGH 2 · MEDIUM 4 · LOW 3 · INFO 2. Wiederkehrend: `review-fix-applied-only-at-cited-site` (F-3, mit
F-4 und F-1 als Gegenstücke zu `slice-263` F-6/F-10 — drei Befunde dieses Reviews sind Klassen, die
R1 von `slice-263` an der Parallel-Achse schon gemeldet hat); `begruendung-trifft-am-gegenstand-nicht-zu`
(F-7, F-11).

## Verdikt

Nicht freigegeben. F-1 und F-2 blockieren; F-3 bis F-6 vor Merge klären. F-3 erreicht mit diesem Slice
den dritten Treffer seiner Register-Zeile — ein Steering-Loop-Signal, kein Einzelbefund.
