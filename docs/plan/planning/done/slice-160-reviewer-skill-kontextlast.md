# Slice slice-160: Der Reviewer-Skill trägt sechzehn Anker

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`.harness/skills/reviewer.md`](../../../../.harness/skills/reviewer.md);
Baseline
[`modul-10-review-harness.md`](../../../../.harness/baseline/v5.12.0/regelwerk/modul-10-review-harness.md)
§Ziel-Form: Reviewer-Skill;
[slice-131](welle-83/slice-131-reviewer-skill-waisen.md) (der Waisen-Zensus);
[slice-147](../done/slice-147-reviewer-anker-reichweite.md) (der Anlass).

**Berührte Spec-Stellen:** — (Rollen-Skill; das Produkt bleibt unberührt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

Der Reviewer-Skill trägt **16** benannte HIGH-/MEDIUM-Anker — 7 HIGH, 2 einzeln
stehende MEDIUM, 7 MEDIUM in einem Sammelpunkt. Jeder einzelne ist gegen einen
belegten Fall geschrieben; zusammen sind sie eine Leseliste, die vor jedem
Review vollständig im Kopf sein müsste.

**Die Sorge ist nicht neu, sie ist nur nie gemessen worden.** Mehrere Slices
tragen den Risikosatz *„der fünfzehnte wird nicht mehr gelesen"* als
Abnahme-Punkt und beantworten dann die Nebenfrage („schärft dieser Anker einen
bestehenden?") statt der gestellten. Die Schwelle war beim sechzehnten Anker
längst überschritten.

**Die erste Frage ist eine Messung, keine Umstrukturierung:** greifen die Anker
überhaupt? Über die abgelegten Reports lässt sich zählen, welcher Anker wie oft
zu einem Befund geführt hat — und welcher seit seiner Einführung **nie**.

## 2. Vorgehen

1. **Zählen, bevor umgebaut wird.** Je Anker: wie viele Befunde in
   `docs/reviews/` lassen sich ihm zuordnen, und seit wann steht er? Ein Anker
   ohne Treffer ist nicht automatisch wertlos — er kann die Klasse verhindert
   haben —, aber die Zahl gehört auf den Tisch, bevor jemand kürzt.
2. **Die Kategorien-Semantik der Baseline gegenprüfen.** `modul-10` gibt
   HIGH/MEDIUM/LOW/INFO vor; ob 16 repo-eigene Anker darunter die Ziel-Form
   noch treffen oder sie überschreiben, ist zu lesen, nicht anzunehmen.
3. **Erst dann die Form entscheiden.** Kandidaten: Zusammenzug verwandter Anker;
   eine Zwei-Ebenen-Form (kurze Prüffragen-Liste vorn, Begründungen hinten);
   Auslagerung der Begründungen in die Registerzeilen, auf die sie ohnehin
   zeigen. Jede Form hat einen Preis — der gehört benannt.
4. **Kein Anker wird gestrichen, ohne dass seine Klasse einen anderen Ort hat.**
   Sonst entsteht die Waise in der Gegenrichtung: die Regel steht in
   `AGENTS.md`, und niemand prüft sie mehr.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine neue Anker-Klasse.** Der Slice räumt, er ergänzt nicht.
- **Keine Änderung an den Hard Rules selbst.** Was in `AGENTS.md` §3/§5 steht,
  bleibt; hier geht es um die Feedback-Hälfte.
- **Kein Zusammenzug ohne Messung.** „Wirkt verwandt" ist kein Grund; zwei
  Anker mit unvereinbaren Arbeitsanweisungen gehören auseinander, auch wenn sie
  dieselbe Kategorie tragen.

## 4. Definition of Done

- [ ] Je Anker eine **Trefferzahl** aus `docs/reviews/`, mit dem Zeitraum seit
      seiner Einführung. — **Nicht erfüllt, und der Haken war falsch gestellt.**
      Der Kanon nennt eine andere Quelle: *„Ein Archiv-Scan ist nicht nötig —
      die Häufung steht im Register"*
      ([`modul-10`](../../../../.harness/baseline/v5.12.0/regelwerk/modul-10-review-harness.md)).
      Aus dem [Beobachtungs-Register](../observations.md) liegen vier Zahlen
      vor — `BEO-009` **6**, `BEO-012` **5**, `BEO-011` **4**, `BEO-004` **3** —,
      für die übrigen zwölf Klassen keine, weil sie dort keinen Eintrag haben.
      Der Zeitraum steht für keine.
- [x] Die Form ist entschieden und ihr **Preis benannt**: zwei Ebenen können
      driften, und sie taten es im ersten Anlauf sofort. Der Preis steht im
      Skill, nicht nur hier.
- [x] Kein gestrichener Anker ohne benannten Ersatz-Ort — **es ist keiner
      gestrichen**. Der Haken ist damit trivial erfüllt, und das ist keine
      Leistung, sondern das Ergebnis von §5.
- [x] `make gates` grün (Exit explizit); unabhängiger Review, Urteil *„nicht
      schließbar"*, vierzehn Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Die Trefferzahl misst die Vergangenheit, nicht den Wert.** Ein Anker ohne
  Treffer kann seine Klasse verhindert haben; wer nach Zahlen kürzt, entfernt
  womöglich genau den, der wirkt. Die Zahl ist ein Eingang in die Entscheidung,
  nicht die Entscheidung. — **Ausgang: eingetreten, und schärfer als
  geschrieben.** Der Punkt warnte, die Zahl könne den *Wert* verfehlen. Sie
  verfehlte den **Gegenstand**: gezählt wurden `befund:`-Felder, die nur 104
  der 214 Reports führen, das `quelle`-Feld war nicht darin — und genau dort
  steht der Anker-Name —, und die handgezählten Reports tragen gar kein
  `befund:`-Feld. Die beiden Zahlen, die als Beleg gegeneinandergestellt
  wurden, hatten **keine Schnittmenge**. Der Fall ist als zweite Instanz von
  [`BEO-020`](../observations.md) gezählt, mit der benannten Abweichung: dort
  war es die *fremde* Menge, hier die *eigene Teilmenge*.
- **Kürzen ist die Bewegung, die sich gut anfühlt.** Ein kürzerer Skill liest
  sich besser und prüft schlechter; der Slice darf nicht am Umfang gemessen
  werden, sondern an der Frage, ob eine Klasse ihren Ort behält. —
  **Ausgang: eingetreten, in einer Form, die der Punkt nicht vorsah.** Gekürzt
  wurde kein Text — gekürzt wurde **Sichtbarkeit**: die erste Fassung der
  Arbeitsfläche trug 7 von 16 Klassen und **2 von 7 HIGH**. Fünf
  merge-blockierende Klassen standen nicht darauf, darunter der
  Stilles-Grün-Pfad. Der Punkt fragte, ob eine Klasse ihren **Ort** behält;
  alle behielten ihn. Die Frage, die gefehlt hat, ist, ob sie noch **gefunden**
  wird.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung ergibt, dass die
Anker greifen und die Last tragbar ist — dann ist die ehrliche Lieferung eine
gestrichene Sorge, keine Umstrukturierung.

**Der Trigger ist zur Hälfte eingetreten und greift trotzdem nicht.** Die
Register-Zähler belegen die erste Hälfte: die Anker greifen. Die zweite ist
**nie gemessen worden** — es gibt keine Last-Größe außer der Zeilenzahl, und
kein Kriterium, ab wann sie „tragbar" ist. Eine Rückführung hätte damit auf
einer ungemessenen Hälfte gestanden, und die Sorge wäre mit derselben
Begründung gestrichen worden, mit der sie entstand. Geliefert ist stattdessen
die Ordnung; die Last-Frage bleibt offen und ohne Messgröße.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Rollen-Skills (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-011`](../observations.md) — die Form gehört aus dem Bestand, nicht aus
  dem einen Anker, der zuletzt dazukam;
  [`BEO-012`](../observations.md) — vor jedem Zitat der Baseline deren
  Geltungsbereich lesen.
- **Nachtlauf-Stand lesen** (`make nightly-state`, dritte Vorprüfung nach
  [`MR-053`](../../../../harness/conventions.md#mr-053)): **ROT**, jüngster
  Lauf `2026-08-27T10:49:23Z`, `head_sha 48cf132`. **Gelesen:** derselbe Lauf,
  den [slice-164](../done/slice-164-nachtlauf-kadenz.md) §7 bereits eingeordnet
  hat — er lief vor den sechs Pin-Hebungen aus
  [slice-161](../done/slice-161-sechs-pins-heben.md), die zum Dispatch noch
  nicht auf `origin/main` lagen. Keine neue Meldung; der nächste Lauf ist die
  Probe darauf, und die steht noch aus.

Slice-ID: slice-160. Betroffene IDs: — (kein `DC-`-Bezug; Rollen-Skill).
Module: — . Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Aufräumen an einem vorhandenen Artefakt der
Review-Schicht.

## 9. Closure-Notiz (nach `done/`)

**Der Slice wollte messen, bevor er kürzt. Die Messung lag bereits vor — an der
Stelle, die der Kanon dafür benennt, und nicht dort, wo ich gesucht habe.**

**Die Häufung steht im Register, nicht im Archiv.**
[`modul-10`](../../../../.harness/baseline/v5.12.0/regelwerk/modul-10-review-harness.md)
sagt es in dem Absatz, der die Ziel-Form des Reviewer-Skills definiert: *„Der
Report ist ein Lauf-Beleg, kein Wissensspeicher … das steuerungsrelevante
Signal ist die Finding-Klasse, die über die Slice-Closure §7 ins
Beobachtungs-Register wandert … Ein Archiv-Scan ist nicht nötig — die Häufung
steht im Register."* Ich habe den Archiv-Scan gefahren und das
[Register](../observations.md) nicht gelesen. Es trägt die Zähler für vier der
sechzehn Klassen: `BEO-009` **6**, `BEO-012` **5**, `BEO-011` **4**, `BEO-004`
**3** — alle vier an oder über der Schwelle. Damit war die Frage von §2
Schritt 1 für die vier beantwortbar, bevor eine Zeile Scan lief.

**Der Scan maß eine andere Menge, als er aussagte.** Er lief über
`befund:`-Felder. Die tragen nur **104** der 214 Reports; das `quelle`-Feld —
genau das Feld, in dem der Anker-Name steht — war nicht darin; **307** der
Zeilen datieren vor der Einführung des Ankers; und die fünfzehn jüngsten
Reports tragen gar kein `befund:`-Feld. Die Handzählung lag damit **vollständig
außerhalb** des Scans. Der Satz *„um eine Größenordnung falsch"* stellte zwei
Zahlen gegeneinander, die keine gemeinsame Grundgesamtheit haben. Gezählt als
zweite Instanz von [`BEO-020`](../observations.md), mit der Abweichung: dort
war die verwechselte Menge eine **fremde**, hier eine **eigene Teilmenge**.

**Zwei Sätze im Skill waren schlicht falsch, und beide führten vom Befund
weg.** *„Die Reports zitieren keine Anker-Namen"* — sieben tun es wörtlich,
fünf davon im `quelle`-Feld eines Findings. *„Der Provenance-Marker ist im
gesamten lebenden Bestand nirgends gesetzt"* —
[ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md) Zeile 77
trägt ihn, liegt **nicht** in `matrix.exempt-paths` und wird ohne ihn zu einem
`matrix-forbidden`-Befund; der Review hat das mit dem Produkt gegen eine
Repo-Kopie gegengeprüft. Die [`.d-check.yml`](../../../../.d-check.yml) sagt
zwei Zeilen unter der Ausnahme-Liste selbst *„neue ADRs ab 0022 tragen den
Marker"*. Der zweite Satz stand im Arbeitsflächen-Teil und hätte den Reviewer
an der **einzigen** Stelle vorbeigeführt, an der die Frage zu stellen ist.

**Die erste Fassung der Arbeitsfläche kürzte, ohne Text zu kürzen.** Sie trug
7 von 16 Klassen und **2 von 7 HIGH**. Nicht darauf standen fünf
merge-blockierende Klassen, darunter der Stilles-Grün-Pfad, den
[`AGENTS.md`](../../../../AGENTS.md) §4 *„die häufigste Form von Harness-Lüge"*
nennt — und für die neun herabgestuften Klassen war **keine** Zahl erhoben
worden. Eine Einstiegs-Ordnung, die an den blockierenden Klassen vorbeiführt,
ist schlechter als keine. Die zweite Fassung trägt alle sechzehn, jede Frage so
gestellt, dass „ja" ein Finding ist; vorher hieß „ja" viermal bestanden,
zweimal Befund.

**Das Dokument ist länger geworden — 172 auf 206 Zeilen —, und das steht jetzt
darin.** *„Die Last ist gesenkt"* hing an der Annahme, der Leser höre nach der
Tabelle auf; genau das untersagt derselbe Absatz zwei Zeilen darüber. Geliefert
ist eine Ordnung, keine Senkung. Der Mess-Absatz ist ganz aus dem Skill heraus:
ein Lauf-Beleg gehört hierher, im Dauerdokument verfällt er still und ohne
Pfleger.

**Kein Anker ist gestrichen — aber der Grund ist ein anderer als der zuerst
notierte.** Nicht *„die Zählung ist unmöglich"*, sondern: die vier Klassen mit
Registerzähler sind produktiv, und für die zwölf ohne gibt es keine Zahl, auf
der eine Streichung stehen könnte. Das ist der Fall, den §2 Schritt 1
ausdrücklich vorsieht — ein Anker ohne Treffer ist nicht automatisch wertlos.

**Was der Review nicht brechen konnte:** die Baseline-Ziel-Form. Alle sechs
Pflichtteile stehen unverändert, die zweite Ebene ist unangetastet, und eine
vorgeschaltete Zusammenfassung ist eine **Ergänzung**, keine Ersetzung einer
Baseline-Regel — also kein `MR-`Eintrag. Das war die offene Frage aus §2
Schritt 2, und sie ist gelesen statt angenommen worden.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 564 Dateien, 0 Befunde),
`make fullbuild` (Exit 0, 48 Anforderungen / 0 Waisen). Ein unabhängiger Review
ist gelaufen; sein Urteil war *„nicht schließbar"*, seine vierzehn Befunde sind
eingearbeitet, und seine vier HIGH sind eigens nachgemessen statt übernommen —
zwei davon widerlegten Sätze, die ich in ein lebendes Dokument geschrieben
hatte. Den einen Pfad-Verweis, den der Lifecycle-Move zu ziehen hatte, fand
nicht ich, sondern `doc-check`: mein Suchfilter schnitt genau die Zeile weg,
die er finden sollte.
