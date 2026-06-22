# Slice slice-043: codepaths `exempt-paths` (Datei-Ventil wie ids)

**Status:** in-progress (Planung; welle-32-codepaths-exempt).

**Welle:** welle-32-codepaths-exempt (Trigger: slice-042-Nebenbefund — das
Modul `codepaths` kennt, anders als `ids`, kein `exempt-paths`; Review-Reports
zitieren naturgemäß `Datei:Zeile`/geplante Pfade und lösen `codepath-missing`
aus).

**Bezug:** Change Request an
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(Modul `codepaths`). Vorbild: das `exempt-paths`-Ventil von
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(slice-018/slice-023). **Kein ADR** (additiv, spiegelt ein bestehendes Ventil;
keine neue Import-Kante — Präzedenz slice-023 ohne ADR).

**Autor:** pt9912. **Datum:** 2026-06-22.

---

## 1. Ziel

Das Modul `codepaths` bekommt das **`exempt-paths`-Ventil** (Glob-Liste von
Dateien, die das Modul nicht prüft) — Parität zum gleichnamigen `ids`-Ventil.
Motiviert vom slice-042-Nebenbefund: Review-Reports unter `docs/reviews/`
zitieren systematisch `Datei:Zeile` und geplante Pfade; ohne Ventil erzeugt
`codepaths` dort `codepath-missing` (slice-042 musste die Report-Pfade als
Klartext schreiben). `ids` löst exakt dasselbe längst über `exempt-paths`
(`docs/reviews/**`); `codepaths` zieht nach.

## 2. Entscheidungen

- **Mechanik (spiegelt `ids`):** `model.CodepathsConfig` bekommt
  `ExemptPaths []string` (wie `model.IDPattern.ExemptPaths`). `CheckCodepaths`
  überspringt eine Datei **ganz**, wenn ihr Pfad ein Glob aus `exempt-paths`
  matcht — via dem bestehenden `ignored(file, …)`/`matchGlob`-Helfer (Glob-
  Syntax wie `scan.ignore`). Datei-Ebene; der Zeilen-Opt-out bleibt der
  bestehende `d-check:ignore`-Marker (komplementär).
- **Decode:** `applyCodepaths` (configyaml) liest `exempt-paths` in den
  `Codepaths`-Raw-Struct und übernimmt sie; Validierung analog `roots`
  (kein Pflicht-Inhalt; leere Liste = Default, byte-identisch).
- **Spec:** Change Request an
  [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
  (Lastenheft **0.23.0**) + Schärfung der `codepaths`-Algorithmus-Sektion in
  `spec/spezifikation.md` (exempt-paths-Semantik, datei-weit).
- **Abwärtskompatibel:** ohne `exempt-paths` ist der Befundsatz **byte-
  identisch** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus);
  Gegentest).
- **Dogfooding (der Payoff):** [`.d-check.yml`](../../../../.d-check.yml)
  `codepaths.exempt-paths: ["docs/reviews/**"]` — Review-Reports lösen kein
  `codepath-missing` mehr aus; künftige Reports dürfen `Datei:Zeile` als
  Inline-Code schreiben.
- **Konsistenz-Parität:** `--suggest-config ai-harness[-init]` und das
  `--print-config`-Gerüst geben `codepaths.exempt-paths` mit aus (sonst Doku-/
  Vorlagen-Drift gegenüber dem `ids`-Block).

## 3. Definition of Done

- [ ] **Lastenheft-CR**
  [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in):
  `exempt-paths` ergänzt; drei Akzeptanzkriterien (Happy/Boundary/Negative) +
  Out-of-Scope; Versions-Bump **0.23.0** + Historie-Zeile.
- [ ] `spec/spezifikation.md` (codepaths-Algorithmus): exempt-paths-Semantik
  (datei-weit, Glob wie `scan.ignore`, unabhängig von `roots`).
- [ ] `model.CodepathsConfig.ExemptPaths`; `configyaml` Decode (Raw-Struct
  `exempt-paths` + `applyCodepaths`) mit Validierung; `CheckCodepaths`
  überspringt exempte Dateien (`ignored`-Helfer).
- [ ] Tests: `configyaml_test` (`codepaths.exempt-paths` akzeptiert/validiert),
  `codepaths`-Rule-Test (exempte Datei übersprungen; **ohne** Ventil
  byte-identisch).
- [ ] `--print-config`/`--suggest-config`: `codepaths.exempt-paths` sichtbar
  (Parität zum `ids`-Block).
- [ ] Dogfooding: [`.d-check.yml`](../../../../.d-check.yml)
  `codepaths.exempt-paths: ["docs/reviews/**"]`.
- [ ] [`CHANGELOG.md`](../../../../CHANGELOG.md) (nutzersichtbar, 0.23.0).
- [ ] `make gates` grün (Coverage-Schwelle gehalten); unabhängiges Review R1;
  Closure (Slice → `done/` mit Roadmap-Flip,
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

## 4. Risiken / offene Punkte

- **Granularität:** `exempt-paths` wirkt **datei-weit** (ganze Datei
  übersprungen) — identisch zum `ids`-Ventil; Zeilen-Granularität bleibt der
  `d-check:ignore`-Marker. Bewusst, keine Vermischung.
- **Determinismus:** ohne gesetztes Ventil byte-identisch
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)) —
  Gegentest gegen den Ist-Befundsatz Pflicht.
- **Vorlagen-Drift:** `--print-config`/`--suggest-config` müssen das neue Feld
  mitführen, sonst weicht die generierte Vorlage vom dokumentierten Schema ab
  (eigener DoD-Punkt).
- **Optionaler Folgeschritt (nicht Teil der DoD):** mit dem Ventil + dem
  `docs/reviews/**`-Dogfooding könnte der slice-042-Report seine Klartext-Pfade
  wieder auf Inline-Code umstellen — kosmetisch, kein Muss.

## 5. Trigger

slice-042-Closure-Nebenbefund (2026-06-22): `codepaths` ohne `exempt-paths`-
Parität zu `ids` zwingt Review-Reports zu Klartext-Pfaden. Kein Blocker, daher
als eigener Mini-Slice ausgekoppelt.

## 6. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Produkt-Code + Spec; Greenfield-Default „Doc führt,
Code folgt").

## 7. Closure-Notiz (nach `done/`)

_(folgt mit der Umsetzung.)_
