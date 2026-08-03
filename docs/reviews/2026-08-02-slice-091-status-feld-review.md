# Review-Report: slice-091 — Slice-`Status:`-Feld abschaffen (welle-67, Etappe D, D-5) — 2026-08-02

**Review-Art:** Plan-Review (gegen die adoptierte Baseline-Form: `modul-05` §Lifecycle
+ `slice.template.md`) mit Code-/Gate-Anteil (Faktentreue „kein Gate liest das Feld",
zwei `Status:`-Achsen unberührt).

**Gegenstand:** slice-091, Commit-Range `737ee50..HEAD` (2 Commits: `49e016f` open,
`a9af858` D-5).

**Skill:** [`reviewer.md`](../../.harness/skills/reviewer.md) @ 1.3.0 · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad) -->
**Modell:** claude-opus-4-8 · **Datum:** 2026-08-03

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [slice-091](../plan/planning/in-progress/slice-091-slice-status-feld-entfernen.md) (Plan)
- Baseline: [`modul-05-planning-harness.md` §Lifecycle als State Machine](../../.harness/baseline/v5.0.0/regelwerk/modul-05-planning-harness.md#lifecycle-als-state-machine)
  und die Vorlage [`slice.template.md`](../../.harness/baseline/v5.0.0/templates/docs/plan/planning/slice.template.md)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules §3.3, §3.5, §5)
- [`MR-013`](../../harness/conventions/MR-013-lifecycle-move-buendelung.md) (Lifecycle-Move-Adaption)
- keine `DC-*`/`LH-*`-Anforderung und keine ADR berührt (reiner Harness-/Prozess-Slice, GF)

---

## Findings

### F-1 — Zweite „Status-Zeile"-Erwähnung in `AGENTS.md` §3.3 nicht nachgezogen

- `kategorie`: MEDIUM
- `quelle`: Hard Rule ([`AGENTS.md`](../../AGENTS.md) §3.3 ↔ §5), Ziel-Konsistenz zur kanonischen
  Quelle [`MR-013`](../../harness/conventions/MR-013-lifecycle-move-buendelung.md)
