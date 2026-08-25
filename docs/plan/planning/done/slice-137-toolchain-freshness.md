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

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

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

- [x] Zwei Targets, ein Skript; `--compare` ruft **nur** den Vergleicher, mit
      drei Fixture-Proben geprüft (gleich ⇒ 0, ungleich ⇒ 3, leer ⇒ `SKIP`/0).
- [x] Beide hängen im **bestehenden** Nachtlauf, drei Schritte, die späteren mit
      `if: always()` — keine Achse verdeckt die andere.
- [x] Je Achse ein konstruierter Verstoß mit **gelesener Ursache**: `1.20.0` bzw.
      `v2.0.0` ⇒ *VERALTET* mit genanntem Upstream-Stand; Rückbau ⇒ `ok`. Der
      Ausfall-Pfad ⇒ `SKIP`, **zweimal** belegt: unerreichbares Repo und — nach
      dem Review — eine Müll-Antwort mit HTTP 200.
- [x] `AGENTS.md` §4 und die Sensors-Tabelle tragen beide Targets;
      `gate-consistency` grün.
- [x] `make gates` Exit 0 (zehn Glieder); unabhängiger Review
      ([Report](../../../reviews/2026-08-25-slice-137-toolchain-freshness-review.md)),
      blockierend mit einem HIGH, alle vier Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Ein fail-open-Sensor, der immer `SKIP` meldet, ist ein grüner Blindgänger.**
  — **Ausgang:** *nicht eingetreten, und der Nachweis war ein echter Lauf.*
  Beide Achsen holen wirklich und melden `ok`. **Das Gegenstück ist dafür
  eingetreten:** der Sensor meldete zu viel statt zu wenig — eine
  HTTP-200-Müll-Antwort galt als Versionsstand und ergab ein falsches
  *VERALTET*. Beide Fehlrichtungen sind jetzt belegt.
- **Die Vergleichsform ist nicht Semver.** — **Ausgang:** *geprüft, nicht
  übernommen.* Gleich/Ungleich trägt für beide Achsen, weil beide Reihen monoton
  sind; der Skript-Kopf sagt das als **Entscheidung** hin, samt der Bedingung,
  unter der sie fiele.
- **Zwei Achsen in einem Nachtlauf-Job.** — **Ausgang:** *erkannt und gelöst,
  bevor es eintrat.* Drei Schritte statt einem, die späteren mit `if: always()`.
  Ein Job, der nach der ersten roten Achse abbricht, verdeckt die zweite.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — alle Lifecycle-Verzeichnisse
sind leer, welle-84 ist geschlossen.

**Rückführungen:** `in-progress` → `next`, falls sich zeigt, dass eine der
Quellen ein Semver-Urteil verlangt statt Gleich/Ungleich — dann ist die
Vergleichsform eine eigene Entscheidung.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Werkzeuge (GF), CI (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25):
  [`BEO-007`](../observations.md) für jeden Beleg-Lauf — Exit direkt lesen.
  [`BEO-011`](../observations.md) für die Aussage, was der Sensor abdeckt: *zwei*
  Achsen, nicht *„die Toolchain"*. [`BEO-010`](../observations.md), weil zwei
  neue Targets in drei Doku-Flächen erscheinen.

Slice-ID: slice-137. Betroffene IDs: — (Harness-Sensoren; keine Anforderung).
Module: Harness-Werkzeuge, CI. Gates: `make gate-consistency`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — zweite Hälfte einer bereits gebauten Achse.

## 9. Closure-Notiz (nach `done/`)

Geliefert: zwei Targets über **einen** parametrierten Sensor, beide im
bestehenden Nachtlauf. Die Blindheit, die
[slice-133](slice-133-baseline-sensor-verdrahten.md) zur Hälfte behoben hat, ist
damit ganz behoben — drei gepinnte Fremd-Bestände, drei Achsen, ein Job.

**Der Fund des Slice ist eine Asymmetrie in meiner eigenen Sorgfalt.** Den
GitHub-Zweig hatte ich strukturell abgesichert: die Endstation des Redirects
*muss* `/releases/tag/` enthalten, sonst gilt der Stand als nicht ermittelbar.
Den `go.dev`-Zweig nicht — dort nahm ich **jede** HTTP-200-Antwort als
Versionsstring. Eine Fehlerseite hätte ein *VERALTET* erzeugt, also das genaue
Gegenteil der Zusage, die im Skript-Kopf, in `AGENTS.md` §4 und in der
Sensors-Zeile gleichlautend steht: **Parse-Ausfall ist SKIP.** Zwei Zweige,
dieselbe Frage, nur einer beantwortet.

**Und die Fehlrichtung ist die seltenere.** Ein fail-open-Sensor scheitert
üblicherweise ins Stille — er meldet nichts und sieht grün aus. Dieser hier
hätte ins Laute gescheitert: ein Alarm ohne Anlass. Beide Richtungen sind jetzt
mit einem konstruierten Fall belegt, und der zweite brauchte einen `curl`-Stub
im `PATH`, weil er anders gar nicht herstellbar ist.

**Eine Aussage ist vom Zitat zur Messung geworden.** *„golang/go publiziert
keine Release-Objekte"* stand dreifach im Repo, aus dem Schwester-Repo
übernommen und nie geprüft. Nachgemessen: `golang/go/releases/latest` landet auf
`.../releases` — **ohne** `/tag/` —, `golangci-lint` dagegen auf
`.../releases/tag/v2.13.1`. Der GitHub-Zweig SKIPpte bei Go also korrekt, statt
falsch zu vergleichen. Die Sonderquelle ist damit begründet und nicht behauptet.

**Selbst gefunden, vor dem Commit:** meine drei neuen Kommentar-Stellen
verletzten `AGENTS.md` §3.7 — eine Slice-Kennung im Skript-Kopf, eine
Befund-Nummer in einem Kommentar, ein Mess-Label mit Datum in der neuen
Quellen-Aussage. Herkunft gehört als **ein** auflösbares Feld in den Kommentar
oder gar nicht hinein; die Begründung selbst steht weiter da, nur ohne ihre
Provenienz.

**Wellenlos, und das war die richtige Form.** Der Slice hat keine
Closure-Bedingung, die mehr beobachtet als seine eigene DoD. Ihn in eine Welle zu
packen hätte einen Trigger erzeugt, der nichts prüft — `modul-06` nennt das
Zeremonie. Sein Zustand ist die Verzeichnis-Position, seine Belege stehen hier
und in `git`.

**Offen und benannt:** Docker-Basis-Images und Action-Pins haben weiterhin keine
Freshness-Achse. Ob sie eine brauchen, ist eigens zu entscheiden — der Sensor
ist parametriert, das Hinzufügen wäre je Achse eine Zeile.
