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

**Verantwortlich:** pt9912 (Implementer-Rolle, beansprucht 2026-08-30).

**Autor:** pt9912. **Datum:** 2026-08-30.

---

## 1. Ziel

**Der Closure-Bindepunkt kann zwei Defekte nicht sehen, die das Produkt längst
findet.** `make verify-closure-notes` fährt `planning` und `structure` über das
`done/`-Profil. Beide lesen den **bereinigten** Abschnitts-Text. Ein
**vergessener Schluss-Fence** verschluckt damit alles dahinter — und die
Bedingung, die darüber etwas zusagt, wird **still** wahr.

**„Das Produkt findet sie längst" ist wörtlich zu nehmen, und das entscheidet
den Wert dieses Slice** (nachgemessen nach dem Review): `spans` steht im
Hauptprofil, `make doc-check` läuft in `gates` **und** im `pre-commit`-Hook, und
die 546 Dateien des Bindepunkts sind eine **Teilmenge** seiner 608. Der Zuwachs
an gefundenen Defekten ist damit **null** — gekauft wird die **Unabhängigkeit**
des Bindepunkts von einem fremden Profil, nicht neue Deckung
([ADR-0077](../../adr/0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md)).

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
Beide über **678 Pfade** gegeneinander gefahren (`find -L` über den Baum: 674 reguläre Dateien plus vier aufgelöste `.claude/rules/`-Symlinks, davon 17 gitignorierte unter `.harness/cache/` — eine **Obermenge** der 661 getrackten):
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
4. **Eine ADR trägt den Entscheid und die Messung** ([ADR-0076](../../adr/0076-spans-am-closure-bindepunkt.md), nach dem Review abgelöst durch [ADR-0077](../../adr/0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md)). Warum `spans` an diesen
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

- [x] `make verify-closure-notes` fährt `spans`; **gemessen vorher und
      nachher**: 546 Dateien, 0 Befunde, Exit 0 unter beiden Rezepten — der
      Bestand bleibt unverändert. Dass `spans` dabei nicht bloß schweigt,
      belegt die Probe darunter, nicht der grüne Lauf.
- [x] **Die neue Deckung ist belegt, nicht behauptet:** ein echter
      `done/`-Slice als Fixture, unverändert **grün** (Exit 0), mit angehängtem
      offenem Fence **rot** (`fence-unclosed`, Exit 1). Die **erste** Probe war
      schon vorher rot und belegte damit nichts — verworfen
      ([`BEO-023`](../observations.md)). Die Verifikation hat eine **Kontrolle**
      ergänzt, die zeigt, dass das Vorher-Grün nicht leer ist.
- [x] Vier Deklarations-Flächen tragen die dritte Modul-Angabe und **alle drei**
      Grund-Codes: `AGENTS.md` §4, Sensors-Tabelle, `##`-Hilfetext,
      Profil-Kopfkommentar — `make gate-consistency` grün (611/0). **Benannte
      Grenze, gemessen:** dieser Lauf hält sie **nicht** gegeneinander —
      `targets` vergleicht Target-Namen, keine Beschreibungstexte.
- [x] [ADR-0077](../../adr/0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md)
      begründet Verortung, **Unabhängigkeit statt Deckung** und die
      **Nicht**-Änderung der Fence-Lexik; sie löst
      [ADR-0076](../../adr/0076-spans-am-closure-bindepunkt.md) ab, deren
      tragende Begründung der Review widerlegt hat. Beide im
      [ADR-Index](../../adr/README.md).
- [x] [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) trägt die
      Neu-Messung als `## Geschichte`; `make adr-check` grün (611/0).
- [x] [slice-178](../open/slice-178-offene-tasks-roh.md) §1 steht auf dem
      Gemessenen: die **ursprünglichen** Zahlen 25/45 sind richtig, falsch war
      die Folgerung „ungerade ⇒ exponiert". [`BEO-020`](../observations.md) auf
      Zähler **4** — die vierte Instanz ist meine eigene Fehl-Korrektur.
- [x] `make gates` (611 Dateien, 0 Befunde, Exit 0) und `make fullbuild` grün;
      [**unabhängiger Review**](../../../reviews/2026-08-30-slice-180-closure-profil-spans-review.md)
      (3 HIGH · 2 MEDIUM · 2 LOW · 2 INFO, blockierend) und
      [**Verifikation**](../../../reviews/2026-08-30-slice-180-closure-profil-spans-verifikation.md)
      (A-1 bis A-8) in eigenen Kontexten gelaufen, alle Befunde eingearbeitet
      und die strittigen Messungen vorher selbst nachgefahren.

