# Trigger-Audit welle-97 — Zweige ADR und Hard Rule, Lese-Schritt des Registers (Architect)

**Rolle:** Architect (Modul 8, §Rollen-Sequenz für eine Welle: Closure-Schritt 2 und 3b, Planner → Architect →
Planner). Rolleninhaber: pt9912.
**Gegenstand:** [welle-97](../plan/planning/welle-97-linelength-eingerueckte-tabellen.md) mit
[slice-293](../plan/planning/done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md) und
[slice-294](../plan/planning/done/slice-294-linelength-eingerueckte-tabellen-kern.md), Range `cb164eaa^..HEAD`
(28 Commits, Spitze `5c8ded99`).
**Eingang:** Audit-Vorlage des Planners; Review und Verifikation von slice-294
([R1](2026-11-24-slice-294-linelength-eingerueckte-tabellen-r1.md),
[Verifikation](2026-11-24-slice-294-linelength-eingerueckte-tabellen-verify.md)). Die repo-weite Verifikation der
Welle entsteht parallel und ist nicht Eingang.
**Nicht Gegenstand:** die Zweige Carveout und bootstrap-aware Gate (prüft der Planner).
**Datum:** 2026-11-24.
**Geändert wurde nichts:** keine ADR, keine Hard Rule, kein Register-Eintrag; dies ist nur das Verdikt.

## Gesamturteil

Ein Re-Evaluierungs-Trigger ist eingetreten und in der Welle beantwortet: [ADR-0114](../plan/adr/0114-linelength-modul.md)
Trigger (b), beantwortet durch die Folge-ADR [ADR-0116](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md)
(`Accepted` in `41ee1279`, teil-supersedes Entscheidung 5 (b)). **Eine weitere Folge-ADR ist nicht nötig.** Kein
anderer Trigger ist eingetreten, **keine Hard-Rule-Zeile** wird entfernt.

