# Review R1 — slice-255: Reviewer-Regeln aus `v6.17.0` in die Reviewer-Skills

- **Review-Art:** Code (Gegenstand ist Harness-Prosa). Geprüft wurde gegen den Slice-Plan
  (`slice-255`, §1 Ziel und Abgrenzung samt Begründung der Lockerung, §3 Spiegel-Liste,
  §6 Risiken), gegen `MR-074` Bewegung 3, gegen die Baseline `v6.17.0` ·
  `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill und die Vorlagen
  `v6.17.0` · `templates/.harness/skills/reviewer.template.md`,
  `…/closure-note-reviewer.template.md`, `templates/docs/reviews/review-report.template.md`
  (je im Delta gegen `v6.13.0`, gelesen aus `957aeedc^`), gegen `MR-025` und die Hard
  Rules `AGENTS.md` §3.6/§3.7 sowie §5 Regel 13/15. Die DoD-Abhakung prüft dieses Review
  nicht.
- **Gegenstand:** `slice-255` · Commit `badcbb35` (`.harness/skills/reviewer.md`
  1.16.0 → 1.17.0, `.harness/skills/closure-note-reviewer.md` 1.0.0 → 1.1.0).
- **Skill:** `reviewer.md` @ 1.16.0 (Fassung `badcbb35^` — der Gegenstand ist der Skill
  selbst; die neue Fassung 1.17.0 ist Prüfgegenstand, nicht Maßstab). Report-Form
  (`pfad` als Kurzzitat) folgt auf Auftrag bereits der neuen Regel.
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-07
- **Eingangs-Kontext:** Slice-Plan `slice-255`; `MR-074` (Bewegung 3); Baseline-Delta
  `v6.13.0` → `v6.17.0` für `modul-10`, Reviewer-Vorlage, Closure-Note-Vorlage,
  Report-Vorlage; Herkunft der bisherigen Failure-Szenario-Regel (`git log -S`: eingeführt
  mit `853bcfe1`, ohne Herkunfts-Anker — kein Retirement-Check fällig). Vorherige Findings
  am selben Gegenstand: R1/R2 zu `slice-235` (Rollen-Agents; Klasse
  „zu allgemeiner Baseline-Anker").
- **Proben:** `make gates` auf dem Stand `badcbb35` gefahren — grün
  (`[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check +
  coverage-gate + semgrep + gate-consistency + planning-check green`). Spiegel-Suche per
  `grep -rn "Datei:Zeile\|<Zeile>\|Kein Finding ohne\|ohne Failure\|nice-to-fix"` außerhalb
  des vendorten Baums.

## Findings

| # | kategorie | befund | quelle | pfad | verifizierbar | klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die unveränderte **Kontext-Eskalation** hebt „dieselbe Beobachtung im Gate-/Sicherheitspfad" eine Stufe — seit der Lockerung macht sie aus einem LOW, der legitim **ohne** Failure-Szenario gemeldet wird (nur Konventions-Anker), ein MEDIUM, das das neue Anti-Pattern ohne Failure-Szenario verbietet; ebenso aus einem anker-losen INFO ein LOW ohne Anker. Failure-Szenario: Doku-Drift in einem Gate-Skript — Reviewer A eskaliert per §Kontext-Eskalation auf MEDIUM (blockierend), Reviewer B lässt sie per Anti-Pattern auf LOW; gleiche Eingabe, anderes Verdikt — genau der Dissens, den der Skill laut Baseline verhindern soll. Vorher bestand der Konflikt nicht, weil *jedes* Finding ein Szenario brauchte. | `v6.17.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill (Drift „gleiche Eingabe → andere Findings/Kategorien") | `.harness/skills/reviewer.md` · „dieselbe Beobachtung im Gate-/Sicherheitspfad steigt eine Stufe" (Z. 197–198, Lesehilfe) | nein — Urteil über Skill-Kohärenz, kein Gate | Lockerung ohne Nachzug der gekoppelten Regel |
| F-2 | MEDIUM | Die neue LOW-Bedingung verlangt einen Anker aus vier Sorten (ADR, Hard Rule, Linter-Regel, Eintrag im Skill), das Output-Schema bietet für `quelle` aber nur `DC-*`, ADR, `MR-*`, Hard-Rule-Name oder „Maintainability" — zwei der vier Anker-Sorten haben dort keinen Wert, und „Maintainability" ist gerade die anker-lose Angabe. Failure-Szenario: ein LOW mit `quelle` „Maintainability" ist schema-konform, und ob er einen Anker hat, lässt sich aus dem Report nicht mehr ablesen — die Regel „ohne Anker kein Finding" ist an ihrem einzigen Ausgang unprüfbar. Zur Auftragsfrage: die drei Bestands-Beispiele (Doku-Drift, latente Wartungsfalle, Ketten-Duplikate) haben keinen ADR-/Hard-Rule-/Linter-Anker; sie sind nur dadurch geankert, dass sie selbst die „Einträge in diesem Skill" sind — wo das im Report steht, sagt der Skill nicht. | Maintainability (Output-Schema ↔ LOW-Zeile desselben Skills); `v6.17.0` · `templates/.harness/skills/reviewer.template.md` §Klassifikation | `.harness/skills/reviewer.md` · „`MR-*`-ID, Hard-Rule-Name oder „Maintainability")" (Z. 220, Lesehilfe) | nein | neue Pflicht ohne Feld im Ausgabe-Schema |
| F-3 | LOW | Der Zusatzsatz „INFO braucht keins von beidem" steht weder in `modul-10` noch in der Vorlage — die Übernahme geht hier über die Baseline hinaus, obwohl Plan und Botschaft „folgt der Baseline-Regel" sagen. Er widerspricht zwei Zeilen höher „Kein Stil-Polizist: Formatierung/Benennung ohne Konventions-Anker ist kein Finding" (ein Stil-INFO ohne Anker wäre nach dem einen erlaubt, nach dem anderen kein Finding) und lässt die Baseline-Bedingung für das außer-Auftrag-INFO weg („INFO-Finding mit Rollen-Verweis"). | Eintrag in diesem Skill (Anti-Pattern „Kein Stil-Polizist"); `v6.17.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill | `.harness/skills/reviewer.md` · „LOW trägt stattdessen einen Konventions-Anker, INFO braucht keins von beidem." (Z. 210) | nein | Übernahme weiter als die Quelle |
| F-4 | LOW | Der Skill definiert `pfad` jetzt zweimal verschieden: das Output-Schema als „Datei · wörtliches … Kurzzitat der Stelle als Anker", der unveränderte Absatz unter §Ablage/Zitier-Form „präzisiert" `pfad` für Baseline-Fundstellen als `` `v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt> `` — ohne Kurzzitat, mit dem Abschnitt als Anker. Ein `§Ziel-Form`-Abschnitt ist dutzende Zeilen lang; ein `pfad` in der Zitier-Form ist damit genau die nicht eindeutig auffindbare Fundstelle, die die neue Regel abschaffen soll. | Eintrag in diesem Skill (Output-Schema `pfad`); `v6.17.0` · `templates/docs/reviews/review-report.template.md` (Spalte „Datei · wörtliches Kurzzitat, Zeile optional") | `.harness/skills/reviewer.md` · „Das `pfad`-Feld eines Findings ist davon **nicht ausgenommen, sondern präzisiert**" (Z. 256) | nein | Spiegel im selben Dokument nicht nachgezogen |
| F-5 | INFO | Plan §1 nennt „drei Reviewer-Regeln"; `MR-074` Bewegung 3 führt vier Elemente, darunter „Kein Stil-Polizist". Der Skill trägt es schon (Bestand, Z. 203–204) — der Plan sagt nicht, dass das vierte Element bereits erfüllt ist. Rollen-Verweis: Planner/Closure (Einlösungs-Vermerk in `MR-074`). | `MR-074` (Bewegung 3) | `docs/plan/planning/in-progress/slice-255-reviewer-regeln-v6170.md` · „Die Reviewer-Skills folgen den drei Reviewer-Regeln der Baseline" | nein | Inventar-Element ohne Vermerk |
| F-6 | INFO | Die übernommenen Sätze (Failure-Szenario, `pfad`-Definition in beiden Skills) sind wortgleich zur Vorlage, tragen aber keine `d-check:cite`-Direktive; ob eine übernommene Regel ein „wörtliches Zitat" im Sinn des DoD-Punkts ist, ist eine DoD-Frage. Rollen-Verweis: Verifier. | `MR-052`/DoD `slice-255` | `.harness/skills/reviewer.md` · „Datei · wörtliches, in der Datei eindeutig auffindbares Kurzzitat der Stelle" | ja — `make doc-check` mit `citations` würde eine gesetzte Direktive prüfen | — |
| F-7 | INFO | Außerhalb des Diffs, Bestand: die Überschrift heißt „Die sechzehn Prüffragen", die Tabelle trägt 18 Zeilen; der Spiegel `.claude/agents/reviewer.md` sagt ebenfalls „sechzehn Prüffragen". Rollen-Verweis: Planner (Folge-Pflege). | Maintainability | `.harness/skills/reviewer.md` · „## Die sechzehn Prüffragen (erste Ebene)" | nein | Zählwort driftet gegen Tabelle |
| F-8 | INFO | Der Closure-Note-Skill behält „Kein Finding ohne Failure-Szenario" für alle Kategorien, der Schwester-Skill nicht. Die Botschaft begründet das („weil die Vorlage sie nicht aendert" — genauer: die Vorlage führt die Regel gar nicht, sie ist Repo-Bestand); der Plan, der die Lockerung begründet, trägt die Gegen-Entscheidung nicht. Bewusste Designnotiz, kein Mangel. | Maintainability | `.harness/skills/closure-note-reviewer.md` · „**Kein Finding ohne Failure-Szenario.**" | nein | Entscheidung nur in der Botschaft |

## Negativbefunde

- geprüft, ohne Befund: Wortlaut der Failure-Szenario-Regel gegen `modul-10`/Vorlage —
  Satz 1–2 wortgleich, die Lockerung selbst reicht genau bis zur Baseline (Überschuss nur
  der Zusatzsatz, F-3).
- geprüft, ohne Befund: LOW-Bedingung gegen die Vorlage — Anker-Sorten identisch
  („Eintrag im Reviewer-Skill" → „Eintrag in diesem Skill").
- geprüft, ohne Befund: `pfad`-Wortlaut in beiden Skills gegen die jeweilige Vorlage
  (Closure-Note: einziges Delta der Vorlage, vollständig übernommen).
- geprüft, ohne Befund: Vollständigkeit Bewegung 3 — die Report-Vorlage hat im Repo keine
  lokale Kopie (§Ablage verweist auf die Baseline-Ziel-Form), es fehlt kein Träger;
  Bewegung 4 (Register-Kennung) ist nicht mitgenommen, wie geplant.
- geprüft, ohne Befund: Abgrenzung — zwei Dateien, keine Report-Umformung, kein Sensor;
  §3.6 nicht berührt (kein Gate), §3.7 nicht berührt (keine Code-/Config-Kommentare).
- geprüft, ohne Befund: Retirement-Frage der gelockerten Regel — eingeführt mit `853bcfe1`,
  kein Herkunfts-Anker, keine Steering-Loop-Herkunft.
- geprüft, ohne Befund: Spiegel (`MR-025`) — `.claude/agents/*.md` und
  `.claude/commands/*.md` beschreiben `pfad` nicht; `Datei:Zeile` im Benutzerhandbuch
  (Z. 134, 2314, 2759), Lastenheft, Spezifikation und Code betrifft das Befund-Format des
  Produkts; der Kommentar `.d-check.yml` · „Review-Reports zitieren naturgemäß
  Datei:Zeile/Pfade" bleibt sachlich richtig (Reports nennen weiter Dateien).
- geprüft, ohne Befund: Commit-Botschaft gegen Diff (Regel 15) — jede Aussage durch den
  Diff bzw. die Spiegel-Suche gedeckt; „make gates gruen" in diesem Lauf nachgefahren, grün.
- geprüft, ohne Befund: Plan-Behauptungen (Regel 13) — „strenger als die Baseline" trifft
  gegen `v6.17.0` zu (`v6.13.0` führte keine Failure-Szenario-Regel); „nicht gemessen"
  ist ehrlich ausgewiesen.
- geprüft, ohne Befund: Version/Datum — Minor-Hebung je Skill bei Regel-Änderung, Datum
  2026-10-07; bestehende `d-check:cite`-Direktiven unberührt und grün (`doc-check`).

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 2 | 2 | 4 |

## Verdikt

Nicht merge-reif ohne Klärung von F-1 und F-2: beide entstehen erst durch die Lockerung
und machen den Skill an genau der neuen Grenze uneindeutig (F-1: Verdikt hängt davon ab,
welche von zwei Regeln der Reviewer liest; F-2: die neue Pflicht hat kein Feld im
Ausgang). F-3 und F-4 sind nice-to-fix im selben Zug; F-5 bis F-8 ohne erwartete Aktion
in diesem Slice.
