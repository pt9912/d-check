# Slice slice-046: `--suggest-config ai-harness` — Auffindbarkeits-Hinweis für situative opt-in-Module

**Status:** done (abgeschlossen, welle-35-suggest-opt-in-hinweis).

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

**Umsetzung.** Schärfung
[`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten):
die `--suggest-config ai-harness[-init]`-Ausgabe (`renderHarnessSuggestion` in
`internal/hexagon/core/app/suggest.go`) trägt nach dem `modules:`-Eintrag einen
YAML-Kommentar, der die nicht aktivierten situativen opt-in-Module
(`external`/`spans`/`hostpaths`/`diagrams`) nennt und auf `--print-config`
verweist — Auffindbarkeit ohne Aktivieren eines inerten Moduls (`diagrams` braucht
repo-spezifische `patterns`/`defined-in`, lässt sich nicht ableiten). Lastenheft
0.26.0 (Out-of-Scope + Auffindbarkeits-AK). **Kein ADR** (Schärfung der
advisory Ausgabe).

**Belege.** `make gates` grün; `TestCLI006_AiHarness_Happy` prüft Kommentar +
fortbestehende Decode-Gültigkeit. Minor-Release **v0.26.0** auf GHCR (Run
`28040897654` grün, Tags `v0.26.0`+`latest`), Digest-Pin
`ghcr.io/pt9912/d-check@sha256:19d53a26d8d82a919015a8befe24f852bd61f2ddea58bd29e3f4cf944a8403f3`.
Handbuch 1.8/v0.26.0, CHANGELOG `[0.26.0]`.

**Review.** R1 ACCEPT (0 HIGH/0 MEDIUM/1 LOW); F-1 behoben — `external` durchgängig
ergänzt, sodass Kommentar == AK == Out-of-Scope.

**Lerneintrag.** Die kuratierte ai-harness-Vorlage lässt situative opt-ins bewusst
weg; ein Modul, das ohne repo-spezifische Config inert ist (`diagrams`), gehört
**nicht** in eine generische Vorlage — der richtige Schnitt ist ein Verweis aufs
Voll-Schema (`--print-config`), nicht stilles Aktivieren.
