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

- [ ] **Range auflösen:** Die Range-Auflösung des CI-Schritts liegt in einem
      Skript unter `tools/harness/`, das ein Push-`before`, das nicht als Commit
      auflösbar ist, auf `origin/<default>..HEAD` zurückführt und nur dann; der
      Schritt in `ci.yml` ruft es.
- [ ] **Beleg:** Eine Probe fährt das Skript gegen ein Wegwerf-Repo für die
      Fälle neuer Branch, normaler Push, Force-Push (Basis fehlt) und Push auf
      den Default-Branch ohne auflösbare Basis; der letzte Fall endet laut, die
      anderen mit der erwarteten Range.
- [ ] **Entscheidung:** Eine ADR hält fest, dass die CI-Range sich ändert, was
      der Rückfall zusagt und was nicht (Force-Push auf einen Feature-Branch
      prüft mehr Commits, nie weniger; auf dem Default-Branch bleibt die
      Nicht-Auflösbarkeit ein Fehler). Kopf und Index stehen auf `Accepted`.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review; Verifikation.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — oder
      „keine Beobachtung angefallen" in §7 notiert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
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

- **Der Rückfall macht aus „Basis fehlt" ein stilles Grün.** Auf dem
  Default-Branch bliebe sonst eine leere Range laut nur, wenn der Wächter sie
  sieht. — **Ausgang:** offen bis Closure; Probe „Push auf den Default-Branch
  ohne auflösbare Basis" trägt es.
- **Ein Rückfall prüft mehr Commits als der Push brachte.** Auf einem
  Feature-Branch mit alter, schon geprüfter Historie meldet `trace-check` dann
  Altlast. — **Ausgang:** offen bis Closure.
- **Die Gate-Änderung ist eine Senkung.** `AGENTS.md` §3.6: die Entscheidung
  gehört in eine ADR, kein PR-Kommentar. — **Ausgang:** offen bis Closure; die
  ADR ist Teil der DoD.

## 7. Closure-Notiz

<!-- Wird bei der Closure gefüllt, vor dem `git mv` nach `done/`. -->

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** keiner mit neuer Schwelle.
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

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
