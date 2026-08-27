# ADR-0065: Die Spiegel-Gleichheit ist der Config-Digest, nicht der Manifest-Digest

**Status:** Accepted

**Datum:** 2026-08-27

**Autor:** pt9912

**Bezug:** [`DC-FA-DIST-002`](../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel), [`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image), [ADR-0064](0064-dockerhub-spiegel-fail-closed.md), [ADR-0002](0002-distribution-ghcr-image.md), [ADR-0014](0014-latest-tag-fuer-stabile-releases.md)

**Supersedes:** [ADR-0064](0064-dockerhub-spiegel-fail-closed.md)

**Schärft:** [`DC-FA-DIST-002`](../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel) — **welche** Größe die zugesagte Gleichheit ist und woran sie gemessen wird.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0064](0064-dockerhub-spiegel-fail-closed.md) sagte zu, das gespiegelte Bild
trage **denselben Manifest-Digest** wie das GHCR-Original, und ließ die Pipeline
das fail-closed prüfen. Der umsetzende Slice hatte die Frage selbst als Risiko
notiert — *„`docker tag` + `push` sollte denselben Manifest-Digest liefern;
‚sollte' ist keine Zusage"* — und sie dann **in** der Pipeline beantwortet statt
**vor** ihr.

Ein unabhängiger Review hat sie beantwortet, und die Antwort ist ein Nein.
Gemessen am Schwester-Repo `d-migrate`, das mit **identischem Verfahren**
spiegelt (ein lokal gebautes Bild, `docker tag` je Registry, `docker push`):

| Größe | Docker Hub | GHCR |
| ----- | ---------- | ---- |
| Image-ID (lokal) | `sha256:e886cdbd…` | `sha256:e886cdbd…` — **gleich** |
| **Config**-Digest (aus der Registry gelesen) | `sha256:e886cdbd…` | `sha256:e886cdbd…` — **gleich** |
| **Manifest**-Digest | `sha256:02b87477…` | `sha256:8f5abc4f…` — **verschieden** |

**Drei von drei** geprüften Tag-Paaren tragen zwei verschiedene
Manifest-Digests. Die Ursache steht im Manifest: gleicher `mediaType`, gleicher
Config-Digest — also identischer Inhalt —, aber **verschiedene Layer-Blobs**
(29 751 109 gegen 30 612 704 Byte für dieselbe Schicht): die Signatur einer
Neu-Kompression beim zweiten Push.

Die Annahme, unter der ADR-0064 stand, ist damit gekippt — und mit ihr die
Entscheidung. Hätte die Pipeline so ausgeliefert, wäre **jedes** Release
gebrochen: beide Pushes gelingen, der Vergleich schlägt fehl, der Job bricht ab,
und weil `Create GitHub Release` danach steht, existierte das Release nicht,
während der Spiegel bereits öffentlich liegt.

## Entscheidung

Wir wählen **den Config-Digest als zugesagte Gleichheits-Größe**, aus beiden
Registries gelesen.

1. **Die Gleichheit ist der Config-Digest.** Er ist die Identität des
   Bild-**Inhalts** und über Registries hinweg stabil, weil er über die
   Konfiguration und die *unkomprimierten* Schicht-Kennungen gebildet wird. Der
   Manifest-Digest hängt an der Blob-Kompression und ist registry-lokal.
2. **Der Manifest-Digest ist ausdrücklich registry-lokal, und das gehört in die
   Doku.** Wer `docker.io/…@sha256:…` pinnt, nimmt den Docker-Hub-Digest; der
   GHCR-Digest löst dort nicht auf. Das ist eine Eigenschaft von
   Registry-Digests, keine des Spiegels — verschwiegen wäre es eine Falle.
3. **Gelesen wird aus den Registries, nicht aus dem lokalen Daemon.** Zwei aus
   demselben lokalen Bild abgeleitete Werte sind trivial gleich; ein Vergleich,
   der nicht fehlschlagen kann, ist ein stiller Grün-Pfad.
4. **Die Zugangsdaten werden geprüft, bevor sie verbraucht werden**, und der
   Login läuft als `run`-Schritt. Eine `uses:`-Action scheitert mit ihrem
   eigenen Text; ein `trap` im *nächsten* Schritt liefe nie an — das
   Negative-Kriterium der Anforderung verlangt aber, dass die Meldung den
   bereits veröffentlichten GHCR-Stand benennt.
5. **Die Darstellung kann das Release nicht rot machen.** Das Zeichen-Limit der
   Kurzbeschreibung ist eine Eigenschaft einer **Repo-Datei** und damit
   hermetisch prüfbar: es wandert nach `make gates`, wo es beim Schreiben
   auffällt statt beim Veröffentlichen. Beide Release-Schritte der Darstellung
   tragen `continue-on-error`.
6. **Gemessen werden Zeichen, nicht Bytes.** Docker Hubs Limit ist ein
   Zeichen-Limit; eine Byte-Messung wäre bei Umlauten strenger als die Regel und
   erzeugte Falsch-Rot an einem Text, den Docker Hub annimmt.
7. **Unverändert übernommen aus [ADR-0064](0064-dockerhub-spiegel-fail-closed.md):**
   Spiegeln heißt kopieren statt neu bauen; fail-closed für den **Spiegel**, mit
   einer Meldung, die den Teil-Zustand benennt; Reihenfolge GHCR zuerst;
   Bild-Name über `vars` überschreibbar; Darstellung aus Repo-Dateien mit der
   Kategorie als Text.

**Korrektur einer Zuschreibung.** [ADR-0064](0064-dockerhub-spiegel-fail-closed.md)
nannte `d-migrate` ein *„fail-open-Muster"*, dessen Release *„von fremder
Verfügbarkeit unabhängig"* bleibe. Am Workflow gemessen trägt nur die erste
Hälfte: `d-migrate` überspringt den Spiegel bei **fehlendem Secret**. Ist das
Secret gesetzt und Docker Hub gestört, hat sein Push-Schritt **kein**
`continue-on-error` — es wird ebenfalls rot. Der Unterschied liegt im fehlenden
Secret, nicht in der Störung.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (ADR-0064 stehen lassen) | keine Arbeit | die Zusage ist am Bestand widerlegt und bräche **jedes** Release |
| B — Manifest-Digest zusagen, `buildx imagetools create` statt `tag`+`push` | wäre eine ausdrückliche Kopier-Operation | ungemessen, ob sie den Digest wirklich erhält — ohne zweite Registry zum Prüfen wäre das die **zweite** ungeprüfte Zusage in Folge, also genau der Fehler, den diese ADR korrigiert |
| C — Gleichheit gar nicht zusagen, beide Digests nur dokumentieren | ehrlich, kein Bruchrisiko | die Zusage stünde ohne Sensor; ein Konsument wüsste nicht, ob er dasselbe Werkzeug zieht — das Muster, das dieses Repo als stillen Grün-Pfad führt |
| D — Zeichen-Limit im Release-Schritt belassen | ein Prüf-Ort weniger | ein Schreibfehler fiele erst beim Veröffentlichen auf und machte ein gültiges Release rot |
| **E — Config-Digest zusagen, aus beiden Registries gelesen (gewählt)** | gemessene statt plausible Eigenschaft; der Vergleich kann fehlschlagen und tut es bei einem echten Fehler | der Manifest-Digest bleibt registry-verschieden — Konsumenten müssen den passenden Pin nehmen, und das muss die Doku sagen |

## Konsequenzen

- **Positiv:** die Zusage ist eine **gemessene** Eigenschaft statt einer
  plausiblen. Der Vergleich liest zwei registry-zurückgemeldete Werte und ist
  damit kein Ritual.
- **Negativ:** ein Konsument, der Digest-Pins zwischen den Registries kopiert,
  läuft auf. Das ist eine Eigenschaft von Registry-Digests, keine des Spiegels —
  die Doku sagt es deshalb ausdrücklich, statt es der Entdeckung zu überlassen.
- **Negativ:** die erste Fassung wurde veröffentlicht, bevor die Zusage gemessen
  war. Der Fehler saß nicht im Nichtwissen — das Risiko war benannt —, sondern
  in der Reihenfolge.
- **Folgepflicht:** [`DC-FA-DIST-002`](../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)
  auf den Config-Digest umschreiben; das Zeichen-Limit als Go-Test in
  `make gates`; der registry-lokale Manifest-Digest in Handbuch und beide
  READMEs; [ADR-0064](0064-dockerhub-spiegel-fail-closed.md) auf
  `Superseded by` setzen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `release.yml`, Schritt *Mirror to Docker Hub* | Config-Digest von `docker.io/…:v<version>` und `ghcr.io/…:v<version>`, **aus den Registries gelesen**, sind gleich — sonst Abbruch mit Nennung des veröffentlichten GHCR-Stands | — (Tag-Push, kein lokales Target) |
| Go-Test `TestDockerHubDescriptionLimit` | `packaging/dockerhub/description.txt` ist nicht leer und höchstens 100 **Zeichen** | `make test` (in `make gates`) |

## Re-Evaluierungs-Trigger

Ändert Docker Hub oder GHCR die Blob-Behandlung so, dass auch der
**Manifest**-Digest über Registries hinweg stabil wird, ist die schärfere Zusage
wieder möglich — und dann erst zu **messen**, nicht anzunehmen.

Zweiter Trigger, unverändert aus
[ADR-0064](0064-dockerhub-spiegel-fail-closed.md): fällt das Release **zweimal**
an einer Docker-Hub-Störung aus, während GHCR grün war, ist die
fail-closed-Entscheidung neu zu stellen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-08-27 | Accepted — ersetzt [ADR-0064](0064-dockerhub-spiegel-fail-closed.md), nachdem ein unabhängiger Review dessen tragende Annahme am Bestand widerlegt hat | [slice-165](../planning/done/slice-165-dockerhub-spiegel.md) |
