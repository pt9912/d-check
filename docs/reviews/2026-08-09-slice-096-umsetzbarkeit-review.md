# Review-Report: slice-096 — Umsetzbarkeit des Moduls `structure` — 2026-08-09

**Review-Art:** Plan/Design — geprüft wird die neu spezifizierte Anforderung
gegen den **vorhandenen Code**, gegen die eigenen Zusagen der Spezifikation und
gegen den geschnittenen Umsetzungspfad. Leitfrage: „lässt sich das so bauen,
und was bricht dabei?" — nicht „ist der Text schön".

**Gegenstand:** Diff-Range `45246bb..acbb419` (drei Commits: Abnahme-Punkte,
CR `structure` + ADR + Folge-Slices), Slice
[slice-096](../plan/planning/in-progress/slice-096-structure-modul-analyse.md).

**Skill:** `.harness/skills/reviewer.md` @ 1.3.0 (`419e82e`, 2026-08-03)
**Modell:** claude-opus-5[1m] · **Datum:** 2026-08-09

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  und [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure),
  [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  und die neuen §2-Schema-Zeilen zu `structure[]`
- [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md),
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md),
  [ADR-0042](../plan/adr/0042-markdown-lexik-folgt-commonmark.md)
- Folge-Slices [slice-099](../plan/planning/open/slice-099-structure-modul.md)
  und [slice-100](../plan/planning/open/slice-099-structure-modul.md)
- Bestehender Code: `internal/hexagon/core/rules/planning.go`,
  `internal/hexagon/core/rules/markdown.go`,
  `internal/hexagon/core/model/finding.go`,
  `internal/hexagon/core/rules/scan.go`,
  `internal/adapter/driven/configyaml/configyaml.go`
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules),
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus),
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)

**Belegläufe:** Fixtures in einem Temp-Verzeichnis außerhalb des Repos, geprüft
mit `docker run --rm --network none -v <fixture>:/repo:ro d-check:latest`
(ausgeliefertes Image, enthält die Closure-Fähigkeit).

---

## Findings

### F-1 — Unbalancierter/verschachtelter Fence verschluckt den Abschnitt, die Bedingung läuft grün

