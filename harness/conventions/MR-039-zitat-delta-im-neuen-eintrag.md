# MR-039 — Ändert die Baseline einen zitierten Wortlaut, hält das ein neuer Eintrag fest — nicht der zitierende

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon gibt das **Prinzip** und lässt
  offen, wo das Delta landet.
  [`modul-02-harness-bootstrap.md` §Freshness-Audit](../../.harness/baseline/v5.18.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)   <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-02-harness-bootstrap.md:283-284 -->
  hält beides fest — *„bestehende werden nicht rückwirkend umgeschrieben"* und
  *„**Rückbau ist ein neuer Eintrag, kein Edit** … Die alte Zeile ist die
  historisch korrekte Aussage über den damaligen Zustand"*. Was er **nicht**
  sagt: wohin die Feststellung gehört, dass ein zitierter Baseline-Satz seit dem
  Bump anders lautet. Die Form-Frage tritt die Rangliste an diesen Speicher ab.
- **Datum:** 2026-08-26
- **Geltungsbereich:** wörtliche Zitate der vendorten Baseline in **allen
  lebenden** Dokumenten dieses Repos — Konventionsspeicher, `AGENTS.md`,
  `harness/README.md`, die Skills und die lebenden Planungs-Dokumente. Nicht
  `done/`, nicht `docs/reviews/`, nicht `conventions/done/`.
- **Adaption:** Beim Baseline-Bump wandert der **Link** auf den neuen Pin
  ([`MR-021`](../conventions.md#mr-021)). Der zitierte **Wortlaut bleibt
  unangetastet** — auch wenn er nicht mehr gilt. Hat sich der Wortlaut geändert,
  hält der **Bump-Eintrag** das fest: welches Dokument welchen Satz zitiert, wie
  er damals lautete und wie er jetzt lautet.

  **Warum nicht im zitierenden Dokument:** Der Kanon verlangt es andersherum.
  Ein bestehender Eintrag wird nicht rückwirkend umgeschrieben, und sein Zitat
  ist die historisch korrekte Aussage über den Stand, gegen den er geschrieben
  wurde. Wer es aktualisiert, macht aus dem Beleg eine Behauptung über eine
  Zukunft, die der Autor nicht kannte.

  **Warum es überhaupt festgehalten wird:** Ohne den Vermerk altert das Zitat
  **still**. Der Pfad daneben löst sauber auf, der Wortlaut gilt nicht mehr, und
  kein Gate sieht die Kombination — die vierte Spiegel-Klasse aus
  [`BEO-008`](../../docs/plan/planning/observations.md).

  **Bestand bei Einführung, gemessen:** In [`MR-033`](../conventions.md#mr-033)
  sind **zwei** Zitate durch die Klärung aus Kurs-Welle 96 veraltet, und beide
  bleiben dort unverändert stehen:

  | Quelle | Zitiert (Stand `v5.11.0`) | Seit `v5.12.0` |
  |---|---|---|
  | `templates/AGENTS.template.md` §3.4 | *„`spec/architecture.md` referenziert Modul-Pfade, aber **keine** Wellen, Slices, Commit-Hashes oder Closure-Daten."* | *„`spec/architecture.md` darf Pfade zu **Code-Modulen** referenzieren (`src/service/`), aber **keine** Wellen, Slices, Commit-Hashes oder Closure-Daten."* |
  | `regelwerk/modul-03-spec.md` §Ziel-Form: Architektur-Sicht | *„Sprach- und meilensteinfrei — referenziert Modul-Pfade, aber keine Wellen, Slices, Commit-Hashes oder Closure-Daten."* | *„Sprach- und meilensteinfrei — **darf** Pfade zu **Code-Modulen** referenzieren (`src/service/`), aber **keine** Wellen, Slices, Commit-Hashes oder Closure-Daten. Die Erlaubnis ist keine Pflicht: Eine Sicht, die ihre Komponenten nur über Rollen und `ARC-*` führt, ist ebenso konform."* |

  **Die Klärung bestätigt die Lesart, unter der
  [`MR-033`](../conventions.md#mr-033) geschrieben wurde** — gemeint sind
  Code-Modul-Pfade. Der Eintrag bleibt damit eine Verschärfung und gilt
  unverändert weiter.
- **Begründung:** Ein Verweis, der zugleich **zitiert**, hat zwei Hälften, und
  beim Bump wandert nur eine. Die andere gehört nicht mitgezogen — sie gehört
  **vermerkt**, an einer Stelle, die den Bump ohnehin beschreibt.
- **Löst auf:** [`MR-038`](../conventions.md#mr-038)
- **Ausgelöst durch Baseline-Stand:** v5.12.0
- **Auflösungs-Trigger:** der Kanon benennt selbst, wo ein geändertes
  Baseline-Zitat vermerkt wird. Dann gilt seiner.
