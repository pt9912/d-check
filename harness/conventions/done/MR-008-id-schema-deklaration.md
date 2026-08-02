# MR-008 — ID-Schema-Deklaration (Nachtrag zur Baseline-Aussage)

- **Status:** Accepted
- **Aufgelöst durch:** Baseline-Stand v5.0.0 (ID-Schema-Deklaration ist Baseline-Default, MR-000)
- **Datum:** 2026-06-11
- **Geltungsbereich:** gesamtes Repo (alle Artefakt-IDs und
  Traceability-Verweise)
- **Adaption:** Nachtrag der im Konventions-Template als Teil von
  MR-000 vorgesehenen ID-Schema-Deklaration, die in der
  Initial-Setzung fehlte (MR-000 ist akzeptiert und bleibt
  unverändert — Nachtrag daher als eigener Eintrag). Deklariertes
  Schema: Anforderungen `DC-FA-<BEREICH>-<NNN>` (Bereichskürzel, siehe
  [`MR-002`](../../conventions.md#mr-002--id-schema-mit-bereichskürzeln-ab-initialer-fassung))
  und `DC-QA-<NN>`; ADRs `ADR-NNNN` (vierstellig, gemäß
  Kurs-ADR-Vorlage `NNNN-titel.template.md` — das Konventions-Template
  nennt abweichend dreistellig `ADR-<NNN>`, eine
  Kurs-Template-Inkonsistenz; Korrektur in der Kurs-Quelle steht aus,
  analog [`MR-006`](../../conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs));
  Konventions-Adaptionen `MR-NNN`; Carveouts `CO-NNN` (bisher
  ungenutzt); Slices `slice-NNN`. Die `ids`-Selbstkonfiguration in
  [`.d-check.yml`](../../../.d-check.yml) prüft die Linkpflicht der
  Kennungen maschinell.
- **Begründung:** Eine undeklarierte ID-Systematik ist eine stille
  Setzung (gleiche Harness-Lüge-Klasse wie ein undeklariertes Gate);
  sichtbar geworden im Template-Vergleich (User-Review, 2026-06-11).
- **Auflösungs-Trigger:** permanent. *(Die in der Adaption vermerkte
  Kurs-Template-Inkonsistenz — `conventions.template.md` dreistellig vs.
  ADR-Vorlage vierstellig — ist mit Baseline `v1.3.0` behoben (ADR-ID
  upstream vierstellig vereinheitlicht); siehe
  [`MR-012`](../../conventions.md#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011).)*
