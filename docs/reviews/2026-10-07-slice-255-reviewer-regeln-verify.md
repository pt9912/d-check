# Verifikation slice-255 — Reviewer-Regeln aus `v6.17.0` in die Reviewer-Skills (DoD)

- **Rolle:** Verifier (Modul 11), Frage „Bauen wir es richtig?". Geprüft wurde gegen §2 DoD und §1 Ziel und Abgrenzung des Slice-Plans, samt der dort vermerkten Plan-Änderung nach R1. Review-Entscheidungen waren nicht der Maßstab.
- **Gegenstand:** `72150573..79173752`, drei Commits: `badcbb35` (feat), `1ff4e5d7` (R1-Report), `79173752` (R1-Einarbeitung F-1 bis F-4 und F-7).
- **Kanon:** `v6.17.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill; `v6.17.0` · `templates/.harness/skills/reviewer.template.md`, `…/closure-note-reviewer.template.md`, `templates/docs/reviews/review-report.template.md`; `MR-074` Bewegung 3.
- **Sensor-Evidence:** Alles selbst gemessen: `make gates` auf HEAD `79173752` mit sauberem Baum, ein Bruch-Test an den Cite-Spannen, ein whitespace-normalisierter Wortlaut-Abgleich Skill gegen Vorlage (`tr -s` plus `grep -F`) und eine Spiegel-Suche per `grep`. Aus dem Implementer-Bericht ist nichts übernommen.
- **Modell-ID:** claude-opus-5-5 (Claude Agent SDK)
- **Datum:** 2026-10-07

---

## DoD-Prüfung je Punkt (§2)

### 1. `.harness/skills/reviewer.md`: LOW-Anker, Failure-Szenario nur für HIGH/MEDIUM, `pfad` als Kurzzitat, Version/Datum, wortnah, `d-check:cite` bei wörtlichem Zitat: **TEILWEISE** (V-1)

- **LOW-Anker:** Der Skill sagt „**LOW** (nice-to-fix) — *mit Konventions-Anker* (ADR, Hard Rule, Linter-Regel, Eintrag in diesem Skill; ohne Anker kein Finding)". Kanon: `reviewer.template.md` · „**LOW** — *mit Konventions-Anker* (ADR, Hard Rule, Linter-Regel, Eintrag im Reviewer-Skill)". Die Anker-Sorten sind identisch, die Formulierung ist wortnah. **Erfüllt.**
- **Failure-Szenario nur HIGH/MEDIUM:** Der Skill sagt „**Kein HIGH- oder MEDIUM-Finding ohne Failure-Szenario:** was sich nicht als konkretes Versagen erzählen lässt, wird nicht als HIGH oder MEDIUM gemeldet. LOW trägt stattdessen einen Konventions-Anker." Kanon: `reviewer.template.md` · „wird nicht als HIGH oder MEDIUM gemeldet." und `modul-10` · „Kein HIGH- oder MEDIUM-Finding ohne Failure-Szenario". Nach Whitespace-Normalisierung sind die Sätze 1–2 ein wörtlicher Teilstring der Vorlage. Der Zusatz aus `badcbb35`, „INFO braucht keins von beidem" (R1 F-3, über die Baseline hinaus), ist entfernt. **Erfüllt.**
- **`pfad` als Kurzzitat:** Das Output-Schema sagt „Datei · wörtliches, in der Datei eindeutig auffindbares Kurzzitat der Stelle als Anker; die Zeile darf als Lesehilfe dazu, ist aber nicht der Anker". Das ist ein wörtlicher Teilstring von `reviewer.template.md` §Output-Schema. **Erfüllt.**
- **„Kein Stil-Polizist"**, das vierte Element von Bewegung 3: Der Bestand „Formatierung/Benennung ohne Konventions-Anker ist kein Finding" deckt sich inhaltlich mit `modul-10` · „Formatierung oder Benennung ohne Konventions-Anker ist kein Finding". Der Plan vermerkt das jetzt in §1 (R1 F-5).
- **Version/Datum:** Kopf `1.16.0 · 2026-09-07` wurde zu `1.17.0 · 2026-10-07`. **Erfüllt.**
- **`d-check:cite`, wo wörtlich zitiert wird:** **nicht erfüllt**, siehe V-1. Die drei übernommenen Stellen sind gemessen wortgleich zur Vorlage, tragen aber keine Direktive. Bestehende Direktiven im selben Dokument decken genau diese Form ab, nämlich eine wortgleiche Übernahme als Fließtext ohne Anführungszeichen (`reviewer.md` · „Eine „geprüft, ohne Befund"-Zeile pro betrachtetem Bereich", cite auf `modul-10:87`).

### 2. `.harness/skills/closure-note-reviewer.md`: `pfad` als Kurzzitat, Version gehoben: **ERFÜLLT** (Cite-Frage wie V-1)

- Der Skill sagt „`docs/plan/planning/done/<slice>.md` · wörtliches Kurzzitat der Stelle als Anker; die Zeile darf als Lesehilfe dazu, ist aber nicht der Anker". Das ist ein wörtlicher Teilstring von `closure-note-reviewer.template.md` §Output-Schema. Laut `MR-074` ist es das einzige Delta der Vorlage, Bewegung 3 ist damit vollständig getragen.
- Version `1.0.0` wurde zu `1.1.0`, das Datum steht auf 2026-10-07.
- Die Failure-Szenario-Pflicht bleibt hier für alle Kategorien („**Kein Finding ohne Failure-Szenario.**"). Das ist Repo-Bestand: Die Vorlage führt die Regel nicht, und der Plan verlangt hier nur `pfad`. Die LOW-Zeile des Closure-Skills („alle drei Inhalte da, aber schwer nachvollziehbar formuliert") steht damit nicht im Widerspruch. Es gibt keinen Befund (R1 F-8, bewusste Designnotiz).

### 3. Spiegel geprüft, `make gates` grün: **ERFÜLLT**

- **Spiegel:** `.claude/agents/reviewer.md` · „achtzehn Prüffragen, Output-Schema, Negativbefund-Pflicht". Die Prüffragen-Tabelle hat gezählt 18 Zeilen (#1–#18), der Skill-Titel heißt „Die achtzehn Prüffragen", und der Absatz unter der Tabelle sagt ebenfalls „achtzehn". Alle drei Stellen sind konsistent (R1 F-7). `grep -n "Prüffragen\|pfad\|Failure"` über `.claude/agents/*.md` und `.claude/commands/*.md` liefert nur diese eine Zeile. Kein Agent-Prompt und kein Command beschreibt `pfad` als `Datei:Zeile`. Repo-weit stehen die Treffer für `Datei:Zeile` außerhalb von Baseline, Reviews und `done/` im Produkt-Befundformat (Lastenheft, Spezifikation, Handbuch, CHANGELOG, ADR-0039) und in `.d-check.yml:630`. Keiner davon beschreibt das Review-`pfad`-Feld. `harness/README.md` §Guides nennt die Skills nur und braucht keinen Nachzug.
- **`make gates`** (selbst gefahren, HEAD `79173752`, sauberer Baum): `EXIT=0`, Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`. Die Einzelbelege lauten:
  - `verify ok (54 Dateien, vollständig)`
  - `d-check: 961 Datei(en) geprüft, 0 Befund(e)`
  - `0 issues.`
  - `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`
  - `Ran 55 rules on 65 files: 0 findings.`
- **Cite-Spannen lösen auf, mit Bruch-Test:** Die vier lebenden Direktiven der beiden Skills sind im `doc-check` grün:
  - `reviewer.md:59` → `AGENTS.md:286`
  - `reviewer.md:234` → `modul-10:87`
  - `closure-note-reviewer.md:13` → `modul-11:83`
  - `closure-note-reviewer.md:86` → `closure-note-reviewer.template.md:84`

  Für den Bruch-Test wurde `reviewer.md:234` testweise auf `:86-86` gesetzt und `closure-note-reviewer.md:86` auf `:83-83`. `make doc-check` meldete dann `2 Befund(e)`, je `citation-mismatch … Zitattext ist kein zusammenhängender Teilstring der Quell-Spanne (Zitat-Fäule)`, und brach mit `make: *** [Makefile:139: doc-check] Fehler 1` ab. Der Sensor wird also aus dem richtigen Grund rot. Beide Dateien wurden danach per `git checkout` zurückgesetzt, `git status` ist sauber. Das grüne `0 Befund(e)` trägt die Aussage, dass die Spannen auflösen.

### 4. Review durchgeführt, Report liegt vor: **ERFÜLLT**

- `docs/reviews/2026-10-07-slice-255-reviewer-regeln-r1.md` (`1ff4e5d7`, eigener Kontext). Verdikt: „Nicht merge-reif ohne Klärung von F-1 und F-2".
- **Die R1-Einarbeitung (`79173752`) wurde gegen den Report geprüft.** Eine R2 gibt es nicht. Deshalb ist die Einarbeitung hier Stelle für Stelle bestätigt:
  - **F-1, eingearbeitet:** Der Skill sagt jetzt „steigt eine Stufe — aber nur, wenn die höhere Stufe ihre Bedingung erfüllt: nach MEDIUM nur mit erzählbarem Failure-Szenario, von INFO nach LOW nur mit Konventions-Anker; sonst bleibt die Stufe." Beide Konfliktpfade aus F-1 sind geschlossen. MEDIUM nach HIGH braucht keine neue Bedingung, weil MEDIUM das Szenario schon trägt.
  - **F-2, eingearbeitet:** `quelle` führt jetzt auch „Linter-Regel, Abschnitt dieses Skills" sowie „„Maintainability" — letzteres ist **kein** Konventions-Anker und trägt keinen LOW". Damit sind alle vier LOW-Anker-Sorten im Ausgabe-Schema darstellbar. Ein kleiner Rest ist in V-2 beschrieben.
  - **F-3, eingearbeitet:** Der Zusatzsatz zu INFO ist entfernt (siehe Punkt 1).
  - **F-4, eingearbeitet:** §Ablage ergänzt um „Das Kurzzitat aus dem Output-Schema kommt dazu: `` `v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt> · „<Kurzzitat>" ``". Die beiden `pfad`-Definitionen sind damit deckungsgleich.
  - **F-5:** im Plan vermerkt. **F-6:** an den Verifier verwiesen und hier als V-1 behandelt. **F-7:** eingearbeitet (siehe Punkt 3). **F-8:** keine Aktion erwartet.

### 5. Closure-Notiz, Register, Risiko-Ausgänge, Paarungen, MR-074-Vermerk: **OFFEN, wie erwartet**, kein Befund

§6 steht auf `(offen)`, §7 trägt `—`. Das ist der Closure-Schritt nach der Verifikation. `make verify-closure-notes` ist nicht anwendbar, solange der Slice nicht in `done/` liegt. `MR-074` ist im Range unverändert, der Einlösungs-Vermerk kommt mit der Closure (Plan §3).

---

## Abgrenzung (§1): gehalten

- Der Range berührt fünf Dateien. Neben den beiden Skills und `.claude/agents/reviewer.md` sind das der Plan und der neue R1-Report. Kein bestehender Report wurde umgeformt: `git diff --stat 72150573..HEAD -- docs/reviews` zeigt nur den neuen R1-Report.
- Es gibt keinen Sensor auf die Report-Form, keine Änderung an Spezifikation, Werkzeug-Teil oder Register-Kennung, also keine Bewegung außer 3. Kein Gate ist berührt (`AGENTS.md` §3.6).
- Die Plan-Änderung nach R1 zieht Kontext-Eskalation, `quelle`, die zweite `pfad`-Definition und das Zählwort mit. Sie ist in §1 vermerkt und vom tatsächlichen Diff gedeckt. Der Diff geht nicht darüber hinaus: Die Prüffragen-Tabelle selbst ist unverändert.
- Die Rückführungs-Bedingung §4 („Umstellung des Prüffragen-Katalogs über die drei Stellen hinaus") ist nicht eingetreten.

---

## Verifier-Befunde

### V-1: MEDIUM

- **Kategorie:** MEDIUM (DoD-Verletzung)
- **Quelle:** DoD `slice-255` Punkt 1 („mit `d-check:cite`, wo wörtlich zitiert wird"); Plan §8 (Sichtung `BEO-ALL/citation-stretched-beyond-scope`: „wörtliche Übernahmen aus der Baseline werden als `d-check:cite` geführt"); `MR-051`
- **Pfad:**
  - `.harness/skills/reviewer.md` · „wird nicht als HIGH oder MEDIUM gemeldet."
  - `.harness/skills/reviewer.md` · „Datei · wörtliches, in der Datei eindeutig auffindbares Kurzzitat der Stelle"
  - `.harness/skills/closure-note-reviewer.md` · „wörtliches Kurzzitat der Stelle als Anker; die Zeile darf als Lesehilfe dazu"
- **Befund:** Alle drei Übernahmen sind gemessen wortgleich zur Vorlage, als whitespace-normalisierter Teilstring von `reviewer.template.md:87-88` bzw. `:99-100` und `closure-note-reviewer.template.md:77-78`. Keine trägt eine `d-check:cite`-Direktive. Der Plan selbst hält in §8 fest, dass wörtliche Übernahmen als `cite` geführt werden. Im selben Skill steht dafür eine Präzedenz: die wortgleiche Fließtext-Übernahme bei §Negativbefunde mit cite auf `modul-10:87`. Der DoD-Punkt lässt sich in dieser Form nicht ehrlich abhaken.
- **Failure-Szenario:** Ein künftiger Baseline-Bump ändert den Wortlaut einer dieser Regeln. Der Bump-Lauf ankert nach `MR-051` die cite-Spannen neu und findet für diese drei Stellen keine. `citations` meldet nichts, weil keine Direktive da ist. Der Skill behauptet dann weiter, wortnah zur Baseline zu sein, ohne es zu sein. Das ist die Zitat-Fäule, die das Modul fangen soll, nur still.
- **Verifizierbar:** ja. Ist eine Direktive gesetzt, prüft `make doc-check` (`citations`) sie. Der Bruch-Test oben zeigt, dass das Modul auf den Skills greift.
- **Ausweg (eines von beiden):** Entweder die Direktiven setzen. Bei der Bullet-Form muss das zitierte Stück dafür als eigener Absatz oder an einer Zeile stehen, die nur den Vorlagen-Teilstring trägt. Oder den DoD-Wortlaut per vermerkter Plan-Änderung mit Begründung auf „wortnah, ohne cite" einengen.
- **Klasse:** DoD-Klausel ohne Lieferung (Cite-Pflicht bei wortgleicher Übernahme)

### V-2: LOW

- **Kategorie:** LOW
- **Quelle:** Eintrag in diesem Skill (LOW-Zeile und Output-Schema `quelle`); Rest von R1 F-2
- **Pfad:**
  - `.harness/skills/reviewer.md` · „Eintrag in diesem Skill; ohne Anker kein Finding"
  - `.harness/skills/reviewer.md` · „letzteres ist **kein** Konventions-Anker und trägt keinen LOW"
- **Befund:** Die LOW-Zeile zählt vier Anker-Sorten abschließend auf: ADR, Hard Rule, Linter-Regel und Eintrag im Skill. Das `quelle`-Feld führt zusätzlich `DC-*`- und `MR-*`-IDs. Weil es „Maintainability" ausdrücklich als einzigen Nicht-Anker markiert, legt es nahe, dass `DC-*` und `MR-*` als Anker taugen. Ob ein LOW, der nur auf einem `MR-*`-Eintrag ruht, zulässig ist, beantwortet der Skill an zwei Stellen verschieden.
- **Failure-Szenario:** Nicht nötig für LOW, der Anker oben genügt. Folge wäre, dass zwei Reviewer denselben LOW mit `quelle` `MR-0xx` verschieden behandeln.
- **Verifizierbar:** nein, das ist ein Urteil über Skill-Kohärenz.
- **Klasse:** abschließende Aufzählung gegen offenes Ausgabe-Feld

### V-3: INFO

- **Kategorie:** INFO
- **Quelle:** `AGENTS.md` §6 Schritt 4 („gehört vor den Code, nicht in den Bericht danach")
- **Pfad:** `docs/plan/planning/in-progress/slice-255-reviewer-regeln-v6170.md` · „*(Plan-Änderung nach R1:"
- **Befund:** Die Plan-Änderung steht im Plan und nicht nur in der Botschaft. Sie kam aber im selben Commit `79173752` wie die Skill-Änderung, die sie deckt, und nicht davor. Inhaltlich ist sie vollständig und entspricht dem Diff. Dazu kommt: Die R1-Einarbeitung der zwei MEDIUMs hat keine eigene R2-Runde. Ihre Bestätigung trägt dieser Bericht (Punkt 4). Erwartete Aktion: keine. Rollen-Verweis: Planner, für die Closure-Notiz.
- **Klasse:** Plan-Änderung im Code-Commit

Negativ geprüft ohne Befund:
- Kontext-Eskalation gegen Anti-Pattern und LOW-Zeile
- `pfad` in Output-Schema gegen §Ablage
- Stil-Polizist gegen LOW und INFO
- die MEDIUM-Zeilen der Prüffragen-Tabelle gegen die Failure-Szenario-Pflicht
- Closure-Skill-Bestand (Failure-Szenario-Pflicht, LOW-Zeile) gegen den Planumfang
- Zählwort „achtzehn" an allen drei Stellen
- Versionen und Daten
- Spiegel in `.claude/agents/` und `.claude/commands/`
- vier lebende cite-Spannen samt Bruch-Test
- Abgrenzung (keine Report-Umformung, kein Sensor, keine andere Bewegung, §3.6 unberührt)

---

## Verdict

- **DoD 1 teilweise erfüllt:** Die drei Regeln stehen korrekt und wortnah. Die cite-Klausel ist nicht eingelöst (V-1).
- **DoD 2–4 erfüllt.**
- **DoD 5 erwartungsgemäß offen.**

`make gates` wurde selbst gefahren und ist grün. Die cite-Spannen lösen auf, und der Bruch-Test zeigt, dass `citations` an den Skills aus dem richtigen Grund rot wird. Die R1-MEDIUMs F-1 und F-2 sind in der Sache geschlossen. Die drei Regeln stehen in beiden Skills und im Agent-Spiegel ohne Widerspruch zu Kontext-Eskalation, `quelle`-Feld, Output-Schema und §Ablage, bis auf den kleinen Rest in V-2.

**Vor dem Abhaken von DoD 1 ist V-1 aufzulösen:** Entweder werden die Direktiven gesetzt, oder der DoD-Wortlaut wird mit Begründung eingeengt. V-2 bitte annehmen oder begründen. V-3 braucht keine Aktion.
