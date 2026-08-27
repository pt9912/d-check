# Slice slice-147: Ein Anker gegen Zitat-Reichweite — Feedforward, nicht nur Feedback

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`BEO-012`](../observations.md) (Zähler 4); Baseline-Regelwerk
[`grundlagen-klassifikation.md` §Die 2×2-Matrix](../../../../.harness/baseline/v5.12.0/regelwerk/grundlagen-klassifikation.md#2x2-matrix)
(Feedforward ↔ Feedback); `.harness/skills/reviewer.md` §Repo-spezifische Anker
pro Kategorie.

**Berührte Spec-Stellen:** — (Harness-Regeltext; keine Anforderung des Produkts).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

[`BEO-012`](../observations.md) steht bei Zähler **4** — eine Quelle wird über
ihren Geltungsbereich hinaus zitiert. **Alle vier Male fand es der zweite
Leser.** Das ist die entscheidende Eigenschaft der Klasse: sie ist im Review
zuverlässig zu finden und beim Schreiben zuverlässig zu übersehen, weil ein
Zitat wie ein Beleg aussieht.

Der Reviewer-Skill trägt für sie **keinen** Anker, während er für
[`BEO-009`](../observations.md) und [`BEO-004`](../observations.md) je einen
MEDIUM-Eintrag mit Arbeitsanweisung führt. Die vierte Instanz zeigt zugleich,
dass Feedback allein nicht reicht: sie entstand in einem Dokument, dessen
Gegenstand diese Klasse **ist**.

## 2. Vorgehen

1. **Zuerst die Kategorie entscheiden**, nicht den Wortlaut. `BEO-009`(b) und
   `BEO-004` stehen als MEDIUM; ob Reichweite dieselbe Stufe trägt, ist ein
   Urteil und gehört begründet.
2. Anker im Reviewer-Skill mit **Prüffrage an den Diff**, nicht mit einer
   Fallliste — die Fallliste ist bei `BEO-004` dreimal unvollständig gewesen.
3. **Die Feedforward-Hälfte prüfen:** trägt `AGENTS.md` §5 die Regel schon, oder
   ist der Skill wieder ihr einziger Ort? Beides ist Pflicht
   (Baseline `modul-09-implementierung.md`: *„Jede Hard Rule liegt in zwei
   Quadranten"*; die 2×2-Matrix selbst steht deskriptiv in
   `grundlagen-klassifikation.md`);
   eine Waise im Skill ist
   genau der Befund, den [slice-131](../done/slice-131-reviewer-skill-waisen.md)
   behandelt hat.
4. **Messen statt annehmen**, ob eine mechanische Hälfte existiert: die
   Wortlaut-Probe aller Blockzitate eines Dokuments gegen die zitierte Quelle
   ist in [slice-140](../done/slice-140-konsumenten-cr.md) einmal von Hand
   gelaufen und hat einen Befund gebracht, den beide Leser übersahen. Sie findet
   die **Tilgung ohne Auslassungszeichen** — nicht die überzogene Reichweite.
   Was sie deckt und was nicht, gehört gemessen, bevor jemand ein Target baut.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Gate auf Reichweite.** Ob ein Satz weiter trägt als seine Quelle, ist
  ein Urteil, kein `grep` — die Registerzeile sagt das seit ihrer Anlage.
- **Kein Nachziehen der vier Belege.** Sie liegen in `done/` und sind Protokoll.
- **Keine Ableitung aus dem Zähler.** Die Registerzeile hält ausdrücklich fest,
  dass die Entscheidung über die Form eigens zu treffen ist; ein Zähler von 4
  ist ein Anlass, kein Beschluss.

## 4. Definition of Done

- [x] Die Kategorie ist **begründet** gewählt, nicht von `BEO-009` abgeschrieben.
- [x] Der Anker steht als **Prüffrage**, nicht als Fallliste.
- [x] Beide Quadranten belegt: Feedforward-Ort benannt, Feedback-Ort benannt —
      oder ausdrücklich ausgewiesen, warum einer entfällt.
- [x] Die Deckung der Wortlaut-Probe ist **gemessen** ausgewiesen: was sie
      findet und was sie nicht findet.
- [x] `make gates` grün (Exit explizit); unabhängiger Review. *(Der Haken stand
      im Feature-Commit, **bevor** der Review lief — `BEO-009`(a), und
      slice-146 zwei Commits vorher genauso. Er trägt jetzt, weil der Review
      gelaufen ist; die Reihenfolge war trotzdem falsch.)*

## 5. Abnahme-Punkte / Risiken

- **Ein Anker mehr im Skill ist kein Schutz, sondern Kontext-Last.** Der Skill
  trägt bereits mehrere HIGH- und MEDIUM-Anker; der fünfzehnte wird nicht mehr
  gelesen. Zu prüfen ist, ob dieser Anker einen bestehenden **schärft**, statt
  neben ihn zu treten. — **Ausgang: eingetreten, und die Prüfung fällt gegen den
  Zusammenzug aus.** Der nächstliegende Kandidat wäre der `BEO-009`(b)-Anker:
  beide melden eine überdehnte Reichweite. Sie prüfen aber verschiedene
  Gegenstände — dort die eigene **Messung**, hier eine **fremde Quelle** —, und
  die Arbeitsanweisungen sind unvereinbar: „suche die N+1-te Form" gegenüber
  „lies das Geltungs-Feld statt des Titels". Ein zusammengezogener Anker trüge
  zwei Prüffragen unter einer Überschrift und würde für beide schlechter
  gelesen. Er steht deshalb daneben, mit ausdrücklichem Verweis auf die geteilte
  Kategorie-Begründung.

  **Die Last selbst bleibt davon unbeantwortet, und das war der eigentliche
  Punkt.** Gezählt: **16** benannte HIGH-/MEDIUM-Anker nach dieser Änderung, 15
  davor — der neue ist der sechzehnte, und die Schwelle des Risikosatzes („der
  fünfzehnte wird nicht mehr gelesen") war damit schon vorher überschritten.
  Ausgang mit Kennung:
  [slice-160](../open/slice-160-reviewer-skill-kontextlast.md).
- **Die Klasse ist beim Schreiben blind — auch mit Anker.** Der Anker wirkt im
  Review, also erst nachdem der Fehler geschrieben ist. Ob die
  Feedforward-Hälfte überhaupt greifen kann, ist die eigentliche Frage dieses
  Slice und darf nicht im Anker untergehen. — **Ausgang: eingetreten, und die
  Feedforward-Hälfte fehlte vollständig.** Gemessen: die Regel stand **nirgends**
  normativ — weder in [`AGENTS.md`](../../../../AGENTS.md), noch in
  [`harness/README.md`](../../../../harness/README.md), noch im
  Konventionsspeicher; sie existierte nur als Beobachtung. Ein Anker allein wäre
  damit genau die Skill-Waise gewesen, die
  [slice-131](../done/slice-131-reviewer-skill-waisen.md) behandelt hat. Sie
  steht jetzt als Hard Rule in §5. **Ob sie greift, ist damit nicht behauptet** —
  vier Instanzen entstanden, während die Klasse im Register stand, eine davon in
  einem Dokument über die Klasse selbst.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Kategorie-Frage eine
Regelwerk-Klärung braucht — dann hängt sie am ausstehenden Konsumenten-CR
([slice-140](../done/slice-140-konsumenten-cr.md), Punkt 4).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Regeltext (GF), Reviewer-Skill (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25):
  [`BEO-012`](../observations.md) ist der Anlass;
  [`BEO-011`](../observations.md) für jede Aussage darüber, was der Anker
  „immer" fängt.

Slice-ID: slice-147. Betroffene IDs: — (kein `DC-`/`ADR-`-Bezug). Module:
Harness-Regeltext. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Regeltext-Arbeit am eigenen Skill.

## 9. Closure-Notiz (nach `done/`)

**Beide Hälften stehen — und der erste Anlauf hat begründet, warum, mit einem
Satz, der die Klasse dieses Slice begeht.** Er lautete: die Regel stehe
**nirgends** normativ. Gemessen waren drei Orte
([`AGENTS.md`](../../../../AGENTS.md),
[`harness/README.md`](../../../../harness/README.md), Konventionsspeicher — dort
steht sie wirklich nicht), geschlossen wurde universal. **Der Kanon trägt sie
wörtlich:**
`grundlagen-source-precedence.md` §Wie weit trägt ein zitierter Satz — *„Sie ist
an **jede** zitierte Aussage zu stellen, auch an einen Satz der Baseline"*. Und
[slice-149](../done/slice-149-baseline-v5120-delta-audit.md) hatte sie mir
vorgelegt: als Delta-Punkt 4 des `v5.12.0`-Audits, mit slice-147 ausdrücklich
als Träger der Feedforward-Hälfte. Die Quelle lag in einem geschlossenen
Artefakt und wurde für nicht existent erklärt.

**Die Änderung bleibt, die Begründung ist eine andere.** Sie steht nicht, weil
die Regel fehlte, sondern weil der Kanon **zwei Quadranten** verlangt
(`modul-09-implementierung.md`: *„Jede Hard Rule liegt in zwei Quadranten"*) und
die verkörperte Form der Ort ist, den ein Implementer liest. Damit ist der
Anker nach dem Kriterium von
[slice-131](../done/slice-131-reviewer-skill-waisen.md) auch **keine Waise**,
sondern ein *„buchstabiert eine gerankte Fundstelle aus"* — dessen Zensus
akzeptiert die vendorte Baseline mehrfach als gültige Fundstelle. Der §5-Absatz
trägt den Kanon-Zeiger, wie §3.7 es vorführt.

**Die Kategorie ist begründet, nicht abgeschrieben.** MEDIUM, und zwar aus
demselben Grund wie `BEO-009`(b): die Quelle ist echt, überdehnt ist nur ihre
Reichweite. Die HIGH-Schwelle greift nicht — es entsteht kein stiller Grün-Pfad
in einem Gate, sondern eine Aussage ohne Fundament.

**Der Anker steht daneben, nicht im BEO-009-Eintrag.** Die Prüfung dafür fiel
gegen den Zusammenzug aus: beide Anker melden überdehnte Reichweite, prüfen aber
verschiedene Gegenstände — die eigene **Messung** gegen eine fremde **Quelle** —
und tragen unvereinbare Arbeitsanweisungen. Zwei Prüffragen unter einer
Überschrift würden für beide schlechter gelesen.

**Die mechanische Hälfte ist gemessen, und sie deckt die Klasse nicht.** Das
Modul `citations` vergleicht den normalisierten Zitattext als
**zusammenhängenden Teilstring** der Quell-Spanne. Es findet damit
Wortlaut-Drift und die **Tilgung ohne Auslassungszeichen** — ein gespleißtes
Zitat ist kein zusammenhängender Teilstring. Es findet **nicht**, was dieser
Slice behandelt: ein **wörtlich korrektes** Zitat, das eine Aussage stützen
soll, die es nicht trägt. Dazu drei benannte Grenzen, die größte zuerst: es
prüft **nur per Direktive ausgezeichnete** Zitate — *„kein Prosa- oder
Voll-Scan"* —, ein wörtliches Zitat **ohne** `d-check:cite` ist damit der
Normalfall im Bestand und überhaupt kein Gegenstand. Ebenso ein Verweis ohne
Zitat. Und Zitate unter **16** normalisierten Zeichen prüft es gar nicht. Hinzu
kommt, dass das Modul heute nicht aktivierbar ist
([slice-152](../in-progress/slice-152-citations-scharfschalten.md)).

**Ein Widerspruch im eigenen Register, gemeldet statt übergangen.** Zwei Zeilen
über `BEO-012` steht [`BEO-017`](../observations.md) in derselben Lage — der
Kanon trägt die Regel — und zieht den **umgekehrten** Schluss: *„Darum hier
keine eigene Regelzeile: sie zu duplizieren verstieße gegen `AGENTS.md` §1."*
Beide können nicht gelten. Aufgelöst zugunsten der Zwei-Quadranten-Pflicht:
`AGENTS.md` §1 verbietet **Duplikation ohne Zeiger**, nicht die verkörperte
Form — §3.7 führt genau diese Bauform vor (Regel plus `Kanon:`-Zeiger).
`BEO-017`s Zeile ist damit die zu enge; sie gehört beim nächsten Anfassen
nachgezogen, nicht von hier aus umgeschrieben.

**Und ein zweiter Widerspruch, an derselben Stelle entstanden.** Der Kanon
verlangt für eine Steering-Loop-Hard-Rule den Herkunfts-Anker
`(seit slice-<NNN>)`, wenn keine Welle läuft;
[`MR-045`](../../../../harness/conventions.md#mr-045) — heute geschrieben —
verbietet Slice-Verweise in `AGENTS.md` so breit, dass er ihn miterfasste. Die
höherrangige Quelle ist der Kanon; aufgelöst als
[`MR-050`](../../../../harness/conventions.md#mr-050), nicht durch einen Edit am
Vorgänger.

**Was der Slice nicht behauptet:** dass die Regel greift. Alle vier Instanzen
entstanden, während die Klasse im Beobachtungs-Register stand, die vierte in
einem Dokument über diese Klasse — **und die fünfte in diesem Slice selbst**.
Sie gehört an den Zähler, mit slice-147 als Beleg. Damit hat der Slice den
einzigen Wirksamkeitsnachweis, den er nicht wollte: der Anker, den er einführt,
hat ihn gefunden.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 535 Dateien, 0 Befunde). Der
Reviewer-Skill steht auf `1.11.0`; [`BEO-012`](../observations.md) trägt beide
Orte statt der offenen Lücke.
