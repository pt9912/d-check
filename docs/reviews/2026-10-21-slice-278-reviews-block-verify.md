# Verifikation — slice-278: der `reviews`-Block gegen die Vorlage von v6.18.1

- **Rolle:** Verifier (Modul 11) — DoD, Plan und Sensor-Seite gegen den Stand
- **Gegenstand:** `docs/plan/planning/in-progress/slice-278-reviews-block-auf-vorlage-pruefen.md` §1, §2, §3 (Messung, Plan-Änderung, Spiegel-Treffer), §4, §6; `harness/sensors/review-coverage.md`
- **Range:** `607648bc..HEAD` (Spitze `1f3b537a`); der Review-Report `2026-10-21-slice-278-reviews-block.md` nur als Kontext
- **Datum:** 2026-10-21

## Messungen (selbst gefahren, Wegwerf-Worktree unter dem Home, danach entfernt)

Image `d-check:latest`, Modul `reviews` allein, `--network none`, Mount `:ro`; Profil `.d-check.yml` mit den 22 `--disable`, Übergangs-Profil mit `--config .d-check.closure.yml --enable planning --enable structure --enable spans --enable reviews` (wie das Makefile-Target).

| Lauf | Kopie / Konfiguration | Plan behauptet | Gemessen | |
|---|---|---|---|---|
| Zählung | Slice-Dateien unter `done/` / davon Stub (Marker) / lebend | 276 / 229 / 47 | 276 / 229 / 47 | ok |
| A | HEAD-Konfiguration, alle Reports | Exit 0, 0 Befunde | beide Profile Exit 0, 0 Befunde | ok |
| B | HEAD-Konfiguration, `docs/reviews/*` entfernt | Exit 1, 42 `review-missing` | beide Profile Exit 1, **42** | ok |
| B′ | wie B, Konfiguration von `607648bc` | (gleich 42) | Exit 1, **42** | ok |
| C | Wortlaut „Review durchgeführt" → „Durchsicht erfolgt" in allen Plänen (42 Dateien), keine Reports, Konfiguration `607648bc` | Exit 0, 0 Befunde | beide Profile **Exit 0, 0 Befunde** | ok |
| D | wie C, Konfiguration HEAD | Exit 1, „keine Review-Zusage unter 47 Kandidat(en)" | beide Profile Exit 1, ein Befund `keine Review-Zusage unter 47 Kandidat(en) — require-promises …` | ok |
| E | gesunder Bestand, HEAD-Konfiguration | Exit 0 | beide Profile Exit 0, 0 Befunde | ok |
| F | gesunder Bestand, `match: name` | Exit 1, 19 Befunde | Exit 1, **19** `review-missing` | ok |
| G | alle lebenden Slices als Stub markiert | Exit 1 „leere Pruefmenge … fail-closed"; mit `skip-allows-empty: true` Exit 0 | Exit 1 `leere Pruefmenge: 0 Kandidat(en), 0 Review-Zusage(n) … fail-closed`; mit dem Schlüssel Exit 0 | ok |
| Zusagen | lebende Slices mit Zusage-Muster (grep-Näherung des Modul-Musters) | 43 von 47 (Sensor-Seite) | 43 von 47 | ok |

Hinweis: Der erste D-Lauf im Profil `main` brach einmal mit „nicht lesbar: open …slice-228…" ab (Exit 2), während die Datei existiert; zwei Wiederholungen ohne Eingriff liefen sauber auf Exit 1. Das war ein Artefakt meines `sed -i` über den Mount, kein Befund am Repo. In der Reports-Zählung lagen 120 Einträge in `docs/reviews/` (119 Reports laut Plan, plus der Report dieses Slice bzw. Unterverzeichnis `archiv/`); die Zusage-Zahlen (42) sind davon unberührt.

## Gates am Repo (ohne Pipe, Exit-Code erfasst)

