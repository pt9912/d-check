# Verifikation slice-256 — Multi-Arch-Image `linux/amd64` + `linux/arm64` (DoD)

- **Rolle:** Verifier (Modul 11), Frage „Bauen wir es richtig?". Geprüft wurde gegen §2 DoD, §1 Ziel und Abgrenzung und §3 Plan des Slice-Plans (samt der vermerkten Plan-Änderung), gegen [`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image)/[`DC-FA-DIST-002`](../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel) (Lastenheft 0.98.0), [`SPEC-066`](../../spec/spezifikation.md#6-externe-verträge) und [ADR-0102](../plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md) (Proposed). Review-Entscheidungen waren nicht der Maßstab.
- **Gegenstand:** `e5ed67ee..22816afd`, fünf Commits: `f2ba9ad8` (Spec/ADR), `91e4e650` (Build/Test), `16ebc6f9` (`release.yml`), `b4a11db2` (R1-Report), `22816afd` (R1-Einarbeitung F-1 bis F-9).
- **Sensor-Evidence:** Alles selbst gemessen, nichts aus dem Implementer-Bericht übernommen: `make gates`, `make image-test`, `make image-test-arm64` (fail-closed-Pfad), `make trace-check`/`make adr-check` über die Range, `make review-coverage`; der Release-Pfad lokal gegen zwei Wegwerf-`registry:2` (`127.0.0.1:5101`, `:5102`) mit einem buildx-Builder `--driver docker-container --driver-opt network=host`; die Schritte *Push to GHCR* und *Mirror to Docker Hub* aus `release.yml` wörtlich per `awk` extrahiert und unter `bash` gefahren. Kein Push in echte Registries, kein Tag. Registries, Builder und die eigenen `localhost:5101`-/Wegwerf-Bilder sind danach entfernt.
- **Modell-ID:** claude-opus-5-5 (Claude Agent SDK)
- **Datum:** 2026-10-07

---

## DoD-Prüfung je Punkt (§2)

### 1. Lastenheft DIST-001/002, Version/Historie, neue ADR samt Index, `SPEC-066`: **ERFÜLLT**

- DIST-001 nennt Index mit `linux/amd64`/`linux/arm64`, Identität zur nativen Ausführung **je Plattform**, „jede Variante ist vor ihrer Veröffentlichung daraufhin geprüft"; neues Kriterium „Boundary (Plattform)"; weitere Plattformen im Out-of-Scope.
- DIST-002 wechselt auf den Index-Digest (Beschreibung und Happy Path); der Out-of-Scope-Punkt „andere Plattform-Matrix" bleibt als Inhalts-Gleichheit stehen, „Gleichheit des Manifest-Digests" entfällt begründet.
- Version `0.97.2` → `0.98.0`, Historie-Zeile vorhanden; kein Abwärts-Token (`make doc-check` grün).
- ADR-0102 neu, `Proposed`, `Supersedes: ADR-0065`, Re-Evaluierungs-Trigger, Index-Zeile in `docs/plan/adr/README.md`. ADR-0065 unverändert (`make adr-check RANGE=e5ed67ee..HEAD` → `0 Befund(e)`); ihr `Superseded`-Übergang ist laut ADR-0102 Folgepflicht beim `Accepted` — kein Befund.
- `SPEC-066`-Zeile nennt Index-Digest, beide Plattformen, Cross-Compile; Historie-Zeile in `spezifikation.md`.

### 2. `Dockerfile` cross-kompiliert; `image-test.sh` prüft benannte Plattform; Target für die Nicht-Host-Plattform im Index: **ERFÜLLT**

