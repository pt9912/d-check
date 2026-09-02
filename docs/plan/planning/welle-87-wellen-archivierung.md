# Welle welle-87: Die Wellen-Archivierung wird nachgerüstet

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach in `docs/plan/planning/`; bei der Closure wandert sie nach `done/` und
bekommt ihre Ergebnisnotiz `welle-87-results.md` daneben.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-02.

---

## 1. Welle-Ziel

**Dieses Repo hat noch keine einzige Welle archiviert — obwohl der Kanon das
für jede Wellen-Closure vorschreibt, seit die Regel eingeführt ist.**
Gemessen bei [slice-188](done/slice-188-register-gegen-neuen-kanon.md):
kein `archiv.zip`, kein `done/<welle-id>/`-Verzeichnis existiert, während
welle-60 bis welle-85 bereits geschlossen sind.

<!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md:238-240 -->
> **Zeitdokumente der Welle archivieren.** Die Slice-Dateien, die sie
> einsammelt, ihr eigener Plan und die Review-Reports dieser Slices wandern
> in ein unveränderliches Archiv `done/<welle-id>/archiv.zip`.

**Die Nachrüst-Ausnahme des Kanons deckt nur die Vergangenheit, nicht die
Zukunft.** Wellen, die vor Einführung der Regel schlossen, müssen nicht
archiviert werden — ein Repo bleibt ohne das konform. [`welle-86`](welle-86-closure-uebergang-durchsetzen.md)
steht aber **aktuell offen**: Schließt sie, greift die Pflicht ungemindert,
und heute führt sie kein Werkzeug aus. Diese Welle baut das Werkzeug **und**
wendet es rückwirkend auf den Alt-Bestand an — bevor welle-86 selbst vor
derselben Lücke steht.

**Das Mehr gegenüber den einzelnen Slice-DoDs:** keine einzelne Slice-DoD
kann belegen, dass die Archivierungs-**Operation als Ganzes** — Auswahl,
ZIP, Stubs, Verweis-Nachzug, in beiden Pfad-Formen — für einen realen
Wellen-Bestand funktioniert. Das ist ein Werkzeug-Belastungstest, kein
Einzelkriterium.

## 2. Trigger (Welle startet)

- [slice-188](done/slice-188-register-gegen-neuen-kanon.md) ist geschlossen
  — er hat die Lücke gemessen und den Anlass geliefert.

## 3. Closure-Trigger (Welle schließt)

- Ein Werkzeug (Go-Programm oder Bash-Skript nach dem Muster von
  `tools/harness/*.sh`, nie ein Host-Skript-Interpreter, `AGENTS.md` §3.1)
  führt die Archivierungs-Operation aus: sammelt die Vorgänge einer Welle,
  baut `done/<welle-id>/archiv.zip`, erzeugt die Stubs nach
  `archiv-stub-slice.template.md`/`archiv-stub-welle.template.md`, lässt
  Review-Reports ohne Stub, zieht eingehende Verweise in beiden Pfad-Formen
  nach.
- Alle vor dieser Welle geschlossenen Wellen (welle-60 bis welle-85) sind
  archiviert — oder es ist begründet dokumentiert, warum eine Teilmenge
  zurückgestellt wird.
- Die wellenlosen Alt-Slices vor der Einführung sind einer Zuordnung
  zugeführt (Kanon verlangt eine Entscheidung: chronologisch nächste
  geschlossene Welle oder ein Sammel-Archiv — diese Welle trifft sie).
- `make gates` und `make fullbuild` grün auf dem archivierten Bestand.
- Closure-Notiz in `welle-87-results.md`.

## 4. Slices in dieser Welle

<!-- Grob geschnitten bei Eröffnung; jeder Slice wird einzeln in open/
geplant, nicht hier im Detail vorweggenommen (Modul 5 §Regeln gegen
typische Fehlannahmen: "Erst plan ich alle Slices"). -->

| Slice | Titel | Bezug |
|---|---|---|
| slice-<NN-A> | Archivierungs-Werkzeug bauen (Sammeln, ZIP, Stubs, Verweis-Nachzug) | `modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 4 |
| slice-<NN-B> | Alt-Bestand archivieren (welle-60…welle-85) und Zuordnung der wellenlosen Alt-Slices entscheiden | dito |

## 5. Abhängigkeiten

- **Wird gebraucht von:** [`welle-86`](welle-86-closure-uebergang-durchsetzen.md)
  — ihre eigene Closure trifft dieselbe Archivierungspflicht, sobald ihre
  vier Vorbedingungen erfüllt sind. `welle-86` bleibt eigenständig und wird
  **nicht** in diese Welle eingesammelt; sie nutzt nur das hier gebaute
  Werkzeug bei ihrer eigenen Closure.
- Keine Blockade in der Gegenrichtung: `welle-86` kann unabhängig weiterlaufen.

## 6. Out-of-Scope für diese Welle

- **`welle-86` selbst schließen.** Eigene Welle, eigene vier Vorbedingungen —
  siehe [`welle-86`](welle-86-closure-uebergang-durchsetzen.md) §1.
- **[`BEO-027`](observations.md) auflösen** (Registerzeilen ohne
  formgültigen Ausgang). Verwandtes, aber eigenständiges Thema — eigener
  Folge-Slice, keine Vermischung mit der Archivierungs-Mechanik.
- **Der `registry`-Modul-Vorschlag** aus [`BEO-001`](observations.md#gestrichene-einträge)
  (ADR-Index/Konventionsspeicher-Index gegen dieselbe Drift). Andere
  Register-Klasse, eigene Entscheidung.

## 7. Vorgelagert (Eröffnungs-Schritt 2 — offene Beobachtungen sichten)

Register-Stand 2026-09-02, höchste Kennung `BEO-027`. Keine Beobachtung
trifft die Archivierungs-Mechanik selbst direkt — zwei liegen in derselben
Sub-Area (`docs/plan/planning/`) und werden bewusst **nicht** mit
eingesammelt (siehe §6): [`BEO-015`](observations.md) (ein Slice bekommt bei
der Closure einen Ausgang, den es nicht gibt — dieselbe Familie wie
`BEO-027`, aber die andere Richtung) und [`BEO-027`](observations.md)
selbst (Registerzeile ohne zugewiesenen Ausgang). Keine Treffer sind
ebenfalls eine Antwort.

## 8. Closure-Notiz

Ergebnis: —
Zähler: —
