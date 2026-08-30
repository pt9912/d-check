# Slice slice-182: Eine erklärte Teilmenge darf die Menge leeren — wenn sie sagt, wie viele

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der Anlass ist ein **eingehender CR**, keine Welle.

**Bezug:** [der eingehende CR 4](../../cr/2026-08-30-cr-a-check-leermenge.md)
(Antrag und Beleg des Absenders);
[CR 3](../../cr/2026-08-30-cr-a-check-structure-teilmenge.md) samt
[Antwort](../../cr/2026-08-30-antwort-a-check-structure-teilmenge.md) (der
Vorgänger, aus dessen Anwendung dieser Antrag kommt);
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(das erweiterte Modul);
[ADR-0075](../../adr/0075-erklaerte-teilmenge-in-structure.md) (die
Nullmengen-Härte, deren Reichweite hier zurückgeschnitten wird).

**Berührte Spec-Stellen:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
seine `.a`-Verfeinerung (Schritt 3, §2-Schema, §4-Grund-Codes) und
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
— letztere **nur zur Abgrenzung**, sie wird nicht geändert.

**Verantwortlich:** — (bis zur Priorisierung).

**Autor:** pt9912. **Datum:** 2026-08-30.

---

## 1. Ziel

**Zwei verschiedene Zustände teilen sich einen Grund-Code.** Seit
[ADR-0075](../../adr/0075-erklaerte-teilmenge-in-structure.md) meldet eine
Regel `section-missing`, wenn `exempt-section-pattern` **alle** Treffer
abzieht. Der Absender von
[CR 4](../../cr/2026-08-30-cr-a-check-leermenge.md) hat gemessen, dass das
seinen Bestand trifft: 19 Anforderungen, 19 grandfathert, 19 passende
Abschnitte — und die Regel meldet rot, wo das abgelöste Skript
`0 neue AC(s) geprueft, 19 grandfathered` mit Exit 0 meldete. **Das Modul
macht mehr rot als der Sensor, den es ablöst.**

| Zustand | Bedeutung | heute |
|---|---|---|
| `section-pattern` trifft **nichts** | Konfigurationsdefekt | `section-missing` — richtig |
| Muster trifft, Ausnahme nimmt **alle** | Bestandszustand: *„noch nichts Neues zu prüfen"* | `section-missing` — falsch |

**Die Härte bleibt richtig, ihre Reichweite ist zu weit.** [ADR-0075](../../adr/0075-erklaerte-teilmenge-in-structure.md) begründet
sie damit, dass ein zu breites Muster die Regel sonst **still** abschaltet. Das
trifft ein **generisches** Muster. Ein **aufzählendes**, das 19 Kennungen
einzeln nennt, kann nicht versehentlich zu breit werden — es kann nur
**veralten**.

## Der Antrag ist angenommen, seine Form nicht — und das ist gemessen

Der Absender beantragt `exempt-may-empty: true` und trägt die Sichtbarkeit
nach `--doctor`. Er nennt den Preis selbst *„schwächer"*.

**Gemessen ist er null.** `--doctor` läuft in **keinem** Gate dieses Repos —
nicht im `Makefile`, nicht in `.github/workflows/`, nicht im
`pre-commit`-Hook —, und `--print-mk` verteilt kein Doctor-Target an
Konsumenten. Eine Zeile dort erreicht nur den, der ohnehin nachsieht; für den
Bestandszustand, um den es geht, sieht niemand nach.

