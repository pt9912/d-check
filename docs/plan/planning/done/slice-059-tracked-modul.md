# Slice slice-059: Modul `tracked` — Getrackt-Status auflösbarer Referenz-Ziele

**Status:** done (welle-48-tracked-modul, Closure 2026-07-03).

**Welle:** welle-48-tracked-modul (Trigger: Auftraggeber-Frage 2026-07-03 —
„Was passiert, wenn ein Dokument ein gitignoriertes Dokument referenziert?" —
plus Fixture-Demo: Erzeuger-Checkout grün, frischer Klon `target-missing`).

**Bezug:** Führt eine **neue Anforderung** im Lastenheft ein
([`DC-FA-TRK-001`](../../../../spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in),
Lastenheft 0.37.0; zugleich
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
9→10 Targets) plus [ADR-0030](../../adr/0030-tracked-referenz-ziele.md)
(Proposed). Dritte Nutzung des VCS-Ports aus
[`slice-053`](../done/slice-053-vcs-modul.md) (`vcs`: Range-Diff,
[`slice-056`](../done/slice-056-commits-modul.md)/`commits`: Messages,
`tracked`: **Index** — ohne Range); Kein-Doppelbefund-Prinzip wie
[`slice-049`](../done/slice-049-pins-modul.md)/`pins`.

**Autor:** pt9912. **Datum:** 2026-07-03.

---

## 1. Ziel

d-check ist im Default git-agnostisch: eine Referenz auf eine **nur lokal
existierende** Datei (untracked/gitignoriert) ist beim Erzeuger grün und auf
jedem frischen Klon `target-missing` — Umgebungs-Drift zwischen Arbeitsbäumen,
die heute erst die CI des nächsten Checkouts fängt (oder Einzelfall-Disziplin
wie das gitignore+`scan.ignore`-Doppel
([`MR-017`](../../../../harness/conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen))
und das Vendoring
([`MR-019`](../../../../harness/conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017))).
**Neu:** ein opt-in Modul `tracked` (16.) prüft die von `links` aufgelösten,
**existierenden** repo-internen Link-/Bild-Ziele gegen den **git-Index** und
meldet `target-untracked` — die Falle wird am Entstehungsort gefangen.

## 2. Entscheidungen

- **Index ist die Wahrheit, keine `.gitignore`-Interpretation:** kein zweiter
  Regel-Interpreter; frisch gestagte Dateien gelten als getrackt
  (WIP-tauglich — `git add` macht neue Doku grün).
- **VCS-Port, dritte Nutzung — ohne Range:** Port-Erweiterung um die
  Index-Abfrage (Präzedenz: `commits` erweiterte um Message-Lesen); Prüfung
  **je gescannter Quell-Datei** über dieselbe Auflösungs-Mechanik wie `links`
  (unabhängig von dessen Aktivierung; Index-Menge einmal je Lauf — R1-M1).
  Eingabe erweitert (`.git`), aber
  lokal/lesend/deterministisch
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)/[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  in der `vcs`-Lesart).
- **Kein Doppelbefund:** nur existierende, aufgelöste **Datei**-Ziele;
  `target-missing` bleibt `links` (`pins`-Prinzip); Verzeichnis-Ziele kein
  Kandidat (Index führt nur Dateien — R1-M3 dokumentiert), Symlink-Referenzen
  kategorisch `links`-Domäne (R2-M2/M3: false-positive hinter getrackten
  Verzeichnis-Symlinks vermieden, Skip beidseitig verriegelt); der Befund
  nennt den **aufgelösten** Zielpfad (Ventil-Parität — R2-M1).
- **Ventil `tracked.exempt-targets`** (Glob über den aufgelösten Zielpfad,
  referenz-weit analog `codepaths.ignore-refs`) für absichtlich untrackte
  Ziele; ohne Eintrag byte-identisch.
- **Fail-closed:** aktives `tracked` ohne lesbares `.git` ⇒ Exit 2 (wie
  `vcs`/`commits`); strikt opt-in, default-aus byte-identisch, diagnose-only.
