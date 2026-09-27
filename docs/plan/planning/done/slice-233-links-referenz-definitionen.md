# Slice slice-233: `links` — Referenz-Definitionen (`[label]: ziel`) mit Datei-Ziel werden geprüft

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos**. Die Closure-Bedingung geht über die eigene DoD
nicht hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:** [`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
(dessen Out-of-Scope nennt Reference-Style-Links),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Determinismus). Change Request eines Konsumenten; Auftraggeber-Entscheid im
Lauf vom 2026-09-27: **jede Definition mit Datei-Ziel**, **standardmäßig an**.
Eine neue ADR begleitet die Änderung (die nächste freie Nummer).

**Berührte Spec-Stellen:**
[`DC-FA-LINK-001.a`](../../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
(Ziel-Menge des Moduls `links`).

**Verantwortlich:** pt9912.

**Autor:** claude-sonnet-5. **Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Das Modul `links` prüft die Zieladresse jeder
Referenz-Definition `[label]: ziel` mit Datei-Ziel und meldet ein totes Ziel
als `target-missing` — unabhängig davon, ob die Definition benutzt wird.

**Befund, gegen den Kanon gelesen.** Der Out-of-Scope-Satz in
Die Link-Anforderung (siehe Bezug) schließt Reference-Style-Links ausdrücklich aus. Der Slice
**streicht** diese Ausnahme für die Definition und ändert damit Vertragstext.

**Abnahme (aus dem Change Request):** ein Fixture-Repo mit totem Ziel ergibt
**vor** der Änderung 0 Befunde und **danach** einen `target-missing`; dasselbe
mit lebendem Ziel bleibt bei 0; die bereits gemeldeten Formen ändern ihr
Verhalten nicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Auflösung der Verwendung** (`[text][label]`, `[label][]`, `[label]`) — die
  Wahl „jede Definition" macht sie für die Existenz-Prüfung überflüssig; der
  Change Request hat sie nicht gemessen. Die Definition wird geprüft, die
  Verwendung nicht aufgelöst, und die Spezifikation sagt das ausdrücklich
  (Out-of-Scope bleibt für die Verwendung stehen).
- **Modul `anchors` bei Definitionen** (`[label]: datei.md#anker`) — im
  Change Request „nicht gemessen"; `links` ignoriert wie bei Inline-Links
  den Fragment-Teil, `anchors` bleibt für Definitionen außen vor und die
  Grenze steht in der Spezifikation.
- **Definition mit Ziel hinter dem Zeilenumbruch** (`[label]:` ⏎ `ziel`) und
  **mehrzeilige Titel** — nicht gemeldet, nicht gemessen; die
  Zeilenumbruch-Form bei Inline-Links liefert
  [slice-232](../done/slice-232-links-ziel-hinter-zeilenumbruch.md).
- **`--repair` für Definitionen** — Schreib-Pfad, anderer Vorgang (siehe
  slice-232).
- **Handbuch, README, CHANGELOG, Release** — Release-Prep-Vorgang
  (`AGENTS.md` §5).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — höchstens drei Liefer-Punkte.

- [x] **Vertrag:** `spec/lastenheft.md` (Link-Anforderung: Out-of-Scope-Satz
      auf die Verwendung eingegrenzt, Beschreibung und je ein Akzeptanzkriterium
      Negative/Happy, Bump mit Historie-Zeile nach
      [MR-032](../../../../harness/conventions/MR-032-historie-vor-accepted.md)),
      `spec/spezifikation.md` (Extraktion der Definition, Form der
      Definitions-Zeile, Grenzen aus §1) und eine neue ADR samt Index-Eintrag.
      **Form nach dem Kanon** (`modul-03-spec.md`, `modul-04-adrs.md`): die
      Historie-Zeile nennt **weder ADR noch Slice** (Decken-Regel), die
      Spezifikation nennt in keinem Abschnitt eine ADR oder einen Slice; die
      ADR trägt `Schärft:` aufwärts, mindestens drei verglichene Alternativen
      mit Trade-off, eine Fitness Function und einen `Re-Evaluierungs-Trigger`.
- [x] **Prüfung:** `links` erkennt Definitions-Zeilen außerhalb von Fences und
      Inline-Code (dieselbe Vorverarbeitung wie bei Inline-Links) und prüft
      das Datei-Ziel mit **derselben** Auflösung, Escape-Prüfung, Symlink-Regel
      und `ignore-refs`-Ventil. Tests: totes Ziel (vorher 0 / nachher 1
      `target-missing`, Rot-Beleg gegen den alten Stand), lebendes Ziel (0),
      externes Schema (0), Definition im Fence (0), unveränderte
      Inline-Kontrollen.
- [x] **Bestandsmessung:** `make doc-check` auf diesem Repo (Dogfooding) gegen
      den Stand vor der Änderung; jeder neue Befund behoben oder mit Grund
      benannt (Ausgabe in §7).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine
      Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Erweiterung der Link-Anforderung (kein neues Kürzel); Bump + Historie |
| `spec/spezifikation.md` (Extraktions-Abschnitt) | update | Definitions-Extraktion und Grenzen |
| `docs/plan/adr/` (neu) + `README.md` | neu / update | „jede Definition" statt „nur verwendete" mit Alternativen; Standard-Aktivierung; `## Re-Evaluierungs-Trigger` |
| `internal/hexagon/core/rules/links.go`, `markdown.go` | update | Definitions-Erkennung, Einspeisung in die bestehende Ziel-Prüfung |
| Tests und Akzeptanz-Fixture | update / neu | Happy/Boundary/Negative nach den Akzeptanzkriterien |

**Reihenfolge der Vertrags-Änderung:** wie im Vorgänger-Plan — Lastenheft auf
`Draft`, die Änderung darf im Slice liegen (Kanon, *Wann die CR-Pflicht
beginnt*; [MR-032](../../../../harness/conventions/MR-032-historie-vor-accepted.md)
verlangt Bump und Historie-Zeile trotzdem) und steht in einem eigenen Commit vor
dem Code.

**Zu klären, bevor Code entsteht:** Wo der Definitions-Befund seine Zeile und
sein `Ziel`-Feld her nimmt (die Definitions-Zeile), und ob `codepaths`,
`external`, `tracked`, `matrix` eine Definition **mitsehen** — sie teilen
`ExtractLinks`. Liefert die Definition `LinkRef`-Werte in dieselbe Liste, ziehen
alle Konsumenten mit; das ist wie bei slice-232 in der ADR zu entscheiden,
nicht als Nebenwirkung hinzunehmen.

## 4. Trigger

**Start** (`next` → `in-progress`): **nach** dem Closure von
[slice-232](../done/slice-232-links-ziel-hinter-zeilenumbruch.md) — beide bumpen das
Lastenheft und die Link-Anforderung; parallele Läufe kollidierten in der
Versions-Nummer und in der Extraktions-Stelle. Bei der Beanspruchung entsteht
der dritte Vorprüfungs-Block (Nachtlauf-Stand).

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): wenn die Konsumenten-Frage (§3) für mehr als
  zwei Module eine eigene Entscheidung verlangt, oder die Bestandsmessung mehr
  Befunde liefert, als ein Nachzug in einer Review-Sitzung trägt.
