# ADR-0056 — `links.resolve-from`: nur wandernde Quellorte, und Ziel-Identität gehört zur Auflösbarkeit

**Status:** Accepted
**Datum:** 2026-08-16
**Autor:** pt9912
**Bezug:** [`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
(Erweiterung — Schnitt-Kriterium aus
[ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md): Einzelmodul-Frage,
kein neues Kürzel); Konventions-Anlass
[`MR-013`](../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)
(deren Invariante hier maschinell wird); Change Request des Konsumenten
`a-check` (CR 2 seiner Werkzeug-Abdeckungs-Analyse)

## Kontext

Der Planning-Lifecycle ist eine Zustandsmaschine über Verzeichnisse, und der
Zustandswechsel ist ein `git mv` **ohne** Inhaltsänderung. Ein präfixloser
Nachbar-Verweis ist am Ist-Ort grün und bricht beim nächsten Wechsel — gemeldet
erst **nach** dem Move, wenn die Reparatur die Move-Regel verletzt.

Die Klasse ist nicht hypothetisch. Der historische Maximalfall: **ein**
`git mv` brach 19 Verweise auf einen Schlag. Und sie wiederholt sich bei jedem
Move — an einem einzigen Tag wurden bei drei Lifecycle-Moves 10, 15 und 14
Verweise von Hand nachgezogen. Die Bestandsmessung gegen den Stand vor dem
damaligen Move liefert 19 positionsabhängige Verweise — 15 davon Teil des
realen Bruchs; dessen übrige vier waren Links eines Review-Reports auf den
verschobenen Slice (Quelle ortsfest, **Ziel** gewandert), eine Klasse, die
diese Fähigkeit strukturell nicht deckt. Die Zahlengleichheit ist Zufall.

Dieselbe Messung hat den naiven Zuschnitt widerlegt: über **alle vier**
Lifecycle-Verzeichnisse gerechnet wären heute 108 Verweise „positionsabhängig" —
die Spitzenreiter sind Wellendokumente und Slices im **Ruheort**, die nie wieder
wandern. Eingeschränkt auf die wandernden Verzeichnisse: 0 von 79.

## Entscheidung

1. **Erweiterung von [`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links), kein neues Modul und kein neues
   Kürzel.** Es ist dieselbe Prüfung (löst ein Verweis auf?) mit erweiterter
   Quellort-Menge — keine neue Frage. Opt-in über `links.resolve-from`; ohne
   den Block ist der Befundsatz byte-identisch.

2. **Konfiguriert werden Gruppen wandernder Verzeichnisse, und nur Dateien
   darin sind Quellen.** Eine Gruppe ist eine Menge von Geschwister-Orten, über
   die Dateien wandern (`dirs`), plus optional die Menge der **ortsfesten**
   Orte, die als hypothetische Ziele mitzählen, deren Dateien aber keine
   Quellen sind (`fixed-dirs`, etwa der Ruheort): eine Datei dort ist am
   Endzustand, ihre Verweise müssen nur vom Ist-Ort auflösen — die Messung
   zeigt sonst 108 Falsch-Positive. Eine Datei aus einem `dirs`-Verzeichnis
   muss dagegen von **jedem** Ort der Gruppe (`dirs` ∪ `fixed-dirs`) auflösen,
   denn dorthin kann sie wandern.

3. **Positionsabhängig heißt: nicht überall auflösbar ODER nicht überall
   dasselbe Ziel.** Ein Verweis, der von jedem Ort auflöst, aber auf
   verschiedene Ziele, meint je nach Ort etwas anderes — das ist dieselbe
   Fehlerklasse in stiller Form. Beide Fälle tragen den **einen** neuen
   Grund-Code `link-position-dependent`; die Meldung nennt den nicht
   auflösenden Ort bzw. die divergierenden Ziele.

