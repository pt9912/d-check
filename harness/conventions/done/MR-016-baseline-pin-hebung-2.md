# MR-016 — Baseline-Pin-Hebung (zweiter Nachtrag zu MR-011)

- **Status:** Accepted
- **Aufgelöst durch:** MR-023 (Pin-Hebung auf v5.0.0)
- **Datum:** 2026-06-25
- **Geltungsbereich:** [§Baseline](../../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen), die gespiegelten
  `lab-regelwerk.zip`-Pointer in [`AGENTS.md`](../../../AGENTS.md) §1 und
  [`harness/README.md`](../../README.md) §Guides sowie die kurs-gekoppelten Trigger
  [`MR-006`](../../conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs),
  [`MR-008`](../../conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage),
  [`MR-014`](../../conventions.md#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template)
- **Adaption:** Der von
  [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt) vorgesehene weitere
  Nachtrag: der Baseline-Pin ist von `v1.3.0` auf den Release-Tag **`v1.4.0`**
  gehoben. Aktualisiert: die §Baseline-Stand-Zeile, der Lehrmaterial-Repo-Link
  und die `agents-regelwerk.md`-Raw-URL + `lab-regelwerk.zip`-Asset-URL in
  §Adoptierte Konventions-Quellen, die gleichen beiden URLs in
  `harness/README.md` §Guides sowie — zip-only seit
  [`MR-015`](../../conventions.md#mr-015--auflösung-der-mr-012-pointer-drift-agentsmd-routet-spiegelt-nicht-mehr) —
  der `lab-regelwerk.zip`-Pointer samt `(v1.4.0)`-Stand-Markierung in
  `AGENTS.md` §1. Das **Regelwerk ist inhaltlich unverändert**: `kurs/de/`
  (Quelle von `agents-regelwerk.md` und Bundle) trägt zwischen `v1.3.0` und
  `v1.4.0` keinen Diff; der Pin wandert, der Inhalt nicht. Der gesamte
  v1.3.0→v1.4.0-Diff liegt in `lab/templates/`. Die beiden an die Kurs-Quelle
  gekoppelten Trigger sind dabei re-evaluiert, am Tag `v1.4.0` verifiziert,
  plus eine teilweise Auflösung von [`MR-014`](../../conventions.md#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template):
  - [`MR-006`](../../conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)-Trigger
    **unverändert eingetreten:** die Spec-Straten-Vorlagen
    (`spezifikation.template.md`, `architecture.template.md`) tragen keinen
    Diff zwischen `v1.3.0` und `v1.4.0`; der zu `v1.3.0` festgestellte Stand
    (kein ADR-Rückzeiger im bindenden Text) hält.
  - [`MR-008`](../../conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)-Konformität
    **gehalten:** das ADR-Template (`NNNN-titel.template.md`) hat sich geändert
    (Abschnitt „Verglichene Alternativen" von Option-A/B/C-Prosa auf eine
    Pro/Contra-**Tabelle**), die ADR-ID-Stelligkeit bleibt jedoch durchgängig
    vierstellig (`ADR-NNNN`).
  - [`MR-014`](../../conventions.md#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template)-Abweichung
    **teilweise aufgelöst:** die dort als Haus-Stil deklarierte ADR-„Verglichene
    Alternativen"-**Tabelle** (statt Option-A/B/C-Prosa) ist mit `v1.4.0`
    **Baseline-Default** geworden — dieser eine Punkt ist damit baseline-konform
    statt Abweichung. Die übrigen MR-014-Punkte (Slice-Abschnitts-Reihenfolge,
    Fitness-Function als Prosa-Bullets, zweispaltige „Geschichte") stehen
    unverändert: `slice.template.md` und die restliche ADR-Template-Struktur
    tragen keinen Diff. Die `@v1.3.0`-Template-Links im MR-014-Body bleiben als
    Vergangenheits-Aussage **unangetastet** (append-only, analog
    [`MR-015`](../../conventions.md#mr-015--auflösung-der-mr-012-pointer-drift-agentsmd-routet-spiegelt-nicht-mehr)).
- **Begründung:** Die Tag-Hebung erfolgt per Nachtrags-MR (kein Überschreiben
  des akzeptierten [`MR-012`](../../conventions.md#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011)-Eintrags),
  wie in [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt) vorgesehen. v1.3.0
  → v1.4.0 ist inhaltlich nur Template-Pflege: die ADR-Alternativen-Tabelle
  (s. o.), `harness.mk` hebt den **d-check-Konsumenten**-Digest-Pin
  (v0.8.0 → v0.23.0 — betrifft den Kurs als d-check-Nutzer, nicht d-checks
  eigene Konventionen) und eine Templates-`README`-Notiz; keine Regelwerk- oder
  sonstige Konventionsänderung.
- **Auflösungs-Trigger:** permanent als Provenienz; die nächste
  Baseline-Version wird erneut per Nachtrags-MR gehoben.
