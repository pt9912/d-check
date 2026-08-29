# Welle welle-86: Der Closure-Übergang trägt seine Vorbedingungen

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach in `docs/plan/planning/`; bei der Closure wandert sie nach `done/` und
bekommt ihre Ergebnisnotiz `welle-86-results.md` daneben.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-08-29.

---

## 1. Welle-Ziel

**Der Übergang `in-progress → done` findet nicht statt, wenn eine seiner
Vorbedingungen fehlt — statt hinterher zu melden, dass er stattgefunden hat.**

Das Regelwerk nennt die Bedingung unmissverständlich:

<!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:33-34 -->
> DoD-Häkchen und Closure-Notiz
> sind die Bedingung dafür, dass die Datei überhaupt nach `done/` darf.

**Gehalten wird davon heute die Hälfte.** `verify-closure-notes` prüft die
Closure-Notiz und die Risiko-Ausgänge. Vier weitere Vorbedingungen desselben
Übergangs hat kein Sensor:

| Vorbedingung | Quelle | heute |
|---|---|---|
| Closure-Notiz vorhanden, Substanz, keine Floskel | `modul-05` | gehalten (`planning`) |
| jedes Risiko trägt einen Ausgang | `modul-05` | gehalten (`structure`) |
| DoD-Häkchen gesetzt | `modul-05` | **offen** |
| Beobachtungs-Register fortgeschrieben | `modul-06` | **offen** |
| Review fand statt | `modul-01`, `modul-08` | **offen** |
| die Prüfung hängt am **Übergang**, nicht am Zustand | `modul-05` | **offen** |

**Der Anlass ist gemessen und liegt in diesem Repo:** slice-168, -169 und -170
gingen mit offenem Review-Haken nach `done/`, ohne Review, ohne Verifikation und
ohne Register-Fortschreibung. Kein Gate hat es gemeldet, weil keines danach
sieht. Über den Bestand: **37** von 169 `done/`-Slices tragen mindestens einen
offenen Haken; **87 von 95** Slices mit Review-Zusage haben einen Report in
`docs/reviews/` — die Konvention wird gelebt, ihr fehlt nur der Wächter.

**Offene Beobachtungen gesichtet** (Eröffnungs-Schritt 2, `modul-06`;
Register-Stand 2026-08-29, höchste Kennung `BEO-021`). Drei betreffen die
Sub-Areas dieser Welle:

- [`BEO-013`](observations.md) — *ein Wächter, der nichts mehr fängt, bleibt
  stehen*: hier eine Stufe früher — eine Regel, die **nie** einen Wächter
  hatte. Die Bestands-Ausnahmen, die diese Welle einführt, sind der künftige
  Kandidat für genau diese Klasse.
- [`BEO-011`](observations.md) — *Regel aus dem Anlass statt aus dem Bestand*:
  die **Dringlichkeit** stammt aus drei Slices einer Sitzung, die **Regel** aus
  37 bzw. 87/95 gemessenen Fällen. Der Unterschied gehört in jede Slice-DoD
  dieser Welle.
- [`BEO-015`](observations.md) — *ein offener Punkt bekommt bei der Closure
  einen Ausgang, den es nicht gibt*: dieselbe Familie urteilsfreier
  Closure-Prüfungen, die diese Welle erweitert.

Keiner der drei hat mit dieser Welle 3× erreicht; sie gehen als Risiko in die
betroffenen Slices, nicht als eigener Slice.

## 2. Trigger (Welle startet)

- [slice-171](in-progress/slice-171-vorpruefungen-belegen.md) ist geschlossen —
  er belegt die Lektüre der Regel, auf die diese Welle sich beruft, und hält den
  WIP-Slot.

## 3. Closure-Trigger (Welle schließt)

Das **Mehr** gegenüber den einzelnen Slice-DoDs: jede DoD belegt *einen*
Sensor; die Welle belegt den **Übergang als Ganzes**.

- Alle vier Slices liegen in `done/`.
- `make fullbuild` grün.
- **Ein konstruierter Test-Slice mit je einer fehlenden Vorbedingung wird beim
  `mv`-Commit abgewiesen** — vier Proben, je eine pro Vorbedingung, mit
  Erwartung und Ergebnis. Das ist der Beleg, den keine einzelne Slice-DoD
  liefert: dass die Bedingungen **am Übergang** greifen und nicht nur im
  Nachhinein melden.
- Closure-Notiz in `done/welle-86-results.md`.

## 4. Slices in dieser Welle

| Slice | Gegenstand | Werkzeug |
|---|---|---|
| [slice-172](open/slice-172-closure-uebergang-waechtern.md) | DoD-Häkchen gesetzt | `structure` (vorhanden) |
| slice-173 | Review-Report-Deckung: jeder `done/`-Slice mit Review-Zusage hat einen Report | neue Fähigkeit (Deckung zweier Mengen) |
| slice-174 | Beobachtungs-Register-Deckung: zitierte `BEO-<NNN>` hat Registerzeile, jede Zeile trägt einen Beleg | `modul-06` nennt die maschinelle Hälfte selbst |
| slice-175 | Bindung an den **Übergang**: der `mv`-Commit nach `done/` wird geprüft, nicht der Zustand danach | `vcs`/`commits`-Port, `pre-commit`-Hook |

Die drei letzten sind **noch nicht angelegt** — sie entstehen, wenn sie
drankommen; wer alle Slices vor der ersten Implementierung plant, plant tote
Slices (`modul-05`).

## 5. Abhängigkeiten

- **slice-175 setzt die drei anderen voraus.** Er bindet die Prüfungen an den
  Übergang; ohne sie bände er nichts. Umgekehrt sind 172–174 einzeln lieferbar
  und wirken schon als Zustands-Prüfung.
- **slice-173 und slice-174 sind unabhängig voneinander**, teilen aber die
  Klasse (Deckung zweier Mengen). Wer zuerst liefert, prägt die Form.

## 6. Out-of-Scope für diese Welle

- **Die Qualität eines Reviews.** Ein Report ist ein Artefakt, kein
  Qualitätsnachweis. Die Welle prüft Deckung — *Mensch urteilt, Maschine prüft
  Deckung* (`modul-06`).
- **Die Rollen-Trennung selbst.** Ob Review und Verifikation in *eigenen
  Kontexten* liefen, ist im Repo nicht sichtbar. Genau der Verstoß der
  auslösenden Sitzung wäre für keinen Sensor erkennbar gewesen — nur seine
  Folgen.
- **Die übrigen Lifecycle-Übergänge** (`open→next`, `next→in-progress`, die
  beiden Rückführungen). Diese Welle nimmt **einen** Übergang; die anderen sind
  ein eigener Schnitt.
- **Das Nachrüsten der 37 Bestands-Slices.** Sie sind Belege ihrer Zeit; ein
  nachträglich gesetzter Haken behauptet einen Review, den es nicht gab.
- **Die `AGENTS.md`-Umschichtung** (Zyklus und Rollen rein, Gate-Referenz raus).
  Der Befund steht — 511 Zeilen, Duplikation gegen `harness/README.md` —, aber
  er ist eine Form-Frage, keine Durchsetzungs-Frage.

## 7. Closure-Notiz
