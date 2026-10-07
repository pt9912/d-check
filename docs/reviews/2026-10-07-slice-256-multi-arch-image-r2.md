# Review R2 — slice-256: Multi-Arch-Image `linux/amd64` + `linux/arm64`

- **Review-Art:** Code. Geprüft gegen den Slice-Plan (`slice-256`, §3 samt
  Plan-Änderung im Lauf, §6 Risiken), gegen
  [ADR-0102](../plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md)
  (Proposed, Fassung nach R1),
  [ADR-0011](../plan/adr/0011-digest-pins-build-gate-images.md), gegen
  [`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image)/[`DC-FA-DIST-002`](../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)
  und die Hard Rules [`AGENTS.md`](../../AGENTS.md) §3.7, §3.9 sowie §5 Regeln 13
  und 15. Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-256`, **nur die Einarbeitungen seit R1** — Range
  `16ebc6f9..61a539b8` ohne `docs/reviews` (`22816afd` R1-Fix, `61a539b8`
  Verify-Fix). Anlass: Verify-Befund V-1 (die Veröffentlichungs-Mechanik lag
  außerhalb der R1-Range).
- **Skill:** `reviewer.md` @ 1.17.0.
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-07
- **Eingangs-Kontext:** Slice-Plan `slice-256`; ADR-0102, ADR-0011; R1-Report
  (`2026-10-07-slice-256-multi-arch-image-r1.md`, F-1..F-9) und Verify-Report
  (`2026-10-07-slice-256-multi-arch-image-verify.md`, V-1..V-5) als vorherige
  Findings am selben Gegenstand.
- **Proben (lokal, kein Push in echte Registries):** Wegwerf-`registry:2` auf
  `127.0.0.1:5105`, danach abgeräumt. `d-check:latest` (amd64) und
  `d-check:arm64` gepusht, per `imagetools create` zu einem Index gefügt
  (`sha256:5de3414c…`).
  - `tools/image-verify-published.sh` mit `REF=<repo>@<index-digest>` **und**
    mit `REF=<repo>:<tag>` — beide grün, Binary je Plattform gleich.
  - `tools/image-publish.sh` mit einem `docker`-Shim im Scratchpad, der nur
    `buildx build` ersetzt (schreibt `containerimage.digest` in die
    Metadaten-Datei) und auf Wunsch `imagetools create` scheitern lässt; alles
    andere ging an das echte `docker`. Gutfall mit `PUBLISH_LATEST=true`:
    beide Tags gesetzt und nachgeprüft, letzte Zeile
    `image-publish: 127.0.0.1:5105/d-check@sha256:5de3414c…`. Bruch über die
    Label-Version (`VERSION=9.9.9`): `FAIL — Gegenprobe rot — KEIN Tag gesetzt`,
    `:v9.9.9` danach `not found`. Bruch in `imagetools create`: Exit 1, **ohne**
    `image-publish: FAIL`-Zeile (siehe R2-1).
  - Daemon-Verhalten (klassischer `overlay2`-Store): zweites
    `docker pull --platform` derselben **Tag**-Referenz für die andere
    Plattform → Exit 0 (Tag umgehängt); derselben **Digest**-Referenz →
    `cannot overwrite digest sha256:5de3414c…`, Exit 1 (siehe R2-2).
  - `sed`-Extraktion des Push-Schritts gegen ein synthetisches Log: die Zeilen
    `(1) Index ohne Tag gepusht — …` und `image-verify-published: OK — …`
    treffen **nicht**, die Schlusszeile trifft; ohne Schlusszeile leer → der
    Leer-Zweig greift. `false | tee` unter `set -euo pipefail` + `ERR`-Trap →
    Trap feuert, Exit 1.
  - Builder-Image-Pins per `imagetools inspect`: `moby/buildkit:v0.33.1` und
    `tonistiigi/binfmt:qemu-v10.2.3` sind Index-Digests (OCI-Index), gleich den
    Pins in `release.yml` (ADR-0011 Entscheidung 3).
  - `make doc-check`, `make gate-consistency`, `make workflow-pins` — je
    `968 Datei(en) geprüft, 0 Befund(e)`.

## Findings

| # | kategorie | befund | quelle | pfad | verifizierbar | klasse |
|---|---|---|---|---|---|---|
| R2-1 | LOW | Die Trap-Meldung des Push-Schritts verweist auf „die Meldung von image-publish oben", die „den Stand" nenne. Das stimmt nur für die zwei `fail()`-Pfade (Digest unlesbar, Gegenprobe rot, Tag-Nachprüfung). Scheitert `docker buildx build` in Schritt (1) oder `imagetools create` in Schritt (3), endet das Skript über `set -e` ohne eigene Zeile — gemessen: nur `ERROR: simulated push failure`, Exit 1. Ebenso nennt die Grenze im Kopf von `image-publish.sh` nur den Ausgang von (2). Dass nach einem Abbruch in (3) `v<version>` bereits gesetzt sein kann und `:latest` noch nicht, steht dort nicht. Kein Fehlverhalten: der Lauf bleibt rot, und jeder gesetzte Tag zeigt auf den geprüften Digest. Ungenau ist nur die Diagnose. | [`AGENTS.md`](../../AGENTS.md) §5 Regel 13; Skill-Anker 18 | `.github/workflows/release.yml` · „die Meldung von image-publish oben nennt den Stand"; `tools/image-publish.sh` · „Grenze: fällt (2), liegt der ungetaggte Index" | ja — Shim-Probe wie oben (`imagetools create` mit Exit 1) | Diagnose-Zusage gilt nur für die eigenen `fail()`-Pfade |
| R2-2 | LOW | Der Kommentar vor dem Plattform-Abruf sagt allgemein: „der Daemon bindet eine Index-Referenz an genau ein lokales Bild, ein zweites Ziehen derselben Referenz für die andere Plattform scheitert". Gemessen gilt das für eine **Digest**-Referenz (`cannot overwrite digest`). Für eine **Tag**-Referenz, die der Skript-Kopf ausdrücklich zulässt („Tag oder `<repo>@<digest>`"), gelingt das zweite Ziehen (Tag umgehängt). Gemessen wurde am klassischen Store, nicht am containerd-Store. Dieselbe Aussage steht in der Botschaft von `22816afd`. Der Code arbeitet in beiden Formen richtig (beide Proben grün); die Begründung trifft nur auf eine zu. | [`AGENTS.md`](../../AGENTS.md) §5 Regeln 13/15 | `tools/image-verify-published.sh` · „ein zweites Ziehen derselben Referenz für die andere Plattform scheitert" | nein (Kommentar) | Begründung an einer Form gemessen, für alle gesagt |
| R2-3 | INFO | F-8 ist im Code geschlossen (eigene Meldung für einen Nicht-Linux-Host, Verweis auf den macOS-Absatz in `releasing.md`, der existiert). Der Kopf von `image-test.sh` nennt als Grenze aber weiter nur binfmt/QEMU für eine fremde Plattform. Die Host-Annahme „Linux" steht nur im Code-Zweig. | Maintainability | `tools/image-test.sh` · „Grenze: eine Plattform, die nicht die des Hosts ist, braucht" | nein | — |
| R2-4 | INFO | Die beiden neuen Builder-Pins (`moby/buildkit`, `tonistiigi/binfmt`) erscheinen nicht in `make versions`, dem Beleg der Reproduzierbarkeits-Pins nach ADR-0011 Entscheidung 5. Das ist keine Verletzung: Entscheidung 5 nennt nur `FROM`-Zeilen und semgrep, und ADR-0102 Entscheidung 8 überträgt die Pin-Pflicht, nicht die Beleg-Pflicht. Die Konsequenzen von ADR-0102 nennen die fehlende Frische-Achse, nicht diese zweite Lücke. | Maintainability | `docs/plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md` · „Die Images der Build-Actions sind digest-gepinnt." | nein | — |

