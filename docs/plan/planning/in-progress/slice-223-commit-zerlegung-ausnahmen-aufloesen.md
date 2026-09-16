# Slice slice-223: Die Ausnahmen zur Commit-Zerlegung auflösen — die Baseline reicht

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`MR-013`](../../../../harness/conventions.md#mr-013),
[`MR-059`](../../../../harness/conventions.md#mr-059),
[`MR-061`](../../../../harness/conventions.md#mr-061),
[`MR-062`](../../../../harness/conventions.md#mr-062),
[`MR-063`](../../../../harness/conventions.md#mr-063),
[`MR-064`](../../../../harness/conventions.md#mr-064) — die sechs Einträge,
deren Fortbestand dieser Slice prüft, sowie
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(`make planning-check` koppelt Ruhe-Marker und Verzeichnis — die Kopplung, aus
der der Eintrag oben seine Notwendigkeit ableitet).

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle).

**Verantwortlich:** pt9912 (Implementer-Rolle).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Ziel.** `AGENTS.md` §3.3 trägt wieder die Kanon-Form: die Regel, ihre zwei
Fälle, ihre Begründung — **ohne** die sieben Ausnahme-Blöcke. Die sechs
`MR`-Einträge, die diese Blöcke tragen, wandern nach
`harness/conventions/done/`, weil ihre Abweichung sich als **nicht nötig**
erweist. Was an ihre Stelle tritt, ist **ein Satz**: wo die Zerlegung nichts
schützt, greift sie nicht.

**Der Anlass ist eine Auftraggeber-Frage** (*„Warum brauchen wir diese
Abweichung? Ist Baseline nicht gut genug?"*), und die Antwort ist gemessen:

| | Vorlage `v6.9.0` | d-check |
|---|---|---|
| §3.3 | **15** Zeilen | **88** (73 davon Ausnahmen, 83 %) |
| §3 gesamt | 110 Zeilen | 328 |
| `AGENTS.md` gesamt | 10 504 B | 40 994 B |

**Nachgemessen bei der Beanspruchung (2026-09-16), gegen die jetzt aktuelle
Vorlage `v6.9.0` statt der ursprünglich zitierten `v6.6.0`** — der Pin ist
zwischen Anlage und Beanspruchung dieses Slice gehoben worden
([slice-224](../done/slice-224-baseline-v690-bump.md)). Die Zahlen
verschieben sich nur geringfügig (die Zeilendifferenz stammt aus dem
Baseline-Bump selbst, byte-genau: siehe [`MR-072`](../../../../harness/conventions.md#mr-072)), der **Befund bleibt
unverändert**: `AGENTS.template.md` §3.3 hat die zwei Kanon-Fälle unverändert
seit `v6.6.0` (byte-identisch geprüft), und die Schlussfolgerung unten trägt
weiter.

**Zwei Kanon-Stellen tragen den Entscheid, und beide sind wörtlich geprüft.**

Erstens der **Ort**: Adaptionen gehören nicht als Text nach `AGENTS.md`.
*„Von außen wird der Index adressiert, nicht die Eintrags-Datei — wer aus
`AGENTS.md`, einem Slice oder einer ADR auf eine Adaption zeigt, verlinkt
`harness/conventions.md#mr-<NNN>`"* (`grundlagen-harness-dateien.md`
§Konventionsspeicher). Die Begründung dort ist der **Lesepfad** und ein
Korrektheits-Risiko: *„Ein aufgelöster Eintrag liest sich wie ein geltender."*
Wir zahlen den Preis derzeit **doppelt** — einmal im Index, einmal
ausgeschrieben.

Zweitens der **Wert**: §3.3 nennt seine Begründung selbst — Rename-Erkennung
und `git log --follow`. Im Geltungsbereich der sechs Einträge wird sie **nicht
verletzt**:

- **[`MR-059`](../../../../harness/conventions.md#mr-059), [`MR-062`](../../../../harness/conventions.md#mr-062), [`MR-063`](../../../../harness/conventions.md#mr-063), [`MR-064`](../../../../harness/conventions.md#mr-064)** (Archiv-Stub-Moves): Der Kanon
  sagt selbst, *„der Volltext eines geschlossenen Slice kommt in keinem
  lesenden Knoten vor … Genau deshalb darf er ins Archiv wandern"*
  (`grundlagen-traceability.md`). Wo niemand den Volltext liest, hat
  `git log --follow` **keinen Konsumenten** — die Zerlegung schützt dort
  nichts. [`MR-059`](../../../../harness/conventions.md#mr-059) deklariert ohnehin `Ersetzt-Baseline-Regel: keine`; er
  weicht von **unserem** §3.3 ab, nicht vom Kanon.
- **[`MR-013`](../../../../harness/conventions.md#mr-013)** (drei Fälle): Die Slice-Datei bleibt im Move-Commit
  **unverändert** — die Rename-Erkennung ist unberührt, der Score bleibt 100 %.
  Deviant ist allein, dass *andere* Dateien mitkommen. Für diese Lage hat der
  Kanon eine Antwort: *„Beide Commits gehören in denselben Push. Zwischen ihnen
  ist das Repo kurz rot; das ist zulässig, solange dieser Zwischenstand nicht
  die **Spitze** eines Push wird."*
- **[`MR-061`](../../../../harness/conventions.md#mr-061)** ist ein **anderer Fall** und braucht einen eigenen
  Auflösungsgrund: Er ist nicht *unnötig*, sondern **erschöpft** — sein
  eigenes Feld sagt *„die Regel ist mit dem Vollzug des einen Commits, den sie
  deckt, bereits erschöpft"*.

**Ein gemeldeter Widerspruch gehört dazu** (`AGENTS.md` §1: *„Melde den
Widerspruch, statt ihn stillschweigend nach einer Seite aufzulösen"*): Der
Block *Ausnahme MR-/Wellen-Lifecycle-Move* steht gegen eine **ausdrückliche**
Kanon-Regel — *„Der `git mv` zieht die Pfad-Berichtigung nach sich, **als
eigener Commit nach dem Umzug**"*. [`MR-013`](../../../../harness/conventions.md#mr-013)s Feld `Ersetzt-Baseline-Regel`
nennt `modul-05`, **nicht** diese Stelle; die Abweichung ist also gegen die
falsche Regel deklariert.

**Plan-Änderung bei der Implementierung (2026-09-16) — [`MR-013`](../../../../harness/conventions.md#mr-013) löst sich nur
zur Hälfte auf.** Die obige Analyse *„[`MR-013`](../../../../harness/conventions.md#mr-013) (drei Fälle): Die Slice-Datei
bleibt im Move-Commit unverändert"* trifft nur auf **zwei** seiner drei Fälle
zu (Slice-Lifecycle-Move, Beanspruchung — dort ändert sich die bewegte Datei
nicht, nur andere Dateien reisen mit; das ist plain Git-Semantik und braucht
keine Ausnahme). Der **dritte** Fall — *Ausnahme MR-/Wellen-Lifecycle-Move*
(`conventions/` → `conventions/done/`) — ändert die **bewegte Datei selbst**
(Link-Tiefen-Fixes), fällt damit unter **keine** der beiden Kanon-Bedingungen
und wurde in dieser Sitzung noch aktiv gebraucht: der Move von [`MR-071`](../../../../harness/conventions/done/MR-071-baseline-v660.md#mr-071--baseline-pin-hebung-auf-v660-dreizehnter-nachtrag-zu-mr-011-nachtrag-zu-mr-023) nach
`conventions/done/` (slice-224) brauchte genau diese Bündelung, weil der
lokale `pre-commit`-Hook `doc-check` auf **jedem** Commit fährt — der
kanonische Zwei-Commit-Weg (*„der `git mv` zieht die Pfad-Berichtigung nach
sich, als eigener Commit nach dem Umzug"*, `grundlagen-traceability.md`)
wäre zwischen den beiden Commits lokal gar nicht committierbar, weil der
Zwischenstand (Datei am neuen Pfad, Links noch für den alten berechnet) den
Hook sofort rot macht. **Konsequenz:** [`MR-013`](../../../../harness/conventions.md#mr-013) wird nicht aufgelöst,
sondern auf genau diesen einen Fall **getrimmt**; sein
`Ersetzt-Baseline-Regel`-Feld wird von `modul-05` (falsch, wie §1 oben schon
zeigte) auf die tatsächlich einschlägige Kanon-Stelle korrigiert. DoD (2)
zählt deshalb **fünf** aufgelöste Einträge (059/061/062/063/064), nicht
sechs. Der gemeldete Widerspruch oben besteht als Befund fort — er trägt
jetzt die Begründung, warum die Abweichung **bleibt**, statt warum sie
fällt.

**Abgrenzung — vier Punkte, jeder mit Grund:**

1. **Kein anderer Abschnitt von `AGENTS.md`.** §3 hat 326 Zeilen gegen 109 in
   der Vorlage, §5 hat 170 gegen 14. Dass dort dieselbe Klasse steckt, ist
   **nicht gemessen** — *ein Folge-Slice übernähme es*, und
   [slice-221](../next/slice-221-agents-md-tabellenzellen.md) wartet ohnehin
   auf einen Neuschnitt gegen die Vorlage.
2. **Keine Änderung an `make planning-check`.** Die Kopplung von Ruhe-Marker
   und Verzeichnis bleibt, wie sie ist; nur die **Commit-Granularität** ändert
   sich — *Schicht-Abgrenzung*.
3. **Kein Umschreiben eingefrorener Artefakte.** `done/`-Slices und
   Review-Reports, die [`MR-013`](../../../../harness/conventions.md#mr-013) zitieren, bleiben unangetastet; die
   Index-Anker reisen mit, das ist ihr Zweck. *Bestand bleibt bewusst stehen.*
4. **`tools/archive-wave` wird nicht geändert.** Es darf weiter einen Commit
   erzeugen — die Frage ist, ob das eine **Ausnahme braucht**, nicht ob das
   Werkzeug anders arbeiten soll. *Es wäre ein anderer Vorgang.*

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** `AGENTS.md` §3.3 trägt die Kanon-Form plus **einen** Satz für
      die Klasse, die die Zerlegung nicht braucht (vollständig ersetzter
      Inhalt), und einen **kurzen, korrekt zitierten** Rest-Block für die
      eine Klasse, die sie weiter braucht (Move-Commit ändert die bewegte
      Datei selbst — s. Plan-Änderung). Sechs der sieben Ausnahme-Blöcke sind
      weg; wo ein Leser mehr braucht, steht der Zeiger auf den
      Konventions-Index.
- [ ] **(2)** **Fünf** `MR`-Einträge liegen in `harness/conventions/done/`,
      mit **zwei verschiedenen** Auflösungsgründen: *Baseline-Konformität* für
      059/062/063/064, *erschöpft* für 061. [`MR-013`](../../../../harness/conventions.md#mr-013) bleibt aktiv, getrimmt auf
      den einen verbleibenden Fall, `Ersetzt-Baseline-Regel` korrigiert. Index-Zeilen
      von der aktiven in die aufgelöste Tabelle bewegt, Anker unverändert —
      Präzedenz: [`MR-014`](../../../../harness/conventions.md#mr-014), [`MR-027`](../../../../harness/conventions.md#mr-027), [`MR-038`](../../../../harness/conventions.md#mr-038).
- [ ] **(3)** **Die geänderte Praxis ist gefahren, nicht behauptet:** ein
      Lifecycle-Übergang in der neuen Zwei-Commit-Form, mit gemessener
      Ausgabe — reiner Move (Rename-Score 100 %), dann die gekoppelten
      Verweise; `make gates` auf dem **zweiten** Commit grün, und der rote
      Zwischenstand ausdrücklich benannt.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Die Probe steht vor der Streichung**, nicht danach: Solange nicht gemessen
ist, dass die Zwei-Commit-Form im Alltag trägt, ist die Auflösung eine
Behauptung.

1. **Zwei-Commit-Form fahren** — an diesem Slice selbst, bei seiner eigenen
   Beanspruchung. Der Beleg für DoD (3) entsteht damit aus dem Vorgang, nicht
   aus einer Attrappe.
2. **Zitier-Stellen zählen**, bevor die Einträge wandern: Wer verweist auf die
   sechs, und über welche Form (Index-Anker oder Pfad)? Ein Pfad-Verweis
   bräche beim Move — der Index-Anker nicht.
3. §3.3 umschreiben, Einträge bewegen, Index-Zeilen umhängen.
4. `make gates`, Handoff.

## 4. Trigger

**Beanspruchung:** slice-222 war zum Anlege-Zeitpunkt geschlossen (der Pin
stand damals auf `v6.6.0`, und die Vorlage war die Grundlage dieses
Entscheids — inzwischen gehoben auf `v6.9.0` durch
[slice-224](../done/slice-224-baseline-v690-bump.md), Befund unverändert,
siehe §1 Nachmessung), WIP-Limit frei, Auftraggeber-Freigabe vom 2026-09-08.
**Diese Zeile war zwischen dem `v6.9.0`-Bump und dieser Beanspruchung kurz
veraltet** — ein Gegenwarts-Satz über einen bereits abgelösten Pin,
gefunden vom unabhängigen Review von slice-224 (F-2), hier bei der
Beanspruchung selbst nachgezogen, wie dort empfohlen.

**Rückführung nach `open/`** (`in-progress→open`): wenn die Probe aus Schritt 1
zeigt, dass die Zwei-Commit-Form im Alltag **nicht** trägt — etwa weil ein
Hook oder die CI den Zwischenstand doch sieht. Dann ist [`MR-013`](../../../../harness/conventions.md#mr-013) berechtigt
und nur seine *Ablage* falsch, und der Slice wird neu geschnitten.

**Rückführung nach `next/`** (`in-progress→next`): wenn die Zähl-Schritte aus
Schritt 2 mehr als eine Handvoll Pfad-Verweise finden, die beim Move brechen
— dann ist der Nachzug der eigentliche Vorgang.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `make gates` grün mit echter Ausgabe, unabhängiger
Review durchgeführt und eingearbeitet, Closure-Notiz geschrieben, Register
fortgeschrieben, jedes Risiko aus §6 mit einem der drei Ausgänge.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Kanon-Satz, der [`MR-013`](../../../../harness/conventions.md#mr-013) erübrigt, ist für einen anderen Fall
  geschrieben.** *„Beide Commits gehören in denselben Push"* steht im Abschnitt
  über die **Anker-Paarung** und meint dort den `MR`-Move. Ihn auf den
  Slice-Lifecycle zu übertragen ist eine **Verallgemeinerung**, kein Zitat —
  und genau die Klasse, die
  [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  (18×) führt. Der Slice muss das benennen, statt den Satz als Beleg
  auszugeben. — **Ausgang:** \<offen\>
- **Eine Abweichung aufzulösen ist teurer als sie zu behalten, wenn sie
  gebraucht wird.** Sechs Einträge zu bewegen und die Praxis umzustellen ist
  irreversibel genug, dass ein Irrtum teuer wird: Käme [`MR-013`](../../../../harness/conventions.md#mr-013) zurück, wäre
  er ein neuer Eintrag mit neuer Nummer, und die Historie läge in zwei
  Richtungen. **Die Probe aus Schritt 1 ist die einzige Absicherung** — und
  sie misst **einen** Übergang, nicht die Klasse. — **Ausgang:** \<offen\>
- **Der gemeldete Widerspruch könnte in die andere Richtung aufzulösen sein.**
  `AGENTS.md` §1 sagt, bei Konflikt gewinnt die höherrangige Quelle — hier der
  Kanon. Aber die Möglichkeit, dass **unsere** Form die bessere ist und der
  richtige Weg ein Change Request an die Baseline wäre, ist mit dieser Regel
  nicht ausgeschlossen; sie ist nur nicht der Default. Der Slice entscheidet
  sich für Konformität, und das ist eine **Wahl**, kein Zwang. — **Ausgang:**
  \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende.

**Die drei Vorprüfungen sind bei der Beanspruchung am 2026-09-16 bestätigt.**
**Sub-Area-Wahl:** eine, `*` (Repo-Default) — der Slice ändert `AGENTS.md`
und den Konventionsspeicher, `tools/harness/` ist nicht berührt. **Register**
(40 Verzeichnisse, gemergter Stand): einschlägig ist
[`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(18×, bereits über der Schwelle) — er trägt §1 und §6 dieses Plans bereits;
kein weiterer Treffer erreicht mit diesem Slice erstmals die Schwelle.
**Nachtlauf:** `make nightly-state` am 2026-09-16 — `upstream-drift.yml`
planmäßig rot (derselbe, mit slice-224 bereits gelesene Fremd-Release-Befund,
inzwischen durch den Bump erledigt), `image-scan.yml` grün, ohne Bezug zu
diesem Slice.

**Ein Eintrag ist jetzt schon absehbar** und formt §6:
[`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(18×) — der ganze Entscheid stützt sich auf zwei Kanon-Zitate, und **eines
davon ist für einen anderen Fall geschrieben**. Das steht als erstes Risiko
da, nicht im Bericht danach.
