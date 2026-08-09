# Slice slice-094: Zähl-Parität der Closure-Note-Struktur (Inline-Code + Satzende-Form)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** ohne Welle — der Slice ist einzeln lieferbar und trägt keine
Closure-Bedingung, die von seiner DoD verschieden wäre.

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

## 3. Abnahme-Punkte

1. **Reicht die Angleichung, oder braucht die Zählung ein Config-Ventil?** Der
   Vorschlag ist **ohne** Ventil: eine feste, an CommonMark orientierte Zählung
   (Inline-Spans sind Code, kein Fließtext) statt eines Schalters, den niemand
   bewusst setzt. Ein Ventil wäre eine zweite Semantik für dieselbe Frage.
2. **Wird der Bestand rot?** Vorab zu messen — dieselbe Disziplin wie bei der
   Schwellen-Wahl in [slice-093](../done/slice-093-closure-note-gate.md). Wird
   er es, ist das ein Befund über die Notizen, kein Grund, die Zählung
   aufzuweichen.

## 4. Definition of Done

- [ ] Zählung angeglichen (Inline-Spans entfernt, Satzende nur vor
      Whitespace/Zeilenende); Spezifikation Schritt C4 und das Akzeptanzkriterium
      nachgezogen.
- [ ] **Paritäts-Mutations-Beleg** gegen die a-check-Fixtures: jede Fixture, die
      das Adopter-Skript rot macht, macht auch das Modul rot — und umgekehrt.
- [ ] `make gates` + `make verify-closure-notes` grün; Release als **Minor**
      (Handbuch-§11-Zeile mit dem „findet mehr"-Hinweis).

## 5. Risiken / offene Punkte

- **Der eigene Bestand könnte rot werden** (Notizen, deren Sätze knapp über der
  Schwelle liegen und Inline-Code enthalten). — **Ausgang:** offen bis zur
  Vorab-Messung in Abnahme-Punkt 2.
- **Konsumenten-Bruch:** ein grüner Lauf kann rot werden. Das ist der zugesagte
  Minor-Charakter, keine Überraschung — gehört aber in die Release-Notiz.
  — **Ausgang:** offen bis zur Release-Prep.

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei. Der Slice ist
Voraussetzung dafür, dass a-check sein `verify-closure-notes`-Skript
**paritätstreu** ablösen kann — er steht deshalb vor
[slice-096](slice-096-structure-modul-analyse.md).

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

_Ausstehend._
