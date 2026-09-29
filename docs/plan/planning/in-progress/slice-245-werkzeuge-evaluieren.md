# slice-245: `tools/harness`-Werkzeuge aus ai-harness-init evaluieren

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle as State Machine.

**Welle:** welle-91.

**Bezug:** [`welle-91`](../welle-91-adoption-ai-harness-init.md),
[`MR-073`](../../../../harness/conventions.md#mr-073),
[`MR-004`](../../../../harness/conventions.md#mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild)
(Gate-Nachweis-Mechanik, deren Pendants hier evaluiert werden).
Keine `DC-*` — falls die stille-Grün-Lücke von vcs/commits einen
Produkt-Fix ergibt, trägt der **Folge-Slice** das Anforderungs-Delta.

**Berührte Spec-Stellen:** — *(Werkzeug-Doku ist kein Spec-Stratum; ein
etwaiger Produkt-Fix ist Folge-Slice)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-29.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die `tools/harness`-Skripte der Schwester sind je evaluiert —
adoptiert (an d-check angepasst, mit `make`-Target) oder abgelehnt
(Begründung):

- **slice-mv.sh** (292 Z.): automatisiert die §3.3-Zweikommits (reiner
  Move-Commit + Verweis-Reparatur als getrennter Commit) — Eingangs- und
  Ausgangs-Verweise je Lifecycle-Wechsel.
- **selbstpruefung.sh** (239 Z.): Negativ-Selbsttest der Commit-Kennungs-
  Hooks im Wegwerf-Klon (ohne Kennung fällt, mit geht durch).
- **history-range-guard.sh** (97 Z.): Vorlauf-Wächter gegen **stilles Grün
  über leerem Prüfbereich** in shallow-Klonen für die history-lesenden
  d-check-Targets (vcs, commits) — die Lücke ist zuerst am eigenen Adapter
  zu verifizieren (Gegenprobe im shallow-Clone), dann Bauform entscheiden
  (Vorlauf-Wächter wie die Schwester oder Produkt-Fix = Folge-Slice).
- **e2e-abdeckung.sh / traeger-fetch.sh**: evaluiert — d-check führt kein
  E2E-Skript und keinen eigenen Träger; voraussichtlich n.a. mit Begründung.

**Ausdrücklich NICHT in diesem Slice:**

- **Kein Produkt-Fix an vcs/commits ohne Anforderung** — die Lücke wird
  belegt, der Fix ist Folge-Slice mit `DC-*`.
- **Keine Übernahme der .mk-Suite als Ganzes** — je Werkzeug einzeln; die
  .mk-Form (doc-gate/enforce) folgt d-checks Makefile-Struktur, nicht der
  der Schwester.

## 2. Definition of Done

- [ ] Je Werkzeug (5) eine belegte Entscheidung: adoptiert (funktioniert,
      Target vorhanden, gates grün) oder abgelehnt (Begründung).
- [ ] Bei history-range-guard: die stille-Grün-Behauptung ist am eigenen
      Adapter verifiziert (Gegenprobe im shallow-Clone, Ausgabe belegt).
- [ ] `make gates` grün.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/*.sh` / `Makefile` (nach Evaluierung) | neu/update | adoptierte Werkzeuge |
| dieser Plan | update | etwaige Anker-Korrekturen |

## 4. Trigger

**Start** (`open` → `in-progress`): direkt beansprucht — die
Welle-Eröffnung (welle-91, Auftraggeber 2026-09-29) ist die Beanspruchung;
der Nachtlauf-Stand wird bei der Beanspruchung gelesen
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Rückführungen:** `in-progress` → `next` (zu groß): ergibt die
stille-Grün-Verifikation einen Produkt-Fix mit eigenem
Anforderungs-Delta, wird er geschnitten und dieser Slice schließt mit der
Bewertung. `in-progress` → `open`: keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag (drei Paarungen
wellenlos hier geprüft).

## 6. Risiken und offene Punkte

- Adoptierte Skripte sind Host-bash — sie laufen außerhalb des
  Docker/make-only-Vertrags (§3.1 gilt der Produkt-Toolchain; die
  Gate-Skripte-Klasse ist POSIX-bash, dieselbe wie `tools/harness/`
  heute). — **Ausgang:** *(offen)*
- Die stille-Grün-Gegenprobe braucht einen shallow-Clone — Mechanik
  (Repositorie-Größe, Netz) einplanen. — **Ausgang:** *(offen)*

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

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das
Harness-Werkzeug selbst (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** keine berührte Sub-Area
mit Treffern in offenen Beobachtungen.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung):
`make nightly-state` wird bei der Beanspruchung gelesen und der Stand in
§7 notiert.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.

### Sub-Area: `*` (Harness-Werkzeug)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — die Werkzeuge der Schwester sind
  dokumentiert, die Pendants (`tools/harness/`) existieren.
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — je Werkzeug eine belegte
  Entscheidung; die stille-Grün-Verifikation ist Gegenprobe am Adapter.
- **Reconciliation-Aufwand:** Keiner — ein Produkt-Fix trägt sein
  Anforderungs-Delta als Folge-Slice.
