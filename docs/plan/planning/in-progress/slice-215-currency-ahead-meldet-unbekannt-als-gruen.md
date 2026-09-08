# Slice slice-215: Der Currency-Zweig meldet „unbekannt" als grün

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`BEO-HARN/check-latest-blind-before-pin`](../observations/BEO-HARN/check-latest-blind-before-pin/observation.md)
(1×, Stand *offen* — der Anlass; von
[slice-214](../done/slice-214-gepinnter-cache-ausserhalb-des-repos.md) beinahe
übersehen), [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette,
[`MR-004`](../../../../harness/conventions.md#mr-004) (Sub-Area `tools/harness/`).

**Berührte Spec-Stellen:** — *(keine; `baseline-freshness` ist kein Produkt-Gate
und trägt keine `DC-*`-Bindung)*

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Der `ahead`-Zweig von `fetch-baseline-cache.sh --check-latest` meldet
*„Pin nicht in der Release-Liste"* **auf stderr und beendet mit 0**. Der
Nachtlauf liest den Exit-Code — also **grün**, während die Currency-Frage
unbeantwortet blieb. Das ist ein **stiller Grün-Pfad** in einem lebenden
Sensor; er wird geschlossen oder begründet beibehalten, und die Bedingung, die
ihn erreichbar macht, wird benannt.

**Der Weg dorthin ist gemessen und datiert.** Die Currency-Prüfung liest
`releases?per_page=100` **ohne Paginierung**. GitHub liefert neueste zuerst;
liegt der Pin außerhalb der jüngsten 100 Releases, findet ihn die Liste nicht
und der Zweig ist `ahead`. **Heute: 58 Releases von 100.** Der Pfad ist also
nicht hypothetisch, sondern terminiert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Ursache der Anlass-Beobachtung nachträglich beweisen.** Sie ist beim
  `v6.0.0`-Bump aufgetreten, nie untersucht und **heute nicht reproduzierbar**:
  Release-Liste und Tag-Liste des Kurs-Repos sind deckungsgleich (58/58), die
  naheliegende Erklärung *„Tag ohne Release-Objekt"* trägt also nicht. Der
  Eintrag bekommt einen **untersuchten** Stand, keine erfundene Ursache.
- **`baseline-freshness` in `gates` ziehen.** Es ist bewusst fail-open und
  netz-gebunden ([`MR-011`](../../../../harness/conventions.md#mr-011)-Kette);
  das zu ändern wäre ein Entscheid mit ADR und eine andere Frage als die hier
  gestellte.
- **Die übrigen Freshness-Achsen.** `freshness-go`, die drei Action-Pins, die
  Digest-Achsen — sie teilen die fail-open-Bauart, aber nicht den `ahead`-Fall;
  ob sie einen eigenen stillen Pfad haben, ist ein anderer Vorgang.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Der `ahead`-Zweig ist **entschieden**: Exit-Code gesetzt (und
      damit im Nachtlauf sichtbar) **oder** begründet bei 0 belassen — die
      Begründung nennt, wer den stderr-Text dann liest.
- [ ] **(2)** Ein **Bruch-Test** belegt die Entscheidung: ein Pin, der nicht in
      der Liste steht, führt zu dem Verhalten, das (1) festlegt — mit echter
      Ausgabe, netzlos oder mit benanntem Netz-Bedarf.
- [ ] **(3)** Das **100er-Fenster** steht in
      [`harness/sensors/baseline-freshness.md`](../../../../harness/sensors/baseline-freshness.md)
      §Grenze, mit der gemessenen Zahl und der Bedingung, unter der es greift.
      Der Registereintrag trägt einen untersuchten Stand.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Zuerst die Form des Gegenstands** ([`AGENTS.md`](../../../../AGENTS.md) §5),
denn (1) ist ein Urteil über einen Ausgang und nicht über eine Zahl.

**Was zählt als *stiller Grün-Pfad*?** Ein Lauf, der (a) einen Zustand erreicht,
in dem die zugesagte Frage **unbeantwortet** bleibt, und (b) mit einem
Exit-Code endet, den der **Konsument** als „geprüft und in Ordnung" liest. Der
Konsument ist hier der Nachtlauf-Job, und er liest ausschließlich den
Exit-Code — `make baseline-freshness` ruft das Skript ohne Auswertung der
Ausgabe. **Ein stderr-Text allein erfüllt (b) nicht als Gegenbeweis**, sondern
ist genau das Merkmal.

**Die vier Currency-Zustände, aus dem Skript gelesen:**

| Zustand | Bedeutung | Exit heute |
|---|---|---|
| `current` | Pin ist der neueste Release-Tag | 0 — richtig |
| `newer` | neuere Tags vorhanden | 3 — sichtbar |
| `skip` | Liste nicht lesbar (Netz/Rate-Limit) | 0 — **fail-open, deklariert** |
| `ahead` | Pin **nicht in der Liste** | 0 — **undeklariert** |

**`skip` und `ahead` sind nicht dasselbe, und darin liegt der Slice.** `skip`
ist die deklarierte fail-open-Wahl: Die Prüfung konnte gar nicht laufen, und
das Repo hat entschieden, deshalb nicht rot zu werden. `ahead` heißt: Die
Prüfung **lief** und ihr Ergebnis ist *„ich finde den eigenen Pin nicht"* —
ein Zustand, der entweder auf ein zurückgezogenes Release oder auf das
Listen-Fenster zeigt. **Beides sind Befunde, keine Ausfälle.**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`tools/harness/fetch-baseline-cache.sh`](../../../../tools/harness/fetch-baseline-cache.sh) | update, falls (1) auf „Exit setzen" fällt | der `ahead`-Zweig |
| [`harness/sensors/baseline-freshness.md`](../../../../harness/sensors/baseline-freshness.md) | update | das Fenster und der Ausgang |
| [`BEO-HARN/check-latest-blind-before-pin`](../observations/BEO-HARN/check-latest-blind-before-pin/state.md) | update | untersuchter Stand |

## 4. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.
[slice-214](../done/slice-214-gepinnter-cache-ausserhalb-des-repos.md) liegt in
`done/` — dort ist der Eintrag aufgefallen, der diesen Slice auslöst.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt sich, dass der Bruch-Test aus DoD (2)
  nur mit einer **Netz-Attrappe** oder einem Umbau der Skript-Struktur zu
  haben ist, ist das eigene Arbeit — dann wird der Test ein eigener Slice und
  dieser schließt mit (1) und (3).
- `in-progress` → `open` (blockiert): Ergibt (1), dass ein Exit-Code den
  Nachtlauf bei **jedem** Rate-Limit rot machen würde, ruht der Slice bis zum
  Entscheid über die fail-open-Linie — die ist
  [`MR-011`](../../../../harness/conventions.md#mr-011)-gebunden und nicht
  nebenbei zu ändern.

## 5. Closure-Trigger

Zwei beobachtbare Kriterien und ein Lerneintrag: (a) der `ahead`-Zweig hat eine
entschiedene und **belegte** Ausgangs-Semantik, mit echter Ausgabe im Slice;
(b) `make gates` ist grün, das Fenster steht in der Sensor-Datei und der
Registereintrag trägt seinen Stand.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Bruch-Test braucht Netz oder eine Attrappe.** Die Currency-Prüfung
  liest die GitHub-API; ein Pin, der nicht in der Liste steht, lässt sich
  netzlos nur mit einer Attrappe erzeugen. `--selftest` deckt heute die
  **Alias-Auflösung**, nicht den Currency-Zweig. **Wenn der Beleg nur mit Netz
  zu haben ist, gehört das gesagt** — ein behaupteter Bruch-Test wäre schlimmer
  als keiner. — **Ausgang:** \<offen\>
- **Ein Exit-Code kann die fail-open-Linie verschieben, ohne dass es auffällt.**
  `skip` (Netz weg) und `ahead` (Pin nicht gefunden) laufen heute beide auf 0.
  Wer nur `ahead` hebt, muss sicher sein, dass kein Netz-Ausfall in diesem
  Zweig landet — sonst wird der Nachtlauf bei jeder API-Störung rot, und die
  bewusste fail-open-Wahl ist still gekippt. — **Ausgang:** \<offen\>
- **Die Anlass-Beobachtung bleibt möglicherweise ungeklärt.** Heute sind
  Release- und Tag-Liste deckungsgleich; die naheliegende Erklärung trägt
  nicht. Ein Slice, der die Ursache **nicht** findet, muss das als Ergebnis
  hinschreiben — und der Registereintrag bleibt dann bei 1× mit untersuchtem
  Stand, statt eine Ursache zu bekommen, die niemand gemessen hat.
  — **Ausgang:** \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende. Dieses Repo führt **drei** Prüfungen — die
zwei kanonischen und, als Adaption, den Nachtlauf-Stand
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:268-269 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Zwei** Sub-Areas, und diesmal wird die zweite **geändert**, nicht nur
gelesen — anders als in slice-214:

- **`tools/harness/`** (Kürzel `HARN`,
  [`MR-004`](../../../../harness/conventions.md#mr-004)) — der **Gegenstand**:
  [`fetch-baseline-cache.sh`](../../../../tools/harness/fetch-baseline-cache.sh)
  und sein Currency-Zweig. Fällt DoD (1) auf „Exit setzen", bekommt die
  Sub-Area einen **Verhaltens**-Diff.
- **`*`** (Repo-Default) — die Sensor-Beschreibung und der Registereintrag.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **38** Verzeichnisse über beide
Kürzel). **Diesmal wurde `BEO-HARN/` einzeln geöffnet und nicht angenommen** —
das ist die Lehre aus slice-214, wo ein „kein Eintrag" ungeprüft dastand.
**Vier** Einträge sind einschlägig:

- [`BEO-HARN/check-latest-blind-before-pin`](../observations/BEO-HARN/check-latest-blind-before-pin/observation.md)
  (1×, Stand *offen*) — **der Anlass**, und der einzige Eintrag unter diesem
  Kürzel. Er bleibt bei 1×: Dieser Slice untersucht ihn, er tritt nicht erneut
  auf.
- [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
  (9×, Ausgang *geplant*) — **die Klasse des Gegenstands.** Der `ahead`-Zweig
  ist genau das: ein Lauf, der wie ein bestandener aussieht und die zugesagte
  Frage nicht beantwortet hat. Der Ableiter — *„wer einen Wächter für eine
  Zusage schreibt, dreht sie einmal um und belegt, dass genau ein Test rot
  wird"* — ist DoD (2) dieses Slice.
- [`gruene-wettlauf-probe-beweist-nichts`](../observations/BEO-ALL/gruene-wettlauf-probe-beweist-nichts/observation.md)
  — für den Bruch-Test: Eine Probe, die grün ist, weil sie den Zustand gar
  nicht erzeugt hat, belegt nichts. §6 führt das als erstes Risiko, weil der
  Currency-Zweig Netz braucht.
- [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  (15×, Stand *gemischt*) — nach slice-214 der wachsamste Eintrag: Dort war die
  Menge um fünf Mitglieder zu klein. Hier ist die Menge klein und benannt
  (**vier** Currency-Zustände, aus dem Skript gelesen), aber die Aussage über
  den **Konsumenten** — *„der Nachtlauf liest nur den Exit-Code"* — gehört am
  Workflow geprüft, nicht angenommen.

**Keiner der vier erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-08 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (jüngster Lauf 2026-09-07T05:33:45Z) und `image-scan.yml`
(2026-09-07T08:21:32Z). **Und dieses Grün ist der Gegenstand des Slice, nicht
seine Kulisse:** `upstream-drift.yml` fährt `make baseline-freshness`. Sein
Grün heißt heute *„Pin ist der neueste Tag"* (Zustand `current`) — es hieße
aber **dasselbe**, wenn der Zustand `ahead` wäre. Genau diese Ununterscheidbarkeit
schließt der Slice.

**Modus-Begründungsblock.** Beide berührten Sub-Areas GF — ein Block je
Sub-Area.

### Sub-Area: `tools/harness/`

- **Modus:** GF ([`MR-004`](../../../../harness/conventions.md#mr-004)).
- **Konventions-Dichte:** hoch — die Skript-Mechanik ist über
  [`MR-004`](../../../../harness/conventions.md#mr-004),
  [`MR-005`](../../../../harness/conventions.md#mr-005) und die
  [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette getragen; die
  fail-open-Linie ist deklariert und darf nicht nebenbei kippen (§6).
- **Phase-Reife:** Phase 5 — das Skript ist seit vielen Bumps in Gebrauch.
- **Evidenz-/Diskrepanz-Risiko:** **hoch für den Beleg, niedrig für den
  Befund.** Dass der Zweig 0 zurückgibt, ist aus dem Skript gelesen und
  eindeutig. Ob sich das netzlos **brechen** lässt, ist offen — `--selftest`
  deckt die Alias-Auflösung, nicht den Currency-Zweig.
- **Reconciliation-Aufwand:** keiner (GF).

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** hoch für die Form der Sensor-Datei; seit slice-213
  sagt [`AGENTS.md`](../../../../AGENTS.md) §5, wie eine Grenze zu prüfen ist.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** niedrig — die Zahl (58 von 100) ist gemessen
  und die Bedingung aus dem Skript gelesen.
- **Reconciliation-Aufwand:** keiner (GF).
