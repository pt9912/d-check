# `make runtime-base-digest` u. a. — meldet, ob ein digest-gepinntes Fremd-Image unter demselben Tag neu gebaut wurde

## Vertrag

Sechs Achsen, ein Sensor: die drei `Dockerfile`-Stages, das semgrep-Gate-Image,
das a-check-Image und das Trivy-Image des CVE-Scans. Sie beantworten eine
**andere Frage** als die Versions-Achsen — nicht „gibt es einen neueren Tag",
sondern „trägt derselbe Tag inzwischen einen anderen Digest".

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
5. **`trivy-digest` läuft nicht im Nachtlauf** — `upstream-drift.yml` führt die
   anderen fünf Achsen; der Scanner-Digest wird nur gemeldet, wenn jemand die Achse
   selbst ruft.

Quelle, Digest-Form, Vergleich, fail-open und Exit legt
[`SPEC-099`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
fest; diese Datei führt sie nicht ein zweites Mal.

## Bindung

kein Gate — meldet, urteilt nicht über den Arbeitsbaum; die Hebung bleibt ein
bewusster Akt. **Netz**, fail-open, **nicht** in `gates`; gerufen vom Nachtlauf
(außer `trivy-digest`, Grenze 5).
[ADR-0011](../../docs/plan/adr/0011-digest-pins-build-gate-images.md) ·
[`SPEC-099`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
