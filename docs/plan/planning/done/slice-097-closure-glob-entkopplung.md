# Slice slice-097: `planning.closure.glob` — eigener Kandidaten-Filter für die Closure-Fähigkeit

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** gemeinsam mit [slice-098](../done/slice-098-closure-note-placeholder.md) —
der Konsument kann sein letztes handgeschriebenes Prüfskript **erst dann**
zurückziehen, wenn **beide** Slices liegen. Das ist eine Closure-Bedingung
jenseits der beiden Slice-DoDs und damit genau der Fall, für den es eine Welle
gibt. **Zuordnung entschieden** mit der welle-69-Closure (2026-08-09): der Slice
bleibt **eigenständig** und geht nicht im `structure`-Modul auf.

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Closure-Fähigkeit, Spezifikation Schritt C2),
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md).
**Change Request** des Konsumenten `ai-harness-course` (CR 1).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Die Closure-Fähigkeit bekommt einen **eigenen** Kandidaten-Filter
(`planning.closure.glob`, Basisname-Glob, Default = `planning.slice-glob`).
Heute teilen sich zwei Fähigkeiten mit **verschiedenen Grundmengen** einen Knopf;
er ist nicht einstellbar, ohne eine von beiden zu verbiegen.

## 2. Der gemessene Befund

Die Aussagen sind verschieden: die Roadmap-Invariante fragt „liegt hier noch
Arbeit?", die Closure-Fähigkeit „ist jedes abgeschlossene Paket dokumentiert?".
Beide lesen heute `planning.slice-glob`.

Nachgemessen am 2026-08-09 gegen v0.52.0, Endzustand *kein Slice mehr in
Arbeit, Ruhe-Marker gesetzt* (muss grün sein):

```text
slice-glob "slice-*.md"  → 0 Befunde                    ✓
slice-glob "*.md"        → roadmap.md  planning-drift    ✗ falsch-rot
```

Ursache: die Roadmap-Datei liegt selbst im gezählten Verzeichnis und matcht
`*.md`. Die Invariante sieht dauerhaft „eine Datei liegt da" und verlangt, dass
der Ruhe-Marker **nie** gesetzt wird — falsch genau dann, wenn die Welle
wirklich ruht.

**Das ist ein Entwurfsfehler aus
[slice-093](welle-68/slice-093-closure-note-gate.md)**, nicht eine fehlende
Fähigkeit: die zweite Fähigkeit hat den Glob der ersten mitbenutzt, weil beide
zufällig Slice-Dateien meinten. Sobald eine Grundmenge weiter ist, bricht die
Kopplung.

## 3. Abnahme-Punkte

1. **Default = `planning.slice-glob`**, nicht ein eigener Literal-Default. Nur so
   ist der Befundsatz ohne den Schlüssel byte-identisch zu v0.52.0.
2. **Der eigene Bestand ist mitbetroffen — was tun wir damit?** In `done/` liegen
   **9** `welle-*-results.md`, jede mit einer Closure-Notiz; sie sind heute
   unsichtbar. Weitet man den Glob, greift sofort der nächste Unterschied: ihre
   Notiz-Überschrift ist eine **H1** (`# Welle NN — … — Closure-Notiz`), das
   Default-`heading-pattern` verlangt H2/H3 ⇒ **9** `closure-note-missing`
   (nachgemessen 2026-08-09). Die Wellen-**Plan**-Dateien sind nicht darunter —
   ihre Notiz ist H2. Ein zehnter Befund der Messung, ein `planning-drift`,
   entsteht **nur** beim Weiten von `slice-glob`; mit dem hier beantragten
   eigenen Schlüssel bleibt die Roadmap-Invariante unberührt, und genau das ist
   der Punkt.
   **Entschieden 2026-08-10, in zwei Schritten — der zweite hat den ersten
   umgekehrt.**

   Gemessen wurde zunächst eine **vierte** Variante, die erst mit dem neuen
   Schlüssel möglich ist: `closure.glob: "*.md"` **und** `heading-pattern` auf
   `^#{1,3}`.

   | Variante | Befunde | geprüfte Dokumente |
   |---|---|---|
   | heute (Filter erbt `slice-glob`) | 0 | 95 |
   | Filter `*.md` | 12 | +15 |
   | Filter `*.md` **und** Muster `^#{1,3}` | 0 | +15 |

   Elf der zwölf sind kein fehlender Abschnitt, sondern eine **H1** gegen ein
   Muster, das H2/H3 verlangt. Der zwölfte war **echt**: ein geschlossenes
   Wellendokument trug in `done/` noch „_Ausstehend._". Ein zweites entging dem
   Befund nur, weil sein „_Ausstehend._" wortreich genug für die Satz-Schwelle
   war. **Beide sind gefüllt** — der Wert der Messung steht unabhängig davon,
   wie die Konfigurations-Frage ausgeht.

   Die dritte Zeile sah nach der besten Antwort aus. Ein unabhängiger Review hat
   gemessen, dass sie ein **Falsch-Negativ** baut: bei den Ergebnisnotizen ist
   die auf `^#{1,3} .*Closure-Notiz` passende Überschrift der Dokument-**Titel**,
   der gemessene Abschnitt also die ganze Datei.

   | Fixture | Ergebnis |
   |---|---|
   | drei `_Ausstehend._`-Abschnitte unter H1-Titel | **kein Befund** |
   | derselbe Platzhalter unter H2 | `closure-note-thin` |

   Damit hätte die Weitung genau die Defekt-Klasse unsichtbar gemacht, für die
   das Gate gebaut ist — und dieselbe, die die Messung eine Stunde zuvor gefunden
   hatte. Zugleich liefe die Floskel-Prüfung über beliebigen Dateitext.

   **Endstand: der eigene Bestand bleibt beim geerbten Filter.** Die Wurzel ist
   eine Artefakt-Verwechslung: eine Wellen-Ergebnisnotiz *enthält* keine
   Closure-Notiz, sie **ist** eine. Sie zu prüfen heißt, ihre Abschnitts-Struktur
   zu prüfen — die Zusage des Moduls `structure`
   ([slice-099](slice-099-structure-modul.md)), nicht die dieser
   Fähigkeit. Ein Kandidaten-Filter kann eine Datei in die Menge holen, aber
   nicht die Frage ändern, die an sie gestellt wird.

   Die 15 Wellen-Dokumente bleiben damit ungeprüft — **benannt**, nicht
   übersehen. Das ist die Grenze, die das Register als **BEO-004** führt, und sie
   ist hier bewusst offen statt falsch geschlossen.

