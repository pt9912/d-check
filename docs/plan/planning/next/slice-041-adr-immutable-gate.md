# Slice slice-041: ADR-Immutable-Gate (Accepted-ADRs nicht still ändern)

**Status:** open (geplant — noch nicht priorisiert).

**Welle:** — (offen; noch keiner Welle zugeordnet).

**Bezug:** Hard Rule `AGENTS.md` §3.5 (eine `Accepted`-ADR wird nicht
inhaltlich überschrieben; Korrekturen via neue ADR mit `Supersedes`).
Schwester-Lücke zu slice-040; dieselbe „aspirativ-vs-bindend"-Klasse wie das
Traceability-Gate ([ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)).

**Autor:** pt9912. **Datum:** 2026-06-21.

---

## 1. Ziel

Ein Wächter, der **inhaltliche Änderungen an `Accepted`-ADRs** erkennt und
ablehnt — erlaubt bleiben nur reine Anhänge unter `## Geschichte` und der
Status-Übergang nach `Superseded by …`. AGENTS §3.5 ist heute dokumentiert,
aber nicht erzwungen.

## 2. Zu entscheiden (im Slice)

- **Mechanik:** CI-Range-Check (wie das Traceability-Gate) — im
  PR/Push-Diff jede geänderte `docs/plan/adr/NNNN-*.md` mit `Status:
  Accepted` flaggen, außer der Diff betrifft nur die `## Geschichte`-Tabelle
  bzw. den Status-Wechsel nach `Superseded`. Alternativen: Hash-Baseline pro
  Accepted-ADR (schwächer, mit-editierbar) oder git-history-Vergleich gegen
  den Accept-Commit (komplex).
- **Granularität:** zeilengenaue Diff-Klassifikation (Geschichte/Status vs.
  Körper) — wie robust?
- **Carveout/Grandfathering:** vor-Einführung bestehende Accepted-ADRs sind
  Baseline (nur ab Einführung geprüft) — analog dem Traceability-Gate.

## 3. Definition of Done (vorläufig)

- [ ] Mechanik entschieden; Wächter (CI-Range-Check und/oder lokales Target)
  implementiert; deterministisch, read-only.
- [ ] Negativ-Selbsttest: eine simulierte Körper-Änderung an einer
  `Accepted`-ADR feuert; ein reiner `## Geschichte`-Anhang feuert **nicht**.
- [ ] Doku-Sync (`harness/README.md` §Sensors / `AGENTS.md` §4 falls
  Make-Target); `make gates` grün; Review R1.
- [ ] ADR nach Bedarf (Enforcement-Topologie, falls CI-seitig wie beim
  Traceability-Gate).

## 4. Risiken / offene Punkte

- **Diff-Klassifikation** (Geschichte/Status vs. Körper) ist die Kernschwierigkeit
  — eine zu grobe Heuristik erzeugt Falsch-Positive bei legitimen Anhängen.
- **Bootstrap:** der Check greift erst ab Einführung (Grandfathering der
  Bestands-ADRs), sonst flaggt der erste Lauf historische Edits.

## 5. Trigger

Nutzer-Audit 2026-06-21: AGENTS §3.5 dokumentiert, aber nicht maschinell
erzwungen. Aus slice-040 (Planning-Konsistenz) als eigener Folge-Slice
ausgekoppelt (andere Mechanik: ADR-Diff statt Verzeichnis-State).

## 6. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Harness-Mechanik/Doku; Greenfield-Default).
