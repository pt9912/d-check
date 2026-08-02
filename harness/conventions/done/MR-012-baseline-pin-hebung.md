# MR-012 — Baseline-Pin-Hebung (Nachtrag zu MR-011)

- **Status:** Accepted
- **Datum:** 2026-06-18
- **Geltungsbereich:** [§Baseline](../../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen),
  [`MR-006`](../../conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs),
  [`MR-008`](../../conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)
- **Adaption:** Der von
  [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt) vorgesehene
  Nachtrag: der Baseline-Pin ist von `v1.2.1` auf den Release-Tag
  **`v1.3.0`** gehoben (Stand-Zeile, Lehrmaterial-Repo-Link,
  `agents-regelwerk.md`-Raw-URL sowie die gespiegelten Pointer in
  `AGENTS.md` und `harness/README.md`). Die beiden an die Kurs-Quelle
  gekoppelten Trigger sind dabei re-evaluiert, am Tag `v1.3.0` verifiziert:
  - [`MR-006`](../../conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)-Trigger
    **eingetreten:** die Spec-Straten-Vorlagen verweisen im bindenden Text
    nicht mehr abwärts auf ADRs (`spezifikation.template.md`: „kein
    ADR-Rückzeiger hier"; `architecture.template.md`: „kein ADR-Bezug in
    dieser Sicht", §ADR-Index entfernt). Die verbleibende ADR-Spalte nur
    in der Historie ist legitime Provenance (Regel 5), deckungsgleich mit
    dem `matrix`-`exclude-sections`-Default dieses Repos.
  - [`MR-008`](../../conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)-Kurs-Inkonsistenz
    **behoben:** die ADR-ID ist upstream durchgängig vierstellig
    (`ADR-NNNN`) — `conventions.template.md`, die `Bezug:`-Zeile des
    ADR-Templates und die `ADR-Bindung`-Klasse des Regelwerks tragen
    v1.3.0 vierstellig. d-checks eigene vierstellige Wahl ist damit
    baseline-konform statt vorgezogene Abweichung.
- **Begründung:** Die Tag-Hebung erfolgt per Nachtrags-MR (kein
  Überschreiben des akzeptierten Eintrags), wie in
  [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt) vorgesehen; die
  gekoppelten Trigger werden dabei re-evaluiert. v1.2.1 → v1.3.0 ist
  inhaltlich nur der ADR-Stelligkeits-Fix (eine Zeile im
  Regelwerk-Bundle) plus Tag-Rewrite — keine sonstige Regeländerung.
- **Auflösungs-Trigger:** permanent als Provenienz; die nächste
  Baseline-Version wird erneut per Nachtrags-MR gehoben.
