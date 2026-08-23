# Slice slice-134: Die verbotene Direktive wirkt — `nolintlint` ins Profil

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-84-durchsetzung](../welle-84-durchsetzung.md).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.2 (Suppression-Verbot;
Auflösungs-Trigger `nolintlint`), §3.6 (Linter-Strenge nur per ADR);
[ADR-0006](../../adr/0006-lint-profil-solid.md) (Lint-Profil);
geschnitten vom Zensus in [slice-132](../done/slice-132-hard-rule-zensus.md).

**Berührte Spec-Stellen:** — (Lint-Profil; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Der Zensus hat §3.2 als **einzige** Regel gemessen, deren verbotene Form nicht
nur ungeprüft bleibt, sondern **wirkt**: ein echter Lint-Verstoß mit passendem
`//nolint` lässt `make lint` mit **Exit 0** durchlaufen. `nolintlint` — der
Linter, der genau das meldet — steht nicht im Profil.

Das ist die schwerste Klasse, die der Zensus gefunden hat: nicht eine Regel ohne
Wächter, sondern eine Regel, deren Umgehung **funktioniert und still bleibt**.

## 2. Vorgehen

1. `nolintlint` in [`.golangci.yml`](../../../../.golangci.yml) aufnehmen,
   Einstellungen bewusst wählen: verlangt wird mindestens eine **Begründung**
   und ein **spezifischer** Linter-Name; ob eine unbenutzte Direktive gemeldet
   wird, ist eigens zu entscheiden.
2. **Am Bestand messen, bevor scharfgeschaltet wird:** wie viele Befunde meldet
   das Profil heute? Bei 0 ist der Zug frei; bei >0 gehört jede Fundstelle
   einzeln beurteilt — geräumt oder als zentrale Ausnahme mit `Why:` geführt.
3. **Bewusstes Brechen — drei Formen, weil zwei Schalter getrennt greifen.**
   Der Zensus-Verstoß trug keine Begründung; erst die **vollständige** Form
   (benannter Linter **und** nachgestellter `// Grund`) passiert alle drei
   Schalter. Zu messen sind deshalb alle drei: nackte `//nolint` ⇒ rot ·
   benannt, aber unbegründet ⇒ rot · vollständig ⇒ **grün**. Der letzte, grüne
   Exit ist der Beleg für die verbleibende Lücke — und die Zusage des Gates ist
   genau die Differenz zwischen den ersten beiden und ihm.
4. `AGENTS.md` §3.2 nachziehen: der Auflösungs-Trigger ist eingelöst, die Regel
   wechselt von *einseitig* auf **teilgedeckt** — nicht auf *gedeckt*, denn
   `nolintlint` prüft die **Form** der Direktive, nicht ihre Berechtigung (§5).
   Der Zensus-Eintrag in
   [slice-132](../done/slice-132-hard-rule-zensus.md) bleibt als
   historische Messung stehen.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Lockerung anderer Linter**, um Platz zu schaffen. Eine Senkung wäre
  nach §3.6 ADR-pflichtig.
- **Keine Sanierung von Alt-Suppressions über die Messung hinaus.** Findet
  Schritt 2 Bestand, wird er benannt und einzeln entschieden — nicht pauschal
  ausgenommen.
- **Kein zweiter Wächter** auf andere Suppression-Formen (`//lint:ignore`,
  Build-Tags). Erst messen, ob es sie gibt.

## 4. Definition of Done

- [x] `nolintlint` ist im Profil, alle drei Schalter scharf, jeder mit seiner
      Begründung im Kommentar.
- [x] Der Bestand ist gemessen: **null** Direktiven, über **sechs** Muster
      (`//nolint`, `// nolint`, `//lint:ignore`, `//nolint:all`,
      `revive:disable`, `staticcheck:ignore`). Nichts zu räumen.
- [x] **Drei** Direktiv-Formen gemessen, jede mit gelesenem Exit **und**
      geprüfter Ursache im Log: nackt ⇒ rot (*should mention specific linter*) ·
      benannt-unbegründet ⇒ rot (*should provide explanation*) · vollständig
      ⇒ **grün**. Beide roten Läufe meldeten **genau einen** Befund, und zwar
      `nolintlint`.
- [x] `AGENTS.md` §3.2 sagt hin, was das Gate trägt und was nicht — samt des
      auf **permanent** zurückgesetzten Triggers.
- [x] `make gates` Exit 0 (neun Glieder, 475 Dateien, 0 Befunde); unabhängiger
      Review ([Report](../../../reviews/2026-08-23-slice-134-nolintlint-review.md)),
      nicht blockierend, beide MEDIUM eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Ein Linter, der Bestand meldet, wird gern per Ausnahme stumm gestellt.**
  — **Ausgang:** *nicht eingetreten, und die Versuchung gab es gar nicht.* Der
  Bestand ist null; es war keine Ausnahme nötig, und keine wurde gesetzt. Der
  `exclusions`-Block blieb unberührt — der Review hat gegengeprüft, dass
  `nolintlint` dort nirgends versehentlich mit ausgenommen ist, auch nicht für
  `_test.go`.
- **`nolintlint` prüft die Form, nicht die Berechtigung.** — **Ausgang:**
  *eingetreten, genau wie geschrieben — und meine eigene DoD hatte es zwei
  Absätze weiter oben trotzdem vergessen.* Sie verlangte, dass der
  Zensus-Verstoß jetzt rot wird; er wird es nicht. Der Widerspruch stand vier
  Zeilen auseinander im selben Dokument, und gefunden hat ihn der Review, nicht
  ich.

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-132](../done/slice-132-hard-rule-zensus.md)
in `done/`, WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Bestandsmessung so viele
Fundstellen zeigt, dass ihre Räumung ein eigener Slice ist.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Lint-Profil (GF), Gate-Landschaft (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-011`](../observations.md) für die Bestandsaussage in Schritt 2 — „null
  Fundstellen" ist eine Messung, keine Erwartung. [`BEO-007`](../observations.md)
  für jeden Beleg-Lauf.

Slice-ID: slice-134. Betroffene IDs: — (Lint-Profil; keine Anforderung).
Module: Lint-Profil, Gate-Landschaft. Gates: `make lint`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Profil-Erweiterung an bestehender Mechanik.

## 9. Closure-Notiz (nach `done/`)

Geliefert: `nolintlint` steht im Profil, alle drei Schalter scharf, und
[`AGENTS.md`](../../../../AGENTS.md) §3.2 sagt hin, was das Gate trägt und was
nicht.

**Das Ergebnis ist kleiner als der Zensus erwarten ließ — und das ist der Wert
dieses Slice.** Drei Direktiv-Formen, gemessen mit gelesener Ursache:

| Form | `make lint` | Meldung |
|---|---|---|
| `//nolint` | **rot** | *should mention specific linter* |
| `//nolint:unused,gochecknoglobals` | **rot** | *should provide explanation* |
| `//nolint:unused,gochecknoglobals // Grund` | **grün** | — |

Die dritte Zeile ist der Befund: **ein echter Lint-Verstoß mit vollständiger
Direktive läuft weiterhin durch.** Der Zensus-Fund besteht fort. Was sich
geändert hat, ist nicht die Möglichkeit der Umgehung, sondern ihr **Preis**: sie
muss ihren Linter nennen, sich begründen und wirken — sie wird sichtbar und
zurechenbar, statt still zu sein. §3.2 steht damit auf *teilgedeckt*, und der
Auflösungs-Trigger ist wieder **permanent**: die Berechtigungsfrage ist ein
Urteil, und ein zweiter Wächter darauf wäre der Heuristik-Wächter, den die Welle
ausschließt.

**Dreimal in einem Slice habe ich aus einem roten Exit den falschen Schluss
gezogen** — und jedes Mal war es dieselbe Bewegung. Im Zensus nannte die Probe
den falschen Linter (rot, aber aus anderem Grund). Hier nannte der Slice-Text
eine Form *wohlgeformt*, die ich nie gefahren hatte (sie wird gemeldet). Und die
DoD verlangte einen Beleg, den §5 desselben Dokuments vier Zeilen weiter
ausschließt. **Ein Exit-Code ist kein Beleg, solange seine Ursache ungelesen
bleibt** — die DoD verlangt das jetzt ausdrücklich.

**Und eine Zurechnung, die ich im selben Commit einmal richtig und einmal falsch
gemacht habe.** Bei [ADR-0006](../../adr/0006-lint-profil-solid.md) habe ich die „24 weiteren Linter" korrekt dem
u-boot-Vorbild zugerechnet und dadurch eine Drift **nicht** gemeldet, die es nie
gab. Drei Zeilen weiter habe ich `nolintlint` in genau diese 24 hineingezählt —
und im selben Kommentarblock geschrieben, er stehe *über das adoptierte Profil
hinaus*. Getrennt: 23 aus jener ADR, `nolintlint` aus §3.2, zusammen 24.

**Offen und benannt:** [`BEO-013`](../observations.md) ist in diesem Slice
entstanden — für Go ist die Frage *„wirkt diese Unterdrückung überhaupt noch?"*
seit heute gemessen, für den eigenen `<!-- d-check:ignore -->` nicht. Zwölf
aktive Marker, elf davon eingefroren in `done/`, einer in einem lebenden
Dokument.
