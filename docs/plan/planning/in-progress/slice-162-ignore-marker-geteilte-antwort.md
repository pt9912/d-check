# Slice slice-162: `d-check:ignore` beantwortet dieselbe Frage anders als `d-check:cite`

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
(Entscheidung 1: ein Konsument, der eine Lexik-Frage selbst beantwortet, ist ein
Defekt — und der Re-Evaluierungs-Trigger *„eine vierte Stelle"*);
[ADR-0060](../../adr/0060-citations-marker-scan-geteilte-prosa-antwort.md) (die
Skopierung, die diesen Slice schneidet);
[ADR-0025](../../adr/0025-codepaths-ignore-refs.md) (das Ventil selbst);
[slice-158](../done/slice-158-citations-inline-code.md) (der Anlass).

**Berührte Spec-Stellen:** [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
[`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[`DC-FA-DIAG-001`](../../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in)
— je die Ventil-Zusage.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

`d-check` kennt zwei Direktiven, und seit
[slice-158](../done/slice-158-citations-inline-code.md) beantworten sie
**dieselbe** Frage verschieden. *„Ist diese Zeile eine Direktive"* liest
`citations` auf dem inline-code-gestrippten Text; die vier Konsumenten des
Ventils lesen weiter roh:
[`codepaths.go`](../../../../internal/hexagon/core/rules/codepaths.go),
[`ids.go`](../../../../internal/hexagon/core/rules/ids.go),
[`versions.go`](../../../../internal/hexagon/core/rules/versions.go),
[`diagrams.go`](../../../../internal/hexagon/core/rules/diagrams.go). Eine
Erwähnung von `d-check:ignore` in Backticks **wirkt** dort als Ventil.

Das ist die Klasse aus
[ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
Entscheidung 1 — *zwei Antworten auf dieselbe Frage in einem Lauf sind ein
stiller Grün-Pfad, den kein Gate sieht* — und zugleich ihr eigener
Re-Evaluierungs-Trigger (*„eine vierte Stelle beantwortet eine Lexik-Frage
selbst"*).

**Gemessen** (Produkt-Lexik: Fence-Automat, absatzweise Spannen gleicher
Backtick-Länge, über 544 getrackte Markdown-Dateien):

| Lage von `d-check:ignore` | Zeilen |
|---|---|
| frei in Prosa (wirkt so oder so) | **63** |
| **nur** in Inline-Code (wirkt heute, würde nach Angleich nicht mehr) | **173** |

**Obergrenze der Sprengweite, gemessen:** wird das Ventil **ganz** abgeschaltet,
meldet der Repo-Lauf **58** Befunde (21 davon in
[`spec/lastenheft.md`](../../../../spec/lastenheft.md)). Der Angleich beträfe
nur die 173 Zeilen, also **höchstens** 58 — wie viele davon, ist die erste
Messung dieses Slice und **nicht** vorweggenommen.

## 2. Vorgehen

1. **Die Zahl unter der Obergrenze messen**, bevor entschieden wird: welche der
   58 Befunde hängen an einer Erwähnung in Inline-Code, welche an einem freien
   Marker? Ein Lauf mit gestripptem Ventil gegen einen mit rohem, Befundsätze
   verglichen.
2. **Die Richtung ist die unangenehme.** `citations` wurde **stiller**, das
   Ventil würde **lauter**: nach dem Angleich melden Zeilen, die heute
   schweigen. Jeder dieser Befunde ist einzeln zu beurteilen — echt (das Ventil
   stand nie legitim dort) oder Falsch-Rot (die Zeile braucht ein Ventil, nur
   ein anderes).
3. **Die Vertrags-Frage.** Vier Anforderungen tragen die Ventil-Zusage. Ist der
   Angleich eine Schärfung je Anforderung oder eine querschnittliche? Und trägt
   [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)s
   Trigger-Frage — *„kann ein **Gate** die Klasse prüfen statt eines Reviews"* —
   hier eine Antwort? Ein Wächter *„kein Konsument liest roh, wo die geteilte
   Antwort existiert"* wäre die strukturelle Reparatur statt der punktuellen.
4. Nur bauen, was die Messung trägt; die Entscheidung gegen einen Angleich wäre
   ebenso auszuweisen wie einer für ihn.
5. Bewusstes Brechen je gedeckter Form, **Ursache gelesen**; Rückbau grün.
6. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Änderung an `citations`.** Dessen Antwort ist entschieden
  ([ADR-0060](../../adr/0060-citations-marker-scan-geteilte-prosa-antwort.md)).
- **Keine Änderung an den drei gescopten Roh-Lesungen** aus
  [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) —
  `versions`-Pins, `pins`-Span, `immutable`-Core beantworten andere Fragen.
- **Kein Umschreiben der 173 Fundstellen.** Wenn ein Befund echt ist, wird die
  Zeile repariert; wenn nicht, das Ventil richtig gesetzt — keine
  Massen-Umformatierung, um ein Gate ruhigzustellen.

## 4. Definition of Done

- [ ] Die Zahl unter der Obergrenze ist **gemessen**, nicht geschätzt.
- [ ] Je aufgedecktem Befund eine Beurteilung: echt oder Falsch-Rot.
- [ ] Die Vertrags-Frage ist entschieden — Angleich oder begründete Skopierung —
      und die Straten hängen zusammen.
- [ ] Die Gate-Frage aus
      [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)s
      Trigger ist **beantwortet**, nicht übergangen.
- [ ] Ein konstruierter Verstoß je gedeckter Form, **Ursache gelesen**.
- [ ] Doku-Currency: Handbuch, `README`-Fassungen, `CHANGELOG` mitgezogen.
- [ ] `make gates` grün (Exit explizit), `make fullbuild` grün; unabhängiger
      Review.

## 5. Abnahme-Punkte / Risiken

- **Der Angleich macht ein Gate lauter, nicht leiser.** Bis zu 58 heute stille
  Zeilen melden danach — und ein Gate, das auf einen Schlag viel meldet, wird
  abgeschaltet statt gelesen. — **Ausgang:**
- **Die Reparatur könnte punktuell statt strukturell ausfallen.** Vier
  Konsumenten einzeln umzustellen löst diesen Fall und nicht die Klasse; der
  fünfte Konsument entstünde morgen wieder roh. — **Ausgang:**
- **Ein Ventil, das enger wird, kann legitime Ausnahmen entwerten.** Wer
  `d-check:ignore` bisher in Backticks setzte, tat es womöglich absichtlich und
  in gutem Glauben. — **Ausgang:**

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung zeigt, dass der
Angleich eine Auftraggeber-Entscheidung über den Ventil-Vertrag verlangt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Spec-Straten (GF), Doku (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-020`](../observations.md) — die gemessene Menge muss die sein, über die
  geredet wird (die 58 sind eine **Obergrenze**, keine Antwort);
  [`BEO-011`](../observations.md) — die Regel gehört aus dem Bestand, nicht aus
  dem Anlass; [`BEO-017`](../observations.md) — ein rotes Gate muss vom
  geprüften Grund kommen.

Slice-ID: slice-162. Betroffene IDs:
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
[`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[`DC-FA-DIAG-001`](../../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in).
Module: `codepaths`, `ids`, `versions`, `diagrams`.
Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Angleich einer vorhandenen Antwort an eine
vorhandene geteilte Antwort.

## 9. Closure-Notiz (nach `done/`)

— (offen)
