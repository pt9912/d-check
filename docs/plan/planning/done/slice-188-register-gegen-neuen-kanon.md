# Slice slice-188: Das Register gegen die neue Beleg-Definition

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Reaktiv: die Quelle des Kanons hat zwei Regeln
gesetzt, die unseren Bestand betreffen; sein Closure-Grund geht über die eigene
DoD nicht hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:** die zwei Quell-Wellen des Adopters `ai-harness-course`
(Bereichssegment; die Schwelle und ihre drei Ausgänge) — **noch in keinem
Release**;
[slice-183](../done/slice-183-baseline-v5150.md) (der Baseline-Bump, der sie
frühestens transportieren kann);
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(die Register-Deckung, die diesen Bestand prüft).

**Berührte Spec-Stellen:** — dieser Slice ändert **Bestand**, keine Spec-Stelle.

**Verantwortlich:** —.

**Autor:** pt9912. **Datum:** 2026-08-31.

---

## 1. Ziel

**Zwei Zähler unseres Registers sind zu hoch, und bis jetzt fehlte die Regel,
die das entscheidet.** Beide Zeilen tragen eine Begründung, beide galten als
*bewusste Abweichung* — der Adopter hat sie gemessen und uns gemeldet, und wir
haben ihm geantwortet, die Anzahl-Prüfung sei deshalb blockiert.

**Die zwei Quell-Wellen lösen das auf, und zwar gegen uns.** Die neue Regel:

- *Zwei Funde **im selben** Vorgang sind eine Gelegenheit, kein zweites
  Auftreten.*
- *Ein Vorkommen **ohne** abgeschlossenen Vorgang bewegt den Zähler nicht — es
  gehört in den Eintrag, benannt, nicht gezählt.*

Angewandt:

| Zeile | heute | nach der Regel | warum |
|---|---|---|---|
| `BEO-023` | 7 | **6** | `slice-178` steht zweimal — ein Vorgang |
| `BEO-004` | 3 | **1** | die übrigen Vorkommen hatten keinen abgeschlossenen Vorgang |

**Der Befund ist damit nicht „der Kanon war unvollständig", sondern „unsere
Zahlen waren falsch".** Sie standen nur ohne Regel unwiderlegbar da. Das ist die
unbequemere und die richtige Lesart, und sie gehört in die Closure-Notiz.

**Zwei weitere Neuerungen treffen denselben Bestand**, gemessen:

- **Belege sind nicht mehr nur Slice-Kennungen** — Welle und Review-Report sind
  ebenfalls abgeschlossene Vorgänge. Heute tragen **alle** 22 Zeilen
  ausschließlich Slice-Kennungen; die Erweiterung kann Zähler **erhöhen**, wo
  ein Vorkommen bisher unbelegt blieb.
- **Der Stand wird eine geschlossene Menge** (*verkörpert · geplant ·
  gestrichen*, darunter `offen`). Heute: **8** von 22 Zellen beginnen in einer
  dieser Formen, **14** sind freie Prosa. Betroffen von der Schwellen-Bindung
  sind die **12** Zeilen mit Zähler ≥ 3.

## 2. Vorgehen

