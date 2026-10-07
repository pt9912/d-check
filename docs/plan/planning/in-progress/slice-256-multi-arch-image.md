# slice-256: Multi-Arch-Image — `linux/amd64` und `linux/arm64` unter einem Index

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[`DC-FA-DIST-002`](../../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel);
[ADR-0002](../../adr/0002-distribution-ghcr-image.md),
[ADR-0011](../../adr/0011-digest-pins-build-gate-images.md) §3 (die Basis-Pins
sind bereits Index-Digests mit `arm64`),
ADR-0065 <!-- d-check:ignore (ADR-0065 ist abgelöst, ein Link wäre matrix-inactive) -->
(Prüfgröße des Spiegels — wird abgelöst); Auftraggeber-Anfrage 2026-10-07
(ein Image, das `docker pull` auf Apple Silicon und arm64-Linux ohne
Emulation wählt).

**Berührte Spec-Stellen:** [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[`DC-FA-DIST-002`](../../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel),
[`SPEC-066`](../../../../spec/spezifikation.md#6-externe-verträge) (Runtime-Basis-Zeile).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Jedes Release veröffentlicht unter seinen Tags einen **OCI-Index**
mit den Plattformen `linux/amd64` und `linux/arm64`; `docker pull` wählt die
passende Variante selbst (macOS auf Apple Silicon über Docker Desktop,
arm64-Linux, Graviton). Beide Varianten sind vor dem Push **nativ-identisch
geprüft** — die Zusage aus DIST-001 gilt je Plattform, nicht nur für die des
Release-Runners. Der Spiegel nach Docker Hub ist eine **Kopie des Index**,
und seine Prüfgröße wird der **Index-Digest** — schärfer als der heutige
Config-Digest, weil er alle Plattformen zugleich bindet.

**Vorab gemessen** (Spike, 2026-10-07, lokal gegen zwei Wegwerf-`registry:2`,
kein Push nach GHCR/Hub):

- Cross-Compile (`FROM --platform=$BUILDPLATFORM` für die Go-Stages,
  `GOOS`/`GOARCH` aus `TARGETOS`/`TARGETARCH`) baut beide Plattformen in 51 s
  **ohne Emulation**; die Runtime-Stage löst ihre Plattform über den
  vorhandenen Index-Pin auf.
- Das `amd64`-Binary des Cross-Builds ist **byte-gleich** zum Binary aus
  `make build` (sha256 `c46c1576…` beide) — das von den Gates geprüfte
  Binary ist das ausgelieferte.
- `docker buildx imagetools create -t <B> <A>` erhält den **Index-Digest
  byte-gleich** über zwei Registries (`sha256:0b369a06…` beidseitig) —
  anders als der heutige `docker tag`/`push`-Spiegel, der neu komprimiert
  (der Grund für ADR-0065). **Gegen Docker Hub ungemessen** — das misst der <!-- d-check:ignore (ADR-0065 ist abgelöst, ein Link wäre matrix-inactive) -->
  Prerelease-Lauf (§2), bevor die Zusage im Lastenheft `Accepted` wird.
- Labels sind je Plattform über `imagetools inspect --format '{{json .Image}}'`
  lesbar; `--provenance=false` hält den Index frei von Attestations-Einträgen
  (`unknown/unknown`), an denen die Plattform-Prüfung sonst hinge.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **CVE-Scan je Plattform** (`tools/image-scan.sh` scannt heute die
  Host-Variante des Index, also nur `amd64`) — Folge-Slice
  slice-257; bis dahin ist `arm64` im
  Nachtlauf blind, als Risiko in §6 geführt.
- **Weitere Plattformen** (`linux/arm/v7`, `s390x`, Windows) — keine Anfrage;
  jede weitere ist eine Erweiterung der Zusage, kein Teil dieser.
- **Native Release-Binaries** (`darwin/arm64` ohne Docker) — DIST-001
  Out-of-Scope, eigene ADR-Frage.
- **Lokaler `make build` bleibt Single-Platform** — die Gates (`doc-check`,
  `trace-check`, …) brauchen ein geladenes Image der Host-Plattform; ein
  Multi-Platform-`--load` setzt den containerd-Image-Store voraus. Gebaut
  wird der Index nur im Release-Pfad.
- **Handbuch, READMEs, `operations.md`, Hub-Overview** — Release-Prep
  (`AGENTS.md` §5 Regel 17), nicht im Feature-Commit.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] Lastenheft: DIST-001 nennt die Plattformen und die Nativ-Identität je
      Plattform; DIST-002 wechselt die Prüfgröße auf den Index-Digest
      (Akzeptanzkriterien neu, Out-of-Scope „andere Plattform-Matrix"
      bleibt als Zusage der Inhalts-Gleichheit); Version und Historie.
      Neue ADR (Multi-Arch-Index, Spiegel per `imagetools create`, löst
      ADR-0065 ab) samt Index; [`SPEC-066`](../../../../spec/spezifikation.md#6-externe-verträge) nachgezogen. <!-- d-check:ignore (ADR-0065 ist abgelöst, ein Link wäre matrix-inactive) -->
- [x] `Dockerfile` cross-kompiliert; `tools/image-test.sh` prüft eine
      benannte Plattform (Binary aus dem Image der Plattform, Ausführung über
      binfmt/QEMU, wo sie nicht die des Hosts ist); ein `make`-Target baut und
      prüft die Nicht-Host-Plattform — im Index geführt.
- [x] `release.yml`: buildx-Index-Build mit `--provenance=false`, beide
      Plattformen **vor** dem Push geprüft, Label-Check je Plattform,
      Konsumenten-Pin = Index-Digest, Spiegel per `imagetools create`,
      Gleichheitsprüfung über den Index-Digest (fail-closed wie bisher); neue
      Actions SHA-gepinnt (`make workflow-pins`). **Prerelease-Lauf**
      (`v0.84.0-rc.1`) grün, mit gemessenem Index-Digest auf beiden
      Registries.
- [x] `make gates` grün; Review durchgeführt, Report unter `docs/reviews/`;
      Verifikation.
- [x] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | DIST-001/002, Version, Historie |
| `spec/spezifikation.md` | update | [`SPEC-066`](../../../../spec/spezifikation.md#6-externe-verträge) (Plattformen; Altdrift `latest` mit) |
| `docs/plan/adr/0102-…md` + `README.md` | neu / update | Entscheidung + Index |
| `Dockerfile` | update | `--platform=$BUILDPLATFORM`, `TARGETOS`/`TARGETARCH` |
| `tools/image-test.sh` | update | Plattform-Parameter; die amd64-Annahme im Kopf entfällt |
| `Makefile` | update | Plattform-Variable für `image-test`, Target für die Nicht-Host-Plattform |
| `harness/README.md`, `harness/sensors/image-test.md` | update | Index-Zeile, Grenze |
| `.github/workflows/release.yml` | update | Index-Build, Prüfung vor Push, Spiegel, Digest |
| `docs/user/releasing.md` | update | Pipeline-Schritte (Betriebs-Doku des Release-Pfads, kein Release-Prep-Gegenstand) |
| `tools/image-publish.sh`, `tools/image-verify-published.sh` | neu | Push ohne Tag, Gegenprobe, dann Tags |
| `.github/dependabot.yml`, `tools/image-scan.sh` | update | Kommentare, die die alte FROM-Form bzw. den Config-Digest nannten |

*(Plan-Änderung im Lauf: die beiden Skripte, die Kommentar-Spiegel in `dependabot.yml` und `tools/image-scan.sh` sowie die awk-Extraktion von `image-digest-axis` im `Makefile` — die FROM-Zeile mit `--platform` hätte `go-base-digest` sonst still auf SKIP gesetzt; nach R1 die Reihenfolge Push ohne Tag → Gegenprobe → Tags und die Digest-Pins der Builder-Images; nach dem Prerelease-Lauf auf Auftraggeber-Wunsch der Abschnitt *Vorabversion* in `releasing.md` — die Datei unterschied Prerelease und stabiles Release nur bei `:latest`.)*

**Reihenfolge im Release-Pfad** (Entwurf, die ADR legt fest): `make ci`
(Gates + `amd64`-image-test wie heute) → QEMU/buildx einrichten → je
Plattform `--load` aus demselben Builder und `image-test` dagegen → Index
`--push` nach GHCR (Cache-Treffer, dieselben Layer) → Gegenprobe: das Binary
jeder Plattform im **gepushten** Index ist byte-gleich zum geprüften →
Label-Check je Plattform → Spiegel `imagetools create` → Index-Digest beider
Registries gleich → GitHub-Release mit Index-Digest.

**Spiegel vor dem Editieren** ([`MR-025`](../../../../harness/conventions.md#mr-025)) —
die Aussage „Prüfgröße Config-Digest" steht außer in DIST-002 und ADR-0065 in: <!-- d-check:ignore (ADR-0065 ist abgelöst, ein Link wäre matrix-inactive) -->
`release.yml` (Kommentare, `config_digest()`), `docs/user/releasing.md`,
`docs/user/benutzerhandbuch.md` §Docker-Image, `docs/user/operations.md`,
`README.md`/`README.de.md`, `packaging/dockerhub/overview.md`; die
amd64-Annahme in `tools/image-test.sh` (Kopf) und `docs/user/releasing.md`
§`image-test` auf macOS. Die Nutzer-Doku zieht die Release-Prep nach.

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt; `in-progress/` leer
(slice-255 geschlossen).

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): die Prüfung vor dem Push lässt sich im
  Release-Pfad nicht ohne eigene Registry-Zwischenstufe bauen — dann wird
  „Index bauen und spiegeln" von „arm64 vor dem Push prüfen" getrennt.
- `in-progress` → `open` (blockiert): der Prerelease-Lauf zeigt, dass Docker
  Hub den Index-Digest beim `imagetools create` **nicht** erhält — dann trägt
  die Lastenheft-Zusage nicht und braucht eine andere Prüfgröße (je
  Plattform der Config-Digest); Rückfrage an den Auftraggeber.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Das stabile Release `v0.84.0` ist Release-Prep, nicht
Closure-Bedingung.

## 6. Risiken und offene Punkte

- **Docker Hub erhält den Index-Digest nicht** (gemessen nur gegen
  `registry:2`). — **Ausgang:** entfallen — gemessen in zwei echten Läufen:
  `v0.84.0-rc.1` (`sha256:bd064c9f…`) und `v0.84.0` (`sha256:e82ef2d2…`)
  tragen auf GHCR und Docker Hub denselben Index-Digest, `:latest` ebenso.
- **`arm64` im Nachtlauf-CVE-Scan blind**, bis slice-257 schließt. —
  **Ausgang:** eingetreten — Folge-Slice slice-257 (CVE-Scan je Plattform).
- **QEMU-Laufzeit im Release-Job** — der `arm64`-image-test läuft emuliert;
  der Job hat 30 min. Gemessen wird im Prerelease-Lauf. — **Ausgang:**
  entfallen — der ganze Release-Job brauchte 4,5 min (`v0.84.0-rc.1`), der
  arm64-Test ist darin enthalten.
- **Ein zweiter Bau zwischen Prüfung und Push** — `--load` je Plattform und
  `--push` des Index sind getrennte buildx-Aufrufe; die Gegenprobe auf das
  gepushte Binary schließt die Lücke nur, wenn sie fail-closed ist. —
  **Ausgang:** entfallen — die Gegenprobe ist fail-closed und läuft vor jedem
  Tag; lokal gebrochen (anderes Binary ⇒ rot, kein Tag), in beiden echten
  Läufen grün mit gleichen Binaries je Plattform.
- **Builder-Images ohne Frische-Achse** — `moby/buildkit` und
  `tonistiigi/binfmt` sind in `release.yml` digest-gepinnt, aber weder
  Dependabot noch der Nachtlauf melden einen neueren Stand. —
  **Ausgang:** weiter offen — [`BEO-ALL/digest-pin-ohne-frische-achse`](../observations/BEO-ALL/digest-pin-ohne-frische-achse/state.md).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Messung vor dem Plan trug den ganzen Slice:
  ein Spike gegen zwei Wegwerf-Registries zeigte vor jeder Spec-Änderung, dass
  `imagetools create` den Index-Digest erhält und das Cross-Compile-Binary
  byte-gleich zu dem der Gates ist — die schärfere Zusage in DIST-002 stand
  damit auf einem Befund, nicht auf einer Annahme, und der echte Lauf
  bestätigte sie gegen Docker Hub. Jeder Fehlerpfad des Release-Pfads ist
  lokal gebrochen worden (anderes Binary, Einzel-Manifest, Attestation,
  Label-Version, unerreichbare Referenz, scheiternder Push).
- **Was ging anders als geplant:** Der erste Entwurf prüfte das gepushte
  Binary **nach** dem Taggen — der Review fand, dass `:latest` dann schon auf
  einen ungeprüften Index zeigen kann (R1-F-1); jetzt wird ohne Tag gepusht,
  geprüft und erst dann getaggt. Die Gegenprobe zog anfangs beide Plattformen
  über dieselbe Digest-Referenz, die der Daemon nur an ein Bild bindet — erst
  der Lauf gegen die Wegwerf-Registry zeigte es. Die Images, die die neuen
  Actions ziehen, waren ungepinnt (R1-F-2). Die FROM-Zeile mit `--platform`
  hätte `go-base-digest` still auf SKIP gesetzt; gefunden beim Auflisten der
  Spiegel, vor dem Commit. `releasing.md` kannte keine Vorabversion — der RC
  folgte der Datei wörtlich und zog Pins und Prosa auf sich; der Abschnitt
  *Vorabversion* ist auf Auftraggeber-Wunsch mitgenommen. Der Statuswechsel von
  ADR-0065 auf `Superseded` machte jeden Link auf sie zu `matrix-inactive`; <!-- d-check:ignore (ADR-0065 ist abgelöst, ein Link wäre matrix-inactive) -->
  lebende Verweise zogen um, die immutable [ADR-0068](../../adr/0068-lokale-workflow-referenzen-ohne-pin.md) kam wie [ADR-0047](../../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md) in
  `matrix.exempt-paths` (getragen von [ADR-0102](../../adr/0102-multi-arch-index-und-spiegel-per-index-digest.md)).
- **Steering-Loop-Eintrag:** keine neue Verkörperung —
  [`BEO-ALL/commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/state.md)
  (verkörpert) bekommt einen Beleg; neu registriert sind
  [`BEO-ALL/plan-aenderung-erst-im-fix-commit`](../observations/BEO-ALL/plan-aenderung-erst-im-fix-commit/state.md)
  (2×, mit dem Beleg aus slice-255) und
  [`BEO-ALL/digest-pin-ohne-frische-achse`](../observations/BEO-ALL/digest-pin-ohne-frische-achse/state.md)
  (1×).
- **Beobachtungs-Register (`../observations/`):** Evidence `slice-256` unter
  den drei genannten Einträgen, dazu `slice-255` unter dem neuen
  Plan-Änderungs-Eintrag.
- **Folge-Slices:** slice-257 (CVE-Scan je Plattform).
- **Risiken aus §6:** je ein Ausgang oben. Trigger-Audit: kein Carveout, kein
  bootstrap-aware Gate, keine Hard Rule mit eingetretenem Trigger;
  [ADR-0102](../../adr/0102-multi-arch-index-und-spiegel-per-index-digest.md)
  ist `Accepted` (ihr erster Re-Evaluierungs-Trigger — Docker Hub erhält den
  Index-Digest nicht — ist durch die zwei Läufe ausgeschlossen),
  ADR-0065 <!-- d-check:ignore (ADR-0065 ist abgelöst, ein Link wäre matrix-inactive) -->
  `Superseded by ADR-0102`. Nachtlauf-Stand <!-- d-check:ignore (Status-Wert, kein Verweis) -->
  ([`MR-053`](../../../../harness/conventions.md#mr-053)): wie in §8.
- **Drei Paarungen:** (a) Anker — kein `liegt in`-Feld; (b) Folge-Slice —
  slice-257 liegt in `open/`; (c) Register — die drei zitierten
  Beobachtungen existieren und tragen Belege.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** eine berührte Sub-Area: die
Distribution des Repos (Build-Rezept, Release-Pfad) unter dem Default `*`
(`ALL`); deklariert. `tools/harness/` (`HARN`) ist nicht berührt.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:**
[`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/state.md)
— der Digest-Pin in Handbuch §2 wird ein Index-Digest; sein Nachzug bleibt
ein ungewächterter Spiegel der Release-Prep, dieser Slice ändert daran
nichts. Keine weiteren Treffer für Image, Release oder Spiegel.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-07 — `upstream-drift` und `image-scan` grün.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
