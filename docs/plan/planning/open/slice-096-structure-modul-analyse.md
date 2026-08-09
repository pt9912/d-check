# Slice slice-096: Modul `structure` — Analyse, Modul-Schnitt und Ablöse-Pfad

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** eigene Welle, bei Freigabe zu eröffnen — dieser Strang trägt eine
Closure-Bedingung **jenseits** seiner DoD: die Analyse allein löst kein
Adopter-Skript ab. Umsetzung, Paritäts-Beleg und der Ablöse-Pfad sind
Folge-Slices.

**Dieser Slice läuft ZUERST** (Auftraggeber-Entscheid 2026-08-09, nach dem
Backlog-Schnitt-Review): [slice-094](slice-094-closure-zaehl-paritaet.md),
[slice-097](slice-097-closure-glob-entkopplung.md) und
[slice-098](slice-098-closure-note-placeholder.md) schärfen alle dieselbe
Fähigkeit, die hier neu geschnitten wird. Sie zuerst einzeln auszuliefern hieße,
dreimal eine Semantik zu versprechen, die dieser Slice gerade neu definiert —
und sie danach zu migrieren. Die ursprüngliche Reihenfolge (094 zuerst) ist
damit **umgekehrt**.

**Bezug:** **Change Request** aus dem Schwester-Repo a-check (CR 1 seiner
Werkzeug-Abdeckungs-Analyse). Berührt
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
und [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) —
deren Closure-Fähigkeit ist der **Spezialfall** des beantragten Moduls.
Formvorbild für den Ablöse-Pfad:
[ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Entscheiden und festhalten, wie d-check **Struktur-Invarianten innerhalb eines
Dokuments** bekommt (bisher deckt der Modulsatz nur **Referenz**-Invarianten ab)
— und wie die in v0.52.0 ausgelieferte, planning-lokale Closure-Fähigkeit darin
aufgeht, **ohne** Config-Bruch bei Adoptern.

## 2. Die Kollision, die zuerst geklärt gehört

Der Antrag beschreibt `require-section` · `non-empty` · fence-treues
`min-sentences` · `forbid-pattern` — das ist, bis auf die Verallgemeinerung über
**beliebige Dokumentklassen**, genau die Fähigkeit, die
[slice-093](../done/slice-093-closure-note-gate.md) als
`planning.closure.*` ausgeliefert hat. Zwei Mechanismen für dieselbe Frage will
dieses Repo nicht.

Gemessen (2026-08-09, gegen eine Kopie des a-check-Bestands mit 76 Slices):
v0.52.0 deckt **eines** der drei beantragten Skripte, und das mit Kalibrierung —
das Adopter-Muster braucht wegen der fehlenden RE2-Lookahead-Unterstützung eine
positiv formulierte Ausschluss-Alternative, und die Floskel-Prüfung ist dort
zeilen-verankert, hier Teilstring. Nicht gedeckt sind die abschnitts-treue
Task-Zählung und die benannten Pflicht-Bausteine.

## 3. Abnahme-Punkte

1. **Modul-Schnitt.** Neues Modul `structure` (Liste von Regeln über
   Datei-Globs) — gegen die Alternative, `planning.closure` weiter auszubauen.
   Kriterium ist das aus
   [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md):
   querschnittlich über Dokumentklassen ⇒ eigenes Kürzel.
2. **Ablöse-Pfad für die Closure-Fähigkeit.** **Präzisierung nach dem
   Schnitt-Review (F-4):** superseded werden kann **nicht** die Anforderung
   [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
   — sie trägt **auch** die Aktiv-Status-Invariante (Roadmap ↔ in-progress), die
   `structure` nicht abdeckt und die bleibt. Wandern kann nur die **zweite
   Fähigkeit**. Der Ablöse-Pfad ist also ein Teil-Supersede, kein voller.
   Und die Analogie zu
   [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md) trägt
   nur halb: dort war der Alias **byte-identisch** (gleiche Grund-Codes), hier
   stünden `closure-note-*` gegen `section-*`. **Zu entscheiden:** bleiben die
   drei bestehenden Grund-Codes erhalten (dann emittiert `structure` sie für
   Closure-Regeln), oder werden sie durch `section-*` ersetzt (dann ist es ein
   Befund-Bruch, den ein Alias nicht abfängt)?
3. **Kardinalität mehrerer Closure-Abschnitte.** Das abzulösende
   Konsumenten-Skript prüft, wie **viele** Closure-Abschnitte ein Dokument trägt,
   und meldet Mehrdeutigkeit; die heutige Fähigkeit liest laut Spezifikation nur
   den **ersten** Treffer. Die Aktiv-Status-Prüfung desselben Moduls ist an
   dieser Stelle bereits fail-closed (mehrfache kanonische Überschrift ⇒
   Befund) — die Asymmetrie ist unbeabsichtigt. **Zu entscheiden:** deckt
   `structure` die Kardinalität ab, oder wird sie als Nicht-Ziel benannt?
4. **Semantik der zwei fehlenden Bausteine.** Abschnitts-treue Task-Zählung
   (Obergrenze innerhalb eines Abschnitts) und benannte Pflicht-Bausteine
   (Happy/Boundary/Negative als hervorgehobene Marken) — beides muss
   fence-treu und abschnitts-treu definiert werden, sonst wiederholt d-check den
   Fehler der abgelösten Skripte.
5. **Grandfathering.** Der Antrag verweist auf `exempt-paths`; zu prüfen, ob das
   für eine Stichtags-Regel („erst ab Slice N") ausreicht oder ob die
   Abschnitts-Regel selbst einen Skopus braucht.

## 4. Definition of Done

- [ ] Abnahme-Punkte 1–5 entschieden und begründet; Change Request in
      [`spec/lastenheft.md`](../../../../spec/lastenheft.md) formuliert
      (Bereichskürzel, Akzeptanzkriterien-Trio je Grund-Code) + begleitende ADR
      mit dem Supersede-/Alias-Pfad.
- [ ] **Abdeckungs-Messung** je beantragter Prüfung: gedeckt / gedeckt nach
      Kalibrierung / nicht gedeckt — mit Beleg, nicht per Lektüre. Die
      Paritäts-Fixtures liegen im Antragsteller-Repo vor und werden beigezogen.
- [ ] Folge-Slices geschnitten (Implementierung, Paritäts-Beleg, Ablösung des
      Alias) und als Dateien in `open/` angelegt — genannt ohne angelegt wäre
      dieselbe Klasse wie ein halluziniertes Gate. **Dazu gehört die
      Entscheidung über die drei bereits liegenden Slices**
      ([094](slice-094-closure-zaehl-paritaet.md),
      [097](slice-097-closure-glob-entkopplung.md),
      [098](slice-098-closure-note-placeholder.md)): gehen sie im `structure`-Schnitt
      auf, bleiben sie eigenständig, oder ändert sich ihr Zuschnitt? Sie stehen
      bis dahin ohne Wellen-Zuordnung.

## 5. Risiken / offene Punkte

- **Contract-Churn bei Adoptern:** `planning.closure` ist seit v0.52.0
  ausgeliefert. Ein Supersede ohne Alias wäre ein Bruch nach wenigen Tagen.
  — **Ausgang:** offen bis Abnahme-Punkt 2.
- **Der Modul-Schnitt könnte zu breit geraten.** „Struktur-Invarianten" ist eine
  Kategorie, keine Prüfung; ohne scharfe Grenze wächst `structure` zum
  Sammelbecken. — **Ausgang:** offen; die Grenze gehört in die ADR, nicht in die
  Implementierung.
- **Fremd-Repo-Abhängigkeit:** die Paritäts-Fixtures liegen nicht hier.
  — **Ausgang:** offen; beizuziehen, nicht nachzubauen.

## 6. Trigger

**Start** (`next` → `in-progress`): Freigabe. **Keine Vorbedingung** — dieser
Slice ist der erste des Strangs. Die frühere Fassung hängte ihn an
[slice-094](slice-094-closure-zaehl-paritaet.md) („erst die Zähl-Parität, dann
messen"); der Schnitt-Review hat gezeigt, dass das die Reihenfolge verkehrt: 094
sagt Deckungsgleichheit einer Semantik zu, die **dieser** Slice gerade neu
definiert, und weil ausgeliefert wird, was zuerst fertig ist, wäre der Vertrag
dann schon draußen (F-3).

**Rückführungen:** `in-progress` → `open`, falls die Messung ergibt, dass der
Antrag mehrere unabhängige Module beschreibt statt eines.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Analyse berührt `spec/` und `docs/plan/`; die Folge-Slices
  zusätzlich `internal/`. Alle unter dem Repo-Default GF
  (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten:** das Register führt **BEO-001**
  (Datei-Register driften unbemerkt gegen ihre Autoritäts-Tabelle). **Bewusst
  festgehalten, weil die Verwechslung naheliegt:** `structure` deckt BEO-001
  **nicht** ab. Dort geht es um eine Referenz **zwischen** Dokumenten (existiert
  eine Datei, die niemand registriert?), hier um die Form **innerhalb** eines
  Dokuments. Wer BEO-001 in diesem Slice erledigt glaubt, lässt die Lücke offen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — der Slice liefert Analyse und Zusage
(Change Request + ADR), der Code folgt in den Folge-Slices. Kein Brownfield: die
abzulösenden Skripte liegen in einem **anderen** Repo und werden nicht
rückdokumentiert, sondern durch eine eigene, spec-first formulierte Fähigkeit
ersetzt.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
