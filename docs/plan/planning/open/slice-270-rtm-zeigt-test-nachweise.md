# slice-270: Die RTM zeigt die Test-Nachweise

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** Auftraggeber-Frage 2026-10-09: `make trace` zeigt nur ADR- und
Slice-Spalten, keine Tests; die Schwester-Repos pg-change-feed und
pgwire-recorder binden ihre Test-Nachweise über `trace.coverage` ein.
Auftraggeber-Entscheid: erzeugt und gewächtert, nicht von Hand gepflegt.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die RTM dieses Repos (`make trace`) zeigt je Anforderung, welche
Tests sie belegen — in zwei Spalten: **Tests** (die Go-Suite von `make test`)
und **E2E** (`make image-test` gegen das gebaute Image). Je Spalte eine
Abdeckungs-Datei mit der Tabelle `Kennung → Test → Datei:Zeile`, eingebunden
über `trace.coverage` in der `.d-check.yml`. Die Tabelle wird aus den
Testquellen **abgeleitet**: ein Test in `make test` erzeugt sie und vergleicht
sie mit der committeten Datei; weicht sie ab, ist `make test` rot.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine neue Produkt-Fähigkeit** — `trace.coverage` gibt es; der Slice ist
  Konfiguration und ein Test dieses Repos über sich selbst.
- **Das Urteil, ob ein Test seine Anforderung wirklich belegt** — die Tabelle
  zeigt die Deklaration; ob sie trägt, bleibt Review und Verifikation.
- **`make bench` und die Bestandsproben** (`blackbox-probe`) als eigene
  Spalten — kein Anlass; ein eigener Schnitt, wenn die RTM sie braucht.
- **Waisen schließen** — der Slice macht Belege sichtbar, er schreibt keine
  neuen Tests.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Zwei Abdeckungs-Dateien (Tests, E2E), aus den Testquellen abgeleitet,
      über `trace.coverage` eingebunden; `make trace` zeigt beide Spalten.
- [ ] Ein Test in `make test` hält jede Datei gegen ihre Ableitung — eine
      neue oder entfernte Deklaration ohne nachgezogene Datei ist rot
      (bewusstes Brechen belegt); `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| docs/user/abdeckung-tests.md, docs/user/abdeckung-e2e.md | neu | die zwei Labels der Coverage-Spalte |
| Ableitung und Wächter (Go-Test im Paket der Repo-Selbsttests, `make test`) | neu | Wächter gegen Drift |
| `make abdeckung` (Werkzeug, kein Gate) | neu | schreibt die Dateien nach dem Muster von `make tidy` |
| `tools/image-test.sh` | update | jede Phase deklariert ihre Kennungen an Ort und Stelle |
| `.d-check.yml` (`trace.coverage`) | update | Einbindung |
| neue ADR (Index), neuer `MR`-Eintrag (Index) | neu | Test-Nachweise entlasten eine Anforderung von der Waise |
| `harness/sensors/test.md`, `harness/README.md` (Werkzeug-Zeile) | update | dritte Zusage des Repos über sich selbst; das neue Target |

**Entscheidungen beim Beanspruchen (Messungen 2026-10-09):**

- **Deklarations-Form:** eine Kennung im Doc-Kommentar unmittelbar über
  `func Test…` — gemessen tragen 356 der 944 Testfunktionen sie schon
  (`awk` über die `//`-Zeilen direkt vor jedem `func Test`); keine Umschreibung
  nötig, die Rückführung aus §4 greift nicht. Gelesen wird mit dem Go-Parser
  (`go/parser`, `go/ast`), nicht mit einem eigenen Muster über den Quelltext
  (Workflow-Skelett Schritt 19); die Kennung selbst mit
  `trace.requirements.id-pattern`. Eine Kennung an anderer Stelle (Testname,
  Fehlermeldung, Kommentar im Rumpf) zählt nicht.
- **Ort:** `docs/user/`, wie in den Schwester-Repos.
- **Bindung:** Datei und Testname, keine Zeilennummer — eine Einfügung über
  einem Test ändert die Tabelle sonst, ohne dass sich an der Abdeckung etwas
  ändert.
- **Entlastung** (Auftraggeber-Entscheid): `trace.coverage` macht eine
  Anforderung mit Test-Nachweis waisenfrei, auch ohne Slice
  ([`DC-FA-COV-001.a`](../../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
  Schritt 5). Das lockert `make completeness-check` (`AGENTS.md` §3.6 ⇒ ADR)
  und weicht vom Baseline-Vorschlag ab, nach dem der Slice entlastet
  (Baseline-Regelwerk `grundlagen-traceability.md` §Die zweite Richtung ⇒
  `MR`). Heute ohne Wirkung: die RTM meldet 0 Waisen.

## 4. Trigger

**Start** (`next` → `in-progress`): `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): die feste Deklarations-Form verlangt, dass
  mehr als eine Handvoll Tests umgeschrieben werden — dann zuerst die Form,
  die Spalte danach.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- **Die Spalte behauptet mehr, als der Test prüft** — eine Kennung im
  Kommentar ist eine Deklaration, kein Beleg; die RTM liest sich danach wie
  ein Nachweis. — **Ausgang:** *(offen)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** geändert werden ein Test im Paket
der Repo-Selbsttests, ein Gate-Skript (`tools/image-test.sh`), das `Makefile`,
die `.d-check.yml`, zwei neue Doku-Dateien, ADR, `MR` und Harness-Doku — alle
unter dem Default `*` (`ALL`); deklariert. Kein Produkt-Code.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** gelesen am 2026-10-09.
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/state.md)
(verkörpert als Schritt 19) — der Leser der Deklarationen ist eine Erkennung:
der Go-Parser statt eines Musters, dazu eine Negativliste (Kennung im
Rumpf, im Testnamen, in einem Kommentar mit Leerzeile davor, in einem
`t.Run`-Namen) und der eigene Bestand vor dem Code;
[`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/state.md)
— die Kennung liest dasselbe `id-pattern` wie die RTM, kein zweites;
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
— die Grenze der Spalte steht in den beiden Dateien selbst.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-09 (`make nightly-state`) — `upstream-drift` grün;
`image-scan` rot im Lauf vor dem Release v0.85.0; `make image-scan` gegen
das veröffentlichte Image danach: keine behebbaren CRITICAL/HIGH auf beiden
Plattformen und beiden Registries.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
