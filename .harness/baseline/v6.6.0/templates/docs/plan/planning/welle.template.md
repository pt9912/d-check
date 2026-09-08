# Welle <welle-id>: <Titel>

> **Template-Hinweis.** Vorlage für eine Welle (Bündel von Slices, das
> gemeinsam geplant und abgeschlossen wird, siehe
> [Baseline-Regelwerk Modul 5](../../../../regelwerk/modul-05-planning-harness.md)
> und [Modul 6](../../../../regelwerk/modul-06-roadmap.md)).
> Kopiere nach `docs/plan/planning/<welle-id>.md` und ersetze
> Platzhalter. Lösche diesen Block.

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<NN>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** M<NN> oder "kein Meilenstein-Bezug".

**Verantwortlich:** <Name>. **Datum:** YYYY-MM-DD.

---

## 1. Welle-Ziel

<!-- BEDIENHINWEIS: Eine Aussage, die sich an einem Lasttest oder
Akzeptanzkriterium spiegelt. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

<…>

## 2. Trigger (Welle startet)

<!-- BEDIENHINWEIS: Was muss vorher passiert sein? -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- <z.B. Welle <welle-vorher-id> done.>
- <z.B. ADR-<NNNN> accepted.>

## 3. Closure-Trigger (Welle schließt)

<!-- BEDIENHINWEIS: Aktion, nicht Termin. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- <z.B. Alle Slices done.>
- <z.B. `make fullbuild` grün.>
- <z.B. Replay-Lauf gegen Golden Set durchläuft.>
- <z.B. Closure-Notiz in `welle-<NN>-results.md`.>

## 4. Slices in dieser Welle

<!-- BEDIENHINWEIS: keine Status-Spalte ergaenzen. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-<NN-A> | <…> | LH-FA-<NN> |
| slice-<NN-B> | <…> | LH-FA-<NN> |

## 5. Abhängigkeiten

<!-- BEDIENHINWEIS: Falls jemand diese Welle aendert — was bricht? -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle <welle-id> (wegen <Vertragspunkt>).
- Wird blockiert von: Welle <welle-id>.

## 6. Out-of-Scope für diese Welle

<!-- BEDIENHINWEIS: explizite Nicht-Inhalte, schuetzt vor Scope-Creep. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- <…>

## 7. Closure-Notiz

<!--
BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md §Verwendung,
Schritt 5) und darf deshalb nichts Tragendes halten.

- Erst nach Welle-Abschluss fuellen; nur die Nummer, nicht die volle Welle-ID.
- Ziel-Form der Ergebnis-Notiz: `welle-results.template.md` — Schwester-Vorlage
  im Template-Verzeichnis, kein Artefakt deines Repos. Sie ist von der
  Ruheort-Regel ausgenommen und faellt mit diesem Kommentar ohnehin weg.
-->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: <Zeiger auf `welle-<NN>-results.md`, Geschwister im Ruheort `done/`>
Zähler: <Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort>
