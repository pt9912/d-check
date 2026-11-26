# ADR-0117: Die neuen Wellen-Namen laufen über die Konfiguration dieses Repos — die Bestands-Wellen bleiben ungewächtert

**Status:** Proposed

**Datum:** 2026-11-26

**Autor:** pt9912

**Bezug:** [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(das Modul `planning`, `planning.waves`),
[ADR-0028](0028-planning-lifecycle-modul.md) (das Modul),
[ADR-0055](0055-wellen-invariante-artefakt-und-grund-codes.md) (Wellen-Invariante `wave-drift` und
Register-Deckung), [`AGENTS.md`](../../../AGENTS.md) §3.6 (Senkung einer Prüfregel nur per ADR),
slice-298 <!-- d-check:status-provenance -->.

**Schärft:** — (das Produkt, seine Defaults und die Spezifikation von `planning.waves` bleiben; entschieden
wird die Konfiguration **dieses** Repos).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Das Baseline-Regelwerk benennt seit `v7.0.0` die Dateien einer Welle englisch (`wave-<Kennung>.md`,
`wave-<Kennung>-results.md`, `archive.zip`, Anker `seit wave-<Kennung>`). Der Auftraggeber hat am 2026-11-25
entschieden: **ab jetzt für Neues, der Bestand bleibt.** Die Schreibweisen stehen damit nebeneinander.

Das Modul `planning` liest **ein** literales Präfix: `planning.waves.glob` und `results-glob` tragen dasselbe
Präfix, und das Modul leitet die Kennung daraus ab (`wavePrefix` in `planning_waves.go`; die Prüfung der
Konfiguration verlangt Gleichheit). Gemessen in einem Wegwerf-Klon, Stand `590dde3d`:

- Mit `welle-*` (heute) ist eine neue `wave-<NN>-….md` **unsichtbar**: `make planning-check` Exit 0, 0 Befunde,
  auch ohne Zeiger in der Roadmap.
- Mit `wave-*` meldet es `wave-drift` für die Datei ohne Zeiger; die **38** Abschluss-Register-Zeilen
  `welle-<NN>` und ihre **38** Ergebnisnotizen sind dann **ungewächtert** — kein Befund, auch wo einer stünde.
- Beide Bevölkerungen zugleich kann nur eine Änderung am Produkt (das Modul liest eine Liste von Präfixen).

## Entscheidung

1. Dieses Repo stellt `planning.waves.glob` auf `wave-*.md` und `results-glob` auf `wave-*-results.md` um. Das
   ist **Konfiguration**; die Defaults des Produkts (`welle-*`) bleiben, ebenso Lastenheft und Spezifikation.
2. **Die Senkung ist benannt:** Für die Bestands-Wellen (`welle-<NN>`) laufen die Wächter des Moduls `planning`
   (Wellen-Invariante, Register-Deckung) nicht mehr. Das gilt **nur** für `planning`; `matrix`, das
   Closure-Profil und das Archiv-Werkzeug bekommen die neue Form **neben** die alte
   (slice-298 <!-- d-check:status-provenance -->, Probe in einer Wegwerf-Kopie).
3. Die Umstellung steht an der Konfiguration mit ihrer Grenze im Kommentar, und der Eintrag zur Koexistenz
   in `harness/conventions/` verweist auf diese ADR, statt die Senkung nachzuerzählen.
4. Die größere Antwort — das Modul `planning` liest mehrere Präfixe — ist ein **Produkt-Merkmal** (Anforderung,
   Spezifikation, ADR, Handbuch, Release) und gehört dem Auftraggeber. Sie löst die Senkung auf und ist **nicht**
   Teil dieser Entscheidung.

## Verglichene Alternativen

- **A. Das Modul liest mehrere Präfixe.** Löst die Senkung, kostet ein Produkt-Merkmal an acht Stellen von
  `wavePrefix`, der Konfigurations-Prüfung und rund 100 Test-Zeilen, dazu Lastenheft, Spezifikation, ADR,
  Handbuch und ein Release. Verworfen für jetzt: der Anlass (neue Namen) trägt es nicht, die Entscheidung liegt
  beim Auftraggeber.
- **B. Den Bestand umbenennen.** Dann gäbe es keine Koexistenz. Verworfen: 345 Pfade mit `welle-<NN>`, 188
  `archiv.zip`, 19 Anker und 38 Register-Zeilen sind Belege (Kommandos und Zahlen:
  [`MR-082`](../../../harness/conventions/MR-082-wave-welle-koexistenz.md) §Begründung und §Auflösungs-Trigger);
  der Auftraggeber hat „Bestand bleibt“ entschieden.
- **C. Die Konfiguration bleibt auf `welle-*`.** Keine Senkung, aber eine neue Welle ist für das Gate
  unsichtbar — ein grünes `planning-check` über einer Menge, die es nicht sieht. Verworfen: das ist die
  schlechtere Senkung, und sie wäre unbenannt.

## Konsequenzen

- Der Bestand ist ungewächtert; ein Defekt in einer geschlossenen Bestands-Welle bleibt für `planning` stumm.
  Der Preis ist klein und fest (die Menge wächst nicht mehr), die Bestands-Wellen sind geschlossen.
- Jede neue Welle ist sichtbar, `make planning-check` bleibt in `make gates`.
- Zwei Schreibweisen stehen nebeneinander, bis der Bestand umbenannt ist (Eintrag zur Koexistenz).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Probe in einer Wegwerf-Kopie des Repos | flache `wave-<NN>-probe.md` ohne Zeiger ⇒ `wave-drift`, mit Zeiger Exit 0 | `make planning-check` |

## Re-Evaluierungs-Trigger

- Der Auftraggeber entscheidet das Produkt-Merkmal „mehrere Präfixe“ — dann entfällt die Senkung (Alternative A).
- Der Bestand wird umbenannt — dann entfällt die Senkung mit der Koexistenz.
- In einer Bestands-Welle wird ein Defekt gefunden, den `planning` gemeldet hätte — dann ist A vorzuziehen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-11-26 | Proposed: Architect-Verdikt zur Hebung der Baseline, Frage 3 | slice-298 <!-- d-check:status-provenance --> |
