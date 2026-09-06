# `make image-scan` — CVE-Scan gegen die publizierten Images, nicht gegen den Arbeitsbaum

## Vertrag

Trivy, digest-gepinnt. **Anderer Gegenstand als alle übrigen Sensoren:** nicht
der Arbeitsbaum, sondern **was Anwender ziehen** — CVEs entstehen ohne Commit,
und gegen sie ist ein push-getriebenes Gate prinzipiell blind.

**Netz ist hier der Zweck, nicht ein Zugeständnis:** Eine gepinnte Vuln-DB
fände nur die CVEs von gestern. Der **Scanner** ist digest-gepinnt, die **DB**
bewusst nicht.

Zwei Läufe je Image: ein Vollbericht über alle Schweregrade, der nie fällt, und
der Entscheidungslauf `CRITICAL`/`HIGH` **mit verfügbarem Fix** — nur der macht
rot.

## Grenze — was das Grün nicht abdeckt

1. **Der Fund-Raum ist klein und gemessen** — fünf OS-Pakete plus die
   Modul-Liste des Binaries. Ein grüner Lauf sagt „nichts Bekanntes in diesem
   Raum", **nicht** „das Image ist sicher".
2. **Beide Trivy-Läufe fahren `--exit-code 0`**, weil Trivy einen echten
   Fehler ebenfalls mit 1 quittiert — gemessen. Die Auswertung übernimmt das
   Skript.

`--selftest` prüft die Auswertung netzlos (sieben Proben); die
Trivy-**Feldnamen** deckt er nicht.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | sauber |
| 1 | behebbare Befunde (`CRITICAL`/`HIGH` mit Fix) |
| 2 | **Scan gescheitert** — ausdrücklich **kein** grüner Befundstand |

**Das sind die Codes des Skripts.** `make` normalisiert jeden fehlgeschlagenen
Recipe auf seinen eigenen Exit 2 — über das Target sind 1 und 2 damit nicht
unterscheidbar; der Nachtlauf liest deshalb die **Ausgabe**, nicht den
Exit-Code.

## Bindung

Netz, bewusst **nicht** in `gates`; gerufen vom Nachtlauf
[`image-scan.yml`](../../.github/workflows/image-scan.yml). Kein Docker-Socket.
[ADR-0066](../../docs/plan/adr/0066-cve-scan-gegen-das-publizierte-image.md)
