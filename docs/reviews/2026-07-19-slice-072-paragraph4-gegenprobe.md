# Abschluss-Gegenprobe — slice-072 §4-Aufgabenorientierung (frisches Voll-Audit nach Remediation)

**Datum:** 2026-07-19 ·
**Auditor:** unabhängiger, adversarialer Doku-Auditor ·
**Gegenstand:** `docs/user/benutzerhandbuch.md`, `## 4. Aufgaben` **vollständig**
(§4.1–§4.16, aktuell Z. 146–857), gelesen im **Worktree-Stand** (die letzten
zwei Remediations-Bausteine B-7 + citations-§5-Titel sind noch uncommittet;
siehe Beleglage) ·
**Kriterium:** `docs/user/benutzerhandbuch-standard.md` §2 („Aufgaben statt
Funktionen"), §5 („Schritt-für-Schritt": Ausgangssituation → Ziel →
Voraussetzungen → nummerierte Schritte → erwartetes Ergebnis → Fehler-Hinweise),
§9/§11 (Versions-/Changelog-Angaben gebündelt), §3 (§5/§6 als legitime
Referenz-Kapitel) ·
**Checkliste:** die 14 Befunde B-1…B-8 / N-1…N-6 des Ursprungs-Audits
[`2026-07-19-slice-072-paragraph4-audit.md`](2026-07-19-slice-072-paragraph4-audit.md)
und das bereits erteilte Cluster-Review-Verdikt (ACCEPT)
[`2026-07-19-slice-072-paragraph4-cluster-review.md`](2026-07-19-slice-072-paragraph4-cluster-review.md).

## Beleglage / Methodik

- §4 vollständig gelesen (Z. 146–857); §5 (`trace`-Referenz Z. 1263–1464), §7
  (Z. 1490–1546) und §11 (Z. 1587–1637) als Verlagerungs- und
  Duplikations-Gegenprobe gelesen.
- **Commit-Zuordnung der Remediation** (`git log`): B-1/B-2 in `0442d3d`,
  B-3/B-4/B-5/B-6 + N-1…N-4 (§4.12-Cluster-Auftrennung) in `c17c34b`
  (Cluster-Review `f6bb802` = ACCEPT), Handbuch-Version 1.42 + §11-Zeile in
  `790299f`. Der **uncommittete Worktree-Diff** (`git diff`, +10/-9) trägt die
  letzten zwei Bausteine: **B-7** (das Fähigkeits-Inventar wandert in §4.4
  **hinter** die Entscheidungsverzweigung, Z. 252–257) und den **citations-§5-Titel**
  auf Task-Form (Z. 1022). Dieses Audit bewertet exakt diesen Worktree-Stand.
- Verlagerungen an der Datei verifiziert (Grep nach `Bis v`/`Ab v` im §4-Bereich,
  WAISE-Definitionsstellen, produkt-innensichtige Opener, Anker-/Test-Abhängigkeit
  der §5-Titeländerung).
- Regress-Sicherheit der uncommitteten Bausteine geprüft: der §4.4-Diff ist reines
  Prosa-Umordnen (kein Fence, kein Config-Block berührt); die §5-Titeländerung hat
  **keinen** Slug-Anker-Link (`grep '](#…citation…)'` leer) und **keine**
  Test-Harness-Zusicherung auf den Alt-Text („Zitat-Verifikation" existiert nur
  noch in ADR/Spec/Changelog/Review-Prosa, in keinem `*_test.go`). Ein Regress an
  `handbook_examples_test.go`/`docexamples_test.go` durch diese zwei Bausteine ist
  damit ausgeschlossen; der committete Cluster ist bereits per `make test` grün
  belegt (Cluster-Review Achse 4).

## Befund-für-Befund

