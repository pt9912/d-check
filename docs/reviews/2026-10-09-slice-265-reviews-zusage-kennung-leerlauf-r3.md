# Review R3: slice-265, `reviews` (Zusage-Muster, benannte Kennungen, Leerlauf und Unterverzeichnisse)

- **Review-Art:** Code. Geprüft wird der Fix-Commit gegen den Slice-Plan (`slice-265` §1, §2, §3 mit
  den Plan-Änderungen nach R1 und R2, §8), gegen die Befunde F-1 bis F-9 aus R2, gegen die
  Entscheidungen (`ADR-0081`, `MR-025`) und gegen die Hard Rules `AGENTS.md` §3.7/§3.8 sowie §5
  Regel 13/15. Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-265` · Fix-Commit `3632d568`, davor die Plan-Änderung `a387b494`.
- **Skill:** `reviewer.md` @ 1.19.0 (`980d6314`)
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** R2-Report `2026-10-09-slice-265-reviews-zusage-kennung-leerlauf-r2.md`;
  `DC-FA-RVW-001` im Lastenheft 0.101.1 (Beschreibung, „Warum diese beiden Formen als Default",
  Grenzen, Kriterien, Historie); Spezifikation `DC-FA-RVW-001.a` Schritte 3/4, Grenzen-Absatz,
  §2-Zeilen `reviews.promise-pattern`/`reviews.match`, `SPEC-081`, §8-Historie. Im Code gelesen:
  `internal/hexagon/core/rules/reviews.go` (vollständig), `reviews_muster_test.go`,
  `internal/hexagon/core/model/config.go`, die Exit-2-Meldung in `configyaml.go`,
  `config_template.go`, `harness/sensors/review-coverage.md` (vollständig), `Makefile`
  (`review-coverage`), `.d-check.yml` (Block `reviews`). Vendored Vorlage `v6.17.0` ·
  `templates/docs/plan/planning/slice.template.md` (DoD-Zeile).
- **Proben** (Wegwerf-Worktree auf `3632d568` im Scratchpad, `make test IMAGE=dcprobe-r3`; danach
  Worktree und Image entfernt, Arbeitsbaum unverändert):
  1. Probetest Default-Muster, je ein Slice ohne passenden Report (Zusage ⇔ ein Befund auf Zeile 3).
     **Zusage:** Vorlagen-Zeile mit Folgezeile; `*`, `+` (offen), `1.`; `[X]`; Tab zwischen Bullet,
     Box und Text; vier Leerzeichen nach der Box; eingerückter Punkt; CRLF-Zeilenende;
     `- [ ] Review durchgeführt?`; `- [x] Vorlage nennt `[ ] Review durchgeführt` als Punkt`;
     `- [x] Kein unabhängiger Review nötig`. **Keine Zusage:** „Review durchgeführt" auf der
     Folgezeile eines Punkts; `- [x] `make gates` grün; Review durchgeführt, Report unter
     `docs/reviews/`;` (wörtlich aus dem Bestand, siehe F-1); `Adaptions-Review durchgeführt`;
     `kein Review durchgeführt`; `**Review durchgeführt**`; NBSP nach der Box; kein Leerraum nach
     der Box.
  2. Probetest `match: name`: `slice-a-gro` / Report `…slice-a-gro` + U+0308 + `sse-r1.md`
     **gedeckt**; `slice-a-gr` / Report mit ungültigem Byte `\xff` bzw. abgeschnittenem `\xc3` danach
     **gedeckt**; `slice-a` / `…slice-a²-r1.md` **gedeckt**; `slice-a` / `…slice-a١-r1.md` (arabische
     Ziffer) und `…slice-aЖ-r1.md` nicht gedeckt; `slice-a` / Name endet mit dem Basisnamen gedeckt;
     `slice-x1` / `…slice-x12-y-r1.md` nicht gedeckt.
  3. Mutation: Default zurück auf `[Uu]nabhängiger Review|Review durchgeführt`, Wortgrenze zurück
     auf ASCII (`r >= 0x80 ||` vor der Unicode-Prüfung). Rot werden genau
     `TestReviewsDefault_VorlagenFormNurAmPunktanfang` und `TestReviewsMatch_NameWortgrenzeUnicode`;
     die übrigen Pakete bleiben `ok`.
  4. Bestand: `grep` auf Checkbox-Zeilen mit „Review durchgeführt" unter `docs/plan/planning/done/`
     (rekursiv): 25 × am Punktanfang, 1 × hinter anderem Text (`slice-256`, Report vorhanden).
     Spiegel-Grep nach „unabhängiger Review" / „Review durchgeführt" / `DefaultPromisePattern` über
     das Repo ohne `done/`, `docs/reviews/`, `.harness/`.

## Stand der R2-Befunde

| R2 | Stand | Beleg |
|---|---|---|
| F-1 (Komposita, Verneinung, naher Negativtest) | **geschlossen** für Komposita und Verneinung; der Fix bringt eine neue Lücke in der Gegenrichtung, siehe F-1 unten. „Review durchgeführt?" ist weiter Zusage, siehe F-4 | Probe 1, Probe 3; neuer Test mit allen drei Formen, mit dem alten Muster rot |
| F-2 (Plan verspricht Byte-Identität) | **in der Sache geschlossen** (§1 Ziel, DoD-Punkt 3 und §8 nennen die zwei Default-Änderungen); zwei der in R2 zitierten Stellen stehen unverändert, siehe F-3 | `a387b494` |
| F-3 (Spiegel des Defaults) | **geschlossen** für die vier gemeldeten Stellen; ein fünfter Spiegel fehlt, siehe F-2 | Exit-2-Meldung, Kopfkommentar, `Makefile`, Sensor-Vertrag |
| F-4 (Wortgrenze nur ASCII) | **geschlossen**; Restformen der Grenze siehe F-5 | Probe 2, Probe 3 |
| F-5 (Vorlage „leer ⇒ Exit 2") | **geschlossen** | „explizit leer ⇒ Exit 2, ohne Wert ⇒ Default" |
| F-6 (Slice-Nummer, „nur noch") | **geschlossen** | `slice-x1`/`slice-x12-y`; „MustCompile trifft hier gueltige Muster" |
| F-7 (Historie-Zählung) | **geschlossen** | Zeile 0.101.1; die Erweiterung „Unterverzeichnisse und Stubs" ist in der Zeile 0.101.0 inhaltlich genannt |
| F-8, F-9 (INFO) | unverändert, nicht blockierend | — |

## Findings

### F-1 (MEDIUM): Der Default verwirft eine Zusage-Form aus dem eigenen Bestand, und die Lastenheft-Begründung trifft am Bestand nicht zu

- **quelle:** `DC-FA-RVW-001` (Lastenheft 0.101.1, „Warum diese beiden Formen als Default" und
  Kriterium „Negative (Default)"); `AGENTS.md` §5 Regel 13; Skill-Fragen 1, 18 und 20;
  `BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`
- **pfad:** `spec/lastenheft.md` · „„`make gates` grün, Review durchgeführt"), oder die Wortfolge";
  `internal/hexagon/core/model/config.go` · „``\[[ xX]\][ \t]+Review durchgeführt` ``“;
  Bestand `slice-256` · „- [x] `make gates` grün; Review durchgeführt, Report unter `docs/reviews/`;“
