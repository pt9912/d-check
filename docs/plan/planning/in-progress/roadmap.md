# Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-06-24.

**Format-Regel:** Die Roadmap ist eine Reihenfolge von **Wellen**,
keine Reihenfolge von Terminen. Termine erscheinen — falls überhaupt —
als Konsequenz der Wellen-Schätzung, nicht als Treiber.

---

## Aktuelle Welle

**welle-38-pins** — [`slice-049`](slice-049-pins-modul.md):
Modul `pins` (Idee 2, content-drift): ein optionaler Content-Pin (Hash des
Ziel-Spans) erkennt, dass sich der **Inhalt** eines verlinkten Abschnitts seit
dem Verlinken geändert hat (Befund `link-stale`) — über `target-missing`/
`anchor-missing` hinaus. Spike erledigt (Kalibrier-Grundlage); doc-first als neue
PIN-Anforderung + begleitender ADR (Fence-Öffnung). Zuletzt abgeschlossen:
welle-37-versions ([`slice-048`](../done/slice-048-versions-modul.md) — opt-in
Modul `versions` (10. Modul): gepinnte `ghcr`-Image-Verweise müssen die aktuelle
Version aus `version.md#aktuell` tragen, sonst `version-stale`; liest Pins auch
in Fenced-Code (gescopte Ausnahme), Ventile `exempt-paths`/`d-check:ignore`,
fail-closed, diagnose-only
([`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[ADR-0019](../../adr/0019-versions-pin-fence-ausnahme.md), Lastenheft 0.28.0);
Plan-Review R1→R3 ACCEPT + Impl-Review (4 Befunde behoben); `make gates` grün
(Coverage 93,90 %); Release **v0.28.0** auf GHCR, Digest-Pin
`ghcr.io/pt9912/d-check@sha256:0bb84b529d3a65bdf9e849dd79cb8e9011bc388ecf9bffc5930f6c96bcc0cba8`;
[Closure](../done/slice-048-versions-modul.md#7-closure-notiz-nach-done)). Davor
welle-36-print-mk-erweiterung
([`slice-047`](../done/slice-047-print-mk-doctor-repair-help-digest.md) — CR an
`--print-mk`: das Fragment bekommt `doc-doctor`/`doc-repair`/`doc-help` plus
`DCHECK_DIGEST` (Digest-Override per `ifeq`); alle Targets `##`-annotiert
([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben),
Lastenheft 0.27.0); kein ADR; R1 ACCEPT;
[Closure](../done/slice-047-print-mk-doctor-repair-help-digest.md#7-closure-notiz-nach-done)).
Davor welle-35-suggest-opt-in-hinweis
([`slice-046`](../done/slice-046-suggest-config-opt-in-hinweis.md) — Schärfung der
`--suggest-config ai-harness[-init]`-Ausgabe: sie nennt die nicht aktivierten
situativen opt-in-Module (`external`/`spans`/`hostpaths`/`diagrams`) in einem
Kommentar mit Verweis auf `--print-config`
([`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten),
Lastenheft 0.26.0); kein ADR; R1 ACCEPT;
[Closure](../done/slice-046-suggest-config-opt-in-hinweis.md#7-closure-notiz-nach-done)).
Davor welle-34-diagram-ids
([`slice-045`](../done/slice-045-diagram-ids.md) — opt-in Modul `diagrams` öffnet
gezielt benannte Diagramm-Fences (Default `mermaid`) und prüft die darin
gefundenen Kennungen auf Existenz in ihrer `defined-in`-Quelle (Befund
`diagram-id-undefined`); Existenz statt Link-Policy, reine Token-Extraktion ohne
Mermaid-Parser ([`DC-FA-DIAG-001`](../../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in),
[ADR-0018](../../adr/0018-diagram-fence-ausnahme.md), Lastenheft 0.25.0);
Fundament R1→R2 ACCEPT + Impl-R1 ACCEPT;
[Closure](../done/slice-045-diagram-ids.md#7-closure-notiz-nach-done)).
Davor welle-33-print-mk-trace
([`slice-044`](../done/slice-044-print-mk-trace-targets.md) — opt-in
`--trace --require-complete` bindet Requirements-Waisen an Exit 1
([`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code));
das `--print-mk`-Fragment trägt zusätzlich `doc-trace`/`doc-complete` plus
`TRACE_FLAGS` ([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben),
Lastenheft 0.24.0); kein ADR; unabhängiges R1 ACCEPT + R2;
[Closure](../done/slice-044-print-mk-trace-targets.md#7-closure-notiz-nach-done)).
Davor welle-32-codepaths-exempt
([`slice-043`](../done/slice-043-codepaths-exempt-paths.md) — Modul `codepaths`
bekam das `exempt-paths`-Ventil (Datei-Glob ohne codepath-Prüfung) wie `ids`;
Change Request [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(Lastenheft 0.23.0), kein ADR; unabhängiges R1 ACCEPT;
[Closure](../done/slice-043-codepaths-exempt-paths.md#7-closure-notiz-nach-done)).
Davor welle-31-requirements-completeness
([`slice-042`](../done/slice-042-requirements-completeness-gate.md) —
Requirements-Completeness-Gate `make completeness-check` failt bei
Requirements-Waisen (`d-check --trace --json` → `orphans > 0`), Closure-
Bindepunkt an `make fullbuild`, bewusst nicht in `make gates`/`ci`
([ADR-0017](../../adr/0017-requirements-completeness-gate.md));
[Closure](../done/slice-042-requirements-completeness-gate.md#7-closure-notiz-nach-done)).
Davor welle-30-adr-immutable
([`slice-041`](../done/slice-041-adr-immutable-gate.md) — ADR-Immutable-Gate
`make adr-check` erzwingt [`AGENTS.md` §3.5](../../../../AGENTS.md#35-adrs-sind-nach-accepted-immutable)
(Accepted-ADRs nur `## Geschichte`-Anhang + Status-Übergang; CI-Range +
pre-commit-Hook; [ADR-0016](../../adr/0016-adr-immutable-gate.md));
[Closure](../done/slice-041-adr-immutable-gate.md#7-closure-notiz-nach-done)).
Davor welle-29-planning-consistency
([`slice-040`](../done/slice-040-planning-consistency-gate.md) — Meta-Gate
`make planning-check` erzwingt Roadmap §Aktuelle Welle ↔
`in-progress/slice-*` (beide Richtungen, fail-closed); kein ADR;
[Closure](../done/slice-040-planning-consistency-gate.md#7-closure-notiz-nach-done)).
Davor welle-28-print-mk (`slice-038` — read-only-Generator `--print-mk` gibt ein
include-bares `d-check.mk` (version-gepinntes Image + `doc-check`-Target) aus;
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
0.22.0, Review R1+R2, kein ADR;
[Closure](../done/slice-038-print-mk.md#8-closure-notiz-nach-done)),
welle-27-rtm-trace (`slice-036`), welle-26-suggest-prefix (`slice-037`),
welle-25-pr-ci-traceability (`slice-039`), welle-23-latest-tag (`slice-034`)
und welle-24-kern-paketschnitt (`slice-035`).
Letztes Release **v0.28.0** auf GHCR (2026-06-24, slice-048: Modul `versions`),
Digest-Pin
`ghcr.io/pt9912/d-check@sha256:0bb84b529d3a65bdf9e849dd79cb8e9011bc388ecf9bffc5930f6c96bcc0cba8`;
davor v0.27.0 (slice-047) und v0.26.0 (slice-046). `slice-035` ist als reiner
Refactor noch in keinem Release.

## Nächste Wellen

| Welle | Trigger | Wichtigste Slices | Geschätzter Aufwand |
|---|---|---|---|

_(Keine Folge-Welle geplant — wartet auf Trigger.)_

## Meilensteine

| Meilenstein | Welle(n) | Trigger | Status |
|---|---|---|---|
| M1: Spec-Fundament steht | welle-01 | Closure-Trigger welle-01 | **erreicht** (2026-06-10) |
| M2: Dogfooding — `d-check` prüft die eigene Doku | welle-02 | slice-004 done, vendorter Bootstrap-Sensor gelöscht | **erreicht** (2026-06-10) |
| M3: erstes GHCR-Release + Pilot-Migration | welle-04 | Image veröffentlicht, ≥1 Repo migriert | **erreicht** (2026-06-12: v0.1.0/v0.2.0 auf GHCR, drei Repos migriert) |

## Abhängigkeitsgraph

```mermaid
flowchart LR
    W1[welle-01-fundament]
    W2[welle-02-mvp]
    W3[welle-03-regelmodule]
    W4[welle-04-distribution-und-migration]
    W5[welle-05-rollout]

    W1 --> W2
    W2 --> W3
    W3 --> W4
    W4 --> W5
```

## Abgeschlossene Wellen

| Welle | Abschluss | Closure-Notiz |
|---|---|---|
| welle-01-fundament | 2026-06-10 | [slice-001 §7](../done/slice-001-adr-fundament.md#7-closure-notiz-nach-done), [slice-002 §7](../done/slice-002-architektur-und-spezifikation.md#7-closure-notiz-nach-done) |
| welle-02-mvp | 2026-06-10 | [slice-003 §7](../done/slice-003-cli-kern-und-links-modul.md#7-closure-notiz-nach-done), [slice-004 §7](../done/slice-004-anchors-modul-und-dogfooding.md#7-closure-notiz-nach-done) |
| welle-03-regelmodule | 2026-06-11 | [slice-005 §7](../done/slice-005-lint-profil-solid.md#7-closure-notiz-nach-done), [slice-006 §7](../done/slice-006-ids-modul.md#7-closure-notiz-nach-done), [slice-007 §7](../done/slice-007-matrix-modul-selbstkonfiguration.md#7-closure-notiz-nach-done), [slice-008 §7](../done/slice-008-external-modul.md#7-closure-notiz-nach-done), [slice-009 §7](../done/slice-009-coverage-und-meta-gates.md#7-closure-notiz-nach-done) |
| welle-07-modul-scope | 2026-06-12 | [slice-017 §7](../done/slice-017-modul-scope.md#7-closure-notiz-nach-done) |
| welle-05-rollout | 2026-06-12 | [slice-014 §7](../done/slice-014-rollout-restliche-repos.md#7-closure-notiz-nach-done) |
| welle-04-distribution-und-migration | 2026-06-12 | [slice-010 §7](../done/slice-010-image-integrationstests-und-repro-belege.md#7-closure-notiz-nach-done), [slice-011 §7](../done/slice-011-ghcr-release-pipeline.md#7-closure-notiz-nach-done), [slice-013 §7](../done/slice-013-codepaths-modul.md#7-closure-notiz-nach-done), [slice-012 §7](../done/slice-012-pilot-migrationen.md#7-closure-notiz-nach-done) |
| welle-06-sensorik | 2026-06-13 | [slice-015 §7](../done/slice-015-spans-modul.md#7-closure-notiz-nach-done), [slice-016 §7](../done/slice-016-hostpaths-modul.md#7-closure-notiz-nach-done); Minor-Release **v0.4.0** auf GHCR (Run `27456611216` grün in 1m54s, OCI-Label `image.version`=0.4.0, Smoke-Lauf 42 Dateien/0 Befunde), Digest-Pin `ghcr.io/pt9912/d-check@sha256:3281ce538272fbfa086c2ee045a058542af0d1653425e995d56bd886ad730d61` |
| welle-08-linkpolitik | 2026-06-13 | [slice-018 §9](../done/slice-018-link-politik-ids.md#9-closure-notiz-nach-done); konfigurierbare `link-policy: prose\|always` für `ids` ([`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids) 0.8.0), kalibriert (d-check 170, u-boot 2, b-trace 2), Dogfooding-Sweep ~127 Links gate-verifiziert, nutzersichtbar dokumentiert |
| welle-09-config-geruest | 2026-06-13 | [slice-019 §7](../done/slice-019-print-config.md#7-closure-notiz-nach-done); Option `--print-config` gibt ein statisches `.d-check.yml`-Gerüst auf stdout ([`DC-FA-CLI-005`](../../../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben) 0.9.0), kein Repo-Zugriff/Schreiben, End-to-End belegt |
| welle-10-config-ableitung | 2026-06-13 | [slice-020 §7](../done/slice-020-suggest-config.md#7-closure-notiz-nach-done); Option `--suggest-config` leitet `ids`-Muster aus Autoritäts-Quellen ab ([`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten) 0.10.0, [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)-Schärfung), read-only; Korpora-Gegentest dokumentiert (d-check round-trippt, b-trace zeigt die Heading-Grenze) |
| welle-11-help | 2026-06-13 | [slice-021 §7](../done/slice-021-help-usage.md#7-closure-notiz-nach-done); reichhaltige `--help` (Synopsis, `[pfad]`-Argument, Config-Pointer auf `--print-config`) — Schärfung [`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel) 0.11.0 |
| welle-12-html-anker | 2026-06-15 | [slice-022 §6](../done/slice-022-html-anker.md#6-closure-notiz-nach-done); Inline-HTML-Anker als gültige Anker-Menge ([`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors) 0.12.0), `anchors`+`codepaths` via gemeinsame `AnchorSet`; zwei Review-Läufe (CR+Code, je R1); Minor-Release **v0.9.0** auf GHCR (Run `27541267343` grün in 2m2s), Digest-Pin `ghcr.io/pt9912/d-check@sha256:5bccf9fb3d1c54639dec3a541771d2ea43db9a0c1c58c28b3f12f20d38133d1b` |
| welle-13-ventil-prosa | 2026-06-16 | [slice-023 §8](../done/slice-023-ventil-prosa.md#8-closure-notiz-nach-done); die `ids`-Ventile `exempt-paths`/`d-check:ignore` gelten für nackte Fließtext-Vorkommen, nicht nur Inline-Code ([`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids) Change Request 0.13.0), als geteilte Exemption-Schicht; Review R1 (1 MEDIUM/4 LOW disponiert); Minor-Release **v0.10.0** auf GHCR (Run `27627924899` grün in 1m52s), Digest-Pin `ghcr.io/pt9912/d-check@sha256:ca49d33f22ecadfd08db03e4487b52b3f2a70dec01a41f2d0f472bfc2012797c` |
| welle-14-matrix-lineage | 2026-06-17 | [slice-024 §8](../done/slice-024-matrix-supersede-lineage.md#8-closure-notiz-nach-done); opt-in `allow-supersede-lineage` (+ `supersede-fields`) nimmt die deklarierte Supersede-Lineage-Kante des Moduls `matrix` von der Status-Prüfung aus ([`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix) Change Request 0.14.0); nur `matrix-inactive`, `matrix-forbidden` unberührt; `matrix` bleibt marker-frei (Entscheidung B2). `make gates` grün (Coverage 95,00 %), Vorher/Nachher am Image belegt (1 → 0 Befunde); Minor-Release **v0.11.0** auf GHCR (Run `27679319866` grün in 1m58s, Tags `v0.11.0`+`latest`), Digest-Pin `ghcr.io/pt9912/d-check@sha256:6ec1c463b5276b3314881839bd800b5e9aab12fa624a35d31618cecb62f17795` |
| welle-15-doctor-repair | 2026-06-18 | [slice-025 §7](../done/slice-025-doctor.md#7-closure-notiz-nach-done), [slice-026 §7](../done/slice-026-repair.md#7-closure-notiz-nach-done); Diagnose-Modus `--doctor` ([`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus) 0.15.0) und Reparatur-Patch `--repair` ([`DC-FA-CLI-008`](../../../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch) 0.16.0); read-only/stdout-only, zwei `--repair`-Stufen (konservativ/breit), `git apply`-Round-Trip belegt; `make gates` grün (Coverage 93,9 %); in Release **v0.12.0** ausgeliefert (Digest-Pin s. welle-16-Zeile) |
| welle-16-image-test-modi | 2026-06-18 | [slice-027 §7](../done/slice-027-image-test-modi.md#7-closure-notiz-nach-done); `make image-test` deckt `--doctor`/`--repair` nativ == Container byte-identisch ab ([`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)-Härtung, [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)); E2E-Lücke vor Release v0.12.0 geschlossen; Minor-Release **v0.12.0** auf GHCR (Run `27779122412` grün in 2m13s, Tags `v0.12.0`+`latest`), Digest-Pin `ghcr.io/pt9912/d-check@sha256:e65654ef8b35c9329f01eeee693bd0c10f583c9e6e01c89f24dd3c2615de32ac` |
| welle-17-benutzerhandbuch | 2026-06-18 | [slice-028 §7](../done/slice-028-benutzerhandbuch.md#7-closure-notiz-nach-done); aufgabenbasiertes Benutzerhandbuch unter docs/user/ nach dem adoptierten (mit-getrackten) Standard — alle Use Cases (v0.12.0), inkl. `--repair`-Pipe-Einzeiler; Self-Review R1, doc-check-rein; abgeleitete Nutzer-Doku (kein Vertrag) |
| welle-18-doctor-json | 2026-06-19 | [slice-029 §7](../done/slice-029-doctor-json.md#7-closure-notiz-nach-done); maschinenlesbare Diagnose `--doctor --json` ([`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus) Change Request 0.17.0) — dritte Ausgabe desselben Fix-Kandidaten-Modells, `findings` zusätzlich mit `reasonText`/`fixCandidate` (oder explizit `null`), `--doctor` nun mit `--json` kombinierbar; **unabhängiger** Review R1 (0 HIGH/MEDIUM/LOW, 2 INFO, INFO-1 in-flight geschlossen); `make gates` grün (Coverage 93,80 %); Minor-Release **v0.17.0** auf GHCR (Run `27806700510` grün in 1m56s, Tags `v0.17.0`+`latest`), Digest-Pin `ghcr.io/pt9912/d-check@sha256:fe8a1ccd718c04005e814aae7d82d32dc8f320688e9b738c85d7d0f9ac08935d` |
| welle-19-suggest-ai-harness | 2026-06-19 | [slice-030 §7](../done/slice-030-suggest-config-ai-harness.md#7-closure-notiz-nach-done); `--suggest-config`-Harness-Vorlage in **zwei Modi** ([`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten) Change Request 0.18.1) — `ai-harness-init` (Voll-Kanon fürs leere Repo) und `ai-harness` (repo-bewusst); kanonische ids-/matrix-/Modulset-Vorlage (Spiegel der `.d-check.yml`), read-only/advisory, deterministisch; Henne-Ei-Aufteilung nach Auftraggeber-Einwand; **unabhängiges** Review R1 (0 HIGH/MEDIUM/LOW, 2 INFO in-flight geschlossen); `make gates` grün (Coverage 93,70 %); Minor-Release **v0.18.0** auf GHCR (Run `27836313193` grün in 1m53s, Tags `v0.18.0`+`latest`), Digest-Pin `ghcr.io/pt9912/d-check@sha256:9c52e2d0e18de32146d0383257d240135288f2d1c25941e0fd08a465b8933e5c` |
| welle-20-yaml-ausgabe | 2026-06-19 | [slice-031 §7](../done/slice-031-yaml-ausgabe.md#7-closure-notiz-nach-done); Ausgabeformat **YAML** (`--yaml`, [`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate) Change Request 0.19.0) — strukturgleich zu `--json`, volle Parität inkl. `--doctor --yaml`; yaml.v3 zusätzlich im report-Adapter ([ADR-0009](../../adr/0009-yaml-im-report-adapter.md)), arch-check R3 erweitert; **unabhängiges** Review R1 (0 HIGH/MEDIUM, 2 LOW in-flight geschlossen); `make gates` grün (Coverage 94,10 %); in Release **v0.19.0** ausgeliefert (Digest s. welle-22-Zeile) |
| welle-21-semgrep-gate | 2026-06-20 | [slice-032 §7](../done/slice-032-semgrep-gate.md#7-closure-notiz-nach-done); hermetisches semgrep-Gate ([ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md)) — gepinntes Image `semgrep/semgrep:1.167.0` + gepinnter Regel-Commit, Umfang `go/lang/security` (55 Regeln, 0 Befunde), Host-XDG-Cache, `--network none`; Anti-Silent-Green erzwingt eine `Ran N rules`-Zeile (leerer Cache ⇒ Exit 2); Review R1 (HIGH-1 stilles Grün behoben); `make gates` grün; in Release **v0.19.0** ausgeliefert (Digest s. welle-22-Zeile) |
| welle-22-digest-pins | 2026-06-20 | [slice-033 §7](../done/slice-033-dockerfile-digest-pins.md#7-closure-notiz-nach-done); alle vier extern bezogenen Images per `@sha256:`-**Manifest-Listen**-Digest (amd64+arm64) inline neben dem Tag gepinnt — drei `Dockerfile`-`FROM` + semgrep ([ADR-0011](../../adr/0011-digest-pins-build-gate-images.md), vereinheitlicht [ADR-0002](../../adr/0002-distribution-ghcr-image.md)/[ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md) ohne Edit); `make versions` belegt alle Pins, `make ci` grün inkl. `image-test`; Review R1 (0 HIGH/1 MEDIUM verifiziert/2 LOW); Minor-Release **v0.19.0** auf GHCR (2026-06-20, Run `27871183930` grün, Tags `v0.19.0`+`latest`, bündelt slice-031/032/033), Digest-Pin `ghcr.io/pt9912/d-check@sha256:6134b8bd963de188858357ba05861a849dfb79dfac774437818f976100909ceb` |
| welle-24-kern-paketschnitt | 2026-06-20 | [slice-035 §7](../done/slice-035-kern-paketschnitt.md#7-closure-notiz-nach-done); Kern (5.212 Z., ein Paket) in drei Pakete `model`/`rules`/`app` geschnitten, Importrichtung `app→rules→model` per arch-check **R6** erzwungen ([ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)); **kein Verhaltens-Delta** (Befundsatz byte-identisch, alle Tests grün), `make gates` grün (Coverage 93 %); Review R1 merge-fähig (0 HIGH/0 MEDIUM/1 LOW/3 INFO); noch in keinem Release ausgeliefert (reiner Refactor) |
| welle-23-latest-tag | 2026-06-21 | [slice-034 §7](../done/slice-034-latest-tag-versoehnen.md#7-closure-notiz-nach-done); [ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md) ratifiziert die `:latest`-für-stabile-Praxis und löst [ADR-0002](../../adr/0002-distribution-ghcr-image.md) §4 (Tagging-Klausel „kein `latest`") teil-ab — Konsum verbindlich per `@sha256:`-Digest ([ADR-0011](../../adr/0011-digest-pins-build-gate-images.md), [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)); **kein Verhaltens-Delta** (Code/Doku fuhren A bereits), `make gates` grün; unabhängiges Review R1 (0 HIGH/0 MEDIUM/1 LOW/2 INFO); kein Release (Doku-Ratifikation) |
| welle-25-pr-ci-traceability | 2026-06-21 | [slice-039 §7](../done/slice-039-pr-ci-traceability-gate.md#7-closure-notiz-nach-done); PR-/Push-CI (`ci.yml`) ruft `make ci` + `make trace-check`; Traceability-Gate (`tools/trace-check.sh` + `commit-msg`-Hook via `make hooks`) erzwingt DC-/ADR-/MR-/slice-ID in Commits ([ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)); unabhängiges Review R1 (2 HIGH/1 MEDIUM behoben) + R2 (HIGH-A behoben), adversarial verifiziert; `make gates` grün; kein Release (Harness-Infra) |
| welle-26-suggest-prefix | 2026-06-21 | [slice-037 §8](../done/slice-037-suggest-config-id-prefix.md#8-closure-notiz-nach-done); `--suggest-config ai-harness[-init]` Kennungs-Präfix parametrisierbar — Flag `--id-prefix`, Ableitung aus dem Lastenheft (`ai-harness`), Platzhalter `<PREFIX>` + TODO statt fixem `DC-` ([`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten) 0.20.0, [ADR-0015](../../adr/0015-suggest-config-id-prefix.md)); **Breaking** (Init ohne Präfix → Platzhalter); unabhängiges Review R1 (2 MEDIUM behoben) + R2 (bestätigt); `make gates` grün; kein Release |
| welle-27-rtm-trace | 2026-06-21 | [slice-036 §8](../done/slice-036-rtm-trace.md#8-closure-notiz-nach-done); read-only-Modus `--trace` gibt eine Requirements Traceability Matrix aus (je Anforderung referenzierende ADRs/Slices + Waisen-Markierung), Default Markdown, optional `--trace --json`/`--yaml` ([`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix), Lastenheft 0.21.0); Doku-only, eigene Ableitung in `app` (ids/matrix liefern nur Findings), kein ADR (additiv); unabhängiges Review R1 (1 LOW behoben) + R2 (LOW-2 behoben); `make gates` grün; kein Release |
| welle-28-print-mk | 2026-06-21 | [slice-038 §8](../done/slice-038-print-mk.md#8-closure-notiz-nach-done); read-only-Generator `--print-mk` gibt ein include-bares `d-check.mk` aus — überschreibbare `DCHECK_IMAGE`-Variable (version-gepinntes Image, beim Tag-Build via `-ldflags -X` eingebettet; Digest via Override) + `doc-check`-Target ([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben), Lastenheft 0.22.0); kein ADR (Henne-Ei beim Digest, Version-Tag-Default konsistent mit Konsum-Pin-Politik); unabhängiges Review R1 (0 HIGH/0 MEDIUM/1 LOW won't-fix) + R2 (bestätigt); `make gates` grün; kein Release |
| welle-29-planning-consistency | 2026-06-21 | [slice-040 §7](../done/slice-040-planning-consistency-gate.md#7-closure-notiz-nach-done); Meta-Gate `make planning-check` erzwingt Roadmap §Aktuelle Welle ↔ `in-progress/slice-*` (beide Richtungen, fail-closed), in `make gates`; kein ADR; unabhängiges Review R1 (1 MEDIUM Heading-Guard behoben); `make gates` grün; kein Release |
| welle-30-adr-immutable | 2026-06-21 | [slice-041 §7](../done/slice-041-adr-immutable-gate.md#7-closure-notiz-nach-done); ADR-Immutable-Gate `make adr-check` erzwingt [`AGENTS.md` §3.5](../../../../AGENTS.md#35-adrs-sind-nach-accepted-immutable) (Accepted-ADRs nur `## Geschichte`-Anhang + Status-Übergang; CI-Range + pre-commit-Hook; [ADR-0016](../../adr/0016-adr-immutable-gate.md)); unabhängiges Review R1 (1 MEDIUM core-Status-Strip behoben); `make gates` grün; kein Release |
| welle-31-requirements-completeness | 2026-06-22 | [slice-042 §7](../done/slice-042-requirements-completeness-gate.md#7-closure-notiz-nach-done); Closure-Meta-Gate `make completeness-check` failt bei Requirements-Waisen (`--trace --json`, `orphans>0`), an `make fullbuild` (bewusst nicht `gates`/`ci`); [ADR-0017](../../adr/0017-requirements-completeness-gate.md); unabhängiges Review R1 (NACHBESSERN, F-1/F-2/F-3 vor Accept gefixt); `make gates` grün; kein Release |
| welle-32-codepaths-exempt | 2026-06-22 | [slice-043 §7](../done/slice-043-codepaths-exempt-paths.md#7-closure-notiz-nach-done); Modul `codepaths` bekam das Datei-Ventil `exempt-paths` (Glob wie `scan.ignore`, datei-weit) wie `ids` ([`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) Change Request, Lastenheft 0.23.0); Dogfooding `.d-check.yml` exemptet `docs/reviews/**`; kein ADR; unabhängiges Review R1 ACCEPT; Minor-Release **v0.23.0** auf GHCR, Digest-Pin `ghcr.io/pt9912/d-check@sha256:68951f5a3dd7ad3404e1996d45327f3df2585c0ef2b0b6bde7ccf790da4ddf6a` |
| welle-33-print-mk-trace | 2026-06-23 | [slice-044 §7](../done/slice-044-print-mk-trace-targets.md#7-closure-notiz-nach-done); opt-in `--trace --require-complete` bindet Requirements-Waisen an Exit 1 (neue [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code); Default-`--trace` bleibt advisory Exit 0); `--print-mk`-Fragment um `doc-trace`/`doc-complete`-Targets + `TRACE_FLAGS` erweitert ([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) Change Request, Lastenheft 0.24.0); kein ADR; unabhängiges R1 ACCEPT + R2 (Doku-Drift F-A behoben); `make gates`+`completeness-check` grün; Minor-Release **v0.24.0** auf GHCR (Run `28008942708` grün in 2m11s, Tags `v0.24.0`+`latest`), Digest-Pin `ghcr.io/pt9912/d-check@sha256:1c28a2b7e0e624763577ecba75b027f384692ecaa8a78a6e353a1a0c1889a4f8` |
| welle-34-diagram-ids | 2026-06-23 | [slice-045 §7](../done/slice-045-diagram-ids.md#7-closure-notiz-nach-done); opt-in Modul `diagrams` öffnet gezielt benannte Diagramm-Fences (Default `mermaid`) und prüft die darin gefundenen Kennungen auf **Existenz** in ihrer `defined-in`-Quelle (Befund `diagram-id-undefined`); Existenz statt Link-Policy (in Fences kein Markdown-Link), reine Token-Extraktion ohne Mermaid-Parser, scoped Fence-Ausnahme ([`DC-FA-DIAG-001`](../../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in), [ADR-0018](../../adr/0018-diagram-fence-ausnahme.md), Lastenheft 0.25.0); doc-first-Fundament R1 (NACHBESSERN→behoben)→R2 ACCEPT + Implementierungs-R1 ACCEPT (F-4-Tests ergänzt); `make gates` grün; Minor-Release **v0.25.0** auf GHCR (Run `28031261024` grün in 2m20s, Tags `v0.25.0`+`latest`), Digest-Pin `ghcr.io/pt9912/d-check@sha256:a2c5428214f1b3c616e0ba2e8d25bf77e4b11bf74470f10c1cd65d748667eb0f` |
| welle-35-suggest-opt-in-hinweis | 2026-06-23 | [slice-046 §7](../done/slice-046-suggest-config-opt-in-hinweis.md#7-closure-notiz-nach-done); Schärfung [`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten): `--suggest-config ai-harness[-init]` nennt die nicht aktivierten situativen opt-in-Module (`external`/`spans`/`hostpaths`/`diagrams`) in einem Kommentar mit Verweis auf `--print-config` (Auffindbarkeit ohne Aktivieren eines inerten Moduls; `diagrams` nicht ableitbar — braucht repo-spezifische `patterns`/`defined-in`); Lastenheft 0.26.0, kein ADR; unabhängiges R1 ACCEPT (F-1 LOW `external` ergänzt); `make gates` grün; Minor-Release **v0.26.0** auf GHCR (Run `28040897654` grün, Tags `v0.26.0`+`latest`), Digest-Pin `ghcr.io/pt9912/d-check@sha256:19d53a26d8d82a919015a8befe24f852bd61f2ddea58bd29e3f4cf944a8403f3` |
| welle-36-print-mk-erweiterung | 2026-06-23 | [slice-047 §7](../done/slice-047-print-mk-doctor-repair-help-digest.md#7-closure-notiz-nach-done); CR [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben): `--print-mk`-Fragment um `doc-doctor` (`--doctor`), `doc-repair` (`--repair`, `git apply`-rein via `@`), `doc-help` (namespaced Self-Doku) + `DCHECK_DIGEST` (Digest-Override per `ifeq`, sticht den Tag) erweitert; alle Targets `##`-annotiert; Lastenheft 0.27.0, kein ADR; unabhängiges R1 ACCEPT (3 INFO; Fragment real `make -n`-validiert, Digest-Override belegt); `make gates` grün; Minor-Release **v0.27.0** auf GHCR (Run `28047075398` grün, Tags `v0.27.0`+`latest`), Digest-Pin `ghcr.io/pt9912/d-check@sha256:2bc2598cbcd3622d98b33864a112fce02150b44776fc930fa404c98bd01668e1` |
| welle-37-versions | 2026-06-24 | [slice-048 §7](../done/slice-048-versions-modul.md#7-closure-notiz-nach-done); neues opt-in Modul `versions` (10. Modul): Versions-Pin-Konsistenz — gepinnte `ghcr`-Image-Verweise gegen die aktuelle Version aus `version.md#aktuell`, Befund `version-stale`, liest auch Fenced-Code (gescopte Fence-Ausnahme), Ventile `exempt-paths`/`d-check:ignore`, fail-closed, diagnose-only ([`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in), [ADR-0019](../../adr/0019-versions-pin-fence-ausnahme.md), Lastenheft 0.28.0); Doku-Boden `version.md` (only-current-anchor) + Release-Prep (alle ghcr-Pins → v0.28.0); Plan-Review R1→R3 ACCEPT + Impl-Review (4 Befunde); `make gates` grün (Coverage 93,90 %); Minor-Release **v0.28.0** auf GHCR (Run `28095582612` grün in 2m21s, Tags `v0.28.0`+`latest`), Digest-Pin `ghcr.io/pt9912/d-check@sha256:0bb84b529d3a65bdf9e849dd79cb8e9011bc388ecf9bffc5930f6c96bcc0cba8` |

## Historische Trigger-Verschiebungen

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| 2026-06-11 | slice-012-Trigger: „slice-011 done" → „slice-011 **und** slice-013 done" | Der [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Vergleichslauf gegen das erweiterte `docs-check.js` zeigte die Inline-Code-Pfad-Prüfung als Konsolidierungs-Lücke; Change Request [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) (Lastenheft 0.3.0) als slice-013 eingeschoben |
