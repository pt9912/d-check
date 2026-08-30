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

Nicht teilweise blind, sondern ganz.

**Die Zahlen dieses Absatzes stimmen; der Schluss daraus stimmt nicht**
(nachgemessen in slice-180, Review und Verifikation unabhängig). Die 25 bzw. 45
Backticks im DoD-Abschnitt von `slice-061` und `slice-076` reproduzieren exakt —
und weil beide Abschnitte **je ein einziger Absatz** sind, ist die
Abschnitts-Summe hier zugleich die Absatz-Summe. **Falsch ist die Folgerung
„ungerade ⇒ exponiert".** Eine ungerade Zahl heißt, dass der **letzte** Backtick
eine Spanne öffnet, die nichts schließt — dann findet die Bereinigung gar keine
vollständige Spanne und entfernt nichts. Verschluckt wird ein Task-Item nur von
einer Spanne, die es **umschließt**: ein Backtick davor **und** einer dahinter.
Positiv-Kontrolle gefahren: an denselben Absatz von `slice-061` ein `- [ ]`
angehängt ⇒ die Regel meldet, trotz der 25.

**Die Exposition dieses Repos ist heute null**, aber aus einem anderen Grund als
„keine unbalancierten Absätze": repo-weit gibt es **sechs** Stellen in fünf
Dateien, an denen ein `- [ ]` durch eine **wohlgeformte** Spanne aus dem
bereinigten Text verschwindet — darunter zwei `done/`-Slices und dieser Plan
selbst. Alle sechs sind bewusste Prosa **über** den Marker; kein echter
unquittierter Haken ist verdeckt. Dass `spans` schweigt, belegt das **nicht**:
es meldet eine ungeschlossene Folge nur, wenn sie an Nicht-Whitespace klebt,
und eine wohlgeformte Spanne ist ohnehin kein Befund.

**Der Slice steht damit auf dem konstruierten Fall, nicht auf einem
Bestands-Fund.** Das ist kein Grund, ihn fallenzulassen — ein Wächter, der an
einem Tippfehler abschaltet, bleibt ein Defekt —, aber die Dringlichkeit ist
eine andere, als dieser Absatz zunächst behauptete.

**Ein erster Korrektur-Versuch in slice-180 lag selbst daneben** und ersetzte
die richtigen Zahlen durch „21 bzw. 4": gezählt war der Bereich zwischen `## 4.`
und `## 5.`, während beide Slices ihre DoD als **§3** führen — also §4
*Risiken*. Beide Fassungen sind
[`BEO-020`](../observations.md); die zweite entstand beim Korrigieren der
ersten.

**Der Preis aus [ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md)
bleibt richtig für Prosa — und wird zu teuer für eine
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


**Beim ersten Bau gefunden und hier aufgehoben.** Review und Verifikation haben
die Bedingung gebaut, gemessen — und einen stillen Ausfall derselben Gestalt
gefunden, gegen die sie gebaut war. Das ist der Grund für die Rückführung (§6).

- **Ein vergessener Schluss-Fence schaltet die Bedingung ab.** Die
  Zeilen-Auswahl fällt über einen Fence-Toggle; eine einzelne unbeendete
  Fence-Zeile blendet alles Folgende aus. Der **häufige** Fall wird im inneren
  Loop von `spans` gefangen (`fence-unclosed`) — das Closure-Profil, für das
  die Bedingung gedacht ist, führt `spans` aber **nicht**, und ein **naiv
  ausgeglichener** Fence entkommt auch `spans`. Damit bleibt der Wächter über
  eine Vorbedingung durch ein Fence-Zeichen abschaltbar.
- **Der Regressions-Test testete nicht, was sein Name sagt — selbst
  nachgemessen.** Das Fixture trug `` `erledigt ` `` mit **zwei** Backticks,
  also eine wohlgeformte Spanne; der **Vorgänger** findet darin den offenen
  Haken (`section-forbidden`, Exit 1). Der Test belegte „Befund auf der
  Item-Zeile", nicht „überlebt die Blindstelle". Für den tatsächlich blinden
  Fall gab es **keinen** Test. Das ist [`BEO-023`](../observations.md) im
  eigenen Wächter, einen Tag nach der Registrierung der Klasse.
- **Die Abschnittsgrenze hielt kein Test.** Die Mutation „über die ganze Datei
  zählen statt nur im Abschnitt" lief grün durch.
