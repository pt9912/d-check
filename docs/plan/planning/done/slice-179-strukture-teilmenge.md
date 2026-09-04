# Slice slice-179: Eine `structure`-Regel erklärt ihre Grundmenge

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der Anlass ist ein **eingehender CR**, keine Welle.

**Bezug:** [der eingehende CR](../../cr/2026-08-30-cr-a-check-structure-teilmenge.md)
(Antrag und Beleg des Absenders);
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(das erweiterte Modul);
[ADR-0070](../../adr/0070-tabellen-klammer-und-spaltenliste.md) (die Präzedenz:
eine Option an einer bestehenden Bedingung statt eines neuen Moduls);
[ADR-0069](../../adr/0069-zellenlaenge-als-strukturbedingung.md) (dieselbe
Bauform).

**Berührte Spec-Stellen:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
und seine `.a`-Verfeinerung in der Spezifikation.

**Verantwortlich:** pt9912 (Implementer-Rolle, beansprucht 2026-08-30).

**Autor:** pt9912. **Datum:** 2026-08-30.

---

## 1. Ziel

**Eine `structure`-Regel prüft heute die vorgefundene Menge; sie kann ihre
Grundmenge nicht erklären.** Ein
[eingehender CR](../../cr/2026-08-30-cr-a-check-structure-teilmenge.md) beantragt
zwei Optionen, die beide nur **verkleinern** — `tasks-ignore-pattern` an
`max-tasks` und `exempt-sections` an der Abschnitts-Auswahl.

von `max-tasks: 3` über die eigenen Slice-Pläne liefert **166 Befunde**: **80**
`section-oversized` über **89** DoD-Abschnitte und 86 Dateien ohne den
Abschnitt. Die jüngsten Slices tragen 5 bis 9 DoD-Items, von denen der
Gate-/Review-/Verifikations-Punkt in **jedem** steht. Nicht der Bestand ist
falsch — der Zähler misst das Falsche.

