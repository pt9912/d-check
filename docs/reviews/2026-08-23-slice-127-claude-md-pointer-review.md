# Review-Report: slice-127 — CLAUDE.md auf einen Pointer eindampfen

**Review-Art:** Code- und Plan-Review (Slice-Plan · Hard Rules · Konventionen gegen den Diff), unabhängiger Reviewer ohne Anteil an der Arbeit
**Gegenstand:** `git diff ef28d3e..HEAD` — die drei Slice-Commits `ebc3299` (angelegt), `4e651c0` (beansprucht, Lifecycle-Move + Ruhe-Marker), `0dcdbd9` (CLAUDE.md 23 → 5 Zeilen)
**Skill:** `.harness/skills/reviewer.md` @ 1.9.0 · **Modell-ID:** `claude-opus-5[1m]` · **Datum:** 2026-08-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [`slice-127-claude-md-pointer.md`](../plan/planning/done/slice-127-claude-md-pointer.md)
- [`CLAUDE.md`](../../CLAUDE.md) vorher (`git show ef28d3e:CLAUDE.md`) und nachher
- [`AGENTS.md`](../../AGENTS.md) §1, §2, §3.1, §3.3, §3.7, §4, §5, §6
- [`harness/README.md`](../../harness/README.md) §Leseordnung, §Source precedence
- [`harness/conventions.md`](../../harness/conventions.md) §Modus-Deklaration; MR-004, MR-005, MR-012, MR-013, MR-015
- `.claude/hooks/pretooluse-command-guard.sh`, `.claude/hooks/stop-require-gates.sh`, `.claude/settings.json`, `.claude/commands/implement-slice.md`
- [`observations.md`](../plan/planning/observations.md) BEO-002, BEO-007, BEO-009, BEO-010, BEO-011
- [`roadmap.md`](../plan/planning/in-progress/roadmap.md) §Offene Wellen; Baseline `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht

**Vom Reviewer selbst gefahren** (Exit-Code direkt aus `$?`, Ausgabe in Datei umgeleitet — BEO-007):
`make doc-check` Exit 0 · `make planning-check` Exit 0 · `make gates` Exit 0 · beide `.claude/hooks/`-Skripte gegen selbst gebaute Payloads.

---

## Findings

### F-1 · MEDIUM — „Konflikt melden" hat im ganzen Repo keine zweite Fundstelle

- `kategorie`: MEDIUM
- `quelle`: BEO-011 (Vollständigkeits-Aussage aus dem Anlass statt aus dem Bestand); AGENTS.md §1
- `pfad`: `git show ef28d3e:CLAUDE.md`:23 · `AGENTS.md`:14-16 · Botschaft `0dcdbd9`, Zuordnungstabelle Zeile 7
- `befund`: Die gestrichene Zeile „Bei Quellen-Konflikt: Konflikt **melden** und der höherrangigen Quelle folgen" trägt zwei Aussagen, die AGENTS.md §1 nicht trägt — die **Meldepflicht** und den Fall **zweier kanonischer Quellen untereinander** (§1 regelt ausschließlich „zwischen dieser Datei und einer kanonischen Quelle"); `harness/README.md`:9-10 und `harness/conventions.md`:19-20 sagen für den eigenen Fall „die Datei wird angepasst", nicht „melden". `grep -rn "Konflikt melden\|melden und der"` über das Repo findet die Pflicht außerhalb des Slice-Plans selbst nirgends, während derselbe Plan sie in §3 („das ist **zu melden**") und DoD-Punkt 4 als geltend voraussetzt.
- `verifizierbar`: nein — kein Gate prüft Prosa-Aussagen; der Befund ist per `grep` über den Wortlaut reproduzierbar.
- `klasse`: `vollstaendigkeits-aussage-ohne-zweite-fundstelle`

### F-2 · MEDIUM — „IDs vor der Implementierung → AGENTS.md §6.3" deckt drei der fünf Positionen nicht

- `kategorie`: MEDIUM
- `quelle`: BEO-011; AGENTS.md §6.3
- `pfad`: `git show ef28d3e:CLAUDE.md`:19-20 · `AGENTS.md`:299 · `.claude/commands/implement-slice.md`:13-18
- `befund`: AGENTS.md §6.3 lautet vollständig „Betroffene Requirement-/ADR-IDs identifizieren" — **Slice-ID**, **betroffene Module** und **auszuführende Gates** stehen dort nicht, und „identifizieren" ist nicht die Nenn-Pflicht („benennen"), die die gestrichene Zeile aufstellte. Die vollständige Fünferliste existiert im Repo nur in `.claude/commands/implement-slice.md` Schritt 7 — weder in einer der drei im Plan genannten Quellen (`AGENTS.md`, `harness/README.md`, `harness/conventions.md`) noch im automatisch geladenen Kontext; sie greift erst, wenn der Nutzer `/implement-slice` aufruft.
- `verifizierbar`: nein — reproduzierbar per `git grep -n "auszuführende Gates" ef28d3e`, das genau zwei Treffer liefert, beide in der alten `CLAUDE.md`.
- `klasse`: `zuordnung-deckt-teilmenge`

### F-3 · MEDIUM — „AGENTS.md §6.8 (wörtlich)" ist nicht wörtlich

- `kategorie`: MEDIUM
- `quelle`: BEO-009 Richtung (a) — die Botschaft behauptet eine Eigenschaft, die die Quelle nicht hat
- `pfad`: Botschaft `0dcdbd9`, Zuordnungstabelle Zeile 6 · `AGENTS.md`:304 · `.claude/commands/implement-slice.md`:26
- `befund`: AGENTS.md §6.8 sagt „keine Erfolgsmeldung ohne Gate-**Ausführung**", die gestrichene Zeile sagte „Kein Erfolg ohne **echte Gate-Ausgabe**" — Ausführung ist der Lauf, Ausgabe ist der vorzuzeigende Beleg; `grep -n "Gate-Ausgabe" AGENTS.md` liefert keinen Treffer. Der wörtliche Zwilling steht in `.claude/commands/implement-slice.md`:26 („Do not claim completion without command output"), also gerade nicht in der als Quelle benannten Datei.
- `verifizierbar`: nein — reproduzierbar per `grep` auf beide Wortlaute.
- `klasse`: `woertlich-behauptet-sinngemaess-belegt`

### F-4 · MEDIUM — der Pointer verspricht eine „Leseordnung", die AGENTS.md nicht führt

- `kategorie`: MEDIUM
- `quelle`: MR-015 (AGENTS.md routet); Maintainability
- `pfad`: `CLAUDE.md`:3-5 · `AGENTS.md` (gesamte Datei) · `harness/README.md`:17-33
- `befund`: Die neue, einzige Sachzeile nennt vier Ziele — „Source Precedence, Leseordnung, Workflow, Gates" —, aber `grep -n "Leseordnung\|Leseliste\|Lesereihenfolge" AGENTS.md` endet mit Exit 1: das Wort kommt in AGENTS.md nicht vor, und §6.1 sagt nur „`harness/README.md` lesen", ohne zu nennen, was dort zu finden ist. Wer §6.1 folgt, landet in `harness/README.md` §Leseordnung Punkt 1, der auf AGENTS.md zurückzeigt — genau der Zirkel, den der Slice in §3 als out-of-scope deklariert; aufgelöst hatte ihn bisher die gestrichene CLAUDE.md-Leseliste, die `harness/README.md` explizit auf Platz 1 stellte.
- `verifizierbar`: nein — kein Gate misst, ob ein Pointer sein Ziel erreicht; reproduzierbar per `grep` (Exit 1).
- `klasse`: `pointer-nennt-nicht-vorhandenes-ziel`

### F-5 · MEDIUM — die Hook-Begründung reicht weiter als die Hooks

- `kategorie`: MEDIUM
- `quelle`: reviewer.md 1.9.0 §MEDIUM *Botschaft verallgemeinert über die Messung hinaus*; BEO-009 Richtung (b); MR-004/MR-005
- `pfad`: `docs/plan/planning/done/slice-127-claude-md-pointer.md`:50-56 · `.claude/hooks/pretooluse-command-guard.sh`:28-33 · `.claude/hooks/stop-require-gates.sh`:17-26 · `.claude/settings.json`:5
- `befund`: (a) Der PreToolUse-Guard trägt nur die **negative** Hälfte der Regel: selbst gefahren gibt er bei `{"tool_input":{"command":"docker run img npm test"}}` und bei `python -c "…"` keine Ausgabe (= frei), blockt aber `pip install foo` — „Nur `make`-Targets für Checks und Gates" ist damit nicht erzwungen, und `"matcher": "Bash"` schließt jedes andere Tool (Datei-Schreibzugriffe, MCP-Aufrufe) von vornherein aus. (b) `stop-require-gates.sh` hat **zwei** bedingungslose Freigabepfade statt des einen, den sein Kommentar als Restlücke nennt: neben dem frischen Klon mit cleanem Tree gibt die Schleifen-Schutz-Klausel bei `{"stop_hook_active": true}` — selbst gefahren — `{"decision":"approve"}` zurück, unabhängig vom Gate-Nachweis. Die Plan-Formulierung „gibt den Stop nur frei, wenn der Repo-**Inhalt** durch einen erfolgreichen `make gates`-Lauf gedeckt ist" trifft für keinen der beiden Pfade zu; BEO-007 hält zusätzlich fest, dass der Stop-Hook „nur das Sitzungs-Ende, nicht den Commit dazwischen" fängt.
- `verifizierbar`: ja — `bash .claude/hooks/pretooluse-command-guard.sh < payload.json` bzw. `bash .claude/hooks/stop-require-gates.sh < payload.json`; beide Läufe sind oben zitiert und wurden vom Reviewer gefahren.
- `klasse`: `mechanik-begruendung-ueber-die-mechanik-hinaus`

### F-6 · LOW — Zuordnungstabelle Zeile 1 und Zeile 2 sind ungenau

- `kategorie`: LOW
- `quelle`: BEO-009 Richtung (a); Maintainability
- `pfad`: Botschaft `0dcdbd9`, Zuordnungstabelle Zeilen 1 und 2 · `harness/README.md`:17-33, :35 · `AGENTS.md`:59-71, :297-299
- `befund`: `harness/README.md` §Leseordnung führt **vier** Zeiger (AGENTS.md · conventions.md · aktiver Slice samt Roadmap · Regelwerk bei Bedarf) und nennt weder „referenzierte ADRs unter `docs/plan/adr/`" noch „referenzierte Anforderungen unter `spec/`"; AGENTS.md §6.2 sagt dazu nur generisch „Relevante kanonische Quelle lesen" — zwei der sechs Leselisten-Punkte haben in den beiden angegebenen Quellen keine explizite Entsprechung. Zeile 2 behauptet außerdem „AGENTS.md §2 (dort steht die **Tabelle** selbst)": §2 ist eine nummerierte Liste, die Tabelle steht in `harness/README.md` §Source precedence.
- `verifizierbar`: nein — reproduzierbar durch Lesen beider Abschnitte.
- `klasse`: `zuordnung-deckt-teilmenge`

### F-7 · LOW — MR-013 wird für eine Richtung zitiert, die sein Geltungsbereich nicht nennt

- `kategorie`: LOW
- `quelle`: MR-013; AGENTS.md §3.3
- `pfad`: Botschaft `4e651c0` · `harness/conventions/MR-013-lifecycle-move-buendelung.md`:6-9 · `AGENTS.md`:108
- `befund`: MR-013 nennt als Geltungsbereich ausdrücklich „der Slice-Lifecycle `docs/plan/planning/in-progress/` → `…/done/`" (plus die MR-/Wellen-Moves), und die AGENTS.md-§3.3-Ausnahme trägt dieselbe Einschränkung; der hier gebündelte Move ist `open/` → `in-progress/`. Die Kopplung, die MR-013 begründet, ist symmetrisch und der Commit ist gate-richtig, aber der deklarierte Geltungsbereich deckt die Beanspruchungs-Richtung nicht — wer §3.3 wörtlich liest, baut hier einen byte-reinen Move-Commit.
- `verifizierbar`: ja — ein byte-reiner `open/`→`in-progress/`-Move ohne Marker-Flip lässt `make planning-check` rot laufen (Marker-Wächter in beide Richtungen).
- `klasse`: `regel-geltungsbereich-schmaler-als-praxis`

### F-8 · INFO — „in dieser Sitzung zweimal ausgelöst" ist vom Reviewer nicht nachprüfbar

- `kategorie`: INFO
- `quelle`: BEO-009 Richtung (a)
- `pfad`: `docs/plan/planning/done/slice-127-claude-md-pointer.md`:53 · Botschaft `ebc3299`
- `befund`: Der Guard schreibt kein Log, und im Repo existiert kein Artefakt, gegen das die Zahl zu halten wäre — die Aussage ist weder zu bestätigen noch zu widerlegen. Sie trägt in §1 argumentatives Gewicht („die Sorge … trägt nicht") und steht damit als einzige Kennzahl des Slice ohne reproduzierbaren Beleg neben vier nachgefahrenen Gate-Zahlen.
- `verifizierbar`: nein — kein Sitzungs-Artefakt im Repo.
- `klasse`: `beleg-ohne-artefakt`

---

## Negativbefunde (geprüft, ohne Befund)

1. **Zeilenweiser Abgleich alt → neu, alle acht Streichungen.** Tragfähig belegt sind: Zeile 3 (Greenfield/Doc führt) → `harness/conventions.md`:157; Zeile 16 (Source Precedence) → `AGENTS.md`:59-71; Zeilen 17-18 (make-only, Toolchain-Verbot) → `AGENTS.md`:75-89, dort sogar breiter (Host-Go ist genannt); Zeile 21 (`make gates`) → `AGENTS.md`:302, dort präziser (Handoff). Die Zeilen 19-20, 22 und 23 tragen F-1 bis F-3.
2. **„23 Zeilen auf 5".** Bestätigt: `git show ef28d3e:CLAUDE.md` = 23 Zeilen, `CLAUDE.md` = 5; `git show --stat 0dcdbd9` meldet 3 insertions / 21 deletions — 23 − 21 = 2 unveränderte Kopfzeilen + 3 neue = 5. Rechnung geht auf.
3. **„452 Dateien, 0 Befunde" (beide Botschaften).** Selbst gefahren: `make doc-check` Exit 0, Ausgabe `d-check: 452 Datei(en) geprüft, 0 Befund(e)`; `make planning-check` Exit 0, identische Zeile.
4. **„acht Gates, Coverage 94,80 % gegen Schwelle 93 %".** Selbst gefahren: `make gates` Exit 0; `coverage-gate: OK — Coverage 94.80% erfüllt Schwelle 93%`; `Makefile`:158 listet acht Gates (`doc-check lint test arch-check coverage-gate semgrep gate-consistency planning-check`) plus `record-gates` (Nachweis, kein Gate). Zahl und Schwelle stimmen.
5. **Eingehende Verweise, alle Dateitypen.** `grep -rn "CLAUDE" .` über `.claude/`, `.github/`, `Makefile`, Skripte, YAML und Markdown: neun Sach-Treffer (acht Review-Reports, `.claude/commands/implement-slice.md`:7); der einzige Nicht-Doku-Treffer ist `CLAUDE_PROJECT_DIR` in `.claude/settings.json` und betrifft die Datei nicht. Kein Treffer trägt einen Anker; `make doc-check` Exit 0 bestätigt, dass kein Verweis ins Leere läuft. `CLAUDE.md` liegt unter `scan.roots: ["."]`, ist also tatsächlich gescannt.
6. **„sieben Review-Reports".** Nachgezählt: acht Report-Dateien nennen CLAUDE.md, davon sieben als Datei-Link und eine (`docs/reviews/2026-08-02-slice-091-status-feld-review.md`:45) als Inhalts-Zitat ohne Link. Die Botschaft zählt beide Klassen getrennt und korrekt.
7. **Das frozene Inhalts-Zitat.** „`AGENTS.md` §3.3 (laut `CLAUDE.md` das zuerst gelesene Doc)" war unter der alten Fassung **falsch** — dort stand AGENTS.md auf Platz 2 hinter `harness/README.md` — und wird durch die neue Fassung („zuerst `AGENTS.md` … lesen") richtig. Die Botschafts-Aussage „wird durch die Kürzung richtiger, nicht falscher" hält.
8. **`.claude/commands/implement-slice.md` bleibt funktionsfähig.** Schritt 1 liest CLAUDE.md (jetzt den Pointer), Schritt 3 liest AGENTS.md ohnehin als eigenen Schritt; die Schritte 2, 4-6 spiegeln die alte Leseliste weiter. Kein toter Verweis, kein Bruch im Kommando-Ablauf.
9. **Wellenlosigkeit.** Baseline `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht: „Wellenlose Arbeit erscheint nicht in der Roadmap — weder beim Start noch beim Abschluss." Es ist keine Welle-Datei offen, kein Zeiger gesetzt; die Form des `**Welle:**`-Feldes ist deckungsgleich mit den Präzedenzen `slice-112` und `slice-121`. Die Begründung („eigene DoD, keine gemeinsame Closure-Bedingung") trifft das Baseline-Kriterium.
10. **Ruhe-Marker-Zustand.** `in-progress/` trägt seit `4e651c0` einen Slice, also muss „Nichts in Arbeit." weichen; der Marker-Wächter misst beide Richtungen, `make planning-check` Exit 0 am HEAD bestätigt den Zustand.
11. **Move-Commit-Form.** `git show --stat 4e651c0` zeigt `{open => in-progress}` mit 0 Zeilen Änderung an der Slice-Datei (R100, Rename-Detection hält) plus genau die eine Roadmap-Zeile. Der Slice-Body (DoD-Haken, Closure-Notiz) ist unberührt geblieben — MR-013-konform bis auf die Richtungs-Frage in F-7.
12. **Pflicht-Vorprüfungen.** §7 „Vorgelagert" trägt beide (Sub-Area · offene Beobachtungen sichten) **vor** der Modus-Begründung §8 — `AGENTS.md`:287-290 erfüllt. BEO-011 ist als einschlägig benannt und mit dem richtigen Auftrag versehen („der Beleg ist zeilenweise zu führen, nicht zu behaupten"), BEO-002 für die Ränder.
13. **Slice-Kopf-Form.** `**Lifecycle:**` statt `**Status:**`, `**Berührte Spec-Stellen:** —` (kein Spec-Bezug), `**Verantwortlich:**` gesetzt — `AGENTS.md` §5 erfüllt; kein Zustandsfeld, das Chronik trägt.
14. **Out-of-Scope-Treue, alle vier Ausschlüsse.** `git diff --stat ef28d3e..HEAD` = drei Dateien: keiner der drei Drift-Defekte ist einzeln nachgezogen; `.d-check.yml` ist unverändert (kein `pins`/`dpin`, kein `FOCUS_DISABLE`-Nachzug, BEO-010 nicht angerissen); weder `AGENTS.md` §6 noch `harness/README.md` §Leseordnung sind angefasst; die neue CLAUDE.md-Zeile stellt keine Regel auf, die nicht in AGENTS.md steht (der Heredoc-Guard-Kandidat ist nicht mitgenommen).
15. **Traceability-Form.** Alle drei Botschaften nennen eine `MR-`-ID und `slice-127`; `make trace-check` ist bewusst nicht Teil von `gates`, die Form ist aber erfüllt.
16. **Kommentar- und Zustandsfeld-Hygiene (`AGENTS.md` §3.7).** Der Diff fügt keinen Code-, Konfigurations- oder Skript-Kommentar hinzu und kein Zustandsfeld; keine Slice-Nummer und kein Mess-Label in einer Kommentarzeile.
17. **Kosmetik ohne Konventions-Anker — kein Finding.** Die Marker-Entfernung hinterlässt in `roadmap.md` nach Zeile 33 zwei aufeinanderfolgende Leerzeilen. Kein Gate und keine Repo-Konvention bindet das (kein markdownlint in `Makefile` oder CI); nach dem Anti-Pattern *Kein Stil-Polizist* nicht gemeldet, nur festgehalten.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --------- | ------ | -------- |
| HIGH      | 0      | — |
| MEDIUM    | 5      | F-1, F-2, F-3, F-4, F-5 |
| LOW       | 2      | F-6, F-7 |
| INFO      | 1      | F-8 |

