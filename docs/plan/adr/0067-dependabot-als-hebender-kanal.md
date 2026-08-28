# ADR-0067: Dependabot als hebender Kanal — mit Kennung statt mit Gate-Ausnahme

**Status:** Accepted

**Datum:** 2026-08-28

**Autor:** pt9912

**Bezug:** [ADR-0066](0066-cve-scan-gegen-das-publizierte-image.md) (der Sensor,
dessen Befunde niemand hebt),
[ADR-0011](0011-digest-pins-build-gate-images.md) (Digest-Pins),
[ADR-0027](0027-commits-traceability-modul.md) (das Traceability-Gate)

**Schärft:** — (CI-Konfiguration ohne Spec-Stratum; Präzedenz
[ADR-0010](0010-semgrep-hermetisches-gate.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0066](0066-cve-scan-gegen-das-publizierte-image.md) hat einen Sensor
gebaut, der beim ersten Lauf **vierzehn behebbare HIGH-Befunde** im
ausgelieferten Image fand. Gehoben hat sie ein Mensch. Hebt beim nächsten Satz
niemand, meldet der Sensor dieselben Befunde jede Nacht weiter — und die Klasse
*„Sensor ohne Adressat"* ist in diesem Repo bereits einmal geschnitten worden.

**Die Kollision steht vor der Konfiguration, nicht hinter ihr.**
`make trace-check` läuft in der PR-CI über **jede** Commit-Message und verlangt
eine `DC-`/`ADR-`/`MR-`/`slice-`-Kennung
([ADR-0027](0027-commits-traceability-modul.md)). Dependabots Botschaften tragen
keine. **Gemessen gegen das echte Gate:**

| Botschaft | `make trace-check` |
| --------- | ------------------ |
| `build(ci): bump actions/checkout from 7.0.1 to 7.1.0` | **rot** |
| `build(deps) [ADR-0067]: bump golang.org/x/net from 0.56.0 to 0.57.0` | **grün** |

Ohne diese Frage vorweg wäre **jeder** Dependabot-PR ab dem ersten Tag rot — ein
Kanal, den man abstellt, statt ihn zu benutzen.

## Entscheidung

Wir wählen **Dependabot für `gomod` und `github-actions`, mit der ADR-Kennung im
Commit-Präfix**.

1. **Die Kennung ist eine Zusage, keine Umgehung.** `commit-message.prefix`
   trägt `[ADR-0067]` — die Kennung **dieses** Entscheids. Die ADR, die
   Dependabot erlaubt, ist der ehrliche Anker für die Commits, die daraus
   entstehen. `commits.exempt-pattern` bleibt **unverändert**.
2. **`gomod`: ja.** Dort lagen die vierzehn Befunde — neun in
   `golang.org/x/crypto`, vier in `golang.org/x/net`, eine in `go-git`. Das ist
   das Ökosystem, für das dieser Kanal existiert.
3. **`github-actions`: ja.** Fremde Actions laufen mit den Zugangsdaten dieses
   Repositories; ihre Aktualisierungen sind sicherheitsrelevant. Die drei
   Frische-Achsen **melden** sie bereits — Dependabot **hebt** sie. Das ist
   keine Doppelung, sondern die zweite Hälfte.
4. **`docker`: nein, und das ist der wichtigste Ausschluss.** Die Basis-Images
   sind **digest**-gepinnt ([ADR-0011](0011-digest-pins-build-gate-images.md))
   und tragen ihre Version in einer Make-Variablen. Ein Dependabot-PR höbe den
   Tag und ließe den Digest stehen — er bräche genau die Kopplung, die den Pin
   trägt. Für diese Frage stehen **fünf** Digest-Achsen im Nachtlauf.
5. **Gruppiert, wöchentlich.** Ein PR je Action bei sechs `uses:`-Einträgen ist
   kein Kanal, sondern eine Flut.
6. **Kein Automerge.** Dependabot öffnet PRs; geprüft wird mit `make ci`, und
   der Merge bleibt ein Akt — dieselbe Linie wie bei den Frische-Achsen, wo die
   Hebung ausdrücklich ein bewusster Akt ist.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, von Hand heben | kein neuer Kanal, volle Kontrolle | der Sensor meldet dann jede Nacht dasselbe; die Hebung hängt daran, dass jemand hinsieht — genau die Lücke, die ADR-0066 offen lässt |
| B — `commits.exempt-pattern` um Dependabot erweitern | eine Zeile, sofort grün | eine **gelockerte Prüfregel** ([`AGENTS.md`](../../../AGENTS.md) §3.6); der Gate verlöre seine Aussage für eine ganze Commit-Klasse, und die Ausnahme überlebte ihren Grund |
| C — Dependabot **ohne** `gomod`, nur `github-actions` | die Actions sind die riskantere Fläche | ließe genau das Ökosystem draußen, in dem die vierzehn Befunde lagen |
| D — zusätzlich `docker` | ein Kanal für alle Fremd-Bestände | bräche die Digest-Pin-Kopplung (Entscheidung 4); die fünf Digest-Achsen decken die Frage bereits |
| **E — `gomod` + `github-actions`, Kennung im Präfix (gewählt)** | der Kanal hebt genau dort, wo gemeldet wird; kein Gate gelockert; die Kennung ist gemessen wirksam | die Kennung behauptet für jeden künftigen Bump einen Bezug zu **diesem** Entscheid — wahr für den Kanal, nicht für den Inhalt |

## Konsequenzen

- **Positiv:** die Befunde des Sensors bekommen einen Weg, der nicht daran
  hängt, dass jemand hinsieht.
- **Positiv:** kein Gate ist gelockert. Die Lösung macht die neue Commit-Klasse
  gate-**konform**, statt sie auszunehmen.
- **Negativ:** die Kennung im Präfix sagt *„gehört zu ADR-0067"* — das gilt dem
  **Kanal**, nicht dem Inhalt des einzelnen Bumps. Wer sie liest, könnte mehr
  Bezug annehmen, als da ist; deshalb steht der Satz hier.
- **Negativ, offen:** Dependabot sieht `go.mod`, nicht das ausgelieferte Bild.
  Die vierzehn Befunde lagen in **indirekten** Abhängigkeiten; ob es sie ohne
  Security-Advisory anfasst, ist **nicht gemessen** und wird hier nicht
  behauptet. Der Sensor bleibt die Wahrheit über das Artefakt.
- **Folgepflicht:** [`AGENTS.md`](../../../AGENTS.md) §5 nennt die
  Dependabot-Form, damit die Kennung nicht wie eine Ausnahme aussieht; die
  `ignore`-Einträge tragen ihren Grund **in der Datei**.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Modul `commits` | jede Commit-Message nennt eine `DC-`/`ADR-`/`MR-`/`slice-`-Kennung — für Dependabot über `commit-message.prefix`, nicht über eine Ausnahme | `make trace-check` (PR-CI, `commit-msg`-Hook) |
| `make ci` | jeder Dependabot-PR durchläuft dieselben Gates wie jeder andere | `make ci` (PR-CI) |

## Re-Evaluierungs-Trigger

Öffnet Dependabot über **drei** Monate keinen einzigen PR, den der Sensor nicht
schon gemeldet hatte, ist der Kanal überflüssig — dann kostet er Rauschen ohne
eigenen Beitrag.

Zweiter Trigger: fasst Dependabot die **indirekten** Go-Abhängigkeiten nicht an,
in denen die bisherigen Befunde lagen, ist die Ökosystem-Wahl neu zu stellen —
dann trägt `gomod` hier nicht, was diese Entscheidung ihm zuschreibt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-08-28 | Accepted — die Traceability-Kollision war vor der Konfiguration gemessen und ohne Gate-Lockerung gelöst | [slice-167](../planning/in-progress/slice-167-dependabot.md) |
