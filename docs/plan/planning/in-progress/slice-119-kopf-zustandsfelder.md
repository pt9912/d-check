# Slice slice-119: Kopf-Zustandsfelder — drei weg, eines bleibt begründet

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-81-zustandsfelder](../welle-81-zustandsfelder.md) (zugeordnet
bei der Eröffnung).

**Bezug:** Baseline-Regelwerk `grundlagen-harness-dateien.md` §Was ein
Kommentar trägt („Die Kopfzeile eines lebenden Registers ist derselbe Fall")
und `modul-03-spec.md` §Ziel-Form: Spezifikation („Kein Kopf-Datum, kein
Kopf-Status") sowie §Ziel-Form: Architektur-Sicht (der Frische-Marker bleibt);
die Ziel-Formen der vendorten Vorlagen für Roadmap, Beobachtungs-Register,
Spezifikation und Sicht.

**Berührte Spec-Stellen:** `spezifikation.md` Kopf — der Verweis zeigt
aufwärts.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Gemessen tragen **vier** Dateien die Zeile `**Status:** Aktiv. **Letzte
Änderung:** …`. Drei verlieren sie: die Roadmap und das Beobachtungs-Register,
weil *Aktiv* kein Zustand ist, den ein Register je wechselt, und ein Datum,
das niemand pflegt, einen behauptet — ihr Zustand ist ihr Inhalt, ihr
Änderungsdatum hält `git`. Die Spezifikation ebenso, weil ihre Historie das
Datum trägt und zwei Felder für eines driften. Die **Sicht behält** ihre
Zeile: dort ist sie der bewusste Frische-Marker, weil die Sicht keine Historie
hat — und genau das steht künftig auch dort.

## 2. Vorgehen

1. **Messen statt annehmen:** alle Vorkommen der Kopfzeile im Baum auflisten
   (der Bestand ist vier); je Datei die Ziel-Form der vendorten Vorlage
   danebenlegen.
2. **Entfernen** in Roadmap, Beobachtungs-Register und Spezifikation —
   ersatzlos, samt der Leerzeile.
3. **Begründen, wo sie bleibt:** die Sicht bekommt einen halben Satz, warum
   ihr Kopf-Datum ein Frische-Marker und kein Protokoll ist — sonst liest der
   nächste Durchgang die Zeile als Rest und entfernt sie.
4. **Rand prüfen:** ob eine Prüf-Regel, ein Test oder eine andere Datei die
   entfernte Zeile erwartet (Abschnitts-Regeln, Muster in der Prüf-Config,
   Verweise in der Autoritäts-Doku) — per `grep` abgeleitet.
5. Unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Der Vertrag bleibt unberührt:** das Lastenheft trägt Version und Status
  weiter — die Vorlage führt beide unverändert.
- **Keine `Stand`-Zellen, kein Drift-Log** (slice-120).
- **Kein Produkt-Code.**

## 4. Definition of Done

- [ ] Die drei Kopfzeilen sind entfernt; die Sicht trägt ihre weiter, mit
      Begründung an Ort und Stelle.
- [ ] Kein Rest verweist auf die entfernten Zeilen (gemessen, nicht
      angenommen).
- [ ] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Eine Abschnitts-Regel könnte die Zeile erwarten** — die Prüf-Config führt
  Regeln über Kopf-Abschnitte; sie werden vor dem Entfernen gelesen. —
  **Ausgang:** *(bei Closure)*
- **Die Sicht-Ausnahme sieht wie eine Inkonsistenz aus**, wenn ihre Begründung
  fehlt; sie steht deshalb in der Datei, nicht nur in dieser Notiz. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): slice-118 in `done/` (die Regel steht im
Briefing).

**Rückführungen:** keine erwartet.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Spec-Straten (GF), Planungs-Register (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): BEO-002
  (Spiegel-Pflicht — die Kopfzeile könnte anderswo erwartet werden, §2
  Schritt 4); BEO-006/009/010 als Arbeitsregeln.

Slice-ID: slice-119. Betroffene IDs: — (Form-Regel der Baseline). Module:
Spezifikation, Sicht, Roadmap, Beobachtungs-Register. Gates:
`make doc-check` (eng), `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Form-Angleichung an die adoptierte
Baseline.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
