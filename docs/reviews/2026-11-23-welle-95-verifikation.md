# Verifikation welle-95 — Closure-Trigger, Schritt 1

**Rolle:** Verifier (Modul 8, §Rollen-Sequenz für eine Welle: Verifier → Planner) ·
**Gegenstand:** Closure-Trigger aus `docs/plan/planning/welle-95-linelength-modul.md` §3, geprüft
auf `f98ca836` mit sauberem Arbeitsbaum · **Bezug:**
[`DC-FA-LLEN-001`](../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in),
[ADR-0114](../plan/adr/0114-linelength-modul.md) · Datum 2026-11-23.
Jeder Beleg stammt aus einem Lauf, den ich selbst gefahren habe. Das Trigger-Audit ist Schritt 2
(Planner/Architect) und hier nicht geprüft. Der Beleg geht über die DoDs der Slices hinaus; die
Slice-Verifikation steht in [2026-11-23-slice-290-verifikation.md](2026-11-23-slice-290-verifikation.md).

## Trigger-Punkte

| Punkt | Urteil | Beleg |
|---|---|---|
| 1 — alle Slices in `done/` | bestätigt | `docs/plan/planning/done/` enthält `slice-289-linelength-anforderung-und-adr.md` und `slice-290-linelength-modulkern.md`; `open/`, `next/` und `in-progress/` enthalten keinen Slice, nur `in-progress/roadmap.md`. Einen Verweis auf `welle-95-linelength-modul` tragen in den Lifecycle-Verzeichnissen nur diese beiden Slices und die Roadmap. Der Move von slice-290 ist `9ec16df2`. |
| 2 — `make fullbuild` grün | bestätigt | Exit 0. Auszug: `verify ok (54 Dateien, vollständig)`; `0 issues.`; `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%`; `Ran 55 rules on 75 files: 0 findings.`; `[gates] … green`; `image-test: OK — Akzeptanzkriterien erfüllt (…, DC-FA-LLEN-001)`; `bench: Median 1074 ms (Pass-Kriterium < 5000 ms — DC-QA-01)`; `58 Anforderung(en), 0 Waise(n).`; verify-closure-notes `1131 Datei(en) geprüft, 0 Befund(e)`. Abschluss: `[fullbuild] green — image-hash sha256:7ffebd74950ad2639eae893e2373848a5719c2e2aebe997915979aaf4fbff4c5` (gleich `docker image inspect d-check:latest` direkt danach). |
| 3 — ohne den Block nichts Neues | bestätigt | `make blackbox-probe REF=v0.91.0` auf `f98ca836`, Exit 0: `byte-identisch über 20 Vergleiche (Vorher 58c9486d, stdout/stderr/Exit getrennt, Kanarienlauf vorher und nachher)`. Verglichen sind die Fixtures `ids`, `links`, `sauber`, `targets` und der Eigenlauf `repo`, je in Default, `--json`, `--yaml` und `--doctor`. Das eigene Profil trägt keinen `linelength`-Block. Zusätzlich lief das Fixture von Phase (9) ohne Konfigurationsdatei gegen `d-check:probe-vorher` und das Nachher-Image (`--network none`, `--user 12345:12345`, `:ro`), in denselben vier Formen: `4/4 gleich`, Exit 0. |
| 4 — die Entscheide tragen am Image | bestätigt | In `make fullbuild` lief `make image-test` mit Phase (9) (`b792b305`) und meldete `image-test: (9) linelength — 40 grün, 41 rot, Codepunkte statt Bytes, Fenced-Code/Tabelle/Wort über der Grenze ausgenommen und je einzeln schaltbar (:ro, fremder Benutzer, ohne Netz)`. Das Skript gibt je Lauf nur diese Zeile aus. Die Einzelergebnisse unten stammen deshalb aus denselben vier Läufen, die ich nachgefahren habe: gleiches Fixture (aus `tools/image-test.sh` extrahiert), gleiche Konfiguration und gleiche Aufrufform, gegen das Image `sha256:7ffebd74…`. |

