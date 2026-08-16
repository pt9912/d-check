# Welle 74 — Dieselbe Klasse, drei andere Lexiken — Closure-Notiz

**Welle:** welle-74-geteilte-lexik-raender
**Abschluss:** 2026-08-16
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Sechs Module beantworten ihre Lexik-Fragen jetzt gleich**
  ([slice-103](slice-103-geteilte-lexik-raender.md)): `citations` (Absatzgrenze),
  `versions` und `pins` (Anker-Auflösung, vollständig — Fence, Inline-Code,
  Tag-Kontext, Duplikat-Slug, Prozent-Dekodierung, Groß-/Kleinschreibung),
  `planning` (Überschrift, Marker, Block-Grenze), `vcs` (Status-Zeile,
  `immutable-when`) und `targets` (Tabellenzeilen).
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md).
- **Zwei stille Grün-Pfade in Gates geschlossen:** eine reale Core-Änderung
  einer `Accepted`-ADR passierte das Immutabilitäts-Gate mit Exit 0 ohne
  Ausgabe, wenn die Datei ihren eigenen Kopf als Beispiel zeigt; und
  `planning-drift` entfiel still, wenn der Aktiv-Block an einer eingerückten H2
  oder einer H1 endete.
- **Eine benannte Grenze statt eines zweiten Wächters:** die Abschnitts-Maske
  von `vcs` läuft auf git-Blobs, die kein scannendes Modul sieht — mit einem
  Re-Evaluierungs-Trigger, der sich beobachten lässt.
- **Release v0.58.0**, Digest `sha256:a4a9275b…b17a`, Pipeline im ersten Anlauf
  grün.

## Was hat funktioniert?

- **Erst messen, dann entscheiden — und die Messung mit dem Produkt fahren.**
  Eine naive Paritätszählung meldete zwei verdächtige Revisions-Blobs; das
  Produkt über dieselben Blobs sagte null. Die Nachrechnung irrte, nicht der
  Bestand.
- **Drei Review-Runden, jede blockierend, jede an einer neuen Stelle.** Genau
  das Muster aus dem Vorgänger-Slice — und die dritte Runde war die
  entscheidende: sie hat alle zwanzig Module plus die App-Schicht geprüft und
  drei gefunden, die in keinem der beiden ersten Reports vorkamen.
- **Der Auftraggeber hat zweimal in laufende Arbeit hineingefragt** und beide
  Male einen Defekt getroffen: einmal die noch nicht ausgeführte Bewegung nach
  `done/`, einmal drei Module, deren Handbuch-Beschreibungen im Release-Prep
  fehlten.

## Was ging anders als geplant?

- **Die Klasse war zweimal „geschlossen" und war es nicht.** Nach Runde eins
  standen zwei Stellen fest, nach Runde zwei vier, nach Runde drei sechs. Jede
  Zwischenbilanz war als vollständig formuliert.
- **Die Reparatur selbst war zweimal halb.** Die Anker-Frage wurde erst zur
  Fence-Hälfte vereinheitlicht, dann zur HTML-Hälfte, und erst in der dritten
  Runde auch in der Slug-Hälfte. „Dieselbe Grundmenge" ist noch keine
  „dieselbe Antwort".
- **Die Richtungs-Aussage war dreimal geschlossen formuliert und dreimal
  unvollständig.** Ursache beim dritten Mal war ein Werkzeug-Fehler von mir: ein
  Batch-Editor schreibt erst, wenn alle Ersetzungen greifen — eine gescheiterte
  Ersetzung hat die Korrektur der anderen still verschluckt.
- **Die Verkörperung hat den Defekt geerbt, den sie ablösen sollte.** Der
  Kopplungs-Test fuhr zwei der drei Konsumenten.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 4 — was aus dieser Welle in die Steuerung
zurückfließt.

- **Eine Aufzählung belegt keinen Klassen-Abschluss.** Dreimal hat eine Liste,
  die „vollständig" hieß, Stellen nicht gekannt. Was trägt, ist entweder eine
  **Kopplung** (ein Test, der alle Konsumenten dieselbe Eingabe beantworten
  lässt) oder eine **erschöpfende Prüfung mit Negativbefund je Kandidat**.
