# Slice slice-101: Unbalancierter Fence verschluckt still den Rest des Abschnitts

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** [welle-70-fence-lexik](../welle-70-fence-lexik.md), eröffnet am
2026-08-09. Die Welle bündelt nur diesen Slice — nicht wegen eines Mehr
gegenüber der DoD, sondern weil `make planning-check` einen Slice in Arbeit ohne
benannte aktive Welle nicht zulässt (die Zwei-Zustands-Kopplung aus
[`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform)).

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Closure-Note-Struktur), [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md)
(die Fence-Lexik und ihre bewusst offen gelassene Grenze).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Einen **ausgelieferten stillen Grün-Pfad** schließen: ein ungeschlossener oder
verschachtelter Fenced-Code-Block verschluckt in v0.52.0 alles hinter sich, und
die Bedingungen der Closure-Note-Struktur laufen darüber grün.

## 2. Der Beleg

Nachgestellt am 2026-08-09 gegen das veröffentlichte Image, mit einer wörtlich
konfigurierten Floskel hinter einem ungeschlossenen Fence:

```text
Floskel hinter ungeschlossenem Fence  → 0 Befund(e), Exit 0
dieselbe Floskel ohne den Fence       → closure-note-boilerplate, Exit 1
```

Ein Autor, der einen Code-Block nicht schließt, schaltet damit unbemerkt die
Prüfung des restlichen Abschnitts ab. Das ist die schwerste Befund-Klasse dieses
Repos — ein Gate, das grün meldet, ohne geprüft zu haben.

**Der Defekt ist älter als die Closure-Fähigkeit.** Der Fence-Automat toggelt
naiv; [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) hat den
längenabgeglichenen Fence-Schluss ausdrücklich offen gelassen („naiver-Toggle-vs-
strikter-Schluss bewusst offen"). Was damals eine vertretbare Grenze war, ist mit
einer Bedingung, die **innerhalb** eines Abschnitts misst, zu einem Silent-Grün
geworden.

## 3. Abnahme-Punkte

1. **Wie weit geht der Fix?** (a) Nur die Closure-/Struktur-Bedingungen behandeln
   einen unbalancierten Fence als „bis Abschnitts-Ende offen" und melden ihn;
   (b) der geteilte Fence-Automat bekommt den längenabgeglichenen Schluss aus
   CommonMark, was **alle** Module berührt; (c) ein eigener Grund-Code für den
   unbalancierten Fence. Zu entscheiden — (b) ist die Wurzel, aber auch die
   Änderung mit der größten Reichweite.
2. **Bestandsmessung vor der Wahl:** wie viele Dateien im eigenen Repo und in den
   Adopter-Repos tragen unbalancierte Fences? Die Zahl entscheidet mit, ob (b)
   ein Minor oder ein Aufräum-Projekt ist.

## 4. Definition of Done

- [ ] Abnahme-Punkte 1–2 entschieden, mit Messung belegt; Vertragsanpassung
      (Lastenheft/Spezifikation) und ggf. ADR.
- [ ] Der oben belegte Fall meldet; Test mutations-echt (Rückbau des Fixes macht
      ihn wieder grün).
- [ ] `make gates` + `make verify-closure-notes` grün; Release als **Minor**
      (d-check findet danach mehr).

## 5. Risiken / offene Punkte

- **Reichweite:** Variante (b) ändert die Lexik für alle Module — dieselbe Klasse
  wie [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md), wo ein
  Differential-Spike gegen einen echten Parser die Grundlage war. — **Ausgang:**
  offen bis Abnahme-Punkt 1.
- **Der Defekt ist ausgeliefert.** Bis zum Fix ist die Zusage der
  Closure-Note-Struktur schwächer als dokumentiert. — **Ausgang:** offen; zu
  entscheiden, ob das in die Release-Notiz von v0.52.0 nachgetragen wird.

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei. **Vor**
[slice-099](../open/slice-099-structure-modul.md) sinnvoll — sonst erbt das neue Modul
den bekannten stillen Grün-Pfad über die geteilte Mechanik.

**Rückführungen:** `in-progress` → `open`, falls die Bestandsmessung zeigt, dass
Variante (b) eine eigene Sanierung nach sich zieht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001**; andere
  Klasse, nichts zu berücksichtigen. **Kandidat für einen neuen Eintrag:** „eine
  bewusst offen gelassene Lexik-Grenze wird zum Silent-Grün, sobald ein neues
  Modul innerhalb der Grenze misst" — bei der Closure dieses Slice zu prüfen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — die Korrektur wird zuerst als Zusage
formuliert (welches Verhalten gilt bei unbalanciertem Fence?), dann geliefert.
Kein Brownfield: es wird kein undokumentierter Bestand inventarisiert, sondern
eine dokumentierte Grenze neu bewertet.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
