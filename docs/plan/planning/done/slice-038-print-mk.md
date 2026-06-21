# Slice slice-038: `--print-mk` — `d-check.mk` ausgeben (include-bare Integration)

**Status:** done (Closure 2026-06-21; Review R1+R2, kein ADR nötig).

**Welle:** welle-28-print-mk (Trigger: a-check-Bootstrap 2026-06-20 — a-check
bindet das Doku-Gate per **handgepflegtem** `d-check.mk` ein; der Pin lebt
dadurch im Konsumenten statt in d-check). *(welle-Nummer von ursprünglich
„welle-27" korrigiert — welle-27 ist welle-27-rtm-trace, slice-036.)*

**Bezug:**
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
(die umgesetzte Anforderung, CR 0.22.0),
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(verteiltes Image),
[`DC-FA-CLI-005`](../../../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
(`--print-config` — read-only-Generator als Vorbild),
[ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)
(Digest-Pins),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(read-only).

**Autor:** pt9912. **Datum:** 2026-06-20.

---

## 1. Ziel

Ein read-only-Modus `d-check --print-mk`, der ein **`d-check.mk`** auf stdout
ausgibt — ein include-bares Makefile-Fragment mit dem **aktuell
digest-gepinnten** d-check-Image und einem `doc-check`-Target. Konsumenten
`include d-check.mk` und liefern ihre `.d-check.yml` — **keine Skript-/
Recipe-Kopie**. Reiht sich in die read-only-Generatoren
(`--print-config`/`--suggest-config`) ein.

## 2. Hintergrund / Symmetrie

a-check definiert dasselbe Muster bereits für sich selbst
(`AC-FA-DIST-001`: Image + `--print-mk` + `a-check.mk`). Dieser Slice zieht
d-check nach: beide Stack-Tools werden „include-and-configure"-Distributables.
a-check pflegt das `d-check.mk` heute interim von Hand — `--print-mk` ersetzt
das, und der Pin lebt (richtig) in d-check.

## 3. Zu entscheiden (im Slice)

- **Digest-Quelle:** woher kennt das Binary den Pin seines *eigenen*
  Release-Images? Build-Zeit eingebettet (`-ldflags`-Var, am Release gesetzt)
  vs. zur Laufzeit unbekannt → Henne-Ei beim Erst-Release (analog
  [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)). Wahrscheinlich
  eingebettete Version + Digest, beim Tag-Build gesetzt.
- **Target-Umfang:** nur `doc-check`, oder zusätzlich ein aggregierendes
  `gates`? Empfehlung: minimal `doc-check`; Konsumenten komponieren selbst.
- **Variablen-Override:** `DCHECK_IMAGE ?= …` (überschreibbar) wie im
  a-check-Interim.
- **Lastenheft-CR:** neue Anforderung (DIST- oder CLI-Familie) für
  `--print-mk` mit Happy/Boundary/Negative; oder Erweiterung von
  [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image).

## 4. Definition of Done (vorläufig)

- [x] Lastenheft-CR (0.22.0) + Spezifikation (`spec/spezifikation.md` §…a)
  für `--print-mk` (Bezug oben).
- [x] Modus im Paket `cli` (wie `--print-config`, repo-frei), gibt
  `d-check.mk` auf stdout aus, **schreibt nichts**
  ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)),
  deterministisch
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- [x] Ausgabe trägt einen **version-gepinnten** Image-Ref (ins Binary via
  `-ldflags -X` eingebettet); Digest via `DCHECK_IMAGE`-Override (Henne-Ei,
  Auftraggeber-Entscheidung) statt eingebettetem Digest.
- [x] Akzeptanztests (`TestCLI038_PrintMK*`); Dogfooding
  `--print-mk | make -n -f - doc-check` parst grün; unbekanntes Flag → Exit 2.
- [x] `docs/user/operations.md` + `CHANGELOG` ergänzt; `make gates` grün;
  **kein ADR** (Version-Tag-Default konsistent mit der ratifizierten
  Konsum-Pin-Politik).
- [x] Unabhängiges Review R1 (0 HIGH/0 MEDIUM/1 LOW) + R2 (bestätigt);
  Closure.

## 5. Risiken / offene Punkte

- **Pin-Aktualität:** der ausgegebene Digest muss das *tatsächlich
  veröffentlichte* Image treffen — Kopplung an die Release-Pipeline
  (`release.yml`), nicht an den Quellbaum.
- **Henne-Ei beim Erst-Release** des Features (Digest erst nach Push
  bekannt) — gleiche Klasse wie
  [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md); dort gelöste
  Mechanik wiederverwenden.

## 6. Trigger

a-check-Bootstrap 2026-06-20: doc-check wurde per handgepflegtem `d-check.mk`
eingebunden (`include`). Das Zielbild ist d-check-eigene Ausgabe via
`--print-mk`, symmetrisch zu a-checks `AC-FA-DIST-001`.

## 7. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (CLI-/Core-/Doku-Arbeit; Greenfield-Default).

## 8. Closure-Notiz (nach `done/`)

**Umsetzung.** Neuer read-only-Generator `--print-mk`
([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)):
gibt ein include-bares `d-check.mk` auf stdout (wie `--print-config`,
repo-frei via `earlyGenerators`). `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v<version>`
— Version via `-ldflags -X` in der Dockerfile-build-Stage eingebettet,
gespeist aus dem bestehenden `make ci VERSION=…` (kein `release.yml`-Edit) —
plus ein `doc-check`-Target. **Digest via `DCHECK_IMAGE`-Override** (Henne-Ei:
das Binary kennt seinen eigenen Digest nicht); Version-Tag-Default konsistent
mit der ratifizierten Konsum-Pin-Politik. **Kein ADR** (additiv, keine neue
Import-Kante; `print_mk.go` importiert nur `fmt`).

**Belege.** `make gates` grün (doc-check, lint, test, arch-check, coverage
94,20 %, semgrep 29/0, gate-consistency); read-only
([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)),
deterministisch
([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
Dogfooding: `d-check --print-mk | make -n -f - doc-check` parst + expandiert;
Build mit `--build-arg VERSION=9.9.9-test` → `:v9.9.9-test` (R1/R2 verifiziert).

**Review R1** (`docs/reviews/2026-06-21-slice-038-print-mk.md`):
0 HIGH/0 MEDIUM/1 LOW/1 INFO. **R2**
(`docs/reviews/2026-06-21-slice-038-r2-verifikation.md`): R1-Kern bestätigt,
keine neuen blockierenden Punkte. **LOW-1** (`--print-mk` kombiniert still mit
`--json`/`--doctor`) won't-fix: **identisch** zum Schwester-Generator
`--print-config` (R2: kein Fall schlimmer; `comboError` greift vor
`earlyGenerators`, daher `--print-mk --json --yaml` → Exit 2). INFO won't-fix.

**Lerneintrag.** Ein Binary kann seinen eigenen Image-Digest nicht
selbst-einbetten (der Digest hasht das Binary) — die ehrliche Lösung ist der
Version-Tag als self-known Default plus `DCHECK_IMAGE`-Override für den
strikten Digest-Pin. Und: ein neuer CLI-Modus kann die Komplexitäts-Schwelle
von `Run` kippen — auslagern (`earlyGenerators`) statt inline.
