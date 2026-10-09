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

- [ ] Reiner Move-Commit (Git erkennt den Rename), danach ein Commit, der jeden
      Verweis nachzieht — die Liste am Repo gezählt, das Kommando im Plan;
      `make doc-check`, `make adr-check` über die Range und `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
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
  Adopter) zeigen auf den alten Pfad. — **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

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