1. **Die Timing-Frage zuerst entscheiden, nicht nebenbei.** Die Regeln liegen
   auf `main` des Adopters und in **keinem Release**; unser Pin steht auf
   `v5.12.0`, und selbst `v5.14.0` enthält sie nicht. Wer sie jetzt zitiert,
   zitiert etwas, das im vendorten Baum nicht existiert.
   **Aber:** die *Abweichung* ist schon unter dem heutigen Pin eine — `v5.12.0`
   verlangt bereits **Anzahl** (so viele Belege wie der Zähler), und 7 ≠ 6 ist
   auch dort falsch. Neu ist nur die **Richtung** der Reparatur. Zu entscheiden
   ist also: unter dem heutigen Pin korrigieren (Richtung aus dem unreleasten
   Kanon vorweggenommen, ausdrücklich als solche benannt) — oder auf ein Release
   warten und den Slice bis dahin ruhen lassen.

   **Entschieden bei Beanspruchung (2026-09-02): die Frage hat sich von selbst
   erledigt.** [slice-189](../done/slice-189-baseline-v5180.md) hat den Pin
   zwischenzeitlich auf `v5.18.0` gehoben — der vendorte Baum trägt die beiden
   zitierten Regeln jetzt wörtlich:

   <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md:105-107 -->
   > **Ein Vorgang zählt einmal — und was keinen hat, zählt gar nicht.** Der
   > Regelfall eines Belegs ist die Slice-Kennung; auch eine Welle und ein
   > Review-Report sind abgeschlossene Vorgänge und taugen als Beleg. Zwei
   > Funde **im selben** Vorgang

   Es gibt also **keine** Vorwegnahme mehr zu benennen — die Korrektur zitiert
   einen Kanon-Abschnitt, der im gepinnten Baum steht, nicht auf `main`. Die
   ursprüngliche Alternative (auf ein Release warten) ist gegenstandslos
   geworden, bevor der Slice sie ziehen musste.
2. **Die zwei Zähler korrigieren**, jeweils mit der Begründung in der Zeile;
   die bisherige Abweichungs-Begründung wird **nicht gelöscht**, sondern
   abgelöst — sie war die ehrliche Notiz ihrer Zeit.
3. **Alle 22 Zeilen gegen die erweiterte Beleg-Definition durchgehen**: gibt es
   Vorkommen, die als Welle oder Review-Report belegbar sind und den Zähler
   **heben**? Das ist Urteil, kein `grep`, und es gehört einzeln begründet.
4. **Die 14 freien Stand-Zellen auf die geschlossene Menge umstellen** — und
   dabei prüfen, welche der 12 Schwellen-Zeilen einen Ausgang tragen, der die
   Form erfüllt.
5. **Erst danach ist die Anzahl-Achse baubar** (Folge-Slice): sie braucht einen
   Bestand, der die Regel erfüllt, sonst startet sie rot.
6. `make gates`; **Review**; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Nicht die Anzahl-Prüfung bauen.** Sie ist der Folge-Slice und braucht diesen
  Bestand.
- **Kein Baseline-Bump.** Der ist [slice-183](../done/slice-183-baseline-v5150.md), und
  er transportiert die neuen Regeln ohnehin nicht.
- **Keine Migration auf die Verzeichnisform.** Das ist der BEO-CR, und er wartet
  auf seine eigenen Bedingungen.
- **Keine Zähler, die aus dem Gedächtnis steigen.** Punkt 3 erhöht nur, wo ein
  abgeschlossener Vorgang benannt werden kann.

## 4. Definition of Done

- [x] Die Timing-Entscheidung aus §2.1 ist **getroffen und begründet** — nicht
      implizit durch das Tun.
- [x] `BEO-023` steht auf **6**, `BEO-004` auf **1**, je mit Begründung in der
      Zeile; die abgelöste Notiz bleibt lesbar.
- [x] Alle **24** aktuell aktiven Zeilen (22 zur Autorenschaft, seither
      `BEO-025`/`BEO-026` dazugekommen — beide bereits unter der neuen Regel
      angelegt) sind gegen die erweiterte Beleg-Definition geprüft; jede
      Erhöhung nennt ihren abgeschlossenen Vorgang.
