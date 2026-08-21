# Review-Report: slice-105 — Chronologie-Monotonie als siebte `structure`-Bedingung — 2026-08-21

**Review-Art:** unabhängiger Erst-Review (Code, gegen die Vertragsflächen) —
frischer Kontext, vor dem Release. Nicht geprüft wird die DoD-Abhakung und die
Gate-Lauf-Bestätigung (getrennter Kontext, Verifikation); der im Commit
behauptete Retro-Beleg (27 Befunde am Vor-welle-73-Stand) ist eine
Verifikations-, keine Review-Aussage.

**Gegenstand:** `slice-105` (Welle welle-77-chronologie-ordnung),
Feat-Commit `9b33d8b`; Arbeitsbaum-Stand `HEAD` = `9b33d8b`-Nachfolge (clean
zum Review-Zeitpunkt).

**Skill:** `.harness/skills/reviewer.md` @ 1.4.0 ·
**Modell:** claude-fable-5 · **Datum:** 2026-08-21

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  (Lastenheft-Fassung 0.61.0: Bedingungs-Tabelle, Absatz „Chronologie-Monotonie",
  fail-closed-Liste, Akzeptanzkriterien, Out-of-Scope, §7-Zeile)
- §[`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritte 1/5/6, §2-Schema (`structure[].table-order`/`table-column`),
  §4-Grund-Code-Tabelle (`section-unordered`/`section-cell-untyped`), §7-Zeile
- [ADR-0057](../plan/adr/0057-structure-tabellen-monotonie.md) (Proposed),
  dazu ADR-0044 (Schnitt-Kriterium), ADR-0049 (Modul-Schnitt), ADR-0054
  (geteilte Lexik), ADR-0005 (Import-Regeln), `DC-QA-02`/`DC-QA-03`
- Slice-Plan `docs/plan/planning/in-progress/slice-105-tabellen-monotonie.md`,
  Wellendokument `docs/plan/planning/welle-77-chronologie-ordnung.md`,
  Hard Rules `AGENTS.md` §3

Geprüfte Implementierung: `internal/hexagon/core/rules/structure_tableorder.go`
(+`_test`), `structure.go`, `markdown.go` (gehobene Tabellen-Lexik),
`planning_waves.go`, `lexikon_kopplung_test.go`, `model/config.go`,
`model/finding.go`, `app/diagnose.go`, `adapter/driven/configyaml/*`
(+`gate_consistency_test.go`), `adapter/driving/cli/config_template.go`,
`.d-check.yml`, `Makefile`.

---

## Findings

### F-1 · MEDIUM — Typ-Mischung: zwei Vertragslesarten, der Code folgt der einen, die andere steht daneben, der Zweig ist ungetestet

- **kategorie:** MEDIUM
- **quelle:** `DC-FA-STRUCT-001` / §`DC-FA-STRUCT-001.a` Schritt 6 / ADR-0057
  Entscheidung 7
- **pfad:** `internal/hexagon/core/rules/structure_tableorder.go:170–174` ·
  `spec/spezifikation.md:2032` und `spec/spezifikation.md:2068–2071` ·
  `spec/lastenheft.md:2346–2350` ·
  `internal/hexagon/core/rules/structure_tableorder_test.go` (Fall
  „Typ-Mischung", nur 2 Datenzeilen)
- **befund:** Die Spezifikation beschreibt dieselbe Prüfung zweimal
  verschieden. Die §6-Bedingungs-Tabelle (Zeile 2032) verlangt, dass „jede
  Schlüsselzelle typisierbar ist und den **Spalten-Typ** trägt"; der
  Schritt-6-Fließtext (2068–2071) sagt „ein anderer Typ als beim typisierbaren
  **Vorgänger** ⇒ `section-cell-untyped` … der Vergleich setzt beim nächsten
  typisierbaren Nachbar-Paar wieder auf", und ADR-0057 E7 formuliert „eine
  kaputte Zelle … **meldet sich selbst**" (Singular). Der Code implementiert
  die reine Paar-Lesart: nach einer Misch-Zelle bleibt sie der
  Vergleichs-Anker (`return &f, &key, lineNo`, Zeile 174). Folge bei der
  Spalte Datum–Version–Datum(–Datum…): **zwei** `section-cell-untyped` —
  eines an der fremden Zelle, eines an der **gesunden** Folge-Zeile, die den
  Spalten-Typ sehr wohl trägt. Nach der Tabellen-Lesart („trägt den
  Spalten-Typ") dürfte die zweite Zeile nicht melden; nach der Paar-Lesart
  muss sie. Kein Test pinnt den Kaskaden-Fall (der Bestandstest deckt nur die
  2-Zeilen-Mischung, bei der beide Lesarten zusammenfallen) — bei einem neuen
  öffentlichen Vertrag ist genau dieser divergente Zweig der ungetestete.
- **verifizierbar:** ja — Unit-Test mit Spalte `2026-08-16` / `v0.60.0` /
  `2026-08-10` (`table-order: desc`): der Lauf liefert heute **zwei**
  `section-cell-untyped` (Zeilen 3 und 4 der Tabelle); die §6-Tabellen-Lesart
  verspricht eines.
- **klasse:** vertrags-mehrdeutigkeit-messmethode
- **Übergabe (Fix-Richtung, außerhalb des Befunds):** eine von zwei kleinen
  Auflösungen — (a) Vertrag auf die Paar-Semantik pinnen: §6-Tabellenzelle von
  „den Spalten-Typ trägt" auf „den Typ ihres typisierbaren Vorgängers trägt"
  schärfen, ADR-E7-Formulierung angleichen, Kaskaden-Testfall (3 Zeilen, 2
  Befunde) festschreiben; oder (b) Code auf die Reset-Semantik angleichen (bei
  Mischung `nil`-Anker wie bei den beiden anderen Fehlerfällen zurückgeben)
  plus Testfall (3 Zeilen, 1 Befund). Beide sind minimal-invasiv; entscheidend
  ist, dass genau eine Lesart übrig bleibt und getestet ist.

### F-2 · LOW — `versionSegments`: der „unerreichbare" Fehlerpfad ist erreichbar und degradiert still

- **kategorie:** LOW
- **quelle:** Maintainability (latente Wartungsfalle) / `DC-FA-STRUCT-001`
  (fail-closed-Anspruch der Bedingung)
- **pfad:** `internal/hexagon/core/rules/structure_tableorder.go:65–77`
  (Kommentar Zeile 71)
- **befund:** Der Kommentar behauptet, der `strconv.Atoi`-Fehlerzweig sei
  „nach dem Muster nicht erreichbar" — das Muster `v?\d+(?:\.\d+)+` erlaubt
  aber beliebig lange Ziffernfolgen; ein Segment jenseits von int64 (z. B. ein
  20-stelliger Zahlen-Token in einer Prosa-Zelle der Schlüsselspalte) lässt
  `Atoi` scheitern, `segs` wird `nil`, und der Schlüssel vergleicht fortan
  still als kleinstmögliche Version (leere Segmentfolge) statt als
  untypisierbar zu melden — ein stiller Sonder-Pfad in einer Bedingung, deren
  erklärte Disziplin „Befund statt Übersprung" ist. Kein Test deckt den Zweig.
- **verifizierbar:** ja — Unit-Test mit Zelle `12345678901234567890.1` in
  einer sonst absteigenden Versions-Spalte: heute kein `section-cell-untyped`,
  und die Zeile gilt als kleinster Schlüssel.
- **klasse:** stiller-degradations-rand
- **Übergabe:** den Fehlerzweig als untypisierbar behandeln (ok=false nach
  oben reichen) oder die Segmentlänge im Muster/Parser begrenzen; den
  Kommentar in jedem Fall korrigieren.

### F-3 · INFO — Kopplungs-Reichweite: Kopf-/Trennzeilen- und Zell-Frage sind funktions-geteilt, aber nicht kopplungs-getestet

- **kategorie:** INFO
- **quelle:** ADR-0054 Entscheidung 4 (Form des Kopplungs-Tests)
- **pfad:** `internal/hexagon/core/rules/markdown.go:183–219` ·
  `internal/hexagon/core/rules/lexikon_kopplung_test.go` (neuer
  `TestTabellenzeilenFrageHatEineAntwort`)
- **befund:** Der neue Kopplungs-Test bindet die **Tabellenzeilen**-Frage über
  alle drei Konsumenten (`targets`, `planning.waves`, `structure`) — wirksam:
  jede abweichende Antwort eines Konsumenten macht ihn rot. Die mit diesem
  Slice erst geteilt gewordenen Antworten `tableHeaderOrSeparator` und
  `tableCells` haben je **zwei** Konsumenten (`planning.waves`, `structure`)
  und sind nur über die gemeinsame Funktion gebunden; keine Kopplung fährt
  eine kopf-/trennzeilen-förmige Eingabe durch beide und vergleicht die
  Antworten. Nach der ADR-0054-Schwelle (Kopplung ab dem dritten Konsumenten)
  ist das vertretbar — die Annahme ist aber nirgends notiert, und der dritte
  Konsument einer dieser beiden Fragen träte ohne Erinnerung ein.
- **verifizierbar:** nein (dokumentationswürdige Annahme; kein Gate schlägt
  heute an)
- **klasse:** kopplungs-schwelle-undokumentiert

### F-4 · INFO — Regel-Identität schließt eine zweite Chronologie-Aussage über denselben Abschnitt aus, unbenannt

- **kategorie:** INFO
- **quelle:** `DC-FA-STRUCT-001` (Regel-Identität/Duplikat ⇒ Exit 2;
  Out-of-Scope-Absatz)
- **pfad:** `spec/lastenheft.md:2256–2257` · `internal/hexagon/core/model/config.go`
  (`Identity` aus `files` + Abschnitts-Selektor)
- **befund:** Die Regel-Identität besteht aus Glob und Abschnitts-Selektor;
  zwei Regeln mit identischer Identität sind Exit 2. Damit ist die
  Chronologie-Bedingung die erste, deren Mehrfach-Instanziierung je Abschnitt
  nicht ausdrückbar ist: wer dieselbe Tabelle in **zwei** Spalten monoton
  zusagen will (etwa Spalte 1 Version und Spalte 2 Datum im
  Release-Register), erhält einen Konfigurations-Duplikat-Abbruch — laut,
  nicht still, aber als Grenze nirgends benannt (das Out-of-Scope nennt nur
  Ordnungs-Aussagen **über** Spalten-Grenzen hinweg, nicht zwei unabhängige
  Ein-Spalten-Zusagen).
- **verifizierbar:** ja — Config mit zwei Regeln gleicher `files`/`section`
  und verschiedenem `table-column` ⇒ Exit 2.
- **klasse:** unbenannte-ausdrucks-grenze

### F-5 · INFO — Zweite Zell-Antwort im Produkt: `app.splitPipeTableLine` neben `rules.tableCells`

- **kategorie:** INFO
- **quelle:** ADR-0054 (eine Frage, eine Antwort) / Maintainability
- **pfad:** `internal/hexagon/core/app/trace_table.go:351–375` ·
  `internal/hexagon/core/rules/markdown.go:214–219`
- **befund:** Die gehobene `tableCells`-Antwort (naiver Split an `|`, exakt
  der Vertragstext aus §`DC-FA-STRUCT-001.a` Schritt 6) trägt den Kommentar
  „die geteilte Zell-Antwort" — im Paket `app` existiert aber eine zweite,
  reichere Zell-Zerlegung (`splitPipeTableLine`: escaped Pipes,
  Backtick-Spans, optionale Rand-Pipes) für den RTM-/trace-Leser. Beide sind
  je für sich vertragstreu (verschiedene §§), aber es sind zwei Antworten auf
  „was sind die Zellen dieser Zeile" im selben Produkt; eine Schlüsselzelle
  mit Pipe im Backtick-Span würde in `structure`/`planning.waves` die
  Spaltenadresse verschieben, im trace-Leser nicht. Vorbestand, durch die
  Hebung nur sichtbarer geworden — die BEO-003-Klasse in latenter Form.
- **verifizierbar:** nein (heute kein abweichender Befund in den sechs
  aktivierten Tabellen; keine enthält Pipes in Zellen)
- **klasse:** BEO-003-latenz

---

## Negativbefunde (geprüft, ohne Befund)

1. **Vertrag ↔ Code, Typisierung erster Treffer:** geprüft, ohne Befund —
   `typeTableKey` nimmt den frühesten Treffer beider Muster
   (`d[0] <= v[0]`-Vorrang; ein Index-Gleichstand ist mit diesen Mustern
   nicht konstruierbar), exakt die „erster Treffer in der rohen
   Zelle"-Zusage; Inline-Code-Backticks und HTML-Anker überleben (Testfall
   `version.md`-Form vorhanden).
2. **Nicht-strikte Monotonie und Richtung je Regel:** geprüft, ohne Befund —
   Vorzeichenlogik `asc: c<0` / `desc: c>0` entspricht ≥/≤; gleiche Schlüssel
   in Folge sind grün (Test); Datum als ISO-String ist bei fester
   Muster-Breite ordnungstreu; Versionsvergleich segmentweise numerisch,
   kürzere Segmentfolge kleiner, `v`-Präfix ordnungsneutral (Tests
   `1.10 > 1.9`, `1.9 < 1.9.1`, `v1.10`).
3. **Je zusammenhängender Tabelle:** geprüft, ohne Befund — `inTable`-Reset
   bei jeder Nicht-Tabellenzeile (auch Fence-Zeile/Fence-Inhalt); zwei durch
   Prosa getrennte, je sortierte Tabellen sind kein Befund (Test). Zwei ohne
   Leerzeile aneinanderstoßende Tabellen verschmelzen zu einer Folge — das
   entspricht sowohl dem Vertragstext („zusammenhängende Folge") als auch dem
   GFM-Rendering.
4. **Kopf-/Trennzeilen-Skip:** geprüft, ohne Befund — geteilte Antwort
   (`tableHeaderOrSeparator`), Kopfzeile = Tabellenzeile unmittelbar vor
   einer Trennzeile, wortgleich zum Vertrag; eine Datenzeile vor einer
   verirrten Trennzeile gilt damit als Kopfzeile — das ist die bestehende,
   vertraglich fixierte Lexik-Antwort, keine Abweichung dieses Slice.
   Kein Blick über das Abschnitts-Ende hinaus (die Section endet stets an
   einer Überschrift oder EOF; eine Überschrift beginnt nie mit `|`).
5. **Leerlauf-Befund:** geprüft, ohne Befund — `dataRows == 0` je Abschnitt
   ⇒ `section-unordered` an der Abschnitts-Überschrift; deckt „gar keine
   Tabelle", „nur Kopf+Trennzeile" und „Tabelle nur im Fence" (Tests für
   Fall 1 und 3; Fall 2 folgt aus demselben Zähler).
6. **Reset nach untypisierbarer Zelle:** geprüft, ohne Befund für die Fälle
   „kein Token" und „zu wenige Zellen" — Anker wird genullt, das nächste
   typisierbare Nachbar-Paar vergleicht wieder (Test
   `TestTableOrderSetztNachUntypisierbarNeuAuf` bestätigt auch, dass
   Nicht-Nachbarn nicht verglichen werden). Der Misch-Fall ist F-1.
7. **Roh-Zeilen-Ausnahme im Vertrag (BEO-004-Frage):** geprüft, ohne Befund —
   die Ausnahme ist an allen Orten benannt, an die sie gehört:
   Lastenheft-Bedingungs-Intro (Zeile 2282–2285) **und**
   Chronologie-Absatz (2340–2344), Spezifikation Schritt 5 (2013–2016)
   **und** Schritt 6 (2055–2056), §2-Schema-Zeile (2530), ADR-0057
   Entscheidung 4 samt Re-Evaluierungs-Trigger („zweiter Anwender"),
   `--print-config`-Vorlage („ROHE Zellen"). Die Bedingung liest keine neuen
   Dateien — dieselben Kandidaten, roher Text; die Fence-Frage bleibt bei der
   geteilten Lexik (`proseLineSet`).
8. **Erreichbarkeit der Grund-Codes:** geprüft, ohne Befund — beide Codes
   werden von Code erzeugt und von Tests getroffen: `section-unordered` über
   Bruch-Zeile (mit Zeilennummern-Assertion) und Leerlauf (mit
   Zeilen- und Meldungs-Assertion); `section-cell-untyped` über alle drei
   Emissionspfade (kein Token / zu wenige Zellen / Typ-Mischung). Ungetestete
   Ränder sind genau die aus F-1 (Kaskade) und F-2 (Atoi-Zweig).
9. **Geteilte Lexik, keine Selbst-Antworten:** geprüft, ohne Befund — grep
   über `internal/` findet die Tabellen-Antworten nur noch in `markdown.go`;
   `targets.go:157`, `planning_waves.go:208/211` und
   `structure_tableorder.go:121/129/154` konsumieren. Die Hebung aus
   `planning_waves.go` ist eine textidentische Verschiebung (Diff geprüft);
   die Bestandstests beider Konsumenten sind unberührt. Einzige weitere
   Zell-Zerlegung ist der app-seitige trace-Leser (F-5, andere Vertragsfläche).
10. **Kopplungs-Test wirksam:** geprüft, ohne Befund — drei Schreibweisen
    (außerhalb Fence / im Fence / eingerückt) durch alle drei Konsumenten
    über deren echte Befund-Pfade; jede einseitige Divergenz (auch eine
    legitime Lexik-Änderung) macht mindestens einen Konsumenten ungleich der
    erwarteten gemeinsamen Antwort und den Test rot. Reichweiten-Grenze in
    F-3 notiert.
11. **Stille Grün-Pfade / fail-closed:** geprüft, ohne Befund — die drei
    neuen Config-Ränder brechen mit Exit 2 (`structureChronologieFehler`,
    Decode-Tests für alle drei plus Positiv-Durchreichung); Leerlauf meldet;
    `exempt-paths`-Nullmenge meldet weiter (Bestand); `sections: one` mit
    Mehrfach-Treffer bricht vor der Messung ab (kein Chronologie-Lauf über
    dem falschen Abschnitt); `sections: each` erzeugt je Abschnitt einen
    eigenen Leerlauf-/Bruch-Befund an dessen Zeilen.
12. **Dedup-Kollisionen:** geprüft, ohne Befund — Tupel ist
    (Datei, Zeile, Regel, Ziel, Grund) (`SortFindings`,
    `model/finding.go:85–117`); pro Datenzeile entsteht höchstens ein
    Chronologie-Befund, Bruch-Zeilen sind nie Überschrift-Zeilen (`|` vs
    `#`), der Leerlauf an der Überschrift kollidiert mit Prosa-Bedingungen
    nur reason-verschieden, und die Regel-Identität im `target` hält
    Regel-Paare auseinander. Zwei verschiedene Reparaturen unter einem Tupel:
    nicht konstruierbar.
13. **Regex-Ränder über F-2 hinaus:** geprüft, ohne Befund — ein Datum
    matcht nie die Versions-RE (kein Punkt), eine Version nie die Datums-RE
    (kein Bindestrich-Format); eingebettete Token (Datum in längerer
    Ziffernfolge, Version in Prosa) sind durch die „erster
    Treffer"-Zusage gedeckt; Kopfzeilen-Wörter erreichen die Typisierung
    nicht (Skip vor der Zell-Lesung, Test); semantisch unmögliche Daten
    (`2026-13-99`) vergleichen als String konsistent — der Vertrag verspricht
    Muster-, keine Kalender-Validierung.
14. **Selbst-Aktivierung `.d-check.yml`:** geprüft, ohne Befund — sechs
    Regeln, alle sechs Heading-Klartexte zeichengenau gegen die Live-Dateien
    verifiziert (`## 7. Historie` ×2, `## Historische
    Trigger-Verschiebungen`, `## Abgeschlossene Wellen`, `## Verlauf`,
    `## 11. Änderungshistorie`); Spaltenwahl stimmt mit den Tabellen überein
    (Roadmap-Wellenregister Spalte 2 = Datum, alle übrigen Spalte 1);
    Richtung je Regel wie Slice §5 (Handbuch als einzige `asc`); jede
    Heading-Umbenennung fiele als `section-missing` laut aus.
15. **Gate-Nebenwirkungen:** geprüft, ohne Befund — `FOCUS_DISABLE` trägt
    `--disable structure` und wirkt auf alle vier fokussierten Gates
    (`gate-consistency`, `planning-check`, `trace-check`, `adr-check`) —
    kein `planning-check`-Nebeneffekt (Abnahme-Punkt 2 des Slice);
    `verify-closure-notes` fährt das eigene Profil `.d-check.closure.yml`,
    dessen drei `structure`-Regeln keine neuen Schlüssel tragen
    (byte-identisches Verhalten, ausdrücklich gemessen statt behauptet);
    `--print-mk` hat sein `doc-structure`-Target unverändert (kein neues
    Target, Target-Zahl konstant).
16. **Spiegel-Vollständigkeit (MR-025), per grep:** geprüft, ohne Befund —
    `AllReasons()` + `reasonTexts()` +2 im selben Commit; §4-Tabelle +2 und
    beidseitig maschinell verriegelt
    (`TestAllReasonsDeckungGegenSpezifikationGrundCodes` parst die
    Spec-Tabelle); §2-Schema +2; Lastenheft (Tabelle, Absatz, fail-closed,
    sechs neue Akzeptanzkriterien, Out-of-Scope, §7-Zeile 0.61.0 **oben**);
    Spezifikation §7-Zeile **oben**; ADR-Index +1;
    `--print-config`-Vorlage +2 (Vorlage bleibt parser-rund, Test
    vorhanden); Netzlos-Modulliste + Guard-Fall rotiert;
    `operations.md`-Enumerationen korrekt unverändert (structure war
    gelistet, keine neue CLI-Option); README/CHANGELOG/`version.md`/Handbuch
    bewusst Release-Prep — deren Fehlen ist hier vertragsgemäß, die
    §11-Positionsregel (chronologisch unter die letzte) ist ab diesem Slice
    maschinell.
17. **Determinismus (DC-QA-02):** geprüft, ohne Befund — keine
    Map-Iteration im Bedingungs-Pfad (das `prose`-Set wird nur punktuell
    abgefragt), Kandidaten stabil sortiert, Befunde in Zeilen-Reihenfolge,
    `SortFindings` kanonisiert; ohne `table-order` ist der Befundsatz
    byte-identisch (Opt-in-Test je Regel; Modul-aus-Fall unverändert im
    Bestand).
18. **Hermetik (DC-QA-03):** geprüft, ohne Befund — keine neuen Eingaben
    außer den regel-benannten Markdown-Dateien, kein git, kein Netz;
    `structure` in `netlessDocModules()` nachgezogen und beidseitig
    verriegelt (fehlendes Modul und Netz-Modul machen den Live-Test rot).
19. **ADR-0005-Import-Regeln:** geprüft, ohne Befund —
    `structure_tableorder.go` importiert nur stdlib und `core/model`.
20. **Referenz-Richtung (SDP)/Marker-Ehrlichkeit:** geprüft, ohne Befund —
    ADR-0057 nennt `slice-105` nur in der Geschichte (matrix-ausgenommen,
    Provenance-Ort); kein Provenance-Marker im Körper, keine getarnte
    Entscheidungsgrundlage; der ADR-Index (kein Matrix-Klassen-Mitglied)
    trägt den Slice als Status-Provenance wie bei den Vorgängern.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 |
| LOW | 1 | F-2 |
| INFO | 3 | F-3, F-4, F-5 |

## Verdikt

**Blockierend (ein MEDIUM), mit klar umrissenem, kleinem Ausweg.** F-1 ist
eine echte Zwei-Lesarten-Stelle im soeben geschriebenen Vertrag — die
§6-Bedingungs-Tabelle („Spalten-Typ") und der Schritt-6-Fließtext
(„Vorgänger-Typ") versprechen bei der Spalte Datum–Version–Datum verschiedene
Befundbilder, der Code liefert das eine, kein Test pinnt es, und vor dem
Release ist der einzige Zeitpunkt, zu dem sich das ohne
Grund-Code-Semantik-Nachtrag klären lässt. Der Fix ist in beiden Richtungen
minutenklein (ein Satz plus ein Testfall, oder zwei Zeilen plus ein
Testfall). F-2 sollte im selben Zug mitgenommen werden (Kommentar ist
faktisch falsch, Zweig widerspricht der erklärten fail-closed-Disziplin),
blockiert für sich allein aber nicht. Alles Übrige — Lexik-Hebung,
Kopplungs-Test, Config-Ränder, Selbst-Aktivierung, Spiegel, Gates,
Determinismus, Hermetik — ist ohne Befund und approve-fähig.
