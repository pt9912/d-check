# Review slice-296 — Entscheidungsvorlage zu den Neuerungen von `v7.0.0`, Runde 2 (letzte)

- **Review-Art:** Doku (gegen Plan, die Befunde der Runde 1, Hard Rules)
- **Gegenstand:** `slice-296`, Fix-Commits `433447b1..582f10b3` (`b620261c` Plan-Notiz, `582f10b3` Vorlage)
- **Skill:** `reviewer.md` @ 1.20.0
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-11-27
- **Eingangs-Kontext:** Report der Runde 1 (F-1 bis F-7), `slice-296` §2, §3 und §6, `wave-98` §5 (Review-Vorgabe),
  `ADR-0002`, `AGENTS.md` §3.5, §3.6 und §5, Baum `v7.0.0`
- **Zusätzlich gefahren** (Wegwerf-Bäume im Scratchpad): echter `git worktree` unter `.claude/worktrees/w1` mit genau den
  vier vorgeschlagenen Zeilen aus 5b — d-check im Standard-, `planning`-, `targets`-, `workflows`-, `mentions`-,
  `reviews`-, Closure- und Completeness-Lauf sowie a-check; Nennung von `harness/image-hash.txt` in
  `harness/README.md` mit und ohne Datei.
- **Nicht gefahren:** `make gates` im Klon mit Arbeitsbaum (auf dem Host liefen fremde Builds, Swap belegt); damit auch
  nicht `lint`, `test`, `coverage-gate`, `semgrep` dort. Die `make test`-Aussagen in B1 und B2 bleiben ungeprüft.

## Stand der Befunde aus Runde 1

| Befund | Stand | Geprüft |
|---|---|---|
| F-1 (5b, Preis) | **geschlossen** | Mit den vier Zeilen alle gefahrenen d-check-Läufe 0 Befunde, a-check 0; „Closure-Profil braucht nichts" bestätigt |
| F-2 (7a, 7b ohne Bedingung) | **geschlossen** | Zeilen 4b, 7a, 7b, 17c und beide Wortlaute tragen Review-Art und Reichweite, 7b den ganzen Absatz |
| F-3 (Ausnahme ohne Gegenstand) | **geschlossen** | 7b sagt es; das genannte `git grep` ist leer; Empfehlung „nichts tun", die Anker-Frage steht in der Übersicht |
| F-4 (fehlende Grenze) | **geschlossen** | neue Zeile 2f, 2d verweist, B1 nennt die Ziffern-Bedingung auf beiden Preisseiten; Messung deckt sich mit Runde 1 |
| F-5 (Klasse der Check-Zeilen) | **geschlossen** | 1b, 3, 18a, 18c, 28 „neu und anwendbar", A6 und A9 bis A20 „Lücke (wahlfrei)"; A6 nennt `pinsync`, das Docker-only-Repo und den A8-Fund; 1c und A9 widersprechen sich nicht mehr |
| F-6 (Zitat als Bedingung) | **geschlossen** | Schein-Bedingungen entfernt; 4d und 18b als Lesart ausgewiesen |
| F-7 (A5) | **geschlossen** | A5 „Lücke", Regel und `ADR-0002` genannt, Empfehlung „übernehmen"; A2, Kurzfassung und Übersicht nachgezogen |

F-8 bis F-12 (LOW, INFO) sind nach der Review-Vorgabe der Welle nicht eingearbeitet; sie bleiben im Report der Runde 1.

## Neue Findings

Keine HIGH, keine MEDIUM. Zwei LOW, zwei INFO.

### R2-1 — LOW — A5 gibt `ADR-0002` und den Preis des Nicht-Wollens nur halb wieder

- **quelle:** `AGENTS.md` §5 Regel 13 und 16
- **pfad:** Vorlage, Zeile A5 · „Sie nicht zu wollen, bräuchte eine Folge-ADR, weil die ADR sie nennt."
- **befund:** Der Satz trägt: `ADR-0002` ist `Accepted` und nennt die Datei, `AGENTS.md` §3.5 lässt eine Korrektur nur als
  neue ADR zu. Zweierlei fehlt. Derselbe Spiegelstrich der Fitness Function kündigt ein CI-Gate auf digest-stabilen
  Image-Bau an, zu dem die Datei das Beleg-Manifest ist; ein solches Gate finde ich nicht
  (`git grep -i 'digest-stabil' -- Makefile .github tools` ist leer), und die fünf Zeilen lösen es nicht ein. Und die
  Datei verlangt auch die Baseline: Wer sie nicht will, braucht dafür zusätzlich den Adaptions-Eintrag, den §So lesen
  für jede bewusst nicht befolgte Regel nennt.
- **verifizierbar:** ja — Lesen von `ADR-0002` §Fitness Function
- **klasse:** grenz-aussage-ohne-ihre-bedingung

### R2-2 — LOW — „Lücke (wahlfrei)" ist eine vierte Antwort, die sich selbst zu widersprechen scheint

