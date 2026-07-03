# Review-Report — slice-058 (arch-check via a-check) — Plan-Review (R1)

**Datum:** 2026-07-03
**Reviewer-Rolle:** unabhängig/adversarial, Fokus **Plan-/Doc-first-Stand**
(slice-058 + ADR-0029, beide committet `1920a05`; Implementierung steht aus).
**Gegenstand:** `docs/plan/planning/in-progress/slice-058-arch-check-via-a-check.md`
und [ADR-0029](../plan/adr/0029-arch-check-via-a-check.md) (Proposed — noch
editierbar, `adr-check` schützt nur `Accepted`) gegen das abzulösende
`tools/arch-check.sh`, den Makefile-/Dockerfile-Bestand und den **realen
a-check-v0.6.0-Stand** (lokaler Checkout des [Schwester-Repos](https://github.com/pt9912/a-check), Spec + Quelltext).
**Baseline:** `.harness/skills/reviewer.md` v1.2.0, `AGENTS.md` §3,
`harness/conventions.md` (MR-007, MR-013, MR-014).
**NICHT geprüft:** DoD-Abhakung (Verifikations-Rolle; Impl existiert noch nicht),
a-check-Repo-interne Qualität (eigener Harness drüben).

**Verifikations-Proben (Quellen-Lektüre statt Gate-Läufe — es gibt noch keinen Diff):**

- a-check `tech-leak`-Auswertung: `internal/hexagon/core/rules.go:50` — Match ist
  `strings.Contains(f.Path, tech.Adapter)`, d. h. **ein** Pfad-Substring je Pattern;
  `matchTech` = Erst-Treffer in Deklarationsreihenfolge (ein zweiter Eintrag mit
  gleichem Pattern feuert nie).
- a-check `composition_root`: `rules.go:18-19` — `continue` **vor jeder** Prüfung;
  SPEC-RULE-001 bestätigt: ausgenommen von **allen** Schicht-Regeln **und** `tech-leak`.
- a-check Scanner: `extract/extract.go` WalkDir skippt **nur `.git`**; jede Datei,
  die einem `languages`-Glob entspricht, wird gescannt — **inkl. `*_test.go`**;
  Glob-Engine (`rules.go:411` ff., `globToRegexp`) kennt `**`/`*`/`?`, **keine
  Negation** → Test-Dateien sind per Config nicht ausschließbar.
- d-check-Bestand: `internal/adapter/driven/git/git_test.go:4` importiert `os`;
  `internal/adapter/driven/report/report.go:14` importiert `gopkg.in/yaml.v3`.
- Skript-Regelzweige: `tools/arch-check.sh:49-57` (R2 = **zwei** Kapseln: `net/http`
  → httpcheck **und** go-git → git-Adapter), `:59-64` (R3 erlaubt yaml in
  **configyaml und report**, ADR-0009).
- `tools/gate-consistency.sh`: beide Richtungen laufen über
  `makefile_targets Makefile` — geparst wird **nur die Datei `Makefile`**,
  include-Fragmente nicht; das `--print-mk`-Fragment definiert das Target
  **`a-check`**, nicht `arch-check` (Beleg: `a-check.mk` im Schwester-Repo-Checkout).
- Kontext-Fakten bestätigt: a-check v0.6.0 getaggt, Digest-Pin
  `sha256:b349a150…` im Fragment, Lauf `--network none` + `:ro`; sieben
  Sprach-Backends (go/cpp/rust/kotlin/java/python/csharp); `.d-check.yml:89`
  trägt genau die vier bestehenden `ignore-refs`-Tombstones.

---

## Findings

### HIGH-1 — Regel-Übersetzung in §2/ADR-0029 lässt die R2-go-git-Kapsel weg und bindet yaml nur an configyaml

- **kategorie:** HIGH
- **quelle:** ADR-0009, ADR-0024, ADR-0005 (R2/R3); `AGENTS.md` §3.6 /
  Plan-eigene „keine stille Lockerung"-Zusage
- **pfad:** `docs/plan/planning/in-progress/slice-058-arch-check-via-a-check.md:46-50`;
  `docs/plan/adr/0029-arch-check-via-a-check.md:66-68`
- **befund:** Die `tech`-Aufzählung („`net/http` → httpcheck-Adapter, `yaml` →
  configyaml-Adapter, `os` → fs-Adapter") weicht in zwei Punkten vom abzulösenden
  Skript ab: (a) R3 erlaubt yaml in **zwei** Adaptern (`configyaml` **und**
  `report`, ADR-0009; `tools/arch-check.sh:59-64`) — `report.go:14` importiert
  yaml real; (b) die zweite R2-Kapsel **go-git → git-Adapter** (ADR-0024;
  `tools/arch-check.sh:52-57`) fehlt vollständig. Zusätzlich deckt die Liste den
  R1-Bann von `syscall`/`io/fs`/`net`(-Sockets)/übrigen `net/*` im Kern nicht ab
  (nur `net/http`/`yaml`/`os` genannt). Verschärfend: ein
  Ein-Pattern-zwei-Adapter-Binding ist in a-check v0.6.0 **nicht ausdrückbar**
  (`tech.Adapter` = ein Pfad-Substring, Erst-Treffer je Pattern) — der R3-Fall ist
  damit dieselbe Klasse „nicht per Config schließbar", die der Plan nur für R4
  erwägt.
- **Failure-Szenario:** Wörtlich umgesetzt ist `make arch-check` auf dem sauberen
  Baum **rot** (yaml im report-Adapter → `tech-leak`); die naheliegende
  „Reparatur" (`adapter: internal/adapter/driven/` breit binden bzw. go-git gar
  nicht deklarieren) ist eine **stille Lockerung** genau der Kapseln, die das
  Skript heute erzwingt — die Harness-Lüge-Klasse, die der Plan selbst
  ausschließt.
- **verifizierbar:** ja — nach Umsetzung: yaml-Mutations-Probe im fs-Adapter bzw.
  go-git-Import außerhalb des git-Adapters muss das Gate rot machen; mit der
  §2-Config tut sie es nicht (bzw. der saubere Baum ist schon rot).

### HIGH-2 — `_test.go`-Scope-Delta fehlt in §4: a-check scannt Test-Dateien, `go list` sah sie nie — konkret rot auf dem sauberen Baum

- **kategorie:** HIGH
- **quelle:** `DC-QA-04`-Analogon (Erkennungs-Differenz zur Alt-Mechanik),
  eskaliert (Gate-Pfad); a-check `AC-QA-02`
- **pfad:** `docs/plan/planning/in-progress/slice-058-arch-check-via-a-check.md:95-101`
  (Risiko-Abschnitt benennt nur Sonderfälle a+b)
- **befund:** Das Skript prüft via `go list -f '{{.Imports}}'` nur die
  **Nicht-Test-Imports** der Pakete; a-check scannt jede `languages`-Glob-Datei
  (`**/*.go`) **inklusive `*_test.go`**, und die Glob-Engine kennt keine Negation
  — Test-Dateien sind per Config nicht ausschließbar. Konkret:
  `internal/adapter/driven/git/git_test.go:4` importiert `os`; mit der geplanten
  Bindung `os` → fs-Adapter ist das ein `tech-leak` → Gate rot auf dem sauberen
  Baum. Der §4-Risikoblock behandelt nur den Präzisions-Trade der Extraktion,
  nicht diese **Scope-Erweiterung**; nach der Plan-eigenen Regel („nicht per
  Config schließbar ⇒ CR ans a-check-Lastenheft, Umstellung wartet") ist das ein
  heute schon feststellbarer CR-Kandidat, kein Proben-Fund.
- **Failure-Szenario:** Die Umstellung wird gebaut, `make arch-check` ist sofort
  rot; unter Zeitdruck wird `os` breiter gebunden oder als `ignore_symbols`
  entschärft — R4 verliert still seine Substanz.
- **verifizierbar:** ja — `a-check`-Lauf mit einer os→fs-`tech`-Bindung gegen den
  aktuellen Baum reproduziert den Befund ohne jede Mutation.

### MEDIUM-1 — `composition_root` ist eine Total-Ausnahme; §4(b) framet sie als sauberes R4-Ventil

- **kategorie:** MEDIUM
- **quelle:** a-check SPEC-RULE-001 (`composition_root` = ausgenommen von
  **allen** Schicht-Regeln **und** `tech-leak`); Skript-R2/R3 gelten heute
  repo-weit, auch für CLI/`cmd`
- **pfad:** `docs/plan/planning/in-progress/slice-058-arch-check-via-a-check.md:98-100`;
  `docs/plan/adr/0029-arch-check-via-a-check.md:114-115`
- **befund:** Der Plan löst die R4-Dreifach-Zone, indem CLI + `cmd` als
  `composition_root` deklariert werden („müssen … sauber herausfallen"). Nach
  a-check-Spec und Implementierung (`rules.go:18-19`, `continue` vor jeder
  Prüfung) fallen diese Pfade damit aus **sämtlichen** Prüfungen — das Skript
  flaggt dort heute aber `net/http` (R2), yaml (R3) und go-git. Das ist ein
  absehbares Paritäts-Delta (Deckungsverlust auf CLI/`cmd`), das §4 nicht als
  solches benennt.
- **Failure-Szenario:** Nach der Umstellung importiert `internal/adapter/driving/cli`
  yaml oder `net/http` — das Alt-Skript hätte rot gezogen, a-check bleibt grün;
  das Delta wurde nie dokumentiert, also gilt es fälschlich als „Parität belegt".
- **verifizierbar:** ja — Mutations-Probe „yaml-Import in `cli`" nach Umsetzung:
  bleibt grün; die DoD-Pflicht „Rest-Deltas explizit gelistet" muss diesen Fall
  tragen.

### MEDIUM-2 — „sechs Mutations-Proben, je Regel eine" verriegelt R2 nur halb (und die R3-Allow-Seite gar nicht)

- **kategorie:** MEDIUM
- **quelle:** slice-057-R3-Lehre (eine Probe muss den konkreten Guard verriegeln);
  `tools/arch-check.sh:49-57` (R2 = zwei unabhängige Kapseln)
- **pfad:** `docs/plan/planning/in-progress/slice-058-arch-check-via-a-check.md:79-81`
  (DoD), `:51-53` (§2)
- **befund:** Die Proben-Zählung folgt der Regel-Nummerierung (R1–R6 → sechs),
  nicht den Verbotszweigen des Skripts. R2 enthält zwei unabhängige Kapseln
  (`net/http`, go-git) — eine Probe lässt die jeweils andere unverriegelt.
  Für R3/R1 fehlt zudem die Allow-Seiten-Probe (yaml im report-Adapter bzw.
  `net/url` im Kern darf **nicht** flaggen) — genau die Fälle, die §4(a) als
  klärungsbedürftig benennt, tauchen in der Proben-Zählung nicht auf.
- **Failure-Szenario:** Alle sechs Proben grün-rot wie erwartet, go-git-Kapsel
  trotzdem nie im Netz — die Parität ist „belegt", ohne belegt zu sein
  (dieselbe Falle, die slice-057-R3 dokumentiert).
- **verifizierbar:** ja — Proben-Matrix je Skript-Verbotszweig (R1, R2a, R2b, R3,
  R4, R5, R6 + Negativ-Proben net/url und yaml@report) statt je Regel-Nummer.

### MEDIUM-3 — „Target-Name bleibt (gate-consistency hält)" ist nur bei bestimmter Verdrahtung wahr: gate-consistency parst nur die Datei `Makefile`, das Fragment-Target heißt `a-check`

- **kategorie:** MEDIUM
- **quelle:** `tools/gate-consistency.sh` (beide Richtungen via
  `makefile_targets Makefile` — kein include-Parsing); `a-check --print-mk`
  (Target `a-check`)
- **pfad:** `docs/plan/planning/in-progress/slice-058-arch-check-via-a-check.md:62-63`, `:40-41`
- **befund:** Der Plan lässt offen, **wo** das Target `arch-check` künftig
  definiert ist. Liegt es (umbenannt) im inkludierten `a-check.mk`, sieht
  `gate-consistency` es nicht → Richtung 1 rot („dokumentiert, Makefile kennt es
  nicht"), obwohl `make arch-check` funktioniert. Die Behauptung
  „`gate-consistency` hält Doku ↔ Makefile" stimmt nur, wenn das Target im
  `Makefile` selbst bleibt (und das Fragment z. B. nur Pin/Variable bzw. ein
  delegiertes `a-check`-Target liefert) — diese Verdrahtungs-Entscheidung gehört
  in §2. Randnotiz derselben Stelle: der Rename `a-check`→`arch-check` ist eine
  Divergenz vom generierten Fragment, die bei jeder Pin-Hebung via erneutem
  `--print-mk` re-appliziert werden muss.
- **Failure-Szenario:** Umsetzung nach Wortlaut („a-check.mk wird eingebunden",
  Target dorthin) → `make gates` rot am Meta-Gate; Fix unter Zeitdruck =
  gate-consistency-Parser aufweichen (Gate-Anfassen ohne Plan-Deckung).
- **verifizierbar:** ja — `make gate-consistency` nach Umsetzung.

---

## Findings (LOW/INFO)

### LOW-1 — Rückbau-Checkliste: `clean`-Target und Makefile-Kopf tragen weitere `arch-check`-Stage-Reste

- **kategorie:** LOW · **quelle:** Maintainability
- **pfad:** `Makefile:5` (Kopfkommentar „Test/arch-check laufen über
  `docker build --target`"), `Makefile:223` (`$(IMAGE):arch-check` im `clean`),
  `Makefile:24-27` (`NO_CACHE_FILTER_ARCH` — in DoD genannt)
- **befund:** Die DoD nennt Kopfkommentar und `NO_CACHE_FILTER_ARCH`, nicht aber
  die `clean`-Zeile; bleibt sie stehen, referenziert `clean` ein nie mehr
  gebautes Image-Tag (toter Rest, kein Fehler).
- **verifizierbar:** ja — `grep -n arch-check Makefile` nach Umsetzung.

### LOW-2 — §2-Layer-Skizze: `model`/`rules` brauchen explizite `role: domain`, und unter `core/` liegen vier Pakete, nicht drei

- **kategorie:** LOW · **quelle:** a-check AC-FA-RULE-006 (Rollen-Inferenz kennt
  nur `core`/`ports`/`adapters`/`app`) · **pfad:**
  `docs/plan/planning/in-progress/slice-058-arch-check-via-a-check.md:46-48`
- **befund:** Werden `model`/`rules`/`app` eigene Layer (nötig für R6-Kanten),
  greift für `model`/`rules` keine Namens-Inferenz — ohne explizites
  `role: domain` dispatchen die Reinheits-Regeln dort nicht (R1-Deckungsverlust;
  eine korrekt platzierte R1-Probe würde es fangen — Proben-Platzierung siehe
  MEDIUM-2). Zudem existiert als viertes Paket `internal/hexagon/core/coretest`
  (Test-Helfer), das eine Layer-Zuordnung braucht; „die drei Kern-Pakete" ist
  layout-unvollständig.
- **verifizierbar:** ja — R1-Mutations-Probe in `core/model` nach Umsetzung.

### INFO-1 — Offene Frage §4(a) ist nach Quellen-Lage positiv beantwortbar: `net/url` flaggt nicht

- **kategorie:** INFO · **quelle:** a-check SPEC-RULE-001 · **pfad:**
  `docs/plan/planning/in-progress/slice-058-arch-check-via-a-check.md:96-98`
- **befund:** `core-impurity` flaggt nur Imports, die auf einen
  `app`/`port`/`adapter`-Layer **oder** ein `tech`-Muster auflösen; `net/url`
  löst auf keinen Layer auf und bleibt unauffällig, solange die
  `net`-Tech-Muster regex-verankert sind (RE2 ohne Lookahead — die verbotene
  `net`-Familie enumerieren statt Substring `net`). Die Sonderfall-Sorge (a)
  reduziert sich damit auf Pattern-Disziplin in der `.a-check.yml`.

---

## Negativbefunde (geprüft, ohne Befund)

- **Kein-Lastenheft-CR-Begründung:** schlüssig — reine Gate-Mechanik ohne
  Produkt-Delta, Präzedenz slice-039/slice-055 trägt; die `DC-QA-03`-Bindung des
  Gates bleibt deklariert. Ohne Befund.
- **Roadmap-/Planning-Konsistenz:** welle-47 §Aktuelle Welle ⟺ slice-058 in
  `in-progress/`; der §Nächste-Wellen-Zeiger ist konsistent auf welle-47
  umgeschrieben; MR-013-Pfad für die Closure benannt. Ohne Befund.
- **ADR-0029-Form:** Teil-Supersede-Konstruktion (Regeln R1–R6 bleiben, nur
  Mechanik) analog ADR-0026/0027; Geschichte-/Index-Annotation an ADR-0005/0012
  erst bei Closure (slice-053-R1-F-4-Lehre eingehalten); MR-014-Hausstil;
  Status Proposed → Korrekturen aus diesem Review sind regelkonform einarbeitbar.
  Ohne Befund.
- **Tombstone-Mechanik:** fünfter `codepaths.ignore-refs`-Eintrag konsistent mit
  ADR-0025; die vier Bestands-Tombstones in `.d-check.yml:89` verifiziert; die
  Atomaritäts-Anforderung (Umbau+`git rm`+Tombstone in einem Commit) ist als
  Risiko korrekt erfasst. Ohne Befund.
- **Pin-/Hermetik-Politik:** Fragment ist digest-gepinnt (v0.6.0,
  `sha256:b349a150…`), Lauf `--network none` + read-only-Mount; Pull-am-Pin =
  Setup wie semgrep (ADR-0010/0011-konform); `make versions`-Ausweis in DoD.
  Ohne Befund.
- **Referenz-Richtung (SDP)/Marker-Ehrlichkeit:** keine Provenance-Marker im
  Gegenstand; ADR→Slice-Nennung nur in der Geschichte (Provenance). Ohne Befund.

---

## Kategorie-Summary

| Kategorie | Anzahl |
| --- | --- |
| HIGH | 2 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 1 |

## Verdikt

**NACHBESSERN (Plan-Ebene, vor Implementierungsbeginn).** Der Plan ist in
Prozess, Bindepunkten und Rückbau-Mechanik sauber; die beiden HIGHs betreffen
die Substanz der Regel-Übersetzung: die §2-/ADR-Abbildung ist gegen das
abzulösende Skript unvollständig (R2b/R3/R1-Bannliste), und das
`_test.go`-Scope-Delta macht den sauberen Baum vorhersagbar rot — beide sind
heute quellen-belegbar und gehören vor die Umstellung (Regel-Übersetzung
korrigieren, Delta-Liste/CR-Kandidaten an a-check in §4 aufnehmen), nicht in
die Proben-Phase. MEDIUM-1 bis MEDIUM-3 sind Plan-Präzisierungen
(composition_root-Delta benennen, Proben-Matrix je Verbotszweig,
gate-consistency-Verdrahtung festlegen). ADR-0029 ist Proposed und darf
entsprechend nachgezogen werden.
