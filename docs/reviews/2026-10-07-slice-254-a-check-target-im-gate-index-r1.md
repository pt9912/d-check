# Review slice-254 — Das Target `a-check` aus `a-check.mk` steht im Gate-Index (R1)

- **Review-Art:** Code-/Konfigurations-Review. Geprüft gegen den Slice-Plan `slice-254` (§1 Ziel und Abgrenzung, §3 Plan samt Spiegel-Liste, §6 Risiken), die Baseline-Regel zum Gate-Index mit mehreren Eigentümern, die Hard Rules (`AGENTS.md` §3.7, §3.8, §4, §5 Regeln 13/15/16) und die Konventionen (`MR-025`, `MR-039`, `MR-074`). **Nicht** geprüft gegen die DoD; das ist Sache der Verifikation in einem getrennten Kontext.
- **Gegenstand:** Commit `289c8193`, 3 Dateien: `.d-check.yml` (Kommentar und `targets.makefiles`), `harness/README.md` (Zeile `arch-check`), `harness/conventions/MR-074-baseline-v6170.md` (Vermerk zu Bewegung 1).
- **Skill:** `reviewer.md` @ 1.16.0 (`957aeedc`)
- **Modell-ID:** claude-opus-5-5 (Claude Agent SDK)
- **Datum:** 2026-10-07
- **Eingangs-Kontext:** Slice-Plan `slice-254`; `MR-074` Bewegung 1 und Absatz *Zitat-Delta*; `MR-039`; `v6.17.0` · `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt (Absatz *Ein Index, mehrere Eigentümer*); `v6.17.0` · `regelwerk/modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin), Absatz *Vorhanden ≠ behauptet*; `v6.17.0` · `templates/AGENTS.template.md` §4; `v6.17.0` · `templates/.d-check.yml` `targets`-Block; `DC-FA-TGT-001` mit Out-of-Scope; ADR-0029. Frühere Befunde zum selben Gegenstand: Review `slice-253` R1, F-1 (Behauptung „keine Werkzeug-Fragmente", record-claim-vs-diff) und F-3 (LOW, das elidierte Vorlagen-Zitat im `targets`-Kommentar).
- **Vom Reviewer gefahrene Sensoren:**
  - `make gate-consistency` auf dem Stand des Commits: `956 Datei(en) geprüft, 0 Befund(e)`.
  - **Bewusstes Brechen (a):** `· `make a-check`` aus der ersten Zelle entfernt. Ergebnis: `a-check.mk:57	a-check	gate-undocumented`, Exit 2 (`make: *** … Fehler 1`).
  - **Bewusstes Brechen (b):** Ein Phantom-Target `a-check-fake` an `a-check.mk` angehängt. Ergebnis: `a-check.mk:61	a-check-fake	gate-undocumented`, Exit 2.
  - Beide Änderungen per Kopie zurückgesetzt.
  - `make gates`: `[gates] … green`, EXIT=0.
  - Arbeitsbaum danach: `git status` leer bis auf diesen Report.

---

## Findings

### F-1 — MEDIUM

- **Kategorie:** MEDIUM. Kontext-Eskalation von LOW, weil dieselbe Beobachtung zum zweiten Mal am selben Artefakt auftritt; der Ort ist die Konfiguration eines Gates.
- **Quelle:** Anker „Quelle über ihren Geltungsbereich hinaus zitiert" (`AGENTS.md` §5 Regel 16); `MR-039`; `v6.17.0` · `templates/AGENTS.template.md` §4
- **Pfad:** `.d-check.yml:931-936` · „d-check dokumentiert jedes Target des Makefile UND des Fragments a-check.mk … dem EINEN Gate-Index (MR-071, aus der adoptierten Vorlage: "Der Gate-Index steht einmal ... Diese Datei fuehrt die Liste nicht.")"
- **Befund:** Der Kommentar wurde in diesem Commit neu geschrieben. Er sagt jetzt ausdrücklich, dass die Targets des Fragments im einen Index stehen, und belegt das mit dem elidierten Vorlagen-Zitat. Genau die ausgelassene Mitte lautet in `v6.17.0`: „Targets aus Werkzeug-Fragmenten stehen in dem Teil des Werkzeugs, den §Sensors verlinkt". Das Zitat stützt also die neue Aussage nur, weil der Satz fehlt, der ihr im Wortlaut widerspricht. Warum die Regel hier nicht greift (siehe F-4), steht nirgends im Kommentar. Der Schutz aus `MR-039` (ein historisches Zitat bleibt unangetastet) trägt nicht für einen Satz, der in diesem Commit neu geschrieben wurde. Auch der Vermerk in `MR-074` schließt den Absatz *Zitat-Delta* nicht. Dort heißt es: „seine Aussage hängt an Bewegung 1 oben". Bewegung 1 ist jetzt als eingelöst vermerkt, ob das Zitat die Aussage noch trägt, sagt der Vermerk aber nicht. **Failure-Szenario:** Der nächste Lauf liest den Kommentar gegen die vendorte Vorlage und findet die ausgelassene Klausel. Er hält das Repo dann für still abweichend und legt `harness/mk/a-check.md` an, oder er meldet eine MR-Pflicht, die nicht besteht.
- **Verifizierbar:** ja. `sed -n 203-207p .harness/baseline/v6.17.0/templates/AGENTS.template.md` neben `.d-check.yml:931-936` lesen. Kein Gate fängt das.
- **Klasse:** elided-quote-hides-counter-clause (Fortsetzung von `slice-253` R1 F-3)

### F-2 — LOW

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §6 Schritt 4 (eine Abweichung vom Plan ist eine Plan-Änderung); Plan `slice-254` §2/§3
- **Pfad:** Plan `slice-254` §2, zweiter Punkt · „der Kommentar nennt die Zahl der Targets gemessen"
- **Befund:** Der Diff entfernt die Zahl und setzt eine Aussage an ihre Stelle. Die Begründung steht in der Commit-Botschaft: die Zahl veraltet sonst wieder. In der Sache ist das richtig und konform zu `AGENTS.md` §5 Regel 13. Der Plan nennt aber weiter die gemessene Zahl, und auch die Spiegel-Zeile in §3 nennt „Target-Zahl". **Failure-Szenario:** Die Verifikation prüft gegen den Wortlaut der DoD, findet keine Zahl und meldet eine DoD-Verletzung. Oder die Abweichung bleibt still bestehen, und Plan und Diff widersprechen sich im `done/`-Bestand.
- **Verifizierbar:** ja. Plan §2 neben `.d-check.yml:931` lesen.
- **Klasse:** plan-form-drift-without-plan-update

### F-3 — LOW

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §3.8 (ein Modul verspricht nur über das, was es scannt); Prüffrage 18; `DC-FA-TGT-001` Out-of-Scope („Auflösen von `include`-Direktiven", `spec/lastenheft.md:3383`)
- **Pfad:** `harness/sensors/gate-consistency.md` §Grenze (außerhalb des Diffs, durch den Diff sichtbar geworden); `.d-check.yml:944` · `makefiles: [Makefile, a-check.mk]`
- **Befund:** Der Slice schließt eine Instanz einer still grünen Lage: Ein `include`-tes Fragment war für den Sensor unsichtbar. Die allgemeine Form der Lücke bleibt offen. `targets` verfolgt kein `include`, die Liste `makefiles` ist kuratiert, und ein künftiges zweites `include` ohne Eintrag dort bleibt wieder ohne Befund. Das steht im Out-of-Scope des Lastenhefts. Die Grenzen-Liste der Sensor-Datei nennt es nicht, und diese Liste liest ein Lauf über das Grün. Die Lücke liegt im Bestand, nicht im Diff, deshalb LOW.
- **Verifizierbar:** ja. Testweise ein zweites `include x.mk` mit einem eigenen Target anlegen, dann zeigt `make gate-consistency` `0 Befund(e)`. Nicht gefahren: Die Grenze ist durch den Vorzustand belegt, `a-check` war unter `makefiles: [Makefile]` ohne Befund.
- **Klasse:** grenze-liste-missing-include-axis

### F-4 — INFO

- **Kategorie:** INFO
- **Quelle:** `v6.17.0` · `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt; `v6.17.0` · `regelwerk/modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin); `v6.17.0` · `templates/.d-check.yml` `targets`-Block
- **Pfad:** Plan `slice-254` §1, erster Ausschluss; `MR-074` Vermerk · „kein Werkzeug-Teil (a-check schreibt keinen)"
- **Befund:** Das Ergebnis trägt, und es braucht **keinen** MR-Eintrag. Die Begründung stützt sich aber auf das schwächere Argument. „Werkzeug-eigen" ist eine Bedingung an einen vorhandenen Werkzeug-Teil und kein Kriterium dafür, ob das Regime überhaupt gilt. Das Kriterium für die Geltung steht im Satz davor: „Ein Bootstrap-Werkzeug, das Make-Fragmente unter `harness/mk/` erzeugt". `a-check.mk` liegt im Wurzelverzeichnis, stammt aus `--print-mk` und ist von Hand angepasst. Für genau diese Klasse zeigt die Baseline den anderen Weg: Den `d-check.mk`-Fall aus `--print-mk` führt sie in `harness/README.md` §Sensors oder namentlich exempt (`modul-13` *Vorhanden ≠ behauptet*). Die Vorlage `.d-check.yml` ergänzt dasselbe Fragment schlicht in `makefiles` („d-check.mk erst ergaenzen, wenn `d-check --print-mk` gelaufen ist") und führt den Zweig für `harness/mk/` getrennt. Mit diesem Anker wäre die Begründung unabhängig davon, ob a-check später einmal einen Teil schreibt.
- **Verifizierbar:** nein (Urteil).
- **Klasse:** applicability-argued-from-form-condition

### F-5 — INFO

- **Kategorie:** INFO
- **Quelle:** `MR-025` (Spiegel vor dem Editieren); Plan `slice-254` §3
- **Pfad:** Plan §3 · „`harness/sensors/arch-check.md` | update, falls gemessen | nennt er das Rezept?"
- **Befund:** Die Sensor-Datei bleibt unverändert. Ihr Vertrag nennt „ein include-bares `a-check.mk`", aber nicht das Target `a-check`. Das ist vertretbar, weil die Index-Zeile das Rezept jetzt nennt. Der Spiegel-Punkt ist im Diff und in der Botschaft aber nicht als entschieden ausgewiesen. Die Closure-Notiz kann das tragen.
- **Verifizierbar:** nein.
- **Klasse:** mirror-decision-unrecorded

---

## Negativbefunde

1. **Frage aus dem Auftrag: MR-Pflicht für den fehlenden Werkzeug-Teil.** Geprüft gegen die Bedingungen der Baseline, gegen *Vorhanden ≠ behauptet* in `modul-13` und gegen den `targets`-Block der Vorlage. Es gibt keine stille Abweichung, und kein MR ist nötig (Begründungsanker siehe F-4).
2. **Index-Zeile, Erkennung:** `targets` erkennt `make a-check` in der ersten Zelle, nachgewiesen durch das bewusste Brechen (a). Die Zeile beginnt mit `|`. Die Target-Zelle trägt nackte Target-Namen ohne Aufruf-Argumente.
3. **Index-Zeile, Bindung:** Die Bindung ist byte-gleich zum Vorzustand (ADR-0005, ADR-0012, ADR-0029, `DC-QA-03`). Die Vertrags-Ergänzung stimmt gegen den Gegenstand: `Makefile:71` lautet `arch-check: a-check`, `a-check.mk:57` ist das Rezept.
4. **Der Sensor hält das Fragment:** Ein neues Target im Fragment meldet `gate-undocumented`, nachgewiesen durch das bewusste Brechen (b). Das in §6 benannte Risiko verhält sich wie beschrieben. Es gibt keinen still grünen Pfad im Diff (Prüffrage 1).
5. **`AGENTS.md` §5 Regel 13, Aussage gegen Gegenstand:** „jedes Target des Makefile UND des Fragments a-check.mk" ist auf dem Stand des Commits wahr (`gate-consistency` 0 Befunde, `exempt-targets: []`). Das Makefile bindet nur ein `include` ein (`Makefile:33`), es gibt keine weitere Datei, die die Aussage übersehen würde. Die Aussage stimmt; nur ihr Beleg-Zitat ist problematisch (F-1).
6. **Zahlen in der Commit-Botschaft:** 57 Regeln im `Makefile`. Gezählt als Regel-Zeilen ohne Spezial-Targets: Die einfache Zählung ergibt 59, abzüglich `.PHONY` und `.NOTPARALLEL` bleiben 57. Gegengeprüft: 52 `.PHONY`-Namen plus 5 Regeln ohne `.PHONY` (`baseline-probe`, `doc-complete`, `hubdesc-pin-freshness`, `mention-coverage`, `review-coverage`) ergeben 57. Das `Makefile` ist im Commit unverändert, also stimmt „schon vorher 57". Mit `a-check` sind es 58. Die Zahlen stimmen, es gibt keine Überdehnung (Prüffrage 8).
7. **`AGENTS.md` §3.7, geänderter Kommentar:** Neuer Text sind Zusage („dokumentiert jedes Target …") und Kopplung („das Rezept hinter arch-check"), dazu ein Rang-Zeiger (Vorlagen-Zitat, zum Inhalt siehe F-1). Keine Review-Historie, keine Slice-Nummer, keine Mess-Labels. Die nun weggefallene Zahl war ein veraltendes Mess-Label, ihr Entfernen ist §3.7-konform.
8. **`MR-074`-Vermerk, Form:** Der Vermerk steht in einem **aktiven** Eintrag (`harness/conventions/`, nicht `done/`) als Zusatz innerhalb von Bewegung 1. Der bestehende Wortlaut ist nicht überschrieben, Kandidat und Ausgang stehen nebeneinander, und `slice-254` dient als auflösbare Kennung. Das ist dieselbe Form, die der Eintrag selbst für die Folge-Slices ankündigt. `MR-045` gilt nur für `AGENTS.md` und `harness/README.md`, nicht für MR-Dateien. Die Kopfzeile „keine davon übernommen" bleibt wörtlich wahr: Übernommen wurde die Werkzeug-Teil-Form nicht. Zum Inhalt siehe F-1 (offener Zitat-Delta-Absatz) und F-4.
9. **Abgrenzung des Plans:** Es gibt kein `harness/mk/`, kein `authority-disjoint` und keine Bewegung 2 bis 6 im Diff. Die Abgrenzung wurde nicht ausgeweitet.
10. **Hard Rules ohne Gegenstand:** §3.1 (kein Host-Werkzeug im Diff), §3.2 (keine Suppression), §3.3 (kein Move), §3.4 (keine Spec-Datei berührt), §3.5 (keine ADR berührt), §3.6 (keine Schwelle gesenkt; der Scan-Umfang des Sensors wird **erweitert**), §3.9 (kein Workflow): ohne Befund.
11. **Spiegel der Konfiguration:** `harness/sensors/gate-consistency.md`, `AGENTS.md` §4 und `harness/README.md` führen die `makefiles`-Liste nicht als Wert. Es gibt keinen weiteren Spiegel, der hätte nachgezogen werden müssen.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 |
| LOW | 2 | F-2, F-3 |
| INFO | 2 | F-4, F-5 |

## Verdikt

**Mit Auflage.** Die Mechanik ist korrekt und zweifach durch bewusstes Brechen belegt. Die Kernfrage (MR-Pflicht) ist negativ entschieden. F-1 blockiert typischerweise: Der neu geschriebene Kommentar stützt eine Aussage mit einem Zitat, dessen ausgelassene Mitte ihr im aktuellen Wortlaut widerspricht. Dieselbe Klasse wurde schon in `slice-253` R1 als LOW gemeldet, deshalb ist eine Wiederholung zu vermeiden. F-2 gehört vor die Verifikation geklärt. F-3 bis F-5 kann der Implementer annehmen oder begründen.
