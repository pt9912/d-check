# Slice slice-098: `closure-note-placeholder` — unausgefüllte Template-Platzhalter erkennen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** welle-71, gemeinsam mit [slice-097](../done/slice-097-closure-glob-entkopplung.md)
(Kennung nachgetragen bei [slice-191](../in-progress/slice-191-alt-bestand-archivieren.md)
— das Feld nannte bislang nur die Zuordnungs-Bedingung, nicht die Kennung) —
erst mit **beiden** ist die Closure-Fähigkeit eine Obermenge des
Konsumenten-Skripts; einzeln bringt jede Deckung, löst den Faden aber nicht auf.
**Zuordnung entschieden** mit der welle-69-Closure (2026-08-09): der Slice
bleibt **eigenständig** und geht nicht im `structure`-Modul auf.

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(neuer Grund-Code in der C-Kette),
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md).
**Change Request** des Konsumenten `ai-harness-course` (CR 2).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Ein neuer Grund-Code `closure-note-placeholder` (opt-in über
`planning.closure.placeholder`, bool, Default `false`) meldet Closure-Notizen,
die noch den **unausgefüllten Rumpf** eines Templates tragen.

## 2. Der gemessene Befund

Ein Template-Rumpf ist syntaktisch vollständig und passiert **alle drei**
bestehenden Prüfungen. Nachgemessen am 2026-08-09 gegen v0.52.0 mit der Notiz
`Ergebnis: <ergebnis>. Belege: <belege>. Offen: <offen>. Ende: <ende>.` —
**0 Befunde**: der Abschnitt existiert, vier Satzende-Zeichen erreichen die
Schwelle, und keine Phrase der Floskel-Liste kommt vor.

