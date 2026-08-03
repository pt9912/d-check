# Slice slice-088: Regelwerk-Migration Etappe D — Planning-Layer-Form (Wellen-Lifecycle · Roadmap · Beobachtungs-Register)

**Status:** Done (welle-67, Planning-Layer abgeschlossen 2026-08-02).

**Welle:** welle-67-baseline-v500-migration (Etappe D, erster von vier
„Mini-Welle"-Slices, nach
[slice-087](../done/slice-087-spec-historie-referenzrichtung.md)).

**Bezug:** Erster Slice der **Etappe D (Form-Konformität)** aus
[slice-083](../done/slice-083-regelwerk-v500-migration-analyse.md) §2.7 /
[slice-085](../done/slice-085-etappe-b-modul-delta.md) §3.2 (Findings D-1…D-11).
Etappe D ist als **Mini-Welle** in vier Slices geschnitten (Nutzer-Entscheid
2026-08-02) und **neu geordnet**, nachdem sich zeigte, dass die **laufende**
welle-67 ohne ihr Baseline-Wellendokument läuft (Nutzer-Hinweis): **slice-088
Planning-Layer-Form** (D-1/D-2/D-3/D-4/D-9), slice-089 Doc-Form (D-8/D-11),
slice-090 Review-Infrastruktur (D-6/D-7/D-10), slice-091 Slice-`Status:`-Feld (D-5).
Dieser Slice zieht die **kohärente Planning-Schicht** nach — Roadmap **+** Welle **+**
Register gehören zusammen. Maßgeblich: das Baseline-Regelwerk `modul-06-roadmap.md`
(Roadmap-Struktur, Wellen-Closure-Prozedur, Beobachtungs-Register) und
`modul-05-planning-harness.md` (Lifecycle, Vorprüfungen). **Kein Release.**

**Autor:** pt9912. **Datum:** 2026-08-02.

---

## 1. Ziel

Die d-check-eigene **Planning-Schicht** an die v5.0.0-Baseline-Form heben: die
laufende **welle-67** bekommt ihr **Wellendokument** (D-2), das
**Beobachtungs-Register** entsteht als stehende Datei (D-3), die **Roadmap** führt die
Baseline-Abschnitte (D-1), Slices tragen die zwei **Vorprüfungen** vor der
Modus-Begründung (D-4), und der **Carveout-/Trigger-Audit** ist als
Welle-Closure-Schritt verankert (D-9). Keine neuen Produkt-Features, kein Release.

## 2. Vorgehen

1. **D-2 Wellendokument (welle-67).** `welle-67-baseline-v500-migration.md` **flach**
   unter `docs/plan/planning/` anlegen (Ziel-Form `welle.template.md`): Welle-Ziel ·
   Start-Trigger · Closure-Trigger · Slices 084–091 (**ohne** Status-Spalte) ·
   Abhängigkeiten · Out-of-Scope · §7-Closure-Notiz **ausstehend** (Pointer erst bei
   Closure, Ruheort-Regel). Lifecycle: flach solange aktiv, `git mv` → `done/` neben
   `welle-67-results.md` bei Welle-Closure (Zustand = Verzeichnis, **kein** Status-Feld).
2. **D-3 Beobachtungs-Register.** `observations.md` als **stehende** Datei flach anlegen
   (Ziel-Form `observations.template.md`): Tabelle
   `| Kennung | Beobachtung | Sub-Area | Zähler | Belege | Stand |` + Sektion
   „Gestrichene Einträge". Start `— keine —` — die verkörperten Lehren stehen bereits in
   `AGENTS.md`/Gates/Konventionen; das Register trägt **offene** (unter-Schwelle)
   Beobachtungen, davon gibt es keine. Die leere Liste **ist** die Aussage.
3. **D-1 Roadmap-Abschnitte.** `roadmap.md` um die fehlenden Baseline-Abschnitte
   ergänzen: `## Meilensteine`, `## Abhängigkeitsgraph` (mermaid), `## Abgeschlossene
   Wellen` (Tabelle mit Zeiger auf `done/welle-NN-results.md`); den Closure-Bestand aus
   der `## Aktuelle Welle`-„Vorgänger"-Prosa in `## Abgeschlossene Wellen` überführen.