- **Es sind vier Bedingungen auf rohem Text, nicht drei.** Die
  **Zellenlängen**-Bedingung liest ebenfalls die rohen Zeilen und sagt das
  selbst. Die Aufzählung stand an vier Stellen falsch, und
  [ADR-0074](../../adr/0074-offene-tasks-auf-rohen-zeilen.md)s
  Re-Evaluierungs-Trigger („bei drei benannten Ausnahmen") war damit im Moment
  der Annahme schon erfüllt — er konnte nie feuern. Korrigiert als
  `## Geschichte` dort, weil der Kern einer `Accepted`-ADR nicht überschrieben
  wird.
- **Die Spezifikation widersprach sich intern.** Die `.a`-Verfeinerung
  Schritt 5 sagte weiter „zwei benannte Ausnahmen", Schritt 6 hatte keine Zeile
  für die neue Bedingung, und Schritt 7 verlangte den Befund auf der
  Abschnitts-Überschrift, während der Code auf der Item-Zeile meldete. Die
  Historie-Zeile fehlte.
- **Die Lexik war kopiert, nicht geteilt.** `openTaskItemRE` war ein wörtliches
  Präfix von `taskItemRE`, ohne Kopplungs-Test — genau die
  [`BEO-003`](../observations.md)-Form, mit der die Bedingung begründet war.
- **Bei Schwelle > 0 meldete sie alle Items**, auch die erlaubten: eine
  Verletzung, drei Befunde, und keiner davon die Reparaturstelle.
- **Die ausgewiesene Inline-Code-Grenze war überzeichnet.** Eine **einzeilige**
  Inline-Spanne meldet gar nicht — das Muster ist zeilen-verankert; nur die
  **mehrzeilige** Spanne zählt mit. Der genannte Preis war höher als der echte.
- **Ungenannte Lücken:** ein Task-Item im **Blockquote** (`> - [ ]`) entkommt,
  und `- [\t]` mit Tab in der Box ebenfalls.
- **[`BEO-016`](../observations.md) wurde in der Sichtung übersehen** — die
  Beobachtung *„Prosa verschwindet in einer absatzweiten Inline-Code-Spanne —
  und kein Gate sieht es"* ist die Klasse dieses Defekts, mit einer Prozedur,
  die der Slice de facto ausführte.

**Was davon bleibt und was neu gemacht wird:** die Bauform trägt — Befund auf
der Item-Zeile, alle vier Listen-Marker, `sections: each`, CRLF, Config-Rand,
100 % Coverage auf der Funktion. Neu zu schreiben sind die Tests für die
Blindstelle und die Abschnittsgrenze; die Zahl „drei" ist überall „vier"; die
Schwellen-Semantik ist zu entscheiden; und die Grenzen sind zu nennen, wie sie
sind, nicht wie sie klingen.

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

**Rückführung eingetreten am 2026-08-29** (`in-progress` → `open`): nicht aus
dem unten genannten Grund — die rohe Lesung erzeugt im eigenen Bestand keine
Falsch-Positive. Die Bedingung war gebaut und gemessen, als Review und
Verifikation den **stillen Fence-Pfad** fanden (§2) und dazu einen
Regressions-Test, der die Blindstelle gar nicht traf. **Entscheid des
Auftraggebers:** erst die **Fence-Lexik** selbst, dann diese Bedingung — nicht
die Blindstelle benennen und nicht `spans` ins Closure-Profil nachziehen.
Die Umsetzung ist zurückgenommen, die Messungen bleiben in §2, und
[ADR-0074](../../adr/0074-offene-tasks-auf-rohen-zeilen.md) bleibt mit drei
`## Geschichte`-Zeilen stehen — eine `Accepted`-Entscheidung verschwindet nicht
mit ihrer Umsetzung, und das Immutable-Gate hat den Löschversuch abgewiesen.

**Dritte Beanspruchung am 2026-08-30** (`open` → `in-progress`): **die
Vorbedingung aus der Rückführung ist überholt, und das ist gemessen.** Der
Entscheid lautete *„erst die Fence-Lexik selbst"*; er fiel, als `spans` im
Closure-Profil **nicht** lief. Zwei Messungen von heute:

- **Die Blindstelle zerfällt in zwei Fälle, und nur einer ist die Lexik-Frage.**
  Ein **vergessener** Schluss-Fence blendet alles Folgende aus — das tut die
  CommonMark-Lesart **genauso**, eine offene Fence läuft dort bis Dateiende. Ein
  Lexik-Wechsel behebt diesen Fall also nicht; ihn fängt `fence-unclosed`, und
  das Closure-Profil fährt `spans` seit
  [slice-180](../done/slice-180-closure-profil-spans.md)
  ([ADR-0077](../../adr/0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md)).
  Nur der **naiv ausgeglichene** Fence — eine Zeile, die der Toggle als
  Schließer liest und CommonMark nicht — ist die Lexik-Divergenz.
- **Diese Divergenz hat im Bestand null Realfälle.** Beide Automaten über alle
  620 Markdown-Dateien gegeneinander gefahren: der **einzige** Treffer ist die
  Prosa von
  [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md), die den
  Unterschied beschreibt. Dieselbe ADR hat die Frage **ausdrücklich offen
  gelassen** und ihre Bedingung benannt: die beiden offenen Punkte seien
  *„**unbelegt** — kein Realfall in den 522 Dateien — und bekommen erst eine
  Regel, wenn einer existiert"*. Es existiert keiner.

  **Zur Reproduzierbarkeit der Messung, weil sie eine Entscheidung trägt:** der
  Zwei-Automaten-Vergleich hängt daran, wie die CommonMark-Seite die
  **Infozeilen-Regel** behandelt. Meine Fassung wendet sie nur auf den Toggle
  an und findet deshalb **einen** Treffer (die ADR-Prosa selbst); der Review hat
  einen Automaten gebaut, der sie beidseits anwendet, und findet **null**.
  **Die tragende Aussage ist dieselbe** — kein Realfall, der ein Verhalten
  änderte —, aber die Trefferzahl ist eine Eigenschaft des Messverfahrens, nicht
  des Bestands, und gehört so gelesen.

**Entscheid des Auftraggebers 2026-08-30:** diesen Slice jetzt bauen; Fall 2
bleibt eine **benannte Grenze** mit dieser Messung, statt einen Lexik-Umbau
gegen die eigene Bedingung von [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) vorzuziehen.

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
  stehen (§3); [`BEO-016`](../observations.md) — Prosa verschwindet in einer
  absatzweiten Inline-Code-Spanne, und kein Gate sieht es: **das ist die
  Klasse dieses Defekts**, und sie war beim ersten Schnitt übersehen.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — bei der dritten Beanspruchung neu gelesen,
  `upstream-drift.yml` zuletzt 2026-08-30T06:08:17Z, `image-scan.yml`
  2026-08-30T09:16:25Z. **Dieser Block trägt bewusst keine
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
