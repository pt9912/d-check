# MR-036 — Die Antwort auf einen ausgehenden CR liegt bei ihrem CR

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon kennt den ausgehenden Change
  Request nicht (siehe [`MR-035`](../conventions.md#mr-035)) und damit auch
  nicht die Antwort darauf. Die Ablage-Frage ist eine Form-Frage, und die tritt
  die Rangliste an diesen Speicher ab
  ([`grundlagen-harness-dateien.md` §Konventionsspeicher](../../.harness/baseline/v5.18.0/regelwerk/grundlagen-harness-dateien.md#harnessconventionsmd-als-konventionsspeicher)).
- **Datum:** 2026-08-26
- **Geltungsbereich:** `docs/plan/cr/` — **eingehende Antworten** auf
  Change Requests, die dieses Repo an die adoptierte Baseline gerichtet hat.
- **Adaption:** Die Antwort liegt als datierte Datei neben ihrem CR, Form
  `<YYYY-MM-DD>-antwort-<gegenstand>.md`, und nennt im Kopf Absender, Datum,
  den Bezug auf den CR und das Ergebnis. Sie wird **wörtlich** abgelegt, nicht
  zusammengefasst: der Wert liegt in der Begründung der Gegenseite, und die
  überlebt keine Paraphrase.

  **Verhältnis zu [`MR-035`](../conventions.md#mr-035):** Der dortige Satz
  *„keine Pflicht, den Ausgang nachzutragen — wer das will, hängt es an den
  Slice, der den CR schreibt"* betrifft den **Vermerk** über den Ausgang
  (angenommen/abgelehnt) und gilt unverändert. Dieser Eintrag betrifft das
  **Antwort-Dokument** und deckt damit eine Klasse, die
  [`MR-035`](../conventions.md#mr-035)s Geltungsbereich — *Change Requests, die
  dieses Repo richtet* — nicht erfasst. Er schärft ihn nicht und löst ihn nicht
  ab; er steht daneben. Beide zusammen sagen: das Dokument gehört ins
  Verzeichnis, der Vermerk darf an den Slice.

  **Warum wörtlich und warum hier:** [`MR-035`](../conventions.md#mr-035)
  entstand aus dem Befund, dass beim vorigen Konsumenten-CR nicht mehr
  feststellbar war, welche Punkte mit welcher Begründung angenommen wurden. Eine
  Antwort, die nur als Zusammenfassung überlebt, lässt genau diese Frage wieder
  offen — und die Teil-Annahme mit **anderer Encodierung** ist der Fall, in dem
  die Begründung mehr wert ist als das Ergebnis.

  **Was die Ablage weiterhin nicht ist:** kein Register, keine Liste, kein
  Index, kein Zähler. Es gibt keine Pflicht, eine Antwort zu bekommen oder eine
  ausbleibende zu vermerken.
- **Ausgelöst durch Baseline-Stand:** v5.11.0
- **Auflösungs-Trigger:** der Kanon benennt selbst einen Ruheort für den
  ausgehenden Change Request und seine Antwort. Dann gilt seiner.
