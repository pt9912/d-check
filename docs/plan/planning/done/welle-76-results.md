# Welle 76 — Ein Verweis löst von jedem Lifecycle-Ort auf — Closure-Notiz

**Welle:** welle-76-ortsfeste-verweise
**Abschluss:** 2026-08-16
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **`links.resolve-from`** ([slice-095](wellenlos/slice-095-links-resolve-from.md)):
  Erweiterung von
  [`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
  (kein neues Modul, kein neues Kürzel — Kriterium aus
  [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md)), opt-in über Gruppen wandernder Verzeichnisse. Eine Datei
  in einem `dirs`-Verzeichnis muss jedes relative Ziel von **jedem** Ort ihrer
  Gruppe auflösen, überall auf **dasselbe** Ziel; beide Fehlarten tragen den
  neuen Grund-Code `link-position-dependent`.
  [ADR-0056](../../adr/0056-resolve-from-wandernde-quellorte.md).
- **Im eigenen Baum scharfgeschaltet** (`.d-check.yml`): die Gruppe
  `open/`/`next/`/`in-progress/` mit Ruheort `done/` als `fixed-dirs` läuft in
  `make gates` — die Closure dieser Welle war ihr erster Lauf unter scharfer
  Prüfung.
- **Release v0.60.0**, Digest `sha256:5892a87b…d3f9`, beide Pipelines im
  ersten Anlauf grün.

## Was hat funktioniert?

- **Die Messung hat den Entwurf getragen, bevor Code existierte:** über alle
  vier Lifecycle-Verzeichnisse gerechnet wären 108 Verweise „positionsabhängig"
  gewesen — eingeschränkt auf die wandernden: 0 von 79. Der Zuschnitt „nur
  `dirs`-Dateien sind Quellen" stand darum im Vertrag, nicht im Nachtrag.
- **Der Retro-Beleg lief mit dem Produkt und wurde ehrlich gezählt:** 19
  Befunde am Stand vor der welle-69-Eröffnung — 15 davon Teil des realen
  19-Link-Bruchs; dessen übrige vier waren **Ziel**-Wanderungen, eine Klasse,
  die die Fähigkeit strukturell nicht deckt. Die Zahlengleichheit ist Zufall,
  und genau das steht so in ADR, Lastenheft und CHANGELOG.
- **Zwei Review-Runden mit klarer Arbeitsteilung:** die erste blockierend
  (15 Befunde — darunter die Ist-Ort-Vorbedingung gegen Doppelbefunde und der
  stille Quellen-Ausfall bei Tippfehler-Orten), die zweite APPROVE mit reinen
  Text-Auflagen, alle vor dem Release geheilt (Lastenheft 0.60.2).

## Was ging anders als geplant?

- **Der erste fail-closed-Zuschnitt fiel in der CI, nicht lokal:** er meldete
  auf jedem frischen Klon das legitim **geleerte** `open/`-Verzeichnis — git
  überträgt leere Verzeichnisse nicht, lokal existierte es noch. Ein einzelner
  fehlender Ort ist von einem Tippfehler nicht unterscheidbar; gemeldet wird
  jetzt die Gruppe ohne einen einzigen existierenden Ort und der Ort, der als
  Datei existiert. Die Rest-Grenze steht im Vertrag, der Klon-Fall ist als Test
  nachgestellt.
- **Ein „zeichengenau"-Überclaim überlebte zwei Heilungen:** die Aussage, der
  Retro-Lauf treffe „exakt die 19 des realen Bruchs", war nach der ehrlichen
  15/19-Zählung im Lastenheft korrigiert — stand aber weiter im CHANGELOG und
  im Config-Kommentar. Erst die Re-Review-Auflage N-2 zog die letzten beiden
  Spiegel nach.
- **Ein abbrechender Batch-Editor verschluckte einen ganzen Vertragsblock:**
  beim 0.60.0-Schnitt des Lastenhefts landeten Historie und Akzeptanzkriterien,
  die Anforderungs-**Beschreibung** selbst nicht — der Editor schrieb erst nach
  der letzten gelungenen Ersetzung, und eine misslungene warf alles davor weg.
  Gefunden über eine ins Leere greifende Assertion, nachgeliefert in 0.60.1;
  seither wird je Ersetzung geschrieben.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 4 — was aus dieser Welle in die Steuerung
zurückfließt.

- **Eine Fähigkeit, die Verzeichnis-Existenz zusagt, muss gegen einen frischen
  Klon entworfen werden.** git überträgt leere Verzeichnisse nicht — ein
  Zustand, der lokal konsistent aussieht, existiert auf jedem Klon anders. Der
  richtige fail-closed-Schnitt ist der, der auf **beiden** Ständen dieselbe
  Antwort gibt; `git archive HEAD` ist die billige Gegenprobe.
- **Ein Werkzeug, das erst am Ende schreibt, macht Teilerfolge unsichtbar.**
  Der Batch-Editor hat angewandte Ersetzungen mit der fehlgeschlagenen
  verworfen — das Ergebnis sah aus wie „nichts passiert", war aber „die Hälfte
  passiert". Schreiben je Ersetzung, und jede Ersetzung mit eigener Assertion.
- **Ein korrigierter Claim braucht seine Spiegel-Liste** — dieselbe Bewegung
  wie [`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
  bei einer Semantik-Änderung: die 15/19-Korrektur war im
  Vertrag angekommen und in CHANGELOG und Config-Kommentar zweimal nicht.

## Beobachtungs-Register (Zeiger)

Der Lese-Schritt dieser Welle: das Register führt **BEO-002** (3×, verkörpert),
**BEO-003** (3×, verkörpert), **BEO-004** (3×, verkörpert), **BEO-005** (1×);
BEO-001 ist gestrichen.

- **BEO-004 hat als Verkörperung gewirkt:** die Re-Review fand mit dem
  Reviewer-Anker („welche Eingaben liest dieses Modul, die es nicht scannt?")
  die Abdeckungs-Stille der Gruppen-Orte — Orte außerhalb des wirksamen
  Scan-Bereichs sind still keine Quellen. Als benannte Grenze in den Vertrag
  übernommen (N-1); der Zähler bleibt bei 3, die Verkörperung zählt nicht, sie
  findet.
- **BEO-002 ist in verwandter Form aufgetreten, ohne den Zähler zu bewegen:**
  der stehengebliebene Überclaim in CHANGELOG und Config-Kommentar ist die
  Spiegel-Klasse von [`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten),
  angewandt auf eine **Korrektur** statt eine
  Semantik-Änderung. Von der Re-Review gefunden (N-2), im Steering-Loop-Eintrag
  oben festgehalten.
- **BEO-005** unberührt: die geplante Chronologie-Welle steht in der Vorschau;
  diese Welle hat keine chronologische Tabelle angefasst, deren Richtung
  kippen konnte (die neuen Historien-Zeilen folgen der sortierten Ordnung).

## Folge-Slices

- **Keiner.** Der Backlog unter `open/` ist leer; die Chronologie-Ordnung
  (**BEO-005**) steht als geplante Welle in der Vorschau der Roadmap. Die
  LOW-/INFO-Reste der Reviews (F-10 Laufzeit-Streuung dokumentieren, F-14/F-15
  redaktionell, N-7 Meldungs-Wortlaut) sind bewusst nicht slice-würdig: keiner
  ändert Verhalten oder Vertrag.

## Verifikation

- **Closure-Trigger erfüllt:** [slice-095](wellenlos/slice-095-links-resolve-from.md) in
  `done/`; der Realdatenbeleg liegt vor (Retro-Lauf mit dem Produkt, die
  Quellort-Hälfte des Bruchs; Ziel-Wanderung als Grenze in
  [ADR-0056](../../adr/0056-resolve-from-wandernde-quellorte.md)); die
  Laufzeit-Zusage gemessen (im Rauschen,
  [`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance));
  Release **v0.60.0** als
  Minor samt Digest-Backfill; `make fullbuild` grün, Image-Hash
  `sha256:59afc1fb…89d4`.
- **Trigger-Audit** über die drei Artefaktklassen: keine offenen Carveouts,
  keine stehengebliebene Gate-Reifestufe (die Fähigkeit ist von Anfang an
  blockierend in `make gates`), und keiner der drei Re-Evaluierungs-Trigger von
  [ADR-0056](../../adr/0056-resolve-from-wandernde-quellorte.md) ist
  eingetreten.
- **Zwei unabhängige Review-Runden** mit eigenen Messungen und
  Mutations-Gegenproben über Dateikopien; die CI-Justierung ist beidseitig
  getestet (Klon-Fall grün, Tippfehler-Fall rot).