| Item | Verdikt | Ein-Satz-Beleg (aktuelle Stelle) |
|---|---|---|
| B-1 | ERLEDIGT | §4.7 `order`/`direction` (Z. 334–356) trägt jetzt einen **Ziel-Titel** („Auch *innerhalb* einer Klasse eine Rangfolge erzwingen") und öffnet aus der **Lesersituation** (Z. 335: „Manchmal sollen selbst Dokumente derselben Klasse nicht abwärts verweisen …"); die Exit-2-Bedingung ist auf einen kurzen Fehler-Hinweis (Standard §5) reduziert. |
| B-2 | ERLEDIGT | §4.7 `token` (Z. 358–384) trägt Ziel-Titel („Auch bare ID-Tokens im Fließtext als Referenz prüfen") und öffnet lesersituativ (Z. 358–360: „Eine Referenz muss kein Markdown-Link sein …"); der Alt-Opener „`matrix` sieht standardmäßig nur Markdown-Links" ist weg. |
| B-3 | ERLEDIGT | Die Grammatik-/Feld-/Tabellen-Binding-Referenz lebt jetzt allein in §5 (Z. 1279–1330); §4.12 hält nur noch einen kurzen Syntax-Rahmen mit §5-Zeiger (Z. 619–621) — der ~67-Z.-Block ist entschärft, keine Verdopplung (Cluster-Review Achse 1/2). Residuum siehe NEU-1 (INFO). |
| B-4 | ERLEDIGT | Modalität ist eigene Aufgabe §4.14 „Nur MUSS-Anforderungen als Gate erzwingen" mit **Ausgangslage→Ziel→Vorgehen→Ergebnis** (Z. 695–722); Klassifikations-Mechanik + Fail-closed-Liste nach §5 (Z. 1384–1403) verwiesen (Z. 720–722). |
| B-5 | ERLEDIGT | Der „Versionsgrenze"-Blockquote ist weg (Grep „Versionsgrenze" leer); im gesamten §4 (Z. 146–857) steht **keine** „Bis vX/Ab vX"-Prosa mehr (nur die Image-Tags `:v0.51.1` in Befehlen); „0 Anforderungen" ist nach §7 (Z. 1531–1546), die Historie führt §11. |
| B-6 | ERLEDIGT | §4.12 ist in §4.12 (RTM/`--trace`) · §4.13 (`trace.coverage`) · §4.14 (`modality`) · §4.15 (`cross-consistency`) · §4.16 (`--print-mk`) aufgetrennt; die **doppelte WAISE-Definition** ist aufgelöst — Ziel nennt WAISE nur mit Vorwärts-Zeiger (Z. 591), die Definition steht **einmal** im Ergebnis (Z. 630–632). |
| B-7 | ERLEDIGT | §4.4: die Entscheidungsfrage „Welche Sie nehmen, hängt von Ihrer Ausgangslage ab" (Z. 228–229) und die Frisch-/Bestehend-Verzweigung (Z. 231/240) stehen **vor** dem Fähigkeits-Inventar, das nun nachgelagert ist (Z. 252–257) — genau die im DoD verlangte Umstellung (uncommitteter Diff). |
| B-8 | ERLEDIGT (war bereits) | Alle §4.x-Titel folgen „Nutzerziel + Flag/Schlüssel in Klammern" (§4.5/§4.11/§4.12–§4.16); kein reiner Flag-Titel auf §4.x-Ebene, und die vormaligen Feature-Name-Sub-Titel in §4.7 sind über B-1/B-2 auf Ziel-Form. |
| N-1 | ERLEDIGT | „Tabellengrenze" (Grammatik + „Bis v0.47.0") ist nach §5 (Z. 1324–1326) und §11 (Z. 1631) verlagert; §4.12 nennt sie nur als §5-Zeiger (Z. 619–621). |
| N-2 | ERLEDIGT | „Direktiven-Marker in einer Tabellenzeile" (+ „Bis v0.48.0") ist nach §5 (Z. 1326–1330) und §11 (Z. 1632) verlagert; §4.12 nur Zeiger. |
| N-3 | ERLEDIGT | „Zugesagte Range-/Aufzählungs-Notationen" (Komma-Kurzform ⇒ Exit 2) ist nach §5 (Z. 1378–1382) verlagert; §4.13 verweist nur (Z. 691–693). |
| N-4 | ERLEDIGT | Versions-/Kompatibilitätsprosa ist aus allen §4-Aufgaben entfernt (siehe B-5); §11 (Z. 1595–1637) führt die gebündelte Historie inkl. neuer Zeile 1.42 (Z. 1637). |
| N-5 | BEWUSST-REFERENZ | `citations`/`sources` haben keine eigene §4-Aufgabe, aber je eine **task-betitelte §5-Referenz** (Z. 1022 „Zitate und Zeilen-Referenzen gegen ihre Quelle prüfen", Z. 1158 „Externe Quellen auf ihren Inhalt pinnen"); die vom Ursprungs-Audit gerügte Inkonsistenz zu `pins` (Z. 1123) ist damit aufgehoben — bewusste Zurückhaltung für fortgeschrittene opt-in-Features, legitim nach Standard §3. |
| N-6 | BEWUSST-REFERENZ | §4.9/§4.11 erklären beide `--json`/`--yaml`; das Ursprungs-Audit stufte die Überlappung selbst als „akzeptable In-Task-Referenz" ein — unverändert, kein Standard-Verstoß. |

