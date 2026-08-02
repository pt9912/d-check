# MR-014 — Slice-/ADR-Doc-Struktur: Repo-Haus-Stil ggü. Baseline-Template

- **Status:** Accepted
- **Aufgelöst durch:** Baseline-Stand v4.0.0 (Doc-Form ist Baseline-Wahl, ADR-Alternativen-Tabelle wurde Default)
- **Datum:** 2026-06-22
- **Geltungsbereich:** Slice-Dateien unter `docs/plan/planning/` und
  ADR-Dateien unter `docs/plan/adr/` (Doc-**Struktur**, nicht Inhalt)
- **Adaption:** Slice- und ADR-Dateien folgen einem repo-lokalen Haus-Stil,
  der von den adoptierten Baseline-Vorlagen
  ([`slice.template.md`](https://github.com/pt9912/ai-harness-course/blob/v1.3.0/lab/templates/docs/plan/planning/slice.template.md),
  [`NNNN-titel.template.md`](https://github.com/pt9912/ai-harness-course/blob/v1.3.0/lab/templates/docs/plan/adr/NNNN-titel.template.md))
  in der **Abschnitts-Struktur** abweicht; Header-Block und Pflichtinhalte
  bleiben deckungsgleich:
  - **Slice:** keine separate „Plan (vor Code)"-Tabelle (der Plan steht im
    Entscheidungen-/Regel-Block, §2); kein eigener Abschnitt „Closure-Trigger"
    (in §5 Trigger bzw. der DoD aufgegangen); Reihenfolge §1 Ziel · §2
    Entscheidungen/Regel · §3 DoD · §4 Risiken · §5 Trigger · §6
    Sub-Area-Modus · §7 Closure-Notiz (das Template führt Sub-Area als §8).
  - **ADR:** „Verglichene Alternativen" als **Tabelle** (Pro/Contra) statt
    Option-A/B/C-Prosa; „Fitness Function" als Prosa-Bullets statt der
    Tooling/Regel/Target-Tabelle; „Geschichte" zweispaltig (Datum/Ereignis)
    statt dreispaltig (der Verweis steht im Ereignis-Text).
- **Begründung:** Die Vorlagen sind ausdrücklich Start-Vorlagen („Kopiere …
  ersetze Platzhalter, lösche den Hinweis-Block"); strukturelle Anpassung ist
  erwartbar. Der Haus-Stil ist über die bisherigen Slices und ADRs konsistent
  gelebt; Konsistenz mit dem Bestand sticht die buchstäbliche Template-Form.
  Sichtbar geworden im Template-Vergleich (Nutzer-Frage 2026-06-22,
  lab-templates v1.3.0) — bis dahin eine **stille Setzung** (gleiche Klasse wie
  eine undeklarierte Konvention). Inhaltlich bindende Elemente bleiben
  unberührt: die Header-Felder **Bezug**/**Schärft**, die Akzeptanz-/DoD-
  Disziplin, die Sub-Area-Modus-Pflicht und die `## Geschichte`-Anhang-Regel
  der ADR-Immutability ([ADR-0016](../../../docs/plan/adr/0016-adr-immutable-gate.md)).
  [`MR-000`](../../conventions.md#mr-000--baseline-aussage) bleibt unverändert; Nachtrag daher als
  eigener Eintrag, analog
  [`MR-008`](../../conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage).
- **Auflösungs-Trigger:** permanent (Doc-Form ist Repo-Wahl; Baseline: „Form
  ist Wahl"). Bei künftiger Baseline-Hebung wird die Struktur-Differenz hier
  re-evaluiert (analog [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt)).
