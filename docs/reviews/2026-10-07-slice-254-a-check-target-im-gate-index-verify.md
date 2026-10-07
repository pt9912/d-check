# Verifikation slice-254 — Das Target `a-check` aus `a-check.mk` steht im Gate-Index

- **Rolle:** Verifier. Die Frage ist: Bauen wir es richtig? Geprüft wird gegen Plan und DoD, nicht gegen den Diff (das ist die Rolle des Reviewers) und nicht gegen den realen Bedarf (Validator).
- **Gegenstand:** Plan `docs/plan/planning/in-progress/slice-254-a-check-target-im-gate-index.md` (DoD in §2, darunter DoD 2 mit vermerkter Plan-Änderung nach R1-F-2). Commits `289c8193` (feat), `a75d8c24` (R1-Report), `a47e0ce2` (R1-Einarbeitung), HEAD = `a47e0ce2`.
- **Bezug:** `MR-074` Bewegung 1 und Absatz *Zitat-Delta*; `v6.17.0` · `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt (Absatz *Ein Index, mehrere Eigentümer*, Z. 206–248); `v6.17.0` · `templates/.d-check.yml` `targets`-Block (Z. 27–50); `v6.17.0` · `regelwerk/modul-13-quality-gates.md` *Vorhanden ≠ behauptet* (Z. 125–129).
- **Modell-ID:** claude-opus-5-5 (Claude Agent SDK)
- **Datum:** 2026-10-07

## Selbst gefahrene Sensoren (echte Ausgabe)

| Lauf | Zustand | Ausgabe | Exit |
| --- | --- | --- | --- |
| `make gate-consistency` | HEAD | `d-check: 957 Datei(en) geprüft, 0 Befund(e)` | 0 |
| Brechen A: `· \`make a-check\`` aus `harness/README.md:73` entfernt | gebrochen | `a-check.mk:57	a-check	gate-undocumented	Makefile-Regel \`a-check\` ohne Deklaration in der Autoritäts-Doku harness/README.md` | 2 |
| Brechen B: `a-check-fake:` an `a-check.mk` angehängt | gebrochen | `a-check.mk:60	a-check-fake	gate-undocumented …` | 2 |
| Brechen C (Vorzustand der Konfiguration): `makefiles: [Makefile]` + `a-check-fake` im Fragment, Index-Zeile bleibt | gebrochen | einziger Befund `harness/README.md:73	a-check	gate-phantom`; **kein** Befund für `a-check-fake` | 2 |
| `make gates` | HEAD | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` (Coverage 94,70 % ≥ 93 %, lint `0 issues.`, semgrep `gesamt: 0 Befund(e)`) | 0 |

Jede Änderung wurde per Kopie aus dem Scratchpad zurückgesetzt. Danach war `git status` leer.

Brechen C belegt beide Richtungen. Der Eintrag `a-check.mk` in `targets.makefiles` trägt die Last: ohne ihn ist die Index-Zeile ein Phantom. Und Grenze 3 stimmt: ein Fragment, das nicht in `makefiles` steht, bringt undeklarierte Targets still mit.

## Prüfung gegen den Gegenstand (`AGENTS.md` §5 Regel 13)

**`.d-check.yml:925-940`, `targets`-Kommentar:**

- *„d-check dokumentiert jedes Target des Makefile UND des Fragments a-check.mk … in harness/README.md §Sensors"*. Wahr. `gate-consistency` meldet 0 Befunde bei `exempt-targets: []`, und Brechen A/B zeigt, dass das Fragment tatsächlich gelesen wird.
- *„(das Rezept hinter arch-check)"*. Wahr: `Makefile:71` lautet `arch-check: a-check …`, `a-check.mk:58-59` ist das Rezept.
- *„a-check.mk liegt im Wurzelverzeichnis und ist an die Repo-Politik angepasst"*. Wahr: `Makefile:33` lautet `include a-check.mk`, und der Kopf des Fragments selbst sagt „Erzeugt aus `a-check --print-mk` und an die Repo-Politik angepasst (ADR-0029)".
- *„die Regel der Baseline für Werkzeug-Teile gilt Fragmenten unter harness/mk/, die ein Werkzeug selbst erzeugt"*. Trifft den Wortlaut von `grundlagen-harness-dateien.md:206-219` („Bootstrap-Werkzeug, das Make-Fragmente unter `harness/mk/` erzeugt"; Bedingung „Das Werkzeug schreibt seinen Teil bei jedem Lauf neu"). Auch die Vorlage `.d-check.yml:35` passt: „d-check.mk erst ergaenzen, wenn `d-check --print-mk` gelaufen ist", Zweig `harness/mk/` getrennt in Z. 38–46.
- *„Grenze: targets folgt keinem include"*. Wahr gegen den Code: `internal/hexagon/core/rules/targets.go:88-111` (`collectMakefileRules`) liest nur die Dateien aus `cfg.Makefiles` zeilenweise gegen `makefileRuleRe`, eine Auswertung von `include` gibt es nicht. Brechen C bestätigt das im Verhalten.
- Das gekürzte Vorlagen-Zitat (R1-F-1) ist entfernt. `grep` über lebende Dateien findet es nur noch als Erwähnung in `MR-074:121`.

**`harness/sensors/gate-consistency.md:21-25`, Grenze 3:**

- *„gelesen werden genau die Dateien in targets.makefiles (heute Makefile und das Fragment a-check.mk)"*. Wahr (Code und Konfiguration).
- *„Ein weiteres eingebundenes Fragment bliebe ohne Befund, bis es dort eingetragen ist"*. Für die Richtung `gate-undocumented` wahr (Brechen C). Nuance siehe V-4.
- *„die Liste der include-Zeilen im Makefile zeigt, was gelesen werden müsste"*. Wahr für den Ist-Stand: genau eine `include`-Zeile (`Makefile:33`), und `a-check.mk` bindet nichts weiter ein.

**Zahlen der Commit-Botschaft `289c8193` (57/58):** Gezählt mit dem Regex des Moduls (`^[A-Za-z][A-Za-z0-9 _-]*:([^=]|$)`) per `grep -E`. Ergebnis: `Makefile` vor dem Slice (`289c8193^`) 57 eindeutige Regeln auf 57 Zeilen, `Makefile` an HEAD 57 (im Slice unverändert), `a-check.mk` 1. Damit „schon vorher 57, jetzt 58 mit dem Fragment": **stimmt**. Die entfernte Zahl „54" war tatsächlich veraltet.

## DoD-Verdikte

| # | DoD-Punkt | Verdikt | Beleg |
| --- | --- | --- | --- |
| 1 | `harness/README.md` §Sensors führt `make a-check` in der `arch-check`-Zeile, mit Vertrag und unveränderter Bindung | **erfüllt** | `harness/README.md:73`. Die erste Zelle trägt beide Targets, und beide werden erkannt (Brechen A: ohne die Zeile `gate-undocumented`; Brechen C: ohne Fragment `gate-phantom` auf genau diese Zeile). Der Vertrag nennt Rezept und Delegation, gemessen an `Makefile:71` und `a-check.mk:57-59`. Die Bindungs-Zelle ist im Diff `289c8193` byte-gleich. |
| 2 | `targets.makefiles` liest `a-check.mk`; der Kommentar sagt gemessen, was der Index deckt (Plan-Änderung nach R1-F-2); Bewusstes Brechen | **erfüllt** | `.d-check.yml:942` lautet `makefiles: [Makefile, a-check.mk]`. Jede Kommentar-Aussage ist gegen den Gegenstand geprüft (oben). Bewusstes Brechen selbst gefahren: A (`a-check.mk:57 a-check gate-undocumented`, Exit 2) und zusätzlich B (erfundenes Target im Fragment ⇒ Befund). Die Plan-Änderung steht kursiv an der DoD und kam vor dem Häkchen, nicht danach. |
| 3 | `make gates` grün | **erfüllt** | selbst gefahren, Exit 0, Ausgabe oben |
| 4 | Review durchgeführt, Report unter `docs/reviews/` | **erfüllt** | `docs/reviews/2026-10-07-slice-254-a-check-target-im-gate-index-r1.md` (getrennter Kontext, Verdikt „mit Auflage"). Einarbeitung siehe unten. |
| 5 | Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | **offen, erwartet** | §7 leer, §6-Ausgang `(offen)`. Das ist kein Befund: der Slice liegt in `in-progress/`. `make verify-closure-notes` greift erst nach dem Move nach `done/`. |

## R1-Einarbeitung gegen den R1-Report

| R1 | Forderung | Stand |
| --- | --- | --- |
| F-1 (MEDIUM) | Das gekürzte Zitat stützt eine Aussage, der seine ausgelassene Mitte widerspricht | **eingearbeitet.** Das Zitat ist aus `.d-check.yml` entfernt, und der tragende Grund steht an seiner Stelle (gegen den Baseline-Wortlaut geprüft). `MR-074` vermerkt die Neufassung (`MR-074:124-125`). Restpunkt zur Form siehe V-1. |
| F-2 (LOW) | Plan nennt eine Zahl, der Diff eine Aussage | **eingearbeitet.** DoD 2 und §3-Spiegel („Target-Aussage") sind nachgezogen, die Plan-Änderung ist vermerkt. |
| F-3 (LOW) | `include`-Achse fehlt in der Grenzen-Liste | **eingearbeitet.** Grenze 3 und eine Grenz-Zeile im `targets`-Kommentar, beide gegen Code und Verhalten bestätigt. |
| F-4 (INFO) | Stärkerer Anker für die Abgrenzung | **im Plan und in `.d-check.yml` eingearbeitet,** im `MR-074`-Vermerk nicht (V-2). |
| F-5 (INFO) | Spiegel-Entscheidung `arch-check.md` nicht ausgewiesen | **auf die Closure-Notiz verschoben** (laut Botschaft `a47e0ce2`). `harness/sensors/arch-check.md:6-7` nennt `a-check.mk` bereits. Die Closure-Notiz muss den Punkt tragen. |

## Neue Befunde

### V-1 — INFO

- **Pfad:** `harness/conventions/MR-074-baseline-v6170.md:119-125`
- **Kern:** Der Absatz *Zitat-Delta* sagt weiter im Präsens „der Kommentar … **zitiert** … mit Auslassung" und hängt dahinter „slice-254 hat ihn neu geschrieben, das gekürzte Zitat ist entfallen" an. Die erste Aussage stimmt am Gegenstand nicht mehr (`.d-check.yml` enthält kein Zitat). Der Absatz liest sich als Chronik eines aktiven Eintrags. Das ist inhaltlich auflösbar und kein DoD-Gegenstand.

### V-2 — INFO

- **Pfad:** `harness/conventions/MR-074-baseline-v6170.md:38-40`
- **Kern:** Der Einlösungs-Vermerk zu Bewegung 1 begründet weiter mit dem schwächeren Argument aus R1-F-4: „kein Werkzeug-Teil (a-check schreibt keinen)". Plan §1 und `.d-check.yml:936-937` stützen sich inzwischen auf das Geltungskriterium: Fragment im Wurzelverzeichnis, vom Repo angepasst, nicht unter `harness/mk/`. Die drei Stellen begründen dasselbe Ergebnis verschieden. Schreibt a-check künftig einmal einen Teil, trüge der Vermerk nicht mehr, die anderen beiden schon.

### V-3 — INFO

- **Pfad:** `harness/README.md:153` (Gate-Taxonomie)
- **Kern:** Die Klasse *Produkt-Gates* führt `arch-check`, nicht aber `a-check`. Das ist vertretbar, denn `a-check` ist das Rezept, an das `arch-check` delegiert, und die Index-Zeile sagt das. Kein Sensor liest die Taxonomie (`targets` liest nur `make X`-Tokens). Nur der Vollständigkeit halber genannt, keine Handlung verlangt.

### V-4 — INFO

- **Pfad:** `harness/sensors/gate-consistency.md:22-24`; `.d-check.yml:938-939`
- **Kern:** „Ein weiteres eingebundenes Fragment bliebe ohne Befund" gilt für die Richtung `gate-undocumented`. Deklariert jemand ein Target eines nicht gelisteten Fragments im Index, meldet der Sensor laut `gate-phantom` (Brechen C: `harness/README.md:73 a-check gate-phantom`). Die Grenze ist damit nicht falsch, nur einseitig formuliert. Die laute Gegenrichtung ist eher ein Vorteil als eine Lücke.

## Negativbefunde

1. **Abgrenzung (§1):** Es gibt kein `harness/mk/`, kein `authority-disjoint` und keine weitere `MR-074`-Bewegung im Diff (`git show --stat` der drei Commits). Die Abgrenzung ist nicht ausgeweitet.
2. **Rückführungs-Bedingung (§4):** „Das Fragment trägt mehr Targets als `a-check`" ist nicht eingetreten. Gemessen: genau 1 Regel in `a-check.mk`.
3. **Hard Rules:** §3.3 (die Plan-Datei wurde in `in-progress/` geändert, ohne Move im selben Commit), §3.4/§3.5 (keine Spec-Datei, keine ADR berührt), §3.6 (der Scan-Umfang wird erweitert, keine Senkung), §3.7 (der neue Kommentar trägt Zusage, Kopplung und Grenze, keine Slice-Nummer und kein Mess-Label). Ohne Befund.
4. **Traceability:** Alle drei Commits nennen `slice-254`, `289c8193` und `a47e0ce2` zusätzlich `MR-074` und `DC-FA-TGT-001`.

## Verdikt

**DoD 1–4 bestätigt, durch eigene Messung und bewusstes Brechen.** DoD 5 ist erwartungsgemäß offen. Es gibt keinen HIGH-, MEDIUM- oder LOW-Befund. V-1 bis V-4 sind INFO. Für die Closure: F-5 aus R1 in der Closure-Notiz entscheiden, V-1/V-2 nach Ermessen beim Closure-Commit glätten.
