# Slice slice-178: Eine Bedingung zählt offene Task-Items auf dem rohen Abschnitt

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der **Anlass** liegt in
[welle-86](../welle-86-closure-uebergang-durchsetzen.md), die **Fähigkeit**
gehört jeder `structure`-Regel.

**Bezug:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(das erweiterte Modul);
[ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md) (der
ausgewiesene Preis der absatzweisen Paarung — hier wird er zu teuer);
[ADR-0057](../../adr/0057-structure-tabellen-monotonie.md) (die Präzedenz: eine
Bedingung, die **rohe** Zeilen liest, mit Begründung);
[ADR-0069](../../adr/0069-zellenlaenge-als-strukturbedingung.md) (dieselbe
Bauform: neue Bedingung, eigener Grund-Code, Lastenheft-Bump).

**Berührte Spec-Stellen:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
und seine `.a`-Verfeinerung in der Spezifikation.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-08-29.

---

## 1. Ziel

**Eine Zusage über offene Task-Items ist auf dem bereinigten Abschnitts-Text
nicht haltbar.**

`forbid-pattern` und `max-tasks` lesen den Text, aus dem Fenced-Code entfernt
und **Inline-Code geleert** wurde. Die Paarung ist **absatzweise** — der
ausgewiesene Preis aus
[ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md). Ein
**einzelner** überzähliger Backtick irgendwo im Absatz macht damit den Rest
unsichtbar.

**Gemessen an einer echten Datei, mit dem DoD-Haken-Wächter aus
[slice-172](../open/slice-172-closure-uebergang-waechtern.md):**

| Eingabe | Ergebnis |
|---|---|
| offener Haken | 1 Befund, Exit 1 |
| derselbe Haken + ein Backtick weiter oben im Absatz | **0 Befunde, Exit 0** |

Nicht teilweise blind, sondern ganz. **Die Exposition ist real:** `slice-061`
und `slice-076` tragen heute ungerade Backtick-Zahlen in ihrem DoD-Abschnitt
(25 bzw. 45).

**Der Preis war für einen Vorlagen-Platzhalter ausgewiesen, hier zahlt ihn eine
Vorbedingung.** Der Entscheid hat die Paarung gegen eine Falsch-Positiv-Klasse
entschieden und den Preis benannt. Diese Abwägung bleibt richtig für Prosa. Für
eine **Closure-Vorbedingung** kippt sie: ein Wächter, den ein Tippfehler
abschaltet, meldet grün, wo er nichts gesehen hat.

**Die Fähigkeit ist zur Hälfte schon da, und das entscheidet den Schnitt.**
`taskItemRE` in `internal/hexagon/core/rules/structure.go` erkennt ein Task-Item
bereits nach der vollen CommonMark-Form — `-`, `*`, `+` **und** geordnete
Listen, `[ ]`, `[x]`, `[X]`, mit führendem Weißraum. Sie ist die Lexik des
Moduls; was fehlt, ist eine Bedingung, die sie auf die **rohen** Zeilen
anwendet.

**Der zweite Befund fällt damit von selbst weg:** ein `forbid-pattern` in der
Konfiguration deckt nur die Bullet-Form, die sein Autor aufgeschrieben hat —
gemessen laufen `* [ ]` und `+ [ ]` durch `- \[ \]` still hindurch. Eine
Bedingung, die die Modul-Lexik nutzt, kann diesen Fehler nicht machen. Das ist
die Klasse, die dieses Repo als [`BEO-003`](../observations.md) führt.

## 2. Vorgehen

1. **Eine neue Bedingung `max-open-tasks` (int ≥ 0)** in
   [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in):
   Obergrenze der **offenen** Task-Items im Abschnitt, gezählt über
   `taskItemRE` auf den **rohen** Abschnitts-Zeilen. Eigener Grund-Code
   (`section-tasks-open`), weil die Deduplikation zwei Verletzungen desselben
   Abschnitts sonst zusammenfallen ließe.