## 5. Abnahme-Punkte / Risiken

- **Ein drittes Modul am Bindepunkt ist ein drittes, das rot werden kann** —
  über **drei** Grund-Codes: `fence-unclosed`, `span-unclosed` und
  `span-nested-link`. Heute 0 Befunde, aber `spans` prüft den ganzen
  `done/`-Bestand mit; ein künftiger Slice mit unbalanciertem Absatz blockiert
  dann die Closure statt nur die Prosa zu trüben. Das ist gewollt und trotzdem
  ein Risiko. — **Ausgang:** weiter offen, und **um einen Grund-Code breiter**,
  als dieser Punkt annahm: `span-nested-link` steht seit dem Review in allen
  vier Deklarations-Flächen. Bestands-Rauschen weiterhin null.
- **Die Messung „0 Abweichungen in 676 Dateien" ist eine Aussage über
  **diesen** Bestand.** Ein Adopter mit anderer Schreibweise kann die Divergenz
  haben; dieser Slice sagt über ihn nichts. — **Ausgang:** weiter offen, und der
  Nenner ist jetzt benannt — 678 Pfade (`find -L`), eine Obermenge der 661
  getrackten. Über einen fremden Bestand sagt der Slice weiterhin nichts, und
  die Messung selbst ist modellierungs-kritisch: unter der falschen
  Öffner-Lesart ergäbe derselbe Baum eine divergente Datei.
- **`spans` deckt den umschließenden Span nicht**, und die DoD sagt das auch
  nicht zu. Wer die Tabelle in §1 quer liest, könnte „Closure-Bindepunkt sieht
  jetzt alles" mitnehmen. — **Ausgang:** weiter offen. Der Review hat den Punkt
  bestätigt und verschärft: nicht nur `spans` deckt den umschließenden Span
  nicht — **kein** Modul kann es, weil das Item dort per CommonMark Code ist.
  Die Nicht-Zusage steht jetzt in vier Flächen statt in dreien.
