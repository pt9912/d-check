# Slice slice-094: Zähl-Parität der Closure-Note-Struktur (Inline-Code + Satzende-Form)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** [welle-72-closure-semantik](../done/welle-72-closure-semantik.md), gemeinsam mit
[slice-104](../done/slice-104-floskel-wortgrenze.md). **Zuordnung entschieden** mit der
welle-69-Closure (2026-08-09): der Slice bleibt **eigenständig**. Die Welle klammert
beide, weil sie die Semantik eines **ausgelieferten** Gates ändern und **beide** die
Floskel-Prüfung lockern — das gehört in **eine** Release-Notiz.

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Closure-Note-Struktur, Schritt C4),
[ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md).
Anlass ist eine gemessene **Parität-Lücke** gegenüber dem Schwester-Repo a-check
(dessen `verify-closure-notes`-Skript soll durch dieses Modul ablösbar werden).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Die Substanz-Zählung der Closure-Note-Struktur so schärfen, dass sie zum
abzulösenden Adopter-Skript **paritätisch** ist: Inline-Code zählt nicht mit, und
ein Satzende-Zeichen zählt nur, wenn ihm Whitespace oder das Zeilenende folgt.

## 2. Der gemessene Befund

Gemessen am 2026-08-09 gegen eine Kopie des a-check-Bestands (76 Slices):

```text
## 6. Closure-Notiz

Siehe `a.md` und `b.md`.
```

Das Adopter-Skript entfernt Inline-Code **und** verlangt Whitespace nach dem
Satzzeichen ⇒ es zählt **1** und meldet bei Schwelle 2 rot. d-check v0.52.0
entfernt nur Fenced-Code und zählt jedes `.` ⇒ es zählt **3** und bleibt
**grün**. Die Abweichung liegt in der gefährlichen Richtung: eine Notiz, die der
Adopter-Sensor als zu dünn meldet, läuft bei d-check durch.

**SemVer: Minor, kein Patch.** d-check findet nach dieser Änderung **mehr** —
dieselbe Einordnung wie bei den Markdown-Lexik-Angleichungen
([ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md)) und der
Tabellengrenze ([ADR-0043](../../adr/0043-tabellengrenze-am-relevanten-header.md)):
ein grüner Konsumentenlauf kann danach rot werden.

**Die Änderung wirkt in ZWEI Richtungen — das war in der ersten Fassung nicht
gesehen (Schnitt-Review F-1).** Die Substanz-Zählung und die Floskel-Prüfung
lesen **denselben** bereinigten Abschnittstext. Wer Inline-Code entfernt,
verschärft nicht nur die Zählung, sondern **lockert** zugleich die
Floskel-Prüfung: eine Phrase in Backticks wird danach nicht mehr gefunden.
Empirisch belegt — heute rot, danach grün.

Sachlich ist das die **bessere** Semantik (eine *zitierte* Floskel ist keine
benutzte; die heutige Fassung meldet jede Notiz, die über Floskeln *schreibt*).
Aber es ist eine Lockerung einer Prüfregel, und
[`AGENTS.md` §3.6](../../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)
macht die ADR-pflichtig. **Auftraggeber-Entscheid 2026-08-09: ADR-Variante** —
die Lockerung wird bewusst angenommen und begründet, statt zwei getrennte
bereinigte Texte zu führen.

## 3. Abnahme-Punkte

1. **Reicht die Angleichung, oder braucht die Zählung ein Config-Ventil?** →
   **Entschieden 2026-08-10: ohne Ventil.** eine feste, an CommonMark orientierte Zählung
   (Inline-Spans sind Code, kein Fließtext) statt eines Schalters, den niemand
   bewusst setzt. Ein Ventil wäre eine zweite Semantik für dieselbe Frage.
2. **Wird der Bestand rot?** → **Gemessen 2026-08-10 mit dem Produkt: nein.**
   Die Schwelle wird von beiden Seiten hochgedreht, bis der erste Slice rot wird:

   | | Minimum im Bestand | bei unserer Schwelle 4 |
   |---|---|---|
   | v0.55.0 | 7 | grün (Abstand 3) |
   | nach der Angleichung | **5** | **grün** (Abstand 1) |

   Keine Sanierung nötig. Der **Abstand zur Schwelle schrumpft aber von 3 auf 1**
   — das ist die Wirkung, nicht ein Nebeneffekt, und gehört in die Release-Notiz.

