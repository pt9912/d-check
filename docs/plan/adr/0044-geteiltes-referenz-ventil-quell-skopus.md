# ADR-0044 — Geteiltes Referenz-Ventil `ignore-refs`: Quell-Skopus, Zwei-Feld-Semantik, Alias

**Status:** Accepted
**Datum:** 2026-07-18
**Autor:** pt9912
**Schärft:** [`DC-FA-REF-001.a`](../../../spec/spezifikation.md#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs)
**Bezug:** [`DC-FA-REF-001`](../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus), [`DC-FA-CODE-001`](../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in), [`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links), [`DC-FA-ANCH-001`](../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors), [`DC-FA-SCAN-001`](../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln), [ADR-0030](0030-tracked-referenz-ziele.md), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).

## Kontext

Ein Konsument (`ai-harness-course`, CR 2026-07-17) pflegt ein Verzeichnis mit
**Template-Dateien**. Deren Referenzen zerfallen in zwei Klassen: **Ziel-Repo-Pfade**
(Platzhalter, die im Quell-Repo bewusst nicht auflösen) und **Kurs-/Doku-Verweise**
(die auflösen und beim Release-Bau unveränderlich auf eine getaggte Blob-URL gepinnt
werden). Weil die erste Klasse massenhaft `target-missing` erzeugt, steht das **ganze
Verzeichnis** in `scan.ignore` — und damit ist auch die zweite Klasse **ungeprüft**.
Genau deren Auflösung wird beim Release eingefroren: wer eine Überschrift umbenennt
und das Template vergisst, liefert einen toten Anker aus, und kein Gate merkt es.
Gemessen (Verzeichnis aus `scan.ignore` entfernt, digest-gepinntes Image): 42 Findings,
davon 8 echte Platzhalter und 34 im Ziel-Repo real — 37 der 42 sind `links`, 5
`codepaths`.

Die vorhandenen Ventile greifen nicht:

- **`scan.ignore`** (Quell-Achse) ist alles-oder-nichts pro Datei — es opfert die
  prüfbare Klasse, um die symbolische loszuwerden.
- **Der Zeilen-Marker `d-check:ignore`** ist hier **aktiv schädlich**: Templates
  werden vom Adopter kopiert, die Marker reisen mit und unterdrücken an Zielposition
  dauerhaft eine dann echte Prüfung.
- **`codepaths.ignore-refs`** ([`DC-FA-CODE-001`](../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
  Ziel-Achse) ist **modul-lokal** (erreicht die 37 `links`-Findings gar nicht) und
  **global** (kein Quell-Skopus — dieselben Pfade würden repo-weit blind).

Es fehlt „Referenz auf **Y** nicht prüfen, **wenn sie in X steht**" — das
**Kreuzprodukt** der zwei vorhandenen Achsen (Quelle × Ziel), nicht eine dritte.

## Entscheidung

1. **Geteiltes Ventil statt Modul-Duplikat.** `ignore-refs` wandert von
   `codepaths`-lokal zu einer geteilten Top-Level-Anforderung
   ([`DC-FA-REF-001`](../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)),
   die `links`/`anchors`/`codepaths` honorieren — das **Ziel-Achsen-Pendant** zu
   `scan.ignore` ([`DC-FA-SCAN-001`](../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln),
   Quell-Achse). Das folgt dem mit dem Auftraggeber festgelegten **Kürzel-Kriterium**:
   *querschnittlich (mehrere Module teilen die Fähigkeit) → neues Bereichskürzel;
   Einzelmodul-Erweiterung → bestehende Anforderung ändern.* `ignore-refs` ist
   querschnittlich (drei Module) → geteilt.

2. **Zwei Felder (`refs`/`keep`) statt `!`-Negation.** Der CR schlug ursprünglich
   einen Schlüssel `ignore-refs-in` mit `!`-Negation und **gitignore-Last-Match** vor.
   Gewählt: `refs` (Ziel-Globs) und `keep` (Ausnahmen), Semantik **ignorieren, wenn
   `refs` matcht ∧ `keep` nicht** — `keep` gewinnt **unbedingt und
   reihenfolge-unabhängig**. Grund: die entscheidende Zeile der Konsumenten-Messung
   war bereits `match(refs) ∧ ¬match(keep)`; der CR-Text hatte dieselbe Semantik nur
   in `!`-Syntax gegossen und dann fälschlich gitignore-Ordnung drangeschrieben — er
   war an der Stelle **schlicht falsch beschriftet**. Der Preis der Zwei-Feld-Form
   (sie kann nicht `ignore`→`keep`→`ignore` alternieren) ist real, aber **gemessen
   null**: alle 24 von `keep` zurückgeholten Ziele existieren, **kein** Fall braucht
   ein Re-Ignore. Reihenfolge-Unabhängigkeit entspricht zudem dem Lesen einer
   YAML-Liste.

3. **`keep` ist konstitutiv, nicht optional.** Gegen den realen Bestand: nur `refs`
   ignoriert 62 Verweise, davon **24 real** (fälschlich blind); `refs ∧ ¬keep`
   ignoriert 38, davon **0 real** — 63 echt geprüft. Ohne `keep` tauscht das Feature
   nur eine Blindstelle gegen eine andere; mit `keep` ist der Schnitt exakt.

4. **`in:` ist Skopus, keine vierte Achse.** Ein danebengestellter `ignore-refs-in`
   hätte die Ziel-Achse **dupliziert**; ein `in:`-Glob auf die **Quelldatei** an
   `ignore-refs` ist dieselbe Achse mit Skopus — das Kreuzprodukt der zwei
   vorhandenen.

5. **Ziel-Globs matchen den aufgelösten Pfad** — keine neue Auflösungs-Semantik,
   dieselbe Ventil-Parität, die [ADR-0030](0030-tracked-referenz-ziele.md) schon
   entschieden hat (Befund-`target` = aufgelöster Pfad).

6. **Alias `codepaths.ignore-refs` bleibt (kein Config-Bruch).** Die modul-lokale
   Liste wirkt weiter wie ein `ignore-refs`-Eintrag ohne `in`/`keep`, skopiert auf
   `codepaths` — byte-identisch. Ob der Alias eine **Deprecation-Frist** bekommt,
   ist ein späterer CR, nicht diese Entscheidung; bis dahin verdoppelt er die
   Config-Oberfläche.

7. **Wirkung: Existenz/Escape/Anker unterdrückt, Symlink bleibt.** Das Ventil
   überspringt für ein Treffer-Ziel die Existenz-, die Repo-Escape- und (bei
   Markdown) die Anker-Prüfung — dieselbe Unterdrückung wie das bisherige
   `codepaths.ignore-refs`. Die Symlink-Ablehnung ist eine ziel-unabhängige
   Sicherheits-Prüfung an **existierenden** Zielen (Defense-in-Depth) und bleibt.

8. **Muster statt Heuristik.** Ein verfälschter Pfad in einer Template-Datei muss
   ein **ERROR** bleiben. Eine „ignoriere, was nicht auflöst"-Heuristik bestünde
   diesen Tippfehler-Test nicht — deshalb explizite `refs`-Muster, additiv und
   fail-visible.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Geteiltes `ignore-refs` + `in`/`keep`** (gewählt) | Kreuzprodukt der Achsen; ein Ventil für drei Module; `keep` schneidet exakt; Alias bricht nichts | Config-Fläche wächst; Alias verdoppelt sie bis zu einer Deprecation |
| `ignore-refs-in` mit `!`-Negation (CR-Original) | nah am Eingereichten | zweiter Glob-Dialekt; gitignore-Ordnung, die gemessen nie gebraucht wird; dupliziert die Ziel-Achse |
| Ganzes Verzeichnis in `scan.ignore` (Status quo) | null Aufwand | opfert 63 real prüfbare Verweise; Blindheit wird beim Release eingefroren |
| Zeilen-Marker je Referenz | vorhanden | reist mit kopierten Templates; unterdrückt an Zielposition eine echte Prüfung |
| Nur `refs` (ohne `keep`) | einfacher | 24 real auflösende Verweise fälschlich blind — tauscht Blindstelle gegen Blindstelle |

**Fitness-Funktion:**

- Der Realdatenbeleg gegen das Konsumenten-Repo meldet **0 Findings bei 63
  tatsächlich geprüften** Verweisen — nicht durch Wegschauen (38 ignoriert, davon 0
  real existierend).
- Ein verfälschter Pfad **und** ein verfälschter Anker in einer Template-Datei
  erzeugen je einen Befund (Tippfehler-Test).
- Dieselben Ziel-Muster außerhalb des `in`-Globs bleiben voll geprüft
  (Skopus-Isolation).
- Ohne Block (weder Top-Level `ignore-refs` noch `codepaths.ignore-refs`) jede
  Prüfung byte-identisch
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)), read-only
  ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

## Konsequenzen

- **Positiv:** ein bislang blinder, prüfbarer Referenz-Typ wird sichtbar; das
  Template-Verzeichnis des Konsumenten bekommt 63 echte Prüfungen statt eines
  pauschalen `scan.ignore`; opt-in und byte-identisch für Nicht-Konsumenten; ein
  Ventil deckt drei Module.
- **Negativ / Kosten:** die Config-Fläche wächst um eine vierte Achse
  (`scan.ignore` · Zeilen-Marker · `exempt-paths` · `ignore-refs`+`in`/`keep`); die
  Handbuch-Doku muss die Achsen **gegeneinander** erklären, nicht nur nebeneinander.
  Der Alias verdoppelt die Ziel-Achsen-Fläche, bis ein späterer CR über eine
  Deprecation entscheidet. `keep` ist eine kuratierte Kante mit eigener Drift-Gefahr
  (wie `matrix.exclude-sections`).
- **Verworfen:** `!`-Negation, danebengestellter `ignore-refs-in`, Nur-`refs`,
  Zeilen-Marker, Status-quo-`scan.ignore` (jeweils oben begründet).

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-07-18 | Proposed. Change Request Konsument `ai-harness-course` (2026-07-17), Design nach zwei Rückfragen verfeinert (zwei Felder statt `!`-Negation; `in:`-Skopus statt vierter Achse). §4-Verortungs-Vorfrage vom Auftraggeber entschieden (neues geteiltes Kürzel), gemeinsames Kürzel-Kriterium mit slice-079. Umsetzender Slice slice-078. |
| 2026-07-18 | **Accepted.** Code über `links`/`anchors`/`codepaths` + Alias umgesetzt (Semantik deckungsgleich mit dieser ADR), keep/in/Wiring per Mutation gepinnt; Realdatenbeleg gegen `ai-harness-course` erbracht (Baseline 42 → Ventil 0 Befunde, zwei injizierte Tippfehler beider Klassen gefangen — nicht durch Wegschauen); unabhängiger Review R1 ACCEPT-WITH-NITS (Nits eingearbeitet). Als v0.49.0 veröffentlicht (slice-078). |
