# slice-268: `releasing.md` zieht nach docs/user/maintainer/

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** Auftraggeber-Wunsch 2026-10-09; Baseline `v6.17.0` ·
`regelwerk/modul-04-adrs.md` §Nachzug ist keine Überschreibung; braucht
slice-267 (Pfad-Nachzug in `Accepted`-ADRs).

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** `docs/user/releasing.md` liegt unter docs/user/maintainer/releasing.md
— ein reiner `git mv`, danach ein eigener Commit, der die eigenen relativen
Links der Datei und alle eingehenden Verweise nachzieht, auch die in
`Accepted`-ADRs und in eingefrorenen Dokumenten (Pfad-Nachzug). Gemessen beim
Schnitt: 12 lebende Stellen, dazu 9 Links aus eingefrorenen Dateien (zwei
ADRs, drei Wellen-Ergebnisnotizen, ein alter CHANGELOG-Eintrag) und weitere
Erwähnungen ohne Link.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Weitere Dateien nach `maintainer/`** — nur `releasing.md` ist gewünscht;
  was sonst Maintainer-Doku ist, wäre ein eigener Schnitt.
- **Inhaltliche Änderungen an `releasing.md`** — der Umzug ändert nur Pfade.
- **Die Prüfung des Pfad-Nachzugs** — slice-267.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] Reiner Move-Commit (Git erkennt den Rename), danach ein Commit, der jeden
      Verweis nachzieht — die Liste am Repo gezählt, das Kommando im Plan;
      `make doc-check`, `make adr-check` über die Range und `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [x] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/releasing.md` → docs/user/maintainer/releasing.md | move | Umzug |
| alle Dateien mit Verweis auf `releasing.md` | update | Pfad-Nachzug |

*(Plan-Änderung vor dem Code-Commit: Die Zählung beim Schnitt las nur Links.
Gemessen im Wegwerf-Klon mit dem reinen Umzug und `make doc-check`
(`docker run … d-check:latest` über den Klon, Befunde nach Grund und Datei
gezählt): 39 `target-missing` — 23 eigene Links der Datei, 16 eingehende
(CHANGELOG, README und README.de je zwei, Benutzerhandbuch zwei,
Docker-Hub-Beschreibung, [ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md) vier, [ADR-0067](../../adr/0067-dependabot-als-hebender-kanal.md) einer in `## Geschichte`,
drei Wellen-Ergebnisnotizen) — und 19 `codepath-missing`: der Pfad als
Inline-Code in zwei `Accepted`-ADRs, in geschlossenen Slices, Register-Belegen,
einem aufgelösten MR, einem alten CHANGELOG-Eintrag und in diesem Plan. Die
Links werden nachgezogen. Die Inline-Code-Erwähnungen sind historischer Text
— in den ADRs dürfen sie sich nicht ändern, Inline-Code ist kein Link —; der
alte Pfad kommt deshalb ins Tombstone-Register `codepaths.ignore-refs`
([ADR-0025](../../adr/0025-codepaths-ignore-refs.md)). Dazu die Erwähnungen,
die kein Doku-Gate liest: `.github/dependabot.yml`, zwei in `release.yml`,
eine in `hub-description.yml`, `tools/image-test.sh`; und der
`matrix.exempt-paths`-Eintrag der `.d-check.yml`.)*

## 4. Trigger

**Start** (`next` → `in-progress`): slice-267 in `done/`; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `open` (blockiert): `adr-check` lässt den Nachzug in einer
  ADR trotz slice-267 nicht durch — dann zurück zu slice-267.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- **Verweise außerhalb des Repos** — Links von außen (Docker-Hub-Beschreibung,
  Adopter) zeigen auf den alten Pfad. — **Ausgang:** *weiter offen* — die
  Docker-Hub-Seite (`packaging/dockerhub/overview.md`) verlinkt die Datei
  nicht; Links von Adoptern sind von hier aus nicht messbar, die Release-Notiz
  nennt den Umzug;
  [`BEO-ALL/externer-link-auf-umgezogene-datei`](../observations/BEO-ALL/externer-link-auf-umgezogene-datei/state.md)
  (1×).

## 7. Closure-Notiz

