# Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-07-15.

**Form:** folgt [Kurs-Modul 6](../../../../.harness/baseline/v1.4.0/regelwerk/modul-06-roadmap.md).

---

## Aktuelle Welle

Keine aktive Welle.

## Nächste Wellen

**Im Backlog (`next/`), auf Aufnahme in eine Welle wartend:** derzeit keiner.

**Im Eingang (`open/`), doc-first aufgesetzt, auf Wellen-Einplanung wartend:**
slice-071 — Trace-Kreuzverweis-Konsistenz-Gate
([`DC-FA-XREF-001`](../../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in),
[ADR-0038](../../adr/0038-trace-cross-consistency.md)); Rückkanten-Vertrag gegen
grid-gyms reale Quellen geerdet, Generator als spätere CR sequenziert.

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