**Die Antwort ist deshalb eine Deklaration statt einer Erlaubnis.** Nicht
*„darf leer sein"*, sondern *„so viele sind es"* — `exempt-expect-count`. Damit
bleibt die Prüfung im **Gate-Lauf**, und genau das, was der Absender als
einziges Risiko benennt (*„es kann nur veralten"*), wird **laut**, statt in
einem Modus zu landen, den niemand fährt.

## 2. Vorgehen

1. **`exempt-expect-count` (int ≥ 0) an `exempt-section-pattern`.** Stimmt die
   Zahl mit den tatsächlich ausgenommenen Abschnitten überein, entsteht **kein**
   Nullmengen-Befund. Stimmt sie nicht, entsteht einer — auch dann, wenn noch
   Abschnitte übrig sind.
2. **Ein neuer Grund-Code, und das ist eine Abweichung vom Antrag.** Er sagt
   *„kein neuer Grund-Code"*; das galt seiner Form. Eine **Zahl** kann nicht
   stimmen, und dieser Zustand verlangt eine andere Reparatur (*Aufzählung oder
   Zahl nachziehen*) als `section-missing` (*Selektor korrigieren*).
   [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
   schreibt für genau diesen Fall einen eigenen Code vor. Vorschlag:
   `section-exempt-mismatch` — die Form folgt `section-heading-mismatch`.
3. **Die Drift ist beidseitig.** Mehr ausgenommen als deklariert ist ebenso ein
   Befund wie weniger. Wer die Aufzählung erweitert, ohne die Zahl zu ziehen,
   hat dieselbe Lücke wie umgekehrt — eine einseitige Prüfung wäre ein halber
   Wächter.
4. **`0` bleibt erlaubt und bedeutet etwas:** *„das Muster soll heute noch
   nichts treffen"* — ein Bestand, der erst wächst. Es zu verbieten nähme dem
   Schlüssel seinen Vorwärts-Fall.
5. **Der Schlüssel greift nur nach dem Abzug.** Trifft schon `section-pattern`
   nichts, bleibt es `section-missing` — genau die Trennung aus §1, und sie
   gehört als Test.
6. **Zwei offene Fragen werden beim Bauen entschieden und begründet**, nicht
   angenommen: ob `exempt-expect-count` in die **Regel-Identität** gehört
   (Vermutung: nein — zwei Regeln mit gleichem Selektor **und** gleicher
   Ausnahme, aber verschiedener Zahl sind ein Widerspruch, kein Paar), und ob
   der neue Befund ein **raw** Finding ist wie der Nullmengen-Befund
   (Vermutung: nein — dort hat die Regel nicht gemessen, hier hat sie eine
   Deklaration widerlegt).
7. **Config-Ränder, fail-closed:** der Schlüssel **ohne**
   `exempt-section-pattern` ⇒ Exit 2 (halbe Aktivierung, Präzedenz
   `table.order-column` ohne `table.order`); ein Wert **< 0** ⇒ Exit 2.
8. **Umkehr-Proben je Zusage** ([`BEO-023`](../observations.md), Zähler **5**):
   je Mutation genau die Tests, die dagegen stehen — und für jeden
   Regressions-Test der Beleg, dass der **Vorzustand** an seinem Fixture
   scheitert. Kein Anker einer Assertion darf still auf eine schwächere Form
   zurückfallen.
9. **ADR** (§3.6: die geprüfte Menge zu verkleinern ist eine **Lockerung**),
   Lastenheft-Bump samt Historie mit CR-Bezug, `.a`-Verfeinerung **und**
   §4-Grund-Code-Tabelle, `--print-config`-Gerüst, Handbuch.
10. **Antwort an den Absender** mit den drei Abweichungen und der
    `--doctor`-Messung, die sie trägt.
11. `make gates`, `make fullbuild`; **Review** und **Verifikation** als
    getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Schweregrad in `model.Finding`.** Der erste Entwurf des Absenders
  hätte ihn stillschweigend eingeführt;
  [`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)
  führt differenzierte Exit-Codes ausdrücklich als Out-of-Scope. Wenn das
  Werkzeug je einen bekommt, ist das ein **eigener** Entscheid.
- **Keine `--doctor`-Zeile.** Sie wäre die Form des Antrags; die gewählte
  braucht sie nicht, und eine Sichtbarkeit ohne Leser ist keine.
- **Kein `exempt-may-empty` daneben.** Zwei Schlüssel für eine Frage sind die
  Verdopplung, die [ADR-0070](../../adr/0070-tabellen-klammer-und-spaltenliste.md)
  für die Tabellen-Bedingungen zurückgebaut hat.
- **Nichts für `exempt-paths`.** Der Absender grenzt es selbst ab: ein
  Datei-Glob ist generisch und kann einen ganzen Baum verschlucken; eine
  Abschnitts-Aufzählung in **einer** Datei kann das nicht.
- **Kein Anwenden auf den eigenen Bestand.** Dieses Repo führt heute keine
  Regel mit `exempt-section-pattern`; ob eine entsteht, ist ein eigener
  Entscheid nach der Fähigkeit.

## 4. Definition of Done

- [ ] `exempt-expect-count` ist im Schema, im `--print-config`-Gerüst, in
      [`spec/lastenheft.md`](../../../../spec/lastenheft.md) (Bump + Historie
      mit CR-Bezug) und in
      [`spec/spezifikation.md`](../../../../spec/spezifikation.md) (Schritt 3,
      §2-Schema, §4-Grund-Code) geführt.
- [ ] **Default byte-identisch, gemessen:** ein Lauf ohne den Schlüssel liefert
      denselben Befundsatz wie vor der Änderung — gegen das Vorgänger-Image,
      nicht gegen einen grünen Lauf.
- [ ] **Die deklarierte Leermenge ist stumm:** N Abschnitte, N ausgenommen,
      `exempt-expect-count: N` ⇒ **kein Befund**, Exit 0. Mit Test.
- [ ] **Die Drift ist beidseitig laut:** weniger **und** mehr ausgenommen als
      deklariert ⇒ `section-exempt-mismatch`. Je ein Test.
- [ ] **Die Trennung hält:** trifft `section-pattern` nichts, bleibt es
      `section-missing` — auch mit gesetztem Schlüssel. Mit Test.
- [ ] **Zwei Config-Ränder:** ohne `exempt-section-pattern` ⇒ Exit 2; Wert < 0
      ⇒ Exit 2. Je ein Test.
- [ ] **Die zwei offenen Fragen sind entschieden und begründet** (Identität,
      raw-vs-hint) — im Code-Kommentar und in der ADR, nicht nur im Kopf.
- [ ] **Umkehr-Proben** je Zusage, jede von den Tests gefangen, die dagegen
      stehen; je Regressions-Test der Beleg, dass der **Vorzustand** an diesem
      Fixture scheitert ([`BEO-023`](../observations.md)).
- [ ] Eine ADR begründet Verortung, den **neuen Grund-Code**, die beidseitige
      Drift und die Abweichung von der beantragten Form; im
      [ADR-Index](../../adr/README.md) eingetragen.
- [ ] Das [Benutzerhandbuch](../../../user/benutzerhandbuch.md) führt den
      Schlüssel samt der Falle, die ihn nötig macht.
- [ ] Der Absender bekommt eine **Antwort** — angenommen in der Sache,
      abgelehnt in der Form, mit der `--doctor`-Messung, die das trägt.
- [ ] `make gates` und `make fullbuild` grün (Exit explizit); **unabhängiger
      Review**; **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Die geprüfte Menge zu verkleinern bleibt eine Lockerung** — jetzt eine mit
  Zahl, aber der Bestandszustand ist ab dann stumm. Wer die Zahl mitzieht, ohne
  die Aufzählung zu prüfen, hat einen Wächter, der nur noch sich selbst
  bestätigt. — **Ausgang:** *(bei Closure)*
- **Eine deklarierte Zahl ist Autoren-Text und altert wie jede.** Der Schlüssel
  verschiebt die Pflege, er nimmt sie nicht ab. — **Ausgang:** *(bei Closure)*
- **Neue Bauform ohne Präzedenz:** gemessen führt kein Schlüssel dieses Moduls
  eine **erwartete Anzahl**. Ob die Form trägt, zeigt erst der zweite Fall. —
  **Ausgang:** *(bei Closure)*
- **Die Abweichung von der beantragten Form ist ein Urteil über einen fremden
  Bestand.** Die `--doctor`-Messung gilt **diesem** Repo; ob der Adopter
  `--doctor` in einem Gate fährt, weiß ich nicht und habe ich nicht gefragt. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt
keinen Slice.

**Rückführungen:** `in-progress` → `open`, falls sich zeigt, dass der Adopter
`--doctor` sehr wohl in einem Gate fährt — dann trägt seine Form, und die
Abweichung wäre eine Bevormundung statt einer Schärfung.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/hexagon/core/` (Kern: Modell und Regel) und
  `spec/` (Anforderung und Verfeinerung). Beide fallen unter den Default `*` =
  **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration).
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-30, höchste Kennung
  `BEO-024`): [`BEO-012`](../observations.md) — **Zähler 8**, breiteste
  Instanz war slice-181 mit fünf Fundstellen derselben überdehnten Aussage:
  dieser Slice zitiert [ADR-0075](../../adr/0075-erklaerte-teilmenge-in-structure.md), CR 4 und [`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes), und **jedes** Zitat
  ist vor dem Schreiben an seinem Geltungsbereich zu prüfen.
  [`BEO-023`](../observations.md) — **Zähler 5**, zuletzt ein Test, dessen
  Anker still auf die schwächere Form zurückfiel: §2 Punkt 8 ist die Antwort
  darauf. [`BEO-013`](../observations.md) — ein Wächter, der nichts mehr
  fängt: eine Zahl, die jemand blind mitzieht, ist genau das, und sie steht als
  erstes Risiko in §5.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen:** entfällt in `open/` — der Block entsteht
  **spätestens bei der Beanspruchung** (`open→in-progress`)
  ([`MR-053`](../../../../harness/conventions.md#mr-053)).

Slice-ID: slice-182. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in).
Module: `structure`. Gates: `make gates`, `make test`, `make doc-check`,
`make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Ein optionaler Konfigurationsschlüssel, ein
neuer Grund-Code, zwei Config-Ränder; kein Fremdsystem, keine Reconciliation,
kein Bestand, der umgestellt werden müsste — dieses Repo führt heute keine
Regel mit `exempt-section-pattern`.

## 9. Closure-Notiz (nach `done/`)
