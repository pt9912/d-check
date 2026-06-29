# Slice slice-055: `make completeness-check` dogfoodet den in-Produkt-Flag — `completeness-check.sh` entfernen

**Status:** done (welle-44-completeness-rückbau).

**Welle:** welle-44-completeness-rückbau (Trigger: Auftraggeber-Audit 2026-06-29 „welche
`tools/*.sh` lassen sich noch in d-check mechanisieren?" — erster, kleinster Hebel).

**Bezug:** Verdrahtet das Closure-Gate `make completeness-check` auf den **bereits
existierenden** in-Produkt-Flag `--trace --require-complete`
([`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code),
slice-044) statt auf `tools/completeness-check.sh`; die RTM-Quelle bleibt
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix).
Eine **neue ADR** löst die „Skript-als-Gate-Quelle"-Teilentscheidung von
[ADR-0017](../../adr/0017-requirements-completeness-gate.md) ab (`Supersedes`; Policy und
Bindepunkt bleiben). Erster Folge-Anwendungsfall des Tombstone-Ventils
[ADR-0025](../../adr/0025-codepaths-ignore-refs.md) (slice-054).

**Autor:** pt9912. **Datum:** 2026-06-29.

---

## 1. Ziel

**Die Lücke:** [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code) (`--trace --require-complete`) wurde in slice-044 geschaffen,
**„damit Konsumenten die `completeness-check.sh`-Parsing-Logik nicht kopieren müssen"**
(Lastenheft §7, 0.24.0). d-check **selbst** aber führt sein Closure-Gate
`make completeness-check` weiter über `bash tools/completeness-check.sh` (ein Skript, das
`d-check --trace --json` aufruft und `orphans` in Bash parst). d-check empfiehlt
Konsumenten den verteilten Flag, kopiert für sich aber weiter die Skript-Logik — genau der
[`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Copy-Drift,
den die Werkzeug-Familie auflösen will. (Es existiert bereits `make doc-complete` =
`--trace --require-complete`, aber nur als „Dogfood, **kein** Gate-Bindepunkt" — die
Closure-Wahrheit hängt noch am Skript.)

**Die Lösung:** `make completeness-check` ruft denselben in-Produkt-Flag wie der
Dogfood/die Konsumenten — `d-check --trace --require-complete` (Image, netzlos,
read-only). Das Skript `tools/completeness-check.sh` wird **entfernt**. Damit isst d-check
für seine eigene Vollständigkeits-Invariante sein eigenes Futter; **eine** Wahrheit
(der Flag), keine Skript-Kopie.

## 2. Entscheidungen

- **`make completeness-check`-Recipe → Image-Flag.** Statt `bash tools/completeness-check.sh`
  ruft das Target `docker run --rm --network none -v <repo>:/repo:ro <image> --trace
  --require-complete` (hermetisch, read-only — wie der bestehende `doc-complete`-Dogfood).
  Bindepunkt unverändert: in `make fullbuild`, **nicht** in `gates`/`ci`
  ([ADR-0017](../../adr/0017-requirements-completeness-gate.md)-Policy bleibt).
- **Skript entfernt.** `git rm tools/completeness-check.sh`. Sein Negativ-Selbsttest
  (orphans-Parsing-Vektoren) ist durch die **akzeptanz-getestete** in-Produkt-Logik
  ([`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code),
  bestehende CLI-Tests) abgelöst — kein Garantieverlust.
- **Waisen-Sichtbarkeit gewahrt (kein Produkt-Change).** Der Flag meldet die **Anzahl**
  (stderr) und Exit 1; die **einzelnen** Waisen-IDs sind in der `--trace`-Tabelle als
  Status **`WAISE`** sichtbar. Das Skript listete sie konzis; die Tabelle zeigt sie mit
  Kontext. Parität ausreichend — **keine** Erweiterung von [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code) nötig.
- **Tombstone (`codepaths.ignore-refs`).** Die immutable
  [ADR-0017](../../adr/0017-requirements-completeness-gate.md)-Geschichte referenziert
  `tools/completeness-check.sh` in Inline-Code → nach `git rm` ein `codepath-missing` an
  uneditierbarer Doku. Eintrag in `.d-check.yml` `codepaths.ignore-refs` löst es —
  **zweiter realer Anwendungsfall** des in slice-054 gebauten Ventils
  ([ADR-0025](../../adr/0025-codepaths-ignore-refs.md)).
- **Neue ADR, kein neues DC.** Die Mechanik-Ablösung (Skript → in-Produkt-Flag) ist eine
  Entscheidung über [ADR-0017](../../adr/0017-requirements-completeness-gate.md)s
  „eine-Skript-Quelle"-Teil → neue ADR mit `Supersedes` (Policy/Bindepunkt unverändert,
  wie [ADR-0024](../../adr/0024-vcs-immutable-gate.md) die Skript-Mechanik von
  [ADR-0016](../../adr/0016-adr-immutable-gate.md) ablöste). **Kein** neues `DC-*` (die Mechanik [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)/[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix) existiert),
  **kein** Lastenheft-/Spezifikations-Eingriff.
- **Kein Release.** Es ändert sich **kein** Produkt-Code (`internal/`/`cmd/` unberührt) →
  das GHCR-Image ist byte-identisch zu v0.34.0. Reiner Harness-Refactor: **kein
  Versions-Bump, kein GHCR-Release, kein digest-backfill**; Commit-Kette
  feat → review → closure-move → closure-body.
- **`doc-complete` unberührt.** Der Konsumenten-Name im `--print-mk`-Fragment bleibt; nur
  d-checks **eigenes** Closure-Gate zieht nach.

## 3. Definition of Done

### 3a. Artefakte (Doc-first → Code → Entfernung)

- [ ] **Doc-first:** neue ADR (`Supersedes` die Skript-Gate-Teilentscheidung von
  [ADR-0017](../../adr/0017-requirements-completeness-gate.md); Bezug
  [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)/[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix),
  [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md),
  [`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding))
  + ADR-Index — inkl. **Korrektur** der [ADR-0025](../../adr/0025-codepaths-ignore-refs.md)-Index-Zeile
  auf `Accepted` (die Datei war seit der slice-054-Closure Accepted, der Index hing auf
  Proposed nach; R1-F-2).
- [ ] **Code:** `Makefile` — `completeness-check`-Recipe auf den Image-Flag; `git rm
  tools/completeness-check.sh`; `.d-check.yml` `codepaths.ignore-refs` +=
  `tools/completeness-check.sh`.
- [ ] **Gate-Doku nachziehen:** [`AGENTS.md`](../../../../AGENTS.md) §4 und
  [`harness/README.md`](../../../../harness/README.md) §Sensors — die
  `completeness-check`-Zeilen von „`tools/completeness-check.sh`/Bash-Parsing" auf
  „in-Produkt `--trace --require-complete`" umschreiben (Doku ↔ Makefile via
  `make gate-consistency`).

### 3b. Verifikation (Korrektheit)

- [ ] `make completeness-check` grün **über den Flag** (0 Waisen, Exit 0) und
  `make gate-consistency` grün (Target existiert + in beiden Doku-Stellen gelistet).
- [ ] `make ci` grün (`doc-check` u. a.) + `make fullbuild` läuft bis `completeness-check`
  grün durch.
- [ ] **Adversariale Waise:** eine künstlich verwaiste Anforderung lässt
  `make completeness-check` **rot** (Exit 1) mit Anzahl-Meldung **und** sichtbarer
  `WAISE`-Zeile in der `--trace`-Tabelle — danach revertiert.
- [ ] Zwei **unabhängige** Reviews (R1 Doc/Harness, R2 Mechanik/Gate) — HIGH/MEDIUM behoben.

### 3c. Validierung (Wirksamkeit)

- [ ] **Copy-Drift geschlossen:** `tools/completeness-check.sh` ist entfernt; d-checks
  Closure-Gate und die Konsumenten-Verteilung nutzen **dieselbe** in-Produkt-Mechanik.
- [ ] **Tombstone-Ventil zum zweiten Mal wirksam:** `make doc-check` bleibt grün, obwohl
  die immutable [ADR-0017](../../adr/0017-requirements-completeness-gate.md)-Referenz auf
  den entfernten Pfad zeigt (Beleg, dass `codepaths.ignore-refs` generisch trägt).
- [ ] **Keine Regression:** Image byte-identisch (kein Produkt-Code), Default-Lauf
  unverändert.

### 3d. Closure

- [ ] Move nach `done/` + Roadmap-Flip
  ([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise));
  die neue ADR → Accepted + ein `## Geschichte`-Append an die immutable
  [ADR-0017](../../adr/0017-requirements-completeness-gate.md) (Skript abgelöst — erlaubter
  Anhang, Core unangetastet) **+ die [ADR-0017](../../adr/0017-requirements-completeness-gate.md)-Index-Zeile
  mit der Teil-Supersede-Notiz auf die neue ADR annotieren** (Form wie
  [ADR-0016](../../adr/0016-adr-immutable-gate.md)/[ADR-0024](../../adr/0024-vcs-immutable-gate.md)
  sie im Index führen; R1-F-1 — sonst bleibt der Index-Status inkonsistent); welle-44 aus
  §Nächste Wellen entfernt. **Kein Release.**

## 4. Risiken / offene Punkte

- **Waisen-Ausgabe weniger konzis:** die `--trace`-Tabelle ist länger als die
  Skript-Waisen-Liste. Akzeptabel — die `WAISE`-Zeilen + die Anzahl-Meldung benennen die
  Waisen eindeutig, mit mehr Kontext. Sollte sich Bedarf für eine knappe Nur-Waisen-Ausgabe
  zeigen, wäre das eine [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)-Erweiterung (eigener Trigger), kein Blocker hier.
- **Schwergewichtiger pro Lauf:** der Flag baut/ruft das Image (wie `adr-check` seit
  slice-053). Akzeptabel — `completeness-check` ist ein Closure-/`fullbuild`-Gate, kein
  Pro-Commit-Lauf; die CI baut das Image ohnehin.
- **[ADR-0017](../../adr/0017-requirements-completeness-gate.md)-Teilentscheidung umkehren:**
  braucht regelkonform eine neue ADR (`Supersedes`); Policy (Waise → FAIL) und Bindepunkt
  (`fullbuild`, nicht `gates`/`ci`) bleiben unberührt.
- **Index-Status-Honesty schlüpft an der Closure durch (R1-F-1/F-2 — Steering):** slice-055
  korrigiert nebenbei die [ADR-0025](../../adr/0025-codepaths-ignore-refs.md)-Index-Zeile
  (Proposed→Accepted, ein slice-054-Closure-Miss; die Datei war schon Accepted) und muss bei
  der eigenen Closure die [ADR-0017](../../adr/0017-requirements-completeness-gate.md)-Index-
  Annotation mitziehen (§3d). Dieselbe wiederkehrende Klasse zum zweiten Mal → Kandidat für
  eine **Mechanisierung** (ein Check „ADR-Datei-Status ↔ Index-Zeilen-Status konsistent",
  selbst d-check-nah, da Doc-Konsistenz) statt manuellem Treffen pro Slice — als Folge-Idee
  notiert, nicht in diesem Slice gebaut.
- **Leerer Scan-Scope ⇒ silent-green (R2-INFO, vorbestehend):** findet `--trace` keine
  Anforderung (fehlende/leere `spec/lastenheft.md` o. Ä.), ist `Orphans==0` ⇒ Exit 0 — für
  d-checks eigenes Repo (35 Anforderungen) unkritisch. **Von slice-055 nicht eingeführt:** das
  abgelöste Skript hatte denselben Pfad (kein `total>0`-Boden bei `{"total":0,"orphans":0}`).
  Ein Vollständigkeits-Boden wäre eine
  [`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)-Erweiterung
  (Produkt-Change, eigener Trigger), bewusst nicht hier.

## 5. Trigger

Auftraggeber 2026-06-29 nach dem v0.34.0-Release: Audit „welche der `tools/*.sh` könnte man
noch ersetzen?" — `completeness-check` als kleinster Hebel (in-Produkt-Mechanik existiert
seit slice-044), vor `trace-check` (welle-45) und `planning-consistency` (welle-46).

## 6. Sub-Area-Modus-Begründung

BF (Build/Harness-Refactor: Makefile-Gate + Skript-Entfernung + Gate-Doku; kein Produkt-Code,
kein Spec-Eingriff). Doc-first bleibt für die ablösende ADR (Entscheidung führt).

## 7. Closure-Notiz (nach `done/`)

**Geliefert:** `make completeness-check` dogfoodet jetzt den in-Produkt-Flag `--trace
--require-complete` ([`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code))
über die geteilte Makefile-Variable `COMPLETE_FLAGS` (eine Quelle mit `doc-complete`); das
Skript `tools/completeness-check.sh` ist per `git rm` entfernt und als
`codepaths.ignore-refs`-Tombstone deklariert (zweiter realer Einsatz des slice-054-Ventils,
[ADR-0025](../../adr/0025-codepaths-ignore-refs.md)). Damit isst d-check für seine eigene
Vollständigkeits-Invariante sein verteiltes Futter — der
[`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Copy-Drift
am Closure-Gate ist geschlossen.

**Doc-first:** [ADR-0026](../../adr/0026-completeness-in-product-gate.md) (Accepted,
`Supersedes` die Skript-als-Gate-Quelle-Teilentscheidung von
[ADR-0017](../../adr/0017-requirements-completeness-gate.md) — Policy „Waise→FAIL" + Bindepunkt
„`fullbuild`, nicht `gates`/`ci`" bleiben); ADR-Index (+ Korrektur der
[ADR-0025](../../adr/0025-codepaths-ignore-refs.md)-Index-Zeile auf Accepted, ein
slice-054-Closure-Miss) + [ADR-0017](../../adr/0017-requirements-completeness-gate.md)-Index-Annotation + Geschichte-Append.

**Verifikation:** `make ci` + `make gate-consistency` grün; eine **adversariale Waise** (synthetische Test-Anforderung)
trieb `make completeness-check` rot — Anzahl-Meldung + sichtbare `WAISE`-Zeile in der
`--trace`-Tabelle, danach per Edit revertiert (lastenheft sauber). Zwei unabhängige Reviews
(R1 doc 0H/1M/1L, R2 mechanik 0H/0M/1L/1I/1 REFUTED — fail-closed-Regression mit Skript-Zitat
widerlegt), alle Befunde behoben; Reports unter `docs/reviews/2026-06-29-slice-055-completeness-doc-r1.md`
und `docs/reviews/2026-06-29-slice-055-completeness-mechanic-r2.md`.

**Validierung:** Copy-Drift geschlossen (eine Mechanik für Gate + Konsumenten); Tombstone-Ventil
zum **zweiten** Mal wirksam (doc-check grün trotz entferntem, von immutabler
[ADR-0017](../../adr/0017-requirements-completeness-gate.md) zitiertem Skript). **Kein Produkt-Code** (`internal/`/`cmd/` unberührt) → Image byte-identisch
zu v0.34.0, **kein Release** (kein Versions-Bump, kein GHCR, kein Tag).

**Bewusst offen (R2-INFO, vorbestehend):** leerer Scan-Scope ⇒ `Orphans==0` ⇒ Exit 0 — von
slice-055 nicht eingeführt (das Skript hatte denselben Pfad); ein `total>0`-Boden wäre eine
[`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)-Erweiterung.
**Steering (R1-F-2):** die wiederkehrende Index-Status-Honesty-Klasse (slice-054 verfehlte den
[ADR-0025](../../adr/0025-codepaths-ignore-refs.md)-Flip, slice-055 fängt ihn) ist Kandidat für eine künftige Mechanisierung (ADR-Datei-Status
↔ Index-Zeilen-Status).
