# Slice slice-188: Das Register gegen die neue Beleg-Definition

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Reaktiv: die Quelle des Kanons hat zwei Regeln
gesetzt, die unseren Bestand betreffen; sein Closure-Grund geht über die eigene
DoD nicht hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:** die zwei Quell-Wellen des Adopters `ai-harness-course`
(Bereichssegment; die Schwelle und ihre drei Ausgänge) — **noch in keinem
Release**;
[slice-183](../done/slice-183-baseline-v5150.md) (der Baseline-Bump, der sie
frühestens transportieren kann);
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(die Register-Deckung, die diesen Bestand prüft).

**Berührte Spec-Stellen:** — dieser Slice ändert **Bestand**, keine Spec-Stelle.

**Verantwortlich:** —.

**Autor:** pt9912. **Datum:** 2026-08-31.

---

## 1. Ziel

**Zwei Zähler unseres Registers sind zu hoch, und bis jetzt fehlte die Regel,
die das entscheidet.** Beide Zeilen tragen eine Begründung, beide galten als
*bewusste Abweichung* — der Adopter hat sie gemessen und uns gemeldet, und wir
haben ihm geantwortet, die Anzahl-Prüfung sei deshalb blockiert.

**Die zwei Quell-Wellen lösen das auf, und zwar gegen uns.** Die neue Regel:

- *Zwei Funde **im selben** Vorgang sind eine Gelegenheit, kein zweites
  Auftreten.*
- *Ein Vorkommen **ohne** abgeschlossenen Vorgang bewegt den Zähler nicht — es
  gehört in den Eintrag, benannt, nicht gezählt.*

Angewandt:

| Zeile | heute | nach der Regel | warum |
|---|---|---|---|
| `BEO-023` | 7 | **6** | `slice-178` steht zweimal — ein Vorgang |
| `BEO-004` | 3 | **1** | die übrigen Vorkommen hatten keinen abgeschlossenen Vorgang |

**Der Befund ist damit nicht „der Kanon war unvollständig", sondern „unsere
Zahlen waren falsch".** Sie standen nur ohne Regel unwiderlegbar da. Das ist die
unbequemere und die richtige Lesart, und sie gehört in die Closure-Notiz.

**Zwei weitere Neuerungen treffen denselben Bestand**, gemessen:

- **Belege sind nicht mehr nur Slice-Kennungen** — Welle und Review-Report sind
  ebenfalls abgeschlossene Vorgänge. Heute tragen **alle** 22 Zeilen
  ausschließlich Slice-Kennungen; die Erweiterung kann Zähler **erhöhen**, wo
  ein Vorkommen bisher unbelegt blieb.
- **Der Stand wird eine geschlossene Menge** (*verkörpert · geplant ·
  gestrichen*, darunter `offen`). Heute: **8** von 22 Zellen beginnen in einer
  dieser Formen, **14** sind freie Prosa. Betroffen von der Schwellen-Bindung
  sind die **12** Zeilen mit Zähler ≥ 3.

## 2. Vorgehen

1. **Die Timing-Frage zuerst entscheiden, nicht nebenbei.** Die Regeln liegen
   auf `main` des Adopters und in **keinem Release**; unser Pin steht auf
   `v5.12.0`, und selbst `v5.14.0` enthält sie nicht. Wer sie jetzt zitiert,
   zitiert etwas, das im vendorten Baum nicht existiert.
   **Aber:** die *Abweichung* ist schon unter dem heutigen Pin eine — `v5.12.0`
   verlangt bereits **Anzahl** (so viele Belege wie der Zähler), und 7 ≠ 6 ist
   auch dort falsch. Neu ist nur die **Richtung** der Reparatur. Zu entscheiden
   ist also: unter dem heutigen Pin korrigieren (Richtung aus dem unreleasten
   Kanon vorweggenommen, ausdrücklich als solche benannt) — oder auf ein Release
   warten und den Slice bis dahin ruhen lassen.
2. **Die zwei Zähler korrigieren**, jeweils mit der Begründung in der Zeile;
   die bisherige Abweichungs-Begründung wird **nicht gelöscht**, sondern
   abgelöst — sie war die ehrliche Notiz ihrer Zeit.
3. **Alle 22 Zeilen gegen die erweiterte Beleg-Definition durchgehen**: gibt es
   Vorkommen, die als Welle oder Review-Report belegbar sind und den Zähler
   **heben**? Das ist Urteil, kein `grep`, und es gehört einzeln begründet.
4. **Die 14 freien Stand-Zellen auf die geschlossene Menge umstellen** — und
   dabei prüfen, welche der 12 Schwellen-Zeilen einen Ausgang tragen, der die
   Form erfüllt.