4. **Eigener Grund-Code statt `target-missing`.** Die Reparatur ist eine
   andere (Pfad präfixieren statt Ziel anlegen), und am Ist-Ort ist nichts
   kaputt — ein `target-missing` wäre eine falsche Aussage über den Ist-Zustand.

5. **Anker bleiben außen vor; die Ziel-Menge ist die der bestehenden
   Auflösung** (Links und Bilder, dieselbe Vorverarbeitung, dieselbe
   Prozent-Dekodierung, dieselben Ventile). Die Anker-Frage hängt am
   Ziel-Dokument, nicht am Quellort.

6. **Die Ziel-Wanderung ist eine benannte Grenze, keine Fähigkeit.** Bricht ein
   Verweis, weil sein **Ziel** wandert und die Quelle ortsfest ist (der
   Review-Report, der auf einen `in-progress/`-Slice zeigt), meldet diese
   Prüfung nichts — sie prüft hypothetische **Quell**-Orte. Bewusst: für
   Lauf-Belege ist `ignore-refs` das etablierte Ventil, lebende Dokumente
   zieht die Move-Regel im selben Commit nach, und eine Ziel-Hypothese müsste
   raten, **wohin** jedes Ziel wandern kann — das ist eine andere, teurere
   Frage.

## Alternativen

- **Alle Lifecycle-Verzeichnisse als Quellen.** Verworfen durch die Messung:
  108 Falsch-Positive auf ortsfesten Dokumenten — die Fähigkeit wäre am
  eigenen Bestand unbrauchbar.
- **Nur Auflösbarkeit prüfen, Ziel-Identität ignorieren.** Verworfen: ein
  Verweis, der überall auflöst, aber auf verschiedene Dateien, ist die stillere
  Hälfte derselben Klasse.
- **`target-missing` wiederverwenden.** Verworfen nach Entscheidung 4.
- **Ein Auto-Rewrite beim Move.** Verworfen: d-check bleibt diagnose-only in
  dieser Klasse; welcher Pfad gemeint ist, ist eine Autoren-Entscheidung.

## Konsequenzen

- **Die Fähigkeit findet heute im eigenen Baum nichts** — die Null ist von
  einer Woche Hand-Nachzügen gekauft, und genau diese Arbeit entfällt künftig:
  der nächste präfixlose Nachbar-Verweis wird **vor** dem Move gemeldet.
- Die Prüf-Menge wächst pro Quell-Datei um |Gruppe|−1 zusätzliche Auflösungen.
  Die Laufzeit-Zusage ([`DC-QA-01`](../../../spec/lastenheft.md#dc-qa-01--performance))
  ist zu messen, nicht zu behaupten; die Rückführung des Slice hängt daran.
- Absichtlich ortsgebundene Verweise laufen über das bestehende Ventil
  `ignore-refs` — kein neues Ventil.

## Re-Evaluierungs-Trigger

- Ein Konsument braucht Gruppen mit **nicht-Geschwister**-Orten (verschiedene
  Tiefen) — dann ist die Auflösungs-Semantik der Gruppe neu zu stellen.
- Die Laufzeit-Messung reißt die Zusage — dann ist eine Kandidaten-Eingrenzung
  (nur relative Verweise ohne `../`-Präfix?) zu prüfen, bevor die Fähigkeit
  beschnitten wird.
- Das Ventil `ignore-refs` erweist sich als zu grob für ortsgebundene Verweise
  (es nimmt referenz-weit aus, nicht orts-weise) — dann ist ein eigenes Ventil
  eine neue Entscheidung.

## Geschichte

- 2026-08-16: Proposed (`slice-095`, nach der Bestandsmessung samt Retro-Beleg).
- 2026-08-16: Accepted (Closure `welle-76`; Release v0.60.0, zwei Review-Runden,
  fail-closed-Zusage nach CI-Realfund an der git-Realität justiert — der
  Kern-Entscheid „nur wandernde Quellorte, Ziel-Identität gehört zur
  Auflösbarkeit" ist unverändert umgesetzt).
