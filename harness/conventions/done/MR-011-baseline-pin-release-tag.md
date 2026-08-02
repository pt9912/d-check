# MR-011 — Baseline auf Release-Tag gepinnt

- **Status:** Accepted
- **Datum:** 2026-06-18
- **Geltungsbereich:** [§Baseline](../../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen)
- **Adaption:** Die Baseline-Aussage und die externen
  Konventions-Quellen waren bislang **unversioniert**: der Stand als
  bloßer Datumsstempel „Template-Set 2026-06" und die Quellen-Pointer
  auf den `main`-Branch (`.../ai-harness-course/main/...`). Beides ist
  auf den Release-Tag **`v1.2.1`** gepinnt — die Stand-Zeile, die
  Raw-URL von `agents-regelwerk.md` und der Lehrmaterial-Repo-Link.
  MR-000 bleibt unverändert (akzeptiert; Nachtrag daher als eigener
  Eintrag, analog
  [`MR-008`](../../conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)).
- **Begründung:** Ein `main`-Pointer ist genau die Doku-Drift, die das
  Regelwerk für Außen-Verweise verbietet (tag-gepinnte URLs als
  Anti-Drift-Default); ein unversionierter Stand macht die
  Baseline-Konformitäts-Aussage von
  [`MR-000`](../../conventions.md#mr-000--baseline-aussage) („keine inhaltlichen Adaptionen
  ggü. Baseline-Default") nicht mehr eindeutig prüfbar, sobald `main`
  weiterläuft. Anlass: User-Direktive 2026-06-18, das aktuelle
  Regelwerk (v1.2.1) zu befolgen.
- **Auflösungs-Trigger:** permanent als Mechanik; der konkrete Tag wird
  bei Adoption einer neueren Baseline-Version per Nachtrags-MR
  hochgezogen (dann auch Re-Evaluierung der an die Kurs-Quelle
  gekoppelten Trigger von
  [`MR-006`](../../conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)
  und [`MR-008`](../../conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)).
