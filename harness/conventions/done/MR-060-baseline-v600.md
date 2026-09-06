# MR-060 — Baseline-Pin-Hebung auf v6.0.0 (zehnter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(keine; Pin-Fortschreibung innerhalb des von
  [`MR-023`](../../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  festgelegten self-contained Bundle-Layouts)*
- **Datum:** 2026-09-03
- **Geltungsbereich:** [§Baseline](../../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen), die
  pin-gebundenen Verweise
  ([`MR-021`](../../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
  in [`AGENTS.md`](../../../AGENTS.md), [`harness/README.md`](../../README.md), den
  aktiven `MR-*`-Dateien, den Spec-Straten, den Planning-Docs und den
  Reviewer-Skills; dazu die vier Aliase unter `.claude/rules/`
  ([`MR-055`](../../conventions.md#mr-055)).
- **Adaption:** Der Baseline-Pin ist von `v5.18.0` auf **`v6.0.0`** gehoben —
  die von [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt)
  vorgesehene Fortschreibung, zehnter Nachtrag der Serie; ersetzt
  [`MR-058`](MR-058-baseline-v5180.md) nach dessen eigenem
  Auflösungs-Trigger. Kein Layout-Wechsel: dasselbe self-contained Bundle,
  dasselbe Materialisierungs-Skript, unverändertes Pfadschema.

  **Drei Tags, Inhalt je Tag gelesen statt angenommen** — netzlos via
  `docker buildx imagetools`/direkter Release-API-Abfrage (`make
  baseline-freshness`s `check-latest` zeigte `v6.0.0` aus ungeklärtem Grund
  nicht, separat zu beobachten, kein Blocker für diese Hebung):

  | Tag | Top-Kurs-Wellen | Gegenstand |
  |---|---|---|
  | `v5.19.0` | 112–113 | wellenlose Zeitdokumente-Archivierung, Erstfassung — **fehlerhaft**, siehe unten |
  | `v5.20.0` | 114 | Korrektur der Erstfassung: echter Träger (Slice-Closure) statt bloßem Halbsatz |
  | `v6.0.0` | 115–116 | Beobachtungs-Register komplett neu gestaltet (Tabelle → Verzeichnis) |

  **Der eigene CR trägt einen Teil dieser Historie.** d-check meldete am
  2026-09-03 einen Widerspruch im Kanon (fehlender Träger für wellenlose
  Zeitdokumente-Archivierung,
  [CR](../../../docs/plan/cr/2026-09-03-cr-ai-harness-course-wellenlose-review-archivierung.md)).
  Die erste Antwort (in `v5.19.0` ausgeliefert) lehnte die beantragte
  sechste Tabellenzeile ab und nannte die Lücke stattdessen einen Halbsatz
  ohne Träger — **falsch**, wie die zweite Antwort (nach `v5.19.0`, vor
  `v6.0.0`) selbst einräumt: das „kein Trigger ohne Beobachtung ist
  Zeremonie"-Argument wurde auf die falsche Frage angewandt. `v5.20.0`
  liefert die Korrektur nach; `v6.0.0` fasst beide Deltas zusammen.

  **Das Bundle-Delta, gezählt statt geschätzt:** von **53** Dateien sind
  **39** unverändert, **14** tragen echten Inhalt: `grundlagen-begriffe.md`,
  `grundlagen-harness-dateien.md`, `grundlagen-traceability.md` (neues
  Diagramm), `modul-05-planning-harness.md`, `modul-06-roadmap.md`,
  `modul-10-review-harness.md`, `regelwerk/README.md`, `SHA256SUMS`
  (Manifest, erwartet), sowie sieben Templates
  (`AGENTS.template.md`, `harness/conventions.template.md`, <!-- d-check:ignore (Kurznamen, Wurzel ist .harness/baseline/v6.0.0/templates/) -->
  `docs/plan/planning/observation.template.md` — umbenannt aus <!-- d-check:ignore (Kurzname, Wurzel ist .harness/baseline/v6.0.0/templates/) -->
  `observations.template.md` —, `README.template.md`,
  `reconciliation.template.md`, `slice.template.md`,
  `welle-results.template.md`).

  **Zwei inhaltlich unabhängige Stränge im Delta:**

  1. **Wellenlose Zeitdokumente-Archivierung — adoptiert (Regelwerk-Ebene).**
     Modul 6 §Das Beobachtungs-Register und §Wellen-Closure-Prozedur sowie
     Modul 5 und Modul 10 bekommen eine sechste Trägerzeile: jeder
     wellenlos geschlossene Slice archiviert sich **selbst**, nach den drei
     Paarungen, Schlüssel `done/slice-<NNN>-archiv.zip` flach neben dem
     Stub. Die Umsetzung in `tools/archive-wave` (ein neuer
     Einzel-Slice-Modus) ist ein eigener, unverbindlicher Folge-Slice — kein
     Nachrüst-Zwang für den bestehenden Bestand.
  2. **Beobachtungs-Register neu gestaltet — adoptiert (slice-194/195,
     dieselbe Welle).** `BEO-<NNN>` (fortlaufende Nummer, zentrale
     Vergabestelle) wird `BEO-<KUERZEL>/<slug>` (Pfad, nachgeschlagen statt
     vergeben); die Ablage wird ein Verzeichnis
     (`observations/<pfad>/{observation.md,state.md,evidence/}`) statt
     einer Tabellen-Datei, der Zähler wird aus der Zahl der
     `evidence/`-Dateien abgeleitet statt gepflegt. Produktcode-Erweiterung
     und Datenmigration sind eigene Slices (welle-88), weil sie als ein
     Slice die eigene Größenregel (Modul 5, ≤ 3 Liefer-Punkte) überschreiten
     würden.

  **Die Spiegel-Klassen aus [`BEO-ALL/pin-bump-mirrors-ungated`](../../../docs/plan/planning/observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md):**
  40 lebende Dateien mit Pfad-/URL-Verweisen auf `v5.18.0` (30 `MR-*.md`,
  2 Reviewer-Skills, `AGENTS.md`, `harness/README.md`,
  `harness/conventions.md`, `spec/architecture.md`, `spec/spezifikation.md`,
  `docs/plan/planning/README.md`, die Roadmap, das Beobachtungs-Register)
  vor der Hebung. Alle 40 gehoben, keines stehen gelassen — kein immutables
  ADR trägt `v5.18.0` mehr **als lebender Verweis**; drei immutable ADRs
  ([ADR-0080](../../../docs/plan/adr/0080-uses-pin-tag-conflict.md),
  [ADR-0081](../../../docs/plan/adr/0081-reviews-modul.md),
  [ADR-0082](../../../docs/plan/adr/0082-uebergangswaechter-reviews-observations.md))
  zitieren den Stand ihrer Zeit weiter und
  sind quell-skopiert ausgenommen (`.d-check.yml`-Tombstone, s. u.). Zwei
  Übersetzungsfehler beim ersten Zensus-Durchlauf gefunden und behoben: der
  relative Pfad `../baseline/v5.18.0/` <!-- d-check:ignore (stale Beispiel-Pfad, illustriert den behobenen Fehler) --> in `.harness/skills/reviewer.md`
  (fünf Vorkommen) matcht nicht das Such-Muster `.harness/baseline/v5.18.0`,
  weil die Datei bereits innerhalb von `.harness/` liegt; ebenso ein
  Release-Download-Link in `AGENTS.md`, der kein `.harness/baseline/`-Pfad-
  Segment trägt.

  **Die vier Aliase unter `.claude/rules/`** sind auf `v6.0.0` gezogen;
  `make baseline-verify` bestätigt Integrität, Manifest-Deckung und
  Alias-Auflösung.

  **`.d-check.yml`-Tombstone ergänzt** für den entfernten `v5.18.0`-Baum:
  vier Glob-Einträge über die eingefrorenen Verzeichnisse
  (`docs/plan/planning/done/**`, `docs/reviews/**`, `docs/plan/adr/**`,
  `harness/conventions/done/**`) plus ein fünfter für den ausgehenden CR
  (`docs/plan/cr/**`, gesendete Kommunikation, ebenfalls nicht nachträglich
  umgeschrieben). Zwölf reale Fundstellen gemessen: vier `done/`-Slices, zwei
  Review-Reports, zwei eingefrorene `MR`-Einträge, drei immutable ADRs und
  der ausgehende CR selbst.

  **cite-Direktiven:** 15 lebende Direktiven außerhalb der eingefrorenen
  Verzeichnisse liegen in 11 Dateien; alle gegen den Datei-Diff geprüft.
  **4 sind neu verankert** (Zeilen um +4 verschoben, alle in
  `modul-05-planning-harness.md`: `:142-142`→`:146-146`,
  `:152-153`→`:156-157`), **0 entfernt**, **11 unverändert bestätigt**
  (Zieltext liegt vor der jeweiligen Einfügestelle oder in einer
  unveränderten Datei). `make doc-check --enable citations` bestätigt: 159
  Dateien, 0 Befunde.

  **Der Adaptions-Review ist durch alle 34 aktiven Einträge gelaufen
  (33 vorbestehende plus dieser): alle 34 bleiben gültig.** Kein
  Abschnittsname ist umbenannt, keine bestehende Regel entfällt oder
  widerspricht dem neuen Delta — beide Stränge sind additiv (neue
  Tabellenzeile, neue Erklärungs-Absätze, umgestellte Ablageform ohne
  Streichung einer bestehenden Zusage).

- **Begründung:** Ein Adopter, der seine Baseline nicht auf einen Tag pinnt,
  auditiert gegen ein bewegliches Ziel; der Pin macht den Stand zitierbar
  und die Abweichung benennbar. Dass er **fortgeschrieben** wird statt zu
  altern, ist die Bedingung dafür, dass der Freshness-Audit etwas zu
  vergleichen hat.
- **Löst auf:** [`MR-058`](MR-058-baseline-v5180.md)
- **Ausgelöst durch Baseline-Stand:** v6.0.0
- **Auflösungs-Trigger:** der Kurs veröffentlicht einen neuen Release-Tag;
  dann Fortschreibung durch den nächsten Nachtrag zu
  [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt).