---

## Verdikt

**Blockierend** — kein HIGH, fünf MEDIUM, zwei LOW, ein INFO.

Der Slice ist handwerklich sauber gefahren: Lifecycle-Move in R100-Form, Ruhe-Marker
korrekt geflippt, alle vier Out-of-Scope-Zusagen gehalten, und jede nachprüfbare
Zahl der drei Botschaften stimmt (452/0, acht Gates, 94,80 %, 23 → 5). Die
Vollständigkeits-Aussage, auf der die **Entfernung** ruht — „trägt nichts Eigenes,
jede Zeile steht schon in AGENTS.md / `harness/README.md` / `harness/conventions.md`" —
hält der zeilenweisen Gegenprobe jedoch an zwei Stellen nicht stand: die Meldepflicht
bei Quellen-Konflikt (F-1) existiert im Repo nirgends sonst, und die Nenn-Pflicht vor
der Implementierung (F-2) steht vollständig nur in `.claude/commands/implement-slice.md`,
also außerhalb der drei genannten Quellen und außerhalb des auto-geladenen Kontexts.
Der Plan hat diesen Ausgang selbst vorgesehen: §6 *Rückführungen* nennt
`in-progress` → `next` für den Fall, „dass CLAUDE.md doch Eigeninhalt trägt, den
AGENTS.md nicht abdeckt".

F-5 betrifft die zweite Säule derselben Begründung: die Hooks tragen weniger, als der
Plan ihnen zuschreibt — der PreToolUse-Guard nur die negative Hälfte der make-only-Regel
und nur für das Bash-Tool, der Stop-Hook mit einem zweiten, im Plan nicht genannten
bedingungslosen Freigabepfad. F-4 trifft das Ergebnis selbst: von den vier Zielen, die
die verbleibende Zeile verspricht, ist eines in AGENTS.md nicht auffindbar.

F-6 bis F-8 blockieren nicht; F-7 ist ein Befund an der Konvention, nicht an diesem
Slice, und F-8 ist ausdrücklich als nicht nachprüfbar markiert statt durchgewinkt.
