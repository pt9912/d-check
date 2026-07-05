# Slice slice-063: Modul `targets` — Deklarations-Konsistenz Doku ↔ Build-Targets

**Status:** done (welle-52-targets-modul, Release **v0.38.0**). Alle DoD-Punkte
erfüllt (§3): Modul + §4-Grund-Codes + Akzeptanztests, Paritäts-Mutations-Beleg
vs. `gate-consistency.sh`, Dogfood + Selbstbezug, Config-Surface. **Zwei
unabhängige Reviews** (Fundament R1 + Impl R2, beide ACCEPT/eingearbeitet, §8);
`make gates`/`make ci`/`make fullbuild` grün (39 Anforderungen, 0 Waisen).
Lifecycle-Move + Roadmap-Flip atomar
([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise));
Closure-Notiz §7.

**Welle:** welle-52-targets-modul.

**Bezug:** [`DC-FA-TGT-001`](../../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(Modul `targets`, Bereich `TGT`) +
[ADR-0031](../../adr/0031-targets-deklarations-konsistenz-modul.md)
(Identitäts-Ausweitung „Doku-Checker → Deklarations-Konsistenz-Checker" +
Makefile-Lese-Entscheidung, macht die Roadmap-Einstufung „gate-consistency nicht
d-check-fähig" hinfällig). 17. Regelmodul.

**Autor:** pt9912. **Datum:** 2026-07-05.

---

## 1. Ziel

`tools/gate-consistency.sh` ist ein **kopiertes** Meta-Gate über die Repo-Familie
(d-check, a-check, belief-agent, …) und **driftet** bereits real auseinander
(a-checks Fassung hat einen Pin-Block, eine Utility-Allowlist und `d-check.mk`
als zweite Target-Quelle — d-checks nicht). Sein prüfwürdiger, cross-repo
**gemeinsamer Kern** — „jedes in der Doku als ` `make X` ` behauptete Gate
existiert real, und jedes reale Gate ist dokumentiert" — ist genau eine
**Doku-Behauptung über die Repo-Realität**, also d-checks Mission („Meta-Gate
gegen Harness-Lügen, kein halluziniertes Gate"). Neu: ein opt-in Modul
`targets` mechanisiert diesen Kern als **verteilbares** Modul (Image +
` .d-check.yml ` + `--print-mk`), sodass jedes Repo **dasselbe Binary mit
eigener Config** fährt statt eine driftende Skript-Kopie. Der Kopie-Drift ist
damit an der Wurzel getrocknet — wie zuvor bei `adr-immutable→vcs`,
`trace→commits`, `planning→planning`.

Präzedenz-Form: **`targets` verhält sich zu (Doku-Tabelle ↔ Makefile-Regeln) wie
`planning` zu (Roadmap-Doku ↔ `in-progress/`-Verzeichnis)** — beide prüfen
hermetisch, ob eine Doku-Behauptung mit der Repo-Struktur übereinstimmt; beide
lesen nur (`ReadFile`/Listing), führen nichts aus. Genau diese Form (die
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
mit slice-057 etablierte) hebt die alte „nicht d-check-fähig, weil Markdown ↔
Makefile ↔ YAML"-Einstufung auf.

## 2. Entscheidungen (Entwurf — R1 eingearbeitet)

- **Zwei Richtungen, zwei Grund-Codes `gate-phantom` / `gate-undocumented`**
  (R1-F-3: **nicht** `target-*` — das kollidiert semantisch mit
  `target-missing` [`links`] / `target-untracked` [`tracked`], die ein
  *Link-Ziel* meinen; hier ist ein *Build-Gate* gemeint. Modulname bleibt
  `targets`; Modul↔Code-Präfix müssen nicht gleich sein — Präzedenz `links` →
  `target-missing`):
  - **Richtung 1 — Phantom-Gate:** jedes ` `make X` `-Token in einer
    konfigurierten Doku-Datei muss eine Regel in einer konfigurierten
    Makefile-Quelle sein; sonst **`gate-phantom`** (Befund an Datei:Zeile der
    Doku-Behauptung).
  - **Richtung 2 — undeklariertes Gate:** jede Makefile-Regel (minus
    `targets.exempt-targets`) muss in der Autoritäts-Doku (`targets.authority`)
    als ` `make X` ` stehen; sonst **`gate-undocumented`** (Befund an
    Datei:Zeile der Makefile-Regel).
- **Tabellen-Scoping ist Vertrag (R1-F-1, fidelity-kritisch):** ` `make X` `
  wird **ausschließlich aus Tabellenzeilen** extrahiert (Zeilen mit `|`-Präfix,
  wie das abgelöste Skript `grep -E '^\|'`) — **nicht** aus Prosa. AGENTS.md
  trägt `make X` auch in Prosa (§3.1 „Richtig: `make gates`", §3.3
  `make planning-check`, …); die zählen **nicht** als Existenz-Behauptung. Ohne
  diese Regel entstünden spuriöse `gate-phantom` aus Prosa (ein in Prosa
  erwähntes, entferntes Target). Der Doc-first-Spec schreibt die
  `^|`-Scoping-Regel explizit fest.
- **Whole-file, kein Sektions-Scoping (R1-F-4):** Richtung 1+2 scannen die
  ganze Datei nach Tabellenzeilen (fidelity ≡ Skript). Die Config kennt **keinen
  Heading-Key** wie `planning.heading` — daher **keine** „§4"/„§Sensors"-
  Präzision in der Spec-Prosa behaupten (ein Tabelleneintrag außerhalb §4
  befriedigt Richtung 2 legitim). Optionaler `authority`-Section-Anker = spätere
  Erweiterung, nicht dieser Slice.
- **Config-Schema** (`targets.*`):
  ```yaml
  targets:
    makefiles: [Makefile]                        # Regelnamen-Quelle(n); Consumer: + das mk-Fragment
    doc-tables: [AGENTS.md, harness/README.md]    # Dateien mit make-X-Tabellen (Richtung 1)
    authority: AGENTS.md                          # Vollständigkeits-Quelle (Richtung 2)
    exempt-targets: []                            # Regeln ohne Doku-Pflicht (s. u.)
  ```
  Leeres `makefiles` **oder** `doc-tables` ⇒ Modul inert (wie `planning` ohne
  Roadmap-Pfad). Fehlt `authority`, entfällt Richtung 2.
- **`exempt-targets` (R1-F-5):** in d-checks **eigener** Config bleibt es
  **leer** — `help`/`clean`/… stehen heute alle in AGENTS.md §4, die
  erschöpfende Makefile→§4-Prüfung soll erhalten bleiben. Ein `[help, clean]` im
  Schema-Beispiel ist reine **Konsumenten-Illustration** (Repos, die Utilities
  bewusst nicht dokumentieren), **nicht** d-checks Config.
- **Hermetisch, reiner Filesystem-Port** — `ReadFile` auf die konfigurierten
  Dateien; **kein** Makefile-*Ausführen*, kein git, kein Netz. Kein neuer Port
  (dieselbe `ReadFile`-Mechanik wie `planning`); der Makefile-„Parser" ist die
  **gleiche Regex-Heuristik** wie im Skript (literale Regelnamen, außer
  Zuweisungen/`.PHONY`) — **keine** Pattern-Rules/variabel benannten Targets
  (Harness-Makefiles nutzen literale Namen; Fidelity ≡ Skript).
- **Strikt opt-in, diagnose-only, fail-closed, default-aus byte-identisch** —
  wie `planning`: nie Default; ohne aktives `targets` byte-identisch
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)); aktiv
  mit fehlender konfigurierter Datei ⇒ Exit 2; kein `--repair`-Hunk; kein Netz/
  kein Schreiben
  ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Scope = Kern (Richtung 1+2).** Der ` .d-check.yml `-Modul-Listen-Check
  (heutige gate-consistency-Prüfung 3) bleibt draußen: repo-spezifische
  Selbst-Konsistenz der Netzlos-Gate-Config, **nicht** der cross-repo-driftende
  Kern.
- **Dogfood + Retirement.** `make gate-consistency` wird auf den Modul-Dogfood
  umgestellt (analog `make planning-check`) und behält nur den kleinen
  Prüfung-3-Rest. Cross-repo-Gewinn: a-check/belief-agent/… ersetzen ihre
  gate-consistency-Kerne durch `--enable targets` + Config (Skript-Kern-
  Tombstone; a-checks Pin-Block bleibt a-check-lokal).
- **Bootstrap-Selbstbezug:** d-check prüft dann seine *eigene* Doku↔Makefile-
  Konsistenz über sein eigenes Binary. Absicherung: modul-interner
  **Negativ-Selbsttest** (Guard entfernt ⇒ Test rot), wie bei
  `planning`/`commits`/`vcs`.

## 3. Definition of Done (R1 eingearbeitet)

- [x] **Modul `targets`** (Kern-Regel + Config + fail-closed) + Akzeptanztests
  (Happy/Boundary/Negative je Richtung, Determinismus, default-aus
  byte-identisch) + Bootstrap-Negativ-Selbsttest.
- [x] **Paritäts-Nachweis (R1-F-2):** der Modul-Befundsatz ≡ der Befundsatz des
  abgelösten `gate-consistency.sh` (Richtung 1+2) auf dem aktuellen Repo-Baum —
  Belege je Richtung (Mutations-Proben: dokumentiertes Phantom-Target,
  undokumentiertes Makefile-Target), wie arch→a-check/vcs/commits. Ein aktives
  Gate wird nur mit belegter Parität stillgelegt.
- [x] **Doc-first:** neue Anforderung (Bereich `TGT`) im Lastenheft (+ Modul-
  Liste + Glossar + Makefile-Fragment-Vertrag + Historie/Version); Spezifikation
  (Algorithmus **inkl. `^|`-Tabellen-Scoping**, Schema `targets.*`, Grund-Codes
  `gate-phantom`/`gate-undocumented`); begleitende ADR (Motivation **mit
  auditierbarem Pointer auf a-checks `gate-consistency.sh`-Fassung**, R1-F-9) +
  ADR-Index.
- [x] **Dogfood + Selbstbezug (R1-F-6):** `make gate-consistency` auf das Modul
  umgestellt (Kern), d-checks eigene ` .d-check.yml ` `targets:`-Config (mit
  leerem `exempt-targets`) aktiviert; **die `make gate-consistency`-Zeilen in
  AGENTS.md §4 und `harness/README.md` §Sensors** an die Umstellung angepasst —
  und genau diese Tabellen prüft `targets` (wechselseitig konsistent).
- [x] **Config-Surface:** `--print-config`/`--suggest-config`/`--print-mk`
  (`doc-targets`) + Benutzerhandbuch §5/§6.
- [x] `make gates`/`make ci`/`make fullbuild` grün; **zwei unabhängige Reviews**
  (dieses Doc-first-Fundament-Review R1 + Impl-Review); Closure-Move + Body +
  **Lerneintrag** (Modul 5); Release (nutzersichtbar: neues Modul) +
  Digest-Backfill.

## 4. Risiken / offene Punkte

- **Slice-Größe (Modul 5, ≤3 DoD):** der Umfang liegt am oberen Rand. Das
  Fundament-Review R1 hat die Größe **nicht** als Finding markiert, und die
  Präzedenz-Modul-Slices (057 `planning`, 059 `tracked`) trugen denselben
  Umfang (Modul + Doc-first + Config + Dogfood + Release) als *ein* Slice. Der
  Dogfood-/Selbstbezug (§3 R1-F-6) ist zudem intrinsisch an das Modul gekoppelt.
  **Vorschlag: ein Slice** (statt 063a/063b); Split bleibt Option, falls
  gewünscht.
- **`gates`-Einbindung:** `targets` als Dogfood in `make gate-consistency`
  (bleibt in `gates`) statt neues `make targets-check`.
- **Autoritäts-Doku (Richtung 2):** `AGENTS.md` §4 als Vollständigkeits-Quelle;
  `harness/README.md` §Sensors nur Richtung-1-Behauptung. Section-Anker als
  spätere Präzision (R1-F-4), nicht dieser Slice.
- **Makefile-Parser-Fidelity:** literale Regelnamen; Pattern-Rules/variable
  Targets nicht erkannt — dokumentierte Grenze, ≡ Skript.

## 5. Trigger

Auftraggeber 2026-07-05: „driftende Skripte aus den verschiedenen Repos
einfangen" + Bestätigung „d-check = declaration-checker, Klasse A regelkonform".
Cross-repo-Beleg: a-checks `gate-consistency.sh` ≠ d-checks (Pin-Block,
Utility-Allowlist, zweite Target-Quelle) — der **auditierbare Pointer** darauf
gehört in die ADR-Motivation (R1-F-9). Klasse-A-Muster
([`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)):
mechanisieren → ein verteiltes Modul → Skript-Kern-Tombstone. Fünf Präzedenzen:
adr-immutable→`vcs`, completeness→Flag, trace→`commits`, planning→`planning`,
arch→a-check.

## 6. Sub-Area-Modus-Begründung

GF (neues Regelmodul im bestehenden Hexagon-Schnitt; Kern-Regel in `rules`,
Filesystem-Port-Nutzung wie `planning`; Doc führt — Inventur prüft nur
Code-Konformität). Kein neuer Adapter/Port, keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

**Ergebnis:** das 17. opt-in Regelmodul `targets` ist released (v0.38.0). Es prüft
**hermetisch** (nur der Filesystem-Port, kein git/Netz/Makefile-Ausführen) die
Deklarations-Konsistenz Doku ↔ Build-Targets: ein in einer Doku-**Tabellenzeile**
behauptetes `make X` ohne Makefile-Regel ⇒ `gate-phantom` (Richtung 1); jede
Makefile-Regel (minus `exempt-targets`) ohne Autoritäts-Doku-Eintrag ⇒
`gate-undocumented` (Richtung 2). fail-closed bei fehlender Datei, default-aus
byte-identisch.

**Retirement (der Kern-Wert).** `make gate-consistency` dogfoodet nun das Modul
für den Doku-↔-Makefile-Kern (via Image, `--enable targets`);
`tools/gate-consistency.sh` ist auf die repo-spezifische
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Modullisten-Prüfung
reduziert. Der cross-repo-driftende Skript-Kern ist damit **verteilbar
mechanisiert** — a-check/belief-agent/… ersetzen ihn künftig durch
`--enable targets` + eigene Config (vierte Tombstone-Klasse nach
`adr-immutable→vcs`, `trace→commits`, `planning→planning`).

**Paritäts-Beleg (DoD-kritisch).** Auf dem echten Repo-Baum sind Modul und Skript
beide grün; unter zwei Mutationen (ein Phantom-Target in AGENTS.md §4, ein
undokumentierter Makefile-Target) feuern **Skript und Modul identisch** auf
dieselben zwei Targets — das Modul zusätzlich mit präziser Datei:Zeile. Der
Impl-Review bestätigte die Regex-Fidelity als **zeichensatz-exakt** zum Skript.

**Config-Surface.** `--print-mk` trägt das elfte Target `doc-targets`
([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben),
10→11 — schloss die einzige Vertrags-Lücke), `--print-config`/`--suggest-config`
führen `targets`, Benutzerhandbuch §5/§6/§4.13/§11.

**Reviews.** Zwei unabhängige Reviews: Fundament R1 (NACHBESSERN, 6 Befunde
eingearbeitet — u. a. `^|`-Tabellen-Scoping als Spec-Vertrag, Paritäts-DoD) +
Impl R2 (ACCEPT, 0H/0M/1L/1I — F-1 Vorkommen-Granularität + F-2
`exempt-targets`-exakt-vs-Glob als Klarstellungen eingearbeitet).

**Commit-Kette:** doc-first-Fundament → Reviews → Plan/welle-52 → Feat-Kern →
Dogfood/Retirement → Config-Surface → Handbuch → Impl-R2 → Release-Prep →
Closure-Move → Closure-Body (Tag `v0.38.0`) → Digest-Backfill.

**Lehren:** (i) ein neues Modul ist gleichzeitig ein **Retirement** (Skript-Kern)
und ein **Selbstbezug** (d-check prüft seine eigene Doku↔Makefile-Tabelle) — die
Dogfood-Config und die `make gate-consistency`-Zeilen in AGENTS.md §4 /
`harness/README.md` §Sensors müssen im selben Zug konsistent bleiben. (ii) Der
`--print-mk`-Target-Count ist ein **stiller Vertrag** (Lastenheft sagte 11, Code
emittierte 10; kein Gate fing es) — beim Modul-Slice den Print-mk-Count mitziehen.
(iii) Modul- und Code-Präfix müssen nicht gleich sein (`targets` → `gate-*`), um
Grund-Code-Kollisionen mit `target-missing`/`target-untracked` zu vermeiden.

## 8. Review-Nachtrag (Design-Review R1)

Unabhängiges Design-Review (0 HIGH / 2 MEDIUM / 4 LOW / 3 INFO) — Einarbeitung:

- **F-1 (MED, Tabellen-Scoping):** `^|`-Scoping als Spec-Vertrag aufgenommen (§2).
- **F-2 (MED, Paritäts-Nachweis):** eigener DoD-Punkt (§3).
- **F-3 (LOW, Reason-Codes):** `target-*` → `gate-phantom`/`gate-undocumented` (§2).
- **F-4 (LOW, §-Scoping):** whole-file Tabellenzeilen; „§4"-Präzision aus der
  Prosa gestrichen, Section-Anker als spätere Option (§2/§4).
- **F-5 (LOW, exempt-Beispiel):** d-checks eigenes `exempt-targets` leer;
  `[help, clean]` als Konsumenten-Illustration deklariert (§2).
- **F-6 (LOW, Sensors/§4-Selbstbezug):** AGENTS §4 + `harness/README.md`
  §Sensors-Updates in die DoD (§3), inkl. der wechselseitigen Konsistenz.
- **F-9 (INFO, a-check-Beleg):** auditierbarer Pointer in die ADR-Motivation (§3/§5).
- **F-7 (INFO, bare IDs) / F-8 (INFO, in-progress-planning-drift):** bereits
  erledigt — die `next/`-Fassung erfindet keine noch nicht vergebenen IDs
  (kein vorgezogenes ADR-/Anforderungs-Token) und liegt in `next/` statt
  `in-progress/`; `make doc-check`/`planning-check` grün.

**Fundament-Review R1** (unabhängig, auf Lastenheft + Spezifikation + ADR;
[Report](../../../reviews/2026-07-05-slice-063-foundation-r1.md)) — **NACHBESSERN**
(0 HIGH / 2 MEDIUM / 2 LOW / 2 INFO), alle eingearbeitet:

- **F-1 (MED):** Ziel-Zählung — das Lastenheft nannte im `--print-mk`-Vertrag
  11 Targets, die Spec-Mechanik-Sektion noch 10 (Strata-Widerspruch); Spec
  **10→11** + `doc-targets` nachgezogen, beide Historien konsistent.
- **F-2 (MED):** Tabellen-Scoping — Algorithmus sagte „getrimmter Anfang"
  (Einrückung erlaubt), Skript/ADR sagen `^\|` (Spalte 0); Algorithmus auf
  Spalte-0-Parität korrigiert.
- **F-3 (LOW):** Lastenheft-Historie trägt jetzt denselben §4-Deferral-Vorbehalt
  wie die Spec (nicht mehr einseitig).
- **F-4 (LOW):** 7. Akzeptanzkriterium „Boundary (exempt)" ergänzt (die einzige
  Skript-Divergenz `targets.exempt-targets`).
- **F-5 (INFO):** Doku-Token-Zeichensatz `X` = `[a-z][a-z0-9_-]*` (Skript-Parität)
  spezifiziert.
- **F-6 (INFO):** Inert nur an `targets.makefiles`; Richtungen entkoppelt
  (Lastenheft + Spec + §2-Schema + ADR).
