# MR-024 — `## Aktuelle Welle`: Ruhe-Marker-/Prosa-Form statt Template-Struktur-Felder

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`modul-06-roadmap.md` §Roadmap-Struktur (Aktuelle Welle)](../../.harness/baseline/v5.0.0/regelwerk/modul-06-roadmap.md#roadmap-struktur-fünf-abschnitte-modul-6)
- **Datum:** 2026-08-02
- **Geltungsbereich:** [`roadmap.md` §Aktuelle Welle](../../docs/plan/planning/in-progress/roadmap.md#aktuelle-welle),
  `make planning-check`
  ([`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) /
  [ADR-0028](../../docs/plan/adr/0028-planning-lifecycle-modul.md))
- **Adaption:** Das Baseline-Template führt `## Aktuelle Welle` mit **Struktur-Feldern**
  (`**Welle-ID:**` / `**Start:**` / `**Geplantes Ende:**` / `**Closure-Trigger:**`) und
  setzt **immer** eine laufende Welle voraus. d-check weicht in zwei Punkten ab:
  - **(a) Prosa + Verweis statt Feld-Duplikat.** Die drei Pflicht-Bestandteile
    (Slice-IDs · Trigger · Closure-Kriterien) stehen **in Prosa** im Abschnitt; die
    Feld-Details (Welle-ID/Start/Geplantes Ende) liegen im **Wellendokument**
    (`docs/plan/planning/welle-<NN>-…md`, seit slice-088). Die Roadmap bleibt
    **Sequenzierungs-Autorität**, das Wellendokument trägt die Welle-Details — keine
    Zweitquelle.
  - **(b) Expliziter Ruhe-Marker für den wellenlosen Zustand.** Wo das Template keinen
    Zustand „keine laufende Welle" kennt, macht d-check ihn **explizit** über den Marker
    „**Keine aktive Welle.**". Er ist **gate-erzwungen**: kein `slice-*` in
    `in-progress/` ⟺ Marker im Abschnitt (das `planning`-Modul prüft den Invariant
    hermetisch).
  Nachtrag zu [`MR-013`](../conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)
  (nutzt den Ruhe-Marker im Move-Commit-Kontext) und
  [`MR-014`](../conventions.md#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template)
  („Doc-Form ist Repo-Wahl");
  [`MR-000`](../conventions.md#mr-000--baseline-aussage) bleibt unverändert.
- **Begründung:** d-checks Identität ist **mechanisieren, wo die Baseline beim Menschen
  bleibt**: der wellenlose Zustand ist im Regelwerk (modul-06 §Wann Arbeit eine Welle
  braucht — „wellenlose Arbeit erscheint nicht in der Roadmap") nur implizit; d-check
  macht ihn seit slice-057 (`planning`-Modul) maschinell prüfbar. Die Feld-Details ins
  Wellendokument zu legen statt in die Roadmap folgt derselben „eine Quelle"-Linie wie
  die Lifecycle-Adaption in `MR-013`.
- **Auflösungs-Trigger:** permanent, solange `make planning-check` den
  Roadmap-↔-in-progress-Invariant erzwingt.
