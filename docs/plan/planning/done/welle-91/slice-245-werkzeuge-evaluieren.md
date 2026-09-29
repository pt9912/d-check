# slice-245: `tools/harness`-Werkzeuge aus ai-harness-init evaluieren

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle as State Machine.

**Welle:** welle-91.

**Bezug:** [`welle-91`](../../welle-91-adoption-ai-harness-init.md),
[`MR-073`](../../../../../harness/conventions.md#mr-073),
[`MR-004`](../../../../../harness/conventions.md#mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild)
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

- [x] Je Werkzeug (5) eine belegte Entscheidung: adoptiert (funktioniert,
      Target vorhanden, gates grün) oder abgelehnt (Begründung).
- [x] Bei history-range-guard: die stille-Grün-Behauptung ist am eigenen
      Adapter verifiziert (Gegenprobe im shallow-Clone, Ausgabe belegt).
- [x] `make gates` grün.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/*.sh` / `Makefile` (nach Evaluierung) | neu/update | adoptierte Werkzeuge |
| dieser Plan | update | etwaige Anker-Korrekturen |

## 4. Trigger

**Start** (`open` → `in-progress`): direkt beansprucht — die
Welle-Eröffnung (welle-91, Auftraggeber 2026-09-29) ist die Beanspruchung;
der Nachtlauf-Stand wird bei der Beanspruchung gelesen
([`MR-053`](../../../../../harness/conventions.md#mr-053)).

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

- **Was hat funktioniert:** Je Werkzeug eine belegte Entscheidung statt einer
  Pauschalübernahme — die fünf Entscheidungen (3 adoptiert, 2 abgelehnt)
  trägt die Tabelle des Verify-Reports
  ([../../../../../docs/reviews/2026-09-29-slice-245-werkzeuge-verify.md](../../../../../docs/reviews/2026-09-29-slice-245-werkzeuge-verify.md),
  Vollauf von `make selbstpruefung` inklusive). Jede Verhaltens-Behauptung
  (Identity-Fallback, Unterordner-Mapping, Wächter-Exit-Codes, leere Range)
  wurde im Wegwerf-Klon gemessen, bevor sie in Doku oder Folge-Plan stand —
  die shallow-Clone-Gegenprobe machte die stille-Grün-Lücke am eigenen
  Adapter sichtbar und schnitt [slice-247](../../open/slice-247-vcs-leere-range-stilles-gruen.md).
- **Was ging anders als geplant:** Die d-check-Anpassung des slice-mv war
  unvollständig — die ausgehende Verweis-Richtung zog die Pfadtiefe nicht
  mit (R1-F-1, Repro des Reviewers), und die README-Zeilen behaupteten ein
  history-range-guard-Target ohne Makefile-Regel (gate-phantom, gefunden im
  repo-weiten Handoff-Lauf). Beides ist behoben
  (09aeb4d5, 8637f08a), bevor slice-mv produktiv eingesetzt hat.
- **Steering-Loop-Eintrag:** keine Verkörperung — zwei Einträge sind gezählt,
  nicht verkörpert: die Klasse „Slice-Nummer im Kommentar" steht bei 2×
  ([BEO-ALL/kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen](../../observations/BEO-ALL/kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen/state.md)),
  die stille-Grün-Klasse wurde neu angelegt
  ([BEO-ALL/stilles-gruen-ueber-leerer-range](../../observations/BEO-ALL/stilles-gruen-ueber-leerer-range/state.md)).
- **Beobachtungs-Register (`../../observations/`):**
  `BEO-ALL/stilles-gruen-ueber-leerer-range/` neu angelegt, Beleg
  `evidence/slice-245.md`; `evidence/slice-245.md` in
  `BEO-ALL/kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen/` ergänzt —
  Zähler steht damit bei 2×.
- **Folge-Slices:** [slice-247](../../open/slice-247-vcs-leere-range-stilles-gruen.md)
  (vcs-Modul meldet stilles Grün über leerer, auflösbare Range) — ist eine
  Datei in `open/`.
- **Risiken aus §6:** Risiko 1 (Host-bash): entfallen — die Skripte bleiben
  in der POSIX-bash-Klasse des Bestands, §3.1 gilt der Produkt-Toolchain
  (R1-Negativbefund). Risiko 2 (shallow-Clone-Mechanik): entfallen — die
  Gegenprobe lief lokal über `file://` ohne Netz und mehrfach wiederholt
  (Verify-Report, eigene Messung). Nachtlauf-Stand ([`MR-053`](../../../../../harness/conventions.md#mr-053)): beide Läufe
  grün am 2026-09-29 (upstream-drift 06:29 UTC, image-scan 09:52 UTC) —
  nichts zu lesen.
- **Drei Paarungen:** (a) Anker — kein `liegt in`-Feld, nichts verkörpert,
  die Paarung trifft nicht zu; (b) Folge-Slice — slice-247 existiert in
  `open/`; (c) Register — beide genannten Verzeichnisse existieren, `evidence/`
  ist je nicht leer.

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
