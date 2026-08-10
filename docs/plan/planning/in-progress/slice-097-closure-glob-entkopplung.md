# Slice slice-097: `planning.closure.glob` — eigener Kandidaten-Filter für die Closure-Fähigkeit

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** gemeinsam mit [slice-098](../open/slice-098-closure-note-placeholder.md) —
der Konsument kann sein letztes handgeschriebenes Prüfskript **erst dann**
zurückziehen, wenn **beide** Slices liegen. Das ist eine Closure-Bedingung
jenseits der beiden Slice-DoDs und damit genau der Fall, für den es eine Welle
gibt. **Zuordnung entschieden** mit der welle-69-Closure (2026-08-09): der Slice
bleibt **eigenständig** und geht nicht im `structure`-Modul auf.

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Closure-Fähigkeit, Spezifikation Schritt C2),
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md).
**Change Request** des Konsumenten `ai-harness-course` (CR 1).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Die Closure-Fähigkeit bekommt einen **eigenen** Kandidaten-Filter
(`planning.closure.glob`, Basisname-Glob, Default = `planning.slice-glob`).
Heute teilen sich zwei Fähigkeiten mit **verschiedenen Grundmengen** einen Knopf;
er ist nicht einstellbar, ohne eine von beiden zu verbiegen.

## 2. Der gemessene Befund

Die Aussagen sind verschieden: die Roadmap-Invariante fragt „liegt hier noch
Arbeit?", die Closure-Fähigkeit „ist jedes abgeschlossene Paket dokumentiert?".
Beide lesen heute `planning.slice-glob`.

Nachgemessen am 2026-08-09 gegen v0.52.0, Endzustand *kein Slice mehr in
Arbeit, Ruhe-Marker gesetzt* (muss grün sein):

```text
slice-glob "slice-*.md"  → 0 Befunde                    ✓
slice-glob "*.md"        → roadmap.md  planning-drift    ✗ falsch-rot
```

Ursache: die Roadmap-Datei liegt selbst im gezählten Verzeichnis und matcht
`*.md`. Die Invariante sieht dauerhaft „eine Datei liegt da" und verlangt, dass
der Ruhe-Marker **nie** gesetzt wird — falsch genau dann, wenn die Welle
wirklich ruht.

**Das ist ein Entwurfsfehler aus
[slice-093](../done/slice-093-closure-note-gate.md)**, nicht eine fehlende
Fähigkeit: die zweite Fähigkeit hat den Glob der ersten mitbenutzt, weil beide
zufällig Slice-Dateien meinten. Sobald eine Grundmenge weiter ist, bricht die
Kopplung.

## 3. Abnahme-Punkte

1. **Default = `planning.slice-glob`**, nicht ein eigener Literal-Default. Nur so
   ist der Befundsatz ohne den Schlüssel byte-identisch zu v0.52.0.
2. **Der eigene Bestand ist mitbetroffen — was tun wir damit?** In `done/` liegen
   **9** `welle-*-results.md`, jede mit einer Closure-Notiz; sie sind heute
   unsichtbar. Weitet man den Glob, greift sofort der nächste Unterschied: ihre
   Notiz-Überschrift ist eine **H1** (`# Welle NN — … — Closure-Notiz`), das
   Default-`heading-pattern` verlangt H2/H3 ⇒ **9** `closure-note-missing`
   (nachgemessen 2026-08-09). Die Wellen-**Plan**-Dateien sind nicht darunter —
   ihre Notiz ist H2. Ein zehnter Befund der Messung, ein `planning-drift`,
   entsteht **nur** beim Weiten von `slice-glob`; mit dem hier beantragten
   eigenen Schlüssel bleibt die Roadmap-Invariante unberührt, und genau das ist
   der Punkt.
   **Entschieden 2026-08-10 — und zwar anders, als die drei Kandidaten es
   anlegten.** Gemessen wurde eine **vierte** Variante, die erst mit dem neuen
   Schlüssel möglich ist: `closure.glob: "*.md"` **und** `heading-pattern` auf
   `^#{1,3}`.

   | Variante | Befunde | geprüfte Dokumente |
   |---|---|---|
   | heute (Filter erbt `slice-glob`) | 0 | 96 |
   | Filter `*.md` | 12 | +15 |
   | Filter `*.md` **und** Muster `^#{1,3}` | **0** | **+15** |

   Elf der zwölf sind kein fehlender Abschnitt, sondern eine **H1** gegen ein
   Muster, das H2/H3 verlangt. Der zwölfte war **echt**: ein geschlossenes
   Wellendokument trug in `done/` noch „_Ausstehend._". Ein zweites entging dem
   Befund nur, weil sein „_Ausstehend._" wortreich genug für die Satz-Schwelle
   war. Beide sind gefüllt.

   „Eng lassen" wäre damit genau die Grenze, die das Register als **BEO-004**
   führt — eine Prüfung, die an Dokumenten vorbeisieht, die sie meint. Die
   Zahlen sind seit der Slice-Anlage gewachsen (9 → 11 `results`-Dateien): der
   Bestand wächst weiter, die Lücke also auch.

