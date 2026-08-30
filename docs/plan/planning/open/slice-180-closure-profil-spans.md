# Slice slice-180: Der Closure-Bindepunkt sieht die Defekte, die er heute übersieht

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der Anlass ist eine Messung, keine Welle.

**Bezug:**
[`DC-FA-SPAN-001`](../../../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
(das Modul, das an den Bindepunkt kommt);
[`DC-FA-CLI-012`](../../../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)
(der Pfad, über den der Bindepunkt sein Profil fährt);
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) (die
Entscheidung, die den Bindepunkt und sein zweites Profil geschaffen hat);
[ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) (deren offener
Punkt (a) diese Messung beantwortet);
[slice-178](../open/slice-178-offene-tasks-roh.md) (wartet auf dieses Ergebnis).

**Berührte Spec-Stellen:** —. Der Slice ändert **keine** Zusage des Produkts;
er ändert, welche Module dieses Repo an seinem eigenen Closure-Bindepunkt
fährt.

**Verantwortlich:** — (bis zur Priorisierung).

**Autor:** pt9912. **Datum:** 2026-08-30.

---

## 1. Ziel

**Der Closure-Bindepunkt kann zwei Defekte nicht sehen, die das Produkt längst
findet.** `make verify-closure-notes` fährt `planning` und `structure` über das
`done/`-Profil. Beide lesen den **bereinigten** Abschnitts-Text. Ein
**vergessener Schluss-Fence** verschluckt damit alles dahinter — und die
Bedingung, die darüber etwas zusagt, wird **still** wahr.

**Gemessen an vier Proben, mit dem heutigen Stand:**

| Probe | `structure` allein | `structure` + `spans` |
|---|---|---|
| offener Haken, nackt | `section-forbidden` | ebenso |
| ungerader Backtick im Absatz | `section-forbidden` | ebenso **plus** `span-unclosed` |
| Span, der das Item **umschließt** | **still** | **still** |
| vergessener Schluss-Fence | **still** | `fence-unclosed` |

