# Slice slice-069: Trace-Handbuch — Definitionssyntax und Nullmengen-Grenze

**Status:** in-progress (welle-58-trace-handbuch-parsergrenzen).

**Welle:** welle-58-trace-handbuch-parsergrenzen (Trigger: reproduzierter
Konsumentenbefund in `m-trace`: 371 tabellarisch definierte Anforderungen
ergaben mit d-check v0.42.0 `0 Anforderungen, 0 Waisen` und Exit 0).

**Bezug:** Dokumentationskorrektur zum bestehenden Verhalten von
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix),
[`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code),
[`DC-FA-COV-001`](../../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
und
[`DC-FA-MOD-001`](../../../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in).
**Kein Change Request** (kein Verhaltens- oder Vertragsdelta), **kein ADR**
(keine neue Architekturentscheidung), **kein Release** (nur abgeleitete
Nutzer-Dokumentation und ihre Verifikation).

**Autor:** pt9912. **Datum:** 2026-07-14.

---

## 1. Ziel

Das Benutzerhandbuch beschreibt die konfigurierbaren Trace-Quellen, lässt
aber die feste Markdown-Grammatik der Anforderungsdefinition offen. Dadurch
legt die Aussage, eine angepasste Regex bilde andere Repo-Konventionen
vollständig ab, fälschlich nahe, auch tabellarische Lastenhefte würden
eingelesen. Das Handbuch macht die tatsächliche Definitionssyntax, die
Modalitätsquelle, die Waisen- und Referenzsemantik sowie die derzeit
fail-open behandelte Nullmenge ausdrücklich und gibt tabellenbasierten
Brownfield-Repos einen sicheren Migrationsweg.

## 2. Entscheidungen

- **Definition ist eine ATX-Überschrift:** Die Requirement-ID muss als erstes
  vollständiges Heading-Token auf `trace.requirements.id-pattern` passen.
  Tabellen, Listen, Fließtext und Setext-Überschriften definieren keine
  Anforderungen. Die Regex konfiguriert nur die ID-Gestalt, nicht die
  Dokumentgrammatik.
- **Modalität kommt aus dem Heading-Body:** `modality: {}` aktiviert die
  Built-in-Defaults; klassifiziert wird nur der Body bis zur nächsten
  gleich- oder höherrangigen Überschrift. Eine Tabellenspalte `Muss` zählt
  nicht.
- **Nullmenge ist kein grünes Vollständigkeitssignal:** Fehlende Quelle,
  unpassendes Format oder null ID-Treffer liefern in v0.42.0 eine leere RTM;
  auch `--require-complete` endet dann mit Exit 0. Das Handbuch verlangt vor
  Gate-Bindung eine Plausibilisierung von `total` gegen den erwarteten Bestand.
- **Waise ist exakt `¬slice ∧ ¬coverage`:** Eine ADR-Referenz allein
  verhindert `WAISE` nicht. Die Definition steht bereits in der ersten
  Ergebnis-Erklärung statt erst im Coverage-Unterabschnitt.
- **Referenzscan wird operational beschrieben:** rekursiv unter `dir`,
  `file-pattern` gegen den Basisdateinamen, Capture-Gruppe 1 plus
  `id-prefix` als Owner-Kennung; Dateien ohne Match werden übersprungen,
  Treffer im gesamten Dateiinhalt dedupliziert und sortiert.
- **Brownfield-Migration statt stiller Fehlannahme:** Tabellenbestände werden
  in Heading-plus-Body-Form überführt oder deterministisch in eine solche
  Trace-Quelle projiziert; die Projektion braucht eine eigene Drift-Prüfung.
- **Produktgrenzen bleiben sichtbar:** Fail-closed bei null Anforderungen
  und native konfigurierbare Tabellenspalten sind Folge-CR-Kandidaten, nicht
  Teil dieses Dokumentations-Slice.

## 3. Definition of Done

- [x] Handbuch öffnet §4.12 mit der exakten Definitionssyntax samt
  positivem und negativem Beispiel.
- [x] Handbuch warnt bei `total: 0` ausdrücklich vor einem irreführenden
  Exit 0 und beschreibt die Gate-Plausibilisierung.
- [x] Modalitäts-, Waisen- und Referenzscan-Semantik sind an einer Stelle
  eindeutig und widerspruchsfrei dokumentiert.
- [x] Brownfield-Migration für tabellarische Lastenhefte ist beschrieben;
  native Tabellenunterstützung wird nicht behauptet.
- [x] Die Handbuch-Verifikation verriegelt die neuen Kernhinweise gegen
  stilles Entfernen; bestehende Trace-E2E-Ausgabe bleibt grün.
- [x] Engster Sensor und `make gates` sind grün; Handbuch-Historie und
  Changelog nennen die Dokumentationskorrektur.

## 4. Risiken / offene Punkte

- **Dokumentation des Mangels ist kein Fix:** Konsumenten müssen bis zu einem
  Folge-CR die erkannte Anzahl selbst plausibilisieren.
- **Generierte Zwischenquelle:** Eine Projektion aus Tabellen kann driften;
  deshalb darf sie nur mit deterministischer Erzeugung und eigener
  Konsistenzprüfung als Gate-Quelle dienen.
- **Versionsbezug:** Die Warnung beschreibt ausdrücklich v0.42.0. Ein späterer
  Fail-closed- oder Tabellen-CR muss diesen Abschnitt gemeinsam mit dem
  Verhalten fortschreiben.

## 5. Trigger

Auftraggeber-Befund 2026-07-14 am Konsumenten `m-trace`: Das Lastenheft
enthält 371 Anforderungen in Tabellenzeilen. Die explizite Trace-Konfiguration
mit `requirements.source: spec/lastenheft.md` und passender ID-Regex liefert
dennoch `0 Anforderungen, 0 Waisen` bei Exit 0. Quellcode-Analyse zeigt die
nicht im Handbuch genannte Heading-Grammatik.

## 6. Sub-Area-Modus-Begründung

GF für `docs/user` und die bestehende Handbuch-Verifikation: Die ausgelieferte
Implementierung ist die Evidenz, das abgeleitete Handbuch wird daran
angeglichen. Kein neuer Adapter, keine BF-/Hybrid-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

Offen bis DoD, Review und Gates erfüllt sind. Self-Review R1
([Report](../../../reviews/2026-07-14-slice-069-trace-handbuch-r1.md)) mit
Verdikt NACHBESSERN (0 HIGH/2 MEDIUM/1 LOW): alle drei Befunde eingearbeitet
(einleitende Waisen-Definition, `modality.require-levels`-Gate-Semantik,
konditionale JSON-/YAML-Felder); Folgereview steht aus.
