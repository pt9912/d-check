# Slice slice-095: `links.resolve-from` — Auflösung unabhängig vom Lifecycle-Verzeichnis

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** bei Start zu eröffnen. Ein Slice „ohne Welle" ist in diesem Repo
nicht einlösbar: `make planning-check` koppelt Ruhe-Marker und `in-progress/`
atomar ([`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform)),
ein Slice in Arbeit verlangt also eine benannte Welle — auch wenn er inhaltlich
keine bräuchte.

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
Closure von [slice-093](welle-68/slice-093-closure-note-gate.md) am 2026-08-09
ist die Klasse **zweimal** eingetreten:

1. Die Links der Review-Reports auf den Slice zeigten nach `in-progress/` und
   brachen mit dem Move.
2. Das Wellendokument verwies auf `done/slice-09….md`; als es selbst nach
   `done/` wanderte, brachen seine eigenen Zeiger.

Und ein drittes Mal, im selben Maßstab: die Eröffnung von
[welle-69](welle-69/welle-69-structure-schnitt.md) am 2026-08-09 verschob **einen**
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
   — **Bestätigt** (beim Wellen-Start; begleitende ADR trägt den Entscheid).
2. **Befund-Form.** Ein eigener Grund-Code (Vorschlag `link-position-dependent`)
   statt `target-missing` — die Reparatur ist eine andere (Pfad präfixieren, nicht
   Ziel anlegen), und am Ist-Ort ist nichts kaputt.
   — **Bestätigt**, mit einer Schärfung aus der Messung: positionsabhängig ist
   ein Verweis auch dann, wenn er von jedem Ort auflöst, aber auf
   **verschiedene** Ziele — er meint dann je nach Ort etwas anderes. Beide
   Fehlarten tragen denselben Code; die Meldung nennt den Ort bzw. die Ziele.
3. **Zählt das Bild-Ziel mit, und was ist mit Ankern?** Vorschlag: dieselbe
   Ziel-Menge wie die bestehende Auflösung, Anker bleiben außen vor (die
   Anker-Prüfung hängt am Ziel-Dokument, nicht am Quellort).
   — **Bestätigt.**

## 3a. Messung: Quell-Menge, Bestand und der Retro-Beleg

Näherungs-Messung über den eigenen Planungs-Baum (Entwurfs-Entscheidung; die
finale Zahl liefert das Produkt nach der Implementierung):

| Quell-Menge | relative Verweise | positionsabhängig |
|---|---|---|
| alle vier Lifecycle-Verzeichnisse | 1875 | **108** — die Spitzenreiter sind Wellendokumente und Slices in `done/` |
| nur die **wandernden** (`open/`, `next/`, `in-progress/`) | 79 | **0** |

**Zwei Entscheide folgen daraus:**

1. **Quellen sind nur die wandernden Verzeichnisse.** Eine Datei in `done/` ist
   am Endzustand ihres Lifecycles — sie wandert nie wieder, ihre Verweise müssen
   nur vom Ist-Ort auflösen. Ohne diese Einschränkung wären 108 der Befunde
   Falsch-Positive auf ortsfesten Dokumenten.
2. **Der Bestand ist heute grün — und das ist das Ergebnis, nicht die
   Ausgangslage.** Die Null ist von einer Woche Hand-Nachzügen gekauft: allein
   heute wurden bei drei Lifecycle-Moves 10, 15 und 14 Verweise nachgezogen,
   jedes Mal **nach** dem Move gemeldet. Der Wert der Fähigkeit ist Prävention;
   der Beleg dafür ist retro.

**Der Retro-Beleg trifft die Slice-Hälfte des historischen Bruchs — und die
Zahlengleichheit ist Zufall, kein Beweis.** Dieselbe Messung gegen den Stand
**vor** der welle-69-Eröffnung: 76 relative Verweise in wandernden Dateien,
davon **19 positionsabhängig** (der später verschobene Slice trug 7, die vier
zurückbleibenden Geschwister den Rest — dieser Slice selbst 2). Der reale
19-Link-Bruch setzte sich anders zusammen: 15 dieser Verweise überlappen, die
übrigen vier des Bruchs waren Links eines **Review-Reports** auf den
verschobenen Slice — die Quelle ortsfest, das **Ziel** gewandert. **Diese
Klasse deckt die Fähigkeit strukturell nicht** (sie prüft hypothetische
Quell-, nicht Ziel-Orte) und deckt sie bewusst nicht: für Review-Reports als
Lauf-Belege ist `ignore-refs` das etablierte Ventil, und lebende Dokumente
zieht die Move-Regel im selben Commit nach. Die Grenze steht in der ADR.

## 4. Definition of Done

- [x] [`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
      um `resolve-from` erweitert (Akzeptanzkriterien inkl. „ohne Block
      byte-identisch"; Lastenheft 0.60.2), Algorithmus als Schritt 6 der
      Spezifikation, begleitende
      [ADR-0056](../../adr/0056-resolve-from-wandernde-quellorte.md).
- [x] Implementierung + Tests; **Realdatenbeleg im eigenen Repo**: die
      **Quellort-Hälfte** der belegten Fälle wäre vor dem Move rot gewesen
      (Retro-Lauf mit dem Produkt: 19 Befunde am Vor-welle-69-Stand); die
      Ziel-Wanderungs-Hälfte ist als Grenze in der ADR benannt.
- [x] `make gates` grün; Release **v0.60.0** als Minor (Digest
      `sha256:5892a87b…d3f9`), Fähigkeit im eigenen Baum scharf in `make gates`.

## 5. Risiken / offene Punkte

- **Falsch-Positive bei absichtlich ortsgebundenen Verweisen** (ein Slice, der
  bewusst auf seinen eigenen Ruheort zeigt). — **Ausgang:** das bestehende
  Ventil `ignore-refs` trägt den Fall (dieselbe Wirkung wie in Schritt 4,
  getestet); kein neues Ventil. Erweist es sich als zu grob, ist das ein
  benannter Re-Evaluierungs-Trigger der ADR.
- **Kombinatorik:** vier Quellorte × alle Verweise erhöht die Prüf-Menge.
  — **Ausgang:** gemessen, nicht behauptet: die Laufzeit des vollen Laufs liegt
  im Rauschen der Vorversion
  ([`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance) unberührt);
  kein Carveout nötig, keine Rückführung.

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe; WIP-Slot frei. Unabhängig von
[slice-094](../done/slice-094-closure-zaehl-paritaet.md) und
[slice-096](welle-69/slice-096-structure-modul-analyse.md) umsetzbar.

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

Geliefert am 2026-08-16 mit Release **v0.60.0**. `links.resolve-from` prüft
Dateien in wandernden Lifecycle-Verzeichnissen gegen jeden Ort ihrer Gruppe;
beide Fehlarten (nicht überall auflösbar, nicht überall dasselbe Ziel) tragen
den neuen Grund-Code `link-position-dependent`. Alle drei Abnahme-Punkte sind
mit ihren notierten Ausgängen umgesetzt; die Messung aus §3a hat den Zuschnitt
getragen (nur `dirs`-Dateien sind Quellen — sonst 108 Falsch-Positive) und der
Retro-Beleg lief mit dem Produkt (19 Befunde am Vor-welle-69-Stand, die
Quellort-Hälfte des realen Bruchs). Zwei unabhängige Review-Runden: die erste
blockierend (15 Befunde, darunter die Ist-Ort-Vorbedingung gegen Doppelbefunde
und der stille Quellen-Ausfall bei Tippfehler-Orten), die zweite APPROVE mit
Text-Auflagen (geheilt vor dem Release, Lastenheft 0.60.2). Dazu ein
**CI-Realfund**: der erste fail-closed-Zuschnitt meldete auf jedem frischen
Klon das legitim geleerte `open/` — git überträgt leere Verzeichnisse nicht;
die Zusage ist an dieser Realität justiert und die Rest-Grenze im Vertrag
benannt. Drei Grenzen sind ausdrücklich Vertrag, nicht Lücke: Ziel-Wanderung
([ADR-0056](../../adr/0056-resolve-from-wandernde-quellorte.md) Entscheidung 6), Scan-Bereichs-Kopplung der Gruppen-Orte und der
einzelne fehlende Ort. Diese Closure selbst war der erste Anwendungsfall: die
Fähigkeit lief scharf, während ihr eigener Slice nach `done/` wanderte.
Wellen-Kontext in [welle-76-results.md](welle-76-results.md).
