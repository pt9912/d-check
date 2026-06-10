# Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-06-10.

**Format-Regel:** Die Roadmap ist eine Reihenfolge von **Wellen**,
keine Reihenfolge von Terminen. Termine erscheinen — falls überhaupt —
als Konsequenz der Wellen-Schätzung, nicht als Treiber.

---

## Aktuelle Welle

**Welle-ID:** welle-01-fundament
**Start:** 2026-06-10
**Slices:** [slice-001](../open/slice-001-adr-fundament.md), [slice-002](../open/slice-002-architektur-und-spezifikation.md)

**Closure-Trigger:** Fundament-ADRs (Sprache, Distribution,
Config-Format) `Accepted`; `spec/spezifikation.md` und
`spec/architecture.md` existieren; Source-Precedence-Tabellen von
„geplant" auf Links umgestellt; `make gates` grün.

## Nächste Wellen

| Welle | Trigger | Wichtigste Slices | Geschätzter Aufwand |
|---|---|---|---|
| welle-02-mvp | welle-01 done | [slice-003](../open/slice-003-cli-kern-und-links-modul.md), [slice-004](../open/slice-004-anchors-modul-und-dogfooding.md) | M |
| welle-03-regelmodule | welle-02 done | Module `ids`, `matrix`, `external`; Gate-Ausbau: `coverage-gate` (bootstrap-aware), `gate-consistency` (Slices werden bei Wellen-Start geschnitten) | M |
| welle-04-distribution-und-migration | welle-03 done | GHCR-Release-Pipeline ([`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)), Reproduzierbarkeits-Belege (`versions`, `fullbuild` mit Image-Hash), Pilot-Migrationen in 3 Repos ([`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)) | M |

## Meilensteine

| Meilenstein | Welle(n) | Trigger | Status |
|---|---|---|---|
| M1: Spec-Fundament steht | welle-01 | Closure-Trigger welle-01 | offen |
| M2: Dogfooding — `d-check` prüft die eigene Doku | welle-02 | slice-004 done, vendorter Bootstrap-Sensor gelöscht | offen |
| M3: erstes GHCR-Release + Pilot-Migration | welle-04 | Image veröffentlicht, ≥1 Repo migriert | offen |

## Abhängigkeitsgraph

```mermaid
flowchart LR
    W1[welle-01-fundament]
    W2[welle-02-mvp]
    W3[welle-03-regelmodule]
    W4[welle-04-distribution-und-migration]

    W1 --> W2
    W2 --> W3
    W3 --> W4
```

## Abgeschlossene Wellen

| Welle | Abschluss | Closure-Notiz |
|---|---|---|
| — noch keine — | | |

## Historische Trigger-Verschiebungen

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| — | | |
