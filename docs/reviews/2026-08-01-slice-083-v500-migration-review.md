# Review-Report: slice-083 (Regelwerk-Migration v1.4.0 → v5.0.0, Delta-Analyse + Etappen-Schnitt) — 2026-08-01

**Review-Art:** Plan-Review (Analyse-/Etappen-Schnitt gegen die kanonischen Quellen), unabhängiger Frischkontext.

**Gegenstand:** `docs/plan/planning/open/slice-083-regelwerk-v500-migration-analyse.md` (uncommittet).

**Skill:** `.harness/skills/reviewer.md` v1.2.0

**Modell:** claude-opus-4-8 · **Datum:** 2026-08-01

**Eingangs-Kontext (geprüfte Verträge/Quellen):**

- v5.0.0-Regelwerk-Bundle (Scratchpad `v5-regelwerk/`): `regelwerk/` (8 grundlagen + 17 Module + README) und `templates/`.
- Vendorter v1.4.0-Von-Stand: `.harness/baseline/v1.4.0/regelwerk/`.
- Kurs-Repo das Schwester-Repo `ai-harness-course` (Tags v1.4.0/v2.0.0/v3.5.2/v4.0.0/v4.1.0/v5.0.0, `CHANGELOG.md`).
- d-check-Ist: `harness/conventions.md`, `AGENTS.md`, `harness/README.md`, `.harness/skills/reviewer.md`, `.d-check.yml`, `tools/harness/fetch-baseline-cache.sh`.

**Methode:** jede Faktenbehauptung gegen die Quelle nachgemessen (git-Tag-Diffs, Byte-/Zeilenzählung, Repo-Greps), die Etappen gegen die tatsächlichen Modul-Fähigkeiten (`links`/`anchors`/`ids`/`codepaths`/`vcs`) geprüft. REFUTE nur mit Zitat.

---

## Findings

### F-1 (MEDIUM) — Bruch 1 (retired `agents-regelwerk.md` + entfallener `konventionen.md`-Quellzeiger) hat keinen eigenen Etappe-Schritt; Etappe A ist als „vollständig spezifiziert" ausgewiesen

