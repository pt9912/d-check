# Slice slice-228: Verbliebenen Bestand slice-200 bis slice-227 archivieren

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos**. Reine Bestandspflege mit einem bereits gebauten
und getesteten Werkzeug (`tools/archive-wave -slice=<id>`, slice-196), kein
neuer Produktcode; der Closure-Grund geht über die eigene DoD nicht hinaus
(Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** slice-200 (Präzedenz — identischer Vorgang für den damaligen
Backlog, 2026-09-04: Dry-Run über den ganzen Bestand vor der Anwendung, je
Slice ein Lauf, ein gebündelter Commit). slice-197/slice-199 (ursprüngliche
Werkzeug-Herkunft, welle-89).

**Berührte Spec-Stellen:** — (Bestandspflege, keine neue Anforderung).

**Verantwortlich:** claude-sonnet-5.

**Autor:** pt9912. **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die 28 flachen wellenlosen `done/`-Slices `slice-200` bis
`slice-227` (samt ihren 37 zugehörigen Review-Reports unter `docs/reviews/`)
per `tools/archive-wave -slice=<id> -apply` nach
`docs/plan/planning/done/wellenlos/` archivieren — der Backlog, der seit dem
letzten Sweep (slice-200, 2026-09-04) angefallen ist.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Änderung des `archive-wave`-Werkzeugs selbst** — anderer Vorgang:
  eine reine Bestandspflege mit einem bereits gebauten und getesteten
  Werkzeug ändert dessen Verhalten nicht. Das gilt ausdrücklich auch für die
  bekannte, in `BEO-ALL/batch-slice-archival-zips-post-rewrite-content`
  offen geführte Order-Abhängigkeit des Zip-Inhalts bei sequenzieller
  Batch-Archivierung (siehe §6) — eine Behebung wäre Produktcode-Arbeit am
  Werkzeug, kein Bestandspflege-Akt.
- **`docs/reviews/archiv/`** — bereits archivierter Bestand bleibt bewusst
  stehen; dieser Slice sammelt nur den flachen Rest ein.
- **Eine nachträgliche Korrektur der Scope-Abgrenzung von slice-197** —
  dessen eigene §3 dort bleibt stehen, wie bereits von slice-200 begründet;
  dieser Slice schließt nur die seither neu entstandene Lücke.
- **Der `REVIEW=`-Modus für eigenständige Reviews** — entfällt: geprüft
  (siehe §3), kein flacher Review-Report unter `docs/reviews/` liegt
  außerhalb des Bereichs `slice-200`–`slice-227`.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] Alle 28 wellenlosen Slices `slice-200`–`slice-227` archiviert (Stub +
      `archiv.zip` je Slice unter `docs/plan/planning/done/wellenlos/`,
      inklusive ihrer 37 zugehörigen Review-Reports); kein flacher Rest aus
      diesem Bereich mehr unter `docs/plan/planning/done/` oder
      `docs/reviews/`.