- [x] Die Stand-Zellen tragen die geschlossene Menge, **mit einer ehrlichen
      Einschränkung statt eines erfundenen Ausgangs:** von den **11** aktuellen
      Schwellen-Zeilen (Zähler ≥ 3 nach der Korrektur) tragen **9** einen
      formgültigen Ausgang (*verkörpert*/*geplant*). Für **4 Fundstellen**
      (`BEO-008`, `BEO-015`, zwei Instanzen von `BEO-020`) war weder ehrlich
      verfügbar — die unabhängige Verifikation hat einen ersten Versuch mit
      „weiter offen" zu Recht als denselben Kanon-Verstoß zurückgewiesen, den
      `BEO-015` selbst beschreibt. Als neue Beobachtung
      [`BEO-027`](../observations.md) registriert statt hier erfunden — ihre
      Auflösung ist ein eigener Folge-Slice, kein Rest dieses.
- [x] `make gates` grün (Exit explizit); **unabhängiger Review**.
- [x] Closure-Notiz mit Lerneintrag; jedes Risiko aus §5 mit Ausgang; die drei
      Paarungen geprüft.

## 5. Abnahme-Punkte / Risiken

- **Zähler zu senken sieht aus wie Vertuschen.** Wer die Historie liest, sieht
  eine Zahl fallen. Die Begründung muss in der Zeile stehen, nicht nur im
  Commit — sonst ist die Korrektur von einer Schönung ununterscheidbar.
  — **Ausgang: entfallen.** Beide Korrekturen tragen die Begründung **in der
  Zeile selbst**, mit der alten Notiz **abgelöst statt gelöscht** davor
  stehend — nachprüfbar ohne den Commit zu lesen. Nicht eingetreten.
- **Punkt 3 ist die Stelle, an der aus einer Korrektur eine Erfindung wird.**
  Einen Zähler zu **heben**, weil sich ein Vorgang „auch als Beleg lesen lässt",
  wäre eine Aussage aus dem Anlass statt aus dem Bestand
  ([`BEO-011`](../observations.md), Zähler 5). — **Ausgang: entfallen.** Jede
  Erhöhung (`BEO-002` auf 9) nennt einen echten, benannten Vorgang
  (slice-188 selbst); keine Zeile wurde ohne Beleg gehoben. Nicht eingetreten.
- **Die Regeln sind unreleast.** Ein Zitat auf sie ist heute nicht auflösbar,
  und ein Slice, der auf eine `main`-Fassung baut, kann veralten, bevor er
  schließt. — **Ausgang: entfallen.** Durch [slice-189](../done/slice-189-baseline-v5180.md)
  gegenstandslos geworden, bevor dieser Slice es ziehen musste — siehe §2
  Punkt 1.
- **Die Stand-Umstellung berührt 14 Zellen mit gewachsener Prosa.** Beim
  Umschreiben geht leicht verloren, was die Zelle *zusätzlich* sagte
  ([`BEO-002`](../observations.md), Zähler 7). — **Ausgang: entfallen.** Jeder
  Ausgang-Marker wurde **angehängt**, nie eine bestehende Zeile umformuliert
  oder gekürzt — gemessen: kein Satz der Alt-Prosa ist verschwunden. **Neu und
  ungeplant eingetreten ist die Schwester-Gefahr:** nicht eine Zelle, sondern
  die **Kopfzeile der Datei** (die Beleg-Regeln in eigener Prosa) blieb beim
  ersten Durchgang stehen — vom Auftraggeber gefunden, als neunte Instanz von
  [`BEO-002`](../observations.md) nachgetragen (Zähler 9).

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — heute hält
[slice-174](../done/slice-174-register-deckung.md) den Slot. **Zusätzlich
die Entscheidung aus §2.1**: entweder ein Release des Adopters, das die zwei
Wellen enthält, **oder** der ausdrückliche Beschluss, unter dem heutigen Pin zu
korrigieren.

**Rückführungen:** `in-progress` → `open`, falls Punkt 3 zeigt, dass die
erweiterte Beleg-Definition mehr Zeilen bewegt als erwartet — dann ist das ein
eigener Schnitt und keine Nebenwirkung dieser Korrektur.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `docs/plan/planning/` (das Register selbst). **Eine**
  Sub-Area, unter dem Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration pro Sub-Area); der Slice berührt weder Code noch Spec. Die
  Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:219-220 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand **2026-09-02** — bei der
  Beanspruchung aufgefrischt, höchste Kennung jetzt `BEO-026`):
  [`BEO-011`](../observations.md) (Zähler **5**, unverändert seit der
  Autorenschaft) — jede **Erhöhung** in Punkt 3 ist eine Aussage, die aus dem
  Bestand kommen muss und nicht aus dem Anlass; [`BEO-002`](../observations.md)
  (Zähler jetzt **8**, war 7 bei der Autorenschaft — zwei weitere Zeilen sind
  seither dazugekommen) — die freien Stand-Zellen sind gewachsene Prosa, und
  beim Umstellen geht leicht verloren, was sie zusätzlich sagten;
  [`BEO-012`](../observations.md) (Zähler jetzt **12**, war 11) — die Frage
  „liegt die Quelle im gepinnten Baum" ist seit dem Baseline-Bump auf
  `v5.18.0` ([slice-189](../done/slice-189-baseline-v5180.md)) mit Ja
  beantwortet (siehe §2 Punkt 1), die Wachsamkeit der Klasse bleibt trotzdem
  für jedes künftige Zitat gültig. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:225-225 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): `image-scan.yml`
  **gruen** (jüngster Lauf 2026-09-02T07:56:37Z). `upstream-drift.yml`
  **ROT** — jüngster Lauf 2026-09-02T05:19:44Z, aber **planmäßig**: die
  Meldung ist der bekannte, informative Fremd-Release-Fund (Go 1.27.0→1.27.1,
  semgrep 1.175.0→1.176.0), kein Zitat-Bruch nach dem Bump und keine
  Regression — nachgelesen im Lauf selbst, nicht nur an der Farbe. Beide
  Achsen sind für diesen Slice ohne Konsequenz: er berührt weder eine
  Toolchain-Version noch eine Zitat-Spanne außerhalb der beiden in §2 Punkt 1
  geprüften. **Dieser Block trägt bewusst keine `cite`-Direktive** — sein Ziel
  ist eine Repo-Adaption, kein Baseline-Abschnitt
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — die eine berührte Sub-Area fällt unter den
Default. Kein Code, keine Spec, keine Reconciliation: der Slice bringt einen
**Bestand** mit einer Regel in Übereinstimmung, die anderswo entschieden wurde.

## 9. Closure-Notiz (nach `done/`)

**Geliefert.** Zwei Zähler-Korrekturen nach der v5.18.0-Kanon-Regel
(„zwei Funde im selben Vorgang sind eine Gelegenheit, kein zweites
Auftreten"; „ein Vorkommen ohne abgeschlossenen Vorgang bewegt den Zähler
nicht"): `BEO-023` 7→6 (Duplikat-Vorgang `slice-178` entfernt), `BEO-004`
3→1 (zwei unbelegte Vorkommen entfernt). Alle 24 aktuell aktiven
Registerzeilen gegen die erweiterte Beleg-Definition (Welle/Review-Report
zählen jetzt auch) durchgegangen — keine weitere Zeile gewinnt einen neuen
Beleg. Elf Schwellen-Zeilen (Zähler ≥ 3) auf die geschlossene Stand-Form
umgestellt, neun davon formgültig (*verkörpert*/*geplant*). Die Kopfzeile
der Registerdatei trägt jetzt die Baseline-Version und die erweiterte
Beleg-Definition, samt einer benannten Grenze: der Kanon sieht eine
Wellen-Archivierung vor (`archiv.zip`, Review-Reports ohne Stub), die
dieses Repo noch nie durchgeführt hat — für die aktuell offene
[`welle-86`](welle-86-closure-uebergang-durchsetzen.md) greift die
Nachrüst-Ausnahme **nicht**, und kein Werkzeug führt die Archivierung aus.

