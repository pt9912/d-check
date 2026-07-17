# Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-07-17.

**Form:** folgt [Kurs-Modul 6](../../../../.harness/baseline/v1.4.0/regelwerk/modul-06-roadmap.md).

---

## Aktuelle Welle

**Welle-ID:** welle-60-trace-cross-consistency

**Slices** (WIP-Limit = 1, Modul 5 — **genau einer in Arbeit**):
[`slice-075`](slice-075-komma-kurzform-fail-closed.md) **in Arbeit** —
Komma-Kurzform fail-closed: `GG-SCN-001, 007` ließ `007` **still** fallen und
verfälschte produktiv verdrahtetes `trace.coverage`; künftig Exit 2 mit
Notations-Hinweis (Change Request, [ADR-0041](../../adr/0041-komma-kurzform-fail-closed.md),
Lastenheft 0.46.0, SemVer-Minor). ·
[`slice-076`](../next/slice-076-markdown-lexik-commonmark.md) in `next/`, wartet
auf den WIP-Slot. · [`slice-071`](../open/slice-071-trace-cross-consistency-gate.md)
in `open/`, **blockiert** (Realdatenbeleg hängt an slice-074);
[`DC-FA-XREF-001`](../../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in),
Code und Reviews liegen vor, v0.44.0/v0.45.0 sind getaggt. ·
[`slice-073`](../done/slice-073-link-transparente-range-fortsetzung.md) in `done/`
(v0.45.1, R2 ACCEPT-WITH-NITS).

**Vorgänger-Trigger:** welle-59-trace-tabellenquellen abgeschlossen
([`slice-070`](../done/slice-070-trace-tabellenquellen-nullmengen-guard.md) in
`done/`, v0.43.1 veröffentlicht).

**Trigger:** Auftraggeber-Befund grid-gym (Trigger 088) — die §27.1-Vorwärts-Zeile
einer Architektur-Anforderung und ihre `Bezug`-Rück-Kanten hatten Schnittmenge
null, von keinem Gate bemerkt.

**Closure-Kriterien von welle-60:**

- Beide Richtungsdifferenzen, `superset`, range-aware Expansion und das
  `exclude-req`-Ventil sind als Akzeptanztests verriegelt; der Default ohne
  `trace.cross-consistency`-Block ist byte-identisch belegt.
- Der reale grid-gym-Drift wird geflaggt, die Mittelschicht-Familien nicht.
  **Offen:** der von [ADR-0038](../../adr/0038-trace-cross-consistency.md)
  Entscheidung 7 geforderte Realdatenbeleg bricht an grid-gyms
  `architecture.md:913` weiter mit Exit 2 ab — er hängt an slice-074.
- **[erfüllt]** Eine verlinkte Range expandiert wie die unverlinkte (slice-073);
  der ausgelieferte `trace.coverage`-Falschbefund ist weg und als Patch (v0.44.1,
  klammer-balanciert nachgebessert in v0.45.1) veröffentlicht. slice-073 in `done/`.
- Die zwei Lexik-Regeln aus slice-076 sind per Mutation gepinnt und als Minor
  veröffentlicht.
- [ADR-0038](../../adr/0038-trace-cross-consistency.md),
  [ADR-0039](../../adr/0039-link-transparente-range-fortsetzung.md) und
  [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) sind `Accepted`;
  unabhängiger, kontext-getrennter Closure-Review liegt vor.
- `make gates` und `make ci` grün, Release samt GHCR-Digest-Backfill dokumentiert.

## Nächste Wellen

**Im Backlog (`next/`), in welle-60 eingeplant, auf den WIP-Slot wartend:**
[`slice-076`](../next/slice-076-markdown-lexik-commonmark.md).

**Im Eingang (`open/`), auf Wellen-Einplanung wartend:**
[`slice-071`](../open/slice-071-trace-cross-consistency-gate.md) (blockiert) ·
[`slice-074`](../open/slice-074-kommentar-suffix-tabellenzeilen.md)
(zurückgestellt, Implementierung zurückgenommen) ·
[`slice-077`](../open/slice-077-stiller-tabellen-uebersprung.md) ·
[`slice-078`](../open/slice-078-ignore-refs-quell-skopus.md) · [`slice-079`](../open/slice-079-zitat-verifikation.md) ·
[`slice-072`](../open/slice-072-handbuch-aufgabenorientierung.md).

**Kandidat (noch kein Slice, auf Freigabe wartend):** der **RTM-Generator** (RTM
aus den Rückwärts-`Bezug`-Kanten erzeugen; von
[ADR-0038](../../adr/0038-trace-cross-consistency.md) Entscheidung 7 als spätere
CR sequenziert, slice-071 ist sein Korrektheits-Harness). Auftraggeber-Nachreichung
2026-07-17: zusätzlich Artefakt-Titel + Kanten-Anmerkung. Freigabe und Scope offen.

Ferner ein `--print-version-md`-Scaffold, das ein `version.md`-Skelett mit
Platzhaltern auf stdout ausgibt (Familie `--print-config`/`--print-mk`/
`--suggest-config`; read-only, deterministisch). Produkt-Feature ⇒ Change Request
(`DC-FA-CLI-*` im Lastenheft) + Slice + Spezifikation-`.a`, **kein** ADR (additive
CLI-Ausgabe). Anlass: Nutzer-Frage 2026-07-04 zum Nachbau von `version.md` in
Fremd-Repos (der Aufbau selbst ist seit Handbuch 1.21 dokumentiert).

## Historische Trigger-Verschiebungen

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| 2026-06-11 | slice-012-Trigger: „slice-011 done" → „slice-011 **und** slice-013 done" | Der [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Vergleichslauf gegen das erweiterte `docs-check.js` zeigte die Inline-Code-Pfad-Prüfung als Konsolidierungs-Lücke; Change Request [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) (Lastenheft 0.3.0) als slice-013 eingeschoben |
| 2026-07-17 | **WIP-Limit wiederhergestellt:** slice-071 `in-progress`→`open` (Blocker), slice-076 `in-progress`→`next`; welle-60 führt nur noch slice-073 in Arbeit. Reihenfolge danach: slice-073 zu Ende (vier offene R1-Befunde + bestätigender Review) → Closure → slice-075 | `in-progress/` trug **drei** Slices gleichzeitig; Modul 5: „WIP-Limit pro Implementer = 1 ist eine harte Größe, kein Vorschlag" und `next→in-progress` verlangt „WIP-Limit frei". Bei slice-076 wurde die Bedingung beim Einplanen schlicht nicht geprüft (`6d60094`); slice-071 war bereits blockiert und hätte nach Modul 5 längst zurückgeführt gehört — beides still, bis der Auftraggeber die Regel einforderte. slice-075 erhält Vorrang vor slice-076, weil er produktiv verdrahtetes `trace.coverage` **verfälscht** (Auftraggeber-Meldung grid-gym), während slice-076 Blindheit ohne Falschaussage ist |
| 2026-07-17 | slice-074 aus welle-60 zurückgestellt (`in-progress/` → `open/`), Implementierung zurückgenommen; slice-076 in welle-60 nachgenommen | Drei unabhängige Reviews belegten an fünf aufeinanderfolgenden Fassungen dieselbe Klasse, zuletzt einen Stilles-Grün-Pfad (R3-F-1). Der Realdatenbeleg für slice-071 ist damit weiter blockiert — offen ausgewiesen statt still weitergeschoben. slice-076 kam aus dem Spike, den die Rücknahme ausgelöst hat |
