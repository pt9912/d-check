# slice-263: Closure-Prüfung über Unterverzeichnisse und Stub-Ausnahme — Produkt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** Befund F-1 (HIGH) aus dem Review von slice-260;
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Closure-Fähigkeit von `planning`),
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md);
Auftraggeber-Entscheid 2026-10-09 (eigener Slice statt Mitnahme in slice-260).

**Berührte Spec-Stellen:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (Closure-Kandidaten,
`planning.closure.dir`), [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(Dateimenge), `spec/spezifikation.md` §2 (Schlüssel der
`planning`- und `structure`-Konfiguration).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Das Produkt kann geschlossene Slices auch in Unterverzeichnissen
prüfen und archivierte Stubs dabei an ihrem Inhalt erkennen — zwei
opt-in-Schlüssel, ohne die der Befundsatz byte-identisch bleibt:

- `planning.closure.recursive` — die Closure-Fähigkeit liest die Kandidaten
  auch aus den Unterverzeichnissen von `planning.closure.dir`
  (der Basisnamen-Filter bleibt `planning.closure.glob`).
- `planning.closure.skip-pattern` und `structure[].skip-pattern` — ein RE2
  gegen den Datei-Inhalt; ein Kandidat, auf den es passt, ist keiner (ein
  Stub trägt `> **ARCHIVIERT**`). Die Nullmengen-Regel gilt nach dem Abzug:
  bleibt kein Kandidat, ist das fail-closed wie bisher.

Gemessen beim Schnitt: unter `done/wellenlos/` liegen 98 Slice-Dateien, 80
davon Stubs mit Marker und Archiv daneben, 18 Volltexte ohne beides; die
Verzeichnisse `done/welle-*/` tragen nur Stubs.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Konfiguration dieses Repos, die vier Altverstöße und die Rücknahme
  der Grenz-Aussagen** — slice-264; geteilt beim Schnitt, weil die
  Stub-Erkennung ein Kriterium im Produkt braucht (Rückführungs-Bedingung aus
  §4 der ersten Fassung).
- **Andere Module mit einem `done/`-Verzeichnis** (`reviews.done-dir`) — beim
  Beanspruchen messen; trifft dieselbe Blindheit zu, ist das ein eigener
  Befund mit eigenem Slice.
- **Ein `**`-Glob für `structure[].files`** — gibt es bereits.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `planning.closure.recursive` und `planning.closure.skip-pattern`:
      Verfeinerung in der Spezifikation, Konfig-Validierung (Exit 2 bei nicht
      kompilierendem Muster), Tests, die ohne die Änderung aus dem richtigen
      Grund rot sind.
- [ ] `structure[].skip-pattern`: dasselbe für das Modul `structure`.
- [ ] Ohne die neuen Schlüssel ist die Ausgabe unverändert
      (`make blackbox-probe REF=<Stand davor>`); `--print-config` und die
      übrigen Spiegel der Konfiguration nachgezogen; `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Kern-Regeln `planning` (Closure-Kandidaten) und `structure` (Dateimenge) | update | Rekursion, Inhalts-Ausnahme |
| Konfig-Modell und YAML-Adapter | update | zwei Schlüssel samt Validierung |
| `spec/spezifikation.md` (Verfeinerung zu [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) Schritt C2, Verfeinerung zu [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in), §2-Schema) | update | Verfeinerung |
| `--print-config`-Vorlage und weitere Spiegel (beim Beanspruchen per grep gelistet, Schritt 17) | update | Konfig-Oberfläche |
| `spec/lastenheft.md` (beide Anforderungen, Versions-Bump und Historie) | update | die Kandidatenmengen sind dort abschließend beschrieben |

*(Plan-Änderung nach R1, vor dem Code: Das Lastenheft beschreibt die
Kandidatenmengen beider Module abschließend („ausschließlich
`planning.closure.dir`", „abgezogen wird `exempt-paths`"); die neuen Schlüssel
gehören deshalb dorthin, nicht nur in die Spezifikation (R1 F-2,
`AGENTS.md` §5 Regel 3) — Version 0.99.0 mit Historie-Zeile. Dazu: die
Fehlermeldung eines unlesbaren `closure.dir` bleibt ohne die Schlüssel
byte-identisch, das oberste Verzeichnis geht ungereinigt an `List` und in die
Meldung (F-1); das `--print-config`-Beispiel trägt `(?m)` (F-3); die stille
Richtung von `skip-pattern` — ein Volltext, der den Marker zitiert, fällt aus —
und der nicht verfolgte Symlink auf ein Unterverzeichnis stehen als Grenze in
der Spezifikation (F-4, F-7); ein Test hält die `SKIP_DIRS`-Ausnahme (F-5); die
Grund-Code-Zeile von `closure-note-missing` und die §2-Zeile
`structure[].files` ziehen nach (F-8); ein Kommentar ohne Bezug auf den
Vorzustand (F-9). F-6 (Botschaft behauptet zu viel) geht in die
Closure-Notiz.)*

*(Plan-Änderung nach R2, vor dem Code: nur Wortlaut der Spec und der Tests —
`skip-pattern` in die Exit-2-Aufzählung von der `structure`-Anforderung (R2 F-1); die
Nullmengen-Aussage gilt der Gesamtmenge, nicht einem einzelnen
Unterverzeichnis, in der Grund-Code-Zeile von `closure-note-missing` und im fail-closed-Absatz des Lastenhefts
(F-2); Chronik-Wörter in den neuen Tests (F-3); Satzbau des umgeschriebenen
Lastenheft-Absatzes und „wie zuvor" in C2 (F-4); Akzeptanzkriterien für Exit 2
bei `skip-pattern` und das unlesbare Unterverzeichnis, Historie der
Spezifikation vervollständigt (F-5).)*

*(Plan-Änderung nach der Verifikation, vor dem Code: die Exit-2-Tests prüfen,
dass die Meldung den Schlüssel nennt (V-2); eine dritte Review-Runde gibt die
Nachzüge nach R2 frei (V-1).)*

## 4. Trigger

**Start** (`next` → `in-progress`): `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): die Spiegel der Konfig-Oberfläche
  (Vorlage, Vorschlag, Handbuch) reißen mehr als die Kern-Änderung auf — dann
  `structure` in einen eigenen Slice.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Produkt-Verhalten — geht mit dem nächsten Release
