# Slice slice-032: semgrep als hermetisches Gate

**Status:** done (Closure 2026-06-20).

**Welle:** welle-21-semgrep-gate (Trigger: Nutzer-Entscheid „semgrep ins
Gate", 2026-06-19; baut auf dem eigenständigen `make semgrep`-Target auf).

**Bezug:**
[ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md) (Entscheidung),
[ADR-0006](../../adr/0006-lint-profil-solid.md) (golangci-lint als erste
statische Analyse),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(reproduzierbar),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(netzlos).

**Autor:** pt9912. **Datum:** 2026-06-19.

---

## 1. Ziel

semgrep wird ein **reproduzierbares Gate mit netzlosem Scan** in
`make gates`: ein **gepinntes, lokal gecachtes** Regelset (kuratierte
Sprach-Packs, **nicht** ins Repo aufgenommen) statt `--config auto`, Scan
mit `--network none` + gepinntem Image. Identische Eingabe ⇒ identische
Befunde, kein Netz im Scan. Ergänzt golangci-lint sprachübergreifend
(bash, Dockerfile).

## 2. Definition of Done

- [x] **Regel-Cache (nicht im Repo):** `semgrep/semgrep-rules` an einem
  **festen Commit-Pin** wird von `tools/semgrep.sh` einmalig in einen Cache
  **außerhalb des Repos** (`~/.cache/d-check/semgrep-rules/<commit>/`,
  Override `SEMGREP_RULES_CACHE`) geholt (wie `go mod`/Image-Pull) und per
  `--config` genutzt; kuratierter Umfang **`go/lang/security`** (bash/
  dockerfile/secrets ausgelassen — Probelauf go+bash+dockerfile: 13 Treffer,
  alle FP). Pin + Quelle stehen in `tools/semgrep.sh`.
- [x] **`tools/semgrep.sh`:** von `--config auto` auf das **gecachte**
  Regelset + `--network none` (+ `--disable-version-check`) umgestellt,
  Image gepinnt; **kein** `--exclude-rule` nötig (`go/lang/security`: 0
  Befunde).
- [x] **`make gates`:** semgrep aufgenommen (damit auch `make ci`); `--error`
  bricht das Gate bei Befund. Eigenständiges `make semgrep` bleibt.
- [x] **Doku/Konsistenz:** AGENTS §4-Zeile auf „Bestandteil von `gates`"
  aktualisiert; `harness/README.md` Sensors-Tabelle um semgrep ergänzt;
  `gates`-Beschreibung nennt semgrep. `make gate-consistency` grün.
- [x] **Verifikation:** `make gates` grün **inkl. semgrep offline**
  (`--network none`, lokale Regeln); reproduzierbar (zwei Läufe identisch).
- [x] **[ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md)** auf `Accepted`; `make gates` grün; unabhängiges
  Review R1; Closure-Notiz.

## 3. Plan (vor Code)

| Datei | Art | Begründung |
|---|---|---|
| `tools/semgrep.sh` | update | Regel-Cache am Commit-Pin holen (außerhalb Repo), `--config` lokal + `--network none`, Image-Pin; Regel-Nachweis statt FP-Ausschluss (0 Befunde) |
| `Makefile` | update | semgrep in `gates` einhängen |
| `AGENTS.md`, `harness/README.md` | update | semgrep als Gate dokumentieren (gate-consistency beide Richtungen) |
| `docs/plan/adr/0010-semgrep-hermetisches-gate.md`, `docs/plan/adr/README.md` | update | ADR auf Accepted, Index-Status |

Die ADR-Entscheidung ist mit dem Nutzer-Entscheid gesetzt; kein
Lastenheft-Bezug (Prozess-/Qualitäts-Gate, kein `DC-*`-Vertrag — wie das
Lint-Profil/`make lint`).

## 4. Trigger

Nutzer-Entscheid „Option 3 / Variante A" (2026-06-19); ADR (Proposed, Bezug oben).

## 5. Closure-Trigger

DoD vollständig inkl. grüner Gates (semgrep offline), Review R1 und der
ADR auf `Accepted`.

## 6. Risiken und offene Punkte

- **Cache-Hol-Schritt (Setup-Netz):** das erste `make gates`/`make semgrep`
  holt den Regel-Cache am Commit-Pin (Netz, wie `go mod`/Image-Pull),
  danach offline; der Cache liegt **außerhalb des Repos** (keine
  Interferenz mit `go list ./...` oder dem d-check-Selbstscan). Pin-Hebung
  ist ein bewusster Commit (Befund-Delta sichtbar), kein Auto-Update.
- **Image-Pull/Cache vs. `--network none`:** Image-Pull und Cache-Holen
  laufen über Host-Netz **vor** dem Scan; der Scan selbst läuft im
  Container netzlos — Hermetik des *Scans* gewahrt (wie die
  Basis-Image-Pulls der übrigen Stages).
- **Rauschen im Gate:** kuratierte Packs gewählt (nicht Voll-Mirror), um
  „grün=Boden" nicht durch False-Positive-Last zu untergraben; neue FPs →
  zentrale Ausnahme mit Why.
- **Determinismus** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  gepinntes Image + gepinntes Regelset ⇒ stabiler Befundsatz.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** `tools/semgrep.sh` von `--config auto`-Standalone auf das
hermetische Gate umgebaut: gepinntes Image `semgrep/semgrep:1.167.0` +
gepinnter Regel-Commit `d41fb34cf74466e2878af5f268ebf54466a04541` aus
`semgrep/semgrep-rules`, einmalig in `~/.cache/d-check/semgrep-rules/<commit>/`
(Override `SEMGREP_RULES_CACHE`) geholt (Host-`git`, idempotent über den
Commit-Schlüssel, atomar via `.tmp`+`mv`), Scan
`docker run --network none --disable-version-check --metrics off
--error --config /rules/go/lang/security`. In `make gates` eingehängt
(vor `gate-consistency record-gates`); `AGENTS.md` §4, `harness/README.md`
Sensors-Tabelle + `gates`-Zeile, [ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md) (Accepted) und ADR-Index nachgezogen.

**Umfang-Entscheidung (Nutzer, 2026-06-20).** `go/lang/security` (55 Regeln,
0 Befunde) statt go+bash+dockerfile (13 Treffer, alle FP) — kein
`--exclude-rule` nötig. Cache-Ablage als Host-XDG-Cache bestätigt.

**Belege.**
- `make gates` grün (Exit 0):
  `doc-check + lint + test + arch-check + coverage-gate + semgrep +
  gate-consistency green`.
- semgrep offline: `Ran 55 rules on 26 files: 0 findings`, **3,1 s**
  (Cache-Hit, `--network none`); zwei Läufe **byte-identisch**
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- Cache liegt außerhalb des Repos (24 MB), zweiter/dritter Lauf ohne
  Fetch-Meldung (offline;
  [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

**Review-Runde R1** (`docs/reviews/2026-06-20-slice-032-semgrep-gate.md`,
Verdikt anfangs *Blockiert*): HIGH-1 — Gate wies nicht nach, dass `--config`
Regeln lud (regel-leeres/umbenanntes Subset ⇒ „0 findings" ⇒ stilles Grün).
Behoben: Scan ohne `exec`, Output capturen, positive `Ran N rules`-Zeile
(N≥1) als Pflicht-Nachweis erzwingen ⇒ leerer Cache liefert jetzt Exit 2
(empirisch belegt) statt Exit 0. MEDIUM-1 — ADR-Index betitelte das Regelset
als „vendortes" (Gegenteil der Entscheidung); korrigiert auf „lokal gecachtes
… Regelset". INFO-1 (Image-Tag statt Digest) won't-fix — konventionskonform
zur Pin-Politik (`GO_VERSION`/`GOLANGCI_LINT_VERSION`); als Deferred notiert.

**Lerneintrag.**
- semgreps Versions-Ping läuft unter `--network none` in einen ~2-min-Timeout
  → `--disable-version-check` ist Pflicht für netzlose Gates (3 s statt 2 min).
- `--error` schützt nur gegen Befunde, **nicht** gegen 0 geladene Regeln; ein
  Gate über externe Regeln braucht einen expliziten „Regeln liefen"-Nachweis,
  sonst ist Abdeckungsverlust ein stilles Grün.

**Deferred (eigener Mini-Slice).** Dockerfile-**Digest**-Pinning (echte
Reproduzierbarkeits-Verschärfung über alle Image-Pins, INFO-1 inklusive).

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Tooling-/Doku-Arbeit; Greenfield-Default).
