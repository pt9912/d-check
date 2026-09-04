# Slice slice-139: Ein Risiko ohne Ausgang darf nicht in `done/` liegen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos** (Baseline-Regelwerk
[`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht](../../../../.harness/baseline/v5.11.0/regelwerk/modul-06-roadmap.md)):
seine Closure-Bedingung wäre seine eigene DoD.

**Bezug:** Baseline-Regelwerk
[`modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst](../../../../.harness/baseline/v5.11.0/regelwerk/modul-05-planning-harness.md)
— *„Ein Slice geht nicht nach `done/`, während ein Risiko ohne Ausgang
dasteht."* Dazu [`BEO-015`](../observations.md);
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) (die
Closure-Notiz-Struktur im Modul `planning`);
[`AGENTS.md`](../../../../AGENTS.md) §4.

**Berührte Spec-Stellen:** — (Harness-Gate; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

Der Kanon verlangt, dass **jeder** offene Punkt beim Übergang nach `done/` genau
einen von drei Ausgängen bekommt, und schließt mit einem harten Satz: *„Ein
Slice geht nicht nach `done/`, während ein Risiko ohne Ausgang dasteht."*

**Diese Regel ist heute vollständig ungewacht** — obwohl ihr Verstoß die
maschinenlesbarste Form hat, die es gibt: ein wörtlicher Vorlagen-Platzhalter in
einer Datei an einem bekannten Ort. Gemessen:

- Ein `done/`-Slice mit unaufgelöstem `*(bei Closure)*` in §5 läuft durch
  `make verify-closure-notes` mit **0 Befunden**.
- Auch nachdem `"(bei Closure)"` in die `boilerplate`-Liste aufgenommen wurde:
  weiterhin **0** — die Liste greift nur in der Closure-Notiz, nicht im
  Risiko-Abschnitt.
- `structure` kennt `non-empty`, `table-order`, `headings-match` — **kein**
  „darf Muster X nicht enthalten".

**Der Bestand ist sauber**, dreifach gemessen: null `*(bei Closure)*`, null
`*(wird mit dem Closure-Body gefüllt)*`, null `<…>` in `done/`-Slices. Der
Wächter startet also grün und wirkt ab dem ersten Verstoß.

## 2. Vorgehen

1. **Ein Skript**, das `done/slice-*.md` gegen die drei Platzhalter-Formen hält
   — die zwei repo-lokalen und die `<…>`-Form der Kanon-Vorlage.
2. **Fail-closed bei leerer Prüfmenge**: findet es keine `done/`-Slices, bricht
   es rot ab. „Nichts gefunden" und „nichts zu prüfen" dürfen im Exit nicht
   gleich aussehen.
3. **Ort:** ein eigenes `make`-Target, gehängt an **`fullbuild`** neben
   `verify-closure-notes` — nicht an `gates`. Begründung: die Regel gilt dem
   **Übergang nach `done/`**, und das ist der Closure-Bindepunkt; dieselbe
   Einordnung, die `completeness-check` und `verify-closure-notes` schon tragen.
4. **Bewusstes Brechen** je Platzhalter-Form: gesetzt ⇒ rot mit **gelesener
   Fundstelle**; Rückbau ⇒ grün. Plus die leere Prüfmenge.
5. `AGENTS.md` §4 und die Sensors-Tabelle nachziehen — sonst
   `gate-consistency`-rot.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Modul-Fähigkeit.** Die abschnittsgenaue Form (`structure` mit einem
  `forbid-match`) ist der spätere, sauberere Weg — Produkt-Delta mit ADR, und
  sie gehört **hinter** einen Zensus über die Closure-Prozedur.
- **Keine Ausweitung auf `boilerplate`.** Den ganzen Slice gegen die
  Floskel-Liste zu halten erzeugte Falsch-Positive: ein Slice-Plan darf
  „fertig" in Prosa enthalten.
- **Keine Prüfung, ob der Ausgang inhaltlich trägt.** Ob *„eingetreten"* die
  richtige Antwort war, ist Urteil. Geprüft wird, dass **überhaupt** einer
  dasteht.

## 4. Definition of Done

- [x] Das Target hält `done/slice-*.md` gegen **vier** Platzhalter-Formen — die
      vierte kam erst durch den Review dazu, und es war ausgerechnet die, um die
      die Regel geht (§9). Fail-closed bei leerer Prüfmenge **und** bei einem
      Leseversagen.
- [x] Es hängt an `fullbuild`, nicht an `gates`; die Begründung **und ihre
      Kehrseite** stehen im Target-Kommentar und in §5.
- [x] **Sechs** konstruierte Verstöße rot gesehen, jeder mit gelesener
      Fundstelle: vier Platzhalter-Formen, die leere Prüfmenge und die
      unlesbare Datei; Rückbau je grün.
- [x] `AGENTS.md` §4 und die Sensors-Tabelle tragen das Target;
      `gate-consistency` grün. **Darüber hinaus:** die Meta-Gates-Klassifikation
      kannte **sechs** heute entstandene Targets nicht — alle eingeordnet.
- [x] `make gates` Exit 0 (zehn Glieder), `make fullbuild` Exit 0 (der Wächter
      läuft darin, 137 Slices); unabhängiger Review
      ([Report](../../../reviews/2026-08-25-slice-139-closure-ausgang-waechter-review.md)),
      blockierend mit **zwei HIGH**, alle vier Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Ein Wächter auf eine Zeichenkette ist so gut wie seine Liste.** —
  **Ausgang:** *eingetreten, sofort, und genau an der teuersten Stelle.* Meine
  Liste hatte drei Formen und ließ die aus, um die die Regel geht. Ich hatte die
  Vorlage nach dem *Wortlaut* meiner eigenen Platzhalter durchsucht statt nach
  **ihren** — der Beleg lag in der Datei, die ich zitiert habe.
- **`fullbuild` statt `gates` heißt: der Verstoß fällt später auf.** —
  **Ausgang:** *eingetreten wie erwartet, benannt, nicht geheilt.* Die
  Einordnung bleibt richtig: die Regel gilt dem Übergang, nicht dem
  Arbeitsbaum. Der Preis steht im Target-Kommentar.
- **Der Wächter prüft nur `done/`.** — **Ausgang:** *eingetreten und enger als
  gedacht.* Er sagt nichts über `open/` — und, wie der Review zeigte, auch
  nichts über die Formen, die der Platzhalter-Erkenner des Produkts abdeckt.
  Beide Grenzen stehen jetzt im Skript-Kopf, die zweite mit
  Auflösungs-Trigger.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls sich zeigt, dass die
Platzhalter-Formen im Bestand uneinheitlich sind — dann ist ihre Vereinheitlichung
ein eigener Slice und der Wächter hinge an einer geratenen Liste.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Werkzeuge (GF), Gate-Landschaft (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25):
  [`BEO-015`](../observations.md) ist der **Anlass** dieses Slice — er schließt
  die Feedback-Hälfte zu genau der Regel, deren vierten Ausgang der Eintrag
  benennt. [`BEO-007`](../observations.md) für jeden Beleg-Lauf.
  [`BEO-010`](../observations.md), weil ein neues Target in drei Doku-Flächen
  erscheint.

Slice-ID: slice-139. Betroffene IDs: — (Harness-Gate; keine Anforderung).
Module: Harness-Werkzeuge, Gate-Landschaft. Gates: `make gate-consistency`,
`make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — neuer Wächter auf eigenem Bestand.

## 9. Closure-Notiz (nach `done/`)

Geliefert: `make closure-outcomes` hängt an `fullbuild` und hält die
`done/`-Slices gegen vier Platzhalter-Formen. Die Drei-Ausgänge-Regel des
Baseline-Regelwerks hat damit eine Feedback-Hälfte — die **urteilsfreie**; ob ein
eingetragener Ausgang inhaltlich trägt, bleibt Urteil und ist ausdrücklich nicht
Gegenstand.

**Der Wächter gegen stille Grün-Pfade hatte selbst einen, zum zweiten Mal in
zwei Slices.** Meine Liste kannte drei Platzhalter-Formen und ließ die aus, um
die die Regel geht: die Vorlage schreibt das Ausgang-Feld als
`<eingetreten: … | entfallen: … | weiter offen: …>`. Ein wörtlich kopiertes,
unaufgelöstes Risiko wäre unentdeckt durchgelaufen. **Der Grund ist der
unangenehme Teil:** ich habe die Vorlage nach dem Wortlaut **meiner eigenen**
Platzhalter durchsucht statt nach **ihren** — und den Beleg dabei in derselben
Datei übersehen, die ich im Slice-Plan zitiere.

**Und ein maskierter Fehler, den `set -euo pipefail` nicht fängt.** Ein
`|| true` hinter einer Pipe verschluckte ein **Leseversagen** vollständig: 0
Befunde, Exit 0, ohne die Datei je gesehen zu haben. Das ist `BEO-007` eine
Ebene tiefer — dort war es der Exit *hinter* der Pipe, hier der Fehler *davor*.
Lesbarkeit und `sed`-Erfolg werden jetzt geprüft; die Probe mit einer
unlesbaren Datei meldet rot.

**Der schwerste Befund ist ein Selbstwiderspruch im Abstand weniger Stunden.**
Das Produkt trägt in `checkClosurePlaceholder` eine fence- und
inline-code-bewusste Platzhalter-Erkennung, deren Kommentar ausdrücklich sagt:
*„dieselbe geteilte Lexik wie überall, **kein Nachbau**"*. Ich habe einen
Nachbau gebaut — und genau dieses Argument am selben Tag benutzt, um ein
fremdes Skript abzulehnen, weil *„wir das im Produkt haben"*.

Die Messung macht das Bild feiner, ohne den Vorwurf zu entkräften: das Muster
des Produkts verlangt **whitespace-freie** Winkelklammern und sieht die
Ausgang-Zeile deshalb nicht; außerdem prüft es nur den **Abschnitt** der
Closure-Notiz. Zwei der vier Formen liegen außerhalb seiner Reichweite. Die
verbleibende Überschneidung ist im Skript-Kopf **benannt** statt übergangen,
mit Auflösungs-Trigger: sobald der Abschnitts-Skopus des Moduls den ganzen Slice
umfasst, fällt dieses Skript ersatzlos. **Deckung wegzunehmen, um Doppelung zu
vermeiden, tauschte eine echte Prüfung gegen ein Reinheitsargument** — das ist
die Abwägung, und sie steht da, damit der nächste Leser sie prüfen kann.

**Ein Nebenfund, der größer war als sein Befund.** Der Review beanstandete zwei
nicht nachgezogene `fullbuild`-Spiegel. Beim Nachziehen zeigte sich, dass die
Meta-Gates-Klassifikation **sechs** an diesem Tag entstandene Targets gar nicht
kannte. Alle sechs sind jetzt eingeordnet, mit ihren Bindepunkten — eine Fläche,
die niemand bewacht, und deshalb genau die Klasse, die
[`BEO-010`](../observations.md) führt.

**Offen und benannt:** Die Platzhalter-Liste ist eine Liste. Ändert die Vorlage
ihre Form, schweigt der Wächter — sie gehört beim nächsten Vorlagen-Bump
mitgeprüft. Das ist keine Restlücke, die man wegkonfigurieren kann, sondern der
Preis eines Zeichenketten-Wächters; die saubere Form wäre der Abschnitts-Skopus
im Produkt.
