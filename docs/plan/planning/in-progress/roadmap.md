# Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-07-17.

**Form:** folgt [Kurs-Modul 6](../../../../.harness/baseline/v1.4.0/regelwerk/modul-06-roadmap.md).

---

## Aktuelle Welle

**Welle-ID:** welle-60-trace-cross-consistency

**Slices:** [`slice-073`](slice-073-link-transparente-range-fortsetzung.md) —
link-transparente Range-Fortsetzung: eine verlinkte Range (`` [`ID`](…)..009 ``)
expandiert nicht, weil die Spec-Verengung „unmittelbar" strukturell mit d-checks
eigener Linkpflicht kollidiert. **Ausgelieferter Defekt seit v0.41.0**
(`trace.coverage` meldet falsche Waisen); Defekt-Fix, kein CR
([ADR-0039](../../adr/0039-link-transparente-range-fortsetzung.md)). Läuft
**vorrangig** — er trifft Bestandskonsumenten und blockiert zugleich, über den
geteilten Parser, den Realdatenbeleg von slice-071.

[`slice-076`](slice-076-markdown-lexik-commonmark.md) — Markdown-Lexik an
CommonMark/GFM angleichen: die Trennzelle verlangt drei Bindestriche statt einem
(GFM), und eine `` ``` ``-Zeile mit Backtick in der Infozeile öffnet fälschlich
einen Fence und blendet **alle** Module bis zum Dateiende. Beides **ausgeliefert
und still**; belegt per Differential gegen goldmark v1.8.4 über 522 reale Dateien
(490 Tabellen ⇒ 8 Abweichungen, alle „d-check ist blind"). Defekt-Fix, kein CR,
aber **SemVer-Minor** — d-check findet danach **mehr**
([ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md)). Unabhängig von
slice-073/slice-071, blockiert keinen von beiden.

[`slice-071`](slice-071-trace-cross-consistency-gate.md) — der
`--trace`-Lauf vergleicht opt-in die Vorwärts-RTM-Tabelle (Anforderung → Design)
gegen die Rückwärts-`Bezug`-Kanten (Design → Anforderung) und meldet je
Anforderung beide Mengendifferenzen. Vertrag:
[`DC-FA-XREF-001`](../../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
(Lastenheft 0.44.2) und [ADR-0038](../../adr/0038-trace-cross-consistency.md).

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
  **Offen:** der von [ADR-0038](../../adr/0038-trace-cross-consistency.md)
  Entscheidung 7 geforderte Realdatenbeleg bricht an grid-gyms
  `architecture.md:913` weiter mit Exit 2 ab — die Direktiven-Zeile ist mit der
  Rücknahme von slice-074 wieder unlesbar. Der Beleg hängt an slice-074/076.
- Eine verlinkte Range expandiert wie die unverlinkte (slice-073); der
  ausgelieferte `trace.coverage`-Falschbefund ist weg und als Patch veröffentlicht.
- [ADR-0038](../../adr/0038-trace-cross-consistency.md) und
  [ADR-0039](../../adr/0039-link-transparente-range-fortsetzung.md) sind
  `Accepted`; unabhängiger, kontext-getrennter Closure-Review liegt vor.
- `make gates` und `make ci` grün, Release samt GHCR-Digest-Backfill dokumentiert.

## Nächste Wellen

**Im Backlog (`next/`), auf Aufnahme in eine Welle wartend:** derzeit keiner.

**Im Eingang (`open/`), auf Wellen-Einplanung wartend:**
[`slice-074`](../open/slice-074-kommentar-suffix-tabellenzeilen.md) —
Direktiven-Zelle in Tabellenzeilen ([ADR-0040](../../adr/0040-kommentar-suffix-in-tabellenzeilen.md)
Proposed). **Aus `in-progress/` zurückgestellt, Implementierung zurückgenommen**
(2026-07-17): drei unabhängige Reviews haben an fünf aufeinanderfolgenden
Fassungen dieselbe Klasse belegt — zuletzt R3-F-1, ein Stilles-Grün-Pfad
(Exit 1 ⇒ Exit 0) in beiden Konsumenten. Der Defekt bleibt offen und
ausgeliefert (Exit 2 auf einer Zeile, die jeder Renderer normal darstellt);
die Rücknahme ist die ehrliche Zwischenlage, nicht die Lösung. Vorbedingung
für den Realdatenbeleg von slice-071. **Der Auftraggeber führt ihn (2026-07-17)
weiterhin als einen von drei offenen Punkten** — die Rücknahme hat den Blocker
bei ihm wiederhergestellt, sie hat ihn nicht erledigt.

Ferner [`slice-077`](../open/slice-077-stiller-tabellen-uebersprung.md) —
stiller Tabellen-Übersprung: eine **irrelevante** Tabelle gleicher Breite
verschluckt ohne Leerzeile die folgende **relevante**; deren Anforderungen
verschwinden lautlos (`0 Waise(n)`/Exit 0 bei zwei echten Waisen, gemessen gegen
das ausgelieferte v0.45.1, **ohne** jeden Marker). **Bewusst ohne ADR erfasst —
die tragende Regel ist offen**, und die naheliegende ist bereits widerlegt. Von
den offenen Tabellen-Defekten der einzige, der still Waisen verschweigt.

[`slice-075`](../open/slice-075-komma-kurzform-fail-closed.md) — Komma-Kurzform
fail-closed ([ADR-0041](../../adr/0041-komma-kurzform-fail-closed.md) Proposed;
**Change Request**, Lastenheft 0.46.0, SemVer-Minor). `GG-SCN-001, 007` deckte
nur die erste Kennung — still. Die Kurzform war nie zugesagt; der Defekt ist das
fehlende Signal, nicht die fehlende Unterstützung. **Vorrangig einzuplanen
(Auftraggeber-Meldung 2026-07-17):** der stille Drop **verfälscht produktiv
verdrahtetes** `trace.coverage` bei grid-gym — der einzige offene Punkt, der
heute still falsche Zahlen liefert und dabei verdrahtet ist.

Ferner slice-072 — Handbuch-
Aufgabenorientierung der §4-Kapitel gegen den
[Benutzerhandbuch-Standard](../../../user/benutzerhandbuch-standard.md) §2
(sieben Audit-Befunde; rein redaktionell — kein Change Request, kein ADR, kein
Release). Ursache ist strukturell: §4.12 wuchs über die Slices 066–071, weil jeder
Slice seine Fähigkeit anhängte, statt eine Aufgabe zu schreiben.

**Kandidat (noch kein Slice, auf Freigabe wartend):** der **RTM-Generator** —
die Vorwärts-RTM aus den Rückwärts-`Bezug`-Kanten **erzeugen**, statt sie nur
gegen sie abzugleichen. Von [ADR-0038](../../adr/0038-trace-cross-consistency.md)
Entscheidung 7 bewusst als **spätere CR** sequenziert; das
`trace.cross-consistency`-Gate aus
[`slice-071`](slice-071-trace-cross-consistency-gate.md) ist sein
Korrektheits-Harness. **Auftraggeber-Nachreichung 2026-07-17 (grid-gym), einer
von drei offenen Punkten:** der Generator soll zusätzlich **Artefakt-Titel** und
**Kanten-Anmerkung** tragen. Produkt-Feature ⇒ Change Request (`DC-FA-*` im
Lastenheft) + Slice + Spezifikation-`.a` + ADR. **Freigabe und Scope offen** —
insbesondere, woher der Artefakt-Titel kommt (die Rück-Sicht bindet heute nur die
erste Spalte als Artefakt-ID) und was eine Kanten-Anmerkung normativ ist.

Ferner ein `--print-version-md`-Scaffold, das ein
`version.md`-Skelett mit Platzhaltern auf stdout ausgibt (Familie `--print-config`/`--print-mk`/
`--suggest-config`; read-only, deterministisch). Produkt-Feature ⇒ Change Request
(`DC-FA-CLI-*` im Lastenheft) + Slice + Spezifikation-`.a`, **kein** ADR (additive CLI-Ausgabe).
Anlass: Nutzer-Frage 2026-07-04 zum Nachbau von `version.md` in Fremd-Repos (der Aufbau selbst ist
seit Handbuch 1.21 dokumentiert).

## Historische Trigger-Verschiebungen

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| 2026-06-11 | slice-012-Trigger: „slice-011 done" → „slice-011 **und** slice-013 done" | Der [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Vergleichslauf gegen das erweiterte `docs-check.js` zeigte die Inline-Code-Pfad-Prüfung als Konsolidierungs-Lücke; Change Request [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) (Lastenheft 0.3.0) als slice-013 eingeschoben |
