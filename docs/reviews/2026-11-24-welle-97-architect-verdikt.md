# Architect-Verdikt — welle-97 (`linelength`, eingerückte Tabellenzeilen), Plan gegen Entscheidungslage

**Rolle:** Architect (Modul 8), frischer Kontext, **vor** Code. Rolleninhaber: pt9912.
**Gegenstand:** [welle-97](../plan/planning/done/welle-97-linelength-eingerueckte-tabellen.md) und ihre zwei
Slice-Pläne [slice-293](../plan/planning/done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md)
und [slice-294](../plan/planning/done/slice-294-linelength-eingerueckte-tabellen-kern.md), Stand `c9379511`.
Die Fragen 1 bis 3 der Welle (§5) hat der Auftraggeber am 2026-11-24 bestätigt, Frage 3 als „höchstens drei
Leerzeichen, kein Listenkontext“.
**Anlass:** [Issue #20](https://github.com/pt9912/d-check/issues/20).
**Datum:** 2026-11-24.
**Geändert wurde kein Plan, keine Spec, kein Code.** Mit diesem Verdikt entsteht die Folge-ADR
[ADR-0116](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) (`Proposed`) samt Index-Zeile. Sie trägt
die Entscheide 1 bis 3 wie bestätigt und Entscheid 4 aus diesem Verdikt.

## Gesamturteil

Die Welle ist klein geschnitten und trägt: vier Entscheide, sechs Kriterien, keine neue Konfiguration. **Frage 4
ist entschieden**, wie der Plan empfiehlt, mit einer Schärfung (F-1). Der Helfer wohnt in der geteilten Lexik und
**schließt `tableRowLine` ein**; `linelength` ruft nur ihn. Kein anderes Modul bewegt sich.

**Ein Befund bricht den Closure-Trigger der Welle, wie er steht** (F-2): Die Probe gegen `v0.92.0` kann nicht
„byte-identisch“ enden, weil das Vorher-Image das eigene Profil seit welle-96 ablehnt.

Die übrigen Befunde sind Text, alle **V1**:
- eine Stelle des Lastenhefts, die der Plan nicht als berührt führt (F-3);
- drei Grenzen ohne Text (F-4);
- der falsche Re-Evaluierungs-Trigger von ADR-0114 (F-5);
- der Status-Nachzug von ADR-0114 (F-6).

**Keine Entscheidung des Auftraggebers nötig.** Entscheid 4 ist die Antwort des Architect auf Frage 4 und
bleibt im Rahmen der bestätigten Entscheide 1 bis 3. Lastenheft §5 bleibt wahr, solange der Helfer in der
geteilten Lexik wohnt (F-3). Erst wenn er in `linelength` wanderte, wäre §5 berührt, und das wäre eine
Entscheidung des Auftraggebers.

Verdikt-Kürzel (Modul 8, Konflikt-Pfad): **V1** = Entscheidung gilt, Plan hat falsch behauptet/ausgelassen ·
**V2** = Folge-Entscheidung nötig · **V3** = Lockerung legitim, aber undokumentiert.

## Befunde

### F-1 — Frage 4: ein Helfer in der geteilten Lexik, der `tableRowLine` einschließt — **V2** (ADR-0116 Entscheid 4), und der Plan ruft ihn „zusätzlich“ — **V1**

- **Wer sich bei der Weitung bewegte**, am Gegenstand gezählt (`grep -n tableRowLine internal/hexagon/core/rules/*.go`
  ohne Tests):
  - `targets.go` Z. 246;
  - `planning_waves.go` Z. 277;
  - `structure_tablecell.go` Z. 43 und 93, `structure_tableorder.go` Z. 129;
  - `tableHeaderOrSeparator` in `markdown.go` Z. 205 (Konsumenten `planning.waves` und `structure`);
  - `linelength.go` Z. 72.

  `targets` hat in der Spezifikation die Zusage „erstes Zeichen ein Pipe, Spalte 0“
  ([spezifikation.md](../../spec/spezifikation.md) Z. 3040 und Schema Z. 4121). Dazu prüft `tableSeparatorRow`
  (`markdown.go` Z. 208–211) das führende `|` selbst an der rohen Zeile. Eine Weitung von `tableRowLine`
  verschöbe drei Module gegen ihre Zusagen und bliebe halb, wenn sie die Trennzeile nicht mitnähme.
  **Entscheid: nicht weiten.**
- **Ein Helfer daneben** ist nach [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 2 eine **andere Frage**, gescopt durch eine ADR. `linelength` fragt „ist die Zeile Teil einer
  Tabelle und deshalb nicht umbrechbar“, die drei anderen „deklariert die Zeile etwas in einer Spalte-0-Tabelle“.
  Entscheidung 1 verbietet nur, dass ein **Konsument** die Frage selbst beantwortet. Deshalb wohnt der Helfer
  in der geteilten Lexik (`markdown.go`), nicht in `linelength.go`.
- **Die Schärfung:** Der Helfer **ruft** `tableRowLine` und erweitert sie um null bis drei Leerzeichen davor,
  statt die Spalte-0-Antwort ein zweites Mal zu schreiben. `linelength` ruft dann **nur** den Helfer.
  - slice-294 §3 Z. 108 sagt dagegen, `linelength` rufe ihn „zusätzlich zu `tableRowLine`“.
  - Mit zwei Aufrufen bliebe die Mutation zu Kriterium 6 (Z. 131, „der Helfer trifft die Spalte 0 nicht mehr“)
    grün, weil `tableRowLine` die Spalte 0 weiter trifft. Der Test führe seine Regel nicht
    ([`BEO-ALL/testfall-faehrt-seine-regel-nicht`](../plan/planning/observations/BEO-ALL/testfall-faehrt-seine-regel-nicht/observation.md)).
- **Kopplung** (ADR-0054 Entscheidung 4): `TestTabellenzeilenFrageHatEineAntwort`
  ([lexikon_kopplung_test.go](../../internal/hexagon/core/rules/lexikon_kopplung_test.go) Z. 207–227) hält heute
  alle vier Konsumenten auf einer Antwort, auch für „eingerueckt“ mit zwei Leerzeichen (Z. 134: `false`). Künftig
  gilt:
  - die drei Spalte-0-Konsumenten weiter auf einer Antwort, die eingerückte Zeile eingeschlossen;
  - `linelength` auf der Antwort des Helfers, in beiden Richtungen.

  Der Kopf-Kommentar Z. 107–113 („… GLEICH beantwortet …, eine eingerueckte Zeile ist keine
  Spalte-0-Tabellenzeile“) ändert sich mit. slice-294 §3 Z. 110 sagt das im Kern; die zwei Punkte gehören
  ausgeschrieben dazu.
- **Plan-Änderung:**
  - slice-294 §3 Z. 107: „der Helfer ruft `tableRowLine` und erweitert sie um null bis drei Leerzeichen“.
  - Z. 108: „`linelength` ruft für (b) nur den Helfer“.
  - Mutation zu Kriterium 6 (Z. 131): „der Helfer ruft `tableRowLine` nicht mehr“.
  - Z. 110: der Kopf-Kommentar des Kopplungs-Tests wird mit angepasst.
  - Welle §5 Frage 4 nennt das Verdikt und ADR-0116 statt der Empfehlung.

### F-2 — Der Closure-Trigger (a1) kann gegen `v0.92.0` nicht „byte-identisch“ enden — **V1**

- Welle §3 Z. 67: „`make blackbox-probe REF=v0.92.0` endet mit ‚byte-identisch‘“. Das Profil dieses Repos schaltet
  seit welle-96 `vcs.follow-renames` ein. Das Vorher-Image `v0.92.0` kennt den Schlüssel nicht, liest die Datei
  strikt und lehnt sie ab. Der Eigenlauf `repo` endet vorher mit Exit 2, nachher mit Exit 0. Gemessen in der
  [Verifikation zu welle-96](2026-11-24-welle-96-verifikation.md) Z. 36–41: „4 von 20 Vergleichen weichen ab
  (Vorher 61a412eb)“. Diese Abweichung kommt aus welle-96, nicht aus dieser Welle.
- Der billigste Weg, der trägt: die Vorher-Seite auf den letzten Stand **vor** dieser Welle setzen, nicht auf
  das letzte Release. Das ist der Self-Close von welle-96, `69c53e0f`. Das Target nimmt jede Revision
  (`git archive`). Dasselbe gilt für (a2) (Z. 76–77, „`git archive v0.92.0`“): mit derselben Vorher-Seite misst
  der Vergleich genau diese Welle.
- **Plan-Änderung:** Welle §3 (a1) `REF=69c53e0f` (Self-Close von welle-96) statt `REF=v0.92.0`, (a2) dasselbe
  Vorher-Image; ein Halbsatz zur Begründung (welle-96 hat das Profil geändert, `v0.92.0` lehnt es ab).

### F-3 — Lastenheft §5 trägt „mit derselben Markdown-Lexik“; der Plan führt die Stelle nicht — **V1**

- §5 ([lastenheft.md](../../spec/lastenheft.md) Z. 4765) nimmt die Zeilenlänge vom globalen Lint-Ausschluss aus,
  weil sie „mit derselben Markdown-Lexik wie die übrigen Module, die Markdown lesen“ läuft. Das war eine
  Entscheidung des Auftraggebers ([ADR-0114](../plan/adr/0114-linelength-modul.md) Kontext).
- slice-293 Kopf „Berührte Spec-Stellen“ (Z. 21–28) und §3 nennen §5 nicht. Die Klasse ist
  [`BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft`](../plan/planning/observations/BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft/observation.md).
- Mit Entscheid 4 bleibt §5 wahr: Die Antwort kommt aus der geteilten Lexik, auf eine eigene, benannte Frage.
  Sie bliebe es nicht mit einer privaten Erkennung in `linelength` (ADR-0116, Alternative F).
- **Plan-Änderung:** slice-293 §3 ein Satz: „Lastenheft §5 geprüft, unverändert: die Antwort kommt aus der
  geteilten Lexik (ADR-0116 Entscheid 4)“. Dazu §1 als ausdrücklich nicht berührt.

### F-4 — Drei Grenzen der Regel stehen in keinem Plantext — **V1**

Die Welle hat zwei davon, das Lastenheft und die Kriterien in slice-293 §3 (Z. 131–138) keine. Je mit Bedingung,
alle drei in ADR-0116 §Konsequenzen:

1. **Prosa-Zeile mit `|`.** Eine Trennzeile wird weiter nicht verlangt. Beginnt eine Zeile nach ein bis drei
   Leerzeichen mit `|`, ohne dass eine Tabelle dasteht (etwa eine Absatz-Fortsetzung), ist sie ausgenommen, wie
   heute in Spalte 0 ([lastenheft.md](../../spec/lastenheft.md) Z. 4634–4635).
2. **Tabelle auf der Zeile der Listenmarke** (`- | a | b |`): Sie beginnt nicht mit Leerzeichen und `|` und bleibt
   gemessen.
3. **Marke von vier Spalten oder verschachtelte Liste** (`10. `): vier und mehr Leerzeichen, bleibt gemessen. Die
   Welle nennt das (§5 Frage 3), die Anforderung muss es tragen.

Das sind Sätze im Punkt „Tabellenzeile heißt …“ (Z. 4634–4638), keine neuen Kriterien. Kriterium 2 deckt die
Schwelle bei vier schon.

- **Plan-Änderung:** slice-293 DoD-Punkt 1 (Z. 78–85): der Punkt „Tabellenzeile heißt …“ nennt die drei Grenzen
  mit ihrer Bedingung; slice-294 §3 bekommt einen Testfall zu 1 (Prosa-Zeile mit drei Leerzeichen und `|`,
  ausgenommen) und zu 2.

### F-5 — Der fällige Trigger von ADR-0114 ist (b), nicht (a) — **V1**

- Welle §3 (c), Z. 94–98, nennt Re-Evaluierungs-Trigger (a) von ADR-0114 fällig: „die geteilte Lexik ihre …
  Tabellen-Antwort ändert“.
- Mit Entscheid 4 ändert sie sie **nicht**; `tableRowLine` bleibt zeichengleich, (a) tritt nicht ein.
- Was das Issue auslöst, ist (b): „ein Konsument eine Zeile ausnehmen muss, die weder Code noch Tabelle noch
  unumbrechbar ist“. Die eingerückte Tabellenzeile ist nach ADR-0114 Entscheidung 5 (b) keine Tabelle. ADR-0114
  verweist dafür auf den Marker; ADR-0116 vergleicht ihn als Alternative G und verwirft ihn wie
  ADR-0114 Entscheidung 9.
- **Plan-Änderung:** Welle §3 (c): Trigger (b), nicht (a); ADR-0116 als Antwort. slice-293 Bezug (Z. 16–17)
  bleibt.

### F-6 — Der Nachzug an ADR-0114 ist nicht geplant — **V1**

- Mit dem Übergang von ADR-0116 nach `Accepted` bekommt ADR-0114 eine Zeile in `## Geschichte` („Entscheidung 5 (b)
  teil-superseded durch ADR-0116“). Ihre Zeile im Index führt „Accepted (teil-superseded: ADR-0116)“ wie bei
  ADR-0002 und ADR-0013.
- Die Status-Zeile **bleibt** `Accepted`. `vcs.head-allow` lässt nur `Accepted` oder `Superseded by ADR-NNNN` zu,
  und die Ablösung ist nur teilweise.
- slice-293 §1 Z. 60–63 überlässt das dem Verdikt; slice-294 DoD-Punkt 3 (Z. 85–88) nennt nur ADR-0116.
- **Plan-Änderung:** slice-294 DoD-Punkt 3 und §3 Z. 114: „im selben Commit wie ADR-0116 `Accepted`:
  Geschichte-Zeile in ADR-0114, Index ‚teil-superseded: ADR-0116‘“.

## Bestätigt, mit akzeptiertem Negativ

- **„Höchstens drei Leerzeichen“ gegen den Gegenstand (Prüffrage 3):**
  - `proseLineSet` fragt nur „außerhalb von Fenced-Code“. Eine eingerückte Zeile ohne Zaun liegt für sie in Prosa,
    gleich wie viele Leerzeichen davor stehen.
  - Der Schnitt bei vier ist derselbe, den die Anforderung für den eingerückten Code-Block ohne Zaun zieht
    ([lastenheft.md](../../spec/lastenheft.md) Z. 4626). Eine Zeile mit vier Leerzeichen bleibt also unter jeder
    Lesart gemessen.
  - Gezählt wird U+0020; `TrimFenceIndent` (Space und Tab) betrifft nur die Zaun-Erkennung, nicht diese Frage.
  - Die Lesart ist eine deterministische Näherung ohne Listenkontext. Gegen CommonMark genau ist sie für die Fälle
    unter `- ` und `1. `, die der Anlass zeigt (Issue: drei Leerzeichen unter `1. `). Die übrigen Fälle stehen als
    Grenzen in ADR-0116.
- **Die vom Planner genannten Ränder:**
  - Tab, U+00A0 und `> |` bleiben Befund (Kriterium 2).
  - Die eingerückte `|`-Zeile im Zaun bei `exempt-code-blocks: false` bleibt Befund (Kriterium 3); der Helfer
    verlangt Prosa, weil er `tableRowLine` einschließt.
  - Ein Wort über der Grenze in einer eingerückten Tabellenzeile: (b) nimmt sie vor (c) aus
    (`linelength.go` Z. 72 vor Z. 80), sie zählt nicht in `U`, wie die Spalte-0-Zeile.
  - Fehlende Entscheidung: keine außer F-1.
- **Spiegel-Zählung (Prüffrage 4):** `git grep -c 'exempt-tables' -- . ':!docs/plan/planning' ':!docs/reviews' ':!CHANGELOG.md'`
  ergibt am Stand `c9379511` **10 Dateien, 27 Zeilen**, wie slice-293 §3 sagt. Was sich ändert und was nicht:
  - Lastenheft (4 Zeilen): Z. 4603 (Beschreibung) und 4707 (Kriterium „Grenzen“) ändern sich; 4706 und 4715
    bleiben.
  - Spezifikation (4): Z. 3736 (Schritt 4) und 4116 (Schema) ändern sich; 3702 und 3711 (Eingabefehler) bleiben.
  - Handbuch (3): Z. 1679 und 1706 ändern sich, 1725 bleibt; dazu Z. 1713 und 3135 ohne den Schlüssel.
  - `config_template.go` Z. 314 ändert sich.
  - `linelength_test.go`: neue Fälle.
  - Ohne Änderung bleiben `configyaml.go`, `configyaml/linelength_test.go`,
    `cli_acceptance_test.go` Z. 3457 (schreibt nur eine Konfiguration), `tools/image-test.sh` (Tabelle in
    Spalte 0) und ADR-0114 (unveränderlich).

  Ein Spiegel-Test, der Text davon hält, ist nicht dabei. `docexamples_test.go` prüft die YAML-Blöcke des
  Handbuchs nur auf Gültigkeit, ein geänderter Kommentar bleibt gültig; `make test` hält
  `docs/user/abdeckung-tests.md` gegen die Ableitung (slice-294 §3 Z. 113). Die Dateitabelle von slice-294 deckt
  alles ab.

  **Akzeptiertes Negativ:** `README.md` (Modul-Bullet `linelength`, Z. 212–218) und `README.de.md` Z. 215 („Zeilen, die mit `|`
  beginnen“) werden ungenau. Sie sind README-Prosa und gehören in die Release-Prep
  ([`AGENTS.md`](../../AGENTS.md) §5 Regel 17), wie slice-294 Z. 144 sagt.
- **(a2) und (b) des Closure-Triggers (Prüffrage 5):** (a2) trägt, mit der Vorher-Seite aus F-2. Das Fixture ohne
  eingerückte `|`-Zeilen belegt die Gleichheit, das zweite Fixture die eine gewollte Abweichung samt Gegenprobe
  `exempt-tables: false`. (b) fährt die Form des Issues und die drei Grenzen am Image.
- **MR-032:** slice-293 DoD-Punkt 1 verlangt Version und **eine** Historie-Zeile je Version, die Lehre aus
  welle-96 R2 L-2. Die Änderung ist eine Verhaltensänderung, also ein Minor-Schritt über 0.114.3.
- **Andere Tabellen-Leser außerhalb der Kopplung:** `planning_observations.go` Z. 152–154 trimmt vor dem `|` und
  liest eine eingerückte Registerzeile; `metrics_table.go` Z. 211–213 verlangt Spalte 0. Beides ist Bestand,
  berührt diese Welle nicht und ist keine Folge von Entscheid 4. **Akzeptiertes Negativ** ohne eigenen Träger;
  den Bestand prüft der nächste Trigger-Audit von ADR-0054 („eine vierte Stelle“).

## Nicht geprüft

- Kein Code gefahren; die Aufrufer sind gezählt, nicht getestet.
- Ob der Bestand des Melders Tabellen mit vier und mehr Leerzeichen trägt — außerhalb des Repos.
- Gates: nur `make doc-check` über diesen Report, ADR-0116 und den Index.

## Zusammenfassung

| Befund | Verdikt | Träger der Korrektur |
|---|---|---|
| F-1 Helfer in der Lexik, schließt `tableRowLine` ein; „zusätzlich“ im Plan | V2 + V1 | ADR-0116 Entscheid 4; Plan 294 §3; Welle §5 |
| F-2 (a1)/(a2) gegen `v0.92.0` scheitern am Profil | V1 | Welle §3 |
| F-3 Lastenheft §5 nicht geführt | V1 | Plan 293 Kopf, §3 |
| F-4 drei Grenzen ohne Text | V1 | Plan 293 DoD 1; Plan 294 §3 |
| F-5 fälliger Trigger von ADR-0114 ist (b) | V1 | Welle §3 (c) |
| F-6 Nachzug an ADR-0114 | V1 | Plan 294 DoD 3, §3 |

Kein V3: keine Schwelle, kein Gate wird gelockert. `linelength` ist opt-in und im eigenen Profil nicht an.

## Plan-Änderungen für den Planner

1. slice-294 §3 Z. 107/108: der Helfer ruft `tableRowLine` und erweitert sie; `linelength` ruft nur den Helfer;
   die Mutation zu Kriterium 6 wird „der Helfer ruft `tableRowLine` nicht mehr“; Z. 110 nennt den Kopf-Kommentar
   des Kopplungs-Tests (F-1).
2. Welle §5 Frage 4: das Verdikt und ADR-0116 statt der Empfehlung (F-1).
3. Welle §3 (a1) `REF=69c53e0f`, (a2) dasselbe Vorher-Image, mit Begründung (F-2).
4. slice-293 Kopf und §3: Lastenheft §5 geprüft und unverändert, mit Grund (F-3).
5. slice-293 DoD 1: drei Grenzen im Punkt „Tabellenzeile heißt …“; slice-294 §3 Testfälle zu Grenze 1 und 2
   (F-4).
6. Welle §3 (c): Trigger (b) von ADR-0114, nicht (a) (F-5).
7. slice-294 DoD 3 und §3 Z. 114: Geschichte-Zeile in ADR-0114 und Index „teil-superseded: ADR-0116“ im selben
   Commit wie ADR-0116 `Accepted` (F-6).