Den Vergleich in Phase (9) habe ich am Skript gelesen, nicht gebrochen: Die Liste der Befunde wird
**als Ganzes** auf Gleichheit geprüft (`cut -f1-3` gegen die erwartete Liste). Jeder zusätzliche oder
fehlende Befund macht den Lauf rot, zum Beispiel die Mehrbyte-Zeile, wenn Bytes gezählt würden.
Ebenso prüft das Skript Länge und Grenze in der Meldung und die Bezugsmengen-Zeile auf stderr
zeichengenau. Das Rot der Regel selbst ist in der Slice-Verifikation per Mutation belegt (M01–M24).

## Probe → Ergebnis am Image

Fixture `probe.md` aus Phase (9) mit `max-chars: 40`. Die Zählung habe ich mit `awk` gemacht, in
Codepunkten unter `C.UTF-8` und in Bytes unter `C`. Gestartet wurde mit
`docker run --rm --network none --user 12345:12345 -v <fixture>:/repo:ro d-check:latest --enable linelength`.

| Probe | Zeile | Länge | Default (Lauf l0) | mit Schalter aus |
|---|---|---|---|---|
| knapp unter der Grenze | Z. 2 | 39 | kein Befund | — |
| genau auf der Grenze | Z. 3 | 40 | kein Befund | — |
| knapp über der Grenze | Z. 4 | 41 | `probe.md:4  max-chars=40  line-chars-exceeded  Zeile hat 41 Zeichen, erlaubt sind 40` | in allen vier Läufen derselbe Befund |
| Mehrbyte an der Grenze | Z. 5 | 40 Codepunkte, 90 Bytes | kein Befund | — |
| Mehrbyte über der Grenze | Z. 6 | 41 Codepunkte, 91 Bytes | Befund, `41 … 40` | in allen vier Läufen derselbe Befund |
| überlange Zeile in Fenced-Code | Z. 8 (Zäune Z. 7/9) | 62 | kein Befund | l1 `exempt-code-blocks: false`: `probe.md:8 … Zeile hat 62 Zeichen, erlaubt sind 40` |
| überlange Tabellenzeile | Z. 10 | 63 | kein Befund | l2 `exempt-tables: false`: `probe.md:10 … Zeile hat 63 Zeichen, erlaubt sind 40` |
| lange URL (Wort mit 55) | Z. 11 | 66 | kein Befund, zählt in U | l3 `exempt-unbreakable: false`: `probe.md:11 … Zeile hat 66 Zeichen, erlaubt sind 40` |

| Lauf | Exit | Befunde | Bezugsmenge auf stderr |
|---|---|---|---|
| l0 Default | 1 | Z. 4, 6 | `linelength: 1 Datei(en), 7 Zeile(n) geprüft, 1 davon mit einem Wort über der Grenze ausgenommen` |
| l1 `exempt-code-blocks: false` | 1 | Z. 4, 6, 8 | `… 10 Zeile(n) geprüft, 1 davon …` |
| l2 `exempt-tables: false` | 1 | Z. 4, 6, 10 | `… 8 Zeile(n) geprüft, 1 davon …` |
| l3 `exempt-unbreakable: false` | 1 | Z. 4, 6, 11 | `… 7 Zeile(n) geprüft, 0 davon …` |
| Fixture ohne Konfiguration, Vorher gegen Nachher | 0 | keine | keine Zeile; byte-gleich in vier Formen |

Das sind genau die Befunde, die `DC-FA-LLEN-001` zusagt: Die Grenze selbst ist befundfrei, gezählt
werden Codepunkte, und die drei Ausnahmen gelten im Default und lassen sich je einzeln abschalten.
Die Bezugsmenge zählt die geprüften Zeilen nach (a) und (b) und dazu die über (c) ausgenommenen.

## Ergebnis

Die vier geprüften Trigger-Punkte sind bestätigt. Image-Hash des Laufs:
`sha256:7ffebd74950ad2639eae893e2373848a5719c2e2aebe997915979aaf4fbff4c5`. `make blackbox-probe` baut
`d-check:latest` danach mit `VERSION=0.0.0-dev` neu; der Hash oben ist der aus `make fullbuild`. Offen
für Schritt 2 und die folgenden Closure-Schritte (Planner/Architect): das Trigger-Audit und
`welle-95-results.md`, die die Tabelle oben übernimmt.