**Zählung der 14 Ursprungsbefunde:** 12 ERLEDIGT (B-1…B-8, N-1…N-4),
2 BEWUSST-REFERENZ (N-5, N-6), **0 OFFEN**.

## Neue Findings (frischer Blick auf ganz §4)

### NEU-1 — §4.12 „Unterstützte Definitionssyntax" öffnet noch referenz-förmig vor „Vorgehen"

- **kategorie:** INFO (Standard §5 Reihenfolge; identisch mit Cluster-Review F-3)
- **quelle:** Standard §5 (Ausgangssituation→Ziel→Vorgehen)
- **pfad:** `docs/user/benutzerhandbuch.md:596`–`621`
- **befund:** Zwischen **Voraussetzung** (Z. 593) und **Vorgehen** (Z. 623) steht
  der Block „**Unterstützte Definitionssyntax:**" (Z. 596–621, ~26 Z. mit zwei
  Markdown-Beispielen) — er öffnet mit „Standardmäßig erkennt `--trace` eine
  Anforderung aus einer ATX-Überschrift …", also leicht produkt-innensichtig. Die
  detaillierte Grammatik ist bereits nach §5 ausgelagert (B-3 erledigt); übrig
  bleibt der erklärende Rahmen, den §4.13/§4.14 in dieser Form **nicht** mehr
  tragen.
- **failure-szenario:** Ein Leser, der nur die Handlungsfolge sucht, liest erst
  ~26 Z. Definitions-Syntax, bevor der Befehl kommt — er wird **verzögert**, nicht
  fehlgeleitet (die Info „Anforderungen kommen aus Headings" ist sachlich richtig
  und für das Ziel nötig). Kein Standard-Verstoß, der den Leser real irreführt;
  vertretbar als „was der Leser einmal wissen muss".

### NEU-2 — §5 trägt eine verstreute Versions-Parenthese, die §11 schon führt

- **kategorie:** INFO (Standard §9/§11 Bündelung; außerhalb §4)
- **quelle:** Standard §9/§11 (Versions-/Changelog-Angaben gebündelt)
- **pfad:** `docs/user/benutzerhandbuch.md:1375`
- **befund:** Der §5-Absatz zur Range-Notation schließt mit „(Bis v0.44.0
  expandierte eine verlinkte Range gar nicht — das erzeugte in `trace.coverage`
  falsche Waisen.)" — dieselbe Aussage steht als §11-Zeile 1.31 (Z. 1626). Es ist
  die **einzige** verbliebene „Bis vX"-Parenthese außerhalb §11 und liegt in der
  **Referenz** §5, nicht in einer §4-Aufgabe.
- **failure-szenario:** Ein Maintainer, der die Versions-Historie in §11 pflegt,
  muss diese Compat-Notiz bei einer Aktualisierung an **zwei** Stellen führen
  (milde Doku-Drift-Gefahr). Kein Leser wird fehlgeleitet; §5 ist ein legitimes
  Referenz-Kapitel, in dem eine knappe Compat-Begründung vertretbar ist —
  optionaler Nachzug, kein Blocker.

