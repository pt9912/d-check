# Slice slice-124: `diagrams` — Ventil-Parität und die fehlenden Schema-Zeilen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-82-config-flaechen](welle-82-config-flaechen.md) (zugeordnet
bei der Eröffnung).

**Bezug:**
[`DC-FA-DIAG-001`](../../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in)
(die zu erweiternde Anforderung),
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
und [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(die beiden Module, die den Zeilen-Marker heute honorieren — die Ziel-Form der
Parität), [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`MR-025`](../../../../harness/conventions.md#mr-025).

**Berührte Spec-Stellen:**
[`DC-FA-DIAG-001.a`](../../../../spec/spezifikation.md#dc-fa-diag-001a--kennungs-konsistenz-in-diagramm-fences-diagrams)
(Algorithmus) und das §2-Schema — dort trägt das Modul heute **gar keine Zeilen**.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

`diagrams` ist das einzige Modul ohne jedes Ventil: weder `exempt-paths` noch
den Zeilen-Marker, den `codepaths` und `ids` honorieren. Wer es aktiviert, hat
nur den modul-lokalen Scan-Scope — ein Beispiel-Diagramm mit erfundener
Kennung in einem Report blockiert sonst über den `pre-commit`-Hook jeden
Commit. Gemessen ist das in welle-80: das eigene Profil musste gescopt werden,
um genau das zu umgehen. Dieser Slice stellt die **Parität** her und trägt
zugleich nach, was seit der Einführung fehlt: die Schlüssel des Moduls stehen
nur im Algorithmus-Abschnitt, **nicht** im §2-Konfigurations-Schema — anders
als bei jedem anderen Modul.

## 2. Vorgehen

1. **CR-Commit zuerst:** Lastenheft
   [`DC-FA-DIAG-001`](../../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in)
   um die beiden Ventile
   erweitern — Beschreibung, Akzeptanzkriterien (Datei per Glob ausgenommen,
   Zeile per Marker ausgenommen, Default byte-identisch ohne beide,
   fail-closed-Rand), §7-Historie.
2. **ADR der Welle** um die dritte Entscheidung ergänzen: Ventil-Parität als
   Prinzip (ein Modul, das Befunde an Zeilen hängt, braucht ein Zeilen-Ventil),
   und die Abgrenzung zum Scope (Scope entfernt die Datei aus der Prüfung, das
   Ventil nur die Referenz).
3. **Spezifikation:** die Ventile im Algorithmus **und** die vollständigen
   §2-Schema-Zeilen des Moduls (`fences`, `patterns[].regex`,
   `patterns[].defined-in`, `scope`, plus die neuen) — die Lücke wird mit
   geschlossen, weil sie derselbe Vertrag ist.
4. **Code + Tests:** Ventile im Modul, Config-Rand; Tests für beide Ventile,
   für ihre Nicht-Wirkung ohne Angabe und für den Fall, der welle-80 zum
   Scoping zwang.
5. **Messen, nicht annehmen:** ob das eigene Profil danach ohne Scope auskommt,
   wird gemessen — die Umstellung selbst ist **nicht** Teil dieses Slice.
6. **Spiegel** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
   Anforderung, Algorithmus, §2-Schema, `--print-config`-Vorlage,
   Config-Kommentar (Handbuch ist Release-Prep — dort steht seit slice-121 die
   Aussage „diagrams hat kein Ventil", die mit dem Release fällt).
7. Unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Das eigene Profil bleibt gescopt** — der Rückbau ist ein eigener Entscheid
  nach der Messung.
- **Kein Handbuch, kein CHANGELOG** (slice-125) — aber der Handbuch-Satz aus
  slice-121 ist dort ausdrücklich zu korrigieren.
- **Keine Default-Änderung.**

## 4. Definition of Done

- [x] CR-Commit (Lastenheft allein) liegt **vor** Spezifikation und Code —
      Lastenheft 0.65.0 allein, danach ADR, dann Spezifikation samt Code.
- [x] Beide Ventile implementiert und getestet; §2-Schema trägt alle Schlüssel
      des Moduls — **mit einer Korrektur aus dem Review:** die
      `scope`-Zeile ist wieder gestrichen, weil die generischen
      `<modul>.scope.*`-Zeilen sie decken und eine eigene die zweite
      Pflegestelle wäre. „Alle Schlüssel" heißt also: alle **modul-eigenen**.
- [x] Default-Beweis byte-identisch — **mit einer benannten Grenze:** nur das
      Datei-Ventil hängt an einem Schlüssel, der Marker hängt am Inhalt. Für
      die Marker-Hälfte gilt die Byte-Identität nur für Bäume ohne die
      Zeichenfolge in einer gelisteten Fence (gemessen: identische Config,
      Vorgänger zwei Befunde, neuer Stand einer).
- [x] Der Fall, der zum Scoping zwang, ist als Test gepinnt — ein
      Beispiel-Diagramm mit **drei** erfundenen Kennungen, freigeschaltet über
      **eine** Marke auf der Öffnungszeile.
- [x] `make gates` grün; unabhängiger Review (kein HIGH, sieben MEDIUM, acht
      LOW, alle eingearbeitet, [Report](../../../reviews/2026-08-22-slice-124-diagrams-ventile-review.md));
      Closure-Notiz; Register gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Ein Ventil in einem Fence ist anders als in Prosa:** der Marker ist ein
  HTML-Kommentar — in einem `mermaid`-Fence ist er kein Kommentar, sondern
  Diagramm-Text. Die Zeilen-Semantik muss das aushalten oder die Grenze
  benennen. — **Ausgang:** aufgelöst, und die Auflösung ist die eigentliche
  Substanz des Slice: der Marker ist ein **Token**, kein HTML-Kommentar; wie
  er vor dem Renderer versteckt wird, ist Sache der Diagramm-Sprache. Dazu kam
  ein Ort, den der Plan nicht vorsah — die **Öffnungszeile** für den ganzen
  Block, weil die intuitive Platzierung sonst still nichts täte. Der Review
  hat die Kehrseite gefunden: die **schließende** Zeile ist kein Ventil-Ort
  und war still folgenlos; jetzt mit Negativtest und Halbsatz benannt.
- **Die §2-Lücke ist älter als dieser Slice** — sie zu schließen ist richtig,
  vergrößert aber den Diff über den Ventil-Kern hinaus. — **Ausgang:** richtig
  gewesen, und teurer als gedacht: von sechs neu geschriebenen Zeilen waren
  **drei** falsch (`scope`-Default in beiden Lesarten, `fences: []` fällt
  entgegen der Präambel auf den Default zurück, `defined-in` sagte „Datei",
  wo der Rand ein Verzeichnis durchließ). Eine Lücke zu schließen heißt, jede
  Zeile einzeln gegen den Rand zu messen — nicht plausibel zu formulieren.
  Die dritte hat Code geändert statt Text.
- **Parität heißt nicht Gleichheit:** was `codepaths`/`ids` können, muss hier
  nicht identisch heißen; die Namen folgen dem Bestand. — **Ausgang:**
  gehalten, aber die Bezugsgröße war falsch: `versions` trägt beide Ventile
  ebenfalls und ist die nächstliegende Präzedenz — Plan, Anforderung und
  Algorithmus nannten nur zwei Module. Korrigiert.

## 6. Trigger

**Start** (`open` → `in-progress`): Wellen-Eröffnung (eingetreten).

**Rückführungen:** `in-progress` → `next`, falls der Zeilen-Marker im Fence
nicht ausdrückbar ist (dann trägt nur das Datei-Ventil, und die Grenze wird
benannt).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Kern (GF), Config-Rand (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): **BEO-004**
  (ein Modul liest Eingaben, die es nicht scannt) ist einschlägig — das Modul
  liest seine `defined-in`-Quelle; BEO-002 als Spiegel-Pflicht,
  BEO-006/009/010 als Arbeitsregeln.

Slice-ID: slice-124. Betroffene IDs:
[`DC-FA-DIAG-001`](../../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus). Module:
`diagrams` (Kern `rules/`), Config-Rand, Spec. Gates: `make test` (eng),
`make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — additive Erweiterung plus Nachtrag einer
Vertrags-Lücke im eigenen Spec-Stratum.

## 9. Closure-Notiz (nach `done/`)

Geliefert sind beide Ventile und die §2-Schema-Zeilen, die dem Modul als
einzigem fehlten. Der Zeilen-Marker wirkt an zwei Orten — auf einer
Diagramm-Zeile für sie, auf der Öffnungszeile für den ganzen Block.

**Die Lehre ist eine Reichweiten-Lehre, und sie ist in dieser Welle die
dritte.** Ich habe geschrieben, `diagrams` sei das einzige Modul, das Befunde
an Zeilen hängt und kein Ventil trägt. Drei weitere tun dasselbe
(`hostpaths`, `pins`, `spans`). In slice-123 war es „als einzige" über eine
Ausnahme, in slice-122 „zwei Paare, zwei Befunde", hier „das einzige Modul".
Jedes Mal war die Aussage aus dem **Anlass** abgeleitet statt aus dem
**Bestand** — und jedes Mal hätte ein `grep` über die Nachbarn genügt. Wer
eine Eigenschaft für exklusiv erklärt, muss die Menge gezählt haben.

**Zwei Grenzen sind erst durch den Review benannt worden**, beide echt: die
Ventile wirken **scan-seitig** und nicht auf der Definitionsmenge (eine im
`defined-in` als illustrativ markierte Kennung definiert weiter), und die
Marker-Hälfte ist **nicht** opt-in — sie hängt am Zeilen-Inhalt, nicht an
einem Schlüssel. Die zweite widerlegt eine Zusage des Wellendokuments für die
Hälfte dieser Erweiterung; sie steht jetzt dort als benannte Grenze.

**Die Messung aus §2.5 sagt weniger, als sie zu sagen schien.** „Ohne Scope,
446 Dateien, null Befunde" ist eine **Bestandsaufnahme** — heute trägt kein
Diagramm außerhalb `spec/` ein `ARC-\d{3}`-Token, und das galt vorher genauso.
Sie kann „die Ventile helfen" nicht von „es hat nie etwas gefeuert" trennen.
Den Wellen-Closure-Trigger erfüllt die konstruierte Gegenprobe (dieselbe
Konfiguration mit `exempt-paths` weist das Vorgänger-Image mit Exit 2 zurück),
nicht diese Zahl. Der Rückbau des Scopes bleibt ein eigener Entscheid.