**Was funktioniert hat.** Die Timing-Frage aus §2 Punkt 1 hat sich von
selbst erledigt — [slice-189](../done/slice-189-baseline-v5180.md) hob den
Baseline-Pin zwischenzeitlich auf `v5.18.0`, und der zitierte Kanon-Abschnitt
steht seither wörtlich im vendorten Baum. Die Methode „nur anhängen, nie
umschreiben" für die Stand-Zellen hat den in §5 befürchteten
Informationsverlust vermieden: kein Satz der Alt-Prosa ist verschwunden.

**Was anders lief.** Zweimal fand nicht der Review oder die Verifikation
den ersten Fund, sondern der Auftraggeber selbst — beide Male am eigenen
Beitrag dieses Slices, nicht an einer Nebensache:

1. **Die Kopfzeile der Registerdatei war die falsche Stelle, an der die
   Semantik-Änderung endete.** Alle Zeilen wurden korrigiert; die
   Zusammenfassung der Beleg-Regeln am Dateianfang, die genau diese Semantik
   in eigener Prosa trägt, blieb unangetastet und beschrieb danach eine
   Regel, die im Rest der Datei schon nicht mehr galt. Neunte Instanz von
   [`BEO-002`](../observations.md) (8→9) — der übersehene Rand ist hier
   nicht eine Historie-Zeile, sondern die Selbstbeschreibung des geänderten
   Artefakts selbst.
