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

- [x] `.github/dependabot.yml` existiert; jedes aufgenommene **und** jedes
      ausgelassene Ökosystem trägt seine Begründung in der Datei. **Die
      Begründung des wichtigsten Ausschlusses war im ersten Anlauf falsch** und
      steht jetzt auf [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)
      §4 statt auf einem technischen Bruch, den es nicht gibt.
- [x] Die Traceability-Kollision ist **gemessen** gelöst — und der Review hat
      die Messung reproduziert **und erweitert**: auch eine gruppierte,
      mehrzeilige Botschaft und ein Dependabot-Merge-Commit laufen grün, und
      `trace-check` prüft die **ganze** Message, nicht nur den Betreff.
- [x] `commits.exempt-pattern` ist **unverändert** — nachgeprüft gegen
      `origin/main`.
- [x] ADR mit Fitness Function; [`AGENTS.md`](../../../../AGENTS.md) §5 nennt
      die Dependabot-Form. **Nachgebessert:** der §3.6-Verweis dort las die
      Regel als Verbot statt als Verfahrenspflicht.
- [x] `make gates` und `make ci` grün (Exit explizit); unabhängiger Review,
      Urteil *„schließbar nach Nacharbeit"*, neun Befunde eingearbeitet.
- [ ] **Nicht erfüllt und nicht erfüllbar durch diesen Slice:** die zweite
      Hälfte des Kanals ist ein **Repository-Schalter**, keine Datei.
      Dependabot-Alerts und Security-Updates sind aus (gemessen); ohne sie
      öffnet ein CVE ohne neues Upstream-Release keinen PR. Als Vorbedingung in
      [`releasing.md`](../../../../docs/user/releasing.md) benannt.

## 5. Abnahme-Punkte / Risiken

- **Eine Kennung im Präfix ist eine Behauptung über jeden künftigen Commit.**
  Sie sagt *„dieser Bump gehört zu [ADR-0067](../../adr/0067-dependabot-als-hebender-kanal.md)"* — und das stimmt für den Kanal,
  nicht notwendig für den Inhalt. Wer sie liest, könnte mehr Bezug annehmen, als
  da ist. — **Ausgang: weiter offen, als Beobachtung notiert.** Der Slice kann
  ihn nicht schließen: er tritt erst ein, wenn jemand einen Bump-Commit liest
  und aus der Kennung einen inhaltlichen Bezug folgert. **Was dagegen steht, ist
  geschrieben, nicht gebaut** — der Satz *„die Kennung gilt dem Kanal, nicht dem
  Inhalt des einzelnen Bumps"* steht in
  [`AGENTS.md`](../../../../AGENTS.md) §5 und in der ADR. Kein Gate fängt eine
  Fehl-Lesung; die Klasse bleibt Urteil.
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

**Die tragende Messung hielt. Alles, was ich um sie herum behauptet habe, hielt
zur Hälfte nicht.**

**Der Kern war richtig gebaut, und zwar in der richtigen Reihenfolge.**
`make trace-check` verlangt von jeder Commit-Message eine Kennung; Dependabots
Botschaften tragen keine. Die Frage stand **vor** der Konfiguration, wurde gegen
das echte Gate gemessen, und die Lösung lockert kein Gate: die Kennung steht im
Präfix, `commits.exempt-pattern` ist unverändert. Der Review hat die Messung
reproduziert **und erweitert** — auch eine gruppierte, mehrzeilige Botschaft und
ein Dependabot-Merge laufen grün.

**Der Kanal hätte die Fundklasse nicht erreicht, für die er gebaut ist.** Das
ist der schwerste Befund, und er trifft die Zusage im Dateikopf. Meine ADR
führte die Reichweite als *„nicht gemessen"* — sie **war** messbar: dreizehn der
vierzehn Befunde lagen `// indirect`, `go.mod` führt zwei direkte gegen zwanzig
indirekte Requires, Version-Updates sehen ohne `allow: dependency-type: all`
**nur** die direkten, und der zweite Weg ist im Repository abgeschaltet
(`automated-security-fixes` `enabled:false`, `vulnerability-alerts` HTTP 404).
Beide Wege zu. Der Kopf versprach *„einen Weg, der nicht daran hängt, dass
jemand hinsieht"* — und hätte zwei direkte Requires bewegt.