| Kommando | Ergebnis |
|---|---|
| `make review-coverage` | `RC=0`, `1192 Datei(en) geprüft, 0 Befund(e)` |
| `make verify-closure-notes` | `RC=0`, `1054 Datei(en) geprüft, 0 Befund(e)` |
| `make gates` | `GATES_RC=0`; `baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `git status` danach leer |

## DoD-Punkte

- **Messung — erfüllt.** Alle Zahlen der Plan-Tabelle (Läufe A–G) und die korrigierten Zählungen (276 / 229 / 47, 42, 19) reproduziert; an Kopien, nicht am Arbeitsbaum.
- **Konfiguration — erfüllt.** `require-promises: true` steht in `.d-check.yml` und `.d-check.closure.yml` (Diff); beide Gate-Läufe grün **und** nicht über leerer Menge (Lauf B: 42 Zusagen sichtbar; Lauf D belegt, dass der Schalter den Leerlauf meldet). `match` und `skip-allows-empty` fehlen mit Begründung im Plan (F: 19 Fehlbefunde; G: nähme die Fail-closed-Antwort, 47 lebende Slices heute) — beide Begründungen durch Messung gedeckt.
- **`make gates` grün — erfüllt** (siehe oben).
- **Review — Report liegt vor**; Verifikation mit diesem Bericht.
- **Offen, keine Verletzung:** Beobachtungs-Register, Closure-Notiz §7, Risiko-Ausgänge §6 (beide „offen bis Closure"; Gegenprobe: Risiko 1 eingetreten im Sinne „Leerlauf besteht latent", nicht als Ist-Zustand — heute 42 sichtbare Zusagen; Risiko 2 entfällt messbar, Lauf E/A: kein Befund über Altbestand), drei Paarungen, `git mv` nach `done/wellenlos/`.

## Bewusstes Brechen (Modul 11) — Wegwerf-Worktree

| Zurückgenommen | Lauf | Ergebnis |
|---|---|---|
| `require-promises: true` in **beiden** Profilen (Konfiguration `607648bc`) bei geändertem Wortlaut | C | beide Profile Exit 0, 0 Befunde — grün über leerer Menge, aus dem behaupteten Grund |
| `require-promises` nur im **Closure**-Profil zurückgenommen (Hauptprofil behält ihn) | Kopplung | `main` Exit 1 (Leerlauf-Befund), `closure` Exit 0 — die Profile antworten verschieden |
| `require-promises` nur im **Hauptprofil** zurückgenommen (Closure behält ihn) | Kopplung | `main` Exit 0, `closure` Exit 1 — umgekehrt verschieden |

Die Kopplung ist damit real, und der Kopplungs-Kommentar beider Dateien („drei Schlüssel", `require-promises` genannt) beschreibt sie zutreffend.

## Sensor-Seite `harness/sensors/review-coverage.md` gegen die Messung

| Aussage | Befund |
|---|---|
| Punkt 1: beide Profile tragen „dieselben drei Schlüssel (`recursive`, `skip-pattern`, `require-promises`)" | stimmt (Diff beider Blöcke) |
| Punkt 2: 42 erkannte Zusagen, 42 `review-missing` ohne Reports | stimmt (B) |
| Punkt 2: Muster trifft nichts mehr → ohne Schalter grün, mit Schalter „keine Review-Zusage unter N Kandidat(en)" | stimmt (C, D) |
| Punkt 2: Grenze „Totalausfall, nicht der einzelne Slice" | konsistent; folgt aus dem Verhalten (D meldet nur bei Null Zusagen). Einen Einzel-Slice-Ausfall habe ich nicht eigens gefahren (LOW 2) |
| Punkt 3: 47 lebende Slices, 43 mit Zusage | stimmt (47 / 43) |
| Punkt 3: `skip-allows-empty` nicht gesetzt; „leere Prüfmenge … fail-closed" | stimmt (G) |
| Punkt 3: Archivieren aller Zusage-Träger lässt „keine Review-Zusage unter N Kandidat(en)" zurück, `skip-allows-empty` greift dort nicht | aus der Zusammensetzung (43 von 47) und dem Wortlaut von D abgeleitet, nicht als eigener Lauf gefahren (LOW 2); G zeigt den benachbarten Fall (alles Stub) |

## Plan-vs-Code (Abgrenzung §1)

- `git diff 607648bc..HEAD -- internal/ spec/` ist leer: keine Moduländerung, keine Spec-Änderung.
- Keine nachgeholten Reports: der Diff legt unter `docs/reviews/` nur den Report dieses Slice an; Messläufe B/C/D nur an Kopien.
- Keine weiteren v6.18.1-Deltas übernommen (`match`, `skip-allows-empty`, Vorlagen-Pfade nicht angefasst).
- Dateien im Diff: `.d-check.yml` (§3), `.d-check.closure.yml` (Plan-Änderung vor dem Code), `harness/sensors/review-coverage.md` (Spiegel-Treffer §3), der Slice-Plan selbst, der Review-Report. Nicht ausdrücklich im Plan genannt: `docs/plan/planning/in-progress/roadmap.md` (Ruhe-Marker „Nichts in Arbeit." entfällt — Folge des Anspruchs, von `planning-check` gekoppelt), `harness/conventions/MR-080-baseline-v6181.md` und `done/wellenlos/slice-277-…` (reine Pfad-Nachzüge `open/` → `in-progress/` der Verweise auf diesen Slice) — LOW 1.

## Befunde

- **HIGH / MEDIUM:** keine.
- **LOW 1:** Drei Dateien im Diff (`roadmap.md`, `MR-080-…`, `done/wellenlos/slice-277-…`) sind weder in §3 noch in einer Plan-Änderungsnotiz genannt. Inhaltlich sind es Lifecycle-Folgen (Ruhe-Marker, Pfad-Nachzug des eigenen Slice), aber `slice-277` liegt in `done/` (eingefroren), und die Zeile dort trägt noch den Text „— `open/`" neben dem neuen Pfad `in-progress/` (`docs/plan/planning/done/wellenlos/slice-277-baseline-v6181-bump.md`, Zeile um 207). Zustandsangabe `open/` ist nun falsch; das ist Historie der Closure-Notiz, aber ein toter Hinweis neben dem korrekten Link.
- **LOW 2:** Zwei Aussagen der Sensor-Seite (Grenze „nicht der einzelne Slice"; „Archivieren aller Zusage-Träger → Befund trotz Vollständigkeit") sind aus Schaltverhalten und Zählung abgeleitet, nicht als eigener Messlauf in der Plan-Tabelle belegt. Sie sind plausibel und mit D/G konsistent; ein Lauf H (ein Slice mit geändertem Wortlaut neben trägen anderen → Exit 0; Zusage-Träger als Stub → Befund) würde sie schließen (`harness/sensors/review-coverage.md` Punkt 2/3).

## Ergebnis

DoD-Liefer-Punkte **Messung** und **Konfiguration** sowie `make gates` bestätigt und reproduziert; Bewusstes Brechen zeigt Rot aus dem behaupteten Grund und die Kopplung der Profile. Abgrenzung eingehalten. Offen (planmäßig): Register, Closure-Notiz, Risiko-Ausgänge, Paarungen, `git mv`. Das Repo ist nach den Läufen sauber (`git status` leer, Worktree entfernt und geprunt).
