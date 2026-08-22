# Slice slice-112: Wellen-Invariante — Kennungs-Mengen im Lastenheft-Wortlaut (0.62.1)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** ohne Welle — Redaktions-Slice ohne eine Closure-Bedingung, die von
seiner DoD verschieden wäre (Baseline-Regelwerk `modul-06-roadmap.md` §Wann
Arbeit eine Welle braucht). Unter `waves.mode: many` ist der wellenlose Slice
ein legitimer, gemessener Zustand der eigenen Roadmap.

**Bezug:**
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
§Wellen-Invariante (Aussage 1/2, `wave-drift`),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Byte-Identität des Befundsatzes bleibt unberührt),
[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
(Proposed — Entscheidung 6 wird im Wortlaut präzisiert und per
`## Geschichte` fortgeschrieben, nicht ersetzt),
[`MR-025`](../../../../harness/conventions.md#mr-025) (Semantik-Fläche ⇒
Spiegel-Liste vor dem Editieren). Anlass: der INFO-Befund F-5 des
[slice-111-Reviews](../../../reviews/2026-08-21-slice-111-waves-mode-review.md)
— als Präzisierungs-Kandidat in `welle-79-results.md` §Folge-Slices benannt,
Auftraggeber-Entscheid 2026-08-22: jetzt, als Mini-Slice, ohne Release.

**Berührte Spec-Stellen:** `spezifikation.md` §[`DC-FA-PLAN-001.a`](../../../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) (W3,
`one`-Satz) und §4 (`wave-drift`-Zeile) — der Verweis zeigt aufwärts, die Spec
nennt diesen Slice nicht.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Der Lastenheft-Wortlaut der ersten Wellen-Aussage sagt, was das Prädikat
tut: **beide** Modi vergleichen **Kennungs-Mengen**, nicht Datei-Mengen —
zwei flache Wellendokumente derselben Kennung sind **ein** Element (unter
`one` „genau eine Kennung", unter `many` Kennungs-Menge gegen
Kennungs-Menge). Heute liest sich Zeile 1/2 als Datei-Menge („genau ein
flaches Wellendokument", „Menge der flachen Wellendokumente"), implementiert
ist seit v0.59.0 die Kennungs-Map aus W2 — und die Spezifikation sagt es in
W2 bereits („zwei Dateien mit derselben Kennung zählen als eine Welle"), W3
und §4 aber nicht. Kein Verhaltens-Unterschied, darum **kein Release**; der
Beleg der präzisierten Aussage ist ein **Pinning-Test** mit zwei
gleich-kennigen Dateien in beiden Modi.

## 2. Vorgehen

1. **Lastenheft-Commit zuerst (Doc führt):** Version 0.62.0 → 0.62.1.
   Zeile 1/2 der Aussage-Tabelle: `one` ⇒ „unter den flachen Wellendokumenten
   liegt **genau eine** Wellen-Kennung (Ruhe-Marker ⟺ keine)", `many` ⇒
   „…gleich der **Kennungs-Menge** der flachen Wellendokumente", plus der
   Satz, der die Vergleichsgröße für beide Modi nennt (Kennungen, nicht
   Dateien; zwei flache Dokumente derselben Kennung sind ein Element). Neues
   Akzeptanzkriterium **Wellen-Boundary (gleiche Kennung, beide Modi)**.
   §7-Historie-Zeile 0.62.1 (Präzisierung ohne Verhaltensänderung; die
   Doppel-Dokument-Frage bleibt ausdrücklich offen).
2. **Spezifikation:** W3 `one`-Satz auf „genau eine Wellen-Kennung (W2)",
   die Zahl-Nennung auf Kennungen, §4-`wave-drift`-Zeile in beiden
   Modus-Hälften, §7-Historie-Zeile, Kopf „Letzte Änderung".
   **[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)**
   Entscheidung 6 im Wortlaut („gegen die Kennungs-Menge der flachen
   Wellendokumente") + `## Geschichte`-Zeile (Status bleibt Proposed).
3. **[`MR-025`](../../../../harness/conventions.md#mr-025)-Spiegel-Liste vor
   dem ersten Edit:** Lastenheft Zeile 1/2 + Akzeptanzkriterien + §7 ·
   Spezifikation W3 + §4-Zeile + §7 + Kopf · [ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md) Entscheidung 6 +
   Geschichte · Benutzerhandbuch §6-Absatz „Zwei Kardinalitäts-Modelle" +
   Kopf (1.54, Stand) + §11-Zeile · CHANGELOG `[Unreleased]` ·
   **nicht** die Meldungs-Klartexte in `planning_waves.go`/`diagnose.go`
   (sie zählen dieselbe Kennungs-Map; eine Textänderung bräche die
   zugesagte Byte-Identität des `one`-Pfads, §3) · Config-Kommentare der
   eigenen Profile (Prosa, bleibt korrekt).
4. **Pinning-Test** `TestWavesGleicheKennungZaehltEinmal`
   (`planning_waves_test.go`): zwei flache Dateien derselben Kennung — unter
   `one` mit genanntem Aktiv-Status **kein** Befund, unter `many` mit einem
   Zeiger dieser Kennung **kein** Befund.
5. **Handbuch 1.54** (Doku-only-Zeile in §11, Software-Version bleibt
   v0.62.0 — Präzedenz 1.29) + **CHANGELOG `[Unreleased]`** (Doku-Präzisierung,
   kein Tag).
6. **Unabhängiger Review** (Wortlaut-Präzision, Spiegel-Vollständigkeit,
   Test-Aussagekraft), Auflagen einarbeiten; **Closure:** Move nach `done/`
   mit Roadmap-Flip (Ruhe-Marker zurück), Closure-Body,
   `make fullbuild` mit explizitem Exit.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Verhaltensänderung** — kein Grund-Code, keine Meldungs- oder
  `--doctor`-Textänderung (der `one`-Pfad bleibt byte-identisch zu v0.62.0),
  kein Release/Tag.
- **Keine Entscheidung, ob zwei gleich-kennige flache Dokumente selbst ein
  Befund sein sollten** (Doppel-Dokument). Der Test pinnt das
  **Ist**-Verhalten; die Wunsch-Semantik wäre ein eigener CR.
- **Kein Mehr-Wellen-Betrieb-Thema** — das ist ein eigener Roadmap-Entscheid
  (am selben Tag getroffen, eigener Commit).

## 4. Definition of Done

- [x] Lastenheft 0.62.1 (Zeile 1/2 + Akzeptanzkriterium + Historie) liegt
      als **erster** Commit vor Spezifikation, ADR und Test in der Historie.
- [x] Spezifikation (W3, §4, §7, Kopf) und
      [ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
      (Entscheidung 6 + Geschichte) nachgezogen;
      [`MR-025`](../../../../harness/conventions.md#mr-025)-Liste aus §2
      Schritt 3 abgehakt.
- [x] Pinning-Test grün in beiden Modi (`make test`), `make gates` grün.
- [x] Handbuch 1.54 + CHANGELOG `[Unreleased]` nachgezogen.
- [x] Unabhängiger Review vor der Closure; Closure-Notiz mit
      Lerneintrag; Beobachtungs-Register gesichtet (§7 zitiert den Stand);
      `make fullbuild` grün (Exit explizit).

## 5. Abnahme-Punkte / Risiken

- **Lesbarkeit der Zeile 1/2:** die Präzisierung darf die Tabellenzeile
  nicht in Prosa ertränken — ein Zusatz-Satz für beide Modi statt je Modus
  eine Klammer. — **Ausgang:** entfallen — ein Zusatz-Satz für beide Modi, die
  Zeile blieb eine Tabellenzeile; der Review meldete keinen Lesbarkeits-Befund.
- **Der Test pinnt eine Toleranz:** zwei gleich-kennige Dateien sind kein
  Befund. Das ist Ist-Verhalten und als Grenze in §3 benannt; der Test
  trägt die Aussage „Kennung ist die Vergleichsgröße", nicht „Doppel-
  Dokumente sind erwünscht". — **Ausgang:** entfallen als Risiko — die Grenze
  steht in §3 und im Akzeptanzkriterium selbst; ob ein Doppel-Dokument
  meldepflichtig wird, ist eine Produktfrage (CR-Kandidat), kein Risiko dieses
  Slice.
- **Spiegel-Vollständigkeit** (BEO-002 ⇒ [`MR-025`](../../../../harness/conventions.md#mr-025)): die Liste in §2 Schritt 3
  ist die Antwort; der Review prüft sie gegen den Diff. — **Ausgang:**
  **eingetreten** — drei Stellen mit der alten Lesart blieben in drei ohnehin
  bearbeiteten Dateien stehen (Review F-1/F-2, MEDIUM), eingearbeitet; Klasse
  BEO-002 (verkörpert), Gegenmittel in
  [`MR-025`](../../../../harness/conventions.md#mr-025) geschärft (Ableiter:
  `grep` nach dem alten Wortlaut).
- **Arbeitsregeln aus dem Register:** BEO-006 (`git status` vor
  pfad-selektiven Commits), BEO-009 (Botschaft erst nach der Probe,
  `git diff --stat` gegen jede Botschafts-Zeile), BEO-007-Regel (Gate-Exits
  bindend). — **Ausgang:** entfallen — BEO-006/009 traten nicht ein; der
  `pre-commit`-Hook (BEO-007 verkörpert) fing zweimal Rot vor dem Commit: drei
  nackte Kennungen im Slice-Dokument, ein als Link geparster Test-Fixture-Pfad
  im Review-Report.

## 6. Trigger

**Start** (`open` → `in-progress`): Auftraggeber-Entscheid 2026-08-22 („jetzt
als Mini-Slice, kein Release") — eingetreten; der Move folgt unmittelbar.

**Rückführungen:** `in-progress` → `next`, falls die Präzisierung eine
**Verhaltensfrage** aufwirft (Doppel-Dokument als Befund) — dann ist es ein
CR, keine Redaktion. `in-progress` → `open`: entfällt (keine externe
Abhängigkeit).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Spec-Straten (`spec/`, GF — Repo-Default),
  Harness-/Planungs-Doku (GF), Benutzerhandbuch (GF), Produkt-Kern-Test
  `internal/hexagon/core/rules` (GF) — keine Sub-Area unter der Schwelle.
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-21): BEO-006,
  BEO-008, BEO-009 offen bei je 2; BEO-007 verkörpert. BEO-008
  (Pin-Hebungs-Spiegel) ist nicht berührt; BEO-006/BEO-009 gelten als
  Arbeitsregeln (§5); BEO-002 wirkt verkörpert als
  [`MR-025`](../../../../harness/conventions.md#mr-025) — die Spiegel-Liste
  steht in §2 Schritt 3.

Slice-ID: slice-112. Betroffene IDs:
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus);
[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
(Wortlaut + Geschichte). Module: Spec (Lastenheft + Spezifikation), ADR,
Handbuch, CHANGELOG, Test des Moduls `planning` (Kern `rules/`, kein
Produktions-Code). Gates: `make test` (eng), `make doc-check`, `make gates`,
`make fullbuild` (Closure).

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — redaktionelle Präzisierung einer eigenen,
spezifizierten Anforderung samt Pinning-Test; kein Legacy-Import, keine
Inventur.

## 9. Closure-Notiz (nach `done/`)

**Geliefert:** Lastenheft 0.62.1 zuerst (Zeile 1/2, Prosa-Absatz, neues
Akzeptanzkriterium, Historie), dann die Spiegel — Spezifikation W3/§2-Schema/§4/§7,
[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
Entscheidung 6 + Geschichte (Proposed), Handbuch 1.54 (§6 zweimal, §11, Kopf),
CHANGELOG `[Unreleased]` — und der Beleg: `TestWavesGleicheKennungZaehltEinmal`
in beiden Modi, Mutations-Gegenprobe per Dateikopie (Datei- statt
Kennungs-Schlüssel ⇒ beide Subtests rot, dazu sieben Bestandstests). Kein
Produktions-Code, kein Release; der wellenlose Slice lief als gemessener
Zustand durch `planning-check` (Marker raus, Liste leer, Bijektion ∅ ⟺ ∅).

**Review** ([Report](../../../reviews/2026-08-22-slice-112-kennungs-menge-review.md)):
APPROVE mit Auflagen — 0 HIGH, 2 MEDIUM, 1 LOW, 1 INFO. Die beiden MEDIUM waren
dieselbe Klasse: drei Stellen mit der alten Lesart in drei Dateien, die der
Slice ohnehin bearbeitete (Lastenheft-Prosa „Zwei Kardinalitäts-Modelle",
Spezifikation-§2-Schema-Zeile, Handbuch-§6-Vorsatz) — eingearbeitet. F-3 (LOW)
bleibt bewusst stehen: die Kurzformen „Kennungs-Bijektion ⟺ Dateien" in
README/AGENTS sind Gate-Beschreibungen, in denen „Kennungs-" die Vergleichsgröße
bereits nennt; im gelebten Betrieb ist eine Welle eine Datei. F-4 (INFO)
trägt §3.

**Was ging anders als geplant:** Die [`MR-025`](../../../../harness/conventions.md#mr-025)-Spiegel-Liste stand im Plan (§2
Schritt 3) und war trotzdem lückenhaft — sie war erinnert, nicht abgeleitet.
Der Ableiter für eine Wortlaut-Präzisierung ist das `grep` nach dem **alten**
Wortlaut über den ganzen Baum; er hätte alle drei Stellen gefunden.

- **Steering-Loop-Eintrag:** Guide geschärft: der Wortlaut-Ableiter steht jetzt
  ausdrücklich in [`MR-025`](../../../../harness/conventions.md#mr-025)
  — liegt in [`harness/conventions/MR-025-spiegel-vor-dem-editieren.md` §Adaption](../../../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md).
  Auslöser: BEO-002 (verkörpert seit slice-099; hier erneut eingetreten, nicht
  weitergezählt).
- **Beobachtungs-Register (`../observations.md`):** BEO-002 um das erneute
  Eintreten ergänzt (Stand-Spalte, Beleg slice-112); keine neue Beobachtung.
- **Folge-Slices:** keiner — die einzige offene Frage (Doppel-Dokument als
  eigener Befund) ist eine Produktfrage für einen Change Request, kein Slice.
- **Risiken aus §6:** alle mit Ausgang (§5) — eines eingetreten, drei entfallen.
- **Drei Paarungen** (wellenloser Slice, darum hier geprüft): Anker — der
  Steering-Loop-Eintrag trägt seinen Herkunfts-Anker ([`MR-025`](../../../../harness/conventions.md#mr-025)); Folge-Slice —
  keiner, also keine Datei in `open/`; Register — BEO-002 zitiert, nicht neu
  formuliert.
