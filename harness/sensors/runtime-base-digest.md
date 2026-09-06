# `make runtime-base-digest` u. a. — meldet, ob ein digest-gepinntes Fremd-Image unter demselben Tag neu gebaut wurde

## Vertrag

Fünf Achsen, ein Sensor: die drei `Dockerfile`-Stages, das semgrep-Gate-Image
und das a-check-Image. Sie beantworten eine **andere Frage** als die
Versions-Achsen — nicht „gibt es einen neueren Tag", sondern „trägt derselbe
Tag inzwischen einen anderen Digest".

**Dass eine Versions-Achse darüber genügte, ist gemessen falsch:**
`make freshness-go` meldete `ok`, während `golang:1.27.0` einen anderen Digest
trug als der Pin. Für das Runtime-Image ist diese Achse zusätzlich die
**einzige** Handhabe — sein Tag `nonroot` führt keine Version.

Quelle ist `docker buildx imagetools inspect`, **registry-agnostisch**:
derselbe Aufruf trifft gcr.io, Docker Hub und GHCR, ohne dass hier Token-Flüsse
je Registry gepflegt werden. Für gcr.io allein ginge auch ein `curl` auf den
`docker-content-digest`-Header — der Grund für `imagetools` ist die **Menge**
der Registries, nicht die Unmöglichkeit des Handbetriebs.

## Grenze — was das Grün nicht abdeckt

1. **Verglichen wird der Digest der adressierten Referenz** — bei einer
   Multi-Plattform-Liste der Listen-Digest, dieselbe Größe, die im
   `Dockerfile` steht. Permanent, gewollt.
2. **Das Urteilswort ist `ABWEICHEND`, nicht `VERALTET`** — Digests haben
   **keine Ordnung**; der Sensor kann nicht sagen, welcher der neuere ist.
   Permanent.
3. **Docker ist Voraussetzung des Zweigs** — `imagetools inspect` kennt keine
   eigene Zeitgrenze, deshalb steht ein `timeout` davor; fehlt Docker, ist das
   ein `SKIP`, kein Befund.
4. **Fail-open** — jede Netz- oder Werkzeugstörung endet als `SKIP` mit Exit 0.
   Ein Sensor, der bei fremder Störung rot wird, wird abgeschaltet.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | Pin entspricht dem Upstream-Digest — oder `SKIP` (fail-open) |
| 3 | `ABWEICHEND`: derselbe Tag trägt upstream einen anderen Digest |

Das sind die Codes des **Skripts**; `make` normalisiert einen fehlgeschlagenen
Recipe auf seinen eigenen Exit 2. Welcher Fall vorliegt, sagt die Ausgabe.

## Bindung

kein Gate — meldet, urteilt nicht über den Arbeitsbaum; die Hebung bleibt ein
bewusster Akt. **Netz**, fail-open, **nicht** in `gates`; gerufen vom Nachtlauf.
[ADR-0011](../../docs/plan/adr/0011-digest-pins-build-gate-images.md)
