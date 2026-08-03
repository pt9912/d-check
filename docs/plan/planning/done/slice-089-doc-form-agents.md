# Slice slice-089: Regelwerk-Migration Etappe D — Doc-Form (AGENTS.md-Currency · ADR-Trigger · Slice-Vorprüfungen)

**Status:** Done (welle-67, Doc-Form abgeschlossen 2026-08-02).

**Welle:** welle-67-baseline-v500-migration (Etappe D, zweiter „Mini-Welle"-Slice, nach
[slice-088](../done/slice-088-etappe-d1-doc-form.md)).

**Bezug:** Etappe D (Form-Konformität) aus
[slice-085](../done/slice-085-etappe-b-modul-delta.md) §3.2 — **D-11** (AGENTS.md ↔
`AGENTS.template.md`), **D-8** (ADR-Re-Evaluierungs-Trigger) und die aus
[slice-088](../done/slice-088-etappe-d1-doc-form.md) delegierte **D-4**-AGENTS-Verankerung.
Mini-Welle: slice-088 Planning-Layer (done) → **slice-089 Doc-Form** → slice-090
Review-Infrastruktur (D-6/D-7/D-10) → slice-091 Slice-`Status:`-Feld (D-5). **Kein Release.**

**Autor:** pt9912. **Datum:** 2026-08-02.

---

## 1. Ziel

`AGENTS.md` an die v5.0.0-Baseline-Form angleichen: **§1-Currency** (D-11 — die von
Etappe A/C überholte Templates-Cache-Prosa raus, `modul-09`-Kanon-Zeiger +
`{regelwerk,templates}`-vendored-Layout rein), die **ADR-Re-Evaluierungs-Trigger**-
Konvention verankern (D-8) und die zwei **Slice-Vorprüfungen** als Konvention aufnehmen
(D-4). AGENTS.md ist der **Source-Precedence-Anker** — vorsichtig, review-pflichtig.

## 2. Vorgehen

1. **D-11 AGENTS §1-Currency.** Die stale Templates-Cache-Prosa (Producer-/Self-Hoster
   „verkörpert keine Templates"; Entpacken nach dem entfallenen Cache-Zweig; das nicht
   mehr existierende Einzel-`lab-templates.zip`) durch die v5.0.0-Realität ersetzen:
   **beide** Bäume `.harness/baseline/<tag>/{regelwerk,templates}/` committet vendored
   (aus dem self-contained Bundle); die Templates als **Referenz-Form** (Ziel der
   `../templates/…`-Regelwerk-Verweise) **und** Kopiervorlage. Den **Kanon-Zeiger** auf
   `modul-09-implementierung.md` §AGENTS.md-Regeln ergänzen. Konsistent mit der
   AGENTS-routet-Regel (Pointer, nicht Spiegel) und der gelebten Haus-Stil-Form (deren
   Adaption in Etappe C als baseline-konform aufgelöst wurde).
