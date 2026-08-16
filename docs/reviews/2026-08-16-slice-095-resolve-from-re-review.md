# Re-Review-Report: slice-095 — `links.resolve-from` (ortsfeste Verweise) — 2026-08-16

**Review-Art:** Code — **bestätigende Re-Review**: geprüft wird, ob die
fünfzehn Befunde des Erst-Reports
([2026-08-16-slice-095-resolve-from-review.md](2026-08-16-slice-095-resolve-from-review.md))
wirklich geheilt sind — und ob die Heilung neue Defekte eingeführt hat. Jedes
Urteil stützt sich auf eigene Läufe oder Code-/Vertragszitate, nicht auf den
Commit-Text. Sonderfall dieses Slice: zwischen Erst-Review und Re-Review lag
ein **roter CI-Lauf**, der den ersten fail-closed-Zuschnitt widerlegte — F-1
ist deshalb **anders** geheilt als vom Erst-Report nahegelegt (Einzelheiten in
der F-1-Zeile). Nicht geprüft wird die DoD-Abhakung (getrennter Kontext,
Verifikation).

**Gegenstand:** Commit-Range `eb72435..3eaf512` (sieben Commits), im
Besonderen die Heilungs-Commits `aed757b` (Review-Heilung) und `3eaf512`
(CI-Realfund-Justierung + nachgelieferte Lastenheft-Beschreibung);
Arbeitsbaum-Stand `3eaf512` (= HEAD, clean).