## R1-Befunde — sachlich geschlossen?

| R1 | Stand nach R2 |
|---|---|
| F-1 | **geschlossen.** Push ohne Tag (`push-by-digest=true,name-canonical=true`), Gegenprobe am Digest, Tags erst danach, jeder gesetzte Tag nachgeprüft (V-4). Bruch gemessen: kein Tag gesetzt. ADR-0102 Entscheidung 3 und Option F nennen den Rest (ungetaggter Index) als Contra. Restunschärfe der Diagnose: R2-1. |
| F-2 | **geschlossen.** Beide Images per Index-Digest gepinnt (gegen die Registry geprüft). Die Scope-Frage ist durch ADR-0102 Entscheidung 8 selbst entschieden. Die fehlende Frische-Achse steht in den Konsequenzen und als offenes Risiko im Plan. |
| F-3 | **geschlossen.** `{ … \|\| true; }` in `index_digest()`. Der Kommentar nennt den Grund richtig: der Leer-Zweig wird unter `ERR`-Trap erreicht. |
| F-4 | **geschlossen.** Grenze 3 in der Sensor-Datei von `image-test` nennt dasselbe QEMU für beide Seiten, die erzwungene Plattform und als Belegtes nur die Vorbedingung (Index trägt `linux/arm64`). |
| F-5 | **geschlossen** (V-2 ebenso): „gleicher Inhalt, ab ADR-0102 auch derselbe Index-Digest". Der Satz schreibt der alten Messung keine falsche Eigenschaft mehr zu. |
| F-6 | **geschlossen:** nur noch `(ADR-0102)` als Rang-Zeiger, dazu eine Zusage samt Grund. |
| F-7 | **geschlossen:** Digest aus der Schlusszeile, Form per `sed` geprüft, Leer-Zweig mit Fehlermeldung (Probe oben). |
| F-8 | **im Code geschlossen**; Kopf: R2-3. |
| F-9 | **geschlossen:** `image-test-arm64` und `image-publish` in `.PHONY`. |

