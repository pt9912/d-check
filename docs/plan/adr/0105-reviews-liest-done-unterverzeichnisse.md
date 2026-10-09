# ADR-0105: Die Review-Deckung liest `done/` samt Unterverzeichnissen

**Status:** Accepted

**Datum:** 2026-10-09

**Autor:** pt9912

**Bezug:** [`DC-FA-RVW-001`](../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in);
[ADR-0081](0081-reviews-modul.md) (Entscheidung 4, die diese ADR ablöst),
[ADR-0082](0082-uebergangswaechter-reviews-observations.md) (der
Übergangs-Wächter, der `reviews` im Closure-Profil führt); `AGENTS.md` §3.6;
slice-264 <!-- d-check:status-provenance -->.

**Supersedes:** [ADR-0081](0081-reviews-modul.md), nur Entscheidung 4.

**Schärft:** — *(keine Spec-Stelle; eine Entscheidung über die Konfiguration
dieses Repos — das Produkt bietet `reviews.recursive` und
`reviews.skip-pattern` bereits an)*

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0081](0081-reviews-modul.md) Entscheidung 4 ließ das Modul `reviews`
nur die Dateien **unmittelbar** in `done/` lesen: Ein archivierter Stub liegt
in einem Unterverzeichnis, trägt keine DoD mehr und fiel damit ohne
Sonderfall aus der Kandidatenmenge.

Diese Annahme trägt nicht mehr. Wellenlose Slices schließen direkt nach
`done/wellenlos/`, und dort liegen Volltexte neben Stubs. Ohne Rekursion
prüft weder `make review-coverage` noch der Übergangs-Wächter
`make verify-closure-notes` die Review-Zusage eines solchen Slice — der
Übergang wird erkannt, die Prüfung fällt still aus. Gemessen beim Schnitt:
24 Volltexte unter `done/wellenlos/`, alle ungeprüft.

Das Produkt kann beides: `reviews.recursive` liest die Unterverzeichnisse,
`reviews.skip-pattern` nimmt eine Datei nach ihrem Inhalt aus
([`DC-FA-RVW-001`](../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)).
Die Konfiguration dieses Repos davon abweichen zu lassen, widerspricht einer
akzeptierten Entscheidung und braucht nach `AGENTS.md` §3.6 eine eigene.

## Entscheidung

1. Beide Profile — `.d-check.yml` (`make review-coverage`) und
   `.d-check.closure.yml` (`make verify-closure-notes`) — setzen
   `reviews.recursive: true`.
2. Der Stub fällt nicht mehr über seine Lage heraus, sondern über seinen
   Marker: `reviews.skip-pattern` trifft die Form, die das Archiv-Werkzeug
   schreibt, `> **ARCHIVIERT** — Volltext:` am Zeilenanfang.
3. Beide Profile tragen dieselben zwei Schlüssel. Weichen sie ab, antworten
   der Übergangs-Wächter und `make review-coverage` verschieden auf dieselbe
   Frage.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon (Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR
(MADR)).

| Option | Pro | Contra |
|---|---|---|
| Nichts tun: nicht rekursiv bleiben | Keine Änderung an einer akzeptierten Entscheidung; der Stub fällt ohne Muster heraus | Die Review-Zusage jedes wellenlosen Slice bleibt ungeprüft, obwohl der Übergang ihn erkennt |
| Rekursion nur im Closure-Profil | Der Übergang prüft den Slice; das Hauptprofil bleibt unberührt | `make review-coverage`, der Lauf zum Untersuchen eines roten Übergangs ([ADR-0082](0082-uebergangswaechter-reviews-observations.md)), meldet dann grün, wo der Übergang rot ist — gemessen im Review dieser Entscheidung |
| Die Volltexte archivieren | Die Lage-Regel stimmt wieder ohne neue Schlüssel | Ein Vorgang am Bestand, kein Teil der Prüfung; neue wellenlose Slices schließen trotzdem als Volltext |
| **Gewählt:** Rekursion in beiden Profilen, Stub über den Marker | Beide Läufe prüfen dieselbe Kandidatenmenge, jeder Volltext zählt | Ein Volltext, der den Marker zitiert, fällt still heraus; das Muster hängt an der Form des Archiv-Werkzeugs |

## Konsequenzen

- Ein Slice unter `done/wellenlos/` mit Review-Zusage und ohne Report macht
  beide Läufe rot. Heute ohne Wirkung: beide melden 0 Befunde.
- **Grenze:** Ein Volltext, der den Marker am Zeilenanfang zitiert, fällt
  still aus der Prüfung; ein Symlink auf ein Unterverzeichnis wird nicht
  verfolgt.
- **Grenze:** `reviews.exempt-paths` matcht mit blankem `path.Match` — `**`
  steht dort für genau ein Segment. Eine Ausnahme in der `**/`-Form trifft nur
  einen Slice eine Ebene unter `done/`; sie nennt den Pfad deshalb
  ausgeschrieben.

## Fitness Function (falls maschinell prüfbar)

Ohne den Report eines Slice unter `done/wellenlos/` melden
`make review-coverage` und `make verify-closure-notes` beide `review-missing`;
mit dem Report schweigen beide.

## Re-Evaluierungs-Trigger

- Das Archiv-Werkzeug ändert die Form des Stub-Markers.
- Ein Volltext muss den Marker am Zeilenanfang zitieren.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-09 | Proposed (Review R2 zu slice-264 fand den Widerspruch zu ADR-0081 Entscheidung 4) |
| 2026-10-09 | Proposed → Accepted (nach Review R1–R3 und Verifikation; die Fitness Function ist im Klon gebrochen) |
