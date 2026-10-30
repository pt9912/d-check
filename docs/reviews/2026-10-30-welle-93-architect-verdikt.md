# Architect-Verdikt — welle-93 (Modul `metrics`), Plan gegen Entscheidungslage

**Rolle:** Architect (Modul 8), frischer Kontext, **vor** Code.
**Gegenstand:** [welle-93](../plan/planning/welle-93-metrics-modul.md) und ihre drei Slice-Pläne
[slice-282](../plan/planning/next/slice-282-metrics-anforderung-und-adr.md),
[slice-283](../plan/planning/open/slice-283-metrics-modulkern-formel-und-grund.md),
[slice-284](../plan/planning/open/slice-284-metrics-messung-aus-git.md) (alle uncommittet).
**Anlass:** GitHub-Issue 13 (opt-in-Modul `metrics`).
**Datum:** 2026-10-30.
**Geändert wurde nichts:** kein Plan, keine ADR; dies ist nur das Verdikt.

## Gesamturteil

Keine bestehende Entscheidung wird von den Plänen oder den fünf Entscheiden verletzt.
Die Richtung (opt-in-Modul, VCS-Port, fail-closed, reine Arithmetik, kein Schreiben) deckt sich mit
[ADR-0024](../plan/adr/0024-vcs-immutable-gate.md), [ADR-0027](../plan/adr/0027-commits-traceability-modul.md)
und [ADR-0030](../plan/adr/0030-tracked-referenz-ziele.md). Es gibt aber **drei Entscheidungen, die vor
dem Code fehlen** (F-2, F-3, F-4), **zwei Plan-Fehler, die `make gates` rot machen oder falsch verorten**
(F-1, F-5) und **eine Lücke im ADR-Lebenslauf** (F-6). Alle Folge-Entscheidungen gehören in die **eine** ADR
aus slice-282, die ohnehin neu entsteht — eine zweite ADR ist nicht nötig.

Verdikt-Kürzel (Modul 8, Konflikt-Pfad): **V1** = Entscheidung gilt, Plan hat falsch behauptet/ausgelassen ·
**V2** = Folge-Entscheidung nötig · **V3** = Lockerung legitim, aber undokumentiert.

## Befunde

### F-1 — Spiegel-Tests widersprechen der Slice-Reihenfolge und der Release-Prep-Abgrenzung — **V1**

- Die Spiegel-Tests halten Lastenheft, README(s), Benutzerhandbuch-Tabelle und `operations.md`-Optionen-Tabelle
  in **beide** Richtungen gegen `model.ValidModules()`
  ([registry_mirror_test.go](../../internal/hexagon/core/app/registry_mirror_test.go) Z. 57–62 „kein bekanntes
  Modul (verwaist?)", Z. 108–160 die sieben Fundorte).
- [slice-282](../plan/planning/next/slice-282-metrics-anforderung-und-adr.md) Z. 140 trägt `metrics` in die
  Glossar-Zeile „Regelmodul" ein, aber die Registrierung in `ValidModules()` liegt erst in
  [slice-283](../plan/planning/open/slice-283-metrics-modulkern-formel-und-grund.md) Z. 72. Ab slice-282 wäre
  `TestModulRegistrySpiegel_LastenheftGlossar` rot (verwaister Eintrag). Der zweite Lastenheft-Fundort, der Modulsatz in
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) (Lastenheft Z. 111), fehlt im
  Plan 282 ganz.
- slice-283 widerspricht sich: §1 Z. 53–56 nennt die Spiegel-Tests als Ausnahme von der Release-Prep, DoD Z. 87–88
  sagt „Benutzerhandbuch und `operations.md` **nicht**". Mit `ValidModules()` + `metrics` werden aber genau die
  Handbuch-Tabelle (Test Z. 158) und die `operations.md`-Optionen-Tabelle (Test Z. 128) und beide READMEs rot.
- **Verdikt:** Die Gate-Lage gilt; die Pläne behaupten falsch. Korrektur: Glossar- und CLI-002-Eintrag wandern von
  slice-282 nach slice-283 (oder 282 benennt `metrics` dort nicht), und slice-283-DoD nimmt die vier tabellen-/listenhaltigen
  Fundorte als Spiegel-Ausnahme ausdrücklich auf (nur die Modulzeile, keine Prosa).

