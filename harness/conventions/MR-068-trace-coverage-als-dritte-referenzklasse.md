# MR-068 — `trace.coverage` ist eine dritte, entlastende Referenzklasse neben dem Kanon-Vorschlag

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:**
  [`grundlagen-traceability.md` §Die zweite Richtung: Anforderung → Beleg](../../.harness/baseline/v6.5.0/regelwerk/grundlagen-traceability.md)
  — der Kanon **schlägt vor**, dass der Slice entlastet und die ADR als Spalte
  danebensteht, und verlangt für einen anderen Schnitt ausdrücklich eine
  Deklaration, *„wie jede Abweichung von der Baseline"*. Er nennt als Beispiel
  genau diesen Fall: eine kuratierte Nachweis-Datei als entlastende Quelle.
- **Datum:** 2026-09-07
- **Geltungsbereich:** `trace.coverage` in
  [`.d-check.yml`](../../.d-check.yml),
  [`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in),
  `make completeness-check`
- **Adaption:** Neben dem Kanon-Schnitt — **Slice entlastet, ADR ist Spalte** —
  führt dieses Repo eine **dritte, opt-in** Referenzklasse: eine kuratierte
  Coverage-Quelle. Ist sie konfiguriert, deckt sie eine Anforderung ebenso ab
  wie ein referenzierender Slice.

  **Der Kanon-Schnitt selbst ist unverändert übernommen**, und das ist keine
  Selbstverständlichkeit: Eine bloße ADR-Referenz **entlastet nicht**
  ([`DC-FA-CLI-011`](../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
  §Out-of-Scope sagt es wörtlich), die `adrs:`-Untersektion konfiguriert eine
  **Anzeigespalte**. Wer sie für entlastend hält, verschiebt die Schwelle von
  `make completeness-check` — eine Gate-Senkung ohne ADR.

  **Warum die dritte Klasse:** Nicht jede Anforderung wird von einem Slice
  eingelöst. Eine Qualitätszusage kann durch einen Test oder eine kuratierte
  Nachweis-Tabelle gedeckt sein, ohne dass je ein Slice sie im Namen führte;
  ohne die Klasse bliebe sie dauerhaft Waise und der Waisen-Zähler damit
  unbrauchbar als Signal.

- **Grenze:** Die Klasse ist **opt-in und kuratiert** — sie deckt, was jemand
  eingetragen hat, und driftet wie jede kuratierte Kante. Das teilt sie mit
  `matrix.exclude-sections` und `trace.cross-consistency`s `exclude-req`, und
  es ist der Preis, den der Kanon meint, wenn er die Deklaration verlangt: Eine
  entlastende Quelle ist eine **Setzung**, keine Naturgesetzlichkeit.
- **Begründung:** Der Kanon kannte den Vorschlag bis `v6.5.0` nicht; die
  Klasse existierte vorher und war nur nicht als Abweichung geführt. Mit der
  neuen Sektion ist sie deklarationspflichtig geworden — dieser Eintrag löst
  das ein, ohne am Verhalten etwas zu ändern.
- **Auflösungs-Trigger:** der Kanon nimmt die kuratierte Nachweis-Quelle in
  seinen Vorschlag auf — dann gilt seiner, und die Abweichung entfällt.
