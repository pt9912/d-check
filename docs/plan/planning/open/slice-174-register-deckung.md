# Slice slice-174: Die Register-Deckung bekommt ihren Wächter

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-86](../welle-86-closure-uebergang-durchsetzen.md) — dritter
von vier Slices. Ihr Closure-Trigger fordert einen Beleg, den keine einzelne
DoD liefert: dass die Vorbedingungen **am Übergang** greifen.

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(das Planning-Modul, mögliche Heimat der Bedingung);
[ADR-0028](../../adr/0028-planning-lifecycle-modul.md) (seine Bauform).
Der Kanon nennt die maschinelle Hälfte selbst — Baseline-Regelwerk
`modul-06-roadmap.md` §Das Beobachtungs-Register.

**Berührte Spec-Stellen:**
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
oder eine neue `DC-FA-*`-Anforderung — die Wahl ist Teil des Slice (§2), und mit
ihr die `.a`-Verfeinerung in der
[Spezifikation](../../../../spec/spezifikation.md).

**Verantwortlich:** —.

**Autor:** pt9912. **Datum:** 2026-08-31.

---

## 1. Ziel

**Das Beobachtungs-Register ist die einzige Deckungs-Prüfung des Steering
Loops, die der Kanon ausdrücklich als maschinell entscheidbar benennt — und
niemand hat sie gebaut.** Der Kanon trennt sauber: *Mensch urteilt, Maschine
prüft Deckung.* Geprüft werden kann, ob eine in `done/` zitierte `BEO-<NNN>`
eine Registerzeile hat; **nicht** geprüft wird die Umkehrung („jede Zeile ist
irgendwo zitiert"), weil die allermeisten unter der Schwelle stehen.

**Gemessen am eigenen Bestand, 2026-08-31 — der Wächter wäre heute grün:**

| Messung | Ergebnis |
|---|---|
| in `done/` zitierte Kennungen | **24** |
| davon ohne Registerzeile | **0** |
| Registerzeilen (Haupttabelle) | 22, dazu 2 gestrichene |

**Ein grüner Wächter ist hier kein Argument gegen ihn, sondern für ihn:** ein
erfundenes `BEO-999` in einer Closure-Notiz fiele heute niemandem auf, und die
Zeile, die es belegen soll, gäbe es nicht. Die Deckung ist gelebte Praxis ohne
Sensor — genau die Lage, die
[welle-86](../welle-86-closure-uebergang-durchsetzen.md) für den ganzen
Closure-Übergang beschreibt.

## 2. Vorgehen

1. **Die Richtung festlegen und begründen:** geprüft wird **Zitat ⇒ Zeile**,
   nicht die Umkehrung. Der Kanon schließt die Umkehrung ausdrücklich aus; das
   gehört in die Anforderung, nicht nur in den Plan.
2. **Die Scan-Menge entscheiden.** Der Kanon sagt „in `done/` zitiert".
   Gemessen wird die Kennung aber auch in [`AGENTS.md`](../../../../AGENTS.md),
   in `MR-`Dateien und in `open/`-Slices. Ob die Prüfung dort mitgilt, ist eine
   Entscheidung mit Folgen — und sie gehört nach
   [`AGENTS.md`](../../../../AGENTS.md) §3.8 in die Anforderung, weil ein Modul
   nur über das verspricht, was es scannt.
3. **Den Träger wählen — zwei Kandidaten, die Wahl ist zu messen, nicht zu
   setzen.** (a) Eine Bedingung im vorhandenen Modul `planning`: es prüft
   bereits „Doku-Behauptung ↔ Repo-Struktur", hermetisch, und das Register
   liegt in seinem Layout. (b) Ein eigenes Modul `registry`: diese Form ist im
   Register selbst vorgezeichnet ([`BEO-001`](../observations.md), gestrichen)
   — *ein Verzeichnis-Muster je Register plus Autoritäts-Datei, eine Richtung,
   ein Grund-Code* — und deckte später den ADR-Index und den
   Konventionsspeicher-Index mit. Kriterium: trägt die dritte Anwendung genug,
   um ein Modul zu rechtfertigen, oder ist sie Spekulation?
4. **Nicht über `ids`.** Der naheliegende Weg — Kennung als Linkpflicht auf
   einen Anker je Registerzeile, wie `version.md` es für Versionen tut — ist
   **gemessen verworfen**: 1450 Nennungen im Baum, davon 1042 nackt, allein 708
   in `done/` und 626 in `docs/reviews/`. Das Retrofit schriebe eingefrorene
   Lauf-Belege um. Die Messung gehört in die ADR, damit der Weg nicht zweimal
   geprüft wird.
5. **Umkehr-Probe vor dem Scharfschalten:** ein konstruiertes `BEO-999` in
   einer Closure-Notiz muss melden, und der heutige Bestand muss grün bleiben —
   beide Zahlen mit Ausgabe.
6. **Vertrag und Doku:** Anforderung (neu oder erweitert), `.a`-Verfeinerung
   samt Grund-Code, ADR für die Entscheide aus 2–4, Aktivierung im eigenen
   Profil, Handbuch.
7. `make gates` und `make fullbuild`; **Review** und **Verifikation** als
   getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Die Beleg-Form-Prüfungen (Form · Anzahl · Lage).** Sie gehören zur selben
  Kanon-Stelle, aber nicht in diesen Schnitt. **Lage** deckt bereits `links`
  (alle 21 Beleg-Zellen sind `done/`-Links). **Form** ist erfüllt, aber
  ungewächtert. **Anzahl** ist die einzige echte Lücke — und sie ginge auf dem
  Bestand rot: [`BEO-023`](../observations.md) trägt Zähler 7 bei sechs
  eindeutigen Slices, [`BEO-004`](../observations.md) Zähler 3 bei einem. Beide
  Abweichungen sind **in der Zeile begründet**, und der Kanon-Geber räumt ein,
  zwei Regeln zu schulden (ein Vorkommen außerhalb einer Slice-Closure; ein
  zweites Vorkommen derselben Klasse im selben Slice). **Der Weg dorthin ist
  benannt, nicht verbaut:** eine Bestands-Ausnahme wie
  [`MR-049`](../../../../harness/conventions.md#mr-049) und
  [`MR-056`](../../../../harness/conventions.md#mr-056) trägt so etwas — sie
  gehört in einen eigenen Slice mit eigener Messung.
- **Die Umkehrung** („jede Zeile ist zitiert"). Der Kanon schließt sie aus.
- **Das Urteil, ob zwei Beobachtungen dieselbe sind.** Mensch.
- **Das Nachrüsten von Bestands-Zitaten.** Es gibt nichts nachzurüsten: 24 von
  24 sind gedeckt.

## 4. Definition of Done

- [ ] Eine in `done/` zitierte `BEO-<NNN>` **ohne** Registerzeile erzeugt einen
      Befund mit eigenem Grund-Code; die Richtung und die **Scan-Menge** stehen
      in der Anforderung, nicht nur im Code.
- [ ] **Umkehr-Probe gemessen:** konstruiertes `BEO-999` ⇒ Befund; eigener
      Bestand ⇒ 0 Befunde. Beide Ausgaben in der Closure-Notiz.
- [ ] Träger-Entscheid (`planning`-Bedingung gegen eigenes Modul) **begründet**
      in einer ADR, samt der gemessenen Verwerfung des `ids`-Wegs.
- [ ] Anforderung, `.a`-Verfeinerung mit Grund-Code, Profil-Aktivierung und
      Handbuch sind nachgezogen.
- [ ] `make gates` und `make fullbuild` grün (Exit explizit); **unabhängiger
      Review**; **Verifikation** — beide in eigenen Kontexten.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register
      fortgeschrieben; jedes Risiko aus §5 mit Ausgang; die drei Paarungen
      geprüft.

## 5. Abnahme-Punkte / Risiken

- **Der Wächter ist von Anfang an grün.** Er verhindert eine Klasse, die noch
  nie eingetreten ist — genau die Lage, in der ein Sensor später für überflüssig
  gehalten und entfernt wird ([`BEO-013`](../observations.md), Zähler 1). Die
  Anforderung muss sagen, **wovor** er schützt, nicht nur was er prüft.
  — **Ausgang:** *(bei Closure)*
- **Die Scan-Menge ist die eigentliche Entscheidung.** Zählt nur `done/`, ist
  ein erfundenes `BEO-999` in einem `open/`-Slice oder in
  [`AGENTS.md`](../../../../AGENTS.md) weiter unsichtbar; zählt alles, meldet
  die Prüfung womöglich in Dokumenten, die gar keine Belege führen sollen.
  — **Ausgang:** *(bei Closure)*
- **Ein eigenes Modul auf Vorrat.** Kandidat (b) rechtfertigt sich mit zwei
  weiteren Registern, die heute niemand prüft. Das ist eine Aussage über
  künftigen Bedarf — [`BEO-011`](../observations.md), Zähler 5, ist genau die
  Klasse „Regel aus dem Anlass statt aus dem Bestand". — **Ausgang:** *(bei
  Closure)*
- **Die Beleg-Form bleibt liegen**, und mit ihr die zwei begründeten
  Abweichungen. Wer den Slice liest, könnte die Deckung für die ganze
  maschinelle Hälfte halten. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt keinen
Slice.

**Rückführungen:** `in-progress` → `next`, falls der Träger-Entscheid auf ein
eigenes Modul fällt — dann wächst der Slice über drei Liefer-Punkte hinaus
(Modul, Vertrag, Profil, Handbuch) und gehört neu geschnitten.
`in-progress` → `open`, falls die Scan-Mengen-Frage eine Kanon-Entscheidung
verlangt, die wir nicht allein treffen können.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/hexagon/core/rules` (die Bedingung), `spec/`
  (Anforderung und Verfeinerung) und `docs/plan/planning/` (das geprüfte
  Artefakt). Alle drei fallen unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration pro Sub-Area). Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-31, höchste Kennung
  `BEO-024`): [`BEO-013`](../observations.md) (Zähler 1) — *ein Wächter, der
  nichts mehr fängt, bleibt stehen*; hier eine Stufe früher, weil dieser
  Wächter von Anfang an nichts fängt, und deshalb als Risiko geführt;
  [`BEO-011`](../observations.md) (Zähler 5) — die Rechtfertigung eines eigenen
  Moduls mit künftigem Bedarf ist genau diese Klasse;
  [`BEO-015`](../observations.md) (Zähler 3) — *ein offener Punkt bekommt bei
  der Closure einen Ausgang, den es nicht gibt*: dieselbe Familie urteilsfreier
  Closure-Prüfungen, die dieser Slice erweitert;
  [`BEO-002`](../observations.md) (Zähler 7) — die Spiegel-Frage trifft hier
  die Anforderung, deren Scan-Menge an mehreren Stellen stehen wird. Die Regel,
  die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): `upstream-drift.yml`
  meldet **ROT** (Lauf 2026-08-31T06:31:40Z), `image-scan.yml` grün
  (2026-08-31T09:48:33Z). Gelesen: die drei roten Schritte waren
  `baseline-freshness` (bekannt, [slice-183](../open/slice-183-baseline-v5140.md)),
  `freshness-a-check` und `go-base-digest` — **die letzten beiden sind seither
  behoben**, der Lauf ist nur noch nicht wiederholt worden. Das ist die
  benannte Grenze des Targets: es liest den **jüngsten** Lauf, nicht sein
  Alter. Keiner der drei berührt diesen Slice. **Dieser Block trägt bewusst
  keine `cite`-Direktive** — sein Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-174. Betroffene IDs:
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(möglicher Träger), [ADR-0028](../../adr/0028-planning-lifecycle-modul.md).
Module: `planning` oder neu. Gates: `make test`, `make planning-check`,
`make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — alle drei berührten Sub-Areas fallen unter
den Default: Doc führt, Code folgt. Die Regel steht im Kanon, der Bestand ist
gemessen und konform (24/24), es gibt keine Reconciliation. Das
Evidenz-Risiko ist niedrig und **gemessen**, nicht geschätzt.

## 9. Closure-Notiz (nach `done/`)
