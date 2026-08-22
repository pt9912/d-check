# Slice slice-115: `ARC-*`-Vergabe in der Architektur-Sicht + `diagrams` als Konsument

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-80-struktur-ids](../welle-80-struktur-ids.md) (zugeordnet bei
der Eröffnung).

**Bezug:**
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(Linkpflicht `ARC-\d{3}` seit slice-113), Modul `diagrams` (Kennungen in
`mermaid`-Fences gegen eine `defined-in`-Quelle — der gebaute Konsument für
genau diese Kennungs-Art), Baseline-Template `architecture.template.md` §1–§3
(„Hier werden die `ARC-*` für Komponenten vergeben — eine Zeile je Kasten des
Diagramms"; §2 nennt dieselben, §3 setzt die Reihe fort),
[`MR-006`](../../../../harness/conventions.md#mr-006) (Referenzrichtung —
Spezifikation nennt keine `ARC-*`), Entscheide D1/D3.

**Berührte Spec-Stellen:** `architecture.md` §1 Komponenten-Übersicht, §2
Schichten und Constraints, §3 externe Abhängigkeiten, §5 Fehlermodelle — der
Verweis zeigt aufwärts.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

`spec/architecture.md` trägt `ARC-<NNN>` fortlaufend je Datei: §1 bekommt
unter dem Komponenten-Diagramm eine **Komponenten-Tabelle** (eine Zeile je
Kasten — heute sieben Rollen: CLI, CORE, FS, HTTP, VCS, CFG, REP — mit
Kennung, Rolle, Zweck), die Kasten-Labels im `mermaid`-Fence tragen die
Kennung mit; §2 Schichten-Tabelle nennt dieselben Kennungen; §3 externe
Abhängigkeiten setzen die Reihe fort; §5 Fehlermodelle bleibt kennungslos —
die Vorlage vergibt an Komponenten und Berührungspunkte, nicht an
Fehlerquellen (Korrektur des ersten Schnitts, die Vorlage sticht). Der
Konsument:
`diagrams` opt-in (`defined-in: spec/architecture.md`) prüft, dass jede im
Diagramm genannte Kennung definiert ist — **vorher am Bestand gemessen**, im
Vergabe-Commit scharf (der `pre-commit`-Hook bindet Commits an Grün). Und eine
Messung, die die Referenzrichtung schützt: `spec/spezifikation.md` nennt
**keine** `ARC-*` (Spezifikation → Architektur wäre abwärts, `matrix`
`spec-straten` `no-downward`).

## 2. Vorgehen

1. **Inventur:** Kästen des §1-Flowcharts (heute 7), Zeilen §2 (7), §3 (4)
   — gemessen; Hard Rule §3.4 im Blick (sprach-/meilensteinfrei:
   Kennungen benennen Rollen, keine Pakete).
2. **Vergabe:** §1 Komponenten-Tabelle `| Kennung | Rolle | Zweck |` mit
   `ARC-NNN` fortlaufend; Diagramm-Labels `ARC-NNN CLI` usw.; §2 erste Spalte nennt die
   Kennung; §3 setzt fort, §5 bleibt kennungslos. Kopf der Datei: ein Satz zur Kennungs-Form
   (Struktur-ID, keine Anforderung — Baseline-Wortlaut).
3. **Konsument:** `diagrams` in `.d-check.yml` opt-in mit `defined-in:
   spec/architecture.md` und dem `ARC-\d{3}`-Muster (Wortgrenzen — der
   Präzedenz-Befund `ARC-\d{2}` gegen eine dreistellige Kennung); Probe rot-vorher (Diagramm mit
   Kennung, Tabelle ohne) / grün-nachher, dazu eine konstruierte Gegenprobe
   (eine nicht definierte Tipp-Kennung im Fence ⇒ Befund).
4. **`matrix`-Messung:** `git grep -n 'ARC-' spec/spezifikation.md` ⇒ null;
   in der Closure-Notiz festgehalten.
5. **Spiegel** ([`MR-025`](../../../../harness/conventions.md#mr-025),
   `grep` nach „Komponenten-Übersicht", „§2 (Reporter-Rolle"): ADR-Verweise auf
   `architecture.md §N` (nur Proposed-ADRs nachziehen — slice-116), Handbuch
   (falls es die Sicht zitiert), `harness/README.md` §Sensors (doc-check-Zeile:
   `diagrams` kommt als Modul dazu — Modul-Liste der `.d-check.yml` ist per
   Go-Test an [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) gebunden: `diagrams` ist netzlos, der Test muss grün
   bleiben), AGENTS §4 (`doc-check`-Zeile nennt die Module?).
6. Unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine `SPEC-*`** (slice-114), **keine ADR-Nachzüge** (slice-116).
- **Keine Technologie-Namen** in der Sicht (Hard Rule §3.4 bleibt).
- **Kein Produkt-Code** am Modul `diagrams`; findet die Probe einen Defekt, ist
  das ein eigener Slice außerhalb der Welle.

## 4. Definition of Done

- [x] §1 Komponenten-Tabelle + Diagramm-Labels, §2 referenzierend, §3
      fortlaufend (§5 kennungslos nach Vorlage); Zählung gemessen in der
      Closure-Notiz.
- [x] `diagrams` opt-in scharf im Vergabe-Commit; Messung rot-vorher /
      grün-nachher + konstruierte Gegenprobe dokumentiert; [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Modul-
      listen-Test grün.
- [x] `matrix`-Messung: keine `ARC-*` in der Spezifikation.
- [x] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Diagramm-Labels mit Kennung** könnten das Flowchart unleserlich machen —
  Form: Kennung als Präfix im Label, Rolle bleibt lesbar. — **Ausgang:**
  entfallen — die Labels tragen die Kennung als Präfix, die Rolle steht
  unverändert dahinter; der Review meldete keinen Lesbarkeits-Befund.
- **`diagrams`-Regex-Härte** (Wortgrenzen) — die konstruierte Gegenprobe ist
  der Wächter. — **Ausgang:** **eingetreten als benannte Grenze.** Die
  Wortgrenzen halten, aber die feste Dreistelligkeit lässt eine vertippte
  zwei- oder vierstellige Kennung im Diagramm still durch (Review F-6);
  als vierte Grenze im Config-Kommentar benannt.
- **Doppelte Wahrheit §1-Tabelle vs. §2-Tabelle** (beide nennen die Kennung):
  §1 vergibt, §2 referenziert — so steht es im Template; der Review prüft, dass
  §2 keine eigene Kennung einführt. — **Ausgang:** entfallen — der Review hat
  die Mengengleichheit gemessen: sieben zu sieben, keine Neuvergabe in §2.

## 6. Trigger

**Start** (`open` → `in-progress`): slice-113 in `done/`; unabhängig von
slice-114 (andere Datei, andere Reihe).

**Rückführungen:** `in-progress` → `next`, falls `diagrams` am eigenen Diagramm
einen Modul-Defekt zeigt (dann erst Modul-Slice).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Spec-Stratum Sicht (`spec/architecture.md`, GF),
  Prüf-Profil (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): BEO-002 ([`MR-025`](../../../../harness/conventions.md#mr-025)),
  BEO-006/009 Arbeitsregeln; BEO-004 (Modul liest, was es nicht scannt —
  `diagrams` liest `defined-in`: die Quelle liegt im Scan, aber die Frage
  „welche Eingaben liest das Modul, die es nicht scannt" wird im Review
  gestellt).

Slice-ID: slice-115. Betroffene IDs:
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(Modul-Listen-Test), [`MR-006`](../../../../harness/conventions.md#mr-006).
Module: Architektur-Sicht, Prüf-Profil (`diagrams`), harness/AGENTS-Zeilen.
Gates: `make doc-check` (eng), `make test` (Modul-Listen-Test), `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Vergabe im eigenen Spec-Stratum nach
Baseline-Form; kein Legacy-Import.

## 9. Closure-Notiz (nach `done/`)

**Geliefert:** elf Kennungen in der Architektur-Sicht — §1 vergibt sieben
Komponenten in einer neuen Komponenten-Tabelle und trägt sie im Kasten-Label
des Flowcharts, §2 nennt dieselben sieben in einer neuen ersten Spalte ohne
Neuvergabe, §3 setzt mit vier externen Berührungspunkten fort. Der Konsument
ist das Modul `diagrams`, gescopt auf die Spec-Straten und in der
Netzlos-Modulliste des Gate-Tests verankert: fällt es aus der Config, wird
`make test` rot. Die Referenz-Richtung ist gemessen — die Spezifikation nennt
keine Kennung der Sicht.

**Review** ([Report](../../../reviews/2026-08-22-slice-115-arc-vergabe-review.md)):
APPROVE mit Auflagen — 0 HIGH, 2 MEDIUM, 2 LOW, 2 INFO, sechzehn
Negativ-Proben. Alle sechs eingearbeitet und je mit dem Produkt belegt.

**Was ging anders als geplant — dreimal:**
1. **Der Plan wollte auch §5 Fehlermodelle bekennen; die Vorlage nicht.**
   Struktur-IDs gehören dort an Komponenten und Berührungspunkte, nicht an
   Fehlerquellen. Die Vorlage sticht den eigenen Plan — die Entscheidung
   steht jetzt im Slice und im Wellendokument, nicht nur in einer
   Commit-Botschaft.
2. **Eine Gegenprobe war grün, wo ich rot erwartet hatte** — und das war der
   Fund, nicht der Fehler: entfernt man nur die Vergabe-Zeile aus §1, bleibt
   die Kennung definiert, weil §2 sie nennt. Das Modul liest Token, nicht
   Struktur; **Nennung und Vergabe sind für es dasselbe**. Der Wächter fängt
   die verwaiste Kennung im Diagramm, nicht die halbe Löschung in der Sicht.
3. **Ein neues Modul in der Config hat einen Spiegel im Makefile**, den ich
   nicht kannte: `FOCUS_DISABLE` wählt für die fokussierten Gates alle
   Datei-Module ab und spiegelt dafür die Modulliste. Nicht nachgezogen hätte
   das neue Modul im `pre-commit`-Hook auf ungestagtes WIP über-gefeuert. Der
   Kommentar sagte es ausdrücklich — gelesen hat ihn der Review.

- **Steering-Loop-Eintrag:** Sensor ergänzt: Kennungen in Diagramm-Fences der
  Spec-Straten sind an ihre Definitionsquelle gebunden — liegt in
  `.d-check.yml §diagrams`. Kein Auslöser aus dem Register.
- **Beobachtungs-Register (`../observations.md`):** keine neue Beobachtung.
  Die Makefile-Spiegel-Lücke ist die verkörperte Klasse BEO-002 („eine
  Änderung wird nur im Dokumentkörper nachgezogen, ihre Ränder bleiben
  stehen"), hier als Ausführungs-Rand statt als Prosa-Rand; zitiert statt neu
  formuliert, der Zähler bleibt.
- **Folge-Slices:** [slice-116](../in-progress/slice-116-adr-neuzugangs-regel.md) —
  sein Trigger (114 und 115 in `done/`) ist mit dieser Closure eingetreten.
- **Risiken aus §6:** alle mit Ausgang (§5) — eines als benannte Grenze
  eingetreten, zwei entfallen.
- **Drei Paarungen:** Wellen-Slice — die Paarungen prüft die Welle-Closure.
