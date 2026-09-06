# Slice slice-206: Modul `mentions` — Implementierung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`DC-FA-MENT-001`](../../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in)
(Anforderung), [ADR-0084](../../adr/0084-mentions-eigenes-modul.md) (Bauform
und Verdikt-Semantik). Beide entstanden in slice-205 <!-- d-check:status-provenance -->; dieser Slice
liefert, was dort ausdrücklich aus dem Umfang genommen wurde.

**Berührte Spec-Stellen:**
[`DC-FA-MENT-001`](../../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in)
— **keine Änderung**; der Slice setzt sie um. Berührt werden zusätzlich
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
(die Modulliste nennt `mentions` bereits) und
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(neuer Konfigurations-Block).

**Verantwortlich:** — · **Autor:** pt9912. **Datum:** 2026-09-06.

---

## 1. Ziel

Das Modul `mentions` bauen: Soll-Menge aus Pfad-Globs gegen Ist-Menge aus
Dokumenten, Befund `artifact-unmentioned` für jedes Mitglied, das in keinem
Dokument vorkommt. Umfang und Verdikt-Semantik sind entschieden; dieser Slice
verhandelt sie nicht neu.

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Kern-Regel + Modul-Registrierung | neu | die Prüfung selbst, samt `validModules()`-Eintrag |
| Konfigurations-Schema | update | `mentions.artifacts`, `mentions.documents`, `mentions.match` |
| Tests | neu | das Akzeptanzkriterien-Trio der Anforderung als Fitness Function |
| [`docs/user/`](../../../user/) | update | Handbuch und Betriebsdoku nennen das neue Modul |

**Der Kalibrierungs-Auftrag ist der eigentliche Gehalt.** slice-205 schließt
mit dem Befund, dass d-checks eigene Gegenprobe auf der richtigen Menge grün
läuft — es gibt **kein eigenes Rauschen**, an dem sich Schwellen und
Ausnahme-Klassen justieren ließen. Dieser Slice muss deshalb an einem
**Fremd-Bestand** messen und das Ergebnis als solches ausweisen. Die drei in
slice-205 benannten Mengen sind der Startpunkt, nicht das Ziel.

## 3. Ausdrücklich NICHT in diesem Slice

- **Eine zweite Quell-Form der Soll-Menge** (Literale statt Pfade). Out-of-Scope
  (2) der Anforderung; sie braucht ihren eigenen Beleg.
- **Die Aufnahme in `make gates`.** Eine neue Modul-Klasse startet als
  eigenständiger Fokus-Lauf — dieselbe Vorsicht wie bei `review-coverage` und
  `trace-check`. Die Aufnahme ist ein späterer, eigener Entscheid.
- **Jede Änderung an [`DC-FA-MENT-001`](../../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) oder [ADR-0084](../../adr/0084-mentions-eigenes-modul.md).** Zeigt die
  Implementierung, dass eine Festlegung nicht trägt, ist das ein Folge-Entscheid
  und keine stille Anpassung.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Das Modul `mentions` existiert, ist opt-in und erfüllt alle
      Akzeptanzkriterien von
      [`DC-FA-MENT-001`](../../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in)
      als Tests — einschließlich beider fail-closed-Fälle und des
      byte-identischen Default-Laufs.
- [ ] **(2)** Ein Fokus-Target fährt es, und die Doku (`docs/user/`, Sensors-
      Tabelle) nennt es — mit `kein Gate` in der Zeile, solange es nicht in
      `gates` läuft.
- [ ] **(3)** Ein **Kalibrierungs-Beleg** an einem Fremd-Bestand liegt vor: je
      Mengen-Wahl die Zahl der Funde **und** das Urteil, wie viele davon Mängel
      sind. Ohne diesen Beleg ist die Ausnahme-Klasse geraten.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §5 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Befund-Ort ist festgelegt, aber ungewohnt.** `file` ist der
  Artefakt-Pfad, `line` der Platzhalter `1` — das Artefakt wird nie geöffnet.
  Sechs Module führen `Line: 1` bereits als Platzhalter; ob die Kombination aus
  nicht geöffneter Datei **und** Platzhalter-Zeile in `--doctor` und `--repair`
  sinnvoll erscheint, ist ungeprüft. — **Ausgang:** <offen>
- **Die Bezugsmenge muss in zwei Ausgabe-Formen erscheinen.** stderr im
  Default, `summary`-Felder unter `--json`/`--yaml`. Die zweite Hälfte
  erweitert ein Struct, das heute zwei Felder trägt; ob das rückwärtskompatibel
  bleibt, entscheidet sich an den vorhandenen Konsumenten. — **Ausgang:** <offen>
- **Kein eigenes Rauschen — geerbt aus slice-205.** Die Ausnahme-Klasse wird am
  Fremd-Bestand justiert, nicht am eigenen. Das weicht von der gelebten Praxis
  ab und ist der Grund für DoD (3). — **Ausgang:** <offen>

## 6. Trigger

**Start** (`open` → `in-progress`): slice-205 liegt in `done/`, Anforderung und
ADR stehen. WIP-Limit frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt sich, dass die Kalibrierung am
  Fremd-Bestand eine eigene Mess-Arbeit ist, wird sie ein eigener Slice vor der
  Implementierung.
- `in-progress` → `open` (blockiert): Trägt eine Festlegung aus
  [ADR-0084](../../adr/0084-mentions-eigenes-modul.md) nicht, ruht der Slice
  bis zum Folge-Entscheid — stillschweigend angepasst wird sie nicht.

**Closure-Trigger.** Zwei beobachtbare Kriterien und ein Lerneintrag: (a) das
Fokus-Target läuft grün und die Akzeptanzkriterien stehen als Tests; (b) der
Kalibrierungs-Beleg aus DoD (3) liegt vor.

## 7. Vorgelagert (vor der Modus-Begründung)

<entsteht spätestens bei der Beanspruchung — ein Plan in `open/` trägt die drei
Vorprüfungen noch nicht>

## 8. Sub-Area-Modus-Begründung

<entsteht mit den Vorprüfungen bei der Beanspruchung>

## 9. Closure-Notiz (nach `done/`)

<wird vor dem `git mv` nach `done/` gefüllt>
