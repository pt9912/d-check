# Slice slice-012: Pilot-Migrationen (drei Tool-Familien)

**Status:** open.

**Welle:** welle-04-distribution-und-migration (Abschluss).

**Bezug:** [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
(Messmethode: Pilot-Migration in mindestens drei Repos);
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(repo-spezifisches Verhalten per Config); Meilenstein M3 (Teil 2:
≥ 1 Repo migriert — hier: drei).

**Autor:** pt9912. **Datum:** 2026-06-11.

---

## 1. Ziel

Drei Schwester-Repos — je ein Vertreter der Shell-Familie
(`verify-doc-refs.sh`), der Python-Familie (`check_refs.py`, inkl.
u-boot-Vollausbau) und der JS-Familie (`docs-check.js`) — nutzen das
veröffentlichte `d-check`-Image statt ihrer Tool-Kopie; die
Vergleichsläufe sind dokumentiert und belegen `DC-QA-04`.

## 2. Definition of Done

- [ ] Pro Pilot-Repo: passende `.d-check.yml` geschrieben
  (Scan-Wurzeln, Module, ggf. ids-/matrix-Regeln äquivalent zur
  Alt-Tool-Abdeckung).
- [ ] Pro Pilot-Repo ein dokumentierter **Vergleichslauf** Alt-Tool
  vs. `d-check` auf demselben Repo-Stand:
  `d-check` meldet mindestens dieselben echten Befunde und erzeugt
  keine False-Positives, die eine bislang grüne CI brechen
  (Akzeptanz-Kern von
  [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools));
  Differenzen werden triagiert: echter Mehr-Befund (Fix im
  Ziel-Repo oder dortige Ignore-Regel) vs. False-Positive
  (Spec-Fortschreibung/Fix hier, vor Abschluss des Slices).
- [ ] Pro Pilot-Repo: CI-/Make-Schritt auf
  `ghcr.io/pt9912/d-check@sha256:…` (Digest-Pin gemäß
  Release-Doku aus slice-011) umgestellt; die Alt-Tool-Kopie ist im
  Ziel-Repo gelöscht oder als deprecated markiert (Entscheidung liegt
  beim Ziel-Repo, wird dokumentiert).
- [ ] Vergleichstabellen (Repo, Familie, Alt-Befunde, d-check-Befunde,
  Differenzen + Triage) in der Closure-Notiz — der zweite und dritte
  Datenpunkt nach dem Eigenlauf aus
  [`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding).
- [ ] `make gates` grün (dieses Repo);
  [`CHANGELOG.md`](../../../../CHANGELOG.md); Closure-Notiz mit
  Steering-Loop-Lerneintrag — damit ist zugleich der
  welle-04-Closure-Trigger erfüllt (M3 erreicht).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `<pilot-repo>/.d-check.yml` (3×, extern) | neu | Familien-äquivalente Konfiguration |
| `<pilot-repo>`-CI/Makefile (3×, extern) | update | Digest-gepinnter `d-check`-Step, Alt-Tool-Ablösung |
| Closure-Notiz dieses Slices | neu | Vergleichsläufe + Triage als `DC-QA-04`-Beleg |
| [`CHANGELOG.md`](../../../../CHANGELOG.md) | update | Migrations-Stand |

## 4. Trigger

slice-011 done — die Pilot-Repos konsumieren das **veröffentlichte**
Image per Digest-Pin, nicht einen lokalen Build.

## 5. Closure-Trigger

DoD vollständig (drei dokumentierte Vergleichsläufe, drei umgestellte
CI-Steps) + Closure-Notiz — erfüllt zugleich den
welle-04-Closure-Trigger und M3.

## 6. Risiken und offene Punkte

- Die Hauptarbeit liegt **außerhalb dieses Repos**; hier landen nur
  die Belege. Die Ziel-Repos haben eigene Gates/Prozesse — Änderungen
  dort folgen deren Konventionen.
- u-boot-Vollausbau (`check_refs.py`) ist der anspruchsvollste
  Vergleich (größter Funktionsumfang der Alt-Familie); Differenzen
  dort können Spec-Fortschreibungen hier auslösen — bewusst als
  letzter der drei Vergleiche einplanen.
- `DC-QA-04` verlangt „keine False-Positives, die eine grüne CI
  brechen" — bei legitimen Mehr-Befunden ist die Triage-Doku
  entscheidend (Mehr-Befund ≠ False-Positive).

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Dieses Repo: GF (Doku-/Beleg-Arbeit). Pilot-Repos: außerhalb des
Geltungsbereichs dieses Harness (deren Modus gilt dort).
