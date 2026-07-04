# Slice slice-060: `--doctor`-Klartext-Vollständigkeit (AllReasons ↔ §4)

**Status:** in-progress (welle-49-doctor-klartexte).

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
  | `planning-drift` | Roadmap-Aktiv-Status und Slice-Bestand inkonsistent (oder Roadmap/Überschrift fehlt — fail-closed) |

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
- [ ] **Belege/Prozess:** `make gates`/`make ci` grün; unabhängiges Review
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
[slice-059](../done/slice-059-tracked-modul.md) §7 („Kandidat für einen
Hygiene-Slice samt Vollständigkeits-Verriegelung §4 ↔ `AllReasons`");
Nutzer-Auftrag 2026-07-04.

## 6. Sub-Area-Modus-Begründung

GF („Doc führt, Code folgt" — der Doc-Vertrag besteht seit Lastenheft
0.15.0/0.17.0; der Code wird nachgezogen und verriegelt). Kein neuer
Adapter, keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

*(bei Closure: Umsetzung, Belege, Lerneintrag, Steering-Loop-Eintrag.)*
