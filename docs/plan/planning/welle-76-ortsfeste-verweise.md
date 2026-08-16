# Welle welle-76-ortsfeste-verweise: Ein Verweis löst von jedem Lifecycle-Ort auf

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-76-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Erweiterung eines bestehenden
Moduls auf Konsumenten-CR).

**Verantwortlich:** pt9912. **Datum:** 2026-08-16.

---

## 1. Welle-Ziel

Die Invariante der Lifecycle-Move-Regel maschinell machen: ein relativer
Verweis in einer Slice-Datei muss aus **jedem** Lifecycle-Verzeichnis auflösen,
nicht nur vom Ist-Ort. [slice-095](in-progress/slice-095-links-resolve-from.md)
liefert `links.resolve-from` — dieselbe Prüfung mit erweiterter
Quellort-Menge, ohne den Schlüssel byte-identisch.

**Der Anlass ist dreifach im eigenen Bestand belegt** (zweimal bei der Closure
von slice-093, einmal die 19 gebrochenen Links bei der Eröffnung von welle-69) —
und er ist in **jeder** Welle seither wieder eingetreten: jede Closure dieser
Woche hat nach dem `git mv` Verweise von Hand nachgezogen, zuletzt 15 in
welle-75. Die Klasse skaliert mit der Zahl der Nachbarn, nicht mit der
Repo-Größe.

## 2. Trigger (Welle startet)

Freigabe des Auftraggebers (2026-08-16), WIP-Slot frei (welle-75 geschlossen,
`in-progress/` trägt nur die Roadmap). Der Slice ist von allen anderen Strängen
unabhängig.

## 3. Closure-Trigger (Welle schließt)

- [slice-095](in-progress/slice-095-links-resolve-from.md) liegt in `done/`.
- **Der Realdatenbeleg liegt vor:** die **Slice-Hälfte** des historischen
  19-Link-Bruchs wäre mit aktivem `resolve-from` **vor** dem Move rot gewesen
  (retro gemessen, mit dem Produkt); die Ziel-Wanderungs-Hälfte (Review-Reports,
  Wellendokument) ist als **Grenze** in der ADR benannt — sie ist eine andere
  Frage, keine offene Lücke.
- **Die Laufzeit-Zusage ist gemessen**, nicht behauptet
  ([`DC-QA-01`](../../../spec/lastenheft.md#dc-qa-01--performance)) — die
  Rückführung des Slice hängt daran.
- Release als **Minor** (opt-in; ohne den Schlüssel byte-identisch), die
  Richtung in der Notiz **offen** formuliert.
- `make fullbuild` grün.

## 4. Slices in dieser Welle

| Slice | Rolle |
|---|---|
| [slice-095](in-progress/slice-095-links-resolve-from.md) | `links.resolve-from`: hypothetische Quellorte, eigener Grund-Code, Ventil-Anschluss |

## 5. Abhängigkeiten

- [`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
  und die Auflösungs-Mechanik liegen vor; erweitert wird eine bestehende
  Anforderung ([ADR-0044](../adr/0044-geteiltes-referenz-ventil-quell-skopus.md)-Kriterium).
- Das bestehende Ventil `ignore-refs` ist der Kandidat für absichtlich
  ortsgebundene Verweise — bevor ein neues erfunden wird.

## 6. Out-of-Scope für diese Welle

- **Anker-Prüfung von hypothetischen Orten.** Die Anker-Frage hängt am
  Ziel-Dokument, nicht am Quellort (Abnahme-Punkt 3 des Slice).
- **Ein Auto-Rewrite der Verweise beim Move.** d-check bleibt diagnose-only in
  dieser Klasse; die Reparatur ist eine Autoren-Entscheidung.
- Die Chronologie-Ordnung (**BEO-005**) — steht als eigene geplante Welle in
  der Vorschau.

## 7. Closure-Notiz

_Ausstehend._
