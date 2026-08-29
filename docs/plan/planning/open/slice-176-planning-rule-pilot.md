# Slice slice-176: Die Planungs-Regeln werden zugestellt, wenn sie gelten

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §1 (das Prinzip: Hard Rules und
**Pointer**, keine Duplikation) und §3 (der Abschnitt, der über die Hälfte der
Datei trägt);
[`MR-054`](../../../../harness/conventions.md#mr-054) (die Beleg-Form, die
dieser Slice übernimmt);
[`MR-043`](../../../../harness/conventions.md#mr-043) (der Werkzeug-Einstieg
importiert `AGENTS.md` — der Mechanismus, den dieser Slice ergänzt);
[`MR-042`](../../../../harness/conventions.md#mr-042) (die Präzedenz für einen
**werkzeug-lokalen** Träger und seine Nicht-Zusagen).

**Berührte Spec-Stellen:** — (Werkzeug-Konfiguration und Konventionsform; keine
Produkt-Anforderung berührt).

**Verantwortlich:** — (bis zur Priorisierung).

**Autor:** pt9912. **Datum:** 2026-08-29.

---

## 1. Ziel

**`AGENTS.md` ist auf 511 Zeilen gewachsen, und die Regeln, die eine Sitzung
gekostet haben, standen nicht darin.**

Gemessen: die Datei trägt 511 Zeilen gegen 236 der Baseline-Vorlage; allein
**§3 Harte Regeln** sind 280 davon — mehr als die Hälfte, und die meisten
dieser Regeln gelten für **einen** Pfad. Die Werkzeug-Dokumentation nennt genau
diese Lage und genau dieses Mittel:

> target under 200 lines per CLAUDE.md file. Longer files consume more context
> and reduce adherence. If your instructions are growing large, use path-scoped
> rules so instructions load only when Claude works with matching files.

**Der Import löst es nicht.** [`CLAUDE.md`](../../../../CLAUDE.md) zieht
`AGENTS.md` über `@AGENTS.md` in jeden Lauf — dieselbe Dokumentation sagt dazu,
Imports helfen der *Organisation*, nicht dem Kontext: sie laden beim Start
vollständig.

**Was gefehlt hat, ist belegt.** In der auslösenden Sitzung liefen drei Slices
(168, 169, 170) durch, ohne dass der Zyklus (`Spec → ADR → Plan → Code →
Review → Verifikation → Closure`) im Kontext war. Er steht in `modul-01`, das
niemand geöffnet hatte — und in `AGENTS.md` steht er nicht, korrekterweise:
§6 beschreibt den **Implementer**-Workflow, und der endet bei Schritt 8. Folge:
Review und Verifikation fielen aus, das Beobachtungs-Register blieb
unangetastet, drei Slices gingen mit offenem DoD-Haken nach `done/`.

**Eine pfad-gebundene Regel hätte den Zyklus geladen, als die erste Slice-Datei
angefasst wurde** — nicht, weil jemand vorher wusste, dass er ihn braucht.

## 2. Vorgehen

1. **Eine Regel als Pilot, nicht fünf.** `.claude/rules/planning.md` mit
   `paths: ["docs/plan/planning/**"]` — der Pfad, dessen Fehlen gemessen ist.
   Die anderen vier Kandidaten (`adr`, `spec`, `workflows`, `hexagon`) folgen
   **erst**, wenn das Muster trägt.
