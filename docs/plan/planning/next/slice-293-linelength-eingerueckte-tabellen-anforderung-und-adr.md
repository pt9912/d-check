# Slice slice-293: Anforderung, Spezifikation und Folge-ADR — `exempt-tables` nimmt eingerückte Tabellenzeilen aus; Entscheide festschreiben

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-97-linelength-eingerueckte-tabellen`](../welle-97-linelength-eingerueckte-tabellen.md).

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
[`slice-294`](../open/slice-294-linelength-eingerueckte-tabellen-kern.md)) — und die vier Entscheide
unten sind **vor dem ersten Code** vom Auftraggeber bestätigt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Tests, Konfigurations-Rand** — gehören
  [`slice-294`](../open/slice-294-linelength-eingerueckte-tabellen-kern.md); dieser Slice schreibt, was der
  Code einlösen soll, damit die Reviews dort gegen eine Anforderung prüfen.
- **Das Benutzerhandbuch (§4.28 und die Modul-Tabelle)** — es bildet den Ist-Zustand ab
  und beschreibt keine Fassung vor ihrem Code; es ändert sich mit dem Code in
  [`slice-294`](../open/slice-294-linelength-eingerueckte-tabellen-kern.md).
- **Der Text von [`ADR-0114`](../../adr/0114-linelength-modul.md)** — eine `Accepted`-ADR wird nicht
  überschrieben ([`AGENTS.md`](../../../../AGENTS.md) §3.5); die Folge-ADR löst Entscheidung 5 (b) und
  die Konsequenz zur eingerückten Tabellenzeile teilweise ab
  ([`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md)). Die Geschichte-Zeile in
  [`ADR-0114`](../../adr/0114-linelength-modul.md) und ihre Index-Zeile („teil-superseded: [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md)“) tragen den Übergang von [`ADR-0116`](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) nach
  `Accepted` und gehören deshalb [`slice-294`](../open/slice-294-linelength-eingerueckte-tabellen-kern.md)
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

- [ ] [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in)
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
- [ ] Die Spezifikation führt das Verfahren (Schritt 4, „Nicht Gegenstand“, „Grenzen, am Gegenstand
      gelesen“) und die Beschreibung von `linelength.exempt-tables` im Schema; jede Aussage über die
      Grenzen ist am Gegenstand gelesen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 13), nicht
      aus der Beschreibung übernommen.
- [ ] Die Folge-ADR (Entscheide samt verglichenen Alternativen, darunter „`tableRowLine` selbst
      weiten“ und „eigener Schalter“) liegt als `Proposed` vor, ist im
      [ADR-Index](../../adr/README.md) geführt und trägt `## Re-Evaluierungs-Trigger`.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

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
| Plan von [`slice-294`](../open/slice-294-linelength-eingerueckte-tabellen-kern.md) | update | die Testfälle, die aus den bestätigten Entscheiden folgen |