- **Config-Surface + Verteilung:** `--print-mk` trägt `doc-tracked`
  (`--enable tracked`, fokussierte `--disable`-Liste aus `ValidModules`
  abgeleitet, ohne Range;
  [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  9→10); `--print-config`/`--suggest-config` führen `tracked` (opt-in-Hinweis).
- **Kein neuer Gate-Bindepunkt in d-check:** Beleg-Lauf gegen das eigene Repo
  bei Closure (grün + adversariale Untracked-Probe rot); kein neues
  Make-Target ohne Bedarf.

## 3. Definition of Done

- [x] **Spec/Doc (doc-first):** neue Anforderung
  [`DC-FA-TRK-001`](../../../../spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)
  (Bereich `TRK` in §3, Versions-Bump 0.37.0 + §7-Historie, `tracked` in
  [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
  + Glossar,
  [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)-Körper
  9→10) + [ADR-0030](../../adr/0030-tracked-referenz-ziele.md) (Proposed, +
  Index) + Spezifikation `.a` + Grund-Code `target-untracked` (§4) +
  Schema-Key `tracked.exempt-targets` (§2).
- [x] **Code:** VCS-Port um Index-Abfrage erweitern (getrackte Pfade;
  gestagte Dateien enthalten), `tracked`-Post-Pass in `run.go` über die
  aufgelösten `links`-Ziele (existierend ∧ ∉ Index ∧ kein
  `exempt-targets`-Match ⇒ `target-untracked`), `model`-Config +
  `validModules()` + Reason, `configyaml` (rawTracked/applyTracked, Glob
  config-zeitig validiert — fail-closed), fail-closed ohne `.git` (Exit 2).
- [x] **Tests:** die sieben Akzeptanzkriterien (Happy/Index-Wahrheit/
  Modul-aus-byte-identisch/Negative/Kein-Doppelbefund/Ventil/fail-closed)
  gegen git-Fixture-Repos; Guards mutations-verifiziert (slice-057-R3-Lehre:
  jede Probe verriegelt genau ihren Guard; Injektion/Aufbau je Probe
  verifiziert — slice-058-Harness-Lehre).
- [x] **Config-Surface:** `--print-mk doc-tracked`; `--print-config`/
  `--suggest-config` (opt-in-Kommentar); Benutzerhandbuch-Abschnitt;
  `FOCUS_DISABLE`-Kommentarlage im Makefile geprüft (`tracked` ist kein
  Default-Modul — Liste wächst nicht).
- [x] **Belege:** `make ci`/`make fullbuild` grün; Beleg-Lauf
  `--enable tracked` gegen das eigene Repo (grün) + adversariale Probe
  (temporär untracktes, verlinktes Ziel ⇒ rot, Revert ⇒ grün); unabhängige
  Reviews (mind. R1 doc + R2 code) vor Closure; CHANGELOG; release-prep
  (Minor v0.37.0); Closure (Move nach `done/` + Roadmap-Flip,
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise));
  bei Closure [ADR-0030](../../adr/0030-tracked-referenz-ziele.md) → Accepted;
  Release v0.37.0 + Digest-Backfill.

## 4. Risiken / offene Punkte

- **go-git-Index-Semantik:** die Port-Erweiterung muss den Index (staged
  inklusive) liefern, nicht nur HEAD-Tree — sonst flaggte das Modul frisch
  gestagte Dateien (AK „Index-Wahrheit" verriegelt das per Test).
- **Performance:** eine Index-Abfrage pro Lauf (Set-Aufbau einmalig), kein
  Per-Referenz-git-Zugriff —
  [`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance) unkritisch,
  im Bench beobachten.
- **Symlink-/Groß-Kleinschreibungs-Ränder:** Index-Pfade sind
  case-sensitiv/slash-normalisiert; die Kandidaten-Pfade kommen aus derselben
  Auflösung wie `links` — Normalisierungs-Parität testen.
- **Konsumenten ohne git-Repo** (Tarball-/Doku-Ordner-Läufe): `doc-tracked`
  wäre dort Exit 2 (fail-closed, korrekt) — Handbuch nennt die Voraussetzung
  explizit.

## 5. Trigger

Auftraggeber-Frage 2026-07-03 („Was passiert, wenn ein Dokument ein
gitignoriertes Dokument referenziert?"), beantwortet per Code-/Spec-Befund
(git-agnostisch by design) + Fixture-Demo (lokal grün / frischer Klon
`target-missing`); Auftrag, die Prüf-Lücke regelkonform als CR + Slice
aufzusetzen. Die Idee war zuvor als offene Notiz festgehalten
(target-untracked via VCS-Port).

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code + Spec; „Doc führt, Code folgt"). Das `tracked`-Modul ist
eine GF-Erweiterung der Kern-Regel-Schicht über den bestehenden VCS-Port
(Port-Erweiterung wie bei `commits`); keine BF-Sub-Area, kein neuer Adapter.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** slice-059 als **16. Regelmodul `tracked`**: der VCS-Port ist um die
Index-Abfrage `TrackedPaths()` erweitert (dritte Nutzung — `vcs`: Range-Diff,
`commits`: Messages, `tracked`: **Index**, ohne Range; gestagte, nie committete
Dateien enthalten = Index-Wahrheit, keine `.gitignore`-Interpretation); das Modul
prüft **je gescannter Quell-Datei** die aufgelösten, existierenden repo-internen
Link-/Bild-**Datei**-Ziele (dieselbe Auflösungs-Mechanik wie `links`, unabhängig
von dessen Aktivierung) und meldet `target-untracked` mit dem **aufgelösten**
Zielpfad. Kein Doppelbefund (`target-missing` bleibt `links`; Verzeichnis-Ziele
kein Kandidat, Symlink-Referenzen kategorisch `links`-Domäne), Ventil
`tracked.exempt-targets` (referenz-weit, segmentweise validiert — Exit 2 bei
kaputtem Glob), fail-closed ohne lesbares `.git` bzw. verdrahteten Port (Exit 2).
CLI öffnet den git-Adapter für `tracked` **ohne** Range-Pflicht. Config-Surface:
`--print-mk doc-tracked` ([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
9→10), `--print-config`/`--suggest-config`, Handbuch 1.19, README (DE zuerst, EN
nachgezogen — sechzehn Module), CHANGELOG. Doc-first:
[`DC-FA-TRK-001`](../../../../spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)
(Lastenheft 0.37.0, Review-Präzisierung 0.37.1) +
[ADR-0030](../../adr/0030-tracked-referenz-ziele.md) (bei Closure Accepted) +
Spezifikation `.a`/§2/§4 gingen dem Code voraus.

**Belege.**
- `make ci` **grün** (doc-check 178/0, lint, test, arch-check via a-check 0 Befunde,
  Coverage, semgrep 0/55, gate-consistency, planning-check; image-test nativ ==
  Container); `make fullbuild` **grün** (38 Anforderungen/**0 Waisen**).
- **AK-Tests** an git-Fixtures (committed/staged/untracked; Mem-FS + fakeVCS + E2E-
  Temp-Repos): alle sieben Kriterien; **Mutations-Verriegelung** beider Fail-closed-
  Guards (Glob-Validierung entfernt ⇒ genau `TestTracked_UngueltigesGlobExit2` rot;
  nil-Port-Guard entfernt ⇒ genau `TestRunWithVCS_TrackedFailClosedOhnePort` rot) —
  von R2 unabhängig reproduziert; nach R2 zusätzlich Bild-Einschluss, Auflösung
  (`./`-Form) und Symlink-Skips verriegelt (die IsImage-Skip-Mutation überlebte
  zuvor — R2-Fund).
- **Beleg-Lauf** `--enable tracked` gegen das eigene Repo: **grün** (Index-Zugriff
  im read-only-Mount); adversariale Probe (untracktes, verlinktes Ziel) ⇒ rot mit
  `target-untracked`, Revert ⇒ grün.
- **Zwei unabhängige Reviews** (Reports
  [r1](../../../reviews/2026-07-03-slice-059-tracked-doc-r1.md)/[r2](../../../reviews/2026-07-03-slice-059-tracked-code-r2.md)):
  R1 doc (5 MEDIUM/3 LOW/1 INFO — u. a. Post-Pass-Wording, CLI-010-Neun-Drift,
  Verzeichnis-Ziel-Doku) und R2 code (4 MEDIUM/2 LOW/2 INFO — u. a. roh-vs-aufgelöst,
  Verzeichnis-Symlink-false-positive per Klon-Gegenprobe, überlebende
  Bild-Skip-Mutation, Glob-Segment-Validierung); **alle Befunde eingearbeitet** und
  regressions-verprobt (Lastenheft 0.37.1).
- Release **v0.37.0** auf GHCR — Pipeline-Run + Digest-Pin folgen als
  Digest-Backfill nach dem Tag.

**Lerneintrag.** (1) Rollen-Teilung trägt: R1 (doc) fand die Spec-Treue-Drifts
(„Post-Pass", Neun-Targets), R2 (code) die semantischen Ränder (Symlink-Alias vs.
Index-Realpfad — nur per Klon-Gegenprobe sichtbar) und die überlebende
Bild-Skip-Mutation; keiner der Funde überlappte. (2) Die slice-057/058-Lehren
zahlten direkt ein: Mutations-Verriegelung je Guard und fail-closed-Proben-Harness
waren von Anfang an Teil der DoD — R2 bestätigte beide unabhängig. (3) Bestands-
Folgepunkt: die `--doctor`-Klartext-Liste (`AllReasons`/`reasonTexts`) hinkte seit
v0.25 sieben Grund-Codes hinterher (hier `target-untracked` ergänzt; der Rest ist
Kandidat für einen Hygiene-Slice samt Vollständigkeits-Verriegelung §4 ↔
`AllReasons`).