- `in-progress` → `open` (blockiert): wenn slice-232 mit einer geänderten
  Extraktions-Schnittstelle schließt, die den Plan hier ungültig macht.

## 5. Closure-Trigger

DoD vollständig (§2) + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Standard-an: Konsumenten mit bisher stillen toten Definitions-Zielen werden
  rot — auch bei **unbenutzten** Definitionen. **Ausgang: entfallen.** Die
  Bestandsmessung (§7) zeigt 0 neue Befunde in diesem Repo. Die
  CHANGELOG-Erwähnung bleibt Standard-Release-Prep für jedes neue
  Standard-an-Verhalten (`AGENTS.md` §5), kein Folge-Slice nötig.
- Eine Definition in einem Dokument, das Ziele bewusst als Platzhalter führt
  (Templates), meldet Fehlalarme. **Ausgang: entfallen.** Kein solcher Fall
  in diesem Repo (Bestandsmessung 0 neue Befunde); `ignore-refs` bleibt als
  Ventil verfügbar, falls ein Konsument ihn braucht.
- Die Form „Definition" ist enger oder weiter als CommonMark (Einrückung bis
  drei Spaces, Label mit Backslash-Escapes, Definition in Blockquote/Liste).
  **Ausgang: eingetreten, teilweise.** Die Form wurde vor der Messung
  in der begleitenden Entscheidungs-Kette
  ([ADR-0095](../../adr/0095-links-referenz-definitionen-titel-delimiter-pflicht.md))
  ausgeschrieben und deckt die drei genannten Fälle ab. Zwei unabhängige
  Review-Runden fanden zusätzlich: (a) eine Titel-Delimiter-Lücke ließ
  gewöhnliche Prosa als Definition mit erfundenem Ziel lesen, behoben; (b)
  ein whitespace-freies, klammerartiges Ziel-Token (`[TODO]: (spaeter)`)
  bleibt unvalidiert, nicht blockierend laut Review-Verdikt. **Ausgang für
  (b): weiter offen** →
  [`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md).

## 7. Closure-Notiz

**Geliefert:** die gemeinsame Link-Extraktion (`ExtractLinks`) erkennt eine
Link-Referenz-Definition (`[label]: ziel "titel"`) unabhängig davon, ob sie
im Dokument verwendet wird; ein totes Dateiziel meldet `target-missing` wie
bei einem Inline-Link, auf der Definitions-Zeile. `spec/lastenheft.md`
(0.93.2), `spec/spezifikation.md`
([`DC-FA-LINK-001.a`](../../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
Schritt 3) und eine
dreigliedrige Entscheidungs-Kette (zwei Nachzüge, je Accepted → Superseded)
mit
[ADR-0095](../../adr/0095-links-referenz-definitionen-titel-delimiter-pflicht.md)
als aktuellem, Accepted Tip tragen die Entscheidung. Fünf der sechs `ExtractLinks`-Konsumenten
behandeln eine Definition wie jeden anderen `LinkRef`; `anchors` überspringt
sie vollständig (neuer Out-of-Scope-Satz).

**Zwei unabhängige Review-Runden fanden je einen echten HIGH-Befund.** R1
(R1-H1): die Erstfassung ließ `NormalizeTarget` naiv am ersten Leerzeichen
abschneiden — eine gewöhnliche Prosazeile (`[TERM]: First In, First Out`)
wurde dadurch als Definition mit erfundenem Ziel „First" gelesen.
Korrigiert durch die aktuelle Fassung der Entscheidungs-Kette
([ADR-0095](../../adr/0095-links-referenz-definitionen-titel-delimiter-pflicht.md)):
der Rest der Zeile hinter dem Ziel-Token muss leer oder ein korrekt
delimitierter Titel (`"…"`, `'…'`, `(…)`) sein, sonst bleibt die **ganze**
Zeile unerkannt. R2 fand dabei außerdem verbotene Review-Befund-Marker in
zwei neuen Kommentaren (R2-H1, trivial behoben) sowie eine eigene, während
der Implementierung **selbst** vor dem ersten Test entdeckte Ungenauigkeit
in der Erstfassung der Kette: ein Backslash-Escape im Label lässt die ganze
Zeile unerkannt, nicht nur die Label-Grenze verschieben — noch vor jedem
Review korrigiert. Alle merge-blockierenden Befunde sind eingearbeitet; ein
MEDIUM-Fund (R2-M1: ein whitespace-freies, klammerartiges Ziel-Token wie
`[TODO]: (spaeter)` bleibt
unvalidiert) ist laut Review-Verdikt nicht blockierend und als offene
Beobachtung registriert (§6).

**Bestandsmessung (DoD-Pflicht, alle sechs `ExtractLinks`-Konsumenten):**
`make doc-check --enable external --enable tracked` gegen den Stand
unmittelbar vor diesem Slice (`6f8d6906`) und gegen die Endfassung: **107 →
107 Befunde**, unverändert (`external-status`, Sandbox ohne Netz). Kein
neuer Befund, unabhängig von R1 und R2 bestätigt.

**Steering-Loop-Eintrag:** neue Beobachtung
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(1×, `offen`): eine Erkennungs-Regex für eine Markdown-Mikrosyntax bleibt
gegen Freitext-Negativfälle ungeprüft, bis ein Review adversarial dagegen
testet — in diesem Slice zweimal in Folge (R1-H1, R2-M1), derselbe
Mechanismus auf zwei Ebenen.

**Risiko-Ausgänge:** zwei von drei Risiken aus §6 *entfallen*, eines
*eingetreten (teilweise)* mit einem Teil-Ausgang *weiter offen* — siehe
dort.

**Register-Sichtung bei Planung:** §8 vermerkte „keine Treffer" —
zutreffend geblieben; die neue Beobachtung war zuvor nicht im Register.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das Produkt
(`spec/`, `internal/`) unter dem Default `*` (Kürzel `ALL`); bereits
deklariert, keine Ausdifferenzierung nötig.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Stichworte Referenz-Definition, Reference-Style, `links`-Modul): **keine
Treffer**. Die im Change Request genannte Beobachtung liegt im Register des
Konsumenten.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung, 2026-09-27):
`make nightly-state` meldet `upstream-drift.yml` **ROT** (Lauf
2026-09-27T06:03:41Z, unverändert seit slice-232), `image-scan.yml` **grün**
(Lauf 2026-09-27T09:18:12Z). Derselbe veraltete Stand wie bei slice-232: vier
der fünf gemeldeten Fremd-Release-Stände sind bereits gehoben, die
Kurs-Baseline (`v6.9.0` → `v6.10.0`) bleibt bewusst zurückgestellt bis
`v6.11.0`. Keiner der fünf Punkte berührt dieses Modul (`links`).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

### Sub-Area: `*` (Produkt: Spezifikation und Code)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — der Out-of-Scope-Satz ist Kanon; seine
  Änderung läuft über Lastenheft-Bump und ADR.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** Mittel — Fehlalarme bei Platzhalter-Zielen
  sind das erwartbare Bestandsrisiko (§6); die Bestandsmessung trägt es.
- **Reconciliation-Aufwand:** Keiner — kein Brownfield-Bestand.
