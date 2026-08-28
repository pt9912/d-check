# Slice slice-167: Der Sensor bekommt einen Kanal, der hebt statt zu melden

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [ADR-0066](../../adr/0066-cve-scan-gegen-das-publizierte-image.md)
(der Sensor, dessen Befunde niemand hebt);
[ADR-0011](../../adr/0011-digest-pins-build-gate-images.md) (die Digest-Pins,
die Dependabot **nicht** anfassen darf);
[ADR-0027](../../adr/0027-commits-traceability-modul.md) (das Traceability-Gate,
an dem jeder Dependabot-PR sonst scheitert).

**Berührte Spec-Stellen:** — (CI-Konfiguration; das Produkt bleibt unberührt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-28.

---

## 1. Ziel

**Die vierzehn Befunde von gestern hat ein Mensch gehoben. Der nächste Satz soll
sich selbst als PR melden.**

[ADR-0066](../../adr/0066-cve-scan-gegen-das-publizierte-image.md) hat den
Sensor gebaut; er **meldet**. Hebt niemand, meldet er dieselben Befunde morgen
wieder — und die Klasse *„Sensor ohne Adressat"* ist in diesem Repo schon
einmal geschnitten worden
([slice-164](../done/slice-164-nachtlauf-kadenz.md)).

## 2. Vorgehen

1. **Zuerst die Kollision, dann die Konfiguration.** `make trace-check` läuft in
   der PR-CI über **jede** Commit-Message und verlangt eine `DC-`/`ADR-`/`MR-`/
   `slice-`-Kennung. **Gemessen:** eine gewöhnliche Dependabot-Botschaft
   (*„build(ci): bump actions/checkout from 7.0.1 to 7.1.0"*) ist **rot**;
   dieselbe mit einer ADR-Kennung im Präfix ist **grün**. Ohne diese Frage
   vorweg wäre jeder Dependabot-PR ab dem ersten Tag rot.
2. **Die Kennung ist eine Zusage, keine Umgehung.** Der Ausweg wäre eine
   Erweiterung von `commits.exempt-pattern` — eine **gelockerte Prüfregel**
   ([`AGENTS.md`](../../../../AGENTS.md) §3.6). Stattdessen trägt
   `commit-message.prefix` die ADR-Kennung **dieses** Entscheids: die ADR, die
   Dependabot erlaubt, ist der ehrliche Anker für die Commits, die daraus
   entstehen.
3. **Die Ökosysteme aus dem Bestand wählen, nicht aus dem Anlass.** Zwei
   Schwester-Repos fahren Dependabot unterschiedlich — eines **nur**
   `github-actions`, eines zusätzlich `docker` mit ausführlichen `ignore`-Regeln.
   Für d-check entscheidet **unsere** Pin-Lage, nicht deren Beispiel.
4. **`docker` bleibt draußen, und das ist der wichtigste Ausschluss.** Die
   Basis-Images sind **digest**-gepinnt
   ([ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)) und tragen ihre
   Version in einer Make-Variablen; ein Dependabot-PR am Tag ließe den Digest
   stehen und bräche die Kopplung. Für genau diese Frage gibt es **fünf**
   Digest-Achsen im Nachtlauf.
5. **Gruppieren, sonst ist es Rauschen.** Ein PR je Action bei sieben externen
   `uses:`-Einträgen ist kein Kanal, sondern eine Flut.
6. ADR; `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Automatik am Merge.** Dependabot öffnet PRs; die Prüfung bleibt
  `make ci`, und der Merge bleibt ein Akt.
- **Kein `docker`-Ökosystem** (siehe §2.4).
- **Keine Abschaffung der Frische-Achsen.** Sie überschneiden sich mit
  Dependabot nur bei den **Action**-Pins; Go-Toolchain, `semgrep`, `a-check`,
  Trivy und die Baseline sind hand-gepinnt und bleiben es.
- **Keine Änderung an `commits.exempt-pattern`.** Der Gate bleibt, wie er ist.

## 4. Definition of Done

- [ ] `.github/dependabot.yml` existiert; jedes aufgenommene **und** jedes
      ausgelassene Ökosystem trägt seine Begründung in der Datei.
- [ ] Die Traceability-Kollision ist **gemessen** gelöst: eine simulierte
      Dependabot-Botschaft läuft gegen das echte `make trace-check` grün, eine
      ohne Kennung rot.
- [ ] `commits.exempt-pattern` ist **unverändert** — die Lösung lockert kein
      Gate.
- [ ] ADR mit Fitness Function; `AGENTS.md` §5 nennt die Dependabot-Commits als
      benannte Form, damit die Kennung nicht wie eine Ausnahme aussieht.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Eine Kennung im Präfix ist eine Behauptung über jeden künftigen Commit.**
  Sie sagt *„dieser Bump gehört zu [ADR-0067](../../adr/0067-dependabot-als-hebender-kanal.md)"* — und das stimmt für den Kanal,
  nicht notwendig für den Inhalt. Wer sie liest, könnte mehr Bezug annehmen, als
  da ist. — **Ausgang:** *(bei Closure)*
- **Ein `ignore`-Eintrag, der seinen Grund überlebt, ist ein stiller
  Grün-Pfad.** Er sieht aus wie eine Entscheidung und ist eine Gewohnheit; kein
  Gate sieht ihn ([`BEO-013`](../observations.md)). — **Ausgang: entfallen.**
  Die Konfiguration enthält **keinen** `ignore`-Schlüssel — das Risiko beschrieb
  die Datei des Schwester-Repos, nicht unsere. Die Klasse trifft stattdessen die
  **Ausschlüsse** (`docker` und die vier hand-gepinnten Bestände), und die
  tragen ihren Grund in der Datei.
- **Dependabot hebt, was es sieht — und es sieht `go.mod`, nicht das
  ausgelieferte Bild.** Die vierzehn Befunde lagen in **indirekten**
  Abhängigkeiten; ob Dependabot sie ohne Security-Advisory anfasst, ist eine
  offene Frage und keine Zusage. — **Ausgang: eingetreten, und die Frage war
  nicht offen, sondern ungemessen.** Gemessen: **dreizehn** von vierzehn lagen
  indirekt; Version-Updates sehen ohne `allow: dependency-type: all` **nur** die
  zwei direkten Requires; und der zweite Weg — Security-Updates — ist im
  Repository **abgeschaltet**. Der Kanal erreichte die Fundklasse in der ersten
  Fassung also **nicht**. Die Konfiguration trägt jetzt `allow`; die zweite
  Hälfte ist ein Repo-Setting und steht als Vorbedingung in
  [`releasing.md`](../../../../docs/user/releasing.md).

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `open`, falls sich zeigt, dass die
Traceability-Kollision **nur** über eine Gate-Lockerung lösbar ist — dann ist
der Entscheid ein anderer (§3.6) und gehört vor die Konfiguration.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** CI-/Dependency-Konfiguration (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-28):
  [`BEO-011`](../observations.md) — die Regel aus dem **Bestand**, nicht aus dem
  Anlass: zwei Schwester-Repos fahren Dependabot verschieden, und **unsere**
  Pin-Lage entscheidet, nicht ihr Beispiel;
  [`BEO-013`](../observations.md) — ein Unterdrückungs-Marker, der nichts mehr
  unterdrückt, bleibt stehen: genau die Bauform eines `ignore`-Eintrags;
  [`BEO-009`](../observations.md) — nicht mehr behaupten, als die Messung trägt.
- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): `upstream-drift.yml`
  **ROT**, jüngster Lauf weiterhin `2026-08-27T10:49:23Z` — **der Lauf des 28.
  ist ausgefallen**, nicht nur verspätet (05:22 UTC gegen Cron 01:00 bei
  historisch +75…+81 min, kein Eintrag in der Lauf-Liste). Die benannte Grenze
  aus [slice-166](../done/slice-166-cve-image-scan.md) §5 ist damit am ersten
  Tag eingetreten. `image-scan.yml` hat noch keinen geplanten Lauf gehabt.

Slice-ID: slice-167. Betroffene IDs: — (kein `DC-`Bezug; CI-Konfiguration).
Module: — . Gates: `make gates`, `make trace-check`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — neue Konfigurationsdatei neben den
vorhandenen CI-Artefakten; kein Fremdsystem, keine Reconciliation.

## 9. Closure-Notiz (nach `done/`)
