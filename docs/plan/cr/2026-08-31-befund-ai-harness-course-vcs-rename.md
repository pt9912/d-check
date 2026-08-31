# Eingehender Befund aus `ai-harness-course` — `vcs` meldet den reinen Rename über `--range` nicht

**Absender:** ai-harness-course (Adopter von d-check).
**Eingegangen:** 2026-08-31, über den Auftraggeber.
**Klasse:** **Werkzeug-Befund** gegen eine zugesagte Anforderung — kein Change
Request, keine neue Fähigkeit beantragt, keine Vertrags-Änderung erbeten.
**Gegenstand:** [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in).
**Zusammenhang:** hängt an keinem offenen Faden; er zieht den Abschnitt
*Bereits gelöst: Immutabilität* des
[BEO-CR](2026-08-31-cr-ai-harness-course-observations-relational.md) zurück.

Dieses Dokument hält den Befund **wie empfangen** fest.

---

## Die Zusage

Die Anforderung nennt drei Befund-Auslöser, darunter wörtlich *gelöschte oder
umbenannte immutable Datei → `core-drift-vcs`*, und zwei **gleichrangige**
Eingabe-Modi: `--range <base>..<head>` für die CI und `--staged` für den
lokalen pre-commit-Hook. Ein Modus-Vorbehalt steht nirgends.

## Die Messung des Absenders

| Lauf | Ergebnis |
|---|---|
| `git mv` gestaged, `--staged` | `core-drift-vcs` auf dem alten Pfad |
| derselbe Rename committet, `--range base..head` | **0 Befunde** |
| Rename **mit** vollständiger Umformulierung, `--range` | `core-drift-vcs` |
| Datei stattdessen gelöscht, `--range` | `core-drift-vcs` |

Gemessen am Image-Pin `v0.67.0`; als wiederholbarer Lauf in seinem Repo unter
`lab/team-sim` (s15a/s15b/s15c).

## Warum das mehr ist als eine Lücke (Begründung des Absenders)

- **Der stille Modus ist der CI-Pfad.** Laut wird nur der lokale Hook — der,
  den man überspringen kann.
- **Die Aussage der Anforderung ist modus-abhängig**, und das steht nirgends.
- Der Fall ist der **normale** Weg, eine Pfadidentität zu verletzen: eine Datei
  wird verschoben, ohne dass jemand ihren Inhalt anfasst.
- **Wo der Dateiname eine Aussage trägt, wird aus der Lücke eine Fälschung.** In
  seinem Register-Entwurf ist der Dateiname eines Belegs die Slice-Kennung, die
  er belegt; ein reiner Rename lässt jede Zählung richtig aussehen und macht den
  Beleg falsch (gemessen als `lab/team-sim` s16c: Zähler unverändert, 0 Befunde,
  andere behauptete Herkunft).

## Ursache — vom Absender gelesen, nicht instrumentiert

Der Tree-Diff läuft mit den Default-Optionen der verwendeten Bibliothek, und
dort ist die Rename-Erkennung gesetzt. Ein Rename mit identischem Inhalt kommt
deshalb als **eine** Änderung mit `From` **und** `To` an und wird als
Modifikation auf dem **neuen** Pfad gebucht; die BASE-Version wird unter dem
neuen Pfad gesucht, nicht gefunden, und der Lauf endet still. Der Kommentar über
der Übersetzungs-Funktion sagt das Gegenteil; für den `--staged`-Pfad stimmt er,
weil dieser eine eigene Übersetzung hat.

## Was der Absender nicht beantragt

Keine neue Fähigkeit, keinen neuen Grund-Code, keine Vertrags-Änderung — und
ausdrücklich **keine** Ähnlichkeits-Erkennung. Die Bitte ist, dass der Rename im
Range-Pfad gar nicht erst als Rename gebucht wird. Welcher Schnitt das leistet,
überlässt er d-check.

## Eine Frage, keine Behauptung

Ob dies eine zweite Instanz von [`BEO-024`](../planning/observations.md) ist,
überlässt der Absender ausdrücklich uns.
