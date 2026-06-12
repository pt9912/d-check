# Slice slice-014: Rollout — die restlichen acht Alt-Tool-Vorkommen

**Status:** open.

**Welle:** welle-05-rollout (Abschluss).

**Bezug:** [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
(Anforderung gilt für **alle zwölf** Quell-Tool-Vorkommen; die
Messmethode — drei Piloten — ist seit
[slice-012](../done/slice-012-pilot-migrationen.md) erfüllt, dieser
Slice vervollständigt den Beleg auf 12/12);
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(repo-spezifisches Verhalten per Config).

**Autor:** pt9912. **Datum:** 2026-06-12.

---

## 1. Ziel

Die acht verbleibenden Alt-Tool-Vorkommen in den Schwester-Repos
nutzen das veröffentlichte `d-check`-Image (Digest-Pin v0.2.0) statt
ihrer Tool-Kopie — mit den in slice-012 erprobten Familien-Rezepten:

- **Shell ×4** (`verify-doc-refs.sh`): b-trace, m-trace, cmake-xray,
  k-deskflight — vier untereinander abweichende Kopien (38–119
  Zeilen, der Drift-Beleg der Konsolidierung).
- **Python ×2** (`check_refs.py`): c-hsm-doc (Familien-Ursprung,
  `HSM-*`-Muster), grid-gym — beide ohne u-boot-Vollausbau.
- **JS ×2** (`docs-check.js`): b-cad, euler-fourier-hilbert
  (Math-Validierung bleibt als Rest-Sensor — Mathematik-Prüfung ist
  laut Lastenheft §5 global out-of-scope).

## 2. Definition of Done

- [ ] Pro Repo: passende `.d-check.yml` (Familien-Rezept aus
  slice-012; `ids`-Muster mit `\b`-Wortgrenzen, Scan-Scope
  dokumentiert gespiegelt oder bewusst erweitert).
- [ ] Pro Repo ein dokumentierter **Vergleichslauf** Alt-Tool vs.
  `d-check` auf demselben Repo-Stand; Differenzen triagiert nach dem
  slice-012-Schema (echter Mehr-Befund → Fix im Ziel-Repo oder
  dortige Ignore-Regel; False-Positive → Spec-Fortschreibung/Fix
  hier, vor Abschluss).
- [ ] Pro Repo: Make-/CI-Schritt auf den v0.2.0-Digest-Pin
  umgestellt; Alt-Tool-Kopie gelöscht, geschrumpft (Fall
  euler-fourier-hilbert: Math-Rest-Sensor; ggf. b-cad) oder
  deprecated — Entscheidung beim Ziel-Repo, dokumentiert.
- [ ] Vergleichstabelle (Repo, Familie, Alt-Befunde, d-check-Befunde,
  Differenzen + Triage) in der Closure-Notiz — zusammen mit
  slice-012 und dem Eigenlauf
  ([`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding))
  ist [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
  damit über alle zwölf Vorkommen belegt.
- [ ] `make gates` grün (dieses Repo);
  [`CHANGELOG.md`](../../../../CHANGELOG.md); Closure-Notiz mit
  Steering-Loop-Lerneintrag — erfüllt zugleich den
  welle-05-Closure-Trigger.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `<repo>/.d-check.yml` (8×, extern) | neu | Familien-äquivalente Konfiguration |
| `<repo>`-Makefile/CI (8×, extern) | update | Digest-gepinnter `d-check`-Step, Alt-Tool-Ablösung |
| Closure-Notiz dieses Slices | neu | 12/12-Vergleichstabelle als `DC-QA-04`-Beleg |
| [`CHANGELOG.md`](../../../../CHANGELOG.md) | update | Rollout-Stand |

## 4. Trigger

slice-012 done (erfüllt 2026-06-12) — Release v0.2.0 mit allen sechs
Modulen ist auf GHCR, die drei Familien-Rezepte sind erprobt, die
Triage-Systematik steht. Kein neues Release nötig.

## 5. Closure-Trigger

DoD vollständig (acht dokumentierte Vergleichsläufe, acht umgestellte
Make-/CI-Steps) + Closure-Notiz — erfüllt zugleich den
welle-05-Closure-Trigger.

## 6. Risiken und offene Punkte

- Erwartung aus dem slice-012-Lerneintrag: Pro Repo können echte
  Mehr-Befunde der CommonMark-Span-Klasse auftauchen (in u-boot
  ~100 versteckte Defekte) — Aufwand pro Repo schwankt entsprechend.
- euler-fourier-hilbert ist der anspruchsvollste Fall (Math-Anteil
  → Rest-Sensor-Schnitt analog Kurs-Repo) — bewusst als letzter.
- Reihenfolge nach aufsteigendem Risiko: Shell ×4 → Python ×2 →
  JS ×2.
- Die Pilot-Commits aus slice-012 (drei Repos) sind noch ungepusht
  (Push macht der Repo-Owner); die dortigen CIs laufen bis dahin auf
  altem Stand.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Dieses Repo: GF (Doku-/Beleg-Arbeit). Ziel-Repos: außerhalb des
Geltungsbereichs dieses Harness (deren Modus gilt dort).
