# ADR-0112: Das Modul `metrics` hält Kennzahl, Tabelle und git-Messung eines Plans in einem Modul

**Status:** Proposed

**Datum:** 2026-10-30

**Autor:** pt9912

**Bezug:** GitHub-Issue 13 (Modul `metrics`) samt den Angaben des Issue-Autors vom
2026-10-30,
[`DC-FA-MET-001`](../../../spec/lastenheft.md#dc-fa-met-001--kennzahl-eines-planungsdokuments-gegen-messung-aus-git-modul-metrics-opt-in),
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[ADR-0024](0024-vcs-immutable-gate.md) und [ADR-0027](0027-commits-traceability-modul.md)
(der VCS-Port und seine Fail-closed-Haltung), [ADR-0072](0072-workflows-modul.md)
(Muster eines Parse-Ports, hier **nicht** gewählt),
slice-282 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-MET-001.a`](../../../spec/spezifikation.md#dc-fa-met-001a--kennzahl-eines-planungsdokuments-gegen-messung-aus-git-metrics)
(neues Modul-Verfahren) und, in der Zuständigkeitsspalte des VCS-Adapters,
[`ARC-005`](../../../spec/architecture.md) und [`ARC-010`](../../../spec/architecture.md)
(Sektion „Komponenten“, Tabellenzellen ohne eigenen Anker).

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v6.18.1/regelwerk/modul-04-adrs.md).

---

## Kontext

Ein Repo führt in seinen Slice-Plänen einen maschinenlesbaren Block mit einer
Größenmarke und prüft ihn mit einem eigenen Programm von rund 5300 Zeilen. Das Issue
verlangt, dass d-check diese Prüfung als opt-in-Modul übernimmt: Obergrenze gegen
Formel, Messung (`git diff --numstat`, Einfügungen plus Löschungen) gegen Obergrenze,
Ausnahme-Eintrag mit Grund. Das Anlass-Repo ist nicht lesbar; die Entscheide stützen
sich auf seine anonymisierten Angaben im Issue, **nicht** auf eine eigene Messung.
Das Programm des Anlass-Repos entfällt nach der Abnahme ersatzlos, es gibt keine
Arbeitsteilung — das Modul trägt die ganze Prüfung samt Tabellenzählung.

Neun Entscheide hat das Issue offengelassen. Drei Randbedingungen des Repos gelten
unabhängig davon: der Kern trägt keine YAML-Bibliothek
([`ARC-002`](../../../spec/architecture.md), durchgesetzt über `.a-check.yml`), git wird
nur über den VCS-Port und ohne git-Binary gelesen, und ein Lauf ohne `metrics`-Konfiguration
bleibt byte-identisch.

## Entscheidung

Wir führen ein **eigenes opt-in-Modul `metrics`** mit diesen Festlegungen (die Nummern
sind die der Entscheide in der Anforderung):

1. **Block:** ein Fenced-Block mit konfiguriertem Info-String und Zeilen-Syntax
   `schluessel: wert` (am ersten `:` getrennt, getrimmt); Pfade und Ausnahmen
   wiederholbar, alles andere höchstens einmal; kein YAML. Plan-Suche: `slice-*.md`
   unmittelbar im konfigurierten Verzeichnis, nicht rekursiv.
2. **Basis:** ein Feld des Blocks (Kurz- oder Voll-Hash), `HEAD` fest. Eine nicht
   auflösbare oder nicht vorfahrende Basis ist ein Eingabefehler.
3. **Formel:** ganzzahlig, zwei Stufen — Schätzung (Summe aus Klassenfaktor mal Zeilen
   und Feldfaktor mal Feld), Obergrenze (genau eine Aufrundung auf ein Vielfaches der
   Einheit) — plus eine Schwellenregel (Pflichtfeld ab Schwelle). Alle Zahlen sind
   Konfiguration.
4. **Messung:** `--no-renames`-Semantik, literale Pfade, rekursive Verzeichnisse,
   Binärdateien zählen 0 mit Hinweis, ausgenommener Teilbaum, nur Eingechecktes.
5. **Fehlerklassen:** ein Eingabefehler eines Plans ist Exit 2 mit Klartext und
   unterdrückt die Befunde dieses Plans; Befunde sind
   `metrics-limit-mismatch`, `metrics-estimate-mismatch`, `metrics-over-limit`,
   `metrics-reason-empty` und `metrics-cut-reason-missing`; Hinweise sind keine Befunde.
6. **Lesen im Kern:** die Zeilen-Syntax wird im Kern gelesen; **kein Parse-Port**.
7. **Fence-Behandlung:** eine auf `metrics` begrenzte Lese-Ausnahme neben der geteilten
   Vorverarbeitung; ein ungeschlossener Fence ist ein Eingabefehler.
8. **Zählen am Adapter:** der VCS-Port erhält **eine** Methode für die Zeilen-Summe
   (Basis, Pfade, ausgenommener Teilbaum); der Kern zählt nicht selbst.
9. **Tabelle im Modul:** genau eine Markdown-Tabelle pro Plan, konfigurierbarer Kopf,
   Klassen mit Faktor; Prüfung 2 („Schätzung gleich Tabelle“, nur bei Messung 0) und der
   Hinweis „Marke ausstehend“ gehören zum Modul.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, das Programm des Anlass-Repos bleibt | kein Aufwand hier | genau der Zustand, den das Issue beklagt: ein eigenes Programm von rund 5300 Zeilen je Repo |
| B — Block als YAML, Parse-Port nach dem Muster von `WorkflowParser` | die Block-Syntax wäre frei erweiterbar | ein neuer Port samt Adapter für ein Format, das nach den Angaben des Anlass-Repos nur flache Zeilen braucht; der Modulkern würde um Port und Adapter größer |
| C — Zeilen im Kern über `AllPaths` und `FileAt` zählen | kein neuer Port-Zuwachs | ein Zeilen-Diff (Binärerkennung, Zeilenenden) im Kern; die `numstat`-Gleichheit wäre dort schwerer zu belegen als am Adapter, der den Diff besitzt |
| D — feste Regel statt konfigurierbarer Formel | einfacher | die Faktoren sind Fachlogik des Repos; ein zweiter Konsument bräuchte Code statt Konfiguration |
| E — Tabellenzählung als spätere, eigene Anforderung | kleineres erstes Modul | der Anlass-Repo-Eigentümer verlangt keine Arbeitsteilung; das Programm könnte nicht entfallen |
| **F — ein Modul mit Zeilen-Syntax im Kern, Zählung am VCS-Port, Tabelle im Modul** | trägt die ganze Prüfung; keine neue Bibliothek im Kern; ein Port-Zuwachs | die Tabelle ist repo-spezifisch und nur über die Konfiguration ausdrückbar; mehr Fläche für Randformen |

## Konsequenzen

- Positiv: das Programm des Anlass-Repos entfällt; ein weiteres Repo mit anderem Block,
  anderer Formel oder anderer Tabelle braucht nur eine Konfiguration.
- Negativ: der VCS-Port wächst um eine Methode, und mit ihr jeder Test-Double des Ports;
  die geteilte Vorverarbeitung erhält eine begrenzte Ausnahme für `metrics`; das Modul
  liest Pfade und git-Revisionen, die es nicht scannt — dort gilt die Zusage der
  Scan-Menge nicht, und die Anforderung benennt die Grenze.
- Folgepflicht: die `numstat`-Gleichheit des Adapters wird am Fixture mit allen Randformen
  belegt, **bevor** diese ADR `Accepted` wird; die Paritäts-Probe liegt beim Eigentümer
  des Anlass-Repos und ist Teil des Wellen-Abschlusses.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check | der Kern importiert keine YAML-Bibliothek und kein git-Paket; die git-Tür bleibt der Adapter | `make arch-check` |
| Spiegel-Test | `model.ValidModules()` stimmt mit Lastenheft, READMEs, Handbuch und `operations.md` überein | `make test` |

## Re-Evaluierungs-Trigger

Die Entscheidung wird neu geprüft, wenn (a) ein zweiter Konsument ein Blockformat
braucht, das die Zeilen-Syntax nicht ausdrückt (dann Parse-Port, Alternative B), (b) die
Fixture-Gleichheit der Zeilen-Summe mit `git diff --numstat` am Adapter nicht
herstellbar ist (dann ist die Messzusage zu korrigieren, vor `Accepted`), oder (c) die
Paritäts-Probe des Anlass-Repos ein abweichendes Urteil je Plan zeigt, das keine der
benannten Randformen erklärt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-30 | Proposed | slice-282 <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
