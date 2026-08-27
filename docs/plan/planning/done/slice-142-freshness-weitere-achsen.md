# Slice slice-142: Zwei ungewachte Pin-Klassen — Basis-Images und Action-Pins

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [slice-137](../done/slice-137-toolchain-freshness.md); `tools/harness/pin-freshness.sh`; [`AGENTS.md`](../../../../AGENTS.md) §3.9, §4.

**Berührte Spec-Stellen:** — (Harness-Sensoren; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

Drei gepinnte Fremd-Bestände sind gewacht (Baseline, Go, golangci-lint). Zwei
Klassen sind es nicht: die **Docker-Basis-Images** (digest-gepinnt im
`Dockerfile`) und die **Action-Pins** der Workflows, deren *Form* seit
`make workflow-pins` geprüft ist, deren **Gültigkeit** aber nicht.

Der Sensor ist parametriert; je Achse wäre es eine Zeile. **Ob sie eine Achse
brauchen, ist trotzdem eine Entscheidung** — ein Sensor, den niemand liest, ist
Aufwand ohne Wirkung.

## 2. Vorgehen

1. **Je Klasse messen, bevor gebaut wird:** wie oft hat sich der Pin
   tatsächlich bewegt, und was wäre die Folge einer verpassten Hebung?
2. **Quellen-Form klären.** Ein Digest hat keinen Release-Tag; die Frage
   *„gibt es einen neueren Digest für denselben Tag"* ist eine andere als
   *„gibt es einen neueren Tag"*.
3. Nur bauen, was die Messung trägt — und die Entscheidung gegen eine Achse
   **ausweisen**, nicht weglassen.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Auto-Bump.** Wie bei den drei bestehenden Achsen: melden, nicht heben.
- **Keine Aufnahme in `gates`.** Netzhaltiges gehört in den Nachtlauf.

## 4. Definition of Done

- [x] Je Klasse eine **Messung** und daraus eine begründete Entscheidung.
- [x] Was gebaut wird, hängt im bestehenden Nachtlauf und ist fail-open.
- [x] Was **nicht** gebaut wird, ist mit Grund ausgewiesen.
- [x] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Eine Achse ist billig zu bauen und teuer zu ignorieren.** Wächst die Zahl
  der Nachtlauf-Meldungen, sinkt die Aufmerksamkeit für jede einzelne. —
  **Ausgang: eingetreten.** Der Nachtlauf wuchs von vier Achsen auf zwölf, und
  **sechs** melden beim ersten Lauf Rückstand. Ein Nachtlauf, der jede Nacht
  dasselbe Rot zeigt, verliert das Signal für die Achse, die *neu* rot wird —
  genau der verwaiste Sensor, gegen den er gebaut wurde. Das Heben ist ein
  bewusster Akt und damit nicht dieser Slice: geschnitten als
  [slice-161](../open/slice-161-sechs-pins-heben.md).
- **Digest-Achsen haben keine Tag-Semantik.** Die Gleich/Ungleich-Form der
  bestehenden Achsen trägt hier womöglich nicht. — **Ausgang: entfallen,
  gemessen.** Sie trägt sogar besser: bei Tags musste eigens begründet werden,
  dass Gleich/Ungleich statt eines Semver-Sorts genügt (beide Reihen sind
  monoton). Ein Digest hat gar keine Ordnung — „anders" ist die **einzige**
  sinnvolle Aussage, und genau die macht der vorhandene Vergleicher. Was die
  Form doch kostete, ist das **Urteilswort**: `VERALTET` behauptete eine
  Richtung, die es dort nicht gibt; die Digest-Achsen melden `ABWEICHEND`.
- **Die Klassen-Definition könnte Mitglieder auslassen.** „Basis-Images" und
  „Action-Pins" sind Namen, keine Aufzählung. — **Ausgang: eingetreten,
  behoben.** Die erste Fassung ließ **vier** Pins aus: zwei Dockerfile-Stages
  (mit der falschen Zusage, sie seien „über die Achsen darüber gewacht"), das
  semgrep-Gate-Image und das a-check-Image. Beim a-check-Pin war der Grund
  strukturell — seine Version stand nur im **Kommentar**, wo kein Sensor sie
  liest; sie steht jetzt als `A_CHECK_VERSION` im Fragment.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung zeigt, dass eine Klasse eine andere Vergleichsform braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Werkzeuge (GF), CI (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25): [`BEO-011`](../observations.md) für die Bewegungs-Messung; [`BEO-007`](../observations.md) für jeden Beleg-Lauf.

Slice-ID: slice-142. Betroffene IDs: — (Harness-Sensoren; keine Anforderung). Module: Harness-Werkzeuge, CI.
Gates: `make gate-consistency`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Erweiterung an bereits gebauter Mechanik.

## 9. Closure-Notiz (nach `done/`)

**Zwölf Achsen statt vier — und die Entscheidung, die dieser Slice zuerst
traf, stand auf einer Messung der falschen Menge.**

**Die erste Messung war vollständig und richtig, nur über den falschen
Bestand.** Gezählt wurde, wie viele verschiedene Werte jeder Pin **in unserer
Historie** trug — also wie oft **wir** gehoben haben. Daraus las sich
*„Action-Pins: null Bewegungen"* und *„`distroless`: nie bewegt"* wie eine
Aussage über upstream, und sie trug eine Entscheidung: die Action-Achse wurde
**abgelehnt**. Ein Bestand, den nur wir ändern, kann über einen Bestand, den
nur andere ändern, nichts sagen. Als Klasse im Register:
[`BEO-020`](../observations.md).

**Die direkte Messung — upstream fragen — kippte beide tragenden Sätze.** Die
Ablehnung stand auf zweien: *„braucht eine eigene Quellen-Form"* und *„ein
alter SHA ist kein Sicherheitsproblem"*. Der erste ist falsch, weil der
**Tag-Kommentar** neben dem SHA genau die Größe trägt, die der vorhandene
`releases/latest`-Zweig vergleicht — die Achse ist ein Extraktor, keine neue
Quelle. Der zweite ist falsch, weil das Pinnen zwar das Umhängen eines Tags
ausschließt, aber zugleich **blind für dessen Behebung** macht: `checkout`
v6.0.2 steht gegen v7.0.1, `login-action` v4.2.0 gegen v4.6.0.

**Und die Zusage, die zwei übrigen Basis-Images seien „über die Achsen darüber
gewacht", war ebenfalls falsch — gemessen.** `make freshness-go` meldet `ok`,
während `golang:1.27.0` einen **anderen Digest** trägt als unser Pin. Eine
Tag-Achse wacht die **Version**, nicht den **Bau**; wer beide Pins führt,
braucht beide Fragen. Das ist keine Feinheit, sondern der Grund, warum die
Digest-Form überhaupt existiert.

**Gebaut sind deshalb zwölf Achsen, in drei Formen:**

| Form | Frage | Gegenstand |
|---|---|---|
| Baseline-Skript | neuerer Tag **+** Content-Drift | vendorte Baseline |
| Versions-Achsen | „gibt es einen neueren Tag" | Go, `golangci-lint`, `semgrep`, `a-check`, `checkout`, `login-action` |
| Digest-Achsen | „trägt derselbe Tag einen anderen Digest" | die drei `Dockerfile`-Stages, `semgrep`, `a-check` |

**Sechs melden beim ersten Lauf Rückstand** — beide Action-Pins, `semgrep`
(1.167.0 gegen 1.175.0), `a-check` (v0.8.0 gegen v0.17.0) und zwei
Basis-Image-Digests. Gehoben wird hier nichts: §3 sagt *melden, nicht heben*.
Dass der Nachtlauf damit **dauerrot** startet, ist der eingetretene §5-Punkt und
als [slice-161](../open/slice-161-sechs-pins-heben.md) geschnitten.

**Drei Änderungen, die das Bauen erzwang, und jede schließt eine eigene Lücke.**
Der `a-check`-Pin trug **keinen Tag** — seine Version stand nur im Kommentar,
wo kein Sensor sie lesen kann; sie steht jetzt als `A_CHECK_VERSION` im
Fragment, und die Referenz führt Tag **und** Digest wie die
`Dockerfile`-Stages. Das `v`-Präfix wird jetzt **symmetrisch auf beiden Seiten**
gestrippt: unsere Pins führen es uneinheitlich (`v2.13.1`, aber `1.167.0`), und
der Kopf-Kommentar behauptete zuvor, das Strippen mache den Vergleich
großzügiger — solange es symmetrisch geschieht, macht es ihn **richtig**. Und
die Digest-Achsen melden `ABWEICHEND` statt `VERALTET`, weil Digests keine
Ordnung haben.

**Vier Kanten am Sensor, beim Einbau gefunden.** Der Werkzeug-Riegel für `curl`
stand **vor** dem Dispatch — für die Digest-Form wäre er ein SKIP wegen eines
Werkzeugs gewesen, das der Zweig gar nicht ruft, also ein stilles Abschalten aus
dem falschen Grund; er steht jetzt je Zweig. `imagetools inspect` kennt keine
eigene Zeitgrenze, deshalb steht ein `timeout` davor — die Fail-open-Zusage des
Kopfes gilt für **jeden** Zweig, und ohne sie wäre eine hängende Verbindung ein
Job-Timeout statt eines SKIP. Die **Pin-Seite** wird jetzt ebenso auf ihre Form
geprüft wie die Upstream-Seite: der Pin kommt aus einer Textextraktion am
`Dockerfile` und ist die fragilere Hälfte. Und die `ok`-Meldung sagte *„ist der
neueste Stand"* — auch das eine Ordnungs-Behauptung; sie sagt jetzt *„entspricht
dem Upstream-Stand"*.

**Was offen bleibt, benannt.** Die Action-Achsen vergleichen den
**Tag-Kommentar** gegen upstream, nicht den **SHA** gegen den Kommentar. Ob der
gepinnte SHA den Commit bezeichnet, den sein Kommentar behauptet, ist eine
andere Frage; sie bleibt die benannte Grenze von `make workflow-pins`
([`AGENTS.md`](../../../../AGENTS.md) §3.9). Und keine Achse trägt eine
**Kadenz**: sie melden jede Nacht, aber wer das Rot wann liest, steht nirgends —
als dritter §5-Punkt von [slice-161](../open/slice-161-sechs-pins-heben.md)
geführt.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 545 Dateien, 0 Befunde),
`make workflow-pins` (4 `uses:`-Einträge, alle SHA-gepinnt mit Tag-Kommentar),
`make arch-check` (0 Befunde, gegen den getaggten a-check-Pin). Alle zwölf
Achsen sind gefahren, mit gelesener Ausgabe; der netzlose `--compare`-Einstieg
ist gegen beide Urteilsworte geprüft. Ein unabhängiger Review ist gelaufen —
sein Urteil war *„nicht schließbar in der vorliegenden Form"*, und die zwei
blockierenden Befunde sind die zwei oben widerlegten Entscheidungs-Aussagen.
