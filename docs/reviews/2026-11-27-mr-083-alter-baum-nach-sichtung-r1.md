# Review-Report: MR-083, der alte Baseline-Baum fällt nach der Sichtung (Runde 1)

- **Review-Art:** Design (neuer Adaptions-Eintrag gegen Baseline-Kanon, lebende Adaptionen und Hard Rules)
- **Gegenstand:** Commit `76bfdae4` — neuer Eintrag MR-083, eine Index-Zeile in `harness/conventions.md`, Grenze 5 in der Sensor-Datei zu `make baseline-verify`
- **Skill:** `reviewer.md` @ 1.20.0 (Stand `76bfdae4`)
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-11-27
- **Eingangs-Kontext:** `v7.0.0` · `regelwerk/modul-02-harness-bootstrap.md` §Freshness-Audit; `v7.0.0` ·
  `regelwerk/modul-06-roadmap.md` §Wann Arbeit eine Welle braucht; `v7.0.0` ·
  `templates/harness/conventions/MR-NNN-title.template.md`; MR-021, MR-023, MR-051, MR-052, MR-055, MR-069, MR-081;
  `AGENTS.md` §3; `SPEC-092`; wave-98 §3 und §5; `BEO-ALL/closure-trigger-misst-was-ein-spaeterer-slice-veraendert`.
  Kein Slice-Plan: der Commit läuft ohne Slice. DoD nicht gelesen.

## Findings

### F-1 — Der angegebene Hand-Befehl läuft nicht

- **kategorie:** MEDIUM (Grundstufe LOW, Doku-Drift; eine Stufe angehoben nach Skill §Kontext-Eskalation, weil die
  Stelle die Grenze eines Gates beschreibt)
- **quelle:** Reviewer-Skill, Anker LOW *Doku-Drift* und *Kontext-Eskalation*
- **pfad:** `harness/sensors/baseline-verify.md` · „`tools/harness/fetch-baseline-cache.sh --verify <alt-tag>`. Permanent.";
  ebenso MR-083, Absatz *Grenzen* · „den alten Baum prüft von Hand"
- **befund:** Das Skript ist ohne Ausführungsrecht getrackt (`git ls-files -s` zeigt `100644`); der Aufruf in der
  genannten Form endet mit „Keine Berechtigung". Er läuft nur mit vorangestelltem `bash`, wie ihn das `Makefile` ruft.
- **Failure-Szenario:** Bei der nächsten Hebung will jemand den liegenden alten Baum prüfen — der Befehl ist der
  einzige dafür genannte Ersatz, weil `make baseline-verify` ihn nicht prüft — und erhält statt eines Urteils einen
  Shell-Fehler. Der Ausfall ist laut, nicht still.
- **verifizierbar:** ja, durch den Aufruf selbst; kein Gate meldet ihn.
- **klasse:** dokumentierter-befehl-laeuft-nicht

### F-2 — Die einzige Abweichung vom Kanon steht ohne Begründung

- **kategorie:** LOW
- **quelle:** `v7.0.0` · `templates/harness/conventions/MR-NNN-title.template.md` (Pflichtfelder *Adaption* und
  *Begründung*)
- **pfad:** MR-083, Feld *Ersetzt-Baseline-Regel* · „der Gang durch die Adaptions-Liste ist keine Bedingung dafür"
- **befund:** Dass für das Entfernen die Sichtung genügt und der Adaptions-Durchgang keine Bedingung ist, steht nur
  im Feld *Ersetzt-Baseline-Regel*; *Adaption* nennt es nicht, und *Begründung* sagt nicht, warum. Begründet werden
  die Reihenfolge Sichtung, Entfernen, Feststellung und das Entfernen selbst.
- **verifizierbar:** nein (Lektüre)
- **klasse:** abweichung-ohne-begruendung

### F-3 — Der Reviewer-Skill sagt „genau einen Tag"

- **kategorie:** LOW
- **quelle:** Reviewer-Skill, Anker LOW *Doku-Drift*
- **pfad:** `.harness/skills/reviewer.md` §Ablage · „Der vendorte Baum … der nächste Bump löscht den alten":
  der Satz lautet „Der vendorte Baum trägt genau **einen** Tag; der nächste Bump löscht den alten"