**Die Fence-Lexik ist nicht der Gegenstand, und das ist gemessen.** Das Produkt
führt zwei Lesarten — den Toggle (`FenceToggle`, benutzt von `proseLines`,
`PreprocessMarkdown` und der Abschnitts-Findung) und die CommonMark-Lesart
(`FenceRun` + `FenceCloses`, benutzt von `spans` und dem Tabellen-Reader).
Beide über **alle 676** Markdown-Dateien dieses Repos gegeneinander gefahren:
**null Abweichungen**. Die Messung *kann* finden — drei konstruierte Fälle
(Infostring hinter dem Schluss-Fence, Zeichenwechsel `` ``` ``→`~~~`, zu kurzer
Schluss-Run) meldet sie, die saubere Datei nicht.
[ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) hat genau diesen
Punkt als *„bewusst offen gelassen … unbelegt — kein Realfall in den 522
Dateien"* geführt; bei 676 gilt es weiter.

**Die dritte Zeile der Tabelle ist kein Defekt, und deshalb steht sie hier
statt in der DoD.** Ein wohlgeformter Code-Span, der ein Task-Item umschließt,
macht das Item nach CommonMark zu **Code**. Die Bereinigung hat recht; nur der
Autor meinte etwas anderes. Kein Modul kann das unterscheiden, weil am Dokument
nichts falsch ist.

## 2. Vorgehen

1. **`spans` an den Closure-Bindepunkt.** Das Makefile-Rezept von
   `verify-closure-notes` bekommt `--enable spans`. **Nicht** in `modules:` des
   Profils: dort steht `modules: []`, und die Aussage *„das Profil aktiviert
   ausschließlich, was die Kommandozeile dazuschaltet"* soll wahr bleiben.
2. **Der Bestand ist vorher gemessen**, nicht danach entdeckt: `spans` meldet
   über das ganze Repo heute **0** Befunde. Der Zuwachs ist damit
   rausch-frei — und das gehört als Zahl in den Commit, nicht als Zusage.
3. **Vier Deklarations-Flächen mitziehen** —
   [`AGENTS.md`](../../../../AGENTS.md) §4 (die Zeile zu
   `verify-closure-notes`), die Sensors-Tabelle in
   [`harness/README.md`](../../../../harness/README.md), der `##`-Hilfetext des
   Targets (den `gate-consistency` gegen die Doku hält) und der Kopfkommentar
   von [`.d-check.closure.yml`](../../../../.d-check.closure.yml).
4. **Eine ADR trägt den Entscheid und die Messung.** Warum `spans` an diesen
   Bindepunkt gehört, warum es **nicht** in `modules:` steht, und warum die
   Fence-Lexik unangetastet bleibt.
5. **[ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) bekommt einen
   `## Geschichte`-Eintrag**: sein offener Punkt (a) ist erneut gemessen — 676
   statt 522 Dateien, weiterhin kein Realfall. Die ADR ist `Accepted` und wird
   nicht im Kern angefasst.
6. **[slice-178](../open/slice-178-offene-tasks-roh.md) §1 wird richtiggestellt.** Seine
   Expositions-Behauptung zählt Backticks **abschnittsweise** (*„slice-061 und
   slice-076 tragen ungerade Backtick-Zahlen (25 bzw. 45)"*); die Paarung läuft
   **absatzweise**. Gemessen trägt keine der beiden Dateien einen unbalancierten
   Absatz — `spans` meldet dort nichts. Das ist
   [`BEO-020`](../observations.md), und die Zeile wird korrigiert, nicht
   gestrichen.
7. `make gates`, `make fullbuild`; **Review** und **Verifikation** als getrennte
   Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Änderung an der Fence-Lexik.** Beide Lesarten bleiben, wie sie sind —
  die Messung hat den Anlass nicht gefunden, den
  [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) als Bedingung
  benannt hat.
- **Keine Roh-Zeilen-Bedingung.** Das ist
  [slice-178](../open/slice-178-offene-tasks-roh.md), und er kommt danach.
- **Kein `non-empty` an den bestehenden Closure-Regeln.** Es machte denselben
  Fall laut, wäre aber eine Bedingungs-Änderung an einer Zusage, die dieser
  Slice nicht führt.
- **Keine Antwort auf den umschließenden Span.** Er ist kein Defekt (§1).

## 4. Definition of Done

- [ ] `make verify-closure-notes` fährt `spans`; **gemessen vorher und
      nachher**: Befundzahl über den Bestand unverändert (heute 0 aus `spans`).
- [ ] **Die neue Deckung ist belegt, nicht behauptet:** eine Probe mit
      vergessenem Schluss-Fence im `done/`-Abschnitt läuft gegen das Profil
      **vorher grün und nachher rot** (`fence-unclosed`).
- [ ] Vier Deklarations-Flächen tragen die dritte Modul-Angabe:
      `AGENTS.md` §4, Sensors-Tabelle, `##`-Hilfetext, Profil-Kopfkommentar —
      `make gate-consistency` grün.
- [ ] Eine ADR begründet Verortung (Rezept statt `modules:`), Deckung und die
      **Nicht**-Änderung der Fence-Lexik; im [ADR-Index](../../adr/README.md)
      eingetragen.
- [ ] [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) trägt die
      Neu-Messung als `## Geschichte`; `make adr-check` grün.
- [ ] [slice-178](../open/slice-178-offene-tasks-roh.md) §1 nennt die richtige
      Zähl-Einheit; [`BEO-020`](../observations.md) ist fortgeschrieben.
- [ ] `make gates` und `make fullbuild` grün (Exit explizit); **unabhängiger
      Review**; **Verifikation** gegen DoD — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Ein drittes Modul am Bindepunkt ist ein drittes, das rot werden kann.**
  Heute 0 Befunde, aber `spans` prüft den ganzen `done/`-Bestand mit; ein
  künftiger Slice mit unbalanciertem Absatz blockiert dann die Closure statt
  nur die Prosa zu trüben. Das ist gewollt und trotzdem ein Risiko. —
  **Ausgang:** *(bei Closure)*
- **Die Messung „0 Abweichungen in 676 Dateien" ist eine Aussage über
  **diesen** Bestand.** Ein Adopter mit anderer Schreibweise kann die Divergenz
  haben; dieser Slice sagt über ihn nichts. — **Ausgang:** *(bei Closure)*
- **`spans` deckt den umschließenden Span nicht**, und die DoD sagt das auch
  nicht zu. Wer die Tabelle in §1 quer liest, könnte „Closure-Bindepunkt sieht
  jetzt alles" mitnehmen. — **Ausgang:** *(bei Closure)*
- **Der Bindepunkt bleibt außerhalb von `gates`.** Ein Defekt fällt erst bei
  `make fullbuild` auf, also spät. Das ist die Bindepunkt-Trennung aus
  [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) und
  keine Regression — aber es begrenzt, was dieser Slice einlöst. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt
keinen Slice.

**Rückführungen:** `in-progress` → `open`, falls die Vorher/Nachher-Probe
zeigt, dass `spans` am Bindepunkt Bestands-Befunde erzeugt, die niemand
einplanen wollte — dann ist die Frage zuerst der Bestand und nicht das Gate.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `harness/` (Gate-Deklaration und Sensors-Tabelle) und
  die Repo-Konfiguration (`Makefile`, `.d-check.closure.yml`). Beide fallen
  unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration). **Kein Produkt-Code** ist berührt — das Modul existiert,
  es wird nur an einem zweiten Ort gefahren.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-30, höchste Kennung
  `BEO-024`): [`BEO-020`](../observations.md) — **Zähler 3, Schwelle
  erreicht**: gemessen wird die eine Menge, ausgesagt wird über eine andere.
  Dieser Slice **korrigiert eine Instanz davon** (slice-178 §1 zählt
  abschnittsweise, der Mechanismus paart absatzweise) und ist zugleich selbst
  exponiert — seine tragende Zahl ist „676 Dateien, 0 Abweichungen", also eine
  Aussage über **diesen** Bestand; sie steht als Risiko in §5.
  [`BEO-013`](../observations.md) — ein Wächter, der nichts mehr fängt: die
  dritte Tabellenzeile in §1 ist genau das, und sie ist als **Nicht**-Ziel
  benannt. [`BEO-023`](../observations.md) — **Zähler 3**: die
  Vorher/Nachher-Probe der DoD ist die Antwort darauf, nicht ein grüner Lauf.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen:** entfällt in `open/` — der Block entsteht
  **spätestens bei der Beanspruchung** (`open→in-progress`)
  ([`MR-053`](../../../../harness/conventions.md#mr-053)).

Slice-ID: slice-180. Betroffene IDs:
[`DC-FA-SPAN-001`](../../../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in),
[`DC-FA-CLI-012`](../../../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben).
Module: `spans` (an einem zweiten Bindepunkt). Gates: `make gates`,
`make gate-consistency`, `make verify-closure-notes`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Eine Rezept-Zeile, vier Deklarations-Flächen
und zwei ADR-Vorgänge; kein Produkt-Code, kein Fremdsystem, keine
Reconciliation, kein Bestand, der umgestellt werden müsste.

## 9. Closure-Notiz (nach `done/`)
