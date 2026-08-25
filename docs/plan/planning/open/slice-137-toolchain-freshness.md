# Slice slice-137: Zwei gepinnte Toolchains, die niemand beobachtet

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos.** Seine Closure-Bedingung wäre seine eigene DoD, und
Baseline-Regelwerk
[`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht](../../../../.harness/baseline/v5.11.0/regelwerk/modul-06-roadmap.md)
nennt das Zeremonie: *„Ein Trigger, der nichts beobachtet, was die Slices nicht
ohnehin belegen, ist Zeremonie."*

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.1 (Docker/make-only), §4
(Gate-Tabelle); [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(das Produkt, nicht der innere Lauf); der Nachtlauf
[`upstream-drift.yml`](../../../../.github/workflows/upstream-drift.yml) und das
Muster von `make baseline-freshness` aus
[slice-133](../done/slice-133-baseline-sensor-verdrahten.md).

**Berührte Spec-Stellen:** — (Harness-Sensoren; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

`GO_VERSION=1.27.0` und `GOLANGCI_LINT_VERSION=v2.13.1` stehen im `Makefile` und
im `Dockerfile`, beide zusätzlich digest-gepinnt. **Nichts prüft, ob upstream
weiter ist.** Gemessen: sieben Fundstellen der beiden Namen in `Makefile` und
Workflows — alle **nutzend** (Build-Arg, `make versions`-Ausgabe), **keine
prüfend**.

Die Folge ist belegt, nicht vermutet: dass der Kurs `v5.11.0` veröffentlicht
hatte, hat uns der Auftraggeber gesagt, nicht der Sensor. Für die Baseline-Achse
ist das seit [slice-133](../done/slice-133-baseline-sensor-verdrahten.md)
behoben; die zwei Toolchain-Achsen sind die verbliebene Hälfte derselben
Blindheit.

**Der teure Teil ist bereits gebaut.** Der Nachtlauf `upstream-drift.yml`
existiert, getrennt von `ci.yml`, mit der Begründung, warum ein Upstream-Ausfall
die Integrations-CI nie rot färben darf. Dieser Slice hängt zwei Targets hinein.

## 2. Vorgehen

1. **Einen parametrierten Sensor** schreiben statt zwei Kopien: Name, gepinnter
   Wert, Quelle. Zwei Quellen-Formen, weil `go` keine GitHub-Releases führt —
   `releases/latest`-Redirect für golangci-lint, `go.dev/VERSION?m=text` für Go.
2. **Fetch, Normalisierung und Vergleich trennen.** Der Vergleicher muss
   **netzlos** aufrufbar sein, sonst ist er nicht prüfbar; die Normalisierung
   ebenso (`go1.27.0` ⇄ `1.27.0`).
3. **Fail-open wie die Baseline-Achse:** Netz- oder Werkzeug-Ausfall ⇒ `SKIP`,
   nicht rot. Mit Zeitgrenzen — ohne sie wäre eine hängende Verbindung genau der
   Fehler, den [slice-133](../done/slice-133-baseline-sensor-verdrahten.md)
   schon einmal einarbeiten musste.
4. **In den bestehenden Nachtlauf hängen**, nicht in einen zweiten. Ein Ausfall
   einer Achse darf die anderen nicht verdecken.
5. **`AGENTS.md` §4 und die Sensors-Tabelle** nachziehen — sonst
   `gate-consistency`-rot.
6. **Bewusstes Brechen** je Achse: gepinnter Wert künstlich veraltet ⇒ rot mit
   **gelesener Ursache**; Rückbau ⇒ grün. Dazu die netzlose Probe des
   Vergleichers.
7. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Auto-Bump.** Der Sensor meldet; die Hebung bleibt ein bewusster Akt —
  und sie zieht beim Go-Bump das golangci-Pendant nach (Auftraggeber-Regel).
- **Nichts in `gates`.** Netzhaltiges gehört nicht in den inneren Lauf.
- **Keine dritte Achse.** Docker-Basis-Images und Action-Pins bleiben draußen;
  ob sie eine Achse brauchen, ist eigens zu entscheiden.
- **Kein d-check-Modul.** Die Pins stehen im `Dockerfile`/`Makefile`, also
  außerhalb der gescannten Markdown-Menge (§3.8).

## 4. Definition of Done

- [ ] Zwei Targets, ein Skript; der **Vergleicher ist netzlos aufrufbar** und
      mit Fixture-Werten geprüft.
- [ ] Beide hängen im **bestehenden** Nachtlauf, keine Achse verdeckt die andere.
- [ ] Je Achse ein **konstruierter Verstoß** mit gelesener Ursache: veralteter
      Pin ⇒ rot, Rückbau ⇒ grün. Plus der Ausfall-Pfad ⇒ `SKIP`, nicht rot.
- [ ] `AGENTS.md` §4 und die Sensors-Tabelle tragen beide Targets;
      `gate-consistency` grün.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein fail-open-Sensor, der immer `SKIP` meldet, ist ein grüner Blindgänger.**
  Die Zusage lautet *„meldet, wenn upstream weiter ist"* — sie ist nur wahr,
  wenn der Fetch-Pfad wirklich läuft. Der Nachweis ist ein **echter** Lauf, kein
  Fixture. — **Ausgang:** *(bei Closure)*
- **Die Vergleichsform ist nicht Semver.** `baseline-freshness` fragt
  Gleich/Ungleich gegen den neuesten Tag; eine monotone Reihe kennt kein „neuer,
  aber älter". Ob das für `go.dev` ebenso trägt, ist zu prüfen und nicht zu
  übernehmen. — **Ausgang:** *(bei Closure)*
- **Zwei Achsen in einem Nachtlauf-Job:** bricht die erste ab, läuft die zweite
  womöglich nicht. Der Slice muss zeigen, dass jede Achse ihr Urteil abgibt. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — alle Lifecycle-Verzeichnisse
sind leer, welle-84 ist geschlossen.

**Rückführungen:** `in-progress` → `next`, falls sich zeigt, dass eine der
Quellen ein Semver-Urteil verlangt statt Gleich/Ungleich — dann ist die
Vergleichsform eine eigene Entscheidung.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Werkzeuge (GF), CI (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-007`](../observations.md) für jeden Beleg-Lauf — Exit direkt lesen.
  [`BEO-011`](../observations.md) für die Aussage, was der Sensor abdeckt: *zwei*
  Achsen, nicht *„die Toolchain"*. [`BEO-010`](../observations.md), weil zwei
  neue Targets in drei Doku-Flächen erscheinen.

Slice-ID: slice-137. Betroffene IDs: — (Harness-Sensoren; keine Anforderung).
Module: Harness-Werkzeuge, CI. Gates: `make gate-consistency`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — zweite Hälfte einer bereits gebauten Achse.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