- **befund:** Das Kriterium „Negative (Default)" erklärt „`make gates` grün, Review durchgeführt" zur
  Nicht-Zusage, und der neue Test hält genau diese Form als Negativfall. Im eigenen Bestand steht
  dieselbe Form als echte Zusage mit Report (`slice-256`, einer von 26 Checkbox-Punkten mit
  „Review durchgeführt"). Die Begründung im Lastenheft nennt nur „Adaptions-Review durchgeführt" und
  „kein Review durchgeführt", die tatsächlich keinen Report zusagen. Für die Wortfolge hinter anderem
  Text gibt sie keinen Grund, und am Gegenstand trifft das Gegenteil zu. Die Commit-Botschaft führt
  „die Wortfolge mitten im Satz" als Fehltreffer des alten Musters. Die Sensor-Datei sagt in Grenze 2,
  die DoD-Punkte dieses Repos trügen die Vorlagen-Form; das gilt für 25 von 26.
- **Failure-Szenario:** Ein Adopter schreibt die DoD-Zeile so wie `slice-256`: „- [x] `make gates` grün;
  Review durchgeführt, Report unter `docs/reviews/`". Der Report fehlt. Andere Slices tragen die
  Vorlagen-Form, also greift auch `require-promises` nicht. `make review-coverage` meldet für diesen
  Slice nichts, und das Gate ist an dieser Stelle still grün. Dieses Repo trifft das, sobald
  `reviews.recursive` gesetzt ist (§1, nach dem Release): `slice-256` liegt unter `done/wellenlos/`
  und fiele dann still aus der Prüfung. Heute fehlt dort kein Report.
- **verifizierbar:** ja, Probe 1 (`bestand-256`) bzw. ein Lauf mit `recursive: true` und einem
  `reviews-dir` ohne Reports: der Befundsatz zählt `slice-256` nicht.
- **klasse:** `erkennungs-regex-nur-gegen-positivfaelle-getestet` (Gegenrichtung: Negativfall nicht
  gegen den Bestand geprüft) · `begruendung-trifft-am-gegenstand-nicht-zu`

### F-2 (LOW): Ein fünfter Spiegel des Defaults ist nicht nachgezogen, und der Plan erklärt ihn zu keinem Spiegel

- **quelle:** `MR-025`; `AGENTS.md` §5 Regel 13; Slice-Plan §8 („`.d-check.yml` … nennen das Modul,
  ohne seine Erkennung zu beschreiben — kein Spiegel")
- **pfad:** `.d-check.yml` · „Review-Zusage (ein DoD-Item, dessen Text "unabhängiger Review" nennt --“
- **befund:** Der Kommentar über dem `reviews`-Block der eigenen Konfiguration beschreibt die Erkennung
  und nennt nur die Phrase. Dieser Block nutzt den Default (Sensor-Vertrag: „setzt kein eigenes Muster
  und nutzt den Default"), und der Default erkennt seit diesem Slice auch die Vorlagen-Form. Der
  Kommentar beschreibt damit nicht mehr, was da ist. Die Aussage in §8, `.d-check.yml` sei kein
  Spiegel, stimmt am Gegenstand nicht. Die Abgrenzung in §1 („Die Konfiguration dieses Repos")
  betrifft das Setzen der neuen Schlüssel, nicht diesen Kommentar.
- **verifizierbar:** ja, der Spiegel-Grep aus §8.
- **klasse:** `spiegel-unvollstaendig` · `review-fix-applied-only-at-cited-site`

### F-3 (LOW): Zwei der in R2 F-2 zitierten Plan-Stellen sind unverändert, die Plan-Änderung erklärt sie für nachgezogen

- **quelle:** `AGENTS.md` §6 Schritt 4, §5 Regel 15; R2 F-2
- **pfad:** Slice-Plan `slice-265` · „Auftraggeber-Entscheid 2026-10-09 (Default der
  Phrase bleibt, ein Release mit slice-263)“; §3 · „die Default-Phrase trifft keinen DoD-Punkt dieses Repos“
- **befund:** Der Kopf `**Bezug:**` nennt die erste Entscheidung („Default der Phrase bleibt") ohne
  Hinweis, dass sie mit der Plan-Änderung nach R1 revidiert wurde. Die §3-Zeile der Sensor-Datei sagt,
  der Default treffe keinen DoD-Punkt dieses Repos. Seit der Default-Erweiterung trifft er 25. Die
  Plan-Änderung nach R2 sagt: „Die Stellen im Plan, die „ohne die Schlüssel unverändert" sagten, sind
  nachgezogen." Das gilt für §1, die DoD und §8, nicht für diese zwei Stellen. Der DoD-Punkt selbst
  ist richtig gefasst, deshalb nur LOW.
- **verifizierbar:** nein (Plan-Text).
- **klasse:** `review-fix-applied-only-at-cited-site`

### F-4 (INFO): Die zwei Alternativen des Defaults sind ungleich verankert

- **quelle:** `DC-FA-RVW-001` („Warum diese beiden Formen als Default"); Maintainability
- **pfad:** `internal/hexagon/core/model/config.go` · „`` `[Uu]nabhängiger Review|` ``“
- **befund:** Die Vorlagen-Form gilt nur hinter der Task-Box, weil eine Verneinung keinen Report
  zusagt. Die Phrase hat weder Anker noch Verneinungs-Ausschluss: „- [x] Kein unabhängiger Review
  nötig" bleibt Zusage (Probe 1), ebenso „- [ ] Review durchgeführt?", der dritte Fall aus R2 F-1. Die
  Fehlerrichtung ist jeweils laut (ein unnötiges `review-missing`), nicht still. Bei der Phrase ist
  das Verhalten seit `ADR-0081` unverändert. Das Lastenheft begründet die Verneinung aber nur für eine
  der beiden Alternativen.
- **verifizierbar:** ja, Probe 1.
- **klasse:** `alternativen-ungleich-begrenzt`

### F-5 (INFO): Die Unicode-Wortgrenze lässt kombinierende Zeichen, ungültige Bytes und hochgestellte Ziffern als Grenze durch

- **quelle:** `DC-FA-RVW-001` („weder Buchstabe noch Ziffer (in jeder Schrift)", Grenze „über einen
  Bindestrich, Unterstrich oder Punkt"); Maintainability
- **pfad:** `internal/hexagon/core/rules/reviews.go` · „`r == utf8.RuneError ||`“
- **befund:** Ein kombinierendes Zeichen (U+0308, Kategorie Mn), ein ungültiges oder abgeschnittenes
  UTF-8-Byte und `²` (Kategorie No) gelten als Grenze. `slice-a-gro` wird deshalb vom Report zu
  `slice-a-gro` + U+0308 + `sse` gedeckt (NFD-Schreibweise von „grösse"), `slice-a` von `slice-a²-…`
  (Probe 2). Wörtlich trägt das der Vertragssatz, denn keines der drei ist Buchstabe oder Ziffer im
  Unicode-Sinn. Die Grenzen-Aufzählung „Bindestrich, Unterstrich oder Punkt" liest sich aber als
  geschlossen. Für reale Slug-Kennungen sind die Fälle sehr selten.
- **verifizierbar:** ja, Probe 2.
- **klasse:** `vertragssatz-weiter-als-code`

## Negativbefunde

- **Default-Muster, Positivfälle (Probe 1):** Vorlagen-Zeile mit Folgezeilen, alle drei Bullet-Formen
  und `1.`, `[X]`, Tab, mehrere Leerzeichen, Einrückung, CRLF: alle Zusage. Die Folgezeilen stören
  nicht, weil der Anker die Box der ersten Zeile ist. Ohne Befund.
- **Default-Muster, Negativfälle:** Komposita, Verneinung, Fettdruck, NBSP und eine fehlende Leerstelle
  nach der Box sind keine Zusage. Die letzten drei sind keine Vorlagen-Form, ein Ausschluss ist also
  vertretbar. Eine inline zitierte Box mitten im Punkt trifft (Probe 1), das ist harmlos. Ohne Befund
  bis auf F-1 und F-4.
- **Tests aus dem richtigen Grund rot (Probe 3):** Beide neuen Tests werden mit dem jeweiligen
  Vorzustand rot und sonst kein Test. Die Botschaft behauptet das für beide („mit dem alten Muster
  rot", „mit ASCII-Grenze rot"), und es stimmt. Ohne Befund.
- **`hasReviewContaining`:** Das Namensende (`DecodeRuneInString("")` ⇒ `RuneError`) deckt, wie
  vorher. `rest = rest[i+1:]` kann mitten in einer Mehrbyte-Sequenz landen. Das ist unschädlich, weil
  `strings.Index` byteweise sucht und der Basisname gültiges UTF-8 ist. Andere Schriften (Kyrillisch,
  arabische Ziffern) setzen den Namen fort (Probe 2). Ohne Befund bis auf F-5.
- **Spiegel (Probe 4):** Exit-2-Meldung, Kopfkommentar der Regel, `Makefile`-Hilfe, Sensor-Vertrag,
  `--print-config`-Vorlage, Spezifikation (Schritt 3, §2, `SPEC-081`) und Lastenheft nennen den
  Default mit Task-Box. README, README.de und Handbuch sind Release-Prep (§1, Regel 17).
  `ADR-0081` ist immutabel, die CR-Dateien sind Belege, `reviews_test.go` verwendet die Phrase als
  Testtext. Ohne Befund bis auf F-2.
- **Lastenheft 0.101.1 und Spezifikation gegen Code:** Die Task-Box-Form `\[[ xX]\][ \t]+` steht in
  der Spezifikation wörtlich wie im Code. „In jeder Schrift" entspricht `unicode.IsLetter`/`IsDigit`.
  Die Grenze nennt Unterstrich und Punkt. Der Versions-Bump 0.101.0 → 0.101.1 ist ein Nachzug ohne
  neue Anforderung. Ohne Befund bis auf F-1 und F-5.
- **Referenz-Richtung (`AGENTS.md` §3.4):** Der Spec-Diff enthält kein Slice-, Wellen-, ADR- oder
  Commit-Token. Ohne Befund.
- **Kommentare (`AGENTS.md` §3.7):** `config.go` (Zusage plus Abgrenzung), `reviews.go` (Kopplung,
  GRENZE mit neutralen Beispielen, ohne „nur noch"), Testkommentare und die Vorlagen-Zeilen tragen eine
  der fünf Klassen, ohne Slice-Nummer und ohne Chronik. Ohne Befund.
- **Hexagon, Netz, Suppression:** Neu sind nur die Standardbibliotheks-Importe `unicode` und
  `unicode/utf8` im Kern. Kein Netz, kein `//nolint`, keine gesenkte Schwelle. Ohne Befund.
- **Lese-Achse (`AGENTS.md` §3.8):** Die Achsen sind unverändert und bleiben benannt. Ohne Befund.
- **Abgrenzung §1:** Handbuch und README sind nicht angefasst, die `reviews`-Schlüssel der
  `.d-check.yml` auch nicht. Das `Makefile` ist vor dem Code in §3 aufgenommen (`a387b494`). Ohne
  Befund.
- **Commit-Botschaft (Regel 15):** Die Zuordnung F-1 bis F-7 stimmt mit dem Diff überein, die
  Mutationsbehauptungen ebenfalls (Probe 3). Überdehnt sind „die Wortfolge mitten im Satz" als
  Fehltreffer (F-1) und „vier Spiegel" als vollständige Menge (F-2). „make review-coverage … 0 Befunde"
  ist nicht nachgemessen; es ist eine Lauf-Bestätigung (Verifier).

## Kategorie-Summary

HIGH 0 · MEDIUM 1 · LOW 2 · INFO 2. Wiederkehrend: `review-fix-applied-only-at-cited-site` (F-2, F-3)
zum dritten Mal in diesem Slice, nach R2 F-1 bis F-3. `erkennungs-regex-nur-gegen-positivfaelle-getestet`
taucht in F-1 in der Gegenrichtung auf: Der Negativfall wurde nicht gegen den Bestand geprüft. Beide
Register-Zeilen stehen nach R2 an der Schwelle. Das ist ein Steering-Loop-Signal für die Closure.

## Verdikt

**Nicht freigegeben.** R2 F-1 bis F-7 sind in der Sache geschlossen. Die zwei neuen Tests sind aus
dem richtigen Grund rot, und Spezifikation und Code stimmen überein. Vor dem Merge zu klären ist F-1:
Der Default und sein Lastenheft-Kriterium erklären eine Form zur Nicht-Zusage, die im eigenen Bestand
eine echte Zusage mit Report ist. Ohne diesen Report wäre das Gate an dieser Stelle still grün, und die
Begründung im Lastenheft deckt diesen Ausschluss nicht. Ob das Muster geändert wird oder der
Ausschluss als benannte Grenze mit zutreffender Begründung bleibt, entscheidet die Implementation. F-2
bis F-5 blockieren nicht.
