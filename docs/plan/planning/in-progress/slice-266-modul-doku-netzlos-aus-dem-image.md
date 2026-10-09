# slice-266: Die Doku der Module ist netzlos aus dem Image lesbar

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** eingehender CR von `sf-connector`
([`docs/plan/cr/2026-10-09-cr-eingehend-sf-connector-reviews-zusage.md`](../../cr/2026-10-09-cr-eingehend-sf-connector-reviews-zusage.md),
Punkt 4): ein netzlos arbeitendes Repo kann die Erkennungsregel eines Moduls
heute nur durch Probieren ermitteln.

**Berührte Spec-Stellen:** eine neue CLI-Anforderung im Lastenheft (Handbuch
und Spezifikation aus dem Binary) samt Verfeinerung in der Spezifikation;
die Kennung vergibt der Feat-Commit.

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Wer nur das Image hat, liest aus ihm, was ein Modul prüft und welche
Schlüssel es nimmt — über eine Ausgabe des Werkzeugs (etwa `--help <modul>`)
oder eine mitgelieferte Datei. Beim Beanspruchen wird entschieden, welche Form
und welche Quelle; die Quelle darf nicht zu einer zweiten Beschreibung neben
Spezifikation und Handbuch werden, die driftet.

**Entschieden beim Beanspruchen (Auftraggeber):** Handbuch und Spezifikation
werden beim Build ins Binary eingebettet, byte-gleich zu ihren Quelldateien;
`--manual <begriff>` gibt die Abschnitte beider Dokumente aus, deren
Überschrift den Begriff nennt, samt Unterabschnitten. Eine zweite Beschreibung
entsteht nicht — die Quelle **ist** die Datei. Der Titel eines Dokuments
liefert es ganz.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Erkennungsregel von `reviews` selbst** — slice-265.
- **Eine Übersetzung oder Kürzung des Handbuchs** — der Gegenstand ist die
  Erreichbarkeit, nicht ein neuer Text.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] Die Doku der Module ist aus dem Image ohne Netz lesbar; die Form ist im
      Lastenheft zugesagt und durch einen Test gehalten.
- [x] Ein Sensor oder Test hält die mitgelieferte Doku gegen ihre Quelle, damit
      sie nicht driftet; `make gates` grün.
- [x] Folge-ADR zu [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md) für das Paket im Modul-Root, bis zur Closure
      `Proposed`; `.a-check.yml` gibt ihm eine Schicht.
- [x] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [x] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft; der CR
      trägt seine Entscheidung zu Punkt 4.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Paket im Modul-Root (Einbettung), CLI (`--manual`), Tests | create/update | Ausgabe aus den eingebetteten Dokumenten |
| `tools/image-test.sh` | update | eine Phase: `--manual` netzlos im Container |
| `spec/lastenheft.md`, `spec/spezifikation.md` | update | Zusage |
| `.a-check.yml`, Folge-ADR zu [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md), ADR-Index | create/update | **Plan-Änderung nach Review R1 F-1:** das Paket im Modul-Root bekommt eine Schicht und eine Entscheidung, sonst bleibt es ungeprüft |
| Hilfe-Ausgabe, `harness/sensors/image-test.md` | update | **Plan-Änderung nach Review R1 F-3/F-4:** die Spiegel der neuen Option |

## 4. Trigger