3. **Neu aufgetaucht: was folgt einem Satzende?** Die Regel „nur vor
   Whitespace oder Zeilenende“ trifft eine Klasse, die der Slice nicht
   vorhergesehen hatte. Gemessen über die 97 eigenen Notizen, **nach**
   Inline-Code-Bereinigung, 4066 Satzende-Zeichen:

   | folgt darauf | Vorkommen | Bewertung |
   |---|---|---|
   | Whitespace / Zeilenende | 1320 | echte Sätze |
   | `.`, `/`, `m`, Ziffern | ~2400 | **Link-Pfade und Versionen** — nie Sätze; sie nicht mehr zu zählen ist der Kern der Angleichung |
   | `*` | **170** | `**Umsetzung.**` — ein **fett gesetztes** Satzende, also ein echter Satz, der nicht mehr zählt |

   Die 170 erklären den Abfall des Minimums von 7 auf 5 fast vollständig.
   **Entschieden: Whitespace-Regel wie spezifiziert, ohne Ausnahme für
   schließende Auszeichnung.** Gründe: es ist die zugesagte **Parität** (der
   ganze Zweck des Slice); die Richtung ist die sichere (zählt weniger ⇒ Gate
   strenger); und eine Ausnahmeliste für `*`, `_`, `` ` ``, `)`, `"` wäre eine
   **dritte** Semantik, die weder der Adopter noch CommonMark definiert. Der
   Preis ist benannt: ein Repo, das viele Sätze fett schließt und knapp über der
   Schwelle liegt, wird rot — genau der zugesagte Minor-Charakter.

## 4. Definition of Done

- [x] Vertrag geliefert: Lastenheft 0.55.0 (Substanz- und Floskel-Bullet, drei
      neue Akzeptanzkriterien), Spezifikation Schritt **C4** (eine Bereinigung
      für alle Bedingungen, Satzende nur vor Whitespace/Zeilenende) und
      [ADR-0053](../../adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md)
      `Proposed` für die Lockerung
      ([`AGENTS.md` §3.6](../../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)).
- [x] **Paritäts-Beleg der Zähl-Semantik** gegen den **realen** Bestand des
      Adopters, nicht gegen Fixtures: **84 von 84** Closure-Notizen ergeben
      **identische** Satzzahlen, und an der Adopter-Schwelle 2 sind beide
      Seiten symmetrisch (0 rot / 0 rot). Gefahren wurde Zahl gegen Zahl —
      d-check mit `min-sentences: 999`, damit **jede** Notiz ihre gezählte Zahl
      meldet, gegen die Shell-Pipeline des Adopters, Zeile für Zeile.

      **Zwei Zwischenstände waren rot und beide lagen an meiner
      Test-Konfiguration, nicht am Modul.** Erst fehlten 7 Dateien, weil sie
      `## N. Closure` ohne „-Notiz“ schreiben. Dann wichen 3 ab, weil mein
      geweitetes Muster auch `## 5. Closure-Trigger` traf — den der Adopter
      ausdrücklich ausschließt. Beides zeigt dieselbe Sache: **die Parität hängt
      am `heading-pattern`, nicht an der Zählung.**

      Nebenbefund, der in die Doku gehört: die Ausschluss-Bedingung des Adopters
      (`Closure`, aber nicht `Closure-Trigger`) ist in **RE2 ausdrückbar**, obwohl
      RE2 keinen Lookahead kennt — `^## .*[Cc]losure([^-]|-Notiz|$)`. Ein
      Adopter, der migriert, braucht dafür keinen Workaround. **Bewusst verengt (Schnitt-Review
      F-2):** eine Parität *in beide Richtungen über alle Prüfungen* ist nicht
      erreichbar, weil das Adopter-Skript zusätzlich die **Anzahl** der
      Closure-Abschnitte prüft und d-check laut Spezifikation nur den ersten
      liest. Diese Lücke ist ein Abnahme-Punkt von
      [slice-096](welle-69/slice-096-structure-modul-analyse.md), nicht dieses Slice.
