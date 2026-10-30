# slice-281: `--print-mk` — `doc-manual` (Handbuch und Spezifikation des Werkzeugs lesen)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Auftraggeber-Wunsch, die Closure-Bedingung ist die DoD
selbst.

**Bezug:** [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
(das `--print-mk`-Fragment),
[`DC-FA-CLI-013`](../../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen)
(`--manual <begriff>`). Anlass: Der Auftraggeber fragt, ob `d-check.mk` ein
`manual`-Target hat — es hat keines. Wer das Fragment einbindet, erreicht die
Hilfe des Werkzeugs über `doc-usage`, Handbuch und Spezifikation des gepinnten Image aber nur
über den selbst zusammengesetzten `docker run`-Aufruf.

**Berührte Spec-Stellen:** [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-30.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Das `--print-mk`-Fragment trägt ein vierzehntes Target `doc-manual`:
`make doc-manual MANUAL=<begriff>` ruft das gepinnte Image mit `--manual "<begriff>"` und
gibt die passenden Abschnitte von Benutzerhandbuch und Spezifikation auf **stdout** aus.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderungen an `--manual` selbst** (Treffer-Regel, Kopfzeilen, Fehlerhinweis). Das Target
  exponiert den Modus, wie `doc-usage` den Hilfe-Modus; die Anforderung dazu steht und
  ändert sich nicht (anderer Gegenstand: Werkzeug statt Fragment).
- **Eine Vorbelegung von `MANUAL`.** Ein leerer Begriff ist nach [`DC-FA-CLI-013`](../../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen) Exit 2;
  eine Vorbelegung wäre eine zweite Wahrheit neben dem Werkzeug und ließe `make doc-manual`
  etwas ausgeben, das niemand gefragt hat.
- **Ein Make-Guard für leeres `MANUAL`** (`$(error …)`). Das Werkzeug beantwortet den Fall
  bereits mit Exit 2; die Fragment-Targets `doc-immutable`/`doc-commits` führen für
  `RANGE` ebenfalls keinen eigenen Guard.
- **Weitere Targets** (`doc-spec`, `doc-readme`). Nicht gewünscht; ein Begriff genügt,
  der die Dokumente selbst wählt.
- **Release-Prep v0.89.0** (CHANGELOG, `version.md`, Pins, Handbuch-Kopf und §11): eigener
  Vorgang nach der Closure
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 17), kein Feature-Commit.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] **Fragment:** `d-check --print-mk` erzeugt das Target `doc-manual` mit `##`-Annotation,
      Recipe `@docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --manual
      "$(MANUAL)"` (Echo unterdrückt, die Ausgabe ist die Nutzlast), ohne `--enable`/`--disable`.
- [x] **Lauf:** `make doc-manual MANUAL=planning` gegen ein eingebundenes Fragment gibt
      auf stdout die Abschnitte beider Dokumente aus (Exit 0); ohne `MANUAL` Exit 2 des
      Werkzeugs, durchgereicht von make.