- **befund:** Nach MR-083 liegen während einer Hebung zwei Tags nebeneinander, und der Pin-Commit löscht den alten
  nicht. Die Zitier-Regel, die der Satz begründet, bleibt davon unberührt.
- **verifizierbar:** nein
- **klasse:** beschreibung-ueberholt-durch-neue-adaption

### F-4 — Der Auflösungs-Trigger von MR-052 liest sich danach als erfüllt

- **kategorie:** INFO
- **quelle:** MR-052
- **pfad:** MR-052, *Auflösungs-Trigger* · „entfällt, sobald ein Bump den alten Baum nicht mehr"
- **befund:** MR-083 legt „Bump" nur für den zitierten Satz aus MR-021 als *Hebung* aus. MR-052 verwendet dasselbe
  Wort; in der Sache gilt sein Trigger nicht, weil die Hebung den Baum weiterhin entfernt.
- **verifizierbar:** nein
- **klasse:** begriff-nur-fuer-eine-quelle-ausgelegt

### F-5 — „Gesichtet" bei neuen oder entfallenen Verzeichnissen

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** MR-083 · „wenn eine committete Datei jede Datei nennt, die"
- **befund:** `diff -rq` meldet ein nur einseitig vorhandenes Verzeichnis als eine Zeile mit dem Verzeichnisnamen,
  nicht je Datei darin. Für einen solchen Fall deckt die Bedingung die Dateien darunter nicht.
- **verifizierbar:** ja, mit `diff -rq` an zwei Probe-Bäumen
- **klasse:** mengen-definition-ueber-werkzeugausgabe

### F-6 — Die Feststellung ist nicht ans Ende der Hebung gebunden

- **kategorie:** INFO
- **quelle:** `BEO-ALL/closure-trigger-misst-was-ein-spaeterer-slice-veraendert`
- **pfad:** MR-083 · „**erst danach** feststellen, was den alten Pin noch nennt."
- **befund:** Geordnet ist Sichtung vor Feststellung. Läuft in derselben Welle nach dem Entfernen noch ein Slice, der
  den alten Tag in eine lebende Datei schreibt, steht sie nicht in der festgestellten Menge. Der Anlass (die Datei der
  Sichtung selbst) ist abgedeckt, und der Eintrag behauptet nicht mehr.
- **verifizierbar:** nein
- **klasse:** menge-am-zwischenstand-festgestellt

### F-7 — Text, den die Regel nicht braucht

- **kategorie:** INFO
- **quelle:** Maintainability (Maßstab des Auftraggebers: kleinste tragende Regel)
- **pfad:** MR-083 · „**Entfernt werden muss weiterhin.**" bis „ein grüner Lauf den Nachzug also nicht", dazu
  „**Nicht gemessen:**"
- **befund:** Der Mess-Absatz stützt das Entfernen, das MR-021 (1) schon verlangt und das MR-083 nicht ändert; der
  Satz *Nicht gemessen* hängt allein an ihm. Der Satz „Die Hebung von MR-081 bleibt, wie sie lief" wiederholt „ab der
  nächsten". Tragend sind die drei Schritte, *Gesichtet*, *Die Feststellung*, *Abgeschlossen*, der Bezug zu MR-021 und
  die Grenze „kein Gate".
- **verifizierbar:** nein
- **klasse:** adaptions-eintrag-mit-nebenbeleg

## Negativbefunde

- **Zitat aus modul-02, Geltungsbereich** — geprüft, ohne Befund. Der Satz steht im Absatz *Der Review vergleicht auch
  die Form* des Freshness-Audits; „Review" ist dort der Audit des neuen Stands, und der Nachbar-Absatz führt den Gang
  durch die Adaptions-Liste als Teil desselben Reviews. Der Eintrag schreibt dem Satz nichts darüber hinaus zu.
- **MR-021** — geprüft, ohne Befund. (1) und die Link-Pin-Regel bleiben; der verschobene Zeitpunkt und die Auslegung
  von „Pflicht-Teil des Bumps" sind ausdrücklich benannt.
- **MR-051, MR-055, MR-069** — geprüft, ohne Befund. Spannen, Aliase und Ventil-Einträge bleiben fällig; ihr Befund
  erscheint beim Entfernen statt beim Pin-Commit, und der Eintrag sagt, dass ein grüner Lauf bis dahin nichts belegt.
