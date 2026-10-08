# slice-260: Festlegungen der lokalen Wächter, Hooks und Prüfer in die Spezifikation

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`MR-074`](../../../../harness/conventions.md#mr-074) (Bewegung 2,
Rest); Folge von slice-259, der den Abschnitt §7 der Spezifikation anlegt.

**Berührte Spec-Stellen:** `spec/spezifikation.md` §7.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die Festlegungen der lokalen Wächter, Hooks und Prüfer, die nicht
in `make gates` laufen und keine eigene Anforderung verfeinern, stehen in §7
der Spezifikation, und ihre Sensor-Dateien verlinken die Kennung: der
Tool-Call-Wächter (was blockiert wird, fail-closed), das Handoff-Gate aus
Stop-Hook und `record-gates` (Inhalts-Hash, Schleifen-Schutz, Freigabe ohne
Nachweis), die git-Hooks (welcher Übergang welche Prüfung auslöst) und
`blackbox-probe` (Kanarienlauf, Abbruch, Vergleich).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Abschnitt selbst und die Gates aus `make gates`** — slice-259.
- **Netz- und Nachtlauf-Werkzeuge** (`image-scan`, die Versions- und
  Digest-Achsen, `baseline-freshness`, `nightly-state`) — slice-262
  übernimmt sie; geteilt beim Schnitt, weil die Abgrenzung neun Werkzeuge mit
  eigener Festlegung ergab, über der Grenze aus §4.
- **Was eine Anforderung durchsetzt** (`image-test`, `trace-check`,
  `adr-check`, die Hooks, soweit sie diese nur rufen) — deren Festlegung ist
  Verfeinerung der Anforderung in §1, nicht §7 (Vorlage §7).
- **Werkzeuge, die nur bewegen oder sagen** (`archive-wave`, `slice-mv`,
  `help`, `clean`, `versions`) — sie treffen keine Festlegung, die man
  fortschreiben müsste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Je Werkzeug mit eigener Festlegung ein §7-Eintrag, am Code geprüft.
- [ ] Die Sensor-Dateien verlinken die Kennung statt Schwelle und Randform zu
      führen; `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft;
      [MR-074](../../../../harness/conventions.md#mr-074) Bewegung 2 mit dem
      Anteil dieses Slice vermerkt (eingelöst, sobald auch slice-262 schließt).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | §7-Einträge |
| `harness/sensors/*.md` der betroffenen Werkzeuge | update | Verweis auf die Kennung |
| `harness/sensors/{guard-probe,hooks,blackbox-probe}.md`, `harness/README.md` (Zeile `record-gates`) | update | die Sensor-Dateien bzw. die Index-Zeile der vier Werkzeuge |
| `.githooks/pre-commit`, `.github/workflows/ci.yml` | update | der Closure-Übergangs-Wächter folgt der Festlegung |

*(Plan-Änderung vor dem Code: Am Code nachgelesen erkennt der
Closure-Übergangs-Wächter in `pre-commit` und in der CI nur einen Slice, der
**direkt** unter `done/` landet (`^docs/plan/planning/done/slice-…\.md$`).
Seit 2026-09-29 schließen Slices nach `done/wellenlos/` bzw.
`done/<welle-id>/` — gemessen über die Move-Commits: 20 Closures seither
(17 nach `done/wellenlos/`, 3 nach `done/welle-91/`),
keine davon direkt unter `done/`; der Wächter hat bei keiner ausgelöst.
`verify-closure-notes` lief nur, weil es von Hand gefahren wurde. Die
Festlegung wird so geschrieben, dass jeder Rename/Add eines
`slice-*.md` irgendwo unter `done/` den Wächter auslöst, und beide Stellen
folgen ihr. Der Ausschluss archivierter Stubs entfällt: Ein Stub-Add löst
einen Lauf über den ganzen Bestand aus, der ohnehin grün sein muss — Kosten,
keine falsche Zusage.)*

## 4. Trigger

**Start** (`next` → `in-progress`): slice-259 in `done/`; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): ein Werkzeug trifft mehr als eine
  Festlegung, die sich nicht in einer §7-Zeile fassen lässt — dann je
  Werkzeug ein Slice.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Keine bekannt.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** geändert werden die Spezifikation und
die Harness-Doku unter dem Default `*` (`ALL`); deklariert. `tools/harness/`
(`HARN`) und `.claude/hooks/` werden nur gelesen — die Festlegung beschreibt
den Wächter, sie ändert ihn nicht; berührt ist `HARN` damit nicht. Ändert
sich das im Lauf (ein Wächter folgt seiner Festlegung, wie das
Coverage-Skript in slice-259), ist das eine Plan-Änderung vor dem Code.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** gelesen am 2026-10-08.
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
(verkörpert, wach; in slice-259 dreimal getroffen) — jede Liste in §7 (was
der Wächter blockiert, welcher Übergang welche Prüfung auslöst) wird am Code
gezählt, nicht aus der Sensor-Datei übernommen;
[`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/state.md)
(2×) — trifft die Sensor-Dateien, die hier auf die Kennung umgestellt werden;
ein dritter Treffer wäre eine Lücke;
[`BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft`](../observations/BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft/state.md)
— die fail-closed-Randformen des Wächters werden als Festlegung geschrieben,
also an ihren Fehlerformen geprüft. Außerhalb dieses Slice gefunden:
[`BEO-ALL/begruendung-traegt-entscheidung-nicht`](../observations/BEO-ALL/begruendung-traegt-entscheidung-nicht/state.md)
steht seit slice-222 bei drei Belegen auf `offen`, ohne Ausgang — gemeldet,
nicht hier gelöst.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-08 aus dem jüngsten Lauf (`make nightly-state`) —
`upstream-drift` grün (11:37 UTC), `image-scan` grün (10:38 UTC).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
