# Review — slice-089 (welle-67 Etappe D „Doc-Form", AGENTS.md-Angleich an v5.0.0)

- **Datum:** 2026-08-02
- **Reviewer:** unabhängiger, kontext-getrennter Frischkontext (Reviewer-Skill v1.2.0)
- **Gegenstand:** slice-089 — AGENTS.md an die v5.0.0-Baseline-Form angeglichen
  (D-11 §1-Currency, D-8 ADR-Re-Evaluierungs-Trigger, D-4 Slice-Vorprüfungen)
- **Commit-Range:** `6560a7d..HEAD` (`968c147` open, `2c42d74` AGENTS)
- **Ziel-Form:** `.harness/baseline/v5.0.0/templates/AGENTS.template.md`
- **Maßgebliche Regel:** `.harness/baseline/v5.0.0/regelwerk/modul-09-implementierung.md`
  §AGENTS.md-Regeln
- **Slice-Plan:** `docs/plan/planning/in-progress/slice-089-doc-form-agents.md`
- **Berührte Dateien im Range:** `AGENTS.md`, `docs/plan/planning/in-progress/roadmap.md`,
  `docs/plan/planning/in-progress/slice-089-doc-form-agents.md` (kein spec/, kein adr/, kein Code)

---

## Findings

### F-1 (LOW, Maintainability / Doku-Drift) — `docs/plan/planning/in-progress/roadmap.md:130` + `docs/plan/planning/in-progress/slice-089-doc-form-agents.md:44`

