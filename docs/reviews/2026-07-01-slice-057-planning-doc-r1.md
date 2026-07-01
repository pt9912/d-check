# Review-Report — slice-057 (Modul `planning`) — Doc-Straten & Cross-Surface (R1)

**Datum:** 2026-07-01
**Reviewer-Rolle:** unabhängig/adversarial, Fokus **Doc-Straten + Cross-Surface-Konsistenz**
(NICHT Go-Code — separater zweiter Reviewer).
**Gegenstand:** neues hermetisches opt-in Modul `planning` (15.), mechanisiert das
entfernte `tools/planning-consistency.sh`. Diff `git diff HEAD` + neue Dateien
(ADR-0028, slice-057, `rules/planning*.go`, `cli_planning_test.go`).
**Baseline:** `.harness/skills/reviewer.md` (Kategorien/Schema/Negativbefund-Pflicht),
`AGENTS.md` §3, `harness/conventions.md` (MR-006, MR-013, MR-014).
**NICHT geprüft:** DoD-Abhakung (Verifikations-Rolle), Go-Modul-Korrektheit (2. Reviewer).

**Verifikations-Läufe (Image `d-check:latest`, `--network none`, read-only):**

- `--print-mk | grep -cE '^doc-[a-z-]+:'` = **9** (doc-check, -trace, -complete,
  -doctor, -repair, -immutable, -commits, **-planning**, -help).
- Default-`doc-check`: `166 Datei(en), 0 Befund(e)`, Exit 0 (Tombstone-Logik trägt).
- Dogfood `--enable planning`: `0 Befund(e)` (slice-057 in-progress ⟺ Roadmap benennt
  aktive Welle → `hasActive==hasSlices`).
- Negativ-Probe (Scratch-Repo): Drift „idle-Marker + Slice vorhanden" ⇒ 1×
  `planning-drift` (Zeile = Heading-Zeile, target = Slice-Verzeichnis); Heading-Tippfehler
  ohne Slice ⇒ 1× `planning-drift` (Zeile 1, fail-closed). Modul ist **nicht** still-grün.
- `--print-config`/`--suggest-config` führen `planning`.

---

## Findings

### MEDIUM-1 — CHANGELOG behauptet Handbuch-Doku, die nicht existiert

- **kategorie:** MEDIUM
- **quelle:** Maintainability / Harness-Ehrlichkeit (Doku-Phantom, gleiche Klasse wie
  ein behauptetes Gate)
- **pfad:** `CHANGELOG.md:68` vs. `docs/user/benutzerhandbuch.md` (unverändert)
- **befund:** Der CHANGELOG-Eintrag 0.36.0 schließt mit „Benutzerhandbuch §5/§6
  dokumentiert das Modul." Das Benutzerhandbuch ist jedoch **vollständig
  un-nachgezogen**: die Modul-Tabelle (§5/§6, Zeilen 805–817) endet bei `commits`
  ohne `planning`-Zeile; die Versions-Historie (Zeile 924) endet bei `1.17 / v0.35.0`
  ohne 0.36.0-Zeile; der `--print-mk`-Abschnitt (Zeile 620–622) nennt weiter „acht"
  Targets, `v0.35.0` und listet `doc-planning` nicht. Die Vorgänger `vcs` (slice-053)
  und `commits` (slice-056) trugen Modul-Zeile **und** Historie-Zeile nach; hier
  fehlt beides, während der CHANGELOG die Erledigung im Präsens behauptet.
- **Failure-Szenario:** Ein Nutzer liest CHANGELOG 0.36.0, folgt dem Hinweis auf
  Handbuch §5/§6 und findet **keinen** `planning`-Eintrag — die Aussage ist zum
  Review-Stand faktisch falsch. (Ob das Handbuch als späterer Release-Prep-Schritt
  nachgezogen wird oder die CHANGELOG-Aussage zu entschärfen ist, muss vor Merge
  entschieden werden — beides zugleich stimmt derzeit nicht.)
- **verifizierbar:** ja (Beobachtung reproduzierbar per `grep -n planning
  docs/user/benutzerhandbuch.md` → 0 Modul-/Target-Treffer). Kein Gate erzwingt die
  CHANGELOG-↔-Handbuch-Konsistenz.

---

## Findings (LOW)

