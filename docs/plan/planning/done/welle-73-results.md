# Welle 73 — Das 20. Regelmodul und die verkörperte Spiegel-Regel — Closure-Notiz

**Welle:** welle-73-structure-umsetzung
**Abschluss:** 2026-08-15
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Modul `structure`** ([slice-099](welle-69/slice-099-structure-modul.md)), das 20.
  Regelmodul: Struktur-Invarianten **innerhalb** eines Dokuments, je Regel eine
  Dokumentklasse über eigene Globs, ein Abschnitt und bis zu sechs Bedingungen
  mit je eigenem Grund-Code.
  [ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md).
- **Die Abschnitts-Mechanik ist geteilt, nicht kopiert.** Die Closure-Fähigkeit
  des Moduls `planning` ist ein **Preset** derselben Semantik; ein
  Kopplungs-Test fährt dieselbe Eingabe durch beide Oberflächen.
- **Eine seit 0.51.0 offene Zusage ist eingelöst:** `closure-note-ambiguous`
  war in Anforderung, Algorithmus, §4-Tabelle und `--doctor`-Klartext deklariert
  und wurde **nie erzeugt**. Neun neue Grund-Codes insgesamt.
- **`--print-mk` trägt ein zwölftes Target** (`doc-structure`), und die
  Config-Vorlage ein `structure`-Gerüst.
- **Drei eigene Regeln aktiviert** — gemessen, nicht geraten: von acht
  plausiblen Kandidaten haben drei den Bestand grün gefunden.
- **Release v0.57.0**, Digest `sha256:e9d52946…91d7`, Pipeline im ersten Anlauf
  grün.

## Was hat funktioniert?

