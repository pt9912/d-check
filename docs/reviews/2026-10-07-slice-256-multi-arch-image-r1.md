# Review R1 — slice-256: Multi-Arch-Image `linux/amd64` + `linux/arm64`

- **Review-Art:** Code. Geprüft gegen den Slice-Plan (`slice-256`, §1 Ziel und
  Abgrenzung, §3 Plan samt Spiegel-Liste, §6 Risiken), gegen
  [ADR-0102](../plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md)
  (Proposed), [ADR-0011](../plan/adr/0011-digest-pins-build-gate-images.md),
  [ADR-0065](../plan/adr/0065-spiegel-gleichheit-ist-der-config-digest.md) (wird
  abgelöst), gegen
  [`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image)/[`DC-FA-DIST-002`](../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)
  (Lastenheft 0.98.0) und die Hard Rules `AGENTS.md` §3.1, §3.5–§3.7, §3.9 sowie §5
  Regel 15. Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-256` · Range `e5ed67ee..16ebc6f9` (`f2ba9ad8` Spec/ADR,
  `91e4e650` Build/Test, `16ebc6f9` `release.yml`/`releasing.md`).
- **Skill:** `reviewer.md` @ 1.17.0.
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-07
- **Eingangs-Kontext:** Slice-Plan `slice-256`; ADR-0102, ADR-0011, ADR-0065;
  DC-FA-DIST-001/002; Action-Metadaten der gepinnten
  `docker/setup-qemu-action@99012661…` und `docker/setup-buildx-action@f87e5991…`
  (`action.yml`, per `gh api` gelesen). Vorherige Findings am selben Gegenstand: R1
  `slice-230` (Gate-Skript-Kommentare), Release-Prep-Reviews `v0.70.0`/`v0.71.0`
  (Spiegel-Prüfgröße, ADR-0065).
- **Proben (nur lesend, kein Push in echte Registries):** `make gate-consistency`,
  `make doc-check`, `make workflow-pins` — je `966 Datei(en) geprüft, 0 Befund(e)`;
  `make image-test` grün (`Plattform amd64 — d-check:latest`, Kriterien 1–4);
  `make image-test-arm64` auf einem Host ohne arm64-binfmt rot mit
  `Binary für arm64 auf diesem Host nicht ausführbar — binfmt/QEMU für arm64 fehlt`
  (make-Exit 2); die `image-digest-axis`-awk-Zeile gegen alle drei Präfixe
  (`gcr.io/distroless`, `golang:`, `golangci/`) — je Referenz und Digest korrekt;
  `index_digest()` aus dem Mirror-Schritt unter `set -euo pipefail` + `ERR`-Trap gegen
  eine nicht erreichbare Referenz (Skript im Scratchpad) — Ausgabe `TRAP-ERR`, Exit 1,
  die Zeile nach der Zuweisung wird **nicht** erreicht.

## Findings

| # | kategorie | befund | quelle | pfad | verifizierbar | klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Gegenprobe „gepushtes Binary == geprüftes" läuft **nach** `--push`, und der Push setzt bei stabilen Releases zugleich `:latest`. ADR-0102 führt für Option F als Pro „keine ungeprüfte Variante unter dem Tag" und verwirft Option E genau mit „eine ungeprüfte Variante läge bereits unter dem Versions-Tag, wenn die Prüfung fällt"; das Lastenheft sagt „jede Variante ist vor ihrer Veröffentlichung daraufhin geprüft". Im Fehlerfall der Gegenprobe trifft das Contra von E auf F ebenso zu — und die Konsequenzen der ADR nennen es nicht (der Code weiß es: die `ERR`-Meldung des Push-Schritts sagt „kann bereits ein Index liegen"). Failure-Szenario: der `amd64`-Teil des Index wird im buildx-Container-Builder (andere BuildKit-Instanz, kalter Cache) neu gebaut, das geprüfte `d-check:latest` stammt aus dem Daemon-Builder von `make ci`; weicht das Binary ab, ist das Release rot — `ghcr.io/…:v<version>` **und** `:latest` zeigen aber bereits auf den ungeprüften Index, und jeder `:latest`-Konsument zieht ihn. | ADR-0102 (§Verglichene Alternativen, Option F/E); [`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image); Skill-Anker 10/18 | `docs/plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md` · „keine ungeprüfte Variante unter dem Tag"; `spec/lastenheft.md` · „jede Variante ist vor ihrer Veröffentlichung daraufhin geprüft"; `Makefile` · „`--push .`" (Ziel `image-publish`) | nein — nur ein Lauf mit erzwungener Binary-Abweichung gegen eine Registry zeigt den Zustand nach dem Abbruch | Zusage „vor Veröffentlichung geprüft" gilt nur im Gutfall |
| F-2 | MEDIUM | Der Release-Pfad bezieht zwei neue Fremd-Images über bewegliche Tags: `setup-buildx-action` startet den Builder (`driver: docker-container`, `use: true` als Default) aus `moby/buildkit:buildx-stable-1`, `setup-qemu-action` registriert binfmt aus `docker.io/tonistiigi/binfmt:latest` (Default `image`). Die Actions sind SHA-gepinnt, die Images, die sie ziehen, nicht. Failure-Szenario: der `arm64`-Teil, der ausgeliefert wird, **und** das geprüfte `d-check:arm64` kommen beide aus demselben ungepinnten BuildKit-Container — die Gegenprobe vergleicht zwei Ausgaben desselben Builders und kann eine Änderung unter dem beweglichen Tag für `arm64` nicht sehen (für `amd64` vergleicht sie gegen den Daemon-Builder und sähe sie); der `arm64`-Test selbst läuft unter einem Emulator aus `:latest`. Gleicher Source-Tree, anderer Builder — die Restlücke, die ADR-0011 für Build-Images schließt. Ob der Release-Builder unter ADR-0011 Entscheidung 1 fällt (Wortlaut „jedes extern bezogene Image"; das `Schärft:`-Feld nennt nur `FROM`-Zeilen und semgrep), ist eine Architect-Frage; ADR-0102 §Konsequenzen nennt die Actions, nicht ihre Images. | ADR-0011 Entscheidung 1; `AGENTS.md` §3.9 (Begründung Supply-Chain); [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) | `.github/workflows/release.yml` · „uses: docker/setup-buildx-action@f87e5991a6d7451dcb8d9637bfbc97413f497069 # v4.4.1" und „uses: docker/setup-qemu-action@99012661954931238ded8c8b007157a8430204e1 # v4.4.0" | nein — `make workflow-pins` prüft die `uses:`-Form, nicht die von der Action gezogenen Images | Pin der Action, nicht ihres Images |
| F-3 | MEDIUM | Die Botschaft von `16ebc6f9` sagt „fehlende Referenz liefert leer (UNGEPRUEFT-Pfad)". Im Mirror-Schritt läuft `index_digest` unter `set -euo pipefail`: scheitert `imagetools inspect`, ist die Pipeline nicht-null, die Zuweisung `GHCR_IDX="$(…)"` bricht über `ERR` ab, und der `UNGEPRUEFT`-Zweig wird nie erreicht (Probe oben: `TRAP-ERR`, Folgezeile nicht erreicht). Gemessen wurde offenbar die Funktion allein, nicht der Schritt. Der Kommentar „Eine leere Antwort ist UNGEPRUEFT, nicht bestätigt" stimmt nur für einen Lauf mit Exit 0 und Nicht-Digest-Ausgabe. Failure-Szenario: Docker Hub antwortet auf das Inspect nicht (Rate-Limit, Propagation) — das Release ist rot (fail-closed bleibt), aber mit der generischen Trap-Meldung statt der zugesagten Diagnose beider Seiten, und die Botschaft behauptet einen Pfad, den der Schritt nicht nimmt. Die Struktur ist Bestand (`config_digest()` trug dieselbe, samt Kommentar „schweigt statt zu scheitern"), die Behauptung ist neu. | `AGENTS.md` §5 Regel 15; Skill-Anker 8 (`BEO-ALL/commit-message-overclaims-work`) | `.github/workflows/release.yml` · „GHCR_IDX=\"$(index_digest" ; Commit `16ebc6f9` · „fehlende Referenz liefert leer (UNGEPRUEFT-Pfad)" | ja — der Schritt-Rumpf unter `bash -euo pipefail` gegen eine unerreichbare Referenz | Messung an der Funktion, Aussage über den Schritt |
| F-4 | MEDIUM | Die Grenze-Liste von `image-test` nennt für `arm64` nur, dass binfmt fehlen kann — nicht, dass **beide** Seiten des Vergleichs (extrahiertes Binary und Container) unter demselben QEMU-User-Mode-Emulator auf einem `amd64`-Host laufen und die Plattform per `--platform` **erzwungen** wird. Das neue Akzeptanzkriterium „Boundary (Plattform)" spricht von einem `linux/arm64`-Host, der die Variante **selbst** wählt; beides misst der Lauf nicht (die Auswahl stützt sich auf das Plattform-Feld im Index, Prüfung (1) von `image-verify-published.sh`). Failure-Szenario: ein Verhalten, das nur auf echter arm64-Hardware oder einem arm64-Kernel mit anderer Seitengröße auftritt, ist für beide Seiten gleich unsichtbar — „nativ == Container" bleibt grün, und wer `make image-test-arm64` grün als Beleg des Kriteriums liest, liest mehr, als gemessen ist. | Skill-Anker 18 (Grenzen-Liste ohne größte Lücke); Skill-Anker 10; [`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image) | `harness/sensors/image-test.md` · „**`make image-test-arm64` braucht binfmt/QEMU für `arm64` auf dem Host.**" | nein | Emulation als Plattform-Beleg ohne benannte Grenze |
| F-5 | LOW | Der Kommentar in `tools/image-scan.sh` begründet die Zwei-Registry-Aufnahme mit „gleicher Config-Digest, gleicher Inhalt" — ein Spiegel der abgelösten Prüfgröße. Er steht weder in der `MR-025`-Liste des Plans (§3) noch in der Release-Prep-Menge (Handbuch, READMEs, `operations.md`, Hub-Overview); der Folge-Slice `slice-257` berührt die Datei, sein Plan nennt den Satz nicht. | [`MR-025`](../../harness/conventions.md#mr-025) | `tools/image-scan.sh` · „gleicher Config-Digest, gleicher Inhalt" | nein | Spiegel außerhalb der aufgelisteten Menge |
| F-6 | LOW | Der Kopf von `image-verify-published.sh` trägt Herkunft als Prosa mit zwei Kennungen, eine davon die von diesem Slice abgelöste ADR: „(ADR-0065 Punkt 3, fortgeführt in ADR-0102)". Nicht HIGH, weil der Satz eine Zusage trägt (Lesen aus der Registry) — abweichend ist nur die Form des Herkunfts-Felds. | `AGENTS.md` §3.7 (Herkunft nur als **ein** auflösbares Feld) | `tools/image-verify-published.sh` · „Punkt 3, fortgeführt in ADR-0102)" | nein | Herkunfts-Prosa im Kommentar |
| F-7 | INFO | Der Push-Schritt prüft die Form von `INDEX` nicht (anders als `index_digest()` im Mirror-Schritt); bei Exit 0 und leerer Ausgabe stünde `ghcr.io/pt9912/d-check@` als Pin im Job-Summary. Fail-closed bleibt es, weil `imagetools create` im Mirror-Schritt daran scheitert — die Meldung nennt dann aber genau diesen kaputten Pin als „BEREITS VEROEFFENTLICHT". | Maintainability | `.github/workflows/release.yml` · „DIGEST=\"$REGISTRY/$IMAGE_NAME@$INDEX\"" | nein | — |
| F-8 | INFO | Auf einem macOS-Host (Docker Desktop) ist das extrahierte Linux-Binary nie ausführbar; `make image-test` endet dort jetzt mit „binfmt/QEMU für arm64 fehlt" — eine Diagnose, die auf Darwin nicht zutrifft. Der frühere Kopf nannte die Host-Annahme ausdrücklich; die neue Grenze im Kopf spricht nur von binfmt. Rollen-Verweis: Release-Prep (Handbuch §`image-test` auf macOS steht ohnehin auf der Liste). | Maintainability | `tools/image-test.sh` · „binfmt/QEMU für $want_arch fehlt" | nein | — |
| F-9 | INFO | `image-test-arm64` und `image-publish` fehlen in der `.PHONY`-Liste, in der die übrigen Targets stehen. | Maintainability | `Makefile` · „.PHONY: nightly-state freshness-semgrep" | nein | — |

## Negativbefunde

- geprüft, ohne Befund: `Dockerfile` — `--platform=$BUILDPLATFORM` auf `deps` vererbt sich
  auf `compile`/`test`/`coverage`/`build`; die drei ersten sollen nativ laufen (Tests
  prüfen die Host-Plattform, wie bisher), `lint` hat ein eigenes `FROM` und kopiert nur
  den plattformneutralen Modul-Cache. `TARGETOS`/`TARGETARCH` leer ⇒ `GOOS=`/`GOARCH=`
  leer ⇒ Go-Default (Host); der Legacy-Builder ist durch `# syntax=docker/dockerfile:1.7`
  ohnehin ausgeschlossen. Die Runtime-Stufe hat kein `RUN`.
- geprüft, ohne Befund: `image-digest-axis` — die awk-Zeile liefert für alle drei Achsen
  Referenz und Digest (Probe oben); `make versions` gibt die `FROM`-Zeile jetzt mit
  `--platform=…` aus, rein kosmetisch. `pin-freshness.sh` prüft den Pin ohnehin auf
  Digest-Form.
- geprüft, ohne Befund: `release.yml`-Reihenfolge — `make ci` (Daemon-Builder) vor
  QEMU/buildx; `setup-buildx-action` setzt den Container-Builder als aktuellen
  (`use: true`), `docker build … --load` in `image-test-arm64` nutzt ihn und lädt das
  Einzel-Plattform-Bild in den Daemon; `image-publish` ruft `docker buildx build` auf
  denselben Builder (Cache-Treffer für `arm64`). Label-Pin vor dem Push an beiden lokalen
  Bildern; `GHCR_DIGEST` hat die Form `ghcr.io/pt9912/d-check@sha256:…` und taugt als
  Quelle für `imagetools create`; `imagetools` liest die Credentials aus der
  Docker-Konfiguration (GHCR-Login der Action, Hub-Login im `run`-Schritt).
- geprüft, ohne Befund: `:latest`-Logik — GHCR über `PUBLISH_LATEST="$IS_STABLE"` und
  `$(filter true,…)`, Hub über `TAGS+=(-t …:latest)` nur bei `IS_STABLE == true`;
  Prerelease setzt keines (ADR-0014).
- geprüft, ohne Befund: `image-verify-published.sh` auf stille Grün-Pfade — (1) `|| fail`
  an der Pipeline, Einzel-Manifest oder Attestations-Eintrag ⇒ Mengen-Ungleichheit ⇒
  Exit 1; (2) Zeilenzahl = 2 erzwungen, leerer Wert ⇒ Exit 1, `while … done < file` ohne
  Subshell, `fail` beendet das Skript; (3) `binary_sha` läuft im `||`-Kontext ohne
  `set -e`, behandelt `create`/`cp` aber explizit, die letzte Pipeline
  `sha256sum | cut` trägt ihren Status unter `pipefail`; eine leere SHA ist ohne
  Fehler nicht erreichbar.
- geprüft, ohne Befund: `image-test.sh` — `${PLAT[@]+"${PLAT[@]}"}` ist das korrekte
  Idiom für leere Arrays unter `set -u`; ELF-Byte 18 (`3e`/`b7`) trifft e_machine;
  kurze oder leere Datei ⇒ leeres `got_mach` ⇒ Exit 1; Exit 126 wird **vor** dem
  Container-Lauf abgefangen (gemessen).
- geprüft, ohne Befund: Gate-Index — `image-test-arm64` als Gate mit Bindepunkt
  `release.yml`, `image-publish` als Werkzeug mit `kein Gate` in der Zeile;
  `make gate-consistency` grün.
- geprüft, ohne Befund: `AGENTS.md` §3.9 — beide neuen `uses:` mit vollem SHA und
  Tag-Kommentar; `make workflow-pins` grün.
- geprüft, ohne Befund: `AGENTS.md` §3.5 — ADR-0065 unverändert; ADR-0102 ist neu,
  `Proposed`, mit Re-Evaluierungs-Trigger und Index-Zeile.
- geprüft, ohne Befund: Adressierungs-Form (Frage 16) — `Schärft:` nennt die
  `DC-*`-Kennungen und `SPEC-066`; Provenance-Marker im `Bezug:` zeigt nur auf den
  Vorgang (Frage 14).
- geprüft, ohne Befund: Spec-Straten (§3.4) — kein Abwärts-Token in Lastenheft/
  Spezifikation; `make doc-check` grün; Lastenheft-Version und Historie (`MR-032`).
- geprüft, ohne Befund: Abgrenzung — Handbuch, READMEs, `operations.md`, Hub-Overview
  unberührt (Release-Prep); `tools/image-scan.sh` nicht angefasst (slice-257); lokaler
  `make build` bleibt Single-Platform.
- geprüft, ohne Befund: Kommentar-Klassen (§3.7) in `Dockerfile`, `Makefile`,
  `release.yml`, `image-test.sh` und `dependabot.yml` — Zusage, Kopplung und Grenze;
  bis auf F-6 kein Befund.
- geprüft, ohne Befund: Botschaften `f2ba9ad8` und `91e4e650` gegen den Diff — die
  genannten Proben sind gedeckt; „image-test-arm64 ohne binfmt rot" in diesem Lauf
  nachgefahren.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 4 | 2 | 3 |

## Verdikt

**Merge-/Closure-blockierend: ja** — vier MEDIUM. F-1 und F-4 betreffen die Reichweite
der Zusage in ADR-0102 und im Lastenheft (Plan-/Entscheidungs-Frage, Rollen-Verweis
Architect); F-2 ist eine Pin-Frage mit Scope-Klärung über ADR-0011 (Architect); F-3 ist
eine Botschafts-/Kommentar-Korrektur. Kein HIGH: alle beschriebenen Fehlerpfade bleiben
fail-closed, keiner endet grün. Der Prerelease-Lauf (`v0.84.0-rc.1`) ist von diesem
Review nicht gedeckt.
