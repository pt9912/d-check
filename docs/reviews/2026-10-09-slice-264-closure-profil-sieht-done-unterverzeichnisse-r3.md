# Review R3 — slice-264: Das Closure-Profil prüft die Slices in den Unterverzeichnissen von `done/`

**Review-Art:** Code (Diff gegen Plan, ADRs und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-264, Plan-Änderung `4c38b637` und Korrektur-Commit `b4cb0202`
gegen die Findings aus Review R2 (`84c5577e`)
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-264 samt Plan-Änderungen; Reviews R1 und R2 zu
slice-264; [ADR-0081](../plan/adr/0081-reviews-modul.md),
[ADR-0105](../plan/adr/0105-reviews-liest-done-unterverzeichnisse.md),
[ADR-0082](../plan/adr/0082-uebergangswaechter-reviews-observations.md), der
[ADR-Index](../plan/adr/README.md) §Konventionen;
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in);
`AGENTS.md` §3.5, §3.6, §3.7, §5 Regeln 4 und 5; Baseline `v6.17.0` ·
`templates/docs/plan/adr/NNNN-titel.template.md` §Verglichene Alternativen.

## Messungen (Kommando und Ergebnis)

Im Wegwerf-Klon des Stands `b4cb0202`.

- `make adr-check RANGE=84c5577e..b4cb0202` → `0 Befund(e)`; ebenso
  `RANGE=fe54da7e..b4cb0202` über den ganzen Slice. Gegenprobe: Entscheidung 4 in
  ADR-0081 umformuliert und committet → `core-drift-vcs` auf
  `docs/plan/adr/0081-reviews-modul.md`, Exit rot. Das Gate liest den neuen Eintrag
  als reinen `## Geschichte`-Anhang und schlägt bei einer Änderung am Kern an.
- Spiegel-Suche nach der Rekursions-Entscheidung:
  `grep -rn -i "nicht rekursiv\|nicht-rekursiv\|unmittelbar"` über die lebenden
  Dateien (ohne `docs/reviews/`, `done/`, vendorte Baseline, CHANGELOG), gefiltert
  auf `review`/`done`. Treffer: ADR-0081 Entscheidung 4 und Fitness-Function-Zeile
  (immutabel, jetzt mit Geschichte-Zeile), Lastenheft und Spezifikation (beschreiben
  den **Default** ohne `reviews.recursive` und das nie rekursive `reviews-dir`),
  `--print-config`-Vorlage (`reviews-dir … (nicht rekursiv)`, richtig). Dazu das
  Benutzerhandbuch (R1 F-8, zurückgestellt). Kein weiterer Spiegel sagt für die
  Konfiguration dieses Repos „nicht rekursiv".

## Stand der R2-Findings

| R2 | Stand | Beleg |
|---|---|---|
| F-9 MEDIUM | behoben, Restpunkt F-14 | ADR-0105 (Proposed) mit `**Supersedes:**`-Feld auf Entscheidung 4, Index-Zeile am Ende der Tabelle, Geschichte-Zeile an ADR-0081; beide Profile und `review-coverage.md` verweisen auf ADR-0105 |
| F-10 LOW | behoben | Grenze 4 entfernt, Nummern nachgezogen |
| F-11 LOW | behoben | GRENZE und ADR-0105 nennen die gemessene Semantik (`**` = genau ein Segment) |
| F-12 LOW | angenommen mit Begründung | Bestandsform (16 Einträge); ein Retrofit ist ein eigener Vorgang. Der Widerspruch zu `observations/README.md` bleibt bestehen, ist aber älter als dieser Slice |
| F-13 INFO | behoben | eigene Historie-Zeile „Nachzug nach Review an §7 `SPEC-095`" |
| R1 F-7 INFO | behoben | Plan §1 und §8 nennen 24, mit der Zählform angegeben |

## Findings

### F-14 — LOW — ADR-0081 trägt „abgelöst" als Tatsache, ADR-0105 ist erst `Proposed`

- **kategorie:** LOW
- **quelle:** [ADR-Index](../plan/adr/README.md) §Konventionen (Status der abgelösten
  ADR und Index-Status werden bei `Accepted` nachgezogen); `AGENTS.md` §3.5 (eine
  Geschichte-Zeile an einer `Accepted`-ADR ist nur anhängbar, nicht korrigierbar)
- **pfad:** `docs/plan/adr/0081-reviews-modul.md` · „**Entscheidung 4 abgelöst** durch [ADR-0105]"
- **befund:** Der Plan führt ADR-0105 ausdrücklich „`Proposed` bis zur Closure"; die
  angehängte Zeile in der immutablen ADR-0081 sagt schon jetzt „abgelöst". Wird
  ADR-0105 im Review oder bei der Abnahme geändert oder verworfen, bleibt der Satz
  stehen und lässt sich nur mit einer weiteren Zeile widerrufen. Die übrigen
  Supersede-Markierungen (Status `Accepted (teil-superseded: ADR-0105)` in der Datei,
  Index-Status von ADR-0081) stehen dagegen richtigerweise noch aus.
