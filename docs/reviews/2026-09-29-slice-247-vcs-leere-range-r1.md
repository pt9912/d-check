# Review R1 — slice-247: vcs-Modul meldet stilles Grün über leerer, auflösbarer Range

- **Review-Art:** Code — geprüft gegen den Slice-Plan (`slice-247`), `DC-FA-VCS-002`,
  `ADR-0024`/`ADR-0027`, `MR-025`, `MR-032`, `AGENTS.md` §3/§5 (Maintainability;
  DoD-Abhakung ist nicht Gegenstand dieses Reviews)
- **Gegenstand:** `slice-247` · Range `5a45b0b2..HEAD` — `5f19f081` (Claim, reiner
  Move), `f4522207` (Implementierung: Lastenheft 0.93.3, Spezifikation
  §DC-FA-VCS-002.a, Adapter-Fix, zwei Tests)
- **Skill:** `reviewer.md` @ 1.16.0
- **Modell-ID:** glm-5.3-flash (Z.ai)
- **Datum:** 2026-09-29
- **Eingangs-Kontext:** Slice-Plan (`slice-247`); `DC-FA-VCS-002` (Lastenheft
  0.93.3), `DC-FA-VCS-001`/`DC-FA-COMMITS-001` als Nachbar-Verträge;
  `ADR-0024` (Modul `vcs`), `ADR-0027` (Modul `commits`); `MR-025`, `MR-032`;
  Hard Rules `AGENTS.md` §3.4/§3.7; Baseline `v6.13.0` ·
  `regelwerk/modul-10-review-harness.md` §Ziel-Form; vorherige Findings am
  Modul: slice-220 (CO-001), slice-245 (stille-Grün-Proben A–D am shallow-Clone)

---

## Findings

### M-1 — Schwester-Sektion §DC-FA-COMMITS-001.a behauptet weiter die gemeinsame Range-Semantik

- **Kategorie:** MEDIUM
- **Quelle:** `MR-025` (Semantik-Änderung: Spiegel) · Prüffrage 11
- **Pfad:** `spec/spezifikation.md:1983` und `spec/spezifikation.md:1991` (neu:
  `spec/spezifikation.md:1958` — §DC-FA-VCS-002.a)
- **Befund:** §DC-FA-COMMITS-001.a sagt, `--range` folge „**dieselbe** Range-Semantik
  wie `vcs`", und erklärt: „Eine **gültige** Range mit 0 Commits ist kein Fehler
  (nichts zu prüfen ⇒ Exit 0, wie `vcs` ohne geänderte Datei)". Seit
  §DC-FA-VCS-002.a in derselben Datei gilt in `vcs` das Gegenteil (auflösbare,
  leere Range ⇒ Exit ≠ 0) — zwei Module derselben Eingabe-Klasse haben jetzt
  verschiedene Range-Verträge, und die Gegenbehauptung steht stehengeblieben im
  Nachbarn. Die Out-of-Scope-Zeile der neuen DC (`spec/lastenheft.md:2207`)
  erklärt den Ausschluss des `commits`-Moduls, benennt aber keinen Grund, der
  die Divergenz trägt — und die Spiegel-Liste des Plans (§6, Risiko 2) nannte
  die Spezifikation als Spiegel, zog sie aber nur für den Neuzugang, nicht für
  die Schwester-Sektion.
- **Verifizierbar:** ja — `make doc-check` bleibt grün (Anker/Links lösen auf);
  der Befund ist semantisch, kein Gate fängt ihn.
- **Klasse:** `modul-divergenz-selbe-eingabeklasse`

