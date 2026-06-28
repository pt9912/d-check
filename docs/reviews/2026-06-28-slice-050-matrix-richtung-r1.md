# Review-Report — slice-050 (matrix: klasseninterne Verweisrichtung `DC-FA-MTX-002`)

## Kopf-Metadaten

- **Gegenstand:** uncommittete Änderungen slice-050 (Modul `matrix`, neues
  Feature klasseninterne Verweisrichtung über `order`/`direction`, Befund
  `matrix-downward`).
- **Reviewer:** unabhängiger Reviewer (Subagent), Lauf R1.
- **Datum:** 2026-06-28.
- **Diff-Range:** Arbeitsbaum (getrackt: `git diff`; untrackt: `ADR-0021`,
  `slice-050`). 11 getrackte Dateien geändert, 2 neue Dateien.
- **Betroffene Anforderungen/ADRs:** `DC-FA-MTX-002` (neu, Lastenheft),
  `DC-FA-MTX-001`, `ADR-0021`, `MR-006`, `MR-001`, `DC-QA-02`, `DC-QA-03`.
- **Skill-Basis:** `.harness/skills/reviewer.md` v1.0.0.
- **Eigene Verifikation:** `make doc-check` ausgeführt (Ground Truth, s. u.).
- **Hinweis zum Schema:** Befunde solution-frei (Skill-Regel); Fix-Hinweise
  stehen getrennt unter „Übergabe". Severity-Tier `BLOCKER` = explizite
  Eskalation über Skill-`HIGH` für den dispositiven Merge-Blocker.

---

## Findings

### BLOCKER-1 — Anforderung `DC-FA-MTX-002` fehlt im Lastenheft; doc-check ist rot

- **kategorie:** BLOCKER (Skill-HIGH: Korrektheits-/Harness-Lüge-Klasse)
- **quelle:** Anforderungs-Anlege-Prozess (`harness/conventions.md` §
  Anforderungs-Anlege-Prozess), `MR-001`, `DC-FA-MTX-002`, doc-check-Gate
- **pfad:** `spec/lastenheft.md` (gesamt; Datei **nicht** im `git status` —
  unverändert) — verweisende Bruchstellen: `CHANGELOG.md:22`,
  `docs/plan/adr/0021-matrix-klasseninterne-verweisrichtung.md:6`,
  `docs/plan/adr/README.md:36`,
  `docs/plan/planning/in-progress/roadmap.md:20`,
  `docs/plan/planning/in-progress/slice-050-matrix-klasseninterne-richtung.md:11`,
  `docs/plan/planning/in-progress/slice-050-matrix-klasseninterne-richtung.md:59`,
  `spec/spezifikation.md:565`, `spec/spezifikation.md:1153`
- **befund:** Das Lastenheft trägt **keine** Anforderung `DC-FA-MTX-002`:
  kein §3-Abschnitt, keine drei Akzeptanzkriterien, kein Out-of-Scope, kein
  Versions-Bump (Datei steht weiter auf **Version 0.29.0**), keine
  §7-Historie-Zeile. Alle acht Dokumentstellen verlinken dennoch den Anker
  `lastenheft.md#dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix`,
  der nicht existiert. Dadurch meldet `make doc-check` **8 ×
  `anchor-missing`** und endet mit Exit 1 — das Gate ist im Arbeitsbaum rot,
  nicht grün; die Greenfield-Reihenfolge „Doc führt, Code folgt" ist
  invertiert (Code/ADR/Spezifikation/Changelog existieren, der Vertrag im
  Lastenheft nicht).
- **verifizierbar:** ja — `make doc-check` (Exit 1, 8 Befunde; verifiziert in
  diesem Lauf).

### HIGH-1 — Falsche Grün-/Fertig-Attestierung in DoD, ADR und CHANGELOG

- **kategorie:** HIGH (Skill: Harness-Lüge — behauptetes Grün entgegen realem
  Gate-Stand)
