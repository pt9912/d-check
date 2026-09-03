# Review-Report: welle-72-closure-semantik (slice-094 + slice-104) — 2026-08-10

**Review-Art:** Code — geprüft wird der Diff gegen Lastenheft, Spezifikation,
[ADR-0053](../plan/adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md)
und beide Slice-Pläne; die DoD-Abhakung ist ausdrücklich **nicht** Gegenstand
(das ist Verifikation, anderer Kontext).

**Gegenstand:** `ae1a7ee..HEAD` (`877e028`) — fünf Commits: Bestandsmessung,
CR + ADR, Implementierung slice-094, Lifecycle-Move, Implementierung slice-104.

**Skill:** `.harness/skills/reviewer.md` @ 1.3.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-08-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [welle-72-closure-semantik](../plan/planning/welle-72-closure-semantik.md), besonders §5 (die Kopplung beider Slices)
- [slice-094](../plan/planning/done/welle-72/slice-094-closure-zaehl-paritaet.md) und [slice-104](../plan/planning/in-progress/slice-104-floskel-wortgrenze.md)
- [ADR-0053](../plan/adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md) (Proposed) und [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (Lastenheft 0.56.0) und [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) Schritt C4/C4b
- [`AGENTS.md`](../../AGENTS.md) Hard Rules, besonders [§3.6](../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)
- das `verify-closure-notes`-Skript des Schwester-Repos (abzulösendes Adopter-Werkzeug)

**Eigene Läufe dieses Reviews:** `make test` (Baseline + 12 Mutationsläufe),
`make gates`, `make adr-check`, `make verify-closure-notes`, dazu Image-Läufe
gegen sechs Fixture-Repos in einem Temp-Verzeichnis und gegen eine Kopie des
Adopter-Bestands. Der Arbeitsbaum ist unverändert (Beleg am Ende).

---

## Findings

### F-1 — Satzende vor CRLF-Zeilenende zählt nicht mehr

- `kategorie`: HIGH
- `quelle`: [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) Schritt C4 („nur vor Whitespace oder Zeilenende"), [`DC-QA-04`](../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
- `pfad`: `internal/hexagon/core/rules/planning.go:360`
- `befund`: Als Folgezeichen akzeptiert die Zählung nur Leerzeichen, Tab und
  Zeilenumbruch. In einer CRLF-Arbeitskopie steht vor dem Zeilenumbruch ein
  Wagenrücklauf, und damit zählt **kein einziges** zeilenschließendes Satzende
  mehr: eine Closure-Notiz mit vier sauberen Sätzen (je einer pro Zeile) ergibt
  gemessen **0** statt 4 und meldet `closure-note-thin`; das abzulösende
  Adopter-Skript zählt dieselbe Notiz mit **4**, weil sein `[[:space:]]` den
  Wagenrücklauf einschließt. Vor dieser Änderung war die Zeilenende-Form für die
  Zählung folgenlos (jedes Satzzeichen zählte). Der Rest des Produkts arbeitet
  auf CRLF korrekt — ein CRLF-Fixture mit Anker- und Datei-Link läuft mit
  `links`/`anchors` befundfrei durch, und die Abschnitts-Überschrift derselben
  CRLF-Notiz wird gefunden (der Befund trägt die Zeilennummer).
- `verifizierbar`: ja — Image-Lauf gegen eine CRLF-Fixture (gemessen: Modul 0,
  Adopter-Pipeline 4). `make verify-closure-notes` im eigenen Repo bleibt grün,
  weil hier keine CRLF-Datei liegt; der Fall ist über die eigenen Gates
  unsichtbar.
- `klasse`: „Zeilenende-Lexik ignoriert CRLF"

### F-2 — Widerrufene Aussage am Nachbar-Schritt: „die Zählung sieht Inline-Code weiterhin"

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in), [ADR-0053](../plan/adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md) Entscheidung 1
- `pfad`: `spec/lastenheft.md:2006` (gleichlautend `spec/spezifikation.md:1793`)
- `befund`: Der Platzhalter-Absatz beider Vertragsdokumente sagt weiter, die
  Substanz-Zählung bleibe unberührt und sehe Inline-Code weiterhin, die engere
  Sicht gelte nur für C4b. Ein Lauf widerlegt das: eine Notiz, deren fünf Sätze
  vollständig in Inline-Code stehen, meldet `closure-note-thin` — genau so, wie
  es die Fitness Function der ADR verlangt. Dieselbe Spezifikation trägt
  zusätzlich zwei gegenläufige Sätze: die §-Historie-Zeile behauptet, C4b leere
  Inline-Code „nicht mehr selbst", der C4b-Text darüber spricht unverändert vom
  **zusätzlichen** Leeren der Spans auf dem C4-Text.
- `verifizierbar`: ja — `--config`-Profil mit einer Notiz, deren Sätze in
  Inline-Code stehen; der Lauf meldet `closure-note-thin`, Exit 1.
- `klasse`: „Widerrufene Fassung an der Nachbar-Bedingung stehengeblieben"

### F-3 — Akzeptanzkriterium „Negative (Floskel)" ist gegen die Umsetzung falsifizierbar

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `pfad`: `spec/lastenheft.md:2062`
- `befund`: Das Kriterium verlangt einen Befund, sobald der **literale
  Teilstring** einer deklarierten Phrase case-insensitiv im bereinigten
  Abschnitts-Text vorkommt. Gemessen mit der Phrase `ok` und der Notiz „Der
  Ablauf ist dokumentiert.": 0 Befunde, Exit 0 — das Kriterium fordert einen
  Befund, den der Lauf nicht liefert. Die Lastenheft-Historie führt exakt diese
  Klasse (Kriterium gegen die Umsetzung falsifizierbar) schon einmal als
  behobenen Defekt, Eintrag 0.53.1.
- `verifizierbar`: ja — Image-Lauf mit `boilerplate` auf eine kurze Phrase und
  einer Notiz, die sie nur wortintern trägt.
- `klasse`: „Akzeptanzkriterium gegen die Umsetzung falsifizierbar"

### F-4 — §2-Schema und §4-Grund-Code-Tabelle tragen die alte Semantik

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
- `pfad`: `spec/spezifikation.md:2324` (dazu `:2322`, `:2435`, `:2436`)
- `befund`: Die normativen Referenztabellen desselben Dokuments, dessen Schritt
  C4 geändert wurde, beschreiben unverändert die widerrufene Fassung: §2 nennt
  `planning.closure.boilerplate` „literale Floskel-**Teilstrings**" und
  `min-sentences` „Satzende-Zeichen … **nach** Entfernen der
  Fenced-Code-Blöcke"; §4 wiederholt beides in den Zeilen zu
  `closure-note-thin` und `closure-note-boilerplate`. Ein Adopter konfiguriert
  über §2 und §4 — er erwartet danach Teilstring-Treffer und eine Zählung, die
  Inline-Code mitzählt; beides liefert der Lauf nicht.
- `verifizierbar`: ja — derselbe Lauf wie F-3 (Teilstring) bzw. F-2
  (Inline-Code).
- `klasse`: „Rand nicht nachgezogen (§2/§4 gegen den geänderten Schritt)"

### F-5 — Die Lockerung ist breiter als „eine Phrase in Backticks"

- `kategorie`: MEDIUM
- `quelle`: [`AGENTS.md` §3.6](../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden), [ADR-0053](../plan/adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md) Entscheidung 2
- `pfad`: `docs/plan/adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md:47`
- `befund`: ADR und Vertragstext beschreiben die angenommene Lockerung
  durchgehend als „eine *zitierte* Floskel ist keine benutzte" — eine Phrase in
  Backticks trifft nicht mehr. Gemessen ist sie breiter, weil die geteilte
  Absatz-Paarung auch **Fließtext** leert. Zwei Fälle nachgestellt, jeweils
  gegen die alte und die neue Fassung: (a) eine benutzte Floskel im Fließtext,
  in deren Absatz zwei einzelne Backtick-Zeichen stehen (etwa eine Notiz, die
  über Backticks schreibt) — alte Fassung meldet
  `closure-note-boilerplate`, neue Fassung meldet **nicht**; (b) ein
  Backtick-Zeichen **oberhalb** der Abschnitts-Überschrift, das mit einem
  Zeichen im Abschnitt paart — die Floskel im gemessenen Abschnitt verschwindet,
  obwohl das öffnende Zeichen außerhalb des Abschnitts steht (die alte,
  abschnittslokale Extraktion konnte das nicht). Der Absatz des Lastenhefts, der
  die Paarungs-Grenze benennt, ist derselbe, der laut F-2 bestreitet, dass sie
  für die Zählung gilt.
- `verifizierbar`: ja — je ein Lauf mit dem aktuellen Image und mit einem
  Image der Vorgänger-Semantik gegen dieselbe Fixture.
- `klasse`: „Lockerung enger beschrieben als gemessen"

### F-6 — Kalibrier-Kommentar im Closure-Profil steht auf der Messung vor der Angleichung

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `.d-check.closure.yml:46`
- `befund`: Der Kommentar über `min-sentences` sagt, das Bestands-Minimum liege
  bei 7 Satzende-Zeichen und der erste Slice werde bei 8 rot. Mit dem gebauten
  Image nachgemessen: Minimum **5**, zweitkleinster Wert **6** — rot also
  bereits ab 6. Die Datei ist in diesem Diff geändert worden (Floskel-Liste plus
  neuer Kommentarblock); der Kalibriersatz zwei Absätze darüber nicht. Wer die
  Schwelle das nächste Mal anhebt, rechnet nach diesem Kommentar mit drei
  Stufen Luft und hat eine.
- `verifizierbar`: ja — Lauf mit hochgedrehtem `min-sentences` gegen den
  `done/`-Bestand.
- `klasse`: „Kalibrier-Kommentar überlebt die Messung, die er zitiert"

### F-7 — Zwei Rückbauten bleiben grün: Weitersuche und Großbuchstaben-Zweig

- `kategorie`: LOW
- `quelle`: [slice-104](../plan/planning/in-progress/slice-104-floskel-wortgrenze.md) DoD (sieben Rückbauten, alle rot)
- `pfad`: `internal/hexagon/core/rules/planning.go:283`
- `befund`: Eigene Mutations-Gegenprobe: (a) die Weitersuche nach einem
  verworfenen Treffer um eine Position wird durch einen Sprung hinter den
  Treffer ersetzt — die gesamte Suite bleibt **grün**, obwohl die Funktion
  danach überlappende Vorkommen verliert (nachgestellt an einer Phrase, deren
  Anfang gleich ihrem Ende ist: Treffer kippt auf kein Treffer). Die
  DoD-Zusage „Weitersuche nach einem verworfenen Treffer" ist damit nur für den
  **nicht überlappenden** Fall gehalten. (b) Der Großbuchstaben-Zweig der
  Wortzeichen-Prüfung ist über die Modul-Oberfläche unerreichbar, weil Text und
  Phrase vorher kleingeschrieben werden; sein Entfernen lässt die Suite ebenfalls
  grün.
- `verifizierbar`: ja — beide Rückbauten auf einer Dateikopie, `make test` grün.
- `klasse`: „Mutations-Zusage deckt nur den nicht-überlappenden Fall"

### F-8 — Die Produkt-Ausgabe `--print-config` beschreibt die widerrufene Semantik

- `kategorie`: LOW
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `pfad`: `internal/adapter/driving/cli/config_template.go:130`
- `befund`: Die Vorlage, die ein Adopter über `--print-config` in sein Repo
  kopiert, erklärt `min-sentences` als „Satzende-Zeichen AUSSERHALB
  Fenced-Code" und `boilerplate` als „literale Floskeln". Das ist Produkt-,
  nicht Doku-Text: er wandert in die Konfigurationsdatei des Konsumenten und
  überlebt dort das Release.
- `verifizierbar`: ja — `--print-config` ausgeben und mit Schritt C4 vergleichen.
- `klasse`: „Produkt-Ausgabe trägt die alte Semantik"

### F-9 — Die ASCII-Wortgrenze steht im Vertrag nur in der günstigen Richtung

- `kategorie`: INFO
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `pfad`: `spec/lastenheft.md:1939`
- `befund`: Der Vertrag nennt Satzzeichen, Anführungszeichen, Bindestrich und
  Zeilenrand als Grenzen; die Spezifikation ergänzt die ASCII-Menge und
  begründet sie an den **helfenden** Formen. Die Gegenrichtung steht nur im
  Slice: gemessen trifft die neu aufgenommene Phrase `fertig` **innerhalb** eines
  Kompositums, sobald ihr unmittelbar ein Nicht-ASCII-Buchstabe vorausgeht. Und
  eine Phrase, die selbst nur aus Nicht-Wortzeichen besteht, trifft überall dort,
  wo sie steht — für solche Phrasen ist die Beschreibung „Wortgrenze in
  RE2-Semantik" aus dem Slice nicht deckungsgleich mit der Umsetzung, wohl aber
  der Vertragssatz „vor und hinter dem Treffer kein Wortzeichen".
- `verifizierbar`: ja — direkter Funktionstest bzw. Lauf mit einer solchen
  Phrase.
- `klasse`: „ASCII-Wortgrenze in der Prosa nur in der günstigen Richtung benannt"

### F-10 — Die Fitness Function ist universell formuliert, belegt ist sie am Korpus

- `kategorie`: LOW
- `quelle`: [ADR-0053](../plan/adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md) §Fitness Function
- `pfad`: `docs/plan/adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md:97`
- `befund`: Zugesagt ist: eine Fixture, die das Adopter-Skript wegen der
  Satzzählung rot macht, macht auch das Modul rot. Fünf konstruierte
  Eingabe-Klassen zählen verschieden — mehrzeiliger Inline-Span (2 gegen 3),
  Span mit doppelter Backtick-Folge (2 gegen 5), CRLF (0 gegen 4, siehe F-1),
  Tilde-Fence (2 gegen 6) und eine Fence-Infozeile, die selbst ein
  Backtick-Zeichen trägt (4 gegen 1). Der letzte Fall verletzt die zugesagte
  Richtung: das Adopter-Skript zählt 1 und ist bei seiner Schwelle 2 rot, das
  Modul zählt 4 und ist bei Schwelle 4 grün. Am realen Bestand tritt keine
  dieser Klassen auf — die Zusage gilt dem Korpus, nicht der Lexik.
- `verifizierbar`: ja — die fünf Fixtures gegen beide Zähler.
- `klasse`: „Fitness Function universell formuliert, korpus-belegt"

## Negativbefunde

- geprüft, ohne Befund: **Paritäts-Beleg 84/84**, eigenständig nachgerechnet.
  Kopie des Adopter-Bestands (84 Slice-Dateien) gegen die nachgebildete
  Shell-Pipeline des Schwester-Repos und gegen das gebaute Image mit
  hochgedrehter Schwelle: **0 Abweichungen**, Minimum 3 auf beiden Seiten, an
  der Adopter-Schwelle 2 beidseitig 0 rot. Die Zahl im Slice stimmt.
- geprüft, ohne Befund: **Bestandsmessung 7 → 5**. Mit dem gebauten Image und
  einem Image der Vorgänger-Semantik gegen den eigenen `done/`-Bestand
  gefahren: Minimum 7 (alt) gegen 5 (neu), bei Schwelle 4 beide grün, 0 Dateien
  unter der Schwelle.
- geprüft, ohne Befund: **die 170 fett gesetzten Satzenden**. Eigene
  Nachmessung über den heutigen Bestand (98 Notizen statt 97 zur Messzeit):
  4148 Satzende-Zeichen, davon 1357 vor Whitespace oder Zeilenende, 178 vor
  einem Auszeichnungszeichen, rund 2566 in Pfaden und Versionsnummern. Die
  Struktur der Tabelle im Slice ist reproduziert; die Zahlen liegen genau um die
  eine hinzugekommene Notiz höher.
- geprüft, ohne Befund: **Richtung der Zählung**. Datei für Datei verglichen
  (98 Notizen, alte gegen neue Semantik): **kein einziger** Wert steigt,
  mittlerer Verlust 33,7, maximaler 101. `closure-note-thin` findet
  ausschließlich mehr — die Minor-Einordnung trägt.
- geprüft, ohne Befund: **Mutations-Echtheit**. Neun eigene Rückbauten auf einer
  Dateikopie, jeder einzeln, jeder rot: Satzende-Form ganz entfernt; Tab aus der
  Whitespace-Menge; Zählung am Textende abgeschaltet; Abschnittstext ohne
  Inline-Code-Bereinigung; Wortgrenzen-Vergleich zurück auf Teilstring; Ziffern
  als Nicht-Wortzeichen; Unterstrich als Nicht-Wortzeichen; nur die linke
  Grenze; nur die rechte Grenze. Zwei weitere bleiben grün — das ist F-7.
- geprüft, ohne Befund: **Ränder der Wortgrenzen-Funktion**. Leere Phrase,
  Phrase gleich dem ganzen Text, Phrase länger als der Text, Regex-Metazeichen
  in der Phrase, überlappende Vorkommen, Phrase mit Zeichen am Rand, die keine
  Wortzeichen sind. Kein Fehlverhalten, kein Absturz, keine Regex-Kompilierung
  je Phrase.
- geprüft, ohne Befund: **`n/a` in der eigenen Liste**. Nachbarschaften
  gemessen: in Wortmitte („ein/aus", „bahn/ampel") trifft es nicht, als
  eigenständiger Ausdruck und in Klammern trifft es. Die Aufnahme ist sicher im
  zugesagten Sinn.
- geprüft, ohne Befund: **Lifecycle-Spur**. `git mv` nach `done/`, alle
  Verweise in Roadmap, Wellendokument, Nachbar-Slices und Ergebnisnotizen
  nachgezogen, Tombstone im Referenz-Ventil für den eingefrorenen
  `in-progress`-Pfad ergänzt. `planning-check` grün, kein toter Link.
- geprüft, ohne Befund: **Modul-Schnitt und Importe**. Beide neuen Helfer liegen
  im `rules`-Paket, keine neue Abhängigkeit, keine Regex-Kompilierung im
  Prüfpfad, `make arch-check` grün.
- geprüft, ohne Befund: **Byte-Identität und Opt-in**. Ohne
  `planning.closure.dir` wird keine Datei geöffnet; keine neue Config-Achse,
  kein Ventil — die ADR-Entscheidung „kein Schalter" ist im Code eingehalten.
- geprüft, ohne Befund: **Gate-Läufe im Ist-Zustand**: `make gates` grün,
  `make verify-closure-notes` 0 Befunde bei unveränderter Schwelle 4,
  `make adr-check` grün, `make test` grün.
- geprüft, ohne Befund: **ADR-Status `Proposed`** nach der Slice-Closure — das
  entspricht der eigenen Praxis (ADR-0050 steht nach abgeschlossener Welle
  ebenfalls auf `Proposed`); kein Bruch mit der Konvention.
- geprüft, ohne Befund: **Referenz-Richtung**. Kein Provenance-Marker im Diff,
  keine Slice-Kennung im ADR-Körper, `matrix` grün.
- geprüft, ohne Befund: **Netz und Seiteneffekte** — alle Läufe netzlos und
  read-only, kein neuer Zugriff außerhalb des Filesystem-Ports.

## Zur Wellen-Frage: deckt **eine** Notiz beide Slices?

Ja, wenn sie **drei** Richtungen nennt statt zwei: `closure-note-thin` findet
mehr (Zähl-Angleichung), `closure-note-boilerplate` findet weniger — und zwar
aus **zwei** unabhängigen Gründen (Inline-Code fällt weg **und** der Vergleich
läuft an Wortgrenzen). Das Wellendokument §5 und beide Slices nennen die zwei
Gründe; die vorbereiteten Historie-Einträge in Lastenheft und Spezifikation
ebenfalls. Die Einordnung als **Minor** ist für beide Richtungen belegt: kein
Wert der eigenen Messung steigt (Zählung), und der Wortgrenzen-Vergleich ist
eine echte Teilmenge des Teilstring-Vergleichs — ein roter Lauf kann grün
werden, ein grüner rot, beides ist mehr als eine Fehlerkorrektur und weniger als
ein Vertragsbruch.

**Offen für die Release-Prep** (kein Finding, aber Voraussetzung des
Wellen-Closure-Triggers): `CHANGELOG.md` trägt einen leeren
`[Unreleased]`-Block; `version.md` (Anker-Wanderung), die `ghcr`-Pins in beiden
READMEs und im Benutzerhandbuch, der Handbuch-Kopfstempel und die
§11-Verlaufszeile (chronologisch unter die letzte), der Aufgaben-Abschnitt
§4.17 samt Konfigurationstabelle in §5/§6 sowie die Modulzeile in
`harness/README.md` beschreiben durchgehend die alte Semantik. Neues Modul und
neue CLI-Option gibt es nicht, `operations.md` ist also nicht betroffen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 4 |
| LOW | 4 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Zeilenende-Lexik ignoriert CRLF · Widerrufene
Fassung an der Nachbar-Bedingung stehengeblieben · Akzeptanzkriterium gegen die
Umsetzung falsifizierbar · Rand nicht nachgezogen (§2/§4 gegen den geänderten
Schritt) · Lockerung enger beschrieben als gemessen · Kalibrier-Kommentar
überlebt die Messung, die er zitiert · Mutations-Zusage deckt nur den
nicht-überlappenden Fall · Produkt-Ausgabe trägt die alte Semantik ·
ASCII-Wortgrenze in der Prosa nur in der günstigen Richtung benannt · Fitness
Function universell formuliert, korpus-belegt

## Verdikt

**Merge-blockierend:** ja — ein HIGH und vier MEDIUM.

**Release-Empfehlung: noch nicht taggen.** Die Sache selbst ist solide: die
Zähl-Parität ist unabhängig nachgerechnet und stimmt auf allen 84
Adopter-Notizen, die Bestandsmessung stimmt, die Richtung der Zählung ist über
98 eigene Notizen monoton, und neun von elf geprüften Rückbauten fallen rot. Was
fehlt, ist die Kante:

1. **F-1 vor dem Tag.** Ein ausgeliefertes Gate bekommt in diesem Release eine
   neue Falsch-Positiv-Klasse, die im eigenen Repo unsichtbar ist und genau die
   Parität verfehlt, für die die Welle angetreten ist. Entweder die Zeilenende-
   Form deckt CRLF ab, oder die Grenze steht benannt im Vertrag und in der
   Release-Notiz — still ausliefern ist die einzige Variante, die nicht geht.
2. **F-2 bis F-4 vor dem Tag.** Drei Vertragsflächen sagen nach diesem Diff das
   Gegenteil der Umsetzung, zwei davon per Lauf falsifizierbar. Das Repo hat
   diese Klasse in derselben Anforderung schon zweimal geheilt.
3. **F-5 gehört in die ADR**, die es nach der Hard Rule ohnehin nur gibt, um
   die Lockerung zu benennen. Eine Lockerung, deren Reichweite größer ist als
   ihre Begründung, ist nicht vollständig entschieden.

F-6 bis F-10 sind Nachzug und können mit der Release-Prep laufen. Danach: **ein**
Minor-Release mit **einer** Notiz, die alle drei Richtungen einzeln nennt.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein
Lauf-Beleg (dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und
ersetzt keine Verifikation.
