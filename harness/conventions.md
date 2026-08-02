# Harness-Konventionen

## Purpose

Diese Datei deklariert die *repo-lokalen* Strukturregeln dieses Repos
gegenüber der adoptierten Harnesskonvention (Baseline). Sie ist der
Default-Ort für:

- **Adaptionen** ggü. der Baseline (mit Begründung und Auflösungs-Trigger).
- **ID-Schema-Deklaration** — welches Präfix-Schema dieses Repo nutzt.
  Der Baseline-Default wird als Teil der [`MR-000`](#mr-000--baseline-aussage)-Aussage
  festgehalten; ein abweichendes Präfix oder Schema ist ein eigener `MR`-Eintrag.
- **Zusatzklassen-Deklarationen** für repo-spezifische
  Bindung-Klassen in der Sensors-Tabelle, die über die vier kanonischen
  hinausgehen (ADR, Carveout, Schwelle, Reproduzierbarkeit).
- **Modus-Deklarationen** pro Sub-Area (Greenfield / Brownfield /
  Hybrid) inklusive Konvergenz-Auftrag bei BF.

Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt die
kanonische Quelle (Source Precedence, [`harness/README.md`](README.md#source-precedence)).
Diese Datei ist konformitätsbringend für *Form*-Fragen, nicht autoritativ
über Inhalt.

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

**Adoptions-Erklärung ([MR-000](#mr-000--baseline-aussage)):** keine inhaltlichen Adaptionen gegenüber dem
Baseline-Default für Verzeichniskonvention, Lifecycle-Regeln und
Carveout-Disziplin; jede spätere Abweichung wird als `MR-<NNN>` nachgetragen.

Jede Adaption liegt als **eigene Datei** unter `harness/conventions/` (aktiv)
bzw. `harness/conventions/done/` (aufgelöst — der **Zustand ist die
Verzeichnis-Position**, kein Status-Feld). Die Voll-Slug-Anker aller Einträge
(auch der aufgelösten) hält der **Anker-Kompatibilitäts-Block** unten, damit die
in-repo-Verweise auf `conventions.md#mr-…` ohne Retarget auflösen.

| Adaption | Zustand |
|---|---|
| [MR-001 — Source Precedence mit eigener Spezifikations-Schicht](conventions/done/MR-001-eigene-spec-schicht.md) | aufgelöst |
| [MR-002 — ID-Schema mit Bereichskürzeln ab initialer Fassung](conventions/done/MR-002-id-schema-bereichskuerzel.md) | aufgelöst |
| [MR-003 — Vendorter Bootstrap-Sensor tools/verify-doc-refs.sh](conventions/done/MR-003-vendorter-bootstrap-sensor.md) | aufgelöst |
| [MR-004 — Gate-Nachweis-Mechanik und .claude-Hooks nach b-cad-Vorbild](conventions/MR-004-gate-nachweis-mechanik.md) | aktiv |
| [MR-005 — Härtung ggü. b-cad: inhaltsbasierter Gate-Nachweis, Sub-Shell-Prüfung](conventions/MR-005-haertung-gate-nachweis.md) | aktiv |
| [MR-006 — Referenzrichtung: Spec-Straten verweisen nie abwärts auf ADRs](conventions/MR-006-referenzrichtung-matrix.md) | aktiv |
| [MR-007 — Auflösung von MR-003: doc-check als Dogfooding](conventions/done/MR-007-aufloesung-mr-003.md) | aufgelöst |
| [MR-008 — ID-Schema-Deklaration (Nachtrag zur Baseline-Aussage)](conventions/done/MR-008-id-schema-deklaration.md) | aufgelöst |
| [MR-009 — Source Precedence ohne docs/user-Rang](conventions/done/MR-009-source-precedence-ohne-docs-user.md) | aufgelöst |
| [MR-010 — Auflösung von MR-009: docs/user-Rang eingefügt](conventions/done/MR-010-aufloesung-mr-009.md) | aufgelöst |
| [MR-011 — Baseline auf Release-Tag gepinnt](conventions/done/MR-011-baseline-pin-release-tag.md) | aufgelöst |
| [MR-012 — Baseline-Pin-Hebung (Nachtrag zu MR-011)](conventions/done/MR-012-baseline-pin-hebung.md) | aufgelöst |
| [MR-013 — Lifecycle-Move-Commit bündelt gekoppelte Verweise](conventions/MR-013-lifecycle-move-buendelung.md) | aktiv |
| [MR-014 — Slice-/ADR-Doc-Struktur: Repo-Haus-Stil ggü. Baseline-Template](conventions/done/MR-014-slice-adr-haus-stil.md) | aufgelöst |
| [MR-015 — Auflösung der MR-012-Pointer-Drift: AGENTS.md routet, spiegelt nicht mehr](conventions/MR-015-agents-md-routet.md) | aktiv |
| [MR-016 — Baseline-Pin-Hebung (zweiter Nachtrag zu MR-011)](conventions/done/MR-016-baseline-pin-hebung-2.md) | aufgelöst |
| [MR-017 — Lokale Baseline-Lese-Form (Cache) aus dem Selbst-Scan ausgenommen](conventions/done/MR-017-cache-selbst-scan.md) | aufgelöst |
| [MR-018 — d-check verkörpert als Producer-/Self-Hoster keine Templates](conventions/done/MR-018-keine-templates.md) | aufgelöst |
| [MR-019 — Regelwerk-Lese-Form committet statt gecacht (Nachtrag zu MR-017)](conventions/done/MR-019-regelwerk-vendored.md) | aufgelöst |
| [MR-020 — Baseline-Template-Propagation per Drift-Audit (template-frei bestätigt)](conventions/done/MR-020-template-propagation.md) | aufgelöst |
| [MR-021 — In-Repo-Verweise auf das vendored Regelwerk sind pin-gebunden](conventions/MR-021-vendored-verweise-pin-gebunden.md) | aktiv |
| [MR-022 — Baseline-Currency-Audit-Modus (Nachtrag zu MR-019)](conventions/done/MR-022-currency-audit.md) | aufgelöst |
| [MR-023 — Baseline-Pin-Hebung auf v5.0.0 samt self-contained Bundle-Layout](conventions/MR-023-baseline-v500.md) | aktiv |

## Anker-Kompatibilität (Baseline-Migration v5.0.0)

<!--
Migrationsspezifischer Kompatibilitäts-Block (welle-67, Etappe C): Mit dem Umzug
der Adaptionen aus dem Inline-Adaptions-Block in Einzeldateien verlassen die
MR-Heading-Anker diese Datei. 188 conventions.md-Voll-Slug-Links in 57 Dateien
(darunter 12 immutable ADRs und done-Slices) wuerden brechen. Dieser Block haelt
je referenziertem MR — auch den aufgeloesten — einen Voll-Slug-Anker, sodass alle
Links ohne Retarget und ohne ADR-Edit aufloesen. Ein frisches v5.0.0-Repo braucht
ihn nicht.
-->

<a id="mr-000--baseline-aussage"></a>
<a id="mr-001--source-precedence-mit-eigener-spezifikations-schicht"></a>
<a id="mr-002--id-schema-mit-bereichskürzeln-ab-initialer-fassung"></a>
<a id="mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh"></a>
<a id="mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild"></a>
<a id="mr-005--härtung-ggü-b-cad-inhaltsbasierter-gate-nachweis-sub-shell-prüfung"></a>
<a id="mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs"></a>
<a id="mr-007--auflösung-von-mr-003-doc-check-als-dogfooding"></a>
<a id="mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage"></a>
<a id="mr-009--source-precedence-ohne-docsuser-rang"></a>
<a id="mr-010--auflösung-von-mr-009-docsuser-rang-eingefügt"></a>
<a id="mr-011--baseline-auf-release-tag-gepinnt"></a>
<a id="mr-012--baseline-pin-hebung-nachtrag-zu-mr-011"></a>
<a id="mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise"></a>
<a id="mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template"></a>
<a id="mr-015--auflösung-der-mr-012-pointer-drift-agentsmd-routet-spiegelt-nicht-mehr"></a>
<a id="mr-016--baseline-pin-hebung-zweiter-nachtrag-zu-mr-011"></a>
<a id="mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen"></a>
<a id="mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates"></a>
<a id="mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017"></a>
<a id="mr-020--baseline-template-propagation-per-drift-audit-template-frei-bestätigt"></a>
<a id="mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden"></a>
<a id="mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019"></a>
<a id="mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout"></a>

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
  nie das Lastenheft (siehe [`MR-001`](#mr-001--source-precedence-mit-eigener-spezifikations-schicht)-Begründung); wer das Lastenheft
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
