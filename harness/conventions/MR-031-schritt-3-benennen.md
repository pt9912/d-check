# MR-031 — Schritt 3 des Agenten-Workflows verlangt Benennen, nicht nur Identifizieren

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`modul-09-implementierung.md`](../../.harness/baseline/v5.18.0/regelwerk/modul-09-implementierung.md)   <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-09-implementierung.md:20-20 -->
  §Minimal Agent Workflow, Schritt 3 — *„Betroffene Requirement-/ADR-IDs
  identifizieren."*
- **Datum:** 2026-08-23
- **Geltungsbereich:** [`AGENTS.md`](../../AGENTS.md) §6 Schritt 3; jeder
  Implementer-Lauf in diesem Repo.
- **Adaption:** Schritt 3 verlangt hier zusätzlich, das Ergebnis **zu
  benennen** statt es nur zu identifizieren, und nennt vier Felder mehr:
  Slice-ID, betroffene `DC-*`-IDs, ADR-IDs, betroffene Module und die
  auszuführenden Gates.

  **Warum das eine Adaption ist und keine Ergänzung:** Der Kanon führt Schritt 3
  als **benannten** Schritt einer nummerierten Folge. Wer ihn verschärft,
  weicht von einem Baseline-Default ab — auch wenn er nur *mehr* verlangt.
  Ohne Eintrag hier hätte der Freshness-Audit nichts zu prüfen und die
  Abweichung wäre beim nächsten Baseline-Bump unsichtbar.

  **Was die Verschärfung trägt:** Der Kanon nennt die *minimalen Eingaben eines
  Implementer-Agenten gegen Halluzination* und stellt fest, dass fehlende   <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-09-implementierung.md:166-166 -->
  Eingaben *„durch Raten ersetzt"* werden. Er verlangt in Schritt 8 den Bericht
  über **gelaufene** Sensors, aber nirgends die **Vorab-Nennung** der Gates,   <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/grundlagen-harness-dateien.md:187-187 -->
  die laufen sollen — obwohl er halluzinierte Gates selbst als *„die häufigste
  Form von Harness-Lüge"* führt. Diese Asymmetrie schließt der Eintrag
  repo-lokal.

  **Grenze, benannt statt zugesagt:** Die Verschärfung liegt allein im
  Feedforward-Quadranten. Kein Gate prüft, ob ein Lauf die fünf Felder genannt
  hat — die Nennung ist Prosa, und eine Prosa-Pflicht ohne Wächter bleibt
  *halb durchgesetzt*
  ([`modul-09`](../../.harness/baseline/v5.18.0/regelwerk/modul-09-implementierung.md)
  §AGENTS.md-Regeln). Das ist hier zu wissen, nicht zu heilen: ein
  Heuristik-Wächter auf Botschafts-Text wäre ein behauptetes Gate.
- **Begründung:** Die Pflicht lebte seit dem Bootstrap in der
  Werkzeug-Einstiegsdatei und im Workflow-Skelett — beides Orte, an denen der
  Kanon seit `v5.11.0` **keine** Festlegung mehr zulässt (§Vollständigkeit).
  Sie brauchte damit einen gerankten Ort; `AGENTS.md` ist er, und die
  Abweichung vom kanonischen Schritt gehört in diesen Speicher.
- **Ausgelöst durch Baseline-Stand:** v5.11.0
- **Auflösungs-Trigger:** die Baseline verlangt die Vorab-Nennung selbst — dann
  ist die Adaption aufgelöst und `AGENTS.md` §6 folgt wieder unverändert dem
  kanonischen Schritt.
