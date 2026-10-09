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

| `Makefile` (`gates`) | update | der Nachweis entsteht nur nach grünen Gliedern, auch unter `make -k` |

*(Plan-Änderung vor dem Code: `record-gates` steht als letzter
**Prerequisite** von `gates`; der Kommentar sagt, `make` breche vorher ab. Unter
`make -k gates` stimmt das nicht — `-k` arbeitet die übrigen Prerequisites
weiter ab, der Nachweis entsteht trotz rotem Glied, und der Stop-Hook gäbe
frei (gemessen an einem Modell-Makefile mit `.NOTPARALLEL`: Glied rot,
`record-gates` läuft, Exit 2). Der Nachweis wandert in das **Rezept** von
`gates`, das unter `-k` erst läuft, wenn alle Prerequisites grün sind; das
Target `record-gates` bleibt als Werkzeug stehen.)*

| neuer Eintrag unter `harness/conventions/`, Index in `harness/conventions.md`, Kopfkommentar `tools/harness/record-gates.sh` | neu/update | die Härtung am Handoff-Gate landet als eigener Eintrag, der [`MR-004`](../../../../harness/conventions.md#mr-004) schärft (Baseline-Regelwerk `modul-13-quality-gates.md` §Guard-Härtung); der Kopfkommentar des Skripts nennt den Ort des Aufrufs |

*(Plan-Änderung vor dem Code: [`MR-004`](../../../../harness/conventions.md#mr-004)
beschreibt `record-gates` als letzten Prerequisite von `gates` — mit dem Fix
aus der vorigen Plan-Änderung stimmt das nicht mehr, und ein akzeptierter
Eintrag wird nicht umgeschrieben.)*

| `.claude/hooks/stop-require-gates.sh`, `harness/sensors/verify-closure-notes.md`, `Makefile` (Hilfe-Zeile `blackbox-probe`) | update | R1-Befunde |

*(Plan-Änderung nach R1, vor dem Code: Der erweiterte Übergangs-Wächter löst
aus, aber `make verify-closure-notes` liest keine Unterverzeichnisse von
`done/` (R1 F-1, gemessen). Das zu beheben braucht eine Produkt-Änderung und
eine Stub-Unterscheidung — es übernimmt slice-263 (Auftraggeber-Entscheid
2026-10-09). Hier werden die Aussagen ehrlich: Hook-Kommentar,
[`SPEC-095`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge),
`hooks.md` und `verify-closure-notes.md` sagen, dass der Lauf auslöst, aber
nur `done/` selbst prüft. Dazu: die Erkennung liest die ganze Diff-Ausgabe,
statt bei `grep -q` per SIGPIPE abzubrechen (F-2); `make -i gates` schreibt
keinen Nachweis — Erkennung und Schreiben stehen in einer Rezeptzeile, weil
`-i` auch deren Abbruch ignorierte (F-3, am Modell-Makefile gegen `-i`,
`-ik`, `--ignore-errors`, `-j2 -i`, `-k`, `-s`, `-w` und eine Variable mit
`i` gefahren); scheitert im Stop-Hook der Hash, blockt er mit Grund statt
ohne Antwort zu enden (F-5); die Grenzen von
[`SPEC-093`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
nennen die gemessenen Durchlass-Klassen (F-4, F-8); `PROBE_FORMS` leer heißt
Default (F-6); Pointer- und Abschnitts-Korrekturen F-7, F-9, F-10, F-11.)*

*(Plan-Änderung nach R2, vor dem Code: Die Verträge von
`verify-closure-notes.md` und `hooks.md` versprechen die Prüfung weiter für
ganz `done/`; sie werden auf die Slices direkt unter `done/` eingeschränkt, die
Grenzen-Liste von `verify-closure-notes.md` nennt die Unterverzeichnisse
(R2-F-1). Der Stop-Hook blockt auch, wenn `git status` scheitert — der
Zustand ist dann so wenig gelesen wie bei einem gescheiterten Hash (R2-F-2,
Entscheidung: nachziehen statt benennen, wie bei F-5). Die §8-Zeile der
Spezifikation vom 2026-10-08 wird wiederhergestellt, der Nachzug steht in
einer neuen Zeile (R2-F-3). slice-263 nennt in §3 die vier Aussagen, die er
mit der Behebung zurücknimmt (R2-F-4). [`SPEC-093`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) fasst den Extraktor-Fall
enger (R2-F-5).)*

*(Plan-Änderung nach der Verifikation, vor dem Code:
[`MR-076`](../../../../harness/conventions.md#mr-076) nimmt die Härtung des
Stop-Hooks auf — Geltungsbereich und Adaption —, sonst stünde sie ohne
Eintrag (V-1). Die Durchlass-Klasse „Flag hinter einem Präfix" kommt in
Grenze 1 von `guard-probe.md` und in den GRENZE-Kommentar des Wächters (V-2).
[`SPEC-094`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
nennt die absichtliche Umgehung `MAKEFLAGS=` (V-3). [`MR-074`](../../../../harness/conventions.md#mr-074) und slice-262
nennen die drei Werkzeuge, die noch keinem Slice zugeordnet sind (V-5).)*

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
**Nachtrag nach den Plan-Änderungen:** `tools/harness/record-gates.sh` und der
Stop-Hook sind jetzt geändert; `HARN` ist berührt, Modus GF wie deklariert
(konventionsgetragen über [`MR-004`](../../../../harness/conventions.md#mr-004), geschärft durch
[`MR-076`](../../../../harness/conventions.md#mr-076)).

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
