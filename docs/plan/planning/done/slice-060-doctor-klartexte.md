# Slice slice-060: `--doctor`-Klartext-Vollständigkeit (AllReasons ↔ §4)

**Status:** done (welle-49-doctor-klartexte, Closure 2026-07-04).

**Welle:** welle-49-doctor-klartexte (Trigger: Bestands-Folgepunkt aus dem
slice-059-Closure-Lerneintrag; Nutzer-Auftrag 2026-07-04, die offenen
Hygiene-Punkte anzugehen).

**Bezug:** Defekt-Fix gegen
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
bzw. Spezifikation
[§`DC-FA-CLI-007.a`](../../../../spec/spezifikation.md#dc-fa-cli-007a--diagnose-modus)
Schritt 3 („für jeden Grund-Code aus §4
genau ein Eintrag, abgesichert durch eine Vollständigkeits-Prüfung gegen die
Reason-Konstanten"). **Kein Change Request** (das Lastenheft bleibt unberührt —
der Vertrag besteht und wird erfüllt statt geändert), **kein ADR** (keine
Architektur-Entscheidung; Mapping-Ergänzung plus Test im bestehenden Schnitt).

**Autor:** pt9912. **Datum:** 2026-07-04.

---

## 1. Ziel

`--doctor` übersetzt jeden Grund-Code über ein festes Mapping in einen
Klartext; die Spezifikation fordert je §4-Grund-Code genau einen Eintrag.
Die kanonische Liste `AllReasons()`
(`internal/hexagon/core/app/diagnose.go`) hinkt seit v0.25 **sieben** Codes
hinterher: `diagram-id-undefined`, `version-stale`, `link-stale`,
`core-drift`, `core-drift-vcs`, `commit-untraceable`, `planning-drift`. Für
Befunde dieser Module zeigt `--doctor` den rohen Code statt des Klartexts
(fail-safe-Pfad von `ReasonText`), und die maschinenlesbare Diagnose
(`--doctor --json`/`--yaml`) trägt den Code im `reasonText`-Feld. **Neu:**
die sieben Klartexte plus eine beidseitige **Verriegelung** `AllReasons()` ↔
§4-Grund-Code-Tabelle der Spezifikation — der Pfad, auf dem die Lücke
siebenmal unbemerkt wuchs (hand-gepflegte Liste ohne Gegen-Autorität), ist
damit zu.

## 2. Entscheidungen

- **Defekt-Fix, kein CR/ADR** (s. Bezug). Die Grund-Codes selbst sind in §4
  vollständig dokumentiert — doc-first hat gehalten; nur das Klartext-Mapping
  blieb zurück.
- **Verriegelung gegen §4 (doc ↔ code), nicht nur gegen Konstanten:** der
  bestehende Deckungs-Test
  (`internal/hexagon/core/app/diagnose_test.go`) verriegelt
  `reasonTexts` ↔ `AllReasons`; die Lücke war die hand-gepflegte
  `AllReasons`-Liste selbst. Der neue Test parst die Grund-Code-Tabelle aus
  `spec/spezifikation.md` §4 (erste Spalte, Backtick-Code) und vergleicht die
  Mengen **beidseitig** mit `AllReasons()`. Weil jedes neue Modul seinen
  Grund-Code doc-first zuerst in §4 einträgt, macht der Test das Vergessen im
  Mapping ab sofort rot — die Spezifikation wird Gegen-Autorität der Liste.
- **Fail-closed (Heading-Guard-Muster aus slice-057):** fehlende oder
  mehrdeutige §4-Überschrift oder leere Code-Tabelle ⇒ Test rot, kein
  stilles Grün.
- **Klartext-Stil wie Bestand** (deutsch, eine Zeile, benennt die Bedingung;
  an die §4-Bedingungsspalte angelehnt, ohne sie zu duplizieren):

  | Code | Klartext |
  |---|---|
  | `diagram-id-undefined` | Kennung im Diagramm-Fence ohne Definition in ihrer defined-in-Quelle |
  | `version-stale` | Versions-Pin weicht von der aktuellen Version ab |
  | `link-stale` | Ziel-Inhalt eines gepinnten Links weicht vom hinterlegten Content-Pin ab |
  | `core-drift` | Core einer gepinnten Datei weicht vom hinterlegten Immutabilitäts-Pin ab |
  | `core-drift-vcs` | Core einer immutablen Datei über die Commit-Range geändert, gelöscht/umbenannt oder mit unzulässigem Status-Übergang |
  | `commit-untraceable` | Commit-Message ohne Traceability-Kennung |
  | `planning-drift` | Roadmap-Aktiv-Status und Slice-Bestand inkonsistent (oder Roadmap/Überschrift fehlt bzw. mehrdeutig — fail-closed) |

- **Kein Verhaltens-Delta außerhalb der Diagnose:** Default-/JSON-/YAML-
  Befundausgabe, Exit-Codes und `--repair` bleiben unberührt
  (`FixCandidateFor` liefert weiter nur für `id-unlinked` einen Kandidaten).
- **Mutations-Beleg je Verriegelung (R3-Lehre slice-057 — jede Probe
  verriegelt genau ihren Guard):** (a) ein Code-**Paar** aus `AllReasons()`
  **und** `reasonTexts()` entfernt (der historische Fehlermodus: das Paar
  bleibt in sich konsistent, nur die Liste hinkt der Spec hinterher) ⇒
  **genau** der neue §4-Test rot, der bestehende Deckungs-Test bleibt grün;
  (b) ein Klartext nur aus `reasonTexts()` entfernt ⇒ der bestehende
  Deckungs-Test rot; (c) §4-Überschrift für den Parser unauffindbar ⇒ der
  fail-closed-Zweig schlägt an (kein stilles Grün bei leerer Menge).
- **Release als Patch v0.37.1** (SemVer: Fehlerbehebung am bestehenden
  Vertrag, keine neue Funktionalität, keine neue Config-Surface) — erster
  Patch-Release des Repos; Release-Prep nach
  [`releasing.md`](../../../user/releasing.md) inklusive des neuen
  Operations-Checklisten-Punkts.

## 3. Definition of Done

- [x] **Code:** `AllReasons()` und `reasonTexts()` in
  `internal/hexagon/core/app/diagnose.go` um die sieben Codes/Klartexte
  ergänzt (über die Reason-Konstanten aus `model`/`rules`, keine
  String-Literale).
- [x] **Verriegelung:** neuer Test parst die §4-Grund-Code-Tabelle der
  Spezifikation und vergleicht beidseitig mit `AllReasons()`; fail-closed
  bei fehlender/mehrdeutiger Überschrift oder leerer Tabelle; der
  bestehende Deckungs-Test `reasonTexts` ↔ `AllReasons` bleibt.
- [x] **Mutations-Belege** (a)–(c) aus §2 erbracht und im Closure
  dokumentiert.
- [x] **Beleg-Lauf:** `--doctor` an einer Probe mit einem der sieben Codes
  zeigt den Klartext statt des rohen Codes; `--doctor --json` trägt
  `reasonText` ≠ Code (Vorher/Nachher).
- [x] **Belege/Prozess:** `make gates`/`make ci` grün; unabhängiges Review
  vor Closure; CHANGELOG (Fixed); release-prep v0.37.1; Closure-Move nach
  `done/` + Roadmap-Flip
  ([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)),
  Closure-Body als Folge-Commit; Release v0.37.1 auf GHCR +
  Digest-Backfill.

## 4. Risiken / offene Punkte

- **Spec-Zugriff im Test:** `make test` läuft im Builder-Image mit vollem
  Build-Kontext (keine `.dockerignore`-Ausnahme) — `spec/spezifikation.md`
  liegt vom Paket aus unter `../../../../spec/spezifikation.md`. Erster
  Test, der die echte Spec liest (bisher nur synthetische Fixtures):
  bewusst, denn genau die Kopplung Doc ↔ Code ist der Zweck; deterministisch
  (gleicher Baum ⇒ gleiches Ergebnis,
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)-konform).
- **Format-Kopplung an die §4-Tabelle:** ändert sich das Tabellen-Format,
  wird der Test rot (fail-closed, gewollt) — die Reparatur ist dann eine
  bewusste Parser-Anpassung, kein stilles Grün.