- **kategorie:** MEDIUM
- **quelle:** §2.2 Bruch 1 / §2.7 Etappe A (Schritt 4) / DoD-Zeile „A vollständig spezifiziert"
- **pfad:** `docs/plan/planning/open/slice-083-...md:100` (Bruch 1), `:334` (A4), `:449` (DoD); Belege `harness/conventions.md:39` + `:37`, `harness/README.md:42`
- **befund:** §2.2 Bruch 1 nennt `agents-regelwerk.md` als in v2.0.0 retired und listet zwei Live-Fundstellen (`harness/README.md:42` §Guides, `harness/conventions.md:39` §Adoptierte). Kein Etappe-Schritt behandelt diese Retirement; A4 lautet „die sechs Live-Fundstellen (§2.2 #4) auf den neuen Pfad ziehen" — #4 sind ausdrücklich die pin-gebundenen Vendor-Pfad-Verweise, nicht die Adoptions-Quellen. Eine wörtliche Ausführung von A4 zieht die `agents-regelwerk.md`-Raw-URL auf einen v5.0.0-Pfad, an dem die Datei seit v2.0.0 nicht mehr existiert (404). Zusätzlich trägt derselbe §Adoptierte-Block (`harness/conventions.md:37`) den Klartext-Zeiger „Konventionen: `kurs/de/grundlagen/konventionen.md`" — diese Quelldatei **entfällt in v5.0.0 ersatzlos** (CHANGELOG Welle 64, gleiche Klasse wie `grundlagen-konventionen.md`); auch sie wird von keiner Etappe als zu ersetzen benannt. Damit ist die als „vollständig spezifiziert" ausgewiesene Etappe A für Bruch 1 unvollständig.
- **verifizierbar:** ja — nach Ausführung von Etappe A zeigt `grep -n agents-regelwerk harness/README.md harness/conventions.md` weiterhin das retirete Asset als „adoptiertes Betriebsregelwerk"; ein `make gates` (links) gegen eine auf v5.0.0 umgehängte `agents-regelwerk.md`/`konventionen.md`-URL bzw. deren Auflösung schlägt an.

### F-2 (MEDIUM, grenzt an HIGH) — Etappe C bricht mit dem `conventions.md`-Split 52 repo-weite `MR`-Anker-Links, davon 12 in immutablen Accepted-ADRs; Schritt-8-Zusagen sind so nicht erreichbar

- **kategorie:** MEDIUM
- **quelle:** §2.3 (Index-statt-Inline) / §2.7 Etappe C (Schritte 2 + 8) / DC-FA-ANCH-001 / ADR-0016+ADR-0024 (Immutabilität)
- **pfad:** `docs/plan/planning/open/slice-083-...md:369` (C2), `:388` (C8); Belege `.d-check.yml:30-31` (ids `MR-\d{3}` target `harness/conventions.md`), `harness/conventions.md` (23× `### MR-NNN`-Heading)
- **befund:** C2 überführt die 23 `### MR-NNN`-Heading-Einträge aus `harness/conventions.md` in Einzeldateien `harness/conventions/MR-NNN-*.md`; die Anker `#mr-nnn-...` verlassen damit `conventions.md`. Repo-weit zeigen **52 Live-Links** aus **16 Dateien** auf `conventions.md#mr-xxx` (gemessen, ohne baseline/reviews/done/slice-083). **12 dieser Quelldateien sind Accepted-ADRs** (0019, 0021, 0022, 0023, 0024, 0026, 0027, 0028, 0029, 0030, 0031, 0046 — je `**Status:** Accepted`, vcs-gegated per `.d-check.yml:113`), die nach ADR-0016/ADR-0024 **nicht editiert** werden dürfen und deren Links folglich nicht retargetet werden können. C8 behauptet dennoch „`make gates` + `make adr-check` grün", „die neuen `conventions/`-Dateien erfüllen die `ids`-Linkpflicht (jeder `MR`-Verweis verlinkt)" und „keine Accepted-ADR wird inhaltlich berührt (die eingefrorene ADR bleibt via Tombstone aus A6)". A6 tombstoned aber nur die **eine** v1.4.0-Pfad-Referenz von ADR-0022, nicht die `MR`-Anker-Links der zwölf ADRs. „Jeder `MR`-Verweis verlinkt" adressiert die `link-policy`, nicht die Anker-Existenz (DC-FA-ANCH-001): nach dem Split lösen die 52 Links nicht mehr auf. Der nötige Umfang (Massen-Retarget der editierbaren Quellen AGENTS.md/README/spec/offene Slices, Massen-`ignore-refs`-Tombstone für die immutablen ADRs, Umstellung des `ids`-`target` von der Einzeldatei auf das Verzeichnis) ist in keiner Etappe benannt.
- **verifizierbar:** ja — nach C2 meldet `make gates` (anchors-Modul) ~52 `anchor-missing` auf `conventions.md#mr-xxx`, davon 12 in Accepted-ADRs; C8 (`make gates` grün) ist ohne die unbenannte Tombstone-/Retarget-Kampagne nicht erreichbar.

### F-3 (LOW) — §2.2 Bruch 6 sagt „in **acht** grundlagen-*-Dateien aufgesplittet"; die Quelle splittet `konventionen.md` in **sechs**

- **kategorie:** LOW
- **quelle:** §2.2 Bruch 6 vs. §2.1 vs. CHANGELOG Welle 63/64
- **pfad:** `docs/plan/planning/open/slice-083-...md:137-138`
- **befund:** Bruch 6 formuliert „die eine große Konventions-Datei ist in acht `grundlagen-*`-Dateien aufgesplittet". Gemessen: `konventionen.md` wird in **sechs** Dateien überführt (`grundlagen-begriffe`/`-bootstrap`/`-harness-dateien`/`-referenz-richtung`/`-source-precedence`/`-traceability`); die **acht** ist der grundlagen-Gesamtbestand (die sechs plus die schon zuvor bestehenden `-durchsetzungsschicht`/`-klassifikation`). §2.1 gibt es korrekt wieder (drei → acht, sechs neu). CHANGELOG Welle 63: „`grundlagen-konventionen.md` (1099 Zeilen) liegt auf sechs Seiten." Interne Zahl-Inkonsistenz, migrationsleitend unkritisch.
- **verifizierbar:** ja — Datei-Diff `.harness/baseline/v1.4.0/regelwerk` (3 grundlagen) vs. v5-Bundle (8 grundlagen, davon 6 aus `konventionen.md`).

### F-4 (LOW) — DoD sagt „Die **vier** Brüche"; §2.2 führt **sechs**

- **kategorie:** LOW
- **quelle:** §3 DoD vs. §2.2
- **pfad:** `docs/plan/planning/open/slice-083-...md:441` (DoD) vs. `:98` (§2.2 „Die sechs Brüche")
- **befund:** Die DoD-Checkbox „Die vier Brüche benannt und je an einer konkreten Fundstelle belegt (§2.2)" ist bei der Neubasierung 2026-08-01 (Brüche 5 + 6 ergänzt) nicht nachgezogen worden; §2.2 heißt „Die sechs Brüche" und listet sechs. Doku-Drift innerhalb des Slice.
- **verifizierbar:** ja — Textvergleich der beiden Zeilen.

### F-5 (LOW) — §2.3-„Fehlende Artefakte" übergeht das Beobachtungs-Register (`observations.md`), einen v5.0.0-Standard, den d-check nicht führt

- **kategorie:** LOW
- **quelle:** §2.3 (Aufzählung „Fehlende Artefakte") / §2.8 / Modul 6 §Beobachtungs-Register / grundlagen-harness-dateien §Verzeichniskonvention
- **pfad:** `docs/plan/planning/open/slice-083-...md:172-176` (fehlende Artefakte); Quelle `v5-regelwerk/regelwerk/modul-06-roadmap.md:77` + `grundlagen-harness-dateien.md:14`
- **befund:** §2.3 zählt als fehlende v5.0.0-Artefakte den `closure-note-reviewer`-Skill und die Review-Report-Kopf-Felder auf (beide bestätigt vorhanden im Bundle, in d-check absent). Nicht genannt: das **Beobachtungs-Register** `docs/plan/planning/observations.md` (der Steering-Loop-Zähler, `BEO-<NNN>`) — v5.0.0 führt es als Standard-Datei (`grundlagen-harness-dateien.md:14`) und als stehende Modul-6-Sektion, auf die die Wellen-Closure ihren Lese-Schritt stützt; in v1.4.0 existiert es noch nicht (`modul-06` dort: 0 Treffer), d-check hat es nicht (`ls docs/plan/planning/observations.md` → absent). Da die Analyse Modul 6 für die fünf Roadmap-Abschnitte + Fünf-Schritt-Closure bereits ausgewertet hat, ist die Auslassung innerhalb des gelesenen Materials; §2.8 disclaimt das volle Modul-Gegenlesen als Etappe B, was die Wirkung begrenzt.
- **verifizierbar:** ja — `ls docs/plan/planning/observations.md` (absent) gegen `grundlagen-harness-dateien.md:14`.

### F-6 (INFO) — Fork-Grenze: der Slice zitiert die **weite** Template-Fassung; `grundlagen-source-precedence.md` gibt eine **engere**

- **kategorie:** INFO
- **quelle:** §2.3/§2.4 (Fork-Grenze) / `templates/harness/conventions/MR-NNN-titel.template.md:18-19` vs. `grundlagen-source-precedence.md:79-88`
- **pfad:** `docs/plan/planning/open/slice-083-...md:213-216` (§2.3), `:247-253` (§2.4 Achse 1)
- **befund:** Der Slice formuliert „wer keine benannte Regel ersetzt, ist ein Fork" (wörtlich Template Z. 18-19) und klassifiziert die repo-eigenen additiven Einträge (Gate-Nachweis-Mechanik MR-004/005, Pointer-Disziplin MR-015) darüber als Forks. `grundlagen-source-precedence.md:79-82` fasst „Fork" enger: eine `MR`, die „die Baseline **pauschal für nicht anwendbar erklärt**, statt eine benannte Regel zu ersetzen" — ein additiver repo-lokaler Eintrag in einem von der Baseline nicht geregelten Bereich erklärt die Baseline nicht pauschal für unanwendbar. Unter der engeren Fassung wären diese Einträge nicht zwingend Forks. Der Slice deklariert die verbindliche Fork-Zuordnung ausdrücklich als Etappe-C-/Abnahme-Detail (§2.8), womit die Wirkung auf INFO begrenzt bleibt; die zwei divergierenden Baseline-Fassungen gehören aber beim Abnahme-Entscheid benannt.
- **verifizierbar:** ja — Zitatvergleich der beiden Quelltexte.

### F-7 (INFO) — „die 23 Adaptionen" zählt `MR-000` (Adoptions-Erklärung) als Adaption mit

- **kategorie:** INFO
- **quelle:** §2.3/§2.7 Etappe C („die 23 Adaptionen") / `grundlagen-source-precedence.md:85-88`
- **pfad:** `docs/plan/planning/open/slice-083-...md:204` + `:311`
- **befund:** `harness/conventions.md` führt MR-000..MR-022 = 23 Einträge. `grundlagen-source-precedence.md:87-88`: `MR-000` ist „weder Fork noch Adaption, sondern die Adoptions-Erklärung selbst". „Die 23 Adaptionen" zählt sie mit; sachlich sind es 22 Adaptionen plus die Adoptions-Erklärung. Rein terminologisch, ohne Migrations-Wirkung.
- **verifizierbar:** ja — `grep -c '^### MR-' harness/conventions.md` (23) gegen die Quell-Definition.

---

## Negativbefunde (je Prüfachse — geprüft, ohne blockierenden Befund)

**Achse 1 — Faktische Delta-Claims (§2.1/§2.2).** Nachgemessen und **bestätigt**: grundlagen 3 → 8 (`konventionen.md` entfällt, 6 Ersatzdateien); Module `03-lastenheft`→`03-spec` und `04-architektur-adrs`→`04-adrs` (v4.0.0-CHANGELOG-Bruchblock); Zeilen **4030 → 3851** (byte-genau: `cat *.md | wc -l` beider Bäume); Version→Welle **v4.0.0=Welle 62, v4.1.0=Welle 63 (additiver Split mit Wegweiser), v5.0.0=Welle 64** (README-`Stand:` je Tag); zwei Majors über v3.5.2 (=Welle 34), vier Major-Sprünge v1→v5, ≥14 Releases (tatsächlich 16). Einzige Ungenauigkeit: F-3 (Bruch-6-„acht").

**Achse 2 — Die sechs Brüche (§2.2).** Bruch 2 (self-contained Bundle, `lab-templates.zip` weg) und Bruch 3 (Skript bricht dreifach) **am Skript verifiziert**: `fetch-baseline-cache.sh:174` entpackt das self-contained Bundle nach `${baseline}/regelwerk` → Doppel-Verschachtelung, `:175` `sha256sum regelwerk/*.md` greift ins Leere; `:179` lädt `lab-templates.zip` (existiert nicht mehr) → `curl -f` bricht unter `set -e`; `:111-116` (`--check-latest`) erzeugt aus derselben Verschachtelung `authenticity=drift` (Fehl-Alarm). Bruch 4 **exakt**: sechs Live-Pfad-Fundstellen (AGENTS.md, harness/README.md, conventions.md, planning/README.md, roadmap.md, reviewer.md) + drei historische (ADR-0022, slice-080, slice-081) — file-gruppierter Grep deckungsgleich. Bruch 5 **bestätigt**: keine Live-d-check-Fundstelle nennt die alten Modul-Dateinamen (die 03/04-Treffer sind der Slice selbst + ADR-0004-Fehltreffer). Bruch 6 **bestätigt**: `grep grundlagen-konventionen` liefert live nur `reviewer.md:6+39`. **Befunde:** F-1 (Bruch 1 ohne Etappe-Schritt) und F-2 (ein vom Split verursachter Bruch, den §2.2 nicht enumeriert).

**Achse 3 — conventions.md-Struktur + neue MR-Felder (§2.3/§2.4/Etappe C).** Gegen `grundlagen-harness-dateien.md:114-169` und das MR-Template **bestätigt**: Index-statt-Inline, ein Eintrag je Datei, `conventions/done/` = Position (kein Status-Feld), Pflichtfeld `Ersetzt-Baseline-Regel` (genau eine, Anker-Link, sonst Fork — Template Z. 15-19), `Status: Accepted` (Template Z. 21), `Löst auf`/`Ausgelöst durch Baseline-Stand`, `(schärft …)`. Der Slice-Satz „Der Default ist die Verzeichnis-Form" ist **korrekt** (`grundlagen-harness-dateien.md:150-152` wörtlich; die „Form ist Wahl"-Freiheit gilt für Repos mit ~zwei Adaptionen, nicht für d-checks 23). Kein Über-/Fehlzitat. Nuance zur Fork-Grenze → F-6.

**Achse 4 — Fork-Reklassifizierung (§2.4).** Die „bleibt"-Einträge (MR-004/005 Gate-Nachweis-Mechanik, MR-015 Pointer-Disziplin) ersetzen tatsächlich keine benannte Baseline-Regel (am Eintragstext geprüft) — unter der Template-Fassung Forks, tragfähig. Keine übersehenen Fork-Kandidaten gefunden, die der Slice als „Adaption" fehlklassifiziert. Definitions-Divergenz → F-6.

**Achse 5 — MR-Bestands-Einschätzungen (§2.4-Tabelle).** Stichprobe **bestätigt**: **MR-001** „entfällt" — v4.0.0 macht drei Spec-Straten obligatorisch (`grundlagen-source-precedence.md:111-114`) und Source Precedence neun Ränge mit `spezifikation.md` als Rang 2 (CHANGELOG Welle 62), womit MR-001s Adaption zum Default wird. **MR-018** „entfällt vollständig" — im v5.0.0-Bundle kein Self-Hosting-/Producer-Abschnitt und kein `harness.mk` (grep: 0 Treffer); die Brücken-Begründung ist erloschen; `§Ein- vs. wiederkehrende Templates` (`templates/README.md:72-87`) bestätigt die fünf co-located Skelette. **MR-014** „entfällt" — die v5.0.0-„Form ist Wahl"-Freiheit ist auf `conventions.md` verengt (alle Fundstellen), während die Slice-/ADR-Form template-gegeben ist (`modul-04-adrs.md:57` „Die Form liefert die Vorlage"); die Basis von MR-014 („Baseline: Form ist Wahl", v1.3.0) trägt für Slice/ADR nicht mehr. Kein Urteil gegen die Quelle widerlegt.

**Achse 6 — Migrationsschritte (§2.7 A–D).** Reihenfolge A→B→C→D korrekt begründet (netzloser Boden zuerst). Reviewer-Retarget ist als A5 präsent (`konventionen.md`→`referenz-richtung.md` inkl. Anker-Verifikation). Das für die drei historischen Verweise vorgeschlagene `ignore-refs`/Tombstone ist **tauglich**: `links` und `anchors` honorieren das geteilte Top-Level-`ignore-refs`-Ventil (DC-FA-REF-001, `links.go:37-42`, `anchors.go:213-217`), unterdrücken damit `target-missing`/`anchor-missing` referenz-weit ohne Quelldatei-Ausschluss — auch für die immutable ADR-0022 (nur `.d-check.yml` wird angefasst, nicht der ADR). **Befunde:** F-1 (A unvollständig ggü. Bruch 1), F-2 (C unterspezifiziert ggü. der `MR`-Anker-Bruchflut).

**Achse 7 — Interne Konsistenz + Ehrlichkeit.** §2.8 ist überwiegend ehrlich (Modul-Delta, Entfall-Kandidaten, Closure-Log, Gate-Änderungen, Fork-Zuordnung als offen ausgewiesen). Der „Kostensenker"-Caveat (§2.1) ist **korrekt eingeschränkt**: a-check-Übertragbarkeit gilt nur bis zur v3.5.2-Untergrenze, der v4/v5-Sprung wird als flottenstand-abhängig (Etappe B) offen gehalten. **Befunde:** F-4 (DoD „vier" vs. „sechs"), F-5 (Beobachtungs-Register nicht in „fehlende Artefakte"), plus die Vollständigkeits-Überzeichnung, die F-1/F-2 belegen (A „vollständig spezifiziert" / C „Prozedur präzise" tragen die zwei benannten Lücken nicht).

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 2 | F-1, F-2 |
| LOW | 3 | F-3, F-4, F-5 |
| INFO | 2 | F-6, F-7 |

---

## Verdikt

**Nicht abnahmereif** (zwei offene MEDIUM). Die Delta-**Messung** der Analyse (§2.1/§2.2 Fakten, §2.3/§2.4/§2.5 Urteile) ist durchweg quellentreu und trägt — hier hält der Slice, was er behauptet. Blockierend sind zwei **Etappen-Lücken**, die die Migration real fehlleiten würden, weil sie in als „vollständig"/„präzise" ausgewiesenen Etappen liegen:

- **F-1:** Bruch 1 (retired `agents-regelwerk.md` + entfallener `konventionen.md`-Quellzeiger) hat keinen Schritt; A4s „auf den neuen Pfad ziehen" ist für retirete/entfallene Assets die falsche Operation.
- **F-2:** Etappe C bricht mit dem `conventions.md`-Split 52 repo-weite `MR`-Anker-Links (12 davon in immutablen Accepted-ADRs) — die C8-Zusagen „`make gates` grün" und „keine Accepted-ADR wird berührt" sind so nicht erreichbar; die nötige Tombstone-/Retarget-/`ids`-Retarget-Kampagne ist unbenannt.

Beide sind durch Ergänzung der Etappen-Beschreibung heilbar (keine Umkehr des Schnitts); der Abnahme-Entscheid sollte sie sowie die Fork-Definitions-Divergenz (F-6) explizit adressieren. LOW/INFO sind nice-to-fix und nicht blockierend.