2. **Sie ist die dritte Bedingung, die einen anderen Text liest** — neben der
   Chronologie-Monotonie (rohe Zeilen,
   [ADR-0057](../../adr/0057-structure-tabellen-monotonie.md)) und der
   Überschriften-Bedingung (Überschriften selbst). Der Kommentar über
   `structureConditions` nennt heute **zwei** und begründet sie; er wird auf
   drei gezogen, statt die dritte stillschweigend danebenzustellen.
3. **Kein generischer „roh"-Schalter.** Ein `raw: true`, das jedes Muster an der
   Bereinigung vorbeiführte, holte die Falsch-Positiv-Klasse zurück, gegen die
   dieser Entscheid gefallen ist — ein Platzhalter **im Fence**, über den ein Slice
   schreibt, meldete dann wieder. Die Antwort ist eine **benannte** Bedingung
   mit benanntem Preis, keine Generalvollmacht.
4. **Der Befund zeigt auf die Zeile des Task-Items**, nicht auf die
   Abschnitts-Überschrift — die Reparatur ist dort, wo der Haken steht.
   Präzedenz: die Zellenlängen-Bedingung
   ([ADR-0069](../../adr/0069-zellenlaenge-als-strukturbedingung.md)) meldet
   ebenso auf ihrer Zeile. Ein Befund **je offenem Item**, nicht einer je Datei.
5. **Vor dem Scharfschalten rot messen**, und zwar gegen genau die Fälle, an
   denen der Vorgänger scheiterte: Backtick-Fall, `* [ ]`, `+ [ ]`, geordnete
   Liste, eingerückt — je eine Probe mit Erwartung und Ergebnis.
6. **ADR** für die Entscheide aus Punkt 2–4, Lastenheft-Bump,
   Spezifikations-Verfeinerung, Handbuch.
7. `make gates`; **Review** und **Verifikation** als getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **`max-tasks` bleibt, wie es ist.** Es zählt **alle** Task-Items auf dem
  bereinigten Text und teilt damit die Blindstelle. Es umzustellen wäre eine
  Verhaltensänderung an einem ausgelieferten Schlüssel; ob sie richtig ist,
  gehört gemessen und eigens entschieden. **Benannt, nicht stillschweigend
  gelassen.**
- **Keine Änderung an der Bereinigung selbst.** Die dortige Abwägung bleibt für
  Prosa richtig; dieser Slice stellt ihr eine Bedingung zur Seite, er kippt sie
  nicht.
- **Kein Anwenden auf den DoD-Haken-Wächter.** Das ist
  [slice-172](../open/slice-172-closure-uebergang-waechtern.md); dieser Slice
  liefert die Fähigkeit, nicht ihren ersten Konsumenten.
- **Keine Antwort auf die offene Design-Frage** aus slice-172 §2 (die
  wellen-eingelöste offene Box). Sie ist eine Planungs-Frage, keine
  Produkt-Frage, und wird dort entschieden.

## 4. Definition of Done

- [ ] `max-open-tasks` ist im Schema, in
      [`spec/lastenheft.md`](../../../../spec/lastenheft.md) (Bump + Historie)
      und in [`spec/spezifikation.md`](../../../../spec/spezifikation.md)
      geführt, mit eigenem Grund-Code; **explizit** < 0 ⇒ Exit 2, mit Test.
- [ ] **Die Blindstelle ist gemessen geschlossen:** derselbe Backtick-Fall, an
      dem der Vorgänger auf 0 Befunde fiel, meldet jetzt — Ausgabe vorher und
      nachher in der Commit-Botschaft.
- [ ] **Alle Bullet-Formen gemessen:** `-`, `*`, `+` und die geordnete Liste,
      je offen und je gehakt; eingerückt und mit Tab-Trenner. Erwartung und
      Ergebnis je Fall.
- [ ] **Fence-Treue bleibt:** ein Task-Item **innerhalb** eines Fenced-Blocks
      zählt **nicht** — sonst meldete ein Slice, der über Task-Items schreibt,
      seine eigene Dokumentation. Gemessen, nicht behauptet.
