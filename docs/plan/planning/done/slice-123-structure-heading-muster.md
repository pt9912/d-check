# Slice slice-123: `structure` — jede Überschrift des Abschnitts matcht ein Muster

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-82-config-flaechen](../welle-82-config-flaechen.md) (zugeordnet
bei der Eröffnung).

**Bezug:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(die zu erweiternde Anforderung),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`MR-025`](../../../../harness/conventions.md#mr-025); Anlass ist die
Ersatz-Konstruktion aus welle-80 — eine ausgeschriebene Präfix-Negation, weil
RE2 keinen Lookahead kennt, mit einem belegten **stillen Falsch-Negativ**.

**Berührte Spec-Stellen:**
[`DC-FA-STRUCT-001.a`](../../../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
(Bedingungen), §2-Schema (`structure[].*`), §4 (`section-*`) — der Verweis zeigt aufwärts.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Die Bedingung „**jede** Überschrift dieses Abschnitts genügt einem Muster" ist
mit den heutigen Schlüsseln nur als Negation ausdrückbar — und die war in
welle-80 nachweislich falsch: sie sprach nicht die Heading-Lexik des Moduls
(führender Weißraum, Tab als Trenner), also entkam eine eingerückte Sektion
still. Ein eigener Schlüssel dreht die Aussage um: **positiv, je Überschrift,
mit der Lexik des Moduls** — das Modul kennt seine Überschriften bereits, es
muss sie nur einzeln prüfen statt den Abschnittstext als Ganzes.

## 2. Vorgehen

1. **CR-Commit zuerst:** Lastenheft
   [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
   um die Bedingung erweitern — Beschreibung, Akzeptanzkriterien (Happy-Path, Negativ je
   verletzender Überschrift mit ihrer Zeile als Befund-Ort, Default
   byte-identisch ohne den Schlüssel, fail-closed bei ungültigem Muster),
   §7-Historie.
2. **ADR der Welle** um die zweite Entscheidung ergänzen: positive
   Je-Überschrift-Prüfung statt Negation im Abschnittstext; Ebenen-Wahl
   (welche Überschriften-Ebenen der Abschnitt umfasst); Befund je Überschrift
   statt einer je Abschnitt, damit die Zeile zeigt, wo es klemmt.
3. **Spezifikation:** Bedingung im Algorithmus, §2-Schema, §4-Zeile
   (bestehender oder neuer Grund-Code — die Entscheidung gehört in die ADR und
   folgt dem Kriterium „andere Reparatur ⇒ eigener Code").
4. **Code + Tests:** die Überschriften des Abschnitts liegen im Modul bereits
   vor (geteilte Lexik) — sie werden einzeln geprüft; Tests für alle
   Akzeptanzkriterien, dazu die Fälle, an denen die alte Negation scheiterte
   (eingerückt, Tab-getrennt, vierte Ebene, Überschrift nur aus Inline-Code).
5. **Das eigene Profil umstellen:** die ausgeschriebene Negation wird durch den
   neuen Schlüssel ersetzt — vorher/nachher gemessen, beide Male grün, und die
   Gegenprobe, an der die Negation still war, wird jetzt rot.
6. **Spiegel** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
   Anforderung, Algorithmus, §2-Schema, §4-Tabelle, Klartexte,
   `--print-config`-Vorlage, Config-Kommentar (Handbuch ist Release-Prep).
7. Unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine weiteren `structure`-Bedingungen** — nur diese eine.
- **Kein Handbuch, kein CHANGELOG** (slice-125).
- **Keine Default-Änderung.**

## 4. Definition of Done

- [x] CR-Commit (Lastenheft allein) liegt **vor** Spezifikation und Code —
      Lastenheft 0.64.0 allein, danach ADR, dann Spezifikation samt Code; die
      Review-Auflagen als eigener Nachzug.
- [x] Der Schlüssel prüft **je Überschrift** mit der Modul-Lexik; Befund nennt
      die verletzende Zeile — über `SectionHeadings`, das dieselbe Erkennung
      benutzt wie die Abschnitts-Findung.
- [x] Das eigene Profil nutzt ihn statt der Negation — **die Gegenprobe ist
      eine andere als geplant.** Der eingerückte Fall ist seit dem
      slice-114-Review nicht mehr still (die Negation wurde damals auf die
      Modul-Lexik korrigiert); die Annahme in §1/§2.5 war überholt. Der Review
      hat den Fall gefunden, an dem sie **heute** still ist: eine Überschrift
      innerhalb eines mehrzeiligen Inline-Code-Spans. Gemessen: alte
      Konstruktion Exit 0, neue Exit 1 auf der Zeile der Überschrift.
- [x] Default-Beweis byte-identisch; `make gates` grün; unabhängiger Review;
      Closure-Notiz; Register gesichtet — der Review war **blockierend** (ein
      HIGH, sechs MEDIUM, vier LOW), alle Befunde eingearbeitet und im
      [Report](../../../reviews/2026-08-22-slice-123-structure-heading-muster-review.md)
      belegt.

## 5. Abnahme-Punkte / Risiken

- **Welche Überschriften gehören zum Abschnitt?** Die Ebenen-Frage entscheidet
  über Falsch-Positive (eine tiefere Ebene, die nie gemeint war). Sie gehört in
  die ADR und in den Vertrag, nicht in den Code. — **Ausgang:** entschieden und
  begründet (Default = Abschnitts-Ebene + 1, `headings-level` wählt eine
  andere). Der Review hat die Frage **je Abschnitt** statt je Regel
  nachgemessen — bei einem Selektor, der zwei Ebenen trifft, prüft jeder
  Abschnitt seine eigene Kind-Ebene. Und er hat eine Grenze gefunden, die der
  Plan nicht sah: zwei Ebenen in **einer** Regel sind heute nicht bloß
  unvorgesehen, sondern durch die Regel-Identität **gesperrt** (Exit 2). Steht
  jetzt als dritte Grenze in der ADR.
- **Ein Befund je Überschrift statt je Abschnitt** ändert die Befund-Zahl —
  das ist gewollt (die Zeile zeigt, wo es klemmt), muss aber zugesagt sein. —
  **Ausgang:** zugesagt in Anforderung, Algorithmus und §4-Zeile; ein eigener
  Grund-Code, weil die Reparatur eine andere ist. Getestet mit zwei
  Verletzungen in einem Abschnitt: zwei Befunde auf **ihren** Zeilen.
- **Der Umstieg des eigenen Profils ist der eigentliche Beweis:** wenn die neue
  Bedingung die alte nicht deckt, zeigt es sich hier. — **Ausgang:** gehalten,
  und der Beweis ist stärker als geplant. Über **22** konstruierte Formen fand
  der Review **keinen** Fall, in dem der neue Schlüssel still bleibt und die
  Negation eine echte Überschrift gemeldet hätte — mit Strukturargument: die
  Bereinigung ersetzt Zeichen nur durch Leerzeichen, sie kann ein
  `SPEC-NNN␣`-Präfix nie erzeugen. Der Befundsatz des eigenen Repos ist vor und
  nach der Umstellung byte-identisch.

## 6. Trigger

**Start** (`open` → `in-progress`): Wellen-Eröffnung (eingetreten).

**Rückführungen:** `in-progress` → `next`, falls die Ebenen-Frage einen zweiten
Schlüssel verlangt (dann erst Vertrag klären).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Kern (GF), Config-Rand (GF), Spec-Straten (GF),
  eigenes Prüf-Profil (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): **BEO-003**
  (geteilte Lexik driftet an den Rändern) ist einschlägig — der neue Schlüssel
  muss die Lexik des Moduls **benutzen**, nicht nachbauen; BEO-002 als
  Spiegel-Pflicht, BEO-006/009/010 als Arbeitsregeln.

Slice-ID: slice-123. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus). Module:
`structure` (Kern `rules/`), Config-Rand, Spec, eigenes Profil. Gates:
`make test` (eng), `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — additive Erweiterung einer eigenen,
spezifizierten Anforderung; sie löst zugleich eine Ersatz-Konstruktion ab.

## 9. Closure-Notiz (nach `done/`)

Geliefert ist die achte `structure`-Bedingung: *jede* Überschrift des
Abschnitts genügt einem Muster — positiv, je Überschrift, auf ihrer Zeile, mit
der Erkennung des Moduls statt einer nachgebauten. Das eigene Profil führt sie
seither statt der ausgeschriebenen Präfix-Negation.

**Die Lehre steht in der Begründung, nicht im Ergebnis.** Der Plan wollte
heilen, was schon geheilt war: die Negation war seit dem slice-114-Review nicht
mehr zu eng. Ich habe acht Formen gemessen, alle identisch, und daraus
geschlossen, die Umstellung sei „verhaltenserhaltend, nicht heilend". Der
Review fand die neunte — eine Überschrift in einem mehrzeiligen
Inline-Code-Span, wo die Bereinigung die Zeile leerräumt und die Negation
schweigt. **Acht Messungen sind acht Messungen, keine Eigenschaft.** Wer einen
Schluss zieht, muss die Menge nennen, über die er gilt, oder die N+1-te Form
suchen.

Der zweite Fund ist ein Name. `heading-pattern` hätte in **einem** Profil zwei
gegensätzliche Rollen getragen: Selektor unter `planning.closure`, Bedingung
unter `structure`. Umbenannt zu `headings-match`/`headings-level`, solange
nichts released ist — danach wäre derselbe Schnitt ein Breaking Change gewesen.
Das ist der eigentliche Wert des Reviews **vor** dem Release, nicht danach.

Drei Widersprüche derselben Klasse kamen dazu: „als einzige" stand an drei
Stellen im Lastenheft, wo Spezifikation und ADR „die zweite" sagten; Schritt 5
der Spezifikation widersprach Schritt 6 desselben Dokuments; zwei
Code-Kommentare trugen dieselbe überholte Invariante. Wer eine Ausnahme zur
zweiten macht, muss jede Stelle finden, die „die einzige" sagt — und das ist
kein `grep` nach dem neuen Wortlaut, sondern nach dem alten.
