# Slice slice-127: Zwei heimatlose Hard Rules nach AGENTS.md umziehen — danach ist CLAUDE.md ein Pointer

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** ohne Welle — ein einzelner Harness-Hygiene-Punkt mit eigener DoD und
ohne gemeinsame Closure-Bedingung mit anderer Arbeit (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** Baseline-Regelwerk
[`modul-09-implementierung.md` §AGENTS.md-Regeln](../../../../.harness/baseline/v5.11.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)
(Zwei-Quadranten-Regel; AGENTS.md gehört in jeden Lauf-Kontext) und
[`grundlagen-source-precedence.md`](../../../../.harness/baseline/v5.11.0/regelwerk/grundlagen-source-precedence.md)
(Konflikt-Hard-Rule); [`AGENTS.md`](../../../../AGENTS.md) §1 und §6;
[`MR-015`](../../../../harness/conventions.md#mr-015) (AGENTS.md **routet**,
spiegelt nicht) und [`MR-012`](../../../../harness/conventions.md#mr-012).

**Berührte Spec-Stellen:** — (Harness-Dateien; keine Anforderung, kein
Spec-Stratum, kein Produkt-Code).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

**Neuschnitt nach der Rückführung.** Der erste Anlauf wollte `CLAUDE.md` auf
einen Pointer kürzen und stützte sich auf die Zensur *„jede Zeile steht schon
woanders"*. Der unabhängige Review hat sie widerlegt: **zwei** Aussagen stehen
im ganzen Repo nur dort. Damit ist die Aufgabe ein **Umzug**, keine Kürzung —
und die Kürzung ist ihre **Folge**, nicht ihr Inhalt.

Das Regelwerk hat für beide einen Ort, und wir haben ihn nicht benutzt:

> *„Jede Regel, der ein Agent folgen muss, steht in einer gerankten Quelle, im
> Konventionsspeicher oder in der adoptierten Baseline … Artefakte außerhalb
> dürfen verweisen, ausführen und einen dort gerankten Ablauf ausbuchstabieren,
> aber nichts festlegen.“*
> — `grundlagen-source-precedence.md` §Vollständigkeit

Die beiden Waisen:

- **Die Konfliktregel.** `CLAUDE.md` sagt *„Konflikt melden und der
  höherrangigen Quelle folgen"*. [`AGENTS.md`](../../../../AGENTS.md) §1 trägt
  davon nur die halbe Aussage — und zwar über einen engeren Fall (*diese Datei*
  gegen eine kanonische Quelle, nicht zwei kanonische gegeneinander).
  **Der Kanon ist dabei stärker als beide:** *„dass bei Konflikt die niedriger
  rangierte Quelle **angepasst** wird, ist universal (Hard Rule)"*, und ein
  Widerspruch *„gehört benannt"*. Der Umzug ist deshalb keine Verschiebung,
  sondern eine **Angleichung**.
- **Die Benenn-Pflicht.** *„Vor der Implementierung benennen: Slice-ID,
  betroffene `DC-*`-IDs, ADR-IDs, betroffene Module, auszuführende Gates."*
  Der kanonische Schritt 3 verlangt nur, Requirement-/ADR-IDs zu
  **identifizieren**. **Das Delta ist schmaler, als es aussieht:** die
  Baseline-Slice-Vorlage trägt `Bezug:`, `Berührte Spec-Stellen:` und
  `## 3. Plan (vor Code)`, deckt also IDs und Plan bereits ab; übrig bleiben
  **Module und Gates vorab benennen**. Auch dafür gibt es einen kanonischen
  Anker an anderer Stelle — die Baseline-Vorlage der Projekt-README verlangt
  einen `Gates:`-Punkt mit derselben Warnung („halluzinierte Gates sind die
  häufigste Form von Harness-Lüge"). Kanonisch ist also die **Sorge**,
  repo-lokal ihre Verortung **pro Slice**. Heute lebt die volle Form nur in
  `.claude/commands/implement-slice.md` — außerhalb jeder gerankten Quelle.

**Der Kanon stützt die Pointer-Form ausdrücklich.** Die Baseline-Vorlage der
Projekt-README schreibt für ihren Harness-Signal-Block: *„Pointer auf die
kanonischen Quellen — **Inhalt nicht wiederholen**."* Dieselbe Regel, eine
Datei weiter.

**Dass die Baseline `CLAUDE.md` nicht kennt, ist kein Loch.** Sie verlangt, dass
AGENTS.md *in jedem Lauf-Kontext* liegt; welche Datei ein Agenten-Tool dafür
automatisch lädt, ist Werkzeug-Sache. Genau deshalb darf dort nichts **nur**
stehen: eine Datei ohne Rang in der Source Precedence ist kein Ort für eine
Hard Rule.

**Korrektur einer Prämisse des ersten Anlaufs:** die dort als „Spannung"
gemeldete Leseordnung ist **kein Konflikt**. Die Baseline weist §Leseordnung
ausdrücklich *„für den neuen Menschen"* aus, und
[`harness/README.md`](../../../../harness/README.md) nennt sie selbst die
„Menschen-Hälfte des Einstiegs"; der 8-Schritt-Pfad in
[`AGENTS.md`](../../../../AGENTS.md) §6 ist der **Agenten**-Workflow. Zwei
Adressaten, zwei Ordnungen. Es gibt nichts aufzulösen und nichts zu melden.

## 2. Vorgehen

1. **[`AGENTS.md`](../../../../AGENTS.md) §1 — Konfliktregel auf den Kanon
   ziehen:** die niedriger rangierte Quelle wird **angepasst** (nicht nur
   „gilt nicht"), und der Widerspruch **gehört benannt**. Gilt auch zwischen
   zwei kanonischen Quellen, nicht nur gegen diese Datei.
2. **[`AGENTS.md`](../../../../AGENTS.md) §6 — Schritt 3 verschärfen:** vor der
   Implementierung Slice-ID, `DC-*`-IDs, ADR-IDs, betroffene Module und
   auszuführende Gates **benennen**. Als Verschärfung des kanonischen Schritts
   markiert und mit Herkunfts-Anker `(seit slice-127)` versehen
   (`modul-09` §AGENTS.md-Regeln: Hard Rules aus dem Steering Loop tragen ihn;
   ohne Welle die Slice-Form).
3. **Kein `MR-`Eintrag.** Das war im ersten Neuschnitt falsch instrumentiert:
   dieses Repo führt für repo-lokale Hard Rules **keine** MRs — §3.2
   (Suppression-Verbot) hat keinen, und der einzige MR, der
   [`AGENTS.md`](../../../../AGENTS.md) berührt, ist
   [`MR-015`](../../../../harness/conventions.md#mr-015) über die **Rolle** der
   Datei, nicht über einen Regel-Inhalt. MRs sind hier für Form- und
   Konventions-Deltas reserviert. Eine Hard Rule kommt nach `modul-09` in
   AGENTS.md mit Herkunfts-Anker — das ist das ganze Instrument.
4. **Erst danach `CLAUDE.md` auf den Pointer reduzieren** — und die Zeile so
   formulieren, dass sie **stimmt**: der erste Anlauf versprach Routing zur
   „Leseordnung", ein Wort, das in AGENTS.md gar nicht vorkommt.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine neue Hard Rule.** Beide Regeln existieren bereits; sie wechseln den
  Ort und werden dabei auf den Kanon angeglichen.
- **Kein Gate für die beiden.** Nach der Zwei-Quadranten-Regel bleiben sie
  damit **halb durchgesetzt** — das ist als Grenze zu benennen, nicht zu
  verschweigen und nicht durch einen Heuristik-Wächter zu übertünchen.
- **Keine `dpin`-Absicherung**, kein Aktivieren des Moduls `pins` (zöge den
  `BEO-010`-Nachzug nach).
- **Keine Änderung an den `.claude/`-Hooks** und keine Aussage darüber, was sie
  abdecken — der erste Anlauf hat sie überdehnt (der Guard greift nur bei
  `Bash`, `stop-require-gates.sh` hat einen zweiten Freigabepfad).

## 4. Definition of Done

- [x] [`AGENTS.md`](../../../../AGENTS.md) trägt beide Regeln. Die Konfliktregel
      in der **kanonischen** Fassung (*die niedriger rangierte wird angepasst*)
      — **nach Review korrigiert:** die Melde-Pflicht steht daneben als
      **unsere** Regel, nicht als Kanon-Zitat; der Kanon nennt sie nur im
      MR-gegen-Baseline-Fall. **Ohne Herkunfts-Anker**, ebenfalls nach Review:
      beide Regeln sind Bestand seit dem Bootstrap, und
      `grundlagen-traceability.md` verbietet das Nachrüsten — *„der leere
      Zustand **ist** die ehrliche Information."*
- [x] `CLAUDE.md` trägt Titel und genau eine Anweisung; die Vorfassung ist
      zeilenweise belegt. **Mit einer Korrektur:** die Tabelle führte acht
      *Zeilen*, nicht acht *Aussagen* — eine **neunte** steckte **innerhalb**
      einer Zeile und ging verloren (die `conventions.md`-Lesepflicht für
      Code-Änderungen). Sie ist zurück, ebenso die Nuance *Gate-Ausgabe* statt
      *Gate-Ausführung*.
- [x] Die Pointer-Zeile verspricht nur, was AGENTS.md einlöst — die
      „Leseordnung" des ersten Anlaufs ist raus.
- [x] `make gates` Exit 0 (acht Gates, 463 Dateien, 0 Befunde); **zwei**
      unabhängige Reviews, einer je Anlauf
      ([erster](../../../reviews/2026-08-23-slice-127-claude-md-pointer-review.md),
      [zweiter](../../../reviews/2026-08-23-slice-127-umzug-review.md)), beide
      blockierend, alle Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Der Umzug macht AGENTS.md länger.** — **Ausgang:** *nicht eingetreten.*
  Der Review hat die Rechtfertigung am Bestand geprüft: beide Aussagen sind
  wirklich Waisen. Netto ist die Datei sogar **kürzer** geworden, weil im
  selben Zug acht Slice-Verweise und drei veraltete Zukunftsform-Sätze
  gefallen sind.
- **Zwei Hard Rules ohne Gate.** — **Ausgang:** *eingetreten wie erwartet und
  benannt.* Beide liegen allein im Feedforward-Quadranten;
  [`MR-031`](../../../../harness/conventions.md#mr-031) schreibt die Grenze
  ausdrücklich hin, statt sie mit einem Heuristik-Wächter zu übertünchen.
- **Die Vollständigkeits-Aussage ist erneut die Nagelprobe.** — **Ausgang:**
  **erneut eingetreten, aber eine Ebene feiner.** Die Tabelle stimmte auf
  Zeilen-Ebene und war auf **Aussagen**-Ebene zu grob: eine Zeile bündelte
  sechs Punkte, eine andere trug zwei Aussagen, und eine neunte steckte
  innerhalb einer Zeile. Der Beleg stand vor dem Löschen — nur in der falschen
  Auflösung.

## 6. Trigger

**Start** (`next` → `in-progress`): **blockiert — der Wartegrund hat am
2026-08-23 gewechselt.** Der ursprüngliche (die CR-Frage) ist **erledigt**: der
Kurs hat den Konsumenten-CR dieses Repos in **Kurs-Welle 94** beantwortet und
mit [`v5.11.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v5.11.0)
veröffentlicht. Die Werkzeug-Einstiegsdatei hat jetzt eine **kanonische Rolle**
(*„kein Bindepunkt, sondern der Einstieg … Sie verweist dorthin und legt nichts
fest"*), und die Rangliste trägt die **Vollständigkeits-Zusage** samt Prüffrage
und der Regel **Waise ⇒ umgezogen, nicht gelöscht**.

**Der zweite Wartegrund ist seit dem 2026-08-23 ebenfalls erledigt:**
[slice-128](../done/slice-128-baseline-v5110-vendoring.md) hat den Pin auf
`v5.11.0` gehoben, und die kanonische Rolle der Werkzeug-Einstiegsdatei samt
Vollständigkeits-Zusage liegt damit im **gepinnten** Baum. Der Slice ist
**startbar**.

**Ein Reihenfolge-Vorbehalt bleibt, und er ist keine Blockade:**
[slice-129](../done/slice-129-baseline-v5110-delta-audit.md) beantwortet unter
anderem, **welche** Artefakte außerhalb der Rangliste normativen Text tragen —
`CLAUDE.md` ist ein bekannter Eintrag dieser Liste, aber ob der einzige, weiß
erst das Audit. Wer zuerst startet, ist eine Entscheidung, keine Ableitung:
vorher heißt einen bekannten Fall früh erledigen, nachher heißt alle Fälle
einmal gemeinsam schneiden.

**Warum pausiert statt durchgezogen:** die Änderung ist in **sieben** Repos
fällig, nicht in einem. Eine Form zu setzen, bevor der Kanon sie kennt, hieße
sie siebenmal zu setzen und danach eventuell siebenmal zu korrigieren.

**Rückführungen:** `in-progress` → `next`, falls der Review eine **dritte**
Aussage findet, die nirgends sonst steht — dann ist die Zensur wieder falsch
gewesen und der Schnitt taugt nicht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Dateien (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-011`](../observations.md) ist die **zentrale** Beobachtung dieses
  Slice — sein erster Anlauf ist die siebte Instanz der Klasse. Die
  Vollständigkeits-Aussage ist deshalb zeilenweise zu belegen, **bevor**
  gelöscht wird, und der Beleg gehört in die Closure-Notiz.
  [`BEO-002`](../observations.md) für die Ränder (Verweise auf CLAUDE.md,
  Zitate ihres Inhalts in frozen Reports).

Slice-ID: slice-127. Betroffene IDs:
[`MR-015`](../../../../harness/conventions.md#mr-015),
[`MR-012`](../../../../harness/conventions.md#mr-012). Module: Harness-Dateien.
Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Harness-Dokumentation nach kanonischer
Form; die Zwei-Quadranten-Regel des Regelwerks gibt den Ort vor.

## 9. Folge-Punkte (außerhalb dieses Slice)

- **Konsumenten-CR an `ai-harness-course` — der Slice wartet darauf.** Zwei
  Punkte, **verschieden belegt**, und die Trennung gehört in den CR:
  - **Form der Werkzeug-Einstiegsdatei** (reiner Pointer, kein Eigeninhalt),
    in `grundlagen-harness-dateien.md` §Verzeichniskonvention. **Evidenz
    stark und gemessen:** 3–477 Zeilen über sieben Repos, vier mit einer Hard
    Rule in einer Datei ohne Rang, und das Repo des Kanons praktiziert die
    Form bereits, ohne dass sein Regelwerk sie nennt.
  - **Schritt 3 des 8-Schritt-Pfads** um die Vorab-Nennung von Modulen und
    Gates erweitern. **Evidenz schwach:** 2 von 7 Repos, davon einer ein
    Klon; der Zensus misst die Verbreitung des *Textes*, nicht den *Bedarf*.
    Das Argument ist eine benennbare **Asymmetrie** — Schritt 8 verlangt den
    Bericht über gelaufene Sensors, nichts verlangt die Vorab-Nennung, obwohl
    die Baseline halluzinierte Gates selbst „die häufigste Form von
    Harness-Lüge" nennt. Im CR **benennen, nicht fordern**; ohne
    Schadensmessung wäre die Forderung eine Aussage aus dem Anlass
    ([`BEO-011`](../observations.md)).
- **Die Prozedur ist übertragbar, das Ergebnis nicht.** Jede der sieben
  Dateien trägt anderen Inhalt (`bernstein` 477 Zeilen), und der erste Anlauf
  dieses Slice ist genau daran gescheitert, die Kürzung ohne Beleg zu
  behaupten. Übertragbar ist deshalb die **Reihenfolge**, nicht die Zeile:
  (1) zeilenweise belegen, wo jede Aussage in einer **gerankten** Quelle
  steht; (2) die Waisen zuerst nach `AGENTS.md` umziehen, dabei am Kanon
  ausrichten; (3) **erst dann** kürzen. Wer mit Schritt 3 anfängt, löscht
  bindende Regeln.
- **`a-check` trägt dieselben zwei Waisen-Regeln** (Klon unserer Fassung, nur
  `AC-*` statt `DC-*`). Fremdes Repo, eigener Slice dort; hier nur benannt.
  Dieselbe Frage stellt sich für `b-trace`, `m-trace` und `ai-harness-init`.

## 10. Closure-Notiz (nach `done/`)

Geliefert: zwei Waisen stehen in `AGENTS.md`, `CLAUDE.md` trägt vier Zeilen und
verweist. Die Reihenfolge, die der Kanon seit Kurs-Welle 94 vorschreibt —
*belegen, umziehen, erst dann kürzen* —, ist eingehalten.

**Dieser Slice hat zwei Anläufe gebraucht, und das ist sein Ergebnis.** Der
erste wollte kürzen und stützte sich auf die Zensur *„jede Zeile steht schon
woanders"*. Sie war falsch; ohne den Review wären zwei Hard Rules gelöscht
worden, und kein Gate hätte es gemeldet. Die Rückführung nach `next/` war kein
Umweg, sondern die einzige Stelle, an der ein falsch geschnittener Slice
billig bleibt.

**Die Lehre des zweiten Anlaufs ist feiner und unangenehmer.** Ich habe den
zeilenweisen Beleg diesmal geführt — und er war trotzdem zu grob. Acht
*Tabellenzeilen* sind nicht acht *Aussagen*: eine bündelte sechs Punkte, eine
trug zwei, und die neunte, die verloren ging, steckte **innerhalb** einer
Zeile (die `conventions.md`-Lesepflicht galt für Code **und** Doku, die
Ersatz-Fundstelle nur für Doku). **Eine Vollständigkeits-Prüfung ist so gut wie
ihre Auflösung** — und die Auflösung wählt, wer prüft, nicht der Gegenstand.

**Zweimal habe ich den Kanon für mich sprechen lassen, wo er es nicht tut.**
Der Anker `(seit …)` datierte einen **Umzug** als Ursprung, obwohl beide Regeln
seit dem Bootstrap bestehen und `grundlagen-traceability.md` das Nachrüsten
ausdrücklich verbietet — *„der leere Zustand **ist** die ehrliche
Information."* Ich habe den Anker daraufhin nicht entfernt, sondern **ersetzt**
(Slice → Baseline-Bump), und erst der zweite Hinweis hat gezeigt, dass jeder
Anker falsch ist. Und *„der Widerspruch gehört benannt"* steht im Kanon nur im
Fall MR-gegen-Baseline — dem einen, in dem gerade **nicht** angepasst wird. Die
Melde-Pflicht ist unsere; sie steht jetzt ohne fremde Autorität da.

**Ein Nachzug, der nicht geplant war und trotzdem hierher gehört:** Auf
Auftraggeber-Hinweis sind acht Slice-Verweise aus `AGENTS.md` und
`harness/README.md` gefallen — normative Einstiegsdateien verweisen auf
Normatives und ADRs, nicht auf Planung. Drei davon waren nicht bloß Verweise,
sondern veraltete **Zukunftsform** („die Datei entsteht mit …"), seit Monaten
Gegenwart. Was bleibt, sind zwei `seit welle-<NN>` — die eine Anker-Form, die
[`AGENTS.md`](../../../../AGENTS.md) §3.7 und der Kanon ausdrücklich zulassen.

**Offen und benannt:** Beide umgezogenen Regeln liegen allein im
Feedforward-Quadranten. Kein Gate prüft, ob ein Lauf die fünf Felder genannt
oder einen Widerspruch gemeldet hat.
[`MR-031`](../../../../harness/conventions.md#mr-031) sagt das hin, statt es zu
heilen: ein Wächter auf Botschafts-Prosa wäre ein behauptetes Gate — und genau
diese Forderung hat der Kurs in unserem eigenen CR zu Recht abgelehnt.