Für [`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../plan/planning/observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md)
(3×) schlage ich **`verkörpert`** vor, als Erweiterung der bestehenden Regel 13 in `AGENTS.md` §5 (ein Halbsatz in
der Tabellenzeile, ein Absatz im Volltext), *permanent*, `seit welle-97` — keine neue Regel (Teil 2).

**Keine Entscheidung des Auftraggebers nötig.**

Gate-Lage über die Welle:
`git diff --stat cb164eaa^..HEAD -- .d-check.yml Makefile .golangci.yml tools/coverage-gate.sh` ⇒ **leer**. Keine
Schwelle, kein Profil, kein Gate geändert. `linelength` ist opt-in und im eigenen Profil nicht an.

## Teil 1 — ADR

| ADR | Trigger | eingetreten? | Beleg |
|---|---|---|---|
| [ADR-0114](../plan/adr/0114-linelength-modul.md) | (a) die geteilte Lexik ändert ihre Fence- oder Tabellen-Antwort | **nein** | `tableRowLine` ist zeichengleich; der Diff von `markdown.go` über die Welle ist reiner Zuwachs (28 Einfügungen, keine Löschung), Verifikation slice-294 Punkt 1 |
| ADR-0114 | (b) ein Konsument muss eine Zeile ausnehmen, die weder Code noch Tabelle noch unumbrechbar ist | **ja, beantwortet** | Issue #20: die eingerückte Tabellenzeile war nach Entscheidung 5 (b) keine Tabelle. Die Antwort ist nicht der Marker, auf den der Trigger verweist, sondern ADR-0116; ADR-0116 vergleicht den Marker als Alternative G und verwirft ihn wie ADR-0114 Entscheidung 9. Der Vermerk steht in der `## Geschichte` von ADR-0114 und im Index („teil-superseded: ADR-0116“). Der Trigger bleibt für andere Zeilenarten bestehen |
| ADR-0114 | (c) Grobheit von (c) · (d) breite Zeichen · (e) andere Dateiart · (f) weitere Lint-Prüfung | **nein** | kein Antrag, kein Review-Befund dazu |
| [ADR-0116](../plan/adr/0116-linelength-eingerueckte-tabellenzeilen.md) | (1) der Bestand des Melders behält Befunde an Tabellen mit vier oder mehr Leerzeichen | **nein** | noch nicht veröffentlicht; der Bestand liegt außerhalb des Repos |
| ADR-0116 | (2) die Lexik bekommt eine Antwort auf den Listenkontext · (4) `> \|` verlangt · (5) eine Prosa-Zeile mit `\|` als zu Unrecht ausgenommen gemeldet | **nein** | keine solche Änderung, kein Antrag, kein Befund |
| ADR-0116 | (3) ein zweiter Konsument will die eingerückte Tabellenzeile | **nein** | `unbreakableTableRow` hat einen Aufrufer außerhalb der Tests, `linelength.go` Z. 75 |
| [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) | (1) ein Fall der Klasse tritt im Bestand ein | **nein** | die neue Antwort ist eine eigene, durch ADR-0116 gescopte Frage in der geteilten Lexik (Entscheidung 2), kein Konsument beantwortet sie selbst. Gebunden nach Entscheidung 4: `TestTabellenzeilenFrageHatEineAntwort` hält die drei Spalte-0-Konsumenten und `tableRowLine` auf einer Antwort, eine zweite Assertion hält den Helfer in beiden Richtungen; die Mutationen „`tableRowLine` geweitet“ und „Helfer ruft `tableRowLine` nicht“ sind rot (Verifikation slice-294) |
| ADR-0054 | (2) `spans` meldet `fence-unclosed` in einer Datei unter `vcs.paths` | **nein** | `make doc-check` an `HEAD` 0 Befunde |
| ADR-0054 | (3) eine **vierte** Stelle beantwortet eine Lexik-Frage selbst | **durch die Welle nicht eingetreten** | Review R1 zu slice-293 fand zwei Leser, die vor dem `\|` trimmen: den RTM-Leser (`SplitPipeTableLine`) und den Beobachtungs-Leser in `planning` (`planning_observations.go` Z. 152–154). Beide sind älter als die Welle; der Trigger ist laut Geschichte von ADR-0054 schon einmal eingetreten und mit der Kopplung beantwortet. Die Welle hat sie sichtbar gemacht und das Lastenheft auf die drei Spalte-0-Konsumenten eingeschränkt. Keine Folge-ADR aus diesem Audit; ob ein Gate die Klasse prüfen kann, bleibt die Frage des Triggers, wenn eine weitere Stelle dazukommt |
| ADR-0113, ADR-0115 | — | **nicht berührt** | `git diff --stat cb164eaa^..HEAD -- docs/plan/adr` nennt nur ADR-0114 (eine Geschichte-Zeile), ADR-0116 und den Index; `pinsync` und `vcs` unverändert |

## Teil 1 — Hard Rules

**`AGENTS.md` §3** — die Datei ist über die Welle unverändert
(`git diff --stat cb164eaa^..HEAD -- AGENTS.md` leer).

| Regel | Auflösungs-Trigger | eingetreten? |
|---|---|---|
| §3.1 Docker/make-only | keiner benannt (Bestand) | — |
| §3.2 Suppression-Verbot | permanent | nein |
| §3.3 `git mv` + Inhalt = zwei Commits | permanent | nein — angewandt: die sechs Lifecycle-Moves der Welle zeigen `R100` |
| §3.4 Architektur sprachfrei, Straten nie abwärts | Closure-Datum keiner; Rollen statt Technologie permanent | nein — `spec/architecture.md` unberührt |
| §3.5 ADRs nach `Accepted` immutable | durchgesetzt (`make adr-check`) | — angewandt: ADR-0114 bekam nur einen Anhang in `## Geschichte`, die Status-Zeile blieb `Accepted`; `make adr-check` grün über `HEAD~1..HEAD` und `--staged` (der gestagte Fall in einem Wegwerf-Klon gefahren) |
| §3.6 Gates nicht ohne ADR lockern | permanent | nein — Gate-Konfiguration unverändert (oben) |
| §3.7 Kommentare tragen eine der fünf Klassen | keiner benannt | — Review R1 zu slice-294 F-1 (LOW): eine *Grenze* mit einer Begründung, die nur einen Teil ihrer Fälle trägt; ins Register (Teil 2), kein Trigger |
| §3.8 Ein Modul verspricht nur über seine Scan-Menge | permanent | nein — der Helfer liest dieselben Zeilen wie zuvor |
| §3.9 Action-Referenzen SHA-gepinnt | durchgesetzt (`make workflow-pins`) | — keine Workflow-Änderung |

**`AGENTS.md` §5** (Volltexte unter `harness/rules/dokumentations-regeln/`):

| Regel | Auflösungs-Trigger | eingetreten? |
|---|---|---|
| 13 Grenzen gegen den Gegenstand prüfen | permanent | nein — angewandt und an zwei Stellen nur von Reviews eingelöst (Teil 2) |
| 15 Botschaft behauptet nicht mehr als die Arbeit | permanent | nein — die Commit-Botschaften der Welle nennen Belege; kein Review-Befund dazu |
| 18 Festlegung zitiert ihre Quelle | permanent | nein — angewandt: die Pläne zitieren die Kriterien des Lastenhefts statt sie nachzuerzählen (slice-293 §3, slice-294 §3) |
| 17 Release-Prep | die Spiegel-Tests halten keinen Fundort in einer README oder im Handbuch mehr | nein — `registry_mirror_test.go` hält die README-Bullets und die Handbuch-Tabelle weiter |
| übrige (1–12, 14, 16) | permanent oder keiner benannt | nein — kein Bezug zur Welle |

**0 offen** in beiden Klassen.

## Teil 2 — Lese-Schritt: `BEO-ALL/grenz-aussage-ohne-ihre-bedingung` (3×)

**Gelesen:** `observation.md`, `state.md` (offen) und die drei Belege
[`slice-291`](../plan/planning/observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/evidence/slice-291.md),
[`slice-293`](../plan/planning/observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/evidence/slice-293.md),
[`slice-294`](../plan/planning/observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/evidence/slice-294.md).

**Die Klasse:** Eine Grenze oder ihre Begründung steht an einer Stelle mit ihrer Bedingung. An einer zweiten Stelle
— Kriterium, Historie, Spezifikation, ADR-Probe, Kommentar, Testfall — steht sie ohne die Bedingung oder für Fälle,
die die Begründung nicht trägt. Alle drei Belege fand ein Review, keiner ein Gate. Die Stellen lagen im Lastenheft,
in der Spezifikation, in einer ADR und in Kommentaren, nur einmal in einem Plan.

**Ist sie schon gedeckt?** Teilweise, an drei Stellen, keine trifft den Kern:

- **Regel 13** („Wer eine Grenze aufschreibt, prüft sie gegen den Gegenstand“) verlangt, eine aufgeschriebene Grenze
  gegen Code und Vertrag zu prüfen. Sie sagt nichts darüber, dass eine **wiedergegebene** Grenze ihre Bedingung
  mitnimmt. Beleg slice-293 (eine Aussage über „die übrigen Module“, die am Gegenstand nur für drei stimmt) fällt
  unter Regel 13. Die Belege slice-291 (Bedingung beim Übertragen verloren) und slice-294 (Begründung für fremde
  Fälle) fallen nicht darunter.
- **Regel 18** deckt das Nacherzählen im **Slice-Plan** und verweist für Spec, ADR und Handbuch ausdrücklich auf die
  Regeln 3, 13 und 16. Die drei Belege liegen fast alle dort, also außerhalb von Regel 18.
- **Das Workflow-Skelett** (`.claude/commands/implement-slice.md`, Schritt 21, „seit slice-265“) verlangt die Suche
  über die **Aussage** statt über die Fundstelle — aber nur vor einem Fix-Commit nach dem Review. Zwei der drei
  Belege entstanden beim ersten Schreiben.
- Der Reviewer-Skill fand alle drei; ein weiterer MEDIUM-Anker dort säße hinter dem Entstehungsort.

**Vorschlag: `verkörpert`, Zielort Regel 13** (Hard Rule, `AGENTS.md` §5, Volltext
`harness/rules/dokumentations-regeln/13-grenzen-gegen-gegenstand-pruefen.md`). Begründung für die kleinste Form:

- **Keine neue Regel.** Regel 18 weist Spec, ADR und Handbuch schon an Regel 13. Eine neue Regel 19 stünde daneben
  und gäbe der Frage „wie schreibe ich eine Grenze“ zwei Orte. Die Erweiterung ist ein Halbsatz in der Tabellenzeile,
  die jeder Lauf liest, und ein Absatz im Volltext.
- **Nicht `gestrichen`.** Die Belege zeigen nicht, dass Regel 13 wirkt: Zwei der drei Fälle liegen außerhalb ihres
  Wortlauts, und alle drei fand erst ein Review. Dreimal in drei Slices hintereinander, trotz der ausdrücklichen
  Planvorgabe der Welle („jede Grenze mit ihrer Bedingung“), ist eine Lücke im Text, kein Zufall.
- **Nicht `geplant`.** Kein Werkzeug fängt die Klasse. Ob zwei Stellen dieselbe Aussage tragen, ist ein Urteil.
  `citations` hält nur ein deklariertes wörtliches Zitat, und eine Grenze wird umformuliert wiedergegeben, nicht
  zitiert. Ein Folge-Slice hätte keinen Gegenstand.

**Mechanisierbar?** Nein. Die Suche über den Inhalt ist ein Urteil, deshalb *permanent*.

**Wortlaut 1 — Tabellenzeile 13 in `AGENTS.md` §5**, ersetzt den Text der Spalte „Regel“ (Z. 323):

```markdown
| 13 | Wer eine Grenze aufschreibt, prüft sie gegen den Gegenstand (Code/Config), nicht gegen ihre eigene Beschreibung; wer sie an einer zweiten Stelle wiedergibt, nimmt ihre Bedingung mit. | [Volltext](harness/rules/dokumentations-regeln/13-grenzen-gegen-gegenstand-pruefen.md) |
```

**Wortlaut 2 — neuer Schlussabsatz im Volltext**, hinter dem bestehenden Anker „seit slice-213“:

```markdown
**Eine Grenze, die an einer zweiten Stelle wiedergegeben wird**, nimmt ihre Bedingung mit — in ein
Kriterium, eine Historie-Zeile, die Spezifikation, eine ADR, einen Kommentar oder einen Testfall. Für
eine Begründung gilt dasselbe: Sie steht nur bei den Fällen, die sie trägt. Vor dem Handoff wird jede
Stelle, die dieselbe Aussage trägt, über ihren **Inhalt** gesucht, nicht über ihren Wortlaut; eine Stelle,
die die Grenze ohne ihre Bedingung trägt, sagt für den ausgenommenen Fall das Gegenteil zu. Im Slice-Plan
trägt das [Regel 18](18-plan-festlegung-zitiert-ihre-quelle.md) (zitieren statt nacherzählen), vor einem
Fix nach dem Review die Suche über die Aussage im Workflow-Skelett; dieser Absatz gilt beim ersten
Schreiben und überall sonst. Urteil, kein `grep`. *(Hard Rule aus dem Steering Loop,
[`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../../../docs/plan/planning/observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md),
seit welle-97; Auflösungs-Trigger: permanent.)*
```

**Wortlaut 3 — Register und Closure-Notiz:**

- `state.md`: `**Stand:** verkörpert — Regel 13 (AGENTS.md §5, Volltext harness/rules/dokumentations-regeln/13-grenzen-gegen-gegenstand-pruefen.md, Schlussabsatz), seit welle-97.`
- Steering-Loop-Eintrag in `welle-97-results.md`: Regel · `BEO-ALL/grenz-aussage-ohne-ihre-bedingung` ·
  Belege slice-291, slice-293, slice-294 · liegt in `harness/rules/dokumentations-regeln/13-grenzen-gegen-gegenstand-pruefen.md`.

**Hinweise für den Commit:**
- Die Tabellenzeile 13 steht in `AGENTS.md` Z. 323, also hinter der einzigen Zeilen-Zitation in die Datei
  (`.harness/skills/reviewer.md` Z. 62, `AGENTS.md:287-287`). Ein längerer Zellentext ändert die Zeilenzahl nicht.
- Der Volltext bekommt einen zweiten Herkunfts-Anker; die Anker-Paarung prüft eine Datei ohne Suffix auf
  `seit welle-97` irgendwo in ihr.
- Der Reviewer-Skill zitiert Regel 13 im MEDIUM-Anker zu Prüffrage 18 nicht wörtlich. Der Skill bleibt deshalb
  unberührt.

## Nicht geprüft

- Die Zweige Carveout und bootstrap-aware Gate (Planner).
- Die repo-weite Verifikation der Welle (`make fullbuild`, Probe am Image, Vorher-/Nachher-Lauf). Sie entsteht
  parallel; Teil 1 stützt sich auf Review und Verifikation von slice-294 und auf das Audit am Stand `5c8ded99`.
- Die übrigen Register-Einträge der Welle (R1 F-1/F-2, V-1/V-2 der Verifikation von slice-294): Zählung des Planners.
