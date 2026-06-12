# Slice slice-014: Rollout — die restlichen neun Alt-Tool-Vorkommen

**Status:** in-progress.

**Welle:** welle-05-rollout (Abschluss).

**Bezug:** [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
(Anforderung gilt für **alle dreizehn** Quell-Tool-Vorkommen —
Inventur-Nachtrag Lastenheft 0.4.0; die
Messmethode — drei Piloten — ist seit
[slice-012](../done/slice-012-pilot-migrationen.md) erfüllt, dieser
Slice vervollständigt den Beleg auf 13/13);
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(repo-spezifisches Verhalten per Config).

**Autor:** pt9912. **Datum:** 2026-06-12.

---

## 1. Ziel

Die neun verbleibenden Alt-Tool-Vorkommen in den Schwester-Repos
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
- **Eigenständige Linie ×1** (`check_markdown_links.py`): bess-ems —
  Inventur-Nachtrag (Lastenheft 0.4.0, Fund Auftraggeber 2026-06-12);
  der Host-Pfad-Prosa-Check bleibt dort als Rest-Sensor.

## 2. Definition of Done

- [x] Pro Repo: passende `.d-check.yml` (Familien-Rezept aus
  slice-012; `ids`-Muster mit `\b`-Wortgrenzen, Scan-Scope
  dokumentiert gespiegelt oder bewusst erweitert).
- [x] Pro Repo ein dokumentierter **Vergleichslauf** Alt-Tool vs.
  `d-check` auf demselben Repo-Stand; Differenzen triagiert nach dem
  slice-012-Schema (echter Mehr-Befund → Fix im Ziel-Repo oder
  dortige Ignore-Regel; False-Positive → Spec-Fortschreibung/Fix
  hier, vor Abschluss).
- [x] Pro Repo: Make-/CI-Schritt auf den v0.2.x-Digest-Pin
  umgestellt; Alt-Tool-Kopie gelöscht, geschrumpft (Fälle
  euler-fourier-hilbert: Math-Rest-Sensor; bess-ems:
  Host-Pfad-Rest-Sensor; ggf. b-cad) oder
  deprecated — Entscheidung beim Ziel-Repo, dokumentiert.
- [x] Vergleichstabelle (Repo, Familie, Alt-Befunde, d-check-Befunde,
  Differenzen + Triage) in der Closure-Notiz — zusammen mit
  slice-012 und dem Eigenlauf
  ([`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding))
  ist [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
  damit über alle dreizehn Vorkommen belegt.
- [x] `make gates` grün (dieses Repo);
  [`CHANGELOG.md`](../../../../CHANGELOG.md); Closure-Notiz mit
  Steering-Loop-Lerneintrag — erfüllt zugleich den
  welle-05-Closure-Trigger.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `<repo>/.d-check.yml` (9×, extern) | neu | Familien-äquivalente Konfiguration |
| `<repo>`-Makefile/CI (9×, extern) | update | Digest-gepinnter `d-check`-Step, Alt-Tool-Ablösung |
| Closure-Notiz dieses Slices | neu | 13/13-Vergleichstabelle als `DC-QA-04`-Beleg |
| [`CHANGELOG.md`](../../../../CHANGELOG.md) | update | Rollout-Stand |

## 4. Trigger

slice-012 done (erfüllt 2026-06-12) — Release v0.2.0 mit allen sechs
Modulen ist auf GHCR, die drei Familien-Rezepte sind erprobt, die
Triage-Systematik steht. Kein neues Release nötig.

## 5. Closure-Trigger

DoD vollständig (neun dokumentierte Vergleichsläufe, neun umgestellte
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

**Abgeschlossen:** 2026-06-12. Mit den drei slice-012-Piloten und dem
Eigenlauf ([`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding))
sind alle **dreizehn** Quell-Tool-Vorkommen migriert —
[`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
ist vollständig belegt. Zwischenrelease **v0.2.1**
(`sha256:7bdef75b…e050eb`) aus dem Adoptions-Befund (s. u.).

### Vergleichsläufe (Rollout, 9 Vorkommen)

| Repo | Familie/Variante | Alt | d-check | Differenzen + Triage | Alt-Tool-Entscheidung | Commits |
|---|---|---|---|---|---|---|
| b-trace | Shell (38 Z.) | 0 | 0 | keine | gelöscht | b51bcb2 |
| m-trace | Shell (92 Z., Ursprung der d-migrate-Variante) | 0 | 0 | keine | gelöscht | 4b0b1c1 |
| cmake-xray | Shell (92 Z., whitespace-identische m-trace-Kopie) | 0 | 0 | keine | gelöscht | 11e70ad |
| k-deskflight | Shell (119 Z., ADR-verankert) | 0 | 2 → 0 | 2 echte Mehr-Befunde (veraltete Roadmap-Anker — Alt schnitt Fragmente ab); Gate-Ablösung per dortiger **ADR 0016**, löst zugleich deren Docker-only-Carveout (Host-Bash) auf | gelöscht | b2ae4fb, 28ab762 |
| c-hsm-doc | Python (159 Z., Familien-Ursprung) | 0 | 0 | keine im Bestand; 4 Folge-Links auf das gelöschte Skript vom ersten d-check-Lauf gefangen und nachgezogen | gelöscht | 64dcc31 |
| grid-gym | Python (154 Z.) | 0 | 4 → 0 | 4 echte Mehr-Befunde: 3 veraltete ADR-Verweise (CONTRIBUTING.md, außerhalb des Alt-Scopes) + 1 auf GitHub kaputter Roadmap-Link **hinter einem mehrzeiligen Commit-Span** (CommonMark-Klasse aus dem slice-012-Lerneintrag, drittes Repo mit dem Muster) | gelöscht inkl. Dockerfile-Stage | 766ae8c |
| euler-fourier-hilbert | JS (520 Z., Familien-Basis) | 0 (318 Links, 1259 Math-Blöcke) | 0 | keine — auch die Slug-Differenz (Alt kollabiert Bindestriche) blieb im Bestand folgenlos | geschrumpft: **Math-Rest-Sensor** (MathJax + GitHub-Quirks, 520 → ~340 Z.) | 4f6b584 |
| b-cad | JS (469 Z., vendorte Kurs-Kopie) | 0 | 8 → 0 | 8 echte Mehr-Befunde aus **adaptierten** codepaths-Wurzeln (die Kopie prüfte die toten Kurs-Präfixe `lab/`/`kurs/`): 2 Pfad-Fixes, 6 begründete Marker (geplante Artefakte, d-migrate-Verweise, Plan-Alternative). Build-Stage-Konvention (kein Bind-Mount) erhalten: `FROM d-check` + `COPY` + `RUN`; dortiges `MR-003` → `MR-007` aufgelöst | gelöscht | 77fbd17, dfbb44c |
| bess-ems | eigenständige Python-Linie (377 Z., Inventur-Nachtrag Lastenheft 0.4.0) | 0 (577 Links) | 1 → 0 | 1 echter Mehr-Befund (Self-Anker: `Daten- und` sluggt zu Doppel-Bindestrich, Alt kollabierte) | geschrumpft: **Host-Pfad-Rest-Sensor** (377 → ~130 Z.) | a6c8b9b, 1085912 |

Bekannte, dokumentierte Abdeckungs-Deltas (jeweils im Bestand
befundfrei, in den Ziel-Configs notiert): Reference-Style-Definitionen,
HTML-/`{#id}`-Anker, ASCII-Slug-Varianten (bess-ems); vertagte
Kennungs-Auflösungen als künftige `ids`-Konfiguration (c-hsm-doc
`HSM-*`, grid-gym `GG-*`/`AC-*`).

### Zusatz: Neu-Adoption pkcs11-course (kein 14. Vorkommen)

Fund Auftraggeber: Repo ohne jeden Doku-Check → Adoption statt
Migration (`.d-check.yml` + `make docs-check`, 83 Dateien, 1
Bestands-Befund gefixt). Der Adoptionslauf fand eine
Robustheits-Lücke, die kein kuratiertes Alt-Tool-Repo zeigen konnte:
unlesbare Laufzeit-Residuen (root-eigene `.gradle/`-,
SoftHSM-Token-Verzeichnisse) brachen den Scan ab, und `scan.ignore`
bot keinen Ausweg. Fix in **v0.2.1**: Ignore-Muster prunen den
Verzeichnis-Abstieg, `.gradle` in den `SKIP_DIRS`
(`DC-FA-SCAN-001`-Spec fortgeschrieben, 43d1c38).

### Lerneintrag (Steering Loop)

1. **Die CommonMark-Span-Blindheit war keine u-boot-Singularität:**
   grid-gym versteckte dahinter einen echten toten Link, b-cads
   vendorte Kopie prüfte tote Präfixe, bess-ems' Slug-Kollabierung
   übersah einen kaputten Anker. Konsolidierung ersetzt nicht nur
   Drift — sie deckt die akkumulierten *Folgen* der Drift auf
   (16 echte Mehr-Befunde in 9 Repos bei durchgängig „grünen" CIs).
2. **Adoptionen sind eine eigene Sensor-Quelle:** Der erste Lauf in
   einem ungepflegten Korpus (pkcs11-course) fand eine
   Robustheits-Klasse (unlesbare Residuen), die in den 13 kuratierten
   Migrations-Repos nie auftrat → v0.2.1. Vergleichsläufe testen
   Parität, Adoptionen testen Robustheit.
3. **Vendoring ohne lebende Quelle verwaist still:** b-cads
   „unverändert übernommene" Kopie hatte ihre Drift-Nachzieh-Pflicht
   (deren `MR-003`) faktisch verloren, als die Quelle migrierte. Der
   Digest-Pin ist das bessere Vendoring — Reproduzierbarkeit ohne
   Pflege-Pflicht.

## 8. Sub-Area-Modus-Begründung

Dieses Repo: GF (Doku-/Beleg-Arbeit). Ziel-Repos: außerhalb des
Geltungsbereichs dieses Harness (deren Modus gilt dort).