- `kategorie`: HIGH
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  §Akzeptanzkriterien „Boundary (fence-treu)"; [ADR-0042](../plan/adr/0042-markdown-lexik-folgt-commonmark.md)
- `pfad`: `spec/spezifikation.md:1748` (Schritt 4 „Bereinigen"), `spec/lastenheft.md:2024`
- `befund`: Die Bereinigung ist auf die vorhandene `FenceToggle`-Lexik
  festgeschrieben, die ein naiver Umschalter ist: ein Öffner ohne Schluss und
  eine Backtick-Fence innerhalb einer Tilde-Fence kippen den Zustand, sodass
  Abschnitts-Inhalt als „im Code" gilt und aus der Messung fällt. Der Vertrag
  sagt nur die eine Richtung zu („was im Fence steht, zählt nicht"); die andere
  („was nicht im Fence steht, zählt") ist nirgends zugesagt, und alle
  Verbots-/Pflicht-Bedingungen (`forbid-pattern`, `require-all`, `require-any`,
  `max-tasks`) melden dann **nichts** statt zu melden. Beleg über dieselbe
  Mechanik im ausgelieferten Preset: eine Closure-Notiz mit einem Satz, danach
  einem ungeschlossenen Backtick-Fence und dahinter der wörtlich konfigurierten
  Floskel liefert `0 Befund(e)`, Exit 0 — die Floskel steht in der Datei und
  wird nicht gefunden. Umgekehrt liefert eine Floskel **innerhalb** einer
  Tilde-Fence, die eine Backtick-Fence enthält, `closure-note-boilerplate` —
  ein Befund auf Code-Inhalt.
- `verifizierbar`: ja — Fixture mit ungeschlossenem Fence im geprüften
  Abschnitt, Lauf mit aktiviertem Modul; erwartet würde ein Befund, geliefert
  wird Exit 0 (bereits gegen `d-check:latest` mit `--enable planning` gezeigt).
- `klasse`: Naiver Fence-Toggle macht Verbots-Bedingung still grün

### F-2 — „mehrere verletzte Bedingungen ⇒ mehrere Befunde" ist mit der globalen Befund-Deduplikation nicht erreichbar

- `kategorie`: HIGH
- `quelle`: [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 5; [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)
- `pfad`: `spec/spezifikation.md:1751` gegen `internal/hexagon/core/model/finding.go:50`
- `befund`: `SortFindings` entfernt Befunde mit identischem Tupel
  (Datei, Zeile, Regel, Ziel, Grund) — die `message` gehört nicht dazu. Für
  `structure` sind bei zwei verletzten Bedingungen desselben Abschnitts alle
  fünf Felder gleich: dieselbe Datei, dieselbe Zeile (Schritt 6 legt nur
  „`line` wie oben" fest, also die Überschriftszeile), `rule` = `structure`,
  `target` = derselbe `files`-Glob, `reason` = `section-constraint`. Die
  Zusage „mehrere verletzte Bedingungen ⇒ mehrere Befunde" kollabiert damit auf
  einen Befund, dessen überlebende Meldung von der Einfüge-Reihenfolge abhängt.
  Dasselbe trifft zwei Regeln mit gleichem `files`-Glob, die beide
  `section-missing` auf Zeile 1 derselben Datei melden. Beleg für die
  Dedup-Mechanik: zwei identische kaputte Links in derselben Zeile ergeben
  **einen** Befund, zwei verschiedene Ziele in einer Zeile ergeben zwei.
- `verifizierbar`: ja — Fixture mit einem Abschnitt, der `min-sentences` **und**
  `forbid-pattern` verletzt; die Zählung in der Zusammenfassung entscheidet.
- `klasse`: Sammel-Grund-Code kollidiert mit Befund-Dedup

### F-3 — Die Marken-Verankerung schließt Listen-Items aus und trifft damit die eigene Motiv-Dokumentklasse nicht

- `kategorie`: HIGH
- `quelle`: [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  Entscheidung 5; [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  §Benannte Marken
- `pfad`: `spec/spezifikation.md:1761`, `spec/lastenheft.md:1986`
- `befund`: Zugesagt ist, dass eine Marke gilt, wenn eine Zeile „nach führendem
  Whitespace" mit der hervorgehobenen Marke beginnt. Ein Listen-Marker ist kein
  Whitespace. Die Dokumentklasse, mit der die Entscheidung begründet wird
  („eine Anforderung muss alle Akzeptanz-Bausteine tragen"), schreibt ihre
  Bausteine in diesem Repo ausschließlich als Listen-Items: 267 Zeilen in
  `spec/lastenheft.md` beginnen mit Bindestrich plus doppeltem Stern, davon 65
  mit der Marke „Boundary" — und **null** Zeilen beginnen bare mit
  „Happy Path" oder „Boundary" in Fettschrift. Eine Regel `require-all` mit
  diesen Marken meldete jede Anforderung als verletzt. Dieselbe Lücke trifft
  Marken in Tabellenzellen (die Zeile beginnt mit einem Pipe-Zeichen) und
  Marken in eingerückten Unterlisten.
- `verifizierbar`: ja — der Paritäts-Beleg aus
  [slice-100](../plan/planning/open/slice-099-structure-modul.md)
  §3 gegen eine Anforderungsdatei dieses Repos schlägt fehl, solange die
  Verankerung Listen-Marker nicht kennt.
- `klasse`: Zeilenanker ignoriert Listen-Marker

### F-4 — Die Kandidaten-Menge des Post-Passes ist nicht definiert

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 2; [`DC-FA-SCAN-001`](../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln)
- `pfad`: `spec/spezifikation.md:1732`
- `befund`: Schritt 2 sagt „die Dateien, deren Wurzel-relativer Pfad auf `files`
  matcht" — aber nicht, aus **welcher** Menge. Der Text stellt das Modul
  ausdrücklich neben `planning` („scannt nicht den Markdown-Baum"), und
  `planning` listet genau ein deklariertes Verzeichnis über den
  Filesystem-Port; ein Glob dagegen braucht einen Baumlauf. Offen bleibt
  damit: gelten `scan.roots` und `scan.ignore`, gelten die `SKIP_DIRS`, und
  sind Nicht-Markdown-Dateien Kandidaten (der einzige vorhandene Baumlauf,
  `DiscoverFiles`, filtert auf `.md` und prunt an `scan.ignore`). Beide
  Auslegungen sind vertragskonform und liefern verschiedene Befundsätze; die
  Scan-Menge-Auslegung nimmt zusätzlich still Dateien aus einer Regel heraus,
  die deren Glob trifft und `scan.ignore` verdeckt — die Regel meldet dann
  Erfolg über eine geschrumpfte Menge. Ein `files`-Glob, der ausschließlich
  Nicht-Markdown trifft, ist unter der einen Auslegung „null Kandidaten"
  (`section-missing`), unter der anderen eine Reihe von Binärdateien, die
  zeilenweise gelesen werden.
- `verifizierbar`: ja — Fixture mit einer Regel, deren Glob eine unter
  `scan.ignore` stehende Datei trifft; Befund vorhanden oder nicht entscheidet
  die Auslegung.
- `klasse`: Post-Pass ohne definierte Eingabemenge

### F-5 — Die Preset-Kopplung ist nur halb nachgezogen: Schema-Zeile und Lastenheft führen weiter den „ersten Treffer"

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in);
  [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  Entscheidung 4 („eine Änderung an nur einer Stelle ist ein Spezifikations-Bug")
- `pfad`: `spec/spezifikation.md:2188`, `spec/lastenheft.md:1876`
- `befund`: Der neue Schritt C3 bricht bei mehr als einem Treffer mit
  `closure-note-ambiguous` ab. Die §2-Schema-Zeile zu
  `planning.closure.heading-pattern` sagt unverändert „der **erste** Treffer
  eröffnet den geprüften Abschnitt", und der Beschreibungstext der Anforderung
  sagt drei Absätze vor der neuen Mehrdeutigkeits-Regel unverändert „in ihr der
  **erste** Abschnitt, dessen Überschrift … passt". Damit stehen in derselben
  Änderung beide Semantiken nebeneinander; der Implementer, der der
  Schema-Tabelle folgt (die für Config-Schlüssel die feinste Auflösung hat),
  liefert genau das stille Verhalten, das die Änderung beenden soll.
- `verifizierbar`: ja — ein Fixture mit zwei passenden Überschriften; nur eine
  der beiden Fassungen liefert `closure-note-ambiguous`.
- `klasse`: Schema-Zeile nicht mit dem Algorithmus mitgeführt

### F-6 — Vergleichsgegenstand von `section` ist an drei Stellen verschieden beschrieben

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
- `pfad`: `spec/lastenheft.md:1970`, `spec/spezifikation.md:1741`, `spec/spezifikation.md:2192`
- `befund`: Das Lastenheft nennt `section` den „Heading-Klartext, exakt", die
  Algorithmus-Sektion einen „Klartext-Vergleich der getrimmten
  Überschriften-Zeile", die Schema-Zeile einen „getrimmten Zeilen-Vergleich".
  Ob eine Regel `section: Akzeptanzkriterien` also auf eine H2 mit diesem Text
  passt oder ob der Rautenpräfix mitgeschrieben werden muss (und die Regel
  damit die Überschriftsebene festnagelt), ist nicht entschieden. Das
  Happy-Path-Akzeptanzkriterium lässt sich nicht als Test schreiben, ohne die
  Antwort zu raten; beide Auslegungen liefern für dieselbe Konfiguration
  entweder einen Treffer oder `section-missing`.
- `verifizierbar`: ja — der Happy-Path-Test entscheidet sich beim ersten
  Fixture; er ist heute mit beiden Fassungen begründbar.
- `klasse`: Vergleichsgegenstand eines Klartext-Schlüssels nicht festgelegt

### F-7 — Ankersemantik von `forbid-pattern` über den mehrzeiligen Abschnitts-Text ist offen

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 5
- `pfad`: `spec/spezifikation.md:1760`, `spec/spezifikation.md:2197`
- `befund`: `section-pattern` wird gegen **eine Zeile** geprüft, `forbid-pattern`
  gegen den **gesamten bereinigten Abschnitts-Text**. Ob dabei der
  Mehrzeilen-Modus gilt, sagt der Vertrag nicht; ohne ihn binden Zeilenanfang
  und Zeilenende in Go-RE2 an Anfang und Ende des ganzen Textes, und der Punkt
  überquert keinen Zeilenumbruch. Ein Muster, das nach dem Vorbild von
  `section-pattern` zeilenverankert geschrieben wird (etwa „Zeile beginnt mit
  offener Checkbox"), trifft dann nie — die Bedingung läuft still grün, obwohl
  der verbotene Text im Abschnitt steht. Zwei Anker-Regime derselben
  Konfigurations-Oberfläche ohne einen Satz, der sie unterscheidet.
- `verifizierbar`: ja — Fixture mit zeilenverankertem `forbid-pattern` und
  passendem Text ab Zeile 2 des Abschnitts.
- `klasse`: Zwei Anker-Regime für RE2-Schlüssel derselben Oberfläche

### F-8 — Das veröffentlichte §2-Schema ist größer als der Kern-Slice, und die Config-Dekodierung ist strikt

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-CONF-001`](../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei);
  [slice-099](../plan/planning/open/slice-099-structure-modul.md) §1
- `pfad`: `spec/spezifikation.md:2195` gegen
  `internal/adapter/driven/configyaml/configyaml.go:1400`
- `befund`: Die Schema-Tabelle führt seit diesem Commit `structure[].max-tasks`,
  `require-all` und `require-any`; geliefert werden sie erst mit
  [slice-100](../plan/planning/open/slice-099-structure-modul.md).
  Die Konfiguration wird mit `KnownFields(true)` dekodiert — ein Schlüssel, den
  die Struktur nicht kennt, ist Exit 2. Ein Release nach dem Kern-Slice weist
  damit eine **spezifikationskonforme** Konfiguration mit Exit 2 ab; der
  aktuelle Stand zeigt die Fehlerform bereits für den ganzen Block
  (`field structure not found`, Exit 2). Umgekehrt wäre die Alternative,
  die Schlüssel im Kern-Slice als akzeptiert-aber-wirkungslos aufzunehmen — ein
  stiller Grün-Pfad. Der Slice-Schnitt trägt in dieser Reihenfolge nur, wenn
  zwischen beiden Slices kein Release liegt; das steht in keinem der beiden
  Slices.
- `verifizierbar`: ja — Konfiguration mit `structure[].max-tasks` gegen ein
  Image, das nur den Kern trägt.
- `klasse`: Schema veröffentlicht vor dem Slice, der es implementiert

### F-9 — „nicht gesetzt" gegen „auf 0 gesetzt" ist nur für die Closure-Fähigkeit ausformuliert

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 1
- `pfad`: `spec/spezifikation.md:2196`, `spec/spezifikation.md:1730`
- `befund`: Für die Closure-Fähigkeit hält der Vertrag ausdrücklich fest, dass
  ein abwesendes `min-sentences` der Default und kein Fehler ist — „die
  Unterscheidung ist Teil der Zusage". Für `structure[]` fehlt dieser Satz,
  obwohl beide neuen Zahl-Schlüssel ihn brauchen und `max-tasks` ihn in die
  **andere** Richtung braucht: `max-tasks: 0` ist ein gültiger und der
  nützlichste Wert („in diesem Abschnitt steht kein offener Task"), während
  abwesend „aus" heißt. Wird der Schlüssel als schlichte Ganzzahl umgesetzt,
  ist der Nullwert vom Default nicht unterscheidbar und die schärfste Setzung
  der Bedingung schaltet sie ab — ohne Meldung. Für `min-sentences` widerspricht
  zusätzlich Schritt 1 („`min-sentences` < 1" ⇒ Exit 2) der Schema-Zeile
  („explizit < 1"), die den Default 0 zulässt.
- `verifizierbar`: ja — Fixture mit `max-tasks: 0` und einem Task-Item im
  Abschnitt; Befund oder Exit 0 entscheidet.
- `klasse`: Nullwert nicht von „nicht gesetzt" unterschieden

### F-10 — `DC-FA-CONF-002` behauptet Universalität, die die Listen-Form von `structure` nicht erfüllen kann

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-CONF-002`](../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope)
- `pfad`: `spec/lastenheft.md:2313`, `spec/spezifikation.md:2191`
- `befund`: Die Anforderung sagt „**Jedes** Regelmodul akzeptiert in der
  Konfigurationsdatei optional einen Schlüssel `<modul>.scope`". `structure`
  ist als bare Liste spezifiziert (`structure[].files` …); ein Unterschlüssel
  `scope` ist damit strukturell unmöglich. Das ist die zweite stille Ausnahme
  neben `sources` (dort steht sie immerhin im Code-Kommentar), und sie ist
  hier nicht folgenlos: `scope` wäre der vorhandene, vertraglich geregelte
  Mechanismus gewesen, um die in F-4 offene Eingabemenge zu benennen. Weder
  die neue Anforderung noch die ADR erwähnen die Kollision.
- `verifizierbar`: ja — eine Konfiguration mit `structure.scope` ist nach der
  Schema-Form ein Syntaxfehler, nach `DC-FA-CONF-002` ein gültiger Schlüssel.
- `klasse`: Universelle Config-Zusage bekommt stille Ausnahme

### F-11 — Kein Slice trägt die CLI-Mit-Modifikation, die jedes bisherige Modul mitgebracht hat

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-CLI-005`](../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
  („dokumentiert die verfügbaren Module und Optionen als Kommentare");
  [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
- `pfad`: `docs/plan/planning/open/slice-099-structure-modul-kern.md:37`,
  `docs/plan/planning/open/slice-100-structure-marken-und-zaehlung.md:38`
- `befund`: Die Historien-Einträge zu `planning`, `targets` und `vcs` führen die
  Erweiterung von `--print-config`, `--suggest-config` und `--print-mk`
  jeweils ausdrücklich als Teil der Modul-Einführung. Der 0.51.0-Eintrag nennt
  sie nicht, und keine der beiden Definitions of Done erwähnt sie — das Modul
  wäre nach beiden Slices in keiner Ausgabe des Werkzeugs sichtbar. Dass genau
  diese Enumerationen driften, ist belegt: die heutige „Verfügbar"-Zeile von
  `--print-config` listet 18 Module und führt `citations` nicht, obwohl es das
  18. Modul ist.
- `verifizierbar`: ja — `--print-config` nach Abschluss beider Slices gegen die
  Modul-Liste aus
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl).
- `klasse`: Modul-Enumerationen in der CLI-Oberfläche nicht mitgeführt

### F-12 — Task-Item-Lexik ist per Beispiel statt per Regel festgelegt

- `kategorie`: LOW
- `quelle`: [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 5
- `pfad`: `spec/spezifikation.md:1757`, `spec/lastenheft.md:1981`
- `befund`: Der Vertrag nennt „Listen-Zeile mit offener oder gesetzter
  Checkbox" bzw. zeigt die Bindestrich-Form. Nicht festgelegt sind: die
  übrigen GFM-Listen-Marker (Stern, Plus, geordnete Listen), die
  Großschreibung des Häkchens, die erlaubte Einrückung und ob eine
  Tabellenzeile mit einer Checkbox zählt. Eine Checkliste, die mit einem Stern
  statt einem Bindestrich geschrieben ist, entginge `max-tasks` — die
  Bedingung meldet nichts, obwohl die Items im Abschnitt stehen. Die
  Markdown-Lexik dieses Repos wird sonst in
  [ADR-0042](../plan/adr/0042-markdown-lexik-folgt-commonmark.md) an GFM
  gebunden, gerade weil handgeschriebene Heuristik still weniger meldet.
- `verifizierbar`: ja — Fixture mit Stern-Listen-Checkboxen und `max-tasks: 0`.
- `klasse`: Markdown-Lexik per Beispiel statt per Regel

### F-13 — Reichweite des Mehrdeutigkeits-Abbruchs ist zwischen Regel und Datei nicht getrennt

- `kategorie`: LOW
- `quelle`: [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 3
- `pfad`: `spec/spezifikation.md:1744`
- `befund`: Der Vertrag sagt „Abbruch für diese **Datei**". Anders als bei der
  Closure-Fähigkeit, wo genau eine Bedingungsmenge je Datei gilt, kann eine
  Datei bei `structure` von mehreren Regeln getroffen werden. Ob eine
  mehrdeutige Überschrift unter Regel A auch die Messung unter Regel B
  abschaltet, ist offen; die abschaltende Lesart nimmt eine erfüllbare Prüfung
  still aus dem Lauf.
- `verifizierbar`: ja — Fixture mit zwei Regeln auf dieselbe Datei, davon eine
  mit mehrdeutiger Abschnitts-Bestimmung.
- `klasse`: Abbruch-Skopus zwischen Regel und Datei nicht getrennt

### F-14 — Die Config-Rand-Liste steht in drei Fassungen, und die Glob-Validierung fehlt

- `kategorie`: LOW
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  §fail-closed
- `pfad`: `spec/lastenheft.md:2009`, `spec/spezifikation.md:1728`
- `befund`: Die Exit-2-Liste des Lastenhefts nennt den leeren Eintrag nur für
  `require-all`/`require-any`, die der Spezifikation zusätzlich für
  `exempt-paths`. In keiner der beiden steht, dass ein ungültiges
  `files`- oder `exempt-paths`-Glob Exit 2 ist — obwohl `planning.slice-glob`
  genau dafür eine eigene Zeile hat („verhindert ein fail-open Silent-Grün").
  Ein ungültiges Glob liefert im vorhandenen Matcher schlicht keinen Treffer;
  bei `exempt-paths` heißt das „Ventil wirkt nicht", bei `files` „Regel läuft
  leer" — beides ohne Hinweis auf die Ursache.
- `verifizierbar`: ja — Konfiguration mit unbalancierter Klammer im Glob;
  Exit 2 oder Befund entscheidet.
- `klasse`: Config-Rand-Liste an drei Stellen mit drei Fassungen

### F-15 — Der Nullmengen-Guard endet an der Modulgrenze

- `kategorie`: INFO
- `quelle`: [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  §Fitness Function („Leerlauf" gegen „Byte-Identität")
- `pfad`: `spec/spezifikation.md:1727`, `spec/lastenheft.md:2020`
- `befund`: Eine Regel ohne Treffer meldet `section-missing` mit der
  ausdrücklichen Begründung, ein leer laufendes Gate dürfe nicht Erfolg melden.
  Ein **aktiviertes Modul ohne Regeln** ist dagegen inert und grün. Die
  Aktivierung ist bei `planning.closure.dir` als „Behauptung, dass dort etwas
  liegt" gewertet worden; `--enable structure` ist dieselbe Behauptung eine
  Ebene höher. Der Fall ist real erreichbar, seit die Konfiguration per
  `--config` aus einer eigenen Datei kommen kann: ein Profil, dem der
  `structure`-Block fehlt, läuft grün durch und sieht wie ein bestandener
  Lauf aus. Bewusst dokumentierte Annahme, kein Widerspruch — aber die
  Begründungsketten der beiden Ebenen widersprechen einander.
- `verifizierbar`: ja — Lauf mit aktiviertem Modul und leerer Regel-Liste.
- `klasse`: Nullmengen-Guard endet an der Modulgrenze

## Negativbefunde

- geprüft, ohne Befund: **Preset-Kopplung Schritt für Schritt gegen den Code.**
  Abschnitts-Bestimmung (`closureHeadingLine`), Abschnitts-Grenze bei gleicher
  oder höherer Ebene (`closureSectionProse`), Fence-Bereinigung (`FenceToggle`)
  und Satzzählung (`countSentenceEnds`) in
  `internal/hexagon/core/rules/planning.go` fallen tatsächlich aus **einer**
  Mechanik; sie nutzen bereits den geteilten ATX-Parser aus
  `internal/hexagon/core/rules/markdown.go`. Zum Preset fehlt genau eine
  Erweiterung — die Trefferzählung statt des Abbruchs beim ersten Treffer.
  Die Kandidaten-Bestimmung (C2 gegen Schritt 2) ist verschieden, wird aber
  vom Kopplungs-Satz ausdrücklich nicht beansprucht.
- geprüft, ohne Befund: **RE2-Fallen im engeren Sinn.** Weder `section-pattern`
  noch `forbid-pattern` verlangen laut Text Lookaround, Rückwärtsreferenzen
  oder gierige Rückverfolgung; beide sind reine Treffer-Prädikate, und die
  Marken sind literal statt regulär. Die Voreinstellung der Closure-Fähigkeit
  ist RE2-fähig. Offen ist nur die Ankersemantik (F-7), nicht die
  Ausdrucksstärke.
- geprüft, ohne Befund: **Netzlos-Modulliste.** Die Messmethode von
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  ist als „alle Module außer `external`, `sources` und `vcs`" formuliert;
  `structure` fällt als hermetisches Modul automatisch hinein, eine
  Vertrags-Ergänzung ist nicht geschuldet. Die repo-lokale
  Modul-Listen-Selbstkonsistenz ist laut
  [`DC-FA-TGT-001`](../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
  §Out-of-Scope ohnehin repo-lokal.
- geprüft, ohne Befund: **Determinismus der Sortierung.** Die
  Sortier-Schlüssel aus `internal/hexagon/core/model/finding.go` tragen das
  Ziel als Tiebreak; zwei Regeln mit verschiedenen `files`-Globs auf derselben
  Datei ordnen deterministisch. Die Kollision entsteht erst bei identischem
  Ziel (F-2).
- geprüft, ohne Befund: **CRLF und fehlender Schluss-Zeilenumbruch.**
  `splitLines` trennt an Zeilenvorschüben, der ATX-Parser und die
  Fence-Erkennung trimmen den Wagenrücklauf weg, und die Abschnitts-Grenze
  läuft bis zum Dateiende. Für eine Marken-Prüfung am Zeilenanfang und eine
  Zeichen-Zählung ist ein anhängender Wagenrücklauf folgenlos.
- geprüft, ohne Befund: **Überschrift in Zeile 1 und Abschnitt am Dateiende.**
  Beide Fälle sind von der Abschnitts-Bestimmung gedeckt; die Zeile 1 als
  Fundort einer echten Überschrift kollidiert nicht mit der Zeile 1 als
  Platzhalter für `section-missing`, weil die Grund-Codes verschieden sind.
- geprüft, ohne Befund: **Modul-Schnitt und Nicht-Supersede.**
  [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) wendet das
  Schnitt-Kriterium aus [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md)
  konsistent an; die Begründung, die stabil zugesagten Grund-Codes nicht zu
  tauschen, deckt sich mit der Codes-Tabelle in
  `spec/spezifikation.md` §4, und der Verzicht auf einen Alias-Slice ist in der
  Welle nachgezogen.
- geprüft, ohne Befund: **Diagnose-only, Opt-in, Hermetik.** `structure`
  verlangt keinen neuen Port; alle Zusagen (kein Reparatur-Hunk, kein git,
  kein Netz, Byte-Identität ohne Aktivierung) sind mit dem vorhandenen
  Filesystem-Port und dem Post-Pass-Muster aus
  `internal/hexagon/core/rules/run.go` erfüllbar.
- geprüft, ohne Befund: **Slice-Reihenfolge.** Der Kern vor den Marken ist
  sachlich richtig — Marken und Zählung setzen die Abschnitts-Bestimmung
  voraus, umgekehrt nicht. Der Kern-Slice ist groß (Modul, Schema, vier
  Grund-Codes, Refactor ausgelieferten Codes), aber die Teile sind gekoppelt:
  das Herausziehen der gemeinsamen Mechanik und `closure-note-ambiguous` sind
  dieselbe Änderung. Der Einwand gegen den Schnitt ist nicht die Größe,
  sondern die Release-Grenze (F-8).
- geprüft, ohne Befund: **`section-constraint` als Sammel-Code, Begründungs-Ebene.**
  Der Widerspruch zur Drei-Codes-Begründung aus
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md) ist
  in [slice-099](../plan/planning/open/slice-099-structure-modul.md) §4
  als offener Punkt benannt und damit nicht still. Bewertbar ist er erst mit
  F-2: die Bündelung ist nicht nur eine Konsumenten-Frage, sondern schlägt in
  der Befund-Menge durch.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 3 |
| MEDIUM | 8 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Naiver Fence-Toggle macht Verbots-Bedingung
still grün · Sammel-Grund-Code kollidiert mit Befund-Dedup · Zeilenanker
ignoriert Listen-Marker · Post-Pass ohne definierte Eingabemenge · Schema-Zeile
nicht mit dem Algorithmus mitgeführt · Vergleichsgegenstand eines
Klartext-Schlüssels nicht festgelegt · Zwei Anker-Regime für RE2-Schlüssel
derselben Oberfläche · Schema veröffentlicht vor dem Slice, der es
implementiert · Nullwert nicht von „nicht gesetzt" unterschieden · Universelle
Config-Zusage bekommt stille Ausnahme · Modul-Enumerationen in der
CLI-Oberfläche nicht mitgeführt · Markdown-Lexik per Beispiel statt per Regel ·
Abbruch-Skopus zwischen Regel und Datei nicht getrennt · Config-Rand-Liste an
drei Stellen mit drei Fassungen · Nullmengen-Guard endet an der Modulgrenze

## Verdikt

**Merge-blockierend:** ja — drei HIGH und acht MEDIUM.

Das Modul ist grundsätzlich baubar: die Preset-Behauptung trägt, der
vorhandene Code liefert Abschnitts-Bestimmung, Fence-Bereinigung und Zählung
bereits aus einer Mechanik, und der einzige echte Zusatz zur Kopplung ist die
Trefferzählung. Die Einwände liegen woanders. F-1 und F-2 sind belegte
Mechanik-Defekte: der eine macht Verbots-Bedingungen bei kaputter
Fence-Balance still grün, der andere macht die ausdrückliche Zusage „mehrere
verletzte Bedingungen ⇒ mehrere Befunde" gegen die globale Deduplikation
unerreichbar. F-3 trifft die Motiv-Dokumentklasse: die zugesagte
Marken-Verankerung passt auf keine einzige Akzeptanz-Zeile dieses Repos.
F-4 bis F-7 und F-9 sind Stellen, an denen zwei vertragskonforme
Implementierungen verschiedene Befundsätze liefern — die
Akzeptanzkriterien zu Happy Path, `forbid-pattern` und `max-tasks` sind ohne
diese Entscheidungen nicht als Test schreibbar.

**Übergabe:** Findings gehen an den Implementer; die Rückkante Review → Plan
ist hier die Regel, nicht die Ausnahme — F-4, F-6, F-7 und F-9 sind
Vertragsentscheidungen, keine Codefragen, und gehören vor
[slice-099](../plan/planning/open/slice-099-structure-modul.md) in die
Spezifikation. Die Finding-Klassen gehen zusätzlich in die Slice-Closure §7
und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg** und ersetzt
keine Verifikation.
