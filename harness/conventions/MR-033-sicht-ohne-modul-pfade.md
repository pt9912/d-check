# MR-033 — Die Architektur-Sicht führt auch keine Modul-Pfade

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`AGENTS.template.md`](../../.harness/baseline/v6.0.0/templates/AGENTS.template.md)
  §3.4 — *„`spec/architecture.md` **referenziert Modul-Pfade**, aber keine
  Wellen, Slices, Commit-Hashes oder Closure-Daten."* Dieselbe Erlaubnis trägt
  [`modul-03-spec.md`](../../.harness/baseline/v6.0.0/regelwerk/modul-03-spec.md)
  §Ziel-Form: Architektur-Sicht (*„Sprach- und meilensteinfrei — referenziert
  Modul-Pfade, aber …"*).
- **Datum:** 2026-08-23
- **Geltungsbereich:** [`spec/architecture.md`](../../spec/architecture.md) und
  die Regel [`AGENTS.md`](../../AGENTS.md) §3.4, die sie trägt.
- **Adaption:** Der Kanon **erlaubt** der Sicht Modul-Pfade; dieses Repo
  **verbietet** sie ihr. `AGENTS.md` §3.4 schreibt *„benennt Schichten und
  Rollen statt Technologie — keine Sprach-/Modul-Pfade"*. Das ist eine
  Verschärfung, und sie war bis zu diesem Eintrag **undeklariert**.

  **Warum sie ein Eintrag sein muss:**
  [`MR-031`](../conventions.md#mr-031) sagt es für genau diese Klasse — *wer
  verschärft, weicht ab, auch wenn er nur mehr verlangt.* Ohne Eintrag ist die
  Abweichung für den Freshness-Audit unsichtbar.

  **Begründung der Strenge:** Die Sicht ist das derivative Stratum; ein
  Modul-Pfad in ihr ist eine Aussage über die **Umsetzung**, und die gehört
  eine Schicht tiefer. Nennt die Sicht Pfade, wandert sie bei jedem
  Verzeichnis-Umbau mit — und das ist genau die Kopplung, die das Stratum
  vermeiden soll. Die sprachkonkrete Übersetzung leben die ADRs, deren
  `Schärft:`-Feld aufwärts zeigt.

  **Was die Adaption *nicht* betrifft:** die vier Meilenstein-Kategorien
  (Wellen, Slices, Commit-Hashes, Closure-Daten) und den ADR-Bezug — die
  verbietet der Kanon ohnehin, dort verschärfen wir nichts.
- **Ausgelöst durch Baseline-Stand:** v5.11.0
- **Auflösungs-Trigger:** der Kanon verbietet der Sicht Modul-Pfade selbst —
  dann ist die Adaption gegenstandslos. Bis dahin gilt sie.
