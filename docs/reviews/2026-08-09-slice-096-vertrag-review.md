# Review-Report: slice-096 (Vertrags- und Konsistenz-Seite) — 2026-08-09

**Review-Art:** Design — geprüft wird die **Vertragsseite** des Change Requests
gegen die Prozess- und Spec-Invarianten: Lastenheft-Zusage ↔
Spezifikations-Algorithmus ↔ ADR-Begründung ↔ Slice-Messung. Die
Implementierbarkeit prüft ein zweiter Reviewer parallel; Code wird hier nur
als **Beleg** herangezogen, wo eine Vertragsaussage sonst nicht falsifizierbar
wäre. Die abzulösenden Adopter-Skripte des Schwester-Repos `a-check` wurden
**nur gelesen**.

**Gegenstand:** [slice-096](../plan/planning/in-progress/slice-096-structure-modul-analyse.md),
Commit-Range `45246bb..acbb419` (drei Commits: Messung + Abnahme-Punkt 1 ·
Abnahme-Punkte 2–5 · CR in Lastenheft/Spezifikation + ADR + Folge-Slices)

**Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) @ 1.3.0 (`419e82e`) ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-08-09

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne diese Liste
ist der Lauf nicht reproduzierbar):

- [slice-096](../plan/planning/in-progress/slice-096-structure-modul-analyse.md)
  (Slice-Plan, Messung §2 + Abnahme-Punkte 1–5 + DoD),
  [slice-099](../plan/planning/open/slice-099-structure-modul-kern.md),
  [slice-100](../plan/planning/open/slice-100-structure-marken-und-zaehlung.md),
  [slice-094](../plan/planning/open/slice-094-closure-zaehl-paritaet.md),
  [slice-097](../plan/planning/open/slice-097-closure-glob-entkopplung.md),
  [slice-098](../plan/planning/open/slice-098-closure-note-placeholder.md)
- [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) (neu, Proposed) samt
  [ADR-Index](../plan/adr/README.md); Bezugs-ADRs
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md),
  [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md),
  [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md)
- [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) (neu),
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (Mit-Modifikation),
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl),
  [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in),
  [`DC-FA-HOST-001`](../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in),
  [`DC-FA-TRK-001`](../../spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in),
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus),
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
  [`DC-QA-04`](../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools);
  §1, §3-Bereichsliste, §6-Glossar, §7-Historie desselben Dokuments
