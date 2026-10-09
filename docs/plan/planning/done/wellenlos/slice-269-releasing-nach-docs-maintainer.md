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

- [x] Ein Move-Commit, den Git als Rename erkennt, mit den eingehenden
      Verweisen; danach der Commit für die eigenen Links — die Liste am Repo
      gezählt, das Kommando im Plan; `make doc-check`, `make adr-check` über
      die Range und `make gates` grün, auch in einem frischen Klon.
- [x] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [x] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| docs/user/maintainer/releasing.md → docs/maintainer/releasing.md | move | Umzug |
| alle Dateien mit Verweis auf die Datei | update | Pfad-Nachzug |
| `.d-check.yml` | update | `matrix.exempt-paths`, Tombstone |
| [`MR-077`](../../../../../harness/conventions/MR-077-releasing-doku-unter-docs-maintainer.md), `harness/conventions.md` | neu, update | der neue Ort weicht vom Baseline-Rang 6 ab |
| `AGENTS.md` §2, `harness/README.md` §Source precedence | update | Rang 6 nennt beide Verzeichnisse |

Gemessen im Wegwerf-Klon mit dem reinen Umzug, dem gelöschten leeren
Verzeichnis und `d-check` über den Klon (Befunde nach Grund und Datei
gezählt): 40 Befunde. 23 eigene Links (12 `target-missing`, 11 `repo-escape`),
16 eingehende `target-missing` in denselben Dateien wie bei slice-268
(CHANGELOG einer, README und README.de je zwei, Benutzerhandbuch zwei,
Docker-Hub-Beschreibung,
[ADR-0014](../../../adr/0014-latest-tag-fuer-stabile-releases.md) vier,
[ADR-0067](../../../adr/0067-dependabot-als-hebender-kanal.md) einer, drei
Wellen-Ergebnisnotizen) und ein `codepath-missing`: die Closure-Notiz von
slice-268 nennt das Zwischen-Verzeichnis in Inline-Code — es kommt ins
Tombstone-Register ([ADR-0025](../../../adr/0025-codepaths-ignore-refs.md)).
Ohne Doku-Gate: `.github/dependabot.yml`, `release.yml` (zwei),
`hub-description.yml`, `tools/image-test.sh`, der Kommentar zum
Tombstone-Eintrag des ursprünglichen Pfads.

Das Zählkommando, im Wegwerf-Klon des Stands vor dem Umzug (nachgetragen
nach der Verifikation, V-1):

```sh
mkdir -p docs/maintainer
git mv docs/user/maintainer/releasing.md docs/maintainer/releasing.md
rmdir docs/user/maintainer
docker run --rm --network none -v "$PWD":/repo:ro d-check:latest \
  | awk -F'\t' 'NF>2 {split($1,a,":"); print a[1], $3}' | sort | uniq -c
```

*(Plan-Änderung vor dem nächsten Code-Commit, nach R1 MEDIUM-1: Mit dem
Umzug verlässt `releasing.md` den Rang 6 der Source Precedence, den
`AGENTS.md` §2 und `harness/README.md` §Source precedence `docs/user/`
zuordnen — „Operations, Releasing". Der Baseline-Default legt die
Releasing-Sicht unter `docs/user/*.md`; der neue Ort ist deshalb eine
Abweichung und bekommt einen `MR`-Eintrag in `harness/conventions.md`; Rang 6
nennt in beiden Tabellen beide Verzeichnisse. Gemessen: keine Scan-Menge der
`.d-check.yml` hängt an `docs/user/` außer einer `structure`-Regel für das
Benutzerhandbuch. (Berichtigt nach R2: `make mention-coverage` prüft nur die
ADR-Dateien gegen ihren Index und sagt über `docs/user/` nichts.))*

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
  Arbeitsbaum allein verdeckt das. — **Ausgang:** *entfallen* — das Risiko
  war echt (39 Befunde mit stehendem Verzeichnis, 40 ohne) und ist für diesen
  Umzug abgefangen: das Verzeichnis ist entfernt, die Nennung steht im
  Tombstone-Register, `doc-check` und der ADR-Check sind im frischen Klon
  grün. Die Klasse steht im Register
  ([`BEO-ALL/leeres-verzeichnis-lokal-verdeckt-ci-befund`](../../observations/BEO-ALL/leeres-verzeichnis-lokal-verdeckt-ci-befund/state.md)).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Releasing-Doku liegt unter
  docs/maintainer/releasing.md, wo der Auftraggeber sie haben wollte; Git
  erkennt den Rename (100 %). Der Schnitt aus dem Register von slice-268 —
  eingehende Verweise im Move-Commit, die bewegte Datei unverändert — stand
  diesmal vor dem Code im Plan und lief im Hook beim ersten Anlauf grün. Der
  zweite Pfad-Nachzug in [ADR-0014](../../../adr/0014-latest-tag-fuer-stabile-releases.md)
  und [ADR-0067](../../../adr/0067-dependabot-als-hebender-kanal.md) ging durch
  `make adr-check`. Im frischen Klon gemessen: `doc-check` und der ADR-Check
  über die Range grün.
