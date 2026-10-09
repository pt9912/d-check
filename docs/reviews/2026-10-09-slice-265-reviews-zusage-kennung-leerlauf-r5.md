# Review R5: slice-265, `reviews` (Zusage-Muster nach R4)

- **Review-Art:** Code. Geprüft wird der Fix-Commit gegen den Slice-Plan (`slice-265` samt
  Plan-Änderung nach R4), gegen den R4-Befund F-1, gegen die Entscheidungen (`ADR-0081`, `MR-025`)
  und gegen die Hard Rules `AGENTS.md` §3.4/§3.7 sowie §5 Regel 13/15. Die DoD-Abhakung prüft dieses
  Review nicht.
- **Gegenstand:** `slice-265` · Fix-Commit `6b861426`, davor die Plan-Änderung `d5ea3203`.
- **Skill:** `reviewer.md` @ 1.19.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** R4-Report `2026-10-09-slice-265-reviews-zusage-kennung-leerlauf-r4.md`;
  `DC-FA-RVW-001` im Lastenheft 0.101.3 (Beschreibung, „Warum diese beiden Formen als Default",
  Kriterien „Happy Path (Vorlagen-Form)" und „Negative (Default)", Historie); Spezifikation
  `DC-FA-RVW-001.a` Schritt 3, §2-Zeile `reviews.promise-pattern`, `SPEC-081`, §8-Historie. Im Code:
  `model.DefaultPromisePattern` samt Kommentar, `reviewPromise` in `reviews.go` (Item = Checkbox-Zeile
  plus Folgezeilen, mit `\n` verbunden), `reviews_muster_test.go`, `config_template.go`, Exit-2-Meldung
  in `configyaml.go`, `Makefile` (`review-coverage`), `.d-check.yml` (`reviews`-Block).
