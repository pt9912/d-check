# Slice slice-104: Floskel-Vergleich an der Wortgrenze statt als Teilstring

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Floskel-Bedingung der Closure-Note-Struktur),
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md)
(Entscheidung 5: Liste per Default leer).
**Change Request** des Auftraggebers, 2026-08-10.

**Autor:** pt9912. **Datum:** 2026-08-10.

---

## 1. Ziel

`planning.closure.boilerplate` vergleicht literale **Teilstrings**. Kurze Phrasen
sind dadurch unbrauchbar — und kurze Phrasen sind genau die, für die die
Bedingung gedacht ist. Der Vergleich soll an **Wortgrenzen** stattfinden.

## 2. Der gemessene Befund

Nachgestellt am 2026-08-10 gegen v0.55.0, `boilerplate: ["ok"]`, mit einer Notiz,
deren einziges Vorkommen in *dokumentiert* liegt:

```text
"Der Ablauf ist dokumentiert."  → closure-note-boilerplate, Exit 1
"Der Ablauf ist beschrieben."   → 0 Befunde, Exit 0
```

**Gemessen über die 96 eigenen Closure-Notizen** (Wortgrenze in RE2-Semantik,
`\w` = `[0-9A-Za-z_]`):

| Phrase | Teilstring | Wortgrenze |
|---|---|---|
| `ok` | **68** | **1** |
| `n/a` | 2 | 0 |
| `fertig` | 3 | 0 |
| `gut` | 2 | 1 |
| `läuft jetzt` | 3 | 3 |
| die fünf **aktuell** konfigurierten Phrasen | 0 | 0 |

Drei Aussagen folgen daraus:

1. **Der CR trifft.** Bei `ok` sind 67 von 68 Treffern Falsch-Positive.
2. **Mehrwortige Phrasen sind verhaltensgleich** (`läuft jetzt`: 3 = 3), und die
   fünf heute konfigurierten Phrasen ändern sich **nicht** — für den eigenen
   Lauf ist die Änderung byte-identisch.
3. **Es macht verworfene Phrasen brauchbar.** `fertig` und `läuft jetzt` wurden
   in [slice-093](../done/slice-093-closure-note-gate.md) ausdrücklich **nicht**
   aufgenommen, weil sie als Teilstring Falschbefunde erzeugten. `fertig` fällt
   mit Wortgrenzen auf 0.

**Wortgrenzen machen kurze Phrasen brauchbar, nicht automatisch sicher.** Der
eine verbleibende `ok`-Treffer ist echt: „image-test ok, benchmark-median
526 ms“ — eine substanzhaltige Notiz. Wer `ok` aufnimmt, kauft weiterhin
diesen Fall mit.

## 3. Abnahme-Punkte

1. **RE2-Portierung.** Der Vorschlag nennt Lookbehind und Lookahead; Go kennt
   sie nicht. Dieselbe Lage wie bei der Platzhalter-Erkennung in
   [slice-098](../done/slice-098-closure-note-placeholder.md), und
   dieselbe Technik: das Grenzzeichen **konsumieren** statt hineinzuschauen.
   Zu entscheiden ist, ob das über eine gebaute Regex je Phrase läuft
   (`regexp.QuoteMeta` + Zeichenklassen) oder über eine Index-Suche mit
   Nachbar-Prüfung — Letzteres kommt ohne Regex-Kompilierung je Phrase aus.
2. **Was ist eine Wortgrenze?** RE2s `\w` ist **ASCII-only**: Umlaute und
   Anführungszeichen zählen als Nicht-Wortzeichen. Für Deutsch hilft das meist
   (gemessen: die beiden `gut`-Treffer stehen hinter `„` und sind echte
   Treffer), aber eine Phrase **innerhalb** eines Wortes mit Umlaut daneben
   würde fälschlich als grenzständig gelten. Am eigenen Bestand gemessen: **kein
   solcher Fall**. Zu entscheiden: ASCII-Grenze übernehmen und die Lage benennen,
   oder auf Unicode-Buchstaben prüfen.
3. **SemVer und Richtung.** Die Änderung findet **weniger** — ein roter
   Konsumentenlauf kann grün werden. Das ist die andere Richtung als bei den
   bisherigen Lexik-Angleichungen und gehört ausdrücklich in die Release-Notiz:
   wer heute auf Teilstring-Treffer baut, verliert sie. Zu entscheiden, ob das
   Minor bleibt.
4. **Eigene Konfiguration.** Welche kurzen Phrasen nehmen wir danach auf?
   Erst messen, dann aufnehmen — dieselbe Disziplin wie in slice-093 und 098.
   `fertig` ist der offensichtliche Kandidat (3 → 0).

## 4. Definition of Done

- [x] Vertragsanpassung in Lastenheft (0.56.0) und Spezifikation (Schritt C4).
      **Keine eigene ADR:** Abnahme-Punkt 2 hat keine eigene Entscheidung
      erzeugt — die ASCII-Grenze folgt der Messung, und den Rahmen für die
      Lockerung derselben Prüfung setzt bereits
      [ADR-0053](../../adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md)
      aus demselben Wellen-Zug.