- `pfad`: `AGENTS.md:104`
- `befund`: D-5 hat die „Status-Zeile"-Erwähnung nur in
  [`MR-013`](../../harness/conventions/MR-013-lifecycle-move-buendelung.md) (Zeile 22–24, jetzt
  „Slice-Body (DoD-Haken + Closure-Notiz; neue Slices ohne Status-Zeile)") und die neue §5-Regel
  („kein `**Status:**`-Feld") gesetzt, aber die **kanonische Spiegel-Stelle** `AGENTS.md:104`
  („Nur der **Slice-Body** (Status-Zeile, Closure-Notiz) bleibt Commit 2") unverändert gelassen.
  `AGENTS.md:106–107` deklariert `MR-013` ausdrücklich als „Kanonisch"; die Zusammenfassung ist
  damit weniger aktuell als ihre eigene Quelle und widerspricht der frisch ergänzten §5-Regel
  desselben Dokuments.
- `verifizierbar`: nein (kein Gate liest diese Prosa; belegbar per `grep -n 'Status-Zeile' AGENTS.md`
  gegen `harness/conventions/MR-013-lifecycle-move-buendelung.md` — die zwei Kanonisch-Spiegel divergieren).
- `klasse`: „Kanonisch-Spiegel divergiert — Regel in einer von zwei gekoppelten Stellen nachgezogen"

**Failure-Szenario:** Ein Implementer liest bei einer künftigen Closure die Hard Rule
`AGENTS.md` §3.3 (laut `CLAUDE.md` das zuerst gelesene Doc), entnimmt „der Slice-Body enthält eine
Status-Zeile", legt für einen neuen Slice eine Status-Zeile in Commit 2 an — und verletzt damit die
drei Zeilen höher hinzugefügte §5-Regel „Slice-Pläne tragen kein `**Status:**`-Feld". Das Doc
widerspricht sich intern; die von §2/§4 des Slice übernommene Aufgabe „die Status-Zeile-Erwähnung
nachziehen" ist auf einer von zwei identischen Stellen unerledigt.

## Negativbefunde

- geprüft, ohne Befund: **Zwei `Status:`-Achsen nicht verwechselt** —
  `git diff 737ee50..HEAD -- harness/conventions/` entfernt **keine** `**Status:**`-Feldzeile;
  [`MR-013`](../../harness/conventions/MR-013-lifecycle-move-buendelung.md):3 trägt weiter
  `**Status:** Accepted`. Keine Datei unter `docs/plan/adr/` im Diff; `make adr-check` grün
  (Image-Digest `sha256:7cda81b5…`, 0 Befunde) → ADR-`**Status:**` immutable unberührt.
- geprüft, ohne Befund: **slice-091-Selbstkonsistenz** — die Datei trägt **kein**
  `**Status:**`-Feld-Header, sondern den `**Lifecycle:**`-Hinweis (Zeile 3–6); §3 Abnahme-Punkt
  ist auf „Entschieden 2026-08-02: template-forward-only" gestellt, §1/§2/§4 durchgängig
  template-forward — **kein** Rest, der noch „retrofit-all" als gewählten Weg behauptet
  (nur als verworfene Alternative benannt).
- geprüft, ohne Befund: **§5-Regel vs. Baseline** — die neue `AGENTS.md`:194–196-Regel deckt sich
  mit [`modul-05` §Lifecycle](../../.harness/baseline/v5.0.0/regelwerk/modul-05-planning-harness.md#lifecycle-als-state-machine)
  („Der Zustand ist das Verzeichnis, nicht ein Kopffeld … ein `Status:`-Feld wäre eine zweite
  Quelle") und mit [`slice.template.md`](../../.harness/baseline/v5.0.0/templates/docs/plan/planning/slice.template.md)
  (führt tatsächlich kein Status-Feld, sondern den `**Lifecycle:**`-Hinweis). Kein Dublett in §5
  (`grep 'Status-Feld' AGENTS.md` liefert nur die neue Stelle), keine Kollision.
- geprüft, ohne Befund: **MR-013-Änderung bricht keinen Verweis** — der Edit ändert nur Prosa im
  Adaptions-Absatz (Zeile 22–24); alle Markdown-Links (`AGENTS.md` §3.3-Anker, `modul-05`-Anker,
  `MR-000`) unverändert. Doc-Selbstscan (`links`/`anchors`/`ids`) 0 Befunde.
- geprüft, ohne Befund: **Faktentreue „kein Gate/Skript liest das Slice-Status-Feld"** —
  `grep -rn 'Status' internal/ tools/` findet nur ADR-Status-Leser (`report`, `git`, `httpcheck`,
  `vcs`), die RTM-Spalte „Status" (`cli_acceptance_test`) und den Roadmap-Aktiv-Status
  (`diagnose.go`, `planning`) — **keinen** Leser eines Slice-`**Status:**`-Kopffelds. Umfang „90
  `done/`-Slices tragen `**Status:**`" per `grep -rl` bestätigt (exakt 90); außer slice-091 (und der
  Roadmap-Prosa) trägt kein `in-progress/`/`next/`/`open/`-Slice das Feld.
- geprüft, ohne Befund: **Gate-Sicherheit / kein Code-Touch** — `git diff --name-only 737ee50..HEAD`
  berührt nur `AGENTS.md`, `roadmap.md`, `slice-091`, `MR-013`; kein `spec/`, `internal/`, `tools/`.
  Doc-Selbstscan `320 Datei(en), 0 Befund(e)`; `make planning-check` 0 (Roadmap-Flip §Aktuelle Welle
  auf slice-091-aktiv koppelt korrekt zum nicht-leeren `in-progress/`); Working-Tree clean.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Kanonisch-Spiegel divergiert — Regel in einer von zwei
gekoppelten Stellen nachgezogen"

## Verdikt

**nicht abnahmereif** — **Merge-blockierend: ja** (1 MEDIUM). Der Kern der Umsetzung ist tragfähig
und faktentreu (zwei `Status:`-Achsen sauber getrennt, kein Gate liest das Feld, slice-091 modelliert
die Ziel-Form, alle Gates grün). Blockierend ist allein F-1: die go-forward-Verankerung ist
unvollständig, weil die **kanonische Spiegel-Stelle** `AGENTS.md` §3.3 dieselbe „Status-Zeile"-Aussage
weiterführt, die `MR-013` bereits abgelöst hat, und damit `AGENTS.md` §3.3 gegen §5 in sich
widersprüchlich ist. Der Fix ist eng lokalisiert (eine Zeile, `AGENTS.md:104`). Ohne ihn erbt die
adoptierte Baseline-Form eine interne Regel-Kollision im zuerst gelesenen Hard-Rules-Doc — bei einem
Slice, dessen erklärter Zweck genau das Nachziehen dieser Erwähnung ist.
