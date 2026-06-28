# Review-Report — slice-050 (matrix: klasseninterne Verweisrichtung) — R2

## Kopf-Metadaten

- **Gegenstand:** uncommittete Änderungen slice-050 (Modul `matrix`, neues
  Feature klasseninterne Verweisrichtung über `order`/`direction`, Grund-Code
  `matrix-downward`), **zweiter Lauf** — Verifikation der R1-Auflösungen plus
  Review der Deltas seit R1.
- **Reviewer:** unabhängiger Reviewer (Subagent), Lauf R2.
- **Datum:** 2026-06-28.
- **Diff-Range:** Arbeitsbaum (`git diff` getrackt; untrackt: `ADR-0021`,
  `slice-050`, R1-/R2-Report). 15 getrackte Dateien geändert.
- **Vorlauf:** R1-Report `2026-06-28-slice-050-matrix-richtung-r1.md`
  (1 BLOCKER, 1 HIGH, 1 MEDIUM, 1 INFO).
- **Skill-Basis:** `.harness/skills/reviewer.md` v1.0.0.
- **Eigene Verifikation:** `make doc-check` ausgeführt (Ground Truth, s. u.).
  Befunde solution-frei; Fix-Hinweise getrennt unter „Übergabe". Tier
  `BLOCKER` = Eskalation über Skill-`HIGH` (dispositiver Merge-Blocker).

---

## Verifikation der R1-Befunde

### R1-BLOCKER-1 — Anforderung fehlte im Lastenheft, doc-check rot → **AUFGELÖST (mit Restpunkt)**

- Die Anforderung steht jetzt im Lastenheft `spec/lastenheft.md:648` als
  vollständige §3-Sektion: Heading **exakt** gleich dem von den acht Verweisern
  genutzten Anker, Beschreibung (`:650`–`:670`), drei Akzeptanzkriterien
  Happy/Boundary/Negative (`:674`–`:676`) und Out-of-Scope (`:678`).
- §7-Historie-Zeile 0.30.0 ergänzt (`spec/lastenheft.md:1106`).
- **Eigener `make doc-check`-Lauf:** Image gebaut, read-only/netzlos
  (`docker run --rm --network none -v …:/repo:ro`), Ergebnis
  `138 Datei(en) geprüft, 0 Befund(e)`, **Exit 0**. Die Anker-Kette
  (`CHANGELOG.md:22`, ADR-0021 `:6`, ADR-Index, Roadmap, slice-050,
  `spec/spezifikation.md:565`/`:1153`) löst auf; die in R1 gemeldeten
  8× anchor-missing sind weg.
- **Restpunkt:** die in R1 mitgeforderte „Versions-Bump"-Teilleistung ist
  **unvollständig** — der Dokument-Header steht weiter auf 0.29.0
  (`spec/lastenheft.md:3`), während §7/CHANGELOG/ADR 0.30.0 behaupten. Als
  eigenständiger Befund MEDIUM-1 unten geführt.

### R1-HIGH-1 — falsche Grün-/Fertig-Attestierung → **AUFGELÖST (mit Restpunkt)**

- doc-check ist real grün (s. o.); die CHANGELOG-Aussage „Lastenheft 0.30.0"
  (`CHANGELOG.md:24`) ist jetzt durch eine existierende §3-Sektion gedeckt.
- DoD: die Punkte Spec/Code/Config-Ausgaben/Dogfood sind `[x]`
  (`slice-050:59,65,70,73`); der Punkt make-gates/CHANGELOG/Closure bleibt
  korrekt `[ ]` (`slice-050:76`) — die Closure (Move nach `done/`,
  Roadmap-Flip) ist tatsächlich noch offen (Slice liegt untrackt in
  `in-progress/`).
- ADR-0021-Geschichte „make gates grün. Status Accepted" (`ADR-0021:77`) ist
  mit dem grünen doc-check und gocognit < Schwelle (s. Negativbefunde)
  vertretbar.
- **Restpunkt:** die DoD-Zeile `:59` und die ADR/CHANGELOG-Aussage „0.30.0"
  attestieren einen Versions-Bump, den der Header nicht hergibt → MEDIUM-1.

### R1-MEDIUM-1 — fehlende Tests für rangfreies Ziel, Gleichrang, Selbstverweis → **AUFGELÖST**

- Neuer Test `TestMatrixDownwardKanten`
  (`internal/hexagon/core/rules/matrix_test.go:320`) deckt die drei Zweige ab
  und trifft sie korrekt:
  - **(a) rangfreies Ziel** (`:329`): Quelle in `order`, Ziel Klassen-Mitglied
    ohne `order`-Treffer → Pfad `!ok2` in `downwardFinding` (`matrix.go:212`).
  - **(b) Gleichrang** (`:336`): ein Glob `spec/*.md` matcht beide Dateien →
    gleicher Rang → Pfad `si >= di` (`si == di`).
  - **(c) Selbstverweis** (`:343`): Datei → sich selbst → `si == di`.
  - Aufwärts (`si > di`) ist zusätzlich in `TestMatrixDownwardRichtung`
    (`:253`) geprüft (architecture→lastenheft).

