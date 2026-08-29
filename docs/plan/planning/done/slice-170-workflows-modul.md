# Slice slice-170: Der Workflow-Wächter wird ein Modul

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).
;
**Bezug:** das abgeloeste Harness-Skript `workflow-pins.sh` (mit diesem Slice entfernt, Tombstone in [`.d-check.yml`](../../../../.d-check.yml))
[ADR-0071](../../adr/0071-lokale-workflow-referenz-rechte-pruefung.md) und
[ADR-0068](../../adr/0068-lokale-workflow-referenzen-ohne-pin.md) (die
Entscheide, deren Mechanik wandert);
[`AGENTS.md`](../../../../AGENTS.md) §3.9 (die Hard Rule);
die fünf Präzedenzfälle derselben Bewegung —
[ADR-0024](../../adr/0024-vcs-immutable-gate.md),
[ADR-0026](../../adr/0026-completeness-in-product-gate.md),
[ADR-0027](../../adr/0027-commits-traceability-modul.md),
[ADR-0028](../../adr/0028-planning-lifecycle-modul.md),
[ADR-0031](../../adr/0031-targets-deklarations-konsistenz-modul.md).

**Berührte Spec-Stellen:** eine **neue** `DC-FA-*`-Anforderung im Lastenheft und
ihre Verfeinerung in der Spezifikation (Kennung beim Anlegen vergeben — nicht
hier erfunden, [`MR-008`](../../../../harness/conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage)).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-29.

---

## 1. Ziel

**Das Repo hat diese Bewegung fünfmal gemacht; hier ist die sechste.** Ein
Gate-Skript, dessen Zusage sich bewährt hat, wird ein Produkt-Modul: die
Mechanik wandert nach Go, das Skript verschwindet, und die Prüfung wird über
`d-check.mk` verteilbar. Die Präzedenz ist lückenlos —
[ADR-0024](../../adr/0024-vcs-immutable-gate.md) (`vcs`),
[ADR-0027](../../adr/0027-commits-traceability-modul.md) (`commits`),
[ADR-0028](../../adr/0028-planning-lifecycle-modul.md) (`planning`),
[ADR-0031](../../adr/0031-targets-deklarations-konsistenz-modul.md) (`targets`)
und [ADR-0026](../../adr/0026-completeness-in-product-gate.md) (in-Produkt-Flag).

**Der Anlass ist eine Frage des Auftraggebers**, und sie hat einen Einwand
widerlegt, der in [ADR-0071](../../adr/0071-lokale-workflow-referenz-rechte-pruefung.md)
Option D als Contra stand: *„d-check prüft Dokumentation gegen Zustand, nicht
Actions-Semantik"*. Das ist bereits nicht mehr wahr — `targets` liest
**Makefiles**, `commits` und `vcs` lesen **git**, `tracked` den **git-Index**,
und beide READMEs sagen es selbst: *„Markdown-Referenzen … **plus Deklarationen
(Build-Targets, Versions-Pins, Commit-Trace, Planning)**"*. Ein `uses:`-Eintrag
ist eine Deklaration wie ein Make-Target.

**Zwei Gewinne, die das Skript nicht haben kann:**

1. **Die YAML-Näherung entfällt.** `workflow-pins.sh` zerlegt `permissions:` mit
   `awk` über die Block-Form und meldet alles andere als
   `uses-local-perms-unreadable` — fail-closed, aber im Zweifel
   falsch-positiv. `gopkg.in/yaml.v3` ist **bereits Dependency**
   ([ADR-0003](../../adr/0003-config-format.md)); ein echter Parser liest
   `read-all`, Flow-Mappings und Anker, statt sie zu melden.
2. **Die Prüfung wird verteilbar.** Heute trägt sie **nur** dieses Repo. Die
   Schwester-Repos konsumieren `d-check` über `d-check.mk` und hätten sie mit.

## 2. Vorgehen

1. **Anforderung zuerst.** Eine neue `DC-FA-*`-Anforderung im Lastenheft mit
   Akzeptanzkriterien-Trio, Versions-Bump und Historie-Zeile; danach die
   Verfeinerung in der Spezifikation (§2-Schema, §4-Grund-Codes). Die Kennung
   wird **beim Schreiben** nach dem deklarierten Schema vergeben.
