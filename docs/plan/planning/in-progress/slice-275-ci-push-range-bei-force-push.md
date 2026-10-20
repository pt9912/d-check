# slice-275: Push-Lauf der CI bei Force-Push — nicht auflösbare Basis fällt auf die Default-Branch-Range zurück

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — reaktive Arbeit, die Closure-Bedingung ist die DoD
selbst.

**Bezug:** [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md) (die
PR-/Push-CI und ihr Traceability-Gate),
[ADR-0027](../../adr/0027-commits-traceability-modul.md) (`trace-check` über
eine Commit-Range),
[`DC-FA-COMMITS-001`](../../../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in),
[`SPEC-102`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
(`history-range-guard`: eine nicht auflösbare Range ist Exit 2).
Anlass: Befund beim Merge des Dependabot-PR #9 — jeder Force-Push eines
Branches (Dependabot `rebase`/`recreate`) macht den Push-Lauf rot, weil
`github.event.before` auf einen Commit zeigt, den der Klon nicht mehr kennt; der
Pflicht-Check desselben Namens blockt dann den Merge.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Der Schritt „Traceability + ADR-Immutable + Slice-Closure-Übergang" der
CI löst die Commit-Range eines Push-Laufs so auf, dass ein Force-Push den Lauf
nicht mehr rot färbt: Ist `github.event.before` im Klon **nicht auflösbar**, gilt
die Range gegen den Default-Branch (`origin/<default>..HEAD`) — wie schon bei
einem neuen Branch. Ein Push auf den Default-Branch selbst bleibt laut, wenn er
keine auflösbare Basis hat.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Merge von PR #9.** Er hängt an einer Policy-Umgehung oder an der
  Wirkung dieses Slice; er ist ein eigener Vorgang über den hebenden Kanal
  ([ADR-0067](../../adr/0067-dependabot-als-hebender-kanal.md)).
- **Eine Lockerung von `history-range-guard`.** Der Wächter hält die
  nicht auflösbare Range laut fest ([`SPEC-102`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)) und bleibt unverändert; der Slice
  ändert, **welche** Range die CI ihm gibt, nicht was er mit einer
  nicht auflösbaren tut.
- **Die Pflicht-Check-Konfiguration auf GitHub** (Ruleset, Branch Protection).
  Sie liegt außerhalb des Repos und ist aus dem Klon nicht auditierbar
  ([`harness/README.md`](../../../../harness/README.md) §Gate-Taxonomie).
- **Weitere `ci.yml`-Schritte.** Nur der Range-Schritt; `make ci` und die
  Release-Pipeline bleiben stehen.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] **Range auflösen:** Die Range-Auflösung des CI-Schritts liegt in einem
      Skript unter `tools/harness/`, das ein Push-`before`, das nicht als Commit
      auflösbar ist, auf `origin/<default>..HEAD` zurückführt und nur dann; der
      Schritt in `ci.yml` ruft es.
- [x] **Beleg:** Eine Probe fährt das Skript gegen ein Wegwerf-Repo für die
      Fälle neuer Branch, normaler Push, Force-Push (Basis fehlt) und Push auf
      den Default-Branch ohne auflösbare Basis; der letzte Fall endet laut, die
      anderen mit der erwarteten Range.
- [x] **Entscheidung:** Eine ADR hält fest, dass die CI-Range sich ändert, was
      der Rückfall zusagt und was nicht (Force-Push auf einen Feature-Branch
      prüft alles, was der Branch gegenüber dem Default-Branch trägt; auf dem Default-Branch bleibt die
      Nicht-Auflösbarkeit ein Fehler). Kopf und Index stehen auf `Accepted`.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review; Verifikation.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — oder
      „keine Beobachtung angefallen" in §7 notiert.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/ci.yml` | update | Der Range-Schritt ruft das Skript statt die Auflösung inline zu tragen |
| `tools/harness/` Skript zur Range-Auflösung | neu | testbar außerhalb von GitHub Actions; trägt die vier Fälle |
| `tools/harness/` Probe des Skripts | neu | Wegwerf-Repo; Proben je Fall, auch der laute Fall |
| `harness/README.md`, `harness/sensors/` | update | Werkzeug-Zeile bzw. Sensor-Seite, falls ein neues `make`-Target entsteht (`make gate-consistency` hält die Tabelle) |
| `docs/plan/adr/0109-…md`, `docs/plan/adr/README.md` | neu / update | die Entscheidung; Index-Eintrag |

