# Slice slice-293: Anforderung, Spezifikation und Folge-ADR — `exempt-tables` nimmt eingerückte Tabellenzeilen aus; Entscheide festschreiben

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-97-linelength-eingerueckte-tabellen`](welle-97-linelength-eingerueckte-tabellen.md).

**Bezug:** [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in)
(die Ausnahme (b) und das Kriterium „Ausnahme (b), Grenzen“ werden geändert, neue Kriterien kommen hinzu),
[`ADR-0114`](../../adr/0114-linelength-modul.md) (Entscheidung 5 (b) und die Konsequenz zur
eingerückten Tabellenzeile werden durch die Folge-ADR teilweise abgelöst),
[`ADR-0054`](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) (Entscheidung 1, 2 und 4:
andere Antwort oder andere Frage, je Konsument eine Assertion).

**Berührte Spec-Stellen:** [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in)
(Absatz „Drei Ausnahmen, je einzeln schaltbar“, Abschnitt „Was Fenced-Code und was eine Tabellenzeile
ist …“ samt Einleitungssatz und dem Punkt „Tabellenzeile heißt …“, Akzeptanzkriterien (b), Out-of-Scope
(4), Version und Historie); `spezifikation.md` §[`DC-FA-LLEN-001.a`](../../../../spec/spezifikation.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength)
Schritt 4, „Nicht Gegenstand“ und „Grenzen, am Gegenstand gelesen“, dazu die Zeile
`linelength.exempt-tables` im Schema [`SPEC-005`](../../../../spec/spezifikation.md#spec-005--d-checkyml);
`architecture.md` bleibt unberührt (sie nennt `linelength` nur in der Modul-Aufzählung). **Lastenheft
§5** ([Globale Out-of-Scope-Punkte](../../../../spec/lastenheft.md#5-globale-out-of-scope-punkte): die
Zeilenlänge läuft „mit derselben Markdown-Lexik wie die übrigen Module, die Markdown lesen“) ist
**geprüft und bleibt unverändert**: die Antwort kommt aus der geteilten Lexik, auf eine eigene,
benannte Frage ([`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4;
[Verdikt](../../../reviews/2026-11-24-welle-97-architect-verdikt.md) F-3). Wanderte der Helfer in
`linelength`, wäre §5 berührt, und das wäre eine Entscheidung des Auftraggebers. Der Verweis
zeigt **aufwärts**: Die Spec nennt diesen Slice nie.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-11-24.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die geänderte Anforderung
[`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in)
steht im Lastenheft, Verfahren und Schema in der Spezifikation, die Entscheidung in einer Folge-ADR zu
[`ADR-0114`](../../adr/0114-linelength-modul.md) (`Proposed` bis zur Closure von
[`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md)) — und die vier Entscheide
unten sind **vor dem ersten Code** vom Auftraggeber bestätigt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Tests, Konfigurations-Rand** — gehören
  [`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md); dieser Slice schreibt, was der
  Code einlösen soll, damit die Reviews dort gegen eine Anforderung prüfen.
- **Das Benutzerhandbuch (§4.28 und die Modul-Tabelle)** — es bildet den Ist-Zustand ab
  und beschreibt keine Fassung vor ihrem Code; es ändert sich mit dem Code in
  [`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md).