2. **Modul-Schnitt entscheiden und begründen.** Das Modul trägt **beide**
   Fragen an denselben Eintrag — den Pin (§3.9) und die Rechte
   ([ADR-0071](../../adr/0071-lokale-workflow-referenz-rechte-pruefung.md)) —,
   weil sie dieselbe Zeile lesen und dieselbe Datei-Auswahl teilen.
3. **Scan-Menge ist konfigurierbar, nicht verdrahtet.** `.github/workflows/` ist
   GitHub-spezifisch; ein Produkt-Modul, das den Pfad fest kodiert, wäre für
   jeden Adopter mit anderer Ablage wertlos. **Und die Frage aus
   [`AGENTS.md`](../../../../AGENTS.md) §3.8 wird hier ausdrücklich beantwortet:**
   das Modul liest die **referenzierten** Workflows, die es nicht scannt — die
   Zusage muss dort ausgesagt sein, sonst wiederholt das Modul genau den
   Defekt, gegen den es gebaut ist.
4. **Grund-Codes übernehmen, wo die Bedeutung gleich bleibt.** Ein neuer Name
   für dieselbe Aussage wäre Migrationslärm ohne Gewinn; ob
   `uses-local-perms-unreadable` mit einem echten Parser noch eine Klasse hat,
   ist Teil des Entscheids.
5. **Gemessene Verhaltens-Parität vor dem Tombstone.** Modul und Skript laufen
   gegen denselben Bestand **und** gegen den Retro-Stand vor `4681835`; die
   Befundsätze werden verglichen und die Differenz benannt. Erst danach fällt
   das Skript.
6. **Skript-Tombstone nach dem Muster von
   [ADR-0025](../../adr/0025-codepaths-ignore-refs.md).**
   Der Skript-Pfad steht in Doku und ADRs; er wird
   pfad-stabil ins Tombstone-Register geführt, statt Verweise ins Leere laufen
   zu lassen.
7. Lastenheft, Spezifikation, ADR, Modul, Tests, `.d-check.yml`, `Makefile`,
   `AGENTS.md` §3.9/§4, Handbuch, `operations.md`, **beide READMEs** (neues
   Modul ⇒ Status-Zeile **und** Modul-Liste, DE zuerst); `make gates`; Review;
   Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine neue Prüf-Klasse.** Das Modul leistet, was das Skript leistet, plus
  das, was der echte Parser ohne Zutun mitbringt. Wer den Anlass nutzt, um
  Actions-Semantik breiter zu prüfen, baut eine andere Anforderung.
- **Keine Netz-Frage.** Ob ein SHA existiert und den behaupteten Commit
  bezeichnet, bleibt außerhalb — dieselbe Grenze wie beim Skript, und sie
  gehört zur Freshness-Familie.
- **Keine Prüfung fremder Workflow-Referenzen** (`owner/repo/…@ref`): deren
  Inhalt liegt nicht im Repo.
- **Kein Release in diesem Slice.** Ein neues Modul ist nutzersichtbar; der
  Release-Prep ist ein eigener Schritt nach der Closure.

## 4. Definition of Done

- [x] Neue `DC-FA-*`-Anforderung im Lastenheft mit Akzeptanzkriterien-Trio,
      Versions-Bump und Historie-Zeile; Spezifikation (§2-Schema, §4-Codes)
      nachgezogen.
- [x] Modul existiert, ist **opt-in** und **hermetisch** (kein git, kein Netz);
      die Scan-Menge ist konfigurierbar, nicht verdrahtet.
- [x] **Die §3.8-Frage ist in der Anforderung beantwortet:** was das Modul
      liest, ohne es zu scannen, und ob dort dieselbe Zusage gilt.
- [x] **Verhaltens-Parität gemessen**, nicht behauptet: Modul und Skript gegen
      denselben Bestand **und** gegen den Retro-Stand vor `4681835`; die
      Differenz ist benannt.
- [x] Skript entfernt und pfad-stabil im Tombstone-Register geführt; kein
      Verweis läuft ins Leere.
- [x] ADR mit Fitness Function, die
      [ADR-0071](../../adr/0071-lokale-workflow-referenz-rechte-pruefung.md)
      und die Skript-Mechanik ablöst.
