# `make image-scan` — CVE-Scan gegen die publizierten Images, nicht gegen den Arbeitsbaum

## Vertrag

Trivy, digest-gepinnt. **Anderer Gegenstand als alle übrigen Sensoren:** nicht
der Arbeitsbaum, sondern **was Anwender ziehen** — CVEs entstehen ohne Commit,
und gegen sie ist ein push-getriebenes Gate prinzipiell blind.

**Netz ist hier der Zweck, nicht ein Zugeständnis:** Eine gepinnte Vuln-DB
fände nur die CVEs von gestern. Der **Scanner** ist digest-gepinnt, die **DB**
bewusst nicht.

Je Registry **und je Plattform** des Index — die Plattformen liest das Skript
aus dem Index selbst, nicht aus einer Kopie der Release-Liste
([ADR-0102](../../docs/plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md));
`IMAGE_SCAN_PLATFORMS` übersteuert sie für einen gezielten Lauf. Je Plattform
drei Läufe: ein **Plattform-Nachweis**, der prüft, dass Trivy wirklich die
verlangte Architektur gescannt hat; ein Vollbericht über alle Schweregrade, der
nie fällt; und der Entscheidungslauf `CRITICAL`/`HIGH` **mit verfügbarem
Fix** — nur der macht rot.

## Grenze — was das Grün nicht abdeckt

1. **Der Fund-Raum ist klein und gemessen** — fünf OS-Pakete plus die
   Modul-Liste des Binaries. Ein grüner Lauf sagt „nichts Bekanntes in diesem
   Raum", **nicht** „das Image ist sicher".
2. **Alle Trivy-Läufe fahren `--exit-code 0`**, weil Trivy einen echten
   Fehler ebenfalls mit 1 quittiert — gemessen. Die Auswertung übernimmt das
   Skript.
3. **Die Plattform steht im Bild, nicht im Flag.** Gemessen: Trivy scannt ein
   Einzel-Manifest-Image bei `--platform linux/arm64` **still als amd64**, mit
   Exit 0. Der Plattform-Nachweis vergleicht deshalb die gemeldete
   Architektur; weicht sie ab oder fehlt die Plattform im Index, gilt der Scan
   als gescheitert (Exit 2), nicht als grün.
4. **Ohne lesbaren Index gibt es keinen Scan.** Ein Ref, dessen Plattformen
   sich nicht aus einem Multi-Plattform-Index lesen lassen (etwa ein
   Einzel-Manifest bis `v0.83.0`), gilt als gescheitert, nicht als gescannt;
   nur `IMAGE_SCAN_PLATFORMS` scannt ihn dann gezielt, und der
   Plattform-Nachweis aus 3. bleibt dabei in Kraft. Gelesen wird der Index mit
   `docker buildx imagetools` — eine Vorbedingung des Laufs; seine letzte
   Meldung steht mit in der Ausgabe und trennt ein Einzel-Manifest von einem
   fehlenden Ref oder einem Netz-Fehler.

`--selftest` prüft die Auswertung netzlos (sieben Proben zur Zählung, vier zur
Architektur); die Trivy-**Feldnamen** deckt er nicht. Fehlt das Feld
`architecture`, bleibt der Nachweis leer und der Scan gilt als gescheitert —
laut, nicht still.

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
