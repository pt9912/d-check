# `make blackbox-probe` — vergleicht das Verhalten vor und nach einer Änderung von außen

## Vertrag

Baut aus einer Git-Referenz (`REF=<ref>`) ein **Vorher**-Image — `git archive`,
Dockerfile und Build-Defaults dieses Stands — und fährt es zusammen mit dem
**Nachher**-Image (`make build`) über dieselben Fixtures und Ausgabeformen.
Verglichen werden **stdout, stderr und Exit getrennt**; zusammengeführte
Streams verschränken sich und zeigen Abweichungen, die keine sind.

- **Fixtures:** jedes Unterverzeichnis von `tools/blackbox-probe/fixtures/`
  (je ein kleines Repository im Default-Zustand eines Moduls) und das Repo
  selbst.
- **Formen:** Standard, `--json`, `--yaml`, `--doctor` (`PROBE_FORMS`
  übersteuert).

Gedacht für die Zusage „ohne Schalter unverändert": vor dem Review laufen
lassen, mit `REF` auf den Stand vor der Änderung.

## Grenze — was das Grün nicht abdeckt

1. **Es vergleicht nur, was Fixtures und Formen abdecken.** Ein Modul ohne
   Fixture ist nur über den Lauf auf dem Repo selbst erfasst; was dessen
   Konfiguration nicht einschaltet, bleibt unverglichen. `ls
   tools/blackbox-probe/fixtures/` zeigt den Ausschnitt.
2. **Es findet einen Unterschied, den niemand erwartet — nicht die falsche
   Erwartung.** Ob „unverändert" die richtige Zusage ist, bleibt beim fremden
   Leser.
3. **Kein Gate.** Eine gewollte Änderung erzeugt Abweichungen; ob eine
   Abweichung gewollt ist, entscheidet der Vorgang, der sie erzeugt.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | byte-identisch über alle Vergleiche |
| 1 | mindestens eine Abweichung — je Fall die Ströme und ein Diff-Auszug |
| 2 | Lauf gescheitert (`REF` fehlt oder ist kein Commit, Vorher-Image nicht baubar, Nachher-Image fehlt, leere Formen oder Fixtures) |

**Das sind die Codes des Skripts.** `make` normalisiert jeden fehlgeschlagenen
Recipe auf seinen eigenen Exit 2; 1 und 2 trennt die **Ausgabe**.

## Bindung

Werkzeug, kein Gate: nicht in `gates`, `ci` oder einem Nachtlauf.
