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
   netzlos und verletzt [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) nicht — die Zusage *„`make gates` bleibt
   offline"* bleibt wörtlich erhalten.
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

- [ ] `--verify` läuft als benanntes Target **in** `make gates`; `make gates`
      bleibt netzlos (Exit explizit).
- [ ] `--check-latest` läuft als benanntes Target **außerhalb** `gates` und in
      einem **eigenen** geplanten Workflow, getrennt von `ci.yml`.
- [ ] **Beide** Verstoß-Richtungen sind rot gesehen — geänderte Datei **und**
      zusätzlich eingelegte Datei —, der Rückbau grün.
- [ ] `AGENTS.md` §4 und die Sensors-Tabelle tragen die neuen Targets;
      `make gate-consistency` grün.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Nachtlauf, in den niemand sieht, ist ein zweiter verwaister Sensor** —
  dieselbe Klasse, die dieser Slice behebt, eine Ebene höher. Der Slice muss
  benennen, **wo** sein Rot sichtbar wird. — **Ausgang:** *(bei Closure)*
- **`gates` um ein Host-Skript zu erweitern berührt §3.1.** Erlaubt sind `git`,
  GNU `make`, `bash` und Docker; `sha256sum`/`find` kommen dazu. Zu prüfen, nicht
  anzunehmen: ob der Gate-Lauf damit noch auf einem frischen Klon durchläuft. —
  **Ausgang:** *(bei Closure)*
- **Der Netz-Teil ist fail-open, der Gate-Teil fail-closed.** Zwei
  Fehlerpolitiken in einem Skript sind eine Verwechslungsquelle; die Targets
  müssen sie im Namen und in der Doku-Zeile trennen. — **Ausgang:** *(bei
  Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [welle-84](../welle-84-durchsetzung.md)
eröffnet, WIP-Limit frei. Hängt **nicht** an
[slice-132](../open/slice-132-hard-rule-zensus.md).

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

*(wird mit dem Closure-Body gefüllt)*
