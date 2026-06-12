# Slice slice-016: Modul `hostpaths` — host-lokale absolute Pfade

**Status:** done.

**Welle:** welle-06-sensorik (Abschluss).

**Bezug:** [`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)
(Change Request 0.6.0),
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
(Modul-Auswahl),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(Konfigurations-Vollvalidierung, `hostpaths.prefixes`),
[ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md) (Layout).

**Autor:** pt9912. **Datum:** 2026-06-12.

---

## 1. Ziel

Das opt-in-Modul `hostpaths` meldet Maschinen-Layout-Leaks
(`hostpath-forbidden`) — dieselbe Klasse, die in den
[`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Läufen
dreimal unabhängig auftrat: als eigener bess-ems-Rest-Sensor, als
8 Host-Pfad-Links im d-migrate-Vergleich und als manuelle
0.2.2-Hygiene-Korrektur im eigenen Lastenheft. Mit dem Modul wird
der bess-ems-Rest-Sensor ablösbar (dortige Entscheidung).

## 2. Definition of Done

- [x] Spezifikation fortgeschrieben: §`DC-FA-HOST-001.a`
  (Muster-Definition: Unix-Präfixliste als erstes Pfad-Segment,
  Windows-Laufwerk/UNC fest; Wortgrenzen-Vorbedingung;
  Satzzeichen-Abtrennung wie `codepaths`-Normalisierung;
  Fence-Ausnahme), `.d-check.yml`-Schema um `hostpaths.prefixes`
  (Constraint: nicht-leere Verzeichnisnamen ohne `/`), Grund-Code
  `hostpath-forbidden` (§4).
- [x] Modul implementiert (Layout nach
  [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md); nutzt
  die geteilten Prosa-/Fence-Helfer; prüft Inline-Code **mit** —
  abweichend von den Strip-Konsumenten, wie spezifiziert); die drei
  Akzeptanzkriterien aus
  [`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)
  als Tests, Determinismus gemäß
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus).
- [x] **Paritäts-Gegentest bess-ems:** auf dem dortigen Stand liefert
  `hostpaths` dieselben Ergebnisse wie der Rest-Sensor
  (`check_markdown_links.py`, Host-Pfad-Teil): beidseitig 0 im
  Bestand, plus ein konstruiertes Fixture mit Unix-/Windows-/
  UNC-Fall, das beide gleich melden — Beleg für die Ablösbarkeit.
- [x] **Kalibrierungslauf** gegen das Golden-Set der dreizehn
  migrierten Repos plus pkcs11-course: Befunde sind echte Leaks
  (Fix im Ziel-Repo anbieten) oder False-Positives (Muster per
  Spec-Fortschreibung verengen, vor Abschluss) — Lehr-Repos mit
  bewussten Pfad-Beispielen sind der Härtetest für die
  Fence-Ausnahme.
- [x] Dogfooding: Selbstkonfiguration
  ([`.d-check.yml`](../../../../.d-check.yml)) aktiviert
  `hostpaths`; eigene Doku befundfrei (Bestand ist geprüft sauber);
  `gate-consistency`-gebundene `DC-QA-03`-Modulliste nachgezogen.
- [x] Angebot an bess-ems dokumentiert: Rest-Sensor-Ablösung durch
  `hostpaths` + Pin-Hebung (Umsetzung ist dortige Entscheidung).
- [x] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md);
  Closure-Notiz mit Steering-Loop-Lerneintrag. Release-Hinweis:
  fließt mit `spans` in das nächste Minor-Release (v0.4.0 —
  v0.3.0 ist seit slice-017 vergeben).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | §`DC-FA-HOST-001.a`, Schema (`hostpaths.prefixes`), Grund-Code |
| Modul `hostpaths` im Hexagon-Kern | neu | Prüflogik (geteilte Prosa-/Fence-Helfer) |
| Tests (AK-Trio, Paritäts-Gegentest, Kalibrierung) | neu | Beleg-Pflicht |
| [`.d-check.yml`](../../../../.d-check.yml) | update | Dogfooding-Aktivierung |
| [`CHANGELOG.md`](../../../../CHANGELOG.md) | update | nutzersichtbares Modul |

## 4. Trigger

Change Request 0.6.0 im Lastenheft (erfüllt 2026-06-12) **und**
Priorisierung durch den Auftraggeber (welle-06). Sinnvoll **nach**
[slice-015](../done/slice-015-spans-modul.md) (erfüllt 2026-06-12 — gemeinsames Minor-Release;
beide Module teilen die Prosa-Helfer — wer zuerst läuft, legt die
Schnittstelle).

## 5. Closure-Trigger

DoD vollständig (insbesondere bess-ems-Paritäts-Gegentest und
Kalibrierungslauf) + Closure-Notiz.

## 6. Risiken und offene Punkte

- **Inline-Code-Einschluss ist die FP-Quelle Nr. 1:** legitime
  How-to-Erwähnungen von Host-Verzeichnissen in Backticks würden
  gemeldet. Bewusste Setzung (dort leben die echten Leaks); die
  Auswege — Fences für Beispiele, `hostpaths.prefixes` für
  Sonderfälle — müssen der Kalibrierungslauf bestätigen. Kippt das,
  ist die Verengung (Inline-Code ausnehmen) eine
  Spec-Fortschreibung, kein Marker.
- Die Default-Präfixliste ist Workspace-geprägt (Development, home,
  Users, …); generische System-Pfade (usr, etc, opt) fehlen bewusst
  — Lehrinhalte erwähnen sie legitim. Erweiterung ist per Config
  möglich, Default bleibt konservativ.
- Lehr-Repos (pkcs11-course, euler, Kurs-Repo) sind die kritischen
  Kalibrierungs-Korpora — bewusst im DoD verankert.

## 7. Closure-Notiz (nach `done/`)

**Abgeschlossen:** 2026-06-12 (e50f4b8). Schließt welle-06 mit dem
v0.4.0-Release (gemeinsam mit
[slice-015](../done/slice-015-spans-modul.md)).

- **Kalibrierungslauf (14 Korpora)** in zwei Iterationen — die erste
  fand zwei False-Positive-Klassen, beide vor Abschluss verengt
  (der DoD-Pfad): (a) Regex-Beispiele wie escapte Punkt-Folgen
  matchten als UNC → Servername beginnt jetzt alphanumerisch
  (Spez); (b) Laufzeit-Orte unter tmp sind überwiegend legitime
  Betriebs-Doku (Log-/Output-Pfade in d-migrate, m-trace, u-boot,
  pkcs11-course) → tmp aus der Default-Präfixliste gestrichen
  (Lastenheft 0.7.2). Die zweite Iteration trennte sauber:
  **echte Workspace-Leaks** (k-deskflight 13×, grid-gym-Archiv 5× —
  jeweils Verweise auf das m-trace-Arbeitsverzeichnis) von
  **gewollter Plattform-Doku** (Windows-/WSL-/UNC-Inhalte in
  cmake-xray, u-boot-Templates, m-trace-WSL-Anleitung,
  pkcs11-Kurs — vertragsgemäße Treffer; Ausweg ist der
  Opt-in-Verzicht bzw. Fences, keine Vertragsänderung).
- **Fixes nach der Audit-Trail-Linie** (grid-gym-CR-Philosophie):
  nur lebende Doku angefasst — k-deskflight `spec/architecture.md`
  (2 Leaks → Repo-Name, b093636). Accepted-ADRs und das
  grid-gym-`done-archive/` (explizit eingefroren) bleiben
  unverändert; deren Ausweg bei Modul-Aktivierung ist
  `hostpaths.scope.ignore` — wofür slice-017 gebaut wurde.
- **Paritäts-Gegentest bess-ems:** Bestand beidseitig 0; Fixture
  (Unix/Windows/UNC, Fence-Ausnahme, URL- und Relativ-Negativfälle)
  — identische Befunde auf identischen Zeilen, d-check normalisiert
  zusätzlich den Span-Backtick aus dem Pfad, den der Rest-Sensor
  mitnahm. **Angebot an bess-ems:** der Rest-Sensor
  (`check_markdown_links.py`, ~130 Z.) ist durch `hostpaths`
  ablösbar; Delta dokumentiert: deren Liste enthält tmp — bei
  Ablösung `hostpaths.prefixes` inklusive tmp setzen, falls die
  schärfere Linie gewünscht bleibt. Umsetzung ist dortige
  Entscheidung (Pin-Bump auf v0.4.0 + Config + Stage-Löschung).

### Lerneintrag (Steering Loop)

Kalibrierung in zwei Iterationen ist der eigentliche Mechanismus:
Die erste Runde hat nicht nur Muster-Fehler (UNC/Regex) gezeigt,
sondern eine *konzeptionelle* Fehleinschätzung des CR — tmp stand
in der Default-Liste, weil der bess-ems-Sensor es so hatte, aber
der breite Korpus bewies das Gegenteil (Laufzeit-Doku ≠ Leak).
Geschärfte Regel: Default-Listen aus Einzel-Repo-Sensoren nicht
ungeprüft generalisieren — der Kalibrierungslauf über den vollen
Korpus entscheidet, was Default wird und was Konfiguration bleibt.
Zweitens bestätigt: die Audit-Trail-Linie (historische Bestände
nicht für neue Sensoren umschreiben) braucht den Modul-Scope aus
slice-017 als Gegenstück — die beiden CRs greifen ineinander.

## 8. Sub-Area-Modus-Begründung

Dieses Repo: GF (Spec führt, Code folgt — Change Request 0.6.0 vor
Implementierung).
