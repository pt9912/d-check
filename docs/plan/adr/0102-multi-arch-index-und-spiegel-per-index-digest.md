# ADR-0102: Das Image ist ein Multi-Arch-Index, und der Spiegel kopiert ihn samt Index-Digest

**Status:** Accepted

**Datum:** 2026-10-07

**Autor:** pt9912

**Bezug:** [`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[`DC-FA-DIST-002`](../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)
(beide geändert, Lastenheft 0.98.0); [ADR-0002](0002-distribution-ghcr-image.md),
[ADR-0011](0011-digest-pins-build-gate-images.md) (§3: die Basis-Pins sind schon
Index-Digests), [ADR-0014](0014-latest-tag-fuer-stabile-releases.md);
slice-256 <!-- d-check:status-provenance -->.

**Supersedes:** [ADR-0065](0065-spiegel-gleichheit-ist-der-config-digest.md)

**Schärft:** [`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
und [`DC-FA-DIST-002`](../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)
— **wie** der Index entsteht, wie jede Plattform vor dem Push geprüft wird und
woran die Spiegel-Gleichheit gemessen wird; dazu die §6-Zeile
[`SPEC-066`](../../../spec/spezifikation.md#6-externe-verträge).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Das Image ist bisher ein einzelnes `linux/amd64`-Manifest. Auf Apple Silicon
und arm64-Linux läuft es nur unter Emulation. Der Auftraggeber verlangt ein
Image, das `docker pull` auf beiden Architekturen ohne Emulation auswählt —
einen **OCI-Index**.

Daran hängen drei Annahmen des Bestands, die alle Single-Arch voraussetzen:

- [ADR-0065](0065-spiegel-gleichheit-ist-der-config-digest.md) sagt als
  Spiegel-Gleichheit **den** Config-Digest zu. Ein Index hat keinen einzelnen;
  die Prüfung in `release.yml` fände in einem Index kein `config`-Feld und
  bräche fail-closed ab — **nach** dem GHCR-Push.
- Der Spiegel tagt und pusht **dasselbe lokale Bild**. Ein Multi-Plattform-Index
  liegt nicht als ein Bild im lokalen Daemon.
- `tools/image-test.sh` vergleicht den Container mit dem Binary **der
  Host-Architektur**. Eine zweite Plattform bliebe ungeprüft.

**Vorab gemessen** (lokal, gegen zwei Wegwerf-Registries `registry:2`, kein Push
nach GHCR oder Docker Hub):

| Messung | Ergebnis |
|---|---|
| Cross-Compile beider Plattformen (`FROM --platform=$BUILDPLATFORM`, `GOOS`/`GOARCH` aus `TARGETOS`/`TARGETARCH`) | 51 s, **keine** Emulation; die Runtime-Stufe löst ihre Plattform über den vorhandenen Index-Pin |
| `amd64`-Binary des Cross-Builds gegen das aus `make build` | sha256 **gleich** (`c46c1576…`) |
| `docker buildx imagetools create -t <B> <A>` | Index-Digest in beiden Registries **gleich** (`sha256:0b369a06…`) |
| Labels je Plattform über `imagetools inspect --format '{{json .Image}}'` | lesbar, `version` je Plattform gesetzt |

ADR-0065 hatte `imagetools create` als Option B verworfen, weil die Erhaltung
des Digests **ungemessen** war. Die Messung liegt jetzt vor — gegen
`registry:2`, nicht gegen Docker Hub. Den Docker-Hub-Teil misst der erste
Prerelease-Lauf, bevor diese ADR `Accepted` wird.

## Entscheidung

Wir wählen einen **Index aus `linux/amd64` und `linux/arm64`, dessen Varianten
vor dem Push geprüft sind und der nach Docker Hub als Index kopiert wird**.

1. **Gebaut wird per Cross-Compile.** Die Go-Stufen laufen auf der
   Build-Plattform und kompilieren mit `GOOS`/`GOARCH` für die Ziel-Plattform;
   die Runtime-Stufe ist Distroless ohne `RUN` und braucht keine Emulation. Der
   Index entsteht nur im Release-Pfad; `make build` und die Gates bleiben bei
   einem Bild der Host-Plattform — ein Multi-Plattform-Bild lässt sich ohne
   containerd-Image-Store nicht lokal laden.
2. **Jede Plattform ist vor dem Push geprüft.** `amd64` wie bisher über
   `make ci`; `arm64` über `make image-test-arm64`: dasselbe Prüfskript, gegen
   die `arm64`-Variante, Binary und Container unter QEMU/binfmt. Das Skript
   prüft dabei die ELF-Architektur des extrahierten Binaries — sonst prüfte ein
   Lauf, der still die Host-Variante zieht, die falsche Plattform.
3. **Das veröffentlichte Binary ist das geprüfte — und erst dann gibt es einen
   Tag.** Prüfung und Push sind getrennte buildx-Aufrufe. Der Index wird
   deshalb zuerst **ohne Tag** gepusht (`push-by-digest`), dann wird je
   Plattform das Binary aus dem **gepushten** Index gezogen und gegen das
   geprüfte verglichen (sha256), und erst danach zeigen `v<version>` und
   gegebenenfalls `:latest` per `imagetools create` auf genau diesen Digest.
   Weicht ein Binary ab, bleibt jeder Tag unberührt; in der Registry liegt
   dann nur der ungetaggte Index, erreichbar über den Digest, den die Meldung
   nennt.
4. **Der Index trägt nur Plattform-Manifeste** (`--provenance=false`,
   `--sbom=false`). Attestations-Einträge (`unknown/unknown`) würden von der
   Plattform-Prüfung als dritte Variante gesehen.
5. **Die Labels werden je Plattform aus der Registry gelesen**, nicht aus dem
   lokalen Daemon; `org.opencontainers.image.version` muss je Plattform der
   Tag-Version entsprechen.
6. **Der Spiegel kopiert den Index** (`docker buildx imagetools create` von der
   GHCR-Referenz `@<index-digest>`), statt ein lokales Bild neu zu pushen. Die
   Gleichheits-Größe ist der **Index-Digest**, aus beiden Registries gelesen —
   er bindet alle Plattformen samt Inhalt und ist schärfer als ein Config-Digest.
7. **Der Konsumenten-Pin ist der Index-Digest** und damit auf beiden Registries
   derselbe. Die registry-lokale Ausnahme aus ADR-0065 Punkt 2 entfällt; die
   Doku sagt das in der Release-Prep.
8. **Die Images der Build-Actions sind digest-gepinnt.** Der buildkit-Builder
   (`setup-buildx-action`) baut das ausgelieferte Binary, das binfmt-Image
   (`setup-qemu-action`) trägt den arm64-Test; beide bezögen sonst einen
   beweglichen Tag. Die Gegenprobe sähe eine Änderung im Builder nicht, weil
   geprüftes und gepushtes Binary aus demselben Builder kommen — dieselbe
   Pin-Pflicht wie für jedes Fremd-Image
   ([ADR-0011](0011-digest-pins-build-gate-images.md)).
9. **Unverändert übernommen aus [ADR-0065](0065-spiegel-gleichheit-ist-der-config-digest.md):**
   Punkte 3 bis 7 — aus den Registries lesen; Zugangsdaten vor dem Verbrauch
   prüfen, Login als `run`-Schritt; die Darstellung macht das Release nicht rot;
   Zeichen statt Bytes; kopieren statt neu bauen, fail-closed mit Nennung des
   GHCR-Stands, GHCR zuerst.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon (Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR
(MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (nur `amd64`) | keine Arbeit | Apple-Silicon- und arm64-Hosts laufen emuliert; die Anfrage bleibt offen |
| B — zwei getrennte Tags je Architektur (`:vX-arm64`) | kein Index, der bisherige Spiegel bliebe | der Konsument wählt selbst; ein Digest-Pin in einem geteilten Makefile passt nur auf eine Architektur |
| C — Index, Build emuliert (ohne `--platform=$BUILDPLATFORM`) | kein Dockerfile-Umbau | die Go-Stufen liefen unter QEMU — um ein Vielfaches langsamer, das 30-min-Budget des Release-Jobs wäre gefährdet |
| D — Index, Spiegel weiter per Config-Digest je Plattform | bewährte Größe | N Vergleiche statt einem, und ein neu gepushter Index hätte wieder registry-lokale Digests — zwei Pins für denselben Inhalt |
| E — Index, arm64 nur nach dem Push prüfen | einfacher Ablauf | eine ungeprüfte Variante läge bereits unter dem Versions-Tag, wenn die Prüfung fällt |
| **F — Cross-Compile-Index, beide Plattformen vor dem Push geprüft, Index ohne Tag gepusht und gegengeprüft, erst dann getaggt, Spiegel als Index-Kopie (gewählt)** | gemessene Eigenschaften; ein Pin für beide Registries; keine ungeprüfte Variante unter einem Tag | zwei Builds (Prüfung, Push) brauchen die Gegenprobe; fällt sie, bleibt ein ungetaggter Index in der Registry; QEMU-Laufzeit für den arm64-Test im Release-Job |

## Konsequenzen

- **Positiv:** ein Pin für beide Registries; native Ausführung auf arm64.
- **Positiv:** die Gegenprobe macht die byte-gleiche Reproduktion des Binaries
  zu einer geprüften Eigenschaft statt einer Annahme.
- **Negativ:** der Release-Job braucht QEMU und einen buildx-Builder — zwei
  neue Actions, SHA-gepinnt und über Dependabot gehoben. Ihre Images
  (`moby/buildkit`, `tonistiigi/binfmt`) sind digest-gepinnt, haben aber
  **keine** Frische-Achse: Dependabot liest `with:`-Eingaben nicht, und der
  Nachtlauf kennt die beiden Pins nicht — eine benannte Lücke: gehoben werden
  sie von Hand.
- **Negativ:** der lokale `make image-test-arm64` braucht binfmt für `arm64` auf
  dem Host; ohne bricht er fail-closed mit Hinweis ab.
- **Offen:** `make image-scan` scannt bis zum Folge-Vorgang nur die
  Host-Variante des Index.
- **Negativ:** mit dem Status `Superseded by ADR-0102` meldet `matrix` jeden
  Link auf ADR-0065 als `matrix-inactive`. Lebende Verweise ziehen auf diese
  ADR um; [ADR-0068](0068-lokale-workflow-referenzen-ohne-pin.md) aber ist
  `Accepted` und nennt ADR-0065 als ihren Anlass. Sie kommt deshalb in
  `matrix.exempt-paths`, wie zuvor ADR-0047 — gemessener Preis: außer den drei
  ADR-0065-Links trägt ihr Körper nur zwei Links auf `AGENTS.md`, die heute
  keinen Befund ergeben.
- **Folgepflicht:** Handbuch, READMEs, `operations.md` und Hub-Overview
  sagen Index-Digest statt Config-Digest; ADR-0065 auf
  `Superseded by` setzen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `tools/image-test.sh` mit `PLATFORM=linux/arm64` | arm64-Binary (ELF-Maschine geprüft) und arm64-Container liefern byte-gleiche Ausgabe und gleichen Exit-Code | `make image-test-arm64` |
| `tools/image-publish.sh` mit `tools/image-verify-published.sh` | Index erst ohne Tag gepusht, Tags erst nach grüner Gegenprobe; der Index trägt genau `linux/amd64` und `linux/arm64`; je Plattform ist das Binary im gepushten Index sha256-gleich zum geprüften; Labels je Plattform gesetzt, `version` = Tag | `make image-publish` (Release-Pfad) |
| `release.yml`, Schritt *Mirror to Docker Hub* | Index-Digest von `docker.io/…:v<version>` und `ghcr.io/…:v<version>`, aus den Registries gelesen, sind gleich — sonst Abbruch mit Nennung des GHCR-Stands | — (Tag-Push) |

## Re-Evaluierungs-Trigger

- Der Prerelease-Lauf zeigt, dass Docker Hub den Index-Digest beim
  `imagetools create` **nicht** erhält — dann trägt Entscheidung 6 nicht, und
  die Gleichheit fällt auf den Config-Digest je Plattform zurück (Option D).
- Der arm64-Test unter QEMU sprengt das Zeitbudget des Release-Jobs.
- Eine weitere Plattform wird verlangt.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-07 | Angelegt als `Proposed`; `Accepted` erst mit der Closure des Vorgangs, nach Review, Verifikation und dem Prerelease-Lauf |
| 2026-10-07 | Nach R1 (vier MEDIUM): Entscheidung 3 auf „ohne Tag pushen, prüfen, dann taggen" umgestellt, Entscheidung 8 (Builder-Images digest-gepinnt) ergänzt, Option F und Konsequenzen nachgezogen. Noch `Proposed`, Körper daher geändert statt angehängt |
| 2026-10-07 | Nach Verifikation (V-2 bis V-4) und R2 (R2-1 bis R2-4) nachgezogen; die Läufe `v0.84.0-rc.1` und `v0.84.0` haben die Docker-Hub-Hälfte von Entscheidung 6 gemessen (derselbe Index-Digest auf beiden Registries) und den arm64-Test unter QEMU grün in 4,5 min. Der Statuswechsel von ADR-0065 brachte die Ausnahme für ADR-0068 (§Konsequenzen). `Accepted` mit der Closure des Vorgangs |
