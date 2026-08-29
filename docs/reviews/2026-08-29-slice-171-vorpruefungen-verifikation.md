# Verifikation slice-171 — DoD und Spec-Konformität

**Gegenstand:** [slice-171](../plan/planning/in-progress/slice-171-vorpruefungen-belegen.md), Stand `4ef5025`.
**Datum:** 2026-08-29. **Verifikation:** unabhängiger Subagent, eigener Kontext.
**Methode:** alle Messungen selbst gefahren (Docker/make-only). Jede temporär geänderte Datei per `git checkout --` zurückgesetzt und mit `sha256sum -c` als byte-identisch bestätigt.

---

## DoD-Tabelle

| # | DoD-Punkt | Status | Beleg |
|---|---|---|---|
| 1 | Zwei kanonische Blöcke tragen je eine `d-check:cite`-Direktive; `doc-check` grün; manipuliertes Zitat **gemessen** rot | **erfüllt** | Direktiven `slice-171…md:153` → `modul-05:213-214`, `:167` → `modul-05:224-225`; Quellzeilen wortgleich. `make doc-check` ⇒ 589 Dateien, 0 Befunde, Exit 0. Bruch-Probe 1 (`Schwelle ≥ 2`→`≥ 3`) ⇒ `citation-mismatch`, Exit 2; Bruch-Probe 2 (`werden notiert`→`werden vermerkt`) ⇒ `citation-mismatch`, Exit 2. Beide zurückgesetzt |
| 2 | Dritter Block **ohne** Direktive, Grund als benannter Entscheid | **erfüllt** | Slice §7, Begründung in §2 Punkt 2 und in MR-054 §„Warum der dritte Block keine Direktive trägt"; kein Marker im Nachtlauf-Block |
| 3 | `citations.scope` nimmt `done/` aus; Nachtrag zu MR-051 | **erfüllt**, mit korrigierter DoD-Prämisse | `.d-check.yml`, neuer `citations:`-Block. Gemessen: wirksame Fehl-Direktive am Ende von `done/slice-170-…md` ⇒ 589 Dateien, 0 Befunde, Exit 0; **malformte** Direktive dort ⇒ ebenfalls Exit 0. Beide Pfade sind ausgenommen, nicht nur der Befund-Pfad |
| 4 | Konvention an **einer** Stelle, „und ist von dort verlinkt" | **teilweise erfüllt** | Eine Stelle: ja. **Verlinkung fehlte dort, wo die Regel gilt** — `AGENTS.md` beschreibt die drei Vorprüfungen und verlinkt MR-053, nicht MR-054; `harness/README.md` nennt `citations`/MR-051, nicht MR-054 |
| 5a | `make gates` grün (Exit explizit) | **erfüllt** | Exit 0; `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` |
| 5b | **unabhängiger Review** | **zum Prüfzeitpunkt nicht belegt** | kein Report in `docs/reviews/`, kein Review-Commit — konsistent mit dem Anlassfall, den der Slice selbst beschreibt |
| 5c | Verifikation gegen DoD/Spec in eigenem Kontext | **erfüllt** | dieser Lauf |

## Spec-Konformität

**DC-FA-CONF-002 — konform.** `citations` ist ein scannendes Modul und gehört
nicht zu den vier, die `scope` ablehnen (`planning`, `targets`, `sources`,
`structure`). `roots: ["."]` erfüllt Existenz und Repo-Grenze.

**DC-FA-CONF-002.a — konform, und die tragende Behauptung ist nachgemessen.**
Die beiden `.harness/`-Ausschlüsse testweise aus `citations.scope` entfernt: die
gemeldete Menge stieg **589 → 655** (+66) — exakt der Delta, den der Commit mit
587 → 653 misst. Damit ist sowohl die „ersetzt statt erbt"-Semantik als auch die
Union-Zählung belegt; die Wiederholung der globalen Ausschlüsse ist Pflicht,
nicht Redundanz. Config zurückgesetzt.

**DC-FA-CITE-001 — konform.** Die Spezifikation nennt `citations.scope`
ausdrücklich und legt fest, dass er die **prüfende** Datei skopiert, nicht das
Ziel. Genau darauf ruht die Konstruktion: `.harness/baseline/**` ist
ausgeschlossen und bleibt trotzdem **zitierbar** — der grüne Lauf beweist es.

`**Berührte Spec-Stellen:** —` im Slice-Kopf ist korrekt: keine
Produkt-Anforderung geändert.

## Abweichungen Zusage ↔ Zustand

1. **Die Zustellung fehlt, und das ist der schwerwiegendste Punkt.** MR-053 ist
   aus `AGENTS.md` §5 **und** `harness/README.md` verlinkt. MR-054 regiert die
   beiden **anderen** Blöcke desselben Absatzes und ist von dort nicht
   erreichbar. Ein Lauf, der nur das Briefing liest, erfährt die neue Pflicht
   nicht.
2. **Die DoD-Prämisse war falsch, und die Umsetzung hat sie korrigiert statt sie
   zu erfüllen.** Die DoD sagt, MR-051s „dort steht keine" gelte „nicht mehr".
   Gemessen: genau **13** Marker-Treffer in den drei Verzeichnissen, **keiner**
   wirksam. MR-054 formuliert deshalb „wird absehbar unwahr" — richtig; der
   DoD-Text selbst überzieht gegenüber dem Befund.
3. **Ein Präzisions-Fehler in der Begründung.** Commit und MR-054 sagen, die 13
   Treffer seien „alle Erwähnungen der Syntax **in Inline-Code**". Mindestens
   einer — `done/slice-152-citations-scharfschalten.md:41` — steht in einem
   Fenced-Block. Der Schluss (keiner wirksam) trägt trotzdem, weil Fences
   ebenfalls ausgenommen sind; die genannte Begründung deckt diesen Treffer
   nicht.
4. **Das zweite Zitat belegt nicht, was seine Anrede behauptet.** Der Block
   schreibt „Die Regel, die diesen Schritt vorschreibt:" und zitiert
   `modul-05:224-225` — die Notations-Nebenregel, nicht den vorschreibenden Satz
   (`:219`). Formal ist die DoD erfüllt; inhaltlich ist es die BEO-012-Klasse,
   die der Slice zwei Zeilen darüber selbst als gesichtet führt.
5. **Die Praxis ist bereits breiter als MR-054s Geltungsbereich.** Der nennt
   „die beiden kanonischen Vorprüfungs-Blöcke jedes neu angelegten
   **Slice-Plans**". Tatsächlich tragen slice-172 §1 und das Wellendokument
   welle-86 bereits Direktiven. Nichts davon ist falsch — aber der
   Geltungsbereich beschreibt die gelebte Form nicht.
6. **Kleinere Punkte.** MR-054s `Ersetzt-Baseline-Regel`-Link trägt das Label
   „§Zwei Schritte vor der Modus-Begründung", die URL aber **kein**
   Anker-Fragment — sie landet am Dateikopf. — MR-051 selbst hat keinen
   Vorwärts-Zeiger auf MR-054; das entspricht der Hauspraxis, ein Leser von
   MR-051 allein bekommt aber die unkorrigierte Geltungsbereich-Aussage.