- **MR-081, wave-98 §5 Frage 9** — geprüft, ohne Befund. Beide beschreiben die gelaufene Hebung; MR-083 gilt ab der
  nächsten.
- **Hard Rules (`AGENTS.md` §3)** — geprüft, ohne Befund. Keine Gate-Senkung (§3.6), kein Spec-Stratum berührt (§3.4),
  die neue Sensor-Grenze trägt *Permanent*, kein Code-Kommentar geändert (§3.7).
- **Form der MR-Vorlage** — geprüft, außer F-2 ohne Befund. Pflichtfelder vorhanden, eine benannte Regel mit
  Abschnitts-Anker, Index-Zeile deckt sich mit dem Eintrag.
- **Nächste Hebung ohne Welle** — durchgespielt, ohne Befund. Ein Slice trägt Pin-Commit, committete Sichtung,
  Entfernen, Feststellung und dann den Move nach `done/`; `v7.0.0` · `regelwerk/modul-06-roadmap.md` macht mehr als
  einen Slice mit gemeinsamer Bedingung zur Welle, die dann Träger ist.
- **Nächste Hebung mit Welle** — durchgespielt, ohne Befund. Einzelne Slices dürfen vor dem Entfernen schließen, die
  Welle nicht; `ls .harness/baseline/` ist von jedem Dritten lesbar.
- **Push-Regel** — durchgespielt, ohne Befund. Der Zustand mit zwei Bäumen ist an Links und Aliasen grün (unten
  nachgefahren). Der Entfern-Commit muss Ventil-Einträge und Nachzüge mitbringen, sonst weist `make doc-check` im
  `pre-commit` ihn ab; ein roter Zwischenstand wird von MR-083 nicht verlangt. Die CI ist nicht geprüft.
- **Erfüllbarkeit von „gesichtet"** — stichprobenhaft, ohne Befund. Die Entscheidungsvorlage zu `v7.0.0` (39
  Datei-Zeilen und eine Sammelzeile mit Kurzformen) hat die Form, die die Bedingung verlangt.
- **Mess-Aussagen des Eintrags** — zur Hälfte nachgefahren, ohne Befund (Wegwerf-Klon von `7e024045`, alter Baum aus
  `git archive 738d2bca`, vorhandenes Image `d-check:latest`, danach entfernt): neuer Link in `README.md` auf den alten
  Baum mit beiden Bäumen 0 Befunde, nach dem Entfernen `target-missing`; zurückgehängter Alias mit beiden Bäumen
  `--verify` Exit 0, nach dem Entfernen `make baseline-verify` Exit 2. **Nicht nachgefahren:** `make gates`,
  `make baseline-probe`, `make verify-closure-notes` mit zwei Bäumen.
- **Behauptete Prüfungen** — geprüft, außer F-1 ohne Befund. Die Sensor-Grenze 5 stimmt mit dem Skript überein (es
  liest nur den Tag der `**Stand:**`-Zeile oder das Argument); der Satz über das Vendoring bei liegendem alten Baum
  ist als Lektüre ausgewiesen und trifft den Code; „Kein Gate prüft Reihenfolge, Sichtung oder Entfernen" ist offen
  gesagt. *Abgeschlossen* hat außer dem Eintrag selbst keinen Träger (kein Command, keine Vorlage), behauptet aber
  auch keinen.
- **`SPEC-092`** — geprüft, ohne Befund. Die Festlegung spricht vom Baum des gepinnten Tags; Grenze 5 widerspricht ihr
  nicht.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 |
| LOW | 2 | F-2, F-3 |
| INFO | 4 | F-4, F-5, F-6, F-7 |

## Verdikt

Die Regel trägt: Das Zitat reicht so weit wie die Zuschreibung, keine lebende Adaption und keine Hard Rule wird
widersprochen, und die drei Bedingungen sind für einen Dritten am Repo entscheidbar und in einer Hebung mit wie ohne
Welle erfüllbar. Offen vor dem Abschluss ist allein F-1 (der Hand-Befehl an beiden Stellen). F-2 bis F-7 blockieren
nicht.
