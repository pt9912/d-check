# Realdatenbeleg slice-071 — `trace.cross-consistency` gegen grid-gym

**Datum:** 2026-07-18. **Slice:** slice-071 (Kreuzverweis-Konsistenz-Gate).
**Prüfling:** d-check v0.48.1 (`d-check:latest` aus HEAD `9de80f6`), Kontrast gegen
v0.48.0. **Gegenstand:** der von ADR-0038 Entscheidung 7 geforderte Lauf gegen die
realen Quellen des Konsumenten grid-gym — bis v0.48.0 an `spec/architecture.md:913`
mit Exit 2 blockiert, entblockt durch die Direktiven-Toleranz aus slice-074.

## 1. Aufbau

Gefahren gegen eine read-only-Arbeitskopie des grid-gym-Repos (Schwester-Repo,
außerhalb dieses Repos; die Kopie wurde nicht zurückgeschrieben — im Original läuft
eine parallele Sitzung). Aufruf wie im dortigen `d-check.mk`:
`docker run --rm --network none -v <kopie>:/repo:ro <image> --trace`.

Config: der **Einmal-Messblock** aus grid-gyms Trigger 088 (`trace.cross-consistency`,
`mode: equal`, Vorwärts-Sicht `docs/plan/traceability.md` §27.1, Rück-Sicht
`spec/architecture.md`, `exclude-req` für Ports/Prinzipien/Meta-Familien), temporär
unter `trace:` eingehängt. Der einzige bewusste Zusatz gegenüber der dort
dokumentierten Config: `backward.sections` um `"17. Testarchitektur"` **erweitert** —
genau der Abschnitt, den der dortige Messlauf mit v0.45.1 weglassen musste, weil
d-check an Zeile 913 abbrach.

## 2. Läufe

| Lauf | Image | `17. Testarchitektur` in `backward.sections` | Exit | Ausgabe |
|---|---|---|---|---|
| Baseline | v0.48.1 | nein | 0 | `161 Differenz(en).` |
| §17 | **v0.48.0** | ja | **2** | `d-check: error: trace.cross-consistency.backward: Tabellenzeile 913 hat 4 statt 3 Zellen` |
| §17 | **v0.48.1** | ja | 0 | `161 Differenz(en).` |

- **Baseline** reproduziert exakt die in Trigger 088 dokumentierte Messung
  (`161 Differenzen = 86 F\B + 75 B\F`) — der Aufbau ist damit als treu belegt.
- **v0.48.0 mit §17** bricht an genau der dokumentierten Stelle ab: die
  E2E-Datenzeile trägt d-checks eigene Ignore-Direktive
  `<!-- d-check:ignore … -->` als überzählige Zelle, die vor slice-074 als 4. Zelle
  gezählt wurde.
- **v0.48.1 mit §17** liest den Abschnitt durch und liefert dasselbe Ergebnis wie
  die Baseline. Die §17-Tabelle hat keine `GG-AR-*`-Erstspalten-ID, trägt also
  **keine** Kante zum Mengenabgleich bei; die Differenzausgabe ist Byte für Byte
  identisch zur Baseline (161 = 161). Allein die **Lesbarkeit** des Abschnitts war
  der Blocker — und genau die stellt slice-074 her.

Damit ist der einzige veränderte Faktor zwischen dem v0.48.0-Abbruch und dem
v0.48.1-Durchlauf die Direktiven-Toleranz an Zeile 913.

## 3. Inhaltliche Fitness (ADR-0038-Fitnessfunktion)

Aus der `## Kreuzverweis-Konsistenz`-Ausgabe des v0.48.1-Laufs:

- **Realer Drift geflaggt (Schnittmenge null):** der Trigger-088-Fall `GG-ARCH-005`
  erscheint mit beiden Richtungen — §27.1 nennt `GG-AR-COMP-CORE`/`-COMP-DOMAIN`
  (`in RTM, ohne Rück-Kante`), die Rück-Kanten `GG-AR-COMP-SCHED`/`-P-005`/`-P-009`
  (`Rück-Kante, ohne RTM-Eintrag`). Ebenso die zweite echte Drift `GG-SIM-009`.
- **Ventil greift:** aus den per `exclude-req` ausgeschlossenen Familien
  (`GG-PRINC`/`GG-CC`/`GG-QA`/`GG-QG`/`GG-COV`/`GG-TESTTYPE`/`GG-ARCHTEST`) erscheint
  **kein** Befund (0 Treffer).
- **1:N grün:** die große Mehrheit der Anforderungen erzeugt keine Differenz; jede
  gemeldete Zeile trägt Richtungslabel und `Datei:Zeile` (161 Zeilen, deterministisch
  sortiert).

## 4. Ein Konsumenten-Befund am Rande (kein d-check-Defekt)

Der `--trace`-Lauf brach zunächst **vor** dem Abgleich an `trace.coverage` ab:
`Komma-Kurzform hinter GG-SCN-001 ist keine zugesagte Notation`. Ursache ist die
Zelle `GG-SCN-001..005, 007, 008` in §27.1 — grid-gym maß mit v0.45.1, wo die
Komma-Aufzählung noch **still verschluckt** wurde; slice-075 (v0.46.0) lehnt sie seit
dem korrekt **fail-closed** ab (genau der als Werkzeug-Artefakt an d-check gemeldete
Punkt aus Trigger 088). Für den isolierten §17-Beleg wurde diese **eine** Zelle in
der Kopie auf die zugesagte Notation `GG-SCN-001..005/007/008` normalisiert — eine
Nacharbeit, die beim Konsumenten ohnehin ansteht, und die den slice-074-Beleg nicht
berührt (beide §17-Läufe nutzten die normalisierte Kopie, sodass Zeile 913 der
einzige Unterschied bleibt).

## 5. Ergebnis

Der Realdatenbeleg ist erbracht: die von slice-074 gelieferte Toleranz ist der exakte
Mechanismus, der den §17-Lauf gegen die realen grid-gym-Quellen entblockt; die
dokumentierte 161-Differenzen-Messung reproduziert auf aktuellem d-check; realer
Drift wird geflaggt, das Ventil greift, das konsistente 1:N läuft grün. Der letzte
offene DoD-Punkt von slice-071 ist damit erfüllt; welle-60 ist abschlussreif.
