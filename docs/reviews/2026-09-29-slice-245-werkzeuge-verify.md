# Verifikations-Report — slice-245 (DoD)

**Verifikations-Art:** DoD/Plan-Konformität (Modul 11 — „Bauen wir es richtig?"; *nicht* der Diff-Review — der steht im R1-Report)

**Gegenstand:** `slice-245` (welle-91) — `tools/harness`-Werkzeuge aus ai-harness-init evaluieren

**Range:** `e6c51c14..HEAD` — `aad1c9e7` (Adoption: 3 Skripte), `a0d53819` (slice-mv Identität-Fallback + DONE-Unterordner), `bff7b762` (slice-247 geschnitten), `8637f08a` (Wächter als Einzellauf-Target), `2402c330` (R1-Report), `09aeb4d5` (R1-Befunde F-1 bis F-6, F-8 eingearbeitet)

**Plan:** `docs/plan/planning/in-progress/slice-245-werkzeuge-evaluieren.md`

**Modell-ID:** glm-5.3-flash · **Datum:** 2026-09-29

---

## DoD §2 — je Punkt geprüft

### DoD 1 — Je Werkzeug (5) eine belegte Entscheidung: **ERFÜLLT**

| Werkzeug | Entscheidung | Beleg (selbst geprüft) |
|---|---|---|
| slice-mv.sh | ADOPTIERT | Target `Makefile:408` + d-check-Anpassung `Makefile:404-407` (`SLICE_MV_DONE_UNTERORDNER`); Funktion in drei Scratch-Repos gemessen (siehe R1-Einarbeitung F-1) |
| history-range-guard.sh | ADOPTIERT | Präludium `Makefile:368` (trace-check) + `Makefile:382` (adr-check); Einzellauf-Target `Makefile:414`; vier Proben im Voll-Klon, zwei im Shallow, Verdrahtungswirkung gemessen (DoD 2) |
| selbstpruefung.sh | ADOPTIERT | Target `Makefile:411`; **Vollauf** (`make selbstpruefung`, Default-Konfiguration inkl. `make gates` im Klon) Exit 0 — alle beobachtbaren Ausgänge in einem Lauf |
| e2e-abdeckung | ABGELEHNT | Gegenstand in der Schwester gelesen (ai-harness-init): `docs/user/e2e-abdeckung.md` + `harness/tools/e2e-abdeckung.sh` gelesen: das Werkzeug liest die Stufen von `harness/tools/full-smoke.sh` aus dem Quelltext — d-check führt kein solches Voll-E2E (grep über `tools/`/`Makefile`: kein Treffer). Begründung trägt |
| traeger-fetch | ABGELEHNT | Skriptkopf gelesen: legt den ai-harness-init-Träger per Fetch aus dem gepinnten Release ab; d-check ist kein Consumer-Repo der Schwester, verteilt sich selbst als GHCR-Image; kein Pendant im Makefile/Gate-Index. Begründung trägt |

Der Vollauf der Selbstprüfung im Detail (alles aus der Ausgabe, nicht behauptet): frischer Klon trägt lokal keinen `core.hooksPath` → `make hooks` setzt `.githooks` → Träger-Abgleich `-ef` bestätigt, dass `.githooks/commit-msg` die gerufene Datei ist → Commit ohne Kennung fällt (Exit 1, HEAD unverändert), Commit mit Kennung geht durch (HEAD trägt die Message) → `make gates` im Klon Exit 0. Rot- und Grün-Fall in einem Lauf — der Träger ist belegt, nicht nur durchgelassen.

### DoD 2 — stille-Grün-Verifikation am eigenen Adapter: **ERFÜLLT** (selbst gemessen)

Eigener Shallow-Klon (`git clone --depth 1 file:///Development/d-check`, `is-shallow = true`, Tiefe 1):

| Probe | Gegenstand | Ergebnis |
|---|---|---|
| K3 | vcs-Modul, `--range HEAD..HEAD` | `d-check: 909 Datei(en) geprüft, 0 Befund(e)`, **Exit 0 — die Lücke, reproduziert** |
| K4 | commits-Modul, dieselbe Range | `error: Range-Basis-Vorfahren nicht lesbar: object not found`, **Exit 2 — laut** |
| K5 | vcs-Modul, `--range HEAD~1..HEAD` | `error: Range-Basis "HEAD~1" nicht auflösbar`, **Exit 2 — laut** |
| K1 | Wächter, `HEAD..HEAD` im Shallow | „aufloesbar, aber LEER (0 Commits)“, **Exit 1** |
| K2 | Wächter, `HEAD~1..HEAD` im Shallow | „NICHT aufloesbar (Basis fehlt im Klon?)“, **Exit 2** |

Voll-Klon-Gegenproben des Wächters: `HEAD~1..HEAD` → Exit 0 („1 Commit(s) — OK“) · `HEAD..HEAD` → Exit 1 · `0000000..HEAD` → Exit 2 · `--staged` ohne Änderung → Exit 0 mit Meldung.

**Aus dem richtigen Grund** (Modul 11, Bewusstes Brechen): `make trace-check RANGE=HEAD..HEAD` im Shallow-Klon → **Exit 2** mit der Wächter-Meldung — dieselbe Range, die K3 ohne Wächter **still grün** zeigte, fällt mit Verdrahtung laut. `make trace-check RANGE=HEAD~1..HEAD` im Shallow → Exit 2 am Wächter; positive Kontrolle `make trace-check RANGE=HEAD~1..HEAD` im Voll-Klon → Exit 0. Und der Hook-Modus: `make trace-check MSGFILE=…` im Voll-Klon (2203 Commits) — ohne Kennung Exit 2 mit `commit-untraceable — Commit-Message ohne DC-/ADR-/MR-/slice-ID` bei gleichzeitig grünem Wächter, mit Kennung Exit 0: die Selbstprüfung fällt **am Modul**, nicht am Vorlauf-Wächter.

### DoD 3 — `make gates` grün: **ERFÜLLT** (selbst gefahren)

Frischer Lauf, Exit 0; die Zusammenfassung in der Ausgabe:

- `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`
- `doc-check`/`planning-check`: 909 Datei(en) geprüft, 0 Befund(e) · `coverage-gate`: 94,60 % ≥ 93 % · `semgrep`: 55 Regeln, 0 findings
- `git status --porcelain` danach: **leer** — seit dem Lauf liegt nichts im Baum.
- Zweitbeleg: `make gates` im frischen Vollklon (aus dem Selbstprüfung-Vollauf) Exit 0 — der Stand ist auch klon-seitig grün.

---

## R1-Einarbeitung — Stichproben am Code

**F-1 (MEDIUM, substantiell) — in beide Richtungen gemessen.** Drei Scratch-Repos mit d-check-Layout (`in-progress/`, Links mit Dateiendung), je zwei Commits erwartet:

| Fall | outgoing nach dem Move | löst auf? |
|---|---|---|
| Post-Fix + `SLICE_MV_DONE_UNTERORDNER=welle-91` | `](../../in-progress/slice-902-nachbar.md)` | ja |
| Post-Fix, flaches `done/` | `](../in-progress/slice-902-nachbar.md)` | ja |
| **Pre-Fix (`8637f08a`) + Unterordner** | `](../in-progress/slice-902-nachbar.md)` | **nein — F-1 reproduziert** |

Eingehend in allen Fällen korrekt (`](done/welle-91/slice-901-scratch.md)` mit Präfix, `](../done/welle-91/slice-901-scratch.md)` praefixlos), Zwei-Commit-Folge (reiner Move, dann Inhalt) eingehalten. Der Pre-Fix-Lauf zeigt die behauptete Lücke exakt — die Prüfung trägt, weil sie ohne den Fix rot läuft.

*Testnotiz:* die ausgehende Ersetzung feuert nur für Ziele mit Dateiendung; ein praefixloses Ziel ohne `.md` löst als Pfad nie auf und wird vom `[ -f ]`-Filter bewusst ausgenommen (Grenze „kein Rateversuch“, Skriptkopf). Erstversuch mit endungslosem Link fiel darum durch — Testdatensatz, kein Befund.

**F-2 (MEDIUM):** `grep slice-2 harness/README.md` → leer. Die drei Werkzeuge-Zeilen tragen keine Planungs-Verweise mehr.

**F-3 (MEDIUM):** keine nackte `(slice-NNN)`-Form mehr in Makefile/Skripten; die verbliebenen Token stecken in der `seit slice-245`-Ankerform (`Makefile:404/408/411/414`, `slice-mv.sh:99`, `selbstpruefung.sh:3`) — die Form, die F-3 selbst als Korrektur-Richtung nennt (Präzedenz `Makefile:313`, `seit slice-215`).

**F-6 (LOW):** `harness/sensors/trace-check.md` und `adr-check.md` tragen je den Abschnitt „Vorlauf-Wächter“ (im Diff von `09aeb4d5` verifiziert).

---

## Plan-vs-Code-Diff

- **Der Stand liefert genau, was der Plan verspricht:** drei Adoptate mit Target, zwei begründete Ablehnungen. Die Abgrenzung „Kein Produkt-Fix an vcs/commits ohne Anforderung“ ist eingehalten — kein Produkt-Code im Range; statt dessen ist [slice-247](../plan/planning/open/slice-247-vcs-leere-range-stilles-gruen.md) in `open/` geschnitten und referenziert die Messwerte (Probe A/B), die ich an eigenen Läufen bestätigt habe.
- **Mehr als der Plan:** nur die prozessverlangten Nebentreffer — `harness/README.md`-Werkzeugezeilen (Gate-Index-Disziplin; `make gate-consistency` in beide Richtungen grün), Sensor-Doku-Paragraphen (R1-F-6), R1-Report, slice-247-Plan. Keine Spec-Berührung.
- **Nichts versprochen und fehlend:** §7 Closure-Notiz steht auf `—` — laut Plan §5 erst vor dem `git mv` nach `done/` zu füllen; die §6-Risiken sind beide durch Belege entlastet (Host-bash-Klasse im R1-Negativbefund; die Shallow-Klon-Mechanik ist mit `file://`-Klon netzfrei gemessen).

## Beobachtungen (nicht blockierend, closure-seitig)

1. `slice-mv.sh:31-39` trägt die Emissions-Prosa der Schwester („Diese Datei wird bei jedem Bootstrap kanonisch neu geschrieben … das mitemittierte Make-Fragment“) — in d-check ist die Datei eine adoptierte Kopie, kein Emissionsprodukt. §3.7-naher Kommentar-Drift, kein DoD-Verstoß; beim Closure-Zug oder als Form-Fix mitnehmen.
2. `slice-mv.sh:312`: die Abschluss-Meldung hardcodet `auf ../$from/ umgehaengt` — im Unterordner-Fall eine Ebene zu niedrig; die Ersetzung selbst nutzt korrekt `$stufen`, nur die Mensch-Meldung weicht.
3. Beim Closure-Zug: slice-247 referenziert slice-245 (`../in-progress/slice-245-werkzeuge-evaluieren.md`, Zeile 12) — der Move nach `done/welle-91/` bricht den Link. Das ist die erste produktive Nutzung von `make slice-mv SLICE=slice-245-werkzeuge-evaluieren TO=done SLICE_MV_DONE_UNTERORDNER=welle-91`, deren eingehende Verweis-Reparatur genau dafür gebaut ist; `make doc-check` sieht den Fall danach.
4. Die DoD-Haken in §2 stehen offen — korrekt für `in-progress/` (MR-056); beim Closure-Zug abhaken.
5. Im Hook-Modus läuft der Wächter mit `HEAD~1..HEAD` (`Makefile:368`) — in einem Shallow-Klon (nur HEAD) blockiert er jeden Commit, laut und mit fetch-depth-Hinweis (fail-closed, kein stiller Zustand); im realen Hook-Szenario (Vollklon) folgenlos. Grenze benannt, kein Befund.

## Verdikt

**DoD vollständig erfüllt.** Alle drei DoD-Punkte sind gegen eigenständig erhobene Belege bestätigt — fünf belegte Werkzeug-Entscheidungen (drei adoptiert und funktionell gemessen, zwei mit gegen die Quellen geprüfter Ablehnung), die stille-Grün-Lücke am eigenen Adapter reproduziert und mit Verdrahtung aus dem richtigen Grund rot, `make gates` frisch und grün mit leerem Baum. R1-MEDIUMs sind wirksam eingearbeitet (F-1 empirisch in beide Richtungen belegt, F-2/F-3 formlich sauber). Der Slice ist DoD-seitig closure-fähig; offen sind nur die closure-seitigen Schritte (DoD-Haken, Volltext der §7-Notiz, Verweis-Nachzug bei slice-247, `verify-closure-notes` beim `fullbuild`).
