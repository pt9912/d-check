# Review R4: slice-265, `reviews` (Zusage-Muster, benannte Kennungen, Leerlauf und Unterverzeichnisse)

- **Review-Art:** Code. Geprüft wird der Fix-Commit gegen den Slice-Plan (`slice-265` mit den
  Plan-Änderungen nach R1 bis R3), gegen die Befunde F-1 bis F-3 aus R3, gegen die Entscheidungen
  (`ADR-0081`, `MR-025`) und gegen die Hard Rules `AGENTS.md` §3.4/§3.7 sowie §5 Regel 13/15.
  Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-265` · Fix-Commit `54edbf59`, davor die Plan-Änderung `a0da941a`.
- **Skill:** `reviewer.md` @ 1.19.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** R3-Report `2026-10-09-slice-265-reviews-zusage-kennung-leerlauf-r3.md`;
  `DC-FA-RVW-001` im Lastenheft 0.101.2 (Beschreibung, „Warum diese beiden Formen als Default",
  Kriterien „Happy Path (Vorlagen-Form)" und „Negative (Default)", Historie); Spezifikation
  `DC-FA-RVW-001.a` Schritt 3, §2-Zeile `reviews.promise-pattern`, `SPEC-081`, §8-Historie. Im Code
  gelesen: `model.DefaultPromisePattern` samt Kommentar, `reviewPromise` in `reviews.go` (Item =
  Checkbox-Zeile plus Folgezeilen, mit `\n` verbunden), `reviews_muster_test.go`, Exit-2-Meldung in
  `configyaml.go`, `config_template.go`, `Makefile` (`review-coverage`), `.d-check.yml`
  (`reviews`-Block), `harness/sensors/review-coverage.md`.
- **Proben** (Wegwerf-Worktree auf `54edbf59` im Scratchpad, Image `dcr4` per `make build`/`make test`;
  danach Worktree, Kopien und Image entfernt, Arbeitsbaum unverändert):
  1. **Bestand.** Kopie von `docs/plan/planning/done/` (rekursiv), Lauf mit `recursive: true`,
     `skip-pattern: '(?m)^> \*\*ARCHIVIERT\*\*'` und einem leeren `reviews-dir`. 31 Volltexte bleiben
     Kandidaten. `grep` findet in ihnen 26 Zeilen mit „Review durchgeführt", alle Checkbox-Punkte. Der
     Lauf meldet genau diese 26 als `review-missing`, je auf derselben Zeile, `slice-256`
     (`` `make gates` grün; Review durchgeführt `` …) eingeschlossen. Zusätzlich die 149 archivierten
     Volltexte aus den 91 `archiv.zip` per `unzip` gelesen: 7 Treffer, alle
     „- [x] Unabhängiger Review durchgeführt", also über die Phrase gedeckt. Keine Fundstelle mit „Review
     durchgeführt" hinter `.`, `:`, `—` oder `(`, keine am Anfang einer Folgezeile.
  2. **Default gegen Einzelfälle** (Image-Lauf, je ein Slice ohne Report). **Zusage:** Vorlagen-Zeile;
     `; Review durchgeführt`; `, Review durchgeführt`; `,Review durchgeführt` ohne Leerraum; `;` + Tab;
     `*`- und `1.`-Bullet, `[X]`; Folgezeile nach `;` (LF und CRLF); `- [ ] Review durchgeführt?`;
     `- [x] Gates grün, Review durchgeführt?`; `- [x] Gates grün, Review durchgeführten …`;
     **`- [ ] Risiko: bei der Closure wird kein` / `      Review durchgeführt`**;
     **`- [x] Gates grün, kein` / `  Review durchgeführt`**. **Keine Zusage:** `kein Review durchgeführt`;
     `Adaptions-Review durchgeführt`; `das Review durchgeführt zu haben`; `Ist das Review durchgeführt?`;
     `, nicht Review durchgeführt`; `und Review durchgeführt`; Fließtext ohne Checkbox (auch mit `;`
     davor); Zeile nach einer Leerzeile; Unterpunkt ohne Box; `**Review durchgeführt**`;
     **`Gates grün. Review durchgeführt`**, **`Abschluss: Review durchgeführt`**,
     **`Gates grün — Review durchgeführt`**, **`(Review durchgeführt)`**.
  3. **Mutation.** `make test` auf `54edbf59` grün. Mit dem Default ohne die Alternative `(?m:^)` wird
     genau `TestReviewsDefault_VorlagenFormNurAmTeilanfang` rot („drei Zusagen ohne Report erwartet",
     nur zwei geliefert). Die übrigen Pakete bleiben `ok`.
  4. **Spiegel-Grep** nach „Task-Box", „Punkt-Teil", „Punktanfang", „unabhängiger Review",
     „Review durchgeführt", `DefaultPromisePattern` über das Repo ohne `done/`, `docs/reviews/`,
     `.harness/`, `docs/plan/cr/`, `docs/plan/adr/`.

## Stand der R3-Befunde

| R3 | Stand | Beleg |
|---|---|---|
| F-1 (Default verwirft die `slice-256`-Form) | **geschlossen**. Jeder der 26 Checkbox-Punkte mit „Review durchgeführt" im Bestand ist Zusage, die 7 archivierten tragen die Phrase. Die Lastenheft-Begründung schließt die Form nicht mehr aus. Der Fix bringt eine neue Lücke in der Gegenrichtung, siehe F-1 unten | Probe 1, Probe 3 |
| F-2 (`.d-check.yml` als Spiegel) | **geschlossen**. Der Kommentar nennt das Default-Muster, §3 und §8 des Plans führen die Datei als Spiegel | `.d-check.yml` Z. 991–996, `a0da941a` |
| F-3 (zwei Plan-Stellen) | **geschlossen**. `**Bezug:**` nennt beide Entscheide, die §3-Zeile der Sensor-Datei behauptet keinen Messstand mehr | `a0da941a` |
| F-4, F-5 (INFO) | unverändert, nicht blockierend. Die Frage-Form reicht jetzt weiter, siehe F-3 unten | Probe 2 |

## Findings

### F-1 (MEDIUM): Die Alternative `(?m:^)` macht eine umbrochene Verneinung zur Zusage, gegen das eigene Negativ-Kriterium

- **quelle:** `DC-FA-RVW-001` (Lastenheft 0.101.2, Kriterium „Negative (Default)"); Spezifikation
  `DC-FA-RVW-001.a` Schritt 3; Skill-Fragen 10 und 11; `AGENTS.md` §5 Regel 13;
  `BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`
- **pfad:** `internal/hexagon/core/model/config.go` · „`` |(?m:^))[ \t]*Review durchgeführt` ``";
  `spec/lastenheft.md` · „die Vorlagen-Form hinter einem Wort („Adaptions-Review durchgeführt", „kein Review durchgeführt""
- **befund:** Das Kriterium „Negative (Default)" und Spezifikation Schritt 3 erklären „kein Review
  durchgeführt", also die Form hinter einem Wort, zur Nicht-Zusage. Steht der Zeilenumbruch eines
  Punkts zwischen dem Wort und „Review", ist sie Zusage. „…, kein" / „Review durchgeführt" und
  „… wird kein" / „Review durchgeführt" melden `review-missing` (Probe 2). Ob derselbe Text Zusage
  ist, hängt damit am Umbruch: in einer Zeile ist „und Review durchgeführt" keine Zusage, umbrochen nach
  „und" ist es eine. Für die Alternative gibt es im Bestand keinen Beleg (Probe 1: keine Fundstelle am
  Anfang einer Folgezeile). Der einzige Positivfall im Test („Gates grün;" / „Review durchgeführt")
  steht hinter `;`. Der neue Test enthält keinen umbrochenen Negativfall. Die Commit-Botschaft sagt
  „Verneinung … bleib[t] ausgeschlossen" und gibt das ohne diese Einschränkung wieder.
- **Failure-Szenario:** Ein Slice-Plan in diesem Repo (harter Umbruch bei rund 80 Zeichen) trägt im
  Risiko-Abschnitt „- [ ] Risiko: bei der Closure wird kein" / „Review durchgeführt, wenn …". Weil
  jeder Checkbox-Punkt der Datei zählt, meldet `make review-coverage` nach der Closure `review-missing`.
  Laut Lastenheft und Spezifikation sagt dieser Punkt keinen Report zu. Abhilfe gibt es nur durch einen
  anderen Umbruch oder ein eigenes Muster. Der Fehler ist laut, nicht still. Heute steht keine solche
  Stelle im Bestand.
- **verifizierbar:** ja, Probe 2 (Image-Lauf auf zwei Ein-Punkt-Slices ohne Report).
- **klasse:** `erkennungs-regex-nur-gegen-positivfaelle-getestet` (Negativfall nur in einem Layout
  geprüft) · `semantik-haengt-am-zeilenumbruch`

### F-2 (INFO): Die geschlossene Liste der Teil-Anfänge ist nur für „hinter einem Wort" begründet

- **quelle:** `DC-FA-RVW-001` („Warum diese beiden Formen als Default", Kriterium „Negative
  (Default)"); Skill-Frage 20
- **pfad:** `spec/lastenheft.md` · „eröffnet (hinter der Task-Box, hinter `;` oder `,`, am Zeilenanfang)"
- **befund:** Hinter `.`, `:`, `—` oder `(` ist „Review durchgeführt" keine Zusage (Probe 2). Für einen
  Leser sind „Gates grün. Review durchgeführt" oder „Abschluss: Review durchgeführt" ebenso der Anfang
  eines Punkt-Teils wie die Form hinter `;`. Die Aufzählung ist geschlossen und damit genau. Die
  Begründung nennt aber nur Verneinung und Kompositum, und das Negativ-Kriterium führt diese Formen
  nicht. Im Bestand steht keine davon, auch nicht in den 149 archivierten Volltexten (Probe 1). Ein
  Adopter mit dieser Schreibweise fällt still heraus, solange `require-promises` nicht gesetzt ist.
  Das deckt Grenze 2 der Sensor-Datei sinngemäß ab.
- **verifizierbar:** ja, Probe 2.
- **klasse:** `ausschluss-ohne-begruendung`

### F-3 (INFO): Die Frage-Form ist an mehr Stellen Zusage als vorher

- **quelle:** R3 F-4; Maintainability
- **pfad:** `internal/hexagon/core/model/config.go` · „`` (?:\[[ xX]\]|[;,]|(?m:^)) ``"
- **befund:** „- [x] Gates grün, Review durchgeführt?" ist jetzt Zusage, unter dem R3-Muster war sie
  es nicht. Ebenso eine Fortsetzung wie „Review durchgeführten …", weil das Muster rechts offen ist.
  Beides geht in die laute Richtung, und keine dieser Formen steht im Bestand. „Ist das Review
  durchgeführt?" bleibt keine Zusage.
- **verifizierbar:** ja, Probe 2.
- **klasse:** `alternativen-ungleich-begrenzt`

## Negativbefunde

- **Bestand (Probe 1):** Alle 26 Checkbox-Punkte mit „Review durchgeführt" unter `done/` (rekursiv,
  ohne Stubs) sind Zusage, auf der richtigen Zeile, `slice-256` eingeschlossen. Ohne Befund.
- **Negativfälle auf einer Zeile (Probe 2):** Verneinung, Kompositum, Form hinter einem Wort,
  Fettdruck, Fließtext ohne Checkbox (auch mit `;` davor), Zeile nach Leerzeile und Unterpunkt ohne
  Box sind keine Zusage. Ohne Befund bis auf F-1.
- **Test aus dem richtigen Grund rot (Probe 3):** Ohne `(?m:^)` fällt genau der neue Test, mit der
  erwarteten Meldung. Die Botschaft „mit dem engeren Muster rot" trifft zu, denn jede engere Variante
  verliert einen der drei Positivfälle. Ohne Befund.
- **Spiegel (Probe 4):** `model`-Kommentar, Exit-2-Meldung, `--print-config`-Vorlage,
  `Makefile`-Hilfe, `.d-check.yml`-Kommentar, Spezifikation (Schritt 3, §2, `SPEC-081`, §8) und
  Lastenheft (Beschreibung, Begründung, beide Kriterien, Historie 0.101.2) nennen dieselben drei
  Teil-Anfänge. Die Sensor-Datei beschreibt das Muster nicht, und ihre Grenze 2 stimmt jetzt am Bestand
  (26 von 26). README, README.de, CHANGELOG und Handbuch sind Release-Prep (Regel 17). `slice-265` §3
  Z. 133 ist die Notiz zur Plan-Änderung nach R2, also Historie des Plans. Ohne Befund.
- **Lastenheft 0.101.2 und Spezifikation gegen Code:** „am Anfang einer Folgezeile" und
  „am Zeilenanfang" meinen dasselbe, weil die erste Zeile mit dem Bullet beginnt. `[ \t]*` deckt die
  Einrückung der Folgezeile, `(?m:^)` greift auch bei CRLF (Probe 2). Der Patch-Bump unter `Draft`
  folgt `MR-032` und der Vorgängerzeile 0.101.1. Ohne Befund bis auf F-1 und F-2.
- **Referenz-Richtung (`AGENTS.md` §3.4):** Der Spec-Diff enthält kein Slice-, Wellen-, ADR- oder
  Commit-Token. Ohne Befund.
- **Kommentare (`AGENTS.md` §3.7):** Der Kommentar über `DefaultPromisePattern` trägt Zusage und
  Abgrenzung, der Test-Kommentar eine Zusage, der `.d-check.yml`-Kommentar eine Zusage, die
  Vorlagen-Zeile eine Zusage mit Abgrenzung. Nirgends stehen Befund-Nummern, Slice-Nummern oder
  Chronik. Ohne Befund.
- **Hexagon, Netz, Suppression:** Keine neuen Importe, kein Netz, kein `//nolint`, keine gesenkte
  Schwelle. Ohne Befund.
- **Plan-Treue:** Die Plan-Änderung nach R3 (`a0da941a`) steht vor dem Code-Commit und nennt die
  drei Teil-Anfänge, den neuen Spiegel und die zwei nachgezogenen Stellen. Der Code tut genau das.
  Handbuch und README sind nicht angefasst. Ohne Befund.
- **Commit-Botschaft (Regel 15):** Die Zuordnung F-1 bis F-3 und die Spiegel-Liste stimmen mit dem
  Diff überein. Überdehnt ist „Verneinung … bleib[t] ausgeschlossen" (F-1). „make test, make lint
  grün; make review-coverage 0 Befunde" ist eine Lauf-Bestätigung (Verifier); `make test` ist hier
  nachgefahren und grün.

## Kategorie-Summary

HIGH 0 · MEDIUM 1 · LOW 0 · INFO 2. Wiederkehrend: `erkennungs-regex-nur-gegen-positivfaelle-getestet`
taucht zum dritten Mal in diesem Slice auf (R2 F-1, R3 F-1, R4 F-1), jedes Mal als Negativfall, der
nur in einer Form geprüft wurde. Die Register-Zeile steht nach R2 an der Schwelle. Das ist ein
Steering-Loop-Signal für die Closure. `review-fix-applied-only-at-cited-site` ist in dieser Runde nicht
aufgetreten.

## Verdikt

**Nicht freigegeben.** R3 F-1 bis F-3 sind geschlossen. Das Default-Muster trifft jeden
Checkbox-Punkt mit „Review durchgeführt" im eigenen Bestand, aktiv wie archiviert, und schließt
Verneinung, Kompositum und die Form hinter einem Wort aus, solange sie auf einer Zeile stehen. Vor dem
Merge zu klären ist F-1. Die Alternative `(?m:^)` hat im Bestand keinen Beleg. Sie macht eine
umbrochene Verneinung zur Zusage, und das widerspricht dem Negativ-Kriterium, das derselbe Commit
schreibt. Ob das Muster geändert wird oder die Grenze in Lastenheft und Spezifikation steht,
entscheidet die Implementation. F-2 und F-3 blockieren nicht.