hinaus.

## 6. Risiken und offene Punkte

- **Ein Volltext, der den Marker zitiert** — in einem Codeblock oder Zitat —,
  fällt mit `skip-pattern` still aus der Prüfung (R1 F-4). — **Ausgang:** *(offen)*
- **Ein Stub ohne Marker** — ein Stub, den ein älteres Werkzeug ohne
  `ARCHIVIERT` schrieb, würde als Volltext geprüft. Gemessen beim Schnitt:
  alle 80 tragen ihn. — **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** geändert werden Produkt-Kern
(`planning`, `structure`), Konfig-Modell, YAML-Adapter, Konfig-Vorlage und die
Spezifikation — alle unter dem Default `*` (`ALL`); deklariert.

**Spiegel vor dem Editieren** (Schritt 17 des Workflow-Skeletts,
[`MR-025`](../../../../harness/conventions.md#mr-025); gemessen mit
`grep -rln "closure\.glob\|closure:\|EffectiveClosureGlob\|ExemptPaths\|exempt-paths"`
über Code, Spezifikation und Doku): Konfig-Modell, YAML-Adapter samt
Validierung, Konfig-Vorlage (`--print-config`), die Kern-Regeln, Schritt C2
und §2-Schema der Spezifikation, `structure`-Schritt 2 und §2-Schema; das
Benutzerhandbuch zieht die Release-Prep nach (`AGENTS.md` §5 Regel 17). Die
ADR der Closure-Kandidaten ist `Accepted` und wird nicht angefasst.
`--suggest-config` schlägt keine `structure`-Regeln und keinen
`closure`-Block mit Unterverzeichnissen vor — kein Spiegel.

**Gemessen zu `reviews.done-dir`** (§1 Abgrenzung): das Modul `reviews` liest
`done-dir` per `List` ohne Abstieg, die Konfig-Vorlage nennt das
ausdrücklich („nicht rekursiv"). Dieselbe Blindheit für die Volltexte unter
`done/wellenlos/` — ein eigener Befund, nicht in diesem Slice; wird bei der
Closure als Folge-Slice geschnitten.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** gelesen am 2026-10-09.
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/state.md)
(verkörpert, wach) — der Slice öffnet eine neue Ziel-Achse (Unterverzeichnisse,
Inhalts-Ausnahme); jede Zusage wird für sie neu geprüft;
[`BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft`](../observations/BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft/state.md)
— der neue Leseweg (Inhalt vor der Kandidatenwahl) wird gegen seine
Fehlerformen gefahren (unlesbare Datei, ungültiges Muster, alles
ausgenommen); [`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
— jede Liste am Code gezählt (Schritt 18).

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-09 aus dem jüngsten Lauf (`make nightly-state`) —
`upstream-drift` grün (2026-10-08 11:37 UTC), `image-scan` grün
(2026-10-08 10:38 UTC).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