**Spiegel vor dem Editieren auflisten**
([`MR-025`](../../../../harness/conventions.md#mr-025), Schritt 17): wer die
Range-Auflösung des CI-Schritts verschiebt, greppt vorher `PUSH_BEFORE`,
`PR_BASE`, `origin/$DEFAULT_BRANCH` und „Commit-Range" über `.github/`,
`harness/`, `docs/user/`, `docs/maintainer/`, `.githooks/` und die Skript-Kommentare
und listet jeden Treffer hier. Der lokale `pre-commit`/`commit-msg`-Hook liest
eine andere Range und ist **kein** Gegenstand — er wird nur gelesen, nicht
geändert.

**Offene Frage für den Plan-Lauf:** Ob das Skript eine eigene Datei bekommt oder
in `history-range-guard.sh` einen Modus erhält. Gegen den Modus spricht, dass der
Wächter bewusst nur prüft und nichts auflöst (siehe Abgrenzung); der Plan-Lauf
legt es vor dem Code fest.

**Plan-Änderung vor dem Code (2026-10-20, Schritt 15):** Die offene Frage ist
entschieden — eigene Dateien: `tools/harness/ci-range.sh` (Auflösung) und
`tools/harness/ci-range-probe.sh` (Probe); `history-range-guard.sh` bleibt
unverändert. Dazu kommen zwei Dateien, die der Plan nicht nannte: das
`Makefile` (Target `ci-range-probe` samt `.PHONY`, damit `make gate-consistency`
die Zeile in `harness/README.md` mit einem Target deckt) und das Workflow
`ci.yml` reicht `github.ref_name` als `REF_NAME` an das Skript durch, damit der
Default-Branch erkennbar ist. Spiegel-Grep (`PUSH_BEFORE`, `PR_BASE`,
`github.event.before`, „Pruefe Commit-Range") über `.github/`, `harness/`,
`docs/user/`, `docs/maintainer/`, `.githooks/`, `tools/`, `Makefile` und die
READMEs: Treffer nur in `ci.yml`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen Slice;
Priorisierung gesetzt (`open` → `next`, `Verantwortlich:`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): Die Probe braucht eine eigene
  Test-Infrastruktur über Skript und Wegwerf-Repo hinaus — dann trennt sich die
  Probe als eigener Slice ab.
- `in-progress` → `open` (blockiert): Ein Rückfall auf die Default-Branch-Range
  lässt sich nicht von „still grün über leerem Prüfbereich" unterscheiden
  (`BEO-ALL/stilles-gruen-ueber-leerer-range`) — dann zuerst klären, ob die
  Nicht-Auflösbarkeit auf dem Default-Branch laut bleibt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft; die ADR steht auf `Accepted`. Ein Dependabot-Rebase
färbt den Push-Lauf danach nicht mehr rot (am nächsten Rebase beobachtbar).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Rückfall macht aus „Basis fehlt" ein stilles Grün.** — **Ausgang:**
  entfallen — auf dem Default-Branch endet der Fall mit Exit 2 (Probe
  `force-push-default`), und eine leere Rückfall-Range fängt
  `history-range-guard`, den `make trace-check` und `make adr-check` zuerst
  rufen. Die Verifikation nahm die Ausnahme zurück und sah die Probe rot
  (Exit 0 über leerer Range).
- **Ein Rückfall prüft mehr Commits als der Push brachte.** — **Ausgang:**
  entfallen — der Rückfall prüft, was der Branch gegenüber dem Default-Branch
  trägt, in der Regel mehr, nach einem Rebase auch weniger als `before..head`
  (die Verifikation maß die Range nach einem Rebase auf einen neueren
  Default-Stand: nur die Branch-eigenen Commits). Altlast eines Feature-Branchs
  meldet `trace-check` als Befund der Historie, nicht des Rückfalls.
- **Die Gate-Änderung ist eine Senkung** (`AGENTS.md` §3.6). — **Ausgang:**
  entfallen — die Entscheidung steht in
  [ADR-0109](../../adr/0109-ci-push-range-faellt-bei-force-push-zurueck.md),
  `Accepted`, mit Re-Evaluierungs-Trigger; sie lockert keine Schwelle, sondern
  wählt die Range für einen sonst unauflösbaren Fall.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Range-Auflösung des CI-Schritts liegt in
  `tools/harness/ci-range.sh`; ein Force-Push auf einen Feature-Branch fällt auf
  `origin/<default>..HEAD` zurück, auf dem Default-Branch und ohne
  `origin/<default>` endet er laut mit Exit 2. `make ci-range-probe` fährt sechs
  Proben und wurde an allen drei Bruchstellen (Rückfall, Default-Ausnahme,
  origin-Wächter) einzeln rot gesehen — vom Implementer und, unabhängig, vom
  Verifier. Unter `set -euo pipefail` bricht `RANGE="$(…)"` bei Exit 2 des
  Skripts den Schritt ab (im Wegwerf-Repo mit gelöschtem `before`-Commit
  belegt).
- **Was ging anders als geplant:** Der Plan nannte weder das `Makefile` (Target
  samt `.PHONY` für `make gate-consistency`) noch den Durchreich-Wert
  `REF_NAME` in `ci.yml`; beides stand in einer Plan-Änderung vor dem Code. Die
  Zusage „mehr Commits, nie weniger" war in Skript, ADR und Plan zu weit
  (Review INFO F-1) und wurde vor `Accepted` berichtigt; die ADR-Fitness
  zählte erst „fünf Fälle und eine Lücke" (Verifikation LOW) und nennt jetzt
  sechs Proben. **Nicht belegt:** `make trace-check`/`adr-check` mit der
  Fallback-Range im Docker-Image und der Fetch-Pfad; der Fetch-Fehler wird
  verschluckt und die Existenzprüfung der Ref ist kein Frische-Beleg (Review
  F-2) — die Folge ist eine größere Range, kein stilles Grün.
- **Steering-Loop-Eintrag:** keiner mit neuer Schwelle.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  angefallen. Der rote Push-Lauf ist ein weiterer Fall „Gate rot, ohne dass die
  geprüfte Eigenschaft verletzt ist", gehört aber nicht unter
  `BEO-ALL/gate-rot-aus-falschem-grund` (dort: eine Probe in einem
  ausgenommenen Abschnitt) — ein Beleg dort wäre eine erzwungene Zuordnung.
- **Folge-Slices:** keiner. Der Merge von PR #9 ist ein eigener Vorgang über
  den hebenden Kanal.
- **Risiken aus §6:** alle drei entfallen (siehe §6). Trigger-Audit: kein
  Carveout, kein bootstrap-aware Gate; [ADR-0109](../../adr/0109-ci-push-range-faellt-bei-force-push-zurueck.md) ist mit diesem Slice `Accepted`
  und trägt ihren Re-Evaluierungs-Trigger; keine Hard Rule mit eingetretenem
  Trigger. Nachtlauf-Stand
  ([`MR-053`](../../../../harness/conventions.md#mr-053)): wie in §8.
- **Drei Paarungen:** (a) Anker — kein Eintrag mit `liegt in`; (b)
  Folge-Slices — keiner genannt; (c) Register — keine zitierte Beobachtung
  außer den in §8 gesichteten, die existieren.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind die CI-Konfiguration, ein
Harness-Skript unter `tools/harness/` und die ADR-Ablage. `tools/harness/` ist
eine eigene Sub-Area (`HARN`, GF); die übrigen laufen unter dem Default `*`
(`ALL`). Beide GF.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; Treffer:

- [`stilles-gruen-ueber-leerer-range`](../observations/BEO-ALL/stilles-gruen-ueber-leerer-range/observation.md)
  (1×): der Rückfall darf nie aus „Basis fehlt" ein Grün über leerem
  Prüfbereich machen — das ist das zentrale Risiko dieses Slice.
- [`gate-rot-aus-falschem-grund`](../observations/BEO-ALL/gate-rot-aus-falschem-grund/observation.md)
  (3×): der rote Push-Lauf ist ein weiterer Fall — ein Gate wird rot, ohne
  dass die geprüfte Eigenschaft verletzt ist.
- [`range-leerfall-mit-lokalem-branchnamen`](../observations/BEO-ALL/range-leerfall-mit-lokalem-branchnamen/observation.md)
  (1×): Range-Auflösung mit lokalen Namen; die Probe läuft gegen ein Wegwerf-Repo,
  nicht gegen den lokalen Klon.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`image-scan` grün; `upstream-drift` rot, **planmäßig gelesen** (neuer
Baseline-Release und veraltete Pins für Go, golangci-lint, semgrep und a-check) —
kein Bezug zu diesem Slice, jeweils eigener Vorgang.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
