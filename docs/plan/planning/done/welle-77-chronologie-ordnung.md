# Welle welle-77-chronologie-ordnung: Eine chronologische Tabelle hält ihre Richtung

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-77-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Erweiterung eines bestehenden
Moduls aus dem eigenen Beobachtungs-Register).

**Verantwortlich:** pt9912. **Datum:** 2026-08-21.

---

## 1. Welle-Ziel

**BEO-005** mechanisieren: eine chronologische Tabelle kippt still ihre
Richtung — aus „unten anhängen“ wird irgendwann „oben einfügen“, und danach
führt dieselbe Tabelle zwei gegenläufige Blöcke; kein Gate liest Reihenfolge.
[slice-105](slice-105-tabellen-monotonie.md) liefert den
**typisierten Monotonie-Vergleich der Schlüsselspalte** als siebte Bedingung im
Modul `structure` (Erweiterung von
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
kein neues Modul — [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md)-Kriterium)
und aktiviert sie nach Messung auf den sechs eigenen Bestandstabellen.

**Die Beweislage ist ausdrücklich schwächer als bei den drei aktivierten
`structure`-Regeln** (Register-Eintrag): der Ist-Bestand ist grün, der Nutzen
ist Prävention, der Beleg ein Retro-Lauf — am Stand vor der Heilung fand der
typisierte Vergleich alle drei gekippten Tabellen (14 · 6 · 7 Verletzungen),
danach null; ein naiver String-Vergleich meldet dagegen drei korrekt sortierte
Tabellen rot. Genau dieser Beleg wird bei der Closure **mit dem Produkt**
wiederholt statt mit dem Mess-Skript behauptet.

## 2. Trigger (Welle startet)

Freigabe des Auftraggebers (2026-08-21, mit Reihenfolge-Entscheid: erst diese
Welle, dann die Baseline-Migration v5.6.0), WIP-Slot frei (welle-76 geschlossen,
`in-progress/` trägt nur die Roadmap). Der Vorschau-Trigger ist seit 2026-08-16
eingetreten: welle-75 ist geschlossen, die Tabellen-Lexik hat mit dieser
Bedingung ihren **dritten** Konsumenten und verdient einen Kopplungs-Test statt
Einzel-Assertionen.

## 3. Closure-Trigger (Welle schließt)

- [slice-105](slice-105-tabellen-monotonie.md) liegt in `done/`.
- **Der Retro-Beleg lief mit dem Produkt:** die drei historisch gekippten
  Tabellen sind am Vor-Heilungs-Stand rot (14 · 6 · 7), der heutige Bestand
  ist null; die naive Gegenprobe (String-Vergleich rot auf korrekt Sortiertem)
  ist als Testfall festgehalten — der Typ ist belegt Pflicht, kein Komfort.
- **Der Kopplungs-Test der Tabellen-Lexik läuft** über alle drei Konsumenten
  (`targets`, `planning.waves`, `structure`) — die Form aus
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md).
- **Register-Schritt zu BEO-005 entschieden:** gestrichen (mechanisiert für
  die aktivierten Tabellen, Rest benannt — dieselbe Form wie BEO-001) oder
  offen mit Begründung.
- Release als **Minor** (opt-in; ohne die neuen Schlüssel byte-identisch).
- `make fullbuild` grün.

## 4. Slices in dieser Welle

| Slice | Rolle |
|---|---|
| [slice-105](slice-105-tabellen-monotonie.md) | siebte `structure`-Bedingung: typisierte Chronologie-Monotonie, geteilte Tabellen-Lexik samt Kopplungs-Test, Selbst-Aktivierung nach Messung |

## 5. Abhängigkeiten

- [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  und das Modul liegen vor ([slice-099](slice-099-structure-modul.md),
  [ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md)); erweitert
  wird eine bestehende Anforderung.
- Die Tabellenzeilen-Lexik ist entdriftet (welle-74) und trägt seit welle-75
  eine Spalten-Adresse in `planning.waves` — sie muss gehoben, nicht neu
  gebaut werden.

## 6. Out-of-Scope für diese Welle

- **Kein Auto-Umsortieren.** d-check bleibt diagnose-only; die Reparatur einer
  gekippten Tabelle ist eine Autoren-Entscheidung.
- **Keine offene Typ-Registry.** Die Typ-Menge ist geschlossen (ISO-Datum,
  Punkt-Version); weitere Typen sind ein eigener Change Request.
- **Keine Ordnung über mehrere Tabellen oder Spalten hinweg** — geprüft wird
  die Schlüsselspalte je zusammenhängender Tabelle (zwei getrennt sortierte,
  gegenläufige Tabellen im selben Abschnitt bleiben eine benannte Grenze).
- Die **Baseline-Migration v5.0.0 → v5.6.0** — nächste Welle in der Vorschau
  (Nutzer-Entscheid 2026-08-21).

## 7. Closure-Notiz

Geschlossen am 2026-08-21. Alle Closure-Trigger erfüllt:
[slice-105](slice-105-tabellen-monotonie.md) liegt in `done/` (Release
**v0.61.0**, Digest `sha256:0e731cfc…9f98`), der Retro-Beleg lief mit dem
Produkt (27 Befunde am Stand vor der welle-73-Heilung — exakt die 14 · 6 · 7
der Skript-Messung über die drei gekippten Tabellen, die drei übrigen null;
heutiger Bestand null), die naive Gegenprobe ist als Testfall festgeschrieben,
der Kopplungs-Test bindet alle drei Konsumenten der Tabellen-Lexik, und der
Register-Schritt ist entschieden: **BEO-005 gestrichen** — die sechs eigenen
Bestandstabellen sind mechanisiert, der Rest steht benannt im
Streichungs-Eintrag. Das Release ist ein Minor (opt-in; ohne `table-order`
byte-identisch). Die Release-Prep dieser Welle war der erste Kunde der eigenen
Fähigkeit: der Handbuch-§11-Eintrag saß im ersten Anlauf falsch herum und wäre
ab diesem Release maschinell rot gewesen. Was wirkte und was anders lief:
[welle-77-results.md](welle-77-results.md).
