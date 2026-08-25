# MR-034 — Die Referenzmatrix bewacht auch die Kante ADR → Welle

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Dieser Eintrag **schärft**
  [`MR-006`](../conventions.md#mr-006) — dessen §Scope-Grenze (C-4) nennt die
  Kante **ADR→Welle** als *„bewusst unbewacht"* und begründet das ausdrücklich
  damit, dass *„d-check Carveout/Welle/Roadmap **nicht als `matrix`-Klassen**
  modelliert; eine Erweiterung wäre ein eigener Change"*. Genau dieser Change
  ist eingetreten.
- **Datum:** 2026-08-23
- **Geltungsbereich:** der `matrix`-Block in
  [`.d-check.yml`](../../.d-check.yml) — die Klasse `welle` und die Regel
  `{from: adr, to: welle, allow: false}`.
- **Adaption:** `MR-006` hält die Scope-Grenze **konditional**: unbewacht,
  *weil* keine Klasse existiert. Mit der Klasse fällt die Bedingung weg, und
  die Kante ist zu bewachen. Die kanonische Referenz-Matrix führt sie als
  flaches Verbot — **ohne** den Provenance-Marker-Ausweg, den `ADR→Slice`
  offenlässt.

  **Was der Eintrag nicht tut:** `MR-006` wird nicht überschrieben. Seine
  Scope-Grenze bleibt für **Carveout** und **Roadmap** unverändert gültig; nur
  die Wellen-Kante verlässt sie.

  **Bestands-Ausnahme, gemessen statt gesetzt:** die Kante trifft heute genau
  **eine** nicht-grandfatherte ADR — eine `Accepted` und damit immutable, deren
  Nennung reine Provenance ist (*„der Pin ist auf `<tag>` gehoben"*). Der
  Marker-Ausweg steht ihr nicht offen, weil ihr Kern nicht mehr geändert werden
  darf. Sie kommt deshalb in `matrix.exempt-paths` — aus demselben Grund wie
  die Reihe `0001`–`0021` davor. Der Preis ist gemessen: ihr Körper trägt
  keinen weiteren Token, den die Ausnahme mit stummschaltete.
- **Ausgelöst durch Baseline-Stand:** v5.11.0
- **Auflösungs-Trigger:** permanent, solange `welle` eine `matrix`-Klasse ist.
  Fällt die Klasse weg, fällt auch dieser Eintrag und `MR-006`s ursprüngliche
  Scope-Grenze gilt wieder unverändert.
