# MR-023 — Baseline-Pin-Hebung auf v5.0.0 samt self-contained Bundle-Layout

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`grundlagen-harness-dateien.md` §Template-Schichtung](../../.harness/baseline/v5.15.0/regelwerk/grundlagen-harness-dateien.md#template-schichtung--was-der-rumpf-trägt-und-was-der-kommentar)
- **Datum:** 2026-08-01
- **Geltungsbereich:** [§Baseline](../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../conventions.md#adoptierte-konventions-quellen), das
  Materialisierungs-Skript
  [`tools/harness/fetch-baseline-cache.sh`](../../tools/harness/fetch-baseline-cache.sh),
  die pin-gebundenen Verweise
  ([`MR-021`](../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
  in [`AGENTS.md`](../../AGENTS.md), [`harness/README.md`](../README.md),
  [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) und den
  Planning-Docs; Nachträge zu
  [`MR-017`](../conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen),
  [`MR-018`](../conventions.md#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates),
  [`MR-019`](../conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017),
  [`MR-022`](../conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019)
- **Adaption:** Der Baseline-Pin ist von `v1.4.0` auf **`v5.0.0`** gehoben — der
  von [`MR-011`](../conventions.md#mr-011--baseline-auf-release-tag-gepinnt) vorgesehene Nachtrag,
  hier über **zwei weitere Major-Sprünge** (`v4.0.0`, `v5.0.0`). Anders als
  [`MR-016`](../conventions.md#mr-016--baseline-pin-hebung-zweiter-nachtrag-zu-mr-011) (reine
  Template-Pflege, Regelwerk unverändert) ist dies eine **inhaltliche Migration**:
  das Regelwerk ist umstrukturiert (Grundlagen 3→8 Dateien, Module umbenannt) und
  das Release-Asset-Layout hat sich geändert. Diese MR deckt **Etappe A**
  (Vendoring + Pin/Pointer):
  - **Self-contained Bundle (beide Bäume vendored).** Das Release liefert seit
    `v5.0.0` **ein** Asset `lab-regelwerk.zip` mit **beiden** Bäumen
    (`regelwerk/` + `templates/`); `lab-templates.zip` entfällt. Beide werden
    jetzt **committet vendored** unter
    `.harness/baseline/<tag>/{regelwerk,templates}/`. Das hebt den Cache-Zweig
    aus
    [`MR-017`](../conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)/[`MR-018`](../conventions.md#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates)
    auf (Templates nicht länger ephemerer `.harness/cache/`-Stand); die
    **vollständige Neufassung** der Cache-/Vendoring-Einträge auf diese Form ist
    Teil der späteren Konventionsspeicher-Etappe.
  - **Skript aufs Bundle gehoben.** `fetch-baseline-cache.sh` entpackt das Bundle
    **tolerant** (Regelwerk am `modul-00`-Marker, Templates als Geschwister),
    prüft eine **Under-Copy-Barriere** (Quelle == vendored) und schreibt das
    `SHA256SUMS`-Manifest über den **tatsächlichen** Bestand **beider** Bäume.
    Der `--check-latest`-Currency-Teil liest jetzt die Release-**Liste** (statt
    `/releases/latest`), der Content-Drift-Teil vergleicht **beide** Bäume; die
    Wortlaut-Angleichung der
    [`MR-019`](../conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)/[`MR-022`](../conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019)-Prosa
    folgt in der Konventionsspeicher-Etappe.
  - **Entfallene Quellzeiger umgeschrieben.** `agents-regelwerk.md` (im Kurs
    zurückgezogen) und die Kurs-`grundlagen/konventionen.md` (in acht
    `grundlagen-*`-Dateien aufgeteilt) werden **nicht** auf tote Ziele
    retargetet, sondern umgeschrieben: die committet vendorte `regelwerk/` (mit
    `README.md` als Index) **ist** die Agenten-Lese-Form; die §Referenz-Richtung
    liegt jetzt in `grundlagen-referenz-richtung.md`.
  - **Eingefrorene Historie via Tombstone.** Drei immutable/`done/`-Verweise auf
    den entfernten `v1.4.0`-Pfad
    ([`ADR-0022`](../../docs/plan/adr/0022-matrix-token-richtung-provenance-marker.md),
    `slice-080`, `slice-081`) werden über das geteilte Referenz-Ventil
    `ignore-refs`
    ([`DC-FA-REF-001`](../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus))
    referenz-weit von der Existenz-Prüfung ausgenommen — kein Editieren
    eingefrorener Doku.
- **Begründung:** Auftraggeber-Vorgabe (2026-07-25): vollständige Migration nach
  `v5.0.0`, der Baseline-Default sticht die repo-lokale Adaption. Etappen-Schnitt
  und Umgang mit den historischen Verweisen sind in der abgenommenen Analyse
  (`slice-083`) begründet. Die Konventionsspeicher-Restrukturierung (Index +
  Datei je MR mit neuen Pflichtfeldern) und die Form-Konformität folgen als
  eigene Etappen.
- **Auflösungs-Trigger:** permanent, solange die Baseline extern gepinnt und
  self-contained (beide Bäume) vendored wird. *(Die Überführung in die
  Datei-je-MR-Form des Konventionsspeichers ist mit Etappe C ausgeführt.)*
