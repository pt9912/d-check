# Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-07-17.

**Form:** folgt [Kurs-Modul 6](../../../../.harness/baseline/v1.4.0/regelwerk/modul-06-roadmap.md).

---

## Aktuelle Welle

**Welle-ID:** welle-60-trace-cross-consistency

**Slice:** [`slice-071`](slice-071-trace-cross-consistency-gate.md) — der
`--trace`-Lauf vergleicht opt-in die Vorwärts-RTM-Tabelle (Anforderung → Design)
gegen die Rückwärts-`Bezug`-Kanten (Design → Anforderung) und meldet je
Anforderung beide Mengendifferenzen. Vertrag:
[`DC-FA-XREF-001`](../../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
(Lastenheft 0.44.0) und [ADR-0038](../../adr/0038-trace-cross-consistency.md).

**Vorgänger-Trigger:** welle-59-trace-tabellenquellen ist abgeschlossen
([`slice-070`](../done/slice-070-trace-tabellenquellen-nullmengen-guard.md) in
`done/`, v0.43.1 veröffentlicht); kein anderer Slice liegt in `in-progress/`.

**Trigger:** Auftraggeber-Befund grid-gym (Trigger 088) — die §27.1-Vorwärts-Zeile
einer Architektur-Anforderung und ihre `Bezug`-Rück-Kanten hatten Schnittmenge
null, von keinem Gate bemerkt.

**Closure-Trigger:**

- Beide Richtungsdifferenzen, `superset`, range-aware Expansion und das
  `exclude-req`-Ventil sind als Akzeptanztests verriegelt; der Default ohne
  `trace.cross-consistency`-Block ist byte-identisch belegt.
- Der reale grid-gym-Drift wird geflaggt, die Mittelschicht-Familien nicht.
- [ADR-0038](../../adr/0038-trace-cross-consistency.md) ist `Accepted`;
  unabhängiger, kontext-getrennter Closure-Review liegt vor.
- `make gates` und `make ci` grün, Release samt GHCR-Digest-Backfill dokumentiert.

## Nächste Wellen

**Im Backlog (`next/`), auf Aufnahme in eine Welle wartend:** derzeit keiner.

**Im Eingang (`open/`), auf Wellen-Einplanung wartend:** slice-073 —
link-transparente Range-Fortsetzung
([ADR-0039](../../adr/0039-link-transparente-range-fortsetzung.md) Proposed;
Defekt-Fix, kein CR). **Ausgelieferter Defekt seit v0.41.0**: eine verlinkte Range
expandiert nicht, `trace.coverage` meldet dadurch falsche Waisen — die
Range-Zusage kollidiert strukturell mit d-checks eigener Linkpflicht. Kandidat für
welle-60, weil der geteilte Parser zugleich den Realdatenbeleg von slice-071
blockiert. Priorität vor slice-071s Rest-Arbeit (trifft Bestandskonsumenten).

Ferner slice-072 — Handbuch-
Aufgabenorientierung der §4-Kapitel gegen den
[Benutzerhandbuch-Standard](../../../user/benutzerhandbuch-standard.md) §2
(sieben Audit-Befunde; rein redaktionell — kein Change Request, kein ADR, kein
Release). Ursache ist strukturell: §4.12 wuchs über die Slices 066–071, weil jeder
Slice seine Fähigkeit anhängte, statt eine Aufgabe zu schreiben.

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