### LOW-1 — DC-FA-CLI-010 Boundary-AK zählt `doc-planning` nicht auf

- **kategorie:** LOW
- **quelle:** `DC-FA-CLI-010` (Makefile-Fragment) — Boundary-Akzeptanzkriterium
- **pfad:** `spec/lastenheft.md:411`
- **befund:** Beschreibung („sowie **neun** … Targets"), Happy-AK (listet `doc-planning`
  „mit `--enable planning` … ohne Range") und Out-of-Scope („gelisteten **neun** …
  `doc-planning`") sind auf 9 gezogen; das **Boundary-AK** enumeriert die Ziel-Modi
  jedoch nur bis `doc-commits` (`… doc-commits → --enable commits … --range $(RANGE)`)
  und nennt `doc-planning` nicht. Der Vorgänger slice-056 hatte `doc-commits` in genau
  dieses Boundary-AK eingetragen (slice-053 zuvor `doc-immutable`) — das Aufzählungs-
  Muster ist „jedes neue Target mit Modus", hier nicht fortgeführt. Damit ist das
  Kriterium intern inkonsistent mit der eigenen „neun"-Zählung.
- **Failure-Szenario:** Ein Verifizierer, der das Boundary-AK als Checkliste der
  Per-Target-Aufruf-Modi nutzt, hat keine Zeile für `doc-planning`s hermetischen
  (Range-losen) Modus. Restrisiko gemindert: das Happy-AK bindet den Fragment-Inhalt
  inkl. „ohne Range" bereits, sodass ein Fehl-Modus dort auffiele — daher LOW
  (Aufzählungs-Drift in einem bindenden AK, kein ungebundener Vertrag).
- **verifizierbar:** nein (Prosa-Vollständigkeit; kein Gate prüft die AK-Enumeration).

---

## Negativbefunde (geprüft, ohne blockierenden Befund)

- **MR-006 / DC-FA-MTX-003 (Referenz-Richtung, Provenance):** Die einzigen bare
  `slice-\d{3}`-Token in ADR-0028 stehen in `## Geschichte` (Zeile 145: `slice-057`,
  `slice-040`); `.d-check.yml` `matrix.exclude-sections` enthält `Geschichte`. Kontext/
  Entscheidung/Konsequenzen tragen Slices nur als Markdown-Link bzw. als `slice-*`
  (Glob, kein `\d{3}`-Token). `matrix`-Logik somit korrekt bedient, `doc-check` grün —
  Provenance ist **in die Geschichte verlagert**, nicht im Körper. Kein Spec-Stratum
  verweist abwärts auf ADRs/Wellen/Slices/Closure-Daten.
- **Spec-Straten-Abwärtsverweise (MR-006 Kern):** `spec/lastenheft.md` und
  `spec/spezifikation.md` referenzieren im neuen Material **keine** ADR/Slice/Welle.
  Die einzige neue AGENTS.md-Referenz (Lastenheft 1279, `AGENTS.md §3.3`) und der
  `MR-013`-Verweis (1316) sind Verweise auf die durchgesetzte **Konvention** — analog
  zum Vorgänger `DC-FA-COMMITS-001`, der auf `harness/README.md §Traceability rules`
  (Rang 9) verweist; AGENTS.md/MR sind **nicht** in der MR-006-Verbotsliste
  (ADRs/Planning-Artefakte). Der neue lange Anker `#33-git-mv…` ist bei Umbenennung
  vom `anchors`-Gate gedeckt (self-healing). Kein Befund.
- **Anforderungs-Anlege-Prozess (DC-FA-PLAN-001):** ID nach Schema, Bereich `PLAN` in
  §3 deklariert; **vier** AK (Happy/Boundary(Modul-aus)/Negative/fail-closed) —
  Superset des 3er-Trios — plus explizite Out-of-Scope-Liste; Versions-Bump 0.36.0 +
  §7-Historie-Zeile; fail-closed benannt; Schärfungs-Richtung gewahrt (ADR schärft
  Spec, Lastenheft direkt geändert). AK decken sich mit §DC-FA-PLAN-001.a
  (`hasActive==hasSlices`, Heading-Guard, fail-closed).
- **Konsistenz Lastenheft↔Spec↔ADR↔Slice↔.d-check.yml:** Grund-Code `planning-drift`,
  Schema-Keys (`planning.roadmap`/`heading`/`marker`/`slice-glob`) und Defaults
  (`## Aktuelle Welle` / „Keine aktive Welle" / `slice-*.md`) sind über alle Schichten
  deckungsgleich. Die Finding-Felder der Spezifikation (§.a Schritt 5: `file`=Roadmap,
  `line`=Heading-Zeile bzw. 1, `target`=Slice-Verzeichnis) sind per Negativ-Probe
  **exakt** so realisiert. „15. Modul" konsistent über alle Surfaces; `validModules`
  = 15 Einträge.
- **8→9-Vollständigkeit (bindender Körper):** `--print-mk` liefert real 9 Targets;
  Beschreibung „neun", Happy-AK, Out-of-Scope „neun", `config_template.go`, `print_mk.go`
  (VIER fmt-Verben), spezifikation §.a „Neun `.PHONY`-Targets" alle konsistent. `grep`
  nach `acht|sieben|sechs` in `spec/` ohne echten Stale-Treffer (nur „macht"/„gecacht"
  u. ä.); der einzige echte Stale-Zähler steht im Benutzerhandbuch (→ MEDIUM-1).
- **Tombstone (`codepaths.ignore-refs`):** `tools/planning-consistency.sh` ergänzt;
  referenz-weit ⇒ deckt **alle** Inline-Refs (immutable slice-040 done, ADR-0028,
  slice-057, roadmap, CHANGELOG, spezifikation §7). `doc-check` grün verifiziert die
  Deckung end-to-end.
- **Config-Surface-Vollständigkeit:** `--print-config` (Verfügbar-Liste + eigener
  `planning`-Block), `--suggest-config` (opt-in-Liste), `--print-mk` (`doc-planning`),
  `config_template.go`, Glossar, `DC-FA-CLI-002`, `model.validModules()` führen alle
  `planning`. **Ausnahme:** Benutzerhandbuch §5/§6 (→ MEDIUM-1).
- **„nachrangig"-Ehrlichkeit:** ADR-0028 benennt den begrenzten, layout-spezifischen
  Verteilungswert in Kontext, Entscheidung, Alternativen-Contra und Konsequenzen;
  Roadmap und Lastenheft §7 spiegeln „bewusst nachrangig". Nicht überverkauft.
- **Bindepunkt `planning-check` in `gates` (checkpoint 8):** Makefile `gates:`
  enthält weiter `planning-check` (`planning-check: build`, `$(DCHECK_RUN) --enable
  planning $(FOCUS_DISABLE)`; `FOCUS_DISABLE` definiert, kein Phantom). AGENTS §4
  („in `gates`"), harness/README §Sensors („in `make gates`", `gates`-Aggregat-Zeile),
  ADR-0028 Entscheidung („anderer Bindepunkt als vcs/commits — in `make gates`")
  konsistent; `adr-check`/`trace-check` bleiben korrekt außerhalb `gates`.
- **ADR-Status/Index:** ADR-0028 `Proposed` (bei Closure → Accepted, DoD), ADR-Index
  ergänzt, `Schärft:` zeigt aufwärts auf §DC-FA-PLAN-001.a, `Bezug:` auf das Lastenheft.
- **MR-014 Doc-Struktur:** ADR-0028 und slice-057 folgen dem Haus-Stil (ADR:
  Alternativen-Tabelle, zweispaltige Geschichte; Slice: §1–§7-Reihenfolge, Closure-Notiz
  §7 Platzhalter). Konsistent mit dem Bestand.

---

## Kategorie-Summary

| Kategorie | Anzahl |
| --- | --- |
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

## Verdikt

**NACHBESSERN.** Kern-Logik, Spec-Straten, Referenz-Richtung, Tombstone und
Config-Surface sind sauber und dogfood-verifiziert (Modul nicht still-grün,
Finding-Felder spec-treu). Blockierend ist **MEDIUM-1**: der CHANGELOG behauptet eine
Benutzerhandbuch-Doku, die nicht existiert — vor Merge/Release entweder Handbuch §5/§6
nachziehen (wie bei `vcs`/`commits`) oder die CHANGELOG-Aussage entschärfen. **LOW-1**
(Boundary-AK ohne `doc-planning`) sollte im selben Zug mitgezogen werden. Beides sind
Cross-Surface-Abgleiche, keine Kern-Defekte.
