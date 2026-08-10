# ADR-0051 — Eigener Kandidaten-Filter für die Closure-Fähigkeit, Default als Verweis

**Status:** Proposed
**Datum:** 2026-08-10
**Autor:** pt9912
**Bezug:** [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Closure-Fähigkeit), [§`DC-FA-PLAN-001.a`](../../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
Schritt C2. **Verfeinert** [ADR-0048](0048-closure-note-struktur-im-planning-modul.md)
Entscheidung 1 — dort trug die Annahme, beide Fähigkeiten teilten eine
Config-Achse, die Begründung gegen ein zweites Modul.
**Change Request** des Konsumenten `ai-harness-course` (CR 1).

## Kontext

[ADR-0048](0048-closure-note-struktur-im-planning-modul.md) hat die
Closure-Struktur als **zweite Fähigkeit** in `planning` untergebracht statt als
eigenes Modul. Eine der tragenden Begründungen: ein zweites Modul hätte
„dieselbe Config-Achse (Slice-Verzeichnis, `slice-glob`) ein zweites Mal
deklariert“.

Die Achse ist aber **nicht** dieselbe. Die beiden Fähigkeiten stellen
verschiedene Fragen:

| Fähigkeit | Frage | Grundmenge |
|---|---|---|
| Lifecycle-Invariante | liegt hier noch **Arbeit**? | Slice-Dateien im Roadmap-Verzeichnis |
| Closure-Struktur | ist jedes **abgeschlossene** Paket dokumentiert? | alles, was den Lifecycle verlässt — auch Wellen- und Etappen-Dokumente |

Solange beide zufällig dasselbe trafen, fiel die Kopplung nicht auf. Sobald eine
Menge weiter sein soll, bricht sie. Gemessen gegen v0.53.0 im Zustand *kein
Slice in Arbeit, Ruhe-Marker gesetzt* (muss grün sein):

```text
slice-glob "slice-*.md"  → 0 Befunde                   ✓
slice-glob "*.md"        → roadmap.md  planning-drift   ✗ falsch-rot
```

Die Roadmap-Datei liegt selbst im gezählten Verzeichnis. Die Invariante sieht
dauerhaft „eine Datei liegt da“ und verlangt, dass der Ruhe-Marker **nie**
gesetzt wird — falsch genau dann, wenn die Welle wirklich ruht.

**Bestandsmessung am eigenen Repo (2026-08-10).** Im Closure-Verzeichnis liegen
110 Markdown-Dateien, von denen der heutige Filter 96 sieht. Ungesehen bleiben
11 `welle-*-results.md` und 4 Wellen-Plan-Dokumente:

| Variante | Befunde | geprüfte Dokumente |
|---|---|---|
| heute (Filter erbt `slice-glob`) | 0 | 96 |
| Filter `*.md` | 12 | +15 |
| Filter `*.md` **und** `heading-pattern` auf `^#{1,3}` | **0** | **+15** |

Elf der zwölf Befunde der mittleren Zeile sind kein fehlender Abschnitt, sondern
eine **H1** gegen ein Default-Muster, das H2/H3 verlangt. Der zwölfte war echt:
ein geschlossenes Wellendokument trug in seiner Closure-Notiz noch
„_Ausstehend._“ — die Messung hat einen realen Rückstand gefunden, bevor
der Schlüssel überhaupt existierte.

## Entscheidung

1. **Eigener Schlüssel `planning.closure.glob`** statt Weiten von
   `planning.slice-glob`. Zwei Fragen mit zwei Grundmengen bekommen zwei Filter.
   Das ist keine Verdopplung, sondern die Auflösung einer Kopplung, die nie
   Absicht war: [ADR-0048](0048-closure-note-struktur-im-planning-modul.md) hat
   den Filter der ersten Fähigkeit mitbenutzt, weil beide **zufällig** dasselbe
   trafen.

2. **Der Default ist ein Verweis, kein Literal.** Nicht gesetzt ⇒ es gilt
   `planning.slice-glob`. Ein kopierter Literal-Default (`slice-*.md`) wäre eine
   zweite Pflegestelle, die still auseinanderläuft — genau die Klasse, die das
   Beobachtungs-Register als **BEO-003** führt. Der Verweis hält beide Mengen
   zusammen, solange niemand sie ausdrücklich trennt, und garantiert nebenbei,
   dass der Befundsatz ohne den Schlüssel **byte-identisch** bleibt.

3. **Explizit leer oder ungültig ⇒ Exit 2**, kein stiller Rückfall auf den
   Default. Den Schlüssel zu setzen ist eine Aussage; eine Aussage, die nichts
   trifft, ist ein Konfigurationsfehler. Dieselbe Begründung wie beim
   **explizit** gesetzten `min-sentences` < 1. Ein *nicht* gesetzter Schlüssel
   ist dagegen keine Aussage und fällt lautlos auf den Verweis zurück — die
   Unterscheidung verlangt, dass die Konfigurations-Schicht *gesetzt* von
   *abwesend* trennt.

4. **Kein neuer Grund-Code.** Es ändert sich die **Kandidaten-Menge**, nicht die
   Aussage über einen Kandidaten. Die drei bestehenden Codes gelten unverändert,
   einschließlich der Nullmengen-Härte: kein Kandidat unter dem gesetzten Glob ⇒
   `closure-note-missing` auf dem Verzeichnis.

5. **Der eigene Bestand wird mitgeweitet** — `closure.glob: "*.md"` **und**
   `heading-pattern` auf `^#{1,3}`. Gemessen null Befunde bei 15 zusätzlich
   geprüften Dokumenten. Das ist die bessere Antwort als „eng lassen“: die
   Wellen-Ergebnisnotizen sind Closure-Notizen im vollen Sinn, und sie ungeprüft
   zu lassen wäre genau die Grenze, die das Register als **BEO-004** führt —
   eine Prüfung, die an Dokumenten vorbeisieht, die sie meint.

6. **SemVer: Minor.** Neuer Config-Schlüssel, rein additiv; ohne ihn ist der
   Befundsatz byte-identisch. Ein Adopter merkt beim Update nichts.

## Verglichene Alternativen

| Alternative | Warum verworfen |
|---|---|
| `planning.slice-glob` weiten | Verbiegt die Lifecycle-Invariante: die Roadmap-Datei matcht mit und der Ruhe-Marker meldet dauerhaft falsch-rot — gemessen |
| Literal-Default `slice-*.md` für den neuen Schlüssel | Zweite Pflegestelle, die still driftet (BEO-003); der Verweis hält beide Mengen ohne Pflegeaufwand zusammen |
| Explizit leerer Glob fällt auf den Default zurück | Übergeht stillschweigend eine gesetzte Aussage — dieselbe Klasse stiller Grün-Pfade, gegen die die Nullmengen-Härte gebaut wurde |
| Eigener Grund-Code für „Kandidat durch Glob ausgeschlossen“ | Es gibt nichts zu melden: ein nicht gematchter Kandidat ist keine Verletzung, sondern nicht Gegenstand |
| Closure-Verzeichnis rekursiv statt per Glob | Andere Achse (Tiefe statt Namensform) und eine viel größere Zusage; der CR verlangt sie nicht |
| Eigenen Bestand eng lassen (Schlüssel nur ausliefern) | Ließe 15 echte Closure-Notizen dauerhaft ungeprüft — und die Messung hat in genau diesen 15 einen realen Rückstand gefunden |

## Konsequenzen

- **Die Begründung von [ADR-0048](0048-closure-note-struktur-im-planning-modul.md)
  Entscheidung 1 verliert eines ihrer Argumente**, nicht aber ihr Ergebnis: dass
  die Config-Achse geteilt sei, stimmt nicht. Dass die **Invariante** dieselbe
  ist — Eintritt und Austritt desselben Lifecycles —, trägt die Entscheidung
  gegen ein zweites Modul weiterhin allein.
- **Zwei Globs laden zu Drift ein.** Der Verweis-Default ist die Gegenmaßnahme:
  wer nichts trennt, pflegt nur ein Muster. Wer trennt, tut es absichtlich.
- **Der eigene Lauf prüft 15 Dokumente mehr** und bleibt grün. Jede künftige
  Abweichung dort ist ein echter Fund.

## Fitness Function

- **Ohne den Schlüssel byte-identisch** zum Stand vor dieser Änderung.
- **`closure.glob` weiten berührt die Roadmap-Invariante nicht** — belegt an
  einem Lauf mit geweitetem `closure.glob` und unverändertem `slice-glob`.
- **Explizit leerer Glob bricht mit Exit 2 ab**, nicht mit dem Default.
- **Der eigene Bestand bleibt bei null**, jetzt über 111 statt 96 Dokumente.

## Re-Evaluierungs-Trigger

- Wenn eine dritte Fähigkeit in `planning` eine vierte Grundmenge braucht, ist
  die Filter-pro-Fähigkeit-Form neu zu bewerten — dann ist die Achse selbst
  falsch geschnitten.
- Wenn der Verweis-Default in der Praxis überrascht (jemand ändert
  `slice-glob` und wundert sich über die Closure-Menge), ist ein Literal-Default
  mit Warnung neu abzuwägen.
- Wenn das Closure-Verzeichnis in Unterordner wächst, wird die Tiefen-Frage
  akut, die diese Entscheidung ausdrücklich offen lässt.

## Geschichte

- 2026-08-10: Proposed (doc-first, `slice-097`).