- **Was ging anders als geplant:** R1 fand, dass die Datei mit dem Umzug
  den Rang 6 der Source Precedence verliert — `docs/user/` war dort
  „Operations, Releasing". Die Plan-Änderung kam vor dem Code:
  [`MR-077`](../../../../../harness/conventions.md#mr-077) deklariert den neuen
  Ort, beide Rangtabellen nennen ihn. Der Slice selbst existiert, weil
  slice-268 ein falsch übermitteltes Ziel umgesetzt hatte; die Korrektur kam,
  als dessen Closure lokal schon committet war, und lief vorwärts statt über
  ein Verwerfen von Commits (Auftraggeber-Entscheid). Zwei eigene
  Arbeitsfehler ohne Folge im Repo: ein `sed` mit leerer Zeilennummer schrieb
  die Plan-Notiz hinter jede Zeile (wiederhergestellt vor dem Commit), und ein
  `| tail` verschluckte den Exit eines roten `doc-check` — der
  `pre-commit`-Hook lehnte ab. Das Zählkommando stand erst nach der
  Verifikation im Plan (V-1).
- **Steering-Loop-Eintrag:** keiner mit Schwelle.
  `BEO-ALL/pipe-swallows-gate-exit-code` ist in welle-79 als Hook verkörpert,
  und der Hook hat hier getragen.
- **Beobachtungs-Register (`../../observations/`):** `evidence/slice-269.md` in
  [`BEO-ALL/pipe-swallows-gate-exit-code`](../../observations/BEO-ALL/pipe-swallows-gate-exit-code/state.md);
  neu
  [`BEO-ALL/leeres-verzeichnis-lokal-verdeckt-ci-befund`](../../observations/BEO-ALL/leeres-verzeichnis-lokal-verdeckt-ci-befund/state.md)
  (1×).
- **Folge-Slices:** keiner. Release v0.85.0 mit slice-263, slice-265,
  slice-267, slice-268 und diesem Slice; die Release-Notiz nennt
  docs/maintainer/releasing.md als neuen Ort.
- **Risiken aus §6:** eines entfallen (siehe §6). Trigger-Audit: kein
  Carveout, kein bootstrap-aware Gate, keine ADR; [`MR-077`](../../../../../harness/conventions.md#mr-077) neu, sein Trigger
  nicht eingetreten; keine Hard Rule mit eingetretenem Trigger.
  Nachtlauf-Stand ([`MR-053`](../../../../../harness/conventions.md#mr-053)): wie
  in §8 — `image-scan` rot bis zum nächsten Release.
- **Drei Paarungen:** (a) Anker — kein Eintrag mit `liegt in`; (b)
  Folge-Slices — keiner genannt; (c) Register — die zwei zitierten
  Beobachtungen existieren und tragen Belege.

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
[`BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise`](../../observations/BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise/state.md)
(1×) — der Schnitt dieses Slice folgt ihr und steht deshalb schon in §1;
[`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`](../../observations/BEO-ALL/pfad-nachzug-gleicher-name-andere-datei/state.md)
(1×) — jeder Nachzug zeigt auf dieselbe Datei;
[`BEO-ALL/externer-link-auf-umgezogene-datei`](../../observations/BEO-ALL/externer-link-auf-umgezogene-datei/state.md)
(1×) — die Release-Notiz nennt den endgültigen Ort.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-09 (`make nightly-state`) — `upstream-drift` grün;
`image-scan` rot (zwei HIGH in der Standardbibliothek des publizierten
Images, behoben ab Go 1.27.2, der Pin ist auf `main`).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
