# Slice slice-046: `--suggest-config ai-harness` — Auffindbarkeits-Hinweis für situative opt-in-Module

**Status:** in-progress (welle-35-suggest-opt-in-hinweis).

**Welle:** welle-35-suggest-opt-in-hinweis (Trigger: Nutzer-Frage nach slice-045 —
„wird `diagrams` in `--suggest-config ai-harness` berücksichtigt?").

**Bezug:** Schärfung an
[`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(Lastenheft **0.26.0**). **Kein ADR** (additiv, Verhaltens-Schärfung der advisory
Ausgabe; kein Architektur-Delta).

**Autor:** pt9912. **Datum:** 2026-06-23.

---

## 1. Ziel

`--suggest-config ai-harness[-init]` führt das kuratierte Default-Modulset
`[links, anchors, ids, matrix, codepaths]` und lässt die situativen opt-in-Module
(`spans`, `hostpaths`, `external`, `diagrams`) bewusst weg. `diagrams` ist dabei
**nicht ableitbar** (es braucht repo-spezifische `patterns`/`defined-in`; mit
Null-Config inert). Statt es still wegzulassen, nennt die Ausgabe die nicht
aktivierten situativen Module in einem **Kommentar mit Verweis auf
`--print-config`** — Auffindbarkeit ohne Aktivieren eines inerten Moduls.

## 2. Entscheidungen

- **Kein Auto-Aktivieren.** `diagrams` zur Out-of-Scope-Liste der nicht
  auto-aktivierten Module ergänzt (`external`/`spans`/`hostpaths`/`diagrams`) —
  konsistent mit der bestehenden Auslassung. **Nur** `diagrams` einzuhängen wäre
  inkonsistent (warum nicht spans/hostpaths?); der **Kommentar** adressiert alle
  situativen Module auf einmal.
- **Kommentar, nicht aktiv.** Ein YAML-Kommentar (`#`) — dekodiert sauber über den
  eigenen Parser (Auffindbarkeits-AK prüft Decode + Substring). Kein Aktivieren
  eines patterns-losen, inerten `diagrams`.
- **Spec:** Schärfung
  [`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten):
  Out-of-Scope-Liste + neue Auffindbarkeits-AK; Version 0.26.0 + §7. Kein
  `--print-config`-Delta (dort steht `diagrams` seit slice-045).

## 3. Definition of Done

- [x] **Lastenheft-CR** (0.26.0): Out-of-Scope `diagrams` + Auffindbarkeits-AK + §7.
- [ ] Code (`internal/hexagon/core/app/suggest.go`): Kommentar-Zeile nach dem
  `modules:`-Eintrag der ai-harness-Vorlage.
- [ ] Test: `TestCLI006_AiHarness_Happy` prüft den Kommentar (`(external, spans,
  hostpaths, diagrams)` + `--print-config`) und die fortbestehende Decode-Gültigkeit.
- [ ] `make gates` grün; unabhängiges Review R1; Closure (Move nach `done/` +
  Roadmap-Flip, [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

## 4. Risiken / offene Punkte

- **Derive-Modus (`--suggest-config <quelle>`)** probt `codepaths`/`spans`/
  `hostpaths` und schlägt feuernde vor; `diagrams` kann dort nie feuern (keine
  patterns). Ein analoger Hinweis dort ist denkbar, aber **out of scope** dieses
  Slice (die Nutzer-Frage galt der ai-harness-Vorlage) — bewusst eng gehalten.
- **Determinismus** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  statischer Kommentar, byte-stabil.

## 5. Trigger

Nutzer-Frage am slice-045-Closure (2026-06-23). Konsequenz aus der bewussten
Kuratierung der ai-harness-Vorlage (Default-Set ohne situative opt-ins) plus der
Nicht-Ableitbarkeit von `diagrams`.

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code `app` + Spec; „Doc führt, Code folgt"). Keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

_(Wird beim Move nach `done/` gefüllt — Belege: `make gates`-Ausgabe,
unabhängiges Review R1, ggf. Release.)_
