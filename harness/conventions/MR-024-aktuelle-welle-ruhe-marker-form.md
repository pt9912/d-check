# MR-024 — `## Aktuelle Welle`: Ruhe-Marker im wellenlosen Zustand (aktive Welle template-konform)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`modul-06-roadmap.md` §Roadmap-Struktur (Aktuelle Welle)](../../.harness/baseline/v5.0.0/regelwerk/modul-06-roadmap.md#roadmap-struktur-fünf-abschnitte-modul-6)
- **Datum:** 2026-08-02
- **Geltungsbereich:** [`roadmap.md` §Aktuelle Welle](../../docs/plan/planning/in-progress/roadmap.md#aktuelle-welle),
  `make planning-check`
  ([`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) /
  [ADR-0028](../../docs/plan/adr/0028-planning-lifecycle-modul.md))
- **Adaption:** Das Baseline-Template führt `## Aktuelle Welle` mit **Struktur-Feldern**
  (`**Welle-ID:**` / `**Start:**` / `**Geplantes Ende:**` / `**Closure-Trigger:**`) und
  setzt **immer** eine laufende Welle voraus. d-check **folgt dieser Feld-Form, solange
  eine Welle läuft** (seit slice-092: die Struktur-Felder stehen im Abschnitt, das
  Wellendokument trägt die weiteren Details). Die **einzige** verbleibende Abweichung ist
  der **wellenlose** Zustand, den das Template nicht kennt:
  - **Expliziter Ruhe-Marker.** d-check macht „keine laufende Welle" **explizit** über den
    Marker „**Keine aktive Welle.**". Er ist **gate-erzwungen**: kein `slice-*` in
    `in-progress/` ⟺ Marker im Abschnitt (das `planning`-Modul prüft den Invariant
    hermetisch). Läuft eine Welle, trägt der Abschnitt stattdessen die Struktur-Felder
    **ohne** Marker (`hasActive == hasSlices` grün) — deshalb war **kein** `planning`-Modul-
    Umbau nötig, um die Template-Form zu erreichen.

  **Nachtrag slice-092:** die frühere Prosa-Form (die drei Pflicht-Bestandteile als
  Fließtext statt Felder) ist **abgelöst** — der Aktiv-Fall ist jetzt template-konform.
  Nachtrag zu [`MR-013`](../conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)
  (nutzt den Ruhe-Marker im Move-Commit-Kontext);
  [`MR-000`](../conventions.md#mr-000--baseline-aussage) bleibt unverändert.
- **Begründung:** d-checks Identität ist **mechanisieren, wo die Baseline beim Menschen
  bleibt**: der wellenlose Zustand ist im Regelwerk (modul-06 §Wann Arbeit eine Welle
  braucht — „wellenlose Arbeit erscheint nicht in der Roadmap") nur implizit; d-check
  macht ihn seit slice-057 (`planning`-Modul) maschinell prüfbar. Die Feld-Details ins
  Wellendokument zu legen statt in die Roadmap folgt derselben „eine Quelle"-Linie wie
  die Lifecycle-Adaption in `MR-013`.
- **Auflösungs-Trigger:** permanent, solange `make planning-check` den
  Roadmap-↔-in-progress-Invariant erzwingt.