- [`spec/spezifikation.md`](../../spec/spezifikation.md) §[`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure),
  §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) (C1–C5),
  §[2-Schema](../../spec/spezifikation.md#2-datenstrukturen-und-schemas),
  §[4-Grund-Codes](../../spec/spezifikation.md#4-grund--und-fehler-codes), §7-Historie
- [`AGENTS.md`](../../AGENTS.md) §2 (Source Precedence), §3 (Hard Rules, insb. §3.4/§3.5), §5 (Dokumentations-Regeln)
- Referenz-Richtungs-Matrix des vendorten Regelwerks
  ([`grundlagen-referenz-richtung.md`](../../.harness/baseline/v5.0.0/regelwerk/grundlagen-referenz-richtung.md))
- **Nur lesend beigezogen:** die drei vermessenen Prüfskripte des Schwester-Repos
  `a-check` (verify-closure-notes, verify-slice-form, verify-ac-form) samt dessen
  `done/`-Bestand und Lastenheft

---

## Findings

### F-1 — Die Messzeile, die den Modul-Schnitt entscheidet, ist unter dem beschlossenen Vertrag nicht ausdrückbar; der sanktionierte Ausweg ist still grün

- `kategorie`: HIGH
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) /
  [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) Entscheidung 1
- `pfad`: `spec/lastenheft.md:1965`, `spec/lastenheft.md:2022`, `spec/lastenheft.md:2041`
- `befund`: Der Vertrag bindet jede Regel an **genau einen** Abschnitt, meldet
  bei mehr als einem Treffer `section-ambiguous` und bricht dann für die Datei ab
  (Akzeptanzkriterium „Negative (mehrdeutig)"), und erklärt „mehr als ein
  Abschnitt pro Regel" zum Out-of-Scope. Messzeile 10 — laut §2 des Slice-Plans
  „die Zeile, die den Schnitt entscheidet" — verlangt die vier Pflicht-Marken
  **je Anforderung**, und Anforderungen sind **wiederkehrende Abschnitte
  innerhalb einer Datei**: das vermessene Skript des Adopters iteriert über
  19 `### AC-`-Überschriften **eines** Lastenhefts. Ein `section-pattern` über
  diese Klasse liefert `section-ambiguous` statt der Marken-Prüfung; der laut
  Out-of-Scope vorgesehene Ausweg („wer zwei Abschnitte prüfen will, schreibt
  zwei Regeln") lässt jede **neu** hinzukommende Anforderung ungeprüft — genau
  die Eigenschaft, die das abzulösende Skript in seinem Kopfkommentar als seinen
  Kern benennt („die Liste wächst NICHT mit: jede künftige AC-ID fällt
  automatisch unter die Regel").
- `verifizierbar`: nein am Gate (der Code existiert noch nicht) — ja am Vertrag:
  die beiden Zusagen stehen im selben Dokument und schließen einander aus; nach
  [slice-100](../plan/planning/open/slice-100-structure-marken-und-zaehlung.md)
  wäre der Nachweis eine `structure`-Regel mit Anforderungs-`section-pattern`
  gegen [`spec/lastenheft.md`](../../spec/lastenheft.md), die `section-ambiguous`
  statt `section-constraint` liefert.
- `klasse`: „Wiederkehrender Abschnitt als Dokumentklasse nicht modellierbar"

### F-2 — `require-any` deckt Messzeile 8 nicht: 36 von 36 realen Lerneintrag-Marken stehen mitten in der Zeile

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
  [`DC-QA-04`](../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
- `pfad`: `spec/lastenheft.md:1986`, `spec/spezifikation.md:1761`
- `befund`: Der Vertrag definiert benannte Marken als **zeilenverankert und
  hervorgehoben** — vorhanden ist eine Marke `M` nur, wenn eine Zeile nach
  führendem Whitespace mit `**M**` oder `**M:**` **beginnt**. Messzeile 8, die
  `require-any` überhaupt erst begründet, prüft im Adopter-Skript dagegen ein
  Vorkommen von `Form: <eine von drei>` **irgendwo** im Abschnitt; im
  `done/`-Bestand des Adopters tragen alle 36 Fundstellen die Form
  `**Lerneintrag — Form: geschärfte Regel.**` — die Marke steht **mitten** in
  der fetten Zeile, keine einzige Zeile beginnt mit `**Form:`. Unter der
  zugesagten Semantik meldet die Regel für jede konforme Notiz
  `section-constraint` (Falsch-Rot); benennt der Adopter stattdessen die Marke
  `Lerneintrag`, prüft die Regel die Alternation über die drei Formen gar nicht
  mehr (Falsch-Grün).
- `verifizierbar`: ja — Paritäts-Beleg aus dem DoD von
  [slice-100](../plan/planning/open/slice-100-structure-marken-und-zaehlung.md)
  gegen die beigezogenen Adopter-Fixtures; er kann unter dieser Marken-Semantik
  nicht grün werden.
- `klasse`: „Marken-Verankerung passt nicht zur belegten Alt-Form"

### F-3 — Die Summenzeile der Messung geht nicht auf, und die falsche Zahl steht im Lastenheft und in der ADR

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) §7-Historie /
  [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) §Kontext
- `pfad`: `docs/plan/planning/in-progress/slice-096-structure-modul-analyse.md:76`,
  `spec/lastenheft.md:2410`, `docs/plan/adr/0049-structure-modul-schnitt-und-preset.md:37`
- `befund`: Die Tabelle führt **elf** Zeilen; ausgezählt ergibt sie 2 gedeckt,
  2 nach Kalibrierung, **5** nicht gedeckt (Zeilen 2, 4, 7, 8, 10), 1 außerhalb,
  1 Ventil. Die Summenzeile behauptet „6 nicht gedeckt" und addiert sich damit
  auf 12 statt 11. Dieselbe „6" ist in die Historie-Zeile 0.51.0 des Lastenhefts
  („11 Prüfungen, davon 2 gedeckt, 2 nach Kalibrierung, 6 ungedeckt") und in den
  Kontext der ADR („elf Prüfungen — 2 heute gedeckt, 2 nach Kalibrierung, 6
  ungedeckt, 1 außerhalb, 1 als Ventil") übernommen worden; die ADR wird mit
  `Accepted` immutabel ([`AGENTS.md`](../../AGENTS.md) §3.5).
- `verifizierbar`: ja — Auszählen der Spalte „Aussage" der Tabelle in §2 des
  Slice-Plans; kein Gate deckt Aggregate über eigene Tabellen ab.
- `klasse`: „Aggregat widerspricht der Einzelaufstellung"

### F-4 — Messzeile 4 ist als „nicht gedeckt" mit dem Befund eines anderen Antrags belegt; das vermessene Skript prüft literale Phrasen, die der bestehende Vertrag abdeckt

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (Floskel-Bedingung)
- `pfad`: `docs/plan/planning/in-progress/slice-096-structure-modul-analyse.md:67`
- `befund`: Zeile 4 („kein Template-Platzhalter") ist als **nicht gedeckt**
  eingestuft, Beleg: „vier Platzhalter-Sätze passieren alle drei Codes". Diese
  vier Sätze stammen aus der Messung von
  [slice-098](../plan/planning/open/slice-098-closure-note-placeholder.md), also
  aus dem Antrag eines **anderen** Konsumenten und über eine
  Winkelklammer-Erkennung. Die Platzhalter-Prüfung des hier vermessenen
  a-check-Skripts ist dagegen eine Liste von fünf **literalen** Phrasen
  (`TODO`, `TBD`, `noch offen`, `beim Abschluss`, `_(folgt)_`) — genau die
  Eingabeform, die `planning.closure.boilerplate` seit v0.52.0 als
  case-insensitive Literal-Teilstrings entgegennimmt. Zeile 4 ist damit für das
  vermessene Skript höchstens „nach Kalibrierung" (Unterschied: der Vertrag
  bereinigt vor der Prüfung Fenced-Code, das Skript nicht), nicht „nicht
  gedeckt"; der Beleg trägt die Aussage der Zeile nicht.
- `verifizierbar`: ja — `d-check --enable planning` mit
  `planning.closure.boilerplate` = den fünf Phrasen gegen den `done/`-Bestand:
  die Klasse fällt heute schon.
- `klasse`: „Beleg stammt aus anderer Quelle als die vermessene Prüfung"

### F-5 — `min-sentences`: der spezifizierte Default unterschreitet die im selben Vertrag zugesagte Exit-2-Grenze

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
- `pfad`: `spec/lastenheft.md:2011`, `spec/lastenheft.md:2029`, `spec/spezifikation.md:2195`
- `befund`: Das Lastenheft nennt als fail-closed-Grund „ein `min-sentences` < 1"
  und das Akzeptanzkriterium „fail-closed (Config-Rand)" fordert Exit 2 für „eine
  Regel … mit `min-sentences` < 1" — ohne die Einschränkung auf einen **explizit
  gesetzten** Wert. Die Spezifikation setzt den Default auf `0` („aus"). Wörtlich
  gelesen bricht damit jede Regel, die `min-sentences` nicht setzt, mit Exit 2 —
  der Normalfall. Die parallele Fähigkeit in
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  führt genau diese Unterscheidung ausdrücklich („ein **abwesendes**
  `min-sentences` ist dagegen der Default (kein Fehler)"); der neue Vertrag
  übernimmt sie nicht.
- `verifizierbar`: ja — ein Test nach dem Akzeptanzkriterium und ein Test nach
  dem §2-Schema widersprechen sich für dieselbe Eingabe (Regel ohne
  `min-sentences`).
- `klasse`: „Default unterschreitet die eigene Validierungsgrenze"

### F-6 — Nullmengen-Guard und `exempt-paths`: die abnahmebindende Lesart lässt eine vollständig ausgehebelte Regel grün laufen

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
  [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md) (Nullmengen-Logik)
- `pfad`: `spec/lastenheft.md:2012`, `spec/lastenheft.md:2027`, `spec/spezifikation.md:1732`
- `befund`: Das Lastenheft knüpft den Leerlauf-Befund an den **Glob**: „Ein Glob
  **ohne** passende Datei ⇒ `section-missing`", und das Ventil-Kriterium sagt für
  eine `exempt-paths`-Datei „kein Befund für dieses Dokument". Die Spezifikation
  bildet die Kandidatenmenge dagegen als „`files`-Treffer **abzüglich
  `exempt-paths`**" und meldet „Null Kandidaten ⇒ `section-missing`". Für eine
  Regel, deren Glob Dateien trifft, die sämtlich per `exempt-paths` ausgenommen
  sind, sagt das Vertrags-Stratum grün, das Technik-Stratum rot. Nach Source
  Precedence ([`AGENTS.md`](../../AGENTS.md) §2) gewinnt die grüne Lesart — eine
  Regel, die nichts mehr prüft, meldet Erfolg.
- `verifizierbar`: ja — eine Regel mit `files` und deckungsgleichem
  `exempt-paths` gegen beide Formulierungen gefahren.
- `klasse`: „Ventil hebelt den Nullmengen-Guard aus"

### F-7 — Die Preset-Kopplung ist die tragende Zusage des CR, wird aber von keinem Akzeptanzkriterium getragen

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) /
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `pfad`: `spec/lastenheft.md:1900`, `spec/lastenheft.md:1995`, `spec/lastenheft.md:2019`
- `befund`: Beide Anforderungen sagen zu, dass die Closure-Fähigkeit ein
  **Preset** derselben Struktur-Semantik ist („gleiche Abschnitts-Bestimmung,
  gleiche Fence-Behandlung, gleiche Zählung … verändert werden darf nur beides
  zugleich"), und der Slice-Plan leitet daraus ab, die beiden könnten „nicht
  auseinanderlaufen, ohne dass ein Test es merkt". Keines der elf neuen
  Akzeptanzkriterien von
  [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  und keines der ergänzten von
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  prüft die Kopplung. Sie existiert als Fitness Function in der ADR und als
  DoD-Punkt in
  [slice-099](../plan/planning/open/slice-099-structure-modul-kern.md) — beides
  Schichten, die keine Abnahme binden. Driften die beiden Oberflächen, fällt kein
  Kriterium.
- `verifizierbar`: ja — die Kriterienliste beider Anforderungen enthält keinen
  Given/When/Then-Satz über zwei Oberflächen.
- `klasse`: „Zugesagte Invariante ohne Akzeptanzkriterium"

### F-8 — „erstes Modul ohne Referenz-Invariante" ist durch den eigenen Bestand widerlegt

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in),
  [`DC-FA-HOST-001`](../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)
- `pfad`: `spec/lastenheft.md:1961`, `spec/lastenheft.md:17`,
  `docs/plan/adr/0049-structure-modul-schnitt-und-preset.md:26`
- `befund`: Der Vertrag stellt `structure` als „das erste Modul, das **keine
  Referenz**-Invariante prüft" vor, und der Kontext der ADR begründet den ganzen
  Antrag damit, der Modulsatz enthalte „**keine einzige Aussage über die Form
  eines Dokuments selbst**". `spans` meldet mit `span-unclosed` eine
  ungeschlossene Backtick-Folge innerhalb eines Absatzes und mit
  `span-nested-link` verschachtelte Link-Syntax — beides Aussagen über die Form
  des Dokuments, ohne dass ein Ziel aufgelöst wird; `hostpaths` verbietet eine
  Pfad-**Gestalt** in Prosa und Inline-Code, ebenfalls ohne Existenz-Prüfung. Die
  Identitäts-Aussage in §1 des Lastenhefts und die Motivations-Aussage der ADR
  sind damit sachlich zu weit gefasst; die ADR wird mit `Accepted` immutabel.
- `verifizierbar`: ja — §4-Grund-Code-Tabelle der Spezifikation, Zeilen
  `span-unclosed`, `span-nested-link`, `hostpath-forbidden`.
- `klasse`: „Identitäts-Aussage ohne Bestandsprüfung"

### F-9 — `section-constraint` ist ein Sammel-Code über sechs Bedingungen; die Unterscheidung liegt in dem Feld, das die Spezifikation ausdrücklich nicht stabil zusagt

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) /
  [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) Entscheidung 3 + 4
- `pfad`: `spec/lastenheft.md:1983`, `spec/lastenheft.md:2023`,
  `docs/plan/adr/0049-structure-modul-schnitt-und-preset.md:122`
- `befund`: Der Vertrag fasst Verletzungen von `non-empty`, `min-sentences`,
  `max-tasks`, `forbid-pattern`, `require-all` und `require-any` in **einem**
  Grund-Code zusammen und verlagert die Unterscheidung in die `message`, die die
  Spezifikation im Befund-Schema als „nicht stabilitätsgarantiert" führt.
  Dieselbe ADR erhebt die Stabilität der Grund-Codes zum Ausschlussgrund gegen
  ein Supersede („jede Konsumenten-CI, die auf den Code filtert, bricht") und
  verwirft in der Alternativen-Tabelle Sammelbefunde als Bauform („Sammelbefunde
  waren schon bei den drei Closure-Codes die verworfene Bauform"). Für einen
  Konsumenten, der wie beschrieben auf `reason` filtert, ist „Abschnitt zu dünn"
  von „verbotenes Muster getroffen" nicht unterscheidbar. Die Konsequenzen-Liste
  der ADR benennt das nicht; registriert ist es nur als offener Punkt in
  [slice-099](../plan/planning/open/slice-099-structure-modul-kern.md) §4 — einer
  Schicht, die keine Entscheidung bindet.
- `verifizierbar`: ja — Befund-Schema §2 der Spezifikation (`message` … nicht
  stabilitätsgarantiert) gegen das Akzeptanzkriterium „Negative (Bedingung
  verletzt)", das genau diesen Feldinhalt zusagt.
- `klasse`: „Sammel-Code trotz verworfener Sammelbefund-Bauform"

### F-10 — Das Vertrags-Stratum begründet eine Zusage mit einem Abwärts-Zeiger auf die Spezifikation, gate-unsichtbar

- `kategorie`: MEDIUM
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.4 / Referenz-Richtung (SDP)
- `pfad`: `spec/lastenheft.md:1999`
- `befund`: Die neue Anforderung begründet im Fließtext: „Sie bleibt bestehen und
  behält ihre Grund-Codes — die sind stabil zugesagt (**§4 der Spezifikation**)".
  Die Referenz-Richtungs-Matrix des vendorten Regelwerks führt für die Zeile
  **Vertrag** (Decke) gegenüber der Spalte **Technik** ein ❌ („Spec-Straten —
  hinein ja, hinaus nie"); die repo-lokale Kodierung in
  [`.d-check.yml`](../../.d-check.yml) verbietet denselben Verweis als
  `matrix-downward`. Weil der Zeiger als **Prosa** ohne Markdown-Link, ohne
  Inline-Code-Pfad und ohne Token auftritt, sehen ihn weder `matrix` noch `links`
  noch `codepaths` — die Verletzung ist gate-unsichtbar, und sie ist
  **normativ** (sie trägt die Zusage, nicht bloß Provenance).
- `verifizierbar`: nein maschinell (kein Gate greift auf Prosa-Zeiger) — ja gegen
  die Matrix-Zeile „Vertrag → Technik ❌".
- `klasse`: „Abwärts-Referenz als Prosa-Zeiger"

### F-11 — Nach der Mehrdeutigkeits-Härtung steht „der erste Treffer" noch an drei Stellen

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `pfad`: `spec/lastenheft.md:1876`, `spec/spezifikation.md:1678`, `spec/spezifikation.md:2188`
- `befund`: Der CR ergänzt `closure-note-ambiguous` mit Abbruch bei mehr als
  einem Treffer. Unverändert stehen geblieben sind: die Beschreibung der
  Anforderung („in ihr **der erste** Abschnitt, dessen Überschrift auf
  `planning.closure.heading-pattern` passt"), der Eingangssatz von Schritt C3
  („wird die **erste** Zeile gesucht") und die §2-Schema-Zeile zu
  `planning.closure.heading-pattern` („der **erste** Treffer eröffnet den
  geprüften Abschnitt"). Die Schema-Zeile ist damit die einzige Stelle, an der
  ein Implementer den alten, jetzt widersprochenen Vertrag ablesen kann; die drei
  Stellen und der neue Absatz beschreiben dieselbe Konfiguration
  gegensätzlich.
- `verifizierbar`: ja — Textvergleich innerhalb derselben zwei Dateien.
- `klasse`: „Teil-Nachzug einer Semantik-Änderung"

### F-12 — Messzeile 11 ordnet eine namentliche ID-Ausnahme innerhalb einer Datei dem Pfad-Ventil zu; ausdrückbar ist sie dort nur als Abschaltung der ganzen Regel

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) §Out-of-Scope /
  [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) Entscheidung 6
- `pfad`: `docs/plan/planning/in-progress/slice-096-structure-modul-analyse.md:74`,
  `docs/plan/planning/in-progress/slice-096-structure-modul-analyse.md:189`
- `befund`: Zeile 11 führt die Grandfathering-Mechanik beider Dokumentklassen als
  eine „Stichtags-Ausnahme" und stuft sie als „Ventil — über die bestehende
  Ausnahme-Mechanik ausdrückbar" ein; Abnahme-Punkt 5 behandelt nur den
  Dateinamen-Stichtag der Slice-Klasse. Für die Anforderungs-Klasse ist die
  Ausnahme des Adopters aber **keine** Stichtags-Regel, sondern eine namentlich
  aufgezählte Liste von 19 Kennungen **innerhalb einer einzigen Datei** — und
  diese 19 sind aktuell der **gesamte** Bestand (19 von 19 Überschriften). Ein
  Pfad-Ventil kann eine Teilmenge von Abschnitten einer Datei nicht adressieren;
  der einzige mögliche `exempt-paths`-Ausdruck ist die Datei selbst, womit die
  Regel vollständig abgeschaltet ist und grün läuft. Die Zeile ist damit für die
  Anforderungs-Klasse falsch als „Ventil" klassifiziert.
- `verifizierbar`: nein am Gate — ja gegen die Ausnahme-Liste des vermessenen
  Skripts und die Zahl der Anforderungs-Überschriften im Lastenheft des Adopters.
- `klasse`: „Ventil-Granularität passt nicht zur belegten Ausnahme"

### F-13 — Die Zeitangabe der ADR zum Contract-Churn stimmt nicht: die abgelöste Fähigkeit ist am selben Tag erschienen

- `kategorie`: LOW
- `quelle`: [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) §Kontext
- `pfad`: `docs/plan/adr/0049-structure-modul-schnitt-und-preset.md:40`,
  `docs/plan/planning/in-progress/slice-096-structure-modul-analyse.md:247`
- `befund`: Der Kontext schreibt, d-check habe „**wenige Tage zuvor**" mit
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md) die
  Closure-Note-Struktur ausgeliefert (v0.52.0); §5 des Slice-Plans spricht von
  einem „Bruch nach wenigen Tagen". Tag, Changelog-Eintrag und Release-Register
  datieren v0.52.0 auf **denselben** Tag wie diesen CR. Die Angabe ist die
  einzige Tatsachenbasis, mit der die ADR die Churn-Größe einordnet, und sie
  wandert mit `Accepted` in ein immutables Dokument.
- `verifizierbar`: ja — [`version.md`](../../version.md) §Aktuell und der
  `0.52.0`-Eintrag in `CHANGELOG.md`.
- `klasse`: „Zeitangabe ohne Quellenabgleich"

### F-14 — Die zweite Modulliste desselben Dokuments bleibt bei 19 Modulen

- `kategorie`: LOW
- `quelle`: [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
- `pfad`: `spec/lastenheft.md:2396`
- `befund`: Der CR ergänzt `structure` in der Modulliste von
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
  (jetzt 20 Namen), lässt die zweite Enumeration derselben Menge im §6-Glossar
  („Regelmodul") aber bei 19 Namen mit `sources` als letztem stehen. Der
  Vorgänger-CR hatte beide Stellen gepflegt und das in seiner Historie-Zeile
  ausdrücklich vermerkt („`sources` in `DC-FA-CLI-002` + Glossar"); die neue
  Historie-Zeile nennt nur noch `DC-FA-CLI-002`. Kein Gate deckt Prosa-Modullisten
  ab — die Netzlos-Modullisten-Integrität in `make test` prüft
  [`.d-check.yml`](../../.d-check.yml), nicht das Lastenheft.
- `verifizierbar`: nein — kein Gate; Textvergleich der beiden Listen.
- `klasse`: „Doppelte Enumeration nur einseitig nachgezogen"

### F-15 — Der Leerlauf-Befund hat keinen definierten `line`-Wert

- `kategorie`: LOW
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
- `pfad`: `spec/lastenheft.md:2028`, `spec/spezifikation.md:1767`
- `befund`: Für den Nullmengen-Fall sagt der Vertrag „`section-missing` auf dem
  Glob" zu; die Befund-Form der Spezifikation regelt `file` und `target` für
  diesen Fall ausdrücklich, für `line` aber nur „wie oben" — und „oben" ist der
  Abschnitts-Fall (`line` = 1 bzw. Zeile der Überschrift). Das
  Befund-Schema verlangt `line` als Integer ≥ 1. Die Parallel-Stelle in
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  pinnt den Wert („bzw. 1"). Ein Test auf den Leerlauf-Befund kann die Zeile
  nicht assertieren, ohne den Vertrag zu erfinden.
- `verifizierbar`: ja — Akzeptanzkriterium „fail-closed (Regel leer)" nennt keinen
  beobachtbaren Zeilenwert.
- `klasse`: „Befund-Feld ohne definierten Wert im Randfall"

### F-16 — Die Aussagen-Zählung übergeht die Nullmengen-Guards, obwohl der CR genau diese Aussage neu in den Vertrag nimmt

- `kategorie`: LOW
- `quelle`: [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
- `pfad`: `docs/plan/planning/in-progress/slice-096-structure-modul-analyse.md:58`
- `befund`: Die Messung zählt bewusst „je Prüfung eine Aussage, nicht je Skript".
  Jedes der drei vermessenen Skripte trägt zusätzlich einen eigenen
  Grundgesamtheits-Guard („null Dateien sind Bestandsverlust, nicht nichts zu
  prüfen") — drei weitere Aussagen derselben Art, die keine Tabellenzeile
  bekommen, obwohl der neue Vertrag sie als `section-missing` auf dem Glob
  ausdrücklich aufnimmt und
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  sie bereits abdeckt. Die Auslassung verschiebt die Bilanz zulasten der Spalte
  „gedeckt" und damit zugunsten der Antrags-Begründung.
- `verifizierbar`: nein — Vergleich der Tabelle mit den drei gelesenen Skripten.
- `klasse`: „Aussagen-Zählung übergeht die Guard-Klasse"

### F-17 — Vorbestehende Zahlen-Drift im Dogfooding-Kommentar, die dieser CR vergrößert

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `.d-check.yml`
- `befund`: Der Kommentar am `trace`-Block begründet den scharfen
  Nullmengen-Guard mit „Heute nachweislich safe (42 Anforderungen, 0 Waisen)".
  Das Lastenheft führt nach diesem CR 48 Anforderungs-Überschriften (47 davor) —
  die Drift bestand also bereits vor dem Diff und wächst mit jeder neuen
  Anforderung. Kein Gate prüft die Zahl; sie ist reine Begründungs-Prosa. Nicht
  Gegenstand dieses Diffs, hier nur als Beobachtung festgehalten.
- `verifizierbar`: ja — `make trace` nennt die Ist-Zahl.
- `klasse`: „Kommentar-Kennzahl ohne Nachzug-Bindung"

## Negativbefunde

- geprüft, ohne Befund: **die tragende Behauptung der ADR** („die Spezifikation
  führt die Grund-Codes als *stabil, maschinenlesbar*, im Unterschied zur
  `message`"). Beide Zitate stehen **wörtlich** so in der Spezifikation (§4
  Einleitungssatz; Befund-Schema §2, Feld `message`). Die Folgerung trägt: ein
  Umbenennen ausgelieferter Codes bräche die einzige als stabil zugesagte Fläche,
  additive Codes tun das nicht — dieselbe Lesart, unter der bereits
  `source-drift`/`source-unreachable` nachgezogen wurden. Einschränkung nur bei
  der Zeitangabe (F-13) und beim Sammel-Code (F-9).
- geprüft, ohne Befund: **Referenz-Richtung der Spezifikation.** Die neuen
  Passagen und die neue §7-Historie-Zeile enthalten **keinen** Token und
  **keinen** Link auf ADRs, Wellen, Slices oder Commit-Hashes (maschinell
  gegengeprüft über den Diff); die ADR-Erwähnung bleibt generisch („in
  begleitender ADR"), wie in den Vorgänger-Zeilen. Auch die neue
  Lastenheft-§7-Zeile ist token-frei.
- geprüft, ohne Befund: **ADR-Form.** `Schärft:` zeigt ausschließlich auf die
  Spezifikation (die neue Algorithmus-Sektion und die Preset-Kopplung), nicht auf
  das Lastenheft; das Lastenheft steht korrekt unter `Bezug:`. Abschnitts-Gerüst
  (Kontext · Entscheidung · Verglichene Alternativen · Konsequenzen · Fitness
  Function · Re-Evaluierungs-Trigger · Geschichte) ist identisch zur
  Haus-Form von
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md).
  `Re-Evaluierungs-Trigger` vorhanden und mit drei prüfbaren Bedingungen belegt.
  Der einzige Slice-Verweis der ADR steht im Abschnitt `Geschichte`, der in
  [`.d-check.yml`](../../.d-check.yml) als `exclude-sections` geführt wird.
- geprüft, ohne Befund: **ADR-Index-Zeile.** Neue Zeile vorhanden, Dateilink
  korrekt, Status `Proposed`, Datum, Bezugs-IDs als Links, Zusammenfassung deckt
  alle sechs Entscheidungen ab.
- geprüft, ohne Befund: **Alternativen-Tabelle.** Sieben Zeilen, darunter beide
  naheliegenden Bauformen (Ausbau von `planning.closure`; Supersede mit Alias)
  und die nicht offensichtliche dritte (ein Modul mit zwei Code-Familien je nach
  Config-Herkunft). Jede mit einem eigenen, nicht wiederholten Verwerfungsgrund.
- geprüft, ohne Befund: **Zahlen der Messung außer F-3.** „480 Zeilen Shell in
  drei Prüfskripten" trifft exakt (170 + 166 + 144); „vier handgeschriebene
  Prüfskripte" trifft (das vierte ist die Referenz-Prüfung); „76 Slices" war ein
  realer Bestandsstand des Adopters am Messtag (heute 79 — die Formulierung
  „gegen eine Kopie" hält das ehrlich); „20. Regelmodul" trifft
  (Auszählung der Liste in
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)).
- geprüft, ohne Befund: **Zuordnung der Messzeilen 1, 2, 3, 5, 6, 7, 9** gegen
  die gelesenen Skripte. Zeile 2 („d-check liest den ersten Treffer") entsprach
  dem Stand vor diesem CR; Zeile 5 (Floskel zeilenverankert vs. Teilstring),
  Zeile 6 (Inline-Code + Satzende-Form) und Zeile 7 (Zählung dateiweit vs.
  abschnittstreu) sind an den Skripten belegbar; Zeile 9 (Dateiname) ist korrekt
  als „außerhalb" eingestuft und deckungsgleich mit der Modul-Grenze.
- geprüft, ohne Befund: **Entscheid über die drei liegenden Slices.** Keiner ist
  durch den neuen Vertrag überflüssig geworden.
  [slice-094](../plan/planning/open/slice-094-closure-zaehl-paritaet.md) schärft
  eine Zählung, die der neue Vertrag als gemeinsame Mechanik festschreibt — die
  Kopplung macht ihn breiter wirksam, nicht entbehrlich.
  [slice-097](../plan/planning/open/slice-097-closure-glob-entkopplung.md)
  betrifft die Kandidatenmenge der Closure-Fähigkeit; `structure` bringt einen
  eigenen Glob mit und löst dort nichts mit.
  [slice-098](../plan/planning/open/slice-098-closure-note-placeholder.md) trägt
  drei Falsch-Positiv-Ausschlüsse als Code, die `forbid-pattern`
  konstruktionsgemäß nicht mitbringt; die DoD sagt das ausdrücklich.
- geprüft, ohne Befund: **Aktiv-Status-Asymmetrie-Behauptung.** Dass die
  Aktiv-Status-Prüfung bei mehrfacher kanonischer Überschrift längst
  fail-closed ist, steht so im Heading-Guard der Anforderung und in der
  §4-Zeile zu `planning-drift` — die Begründung für `closure-note-ambiguous`
  trägt.
- geprüft, ohne Befund: **Inert-Semantik und `DC-QA-02`/`DC-QA-03`-Bindung.** Die
  Zusagen „ohne aktives Modul byte-identisch", „hermetisch, kein git, kein Netz",
  „diagnose-only" sind formuliert und je mit einem Kriterium belegt; die
  Messmethode von
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  ist als „alle Module außer `external`, `sources`, `vcs`" formuliert und
  braucht durch ein weiteres hermetisches Modul keine Änderung.
- geprüft, ohne Befund: **Versions-Bump und Historie-Form.** 0.50.0 → 0.51.0 für
  eine neue Anforderung entspricht der Praxis der Vorgänger-Zeilen; die
  Historie-Zeile nennt Anlass, Grund-Codes, Out-of-Scope und die
  Arbeitsteilung zur Spezifikation/ADR.
- geprüft, ohne Befund: **Commit-Zuschnitt und CHANGELOG.** Der CR-Commit fasst
  nur Spec, ADR, Index und Planung an — kein `CHANGELOG.md`, keine
  README-/Handbuch-Modullisten. Das entspricht dem Zuschnitt des vorangegangenen
  CR-Commits, dessen nutzersichtbare Änderung erst mit der Implementierung
  verbucht wurde.
- geprüft, ohne Befund: **`DC-FA-CONF-001`.** Beschreibt Modul-Parameter
  generisch und braucht für einen neuen Modul-Block keine Ergänzung.
- geprüft, ohne Befund: **Anker-Form der neuen Verweise.** Die
  `#dc-fa-struct-001a`-/`#dc-fa-struct-001--…`-Slugs folgen dem
  GitHub-Slug-Verfahren des Repos; die Kennungen tragen durchweg Links
  (Linkpflicht).
- geprüft, ohne Befund: **Nicht-Ziel-Grenze und `BEO-001`-Abgrenzung.** Die
  Out-of-Scope-Formulierung („Aussagen über den Ort eines Dokuments") deckt die
  in der Messung identifizierte Zeile 9 und ist im Slice-Plan wie in der ADR
  identisch begründet; die ausdrückliche Abgrenzung gegen die offene Beobachtung
  im Register ist in allen drei Slice-Plänen konsistent gehalten.
- geprüft, ohne Befund: **Host-Pfad-Hygiene der neuen Dokumente.** Weder die ADR
  noch die drei Slice-Pläne noch die Spec-Passagen nennen einen Pfad des
  Schwester-Repos; der Adopter wird durchweg ohne Verzeichnis benannt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 11 |
| LOW | 4 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Wiederkehrender Abschnitt als Dokumentklasse
nicht modellierbar · Marken-Verankerung passt nicht zur belegten Alt-Form ·
Aggregat widerspricht der Einzelaufstellung · Beleg stammt aus anderer Quelle
als die vermessene Prüfung · Default unterschreitet die eigene
Validierungsgrenze · Ventil hebelt den Nullmengen-Guard aus · Zugesagte
Invariante ohne Akzeptanzkriterium · Identitäts-Aussage ohne Bestandsprüfung ·
Sammel-Code trotz verworfener Sammelbefund-Bauform · Abwärts-Referenz als
Prosa-Zeiger · Teil-Nachzug einer Semantik-Änderung · Ventil-Granularität passt
nicht zur belegten Ausnahme · Zeitangabe ohne Quellenabgleich · Doppelte
Enumeration nur einseitig nachgezogen · Befund-Feld ohne definierten Wert im
Randfall · Aussagen-Zählung übergeht die Guard-Klasse · Kommentar-Kennzahl ohne
Nachzug-Bindung

## Verdikt

**Merge-blockierend: ja.** Ein HIGH und elf MEDIUM.

Der Schnitt selbst ist tragfähig und sauber begründet: das Kriterium ist
[ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md) entnommen
statt neu erfunden, die Nicht-Supersede-Begründung hält der wörtlichen Prüfung
stand, und die Entscheidung über die drei liegenden Slices trägt. Blockierend
sind nicht der Schnitt, sondern drei Klassen von Vertragsdefekten:

1. **Der Vertrag kann die Fälle nicht ausdrücken, mit denen er begründet wird**
   (F-1 für die entscheidende Messzeile 10, F-2 für Messzeile 8, F-12 für
   Messzeile 11). Das ist der schwerste Befund, weil die Begründung des
   Modul-Schnitts und die Abnahme des Ergebnisses auseinanderfallen: der
   Adopter behielte am Ende genau die Skripte, deren Ablösung den Antrag
   ausgelöst hat.
2. **Zusagen ohne prüfbaren Anker** (F-7 Preset-Kopplung, F-9 Sammel-Code auf
   einem instabilen Feld, F-15 unbestimmtes Befund-Feld) — jeweils Stellen, an
   denen ein Gate grün laufen kann, ohne die Aussage geprüft zu haben.
3. **Innere Widersprüche zwischen den beiden Spec-Straten** (F-5 Default vs.
   Validierungsgrenze, F-6 Ventil vs. Nullmengen-Guard, F-11 „erster Treffer")
   sowie Behauptungen, die der eigene Bestand widerlegt (F-8, F-3, F-4, F-13).

**Zeitkritisch:** F-3, F-8 und F-13 landen mit dem `Accepted`-Übergang in einem
immutablen Dokument ([`AGENTS.md`](../../AGENTS.md) §3.5); sie sind **vor** der
Statusänderung zu klären, nicht danach.

**Übergabe:** Die Findings gehen an den Implementer (Rückkante Review → Plan:
F-1, F-2 und F-12 berühren Abnahme-Punkt 1, 4 und 5 des Slice-Plans, nicht nur
den Text). Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und
von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser
Skill, dieses Modell, dieses Verdikt) und ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat.
