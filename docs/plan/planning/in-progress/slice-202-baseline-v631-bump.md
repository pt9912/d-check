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

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; Gate-Läufe und die
Closure-Pflichten darunter zählen nicht mit.

- [ ] **(1)** Vendorter Baum steht auf `v6.3.1`: alter `v6.0.0`-Baum entfernt,
      neuer materialisiert samt `SHA256SUMS`, die acht Baseline-Aliase unter
      `.claude/rules/` umgehängt; `make baseline-verify` grün (Integrität,
      Manifest-Deckung, Alias-Auflösung).
- [ ] **(2)** Alle **lebenden** pin-gebundenen Verweise auf `v6.3.1` gehoben
      (Zensus nach [`MR-021`](../../../../harness/conventions.md#mr-021),
      inklusive der beiden benannten Übersetzungsfehler-Klassen:
      `../baseline/` <!-- d-check:ignore (Muster-Praefix, illustriert die Fehlerklasse) -->-relative Pfade in `.harness/skills/` und Release-URLs
      ohne `.harness/baseline/`-Segment); `d-check:cite`-Spannen geprüft und
      wo nötig neu verankert
      ([`MR-051`](../../../../harness/conventions.md#mr-051));
      `.d-check.yml`-Tombstone für den entfernten Baum ergänzt.
- [ ] **(3)** [`MR-065`](../../../../harness/conventions.md) geschrieben —
      mit Delta-Tabelle je Tag, Adoptions-Entscheid je Strang und
      Adaptions-Review über alle aktiven Einträge —,
      [`MR-060`](../../../../harness/conventions.md#mr-060) nach
      `conventions/done/` bewegt, beide Index-Tabellen fortgeschrieben.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder eine weitere Datei in einem vorhandenen `evidence/`;
      kein Zähler wird gesetzt, er folgt aus den Dateien.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      im wellenlosen Betrieb hier geprüft, nach dem `git mv`.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/baseline/v6.3.1/` | neu | der gehobene vendorte Baum (Regelwerk + Templates + `SHA256SUMS`) |
| `.harness/baseline/v6.3.1/` | entfällt | sonst stiller Stale-Content neben dem neuen Pin ([`MR-021`](../../../../harness/conventions.md#mr-021)) |
| `.claude/rules/*.md` (Baseline-Aliase) | update | Symlink-Ziele auf den neuen Tag ([`MR-055`](../../../../harness/conventions.md#mr-055)) |
| `harness/conventions.md` | update | §Baseline, §Adoptierte Konventions-Quellen, beide Index-Tabellen |
| neue Eintrags-Datei unter `harness/conventions/` | neu | der Hebungs-Eintrag [`MR-065`](../../../../harness/conventions.md) |
| Eintrags-Datei von [`MR-060`](../../../../harness/conventions.md#mr-060) | `git mv` → `done/` | aufgelöst durch den Nachfolger |
| `AGENTS.md`, `harness/README.md`, `spec/*`, Skills, Planning-Docs | update | pin-gebundene Verweise (Menge bestimmt der Zensus) |
| `.d-check.yml` | update | Tombstone-Globs für den entfernten Baum |

## 4. Trigger

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

## 5. Closure-Trigger

Zwei beobachtbare Kriterien **und** ein Lerneintrag: (a) `make gates` grün
auf dem Endstand, `make baseline-verify` bestätigt Integrität und
Alias-Auflösung gegen `v6.3.1`, `make baseline-freshness` meldet den Pin als
aktuell; (b) kein lebender Verweis trägt mehr `v6.0.0` — gemessen, nicht
behauptet. Dazu die Closure-Notiz mit Steering-Loop-Eintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang.

- **Vier Tags auf einmal überschreiten das Review-Limit einer Sitzung**
  ([`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md),
  Zähler 2×). Das gemessene Delta ist klein (13 Dateien mit echtem Inhalt,
  davon vier mit mehr als 10 Zeilen), aber der Zensus berührt viele Dateien
  mechanisch. — **Ausgang:** <offen>
- **Der Zensus übersieht eine Spiegel-Klasse**
  ([`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md),
  Zähler 5×; kein Sensor hält die Menge). Beim `v6.0.0`-Bump wurden zwei
  Klassen erst im zweiten Durchlauf gefunden. — **Ausgang:** <offen>
- **Der mechanische Rewrite trifft eine eingefrorene Klasse**
  ([`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md),
  Zähler 1×): immutable ADRs, `done/`-Slices, Reviews und der ausgehende CR
  zitieren den Stand ihrer Zeit und dürfen **nicht** gehoben werden. — **Ausgang:** <offen>
- **Die Adoptions-Stränge werden als „mitgemacht" behauptet, obwohl nur der
  Pin steht.** Das Delta trägt mit `harness/sensors/` <!-- d-check:ignore (Kanon-Form aus v6.3.1, im Repo noch nicht angelegt) --> eine Form-Änderung, die
  d-checks überladene Sensors-Tabelle unmittelbar betrifft; sie ist in diesem
  Slice **nicht** umgesetzt — sie kann es nicht sein: die neuen Vorlagen liegen erst **nach** diesem Bump im Repo. — **Ausgang:** <offen; vorgesehener Ausgang: eingetreten, Folge-Slice slice-203 (Adoption der Template-Deltas)>

## 7. Closure-Notiz

<wird vor dem `git mv` nach `done/` gefüllt>

## 8. Sub-Area-Modus-Begründung

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

<!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md:234-236 -->

> **Keine Treffer sind ebenfalls eine Antwort** und werden notiert. Gelesen
> wird der **gemergte** Stand: Das Register ist beim Lesen so alt wie der
> letzte Merge

Register durchgegangen (gemergter Stand, 34 Verzeichnisse). Treffer je
berührter Sub-Area:

- **`tools/harness/` (`BEO-HARN`):** ein Eintrag,
  [`check-latest-blind-before-pin`](../observations/BEO-HARN/check-latest-blind-before-pin/observation.md)
  (1×, offen) — und er ist mit diesem Slice **nicht** eingetreten:
  `make baseline-freshness` meldete alle vier neuen Tags korrekt, obwohl der
  Pin noch dahinter lag. Das ist eine Gegenbeobachtung zum einzigen Beleg,
  kein zweiter. Sie wird in §7 notiert, ohne den Zähler zu bewegen — eine
  Nicht-Wiederholung ist kein Auftreten.
- **`*` (`BEO-ALL`):** vier Einträge mit Bezug zu diesem Vorgang, alle in §6
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

**Modus:** beide berührten Sub-Areas **GF** (Greenfield, Repo-Default) — Doc
führt, Code folgt. Kein Produkt-Code, keine Reconciliation. Die einzige Achse
mit Substanz ist das **Evidenz-/Diskrepanz-Risiko**, und es liegt nicht im
eigenen Bestand, sondern im Fremd-Delta: ob eine aktive Adaption noch trägt,
entscheidet der neue Baseline-Stand. Der Adaptions-Review in §2 (3) ist die
Antwort darauf.