- **Der Bindepunkt bleibt außerhalb von `gates`.** Ein Defekt fällt erst bei
  `make fullbuild` auf, also spät. Das ist die Bindepunkt-Trennung aus
  [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) und
  keine Regression — aber es begrenzt, was dieser Slice einlöst. —
  **Ausgang:** **eingetreten, und anders als erwartet** — dieser Punkt sagte,
  der Defekt falle erst bei `make fullbuild` auf. Gemessen fällt er beim
  **Commit** auf, weil `make doc-check` in `gates` und im `pre-commit`-Hook
  läuft und die Scan-Menge des Bindepunkts eine Teilmenge ist. Der Zuwachs an
  gefundenen Defekten ist null. Aufgefangen von
  [ADR-0077](../../adr/0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md),
  die die Begründung auf die **Unabhängigkeit** des Bindepunkts stellt und
  [ADR-0076](../../adr/0076-spans-am-closure-bindepunkt.md) ablöst.

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

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-30T06:08:17Z,
  `image-scan.yml` 2026-08-29T10:07:43Z. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption, kein
  Baseline-Abschnitt ([`MR-054`](../../../../harness/conventions.md#mr-054)).

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

**Geliefert, und die Begründung dafür ist eine andere als die geplante.** Der
Closure-Bindepunkt fährt `spans`; ein vergessener Schluss-Fence macht dort
keine Zusage mehr still wahr. **Der Zuwachs an gefundenen Defekten ist
null** — das Hauptprofil meldet denselben `fence-unclosed` schon beim Commit,
und die 546 Dateien des Bindepunkts sind eine Teilmenge seiner 608. Gekauft ist
die **Unabhängigkeit** von einem fremden Profil, nicht neue Deckung.

**Die tragende Begründung war falsch, und ein Review hat sie widerlegt.**
[ADR-0076](../../adr/0076-spans-am-closure-bindepunkt.md) nannte den
vergessenen Fence *„die einzige, die heute niemand sieht"*. Der Slice-Plan
hatte es in §1 genauer — *„zwei Defekte, die das Produkt **längst** findet"* —,
und auf dem Weg in die ADR ist die Einschränkung verlorengegangen. Weil das der
**Grund** der Entscheidung ist und kein Detail, geht die Korrektur nicht als
Errata durch, sondern als Folge-ADR
([ADR-0077](../../adr/0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md),
`supersedes`) — Entscheid des Auftraggebers aus drei vorgelegten Optionen, die
dritte war die Rücknahme.

**Zwei weitere Zusagen derselben ADR hielten nicht.** Sie schrieb, `make
gate-consistency` halte den `##`-Hilfetext gegen die Doku; gemessen vergleicht
`targets` Target-**Namen**, und der Lauf bleibt grün, wenn man den Hilfetext
verfälscht **oder** `--enable spans` aus dem Rezept entfernt. Es hält also
**kein** Sensor die vier Deklarations-Flächen gegen das Rezept — das steht
jetzt als benannte Grenze statt als Zusage. Und `non-empty: true` macht den
Fall, auf dem der Entscheid steht, **nicht** laut; die verworfene Alternative
war schwächer, als ich ihr zugestand.

**Die teuerste Lehre ist eine über das Korrigieren.** Meine Richtigstellung an
[slice-178](../open/slice-178-offene-tasks-roh.md) §1 ersetzte zwei
**richtige** Zahlen durch zwei falsche: mein Zähl-Bereich lief von `## 4.` bis
`## 5.`, während beide Slices ihre DoD als **§3** führen — gezählt war §4
*Risiken*. Falsch war an der ursprünglichen Aussage nie die Zahl, sondern die
Folgerung *„ungerade ⇒ exponiert"*: eine ungerade Backtick-Zahl heißt, dass der
letzte Backtick nichts schließt, also wird gar nichts entfernt. Verschluckt
wird ein Item nur von einer Spanne, die es **umschließt**. Positiv-Kontrolle
gefahren — an denselben Absatz ein `- [ ]` angehängt, die Regel meldet, trotz
der 25 Backticks.

**Und der zweite Fehlschluss derselben Art im selben Absatz:** *„`spans` meldet
nichts, also ist die Exposition null"*. `spans` schweigt bei einer
**wohlgeformten** Spanne konstruktionsbedingt; sein Schweigen belegt keine
Abwesenheit. Das Fazit stimmt trotzdem — repo-weit verschwinden sechs
Task-Items durch wohlgeformte Spannen, alle bewusste Prosa über den Marker.
Beides ist [`BEO-020`](../observations.md), und beides fand der zweite Leser.

**Was der Slice über das Werkzeug hinaus gezeigt hat:** eine Messung, die eine
Zahl korrigiert, muss zuerst den **Bereich** prüfen, aus dem die alte stammt.
Eine falsche Zahl und ein falscher Ausschnitt sehen im Ergebnis identisch aus —
und wer den Ausschnitt nicht nennt, kann den Fehler nicht sehen. Das ist die
Prozedur, die `BEO-020` seit heute zusätzlich trägt.

**Die Fence-Frage bleibt, wo sie war.** Beide Lesarten des Produkts liegen über
678 Pfaden ohne eine einzige Abweichung; der offene Punkt (a) aus
[ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) ist erneut
gemessen und bleibt unbelegt. Die Verifikation hat dabei belegt, dass die
Messung **modellierungs-kritisch** ist: unter der falschen Öffner-Lesart ergäbe
derselbe Baum eine divergente Datei — die veröffentlichte Null ist die der
richtigen.

**Fortgeschrieben:** [`BEO-020`](../observations.md) auf Zähler **4** (die
vierte Instanz entstand beim Korrigieren der dritten),
[`BEO-023`](../observations.md) auf **4** — dort als **Antwort** statt als
Defekt: die erste Vorher/Nachher-Probe war schon vorher rot und belegte nichts.

**Was offen bleibt und wohin es gehört:** dass kein Sensor das Rezept gegen die
vier Deklarations-Flächen hält, ist benannt und nicht geschlossen — ein
eigener Schnitt, falls er je gewollt ist.
[slice-178](../open/slice-178-offene-tasks-roh.md) ist damit entsperrt, seine
Dringlichkeit aber gesenkt: er steht auf dem konstruierten Fall, nicht auf
einem Bestands-Fund.