2. **Die „Nachrüst-Ausnahme ist konform"-Formulierung klang stärker, als sie
   ist.** Sie deckt nur Wellen, die vor der Regel-Einführung schlossen — nicht
   `welle-86`, die aktuell offen ist. Nachgeschärft, mit der konkreten
   Konsequenz benannt (kein Werkzeug, nächste Closure ist die erste betroffene).

Die **unabhängige Verifikation** fand einen dritten, härteren Fund, den weder
ich noch der Review sahen: **„weiter offen" ist oberhalb der Zähler-Schwelle
kein gültiger Ausgang** — der Kanon lässt dort nur *verkörpert* oder
*geplant (mit Kennung)* zu. Vier Fundstellen (`BEO-008`, `BEO-015`, zwei
Instanzen von `BEO-020`) hatten ehrlich keins von beidem verfügbar, und ein
erster Versuch trug dort „weiter offen" ein — genau der Kanon-Verstoß, den
`BEO-015` selbst beschreibt (*„ein offener Punkt bekommt bei der Closure
einen Ausgang, den es nicht gibt"*), nur diesmal als **fehlender** statt
**erfundener** vierter Ausgang. Statt einen Ausgang zu fabrizieren, ist die
Lücke als neue Beobachtung [`BEO-027`](../observations.md) registriert —
ihre Auflösung (ein Sensor, der Schwellen-Zeilen gegen die zwei zulässigen
Ausgangs-Formen hält) ist ein eigener Folge-Slice.

**Steering-Loop-Einträge.** Zwei Register-Zähler erhöht, keiner davon neu
außer dem Fund selbst:

- [`BEO-002`](../observations.md) **auf 9** (Semantik-Änderung im Körper
  nachgezogen, Rand blieb stehen) — neu daran: der übersehene Rand ist die
  Selbstbeschreibung des Artefakts, nicht eine Historie- oder Index-Zeile.
- [`BEO-027`](../observations.md) **neu, Zähler 1** — eine Registerzeile
  übersteht ihre Schwellen-Lese-Schritte ohne zugewiesenen Ausgang; vier
  Fundstellen in einem Vorgang, nach der neuen Vorgangs-Zählung eine
  Gelegenheit. Abgrenzung zu [`BEO-015`](../observations.md) in der Zeile
  selbst: dort ein erfundener vierter Ausgang, hier ein fehlender.

**Register-Paarung geprüft:** beide zitierten Kennungen
([`BEO-002`](../observations.md), [`BEO-027`](../observations.md)) existieren
als Registerzeile mit mindestens einem Beleg. **Anker-Paarung:** entfällt —
dieser Slice verkörpert selbst keine Regel mit `seit slice-188`-Anker.
**Folge-Slice-Paarung:** kein genannter Folge-Slice trägt bereits eine
Datei — die Auflösung von `BEO-027` und die Wellen-Archivierung sind als
künftige Arbeit benannt, nicht als existierende Kennung zitiert.
