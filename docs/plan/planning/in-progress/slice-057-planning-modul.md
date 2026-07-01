# Slice slice-057: Modul `planning` — Planning-Lifecycle-Konsistenz (hermetisch)

**Status:** in-progress (welle-46-planning-modul).

**Welle:** welle-46-planning-modul (Trigger: Auftraggeber-`tools/*.sh`-Audit — der
**letzte** Kandidat, `tools/planning-consistency.sh`; anders als die drei Vorgänger
rein **hermetisch**, ohne git).

**Bezug:** Führt eine **neue Anforderung** im Lastenheft ein (Modul `planning`,
`planning-drift`) plus [ADR-0028](../../adr/0028-planning-lifecycle-modul.md) (erste
Planning-Gate-ADR — das slice-040-Skript war ADR-los). Dieselbe Skript-Ablösungs-Linie
wie [`slice-053`](../done/slice-053-vcs-modul.md)/[`slice-055`](../done/slice-055-completeness-rueckbau.md)/[`slice-056`](../done/slice-056-commits-modul.md),
aber **ohne** VCS-Port (nur Filesystem). Verteilung wie der Rest des Werkzeugs
(gepinntes Image, kein kopiertes Skript —
[`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Linie).

**Autor:** pt9912. **Datum:** 2026-07-01.

---

## 1. Ziel

`tools/planning-consistency.sh` ([slice-040](../done/slice-040-planning-consistency-gate.md))
erzwingt die Planning-Lifecycle-Invariante („Roadmap §Aktuelle Welle ↔ `in-progress/slice-*`",
[`AGENTS.md` §3.3](../../../../AGENTS.md#33-git-mv--inhaltsänderung--zwei-commits)) — volle
Prüfung, aber als **Skript** nur per Kopie in Schwester-Repos nutzbar (Copy-Drift,
[`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse).
Es ist das **letzte** unmechanisierte Familien-Skript (nach `adr-immutable-check.sh` → `vcs`,
`completeness-check.sh` → in-Produkt-Flag, `trace-check.sh` → `commits`). **Neu:** dieselbe
Prüfung **verteilbar im Image** — ein opt-in Modul `planning` prüft `hasActive == hasSlices`,
sonst `planning-drift`. **Hermetisch** (nur Roadmap-Datei + Verzeichnis-Listing, **kein**
git, **kein** Netz, read-only) — d-check-nativ wie `codepaths`. d-check **dogfooded** das
Modul für seine eigene Roadmap (`make planning-check` läuft darauf um).

## 2. Entscheidungen

- **Hermetisches Modul `planning`** (15.), **kein** VCS-Port/git — der saubere Unterschied
  zu `vcs`/`commits`: nur der Filesystem-Port (`ReadFile` + Verzeichnis-`List`). Relativiert
  weder [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) noch
  [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).
- **Post-Pass, kein Datei-Scan.** Es prüft **eine** deklarierte Roadmap + ihr
  Verzeichnis-Listing (nicht den Markdown-Baum) — ein Post-Pass in `run.go` wie `vcs`/`commits`,
  aber mit `fsys` statt VCS-Port.
- **Config, damit verteilbar.** `planning.roadmap` (Pflicht zum Aktivieren, leer ⇒ inert),
  `planning.heading`/`marker`/`slice-glob` mit Konventions-Defaults (`## Aktuelle Welle` /
  „Keine aktive Welle" / `slice-*.md`).
- **Heading-Guard fail-closed** (R1-MEDIUM-Lehre von slice-040): fehlt die kanonische
  Überschrift exakt (bzw. die Roadmap-Datei), ⇒ `planning-drift` (kein stilles „aktiv").
- **Bindepunkt unverändert.** `make planning-check` bleibt Meta-Gate **in `make gates`**
  (Arbeitsbaum-Bindepunkt) — anders als `adr-check`/`trace-check`. Nur die Mechanik (Skript
  → Modul) wechselt.
- **Dogfood-Ersatz + Skript entfernt.** `make planning-check` ruft das Image mit `--enable
  planning` (`FOCUS_DISABLE`); `planning:`-Block in [`.d-check.yml`](../../../../.d-check.yml),
  **außerhalb** der Default-`modules` (default-aus byte-identisch). `tools/planning-consistency.sh`
  wird **atomar** per `git rm` entfernt, die immutable slice-040-Inline-Referenz als
  `codepaths.ignore-refs`-Tombstone ([ADR-0025](../../adr/0025-codepaths-ignore-refs.md), vierter Fall).
- **Config-Surface + Verteilung.** `--print-mk` trägt `doc-planning` (`--enable planning`,
  hermetisch ohne `--range`; [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  8→9); `--print-config`/`--suggest-config` führen `planning`.

## 3. Definition of Done

- [x] **Spec/Doc (doc-first):** neue Anforderung
  [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (Bereich `PLAN` in §3, Versions-Bump 0.36.0 + §7-Historie, `planning` in
  [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) + Glossar)
  + [ADR-0028](../../adr/0028-planning-lifecycle-modul.md) (Proposed, + Index) + spezifikation
  `.a` + Grund-Code `planning-drift` (§4) + Schema-Keys `planning.*` (§2) + [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)-Körper 8→9.
- [ ] **Code:** `rules/planning.go` (`CheckPlanning(fsys, cfg)`: Roadmap lesen, Heading-Block,
  Marker, in-progress-Listing, `hasActive==hasSlices`, Heading-Guard fail-closed, `planning-drift`),
  `model.PlanningConfig` + `validModules()` + `finding.go`-Reason, `configyaml` (rawPlanning/applyPlanning/scope),
  `run.go`-Post-Pass (hermetisch, `fsys`). Tests (die 3 Skript-Richtungen + 2 Gegenproben +
  Heading-Guard via Mem-FS); Modul-aus byte-identisch.
- [ ] **Gate-Umbau:** `make planning-check` auf `--enable planning` (Image, `FOCUS_DISABLE`);
  `planning:`-Block in `.d-check.yml`; `tools/planning-consistency.sh` per `git rm` + Tombstone;
  `harness/README.md` §Sensors + `AGENTS.md` §4 nachziehen; `--print-mk doc-planning` +
  print-config/suggest; `gate-consistency` grün.
- [ ] `make ci`/`make fullbuild` grün; `make planning-check` (Dogfood); unabhängige Reviews;
  CHANGELOG; release-prep (v0.36.0); Closure (Move nach `done/` + Roadmap-Flip,
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise));
  bei Closure [ADR-0028](../../adr/0028-planning-lifecycle-modul.md) → Accepted; Release v0.36.0.

## 4. Risiken / offene Punkte

- **`make planning-check` braucht jetzt das Image** (statt Skript) — schwergewichtiger pro Lauf;
  `make gates` baut das Image ohnehin (`doc-check`). Bootstrap: der neue Gate greift erst grün,
  wenn das Image das `planning`-Modul trägt; Umstellung + Skript-Löschung atomar mit dem Modul.
- **Verlust des Per-Lauf-Selbsttests:** das Skript testete sich bei jedem Aufruf (3 Richtungen);
  der Selbsttest wandert in `make test` — derselbe Trade wie
  [ADR-0027](../../adr/0027-commits-traceability-modul.md).
- **Layout-spezifischer Verteilungswert:** bewusst nachrangig (die „Keine aktive Welle"-Konvention
  ist harness-spezifisch) — die Config-Keys machen es für Adopter parametrierbar.
- **Dogfood-Selbstbezug:** `planning` prüft d-checks eigene Roadmap; solange slice-057 in
  `in-progress/` liegt, MUSS die Roadmap eine aktive Welle benennen (bei Closure Flip auf idle).

## 5. Trigger

Auftraggeber-`tools/*.sh`-Audit (2026-06-29, Roadmap welle-46): `planning-consistency.sh` ist
das letzte Familien-Skript. Anders als die drei Vorgänger hermetisch (Markdown↔Dateisystem,
ohne git) — d-check-fit. Nachrangig priorisiert (kleinerer Verteilungswert), jetzt umgesetzt
auf Auftraggeber-Wunsch.

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code + Spec; „Doc führt, Code folgt"). Das `planning`-Modul ist eine
GF-Erweiterung der Kern-Regel-Schicht (nutzt den bestehenden Filesystem-Port); keine
BF-Sub-Area, kein neuer Adapter.

## 7. Closure-Notiz (nach `done/`)

_(folgt bei Closure — Umsetzung, Dogfood-Belege, Gate-Ausgaben, Reviews, Release.)_
