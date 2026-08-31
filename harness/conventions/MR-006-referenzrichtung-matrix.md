# MR-006 — Referenzrichtung: Spec-Straten verweisen nie abwärts auf ADRs

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)](../../.harness/baseline/v5.15.0/regelwerk/grundlagen-referenz-richtung.md#referenz-richtung-sdp-wer-darf-wen-referenzieren)
- **Datum:** 2026-06-10
- **Geltungsbereich:** `spec/*.md`, [`AGENTS.md` §3.4](../../AGENTS.md#34-architektur-sprach-meilensteinfrei-spec-straten-nie-abwärts)
- **Adaption:** Das adoptierte Template-Set 2026-06 sah ADR-Verweise
  in `spezifikation.md` (ADR-Spalte in Defaults/Historie) und
  `architecture.md` (ADR-Spalte, §ADR-Index) vor. Das ist als Fehler
  der Kurs-Vorlagen identifiziert; die Korrektur erfolgt in der
  Kurs-Quelle (Entscheidung User, 2026-06-10). d-check zieht vor:
  **kein Spec-Stratum (Rang 1–3) verweist abwärts auf ADRs oder
  Planning-Artefakte**; Traceability läuft ausschließlich über die
  `Schärft:`-Felder der ADRs (aufwärts). Die spätere
  matrix-Selbstkonfiguration kodiert das als
  `{from: spec-strata, to: adr/slice, allow: false}`. **Scope-Grenze gegen die
  v5.0.0-8×8-Matrix (C-4):** d-checks `matrix` bewacht die **Spec-Decke** (alle
  drei Spec-Straten, seit `v4.0.0` Baseline-Default) und die markierte
  **ADR→Slice**-Kante; die weiteren ❌-Kanten der 8×8-Matrix —
  **ADR→Carveout/Welle/Roadmap** und **Slice→Roadmap** — sind **bewusst
  unbewacht** (d-check modelliert Carveout/Welle/Roadmap nicht als
  `matrix`-Klassen; eine Erweiterung wäre ein eigener Change).
- **Begründung:** Stable Dependencies — die Lösungsbeschreibung muss
  Entscheidungs-Revisionen (Supersede) überleben, ohne selbst
  angefasst zu werden; die Richtung der Begründung ist ADR → Spec,
  nie umgekehrt. Konsistent mit u-boots Checker („view spec may not
  link down").
- **Auflösungs-Trigger:** permanent, solange d-check die Referenz-Richtung via
  `matrix`-Modul mechanisiert und dabei die **C-4-Scope-Grenze** trägt (nur
  Spec-Decke + markierte ADR→Slice-Kante, nicht die volle 8×8-Matrix). *Die
  ursprüngliche Adaption — das Vorziehen der Spec-Decke-Regel vor die Kurs-Quelle —
  ist eingetreten (Spec-Straten-Vorlagen seit `v1.3.0` korrigiert, Spec-Decke seit
  `v4.0.0` Baseline-Default) und bleibt als Provenienz.*
