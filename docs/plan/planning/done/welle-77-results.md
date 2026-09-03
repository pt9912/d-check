# Welle 77 — Eine chronologische Tabelle hält ihre Richtung — Closure-Notiz

**Welle:** welle-77-chronologie-ordnung
**Abschluss:** 2026-08-21
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Die Chronologie-Monotonie** ([slice-105](welle-77/slice-105-tabellen-monotonie.md)):
  siebte Bedingung von
  [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  (kein neues Modul, kein neues Kürzel — Kriterium aus
  [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md)):
  `table-order`/`table-column` vergleichen die Schlüsselspalte jeder
  zusammenhängenden Tabelle **typisiert** (ISO-Datum, Punkt-Version) und
  nicht-strikt monoton, gelesen auf den **rohen** Zellen — die erste benannte
  Ausnahme von der Abschnitts-Bereinigung. Zwei Grund-Codes
  (`section-unordered`, `section-cell-untyped`), drei Exit-2-Ränder.
  [ADR-0057](../../adr/0057-structure-tabellen-monotonie.md).
- **Die Tabellen-Lexik hat ihre Kopplung:** Trennzeilen-/Kopfzeilen-Erkennung
  und Zell-Splitting wohnen beim geteilten Tabellenzeilen-Prädikat; der
  Kopplungs-Test fährt dieselbe Eingabe durch alle drei Konsumenten
  (`targets`, `planning.waves`, `structure`) — die Form aus
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md), fällig
  mit dem dritten Konsumenten (der Vorschau-Trigger dieser Welle).
- **Im eigenen Baum scharfgeschaltet** (`.d-check.yml`, Inner-Loop): die sechs
  chronologischen Bestandstabellen — beide §7-Historien, beide
  Roadmap-Register, `version.md`, Handbuch-§11 (als einzige aufsteigend).
- **Release v0.61.0**, Digest `sha256:0e731cfc…9f98`, beide Pipelines im
  ersten Anlauf grün.

## Was hat funktioniert?

- **Der Register-Eintrag war der Entwurf.** BEO-005 trug seit welle-73 die
  drei Entscheidungen, die keine Details sind (rohe Zeilen · Zell-Adresse via
  geteilter Lexik · Richtung je Regel plus Befund statt Übersprung) — der
  Slice hat sie umgesetzt, nicht neu erfunden. Ein gut geschriebener
  Beobachtungs-Eintrag ist billiger als jede spätere Analyse.
- **Der Retro-Beleg traf die Skript-Messung exakt:** das Produkt meldet am
  Stand vor der welle-73-Heilung 27 Befunde — 14 · 6 · 7 über genau die drei
  historisch gekippten Tabellen, null über die drei übrigen, null am heutigen
  Bestand. Die naive Gegenprobe (String-Vergleich rot auf korrekt Sortiertem)
  ist als Testfall festgeschrieben, nicht als Anekdote.
- **Die Fähigkeit fing ihren ersten Fall vor dem ersten Release:** der
  Handbuch-§11-Eintrag der eigenen Release-Prep saß falsch herum (oben statt
  chronologisch unten — die BEO-005-Geste in Reinform, beim Verfasser dieses
  Slice). Vor dem Commit gedreht; ab v0.61.0 wäre er maschinell rot.

## Was ging anders als geplant?

- **Der frisch geschriebene Vertrag trug zwei Lesarten derselben Prüfung:**
  die §6-Bedingungs-Tabelle sagte „Spalten-Typ", der Fließtext
  „Vorgänger-Typ" — bei Datum–Version–Datum verschiedene Befundbilder, der
  Code lieferte das eine, der divergente Zweig war ungetestet (der
  MEDIUM-Befund des Erst-Reviews). Gepinnt auf die Paar-Lesart mit
  Anker-Reset, Kaskaden-Fall als Test. Wer denselben Vertrag als Tabelle
  **und** Fließtext formuliert, schreibt zwei Verträge.
- **Ein als unerreichbar kommentierter Fehlerpfad war erreichbar:** `\d+`
  kennt keine int-Grenze; ein Überlauf-Segment verglich still als
  kleinstmögliche Version statt zu melden. Jetzt untypisierbar (Befund), der
  Zweig getestet.
