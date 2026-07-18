# Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-07-18.

**Form:** folgt [Kurs-Modul 6](../../../../.harness/baseline/v1.4.0/regelwerk/modul-06-roadmap.md).

---

## Aktuelle Welle

**Keine aktive Welle.** welle-60-trace-cross-consistency ist **abgeschlossen**:
slice-071 ist geschlossen (`done/`) — der letzte offene DoD-Punkt, der von
[ADR-0038](../../adr/0038-trace-cross-consistency.md) Entscheidung 7 geforderte
**Realdatenbeleg** gegen grid-gyms reale §27.1-↔-`Bezug`-Quellen, ist **erbracht**
(d-check v0.48.1): der Drift wird geflaggt (`GG-ARCH-005`/`GG-SIM-009`), das
`exclude-req`-Ventil greift (0 Befunde), das konsistente 1:N läuft grün, die
dokumentierte 161-Differenzen-Messung reproduziert; der zuvor an
`architecture.md:913` blockierende `17. Testarchitektur`-Abschnitt läuft mit der
Direktiven-Toleranz aus slice-074 durch. Die begründende ADR ist `Accepted`.
`in-progress/` ist leer, der WIP-Slot frei.

**Stand von welle-60** (abgeschlossen):
[`slice-073`](../done/slice-073-link-transparente-range-fortsetzung.md) **done**
(v0.45.1, R2 ACCEPT-WITH-NITS). ·
[`slice-075`](../done/slice-075-komma-kurzform-fail-closed.md) **done**
(Komma-Kurzform fail-closed, R1 ACCEPT-WITH-NITS, **v0.46.0 veröffentlicht**). ·
[`slice-076`](../done/slice-076-markdown-lexik-commonmark.md) **done**
(Markdown-Lexik CommonMark/GFM, ACCEPT-WITH-NITS, **v0.47.0 veröffentlicht**). ·
[`slice-071`](../done/slice-071-trace-cross-consistency-gate.md) **done**
(Realdatenbeleg erbracht, [ADR-0038](../../adr/0038-trace-cross-consistency.md)
`Accepted`, R1–R4 + Closure-Review; **v0.44.0/v0.45.0 veröffentlicht**),
[`DC-FA-XREF-001`](../../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in).

**Vorgänger-Trigger:** welle-59-trace-tabellenquellen abgeschlossen
([`slice-070`](../done/slice-070-trace-tabellenquellen-nullmengen-guard.md) in
`done/`, v0.43.1 veröffentlicht).

**Trigger:** Auftraggeber-Befund grid-gym (Trigger 088) — die §27.1-Vorwärts-Zeile
einer Architektur-Anforderung und ihre `Bezug`-Rück-Kanten hatten Schnittmenge
null, von keinem Gate bemerkt.

**Closure-Kriterien von welle-60:**

- **[erfüllt]** Beide Richtungsdifferenzen, `superset`, range-aware Expansion und das
  `exclude-req`-Ventil sind als Akzeptanztests verriegelt; der Default ohne
  `trace.cross-consistency`-Block ist byte-identisch belegt (slice-071).
- **[erfüllt]** Der reale grid-gym-Drift wird geflaggt (`GG-ARCH-005`/`GG-SIM-009`),
  die per `exclude-req` ausgeschlossenen Familien nicht. Der von
  [ADR-0038](../../adr/0038-trace-cross-consistency.md) Entscheidung 7 geforderte
  Realdatenbeleg ist erbracht (slice-071, d-check v0.48.1); der zuvor an grid-gyms
  `architecture.md:913` mit Exit 2 abbrechende `17. Testarchitektur`-Abschnitt läuft
  mit der Direktiven-Toleranz aus slice-074 durch.
- **[erfüllt]** Eine verlinkte Range expandiert wie die unverlinkte (slice-073);
  der ausgelieferte `trace.coverage`-Falschbefund ist weg und als Patch (v0.44.1,
  klammer-balanciert nachgebessert in v0.45.1) veröffentlicht. slice-073 in `done/`.
- **[erfüllt]** Die zwei Lexik-Regeln aus slice-076 sind per Mutation gepinnt und
  als Minor (v0.47.0) veröffentlicht. slice-076 in `done/`.
- **[erfüllt]** [ADR-0038](../../adr/0038-trace-cross-consistency.md),
  [ADR-0039](../../adr/0039-link-transparente-range-fortsetzung.md) und
  [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) sind `Accepted`;
  unabhängiger, kontext-getrennter Closure-Review liegt vor.
- **[erfüllt]** `make gates` und `make ci` grün, Release samt GHCR-Digest-Backfill
  dokumentiert.

**welle-60 ist damit vollständig abgeschlossen** — alle Closure-Kriterien erfüllt.

## Nächste Wellen

**Im Backlog (`next/`):** leer.

**Im Eingang (`open/`), auf Wellen-Einplanung wartend:**
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
| 2026-07-18 | slice-071 wieder aufgenommen (`open/`→`in-progress/`), welle-60 wieder aktiv | Blocker aufgelöst: die Direktiven-Toleranz aus slice-074 (v0.48.1) lässt den `17. Testarchitektur`-Abschnitt durchlaufen, statt an `architecture.md:913` abzubrechen. WIP-Slot frei (Modul 5), daher `open→next→in-progress` in einem Zug. Der von [ADR-0038](../../adr/0038-trace-cross-consistency.md) Entscheidung 7 geforderte Realdatenbeleg gegen grid-gym ist damit fahrbar |
