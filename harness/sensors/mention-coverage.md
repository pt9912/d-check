# `make mention-coverage`

## Vertrag

Hält, dass jedes Mitglied einer konfigurierten **Soll-Menge** von Artefakten in
mindestens einem Dokument einer konfigurierten **Ist-Menge** vorkommt — Modul
`mentions` ([`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in),
[ADR-0084](../../docs/plan/adr/0084-mentions-eigenes-modul.md)), dogfooded über
das eigene Image. Beide Mengen sind Pfad-Globs; der Befund heißt
`artifact-unmentioned`.

**Die Achse ist eine andere als die der RTM.** Die RTM misst, ob eine
Anforderung **verfolgt** ist, und arbeitet über Kennungen; diese Achse misst,
ob ein Artefakt **erwähnt** ist, und arbeitet über Pfade — schema-frei.

**Was dieses Repo damit hält** (die Mengen-Wahl ist das Urteil, nicht die
Differenz — sie steht begründet in [`.d-check.yml`](../../.d-check.yml)): jede
ADR steht im ADR-Index, und jede Datei unter `harness/sensors/` ist aus der
Sensors-Tabelle verlinkt. Beides sind Invarianten, die dieses Repo führt und
die zuvor **kein** Gate hielt.

## Grenze

**Das Modul liest zwei Eingaben und scannt nur eine.** Die Ist-Dokumente werden
als Text gelesen; die Soll-Artefakte werden **nie geöffnet** — sie werden als
Pfade aus dem Dateibaum aufgesammelt, und ihre Mitgliedschaft ist eine Aussage
über das Verzeichnis, nicht über ihren Inhalt.

**Geprüft wird das wörtliche Vorkommen einer Zeichenkette**, nicht ihre Rolle:
eine Nennung im Fenced Block zählt, eine korrekte Beschreibung ohne Namensnennung
zählt nicht. Ob eine Erwähnung *trägt*, bleibt Urteil.

**Die Erkennungsform ist eine Wahl mit gemessenem Preis.** `match: path`
(Default) sucht den Pfad ab der Scan-Wurzel, `basename` nur den Dateinamen. Für
**Markdown-Dokumente, die relativ verlinken, ist `path` die falsche Form** —
gemessen an diesem Repo: 84 von 84 ADRs und 23 von 23 Sensor-Dateien wären
gemeldet worden, ausnahmslos falsch, weil der Index `0084-….md` schreibt und
nicht `docs/plan/adr/0084-….md`. Der Default bleibt trotzdem der strengere,
weil ein bloßer Dateiname kollidieren kann; wer relativ verlinkt, wählt
`basename` **und prüft, ob seine Basisnamen eindeutig sind**.

## Ausgabe und Ausgänge

Die Zusammenfassung nennt die **Bezugsmenge** (`N von M Artefakt(en) erwähnt,
über D Dokument(e)`) — ein Lauf ohne Befund sagt damit, worüber er nichts
gefunden hat. Im Default-Modus steht sie auf stderr, unter `--json`/`--yaml`
als `summary.notes`.

| Ausgang | Bedeutung |
| --- | --- |
| Exit 0 | jedes Mitglied der Soll-Menge kommt in der Ist-Menge vor |
| Exit 1 | mindestens ein `artifact-unmentioned` — der Befund nennt den Artefakt-Pfad; `line` ist der **Vertrags-Platzhalter** `1`, weil das Artefakt nicht geöffnet wird |
| Exit 2 | fail-closed: eine der beiden Mengen fehlt in der Konfiguration, oder ihr Glob trifft nichts |

## Sperren

Netzlos (`--network none`), read-only gemountet, hermetisch — kein git, kein
Netz, kein Schreibpfad.

## Bindung

[ADR-0084](../../docs/plan/adr/0084-mentions-eigenes-modul.md) ·
[`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in).
**Bewusst nicht in `gates`/`ci`** — eine neue Modul-Klasse startet als
eigenständiger Fokus-Lauf, dieselbe Vorsicht wie bei
[`review-coverage`](review-coverage.md) und [`trace-check`](trace-check.md).
