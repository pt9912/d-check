# Slice slice-203: Sensors-Sektion auf die v6.3.1-Form — Träger, Nicht-Gates, Rollenwechsel

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`MR-065`](../../../../harness/conventions.md#mr-065) (der Bump, der
die Vorlagen brachte). Ob die Adoption eigene `MR`-Einträge braucht,
entscheidet der Slice — dieser Plan setzt es nicht voraus.

**Berührte Spec-Stellen:** — (Harness-Form, keine Spec-Stelle).

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-06.

---

## 1. Ziel

Die Sensors-Sektion in [`harness/README.md`](../../../../harness/README.md)
auf die mit `v6.3.1` adoptierte Form bringen: überladene Verträge wandern in
eigene Träger-Dateien, Nicht-Gates werden als solche **in der Zeile**
gekennzeichnet, und der Minimal Agent Workflow benennt Schritt 8 als
Rollenwechsel statt als Abschluss.

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/sensors/` <!-- d-check:ignore (wird mit diesem Slice angelegt) --> | neu | Träger-Verzeichnis; je Datei ein Target, Ziel-Form ist die vendorte `gate.template.md` |
| `harness/README.md` §Sensors | update | die drei ausgelagerten Zellen werden Index-Zeile mit Link; Nicht-Gates in eine zweite Tabelle |
| [`AGENTS.md`](../../../../AGENTS.md) §6 | update | Schritt 8 ist Rollenwechsel, kein Abschluss; kein Self-Review |

**Gemessen, nicht geschätzt** — Zeichen je Tabellenzeile, absteigend:
`verify-closure-notes` 5668 · `test` 2203 · `hooks` 2124 · `doc-check` 1908 ·
`gate-consistency` 1769 · … · `freshness-trivy`/`trivy-digest` 130. **30 von
31** Zeilen liegen über 600 Zeichen; nur die letzte passt in einen Satz.

## 3. Ausdrücklich NICHT in diesem Slice

- **Die vollständige Auslagerung aller 30 Kandidaten.** Dieser Slice nimmt die
  **drei über 2000 Zeichen** (`verify-closure-notes`, `test`, `hooks`) — den
  Ausreißer und seine zwei nächsten. Die übrigen 27 bleiben als **benannter
  Bestand** über der Grenze stehen; sie sind Folge-Slice, nicht Versehen. Grund
  ist die Größenregel: 30 Verträge umzuschichten ist in **einer**
  Review-Sitzung nicht prüfbar
  ([`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md),
  Zähler 2×) — und dieselbe Beobachtung hat genau dafür schon zweimal gefeuert.
- **Die Vollständigkeits-Zeilen mit eingefrorenen Zahlen.** Der Kanon verlangt
  dort das **Kommando** statt der Zahl; das betrifft dieselben 27 Zellen und
  wandert mit ihnen in den Folge-Slice. Die drei ausgelagerten Träger tragen
  die neue Form von Anfang an.
- **Ein `MR`-Eintrag für die Form.** Die Sensors-Datei ist Baseline-Default ab
  `v6.3.1`, keine Adaption — ein `MR` wäre nur nötig, wenn dieses Repo davon
  abwiche.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; Gate-Läufe und die
Closure-Pflichten darunter zählen nicht mit.

- [ ] **(1)** `harness/sensors/` <!-- d-check:ignore (entsteht mit diesem Slice) --> trägt die drei Verträge über 2000 Zeichen
      (`verify-closure-notes`, `test`, `hooks`) nach der vendorten
      `gate.template.md`; ihre Index-Zellen sind auf **einen Satz** gekürzt und
      verlinken die Datei. Der Link ist die einzige geprüfte Fassung der
      Zuordnung — `make doc-check` hält ihn.
- [ ] **(2)** Nicht-Gates stehen in einer **zweiten Tabelle** unter der
      Gate-Tabelle und tragen `kein Gate` **in der Bindung-Spalte**, nicht in
      Prosa daneben. Die Zuordnung ist gemessen: Kriterium ist, **worüber** ein
      Target urteilt — Zustand des Repos (Gate) gegen Vorbedingungen des
      eigenen Laufs (Werkzeug).
- [ ] **(3)** [`AGENTS.md`](../../../../AGENTS.md) §6 Schritt 8 benennt den
      Rollenwechsel (Handoff an Reviewer, kein Self-Review) mit Verweis auf den
      Reviewer-Skill.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §5 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang.

- **Die Auslagerung wird zur Umschichtung ohne Gewinn.** Der Kanon warnt
  ausdrücklich davor, das neue Verzeichnis zur Halde zu machen, die die Sektion
  vorher war: hinein gehört **Grenze und Bedienvertrag**, nicht der
  Deckungsnachweis des Werkzeugs (der lebt in ADR, Spec-Zeile, Skriptkopf).
  Prüfbar am Ergebnis: steht in einer Träger-Datei ein Satz über *welcher Test
  welche Hälfte deckt*, ist er am falschen Ort. — **Ausgang:** <offen>
- **Der Rest-Bestand wird unsichtbar.** Nach diesem Slice sind drei Zellen
  geheilt und 27 nicht — wer die Sektion liest, sieht eine halb umgestellte
  Form und hält sie womöglich für fertig. — **Ausgang:** <offen>
