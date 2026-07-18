# Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-07-18.

**Form:** folgt [Kurs-Modul 6](../../../../.harness/baseline/v1.4.0/regelwerk/modul-06-roadmap.md).

---

## Aktuelle Welle

**Aktive Welle: welle-61-referenz-ventil-quell-skopus.** slice-078 ist in Arbeit
(`in-progress/`): das Referenz-Ventil `ignore-refs` bekommt einen **Quell-Skopus**
(`in:`) plus `refs`/`keep`-Felder (ignorieren, wenn `refs` matcht ∧ `keep` nicht) und
wandert von modul-lokal (`codepaths`) zu **querschnittlich** für
`links`/`anchors`/`codepaths`. Trigger: Konsumenten-CR `ai-harness-course` — die
Ziel-Repo-Platzhalter in Template-Verzeichnissen erzwingen heute das ganze
Verzeichnis in `scan.ignore` und machen damit auch die **echten** Verweise blind,
deren Auflösung beim Release unveränderlich eingefroren wird.

**§4-Vorfrage entschieden (Auftraggeber, 2026-07-18):** das Ventil wohnt als
**neues, geteiltes Bereichskürzel** — das Ziel-Achsen-Pendant zum bestehenden
[`DC-FA-SCAN-001`](../../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln)
der Quell-Achse; `links`/`anchors`/`codepaths` verweisen darauf,
`codepaths.ignore-refs` bleibt **Alias** (kein Config-Bruch), die Ventil-Spezifikation
steht **einmal**.
[`MR-002`](../../../../harness/conventions.md#mr-002--id-schema-mit-bereichskürzeln-ab-initialer-fassung)
deckt die Deklaration des neuen Kürzels. Doc-first: Lastenheft → ADR → Spezifikation
führen, Code folgt.

**Vorgänger:** welle-60-trace-cross-consistency **abgeschlossen** (alle
Closure-Kriterien erfüllt) —
[`slice-071`](../done/slice-071-trace-cross-consistency-gate.md) (Kreuzverweis-
Konsistenz, Realdatenbeleg erbracht,
[ADR-0038](../../adr/0038-trace-cross-consistency.md) `Accepted`, v0.44.0/v0.45.0),
[`slice-073`](../done/slice-073-link-transparente-range-fortsetzung.md) (v0.45.1),
[`slice-075`](../done/slice-075-komma-kurzform-fail-closed.md) (v0.46.0),
[`slice-076`](../done/slice-076-markdown-lexik-commonmark.md) (v0.47.0); die
Wurzel-/Aufsatz-Kette
[`slice-077`](../done/slice-077-stiller-tabellen-uebersprung.md) (v0.48.0) →
[`slice-074`](../done/slice-074-kommentar-suffix-tabellenzeilen.md) (v0.48.1) löste
die Tabellen-Reader-Klasse.

## Nächste Wellen

**Im Backlog (`next/`):** leer.

**Im Eingang (`open/`), auf Wellen-Einplanung wartend:**
[`slice-079`](../open/slice-079-zitat-verifikation.md) ·
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
| 2026-07-18 | **welle-61-referenz-ventil-quell-skopus eröffnet**; slice-078 `open`→`in-progress` (WIP-Slot frei nach welle-60-Abschluss) | §4-Vorfrage vom Auftraggeber entschieden: das erweiterte `ignore-refs`-Ventil (Quell-Skopus `in:`, `refs`/`keep`) wohnt als **neues geteiltes Bereichskürzel** (Ziel-Achsen-Pendant zu [`DC-FA-SCAN-001`](../../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln)), nicht als Änderung dreier Anforderungen — vermeidet die Verdreifachung der Ventil-Spezifikation, `codepaths.ignore-refs` bleibt Alias. Konsumenten-CR `ai-harness-course` |
