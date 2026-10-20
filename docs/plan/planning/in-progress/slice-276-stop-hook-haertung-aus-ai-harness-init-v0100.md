# slice-276: Stop-Hook und Gate-Nachweis auf die Fassung von ai-harness-init v0.10.0 — streng gelassen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`MR-004`](../../../../harness/conventions.md#mr-004) (Gate-Nachweis-Mechanik
und Hooks), [`MR-005`](../../../../harness/conventions.md#mr-005) (Härtung:
inhaltsbasierter Nachweis), [`MR-076`](../../../../harness/conventions.md#mr-076)
(Nachweis im Rezept, Stop-Hook blockt bei unlesbarem Zustand),
[`SPEC-094`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
(das Handoff-Gate), [`AGENTS.md`](../../../../AGENTS.md) §3.6 (Gates werden nur
per ADR gelockert). Anlass: `ai-harness-init` v0.10.0 liefert den Stop-Hook mit
einem HEAD-Stempel und einem Fail-closed-`trap`, dessen **Default** das Handoff-Gate
lockert; die Entscheidung des Auftraggebers lautet: übernehmen, **streng lassen**.

**Berührte Spec-Stellen:** [`SPEC-094`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Stop-Hook, Gate-Nachweis und Selbstprüfung tragen die Mechanik von
`ai-harness-init` v0.10.0 — HEAD-Stempel, `trap` auf Exit 2, strenger Zweig,
Stop-Hook-Zweige in der Selbstprüfung —, und der Hook **ist ohne jeden Schalter
streng**: das Verhalten des Hooks bleibt das heutige, jede
Inhaltsänderung nach dem letzten grünen Lauf blockt weiter.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Lockerung des Defaults.** Der Upstream-Default lässt ein Turn-Ende ohne
  neuen HEAD auch mit ungedeckter Änderung frei. Das ist eine Lockerung eines
  Handoff-Gates und braucht nach `AGENTS.md` §3.6 eine ADR und die ausdrückliche
  Entscheidung des Auftraggebers; sie ist nicht getroffen. Wer sie will, legt
  `.harness/stop-gate-locker` an, per eigener ADR.
- **`slice-mv.sh`.** Es weicht in beide Richtungen ab (unsere Unterordner-Anpassung und
  Stufen-Kette, Upstreams `quelle_finden`); die Übernahme ist ein Merge und ein
  anderer Vorgang an einem anderen Werkzeug.
- **`extract-command.awk` und `pretooluse-command-guard.sh`.** Unsere Fassungen
  sind weiter als Upstream (Strings ohne Trenner, verkürzte Schlüssel,
  Host-Interpreter-Sperre); sie bleiben. Der Rückfluss nach Upstream wäre ein
  Change Request, kein Teil dieses Slice.
- **Die Baseline-Hebung auf v6.18.1** und die neuen Werkzeuge (`e2e-abdeckung`,
  `traeger-fetch`, `span-*`). Eigene Vorgänge mit eigenem `MR`.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] **Hook und Nachweis:** `stop-require-gates.sh` und `record-gates.sh` tragen
      HEAD-Stempel, `trap` auf Exit 2 und den lockeren Zweig, der nur mit
      `.harness/stop-gate-locker` greift; die Datei gibt es nicht. Das Verhalten ohne Variable
      ist das strenge: ein Turn-Ende mit ungedecktem Inhalt blockt, auch ohne
      neuen HEAD.
- [x] **Selbstprüfung:** `selbstpruefung.sh` trägt die Stop-Hook-Zweige und die
      Marker der Fassung v0.10.0, belegt für d-check; `make selbstpruefung` läuft
      grün und beobachtet beide Ausgänge der Stop-Hook-Prüfung.
- [x] **Entscheidung und Spiegel:** Ein neuer `MR` hält fest, was übernommen ist
      und was nicht (der lockere Default), und schärft
      [`MR-076`](../../../../harness/conventions.md#mr-076); [`SPEC-094`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) und die
      Sensor-Zeile in `harness/README.md` sind nachgezogen.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8, kein
      Self-Review; Verifikation.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — oder
      „keine Beobachtung angefallen" in §7 notiert.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.claude/hooks/stop-require-gates.sh` | update | Fassung v0.10.0: `trap`, `head_wert`, strenger Zweig; Meldungstexte bleiben die der Fassung |
| `tools/harness/record-gates.sh` | update | schreibt den HEAD-Stempel `.harness/state/gates-passed.head` neben den Hash; `head_wert` genauso wie der Hook (Kopplung) |
| `.harness/stop-gate-locker` | — (nicht angelegt) | Lockerungs-Schalter; dieses Repo legt ihn nicht an, Streng ist der Default ohne Schalter |
| `tools/harness/selbstpruefung.sh` | update | Fassung v0.10.0 mit den Markern und den Stop-Hook-Zweigen; Belegung für d-check (`.githooks/commit-msg`, `make hooks`, `make gates`) |
| `harness/conventions/MR-<neu>-….md`, `harness/conventions.md` | neu / update | die Übernahme samt Grenze; Index-Zeile |
| `spec/spezifikation.md` | update | [`SPEC-094`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) und Historie: HEAD-Stempel, strenger Modus, Fail-closed |
| `harness/README.md` | update | Sensor-Zeile `selbstpruefung`, Grenzen-Absatz zum Stop-Hook |

**Spiegel vor dem Editieren auflisten**
([`MR-025`](../../../../harness/conventions.md#mr-025), Schritt 17): vor dem Code
`grep` nach `gates-passed`, `stop-require-gates`, `record-gates`,
`stop_hook_active` über `harness/`, `docs/user/`, `docs/maintainer/`, `.claude/`,
`tools/`, `README.md`, `README.de.md`, `AGENTS.md`, `Makefile` und `spec/`; Treffer
beim Anlegen: `harness/conventions.md`, `harness/README.md`, [`MR-004`](../../../../harness/conventions.md#mr-004), [`MR-076`](../../../../harness/conventions.md#mr-076),
`.claude/settings.json`, `.claude/commands/plan-welle.md`,
`.claude/commands/close-welle.md`, `Makefile`, `spec/spezifikation.md`. Der Plan-Lauf
liest jeden und listet, ob er den alten Mechanismus beschreibt.

**Plan-Änderung nach Review-Runde 1 (2026-10-20, Schritt 20):** Der Report
`docs/reviews/2026-10-20-slice-276-stop-hook-streng.md` meldet ein HIGH
(`schalter-datei-ohne-halter`): der strenge Modus hing allein an der Datei
`.harness/stop-gate-streng`, und ihre Entfernung (`rm`, `git stash`,
Branch-Wechsel, toter Symlink) lockert das Gate ohne Signal — das Plan-Risiko
„Datei als einziger Halter" ist eingetreten. Der Schnitt ändert sich, die
Abgrenzung in §1 bleibt.

- **Neu: Streng ist der Default, das Lockern der ausdrückliche Schalter.** Der Hook
  ist streng, **solange** die Datei `.harness/stop-gate-locker` fehlt; nur sie
  schaltet den lockeren Zweig frei, und `STOP_GATE_STRENG=1` hebt sie auf. Eine
  gelöschte oder fehlende Datei kann damit nichts lockern; wer lockern will, muss
  eine Datei anlegen, und das ist eine ADR-pflichtige Entscheidung ([MR-079](../../../../harness/conventions.md#mr-079)). Die
  Datei `.harness/stop-gate-streng` entfällt.
- **Betroffen:** `.claude/hooks/stop-require-gates.sh` (Modus-Bildung),
  `tools/harness/selbstpruefung.sh` (Abschnitt 5 misst: streng ohne Datei, locker mit
  `.harness/stop-gate-locker`, streng mit Variable trotz Datei, Commit ohne Nachweis),
  `.harness/stop-gate-streng` (gelöscht), [`MR-079`](../../../../harness/conventions.md#mr-079), [`SPEC-094`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge), [`SPEC-103`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge),
  `harness/README.md`. Alles in §3 genannt, außer dass die Datei jetzt der
  Lockerungs-Schalter ist.
- **Aussagen-Suche:** `grep -rn "stop-gate-streng"` über `.claude/`, `tools/`,
  `harness/`, `spec/`, `docs/plan/planning/in-progress/` und `.harness/`: Treffer in
  Hook, `record-gates.sh`, `selbstpruefung.sh`, [`MR-079`](../../../../harness/conventions.md#mr-079), `harness/README.md`,
  [`SPEC-094`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge), [`SPEC-103`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) und diesem Plan — jeder wird nachgezogen.
- **INFO 1 des Reports:** Die Selbstprüfung unterscheidet im Abschnitt „Commit ohne
  Nachweis", über welchen Pfad der Hook blockt (streng: über den Hash; locker: neuer
  `HEAD` plus abweichender Hash) und behauptet nicht beides für beide Zweige.
- **Zusätzlich:** Der Fix schreibt einen neuen Modus-Weg; nach ihm läuft eine zweite
  `reviewer`-Runde über genau die Fix-Commits (Schritt 14).

**Rot-Beleg** (Modul 11 §Bewusstes Brechen): Der strenge Zweig und das `trap`
werden testweise entfernt; `make selbstpruefung` bzw. ein Wegwerf-Repo-Lauf muss
dann aus dem richtigen Grund rot werden (ein Turn-Ende mit ungedecktem Inhalt gibt
frei, ein unerwarteter Fehler endet nicht mit Exit 2).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen Slice;
Priorisierung gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): Die Übernahme der Selbstprüfung zieht mehr
  Anpassung nach sich als die Marker (etwa Netz- oder Docker-Annahmen der neuen
  Fassung) — dann trennt sich die Selbstprüfung als eigener Slice ab; der Hook-Teil
  bleibt hier. Wächst der Diff über eine Review-Sitzung, gilt die Form von
  [`MR-066`](../../../../harness/conventions.md#mr-066).
- `in-progress` → `open` (blockiert): Der strenge Modus verhält sich in der
  Selbstprüfung nicht wie der heutige Hook — dann zuerst klären, was der Hook
  heute tut, statt die Selbstprüfung anzupassen.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft; `make selbstpruefung` ist grün.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der strenge Modus verhält sich nicht wie der heutige Hook.** Eine
  Unterschiedlichkeit wäre eine stille Lockerung. — **Ausgang:** entfallen — die
  Verifikation fuhr 18 Lagen im Wegwerf-Repo (die neun geforderten und acht weitere,
  darunter Stempel gleich `HEAD`, gelöschter Stempel, Nachweis als Verzeichnis oder
  ohne Leserecht) gegen den Hook von `2032fe1d` im Default: deckungsgleich. Die
  einzige Abweichung ist eine Verschärfung (außerhalb eines Repositorys Exit 2 statt
  Exit 1, ohne `git` im `PATH` ein Block mit Grund statt Exit 127).
- **Eine Datei `.harness/stop-gate-streng` ist das Einzige, was den Default
  hält.** Fehlt sie in einem frischen Klon oder wird sie gelöscht, gilt der lockere
  Default. — **Ausgang:** eingetreten und im Design behoben — der Review der ersten
  Fassung fand genau das (Datei weg ⇒ Hook locker); streng ist jetzt der Default ohne
  Schalter, nur eine angelegte `.harness/stop-gate-locker` lockert (Plan-Änderung
  oben, [`MR-079`](../../../../harness/conventions.md#mr-079)). Die Datei
  `.harness/stop-gate-streng` entfällt; Löschung, Branch-Wechsel, `git stash` und
  toter Symlink lockern nichts (Probe der Verifikation).
- **Die neue `selbstpruefung.sh` setzt Annahmen (Netz, Docker-Builds im Klon),
  die die alte nicht setzte.** — **Ausgang:** entfallen — der Abschnitt (5) braucht
  nichts über `git`, `bash` und den Hook im Klon hinaus; der Lauf dauert gemessen
  rund zwei Minuten wie vorher, getragen vom `make gates` im Klon.

## 7. Closure-Notiz

- **Was hat funktioniert:** Stop-Hook, `record-gates` und Selbstprüfung tragen die
  Mechanik von `ai-harness-init` v0.10.0 — HEAD-Stempel, `trap` auf Exit 2, einen
  lockeren Zweig, die Stop-Hook-Zweige in `make selbstpruefung` —, und der Hook ist
  **ohne jeden Schalter streng**. Im Default entscheidet er in 18 Lagen wie der alte;
  die Selbstprüfung beobachtet die vier Zweige (streng ohne Schalter, locker mit der
  Datei, streng mit der Variable trotz der Datei, Commit ohne Nachweis). Rot-Belege:
  ohne den Default-Zweig, ohne die Aufhebung durch die Variable und ohne `trap` wird je
  die benannte Lage rot (Verifikation).
- **Was ging anders als geplant:** (1) Meine Empfehlung, den HEAD-Stempel löse das
  wiederholte „Run `make gates` again", war zu optimistisch: der Upstream-Default ist
  eine Lockerung, der strenge Modus verhält sich wie unser Hook — der Nutzen ist
  Härtung, keine Entlastung. Der Auftraggeber entschied: übernehmen, streng lassen.
  (2) Die erste Fassung hielt den strengen Modus über die Datei
  `.harness/stop-gate-streng`; der Review (HIGH) maß, dass ihre Entfernung das Gate
  still lockert. Der Schnitt wurde umgekehrt — streng ohne Schalter, locker nur mit
  `.harness/stop-gate-locker` —, daraufhin lief Runde 2 über die Fix-Commits.
  (3) **Nicht belegt:** das `trap` misst keine Zeile der Selbstprüfung (nur die
  Wegwerf-Probe der Verifikation), und die Selbstprüfung sieht eine unversionierte
  Lockerungs-Datei im Arbeitsbaum nicht (sie prüft den Klon von `HEAD`); beides
  steht in [`MR-079`](../../../../harness/conventions.md#mr-079).
- **Steering-Loop-Eintrag:** keiner mit neuer Schwelle.
- **Beobachtungs-Register (`../observations/`):** neu
  [`BEO-ALL/schalter-datei-ohne-halter`](../observations/BEO-ALL/schalter-datei-ohne-halter/state.md)
  (Beleg `evidence/slice-276.md`, 1×); `evidence/slice-276.md` in
  [`BEO-ALL/adopter-fix-lockert-waechter-ohne-adr`](../observations/BEO-ALL/adopter-fix-lockert-waechter-ohne-adr/state.md)
  (jetzt 2×): die Upstream-Lockerung fiel diesmal vor dem Code auf.
- **Folge-Slices:** keiner. Die Baseline-Hebung auf v6.18.1 und die Tool-Pins sind
  eigene Vorgänge.
- **Risiken aus §6:** alle drei aufgelöst (siehe §6). Trigger-Audit: kein
  Carveout, kein bootstrap-aware Gate, keine ADR; keine Hard Rule mit eingetretenem
  Trigger. Nachtlauf-Stand
  ([`MR-053`](../../../../harness/conventions.md#mr-053)): wie in §8.
- **Drei Paarungen:** (a) Anker — kein Eintrag mit `liegt in`; (b)
  Folge-Slices — keiner genannt; (c) Register — die zitierten Beobachtungen
  existieren.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `tools/harness/`
(`HARN`, GF, konventionsgetragen über
[`MR-004`](../../../../harness/conventions.md#mr-004)) und die Werkzeug-Konfiguration
unter `.claude/` und `.harness/` sowie die Doku — unter dem Default `*` (`ALL`).
Beide GF.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; Treffer:

- [`adopter-fix-lockert-waechter-ohne-adr`](../observations/BEO-ALL/adopter-fix-lockert-waechter-ohne-adr/observation.md)
  (1×): genau die Lage dieses Slice — die Übernahme einer Fremd-Fassung brächte
  einen lockeren Default mit; deshalb ist der strenge Modus gesetzt und die
  Lockerung ausgeschlossen (§1).
- [`haertung-kippt-fehlerpolitik-ungeprueft`](../observations/BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft/observation.md)
  (5×, verkörpert, Prüffrage 19 und Schritt 14): ein neuer Leseweg einer Härtung
  wird gegen seine Fehlerformen gefahren; der Review läuft über den Code-Commit.
- [`pipe-swallows-gate-exit-code`](../observations/BEO-ALL/pipe-swallows-gate-exit-code/observation.md)
  (4×): `trap`, `pipefail` und Exit-Codes des Hooks sind genau dieser Fall.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`image-scan` grün; `upstream-drift` rot, **planmäßig gelesen** (neuer
Baseline-Release `v6.18.1` und veraltete Pins für Go, golangci-lint, semgrep,
a-check) — kein Bezug zu diesem Slice, die Baseline-Hebung ist ein eigener Vorgang.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