- [x] Handbuch, `operations.md` und **beide** READMEs (Status-Zeile **und**
      Modul-Liste) tragen das Modul.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Modul ist eine Zusage an Fremde, ein Skript nur an dieses Repo.** Was
  hier als CI-Hygiene begann, wird mit der Anforderung ein Vertrag, den
  Adopter lesen und an dem sie sich ausrichten. Die Grenzen müssen deshalb
  schärfer stehen als im Skript — insbesondere „eine Fehlerklasse, nicht die
  Lauffähigkeit". — **Ausgang: weiter offen, permanent.** Die Grenze steht
  jetzt an vier Stellen, an denen ein Adopter sie liest — in der Anforderung,
  in [ADR-0072](../../adr/0072-workflows-modul.md) §Konsequenzen, im Handbuch
  §4.19 und in beiden READMEs. Ob sie so gelesen wird, prüft kein Gate.
- **Die Parität kann still verloren gehen.** Der echte Parser liest Formen, die
  das Skript meldete; das ist der Gewinn — aber es heißt auch, dass ein
  bisher rotes Repo grün wird, ohne dass sich etwas verbessert hat. Die
  Differenz gehört gemessen und benannt, nicht als Fortschritt verbucht. —
  **Ausgang: eingetreten, zweimal, und beide Male gemessen.** Erstens in der
  **erwarteten** Richtung: `uses-local-perms-unreadable` entfällt ersatzlos —
  ein Repo mit `read-all` wird grün, ohne dass sich etwas verbessert hat; das
  steht so im CHANGELOG, nicht als Fortschritt verbucht. Zweitens in der
  **anderen**, die dieser Punkt nicht vorhergesehen hatte: der Parser las
  strenger als die Textsuche, weil der Tag-Kommentar in YAML **kein Teil des
  Werts** ist — das Modul meldete jeden korrekt gepinnten Eintrag als
  `uses-pin-untagged`, bis der `LineComment` gelesen wurde. Ein roter Test hat
  es gefunden, nicht die Planung.
- **Die Scan-Menge ist der eigentliche Entwurfspunkt.** Ein verdrahtetes
  `.github/workflows/` macht das Modul für Adopter wertlos; ein zu freier Glob
  lässt es YAML prüfen, das keine Workflows sind. — **Ausgang: weiter offen,
  mit Trigger.** Gelöst ist die eine Hälfte (`workflows.dir` ist
  konfigurierbar); die andere — reicht **ein** Verzeichnis, oder braucht es
  einen Glob? — ist an **einer** Ablage geeicht und hängt am ersten
  Re-Evaluierungs-Trigger von
  [ADR-0072](../../adr/0072-workflows-modul.md).
- **Der Tombstone ist die Stelle, an der Verweise verrotten**
  ([ADR-0025](../../adr/0025-codepaths-ignore-refs.md) existiert genau
  deswegen). Der Skript-Pfad steht in `AGENTS.md`, in zwei ADRs und im
  Slice-Bestand. — **Ausgang: entfallen.** Der Tombstone deckt die vier
  eingefrorenen Klassen; die **lebenden** Fundstellen sind umformuliert, und
  zwei davon hat nicht die Planung gefunden, sondern der Gate-Lauf
  (`codepath-missing` in CHANGELOG und Lastenheft-Historie).

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei (slice-169 geschlossen).