5. **Erst danach ist die Anzahl-Achse baubar** (Folge-Slice): sie braucht einen
   Bestand, der die Regel erfüllt, sonst startet sie rot.
6. `make gates`; **Review**; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Nicht die Anzahl-Prüfung bauen.** Sie ist der Folge-Slice und braucht diesen
  Bestand.
- **Kein Baseline-Bump.** Der ist [slice-183](../done/slice-183-baseline-v5150.md), und
  er transportiert die neuen Regeln ohnehin nicht.
- **Keine Migration auf die Verzeichnisform.** Das ist der BEO-CR, und er wartet
  auf seine eigenen Bedingungen.
- **Keine Zähler, die aus dem Gedächtnis steigen.** Punkt 3 erhöht nur, wo ein
  abgeschlossener Vorgang benannt werden kann.

## 4. Definition of Done

- [ ] Die Timing-Entscheidung aus §2.1 ist **getroffen und begründet** — nicht
      implizit durch das Tun.
- [ ] `BEO-023` steht auf **6**, `BEO-004` auf **1**, je mit Begründung in der
      Zeile; die abgelöste Notiz bleibt lesbar.
- [ ] Alle 22 Zeilen sind gegen die erweiterte Beleg-Definition geprüft; jede
      Erhöhung nennt ihren abgeschlossenen Vorgang.
- [ ] Die Stand-Zellen tragen die geschlossene Menge; die 12 Schwellen-Zeilen
      haben einen formgültigen Ausgang.
- [ ] `make gates` grün (Exit explizit); **unabhängiger Review**.
- [ ] Closure-Notiz mit Lerneintrag; jedes Risiko aus §5 mit Ausgang; die drei
      Paarungen geprüft.

## 5. Abnahme-Punkte / Risiken

- **Zähler zu senken sieht aus wie Vertuschen.** Wer die Historie liest, sieht
  eine Zahl fallen. Die Begründung muss in der Zeile stehen, nicht nur im
  Commit — sonst ist die Korrektur von einer Schönung ununterscheidbar.
  — **Ausgang:** *(bei Closure)*
- **Punkt 3 ist die Stelle, an der aus einer Korrektur eine Erfindung wird.**
  Einen Zähler zu **heben**, weil sich ein Vorgang „auch als Beleg lesen lässt",
  wäre eine Aussage aus dem Anlass statt aus dem Bestand
  ([`BEO-011`](../observations.md), Zähler 5). — **Ausgang:** *(bei Closure)*
- **Die Regeln sind unreleast.** Ein Zitat auf sie ist heute nicht auflösbar,
  und ein Slice, der auf eine `main`-Fassung baut, kann veralten, bevor er
  schließt. — **Ausgang:** *(bei Closure)*
- **Die Stand-Umstellung berührt 14 Zellen mit gewachsener Prosa.** Beim
  Umschreiben geht leicht verloren, was die Zelle *zusätzlich* sagte
  ([`BEO-002`](../observations.md), Zähler 7). — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — heute hält
[slice-174](../done/slice-174-register-deckung.md) den Slot. **Zusätzlich
die Entscheidung aus §2.1**: entweder ein Release des Adopters, das die zwei
Wellen enthält, **oder** der ausdrückliche Beschluss, unter dem heutigen Pin zu
korrigieren.

**Rückführungen:** `in-progress` → `open`, falls Punkt 3 zeigt, dass die
erweiterte Beleg-Definition mehr Zeilen bewegt als erwartet — dann ist das ein
eigener Schnitt und keine Nebenwirkung dieser Korrektur.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `docs/plan/planning/` (das Register selbst). **Eine**
  Sub-Area, unter dem Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration pro Sub-Area); der Slice berührt weder Code noch Spec. Die
  Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:219-220 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-31, höchste Kennung
  `BEO-024`): [`BEO-011`](../observations.md) (Zähler 5) — jede **Erhöhung** in
  Punkt 3 ist eine Aussage, die aus dem Bestand kommen muss und nicht aus dem
  Anlass; [`BEO-002`](../observations.md) (Zähler 7) — die 14 umzuschreibenden
  Stand-Zellen sind gewachsene Prosa, und beim Umstellen geht leicht verloren,
  was sie zusätzlich sagten; [`BEO-012`](../observations.md) (Zähler 11) — die
  neuen Regeln liegen **unreleast** auf `main`, und sie zu zitieren, als
  stünden sie in unserem Pin, wäre genau diese Klasse. Die Regel, die diesen
  Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:225-225 -->
  > **Offene Beobachtungen sichten.**

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — die eine berührte Sub-Area fällt unter den
Default. Kein Code, keine Spec, keine Reconciliation: der Slice bringt einen
**Bestand** mit einer Regel in Übereinstimmung, die anderswo entschieden wurde.

## 9. Closure-Notiz (nach `done/`)
