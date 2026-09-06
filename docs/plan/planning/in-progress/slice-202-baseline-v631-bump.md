# Slice slice-202: Baseline-Pin-Hebung auf v6.3.1

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre. Die Adoptions-Stränge des Deltas sind
Folge-Slices mit eigenen DoDs, kein repo-weites *Mehr* über diesen Slice
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette
(Pin-Fortschreibung), [`MR-021`](../../../../harness/conventions.md#mr-021)
(pin-gebundene Verweise), [`MR-051`](../../../../harness/conventions.md#mr-051)
(cite-Spannen), [`MR-055`](../../../../harness/conventions.md#mr-055)
(Symlink-Aliase), [`MR-060`](../../../../harness/conventions.md#mr-060)
(der abzulösende Vorgänger).

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle; der
vendorte Baum und der Konventionsspeicher liegen außerhalb der Spec-Straten).

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-06.

---

## 1. Ziel

Den vendorten Baseline-Baum von `v6.0.0` auf **`v6.3.1`** heben — vier
Releases, deren Inhalt gelesen statt angenommen wird —, alle pin-gebundenen
Verweise nachziehen und die Adoptions-Entscheidung je Delta-Strang
festhalten. **Abgrenzung:** Dieser Slice hebt den Stand und
entscheidet je Delta-Strang, ob und wie adoptiert wird; die **Umsetzung** der
vier Stränge in der verkörperten Form ist slice-203 — sie setzt die mit diesem
Bump eintreffenden Vorlagen voraus.

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/baseline/v6.3.1/` | neu | der gehobene vendorte Baum (Regelwerk + Templates + `SHA256SUMS`) |
| `.harness/baseline/v6.0.0/` | entfällt | sonst stiller Stale-Content neben dem neuen Pin ([`MR-021`](../../../../harness/conventions.md#mr-021)) |
| `.claude/rules/*.md` (Baseline-Aliase) | update | Symlink-Ziele auf den neuen Tag ([`MR-055`](../../../../harness/conventions.md#mr-055)) |
| `harness/conventions.md` | update | §Baseline, §Adoptierte Konventions-Quellen, beide Index-Tabellen |
| neue Eintrags-Datei unter `harness/conventions/` | neu | der Hebungs-Eintrag [`MR-065`](../../../../harness/conventions.md) |
| Eintrags-Datei von [`MR-060`](../../../../harness/conventions.md#mr-060) | `git mv` → `done/` | aufgelöst durch den Nachfolger |
| `AGENTS.md`, `harness/README.md`, `spec/*`, Skills, Planning-Docs | update | pin-gebundene Verweise (Menge bestimmt der Zensus) |
| `.d-check.yml` | update | Tombstone-Globs für den entfernten Baum |

## 3. Ausdrücklich NICHT in diesem Slice

- **Die Umsetzung der vier Delta-Stränge in der verkörperten Form.** Sie ist
  [slice-203](../open/slice-203-v631-template-adoption.md) — allen voran die
  Auslagerung überladener Sensors-Einträge, die d-checks eigene §Sensors-Tabelle
  unmittelbar betrifft. Die Abgrenzung ist eine **Reihenfolge-Bedingung**, keine
  Bequemlichkeit: die neuen Vorlagen liegen erst nach diesem Bump im Repo.
- **Das Nachrüsten der eingefrorenen Verweise.** `done/`-Slices, Reviews,
  immutable ADRs, aufgelöste `MR`-Einträge und der ausgehende CR zitieren den
  Stand ihrer Zeit und werden nicht gehoben.
- **Jede Änderung an Produkt-Code.** Der Bump berührt Harness-Form und
  Konventionsspeicher, nicht `internal/` oder `cmd/`.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; Gate-Läufe und die
Closure-Pflichten darunter zählen nicht mit.

- [x] **(1)** Vendorter Baum steht auf `v6.3.1`: alter `v6.0.0`-Baum entfernt,
      neuer materialisiert samt `SHA256SUMS`, die acht Baseline-Aliase unter
      `.claude/rules/` umgehängt; `make baseline-verify` grün (Integrität,
      Manifest-Deckung, Alias-Auflösung).
- [x] **(2)** Alle **lebenden** pin-gebundenen Verweise auf `v6.3.1` gehoben
      (Zensus nach [`MR-021`](../../../../harness/conventions.md#mr-021),
      inklusive der beiden benannten Übersetzungsfehler-Klassen:
      `../baseline/` <!-- d-check:ignore (Muster-Praefix, illustriert die Fehlerklasse) -->-relative Pfade in `.harness/skills/` und Release-URLs
      ohne `.harness/baseline/`-Segment); `d-check:cite`-Spannen geprüft und
      wo nötig neu verankert
      ([`MR-051`](../../../../harness/conventions.md#mr-051));
      `.d-check.yml`-Tombstone für den entfernten Baum ergänzt.
- [x] **(3)** [`MR-065`](../../../../harness/conventions.md) geschrieben —
      mit Delta-Tabelle je Tag, Adoptions-Entscheid je Strang und
      Adaptions-Review über alle aktiven Einträge —,
      [`MR-060`](../../../../harness/conventions.md#mr-060) nach
      `conventions/done/` bewegt, beide Index-Tabellen fortgeschrieben.
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder eine weitere Datei in einem vorhandenen `evidence/`;
      kein Zähler wird gesetzt, er folgt aus den Dateien.
- [x] Jedes Risiko aus §5 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      im wellenlosen Betrieb hier geprüft, nach dem `git mv`.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Vier Tags auf einmal überschreiten das Review-Limit einer Sitzung**
  ([`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md),
  Zähler 2×). — **Ausgang:** entfallen — das gemessene Delta war klein (13
  Dateien mit echtem Inhalt), und der unabhängige Review prüfte alle vier
  Commits in **einem** Lauf mit vollständigem Verdikt und eigenen
  Gate-Läufen. Die Grenze wurde nicht erreicht.
- **Der Zensus übersieht eine Spiegel-Klasse**
  ([`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md),
  Zähler 5×; kein Sensor hält die Menge). — **Ausgang:** entfallen — der
  Review prüfte die Vorwärts-Richtung eigens (Negativbefund 3): kein lebender
  Verweis trägt mehr `v6.0.0`. Übersehen wurde nichts; der Fehler lag in der
  **Gegenrichtung** und ist beim nächsten Punkt geführt.
- **Der mechanische Rewrite trifft eine eingefrorene Klasse**
  ([`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md),
  Zähler 1×). — **Ausgang:** eingetreten — als Review-Befund F-4: der Zensus
  hob die §2-Zeile mit, die den zu **entfernenden** Baum identifiziert.
  Behoben vor der Closure; Beleg `evidence/slice-202.md` im Register, Zähler
  steht damit bei 2×. Neue Ausprägung, benannt: nicht ein eingefrorenes
  *Verzeichnis*, sondern eine identifizierende Nennung in einem **lebenden**
  Dokument — kein Gate meldet sie, weil beide Pfade auflösen.
- **Die Adoptions-Stränge werden als „mitgemacht" behauptet, obwohl nur der
  Pin steht.** — **Ausgang:** entfallen — der Review prüfte die Frage eigens
  (Negativbefund 14) und fand keine solche Behauptung: die Abgrenzung ist in
  §1, DoD (2)/(3) und [`MR-065`](../../../../harness/conventions.md#mr-065)
  gleichlautend als Reihenfolge-Bedingung geführt. Der Folge-Slice
  [slice-203](../open/slice-203-v631-template-adoption.md) liegt in `open/`.

## 6. Trigger

**Start** (`open` → `in-progress`): der Kurs hat `v6.3.1` publiziert und
`make baseline-freshness` meldet den Rückstand — beides am 2026-09-06
gemessen. WIP-Limit frei (alle Lifecycle-Verzeichnisse leer).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt der Zensus, dass ein Delta-Strang
  nicht als Folge-Slice abtrennbar ist, sondern die Hebung selbst blockiert
  (etwa weil eine aktive Adaption dem neuen Stand widerspricht), geht der
  Slice zurück zur Zerlegung — Hebung und Adoption werden dann getrennte
  Slices.
- `in-progress` → `open` (blockiert): Erweist sich das Release als
  zurückgezogen oder sein Bundle als integritäts-fehlerhaft
  (`sha256sum -c` rot), ruht der Slice bis zu einem tragfähigen Tag.

**Closure-Trigger.**

Zwei beobachtbare Kriterien **und** ein Lerneintrag: (a) `make gates` grün
auf dem Endstand, `make baseline-verify` bestätigt Integrität und
Alias-Auflösung gegen `v6.3.1`, `make baseline-freshness` meldet den Pin als
aktuell; (b) kein lebender Verweis trägt mehr `v6.0.0` — gemessen, nicht
behauptet. Dazu die Closure-Notiz mit Steering-Loop-Eintrag.

## 7. Vorgelagert (vor der Modus-Begründung)

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md:223-224 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

Zwei Sub-Areas berührt, beide in der Modus-Deklaration geführt:

- **`*` (Repo-Default).** Erfüllt die Schwelle über eigene Konventionsregeln
  (der gesamte Adaptions-Block), eigene Sensoren (`make gates`) und eigene
  Artefaktklassen. Der Konventionsspeicher, `AGENTS.md`, die Spec-Straten und
  die Planning-Docs fallen hierunter.
- **`tools/harness/`.** Erfüllt sie über eigene Mechanik
  ([`MR-004`](../../../../harness/conventions/MR-004-gate-nachweis-mechanik.md)),
  eigenen Sensor (`make baseline-verify`, `make baseline-probe`) und eigene
  Testfläche. Berührt, weil
  [`fetch-baseline-cache.sh`](../../../../tools/harness/fetch-baseline-cache.sh)
  die Materialisierung fährt.

Keine der beiden ist zu grob geschnitten; eine Ausdifferenzierung ist nicht
nötig, weil die Hebung beide nur an ihren bereits deklarierten Kanten berührt.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md:229-229 -->

> **Offene Beobachtungen sichten.**

Register durchgegangen (gemergter Stand, 33 Verzeichnisse). Treffer je
berührter Sub-Area:

- **`tools/harness/` (`BEO-HARN`):** ein Eintrag,
  [`check-latest-blind-before-pin`](../observations/BEO-HARN/check-latest-blind-before-pin/observation.md)
  (1×, offen) — und er ist mit diesem Slice **nicht** eingetreten:
  `make baseline-freshness` meldete alle vier neuen Tags korrekt, obwohl der
  Pin noch dahinter lag. Das ist eine Gegenbeobachtung zum einzigen Beleg,
  kein zweiter. Sie wird in §9 notiert, ohne den Zähler zu bewegen — eine
  Nicht-Wiederholung ist kein Auftreten.
- **`*` (`BEO-ALL`):** vier Einträge mit Bezug zu diesem Vorgang, alle in §5
  als Risiko geführt statt hier nur gezählt —
  [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
  (5×, über der Schwelle, ohne formgültigen Ausgang),
  [`large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)
  (2×),
  [`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
  (1×) und
  [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  (6×, gemischt). Der letzte ist beim Delta-Messen bereits **eingetreten und
  abgefangen worden**: Ein naiver Datei-Diff meldete alle 53 Bundle-Dateien
  als geändert, weil jede eine `<!-- Quelle: … blob/<tag>/… -->`-Zeile trägt;
  das echte Delta sind 13 Dateien. Wer die erste Zahl berichtet hätte, hätte
  eine fremde Menge behauptet statt der gemessenen.

Keiner der Einträge erreicht mit diesem Slice die Schwelle von 3× erstmalig;
`pin-bump-mirrors-ungated` steht bereits darüber und trägt seinen Ausgang
über [`registerzeile-ohne-ausgang-nach-schwelle`](../observations/BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle/observation.md).

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-06 gelesen. `upstream-drift.yml` **ROT**
(Lauf 05:23:44Z) — Ursache gelesen statt weggeklickt: allein die
Baseline-Currency-Achse, also genau der Anlass dieses Slice und damit eine
**planmäßige** Meldung. Zur Absicherung alle vierzehn übrigen Pin-Achsen
lokal nachgefahren (vier Versions-Achsen, Trivy, drei Action-Pins, sechs
Digest-Achsen): durchweg `ok`. `image-scan.yml` grün (07:56:19Z).

## 8. Sub-Area-Modus-Begründung

**Modus:** beide berührten Sub-Areas **GF** (Greenfield, Repo-Default) — Doc
führt, Code folgt. Kein Produkt-Code, keine Reconciliation. Die einzige Achse
mit Substanz ist das **Evidenz-/Diskrepanz-Risiko**, und es liegt nicht im
eigenen Bestand, sondern im Fremd-Delta: ob eine aktive Adaption noch trägt,
entscheidet der neue Baseline-Stand. Der Adaptions-Review in §4 (3) ist die
Antwort darauf.

## 9. Closure-Notiz (nach `done/`)

**Was hat funktioniert.** Die Delta-Messung. Ein naiver Datei-Diff meldet
alle 53 Bundle-Dateien als geändert, weil jede eine `Quelle:`-Zeile mit dem
Tag trägt; erst `diff -I '<!-- Quelle:'` zeigt die **13**, um die es geht.
Ohne diesen Schritt wäre der Slice mit einem vierfach überhöhten Delta
gestartet und hätte vier Tags für unlesbar gehalten. Ebenso getragen hat die
Zwei-Commit-Trennung beim [`MR-060`](../../../../harness/conventions.md#mr-060)-Move — inklusive der Erkenntnis, dass sie
**scheitern kann, ohne dass ein Gate es merkt**: der `pre-commit`-Hook prüft
den Arbeitsbaum, nicht den Commit-Stand.

**Was ging anders als geplant.** Der `git add` des Move-Commits scheiterte am
bereits durch `git mv` entfernten Pfad und riss die Kette mit; der Commit trug
danach nur den Rename, während der Arbeitsbaum grün war. Genau der gate-rote
Zwischenstand, den [`MR-013`](../../../../harness/conventions.md#mr-013)
ausschließt — sichtbar erst, weil der Commit-Stand eigens gegen `git show`
geprüft wurde. Über einen `--amend` geheilt (lokal, ungepusht).

Der unabhängige Review meldete **0 HIGH, 7 MEDIUM, 1 LOW, 1 INFO** und
blockierte. Bemerkenswert ist die Verteilung: Die **mechanische** Hälfte
bestätigte er in achtzehn Negativbefunden als sauber — Baum, Aliase, Zensus in
beide Richtungen, Tombstone-Entscheidung, Move-Konformität, alle zehn Gates in
seinem eigenen Lauf. Blockiert hat ausschließlich die **Selbstauskunft**
darüber: drei falsche Zahlen, eine überhobene Zeile, ein Beleg auf einer
Nebenregel, zwei Begründungen ohne Gegenstand. Alle sieben waren an Text
korrigierbar, keiner verlangte, den Bump zu wiederholen.

**Der Übergangs-Wächter hat einen zweiten Fund geliefert — nach dem Review.**
Der Plan war nach der **Baseline-Vorlage** geschrieben, nicht nach d-checks
**Haus-Form**; beim `git mv` nach `done/` meldete `verify-closure-notes`
prompt `section-missing` für `## 5. Abnahme-Punkte / Risiken`, an dem
[`MR-049`](../../../../harness/conventions.md#mr-049)s Ausgangs-Wächter hängt.
Die Datei wurde auf die neunteilige Haus-Form umgebaut und der Move
wiederholt. Zwei Dinge sind daran bemerkenswert: Der unabhängige Review hat
den Formbruch **nicht** gefunden — er prüfte gegen Plan und Kanon, nicht gegen
das Closure-Profil —, und der Wächter fand ihn genau an dem Bindepunkt, für
den [ADR-0082](../../../../docs/plan/adr/0082-uebergangswaechter-reviews-observations.md)
ihn dorthin gelegt hat. Der Fund wird deshalb **benannt, nicht registriert**:
Die Klasse ist bereits verkörpert, und sie hat funktioniert. Was fehlt, ist
kein Sensor, sondern der Blick in
[`AGENTS.md`](../../../../AGENTS.md) §1 — dort steht, dass die gelebte
Slice-Struktur eine Haus-Stil-Form ist und die Vorlage nur die Referenz-Form.

**Steering-Loop-Eintrag:** Die drei Zahl-Befunde (F-1/F-2/F-3) sind **nicht**
drei Einzelfehler. Für die **fremde** Menge — das Bundle-Delta — war eine
Abfangstelle gebaut worden, und sie griff; falsch waren genau die **selbst
erzeugten** Mengen, für die keine gebaut war. Der Reflex „fremde Zahlen sind
unsicher, eigene kenne ich" ist die Lücke, nicht die Arithmetik. Eingetragen
als weiterer Beleg bei
[`BEO-ALL/eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
(Zähler 7×) — dort ist die Teilmengen-Variante bereits als verwandte Form
benannt. **Verkörpert wurde mit diesem Slice nichts**; der Eintrag ist
gezählt, nicht verkörpert.

**Neu registriert:**
[`BEO-ALL/begruendung-traegt-entscheidung-nicht`](../observations/BEO-ALL/begruendung-traegt-entscheidung-nicht/observation.md)
(1×) — aus F-7. Die Klasse ist eigenständig: Der Tombstone traf die richtige
Entscheidung (**ein** Glob) und begründete sie mit einem Sachverhalt, der
nicht zutrifft. Weil das Ergebnis stimmt, meldet kein Sensor etwas; falsch ist
nur der Satz daneben — und der ist die einzige Fassung, aus der die nächste
Migration ihre Reichweite ableitet.

**Beobachtungs-Register (`../observations/`):** vier Belege
`evidence/slice-202.md` — bei
[`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
(7×),
[`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
(2×),
[`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(13×) und dem neuen Eintrag oben (1×). **Eine Gegenbeobachtung, bewusst ohne
Beleg:**
[`BEO-HARN/check-latest-blind-before-pin`](../observations/BEO-HARN/check-latest-blind-before-pin/observation.md)
trat **nicht** ein — `check-latest` meldete alle vier neuen Tags korrekt,
obwohl der Pin dahinter lag. Eine Nicht-Wiederholung ist kein Auftreten und
bewegt den Zähler nicht; sie steht hier, statt still zu bleiben.

**Nicht behoben, mit Begründung:** F-8 (LOW) — [`MR-055`](../../../../harness/conventions.md#mr-055)s „vier tote Aliase"
beschreibt den Bestand von slice-176 und ist für seinen Zeitpunkt korrekt;
heute sind es acht. Ein akzeptierter Adaptions-Eintrag wird nicht nachträglich
umgeschrieben (Adaptions-Block-Disziplin), und eine Mess-Aussage von damals
ist kein Defekt von heute.

**Folge-Slices:** [slice-203](../open/slice-203-v631-template-adoption.md)
(Adoption der v6.3.1-Template-Deltas) — ist eine Datei in `open/`.

**Risiken aus §5:** vier, jedes mit genau einem Ausgang — drei entfallen mit
Begründung, eines eingetreten und behoben; siehe §5.

**Verifikation.** `make gates` Exit 0 (zehn Gates, 625 Dateien, 0 Befunde) ·
`make baseline-verify` ok (54 Dateien, vollständig) · `make baseline-freshness`
meldet `v6.3.1` als neuesten Tag **und** den Inhalt unverändert — damit ist die
rote Nachtlauf-Achse geschlossen · unabhängiger Review (0/7/1/1, alle sieben
MEDIUM behoben, LOW mit Begründung offen gelassen), Report unter
[`docs/reviews/`](../../../reviews/2026-09-06-slice-202-baseline-v631-review.md).

**Drei Paarungen:** nach dem `git mv` geprüft — Ergebnis im letzten
DoD-Häkchen.

