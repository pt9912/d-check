# Slice slice-098: `closure-note-placeholder` — unausgefüllte Template-Platzhalter erkennen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** gemeinsam mit [slice-097](slice-097-closure-glob-entkopplung.md) —
erst mit **beiden** ist die Closure-Fähigkeit eine Obermenge des
Konsumenten-Skripts; einzeln bringt jede Deckung, löst den Faden aber nicht auf.
**Die Zuordnung ist offen bis zum Abschluss von**
[slice-096](../done/slice-096-structure-modul-analyse.md) — dessen Schnitt entscheidet,
ob die Erkennung hier oder im `structure`-Modul wohnt (Auftraggeber-Entscheid
2026-08-09).

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

- [ ] Grund-Code + `planning.closure.placeholder` in Lastenheft, Spezifikation
      (C-Kette, Nachfilter benannt) und Config-Schema; begleitende ADR nur, falls
      die Erkennungs-Grenze eine eigene Entscheidung braucht.
- [ ] Akzeptanzkriterien des CR als Tests: (1) ohne den Schlüssel
      byte-identisch; (2) die vier Platzhalter-Sätze ⇒ **genau ein** Befund
      (erster Treffer); (3) eine Notiz mit `p95 < 1 s`, `Recall > 0,9`,
      `vector<float>`, einem Autolink und einem HTML-Tag ⇒ **kein** Befund;
      (4) kombinierbar mit `closure-note-thin`.
- [ ] `make gates` + `make verify-closure-notes` grün; die drei
      Falsch-Positiv-Klassen einzeln als Test, nicht als Sammel-Fixture.
- [ ] **Release** (Minor: neuer Grund-Code), Digest-Backfill. Ohne
      veröffentlichte Version erreicht die Welle ihren Zweck nicht
      (Schnitt-Review F-6).

## 5. Risiken / offene Punkte

- **Die RE2-Portierung könnte die erprobte Semantik verfehlen** — die Vorlage ist
  gegen reale Falsch-Positive gehärtet, die portierte Fassung ist es zunächst
  nicht. — **Ausgang:** offen; die drei Klassen werden einzeln getestet, nicht
  gebündelt.
- **Winkelklammern sind in technischer Prosa häufig.** Ein Falsch-Positiv in
  einer Closure-Notiz ist teurer als ein übersehener Platzhalter, weil es das
  Gate unglaubwürdig macht. — **Ausgang:** offen; Default `false` ist die
  Absicherung, bis der eigene Bestand gemessen ist.
- **Der eigene Bestand ist ungemessen.** Vor dem Scharfschalten in der eigenen
  Konfiguration zu prüfen — dieselbe Disziplin wie bei der Floskel-Liste in
  [slice-093](../done/slice-093-closure-note-gate.md). — **Ausgang:** offen.

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei — **und**
[slice-096](../done/slice-096-structure-modul-analyse.md) in `done/`. Von
[slice-097](slice-097-closure-glob-entkopplung.md) unabhängig: die gemeinsame
Abnahmebedingung des Konsumenten verlangt beide, aber keine Reihenfolge.

**Rückführungen:** `in-progress` → `open`, falls die Messung am eigenen Bestand
zeigt, dass die Erkennung ohne ein zusätzliches Ventil nicht falsch-positiv-frei
ist.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001**
  (Datei-Register driften unbemerkt gegen ihre Autoritäts-Tabelle). Andere
  Klasse — hier geht es um den Inhalt **eines** Abschnitts. Nichts zu
  berücksichtigen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — der Vertrag (Grund-Code, Nachfilter,
Akzeptanzkriterien) wird zuerst geschrieben, die Erkennung liefert ihn. Kein
Brownfield: die Konsumenten-Regex ist **Vorlage**, nicht rückzudokumentierender
Bestand — sie wird portiert und neu belegt, nicht übernommen.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
