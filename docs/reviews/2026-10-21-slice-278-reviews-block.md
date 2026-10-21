# Review slice-278 — `reviews`-Block gegen die Vorlage von v6.18.1

- **Review-Art:** Code/Konfiguration (gegen Plan, ADR, Hard Rules)
- **Gegenstand:** `slice-278`, Diff `607648bc..HEAD` (Messung und Plan-Änderung, Konfiguration, Sensor-Seite)
- **Skill:** `reviewer.md` @ 1.20.0
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-10-21
- **Zusätzlich gefahren:** Wegwerf-Kopie des Repos unter dem Home (Docker mountet `/tmp` nicht), Modul `reviews` allein: alle Reports; keine Reports; geänderter Wortlaut mit der committeten und mit der alten Konfiguration; beide Profile; `match: name`. Die Kopie ist entfernt (`git worktree remove`, `prune`).

## Findings

### F-1 — MEDIUM — `match: name` ergibt 19 Fehlbefunde, nicht 42

- **befund:** Die Messungs-Tabelle des Plans nannte für Lauf F 42 Befunde; das ist die Zahl aus Lauf B (keine Reports). Mit allen Reports und `match: name` meldet das Modul **19** `review-missing`; die übrigen 23 Reports tragen den vollen Slug und treffen. Die Entscheidung, `match: name` nicht zu setzen, bleibt richtig (der Lauf ist rot), nur Zahl und Verallgemeinerung („Reports heißen nach `slice-<NNN>`") waren falsch.
- **verifizierbar:** ja — derselbe Lauf; vom Implementer nachgezählt (19).
- **klasse:** messung-zahl-aus-anderem-lauf

### F-2 — MEDIUM — „22 lebende Slices" und „254 Stubs" sind falsch (richtig: 47 und 229)

- **pfad:** `harness/sensors/review-coverage.md` Punkt 3; Plan (Zeilen B und „heute liegen")
- **befund:** Unter `done/` (rekursiv) liegen 276 `slice-*.md`, davon 229 mit dem Stub-Marker und 47 ohne. Lauf D bestätigt „unter 47 Kandidat(en)". Die 22 waren 276 − 254, und die 254 stammten aus einem `grep -rl` über alle Dateien statt nur Slices. Der Satz steht im Sensor-Vertrag und belegt dort die Grenze des Ruhezustands.
- **verifizierbar:** ja — `grep -rL` auf den Marker unter `done/`.
- **klasse:** messung-zahl-aus-anderem-lauf

### F-3 — LOW — Kopplungs-Kommentar im Closure-Profil nennt `require-promises` nicht

- **pfad:** `.d-check.closure.yml`, Kommentar über dem `reviews`-Block
- **befund:** Der Kommentar nannte nur `recursive`, `skip-pattern` und die Bestands-Ausnahme; der neue Schlüssel wirkte wie ein ungepflegter Spiegel.
- **klasse:** kopplungskommentar-nicht-mitgezogen

### F-4 — INFO — Grenze der Sensor-Doku: Archivierung lässt nur Slices ohne Zusage liegen

- **befund:** `require-promises` zählt Zusagen nach dem `skip-pattern`. Heute tragen 43 von 47 lebenden Slices eine Zusage, 4 nicht. Archiviert eine Wellen-Closure alle, die eine tragen, melden beide Profile „keine Review-Zusage unter N Kandidat(en)", und `skip-allows-empty` greift dort nicht. Ein Szenario, kein Fehler; es steht jetzt in der Grenze.

### F-5 — INFO — Reihenfolge der Spiegel-Liste

- **befund:** Die Spiegel-Liste für `review-coverage.md` stand vor ihrem Edit (erfüllt); der Plan trägt unter dem Spiegel-Absatz noch die Anweisung „Treffer vor dem ersten Edit im Plan listen" als unausgefüllten Text. Keine Wirkung.

## Geprüft, ohne Befund

- **Messung (a)–(d):** mit allen Reports Exit 0 in beiden Profilen; ohne Reports 42 `review-missing` in beiden (= 42 erkannte Zusagen); bei geändertem Wortlaut ohne Reports: **committete Konfiguration Exit 1** mit genau einem Befund „keine Review-Zusage unter 47 Kandidat(en) — require-promises …", **Konfiguration von `607648bc` Exit 0** (stilles Grün) — die Änderung wirkt, und die gleiche Antwort in beiden Profilen belegt die Kopplung.
- **`skip-allows-empty`:** die Begründung des Verzichts stimmt mit `reviews.go` und ADR-0106 überein (ohne den Schalter bleibt die Leere nach `skip-pattern` fail-closed).
- **Sensor-Seite:** „Totalausfall, nicht der einzelne Slice" und „null Zusagen ist hier ein Befund" stimmen mit `reviewLeerlauf` überein.
- **`.d-check.yml`-Kommentar:** trägt Zusage, Grenze und Kopplung; Herkunft als auflösbare Felder (`ADR-0081` und die Kopplungs-ADRs).
- **AGENTS.md §3.6:** `require-promises: true` ist eine Verschärfung, keine Lockerung.
- **Abgrenzung §1:** keine Moduländerung, keine Reports nachgeholt.
- **Schritt 15:** die Plan-Notiz zur Closure-Datei steht vor dem Config-Commit.