## 4. Definition of Done

- [x] `planning.closure.glob` in Lastenheft (0.53.0), Spezifikation (Schritt C2)
      und Config-Schema; Default = **Verweis** auf `planning.slice-glob`, kein
      kopiertes Literal. Begründung in
      [ADR-0051](../../adr/0051-eigener-kandidaten-filter-closure.md) `Proposed`;
      [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) §Geschichte
      um den Verfeinerungs-Zeiger ergänzt (ihr Rumpf ist `Accepted` und damit
      immutabel).
- [x] Akzeptanzkriterien des CR als Tests, alle drei plus der Default-Verweis.
      **Fünf** Rückbauten geprüft, alle rot: Filter zurück auf `slice-glob`;
      Default als Literal statt Verweis; leerer Glob fällt still zurück;
      Glob-Validierung entfernt; Wert nicht durchgereicht. End-to-End gegen das
      gebaute Image: eigener Lauf 326 Dateien / 0 Befunde, expliziter leerer
      Glob bricht mit **Exit 2** und einer Meldung ab, die den Schlüssel nennt.
- [x] `make gates` grün; Abnahme-Punkt 2 beantwortet. `.d-check.closure.yml`
      bleibt beim geerbten Filter — die Weitung war der erste Anlauf und ist nach
      dem Review zurückgenommen (§3, Abnahme-Punkt 2).
- [x] **Unabhängiger Review** (Frischkontext) — 0 HIGH, 2 MEDIUM, 6 LOW, 1 INFO;
      merge-blockierend. Byte-Identität, Entkopplung, Config-Rand und
      SemVer-Einordnung sind **belegt** (Alt-Image gegen HEAD-Image auf
      identischem Baum, je byte-gleich auf stdout und stderr). Blockierend waren
      zwei Ränder: das Falsch-Negativ der eigenen Konfiguration und ein
      **doppeltes** Akzeptanzkriterium zur Nullmengen-Härte, dessen ältere
      Fassung noch `planning.slice-glob` nannte und gegen die Umsetzung
      falsifizierbar war. Beide behoben, dazu die LOW-Ränder (Bestandszahlen
      95/110 statt 96/111, YAML-`null` als abwesend ausgewiesen, §4-Zeile und
      `--doctor`-Klartext sagen „Kandidat" statt „Slice", Glob in der
      Nullmengen-Meldung gequotet) und ein **sechster** Rückbau, den die Tests
      nicht fingen: den Config-Rand an `closure.dir` zu koppeln blieb grün,
      obwohl er auch bei inerter Fähigkeit greift.
