# Welle welle-74-geteilte-lexik-raender: Dieselbe Klasse, drei andere Lexiken

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-74-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Defekt-Klasse in ausgelieferten
Modulen).

**Verantwortlich:** pt9912. **Datum:** 2026-08-16.

---

## 1. Welle-Ziel

Die Klasse schließen, die [slice-101](slice-101-fence-unbalanciert.md) für
die **Fence**-Lexik geschlossen hat, dort aber ausdrücklich nur für sie:
*eine geteilte Lexik driftet an den Rändern, weil jeder Konsument sie selbst
vorbereitet.* [slice-103](slice-103-geteilte-lexik-raender.md) trägt
die drei Befunde, die in slice-101 nicht hingehörten.

**Das Mehr gegenüber der Slice-DoD ist eine Entscheidung, keine Lieferung.**
Das Beobachtungs-Register führt die Klasse als **BEO-003** mit Zähler **2**
([slice-101](slice-101-fence-unbalanciert.md),
[slice-099](slice-099-structure-modul.md)). Diese Welle bringt sie
entweder zum Abschluss — oder auf **3**, und dann verlangt die Register-Regel
die Verkörperung statt eines weiteren Zählschritts. Diese Welle entscheidet das
ausdrücklich, statt es der nächsten zu überlassen: **BEO-002** und **BEO-004**
sind in welle-73 genau daran gescheitert, zweimal an der Schwelle vorbeigezählt
worden zu sein.

Die Reihenfolge folgt der Schwere, nicht der Slice-Nummerierung: Fall 1 ist eine
**falsche Antwort in einem ausgelieferten Modul** (Exit 0 statt des zugesagten
fail-closed Exit 2), Fall 2 sind **zwei Antworten auf dieselbe Frage** in
einem Lauf, Fall 3 ist eine **benannte Grenze** ohne heutigen Schaden.

## 2. Trigger (Welle startet)

Freigabe des Auftraggebers (2026-08-16), WIP-Slot frei (welle-73 geschlossen,
`in-progress/` trägt nur die Roadmap), und die Start-Bedingung des Slice ist
erfüllt: [slice-101](slice-101-fence-unbalanciert.md) liegt in `done/`.

## 3. Closure-Trigger (Welle schließt)

- [slice-103](slice-103-geteilte-lexik-raender.md) liegt in `done/`
  — oder seine nach Abnahme-Punkt 2 abgespaltenen Folge-Slices sind geschnitten
  und der verbleibende Teil liegt in `done/`.
- **Die Bestandsmessung liegt vor, bevor der Schnitt entschieden wird.** In
  slice-101 hat sie die Entscheidung gedreht; hier entscheidet sie über
  „ein Slice oder drei".
- **BEO-003 ist entschieden** — geschlossen oder auf 3 mit benannter Form.
- Release als **Minor**, falls das Ergebnis Konsumenten erreicht; die Richtung
  („findet mehr" / „findet weniger") steht in der Notiz.
- `make fullbuild` grün.

## 4. Slices in dieser Welle

| Slice | Rolle |
|---|---|
| [slice-103](slice-103-geteilte-lexik-raender.md) | Die drei Lexik-Befunde: `citations`-Absatzbildung, Anker-Auflösung in `headingSection`, git-Revisionen in `vcs` |

Ob es bei einem Slice bleibt, entscheidet Abnahme-Punkt 2 **nach** der Messung.
Die Welle ist so geschnitten, dass eine Abspaltung sie nicht sprengt: ein
Folge-Slice, der aus Abnahme-Punkt 2 entsteht, gehört in **diese** Welle, solange
er dieselbe Klasse trägt.

## 5. Abhängigkeiten

- Die drei Verträge liegen vor:
  [`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in),
  [`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) mit
  [`DC-FA-PIN-001`](../../../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in),
  [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in).
  Diese Welle ändert Antworten, sie erfindet keine Fähigkeit.
- **Die drei Repos für die Bestandsmessung** müssen vorliegen (dieselbe Methode
  wie in slice-101: der reale Bestand, nicht Fixtures).

## 6. Out-of-Scope für diese Welle

- **Die Tabellen-Lexik.** Sie fehlt zwei anderen Kandidaten
  ([slice-102](../in-progress/slice-102-wellen-lifecycle-invariante.md) Aussagen 3/4 und
  der Ordnungs-Bedingung aus **BEO-005**) und gehört in die Welle ihrer beiden
  Konsumenten. — **Die Begründung dafür war falsch:** dieser Punkt nannte sie
  einen „neuen, nicht driftenden“ Rand. Die dritte Review-Runde hat gezeigt, dass
  `targets` Tabellenzeilen **roh** liest und ein Beispiel im Code-Block dadurch
  ein Target als dokumentiert gelten ließ — die Lexik driftet bereits. Der
  Defekt ist in diesem Slice repariert; **out-of-scope bleibt nur der Ausbau**
  zu einer adressierbaren Tabellen-Lexik (Spalten, Typen, Ordnung).
- **Die git-Achse für scannende Module erreichbar machen.** Fall 3 ist LOW;
  möglich, dass die richtige Lieferung eine benannte Grenze ist und kein Umbau.
- [slice-095](../open/slice-095-links-resolve-from.md).

## 7. Closure-Notiz

Geschlossen am 2026-08-16 mit **v0.58.0**. Alle fünf Closure-Trigger sind
erfüllt: [slice-103](slice-103-geteilte-lexik-raender.md) liegt in `done/`, die
Bestandsmessung lag **vor** dem Schnitt-Entscheid, **BEO-003 ist entschieden**
(auf 3 und als Kopplungs-Test verkörpert), das Release samt Digest-Backfill ist
draußen, und `make fullbuild` ist grün.

Die vollständige Notiz — geliefert, gelernt, Register-Lese-Schritt und
Trigger-Audit — steht in [`welle-74-results.md`](welle-74-results.md).
