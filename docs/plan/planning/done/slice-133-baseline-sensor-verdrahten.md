# Slice slice-133: Der Baseline-Sensor ist gebaut und wird nie aufgerufen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-84-durchsetzung](../welle-84-durchsetzung.md).

**Bezug:** [`tools/harness/fetch-baseline-cache.sh`](../../../../tools/harness/fetch-baseline-cache.sh);
[`AGENTS.md`](../../../../AGENTS.md) §4 (Gate-Tabelle — Autorität über die
Targets) und §3.1 (Docker/make-only);
[`MR-011`](../../../../harness/conventions.md#mr-011)-Pin-Kette;
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(`make gates` bleibt netzlos).

**Berührte Spec-Stellen:** — (Harness-Verdrahtung; keine Anforderung ändert ihre
Aussage).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

`fetch-baseline-cache.sh` trägt drei Fähigkeiten:

- **`--verify`** — `sha256sum -c` über `SHA256SUMS` **plus** Manifest-Deckung
  (Dateien auf Platte == Manifest-Zeilen). Die zweite Hälfte ist die tragende:
  ohne sie passiert eine untermengige, in sich konsistente `SHA256SUMS` grün.
  **Netzlos.**
- **`--check-latest (A)`** — Currency gegen den Pin, gelesen aus der
  Release-**Liste**, ausdrücklich nicht aus `/releases/latest`, „der Prereleases
  überspringt und einen zurückgezogenen Pin verbirgt". **Netz.**
- **`--check-latest (B)`** — Content-Drift am gepinnten Tag: die Bytes beider
  Bäume des Release-Assets gegen das committete `SHA256SUMS`. **Netz.**

**Kein `make`-Target, kein Workflow und kein Hook ruft das Skript auf** —
gemessen über das ganze Repo. Die Folge ist nicht theoretisch: dass der Kurs
`v5.11.0` veröffentlicht hatte, hat uns der Auftraggeber gesagt, nicht der
Sensor, der dafür gebaut wurde.

Dieser Slice steckt ihn ein. Er baut **nichts Neues**.

## 2. Vorgehen

1. **`--verify` als Gate**: eigenes Target, Glied von `make gates`. Es ist
   netzlos — sein Pfad ruft kein Netz-Werkzeug. Dass `make gates` offline bleibt,
   ist dabei eine Eigenschaft **dieses Repos**;
   [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
   gilt dem **Produkt** und wird hier nicht berührt.
2. **`--check-latest` als eigenes Target**, ausdrücklich **nicht** in `gates`
   (Netz), analog `trace-check`/`adr-check`, die aus demselben Grund draußen
   stehen.
3. **Nachtlauf**: ein eigener Workflow neben `ci.yml`, `schedule:` täglich. Er
   ist **von `ci.yml` getrennt**, damit ein Upstream-Ausfall nie die CI rot
   färbt.
4. **`AGENTS.md` §4 und die Sensors-Tabelle in `harness/README.md` nachziehen** —
   `make gate-consistency` prüft Doku↔Makefile in **beide** Richtungen
   (`gate-phantom`/`gate-undocumented`); ohne den Nachzug ist der Slice
   gate-rot.
5. **Bewusstes Brechen**: eine Datei im vendorten Baum verändern ⇒ `--verify`
   rot; eine Datei zusätzlich einlegen ⇒ **ebenfalls** rot (die
   Manifest-Deckung); Rückbau ⇒ grün. Beide Richtungen, weil nur die zweite die
   Zusage trägt, die `sha256sum -c` allein nicht hält.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Verallgemeinerung auf die Toolchain-Pins.** `GO_VERSION` und
  `GOLANGCI_LINT_VERSION` bleiben ungewacht; sie stehen im `Dockerfile` und
  damit außerhalb der gescannten Menge — ihr Sensor ist eine eigene
  Entscheidung, die der Zensus schneidet.
- **Kein Auto-Bump.** Der Sensor **meldet**; die Hebung bleibt ein bewusster
  Akt ([`MR-011`](../../../../harness/conventions.md#mr-011)-Kette).
- **Kein Fail-closed im Netz-Teil.** Netz-, Werkzeug- oder Manifest-Ausfall
  bleibt `SKIP` je Teil — ein Sensor, der bei Netzstörung rot wird, wird
  abgeschaltet, und ein abgeschalteter Wächter ist schlechter als ein löchriger.

## 4. Definition of Done

- [x] `make baseline-verify` ist **erstes Glied** von `make gates`; der
      `--verify`-Pfad ruft kein Netz-Werkzeug, der innere Lauf bleibt offline.
- [x] `make baseline-freshness` läuft **außerhalb** `gates` und im eigenen
      [`upstream-drift.yml`](../../../../.github/workflows/upstream-drift.yml),
      täglich 01:00 UTC, getrennt von `ci.yml`.
- [x] **Vier** Verstoß-Formen rot gesehen, nicht zwei: geänderte · im Baum
      eingelegte · **als Geschwister der Bäume eingelegte** · gelöschte Datei —
      alle Exit 2, Rückbau Exit 0. Die dritte hat erst der Review gefunden
      (§9), und sie ging vorher **grün** durch.
- [x] `AGENTS.md` §4 und die Sensors-Tabelle tragen beide Targets;
      `make gate-consistency` grün.
- [x] `make gates` Exit 0 (neun Glieder, 471 Dateien, 0 Befunde),
      `make fullbuild` Exit 0; unabhängiger Review
      ([Report](../../../reviews/2026-08-23-slice-133-baseline-sensor-review.md)),
      blockierend, alle sechs Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Ein Nachtlauf, in den niemand sieht, ist ein zweiter verwaister Sensor.**
  — **Ausgang:** *nicht eingetreten.* GitHub schickt bei einem fehlgeschlagenen
  Workflow-Lauf eine **Mail**; der Job fällt rot aus statt nur zu protokollieren,
  und damit erreicht die Meldung einen Menschen, ohne dass jemand in die
  Actions-Übersicht sehen müsste. Der Zustellweg ist eine Konto-Einstellung des
  Empfängers, kein Repo-seitiger Mechanismus — für diesen Betrieb reicht das.
- **`gates` um ein Host-Skript zu erweitern berührt §3.1.** — **Ausgang:**
  *eingetreten, aber anders als vermutet.* Der Satz *„Der Host braucht **nur**
  `git`, GNU `make`, `bash` und Docker"* war **schon vorher falsch**:
  `working-tree-hash.sh`, `coverage-gate.sh` und `semgrep.sh` hängen längst in
  `gates` und rufen `awk`, `sha256sum`, `sort`, `grep`, `tr`. Neu kommt allein
  `find` dazu. §3.1 nennt jetzt die **Klasse** statt einer Liste — eine
  Aufzählung wäre der nächste ungewachte Spiegel.
- **Der Netz-Teil ist fail-open, der Gate-Teil fail-closed.** — **Ausgang:**
  *eingetreten, und die Trennung war anfangs nur behauptet.* Die Namen und die
  Doku-Zeilen trennen sauber, aber der fail-open-Pfad hatte **keine
  Zeitgrenze**: ein Fehlschlag wurde zu `SKIP`, eine **hängende** Verbindung
  wäre erst von der Job-Decke abgeräumt worden — und hätte den Nachtlauf rot
  gefärbt, also genau das Gegenteil der Zusage. Gefunden hat es der Review.

## 6. Trigger

**Start** (`open` → `in-progress`): [welle-84](../welle-84-durchsetzung.md)
eröffnet, WIP-Limit frei. Hängt **nicht** an
[slice-132](../in-progress/slice-132-hard-rule-zensus.md).

**Rückführungen:** `in-progress` → `next`, falls das Einhängen in `gates` die
Netzlos- oder Werkzeug-Zusage berührt — dann ist es ein Spec-Thema und kein
Verdrahtungs-Slice.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Werkzeuge (GF), Gate-Landschaft (GF), CI (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-007`](../observations.md) für jeden Beleg-Lauf — der Exit gehört
  gelesen, nicht hinter eine Pipe. [`BEO-010`](../observations.md) für die
  Doku-Spiegel der Gate-Liste: ein neues Target erscheint in **mehreren**
  Tabellen, und `gate-consistency` prüft Namen, nicht Mengen.
  [`BEO-011`](../observations.md) für jede Aussage darüber, was der Sensor
  abdeckt.

Slice-ID: slice-133. Betroffene IDs:
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(Netzlos-Zusage bleibt unberührt — zu belegen, nicht zu ändern). Module:
Harness-Werkzeuge, Gate-Landschaft, CI. Gates: `make gate-consistency`,
`make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Verdrahtung vorhandener Mechanik.

## 9. Closure-Notiz (nach `done/`)

Geliefert: `make baseline-verify` ist erstes Glied von `make gates`,
`make baseline-freshness` läuft im eigenen Nachtlauf, und beide stehen in den
zwei Doku-Tabellen, die `gate-consistency` gegen das `Makefile` hält. Gebaut
wurde **nichts** — das Skript konnte all das schon.

**Der Slice hat gefunden, wonach er gar nicht gesucht hat.** Er trat an, um
einen verwaisten Sensor einzustecken. Der Review hat gezeigt, dass dieser Sensor
**selbst einen stillen Grün-Pfad trug** — und zwar genau eine Verzeichnisebene
über der Lücke, die er schließt. Die Manifest-Deckung zählte per `find` nur
innerhalb von `regelwerk/` und `templates/`; eine Datei, die als **Geschwister**
der beiden Bäume abgelegt wird, liegt in keinem Baum und in keiner
Manifest-Zeile — beide Zahlen blieben gleich, und der Gate meldete grün. Meine
eigene Begründung stand eine Zeile darüber: *„ohne die Manifest-Deckung wäre
‚prüft die Integrität' überdehnt."* Sie war es noch immer.

**Das ist die Lehre, und sie ist unbequemer als die Reparatur.** Ein Sensor
einzustecken heißt nicht, ihn geprüft zu haben. Seit dem 25. Juni 2026 stand
dieses Skript im Repo und wurde von nichts gerufen; die erste Messung, die seine
Zusage wirklich befragt hat, war ein **konstruierter Verstoß** — nicht einer der
Läufe, die es nie gab.

**Zum vierten Mal an einem Tag eine Quelle überdehnt** ([`BEO-012`](../observations.md)).
Ich habe [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
zitiert, um zu begründen, dass `make gates` netzlos bleibt. Jene Anforderung
sagt etwas über das **Produkt** — *„Das Tool schreibt nie … und öffnet keine
Netzwerkverbindungen"* —, und ihre Messmethode ist ein Container-Lauf. Über den
inneren Lauf sagt sie nichts. Dass `gates` offline bleibt, ist eine Eigenschaft
**dieses Repos**, und sie steht jetzt ohne geliehene Autorität da.

**Und eine Zuschreibung, die aus einer zu weiten Messung kam.** Die
Werkzeug-Liste in `AGENTS.md` §3.1 hatte ich über **vier** Skripte gemessen und
dann **zweien** zugeschrieben; `record-gates.sh` ruft keines davon, `sed` ruft in
`gates` niemand, `tr` fehlte. Statt die Liste zu korrigieren nennt §3.1 jetzt die
Klasse — eine Aufzählung wäre der nächste Spiegel ohne Wächter
([`BEO-010`](../observations.md)).

**Offen bleibt nichts, das dieser Slice zu tragen hätte.** Der Workflow-Kopf
sagt hin, wo sein Rot erscheint, und der fehlgeschlagene Lauf meldet sich per
Mail — der Sensor erreicht einen Menschen, nicht nur ein Protokoll.
