# slice-247: vcs-Modul meldet stilles Grün über leerer, auflösbarer Range

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle as State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst; es gibt
keine repo-weite Beobachtung darüber hinaus (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** [`slice-245`](../done/welle-91/slice-245-werkzeuge-evaluieren.md)
— dessen stille-Grün-Verifikation am eigenen Adapter (shallow-Clone, Probe
A: `HEAD..HEAD` ⇒ `0 Befund(e)`, Exit 0; Gegenprobe B: das commits-Modul
bricht auf derselben Range laut ab, Exit 2),
[ADR-0024](../../adr/0024-vcs-immutable-gate.md) (Modul `vcs`),
[`ADR-0027`](../../adr/0027-commits-traceability-modul.md) (Modul `commits`,
als laute Vergleichsstelle). Das Anforderungs-Delta (`DC-*`) entsteht in
diesem Slice — im Lastenheft, nie per ADR (Dokumentations-Regel 3).

**Berührte Spec-Stellen:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
(Lastenheft — das neue Delta entsteht daneben); Spezifikation, Abschnitt zu
den situativen Range-Modulen `vcs`/`commits`.

**Verantwortlich:** —.

**Autor:** pt9912. **Datum:** 2026-09-29.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Das vcs-Modul unterscheidet den Zustand „Range aufgelöst, aber
leer — nichts geprüft" nicht von „wirklich nichts zu melden" und meldet
erstern stille grün (Exit 0, `0 Befund(e)`). Der Slice trägt das
Anforderungs-Delta ins Lastenheft und fixt das Produkt: eine leere,
auflösbare Range ist ein ausdrücklich zu meldender Zustand mit Exit ≠ 0 —
fail-closed, dieselbe Klasse wie die von der Spezifikation verlangte
ungültige Range.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **commits-Modul**: bricht auf leerer Range bereits laut ab (slice-245,
  Probe B) — Bestand bleibt bewusst stehen; keine Änderung nötig.
- **history-range-guard**: deckt den Leerfall bereits laut (Exit 1 mit
  fetch-depth-Hinweis, slice-245, Probe D) — Bestand bleibt bewusst
  stehen; der Vorlauf-Wächter wird durch den Fix nicht überflüssig, denn
  er greift einen Vorgang früher (vor dem Containerstart) und deckt beide
  history-lesenden Targets.
- **Range-Semantik über den Leerfall hinaus** (Syntax, Expansion,
  `--staged`): ein anderer Vorgang — die Spezifikation hält diese Semantik
  schon fest; hier ändert sich ausschließlich der Leerfall.
- **CI-Workflows anderer Repos**: Konsumenten mit shallow checkout lösen
  ihren Stand selbst (fetch-depth) — die Grenze des Produkts ist laut
  Melden, nicht das Aufsetzen fremder Pipelines.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen
lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den
Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Neue `DC-*`-Anforderung im Lastenheft (nur dort; Historie-Zeile und
      Versions-Bump per [`MR-032`](../../../../harness/conventions.md#mr-032)),
      die den Leerfall als laut zu meldenden Zustand festlegt; die
      Spezifikation trägt die Grund-Code-Form des Leerfalls.
- [ ] Der Fix ist umgesetzt: leere, auflösbare Range im vcs-Modul ⇒
      Exit ≠ 0 mit benannter Meldung; der Test, der das fordert, lief ohne
      den Fix aus dem richtigen Grund rot (Bewusstes Brechen, Modul 11 —
      Gegenprobe ist derselbe shallow-Clone-Aufbau wie in slice-245).
- [ ] Regression belegt: das commits-Modul und der
      history-range-guard behalten ihr laut-Verhalten (Stände von
      slice-245, Proben B und D).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      Rollenwechsel nach Schritt 8, kein Self-Review.
- [ ] Closure-Notiz mit Lerneintrag; jedes Risiko aus §6 mit Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | neue `DC-*` + Historie-Zeile ([MR-032](../../../../harness/conventions.md#mr-032)) |
| `spec/spezifikation.md` | update | Leerfall-Grund-Code in den Abschnitt zu den situativen Range-Modulen |
| `internal/adapter/driven/git` (Adapter hinter `vcs`) | update | leere Range statt stiller Befundlosigkeit laut abbrechen |
| Testdatei zum Adapter | neu/update | Happy/Negative/Boundary — nach der neuen `DC-*`; Negative-Fall ist slice-245s shallow-Clone-Aufbau |

**Ansatz:** Der Adapter kennt nach der Range-Auflösung die Commit-Zahl
bereits — die leere Range ist dort ein eigener Zweig neben der unauflösbaren
Range (die bereits Exit 2 liefert, slice-245, Probe C), nicht ein
nachgelagerter Sonderfall.

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt, `Verantwortlich:`
gesetzt; der Nachtlauf-Stand wird bei der Beanspruchung gelesen
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): ergibt die Lastenheft-Arbeit, dass die
  Leerfall-Semantik über `vcs` hinaus (z. B. `spans`, `citations`)
  nachgewiesen leerer Prüfbereiche verallgemeinert werden muss, wird der
  Zuschnitt neu geschnitten.
- `in-progress` → `open` (blockiert): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Konsumenten, die heute eine leere Range fahren (shallow checkout ohne
  fetch-depth), bekommen künftig laut statt grün — das ist der Zweck
  (fail-closed), dennoch: die eigenen Workflows sind auf fetch-depth zu
  prüfen, bevor der Fix committet. — **Ausgang:** *(offen)*
- Die neue `DC-*` ändert zugesagte Semantik (Grund-Code der
  Range-Behandlung) — die Spiegel (Spezifikation, Benutzerhandbuch-Beispiele,
  `--print-mk`-Form) werden vor dem Editieren aufgelistet
  ([`MR-025`](../../../../harness/conventions.md#mr-025)). — **Ausgang:** *(offen)*

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
Produktmodul `vcs` samt Spec-Stratum (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** zum Planungszeitpunkt keine
offene Beobachtung mit Treffern für die Sub-Area. Der stille-Grün-Fund
(selber Klasse wie der ANLASS des history-range-guard der Schwester) geht
bei slice-245s Closure ins Beobachtungs-Register und zählt für künftige
Planungen; [slice-220](../done/wellenlos/slice-220-vcs-pfadmenge-statt-diff.md)
heilte die Objekt-Datenbank-Achse ([`CO-001`](../../carveouts/done/CO-001-vcs-range-stiller-skip.md)),
nicht den Leerfall.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung):
`make nightly-state` wird bei der Beanspruchung gelesen und der Stand in
§7 notiert.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.

### Sub-Area: `*` (Produktmodul `vcs` samt Spec-Stratum)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  und [`ADR-0024`](../../adr/0024-vcs-immutable-gate.md) tragen das
  Modul; der Leerfall ist bisher ungesetzte Schärfe, kein Widerspruch.
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — der Fund ist empirisch belegt
  (slice-245, Probe A), der Gegenstand ein enger Adapter-Zweig.
- **Reconciliation-Aufwand:** Keiner — das Anforderungs-Delta entsteht
  innerhalb dieses Slice im Lastenheft.