2. **D-8 ADR-Trigger-Konvention.** In §5 (Dokumentations-Regeln) die Regel ergänzen:
   **neue ADRs tragen die Sektion `## Re-Evaluierungs-Trigger`** (oder „permanent"). Die
   **20** Accepted-ADRs ohne Sektion sind **immutable/grandfathered** (nicht
   retrofit-bar — das Trigger-Feld liegt im ADR-Core, `adr-check` schlägt bei nachträglicher
   Ergänzung an); der Welle-Closure-Trigger-Audit (slice-088-Wellendokument §3 / D-9)
   prüft die Trigger.
3. **D-4 Slice-Vorprüfungen.** In §5 (Dokumentations-Regeln) als Slice-Plan-Form-Regel
   verankern: jeder Slice-Plan trägt **vor** der Sub-Area-Modus-Begründung „Sub-Area
   prüfen · offene Beobachtungen sichten" (`observations.md`) — Zeiger auf die
   Baseline-Slice-Form, nicht Duplikat der Regel. (§5 trägt bereits die
   Slice-Lifecycle-/Roadmap-Regeln; §6 bleibt der Implementierungs-Workflow.)
4. **Gate.** `make gates` + `make adr-check` grün (**keine** `Accepted`-ADR berührt);
   unabhängiger Frischkontext-Review.

## 3. Abnahme-Punkte / Entscheidungen (mit Default)

1. **AGENTS §1 Templates-Framing (D-11).** Die stale Prosa („d-check verkörpert keine
   Templates; Cache-Staging") ist von v5.0.0 überholt (Templates sind jetzt vendored).
   → **Default:** an die Template-Sprache angleichen (Templates vendored als Referenz-
   **und** Kopiervorlage) + den gelebten Haus-Stil als **Form-Wahl** benennen; **keine**
   Re-Aktivierung des Cache-Modells, kein Widerspruch zur AGENTS-routet-Regel.

## 4. Definition of Done

- [x] `AGENTS.md` §1 v5.0.0-konform (`modul-09`-Kanon-Zeiger, `{regelwerk,templates}`,
  **keine** stale Cache-/`lab-templates.zip`-/„keine-Templates"-Prosa).
- [x] `AGENTS.md` §5 trägt die ADR-Re-Evaluierungs-Trigger-Konvention (D-8);
  immutable-Grandfathering benannt.
- [x] `AGENTS.md` §5 trägt die Slice-Vorprüfungen-Form-Regel (D-4).
- [x] `make gates` + `make adr-check` grün; unabhängiger Frischkontext-Review.

## 5. Risiken / offene Punkte

- **AGENTS = Source-Precedence-Anker** — der §1-Umbau ist vorsichtig und
  review-pflichtig; die „routet, nicht spiegelt"-Form (Pointer statt Duplikat) einhalten.
- **Stale-Links vs. stale-Aussage:** die Ziel-Anker der aufgelösten MRs existieren im
  Index — die **Aussage** ist überholt, also **umschreiben**, nicht nur re-linken.

## 6. Trigger

Abschluss von [slice-088](../done/slice-088-etappe-d1-doc-form.md) (Planning-Layer).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** berührt *Harness/Prozess*-Doku (`AGENTS.md`) — greenfield, GF.
- **Offene Beobachtungen sichten:** `observations.md` = `— keine —`; nichts zu
  berücksichtigen.

## 8. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc/Prozess führt. Berührt `AGENTS.md` (Source-Precedence-Anker);
greenfield-Currency-Angleich an die adoptierte Baseline, ohne Brownfield-Spec.

## 9. Closure-Notiz (nach `done/`)

Umgesetzt: `AGENTS.md` (der Source-Precedence-Anker) ist an die v5.0.0-Baseline-Form
angeglichen.

- **D-11 §1-Currency:** die von Etappe A/C überholte Templates-Cache-Prosa entfernt
  (Producer-/Self-Hoster „verkörpert keine Templates", `lab-templates.zip`, der
  entfallene Cache-Zweig, der aufgelöste Verweis auf die frühere Templates-Adaption) und
  durch die v5.0.0-Realität
  ersetzt: **beide** Bäume `{regelwerk,templates}` committet vendored, Templates als
  **Referenz-/Kopiervorlage**, die gelebte **Haus-Stil-Form** (in Etappe C als
  baseline-konform aufgelöst) benannt; **Kanon-Zeiger** auf `modul-09` ergänzt.
- **D-8 §5:** neue ADRs tragen `## Re-Evaluierungs-Trigger` (oder „permanent"); die **20**
  `Accepted`-ADRs ohne Sektion sind immutable/**grandfathered** (das Feld liegt im
  ADR-Core → `make adr-check` bräche bei Nachtrag); der Welle-Closure-Trigger-Audit
  bestätigt/revidiert sie.
- **D-4 §5:** jeder Slice-Plan trägt vor der Modus-Begründung die zwei Vorprüfungen
  (Sub-Area · offene Beobachtungen im Register sichten) — die aus slice-088 delegierte
  Verankerung; in **§5** (nicht §6 — der sachlich richtige Ort neben den Slice-Lifecycle-/
  Roadmap-Regeln; Review-LOW F-1).

**Review:** unabhängiger Frischkontext-Review
(`docs/reviews/2026-08-02-slice-089-agents-doc-form-review.md`): **abnahmereif**, HIGH 0 /
MEDIUM 0 / LOW 1 / INFO 0 (F-1: §5-vs-§6-Drift in den Planungs-Texten → auf §5 gezogen).
Verifiziert: 0 stale Begriffe in `AGENTS.md`, beide Bäume vendored, `modul-09`-Anker löst,
routet-nicht-spiegelt gewahrt, „20 grandfathered" exakt, die entfernte
Templates-Adaptions-Fundstelle ohne Kollateral. `make gates` + `make adr-check` grün, **keine** ADR/Spec/Code berührt.

**Anschluss:** **slice-090** (Review-Infrastruktur D-6/D-7/D-10), dann slice-091
(Slice-`Status:`-Feld D-5). Danach schließt welle-67.
