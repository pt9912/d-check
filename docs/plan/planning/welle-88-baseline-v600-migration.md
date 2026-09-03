# Welle welle-88: Baseline-Migration auf v6.0.0

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach in `docs/plan/planning/`; bei der Closure wandert sie nach `done/` und
bekommt ihre Ergebnisnotiz `welle-88-results.md` daneben.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-03.

---

## 1. Welle-Ziel

**Der vendorte Baseline-Baum steht auf `v5.18.0`, upstream ist bei `v6.0.0`
(Major).** Gemessen: von `v5.18.0` bis `v5.20.0` ist das Delta bereits
gelesen (drei echte Regelwerk-Dateien — `modul-05`, `modul-06`, `modul-10` —
tragen die wellenlose Zeitdokumente-Archivierung, Anlass war der eigene CR
[`2026-09-03-cr-ai-harness-course-wellenlose-review-archivierung.md`](../cr/2026-09-03-cr-ai-harness-course-wellenlose-review-archivierung.md)).
Von `v5.18.0` bis `v6.0.0` kommt ein zweites, unabhängiges Delta dazu: das
**Beobachtungs-Register wird neu gestaltet** — von einer stehenden Tabellen-Datei
(`docs/plan/planning/observations.md`, Kennung `BEO-<NNN>`, gepflegter Zähler)
zu einer Verzeichnis-Ablage (`docs/plan/planning/observations/`, Kennung = <!-- d-check:ignore (Ziel-Form, existiert erst nach slice-195) -->
Pfad `BEO-<KUERZEL>/<slug>`, abgeleiteter Zähler aus `evidence/`-Dateien).

**Zwei inhaltlich unabhängige Stränge, drei Slices — nicht zwei.** Die
Zeitdokumente-Archivierung ist ein reiner Regelwerk-Nachzug (ein Slice). Die
Register-Neugestaltung berührt drei Schichten zugleich — Produktcode (Modul
`planning`), eine Architekturentscheidung (Kürzel-Deklaration, Speicherform)
und eine Datenmigration (~27 Bestandseinträge, alle lebenden Zitate) — und
überschreitet damit als ein Slice die eigene Größenregel (Baseline-Regelwerk
`modul-05-planning-harness.md` §Ziel-Form: Slice, ≤ 3 Liefer-Punkte, ≤ 2
Schichten). Sie wird deshalb in **Fähigkeit** (slice-194) und **Migration**
(slice-195) geschnitten — schichtweise, nicht künstlich, da die Migration die
Fähigkeit voraussetzt.

**Warum das eine Welle ist, kein Bündel wellenloser Slices.** Das *Mehr*
gegenüber den drei einzelnen Slice-DoDs: erst wenn alle drei geschlossen sind,
ist der Baseline-Pin **konsistent** mit dem, was das Produkt selbst über sein
eigenes Beobachtungs-Register aussagt (Regelwerk sagt Verzeichnis-Form,
Produktcode und Bestand sagen noch Tabelle — dieser Widerspruch wäre nach
nur slice-193 real, bliebe unbemerkt ohne einen Schritt, der genau das
prüft). Der Closure-Trigger unten beobachtet diese Konsistenz, keine der
drei DoDs allein tut das.

## 2. Trigger (Welle startet)

Kein Vorwellen-Trigger nötig — der Anlass ist ein Upstream-Release
(`v6.0.0`), unabhängig vom Repo-Zustand. `in-progress/` ist leer, das
WIP-Limit ist frei.

## 3. Closure-Trigger (Welle schließt)

- slice-193, slice-194, slice-195 alle in `done/`.
- `make baseline-freshness` bestätigt Content-Match am gepinnten Tag
  (`v6.0.0`, Bytes == vendored `SHA256SUMS`).
- `make gates` und `make fullbuild` grün auf dem Endstand.
- Kein lebendes `BEO-<NNN>`-Zitat mehr im Repo (Register-Paarung hält gegen
  die neue Pfad-Form); die eingefrorenen Bestände (`done/`, `docs/reviews/`,
  `harness/conventions/done/`) bleiben unangetastet.
- Closure-Notiz `welle-88-results.md` geschrieben.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| slice-193 | Baseline-Pin auf v6.0.0 — Pfade, Cite-Spannen, wellenlose Archivierung adoptiert | — |
| slice-194 | Beobachtungs-Register: neue Architektur (ADR, Kürzel-Deklaration, Produktcode) | — |
| slice-195 | Beobachtungs-Register: Datenmigration (Bestand + lebende Zitate) | — |

## 5. Abhängigkeiten

- slice-195 wird blockiert von slice-194 (die Migration braucht die
  Produktcode-Fähigkeit, die den neuen Verzeichnis-Zustand überhaupt prüfen
  kann — sonst liefe `make gates` nach der Migration gegen ein Format, das
  kein Sensor kennt).
- slice-193 ist unabhängig von den beiden anderen, läuft aber sinnvollerweise
  zuerst: der Pin selbst muss stehen, bevor seine beiden Regelwerk-Deltas
  (Archivierung, Register-Neugestaltung) einzeln bearbeitet werden.

## 6. Out-of-Scope für diese Welle

- Jedes Release nach `v6.0.0` — eigener Anlass, eigene Welle oder eigener
  wellenloser Slice.
- Rückwirkende `BEO-<KUERZEL>/<slug>`-Neuvergabe in eingefrorenen
  Dokumenten (`done/`, `docs/reviews/`, `harness/conventions/done/`) — sie
  zitieren den Stand ihrer Zeit, keine Migration.
- Die in [`BEO-027`](observations.md) benannte Sensor-Lücke (Registerzeile
  ohne zugewiesenen Ausgang oberhalb der Schwelle) — eigener Folge-Slice bei
  Bedarf, von dieser Welle nicht miterledigt, auch wenn die neue
  `state.md`-Form sie strukturell entschärfen könnte.
- Eine Anpassung von `tools/archive-wave` an das neue Register-Format — das
  Werkzeug liest das Register heute nicht; erst wenn es das täte, wäre das
  hier relevant.

## 7. Closure-Notiz

<!-- wird erst nach Welle-Abschluss gefüllt -->

Ergebnis: `done/welle-88-results.md`
Zähler: [`observations.md`](observations.md) (Stand vor der Migration) bzw.
`observations/` (danach)
