# Slice slice-156: Der Wächter antwortet in der Form, die das Werkzeug heute führt

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`MR-042`](../../../../harness/conventions.md#mr-042) (der Wächter in
seiner heutigen Fassung); [`MR-005`](../../../../harness/conventions.md#mr-005)
(der Nachweis-/Guard-Mechanik-Eintrag);
[`AGENTS.md`](../../../../AGENTS.md) §3.1.

**Berührte Spec-Stellen:** — (Durchsetzungsschicht; das Produkt bleibt
unberührt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

Der Wächter antwortet mit `{"decision": "block", "reason": …}` auf oberster
Ebene. Für `PreToolUse` ist das die **veraltete** Form; die aktuelle ist
`hookSpecificOutput.permissionDecision`. Die alte Form wirkt heute nur über eine
Abwärtskompatibilitäts-Abbildung des Werkzeugs — sie ist gemessen wirksam, aber
sie ist geliehene Zeit.

**Warum das jetzt schwerer wiegt als vorher.** Bis zur bash-Übersetzung gab es
einen zweiten, formunabhängigen Block-Kanal: der Fail-closed-Fall „Interpreter
fehlt" endete mit `exit 2`, und ein Nicht-Null-Exit blockiert unabhängig von
jedem JSON. Diesen Kanal gibt es nicht mehr — **jeder** Block-Pfad läuft jetzt
über die veraltete Antwortform. Fällt die Abbildung weg, fällt der Wächter
still aus: er läuft, er urteilt, und niemand hört ihn.

**Der Fall ist eingetreten, nicht befürchtet — gemessen am 2026-08-27.** Ein
`pip --version` lief über das Bash-Werkzeug durch. Dieselbe Eingabe, direkt
gegen den Wächter gefahren, wird geblockt; die Ausgabe ist **byte-identisch** zu
der Fassung, die früher am selben Tag noch durchgesetzt wurde — beim Autor wie
beim Review-Sub-Agenten, beide mit genau diesem `reason`-Text. Der Wächter
urteilt also richtig, und niemand handelt danach.

**Was daraus NICHT folgt:** *warum*. Die Antwortform ist der wahrscheinlichste
Grund, weil sie die einzige deklariert veraltete Stelle im Pfad ist — bewiesen
ist das nicht. Ein Wechsel der Form, der die Ursache nicht trifft, sähe
hinterher genauso aus wie einer, der sie trifft. Deshalb steht die Messung in
Schritt 1 vor dem Bauen, und deshalb ist der zweite Kanal aus Schritt 3 keine
Kür.

**Getrennt gehalten wird die Prüfung von der Antwortform.** Dieser Slice fasst
die Prüflogik nicht an.

## 2. Vorgehen

1. **Erst messen, was heute wirkt** — beide Formen gegen das Werkzeug fahren und
   die Wirkung lesen, nicht die Dokumentation zitieren. Ein Wächter, dessen
   Antwort niemand liest, ist die teuerste Form des Grünseins.
2. Die Antwort auf `hookSpecificOutput.permissionDecision` umstellen, mit
   `permissionDecisionReason` als Träger des heutigen `reason`-Textes.
3. **Den formunabhängigen Kanal zurückholen** oder ausdrücklich verwerfen: ein
   Nicht-Null-Exit für die Fail-closed-Fälle wäre der zweite Riegel, der heute
   fehlt. Ob das Werkzeug beides zugleich sinnvoll verarbeitet, gehört gemessen.
4. `make guard-probe` erweitern: die Proben lesen heute auf das Wort `decision`.
   Sie müssen auf die **wirksame** Form prüfen, nicht auf eine Zeichenkette, die
   zufällig in beiden Formen vorkommt.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Änderung an der Prüfung.** Segmentierung, Denylists, Rekursion,
  Extraktor bleiben, wie sie sind.
- **Keine Lockerung.** Was heute blockiert, blockiert danach.

## 4. Definition of Done

- [ ] Die Antwortform ist die aktuelle; die Wirkung ist **gemessen**, nicht aus
      der Dokumentation geschlossen.
- [ ] Der Fall „Fail-closed ohne wirksame Antwortform" ist entweder durch einen
      zweiten Kanal gedeckt oder als benannte Grenze eingetragen.
- [ ] Die Proben prüfen die wirksame Form; ein Verdikt `crash` bleibt erhalten.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Eine Antwortform, die im Testlauf wirkt und im Ernstfall nicht, ist
  schlimmer als die alte** — sie sieht richtig aus. Die Messung entscheidet,
  nicht die Dokumentation. — **Ausgang:** *(bei Closure)*
- **Zwei Kanäle können sich widersprechen.** Exit-Code und JSON-Antwort
  zugleich zu senden, kann eine Form entwerten. Gehört gemessen, bevor es
  gebaut wird. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls sich zeigt, dass die neue Form
in dieser Werkzeug-Version nicht zuverlässig wirkt — dann bleibt die alte
stehen, und die Messung ist das Ergebnis.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Durchsetzungsschicht (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-007`](../observations.md) für jeden Exit hinter einer Pipe;
  [`BEO-009`](../observations.md) für jede Aussage darüber, dass die neue Form
  „wirkt" — die gemessene Menge trägt den Schluss, nicht der Eindruck.

Slice-ID: slice-156. Betroffene IDs: — (kein `DC-`-Bezug;
Durchsetzungsschicht). Module: — . Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Formwechsel an einem vorhandenen Artefakt
der Durchsetzungsschicht.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
