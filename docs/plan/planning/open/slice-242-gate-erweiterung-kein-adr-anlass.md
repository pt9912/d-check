# slice-242: Adoption „Gate-Erweiterung ist kein ADR-Anlass" (modul-04, v6.13.0) in AGENTS.md §3.6

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle as State Machine.

**Welle:** ohne Welle — die Klarstellung einer Hard Rule um einen
Ergänzungs-Fall ist kein repo-weites Mehr (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** die Hebung
[`MR-073`](../../../../harness/conventions.md#mr-073--baseline-pin-hebung-auf-v6130-fünfzehnter-nachtrag-zu-mr-011-nachtrag-zu-mr-023)
liefert das Delta (Baseline-Regelwerk `modul-04-adrs.md` v6.13.0:
„Eine Gate-Erweiterung ist nicht automatisch ein ADR-Anlass"),
[`MR-051`](../../../../harness/conventions.md#mr-051),
[`MR-054`](../../../../harness/conventions.md#mr-054),
[`MR-053`](../../../../harness/conventions.md#mr-053).
Keine `DC-*` — keine Produkt-Anforderung berührt.

**Berührte Spec-Stellen:** — *(AGENTS.md ist kein Spec-Stratum;
`spec/` bleibt unverändert)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-29.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Hard Rule `AGENTS.md` §3.6 („Gates dürfen nicht ohne ADR
gelockert werden") um den Baseline-Ergänzungs-Fall erweitern
(`modul-04-adrs.md` v6.13.0): die Aufnahme eines **bereits existierenden,
unabhängig lauffähigen Wächters** in `make gates` ist **keine neue
Entscheidung** — ein Verweis auf die ADR genügt, die den Wächter
ursprünglich trägt; existiert keine, trägt schon die *Einführung* des
Wächters eine. Eine **neue Fehlerklasse**, ein **neuer Scope** oder ein
**Widerspruch** zu einer bestehenden ADR braucht weiterhin eine eigene.
Im Zweifel: ADR.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Keine Lockerung des Senkungs-Verbots** — Satz 1 („Jede
  Schwellen-Senkung … ist ein ADR, kein PR-Kommentar") bleibt byte-fest;
  die Ergänzung ist additionale, keine Umformulierung.
- **Keine Rückwirkung auf bestehende Gates** — die Zusage gilt ab
  Veröffentlichung für künftige Aufnahmen; bestehende Bindungen (Tabelle
  `harness/README.md` §Sensors) bleiben, wie sie sind.
- **Keine Terminologie-Übernahme ungeprüft** — der Baseline-Wortlaut sagt
  „PR-blockierenden Satz"; d-checks `make gates` ist lokal und CI-geprüft
  (§3.1, `harness/README.md` §Sensors) — die Übertragung auf den eigenen
  Betrieb wird im Vollzug formuliert, nicht zitiert.

## 2. Definition of Done

- [ ] `AGENTS.md` §3.6 trägt die Klarstellung: Aufnahme eines existierenden
      Wächters = Verweis auf die tragende ADR genügt; neue Fehlerklasse /
      neuer Scope / Widerspruch = eigene ADR; „Im Zweifel: ADR".
- [ ] Der Senkungs-Satz (§3.6, erste Aussage) ist byte-identisch geblieben
      — Gegenprobe per `git diff` am Commit.
- [ ] `make gates` grün.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §3.6 | update | eine Ergänzung nach dem Senkungs-Satz, vor der Grenze-Notiz |
| dieser Plan | update | etwaige Anker-Korrekturen, falls der Vollzug Zeilen verschiebt |

## 4. Trigger

**Start** (`open` → `in-progress`): direkt beansprucht — die
Auftraggeber-Freigabe „ja bitte" zu den vier Deltas der v6.13.0-Hebung
(2026-09-29) ist die Beanspruchung selbst; der Nachtlauf-Stand wird bei
der Beanspruchung gelesen ([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): keiner erwartet — eine Ergänzung an
  einer Stelle.
- `in-progress` → `open` (Blocker): keiner bekannt; der Delta-Wortlaut ist
  vendiert und gelesen.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag (drei Paarungen
wellenlos hier geprüft).

## 6. Risiken und offene Punkte

- Die Ergänzung lockert versehentlich das Senkungs-Verbot — DoD 2
  (Gegenprobe `git diff`) deckt. — **Ausgang:** *(offen)*
- „unabhängig lauffähig" ist ein Urteil, kein Muster — die Grenze wird im
  Vollzug am eigenen Bestand belegt (welche Targets der Klasse schon
  angehören). — **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das
Harness-Werkzeug selbst (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** keine berührte Sub-Area
mit Treffern in offenen Beobachtungen; der stehende
[`BEO-ALL/pin-bump-mirrors-ungated`](../../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
trägt Pin-Hebungen, nicht ADR-Anlass-Fragen.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung):
`make nightly-state` wird bei der Beanspruchung gelesen und der Stand in
§7 notiert.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.

### Sub-Area: `*` (Harness-Werkzeug)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — der Delta-Wortlaut ist vendiert, die
  Stelle (§3.6) ist eingrenzend, die Grenze (Senkungs-Verbot) ist
  byte-prüfbar.
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — eine Ergänzung, eine
  Gegenprobe.
- **Reconciliation-Aufwand:** Keiner — Graduation-Trigger bleibt
  [`BEO-ALL/pin-bump-mirrors-ungated`](../../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md).
