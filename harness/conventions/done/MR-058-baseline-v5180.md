# MR-058 — Baseline-Pin-Hebung auf v5.18.0 (neunter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(keine; Pin-Fortschreibung innerhalb des von
  [`MR-023`](../../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  festgelegten self-contained Bundle-Layouts)*
- **Datum:** 2026-09-01
- **Geltungsbereich:** [§Baseline](../../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../../conventions.md#adoptierte-konventions-quellen), die
  pin-gebundenen Verweise
  ([`MR-021`](../../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
  in [`AGENTS.md`](../../../AGENTS.md), [`harness/README.md`](../../README.md), den
  aktiven `MR-*`-Dateien, den Spec-Straten und den Planning-Docs; dazu die
  vier Aliase unter `.claude/rules/`
  ([`MR-055`](../../conventions.md#mr-055)); und zusätzlich §Die
  MR-013-Kollision aus MR-057, aufgelöst statt weiter offen gemeldet.
- **Adaption:** Der Baseline-Pin ist von `v5.15.0` auf **`v5.18.0`** gehoben —
  die von [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt)
  vorgesehene Fortschreibung, neunter Nachtrag der Serie; ersetzt
  [`MR-057`](../../conventions.md#mr-057) nach dessen eigenem
  Auflösungs-Trigger. Kein Layout-Wechsel: dasselbe self-contained Bundle,
  dasselbe Materialisierungs-Skript, unverändertes Pfadschema.

  **Drei Tags, Inhalt je Tag gelesen statt angenommen** — netzlos aus dem
  Kurs-Klon, `git show <tag>:CHANGELOG.md | head -1` je Tag:

  | Tag | Top-Kurs-Welle | Gegenstand |
  |---|---|---|
  | `v5.16.0` | 109 | die Wellen-Closure archiviert ihre Zeitdokumente |
  | `v5.17.0` | 110 | die Trennung von Umzug und Berichtigung bekommt ihre Bedingung |
  | `v5.18.0` | 111 | kein Zwang zum Nachrüsten ist kein Verbot |

  **Das Bundle-Delta, gezählt statt geschätzt:** von **51** Dateien sind
  **42** unverändert, **1** trägt nur den Versions-Stempel
  (`regelwerk/README.md`). Bleiben **8** mit echtem Inhalt:
  `grundlagen-traceability.md` (+11), `modul-05-planning-harness.md` (+6),
  `modul-06-roadmap.md` (+56/-7), `modul-08-agentenrollen.md` (+6/-5),
  `modul-10-review-harness.md` (+8), `templates/README.md` (+4/-1) und zwei
  **neue** Templates
  (`archiv-stub-slice.template.md`, `archiv-stub-welle.template.md`).

  **Der Zeitdokumente-Archiv-Mechanismus (Wellen 109 + 111) wird nicht
  adoptiert.** Ein neuer Schritt 4 der Wellen-Closure-Prozedur
  ([`modul-06-roadmap.md` §Wellen-Closure-Prozedur](../../../.harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md#wellen-closure-prozedur-modul-6))
  archiviert geschlossene Slices und Review-Reports in `done/<welle-id>/`,
  mit gekürztem Stub an der alten Stelle bei Slices. Kurs-Welle 111 macht das
  ausdrücklich **optional** — kein Nachrüst-Zwang, kein Verbot. d-check fährt
  überwiegend wellenlos; der Adoptions-Aufwand (drei Regelwerk-Module plus
  zwei neue Vorlagen) steht in keinem Verhältnis zum Nutzen im aktuellen
  Betrieb. Eigener Folge-Slice, sobald eine Wellen-Closure ansteht.

  **Die Spiegel-Klassen aus [`BEO-008`](../../../docs/plan/planning/observations.md):**
  44 lebende Vorkommen von `v5.15.0` vor der Hebung (Pfad-Verweise), dazu 3
  eingefrorene (Aussagen über die Vergangenheit in zwei offenen Plänen und
  diesem Eintrags-Vorgänger). Alle 44 gehoben, keines stehen gelassen — kein
  immutables ADR trägt `v5.15.0`.

  **Die vier Aliase unter `.claude/rules/`** sind auf `v5.18.0` gezogen;
  `make baseline-verify` bestätigt Integrität, Manifest-Deckung und
  Alias-Auflösung.

  **Ein Bulk-Zensus-Fehler aus [`MR-057`](../../conventions.md#mr-057) ist dabei
  aufgefallen und behoben:** dessen Pfad-Hebung hatte
  `harness/conventions/done/` nicht ausgenommen und zwei **eingefrorene**
  Einträge ([`MR-038`](../../conventions/done/MR-038-zitate-pin-gebunden.md),
  [`MR-041`](../../conventions/done/MR-041-guard-node-und-eigene-toolchain.md))
  versehentlich von `v5.12.0` auf `v5.15.0` gehoben — unsichtbar, solange der
  `v5.12.0`-Baum noch da war, und erst mit dessen Entfernung durch diesen
  Bump als `target-missing` sichtbar geworden. Beide auf `v5.12.0`
  zurückgesetzt, `.d-check.yml`s Tombstone-Block für die
  `v5.12.0`-Migration um `harness/conventions/done/**` ergänzt — derselbe
  Fehler wiederholt sich sonst bei jedem künftigen Bump.

  **cite-Direktiven:** von 28 lebenden Direktiven außerhalb der
  eingefrorenen Verzeichnisse liegen 8 in den fünf real geänderten Dateien;
  alle acht gegen den Datei-Diff geprüft. **6 sind neu verankert** (Zeilen um
  +5/+6 verschoben, je nach Lage relativ zu den beiden Einfügepunkten in
  `grundlagen-traceability.md` und `modul-05-planning-harness.md`), **0
  entfernt**, **2 unverändert bestätigt** (Zieltext liegt vor der jeweiligen
  Einfügestelle).

  **Der Adaptions-Review ist durch alle 33 lebenden Einträge gelaufen: alle
  33 bleiben gültig.** Kein Abschnittsname ist umbenannt (anders als beim
  `v5.15.0`-Bump), keine Regel entfällt, keine neue Ergänzung berührt eine
  bestehende Zeile.

## Die MR-013-Kollision aus MR-057, aufgelöst

[`MR-057`](../../conventions.md#mr-057) meldete eine Kollision zwischen
[`MR-013`](../../conventions.md#mr-013) und Kurs-Welle 103
([`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln](../../../.harness/baseline/v5.18.0/regelwerk/grundlagen-traceability.md#herkunfts-anker-für-steering-loop-regeln)),
statt sie aufzulösen — die neue Bedingung lag in keinem Release. Kurs-Welle
110 (jetzt in `v5.17.0`/`v5.18.0`) beantwortet genau diese Meldung, mit dem
eigenen CR als Beispiel im Kurs-Changelog zitiert: *„Ein Adopter hat eine
Kollision gemeldet statt sie aufzulösen — und beim Nachmessen war die Regel
nicht falsch, sondern unvollständig."* Die fehlende Bedingung:

<!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/grundlagen-traceability.md:113-115 -->
> **Beide Commits gehören in denselben Push.** Zwischen ihnen ist das Repo
> kurz rot; das ist zulässig, solange dieser Zwischenstand nicht die
> **Spitze** eines Push wird.

**Die naheliegende Lesart ist zu breit.** Kurs-Welle 110 endet mit *„der
Adopter kann seine Adaption auflösen, statt ihr einen Nachfolger zu geben"*
— das läse sich als Auftrag, [`MR-013`](../../conventions.md#mr-013)s Bündelung
insgesamt zu streichen. Nachgemessen am eigenen Text von
[`MR-013`](../../conventions.md#mr-013) zeigt sich: Die Kollision trifft nur
**zwei** der drei dort gebündelten Move-Klassen.

- **Slice-Lifecycle-Move und Beanspruchung.** [`MR-013`](../../conventions.md#mr-013)s
  eigene Begründung nennt den Anlass: sichtbar bei slice-040 (2026-06-21),
  eine **Push-CI**, die auf dem reinen Move-Commit rot lief
  (`target-missing` + `make planning-check`). Das ist exakt Kurs-Welle 110s
  Fall — ein Zwischen-Commit, der zur geprüften Spitze eines Push wurde. Die
  jetzt zitierbare Bedingung **bestätigt** die Praxis: beide Commits landen
  hier ohnehin nie in getrennten Pushes.
- **MR-/Wellen-Lifecycle-Move.** [`MR-013`](../../conventions.md#mr-013) trägt
  für diese dritte Klasse eine **andere, eigenständige** Begründung: eine
  nach `conventions/done/` bzw. `done/` wandernde Datei trägt relative
  Verweise, die vom neuen Ort eine Ebene tiefer auflösen müssen — ein
  byte-reiner Move-Commit wäre `doc-check`-rot. Das ist kein Risiko am
  gepushten Zwischenstand, sondern eine **lokale** Zusage:
  `make hooks`s `pre-commit`-Hook lässt seit welle-79 keinen roten
  Gate-Exit passieren, und `doc-check` (Modul `links`) meldet
  `target-missing`, sobald irgendein Dokument den (jetzt eine Ebene zu
  flachen) alten Pfad referenziert — unabhängig davon, ob und wann gepusht
  wird. Kurs-Welle 110 adressiert diesen Fall nicht; sie löst eine Bedingung
  für Push-Sichtbarkeit, nicht für lokale Commit-Zulässigkeit.

**Ausgang: Die Meldung ist aufgelöst, die Praxis bleibt unverändert.** Die
Kollision war eine **Vermischung zweier Begründungen unter einem Eintrag**,
kein Widerspruch, der eine Seite zum Weichen zwingt. [`MR-013`](../../conventions.md#mr-013)
selbst braucht keine inhaltliche Änderung — seine beiden Begründungen waren
immer schon getrennt geschrieben, nur ihre Trennschärfe gegenüber dem Kanon
war unklar, bis die fehlende Bedingung zitierbar wurde.
[`AGENTS.md`](../../../AGENTS.md) §3.3 bedarf aus demselben Grund keines
Nachzugs: seine drei Ausnahme-Absätze beschreiben bereits exakt diese
Zweiteilung, ohne sie als solche zu benennen.

- **Begründung:** Ein Adopter, der seine Baseline nicht auf einen Tag pinnt,
  auditiert gegen ein bewegliches Ziel; der Pin macht den Stand zitierbar
  und die Abweichung benennbar. Dass er **fortgeschrieben** wird statt zu
  altern, ist die Bedingung dafür, dass der Freshness-Audit etwas zu
  vergleichen hat.
- **Löst auf:** [`MR-057`](../../conventions.md#mr-057)
- **Ausgelöst durch Baseline-Stand:** v5.18.0
- **Auflösungs-Trigger:** der Kurs veröffentlicht einen neuen Release-Tag;
  dann Fortschreibung durch den nächsten Nachtrag zu
  [`MR-011`](../../conventions.md#mr-011--baseline-auf-release-tag-gepinnt).
