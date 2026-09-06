# `make <target>` — <ein Satz: was das Target tut>

> **Template-Hinweis.** Vorlage für **eine** Sensor-Datei. Kopiere nach
> `harness/sensors/<target>.md`, ersetze alle `<Platzhalter>` und lösche diesen
> Block **und den Abschnitt „Regeln dieser Datei" darunter** — beide gehören zur
> Anleitung, nicht zur Zielform. Diese Datei entsteht **nur, wenn ein Target mehr
> braucht als einen Satz** — Deckungsgrenze, Ausgabe-Bedeutung, Exit-Codes,
> Abbruch-Bedingungen. Passt der Vertrag in einen Satz, ist die Zeile in
> [`../README.md` §Sensors](../README.md#sensors-feedback-gates) vollständig und
> diese Datei überflüssig; ob der Überhang schon unter der Tabelle steht oder in
> die Zelle gedrängt wurde, ist dieselbe Sache (Baseline-Regelwerk
> `grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt).

Regeln dieser Datei:

- **Die Index-Zeile bleibt.** Diese Datei ersetzt die Tabellenzeile nicht, sie
  vertieft sie; die Target-Zelle verlinkt hierher — und dieser Link ist das
  einzige, was ein Sensor an der Zuordnung prüfen kann. Von außen — aus
  `AGENTS.md`, einer ADR, einem Slice — wird **diese Datei direkt** adressiert,
  nicht der Index: sie wandert nie.
- **Kein Status- und kein Datumsfeld.** Der Zustand des Gates ist sein Lauf, und
  der lebt in CI; das Änderungsdatum hält `git`.
- **Kein Lifecycle-Verzeichnis.** Es gibt kein `sensors/done/`: Wird das Target
  retiriert, verschwinden Zeile und Datei. Ein lebendes Artefakt, das dann ins
  Leere zeigt, behauptet eine Deckung, die es nicht mehr gibt — der rote
  Link-Sensor ist die richtige Antwort darauf.
- **Was hier NICHT steht:** womit das Werkzeug selbst gedeckt ist — welcher Test
  welche Hälfte trägt, welcher Mutations-Fall welchen Zweig bewacht. Das ist die
  Frage „ist das Werkzeug richtig?", und ihre Antwort lebt bei ihm: in seiner
  ADR, seiner Spec-Zeile, seinem Skriptkopf. Diese Datei sagt, wie ein Lauf zu
  **lesen** ist.

---

## Vertrag

<Was wäre verletzt, wenn dieses Target rot wird — derselbe Satz wie in der
Index-Zelle, hier ausgeschrieben. Für ein Nicht-Gate stattdessen: was es tut
(bewegt · misst · sagt) und warum es kein Gate ist — es urteilt nicht über den
Zustand des Repos, sondern über die Vorbedingungen seines eigenen Laufs.>

## Grenze — was das Grün nicht abdeckt

<Der Prüfbereich, und wo er enger ist als der Bereich, über den das Grün
gelesen wird. Je Loch ein Punkt, und dazu, ob es heilbar ist oder permanent.>

1. **<Loch 1>** — <warum es außerhalb liegt; heilbar durch <…> / permanent>
2. **<Loch 2>** — <…>

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`<kommando, das den Prüfbereich zählt>`. Eine eingefrorene Zahl stünde hier
falsch, sobald jemand committet.

<Trägt der Lauf eine Vollständigkeits-Zeile („N Dateien geprüft, 0 Befunde"),
gehört hierher der Satz, worüber sie eine Aussage macht — über den Ausschnitt,
nicht über das Repo.>

## Ausgabe und Ausgänge

<Nur, wenn der Lauf mehr als grün/rot sagt.>

| Exit | Bedeutung |
|---|---|
| 0 | <…> |
| <n> | <…> |

<Unterscheidet der Lauf Fehlschlag-Ursachen — geprüfter Baum vs. Umgebung/Leitung
—, steht hier, woran er sie auseinanderhält und in welcher Zeile er das sagt.>

## Sperren

<Woran der Lauf abbricht, bevor er etwas tut. Je Sperre ihr Name, wie ihn die
Abbruch-Meldung nennt, und was zu tun ist.>

- `<sperre>` — <Bedingung> → <was den Weg frei macht>

## Bindung

<`ADR-<NNNN>` · `CO-<NNN>` · `LH-<…>` · Schwelle · Image-Hash — dieselbe
Angabe wie in der Index-Zelle. Für ein Nicht-Gate: `kein Gate` plus die ID,
die seine Existenz begründet.>
