# Slice slice-124: `diagrams` — Ventil-Parität und die fehlenden Schema-Zeilen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-82-config-flaechen](../welle-82-config-flaechen.md) (zugeordnet
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

- [ ] CR-Commit (Lastenheft allein) liegt **vor** Spezifikation und Code.
- [ ] Beide Ventile implementiert und getestet; §2-Schema trägt **alle**
      Schlüssel des Moduls.
- [ ] Default-Beweis byte-identisch; der welle-80-Fall ist als Test gepinnt.
- [ ] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Ein Ventil in einem Fence ist anders als in Prosa:** der Marker ist ein
  HTML-Kommentar — in einem `mermaid`-Fence ist er kein Kommentar, sondern
  Diagramm-Text. Die Zeilen-Semantik muss das aushalten oder die Grenze
  benennen. — **Ausgang:** *(bei Closure)*
- **Die §2-Lücke ist älter als dieser Slice** — sie zu schließen ist richtig,
  vergrößert aber den Diff über den Ventil-Kern hinaus. — **Ausgang:** *(bei
  Closure)*
- **Parität heißt nicht Gleichheit:** was `codepaths`/`ids` können, muss hier
  nicht identisch heißen; die Namen folgen dem Bestand. — **Ausgang:** *(bei
  Closure)*

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

*(wird mit dem Closure-Body gefüllt)*
