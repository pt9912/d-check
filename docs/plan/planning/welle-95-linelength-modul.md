# Welle welle-95: Modul `linelength` — höchstens N Zeichen je Zeile in Markdown-Dateien

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-11-22.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Ein Repo, das für seine Markdown-Dateien eine Zeilenlänge vorgibt, prüft sie mit
d-check: das neue opt-in-Modul `linelength` liest die Markdown-Dateien der
Scan-Menge (`scan.roots`/`scan.ignore`), zählt je Zeile die Zeichen und meldet jede
Zeile über der konfigurierten Obergrenze mit Datei, Zeile, gemessener Länge und
Grenze. Andere Dateiarten liest das Modul nicht — die rohe Größe beliebiger Dateien
bleibt beim Modul `file`. Der Lauf bleibt hermetisch: nur Dateizugriff, kein Netz,
kein git, keine Schreibzugriffe. Was als Zeichen zählt und welche Zeilen (Fenced-Code,
Tabellen, lange URLs) ausgenommen sind, entscheidet der erste Slice vor dem Code.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Der Auftraggeber hat die Zeilenlängen-Prüfung für Markdown-Dateien als Welle
  angefordert (2026-11-22, in der Sitzung; kein GitHub-Issue).
- `origin/main` und der lokale Hauptzweig stehen gleich (`git fetch`, Abstand 0/0),
  `next/` ist leer, `in-progress/` trägt nur die Roadmap, `open/` nur die zwei Slices
  dieser Welle.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle sind `done/`.
- `make fullbuild` grün.
- **Das Mehr gegenüber den Slice-DoDs** (repo-weit, in keiner einzelnen DoD):
  - **Ohne den Block nichts Neues:** `make blackbox-probe` gegen das Vorher-Image
    der letzten Release-Version (`REF=v0.91.0`) zeigt für den Eigenlauf und ein
    Fixture ohne `linelength`-Konfiguration **dieselbe** Ausgabe (stdout, stderr und
    Exit getrennt). Ein neues opt-in-Modul darf keinen bisher grünen Lauf
    verschieben; jeder Slice belegt seine Hälfte, erst der gemeinsame Lauf belegt,
    dass keine die andere verschoben hat.
  - **Die Entscheide tragen am gebauten Image:** ein Fixture mit je einer Zeile
    knapp unter, genau auf und knapp über der Grenze, einer Zeile mit
    Mehrbyte-Zeichen, einer überlangen Zeile in Fenced-Code, in einer Tabelle und
    mit einer langen URL liefert am Image genau die Befunde, die die Anforderung
    zusagt (`make image-test`); die Tabelle Probe → Ergebnis steht in
    `welle-95-results.md`.
- Das Trigger-Audit (seit slice-241) ist über die Welle-Closure vollzogen und in
  `welle-95-results.md` belegt.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-289`](done/slice-289-linelength-anforderung-und-adr.md) | Anforderung, Spezifikation und ADR für das Modul `linelength`; Entscheide festschreiben | neue Anforderung des Moduls (Kennung vergibt der Slice), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) |
| [`slice-290`](next/slice-290-linelength-modulkern.md) | Modulkern: Konfiguration lesen, Zeilen messen, Ausnahmen anwenden, Grund-Code | die Anforderung aus [`slice-289`](done/slice-289-linelength-anforderung-und-adr.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine Welle.
- Wird blockiert von: keiner.
- **Reihenfolge der Slices ist fachlich bindend:** erst
  [`slice-289`](done/slice-289-linelength-anforderung-und-adr.md) (Anforderung und
  Entscheide stehen, bevor Code entsteht), dann
  [`slice-290`](next/slice-290-linelength-modulkern.md). Das WIP-Limit 1 hält sie
  ohnehin nacheinander.
- **Entscheide, die [`slice-289`](done/slice-289-linelength-anforderung-und-adr.md)
  vor dem Schreiben der Anforderung festschreibt** (mit Empfehlung; ohne Entscheid
  des Auftraggebers beginnt kein Code): Modulname und Kürzel, eigenes Modul statt
  Bedingung in `file` oder `structure`, Zähleinheit (Zeichen, Bytes oder
  Anzeigebreite), welche Zeilen ausgenommen sind und ob schaltbar, Name des
  Schwellen-Schlüssels (Abstand zu `max-lines`), Grund-Code und Befund-Ort.
- **Entschieden mit dem Auftraggeber (2026-11-22):** nur Markdown-Dateien; die acht
  Entscheide in [`slice-289`](done/slice-289-linelength-anforderung-und-adr.md) §3 mit den
  Empfehlungen bestätigt.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Andere Dateiarten** (Go, YAML, Skripte, …) — vom Auftraggeber ausgeschlossen
  (2026-11-22); die Anforderung nennt es in ihrem Out-of-Scope.
- **Automatisches Umbrechen oder ein Reparatur-Patch** — d-check ist ein Lese-Werkzeug
  ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
  ob `--fix-patch` das Modul später kennt, ist ein eigener Vorgang.
- **Das Modul im eigenen Profil scharf schalten** — die Repo-Doku führt bewusst
  Absätze und Tabellenzeilen auf einer Zeile; eine Grenze dort wäre eine Umformatierung
  des Bestands, kein Teil der Fähigkeit. Bleibt eine Entscheidung des Auftraggebers
  nach der Welle.
- **Release, Tag, Release-Prep** (CHANGELOG, READMEs, Handbuch-Kopf, Digest-Pins) —
  eigener Vorgang nach der Welle ([`AGENTS.md`](../../../AGENTS.md) §5 Zeile 17).
- **Die roten Nachtläufe vom 2026-11-22** (fünf veraltete Pins, behebbare CVEs im
  publizierten Image) — eigene Pin-Hebungen außerhalb der Welle; sie berühren das
  Modul nicht, gehören aber vor das nächste Release.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-95-results.md` (Geschwister im Ruheort `done/`, entsteht bei der Closure)
Zähler: [Beobachtungs-Register](observations/README.md)