- **Die Regel stand vor dem Slice.** [`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
  ist die Verkörperung von **BEO-002** und wurde **vor** der Umsetzung
  geschrieben, damit sie an diesem Slice zum ersten Mal wirkt statt an ihm zum
  vierten Mal zu fehlen. Sie hat gewirkt: die Ziffer in der
  Out-of-Scope-Zeile, der Intro-Satz beider READMEs und `operations.md` standen
  auf der Liste und wären sonst stehengeblieben.
- **Erst messen, dann entscheiden.** Die Regel-Auswahl kam aus einer Messung des
  eigenen Bestands; die plausibelste Kandidatin („abgeschlossener Slice ohne
  offene Task-Boxen") meldete 32-mal und war **jedes Mal falsch** — die Welle
  löst den Punkt ein, nicht der Slice.
- **Zwei Review-Runden, beide blockierend.** Die zweite hat die Erreichbarkeit
  **methodisch** statt stichprobenartig geprüft: alle 44 `AllReasons`-Konstanten
  gegen ihre Emissionsstelle. Kein zweiter toter Grund-Code.

## Was ging anders als geplant?

- **Die Spiegel-Liste war selbst lückenhaft** — vier Lücken, die der Review
  fand, und eine fünfte, die die Re-Review fand. Die fünfte ist die
  interessanteste: sie lag in **genau der Datei**, die ohnehin bearbeitet wurde,
  und war schon vorher veraltet. Ein Datei-`grep` hätte sie nie gefunden.
- **Ein Grund-Code war vollständig deklariert und tot.** Der
  `AllReasons`-↔-§4-Lockstep-Test prüft **Katalog-Abdeckung**, nicht
  **Erreichbarkeit** — und mein eigener Test zementierte das alte Verhalten mit
  der Begründung „bis der Grund-Code existiert", während er im selben Diff
  entstand.
- **Jeder Datumsstempel des Slice war falsch.** Die Welle lief über zwei
  Kalendertage; zwölf Stellen trugen den Tag des vorigen Releases, weil sie von
  der Zeile darüber abgeschrieben waren.
- **Zwei Historie-Tabellen waren nicht chronologisch** (Nutzer-Befund):
  Spezifikation und Lastenheft tragen je einen alten, unten angehängten Block
  und einen neueren, oben eingefügten — die Richtung kippte irgendwann still.
- **Das Closure-Log der Roadmap fehlte drei Wellen lang.** welle-70, welle-71
  und welle-72 haben ihre Zeile in §Abgeschlossene Wellen nie bekommen, obwohl
  alle drei ihre `done/welle-NN-results.md` geschrieben haben.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 4 — was aus dieser Welle in die Steuerung
zurückfließt.

- **Wo eine Aufzählung eine Menge spiegelt, gehört sie an ihre Quelle gebunden.**
  Eine Liste im Konventionsspeicher ist der Notbehelf für alles, was sich nicht
  binden lässt — die Modul-Enumeration der Config-Vorlage ließ sich binden und
  hängt jetzt an der Registry.
- **Ein Katalog-Test prüft Abdeckung, nicht Erreichbarkeit.** Ein Grund-Code
  kann vollständig deklariert, dokumentiert und tot sein.
- **Das Datum kommt aus dem Kalender, nicht aus der Zeile darüber** — jetzt ein
  eigener Punkt der [Release-Prep-Checkliste](../../../user/releasing.md#release-prep-vor-dem-tag).
- **Append-Logs sind gate-unsichtbar.** Weder ihre Vollständigkeit noch ihre
  Richtung prüft heute etwas; beides ist an einem Tag an drei Stellen
  aufgefallen.

## Beobachtungs-Register (Zeiger)

Der Lese-Schritt dieser Welle: das Register führt **BEO-001** (jetzt 2×),
**BEO-002** (3×, verkörpert), **BEO-003** (jetzt 2×), **BEO-004** (3×,
verkörpert) und neu **BEO-005**.

- **BEO-001** ist **eingetreten — und zwar genau dort, wo die Beobachtung es
  vorausgesagt hatte:** „dieselbe Bauform tragen der Konventionsspeicher-Index
  und das **Wellen-Register der Roadmap**; beide sind heute konsistent, aber
  ebenso ungeprüft". Drei Wellen-Closures ohne ihre Zeile. Der Zähler steigt auf
  2, und der `registry`-Modul-Vorschlag ist damit gemessen statt hypothetisch.
- **BEO-002** ist als [`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
  verkörpert und bleibt **offen**: die Regel ist eine für Menschen, sie kann
  weiter verfehlt werden — und sie war beim ersten Einsatz selbst lückenhaft.
  Kein Zählschritt, sondern eine Schärfung (der Spiegel ist die **Stelle**,
  nicht die Datei).
- **BEO-003** ist **eingetreten**: dieselbe geteilte Mechanik trägt zwei
  Wort-Begriffe (Marke unicode-weit, Floskel ASCII), und nur einer stand im
  Vertrag. Zähler auf 2; die Klasse ist damit über zwei Slices belegt und
  [slice-103](welle-74/slice-103-geteilte-lexik-raender.md) bleibt geschnitten.
- **BEO-004** stand seit welle-70 bei 3 und ist in dieser Closure **verkörpert**
  worden — nicht in der vorgeschlagenen Form (ein Slice-Template gibt es in
  diesem Repo nicht, [`MR-018`](../../../../harness/conventions.md#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates)), sondern als Anker im Reviewer-Skill: dort ist
  die Klasse dreimal gefunden worden, dort greift sie.
- **BEO-005** ist neu: eine chronologische Tabelle kippt still ihre Richtung.

## Folge-Slices

- [slice-103](welle-74/slice-103-geteilte-lexik-raender.md) — dieselbe
  Drift-Klasse in anderen Lexiken; durch BEO-003 zum zweiten Mal belegt.
- [slice-095](wellenlos/slice-095-links-resolve-from.md) und
  [slice-102](wellenlos/slice-102-wellen-lifecycle-invariante.md) lagen bei
  dieser Closure unverändert in `open/` (beide inzwischen umgesetzt —
  die Links zeigen auf ihre heutigen Orte, die Aussage gilt dem Closure-Zeitpunkt).
- **Kein** Folge-Slice für BEO-001: der `registry`-Vorschlag steht im Register
  und bleibt dort, bis er eine Welle bekommt.

## Verifikation

- **Closure-Trigger erfüllt:** [slice-099](welle-69/slice-099-structure-modul.md) in
  `done/`; [`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten) hat am Slice gewirkt und ist mit Bilanz abgehakt; Release
  **v0.57.0** samt Digest-Backfill; `make fullbuild` grün, Image-Hash
  `sha256:a908533a…1e84b`.
- **Trigger-Audit** über die drei Artefaktklassen: keine offenen Carveouts,
  keine stehengebliebene Gate-Reifestufe (`structure` ist von Anfang an
  blockierend im Closure-Gate), und keiner der drei Re-Evaluierungs-Trigger von
  [ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md) ist
  eingetreten — insbesondere braucht keine dritte Oberfläche dieselbe
  Struktur-Semantik.
- Das Closure-Log der Roadmap ist in derselben Bewegung **nachgetragen**: die
  drei fehlenden Wellen-Zeilen stehen, und die Drift-Tabelle ist chronologisch.
