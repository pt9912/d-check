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
- **Stand:** Template-Set 2026-06
- **Datum der Adoption:** 2026-06-10

## Adoptierte Konventions-Quellen

- **Extern (Lehrmaterial):**
  [`ai-harness-course`](https://github.com/pt9912/ai-harness-course)
  (Templates: `lab/templates/`, Konventionen:
  `kurs/de/grundlagen/konventionen.md`).
- **Extern (Agenten-Destillat):**
  [`agents-regelwerk.md`](https://raw.githubusercontent.com/pt9912/ai-harness-course/main/kurs/de/agents-regelwerk.md)
  — operatives Regelwerk für Code-Agenten ohne Didaktik; derivativ,
  bei Konflikt gilt das Lehrmaterial.
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
  als Provenienz stehen).

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
- **Auflösungs-Trigger:** permanent.

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

| Klasse | Form | Bedeutung | Beispiel |
|---|---|---|---|
| DC-Bindung | `DC-…` | Gate prüft eine konkrete Lastenheft-Anforderung | [`DC-QA-02`](../spec/lastenheft.md#dc-qa-02--determinismus) für den Determinismus-Test in `make test`; [`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) für `make arch-check` und das Netzlos-Gate in `make doc-check` |

## Modus-Deklaration pro Sub-Area

| Sub-Area (Pfad / Modul) | Modus | Begründung | Graduation-Bedingung / Folge-Slice |
|---|---|---|---|
| `*` (Default für gesamtes Repo) | Greenfield | Projekt startet spec-first; Doc führt, Code folgt | n/a (GF) |
| `tools/harness/` | Greenfield | adoptierte Harness-Mechanik, konventionsgetragen über [`MR-004`](#mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild) | n/a (GF) |
