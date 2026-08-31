# Verifikation slice-185 — `doc-usage` im `--print-mk`-Fragment

**Art:** Verifikation.
**Gegenstand:** slice-185, Range `78ff14e..HEAD` — `4db8e2b` (eigenständige Vorab-Korrektur, im Bereich, nicht Teil der DoD), `ecf5ec3` (Plan angelegt), `39258a4` (Claim-Move `open/` → `in-progress/`), `8e75841` (Implementierung), `49dbf9d` (Review-Report — **nicht gelesen**, siehe unten), `920f97c` (Korrektur M-1 nach Review).
**Modell-ID:** `claude-sonnet-5`.
**Datum:** 2026-08-31.
**Rolle:** unabhängiger Verifier — geprüft wird gegen **DoD und Spec** (Baseline-Regelwerk [`modul-08-agentenrollen.md` §Welche Rolle braucht welche Artefaktklasse](../../.harness/baseline/v5.12.0/regelwerk/modul-08-agentenrollen.md#welche-rolle-braucht-welche-artefaktklasse-modul-8)), nicht gegen Plan/ADR-Maintainability. Der Review-Report `docs/reviews/2026-08-31-slice-185-code-r1.md` wurde bewusst **nicht** gelesen, um einen von der Review-Rolle unabhängigen Blick zu halten.
**Prüfgrundlage:** Slice-Plan `docs/plan/planning/done/slice-185-print-mk-doc-usage.md` §4 (DoD), §2/§3; [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) (0.80.0) und [`DC-FA-CLI-001`](../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel) im Lastenheft; §[`DC-FA-CLI-010.a`](../../spec/spezifikation.md#dc-fa-cli-010a--makefile-fragment) in der Spezifikation; das tatsächlich erzeugte Fragment aus einem lokal aus HEAD gebauten Image (`d-check:latest`, Image-ID `sha256:5c63708288bd…`, per `make build` neu gebaut — identisch zum bereits vorhandenen Stand, was belegt, dass sich der Code seit `8e75841` nicht mehr verändert hat).
**Nicht verändert:** kein Repo-Artefakt außer dieser Datei. Für die Umkehr-Probe wurde `internal/adapter/driving/cli/print_mk.go` temporär mutiert und per `git checkout` exakt zurückgestellt (Beleg unten); `git status --short` war danach und ist bei Abgabe leer bis auf diese Report-Datei.

---

## 1. DoD-Tabelle (§4 des Slice-Plans)

| # | Behauptet | Gemessen | Verdikt |
|---|---|---|---|
| 1 | `doc-usage` steht im Fragment — `@`-echo-unterdrückt, `##`-annotiert, mit `--help` gegen `$(DCHECK_REF)`; `make doc-help` listet es auf | Fragment aus `d-check:latest --print-mk` (netzlos) enthält `.PHONY: doc-usage` / `doc-usage: ## Aufruf und Optionen von d-check selbst (--help)` / `\t@docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --help`. In einem Wegwerf-Makefile (`include d-check.mk`, `DCHECK_IMAGE := d-check:latest`) gefahren: `make doc-usage` → Exit 0, **kein** `docker run …`-Recipe-Echo im Ausgabestrom (0 Treffer für `docker run` in der abgefangenen Ausgabe). `make doc-usage > out.txt` → `out.txt` **0 Byte**, die Hilfe erscheint trotzdem (auf stderr). `make doc-help` listet `doc-usage  Aufruf und Optionen von d-check selbst (--help)` alphabetisch zwischen `doc-tracked` und `doc-help`-Nachbarn | **erfüllt** |
| 2 | Der Vertrag führt dreizehn Targets: Lastenheft (Beschreibung, beide Akzeptanzkriterien, Out-of-Scope, Bump + Historie) und Spezifikation (Punkt 5: Zahl **und** Aufzählung, Historie) | Lastenheft: Version `0.79.0` → **`0.80.0`**; Beschreibung nennt `doc-usage` mit stderr-Eigenschaft; **Happy-Path**- und **Boundary**-Akzeptanzkriterium führen `doc-usage`; Out-of-Scope zählt „dreizehn" und nennt `doc-usage` in der Klammer-Liste; Historie-Zeile `0.80.0` vorhanden. Spezifikation: Punkt 5 „Dreizehn `.PHONY`-Targets"; eigener Aufzählungspunkt für `doc-usage` (Modus statt Regelmodul, `@`-Suppression, stderr, Mount-Begründung); zwei Historie-Zeilen `2026-08-31` (Zusatz-Target **und** die separat nachgezogene `doc-structure`-Lücke). `git diff 78ff14e..HEAD -- spec/lastenheft.md spec/spezifikation.md` bestätigt alle Fundstellen wörtlich | **erfüllt** |
| 3 | Handbuch §4.16 nennt dreizehn Targets, `doc-usage` in der Liste, und den Ausgabe-Strom — gemessen, nicht gelesen. §11-Zeile und Handbuch-Kopf bleiben der Release-Prep vorbehalten | §4.16 (Zeilen 824–888): „dreizehn `##`-annotierte Targets" mit `doc-usage` in der Klammer-Liste; Fließtext „`make doc-usage` zeigt Aufruf und Optionen … **Die Hilfe erscheint dabei auf `stderr`, nicht auf stdout** … ein `make doc-usage > datei` fängt nichts ein und meldet trotzdem Erfolg" — das ist exakt das, was Zeile 1 dieser Tabelle unabhängig nachgemessen hat. Handbuch-Kopf (`Handbuch-Version 1.64`, `Software-Version v0.70.0`, `Stand 2026-08-30`) und §11 (`## 11. Änderungshistorie`, Zeile 2370) liegen **außerhalb** jedes Diff-Hunks von `78ff14e..HEAD` | **erfüllt** |
| 4a | Umkehr-Probe (BEO-023): ohne den Template-Block wird genau der neue Test rot | **Selbst reproduziert**, nicht nur gelesen: `doc-usage`-Block aus `print_mk.go` entfernt, `docker build --target test .` (volles `go test ./...` in Docker) → **genau ein** `--- FAIL: TestCLI185_PrintMK_UsageTarget`, alle übrigen Pakete `ok`. Datei danach exakt auf HEAD zurückgesetzt (`git checkout --`), `git status --short` leer, `git diff 78ff14e..HEAD -- internal/adapter/driving/cli/print_mk.go` unverändert zum Ausgangsbefund | **erfüllt** |
| 4b | … mit Ausgabe in der Closure-Notiz | §9 des Slice-Plans ist eine leere Überschrift (der Slice liegt in `in-progress/`) | **noch nicht fällig** (Closure-Schritt) |
| 5 | Die Spiegel-Liste aus §2.1 steht in der Closure-Notiz — vollständig abgearbeitet, mit der Zahl der angefassten Stellen | §9 leer | **noch nicht fällig** |
| 6 | `make gates` grün (Exit explizit); unabhängiger Review; Verifikation — beide in eigenen Kontexten | `make gates` selbst gefahren: **Exit 0**, Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; semgrep „Ran 55 rules on 55 files: 0 findings.". Review-Commit `49dbf9d` existiert in eigenem Kontext (Inhalt bewusst ungelesen); diese Datei ist die Verifikation, in ihrem eigenen, vom Review getrennten Kontext | **erfüllt** |
| 7 | Closure-Notiz mit Lerneintrag; Register fortgeschrieben; jedes Risiko aus §5 mit Ausgang; die drei Paarungen | §9 leer, §5-Risiken tragen alle vier `*(bei Closure)*`, `docs/plan/planning/observations.md` unverändert im Diff (`git diff 78ff14e..HEAD -- docs/plan/planning/observations.md` leer) | **noch nicht fällig** |

**Zusammenfassung:** 5 Punkte **erfüllt** (1, 2, 3, 4a, 6), 3 Punkte **noch nicht fällig** (4b, 5, 7 — Closure-Schritte, korrekt: der Slice liegt in `in-progress/`), **0 nicht erfüllt**.

## 2. Akzeptanzkriterien `DC-FA-CLI-010` gegen das laufende Produkt

Alle Läufe: lokal aus HEAD gebautes Image `d-check:latest`, `--network none`, kein `.git`-Zugriff nötig.

| Akzeptanzkriterium (Lastenheft 0.80.0) | Zusage | Gemessen | Verdikt |
|---|---|---|---|
| **Happy Path** | Fragment mit `DCHECK_IMAGE`, `DCHECK_DIGEST`, `TRACE_FLAGS` und den dreizehn Targets inkl. `doc-usage` (`--help`, Recipe-Echo unterdrückt, Hilfe auf **stderr**, Exit 0), Handbuch-URL im Kopf | `d-check:latest --print-mk` → **Exit 0**, 75 Zeilen; alle dreizehn `.PHONY:`-Blöcke vorhanden (gezählt); `doc-usage`-Recipe exakt `\t@docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --help`; Kopf trägt `https://github.com/pt9912/d-check/blob/main/docs/user/benutzerhandbuch.md` | **erfüllt** |
| **Boundary** — `doc-usage` → `--help` mit unterdrücktem Recipe-Echo, Ausgabe auf stderr | wie beim direkten Aufruf | `make doc-usage` (Fragment `include`-t): Exit 0, kein `docker run`-Recipe-Echo. Direkter Produktlauf `d-check:latest --help`: **0 Byte stdout, 2488 Byte stderr, Exit 0** — Byte für Byte identisch zur Zahl aus der Lastenheft-Historie-Zeile 0.80.0. `make doc-usage > out.txt`: `out.txt` 0 Byte, Hilfetext erscheint trotzdem auf dem Terminal (stderr) | **erfüllt** |
| **Negative** — unbekanntes Flag bei `--print-mk` ⇒ Exit 2 | Nutzungsfehler | `d-check:latest --print-mk --nope` → `d-check: error: flag provided but not defined: -nope`, **Exit 2** | **erfüllt** (unverändert von diesem Slice, gegengeprüft) |
| Spezifikation Punkt 5 — „sechs der dreizehn tragen eine Fokus-Liste, sieben nicht" | Zahlenaussage nach M-1-Korrektur | Am erzeugten Fragment nachgezählt: `--disable` tragen `doc-immutable`, `doc-commits`, `doc-planning`, `doc-tracked`, `doc-targets`, `doc-structure` (**6**); ohne `--enable`/`--disable`: `doc-check`, `doc-trace`, `doc-complete`, `doc-doctor`, `doc-repair`, `doc-usage`, `doc-help` (**7**) — 6+7 = 13 | **erfüllt** |

**Kein Akzeptanzkriterium ist unerfüllt.**

## 3. Plan-vs-Code-Diff

**Richtung 1 — im Diff, aber vom Plan nicht angekündigt:**

- `docs/user/operations.md` (Commit `8e75841`): die `--print-mk`-Zeile der CLI-Optionstabelle nannte noch drei Targets aus der 0.24.0-Ära (`doc-check`, `doc-trace`, `doc-complete`) und keine `DCHECK_DIGEST`-Variable. Weder §2 (die vier gelisteten Spiegel-Orte: Lastenheft, Spezifikation, `print_mk.go`, Handbuch) noch §3 noch §4 des Slice-Plans nennen `operations.md`. Die Commit-Botschaft von `8e75841` legt den Fund offen (*„Und ein Spiegel, den die Vorab-Liste des Plans nicht erfasst hatte: operations.md nannte drei Targets von zwölf…"*) und ersetzt die veraltete Aufzählung durch einen alterungsfreien Zeiger auf das erzeugte Fragment. Kein DoD-Punkt verlangt diese Datei, keine Ausschluss-Klausel aus §3 verbietet sie — die Erweiterung liegt im Sinn des Slice-Ziels („die Aufzählung ist der Gegenstand"), ist aber eine **unangekündigte** Erweiterung des Diff-Umfangs gegenüber dem geschriebenen Plan.
- `docs/plan/planning/in-progress/roadmap.md` (Commit `39258a4`, −2 Zeilen „Nichts in Arbeit."): das ist die **deklarierte** Ausnahme aus `AGENTS.md` §3.3 / [`MR-013`](../../harness/conventions.md#mr-013), die den Lifecycle-Move-Commit begleitet — kein unangekündigter Fund, sondern Konvention.

**Richtung 2 — vom Plan zugesagt, aber im Diff nicht wiederzufinden:** keine Lücke gefunden. Alle drei Liefer-Punkte aus §2 (Fragment · Vertrag · Handbuch) sind im Diff; die Umkehr-Probe ist nicht als Commit sichtbar (erwartungsgemäß — sie ist eine Mess-Prozedur, kein Artefakt), aber von mir unabhängig reproduziert (§1 Zeile 4a).

**§3 „Ausdrücklich NICHT" — alle fünf Grenzen gehalten:**

1. `doc-help` unverändert (Recipe, Annotation, Position identisch vor/nach dem Diff, nur `doc-usage` davor eingefügt).
2. Keine Änderung der `--help`-Ausgabe selbst — nur `internal/adapter/driving/cli/print_mk.go` und sein Test geändert, kein Flag-Definitions-Code berührt.
3. Kein `2>&1` im neuen Recipe — bestätigt, `doc-usage`-Zeile trägt keine Strom-Umlenkung.
4. `-v "$(CURDIR):/repo:ro"`-Mount im `doc-usage`-Recipe vorhanden, obwohl `--help` ihn nicht braucht.
5. Kein ADR-Datei-Diff, kein `CHANGELOG.md`-Diff im gesamten Range (`git diff 78ff14e..HEAD --stat -- CHANGELOG.md README.md README.de.md docs/plan/adr/` liefert nichts).

## 4. Spec-Konformität (§3.4 Abwärts-Sperre)

`git diff 78ff14e..HEAD -- spec/lastenheft.md spec/spezifikation.md` — die hinzugefügten Zeilen enthalten **kein** `ADR-\d{4}`-, `slice-\d{3}`- oder `welle-\d{2}`-Token (per `grep` gegengeprüft). Beide Historie-Zeilen benennen weder ADR noch Slice; §3.4 ist gehalten.

## 5. Gate-Läufe (echte Ausgabe)

| Lauf | Ergebnis |
|---|---|
| `docker build --target runtime … -t d-check:latest .` (`make build`) | Neubau aus HEAD, Image-ID **identisch** zum vorher schon vorhandenen `d-check:latest` (`sha256:5c63708288bd…`) — Beleg, dass der Code seit `8e75841` unverändert ist |
| `d-check:latest --print-mk > d-check.mk` | **Exit 0**, 75 Zeilen, dreizehn `.PHONY`-Blöcke |
| `d-check:latest --help` | **Exit 0**, `0` Byte stdout, `2488` Byte stderr |
| `d-check:latest --print-mk --nope` | **Exit 2**, `flag provided but not defined: -nope` |
| `make doc-usage` (Wegwerf-Makefile, `DCHECK_IMAGE := d-check:latest`) | **Exit 0**, kein Recipe-Echo, Hilfe auf stderr |
| `make doc-usage > out.txt` | **Exit 0**, `out.txt` 0 Byte |
| `make doc-help` | **Exit 0**, `doc-usage` alphabetisch gelistet |
| `make doc-check` (Wegwerf-Makefile, Sanity) | **Exit 0**, `0 Datei(en) geprüft, 0 Befund(e)` — Nachbar-Target unbeeinflusst |
| Umkehr-Probe: `docker build --target test .` mit entferntem `doc-usage`-Block | **Exit 1** (Build scheitert am `go test`), genau `--- FAIL: TestCLI185_PrintMK_UsageTarget`, sonst alles `ok` |
| `git checkout -- internal/adapter/driving/cli/print_mk.go` | Datei exakt auf HEAD zurückgestellt, `git status --short` leer |
| `make gates` | **Exit 0** — `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; Coverage 94,70 % ≥ 93 %; semgrep 55 Regeln, 0 Befunde |

## 6. Verdikt

**Bestanden.** Alle fünf fälligen DoD-Punkte (1, 2, 3, 4a, 6) sind gegen das laufende Binary erfüllt, nicht nur gegen die Tests — inklusive einer selbst reproduzierten Umkehr-Probe, die exakt den behaupteten einen Test rot schaltet. Die drei übrigen DoD-Punkte (4b, 5, 7) sind Closure-Schritte und planmäßig noch offen, da der Slice in `in-progress/` liegt. Beide Akzeptanzkriterien von `DC-FA-CLI-010` (0.80.0) halten am Produkt, ebenso die Negative-AK (unverändert). Die einzige Abweichung zwischen Plan und Diff ist eine **positive**, im Commit selbst offengelegte Erweiterung (`operations.md`), die keinem DoD-Punkt widerspricht und keine der fünf Ausschluss-Klaumerzen aus §3 verletzt. `make gates` ist mit echtem Exit 0 gegengeprüft.
