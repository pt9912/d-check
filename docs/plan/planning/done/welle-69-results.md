# Welle 69 — Struktur-Invarianten: Schnitt und Ablöse-Pfad — Closure-Notiz

**Welle:** welle-69-structure-schnitt
**Abschluss:** 2026-08-09
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Der Modul-Schnitt ist entschieden und vertraglich festgehalten:** neue
  Anforderung
  [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  (Lastenheft 0.51.0) mit Algorithmus-Sektion, Config-Schema und
  Akzeptanzkriterien je Fall.
- **[ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md) `Accepted`**
  mit neun Entscheidungen und zwölf verglichenen Alternativen — darunter die
  tragende: es wird **nichts superseded**, weil Grund-Codes stabil zugesagt sind;
  die Closure-Fähigkeit wird stattdessen als **Preset** derselben Semantik
  spezifiziert.
- **§1 des Lastenhefts benennt die zweite Frage** („ist dieses Dokument selbst
  richtig gebaut?") als eigene Kategorie. Sie lief mit `spans`/`hostpaths` längst
  mit, war aber nie ausgesprochen — genau die Klasse, die den Antrag ausgelöst hat.
- **Zwei Folge-Slices geschnitten und angelegt:**
  [slice-099](../open/slice-099-structure-modul.md) (Implementierung) und
  [slice-101](../done/slice-101-fence-unbalanciert.md).
- **Die drei liegenden Closure-Slices haben ihre Entscheidung:** alle drei
  bleiben eigenständig — das war das *Mehr* dieser Welle gegenüber der Slice-DoD.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- **Drei unabhängige Reviews statt eines.** Vertrag und Umsetzbarkeit fanden
  **verschiedene** Klassen; die bestätigende Re-Review fand, was die Heilung
  liegen ließ. Kein einzelner Lauf hätte alle drei Ebenen erreicht.
- **Die ADR bis zuletzt auf `Proposed` zu lassen.** Zwei Korrekturrunden liefen
  über sie; als `Accepted` wären vier Fehler eingefroren gewesen — darunter zwei
  zurückgezogene Entscheidungen in der Index-Zeile.
- **Je Prüfung messen statt je Skript.** Ein Skript ist eine Datei, keine
  Aussage. Erst die Zerlegung in elf Prüfungen machte sichtbar, dass **eine**
  Zeile den Modul-Schnitt entscheidet — und später, dass genau diese Zeile unter
  dem ersten Vertrag unprüfbar gewesen wäre.
- **Ein ausgelieferter Defekt wurde als eigener Slice ausgelagert**, statt ihn in
  die laufende Arbeit zu ziehen: der stille Grün-Pfad im Fence-Automaten gehört
  nicht in einen Schnitt-Slice.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile möglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **Der erste Vertrag konnte seine eigene Begründung nicht ausdrücken.** Der
  Schnitt wurde an einer Messzeile entschieden; der Vertrag machte genau diesen
  Fall unprüfbar. Konsequenz: ein Kardinalitäts-Modus, und die Lehre, ein
  begründendes Beispiel gegen den **fertigen** Vertrag zu halten, nicht gegen die
  Absicht.
- **Zwei Reviews reichten nicht.** Die Heilung der ersten Runde erzeugte zwölf
  neue Befunde — alle desselben Musters: Körper korrigiert, Ränder stehen
  gelassen. Konsequenz: [BEO-002](../observations.md) im Register.
- **Ich habe in eine abgeschlossene Beweislage eingegriffen.** Um das Gate grün
  zu halten, hatte ich Links in zwei fertigen Review-Reports umgebogen — mit dem
  Ergebnis, dass ein Label nicht mehr zu seinem Ziel passte. Konsequenz: zurück
  auf den Ursprungstext, Auflösung über das Referenz-Ventil mit Quell-Skopus, das
  dieses Repo für eingefrorene Verweise besitzt.
- **Aus zwei Folge-Slices wurde einer.** Die geplante Trennung in Kern und
  Marken/Zählung scheitert an der Release-Grenze: das veröffentlichte Schema
  führt alle Schlüssel, und die Dekodierung ist strikt.
- **Das geplante Wellen-Ende (2026-08-12) wurde unterschritten** — die Welle
  schloss am Starttag. Kein Closure-Kriterium war betroffen; die Schätzung war
  schlicht zu vorsichtig für eine reine Schnitt-Welle.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`).

- — keine — . [BEO-002](../observations.md) steht bei 2×, [BEO-001](../observations.md)
  bei 1×. Beide bleiben offen und warten; der Lese-Schritt hat sie gesehen.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Zähler steht in [`observations.md`](../observations.md).
Was in dieser Welle **3×** erreicht hat, steht oben unter
*Steering-Loop-Einträge*.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren.

- [slice-101](../done/slice-101-fence-unbalanciert.md) — der ausgelieferte stille
  Grün-Pfad im Fence-Automaten; **bindende** Vorbedingung für 099.
- [slice-099](../open/slice-099-structure-modul.md) — die Implementierung des
  Moduls samt Preset-Kopplung und neun Grund-Codes.
- [slice-094](../open/slice-094-closure-zaehl-paritaet.md),
  [slice-097](../done/slice-097-closure-glob-entkopplung.md),
  [slice-098](../done/slice-098-closure-note-placeholder.md) — bleiben
  eigenständig, Zuschnitt durch diese Welle entschieden.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Der einzige Slice dieser Welle (096) liegt in `done/`; die Entscheidung über
  094/097/098 ist im Slice festgehalten.
- `make fullbuild` grün — Image-Hash
  `sha256:c42dfc4ba1f8388afd8c595848e6d9a4ac9a1adeef13cb8bee0ab634d0272b11`.
- `make gates` grün: 344 Datei(en) / 0 Befund(e), Coverage 94,30 % (Schwelle
  93 %), keine offenen Carveouts.
- `make verify-closure-notes` grün: 315 Datei(en) / 0 Befund(e) — einschließlich
  der Closure-Notiz dieses Slice.
- `make completeness-check`: 48 Anforderung(en), 0 Waise(n) — die neue
  Anforderung ist von einem Slice referenziert.
- `make bench`: Median 739 ms (Kriterium < 5000 ms).
- **Trigger-Audit** (Schritt 2, alle drei Artefaktklassen): keine offenen
  Carveouts; keine stehengebliebene Gate-Reifestufe (Coverage-Ist über Soll);
  kein ADR-Re-Evaluierungs-Trigger eingetreten — geprüft insbesondere
  [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md), dessen
  Schnitt-Kriterium diese Welle **angewandt** und damit bestätigt hat, und
  [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md), dessen
  Re-Evaluierungs-Trigger „ein zweites Modul braucht dieselbe Closure-Achse"
  **eingetreten** ist und mit der Preset-Kopplung beantwortet wurde — ohne
  Supersede, daher als Verfeinerung in
  [ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md) statt als
  Folge-ADR mit `supersedes`.
- **Kein Release** in dieser Welle: sie liefert Schnitt und Vertrag, kein Produkt.