- `Dockerfile`: `FROM --platform=$BUILDPLATFORM golang:…@sha256:…` und `GOOS=${TARGETOS} GOARCH=${TARGETARCH}` in `build`.
- **Byte-Gleichheit des amd64-Binaries über zwei Builder** (nachgemessen, nicht übernommen): `d-check:latest` aus `make build` (Daemon-Builder) und die `linux/amd64`-Variante aus dem Container-Builder liefern beide `c46c1576799d…`; ebenso `d-check:arm64` (Daemon) und die gepushte `linux/arm64`-Variante (Container-Builder) beide `57b87de8a0e5…`.
- `make image-test` grün: `image-test: Plattform amd64 — d-check:latest`, Kriterien (1)–(4), `image-test: OK — DC-FA-DIST-001-Akzeptanzkriterien erfüllt`.
- `make image-test-arm64` (Host ohne arm64-binfmt): Bild entsteht (`docker image inspect d-check:arm64` → `arm64 0.0.0-dev`), dann `image-test: Plattform arm64 — d-check:arm64` / `image-test: FAIL — Binary für arm64 auf diesem Host nicht ausführbar — binfmt/QEMU für arm64 fehlt`, make-Exit 2. Fail-closed, kein stilles Überspringen. Der grüne arm64-Lauf ist lokal nicht belegbar (kein QEMU, Installation nicht erlaubt) — er gehört zum Prerelease-Lauf.
- **Bewusstes Brechen der ELF-Prüfung:** Wegwerf-Bild `FROM scratch` mit dem **amd64**-Binary, gebaut als `--platform linux/arm64` (Image-Metadaten `arm64`). `IMAGE_REF=… PLATFORM=linux/arm64 bash tools/image-test.sh` → `image-test: FAIL — Binary in vfy-elf-mismatch:arm64 ist ELF-Maschine 0x3e, verlangt arm64 (0xb7)`, Exit 1. Die Prüfung, die ADR-0102 Entscheidung 2 trägt („sonst prüfte ein Lauf, der still die Host-Variante zieht, die falsche Plattform"), wird aus dem richtigen Grund rot.
- `make image-test-arm64` steht in `harness/README.md` §Sensors als Gate mit Bindepunkt `release.yml`; `make image-publish` als Werkzeug mit `kein Gate`. `make gate-consistency` grün (in `gates`).
- Plan-Änderung (awk in `image-digest-axis`): nachgemessen — die neue awk-Zeile liefert für `golang:`, `golangci/`, `gcr.io/distroless` je Referenz und Digest; die alte Form (`index($2,p)==1`) liefert für `golang:` **leer**, hätte `go-base-digest` also still auf SKIP gesetzt. Der Nachzug ist nötig und korrekt.

### 3. `release.yml`: Index-Build, Prüfung vor Push, Label-Check je Plattform, Index-Digest-Pin, Spiegel, Gleichheit, SHA-Pins: **ERFÜLLT bis auf den Prerelease-Lauf — OFFEN, ausstehend**

Reihenfolge in `release.yml`: `make ci` (amd64-Test) → QEMU/buildx → `make image-test-arm64` → GHCR-Login → Label-Pin an beiden lokalen Bildern → `make image-publish` (Push ohne Tag → Gegenprobe → Tags) → Hub-Login → Spiegel → GitHub-Release. Selbst gefahren gegen `registry:2`:

| Pfad | Lauf | Ergebnis |
|---|---|---|
| grün | `BUILDX_BUILDER=vfy256 make image-publish PUBLISH_REPO=localhost:5101/d-check` | Exit 0; `(1) Index ohne Tag gepusht — …@sha256:a61a6b79…`, `(1) Plattformen — linux/amd64 linux/arm64`, `(2) Labels — je Plattform gesetzt, version=0.0.0-dev`, `(3)` beide Binaries gleich, `(3) getaggt — v0.0.0-dev`; Registry-Tags `["v0.0.0-dev"]` |
| grün, stabil | dasselbe mit `PUBLISH_LATEST=true` | Exit 0; **derselbe** Index-Digest `a61a6b79…` (Wiederholung reproduzierbar); `v0.0.0-dev` und `latest` zeigen beide auf `sha256:a61a6b79…` |
| Bruch: anderes Binary (arm64-Bild mit `VERSION=9.9.9` als `TESTED_ARM64`) | `tools/image-publish.sh` gegen `localhost:5101/d-check-brk` | Exit 1; `FAIL — linux/arm64: Binary im gepushten Index (57b87de8…) ist nicht das geprüfte aus d-check:arm64-brk (fcf342b9…)` → `FAIL — Gegenprobe rot — KEIN Tag gesetzt; der ungetaggte Index liegt unter …@sha256:a61a6b79…`. Registry: **kein Tag** (`NAME_UNKNOWN` auf `tags/list`), Index per Digest erreichbar (HTTP 200). Schließt R1 F-1 nachweislich. |
| Bruch: Einzel-Manifest statt Index | `image-verify-published.sh` mit `REF=…@<amd64-Manifest>` | Exit 1; `FAIL — … nicht lesbar oder kein Multi-Plattform-Index` |
| Bruch: Attestations-Einträge | Index mit `--provenance=true` nach `:5102` | Exit 1; `FAIL — Index … trägt [linux/amd64 linux/arm64 unknown/unknown unknown/unknown], verlangt [linux/amd64 linux/arm64]` |
| Bruch: Label-Version | `VERSION=1.2.3` gegen den `0.0.0-dev`-Index | Exit 1; `FAIL — org.opencontainers.image.version=[0.0.0-dev] auf linux/amd64 entspricht nicht der Tag-Version [1.2.3]` |
| Bruch: Registry nicht erreichbar | `REF=localhost:5109/nix:v1` | Exit 1; `FAIL — … nicht lesbar oder kein Multi-Plattform-Index` |
| Push-Schritt (extrahiert) | mit dem Log des grünen Laufs | `digest=localhost:5101/d-check@sha256:a61a6b79…` im `GITHUB_OUTPUT` |
| Push-Schritt, letzte Zeile fehlt (F-7) | Log ohne die Pin-Zeile | Exit 1; `::error::Digest-Pin nicht aus der Ausgabe von image-publish lesbar.` |
| Spiegel (extrahiert), stabil | `:5101` → `:5102` | Exit 0; `Spiegel traegt denselben Index: sha256:a61a6b79…` — Index-Digest über zwei Registries erhalten |
| Spiegel, Prerelease (`IS_STABLE=false`) | nach `:5102/d-check-pre` | Tags `["v0.0.0-dev"]`, **kein** `latest` (DIST-002 Boundary) |
| Spiegel, Hub unerreichbar (F-3; `create` durch `true` ersetzt, sonst wörtlich) | `DOCKERHUB_IMAGE=localhost:5109/…` | Exit 1; `::error::Index-Digest nicht lesbar (ghcr.io=[sha256:a61a6b79…] docker.io=[<leer>]). Die Gleichheit ist damit UNGEPRUEFT…` — der Zweig wird erreicht, nicht mehr der generische ERR-Trap |
| Spiegel, ungleich | GHCR `v9.9.9` → `a61a6b79…`, Hub `v9.9.9` → anderer Index | Exit 1; `::error::Index verschieden: ghcr.io=[…a61a6b79…] docker.io=[…e7126996…]. DC-FA-DIST-002 sagt denselben Index-Digest zu. BEREITS VEROEFFENTLICHT: …` |

- `--provenance=false --sbom=false` im Build (`tools/image-publish.sh`); der Attestations-Bruch oben zeigt, warum.
- Label-Check: vor dem Push an `d-check:latest` und `d-check:arm64` (Daemon), nach dem Push je Plattform aus der Registry (`image-verify-published.sh` (2)).
- Neue Actions SHA-gepinnt mit Tag-Kommentar; `make workflow-pins` in `gates` grün. Die Builder-Images (R1 F-2) gegen Docker Hub gelesen: `moby/buildkit:v0.33.1` → `sha256:cec9f139…`, `tonistiigi/binfmt:qemu-v10.2.3` → `sha256:400a4873…` — beide identisch zu den Pins in `release.yml`.
- **Prerelease-Lauf `v0.84.0-rc.1`: offen, ausstehend** — kein Verifier-Befund. Er allein belegt (a) den grünen arm64-Test unter QEMU im Runner samt Laufzeit (Risiko 3), (b) die Erhaltung des Index-Digests **gegen Docker Hub** (Risiko 1, Re-Evaluierungs-Trigger ADR-0102), (c) das Verhalten der gepinnten Builder-Images in der Action. Der Haken an DoD 3 wartet darauf.

### 4. `make gates` grün; Review mit Report; Verifikation: **ERFÜLLT** (siehe V-1)

- `make gates` selbst gefahren, HEAD `22816afd`, sauberer Baum: Exit 0, Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; Einzelbelege u. a. `d-check: 967 Datei(en) geprüft, 0 Befund(e)`, `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`, `Ran 55 rules on 65 files: 0 findings.`
- `make trace-check RANGE=e5ed67ee..HEAD` → `5 Commit(s) — OK`, `0 Befund(e)`; `make review-coverage` → `0 Befund(e)`.
- R1-Report `docs/reviews/2026-10-07-slice-256-multi-arch-image-r1.md` (Verdikt „blockierend: ja", vier MEDIUM). **Einarbeitung in `22816afd` Stelle für Stelle geprüft** (keine R2):
  - **F-1 eingearbeitet und belegt:** `tools/image-publish.sh` pusht `push-by-digest=true`, prüft, taggt danach per `imagetools create`; ADR-0102 Entscheidung 3 und Option F nachgezogen. Bruch-Test oben: Gegenprobe rot ⇒ kein Tag.
  - **F-2 eingearbeitet:** beide Images digest-gepinnt (Pins gegen Hub gelesen, oben); ADR-0102 Entscheidung 8 und §Konsequenzen benennen die fehlende Frische-Achse; neues Risiko §6.
  - **F-3 eingearbeitet und belegt:** `index_digest()` kapselt `|| true` innerhalb der Pipeline; der `UNGEPRUEFT`-Zweig wird unter `set -euo pipefail` + ERR-Trap erreicht (Tabelle oben).
  - **F-4 eingearbeitet:** Grenze 3 in `harness/sensors/image-test.md` (QEMU auf beiden Seiten, Plattform erzwungen, Wahl auf echtem arm64 nicht belegt).
  - **F-5 eingearbeitet**, aber mit Rest, siehe V-2. **F-6 eingearbeitet:** der Kopf von `image-verify-published.sh` trägt nur noch `(ADR-0102)` als Rang-Zeiger.
  - **F-7 eingearbeitet und belegt** (Push-Schritt, Tabelle oben). **F-8 eingearbeitet** (Nicht-Linux-Zweig in `image-test.sh`; auf Linux nicht ausführbar, nur gelesen). **F-9 eingearbeitet** (`.PHONY` trägt `image-test-arm64 image-publish`).

### 5. Closure-Notiz, Register, Risiko-Ausgänge, Paarungen: **OFFEN, wie erwartet**, kein Befund

§6 steht auf `(offen)`, §7 trägt `—`; `make verify-closure-notes` greift erst in `done/`. Für die Closure belegt dieser Bericht: Risiko 4 („zweiter Bau zwischen Prüfung und Push … nur, wenn fail-closed") — die Gegenprobe **ist** fail-closed und lässt keinen Tag zurück (Bruch-Test). Risiken 1 und 3 hängen am Prerelease-Lauf; Risiko 2 an `slice-257` (liegt in `open/`); Risiko 5 bleibt ohne Sensor.

---

## Akzeptanzkriterien gegen die Belege

| Kriterium | Beleg | Stand |
|---|---|---|
| DIST-001 Happy (amd64) | `make image-test` (1) | belegt |
| DIST-001 Boundary `:ro` (amd64) | `make image-test` (2) | belegt |
| DIST-001 Negative (amd64) | `make image-test` (3) | belegt |
| DIST-001 Kriterien für arm64 | `make image-test-arm64` | lokal nur der fail-closed-Pfad; grün erst im Prerelease (QEMU) |
| DIST-001 Boundary (Plattform) | Index trägt genau `linux/amd64`/`linux/arm64` (`image-verify-published.sh` (1)), arm64-Variante vor dem Push geprüft | **Vorbedingung** belegt; dass ein echter arm64-Host selbst wählt, ist nicht gemessen — als Grenze 3 benannt |
| DIST-001 „vor Veröffentlichung geprüft" | Push ohne Tag, Gegenprobe, Tags danach | belegt inkl. Bruch (kein Tag) |
| DIST-002 Happy (Index-Digest gleich, aus beiden Registries) | Spiegel-Schritt gegen `registry:2` | belegt gegen `registry:2`; **gegen Docker Hub offen** (Prerelease) |
| DIST-002 Boundary (Prerelease ohne `latest`) | Spiegel-Schritt mit `IS_STABLE=false` | belegt gegen `registry:2` |
| DIST-002 Negative (Zugangsdaten) | Login-Schritt unverändert | nur gelesen; im Range nicht berührt |

---

## Abgrenzung (§1): gehalten

- Handbuch, `README.md`/`README.de.md`, `operations.md`, `packaging/dockerhub/overview.md` sind im Range unberührt (Release-Prep). `overview.md` sagt weiter „config digest … manifest digest is registry-local" — das ist die geplante Release-Prep-Nachzugsliste aus §3, kein Befund.
- `make build` bleibt Single-Platform (`docker build` ohne `--platform`).
- `tools/image-scan.sh`: nur ein Kommentar (Plan-Änderung), keine Scan-Logik — `slice-257` bleibt Adresse.
- Jede Datei im Diff steht in §3 oder in der Plan-Änderung. Rückführungs-Bedingungen §4 nicht eingetreten (die Prüfung vor dem Push kommt ohne eigene Zwischen-Registry aus; Docker Hub ist ungemessen).

---

## Verifier-Befunde

### V-1: LOW

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §6 Schritt 8 (Handoff an den Reviewer, kein Self-Review); DoD Punkt 4 („Review durchgeführt")
- **Pfad:** `tools/image-publish.sh` · „Build beider Plattformen und Push OHNE Tag (push-by-digest)"
- **Befund:** Die Veröffentlichungs-Mechanik selbst — `tools/image-publish.sh` (66 Zeilen) samt Umbau von `image-verify-published.sh` (Abruf je Manifest-Digest) und Push-Schritt — entstand erst in `22816afd` und liegt damit außerhalb der R1-Range `e5ed67ee..16ebc6f9`. Funktional ist sie hier bestätigt (grün, sieben Bruchpfade), aber kein Review hat sie auf Maintainability, Kommentar-Klassen und Hard Rules gelesen. Der R1-Verdikt war „blockierend"; die Einarbeitung ist der größte Code-Anteil des Slice.
- **Verifizierbar:** ja — eine R2 über `16ebc6f9..22816afd`.
- **Erwartete Aktion:** R2 über den Fix-Commit, oder in der Closure-Notiz ausdrücklich begründen, warum dieser Bericht sie ersetzt.
- **Klasse:** Fix-Commit größer als die Findings, die ihn auslösten

### V-2: LOW

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §5 Regeln 13 und 15 (Grenze gegen den Gegenstand; nicht mehr behaupten, als gemessen ist)
- **Pfad:** `tools/image-scan.sh` · „gleicher Index-Digest (ADR-0102), gleicher Inhalt"
- **Befund:** Der Absatz beginnt „GEMESSEN vor der Aufnahme … beide Refs melden denselben Befundsatz. Das ist zu erwarten -- gleicher Index-Digest …". Gemessen wurde an Single-Arch-Bildern, deren Manifest-Digests je Registry verschieden waren (der Grund für ADR-0065); die heute publizierten Tags (bis `v0.83.0`) sind weiter solche. Der neue Text schreibt der vergangenen Messung eine Eigenschaft zu, die damals nicht bestand und bis zum ersten Index-Release nicht besteht. R1 F-5 ist damit in der Richtung, nicht in der Sache geschlossen.
- **Verifizierbar:** nein (Kommentar).
- **Erwartete Aktion:** die Begründung zeitlos fassen („gleicher Inhalt; ab ADR-0102 derselbe Index-Digest"), oder den Satz an `slice-257` übergeben, der die Datei ohnehin anfasst.
- **Klasse:** Kommentar-Spiegel auf neuen Stand umgeschrieben, Messung alt

### V-3: INFO

- **Kategorie:** INFO
- **Quelle:** ADR-0102 §Konsequenzen
- **Pfad:** `docs/plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md` · „Handbuch, READMEs, `operations.md`, Hub-Overview und\n  `releasing.md` sagen Index-Digest"
- **Befund:** `releasing.md` ist in diesem Slice bereits nachgezogen (Pipeline-Schritte 2–8, Index-Digest); die Folgepflicht-Liste nennt die Datei weiterhin. Da die ADR noch `Proposed` ist, ist das vor `Accepted` noch korrigierbar.
- **Erwartete Aktion:** `releasing.md` aus der Folgepflicht nehmen, bevor die ADR `Accepted` wird.
- **Klasse:** Folgepflicht-Liste hinter dem Diff

### V-4: INFO

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Pfad:** `tools/image-publish.sh` · „tagged=\"$(docker buildx imagetools inspect"
- **Befund:** Schritt (3) bestätigt nur `v$VERSION` gegen den geprüften Digest, nicht `:latest`. Gemessen zeigen beide auf denselben Digest, und ein scheiterndes `imagetools create` bricht das Skript ohnehin; eine Nachprüfung von `:latest` fehlt nur für den Fall eines teilweise wirksamen `create`.
- **Erwartete Aktion:** keine nötig; optional dieselbe Prüfung für `latest` bei `PUBLISH_LATEST=true`.
- **Klasse:** —

### V-5: INFO

- **Kategorie:** INFO
- **Quelle:** `AGENTS.md` §6 Schritt 4 („gehört vor den Code, nicht in den Bericht danach")
- **Pfad:** `docs/plan/planning/in-progress/slice-256-multi-arch-image.md` · „*(Plan-Änderung im Lauf:"
- **Befund:** Die Plan-Änderung ist vollständig und vom Diff gedeckt, wurde aber erst in `22816afd` geschrieben — `image-verify-published.sh`, die awk-Änderung und der `dependabot.yml`-Kommentar kamen schon mit `91e4e650`. Dieselbe Klasse stand bei `slice-255` (dort V-3); zweites Auftreten in Folge — Kandidat für das Beobachtungs-Register bei der Closure.
- **Erwartete Aktion:** keine am Code; Rollen-Verweis Planner (Closure-Notiz, Register).
- **Klasse:** Plan-Änderung im Code-Commit

Negativ geprüft ohne Befund:
- `Dockerfile` Cross-Compile gegen Byte-Gleichheit über zwei Builder (beide Plattformen)
- ELF-Prüfung in `image-test.sh` samt Bruch-Test
- fail-closed von `image-test-arm64` ohne binfmt
- `image-publish.sh`: Push ohne Tag, Gegenprobe, Tags danach; Reproduzierbarkeit des Index-Digests über zwei Läufe
- `image-verify-published.sh`: fünf Bruchpfade (Binary, Einzel-Manifest, Attestation, Label-Version, unerreichbar)
- Push-Schritt: Pin-Extraktion und Form-Prüfung
- Spiegel-Schritt: Erhaltung des Index-Digests über zwei Registries, Prerelease ohne `latest`, `UNGEPRUEFT`-Zweig, Ungleichheit
- Builder-Image-Pins gegen die Registry
- `image-digest-axis` alte/neue Form
- Gate-Index (`image-test-arm64` als Gate, `image-publish` als Werkzeug), `gate-consistency`, `workflow-pins`
- ADR-0065 unverändert (`adr-check`), Traceability (`trace-check`)
- Abgrenzung (Release-Prep-Dateien, `make build`, `image-scan.sh`-Logik)

---

## Verdict

- **DoD 1, 2 erfüllt.**
- **DoD 3 erfüllt bis auf den Prerelease-Lauf `v0.84.0-rc.1` — offen, ausstehend** (kein Befund; er trägt Docker-Hub-Digest, arm64 grün unter QEMU, Laufzeit).
- **DoD 4 erfüllt**; V-1 empfiehlt eine R2 über den Fix-Commit.
- **DoD 5 erwartungsgemäß offen.**

`make gates` selbst gefahren und grün. Die Zusage „unter keinem Tag liegt eine ungeprüfte Variante" ist lokal im Gut- und im Bruchfall belegt, ebenso die Erhaltung des Index-Digests beim Spiegeln über zwei `registry:2`. Alle neun R1-Befunde sind eingearbeitet, F-1, F-3 und F-7 mit Bruch-Beleg. Kein HIGH, kein MEDIUM. V-1 und V-2 bitte annehmen oder begründen; V-3 vor `Accepted` der ADR mitnehmen; V-4, V-5 ohne Pflicht-Aktion.
