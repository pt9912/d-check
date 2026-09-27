# Slice slice-232: `links` — Ziel hinter dem Zeilenumbruch (`](` ⏎ `ziel)`) in der gemeinsamen Extraktion

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos**. Die Closure-Bedingung geht über die eigene DoD
nicht hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:** [`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
(Beschreibung, Akzeptanzkriterien),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Determinismus). Change Request eines Konsumenten; Auftraggeber-Entscheide zum
Zuschnitt (Reichweite, Aktivierung) im Lauf vom 2026-09-27: **gemeinsame
Extraktion**, **standardmäßig an**. Eine neue ADR begleitet die Änderung
(Nummer bei der Beanspruchung: die nächste freie).

**Berührte Spec-Stellen:**
[`DC-FA-LINK-001.a`](../../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
(Schritt 3, die dort als „normative Grenze für alle Module" geführte
Zeilenbasiertheit).

**Verantwortlich:** —

**Autor:** claude-sonnet-5. **Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Ein gültiger Markdown-Link, dessen Zieladresse hinter einem
Zeilenumbruch nach `](` steht, wird von der gemeinsamen Link-Extraktion
erkannt; das Modul `links` meldet ein totes Ziel dann wie bei der Inline-Form
(`target-missing`).

**Befund, gegen den Code gelesen.** `parseLinkAt` in
`internal/hexagon/core/rules/markdown.go` arbeitet auf **einer** Zeile:
die Klammer-Suche nach `)` schließt innerhalb der Zeile nicht, der Link
entfällt. Das ist keine Schlamperei, sondern die dokumentierte Grenze in
der Spezifikation (Berührte Spec-Stellen, Schritt 3) — der Slice **ändert Vertragstext**, nicht nur Code.

**Abnahme (aus dem Change Request):** ein Fixture-Repo mit je einer Datei pro
Form und totem Ziel-Pfad ergibt **vor** der Änderung 0 Befunde und **danach**
je einen `target-missing`; dieselben Formen mit lebendem Ziel bleiben bei 0
Befunden; Inline, Titel-Form, Spitzklammer und Klammern im Ziel ändern ihr
Verhalten nicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Referenz-Definitionen** (`[name]: ziel`) — anderer Vertragssatz: sie stehen
  heute im Out-of-Scope der Link-Anforderung (siehe Bezug) und brauchen dort eine
  Streichung, keine Schärfung der Extraktion. Übernimmt
  [slice-233](../open/slice-233-links-referenz-definitionen.md).
- **Linktext über Zeilenumbruch** (`[lang\ntext](ziel)`) — der
  Change Request nennt nur die Lücke zwischen `](` und der Adresse; die
  Lücke im Linktext ist nicht gemessen. Die Grenze in der Spezifikation wird
  **verengt** und bleibt für den Linktext stehen, statt sie ohne Anlass zu
  streichen.
- **Titel hinter dem Zeilenumbruch** (`](ziel` ⏎ `"titel")`) — dieselbe
  Begründung: nicht gemeldet, nicht gemessen.
- **Bilder, Autolinks, Referenz-Links in Vollform, Modul `anchors` bei der
  neuen Form** — im Change Request ausdrücklich „nicht gemessen". Bilder
  teilen den Parser und ziehen mit; **als Zusage** wird nur der Inline-Link
  geführt, das Verhalten für Bilder wird gemessen und benannt, nicht
  behauptet.
- **`--repair` für die neue Form** — Fix-Kandidaten schreiben auf Byte-Spannen
  einer Zeile (`ExtractLinkSpans`); eine zeilenübergreifende Spanne fällt aus
  dieser Mechanik. Die neue Form wird dort **nicht** repariert, und die Grenze
  steht in der Spezifikation (anderer Vorgang: Schreib-Pfad, nicht Melde-Pfad).
- **Handbuch, README, CHANGELOG, Release** — Release-Prep-Vorgang, kein
  Feature-Commit (`AGENTS.md` §5).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — höchstens drei Liefer-Punkte.

- [ ] **Vertrag:** `spec/lastenheft.md` (Link-Anforderung: Beschreibung, zwei
      Akzeptanzkriterien Negative/Happy für die neue Form, Versions-Bump mit
      Historie-Zeile nach [MR-032](../../../../harness/conventions/MR-032-historie-vor-accepted.md)),
      `spec/spezifikation.md` (Extraktions-Abschnitt, Schritt 3: Grenze verengt,
      Behandlung von Zeilennummer und `--repair` benannt) und eine neue
      ADR samt Index-Eintrag.
- [ ] **Extraktion:** die gemeinsame Extraktion erkennt `](` + Whitespace mit
      **einem** Zeilenumbruch + Adresse; alle Konsumenten von `ExtractLinks`
      und `ExtractLinkSpans` sind gelesen, ihr Verhalten für die neue Form in
      der ADR benannt. Tests: Fixture je Form mit totem Ziel (vorher 0 /
      nachher je 1 `target-missing`, **Rot-Beleg gegen den alten Stand**) und
      mit lebendem Ziel (0), plus unveränderte Inline-/Titel-/Spitzklammer-/
      Klammer-Kontrollen.
- [ ] **Bestandsmessung:** `make doc-check` auf diesem Repo (Dogfooding) und die
      Befundzahl der Konsumenten-Module gegen den Stand vor der Änderung;
      jeder neue Befund ist behoben oder mit Grund benannt (Ausgabe in §7).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8, kein
      Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine
      Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Link-Anforderung um die Form erweitern (Erweiterung einer bestehenden Anforderung, kein neues Kürzel — Einzelmodul-Frage nach dem Schnitt-Kriterium der Historie); Bump + Historie |
| `spec/spezifikation.md` (Extraktions-Abschnitt) | update | Schritt 3: „zeilenbasiert" wird eingegrenzt; Befund-Zeile der neuen Form festgelegt |
| `docs/plan/adr/` (neu) + `README.md` | neu / update | Begründung der gemeinsamen Extraktion und der Standard-Aktivierung; Alternativen (nur `links`, Opt-in) mit dem Grund ihrer Verwerfung; `## Re-Evaluierungs-Trigger` |
| `internal/hexagon/core/rules/markdown.go` | update | Erkennung im Parser bzw. der Vorverarbeitung |
| `internal/hexagon/core/rules/markdown_test.go`, `links` -Tests, Akzeptanz-Fixture | update / neu | Happy/Boundary/Negative nach den Akzeptanzkriterien |

**Entscheidungen, die die ADR trägt (vor dem Code zu treffen):**

- **Wo die Erkennung sitzt.** `LinkSpan` ist zeilenlokal (Byte-Positionen in
  einer Zeile); `ids`, `pins`, `sources`, `repair` und
  `planning_observations` rechnen damit. Kandidaten: (a) die Vorverarbeitung
  faltet `](`-Zeile und Folgezeile, positionserhaltend; (b) ein eigener
  absatzweiser Pfad neben `forEachLink`. Gewählt wird nach dem
  Blast-Radius je Konsument, nicht nach Bequemlichkeit.
- **Welche Zeile ein Befund nennt** — die der `](`-Zeile oder die der Adresse;
  gilt für alle Module gleich.
- **Was `pins` mit dem Marker der Folgezeile tut** — der `dpin`-Marker bindet
  an den Link **derselben Zeile**; die neue Form ändert das nicht still.

## 4. Trigger

**Start** (`next` → `in-progress`): keine Abhängigkeit. Bei der Beanspruchung
entsteht der dritte Vorprüfungs-Block (Nachtlauf-Stand, `make nightly-state`).
[slice-233](../open/slice-233-links-referenz-definitionen.md) startet **nach** dem
Closure dieses Slice: beide bumpen die Link-Anforderung und das Lastenheft.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn die
  Konsumenten-Lektüre zeigt, dass mehr als zwei Module ihre Zeilenlokalität
  nicht kampflos aufgeben (dann: Vertrag + ADR als erster Slice, Umsetzung
  als zweiter), oder die Bestandsmessung mehr neue Befunde liefert, als eine
  Review-Sitzung als Nachzug trägt.
- `in-progress` → `open` (blockiert): wenn die Wahl der Befund-Zeile eine
  Auftraggeber-Entscheidung braucht, die im Lauf nicht fällt.

## 5. Closure-Trigger

DoD vollständig (§2) + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Standard-an: Konsumenten mit bisher stillen toten Zielen werden rot.
  **Ausgang:** bei Closure zu vergeben (Release-Prep nennt es im CHANGELOG).
- Ein Konsument (`ids`, `pins`, `repair`) rechnet still mit der Zeilenlokalität
  und liefert bei der neuen Form ein falsches, aber grünes Ergebnis.
  **Ausgang:** bei Closure zu vergeben.
- Der Befund-Zeilen-Entscheid verschiebt Zeilennummern für Bestandsbefunde
  der neuen Form; der Determinismus (siehe Bezug) bleibt gewahrt, die
  Byte-Identität gilt für Eingaben **ohne** die neue Form. **Ausgang:** bei
  Closure zu vergeben.

## 7. Closure-Notiz

*(Bei der Closure zu füllen.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das Produkt
(`spec/`, `internal/`) unter dem Default `*` (Kürzel `ALL`,
`harness/conventions.md` §Modus-Deklaration); bereits deklariert, keine
Ausdifferenzierung nötig.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register unter
`docs/plan/planning/observations/` durchgegangen (Stichworte Zeilenumbruch,
Link-Extraktion, Referenz-Definition, `ExtractLinks`): **keine Treffer**. Die
im Change Request genannte Beobachtung liegt im Register des Konsumenten,
nicht in diesem.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

### Sub-Area: `*` (Produkt: Spezifikation und Code)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — Lastenheft, Spezifikation und ADRs sind
  Kanon; die Grenze steht ausdrücklich im Extraktions-Abschnitt der Spezifikation.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** Mittel — die Änderung hebt eine
  dokumentierte Grenze auf, die elf Konsumenten der Extraktion implizit
  kennen (§6); Doc führt, der Code folgt.
- **Reconciliation-Aufwand:** Keiner — kein Brownfield-Bestand.
