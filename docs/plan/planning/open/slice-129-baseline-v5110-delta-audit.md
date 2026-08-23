# Slice slice-129: Etappe B — Delta-Audit über acht Kurs-Wellen, je Welle eine Antwort

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-83-baseline-v5110-migration](../welle-83-baseline-v5110-migration.md)
(zugeordnet bei der Eröffnung).

**Bezug:** Baseline-Regelwerk
[`modul-02-harness-bootstrap.md` §Freshness-Audit](../../../../.harness/baseline/v5.9.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)
(darunter die **Bestands-Stichprobe, die auch bei aktuellem Pin läuft**);
[`AGENTS.md`](../../../../AGENTS.md) §1 (breiterer Pflicht-Blick beim
Drift-Audit); [slice-128](../open/slice-128-baseline-v5110-vendoring.md) (liefert den
Baum, gegen den geprüft wird).

**Berührte Spec-Stellen:** — (das Audit **liest**; was es findet, schneidet
eigene Slices).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Acht Kurs-Wellen (87–94) liegen zwischen unserem alten und dem neuen Pin. Dieses
Audit beantwortet sie **einzeln**: konform ohne Handlung (mit Beleg), Handlung
nötig (mit Slice), oder nicht anwendbar (mit Begründung). **Keine Welle ohne
Zeile** — eine übersprungene ist eine stille Annahme, und stille Annahmen sind
die Klasse, an der die Vorgänger-Welle achtmal gescheitert ist.

Drei Wellen sind vorab als wahrscheinlich folgenreich markiert — **als
Verdacht, nicht als Befund**:

- **Kurs-Welle 90** („ab `Accepted` zählt jede Zeile") berührt vermutlich
  `make adr-check` und das Modul `vcs` — also unsere ADR-Immutabilität.
- **Kurs-Welle 93** („AGENTS.md §4 wird die Autorität über die Targets")
  berührt vermutlich das Modul `targets`, `make gate-consistency` und unsere
  eigene §4-Tabelle.
- **Kurs-Welle 94** („eine Rangliste ordnet, jetzt deckt sie auch ab") bringt
  die **Vollständigkeits-Zusage**. Eine Verletzung ist bereits bekannt und hat
  ihren Slice ([slice-127](../next/slice-127-claude-md-pointer.md)); die Frage
  hier ist, ob es **weitere** gibt — Skill-Dateien, emittierte Fragmente,
  `.claude/commands/`.

Die übrigen fünf sind damit **nicht** als folgenlos erklärt. Sie bekommen
dieselbe Zeile wie die drei.

## 2. Vorgehen

1. **Je Kurs-Welle 87–94 eine Zeile** mit Antwort und Beleg. Quelle ist das
   Kurs-CHANGELOG **und** der Regelwerks-Diff, nicht nur die Überschrift — eine
   Wellen-Überschrift ist eine Zusammenfassung, kein Vertrag.
2. **Die Bestands-Stichprobe fahren**, die das Freshness-Audit auch bei
   aktuellem Pin verlangt — sie prüft, ob der gelebte Bestand dem Regelwerk
   entspricht, nicht nur ob der Pin stimmt.
3. **Die Vollständigkeits-Zusage auf das eigene Repo anwenden** (lesend):
   welche Artefakte außerhalb der Rangliste tragen normativen Text? Die
   Prüffrage lautet *„Steht jede Aussage dieser Datei auch in einer gerankten
   Quelle?"* — Ergebnis ist eine **Liste**, keine Bereinigung.
4. **Etappe C schneiden:** aus den Handlungs-Zeilen werden Slices, mit
   Drift-Log-Eintrag in der Roadmap. Was zu groß ist, wird als eigene Welle
   ausgewiesen statt angehängt.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Umsetzung.** Das Audit liest und schneidet; es ändert weder Code noch
  Konfiguration noch Doku außerhalb seines eigenen Ergebnisses.
- **Kein Vorgriff auf [slice-127](../next/slice-127-claude-md-pointer.md).**
- **Keine Sammel-Antwort.** „Die übrigen fünf sind folgenlos" ist genau die
  Aussage, die dieses Audit verbietet.

## 4. Definition of Done

- [ ] Acht Zeilen, eine je Kurs-Welle 87–94, jede mit Antwort **und** Beleg;
      keine Sammel-Zeile.
- [ ] Bestands-Stichprobe gefahren und ihr Ergebnis benannt.
- [ ] Liste der Artefakte außerhalb der Rangliste, die normativen Text tragen —
      vollständig **belegt**, nicht behauptet (`BEO-011`).
- [ ] Etappe C geschnitten oder als Folge-Welle ausgewiesen; Drift-Log-Zeile
      gesetzt.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Audit über acht Wellen verführt zur Sammel-Antwort.** Je länger die
  Liste, desto attraktiver „betrifft uns nicht". Genau dagegen steht die
  Ein-Zeile-je-Welle-Regel. — **Ausgang:** *(bei Closure)*
- **Die drei Verdachts-Wellen könnten das Audit dominieren** und die fünf
  übrigen zur Formsache machen. — **Ausgang:** *(bei Closure)*
- **Die Vollständigkeits-Liste ist selbst eine Vollständigkeits-Aussage** — und
  damit die Form, die in welle-82 achtmal gekippt ist. Sie braucht denselben
  zeilenweisen Beleg, den sie einfordert. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-128](../open/slice-128-baseline-v5110-vendoring.md)
in `done/` — das Audit liest den **neuen** Baum.

**Rückführungen:** `in-progress` → `next`, falls das Audit eine Produkt-
Konsequenz findet, die vor dem Rest des Audits gehört (etwa ein rotes Gate).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Konventionen (GF), Gate-Mechanik (GF),
  Nutzer-Doku (GF, nur lesend).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23): **`BEO-011`**
  ist die zentrale — dieses Audit produziert zwei Vollständigkeits-Aussagen
  (acht Wellen abgedeckt; alle normativen Artefakte gelistet), und genau solche
  Aussagen sind die Klasse. **`BEO-002`** für die Ränder jeder Regel, die das
  Audit als „konform" abhakt. **`BEO-009`** für jede Zahl im Ergebnis.

Slice-ID: slice-129. Betroffene IDs:
[`MR-021`](../../../../harness/conventions.md#mr-021). Module:
Harness-Konventionen, Gate-Mechanik. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Audit nach der Form der v5.6.0-Migration,
die dieselbe Aufgabe über sechs Stufen gelöst hat.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