- **Wording der Klartexte:** Review-Punkt (Endnutzer-Verständlichkeit,
  keine Duplikation der §4-Bedingungsspalte).

## 5. Trigger

Bestands-Folgepunkt aus dem Closure-Lerneintrag von
[slice-059](welle-48/slice-059-tracked-modul.md) §7 („Kandidat für einen
Hygiene-Slice samt Vollständigkeits-Verriegelung §4 ↔ `AllReasons`");
Nutzer-Auftrag 2026-07-04.

## 6. Sub-Area-Modus-Begründung

GF („Doc führt, Code folgt" — der Doc-Vertrag besteht seit Lastenheft
0.15.0/0.17.0; der Code wird nachgezogen und verriegelt). Kein neuer
Adapter, keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** Defekt-Fix wie geplant: `AllReasons()` und `reasonTexts()`
(`internal/hexagon/core/app/diagnose.go`) tragen die sieben seit v0.25
fehlenden Grund-Codes (`diagram-id-undefined`, `version-stale`, `link-stale`,
`core-drift`, `core-drift-vcs`, `commit-untraceable`, `planning-drift`) über
die Reason-Konstanten; `--doctor` zeigt für diese Module den Klartext statt
des rohen Codes, `--doctor --json`/`--yaml` trägt ihn im `reasonText`-Feld.
Der neue Test `TestAllReasonsDeckungGegenSpezifikationGrundCodes`
(`internal/hexagon/core/app/diagnose_test.go`) verriegelt `AllReasons()`
beidseitig gegen die §4-Grund-Code-Tabelle der Spezifikation — fail-closed
bei unlesbarer Spec, fehlender/mehrdeutiger Überschrift, leerer Tabelle und
(seit R1-LOW-1) jeder Tabellen-Body-Zeile ohne Backtick-Code. Kein CR, kein
ADR, keine neue Config-Surface. Commit-Kette: doc-first `1e5cf41` → feat
`f1017f8` → R1 `7e64fb5` → release-prep `7053583` → closure-move `a26915e`
→ closure-body (dieser Commit) → digest-backfill.

**Belege.**
- `make gates` und `make ci` **grün** (doc-check 180/0, lint, test,
  arch-check via a-check, coverage-gate, semgrep 0/55, gate-consistency,
  planning-check; image-test nativ == Container byte-identisch).
- **Vier Mutations-Belege** (R3-Lehre slice-057, jede Probe verriegelt genau
  ihren Guard): (a) `link-stale`-**Paar** aus `AllReasons()`+`reasonTexts()`
  entfernt (historischer Fehlermodus) ⇒ **genau** der neue §4-Test rot
  („§4-Grund-Code "link-stale" fehlt in AllReasons()"), Deckungs-Test grün;
  (b) Klartext allein entfernt ⇒ genau der Deckungs-Test rot (drei
  Fehlerzeilen); (c) Überschrift für den Parser unauffindbar ⇒
  fail-closed-Fatal statt stillem Grün; (d) §4-Zeile ohne Backticks ⇒
  lauter Fatal mit Zeilen-Zitat (R1-LOW-1-Guard).
- **Beleg-Lauf** (planning-drift-Probe, fehlende Roadmap): **v0.37.0-Image**
  zeigt `Z. 1 · planning-drift [planning]` (roher Code), **lokaler Build**
  zeigt den Klartext; `--doctor --json` trägt `reasonText` ≠ Code.
- **Unabhängiges Review R1**
  ([Report](../../../reviews/2026-07-04-slice-060-doctor-klartexte-r1.md)):
  **ACCEPT**, 0 HIGH/0 MEDIUM/2 LOW/1 INFO — alle eingearbeitet (LOW-1
  Zeilen-Format-Guard + Mutations-Beleg (d), LOW-2 planning-drift-Klartext
  nennt „mehrdeutig", INFO-1 Ein-Tabellen-Annahme dokumentiert).
- Release **v0.37.1** auf GHCR (erster **Patch**-Release; Tag auf dem
  Closure-Stand `5f1420d`, Release-Run 28695142320 grün, Push-CI
  28695142147 grün), Digest-Pin
  `ghcr.io/pt9912/d-check@sha256:3bbdb19bb73200fa37e30eff961cd5429a44e9e945fff3fb65ba7dc4b3cd88dd`.

**Lerneintrag.** (1) Die R1-Beobachtung trägt über den Slice hinaus: ein
Mengen-Verriegelungs-Test braucht neben dem Leere-Menge-Guard auch einen
**Zeilen-Format-Guard**, sonst kehrt der historische Fehlermodus über eine
einzelne still ausgelassene Zeile zurück. (2) Der Deckungs-Test von
slice-025 war korrekt und trotzdem machtlos: er prüfte das **Paar**
`reasonTexts` ↔ `AllReasons`, nicht die Liste gegen ihre Autorität — eine
hand-gepflegte kanonische Liste braucht eine **Gegen-Autorität außerhalb
ihrer selbst** (hier: die doc-first gepflegte §4-Tabelle). Steering-Loop:
Feedback (siebenfache stille Lücke, slice-059-Lerneintrag) wurde als
Sensor-Verschärfung verkörpert, nicht nur behoben.
