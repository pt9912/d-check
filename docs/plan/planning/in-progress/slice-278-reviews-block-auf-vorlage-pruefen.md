# slice-278: Der `reviews`-Block dieses Repos gegen die Vorlage von v6.18.1 messen — sieht der Lauf Zusagen?

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — reaktive Arbeit, die Closure-Bedingung ist die DoD
selbst.

**Bezug:** [`DC-FA-RVW-001`](../../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
(das Modul `reviews`), [ADR-0081](../../adr/0081-reviews-modul.md),
[`MR-080`](../../../../harness/conventions.md#mr-080) (der Anlass: die Vorlage
`templates/.d-check.yml` der Baseline v6.18.1 trägt Schlüssel, die der `reviews`-Block
dieses Repos nicht führt). Anlass: Befund aus der Messung des Deltas in slice-277 — der
`reviews`-Block der [`.d-check.yml`](../../../../.d-check.yml) trägt `done-dir`,
`reviews-dir`, `recursive` und `skip-pattern`, nicht aber `require-promises`, `match`
und `skip-allows-empty`; `make review-coverage` meldet 0 Befunde, **ob der Lauf eine
nichtleere Menge von Zusagen sieht, ist nicht gemessen**.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Messen, wie viele Zusagen das Modul `reviews` in diesem Repo erkennt, und —
ergibt die Messung einen Leerlauf oder eine Fehl-Zuordnung — den `reviews`-Block so
setzen, dass er eine nichtleere Menge sieht und Fehlen eines Reports meldet.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderungen am Modul `reviews`.** Das Modul ist durch zwei Releases und die CRs
  schon nachgezogen; hier geht es um die **Konfiguration dieses Repos**, nicht um das
  Produkt (anderer Vorgang an anderem Gegenstand).
- **Alle übrigen Deltas der Baseline v6.18.1** — Vorlagen-Pfade und Review-Report-
  Dateinamen sind in [`MR-080`](../../../../harness/conventions.md#mr-080) genannt und
  warten auf eigene Entscheidungen.
- **Nachträgliche Reports für geschlossene Slices.** Meldet die Messung Slices ohne
  Report, ist das ein Befund über den Bestand: er wird benannt und je nach Ergebnis
  Carveout oder Folge-Slice, nicht hier nachgeholt.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] **Messung:** Die Zahl der erkannten Zusagen und die Zuordnung jedes
      geschlossenen Slice zu seinem Report stehen im Slice — am Lauf mit der
      heutigen Konfiguration und mit `require-promises: true`, an Kopien der
      Konfiguration, nicht am Arbeitsbaum.
- [ ] **Konfiguration:** Der `reviews`-Block trägt die Schlüssel, die die Messung
      als nötig ausweist (`require-promises`, `match`, `skip-allows-empty`), oder der
      Slice begründet je Schlüssel, warum er fehlen darf; `make review-coverage` ist
      grün und **nicht über leerer Menge**.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8, kein
      Self-Review; Verifikation.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — oder
      „keine Beobachtung angefallen" in §7 notiert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.d-check.yml` | update | `reviews`-Block nach Messung |
| `docs/reviews/…` | — | liest nur; Messung gegen die vorhandenen Reports |
| `harness/README.md` | update (falls gemessen nötig) | Sensor-Zeile `review-coverage` |

**Messung (2026-10-21, an einer Wegwerf-Kopie des Repos unter dem Home — Docker
mountet `/tmp` nicht —, nach der Messung entfernt; jeder Lauf mit dem Image
`d-check:latest`, Modul `reviews` allein):**

| Lauf | Kopie | Ergebnis |
|---|---|---|
| A | heutige Konfiguration, alle 119 Reports | Exit 0, 0 Befunde |
| B | heutige Konfiguration, **keine** Reports | Exit 1, **42** `review-missing` — das Modul sieht 42 Zusagen (276 `done/`-Dateien, davon 254 Stubs) |
| C | Wortlaut der Review-Zeile geändert („Durchsicht erfolgt"), keine Reports | **Exit 0, 0 Befunde — grün über leerer Menge** |
| D | wie C, mit `require-promises: true` | Exit 1, ein Befund „keine Review-Zusage unter 47 Kandidat(en)" |
| E | gesunder Bestand, mit `require-promises: true` | Exit 0, 0 Befunde (keine Regression) |
| F | gesunder Bestand, mit `match: name` | Exit 1, **42** Befunde — Reports heißen nach `slice-<NNN>`, nicht nach dem vollen Slug |
| G | alle lebenden Slices als Stub markiert | Exit 1 „leere Pruefmenge … fail-closed"; mit `skip-allows-empty: true` Exit 0 |

Ergebnis je Schlüssel: **`require-promises: true` wird gesetzt** (D belegt das stille
Grün von C, E belegt: heute folgenlos). **`match: name` bleibt weg** — F zeigt 42
Fehlbefunde; die Vorlage empfiehlt es für Repos mit Slug-Kennungen, d-check führt
`slice-<NNN>`. **`skip-allows-empty` bleibt weg** — G: nur der Fall „alles archiviert"
betrifft ihn, heute liegen 22 lebende Slices; wer den Schlüssel setzt, nimmt dem Modul
die Fail-closed-Antwort, wenn ein falsches `skip-pattern` alles ausnimmt
([ADR-0106](../../adr/0106-skip-allows-empty-erklaert-den-ruhezustand.md)); die Entscheidung
fällt an dem Tag, an dem alles archiviert ist.

**Plan-Änderung vor dem Code (2026-10-21, Schritt 15):** Zusätzlich zur Tabelle trägt die
Änderung `.d-check.closure.yml`: der `reviews`-Block dort ist an den von `.d-check.yml`
**gekoppelt** (Kommentar im Block: weichen die Schlüssel ab, antworten
`make review-coverage` und der Übergangs-Wächter `make verify-closure-notes` verschieden,
[ADR-0082](../../adr/0082-uebergangswaechter-reviews-observations.md)/[ADR-0105](../../adr/0105-reviews-liest-done-unterverzeichnisse.md)); `require-promises` muss in beiden stehen. Dazu der Kopplungs-Kommentar
in `.d-check.yml` („dieselben zwei Schlüssel" → drei).

**Spiegel vor dem Editieren auflisten**
([`MR-025`](../../../../harness/conventions.md#mr-025), Schritt 17): `grep` nach
`require-promises`, `skip-allows-empty`, `reviews.match` und `review-coverage` über
`harness/`, `docs/user/`, `docs/maintainer/`, `.claude/`, `tools/`, `Makefile`,
`spec/` und `.d-check*.yml`; Treffer vor dem ersten Edit im Plan listen.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen Slice;
Priorisierung gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): Die Messung deckt eine Fehl-Zuordnung auf, die
  eine Änderung am Modul braucht — dann trennt sich der Produkt-Teil als eigener Slice
  ab, dieser bleibt Konfiguration.
- `in-progress` → `open` (blockiert): Der Lauf mit `require-promises: true` meldet
  Befunde, die Reports für geschlossene Slices verlangen — dann zuerst entscheiden,
  Carveout oder Nachholen (Abgrenzung).

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Messung ergibt einen Leerlauf, der seit Monaten bestand** — der Lauf war
  grün über einer leeren Menge. — **Ausgang:** offen bis Closure.
- **`require-promises: true` macht den Lauf rot über Altbestand.** — **Ausgang:**
  offen bis Closure; Abgrenzung (Carveout oder Folge-Slice).

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

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Repo-Konfiguration
(`.d-check.yml`) unter dem Default `*` (`ALL`, GF); keine eigene Sub-Area.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; Treffer:

- [`skip-allows-empty-macht-umzug-still`](../observations/BEO-ALL/skip-allows-empty-macht-umzug-still/observation.md)
  (1×) und
  [`name-zuordnung-deckt-praefix-slices`](../observations/BEO-ALL/name-zuordnung-deckt-praefix-slices/observation.md)
  (1×): beide betreffen genau die Zuordnung und den Ruhezustand des Moduls.
- [`stilles-gruen-ueber-leerer-range`](../observations/BEO-ALL/stilles-gruen-ueber-leerer-range/observation.md)
  (1×): die Klasse „grün über leerer Menge", hier für Zusagen statt Commits.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`image-scan` grün; `upstream-drift` rot, **planmäßig gelesen** (der Baseline-Pin
steht seit slice-277 auf `v6.18.1`; veraltet sind Go, golangci-lint, semgrep und
a-check, eigene Vorgänge).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
