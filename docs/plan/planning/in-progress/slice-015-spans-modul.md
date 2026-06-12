# Slice slice-015: Modul `spans` — Markdown-Span-Artefakte

**Status:** in-progress.

**Welle:** welle-06-sensorik (per Roadmap-Fortschreibung;
Start bei Priorisierung durch den Auftraggeber).

**Bezug:** [`DC-FA-SPAN-001`](../../../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
(Change Request 0.5.0),
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
(Modul-Auswahl),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(Konfigurations-Vollvalidierung),
[ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md) (Layout).

**Autor:** pt9912. **Datum:** 2026-06-12.

---

## 1. Ziel

Das opt-in-Modul `spans` meldet die zwei in den
[`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Vergleichsläufen
empirisch belegten Artefakt-Klassen direkt an der Ursache:
ungeschlossene Code-Spans (`span-unclosed`) und verschachtelte
Link-Artefakte (`span-nested-link`). Bisher wurden diese Fälle nur
indirekt sichtbar — als `id-unlinked`-Folgefehler ganzer Absätze
(~100 Stellen in u-boot, je eine Kaskade in grid-gym und bess-ems)
oder gar nicht (Repos ohne `ids`-Konfiguration).

## 2. Definition of Done

- [x] Spezifikation fortgeschrieben: §`DC-FA-SPAN-001.a`
  (Opener-Klassifikation, Absatz-Semantik identisch zur
  Vorverarbeitung §`DC-FA-LINK-001.a` Schritt 2, Muster-Definition
  `span-nested-link` nach Inline-Code-Stripping), `modules`-Liste im
  `.d-check.yml`-Schema, Grund-Codes-Tabelle (§4) um `span-unclosed`
  und `span-nested-link` ergänzt.
- [x] Modul implementiert (Layout nach
  [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md); teilt
  `proseParagraphs`/`forEachInlineCodeSpan` mit der Vorverarbeitung —
  keine Drittkopie der Span-Logik); die drei Akzeptanzkriterien aus
  [`DC-FA-SPAN-001`](../../../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
  als Tests, Determinismus gemäß
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus).
- [x] **Kalibrierungslauf** gegen das Golden-Set der dreizehn
  migrierten Schwester-Repos (aktueller Stand): alle müssen
  span-befundfrei sein oder ausschließlich echte Befunde zeigen —
  die „alleinstehend literal"-Heuristik wird daran kalibriert,
  Abweichungen vor Abschluss triagiert (False-Positive ⇒
  Spec-Fortschreibung hier).
- [x] **Gegentest** gegen die historischen Stände vor den
  Rollout-Fixes (u-boot vor 470be86, grid-gym vor 766ae8c, bess-ems
  vor a6c8b9b): das Modul muss die dort gefixten Artefakt-Stellen
  als `span-unclosed` melden — der Beleg, dass der Sensor die
  bekannte Fehlerklasse trifft.
- [x] Dogfooding: Selbstkonfiguration
  ([`.d-check.yml`](../../../../.d-check.yml)) aktiviert `spans`;
  eigene Doku befundfrei; `gate-consistency`-gebundene
  `DC-QA-03`-Modulliste nachgezogen (netzloser Lauf aller Module
  außer `external`).
- [x] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md);
  Closure-Notiz mit Steering-Loop-Lerneintrag. Release-Hinweis:
  neues Modul ⇒ nächstes Minor-Release (v0.4.0 — v0.3.0 ist seit slice-017 vergeben —, gemeinsam mit
  [slice-016](../open/slice-016-hostpaths-modul.md)), Konsumenten-Pins
  sind Routine-Hebungen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | §`DC-FA-SPAN-001.a`, Schema, Grund-Codes |
| Modul `spans` im Hexagon-Kern | neu | Prüflogik (geteilte Span-/Absatz-Helfer) |
| Tests (AK-Trio, Kalibrierung, Gegentest) | neu | Beleg-Pflicht |
| [`.d-check.yml`](../../../../.d-check.yml) | update | Dogfooding-Aktivierung |
| [`CHANGELOG.md`](../../../../CHANGELOG.md) | update | nutzersichtbares Modul |

## 4. Trigger

Change Request 0.5.0 im Lastenheft (erfüllt 2026-06-12) **und**
Priorisierung durch den Auftraggeber (Start welle-06). Kein
Release-Vorlauf nötig — implementiert wird gegen den lokalen Stand,
das Minor-Release folgt mit dem Slice-Abschluss.

## 5. Closure-Trigger

DoD vollständig (insbesondere Kalibrierungslauf 13/13 und
historischer Gegentest) + Closure-Notiz.

## 6. Risiken und offene Punkte

- **False-Positive-Kalibrierung** ist das Kernrisiko: die
  Unterscheidung „Autorenfehler" (Opener klebt an Nicht-Whitespace,
  ungeschlossen) vs. „beabsichtigt literal" (alleinstehend) ist eine
  Konvention, keine CommonMark-Wahrheit. Das Golden-Set (DoD) ist
  der Sensor dafür; bleibt die Heuristik zu scharf, wird die
  Erkennungsregel per Spec-Fortschreibung verengt — es gibt bewusst
  keinen Opt-out-Marker.
- `span-nested-link` ist im Bestand selten (nur u-boot-Rest-Lint
  kannte es); der Nutzen liegt in der Zukunftssicherung, der
  Aufwand ist klein (ein Muster auf vorverarbeiteten Zeilen).
- Abgrenzung: Emphasis-Artefakte (`*`/`_`) bleiben out-of-scope —
  bewusste Nicht-Behauptung, im Lastenheft dokumentiert.

## 7. Closure-Notiz (nach `done/`)

**Abgeschlossen:** 2026-06-12 (540e5c7). Release folgt gebündelt mit
[slice-016](../open/slice-016-hostpaths-modul.md) als v0.4.0.

- **Kalibrierungslauf (14 Korpora: 13 migrierte Repos +
  pkcs11-course):** 17 echte Artefakte, 1 False-Positive-Klasse.
  Echte Befunde — d-migrate 2× und m-trace 7× doppelt
  verschachtelte Links (`[[…](x)](x)`), u-boot 6×
  Doppel-Backtick-Tippfehler plus 1× escaped Backticks im Code-Span
  (CommonMark kennt kein Backtick-Escape in Spans), grid-gym 2
  ungeschlossene Datei-Spans — alle in den Ziel-Repos gefixt
  (8a47a88a, 439c2a9, 78995d1, 614009d); pkcs11-course:
  NuGet-Cache als Laufzeit-Residuum in `scan.ignore` (102ace9).
  **False-Positive:** das Badge-Muster `[![…](…)](…)` (legales
  Markdown, gefunden in einer vendorten Paket-README) → Lastenheft
  0.7.1 + Code-Fix vor Abschluss, genau der DoD-Pfad.
- **Historischer Gegentest:** u-boot-Stand vor den
  slice-014-Reparaturen (470be86^): **14 `span-unclosed`-Befunde**
  (davon 4 im CHANGELOG — die dokumentierte Titel-Span-Klasse),
  heutiger Stand 0. **DoD-Korrektur:** Der Plan nannte auch
  grid-gym (766ae8c^) und bess-ems (a6c8b9b^) als
  Gegentest-Stände — sachlich falsch: deren slice-014-Befunde waren
  ein *balancierter* mehrzeiliger Span (Parser-Thema, kein
  Artefakt) bzw. eine Slug-Kollabierung; beide gehören nicht zur
  span-Klasse. Der Gegentest läuft daher allein über u-boot.
- **Befund-Semantik dokumentiert:** Das Pairing macht den ersten
  Opener eines gekippten Absatzes oft zum (falsch) geschlossenen
  Span — der Befund zeigt auf die *übrigbleibende* ungeschlossene
  Folge des Absatzes, nicht zwingend auf den Tippfehler selbst
  (deterministisch, im Test verankert).

### Lerneintrag (Steering Loop)

Die Kalibrierungs-DoD hat sich doppelt bezahlt gemacht: (1) Der
einzige False-Positive der neuen Befund-Klasse (Badges) steckte
ausgerechnet in einem *Laufzeit-Residuum* — ohne den
pkcs11-Adoptions-Korpus wäre er erst beim ersten Badge-nutzenden
Konsumenten explodiert; ungepflegte Korpora sind als
Kalibrierungs-Material wertvoller als kuratierte. (2) Drei der vier
Befund-Quellen waren *neue* Fundstellen jenseits der
slice-014-Sweeps — ein Symptom-getriebener Fix-Sweep (damals:
id-unlinked-Folgefehler) findet nur, was Symptome zeigt; der
Ursachen-Sensor findet den Rest. Geschärfte Regel: DoD-Behauptungen
über historische Stände beim Schreiben des Plans verifizieren —
zwei der drei genannten Gegentest-Stände gehörten nicht zur
Befund-Klasse (im Closure korrigiert).

## 8. Sub-Area-Modus-Begründung

Dieses Repo: GF (Spec führt, Code folgt — Change Request 0.5.0 vor
Implementierung).
