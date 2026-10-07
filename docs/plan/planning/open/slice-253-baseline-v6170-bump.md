# slice-253: Baseline-Pin auf `v6.17.0` — mechanisch, ohne Urteil über den Delta

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`MR-073`](../../../../harness/conventions.md#mr-073) (Vorgänger in
der Pin-Serie, Auflösungs-Trigger „die nächste Pin-Hebung"),
[`MR-021`](../../../../harness/conventions.md#mr-021) (pin-gebundene Verweise),
[`MR-039`](../../../../harness/conventions.md#mr-039) (Zitat-Delta im neuen
Eintrag), [`MR-051`](../../../../harness/conventions.md#mr-051)
(`d-check:cite`-Neu-Ankern), [`MR-055`](../../../../harness/conventions.md#mr-055)
(Symlinks), [`MR-069`](../../../../harness/conventions.md#mr-069)
(`ignore-refs`-Stufe), [`MR-070`](../../../../harness/conventions.md#mr-070)
(Frozen-Klassen vor der Ersetzung). Anlass: Nachtlauf `upstream-drift`
(neuere Baseline-Tags) und das Release
[`v6.17.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.17.0).

**Berührte Spec-Stellen:** — *(pin-gebundene Verweise in den Rolle-Zeilen der
Spec-Straten sind Pfad-/Versionsnennungen, keine Spec-Aussage)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Den committet vendorten Baseline-Bestand auf
[`v6.17.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.17.0)
heben (den aktuellen Tag, nicht die Zwischen-Tags), den Pin in
`harness/conventions.md` §Baseline fortschreiben, alle lebenden
pin-gebundenen Referenzen retargeten und die Hebung als neuen MR-Eintrag
führen — Frozen-Klassen und Zitat-Delta gemessen am echten Vorzustand.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Keine Adoption der inhaltlichen Deltas** — jede ist eine Konventions-
  oder Produkt-Entscheidung mit eigener Begründung; der Bump ist mechanisch
  (anderer Vorgang). Der neue MR-Eintrag nennt die Schlagzeilen; Adoptieren
  gehört in eigene Folge-Slices nach Auftraggeber-Freigabe.
- **Frozen-Dateien bleiben byte-stabil** — `done/`-Slices und -Wellen,
  `docs/reviews/`, `Accepted`-ADRs, eingehende CRs, `conventions/done/MR-*`,
  CHANGELOG- und Spec-Historie nennen den alten Pin als
  Vergangenheits-Aussage korrekt (Bestand bewusst).
- **`tools/harness/selbstpruefung.sh`** — sein Stand-Vermerk nennt die
  Version des **Schwester-Werkzeugs**, aus dem es adoptiert ist, nicht den
  Baseline-Pin dieses Repos (anderer Gegenstand).
- **Dependabot-PR #6** — eigener Vorgang nach diesem Slice.
- **Kein Release** — der Bump berührt keine Distribution-Fläche.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Vendierter Baum steht auf `v6.17.0` (`SHA256SUMS` generiert und
      verifiziert), der `v6.13.0`-Baum ist entfernt; §Baseline-Pin nachgezogen;
      der Vorgänger-Eintrag nach `conventions/done/`, neuer MR-Eintrag im Index.
- [ ] Alle lebenden pin-gebundenen Referenzen nennen `v6.17.0` (Pfad-Verweise,
      Release-/Tree-URLs, bare Nennungen in lebenden Dokumenten, Symlinks unter
      `.claude/rules/`); `d-check:cite`-Spannen neu geankert; Zitat-Delta und
      Frozen-Liste im MR-Eintrag gemessen.
- [ ] `make gates` grün; `make baseline-freshness` meldet den Pin aktuell.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      Rollenwechsel nach Schritt 8, kein Self-Review.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/baseline/{v6.13.0 → v6.17.0}/` | neu + entfernen | Re-Vendor beider Bäume + `SHA256SUMS`; alter Baum weg |
| `harness/conventions.md` | update | §Baseline-Pin, §Adoptierte Konventions-Quellen, MR-Index |
| neue MR-Datei unter `harness/conventions/` | neu | Hebung-Eintrag (Delta, Frozen-Liste, Zitat-Delta, Cite-Neu-Ankern) |
| Vorgänger-MR-Datei der Pin-Serie | move | reiner `git mv` nach `conventions/done/`, dann Link-Tiefen |
| `AGENTS.md`, `harness/README.md`, `harness/rules/*`, aktive `MR-*`, `.claude/agents/reviewer.md`, `roadmap.md`, `planning/README.md`, `observations/README.md`, `spec/architecture.md`, `spec/spezifikation.md` (Rolle-Zeile) | update | lebende Token-Swaps, je Datei geprüft (Mischfälle zeilenweise) |
| `.claude/rules/*.md` (Baseline-Symlinks) | update | Aliase auf `v6.17.0` |
| `.d-check.yml` | update | nur, falls gemessen: nächste `ignore-refs`-Stufe |
| dieser Plan | update | die beiden `d-check:cite`-Spannen auf den neuen Baum |

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt; `in-progress/`
leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): ändert der Delta die Struktur des
  Bundles (neue Bäume, entfallene Dateien, die verlinkt sind), sodass der
  Bump nicht mehr mechanisch ist.
- `in-progress` → `open` (blockiert): das Release-Bundle ist nicht ladbar.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Ein lebender Verweis zeigt auf einen Anker oder eine Datei, die es im neuen
  Baum nicht mehr gibt (Umbenennung im Delta). — **Ausgang:** *(offen)*
- Eine Mischfundstelle wird pauschal ersetzt und fälscht eine
  Vergangenheits-Aussage. — **Ausgang:** *(offen)*

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

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: die
Harness-Mechanik (`tools/harness/`, Kürzel `HARN`) samt der pin-gebundenen
Verweise im ganzen Repo (`*`, `ALL`); beide deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** zwei Einträge treffen den
Gegenstand —
[`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/state.md)
(die Frozen-Klassen werden vor der Ersetzung gelistet, nach
[`MR-070`](../../../../harness/conventions.md#mr-070)) und
[`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/state.md)
(Spiegel-Klassen gemessen statt angenommen).

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-07 — `upstream-drift` rot mit der Baseline-Meldung, die
diesen Slice auslöst; dazu semgrep und golang-Digest (je eigener Vorgang);
`image-scan` grün.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
