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

**Was dieses Repo damit hält** (die Mengen-Wahl steht begründet in
[`.d-check.yml`](../../.d-check.yml)): **jede ADR steht im ADR-Index** — die
Regel, die [`AGENTS.md`](../../AGENTS.md) §5 ausdrücklich führt und die zuvor
**kein** Gate hielt.

## Grenze

**Ein Block ist EIN Paar.** Die Ist-Menge wird als **Vereinigung** gelesen: Ein
Mitglied gilt als erwähnt, sobald es in *irgendeinem* Dokument steht. Wer zwei
voneinander unabhängige Invarianten in einen Block legt, hält **keine** von
beiden — gemessen: eine erste Fassung dieser Konfiguration führte ADRs und
Sensor-Dateien gemeinsam und meldete `108 von 108` auch dann noch, als eine ADR
vollständig aus ihrem Index entfernt war. **Die Sensors-Tabellen-Invariante ist
deshalb hier ausdrücklich NICHT gehalten**; sie bräuchte einen zweiten Lauf.

**Das Modul liest zwei Eingaben und scannt nur eine.** Die Ist-Dokumente werden
als Text gelesen; die Soll-Artefakte werden **nie geöffnet** — sie werden als
Pfade aus dem Dateibaum aufgesammelt. **Und eine zweite Achse:** aufgesammelt
wird aus dem **ganzen** Baum ab der Repository-Wurzel, unter der Skip-Liste und
`scan.ignore`, aber **nicht** eingeschränkt auf `scan.roots`.

**Geprüft wird eine Zeichenkette, nicht eine Rolle:** eine Nennung im Fenced
Block zählt, eine korrekte Beschreibung ohne Namensnennung zählt nicht. Ob eine
Erwähnung *trägt*, bleibt Urteil.

**Die Erkennungsform ist am Bestand zu wählen, und beide Prüfungen gehören
dazu.** `match: path` (Default) sucht den Pfad ab der Scan-Wurzel, `basename`
nur den Dateinamen. Für **Markdown-Dokumente, die relativ verlinken, ist `path`
die falsche Form** — gemessen: **84 von 84** ADRs und **24 von 24**
Sensor-Dateien wären gemeldet worden, ausnahmslos falsch. An einem
Fremd-Bestand lieferten dagegen **beide** Formen dasselbe Ergebnis, weil jenes
Handbuch mit vollem Pfad nennt. Wer `basename` wählt, prüft **zwei** Dinge:
dass die Basisnamen eindeutig sind **und** dass keiner Endstück eines anderen
ist — die zweite fehlte einmal und deckte einen echten Fehlalarm zu.

## Ausgabe und Ausgänge

Die Zusammenfassung nennt die **Bezugsmenge** (`N von M Artefakt(en) erwähnt,
über D Dokument(e)`) — ein Lauf ohne Befund sagt damit, worüber er nichts
gefunden hat. Im Default-Modus auf stderr, unter `--json`/`--yaml` als
`summary.notes`.

| Ausgang | Bedeutung |
| --- | --- |
| Exit 0 | jedes Mitglied der Soll-Menge kommt in der Ist-Menge vor |
| Exit 1 | mindestens ein `artifact-unmentioned` — der Befund nennt den Artefakt-Pfad; `line` ist der **Vertrags-Platzhalter** `1`, weil das Artefakt nicht geöffnet wird |
| Exit 2 | fail-closed: eine der beiden Mengen fehlt, ihr Glob trifft nichts, oder ein Verzeichnis unter der Wurzel ist nicht lesbar |

## Sperren

Netzlos (`--network none`), read-only gemountet, hermetisch — kein git, kein
Netz, kein Schreibpfad.

## Bindung

[ADR-0084](../../docs/plan/adr/0084-mentions-eigenes-modul.md) ·
[`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in).
**Bewusst nicht in `gates`/`ci`** — eine neue Modul-Klasse startet als
eigenständiger Fokus-Lauf, dieselbe Vorsicht wie bei
[`review-coverage`](review-coverage.md) und [`trace-check`](trace-check.md).
