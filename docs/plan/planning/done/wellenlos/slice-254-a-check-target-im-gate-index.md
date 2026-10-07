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

- [x] `harness/README.md` §Sensors führt `make a-check` in der Zeile von
      `make arch-check`, mit Vertrag (Rezept aus dem Fragment, `arch-check`
      delegiert dorthin) und unveränderter Bindung.
- [x] `.d-check.yml` `targets.makefiles` liest `a-check.mk`; der Kommentar
      sagt gemessen, was der Index deckt. *(Plan-Änderung: statt einer Zahl,
      die schon vor dem Slice veraltet war, die Aussage — nach R1-F-2.)*
      Bewusstes Brechen: ohne die
      Index-Zeile meldet `make gate-consistency` `gate-undocumented` für
      `a-check` aus dem Fragment.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      Rollenwechsel nach Schritt 8, kein Self-Review.
- [x] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
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
  beachten. — **Ausgang:** entfallen — gewollt laut: ein neues Fragment-Target ohne Index-Zeile meldet `gate-undocumented` (Brech-Probe des Verifiers); die a-check-Hebung prüft es ohnehin mit `--print-mk` gegen das Fragment.

## 7. Closure-Notiz

- **Was hat funktioniert:** Der Rot-Beleg ging der Änderung voraus: nur die
  Konfiguration erweitert, und `make gate-consistency` meldete
  `a-check.mk:57 a-check gate-undocumented`; mit der Index-Zeile 0 Befunde.
  Der Verifier brach in drei Richtungen (Zeile weg, erfundenes Target, alte
  Konfiguration ⇒ `gate-phantom`) und bestätigte jede Aussage gegen den Code.
- **Was ging anders als geplant:** Der Plan wollte die Target-Zahl im
  Kommentar „gemessen" nachziehen — gemessen war sie schon vorher falsch (54
  statt 57); sie ist durch eine Aussage ersetzt (Plan-Änderung nach R1-F-2).
  Und der neu geschriebene Kommentar belegte sich zunächst mit dem gekürzten
  Vorlagen-Zitat, dessen ausgelassene Mitte in `v6.17.0` das Gegenteil sagt
  (R1-F-1) — dieselbe Klasse wie im Hebung-Slice (R1-F-3), jetzt im eigenen
  neuen Text. `harness/sensors/arch-check.md` blieb unverändert, weil er
  `a-check.mk` bereits nennt (R1-F-5).
- **Steering-Loop-Eintrag:** keine neue Verkörperung —
  [`BEO-ALL/citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/state.md)
  (ein Zitat trägt mehr, als sein Geltungsbereich hergibt) bekommt einen
  Beleg.
- **Beobachtungs-Register (`../observations/`):** Evidence `slice-254` unter
  dem genannten Eintrag.
- **Folge-Slices:** keine.
- **Risiken aus §6:** entfallen — ein neues Fragment-Target ohne Index-Zeile
  meldet laut, die nächste a-check-Hebung prüft das Fragment ohnehin.
  Trigger-Audit: kein Carveout, kein bootstrap-aware Gate, keine ADR und keine
  Hard Rule mit eingetretenem Trigger; [MR-074](../../../../harness/conventions.md#mr-074) Bewegung 1 als eingelöst
  vermerkt. Nachtlauf-Stand
  ([`MR-053`](../../../../harness/conventions.md#mr-053)): wie in §8.
- **Drei Paarungen:** (a) Anker — kein `liegt in`-Feld; (b) Folge-Slice —
  keine; (c) Register — die zitierte Beobachtung existiert und trägt Belege.

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