**Start** (`next` → `in-progress`): `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): die gewählte Quelle verlangt einen Umbau
  des Handbuchs — dann erst die Quelle, dann die Ausgabe.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Produkt-Verhalten — geht mit dem nächsten Release
hinaus.

## 6. Risiken und offene Punkte

- **Zweite Beschreibung** — eine eigene Modul-Doku im Binary driftete gegen
  Spezifikation und Handbuch. — **Ausgang:** entfallen — die Dokumente sind
  eingebettet, nicht nachgeschrieben; die Ausgabe ist die Quelle (Verifikation:
  Abschnitte byte-gleich).

## 7. Closure-Notiz

- **Was hat funktioniert:** `d-check --manual <begriff>` gibt die Abschnitte
  von Benutzerhandbuch und Spezifikation aus, deren Überschrift den Begriff
  nennt — ohne Netz und ohne Mount, aus dem Binary selbst. Die Dokumente sind
  per `go:embed` eingebettet; eine zweite Beschreibung gibt es nicht. Die
  Abschnitts-Suche nutzt die Überschriften-Lexik von `structure`. Die
  Verifikation maß im Image byte-gleiche Abschnitte und alle Lastenheft-Fälle;
  Phase 5 des Image-Tests wird gegen v0.85.0 aus dem richtigen Grund rot.
- **Was ging anders als geplant:** `go:embed` erreicht nur Dateien unter dem
  eigenen Verzeichnis — das Paket musste in den Modul-Root, außerhalb des
  Layouts von [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md). R1 fand
  es ungeprüft; daraus wurden eine Schicht in `.a-check.yml` und
  [ADR-0107](../../adr/0107-mitgelieferte-dokumente-im-modul-root.md), als
  Plan-Änderung vor dem Code. Die Gegenprobe zeigte, dass a-check den
  Importpfad des Modul-Roots nicht auflöst; die Kante und den Inhalt des
  Pakets halten jetzt zwei Tests. Ein abgebrochenes `git add` ließ eine
  Commit-Botschaft stehen: `d04d5c8a` trägt die Löschung von
  `manual_test.go` unter der Botschaft der Plan-Änderung (gepusht, im
  Folge-Commit benannt) — Schritt 21 war nicht ausgeführt.
- **Steering-Loop-Eintrag:** keiner mit neuer Schwelle.
  `BEO-ALL/path-scoped-commit-carries-staged-rest` ist als Schritt 21 im
  Workflow-Skelett verkörpert und trat erneut auf — die Regel stand, sie wurde
  in einer `&&`-Kette übergangen.
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-266.md` in
  [`BEO-ALL/path-scoped-commit-carries-staged-rest`](../observations/BEO-ALL/path-scoped-commit-carries-staged-rest/state.md);
  neu
  [`BEO-ALL/a-check-sieht-modul-root-import-nicht`](../observations/BEO-ALL/a-check-sieht-modul-root-import-nicht/state.md)
  (1×).
- **Folge-Slices:** keiner. Produkt-Verhalten — geht mit dem nächsten Release
  hinaus; Handbuch und README nennen `--manual` in der Release-Prep. Der CR von
  `sf-connector` trägt die Entscheidung zu Punkt 4.
- **Risiken aus §6:** entfallen (siehe §6). Trigger-Audit: kein Carveout, kein
  bootstrap-aware Gate; [ADR-0107](../../adr/0107-mitgelieferte-dokumente-im-modul-root.md)
  neu und `Accepted`, ihre Trigger nicht eingetreten; keine Hard Rule mit
  eingetretenem Trigger. Nachtlauf-Stand
  ([`MR-053`](../../../../harness/conventions.md#mr-053)): wie in §8.
- **Drei Paarungen:** (a) Anker — kein Eintrag mit `liegt in`; (b)
  Folge-Slices — keiner genannt; (c) Register — die zwei zitierten
  Beobachtungen existieren und tragen Belege.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind die CLI, ein neues Paket
für die Einbettung, das Image-Testskript und die beiden Spec-Straten — alle
unter dem Default `*` (`ALL`); keine eigene Konvention, kein eigener Modus.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Keine offene Beobachtung
betrifft die Erreichbarkeit der Doku. Mittelbar berührt ist
[`BEO-ALL/semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md):
eine neue Option hat Spiegel (Hilfe-Ausgabe, Lastenheft, Spezifikation,
Handbuch in der Release-Prep), die vor dem Editieren aufgelistet werden.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`upstream-drift` grün (2026-10-09T07:17Z). `image-scan` rot (2026-10-09T10:37Z)
— der Lauf liegt vor dem Release v0.85.0, das die gemeldeten CVEs behebt.
Der Slice ändert das Image (größeres Binary), nicht seine Basis.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