- **Der Text von [`ADR-0114`](../../adr/0114-linelength-modul.md)** — eine `Accepted`-ADR wird nicht
  überschrieben ([`AGENTS.md`](../../../../AGENTS.md) §3.5); die Folge-ADR löst Entscheidung 5 (b) und
  die Konsequenz zur eingerückten Tabellenzeile teilweise ab
  ([`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md)). Die Geschichte-Zeile in
  [`ADR-0114`](../../adr/0114-linelength-modul.md) und ihre Index-Zeile („teil-superseded: [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md)“) tragen den Übergang von [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) nach
  `Accepted` und gehören deshalb [`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md)
  ([Verdikt](../../../reviews/2026-11-24-welle-97-architect-verdikt.md) F-6).
- **Lastenheft §5** — geprüft und unverändert (Kopf „Berührte Spec-Stellen“): kein Satz dort ändert
  sich, solange der Helfer in der geteilten Lexik wohnt.
- **Zitate (`> |`), ein eigener Schalter, eine Listen-Erkennung** — Frage 2, 1 und 3 der Welle; die
  kleinste Antwort lässt sie weg, und eine Anforderung dafür wäre ein eigener Vorgang.
- **CHANGELOG, READMEs** — Release-Prep ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17), kein
  Feature-Commit.

Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in)
      ist im Lastenheft geändert (Version gehoben, Historie-Zeile nach
      [`MR-032`](../../../../harness/conventions.md#mr-032), **eine Zeile je Version**): die
      Ausnahme (b) und der Punkt „Tabellenzeile heißt …“ nennen die Einrückung bis drei
      Leerzeichen; die Zusage „dieselbe Antwort … keine zweite“ und der Out-of-Scope-Satz (4) sind
      **umgeschrieben**, nicht neben einem Widerspruch stehen gelassen; das Kriterium „Ausnahme (b),
      Grenzen“ ist **geändert**, nicht ein zweites daneben gestellt; die sechs Kriterien aus §3
      stehen da, jede Grenze mit ihrer Bedingung. **Der Punkt „Tabellenzeile heißt …“ nennt drei
      Grenzen, je mit ihrer Bedingung** (keine neuen Kriterien; Kriterium 2 deckt die Schwelle bei vier
      schon; Wortlaut in [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md)
      §Konsequenzen, [Verdikt](../../../reviews/2026-11-24-welle-97-architect-verdikt.md) F-4):
      (1) eine Prosa-Zeile, die nach ein bis drei Leerzeichen mit `|` beginnt, ohne dass eine Tabelle
      dasteht, ist ausgenommen — eine Trennzeile wird weiter nicht verlangt, wie heute in Spalte 0;
      (2) eine Tabelle auf der Zeile der Listenmarke (`- | a | b |`) beginnt nicht mit Leerzeichen
      und `|` und bleibt gemessen; (3) eine Tabelle unter einer Marke von vier Spalten (`10. `) oder
      in einer verschachtelten Liste steht mit vier und mehr Leerzeichen und bleibt gemessen.
- [x] Die Spezifikation führt das Verfahren (Schritt 4, „Nicht Gegenstand“, „Grenzen, am Gegenstand
      gelesen“) und die Beschreibung von `linelength.exempt-tables` im Schema; jede Aussage über die
      Grenzen ist am Gegenstand gelesen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 13), nicht
      aus der Beschreibung übernommen.
- [x] Die Folge-ADR (Entscheide samt verglichenen Alternativen, darunter „`tableRowLine` selbst
      weiten“ und „eigener Schalter“) liegt als `Proposed` vor, ist im
      [ADR-Index](../../adr/README.md) geführt und trägt `## Re-Evaluierungs-Trigger`.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in): die Sätze zu Ausnahme (b), die Zusage „dieselbe Antwort … keine zweite“, der Punkt „Tabellenzeile heißt …“, die Kriterien „Ausnahme (b), …“, Out-of-Scope (4), Version und Historie (Spiegel-Zählung unten) |
