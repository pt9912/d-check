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
  (`kurs/de/` — Konventionen in `grundlagen/`, Templates in `lab/templates/`).
  Kanonische Quelle; bei Konflikt maßgeblich.
- **Vendored Baseline (Regelwerk + Templates):** aus dem self-contained
  Release-Asset
  [`lab-regelwerk.zip`](https://github.com/pt9912/ai-harness-course/releases/download/v5.0.0/lab-regelwerk.zip)
  entpackt nach [`.harness/baseline/v5.0.0/`](../.harness/baseline/v5.0.0/regelwerk/)
  (`{regelwerk,templates}/` + `SHA256SUMS`) — der **netzlose** Lesepfad,
  materialisiert/verifiziert per
  [`fetch-baseline-cache.sh`](../tools/harness/fetch-baseline-cache.sh).
- **In-Repo (verkörperte Form):** die eigenen Artefakte dieses Repos (`AGENTS.md`,
  `harness/README.md`, ADRs, Slices, dieser Konventionsspeicher) — kopiert und
  ausgefüllt aus den vendored `templates/` (Referenz-/Ziel-Form).

## Adaptions-Block

Regeln dieser Sektion: Diese Datei trägt den **Index**, nicht die Einträge. Jede
Adaption ist eine eigene Datei unter `harness/conventions/`, kopiert aus der
vendored Vorlage `harness/conventions/MR-<NNN>-titel.template.md`; ist ihr
Auflösungs-Trigger eingetreten, wandert sie per `git mv` nach `conventions/done/`.
Der **Zustand ist die Verzeichnis-Position**, kein Status-Feld — was hier steht,
liest **jeder** Agentenlauf, aufgelöste Adaptionen gehören nicht in diesen Pfad
([Baseline-Regelwerk §Konventionsspeicher](../.harness/baseline/v5.0.0/regelwerk/grundlagen-harness-dateien.md#harnessconventionsmd-als-konventionsspeicher)).
Die Voll-Slug-Anker aller Einträge (auch der aufgelösten) hält der
[Anker-Kompatibilitäts-Block](#anker-kompatibilität-baseline-migration-v500) unten
— **Migrations-Schuld**, damit die in-repo-Verweise auf `conventions.md#mr-…` ohne
Retarget auflösen.

### MR-000 — Baseline-Aussage

Bleibt hier: keine Adaption, sondern die Adoptions-Erklärung — sie gilt für jeden Lauf.

- **Status:** Accepted
- **Datum:** 2026-06-10
- **Geltungsbereich:** gesamtes Repo
- **Ersetzt-Baseline-Regel:** — *(keine; dieser Eintrag ist die Adoptions-Erklärung,
  keine Adaption)*
- **Adaption:** *keine inhaltlichen Adaptionen ggü. Baseline-Default für
  Verzeichniskonvention, Lifecycle-Regeln, Carveout-Disziplin, ID-Schema
  (`DC-FA-*`, `DC-QA-*`, `ADR-NNNN`, `CO-NNN`, `slice-NNN`, `MR-NNN`; Präfix `DC`).*
- **Begründung:** Initial-Setzung. Spätere Adaptionen werden als `MR-<NNN>` nachgetragen.
- **Auflösungs-Trigger:** permanent.

### Aktive Adaptionen

Eine Zeile je Datei in `harness/conventions/`; Geltungsbereich und
Ersetzt-Baseline-Regel stehen hier, damit ein Agent ohne Öffnen entscheiden kann,
ob der Eintrag ihn betrifft.

| MR | Titel | Geltungsbereich | Ersetzt-Baseline-Regel |
|---|---|---|---|
| [MR-004](conventions/MR-004-gate-nachweis-mechanik.md) | Gate-Nachweis-Mechanik + `.claude`-Hooks | `tools/harness/`, `.claude/`, `make record-gates` | `grundlagen-durchsetzungsschicht` §Drei Bindepunkte |
| [MR-005](conventions/MR-005-haertung-gate-nachweis.md) | Härtung: Content-Hash + Sub-Shell-Guard | `working-tree-hash.sh`, `.claude/hooks/` | `grundlagen-durchsetzungsschicht` §Vier Design-Eigenschaften |
| [MR-006](conventions/MR-006-referenzrichtung-matrix.md) | Referenzrichtung/Matrix (+ C-4-Scope-Grenze) | `spec/`-Straten, `matrix`-Config | `grundlagen-referenz-richtung` §SDP |
| [MR-007](conventions/MR-007-aufloesung-mr-003.md) | doc-check als Dogfooding | `make doc-check`, `.d-check.yml` | `modul-13` §Hard Rule (Doku-Disziplin) |
| [MR-013](conventions/MR-013-lifecycle-move-buendelung.md) | Lifecycle-Move-Commit bündelt Verweise | Slice-Lifecycle, `make planning-check` | `modul-05` §Lifecycle als State Machine |
| [MR-015](conventions/MR-015-agents-md-routet.md) | AGENTS.md routet (spiegelt nicht) | `AGENTS.md` §1 | `grundlagen-harness-dateien` §Template-Schichtung |
| [MR-021](conventions/MR-021-vendored-verweise-pin-gebunden.md) | vendored-Verweise pin-gebunden | Live-Links auf die vendored Baseline | `grundlagen-harness-dateien` §Verzeichniskonvention |
| [MR-023](conventions/MR-023-baseline-v500.md) | Baseline-Pin `v5.0.0` + self-contained Bundle | §Baseline, `fetch-baseline-cache.sh` | `grundlagen-harness-dateien` §Template-Schichtung |

### Aufgelöste Adaptionen

Eine Zeile je Datei in `conventions/done/` — nur ID und Nachfolger, damit die Kette
auffindbar bleibt, ohne gelesen zu werden.

| MR | aufgelöst durch |
|---|---|
| [MR-001](conventions/done/MR-001-eigene-spec-schicht.md) | Baseline-Stand `v4.0.0` (drei Spec-Straten Default) |
| [MR-002](conventions/done/MR-002-id-schema-bereichskuerzel.md) | Baseline-Stand `v5.0.0` (ID-Schema = Default) |
| [MR-003](conventions/done/MR-003-vendorter-bootstrap-sensor.md) | [MR-007](conventions/MR-007-aufloesung-mr-003.md) |
| [MR-008](conventions/done/MR-008-id-schema-deklaration.md) | Baseline-Stand `v5.0.0` (ID-Schema = Default) |
| [MR-009](conventions/done/MR-009-source-precedence-ohne-docs-user.md) | [MR-010](conventions/done/MR-010-aufloesung-mr-009.md) |
| [MR-010](conventions/done/MR-010-aufloesung-mr-009.md) | Baseline-Stand `v5.0.0` (docs/user-Rang Baseline-konform) |
| [MR-011](conventions/done/MR-011-baseline-pin-release-tag.md) | [MR-012](conventions/done/MR-012-baseline-pin-hebung.md) |
| [MR-012](conventions/done/MR-012-baseline-pin-hebung.md) | [MR-016](conventions/done/MR-016-baseline-pin-hebung-2.md) |
| [MR-014](conventions/done/MR-014-slice-adr-haus-stil.md) | Baseline-Stand `v4.0.0` (Doc-Form ist Baseline-Wahl) |
| [MR-016](conventions/done/MR-016-baseline-pin-hebung-2.md) | [MR-023](conventions/MR-023-baseline-v500.md) |
| [MR-017](conventions/done/MR-017-cache-selbst-scan.md) | [MR-019](conventions/done/MR-019-regelwerk-vendored.md) |
| [MR-018](conventions/done/MR-018-keine-templates.md) | Baseline-Stand `v5.0.0` (Bundle vendored beide Bäume) |
| [MR-019](conventions/done/MR-019-regelwerk-vendored.md) | Baseline-Stand `v5.0.0` (Vendoring = Default) |
| [MR-020](conventions/done/MR-020-template-propagation.md) | Baseline-Stand `v5.0.0` (Template-Schichtung = Default) |
| [MR-022](conventions/done/MR-022-currency-audit.md) | Baseline-Stand `v5.0.0` (Freshness-Audit = Default) |

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
<!-- Retiriert (welle-67, Etappe C): §Anforderungs-Anlege-Prozess — der Inhalt war Baseline-Duplikat (grundlagen/modul-03-spec); eingefrorene Verweise lösen über diesen Anker auf. -->
<a id="anforderungs-anlege-prozess"></a>

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