**Keine neuen HIGH/MEDIUM.** Die neuen §4.13–§4.15 erfüllen real das
Ausgangslage→Ziel→Vorgehen→Ergebnis-Muster (§4.13 Z. 666–693, §4.14 Z. 697–722,
§4.15 Ausgangslage/Ziel + Schritt 1/2/3 Z. 726–798); die zwei Grep-Treffer für
mögliche Innensicht-Opener (Z. 349 „d-check prüft dann auch …", Z. 703 „… er
klassifiziert jede Anforderung …") stehen **nach** der imperativen Vorgehens-Zeile
und dem lesersituativen Opener — es sind Ergebnis-/Mechanik-Erklärungen, keine
produkt-innensichtigen Sektions-Öffnungen.

## Was bewusst Referenz bleibt (kein Mangel — Standard §3 erlaubt §5/§6 als Referenz)

- **§5 vollständig** (`trace`/`coverage`/`modality`/`cross-consistency`-Feld-,
  Grammatik- und Fail-closed-Referenz, Z. 1263–1464) und **§6** (Modul-Tabelle,
  Z. 1466–1488) — Standard §3 Punkt 5 „Konfiguration". Hierhin ist die §4-Referenz
  bewusst gewandert.
- **§4.9 „Die drei Ausgaben im Vergleich"** (Z. 473–479) — Entscheidungs-Hilfstabelle
  in der Aufgabe (Slice-DoD explizit als Referenz gebilligt).
- **§4.10 „Was `--repair` behebt"** (Z. 502–508) — Reparatur-Entscheidungstabelle
  in der Aufgabe (Slice-DoD explizit gebilligt).
- **§4.15-Kreuzverweis-Aufgabe (`trace.cross-consistency`, Z. 724–798)** — die
  slice-071-Schablone (Ausgangslage→Ziel→Schritt 1/2/3→Ergebnis-Deutung→Hinweis),
  vom Slice §2 ausdrücklich als **Vorbild** benannt; sie ist das Modell, an dem die
  übrige Auftrennung ausgerichtet wurde.
- **§4.16 Target-/`--print-mk`-Aufzählung** (Z. 813–856) — in-Task-Referenz auf die
  emittierten Makefile-Targets. Legitim.
- **N-6 §4.9/§4.11-`--json`/`--yaml`-Überlappung** — in-Task-Referenz, vom
  Ursprungs-Audit als „akzeptabel" eingestuft.
- **N-5: keine §4-Aufgabe für `citations`/`sources`** — bewusste Zurückhaltung;
  beide sind in §5 task-betitelt dokumentiert (Z. 1022 / Z. 1158).
- **Kurze YAML-„so erreichen Sie X"-Beispiele** in §4.7/§4.13/§4.14/§4.15 —
  aufgabengerecht (das Beispiel selbst, nicht feldweise Prosa).

## Nebenbefund zur Doku-Currency (bereits erledigt)

Die drei vom Cluster-Review als F-1 (LOW) markierten stale §4.12-Changelog-Verweise
sind nachgezogen: die §11-Zeilen 1.35/1.36/1.37 (Z. 1630–1632) verweisen für die
nach §5 verlagerten Grammatik-Sonderfälle jetzt korrekt auf **§5**. Kein offener
Rest.

## Verdikt

**Gegenprobe BESTANDEN.** Kein offenes HIGH, kein offenes MEDIUM. Von den 14
Ursprungsbefunden sind 12 belegbar ERLEDIGT und 2 (N-5, N-6) legitime,
ausdrücklich benannte In-Referenz-/Ermessens-Entscheidungen. Die zwei neuen
Findings sind INFO (NEU-1 = ein referenz-förmiger, aber sachlich nötiger
Syntax-Rahmen vor dem §4.12-Vorgehen; NEU-2 = eine verstreute Compat-Parenthese in
§5) — beide verzögern höchstens, führen keinen Leser fehl und blockieren den
Slice-Abschluss nicht. Die vom Slice-DoD als „bewusst Referenz" gelisteten Punkte
(§5/§6, §4.9-Vergleichstabelle, §4.10 `--repair`-Tabelle, §4.15-Kreuzverweis-Vorbild)
sind oben ausdrücklich als solche bestätigt.
