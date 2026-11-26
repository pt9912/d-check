# Slice slice-297: Die Spezifikation heißt `specification.md` — reiner Rename, Link-Nachzug, Produkt-Pfad und zwei Tests

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-98-baseline-v700`](../welle-98-baseline-v700.md).

**Bezug:** [`DC-FA-CLI-013`](../../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen)
(der Vertrag bleibt im Wortlaut; die Kopfzeile von `--manual` trägt den Pfad),
[`ADR-0107`](../../adr/0107-mitgelieferte-dokumente-im-modul-root.md) (`manual.go` bettet die Spezifikation ein),
[`ADR-0103`](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) und
[`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) (der ADR-Gate lässt den Link-Nachzug auf eine umbenannte Datei
durch), [`ADR-0025`](../../adr/0025-codepaths-ignore-refs.md) (das Tombstone-Register von `codepaths`),
[`MR-069`](../../../../harness/conventions.md#mr-069) (`ignore-refs` als deklarierte Gate-Senkung),
[`MR-070`](../../../../harness/conventions.md#mr-070) (Frozen-Klassen vor mechanischer Ersetzung),
[`MR-032`](../../../../harness/conventions.md#mr-032) (Lastenheft-Historie, Zeile 4762),
[`AGENTS.md`](../../../../AGENTS.md) §3.3 (Move und Inhalt in zwei Commits). Anlass: Entscheid 1 der
[Welle](../welle-98-baseline-v700.md) (Namensregel, „ab jetzt, Bestand bleibt“).

**Berührte Spec-Stellen:** die Datei `spec/spezifikation.md` selbst (ihr Name und ihre 161 Selbst-Verweise) ·
[`DC-FA-CLI-013.a`](../../../../spec/specification.md#dc-fa-cli-013a--handbuch-und-spezifikation-ausgeben---manual)
Schritt 5 (die Kopfzeile `==> <pfad>:<zeile>` — der Wortlaut bleibt, die Ausgabe trägt den neuen Pfad) ·
`spec/architecture.md` Z. 26 (der Link auf die Spezifikation) · `spec/lastenheft.md` Z. 4762 (Auftrag des Auftraggebers,
[Verdikt](../../../reviews/2026-11-26-welle-98-architect-verdikt.md) Frage 6). Der Verweis zeigt **aufwärts**: Die
Spec nennt diesen Slice nie.

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-11-25.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Spezifikation heißt `specification.md` (unter `spec/`): ein **reiner** Rename, jeder lebende Verweis und jedes
Link-Ziel nachgezogen, `manual.go` und die zwei Tests, die den Namen tragen, folgen, und die eingefrorenen Nennungen
stehen unter einem benannten Ventil — `make doc-check` ohne Befund, `make adr-check` über die ganze Spanne Exit 0,
`make test` grün.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Bestand an anderen Namen.** `welle-<NN>`, `archiv.zip`, `seit welle-` bleiben; die **neuen** Namen der Wellen,
  Archive und Anker gehören in [`slice-298`](../open/slice-298-neue-namen-wave-archive-anker.md) — der Rename kommt ohne
  sie aus und sie ohne den Rename (zwei Hälften, zwei Slices; sonst vier Liefer-Punkte).
- **Keine Änderung des Wortlauts von [`DC-FA-CLI-013`](../../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen)** und der Spezifikation §[`DC-FA-CLI-013.a`](../../../../spec/specification.md#dc-fa-cli-013a--handbuch-und-spezifikation-ausgeben---manual), kein Standardwert, kein
  Grund-Code, keine Logik. Die Ausgabe von `--manual` ändert sich als **Folge** (Kopfzeile, Link-Ziele der Spezifikation,
  eine Handbuch-Zeile); die Klassen der Abweichung sind in der Welle (d3) benannt und werden dort gemessen.
- **Keine Konsumenten-Beispiele.** `config_template.go` Z. 56–57, `suggest.go` Z. 379–380 und
  `benutzerhandbuch.md` Z. 372–373 zeigen, wie **ein Konsument** seine Spec-Straten nennt, nicht diese Datei; sie
  bleiben in jedem Fall (die Vorlage des Regelwerks nennt dort künftig `specification.md` — das ist Sache der
  [Sichtung](../open/slice-296-baseline-v700-neuerungen-sichtung.md), Frage B).
- **Keine Kommentarzeilen** (Verdikt Frage 7: entschieden, sie bleiben): 90 Zeilen in 33 Go-Dateien, dazu je eine in `Dockerfile`
  und `Makefile`, zwei in `tools/` sind Rang-Zeiger ohne Wirkung; der Nachzug ist die größere Antwort und ein eigener Akt.
- **Das Lastenheft, außer Zeile 4762** (Auftrag des Auftraggebers, Verdikt Frage 6). Jede weitere Änderung dort ist nicht Teil des Slice.
- **Der Pin-Text der Spezifikation (Zeile 8)** gehört dem Bump ([`slice-295`](../done/slice-295-baseline-v700-bump.md)).
- **CHANGELOG, READMEs, Handbuch-Kopf, Release** — Release-Prep ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17); der
  Hinweis auf die geänderte Kopfzeile von `--manual` steht in der Closure-Notiz, nicht im Feature-Commit. **Ausnahme,
  weil sie Links sind:** die zwei Links der READMEs und die drei von `operations.md` müssen im selben Zug fallen, sonst
  meldet `links`; sie sind Link-Ziele, kein Handbuch-Kopf.
- **Eingefrorene Dateien über ihre Link-Ziele hinaus** — `done/`, `docs/reviews/`, `CHANGELOG.md`, `conventions/done/`,
  der Kern der ADRs nennen den alten Namen als Vergangenheits-Aussage und bleiben byte-stabil.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Rename und Link-Ziele in einem Commit:** `specification.md` (unter `spec/`) entsteht per **reinem** `git mv` (R100); derselbe
      Commit trägt die Link-Ziel-Nachzüge der 64 ADRs (und, nach Frage 5, der 39 eingefrorenen Nicht-ADR-Dateien).
      `git diff --cached --name-status -M` vor dem Commit nennt genau **eine** `R100`-Zeile und sonst nur `M`;
      `make adr-check STAGED=1` endet mit Exit 0, und `make adr-check RANGE=<Vor-Bump>..HEAD` ebenfalls.
- [x] **Inhalt und Code:** Die Spezifikation trägt ihre Selbst-Verweise auf den neuen Namen (161 Zeilen), jede lebende
      Fundstelle (die Klassen und Zahlen in §3) nennt `specification.md`, `.d-check.yml`, `manual.go`, die zwei Tests und
      `tools/image-test.sh` folgen, und der **eine** `codepaths.ignore-refs`-Eintrag für die eingefrorenen
      Inline-Code-Nennungen steht mit seiner Grenze im Kommentar; `make doc-check` meldet 0 Befunde, `make test` ist grün,
      und beide Tests laufen **ohne** den Nachzug aus dem richtigen Grund rot (die Meldung ist gelesen, §3).
- [x] **Lastenheft** (Auftrag, Verdikt Frage 6): Zeile 4762 nennt `specification.md` (unter `spec/`), `**Version:**` 0.115.1 → 0.115.2,
      **eine** Historie-Zeile in der Form der übrigen („Redaktionell: …; keine inhaltliche Änderung“, Verweis `—`) nach
      [`MR-032`](../../../../harness/conventions.md#mr-032); kein `DC-*` ändert sich, kein weiterer Satz.
- [x] `make gates` grün; `make image-test` grün (die Zeile 163 von `tools/image-test.sh` prüft die Kopfzeile wörtlich).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

**Die Zählung nennt ihre Form** ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 14): gezählt wird die **git-getrackte
Datei bzw. Zeile**, die die Zeichenkette `spezifikation.md` trägt (`git grep`), eingeteilt nach der **Eigenschaft** der
Datei ([`MR-070`](../../../../harness/conventions.md#mr-070): eingefroren ist, was ein Beleg, eine `Accepted`-ADR, ein
Lauf-Bericht ist — nicht, was in `done/` liegt), nicht nach dem Verzeichnis. Stand `590dde3d`; der Implementer
**wiederholt** die Zählung vor der ersten Änderung. Die Zahlen stammen aus dem Probelauf in einem Wegwerf-Klon
([welle-98 §1](../welle-98-baseline-v700.md)), dort auch ihre Messung am Gate.

| Datei / Komponente | Änderungs-Art | Begründung und Zahl |
|---|---|---|
| `spec/spezifikation.md` → `specification.md` (unter `spec/`) | move | **eine** Datei, reiner `git mv`, R100 — **nicht baubar** allein (`go:embed` in `manual.go`), deshalb beide Commits in **einen** Push |
| 64 `Accepted`-/`Proposed`-/`Superseded`-ADRs | update (**nur** Link-Ziel) | **119** Zeilen (`git grep -c '\](\([^)]*/\)\?spezifikation\.md' -- docs/plan/adr`); im **selben** Commit wie der Move — gemessen: zwei Commits geben **60** `core-drift-vcs` am Hook, ein Commit **0**. Die 8 ADRs, die den Namen nur im Inline-Code nennen, bleiben unberührt (immutabel, Ventil unten) |
| 39 eingefrorene Nicht-ADR-Dateien mit Link-Ziel (`done/` 23, `docs/reviews` 13, `docs/plan/cr` 2, `conventions/done` 1) | update (**nur** Link-Ziel; Frage 5, kleinste Antwort) | **80** Zeilen; im Move-Commit. **Wortgrenze:** zwei Dateinamen enden auf `spezifikation.md` (`slice-002-architektur-und-spezifikation.md`, `slice-114-spec-vergabe-spezifikation.md`) — der Tausch greift nur nach `/`, `(` oder Zeilenanfang |
| `specification.md` (unter `spec/`) | update | **161** Zeilen (158 als Link-Ziel auf sich selbst); zweiter Commit — **Inhaltsänderung nach dem Move** (§3.3) |
| 27 lebende Markdown-Dateien (die 29 Dateien mit 53 Zeilen ohne [`MR-080`](../../../../harness/conventions.md#mr-080), das dann in `conventions/done/` liegt, und ohne das Lastenheft, Frage 6) | update | 51 Zeilen tragen den Namen, **49** werden geändert (zwei sind Konsumenten-Beispiele): `harness/README.md` 7, `harness/sensors/*.md` 14 Dateien mit 24, `harness/conventions.md` und vier aktive MR-Dateien (025, 076, 0098, 006) mit 7, `harness/rules/docker-make-only.md` 2, `AGENTS.md` 2, `README.md` und `README.de.md` je 1, `docs/user/operations.md` 3, `docs/user/benutzerhandbuch.md` **1** (Z. 2792), `spec/architecture.md` 1. **Nicht** geändert: `benutzerhandbuch.md` Z. 372–373 (Konsumenten-Beispiel) |
| [`.d-check.yml`](../../../../.d-check.yml) | update | **6** Zeilen (Z. 428, 499, 503, 739, 827 tragen den Pfad, Z. 970 einen Kommentar) und **ein** Eintrag in `codepaths.ignore-refs` (Z. 657): `spec/spezifikation.md` — **repo-weit** wirksam (benannte Grenze, [`ADR-0025`](../../adr/0025-codepaths-ignore-refs.md)); gemessen deckt er die **104** Inline-Code-Nennungen in eingefrorenen Dateien (`done/` 47, ADRs 44, Register-Belege 5, `conventions/done` 5, `CHANGELOG.md` 2, `docs/plan/cr` 1) und bringt `doc-check` von 123 auf 0 Befunde |
| `manual.go` | update | **2** Zeilen: `SpezifikationPfad` und `//go:embed`; ohne sie bricht der Build |
| `internal/adapter/driving/cli/cli_manual_test.go` Z. 21 | update | erwartet die Kopfzeile `==> spec/spezifikation.md:`; **rot ohne den Nachzug** (gemessen: „Ausgabe ohne `==> spec/spezifikation.md:`“) |
| `internal/hexagon/core/app/diagnose_test.go` Z. 73 | update | liest die Spezifikation am Pfad; **rot ohne den Nachzug** (gemessen: „Spezifikation nicht lesbar (fail-closed)“) |
| `tools/image-test.sh` Z. 163 | update | `grep -q '^==> spec/spezifikation.md:'` — funktional; `make image-test` ist sein Sensor |
| `spec/lastenheft.md` Z. 4762 und Historie | update (Frage 6, Auftrag) | **1** Zeile + Version (0.115.1 → 0.115.2) + **1** Historie-Zeile; das Lastenheft ist Rang 1, die Änderung führt der Implementer im Auftrag des Auftraggebers |
| Kommentare: 90 Zeilen in 33 Go-Dateien, `Dockerfile` 1, `Makefile` 1, `tools/*.sh` 2 | **unverändert** (Frage 7) | `git grep -n 'spezifikation\.md' -- '*.go' Dockerfile Makefile tools`; benanntes Negativ im Verdikt: ohne Wirkung auf Verhalten, §3.7 für 94 Zeilen |
| `config_template.go` (2), `suggest.go` (2), `benutzerhandbuch.md` Z. 372–373 | **unverändert** | Konsumenten-Beispiele (§1) |

**Plan-Änderung (Implementer, vor dem Code; Stand `f44c6e65`, die Zählung ist wiederholt).** Die Wiederholung nennt dieselben
Zahlen wie die Tabelle (64 ADRs mit 119 Zeilen, 39 eingefrorene Dateien mit 80 Zeilen, 161 Selbst-Verweise, 49 geänderte lebende
Zeilen) und dazu **sechs Dateien mit je einem Link-Ziel, die keine Zeile der Tabelle trägt**. Ohne ihren Nachzug meldet `links`
`target-missing`; ihre Antwort ist der Nachzug, nicht das Ventil:

| Dateien | Zeilen | Eigenschaft | Kommando |
|---|---|---|---|
| vier Dateien unter `docs/plan/planning/observations/` (`BEO-ALL/schwelle-mit-mehreren-traegern-ohne-gleichhalter/evidence/slice-259.md`, `BEO-ALL/spezifikation-und-code-weichen-am-rand-ab/evidence/slice-289.md`, `BEO-ALL/grenze-der-scan-menge-nur-bei-einem-modul-benannt/evidence/slice-289.md`, `BEO-HARN/commit-msg-hook-sperrt-junges-repo/observation.md`) | 4 | Register-Belege und eine Beobachtung: der Link zeigte immer auf die **lebende** Datei (Frage 5, dieselbe Begründung wie bei den 39) | `git grep -c '\](\([^)]*/\)\?spezifikation\.md' -- docs/plan/planning/observations` |
| dieser Plan (Z. 27 und 58) und die Welle `welle-98-baseline-v700.md` (Z. 151) | 3 | lebende Pläne mit Link auf den Abschnitt der Spezifikation zu [`DC-FA-CLI-013`](../../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen); ihre Inline-Code-Nennungen des alten Namens beschreiben den Gegenstand und bleiben (Ventil) | `git grep -c '\](\([^)]*/\)\?spezifikation\.md' -- docs/plan/planning/in-progress docs/plan/planning/open docs/plan/planning/welle-98-baseline-v700.md` |

Alle sieben Zeilen: **nur** das Link-Ziel, nach der Wortgrenze (`](`, optional ein Pfad bis `/`, dann `spezifikation.md`). Sie reisen im
**zweiten** Commit, nicht im Move-Commit — dessen Menge (64 + 39 = 103 `M`) bleibt wie in §3. Die Summe aller Link-Ziele geht auf:
119 + 80 + 4 + 3 + 42 (lebende Markdown-Dateien) + 158 (Selbst-Verweise) = 406
(`git grep -h '\](\([^)]*/\)\?spezifikation\.md' | wc -l`).

**Weitere Berichtigungen der Tabelle, ohne Folge für den Umfang:**

- **Zeilennummern der Konfiguration** haben sich mit [`slice-295`](../done/slice-295-baseline-v700-bump.md) verschoben: die Pfad-Zeilen der
  [`.d-check.yml`](../../../../.d-check.yml) stehen bei 446, 517, 521, 757 und 845, der Kommentar bei 988, die Liste
  `codepaths.ignore-refs` bei 675 (`git grep -n 'spezifikation\.md' -- .d-check.yml`). Es bleiben **6** Zeilen und **ein** Eintrag.
- **Kommentarzeilen:** 90 Zeilen in 33 Dateien zählen die Quelldateien ohne Tests (so die Tabelle); mit den Testdateien sind es
  **95 in 38 Dateien** (`git grep -n 'spezifikation\.md' -- '*.go' | grep -v 'go:embed' | grep -c '^[^:]*:[0-9]*:[[:space:]]*//'`).
  Die fünf Zeilen in Testdateien bleiben aus demselben Grund wie die übrigen
  ([Verdikt](../../../reviews/2026-11-26-welle-98-architect-verdikt.md) Frage 7): Rang-Zeiger ohne Wirkung. Dazu tragen 6 Zeilen in
  zwei Testdateien (`cli_acceptance_test.go`, `matrix_test.go`) den Namen als **Fixture-Inhalt** eines erfundenen Repos — sie prüfen
  Module, nicht diese Datei, und bleiben.
- **Zerlegung feiner als in §3:** Commit 2 der Reihenfolge teilt sich nach Klasse in drei — (2) lebende Dateien, Konfiguration, Ventil,
  `manual.go`, Tests, `tools/image-test.sh` und die sechs Dateien oben; (3) die Selbst-Verweise der Spezifikation (Inhalt der bewegten
  Datei, [`AGENTS.md`](../../../../AGENTS.md) §3.3 Fall 1); (4) das Lastenheft (Rang 1, Auftrag). Die Zwischenstände (1) und (2) sind einzeln rot
  (Build beziehungsweise Selbst-Links), die Spitze des Push ist (4); alle Commits gehören in **einen** Push.
- **Messungen des Closure-Triggers (d) der Welle**, damit klar ist, wer sie trägt: **diesen Slice** tragen `make test` (die zwei Tests,
  je mit Mutation) und `make image-test` (die Kopfzeile, `tools/image-test.sh` Z. 163). **Die Welle** trägt am Endstand
  `make blackbox-probe REF=738d2bca` (der Commit vor dem ersten Bump-Commit) und den Vergleich von `--manual` je Begriff mit den vier
  Abweichungsklassen; sie sind dem Slice nicht aufgetragen (§3, Gates).

**Zweite Plan-Änderung (Implementer, vor dem Code):** drei Festlegungen, die die Tabelle offen ließ, und die Folge für die Welle.

- **Zwei Zeilen der Spezifikation bleiben: die Konsumenten-Beispiele.** Die Spezifikation trägt in ihrem Konfigurations-Beispiel
  (`matrix.classes`, Z. 233 und 234) dieselben zwei Zeilen wie `benutzerhandbuch.md` Z. 372–373, `config_template.go` und `suggest.go` —
  ein Konsument nennt seine Spec-Straten, nicht diese Datei (§1, „Keine Konsumenten-Beispiele“). Ein Beispiel, das an vier Stellen
  gespiegelt steht ([`MR-025`](../../../../harness/conventions.md#mr-025)), ändert sich an allen vier oder an keiner; die Sichtung entscheidet es
  ([`slice-296`](../open/slice-296-baseline-v700-neuerungen-sichtung.md), Frage B). Damit trägt die Spezifikation **159** geänderte Zeilen
  (158 Link-Ziele und die Prosa-Zeile 920, die ihren eigenen Pfad nennt), nicht 161; die Zahl 161 zählt die Zeichenkette
  (`grep -c 'spezifikation\.md' <die Spezifikation>`) und nicht den Gegenstand. Danach nennt die Datei den alten Namen **2**-mal.
- **Das Ventil deckt 119 Befunde, nicht 104.** Gemessen am Arbeitsbaum ohne den `ignore-refs`-Eintrag: 119 Befunde `codepath-missing`, mit dem
  Eintrag 0 (`make doc-check`). Davon 109 in eingefrorenen Dateien (`done/` 51, ADRs 44, Register-Belege 5, `conventions/done` 6, `CHANGELOG.md` 2,
  `docs/plan/cr` 1) und 10 in zwei lebenden Klassen, die den alten Namen als **Gegenstand** nennen: die Pläne dieser Welle (9) und der Eintrag
  [`MR-081`](../../../../harness/conventions.md#mr-081) (1, Vergangenheits-Aussage über die Pin-Zeilen). Das Ventil ist deshalb auch für sie repo-weit wirksam; der
  Kommentar am Eintrag nennt es.
- **Folge für (d3) der Welle.** Die Zeile 920 der Spezifikation (der Pfad in der Prosa zu `--manual`) steht in der Ausgabe von `--manual` zu den Begriffen
  `manual` und `Spezifikation` (am Vorher-Image gemessen: je 1 Treffer von `spec/spezifikation.md`) und wechselt mit. Sie ist **keine** der vier
  Klassen von (d3) (Kopfzeile · Link-Ziel · Pin · Handbuch-Zeile) und wäre dort ein Befund; die Welle führt sie als fünfte Klasse
  (*ein Pfad in der Prosa der Spezifikation, Z. 920*), die Messung bleibt bei ihr. Der Slice ändert die Zeile, weil sie sonst eine Datei nennt, die es nicht gibt.

**Reihenfolge und Commit-Zerlegung** ([`AGENTS.md`](../../../../AGENTS.md) §3.3, [`BEO-ALL/path-scoped-commit-carries-staged-rest`](../observations/BEO-ALL/path-scoped-commit-carries-staged-rest/observation.md)):

1. **Vorbedingung, kein Commit dieses Slice:** das Architect-Verdikt zu den Fragen der Welle liegt vor,
   [`slice-295`](../done/slice-295-baseline-v700-bump.md) ist `done/`.
2. **Zählung wiederholen** (die Form oben), Frozen-Klassen **vorher** listen, bevor ein Tausch fällt.
3. **Arbeitsbaum vollständig vorbereiten, dann stufenweise stagen:** Alle Änderungen liegen im Arbeitsbaum, weil der
   `pre-commit`-Hook den **Arbeitsbaum** prüft und `make doc-check` und der Image-Build (`go:embed`) sonst am
   Zwischenstand brechen. **Commit 1:** `git mv` plus die Link-Ziel-Nachzüge (64 ADRs, 39 eingefrorene Dateien) —
   `git diff --cached --name-status -M` nennt eine `R100`-Zeile und 103 `M`-Zeilen (**kein** `D`/`A`-Paar,
   sonst ist die Rename-Erkennung weg); `make adr-check STAGED=1` Exit 0. **Commit 2:** alles andere der Tabelle.
   Beide Commits in **denselben Push**: Der erste ist allein nicht baubar. **Lesart, keine Ausnahme** ([Verdikt](../../../reviews/2026-11-26-welle-98-architect-verdikt.md), Frage 4):
   [`AGENTS.md`](../../../../AGENTS.md) §3.3 schützt die Rename-Erkennung der **bewegten** Datei, sie bleibt `R100`; die Link-Ziele
   **anderer** Dateien reisen mit, und die Spezifikation selbst (Selbst-Links) folgt im zweiten Commit. Ein eigener MR entfällt.
4. Review (höchstens zwei Runden), Verifikation, Closure. `git commit -F <datei> -- <Pfade>`, vorher
   `git diff --cached --stat`.

**Gates:** `make doc-check`, `make adr-check STAGED=1` (je Commit) und `make adr-check RANGE=<Vor-Bump>..HEAD`,
`make test`, `make image-test`, `make gates`; **nicht** hier: `make blackbox-probe` und der Vergleich von `--manual` — sie sind
Closure-Trigger der Welle (§3 (d)).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-295`](../done/slice-295-baseline-v700-bump.md) liegt in `done/`;
`in-progress/` trägt keinen Slice (nur die Roadmap); das **Architect-Verdikt** zu den zwölf Fragen in
[welle-98 §5](../welle-98-baseline-v700.md) (für diesen Slice tragend: 4, 5, 6, 7) liegt vor: [Verdikt](../../../reviews/2026-11-26-welle-98-architect-verdikt.md);
Priorisierung gesetzt, `Verantwortlich:` beim Übergang `open` → `next` gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Wiederholung der Zählung weicht von §3 ab — mehr Dateien, die
  den Namen **nicht** als Link-Ziel tragen und nicht unter das Ventil fallen — oder der Diff wächst über eine
  Review-Sitzung; dann gilt die Form von [`MR-066`](../../../../harness/conventions.md#mr-066).
- `in-progress` → `open` (blockiert): `make adr-check STAGED=1` meldet für den Move-Commit `core-drift-vcs` (die
  Messung dieses Plans trägt nicht mehr, etwa weil `vcs.follow-renames` sich geändert hat), oder das Architect-Verdikt zu
  Frage 4 oder 5 steht aus.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (§2) und `make gates` grün; Review- und Verifikations-Report liegen unter `docs/reviews/`;
Closure-Notiz mit Lerneintrag geschrieben. **Review-Vorgabe der Welle:** höchstens zwei Runden, eingearbeitet
werden nur HIGH und MEDIUM; jede Nachbesserung zitiert die Stelle, die sie trägt
([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 18). Was die **Welle** einlöst — das Produkt ändert sich an genau den
benannten Stellen, die Abweichungsklassen von `--manual`, die Gegenmessung am Endstand — steht in
[ihrem Closure-Trigger](../welle-98-baseline-v700.md) und nicht in dieser DoD
([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 10).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die mechanische Ersetzung trifft eine Datei, die sie nicht treffen soll** — ein Dateiname, der auf `spezifikation.md`
  endet, ein Konsumenten-Beispiel, eine Vergangenheits-Aussage
  ([`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md),
  5×). Im Probelauf ist es **eingetreten** (ein Link in `welle-80-results.md` zeigte danach auf einen Namen, den es nicht
  gibt). — **Ausgang:** *eingetreten.* Die zwei Messungen des Plans halten — `make doc-check` meldet 0 Befunde
  (Verifikation M2: 1325 Dateien), `git diff -U0 f44c6e65..HEAD -- docs/user/benutzerhandbuch.md` nennt genau einen Hunk,
  `@@ -2792 +2792 @@`, die zwei Dateien, die auf `spezifikation.md` enden, stehen nicht im Diff, und die Substitutionsprobe
  fand 103 von 103 Dateien des Move-Commits identisch zur reinen Ersetzung des Link-Ziels —, aber sie haben die Klasse
  „Vergangenheits-Aussage“ nicht gefangen: Der Review fand in [`MR-006`](../../../../harness/conventions.md#mr-006) eine Zeile,
  die den Namen als Aussage über das Template-Set von 2026-06 führte und von `340f3f15` mitersetzt wurde (R1 F-2, LOW;
  behoben in `5e607ad3`). Der Ausgang „entfallen“ nach dem Wortlaut des Plans wäre formal wahr und in der Sache falsch.
  Evidenz-Datei `evidence/slice-297.md` unter dem genannten Eintrag, jetzt 6×.
- **Der Move-Commit allein ist nicht baubar** und wird die Spitze eines Push. — **Ausgang:** *weiter offen.* Nichts ist
  gepusht (`git log origin/main..HEAD` zeigt die Commits dieses Slice an der Spitze von 34 ungepushten), die Bedingung des
  Plans, ein Push und die CI am Spitzen-Commit, ist zur Closure nicht erfüllbar; „eingetreten“ wäre ebenso falsch, denn
  der Zwischenstand ist nicht die Spitze eines Push geworden. Die Spitze `5e607ad3` ist lokal grün (Verifikation M1,
  M7); die Lesung der CI nach dem Push trägt die Push-Pflicht in §7.
- **Ein Test läuft grün, ohne die Regel zu fahren**
  ([`BEO-ALL/testfall-faehrt-seine-regel-nicht`](../observations/BEO-ALL/testfall-faehrt-seine-regel-nicht/observation.md),
  2×). — **Ausgang:** *entfallen.* Die Verifikation fuhr je Test die Mutation im Wegwerf-Klon und las die Meldung; beide
  waren rot, aus dem behaupteten Grund (§7).
- **Das repo-weite Ventil verdeckt künftig einen lebenden Verweis auf den alten Pfad** (benannte Grenze von `ignore-refs`).
  — **Ausgang:** *entfallen.* Der Kommentar am Eintrag in der [`.d-check.yml`](../../../../.d-check.yml) trägt die Grenze mit ihrer
  Bedingung (eine `GRENZE:`-Zeile: der Eintrag wirkt repo-weit, eine lebende Inline-Code-Nennung des alten Pfads bleibt
  still, ein Link darauf wird von `links` gemeldet), in der Form des Vorgängers für `docs/user/releasing.md`
  ([`ADR-0025`](../../adr/0025-codepaths-ignore-refs.md)); der Reviewer fuhr die Probe (lebende Nennung still, lebender
  Link `target-missing`, ein anderer toter Pfad `codepath-missing`).
- **Eine Zahl dieses Plans stimmt nicht** — Zählung am Proxy statt am Gegenstand
  ([`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md),
  10×). — **Ausgang:** *eingetreten.* Reviewer (F-3) und Verifier (V-2) zählten das Ventil nach: 120 Befunde statt 119,
  elf lebende Nennungen statt zehn, zehn Nennungen in den Plänen der Welle statt neun (die Paare mit Kommando in §7); alle
  übrigen Zahlen aus §3 reproduzierten sich. Behoben im Wellen-Plan (`0315fc57`) und in der Closure-Notiz; die Zahlen in
  §3 bleiben als Plan-Stand stehen, die Plan-Änderungen überschreiben die Tabelle ausdrücklich. Evidenz-Datei
  `evidence/slice-297.md` unter dem genannten Eintrag, jetzt 12×.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** Der Zuschnitt „reiner Rename, Link-Ziele, kein Wortlaut“ hielt. Der Move-Commit `56c8f62f` ist
  **eine** `R100`-Zeile und 103 `M`, ohne `D`/`A`-Paar; je Datei gegengerechnet (der alte Stand mit ersetztem Link-Ziel hinter
  `](` gegen den neuen) sind 103 von 103 identisch, und `make adr-check RANGE=f44c6e65..HEAD` endet über die Spanne und
  beide Teilspannen mit 0 Befunden. Die Zerlegung nach [`AGENTS.md`](../../../../AGENTS.md) §3.3 trug, wie der Plan sie vorab
  sagte: einzeln rot sind nur `56c8f62f` (nicht baubar, `doc-check` Exit 2 „konfiguriertes ids-Target existiert nicht“) und
  `340f3f15` (212 Befunde, 205 `target-missing` und 7 `matrix-forbidden`, alle in der Spezifikation), die sieben übrigen
  Commits der Spanne melden 0. Eine Review-Runde (1 MEDIUM, 2 LOW, 3 INFO) und eine Fixrunde.
- **Die Mutation je Test (Risiko 3), im Wegwerf-Klon bei `5e607ad3`, die Meldung gelesen:** Mit dem alten Namen in
  `cli_manual_test.go` Z. 21 endet `make test` mit Exit 2, `--- FAIL: TestManual_GibtAbschnitteAus`, Meldung „Ausgabe ohne
  ==> spec/spezifikation.md:“ — die Kopfzeile des Produkts trägt den neuen Namen, der Test erwartete den alten, also aus dem
  behaupteten Grund. Mit dem alten Namen in `diagnose_test.go` Z. 73: Exit 2,
  `--- FAIL: TestAllReasonsDeckungGegenSpezifikationGrundCodes`, „Spezifikation nicht lesbar (fail-closed): open
  ../../../../spec/spezifikation.md: no such file or directory“ (Meldung wie im Plan). Mit dem alten Namen in
  `tools/image-test.sh` Z. 163: `make image-test` Exit 2, „image-test: FAIL — --manual: kein Abschnitt aus der
  Spezifikation“. Mit dem alten Pfad im `//go:embed` von `manual.go`: der Build bricht, „pattern spec/spezifikation.md: no
  matching files found“ — der Grund, warum der Move allein nicht baubar ist.
- **Was ging anders als geplant:** **(1) Sechs Dateien trugen ein Link-Ziel, das keine Zeile der Tabelle führte** (vier
  Register-Belege, zwei Pläne); Plan-Änderung vor dem Code (`e0a48bb9`), dazu die feinere Zerlegung in vier Commits.
  **(2) Das Konsumenten-Beispiel der Spezifikation (zwei Zeilen im Konfigurations-Beispiel) bleibt**, damit trägt sie 159
  geänderte Zeilen statt 161; das Ventil deckt 119 Befunde statt 104; die Prosa-Zeile 920 (der Pfad zu `--manual`) ist eine
  fünfte Klasse der Ausgabe von `--manual` — zweite Plan-Änderung vor dem Code (`8a458a27`). **(3) Review F-1 (MEDIUM): die
  „fünfte Klasse“ hatte keinen Träger.** Der Slice-Plan sagte, die Welle führe sie; der Wellen-Plan führte weiter vier Klassen
  und den Satz „Eine fünfte Klasse ist ein Befund“, die Gegenmessung am Endstand hätte die Zeile 920 als Abweichung
  gefunden. Behoben in `0315fc57` (die Welle führt fünf Klassen, „Eine sechste Klasse ist ein Befund“, die Endzahlen des
  Ventils). **(4) Review F-2 (LOW): eine Vergangenheits-Aussage in einer lebenden Datei wurde mitersetzt.**
  [`MR-006`](../../../../harness/conventions.md#mr-006) (Datum 2026-06-10) sagt, welchen Namen das damals adoptierte
  Template-Set trug; `340f3f15` setzte den neuen. Behoben in `5e607ad3` (gegen `f44c6e65` netto unverändert, das Diff der Datei ist
  leer). **(5) Review F-3 (LOW) und Verifikation V-2: Zahlen des Plans** — Lernsignal 1.
- **Vorbehalt zum DoD-Punkt „Review durchgeführt“:** Die Fixrunde (`0315fc57`, `5e607ad3`) bekam keine zweite Review-Runde. Anders
  als bei [`slice-295`](../done/slice-295-baseline-v700-bump.md) änderte sie **keine Gate-Konfiguration**: ein Satz-Nachzug im
  Wellen-Plan und eine Zeile in [`MR-006`](../../../../harness/conventions.md#mr-006). Der Verifier las beide und bestätigte F-1 und F-2 als behoben; der Haken steht mit
  diesem Vorbehalt.
- **Steering-Loop-Eintrag:** mit diesem Slice ist keine Regel, kein Sensor und keine benannte Spec-Lücke verkörpert; zwei
  Beobachtungen sind belegt (Register unten). **Lernsignal 1, mit Ursache:** Die Zahl des Ventils stimmte am Endstand
  nicht (119 statt 120), **weil** die Plan-Änderung, die sie nennt, den alten Pfad selbst in Inline-Code schreibt und damit
  zu der Menge gehört, die das Ventil deckt: Gezählt wurde vor dem Schreiben des Absatzes, gelesen wurde nach ihm. Gefunden haben
  es Reviewer (F-3) und Verifier (V-2), die das Kommando fuhren, nicht der Zählende. **Die Paare, getrennt, am Stand
  `5e607ad3`** (Form: die `codepath-missing`-Befunde von `d-check` im Wegwerf-Klon **ohne** den Eintrag in
  `codepaths.ignore-refs`, `docker run --rm --network none -v <Klon>:/repo:ro d-check:latest`, nach Pfad gruppiert;
  Verifier und Planner fuhren es unabhängig, beide 120): **Befunde insgesamt 119 → 120; Nennungen in den Plänen der Welle
  9 → 10** (`welle-98` 4, dieser Plan 4, `slice-296` 2); **lebende Nennungen insgesamt 10 → 11** (die zehn Nennungen der
  Pläne und [`MR-081`](../../../../harness/conventions.md#mr-081) mit einer); **eingefroren unverändert 109** (`done/` 51, ADRs 44, `conventions/done` 6, Register-Belege
  5, `CHANGELOG.md` 2, `docs/plan/cr` 1). Die Summe 11 + 109 ist 120. Die Einteilung gilt für diesen Stand. **Lernsignal 2,
  mit Ursache:** Review F-2 traf dieselbe Klasse wie [`MR-070`](../../../../harness/conventions.md#mr-070), eine Ebene
  tiefer: nicht eine Datei der Frozen-Liste, sondern **eine Zeile** in einer lebenden Datei, die einen Namen als Aussage über
  einen Fremdbestand führt. Der Plan zählte [`MR-006`](../../../../harness/conventions.md#mr-006) unter den 49 lebenden Zeilen, ohne sie gelesen zu haben, **weil** die
  Einteilung nach Datei-Eigenschaft eine Zeile in einer lebenden Datei nicht sehen kann — die Grenze, die der Eintrag und
  die Beobachtung schon benennen, hier zum sechsten Mal.
- **Beobachtungs-Register (`../observations/`):** je eine `evidence/slice-297.md` in
  [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (12×, `verkörpert`; Risiko 5: F-3 und V-2, die Zahl zählt ihre eigene Nennung) und in
  [`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
  (6×, `verkörpert` in [`MR-070`](../../../../harness/conventions.md#mr-070); Risiko 1: F-2). Benannt, nicht gezählt: F-1 (ein
  Folgeakt am Plan einer anderen Stelle, im Slice-Plan benannt, ohne dass dort ein Träger stand; behoben, ein erster Fall
  ohne Beobachtung im Register); F-4 und V-3 (das Ventil deckt elf lebende Nennungen mit, wo [`ADR-0025`](../../adr/0025-codepaths-ignore-refs.md) für lebende
  Dateien den Zeilen-Marker führt; gemessen ohne Wirkung, der Eintrag ist wegen der 109 eingefrorenen ohnehin repo-weit);
  F-5 (die Historie der Spezifikation führt keine Zeile zum Namenswechsel, das Lastenheft schon; keine Regel verlangt sie);
  F-6 (das Konsumenten-Beispiel nennt weiter den alten Namen, und `--suggest-config` schlägt ihn vor — Träger ist die
  Sichtung, [`slice-296`](../open/slice-296-baseline-v700-neuerungen-sichtung.md) Frage B); V-1 (die zwei einzeln roten
  Zwischenstände, in der Push-Pflicht). Keine davon ist umgesetzt.
- **Architektur-Beobachtung:** Ein Rename trägt drei Gegenstände, und drei Wächter urteilen je über einen: die **Link-Ziele**
  (`links`, `anchors`, dazu der ADR-Gate: ein Commit meldet 0 `core-drift-vcs`, zwei Commits meldeten 60, gemessen im Plan —
  deshalb reisen die Link-Ziele der ADRs im Move-Commit), die **Inline-Code-Nennungen** (`codepaths`, das repo-weite Ventil
  mit 120 Befunden dahinter) und die **Kommentar- und Prosa-Nennungen**, die **kein** Wächter hält: 95 Go-Kommentarzeilen in
  38 Dateien, je eine in `Dockerfile` und `Makefile`, zwei in `tools/` — das benannte Negativ des Architect-Verdikts zu
  Frage 7, ohne Wirkung auf Verhalten. Dazu ein Spiegel, den keiner dieser Wächter hält: ein Beispiel, das an vier Stellen
  gleich steht ([`MR-025`](../../../../harness/conventions.md#mr-025)), ändert sich an allen oder an keiner; hier an keiner,
  und nur die Sichtung entscheidet es.
- **Push-Pflicht:** Die Commits `e0a48bb9` bis `5e607ad3` (neun) gehören in **einen** Push, **zusammen** mit den Commits davor
  und danach, die noch vor `origin/main` liegen: `56c8f62f` (nicht baubar) und `340f3f15` (212 Befunde) sind einzeln rot — der
  Zwischenstand, den der Plan vorab benennt, solange er nicht die Spitze eines Push wird. Nach dem Push die CI an der
  **Spitze** lesen, nicht aus dem lokalen Stand schließen; Risiko 2 bleibt bis dahin offen.
- **Release-Prep-Notiz für das `CHANGELOG`** (nicht im Feature-Commit,
  [`AGENTS.md`](../../../../AGENTS.md) §5 Regel 17): Die Ausgabe von `d-check --manual` ändert sich an vier Stellen — die
  Kopfzeile `==> spec/specification.md:<Zeile>` statt des deutschen Namens, die Link-Ziele der Spezifikation in der Ausgabe
  (158), die Handbuch-Zeile 2792 und die Prosa-Zeile 920 der Spezifikation (der Pfad zu `--manual`). Der Vertrag
  [`DC-FA-CLI-013`](../../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen) bleibt im
  Wortlaut; ein Konsument, der die Kopfzeile parst, sieht den neuen Namen.
- **Welle-Closure, nicht dieser Slice:** `make blackbox-probe REF=738d2bca`, der Vergleich von `--manual` je Begriff mit den **fünf**
  Klassen und `make fullbuild` sind hier **nicht** gefahren; die Verifikation zeigte nur, dass Kopfzeile und Zeile 920 den neuen
  Pfad tragen (`docker run … d-check:latest --manual manual`). Sie gehören dem Closure-Trigger von
  [`welle-98`](../welle-98-baseline-v700.md).
- **Folge-Slices:** keine neuen. [`slice-296`](../open/slice-296-baseline-v700-neuerungen-sichtung.md) (Frage B, das
  Konsumenten-Beispiel) und [`slice-298`](../open/slice-298-neue-namen-wave-archive-anker.md) (der Umzug von `welle-98`, der die
  Verweise dieses Slice mitzieht) bestehen schon in `open/`. Für die Kommentarzeilen mit altem Namen ist **kein** Slice
  geschnitten: das Verdikt (Frage 7) hat sie entschieden, der Nachzug ist die größere Antwort und ein eigener Akt.
- **Risiken aus §6:** Risiko 1 *eingetreten* (Beleg: die Evidenz-Datei zu F-2), 2 *weiter offen* (bis Push und CI), 3 *entfallen*,
  4 *entfallen*, 5 *eingetreten* (Beleg: die Evidenz-Datei zu F-3 und V-2) — Begründung je Risiko in §6.
- **Die drei Paarungen**, hier vorab geprüft (getragen von der Closure von welle-98): (a) Anker — kein Steering-Loop-Eintrag
  mit `liegt in`-Feld, keine Paarung auszulösen; (b) Folge-Slice — `slice-296` und `slice-298` bestehen als Dateien in `open/`;
  (c) Register — beide genannten Verzeichnisse existieren und tragen nicht leere `evidence/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:416-417 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Die Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md) führt zwei Zeilen: den Default `*` (`ALL`) und
`tools/harness/` (`HARN`), beide Greenfield. Berührt ist der Default `*`: `spec/`, `docs/`, `harness/`, die
Konfiguration im Repo-Wurzelverzeichnis und der Produkt-Quelltext (`manual.go`, zwei Tests) tragen keine eigene
Deklaration (gemessen an der Tabelle §Modus-Deklaration pro Sub-Area); `tools/image-test.sh` liegt unter `tools/`, nicht
unter `tools/harness/`. `tools/harness/` bleibt **unberührt**.

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:422-422 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stand 2026-11-25,
`ls docs/plan/planning/observations/BEO-ALL/`; Zähler = Dateien unter `evidence/`), Sub-Area `*`. Treffer:
[`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
(verkörpert in [`MR-070`](../../../../harness/conventions.md#mr-070), 5× — Risiko 1: die Frozen-Klassen werden über die
Eigenschaft gelistet, **vor** dem ersten Tausch; im Probelauf trat die Wortgrenzen-Falle auf),
[`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/observation.md)
(verkörpert, 4× — der Spiegel der Spezifikation liegt in Handbuch, READMEs, `harness/sensors/` und Konfiguration; §3 zählt
sie),
[`BEO-ALL/testfall-faehrt-seine-regel-nicht`](../observations/BEO-ALL/testfall-faehrt-seine-regel-nicht/observation.md)
(offen, 2× — Risiko 3: die zwei Tests tragen ihre Mutation),
[`BEO-ALL/path-scoped-commit-carries-staged-rest`](../observations/BEO-ALL/path-scoped-commit-carries-staged-rest/observation.md)
(verkörpert, 5× — `git diff --cached --name-status -M` vor Commit 1, §3),
[`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
(verkörpert, 10× — Risiko 5),
[`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md)
(verkörpert in Regel 13, 3× — die Grenze des repo-weiten Ventils steht mit ihrer Bedingung im Kommentar),
[`BEO-ALL/commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/observation.md)
(verkörpert, 26× — die Botschaft behauptet nur, was der Diff trägt). **Kein Treffer:** keine Beobachtung zu einem
umbenannten Spec-Stratum. Der Zähler folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-25: `image-scan` **grün** (jüngster Lauf 2026-11-25T12:59:28Z). `upstream-drift`
**rot** (jüngster Lauf 2026-11-25T09:26:13Z), der Log war über `gh run view` nicht abrufbar; gelesen ist
`make baseline-freshness` (Exit 3, neuer Release `v7.0.0`) — der Gegenstand der Welle, **planmäßig gelesen**. Zum
Zeitpunkt dieses Slice ist [`slice-295`](../done/slice-295-baseline-v700-bump.md) `done/`; der Stand wird beim Start neu gelesen.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
