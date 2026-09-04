# Welle 75 — Die Roadmap-Aussage gegen das Verzeichnis — Closure-Notiz

**Welle:** welle-75-wellen-register
**Abschluss:** 2026-08-16
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Die Wellen-Invariante** ([slice-102](wellenlos/slice-102-wellen-lifecycle-invariante.md)):
  dritte Fähigkeit des Moduls `planning`, vier Aussagen mit je eigenem
  Grund-Code (`wave-drift` · `wave-preview-exists` · `wave-results-missing` ·
  `wave-unregistered`), opt-in über `planning.waves.dir`.
  [ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md).
- **Der Aktiv-Status wird einmal bestimmt** (`planningActiveStatus`), die
  Tabellenzeilen-Erkennung ist ein geteiltes Prädikat (`tableRowLine`) — beide
  Anschlüsse an [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  standen vor der ersten Zeile Code im Vertrag.
- **Im eigenen Baum scharfgeschaltet** (`.d-check.yml`): die Fähigkeit läuft in
  `make gates`, der Bestand ist grün.
- **Release v0.59.0**, Digest `sha256:b1cc65df…b493`, Pipeline im ersten Anlauf
  grün.

## Was hat funktioniert?

- **Die Messung hat den Entwurf zweimal gedreht, bevor Code existierte:** das
  verpflichtende Artefakt einer geschlossenen Welle ist die **Ergebnisnotiz**
  (gegen das Plan-Dokument geprüft hätte die Aussage 19-mal falsch gemeldet),
  und die Vorschau-Aussage liest die **erste Spalte** statt der Zeile (die
  Trigger-Spalte darf andere Wellen nennen — die eigene Roadmap tat es am
  selben Tag).
- **Der Beleg lief mit dem Produkt über den realen Bestand:** eigener Baum 0,
  Schwester-Repo elf robuste `wave-results-missing`.
- **Zwei Review-Runden, beide blockierend, beide produktiv:** die erste fand den
  HIGH in der Motivations-Richtung des Slice, die zweite wandte die
  Vier-Codes-Begründung der eigenen ADR gegen die Heilung.

## Was ging anders als geplant?

- **Der HIGH der Erst-Review traf die eigene Motivation:** ein unlesbares
  `waves.dir` schaltete die Fähigkeit im Ruhe-Zustand still ab — mit einem
  Pfad-Tippfehler wäre genau die zweimal real eingetretene Verletzung dauerhaft
  unsichtbar gewesen, die der Slice beheben sollte.
- **Meine eigene Messaussage war falsch:** der zwölfte Schwester-Repo-Befund
  war ein Artefakt des Default-Markers meiner Probe-Konfiguration, keine
  Bestands-Verletzung. Der Reviewer fand es, die Nachmessung mit
  konsument-gerechtem Marker bestätigte elf.
- **Die Heilung erzeugte einen neuen Defekt derselben Bauart, die die ADR
  ausschließt:** `wave-drift` wurde zum Sammelcode für drei Bedeutungen, und
  zwei davon kollabierten in der Befund-Deduplikation zu einem Befund — die
  Ziele trennen die Bedeutungen jetzt.
- **Eine eigene Gegenprobe war schlecht konstruiert** (Kennung in der
  Vorschau-Kopfzeile, wo ohne Datei nichts meldet) und blieb grün, obwohl der
  Rückbau das Produktverhalten kippte — erst die Kennung in der
  Abschluss-Kopfzeile machte sie rot.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 4 — was aus dieser Welle in die Steuerung
zurückfließt.

- **Eine Probe-Konfiguration ist Teil der Messung.** Wer einen fremden Bestand
  mit den eigenen Defaults misst, misst auch die Differenz der Konventionen —
  und liest sie leicht als Befund. Konsument-gerechte Schlüssel (hier: der
  Ruhe-Marker) gehören in die Probe, und ein überraschender Befund verdient
  eine Kontroll-Messung mit angepasster Config.
- **Ein Sammel-Grund-Code braucht getrennte Befund-Ziele.** Wer mehrere
  Bedeutungen unter einem Code meldet, muss die Deduplikations-Achse (Datei,
  Zeile, Regel, Ziel, Grund) je Bedeutung eindeutig halten — sonst verschluckt
  sie eine Reparatur.
- **Fail-closed heißt: auch im konsistenten Zustand.** Der stille Pfad lag
  genau dort, wo die leere Menge die Erwartung erfüllte. Die Frage an jeden
  neuen Rand ist nicht „meldet er im Fehlerfall?", sondern „gibt es einen
  Zustand, in dem der Fehler wie Konsistenz aussieht?".
- **Eine Gegenprobe gehört dorthin, wo der Rückbau das Ergebnis kippt** — nicht
  dorthin, wo er zufällig folgenlos ist. Die Kopfzeilen-Probe war grün, weil
  das Fixture den Effekt nicht erreichte.

## Beobachtungs-Register (Zeiger)

Der Lese-Schritt dieser Welle: das Register führt **BEO-001** (2×),
**BEO-002** (3×, verkörpert), **BEO-003** (3×, verkörpert), **BEO-004** (3×,
verkörpert), **BEO-005** (1×).

- **BEO-001 ist GESCHLOSSEN.** Die Beobachtung war: ein Datei-Register und
  seine Autoritäts-Tabelle driften unbemerkt, und kein Gate deckt die Richtung
  „Artefakt ⇒ registriert". Für das **Wellen-Register** — die Instanz, an der
  die Beobachtung zweimal eingetreten ist — prüft jetzt `wave-unregistered`
  genau diese Richtung, im eigenen Baum scharf in `make gates`. Die beiden
  anderen genannten Bauformen (ADR-Index, Konventionsspeicher-Index) bleiben
  ungeprüft; sie sind im Register-Eintrag als Rest ausgewiesen und wandern in
  die Streichungs-Begründung — der `registry`-Modul-Vorschlag bleibt dort als
  Option benannt, falls die Klasse an einer dritten Instanz eintritt.
- **BEO-005** unberührt: die Tabellenzeilen-Lexik hat jetzt ihren zweiten
  Konsumenten mit geteiltem Prädikat; die Ordnungs-Bedingung wartet auf ihre
  eigene Welle, und beim dritten Konsumenten ist der Kopplungs-Test fällig.

## Folge-Slices

- [slice-095](wellenlos/slice-095-links-resolve-from.md) lag bei dieser
  Closure unverändert in `open/` (inzwischen umgesetzt — der Link zeigt auf
  seinen heutigen Ort, die Aussage gilt dem Closure-Zeitpunkt); die
  Chronologie-Ordnung (**BEO-005**) steht als geplante Welle in der Vorschau.
- **Kein** Folge-Slice aus den Reviews: F-13 (doppelte Roadmap-Lesung) ist eine
  benannte Ineffizienz ohne Divergenz-Risiko, seit die Wellen-Fähigkeit bei
  unlesbarer Roadmap schweigt.

## Verifikation

- **Closure-Trigger erfüllt:** [slice-102](wellenlos/slice-102-wellen-lifecycle-invariante.md)
  in `done/`; die Bestandsmessung lag **vor** jeder Scharfschaltung (je Aussage,
  §3a); **BEO-001 ist entschieden** (geschlossen für die eingetretene Instanz);
  Release **v0.59.0** samt Digest-Backfill; `make fullbuild` grün, Image-Hash
  `sha256:c8449147…beac`.
- **Trigger-Audit** über die drei Artefaktklassen: keine offenen Carveouts,
  keine stehengebliebene Gate-Reifestufe (die neue Fähigkeit ist von Anfang an
  blockierend in `make gates`), und keiner der drei Re-Evaluierungs-Trigger von
  [ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md) ist
  eingetreten.
- **Zwei unabhängige Review-Runden** mit eigenen Messungen, Mutations-Gegenproben
  über Dateikopien und unabhängiger Replikation der Elf-robust-Zahl.
