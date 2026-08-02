# MR-020 — Baseline-Template-Propagation per Drift-Audit (template-frei bestätigt)

- **Status:** Accepted
- **Datum:** 2026-06-26
- **Geltungsbereich:** [`docs/plan/planning/README.md`](../../../docs/plan/planning/README.md)
  (§Lifecycle, Closure-Notiz §7), der Templates-Staging-Cache
  `.harness/cache/<tag>/templates/`, die Baseline-Pin-Bump-Prozedur; Nachtrag zu
  [`MR-018`](../../conventions.md#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates)
- **Adaption:** d-check bleibt **template-frei**
  ([`MR-014`](../../conventions.md#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template),
  [`MR-018`](../../conventions.md#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates))
  — bestätigt, nicht aufgehoben. Damit Baseline-Template-Verbesserungen den
  gelebten Haus-Stil dennoch erreichen, wird der dort angelegte
  Templates-Staging-Cache (`.harness/cache/<tag>/templates/`, ephemer — nur diese
  Rolle bleibt dem Cache nach
  [`MR-019`](../../conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017))
  beim **Baseline-Pin-Bump** als **Drift-Audit** ausgeführt: die staged
  Baseline-Skelette werden gegen den Haus-Stil verglichen; substanzielle
  Verbesserungen (neue Pflichtfelder, Pointer, …) werden eingezogen, rein
  stilistische Template-Eigenheiten bleiben außen vor. Erster eingezogener Fall:
  der §7-Closure-Notiz-**Steering-Loop-Eintrag** verweist auf die kanonische
  Definition im **vendorten** Regelwerk
  (`.harness/baseline/<tag>/regelwerk/grundlagen-klassifikation.md` §Steering
  Loop, in-repo auflösbar seit
  [`MR-019`](../../conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)),
  verankert in [`docs/plan/planning/README.md`](../../../docs/plan/planning/README.md).
- **Begründung:** Nutzer-Frage 2026-06-26: „Eigene Wege gehen heißt, Baseline-
  Probleme nochmals lösen." Stimmt — der Fork tauscht *automatische Vererbung*
  gegen *gelebten Haus-Stil*; bei Slice-50 sind die fertigen Slices die bessere
  Vorlage. Der Defekt war nicht der Fork, sondern dass das in
  [`MR-018`](../../conventions.md#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates)
  angelegte Staging nie *ausgeführt* wurde. Auslöser: der Kurs ergänzte
  `slice.template.md` §7 um einen Steering-Loop-Pointer, den d-check ohne Audit
  nicht mitbekam. Bewusst **kein Gate** — Template-Drift ist eine Erkenntnis-,
  keine Laufzeit-Eigenschaft (analog zur Regelwerk-Lesedisziplin); die Disziplin
  liegt am Pin-Bump.
- **Auflösungs-Trigger:** permanent, solange d-check template-frei ist; jeder
  Baseline-Pin-Bump führt den Drift-Audit aus.
