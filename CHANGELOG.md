# Changelog

Alle nennenswerten Änderungen an diesem Projekt werden in dieser Datei
dokumentiert. Das Format folgt [Keep a Changelog](https://keepachangelog.com/de/1.1.0/),
die Versionierung folgt [SemVer](https://semver.org/lang/de/).

## [Unreleased]

### Added

- slice-095 — **Ortsfeste Verweise** (`links.resolve-from`, opt-in;
  `DC-FA-LINK-001` §Ortsfeste Verweise,
  [ADR-0056](docs/plan/adr/0056-resolve-from-wandernde-quellorte.md); Change
  Request des Konsumenten a-check). Wo Dateien per `git mv` zwischen
  Geschwister-Verzeichnissen **wandern** — der Planning-Lifecycle ist der
  Anlassfall —, muss ein relativer Verweis von **jedem** Ort der Gruppe
  auflösen, und überall auf **dasselbe** Ziel. Sonst
  `link-position-dependent`: die Reparatur ist das **Präfixieren des Pfads**,
  nicht das Anlegen des Ziels — am Ist-Ort ist nichts kaputt. Dateien in
  `fixed-dirs` (etwa dem Ruheort) sind hypothetische Ziele, aber keine
  Quellen; ein Ziel, das schon am Ist-Ort fehlt, meldet weiter nur
  `target-missing` (kein Doppelbefund). Ohne den Block byte-identisch.
  **Der Beleg lief mit dem Produkt:** der heutige Bestand ist grün (die Null
  ist von einer Woche Hand-Nachzügen gekauft — allein am Vortag 10, 15 und 14
  nachgezogene Verweise je Move), und der Retro-Lauf gegen den Stand vor der
  welle-69-Eröffnung meldet exakt die **19** des realen Bruchs, mit
  identischer Verteilung. Laufzeit unverändert (im Rauschen, `DC-QA-01`).

## [0.59.0] — 2026-08-16

### Added

- slice-102 — **Wellen-Invariante: die Roadmap-Aussage gegen das Verzeichnis**
  (`DC-FA-PLAN-001` §Wellen-Invariante,
  [ADR-0055](docs/plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md)).
  Dieselbe Lifecycle-Invariante, die `planning` auf der Slice-Ebene prüft, eine
  Ebene höher: der Aktiv-Status-Abschnitt nennt eine Welle genau dann, wenn
  **genau ein** flaches Wellendokument existiert (`wave-drift`); eine
  Vorschau-Zeile nennt keine Welle, die schon eine Datei hat
  (`wave-preview-exists`); jede Zeile des Abschluss-Registers hat ihre
  **Ergebnisnotiz** (`wave-results-missing`) — **und jede Ergebnisnotiz ihre
  Zeile** (`wave-unregistered`). Vier Grund-Codes für vier Reparaturen.
  **Opt-in innerhalb** des opt-in Moduls über `planning.waves.dir`; ohne den
  Schlüssel wird kein Wellendokument geöffnet und der Befundsatz ist
  byte-identisch.
  **Zwei Festlegungen sind gemessen, nicht gesetzt:** das verpflichtende
  Artefakt einer geschlossenen Welle ist die **Ergebnisnotiz** — gegen das
  Plan-Dokument geprüft meldet die Aussage über zwei reale Planungs-Bäume
  19-mal, weil ältere Wellen vor dieser Konvention geschlossen wurden. Und die
  Vorschau-Aussage liest die **erste Spalte** und überspringt Zeilen ohne
  Kennung: geplante Wellen tragen Namen, und die Trigger-Spalte darf andere
  Wellen nennen. **Beim Einführen mit Rückständen rechnen** — im eigenen Baum
  0 Befunde, im Schwester-Repo elf echte (fehlende Ergebnisnotizen); ein
  konsument-gerechter `planning.marker` gehört zur Einführung dazu, sonst
  kommt ein Artefakt-Befund hinzu.

## [0.58.0] — 2026-08-16

### Fixed

- slice-103 — **Geteilte Lexik: drei Konsumenten beantworteten eine Lexik-Frage
  selbst.** `citations` bildete seinen Absatz allein an Leerzeilen — ein
  Fenced-Block zwischen Direktive und Zitat trennte dort **nicht**, und dieselbe
  Datei lieferte je nach Trennzeichen zwei entgegengesetzte Ergebnisse: mit
  Block **Exit 0** (oder ein falsches `citation-mismatch`), ohne Block den
  zugesagten fail-closed-Abbruch. Die Anker-Auflösung von `versions.current-from`
  und `pins` erkannte HTML-Anker **auch innerhalb** von Fences, während `anchors`
  denselben Anker im selben Lauf für nicht existent hielt. Beides folgt jetzt der
  geteilten Antwort. **Die Änderung wirkt in beide Richtungen.** Sie **findet
  mehr**: eine bisher grüne Direktive hinter einem Code-Block bricht fail-closed
  ab (Exit 2), und eine Zeichenfolge, die nur wie ein Anker aussieht — in
  Inline-Code, als `data-id`, als `name` an einem beliebigen Element oder ganz
  ohne Tag —, löst nicht mehr auf. Sie **findet weniger**, unter anderem hier: ein
  von einem Fenced-Block unterbrochenes `>`-Zitat wird gekürzt gelesen und meldet
  sein `citation-mismatch` nicht mehr; ein `dpin`, dessen Ziel-Anker nur in einem
  Fence steht, verliert seinen Drift-Schutz **kommentarlos** (`pins` schweigt zu
  unauflösbaren Zielen); ein **anders geschriebener** Anker trifft nicht mehr, weil
  die Auflösung jetzt case-sensitiv ist wie in `anchors` (`#aktuell` findet
  `id="Aktuell"` nicht mehr — bei `versions` als Exit 2, bei `pins` wieder still);
  und ein `-1`-Anker adressiert bei konkurrierenden Kennungen einen **anderen
  Span**. **Diese Aufzählung ist offen** — sie nennt die gemessenen Fälle, nicht
  alle möglichen; in drei Review-Runden ist sie dreimal unvollständig gewesen. Am
  eigenen Bestand gemessen ist **kein** heutiger Fall betroffen. Der rohe Ziel-Span von `pins` und die Fence-Ausnahme von `versions`
  bleiben unberührt — das sind andere **Fragen**, keine anderen **Antworten**.
- slice-103 — **Die Anker-Frage hat jetzt EINE Antwort, auch in der
  Slug-Hälfte:** `versions.current-from` und `pins` lösen Duplikat-Slugs
  (`#alt-1`) und prozent-kodierte Fragmente (`#a%20b`) jetzt so auf wie
  `anchors` — vorher brach derselbe Lauf mit Exit 2 ab, während `anchors`
  denselben Anker für gültig hielt. **Findet mehr.**
- slice-103 — **Modul `planning`: zwei Falsch-Rot weniger.** Der
  Aktiv-Status-Guard las Rohzeilen: eine Roadmap, die ihren eigenen Abschnitt in
  einem Beispiel-Block zeigt, galt als „Überschrift mehrfach vorhanden“, und eine
  Raute-Zeile in einem Beispiel-Block beendete den Aktiv-Block vorzeitig. Beides
  meldete `planning-drift`, wo nichts driftete. **Findet weniger** — und zwar
  Fehlmessungen. In der Gegenrichtung **findet es mehr**: der Aktiv-Block endet
  jetzt an der geteilten Abschnittsgrenze statt an einer rohen `## `-Prüfung —
  eine eingerückte H2, eine tab-getrennte H2 und eine H1 beenden ihn ebenfalls,
  ein Ruhe-Marker dahinter zählt nicht mehr zum Block, und `planning-drift`
  entfiel dort bisher **still**.
- slice-103 — **Modul `vcs`: der Immutabilitäts-Kopf wird nicht mehr aus einem
  Code-Beispiel gelesen.** Status-Zeile und `immutable-when` liefen über die
  **rohen** Zeilen, während die Abschnitts-Maske derselben Funktion
  fence-bewusst ist. Eine Datei, die ihren eigenen Kopf als Beispiel zeigt —
  Vorlagen und Konventionsdateien tun das —, galt dadurch als immutabel oder
  verschob die gestrippte Status-Zeile, und eine reale Core-Änderung passierte
  das Gate **mit Exit 0 und ohne Ausgabe**. **Findet mehr**, und schließt ein
  stilles Grün in einem Gate. Gemessen: 0 von 54 eigenen ADRs betroffen.
- slice-103 — **Modul `targets`: eine Tabellenzeile im Code-Block dokumentiert
  nichts mehr.** Ein Doku-Beispiel, das eine Ziel-Tabelle zeigt, ließ ein
  undokumentiertes Target als dokumentiert gelten; `gate-undocumented` entfiel
  still. **Findet mehr.**
- slice-103 — **Benannte Grenze statt zweitem Wächter:** die Abschnitts-Maske
  des Moduls `vcs` wird auf git-Blobs gerechnet, die kein scannendes Modul sieht;
  ein unbalancierter Fence in einer Revision verschiebt sie, und für bereits
  existierende Commits ist die Sicht rückwirkend blind. **Die gefährliche
  Richtung ist das stille Grün:** liegt die Öffnung im ausgenommenen Abschnitt,
  läuft die Ausnahme bis zum Dateiende und eine reale Core-Änderung passiert mit
  Exit 0 ohne Ausgabe. Gemessen: **0 von 152** immutablen Revisions-Blobs. Die
  Grenze steht jetzt im Vertrag statt in einem Review-Report, mit einem Trigger,
  der sich beobachten lässt (`fence-unclosed` in einer `vcs.paths`-Datei).

- slice-099 (Closure) — **Modul `structure`, unlesbarer Dateibaum:** scheiterte
  der Baum-Aufbau, meldete das Modul **nichts** — ununterscheidbar von „alle
  Regeln erfüllt", und das als einziger stiller Pfad eines Moduls, dessen
  sämtliche anderen Randbedingungen fail-closed melden. Jetzt meldet **jede**
  Regel `section-missing` mit ihrer Identität. Der Zweig ist über die CLI heute
  nicht erreichbar (der Scanner öffnet die Wurzel vorher); er wird es, sobald
  ein Filesystem-Port Teil-Lesefehler durchreicht.

### Documentation

- slice-099 (Closure) — vier Vertragsflächen nachgezogen: die **Symlink-Grenze**
  gilt ausdrücklich auch für **Datei**-Symlinks (bisher stand nur der
  Verzeichnis-Symlink in §1), der Wort-Begriff der **Marke** ist als
  unicode-weit ausgewiesen und gegen die ASCII-Wortgrenze der Floskel
  abgegrenzt, die Schema-Zeile zu `structure[].min-sentences` trägt die
  Zähl-Semantik (Inline-Code geleert, Satzende nur vor Whitespace/Zeilenende),
  und die Modul-Tabelle des Handbuchs führt `closure-note-ambiguous`.

## [0.57.0] — 2026-08-15

### Added

- slice-099 — **Modul `structure`**, das 20. Regelmodul (`DC-FA-STRUCT-001`,
  [ADR-0049](docs/plan/adr/0049-structure-modul-schnitt-und-preset.md)):
  Struktur-Invarianten **innerhalb** eines Dokuments. Je Regel eine
  Dokumentklasse über **eigene** Globs (unabhängig von `scan.roots`/`scan.ignore`
  — deshalb kein `<modul>.scope`), ein Abschnitt (Klartext **oder** RE2) und bis
  zu sechs Bedingungen mit **je eigenem** Grund-Code: `section-empty`,
  `section-thin`, `section-oversized`, `section-forbidden`,
  `section-pattern-missing`, `section-marker-missing`. Dazu zwei Struktur-Codes:
  `section-missing` (kein passender Abschnitt — **oder** die Regel trifft keine
  Datei, auch nach Abzug von `exempt-paths`) und `section-ambiguous` (mehrere
  Treffer bei `sections: one`). `sections: each` prüft jeden Treffer einzeln.
  Hermetisch, diagnose-only, opt-in; ohne Regeln wird keine Datei geöffnet.
- slice-099 — **`closure-note-ambiguous`**: mehrere Closure-Notiz-Überschriften
  in einem Dokument. Die Spezifikation sagte die Härte seit v0.51.0 zu, ohne dass
  sie implementiert war — ohne eindeutigen Abschnitt wird jetzt **nicht**
  gemessen, statt still den ersten zu nehmen. **Verhaltensänderung für
  Konsumenten:** ein Dokument mit zwei passenden Überschriften meldete bisher
  `closure-note-thin`/`-boilerplate`/`-placeholder` über den **ersten**
  Abschnitt; jetzt meldet es `closure-note-ambiguous` an der Zeile des
  **zweiten** und wird nicht gemessen. **Die Richtung schließt „vorher grün,
  jetzt rot" ein:** war der erste Abschnitt substanzhaltig, meldete der Lauf
  bisher **nichts** — jetzt meldet er die Mehrdeutigkeit. Wer zwei passende
  Überschriften in einem Dokument führt (etwa einen stehengebliebenen
  Vorlagen-Abschnitt), bekommt beim Update einen neuen Befund.
- slice-099 — **zwölftes `--print-mk`-Target `doc-structure`**
  (`DC-FA-CLI-010`, 11 → 12).

### Changed

- slice-099 — die **Abschnitts-Mechanik ist geteilt**: die Closure-Note-Struktur
  des Moduls `planning` ist ein **Preset** derselben Semantik und nutzt dieselbe
  Abschnitts-Findung, Kardinalitäts-Behandlung und Bereinigung. Für Konsumenten
  ändert sich daran nichts — ein Kopplungs-Test fährt dieselbe Eingabe durch
  beide Oberflächen und vergleicht Zeile **und** gemessene Zahl.

## [0.56.0] — 2026-08-10

### Changed

**Diese Version ändert die Semantik eines ausgelieferten Gates in drei
Richtungen.** Wer `planning.closure` nutzt, sollte alle drei lesen; ohne
`closure.dir` ist nichts davon aktiv.

- slice-094 — **`closure-note-thin` findet mehr.** Der Closure-Abschnitt wird
  jetzt **einmal** bereinigt (Fenced-Code **und** Inline-Code), und ein
  Satzende-Zeichen zählt nur, wenn ihm Whitespace oder das Zeilenende folgt
  (`DC-FA-PLAN-001`,
  [ADR-0053](docs/plan/adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md)).
  Anlass ist die **Parität** zum handgeschriebenen Prüfskript, das dieses Modul
  ablösen soll: eine Notiz, die dessen Sensor als zu dünn meldete, lief bei
  d-check durch. Gemessen an einem realen Bestand tragen mehr als die Hälfte
  aller Satzende-Vorkommen keinen Satz — Punkte aus Link-Pfaden und
  Versionsnummern. **Ein grüner Lauf kann rot werden.** Belegt: über die 84
  Closure-Notizen des Schwester-Repos zählen Skript und Modul identisch.
- slice-094 — **`closure-note-boilerplate` findet weniger.** Beide Bedingungen
  lesen denselben bereinigten Text; eine Floskel **in Backticks** trifft damit
  nicht mehr. Das ist beabsichtigt — eine *zitierte* Floskel ist keine benutzte
  —, aber es ist eine **Lockerung**: wer eine Floskel zitiert stehen hat,
  verliert einen bestehenden Befund. Sie reicht weiter als der Wortlaut nahelegt:
  geleert wird **jeder** Inline-Code-Span nach CommonMark-Paarung, auch ein
  unbeabsichtigter aus zwei einzelnen Backticks im selben Absatz.
- slice-104 — **`closure-note-boilerplate` vergleicht an Wortgrenzen** statt als
  Teilstring. Kurze Phrasen waren sonst unbrauchbar: `ok` traf auch in
  *dokumentiert* und *Protokoll* (gemessen: 68 Treffer, davon einer echt).
  Mehrwortige Phrasen sind verhaltensgleich; kurze werden erst dadurch
  benutzbar. Als Wortzeichen gilt die **ASCII**-Menge — ein Umlaut ist keines,
  eine Phrase mit angrenzendem Umlaut gilt also als grenzständig und trifft.
  **Auch das findet weniger.**

### Fixed

- slice-104 — **CRLF:** die Satzzählung akzeptierte das `\r` vor dem
  Zeilenumbruch nicht. In einer CRLF-Arbeitskopie zählte damit kein
  zeilenschließendes Satzende, und eine Notiz mit vier sauberen Sätzen meldete
  `closure-note-thin`. Nur in dieser Version relevant, weil die Zählung vorher
  jedes `.` nahm.

## [0.55.0] — 2026-08-10

### Added

- slice-098 — **`closure-note-placeholder`** als vierte, **opt-in** Bedingung der
  Closure-Note-Struktur (`DC-FA-PLAN-001`,
  [ADR-0052](docs/plan/adr/0052-platzhalter-erkennung-inline-code.md)):
  `planning.closure.placeholder: true` meldet den **unausgefüllten Rumpf** einer
  Vorlage. Anlass: ein Template-Rumpf ist syntaktisch vollständig und passiert
  alle drei bestehenden Bedingungen — Abschnitt vorhanden, Schwelle erreicht,
  keine Floskel. Erkannt wird die Auszeichnungs-Form (öffnende Winkelklammer
  ohne vorausgehendes Wortzeichen, Inneres **ohne Whitespace** — ein Feldname,
  kein Satz). Drei Nachfilter im Code statt im Muster: Autolink/Adresse,
  HTML-Tag-Namen und Markdown-Linkziele in Winkelklammern. **Inline-Code zählt
  nicht** — dort wird Syntax gezeigt; am eigenen Bestand gemessen liegen alle
  zwölf Treffer dort und keiner außerhalb. Gemeldet wird der **erste** Treffer je
  Kandidat, an seiner Zeile. Ohne den Schalter ist der Befundsatz
  byte-identisch.
  **Zwei Grenzen sind benannt:** ein **eingerückter** Code-Block (vier
  Leerzeichen) ist in d-check nirgends modelliert, und eine ungerade
  Backtick-Zahl im Absatz verschiebt die Inline-Code-Paarung — beides
  Eigenschaften der geteilten Markdown-Lexik, nicht dieser Bedingung.
  **Nicht** mit geändert: die Substanz-Zählung sieht Inline-Code weiterhin.

## [0.54.0] — 2026-08-10

### Added

- slice-097 — **`planning.closure.glob`** als eigener Kandidaten-Filter der
  Closure-Fähigkeit (`DC-FA-PLAN-001`,
  [ADR-0051](docs/plan/adr/0051-eigener-kandidaten-filter-closure.md)). Der
  Default ist ein **Verweis** auf `planning.slice-glob`, kein kopiertes Literal:
  ohne den Schlüssel ist der **Befundsatz** byte-identisch, und es gibt genau ein
  Muster zu pflegen. Anlass: die beiden `planning`-Fähigkeiten stellen
  verschiedene Fragen („liegt hier noch Arbeit?" gegen „ist jedes abgeschlossene
  Paket dokumentiert?") mit verschiedenen Grundmengen und teilten sich einen
  Schlüssel, solange beide zufällig dasselbe trafen; wer die eine Menge weitet,
  verbiegt die andere — im Grenzfall matcht die Roadmap-Datei selbst und der
  Ruhe-Marker meldet dauerhaft falsch-rot. Ein **explizit** leerer oder
  ungültiger Glob bricht mit Exit 2 ab statt still auf den Default
  zurückzufallen; ein YAML-`null` gilt als abwesend. Kein neuer Grund-Code — es
  ändert sich die Kandidaten-Menge, nicht die Aussage über einen Kandidaten.

### Changed

- slice-097 — zwei **Begleit-Texte** haben sich geändert; beide sind auch **ohne**
  den neuen Schlüssel erreichbar und ausdrücklich nicht stabilitätsgarantiert,
  aber wer auf Text greppt, sieht sie: die `message` der Nullmengen-Meldung führt
  den Glob jetzt in Anführungszeichen (sichtbar in `--json`), und der
  `--doctor`-Klartext zu `closure-note-missing` sagt **„Kandidat" statt „Slice"**
  — die Kandidaten-Menge muss seit diesem Release nicht mehr aus Slice-Dateien
  bestehen. Der Befundsatz (Datei, Zeile, Regel, Ziel, Grund) ist unverändert.

## [0.53.0] — 2026-08-10

### Added

- slice-101 — **`fence-unclosed`** als dritte Artefakt-Klasse des Moduls `spans`
  (`DC-FA-SPAN-001`,
  [ADR-0050](docs/plan/adr/0050-fence-unclosed-in-spans.md)): eine Fence-Öffnung
  ohne Schluss bis zum **Dateiende**. Anlass war ein **ausgelieferter stiller
  Grün-Pfad** — hinter einem offenen Fence übersprang die Vorverarbeitung den
  Rest, und Gates meldeten grün, ohne geprüft zu haben (reproduziert gegen das
  veröffentlichte Image v0.52.0). Ausgewertet werden **beide** Schluss-Lesarten
  des Produkts (naiver Toggle und längenabgeglichener CommonMark-Schluss); ein
  Befund entsteht, sobald **eine** von beiden am Dateiende offen endet — genau
  dann überspringt mindestens ein Modul den Rest. Dateiweit statt absatzweise
  (ein Fence *ist* eine Absatzgrenze), genau **ein** Befund je Datei, an der
  Öffnungszeile — eine **Fundstelle**, nicht zwingend die Reparaturstelle: welche
  von mehreren Öffnungen fehlt, ist grundsätzlich nicht entscheidbar. Kein neuer
  Config-Schlüssel; `spans` ist als Ganzes opt-in. Die Fence-**Paarung** selbst
  bleibt unverändert ([ADR-0042](docs/plan/adr/0042-markdown-lexik-folgt-commonmark.md))
  — sie löst den Fall nicht: ein nie geschlossener Fence ist unter jeder
  Paarungsregel offen.
  **Grenze:** die Prüfung greift nur auf Dateien im Scan-Scope. Post-Pässe über
  selbst benannte Verzeichnisse und Zieldateien außerhalb der Scan-Wurzeln, aus
  denen Module lesen (`matrix`, `anchors`, `codepaths`, `diagrams`, `versions`,
  `pins`), bleiben unerfasst; bei `pins` ist die Folge still. Der eingerückte
  Code-Block ist in keiner Lesart modelliert.

### Changed

- slice-101 — die Fence-Lexik trimmt an **allen** fünf Konsumenten identisch
  (Space und Tab, nicht unicode-weit). Das betrifft auch das mit 0.52.0
  ausgelieferte Closure-Gate des Moduls `planning`: hinter einer mit
  Unicode-Whitespace eingerückten Fence-Zeile verstummen `closure-note-thin` und
  `closure-note-missing`, während `closure-note-boilerplate` neu meldet.
  Fachlich war das immer die zugesagte Lexik (Spezifikation Schritt C4);
  gemessen wurde bis dahin eine andere. **Wer hier Befunde verliert, verliert
  Fehlmessungen.**
- slice-101 — der `--doctor`-Klartext zu `fence-unclosed` nennt die geltende
  Reichweite („mindestens ein Modul überspringt alles dahinter" statt „jedes
  Modul"). Nur Meldungstext, kein Grund-Code- oder Verhaltens-Delta.

## [0.52.0] — 2026-08-09

### Added

- slice-093 — **Closure-Note-Struktur** als zweite Fähigkeit des Moduls
  `planning` (`DC-FA-PLAN-001`,
  [ADR-0048](docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md)):
  opt-in über `planning.closure.dir` prüft d-check je abgeschlossenem Slice den
  Closure-Notiz-Abschnitt — vorhanden (`closure-note-missing`), genug
  Satzende-Zeichen **außerhalb** der Fenced-Code-Blöcke
  (`closure-note-thin`, Schwelle `planning.closure.min-sentences`, Default 4),
  keine der literalen `planning.closure.boilerplate`-Phrasen
  (`closure-note-boilerplate`, Liste **leer** per Default). Der Abschnitt ist der
  erste Treffer auf `planning.closure.heading-pattern` (Default
  `^#{2,3} .*Closure-Notiz`) bis zur nächsten gleich-/höherrangigen Überschrift.
  Hermetisch, diagnose-only, fail-closed bei fehlendem Verzeichnis, ungültigem
  Muster, `min-sentences` < 1 und leerem Floskel-Eintrag. **Kein neues Modul**;
  ohne `closure.dir` wird keine Slice-Datei geöffnet und der Befundsatz ist
  byte-identisch. Zugesagt ist **Struktur, nicht Bedeutung** — die semantische
  Schicht ist der neue Skill `.harness/skills/closure-note-reviewer.md`.
- slice-093 — neue CLI-Option **`--config <datei>`** (`DC-FA-CLI-012`): sie
  ersetzt für den Lauf die konventionelle `.d-check.yml` der Scan-Wurzel durch
  die genannte Datei (Pfad relativ zur Wurzel, muss **innerhalb** liegen).
  Gleiche strikte Validierung; fehlende Datei, Verzeichnis oder ein Pfad
  außerhalb der Wurzel ⇒ **Exit 2 ohne Rückfall** auf Defaults oder auf die
  konventionelle Datei. Sie **ersetzt**, sie ergänzt nicht — damit fährt ein Repo
  zwei disjunkte Prüf-Profile (Inner-Loop vs. Closure), ohne die Modulwahl auf
  der Kommandozeile nachzubauen. Ohne die Option unverändert byte-identisch.

### Fixed

- Befunde, die die Konfiguration als Herkunft nennen (Config-Pins des Moduls
  `sources`), tragen jetzt die **tatsächlich geladene** Datei statt hart
  `.d-check.yml` — unter `--config` zeigte die Meldung sonst auf eine Datei, die
  der Lauf nie gelesen hat.

## [0.51.1] — 2026-07-19

### Changed

- slice-081 — Modul `pins` (dpin): der `link-stale`-Befund führt den errechneten
  Ziel-Span-Hash jetzt **vollständig** (64 Hex) statt gekürzt (`shortHash`), damit
  ein Content-Pin praktisch anlegbar ist („einmal `--enable pins` laufen,
  gemeldeten Hash in den `<!-- dpin: sha256:… -->`-Marker kopieren"). Neue
  Handbuch-Aufgaben-Sektion §5 „Einen repo-internen Link-Inhalt gegen Drift
  pinnen". Nur die (nicht stabilitätsgarantierte) Befund-Meldung — kein
  Verhaltens-, Grund-Code- oder Vertragsdelta (`DC-FA-PIN-001` unverändert).

## [0.51.0] — 2026-07-19

### Added

- slice-080 — neues opt-in Regelmodul `sources` (19. Modul, **Netz**,
  `DC-FA-SRC-001`, [ADR-0046](docs/plan/adr/0046-sources-upstream-content-drift.md)):
  Upstream-Content-Drift externer Quellen. Eine auf einen `sha256` **gepinnte
  externe `http(s)`-Quelle** — per Marker `<!-- source-pin: [zip] sha256:… -->` am
  unmittelbar vorausgehenden Link **oder** per Config-Block `sources:` mit Einträgen
  `{url, sha256, unpack: none|zip}` — wird geholt, gehasht und verglichen. Abweichung
  → `source-drift` (die Meldung führt den **vollen Ist-`sha256`** als Re-Pin-Vorlage);
  Netzfehler, HTTP-Status ≥ 400, Timeout, Größenlimit **oder** kein gültiges Zip →
  `source-unreachable` (getrennt von `source-drift`). Zwei Quelltypen: Einzeldatei
  (Hash der rohen Antwort-Bytes) und Archiv (`unpack: zip` bzw. Marker-Keyword `zip`
  → pfad-sortiertes, gegen die Zip-Eintrags-Reihenfolge invariantes Content-Manifest).
  `sources` prüft nur absolute `http`/`https`-Ziele (kein Doppelbefund; repo-intern
  bleibt `pins`/`links`-Domäne); Content-Hash-Geschwister von `pins` (in-repo) und
  `external` (Netz-Erreichbarkeit).
- **Zweite Netz-Tür.** Neben `external` ist `sources` das einzige weitere Netz-Modul;
  beide sind strikt opt-in und nie im Default-Lauf — `DC-QA-03` (Seiteneffektfreiheit
  und Netz-Sparsamkeit) ist um `sources` erweitert. Der `sha256` wird case-insensitiv
  geführt; der Fetch ist größenbegrenzt (Body ≤ 64 MiB, Entpack ≤ 256 MiB bzw.
  ≤ 10 000 Einträge, Zip-Bomben-Schutz, Redirects ≤ 5); diagnose-only, kein
  `--repair`-Hunk. Ohne aktives `sources` ist der Befundsatz byte-identisch. Eine
  malformte Direktive (kein `sha256:<hex>`) oder ein ungültiger Config-Eintrag
  (fehlende `url`/`sha256`, unbekanntes `unpack`) bricht fail-closed mit Exit 2 ab
  (`DC-FA-SRC-001`, ADR-0046).

## [0.50.0] — 2026-07-18

### Added

- **Zitat-Verifikation** (`DC-FA-CODE-001`/`DC-FA-CITE-001`, §5/§6): d-check prüft
  `datei:zeile`-**Zitate**, die bisher still ins Leere zeigen konnten.
  - `codepaths.check-lines` (opt-in): verifiziert die Zeilen-Referenz eines
    Inline-Code-Pfads (`datei:<von>-<bis>`) — existierendes Ziel mit ≥ `bis` Zeilen
    (sonst `citation-out-of-range`) und `von ≤ bis` (sonst `citation-inverted-range`).
    Default aus **byte-identisch** (das Zeilen-Suffix wird wie bisher verworfen).
  - Neues, direktiven-getriebenes **18. Modul `citations`** (opt-in): die Direktive
    `<!-- d-check:cite <pfad>:<von>-<bis> -->` markiert das folgende Zitat
    (`>`-Blockquote **oder** inline `„…"`/`"…"`); Quell-Spanne und Zitattext werden
    whitespace-normalisiert, das Zitat muss ein zusammenhängender **Teilstring** der
    Quelle sein (sonst `citation-mismatch`). So besteht ein re-wrapped inline-Zitat,
    jede echte Wort-Abweichung bricht. Mindestlänge 16 Zeichen (schwache
    Kurz-Teilstring-Diskriminierung ausgenommen); Repo-Escape des Ziels ist ein
    Befund; eine malformte Direktive bzw. ein fehlendes Zitat ist fail-closed (Exit 2).

## [0.49.0] — 2026-07-18

### Added

- **Geteiltes Referenz-Ventil `ignore-refs` mit Quell-Skopus** (`DC-FA-REF-001`, §5):
  das bisher modul-lokale `codepaths.ignore-refs` wird zur **querschnittlichen**
  Top-Level-Fähigkeit, die `links`, `anchors` und `codepaths` gemeinsam honorieren.
  Neue Felder je Eintrag: `in` (Glob auf die **Quelldatei** — Quell-Skopus), `refs`
  (Globs auf das **aufgelöste Ziel**) und `keep` (Ausnahmen; ein von `refs`
  getroffenes Ziel bleibt geprüft, wenn `keep` es trifft — **reihenfolge-unabhängig**,
  kein gitignore-Last-Match). Das Ziel-Achsen-Pendant zu `scan.ignore`: es löst die
  **Template-Verzeichnis-Falle** — Ziel-Repo-Platzhalter, die im Quell-Repo massenhaft
  `target-missing` erzeugen, ohne die echten (Kurs-/Template-internen) Verweise blind
  zu machen. Der modul-lokale Schlüssel `codepaths.ignore-refs` bleibt als **Alias**
  gültig (kein Config-Bruch); ohne Block ist der Befundsatz byte-identisch. Ungültige
  `in`/`refs`/`keep`-Globs ⇒ Exit 2 (fail-closed). Die vier Ventil-Achsen (`scan.ignore`
  · `d-check:ignore` · `exempt-paths` · `ignore-refs`) sind im Handbuch gegeneinander
  erklärt.

## [0.48.1] — 2026-07-18

### Fixed

- **Direktiven-Toleranz in Tabellenzeilen** (`trace.requirements.format: table`,
  `trace.cross-consistency`): d-checks eigene Ignore-Direktive
  `<!-- d-check:ignore … -->` hinter der letzten Pipe einer Tabellenzeile brach den
  header-gebundenen Reader mit **Exit 2** ab — die eigene Konvention machte den
  eigenen Reader blind. Eine Datenzeile mit **genau einer** überzähligen,
  ganzzelligen HTML-Kommentar-Zelle wird jetzt auf Header-Breite gelesen; zwei
  überzählige Zellen oder eine Nicht-Kommentar-Zelle bleiben Exit 2. Sicher auf der
  Tabellengrenze aus v0.48.0 — eine relevante Folgetabelle wird nicht verschluckt
  (der zuvor fünffach reproduzierte stille Übersprung ist strukturell zu). Richtung
  **rot→grün** (ein spuriorer Abbruch wird lesbar): kein grüner Lauf wird rot, daher
  SemVer-**Patch** ([ADR-0040](docs/plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md),
  Spezifikation §DC-FA-REQ-001.a Schritt 5).

## [0.48.0] — 2026-07-18

### Changed

- **Tabellengrenze am relevanten Header** (`trace.requirements.format: table` und
  `trace.cross-consistency`): eine irrelevante Tabelle (ihr Header bindet keine
  konfigurierte Rolle) verschluckt die unmittelbar folgende **relevante** Tabelle
  nicht mehr **still**. Ein Header, der eine Rolle bindet, beendet die laufende
  Tabelle — damit wird jede relevante Tabelle erkannt, und ihre Anforderungen bzw.
  Rück-Kanten können nicht mehr lautlos in einer vorangehenden verschwinden. Ein
  rollenloser Header (z. B. eine `| - | - |`-Datenzeile) beendet **nicht** — die
  Gegenprobe, an der fünf rein strukturelle Fassungen scheiterten.

  **Wirkung auf bestehende Läufe:** d-check **findet mehr** — eine bisher
  verschluckte relevante Tabelle liefert nun Anforderungen (neue Waisen möglich);
  ein mehrdeutig-relevanter Header (doppelte Rollen-Spalte) hinter einer
  irrelevanten Tabelle bricht jetzt laut mit **Exit 2** ab, statt still verschluckt
  zu werden. Ein heute grüner Konsumentenlauf kann danach rot sein — laut und in
  der sicheren Richtung. **Defekt-Fix, kein Lastenheft-Change** (die Regel ist
  Spezifikations-Sache), aber SemVer-**Minor**
  ([ADR-0043](docs/plan/adr/0043-tabellengrenze-am-relevanten-header.md),
  Spezifikation §DC-FA-REQ-001.a Schritt 5).

## [0.47.0] — 2026-07-18

### Changed

- **Markdown-Lexik an CommonMark/GFM angeglichen** — zwei belegte, still
  ausgelieferte Blindstellen geschlossen (Differential-Spike gegen goldmark
  v1.8.4 über 522 reale Dateien: 8 Abweichungen, alle „d-check sieht eine
  Tabelle nicht, die jeder Renderer zeigt"):
  - **Trennzelle folgt GFM:** eine Tabellen-Trennzeile braucht nur **einen**
    Bindestrich (`| - |`), nicht drei. Jede reale Tabelle mit `| -- |` war zuvor
    für d-check keine Tabelle — ihre Anforderungen, Links und IDs existierten
    nicht.
  - **Fence-Infozeilen-Regel (CommonMark):** öffnet eine Backtick-Fence-Zeile
    ihren Block mit einer Infozeile, die selbst noch einen Backtick enthält, ist
    sie kein Fence-Öffner, sondern Fließtext. Zuvor blendete ein erklärender Satz
    *über* einem Codeblock den **gesamten Rest der Datei für alle Module** aus —
    ein kaputter Link dahinter verschwand lautlos (Exit 1 ⇒ Exit 0, gemessen).
    Für `~~~`-Blöcke gilt die Regel nicht.

  **Wirkung auf bestehende Läufe:** d-check **findet mehr** — bisher unsichtbare
  Tabellen liefern Anforderungen (neue Waisen möglich), bisher unsichtbare Prosa
  liefert Links und IDs (neue Befunde möglich). Ein heute grüner Konsumentenlauf
  kann danach rot sein — laut, nicht still, und in der sicheren Richtung. Beide
  Regeln sind per Mutation gepinnt. **Defekt-Fix, kein Lastenheft-Change** (die
  Regeln sind Spezifikations-Sache), aber SemVer-**Minor**, weil d-check danach
  mehr findet ([ADR-0042](docs/plan/adr/0042-markdown-lexik-folgt-commonmark.md),
  Spezifikation §DC-FA-REQ-001.a Schritt 3 / §DC-FA-LINK-001.a Schritt 1).

## [0.46.0] — 2026-07-17

### Added

- **Komma-Kurzform in Coverage-Quellen ist fail-closed** (`trace.coverage`,
  `trace.cross-consistency`). Folgt einer Kennung — oder ihrer Range/Enum-Notation
  — ein Komma und **unmittelbar Ziffern** (`GG-SCN-001, 007`, ebenso
  `GG-SCN-001..005, 007, 008`), bricht d-check mit **Exit 2** und einem Hinweis auf
  die zugesagten Notationen ab, statt die Kurzform still fallen zu lassen. Bis
  v0.45.1 verschwand `007` lautlos und erzeugte in `trace.coverage` eine falsche
  Waise — bei einem produktiv verdrahteten Konsumenten. Die Kurzform war nie eine
  zugesagte Notation; der Defekt war das fehlende Signal. Ein Komma **vor** einer
  vollständigen Kennung (`GG-SCN-001, GG-SCN-007`, auch hinter einer Range) bleibt
  unberührt ([ADR-0041](docs/plan/adr/0041-komma-kurzform-fail-closed.md),
  Lastenheft 0.46.0).

  **Wirkung auf bestehende Läufe:** eine Quelle mit Komma-Kurzform **oder** einer
  Prosa-Zahl direkt hinter einer Kennung (`GG-QA-001, 2026`) läuft künftig auf
  Exit 2 — laut und in Sekunden behebbar. Das ist ein **Vertrags-Zuwachs** (neues
  Akzeptanzkriterium), daher SemVer-**Minor**. Wer nur `..`/`/`-Notation oder
  komma-getrennte volle Kennungen nutzt, ist nicht betroffen.

## [0.45.1] — 2026-07-17

### Fixed

- **Keine Falsch-Expansion mehr bei Klammern im Linkziel** (`trace.coverage`,
  `trace.cross-consistency`). Die Link-Transparenz aus v0.44.1 grenzte das Linkziel
  per Regex an der **ersten** `)` ab — eine zweite Link-Definition neben dem
  kanonischen Reader, der klammer-**balanciert** liest. Enthielt ein Ziel eine
  Klammer, landete der URL-Rest im Range-Parser: `` [`GG-QA-001`](…/Rev(2)/002/003.md) ``
  expandierte die Pfadsegmente `/002/003` als Enum — in einer Zelle **ohne jede
  Range-Notation**. Folge: **falsche Deckung, die Waisen versteckt**
  (`--require-complete` Exit 1 → Exit 0).

  **Betroffen: v0.44.1 und v0.45.0.** Wer ein Linkziel mit Klammern in einer
  `trace.coverage`-Quelle führt, sollte aktualisieren — der Defekt macht Läufe
  fälschlich **grün**, ist also still. Alle anderen Formen sind unverändert.

- **Klammern im Linkziel brechen die Range-Fortsetzung nicht mehr.** Umgekehrte
  Richtung derselben Wurzel: `` [`GG-QA-001`](https://x.org/A_(b))..003 `` expandiert
  jetzt. Damit ist die strukturelle Kollision zwischen Range-Notation und
  Linkpflicht auch für Ziele mit Klammern aufgelöst — `ids` sah den Link als
  erfüllt, der Range-Parser nicht.

### Changed

- Die Link-Abgrenzung des Range-Parsers kommt jetzt aus `rules.LinkSuffixEnd` —
  **eine** Link-Definition im Repo statt zwei
  ([ADR-0039](docs/plan/adr/0039-link-transparente-range-fortsetzung.md)).

## [0.45.0] — 2026-07-17

### Added

- **`trace.cross-consistency.forward.req-pattern`** (RE2, Default
  `trace.requirements.id-pattern`) — symmetrisch zum vorhandenen
  `backward.req-pattern`. Er erkennt die Anforderungs-IDs der Vorwärts-ID-Spalte
  und trennt damit den **Vergleichs-Scope vom RTM-Scope**: welche Anforderungen
  der Abgleich vergleicht, entscheidet das Muster — **nicht**, ob eine Anforderung
  in der RTM steht. `--print-config` führt den Schlüssel samt Warnung.

### Fixed

- **Kein Falschbefund mehr bei bewusst gescopter RTM.** Bis v0.44.1 las die
  Vorwärts-Sicht ihre IDs **still** über `trace.requirements.id-pattern`, während
  die Rück-Sicht ihr eigenes Muster nutzte — die Kopplung stand weder im Vertrag
  noch in der Config-Oberfläche. Schließt ein Repo eine Familie bewusst aus der
  RTM aus (etwa Architektur-Meta, das keine Anforderung ist), war `F(R)` für sie
  leer: **jede** Rück-Kante wurde als „Rück-Kante, ohne RTM-Eintrag" gemeldet, und
  die eigentliche `F \ B`-Differenz **verschwand**. Der Lauf sah aus wie ein
  Treffer, war aber ein Nebeneffekt der leeren Sicht. Wer bisher keine gescopte RTM
  fährt, ist nicht betroffen (byte-identisch).

### Changed

- `DC-FA-XREF-001` (Lastenheft 0.45.0) hält nun ausdrücklich fest: **die
  Vergleichs-Schlüsselmenge ist nicht die RTM-Anforderungsmenge**, und der
  Default-Rückfall auf `requirements.id-pattern` ist eine sichtbare
  Konfigurationsentscheidung, keine Ableitung.

## [0.44.1] — 2026-07-17

### Fixed

- **Verlinkte Range-/Enum-Fortsetzungen expandieren wieder** (`trace.coverage`,
  `trace.cross-consistency`): Der geteilte Range-Parser überspringt hinter einer
  Kennung **genau ein** Markdown-Link-Suffix und liest die Fortsetzung dahinter.
  Bis v0.44.0 lieferte `` [`GG-UI-001`](…)..003 `` nur `GG-UI-001`, während das
  unverlinkte `GG-UI-001..003` korrekt expandierte — die Range-Zusage kollidierte
  strukturell mit d-checks eigener Linkpflicht (`ids` mit `link-policy: always`).
  Der Defekt bestand in `trace.coverage` **seit v0.41.0** und erzeugte dort
  **falsche Waisen**, die unter `--require-complete` fälschlich gateten
  ([ADR-0039](docs/plan/adr/0039-link-transparente-range-fortsetzung.md)).

  **Wirkung auf bestehende Läufe:** Wer verlinkte Ranges führt, sieht **weniger**
  Waisen bzw. Differenzen — ein fälschlich roter Lauf wird grün.

  **Richtigstellung (0.45.1):** Die ursprüngliche Zusage „Wer keine verlinkten
  Ranges nutzt, ist nicht betroffen (byte-identisch)" war **falsch** — bei
  Klammern im Linkziel expandierte auch eine Zelle ohne Range-Notation und
  versteckte Waisen. Behoben in 0.45.1; die Zusage gilt erst ab dort.

  Bewusst eng gefasst: übersprungen wird **nur** ein Link-Suffix — nicht
  Whitespace, Emphasis, ein zweites Suffix oder Text zwischen `)` und der
  Fortsetzung. Jede weitere Toleranz würde die Autor-Absicht raten. Die
  Fail-closed-Fälle (`AAA>BBB`, abweichende Ziffern-Breite) gelten unverändert
  auch hinter einem Link. (Die hier ursprünglich behauptete Aussage „Enthält das
  Linkziel selbst eine Klammer, greift die Regel nicht" traf **nicht** zu — siehe
  0.45.1.)

## [0.44.0] — 2026-07-17

### Added

- **Kreuzverweis-Konsistenz zweier Traceability-Sichten** (`trace.cross-consistency`,
  opt-in): Der `--trace`-Lauf vergleicht zusätzlich eine **Vorwärts**-RTM-Tabelle
  (Anforderung → Design-Artefaktmenge) gegen die **Rückwärts**-`Bezug`-Kanten
  (Design → Anforderung, die Quelle der Wahrheit) und meldet je Anforderung beide
  Mengendifferenzen mit Richtungslabel und `Datei:Zeile`. Modi `equal` (beide
  Richtungen gaten) und `superset` (nur „Rück-Kante ohne RTM-Eintrag"); Ventil
  `exclude-req` für Ableitungssprünge in Mittelschichten; `artifact-id-column: first`
  nimmt die erste Spalte, wenn die ID-Header über die Tabellen heterogen sind.
  Beide Sichten laufen über den vorhandenen header-gebundenen Tabellen-Reader und
  die range-aware Span-Semantik — kein neuer Parser.
- Der Abgleich ist **advisory** (`--trace` bleibt Exit 0) und gatet allein über das
  globale `--require-complete` (≥ 1 Differenz ⇒ Exit 1) — kein block-lokaler
  Schalter. `--print-config` führt den neuen Block.
- **Fail-closed** (Exit 2): fehlendes `forward`/`backward`, unbekannter `mode`,
  nicht kompilierbares Regex, leeres Pflichtfeld, fehlende Sicht-Quelle, keine
  Tabelle mit den konfigurierten Headern, mehrfacher Rollen-Header, Sektionsname
  ohne Heading-Treffer, Zellenzahl-Bruch — sowie ein **vakuumer Abgleich**: greifen
  die Muster am Inhalt vorbei oder verschluckt `exclude-req` jede Anforderung, kann
  der Lauf konstruktionsbedingt nie eine Differenz melden; ein `0 Differenz(en)`
  behauptete dann eine nie geprüfte Konsistenz.

### Unchanged

- Ohne `trace.cross-consistency`-Block ist die RTM in allen drei Formaten
  **byte-identisch** (Markdown, `--json`, `--yaml`); es wird nichts geschrieben und
  kein Netz berührt. Eine einseitig leere **Vorwärts**-Sicht bei gepflegten
  Rück-Kanten ist **kein** Fehler, sondern meldet jede Rück-Kante laut — der
  erwartete Zustand, solange die RTM-Tabelle noch nicht auf konkrete IDs
  restrukturiert ist.

## [0.43.1] — 2026-07-15

### Fixed

- `--trace` Tabellen-Parser: ein ungeschlossener Backtick-Span bleibt literal,
  die folgende Pipe trennt wieder Zellen — verhindert falsche Zeilenbreiten-Fehler.
- `trace.requirements.table`: die Exklusivität von `text-column`/`text-columns`
  ist präsenz- statt inhaltsbasiert (auch `text-column: ''` bzw. `text-columns: []`
  neben der anderen Form ist fail-closed).
- Fehlerpräzedenz: die Header-Prüfung (unbenutzte `text-columns`-Alternative) greift
  vor der Duplicate-ID-Meldung (DC-FA-REQ-001.a „Header → Duplicate-ID“).

## [0.43.0] — 2026-07-14

### Added

- [`slice-070`](docs/plan/planning/done/slice-070-trace-tabellenquellen-nullmengen-guard.md)
  ergänzt für `--trace` native Markdown-Pipe-Tabellen als Anforderungsquelle
  ([`DC-FA-REQ-001`](spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen),
  [ADR-0037](docs/plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md)).
  `trace.requirements.format: table` bindet ID-, einen oder mehrere alternative
  Text-Header und eine optionale Modalitätsspalte über exakte Namen;
  `duplicate-ids` bietet die expliziten Brownfield-Politiken `first`/`last`
  neben dem sicheren Default `error`. Escaped Pipes und Pipes in einzeiligen
  Code-Spans bleiben Teil der Zelle.

### Changed

- Eine nichtleer explizite `trace.requirements.source` oder der Tabellenmodus
  bricht bei fehlender Quelle beziehungsweise null erkannten Anforderungen nun
  mit Exit 2 ab. `source: ""` und der unkonfigurierte Heading-Default behalten
  das bisherige Verhalten byte-identisch; mehrdeutige Tabellen-Header,
  fehlerhafte Zeilenbreiten und doppelte IDs werden fail-closed abgewiesen.

### Fixed

- [`slice-069`](docs/plan/planning/done/slice-069-trace-handbuch-parsergrenzen.md)
  präzisiert die Dokumentation von `--trace` für v0.42.0: Anforderungen
  werden nur aus ATX-Überschriften mit der ID als erstem vollständigem Token
  definiert; Tabellen-/Listen-/Fließtext, Body-only-Modalität, die exakte
  Waisen- und Referenzscan-Semantik sowie die leere RTM mit Exit 0 sind nun
  ausdrücklich beschrieben. Eine Brownfield-Anleitung zeigt die Migration
  tabellarischer Lastenhefte, ohne native Tabellenunterstützung zu behaupten.

## [0.42.0] — 2026-07-11

### Added

- slice-068 — neue opt-in **Modalitäts-Klassifikation `trace.requirements.modality`**
  ([`DC-FA-MOD-001`](spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in),
  [ADR-0036](docs/plan/adr/0036-trace-modality-klassifikation.md)) für die
  Requirements Traceability Matrix. d-check klassifiziert jede Anforderung nach
  RFC-2119-Stufe (MUSS/SOLLTE/KANN) anhand **konfigurierbarer Modal-Verb-
  Schlüsselwörter** (Built-in DE+EN-Defaults; `levels` Stufe→Keywords,
  `require-levels` welche Stufen gaten, Default `[must]`) über den ersten/
  längsten Treffer im **normalisierten** Anforderungs-Body (`MUSS NICHT`=may vor
  `MUSS`=must, `DARF NICHT`=must; whitespace-/emphasis-normalisiert, wortgrenzen-
  genau); kein Treffer ⇒ Stufe `unknown`. Neue konditionale **Modality**-Spalte;
  `--require-complete`
  ([`DC-FA-CLI-011`](spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code))
  bricht dann **nur** bei Waisen der `require-levels`-Stufen — SOLLTE/KANN/`unknown`
  advisory. Fail-closed (leeres Keyword / gleiches Keyword in zwei Stufen /
  reservierter Name `unknown` / ungültiges `require-levels` ⇒ Exit 2). **Ohne
  `modality` byte-identisch** (keine Spalte, kein Feld, alle Waisen gaten;
  `DC-QA-02`). Anlass: grid-gyms 10 Coverage-Rest-„Waisen" sind 5× KANN + 4×
  Nicht-Ziele + 1× DARF NICHT; an den Realdaten gaten mit `modality` nur noch die
  2 echten MUSS-Lücken (`GG-MVP-004` `DARF NICHT`, `GG-NONGOAL-005` „muessen").

## [0.41.0] — 2026-07-11

### Added

- slice-067 — neue opt-in **Coverage-Referenzklasse `trace.coverage`**
  ([`DC-FA-COV-001`](spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in),
  [ADR-0035](docs/plan/adr/0035-trace-coverage-quellen.md)) für die Requirements
  Traceability Matrix. Eine **Liste** kuratierter Quellen liest Deckungs-Matrizen
  (z. B. eine ausgelagerte Traceability-Datei) als **eigene Coverage-Dimension**
  ein — je Quelle `files` (explizite Pfade, keine `dir`/`file-pattern` → keine
  ADR-Kontamination), `label` (Owner-Kennung in einer eigenen **Coverage**-Spalte),
  `ranges` (Default true; expandiert `<FAM>-AAA..BBB` und `<FAM>-AAA/BBB/CCC`
  breiten-erhaltend, gegen `requirements.id-pattern` validiert) sowie `sections`
  (Whitelist) / `exclude-sections` (Blacklist) über die
  `matrix.exclude-sections`-Span-Semantik (voller Heading-Klartext). Eine
  Anforderung ist damit **Waise** nur ohne Slice **und** ohne Coverage
  ([`DC-FA-CLI-011`](spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
  angepasst); `--json`/`--yaml` tragen ein `coverage`-Feld. Fail-closed (fehlende
  Datei / leeres `label` / Sektionsname ohne Heading-Treffer / ungültige Range
  ⇒ Exit 2). **Ohne `trace.coverage` ist die RTM byte-identisch** (keine
  Coverage-Spalte, kein Feld; `DC-QA-02`). Anlass: Konsument grid-gym — 171
  „Waisen" waren zu ≥122 anderswo (ADR/Traceability-Matrix/Wellen) belegt; mit
  `trace.coverage` (Range + `exclude-sections`) sinken sie an den realen Daten
  von 113 auf 10.

## [0.40.0] — 2026-07-11

### Changed

- slice-066 — die **Requirements Traceability Matrix** (`--trace`,
  [`DC-FA-CLI-009`](spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix),
  [ADR-0034](docs/plan/adr/0034-trace-konfigurierbare-quellen.md)) ist über einen
  opt-in **`trace`-Config-Block** in `.d-check.yml` quell- und
  kennungs-konfigurierbar. Bislang waren die vier RTM-Annahmen hart an d-checks
  eigene Konvention gebunden — Anforderungs-Quelldatei + Kennungs-Gestalt
  (`-FA-`/`-QA-`) sowie die Slice-/ADR-Dateinamen (`slice-NNN-…`/`NNNN-…`). Neu
  überschreibt `trace.requirements.source`/`.id-pattern`,
  `trace.adrs.dir`/`.file-pattern`/`.id-prefix` und
  `trace.slices.dir`/`.file-pattern`/`.id-prefix` diese Achsen (Capture-Gruppe 1
  der `file-pattern` = Owner-Kennung). **Jedes Feld ist optional; ohne
  `trace`-Block ist die RTM byte-identisch** zum bisherigen Verhalten
  ([`DC-QA-02`](spec/lastenheft.md#dc-qa-02--determinismus)); fail-closed bei
  ungültiger Regex oder `file-pattern` ohne Capture-Gruppe (Exit 2). Damit bildet
  die RTM auch Konsumenten-Repos mit abweichender Kennungs-/Datei-Konvention
  vollständig ab (Anlass: grid-gym sah 6 von 243 Anforderungen).
  `--require-complete`
  ([`DC-FA-CLI-011`](spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code))
  erbt die konfigurierten Quellen. `--print-config` führt einen kommentierten
  `trace`-Block.

## [0.39.0] — 2026-07-06

### Changed

- slice-065 — `--suggest-config ai-harness[-init]` gleicht das vorgeschlagene
  Modulset an die **gelebte** Dogfood-Konvention an
  ([`DC-FA-CLI-006`](spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten),
  [ADR-0033](docs/plan/adr/0033-ai-harness-template-modulset.md)): `spans` und
  `hostpaths` gehören jetzt zur **fixen** Aktiv-Menge, dazu ein **repo-bewusster
  `planning`-Block** (aktiv bei vorhandener Roadmap bzw. im Voll-Kanon
  `ai-harness-init`, sonst auskommentiert). `vcs`/`commits` (die eine
  Commit-Range brauchen) werden auf `--print-mk` verwiesen statt ins statische
  `modules` gelegt; `versions`/`targets` bleiben bewusst vertagt
  (repo-spezifische `pin-pattern`/`authority`). Ein Eignungs-Kriterium K1–K4
  macht die Modul-Aufnahme explizit; die kanonische Vorlage der Spezifikation
  (`DC-FA-CLI-006.a`) deckt die emittierte Ausgabe nun **1:1**. Betrifft nur die
  `ai-harness`-Vorlage — nicht den generischen Quellen-Modus und nicht die eigene
  `.d-check.yml`.

## [0.38.0] — 2026-07-05

### Added

- slice-063 — neues opt-in Regelmodul `targets` (17.), das die
  Deklarations-Konsistenz zwischen Doku und Build-Targets prüft
  ([`DC-FA-TGT-001`](spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in),
  [ADR-0031](docs/plan/adr/0031-targets-deklarations-konsistenz-modul.md)):
  ein in einer Doku-**Tabellenzeile** als `` `make X` `` behauptetes Target
  ohne Makefile-Regel ⇒ `gate-phantom` (Richtung 1); jede Makefile-Regel
  (minus `exempt-targets`) ohne Eintrag in der Autoritäts-Doku ⇒
  `gate-undocumented` (Richtung 2). **Hermetisch** (nur der Filesystem-Port,
  kein git/Netz/Makefile-Ausführen), fail-closed bei fehlender Datei,
  default-aus byte-identisch.
- `--print-mk` trägt das elfte Target `doc-targets`, `--print-config`/
  `--suggest-config` führen `targets`
  ([`DC-FA-CLI-010`](spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)).

### Changed

- Das Meta-Gate `make gate-consistency` dogfoodet nun das Modul `targets`
  für den Doku-↔-Makefile-Kern (via Image); `tools/gate-consistency.sh` ist
  auf die repo-spezifische DC-QA-03-Modullisten-Prüfung reduziert — der
  cross-repo-driftende Skript-Kern ist mechanisiert und verteilbar
  ([ADR-0031](docs/plan/adr/0031-targets-deklarations-konsistenz-modul.md)).

## [0.37.1] — 2026-07-04

### Fixed

- slice-060 — `--doctor`-Klartext-Vollständigkeit
  ([`DC-FA-CLI-007`](spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)):
  die Diagnose zeigte für sieben seit v0.25 hinzugekommene Grund-Codes
  (`diagram-id-undefined`, `version-stale`, `link-stale`, `core-drift`,
  `core-drift-vcs`, `commit-untraceable`, `planning-drift`) den rohen Code
  statt des Klartexts; `--doctor --json`/`--yaml` trug ihn im
  `reasonText`-Feld. Die sieben Klartexte sind ergänzt; ein neuer Test
  verriegelt die kanonische Grund-Code-Liste beidseitig gegen die
  §4-Grund-Code-Tabelle der Spezifikation (fail-closed) — die Lücke kann
  nicht wieder still wachsen.

## [0.37.0] — 2026-07-03

### Added

- slice-059 — neues opt-in Modul `tracked` (16. Regelmodul,
  [`DC-FA-TRK-001`](spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)):
  prüft jedes **auflösbare, existierende** Link-/Bild-Ziel gegen den
  **git-Index** — ein untracktes/gitignoriertes Ziel wäre auf jedem frischen
  Klon `target-missing`; Befund `target-untracked` fängt die Umgebungs-Drift
  am Entstehungsort. Index-Wahrheit (gestagt = getrackt, keine
  `.gitignore`-Interpretation), kein Doppelbefund (fehlende Ziele bleiben
  `links`), Ventil `tracked.exempt-targets` (referenz-weit); liest `.git`
  read-only über den VCS-Port (reine-Go, **ohne** Range), fail-closed ohne
  `.git` (Exit 2), strikt opt-in/default-aus byte-identisch, diagnose-only.
- `--print-mk`: das Fragment trägt zusätzlich **`doc-tracked`**
  (`--enable tracked` + fokussierte `--disable`-Liste, ohne Range;
  [`DC-FA-CLI-010`](spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  9→10 Targets); `--print-config`/`--suggest-config` führen `tracked`.

## [0.36.0] — 2026-07-01

### Added

- slice-057 — neues opt-in Modul `planning` (15. Regelmodul,
  [`DC-FA-PLAN-001`](spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)):
  prüft **hermetisch** (nur Roadmap-Datei + Verzeichnis-Listing, **kein** git,
  **kein** Netz) die Planning-Lifecycle-Invariante — die Roadmap trägt den
  Ruhe-Marker (`planning.marker`, „Keine aktive Welle") in ihrem
  `planning.heading`-Abschnitt genau dann, wenn kein `slice-*` (`planning.slice-glob`)
  im Roadmap-Verzeichnis liegt (`hasActive == hasSlices`), sonst `planning-drift`.
  Fail-closed bei fehlender kanonischer Überschrift/Roadmap-Datei (Heading-Guard);
  strikt opt-in, diagnose-only, default-aus byte-identisch; `heading`/`marker`/
  `slice-glob` überschreibbar (parametrierbar für Adopter)
  ([ADR-0028](docs/plan/adr/0028-planning-lifecycle-modul.md), Lastenheft 0.36.0).
  Das **Gate `make planning-check`** läuft dogfood auf das Modul um (Image,
  `--enable planning`); `--print-config`/`--suggest-config` führen `planning`, und
  `--print-mk` trägt ein `doc-planning`-Target (hermetisch, ohne Range;
  [`DC-FA-CLI-010`](spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  8→9). Benutzerhandbuch §5/§6 dokumentiert das Modul.

### Removed

- slice-057 — `tools/planning-consistency.sh` entfernt (das **letzte** Gate-Skript
  des `tools/*.sh`-Audits), abgelöst durch das Modul `planning`. Die immutable
  [slice-040](docs/plan/planning/done/slice-040-planning-consistency-gate.md)-Inline-Referenz
  ist über `codepaths.ignore-refs` als Tombstone deklariert (vierter Fall nach
  `adr-immutable-check.sh`/`completeness-check.sh`/`trace-check.sh`); der
  3-Richtungs-Negativ-Selbsttest lebt als Modul-Akzeptanztest (`make test`) weiter.

## [0.35.0] — 2026-07-01

### Added

- slice-056 — neues opt-in Modul `commits` (14. Regelmodul,
  [`DC-FA-COMMITS-001`](spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)):
  prüft, dass jede geprüfte **Commit-Message** eine Traceability-Kennung nach
  `commits.id-patterns` (`ADR-`/`MR-`/`DC-`/`slice-`) auf einer Inhalts-Zeile
  trägt, sonst `commit-untraceable`. Liest die Commit-Messages über **denselben
  reine-Go-VCS-Port** wie `vcs` (erweitert um `CommitMessages`; **ohne**
  git-Binary → distroless bleibt, **ohne** Netz); zwei Quellen, eine Prüfung:
  `--range <base>..<head>` (CI/Push, Nicht-Merge-Commits) und
  `--commit-msg <datei|->` (commit-msg-Hook, einzelne Pending-Message via stdin).
  Uniforme `#`-/scissors-Bereinigung (Kennung auf Inhalts-Zeile, nicht im
  Kommentar), Betreff-Ausnahme `commits.exempt-pattern` (Selbstkonfig
  `^(Merge |Revert )`). Strikt opt-in (nie Default, wie `external`/`vcs`),
  fail-closed ohne `.git`/Range/Message, diagnose-only; default-aus byte-identisch
  ([ADR-0027](docs/plan/adr/0027-commits-traceability-modul.md), Lastenheft 0.35.0).
  Das **Gate `make trace-check`** läuft dogfood auf das Modul um (Image,
  `--enable commits` bzw. `--commit-msg -`). Config-Surface vollständig nachgezogen:
  `--print-config`/`--suggest-config` führen `commits` in der Verfügbar-/opt-in-Liste,
  und `--print-mk` trägt ein `doc-commits`-Target (verteilte Commit-Traceability,
  parallel zu `doc-immutable`; [`DC-FA-CLI-010`](spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  7→8 Targets). Benutzerhandbuch §5/§6 dokumentiert das Modul.

### Removed

- slice-056 — `tools/trace-check.sh` entfernt (das **letzte** Gate-Skript der
  Familie), abgelöst durch das Modul `commits`. Die immutable
  [ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md)-Inline-Referenz
  ist über `codepaths.ignore-refs` als Tombstone deklariert (dritter Fall nach
  `adr-immutable-check.sh`/`completeness-check.sh`);
  [ADR-0027](docs/plan/adr/0027-commits-traceability-modul.md) supersedet die
  Skript-Mechanik von [ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md)
  (Policy, Bindepunkt und CI-Topologie unverändert). Der Negativ-Selbsttest lebt
  als Modul-Akzeptanztest (`make test`) weiter.

## [0.34.0] — 2026-06-29

### Added

- slice-054 — Modul `codepaths`
  ([`DC-FA-CODE-001`](spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in))
  um `codepaths.ignore-refs` erweitert: eine Glob-Liste nimmt **bestimmte
  aufgelöste Ziel-Pfade** referenz-weit (datei-/zeilen-unabhängig) von der
  Existenz-/Escape-/Anker-Prüfung aus — ein **Tombstone-Register** bewusst
  entfernter Artefakte, die immutable oder historische Doku noch in Inline-Code
  zitiert. Löst die Frozen-Doc-Refactoring-Falle (eingefrorene Verweise auf
  refaktorierten/gelöschten Code dangeln sonst als `codepath-missing` an
  uneditierbarer Doku) — ohne Edit an immutabler Doku, ohne ganze Doc-Klassen aus
  dem Check zu nehmen. Bewusster Akt mit Gate (ohne Eintrag bleibt
  `codepath-missing`); Default leer → byte-identisch
  ([ADR-0025](docs/plan/adr/0025-codepaths-ignore-refs.md), Lastenheft 0.34.0).
  Dritte Ventil-Achse neben `d-check:ignore` (Zeile) und `exempt-paths` (Datei).

### Removed

- slice-054 — `tools/adr-immutable-check.sh` entfernt (in v0.33.0 durch das Modul
  `vcs` abgelöst und nur noch als „pfad-stabiler Fallback" gehalten). Die
  immutablen [ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md)/[ADR-0024](docs/plan/adr/0024-vcs-immutable-gate.md)-Inline-Referenzen
  sind über `codepaths.ignore-refs` als Tombstones deklariert;
  [ADR-0025](docs/plan/adr/0025-codepaths-ignore-refs.md) nimmt die
  „Skript-behalten"-Teilentscheidung von ADR-0024 zurück (deren VCS-Port/Modul-
  `vcs`-Kern bleibt gültig).

## [0.33.0] — 2026-06-29

### Added

- slice-053 — Neues opt-in-Regelmodul `vcs` (13. Modul,
  [`DC-FA-VCS-001`](spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)):
  git-historienbasierte Immutabilität des **Core** über eine Commit-Range
  (`core(BASE)` vs. `core(HEAD)`), geliefert über `--range <base>..<head>`
  (CI/Push) oder `--staged` (pre-commit). Liest das read-only `.git` über einen
  eigenen, **reine-Go-VCS-Port** (go-git; **ohne git-Binary** → das
  distroless-Image bleibt unangetastet, **ohne Netz**); erweiterter Eingabe-Scope
  (git + Range), aber lokal/lesend/deterministisch (`DC-QA-02`/`DC-QA-03`
  unberührt). Geprüft wird jede in der Range geänderte Datei der Klasse
  `vcs.paths` mit immutabler BASE (`vcs.immutable-when`): Körper-Drift,
  unzulässiger Status-Übergang (`vcs.head-allow`) oder Löschung/Umbenennung →
  `core-drift-vcs`. Core-Semantik in Parität zum bisherigen
  `adr-immutable-check.sh` (nur Kopf-Status-Zeile gestrippt, `exclude-sections`).
  Strikt opt-in, fail-closed ohne `.git`/Range, diagnose-only. **Die volle
  git-Garantie, die der hermetische `immutable`-Pin (v0.32.0) bewusst der
  VCS-Stufe überließ** — beide koexistieren als Defense-in-Depth
  ([ADR-0024](docs/plan/adr/0024-vcs-immutable-gate.md), Lastenheft 0.33.0).
- slice-053 — `--print-mk` (`DC-FA-CLI-010`) trägt ein `doc-immutable`-Target:
  Schwester-Repos beziehen die git-Immutabilität über das gepinnte Image, ohne
  ein Skript zu kopieren (RANGE/STAGED-getrieben, auf `vcs` fokussiert) — der
  Verteilungs-Kern hinter dem Modul.

### Changed

- slice-053 — Das `adr-check`-Gate (Accepted-ADR-Immutabilität) läuft jetzt über
  das Modul `vcs` (Dogfood, im Image verteilt) statt über
  `tools/adr-immutable-check.sh` (bleibt pfad-stabiler Fallback); `pre-commit`-Hook
  und CI rufen `make adr-check` ([ADR-0024](docs/plan/adr/0024-vcs-immutable-gate.md)
  löst die Skript-Mechanik von [ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md)
  ab). Neues `make tidy` (go.mod/go.sum-Pflege in Docker).
- slice-053 — Config-Surface-Bereinigung: `--print-config` (`DC-FA-CLI-005`) führt
  jetzt alle Module (`pins`/`immutable`/`vcs` ergänzt) und `--suggest-config
  ai-harness` (`DC-FA-CLI-006`) nennt die situativen opt-in-Module vollständig
  (`versions`/`pins`/`immutable`/`vcs` nachgezogen).

## [0.32.0] — 2026-06-28

### Added

- slice-052 — Neues opt-in-Regelmodul `immutable` (12. Modul, `DC-FA-IMM-001`):
  Immutabilitäts-Pin gegen Core-Drift. Eine Datei mit dem Inline-Marker
  `<!-- immutable: sha256:<hex> -->` wird gegen den whitespace-normalisierten
  **Core** gehasst — den Datei-Inhalt **ohne** die Marker-Zeile und ohne die per
  `immutable.exclude-sections` benannten Abschnitte (für ADRs typisch
  `Geschichte`); Abweichung → Grund-Code `core-drift`. **Hermetisch** (kein git,
  rein im read-only gescannten Arbeitsbaum); die git-historienbasierte
  `core(BASE)`-vs-`core(HEAD)`-Form bleibt einem späteren opt-in VCS-Adapter
  vorbehalten. Mechanik erbt die `pins`-Normalisierung; Marker auf der
  vorverarbeiteten Zeile (in Fenced-/Inline-Code inert), erster Marker je Datei.
  Diagnose-only, opt-in pro Datei, default-off byte-identisch
  ([`DC-FA-IMM-001`](spec/lastenheft.md#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in),
  [ADR-0023](docs/plan/adr/0023-immutable-core-pin.md), Lastenheft 0.32.0).

## [0.31.0] — 2026-06-28

### Added

- slice-051 — Modul `matrix` um die **token-basierte** Referenz-Richtung erweitert
  (`DC-FA-MTX-003`): eine Klasse kann zusätzlich zu `paths` ein `token`-Regex
  tragen, mit dem `matrix` verbotene Referenzen auch als **bare ID-Token** im
  Prosa-Körper fängt (nicht nur als Link) — `matrix-forbidden` in Token-Form. Ein
  **Provenance-Marker** `<!-- d-check:status-provenance -->` auf derselben Zeile
  nimmt eine verbotene Token-Referenz aus (deklarierte Provenance/Verifikations-
  Zeiger statt Entscheidungsgrundlage) — `matrix`' erster Zeilen-Marker, eng
  begrenzte Umkehr der „nur strukturelle Ausnahmen"-Haltung. Neues
  `matrix.exempt-paths` überspringt ganze Dateien (Grandfathering immutabler,
  vor Einführung `Accepted`-ADRs). Token in Markdown-Links und Fenced-Code zählen
  nicht; ohne `token`/`exempt-paths` byte-identisch
  ([`DC-FA-MTX-003`](spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix),
  [ADR-0022](docs/plan/adr/0022-matrix-token-richtung-provenance-marker.md),
  Lastenheft 0.31.0).

## [0.30.0] — 2026-06-28

### Added

- slice-050 — Modul `matrix` um die **klasseninterne Verweisrichtung** erweitert
  (`DC-FA-MTX-002`): eine Dokumentklasse kann zusätzlich zu `paths` ein `order`
  (Liste von Pfad-Globs, autoritativste Schicht zuerst; Rang = Index des ersten
  Treffers) und `direction: no-downward` tragen. Ein **klasseninterner** Verweis
  von einer höher- auf eine niederrangige Schicht (auch transitiv) erzeugt den
  neuen Befund `matrix-downward`. Damit ist die Source-Precedence-Schichtung
  *innerhalb* eines Stratums (`architecture → spezifikation → lastenheft`)
  prüfbar, additiv zu den Klassen-Paar-Regeln (`DC-FA-MTX-001`); Globs
  generalisieren auf Spec-Verzeichnisse mit vielen Dateien. Fehlkonfiguration ist
  fail-closed (`order`/`direction` nur gemeinsam, unbekannter `direction`-Wert ⇒
  Konfigurationsfehler); ohne beide Felder ist der Befundsatz byte-identisch
  ([`DC-FA-MTX-002`](spec/lastenheft.md#dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix),
  [ADR-0021](docs/plan/adr/0021-matrix-klasseninterne-verweisrichtung.md),
  Lastenheft 0.30.0).

## [0.29.0] — 2026-06-24

### Added

- slice-049 — neues opt-in Regelmodul `pins` (11. Modul): Content-Pin gegen
  inhaltlichen Drift. Ein Link mit Inline-Marker `<!-- dpin: sha256:<hex> -->`
  (gebunden an den unmittelbar vorausgehenden Link derselben Zeile) wird gegen den
  whitespace-normalisierten **rohen** Ziel-Span (ganze Datei oder Heading-Section,
  inkl. Fenced-Code) gehasht; Drift → Befund `link-stale`. Nur auflösbare,
  repo-interne Ziele (struktureller Befund bleibt `links`/`anchors`, kein
  Doppelbefund); diagnose-only
  ([`DC-FA-PIN-001`](spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in),
  [ADR-0020](docs/plan/adr/0020-content-pin-fence-ausnahme.md), Lastenheft 0.29.0).

## [0.28.0] — 2026-06-24

### Added

- slice-048 — neues opt-in Regelmodul `versions` (zehntes Modul): prüft, dass alle
  gepinnten `ghcr`-Image-Verweise die aktuelle Version aus `version.md#aktuell`
  tragen, sonst Befund `version-stale`; liest die Pins **auch in Fenced-Code**
  (gescopte Fence-Ausnahme), Ventile `exempt-paths`/`d-check:ignore`, fail-closed
  bei unauflösbarer Quelle, diagnose-only
  ([`DC-FA-VER-001`](spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
  [ADR-0019](docs/plan/adr/0019-versions-pin-fence-ausnahme.md), Lastenheft 0.28.0).
- Release-Register `version.md` (only-current-anchor): kanonische Quelle der
  aktuellen Version; `--print-config` führt den `versions:`-Block.

### Changed

- Dogfooding: `.d-check.yml` aktiviert `versions` — die `ghcr`-Image-Pins in
  README und Benutzerhandbuch sind ab jetzt gateguarded.

## [0.27.0] — 2026-06-23

### Added

- slice-047 — `--print-mk`-Fragment um drei Targets + eine Variable erweitert:
  `doc-doctor` (`--doctor`-Diagnose), `doc-repair` (`--repair`-Patch, Recipe-Echo
  unterdrückt für `git apply`-reine stdout), `doc-help` (namespaced Self-Doku der
  `doc-*`-Targets via `##`-Annotationen) sowie `DCHECK_DIGEST` (Digest-Override per
  `ifeq`, sticht den Tag von `DCHECK_IMAGE`). Alle Targets `##`-annotiert
  ([`DC-FA-CLI-010`](spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  Change Request, Lastenheft 0.27.0).

## [0.26.0] — 2026-06-23

### Changed

- slice-046 — `--suggest-config ai-harness[-init]`: die Ausgabe nennt die nicht
  aktivierten situativen opt-in-Module (`external`, `spans`, `hostpaths`,
  `diagrams`) jetzt in einem Kommentar mit Verweis auf `d-check --print-config`
  (Auffindbarkeit ohne Aktivieren eines inerten Moduls — `diagrams` braucht
  repo-spezifische `patterns`/`defined-in`). Schärfung
  ([`DC-FA-CLI-006`](spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten),
  Lastenheft 0.26.0).

## [0.25.0] — 2026-06-23

### Added

- slice-045 — Modul `diagrams` (opt-in): öffnet gezielt benannte
  Diagramm-Fences (Default `mermaid`) und prüft die darin gefundenen Kennungen
  auf **Existenz** in ihrer `defined-in`-Quelle (Befund `diagram-id-undefined`).
  Reine Token-Extraktion ohne Mermaid-Parser, **Existenz statt Link-Policy** (in
  Fences kein Markdown-Link möglich), read-only/netzlos (`DC-QA-03`),
  deterministisch (`DC-QA-02`), Default aus (byte-identisch). Fängt Drift/Typos
  in Diagramm-Kennungen (z. B. Architektur-IDs in `mermaid`), die bei opaken
  Fences heute unsichtbar bleiben
  ([`DC-FA-DIAG-001`](spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in),
  [ADR-0018](docs/plan/adr/0018-diagram-fence-ausnahme.md)).

## [0.24.0] — 2026-06-23

### Added

- slice-044 — Option `--require-complete` (nur mit `--trace`): bindet die
  RTM-Waisen-Markierung an den Exit-Code — ≥1 Requirements-Waise ⇒ **Exit 1**
  statt 0, sonst 0; die RTM bleibt unverändert auf stdout, der Default-`--trace`
  bleibt advisory (Exit 0). Erlaubt Konsumenten ein Vollständigkeits-Gate im
  eigenen Makefile, ohne Parsing-Logik zu kopieren. read-only (`DC-QA-03`),
  deterministisch (`DC-QA-02`)
  ([`DC-FA-CLI-011`](spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)).

### Changed

- slice-044 — `--print-mk`-Fragment: zusätzlich die Targets `doc-trace`
  (advisory RTM) und `doc-complete` (`--trace --require-complete`, das
  Vollständigkeits-Gate) plus eine überschreibbare `TRACE_FLAGS`-Variable
  ([`DC-FA-CLI-010`](spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  Change Request 0.24.0).

## [0.23.0] — 2026-06-22

### Changed

- slice-043 — Modul `codepaths`: neues Ventil `exempt-paths` (Glob-Liste,
  Syntax wie `scan.ignore`) nimmt **ganze Dateien** von der codepath-Prüfung
  aus — Parität zum gleichnamigen `ids`-Ventil; datei-weit, komplementär zum
  `d-check:ignore`-Marker. Abwärtskompatibel: ohne `exempt-paths`
  byte-identisch ([`DC-QA-02`](spec/lastenheft.md#dc-qa-02--determinismus)).
  Dogfooding: die eigene `.d-check.yml` nimmt `docs/reviews/**` aus
  ([`DC-FA-CODE-001`](spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
  Change Request 0.23.0).

## [0.22.0] — 2026-06-22

### Added

- slice-038 — Option `--print-mk`: gibt ein include-bares `d-check.mk`
  (überschreibbare `DCHECK_IMAGE`-Variable mit version-gepinntem Image +
  `doc-check`-Target) auf stdout aus — Konsumenten `include`-n statt
  Recipe/Skript zu kopieren; der Image-Ref ist die ins Binary eingebettete
  Release-Version (Digest via `DCHECK_IMAGE`-Override). read-only
  (`DC-QA-03`), deterministisch (`DC-QA-02`)
  ([`DC-FA-CLI-010`](spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)).
- slice-036 — Option `--trace`: gibt eine **Requirements Traceability
  Matrix** (je Anforderung die referenzierenden ADRs/Slices + Waisen-
  Markierung) auf stdout aus — Default Markdown-Tabelle, mit `--json`/`--yaml`
  maschinenlesbar; read-only (`DC-QA-03`), deterministisch (`DC-QA-02`),
  kein Dokument erzeugt; Doku-Domäne (Lastenheft/ADR/Planning)
  ([`DC-FA-CLI-009`](spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)).

### Changed

- slice-037 — `--suggest-config ai-harness[-init]`: das Anforderungs-`ids`-
  Muster ist nicht mehr fix `DC-`. Neues Flag `--id-prefix <PREFIX>` setzt
  das Präfix explizit; der Modus `ai-harness` leitet es aus
  `spec/lastenheft.md` ab (mehrere verschiedene Präfixe ⇒ Fehler).
  **Breaking:** ohne Präfix/Ableitung (typisch `ai-harness-init` im leeren
  Repo) erscheint ein markierter Platzhalter `<PREFIX>` + `# TODO` statt
  eines stillen `DC-`
  ([`DC-FA-CLI-006`](spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten),
  `ADR-0015`).
- slice-034 — Distribution: der `:latest`-Tag zeigt explizit auf das
  **neueste stabile** Release; Vorabversionen (Prereleases) erhalten kein
  `:latest`. `:latest` ist Komfort-Einstieg — für reproduzierbare Läufe
  weiterhin auf eine feste Version oder den `@sha256:`-Digest pinnen
  (ratifiziert `ADR-0002` §4; `ADR-0014`).

## [0.19.0] — 2026-06-20

### Added

- slice-031 — YAML-Ausgabeformat `--yaml`: gibt die Befunde strukturgleich
  zu `--json` als YAML auf stdout aus (`findings`/`summary`/`exitCode`);
  `--doctor --yaml` analog `--doctor --json`. Deterministisch (`DC-QA-02`),
  read-only (`DC-QA-03`)
  ([`DC-FA-CLI-004`](spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)).

### Changed

- slice-032 — semgrep als hermetisches Security-/Static-Analysis-Gate in
  `make gates`: gepinntes Image (`semgrep/semgrep:1.167.0`) + gepinntes,
  lokal außerhalb des Repos gecachtes `go/lang/security`-Regelset, netzloser
  Scan (`--network none`); `--error` bricht das Gate bei Befund. Ergänzt
  golangci-lint sprachübergreifend; reproduzierbar (`DC-QA-02`), netzlos
  (`DC-QA-03`) (`ADR-0010`).
- slice-033 — alle Build- und Gate-Images per `@sha256:`-Digest gepinnt
  (alle Dockerfile-`FROM` — golang, golangci-lint, distroless — und das
  semgrep-Image; Manifest-Listen-Digest amd64+arm64, inline neben dem Tag);
  `make versions` belegt die Pins. Schließt die `ADR-0002`-§1-Digest-Drift
  und vereinheitlicht die Image-Pin-Politik (`ADR-0011`).

## [0.18.0] — 2026-06-19

### Added

- slice-030 — `--suggest-config ai-harness` / `ai-harness-init`: schlägt
  ein an die ai-harness-course-Konvention (Baseline v1.3.0) angelehntes
  `.d-check.yml` vor — kanonische `ids`-Muster, `matrix`-Klassen samt
  Referenzrichtung, Standard-Modulset und Scan-Scope. **Zwei Modi:**
  `ai-harness-init` (Voll-Kanon, alle Blöcke aktiv — Zielbild fürs leere
  Repo, läuft nach Struktur-Anlage) und `ai-harness` (repo-bewusst — nur
  vorhandene Pfade aktiv, fehlende auskommentiert mit Hinweis). Read-only,
  advisory, deterministisch (`DC-QA-02`); mit echten Quellen kombinierbar
  ([`DC-FA-CLI-006`](spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)).

### Changed

- Spezifikation/Dokumentation präzisiert: `--repair-broad` löst nur
  **Verschiebungen** auf (gleicher Name), keine Umbenennungen; das
  Out-of-Scope von `DC-FA-CLI-008` benennt die nicht-reparierbaren
  Befundarten und schließt git-historienbasierte Move-/Rename-Erkennung
  aus; die Reparatur-Ableitbarkeit ist als Entscheidung in `ADR-0008`
  festgehalten (Handbuch §4.10).

## [0.17.0] — 2026-06-19

### Added

- slice-029 — Maschinenlesbare Diagnose `--doctor --json`: dieselbe
  Diagnose wie `--doctor`, aber als JSON-Dokument auf stdout. Die
  `findings` tragen je Eintrag zusätzlich `reasonText` (Grund-Klartext)
  und `fixCandidate` (`{original, replacement, note}` oder explizit
  `null`, wo kein eindeutiger Fix existiert); `summary`/`exitCode` wie
  bei `--json`. Ein drittes Rendering desselben Fix-Kandidaten-Modells
  neben Prosa (`--doctor`) und Patch (`--repair`); deterministisch
  (`DC-QA-02`), read-only (`DC-QA-03`).

### Changed

- `--doctor` ist nun **mit `--json` kombinierbar** (zuvor Nutzungsfehler).
  Nutzungsfehler bleiben nur `--repair`+`--json` und `--doctor`+`--repair`
  ([`DC-FA-CLI-007`](spec/lastenheft.md#dc-fa-cli-007--diagnose-modus),
  [`DC-FA-CLI-004`](spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)).

## [0.12.0] — 2026-06-18

### Added

- slice-025 — Diagnose-Modus `--doctor`: erklärende, nach Datei
  gruppierte Klartext-Diagnose auf stdout statt der knappen Befund-Zeilen,
  mit Fix-Kandidaten wo eindeutig ableitbar (in dieser Version
  `id-unlinked` → Markdown-Link auf das ids-Definitions-`target`). Grund-
  Klartext für alle 14 Grund-Codes, abgesichert durch eine
  Vollständigkeits-Prüfung gegen die Reason-Konstanten. Read-only,
  stdout-only; `--doctor` ist nicht mit `--json` kombinierbar
  (Nutzungsfehler, Exit 2). Das Fix-Kandidaten-Modell ist die
  wiederverwendbare Eingabe für den folgenden Patch-Modus `--repair`
  (slice-026)
  ([`DC-FA-CLI-007`](spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)).
- slice-026 — Reparatur-Patch `--repair`: gibt einen unified diff auf
  stdout aus (`git apply`-kompatibel), der ableitbare Befunde behebt;
  schreibt selbst nichts. **Konservativ** (Default) nur eindeutige Fixes
  (`id-unlinked` → Definitions-Link, nur nackte Prosa-Vorkommen — keine
  Inline-Code- oder Mehrdeutigkeits-Reparatur); **breit** (`--repair-broad`,
  opt-in) zusätzlich Best-Guess (`target-missing` → Datei eindeutig
  gleichen Basisnamens), review-pflichtig mit Marker auf stderr, sodass
  der Patch `git apply`-rein bleibt. Nicht mit `--json`/`--doctor`
  kombinierbar; deterministisch (`DC-QA-02`), read-only (`DC-QA-03`).
  Wiederverwendung des Fix-Kandidaten-Modells aus slice-025
  ([`DC-FA-CLI-008`](spec/lastenheft.md#dc-fa-cli-008--reparatur-patch)).

### Changed

- slice-027 — `make image-test` deckt nun auch `--doctor`/`--repair` ab
  (nativ == Container byte-identisch, stdout + stderr + Exit-Code); E2E-
  Härtung des
  [`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)-
  Vertrags für die neuen Ausgabe-Modi (`DC-QA-02`-Parität).

## [0.11.0] — 2026-06-17

### Added

- slice-024 — Modul `matrix`: opt-in `allow-supersede-lineage` (mit
  `supersede-fields`) nimmt die Supersede-Lineage-Kante von der
  Status-Prüfung aus — eine ablösende Datei darf auf das von ihr
  abgelöste (inaktive) Dokument verweisen, ohne `matrix-inactive` zu
  erzeugen, sofern sie die Ablösung über ein deklariertes Feld benennt
  (Match über Linktext bzw. Zielpfad der Referenz)
  ([`DC-FA-MTX-001`](spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
  Change Request 0.14.0). Wirkt nur auf `matrix-inactive`;
  `matrix-forbidden` (Klassen-Regeln) bleibt unberührt. `matrix` trägt
  bewusst **keinen** `d-check:ignore`-Marker — legitime Ausnahmen sind
  strukturell (`exclude-sections`, `allow-supersede-lineage`).
  Abwärtskompatibel: Default aus ⇒ Befundsatz byte-identisch.

## [0.10.0] — 2026-06-16

### Changed

- slice-023 — die `ids`-Ventile `exempt-paths` und `d-check:ignore`
  gelten jetzt für **alle** Vorkommen eines Musters — nackt im Fließtext
  **wie** in Inline-Code — und unabhängig von der `link-policy`
  ([`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
  Change Request 0.13.0). Bisher griffen beide nur auf die
  `always`-Inline-Code-Vorkommen; eine nackte Kennung in einer
  `exempt-paths`-Datei (oder auf einer `d-check:ignore`-Zeile) wurde
  weiterhin als `id-unlinked` gemeldet. Jetzt ein Ganzdatei- bzw.
  Ganzzeilen-Carve-out. Abwärtskompatibel: Configs ohne gesetzte Ventile
  sind byte-identisch; die Wirkung geht nur in Richtung *weniger* Befunde
  in explizit ausgenommenen Dateien/Zeilen.

## [0.9.0] — 2026-06-15

### Added

- slice-022 — Inline-HTML-Anker als gültige Anker-Menge
  ([`DC-FA-ANCH-001`](spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
  Schärfung 0.12.0): das Modul `anchors` (und mittelbar `codepaths`)
  akzeptiert zusätzlich zu Heading-Slugs die Inline-HTML-Anker der
  Zieldatei — `id` an beliebigem Element und `name` an `<a>`
  (GitHub-Parität, wörtlicher/case-sensitiver Vergleich). HTML in
  Code-Auszeichnung (Fenced-Block oder Inline-Code) zählt nicht.
  Abwärtskompatibel: reduziert Falsch-Befunde `anchor-missing`, erzeugt
  nie neue.

## [0.8.0] — 2026-06-13

Reichhaltige `--help` (Schärfung `DC-FA-CLI-001`, slice-021): die Hilfe
nennt Synopsis und das Pfad-Argument und verweist fürs Config-Format
auf `--print-config`.

### Changed

- slice-021 — reichhaltige `--help`/`-h`
  ([`DC-FA-CLI-001`](spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
  Schärfung 0.11.0): die Hilfe nennt jetzt eine Kurzbeschreibung, die
  Synopsis `d-check [optionen] [pfad]` und beschreibt das bislang
  verschwiegene Pfad-Argument (Scan-Wurzel, Default cwd); für das
  Config-Format verweist sie auf `--print-config`/`--suggest-config`
  (kein Format-Duplikat). Exit 0 / stderr unverändert.

## [0.7.0] — 2026-06-13

Konfigurations-Vorschlag aus Autoritäts-Dokumenten (`--suggest-config`,
Change Request 0.10.0, slice-020) — inkl. Review R1.

### Added

- slice-020 — Option `--suggest-config`
  ([`DC-FA-CLI-006`](spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten),
  Change Request 0.10.0): liest benannte Autoritäts-Quellen und schlägt
  ein `.d-check.yml` vor — leitet je Quelle ein `ids`-Muster aus den in
  Überschriften **definierten** Kennungen ab (Präfix-Alternation,
  Round-Trip-Garantie; Quell-Kennungen als Kommentar) und schlägt
  opt-in-Module nach Signal vor. **Liest, schreibt nie** (read-only-
  Vertrag; Umleiten macht der Aufrufer). Bewusste Grenze: erkennt nur
  großgeschriebene Heading-Token-IDs — Scaffold, kein Orakel, der
  Mensch verengt/ergänzt. Dazu Schärfung von
  [`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids):
  Muster-Ableitung bleibt für die **Prüfung** Out-of-Scope, nur der
  advisory Scaffold-Modus leitet ab. Korpora-Gegentest dokumentiert
  (d-check round-trippt, b-trace zeigt die Heading-Grenze).

### Fixed

- Review R1 zu slice-020 (`/code-review`): das `--suggest-config`-Gerüst
  nimmt jetzt `ids` in die Modul-Liste auf (sonst waren die abgeleiteten
  Muster im erzeugten Config inaktiv — gültiges YAML, semantisch
  wirkungslos); Modul-Probelauf nutzt denselben Scope (`roots: ["."]`)
  wie das Gerüst; `target` wird gequotet (Quellpfade mit `:`/`#` brechen
  das YAML nicht mehr); Heading-Token-Extraktion strippt Links und
  Satzzeichen (`ADR-0001:` wird erkannt); leere Quellenliste ist ein
  Nutzungsfehler. Report unter `docs/reviews/`.

## [0.6.0] — 2026-06-13

Konfigurations-Startgerüst (`--print-config`, Change Request 0.9.0,
slice-019): neue Repos kommen ohne Handarbeit zu einer `.d-check.yml` —
das Werkzeug gibt aus, der Aufrufer leitet um; read-only bleibt.

### Added

- slice-019 — Option `--print-config`
  ([`DC-FA-CLI-005`](spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben),
  Change Request 0.9.0): gibt ein kommentiertes `.d-check.yml`-Startgerüst
  auf stdout aus und endet mit Exit 0 — **kein Scan, schreibt nichts**
  (read-only-Vertrag bleibt; Anlegen via `d-check --print-config >
  .d-check.yml`). Das Gerüst ist statisch, deterministisch und
  dekodiert über den eigenen Parser; es macht die verfügbaren Module
  und Optionen als Kommentare sichtbar. Senkt die Adoptions-Reibung in
  neuen Repos ohne Konfiguration.

## [0.5.0] — 2026-06-13

Konfigurierbare Link-Politik für das Modul `ids` (Change Request 0.8.0,
slice-018): „gut verlinkte Dokumente" wird ein im `.d-check.yml`
konfigurierbares, gemessenes Property.

### Added

- slice-018 — konfigurierbare Link-Politik `ids.patterns[].link-policy`
  ([`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
  Change Request 0.8.0): `link-policy: always` macht auch
  Inline-Code-Vorkommen einer Kennung linkpflichtig — „gut verlinkt"
  als gemessenes, konfigurierbares Property statt menschlicher
  Aufmerksamkeit. Default `prose` (byte-identisch, opt-in fürs Gating).
  Zwei Ventile: `exempt-paths` (Glob-Liste je Muster) und der
  Zeilen-Marker `d-check:ignore` (Geltungsbereich von `codepaths` auf
  `ids` erweitert — illustrative Beispiel-IDs). Kalibriert über die
  drei ids-Repos (d-check, u-boot, b-trace); Dogfooding aktiv (d-check
  setzt `always` und verlinkte den eigenen Befundsatz). Nutzersichtbar
  dokumentiert in [`docs/user/operations.md`](docs/user/operations.md).

### Changed

- Dogfooding-Sweep: d-checks eigene Doku auf `link-policy: always`
  umgestellt; Inline-Code-Kennungen in Slices, ADRs, AGENTS, harness
  und Spezifikation als Links ausgeführt (Sektions-Referenzen `.a` auf
  ihre Spez-Anker), `exempt-paths` für CHANGELOG + Reviews,
  Beispiel-IDs mit `d-check:ignore`.

## [0.4.0] — 2026-06-13

Welle-06-sensorik: zwei opt-in-Sensormodule — `spans` (`DC-FA-SPAN-001`,
Markdown-Span-Artefakte) und `hostpaths` (`DC-FA-HOST-001`, host-lokale
absolute Pfade), je über 14 Korpora kalibriert und gegen die Alt-Sensoren
gegengeprüft; schließt welle-06 (slice-015, slice-016).

### Added

- slice-016 — Modul `hostpaths`
  ([`DC-FA-HOST-001`](spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in),
  Change Request 0.6.0, opt-in): meldet host-lokale absolute Pfade
  (Maschinen-Layout-Leaks) in Prosa **und Inline-Code**
  (`hostpath-forbidden`); Unix-Präfixliste konfigurierbar via
  `hostpaths.prefixes` (Default ohne tmp — Lastenheft 0.7.2 aus dem
  Kalibrierungs-Befund: Laufzeit-Doku ist legitim),
  Windows-/UNC-Muster fest (UNC-Servername alphanumerisch — Schutz
  vor Regex-Beispiel-Treffern), Fences ausgenommen, kein
  Opt-out-Marker. Paritäts-Gegentest gegen den
  bess-ems-Rest-Sensor: identische Befunde auf identischen Zeilen;
  Kalibrierung über 14 Korpora trennte echte Workspace-Leaks
  (k-deskflight-Spec gefixt) von gewollter
  Windows-/WSL-Plattform-Doku (Opt-in-Entscheidung der Repos).

- slice-015 — Modul `spans`
  ([`DC-FA-SPAN-001`](spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in),
  Change Request 0.5.0, opt-in): meldet ungeschlossene Code-Spans,
  die an Nicht-Whitespace kleben (`span-unclosed`, absatzweise —
  alleinstehende Backticks bleiben literal) und Link-Syntax im
  Linktext (`span-nested-link`; Badge-Muster `[![…](…)](…)` ist
  legal — Lastenheft 0.7.1 aus dem Kalibrierungs-Befund). Dogfooding
  aktiv; Kalibrierung über 14 Korpora fand 17 echte Artefakte (in
  den Ziel-Repos gefixt), historischer Gegentest: 14 Befunde auf dem
  u-boot-Stand vor den slice-014-Reparaturen, 0 danach.

## [0.3.0] — 2026-06-12

Modul-lokaler Scan-Scope (`<modul>.scope`, Change Request des
Erst-Bedarfsträgers grid-gym) — dazu der dokumentierte Abschluss des
13/13-Migrations-Rollouts.

### Added

- slice-017 — Modul-lokaler Scan-Scope
  ([`DC-FA-CONF-002`](spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope),
  Change Request des Erst-Bedarfsträgers grid-gym): optionaler
  Schlüssel `<modul>.scope` (`roots` Pflicht, `ignore` optional)
  ersetzt für genau dieses Modul den globalen Scan-Scope — eigener
  Discover-Lauf mit den bekannten Scan-Regeln, Lauf über die
  Vereinigungsmenge mit Einmal-Lese-Garantie; ohne `scope`
  byte-identisches Verhalten (belegt gegen v0.2.1 auf zwei Korpora).
  Konsumenten-Abnahme grid-gym: `ids` kuratiert auf `spec/` +
  `docs/user/` → 311 statt 2378 Befunde, `links`/`anchors`
  unverändert global.

- slice-014 — Rollout abgeschlossen
  ([`DC-QA-04`](spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
  vollständig, 13/13): alle verbleibenden neun Alt-Tool-Vorkommen
  migriert (Shell: b-trace, m-trace, cmake-xray, k-deskflight;
  Python: c-hsm-doc, grid-gym; JS: euler-fourier-hilbert, b-cad;
  eigenständige Linie: bess-ems — Inventur-Nachtrag Lastenheft
  0.4.0). 16 echte Mehr-Befunde in den Ziel-Repos gefixt;
  Rest-Sensoren für Math-Validierung, Host-Pfad-Prosa und
  Modul-Nummern verbleiben dort. Zusatz: Neu-Adoption pkcs11-course
  (Auslöser der v0.2.1-Scan-Härtung). Vergleichstabellen in der
  slice-014-Closure-Notiz; schließt welle-05.

## [0.2.1] — 2026-06-12

Scan-Härtung aus der pkcs11-course-Adoption (slice-014) plus der
dokumentierte Rollout-Stand.

### Fixed

- `scan.ignore`-Muster prunen jetzt den Verzeichnis-Abstieg: ein
  vollständig ignorierter Teilbaum (`pfad/**` oder direkt matchendes
  Muster) wird nicht betreten — unlesbare ignorierte Verzeichnisse
  (z. B. root-eigene Laufzeit-Residuen wie SoftHSM-Tokens) brechen
  den Lauf nicht mehr mit Exit 2 ab
  ([`DC-FA-SCAN-001`](spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln)).
- `SKIP_DIRS` um `.gradle` ergänzt (Parität zur JS-Alt-Familie);
  Spezifikation §3 inkl. Querverweis aus dem Config-Schema.

### Added

- slice-012 — Pilot-Migrationen
  ([`DC-QA-04`](spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)):
  drei Schwester-Repos prüfen ihre Doku jetzt digest-gepinnt über das
  GHCR-Image — d-migrate (Shell-Familie, v0.1.0-Digest, Alt-Skript
  gelöscht), ai-harness-course (JS-Familie, v0.2.0-Digest,
  `docs-check.js` auf den Modul-Nummern-Rest-Sensor geschrumpft) und
  u-boot (Python-Vollausbau, v0.2.0-Digest, `check_refs.py`
  deprecated). Vergleichsläufe und Triage in der Closure-Notiz des
  Slices; schließt welle-04 und Meilenstein M3.

## [0.2.0] — 2026-06-12

Modul `codepaths` (Change Request 0.3.0) und der absatzweise
Inline-Code-Parser aus dem `DC-QA-04`-Gegentest — damit enthält das
Image alle sechs Regelmodule (slice-013, slice-012-Vorlauf).

### Added

- slice-013 — Modul `codepaths` (`DC-FA-CODE-001`, opt-in): explizite
  Pfade in Inline-Code (`./`, `../`, konfigurierte Wurzel-Präfixe via <!-- d-check:ignore (Syntax-Beispiel) -->
  `codepaths.roots`) werden auf Existenz, Repo-Escape und — bei
  Markdown-Zielen mit Fragment — Anker geprüft; Wert-Normalisierung
  (Anführungszeichen, Satzzeichen, `Datei:Zeile`-Suffix), Headings
  ausgenommen, Zeilen-Opt-out `<!-- d-check:ignore (Begründung) -->`
  wirkt nur auf dieses Modul. Dogfooding aktiv (eigene Doku
  befundfrei; 16 begründete Marker an historischen/Beispiel-Pfaden).
- Review-Infrastruktur nach Digest-Welle-18-Konvention:
  Reviewer-Skill (`.harness/skills/reviewer.md` — Kategorien-Anker,
  Output-Schema, Negativbefund-Pflicht) und Report-Ablage
  `docs/reviews/` (ein Report pro Lauf); `.gitignore` auf
  `.harness/state/` verengt (Skills sind versionierte
  Harness-Mechanik).

### Changed

- Lastenheft 0.3.0/0.3.1 — Change Request `DC-FA-CODE-001`: neues
  opt-in-Modul `codepaths` (explizite Pfade in Inline-Code,
  Zeilen-Opt-out `d-check:ignore` nur für dieses Modul) inkl.
  Review-R1-Präzisierungen (Wert-Normalisierung, Anker-Prüfung bei
  Markdown-Zielen); Umsetzung folgt mit slice-013.

### Fixed

- Inline-Code-Erkennung absatzweise statt zeilenweise (CommonMark:
  Spans dürfen Zeilenumbrüche enthalten; Absatzgrenzen sind
  Leerzeile/Fence, ungeschlossene Backtick-Folgen sind literal und
  brechen den Scan nicht mehr ab). Behebt False-Positive-
  `id-unlinked`-Befunde auf korrekt verlinkten Kennungen nach
  Span-Fortsetzungszeilen — gefunden im
  [`DC-QA-04`](spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Gegentest
  gegen den u-boot-Vollausbau (slice-012); Spezifikation
  §`DC-FA-LINK-001.a` Schritt 2 fortgeschrieben.

## [0.1.0] — 2026-06-11

Erster Release: alle fünf Regelmodule, Gate-Vollausbau, Distribution
via `ghcr.io/pt9912/d-check` (slice-011 — `DC-FA-DIST-001`).

### Added

- GHCR-Release-Pipeline
  ([`.github/workflows/release.yml`](.github/workflows/release.yml)):
  Tag-Push `v*` → SemVer-Validate → `make ci` → OCI-Label-Pin
  (`org.opencontainers.image.version` muss der Tag-Version
  entsprechen) → Push mit Semver-Tag (+ `latest` nur für stabile
  Releases) → Digest-Pin im Job-Summary und GitHub-Release;
  Betriebs-/Release-Doku unter `docs/user/` (löst `MR-009`:
  `docs/user` ist jetzt Rang 6 der Source Precedence).
- MIT-Lizenz als Repository-Lizenz ergänzt ([`LICENSE`](LICENSE)).
- Lastenheft 0.1.0 ([`spec/lastenheft.md`](spec/lastenheft.md)):
  Konsolidierung von 12 Quell-Tools, Modul-Schnitt (`links`, `anchors`,
  `ids`, `matrix`, `external`), Docker-Distribution.
- Greenfield-Harness-Bootstrap nach AI-Harness-Kurs: `AGENTS.md`,
  `harness/`, Planning-Struktur mit Roadmap und Slices 001–004,
  Makefile-Gates (`doc-check`, `gates`), `.claude`-Hooks.
- Fundament-ADRs 0001–0004 (slice-001): Go, GHCR-Image auf
  distroless/static (Binary-Distribution vertagt mit Trigger),
  striktes YAML via yaml.v3, Architektur Hexagon light.
- Spec-Straten 2+3 (slice-002): `spec/spezifikation.md` (Prüflauf-,
  Slug-, Modul-Algorithmen; `--json`- und `.d-check.yml`-Schema;
  Defaults; Grund-Codes) und `spec/architecture.md` (Hexagon-Schnitt,
  Import-Constraints als arch-check-Grundlage, Sequenzen).

- slice-003 — erster Go-Code: CLI-Kern (`d-check [pfad]`,
  `--enable/--disable/--json`), Scanner mit Default-Wurzeln/Ignores,
  Modul `links` (Linkziele, Repo-Escape, Symlink-Vorrang,
  RFC-3986-Dekodierung), strikte `.d-check.yml`-Validierung,
  Text-/JSON-Reporter; Layout nach `ADR-0005` (hexagon-/adapter-Ordner,
  u-boot-Konvention); Dockerfile-Stages + Make-Gates `lint`, `test`,
  `arch-check` (Fitness Function R1–R5), Runtime-Image
  distroless/static mit Selbst-Smoke-Test (`make run`).

- slice-004 — Modul `anchors` (GitHub-Slug-Verfahren inkl.
  Duplikat-Suffixen, Fragment-Dekodierung, Schweigen bei fehlender
  Zieldatei) und **Dogfooding**: `make doc-check` läuft über `d-check`
  selbst (`scan.roots: ["."]`, Module links+anchors — erstmals mit
  Anker-Validierung); vendorter Bootstrap-Sensor gelöscht
  (`MR-003` → `MR-007` aufgelöst); Vergleichslauf als erster
  `DC-QA-04`-Datenpunkt.

- slice-005 — SOLID-nahes Lint-Profil (`ADR-0006`, u-boot-Parität ohne
  depguard): 5 Default- + 23 Linter mit u-boot-Kalibrierung,
  gomodguard-Anti-Module, Why-kommentierte Ausnahmen; Code-Refactoring
  statt Carveouts (Globals → Funktionen, Komplexitäts-Splits in
  cli/configyaml/core) — lint-clean ohne //nolint.

- slice-006 — Modul `ids` (`DC-FA-ID-001`): Linkpflicht für Kennungen
  nach konfigurierten Regex-Mustern (Reihenfolge = Präzedenz, erstes
  Match gewinnt pro Vorkommen); „verlinkt" = Vorkommen im Linktext
  eines Markdown-Links, Ziel-Klammern und Bildreferenzen sind
  linkpflichtfrei (kein Fließtext); Grund-Code `id-unlinked`;
  Config-Constraint `ids.patterns[].target` muss existieren (Exit 2).

- slice-007 — Modul `matrix` (`DC-FA-MTX-001`): Dokumentklassen per
  Glob (Reihenfolge = Präzedenz), Referenzregeln pro Klassen-Paar,
  Status-Bedingungen (`**Status:**`-Zeile vor Status-Heading,
  Präfix-Match case-insensitiv, ohne Status aktiv) und
  `exclude-sections` (Provenance-Ausnahme); Grund-Codes
  `matrix-forbidden`/`matrix-inactive`. **Dogfooding-
  Selbstkonfiguration:** die eigene `.d-check.yml` aktiviert
  `ids` + `matrix` (Muster `ADR-*`/`MR-*`/`DC-*`; `MR-006`-
  Referenzrichtung maschinell kodiert); ids-Fortschreibung aus dem
  Selbstlauf: Headings und Definitions-Ort des Musters sind
  linkpflichtfrei; ~50 nackte Kennungen der eigenen Doku verlinkt
  bzw. als Code-Span fixiert.

- slice-008 — Modul `external` (`DC-FA-EXT-001`, opt-in): HTTP-Port im
  Hexagon + `httpcheck`-Adapter (HEAD mit GET-Fallback bei 405/501,
  Timeout konfigurierbar, Redirect-Limit 5, begrenzte Parallelität,
  eine Prüfung pro URL); Grund-Codes
  `external-status`/`external-timeout`/`external-redirects`;
  `make doc-check` läuft jetzt mit `--network none` und ist damit die
  automatisierte `DC-QA-03`-Messmethode (netzloser Lauf aller Module
  außer `external`); Interim-Mechanismus
  `isImplemented`/`SkippedModules` entfernt — alle fünf
  Vertragsmodule sind lauffähig.

- slice-010 — Image-Integrationstests + Reproduzierbarkeits-Belege:
  `make image-test` automatisiert die `DC-FA-DIST-001`-Akzeptanzkriterien
  gegen das lokal gebaute Image (Befund-Ausgabe und Exit-Code nativ
  vs. Container byte-identisch, read-only-Mount, fehlender Mount →
  Exit 2 mit Hinweis); `make versions` (Pins + Runtime-Image-ID),
  `make ci` (gates + image-test — Target der Release-Pipeline) und
  `make fullbuild` (volle Closure inkl. Benchmark, schließt mit dem
  Image-Hash); die „Nicht behauptet"-Listen sind leer.
- slice-009 — Gate-Endausbau (welle-03-Abschluss): `make coverage-gate`
  (Coverage über `./internal/...` per `-coverpkg`, Kalibrierungs-Bindung
  85 % → 90 % bei welle-03 done; Ist-Stand 92,9 %),
  `make gate-consistency` (Meta-Gate gegen Harness-Lügen: dokumentierte
  Targets ↔ Makefile in beide Richtungen, `DC-QA-03`-Modulliste der
  `.d-check.yml`, Selbsttest mit Phantom-Target bei jedem Lauf) — beide
  in `make gates` aggregiert; `make bench` mit `DC-QA-01`-Benchmark
  (Spez §`DC-QA-01.a` eingelöst: 1.000-Dateien-Fixture, Median aus 3
  Container-Läufen; gemessen: 551 ms ≪ 5 s).

### Changed

- Review R1 zu slice-009/010 (Gate-Infrastruktur): Benchmark-Median
  aus `RUNS` abgeleitet statt hart verdrahtet (latente
  Stilles-Grün-Falle); `--cpus 2` im Benchmark-Lauf + Spez-Präzisierung
  (2-vCPU-Normierung aus `DC-QA-01`); Meta-Gate-Parser erkennt
  Mehrfach-Target-Zeilen und schließt Variablen-Zuweisungen aus (mit
  Parser-Selbsttest); `fullbuild: ci bench` statt Kettenduplikat;
  `make versions` ohne Stage-FROM-Rauschen; drei dokumentierte
  Annahmen (QA-03-Listenformat fail-closed, amd64-Kopplung des
  image-test, Bench-Fixture bleibt zur Inspektion liegen);
  93-%-Kalibrierung als Nachtrag in der slice-009-Closure-Notiz.

- Review R1 zu slice-008: Fragmente werden vor Prüfung/Dedupe entfernt
  (eine Prüfung pro Ressource, Befund am Original-Linkziel);
  Schema-Vergleich case-insensitiv (kein stiller Gap zwischen `links`
  und `external` bei `HTTP://`); explizit gesetzte 0 in
  `external.timeout-seconds`/`parallel` ist Konfigurationsfehler statt
  stillem Default; GET-Fallback drained den Body begrenzt (64 KB);
  HTTP-Adapter wird nur noch bei aktivem `external` verdrahtet
  (strukturelle Opt-in-Absicherung); Timeout-Semantik (pro Request)
  spezifiziert; QA-03-Config-Kopplung als gate-consistency-Auftrag in
  slice-009 eingetragen.
- Review R1 zu slice-007: Status-Extraktion fence-aware (Fence-Inhalt
  ist kein Statuswert) und nur für Markdown-Ziele (kein Voll-Read von
  Binärdateien); Doppel-Befund forbidden+inactive als unabhängige
  Verletzungen spezifiziert; gemeinsamer Fence-Scanner (`proseLines`)
  und gemeinsame Ziel-Auflösung (`localTarget`) statt Drittkopien;
  `exclude-sections` der Selbstkonfiguration um die realen
  nummerierten Headings („7. Historie") ergänzt.
- Review R1 zu slice-006: Config-Pfade (`scan.roots`,
  `ids.patterns[].target`) dürfen die Repo-Wurzel nicht verlassen
  (Exit 2 statt stillem Escape); Leerstring-matchende ids-Regexe sind
  Konfigurationsfehler; Inline-Code-Stripping positionserhaltend
  (keine Phantom-Kennungen durch Text-Verschmelzung); zeilenbasierte
  Link-Extraktion als normative Grenze dokumentiert.
- Harness-Hooks gehärtet (`MR-005`): Gate-Nachweis inhaltsbasiert
  (Commit ohne Gate-Lauf wird vom Stop-Hook nicht mehr freigegeben),
  PreToolUse-Guard prüft `bash/sh -c`-Sub-Shell-Strings rekursiv.
- Referenzrichtungs-Korrektur (`MR-006`): ADR-Abwärtsverweise aus
  `spec/spezifikation.md` und `spec/architecture.md` entfernt
  (Kurs-Template-Fehler; Spec-Straten verweisen nie abwärts,
  Traceability über die `Schärft:`-Felder der ADRs).
- `spec/architecture.md` sprachneutral umformuliert (Schichten/Rollen
  statt Modul-Pfade und Imports; sprachkonkrete Übersetzung lebt in
  `ADR-0004`) — Template-Hard-Rule „sprach- und meilensteinfrei" wieder
  voll erfüllt.
- Lastenheft 0.2.1 (redaktionell): Beispiel-Kennungen auf fiktive
  Nummern (`ADR-0042`, `ADR-0099`) — keine Kollision mit real
  entstandenen eigenen ADRs.
