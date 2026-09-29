# Review slice-240 — Baseline-Pin-Hebung `v6.9.0` → `v6.13.0` (R1)

- **Review-Art:** Code/Doku-Review — geprüft gegen Slice-Plan `slice-240` (§1 Ziel und Abgrenzung, §3 Plan, §6 Risiken), die Hard Rules (`AGENTS.md` §3) und die Konventionen (`MR-051`, `MR-055`, `MR-069`, `MR-070`); **nicht** gegen die DoD (Verifikation, getrennter Kontext).
- **Gegenstand:** `604bc25d..d8a7baf2` — 3 Commits: `60aa9b18` (reiner Move MR-072 → `conventions/done/`, 100 % Similarity), `94b80f35` (MR-072 Link-Tiefen ±9), `d8a7baf2` (Baum-Swap 55/55 Dateien, §Baseline-Pin, MR-073, MR-Index, 8 Symlinks, 4 `ignore-refs`-Einträge, Retargets in 44 lebenden Dokumenten).
- **Skill:** `reviewer.md` @ 1.16.0
- **Modell-ID:** glm-5.3-flash (Claude Agent SDK)
- **Datum:** 2026-09-29
- **Eingangs-Kontext:** Slice-Plan `slice-240`; `AGENTS.md` §3; `harness/conventions.md` §Baseline samt aktiven Adaptionen (`MR-021`, `MR-039`, `MR-051`, `MR-055`, `MR-069`, `MR-070`, `MR-072`); `v6.13.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill; vorherige Findings am gleichen Vorgang: `BEO-ALL/pin-bump-mirrors-ungated` (7 Vorkommen, u. a. Review-F-1 des Vorgängers — unerkannte ninth frozen-Klasse) und MR-072s eigener Zitat-Delta-Vermerk.

---

## Findings

### F-1 — MEDIUM

- **Kategorie:** MEDIUM
- **Quelle:** `MR-051` (Re-Ankern-Pflicht) · Anker „Botschaft verallgemeinert über die Messung hinaus" (`BEO-ALL/commit-message-overclaims-work`(b), `AGENTS.md` §5)
- **Pfad:** `harness/conventions/MR-073-baseline-v6130.md:67-73`
- **Befund:** MR-073 behauptet „**Kein `d-check:cite`-Neu-Ankern**, kein Zitat-Delta — anders als beim Vorgänger" und begründet das mit genau zwei Spannen (denen des slice-240-Plans). Der Diff zeigt **neun** neu geankerte Spannen: die zwei Plan-Spannen (`modul-05:363-364→373-374`, `369→379`) plus **sieben weitere in vier Dateien** — `.harness/skills/reviewer.md` (`modul-05:373-374`, `modul-05:379`, `grundlagen-durchsetzungsschicht.md:66`), `harness/conventions/MR-031-schritt-3-benennen.md` (`modul-09:203`), `harness/conventions/MR-005-haertung-gate-nachweis.md` (`grundlagen-durchsetzungsschicht.md:116-118`), `harness/conventions/MR-049-ausgangs-wortschatz.md` (`modul-05:180-181`, `modul-05:170`) — dazu fünf weitere Direktiven mit reiner Tag-Swap-Änderung. Das Neu-Ankern selbst ist `MR-051`-konform und inhaltlich verifiziert (die Ziellinien 170 und 180-181 tragen die zitierten Sätze wörtlich); falsch ist die **Abwesenheits-Behauptung** im Konventionsspeicher — dieselbe Klasse, die der Vorgänger-Review als F-1 am selben Vorgang fand (dort: eine neunte frozen-Klasse, erst im Review deklariert).
- **Verifizierbar:** ja — `git diff 604bc25d..d8a7baf2` gefiltert auf `d-check:cite` zeigt die neun Spannen-Änderungen; `make doc-check` bestätigt die *Substanz* (alle Spannen lösen korrekt auf), fängt aber nicht die Falschbehauptung im MR-Text.
- **Klasse:** record-claim-vs-diff (Anker 8: N+1-te Form — „kein Neu-Ankern" widerlegt durch die siebte Spanne)

### F-2 — LOW

- **Kategorie:** LOW
- **Quelle:** `MR-070` (Frozen-Klassen vor mechanischer Ersetzung — Eigenschaft, nicht Verzeichnis) · Anker „Quelle über ihren Geltungsbereich hinaus zitiert" (hier: die frozen-Liste als Datei-Eigenschaft gelesen)
- **Pfad:** `harness/conventions/MR-073-baseline-v6130.md:47-56`
- **Befund:** MR-073 zählt MR-056 („Zitat-Delta-Vermerk") und `spec/spezifikation.md` („2 Spec-Historie-Zeilen-Träger") zu den 29 Dateien, die „**ganz frozen**" blieben. Beide sind tatsächlich Mischfundstellen: MR-056s Baseline-Link und `spezifikation.md`'s §Rolle-Pointer wurden je in einer Zeile auf `v6.13.0` retargetet (nur die Historie-/Zitat-Prosa blieb bei `v6.9.0`). Die zeilenweise Trennung selbst ist korrekt ausgeführt — die Klassenzugehörigkeit im messenden Eintrag ist falsch; die 36/8/29-Zählung bleibt nur deshalb arithmetisch stimmig, weil beide Dateien ihre `v6.9.0`-Vergangenheits-Aussage behalten haben.
- **Verifizierbar:** ja — `git diff 604bc25d..d8a7baf2 -- harness/conventions/MR-056-dod-haken-waechter.md spec/spezifikation.md` zeigt je eine Retarget-Zeile bei als frozen gelisteten Dateien.
- **Klasse:** frozen-list-misclassification (dieselbe Fundstellen-Klasse wie F-1, gleiche Auflösungsstelle)

---

## Negativbefunde (geprüft, ohne Befund)

1. **Vendorter Baum + Integrität:** 55 Dateien vor wie nachher (Pfade 1:1, `modul-00`/`modul-04` als Neu-/Entfall-Paare unter 50 %-Similarity korrekt erkannt); `SHA256SUMS` (54 Einträge) mit `sha256sum -c` Exit 0 verifiziert. Die Delta-Messung „28 (19 Regelwerk, 8 Templates, SHA256SUMS)" ist methodisch sauber: Stichprobe `grundlagen-begriffe.md` zeigt, dass reine `<!-- Quelle: -->`-URL-Swaps korrekt als delta-frei zählen. MR-073s vier inhaltliche Schlagzeilen sind als Abgrenzung (nicht Adoption) korrekt geführt.
2. **Frozen-Träger unberührt:** 13 `done/`-Slices (inkl. Stub slice-224), 2 Review-Reports, 3 aufgelöste MR-Dateien, CO-002, CR 2026-09-17, Evidence slice-225, ADR-0085, `spec/lastenheft.md`, 3 Code-Dateien unter `internal/hexagon/core/`, `.d-check.closure.yml` — keiner davon im Diff.
3. **Mischfundstellen-Behandlung (die 7 übrigen):** `AGENTS.md`, `harness/README.md`, `harness/conventions.md`, `roadmap.md`, `observations/README.md`, `MR-021`, `MR-049`, `MR-053` — durchgehend zeilenweise getrennt, kein Über-Swap auf frozen Aussagen; der MR-072-Index-Eintrag (Risiko 4 des Plans) wurde korrekt als aufgelöste Zeile mit Vergangenheits-Aussage „auf v6.9.0" neu gebaut statt überschrieben. Roadmaps neu geankerte `§Roadmap-Struktur`-Aussage steht zeilenidentisch in `v6.13.0` (`modul-06-roadmap.md`:62, im Delta unberührt).
4. **`ignore-refs`-Wachstum (`MR-069`-Klasse):** Genau **vier** Markdown-Links auf den entfernten `v6.9.0`-Baum gemessen (grep-verified: ADR-0085, MR-059, MR-061, MR-062) — die Begründungslast trägt der Datei-Kommentar in `.d-check.yml` (Präzedenz ADR-0084, `MR-052`-Historie-Rückgriff, Nicht-Feuern der übrigen frozen Nennungen) und MR-073. Keine undeklarierte Gate-Senkung (`AGENTS.md` §3.6 — das Ventil ist per MR-069 selbst deklariert).
5. **MR-072-Zug (`AGENTS.md` §3.3):** Move-Commit `60aa9b18` ist 0-Zeilen-rein (100 % Similarity), die Link-Tiefen-Korrektur ±9 in `94b80f35` enthält ausschließlich Tiefenanpassungen — Regelfall 1 sauber zerlegt.
6. **Symlinks (`MR-055`):** 8 gemessen statt 7 geplant, Abweichung in MR-073 dokumentiert; alle 8 lösen gegen `v6.13.0` auf.
7. **Zitat-Delta (`MR-039`):** `AGENTS.md`-Zitat „Halluzinierte Gates …" wörtlich in `v6.13.0` · `regelwerk/modul-13-quality-gates.md`:97; MR-056s Zitatzeile ist in `v6.9.0`:36 und `v6.13.0`:36 identisch — „kein Zitat-Delta" verifiziert.
8. **Spec-Straten (`MR-006`, `AGENTS.md` §3.4):** Die beiden Spec-Retargets (`architecture.md`, `spezifikation.md`) sind §Rolle-Pointer; keine ADR-/Wellen-/Slice-/Commit-Token geflossen (korroboriert durch grünes `make doc-check` 885/0 laut Implementer-Bericht).
9. **Hard Rules §3.1/§3.2/§3.9:** Kein Go-/Host-Tool-Pfad, keine Inline-Suppression, keine Workflow-Änderung — nichts geprüft, nichts beanstandet.
10. **Gates:** `make gates` grün, `doc-check` 885/0, `baseline-freshness` Exit 0 — laut Implementer-Bericht, hier nicht wiederholt (Verifier-Rolle); unabhängig davon wurde die Baseline-Integrität per `sha256sum -c` nachgeprüft.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 |
| LOW | 1 | F-2 |
| INFO | 0 | — |

## Verdikt

**MEDIUM blockiert typischerweise — hier mit engem Auflösungspfad.** Der mechanische Vorgang selbst (Swap, Pin, Retargets, Ventil, Move) ist sauber, messrichtig und durchgehend gate-konform; beide Findings liegen im **Neu-MR-Eintrag** (MR-073), nicht im Vollzug, und lösen sich in einer einzigen Doku-Korrektur: MR-073s Absatz „Kein `d-check:cite`-Neu-Ankern" auf die gemessenen neun Spannen korrigieren und MR-056 samt `spec/spezifikation.md` aus der „ganz frozen"-Liste in die Mischfundstellen-Klasse verschieben (36/8/29 → 34/10/27). Das gehört vor dem Closure-Move in den Feat-Commit-Kontext; danach ist der Slice review-seitig frei. Die wiederkehrende Klasse ( MR-Eintrag behauptet Abwesenheit, der Diff zählt Vorkommen — zweite Instanz am selben Vorgang nach dem Vorgänger-Review) bleibt über `BEO-ALL/pin-bump-mirrors-ungated` beobachtet; eine dritte Wiederholung wäre ein Steering-Loop-Signal.