- [ ] Ein Befund **je offenem Item**, auf **seiner** Zeile; zwei offene Items
      in einer Datei ⇒ zwei Befunde, nicht einer.
- [ ] **Umkehr-Probe** ([`BEO-023`](../observations.md)): je Zusage eine
      Mutation, die genau einen Test rot macht — die Probe kostet einen Lauf und
      ist der einzige Beleg, dass der Wächter beißt.
- [ ] Eine ADR begründet die drei Entscheide aus §2 und ist im
      [ADR-Index](../../adr/README.md) eingetragen.
- [ ] Das [Benutzerhandbuch](../../../user/benutzerhandbuch.md) führt die
      Bedingung dort, wo es die übrigen `structure`-Schlüssel führt.
- [ ] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Zwei Bedingungen über dieselbe Frage, verschiedene Antworten.**
  `max-tasks` (bereinigt, alle Items) und `max-open-tasks` (roh, offene Items)
  stehen nebeneinander; wer den falschen greift, bekommt stillschweigend die
  schwächere Zusage. — **Ausgang:** *(bei Closure)*
- **Die rohe Lesung holt Falsch-Positive zurück, die die Bereinigung
  fernhält** — ein `- [ ]` in einem Inline-Code-Beispiel zählt jetzt mit. Die
  Fence-Ausnahme deckt den häufigen Fall, die Inline-Form nicht. —
  **Ausgang:** *(bei Closure)*
- **Die Fähigkeit entsteht für einen einzigen Konsumenten**
  ([`BEO-011`](../observations.md)): slice-172. Ob weitere Regeln offene
  Task-Items zählen wollen, ist nicht gemessen. — **Ausgang:** *(bei Closure)*
- **Der Grund-Code-Raum wächst um einen weiteren `section-*`-Code.** Jede neue
  Bedingung bringt einen; die Menge ist inzwischen zweistellig, und ein Leser
  muss sie unterscheiden können. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): unmittelbar — der Slice entsteht als
Vorbedingung von
[slice-172](../open/slice-172-closure-uebergang-waechtern.md), der dafür zum
zweiten Mal nach `open/` zurückgeführt wurde.

**Rückführungen:** `in-progress` → `open`, falls sich zeigt, dass die rohe
Lesung im eigenen Bestand mehr Falsch-Positive erzeugt als die Blindstelle
Falsch-Negative — dann ist die Abwägung eine andere, und [ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md)s Preis bleibt
der bessere.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/hexagon/core/` (Kern: Modell und Regel) und
  `spec/` (Anforderung und Verfeinerung). Beide fallen unter den Default `*` =
  **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration); eine eigene Deklaration führt nur `tools/harness/`.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-29, höchste Kennung
  `BEO-023`): [`BEO-023`](../observations.md) — ein Wächter, der nie fangen
  konnte: der Anlass dieses Slice ist genau das, und die Umkehr-Probe steht
  deshalb als DoD-Punkt; [`BEO-003`](../observations.md) — eine geteilte Lexik
  driftet, weil jeder Konsument sie selbst vorbereitet: das
  Konfigurations-`forbid-pattern` hat die Bullet-Formen selbst nachgebaut und
  zwei verfehlt, die Modul-Lexik kennt alle; [`BEO-011`](../observations.md) —
  Regel aus dem Anlass: die Fähigkeit entsteht für **einen** Konsumenten, und
  das steht als Risiko in §5; [`BEO-013`](../observations.md) — ein Wächter, der
  nichts mehr fängt: `max-tasks` teilt die Blindstelle und bleibt bewusst
  stehen (§3). Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-29T07:34:35Z,
  `image-scan.yml` 2026-08-29T10:07:43Z. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption, kein
  Baseline-Abschnitt ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-178. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in).
Module: `structure`. Gates: `make gates`, `make test`, `make doc-check`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Eine neue optionale Bedingung über eine
vorhandene Modul-Lexik; kein Fremdsystem, keine Reconciliation, kein Bestand,
der umgestellt werden müsste.

## 9. Closure-Notiz (nach `done/`)