**„Nicht gemessen" war die bequeme Formulierung, nicht die ehrliche.** Sie
klingt nach Sorgfalt und ist hier das Gegenteil: die Frage stand offen, weil ich
sie nicht gestellt habe, nicht weil sie sich nicht stellen ließ. Zwei
API-Aufrufe und ein Blick in `go.mod` hätten gereicht.

**Die Begründung des wichtigsten Ausschlusses war gegen die Quelle falsch —
beide Sätze.** Ich schrieb, ein Dependabot-PR höbe den Tag und ließe den Digest
stehen, und Dependabot stelle die Frage *„anderer Bau desselben Tags"* gar
nicht. Der Review hat `dependabot-core` gelesen: der Docker-Updater hebt Tag
**und** Digest gemeinsam und stellt die Digest-Frage ausdrücklich. **Der
Ausschluss steht trotzdem** — aber auf
[ADR-0011](../../adr/0011-digest-pins-build-gate-images.md) §4, einer **Policy**
statt eines Bruchs, den es nicht gibt. Dazu eine dritte Tatsache, die ich nicht
kannte: `FROM golang:${GO_VERSION}@sha256:…` parst Dependabot ohnehin nicht.

**Die Labels wären wirkungslos gewesen.** `dependencies` und `ci` existieren im
Repo nicht; gesetzte Labels **ersetzen** die Defaults, und unbekannte werden
still ignoriert — die PRs wären ganz ohne Label gekommen. Genau die Bauform,
vor der das Schwester-Repo warnt und die
[`BEO-013`](../observations.md) führt — die ich in §7 **als gesichtet
deklariert** hatte.

**Denselben Fehler zum zweiten Mal: die Index-Zeile am Dateiende.** `>>` hängt
an, statt einzufügen; bei [ADR-0066](../../adr/0066-cve-scan-gegen-das-publizierte-image.md) ebenso. Zweimal heißt mechanisieren statt
aufpassen — eine `structure`-Regel verbietet jetzt ADR-Tabellenzeilen unter
`## Konventionen`, in beide Richtungen gegengeprobt. **Und die erste Fassung der
Regel feuerte nicht:** mein `^`-Anker zielte ins Leere, weil `forbid-pattern`
den Abschnittstext als Ganzes prüft. Eine Regel, die man nicht gegenprobt, ist
eine Hoffnung — das steht als zweite Grenze im Kommentar.

**Was der Review nicht brechen konnte, ist bemerkenswerter als das, was er
brach.** Er hat den Präfix-Mechanismus in `dependabot-core` **nachgeschlagen**
(die schließende eckige Klammer steht dort ausdrücklich in der Zeichenklasse),
das 50-Zeichen-Limit gegen das SchemaStore-Schema geprüft und empirisch am
Schwester-Repo belegt, dass ein Dependabot-Action-Bump SHA **und**
Tag-Kommentar gemeinsam ersetzt — `make workflow-pins` ist also nicht
gefährdet. **Alle drei hatte ich behauptet, ohne sie zu belegen.** Sie stimmten;
das war Glück, nicht Methode.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 576 Dateien, 0 Befunde),
`make ci` (Exit 0, image-test 4/4), `make trace-check` gegen vier simulierte
Botschaften in beide Richtungen mit gelesener Ausgabe, die `structure`-Regel
gegen eine konstruierte Fehlplatzierung, und `git diff origin/main --
.d-check.yml` als Beleg, dass `exempt-pattern` unberührt ist.

**Was offen bleibt, ist kein Rest, sondern die halbe Zusage:** Dependabot-Alerts
und Security-Updates sind Repository-Schalter, keine Dateien. Bis sie an sind,
hebt der Kanal, was **veraltet** ist — nicht, was **verwundbar** ist. Der
Unterschied steht in
[`releasing.md`](../../../../docs/user/releasing.md) §Vorbedingungen.
