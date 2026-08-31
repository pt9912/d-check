# MR-045 — `AGENTS.md` und `harness/README.md` tragen keine Slice-Verweise (schärft MR-013)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon weist `AGENTS.md` die **Hard
  Rules und Pointer** zu und die zeitliche Schicht der Planung
  ([`modul-09-implementierung.md` §AGENTS.md-Regeln](../../.harness/baseline/v5.15.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9),
  [`modul-05-planning-harness.md` §Lifecycle als State Machine](../../.harness/baseline/v5.15.0/regelwerk/modul-05-planning-harness.md#lifecycle-als-state-machine));
  dass ein Briefing **keinen** Slice nennt, folgt daraus, steht aber nirgends
  ausdrücklich.
- **Datum:** 2026-08-27
- **Geltungsbereich:** [`AGENTS.md`](../../AGENTS.md) und
  [`harness/README.md`](../README.md).
- **Adaption:** Beide Dateien tragen **keine** `slice-<NNN>`-Verweise. Was ein
  Slice erarbeitet hat, steht dort als Regel, als Grenze oder als Zeiger auf
  eine ADR bzw. einen `MR`-Eintrag — nie als Verweis auf das Planungs-Artefakt.

  **Folge für [`MR-013`](../conventions.md#mr-013):** dessen Aufzählung der
  gekoppelten Verweise nennt neben der Roadmap auch `AGENTS.md` §4 und
  `harness/README.md` §Sensors. Für diese beiden Ziele ist die Klausel **leer**
  — sie können keinen Slice-Pfad tragen, also gibt es beim Lifecycle-Move dort
  nichts nachzuziehen. Die Klausel gilt unverändert für alle übrigen Ziele
  (Roadmap, andere Slices, Wellendokumente, `docs/plan/`).
- **Begründung:** Ein Briefing, das einen Slice nennt, altert mit ihm: der Slice
  wandert durch vier Verzeichnisse und ist am Ende geschlossen, während die
  Regel weiter gilt. Der Verweis wird damit entweder zur Pflege-Last oder zur
  Lüge — und im Zweifel liest ein Lauf die Regel als so vorläufig wie ihren
  Beleg.

  **Warum als eigener Eintrag und nicht als Korrektur an
  [`MR-013`](../conventions.md#mr-013):** dessen Aufzählung war für den Stand
  richtig, gegen den sie geschrieben wurde. Einträge werden nicht rückwirkend
  umgeschrieben
  ([`modul-02-harness-bootstrap.md` §Freshness-Audit](../../.harness/baseline/v5.15.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)).
- **Auflösungs-Trigger:** permanent.