| `spec/spezifikation.md` | update | §[`DC-FA-LLEN-001.a`](../../../../spec/spezifikation.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength) Schritt 4, „Nicht Gegenstand“, „Grenzen, am Gegenstand gelesen“, die Beschreibung von `linelength.exempt-tables` im Schema, Historie |
| Folge-ADR zu [`ADR-0114`](../../adr/0114-linelength-modul.md) | liegt vor: [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md), `Proposed`, vom Architect mit dem Verdikt geschrieben | Entscheidung 5 (b) und die Konsequenz zur eingerückten Tabellenzeile werden teilweise abgelöst (`Supersedes`-Teil); der Architect gleicht sie an, wenn das Schreiben der Anforderung eine Entscheidung verschiebt |
| **Lastenheft §5** | geprüft, **unverändert** | die Antwort kommt aus der geteilten Lexik ([`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4); die Zusage dort „mit derselben Markdown-Lexik wie die übrigen Module“ bleibt wahr ([Verdikt](../../../reviews/2026-11-24-welle-97-architect-verdikt.md) F-3) |
| `docs/plan/adr/README.md` | update | ADR-Index ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 4) |
| Plan von [`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md) | update | die Testfälle, die aus den bestätigten Entscheiden folgen |

**Die kleinsten Antworten stehen einmal** — im Wortlaut der Welle
([§5](welle-97-linelength-eingerueckte-tabellen.md), Fragen 1 bis 4), mit der Stelle, die sie trägt.
Dieser Plan wiederholt sie nicht, damit keine zweite Fassung neben der ersten steht
([`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md));
der Auftraggeber hat die Fragen 1 bis 3 am 2026-11-24 bestätigt. **Der eine Punkt, der vom Wortlaut
des Issues abweicht, ist Frage 3:** höchstens drei Leerzeichen statt „2 bis 5“, weil die geteilte
Lexik keinen Listenkontext hergibt. Frage 4 ist im
[Architect-Verdikt](../../../reviews/2026-11-24-welle-97-architect-verdikt.md) entschieden.

**Die sechs Akzeptanzkriterien, die die Anforderung tragen soll** (Obergrenze der Welle: etwa acht;
jedes mit der Bedingung, unter der es gilt — der Lastenheft-Text und das Kriterium in
[`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md) §3 müssen **dieselbe** Bedingung
nennen):

| # | Kriterium (Given … then) | Änderung am Bestand |
|---|---|---|
| 1 | **Ausnahme (b), eingerückt, an und aus.** Given eine überlange `\|`-Zeile mit ein bis drei Leerzeichen davor, when der Lauf im Default erfolgt, then kein Befund; with `exempt-tables: false`, then ein Befund. | neu |
| 2 | **Ausnahme (b), Grenzen.** Given je eine überlange `\|`-Zeile mit vier Leerzeichen davor, mit einem Tab davor, mit einem geschützten Leerzeichen (U+00A0) davor, eine `> \|`-Zeile und eine Zeile `a \| b` ohne führende Pipe, when der Lauf im Default erfolgt, then je ein Befund — „Leerzeichen“ ist hier U+0020 und kein anderes Zeichen — anders als der Leerraum der Wortgrenze in Ausnahme (c), zu dem auch der Tab gehört. | **ersetzt** das heutige Kriterium „Ausnahme (b), Grenzen“ (dessen erster Fall „eingerückte `\|`-Zeile ⇒ Befund“ nicht mehr gilt) |
| 3 | **Ausnahme (b) im Zaun.** Given `exempt-code-blocks: false`, `exempt-tables: true` und eine überlange `\|`-Zeile mit drei Leerzeichen davor in einem Zaun, then ein Befund. | neu (das Gegenstück ohne Einrückung steht im heutigen Kriterium und bleibt in 2) |
| 4 | **Bezugsmenge.** Given eine Datei, deren eine Zeile mit drei Leerzeichen vor `\|` beginnt, when der Lauf im Default erfolgt, then zählt die Bezugsmengen-Zeile diese Zeile nicht zu den geprüften Zeilen; with `exempt-tables: false` zählt sie. | neu — die Zeile `linelength: <D> Datei(en), <Z> Zeile(n) geprüft …` ist heute das, was nach `scope`, `exempt-paths`, (a) und (b) bleibt |
| 5 | **Andere Module unverändert.** Given eine Zeile mit ein bis drei Leerzeichen vor `\|`, then liest keines der Module `targets`, `planning.waves` und `structure` sie als Tabellenzeile (Antwort heute: nein). | neu |
| 6 | **Rest unverändert.** Given Dateien ohne eine Zeile mit ein bis drei Leerzeichen vor `\|`, when der Lauf im Default erfolgt, then nimmt Ausnahme (b) genau die Zeilen aus, die mit `\|` beginnen (außerhalb von Fenced-Code, wie bisher). | neu — die Ausgabe-Gleichheit gegen den Stand davor prüft die Welle ([§3 (a)](welle-97-linelength-eingerueckte-tabellen.md)), nicht die Anforderung |

**Spiegel-Zählung vor dem Editieren** ([`MR-025`](../../../../harness/conventions.md#mr-025), Form
vor Zahl — [`AGENTS.md`](../../../../AGENTS.md) §5 Regel 14). Die Aussage „eine eingerückte
`|`-Zeile ist keine Tabellenzeile“ und die Schalter-Beschreibung stehen **nicht** nur an den Stellen,
die der Schlüssel `exempt-tables` nennt:

```bash
# Stellen, die den Schlüssel nennen (lebender Bestand: ohne Planung, Reviews, CHANGELOG)
git grep -c 'exempt-tables' -- . ':!docs/plan/planning' ':!docs/reviews' ':!CHANGELOG.md'            # 10 Dateien, 27 Zeilen am 2026-11-24
# Stellen, die die Aussage ohne den Schlüssel tragen
git grep -n -E 'eingerückte[rn]? (Tabellenzeile|`\|`)|eingerückte oder zitierte|dieselbe Antwort|keine zweite|zweite Antwort' \
  -- spec/lastenheft.md spec/spezifikation.md docs/user/benutzerhandbuch.md docs/plan/adr/0114-linelength-modul.md
```

Von den 27 Zeilen sind 4 Beschreibungen, die sich mit diesem Slice ändern (Lastenheft: Ausnahme (b)
in der Beschreibung und das Kriterium „(b), Grenzen“; Spezifikation: Schritt 4 und die
Schema-Zeile); das Kriterium „(b), an und aus“ bleibt, Kriterium 1 steht daneben. Die übrigen nennen den Schlüssel am Konfigurations-Rand, in Tests, im
Image-Test, im Handbuch (zwei Beschreibungen: [`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md))
oder in der `Accepted`-ADR (unveränderlich). Die zweite Suche findet die Aussage ohne den
Schlüssel: im Lastenheft die Zusage „dieselbe Antwort … keine zweite“, den Punkt „Tabellenzeile
heißt …“ und Out-of-Scope (4); in der Spezifikation die „Grenzen, am Gegenstand gelesen“; im
Handbuch §4.28 den Satz zur eingerückten `|`-Zeile (Slice 294); in der ADR Kontext und Konsequenzen
(unveränderlich). Der Implementer zählt vor dem Editieren **neu** — die Zahlen sind ein Stand, kein
Vertrag.

### Plan-Nachtrag des Implementers (2026-11-24, vor dem Inhalts-Commit)

**Gemessen, bevor editiert wurde** ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 13 und 14): was die
Lexik heute für eine eingerückte `|`-Zeile liefert. Quelle des Messbilds ist das gebaute Image
(`d-check:latest`, erzeugt 2026-11-24T06:50:53+01:00; `markdown.go` und `linelength.go` ändern sich zuletzt
am 2026-11-23, der Code des Moduls ist im Image also der Stand von `HEAD`). Eine Datei mit je einer überlangen
Zeile (`max-chars: 10`, `exempt-unbreakable: false`, damit nur (b) antwortet):

```bash
# Fixture: ein Verzeichnis mit .d-check.yml (linelength: max-chars 10, exempt-unbreakable false)
# und f.md mit je einer Zeile  `| a | b | c |`  mit 1, 2, 3, 4 Leerzeichen, Tab, U+00A0, `> `,
# `a | b | c | d | e`, `- | a | b | c |`, `10. | a | b | c |`, Spalte 0, `   | prosa ohne tabelle`
# und einem Zaun mit zwei `|`-Zeilen (drei Leerzeichen, Spalte 0).
docker run --rm -v "$FIXTURE":/repo:ro d-check:latest -enable linelength
```

| Zeile | Default (heute) | `exempt-tables: false` | `exempt-code-blocks: false` (heute) |
|---|---|---|---|
| 1, 2, 3 Leerzeichen vor `\|` | **Befund** | Befund | Befund |
| 4 Leerzeichen, Tab, U+00A0 vor `\|` | Befund | Befund | Befund |
| `> \|`, `a \| b \| c`, `- \| a \| b \|`, `10. \| a \| b \|` | Befund | Befund | Befund |
| Spalte 0 (`\| a \| b \| c \| d \|`) | kein Befund | Befund | kein Befund außerhalb, im Zaun Befund |
| Prosa-Zeile, drei Leerzeichen, `\|` ohne Tabelle | **Befund** | Befund | Befund |
| Prosa-Zeile `\| prosa ohne tabelle` in Spalte 0 | kein Befund | — | — |
| drei Leerzeichen und `\|` im Zaun | kein Befund (Zaun) | kein Befund (Zaun) | Befund |

Gelesen: Heute misst `linelength` jede eingerückte `|`-Zeile (erste und fünfte Tabellenzeile; diese
Fälle ändern sich mit der Anforderung), und die Spalte-0-Prosa-Zeile ist schon heute ausgenommen — Grenze 1 („wie
heute in Spalte 0“) stimmt. Die Zeilen mit vier Leerzeichen, Tab, U+00A0, `> `, ohne führende Pipe und
`- |` bleiben nach der Änderung gemessen und ändern ihr Messbild nicht. Die Bezugsmenge nennt am Fixture
25 geprüfte Zeilen im Default und 26 mit `exempt-tables: false` (die Spalte-0-Zeile). Die Kopplung der
anderen Module steht in `lexikon_kopplung_test.go` (`tabellenFaelle`, Fall „eingerueckt“ mit zwei
Leerzeichen, Antwort „keine Tabellenzeile“ für alle vier Konsumenten); ein eigenes Messbild für `targets`,
`planning.waves` und `structure` entsteht dort mit `slice-294`, nicht hier. **Nicht gemessen:** dass eine
Tabelle unter `10. ` oder in einer verschachtelten Liste mit vier und mehr Leerzeichen steht — das ist die
Aussage der ADR und des Verdikts über CommonMark, kein Messergebnis des Werkzeugs.

**Zählung vor dem Editieren** (zweite Suche aus §3, neu gezählt am 2026-11-24): `git grep -c 'exempt-tables'`
mit den Ausschlüssen aus §3 liefert **12 Dateien, 35 Zeilen**; die beiden Mehrzeilen gegen die 10 Dateien,
27 Zeilen des Plans sind [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) (7) und der
ADR-Index (1), beide seit dem Plan hinzugekommen. Die Stände der anderen Dateien decken sich mit dem Plan.

**Wortlaut gegen die Plan-Fassung** (keine Plan-Änderung, zwei Ausformungen):

- Kriterium 2 behält den zweiten Fall des heutigen Kriteriums (Spalte-0-`|`-Zeile im Zaun bei
  `exempt-code-blocks: false`, ⇒ Befund): „das Gegenstück ohne Einrückung bleibt in 2“ (§3, Kriterium 3).
- Out-of-Scope (4) nimmt zusätzlich die **Listen-Erkennung** für die Einrückung vor `|` auf, mit der Bedingung
  „die geteilte Lexik kennt keinen Listenkontext“ ([`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md)
  Alternative D und Entscheid 3). Die Fassung wiederholt, was §1 und die Alternativen der ADR schon ausschließen.
- Die drei Grenzen stehen als Sätze im Punkt „Tabellenzeile heißt …“ und in der Spezifikation („Grenzen, am
  Gegenstand gelesen“), **ohne** neue Kriterien (F-4); die Fälle `- | a | b |` und die Prosa-Zeile prüft
  `slice-294` als Zusatz zu Kriterium 2.
- Gegen Lastenheft §5 gehalten ([`BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft`](../observations/BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft/observation.md)):
  keiner der sechs Ausschlüsse (Architektur-/Build-/Mathematik-Prüfung, Auto-Fix, Nicht-Markdown-Formate,
  Rechtschreib-/Stil-/Lint-Prüfung) trifft die Änderung; die Ausnahme „mit derselben Markdown-Lexik wie die
  übrigen Module, die Markdown lesen“ bleibt wahr, weil die Frage in der geteilten Lexik wohnt (Verdikt F-3).

### Plan-Nachtrag des Implementers (2026-11-24, nach dem Inhalts-Commit)

**Die Aussage, über ihren Inhalt gesucht** (nicht an der zitierten Stelle; [`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md)),
vor und nach dem Editieren mit demselben Kommando:

```bash
git grep -n -E 'eingerückte[rn]? (Tabellenzeile|`\|`)|eingerückte oder zitierte|dieselbe Antwort|keine zweite|zweite Antwort|mit `\\?\|` beginn|Spalte 0|Spalte-0' \
  -- spec/lastenheft.md spec/spezifikation.md
```

| Stelle (vor dem Editieren) | Aussage | Nach dem Editieren |
|---|---|---|
| Lastenheft, Absatz „Drei Ausnahmen“, (b) | „… die mit `\|` beginnen, werden nicht gemessen“ | umgeschrieben: „nach null bis drei Leerzeichen (U+0020) mit `\|` beginnen“ |
| Lastenheft, Einleitungssatz des Abschnitts | „dieselbe Antwort … keine zweite“ | umgeschrieben: gilt für Fenced-Code; die Tabellenzeile ist eine eigene, benannte Frage der Lexik |
| Lastenheft, Punkt „Tabellenzeile heißt …“ | eingerückte Zeile wird gemessen | umgeschrieben, drei Grenzen mit Bedingung |
| Lastenheft, Kriterium „Ausnahme (b), Grenzen“ | eingerückte `\|`-Zeile ⇒ Befund | **ersetzt** (kein zweites daneben) |
| Lastenheft, Out-of-Scope (4) | „eigene, strengere Lesart … von Tabellen … zweite Antwort“ | umgeschrieben, Listenkontext aufgenommen |
| Spezifikation, Schritt 4 | „`tableRowLine` gilt“ | umgeschrieben: die Frage von `linelength` schließt `tableRowLine` ein |
| Spezifikation, „Nicht Gegenstand“ | „eine eigene Lesart von … Tabellen“ | umgeschrieben (strengere Lesart, Listenkontext) |
| Spezifikation, „Grenzen, am Gegenstand gelesen“ | „eine eingerückte oder zitierte Tabellenzeile … nicht“ | umgeschrieben, Bedingung je Grenze |
| Spezifikation, Schema `linelength.exempt-tables` | „die mit `\|` beginnen“ | umgeschrieben |
| Lastenheft 4356 (`metrics`), Spezifikation 2419, 2801, 3040, 3555, 4134 (`targets`, `planning.waves`, `structure`, `metrics`) | Spalte 0, „dieselbe Antwort wie bei `targets`“ | **unberührt, bleibt wahr**: sie beschreiben `tableRowLine` und die anderen Module; beides ändert sich nicht |
| Lastenheft 3246 | „keine zweite Regel“ (anderer Gegenstand) | unberührt |
| Handbuch §4.28 (drei Stellen), Modul-Tabelle, Änderungs-Tabelle | Ist-Zustand zu `exempt-tables` | gehört [`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md) |
| `config_template.go`, READMEs | Kommentar bzw. Modul-Bullet | `slice-294` bzw. Release-Prep |
| [`ADR-0114`](../../adr/0114-linelength-modul.md) | Entscheidung 5 (b), Konsequenz | unveränderlich (`Accepted`) |

Nach dem Editieren zeigt dieselbe Suche nur noch Treffer, die entweder die **umgeschriebene** Fassung sind
(Lastenheft Z. 4604, 4618–4622, 4642, 4724, 4743, 4747; Spezifikation Z. 3785, 4129), eine **unberührte** Zeile
der Tabelle oben oder eine **Historie-Zeile früherer Fassungen** (Chronik, nicht umgeschrieben); keine Zeile trägt
mehr die alte Aussage „eine eingerückte Tabellenzeile wird gemessen“ oder „keine zweite Antwort“ für die
Tabellenzeile. Das Risiko „eine Stelle stünde gegen die andere“ (§6, Risiko 4) ist damit am Text gelesen; ob
das Review einen Widerspruch findet, entscheidet Risiko 1.

**Kriterien-Zählung:** sechs Kriterien in der Anforderung, davon fünf neu und eines ersetzt (kein zweites daneben
gestellt); das Kriterium „Ausnahme (b), an und aus“ bleibt. Zuordnung zu den Entscheiden der ADR:
Entscheid 1 (kein neuer Schlüssel; `false` misst) — Kriterium 1 und 4; Entscheid 2 (`> |` gemessen) —
Kriterium 2; Entscheid 3 (höchstens drei Leerzeichen; Tab, U+00A0 und vier gemessen; im Zaun gemessen) —
Kriterium 1, 2 und 3; Entscheid 4 (eigene Frage in der Lexik, andere Module unverändert) — Kriterium 5 und 6 und
die umgeschriebenen Sätze. Die drei Grenzen und die Prosa-Zeile stehen als Sätze im Punkt „Tabellenzeile
heißt …“ ohne eigenes Kriterium (F-4).

### Plan-Nachtrag des Implementers (2026-11-24, nach Review R1, Befund F-1)

Der Befund ([Review R1](../../../reviews/2026-11-24-slice-293-linelength-eingerueckte-tabellen-r1.md) F-1)
trifft eine Aussage, die ich ohne ihre Bedingung geschrieben hatte: „die übrigen Module lesen eine
eingerückte Zeile weiter nicht als Tabellenzeile“. Am Gegenstand gelesen: `declaredObservationIDs` in
`planning_observations.go` trimmt die Zeile (`strings.TrimSpace`) vor dem Test auf `|`, und
`SplitPipeTableLine` in `markdown_table.go` trimmt (`strings.TrimSpace`) — beide lesen eine eingerückte
Zeile als Tabellenzeile; die Aussage gilt nur für die drei Konsumenten von `tableRowLine`. Die Messung
unter `--trace` ist die des Reviews (am eigenen Fixture ohne die Konfiguration des Leser-Musters
nicht nachgestellt). Keine Plan-Änderung über §3 hinaus: die Kürzung wiederholt die Eingrenzung, die
Kriterium 5 und die Spezifikation schon tragen.

**Die Aussage, über ihren Inhalt gesucht**, vor dem Editieren:
`git grep -n -E 'übrigen Module|weiter nicht als Tabellenzeile|andere Module|anderen Module|übrige Module' -- spec/lastenheft.md spec/spezifikation.md`.

| Stelle | Aussage | Ergebnis |
|---|---|---|
| Lastenheft, Absatz zur geteilten Lexik (3 Stellen) | „fragt `linelength` etwas anderes als die übrigen Module“; „Spalte-0-Antwort der übrigen Module“; „die übrigen Module lesen eine eingerückte Zeile weiter nicht …“ | **auf `targets`, `planning.waves` und `structure` gekürzt** |
| Lastenheft, Historie 0.115.0 | derselbe Satz | **bleibt byte-gleich** (Chronik, nicht umgeschrieben); die Präzisierung steht in der neuen Zeile 0.115.1 |
| Spezifikation, Schritt 4 | „Antwort von `tableRowLine`, die die übrigen Module lesen“ | **auf die drei Module gekürzt** |
| Spezifikation, Historie (Zeile von 0.115.0) | nennt schon die drei Module | unberührt; eigene Zeile für den Nachzug |
| Lastenheft Kriterium „andere Module unverändert“ | nennt die drei Module | unberührt, trug die Eingrenzung schon |
| Lastenheft Absatz, „für Fenced-Code … die übrigen Module“ | Fenced-Code, nicht Tabellenzeile | unberührt: anderer Gegenstand, vom Review nicht beanstandet |
| Lastenheft 1573, 1613, 1699, 1807, 1857, 1929, 4785, 4812; Spezifikation 90, 971, 1424, 1751–1842, 2443, 3617, 3674, 3766, Historie 4427, 4445 | andere Gegenstände (Opt-out-Marker, Vorverarbeitung, Eingabefehler, §5) | unberührt |

Version: Lastenheft 0.115.1 mit **eigener** Historie-Zeile „präzisiert“, die Zeile 0.115.0 bleibt byte-gleich
([`MR-032`](../../../../harness/conventions.md#mr-032); die Lehre aus welle-96 R2 L-2).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die Welle
[`welle-97-linelength-eingerueckte-tabellen`](welle-97-linelength-eingerueckte-tabellen.md) ist
offen, das Architect-Verdikt zu Frage 4 liegt unter `docs/reviews/`, und das WIP-Limit ist frei
(`in-progress/` trägt nur die Roadmap). Der Auftraggeber hat die Fragen 1 bis 3 der Welle
bestätigt oder geändert, **bevor** die Anforderung geschrieben wird.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): das Verdikt verlangt, `tableRowLine`
  selbst zu weiten, oder die Bestätigung weitet Frage 3 (Einrückung ab vier Leerzeichen mit
  Listen-Erkennung) — beides trüge mehr als die sechs Kriterien und gehört in einen neu
  geschnittenen Plan, nicht in diesen.
- `in-progress` → `open` (blockiert): der Auftraggeber widerruft die Antwort zu Frage 1, 2 oder 3;
  ohne tragenden Entscheid beginnt kein Code.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Lastenheft, Spezifikation und Folge-ADR liegen vor und `make gates` ist grün),
Review- und Verifikations-Report liegen unter `docs/reviews/`, Closure-Notiz mit Lerneintrag
geschrieben. **Review-Vorgabe der Welle:** höchstens zwei Runden, eingearbeitet werden nur HIGH und
MEDIUM; jede Nachbesserung zitiert die Stelle, die sie trägt
([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 18).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Anforderung sagt an mehreren Stellen, eine eingerückte Tabellenzeile werde gemessen, und
  verspricht „dieselbe Antwort … keine zweite“ (Absatz „Was Fenced-Code und was eine Tabellenzeile
  ist …“, Punkt „Tabellenzeile heißt …“, Kriterium „Ausnahme (b), Grenzen“, Out-of-Scope (4));
  die Änderung muss sie umschreiben, nicht daneben eine widersprechende Zeile stellen
  ([`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`](../observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/observation.md),
  verkörpert, 3×) — **Ausgang:** *entfallen.* Review R1 las die Stellen nebeneinander und fand
  keinen Widerspruch zwischen ihnen; sein einziger Befund (F-1) war ein Satz, der über den
  Gegenstand hinausging, kein Widerspruch (Verifikation, Risiko 1).
- Die Anforderung wächst über sechs Kriterien (Obergrenze der Welle: etwa acht), weil eine
  Randbedingung Randfälle mitbringt
  ([`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](../observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md),
  1×) — **Ausgang:** *entfallen.* Das Lastenheft führt sieben Kriterien der Familie (b) bei der
  Obergrenze acht, die Folge-ADR vier Entscheide (Verifikation, Risiko 2).
- Die enge Antwort zu Frage 3 (höchstens drei Leerzeichen) lässt eine Tabelle unter einer
  Vier-Spalten-Marke oder in einem verschachtelten Listenpunkt als Befund stehen; ob der Bestand des
  Melders solche Zeilen trägt, ist von hier nicht gemessen — **Ausgang:** *weiter offen.* Der
  Auftraggeber hat die enge Fassung am 2026-11-24 bestätigt, die Anforderung nennt die Grenze mit
  ihrer Bedingung. Der Bestand des Melders bleibt ungemessen; trägt er nach der Hebung Zeilen dieser
  Form, ist die Listen-Erkennung ein eigener Vorgang. Register:
  [`BEO-ALL/enge-einrueckungs-grenze-laesst-listen-tabellen-stehen`](../observations/BEO-ALL/enge-einrueckungs-grenze-laesst-listen-tabellen-stehen/observation.md).
- Lastenheft und Spezifikation führen dieselbe Aussage je in eigenen Worten; ein Satz, der nur an einer
  Stelle umgeschrieben wird, stünde gegen den anderen — **Ausgang:** *entfallen.* Die zweite Suche
  nach dem Editieren zeigt keine der genannten Stellen mehr mit der alten Aussage; jede ist
  umgeschrieben oder als unberührt begründet (Verifikation, Risiko 4).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** Der Zuschnitt aus welle-96 trug. Die Anforderung blieb bei sieben
  Kriterien der Familie (b) und vier Entscheiden, die Review-Runde fand einen Befund (einen zu weiten
  Satz) statt vier, und es genügte eine Runde. Die Messung am Image vor dem Schreiben (Fixture mit
  `max-chars: 10`) hat den Ist-Zustand festgelegt, statt ihn zu vermuten; Reviewer und Verifier haben
  sie unabhängig nachgefahren. Der Architect entschied Frage 4 am Code: `tableRowLine` wird nicht
  geweitet, weil `targets`, `planning.waves` und `structure` sonst gegen ihre Zusage „Spalte 0“
  rutschten.
- **Was ging anders als geplant:** **(1) Ein Satz ging über den Gegenstand hinaus** („die übrigen
  Module lesen eine eingerückte Zeile weiter nicht als Tabellenzeile“): der RTM-Leser unter
  `--trace` und der Beobachtungs-Leser in `planning` trimmen die Zeile und lesen sie als Tabelle.
  Gekürzt auf die drei Module; das Lastenheft steht seitdem auf 0.115.1, die Zeile 0.115.0 blieb
  unverändert. **(2) Die Spiegel-Zählung wuchs** von 10 Dateien und 27 Zeilen auf 12 und 35, weil
  [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) und ihr Index-Eintrag nach dem Plan dazukamen. **(3) Die Zählung der „Ausformungen“**
  stand in Plan-Nachtrag und Botschaft verschieden (zwei und drei, R1 F-2).
- **Steering-Loop-Eintrag:** mit diesem Slice ist keine Regel, kein Sensor und keine benannte
  Spec-Lücke verkörpert.
  [`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md)
  erreicht mit diesem Slice 2×.
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-293.md` in
  [`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md)
  (2×) und in
  [`BEO-ALL/commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/observation.md)
  (schon verkörpert); neu angelegt, 1×:
  [`BEO-ALL/enge-einrueckungs-grenze-laesst-listen-tabellen-stehen`](../observations/BEO-ALL/enge-einrueckungs-grenze-laesst-listen-tabellen-stehen/observation.md).
  Benannt, nicht gezählt: die Schema-Zeile von `linelength.exempt-tables` führt eine kürzere
  Aufzählung der Grenzen als Lastenheft und Spezifikation (R1 F-3); die Historie-Zeile 0.115.0 trägt
  den zu weiten Satz weiter und bleibt als Chronik stehen, die Zeile 0.115.1 stellt ihn richtig (V-2);
  der Satz „für Fenced-Code … die übrigen Module“ steht ohne Eingrenzung, betrifft aber eine andere
  Frage (V-1); die drei Grenzen und die Prosa-Zeile haben keine eigenen Kriterien, ihre Tests stehen
  in [`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md) §3 (V-3, Verdikt F-4).
- **Folge-Slices:** keine neuen. [`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md)
  trägt Helfer, Tests je Kriterium, Handbuch und den Übergang von
  [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) auf `Accepted`.
- **Risiken aus §6:** drei *entfallen*, eines *weiter offen* — siehe §6.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist der Default `*` (`ALL`) —
Lastenheft, Spezifikation, ADR; keine Sub-Area mit eigener Modus-Deklaration
(`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stand 2026-11-24,
`ls docs/plan/planning/observations/BEO-ALL/`), Sub-Area `*`. Zähler = Dateien unter `evidence/`.
Treffer:
[`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`](../observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/observation.md)
(verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §5 Regel 18, 3× — Risiko 1 in §6; die
Festlegungen stehen in der Welle mit ihrer Stelle),
[`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](../observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md)
(offen, 1× — Risiko 2; der Zuschnitt dieser Welle folgt ihrem Ableiter),
[`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md)
(offen, 1× — jede Grenze der Anforderung trägt ihre Bedingung, §3),
[`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md)
(offen, 2× — jede Randbedingung bekommt ein Kriterium, das [`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md) mit einem Test belegt),
[`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md)
(offen, 2× — „dieselbe Lexik“ stimmt nur für eine Familie; Frage 4 der Welle),
[`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md)
(verkörpert in welle-74, 6× — der Kopplungs-Test, Kriterium 5),
[`BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft`](../observations/BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft/observation.md)
(offen, 1× — das Review liest die geänderte Anforderung gegen Lastenheft §5, dessen Ausnahme „mit
derselben Markdown-Lexik wie die übrigen Module“ sagt) und
[`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/observation.md)
(verkörpert, 4× — die Spiegel-Suche in §3 nimmt Handbuch und Konfigurations-Vorlage auf; sie gehören
[`slice-294`](../done/slice-294-linelength-eingerueckte-tabellen-kern.md)). **Kein Treffer:**
[`BEO-ALL/prosa-verschwindet-in-inline-code-spanne`](../observations/BEO-ALL/prosa-verschwindet-in-inline-code-spanne/observation.md)
(offen, 2×) — `linelength` misst die rohe Zeile über `proseLineSet`, das nur Fenced-Code
ausschließt und keine Inline-Code-Spanne leert; der Zähler bleibt unberührt.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-24: `image-scan` **grün** (jüngster Lauf 2026-11-23T12:51:15Z).
`upstream-drift` **rot** (jüngster Lauf 2026-11-23T08:29:02Z), und dieser Lauf liegt **vor** den
Pin-Hebungen vom 2026-11-23; sein Log war über `gh run view <id> --log` nicht abrufbar (leere
Ausgabe) — ob die Meldung planmäßige Fremd-Drift ist, ist damit **nicht gelesen**, nur nicht
neuer als die Hebungen. Beim Start dieses Slice neu lesen. Berührt dieses Modul nicht (die Läufe
messen Pins und Images, nicht `linelength`).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