Das ist genau die Ausgangslage, die das Handbuch für §4.17 beschreibt
(„niemand merkt, wenn eine davon ein zurückgelassenes `_Ausstehend._` ist") —
und sie bleibt grün.

## 3. Der heikle Teil: die Erkennung ist nicht portierbar

Eine naive Fassung (`<` … `>`) ist unbrauchbar: Closure-Notizen tragen
QA-Messungen wie `p95 < 1 s` und `Recall > 0,9`. Der Konsument liefert eine an
drei realen Falsch-Positiv-Klassen erprobte Regex mit **Lookbehind und
Lookahead** — und **die kann d-check nicht übernehmen**. Go nutzt RE2; das
Muster wird fail-closed abgelehnt:

```text
d-check: error: .d-check.yml: planning.closure.heading-pattern "(?<![\w/])<(?![\s!/])…"
ist kein gültiges Regex: error parsing regexp: invalid named capture
```

Das ist kein Blocker, aber eine **Portierung mit eigener Testlast**, keine
Kopie. Beide Lookarounds sind feste, negative Zeichen-Prüfungen und damit in
RE2 ausdrückbar, indem man das Zeichen **konsumiert** statt hineinzuschauen:

- `(?<![\w/])` → das Vorzeichen mitmatchen (`(^|[^\w/])`) und per Capture-Gruppe
  wieder abziehen;
- `(?![\s!/])` → als explizite Zeichenklasse am ersten Zeichen nach `<`;
- der Ein-Zeichen-Fall (`<x>`) fällt sonst durchs Raster und braucht eine
  eigene Alternative.

**Was NICHT in die Regex gehört:** die Nachfilter. Verwerfen, wenn der Inhalt
`://` oder `@` enthält (Autolinks) oder — klein und ohne Trailing-Slash — in
einer HTML-Tag-Liste steht. Das ist Code, kein Muster; in die Regex gepresst
wird es unlesbar und unprüfbar.

**Konsumierte Vorzeichen überlappen:** zwei benachbarte Platzhalter können sich
gegenseitig verdecken. Für den Vertrag ist das folgenlos (gemeldet wird der
**erste** Treffer, wie bei der Floskel), gehört aber in die Tests.

## 4. Definition of Done

- [x] Grund-Code + `planning.closure.placeholder` in Lastenheft (0.54.0),
      Spezifikation (Schritt **C4b**, Nachfilter ausdrücklich als Code benannt)
      und Config-Schema. **ADR nötig** —
      [ADR-0052](../../adr/0052-platzhalter-erkennung-inline-code.md) `Proposed`:
      die Erkennungs-Grenze trägt zwei eigene Entscheidungen (Inline-Code ist
      keine Prosa; die Substanz-Zählung bleibt bewusst unberührt).
- [x] Alle vier Akzeptanzkriterien als Tests, die drei Falsch-Positiv-Klassen
      **einzeln** (neun Unterfälle: Vergleichszeichen, Generic, Autolink,
      Mail-Adresse, HTML-Tag mit Attribut, selbstschließendes Tag, in
      Inline-Code gezeigte Meta-Syntax, HTML-Kommentar, schließendes Tag).
      Dazu Randformen (Ein-Zeichen-Platzhalter, verworfener Treffer direkt vor
      einem echten) und die Abschnitts-Grenze.
      **Acht Rückbauten geprüft.** Zwei blieben zunächst grün und haben Arbeit
      erzeugt: der Config-Schalter war ungetestet (jetzt rot), und der
      „Überlappungs-Fix" stellte sich als **äquivalenter Mutant** heraus — durch
      das Re-Slicing greift die `^`-Alternative am neuen Anfang, beide Varianten
      sind gleich. Der Kommentar behauptete ein Problem, das nicht besteht, und
      ist korrigiert; die Zusage, dass eine Verwerfung die Suche nicht beendet,
      ist stattdessen direkt bewacht.
- [x] **Unabhängiger Review** (Frischkontext) — 0 HIGH, 5 MEDIUM, 5 LOW, 2 INFO;
      merge-blockierend. Byte-Identität, Abschnitts-Grenze und SemVer-Einordnung
      **belegt** (Alt-Image gegen HEAD, Klartext/`--json`/`--doctor`). Der
      schwerste Befund traf eine Zusage, die ich selbst geschrieben hatte:
      „Vergleichszeichen sind durch die Form ausgeschlossen" hielt nur für die
      Schreibweise **mit** Leerzeichen — `<1 s und der Recall >0,9` meldete.
      Behoben durch Verengung (Inneres frei von Whitespace), was zugleich
      HTML-Tags mit Attributen erledigt; dazu ein dritter Nachfilter für
      Winkelklammer-Linkziele. Ferner: §C5 der Spezifikation war auf der Fassung
      **vor** C4b stehengeblieben und widersprach ihr, der Substanz-Bullet des
      Lastenhefts behauptete weiter, ein Platzhalter falle der Zählung auf, und
      **33 von 35** Einträgen der HTML-Tag-Liste ließen sich löschen, ohne dass
      ein Test rot wurde (jetzt über die Liste iteriert). Zwei Grenzen sind
      **benannt statt geschlossen**: der eingerückte Code-Block und die ungerade
      Backtick-Parität — beides Eigenschaften der geteilten Lexik.
- [x] `make gates` + `make verify-closure-notes` grün. **Der eigene Bestand ist
      gemessen und die Erkennung scharfgeschaltet** (`placeholder: true` in
      `.d-check.closure.yml`): 0 Befunde über 96 Closure-Notizen. End-to-End
      gegen das gebaute Image: der Template-Rumpf meldet an der Zeile seines
      ersten Treffers, ohne den Schalter verschwindet der Befund.
- [ ] **Release** (Minor: neuer Grund-Code), Digest-Backfill. Ohne
      veröffentlichte Version erreicht die Welle ihren Zweck nicht
      (Schnitt-Review F-6).

## 5. Risiken / offene Punkte

- **Die RE2-Portierung könnte die erprobte Semantik verfehlen.** — **Ausgang:
  belegt.** Neun Falsch-Positiv-Unterfälle einzeln getestet, dazu die Messung am
  eigenen Bestand. Die Portierung ersetzt beide Lookarounds durch Konsumieren
  bzw. eine Zeichenklasse; die Nachfilter liegen im Code, nicht im Muster.
- **Winkelklammern sind in technischer Prosa häufig.** Ein Falsch-Positiv in
  einer Closure-Notiz ist teurer als ein übersehener Platzhalter, weil es das
  Gate unglaubwürdig macht. — **Ausgang: entschärft.** Alle zwölf Treffer der
  portierten Fassung lagen in Inline-Code; die Einschränkung darauf bringt den
  eigenen Bestand auf null. Default `false` bleibt trotzdem — die Form ist
  schreibkultur-abhängig.
- **Der eigene Bestand ist ungemessen.** — **Ausgang: gemessen, dann
  scharfgeschaltet.** Drei Fassungen über die 96 Closure-Notizen: naiv 24
  Treffer, portiert 12, portiert ohne Inline-Code **0** — und **keiner** der 24
  war ein echter Platzhalter. Die Messung hat die Erkennung entschieden, nicht
  bloß bestätigt.

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei — **und**
[slice-096](welle-69/slice-096-structure-modul-analyse.md) in `done/`. Von
[slice-097](../done/slice-097-closure-glob-entkopplung.md) unabhängig: die gemeinsame
Abnahmebedingung des Konsumenten verlangt beide, aber keine Reihenfolge.

**Rückführungen:** `in-progress` → `open`, falls die Messung am eigenen Bestand
zeigt, dass die Erkennung ohne ein zusätzliches Ventil nicht falsch-positiv-frei
ist.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten** (bei Slice-Beginn erneut gelesen — das
  Register führt inzwischen **vier** Einträge, nicht mehr nur BEO-001):
  - **BEO-001** (Datei-Register driften gegen ihre Autoritäts-Tabelle): andere
    Klasse, nichts zu berücksichtigen.
  - **BEO-002** (Semantik-Änderung nur im Körper, Ränder bleiben stehen):
    einschlägig als **Arbeitsregel** — der neue Grund-Code hat Ränder in §4, im
    `--doctor`-Klartext, im Handbuch und in der `print-config`-Vorlage.
  - **BEO-003** (geteilte Lexik driftet an den Rändern): **einschlägig.** Die
    Inline-Code-Paarung ist geteilte Lexik; diese Bedingung wird sie benutzen,
    nicht nachbauen.
  - **BEO-004** (Modul-Grenze nur auf der Quell-Achse): hier ohne Wirkung — die
    Bedingung liest keine zusätzliche Eingabe, sie sieht denselben Abschnitt
    enger.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — der Vertrag (Grund-Code, Nachfilter,
Akzeptanzkriterien) wird zuerst geschrieben, die Erkennung liefert ihn. Kein
Brownfield: die Konsumenten-Regex ist **Vorlage**, nicht rückzudokumentierender
Bestand — sie wird portiert und neu belegt, nicht übernommen.

## 9. Closure-Notiz (nach `done/`)

Geliefert ist der Grund-Code `closure-note-placeholder` als vierte, **opt-in**
Bedingung der Closure-Note-Struktur, ausgeliefert mit **v0.55.0**
([`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
Lastenheft 0.54.1, [ADR-0052](../../adr/0052-platzhalter-erkennung-inline-code.md)
`Accepted`). Der eigene Bestand ist gemessen und die Bedingung scharfgeschaltet:
96 Closure-Notizen, null Befunde.

**Die Messung hat die Erkennung entschieden, nicht bloß bestätigt.** Drei
Fassungen über denselben Bestand: naiv 24 Treffer, portiert 12, portiert ohne
Inline-Code **null** — und **keiner** der 24 war ein echter Platzhalter. Alle
zwölf Treffer der portierten Fassung lagen in Inline-Code. Daraus wurde die
tragende Entscheidung: Inline-Code ist keine Prosa, weil dort Syntax **gezeigt**
wird. Ohne sie wäre das Gate am ersten Tag unglaubwürdig gewesen.

**Eine Messung belegt, was da ist, nicht was möglich ist.** Das war die
teuerste Lehre dieses Slice. Ich hatte in den Vertrag geschrieben,
Vergleichszeichen seien „bereits durch die Form ausgeschlossen“ — und das
stimmte nur für `p95 < 1 s` mit Leerzeichen. Der Review fand die enge
Schreibweise: `Die Latenz blieb <1 s und der Recall >0,9` meldete, weil der
Treffer vom ersten Winkel bis zum späteren spannt. Mein eigener Bestand kennt
diese Form nicht, also hat meine Messung sie nicht gefunden.

Geheilt wurde durch **Verengung** statt durch einen weiteren Nachfilter: das
Innere muss frei von Whitespace sein — ein Platzhalter ist ein *Feldname*, kein
Satz. Das erledigt zugleich HTML-Tags mit Attributen und macht die
Nachfilter-Last kleiner statt größer. Mehrwortige Formen sind damit ein
Re-Evaluierungs-Trigger, keine Erweiterung des Musters.

**Zwei Zusagen waren praktisch unbewacht.** Von 35 Einträgen der HTML-Tag-Liste
ließen sich **33** löschen, ohne dass ein Test rot wurde; dabei fiel auf, dass
`<section>` gar nicht darin stand und meldete. Jetzt iteriert ein Test über die
Liste. Und `closureSectionEnd` war eine **Kopie** der Grenz-Logik aus
`closureSectionProse` — genau die Klasse, die der Slice-Plan selbst als
einschlägig erklärt hatte. Jetzt ruft letztere erstere: eine Grenze, eine
Implementierung.

**Ein Test hat einen echten Defekt gefunden.** `FindAllStringSubmatch` ist
nicht-überlappend: ein vom Nachfilter verworfener Treffer beendete die Suche für
die restliche Zeile. Ersetzt durch eine Schleife auf dem Reststring. Von den
acht Rückbauten blieb außerdem einer grün, der sich als **äquivalenter Mutant**
herausstellte — kein Testfehler, sondern ein Kommentar, der ein Problem
behauptete, das nicht besteht.

**Zwei Grenzen sind benannt statt geschlossen:** ein mit vier Leerzeichen
eingerückter Code-Block ist in d-check **nirgends** modelliert, und eine
ungerade Backtick-Zahl im Absatz verschiebt die Inline-Code-Paarung. Beides sind
Eigenschaften der geteilten Lexik; sie hier allein zu reparieren hieße, eine
sechste Sicht auf denselben Text zu bauen. Sie stehen im **Handbuch**, nicht nur
in der ADR — wer sie nicht kennt, hält den ersten Befund für einen Fehlalarm.

**Ausdrücklich nicht mitgeändert:** die Substanz-Zählung sieht Inline-Code
weiterhin. Sie zu verengen bewegt eine ausgelieferte Schwelle in beide
Richtungen und ist eine eigene Entscheidung mit eigener Messung
([slice-094](../done/slice-094-closure-zaehl-paritaet.md)).
