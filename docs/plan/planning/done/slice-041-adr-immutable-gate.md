# Slice slice-041: ADR-Immutable-Gate (Accepted-ADRs nicht still ändern)

**Status:** done (abgeschlossen, welle-30-adr-immutable).

**Welle:** welle-30-adr-immutable (Trigger: Nutzer-Audit 2026-06-21 —
`AGENTS.md` §3.5 „ADRs sind nach `Accepted` immutable" ist dokumentiert,
aber nicht maschinell erzwungen; Schwester-Lücke zu slice-040).

**Bezug:** Hard Rule [`AGENTS.md` §3.5](../../../../AGENTS.md#35-adrs-sind-nach-accepted-immutable)
(eine `Accepted`-ADR wird nicht inhaltlich überschrieben; Korrekturen via
neue ADR mit `Supersedes`). Mechanik-Präzedenz: das Traceability-Gate
([ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)) — CI-Range-Check
+ lokaler Hook, eine Skript-Wahrheit. Topologie/Policy dieses Slice in einer
eigenen Prozess-ADR (entsteht mit diesem Slice).

**Autor:** pt9912. **Datum:** 2026-06-21.

---

## 1. Ziel

Ein Wächter, der **inhaltliche Änderungen an `Accepted`-ADRs** erkennt und
ablehnt — erlaubt bleiben nur reine Anhänge unter `## Geschichte` und der
`**Status:**`-Übergang (nach `Superseded by ADR-NNNN` bzw. annotiertes
`Accepted`). §3.5 ist heute dokumentiert, aber nicht erzwungen.

## 2. Entscheidungen (Nutzer, 2026-06-21)

- **Topologie:** CI-Range + lokaler Hook (zwei Quadranten, **eine**
  Skript-Wahrheit unter `tools/`, exakt wie das Traceability-Gate). CI prüft
  die Push-/PR-Range; ein lokaler `pre-commit`-Hook prüft den staged-Diff
  gegen `HEAD`. Klon-unabhängiger Backstop = CI.
- **ADR:** ja, eine eigene Prozess-/Durchsetzungs-ADR (pinnt die
  Immutable-Policy + Topologie; schärft keine Spec-Stelle — analog
  [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)).
- **Methode (robust statt roher Zeilen-Diff):** „immutable core" einer
  ADR-Datei = Inhalt **ohne** den `## Geschichte`-Abschnitt **und ohne** die
  `**Status:**`-Zeile. Geprüft wird je **modifizierte**
  `docs/plan/adr/NNNN-*.md`, deren **BASE**-Version `Status: Accepted` trägt:
  weicht `core(BASE)` von `core(HEAD)` ab → FAIL. Zusätzlich: die
  HEAD-`**Status:**`-Zeile muss zulässig bleiben (`Accepted…` oder
  `Superseded by ADR-NNNN`) — kein stiller Status-Rückfall.
  Gelöschte/umbenannte Accepted-ADRs → FAIL.
- **Grandfathering:** automatisch — nur der Range-Diff wird geprüft, nicht
  die Historie. `Proposed`-ADRs sind frei editierbar (BASE nicht `Accepted`).
  Der Index (`docs/plan/adr/README.md`) ist **nicht** betroffen
  (Supersede-Annotation lebt dort, nicht in der ADR-Datei).

## 3. Definition of Done

Geplante Artefakte (Pfad gefenced, da noch nicht existent):

```text
tools/adr-immutable-check.sh   # Waechter: core-compare + Negativ-Selbsttest
.githooks/pre-commit           # ruft das Skript im --staged-Modus
```

- [ ] `tools/adr-immutable-check.sh`: Modi `--range BASE..HEAD` (CI),
  `--staged` (Hook), `--self-test`, leer (Selbsttest + `HEAD~1..HEAD`).
  **Negativ-Selbsttest** (ohne git, auf der `classify`-Kernlogik): Körper-
  Änderung feuert; Geschichte-Anhang + Superseded-Übergang feuern **nicht**;
  Status-Rückfall feuert; `Proposed`-BASE ist frei.
- [ ] `make adr-check` (dünner Wrapper, `RANGE=a..b` wie `make trace-check`);
  **nicht** Teil von `gates`/`ci` (Commit-/Diff-Bindepunkt) — `.PHONY` + help.
- [ ] `.githooks/pre-commit` ruft `adr-immutable-check.sh --staged` (eine
  Wahrheit; `make hooks` hängt `core.hooksPath` bereits ein).
- [ ] `.github/workflows/ci.yml`: ADR-Immutable über dieselbe Range wie
  `trace-check` (billiger Schritt vor dem Docker-Build).
- [ ] Doku-Sync: `harness/README.md` §Sensors **und** `AGENTS.md` §4
  (sonst rotes `make gate-consistency`); `make hooks`-Zeile um `pre-commit`
  ergänzt.
- [ ] Prozess-ADR (Proposed → Accepted nach Review) + ADR-Index.
- [ ] `make gates` grün; Akzeptanz: Wegwerf-Fixture (Körper-Edit an Accepted)
  → Exit ≠ 0. Unabhängiges Review R1; Closure.

## 4. Risiken / offene Punkte

- **Geschichte-Abschnitts-Grenze:** die `core`-Extraktion muss den
  `## Geschichte`-Abschnitt robust abgrenzen (bis zur nächsten `## `-H2 bzw.
  EOF); eine umbenannte Geschichte-Überschrift würde den Abschnitt zum Körper
  zählen (fail-closed: dann feuert das Gate eher, kein silent-green — analog
  Heading-Guard aus [slice-040](../done/slice-040-planning-consistency-gate.md)).
- **Status-Policy v1:** erlaubt sind `Accepted…` und `Superseded by ADR-NNNN`;
  exotische Status-Werte feuern. Verfeinerung (z. B. Teil-Supersede-Grammatik)
  bleibt Folgepunkt.
- **Branch Protection** (Pflicht-Check auf dem Default-Branch) liegt außerhalb
  des Repos — ohne sie ist die CI nur *advisory* (ehrliche Restlücke, wie bei
  [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)).

## 5. Trigger

Nutzer-Audit 2026-06-21: `AGENTS.md` §3.5 dokumentiert, aber nicht
maschinell erzwungen. Aus slice-040 (Planning-Konsistenz) als eigener
Folge-Slice ausgekoppelt (andere Mechanik: ADR-Diff statt Verzeichnis-State).

## 6. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Harness-Mechanik/Doku; Greenfield-Default).

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** Neues ADR-Immutable-Gate `make adr-check`
(`tools/adr-immutable-check.sh`) erzwingt
[`AGENTS.md` §3.5](../../../../AGENTS.md#35-adrs-sind-nach-accepted-immutable):
eine `Accepted`-ADR wird nicht inhaltlich überschrieben. Methode „immutable
core" = Datei **ohne** `## Geschichte`-Abschnitt **und ohne** die
Kopf-`**Status:**`-Zeile; `core(BASE) ≠ core(HEAD)` bei BASE-Accepted → FAIL,
plus HEAD-Status-Guard und Löschen-/Umbenennen-Schutz. Zwei Quadranten, **eine**
Skript-Wahrheit (wie das Traceability-Gate): lokaler `pre-commit`-Hook
(staged-Diff) + PR-/Push-CI (Range, derselbe RANGE wie `make trace-check`);
**nicht** Teil von `make gates`/`ci` (Diff-/Commit-Zeit-Bindepunkt). Policy +
Topologie in [ADR-0016](../../adr/0016-adr-immutable-gate.md) (Accepted).

**Belege.** `make gates` grün (doc-check, lint, test, arch-check, coverage,
semgrep, gate-consistency, planning-check); `make adr-check` Selbsttest (7
Fälle) + drei adversariale Pfade grün (self-test, `--staged` an echter
Accepted-ADR, `--range` via Wegwerf-Commit). Read-only/deterministisch (liest
nur git-Objekte). Dogfooding: die Durchsetzungs-ADR reifte Proposed→Accepted
unter dem laufenden Gate (BASE-Proposed = frei) und ist ab Accept selbst
gate-geschützt.

**Review R1** ([`2026-06-21-slice-041-r1.md`](../../../reviews/2026-06-21-slice-041-r1.md)):
**MERGE-FÄHIG**, 1 MEDIUM / 2 LOW / 2 NIT. **MEDIUM behoben** (`4ce7b70`):
`core()` strippte jede `**Status:**`-Zeile dateiweit → Edit an einer
Körper-Zeile mit `**Status:**`-Beginn rutschte durch (latent); Fix: Status nur
im Kopf-Block straffen, `classify` über die erste Status-Zeile. **LOW behoben:**
Selbsttest-Fälle 6 (core-resume nach `## Geschichte`) + 7 (Körper-Status-Regress);
`--range` ohne `..` → Exit 2. **NIT behoben:** Index + §3.5 auf vierstellig
`ADR-NNNN`. Ein NIT (Meldungs-Wortlaut) won't-fix. **Kein R2** (Nutzer-Entscheid).

**Lerneintrag.** Ein Doku-Immutability-Gate, das eine Markierungszeile
(`**Status:**`) dateiweit strafft, öffnet ein silent-green, sobald dieselbe
Markierung legitim im Körper vorkommt — Struktur-Anker gehören an ihre
Position gebunden (Kopf-Block), nicht global. Dieselbe Klasse wie der
slice-040-Heading-Guard: ein fail-closed-Gate muss seine Anker **verorten**.
