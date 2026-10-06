# ADR-0099: `targets.makefiles` nimmt Glob-Muster an; ein Glob ohne Treffer ist Exit 2

**Status:** Proposed

**Datum:** 2026-10-06

**Autor:** pt9912

**Bezug:** [`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(erweitert, Lastenheft 0.95.0); Anlass ist der eingehende
[CR von `ai-harness-init`](../cr/2026-10-06-cr-eingehend-ai-harness-init-targets-makefiles-glob.md);
Schnitt-Kriterium (Einzelmodul-Frage ⇒ bestehende Anforderung ändern) aus
[ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md); Form der
Konfigurations-Weitung wie [ADR-0058](0058-konfigurations-flaechen-additiv-weiten.md);
slice-250 <!-- d-check:status-provenance -->.

**Schärft:**
[`spec/spezifikation.md` §DC-FA-TGT-001.a](../../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
(Schritt 1a) und die Schema-Zeile `targets.makefiles` unter
[`SPEC-005`](../../../spec/spezifikation.md#spec-005--d-checkyml).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Modul `targets` liest jeden `targets.makefiles`-Eintrag wörtlich. Ein
Adopter, der seine make-Targets auf ein Root-Makefile und eine wachsende
Menge von Fragmenten verteilt, muss jedes Fragment einzeln listen; ein
vergessenes Fragment ist still ungeprüft. Ein Glob in der Konfiguration ist
die einzige Angabe, die mit dem Bestand mitwächst. Heute passiert ein Glob
den Config-Rand und scheitert erst im Lauf, weil er als Dateiname gelesen
wird.

Das Repo kennt zwei Vorbilder für eine Glob-Menge ohne Treffer: das Modul
`file` meldet einen **Befund** (`file-no-match`), das Modul `mentions`
**bricht ab** (Exit 2).

## Entscheidung

1. **Ein Eintrag mit Glob-Zeichen (`*`, `?`, `[`) ist ein Muster** und
   expandiert per `matchGlob` (segmentweise, `**`) gegen die Repo-Wurzel —
   dieselbe Semantik wie die übrigen Glob-Schlüssel. Ein Eintrag ohne
   Glob-Zeichen bleibt ein wörtlicher Pfad; ohne Glob-Eintrag ist der
   Befundsatz byte-identisch.
2. **Ein Glob ohne Treffer ist Exit 2, kein Befund.** Maßgeblich ist das
   Modul selbst, nicht der Nachbar: eine fehlende wörtliche Makefile-Datei
   ist in `targets` seit jeher Exit 2. Ein Glob ist dieselbe Behauptung
   („hier liegen Makefiles"), nur mit mehreren Kandidaten; sie anders zu
   behandeln, hieße zwei Fehlerformen für einen Konfigurationsfehler. Ein
   Befund bräuchte zudem einen neuen Grund-Code und sähe aus wie ein
   Doku-Mangel.
3. **Eine mehrfach erfasste Datei zählt einmal** — erste Nennung in
   Konfigurations-Reihenfolge, Treffer je Muster sortiert. Das hält den
   Befundsatz deterministisch und verhindert doppelte Befunde je Regelzeile.
4. **Gewandert wird ab dem festen Verzeichnis-Präfix des Musters**, mit
   denselben übersprungenen Verzeichnissen wie das Modul `file`;
   Symlinks werden nicht verfolgt. Die Auflösung ist unabhängig von
   `scan.roots`/`scan.ignore` — wie ein wörtlicher Eintrag; die Grenze steht
   in der Anforderung (`AGENTS.md` §3.8).
5. **Ein ungültiges Muster ist beim Laden Exit 2** (segmentweise Prüfung
   wie bei `tracked.exempt-targets`); die bisherige Pfad-Regel (relativ,
   kein `..`) gilt unverändert auch für Muster.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Leerer Glob als Befund** (wie `file`) | im Bericht sichtbar, Exit 1 | neuer Grund-Code; derselbe Konfigurationsfehler hätte im selben Modul zwei Formen (wörtlich: Exit 2, Glob: Befund) |
| **Leerer Glob still inert** | Bootstrap-Phase ohne Fragmente bleibt grün | genau die Lücke, die der CR schließen will: das Modul prüfte unbemerkt nichts |
| **`include`-Direktiven auswerten** | keine zweite Angabe in der Konfiguration | vom CR ausgeschlossen, im Lastenheft Out-of-Scope; liest Makefile-Semantik statt einer prüfbaren Angabe |
| **Glob mit Exit 2 bei leerer Menge** (gewählt) | konsistent zum wörtlichen Fall im selben Modul; fail-closed | ein Adopter mit noch leerem Fragment-Verzeichnis muss den Glob erst eintragen, wenn es ein Fragment gibt |

## Konsequenzen

- `targets.makefiles` trägt Pfade und Muster; die Schema-Zeile und das
  `--print-config`-Gerüst nennen beides.
- Ein Adopter mit Fragment-Verzeichnis listet ein Muster statt jeder Datei.
- Ein Muster wie `**/*.mk` sieht in fest übersprungene Verzeichnisse
  (`build`, `vendor`, …) nicht hinein — benannte Grenze.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test `TestCheckTargetsMakefileGlob` | Expansion, Fundstelle im Fragment, Dublette einmal, `*` bleibt im Segment | `make test` |
| Go-Test `TestCheckTargetsMakefileGlobDoppelstern` | `**` erfasst Unterverzeichnisse | `make test` |
| Go-Test `TestCheckTargetsMakefileGlobLeer` | Glob ohne Treffer ⇒ Fehler (Exit 2) | `make test` |
| Go-Test `TestDecode_TargetsMakefilesGlob` | ungültiges Muster, führender `/`, `..` ⇒ Konfigurationsfehler | `make test` |

## Re-Evaluierungs-Trigger

Ein Adopter meldet, dass Exit 2 bei leerem Glob seine Bootstrap-Phase
blockiert und ein Ventil (etwa ein ausdrücklich optionales Muster) nötig
wird. Ohne das: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-06 | Angelegt als `Proposed`; `Accepted` erst mit der Closure des Vorgangs, nach Review und Verifikation |