### R1-INFO-1 — Invariante „nie beide Befunde an einer Kante" → **AUFGELÖST**

- `slice-050:84`–`:90` expliziert jetzt den Selbstregel-Vorbehalt: solange die
  Klassen-Paar-`rules` klassenübergreifend bleiben, kann eine Kante nicht beide
  auslösen; eine zusätzliche klasseninterne Selbstregel `{from: X, to: X,
  allow: false}` ließe an einer Abwärtskante beide feuern — „kein Defekt,
  sondern zwei verschiedene Verletzungen (Annahme aus dem Impl-Review
  expliziert)".

---

## Neue Findings (Deltas seit R1)

### MEDIUM-1 — Dokument-Header `Version` nicht auf 0.30.0 gebumpt (Rest des R1-BLOCKERs)

- **kategorie:** MEDIUM (Skill: Konsistenz-Lücke zwischen Quellen derselben
  Sache; zugleich Restleistung des R1-BLOCKERs / der R1-HIGH-Attestierung)
- **quelle:** Anforderungs-Anlege-Prozess (Versions-Bump-Konvention),
  Maintainability
- **pfad:** `spec/lastenheft.md:3` (`**Version:** 0.29.0`) gegen
  `spec/lastenheft.md:1106` (§7-Zeile 0.30.0), `CHANGELOG.md:24`
  („Lastenheft 0.30.0"), `ADR-0021:6`/Geschichte, `slice-050:59` (DoD
  „Versions-Bump 0.30.0")
- **befund:** Der Header trägt weiter `0.29.0`, während §7-Historie, CHANGELOG,
  ADR-0021 und die DoD-Abhakung übereinstimmend 0.30.0 behaupten. Die
  Bump-Konvention (über zehn Vorslices belegt, z. B. slice-049: Header
  0.28.0→0.29.0 **im selben** doc-first-Commit wie die §7-Zeile) ist hier nur
  halb ausgeführt: §7 ja, Header nein. Das kanonische Versionsfeld des
  Lead-Dokuments widerspricht vier Quellen, die den Bump als erledigt
  ausweisen.
- **verifizierbar:** nein (kein Gate liest den Header — doc-check grün; die
  `versions`-Konfig nimmt `spec/lastenheft.md` explizit aus, `.d-check.yml`
  `versions.exempt-paths`); belegbar per Diff Header gegen §7/CHANGELOG/ADR.

### LOW-1 — bereits veröffentlichte 0.29.0-Historiezeile beim Edit verfälscht

- **kategorie:** LOW (Skill: Doku-Drift)
- **quelle:** Maintainability
- **pfad:** `spec/lastenheft.md:1107`
- **befund:** Beim Einfügen der 0.30.0-Zeile wurde die unmittelbar folgende,
  bereits mit v0.29.0 freigegebene 0.29.0-Zeile mitgeändert und enthält nun die
  Dopplung „Content-Pin gegen Content-Pin gegen inhaltlichen Drift" (am HEAD:
  „Content-Pin gegen inhaltlichen Drift"). Eine eingefrorene Historie-Zeile
  wurde unbeabsichtigt korrumpiert; rein redaktioneller Schaden, kein
  Verhaltensbezug.
- **verifizierbar:** nein (kein Gate); belegbar per `git diff spec/lastenheft.md`
  (Hunk `@@ -1069 +1103 @@`).

---

## Negativbefunde (geprüft, ohne Befund)

- **suggest.go `renderHarnessMatrix` (`:453`–`:488`):** die geschichtete
  Klasse `spec-straten` emittiert `order` **und** `direction` gemeinsam in
  Blockform (`:471`/`:472`) — fail-closed-konsistent (nie `order` ohne
  `direction`). Die Inline-Kommentare (`# autoritativste Schicht zuerst`,
  `# Abwärtsverweis ⇒ matrix-downward`) sind YAML-gültige Zeilen-Kommentare und
  brechen das Decoden nicht. Die repo-aware-Auskommentierung (`pre`-Präfix,
  `:469`–`:472`) greift für **jede** der vier neuen Blockzeilen. Kein Befund.
- **gocognit `renderHarnessMatrix`:** kognitive Komplexität ≈ 8 gegen Schwelle
  `min-complexity: 20` (`.golangci.yml:70`) — kein Lint-Bruch. Kein Befund.
- **config_template.go `--print-config` (`:39`–`:44`):** das `order`/`direction`-
  Beispiel ist **vollständig auskommentiert** (innerhalb des ohnehin
  kommentierten matrix-Blocks), Schema-Keys decken sich mit `configyaml`/
  `matrix.go`; das Template dekodiert unverändert (kein aktives `order` ohne
  `direction`). Kein Befund.
- **Benutzerhandbuch §4.7 (`:317`–`:338`):** `order` (Glob-Liste, autoritativste
  Schicht zuerst), First-Match-Rang, transitiver Abwärtsverweis →
  `matrix-downward`, rangfrei = nicht geprüft, fail-closed (Exit 2),
  Default-aus = unverändert — deckt sich mit der Implementierung; das
  YAML-Beispiel (`:326`–`:333`) ist verhaltenskonform. Kein Befund.
- **Handbuch-Versions-Pins v0.29.0:** entsprechen der aktuell freigegebenen
  Version (`version.md#aktuell` = v0.29.0); der 0.30.0-Release samt Pin-Bump
  steht erst in der Closure/Release-Prep an (Muster wie slice-049). Vor Release
  korrekt, `versions`-Gate grün. Kein Befund.
- **Lastenheft-Sektionstext DC-FA-MTX-002 (`:648`–`:678`):** vollständig und
  widerspruchsfrei zur Spezifikation Schritt 5 (`spec/spezifikation.md:564`),
  zu ADR-0021 und zur Implementierung (`matrix.go`/`configyaml.go`); keine
  Dopplung/Konflikt im Sektionstext selbst. Kein Befund.
- **Kohärenz Code ↔ Spec ↔ ADR ↔ Dogfood:** `matrix.go downwardFinding`
  (`:205`) deckt Schritt 5 exakt (gleiche Klasse, beidseitig rangbehaftet,
  `si >= di` schweigt, rangfrei/klassenübergreifend ausgenommen, Meldung nennt
  beide Ränge); `validateMatrixDirection` (`configyaml.go:417`) erzwingt die
  fail-closed-Kopplung; `model.DirectionNoDownward` (`config.go`) und das
  Doctor-Mapping (`diagnose.go`, `AllReasons` + `reasonTexts`) sind ergänzt.
  Kein Befund.
- **Dogfood nicht vacuously grün:** `.d-check.yml` `spec-straten` trägt
  `order`/`direction`; `spec/spezifikation.md` (Rang 1) verlinkt 75× aufwärts
  auf `lastenheft.md` (Rang 0) → korrekt schweigend, das Lastenheft hat **keinen**
  Abwärtslink auf spezifikation/architecture. Die Regel wird real ausgeübt und
  ist im grünen Lauf nicht bug-maskiert. Kein Befund.

---

## Kategorie-Summary

| Severity | Anzahl |
| --- | --- |
| BLOCKER | 0 |
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

(Status R1: BLOCKER/HIGH/MEDIUM/INFO je 1 — alle aufgelöst; das MEDIUM-1 hier
ist der unvollständig gebliebene Versions-Bump-Rest des R1-BLOCKERs.)

---

## Verdikt

**Bedingt mergebar — ein MEDIUM vorher schließen.** Der dispositive R1-BLOCKER
ist substanziell aufgelöst: die Anforderung existiert als kanonische
§3-Sektion, die Anker-Kette löst auf, `make doc-check` ist **grün** (Exit 0,
0 Befunde). Feature-Substanz, Tests (inkl. der drei R1-MEDIUM-Zweige),
fail-closed-Config, Determinismus, suggest/print-config/Handbuch und die
Dogfood-Konfig sind kohärent und korrekt. Es bleibt **ein** sauber benennbarer
Rest: der Lastenheft-Header steht weiter auf 0.29.0, obwohl §7/CHANGELOG/ADR/DoD
0.30.0 attestieren (MEDIUM-1) — eine Konsistenz-Lücke, die kein Gate fängt und
die genau die in R1 geforderte Versions-Bump-Teilleistung offen lässt. LOW-1
(verfälschte 0.29.0-Historiezeile) gleich miterledigen. Kein BLOCKER/HIGH mehr;
nach dem Header-Bump ist der Slice merge- (und closure-)reif.

### Übergabe (Fix-Hinweise, außerhalb der Befund-Felder)

- MEDIUM-1: `spec/lastenheft.md:3` `**Version:** 0.29.0` → `0.30.0` (deckt sich
  dann mit §7/CHANGELOG/ADR/DoD).
- LOW-1: in `spec/lastenheft.md:1107` die Dopplung „Content-Pin gegen
  Content-Pin gegen" auf „Content-Pin gegen" zurücksetzen (HEAD-Stand der Zeile).