- **Die Nicht-Gate-Zuordnung ist ein Urteil, kein `grep`.** Ob `make bench`
  über den Zustand des Repos urteilt oder über seinen eigenen Lauf, entscheidet
  kein Muster; eine falsche Einordnung macht aus einem Gate ein Werkzeug oder
  umgekehrt — und das erste ist die gefährlichere Richtung. — **Ausgang:** <offen>

## 6. Trigger

**Start** (`open` → `in-progress`): [`MR-065`](../../../../harness/conventions.md#mr-065)
ist geschlossen und die Vorlagen liegen unter `.harness/baseline/v6.3.1/templates/`.
WIP-Limit frei (alle Lifecycle-Verzeichnisse leer).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt sich beim Schreiben, dass schon die
  **drei** Verträge nicht in einer Review-Sitzung prüfbar sind, wird die
  Auslagerung auf den einen Ausreißer reduziert und der Rest zieht mit dem
  Bestand mit.
- `in-progress` → `open` (blockiert): Erweist sich die Nicht-Gate-Zuordnung als
  Entscheid, der eine eigene Grundlage braucht (welche Targets urteilen worüber),
  ruht der Slice bis zu diesem Entscheid.

**Closure-Trigger.** Zwei beobachtbare Kriterien **und** ein Lerneintrag:
(a) `make gates` grün, und jede angelegte Träger-Datei ist aus ihrer
Index-Zeile heraus verlinkt — der Link-Sensor hält die Zuordnung, eine bloße
Namenskonvention täte das nicht; (b) keine der drei ausgelagerten Index-Zellen
ist länger als ein Satz — gemessen, nicht geschätzt.

## 7. Vorgelagert (vor der Modus-Begründung)

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md:223-224 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area berührt: `*` (Repo-Default). Sie erfüllt die Schwelle über
eigene Konventionsregeln (der Adaptions-Block), eigene Sensoren (`make gates`)
und eigene Artefaktklassen. `tools/harness/` ist **nicht** berührt — der Slice
ändert kein Skript, nur Doku-Form; die Pfad-Nähe von `harness/sensors/` <!-- d-check:ignore (entsteht mit diesem Slice) --> zu
`tools/harness/` ist Namensähnlichkeit, keine Berührung. Eine
Ausdifferenzierung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md:229-229 -->

> **Offene Beobachtungen sichten.**

Register durchgegangen (gemergter Stand, 34 Verzeichnisse). Vier Einträge
betreffen `*` und diesen Slice:

- [`large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)
  (2×, offen) — **treibt den Zuschnitt**, statt nur notiert zu werden: Der
  gemessene Bestand von 30 Kandidaten hätte den Slice über die Grenze gebracht,
  an der dieser Eintrag zweimal gefeuert hat. Der Schnitt auf drei Träger in §3
  ist die Antwort darauf. Ein dritter Beleg entsteht damit **nicht** — die
  Beobachtung wurde gelesen, bevor sie eintreten konnte.
- [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
  (6×, geplant) — unmittelbar im Gegenstand: Die auszulagernden Verträge sind
  genau die Stellen, an denen die Sektion Prüfungen beschreibt. Wer Prosa
  verschiebt, verschiebt auch ihre Behauptungen; jede übernommene Grenz-Aussage
  ist beim Umzug gegen das Target zu halten, nicht nur zu kopieren.
- [`begruendung-traegt-entscheidung-nicht`](../observations/BEO-ALL/begruendung-traegt-entscheidung-nicht/observation.md)
  (1×, offen, aus slice-202) — dieselbe Gefahr eine Ebene tiefer: Eine
  Träger-Datei besteht überwiegend aus Begründungen, und keine davon prüft ein
  Sensor.
- [`module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
  (1×, offen) — die Träger-Dateien tragen je einen Abschnitt *Grenze*; genau
  die Frage, die dieser Eintrag stellt.

Keiner der vier erreicht mit diesem Slice die Schwelle von 3× erstmalig.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-06 gelesen. `image-scan.yml` **grün**
(07:56:19Z). `upstream-drift.yml` **ROT** (05:23:44Z) — und hier greift die
benannte Grenze des Targets: es liest den **jüngsten** Lauf, nicht sein
**Alter**. Der Lauf datiert **vor** der Pin-Hebung aus slice-202; seine
Ursache — der Baseline-Rückstand — ist seither behoben, `make baseline-freshness`
meldet lokal beide Teile grün. Kein offener Befund, aber auch kein grüner
Nachtlauf, bis der nächste läuft.

## 8. Sub-Area-Modus-Begründung

**Modus:** die eine berührte Sub-Area `*` ist **GF** (Greenfield,
Repo-Default) — Doc führt, Code folgt. Kein Produkt-Code, keine
Reconciliation. **Konventionen-Dichte** hoch: die Form ist vollständig als
vendorte Vorlage verankert. **Phase-Reife** hoch: die Sektion existiert, ist
gewachsen und wird nur umgeformt. **Evidenz-/Diskrepanz-Risiko** liegt nicht
im Code-vs-Doku-Abstand, sondern in der **Aussagen-Treue beim Umzug** — die
vier gesichteten Beobachtungen oben benennen genau das.

## 9. Closure-Notiz (nach `done/`)

<wird vor dem `git mv` nach `done/` gefüllt>