- **Sonst nichts.** Geplantes Ende 2026-08-22, geschlossen am 2026-08-21;
  eine Welle, ein Slice, zwei Review-Runden (blockierend → APPROVE ohne
  Auflagen).

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 4 — was aus dieser Welle in die Steuerung
zurückfließt.

- **Eine Bedingungs-Tabelle neben einem Bedingungs-Fließtext ist eine
  Spiegel-Beziehung im Kleinen** — dieselbe Klasse wie
  [`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten),
  nur innerhalb **eines** Dokuments und **eines** Commits: beide Stellen
  entstanden zusammen und divergierten trotzdem. Der billigste Wächter ist der
  Testfall, der genau den Fall pinnt, in dem die Lesarten auseinanderfallen.
- **„Nach dem Muster nicht erreichbar" ist eine Behauptung über das Muster,
  nicht über die Eingabe.** Ein unbegrenztes `\d+` macht jeden
  Parser-Fehlerzweig erreichbar. Wer einen Fehlerzweig als unerreichbar
  kommentiert, schreibt besser den Befund, den er im Zweifel melden würde —
  der Kommentar war falsch, der Befund ist jetzt echt.

## Beobachtungs-Register (Zeiger)

Der Lese-Schritt dieser Welle: **BEO-005 gestrichen** — die eingetretene
Instanz (die sechs eigenen chronologischen Bestandstabellen) ist mechanisiert,
die siebte `structure`-Bedingung läuft scharf in `make gates`; was die
Mechanisierung nicht deckt (Tabellen außerhalb der sechs aktivierten Regeln,
gesplittete Tabellen, fremde Repos ohne Aktivierung), steht benannt im
Streichungs-Eintrag. Das Register führt damit erstmals keine
**unverkörperte** Beobachtung mehr: BEO-001 gestrichen (welle-75),
BEO-002/003/004 verkörpert (und als Klassen ausdrücklich weiter offen),
BEO-005 gestrichen. Der Zähler läuft weiter — der nächste geschlossene Slice
schreibt wieder, und die Verkörperungen finden weiter, statt zu zählen.

## Folge-Slices

- **Keiner.** `open/` ist leer; die nächste Welle (Baseline-Migration
  v5.0.0 → v5.6.0, Reihenfolge-Entscheid des Auftraggebers vom 2026-08-21)
  steht in der Vorschau der Roadmap und schneidet ihre Slices bei der
  Eröffnung. Die INFO-Reste der Reviews sind benannte Grenzen in
  [ADR-0057](../../adr/0057-structure-tabellen-monotonie.md) mit
  Re-Evaluierungs-Triggern, keine Slices.

## Verifikation

- **Closure-Trigger erfüllt:** [slice-105](welle-77/slice-105-tabellen-monotonie.md)
  in `done/`; Retro-Beleg mit dem Produkt (27 = 14 · 6 · 7 am
  Vor-Heilungs-Stand, null heute) samt naiver Gegenprobe als Testfall; der
  Kopplungs-Test läuft; der Register-Schritt ist entschieden (gestrichen);
  Release **v0.61.0** als Minor samt Digest-Backfill; `make fullbuild` grün,
  Image-Hash `sha256:1053f729…6db5`.
- **Trigger-Audit** über die drei Artefaktklassen: keine offenen Carveouts,
  keine stehengebliebene Gate-Reifestufe (die Bedingung läuft von Anfang an
  blockierend in `make gates`), kein eingetretener Re-Evaluierungs-Trigger —
  ausdrücklich geprüft für [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md):
  die zweite Zell-Lesart des trace-Lesers ist keine „vierte Stelle, die eine
  Lexik-Frage selbst beantwortet", sondern per
  [ADR-0057](../../adr/0057-structure-tabellen-monotonie.md) eine **andere
  Frage** auf einer anderen Vertragsfläche — die Einordnung ist dort
  ausdrücklich revidierbar deklariert.
- **Zwei unabhängige Review-Runden:** die erste blockierend (1 MEDIUM,
  1 LOW, 3 INFO — alle geheilt), die zweite bestätigend APPROVE ohne
  Auflagen; 20 Negativbefund-Zeilen im Erst-Report.
