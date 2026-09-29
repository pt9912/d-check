# Review slice-242 — Adoption „Gate-Erweiterung ist kein ADR-Anlass" (R1)

- **Review-Art:** Code/Doku-Review — geprüft gegen den Slice-Plan `slice-242` (§1 Ziel und
  Abgrenzung, §6 Risiken), die Hard Rules (`AGENTS.md` §3), die Konventionen (`MR-073`,
  `MR-051`) und das Delta `v6.13.0` · `regelwerk/modul-04-adrs.md` §Kernidee; **nicht** gegen
  die DoD (Verifikation, getrennter Kontext).
- **Gegenstand:** Commit `d19b2095` (slice-242). Der angekündigte Umfang („§3.6-Klarstellung +
  reviewer-cite Re-Anker") liegt **nicht vollständig in diesem Commit** — die §3.6-Klarstellung
  trägt der Nachbar-Commit `088c3699` mit slice-241-Beschriftung; das bestimmt F-1. Beide
  Commits sind unges pushed (`origin/main` = `e4dc7f85`).
- **Skill:** `reviewer.md` @ 1.16.0
- **Modell-ID:** glm-5.3-flash (Claude Agent SDK)
- **Datum:** 2026-09-29
- **Eingangs-Kontext:** Slice-Pläne slice-242 und slice-241; `AGENTS.md` §3/§5;
  `harness/conventions.md` (`MR-073`, `MR-051`); `v6.13.0` ·
  `regelwerk/modul-04-adrs.md` §Kernidee (Delta, Zeilen 40–48); vorherige Findings am
  gleichen Vorgang: slice-240 R1 F-1 (`record-claim-vs-diff`).

---

## Findings

### F-1 — MEDIUM

- **Kategorie:** MEDIUM — vor dem Push zu lösen (Push-Blocker in der Historie, nicht im Ist-Zustand)
- **Quelle:** Dokumentations-Regel 15
  ([`harness/rules/dokumentations-regeln/15-commit-botschaft-nicht-mehr-behaupten.md`](../../harness/rules/dokumentations-regeln/15-commit-botschaft-nicht-mehr-behaupten.md))
  · Anker 8-Familie `BEO-ALL/commit-message-overclaims-work`, erste Richtung (behauptete
  Arbeit fand in dem Commit nicht statt) — hier MEDIUM statt HIGH-nah, weil die Arbeit
  existiert und am HEAD korrekt ist; nur ihre Commit-Zuordnung ist vertauscht.
- **Pfad:** `d19b2095` (Botschaft vs. Diff) · `088c3699` (`AGENTS.md:206-210` — die
  slice-242-Ziel-Ergänzung unter slice-241-Beschriftung)
- **Befund:** Die §3.6-Klarstellung — der Ziel-Kern des slice-242 — ist im
  slice-241-Commit `088c3699` („Trigger-Audit als Closure-Schritt adoptiert
  (slice-241, MR-073, MR-051)") gelandet; sein AGENTS.md-Diff ist genau diese Ergänzung
  (1 DELETE / 6 INSERT). Der slice-242-Commit `d19b2095` selbst trägt nur den
  reviewer-cite Re-Anker (`AGENTS.md:281-281` → `286-286`), behauptet in seiner Botschaft
  aber „Gate-Erweiterung-Klarstellung in 3.6" — `git show d19b2095 -- AGENTS.md` ist leer.
  Der slice-241-Plan nennt die §3.6-Änderung nirgends. Beide Botschaften tragen IDs, daher
  bleibt `make trace-check` grün — die falsche Zuordnung fängt kein Gate, und ab Push ist
  sie dauerhaft: `git log -- AGENTS.md` würde die Gate-Erweiterung slice-241 zuschreiben,
  dessen Plan und Closure davon nichts wissen, während slice-242s Commit sie nur behauptet.
- **Verifizierbar:** ja — `git show d19b2095 --stat` (1 Datei:
  `.harness/skills/reviewer.md`); `git show 088c3699 -- AGENTS.md` (die Ergänzung);
  `git rev-parse origin/main` (= `e4dc7f85`, beide Commits lokal);
  grep nach „3.6"/„Gate-Erweiterung" im slice-241-Plan → kein Treffer.
- **Klasse:** commit-boundary-cross-slice (Botschaft-vs-Diff; verwandt mit
  `record-claim-vs-diff`, slice-240 R1 F-1)

### F-2 — INFO

- **Kategorie:** INFO
- **Quelle:** Maintainability — Plan-Risiko 2 („die Grenze wird im Vollzug am eigenen
  Bestand belegt")
- **Pfad:** `AGENTS.md:207` (Stand `d19b2095`)
- **Befund:** Der übernommene Delta-Begriff „unabhängig lauffähig" bleibt im Vollzug
  unbelegt: der §3.6-Text nennt keine Targets des eigenen Bestands, die die Klasse schon
  erfüllen, wie Plan-Risiko 2 vorgesehen hatte. Der Plan führt den Ausgang als offen und
  weist ihn in die Closure-Notiz — dort ist der Beleg zu liefern, den der §3.6-Text nicht
  hergibt; andernfalls wäre das der Anker-8-Familie (Vollzug behauptet mehr als belegt).
- **Verifizierbar:** ja — `grep -n "unabhängig lauffähig" AGENTS.md` (nur der adoptierte
  Begriff, keine Bestands-Nennung); §7 der slice-242-Planung (Closure-Notiz) steht noch aus.
- **Klasse:** unbelegter-Urteilsbegriff-bei-adoption

---

## Negativbefunde (geprüft, ohne Befund)

1. **Senkungs-Satz (Prüfpunkt 1):** byte-identisch — gegen den tatsächlichen
   Änderungs-Commit `088c3699`: „Jede Schwellen-Senkung (Coverage, Linter-Strenge,
   Prüfregel) ist ein ADR, kein PR-Kommentar." bleibt als Satz wörtlich erhalten; die
   Ergänzung ist additiv (Zeilen-/Absatz-Ende). Die im Auftrag genannte Gegenprobe gegen
   `d19b2095` ist leer — der Commit berührt AGENTS.md nicht — und prüft an dieser Stelle
   nichts; die substantive Gegenprobe führt F-1 (F-1, nicht DoD-Verstoß).
2. **Delta-Treue (Prüfpunkt 2):** `v6.13.0` · `regelwerk/modul-04-adrs.md` §Kernidee
   (Zeilen 40–48): alle fünf Tragflächen übernommen — Aufnahme eines existierenden,
   unabhängig lauffähigen Wächters = kein ADR-Anlass; Verweis auf die tragende ADR
   genügt; existiert keine, trägt die *Einführung* des Wächters; neue Fehlerklasse /
   neuer Scope / Widerspruch = weiterhin eigene ADR; „Im Zweifel: ADR". Keine
   inhaltliche Drift; die Baseline-Formel „nicht seine spätere Aufnahme" ist als
   Redundanz zur ersten Aussage korrekt gestrichen, „PR-blockierenden Satz" als bewusste
   Abgrenzung (Nr. 3).
3. **Terminologie-Abstand (Prüfpunkt 3):** `grep "PR-blockier" AGENTS.md` → 0 Treffer;
   die Übertragung auf den eigenen Betrieb (`make gates`, lokal + CI) ist im Vollzug
   formuliert, nicht zitiert — die Plan-Zusage hält.
4. **Re-Anker (Prüfpunkt 4):** `AGENTS.md:286` trägt „Halluzinierte Gates sind die
   häufigste Form von Harness-Lüge" wörtlich — komplett in Zeile 286. Die
   +5-Zeilen-Verschiebung aus der §3.6-Ergänzung (6 INSERT / 1 DELETE) rechnet
   281 → 286 exakt nach; der Re-Anker ist inhaltlich richtig, sitzt nur im falschen
   Commit (F-1).
5. **Keine Gate-Senkung:** die Klarstellung senkt keine Schwelle und unterdrückt kein
   Gate — der Senkungs-Satz und „Im Zweifel: ADR" bleiben; die Ergänzung folgt dem
   Kanon-Delta (Kanon vor Briefing), nicht einer lokalen Lockerung. Anker 4 (Suppression
   ohne ADR) greift nicht.
6. **Sensors:** `make doc-check` eigenständig nachgefahren (Docker, `--network none`):
   „889 Datei(en) geprüft, 0 Befund(e)" — deckt die Gates-Aussage des Auftrags für den
   Doku-/cite-Teil des Umfangs.
7. **Sechzehn-Fragen-Sweep:** reiner Doku-Diff — keine Gate-Skripte, Module, Imports,
   Netz-Zugriffe, Code-Kommentare oder Zustandsfelder berührt (Fragen 1–7, 10–18 ohne
   Gegenstand); es entsteht keine Messmethode und kein neuer öffentlicher Vertrag.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 (Commit-Grenze / Botschaft-vs-Diff) |
| LOW | 0 | — |
| INFO | 1 | F-2 (unbelegter Urteilsbegriff — Closure-Sichtung) |

## Verdikt

MEDIUM blockiert typischerweise — hier als **Push-Blocker in der Historie, nicht im
Ist-Zustand** begründet: der Arbeitsbaum an `d19b2095` ist inhaltlich vollständig,
delta-treu, terminologie-sauber und gate-grün (`make doc-check` 889/0). Beanstandet wird
die Commit-Grenze samt Botschafts-Aussage (F-1) — die beiden Commits sind noch lokal
(`origin/main` = `e4dc7f85`) und damit korrigierbar, bevor die falsche Zuordnung Bestand
wird. F-2 ist ein Sichtungspunkt für die Closure-Notiz, kein Änderungsbedarf am Text.
