# `make test` — fährt die Go-Testsuite des Hauptmoduls, inklusive zweier Zusagen des Repos über sich selbst

## Vertrag

Die Akzeptanzkriterien der bezogenen `DC-FA-*` liegen als Tests vor und laufen
grün — darunter der **Determinismus-Test** ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)),
eine Produkteigenschaft. Dazu **zwei** Zusagen, die keinen Produktcode prüfen,
sondern **Aussagen des Repos über sich selbst**:

- **Netzlos-Modullisten-Integrität** der [`.d-check.yml`](../../.d-check.yml)
  als getippter Go-Test statt als Skript-Vergleich
  ([ADR-0032](../../docs/plan/adr/0032-gate-consistency-tombstone.md)): die in
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  zugesagte Modulliste deckt sich mit der konfigurierten.
- **Wohlgeformtheit der Durchsetzungs-Konfiguration**
  ([`MR-048`](../conventions/MR-048-gate-ueber-werkzeug-datei.md)):
  `.claude/settings.json` ist gültiges JSON, jeder verdrahtete Hook-Pfad
  existiert, der Befehls-Wächter hängt am `Bash`-Werkzeug, und die
  Permission-Sperrliste deckt jeden Namen der Wächter-Sperrliste als **ganze**
  Befehlsklasse.

## Grenze — was das Grün nicht abdeckt

Die zweite Zusage ist die mit den meisten Löchern, und sie sind gewollt:

1. **Kein Beleg, dass die Durchsetzung läuft** — nur dafür, dass die Dateien
   wohlgeformt sind und aufeinander zeigen. Permanent: ob ein Hook im
   laufenden Werkzeug greift, ist von außen nicht prüfbar.
2. **Fehlende Dateien werden übersprungen, nicht eingefordert** — ein Repo-Gate
   über eine **Werkzeug**-Datei prüft Wohlgeformtheit, nicht Anwesenheit
   ([`MR-048`](../conventions/MR-048-gate-ueber-werkzeug-datei.md)).
   **Nicht permanent:** der Eintrag führt dafür einen Auflösungs-Trigger —
   nimmt der Kanon die Werkzeug-Artefakte in die Bindepunkte auf, sind sie
   Repo-Pflicht und die Skip-Hälfte entfällt. Bis dahin gilt:
   die Werkzeug-Wahl ist keine Repo-Invariante.
3. **Die Hook-Pfad-Prüfung ist eine Teilketten-Suche, keine Shell-Analyse** —
   ein Pfad, der in einer Kommandokette steckt, wird gefunden; was die Kette
   damit tut, nicht. Heilbar nur mit einem Shell-Parser, und der wäre teurer
   als die Zusage wert ist.
4. **`.claude/settings.local.json` ist nicht Gegenstand** — sie ist
   klon-lokal und gehört keinem Repo-Vertrag an. Permanent.

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`make test` nennt die gelaufenen Pakete. Der Lauf deckt das **Hauptmodul**;
`tools/archive-wave/` führt ein eigenes `go.mod` und läuft unter
`make archive-wave-test`.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | alle Tests grün |
| 1 | mindestens ein Test rot — die Ausgabe nennt Paket und Fall |

**Das sind die Codes des Produkts, nicht die des Targets.** `make test` ist ein
Docker-Build-Recipe; GNU Make normalisiert jeden fehlgeschlagenen Recipe auf
seinen eigenen Exit 2 — über das Target ist die 1 nicht beobachtbar. Welcher
Fall vorliegt, sagt die **Ausgabe**.

## Bindung

Bestandteil von `make gates`.
[`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) ·
[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) ·
[ADR-0032](../../docs/plan/adr/0032-gate-consistency-tombstone.md) ·
[`MR-048`](../conventions/MR-048-gate-ueber-werkzeug-datei.md)