## 4. Definition of Done

- [x] `planning.closure.glob` in Lastenheft (0.53.0), Spezifikation (Schritt C2)
      und Config-Schema; Default = **Verweis** auf `planning.slice-glob`, kein
      kopiertes Literal. Begründung in
      [ADR-0051](../../adr/0051-eigener-kandidaten-filter-closure.md) `Proposed`;
      [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) §Geschichte
      um den Verfeinerungs-Zeiger ergänzt (ihr Rumpf ist `Accepted` und damit
      immutabel).
- [x] Akzeptanzkriterien des CR als Tests, alle drei plus der Default-Verweis.
      **Fünf** Rückbauten geprüft, alle rot: Filter zurück auf `slice-glob`;
      Default als Literal statt Verweis; leerer Glob fällt still zurück;
      Glob-Validierung entfernt; Wert nicht durchgereicht. End-to-End gegen das
      gebaute Image: eigener Lauf 326 Dateien / 0 Befunde, expliziter leerer
      Glob bricht mit **Exit 2** und einer Meldung ab, die den Schlüssel nennt.
- [x] `make gates` grün; Abnahme-Punkt 2 beantwortet und `.d-check.closure.yml`
      entsprechend gesetzt — der eigene Lauf prüft jetzt **111 statt 96**
      Dokumente und bleibt bei null Befunden.
- [ ] **Release** (Minor: neuer Config-Schlüssel), Digest-Backfill. Ohne
      veröffentlichte Version erreicht die Welle ihren Zweck nicht — der
      Konsument kann sein Skript erst gegen ein Release zurückziehen
      (Schnitt-Review F-6).

## 5. Risiken / offene Punkte

- **Neun eigene Befunde bei geweitetem Glob** (siehe Abnahme-Punkt 2).
  — **Ausgang: entfallen.** Es waren inzwischen zwölf, aber elf davon lösen sich
  mit dem Überschriften-Muster auf; der zwölfte war ein echter Rückstand und ist
  behoben. Eine Sanierung der Wellen-Notizen war nicht nötig.
- **Zwei Globs, die fast immer gleich sind,** laden zu Drift ein (einer wird
  gepflegt, der andere nicht). — **Ausgang: adressiert und geprüft.** Der Default
  ist ein **Verweis**, kein Literal: wer nichts trennt, pflegt genau ein Muster.
  Die Gegenprobe „Default als Literal statt Verweis" ist rot — die Eigenschaft
  ist testgehalten, nicht bloß zugesagt. Es bleibt der Restrisiko-Fall, dass
  jemand `slice-glob` ändert und die Closure-Menge unbeabsichtigt mitwandert;
  das steht als Re-Evaluierungs-Trigger in
  [ADR-0051](../../adr/0051-eigener-kandidaten-filter-closure.md).

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei — **und**
[slice-096](../done/slice-096-structure-modul-analyse.md) in `done/`, weil dessen
Schnitt über den Zuschnitt dieses Slice mitentscheidet.

**Rückführungen:** `in-progress` → `open`, falls Abnahme-Punkt 2 eine
Bestands-Sanierung nach sich zieht, die eigenständig geschnitten gehört.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten** (bei Slice-Beginn erneut gelesen — das
  Register führt inzwischen **vier** Einträge, nicht mehr nur BEO-001):
  - **BEO-001** (Datei-Register driften gegen ihre Autoritäts-Tabelle): andere
    Klasse, nichts zu berücksichtigen.
  - **BEO-003** (geteilte Lexik driftet an den Rändern, weil jeder Konsument sie
    selbst vorbereitet): **einschlägig als Warnung**. Ein kopierter
    Literal-Default wäre genau diese Klasse gewesen — deshalb der Verweis.
  - **BEO-004** (Modul-Grenze nur auf der Quell-Achse gedacht): **einschlägig
    als Frage.** Die Closure-Fähigkeit ist ein Post-Pass über ein selbst
    benanntes Verzeichnis; dieser Slice ändert genau, welche Dateien darin sie
    liest. Die Register-Frage „welche Eingaben liest dieses Modul, die es nicht
    scannt?" hat Abnahme-Punkt 2 beantwortet — und dabei 15 übersehene
    Dokumente gefunden.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Change Request und Spezifikations-Schärfung
schreiben die Zusage, der Go-Code liefert sie.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
