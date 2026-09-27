# ADR-0096: `AGENTS.md`-Abschnitts-Körper wandern nach `harness/rules/`, Überschriften und Anker bleiben stehen

**Status:** Accepted

**Datum:** 2026-09-27

**Autor:** pt9912

**Bezug:** Auftraggeber-Entscheid 2026-09-27 (Schwelle 400 Zeilen,
Werkzeug `file[].max-lines`, Ziel-Dateien `AGENTS.md` +
`harness/README.md`); Anlass ist die offene Beobachtung
[`BEO-ALL/briefing-datei-ueberschreitet-lade-budget`](../planning/observations/BEO-ALL/briefing-datei-ueberschreitet-lade-budget/observation.md)
(3× erreicht mit slice-237); slice-239 <!-- d-check:status-provenance -->.

**Schärft:** — reine Konfigurations- und Struktur-Entscheidung, keine neue
oder geänderte `DC-*`-Anforderung ([`DC-FA-FILE-001`](../../../spec/lastenheft.md#dc-fa-file-001--zeilen--und-byte-obergrenzen-einer-ganzen-datei-modul-file-opt-in)
existiert bereits seit ADR-0088 und wird hier nur genutzt, nicht geändert).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`AGENTS.md` wird in jeden Agentenlauf geladen und wächst mit jeder
verkörperten Regel, ohne von selbst zu schrumpfen — vor diesem ADR 575
Zeilen. Der Auftraggeber hat die Größe als Mangel benannt;
[`DC-FA-FILE-001`](../../../spec/lastenheft.md#dc-fa-file-001--zeilen--und-byte-obergrenzen-einer-ganzen-datei-modul-file-opt-in)
(Modul `file`) existiert bereits, war aber nie auf `AGENTS.md` aktiviert
(reiner Sensor-**Besitz** ohne **Tragen**). Eine Aktivierung bei der
heutigen Zeilenzahl liefe sofort rot — die Schwelle 400 setzt eine Kürzung
voraus, keine bloße Konfigurationszeile.

Der naheliegende Umweg, nur die Chronik/Historie zu streichen, bringt laut
Beobachtung nur einen kleinen Teil zurück; der Rest ist Regeltext. Dieses
Repo kennt bereits ein Muster für genau dieses Problem:
`harness/conventions.md`s Adaptions-Block hält im Hauptdokument nur eine
Index-Zeile je `MR-<NNN>`, der Volltext liegt in einer eigenen Datei unter
`harness/conventions/`. `harness/sensors/<gate>.md` löst dieselbe Aufgabe
für Gate-Verträge. `AGENTS.md`s Hard-Rule- und Dokumentations-Regel-Bodies
haben bisher kein Äquivalent.

## Entscheidung

### Ein neues, drittes Verzeichnis `harness/rules/` trägt die Regel-Körper

Für einen ausgewählten Abschnitt (`AGENTS.md` §3.1, §3.7, §3.9 sowie die
gesamte Liste in §5) wird der **Körper** — Begründung, Durchsetzungs-Detail,
Grenzen, Bestandsgrenze, in §5 zusätzlich der komplette Regeltext — in eine
eigene Datei unter `harness/rules/<slug>.md` verschoben (für §5, das
mehrere unabhängige Regeln in einer Liste führt, eine Datei je Regel unter
`harness/rules/dokumentations-regeln/<slug>.md`). In `AGENTS.md` bleiben:

- die **Überschrift** unverändert (Wortlaut und damit der GitHub-Slug-Anker
  bleiben stabil — kein Nachzug für die ~24 Dateien nötig, die
  `AGENTS.md#3.x` bzw. `AGENTS.md §5` heute zitieren);
- der **operative Kern** in ein bis wenigen Sätzen (das „Falsch"/„Richtig",
  die Kernregel selbst — was ein Implementer beim schnellen Lesen ohne
  weiteren Klick braucht);
- ein **Pointer-Satz** auf die neue Datei für Begründung/Durchsetzung/
  Grenzen.

§5 wird zusätzlich von einer Fließtext-Liste zu einer **Index-Tabelle**
(Regel-Stichwort, Ein-Satz-Kurzfassung, Datei-Link) — dieselbe Form wie die
„Aktive Adaptionen"-Tabelle in `harness/conventions.md`.

Kein Wortlaut geht verloren: Der ausgelagerte Text wird **vollständig**
(inklusive seiner `d-check:cite`- und ID-Referenzen) in die neue Datei
kopiert, nicht zusammengefasst.

### `file[].max-lines: 400` für `AGENTS.md` und `harness/README.md`

Zwei `file[]`-Regeln in `.d-check.yml`, je Ziel-Datei eine (das Modul kennt
kein Oder-Glob). `harness/README.md` (233 Zeilen) unterschreitet die
Schwelle bereits und braucht keine Kürzung — die zweite Regel hält nur
künftiges Wachstum derselben Ladebudget-Klasse.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **`structure[].max-lines` je Abschnitt** statt `file[].max-lines` | feingranularer (Budget pro `###`-Abschnitt statt für die ganze Datei) | zählt den **bereinigten** Text, nicht die geladene Rohdatei — verfehlt genau die Beobachtung; bräuchte einen Abschnitts-Selektor je Regel statt einer Datei-weiten Zusage |
| **Nur Chronik/Historie kürzen, Struktur unverändert** | kein neues Verzeichnis, kein Anker-Risiko | liefert laut Beobachtung nur einen kleinen Teil der nötigen Kürzung; der Rest ist Regeltext, keine Historie |
| **Auslagerung nach `harness/rules/`, Überschriften/Anker stabil** (gewählt) | folgt einem im Repo bereits etablierten Muster (`conventions.md`-Adaptions-Block, `harness/sensors/`); kein Nachzug für zitierende Bestandsdateien nötig, da Anker gleich bleiben | ein Implementer braucht für Begründung/Grenzen einen zweiten Klick — dieselbe Abwägung, die `conventions.md` für Adaptionen bereits getroffen hat |

## Konsequenzen

- Neues Verzeichnis `harness/rules/` (plus Unterverzeichnis
  `harness/rules/dokumentations-regeln/` für die §5-Liste).
- `AGENTS.md` schrumpft von 575 auf unter 400 Zeilen, ohne Wortlaut zu
  verlieren — nur der Ort wechselt.
- `.d-check.yml` bekommt zwei neue `file[]`-Regeln.
- Kein Konsument außerhalb von `AGENTS.md` und den neuen Dateien ändert
  sich: alle bestehenden Anker (`AGENTS.md#3.x`, `AGENTS.md §5`) bleiben
  gültig, weil die Überschriften stehen bleiben.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `d-check --enable file` | `AGENTS.md` ≤ 400 Zeilen | `make doc-check` / `make gates` |
| `d-check --enable file` | `harness/README.md` ≤ 400 Zeilen | `make doc-check` / `make gates` |
| `d-check` (Default-Module) | alle bestehenden Anker/Links auf `AGENTS.md#3.x`/`§5` lösen weiterhin auf | `make gates` |

## Re-Evaluierungs-Trigger

`AGENTS.md` nähert sich der 400-Zeilen-Schwelle erneut (z. B. bei
80 % / 320 Zeilen) — dann entweder eine weitere Regel auslagern oder die
Schwelle als bewusste Auftraggeber-Entscheidung anheben (`AGENTS.md` §3.6).
Ohne das: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-27 | Proposed → Accepted (`slice-239`, Auftraggeber-Entscheid) |