**Der zweite Bedarf ist hier ebenso belegt.** Grandfathering läuft in diesem
Repo über **Datei**-Globs ([`MR-049`](../../../../harness/conventions.md#mr-049)),
weil es nichts Feineres gibt. Eine Regel mit `sections: each` über **eine**
Datei — etwa das Lastenheft — hätte keinen Hebel.

## Der Antrag hat zwei gemessene Defekte, und beide sind billig zu beheben

Die Messung unten ist **nachgezogen**, nachdem die Semantik feststand: der
erste Anlauf zählte gegen die **rohen** Zeilen und kam auf 24 gegen 23
Treffer. Das Werkzeug zählt aber gegen den **bereinigten** Text — und damit
verschiebt sich das Bild so weit, dass die alte Zahl nicht bloß ungenau war,
sondern die falsche Frage beantwortete. Gefahren mit dem gebauten Stand über
die **89** DoD-Abschnitte dieses Repos (**444** Task-Items, `max-tasks: 0`,
damit jeder Abschnitt seine Zahl meldet):

| Musterform | Items ignoriert | davon falsch |
|---|---|---|
| CR-Ausdruck, frei | 13 | **2** |
| **derselbe** Ausdruck, verankert (`^…`) | **1** | **1** |
| verankert **und** auf Text, der die Bereinigung überlebt | 26 | 0 |

**Erster Defekt: das freie Muster nimmt echte Zusagen mit.** Zwei der 13 sind
**Liefer**-Punkte, die die Closure-Notiz nur als ihren **Ort** nennen — *„§2/§3/§4/§6
tragen `SPEC-NNN` fortlaufend; Zählung in der Closure-Notiz"* und
*„[ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)-§Kern-Messung in der Closure-Notiz"*. Ein freies Substring-Muster
entfernt sie aus der Zählung, ohne dass jemand es sieht — die Gestalt aus
[`BEO-023`](../observations.md): ein Filter, der still das Falsche entfernt,
sieht aus wie einer, der richtig filtert. **Verankern allein rettet den
Ausdruck aber nicht:** derselbe Ausdruck mit `^` behält genau einen Treffer, und
auch der ist ein Liefer-Punkt. Tragfähig wird erst die Kombination —
verankert **und** auf Text gerichtet, den die Bereinigung übrig lässt (dritte
Zeile). Das ist die Korrektur an der ersten Fassung dieser Messung, die „frei
13" neben „verankert 26" stellte, als wäre es derselbe Ausdruck; ein `^` kann
die Treffermenge nur verkleinern. Gefunden von Review und Verifikation
unabhängig voneinander — die Klasse ist
[`BEO-020`](../observations.md).

**Zweiter Defekt, und er ist der schärfere: die erste Alternative des CR ist
hier komplett wirkungslos.** `make gates` trifft **null** von 444 Items —
gemessen, nicht geschätzt. Der Grund ist die Zusage, die der CR selbst
verlangt: die Zählung liest den **bereinigten** Text, und in **86 der 89**
Abschnitten steht die Wendung in Inline-Code (`` `make gates` ``, null
Fundstellen ohne Backticks). Was fence- und inline-code-treu gezählt wird, ist
auch fence- und inline-code-treu **ignorierbar** — wer ein Muster auf einen
Ausdruck in Backticks richtet, richtet es auf Leerzeichen. Das gehört in die
Antwort an den Absender und ins Handbuch, nicht in eine Fußnote.

## 2. Vorgehen

1. **Beide Optionen aufnehmen, wie beantragt** — als Optionen an bestehenden
   Regeln, ohne neue Grund-Codes, Default byte-identisch. Die Verortung folgt
   [ADR-0070](../../adr/0070-tabellen-klammer-und-spaltenliste.md): eine Option
   an der Bedingung, kein neues Modul.
2. **Die Überdeckung sichtbar machen.** Die Meldung von `section-oversized`
   nennt die ignorierten mit: *„Abschnitt trägt N Task-Items (M ignoriert),
   erlaubt sind K"*. Ein zu breites Muster fällt damit auf, sobald die Regel
   meldet. **Die Grenze bleibt und gehört benannt:** greift das Muster so
   breit, dass die Regel gar nicht mehr meldet, sieht es niemand.
3. **Das Beispiel im Handbuch wird verankert gezeigt**, mit der Messung aus §1
   als Begründung — nicht als Stilfrage, sondern als gemessener Unterschied
   zwischen 23 und 24 Treffern.
4. **Der Name `exempt-sections` ist ein Entscheid, kein Erbe.** `vcs` führt
   bereits `exclude-sections` — **literale** Namen, andere Semantik, anderes
   Modul. Zwei ähnliche Namen für zwei Dinge brauchen eine Zeile, die sagt
   warum; die Alternative ist ein dritter Name.
5. **Die Frage, die der CR nicht stellt, wird beantwortet: leert
   `exempt-sections` die Kandidatenmenge, meldet die Regel `section-missing`** —
   wie bei `exempt-paths`, wo das Lastenheft es ausdrücklich zusagt (*„auch dann,
   wenn erst `exempt-paths` die Menge geleert hat"*). Ohne diese Antwort
   schaltete ein zu breites `exempt-sections` die Regel **still** ab.
6. **Fence-Treue für beide Muster**, wie vom Absender verlangt — und gemessen,
   nicht zugesagt.
7. **Umkehr-Proben** ([`BEO-023`](../observations.md)): je Zusage eine Mutation,
   die genau einen Test rot macht. Die Fixtures des Absenders sind der
   Ausgangspunkt, nicht der Endpunkt: seine Tabelle deckt den Vorher/Nachher-Fall
   und den Default, nicht die Überdeckung und nicht die leere Menge.
8. **ADR**, Lastenheft-Bump samt Historie-Zeile mit CR-Bezug (wie seinerzeit bei
   `--yaml`), Spezifikations-Verfeinerung, Handbuch.
9. `make gates`; **Review** und **Verifikation** als getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Anwenden auf den eigenen Bestand.** Die 80 `section-oversized` sind der *Anlass*,
  nicht der Auftrag: ob dieses Repo `max-tasks` über seine Slice-Pläne scharf
  schaltet, ist ein eigener Entscheid nach der Fähigkeit.
- **Kein Grandfathering-Begriff im Werkzeug.** Der Absender grenzt es selbst
  ab: „ab einer Nummer" ist über Globs bzw. RE2 ausdrückbar.
- **Keine abschnitts-übergreifenden Bedingungen.** Ebenfalls vom Absender
  abgegrenzt; andere Frage, eigener Antrag.
- **Keine Änderung an `exclude-sections` des Moduls `vcs`.** Die Namensfrage
  wird entschieden und begründet, nicht durch Umbenennen einer ausgelieferten
  Konfiguration gelöst.

## 4. Definition of Done

- [x] `tasks-ignore-pattern` und `exempt-section-pattern` sind im Schema, im
      `--print-config`-Gerüst, in
      [`spec/lastenheft.md`](../../../../spec/lastenheft.md) (0.76.0 + Historie
      mit CR-Bezug) und in
      [`spec/spezifikation.md`](../../../../spec/spezifikation.md) (Schritt
      1/3/6, §2-Schema, Historie) geführt; nicht kompilierendes RE2 ⇒ Exit 2,
      mit Test — dazu `tasks-ignore-pattern` **ohne** `max-tasks` als halbe
      Aktivierung.
- [x] **Default byte-identisch, gemessen:** derselbe `max-tasks`-Lauf gegen das
      Image **vor** der Änderung und danach liefert **166** Befunde, `diff`
      leer. Ein Test hält zusätzlich die Meldung ohne Muster als Literal.
- [x] **Die Überdeckung ist sichtbar, mit zwei benannten Grenzen:** die
      Meldung nennt die Zahl der ignorierten Items — gemessen an einem
      Probe-Repo: *„Abschnitt trägt 4 Task-Items (3 ignoriert), erlaubt sind
      3"*, dieselbe Zeile in `--doctor`. Ein Muster ohne Treffer
      (**0 ignoriert**) ist als Test gefangen; der Fall *„das Muster nimmt
      alles"* ist auf **Modell**-Ebene geprüft und über Konfiguration
      **unerreichbar** (`max-tasks < 0` ⇒ Exit 2) — das ist die erste Grenze:
      die Sichtbarkeit greift, solange die Regel meldet. Die zweite: ein
      `hint` **ersetzt** die Meldung samt der Zahl, und auch das ist als Test
      gefangen.
- [x] **Die leere Menge ist beantwortet:** `exempt-section-pattern`, das alle
      Abschnitte trifft ⇒ `section-missing` mit Schlüssel und Zahl in der
      Meldung, nicht stilles Grün — und **auch neben einem `hint`**, weil dort
      die Regel nicht gemessen hat ([ADR-0073](../../adr/0073-befund-erlaeuterung-fuer-menschen.md)).
      Mit zwei Tests und end-to-end gefahren.
- [x] **Ein gesetztes `exempt-section-pattern` trägt die Regel-Identität.**
      Ohne das wäre die Paarung, die Grandfathering braucht — Bedingung A für
      alle Abschnitte, B nur für die übrigen —, ein Konfigurations-Duplikat
      (Exit 2, gemessen). Ein leeres Muster geht **nicht** ein, damit kein
      Bestands-`target` wandert; beides als Test.
- [x] **Fence-Treue gemessen** für beide Muster: das Task-Item im Fence zählt
      nicht, die Überschrift im Fence wird weder gewählt noch ausgenommen (die
      Probe dafür ist **diskriminierend** gebaut und gegen die
      Fence-Blindheits-Mutation als rot gemessen). **Inline-Code-Treue gilt
      nur für das Item-Muster** — ein Muster darauf trifft Leerzeichen
      (`make gates`: **null** von 444 Items); `exempt-section-pattern` **sieht**
      Inline-Code, weil es die Zeichenkette mit `section-pattern` teilt. Die
      halbe Zusage steht als solche im Lastenheft und in der Antwort.
- [x] **Umkehr-Proben:** jede Zusage ist von mindestens einer Probe gedeckt,
      und die drei tragenden fahren den **Vorzustand ohne den Schlüssel in
      derselben Funktion** ([`BEO-023`](../observations.md)). **Gemessen an 13
      Mutationen:** jede wird gefangen; drei davon von genau einem Test, die
      übrigen von mehreren — *„genau ein Test je Mutation"* wäre eine andere
      und falsche Zusage. Die eine Probe, die zunächst **gar nicht** fangen
      konnte (Fence-Treue der Abschnitts-Ausnahme), ist repariert und gegen die
      Fence-Blindheits-Mutation als rot gemessen.
- [x] [ADR-0075](../../adr/0075-erklaerte-teilmenge-in-structure.md) begründet
      Verortung, Namen, Muster-Ziele und die Sichtbarkeits-Zusage; im
      [ADR-Index](../../adr/README.md) eingetragen.
- [x] Das [Benutzerhandbuch](../../../user/benutzerhandbuch.md) führt beide
      Schlüssel, das Beispiel **verankert** — samt der **zwei gemessenen
      Fallen** (freies Muster, Muster auf Inline-Code).
- [x] `make gates` grün (Exit explizit) — 606 Dateien, 0 Befunde, Exit 0;
      [**unabhängiger Review**](../../../reviews/2026-08-30-slice-179-strukture-teilmenge-review.md)
      (1 HIGH · 4 MEDIUM · 4 LOW · 1 INFO, blockierend) und
      [**Verifikation**](../../../reviews/2026-08-30-slice-179-strukture-teilmenge-verifikation.md)
      gegen DoD/Spec (A-1 bis A-7) in eigenen Kontexten gelaufen, alle Befunde
      eingearbeitet und die strittigen Zahlen vorher selbst nachgemessen.
- [x] Der Absender bekommt eine
      [**Antwort**](../../cr/2026-08-30-antwort-a-check-structure-teilmenge.md)
      — angenommen mit drei Messungen dagegen, samt den drei Schritten, die er
      beim Umstellen zu gehen hat.

## 5. Abnahme-Punkte / Risiken

- **Ein Ignorier-Muster ist Autoren-Text und kann still zu breit sein.** Die
  Sichtbarkeits-Zusage (§2 Punkt 2) fängt den Fall nur, **solange die Regel
  meldet**. Wer alles ignoriert, sieht nichts. — **Ausgang:** weiter offen, und
  die Sichtbarkeit ist seit dem Review **schwächer**, als dieser Punkt annahm:
  ein `hint` an derselben Regel ersetzt die Meldung samt der Zahl. Beide Grenzen
  stehen jetzt in Lastenheft, Spezifikation und Handbuch; für den
  Nullmengen-Befund ist die Lücke geschlossen.
- **Zwei ähnliche Namen in einem Werkzeug** (`exempt-section-pattern` in `structure`,
  `exclude-sections` in `vcs`). Wer den falschen greift, bekommt eine andere
  Semantik — literal gegen RE2. — **Ausgang:** entfallen — der Name
  `exempt-sections` wurde nicht vergeben. Der Entscheid des Auftraggebers fiel
  auf `exempt-section-pattern`, und damit gibt es die Verwechslung nicht:
  `-pattern` bedeutet in diesem Modul RE2, `exclude-sections` in `vcs`/`sources`
  ist eine Liste literaler Überschriften. Die Begründung steht in
  [ADR-0075](../../adr/0075-erklaerte-teilmenge-in-structure.md) §2 und in der
  Antwort an den Absender.
- **Die Fähigkeit entsteht aus einem fremden Bestand**
  ([`BEO-011`](../observations.md)). Der Anlass ist hier nachgemessen (80 von 89
  Abschnitten), die **Form** stammt aber vom Absender; ob sie zu diesem Repo passt,
  zeigt erst die erste eigene Regel. — **Ausgang:** weiter offen, und die Frage
  ist schärfer geworden. Die Form stammt vom Absender, und zwei seiner drei
  Beispiel-Muster tragen auf **diesem** Bestand nicht — das ist gemessen und ihm
  mitgeteilt. Ob die Fähigkeit hier passt, zeigt weiterhin erst die erste eigene
  Regel; die Entscheidung darüber ist ausdrücklich nicht Teil dieses Slice (§3).
- **Die geprüfte Menge zu verkleinern ist per Konstruktion eine Lockerung.**
  Sie ist als Option ausgeführt und damit opt-in — aber jede Konfiguration, die
  sie setzt, senkt ihre eigene Zusage, ohne dass ein Gate es meldet. —
  **Ausgang:** weiter offen, und dieser Punkt bleibt es dauerhaft. Er ist kein
  Rest, sondern die Form der Zusage: die Sichtbarkeits-Meldung ist das Beste,
  was das Werkzeug dagegen tun kann, und sie hat zwei benannte Grenzen. Das
  Urteil, ob eine Teilmenge berechtigt ist, liegt beim Konfigurations-Autor und
  in keinem Gate.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei, beansprucht am 2026-08-30 —
[slice-176](../done/slice-176-planning-rule-pilot.md) ist geschlossen, und die
beiden anderen offenen Slices warten auf Vorbedingungen (slice-178 auf die
Fence-Messung, slice-172 auf slice-178); dieser auf nichts.

**Rückführungen:** `in-progress` → `open`, falls sich zeigt, dass die
Sichtbarkeits-Zusage nicht tragfähig ist — dann ist die Frage, ob eine
Lockerung ohne Sensor überhaupt in dieses Werkzeug gehört, und der Befund ist
ein anderer.

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

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-30, höchste Kennung
  `BEO-024`): [`BEO-023`](../observations.md) — ein Wächter, der nie fangen
  konnte: der gemessene Defekt des Antrags ist genau das, und die
  Sichtbarkeits-Zusage ist die Antwort darauf;
  [`BEO-011`](../observations.md) — Regel aus dem Anlass: die Form stammt aus
  einem **fremden** Bestand, der Anlass ist hier nachgemessen, und der
  Unterschied steht als Risiko in §5; [`BEO-013`](../observations.md) — ein
  Wächter, der nichts mehr fängt: ein zu breites Ignorier-Muster erzeugt genau
  ihn, ohne je rot zu werden. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-29T07:34:35Z,
  `image-scan.yml` 2026-08-29T10:07:43Z. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption, kein
  Baseline-Abschnitt ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-179. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in).
Module: `structure`. Gates: `make gates`, `make test`, `make doc-check`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Zwei optionale Konfigurationsschlüssel über
vorhandene Bedingungen; kein Fremdsystem, keine Reconciliation, kein Bestand,
der umgestellt werden müsste.

## 9. Closure-Notiz (nach `done/`)

**Geliefert, und der Antrag war beim Bauen der kleinere Teil.** Beide
Optionen sind aufgenommen wie beantragt — als Optionen an bestehenden
Bedingungen, ohne neues Modul, ohne neuen Grund-Code, Default byte-identisch
(gemessen gegen das Vorher-Image: 166 Befunde, `diff` leer, in fünf
Ausgabeformen). Der Aufwand lag woanders: **drei** Messungen gegen den Antrag
und, nach dem Review, **eine gegen die eigene Messung**.

**Die tragende Zahl war falsch, und sie war arithmetisch unmöglich.** Die
Tabelle stellte *„frei 13 Treffer"* neben *„verankert 26"* unter die
Überschrift „über dieselben 444 Items" — und maß zwei **verschiedene**
Ausdrücke. Ein `^` kann die Treffermenge nur verkleinern. Nachgemessen: derselbe
Ausdruck verankert trifft **eines**, und auch das ist ein Liefer-Punkt; die 26
stammen von einem korpus-eigenen Ausdruck. **Die Empfehlung wird dadurch nicht
schwächer, sondern präziser:** verankern allein genügt nicht, das Muster muss
zusätzlich auf Text zielen, den die Bereinigung übrig lässt. Genau daran
scheitert das CR-Beispiel doppelt — und das ist jetzt der Kern der Antwort
statt einer Randnotiz. Die Zahl stand da bereits in einer **immutablen** ADR,
in beiden Spec-Straten, im Handbuch und im ausgehenden Dokument; die Korrektur
läuft über `## Geschichte`, weil eine `Accepted`-Entscheidung nicht
überschrieben wird.

**Die zweite falsche Zahl bezeichnete die falsche Menge.** *„86 Slice-Pläne mit
`## 4. Definition of Done`"* war der `section-missing`-Zähler desselben Laufs —
also die Dateien **ohne** den Abschnitt. Es sind **89** Abschnitte in 175
Dateien. Alles, was daran hängt (444 Items, 80 Befunde, 166 im
Identitäts-Lauf), reproduziert unverändert; nur die Bezugsmenge war falsch
benannt. Beide Fehler sind dieselbe Klasse — [`BEO-020`](../observations.md),
gemessen wird die eine Menge, ausgesagt wird über eine andere —, und beide fand
der zweite Leser, nicht der Schreibende.

**Zwei Verhaltens-Lücken, die kein Gate gefunden hätte.** Ein `hint` löschte
**beide** neuen Sichtbarkeits-Zusagen; der Nullmengen-Befund gehört zur Klasse
*„die Regel hat nicht gemessen"*, die
[ADR-0073](../../adr/0073-befund-erlaeuterung-fuer-menschen.md) ausdrücklich vom
verfassten Hinweis ausnimmt, und trägt jetzt seine eigene Meldung. Für die
`section-oversized`-Zahl gewinnt der Hinweis zu Recht — das ist eine **zweite**
benannte Grenze der Sichtbarkeits-Zusage. Und das gesetzte
`exempt-section-pattern` ging nicht in die **Regel-Identität** ein, womit
ausgerechnet die Paarung, für die der Schlüssel beantragt wurde — Bedingung A
für alle Abschnitte, B nur für die übrigen —, als Konfigurations-Duplikat mit
Exit 2 endete.

**Der Test gegen [`BEO-023`](../observations.md) war selbst ein
`BEO-023`-Fall.** Die Fence-Probe der Abschnitts-Ausnahme richtete ihr Muster
auf die **gefencte** Überschrift; damit war *„nie gewählt"* von *„gewählt und
dann ausgenommen"* am Ergebnis nicht unterscheidbar, und eine
Fence-Blindheits-Mutation lief grün durch. Der Kommentar **über** diesem Test
schrieb *„eine geerbte Zusage ohne Probe ist eine Erinnerung"*. Die Klasse war
beim Schreiben benannt und im selben Atemzug begangen. Reparatur: eine Zeile;
die Mutation macht ihn jetzt rot, gemessen.

**Eine Zuschreibung an den Absender war unfair und ist zurückgenommen.** Sein
Abschnitts-Muster ist gegen **seine** deklarierte Semantik korrekt — er schreibt
ausdrücklich *„Abschnitte, deren Überschriftstext dieses RE2 trifft"*. Dass es
hier nicht greift, folgt aus **unserer** Entscheidung, und am Antrag wurde
nichts gemessen. Das stand als „dritter gemessener Punkt" in ADR, Code-Kommentar
und Antwort; die Entscheidung bleibt, ihre Begründung steht jetzt auf der
Konsistenz **innerhalb** der Regel statt auf einem fremden Fehler
([`BEO-012`](../observations.md)).

**Was der Slice über das Werkzeug hinaus gezeigt hat:** die Fence-Treue, die
dieses Modul auszeichnet, ist bei einem **Ignorier**-Muster eine Falle. Wer
fence-treu zählt, ignoriert auch fence-treu — und wer sein Muster auf einen
Ausdruck in Backticks richtet, richtet es auf Leerzeichen. Das gilt nicht nur
hier: die Slice-Vorlage der Baseline schreibt `` - [ ] `make gates` grün ``
selbst, das Beispiel-Muster des Absenders ist also in **jedem** Repo inert, das
die Vorlage benutzt. Diese Verallgemeinerung stammt aus der Verifikation und
ist stärker belegt als das, was die Dokumente behaupten.

**Fortgeschrieben:** [`BEO-020`](../observations.md) und
[`BEO-023`](../observations.md) auf Zähler **3** — beide haben damit die
Schwelle erreicht; die Verkörperung fällt an die Closure von
[welle-86](welle-86/welle-86-closure-uebergang-durchsetzen.md) (Lese-Schritt).
[`BEO-012`](../observations.md) ist als achte Instanz nicht gezählt, weil die
Zuschreibung im selben Slice entstand und zurückgenommen wurde, bevor sie
wirkte — sie steht hier, nicht im Register.

**Was offen bleibt und wohin es gehört:** ob dieses Repo `max-tasks` über seine
eigenen Slice-Pläne scharf schaltet, ist ein eigener Entscheid nach der
Fähigkeit (§3) — mit der eben gemessenen Vorwarnung, dass ein tragfähiges
Muster hier auf `^grün \(Exit explizit\)` zielen müsste und damit nur 26 der
444 Items erreicht. Der `CHANGELOG`-Eintrag gehört in die Release-Prep, wie
schon der von slice-177.