- [x] Implementierung samt **fünf** geprüften Rückbauten (alle rot). Zwei blieben
      zunächst grün und haben Arbeit erzeugt: der **Tab** als Whitespace war
      ungetestet, und der Zeilenende-Zweig ist über die Modul-Oberfläche gar
      **nicht erreichbar** — der Abschnittstext endet immer auf einem
      Zeilenumbruch. Die Zusage gilt trotzdem der Funktion, also prüft sie jetzt
      ein direkter Tabellen-Test.
- [x] `make gates` + `make verify-closure-notes` grün.
- [ ] **Release als Minor** — **Wellen-Trigger, nicht Slice-Trigger.**
      [welle-72](../done/welle-72-closure-semantik.md) trägt beide Änderungen in
      **einem** Release mit **einer** Notiz; zwei aufeinanderfolgende Releases
      würden einem Konsumenten seine Gate-Semantik in zwei Schritten
      verschieben. Die Notiz muss den „findet mehr"-Hinweis (`closure-note-thin`)
      **und** den Hinweis auf die gelockerte Floskel-Prüfung tragen.

## 5. Risiken / offene Punkte

- **Der eigene Bestand könnte rot werden.** — **Ausgang: gemessen, wird er
  nicht.** Minimum fällt von 7 auf 5, unsere Schwelle ist 4. Der Abstand
  schrumpft aber von 3 auf 1 — die nächste dünne Notiz fällt eher auf, und das
  ist gewollt.
- **Konsumenten-Bruch:** ein grüner Lauf kann rot werden. Das ist der zugesagte
  Minor-Charakter, keine Überraschung — gehört aber in die Release-Notiz.
  — **Ausgang:** offen bis zur Release-Prep.
- **Die Gegenrichtung ist ein stiller Verlust:** ein Repo, das eine Floskel in
  Backticks stehen hat, verliert einen bestehenden Befund, ohne es zu merken.
  — **Ausgang: belegt und angenommen.** In **beide** Richtungen am Lauf
  nachgestellt: dieselbe Notiz mit einer Floskel in Backticks meldet mit der
  alten Fassung `closure-note-boilerplate` und mit der neuen **nicht**. Ein Test
  hält beide Seiten fest (zitierte Floskel trifft nicht, benutzte trifft weiter);
  die ADR benennt die Klasse für den Release-Text.
- **Fremd-Repo-Abhängigkeit:** die Paritäts-Fixtures liegen nicht in diesem
  Repo (Schnitt-Review F-7). — **Ausgang: erledigt, und besser als geplant.**
  Statt Fixtures wurde der **reale** Bestand beigezogen (84 Closure-Notizen) und
  gegen die **echte** Shell-Pipeline des Adopters gerechnet — nicht gegen eine
  Nacherzählung von ihr — ist das beim Start nicht möglich, wird der Paritäts-Beleg
  durch eigene Fixtures ersetzt und die Zusage entsprechend verengt.

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei — **und**
[slice-096](welle-69/slice-096-structure-modul-analyse.md) in `done/`. Die frühere
Fassung stellte diesen Slice davor; das ist nach dem Schnitt-Review umgekehrt
(F-3): er sagt Deckungsgleichheit einer Semantik zu, die slice-096 gerade neu
schneidet.

**Rückführungen:** `in-progress` → `next`, falls die Vorab-Messung zeigt, dass
die Angleichung eine Bestands-Sanierung nach sich zieht (dann ist das ein
eigener Slice, nicht ein größerer hier).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** berührt Produkt-Code (`internal/`) und Spec (`spec/`),
  beide unter dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001**
  (Datei-Register driften unbemerkt gegen ihre Autoritäts-Tabelle). Andere
  Klasse — dort geht es um eine Referenz **zwischen** Dokumenten, hier um die
  Zählung **innerhalb** eines Abschnitts. Nichts zu berücksichtigen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — die Zusage wird zuerst in Lastenheft und
