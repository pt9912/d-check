# Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-07-15.

**Form:** folgt [Kurs-Modul 6](../../../../.harness/baseline/v1.4.0/regelwerk/modul-06-roadmap.md).

---

## Aktuelle Welle

**Welle-ID:** welle-59-trace-tabellenquellen

**Slice:**
[`slice-070`](slice-070-trace-tabellenquellen-nullmengen-guard.md) — native,
über Header-Namen konfigurierte Markdown-Pipe-Tabellen als
`--trace`-Anforderungsquelle plus fail-closed Nullmengen-Guard bei expliziter
Quelle; Heading-Default byte-identisch. Vertrag: Lastenheft-CR 0.43.0,
[ADR-0037](../../adr/0037-trace-tabellenquellen-nullmengen-guard.md).

**Vorgänger-Trigger:** welle-58-trace-handbuch-parsergrenzen ist nach
ACCEPT-Folgereview abgeschlossen; kein anderer Slice liegt in `in-progress/`.

**Closure-Trigger:**
- slice-070 in `done/`.

## Nächste Wellen

**Im Backlog (`next/`), auf Aufnahme in eine Welle wartend:** derzeit keiner.

**Kandidat (noch kein Slice, auf Freigabe wartend):** ein `--print-version-md`-Scaffold, das ein
`version.md`-Skelett mit Platzhaltern auf stdout ausgibt (Familie `--print-config`/`--print-mk`/
`--suggest-config`; read-only, deterministisch). Produkt-Feature ⇒ Change Request
(`DC-FA-CLI-*` im Lastenheft) + Slice + Spezifikation-`.a`, **kein** ADR (additive CLI-Ausgabe).
Anlass: Nutzer-Frage 2026-07-04 zum Nachbau von `version.md` in Fremd-Repos (der Aufbau selbst ist
seit Handbuch 1.21 dokumentiert).

## Historische Trigger-Verschiebungen

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| 2026-06-11 | slice-012-Trigger: „slice-011 done" → „slice-011 **und** slice-013 done" | Der [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Vergleichslauf gegen das erweiterte `docs-check.js` zeigte die Inline-Code-Pfad-Prüfung als Konsolidierungs-Lücke; Change Request [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) (Lastenheft 0.3.0) als slice-013 eingeschoben |