## Negativbefunde

- geprüft, ohne Befund: **stille Grün-Pfade in `image-publish.sh`.** Build-Abbruch
  über `set -e`; unlesbarer Digest → `fail`; jede Abweichung der Gegenprobe → `fail`
  vor jedem Tag; Tag-Nachprüfung je gesetztem Tag gegen den geprüften Digest, leere
  Antwort → `fail`. Kein Pfad endet mit Exit 0 ohne geprüften Digest unter den Tags.
- geprüft, ohne Befund: **stille Grün-Pfade in `image-verify-published.sh`.** Fehlendes
  Plattform-Manifest → `fail`. Der Abruf über `<repo>@<manifest-digest>` ist
  inhaltsadressiert: auch wo der Daemon Layer lokal wiederverwendet, ist das gelesene
  Binary das unter diesem Digest gepushte. Die Ableitung von `repo` ist für Tag-Form
  (auch mit Registry-Port) und Digest-Form korrekt (Proben mit `127.0.0.1:5105`).
- geprüft, ohne Befund: **Push-Schritt in `release.yml`.** `make … | tee` unter
  `pipefail` erreicht den Trap. Die Extraktion trifft nur die Schlusszeile, nicht die
  Zwischenzeilen mit gleichem Digest. Leer-Zweig mit Fehlermeldung.
- geprüft, ohne Befund: **Spiegel-Schritt.** `index_digest()` mit `|| true`; die
  Gleichheit wird aus beiden Registries gelesen, eine leere Seite → `UNGEPRUEFT` und
  Exit 1.
- geprüft, ohne Befund: **SHA-Pin-Regel §3.9 und ADR-0011 Entscheidung 3.** `uses:`
  unverändert SHA-gepinnt; die neuen `with:`-Images sind Index-Digests mit lesbarem
  Tag.
- geprüft, ohne Befund: **Makefile.** `image-publish` reicht Variablen als
  Einzel-Token durch (`PROGRESS_FLAG := --progress=plain`). Die Prüfung auf ein
  fehlendes `PUBLISH_REPO` steht vor jedem Netzzugriff. Kommentar und
  `##`-Beschreibung geben die Reihenfolge Push ohne Tag → Gegenprobe → Tags wieder.
- geprüft, ohne Befund: **Kommentar-Klassen (§3.7)** in `image-publish.sh`,
  `image-verify-published.sh`, `image-test.sh`, `image-scan.sh` (geänderter Satz),
  `release.yml` und `Makefile`. Sie tragen Zusage, Kopplung (Schlusszeile ↔
  `sed`-Extraktion) und Grenze; keine Review-Historie, keine Slice-Nummer neu.
  Ausnahme ist die Aussage-Genauigkeit in R2-1/R2-2.
- geprüft, ohne Befund: **Doku gegen Code.** `releasing.md` Schritt 5, die
  Werkzeug-Zeile von `make image-publish` im Gate-Index, ADR-0102 Entscheidungen 3/8,
  Option F, Konsequenzen und Fitness-Function-Zeile. Die Plan-Änderung im Lauf deckt
  `image-publish.sh` und die Builder-Pins.
- geprüft, ohne Befund: **Botschaft von `61a539b8`** gegen den Diff (V-2, V-3, V-4;
  „beide Tags geprueft" in diesem Lauf mit `PUBLISH_LATEST=true` nachgefahren).
  Die Botschaft von `22816afd` ist bis auf die Daemon-Begründung (R2-2) gedeckt.
- geprüft, ohne Befund: **ADR-Status.** ADR-0102 ist `Proposed`, die Körper-Änderung
  ist zulässig und in `## Geschichte` vermerkt.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 2 | 2 |

## Verdikt

**Merge-/Closure-blockierend: nein.** Alle neun R1-Befunde sind sachlich geschlossen,
F-8 mit einem Rest im Skript-Kopf (R2-3). Die neue Veröffentlichungs-Mechanik hat
keinen stillen Grün-Pfad. Die Zusage „unter keinem Tag liegt eine ungeprüfte
Variante" ist im Gut- und in zwei Bruchfällen gemessen. R2-1 und R2-2 sind
Aussage-Genauigkeiten in Meldung bzw. Kommentar: annehmen oder begründen. R2-3 und
R2-4 brauchen keine Aktion. Den Prerelease-Lauf (`v0.84.0-rc.1`) gegen GHCR und
Docker Hub deckt dieses Review nicht.