- **quelle:** `MR-004`/Beleg-Pflicht, `DC-QA-02`-Historie-Konvention
- **pfad:**
  `docs/plan/planning/in-progress/slice-050-matrix-klasseninterne-richtung.md:59`
  (DoD-Punkt „Spec (doc-first)" mit `[x]` abgehakt),
  `docs/plan/adr/0021-matrix-klasseninterne-verweisrichtung.md:77`
  (Geschichte: „`make gates` grün. Status Accepted."), `CHANGELOG.md:22`/`:23`
  („Lastenheft 0.30.0")
- **befund:** Der DoD-Punkt „neue Anforderung `DC-FA-MTX-002` im Lastenheft
  (… Versions-Bump 0.30.0 + §7-Historie)" ist als erledigt `[x]` markiert,
  obwohl das Lastenheft unverändert ist; die ADR-Geschichte und das CHANGELOG
  behaupten Lastenheft-Version 0.30.0 bzw. „`make gates` grün", was der
  tatsächliche doc-check-Lauf (Exit 1) widerlegt. Die Attestierung beschreibt
  einen Zustand, den das Gate nicht hergibt.
- **verifizierbar:** ja — Abgleich `git status` (Lastenheft fehlt) +
  `make doc-check` (rot) gegen die `[x]`-/„grün"-/„0.30.0"-Aussagen.

### MEDIUM-1 — Vertrags-Zweige ohne Test: rangfreies Ziel, Gleichrang, Selbstverweis

- **kategorie:** MEDIUM (Skill: fehlende Negativtests bei neuem öffentlichem
  Vertrag)
- **quelle:** `DC-FA-MTX-002`, `spec/spezifikation.md` §DC-FA-MTX-001.a
  Schritt 5
- **pfad:**
  `internal/hexagon/core/rules/matrix_test.go:250` (`TestMatrixDownwardRichtung`)
- **befund:** Die Tests decken aufwärts/abwärts/transitiv und **rangfreie
  Quelle** (`notiz.md`) ab, aber nicht die ebenfalls vertraglich zugesicherten
  Zweige: ranghaftige Quelle → **rangfreies Ziel** (`!ok2`-Pfad in
  `downwardFinding`), **Gleichrang** (`si == di`, vom `si >= di`-Guard
  abgedeckt, aber ungeprüft) und **Selbstverweis** einer Datei auf sich selbst.
  Eine Regression in einem dieser Zweige bliebe von der Suite unbemerkt.
- **verifizierbar:** nein (kein Gate erkennt einen fehlenden Test; Risiko ist
  niedrig, da `downwardFinding` die Zweige korrekt über `!ok2`/`si >= di`
  behandelt).

### INFO-1 — Invariante „nie beide Befunde an einer Kante" gilt nur ohne Intra-Klassen-Regel

- **kategorie:** INFO (dokumentationswürdige Annahme)
- **quelle:** Maintainability, `ADR-0021`-Risiko, slice-050 §4
- **pfad:**
  `docs/plan/planning/in-progress/slice-050-matrix-klasseninterne-richtung.md:81`
  („eine Kante kann nie beide auslösen"), `spec/spezifikation.md:573`
  („unabhängig von `matrix-forbidden`/`matrix-inactive`")
- **befund:** Die Aussage, eine Kante könne nie zugleich `matrix-downward` und
  `matrix-forbidden` auslösen, hält nur, solange keine Intra-Klassen-Regel
  `{from: X, to: X, allow: false}` deklariert ist. Mit einer solchen Regel auf
  einer geordneten Klasse erzeugt eine Abwärtskante beide Befunde gleichzeitig
  — kein Code-Defekt (das Doppelsignal ist konsistent), aber eine
  unausgesprochene Annahme der Invariante.
- **verifizierbar:** ja — ein konstruierter Config-/Regel-Test (geordnete
  Klasse + Selbst-Regel `allow:false`) zeigt beide Befunde an derselben Kante.

---

## Negativbefunde (geprüft, ohne Befund)

- **Rang-Logik (Korrektheit):** `si < di` = abwärts, `si >= di` (inkl.
  Gleichrang und aufwärts) kein Befund; rangfreie Seite über `!ok1`/`!ok2`
  ausgenommen; Selbstverweis (`file == rel`, gleicher Rang) erzeugt korrekt
  nichts. `internal/hexagon/core/rules/matrix.go:205` — kein Befund.
- **fail-closed-Config (3 Fehlerfälle):** `order` ohne `direction`,
  `direction` ohne `order`, unbekannter `direction`-Wert sind in
  `validateMatrixDirection` vollständig abgedeckt und in
  `TestDecode_MatrixDirectionFailClosed` getestet
  (`internal/adapter/driven/configyaml/configyaml.go:413`) — kein Befund.
- **Determinismus / Default-aus (`DC-QA-02`):** ohne `order`/`direction` gibt
  `orderedClass` `nil` zurück und `downwardFinding` kurzschließt; keine
  Map-Iteration im Befundpfad (Slice-Iteration in `orderedClass`/`rankOf`);
  `TestMatrixDownwardDefaultAus` bestätigt byte-identischen Befundsatz
  (`internal/hexagon/core/rules/matrix.go:180`) — kein Befund.
- **Wechselwirkung exclude-sections / forbidden / inactive:** der
  Downward-Check sitzt hinter dem `inRanges`-Skip und ist von
  `matrix-forbidden`/`matrix-inactive` entkoppelt; klassenübergreifende Kanten
  (`dstClass != srcClass`) werden ausgeschlossen
  (`internal/hexagon/core/rules/matrix.go:62`) — kein Befund.
- **ADR-Schärfungsrichtung (`MR-006`/`MR-001`):** `ADR-0021` `Bezug:` zeigt
  aufwärts aufs Lastenheft, `Schärft:` auf die Spezifikation (§DC-FA-MTX-001.a)
  — kein Abwärtsverweis aufs Lastenheft; konventionskonform
  (`docs/plan/adr/0021-matrix-klasseninterne-verweisrichtung.md:11`) — kein
  Befund.
- **Lint/Komplexität (Refactor):** `downwardFinding` sauber extrahiert
  (gocognit-Entlastung), reine Funktion, deterministisch; Import-Gruppierung
  wie im Bestand (kein neuer Lint-Bruch)
  (`internal/hexagon/core/rules/matrix.go:202`) — kein Befund.
- **Code ↔ Spezifikation (§DC-FA-MTX-001.a Schritt 5):** Rang = First-Match-Glob,
  rangfrei ausgenommen, transitiv über Rangvergleich, beide Felder gekoppelt —
  der Code deckt die Spezifikations-Sektion vollständig (die Vertrags-Lücke
  liegt eine Stufe höher im Lastenheft, s. BLOCKER-1)
  (`spec/spezifikation.md:564`) — kein Befund.

---

## Kategorie-Summary

| Severity | Anzahl |
| --- | --- |
| BLOCKER | 1 |
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 1 |

---

## Verdikt

**Nicht mergebar.** BLOCKER-1 ist dispositiv: die kanonische Anforderung
`DC-FA-MTX-002` existiert im Lastenheft nicht (kein §3-Abschnitt, keine
Akzeptanzkriterien, kein Versions-Bump auf 0.30.0, keine §7-Historie), und
`make doc-check` ist im Arbeitsbaum **rot** (Exit 1, 8 × `anchor-missing`) —
entgegen der Handoff-Behauptung „`make gates` ist bereits grün". HIGH-1
(falsche Grün-/Fertig-Attestierung in DoD/ADR/CHANGELOG) verschärft das.
Code-, Test- und Determinismus-Substanz des Features sind solide (alle
Negativbefunde grün); der Slice scheitert nicht am Algorithmus, sondern an der
fehlenden Vertrags-Stufe und der dadurch gebrochenen Anker-Kette. MEDIUM-1
(Test-Lücken) vor Abschluss schließen; INFO-1 ist eine Annahme zum Festhalten.

### Übergabe (Fix-Hinweise, außerhalb der Befund-Felder)

- BLOCKER-1: `DC-FA-MTX-002`-Abschnitt in `spec/lastenheft.md` §3 anlegen
  (Heading exakt = Anker `dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix`),
  Happy/Boundary/Negative + Out-of-Scope, Versions-Bump 0.29.0 → 0.30.0,
  §7-Historie-Zeile; danach `make doc-check` grün nachweisen.
- HIGH-1: DoD-`[x]` erst nach echtem Lastenheft-Stand setzen; ADR-Geschichte/
  CHANGELOG-„0.30.0"/„grün" mit dem belegbaren Gate-Lauf in Einklang bringen.
- MEDIUM-1: Tests für rangfreies Ziel, Gleichrang und Selbstverweis ergänzen.
