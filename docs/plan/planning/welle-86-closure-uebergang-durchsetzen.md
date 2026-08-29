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

- [slice-171](done/slice-171-vorpruefungen-belegen.md) ist geschlossen —
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
| slice-175 | Bindung an den **Übergang**: der `mv`-Commit nach `done/` wird geprüft, nicht der Zustand danach | `.githooks/pre-commit` ruft d-check im `STAGED=`-Modus |

Die drei letzten sind **noch nicht angelegt** — sie entstehen, wenn sie
drankommen; wer alle Slices vor der ersten Implementierung plant, plant tote
Slices (`modul-05`).

### Der Träger von slice-175 ist der git-Hook, nicht der Werkzeug-Hook

**Festgehalten bei der Eröffnung, damit der Slice nicht am falschen Ort
ansetzt.** Das Repo führt zwei Hook-Familien mit sehr verschiedener Reichweite:
`.claude/hooks/` ist **werkzeug-lokal** — [`MR-042`](../../../harness/conventions.md#mr-042)
buchstabiert das aus (*keine CI ruft ihn, ein Lauf ohne dieses Werkzeug ist
ungebunden*). `.githooks/` läuft für **jedes** Werkzeug und jeden Menschen, der
`make hooks` ausgeführt hat.

Eine Durchsetzung im Werkzeug-Hook wäre genauso lokal wie die Zustellung, die
sie flankiert. slice-175 nimmt deshalb `.githooks/pre-commit`.

**Das Muster liegt schon vor und wird nicht neu erfunden.** Der Hook fährt heute
`make adr-check STAGED=1` und `make doc-check` — er **implementiert nichts**, er
ruft d-check. Und der `STAGED=`-Modus beweist, dass das Produkt den gestagten
Zustand lesen kann. Für den Closure-Übergang folgt daraus nur eine Bedingung
davor: erkennen, dass *dieser* Commit einen Slice schließt (gestagter Rename
nach `done/`), und dann die Prüfung fahren, die sonst erst in `make fullbuild`
läuft.

**Drei Orte, eine Logik.** Die Regel lebt im d-check-Modul; Hook und CI rufen
sie, keiner implementiert sie nach — sonst driften drei Fassungen
auseinander. Und weil `--no-verify` den Hook umgeht, braucht dieselbe Prüfung
die **CI-Hälfte**; das Repo macht das bei `adr-check` bereits so
(`pre-commit` **und** PR-CI). Ohne sie ist der Übergang nur höflich gesichert.

**Die vierte Schicht liegt außerhalb dieser Welle:** die pfad-gebundene
Zustellung ([slice-176](open/slice-176-planning-rule-pilot.md))
verhindert keinen Verstoß, sondern die **Überraschung** — ein Hook, der
blockiert, ohne dass der Autor wusste warum, kostet einen ganzen Zyklus. Sie
gehört zum selben Bild, aber nicht zum Closure-Trigger dieser Welle.


### Nachgetragen: was der Hooks-Guide dazu sagt

**Der Abschnitt oben entstand, ohne den Hooks-Guide gelesen zu haben.** Er ist
inzwischen gelesen; die Entscheidung für den git-Hook bleibt, zwei Punkte kommen
dazu.

**Die Entscheidung steht.** Der Guide bestätigt, was
[`MR-042`](../../../harness/conventions.md#mr-042) über die Reichweite sagt: ein
Hook unter `.claude/` läuft im Werkzeug, nicht im Repo. Für eine
**Repo-Invariante** bleibt `.githooks/` der Träger, und die CI-Hälfte bleibt
nötig, weil `--no-verify` den git-Hook umgeht.

**Was fehlte: der `Stop`-Hook fängt früher.** Er löst aus, wenn Claude eine
Runde beendet, und kann sie mit `decision: "block"` **verweigern** — der
mitgegebene Grund geht an das Modell zurück, das dann weiterarbeitet. Das ist
ein anderer Bindepunkt als der Commit: er greift **vor der Übergabe**, nicht
beim Schreiben. Das Repo fährt das Muster bereits
([`stop-require-gates.sh`](../../../.claude/hooks/stop-require-gates.sh)) — nur
für Gates, nicht für den Closure-Übergang. **Er ersetzt den git-Hook nicht**
(dieselbe Werkzeug-Grenze), aber er ergänzt ihn um die Stelle, an der ein
Verstoß noch billig ist. slice-175 entscheidet, ob er beide Hälften nimmt.

**Und: die Begründung in §6 war zu stark.** Dort steht, die **Qualität** eines
Reviews sei out of scope, weil *„Mensch urteilt, Maschine prüft Deckung"*. Der
Guide kennt **prompt-** und **agent-basierte** Hooks ausdrücklich für
Entscheidungen, „die Urteilsvermögen erfordern". Der Ausschluss bleibt richtig —
diese Welle prüft Deckung —, aber sein Grund ist jetzt eine **Wahl** und keine
Unmöglichkeit mehr. Wer ihn später aufhebt, hat einen Weg.

**Nachtrag zum Nachtrag: der Rules-Kanal fällt weg.** Der Abschnitt oben nannte
als vierte Schicht `.claude/rules/` mit `paths`-Frontmatter. Nach
Werkzeug-Auskunft (2026-08-29) hängt eine solche Regel im **Auto-Modus** nur
ein, wenn die Datei über die **dedizierten** Werkzeuge angefasst wird — jeder
Shell-Zugriff geht daran vorbei. In diesem Repo, dessen Auto-Modus zur Shell
rät, hätte der Kanal im Anlassfall geschwiegen. **Auch die Zustellung läuft
deshalb über Hooks**, die das Harness unabhängig vom Werkzeug ausführt;
[slice-176](open/slice-176-planning-rule-pilot.md) ist entsprechend neu
geschnitten, die Klasse steht als [`BEO-024`](observations.md) im Register.

**Ein Rand, den slice-175 kennen muss:** hängen mehrere Hooks am selben Event,
gewinnt die restriktivste Antwort (`deny` vor `defer` vor `ask` vor `allow`).
Ein zweiter Hook kann eine bestehende Erlaubnis also verschärfen, nicht
aufweichen.

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
