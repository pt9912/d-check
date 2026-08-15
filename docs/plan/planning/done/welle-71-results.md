# Welle 71 — Die Closure-Fähigkeit wird Obermenge des Konsumenten-Skripts — Closure-Notiz

**Welle:** welle-71-closure-konsumenten-paritaet
**Abschluss:** 2026-08-10
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **`planning.closure.glob`** ([slice-097](slice-097-closure-glob-entkopplung.md),
  Release **v0.54.0**, [ADR-0051](../../adr/0051-eigener-kandidaten-filter-closure.md)):
  ein eigener Kandidaten-Filter für die Closure-Fähigkeit, Default als
  **Verweis** auf `slice-glob`.
- **`closure-note-placeholder`** ([slice-098](slice-098-closure-note-placeholder.md),
  Release **v0.55.0**, [ADR-0052](../../adr/0052-platzhalter-erkennung-inline-code.md)):
  die vierte, opt-in Bedingung — der unausgefüllte Rumpf einer Vorlage.
- **Das Wellen-Ziel ist erreicht:** die Closure-Fähigkeit deckt jetzt beide
  Lücken, um derentwillen der Konsument sein handgeschriebenes Prüfskript noch
  fuhr. Beide Bausteine sind **veröffentlicht** — ohne Release hätte die Welle
  ihren Zweck nicht erfüllt, das war von Anfang an ein Closure-Kriterium.

## Was hat funktioniert?

- **Erst messen, dann entscheiden — zum dritten und vierten Mal.** In beiden
  Slices hat die Bestandsmessung die Entscheidung **gedreht**, nicht bestätigt.
  Bei 097 fand sie einen realen Rückstand, bevor der Schlüssel überhaupt
  existierte; bei 098 entschied sie die Erkennungs-Grenze (24 → 12 → 0 Treffer).
- **Der unabhängige Review vor jedem Release.** Vier Runden über zwei Slices,
  **alle vier blockierend**, jede mit einem Befund, den ich nicht gesehen hatte —
  darunter zweimal eine Zusage, die ich selbst geschrieben hatte.
- **Die geteilte Lexik benutzen statt nachbauen.** 098 setzt auf
  `PreprocessMarkdown` auf; kein zweiter Inline-Code-Scanner, keine zweite
  Absatz-Semantik.

## Was ging anders als geplant?

- **Der Zuschnitt der Welle wich von der Roadmap ab.** Geplant waren drei Slices
  (094/097/098), geliefert wurden zwei. slice-094 blieb bewusst draußen: die
  beiden hier sind rein **additiv**, 094 ändert die Semantik eines
  ausgelieferten Gates. Das ist eine andere Risiko-Klasse und verdient eine
  eigene Welle mit eigener Messung.
- **Meine Abnahme-Entscheidung in 097 war falsch und wurde vom Review
  gekippt.** Ich hatte die eigene Prüfmenge geweitet, weil „mehr geprüft“
  besser klang — gemessen baute es ein **Falsch-Negativ**: bei den
  Ergebnisnotizen wird der Dokument-Titel zur Abschnitts-Überschrift, der
  gemessene Abschnitt zur ganzen Datei.
- **Eine Messung belegt, was da ist, nicht was möglich ist.** In 098 stimmte
  meine Vertragszusage zu den Vergleichszeichen nur für die Schreibweise mit
  Leerzeichen; die enge Form kennt der eigene Bestand nicht, also fand die
  Messung sie nicht.
- **Zwei Zusagen waren praktisch unbewacht** (33 von 35 HTML-Tag-Einträgen
  löschbar; eine kopierte Grenz-Logik). Beide sind jetzt testgehalten bzw.
  zusammengeführt.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 4 — was aus dieser Welle in die Steuerung
zurückfließt.

- **Eine Bestandsmessung am eigenen Repo ist eine Aussage über den Bestand, nicht
  über den Raum der Eingaben.** Wo eine Zusage eine *Form* ausschließt, gehört
  eine konstruierte Gegenprobe dazu — nicht nur die Korpus-Zahl.
- **Listen im Code brauchen einen Test, der über die Liste iteriert.** Eine
  Stichprobe hält 2 von 35 Einträgen.
- **Bei einer Weitung der Prüfmenge fragen: welcher Abschnitt wird danach
  gemessen?** Mehr Kandidaten sind nicht automatisch mehr Prüfung.

## Beobachtungs-Register (Zeiger)

Der Lese-Schritt dieser Welle: das Register führt **BEO-001** (1×),
**BEO-002** (2×), **BEO-003** (1×) und **BEO-004** (3×).

- **BEO-003** war in beiden Slices **einschlägig und wirksam**: in 097 verhinderte
  sie den kopierten Literal-Default, in 098 den zweiten Inline-Code-Scanner. Der
  Zähler bleibt bei 1 — die Klasse ist hier **vermieden**, nicht eingetreten.
- **BEO-004** steht weiter bei 3 und damit an der Verkörperungs-Schwelle. Diese
  Welle hat sie zweimal berührt: 097 hat die Grenze der Closure-Fähigkeit
  ausgemessen, 098 hat zwei Lexik-Grenzen benannt statt geschlossen. Beides sind
  **benannte** Grenzen, keine neuen Belege — der Zähler wächst nicht.
- **Nichts zu streichen.** Die Verkörperung von BEO-004 bleibt der nächsten
  Planung überlassen; sie ist kein Slice dieser Welle.

## Folge-Slices

- [slice-094](../done/slice-094-closure-zaehl-paritaet.md) und
  [slice-104](../done/slice-104-floskel-wortgrenze.md) — beide ändern die
  Semantik des **ausgelieferten** Closure-Gates und sind als gemeinsame
  Folge-Welle vorgemerkt.
- [slice-099](slice-099-structure-modul.md) — die Abdeckung der
  Wellen-Dokumente, die 097 ausdrücklich **nicht** über den Kandidaten-Filter
  gelöst hat.

## Verifikation

- **Closure-Trigger erfüllt:** beide Slices dieser Welle liegen in `done/`, und
  die Konsumenten-Bedingung ist mit zwei veröffentlichten Releases eingelöst.
- `make fullbuild` grün; `make ci` vor jedem Tag grün.
- Releases **v0.54.0** (Digest `sha256:8ee8de4c…4afb`, Lauf 31374890120) und
  **v0.55.0** (Digest `sha256:b8b17b08…3325`, Lauf 31386249345), beide Pipelines
  erfolgreich, beide Digest-Backfills erledigt.
- **Trigger-Audit** über die drei Artefaktklassen: keine offenen Carveouts; keine
  stehengebliebene Gate-Reifestufe; von den Re-Evaluierungs-Triggern der beiden
  neuen ADRs ist keiner eingetreten — insbesondere hat kein Konsument eine
  Platzhalter-Form gemeldet, die die Erkennung verfehlt, und der Verweis-Default
  hat niemanden überrascht.