Spezifikation geschärft, der Go-Code liefert sie. Kein Brownfield: es wird kein
bestehender undokumentierter Code inventarisiert, sondern eine bestehende,
dokumentierte Zusage präzisiert.

## 9. Closure-Notiz (nach `done/`)

Geliefert ist die Zähl-Parität: der Closure-Abschnitt wird **einmal** bereinigt
(Fences **und** Inline-Code), alle Bedingungen lesen diesen einen Text, und ein
Satzende zählt nur vor Whitespace oder Zeilenende
([`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
Lastenheft 0.55.0, Spezifikation Schritt C4,
[ADR-0053](../../adr/0053-eine-bereinigung-fuer-alle-closure-bedingungen.md)).

**Der Kern war nicht die Zählung, sondern dass es eine Bereinigung gibt.** Die
ursprüngliche Aufgabe klang nach zwei Detail-Regeln. Umgesetzt ist eine
Struktur-Entscheidung: vier Bedingungen, ein bereinigter Text. Die Alternative —
zwei getrennte Texte, einer je Richtung — hätte die Lockerung der Floskel-Prüfung
vermieden und dafür zwei Semantiken für denselben Abschnitt eingeführt. Genau die
Klasse, die das Register als **BEO-003** führt.

**Die Lockerung ist angenommen, nicht umgangen.** Eine *zitierte* Floskel ist
keine benutzte; die alte Fassung meldete jede Notiz, die **über** Floskeln
schreibt — diese Repo-Dokumentation eingeschlossen. Sie ist trotzdem eine
Lockerung, [`AGENTS.md` §3.6](../../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)
macht sie ADR-pflichtig, und die ADR benennt sie samt der Richtung, die
stillschweigend verschwindet, wenn man sie nicht hinschreibt.

**Die Messung hat einen Abnahme-Punkt erzeugt, den der Slice nicht hatte.** Von
4066 Satzende-Zeichen im eigenen Bestand tragen rund 2400 keinen Satz — Punkte
aus Link-Pfaden und Versionsnummern. Aber 170 sind echt und fallen mit weg:
`**Umsetzung.**`, ein fett gesetztes Satzende. Sie erklären den Abfall des
Minimums von 7 auf 5 fast vollständig. Entschieden wurde **für** die strengere
Regel — Parität ist der Zweck, die Richtung ist die sichere, und eine
Ausnahmeliste für `*`, `_`, `` ` ``, `)` wäre eine dritte Semantik, die weder
der Adopter noch CommonMark definiert.

**Der Paritäts-Beleg fiel stärker aus als die DoD verlangte.** Statt Fixtures
wurde der **reale** Bestand des Adopters beigezogen und gegen dessen **echte**
Shell-Pipeline gerechnet: 84 von 84 Notizen, identische Satzzahlen, an der
Adopter-Schwelle symmetrisch. Zwei Zwischenstände waren rot und **beide lagen an
meiner Test-Konfiguration** — erst ein zu enges Überschriften-Muster (7 Dateien
fehlten), dann ein zu weites, das auch `## 5. Closure-Trigger` traf. Die Lehre
steht im Slice: **die Parität hängt am `heading-pattern`, nicht an der Zählung.**
Nebenbefund fürs Handbuch: die Ausschluss-Bedingung des Adopters ist in RE2
ausdrückbar, obwohl RE2 keinen Lookahead kennt.

**Zwei Rückbauten blieben grün und haben Arbeit erzeugt.** Der Tab als
Whitespace war ungetestet; und der Zeilenende-Zweig ist über die
Modul-Oberfläche **gar nicht erreichbar**, weil der Abschnittstext immer auf
einem Zeilenumbruch endet. Kein toter Code — die Zusage gilt der Funktion —,
aber über das Modul nicht prüfbar. Ein direkter Tabellen-Test hält sie jetzt.
Beim Schreiben waren zwei meiner eigenen Erwartungen falsch; der Code hatte
recht.

**Offen und bewusst an die Welle abgegeben:** das Release. Es trägt beide
Änderungen dieser Welle in **einer** Notiz — zwei aufeinanderfolgende Releases
würden einem Konsumenten seine Gate-Semantik in zwei Schritten verschieben.