**Skill:** `.harness/skills/reviewer.md` @ 1.4.0 ·
**Modell:** claude-fable-5 · **Datum:** 2026-08-16

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`DC-FA-LINK-001`](../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
  in der Lastenheft-Fassung 0.60.1 — **zuerst geprüft:** die in 0.60.0
  fehlende Beschreibung „Ortsfeste Verweise" steht jetzt wirklich im
  Lastenheft (`spec/lastenheft.md:1009`–`1033`, samt beider benannter
  Grenzen); außerdem
  [`DC-FA-LINK-002`](../../spec/lastenheft.md#dc-fa-link-002--symlink-ablehnung),
  [`DC-FA-REF-001`](../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
  (neuer Ventil-Wortlaut),
  [`DC-QA-01`](../../spec/lastenheft.md#dc-qa-01--performance),
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus),
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
- §[`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  Schritt 6 (Fassung nach beiden Heilungs-Commits), das §2-Schema
  `links.resolve-from`, die §4-Zeile `link-position-dependent`,
  §[`DC-FA-REF-001.a`](../../spec/spezifikation.md#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs)
  und der Ventil-Konsument in
  §[`DC-FA-CODE-001.a`](../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)
- [ADR-0056](../plan/adr/0056-resolve-from-wandernde-quellorte.md) (Proposed,
  jetzt mit Entscheidung 6: Ziel-Wanderung als Grenze),
  [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md),
  [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md),
  [`MR-013`](../../harness/conventions/MR-013-lifecycle-move-buendelung.md),
  [`AGENTS.md`](../../AGENTS.md) §3, [`CLAUDE.md`](../../CLAUDE.md)
- Der Slice-Plan
  [slice-095](../plan/planning/in-progress/slice-095-links-resolve-from.md)
  (§3a inkl. der korrigierten Überclaim-Passage) und das Wellendokument
  [welle-76](../plan/planning/welle-76-ortsfeste-verweise.md)
  (Closure-Trigger)

**Läufe dieses Re-Reviews.** Baseline `make test` grün (Exit 0); das Image aus
dem HEAD-Stand frisch gebaut. **Sechs Mutationsläufe**, jeder über eine
Dateikopie (Sicherung in einem Scratch-Verzeichnis außerhalb des Repos,
mutieren → `make test` → Original byte-identisch zurückschreiben), Ergebnis am
**Exit-Code** abgelesen, nie per `git checkout`: M-1 (`dirs < 2`-Config-Check
aus), M-2 (Divergenz-Sortierung aus), M-3 (`path.Clean` aus dem
Gruppen-Vergleich), M-4 (fail-closed-Zweig „kein Ort existiert" aus), M-5
(Symlink-Vorbedingung aus), M-6 (die CI-Abschwächung zurückgedreht auf „jeder
fehlende Ort meldet"). Rund **fünfzehn Fixture-Läufe** gegen das HEAD-Image
(netzlos, read-only) in dem Scratch-Verzeichnis, darunter die
**Klon-Replikation** (`git archive HEAD`, entpackt, eigenes `.d-check.yml`),
die Voll-Tippfehler-, Datei-Ort-, Einzel-Tippfehler- und Probe-Datei-Fälle,
zwei Scope-/Scan-Stille-Proben, fünf Config-Rand-Proben, `--print-config`,
`--doctor` über zwei Befund-Fixtures, ein Divergenz-Fixture mit
Doppel-Lauf-Byte-Vergleich und die **Retro-Replikation** (`git archive` des
Stands vor `2a94a408`). Alle Probe-Configs lagen ausschließlich im
Scratch-Verzeichnis und sind gelöscht; keine Repo-Datei außer diesem Report
wurde dauerhaft geändert. Abschließend `make gates` und
`make verify-closure-notes` über den Baum samt diesem Report.
`git status --short` ist am Ende leer bis auf diesen Report.

---

## Urteils-Tabelle F-1 … F-15

| Befund | Urteil | Beleg (eigener Lauf / Zitat) |
|---|---|---|
| F-1 (HIGH) Tippfehler-Ort schaltet Quellen-Rolle still ab | **teilweise geheilt — anders als vom Erst-Report nahegelegt, mit benannter Grenze** | Neuer Post-Pass `CheckResolveFromDirs` (`internal/hexagon/core/rules/links_resolvefrom.go:19`). Kehrseite belegt: Voll-Tippfehler-Gruppe ⇒ genau **1 Befund, Exit 1**; Ort als **Datei** ⇒ 1 Befund, Exit 1. Klon-Replikation (`git archive HEAD` + HEAD-Image): **384 Dateien, 0 Befunde, Exit 0** — das legitim geleerte `open/` fehlt im Klon und meldet nicht. Der Erst-Review-Fall (**Einzel**-Tippfehler `in-progres`) ist jetzt **per Design still** (eigener Lauf: 0 Befunde; Probe-Datei mit präfixlosem Nachbar-Verweis meldet mit korrekter Gruppe 1, mit Tippfehler 0) — die Grenze steht wörtlich auf Lastenheft (`spec/lastenheft.md:1027`–`1030`), Schritt 6 (`spec/spezifikation.md:911`–`915`), 0.60.1-Historie, Code-Kommentar und in zwei Tests; einen dritten Zuschnitt, der Einzel-Tippfehler von legitim geleerten Orten trennt, habe ich nicht gefunden (Negativbefund unten). Guards: M-4 und M-6 sterben an je genau ihrem Test. **Rest:** die dritte Stille-Variante des Erst-Befunds (Ort außerhalb `scan.roots`/`scan.ignore` bzw. `links.scope`) ist weder geheilt noch benannt → N-1 |
| F-2 (MEDIUM) Exit-2-Zusagen ungetestet | **geheilt** | `TestDecode_ResolveFromFehler` (sieben Rand-Fälle, Haus-Muster der Negativtabellen) + `TestDecode_ResolveFromHappy` in `internal/adapter/driven/configyaml/configyaml_test.go:297`; Mutation M-1 (`dirs < 2`-Check entfernt) stirbt an genau diesem Test (`make test` Exit 2, benannter FAIL) |
| F-3 (MEDIUM) Ist-Ort-Vorbedingung nur im Code | **geheilt** | Schritt 6 trägt sie normativ: „**Vorbedingung ist der saubere Ist-Ort** … dieser Schritt schweigt, sonst trüge dieselbe Referenz zwei Befunde" (`spec/spezifikation.md:896`–`899`, jetzt inkl. `repo-escape`/`symlink`); die Lastenheft-Beschreibung sagt dasselbe (`spec/lastenheft.md:1022`–`1024`). Produktverhalten unverändert (Test bestand schon) |
| F-4 (MEDIUM) [`DC-FA-REF-001`](../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)-Wortlaut widerspricht aktiv | **geheilt** (auf der monierten Fläche) | Die abschließende Aufzählung nennt jetzt die „Auflösungs-Klasse … (`target-missing`/`anchor-missing`/`codepath-missing` und `link-position-dependent`)" (`spec/lastenheft.md:950`–`954`); der codepaths-Spiegel in §[`DC-FA-CODE-001.a`](../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) ist nachgezogen (`spec/spezifikation.md:1148`–`1150`). **Rest:** die Wirkungs-Aufzählung in §[`DC-FA-REF-001.a`](../../spec/spezifikation.md#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs) selbst blieb bei vier Codes → N-3 |
| F-5 (MEDIUM) `--print-config` ohne `resolve-from` | **geheilt** | Die emittierte Vorlage des Images führt den auskommentierten `links.resolve-from`-Block mit `dirs`/`fixed-dirs`, Quellen-Semantik und `>= 2`-Hinweis (live verifiziert); das Kopf-Kommentar nennt `links` nicht mehr optionslos |
| F-6 (MEDIUM) „zeichengenau"-Überclaim | **teilweise geheilt** | Vier Flächen korrigiert und von mir nachgemessen: Slice §3a („die Zahlengleichheit ist Zufall, kein Beweis", 15/19-Überlappung, Ziel-Wanderung benannt), ADR-0056 §Kontext, welle-76-Trigger, Lastenheft-0.60.1-Zeile. Eigene Retro-Replikation: exakt **19** `link-position-dependent`, Verteilung 7/4/3/3/2 (dieser Slice selbst 2); Real-Bruch-Zusammensetzung am Commit `2a94a408` gegengeprüft (15 in wandernden Quellen: 7/3/2/2/1, plus 4 im Review-Report [2026-08-09-backlog-schnitt-review.md](2026-08-09-backlog-schnitt-review.md)). **Nicht korrigiert:** CHANGELOG `[Unreleased]` und der `.d-check.yml`-Kommentar tragen den Überclaim weiter → N-2 |
| F-7 (MEDIUM) Closure-Trigger nicht einlösbar | **teilweise geheilt** | Der welle-76-Trigger verspricht nicht mehr „die drei belegten Move-Brüche wären rot gewesen", sondern nur die **Slice-Hälfte** des 19er-Bruchs (retro gemessen — von mir repliziert); die Ziel-Wanderungs-Hälfte ist als Grenze deklariert und in [ADR-0056](../plan/adr/0056-resolve-from-wandernde-quellorte.md) Entscheidung 6 verankert; die Slice-DoD-Zeile ist analog präzisiert. **Rest:** die Zuordnung des Wellendokument-Falls zur „Ziel-Wanderungs-Hälfte" ist unpräzise und von Entscheidung 6 nicht gedeckt → N-4 |
| F-8 (LOW) Doppelbefund bei Symlink-Ziel | **geheilt** | Die Vorbedingung prüft jetzt `symlinkInPath` (`internal/hexagon/core/rules/links_resolvefrom.go:111`–`113`); `TestResolveFromKeinDoppelbefundBeiSymlink` pinnt genau ein `symlink`, kein `link-position-dependent`; Mutation M-5 stirbt an genau diesem Test |
| F-9 (LOW) Unterverzeichnis-Grenze unbenannt | **teilweise geheilt** | Beide Vertragsflächen sagen die Grenze jetzt: „deren **unmittelbare** Dateien Quellen sind — Unterverzeichnisse wandern nicht als Einheit mit" (`spec/lastenheft.md:1014`–`1015`; `spec/spezifikation.md:904`–`906`). **Rest:** kein Akzeptanzkriterium und kein Test pinnt sie (der exakte Vergleich in `resolveFromGroupOf` ist unverändert und weiter nur indirekt bewacht) — als LOW tragbar |
| F-10 (LOW) `..`-Substring statt Segment; Duplikat-Meldung | **nicht geheilt** (offen) | Eigene Läufe: `dirs: [a..b, c]` ⇒ Exit 2 (die Zusage nennt weiter nur das `..`-**Segment**, `spec/spezifikation.md:919`–`920`); `dirs: [a, a]` ⇒ Exit 2 mit der Meldung „dirs-Mitglied **mehrerer Gruppen**" (`internal/adapter/driven/configyaml/configyaml.go:533`, `:540`). Der Erst-Verdikt hatte LOW als Release-Prep-/Closure-abtragbar eingeordnet; die Heilungs-Commits beanspruchen F-10 nicht — als LOW weiter tragbar, aber offen |
| F-11 (LOW) `Kind` vor `escaped` | **geheilt** | Der Hypothesen-Zweig prüft `escaped` zuerst und fragt einen Pfad außerhalb der Wurzel nicht ab (`internal/hexagon/core/rules/links_resolvefrom.go:118`–`124`, mit Kommentar); der Ist-Ort-Zweig ist nur mit nicht-escapenden Zielen erreichbar (`CheckResolveFrom` filtert vorher über dieselbe Auflösung). Per Lektüre verifiziert — wie im Erst-Report festgestellt, sieht kein Gate den Zugriff |
| F-12 (LOW) Sortierung/`path.Clean` unbewacht | **geheilt** | `TestResolveFromDivergenzMeldungSortiert` (Meldungs-Assertion über alle vier Ziele) und `TestResolveFromDirMitSchlussSlash`; Mutationen M-2 und M-3 sterben an je genau ihrem Test (`make test` Exit 2, benannte FAILs) |
| F-13 (LOW) „§Ortsfeste Verweise" ohne Gegenstück; Beschreibung fehlt | **geheilt** | Die Beschreibung ist nachgeliefert (`spec/lastenheft.md:1009`–`1033`, als eigener „Ortsfeste Verweise"-Absatz der Anforderung, inkl. Vorbedingung und **beider** Grenzen); die 0.60.1-Historien-Zeile benennt den Batch-Editor-Verlust ehrlich; die „§Ortsfeste Verweise"-Referenzen in CHANGELOG, `.d-check.yml` und Code haben jetzt ein reales Gegenstück |
| F-14 (INFO) `roadmap.md` ortsfest im wandernden Ort | **nicht geheilt** (offen) | Keine Fläche dokumentiert die Annahme (geprüft: Lastenheft-Beschreibung, Schritt 6, ADR-0056, `.d-check.yml`-Kommentar). Heute weiterhin grün (Klon-Lauf 0 Befunde); als INFO tragbar |
| F-15 (INFO) Ablage-Kosmetik | **nicht geheilt** (offen) | `ReasonLinkPositionDependent` steht weiter unter dem structure-Block-Kommentar (`internal/hexagon/core/model/finding.go:47`–`51`); die §4-Zeile weiter zwischen `closure-note-*` und `wave-*` (`spec/spezifikation.md:2597`). Als INFO tragbar |

## Neue Findings

### N-1 — F-1-Rest: ein Gruppen-Ort außerhalb der Scan-/Scope-Abdeckung verliert die Quellen-Rolle weiterhin still — der neue fail-closed-Pass prüft Existenz, nicht Abdeckung, und keine Fläche benennt die Grenze

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-LINK-001`](../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
  · Reviewer-Anker „Modul-Grenze auf der Ziel-Achse" (die Gruppen-Orte sind
  selbst benannte Verzeichnisse, die der Post-Pass direkt liest, während die
  Quellen-Zusage nur für **gescannte** Dateien gilt)
- `pfad`: `internal/hexagon/core/rules/links_resolvefrom.go:24` (`fsys.Kind`
  direkt, unabhängig vom Scan); `internal/hexagon/core/rules/run.go:150`–`152`
  (Post-Pass nur an `active["links"]` gebunden);
  `spec/spezifikation.md:893` („für jede **gescannte** Datei")
- `befund`: Zwei Scratch-Läufe gegen das HEAD-Image, jeweils mit existierenden
  Gruppen-Orten und einer Probe-Datei, die mit voller Abdeckung **1 Befund**
  erzeugt: `scan.roots: ["spec"]` ⇒ 3 Dateien, **0 Befunde, Exit 0**;
  `links.scope.roots: ["spec"]` ⇒ 385 Dateien, **0 Befunde, Exit 0**. Die
  Orte existieren, der fail-closed-Pass schweigt korrekt — aber keine Datei
  der Gruppe wird je als Quelle geprüft, und genau die Klasse, die der Slice
  maschinell machen soll, bleibt dort dauerhaft unsichtbar. Das ist die im
  Erst-F-1 mitbenannte dritte Stille-Variante; die Heilung deckt Tippfehler
  (Voll-Ausfall) und Datei-Orte, die Abdeckungs-Stille steht auf keiner der
  beiden Vertragsflächen und in keinem Test. Die eigene Konfiguration ist
  nicht betroffen (`roots: ["."]` deckt die Gruppe — deshalb keine
  Gate-Pfad-Eskalation).
- `verifizierbar`: ja — die zwei Scratch-Läufe (in diesem Re-Review
  ausgeführt).
- `klasse`: config-zeiger-fail-open (Rest auf der Abdeckungs-Achse)

### N-2 — F-6-Rest: CHANGELOG `[Unreleased]` und der `.d-check.yml`-Kommentar behaupten weiter „exakt die 19 des realen Bruchs" — im CHANGELOG sogar „mit identischer Verteilung"

- `kategorie`: MEDIUM
- `quelle`: [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (dieselbe Aussage auf allen Flächen) · Messaussagen-Ehrlichkeit (Präzedenz
  wie Erst-F-6)
- `pfad`: `CHANGELOG.md:26`–`27` („meldet exakt die **19** des realen Bruchs,
  mit identischer Verteilung"); `.d-check.yml:29` („meldet das Produkt exakt
  die 19 des realen Bruchs")
- `befund`: Die Heilung hat Slice §3a, ADR-Kontext, welle-76-Trigger und
  Lastenheft-Historie auf die ehrliche 15/19-Überlappung korrigiert — die
  beiden verbliebenen Flächen sagen weiter die widerlegte Identität (die
  Retro-19 und die 19 des realen Bruchs überlappen zu 15; die reale
  Verteilung war 7/8/4 über Slice/Geschwister/Review-Report, die Retro-
  Verteilung ist 7/4/3/3/2 ohne Report-Anteil — beides in diesem Re-Review
  nachgemessen). Der CHANGELOG-Eintrag ist die künftige Release-Fläche und
  widerspricht damit vier korrigierten Flächen; die Heilungs-Commits haben
  beide Dateien nicht angefasst (Range-Diff leer).
- `verifizierbar`: ja — Textabgleich plus Retro-Replikation und
  `2a94a408`-Diff (beide in diesem Re-Review ausgeführt).
- `klasse`: messaussage-praezision (Rest auf zwei Flächen)

### N-3 — F-4-Rest: die Wirkungs-Aufzählung in §[`DC-FA-REF-001.a`](../../spec/spezifikation.md#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs) zählt vier unterdrückte Codes — das Produkt unterdrückt fünf

- `kategorie`: MEDIUM
- `quelle`: §[`DC-FA-REF-001.a`](../../spec/spezifikation.md#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs)
  (die kanonische Ventil-Fläche der technisch verbindlichen Quelle) ·
  [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
- `pfad`: `spec/spezifikation.md:823`–`825` („Wirkung. Ignoriert → das Modul
  überspringt … die **Existenz**-, die **Repo-Escape**- und … die
  **Anker**-Prüfung (kein
  `target-missing`/`codepath-missing`/`repo-escape`/`anchor-missing`)") gegen
  `internal/hexagon/core/rules/links_resolvefrom.go:70`–`72` und
  `spec/lastenheft.md:950`–`954`
- `befund`: Die Heilung hat das Lastenheft
  ([`DC-FA-REF-001`](../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus))
  und den codepaths-Konsumenten nachgezogen, die zentrale
  Wirkungs-Definition des Ventil-Algorithmus aber nicht: wer dort nachliest,
  was das Ventil unterdrückt, kommt ohne den fünften Code
  `link-position-dependent` heraus — dieselbe Lücken-Form, für die Erst-F-3/F-4
  MEDIUM vergaben, jetzt zwischen ranghöchster und technisch verbindlicher
  Fläche. Der Satz in Schritt 6 („das Ventil gilt wie in Schritt 4") sagt die
  Unterdrückung nur auf der Konsumenten-Seite.
- `verifizierbar`: ja — Textabgleich; das Produktverhalten selbst ist
  getestet (`TestResolveFromVentil`).
- `klasse`: vertragsflaeche-luecke

### N-4 — F-7-Rest: der welle-76-Trigger ordnet den Wellendokument-Fall der „Ziel-Wanderungs-Hälfte" zu — er ist eine Quell-Wanderung in Eltern/Kind-Geometrie, die Entscheidung 6 nicht deckt und keine Gruppe ausdrücken kann

- `kategorie`: LOW
- `quelle`: Trigger-Ehrlichkeit (wie Erst-F-7) ·
  [ADR-0056](../plan/adr/0056-resolve-from-wandernde-quellorte.md)
  Entscheidung 6 (definiert die Grenze als „Quelle **ortsfest**, Ziel
  wandert")
- `pfad`: `docs/plan/planning/welle-76-ortsfeste-verweise.md:41`
  („die Ziel-Wanderungs-Hälfte (Review-Reports, **Wellendokument**) ist als
  Grenze in der ADR benannt");
  `docs/plan/planning/in-progress/slice-095-links-resolve-from.md:122`–`125`
- `befund`: Der Wellendokument-Fall (Fall 2 des Slice §2: das flache
  Wellendokument verwies auf `done/…` und brach beim **eigenen** Move) hat
  eine wandernde Quelle, kein wanderndes Ziel — Entscheidung 6 deckt ihn
  wörtlich nicht, und als `resolve-from`-Gruppe ist er nicht ausdrückbar
  (der flache Ort ist Eltern-, nicht Geschwister-Verzeichnis des Ruheorts;
  nur der Re-Evaluierungs-Trigger „nicht-Geschwister-Orte" streift die
  Geometrie). Der Trigger bleibt damit für diesen einen Fall unpräzise: er
  verweist auf eine ADR-Grenze, die den Fall nicht benennt. Kein
  Produktfehler; die nächste Closure dieser Welle führt die Klasse am
  eigenen Wellendokument erneut vor ([`MR-013`](../../harness/conventions/MR-013-lifecycle-move-buendelung.md)
  fängt sie prozessual).
- `verifizierbar`: ja — Textabgleich Trigger ↔ Entscheidung 6 ↔ Slice §2
  Fall 2.
- `klasse`: trigger-praezision

### N-5 — Die §4-Zeile und der `--doctor`-Klartext decken den dritten Produzenten der Kennung nicht: für den fail-closed-Gruppen-Ort-Befund ist die Diagnose faktisch falsch

- `kategorie`: LOW
- `quelle`: §[`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  Schritt 6 („Beides meldet … über denselben Grund-Code") ·
  [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
- `pfad`: `spec/spezifikation.md:2597` (§4-Zeile: nur „relativer Verweis
  einer Datei …"); `internal/hexagon/core/app/diagnose.go:132`
- `befund`: Am Voll-Tippfehler-Fixture zeigt `--doctor` für den
  Gruppen-Ort-Befund den Klartext „relativer Verweis löst nicht von jedem Ort
  seiner resolve-from-Gruppe auf dasselbe Ziel auf" — die Stelle ist aber ein
  Config-Ort (`docs/plan/planning/opne`, Zeile 1), kein Verweis; die §4-Zeile
  beschreibt dieselbe Kennung ebenfalls nur über die Referenz-Hälfte. Wer den
  fail-closed-Befund nachschlägt, bekommt die falsche Reparatur-Richtung
  (Pfad präfixieren statt Config korrigieren). Die Detail-Meldung des Befunds
  selbst ist korrekt; Schritt 6 dokumentiert den Produzenten — nur die beiden
  Kennungs-Nachschlage-Flächen hinken (dieselbe Klasse wie N-4 des
  slice-102-Re-Reviews, dort am `wave-drift`-Sammelcode).
- `verifizierbar`: ja — `--doctor`-Lauf über das Voll-Tippfehler-Fixture (in
  diesem Re-Review ausgeführt).
- `klasse`: vertragsflaeche-luecke (Kennungs-Nachschlag)

### N-6 — Die Spezifikations-Historie §7 trägt keine Zeile für den Heilungs-Nachzug — der Schritt-6-Vertrag hat sich zweimal geändert, die Historie sagt nichts

- `kategorie`: LOW
- `quelle`: Haus-Muster der §7-Historie (zwei „Nachzug nach
  Review/Re-Review"-Zeilen vom selben Tag stehen als Präzedenz direkt
  darunter) · [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
- `pfad`: `spec/spezifikation.md:2637` (die einzige Schritt-6-Zeile —
  unverändert die Einführungs-Zeile); der Range-Diff `1aa39cc..3eaf512`
  enthält keine Historien-Zeile
- `befund`: Beide Heilungs-Commits haben Schritt 6 substanziell fortgeschrieben
  (Ist-Ort-Vorbedingung, Unterverzeichnis-Grenze, fail-closed-Rand samt
  CI-Justierung und benannter Grenze) — die append-only-Historie der
  technisch verbindlichen Quelle führt davon nichts; das Lastenheft hat seine
  ehrliche 0.60.1-Zeile, die Spezifikation nicht. Wer die Historie als
  Drift-Spur liest, hält Schritt 6 für unverändert seit Einführung.
- `verifizierbar`: ja — Range-Diff (in diesem Re-Review geprüft).
- `klasse`: historie-luecke

### N-7 — Lastenheft-Wortlaut „oder ein Ort als Datei existiert" ist weiter als der Code: geprüft werden nur `dirs`-Orte, `fixed-dirs` nicht

- `kategorie`: INFO
- `quelle`: dokumentationswürdige, undokumentierte Annahme (Reviewer-Anker
  INFO)
- `pfad`: `spec/lastenheft.md:1029`–`1030` („wenn **kein** `dirs`-Ort der
  Gruppe existiert oder **ein Ort** als Datei existiert") gegen
  `internal/hexagon/core/rules/links_resolvefrom.go:23`–`31` (Schleife nur
  über `g.Dirs`)
- `befund`: „Ein Ort" liest sich als jeder Ort der Gruppe; der Datei-Check
  läuft nur über `dirs`. Ein `fixed-dirs`-Ort als Datei oder mit Tippfehler
  wird nur indirekt laut — über Referenz-Befunde jeder Quelle („löst von
  … nicht auf"), und bei einer Gruppe ohne relative Quell-Verweise gar nicht.
  In der Praxis fail-loud, aber mit Referenz- statt Config-Diagnose; die
  Spezifikations-Fassung (`spec/spezifikation.md:907`–`910`) ist enger
  formuliert und deckt den Code. Kein beobachteter Fehlbefund.
- `verifizierbar`: ja — Wortlaut-Abgleich; konstruierbar über ein Fixture mit
  `fixed-dirs`-Tippfehler.
- `klasse`: zusagen-praezision-rand

## Negativbefunde (geprüft, ohne Befund)

- **Klon-Replikation (der CI-Realfall):** `git archive HEAD` in ein
  Scratch-Verzeichnis entpackt, HEAD-Image darüber (netzlos, read-only):
  **384 Dateien, 0 Befunde, Exit 0** — `open/` fehlt im Klon (kein tracked
  Inhalt) und meldet nicht, `next/` überlebt per `.gitkeep`. Der rote
  CI-Lauf ist mit dem neuen Zuschnitt nicht reproduzierbar. Ohne Befund.
- **Fail-closed-Kehrseite und Determinismus des neuen Passes:**
  Voll-Tippfehler ⇒ genau 1 Befund; Ort als Datei ⇒ genau 1 Befund; beide
  Meldungen entstehen in Config-Reihenfolge (Gruppen, dann Orte — keine
  Map-Iteration). Rand: fällt ein Datei-Ort mit dem „kein Ort
  existiert"-Befund auf dasselbe Tupel, überlebt deterministisch die
  Datei-Diagnose (Dedup über das 5-Tupel, stabile Sortierung) — nie still,
  nie flackernd. Ohne Befund.
- **Dritter Zuschnitt gesucht, keiner gefunden:** ein
  Ähnlichkeits-Vergleich („ähnlich benanntes Verzeichnis existiert") wäre
  Heuristik — die Klasse, die das Haus in
  [`DC-FA-REF-001`](../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
  ausdrücklich verworfen hat; „existierte in früherer Revision" wäre eine
  git-Achse in einem hermetischen Modul
  ([`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Schnitt
  der git-Module ist strikt opt-in). Die benannte Grenze ist der ehrliche
  Stand; sie ist auf allen vier Flächen (Schritt 6, Lastenheft, 0.60.1-Zeile,
  Code-Kommentar) konsistent formuliert und beidseitig getestet
  (`TestResolveFromEinzelnerLeererOrtIstGruen` /
  `TestResolveFromGruppeOhneExistentenOrtFailClosed`). Negativbefund im Sinn
  des Auftrags.
- **Post-Pass-Verdrahtung:** `CheckResolveFromDirs` läuft genau einmal je
  Lauf (ein `RunWithVCS` je CLI-Invocation, `--suggest` läuft ohne
  `resolve-from`-Config); er läuft auch bei skopiertem `links`
  (Scope-Fixture mit Voll-Tippfehler: 1 Befund trotz `scope.roots: ["spec"]`
  — fail-closed bleibt scharf); kein Doppel über `--doctor`/`--json`. Ohne
  Befund (die Scope-Stille der **Datei**-Hälfte ist N-1).
- **Symlink-Rand des neuen Passes:** ein Gruppen-Ort, der Symlink ist, zählt
  weder als Verzeichnis noch als Datei — eine Voll-Symlink-Gruppe meldet
  fail-closed („kein dirs-Verzeichnis existiert"); konsistent zur
  Haus-Symlink-Ablehnung
  ([`DC-FA-LINK-002`](../../spec/lastenheft.md#dc-fa-link-002--symlink-ablehnung)).
  Ohne Befund.
- **Escape-Disziplin:** der Hypothesen-Zweig wertet `escaped` vor jedem
  `Kind`-Aufruf aus; der Ist-Ort-Zweig ist nur mit nicht-escapenden Zielen
  erreichbar. Kein Metadaten-Zugriff außerhalb der Wurzel mehr (F-11). Ohne
  Befund.
- **Determinismus ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)):**
  Divergenz-Fixture, zwei Image-Läufe: stdout und stderr je byte-identisch;
  die Divergenz-Meldung nennt die Ziele sortiert (JSON-Beleg); die
  Auflösbarkeits-Meldung nennt den ersten Ort in Config-Reihenfolge. Ohne
  Befund.
- **Mutations-Echtheit (sechs Stichproben):** M-1 →
  `TestDecode_ResolveFromFehler`, M-2 → `TestResolveFromDivergenzMeldungSortiert`,
  M-3 → `TestResolveFromDirMitSchlussSlash`, M-4 →
  `TestResolveFromGruppeOhneExistentenOrtFailClosed`, M-5 →
  `TestResolveFromKeinDoppelbefundBeiSymlink`, M-6 →
  `TestResolveFromEinzelnerLeererOrtIstGruen` — jede Mutation `make test`
  Exit 2 mit genau dem benannten FAIL, kein Compile-Artefakt; die zwei im
  Erst-Report offenen Ränder (F-12) sind jetzt tödlich, und die
  CI-Abschwächung ist in **beide** Richtungen gepinnt. Ohne Befund.
- **Retro-Beleg:** eigene Replikation (Archiv des Stands vor `2a94a408`,
  HEAD-Image, Gruppen-Config): exakt **19** `link-position-dependent`,
  Verteilung 7/4/3/3/2, dieser Slice selbst 2 — deckungsgleich mit den
  korrigierten §3a-/ADR-Aussagen; die Real-Bruch-Zusammensetzung (15
  wandernde Quellen + 4 Review-Report-Links) am Commit-Diff verifiziert.
  Ohne Befund auf den vier korrigierten Flächen (Rest: N-2).
- **Lastenheft-Beschreibung gegen den Algorithmus:** die nachgelieferte
  Beschreibung widerspricht Schritt 6 und dem Code in keinem geprüften Punkt
  (Quellen-/`fixed-dirs`-Semantik, Unterverzeichnis-Grenze, Vorbedingung,
  beide Fehlarten, beide Grenzen, Byte-Identität ohne Block); einzige
  Wortlaut-Weite: N-7 (INFO). Ohne weiteren Befund.
- **`--print-config`:** die Vorlage führt den Block schema-konform (zwei
  `dirs`, ein `fixed-dirs`) und dekodiert über den eigenen Parser weiterhin
  fehlerfrei (Template-Test unverändert grün in `make test`). Ohne Befund.
- **Hard Rules und Range-Hygiene:** keine Inline-Suppression, keine
  Gate-Lockerung, keine Änderung an einer `Accepted`-ADR im Range
  (ADR-0056 bleibt Proposed und damit fortschreibbar); die Heilung liegt in
  zwei Commits mit ehrlichen Commit-Texten (der zweite deklariert den
  CI-Realfund und den Batch-Editor-Verlust selbst). Ohne Befund.
- **Nicht nachgezogene Nutzer-Doku:** Handbuch, READMEs und `operations.md`
  tragen `resolve-from` weiterhin nicht — wie im Erst-Report erklärtermaßen
  Release-Prep-Arbeit, kein Befund dieses Reports (die
  Release-Prep-Checkliste greift dort).
- **Gates:** `make test` grün (Baseline und nach allen Rückschreibungen),
  `make gates` grün über den Baum **samt diesem Report** (385 Dateien, 0
  Befunde im doc-check), `make verify-closure-notes` grün (355 Dateien, 0
  Befunde). `git status --short` am Ende: nur dieser Report.

## Kategorie-Summary (neue Findings)

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 3 | N-1, N-2, N-3 |
| LOW | 3 | N-4, N-5, N-6 |
| INFO | 1 | N-7 |

Erst-Befunde: 8 geheilt (F-2, F-3, F-4, F-5, F-8, F-11, F-12, F-13),
4 teilweise geheilt (F-1, F-6, F-7, F-9), 3 nicht geheilt/offen
(F-10, F-14, F-15 — alle LOW/INFO, vom Erst-Verdikt als Release-Prep-/
Closure-abtragbar eingeordnet).

## Verdikt

**APPROVE mit Auflagen (Text-Nachzug vor dem Release).** Die Blocker des
Erst-Reports sind substanziell geheilt und halten eigener Nachmessung stand:
der fail-closed-Pass fängt Voll-Tippfehler und Datei-Orte, der Klon-Fall ist
grün, die CI-Abschwächung ist ehrlich benannt, konsistent formuliert und in
beide Richtungen getestet — einen besseren Zuschnitt als die benannte Grenze
habe ich nicht gefunden. Ist-Ort-Vorbedingung (inkl. Symlink), Ventil-Wortlaut
der ranghöchsten Fläche, Negativtabelle, Template, Sortierungs- und
Normalisierungs-Guards: alles verifiziert, sechs Mutanten sterben an genau
ihren Tests. Was bleibt, ist ausschließlich Text-/Vertragsarbeit: die
Abdeckungs-Stille als Grenze oder Rand binden (N-1), zwei Flächen vom
15/19-Überclaim befreien (N-2), die zentrale Ventil-Wirkungs-Aufzählung um den
fünften Code ergänzen (N-3) — kein Code-Zweig muss sich dafür ändern.

**Release-Empfehlung:** Minor-Release **nach** dem Text-Nachzug von N-1, N-2
und N-3 im noch offenen Slice (danach sind sie immutabel bzw. Release-Fläche);
N-4–N-6 und die offenen Erst-Reste F-10/F-14/F-15 sind mit Release-Prep oder
Closure-Notiz abtragbar, N-7 ist eine Formulierungs-Präzisierung bei
Gelegenheit. Die Nutzer-Doku (Handbuch, READMEs, `operations.md`,
CHANGELOG-Feinschliff) gehört ohnehin in den Release-Prep-Commit.