- [x] **Anforderung und Spiegel:** Lastenheft 0.109.0 ([`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben): vierzehn Targets,
      Akzeptanzkriterien, Out-of-Scope, Historie), Spezifikation (Punkt 5, neuer Punkt,
      Zählung, Historie).
- [x] **Test:** `TestCLI281_PrintMK_ManualTarget` ist rot ohne das Target (Rot-Beleg), grün mit.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8, kein
      Self-Review; Verifikation.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — oder
      „keine Beobachtung angefallen" in §7 notiert.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapter/driving/cli/print_mk.go` | update | Target im Template, Kopfkommentar „vierzehn" |
| `internal/adapter/driving/cli/cli_acceptance_test.go` | update | `TestCLI281_PrintMK_ManualTarget` neben dem `doc-usage`-Test |
| `spec/lastenheft.md` | update | [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben), 0.109.0, Historie |
| `spec/spezifikation.md` | update | [`DC-FA-CLI-010.a`](../../../../spec/spezifikation.md#dc-fa-cli-010a--makefile-fragment) Punkt 5, neuer Punkt, Zählung, Historie |
| `docs/user/benutzerhandbuch.md` | Release-Prep | §4.16 (Zahl, Aufzählung, Beispiel, Hinweis-Absatz) und Kopf/§11 — siehe Plan-Änderung |
| `docs/user/operations.md` | Release-Prep | Zeile `--print-mk` — siehe Plan-Änderung |
| `docs/user/abdeckung-tests.md` | regenerieren | `make abdeckung` nach dem neuen Test |

Kein ADR: derselbe Zuschnitt wie die Erweiterung um `doc-usage` (eine Anforderungs-Änderung
in [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben), keine neue Architekturentscheidung; keine ADR trägt `doc-usage`).

**Plan-Änderung vor dem Code (2026-10-30, Schritt 15):** Handbuch (§4.16, Kopf, §11) und
`docs/user/operations.md` reisen nicht im Feature-Commit, sondern in der Release-Prep —
[`AGENTS.md`](../../../../AGENTS.md) §5 Regel 17 hält Handbuch-Kopf und -Pflege am Release fest, und der
Feature-Commit fasst Handbuch und READMEs nicht an. Der Wortlaut liegt vor (am gebauten Image
gelaufen, Exit 0/2, Begriff mit Leerzeichen, Begriff mit `"` bricht die Shell) und wird dort
eingespielt. Zwischen Feature-Commit und Release-Prep nennen die Nutzer-Dokumente noch dreizehn
Targets — der Zustand des letzten Releases, bis der nächste erscheint.

**Spiegel vor dem Editieren** ([`MR-025`](../../../../harness/conventions.md#mr-025), Schritt 17;
Zählung: `grep -rniE 'dreizehn|thirteen|doc-usage' .` über das Repo ohne `.git`, Frozen-Klassen
`docs/plan/planning/done`, `docs/reviews`, `docs/plan/adr` und `CHANGELOG.md`):

- Lastenheft: Beschreibung („dreizehn `##`-annotierte Targets"), Aufzählung bei `doc-usage`/`doc-help`,
  Happy Path, Boundary, Out-of-Scope („jenseits der gelisteten dreizehn", Aufzählung),
  Kopf `Version`, Historie. Nicht betroffen: „dreizehn Quell-Tools" ([`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)).
- Spezifikation: Punkt 5 („Dreizehn `.PHONY`-Targets"), `doc-usage`-Punkt mit der Zählung
  („sechs der dreizehn tragen eine Fokus-Liste, sieben nicht"), Historie.
- Produkt: `print_mk.go` Kopfkommentar (Zeilen 22–25).
- Handbuch §4.16: Zeile 883 „dreizehn", 886 Aufzählung, Fragment-Beispiel (um 920), Hinweis-Absatz
  (um 932); `operations.md` Zeile 41.
- Nicht betroffen: README (nennt `--print-mk` nicht), `harness/`, `.claude/`, `tools/`, `Makefile`
  (nennt `print-mk` nur in Kommentaren zu eigenen Targets), die vendorte Baseline (eingefroren),
  CHANGELOG-Altbestand (Chronik).

**Zählung am erzeugten Fragment** (Schritt 18): `docker run --rm d-check:latest --print-mk |
grep -c '^\.PHONY'` ergibt heute **13**; nach der Änderung **14**. Davon tragen Fokus-Listen
(`grep -c -- '--disable'`) sechs Targets, `doc-manual` trägt keine: sechs von vierzehn, acht nicht.
Die Zahlen werden nach dem Code am Lauf gegengelesen, nicht aus der Beschreibung übernommen.

**Erkennungs-Regel (Schritt 19):** keine — es entsteht kein Muster, das entscheidet, was ein
Befund ist.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen Slice;
Priorisierung gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): Der Lauf zeigt, dass `--manual` den Mount oder das
  Anführungszeichen-Verhalten anders braucht als `--help` und das Werkzeug geändert werden
  müsste — dann trennt sich die Werkzeug-Änderung als eigener Slice ab.
- `in-progress` → `open` (blockiert): keine absehbare Blockade.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Ein Begriff mit Leerzeichen oder Sonderzeichen zerfällt in der Shell.** Das Recipe setzt
  `"$(MANUAL)"` in Anführungszeichen; ein Begriff mit `"` oder `$` bleibt eine Grenze von make
  und Shell. — **Ausgang:** entfallen — gemessen: `MANUAL="Closure-Notizen auf Substanz"` liefert Exit 0 und 5049 Byte (ein Begriff mit Leerzeichen bleibt einer); `MANUAL="Schritt 1"` kommt als ein Argument an und endet mit Exit 2, weil keine Überschrift so heißt. Ein `"` im Begriff bricht die Shell (`Unterminated quoted string`) — das ist keine Überraschung, sondern die im Lastenheft benannte Out-of-Scope-Grenze.
- **Der leere Fall wird von make als gesetzt gelesen.** `--manual ""` muss Exit 2 liefern, nicht
  den Hilfetext. — **Ausgang:** entfallen — gemessen: ohne `MANUAL` meldet das Werkzeug `--manual braucht einen Begriff`, make endet mit Exit 2, stdout bleibt leer.

## 7. Closure-Notiz

- **Was hat funktioniert:** Das Target ist die kleinste Form: ein Recipe-Eintrag nach dem Vorbild von
  `doc-usage`, am gebauten Image gelaufen statt aus der Beschreibung abgeleitet — vierzehn `.PHONY`-Targets,
  sechs mit `--disable` (heute gezählt, nicht aus dem Lastenheft übernommen). Der Rot-Beleg kam zweifach:
  der Test war rot ohne das Target (`d-check.mk ohne doc-manual-Target`), und der Verifier nahm testweise die
  Anführungszeichen um `$(MANUAL)` weg und sah den Test aus genau diesem Grund rot werden. Die Spiegel standen
  vor dem Editieren im Plan (Schritt 17); danach blieb keine alte Zahl außerhalb der deklarierten Ausnahme stehen.
- **Was ging anders als geplant:** Der Plan wies Handbuch und `operations.md` dem Feature-Commit zu, gegen
  [`AGENTS.md`](../../../../AGENTS.md) §5 Regel 17 (Pflege des Handbuchs am Release). Ich habe es vor dem ersten Commit
  bemerkt, den Wortlaut als Patch zurückgelegt und die Abweichung als Plan-Änderung vor dem Code notiert. Der Review
  fand weder Befund noch Folgearbeit (nur drei INFO, zwei davon Testkommentare, behoben). Der Verifier-LOW zum Begriff mit
  Leerzeichen war gemessen, aber noch nicht im Plan notiert; §6 trägt die Messung jetzt.
- **Steering-Loop-Eintrag:** keiner mit neuer Schwelle.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angefallen, die ein Beleg wäre. Die
  Fehlzuordnung der Handbuch-Dateien im Plan wurde vor dem Commit selbst bemerkt und ist einmalig:
  benannt, nicht gezählt. Schritt 17 wurde diesmal befolgt
  ([`guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/state.md)
  bleibt bei 4×, verkörpert).
- **Folge-Slices:** keiner. Die Release-Prep v0.89.0 ist ein eigener Vorgang nach der Closure: Handbuch §4.16,
  Kopf und §11 sowie `operations.md` Zeile 41 (Wortlaut liegt als Patch vor), dazu CHANGELOG, `version.md`
  und die Pins.
- **Risiken aus §6:** beide entfallen (siehe §6). Trigger-Audit: kein Carveout, kein bootstrap-aware Gate, keine ADR
  und keine Hard Rule mit eingetretenem Trigger; keine Schwelle gesenkt (§3.6). Nachtlauf-Stand
  ([`MR-053`](../../../../harness/conventions.md#mr-053)): wie in §8.
- **Drei Paarungen:** (a) Anker — kein `liegt in`-Feld; (b) Folge-Slice — keiner genannt; (c) Register — die
  zitierte Beobachtung existiert und trägt Belege.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind das Produkt und seine Spezifikation unter dem
Default `*` (`ALL`, GF); keine eigene Sub-Area.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; Treffer:

- [`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
  (9×, verkörpert): die Spezifikation von `doc-usage` wurde damals aus dem Anlass geschrieben, nicht
  aus dem Bestand — hier ist die Gegenprobe die Zählung am erzeugten Fragment (§3).
- [`guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/observation.md)
  (4×, verkörpert): Handbuch und `operations.md` sind Spiegel der Zahl „dreizehn" (§3 listet sie).

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`image-scan` rot am 2026-10-29 — **planmäßig gelesen**: das `HIGH` (`stdlib`, `CVE-2026-111943`,
Fix in Go 1.27.3) betraf das damalige `:latest` (v0.87.1, Go 1.27.2); ein frischer
`make image-scan` gegen das heutige `:latest` (v0.88.0, Go 1.27.3) meldet 0 `HIGH`/`CRITICAL`,
Exit 0. `upstream-drift` rot am 2026-10-29 06:40Z (vor dem Merge von PR 10): `a-check` 0.23.2 → 0.23.3
und das Basis-Image `golang:1.27.3` unter demselben Tag neu gebaut — eigene Vorgänge, nicht dieser Slice.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