- **Eine Richtungs-Aufzählung gehört offen formuliert.** „Und weniger an keiner
  Stelle" ist eine universelle Aussage; belegbar ist immer nur „an diesen
  gemessenen Stellen".
- **Eine geänderte Zusage zählt wie ein neues Feature** — sie steht nicht in
  einem neuen Abschnitt, sondern mitten im alten. Jetzt ein Punkt der
  [Release-Prep-Checkliste](../../../user/releasing.md#release-prep-vor-dem-tag);
  im selben Release trotzdem einmal gerissen und vom Auftraggeber gefunden.
- **Eine Mutations-Gegenprobe wird am Exit-Code geprüft, nicht am
  grep-Muster.** Ein Rückbau, der einen Compile-Fehler erzeugt, sieht sonst aus
  wie „kein Test fängt das". Zweimal passiert, beide Male in derselben Sitzung.

## Beobachtungs-Register (Zeiger)

Der Lese-Schritt dieser Welle: das Register führt **BEO-001** (2×),
**BEO-002** (3×, verkörpert), **BEO-003** (jetzt 3×, **verkörpert**),
**BEO-004** (3×, verkörpert) und **BEO-005** (1×).

- **BEO-003** ist dreimal eingetreten und in dieser Welle **verkörpert** — nicht
  als vierte Liste, sondern als **Kopplungs-Test**
  (`TestAnkerFrageHatEineAntwort`): dieselbe Eingabe durch alle drei
  Konsumenten derselben Frage, Fehlschlag bei abweichender Antwort. Er prüft die
  **Übereinstimmung**, nicht den Aufruf. Offen bleibt die Klasse dort, wo noch
  kein Kopplungs-Test existiert — die Absatz-, Überschriften- und
  Tabellen-Achse tragen bis heute nur Einzel-Assertionen.
- **BEO-005** ist unberührt; die Tabellen-Lexik, die sie brauchte, ist in dieser
  Welle nur **entdriftet**, nicht ausgebaut.
- **BEO-001**, **BEO-002**, **BEO-004** unberührt.

## Folge-Slices

- [slice-095](../open/slice-095-links-resolve-from.md) und
  [slice-102](../done/slice-102-wellen-lifecycle-invariante.md) lagen bei
  dieser Closure unverändert in `open/` (slice-102 ist inzwischen in Arbeit —
  der Link zeigt auf seinen heutigen Ort, die Aussage gilt dem Closure-Zeitpunkt).
- **Kein** Folge-Slice aus Abnahme-Punkt 2: die Messung hat den Schnitt auf
  einen Slice gedreht, und die drei Nachträge der Review-Runden gehörten
  derselben Klasse an.

## Verifikation

- **Closure-Trigger erfüllt:** [slice-103](slice-103-geteilte-lexik-raender.md)
  in `done/`; die Bestandsmessung lag **vor** dem Schnitt-Entscheid; **BEO-003
  ist entschieden** (auf 3, verkörpert); Release **v0.58.0** samt
  Digest-Backfill; `make fullbuild` grün, Image-Hash
  `sha256:cb4e1208…5312`.
- **Trigger-Audit** über die drei Artefaktklassen: keine offenen Carveouts,
  keine stehengebliebene Gate-Reifestufe, und von den
  Re-Evaluierungs-Triggern der Welle-ADR ist einer **eingetreten** — „eine
  vierte Stelle beantwortet eine Lexik-Frage selbst" — und in derselben Welle
  beantwortet: die Antwort ist der Kopplungs-Test, nicht ein weiterer
  Listen-Anlauf.
- **Drei unabhängige Review-Runden**, alle drei blockierend, alle drei mit
  eigenen Messungen und Mutations-Gegenproben; die dritte hat zusätzlich eine
  **Negativbefund-Zeile je Modul** geschrieben — die erste vollständige
  Aufzählung dieser Klasse im Repo.
