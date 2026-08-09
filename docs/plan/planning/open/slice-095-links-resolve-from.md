# Slice slice-095: `links.resolve-from` — Auflösung unabhängig vom Lifecycle-Verzeichnis

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** ohne Welle — einzeln lieferbar, keine Closure-Bedingung jenseits der DoD.

**Bezug:** [`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
(Erweiterung, kein neues Modul);
[`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)
(die Lifecycle-Move-Regel, deren Invariante hier maschinell wird).
**Change Request** aus dem Schwester-Repo a-check (CR 2 seiner
Werkzeug-Abdeckungs-Analyse).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Ein relativer Verweis in einer Slice-Datei muss aus **jedem**
Lifecycle-Verzeichnis auflösen — nicht nur vom aktuellen Ort. Die Option
`links.resolve-from` prüft zusätzliche, hypothetische Quellorte; ohne sie bleibt
das Verhalten byte-identisch.

## 2. Warum das eine echte Invariante ist

Zwei Regeln zusammen erzeugen sie: der Lifecycle ist eine Zustandsmaschine über
Verzeichnisse, und der Wechsel ist ein `git mv` **ohne** Inhaltsänderung. Ein
präfixloser Nachbar-Verweis ist am Ist-Ort grün und bricht beim nächsten
Wechsel — sichtbar erst dann, wenn man ihn nicht mehr reparieren darf, ohne die
Move-Regel zu verletzen.

**Der Beleg liegt im eigenen Repo, nicht nur beim Antragsteller.** Bei der
Closure von [slice-093](../done/slice-093-closure-note-gate.md) am 2026-08-09
ist die Klasse **zweimal** eingetreten:

1. Die Links der Review-Reports auf den Slice zeigten nach `in-progress/` und
   brachen mit dem Move.
2. Das Wellendokument verwies auf `done/slice-09….md`; als es selbst nach
   `done/` wanderte, brachen seine eigenen Zeiger.

Und ein drittes Mal, im selben Maßstab: die Eröffnung von
[welle-69](../done/welle-69-structure-schnitt.md) am 2026-08-09 verschob **einen**
Slice von `open/` nach `in-progress/` — und brach damit **19 Links** auf einen
Schlag. Betroffen waren die vier Nachbar-Slices in `open/` (präfixlose
Geschwister-Verweise), ein Review-Report und der verschobene Slice selbst, dessen
Verweise auf die zurückgebliebenen Geschwister nun ins Leere zeigten.

Alle drei Male wurden die Verweise von Hand nachgezogen, weil `doc-check` sie
**nach** dem Move meldete. Genau das soll vorher auffallen — und die 19 zeigen,
dass die Klasse nicht mit der Repo-Größe skaliert, sondern mit der Zahl der
Nachbarn.

## 3. Abnahme-Punkte

1. **Erweiterung statt neues Kürzel.** Nach dem Kriterium aus
   [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md)
   (querschnittlich → neues Kürzel, Einzelmodul → bestehende Anforderung ändern)
   gehört das in
   [`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links):
   dieselbe Prüfung mit erweiterter Quellort-Menge, keine neue Frage.
   **Zu bestätigen.**
2. **Befund-Form.** Ein eigener Grund-Code (Vorschlag `link-position-dependent`)
   statt `target-missing` — die Reparatur ist eine andere (Pfad präfixieren, nicht
   Ziel anlegen), und am Ist-Ort ist nichts kaputt.
3. **Zählt das Bild-Ziel mit, und was ist mit Ankern?** Vorschlag: dieselbe
   Ziel-Menge wie die bestehende Auflösung, Anker bleiben außen vor (die
   Anker-Prüfung hängt am Ziel-Dokument, nicht am Quellort).

## 4. Definition of Done

- [ ] [`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
      um `resolve-from` erweitert (Akzeptanzkriterien inkl. „ohne Block
      byte-identisch"), Algorithmus und Grund-Code in der Spezifikation,
      begleitende ADR.
- [ ] Implementierung + Tests; **Realdatenbeleg im eigenen Repo**: die beiden
      oben belegten Fälle wären vor dem Move rot gewesen.
- [ ] `make gates` grün; Release als Minor (d-check findet danach mehr).

## 5. Risiken / offene Punkte

- **Falsch-Positive bei absichtlich ortsgebundenen Verweisen** (ein Slice, der
  bewusst auf seinen eigenen Ruheort zeigt). — **Ausgang:** offen; das
  bestehende Ventil `ignore-refs` ist der Kandidat, bevor ein neues erfunden wird.
- **Kombinatorik:** vier Quellorte × alle Verweise erhöht die Prüf-Menge.
  — **Ausgang:** offen; zu messen, ob die Laufzeit-Zusage
  ([`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance)) berührt wird.

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei. Unabhängig von
[slice-094](slice-094-closure-zaehl-paritaet.md) und
[slice-096](../done/slice-096-structure-modul-analyse.md) umsetzbar.

**Rückführungen:** `in-progress` → `open`, falls die Laufzeit-Messung ein
Carveout nötig macht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001**
  (Datei-Register driften unbemerkt gegen ihre Autoritäts-Tabelle). Verwandt,
  aber verschieden: BEO-001 fragt „wird diese Datei **irgendwo** referenziert?",
  dieser Slice fragt „löst dieser Verweis **von überall** auf?". Nichts zu
  berücksichtigen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Change Request und ADR schreiben die Zusage,
die Implementierung liefert sie; spec-first wie jeder d-check-Slice.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