- **quelle:** `slice-296` §3 Schritt 4 (Antworten *befolgt · Lücke · nicht anwendbar*)
- **pfad:** Vorlage, §So lesen · „*Wahlfrei* heißt: anwendbar, aber nicht verlangt."
- **befund:** Erklärt ist *wahlfrei*, nicht, was „Lücke" an etwas nicht Verlangtem heißt. Frage A zeigt damit in der
  Antwort-Spalte dasselbe Wort für echte Lücken (A3, A5, A7, A8) und für vierzehn Zeilen, die keine Pflicht berühren
  (A6, A9 bis A20), darunter sechs Abschnitte einer Sensor-Vorlage.
- **Wirkung auf die Entscheidung:** keine gefunden — A6 sagt „Keine Pflicht" und „nicht übernehmen", alle übrigen
  folgen A6, die Kurzfassung und die Übersicht führen nur die echten Lücken zur Entscheidung.
- **verifizierbar:** ja — Lesen der Tabelle
- **klasse:** vorlage-setzt-sitzungswissen-voraus

### R2-3 — INFO — Wortlaute 2f und A5 brechen vor dem Teil ab, auf den die Zeile sich beruft

2f beruft sich auf die „Pflicht, die Lücke zu benennen"; die zwei Zitate enden mitten im Satz und enthalten die
Pflicht („im Gate-Index als Lücke benennen") nicht. Wortlaut A5 endet vor dem Dateinamen und enthält „verlangt"
nicht. Beides steht an der Quelle so, wie die Zeilen es sagen. **klasse:** zitat-endet-vor-der-aussage

### R2-4 — INFO — 7b stellt eine eigene Frage ohne eigene Empfehlung und ohne eigenes Feld

„Größer, und eine eigene Frage an Sie" (HIGH-Anker im Reviewer-Skill): Die Zeile hat ein Entscheidungsfeld und die
Empfehlung „nichts tun", die dem Kleinsten gilt. Der Auftraggeber muss seine Antwort auf die zweite Frage kenntlich
machen. **klasse:** zeile-buendelt-neuerungen

## Antworten auf die vier Prüfpunkte

- **(a) A5, Folge-ADR:** Die Aussage stimmt. Kapitel (`v7.0.0` · `regelwerk/modul-14-docker-harness.md`
  §Multi-Stage-Build, „verlangt:") und `ADR-0002` nennen die Datei; die ADR steht in der Source Precedence über Briefing
  und Harness-Doku, ein Adaptions-Eintrag allein überschreibt sie nicht. Unvollständig ist sie nach R2-1. Zum Preis:
  Eine Nennung der Datei in `harness/README.md` färbt `make doc-check` rot (`codepath-missing`), solange sie beim
  Prüflauf nicht existiert; „die Datei beim Bau schreiben" deckt das.
- **(b) 5b, vier Zeilen in drei Dateien:** Der Preis trägt in allem, was ich gefahren habe; die Bedingung (ein Werkzeug
  legt den Baum ins Repo) steht im Zeilentitel, die vier Zeilen gelten dem Pfad `.claude/worktrees/`. „`make gates`
  grün" habe ich nicht nachgefahren.
- **(c) Summen, Verweise, Übersicht:** ohne Befund. 40 Dateinummern, eine Zeile mehr (2f); A5 und 7b stehen in der
  Übersicht wie in ihren Zeilen; 2f folgt B1, 17b folgt 7b, 17c folgt 7a, die Check-Zeilen folgen A6.
- **(d) „Lücke (wahlfrei)" ohne Sitzungsgeschichte:** richtig entscheidbar, holprig zu lesen (R2-2).

## Negativbefunde

- **Fix-Diff gegen die Plan-Notiz:** geprüft, ohne Befund — geändert sind die dort genannten Stellen, sonst nichts.
- **Neue Messungen in 5b:** geprüft — die Ausgangslage (4054 ohne Ausschluss, 100 `citation-out-of-range` und a-check rot
  mit nur `scan.ignore` und `.gitignore`) deckt sich mit der Messung aus Runde 1.
- **Bereits befolgt / befolgt:** 2d, 2e, 9b, 10a, A1, A2 im Fix nur um Verweise ergänzt; ohne neuen Befund.
- **Keine Übernahme eingeschleust:** geprüft, ohne Befund — der Diff berührt Vorlage und Plan; die Spalte des
  Auftraggebers ist leer.
- **Hard Rules** (`AGENTS.md` §3): geprüft, ohne Befund — kein Code, keine Konfiguration, keine ADR, keine Spec im Diff.

## Zu den Risiken aus `slice-296` §6 (Eingabe für die Closure)

Die in Runde 1 gemeldeten Abweichungen zu den Risiken 1, 3 und 5 sind in der Vorlage behoben; eingetreten waren sie.
Risiko 2 (Bündel) steht unverändert (F-10, R2-4). Risiko 4: nicht eingetreten.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 2 (R2-1, R2-2) | 2 (R2-3, R2-4) |

## Verdikt

**Abnahmefähig.** F-1 bis F-7 sind geschlossen, der Fix bringt nichts, was nach der Review-Vorgabe der Welle noch
einzuarbeiten wäre. Die LOW- und INFO-Befunde blockieren nicht.
