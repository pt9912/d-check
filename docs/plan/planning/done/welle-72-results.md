# Welle 72 — Die Semantik eines ausgelieferten Gates nachziehen — Closure-Notiz

**Welle:** welle-72-closure-semantik
**Abschluss:** 2026-08-10
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Zähl-Parität** ([slice-094](welle-72/slice-094-closure-zaehl-paritaet.md)): der
  Closure-Abschnitt wird **einmal** bereinigt (Fenced-Code **und** Inline-Code),
  alle Bedingungen lesen diesen einen Text, ein Satzende zählt nur vor
  Whitespace oder Zeilenende.
  [ADR-0053](../../adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md).
- **Wortgrenzen** ([slice-104](welle-72/slice-104-floskel-wortgrenze.md)): die
  Floskel-Bedingung vergleicht an Wortgrenzen statt als Teilstring.
- **Ein Release, eine Notiz** — **v0.56.0**. Das war das Wellen-Ziel jenseits
  beider DoDs, und es ist eingelöst: der CHANGELOG-Eintrag nennt **alle drei**
  Richtungen einzeln, mit Anlass und Preis, statt sie über zwei Releases zu
  verteilen.

## Was hat funktioniert?

- **Die Welle als Klammer für eine Risiko-Klasse.** Beide Slices ändern die
  Semantik eines **ausgelieferten** Gates; sie zusammenzunehmen war keine
  Bequemlichkeit, sondern die einzige Form, in der ein Konsument die Änderung
  **einmal** lesen kann.
- **Erst messen, dann entscheiden — zum fünften und sechsten Mal.** Beide
  Messungen haben Entscheidungen erzeugt, die im Plan nicht standen: die 170
  fett gesetzten Satzenden (slice-094) und die Auswahl der neu aufnehmbaren
  Phrasen (slice-104).
- **Der Paritäts-Beleg gegen den realen Adopter-Bestand**, nicht gegen Fixtures:
  84 von 84 Notizen, identische Zählung, symmetrisch an der Schwelle. Vom
  Reviewer mit **eigener** Methode nachgemessen und bestätigt.

## Was ging anders als geplant?

- **Beim Wellen-Schnitt fiel eine Kopplung auf, die in keinem der beiden Slices
  stand:** beide lockern die **Floskel-Prüfung**, aus verschiedenen Richtungen.
  Erst dadurch wurde klar, warum sie in **eine** Release-Notiz gehören.
- **CRLF war der schwerste Befund — und dieses Repo kann ihn strukturell nicht
  finden.** Der eigene Bestand ist LF-only; kein Gate sieht die Klasse je.
  **Zweite CRLF-Regression an einem Tag** (die erste in slice-101), beide Male
  von einer konstruierten Gegenprobe gefunden, nie vom Korpus.
- **Drei Vertragsflächen standen nach der Angleichung auf der widerrufenen
  Fassung** und waren gegen den Lauf falsifizierbar. Dieselbe Rand-Klasse wie
  **BEO-002**, diesmal dreifach in einem Slice.
- **Zwei Zusagen waren breiter als belegt** (die Reichweite der Lockerung, die
  Universalität der Parität). Beide sind jetzt so formuliert, wie sie belegt
  sind — nicht zurückgenommen, sondern skopiert.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 4 — was aus dieser Welle in die Steuerung
zurückfließt.

- **Eine Semantik-Änderung an einem ausgelieferten Gate gehört mit ihren
  Geschwistern in ein Release.** Zwei Schritte hintereinander verschieben dem
  Konsumenten zweimal den Boden.
- **Zeilenende-Formen brauchen eine konstruierte Gegenprobe.** Ein LF-only-Repo
  kann CRLF-Defekte nicht messen — zweimal an einem Tag belegt.
- **Wer eine Zusage formuliert, prüfe ihren Skopus:** „eine zitierte Floskel“
  und „jeder Inline-Code-Span“ sind nicht dasselbe, und
  „Parität“ am gemessenen Bestand ist nicht Parität überhaupt.

## Beobachtungs-Register (Zeiger)

Der Lese-Schritt dieser Welle: das Register führt **BEO-001** (1×),
**BEO-002** (2×), **BEO-003** (1×) und **BEO-004** (3×).

- **BEO-002** ist in slice-104 **eingetreten** — drei Vertragsflächen blieben
  nach der Semantik-Änderung stehen. Der Zähler steigt auf **3** und erreicht
  damit die Verkörperungs-Schwelle: die Klasse ist jetzt zweimal so oft
  aufgetreten wie sie gezählt wurde, und die naheliegende Form ist eine Regel in
  der Autoritäts-Doku, nicht ein weiterer Zählschritt.
- **BEO-003** war in slice-094 **einschlägig und wirksam**: sie hat die
  Alternative „zwei getrennte bereinigte Texte“ verworfen. Vermieden, nicht
  eingetreten — der Zähler bleibt bei 1.
- **BEO-004** unberührt.

## Folge-Slices

- [slice-099](welle-69/slice-099-structure-modul.md) — das Modul `structure`;
  Trigger seit welle-70 eingetreten.
- [slice-103](slice-103-geteilte-lexik-raender.md) — dieselbe
  Drift-Klasse in anderen Lexiken.

## Verifikation

- **Closure-Trigger erfüllt:** beide Slices in `done/`, **ein** Release
  (**v0.56.0**) mit **einer** Notiz, die jede Richtung einzeln nennt.
- Der eigene Bestand wurde **vor** jeder Änderung gemessen und ist danach grün
  (Minimum fiel von 7 auf 5, Schwelle bleibt 4).
- `make fullbuild` grün; `make ci` vor dem Tag grün.
- **Trigger-Audit** über die drei Artefaktklassen: keine offenen Carveouts, keine
  stehengebliebene Gate-Reifestufe, kein eingetretener Re-Evaluierungs-Trigger
  der ADR dieser Welle — insbesondere hat kein Konsument gemeldet, dass ihm durch die
  gelockerte Floskel-Prüfung ein echter Fall entgeht.