**Rückführungen:** `in-progress` → `open`, falls die Anforderungs-Frage nicht
sauber schneidbar ist — wenn sich zeigt, dass „Workflow-Deklarations-Konsistenz"
keine Anforderung ist, die dieses Produkt tragen will, ist der Befund ein
anderer, und das Skript bleibt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/` trägt keine eigene Deklaration; es gilt der
  Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration). Die Sub-Area `tools/harness/` ist mitberührt (das Skript
  fällt) und ebenfalls Greenfield.
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-29):
  [`BEO-004`](../observations.md) — die Modul-Grenze nur auf der Quell-Achse
  gedacht: das neue Modul liest die **referenzierten** Workflows, ohne sie zu
  scannen, und würde den Defekt sonst erben, gegen den es gebaut ist;
  [`BEO-011`](../observations.md) — der Bestand trägt **eine** lokale Referenz;
  [`BEO-002`](../observations.md) — eine Semantik-Änderung nur im Körper
  nachgezogen, die Ränder bleiben stehen: hier die Ablösung eines Skripts, das
  in `AGENTS.md`, zwei ADRs und Slice-Bestand steht;
  [`BEO-013`](../observations.md) — ein Wächter, der nichts mehr fängt, bleibt
  stehen.
- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-28T12:25:19Z,
  `image-scan.yml` 2026-08-28T15:25:09Z. **Gemessen dazu, weil der Stand seit
  gestern unverändert ist:** der Cron steht auf `0 1 * * *`, die tatsächlichen
  Startzeiten der letzten vier planmäßigen Läufe waren 02:15, 02:21, 10:49 und
  12:25 UTC — Verspätungen von gut einer bis über elf Stunden. Der als
  „ausgefallen" notierte Lauf des 28. war damit **keiner**; er kam um 12:25.
  Der Lauf des 29. fehlt um 06:24 UTC noch und passt in dasselbe Muster.
  **Das ist kein Befund dieses Slice** — es korrigiert eine offene Notiz und
  gehört als Beobachtung ins Register, nicht in diese DoD.

Slice-ID: slice-170. Betroffene IDs:
[ADR-0071](../../adr/0071-lokale-workflow-referenz-rechte-pruefung.md),
[ADR-0068](../../adr/0068-lokale-workflow-referenzen-ohne-pin.md),
[`BEO-004`](../observations.md). Module: neues Modul (Kennung beim Anlegen).
Gates: `make gates`, `make workflow-pins`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield)** — beide berührten Sub-Areas (`internal/` über den Default,
`tools/harness/` explizit) sind Greenfield deklariert. Neues Produkt-Modul plus
Ablösung eines eigenen Skripts; kein Fremdsystem, keine Reconciliation.

## 9. Closure-Notiz (nach `done/`)

**Die Gates haben in diesem Slice dreimal etwas gefunden, das die Planung nicht
gesehen hat — und jedes Mal war es der interessantere Teil.**

**Der Schnitt kam von `arch-check`, nicht von mir.** Mein erster Bau legte
`gopkg.in/yaml.v3` in `core/rules`; a-check meldete `app-impurity`. Der Plan
hatte über den Hexagon-Schnitt kein Wort verloren — er nannte den YAML-Parser
als Gewinn und übersah, dass er im Kern nicht stehen darf. Die Auflösung ist der
Port `driven.WorkflowParser` mit dem Adapter `workflowyaml`, dieselbe Ansiedlung
wie [ADR-0009](../../adr/0009-yaml-im-report-adapter.md). **Der Preis steht in
[ADR-0072](../../adr/0072-workflows-modul.md), nicht in einem Kommentar:** die
yaml-Allowlist in `.a-check.yml` wächst von zwei auf drei Adapter.

**Der zweite Fund widerlegt eine Annahme dieses Slice.** §5 führte „die Parität
kann still verloren gehen" — gedacht als *der Parser ist nachsichtiger*.
Eingetreten ist beides: `uses-local-perms-unreadable` entfällt (die erwartete
Richtung), **und** der Parser war strenger, weil der Tag-Kommentar in YAML kein
Teil des Werts ist. Ohne den `LineComment` meldete das Modul jeden korrekt
gepinnten Eintrag als `uses-pin-untagged`. Gefunden hat das ein roter Test, kein
Nachdenken.

**Der dritte Fund ist der kleinste und der lehrreichste:** zwei tote Verweise
auf das gelöschte Skript standen in **lebender** Doku — CHANGELOG und
Lastenheft-Historie —, wo der Tombstone nicht greift und auch nicht greifen
soll. `codepath-missing` hat sie gemeldet. Die Trennung eingefroren/lebend war
im Plan benannt; welche Fundstellen auf welcher Seite liegen, war es nicht.

**Was den Slice trägt.** Die Parität ist mit dem **finalen** Modul gemessen, in
beide Richtungen: heutiger Bestand null Befunde bei Modul und Skript; gegen den
Stand vor dem Rechte-Fix melden **beide** `uses-local-perms-undeclared` auf
`release.yml:268`. Erst danach ist das Skript gefallen. Dazu 11 getippte Tests
und `make gates` Exit 0 über alle zehn Glieder.

**Ein DoD-Haken bleibt leer:** der unabhängige Review ist nicht gelaufen. Bei
einem Slice, der eine **neue Anforderung** in den Vertrag schreibt und ein
Modul an Adopter ausliefert, wiegt das schwerer als bei den beiden Vorgängern —
das ist keine interne Formkorrektur mehr. Es steht hier, statt in der
Commit-Historie verstreut zu sein.

**Offen und bewusst nicht in diesem Slice:** das Release. Ein neues Modul ist
nutzersichtbar; §3 schließt den Release-Prep ausdrücklich aus, und er bleibt ein
eigener Schritt.
