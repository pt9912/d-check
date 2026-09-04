# Review slice-160 — Die Messung maß eine andere Menge, als sie aussagte

**Gegenstand:** [slice-160](../plan/planning/done/slice-160-reviewer-skill-kontextlast.md), Stand des Feat-Commits vor der Nacharbeit.
**Datum:** 2026-08-27. **Reviewer:** unabhängiger Subagent.

---

## Urteil

**Nicht schließbar.** Vier HIGH. Der Kern: die Messung, die den ganzen Entscheid
trägt, hat **eine andere Menge gemessen als die, über die sie redet**, und ihre
beiden Zahlen stammen aus **disjunkten** Mengen. Der daraus gezogene Satz stand
im Reviewer-Skill — also in dem Dokument, nach dem künftige Reviews arbeiten.

## Befunde

| # | Rang | Befund |
|---|---|---|
| F-1 | **HIGH** | *„Die Reports zitieren keine Anker-Namen"* ist am Bestand widerlegt: **sieben** Reports nennen den Anker wörtlich, **fünf** davon im `quelle`-Feld eines Findings. Die stabile Kennung ist ohnehin die `BEO-`ID |
| F-2 | **HIGH** | Die *„687 Befund-Zeilen der 214 Reports"* sind weder 214 Reports noch nur Befund-Zeilen: nur **104** Reports führen ein `befund:`-Feld, das `quelle`-Feld war nicht im Scan, **307** Zeilen datieren vor der Anker-Einführung — und **keiner** der handgezählten Reports trägt ein `befund:`-Feld. Die beiden Zahlen haben **keine Schnittmenge** (`BEO-020`) |
| F-3 | **HIGH** | *„Frage 5 hat heute keine Instanz"* ist falsch: `docs/plan/adr/0022-…:77` trägt den Marker, ist **nicht** in `matrix.exempt-paths` und wird ohne ihn zu einem `matrix-forbidden`-Befund (Produkt-Gegenprobe). Der Absatz führte den Leser an der einzigen Stelle vorbei, an der die Frage zu stellen ist |
| F-4 | **HIGH** | Die *„Arbeitsfläche"* führt **7 von 16** Ankern und **2 von 7 HIGH** — fünf merge-blockierende Klassen fehlen, darunter der Stilles-Grün-Pfad. In derselben Scan-Menge, in der Frage 3 **einmal** vorkam, trifft *„still"* **44**-mal. Für die neun herabgestuften Klassen wurde **keine** Zahl erhoben |
| F-5 | MEDIUM | *„Die Last ist gesenkt"* ist widerlegt: 172 → 203 Zeilen (+18,1 % Byte). Die einzige Lesart, unter der sie sänke, verbietet derselbe Absatz zwei Zeilen darüber |
| F-6 | MEDIUM | `modul-10` sagt wörtlich *„Ein Archiv-Scan ist nicht nötig — die Häufung steht im Register"*. Der Slice fuhr den Archiv-Scan und las das Register nie; es trägt die Zähler für vier der Anker |
| F-7 | MEDIUM | Der **Rückführungs-Trigger** des Slice (§6) ist zur Hälfte eingetreten und wird nirgends erwähnt |
| F-8 | MEDIUM | DoD-Haken 1 verlangt *je Anker* eine Trefferzahl **mit Zeitraum**; geliefert für 3 von 16, der Zeitraum für keinen |
| F-9 | MEDIUM | Die Grundgesamtheit wechselt stillschweigend von **16** (Slice §1) auf **7** (Botschaft) |
| F-10 | MEDIUM | Der benannte Preis *„zwei Ebenen können driften"* ist bei Geburt eingetreten: Tabellenzeile 7 lässt die **Vorbedingung** des Ankers weg und meldet dort, wo das Ziel gar keine Kennung führt |
| F-11 | LOW | Gemischte Polarität: „ja" heißt viermal bestanden, zweimal Befund, einmal gar keine Ja/Nein-Frage |
| F-12 | LOW | *„Gemessen, bevor gekürzt wurde"* beschreibt einen Vorgang, der nicht stattfand — Chronik-Form an einem Zustands-Ort |
| F-13 | LOW | *„hat **heute** keine Instanz"* ist ein gleitender Zeitbezug ohne Pfleger und ohne Auflösungs-Trigger |
| F-14 | LOW | Keine der drei Zahlen ist aus dem Skill reproduzierbar; **18** naheliegende Definitionen von „Befund-Zeile" ergeben keine 687 |

**Negativbefunde.** Die Handzählung stimmt für ihre acht — alle acht
nachgeprüft; richtig wäre sogar **15 von 15**, der genannte Nenner ist eine
Format-Auswahl. Alle vier heute angelegten ADRs nennen die Kennung im
`Schärft:`-Feld. Die Baseline-Ziel-Form ist **unverletzt** (alle sechs
Pflichtteile vorhanden, zweite Ebene unangetastet, Version gesetzt) — eine
vorgeschaltete Zusammenfassung ist eine **Ergänzung**, keine Ersetzung: **kein
`MR-`Eintrag nötig**. `make gates` und `make fullbuild` selbst gefahren, beide
Exit 0.

## Erledigung

Alle vierzehn Befunde sind eingearbeitet; die vier HIGH sind **eigens
nachgemessen** statt übernommen.

- **F-1**, **F-3** Beide Sätze **gestrichen** — sie waren nicht abzuschwächen,
  sondern falsch. Die Gegenprobe zu F-3 fand die `.d-check.yml` selbst: sie
  sagt zwei Zeilen unter der `exempt-paths`-Liste *„neue ADRs ab 0022 tragen
  den Marker"*.
- **F-2**, **F-12**, **F-13**, **F-14** Der ganze Mess-Absatz ist **aus dem
  Skill heraus**. Ein Lauf-Beleg gehört in die Closure-Notiz; im Dauerdokument
  verfällt er still.
- **F-4**, **F-9** Die erste Ebene trägt jetzt **alle sechzehn** HIGH- und
  MEDIUM-Klassen; LOW und INFO bleiben unten, und das steht dort.
- **F-5** Durch das gemessene Delta ersetzt (172 → 206) und der Nutzen als das
  benannt, was er ist: eine Einstiegs-Ordnung, keine Senkung.
- **F-6** Die Register-Zähler stehen in der Botschaft und in §9: `BEO-009`
  sechs, `BEO-012` fünf, `BEO-011` vier, `BEO-004` drei.
- **F-10** Zeile 16 trägt die Vorbedingung.
- **F-11** Jede der sechzehn Fragen ist so gestellt, dass „ja" ein Finding ist.
- **F-7**, **F-8** In §4 und §6 des Slice abgearbeitet.
