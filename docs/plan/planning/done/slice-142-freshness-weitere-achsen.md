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
  **Ausgang: entfallen.** Der Nachtlauf wächst um **eine** Achse, nicht um
  zwei; die zweite ist mit Messung abgelehnt. Und die eine gebaute meldet beim
  ersten Lauf einen echten Rückstand — sie beginnt nicht als Rauschen.
- **Digest-Achsen haben keine Tag-Semantik.** Die Gleich/Ungleich-Form der
  bestehenden Achsen trägt hier womöglich nicht. — **Ausgang: entfallen,
  gemessen.** Sie trägt sogar besser: bei Tags musste eigens begründet werden,
  dass Gleich/Ungleich statt eines Semver-Sorts genügt (beide Reihen sind
  monoton). Ein Digest hat gar keine Ordnung — „anders" ist die **einzige**
  sinnvolle Aussage, und genau die macht der vorhandene Vergleicher. Die neue
  Form brauchte deshalb keine neue Vergleichs-Semantik, nur eine neue Quelle.

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

**Eine Achse gebaut, eine mit Messung abgelehnt — und die gebaute war beim
ersten Lauf schon fällig.**

**Die Messung zuerst, wie §2 es verlangt.** Gezählt wurde, wie oft sich jeder
Pin seit seiner Einführung bewegt hat:

| Klasse | verschiedene Werte in der Historie | gewacht |
|---|---|---|
| Go-Basis-Image | **drei** Digests | Tag ja, Digest bei gleichem Tag nein |
| Lint-Basis-Image | **zwei** Digests | dito |
| `distroless/static-debian12:nonroot` | **einer**, nie bewegt | **gar nicht** |
| Action-Pins | **null** Bewegungen | Form ja, Gültigkeit nein |

**Die Prämisse des Slice war falsch.** §1 sagte, der Sensor sei parametriert und
je Achse *„wäre es eine Zeile"*. Er trägt zwei Quellen-Formen — beide
beantworten *„gibt es einen neueren Tag"*. Ein Digest hat keinen Tag, ein
Action-SHA keinen Release: **beide** Kandidaten brauchen eine **dritte
Quellen-Form**, keinen Parameter.

**Gebaut: die Digest-Achse, und nur für einen Pin.** Der Tag
`distroless/static-debian12:nonroot` trägt **keine Version** — für ihn gibt es
keine Tag-Frische-Achse, weil es keinen Tag gibt, der sich bewegt; der Digest
ist die einzige Handhabe. Die beiden anderen Basis-Images hängen an den zwei
Toolchain-Variablen und sind über die vorhandenen Achsen gewacht; ihre
Restlücke — ein Neubau **desselben** Tags — ist damit kleiner und schließt sich
bei jedem Bump von selbst.

**Der erste Lauf meldet `VERALTET`:** Pin `sha256:d093aa3e…`, upstream
`sha256:afa5c872…`. Der Tag wurde neu gebaut, und zweieinhalb Monate lang hat
das niemand gesehen. Gehoben wird hier nichts — §3 sagt *melden, nicht heben*,
wie bei den drei bestehenden Achsen.

**Nicht gebaut: die Action-Achse, und der Grund ist gemessen.** Kein
Action-Pin hat sich je bewegt, und die Gefahr, gegen die ein Pin schützt — ein
umgehängter Tag — ist durch das Pinnen strukturell ausgeschlossen;
`make workflow-pins` hält die Form. Was bliebe, ist *„der SHA ist alt"*, und das
ist kein Sicherheitsproblem, sondern ein Aktualitäts-Wunsch ohne Dringlichkeit.
Eine Achse, die nie meldet, senkt die Aufmerksamkeit für die, die melden.

**Was die Ablehnung offen lässt:** die vom Slice genannte *Gültigkeit* — ob der
gepinnte SHA den Commit bezeichnet, den sein Tag-Kommentar behauptet. Das ist
eine andere Frage als Frische, sie bräuchte eine vierte Quellen-Form, und sie
ist hier weder gemessen noch entschieden. [`AGENTS.md`](../../../../AGENTS.md)
§3.9 führt sie bereits als benannte Grenze von `make workflow-pins`; dort
bleibt sie.

**Zwei Nebenbefunde am Sensor, beim Einbau gefunden.** Der Werkzeug-Riegel für
`curl` stand **vor** dem Dispatch — für die Digest-Form wäre er ein SKIP wegen
eines Werkzeugs gewesen, das der Zweig gar nicht ruft, also ein stilles
Abschalten aus dem falschen Grund. Er steht jetzt je Zweig. Und der
Kopf-Kommentar sagte „zwei Quellen-Formen" sowie eine Werkzeug-Liste ohne
Docker; beides ist nachgezogen.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 544 Dateien, 0 Befunde) —
`gate-consistency` grün belegt, dass das neue Target in `Makefile`,
[`AGENTS.md`](../../../../AGENTS.md) §4 und
[`harness/README.md`](../../../../harness/README.md) §Sensors konsistent
deklariert ist. Die neue Achse ist gefahren, mit gelesener Meldung.
