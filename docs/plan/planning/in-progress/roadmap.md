# Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-06-17.

**Format-Regel:** Die Roadmap ist eine Reihenfolge von **Wellen**,
keine Reihenfolge von Terminen. Termine erscheinen — falls überhaupt —
als Konsequenz der Wellen-Schätzung, nicht als Treiber.

---

## Aktuelle Welle

**Keine aktive Welle.** welle-14-matrix-lineage ist abgeschlossen (siehe
„Abgeschlossene Wellen"); das Release `v0.11.0` (GHCR-Digest-Pin) steht
noch aus (Review vor Release). Die nächste Welle wartet auf ihren Trigger
(Change Request im Lastenheft oder Priorisierung durch den
Auftraggeber).

## Nächste Wellen

| Welle | Trigger | Wichtigste Slices | Geschätzter Aufwand |
|---|---|---|---|

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
| welle-14-matrix-lineage | 2026-06-17 | [slice-024 §8](../done/slice-024-matrix-supersede-lineage.md#8-closure-notiz-nach-done); opt-in `allow-supersede-lineage` (+ `supersede-fields`) nimmt die deklarierte Supersede-Lineage-Kante des Moduls `matrix` von der Status-Prüfung aus ([`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix) Change Request 0.14.0); nur `matrix-inactive`, `matrix-forbidden` unberührt; `matrix` bleibt marker-frei (Entscheidung B2). `make gates` grün (Coverage 95,00 %), Vorher/Nachher am Image belegt (1 → 0 Befunde). **Release v0.11.0 (GHCR-Digest-Pin) folgt nach Review** |

## Historische Trigger-Verschiebungen

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| 2026-06-11 | slice-012-Trigger: „slice-011 done" → „slice-011 **und** slice-013 done" | Der [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Vergleichslauf gegen das erweiterte `docs-check.js` zeigte die Inline-Code-Pfad-Prüfung als Konsolidierungs-Lücke; Change Request [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) (Lastenheft 0.3.0) als slice-013 eingeschoben |