### F-2 — Block-Syntax „YAML" und Kern ohne YAML-Bibliothek — **V2** (Entscheid 6 in der slice-282-ADR)

- Entscheid 1 empfiehlt einen Fenced-**YAML**-Block ([slice-282](../plan/planning/next/slice-282-metrics-anforderung-und-adr.md) Z. 51).
  slice-283 §3 Z. 103 liest ihn in einer Datei des Regel-Pakets („Fence mit konfiguriertem Info-String") — ohne Port, ohne Adapter.
- Der Kern darf keine YAML-Bibliothek tragen: [architecture.md](../../spec/architecture.md) Z. 80 (ARC-002, Verbotsspalte
  „YAML-Bibliothek"); durchgesetzt in [.a-check.yml](../../.a-check.yml) (Eintrag zu `gopkg.in/yaml.v3`, nur
  `configyaml`/`report`/`workflowyaml`). Präzedenz für den Ausweg: der Port in
  [workflow.go](../../internal/hexagon/port/driven/workflow.go) Z. 1–7 mit
  [ADR-0072](../plan/adr/0072-workflows-modul.md) und [ADR-0009](../plan/adr/0009-yaml-im-report-adapter.md).
- Zwei tragfähige Wege, die der Auftraggeber wählt: (a) Parse-Port + Adapter nach `WorkflowParser`-Muster (macht slice-283 größer, vierte Schicht, s. F-7),
  (b) bewusst eingeschränkte Zeilen-Syntax (`schlüssel: ganzzahl`), vom Kern selbst gelesen, ohne YAML-Bibliothek. Der Plan nennt keinen von beiden.
- **Verdikt:** Entscheidung fehlt; vor dem Code in der ADR als eigener Entscheid mit Alternativen-Tabelle fassen.
  slice-283 §3 bleibt bis dahin nicht umsetzbar, ohne [architecture.md](../../spec/architecture.md) zu verletzen.

### F-3 — Fence-Lexik: gescopte Fence-Ausnahme fehlt — **V2** (Entscheid 7 in derselben ADR)

- Die gemeinsame Vorverarbeitung entfernt Fenced-Code für **alle** Module; ein Modul, das in Fences liest, braucht eine gescopte
  Ausnahme per ADR ([ADR-0018](../plan/adr/0018-diagram-fence-ausnahme.md) Kontext/Entscheidung 1,
  [ADR-0019](../plan/adr/0019-versions-pin-fence-ausnahme.md), [ADR-0020](../plan/adr/0020-content-pin-fence-ausnahme.md)).
- [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 1 und 2: Eine **eigene** Antwort auf eine geteilte
  Lexik-Frage („wo beginnt/endet ein Fence?") ist ein Defekt; zulässig ist nur eine **andere Frage**, per ADR gescopt.
  [ADR-0042](../plan/adr/0042-markdown-lexik-folgt-commonmark.md) Entscheidung 2 (Infozeile mit Backtick ist kein Fence) ist Teil der geteilten Antwort.
- slice-283 erkennt das Risiko (§6 Z. 167–169) und listet Fälle in der Negativliste (§3 Z. 121–128), trifft aber **keine** Entscheidung, ob das Modul
  einen eigenen Fence-Automaten führt (dann: ADR-0018-Form, Isolation zu `spans`) oder die geteilte Fence-Erkennung wiederverwendet (ADR-0054-Form).
- **Verdikt:** Entscheidung fehlt; die ADR benennt die Ausnahme, ihren Scope (nur Info-String = Blockname, nur Plan-Verzeichnis) und dass die Fence-Grenzen
  aus der geteilten Lexik übernommen werden. **Teil-Negativ, akzeptiert:** „Plan ohne Block ⇒ kein Befund" ist durch das Issue gedeckt (Prüfung „pro Plan mit Block").
  Ein **ungeschlossener** Fence darf aber nicht still als „ohne Block" gelten, sondern ist `metrics-block-invalid` — das gehört in die Anforderung, nicht in einen Folge-Slice.

### F-4 — Port-Erweiterung, Alternativen, Spec-Schärfung fehlen im ADR-Umfang — **V2** (Entscheid 8 in derselben ADR)

- Der ADR-Umfang in [slice-282](../plan/planning/next/slice-282-metrics-anforderung-und-adr.md) Z. 116 deckt „Entscheide 1–4". Die neue VCS-Port-Methode
  ([slice-284](../plan/planning/open/slice-284-metrics-messung-aus-git.md) Z. 69, 102) ist keiner davon. Präzedenz: Port-Erweiterungen stehen in der Modul-ADR
  ([ADR-0027](../plan/adr/0027-commits-traceability-modul.md) Schärft-Feld, [ADR-0030](../plan/adr/0030-tracked-referenz-ziele.md) „Port-Erweiterungs-Präzedenz").
- Eine Alternative fehlt: Der Kern könnte aus den vorhandenen Methoden `AllPaths`/`FileAt`
  ([vcs.go](../../internal/hexagon/port/driven/vcs.go) Z. 33, 35) selbst zählen; dann braucht der Port keine neue Methode, der Kern aber einen Zeilen-Diff.
  Die Wahl ist nicht begründet.
- [architecture.md](../../spec/architecture.md) Z. 83 (ARC-005 Zuständigkeit: „Modul `vcs`, `commits`"), Z. 99 (ARC-010) und Z. 153 (Fehlerzeile) beschreiben, was der Adapter leistet;
  slice-282 „Berührte Spec-Stellen" (Z. 20–22) nennt die Architektur-Sicht nicht, obwohl ADR-0027 sie bei der letzten Port-Erweiterung ausdrücklich schärfte. (`tracked` hat die Zeilen
  nicht nachgezogen — die Drift ist Bestand, kein Grund, sie zu wiederholen.) Die Sicht bleibt dabei sprach-/meilensteinfrei ([AGENTS.md](../../AGENTS.md) §3.4).
- Der Zeilen-Diff kommt aus go-git; die Bibliothek `sergi/go-diff` steht schon als indirekte Abhängigkeit in [go.mod](../../go.mod) Z. 23. **Nicht geprüft:** ob die
  gepinnte Version die benötigte Zeilen-Statistik samt Binär-Erkennung liefert — das ist genau die Fixture-Gleichheit aus slice-284 DoD Z. 69–74.
- **Verdikt:** Entscheidung fehlt; die ADR führt Port-Methode (Signatur, Fehlerverhalten, Alternative „Kern zählt selbst"), Verdrahtung (F-5) und die
  Schärfung der Architektur-Sicht. Die Gleichheit mit dem Kommando bleibt, wie geplant, eine in slice-284 belegte Zusage.

### F-5 — Falsch verortete Verdrahtung in slice-284 — **V1**

- slice-284 §3 Z. 105 nennt `internal/hexagon/core/app/` als „Aufrufer von `RunWithVCS`". Der Aufrufer ist
  [cli.go](../../internal/adapter/driving/cli/cli.go) Z. 656; die Modul-Verteilung liegt in
  [run.go](../../internal/hexagon/core/rules/run.go) Z. 144–157 (`runPostPasses`), und das Paket `app` ruft `RunWithVCS` nirgends.
- Entscheidender: `resolveVCS` ([cli.go](../../internal/adapter/driving/cli/cli.go) Z. 715–745) kennt zwei Fälle — Adapter **mit** Range (`vcs`/`commits`) und Adapter **ohne**
  Range (`tracked`). `metrics` ist ein dritter: Adapter ohne `--range`, Basis aus dem Block, Head fest `HEAD` (Entscheid 2). Ohne Eingriff in `gitModuleActive`/`resolveVCS`
  und die Meldung Z. 707 („…(vcs/commits) braucht --range…") bekäme `metrics` entweder keinen Port (nil ⇒ Absturz statt Exit 2) oder einen unbegründeten Zwang zu `--range`.
  slice-284 führt in `internal/adapter/driving/cli/` nur den Image-Test.
- **Verdikt:** Plan hat falsch verortet; Korrektur der Tabelle (`cli.go` `resolveVCS` + `run.go` `runPostPasses`, nicht `app/`) und der Verdrahtungs-Satz in die ADR (F-4).
  Die Zusage „`--range` bleibt unberührt" (slice-284 Z. 56–58) ist damit konkret zu prüfen: ein Lauf `--enable vcs,metrics --range …` benutzt zwei verschiedene Köpfe (Range-Head vs. `HEAD`).

### F-6 — ADR-Status über drei Slices ungeregelt — **V1**

- slice-284 Z. 22 verspricht, die `numstat`-Gleichheit „belegt oder korrigiert" zu haben; slice-282/283 nennen die ADR „aktive Entscheidungsgrundlage"
  ([slice-283](../plan/planning/open/slice-283-metrics-modulkern-formel-und-grund.md) Z. 18). Wann sie `Accepted` wird, sagt kein Plan.
- Eine `Accepted`-ADR wird nicht inhaltlich geändert ([AGENTS.md](../../AGENTS.md) §3.5, [ADR-0016](../plan/adr/0016-adr-immutable-gate.md)): Fällt die Gleichheit in
  slice-284 anders aus als in Entscheid 4 gesetzt, wäre nach vorzeitigem `Accepted` eine Folge-ADR mit `supersedes` nötig, wo ein Edit am Entwurf genügt hätte.
  Bestandsmuster: Annahme mit der Closure des **umsetzenden** Slice ([ADR-0030](../plan/adr/0030-tracked-referenz-ziele.md) Geschichte, [ADR-0024](../plan/adr/0024-vcs-immutable-gate.md) Geschichte).
- **Verdikt:** Plan ergänzen: ADR bleibt `Proposed` bis zur Closure von slice-284; der Status-Übergang ist DoD-Punkt dort (erlaubt für `make adr-check`). Das ist kein Gate-Thema, nur eine Plan-Lücke.

### F-7 — Schnitt in drei Slices gegen die Hexagon-Schichten — **bestätigt, mit akzeptiertem Negativ**

- Der Schnitt folgt dem Lieferwert (hermetische Formel/Grund vs. Messung aus git), nicht den Schichten — so verlangt es Baseline-Regelwerk
  [modul-05](../../.harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md) §Ziel-Form: Slice („Schnitt nach Lieferwert"). Ein Port-/Adapter-Vorab-Slice wäre ein
  Schicht-Zombie ohne Konsumenten. Reihenfolge 282 → 283 → 284 ist fachlich zwingend ([welle-93](../plan/planning/welle-93-metrics-modul.md) §5).
- Beide Modul-Slices berühren mehr als zwei Schichten (283: Kern, `configyaml`-Adapter, CLI-Gerüst; 284: Port, git-Adapter, Kern, CLI). Das ist bei jedem opt-in-Modul der Fall
  (Präzedenz ADR-0027/ADR-0030) und in beiden Plänen als Rückführungs-Trigger vorab benannt (283 §4 Z. 141–144 mit [`MR-066`](../../harness/conventions.md#mr-066)).
  **Akzeptiert** mit dieser Begründung; es wird kein eigener Träger angelegt. Wählt der Auftraggeber in F-2 Weg (a), wächst slice-283 um Port und Adapter — dann ist der
  benannte Rückführungs-Trigger der Normalfall, und die Block-Lexik sollte vorab als eigener Slice geschnitten werden.
- Zwischenzustand nach slice-283 (Modul aktiv, git-Hälfte fehlt, Exit-2-Zusage ohne `.git` noch nicht eingelöst) ist auf die offene Welle beschränkt; kein Release dazwischen
  ([welle-93](../plan/planning/welle-93-metrics-modul.md) §6). Akzeptiertes Negativ.

## Bestätigte Bezüge (kein Befund)

- **Hermetik und Determinismus:** Ganzzahl-Arithmetik, feste Rundung, kein Gleitkomma — deckt [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) (Lastenheft Z. 4193–4196).
- **Schreibfreiheit und Netzlosigkeit:** nur lesender git-Zugriff über den Port, wie [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  (Z. 2188) und [`DC-FA-COMMITS-001`](../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in) (Z. 2325);
  deckt [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (Z. 4198–4201). `go-git` bleibt einzige git-Tür (R2b in [.a-check.yml](../../.a-check.yml)); der
  Port gibt die fail-closed-Haltung vor ([vcs.go](../../internal/hexagon/port/driven/vcs.go) Z. 18–21, 40–42).
- **Fail-closed ohne `.git`, nicht auflösbare Basis ⇒ Exit 2:** gedeckt durch [architecture.md](../../spec/architecture.md) Z. 153 und [ADR-0030](../plan/adr/0030-tracked-referenz-ziele.md) („Fail-closed"). Dass **eine** unauflösbare
  Basis den ganzen Lauf abbricht (slice-284 Z. 79–80), ist konsequent und gehört als Satz in die ADR (Entscheid 8). Nach dem Gate-Netzlos-Test gilt: `metrics` bleibt außerhalb der Live-`modules:`-Liste (slice-283 Z. 57–63),
  daher greift die Klassifikations-Richtung in [gate_consistency_test.go](../../internal/adapter/driven/configyaml/gate_consistency_test.go) Z. 60–71 nicht; sie wäre bei einer späteren Aufnahme neu zu entscheiden.
- **Entscheid 3 (lineare Formel, kein Parser) und Entscheid 5 (Grund-Codes):** keine bestehende ADR berührt; Spezifikations-Sache in slice-282. Ob die Fachlogik des Anlass-Repos „Zustände" verlangt, entscheidet dessen Gegenprobe (slice-282 Z. 61–62, Rückführung Z. 169–171) — nicht vorab zu entscheiden.
- **Entscheid 4, Rename-Erkennung aus:** Das Issue nennt das Kommando wörtlich (`git diff --numstat <basis>..HEAD`). Ob dessen Standard-Verhalten (Rename-Erkennung an oder aus) mit `--no-renames` zusammenfällt,
  habe ich **nicht** gemessen. Die Empfehlung der Pläne ist deshalb zu Recht nur Empfehlung; maßgeblich ist, was das Programm des Anlass-Repos tut (Paritäts-Probe in [welle-93](../plan/planning/welle-93-metrics-modul.md) §3).
- **Abwärts-Sperre:** Die Pläne lassen die Spec-Straten nie auf Slices oder ADRs zeigen (slice-282 Z. 23–25); `DC-*` entstehen nur im Lastenheft ([AGENTS.md](../../AGENTS.md) §5 Zeile 3).

## Nicht geprüft

- Der Inhalt der Pläne des Anlass-Repos (liegen nicht vor) und damit, ob eine lineare Formel trägt.
- Das Verhalten von `git diff --numstat` und von go-gits Zeilen-Statistik an Randformen (Binärdateien, Zeilenenden, Renames) — das ist der Messauftrag von slice-284.
- Ob `make gates` heute gegen die Pläne grün ist; geprüft wurde nur gegen die Entscheidungs- und Gate-Lage per Lesen.
- Unveränderte ADRs außerhalb der genannten (u. a. die gesamte Linkprüfungs-/`structure`-Familie) und die Prosa-Spiegel außerhalb der sieben Fundorte der Spiegel-Tests.
- Die Reihenfolge-Aussage zu `FOCUS_DISABLE` im Makefile (slice-283 Z. 57–63): nur gelesen, der Test `TestFocusDisable_DecktDCheckYmlModules` nicht gefahren.

## Zusammenfassung

| Befund | Verdikt | Träger der Korrektur |
|---|---|---|
| F-1 Spiegel-Reihenfolge/Release-Prep | V1 Plan falsch | Plan 282/283 |
| F-2 YAML im Kern | V2 Folge-Entscheid | ADR aus slice-282 (Entscheid 6) |
| F-3 Fence-Ausnahme | V2 Folge-Entscheid | ADR aus slice-282 (Entscheid 7) |
| F-4 Port-Methode, Alternative, Architektur-Sicht | V2 Folge-Entscheid | ADR aus slice-282 (Entscheid 8), Spec-Stellen in 282 |
| F-5 Verdrahtung `cli.go`/`run.go` | V1 Plan falsch | Plan 284 |
| F-6 ADR-Status `Proposed` bis 284 | V1 Plan-Lücke | Plan 282/284 |
| F-7 Dreier-Schnitt | bestätigt, Negativ akzeptiert | — |

Kein Befund fällt unter V3 (undokumentierte Lockerung): keine Schwelle, kein Gate wird gesenkt.