- **Was hat funktioniert:** `releasing.md` liegt unter
  `docs/user/maintainer/`; Git erkennt den Rename (R100). Jeder Link folgt der
  Datei — 23 eigene, 16 eingehende, in den `Accepted`-ADRs 0014 und 0067 nur
  das Ziel, ohne Folge-ADR (Pfad-Nachzug nach
  [ADR-0103](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md)); die
  Verifikation löste alle 39 alten und neuen Ziele auf dieselbe Datei auf. Die
  19 Inline-Code-Nennungen des alten Pfads bleiben als historischer Text im
  Tombstone-Register von `codepaths`
  ([ADR-0025](../../adr/0025-codepaths-ignore-refs.md)). Erster echter Lauf des
  Pfad-Nachzugs aus slice-267, im Hook und in der CI grün.
- **Was ging anders als geplant:** Die Zählung beim Schnitt las nur Links;
  der Umzug im Wegwerf-Klon zeigte die 19 Inline-Code-Nennungen dazu — die
  Plan-Änderung stand vor dem Code. Anders als in §1 und der DoD beschrieben
  reisen die eingehenden Verweise im Move-Commit mit: getrennt gestagt, prüfte
  der gestagte ADR-Check den Nachzug in [ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md) gegen einen BASE-Stand ohne den
  alten Pfad und meldete Drift. Der Schnitt steht nur in der Commit-Botschaft,
  nicht als Plan-Änderung, und `AGENTS.md` §3.3 erlaubt das Mitreisen dem
  Wortlaut nach nur beim Übergang nach `done/` (R1 LOW-1, Verifikation V-1).
  Der Tombstone-Eintrag wirkt repo-weit; seine Grenze steht jetzt im Kommentar
  (R1 LOW-2).
- **Steering-Loop-Eintrag:** keiner mit Schwelle. Die Lücke zwischen §3.3 und
  dem ADR-Gate ist neu im Register — ändert sie §3.3, ist das eine Frage an
  den Auftraggeber, nicht an diesen Slice.
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-268.md` in
  [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/state.md);
  neu
  [`BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise`](../observations/BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise/state.md)
  und
  [`BEO-ALL/externer-link-auf-umgezogene-datei`](../observations/BEO-ALL/externer-link-auf-umgezogene-datei/state.md)
  (je 1×).
- **Folge-Slices:** keiner. Release v0.85.0 mit slice-263, slice-265,
  slice-267 und diesem Slice; die Release-Notiz nennt den Umzug.
- **Risiken aus §6:** eines weiter offen (Register, siehe §6). Trigger-Audit:
  kein Carveout, kein bootstrap-aware Gate, keine ADR, keine Hard Rule mit
  eingetretenem Trigger. Nachtlauf-Stand
  ([`MR-053`](../../../../harness/conventions.md#mr-053)): wie in §8 —
  `image-scan` rot bis zum nächsten Release.
- **Drei Paarungen:** (a) Anker — kein Eintrag mit `liegt in`; (b) Folge-Slices
  — keiner genannt; (c) Register — die drei zitierten Beobachtungen existieren
  und tragen Belege.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** geändert werden Betriebs-Doku,
README, Handbuch, Planungs- und Register-Dateien, zwei ADRs (nur Link-Ziele),
Workflows, ein Gate-Skript und die `.d-check.yml` — alle unter dem Default
`*` (`ALL`); deklariert.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** gelesen am 2026-10-09.
[`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`](../observations/BEO-ALL/pfad-nachzug-gleicher-name-andere-datei/state.md)
(1×) — jeder Nachzug in einer `Accepted`-ADR zeigt auf dieselbe umgezogene
Datei, im Review gegen das alte Ziel lesen;
[`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/state.md)
([`MR-070`](../../../../harness/conventions.md#mr-070)) — die Frozen-Klassen
sind vor der Ersetzung über ihre Eigenschaft aufgelistet (Plan-Änderung in §3);
[`BEO-ALL/path-scoped-commit-carries-staged-rest`](../observations/BEO-ALL/path-scoped-commit-carries-staged-rest/state.md)
(verkörpert als Schritt 21) — Move-Commit und Nachzug-Commit getrennt stagen.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-09 aus dem jüngsten Lauf (`make nightly-state`) —
`upstream-drift` grün (2026-10-09 07:17 UTC); `image-scan` rot (2026-10-09
10:37 UTC): zwei behebbare HIGH in der Standardbibliothek des publizierten
Images (Go 1.27.1), behoben ab Go 1.27.2 — der Pin ist auf `main` gehoben,
das nächste Release liefert ihn aus.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
