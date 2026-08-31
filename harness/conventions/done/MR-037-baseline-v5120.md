# MR-037 — Baseline-Pin-Hebung auf v5.12.0 (siebter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(keine; Pin-Fortschreibung innerhalb des von
  [`MR-023`](../../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  festgelegten self-contained Bundle-Layouts)*
- **Datum:** 2026-08-26
- **Geltungsbereich:** [§Baseline](../../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen), die
  pin-gebundenen Verweise
  ([`MR-021`](../../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
  in [`AGENTS.md`](../../../AGENTS.md), [`harness/README.md`](../../README.md), den
  aktiven `MR-*`-Dateien, [`.harness/skills/reviewer.md`](../../../.harness/skills/reviewer.md),
  der Prüf-Konfiguration, den Spec-Straten und den Planning-Docs
- **Adaption:** Der Baseline-Pin ist von `v5.11.0` auf **`v5.12.0`** gehoben
  (Kurs-Tag vom 2026-08-26) — die von
  [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt)
  vorgesehene Fortschreibung, siebter Nachtrag der Serie; ersetzt
  [`MR-030`](../../conventions.md#mr-030) nach dessen eigenem
  Auflösungs-Trigger. Kein Layout-Wechsel: dasselbe self-contained Bundle,
  dasselbe Materialisierungs-Skript, unverändertes Pfadschema.

  **Vier Kurs-Wellen, und alle vier sind die Antwort auf den Konsumenten-CR
  dieses Repos** (abgelegt nach [`MR-036`](../../conventions.md#mr-036)). Das
  Delta ist gemessen, nicht gezählt: von **52** Bundle-Dateien unterscheiden
  sich **29**, und die Subtraktion steht hier statt im Kopf des Lesers — **eine**
  ist das Manifest `SHA256SUMS`, **22** ändern ausschließlich den
  Versions-Stempel (maschinell geprüft: kein einziges Zwei-Zeilen-Delta trägt
  etwas anderes), **eine** ist `regelwerk/README.md` mit Stempel plus
  `**Stand:**`-Zeile. Bleiben **fünf** mit echtem Regel-Inhalt. Die zwei
  Nicht-Markdown-Vorlagen des Bundles (`templates/.d-check.yml`,
  `templates/Makefile`) sind **unverändert** — eigens geprüft, weil die
  Delta-Schleife über `*.md` lief und sie sonst ungesehen durchgefallen wären.

  **Fünf Dateien tragen echten Regel-Inhalt:**

  | Datei | Umfang |
  |---|---|
  | `regelwerk/modul-05-planning-harness.md` | +11 |
  | `regelwerk/modul-13-quality-gates.md` | +8 |
  | `regelwerk/grundlagen-source-precedence.md` | +6 |
  | `regelwerk/modul-03-spec.md` | 5 für 2 |
  | `templates/AGENTS.template.md` | 3 für 2 |

  **Die Spiegel-Klassen aus [`BEO-008`](../../../docs/plan/planning/observations.md),
  je mit Zahl:** Pfad-Verweise **118** gesamt, davon **42 lebend** gehoben und
  76 eingefroren; Release-/Tree-URLs **5 lebend**, alle gehoben; nackte
  Nennungen **24**, davon **4 gehoben** und 20 als Aussagen über die
  Vergangenheit stehengelassen (`Ausgelöst durch Baseline-Stand`, die
  CR-Dokumente, das Drift-Log, die abgelöste Vorgänger-Adaption).

  **Die Messvorschrift gehört zur Zahl**, sonst ist sie nicht nachrechenbar:
  gezählt wurden **Vorkommen** (nicht Dateien) in den von `git ls-files`
  geführten `*.md`/`*.yml`/`*.sh`/`Makefile`-Dateien, **vor** der Hebung, unter
  Ausschluss der vier eingefrorenen Präfixe (`docs/plan/planning/done/`,
  `docs/reviews/`, `harness/conventions/done/`, `.harness/baseline/`); die
  dritte Klasse ist der Rest nach Abzug der Treffer von Klasse 1 und 2. Eine
  andere Abgrenzung ergibt eine andere Zahl — repo-weit ohne jeden Ausschluss
  sind es 60.

  **Die erste Zählung der ersten Klasse war zu eng** und hätte den
  Reviewer-Skill auf einen gelöschten Baum zeigen lassen: das Muster verlangte
  `.harness/baseline/…`, der Skill verweist relativ mit `../baseline/…`. Fünf
  Vorkommen, ein Dokument — gefunden, weil die Klasse-3-Liste sie zeigte.
- **Begründung:** Ein Adopter, der seine Baseline nicht auf einen Tag pinnt,
  auditiert gegen ein bewegliches Ziel; der Pin macht den Stand zitierbar und
  die Abweichung benennbar. Dass er **fortgeschrieben** wird statt zu altern,
  ist die Bedingung dafür, dass der Freshness-Audit überhaupt etwas zu
  vergleichen hat. Diese Hebung ist zusätzlich die erste, deren Inhalt das Repo
  selbst erbeten hat.
- **Löst auf:** [`MR-030`](../../conventions.md#mr-030)
- **Ausgelöst durch Baseline-Stand:** v5.12.0
- **Auflösungs-Trigger:** der Kurs veröffentlicht einen neuen Release-Tag; dann
  Fortschreibung durch den nächsten Nachtrag zu
  [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt).
