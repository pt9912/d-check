# Slice slice-058: arch-check via Schwester-Tool a-check

**Status:** in-progress (welle-47-arch-check-a-check).

**Welle:** welle-47-arch-check-a-check (Trigger: der Roadmap-§Nächste-Wellen-Zeiger
„`arch-check` (Go-Importe → Schwester-Projekt a-check)" ist einlösbar — a-check ist
released und liefert Image + `a-check.mk` + `.a-check.yml`).

**Bezug:** [ADR-0029](../../adr/0029-arch-check-via-a-check.md) (Supersedes die
Fitness-Function-Mechanik von [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md)/[ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md);
Regeln R1–R6 bleiben). **Bewusst kein Lastenheft-CR:** das d-check-Produkt ändert sich
nicht — der [Anforderungs-Anlege-Prozess](../../../../harness/conventions.md#anforderungs-anlege-prozess)
greift nur für `DC-*`-Anforderungen, Gate-Mechanik ist ADR-Domäne (Präzedenz
[slice-039](../done/slice-039-pr-ci-traceability-gate.md)/[slice-055](../done/slice-055-completeness-rueckbau.md));
die [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Bindung
des Gates bleibt. Dieselbe Skript-Ablösungs-Linie wie
[`slice-053`](../done/slice-053-vcs-modul.md)/[`slice-055`](../done/slice-055-completeness-rueckbau.md)/[`slice-056`](../done/slice-056-commits-modul.md)/[`slice-057`](../done/slice-057-planning-modul.md),
aber erstmals durch das **Schwester-Tool** statt durch ein d-check-Modul
([`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Linie:
verteilen statt kopieren).

**Autor:** pt9912. **Datum:** 2026-07-03.

---

## 1. Ziel

`tools/arch-check.sh` (Dockerfile-Stage, `go list`) erzwingt die Import-Regeln R1–R5
([ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md)) + R6
([ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)) — als handgeschriebenes,
Go-spezifisches Skript, dessen Familie in vier Schwester-Repos divergent kopiert ist
(Copy-Drift-Klasse). **Neu:** `make arch-check` konsumiert das
[a-check](https://github.com/pt9912/a-check)-Image (digest-gepinnt, netzlos, read-only)
über ein include-bares `a-check.mk` plus eine repo-eigene `.a-check.yml`, die R1–R6
deklarativ ausdrückt. Das Skript und die Dockerfile-Stage entfallen; die
Schwester-Beziehung wird symmetrisch (a-check konsumiert d-check bereits via `d-check.mk`).

## 2. Entscheidungen

Plan-Review R1 ([Report](../../../reviews/2026-07-03-slice-058-arch-check-plan-r1.md),
Verdikt NACHBESSERN) ist eingearbeitet; die `(R1-…)`-Anker unten benennen das jeweils
adressierte Finding.

- **Konsum-Form:** `a-check.mk` (aus `a-check --print-mk`, an die Repo-Politik angepasst)
  wird ins Makefile eingebunden; das Image ist `@sha256:`-digest-gepinnt
  ([ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)-Politik, aktueller Stand
  v0.6.0); der Lauf ist `docker run --network none` mit read-only-Mount —
  [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-konform
  wie `make semgrep` ([ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md)).
- **gate-consistency-Verdrahtung (R1-MEDIUM-3):** `tools/gate-consistency.sh` parst in
  beiden Richtungen **nur die Datei `Makefile`** (keine includes), und das generierte
  Fragment-Target heißt `a-check`, nicht `arch-check`. Daher: das Gate-Target
  **`arch-check` bleibt im `Makefile` selbst definiert** und delegiert an das
  Fragment-Target `a-check` (das Fragment liefert Pin + Basis-Target unverändert).
  Die Delegation ist zugleich das Ventil gegen Fragment-Divergenz: bei Pin-Hebung wird
  `a-check.mk` per `--print-mk` neu erzeugt, ohne das `arch-check`-Target anzufassen.
- **Regel-Übersetzung (R1-HIGH-1/LOW-2 — vollständig gegen das Skript, nicht gegen die
  Regel-Kurzform):** `.a-check.yml` bildet **alle Verbotszweige** von
  `tools/arch-check.sh` ab — `layers` für Kern/Ports/Adapter **und** die Kern-Pakete
  (`model`/`rules`/`app` **plus** Test-Helfer `coretest`; `model`/`rules` brauchen
  **explizite** `role: domain`, die Namens-Inferenz kennt sie nicht), `edges` als
  Richtungs-Allowlist (R6 + Hexagon-Richtung), `tech` für **beide** R2-Kapseln
  (`net/http` → httpcheck **und** go-git → git-Adapter), R3 (`yaml` → configyaml
  **und** report, [ADR-0009](../../adr/0009-yaml-im-report-adapter.md)), R4
  (`os` → fs-Adapter) **und** die restliche R1-Bannliste (verbotene `net`-Familie
  **enumeriert ohne `net/url`** — die Rollen-Reinheit verbietet im Kern jedes
  tech-Muster kategorisch, ein „Pass-Eintrag" für `net/url` ist mechanisch
  unmöglich (R2-proben-belegt), die Ausnahme geht nur über Nicht-Nennung;
  dazu `syscall`, `io/fs`, `os/`); alle `adapter`-Werte enden auf `/`
  (Substring-Match — sonst erbt ein Präfix-Vetter wie `gitlab` die
  `git`-Kapsel, R2-MEDIUM-2); `composition_root` für CLI + `cmd`
  (die R4-Ausnahme-Zone).
- **Paritäts-Beleg (Pflicht, vor der Umstellung; R1-MEDIUM-2):** adversariale
  **Proben-Matrix je Skript-Verbotszweig** — R1, R2a (`net/http`), R2b (go-git), R3,
  R4, R5, R6: injizierter verbotener Import ⇒ `make arch-check` rot, Revert ⇒ grün —
  **plus Allow-Gegenproben**: `net/url` im Kern und `yaml` im report-Adapter dürfen
  **nicht** flaggen; die R1-Probe liegt in `core/model` (verriegelt die explizite
  Rollen-Zuordnung mit). Zählung je Verbotszweig, nicht je Regel-Nummer
  (slice-057-R3-Lehre: sonst bleibt eine Kapsel unverriegelt). Rest-Deltas der
  Text-Heuristik gegenüber `go list` werden dokumentiert; nicht per Config schließbare
  Deltas werden Change Request an das a-check-Lastenheft (Schwester-Repo) — keine
  stille Lockerung
  ([`AGENTS.md` §3.6](../../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)).
- **Umstellungs-Vorbedingung (R1-HIGH-1/HIGH-2/MEDIUM-1 — drei heute quellen-belegte
  a-check-v0.6.0-Deltas ⇒ CR ans Schwester-Repo, Umstellung wartet):**
  (a) `tech` bindet ein Pattern an **einen** Adapter-Pfad (Erst-Treffer) — die
  R3-Zwei-Adapter-Erlaubnis ist nicht ausdrückbar; (b) der Scanner erfasst
  **`*_test.go`** (Glob-Engine ohne Negation), das Skript prüfte nur
  Nicht-Test-Imports — schon `os` in einem Adapter-Test macht den sauberen Baum rot;
  (c) `composition_root` ist eine **Total-Ausnahme** (auch von `tech`) — auf CLI/`cmd`
  ginge die heutige R2-/R3-Deckung verloren; wird (c) nicht per a-check-Feature
  schließbar, wird der Deckungsverlust als Rest-Delta **explizit gelistet**, nicht
  verschwiegen.
- **Rückbau atomar:** `tools/arch-check.sh` per `git rm` + Dockerfile-Stage `arch-check`
  (samt `--no-cache-filter`-Variable) entfernt; die immutablen Inline-Referenzen
  ([ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md) §Fitness Function,
  [ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)) werden der **fünfte**
  `codepaths.ignore-refs`-Tombstone ([ADR-0025](../../adr/0025-codepaths-ignore-refs.md)).
- **Bindepunkt unverändert:** `arch-check` bleibt Produkt-Gate in `make gates`/`ci`;
  Target-Name bleibt (`gate-consistency` hält Doku ↔ Makefile).
- **Pin-Pflege:** der a-check-Pin erscheint in `make versions`; Hebung = bewusster Commit.
- **Kein d-check-Produkt-Code, kein Release** — das d-check-Image bleibt byte-identisch
  (wie [`slice-055`](../done/slice-055-completeness-rueckbau.md)).

## 3. Definition of Done

- [x] **Doc-first:** [ADR-0029](../../adr/0029-arch-check-via-a-check.md) (Proposed, +
  Index), Roadmap welle-47 aktiv; **kein** Lastenheft-CR (Begründung: §Bezug).
- [ ] **Vorbedingung erfüllt:** die drei v0.6.0-Deltas (§2) per CR im a-check-Repo
  adressiert und in einem Release geliefert; a-check-Pin auf dieses Release gehoben —
  **vor** dem Makefile-Umbau.
- [ ] **Config:** `.a-check.yml` (alle Skript-Verbotszweige als
  `layers`/`edges`/`tech`/`composition_root`, inkl. `role: domain` für
  `model`/`rules`, `coretest`-Zuordnung, beide R2-Kapseln, R3-Doppel-Erlaubnis,
  R1-Bannliste) + `a-check.mk` (digest-gepinnt, `--network none`, read-only-Mount,
  `##`-Help-Annotation).
- [ ] **Makefile/Dockerfile:** `arch-check`-Target bleibt im `Makefile` definiert und
  delegiert an das Fragment-Target `a-check` (gate-consistency parst nur `Makefile`);
  die Dockerfile-Stage `arch-check` + `NO_CACHE_FILTER_ARCH` + die
  `$(IMAGE):arch-check`-Zeile im `clean`-Target + der Makefile-/Dockerfile-Kopfkommentar
  bereinigt; `make versions` weist den a-check-Pin aus.
- [ ] **Rückbau:** `tools/arch-check.sh` per `git rm`; fünfter
  `codepaths.ignore-refs`-Tombstone in `.d-check.yml`; `make doc-check` bleibt grün.
- [ ] **Paritäts-Beleg:** Proben-Matrix je Verbotszweig (R1 in `core/model`, R2a, R2b,
  R3, R4, R5, R6: Verstoß ⇒ rot, Revert ⇒ grün) **plus** Allow-Gegenproben (`net/url`
  im Kern, `yaml` im report-Adapter: kein Befund), Ausgaben im Review-Report/der
  Closure-Notiz; Rest-Deltas explizit gelistet (insb. `composition_root`-Deckung auf
  CLI/`cmd`, falls ungelöst).
- [ ] **Doku-Currency:** [`harness/README.md`](../../../../harness/README.md) §Sensors
  (arch-check-Zeile: Mechanik + Bindung) + [`AGENTS.md`](../../../../AGENTS.md) §4 +
  Makefile-Kopfkommentar; `make gate-consistency` grün.
- [ ] **Gates + Review:** `make gates` und `make ci` grün; mindestens ein unabhängiger
  Review (R1) vor Closure; Closure-Move nach `done/` + Roadmap-Flip
  ([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise));
  bei Closure [ADR-0029](../../adr/0029-arch-check-via-a-check.md) → Accepted +
  Geschichte-Anhänge/Index-Annotation an
  [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md)/[ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)
  (teil-superseded).

## 4. Risiken / offene Punkte

- **Präzisions-Trade `go list` → Text-Heuristik:** a-check extrahiert Importe
  text-heuristisch (dokumentierte Grenze, dort `AC-QA-02`); Build-Tags/Transitive sieht
  nur `go list`. Nach Quellen-Lage bereits geklärt (R1-Review): (a) `net/url` flaggt
  **nicht** (Kern-Reinheit trifft nur Layer-/`tech`-Auflösungen) — Bedingung ist
  Pattern-Disziplin in der `.a-check.yml` (verbotene `net`-Familie enumerieren, RE2
  ohne Lookahead); (b) die R4-Zone läuft über `composition_root`, kostet dort aber die
  R2-/R3-Deckung (Total-Ausnahme). Die drei harten v0.6.0-Deltas stehen in §2 als
  CR-Vorbedingung — Umstellung wartet auf das liefernde a-check-Release.
- **Über-Deckungs-Richtung (R2-INFO-1, benannt):** die `edges`-Allowlist und die
  tech-Bindungen wirken auch **außerhalb** des Kerns (das Skript regulierte dort
  nicht) — strenger als R1–R6, keine Lockerung; ein legitimer neuer Querbezug
  braucht künftig eine bewusste `edges`-/`tech`-Erweiterung statt stillen
  Durchrutschens. Rest-Delta bleibt: **neue** `net/*`-stdlib-Subpakete brauchen
  einen Enumerations-Nachtrag in der `.a-check.yml`.
- **Externe Release-Abhängigkeit:** das Gate hängt am a-check-Release-Stand (wie semgrep
  am gepinnten Regelset) — Digest-Pin macht es reproduzierbar, Pin-Hebung bewusst.
- **Bootstrap-Reihenfolge:** Umstellung + Skript-Löschung + Tombstone müssen in einem
  Commit landen, sonst ist der Zwischenstand gate-rot (`codepath-missing` bzw. doppelte
  Mechanik).
- **Kein Dogfood-Selbstbezug:** anders als die Modul-Ablösungen prüft hier kein
  d-check-Feature — der Slice liefert bewusst kein d-check-Release.

## 5. Trigger

Nutzer-Frage 2026-07-03 („Verwenden wir schon das Tool a-check?") + Auftrag, die Adoption
regelkonform aufzusetzen. Sachlage: a-check v0.6.0 ist released (GHCR, sieben
Sprach-Backends, `--print-mk`-Konsum-Fragment); der Roadmap-§Nächste-Wellen-Zeiger
„arch-check → Schwester-Projekt a-check" (seit dem `tools/*.sh`-Audit 2026-06-29) wird
damit einlösbar.

## 6. Sub-Area-Modus-Begründung

GF („Doc führt, Code folgt": ADR vor Umbau). Kein d-check-Produkt-Code — Gate-/
Harness-Infrastruktur wie [slice-039](../done/slice-039-pr-ci-traceability-gate.md)/[slice-040](../done/slice-040-planning-consistency-gate.md)/[slice-055](../done/slice-055-completeness-rueckbau.md);
keine BF-Sub-Area. Die konsumierte Prüf-Logik lebt im Schwester-Repo (a-check, eigener
Harness); hier entstehen nur Config (`.a-check.yml`), Makefile-Verdrahtung und Rückbau.

## 7. Closure-Notiz (nach `done/`)

_(folgt bei Closure — Umsetzung, Paritäts-/Mutations-Belege je R1–R6, Gate-Ausgaben,
Review, Rest-Deltas.)_
