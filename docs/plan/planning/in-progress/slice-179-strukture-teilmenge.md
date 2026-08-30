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

**Der Anlass gilt in diesem Repo genauso, und das ist nachgemessen.** Ein Lauf
von `max-tasks: 3` über die eigenen Slice-Pläne liefert **160 Befunde bei 223
Dateien**. Die jüngsten Slices tragen 5 bis 9 DoD-Items, von denen der
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
die **86** DoD-Abschnitte dieses Repos (**444** Task-Items, `max-tasks: 0`,
damit jeder Abschnitt seine Zahl meldet):

| Musterform | Items ignoriert | davon falsch |
|---|---|---|
| frei, wie im CR | 13 | **2** |
| verankert (das Item **beginnt** mit dem Ausdruck) | 26 | 0 |

**Erster Defekt: das freie Muster nimmt echte Zusagen mit.** Zwei der 13 sind
**Liefer**-Punkte, die die Closure-Notiz nur als ihren **Ort** nennen — *„§2/§3/§4/§6
tragen `SPEC-NNN` fortlaufend; Zählung in der Closure-Notiz"* und
*„[ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)-§Kern-Messung in der Closure-Notiz"*. Ein freies Substring-Muster
entfernt sie aus der Zählung, ohne dass jemand es sieht — die Gestalt aus
[`BEO-023`](../observations.md): ein Filter, der still das Falsche entfernt,
sieht aus wie einer, der richtig filtert. Die verankerte Form trifft beide
nicht, weil keiner mit dem Ausdruck **beginnt**.

**Zweiter Defekt, und er ist der schärfere: die erste Alternative des CR ist
hier komplett wirkungslos.** `make gates` trifft **null** von 444 Items —
gemessen, nicht geschätzt. Der Grund ist die Zusage, die der CR selbst
verlangt: die Zählung liest den **bereinigten** Text, und in **allen** 86
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

- **Kein Anwenden auf den eigenen Bestand.** Die 160 Befunde sind der *Anlass*,
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
- [x] **Die Überdeckung ist sichtbar:** die Meldung nennt die Zahl der
      ignorierten Items — gemessen an einem Probe-Repo: *„Abschnitt trägt 4
      Task-Items (3 ignoriert), erlaubt sind 3"*, dieselbe Zeile in `--doctor`.
      Ein zu breites Muster (**7 ignoriert**) und eines ohne Treffer
      (**0 ignoriert**) sind je als Test gefangen.
- [x] **Die leere Menge ist beantwortet:** `exempt-section-pattern`, das alle
      Abschnitte trifft ⇒ `section-missing` mit Schlüssel und Zahl in der
      Meldung, nicht stilles Grün. Mit Test und end-to-end gefahren.
- [x] **Fence-Treue gemessen** für beide Muster: das Task-Item im Fence zählt
      nicht, die Überschrift im Fence wird weder gewählt noch ausgenommen —
      **und die teure Hälfte steht dabei:** ein Muster auf einen Ausdruck in
      **Inline-Code** trifft Leerzeichen (`make gates`: **null** von 444 Items).
- [x] **Umkehr-Proben** je Zusage, jede von genau einem Test gefangen — die
      drei tragenden fahren den **Vorzustand ohne den Schlüssel in derselben
      Funktion** ([`BEO-023`](../observations.md)).
- [x] [ADR-0075](../../adr/0075-erklaerte-teilmenge-in-structure.md) begründet
      Verortung, Namen, Muster-Ziele und die Sichtbarkeits-Zusage; im
      [ADR-Index](../../adr/README.md) eingetragen.
- [x] Das [Benutzerhandbuch](../../../user/benutzerhandbuch.md) führt beide
      Schlüssel, das Beispiel **verankert** — samt der **zwei gemessenen
      Fallen** (freies Muster, Muster auf Inline-Code).
- [ ] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.
- [x] Der Absender bekommt eine
      [**Antwort**](../../cr/2026-08-30-antwort-a-check-structure-teilmenge.md)
      — angenommen mit drei Messungen dagegen, samt den drei Schritten, die er
      beim Umstellen zu gehen hat.

## 5. Abnahme-Punkte / Risiken

- **Ein Ignorier-Muster ist Autoren-Text und kann still zu breit sein.** Die
  Sichtbarkeits-Zusage (§2 Punkt 2) fängt den Fall nur, **solange die Regel
  meldet**. Wer alles ignoriert, sieht nichts. — **Ausgang:** *(bei Closure)*
- **Zwei ähnliche Namen in einem Werkzeug** (`exempt-sections` in `structure`,
  `exclude-sections` in `vcs`). Wer den falschen greift, bekommt eine andere
  Semantik — literal gegen RE2. — **Ausgang:** *(bei Closure)*
- **Die Fähigkeit entsteht aus einem fremden Bestand**
  ([`BEO-011`](../observations.md)). Der Anlass ist hier nachgemessen (160 von
  223), die **Form** stammt aber vom Absender; ob sie zu diesem Repo passt,
  zeigt erst die erste eigene Regel. — **Ausgang:** *(bei Closure)*
- **Die geprüfte Menge zu verkleinern ist per Konstruktion eine Lockerung.**
  Sie ist als Option ausgeführt und damit opt-in — aber jede Konfiguration, die
  sie setzt, senkt ihre eigene Zusage, ohne dass ein Gate es meldet. —
  **Ausgang:** *(bei Closure)*

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