- **verifizierbar:** nein — Reihenfolge-Urteil; `grep -n "Status:" docs/plan/adr/0105-*.md`
  gegen die Geschichte-Zeile in ADR-0081.
- **klasse:** `immutabler-anhang-vor-annahme`

### F-15 — LOW — Die Alternativen-Tabelle von ADR-0105 folgt nicht der Vorlagen-Form

- **kategorie:** LOW
- **quelle:** Baseline `v6.17.0` · `templates/docs/plan/adr/NNNN-titel.template.md`
  §Verglichene Alternativen · „mindestens drei Optionen mit Pro/Contra"; `AGENTS.md` §5 Regel 4/5 (neue ADR nach Vorlage)
- **pfad:** `docs/plan/adr/0105-reviews-liest-done-unterverzeichnisse.md` · „| Alternative | Warum nicht |"
- **befund:** Die Tabelle führt drei verworfene Optionen nur mit Contra; die gewählte
  Option steht nicht als Zeile, ein Pro fehlt durchgehend. ADR-0103 und ADR-0104 tragen
  dieselbe Form, ADR-0100 bis ADR-0102 die der Vorlage — der Bestand ist gespalten,
  die Vorlage nicht.
- **verifizierbar:** ja — `sed -n '/## Verglichene Alternativen/,/## Konsequenzen/p'` auf die Datei.
- **klasse:** `adr-alternativen-ohne-pro-contra`

### F-16 — INFO — Bei `Accepted` fällig: Status und Index-Zeile von ADR-0081

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `docs/plan/adr/README.md` · „Review-Report-Deckung wird das Modul `reviews` | Accepted"
- **befund:** Nach Bestand (ADR-0057/ADR-0070, ADR-0002/ADR-0014) trägt eine teilweise
  abgelöste ADR `Accepted (teil-superseded: ADR-NNNN)` in Datei und Index. Für ADR-0081
  ist das mit der Annahme von ADR-0105 nachzuziehen; heute korrekt offen.
- **verifizierbar:** ja — `grep -n "0081" docs/plan/adr/README.md`.
- **klasse:** `teil-supersede-markierung-ausstehend`

## Negativbefunde

- **Teil-Ablösung, Form:** `**Supersedes:**` mit Einschränkung „nur Entscheidung 4",
  wie ADR-0014 (`ADR-0002 §4 — nur …`); Index-Zeile chronologisch am Ende, Status
  `Proposed`, Titel mit `(supersedes …)` wie bei ADR-0092 bis ADR-0095; Bezug,
  `Schärft: —` mit Begründung (Konfigurations-Entscheidung ohne Spec-Stratum),
  Re-Evaluierungs-Trigger, Fitness Function, Geschichte vorhanden — ohne Befund (Tabelle: F-15).
- **Provenance-Marker in ADR-0105:** `slice-264 <!-- d-check:status-provenance -->` im
  Bezug-Feld zeigt den Entstehungsort und begründet keine Entscheidung — ohne Befund.
- **`make adr-check`:** hält die Änderung als Geschichte-Anhang, schlägt bei einer
  Kern-Änderung an (Messung) — ohne Befund.
- **Inhalt ADR-0105 gegen Konfiguration:** Entscheidung 1–3 entsprechen beiden
  `reviews`-Blöcken (gleiche Schlüssel, gleiches Muster); die Glob-Grenze stimmt mit der
  R2-Probe überein — ohne Befund.
- **Kommentare (`AGENTS.md` §3.7):** die neuen und geänderten Kommentare in beiden Profilen
  tragen Kopplung und Grenze mit ADR-Feld, keine Chronik — ohne Befund.
- **Spec-Historie:** die neue Zeile steht oben, chronologisch richtig, und nennt
  Kennung und Anker — ohne Befund.
- **Plan-Änderung vor dem Code:** `4c38b637` liegt vor `b4cb0202` und nennt ADR-0105,
  Geschichte und Index — ohne Befund.
- **Weitere „nicht rekursiv"-Spiegel:** keine außer dem Handbuch (siehe Messung).

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 2 | 1 |

## Verdikt

**Kein blockierender Befund.** R2 ist abgearbeitet; F-14 und F-15 sind klein und
lassen sich vor der Annahme von ADR-0105 nachziehen oder mit Begründung annehmen.
F-16 gehört zur Annahme selbst.