- [x] **Bestätigende Re-Review** (Frischkontext) — beide Erst-MEDIUM als geheilt
      bestätigt; blockierend war, was die Heilung neu aufmachte: die
      Byte-Identitäts-Zusage. `strconv.Quote` in der Nullmengen-`message` und der
      geänderte `--doctor`-Klartext liegen auf Pfaden, die **ohne** den neuen
      Schlüssel erreichbar sind. Präzisiert statt zurückgenommen: die Zusage gilt
      dem **Befundsatz** (Datei, Zeile, Regel, Ziel, Grund), nicht den
      Begleit-Texten — beide sind laut Spezifikation nicht
      stabilitätsgarantiert, gehören aber in die Release-Notiz und stehen jetzt
      dort. Dazu fünf Ränder, zwei davon **gegen** den Reviewer nachgemessen und
      zwei **für** ihn:
      - Das Satzende-Minimum steht als **7** in der Config, nicht als 5 — mit dem
        Produkt gemessen (bei `min-sentences: 8` wird der erste Slice rot). Meine
        eigene Nachrechnung sagte 5 und war falsch.
      - Die Alternativen-Begründung „`path.Match` kennt keine Suffix-Negation"
        war **sachlich falsch** (`[^…]` existiert). Der Reviewer hat ein
        funktionierendes Muster gezeigt; die Nachmessung zeigt aber, dass es
        still **2 von 95** Slices verliert, deren drittes Namenssegment mit `r`
        beginnt. Die Zeile nennt jetzt diesen Grund statt eines falschen.
      - Bestandszahl in der ADR-Tabelle (96 → 95), §4-Zeile um den Leerlauf-Fall
        ergänzt, `--print-config`-Vorlage bildet den Verweis-Default nicht mehr
        als Literal ab.
      - **Siebter** Rückbau: auch `heading-pattern` und `boilerplate` greifen bei
        inerter Fähigkeit; das war ungetestet. Jetzt rot.
- [ ] **Release** (Minor: neuer Config-Schlüssel), Digest-Backfill. Ohne
      veröffentlichte Version erreicht die Welle ihren Zweck nicht — der
      Konsument kann sein Skript erst gegen ein Release zurückziehen
      (Schnitt-Review F-6).

## 5. Risiken / offene Punkte

- **Neun eigene Befunde bei geweitetem Glob** (siehe Abnahme-Punkt 2).
  — **Ausgang: entfallen.** Es waren inzwischen zwölf, aber elf davon lösen sich
  mit dem Überschriften-Muster auf; der zwölfte war ein echter Rückstand und ist
  behoben. Eine Sanierung der Wellen-Notizen war nicht nötig.
- **Zwei Globs, die fast immer gleich sind,** laden zu Drift ein (einer wird
  gepflegt, der andere nicht). — **Ausgang: adressiert und geprüft.** Der Default
  ist ein **Verweis**, kein Literal: wer nichts trennt, pflegt genau ein Muster.
  Die Gegenprobe „Default als Literal statt Verweis" ist rot — die Eigenschaft
  ist testgehalten, nicht bloß zugesagt. Es bleibt der Restrisiko-Fall, dass
  jemand `slice-glob` ändert und die Closure-Menge unbeabsichtigt mitwandert;
  das steht als Re-Evaluierungs-Trigger in
  [ADR-0051](../../adr/0051-eigener-kandidaten-filter-closure.md).

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei — **und**
[slice-096](../done/slice-096-structure-modul-analyse.md) in `done/`, weil dessen
Schnitt über den Zuschnitt dieses Slice mitentscheidet.

**Rückführungen:** `in-progress` → `open`, falls Abnahme-Punkt 2 eine
Bestands-Sanierung nach sich zieht, die eigenständig geschnitten gehört.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten** (bei Slice-Beginn erneut gelesen — das
  Register führt inzwischen **vier** Einträge, nicht mehr nur BEO-001):
  - **BEO-001** (Datei-Register driften gegen ihre Autoritäts-Tabelle): andere
    Klasse, nichts zu berücksichtigen.
  - **BEO-003** (geteilte Lexik driftet an den Rändern, weil jeder Konsument sie
    selbst vorbereitet): **einschlägig als Warnung**. Ein kopierter
    Literal-Default wäre genau diese Klasse gewesen — deshalb der Verweis.
  - **BEO-004** (Modul-Grenze nur auf der Quell-Achse gedacht): **einschlägig
    als Frage.** Die Closure-Fähigkeit ist ein Post-Pass über ein selbst
    benanntes Verzeichnis; dieser Slice ändert genau, welche Dateien darin sie
    liest. Die Register-Frage „welche Eingaben liest dieses Modul, die es nicht
    scannt?" hat Abnahme-Punkt 2 beantwortet — und dabei 15 übersehene
    Dokumente gefunden.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Change Request und Spezifikations-Schärfung
schreiben die Zusage, der Go-Code liefert sie.

## 9. Closure-Notiz (nach `done/`)

