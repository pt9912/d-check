# Slice slice-048: Modul `versions` — Versions-Pin-Konsistenz (Idee 1)

**Status:** done (abgeschlossen, welle-37-versions).

**Welle:** welle-37-versions (Trigger: Auftraggeber-Idee — „statt nacktem
Versionsstring ein Markdownlink auf seine Definition", weiterentwickelt zu
„nicht vergessen, die Versionsnummer beim Release anzupassen").

**Bezug:** Idee 1 dieser Session. Führt eine **neue VER-Anforderung** im
Lastenheft ein (Modul `versions`, Versions-Pin-Konsistenz) plus einen
**begleitenden ADR** (Fence-Öffnung für Versions-Pins, Mechanik-Präzedenz
[ADR-0018](../../adr/0018-diagram-fence-ausnahme.md) für das opt-in Modul
`diagrams`). Verteilung wie der Rest des Werkzeugs, nicht als kopiertes Skript
([`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Linie),
konsumiert über das bestehende `doc-check`-Target
([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)-`--print-mk`-Fragment)
— **kein** neues `--print-mk`-Target.

**Autor:** pt9912. **Datum:** 2026-06-24.

---

## 1. Ziel

Versionsnummern, die über die Doku verstreut hart gepinnt sind (vor allem
`ghcr.io/…:vX.Y.Z` in Kommando-Beispielen, ~18× im Benutzerhandbuch + README),
sollen beim Release nicht mehr still veralten. Zwei sich ergänzende Mechaniken:

1. **Doku-Boden (dieser Slice, Schritt A — erledigt):** ein Release-Register
   [`version.md`](../../../../version.md#aktuell) als kanonischer, auflösender
   Ziel-Ort. **Nur die aktuelle Version trägt einen HTML-Anker** — dadurch
   *bricht* jeder Markdown-Link-Pin auf eine veraltete Version beim Release
   (Anker wandert) und wird vom **bestehenden** `anchors`-Gate gefangen
   (Anker-Kaskade, kein neues Feature nötig). Der Handbuch-Header-Pin ist
   bereits als Markdown-Link auf `version.md#<aktuelle>` ausgeführt.
2. **Werkzeug-Regel (Schritt C — erledigt, `43f9ed0`):** das opt-in Modul
   `versions` prüft alle Versions-Pins gegen die deklarierte aktuelle Version
   (Befund bei Abweichung) — auch in Fenced-Code (die ~18 `docker run`-Beispiele),
   die die Vorverarbeitung sonst entfernt. In `.d-check.yml` aktiviert, Dogfooding
   live.

## 2. Entscheidungen

- **d-check-Feature statt kopiertem Shell-Skript.** Ein `tools/*.sh` würde über
  die Schwester-Repos per Datei-Kopie driften — genau die
  [`MR-003`](../../../../harness/conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh)→[`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Lektion
  (vendorter Sensor → durch Dogfooding ersetzt). Die Prüf-Logik lebt im
  gepinnten Image, Config in `.d-check.yml`; konsumiert über das bestehende
  `doc-check`-Target — **kein** neues `--print-mk`-Target (das Modul läuft als
  Teil von `doc-check`, sobald in `.d-check.yml` aktiviert; weitere
  Fragment-Targets sind bewusst ausgegrenzt).
- **Generisch, nicht d-check-spezifisch.** „Alle Versions-Pins == aktuelle
  Version" ist das Problem jedes Repos mit versioniertem Artefakt → eine
  generische Regel, kein Repo-Skript.
- **Wahrheitsquelle konfigurierbar:** Default `version.md#aktuell` (in-tree,
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)-sicher);
  der Git-Tag scheidet aus (außerhalb des read-only gemounteten Baums).
- **Fence-Öffnung** für die Pin-Muster — bewusste, eng begrenzte Ausnahme von
  der Fence-Opazität; Präzedenz
  [ADR-0018](../../adr/0018-diagram-fence-ausnahme.md) (öffnet gelistete
  Diagramm-Fences). Breiter als `diagrams` → eigener ADR.
- **Eigenes Modul `versions`** (nicht mit dem geplanten content-drift-Modul
  `pins` verschmolzen): andere Mechanik (Wert-Gleichheit vs. Span-Hash), je
  einzeln opt-in und testbar.

## 3. Definition of Done

- [x] **Schritt A (Doku-Boden):** [`version.md`](../../../../version.md#aktuell)
  (Register, only-current-anchor, `## Aktuell` + `## Verlauf` der 22 Tags) +
  Handbuch-Header-Pin als Markdown-Link; `make gates` grün (Anker-Kaskade aktiv
  via bestehendes `anchors`-Gate).
- [x] **Schritt C (Spec):**
  [`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
  (Bereichskürzel `VER` in §3, Versions-Bump 0.28.0 + §7-Historie) +
  [ADR-0019](../../adr/0019-versions-pin-fence-ausnahme.md) (Fence-Öffnung,
  Status Proposed) + ADR-Index; doc-first vor Code.
- [x] **Schritt C (Code):** Modul `versions` (`rules/versions.go`: Fence-offener
  `pin-pattern`-Scan, `current-from`-Auflösung via Heading-Section/HTML-Anker/
  ganze Datei, `version-stale`, Ventile `exempt-paths`/`d-check:ignore`,
  current-from-Datei selbst-ausgenommen, fail-closed), Config-Adapter +
  rules-Registry verdrahtet, `.d-check.yml`-Selbstkonfiguration (Konsum über das
  bestehende `doc-check`-Target, **kein** neues `--print-mk`-Target); 10 Tests
  (Happy/Negative-in-Fence/Ventile/Modul-aus + current-from-Fälle). `make gates`
  grün, Dogfooding live (~18 ghcr-Pins gegen `version.md#aktuell`).
- [ ] `make gates` grün; unabhängiges Review; Closure (Move nach `done/` +
  Roadmap-Flip, [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

## 4. Risiken / offene Punkte

- **Fence-Öffnung breiter als `diagrams`** (Muster-Scan über Fences statt
  gelistete Diagramm-Fences) — Begründung gehört in den ADR; opt-in default-off
  hält die [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)-Abwärtskompatibilität
  (ohne `versions`-Block byte-identisch).
- **Release-Prep v0.28.0:** das aktive `versions`-Gate bewacht jetzt den eigenen
  Release — beim Bump müssen `version.md#aktuell` **und** alle ~18 ghcr-Pins
  gemeinsam gezogen werden, sonst `version-stale`. (Die Markdown-Link-Pins sichert
  zusätzlich die Anker-Kaskade.)
- **Historische Pins** (z. B. `:v0.1.0` in `done/`-Slices) dürfen nicht
  mitgebumpt werden → exempt-paths (wie `ids` `done/`/CHANGELOG ausnimmt).

## 5. Trigger

Auftraggeber-Idee (2026-06-24), nach Spike/Design-Diskussion bestätigt:
Reihenfolge A → C (dieser Slice) → D (Modul `pins`, eigener Slice).

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code + Spec; „Doc führt, Code folgt"). `version.md` ist abgeleitete
Doku (Register), kein Vertrag. Keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** Idee 1 in vier Schritten: **(A)** Doku-Boden — Release-Register
[`version.md`](../../../../version.md#aktuell) (only-current-anchor: nur die
aktuelle Version trägt einen `<a id>`-Anker, der beim Release wandert → veraltete
Markdown-Link-Pins brechen via bestehendem `anchors`-Gate) + Handbuch-Header-Pin
als Link. **(C-Spec)**
[`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
(Lastenheft 0.28.0, Bereich `VER`) + [ADR-0019](../../adr/0019-versions-pin-fence-ausnahme.md)
(Fence-Öffnung, breiter als [ADR-0018](../../adr/0018-diagram-fence-ausnahme.md)) +
spezifikation `.a`/Schema/Grund-Code.
**(C-Code)** Modul `versions` (`rules/versions.go`: Fence-offener Pin-Scan,
`current-from`-Auflösung, `version-stale`, Ventile, fail-closed) + Config-Adapter/
run-Wiring + `.d-check.yml`-Dogfooding + 10 Tests. **(Release)** v0.28.0.
Meta-Gate-Shell-Skript bewusst verworfen (Copy-Drift über die Repo-Familie,
[`MR-003`](../../../../harness/conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh)→[`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding));
Verteilung im gepinnten Image, Konsum über `doc-check`.

**Belege.**
- `make gates` **grün** (doc-check inkl. live `versions`-Dogfooding, lint, test,
  arch-check, Coverage 93,90 %, semgrep, gate-consistency, planning-check).
- Plan-/Fundament-Review **R1→R2→R3 ACCEPT** (4→2→2-LOW Befunde, alle behoben) +
  **Impl-Review** (4 Befunde: `--print-config`-Template, Adapter-Tests,
  Lastenheft-Historie-Exempt, Slice-Doc).
- Dogfooding live: ~18 `ghcr`-Image-Pins gegen `version.md#aktuell` gateguarded;
  die Release-Prep v0.28.0 zog `version.md` + alle Pins gemeinsam — das Gate
  erzwang es (eat-your-own-dogfood).
- Release **v0.28.0** auf GHCR (Run `28095582612` grün in 2m21s, Tags
  `v0.28.0`+`latest`), Digest-Pin
  `ghcr.io/pt9912/d-check@sha256:0bb84b529d3a65bdf9e849dd79cb8e9011bc388ecf9bffc5930f6c96bcc0cba8`.

**Lerneintrag.** Das `versions`-Gate bewacht ab sofort den eigenen Release: Bump =
`version.md#aktuell` + `<a id>`-Anker verschieben + **alle** `ghcr`-Pins gemeinsam,
sonst `version-stale` (in [`docs/user/releasing.md`](../../../user/releasing.md)
§Release-Prep dokumentiert). Generische Regel (jedes Repo mit versioniertem
Artefakt), kein Repo-Skript — der entscheidbare Kern ist Wert-Gleichheit, nicht
Semantik.