### M-2 — Plan-Ausgrenzung des `commits`-Moduls zitiert Probe B über ihren Geltungsbereich hinaus

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 9 (`BEO-ALL/citation-stretched-beyond-scope`)
- **Pfad:** `slice-247` · §1 — `docs/plan/planning/in-progress/slice-247-vcs-leere-range-stilles-gruen.md:46`
- **Befund:** Die Abgrenzung begründet „Bestand bleibt bewusst stehen; keine
  Änderung nötig" mit „commits-Modul: bricht auf leerer Range bereits laut ab
  (slice-245, Probe B)". Probe B ist ein shallow-Klon: dort bricht `commits`
  ab, weil der Vorfahren-Walk an der Shallow-Grenze scheitert („Range-Basis-
  Vorfahren nicht lesbar", `internal/adapter/driven/git/git.go:200`) — nicht
  wegen eines Leerfall-Checks. Auf vollem Fixture ist `commits` auf leerer,
  auflösbarer Range still grün; das Vertrag-Testat
  `internal/hexagon/core/rules/commits_test.go:136`
  (`TestCheckCommitsEmptyRange`: 0 Befunde, Exit 0) sagt das Gegenteil der
  Plan-Behauptung. Die Begründung trägt in ihrer allgemeinen Form nicht; der
  stille-Grün-Zustand, dessen Klasse DC-FA-VCS-002 gerade für `vcs`
  beseitigt, besteht im `commits`-Modul (volles Fixture) fort.
- **Verifizierbar:** ja — `make test` (TestCheckCommitsEmptyRange) bzw. Probe:
  Modul `commits` mit `--range R..R` auf vollem Klon ⇒ Exit 0.
- **Klasse:** `BEO-ALL/citation-stretched-beyond-scope`

### M-3 — DC-Grenzen-Liste nennt nicht die Shallow-Grenze über den Leerfall hinaus

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 18 (`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`)
- **Pfad:** `spec/lastenheft.md:2207` (Out-of-Scope) · `spec/spezifikation.md:1958-1968`
  · `internal/adapter/driven/git/git.go:117`
- **Befund:** Der neue Check walkt die Vorfahren der Basis **komplett**
  (`a.ancestors`, git.go:117). In einem shallow-Klon mit Tiefe ≥ 2 scheitert
  dieser Walk an der Shallow-Grenze — `AllPaths` bricht dann mit „Range-Basis-
  Vorfahren nicht lesbar" ab, **auch bei nicht-leerer Range**, die vor dem Fix
  nur die beiden Endpunkt-Commits/-Trees brauchte. Abhilfe ist ausschließlich
  `fetch-depth: 0`. Weder die Out-of-Scope-Liste (drei Einträge: `commits`,
  history-range-guard, Range-Scoping) noch §DC-FA-VCS-002.a („der Adapter löst
  beide Refs zu Commits auf" — der Ancestry-Walk bleibt unerwähnt) nennen diese
  Grenze. Das Ergebnis ist fail-closed (kein falsches Urteil); die Lücke ist
  die undeklarierte Reichweite des neuen Laut-Seins.
- **Verifizierbar:** ja — Probe am Image: shallow-Klon Tiefe 2 mit nicht-leerer
  Range ⇒ vor dem Fix Befunde, nach dem Fix Exit 2 „Range-Basis-Vorfahren
  nicht lesbar".
- **Klasse:** `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`

### L-1 — Test-Plan verspricht shallow-Clone-Aufbau, geliefert wird ein volles Fixture

- **Kategorie:** LOW
- **Quelle:** Maintainability (Plan-Gebundenheit, `AGENTS.md` §6 Schritt 4)
- **Pfad:** `internal/adapter/driven/git/git_test.go:652` gegen `slice-247` · §3 (Plan-Zeile 89)
- **Befund:** Plan §3: „Negative-Fall ist slice-245s shallow-Clone-Aufbau".
  `TestAllPathsLeereRange` deckt Gleichheit und Inversion auf einem
  vollständigen Fixture ab; der shallow-Klon — der benannte Prüffall der
  DC — hat keinen automatisierten Test. Für `HEAD..HEAD` ist das semantisch
  gleichwertig (der Gleichheits-Check läuft vor jedem Parent-Walk), aber die
  Inversion im shallow-Klon läuft über einen anderen Fehlerpfad (M-3) und ist
  ungetestet. Die Abweichung ist weder im Code noch im Plan deklariert.
- **Verifizierbar:** ja — `make test`: kein Test erzeugt `.git/shallow`.
- **Klasse:** `plan-abweichung-undeclariert`

### L-2 — `Verantwortlich:` trotz Claim ungesetzt

- **Kategorie:** LOW
- **Quelle:** `v6.13.0` · `regelwerk/modul-05-planning-harness.md` §Trigger je Lifecycle-Übergang
- **Pfad:** `slice-247` · Kopf (Plan-Zeile 25)
- **Befund:** Der Claim (`5f19f081`, reiner Move nach `in-progress/`) hat das
  Feld `**Verantwortlich:**` nicht gesetzt — der Start-Trigger des Plans
  selbst (§4) verlangt „`Verantwortlich:` gesetzt". Das Feld ist Deklaration,
  kein Sensor hält es; vor der Closure zu korrigieren.
- **Verifizierbar:** nein (Deklaration).
- **Klasse:** `lifecycle-feld-vergessen`

### I-1 — Handbuch-Fehlerbild für die neue Meldung ist ein Release-Prep-Restposten

- **Kategorie:** INFO
- **Quelle:** `MR-025` · Dokumentations-Regel 17
- **Pfad:** `docs/user/benutzerhandbuch.md:1795-1806`
- **Befund:** Der Handbuch-Abschnitt zum Modul `vcs` listet die Exit-2-Fehlerbilder
  auf; die neue Meldung „Range-Leerfall …" fehlt dort. Die Pflege gehört per
  Regel 17 in die Release-Prep, nicht in den Feature-Commit — als Restposten
  führen, damit sie nicht still veraltet. Die `--print-mk`-Form braucht
  dagegen keinen Spiegel (das Fragment trägt keine Semantik-Texte) — geprüft.
- **Verifizierbar:** nein (Prosa).
- **Klasse:** `doku-drift-release-prep`

---

## Negativbefunde (geprüft, ohne Befund)

- **Leerfall-Äquivalenz:** Leer ⟺ Spitze ∈ ancestors(Basis) inklusive
  Gleichheit — korrekt implementiert (`ancestors()` schließt h selbst ein,
  git.go:229-252). Kein Befund.
- **Gleichheits-Fall vor dem Walk:** geprüft (git.go:114 vor 117) und tragend —
  ohne ihn liefe der shallow-`HEAD..HEAD`-Fall in „Range-Basis-Vorfahren nicht
  lesbar" statt in die benannte Leerfall-Meldung. Kein Befund.
- **Unauflösbare Basis/Spitze bleibt laut:** ResolveRevision-Fehler
  (git.go:107-112) ⇒ Exit 2 wie vor dem Fix; Kern-Negativtest vorhanden
  (`internal/hexagon/core/rules/vcs_test.go:246`). Kein Befund.
- **`TestAllPathsUnlesbarerUnterbaum`-Fix:** die Umstellung auf eine
  nicht-leere Range (git_test.go:317-321) erhält die Regression — ohne den
  zweiten Commit würde der Test am Gleichheits-Check hängen und aus dem
  falschen Grund grün; so erreicht der Fehlerpfad `walkTree` wie zuvor. Die
  Begründung trägt. Kein Befund.
- **§3.4 Spec-Straten:** grep über Lastenheft und Spezifikation — keine
  Slice-Referenz in den neuen Zeilen; die bestehende slice-220-Stelle
  (spezifikation.md:1871) trägt den Provenance-Marker. Kein Befund.
- **`MR-032`:** Versions-Bump 0.93.2 → 0.93.3 plus Historie-Zeile
  (lastenheft.md:3969/3972) vor der Closure. Kein Befund.
- **`AGENTS.md` §3.7:** neue Kommentare tragen `DC-FA-VCS-002` als auflösbares
  Herkunfts-Feld, keine Review-Historie, keine Slice-Nummern; die Meldungen
  folgen der bestehenden Hausform (vgl. „Range-Basis %q nicht auflösbar",
  git.go:108, Bestand). Kein Befund.
- **`ADR-0005`/Hexagon:** keine neuen Importe, keine Import-Richtung gebrochen. Kein Befund.
- **`DC-QA-03`:** alle Lesevorgänge lokal über die Objektdatenbank, kein Netz
  außerhalb `external`. Kein Befund.
- **commits-Vertrag im Code unberührt:** `CommitMessages` ist im Diff nicht
  angetastet — die Änderung bleibt beim `vcs`-Adapter. Kein Befund
  (die *Begründungs*-Hälfte ist M-2).

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 (M-1, M-2, M-3) |
| LOW | 2 (L-1, L-2) |
| INFO | 1 (I-1) |

## Verdikt

**MEDIUM blockiert typischerweise** — hier: die drei MEDIUM blockieren die
Closure, nicht die Code-Abnahme. Der Fix selbst ist semantisch korrekt
(Leer ⟺ Spitze ∈ ancestors(Basis) inkl. Gleichheit), fail-closed, deterministisch
(DC-QA-02: der Walk ist input-agnostisch) und hält den Kostenrahmen des
Schwestermoduls ein; die DC-Form deckt Gleichheit und Inversion ab, und die
Test-Abdeckung (Gleichheit, invertiert, Normalfall, unlesbarer Unterbaum) ist
soweit vorhanden. Die drei MEDIUM sind Spec-/Plan-Artefakte: Spiegel-Nachzug in
§DC-FA-COMMITS-001.a (M-1), Korrektur der Plan-Begründung in §1 (M-2) und
Nachtrag der Shallow-Grenze in §DC-FA-VCS-002.a bzw. der Out-of-Scope-Liste (M-3)
— kein Code-Rework. Nach deren Einarbeitung ist der Slice review-seitig frei
für die Verifikation.