Geliefert ist `planning.closure.glob` als eigener Kandidaten-Filter der
Closure-Fähigkeit, ausgeliefert mit **v0.54.0**
([`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
Lastenheft 0.53.1, [ADR-0051](../../adr/0051-eigener-kandidaten-filter-closure.md)
`Accepted`). Der Default ist ein **Verweis** auf `planning.slice-glob`: ohne den
Schlüssel bleibt der Befundsatz byte-identisch, und es gibt genau ein Muster zu
pflegen.

**Der Anlass war ein Entwurfsfehler, nicht eine fehlende Fähigkeit.**
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) hatte gegen
ein zweites Modul argumentiert, es hätte „dieselbe Config-Achse … ein zweites Mal
deklariert". Die Achse war nie dieselbe: die eine Fähigkeit zählt, was **noch in
Arbeit** ist, die andere prüft, was **abgeschlossen** ist. Solange beide zufällig
dasselbe trafen, fiel es nicht auf. Das Ergebnis jener Entscheidung trägt
weiterhin — es ruht auf der geteilten Lifecycle-**Invariante**, nicht auf der
Config-Achse.

**Die Bestandsmessung hat einen realen Rückstand gefunden, bevor der Schlüssel
existierte.** Im Ruheort liegen 110 Markdown-Dateien, gesehen wurden 95. Von den
15 ungesehenen trug eine — ein von mir eine Stunde zuvor geschlossenes
Wellendokument — in ihrer Closure-Notiz noch „_Ausstehend._"; eine zweite entging
dem Befund nur, weil ihr Platzhalter wortreich genug für die Satz-Schwelle war.
Beide sind gefüllt. Der Wert der Messung steht unabhängig davon, wie die
Konfigurations-Frage ausging.

**Und sie ging anders aus, als ich zuerst entschied.** Ich hatte den eigenen
Bestand mitgeweitet — Glob auf `*.md`, Überschriften-Muster auf `^#{1,3}`, null
Befunde bei 15 zusätzlich erfassten Dokumenten. Der Review hat gemessen, dass das
ein **Falsch-Negativ** baut: bei den Ergebnisnotizen ist die passende Überschrift
der Dokument-**Titel**, der gemessene Abschnitt also die ganze Datei. Ein
Dokument mit drei unausgefüllten Abschnitten bleibt damit grün, während derselbe
Platzhalter unter H2 meldet. Die Weitung hätte genau die Defekt-Klasse unsichtbar
gemacht, für die das Gate gebaut ist — und dieselbe, die die Messung eben gefunden
hatte.

Die Wurzel war eine **Artefakt-Verwechslung**: eine Wellen-Ergebnisnotiz enthält
keine Closure-Notiz, sie **ist** eine. Ein Kandidaten-Filter kann eine Datei in
die Menge holen, aber nicht die Frage ändern, die an sie gestellt wird. Die 15
Wellen-Dokumente bleiben ungeprüft — **benannt**, nicht übersehen; ihre Abdeckung
gehört zum Modul `structure`.

**Zwei Review-Runden, sieben Rückbauten.** Die zweite Runde blockierte auf dem,
was die erste Heilung neu aufmachte: `strconv.Quote` und der geänderte
`--doctor`-Klartext liegen auf Pfaden, die man auch **ohne** den neuen Schlüssel
erreicht, und die ADR sagte pauschal „byte-identisch". Präzisiert statt
zurückgenommen — die Zusage gilt dem **Befundsatz**, die beiden Begleit-Texte
stehen in der Release-Notiz. Zurücknehmen wäre falsch gewesen: der Klartext
*muss* „Kandidat" sagen, seit die Kandidaten-Menge nicht mehr aus Slice-Dateien
bestehen muss.

**Zwei Ränder habe ich gegen den Reviewer nachgemessen — beide gingen an ihn.**
Das Satzende-Minimum steht als **7** in der Config, nicht als 5: meine eigene
Nachrechnung war falsch, weil ihre Fence-Behandlung vom Produkt abwich; mit dem
Produkt gemessen (Schwelle hochdrehen, bis etwas rot wird) sind es 7. Und meine
Begründung „`path.Match` kennt keine Suffix-Negation" war schlicht falsch —
`[^…]` steht im eigenen Handbuch. Sein Gegenbeispiel funktioniert, verliert aber
still **2 von 95** Slices; die Alternativen-Zeile nennt jetzt diesen Grund statt
eines erfundenen.

**Fürs Register:** die Sichtung in §7 nannte bei Slice-Beginn nur BEO-001, das
Register führte da schon vier Einträge. **BEO-003** war die Warnung, die den
Literal-Default verhindert hat; **BEO-004** die Frage, die Abnahme-Punkt 2
beantwortet hat — und deren Antwort am Ende lautete: die Grenze bleibt offen,
aber sie ist benannt.
