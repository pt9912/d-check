# Review-Report: Backlog-Schnitt slice-094…098 — 2026-08-09

**Review-Art:** **Plan-/Schnitt-Review** — geprüft wird nicht der Inhalt der
einzelnen Slice-Pläne, sondern der **Schnitt zwischen** ihnen: Slice-Grenzen,
Reihenfolge, Wellen-Zuordnung, Ablöse-Pfad, CR-Abdeckung und die
Entscheidungen tragenden Messwerte. Gegen Baseline-Regelwerk
`modul-05-planning-harness.md` / `modul-06-roadmap.md`, gegen
[`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
und die referenzierten ADRs.

**Gegenstand:** die fünf Slices in `docs/plan/planning/open/` —
[slice-094](../plan/planning/open/slice-094-closure-zaehl-paritaet.md),
[slice-095](../plan/planning/done/slice-095-links-resolve-from.md),
[slice-096](../plan/planning/in-progress/slice-096-structure-modul-analyse.md),
[slice-097](../plan/planning/open/slice-097-closure-glob-entkopplung.md),
[slice-098](../plan/planning/open/slice-098-closure-note-placeholder.md);
Arbeitsbaum bei HEAD `e03afea`, Working-Tree sauber.

**Skill:** `.harness/skills/reviewer.md` @ 1.3.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-08-09

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- die fünf Slice-Pläne oben; zuletzt geschlossen
  [slice-093](../plan/planning/done/slice-093-closure-note-gate.md)
- [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
  (die Fähigkeit, um die vier der fünf Slices kreisen),
  [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md)
  (Präzedenz für Ablöse-/Alias-Pfad),
  [ADR-0042](../plan/adr/0042-markdown-lexik-folgt-commonmark.md) und
  [ADR-0043](../plan/adr/0043-tabellengrenze-am-relevanten-header.md)
  (SemVer-Präzedenz „findet mehr")
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (Lastenheft) und
  [§DC-FA-PLAN-001.a](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  (Spezifikation, Schritte C1–C5);
  [`DC-FA-LINK-001`](../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
- Baseline-Regelwerk
  [`modul-05-planning-harness.md`](../../.harness/baseline/v5.0.0/regelwerk/modul-05-planning-harness.md)
  und [`modul-06-roadmap.md`](../../.harness/baseline/v5.0.0/regelwerk/modul-06-roadmap.md)
- die vier Change Requests: zwei aus dem Schwester-Repo a-check (dessen
  `slice-073`-§8, nur gelesen), zwei des Konsumenten `ai-harness-course`
  (im Auftrag übergeben, in keiner Datei dieses Repos)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules), Roadmap und
  [Beobachtungs-Register](../plan/planning/observations.md)

**Eigene Läufe** (alle über `docker run --rm --network none -v <fixture>:/repo:ro
d-check:latest`, Fixtures ausschließlich in einem Temp-Verzeichnis, kein
Host-`go`, das Repo nicht verändert):

- **L1** — Spiegel von `docs/plan/planning/` (104 Dateien in `done/`,
  Roadmap in `in-progress/`), Closure-Profil mit Default-`slice-glob`:
  `105 Datei(en) geprüft, 0 Befund(e)`, Exit 0.
- **L2** — derselbe Spiegel mit `slice-glob: "*.md"`: **11 Befunde**, Exit 1.
- **L3** — Minimal-Fixture zur Zähl- und Floskel-Semantik von v0.52.0.
- **L4** — `heading-pattern` mit Lookbehind/Lookahead (RE2-Probe).

---

## Findings

### F-1 — Inline-Code-Bereinigung trifft einen **geteilten** Vorfilter; die Lockerungs-Richtung ist weder benannt noch ADR-gedeckt

- `kategorie`: HIGH
- `quelle`: [`AGENTS.md` §3.6](../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden),
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
  (Entscheidung 3),
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `pfad`: `docs/plan/planning/open/slice-094-closure-zaehl-paritaet.md:22` (§1),
  `:42` (SemVer-Absatz), `:61` (DoD-Punkt 1)
- `befund`: Der Slice beschreibt die Änderung durchgängig als Eingriff in die
  *Substanz-Zählung* und leitet daraus „SemVer: Minor, kein Patch — d-check
  findet danach **mehr**" ab; die Spezifikation kennt aber **einen** bereinigten
  Abschnittstext, den `closure-note-thin` **und** `closure-note-boilerplate`
  teilen. In L3 belegt: eine Floskel innerhalb eines Inline-Code-Spans
  (`` `alles gut` ``) erzeugt heute `closure-note-boilerplate`, Exit 1 — nach
  einer Bereinigung, die Inline-Spans aus dem geteilten Text entfernt, wäre
  derselbe Bestand grün. Der DoD nennt Lastenheft-, Spezifikations- und
  Akzeptanzkriterien-Nachzug, aber **keine ADR**, obwohl §3.6 jede Senkung einer
  Prüfregel ADR-pflichtig macht.
- `verifizierbar`: ja — Fixture mit einem konfigurierten `boilerplate`-Eintrag,
  der ausschließlich innerhalb eines Backtick-Spans vorkommt; `make
  verify-closure-notes` bzw. `--enable planning` vor und nach der Änderung.
- `klasse`: „Geteilter Vorfilter einseitig geändert — Verschärfung dokumentiert,
  Lockerung nicht"

### F-2 — Der bidirektionale Paritäts-DoD von slice-094 ist im eigenen Schnitt nicht erfüllbar

- `kategorie`: MEDIUM
- `quelle`: [`DC-QA-04`](../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools),
  [§DC-FA-PLAN-001.a](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
- `pfad`: `docs/plan/planning/open/slice-094-closure-zaehl-paritaet.md:64` (DoD-Punkt 2)
- `befund`: Der DoD verlangt „jede Fixture, die das Adopter-Skript rot macht,
  macht auch das Modul rot — **und umgekehrt**". Das Adopter-Skript prüft neben
  der Satzzahl auch die **Kardinalität** der Closure-Abschnitte (mehr als einer
  ⇒ Befund, „ein zweiter ist typischerweise ein stehengebliebener Platzhalter")
  und führt dafür eine eigene Selbsttest-Fixture. Die Spezifikation bindet die
  Fähigkeit dagegen auf den **ersten** passenden Abschnitt; in L3 mit einer
  Datei aus gutem erstem und zurückgelassenem zweitem Closure-Abschnitt
  reproduziert: `0 Befund(e)`, Exit 0. slice-094 ändert nur die Zählung, der
  Fall bleibt danach grün. Weder slice-094 noch
  [slice-096](../plan/planning/in-progress/slice-096-structure-modul-analyse.md) §2
  („Nicht gedeckt sind die abschnitts-treue Task-Zählung und die benannten
  Pflicht-Bausteine") führen die Kardinalität als Lücke.
- `verifizierbar`: ja — die genannte Zwei-Abschnitts-Fixture gegen `--enable
  planning`; und `make verify-closure-notes` bleibt darauf grün.
- `klasse`: „Bidirektionale Parität zugesagt, Prüf-Kardinalität nicht im Schnitt"

### F-3 — Zwei Slices desselben Strangs vertreten gegensätzliche Haltungen zum Fremd-Skript, und die festgelegte Reihenfolge entscheidet den Streit vorweg

- `kategorie`: MEDIUM
- `quelle`: Maintainability; Baseline-Regelwerk `modul-05-planning-harness.md`
  §Ziel-Form: Slice (Schnitt nach Lieferwert)
- `pfad`: `docs/plan/planning/open/slice-096-structure-modul-analyse.md:63`
  (§3 Abnahme-Punkt 3), `:119` (§8) gegen
  `docs/plan/planning/open/slice-094-closure-zaehl-paritaet.md:64` (DoD-Punkt 2)
- `befund`:
  [slice-096](../plan/planning/in-progress/slice-096-structure-modul-analyse.md)
  verlangt, die fehlenden Bausteine „fence-treu und abschnitts-treu" **neu** zu
  definieren, „sonst wiederholt d-check den Fehler der abgelösten Skripte", und
  begründet den GF-Modus damit, dass die Fremd-Skripte „nicht rückdokumentiert,
  sondern durch eine eigene, spec-first formulierte Fähigkeit ersetzt" werden.
  [slice-094](../plan/planning/open/slice-094-closure-zaehl-paritaet.md) sagt
  für dieselbe Familie das Gegenteil zu: Deckungsgleichheit mit dem Skript in
  **beide** Richtungen. Da der Start-Trigger von slice-096 „slice-094 in
  `done/`" lautet, ist die Paritäts-Haltung zum Zeitpunkt der Abnahme-Frage
  bereits als Vertrag ausgeliefert und per SemVer-Minor bei Konsumenten
  angekommen.
- `verifizierbar`: nein — Plan-Konsistenz, kein Gate-Lauf entscheidet sie.
- `klasse`: „Reihenfolge entscheidet eine offene Abnahme-Frage vorweg"

### F-4 — Der Alias-/Supersede-Pfad überträgt ADR-0044 auf einen Fall, der in zwei Punkten anders liegt

- `kategorie`: MEDIUM
- `quelle`: [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md)
  (Entscheidungen 1 und 6),
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
  (Entscheidungen 1 und 3),
  [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `pfad`: `docs/plan/planning/open/slice-096-structure-modul-analyse.md:56`
  (§3 Abnahme-Punkt 2)
- `befund`: Der Präzedenzfall trägt in zwei Punkten nicht so weit, wie der
  Abnahme-Punkt ihn spannt. (a) **Ausgabe-Fläche.** Der Alias in
  [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md) war
  ausdrücklich *byte-identisch* — dieselbe Unterdrückung, dieselben Befunde. Hier
  meldet die Quelle `closure-note-missing` / `closure-note-thin` /
  `closure-note-boilerplate` und das beantragte Zielmodul `section-missing` /
  `section-empty` / `section-constraint`; ein Alias auf der Config-Achse lässt
  offen, welches Grund-Code-Vokabular ein Adopter danach in seiner Ausgabe sieht
  und in seinen Ventilen adressiert. (b) **Supersede-Ziel.** ADR-0044 hat die
  Host-Anforderung **nicht** superseded, sondern nur das Ventil in ein neues
  geteiltes Kürzel gehoben. Der Abnahme-Punkt formuliert dagegen „die bestehende
  Anforderung wird per ADR **superseded**" — die bestehende Anforderung trägt
  aber neben der Closure-Fähigkeit auch die Aktiv-Status-Invariante
  (Roadmap ↔ `in-progress/`), die das beantragte Modul nicht abdeckt (Struktur
  *innerhalb* eines Dokuments statt Konsistenz *zwischen* Artefakten).
- `verifizierbar`: ja für (a) — jeder Lauf gegen ein Adopter-Repo nach der
  Ablösung zeigt, welches Vokabular gemeldet wird; nein für (b).
- `klasse`: „Alias-Präzedenz übertragen, ohne Ausgabe-Fläche und Supersede-Ziel
  mitzuprüfen"

### F-5 — Die Bestandsmessung in slice-097 ist auf eine Befundklasse verkürzt; der geweitete Glob erzeugt zehn Closure-Befunde in zwei Klassen

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in);
  Reviewer-Skill §MEDIUM (Spec-Treue-Lücke einer Messmethode)
- `pfad`: `docs/plan/planning/open/slice-097-closure-glob-entkopplung.md:59`
  (§3 Abnahme-Punkt 2)
- `befund`: Der Abnahme-Punkt nennt „**9** `welle-*-results.md` … ⇒ neun
  `closure-note-missing`" und leitet daraus genau drei Optionen ab (Muster auf
  `^#{1,3}` weiten · die Wellen-Notizen umbauen · den Glob eng lassen). Lauf L2
  meldet **11** Befunde: neun `closure-note-missing` (die Ergebnis-Notizen,
  H1-Überschrift — insoweit bestätigt), ein `planning-drift` (das Falsch-Rot der
  Roadmap-Invariante, ebenfalls bestätigt) **und** ein zehnter Closure-Befund
  `closure-note-thin` auf `docs/plan/planning/done/welle-68-planning-roadmap-harness.md:58`,
  dessen §7 im Ruheort `done/` noch `_Ausstehend._` trägt. Dieser gehört einer
  anderen Klasse an (Substanz, nicht Überschriften-Ebene) und wird von keiner
  der drei genannten Optionen berührt. Der Bestand führt außerdem **zwei**
  Wellen-Plandateien in `done/`, die die Zählung „9" nicht kennt.
- `verifizierbar`: ja — L2 ist die Reproduktion; im Repo äquivalent über das
  Closure-Profil mit gesetztem weiten Glob.
- `klasse`: „Bestandsmessung auf eine Befundklasse verkürzt"

### F-6 — Der gemeinsamen Welle 097+098 fehlt in beiden DoDs der Schritt, der ihren erklärten Zweck einlöst

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
  braucht und §Welle ≠ Meilenstein ≠ Release
- `pfad`: `docs/plan/planning/open/slice-097-closure-glob-entkopplung.md:69`
  (§4) und `docs/plan/planning/open/slice-098-closure-note-placeholder.md:73` (§4)
- `befund`: Beide Slices begründen die gemeinsame Welle damit, dass der
  Konsument sein Skript „erst dann zurückziehen" kann bzw. die Fähigkeit „erst
  mit beiden" eine Obermenge ist. Beide DoDs enden bei „`make gates` +
  `make verify-closure-notes` grün"; ein Release-Schritt fehlt, obwohl beide die
  öffentliche Konfigurations-Fläche erweitern und die Fähigkeit den Konsumenten
  nur über ein veröffentlichtes Image erreicht. Die beiden anderen
  Produkt-Slices desselben Backlogs führen ihn ausdrücklich
  ([slice-094](../plan/planning/open/slice-094-closure-zaehl-paritaet.md) §4:
  „Release als **Minor**";
  [slice-095](../plan/planning/done/slice-095-links-resolve-from.md) §4:
  „Release als Minor"). Eine Welle-DoD, die die Lücke auffangen könnte, gibt es
  nicht — die Welle ist angekündigt, nicht eröffnet.
- `verifizierbar`: nein — DoD-Vollständigkeit ist Plan-Eigenschaft.
- `klasse`: „Wellen-Zweck extern, Wellen-DoD ohne den Schritt, der ihn erreicht"

### F-7 — slice-094 macht eine Fremd-Repo-Fixture-Menge zur DoD-Bedingung, ohne sie als Risiko zu führen

- `kategorie`: LOW
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Offene Risiken
  werden bei Closure aufgelöst
- `pfad`: `docs/plan/planning/open/slice-094-closure-zaehl-paritaet.md:64`
  (DoD-Punkt 2) gegen `:69` (§5)
- `befund`: Der DoD verlangt den Paritäts-Beleg „gegen die a-check-Fixtures";
  §5 führt nur Bestands-Rot und Konsumenten-Bruch als Risiken. Die
  Fixture-Menge liegt nicht in diesem Repo — im Antragsteller-Repo existiert sie
  als Selbsttest-Block **innerhalb** des abzulösenden Skripts, nicht als
  beistellbares Artefakt mit eigener Kennung. Die einzige Rückführung des Slice
  (`in-progress` → `next`) ist auf den Fall „Bestands-Sanierung" formuliert und
  greift für eine ausbleibende Fixture-Beistellung nicht.
  [slice-096](../plan/planning/in-progress/slice-096-structure-modul-analyse.md) §5
  führt für dieselbe Menge genau dieses Risiko („Fremd-Repo-Abhängigkeit …
  beizuziehen, nicht nachzubauen").
- `verifizierbar`: nein
- `klasse`: „Blockierende Fremd-Abhängigkeit ohne Risiko-Ausgang"

### F-8 — Die angekündigten Wellen existieren nur in den Slice-Köpfen; die Vorschau-Fläche der Roadmap ist leer, und ein Wellen-Trigger zeigt auf einen wellenlosen Slice

- `kategorie`: LOW
- `quelle`: Baseline-Regelwerk `modul-06-roadmap.md` §Roadmap-Struktur (Bullet
  *Nächste Wellen*: Abhängigkeit als beobachtbare Bedingung in der
  `Trigger`-Spalte **und** als gerichtete Kante im Abhängigkeitsgraphen) sowie
  §Wann Arbeit eine Welle braucht (wellenlose Arbeit erscheint nicht in der
  Roadmap)
- `pfad`: `docs/plan/planning/in-progress/roadmap.md:32` und `:50`;
  `docs/plan/planning/open/slice-096-structure-modul-analyse.md:7` (Kopf-Feld)
  und `:97` (§6)
- `befund`: Drei Slices deklarieren im Kopf-Feld zwei künftige Wellen
  (eine für den `structure`-Strang, eine gemeinsame für 097+098). Die Roadmap
  führt in *Nächste Wellen* „— keine geplante Welle —" und im
  Abhängigkeitsgraphen „keine geplante Folge-Welle". Zugleich ist der
  Start-Trigger der ersten Welle „slice-094 in `done/`" — slice-094 ist als
  wellenlose Arbeit deklariert und erscheint nach derselben Regel nicht in der
  Roadmap; die vom Regelwerk verlangte gerichtete Kante hätte damit keinen
  Knoten, auf den sie zeigen kann.
- `verifizierbar`: nein — `planning-drift` prüft nur die Aktiv-Status-Kopplung,
  nicht die Vorschau-Sektion (in L1 grün).
- `klasse`: „Wellen-Zusage nur im Slice-Kopf, Vorschau-Fläche leer"

### F-9 — `planning.closure.glob` geht in dem beantragten Modul nicht als Semantik, sondern als Achse auf

- `kategorie`: INFO
- `quelle`: [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md)
  (Kürzel-Kriterium: querschnittlich ⇒ eigenes Kürzel); die zwei Change
  Requests
- `pfad`: `docs/plan/planning/open/slice-097-closure-glob-entkopplung.md:24` (§1)
- `befund`: Die Begründung, vor dem Modul-Schnitt zu bauen, lautet für alle drei
  Vorab-Slices gleich („zwei Konsumenten brauchen sie jetzt; die Semantik trägt
  in `structure` weiter"). Für die Zähl-Regel (slice-094) und den neuen
  Grund-Code (slice-098) trifft das zu — beides sind Aussagen *über einen
  Abschnitt*, die ein Nachfolgemodul unverändert übernehmen kann. Der
  Kandidaten-Filter ist von anderer Art: er ist die **Datei-Auswahl-Achse**, die
  das beantragte Modul konstitutiv je Regel führt (Pfad-Glob `files:`), während
  der Slice ihn ausdrücklich als **Basisname**-Glob neben `closure.dir`
  einführt („kein Pfad-Glob"). Die Übertragung ist damit keine Übernahme,
  sondern eine Umrechnung über zwei Glob-Domänen — und der Schlüssel ist im
  Zielmodul per Konstruktion redundant. Der Punkt ist in keinem der beiden
  Slices festgehalten.
- `verifizierbar`: nein — Entwurfs-Beobachtung.
- `klasse`: „Vorab-Fähigkeit dupliziert die Achse des Nachfolge-Moduls"

## Negativbefunde

- **geprüft, ohne Befund — Schnitt von slice-095:** einziger Slice außerhalb des
  Closure-Strangs; berührt eine andere Anforderung
  ([`DC-FA-LINK-001`](../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)),
  ist von 094/096 unabhängig deklariert und verortet die Erweiterung nach dem
  Kürzel-Kriterium aus
  [ADR-0044](../plan/adr/0044-geteiltes-referenz-ventil-quell-skopus.md)
  korrekt in der bestehenden Anforderung (Einzelmodul, keine neue Frage).
- **geprüft, ohne Befund — CR-Abdeckung (Vollständigkeit):** die vier Change
  Requests sind 1:1 auf Slices abgebildet (a-check-CR 1 → slice-096, a-check-CR 2
  → slice-095, Konsumenten-CR 1 → slice-097, Konsumenten-CR 2 → slice-098);
  slice-094 ist kein CR, sondern eine selbst gemessene Paritätslücke. Die vier
  Akzeptanzkriterien von Konsumenten-CR 2 und die drei von Konsumenten-CR 1
  stehen wörtlich in den jeweiligen DoDs, einschließlich der Nachfilter und der
  fail-closed-Zusagen.
- **geprüft, ohne Befund — Doppelung zwischen den Slices:** kein Slice sagt eine
  Lieferung zu, die ein anderer ebenfalls zusagt. Die einzige Berührung ist der
  Paritäts-Beleg (slice-094 für die Zählung, ein Folge-Slice aus slice-096 für
  das Gesamtmodul) — verschiedene Gegenstände, keine doppelte Zusage.
- **geprüft, ohne Befund — Wellen-Regel bei 094 und 095:** beide sind einzelne
  Slices, deren Closure-Trigger die eigene DoD abschriebe; die repo-weiten
  Gate-Bedingungen stehen bereits *in* ihren DoDs. „Ohne Welle" ist die
  Anwendung der Regel, nicht ihre Umgehung.
- **geprüft, ohne Befund — Sichtungs-Schritt (§7) in allen fünf Slices:** jeder
  Plan geht das [Beobachtungs-Register](../plan/planning/observations.md) durch
  und notiert das Ergebnis; slice-096 hält die naheliegende Verwechslung
  ausdrücklich fest, statt sie stillschweigend zu verneinen. Das genügt
  `modul-05` §Zwei Schritte („keine Treffer sind ebenfalls eine Antwort").
- **geprüft, ohne Befund — WIP und Lifecycle-Lage:** `in-progress/` enthält
  keine Slice-Datei, alle fünf liegen in `open/`; die Roadmap trägt den
  Ruhe-Marker. L1 bestätigt die Kopplung mit `0 Befund(e)`.
- **geprüft, ohne Befund — SemVer-Einordnung von slice-094 in der Richtung, die
  sie benennt:** „findet mehr ⇒ Minor, kein Patch" deckt sich wörtlich mit
  [ADR-0042](../plan/adr/0042-markdown-lexik-folgt-commonmark.md) (Entscheidung
  6: „Ein Konsumentenlauf, der heute grün ist, kann danach rot sein … Das ist
  kein Patch") und
  [ADR-0043](../plan/adr/0043-tabellengrenze-am-relevanten-header.md)
  („SemVer-Minor: d-check **findet mehr**"). Beide verlangen zusätzlich die
  Release-Notiz, die der DoD führt. Die Einordnung ist richtig; unvollständig
  ist nur die zweite Richtung (F-1).
- **geprüft, ohne Befund — Messwert (a), neun Wellen-Ergebnisnotizen:** `done/`
  führt genau neun `welle-*-results.md`, alle neun tragen eine
  H1-Closure-Überschrift der Form `# Welle NN — … — Closure-Notiz`. Die Zahl und
  die Überschriften-Ebene stimmen; unvollständig ist die daraus abgeleitete
  Befundzahl (F-5).
- **geprüft, ohne Befund — Messwert (b), Falsch-Rot der Roadmap-Invariante:** in
  L2 exakt reproduziert — `docs/plan/planning/in-progress/roadmap.md:12`,
  `planning-drift`, verursacht dadurch, dass die Roadmap-Datei selbst im
  gezählten Verzeichnis liegt und `*.md` matcht, während der Ruhe-Marker korrekt
  gesetzt ist. Die Ursachenbeschreibung des Slice trifft zu.
- **geprüft, ohne Befund — Messwert (c), RE2-Ablehnung:** L4 bestätigt beide
  Lookaround-Formen fail-closed. Lookbehind ⇒ `error parsing regexp: invalid
  named capture`, Lookahead ⇒ `error parsing regexp: invalid or unsupported Perl
  syntax`, jeweils Exit 2 vor dem Lauf. Die im Slice zitierte Fehlerzeile
  entspricht dem tatsächlichen Format (bis auf den Profil-Dateinamen, der bei
  einem Lauf über das Closure-Profil abweicht).
- **geprüft, ohne Befund — Messwert aus slice-094 §2:** die Notiz „Siehe `a.md`
  und `b.md`." zählt in v0.52.0 drei Satzende-Zeichen und bleibt bei Schwelle 2
  grün (L3, mit Schwelle 3 grün, mit Schwelle 4 `closure-note-thin`). Die
  Zähl-Regel des Adopter-Skripts (Fenced-Blöcke vollständig entfernen,
  Inline-Spans entfernen, nur Satzzeichen vor Whitespace/Zeilenende zählen) ist
  im Skript so implementiert wie im Slice beschrieben.
- **geprüft, nicht widerlegt — Bestands-Risiko von slice-094:** eine
  Näherungs-Nachbildung der neuen Zählung über die 93 `done/`-Slices ergibt
  keine Datei unter der Schwelle 4; die Nachbildung liegt für die *heutige*
  Regel jedoch bei Minimum 7, wo d-check 5 misst, weicht also im Parser ab. Das
  in §5 offen geführte Risiko ist damit weder belegt noch ausgeräumt — die im
  Slice vorgesehene Vorab-Messung bleibt nötig.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 5 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Geteilter Vorfilter einseitig geändert —
Verschärfung dokumentiert, Lockerung nicht · Bidirektionale Parität zugesagt,
Prüf-Kardinalität nicht im Schnitt · Reihenfolge entscheidet eine offene
Abnahme-Frage vorweg · Alias-Präzedenz übertragen, ohne Ausgabe-Fläche und
Supersede-Ziel mitzuprüfen · Bestandsmessung auf eine Befundklasse verkürzt ·
Wellen-Zweck extern, Wellen-DoD ohne den Schritt, der ihn erreicht ·
Blockierende Fremd-Abhängigkeit ohne Risiko-Ausgang · Wellen-Zusage nur im
Slice-Kopf, Vorschau-Fläche leer · Vorab-Fähigkeit dupliziert die Achse des
Nachfolge-Moduls

## Verdikt

**Merge-blockierend:** ja — hier gelesen als: die Freigabe von `open/` nach
`next/` sollte für den Closure-Strang nicht ohne Klärung von F-1 bis F-4
erfolgen. F-1 betrifft eine Prüfregel im Gate-Pfad und die Frage, ob der
geplante Release eine ADR braucht; F-2 macht einen DoD-Punkt unerfüllbar, F-3
und F-4 entscheiden über Reihenfolge und Ablöse-Pfad. F-5 bis F-9 sind vor der
jeweiligen Umsetzung zu klären, blockieren die Freigabe der übrigen Slices aber
nicht: [slice-095](../plan/planning/done/slice-095-links-resolve-from.md) ist
vom Strang unabhängig und trägt keinen Befund.

Zur ausdrücklich gestellten Schnitt-Frage, soweit sie nicht in ein Finding
fällt: der Schnitt in vier statt einen Slice ist tragfähig — ein einzelner
`structure`-Slice überschritte die Größenregel aus `modul-05` (drei
Liefer-Punkte, zwei Schichten) deutlich, und die drei Vorab-Fähigkeiten sind
einzeln lieferbar. Die Begründung „die Semantik trägt in `structure` weiter"
gilt für zwei der drei (F-9). Die Reihenfolge „094 vor 096" ist als
Messgrundlagen-Argument nachvollziehbar, entscheidet aber eine Vertragsfrage
mit, die 096 sich vorbehält (F-3).

**Übergabe:** Findings gehen an den Planner (Rückkante Review → Plan); die
**Finding-Klassen** gehen bei der Closure der betroffenen Slices in
§7 und von dort ins [Beobachtungs-Register](../plan/planning/observations.md).
Dieser Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation — die
DoD-/Spec-Konformität der einzelnen Slices prüft der Verifier separat, und ein
Detail-Plan-Review der Slice-Texte steht praxisgemäß erst beim Übergang nach
`in-progress/` an.
