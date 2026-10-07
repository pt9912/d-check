# slice-254: Das Target `a-check` aus `a-check.mk` steht im Gate-Index

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`MR-074`](../../../../harness/conventions.md#mr-074) (Bewegung 1:
werkzeug-eigener Teil des Gate-Index — dort als Kandidat benannt),
[`DC-FA-TGT-001`](../../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(das Modul, das den Index hält; hier nur genutzt),
[ADR-0029](../../adr/0029-arch-check-via-a-check.md) (a-check als
Architektur-Gate, Fragment per `--print-mk`), Auftraggeber-Freigabe
2026-10-07.

**Berührte Spec-Stellen:** — *(Konfigurations- und Index-Arbeit, keine
Spec-Aussage)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Das einzige Target des Fragments `a-check.mk`, `a-check`, ist im
Gate-Index (`harness/README.md` §Sensors) deklariert, und der
Deklarations-Sensor (`make gate-consistency`) liest das Fragment mit
(`targets.makefiles: [Makefile, a-check.mk]`) — ein künftiges Target im
Fragment ohne Index-Zeile meldet dann `gate-undocumented`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein werkzeug-eigener Teil für a-check im Werkzeug-Verzeichnis der Baseline-Form** — die
  Baseline-Bedingung *werkzeug-eigen* verlangt, dass das Werkzeug den Teil bei
  jedem Lauf selbst schreibt; a-check tut das nicht (`--print-mk` liefert nur
  das Fragment), und `a-check.mk` ist an die Repo-Politik angepasst, gehört
  also dem Repo. Seine Targets stehen deshalb dort, wo die Baseline die Targets
  des Repos führt: in `harness/README.md` §Sensors. Tragender Anker: die
  Baseline-Regel für Werkzeug-Teile gilt einem Werkzeug, das Fragmente unter
  dem Werkzeug-Verzeichnis selbst erzeugt; ein `--print-mk`-Fragment im
  Wurzelverzeichnis führt die Baseline im Haupt-Index (wie ihr `d-check.mk`-Weg
  in der Vorlage `.d-check.yml`) — R1-F-4.
- **Disjunktheits-Schalter** (`authority-disjoint`) — mit einer einzigen
  Autoritäts-Datei wirkungslos.
- **Die übrigen Bewegungen aus dem Hebung-Eintrag** — eigene Slices (Reviewer-Regeln,
  Spezifikation).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `harness/README.md` §Sensors führt `make a-check` in der Zeile von
      `make arch-check`, mit Vertrag (Rezept aus dem Fragment, `arch-check`
      delegiert dorthin) und unveränderter Bindung.
- [ ] `.d-check.yml` `targets.makefiles` liest `a-check.mk`; der Kommentar
      sagt gemessen, was der Index deckt. *(Plan-Änderung: statt einer Zahl,
      die schon vor dem Slice veraltet war, die Aussage — nach R1-F-2.)*
      Bewusstes Brechen: ohne die
      Index-Zeile meldet `make gate-consistency` `gate-undocumented` für
      `a-check` aus dem Fragment.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      Rollenwechsel nach Schritt 8, kein Self-Review.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/README.md` | update | `make a-check` in der `arch-check`-Zeile |
| `.d-check.yml` | update | `targets.makefiles` + Kommentar |
| `harness/sensors/arch-check.md` | update, falls gemessen | nennt er das Rezept? |

**Spiegel vor dem Editieren** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
die Sensor-Datei `harness/sensors/arch-check.md`, der Kommentar zum
`targets`-Block in `.d-check.yml` (Target-Aussage), Bewegung 1 des Hebung-Eintrags (Kandidat
→ eingelöst, im Closure vermerkt).

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): das Fragment trägt mehr Targets als
  `a-check`, die einzeln entschieden werden müssen.
- `in-progress` → `open` (blockiert): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Ein künftiges `--print-mk`-Fragment bringt weitere Targets; ohne Index-Zeile
  meldet der Sensor sie — gewollt, aber bei der nächsten a-check-Hebung zu
  beachten. — **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** eine berührte Sub-Area: der Harness
des Repos (`*`, `ALL`); deklariert.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:**
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
— die Behauptung „einziges Target" wird gegen das Fragment gemessen, nicht
übernommen.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-07 — beide Läufe grün (`upstream-drift` manuell gestartet
nach den Pin-Hebungen).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
