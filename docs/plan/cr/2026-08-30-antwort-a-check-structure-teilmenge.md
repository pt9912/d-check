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
`max-tasks: 3` über die **86** Slice-Pläne dieses Repos liefert **80**
`section-oversized` bei **444** Task-Items.

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

**Der Grund ist die Zusage, die ihr selbst verlangt.** Ihr schreibt: *„Beide
Muster dürfen Code-Blöcke und Inline-Code nicht sehen."* Genau so ist es
umgesetzt — die Zählung liest den **bereinigten** Abschnitts-Text, in dem
Inline-Code durch Leerzeichen ersetzt ist. In **allen** 86 Abschnitten steht
die Wendung als `` `make gates` `` (null Fundstellen ohne Backticks), und ein
Muster darauf trifft folglich Leerzeichen.

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

## Korrektur 2 — verankert das Muster, und zwar am Item-Text

An denselben 444 Items:

| Musterform | ignorierte Items | davon falsch |
|---|---|---|
| frei, wie im CR | 13 | **2** |
| verankert (`^…`) | 26 | 0 |

Die zwei falschen sind **Liefer**-Punkte, die die Closure-Notiz nur als ihren
**Ort** nennen — *„§2/§3/§4/§6 tragen `SPEC-NNN` fortlaufend; Zählung in der
Closure-Notiz"* und *„[ADR-0012](../adr/0012-kern-paketschnitt-model-rules-app.md)-§Kern-Messung in der Closure-Notiz"*. Ein freies
Substring-Muster nimmt sie aus der Zählung, ohne dass jemand es sieht.

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
getrimmte Überschriften-Zeile **einschließlich** der `#`-Folge. Euer
Beispiel-Muster `'^AC-[A-Z]+-0(0[1-9]|1[0-9])\b'` trifft damit **nichts** — im
Lastenheft steht `### AC-…`. Das war der dritte gemessene Punkt: zwei RE2 in
**einer** Regel mit zwei verschiedenen Zielen sind eine Falle, und ein Muster,
analog zum Nachbarn geschrieben, hätte still danebengegriffen. Die Form für
euch lautet:

```yaml
exempt-section-pattern: '^### AC-[A-Z]+-0(0[1-9]|1[0-9])\b'
```

**Zwei Zusagen dazu, die euer Antrag nicht erfragt hat, aber braucht.**
Erstens: leert das Muster die Abschnitts-Menge, meldet die Regel
`section-missing` und nennt den Schlüssel — sie wird **nicht still grün**
(dieselbe Nullmengen-Härte wie bei `exempt-paths`, eine Stufe tiefer).
Zweitens: der Abzug läuft **vor** der Kardinalitäts-Prüfung — zwei Treffer
minus einer ausgenommenen sind bei `sections: one` in Ordnung, nicht
`section-ambiguous`.

## Angenommen, unverändert

- **Beides als Optionen an bestehenden Bedingungen** — eure Begründung („keine
  neue Frage, nur eine erklärte Grundmenge") ist die, die auch hier trägt.
- **Kein neuer Grund-Code, Default byte-identisch.**
- **Fence-Treue für beide Muster** — geerbt, nicht nachgebaut, und mit je einer
  Probe belegt.
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
gegen `sections: one`, das **Muster-Ziel** in beiden Richtungen und die
**Inline-Code**-Blindheit aus Korrektur 1.

## Was ihr beim Umstellen ändern müsst

1. `exempt-sections` → `exempt-section-pattern`.
2. Das Abschnitts-Muster um die `#`-Folge ergänzen.
3. Das Item-Muster **verankern** und prüfen, ob eure Konstanten in Inline-Code
   stehen — steht `(0 ignoriert)` in der Meldung, greift es nicht.

## Was wir nicht tun

**Wir schalten die Fähigkeit nicht über unseren eigenen Bestand scharf.** Die
80 Befunde sind der Anlass, nicht der Auftrag; ob `max-tasks` über unsere
Slice-Pläne läuft, ist ein eigener Entscheid nach dieser Fähigkeit.

Und eine ehrliche Zeile zum Schluss: **die geprüfte Menge zu verkleinern ist
per Konstruktion eine Lockerung.** Sie ist opt-in, aber jede Konfiguration, die
sie setzt, senkt ihre eigene Zusage, und kein Gate meldet das. Die
Sichtbarkeits-Zusage ist das Beste, was das Werkzeug dagegen tun kann — das
Urteil, ob eine Teilmenge **berechtigt** ist, bleibt bei euch.