- **Proben** (Wegwerf-Worktree auf `6b861426` im Scratchpad, Image `dcr5` per `make build`; danach
  Worktree, Kopien und Image entfernt, Arbeitsbaum unverändert):
  1. **Bestand.** Kopie von `docs/plan/planning/done/` (rekursiv), Lauf mit `recursive: true`,
     `skip-pattern: '(?m)^> \*\*ARCHIVIERT\*\*'` und leerem `reviews-dir`. 31 Volltexte sind
     Kandidaten, `grep` findet darin 26 Zeilen mit „Review durchgeführt", alle Checkbox-Punkte. Der
     Lauf meldet 26 `review-missing`; die Menge Datei:Zeile ist mit der `grep`-Liste identisch
     (`diff` leer), `slice-256` Z. 103 eingeschlossen. Zusätzlich die Volltexte aus den 171
     `archiv.zip` per `unzip -p` gelesen: Jede Zeile, die mit „Review durchgeführt" beginnt, steht in
     Fließtext (§5 Closure-Trigger), in keinem Checkbox-Punkt.
  2. **Default gegen Einzelfälle** (Image-Lauf, 24 Ein-Punkt-Slices ohne Report). **Zusage (11/11
     gemeldet):** Vorlagen-Zeile; `; Review durchgeführt`; `, Review durchgeführt`; `;` + Umbruch;
     `,` + Umbruch; `;` + CRLF; `, ` + Leerzeichen + Umbruch; `;` + Tab + Umbruch + Tab; Task-Box allein
     auf der Zeile + Umbruch; `* [X]`; Phrase. **Keine Zusage (13/13 still):** `kein Review
     durchgeführt`; `… wird kein` / `Review durchgeführt`; `…, kein` / `Review durchgeführt`;
     `… und` / `Review durchgeführt`; `…, nicht` / `Review durchgeführt`; `Adaptions-Review
     durchgeführt`; `; Adaptions-` / `Review durchgeführt`; Folgezeile ohne Satzzeichen davor
     (`` `make gates` grün `` / `Review durchgeführt`); Leerzeile dazwischen; Bullet ohne Box;
     Unterpunkt ohne Box; `das Review durchgeführt zu haben`; Phrase über einen Umbruch
     (`; unabhängiger` / `Review (Report)`).
  3. **Mutation.** `make test` auf `6b861426` grün (Exit 0). Mit der R4-Alternative `(?m:^)` zurück
     fällt genau `TestReviewsDefault_VorlagenFormNurAmTeilanfang` auf `slice-007-g.md` („keine der
     vier Formen ist eine Zusage"). Mit `[ \t]*` statt `[ \t\r\n]*` fällt derselbe Test im
     Positivteil (zwei statt drei Zusagen, `slice-003-c.md` fehlt). Beide Erweiterungen des Fixes
     sind damit aus dem richtigen Grund gehalten.
  4. **Spiegel-Grep** nach „Zeilenanfang", „Folgezeile", „Punkt-Teil", „Teil des Punkts",
     „Teils eines Punkts" über das Repo ohne `done/`, `docs/reviews/`, `.harness/`, `docs/plan/cr/`,
     `docs/plan/adr/`.

## Stand der R4-Befunde

| R4 | Stand | Beleg |
|---|---|---|
| F-1 (umbrochene Verneinung wird Zusage) | **geschlossen**. Die Alternative `(?m:^)` ist entfallen, `[ \t\r\n]*` hinter Box/`;`/`,` trägt den Umbruch. Alle fünf umbrochenen Negativformen bleiben still, alle Umbruch-Positivfälle hinter `;`/`,` bleiben Zusage. Der Test hält die umbrochene Verneinung und fällt mit der alten Alternative | Probe 2, Probe 3 |
| F-2 (INFO, geschlossene Liste nur für „hinter einem Wort" begründet) | unverändert, nicht blockierend. Die Liste ist um einen Fall enger geworden: eine Folgezeile, die ohne Satzzeichen davor mit „Review durchgeführt" beginnt, ist keine Zusage mehr. Das steht so in Lastenheft und Spezifikation; im Bestand gibt es keinen solchen Checkbox-Punkt | Probe 1, Probe 2 |
| F-3 (INFO, Frage-Form) | unverändert, nicht blockierend | — |

## Findings

### F-1 (INFO): Die Phrase gilt nicht über einen Zeilenumbruch, die Vorlagen-Form jetzt schon

- **quelle:** `DC-FA-RVW-001` (Beschreibung: „die Phrase „unabhängiger Review""); Skill-Frage 11
  sinngemäß innerhalb eines Musters; Maintainability
- **pfad:** `internal/hexagon/core/model/config.go` · „`` [Uu]nabhängiger Review|(?:\[[ xX]\]|[;,])[ \t\r\n]*Review durchgeführt ``"
- **befund:** Die erste Alternative verlangt ein wörtliches Leerzeichen. Ein Checkbox-Punkt
  „…; unabhängiger" / „Review (Report …)" ist deshalb keine Zusage (Probe 2), obwohl er in Markdown
  als „unabhängiger Review" gelesen wird. Die Vorlagen-Form toleriert den Umbruch seit diesem Fix.
  Im archivierten Bestand stehen vier solche Punkte (die DoD-Punkte von `slice-129`, `slice-131`,
  `slice-148`, `slice-150`, je „`make gates` Exit 0 …; unabhängiger" / „Review ([Report] …)"); unter
  den 31 lebenden Volltexten keiner. Die Lücke bestand schon vor diesem Slice, der Fix führt sie
  nicht ein. Sie ist still: Ein Adopter mit hartem Umbruch bei rund 80 Zeichen, dessen Punkt so
  umbricht und dem der Report fehlt, bekommt keinen Befund, solange `require-promises` nicht gesetzt
  ist. Weder die Grenze der Sensor-Datei noch das Negativ-Kriterium nennt diesen Fall.
- **verifizierbar:** ja, Probe 2 (Image-Lauf, Ein-Punkt-Slice ohne Report).
- **klasse:** `alternativen-ungleich-begrenzt`

## Negativbefunde

- **R4 F-1 (Probe 2, 3):** Umbrochene Verneinung („kein", „nicht", „und", Kompositum) ist keine
  Zusage mehr; die Mutation zeigt, dass der neue Testfall genau daran hängt. Ohne Befund.
- **Bestand (Probe 1):** 26 von 26 Checkbox-Punkten mit „Review durchgeführt" unter `done/`
  (rekursiv, ohne Stubs) bleiben Zusage, jeweils auf der richtigen Zeile. Der Wegfall von `(?m:^)`
  kostet im Bestand keinen Punkt, auch nicht in den Archiven. Ohne Befund.
- **Neue Fehler durch `[ \t\r\n]*` (Probe 2):** Der Umbruch wird nur hinter Box, `;` und `,`
  übersprungen; eine Leerzeile beendet das Item vorher (`reviewPromise`), ein Unterpunkt ohne Box
  bleibt still, CRLF und Tab sind gedeckt. Kein neuer Positiv- oder Negativfehler gefunden. Ohne
  Befund bis auf F-1 (bestehend, nicht neu).
- **Lastenheft 0.101.3 und Spezifikation gegen Code:** Beschreibung, Begründung, „Happy Path
  (Vorlagen-Form)" („auch mit Zeilenumbruch nach dem `;`") und „Negative (Default)" („„kein Review
  durchgeführt" auch über einen Zeilenumbruch") entsprechen dem Muster und den Tests. Spezifikation
  Schritt 3 und §2-Zeile ebenso. `SPEC-081`, die Exit-2-Meldung, die `Makefile`-Hilfe und der
  `.d-check.yml`-Kommentar sagen nur „am Anfang eines Punkt-Teils" ohne Aufzählung und bleiben
  damit richtig. Der Patch-Bump unter `Draft` folgt `MR-032`. Ohne Befund.
- **Spiegel (Probe 4):** Ein Treffer auf „am Zeilenanfang" oder „Folgezeile" in Bezug auf die
  Vorlagen-Form steht außerhalb der eingefrorenen Verzeichnisse nur noch in `slice-265` §3 (Notizen
  zu den Plan-Änderungen nach R3 und R4, also Plan-Historie) und in den Historie-Zeilen 0.101.2 und
  der Spezifikations-Historie (Chronik, gewollt). Die übrigen Treffer betreffen andere Module.
  README, CHANGELOG und Handbuch sind Release-Prep (Regel 17). Ohne Befund.
- **Referenz-Richtung (`AGENTS.md` §3.4):** Der Spec-Diff enthält kein Slice-, Wellen-, ADR- oder
  Commit-Token. Ohne Befund.
- **Kommentare (`AGENTS.md` §3.7):** Der Kommentar über `DefaultPromisePattern` trägt Zusage und
  Abgrenzung, der Test-Kommentar eine Zusage, die Vorlagen-Zeile in `config_template.go` Zusage und
  Abgrenzung. Keine Befund-Nummer, keine Slice-Nummer, keine Chronik. Ohne Befund.
- **Hexagon, Netz, Suppression:** Keine neuen Importe, kein Netz, kein `//nolint`, keine gesenkte
  Schwelle. Ohne Befund.
- **Plan-Treue:** Die Plan-Änderung nach R4 (`d5ea3203`) steht vor dem Code-Commit, nennt den
  Wegfall der Alternative, die Umbruch-Toleranz hinter `;`/`,` und den neuen Test. Der Code tut genau
  das und nicht mehr. Ohne Befund.
- **Commit-Botschaft (Regel 15):** „hatte im Bestand keinen Beleg" stimmt (Probe 1, auch Archive).
  „hinter der Task-Box oder hinter ; oder , gilt die Form auch über einen Zeilenumbruch" stimmt
  (Probe 2). „mit der alten Alternative rot" stimmt (Probe 3). Die Spiegel-Liste stimmt mit dem
  Diff überein und lässt nur Stellen aus, deren Wortlaut ohnehin richtig bleibt. „make test, make
  lint grün; make review-coverage 0 Befunde" ist eine Lauf-Bestätigung (Verifier); `make test` ist
  hier nachgefahren und grün. Ohne Befund.

## Kategorie-Summary

HIGH 0 · MEDIUM 0 · LOW 0 · INFO 1 (neu, Lücke schon vor dem Slice vorhanden) · aus R4 offen: INFO 2.
Die Klasse `erkennungs-regex-nur-gegen-positivfaelle-getestet` tritt in dieser Runde nicht wieder
auf: Der Negativfall steht jetzt auch im umbrochenen Layout im Test.

## Verdikt

**Freigegeben.** R4 F-1 ist geschlossen, und der Fix bringt keinen neuen Fehler. F-1 dieser Runde
ist eine schon vorher bestehende, stille Lücke der Phrasen-Alternative. Sie blockiert nicht und
gehört in die Closure-Notiz: entweder als Risiko mit Ausgang oder als benannte Grenze.
