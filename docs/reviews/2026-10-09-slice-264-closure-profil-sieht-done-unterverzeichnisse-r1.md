# Review R1 — slice-264: Das Closure-Profil prüft die Slices in den Unterverzeichnissen von `done/`

**Review-Art:** Code (Diff gegen Plan, ADRs und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-264, Commit `c55bae02`
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-264;
[ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md),
[ADR-0059](../plan/adr/0059-closure-waechter-weicht-structure-regel.md),
[ADR-0081](../plan/adr/0081-reviews-modul.md),
[ADR-0082](../plan/adr/0082-uebergangswaechter-reviews-observations.md);
[`MR-049`](../../harness/conventions.md#mr-049),
[`MR-056`](../../harness/conventions.md#mr-056),
[`MR-025`](../../harness/conventions.md#mr-025);
[`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Spezifikation Schritt C2),
[`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(Spezifikation Schritt 2),
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in);
`AGENTS.md` §3.3, §3.5, §3.7, §3.8, §5; Baseline `v6.17.0` ·
`regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst;
vorherige Findings am selben Gegenstand: Reviews R1–R3 zu slice-263, R1–R5 zu slice-265.

## Messungen (Kommando und Ergebnis)

Alle im Wegwerf-Klon des Stands `c55bae02`, Image aus `make build` desselben Stands.

- **Ist-Lauf:** `make verify-closure-notes` → `974 Datei(en) geprüft, 0 Befund(e)`.
- **`**/`-Ausnahmen unter `done/wellenlos/`:** vier Kopien von slice-243 unter
  `done/wellenlos/` — `slice-135-probe.md` und `slice-145-probe.md` mit Ausgang
  `behoben`, `slice-165-probe2.md` und `slice-171-probe2.md` mit einem offenen
  DoD-Haken. Ergebnis: genau zwei Befunde — `section-forbidden` auf slice-145,
  `section-open-tasks-marker-missing` auf slice-171; slice-135 und slice-165
  still. Die Ausnahmen greifen im Unterverzeichnis (`matchGlob`: `**` steht
  auch für null Segmente, `*` je Segment über `path.Match`).
- **Stub-Marker gegen den Bestand:** jede `slice-*.md` unterhalb von `done/*/`
  außerhalb `wellenlos/` trägt `^> **ARCHIVIERT`; alle 254 Zeilen, die das
  Muster im Baum trifft, lauten `> **ARCHIVIERT** — Volltext:` und stehen in
  Stubs. Unter `done/wellenlos/` 24 Volltexte ohne Marker (auch am Plan-Commit
  `fe54da7e`: 24).
- **Zwei Läufe, eine Frage:** die vier Reports zu slice-263 aus `docs/reviews/`
  entfernt. `verify-closure-notes` → `review-missing` auf
  `done/wellenlos/slice-263-…md`, Exit rot. `make review-coverage` über
  denselben Stand → `1106 Datei(en) geprüft, 0 Befund(e)`.

## Findings

### F-1 — MEDIUM — Zwei Läufe beantworten die Review-Deckung jetzt verschieden, und die Spiegel sagen es nicht

- **kategorie:** MEDIUM
- **quelle:** [ADR-0082](../plan/adr/0082-uebergangswaechter-reviews-observations.md)
  Entscheidung 4 (`make review-coverage` bleibt als fokussierter Lauf „zum
  Debuggen"); [`MR-025`](../../harness/conventions.md#mr-025); Skill-Prüffragen 11 und 20
- **pfad:** `harness/sensors/review-coverage.md` · „die Konfiguration dieses Repos
  setzt es noch nicht — die Volltexte unter `done/wellenlos/` bleiben bis dahin
  ungeprüft"; `.d-check.yml` · „GRENZE: beide Verzeichnisse werden NICHT rekursiv gescannt"
- **befund:** Mit `reviews.recursive` im Closure-Profil meldet `verify-closure-notes`
  einen fehlenden Report zu einem Slice unter `done/wellenlos/`, `make review-coverage`
  über denselben Stand bleibt grün (Messung „Zwei Läufe"). Wer den roten
  Übergang mit dem Lauf debuggt, den ADR-0082 dafür vorsieht, sieht Grün; die
  Sensor-Datei sagt weiter, die Konfiguration dieses Repos setze die Schlüssel
  nicht, und das Hauptprofil nennt die Abweichung nicht — benannt ist sie nur im
  Closure-Profil. Da wellenlose Slices inzwischen direkt nach `done/wellenlos/`
  schließen, ist das der Regelfall, nicht ein Rand.
- **verifizierbar:** ja — die Probe oben; `grep -n "setzt es noch nicht" harness/sensors/review-coverage.md`.
- **klasse:** `zwei-profile-zwei-antworten`

### F-2 — MEDIUM — slice-242 R2: der nachgetragene Ausgang widerspricht der §7-Notiz, aus der er stammen soll

- **kategorie:** MEDIUM
- **quelle:** [`MR-049`](../../harness/conventions.md#mr-049); Baseline `v6.17.0` ·
  `regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst;
  `AGENTS.md` §5 Regeln 15 und 16
- **pfad:** `docs/plan/planning/done/wellenlos/slice-242-gate-erweiterung-kein-adr-anlass.md` ·
  „**Ausgang:** entfallen — getragen als benannte Grenze"
- **befund:** §7 derselben Datei führt R2 als „bewusst offener Punkt, Ausgang bei der
  nächsten Aufnahme notieren" — das ist der Kanon-Ausgang *weiter offen* (mit
  Register-Eintrag), nicht *entfallen*. Der Nachtrag setzt ein neues Urteil, das
  §7 nicht trägt, und stützt es auf `AGENTS.md` §3.6 „Kein Gate prüft das", das
  die Prüfbarkeit der ganzen Regel betrifft, nicht den ausstehenden
  Bestands-Beleg für „unabhängig lauffähig"; die Commit-Botschaft („behoben aus
  der jeweiligen Paragraph-7-Notiz") und der Satz „aus §7 nachgetragen" in der
  Datei behaupten die Herkunft trotzdem. Folge: ein offener Punkt verschwindet
  aus dem Zähler, ohne dass ein Register-Eintrag ihn hält.
- **verifizierbar:** nein — Urteil; `sed -n '/^## 7/,/^## 8/p'` auf die Datei zeigt den Widerspruch.
- **klasse:** `nachtrag-widerspricht-quelle`

### F-3 — MEDIUM — slice-240: „entfallen — eingetreten" setzt das falsche der drei Wörter

- **kategorie:** MEDIUM
- **quelle:** [`MR-049`](../../harness/conventions.md#mr-049) (der Wortschatz ist
  geschlossen, damit *welcher der drei* an der Form lesbar ist); Baseline `v6.17.0` ·
  `regelwerk/modul-05-planning-harness.md` §Offene Risiken · „ob das Risiko wirklich
  nicht mehr eintreten kann"
- **pfad:** `docs/plan/planning/done/wellenlos/slice-240-baseline-v6130-bump.md` ·
  „**Ausgang:** entfallen — eingetreten und im Slice aufgefangen"
- **befund:** Drei der vier Risiken sind laut §7 eingetreten und im Slice behoben;
  der Ausgang trägt dennoch *entfallen* und schreibt *eingetreten* in den
  Begründungstext. Die urteilsfreie Hälfte, die das Wort prüft, liest damit „konnte
  nicht eintreten", wo es eintrat — der Bestand führt denselben Fall als
  „eingetreten — … behoben im selben Slice" (slice-231, slice-236). Dazu ist das
  erste Risiko (Über-Hebung) in §7 gar nicht geführt — §7 nummeriert vier andere
  Risiken (sein R4 ist ein Zitat-Delta, das §6 nicht kennt); der Nachtrag leitet
  *entfallen* aus dem Schweigen ab.
- **verifizierbar:** ja — `grep -n "entfallen — eingetreten" docs/plan/planning/done/wellenlos/slice-240-*.md`.
- **klasse:** `ausgang-wort-gegen-inhalt`

### F-4 — LOW — Die Kopplung „dieselbe Bestands-Ausnahme wie im Hauptprofil" trägt im Closure-Profil nicht mehr

- **kategorie:** LOW
- **quelle:** [ADR-0082](../plan/adr/0082-uebergangswaechter-reviews-observations.md)
  Konsequenzen (die Ausnahme lebt an zwei Stellen, beide nachziehen); Benutzerhandbuch
  §Glob-Syntax, Absatz „Eine Ausnahme, gemessen und benannt"
- **pfad:** `.d-check.closure.yml` · „Die fuenf Eintraege sind mit slice-200 entfernt --
  siehe die Begruendung im Hauptprofil"
- **befund:** Die Begründung im Hauptprofil ist die Nicht-Rekursion; im Closure-Profil
  ist der Block jetzt rekursiv, die fünf Slices sind dort nur über ihren Stub-Status
  aus der Menge. Und die `exempt-paths` von `reviews` werden mit blankem `path.Match`
  geprüft: ein Eintrag, der aus dem Hauptprofil übernommen wird, trifft keinen Slice
  unter `done/wellenlos/`, und die `**/`-Form, die die `structure`-Regeln derselben
  Datei tragen, trifft dort gar nichts. Latent — der Block trägt heute keine Ausnahme.
- **verifizierbar:** nein — erst beim nächsten Ausnahme-Eintrag.
- **klasse:** `kopplung-ueber-ungleiche-glob-semantik`

### F-5 — LOW — `SPEC-095` sagt „jeden Volltext", die Grenze derselben Spezifikation nicht

- **kategorie:** LOW
- **quelle:** Spezifikation [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C2, „Zwei Grenzen"; `AGENTS.md` §5 Regel 13
- **pfad:** `spec/spezifikation.md` · „Der ausgelöste Lauf prüft jeden Volltext unter `done/` samt Unterverzeichnissen"
- **befund:** C2 benennt, dass ein Volltext, der den Marker zitiert, und ein Slice hinter
  einem Verzeichnis-Symlink still aus der Prüfung fallen; `SPEC-095` behauptet
  ausnahmslos „jeden". Die Sensor-Datei führt beide Grenzen (Grenze 8), die
  höherrangige Stelle nicht.
- **verifizierbar:** nein — Text gegen Text.
- **klasse:** `zusage-ohne-eigene-grenze`

### F-6 — LOW — Das Stub-Muster ist weiter als die Stub-Form, die es treffen soll

- **kategorie:** LOW
- **quelle:** Spezifikation [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C2 · „das Muster gehört deshalb so eng gefasst, dass es nur die Form des Stubs trifft"
- **pfad:** `.d-check.closure.yml` · `skip-pattern: '(?m)^> \*\*ARCHIVIERT'`
- **befund:** Alle 254 Stub-Marker im Baum lauten `> **ARCHIVIERT** — Volltext:`; das
  Muster endet vor dem schließenden `**` und dem Rest dieser Form und nimmt damit jede
  Zitat-Zeile aus, die mit dem fett gesetzten Wort beginnt. Heute trifft es keinen
  Volltext (Messung oben); ein Slice, der über die Archivierung schreibt, ist der
  naheliegende erste.
- **verifizierbar:** ja — `grep -rhn '^> \*\*ARCHIVIERT' docs/plan/planning/done | sed 's/^[0-9]*://' | cut -c1-30 | sort | uniq -c`.
- **klasse:** `ausnahme-muster-weiter-als-gegenstand`

### F-7 — INFO — Die Plan-Zahl der Volltexte unter `done/wellenlos/` ist um eins zu klein

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** Slice-Plan slice-264 §8 · „Unter `done/wellenlos/` liegen 23 Volltexte"
- **befund:** Am Plan-Commit `fe54da7e` liegen dort 24 Dateien ohne Marker (slice-240–243,
  247–260, 263, 265, 267–270). Ohne Folge für den Diff; die Zahl steht auch in §1.
- **verifizierbar:** ja — Zählung je `git show fe54da7e:<datei>` gegen den Marker.
- **klasse:** `messzahl-ohne-kommando`

### F-8 — INFO — Die Glob-Ausnahme im Benutzerhandbuch gilt für `reviews` nicht mehr ohne Einschränkung

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `docs/user/benutzerhandbuch.md` · „Beide Module listen ihr Verzeichnis **nicht rekursiv**"
- **befund:** Seit `reviews.recursive` existiert, stimmt der Satz nur noch ohne den
  Schlüssel; dieses Repo setzt ihn jetzt im Closure-Profil. Älter als dieser Diff, hier
  erstmals am eigenen Bestand wirksam (siehe F-4).
- **verifizierbar:** nein.
- **klasse:** `handbuch-glob-ausnahme-veraltet`

## Negativbefunde

- **`**/`-Ausnahmen (`MR-049`, `MR-056`):** greifen unter `done/wellenlos/` wie direkt
  unter `done/`, gemessen; die festen Ziffernzahlen bleiben erhalten, kein `*`-Rest-Fresser.
- **Stub-Erkennung gegen den Bestand:** jeder Stub trägt den Marker, kein Volltext trifft
  das Muster — ohne Befund (Weite: F-6).
- **„Dieselbe Kandidatenmenge wie `planning.closure`":** trifft zu — `slice-glob`-Default
  `slice-*.md`, dieselben `SKIP_DIRS`, dasselbe Muster.
- **Grenze „Wellen-Ergebnisnotizen nur flach":** alle `welle-*-results.md` liegen direkt in
  `done/`, keine darunter — ohne Befund.
- **Symlink-Grenze:** Sensor-Datei und Spezifikation C2 / STRUCT Schritt 2 sagen dasselbe — ohne Befund.
- **Weitere Stellen mit „nur `done/` direkt" für den Closure-Lauf:** `AGENTS.md`,
  `harness/README.md`, `Makefile` (Kommentar und `##`-Text von `verify-closure-notes`),
  `.github/workflows/ci.yml`, `docs/user/` — ohne Befund; die zwei verbleibenden Stellen
  betreffen `reviews` und stehen in F-1.
- **Kommentare (`AGENTS.md` §3.7):** die neuen YAML-Kommentare (Zusage + Grenze im
  `closure`-Block, Kopplung über dem `structure`-Block, Abgrenzung im `reviews`-Block) und
  der neue `pre-commit`-Satz (Zusage + Rang-Zeiger) tragen je eine Klasse; keine Slice-Nummer,
  kein Mess-Label neu eingeführt.
- **`pre-commit`/CI-Erkennung:** unverändert, nur der Kommentar — ohne Befund.
- **Spec-Straten (`AGENTS.md` §3.4):** die `SPEC-095`-Zeile und die Historie-Zeile tragen
  keinen Slice-, Wellen- oder ADR-Verweis — ohne Befund.
- **ADRs (`AGENTS.md` §3.5):** keine ADR berührt.
- **Gate-Lockerung (`AGENTS.md` §3.6):** keine Schwelle gesenkt; die Kandidatenmenge wächst.
- **slice-241, slice-243:** die nachgetragenen Ausgänge decken sich mit den §7-Notizen — ohne Befund.
- **Plan-Abgrenzung:** die Archivierung der Volltexte unter `done/wellenlos/` ist nicht
  mitgenommen; die Produkt-Schlüssel nicht berührt. Der `reviews`-Block steht nicht
  wörtlich in Plan §3, ist aber in §8 gemessen — kein Abgrenzungsbruch.
- **Hexagon, Netz, Suppression:** kein Go-Code im Diff — nicht anwendbar.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 3 | 3 | 2 |

Wiederkehrende Klasse: Spiegel einer Scope-Änderung nicht vollständig nachgezogen (F-1,
F-5, F-8) — dieselbe Lage wie R1 F-7 zu slice-260.

## Verdikt

**Nachziehen vor Closure.** F-1 bis F-3 blockieren: F-1 lässt zwei Läufe auf dieselbe
Frage verschieden antworten, F-2 und F-3 setzen Ausgänge, die die Quelle nicht trägt, in
eingefrorenen Dateien, deren Form ein Gate als Wahrheit liest.
