# Slice slice-161: Sechs gemeldete Pin-Rückstände heben — der Nachtlauf ist dauerrot

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [slice-142](../done/slice-142-freshness-weitere-achsen.md) (die Achsen, die melden); [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md) (Digest-Pin-Politik); [ADR-0029](../../adr/0029-arch-check-via-a-check.md) (a-check-Fragment); [`AGENTS.md`](../../../../AGENTS.md) §3.9, §4.

**Berührte Spec-Stellen:** — (Pin-Stand; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

Der Nachtlauf wacht seit [slice-142](../done/slice-142-freshness-weitere-achsen.md)
über **zwölf** Achsen. **Sechs davon melden beim ersten Lauf Rückstand:**

| Achse | Pin | upstream |
|---|---|---|
| `checkout-pin-freshness` | `v6.0.2` | `v7.0.1` |
| `login-pin-freshness` | `v4.2.0` | `v4.6.0` |
| `freshness-semgrep` | `1.167.0` | `1.175.0` |
| `freshness-a-check` | `v0.8.0` | `v0.17.0` |
| `runtime-base-digest` | `sha256:d093aa3e…` | anderer Bau desselben Tags |
| `go-base-digest` | `sha256:65b6f280…` | anderer Bau desselben Tags |

Ein Nachtlauf, der **jede Nacht dasselbe Rot** zeigt, ist derselbe verwaiste
Sensor, gegen den er gebaut wurde: er verliert das Signal für die nächste Achse,
die *neu* rot wird. Die Achsen melden nur — das Heben ist ein bewusster Akt, und
dieser Slice ist er.

## 2. Vorgehen

1. **Je Pin eine eigene Entscheidung**, nicht ein Sammel-Bump. Ein
   Major-Sprung (`checkout` v6→v7, `a-check` v0.8→v0.17) ist eine andere
   Entscheidung als ein Digest-Nachzug am selben Tag.
2. **Form je Klasse aus [`BEO-008`](../observations.md):** eine Pin-Hebung hat
   Spiegel. Vor dem Bump je Pin die Spiegel zählen, nicht nur die grep-baren.
3. **Der `a-check`-Sprung ist der teuerste** — neun Minor-Releases, und das Gate
   fährt gegen unser Repo. Erst gegen den neuen Digest laufen lassen, dann
   pinnen; ein rotes `arch-check` aus fremder Regel-Verschärfung ist eine
   eigene Entscheidung, kein Nebeneffekt.
4. **Nach jedem Bump die zugehörige Achse fahren** und ihre echte Ausgabe lesen.
5. `make gates`, `make ci`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Auto-Bump-Mechanismus.** Die Achsen melden; wer hebt, entscheidet.
- **Keine Aufnahme netzhaltiger Achsen in `gates`.** Der Nachtlauf bleibt der
  Bindepunkt.
- **Keine neuen Achsen.** Die Menge ist mit slice-142 abgeschlossen.

## 4. Definition of Done

- [x] Je gemeldeter Pin **entschieden**: alle sechs gehoben, je mit eigener
      Begründung.
- [x] Je Hebung die Spiegel aus [`BEO-008`](../observations.md) mitgezogen —
      **drei** lebende Fundstellen, nicht zwei (die dritte fand erst der Review).
- [x] Der Nachtlauf ist nach dem Slice **nicht mehr dauerrot**: alle zwölf
      Achsen melden `ok`, einzeln gefahren und vom Review nachgefahren.
- [x] `make gates` und `make ci` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Major-Sprung kann ein Gate rot machen.** `a-check` v0.8→v0.17 sind neun
  Releases fremder Regel-Entwicklung; `actions/checkout` v6→v7 kann
  Runner-Verhalten ändern. — **Ausgang: entfallen, gemessen.** `a-check`
  v0.17.0 meldet über dieses Repo 0 Befunde **und** dieselben Befundzeilen wie
  v0.8.0 über eine Proben-Matrix aus **sieben** Verbotszweigen und **drei**
  Allow-Gegenproben — das Grün kommt aus dem Bestand. Und der v7-Bruch trifft
  ausschließlich Workflows, die Fork-PRs unter `pull_request_target` oder
  `workflow_run` auschecken; keiner unserer drei nutzt einen der Trigger.
- **Ein Digest-Nachzug am selben Tag ist nicht folgenlos.** Der Tag ist
  derselbe, der Inhalt nicht; ein neu gebautes Basis-Image kann eine andere
  libc-Ecke tragen. — **Ausgang: eingetreten, in der einen Hälfte.** Nicht
  libc — `distroless/static` bringt keine mit, und das Binary ist
  `CGO_ENABLED=0`. Wohl aber der **CA-Trust-Store**: das Wurzel-Bündel des
  ausgelieferten Images wächst von **142 auf 150** Zertifikate (216 591 →
  224 449 Byte). Wer das Modul `external` gegen HTTPS-Ziele fährt, vertraut ab
  diesem Image einer anderen Menge von Wurzeln. Der erste Anlauf hat das nicht
  gemessen und den Nachzug als Nicht-Ereignis behandelt.
- **Das Heben schließt die Lücke nicht, es verschiebt sie.** Morgen ist der
  nächste Release da. Was dieser Slice nicht leistet, ist eine **Kadenz** — wer
  das Rot wann liest. — **Ausgang: weiter offen.** Der DoD-Zweig
  *„nicht mehr dauerrot"* ist erfüllt, also greift der Adressaten-Zweig nicht;
  die Lücke steht seit slice-142 als benannte Grenze im Kopf von
  [`upstream-drift.yml`](../../../../.github/workflows/upstream-drift.yml).
  Fällig ist sie trotzdem, und **jetzt** ist sie billig: geschnitten als
  [slice-164](../in-progress/slice-164-nachtlauf-kadenz.md).

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls ein Bump ein Gate rot macht und
die Ursache eine eigene Entscheidung braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Werkzeuge (GF), CI (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27): [`BEO-008`](../observations.md) für die Pin-Spiegel; [`BEO-020`](../observations.md) für jede Aussage darüber, wie oft sich ein Pin bewegt; [`BEO-007`](../observations.md) für jeden Beleg-Lauf.

Slice-ID: slice-161. Betroffene IDs: [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md), [ADR-0029](../../adr/0029-arch-check-via-a-check.md). Module: Harness-Werkzeuge, CI.
Gates: `make arch-check`, `make semgrep`, `make gates`, `make ci`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Pflege an bereits gebauter Mechanik.

## 9. Closure-Notiz (nach `done/`)

**Sechs Pins gehoben, zwölf von zwölf Achsen grün — und die erste Zeile der
Commit-Botschaft war trotzdem falsch.**

**Die Sache selbst.** Je Pin eine eigene Entscheidung, wie §2 es verlangt:
`a-check` v0.8.0 → v0.17.0, `semgrep` 1.167.0 → 1.175.0, `actions/checkout`
v6.0.2 → v7.0.1, `docker/login-action` v4.2.0 → v4.6.0 und die zwei
Basis-Image-Digests. Danach melden alle zwölf Achsen `ok` — einzeln gefahren,
Ausgabe gelesen, und vom unabhängigen Review nachgefahren.

**Beim teuersten Pin hätte Grün nicht gereicht.** `a-check` v0.17.0 meldet über
dieses Repo 0 Befunde; das könnte auch heißen, dass die neue Fassung weniger
prüft. Die Gegenprobe war deshalb ein **konstruierter Verstoß**, den beide
Fassungen an derselben Zeile melden. **Der Review hat gezeigt, dass diese Probe
zu schmal war:** [ADR-0029](../../adr/0029-arch-check-via-a-check.md) verlangt
eine **Matrix je Verbotszweig**, und
[`harness/README.md`](../../../../harness/README.md) trägt sie als **lebende**
Zusage. Die Matrix ist nachgeholt — **sieben** Verbotszweige, **drei**
Allow-Gegenproben, byte-gleiche Befundzeilen. Sie ist die Messung des Reviews,
nicht meine; der Beleg liegt im
[Report](../../../reviews/2026-08-27-slice-161-sechs-pins-heben-review.md).

**Der Major-Sprung war geprüft, nicht gewagt.** Die v7.0.0-Release-Notes nennen
genau einen Bruch — das Auschecken von Fork-PRs unter `pull_request_target` und
`workflow_run` —, und keiner unserer drei Workflows nutzt einen der beiden
Trigger. Der Review ist weiter gegangen als ich (Notes von v7.0.0 **und**
v7.0.1, `action.yml`-Diff, Quelltext des Sicherheits-Helfers) und fand keinen
verschwiegenen zweiten Bruch; `node24` lief schon in v6.0.2.

**Die falsche Zeile, und sie steht in einem gepushten Commit.** Die Botschaft
beginnt mit *„Der Nachtlauf war seit slice-142 dauerrot"*. Gemessen gab es
**genau einen** roten Lauf, rund zwei Stunden vor dem Commit; die drei
Nachtläufe davor waren grün, weil die zwölf Achsen erst am selben Vormittag
landeten. Der Slice-Plan §1 sagt es richtig — *„melden **beim ersten Lauf**
Rückstand"* —, und ich habe daraus eine Dauerzustands-Behauptung gemacht.
Korrigieren durch Amenden ging nicht: der Commit war beim Review-Ergebnis
bereits auf `origin/main`. Die richtige Zahl steht deshalb hier, wo sie beim
nächsten Lesen zuerst auffällt. Klasse: [`BEO-020`](../observations.md) —
gemessen war *ein* Lauf, geredet wurde über eine Dauer.

**Ein dritter Pin-Spiegel war übersehen.** Der Zensus nach
[`BEO-008`](../observations.md) fand zwei lebende Fundstellen; es waren drei.
Die dritte stand im Kopf-Kommentar von
[`pin-freshness.sh`](../../../../tools/harness/pin-freshness.sh) und erklärte
die `v`-Normalisierung mit **demselben Beispielwert**, der eine Datei weiter
gehoben wurde — nach dem Bump widersprachen sich beide Stellen. Statt den Wert
zu heben, ist das Beispiel **entfernt**: es war ein vierter Spiegel, der beim
nächsten Bump wieder still driften würde. Die Erklärung trägt ohne ihn.

**Eine Prozedur im eigenen Kopf war übersprungen.**
[`a-check.mk`](../../../../a-check.mk) sagt: *Pin-Hebung ist ein bewusster
Commit; dabei das Fragment per `--print-mk` neu erzeugen.* Getan war nur der
Variablen-Tausch. Das v0.17.0-Fragment trägt drei Neuerungen — eine
Runtime-Indirektion, ein Target `a-check-graph`, die erweiterte
`.PHONY`-Zeile. **Keine ist adoptiert**, und beide Gründe stehen jetzt im Kopf
statt nirgends: die Indirektion zahlt sich nur repo-weit aus, und ein neues
Target ist gate-consistency-pflichtig, also ein eigener Entscheid. Benannter
Kandidat, kein Versehen.

**Eine ADR ist durch die Hebung teilweise überholt — und das ist kein
ADR-Fall.** [ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md)
§Entscheidung nennt die gepinnte Fassung im **Indikativ**.
[ADR-0011](../../adr/0011-digest-pins-build-gate-images.md) §4 erklärt die
Hebung ausdrücklich zum bewussten Commit, und §3.5 lässt `## Geschichte`-Anhänge
zu; eine Zeile dort hält fest, was gilt und was nicht mehr. Der Review hat
dabei die Unterscheidung benannt, die im Zensus fehlte: *eingefroren und weiter
wahr* (die `done/`-Slices, die Reporte, die Geschichte-Zeile von [ADR-0029](../../adr/0029-arch-check-via-a-check.md), die
Mindest-Versions-Aussagen in [`.a-check.yml`](../../../../.a-check.yml)) gegen
*eingefroren und jetzt falsch*. Nur die zweite Klasse braucht eine Handlung.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 554 Dateien, 0 Befunde,
Coverage 94,90 %), `make ci` (Exit 0, `image-test` gegen das neu gebaute
Runtime-Image), **alle zwölf Achsen einzeln** (je Exit 0, Ausgabe gelesen),
`make adr-check` über die Range (0 Befunde). Der Trust-Store-Befund ist eigens
nachgemessen: beide Images exportiert, Dateilisten und Zertifikatszahl
verglichen. Ein unabhängiger Review ist gelaufen; sein Urteil war *„schließbar
nach Nacharbeit"*, und seine acht Befunde sind eingearbeitet — das eine HIGH
ist die falsche erste Zeile oben.
