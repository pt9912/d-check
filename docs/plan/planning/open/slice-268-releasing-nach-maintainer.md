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

**Verantwortlich:** —

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Doku unter dem Default `*` (`ALL`);
beim Beanspruchen neu prüfen.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** beim Beanspruchen neu lesen.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
beim Beanspruchen aus dem jüngsten Lauf lesen.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
