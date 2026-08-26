# MR-038 — Ein Zitat der Baseline ist pin-gebunden wie ein Link, aber es wird ergänzt statt ersetzt (schärft MR-021)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon regelt den Fall nicht: er hält für
  `MR-<NNN>` die **Append-only-Disziplin** fest — *„Wird eine `MR-<NNN>` durch
  ein Baseline-Update gegenstandslos, wird sie nicht überschrieben"*
  ([`grundlagen-source-precedence.md` §Source Precedence](../../.harness/baseline/v5.12.0/regelwerk/grundlagen-source-precedence.md))
  — und schützt bei ADRs den **inhaltlichen** Kern
  ([`modul-04-adrs.md`](../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md)).
  Ob ein **Zitat** der Baseline zum geschützten Kern gehört oder zur
  pin-gebundenen Referenz, sagt keine der beiden Stellen. Die Form-Frage tritt
  die Rangliste an diesen Speicher ab.
- **Datum:** 2026-08-26
- **Geltungsbereich:** Zitate der vendorten Baseline in **aktiven** Einträgen
  unter `harness/conventions/` — insbesondere das Pflichtfeld
  `Ersetzt-Baseline-Regel`.
- **Adaption:** Beim Baseline-Bump wandert der **Link** auf den neuen Pin
  ([`MR-021`](../conventions.md#mr-021)). Hat sich der zitierte **Wortlaut**
  dabei geändert, wird das Zitat **ergänzt, nicht ersetzt**: die aktuelle
  Fassung wird genannt, und die Fassung, gegen die der Eintrag geschrieben
  wurde, bleibt mit ihrem Pin daneben stehen.

  **Warum nicht ersetzen:** Der Eintrag ist ein Geschichtsdokument. Ersetzt man
  den Wortlaut, behauptet er, gegen einen Text geschrieben worden zu sein, den
  es damals nicht gab.

  **Warum nicht einfach stehenlassen:** Das Feld `Ersetzt-Baseline-Regel` hat
  **zwei** Aufgaben — es nennt die Regel, von der dieses Repo abweicht (eine
  Aussage über den **gepinnten** Stand), und es hält die Grundlage der
  Entscheidung fest (eine Aussage über **damals**). Friert nur die zweite ein,
  taugt das Feld für den Freshness-Audit nicht mehr: er vergleicht gegen einen
  Wortlaut, der nicht mehr gilt.

  **Warum das kein Überschreiben ist:** Es wird nichts entfernt und keine
  Entscheidung geändert — es kommt etwas hinzu. Dieselbe Bewegung wie der
  `## Geschichte`-Anhang einer ADR, die der Kanon ausdrücklich zulässt.

  **Anwendungsfall bei Einführung:** genau **einer** — die Klärung aus
  Kurs-Welle 96 hat den Satz geändert, den
  [`MR-033`](../conventions.md#mr-033) zitiert. Alle übrigen Zitate aktiver
  Einträge sind gegen `v5.12.0` wortgleich (gemessen, normalisiert).
- **Begründung:** Ein Verweis, der zugleich **zitiert**, hat zwei Hälften, und
  beim Bump wandert nur eine. Die andere altert still: der Pfad löst sauber
  auf, der Wortlaut daneben gilt nicht mehr, und kein Gate sieht die
  Kombination — die vierte Spiegel-Klasse aus
  [`BEO-008`](../../docs/plan/planning/observations.md).
- **Ausgelöst durch Baseline-Stand:** v5.12.0
- **Auflösungs-Trigger:** der Kanon regelt selbst, wie ein Zitat der Baseline
  einen Bump übersteht. Dann gilt seiner.
