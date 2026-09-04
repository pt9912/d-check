# MR-064 — Bündelung mehrerer Einzel-Slice-Archiv-Moves in einem Commit ist zulässig (Nachtrag zu MR-062)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Schärft
  [MR-062](MR-062-wellenloser-slice-archiv-move.md), dessen eigener
  „Grenze"-Absatz die Bündelungs-Frage ausdrücklich offenließ: *„Ob **mehrere**
  unabhängige Einzel-Slice-Archive in einem gemeinsamen Commit gebündelt
  werden, ist eine andere Frage — diese Regel begründet nur, dass **je Slice**
  genau ein Commit die Untergrenze ist, nicht ob mehrere solcher Commits
  zusammengefasst werden dürfen."*
- **Datum:** 2026-09-04
- **Geltungsbereich:** `tools/archive-wave`s Einzel-Slice-Modus
  (`ApplySlice()`) — die Frage, ob mehrere `-slice=<id> -apply`-Läufe in
  **einem** gemeinsamen Commit landen dürfen, statt je Slice einen eigenen
  Commit zu erzwingen.
- **Adaption:** Zulässig, wenn die Operation über alle gebündelten Slices
  mechanisch uniform und werkzeug-verifiziert bleibt — ohne inhaltliches
  Unterscheidungsurteil zwischen den einzelnen Läufen. Dieselbe Bündelungs-
  Klausel wie [MR-063](MR-063-eigenstaendiger-review-archiv-move.md) sie für
  den Review-Modus trägt, hier nachträglich für den Geltungsbereich benannt,
  den MR-062 selbst offenließ.
- **Begründung:** Die Praxis bestand bereits **zweimal**, ohne eine
  benannte Grundlage: slice-197 bündelte 45 Einzel-Slice-Archive in einem
  Commit (`0fe76a7`/`ab281cd`), slice-200 bündelte 7 weitere
  (`ce5fb50`) — beide Male, ohne dass MR-062 die Frage beantwortete, die
  sein eigener Text als offen benannte. Unabhängiger Review von slice-200
  (2026-09-04) fand dieselbe Lücken-Klasse, die MR-062/MR-063 bereits
  zweimal für andere Geltungsbereiche schließen mussten: eine gelebte
  Praxis ohne deklarierte Grundlage. Diese Regel liefert sie nachträglich,
  statt einen dritten Korrektur-Commit zu erzwingen.
- **Grenze:** Diese Regel beantwortet die Bündelungs-Frage **ausschließlich**
  für den Einzel-Slice-Modus. Sie ist keine Blankovollmacht für beliebige
  Content-Move-Bündelungen — jeder andere Modus (Welle, Review) trägt seine
  eigene Bündelungs-Aussage, soweit vorhanden ([MR-063](MR-063-eigenstaendiger-review-archiv-move.md)
  für den Review-Modus; der Wellen-Modus bündelt per Konstruktion bereits
  alle Slices einer Welle in einem Commit, MR-059).
- **Auflösungs-Trigger:** permanent, solange `tools/archive-wave`s
  Einzel-Slice-Modus je Lauf unabhängig aufrufbar bleibt.
