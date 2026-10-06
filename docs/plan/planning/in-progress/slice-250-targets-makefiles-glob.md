# slice-250: `targets.makefiles` nimmt Glob-Muster an

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`DC-FA-TGT-001`](../../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(Scope — erweitert, kein neues Kürzel: Einzelmodul-Frage nach
[ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md));
Anlass ist der eingehende
[CR von `ai-harness-init`](../../cr/2026-10-06-cr-eingehend-ai-harness-init-targets-makefiles-glob.md);
Vorbild der Konfigurations-Weitung
[ADR-0058](../../adr/0058-konfigurations-flaechen-additiv-weiten.md);
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(ohne Glob-Eintrag byte-identisch).

**Berührte Spec-Stellen:** [§DC-FA-TGT-001.a](../../../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
(Schritte 1 und 2), [`SPEC-005`](../../../../spec/spezifikation.md#spec-005--d-checkyml)
(Schema-Zeile `targets.makefiles`).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Ein `targets.makefiles`-Eintrag mit Glob-Zeichen (`*`, `?`, `[`)
expandiert gegen die Repo-Wurzel (`matchGlob`, segmentweise, `**`), jede
Treffer-Datei wird wie ein wörtlicher Eintrag gelesen; ein Glob ohne Treffer
bricht mit Exit 2 ab, eine doppelt erfasste Datei zählt einmal, ein
ungültiges Glob ist beim Laden Exit 2.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Globs für `targets.doc-tables` und `targets.authority`** — vom CR
  ausdrücklich ausgenommen, und ohne Anlass wäre es Vorsorge; `authority`
  ist zudem eine einzelne Datei mit eigener Semantik.
- **`include`-Auflösung im Makefile** — anderer Vorgang: der CR schließt sie
  aus, und das Lastenheft führt sie als Out-of-Scope (statische Heuristik,
  kein Ausführen).
- **Regel-Erkennung und `exempt-targets`** — Bestand bleibt bewusst stehen,
  vom CR so verlangt.
- **Skopierung über `scan.roots`/`scan.ignore`** — Bestand: wörtliche
  Einträge sind heute unabhängig davon, der Glob ist dieselbe ausdrückliche
  Angabe. Die Grenze wird in die Anforderung geschrieben (`AGENTS.md` §3.8).
- **Benutzerhandbuch, README, CHANGELOG** — Release-Prep
  (`AGENTS.md` §5 Regel 17).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Liefer-Punkte (3):

- [ ] Glob-Expansion im Kern: Treffer sortiert, Dubletten (wörtlich + Glob,
      zwei Globs) einmal gelesen, Fundstelle = echter Dateipfad; Glob ohne
      Treffer ⇒ Exit 2; Tests nach den sechs CR-Akzeptanzkriterien,
      darunter das rote Gegenbeispiel (`gate-undocumented` in einer per Glob
      erfassten Datei) — die Tests liefen ohne die Änderung aus dem
      richtigen Grund rot (Bewusstes Brechen, Modul 11).
- [ ] Config-Rand: Glob-Einträge segmentweise validiert (ungültiges Glob ⇒
      Exit 2), bestehende Pfad-Regel unverändert; `--print-config`-Gerüst mit
      Glob-Beispiel.
- [ ] Lastenheft (Erweiterung, Bump + Historie nach
      [`MR-032`](../../../../harness/conventions.md#mr-032)), Spezifikation,
      Schema; ADR; Antwort-Vermerk im CR.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      Rollenwechsel nach Schritt 8, kein Self-Review.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Anforderung erweitert, Akzeptanzkriterien, Out-of-Scope, Bump + Historie |
| `spec/spezifikation.md` | update | Schritte 1/2, Schema-Zeile, Historie |
| `internal/hexagon/core/rules/targets.go` | update | Expansion vor `collectMakefileRules` |
| `internal/adapter/driven/configyaml/configyaml.go` | update | Glob-Validierung in `applyTargets` |
| `internal/adapter/driving/cli/config_template.go` | update | Gerüst-Beispiel |
| `internal/hexagon/core/rules/targets_test.go`, `configyaml_test.go` | update | Akzeptanzkriterien 1–6 des CR |
| `docs/plan/adr/0099-…`, `docs/plan/adr/README.md` | neu / update | Entscheidungen |

**Spiegel vor dem Editieren** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
Lastenheft (Beschreibung „fail-closed" und Out-of-Scope der Anforderung),
Spezifikation (Schritte 1/2, Schema-Zeile `targets.makefiles`),
`--print-config`-Gerüst, `--doctor`-Klartexte (betroffen nur, falls eine
neue Meldung entsteht — Exit 2 ist kein Grund-Code), Benutzerhandbuch und
`operations.md` (Release-Prep).

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt; `in-progress/`
leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): verlangt die Expansion einen neuen
  Port oder eine Änderung am Filesystem-Adapter, wird der Zuschnitt neu
  geschnitten.
- `in-progress` → `open` (blockiert): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Ein Glob wie `**/*.mk` würde über die festen Überspring-Verzeichnisse
  (`build`, `vendor`, …) hinwegsehen — dieselbe Grenze wie beim Modul
  `file`; sie muss in der Anforderung stehen, nicht nur im Code. —
  **Ausgang:** *(offen)*
- Exit 2 bei leerem Glob kann einen Adopter in der Bootstrap-Phase treffen,
  in der das Fragment-Verzeichnis noch leer ist. — **Ausgang:** *(offen)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das Produkt
samt Spec-Stratum (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** drei offene Einträge der
Sub-Area treffen den Gegenstand —
[`BEO-ALL/stilles-gruen-ueber-leerer-range`](../observations/BEO-ALL/stilles-gruen-ueber-leerer-range/state.md)
(ein leerer Glob ist dieselbe Klasse: Prüfung ohne Gegenstand — deshalb
Exit 2),
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/state.md)
(2× — die Tests tragen die Nachbar-Fälle: wörtliche Einträge unverändert,
Dubletten) und
[`BEO-ALL/semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/state.md)
(Spiegel in §3 vorab gelistet).

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-06 — `image-scan` grün; `upstream-drift` rot mit denselben
drei Fremd-Meldungen wie am Vortag (Baseline-Release, semgrep-Version,
golang-Basis-Digest). Keine berührt diesen Gegenstand.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