2. **Die Regel stellt zu, sie dupliziert nicht.** Sie trägt den Zyklus, die
   Lifecycle-Kanten, die Closure-Pflichten und die Rollen-Trennung in
   verdichteter Form — jeweils mit **zwei** Bindungen an die Quelle: einer
   `d-check:cite`-Direktive für den Wortlaut ([`MR-054`](../../../../harness/conventions.md#mr-054))
   **und** einem Markdown-Link auf den Regelwerk-Abschnitt für die Vollform.
   Damit ist sie keine dritte Wahrheitsquelle, sondern die kontextsensitive
   Auslieferung der zweiten.
3. **Die Verweise sind gate-geprüft, und das ist der Grund für sie.**
   `.claude/rules/` liegt im Scan-Bereich (`scan.roots: ["."]`), also halten
   `links`, `anchors` und `citations` die Regel gegen ihre Quelle. Eine
   Zustellung, die still von ihrer Quelle abdriftet, wäre schlimmer als keine.
4. **`AGENTS.md` gibt erst danach ab — und nur, was die Regel wirklich trägt.**
   Die Verdichtung ist ein **eigener Schritt nach der Messung**, kein
   Vorgriff: erst zeigen, dass die Zustellung greift, dann kürzen.
5. `make gates`; **Review** und **Verifikation** als getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Kürzen von `AGENTS.md`.** Solange nicht gemessen ist, dass die Regel
  greift, wäre jede gestrichene Zeile ein Verlust ohne Ersatz. Die Datei bleibt
  in diesem Slice unverändert.
- **Keine weiteren vier Regeln.** Wer fünf Dateien anlegt, bevor eine getragen
  hat, verteilt einen ungeprüften Entwurf über das halbe Repo.
- **Kein Ersatz für Gates.** Die Werkzeug-Dokumentation sagt es selbst: *„Claude
  treats them as context, **not enforced configuration**."* welle-86 bleibt
  davon unberührt.
- **Keine Aussage über andere Werkzeuge.** Die Regel wirkt in **einem**
  Werkzeug; `AGENTS.md` bleibt die werkzeug-neutrale Datei.

## 4. Definition of Done

- [ ] `.claude/rules/planning.md` existiert, trägt `paths`-Frontmatter auf
      `docs/plan/planning/**` und bleibt **unter 60 Zeilen** — eine Zustellung,
      die selbst zu lang ist, verfehlt ihren Zweck.
- [ ] Jeder normative Satz der Regel trägt **beide** Bindungen: eine
      `d-check:cite`-Direktive auf den Wortlaut **und** einen Link auf den
      Regelwerk-Abschnitt; `make doc-check` prüft beides grün.
- [ ] **Gemessen, dass die Zustellung greift:** ein Lauf, der eine Datei unter
      `docs/plan/planning/` liest, führt die Regel im Kontext — belegt über den
      `InstructionsLoaded`-Hook oder `/context`, mit der Ausgabe im Slice.
- [ ] Die **Nicht-Zusagen** stehen geschrieben, nicht implizit: werkzeug-lokal,
      kein Gate, kein Ersatz für `AGENTS.md`.
- [ ] Der Nachfolge-Entscheid ist benannt — ob und wann die vier weiteren
      Regeln entstehen und was `AGENTS.md` dann abgeben kann.
- [ ] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Ein werkzeug-lokaler Träger ist ungebunden, sobald das Werkzeug wechselt.**
  [`MR-042`](../../../../harness/conventions.md#mr-042) hat das für den
  Tool-Call-Wächter schon ausbuchstabiert: *keine CI ruft ihn, ein Lauf ohne
  dieses Werkzeug ist ungebunden.* Heute arbeitet dieses Repo nur mit einem
  Werkzeug — das ist eine **Momentaufnahme**, keine Eigenschaft. —
  **Ausgang:** *(bei Closure)*
- **Zwei Orte für dieselbe Aussage sind die Drift, die dieses Repo als
  [`BEO-010`](../observations.md) führt.** Die Regel und `AGENTS.md` sagen
  beide etwas über Planung. Solange `AGENTS.md` nichts abgibt (§3), ist das
  **Duplikation** — und der Slice endet in genau diesem Zustand. Die Auflösung
  liegt im Nachfolger, und bis dahin ist der Zustand benannt statt übersehen. —
  **Ausgang:** *(bei Closure)*
- **Die Bindung an das Regelwerk kostet beim Bump doppelt:** der Link trägt den
  `<tag>` ([`MR-021`](../../../../harness/conventions.md#mr-021)), die
  `cite`-Spanne die Zeilennummer ([`MR-051`](../../../../harness/conventions.md#mr-051)).
  Jede neue Regel-Datei erhöht den Bump-Aufwand — und dieser Slice legt die
  erste an. — **Ausgang:** *(bei Closure)*
- **Die Wirkung ist schwer zu messen.** Dass eine Regel *geladen* wurde, ist
  belegbar; dass sie *gewirkt* hat, nicht. Der Anlassfall lässt sich nicht
  wiederholen — er bestand darin, dass niemand etwas vermisste. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — heute hält
[slice-171](../in-progress/slice-171-vorpruefungen-belegen.md) den Slot.

**Rückführungen:** `in-progress` → `open`, falls die Messung zeigt, dass die
Regel **nicht** beim Lesen einer passenden Datei geladen wird — dann trägt die
Annahme nicht, auf der der ganze Schnitt steht, und der Befund ist ein anderer.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `.claude/` (Werkzeug-Konfiguration) und
  `docs/plan/planning/` (der adressierte Pfad). Beide fallen unter den Default
  `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration); eine eigene Deklaration führt nur `tools/harness/`.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-29, höchste Kennung
  `BEO-021`): [`BEO-010`](../observations.md) — eine Liste mit Spiegeln
  außerhalb ihrer Datei: die Regel und `AGENTS.md` §3 wären genau das, solange
  keine abgibt, und das steht als Risiko in §5;
  [`BEO-012`](../observations.md) — eine Quelle über ihren Geltungsbereich
  hinaus zitiert: eine verdichtete Zustellung ist genau diese Gefahr, weshalb
  jeder Satz **zwei** Bindungen an die Quelle trägt;
  [`BEO-011`](../observations.md) — Regel aus dem Anlass: der Pilot-Pfad ist
  aus **einem** Vorfall gewählt, und genau deshalb bleibt es bei einer Regel
  statt fünf. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:224-225 -->
  > **Keine Treffer sind ebenfalls eine
  > Antwort** und werden notiert.

- **Nachtlauf-Stand lesen:** entfällt in `open/` — der Block entsteht
  **spätestens bei der Beanspruchung** (`open→in-progress`), weil ein zum
  Planungszeitpunkt gelesener Stand bis dahin veraltet wäre
  ([`MR-053`](../../../../harness/conventions.md#mr-053)).

Slice-ID: slice-176. Betroffene IDs:
[`MR-054`](../../../../harness/conventions.md#mr-054),
[`MR-042`](../../../../harness/conventions.md#mr-042),
[`MR-043`](../../../../harness/conventions.md#mr-043). Module: `links`,
`anchors`, `citations`. Gates: `make gates`, `make doc-check`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Eine Werkzeug-Konfigurationsdatei plus ihre
Bindung an vorhandene Quellen; kein Produkt-Code, kein Fremdsystem, keine
Reconciliation.

## 9. Closure-Notiz (nach `done/`)