- [x] Alle vier Test-Zusagen erfüllt, dazu ein direkter Test auf die
      Grenz-Funktion für die Ränder, die über das Modul schwer zu stellen sind
      (leere Phrase, Treffer am Textrand, mehrere Vorkommen, `okok`).
- [x] Mutations-Gegenprobe über eine Dateikopie: **sieben** Rückbauten, alle rot
      (zurück auf Teilstring; nur eine der beiden Grenzen geprüft, je einzeln;
      Ziffern bzw. Unterstrich als Nicht-Wortzeichen; Weitersuche nach einem
      verworfenen Treffer; Textrand als Wortzeichen).
- [x] `make gates` + `make verify-closure-notes` grün; die eigene Liste **nach
      Messung** ergänzt. Je Kandidat ein Lauf gegen den eigenen Bestand:

      | Phrase | Treffer | aufgenommen |
      |---|---|---|
      | `fertig` | 0 | **ja** |
      | `n/a` | 0 | **ja** |
      | `erledigt` | 1 | nein — der Treffer ist echte Prosa |
      | `ok` | 1 | nein — „image-test ok“ ist substanzhaltig |
      | `läuft` | 11 | nein |

      Damit sind genau die zwei Phrasen aufgenommen, die als Teilstring
      **verworfen werden mussten**. Die Messung bestätigt auch die Warnung aus
      §2: Wortgrenzen machen kurze Phrasen brauchbar, nicht automatisch sicher.
- [x] **Unabhängiger Review über beide Wellen-Slices** (Frischkontext) —
      1 HIGH, 4 MEDIUM, 4 LOW, 1 INFO; merge-blockierend. Der Paritäts-Beleg
      (84/84), der Abfall 7 → 5 und die 170 fett gesetzten Satzenden sind mit
      **eigener** Methode nachgemessen und bestätigt; die Richtung ist monoton
      (kein Wert steigt). Alle Befunde geheilt, siehe §5.
- [ ] **Release** — **Wellen-Trigger**, gemeinsam mit
      [slice-094](../done/slice-094-closure-zaehl-paritaet.md): **eine** Notiz für
      beide Änderungen, mit **jeder** Richtung einzeln.

## 5. Risiken / offene Punkte

- **Ein ausgeliefertes Gate ändert sein Verhalten.** Anders als bei den letzten
  beiden Slices ist das nicht additiv. Die Richtung ist zwar die harmlosere
  (weniger Falsch-Positive), aber ein Konsument, dessen Notiz heute rot ist,
  bekommt sie ohne Zutun grün. — **Ausgang:** offen; Abnahme-Punkt 3.
- **Die ASCII-Grenze** könnte in einem anderssprachigen Repo anders wirken als
  hier gemessen. — **Ausgang: benannt, nicht geschlossen.** Der Vertrag nennt
  sie jetzt in **beiden** Richtungen: ein Umlaut ist kein Wortzeichen, eine
  Phrase mit angrenzendem Umlaut gilt damit als grenzständig und **trifft**. Am
  eigenen Bestand kommt der Fall nicht vor.
- **CRLF war der schwerste Befund des Reviews** (HIGH) und im eigenen Repo
  unsichtbar: `countSentenceEnds` akzeptierte nur Space, Tab und Zeilenumbruch,
  nicht das `\r` davor. In einer CRLF-Arbeitskopie zählte damit **kein**
  zeilenschließendes Satzende — eine Notiz mit vier sauberen Sätzen meldete
  `closure-note-thin`, dieselbe mit LF blieb grün. Das Adopter-Skript zählt `\r`
  als Whitespace, die Parität wäre also auch gebrochen gewesen.
  — **Ausgang: behoben**, mit Testfall über Modul und Funktion. **Zweite
  CRLF-Regression an einem Tag** — die erste steckte in slice-101.

## 6. Trigger

**Start** (`open` → `next` → `in-progress`): Freigabe; WIP-Slot frei.

**Welle:** [welle-72-closure-semantik](../done/welle-72-closure-semantik.md), gemeinsam
mit [slice-094](../done/slice-094-closure-zaehl-paritaet.md); dort **nach**
094, das die ADR für die Lockerung trägt. Beide ändern die
**Semantik eines ausgelieferten Gates** statt etwas hinzuzufügen, beide brauchen
eine eigene Bestandsmessung und eine konsumentensichtbare Release-Notiz. Sie
zusammen zu schneiden bündelt genau eine Risiko-Klasse — dieselbe Begründung,
mit der slice-094 aus welle-71 herausgehalten wurde.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** bei der Planung erneut lesen. Voraussichtlich
  einschlägig ist die Klasse „geteilte Lexik driftet an den Rändern“ — ein
  Wortgrenzen-Begriff, der nur hier lebt, wäre genau das.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — die geänderte Zusage wird zuerst formuliert
(was heißt Treffer?), dann geliefert.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
