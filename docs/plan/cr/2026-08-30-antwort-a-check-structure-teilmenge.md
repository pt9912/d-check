# Antwort auf den `a-check`-CR vom 2026-08-30 — `structure`: die geprüfte Menge deklarierbar machen

**Absender der Antwort:** d-check
**Datum:** 2026-08-30
**Bezug:** [eingehender CR vom 2026-08-30](2026-08-30-cr-a-check-structure-teilmenge.md),
zwei Optionen gegen [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
**Ergebnis:** **beide angenommen**, einer davon unter anderem Namen; zwei
gemessene Korrekturen am Antragstext
**Einordnung:** MINOR — Lastenheft `0.76.0`, Begründung in
[ADR-0075](../adr/0075-erklaerte-teilmenge-in-structure.md)

---

## Vorab

Der Antrag trägt, und er trägt aus dem Grund, den er selbst nennt: die
Slice-Vorlage der Baseline liefert eine DoD mit konstanten Punkten aus, und
wer sie benutzt **und** die Größen-Regel prüft, bekommt Falsch-Positive auf
jedem neuen Slice. Das ist hier nachgemessen und nicht abgeschrieben:
`max-tasks: 3` über die **89** DoD-Abschnitte dieses Repos (in 175
Slice-Plänen) liefert **80** `section-oversized` bei **444** Task-Items.

Die Verortung ist die beantragte — Optionen an bestehenden Bedingungen, kein
neues Modul, kein neuer Grund-Code, Default byte-identisch. Letzteres ist
gemessen: derselbe `max-tasks`-Lauf gegen das Image **vor** der Änderung und
danach liefert **166** Befunde, `diff` leer.

Zwei Dinge sind anders als beantragt, und beide kommen aus Messungen an eurem
eigenen Beispiel-Muster. Sie stehen hier vollständig, weil ihr sie in eurer
Konfiguration braucht.

## Korrektur 1 — euer Beispiel-Muster ist auf diesem Korpus wirkungslos

`tasks-ignore-pattern: '(make gates|Closure-Notiz|Beobachtungs-Register|Risiko aus)'`,
je Alternative einzeln gegen die 444 Items gefahren:

| Alternative | ignorierte Items |
|---|---|
| `make gates` | **0** |
| `Closure-Notiz` | 12 |
| `Beobachtungs-Register` | 1 |
| `Risiko aus` | 0 |

**Der Grund ist die Zusage, die ihr selbst verlangt — für das Item-Muster.**
Ihr schreibt: *„Beide Muster dürfen Code-Blöcke und Inline-Code nicht sehen."*
Für `tasks-ignore-pattern` ist es genau so umgesetzt: die Zählung liest den
**bereinigten** Abschnitts-Text, in dem Inline-Code durch Leerzeichen ersetzt
ist. In **86 der 89** DoD-Abschnitte steht die Wendung als `` `make gates` ``
(424 Vorkommen, **null** davon ohne Backticks), und ein Muster darauf trifft
folglich Leerzeichen.

**Für das Abschnitts-Muster konnten wir euch diese Zusage nicht geben, und
warum, steht weiter unten:** `exempt-section-pattern` teilt seine Zeichenkette
mit `section-pattern`, und das sieht Inline-Code. Fence-treu ist es;
inline-code-treu nicht. Wir haben die halbe Zusage lieber benannt als
stillschweigend behauptet.

Das ist keine Panne der Umsetzung, sondern die Kehrseite der Fence-Treue: **was
fence-treu gezählt wird, ist auch fence-treu ignorierbar.** Es hat uns beim
Bauen genauso erwischt, deshalb steht es jetzt an drei Stellen — im Lastenheft,
im Handbuch und in der Meldung selbst.

**Die Diagnose ist eingebaut:** ist ein Muster gesetzt, nennt die Meldung die
Zahl der ignorierten Items — **auch bei null**.

```
docs/slice-x.md:3	docs/slice-*.md :: ## 4. Definition of Done	section-oversized	Abschnitt trägt 4 Task-Items (3 ignoriert), erlaubt sind 3
```

`(0 ignoriert)` heißt: das Muster wirkt nicht. **Benannte Grenze:** das greift
nur, **solange die Regel meldet** — wer so breit ignoriert, dass die Schwelle
nie fällt, sieht nichts.

## Korrektur 2 — verankern reicht nicht; das Muster muss auch auf den richtigen Text zielen

An denselben 444 Items, drei Musterformen:

| Musterform | ignorierte Items | davon falsch |
|---|---|---|
| euer Ausdruck, frei | 13 | **2** |
| **derselbe** Ausdruck, verankert (`^…`) | **1** | **1** |
| verankert **und** auf Text, der die Bereinigung überlebt | 26 | 0 |

Die zwei falschen der ersten Zeile sind **Liefer**-Punkte, die die
Closure-Notiz nur als ihren **Ort** nennen — *„§2/§3/§4/§6 tragen `SPEC-NNN`
fortlaufend; Zählung in der Closure-Notiz"* und
*„[ADR-0012](../adr/0012-kern-paketschnitt-model-rules-app.md)-§Kern-Messung in
der Closure-Notiz"*. Ein freies Substring-Muster nimmt sie aus der Zählung,
ohne dass jemand es sieht.

**Die mittlere Zeile ist die unangenehme, und sie gehört euch gesagt:**
verankert man **euren** Ausdruck, bleibt genau **ein** Treffer übrig — und der
ist ebenfalls ein Liefer-Punkt (*„Beobachtungs-Register trägt die Klasse als
`BEO-011`, Zähler 3 …"*). Auf diesem Korpus ist euer Ausdruck also in **beiden**
Formen unbrauchbar, und der Grund ist Korrektur 1: `make gates` und
`Closure-Notiz` stehen fast überall in Backticks oder mitten im Satz. Die
dritte Zeile misst ein Muster, das wir für **diesen** Korpus gewählt haben
(`^grün \(Exit explizit\)` — der Text, der nach dem Leeren der Inline-Code-Spanne
am Item-Anfang übrig bleibt). Es ist kein Vorschlag für euren Bestand, sondern
der Beleg, dass die Kombination *verankert + überlebender Text* trägt.

**Erste Fassung dieser Antwort war hier falsch** und stellte „frei 13" neben
„verankert 26", als wäre es derselbe Ausdruck — ein `^` kann die Treffermenge
nur verkleinern. Der Fehler ist von einem unabhängigen Review und einer
Verifikation unabhängig voneinander gefunden worden; die Tabelle oben ist die
nachgemessene.

**Damit `^` benutzbar ist, sieht das Muster den Item-Text, nicht die Zeile.**
Gegen die rohe Zeile bezeichnete `^` immer den Listen-Marker, und die
verankerte Form wäre unschreibbar. Verglichen wird also der **getrimmte Text
hinter Listen-Marker und Checkbox**:

```yaml
tasks-ignore-pattern: '^\*\*Konstante:'    # trifft
tasks-ignore-pattern: '^- \[ \] \*\*Konstante:'   # trifft NICHT
```

## Der zweite Schlüssel heißt `exempt-section-pattern`

Nicht `exempt-sections`. Zwei Gründe, beide aus dem Bestand:

1. **In `structure` bedeutet das Suffix `-pattern` RE2** —
   `section-pattern`, `forbid-pattern`, `require-pattern`. Euer eigener erster
   Schlüssel, `tasks-ignore-pattern`, folgt dieser Konvention; der zweite tat
   es nicht.
2. **`exclude-sections` ist in zwei anderen Modulen vergeben** (`vcs`,
   `sources`), dort als **Liste literaler** Überschriften. Zwei ähnliche Namen
   für zwei Formen sind eine Verwechslung, die niemand braucht.

**Und das Muster sieht dieselbe Zeichenkette wie `section-pattern`:** die
getrimmte Überschriften-Zeile **einschließlich** der `#`-Folge. Das ist eine
Abweichung von eurem Antrag, und sie geht auf **unsere** Kappe, nicht auf eure:
ihr habt euren Vergleichsgegenstand ausdrücklich deklariert (*„Abschnitte,
deren **Überschriftstext** dieses RE2 trifft"*), und euer Beispiel ist gegen
diese Semantik korrekt geschrieben. Wir haben uns dagegen entschieden, weil
zwei RE2 in **einer** Regel dann zwei verschiedene Zeichenketten läsen — ein
Muster, analog zum Nachbarn `section-pattern` geschrieben, griffe still
daneben. Der Preis ist, dass euer Beispiel-Muster so wie es dasteht **nichts**
trifft. Die Form für euch lautet:

```yaml
exempt-section-pattern: '^### AC-[A-Z]+-0(0[1-9]|1[0-9])\b'
```

**Drei Zusagen dazu, die euer Antrag nicht erfragt hat, aber braucht.**
Erstens: leert das Muster die Abschnitts-Menge, meldet die Regel
`section-missing` und nennt den Schlüssel — sie wird **nicht still grün**
(dieselbe Nullmengen-Härte wie bei `exempt-paths`, eine Stufe tiefer), und
dieser Befund behält seinen Text auch neben einem `hint`. Zweitens: der Abzug
läuft **vor** der Kardinalitäts-Prüfung — zwei Treffer minus einer
ausgenommenen sind bei `sections: one` in Ordnung, nicht `section-ambiguous`.
Drittens: ein **gesetztes** Muster gehört zur **Regel-Identität**. Ohne das
wäre die Paarung, die Grandfathering eigentlich braucht — Bedingung A für alle
Abschnitte, Bedingung B nur für die übrigen —, ein Konfigurations-Duplikat mit
Exit 2 gewesen; genau die Form also, für die ihr den Schlüssel beantragt habt.

## Angenommen, unverändert

- **Beides als Optionen an bestehenden Bedingungen** — eure Begründung („keine
  neue Frage, nur eine erklärte Grundmenge") ist die, die auch hier trägt.
- **Kein neuer Grund-Code, Default byte-identisch.**
- **Fence-Treue für beide Muster** — geerbt, nicht nachgebaut, und mit je einer
  Probe belegt. **Inline-Code-Treue nur für das Item-Muster**, siehe
  Korrektur 1; die halbe Zusage ist benannt, nicht stillschweigend gegeben.
- **Eure Abgrenzung** — Grandfathering „ab Nummer N" als Werkzeug-Begriff und
  abschnitts-übergreifende Bedingungen bleiben draußen. Beides sehen wir
  genauso.

## Eure Paritäts-Tabelle, gefahren

Alle fünf Zeilen sind als Tests umgesetzt
(`internal/hexagon/core/rules/structure_teilmenge_test.go`) und laufen in
`make test`. Die drei tragenden fahren ihre **Umkehr in derselben Funktion**:
ohne den Schlüssel ist der Vorzustand mitgeprüft — ein Regressions-Test ohne
belegte Regression ist keiner.

Ergänzt haben wir, was eure Tabelle nicht deckt: die **Überdeckung** (Meldung
nennt die ignorierten, auch bei null), die **leere Menge**, die Reihenfolge
gegen `sections: one`, das **Muster-Ziel** in beiden Richtungen, die
**Inline-Code**-Blindheit aus Korrektur 1, das Zusammenspiel mit `hint` in
beiden Richtungen und die **Regel-Identität**.

**Eine Zusage war zuerst unbelegt, und ein Review hat es gefunden:** unsere
Fence-Probe für das Abschnitts-Muster richtete ihr Ausnahme-Muster auf die
**gefencte** Überschrift — damit war „nie gewählt" von „gewählt und dann
ausgenommen" am Ergebnis nicht unterscheidbar, und ein fence-blinder Scanner
wäre grün durchgelaufen. Sie zeigt jetzt auf die **echte** Überschrift; gegen
eine Fence-Blindheits-Mutation wird sie rot, gemessen.

## Was ihr beim Umstellen ändern müsst

1. `exempt-sections` → `exempt-section-pattern`.
2. Das Abschnitts-Muster um die `#`-Folge ergänzen.
3. Das Item-Muster **verankern** — und zuerst prüfen, ob eure Konstanten in
   Inline-Code stehen. Auf unserem Korpus wäre euer Ausdruck in **beiden**
   Formen unbrauchbar; wählt den Text, der nach dem Leeren der
   Inline-Code-Spannen am Item-Anfang übrig bleibt. Steht `(0 ignoriert)` in
   der Meldung, greift das Muster nicht.
4. Wenn ihr die Überdeckung sehen wollt, **keinen `hint`** an dieselbe Regel
   hängen — er ersetzt die Meldung samt der Zahl.

## Was wir nicht tun

**Wir schalten die Fähigkeit nicht über unseren eigenen Bestand scharf.** Die
80 `section-oversized` sind der Anlass, nicht der Auftrag; ob `max-tasks` über unsere
Slice-Pläne läuft, ist ein eigener Entscheid nach dieser Fähigkeit.

Und eine ehrliche Zeile zum Schluss: **die geprüfte Menge zu verkleinern ist
per Konstruktion eine Lockerung.** Sie ist opt-in, aber jede Konfiguration, die
sie setzt, senkt ihre eigene Zusage, und kein Gate meldet das. Die
Sichtbarkeits-Zusage ist das Beste, was das Werkzeug dagegen tun kann — das
Urteil, ob eine Teilmenge **berechtigt** ist, bleibt bei euch.
