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

- [x] Die drei Kopfzeilen sind entfernt; die Sicht trägt ihre weiter, mit
      Begründung an Ort und Stelle.
- [x] Kein Rest verweist auf die entfernten Zeilen (gemessen, nicht
      angenommen).
- [x] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Eine Abschnitts-Regel könnte die Zeile erwarten** — die Prüf-Config führt
  Regeln über Kopf-Abschnitte; sie werden vor dem Entfernen gelesen. —
  **Ausgang:** entfallen — gemessen: keine Regel, kein Test und keine
  Autoritäts-Doku hängt an den entfernten Zeilen. Der Review hat es
  gegengeprüft und zusätzlich belegt, dass **kein Gate** die Zeile in
  irgendeiner Richtung hält.
- **Die Sicht-Ausnahme sieht wie eine Inkonsistenz aus**, wenn ihre Begründung
  fehlt; sie steht deshalb in der Datei, nicht nur in dieser Notiz. —
  **Ausgang:** **eingetreten — und die erste Fassung der Begründung war selbst
  falsch.** Sie behauptete einen benannten Trigger, den der Bestand widerlegt
  (94 Commits am Vertrag zwischen zwei Marker-Hebungen, der Marker stand
  still). Jetzt sagt sie, was wahr ist.

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

**Geliefert:** von vier gemessenen Kopf-Zustandszeilen sind drei ersatzlos
entfernt — die beiden lebenden Register, weil *Aktiv* kein Zustand ist, den
ein Register je wechselt, und das Technik-Stratum, weil seine Historie das
Datum trägt. Die Sicht behält ihre und sagt jetzt an Ort und Stelle, warum:
sie führt keine Historie, also ist die Zeile ein Frische-Marker. Nebenbei
tragen beide Spec-Straten die Rollen-Zeile der Vorlage nach, die im Bestand
fehlte.

**Review** ([Report](../../../reviews/2026-08-22-slice-119-kopf-zustandsfelder-review.md)):
merge-blockierend — 0 HIGH, 1 MEDIUM, 3 LOW, alle eingearbeitet. Die
Negativ-Proben belegen die Zusage von der anderen Seite: **kein Gate** hält
die Zeile in irgendeiner Richtung — auch das Entfernen der Sicht-Zeile bliebe
grün. Genau deshalb tragen Briefing und Reviewer-Skill die Regel.

**Was ging anders als geplant — die Begründung war selbst ein Chronik-Feld:**
Ich hatte an die Sicht geschrieben, „ein benannter Trigger pflegt sie: jede
Änderung darüber zieht sie nach". Der Kanon sagt die **Gegenrichtung**, und
der Review hat meine Behauptung am Bestand widerlegt: zwischen zwei
Marker-Hebungen fassten **94 Commits** Lastenheft oder Spezifikation an, der
Marker stand still. Ich hatte einen Trigger erfunden, den niemand fährt —
dieselbe Klasse, gegen die der Anker dieser Welle gerichtet ist, in meiner
eigenen Ausnahme-Begründung. Wer eine Ausnahme schreibt, muss ihren Träger
benennen können; kann er es nicht, ist es keine Ausnahme, sondern eine
Behauptung. Jetzt steht dort: die Zeile wird gesetzt, wenn jemand die Sicht
gegen den Code hält — steht sie still, ist das die ehrliche Aussage „seither
nicht geprüft".

- **Steering-Loop-Eintrag:** kein neuer Träger — die Regel liegt seit
  slice-118 im Briefing und im Reviewer-Skill; dieser Slice wendet sie an.
- **Beobachtungs-Register (`../observations.md`):** keine neue Beobachtung.
  Der erfundene Trigger ist die Klasse, die der neue HIGH-Anker beschreibt —
  er hat beim ersten Lauf gegen den eigenen Text getroffen, was für ihn
  spricht und gegen mich.
- **Folge-Slices:** [slice-120](../open/slice-120-register-und-drift-log.md) —
  letzter Slice der Welle.
- **Risiken aus §6:** beide mit Ausgang (§5) — eines entfallen, eines
  eingetreten.
- **Drei Paarungen:** Wellen-Slice — die Paarungen prüft die Welle-Closure.