4. **D-4 Slice-Vorprüfungen.** Die zwei Vorprüfungen (`modul-05` §Zwei Schritte vor der
   Modus-Begründung / `modul-06` Sichtungs-Schritt) in die Slice-Form aufnehmen:
   „**Vorgelagert** — Sub-Area prüfen · offene Beobachtungen sichten" **vor** der
   §Sub-Area-Modus-Begründung; für slice-088 selbst nachgezogen (§7). Die
   Konventions-Verankerung für **alle** künftigen Slices (`AGENTS.md`
   §Slice-Anlege-Prozess) zieht slice-089 (D-11, AGENTS-Angleich) mit — hier nur das
   gelebte Modell.
5. **D-9 Carveout-/Trigger-Audit.** Den Welle-Closure-Trigger-Audit
   (Carveout · bootstrap-aware Gate · ADR-Re-Eval, `modul-06` Closure-Schritt 2) im
   welle-67-Closure-Trigger verankern; 0 aktive Carveouts → **latent**, aber benannt.
6. **Gate.** `make gates` (inkl. planning-check) + `make adr-check` grün; unabhängiger
   Frischkontext-Review.

## 3. Abnahme-Punkte / Entscheidungen (mit Default)

1. **`## Aktuelle Welle`-Form (D-1) — Ruhe-Marker ↔ Template.** Das Template führt
   Struktur-Felder (Welle-ID/Start/Geplantes Ende/Closure-Trigger) und nennt immer eine
   laufende Welle; d-checks `planning`-Modul
   ([`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in))
   **erzwingt** den „Keine aktive Welle"-Ruhe-Marker in genau diesem Abschnitt.
   → **Umgesetzt 2026-08-02:** Ruhe-Marker + Prosa-Form als **deklarierte Adaption**
   behalten — ein neuer Eintrag im Konventionsspeicher hält fest: die drei
   Pflicht-Bestandteile stehen in Prosa (Feld-Details im Wellendokument), der wellenlose
   Zustand explizit über den Ruhe-Marker („Keine aktive Welle."), gate-erzwungen vom
   `planning`-Modul. Die Mechanisierung ist d-check-Identität — kein Umbau des Moduls.
2. **Register-Start (D-3):** leer (`— keine —`) vs. retroaktives Backfill.
   → **Default:** leer starten (Baseline: „damit fängt jedes Repo an"; die
   wiederkehrenden Lehren sind bereits verkörpert).

## 4. Definition of Done

- [x] `welle-67-baseline-v500-migration.md` flach angelegt, template-konform, §7 ausstehend.
- [x] `observations.md` als stehende Datei angelegt (`— keine —`, template-konform).
- [x] `roadmap.md` führt die Baseline-Abschnitte; Closure-Bestand in `## Abgeschlossene Wellen`.
- [x] Slice-Form trägt die zwei Vorprüfungen (D-4); welle-67-Closure-Trigger nennt den Audit (D-9).
- [x] `## Aktuelle Welle`-Form-Entscheid (Abnahme-Punkt 1) umgesetzt (Adaptions-Eintrag).
- [x] `make gates` + `make adr-check` grün; unabhängiger Frischkontext-Review.

## 5. Risiken / offene Punkte

- **planning-check-Kopplung:** Roadmap-Umbau + neue flache `welle-*.md`/`observations.md`
  dürfen den Ruhe-Marker-/Heading-Guard nicht brechen (das `planning`-Modul prüft den
  `## Aktuelle Welle`-Heading exakt und die in-progress-↔-Ruhe-Invariante).
- **Ruheort-Regel:** die §7-Pointer des Wellendokuments lösen von `done/` auf — daher
  **erst bei Closure** schreiben, sonst `target-missing` von flach.
- **Retroaktivität:** welle-67 läuft schon; das Wellendokument entsteht mitten drin
  (statt bei Eröffnung) — bewusst, um die laufende Welle konform zu dokumentieren.

## 6. Trigger

Abschluss von [slice-087](../done/slice-087-spec-historie-referenzrichtung.md)
(C-3-Nachzug) + Nutzer-Hinweis, dass welle-67 ohne Wellendokument läuft.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** berührt *Harness/Prozess* (Planning-Layout) — greenfield, GF.
- **Offene Beobachtungen sichten:** das Register `observations.md` entsteht **in diesem
  Slice** und startet leer (`— keine —`); keine offene Beobachtung zu berücksichtigen.

## 8. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc/Prozess führt. Berührt die *Harness/Prozess*-Doku (Roadmap,
Wellendokument, Register); greenfield-Form-Angleich an die adoptierte Baseline, ohne
Brownfield-Spec.

## 9. Closure-Notiz (nach `done/`)

Umgesetzt: die d-check-**Planning-Schicht** ist auf die v5.0.0-Baseline-Form gehoben.

- **D-2 Wellendokument:** `welle-67-baseline-v500-migration.md` flach angelegt — die
  **laufende** Welle nachträglich konform dokumentiert (Ziel · Start-Trigger slice-083 ·
  Closure-Trigger · Slices 084–091 **ohne** Status-Spalte · Out-of-Scope · §7 ausstehend
  bis Closure/Ruheort-Regel). Zustand = Verzeichnis, `git mv` → `done/` bei Welle-Closure.
- **D-3 Beobachtungs-Register:** `observations.md` als **stehende** Datei (`— keine —`;
  die verkörperten Lehren stehen in `AGENTS.md`/Gates/Konventionen).
- **D-1 Roadmap:** **sechs** Baseline-Abschnitte (Aktuelle Welle · Nächste Wellen ·
  Meilensteine · Abhängigkeitsgraph · Abgeschlossene Wellen · Historische
  Trigger-Verschiebungen); die „Vorgänger"-Prosa ins Closure-Log migriert; für
  welle-60…66 je eine minimale, **retroaktiv markierte** Ergebnis-Notiz (Nutzer-Entscheid)
  — mit **Reklassifikations-Vorbehalt** (61…66 waren überwiegend wellenlose
  Einzel-Slice-Arbeit).
- **Abnahme-Punkt 1 (Ruhe-Marker ↔ Template):** als **deklarierte Adaption** behalten —
  ein neuer Konventionsspeicher-Eintrag hält die Ruhe-Marker-/Prosa-Form fest
  (Feld-Details im Wellendokument, wellenloser Zustand explizit, gate-erzwungen vom
  `planning`-Modul); **kein** Modul-Umbau.
- **D-4:** die zwei Vorprüfungen sind in diesem Slice modelliert (§7); die
  AGENTS-Konventions-Verankerung koppelt an slice-089. **D-9:** der Trigger-Audit ist im
  welle-67-Closure-Trigger benannt (0 aktive Carveouts, latent).

**Design-Wendung:** der Slice wurde von „Doc-Form" auf **„Planning-Layer-Form"
umgewidmet** (Nutzer-Hinweis), nachdem sich zeigte, dass welle-67 **ohne** ihr
Baseline-Wellendokument lief — Roadmap-Form + Wellen-Lifecycle + Register sind **ein**
kohärenter Planning-Layer und gehören zuerst.

**Review:** unabhängiger Frischkontext-Review
(`docs/reviews/2026-08-02-slice-088-planning-layer-review.md`): **abnahmereif**, HIGH 0 /
MEDIUM 0 / LOW 0 / INFO 1 (Alt-Wellen-Reklassifikation → im Closure-Log-Vorbehalt
benannt). `make gates` + `make adr-check` grün, **keine** ADR/Spec/Code berührt.

**Lehre:** Scratchpad-/tmp-Skripte sind **session-flüchtig** — vorbereitete
Transformationen in der Session ausführen, nicht als „Plan für morgen" liegenlassen
(Nutzer fing es). Und: D-1/D-2/D-3 sind **ein** Planning-Layer, nicht trennen.

**Anschluss:** **slice-089** (Doc-Form: AGENTS.md D-11 + ADR-Trigger D-8), dann slice-090
(Review-Infrastruktur D-6/D-7/D-10) und slice-091 (Slice-`Status:`-Feld D-5). Danach
schließt welle-67.
