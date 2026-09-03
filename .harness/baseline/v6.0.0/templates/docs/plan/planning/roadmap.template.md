# Roadmap

> **Template-Hinweis.** Vorlage für die Roadmap des Repos. Kopiere nach
> `docs/plan/planning/in-progress/roadmap.md` (oder dem in deinem Repo
> üblichen Pfad) und ersetze Platzhalter. Lösche diesen Block.

**Format-Regel:** Die Roadmap ist eine Reihenfolge von **Wellen**,
keine Reihenfolge von Terminen (siehe
Baseline-Regelwerk `modul-06-roadmap.md`).
Termine werden — falls überhaupt — als Konsequenz der Wellen-Schätzung
gezeigt, nicht als Treiber.

---

## Offene Wellen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte — *Offene Wellen* ist **derivativ**: Der
Zustand sind die flachen Welle-Dateien; woran gearbeitet wird, sagt das
`Welle:`-Feld der Slices in `in-progress/`. Ziel, Trigger und
Closure-Kriterien stehen in der Welle-Datei, nicht hier.

- [<welle-NN-titel>](../<welle-NN-titel>.md)

<!-- BEDIENHINWEIS: Zwei unabhängige Aussagen in diesem Block. Die Liste oben
folgt den Dateien (ein Zeiger je offener Welle-Datei). Trägt in-progress/
keinen Slice, kommt der Ruhe-Marker ZUSÄTZLICH dazu — nicht an Stelle der
Liste; beides zugleich ist der Normalfall direkt nach der Wellen-Eröffnung.
Zwei Aussagen, zwei Wächter: die Marker-Hälfte (Marker genau dann, wenn
in-progress/ keinen Slice trägt) und die Listen-Hälfte (Bijektion Zeiger <->
flache Welle-Dateien, in beide Richtungen; der Marker geht nicht ein). Die
Listen-Hälfte braucht einen Sensor, der das Kardinalitäts-Modell kennt — ein
Ein-Wellen-Wächter hält den Block gegen GENAU EINE Datei und meldet legitime
Zustände (mehrere offene Wellen; eine Welle eröffnet, nichts beansprucht) als
Drift. Welche Hälfte dein Sensor prüft, musst du wissen; eine ungewächterte
Hälfte ist zulässig, wenn sie benannt ist.

WORTLAUT des Markers: Baseline-Regelwerk modul-06-roadmap.md, Bullet
"Offene Wellen". Hier bewusst NICHT zitiert, und das ist eine Regel, keine
Nachlässigkeit: Ein Doku-Sensor matcht den Marker als Substring dieses
Blocks, also matcht sich jeder Regel-, Hinweis- oder Beispieltext selbst,
der den Wortlaut literal trägt — der Block meldete "Ruhe" bei beanspruchtem
Slice. Aus demselben Grund gehört der Wortlaut in keine Sektions-Regel-Zeile
oben. Wer den Sensor selbst baut: Code-Fences beim Matchen aus dem Block
nehmen, sonst schlägt ein Beispiel-Auszug durch. -->

## Nächste Wellen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte, Bullet *Nächste Wellen* — die geordnete
Vorschau: je Zeile Welle, Trigger als beobachtbare Bedingung, wichtigste Slices
und geschätzter Aufwand (S/M/L, kein Termin).

| Welle | Trigger | Wichtigste Slices | Geschätzter Aufwand |
|---|---|---|---|
| <welle-N+1> | Welle <N> done | <…> | S/M/L |
| <welle-N+2> | Welle <N+1> done + ADR-<NNNN> accepted | <…> | S/M/L |

## Meilensteine

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Welle ≠ Meilenstein ≠ Release.

<!--
Externe Versprechen oder interne Trigger-Punkte.
"M2: erste lauffähige Fassung" ist ein Meilenstein.
Status: "offen" oder "erreicht YYYY-MM-DD" plus Beleg als auflösbarer
Anker (z. B. Tag, Workflow-Lauf, Ergebnis-Notiz). Erreichte Meilensteine
bleiben hier in der Tabelle —
sie gehören nicht ins Drift-Log, und die Status-Zelle erzählt nicht, wie es
dazu kam (Baseline-Regelwerk grundlagen-harness-dateien.md §Was ein
Kommentar trägt, "Dieselbe Regel für Zustandsfelder").
-->

| Meilenstein | Welle(n) | Trigger | Status |
|---|---|---|---|
| M1 | <welle-NN> | <…> | erreicht / offen |

## Abhängigkeitsgraph

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte, Bullet *Nächste Wellen* — die Abhängigkeit
steht als beobachtbare Bedingung in der `Trigger`-Spalte **und** als gerichtete
Kante hier; eine Welle, die ohne fertige Vorgängerin nicht starten kann, ist
eine Phantom-Welle.

```mermaid
flowchart LR
    W1[Welle 1]
    W2[Welle 2]
    W3[Welle 3]
    W4[Welle 4]
    
    W1 --> W2
    W1 --> W3
    W2 --> W4
    W3 --> W4
```

## Abgeschlossene Wellen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

| Welle | Abschluss | Closure-Notiz |
|---|---|---|
| <welle-NN> | YYYY-MM-DD | [`welle-NN-results.md`](../done/welle-NN-results.md) |

## Historische Trigger-Verschiebungen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte, Bullet *Historische Trigger-Verschiebungen*
— das Drift-Log: jede Umplanung mit Datum, Änderung, Grund. Leer heißt starre
Roadmap, jede Zeile voll heißt treibende.

<!--
Wenn Wellen umgeplant wurden: Datum, Grund, neue Reihenfolge.
Steering-Loop-relevant.
NUR Umplanungen: Trigger verschoben, präzisiert oder ersetzt; Slice oder
Welle umgehängt. KEINE Schließungen (die stehen im Closure-Log oben) und
KEINE erreichten Meilensteine (Status-Spalte) — sonst führt diese Tabelle ein
zweites Closure-Log, und zwei Logs driften.
-->

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| YYYY-MM-DD | <…> | <…> |