- [x] `make gates` grün auf dem Endstand.
- [ ] `make fullbuild` grün auf dem Endstand.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: — kein öffentlicher Vertrag berührt (reine
      Bestandspflege, keine Schnittstellen-Änderung).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — zweites
      Auftreten in `BEO-ALL/batch-slice-archival-zips-post-rewrite-content/`
      (siehe §6).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      Repo ohne Wellen-Betrieb, hier geprüft.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area?

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/planning/done/slice-{200..227}-*.md` | move (Werkzeug) | Stub ersetzt Volltext, `tools/archive-wave -slice=<id> -apply` |
| `docs/reviews/2026-09-0{4,6,7,8}-slice-{200..227}-*.md`, `2026-09-1{6,7}-slice-{219..227}-*.md` (37 Dateien) | move (Werkzeug) | ins jeweilige `archiv.zip` des Slice eingesammelt, kein eigener Stub |
| `docs/plan/planning/done/wellenlos/slice-{200..227}-*.md` + `-archiv.zip` | neu (Werkzeug) | Ziel der Move-Operation |
| repo-weite Querverweise auf die bewegten Pfade | update (Werkzeug) | `RewriteRepo()` zieht sie automatisch nach |
| neue Evidence-Datei in `BEO-ALL/batch-slice-archival-zips-post-rewrite-content/evidence/` | neu | zweites Auftreten der bekannten Order-Abhängigkeit (§6) |
| `docs/plan/planning/in-progress/slice-228-*.md` → `docs/plan/planning/done/` | neu, dann move | dieser Slice-Plan selbst |

**Vorab geprüft:** kein flacher Review-Report unter `docs/reviews/` liegt
außerhalb `slice-200`–`slice-227` (37 Treffer, 0 Ausreißer — Grep gegen alle
37 Dateinamen).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Nutzer-Anfrage „Warum werden die Slices
in docs/plan/planning/done nicht archiviert?" beantwortet, Nutzer bestätigt
Ausführung explizit („ja", 2026-09-17); WIP-Limit frei (`open/`, `next/`,
`in-progress/` waren leer bis auf `roadmap.md`).

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls der
  Voll-Dry-Run über den 28er-Bestand einen strukturellen Fehlerpfad zeigt,
  der eine Werkzeug-Änderung verlangt (siehe §1-Ausschluss) — dann ist der
  Umfang größer als reine Bestandspflege.
- `in-progress` → `open` (blockiert): falls `make gates` oder
  `make fullbuild` auf dem archivierten Endstand strukturell rot bleibt und
  kein Carveout trägt.

## 5. Closure-Trigger

DoD vollständig, `make gates` und `make fullbuild` auf dem Endstand grün,
unabhängiger Review und unabhängige Verifikation abgeschlossen,
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Order-Abhängigkeit des Zip-Inhalts bei sequenzieller Batch-Archivierung**
  ohne Zwischen-Commit — dieselbe Klasse, die der unabhängige Review von
  slice-200 fand (`docs/plan/planning/observations/BEO-ALL/batch-slice-archival-zips-post-rewrite-content/`,
  bislang 1×): ein früherer Lauf in diesem Batch kann per `RewriteRepo()`
  einen Verweis in einem noch nicht archivierten, später folgenden Slice
  umschreiben, bevor dieser selbst gezippt wird — inhaltlich harmlos
  (korrekter Pfad-Nachzug), aber Zip-Inhalt wird von der
  Verarbeitungsreihenfolge abhängig (Determinismus-Anliegen,
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
  Bereits im Dry-Run von `slice-206` beobachtet (Verweis-Fix auf
  `slice-205` angekündigt, das selbst noch im selben Batch archiviert
  wird). **Ausgang:** weiter offen → zitiert
  `BEO-ALL/batch-slice-archival-zips-post-rewrite-content` (2. Auftreten,
  nicht neu angelegt).
- **Verwaiste `exempt-paths`/Konfig-Einträge**, die auf jetzt archivierte
  flache Pfade zeigen (dieselbe Klasse wie slice-200 Fund F-4). **Ausgang:**
  eingetreten — zehn gegenstandslose `ignore-refs`-Einträge in
  `.d-check.yml` gefunden (Kommit `14b1527c`) und im selben Slice entfernt,
  kein Carveout/Folge-Slice nötig; `make doc-check`/`make gates` bestätigen
  den bereinigten Stand grün.
- **`make fullbuild`-DoD-Haken-Timing**: das Closure-Profil des
  nicht-rekursiven `done/slice-*.md`-Globs kann vor dem `git mv` dieses
  Slices selbst auf eine andere Menge treffen als danach (dasselbe Muster
  wie bei slice-200 dokumentiert). **Ausgang:** entfallen — der Haken wird
  wie bei slice-200 erst nach dem bestätigten Endstand gesetzt, dieselbe
  Vorsichtsmaßnahme wird von vornherein übernommen.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — wird vor dem `git mv` nach done/ gefüllt. -->

- **Was hat funktioniert:** Der Dry-Run-vor-Apply-Ablauf aus slice-200 trug
  unverändert: 28 Einzel-Läufe (`make archive-wave SLICE=<id> APPLY=1`),
  `make gates` und `make doc-check` auf dem Endstand grün, ein gebündelter
  Commit für den reinen Move.
- **Was ging anders als geplant:** Die im §6-Risiko vorab benannte
  Order-Abhängigkeit (Zip-Inhalt bei sequenzieller Archivierung) trat wie
  erwartet erneut auf (mehrfach, als **ein** Vorgang gezählt). Zusätzlich,
  nicht vorab benannt: die Archivierung machte zehn `ignore-refs`-Einträge
  in `.d-check.yml` gegenstandslos (dieselbe Klasse wie slice-200 F-4,
  aber ein anderer Konfig-Block als dort — dort `reviews.exempt-paths`,
  hier `ignore-refs`); im selben Slice bereinigt (Kommit `14b1527c`).
- **Beobachtungs-Register (`../observations/`):**
  `evidence/slice-228.md` in
  `BEO-ALL/batch-slice-archival-zips-post-rewrite-content/` ergänzt —
  Zähler steht damit bei 2×.
- **Folge-Slices:** keine.
- **Risiken aus §6:** alle drei mit Ausgang — siehe §6 (1× weiter offen,
  1× eingetreten/behoben, 1× entfallen).
- **Drei Paarungen:** wird nach dem `git mv` dieses Slice-Plans geprüft
  (siehe DoD).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**;
bedingt ist allein der Modus-Block am Ende.

Dieses Repo führt **drei** Prüfungen — die zwei kanonischen und, als
Adaption, den Nachtlauf-Stand ([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Default) — `docs/plan/planning/`, `docs/reviews/`
und `docs/plan/planning/observations/` tragen keine eigene
Modus-Deklaration in `harness/conventions.md` und fallen damit unter den
Default. `tools/archive-wave/` wird nur **aufgerufen**, nicht geändert
(§1-Ausschluss) — keine zweite Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, 46 Verzeichnisse). **Ein** Treffer
für Sub-Area `*`, direkt einschlägig:
`BEO-ALL/batch-slice-archival-zips-post-rewrite-content` (1× aus slice-200)
— siehe §6.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-17 gelesen: `image-scan.yml` grün
(2026-09-17T08:41:27Z). `upstream-drift.yml` **rot**
(2026-09-17T05:38:51Z) — zwei Teilziele betroffen (`freshness-semgrep`,
`go-base-digest`), beide Fremd-Bestand-Freshness-Meldungen ohne Bezug zu
Planning/`archive-wave`; kein Bezug zu diesem Slice.

**Modus-Begründung:** die einzige berührte Sub-Area ist GF (Default) — kein
Begründungsblock nötig.