**Die kleinsten Antworten stehen einmal** — im Wortlaut der Welle
([§5](../welle-97-linelength-eingerueckte-tabellen.md), Fragen 1 bis 4), mit der Stelle, die sie trägt.
Dieser Plan wiederholt sie nicht, damit keine zweite Fassung neben der ersten steht
([`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md));
der Auftraggeber hat die Fragen 1 bis 3 am 2026-11-24 bestätigt. **Der eine Punkt, der vom Wortlaut
des Issues abweicht, ist Frage 3:** höchstens drei Leerzeichen statt „2 bis 5“, weil die geteilte
Lexik keinen Listenkontext hergibt. Frage 4 ist im
[Architect-Verdikt](../../../reviews/2026-11-24-welle-97-architect-verdikt.md) entschieden.

**Die sechs Akzeptanzkriterien, die die Anforderung tragen soll** (Obergrenze der Welle: etwa acht;
jedes mit der Bedingung, unter der es gilt — der Lastenheft-Text und das Kriterium in
[`slice-294`](../open/slice-294-linelength-eingerueckte-tabellen-kern.md) §3 müssen **dieselbe** Bedingung
nennen):

| # | Kriterium (Given … then) | Änderung am Bestand |
|---|---|---|
| 1 | **Ausnahme (b), eingerückt, an und aus.** Given eine überlange `\|`-Zeile mit ein bis drei Leerzeichen davor, when der Lauf im Default erfolgt, then kein Befund; with `exempt-tables: false`, then ein Befund. | neu |
| 2 | **Ausnahme (b), Grenzen.** Given je eine überlange `\|`-Zeile mit vier Leerzeichen davor, mit einem Tab davor, mit einem geschützten Leerzeichen (U+00A0) davor, eine `> \|`-Zeile und eine Zeile `a \| b` ohne führende Pipe, when der Lauf im Default erfolgt, then je ein Befund — „Leerzeichen“ ist hier U+0020 und kein anderes Zeichen — anders als der Leerraum der Wortgrenze in Ausnahme (c), zu dem auch der Tab gehört. | **ersetzt** das heutige Kriterium „Ausnahme (b), Grenzen“ (dessen erster Fall „eingerückte `\|`-Zeile ⇒ Befund“ nicht mehr gilt) |
| 3 | **Ausnahme (b) im Zaun.** Given `exempt-code-blocks: false`, `exempt-tables: true` und eine überlange `\|`-Zeile mit drei Leerzeichen davor in einem Zaun, then ein Befund. | neu (das Gegenstück ohne Einrückung steht im heutigen Kriterium und bleibt in 2) |
| 4 | **Bezugsmenge.** Given eine Datei, deren eine Zeile mit drei Leerzeichen vor `\|` beginnt, when der Lauf im Default erfolgt, then zählt die Bezugsmengen-Zeile diese Zeile nicht zu den geprüften Zeilen; with `exempt-tables: false` zählt sie. | neu — die Zeile `linelength: <D> Datei(en), <Z> Zeile(n) geprüft …` ist heute das, was nach `scope`, `exempt-paths`, (a) und (b) bleibt |
| 5 | **Andere Module unverändert.** Given eine Zeile mit ein bis drei Leerzeichen vor `\|`, then liest keines der Module `targets`, `planning.waves` und `structure` sie als Tabellenzeile (Antwort heute: nein). | neu |
| 6 | **Rest unverändert.** Given Dateien ohne eine Zeile mit ein bis drei Leerzeichen vor `\|`, when der Lauf im Default erfolgt, then nimmt Ausnahme (b) genau die Zeilen aus, die mit `\|` beginnen (außerhalb von Fenced-Code, wie bisher). | neu — die Ausgabe-Gleichheit gegen den Stand davor prüft die Welle ([§3 (a)](../welle-97-linelength-eingerueckte-tabellen.md)), nicht die Anforderung |

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
Image-Test, im Handbuch (zwei Beschreibungen: [`slice-294`](../open/slice-294-linelength-eingerueckte-tabellen-kern.md))
oder in der `Accepted`-ADR (unveränderlich). Die zweite Suche findet die Aussage ohne den
Schlüssel: im Lastenheft die Zusage „dieselbe Antwort … keine zweite“, den Punkt „Tabellenzeile
heißt …“ und Out-of-Scope (4); in der Spezifikation die „Grenzen, am Gegenstand gelesen“; im
Handbuch §4.28 den Satz zur eingerückten `|`-Zeile (Slice 294); in der ADR Kontext und Konsequenzen
(unveränderlich). Der Implementer zählt vor dem Editieren **neu** — die Zahlen sind ein Stand, kein
Vertrag.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die Welle
[`welle-97-linelength-eingerueckte-tabellen`](../welle-97-linelength-eingerueckte-tabellen.md) ist
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
  verkörpert, 3×) — **Ausgang:** *entfallen*, wenn das Review diese Stellen nebeneinander gelesen
  und keinen Widerspruch gefunden hat; sonst *eingetreten* mit dem Fix-Commit als Beleg und einer
  Evidenz-Datei im Register.
- Die Anforderung wächst über sechs Kriterien (Obergrenze der Welle: etwa acht), weil eine
  Randbedingung Randfälle mitbringt
  ([`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](../observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md),
  1×) — **Ausgang:** *entfallen*, wenn Lastenheft und Folge-ADR bei höchstens acht Kriterien und den
  vier Entscheiden der Welle bleiben; sonst *eingetreten* mit einer Evidenz-Datei im Register (dann
  2×).
- Die enge Antwort zu Frage 3 (höchstens drei Leerzeichen) lässt eine Tabelle unter einer
  Vier-Spalten-Marke oder in einem verschachtelten Listenpunkt als Befund stehen; ob der Bestand des
  Melders solche Zeilen trägt, ist von hier nicht gemessen — **Ausgang:** *weiter offen*, wenn der
  Auftraggeber die enge Fassung bestätigt (die Anforderung nennt die Grenze mit ihrer Bedingung; der
  Register-Eintrag entsteht bei der Closure); *entfallen*, wenn er stattdessen die Listen-Erkennung
  verlangt — das ist dann eine Plan-Änderung vor dem Code (§4, Rückführung), nicht Teil dieses
  Slice.
- Lastenheft und Spezifikation führen dieselbe Aussage je in eigenen Worten; ein Satz, der nur an einer
  Stelle umgeschrieben wird, stünde gegen den anderen — **Ausgang:** *entfallen*, wenn die zweite
  Suche aus §3 nach dem Editieren keine der dort genannten Stellen mehr mit der alten Aussage zeigt
  (jede umgeschrieben oder als unberührt begründet).

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

*(Wird bei der Closure gefüllt.)*

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
(offen, 2× — jede Randbedingung bekommt ein Kriterium, das [`slice-294`](../open/slice-294-linelength-eingerueckte-tabellen-kern.md) mit einem Test belegt),
[`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md)
(offen, 2× — „dieselbe Lexik“ stimmt nur für eine Familie; Frage 4 der Welle),
[`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md)
(verkörpert in welle-74, 6× — der Kopplungs-Test, Kriterium 5),
[`BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft`](../observations/BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft/observation.md)
(offen, 1× — das Review liest die geänderte Anforderung gegen Lastenheft §5, dessen Ausnahme „mit
derselben Markdown-Lexik wie die übrigen Module“ sagt) und
[`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/observation.md)
(verkörpert, 4× — die Spiegel-Suche in §3 nimmt Handbuch und Konfigurations-Vorlage auf; sie gehören
[`slice-294`](../open/slice-294-linelength-eingerueckte-tabellen-kern.md)). **Kein Treffer:**
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
