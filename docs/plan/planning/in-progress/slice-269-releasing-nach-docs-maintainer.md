# slice-269: `releasing.md` zieht nach docs/maintainer/

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** Auftraggeber-Korrektur 2026-10-09: das Ziel ist
docs/maintainer/releasing.md, nicht docs/user/maintainer/releasing.md, wohin
slice-268 die Datei verschoben hat; Baseline `v6.17.0` ·
`regelwerk/modul-04-adrs.md` §Nachzug ist keine Überschreibung.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die Datei liegt unter docs/maintainer/releasing.md. Ein Commit
verschiebt sie und nimmt die eingehenden Verweise mit, die bewegte Datei darin
unverändert (Git erkennt den Rename; in `Accepted`-ADRs nur Link-Ziele,
Pfad-Nachzug); ein zweiter zieht ihre eigenen Links nach. Das leere
Verzeichnis docs/user/maintainer/ entfällt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Weitere Dateien nach docs/maintainer/** — nur `releasing.md` ist
  gewünscht; was sonst Maintainer-Doku ist, wäre ein eigener Schnitt.
- **Inhaltliche Änderungen an `releasing.md`** — der Umzug ändert nur Pfade.
- **Die Nennungen des Zwischen-Pfads in slice-268 und dessen Reports** — sie
  beschreiben, was slice-268 geliefert hat, und bleiben stehen.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Ein Move-Commit, den Git als Rename erkennt, mit den eingehenden
      Verweisen; danach der Commit für die eigenen Links — die Liste am Repo
      gezählt, das Kommando im Plan; `make doc-check`, `make adr-check` über
      die Range und `make gates` grün, auch in einem frischen Klon.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| docs/user/maintainer/releasing.md → docs/maintainer/releasing.md | move | Umzug |
| alle Dateien mit Verweis auf die Datei | update | Pfad-Nachzug |
| `.d-check.yml` | update | `matrix.exempt-paths`, Tombstone |

Gemessen im Wegwerf-Klon mit dem reinen Umzug, dem gelöschten leeren
Verzeichnis und `d-check` über den Klon (Befunde nach Grund und Datei
gezählt): 40 Befunde. 23 eigene Links (12 `target-missing`, 11 `repo-escape`),
16 eingehende `target-missing` in denselben Dateien wie bei slice-268
(CHANGELOG einer, README und README.de je zwei, Benutzerhandbuch zwei,
Docker-Hub-Beschreibung,
[ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md) vier,
[ADR-0067](../../adr/0067-dependabot-als-hebender-kanal.md) einer, drei
Wellen-Ergebnisnotizen) und ein `codepath-missing`: die Closure-Notiz von
slice-268 nennt das Zwischen-Verzeichnis in Inline-Code — es kommt ins
Tombstone-Register ([ADR-0025](../../adr/0025-codepaths-ignore-refs.md)).
Ohne Doku-Gate: `.github/dependabot.yml`, `release.yml` (zwei),
`hub-description.yml`, `tools/image-test.sh`, der Kommentar zum
Tombstone-Eintrag des ursprünglichen Pfads.

*(Plan-Änderung vor dem nächsten Code-Commit, nach R1 MEDIUM-1: Mit dem
Umzug verlässt `releasing.md` den Rang 6 der Source Precedence, den
`AGENTS.md` §2 und `harness/README.md` §Source precedence `docs/user/`
zuordnen — „Operations, Releasing". Der Baseline-Default legt die
Releasing-Sicht unter `docs/user/*.md`; der neue Ort ist deshalb eine
Abweichung und bekommt einen `MR`-Eintrag in `harness/conventions.md`; Rang 6
nennt in beiden Tabellen beide Verzeichnisse. Gemessen: keine Scan-Menge der
`.d-check.yml` hängt an `docs/user/` außer einer `structure`-Regel für das
Benutzerhandbuch, und `make mention-coverage` deckt alle Artefakte.)*

## 4. Trigger

**Start** (`next` → `in-progress`): slice-268 in `done/`; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `open` (blockiert): `adr-check` lässt den zweiten Nachzug in
  einer ADR nicht durch.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- **Ein frischer Klon sieht das leere Verzeichnis nicht** — lokal bleibt es
  nach dem Umzug stehen, die CI kennt es nicht; eine Prüfung im lokalen
  Arbeitsbaum allein verdeckt das. — **Ausgang:** *(offen)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** dieselben Dateiklassen wie slice-268
— Betriebs-Doku, README, Handbuch, Planungs-Dateien, zwei ADRs (nur
Link-Ziele), Workflows, ein Gate-Skript, `.d-check.yml` — alle unter dem
Default `*` (`ALL`); deklariert.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** gelesen am 2026-10-09.
[`BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise`](../observations/BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise/state.md)
(1×) — der Schnitt dieses Slice folgt ihr und steht deshalb schon in §1;
[`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`](../observations/BEO-ALL/pfad-nachzug-gleicher-name-andere-datei/state.md)
(1×) — jeder Nachzug zeigt auf dieselbe Datei;
[`BEO-ALL/externer-link-auf-umgezogene-datei`](../observations/BEO-ALL/externer-link-auf-umgezogene-datei/state.md)
(1×) — die Release-Notiz nennt den endgültigen Ort.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-09 (`make nightly-state`) — `upstream-drift` grün;
`image-scan` rot (zwei HIGH in der Standardbibliothek des publizierten
Images, behoben ab Go 1.27.2, der Pin ist auf `main`).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