Die D-4-Regel („jeder Slice-Plan trägt vor der Sub-Area-Modus-Begründung die zwei
Vorprüfungen") ist in AGENTS.md unter **§5 Dokumentations-Regeln** verankert
(`AGENTS.md:194`). Der Slice-Plan §2 Schritt 3 und DoD-Punkt 3 sowie der committete
Roadmap-Closure-Log-Eintrag benennen dagegen **§6 Minimal Agent Workflow** als Ort
(„Slice-Vorprüfungen in §6 [D-4]").

- **Failure-Szenario:** Der Verifier prüft den DoD-Punkt „AGENTS.md §6 trägt die
  Slice-Vorprüfungen (D-4)" gegen §6 (Minimal Agent Workflow, `AGENTS.md:200`), findet
  die Regel dort nicht (sie steht in §5) und liest entweder einen falschen DoD-Fail
  oder muss die Fundstelle nachschlagen. Die committete Roadmap-Zeile weist einen
  späteren Leser navigierend auf die falsche Sektion.
- **Keine Defekt-Aussage über AGENTS.md selbst:** §5 (das u. a. die Slice-Lifecycle-
  und Roadmap-Regeln trägt) ist ein sachlich stimmiger Ort für eine
  Slice-Plan-Struktur-Regel — der Drift liegt zwischen den Planungs-Layer-Texten
  (Plan/DoD/Roadmap: §6) und dem Artefakt (§5), nicht in der Regel selbst.
- **verifizierbar:** ja — `grep -nE '^## [0-9]\.' AGENTS.md` zeigt §5 (Z. 168) vor der
  Vorprüfungs-Zeile (Z. 194) vor §6 (Z. 200); der Roadmap-/DoD-Wortlaut steht im Diff.

---

## Negativbefunde (geprüft, ohne Befund)

1. **Currency-Vollständigkeit (D-11) — ohne Befund.** `grep` über das ganze `AGENTS.md`
   nach „verkörpert keine", „co-located", „lab-templates.zip", „.harness/cache",
   „MR-018", „Autorenquelle", „Cache-Staging", „Adoptions-/Drift-Audit-Staging" liefert
   **keinen** Treffer — die stale Templates-Cache-Prosa ist restlos entfernt. Beide
   Bäume sind committet vendored (`git ls-files .harness/baseline/v5.0.0/` = 51 Dateien;
   `regelwerk/` und `templates/` beide präsent, dazu `SHA256SUMS`). Der neue Kanon-Anker
   `#agentsmd-regeln-modul-9` löst auf (Heading „### AGENTS.md-Regeln (Modul 9)" in
   `modul-09-implementierung.md`; doc-check 0 Befunde). Die Aussage „aus demselben
   self-contained Bundle" ist faktentreu: `tools/harness/fetch-baseline-cache.sh` (Kopf)
   und `harness/conventions/MR-023-baseline-v500.md` belegen ein einziges Asset
   `lab-regelwerk.zip` mit beiden Bäumen; `lab-templates.zip` existiert seit v5.0.0
   nicht mehr.
2. **Routet-nicht-Spiegelt — ohne Befund.** Die neue §1-Templates-Prosa gibt die
   **Template-Form** wieder (`AGENTS.template.md` §1 trägt dieselbe Zwei-Rollen-Aussage
   „Referenz-Form + Kopiervorlage"), nicht kanonischen **Inhalt**; Stand/Provenance wird
   ausdrücklich an `harness/conventions.md` (§Adoptierte Konventions-Quellen bzw.
   §Baseline) geroutet. Kein Widerspruch zu `modul-09` §AGENTS.md-Regeln (Hard Rules +
   Pointer) oder zur Source-Precedence-Ordnung (§2 unverändert).
3. **Faktentreue — ohne Befund.** `harness/conventions/done/MR-014-slice-adr-haus-stil.md`
   (unter `done/`) trägt „Aufgelöst durch: Baseline-Konformität (Form ist Baseline-Wahl)"
   — die §1-Aussage „Haus-Stil-Form in Etappe C als baseline-konforme Form-Wahl aufgelöst,
   nicht als Fork" ist gedeckt. Die „20 grandfathered"-Aussage stimmt: 47 Accepted-ADRs,
   davon 27 mit exaktem `## Re-Evaluierungs-Trigger`-Heading und 20 ohne. Die Behauptung
   „nachträgliches Ergänzen bräche `make adr-check`" ist korrekt: adr-check (Modul `vcs`)
   exemptiert ausschließlich `## Geschichte`-Anhänge und den `**Status:**`-Übergang; ein
   neues `## `-Heading im ADR-Core einer Accepted-ADR fällt nicht darunter.
4. **D-8/D-4-Korrektheit — ohne Befund.** Kein Dublett, keine Kollision: D-8 (§5) ist
   widerspruchsfrei zu §3.5 (ADR-Immutabilität) und zu „Neue ADRs müssen den ADR-Index
   aktualisieren". Das von D-4 referenzierte Register ist vorhanden
   (`docs/plan/planning/observations.md`, Tabelle `— keine —`). Die Modul-Zeiger stimmen:
   `modul-05-planning-harness.md` §Ziel-Form trägt beide Vorprüfungen („Sub-Area-Wahl
   prüfen" / „Offene Beobachtungen sichten"), `modul-06-roadmap.md` trägt den
   Trigger-Audit und das Beobachtungs-Register.
5. **Keine Kollateral-Schäden — ohne Befund.** Der `mr-018`-Anker existiert weiter in
   `harness/conventions.md` (Zeile als „v5.0.0 (Bundle vendored beide Bäume)" aufgelöst
   markiert); die Verweise aus `MR-023`/`MR-019`/`MR-020` und der Alt-Review-Datei zeigen
   auf `conventions.md`, nicht auf die aus AGENTS.md entfernte Fundstelle — kein
   `target-missing`. Keine Accepted-ADR, keine Spec, kein Code im Range berührt.
6. **Gates — ohne Befund (bestätigt).** Netzloser Dogfood `d-check:latest --network none`:
   316 Dateien, 0 Befunde. `make adr-check RANGE=6560a7d..HEAD`: 0 Befunde. `make gates`:
   grün (doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency
   + planning-check).

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

---

## Verdikt: **abnahmereif**

Kein HIGH/MEDIUM. Der einzige Befund (F-1) ist eine LOW-Doku-Drift zwischen den
Planungs-Layer-Texten (Plan/DoD/Roadmap sagen §6) und der tatsächlichen Verankerung
(§5) — blockiert nach Reviewer-Skill (nur HIGH/MEDIUM blockieren typischerweise) den
Merge nicht. Empfehlung an die Implementation/Verifikation: entweder den §6-Wortlaut in
Plan-DoD und Roadmap-Zeile auf §5 nachziehen oder die Regel nach §6 verschieben — Wahl
der Implementation. Die drei Kern-Angleichungen (D-11 Currency, D-8 ADR-Trigger, D-4
Vorprüfungen) sind inhaltlich korrekt, faktentreu, form-konform zur Ziel-Vorlage und
gate-grün.
