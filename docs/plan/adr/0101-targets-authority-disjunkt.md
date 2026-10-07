# ADR-0101: `targets` prüft opt-in die Disjunktheit der Autoritäts-Dateien (`gate-declared-twice`)

**Status:** Proposed

**Datum:** 2026-10-07

**Autor:** pt9912

**Bezug:** [`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(erweitert, Lastenheft 0.97.0); löst den Re-Evaluierungs-Trigger von
[ADR-0100](0100-targets-authority-liste.md) ein (die Baseline nimmt die
Disjunktheits-Regel an — gemeldet im
[Hinweis der Baseline](../cr/2026-10-07-hinweis-eingehend-ai-harness-course-disjunktheit.md));
Schnitt-Kriterium aus [ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md);
slice-252 <!-- d-check:status-provenance -->.

**Schärft:**
[`spec/spezifikation.md` §DC-FA-TGT-001.a](../../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
(Schritt 5a), die Schema-Zeile `targets.authority-disjoint` unter
[`SPEC-005`](../../../spec/spezifikation.md#spec-005--d-checkyml) und die
Grund-Code-Zeile [`SPEC-088`](../../../spec/spezifikation.md#4-grund--und-fehler-codes).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Seit ADR-0100 misst `gate-undocumented` gegen die Vereinigung mehrerer
Autoritäts-Dateien. Eine Vereinigung zählt ein doppelt genanntes Target
einmal; der Sensor bleibt dabei still. ADR-0100 hielt die Disjunktheit
bewusst offen, weil die Regel nur in einem offenen CR an die Baseline stand.
Seit `v6.16.0` führt die Baseline sie: „Kein Target steht in zwei Teilen",
mit dem Zusatz, dass die Disjunktheit eine eigene Prüfung braucht und ohne
sie als benannte Grenze geführt wird.

## Entscheidung

1. **Opt-in über `targets.authority-disjoint: true`.** Die Regel gilt ab einer
   Baseline-Version; Adopter älterer Baselines sollen nicht rot werden. Ohne
   den Schalter ist der Befundsatz byte-identisch.
2. **Eigener Grund-Code `gate-declared-twice`** — die Verletzung ist weder ein
   Phantom noch ein undokumentiertes Gate, sondern eine dritte Aussage über
   den Index.
3. **Die erste Nennung in Konfigurations-Reihenfolge ist die Heimat; gemeldet
   wird jedes Vorkommen in einer späteren Datei.** Das macht die Fundstelle
   deterministisch und legt sie dorthin, wo sie in der Regel behoben wird:
   der Teil des Werkzeugs steht hinter dem des Repos, und das Werkzeug
   regeneriert ihn.
4. **`exempt-targets` wirkt nicht.** Die Ausnahme befreit von der
   Doku-Pflicht, nicht von der Regel, dass zwei Zeilen dasselbe Target
   auseinanderlaufen lassen.
5. **Doppelungen innerhalb einer Datei sind kein Fall** — die Regel spricht
   von zwei Teilen. Dieselbe Datei in zwei Schreibweisen ist ein Teil.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Default-an** | jede Doppelung sofort sichtbar | macht Adopter älterer Baselines rot ohne Grundlage in ihrer Konvention |
| **Befund an jeder Fundstelle, auch der ersten** | symmetrisch | ein Fehler, zwei Befunde; die Behebung findet nur in einer Datei statt |
| **`exempt-targets` respektieren** | eine Ausnahme-Achse weniger zu erklären | ein ausgenommenes Target in zwei Teilen läuft genauso auseinander |
| **opt-in, Heimat = erste Nennung, Befund an späteren** (gewählt) | deterministisch, ein Befund je überzähliger Zeile, rückwärtskompatibel | wer die Reihenfolge der Liste umstellt, verschiebt die Fundstelle |

## Konsequenzen

- Der Kurs kann seine Probe der heutigen Stille (Team-Sim s30d) auf „laut"
  drehen; die benannte Grenze wird zum opt-in-Sensor.
- Die Beobachtung zur ungeprüften Disjunktheit bekommt ihren Ausgang.
- Schema, §4, `AllReasons` und `--doctor` kennen den neuen Grund-Code.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test `TestCheckTargetsAuthorityDisjunkt` | Befund an späteren Dateien, Heimat in der Meldung, auch für ein exemptes Target, drei Dateien | `make test` |
| Go-Test `TestCheckTargetsAuthorityDisjunktGrenzen` | ohne Schalter, eine Datei, gleiche Datei zweifach ⇒ still | `make test` |
| Go-Test `TestDecode_TargetsAuthorityDisjoint` | Schalter durchgereicht, Default aus, kein Bool ⇒ Fehler | `make test` |
| Go-Test `TestAllReasonsDeckungGegenSpezifikationGrundCodes` | Grund-Code in `AllReasons` und §4 | `make test` |

## Re-Evaluierungs-Trigger

Die Baseline macht die Disjunktheit zur Pflicht für jedes Repo mit mehreren
Index-Teilen — dann ist Default-an zu prüfen. Ohne das: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-07 | Angelegt als `Proposed`; `Accepted` erst mit der Closure des Vorgangs, nach Review und Verifikation |
