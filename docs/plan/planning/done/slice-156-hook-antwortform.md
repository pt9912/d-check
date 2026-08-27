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

- [x] Die Antwortform ist die aktuelle; die Wirkung ist **gemessen**, nicht aus
      der Dokumentation geschlossen.
- [x] Der Fall „Fail-closed ohne wirksame Antwortform" ist entweder durch einen
      zweiten Kanal gedeckt oder als benannte Grenze eingetragen.
- [x] Die Proben prüfen die wirksame Form; ein Verdikt `crash` bleibt erhalten.
- [x] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Eine Antwortform, die im Testlauf wirkt und im Ernstfall nicht, ist
  schlimmer als die alte** — sie sieht richtig aus. Die Messung entscheidet,
  nicht die Dokumentation. — **Ausgang: eingetreten, und zwar an der
  Proben-Seite statt an der Antwortform.** Der Review hat gezeigt, dass die
  Probenmenge stdin ausschließlich als Pipe bespielte — also genau die eine
  Form, unter der auch der defekte Leseweg funktionierte. Sie wäre an dem Tag
  grün geblieben, an dem der Wächter jeden Befehl durchließ. Der Ausgang ist
  gebaut, nicht bloß benannt: fünf Leseweg-Proben (Pipe, Datei, Here-String,
  `/dev/null`, geschlossen), und die Gegenprobe mit wiederhergestellter
  Regression meldet `FAIL crash stdin GESCHLOSSEN`. Die Klasse dahinter führt
  [`BEO-018`](../observations.md).
- **Zwei Kanäle können sich widersprechen.** Exit-Code und JSON-Antwort
  zugleich zu senden, kann eine Form entwerten. Gehört gemessen, bevor es
  gebaut wird. — **Ausgang: entfallen, gemessen.** Beide zugleich sind
  konfliktfrei: der `permissionDecisionReason` kam beim Aufrufer wörtlich an,
  `stderr` blieb leer — das stdout-JSON wurde also trotz `exit 2` gelesen. Die
  Messung lief **vor** dem Einbau, wie §2 Schritt 3 es verlangt.

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

**Der Slice hat seine eigene Prämisse widerlegt.** Er trat an, weil der Wächter
nicht mehr durchgesetzt wurde, und benannte die veraltete Antwortform als
wahrscheinlichsten Grund — ausdrücklich als Vermutung, mit der Auflage, vor dem
Bauen zu messen. Die Messung hat sie gekippt: **beide** Antwortformen wirken,
die veraltete über die Abwärtskompatibilitäts-Abbildung, die aktuelle direkt.

**Die Ursache war eine Regression aus [slice-155](slice-155-guard-ohne-node.md).**
Um dort einen Review-Befund zu schließen (`cat` als zweite Host-Abhängigkeit)
wurde die Eingabe auf `$(</dev/stdin)` umgestellt. Das funktioniert über eine
Pipe und scheitert, sobald `stdin` keine wieder-öffenbare Datei ist — unter
`set -e` endet der Wächter dann **ohne Ausgabe**. Belegt per Ablauf-Spur im Hook
selbst: die Marke vor der Lesezeile erscheint, die dahinter nicht. Ein Fix für
*fail-closed* hat den Wächter in jedem Fall **fail-open** gemacht, und das über
einen ganzen Arbeitstag, während an derselben Datei seine Grenzen geschärft
wurden. Bemerkt wurde es durch Zufall, an einem Befehl, der versehentlich einen
blockierten Interpreter trug.

**Der Formwechsel bleibt trotzdem — als Hygiene, nicht als Fix.** So steht er
auch in [`MR-044`](../../../../harness/conventions.md#mr-044). Dazu der zweite
Kanal: jeder Block trägt jetzt zusätzlich einen Nicht-Null-Exit, damit keine
Ablehnung mehr allein an einer Auslegung hängt, die der Wächter nicht
kontrolliert. Was der Eintrag ausdrücklich **nicht** behauptet: dass zwei Kanäle
die Klasse schließen. Sie decken den Ausfall der Antwortform, nicht den Ausfall
des Wächters *vor* der Antwort.

**Die teuerste Lehre liegt nicht im Wächter, sondern in den Proben.** Sie hatten
für ihn gebürgt und hätten den Ausfall nicht gefangen: sie bespielten `stdin`
nur als Pipe und meldeten ihr Urteil in Prosa, ohne es in den Exit zu legen.
Beides ist geschlossen, und beides ist per Gegenprobe belegt statt behauptet.
Die Klasse führt [`BEO-018`](../observations.md), Zähler 1.

**Zwei eigene Aussagen waren zu weit.** Die Botschaft von `04f7f2e` nennt
„452 → 424 Zeilen"; gemessen waren es **468 → 424** — die Ausgangszahl stammte
aus der Erinnerung an einen früheren Stand, nicht aus einer Messung
([`BEO-009`](../observations.md) Richtung a). Und `harness/README.md` schrieb
die Umgehungs-Tabelle [`MR-044`](../../../../harness/conventions.md#mr-044) zu,
die keine trägt; sie steht in
[`MR-042`](../../../../harness/conventions.md#mr-042). Beides korrigiert, die
Botschaft im Verlauf bleibt stehen.

**Ein Widerspruch ist aufgelöst statt befolgt.** Die Begründung, mit der der
Slice-Verweis aus `AGENTS.md` fiel, widerspricht der Aufzählung in
[`MR-013`](../../../../harness/conventions.md#mr-013), die `AGENTS.md` §4 und
`harness/README.md` §Sensors als Träger von Slice-Pfad-Verweisen führt. Beide
tragen null. Aufgelöst als
[`MR-045`](../../../../harness/conventions.md#mr-045) — für diese zwei Ziele ist
die Klausel leer, für alle übrigen gilt sie unverändert.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 528 Dateien, 0 Befunde),
`make guard-probe` (38 Proben, 0 Fehlschläge, **Exit 0** — das Target trägt sein
Urteil jetzt selbst), `make planning-check` (Exit 0). Ein unabhängiger Review
ist gelaufen; seine vier MEDIUM und fünf LOW sind in `5a46acc` eingearbeitet.
