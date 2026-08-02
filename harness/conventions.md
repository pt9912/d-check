# Harness-Konventionen

## Purpose

Diese Datei deklariert die *repo-lokalen* Strukturregeln dieses Repos
gegenüber der adoptierten Harnesskonvention (Baseline):

- **Adaptionen** ggü. der Baseline (mit Begründung und Auflösungs-Trigger).
- **ID-Schema-Deklaration** — welches Präfix-Schema dieses Repo nutzt
  ([`MR-008`](#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)).
- **Zusatzklassen-Deklarationen** für die Sensors-Bindung.
- **Modus-Deklarationen** pro Sub-Area (Greenfield / Brownfield /
  Hybrid) inklusive Konvergenz-Auftrag bei BF.

Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt
die kanonische Quelle (Source Precedence, siehe
[`README.md`](README.md)). Diese Datei ist konformitätsbringend für
*Form*-Fragen, nicht autoritativ über Inhalt.

## Baseline

- **Konvention:** AI-Harness-Kurs
- **Stand:** [`v5.0.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v5.0.0)
  (Release-Tag, gepinnt 2026-08-01; ursprünglich unversioniert adoptiert
  als „Template-Set 2026-06", nachgezogen mit
  [`MR-011`](#mr-011--baseline-auf-release-tag-gepinnt), von v1.2.1 auf
  v1.3.0 gehoben mit [`MR-012`](#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011),
  von v1.3.0 auf v1.4.0 mit
  [`MR-016`](#mr-016--baseline-pin-hebung-zweiter-nachtrag-zu-mr-011), von v1.4.0
  auf v5.0.0 (zwei weitere Majors) mit
  [`MR-023`](#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout))
- **Datum der Adoption:** 2026-06-10

## Adoptierte Konventions-Quellen

- **Extern (Lehrmaterial):**
  [`ai-harness-course@v5.0.0`](https://github.com/pt9912/ai-harness-course/tree/v5.0.0)
  (Templates: `lab/templates/`, Konventionen: die nach Themen aufgeteilten
  `kurs/de/grundlagen/grundlagen-*.md` — die frühere Sammeldatei
  `konventionen.md` entfällt seit dem Regelwerk-Umbau).
- **Extern (Regelwerk):** das operative Betriebsregelwerk für Code-Agenten (ohne
  Didaktik) ist seit dem Kurs-Umbau **kein** separates `agents-regelwerk.md` mehr,
  sondern das nach Modulen und Grundlagen aufgeteilte Regelwerk selbst,
  ausgeliefert als self-contained Release-Bundle
  [`lab-regelwerk.zip`](https://github.com/pt9912/ai-harness-course/releases/download/v5.0.0/lab-regelwerk.zip)
  (ein Asset, **beide** Bäume). Derivativ, bei Konflikt gilt das Lehrmaterial;
  pro Session lädt man **einen** Abschnitt, ohne das ganze Regelwerk im Kontext
  zu halten
  ([`MR-023`](#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)).
- **Lokale Lese-Form:** **beide** Bäume sind seit `v5.0.0` **committet vendored**
  unter `.harness/baseline/<tag>/{regelwerk,templates}/` (aktuell
  [`v5.0.0`](../.harness/baseline/v5.0.0/regelwerk/); entpacktes, self-contained
  `lab-regelwerk.zip`) samt `.harness/baseline/<tag>/SHA256SUMS`-Manifest über
  beide Bäume — netzlos auf jedem Checkout präsent
  ([`MR-019`](#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)
  vendored das Regelwerk; mit dem self-contained Bundle
  [`MR-023`](#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  sind auch die **Templates** vendored statt ephemerer Cache — die Neufassung der
  Cache-Einträge
  [`MR-017`](#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)/[`MR-018`](#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates)
  folgt in der Konventionsspeicher-Etappe). Aus dem Selbst-Scan ausgenommen;
  materialisiert per
  [`tools/harness/fetch-baseline-cache.sh`](../tools/harness/fetch-baseline-cache.sh).
- **Konventions-Vorbilder (Implementierung):**
  [`u-boot`](https://github.com/pt9912/u-boot) — Hexagon-Ordnerkonvention
  ([ADR-0005](../docs/plan/adr/0005-modul-layout-hexagon-ordner.md)),
  Dockerfile-/Makefile-Muster, Pin-Politik;
  [`b-cad`](https://github.com/pt9912/b-cad) — Gate-Nachweis- und
  Hook-Mechanik ([`MR-004`](#mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild)/[`MR-005`](#mr-005--härtung-ggü-b-cad-inhaltsbasierter-gate-nachweis-sub-shell-prüfung));
  [`d-migrate`](https://github.com/pt9912/d-migrate) — Ursprung des
  vendorten Bootstrap-Sensors ([`MR-003`](#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh), aufgelöst).
- **In-Repo (verkörperte Form):** `AGENTS.md`, `harness/README.md`,
  Verzeichniskonvention `spec/` + `docs/plan/` + `harness/`.

## Adaptions-Block

### MR-000 — Baseline-Aussage

- **Datum:** 2026-06-10
- **Geltungsbereich:** gesamtes Repo
- **Adaption:** keine inhaltlichen Adaptionen ggü. Baseline-Default
  für Verzeichniskonvention, Lifecycle-Regeln, Carveout-Disziplin.
- **Begründung:** Initial-Setzung. Spätere Adaptionen werden als
  `MR-<NNN>` nachgetragen.
- **Auflösungs-Trigger:** permanent.

### MR-001 — Source Precedence mit eigener Spezifikations-Schicht

- **Datum:** 2026-06-10
- **Geltungsbereich:** [`harness/README.md` §Source precedence](README.md#source-precedence)
- **Adaption:** Die Source-Precedence-Tabelle führt
  `spec/spezifikation.md` als eigenen **Rang 2** zwischen Lastenheft
  (Rang 1) und Architektur (Rang 3). Der Kurs-Default setzt zwei
  Spec-Ränge; dieses Repo nutzt drei. Beide Dateien der Ränge 2–3
  entstehen mit slice-002; bis dahin sind sie in den Tabellen als
  „geplant" markiert und nicht verlinkt.
- **Begründung:** Spec-Stratifizierung mit drei Spec-Dateien; die
  ADR-Schärfungs-Regel („ADR darf Spezifikation schärfen, nicht
  Lastenheft") soll strukturell sichtbar sein.
- **Auflösungs-Trigger:** permanent.

### MR-002 — ID-Schema mit Bereichskürzeln ab initialer Fassung

- **Datum:** 2026-06-10
- **Geltungsbereich:** [`spec/lastenheft.md`](../spec/lastenheft.md), alle Traceability-Verweise
- **Adaption:** Funktionale Anforderungen verwenden von Beginn an
  Bereichskürzel: `DC-FA-<BEREICH>-<NNN>` (z. B. [`DC-FA-LINK-001`](../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links))
  statt des zweistelligen Kurs-Defaults `<PREFIX>-FA-<NN>`.
  Nichtfunktionale Anforderungen bleiben beim Kurs-Default
  (`DC-QA-<NN>`).
- **Begründung:** Das Lastenheft konsolidiert zwölf Quell-Tools und hat
  dadurch von Anfang an viele Funktionsbereiche; das Kurs-Beispiel
  (DocSearch) zeigt, dass eine spätere Schema-Migration teurer ist als
  ein Bereichsschema ab Welle 1.
- **Auflösungs-Trigger:** permanent.

### MR-003 — Vendorter Bootstrap-Sensor `tools/verify-doc-refs.sh`

- **Datum:** 2026-06-10
- **Geltungsbereich:** `tools/verify-doc-refs.sh` (gelöscht mit slice-004, siehe [`MR-007`](#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)), `make doc-check` <!-- d-check:ignore (historisch: gelöscht) -->
- **Adaption:** Bis `d-check` sich selbst prüfen kann, läuft
  `make doc-check` über ein aus `d-migrate` vendortes Shell-Skript
  (Markdown-Linkziel-Prüfung). Das ist Fremd-Code ohne eigene Spec in
  diesem Repo (Sub-Area in BF, siehe Modus-Tabelle).
- **Begründung:** Ein Doku-Repo ohne Doku-Sensor wäre ein blinder
  Bootstrap; das vendorte Skript ist dependency-frei (bash/awk) und
  deckt den Kern von
  [`DC-FA-LINK-001`](../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
  ab.
- **Auflösungs-Trigger:** slice-004 — `make doc-check` läuft über
  `d-check` selbst (Dogfooding), das Skript wird gelöscht
  ([Slice-Plan](../docs/plan/planning/done/slice-004-anchors-modul-und-dogfooding.md)).

### MR-004 — Gate-Nachweis-Mechanik und `.claude`-Hooks nach b-cad-Vorbild

- **Datum:** 2026-06-10
- **Geltungsbereich:** [`tools/harness/`](../tools/harness/), [`.claude/`](../.claude/), `make record-gates`
- **Adaption:** Übernahme der Working-Tree-Hash-Mechanik
  (`record-gates` als letzter `gates`-Prerequisite, Stop-Hook
  vergleicht den Hash) und der `.claude`-Hooks (PreToolUse-Guard,
  Stop-Gate) aus dem Repo `b-cad`.
- **Begründung:** Bewährte Mechanik gegen „Erfolgsmeldung ohne
  Gate-Lauf"; keine Logik-Dopplung zwischen Makefile und Hook.
- **Auflösungs-Trigger:** permanent.

### MR-005 — Härtung ggü. b-cad: inhaltsbasierter Gate-Nachweis, Sub-Shell-Prüfung

- **Datum:** 2026-06-10
- **Geltungsbereich:** [`tools/harness/working-tree-hash.sh`](../tools/harness/working-tree-hash.sh), [`.claude/hooks/`](../.claude/hooks/)
- **Adaption:** Zwei Abweichungen von der per [`MR-004`](#mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild)
  übernommenen b-cad-Mechanik:
  (a) Der Working-Tree-Hash ist **inhaltsbasiert** (sha256 über alle
  getrackten + untracked Dateiinhalte) statt diff-basiert. Damit gilt
  der Gate-Nachweis über Commits hinweg (gleicher Inhalt = gleicher
  Hash), und ein Commit *ohne* Gate-Lauf macht den Stop-Hook nicht
  mehr grün. Restlücke bleibt: frischer Klon bzw. gelöschter
  `.harness`-State mit cleanem Tree wird freigegeben — dort ist CI das
  Netz.
  (b) Der PreToolUse-Guard prüft Sub-Shell-Strings (`bash -c "…"`,
  `sh -c '…'`) rekursiv (Tiefe ≤ 3, darüber fail-closed).
- **Begründung:** Review-R2-Beobachtungen (User): Commit-Bypass des
  Stop-Hooks und Guard-Umgehung via `bash -c`.
- **Auflösungs-Trigger:** permanent. Rückport beider Härtungen nach
  b-cad steht aus.

### MR-006 — Referenzrichtung: Spec-Straten verweisen nie abwärts auf ADRs

- **Datum:** 2026-06-10
- **Geltungsbereich:** `spec/*.md`, [`AGENTS.md` §3.4](../AGENTS.md#34-architektur-sprach-meilensteinfrei-spec-straten-nie-abwärts)
- **Adaption:** Das adoptierte Template-Set 2026-06 sah ADR-Verweise
  in `spezifikation.md` (ADR-Spalte in Defaults/Historie) und
  `architecture.md` (ADR-Spalte, §ADR-Index) vor. Das ist als Fehler
  der Kurs-Vorlagen identifiziert; die Korrektur erfolgt in der
  Kurs-Quelle (Entscheidung User, 2026-06-10). d-check zieht vor:
  **kein Spec-Stratum (Rang 1–3) verweist abwärts auf ADRs oder
  Planning-Artefakte**; Traceability läuft ausschließlich über die
  `Schärft:`-Felder der ADRs (aufwärts). Die spätere
  matrix-Selbstkonfiguration kodiert das als
  `{from: spec-strata, to: adr/slice, allow: false}`.
- **Begründung:** Stable Dependencies — die Lösungsbeschreibung muss
  Entscheidungs-Revisionen (Supersede) überleben, ohne selbst
  angefasst zu werden; die Richtung der Begründung ist ADR → Spec,
  nie umgekehrt. Konsistent mit u-boots Checker („view spec may not
  link down").
- **Auflösungs-Trigger:** sobald das Kurs-Template-Set korrigiert
  ist, wird dieser Eintrag zur reinen Baseline-Konformität (bleibt
  als Provenienz stehen). *(Eingetreten mit Baseline `v1.3.0`;
  Spec-Straten-Vorlagen korrigiert — siehe
  [`MR-012`](#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011).)*

### MR-007 — Auflösung von MR-003: doc-check als Dogfooding

- **Datum:** 2026-06-10
- **Geltungsbereich:** `make doc-check`, [`.d-check.yml`](../.d-check.yml)
- **Adaption:** Der Auflösungs-Trigger von
  [`MR-003`](#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh)
  ist eingetreten: `make doc-check` läuft über `d-check` selbst
  (Runtime-Image, read-only-Mount; Module `links` + `anchors` über
  die gesamte Repo-Wurzel via `scan.roots: ["."]`). Das vendorte
  Skript `tools/verify-doc-refs.sh` ist gelöscht; der <!-- d-check:ignore (historisch: gelöscht) -->
  Geltungsbereich-Link in MR-003 wurde dafür auf einen Code-Span
  umgestellt (Form-, keine Inhaltsänderung). Vergleichslauf
  (erster Datenpunkt für
  [`DC-QA-04`](../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)):
  Alt-Skript 0 broken links, `d-check` 0 Befunde bei 23 Dateien —
  bei strikt größerer Abdeckung (zusätzlich Anker-Validierung und
  Bildreferenzen).
- **Begründung:** Dogfooding-Ziel von slice-004; die BF-Sub-Area aus
  der Modus-Tabelle ist damit graduiert (gelöscht).
- **Auflösungs-Trigger:** permanent (Dogfooding ist der Zielzustand).

### MR-008 — ID-Schema-Deklaration (Nachtrag zur Baseline-Aussage)

- **Datum:** 2026-06-11
- **Geltungsbereich:** gesamtes Repo (alle Artefakt-IDs und
  Traceability-Verweise)
- **Adaption:** Nachtrag der im Konventions-Template als Teil von
  MR-000 vorgesehenen ID-Schema-Deklaration, die in der
  Initial-Setzung fehlte (MR-000 ist akzeptiert und bleibt
  unverändert — Nachtrag daher als eigener Eintrag). Deklariertes
  Schema: Anforderungen `DC-FA-<BEREICH>-<NNN>` (Bereichskürzel, siehe
  [`MR-002`](#mr-002--id-schema-mit-bereichskürzeln-ab-initialer-fassung))
  und `DC-QA-<NN>`; ADRs `ADR-NNNN` (vierstellig, gemäß
  Kurs-ADR-Vorlage `NNNN-titel.template.md` — das Konventions-Template
  nennt abweichend dreistellig `ADR-<NNN>`, eine
  Kurs-Template-Inkonsistenz; Korrektur in der Kurs-Quelle steht aus,
  analog [`MR-006`](#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs));
  Konventions-Adaptionen `MR-NNN`; Carveouts `CO-NNN` (bisher
  ungenutzt); Slices `slice-NNN`. Die `ids`-Selbstkonfiguration in
  [`.d-check.yml`](../.d-check.yml) prüft die Linkpflicht der
  Kennungen maschinell.
- **Begründung:** Eine undeklarierte ID-Systematik ist eine stille
  Setzung (gleiche Harness-Lüge-Klasse wie ein undeklariertes Gate);
  sichtbar geworden im Template-Vergleich (User-Review, 2026-06-11).
- **Auflösungs-Trigger:** permanent. *(Die in der Adaption vermerkte
  Kurs-Template-Inkonsistenz — `conventions.template.md` dreistellig vs.
  ADR-Vorlage vierstellig — ist mit Baseline `v1.3.0` behoben (ADR-ID
  upstream vierstellig vereinheitlicht); siehe
  [`MR-012`](#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011).)*

### MR-009 — Source Precedence ohne `docs/user`-Rang

- **Datum:** 2026-06-11
- **Geltungsbereich:** [`harness/README.md` §Source precedence](README.md#source-precedence),
  [`AGENTS.md` §2](../AGENTS.md#2-kanonische-quellen-source-precedence)
- **Adaption:** Der Template-Default führt neun Ränge inkl. Rang 6
  `docs/user/*` (Operations, Quality, Releasing); d-check führt acht
  Ränge ohne `docs/user`, weil kein Operations-Doku-Stratum existiert
  (CLI-Tool vor dem ersten Release). Sichtbar geworden im
  Template-Vergleich (User-Review, 2026-06-11) — bis dahin eine
  stille Abweichung.
- **Begründung:** Ein Rang für nicht existierende Dateien wäre ein
  halluzinierter Eintrag (gleiche Klasse wie ein behauptetes Gate);
  die Rangordnung ist laut Baseline projektspezifische Wahl, die hier
  deklariert wird.
- **Auflösungs-Trigger:** welle-04 — mit der Release-Pipeline
  entsteht Betriebs-/Releasing-Doku; der `docs/user`-Rang wird dann
  eingefügt und dieser Eintrag als aufgelöst markiert. *(Eingetreten
  mit slice-011, siehe
  [`MR-010`](#mr-010--auflösung-von-mr-009-docsuser-rang-eingefügt).)*

### MR-010 — Auflösung von MR-009: `docs/user`-Rang eingefügt

- **Datum:** 2026-06-11
- **Geltungsbereich:** [`harness/README.md` §Source precedence](README.md#source-precedence),
  [`AGENTS.md` §2](../AGENTS.md#2-kanonische-quellen-source-precedence),
  `docs/user/`
- **Adaption:** Der Auflösungs-Trigger von
  [`MR-009`](#mr-009--source-precedence-ohne-docsuser-rang) ist
  eingetreten: mit der GHCR-Release-Pipeline (slice-011) existiert
  Betriebs-/Releasing-Doku (`docs/user/releasing.md`,
  `docs/user/operations.md`). Der `docs/user`-Rang ist als Rang 6 in
  beide Source-Precedence-Tabellen eingefügt (Template-Default
  wiederhergestellt, neun Ränge); die nachfolgenden Ränge rücken um
  eins.
- **Begründung:** Baseline-Konformität, sobald die Dateien real
  existieren — kein Rang für Phantome, kein Phantom für Ränge.
- **Auflösungs-Trigger:** permanent (Baseline-Konformität).

### MR-011 — Baseline auf Release-Tag gepinnt

- **Datum:** 2026-06-18
- **Geltungsbereich:** [§Baseline](#baseline), [§Adoptierte
  Konventions-Quellen](#adoptierte-konventions-quellen)
- **Adaption:** Die Baseline-Aussage und die externen
  Konventions-Quellen waren bislang **unversioniert**: der Stand als
  bloßer Datumsstempel „Template-Set 2026-06" und die Quellen-Pointer
  auf den `main`-Branch (`.../ai-harness-course/main/...`). Beides ist
  auf den Release-Tag **`v1.2.1`** gepinnt — die Stand-Zeile, die
  Raw-URL von `agents-regelwerk.md` und der Lehrmaterial-Repo-Link.
  MR-000 bleibt unverändert (akzeptiert; Nachtrag daher als eigener
  Eintrag, analog
  [`MR-008`](#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)).
- **Begründung:** Ein `main`-Pointer ist genau die Doku-Drift, die das
  Regelwerk für Außen-Verweise verbietet (tag-gepinnte URLs als
  Anti-Drift-Default); ein unversionierter Stand macht die
  Baseline-Konformitäts-Aussage von
  [`MR-000`](#mr-000--baseline-aussage) („keine inhaltlichen Adaptionen
  ggü. Baseline-Default") nicht mehr eindeutig prüfbar, sobald `main`
  weiterläuft. Anlass: User-Direktive 2026-06-18, das aktuelle
  Regelwerk (v1.2.1) zu befolgen.
- **Auflösungs-Trigger:** permanent als Mechanik; der konkrete Tag wird
  bei Adoption einer neueren Baseline-Version per Nachtrags-MR
  hochgezogen (dann auch Re-Evaluierung der an die Kurs-Quelle
  gekoppelten Trigger von
  [`MR-006`](#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)
  und [`MR-008`](#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)).

### MR-012 — Baseline-Pin-Hebung (Nachtrag zu MR-011)

- **Datum:** 2026-06-18
- **Geltungsbereich:** [§Baseline](#baseline), [§Adoptierte
  Konventions-Quellen](#adoptierte-konventions-quellen),
  [`MR-006`](#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs),
  [`MR-008`](#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)
- **Adaption:** Der von
  [`MR-011`](#mr-011--baseline-auf-release-tag-gepinnt) vorgesehene
  Nachtrag: der Baseline-Pin ist von `v1.2.1` auf den Release-Tag
  **`v1.3.0`** gehoben (Stand-Zeile, Lehrmaterial-Repo-Link,
  `agents-regelwerk.md`-Raw-URL sowie die gespiegelten Pointer in
  `AGENTS.md` und `harness/README.md`). Die beiden an die Kurs-Quelle
  gekoppelten Trigger sind dabei re-evaluiert, am Tag `v1.3.0` verifiziert:
  - [`MR-006`](#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)-Trigger
    **eingetreten:** die Spec-Straten-Vorlagen verweisen im bindenden Text
    nicht mehr abwärts auf ADRs (`spezifikation.template.md`: „kein
    ADR-Rückzeiger hier"; `architecture.template.md`: „kein ADR-Bezug in
    dieser Sicht", §ADR-Index entfernt). Die verbleibende ADR-Spalte nur
    in der Historie ist legitime Provenance (Regel 5), deckungsgleich mit
    dem `matrix`-`exclude-sections`-Default dieses Repos.
  - [`MR-008`](#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)-Kurs-Inkonsistenz
    **behoben:** die ADR-ID ist upstream durchgängig vierstellig
    (`ADR-NNNN`) — `conventions.template.md`, die `Bezug:`-Zeile des
    ADR-Templates und die `ADR-Bindung`-Klasse des Regelwerks tragen
    v1.3.0 vierstellig. d-checks eigene vierstellige Wahl ist damit
    baseline-konform statt vorgezogene Abweichung.
- **Begründung:** Die Tag-Hebung erfolgt per Nachtrags-MR (kein
  Überschreiben des akzeptierten Eintrags), wie in
  [`MR-011`](#mr-011--baseline-auf-release-tag-gepinnt) vorgesehen; die
  gekoppelten Trigger werden dabei re-evaluiert. v1.2.1 → v1.3.0 ist
  inhaltlich nur der ADR-Stelligkeits-Fix (eine Zeile im
  Regelwerk-Bundle) plus Tag-Rewrite — keine sonstige Regeländerung.
- **Auflösungs-Trigger:** permanent als Provenienz; die nächste
  Baseline-Version wird erneut per Nachtrags-MR gehoben.

### MR-013 — Lifecycle-Move-Commit bündelt gekoppelte Verweise

- **Datum:** 2026-06-21
- **Geltungsbereich:** [`AGENTS.md` §3.3](../AGENTS.md#33-git-mv--inhaltsänderung--zwei-commits),
  der Slice-Lifecycle `docs/plan/planning/in-progress/` → `…/done/`,
  `make planning-check`
- **Adaption:** Erste Adaption einer **Lifecycle-Regel** (Nachtrag zu
  [`MR-000`](#mr-000--baseline-aussage), das „keine Adaptionen … für
  Lifecycle-Regeln" vermerkt; MR-000 bleibt unverändert). Das Baseline-§3.3
  trennt reinen Move (Commit 1, `R100`-Rename) vom Inhalt (Commit 2). Seit
  `make planning-check` (slice-040) den Roadmap-Zustand **atomar** an den
  in-progress-Stand koppelt (kein `slice-*` in `…/in-progress/` ⟺ Roadmap
  §Aktuelle Welle trägt „Keine aktive Welle"), wäre ein *byte*-reiner
  Move-Commit beim Lifecycle-Move zwangsläufig gate-rot (leeres
  `in-progress/` bei noch aktiver Roadmap; zusätzlich `target-missing` auf
  jeden Verweis, der den Slice über seinen `in-progress/`-Pfad verlinkt).
  Adaption: der `git mv`-Commit dieses Moves trägt **zusätzlich** (a) den
  Roadmap-Flip §Aktuelle Welle und (b) alle Pfad-Verweise auf den Slice
  (Roadmap, `AGENTS.md` §4, `harness/README.md` §Sensors) von
  `in-progress/` nach `done/`. Der **Slice-Body** (Status-Zeile,
  Closure-Notiz) bleibt im zweiten Commit; weil die Slice-Datei im
  Move-Commit unverändert ist, hält die Rename-Detection (`R100`) und damit
  die `git log --follow`-Begründung des Baseline-§3.3.
- **Begründung:** Sichtbar 2026-06-21 — die PR-/Push-CI prüft den Push-Tip,
  der ein Zwischen-Commit sein kann; sie lief auf dem reinen Move-Commit von
  slice-040 rot (`target-missing` + `planning-check`). Die
  Per-Commit-Grün-Regel (grün = Boden, nicht Decke) und die
  Rename-Detection schließen sich nur scheinbar aus: die Kopplung betrifft
  **fremde** Dateien, nicht den Slice-Body. slice-040 führte
  `make planning-check` ein und löste damit die Kollision aus.
- **Auflösungs-Trigger:** permanent, solange `make planning-check` den
  Roadmap-↔-in-progress-Invariant erzwingt.

### MR-014 — Slice-/ADR-Doc-Struktur: Repo-Haus-Stil ggü. Baseline-Template

- **Datum:** 2026-06-22
- **Geltungsbereich:** Slice-Dateien unter `docs/plan/planning/` und
  ADR-Dateien unter `docs/plan/adr/` (Doc-**Struktur**, nicht Inhalt)
- **Adaption:** Slice- und ADR-Dateien folgen einem repo-lokalen Haus-Stil,
  der von den adoptierten Baseline-Vorlagen
  ([`slice.template.md`](https://github.com/pt9912/ai-harness-course/blob/v1.3.0/lab/templates/docs/plan/planning/slice.template.md),
  [`NNNN-titel.template.md`](https://github.com/pt9912/ai-harness-course/blob/v1.3.0/lab/templates/docs/plan/adr/NNNN-titel.template.md))
  in der **Abschnitts-Struktur** abweicht; Header-Block und Pflichtinhalte
  bleiben deckungsgleich:
  - **Slice:** keine separate „Plan (vor Code)"-Tabelle (der Plan steht im
    Entscheidungen-/Regel-Block, §2); kein eigener Abschnitt „Closure-Trigger"
    (in §5 Trigger bzw. der DoD aufgegangen); Reihenfolge §1 Ziel · §2
    Entscheidungen/Regel · §3 DoD · §4 Risiken · §5 Trigger · §6
    Sub-Area-Modus · §7 Closure-Notiz (das Template führt Sub-Area als §8).
  - **ADR:** „Verglichene Alternativen" als **Tabelle** (Pro/Contra) statt
    Option-A/B/C-Prosa; „Fitness Function" als Prosa-Bullets statt der
    Tooling/Regel/Target-Tabelle; „Geschichte" zweispaltig (Datum/Ereignis)
    statt dreispaltig (der Verweis steht im Ereignis-Text).
- **Begründung:** Die Vorlagen sind ausdrücklich Start-Vorlagen („Kopiere …
  ersetze Platzhalter, lösche den Hinweis-Block"); strukturelle Anpassung ist
  erwartbar. Der Haus-Stil ist über die bisherigen Slices und ADRs konsistent
  gelebt; Konsistenz mit dem Bestand sticht die buchstäbliche Template-Form.
  Sichtbar geworden im Template-Vergleich (Nutzer-Frage 2026-06-22,
  lab-templates v1.3.0) — bis dahin eine **stille Setzung** (gleiche Klasse wie
  eine undeklarierte Konvention). Inhaltlich bindende Elemente bleiben
  unberührt: die Header-Felder **Bezug**/**Schärft**, die Akzeptanz-/DoD-
  Disziplin, die Sub-Area-Modus-Pflicht und die `## Geschichte`-Anhang-Regel
  der ADR-Immutability ([ADR-0016](../docs/plan/adr/0016-adr-immutable-gate.md)).
  [`MR-000`](#mr-000--baseline-aussage) bleibt unverändert; Nachtrag daher als
  eigener Eintrag, analog
  [`MR-008`](#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage).
- **Auflösungs-Trigger:** permanent (Doc-Form ist Repo-Wahl; Baseline: „Form
  ist Wahl"). Bei künftiger Baseline-Hebung wird die Struktur-Differenz hier
  re-evaluiert (analog [`MR-011`](#mr-011--baseline-auf-release-tag-gepinnt)).

### MR-015 — Auflösung der MR-012-Pointer-Drift: AGENTS.md routet, spiegelt nicht mehr

- **Datum:** 2026-06-22
- **Geltungsbereich:** [`AGENTS.md`](../AGENTS.md) §1,
  [`MR-011`](#mr-011--baseline-auf-release-tag-gepinnt),
  [`MR-012`](#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011),
  §[Adoptierte Konventions-Quellen](#adoptierte-konventions-quellen)
- **Adaption:** Nachtrag zu
  [`MR-012`](#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011), das den
  Baseline-Pin u. a. „in den gespiegelten Pointern in AGENTS.md und
  harness/README.md" hob. Mit der zip-only-Umstellung von `AGENTS.md` §1
  (Commit `f46326a`, 2026-06-22 — Lese-Form ausschließlich `lab-regelwerk.zip`)
  wurde der `agents-regelwerk.md`-Link **aus AGENTS.md entfernt**: AGENTS.md
  **routet** dort für Quelldatei und Stand auf diese Datei
  (§[Adoptierte Konventions-Quellen](#adoptierte-konventions-quellen) /
  §[Baseline](#baseline)), statt die Raw-URL zu spiegeln. Die von
  [`MR-011`](#mr-011--baseline-auf-release-tag-gepinnt)/[`MR-012`](#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011)
  gepinnte `agents-regelwerk.md`-Raw-URL lebt damit nur noch in
  §[Adoptierte Konventions-Quellen](#adoptierte-konventions-quellen) und in
  [`harness/README.md`](README.md) §Guides. Beide bleiben als
  Vergangenheits-Aussage **unverändert** (eigener Eintrag, analog
  [`MR-008`](#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)).
- **Begründung:** Die Pointer-Liste in
  [`MR-012`](#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011) war am 2026-06-18
  korrekt; die spätere zip-only-Korrektur machte den AGENTS.md-Teil
  navigatorisch stale (wer „den gespiegelten AGENTS.md-Pointer" sucht, findet
  ihn nicht mehr). Eine undeklarierte Pointer-Verschiebung in einem
  MR-gepinnten Bereich ist dieselbe stille-Setzung-Klasse, die das Pinning
  verhindern soll; der Nachtrag stellt die Provenienz eindeutig wieder her. Die
  zip-Änderung selbst brauchte kein eigenes MR (reine Lese-Form-/Wortlaut-
  Korrektur, Nutzer-Entscheid 2026-06-22); ihre **Wirkung** auf die gepinnte
  Pointer-Liste wird hier nachgezogen.
- **Auflösungs-Trigger:** permanent (Provenienz).

### MR-016 — Baseline-Pin-Hebung (zweiter Nachtrag zu MR-011)

- **Datum:** 2026-06-25
- **Geltungsbereich:** [§Baseline](#baseline), [§Adoptierte
  Konventions-Quellen](#adoptierte-konventions-quellen), die gespiegelten
  `lab-regelwerk.zip`-Pointer in [`AGENTS.md`](../AGENTS.md) §1 und
  [`harness/README.md`](README.md) §Guides sowie die kurs-gekoppelten Trigger
  [`MR-006`](#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs),
  [`MR-008`](#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage),
  [`MR-014`](#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template)
- **Adaption:** Der von
  [`MR-011`](#mr-011--baseline-auf-release-tag-gepinnt) vorgesehene weitere
  Nachtrag: der Baseline-Pin ist von `v1.3.0` auf den Release-Tag **`v1.4.0`**
  gehoben. Aktualisiert: die §Baseline-Stand-Zeile, der Lehrmaterial-Repo-Link
  und die `agents-regelwerk.md`-Raw-URL + `lab-regelwerk.zip`-Asset-URL in
  §Adoptierte Konventions-Quellen, die gleichen beiden URLs in
  `harness/README.md` §Guides sowie — zip-only seit
  [`MR-015`](#mr-015--auflösung-der-mr-012-pointer-drift-agentsmd-routet-spiegelt-nicht-mehr) —
  der `lab-regelwerk.zip`-Pointer samt `(v1.4.0)`-Stand-Markierung in
  `AGENTS.md` §1. Das **Regelwerk ist inhaltlich unverändert**: `kurs/de/`
  (Quelle von `agents-regelwerk.md` und Bundle) trägt zwischen `v1.3.0` und
  `v1.4.0` keinen Diff; der Pin wandert, der Inhalt nicht. Der gesamte
  v1.3.0→v1.4.0-Diff liegt in `lab/templates/`. Die beiden an die Kurs-Quelle
  gekoppelten Trigger sind dabei re-evaluiert, am Tag `v1.4.0` verifiziert,
  plus eine teilweise Auflösung von [`MR-014`](#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template):
  - [`MR-006`](#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)-Trigger
    **unverändert eingetreten:** die Spec-Straten-Vorlagen
    (`spezifikation.template.md`, `architecture.template.md`) tragen keinen
    Diff zwischen `v1.3.0` und `v1.4.0`; der zu `v1.3.0` festgestellte Stand
    (kein ADR-Rückzeiger im bindenden Text) hält.
  - [`MR-008`](#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)-Konformität
    **gehalten:** das ADR-Template (`NNNN-titel.template.md`) hat sich geändert
    (Abschnitt „Verglichene Alternativen" von Option-A/B/C-Prosa auf eine
    Pro/Contra-**Tabelle**), die ADR-ID-Stelligkeit bleibt jedoch durchgängig
    vierstellig (`ADR-NNNN`).
  - [`MR-014`](#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template)-Abweichung
    **teilweise aufgelöst:** die dort als Haus-Stil deklarierte ADR-„Verglichene
    Alternativen"-**Tabelle** (statt Option-A/B/C-Prosa) ist mit `v1.4.0`
    **Baseline-Default** geworden — dieser eine Punkt ist damit baseline-konform
    statt Abweichung. Die übrigen MR-014-Punkte (Slice-Abschnitts-Reihenfolge,
    Fitness-Function als Prosa-Bullets, zweispaltige „Geschichte") stehen
    unverändert: `slice.template.md` und die restliche ADR-Template-Struktur
    tragen keinen Diff. Die `@v1.3.0`-Template-Links im MR-014-Body bleiben als
    Vergangenheits-Aussage **unangetastet** (append-only, analog
    [`MR-015`](#mr-015--auflösung-der-mr-012-pointer-drift-agentsmd-routet-spiegelt-nicht-mehr)).
- **Begründung:** Die Tag-Hebung erfolgt per Nachtrags-MR (kein Überschreiben
  des akzeptierten [`MR-012`](#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011)-Eintrags),
  wie in [`MR-011`](#mr-011--baseline-auf-release-tag-gepinnt) vorgesehen. v1.3.0
  → v1.4.0 ist inhaltlich nur Template-Pflege: die ADR-Alternativen-Tabelle
  (s. o.), `harness.mk` hebt den **d-check-Konsumenten**-Digest-Pin
  (v0.8.0 → v0.23.0 — betrifft den Kurs als d-check-Nutzer, nicht d-checks
  eigene Konventionen) und eine Templates-`README`-Notiz; keine Regelwerk- oder
  sonstige Konventionsänderung.
- **Auflösungs-Trigger:** permanent als Provenienz; die nächste
  Baseline-Version wird erneut per Nachtrags-MR gehoben.

### MR-017 — Lokale Baseline-Lese-Form (Cache) aus dem Selbst-Scan ausgenommen

- **Datum:** 2026-06-25
- **Geltungsbereich:** [`.d-check.yml`](../.d-check.yml) `scan.ignore`,
  das Materialisierungs-Skript `tools/harness/fetch-baseline-cache.sh`,
  [§Adoptierte Konventions-Quellen](#adoptierte-konventions-quellen),
  [`AGENTS.md`](../AGENTS.md) §1, der gitignorierte Pfad `.harness/cache/`
- **Adaption:** Die in [`AGENTS.md`](../AGENTS.md) §1 vorgesehene Lese-Form der
  adoptierten Baseline (Bundle „herunterladen, entpacken, nur den benötigten
  Abschnitt laden") wird lokal **materialisiert** nach dem Pfadschema
  `.harness/cache/<tag>/regelwerk/` (entpacktes `lab-regelwerk.zip`) und
  `.harness/cache/<tag>/templates/` (entpacktes `lab-templates.zip`); aktueller
  `<tag>` = `v1.4.0`; materialisiert wird reproduzierbar per
  `tools/harness/fetch-baseline-cache.sh` (zieht die beiden Release-Assets und
  entpackt; Tag ohne Argument aus dieser §Baseline-Stand-Zeile abgeleitet —
  kein Drift). Der Cache ist **gitignored** ([`.gitignore`](../.gitignore)
  `.harness/cache/`) — ephemer, kein Repo-Vertrag — und wird über
  `scan.ignore: [".harness/cache/**"]` in [`.d-check.yml`](../.d-check.yml) aus
  dem Dogfooding-Selbst-Scan **ausgenommen**. Grund: der Cache trägt
  Fremdinhalt (die Kurs-Docs referenzieren ihre eigenen `ADR-`/`MR-`-IDs und
  Modulpfade, die in *diesem* Repo nicht existieren); ohne Ausnahme meldete
  `make doc-check` sie als `id-unlinked`/`codepath-missing`. Selbe Klasse wie
  die eingebauten `SKIP_DIRS` (`.git`, `vendor`, `node_modules` —
  Fremd-/Generiertes), daher **keine** Gate-Lockerung im Sinne von
  [`AGENTS.md` §3.6](../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden):
  ausgenommen wird Nicht-Repo-Inhalt, keine Repo-Doku verliert Deckung.
- **Begründung:** Sichtbar geworden bei der v1.4.0-Adoption
  ([`MR-016`](#mr-016--baseline-pin-hebung-zweiter-nachtrag-zu-mr-011)), als der
  Cache erstmals befüllt wurde (Nutzer-Auftrag 2026-06-25). Die Konvention ist
  tag-generisch (`<tag>`), nicht v1.4.0-spezifisch — darum ein eigener Eintrag
  statt Bündelung in den Pin-Bump-Nachtrag.
- **Auflösungs-Trigger:** permanent, solange die Baseline-Lese-Form lokal
  gecacht wird.

### MR-018 — d-check verkörpert als Producer-/Self-Hoster keine Templates

- **Datum:** 2026-06-25
- **Geltungsbereich:** [`AGENTS.md`](../AGENTS.md) §1,
  [§Adoptierte Konventions-Quellen](#adoptierte-konventions-quellen), der
  `.harness/cache/<tag>/templates/`-Cache (Schärfung von
  [`MR-017`](#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)),
  die Autoren-Quelle wiederkehrender Artefakte
- **Adaption:** d-check verkörpert **keine** co-located `*.template.md` und
  weicht damit bewusst von der Baseline-Regel **§Ein- vs. wiederkehrende
  Templates** ab (`lab/templates/README.md`: die wiederkehrenden Skelette — ADR,
  Slice, Welle, Carveout, Review-Report — bleiben co-located, jede neue Instanz
  wird daneben kopiert; die Singletons werden beim Bootstrap einmal gefüllt).
  **d-check ist Producer-/Self-Hoster** der Harness-Werkzeuge: es autoriert seine
  wiederkehrenden Artefakte **nativ im Haus-Stil**
  ([`MR-014`](#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template))
  aus dem gelebten Bestand seiner ADRs und Slices, nicht aus einem co-located
  Baseline-Skelett. Der `lab-templates.zip`-Cache
  (`.harness/cache/<tag>/templates/`) ist daher **Adoptions-/Drift-Audit-Staging,
  nicht Autorenquelle** — niemand autoriert aus dem ephemeren, gitignorierten
  Cache.
- **Begründung:** Sichtbar im Schwester-Repo-Vergleich (Nutzer-Frage
  2026-06-25): das Consumer/Adopter-Repo **bedrock-eu-guard** verkörpert die fünf
  wiederkehrenden Skelette co-located (folgt §Ein- vs. wiederkehrende); d-check
  als Producer tut es nicht. Der **Kurs stützt die Producer-Lesart für
  `harness.mk` bereits live** (Stand > `v1.4.0`): seine `lab/templates/README.md`
  §Self-Hosting-/Producer-Fall nennt „das Tool-Repo selbst (d-check), das seinen
  Doku-Gate via `make doc-check` direkt dogfooded" und nimmt es von der
  `harness.mk`-Adoption aus. d-checks **gepinnte v1.4.0-Baseline trägt diesen
  Abschnitt noch nicht** (er ist post-v1.4.0); diese MR **überbrückt** bis zum
  nächsten Baseline-Bump und zieht die Producer-Logik zugleich auf die
  wiederkehrenden **Dokument**-Skelette weiter — den Schritt, den der Kurs auch
  live (noch) nicht ausspricht. Ohne die Deklaration bliebe d-checks
  Template-Freiheit eine **stille Setzung**, und ein Agent könnte fälschlich aus
  dem Audit-Cache autoren.
- **Auflösungs-Trigger:** Re-Evaluation beim nächsten Baseline-Bump (Pin auf den
  post-v1.4.0-Kurs, der den Self-Hosting-/Producer-Fall in der Templates-README
  bereits kanonisiert): die MR wird gegen den dann aktuellen Kanon geprüft und
  zur reinen Provenienz — oder aufgelöst, falls der Kurs die Producer-Lesart auch
  für die wiederkehrenden Dokument-Skelette übernimmt.

### MR-019 — Regelwerk-Lese-Form committet statt gecacht (Nachtrag zu MR-017)

- **Datum:** 2026-06-26
- **Geltungsbereich:** [§Adoptierte Konventions-Quellen](#adoptierte-konventions-quellen)
  (Lokale Lese-Form), [`.d-check.yml`](../.d-check.yml) `scan.ignore`, das Skript
  [`tools/harness/fetch-baseline-cache.sh`](../tools/harness/fetch-baseline-cache.sh),
  [`AGENTS.md`](../AGENTS.md) §1, der neu committete Pfad `.harness/baseline/<tag>/`;
  Nachtrag zu
  [`MR-017`](#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)
- **Adaption:** Die **Regelwerk**-Lese-Form der adoptierten Baseline wird nicht
  mehr nur im gitignorierten Cache materialisiert, sondern **committet vendored**
  nach `.harness/baseline/<tag>/regelwerk/` (entpacktes `lab-regelwerk.zip`),
  zusammen mit einem committeten `.harness/baseline/<tag>/SHA256SUMS` über die
  vendorten Dateien; aktueller `<tag>` = `v1.4.0`.
  `tools/harness/fetch-baseline-cache.sh` schreibt das Regelwerk in diesen
  Vendor-Pfad, (re)generiert SHA256SUMS und verifiziert; der `--verify`-Modus
  prüft das committete Regelwerk **netzlos** gegen das Manifest
  (CI/Audit/frischer Checkout). Der Vendor-Pfad ist über
  `scan.ignore: [".harness/baseline/**", …]` in [`.d-check.yml`](../.d-check.yml)
  aus dem Dogfooding-Selbst-Scan ausgenommen — **dieselbe** Begründung wie für den
  Cache in
  [`MR-017`](#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen):
  Fremdinhalt (die Kurs-Docs referenzieren eigene `ADR-`/`MR-`-IDs und Modulpfade,
  die in *diesem* Repo nicht existieren); **keine** Gate-Lockerung im Sinne von
  [`AGENTS.md` §3.6](../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden),
  da Nicht-Repo-Inhalt ausgenommen wird, keine Repo-Doku Deckung verliert. Das
  **Template**-Set bleibt unverändert ephemerer Cache unter
  `.harness/cache/<tag>/templates/` (nur Adoptions-/Drift-Audit-Staging,
  [`MR-018`](#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates)) —
  template-frei wird nur das Regelwerk vendored, nicht die Templates.
  [`MR-017`](#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)
  bleibt als Provenienz stehen; seine Cache-Aussage gilt fortan nur noch für die
  Templates.
- **Begründung:** Sichtbar geworden über die Steering-Loop-Verweis-Frage
  (2026-06-26): ein Pointer auf die kanonische Regelwerk-Definition (z. B.
  `grundlagen-klassifikation.md` §Steering Loop) löste in keiner Zielumgebung
  in-repo auf, weil die Lese-Form ephemer/gitignored war — auf frischem Checkout
  oder ohne Netz schlicht abwesend. Schärfer noch: d-check pinnt seine eigenen
  Release-Digests (`sha256:…`) und erkennt mit dem `pins`-Modul inhaltlichen Drift
  verlinkter Quellen, konsumierte aber seine **eigene** Baseline-Lese-Form per
  `curl` von einem Release-Asset **ohne Content-Hash** — der Pin hielt den Tag,
  nicht die Bytes. Das Vendoring schließt drei Lücken in einem: Präsenz (jeder
  Checkout hat das Regelwerk), Offline-Auditierbarkeit und Integrität/Provenienz
  (Bytes via SHA256SUMS + git-Historie statt unverifiziertem Netz-Fetch). Es
  bestätigt das Verkörperungs-Prinzip (Kurs Modul 0: Per-Lauf-/Schwellen-
  Relevantes gehört verkörpert, nicht extern nachgeladen) und stützt
  [`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  (netzlose Lese-Form). Re-Vendor + neues Manifest sind ein bewusster Akt am
  Baseline-Pin-Bump, kein laufender Drift.
- **Auflösungs-Trigger:** permanent, solange die Baseline-Regelwerk-Lese-Form
  lokal vendored wird; der nächste Baseline-Pin-Bump re-vendored und erneuert
  SHA256SUMS.

### MR-020 — Baseline-Template-Propagation per Drift-Audit (template-frei bestätigt)

- **Datum:** 2026-06-26
- **Geltungsbereich:** [`docs/plan/planning/README.md`](../docs/plan/planning/README.md)
  (§Lifecycle, Closure-Notiz §7), der Templates-Staging-Cache
  `.harness/cache/<tag>/templates/`, die Baseline-Pin-Bump-Prozedur; Nachtrag zu
  [`MR-018`](#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates)
- **Adaption:** d-check bleibt **template-frei**
  ([`MR-014`](#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template),
  [`MR-018`](#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates))
  — bestätigt, nicht aufgehoben. Damit Baseline-Template-Verbesserungen den
  gelebten Haus-Stil dennoch erreichen, wird der dort angelegte
  Templates-Staging-Cache (`.harness/cache/<tag>/templates/`, ephemer — nur diese
  Rolle bleibt dem Cache nach
  [`MR-019`](#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017))
  beim **Baseline-Pin-Bump** als **Drift-Audit** ausgeführt: die staged
  Baseline-Skelette werden gegen den Haus-Stil verglichen; substanzielle
  Verbesserungen (neue Pflichtfelder, Pointer, …) werden eingezogen, rein
  stilistische Template-Eigenheiten bleiben außen vor. Erster eingezogener Fall:
  der §7-Closure-Notiz-**Steering-Loop-Eintrag** verweist auf die kanonische
  Definition im **vendorten** Regelwerk
  (`.harness/baseline/<tag>/regelwerk/grundlagen-klassifikation.md` §Steering
  Loop, in-repo auflösbar seit
  [`MR-019`](#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)),
  verankert in [`docs/plan/planning/README.md`](../docs/plan/planning/README.md).
- **Begründung:** Nutzer-Frage 2026-06-26: „Eigene Wege gehen heißt, Baseline-
  Probleme nochmals lösen." Stimmt — der Fork tauscht *automatische Vererbung*
  gegen *gelebten Haus-Stil*; bei Slice-50 sind die fertigen Slices die bessere
  Vorlage. Der Defekt war nicht der Fork, sondern dass das in
  [`MR-018`](#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates)
  angelegte Staging nie *ausgeführt* wurde. Auslöser: der Kurs ergänzte
  `slice.template.md` §7 um einen Steering-Loop-Pointer, den d-check ohne Audit
  nicht mitbekam. Bewusst **kein Gate** — Template-Drift ist eine Erkenntnis-,
  keine Laufzeit-Eigenschaft (analog zur Regelwerk-Lesedisziplin); die Disziplin
  liegt am Pin-Bump.
- **Auflösungs-Trigger:** permanent, solange d-check template-frei ist; jeder
  Baseline-Pin-Bump führt den Drift-Audit aus.

### MR-021 — In-Repo-Verweise auf das vendored Regelwerk sind pin-gebunden

- **Datum:** 2026-06-26
- **Geltungsbereich:** alle Markdown-Links auf `.harness/baseline/<tag>/…` in der
  Live-Doku (aktuell [`harness/README.md`](README.md) §Guides/§Sensors,
  [`docs/plan/planning/README.md`](../docs/plan/planning/README.md),
  [§Adoptierte Konventions-Quellen](#adoptierte-konventions-quellen)); die
  Baseline-Pin-Bump-Prozedur; Nachtrag zu
  [`MR-019`](#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)/[`MR-020`](#mr-020--baseline-template-propagation-per-drift-audit-template-frei-bestätigt)
- **Adaption:** Seit
  [`MR-019`](#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)
  ist das vendored Regelwerk ein in-repo auflösbares **Link-Ziel** — die
  Live-Doku verweist auf konkrete Regelwerk-Dateien (Lesestoff: Modul-/
  Grundlagen-Verweise) statt nur auf externe Kurs-URLs (die als **Provenienz**
  bleiben). Diese Links tragen den **konkreten** Pin (`…/v1.4.0/…`), nicht
  `<tag>` — sie sind damit **pin-gebunden**. Regel: Der Baseline-Pin-Bump-
  Drift-Audit ([`MR-020`](#mr-020--baseline-template-propagation-per-drift-audit-template-frei-bestätigt))
  (1) entfernt das alte `.harness/baseline/<alt-tag>/` und (2) zieht alle
  vendored-Pfad-Links auf den neuen Tag. Wird (2) vergessen, schlägt nach (1)
  `make doc-check` mit `target-missing` an — die Pin-Kopplung ist damit
  **gate-erzwungen**, kein stiller Drift. (Bliebe das alte Tag-Verzeichnis
  stehen, läge stiller Stale-Content vor; darum ist (1) Pflicht-Teil des Bumps.)
- **Begründung:** Nutzer-Entscheid 2026-06-26, das vendored Regelwerk als
  Lesestoff zu verlinken (§Guides-Lese-Form + Modul-13/14-Verweise +
  §Adoptierte-Aktuell-Link). Der Nutzen (klickbar, netzlos, offline auffindbar)
  hat als Preis die Pin-Bindung; die Regel macht den Preis explizit und delegiert
  die Durchsetzung an das vorhandene `links`-Gate, statt einen neuen Sensor zu
  bauen (Steering-Loop-Ökonomie: kein Gate für etwas, das ein vorhandenes Gate
  schon fängt).
- **Auflösungs-Trigger:** permanent, solange in-repo-Verweise auf das vendored
  Regelwerk bestehen.

### MR-022 — Baseline-Currency-Audit-Modus (Nachtrag zu MR-019)

- **Datum:** 2026-07-19
- **Geltungsbereich:** [`tools/harness/fetch-baseline-cache.sh`](../tools/harness/fetch-baseline-cache.sh)
  (neuer Modus `--check-latest`), [`AGENTS.md`](../AGENTS.md) §1; Nachtrag zu
  [`MR-019`](#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)
- **Adaption:** `fetch-baseline-cache.sh` erhält einen dritten Modus
  `--check-latest` mit **zwei** Upstream-Prüfungen (beide Netz, informativ):
  **(A) Currency** — der in [§Baseline](#baseline) gepinnte Tag gegen das
  **neueste stabile** Release
  (`https://api.github.com/repos/pt9912/ai-harness-course/releases/latest`;
  GitHub blendet dort Prereleases/Drafts aus, passend zur Re-Adopt-Semantik).
  **(B) Content-Drift am gepinnten Tag** — das Skript lädt `lab-regelwerk.zip`
  des **gepinnten** Tags und vergleicht dessen Bytes (dasselbe
  `sha256sum regelwerk/*.md`-Manifest) gegen das committete
  `.harness/baseline/<tag>/SHA256SUMS`; eine Abweichung heißt, der Tag wurde
  **verschoben** oder das Asset neu hochgeladen. Ausgang (schlimmster Fall):
  aktuell & authentisch → `exit 0`; neuerer Tag (`sort -V`) → `exit 3`
  (**Signal**, kein Fehler); Content-Drift am gepinnten Tag → `exit 4`
  (**Provenienz-Alarm**); nicht erreichbare Teile (Netz/API/Rate-Limit,
  fehlendes Werkzeug/Manifest) → **SKIP** je Teil (`exit 0`, sofern der andere
  Teil nicht `3`/`4` meldet). Der Modus ist bewusst **kein Gate**
  (`--network`-abhängig wie das Re-Vendoring) und bewusst **nicht fail-closed**
  (Gegenstück zu `--verify`, das netzlose Integrität prüft und sehr wohl
  fail-closed ist): ein nicht erreichbares Upstream darf keinen Lauf blockieren.
- **Begründung:** Übernommen aus dem Kurs-Beispiel
  `lab/example/tools/check_regelwerk_drift.py` (inhaltsbasierter Drift-Sensor der
  adoptierten Form-Quelle), auf d-checks Tag-Pin-Modell übersetzt. d-check pinnt
  einen Release-Tag und vendored ihn hash-verifiziert (`--verify`,
  [`MR-019`](#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017));
  `--verify` prüft aber nur die **vendorten** Bytes gegen ihr **eigenes**
  Manifest — nie gegen Upstream. Es blieben zwei tote Winkel, die
  `--check-latest` schließt: **Currency** (Teil A — liegt der Pin hinter
  Upstream? Ein grüner Integritäts-Check verdeckte das) und **Authentizität des
  gepinnten Tags** (Teil B — d-check *setzt* auf Tag-Immutabilität, aber »Tag
  verschoben / Asset neu« bemerkt `--verify` nicht). Teil B **verifiziert** genau
  die Immutabilitäts-Annahme aus
  [`MR-011`](#mr-011--baseline-auf-release-tag-gepinnt) (prüfen statt vertrauen)
  und ist der Content-Hash-Drift-Kern des Kurs-Beispiels, auf den gepinnten Tag
  angewandt. Beides **ohne** den Grundsatz aus
  [`MR-019`](#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)
  zu unterlaufen (»Re-Vendor + neues Manifest sind ein bewusster Akt am
  Baseline-Pin-Bump, kein laufender Drift«): der Modus **automatisiert nichts**,
  er meldet nur — der Re-Adopt bleibt der bewusste manuelle Akt
  ([`MR-020`](#mr-020--baseline-template-propagation-per-drift-audit-template-frei-bestätigt) /
  [`MR-021`](#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden)).
- **Auflösungs-Trigger:** permanent, solange die Baseline extern gepinnt und
  lokal vendored wird.

### MR-023 — Baseline-Pin-Hebung auf v5.0.0 samt self-contained Bundle-Layout

- **Datum:** 2026-08-01
- **Geltungsbereich:** [§Baseline](#baseline), [§Adoptierte
  Konventions-Quellen](#adoptierte-konventions-quellen), das
  Materialisierungs-Skript
  [`tools/harness/fetch-baseline-cache.sh`](../tools/harness/fetch-baseline-cache.sh),
  die pin-gebundenen Verweise
  ([`MR-021`](#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
  in [`AGENTS.md`](../AGENTS.md), [`harness/README.md`](README.md),
  [`.harness/skills/reviewer.md`](../.harness/skills/reviewer.md) und den
  Planning-Docs; Nachträge zu
  [`MR-017`](#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen),
  [`MR-018`](#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates),
  [`MR-019`](#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017),
  [`MR-022`](#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019)
- **Adaption:** Der Baseline-Pin ist von `v1.4.0` auf **`v5.0.0`** gehoben — der
  von [`MR-011`](#mr-011--baseline-auf-release-tag-gepinnt) vorgesehene Nachtrag,
  hier über **zwei weitere Major-Sprünge** (`v4.0.0`, `v5.0.0`). Anders als
  [`MR-016`](#mr-016--baseline-pin-hebung-zweiter-nachtrag-zu-mr-011) (reine
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
    [`MR-017`](#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)/[`MR-018`](#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates)
    auf (Templates nicht länger ephemerer `.harness/cache/`-Stand); die
    **vollständige Neufassung** der MR-017/018/019-Einträge auf diese Form ist
    Teil der späteren Konventionsspeicher-Etappe.
  - **Skript aufs Bundle gehoben.** `fetch-baseline-cache.sh` entpackt das Bundle
    **tolerant** (Regelwerk am `modul-00`-Marker, Templates als Geschwister),
    prüft eine **Under-Copy-Barriere** (Quelle == vendored) und schreibt das
    `SHA256SUMS`-Manifest über den **tatsächlichen** Bestand **beider** Bäume.
    Der `--check-latest`-Currency-Teil liest jetzt die Release-**Liste** (statt
    `/releases/latest`), der Content-Drift-Teil vergleicht **beide** Bäume; die
    Wortlaut-Angleichung der
    [`MR-019`](#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)/[`MR-022`](#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019)-Prosa
    folgt in der Konventionsspeicher-Etappe.
  - **Entfallene Quellzeiger umgeschrieben.** `agents-regelwerk.md` (im Kurs
    zurückgezogen) und die Kurs-`grundlagen/konventionen.md` (in acht
    `grundlagen-*`-Dateien aufgeteilt) werden **nicht** auf tote Ziele
    retargetet, sondern umgeschrieben: die committet vendorte `regelwerk/` (mit
    `README.md` als Index) **ist** die Agenten-Lese-Form; die §Referenz-Richtung
    liegt jetzt in `grundlagen-referenz-richtung.md`.
  - **Eingefrorene Historie via Tombstone.** Drei immutable/`done/`-Verweise auf
    den entfernten `v1.4.0`-Pfad
    ([`ADR-0022`](../docs/plan/adr/0022-matrix-token-richtung-provenance-marker.md),
    `slice-080`, `slice-081`) werden über das geteilte Referenz-Ventil
    `ignore-refs`
    ([`DC-FA-REF-001`](../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus))
    referenz-weit von der Existenz-Prüfung ausgenommen — kein Editieren
    eingefrorener Doku.
- **Begründung:** Auftraggeber-Vorgabe (2026-07-25): vollständige Migration nach
  `v5.0.0`, der Baseline-Default sticht die repo-lokale Adaption. Etappen-Schnitt
  und Umgang mit den historischen Verweisen sind in der abgenommenen Analyse
  (`slice-083`) begründet. Die Konventionsspeicher-Restrukturierung (Index +
  Datei je MR mit neuen Pflichtfeldern) und die Form-Konformität folgen als
  eigene Etappen.
- **Auflösungs-Trigger:** die Konventionsspeicher-Migration (spätere Etappe)
  fasst diese MR in die Datei-je-MR-Form und gleicht die
  MR-017/018/019/022-Prosa vollständig auf das v5.0.0-Layout an.

## Anforderungs-Anlege-Prozess

Neue oder geänderte `DC-*`-Anforderungen entstehen **nur** in
[`spec/lastenheft.md`](../spec/lastenheft.md) (vertraglich,
Change-Request-Charakter — Baseline-Regel der Spec-Stratifizierung;
Rang-Struktur dieses Repos: [`MR-001`](#mr-001--source-precedence-mit-eigener-spezifikations-schicht)).
Pflicht-Bausteine pro Anforderung:

- **ID gemäß Schema-Konvention** im Lastenheft §3
  (`DC-FA-<BEREICH>-<NNN>`, siehe
  [`MR-002`](#mr-002--id-schema-mit-bereichskürzeln-ab-initialer-fassung));
  ein neues Bereichskürzel wird dort in der Schema-Konvention
  deklariert. Nichtfunktionale Anforderungen: `DC-QA-<NN>`.
- **Drei Akzeptanzkriterien** (Happy/Boundary/Negative im
  Given/When/Then-Stil) plus explizite **Out-of-Scope**-Liste.
- **Versions-Bump + Historie-Zeile** im Lastenheft (§7).
- **Schärfungs-Richtung:** ADRs dürfen die Spezifikation schärfen,
  nie das Lastenheft (siehe `MR-001`-Begründung); wer das Lastenheft
  ändern will, ändert es direkt — als Change Request, nicht per ADR.
- **Beleg-Pflicht:** Test, Gate, Demo oder ADR folgt mit dem
  umsetzenden Slice
  ([`harness/README.md` §Traceability rules](README.md#traceability-rules)).

## Zusatzklassen-Deklaration für Sensors-Bindung

Zusätzlich zu den vier kanonischen Klassen (ADR, Carveout, Schwelle,
Reproduzierbarkeit):

| Klasse     | Form   | Bedeutung                                       | Beispiel                                                                                                                                                                                                                                                          |
| ---------- | ------ | ----------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| DC-Bindung | `DC-…` | Gate prüft eine konkrete Lastenheft-Anforderung | [`DC-QA-02`](../spec/lastenheft.md#dc-qa-02--determinismus) für den Determinismus-Test in `make test`; [`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) für `make arch-check` und das Netzlos-Gate in `make doc-check` |

## Modus-Deklaration pro Sub-Area

| Sub-Area (Pfad / Modul)         | Modus      | Begründung                                                                                                                          | Graduation-Bedingung / Folge-Slice |
| ------------------------------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------- |
| `*` (Default für gesamtes Repo) | Greenfield | Projekt startet spec-first; Doc führt, Code folgt                                                                                   | n/a (GF)                           |
| `tools/harness/`                | Greenfield | adoptierte Harness-Mechanik, konventionsgetragen über [`MR-004`](#mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild) | n/a (GF)                           |
