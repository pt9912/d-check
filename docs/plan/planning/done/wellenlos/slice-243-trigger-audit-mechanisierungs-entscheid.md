# slice-243: Mechanisierungs-Entscheid Trigger-Audit — je Klasse prüfbar oder Prosa

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle as State Machine.

**Welle:** ohne Welle — eine abgrenzte Evaluations-Entscheidung ist kein
repo-weites Mehr (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** die Adoptions-Kette [`MR-073`](../../../../../harness/conventions.md#mr-073)
(Lieferant des Deltas) und [`slice-241`](slice-241-trigger-audit-adoption.md)
(Adoption als Prosa-Schritt); Baseline-Regelwerk `modul-06-roadmap.md`
v6.13.0 §Closure („Erreicht dieselbe Fehlerklasse trotzdem ein **viertes**
Mal die Schwelle, gilt die Prosa-Form als ausgeschöpft: Der neue
Steering-Loop-Eintrag benennt dann entweder einen mechanischen Sensor …
oder begründet explizit, warum keiner möglich ist").
Keine `DC-*` — falls die Evaluierent einen Produkt-Sensor ergibt, wird
dessen Anforderungs-Delta im Slice benannt, aber erst im Folge-Slice
vergeben.

**Berührte Spec-Stellen:** — *(die Evaluierung berührt kein Spec-Stratum;
eine produktive Konsequenz wäre Folge-Slice)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-29.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Mechanisierungs-Frage für den Trigger-Audit (slice-241
adoptierte ihn als Prosa-Schritt) ist **je Klasse beantwortet**: für jede
der vier Trigger-Klassen (Carveout · bootstrap-aware Gate · ADR · Hard
Rule) eine Bewertung „mechanisch prüfbar — Bauform X" oder „Prosa-Verbleib
— Begründung", und je nach Gesamtbild die Umsetzung (kleiner Sensor im
eigenen Profil) oder die dokumentierte Ablehnung. Der
Baseline-Precedent ist die vierte-Mal-Schwelle: **die Prosa-Form ist nicht
ausgeschöpft** (kein wiederholter Audit-Fehler) — die Evaluierung stellt
die Bauformen bereit und entscheidet je Klasse, ob Verkörperung *jetzt*
sinnvoll ist oder bei Bedarf folgt; beides ist ein zulässiger Ausgang.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Produkt-Sensor ohne Anforderung** — ergibt die Bewertung einen
  neuen d-check-Modul-/Config-Umfang, trägt ihn ein **Folge-Slice** mit
  `DC-*`/ADR; dieser Slice trifft nur die Bauform-Entscheidung.
- **Keine Automatisierung des Urteils** — „Re-Evaluierungs-Trigger
  ausgelöst" ist Inhaltstatsache; mechanisch prüfbar ist höchstens der
  **Zustand** (z. B. offene Carveouts, Trigger-Zeilen-Form), nie die
  Urteilsbedeutung.
- **Keine Retro-Prüfung des Bestands** — der Audit-Horizont bleibt die
  laufende Closure (Abgrenzung von slice-241, unangetastet).

## 2. Definition of Done

- [x] Je Klasse eine belegte Zeile „mechanisch prüfbar (Bauform) /
      Prosa-Verbleib (Begründung)" — vier Zeilen, jede mit Gegenstand
      (Datei/Config-Mechanismus, nicht bloße Behauptung).
- [x] Die Entscheidung ist getragen: bei „mechanisch prüfbar" — Umsetzung
      im Slice oder benannter Folge-Slice mit Anforderungs-Umris; bei
      „Prosa-Verbleib" — Begründung, die gegen die
      vierte-Mal-Schwelle argumentiert.
- [x] `make gates` grün.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Bewertonsträger (Closure-Notiz dieses Slices; bei „Sensor" zusätzlich Folge-Slice-Plan) | update/neu | die vier Zeilen + Entscheidung |
| dieser Plan | update | etwaige Anker-Korrekturen, falls der Vollzug Zeilen verschiebt |

## 4. Trigger

**Start** (`open` → `in-progress`): direkt beansprucht — die
Auftraggeber-Anfrage (2026-09-29) ist die Beanspruchung selbst; der
Nachtlauf-Stand wird bei der Beanspruchung gelesen
([`MR-053`](../../../../../harness/conventions.md#mr-053)).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): die Bewertung ergibt für mehr als eine
  Klasse eine Produkt-Bauform, die je einen Folge-Slice mit eigenem
  Anforderungs-Delta trägt — dann wird die Umsetzung geschnitten und
  dieser Slice schließt mit den vier Bewertungen.
- `in-progress` → `open` (Blocker): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag (drei Paarungen
wellenlos hier geprüft).

## 6. Risiken und offene Punkte

- Vorzeitige Mechanisierung gegen die vierte-Mal-Schwelle — der Baseline-
  Wortlaut verlangt den Sensor erst bei ausgeschöpfter Prosa; eine
  Zwangs-Begründung „warum keiner möglich ist" ist **nicht** gefordert,
  solange die Prosa trägt. — **Ausgang:** *(offen)*
- Die Bewertung verwechselt Zustands-Prüfung (mechanisch) mit
  Urteils-Bedeutung (nicht mechanisch) — DoD 1 verlangt je Zeile den
  Gegenstand, nicht die Behauptung. — **Ausgang:** *(offen)*

## 7. Closure-Notiz

- **Was hat funktioniert:** Die vier Zeilen gegen den Baum bewertet, nicht
  gegen die Beschreibung — der Carveout-Bestand (`carveouts/done/` CO-001,
  CO-002; kein offener), die Schwellen-Form (kalibrierte Konstanten),
  die 76 Re-Evaluierungs-Trigger-Sektionen, die §3-Trigger-Zeilen.
- **Was ging anders als geplant:** — *(der Evaluierungs-Vollzug traf die
  Plan-Lage: keine Überraschung, keine Mechanisierungs-Lücke gefunden)*
- **Bewertung je Klasse (DoD 1):**
  - **Carveout — Prosa-Verbleib.** Prüfbar wäre der Zustand
    (Verzeichnis-Position, gemessen: `carveouts/done/` only); die
    Audit-Frage „Auflösungs-Trigger eingetreten?" steht im Carveout-Körper
    als je einzigartiges Konditional. Ein Gate auf „keine offenen
    Carveouts" kriminalisiert legitime Zustände.
  - **bootstrap-aware Gate — n.a.** d-checks Schwellen sind kalibrierte
    Konstanten ohne Hochschalt-Trigger (V-Beleg slice-241); nichts zu
    mechanisieren.
  - **ADR — Prosa-Verbleib.** 76 Sektionen über 96 ADR-Dateien
    ([ADR-0001](../../../../../docs/plan/adr/0001-implementierungssprache.md)
    bis
    [ADR-0096](../../../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md),
    ohne Lücken) plus README-Index tragen eine
    `## Re-Evaluierungs-Trigger`-Zeile (gemessen: `grep -l` über
    `docs/plan/adr/`); die Konditionale sind je einzigartig (Stichprobe
    [ADR-0048](../../../../../docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md),
    [ADR-0066](../../../../../docs/plan/adr/0066-cve-scan-gegen-das-publizierte-image.md),
    [ADR-0072](../../../../../docs/plan/adr/0072-workflows-modul.md) —
    R1-Report); „eingetreten?“ ist Urteil.
    Die Linkpflicht (ids) deckt die Erreichbarkeit, nicht die
    Bedeutungsfälle. *(R1-F-1: die erste Fassung nannte „97 Konditionale“
    — das war die .md-Gesamtmenge inkl. README-Index, keine
    Trigger-Zählung; V-2-Korrektur des Verifiers.)*
  - **Hard Rule — Prosa-Verbleib.** Existenz einer Trigger-Zeile wäre
    structure-Form prüfbar — kriminalisiert aber den grandfathered
    Alt-Bestand (§3.1/§3.5/§3.7/§3.9) und beantwortet doch nicht
    „ausgelöst?".
- **Steering-Loop-Eintrag:** keiner — die vierte-Mal-Schwelle ist nicht
  erreicht (kein Audit-Fehler bisher); der Prosa-Verbleib argumentiert
  exakt damit (DoD 2).
- **Beobachtungs-Register (`../../observations/`):** kein neuer Eintrag — der
  Audit-Schritt trägt seinen Vollzugs-Beleg in der Closure-Notiz
  (slice-241-Muster).
- **Folge-Slices:** keine — bei der ersten Audit-Auffälligkeit kehrt die
  Frage zurück (dann mit belegtem Anlass, gegen die
  vierte-Mal-Schwelle).
- **Risiken aus §6:** R1 (vorzeitige Mechanisierung) — abgewehrt, der
  Prosa-Verbleib argumentiert gegen die Schwelle; R2
  (Zustands- vs. Urteils-Prüfung) — je Zeile getrennt belegt.
- **Drei Paarungen:** Lerneintrag „der Audit hat vier Klassen, keine
  mechanisierbar ohne Anlass" — Folge-Slice: keiner — Register: kein
  neuer Eintrag.
- **Nachtlauf-Stand bei der Beanspruchung** ([`MR-053`](../../../../../harness/conventions.md#mr-053)):
  beide Nachtläufe grün (upstream-drift 2026-09-29, image-scan
  2026-09-28).

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das
Harness-Werkzeug selbst (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** keine berührte Sub-Area
mit Treffern in offenen Beobachtungen; die benachbarten stehenden Einträge
decken Pin-Hebungen und Register-Ausgänge, nicht die Mechanisierungs-Frage.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung):
`make nightly-state` wird bei der Beanspruchung gelesen und der Stand in
§7 notiert.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.

### Sub-Area: `*` (Harness-Werkzeug)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — die Bewertungs-Achsen sind durch die
  Baseline vorgeben (mechanische Form seit slice-122 baubar:
  `versions.patterns`), die vier Klassen sind terminologisch fest.
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — Evaluierung mit dokumentierter
  Entscheidung; der irrende Teil ist rückholbar.
- **Reconciliation-Aufwand:** Keiner — ein produktiver Ausgang trägt sein
  Anforderungs-Delta als Folge-Slice.
