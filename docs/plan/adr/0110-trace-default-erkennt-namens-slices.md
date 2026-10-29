# ADR-0110: Die Default-Erkennung des Trace führt Slices mit Namens-Kennung

**Status:** Proposed

**Datum:** 2026-10-29

**Autor:** pt9912

**Bezug:** [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(Scope der RTM);
[`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(der Exit-Code, an dem das Feld es bemerkt hat);
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Byte-Identität für den Ziffern-Fall);
[ADR-0108](0108-trace-requirements-als-liste-von-quellen.md) (dieselbe Trace-Konfiguration);
slice-279 <!-- d-check:status-provenance -->.

**Schärft:** [`DC-FA-CLI-009.a`](../../../spec/spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix)
Schritt 2 und die Schema-Zeile `trace.slices.file-pattern`.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Trace leitet die Kennung eines Slice aus dem Dateinamen ab. Der Default
erkennt nur `slice-NNN-titel.md`. Ein Repo, dessen Slices `slice-qa-sicherheit.md`
heißen, bekommt jeden davon als Waise gemeldet, und `--trace --require-complete`
endet mit Exit 1, obwohl die Anforderungen gedeckt sind. Der Schlüssel
`trace.slices.file-pattern` löst das, ist aber im Default nicht sichtbar:
gemessen am Image hat dieselbe Datei mit gesetztem Muster `ok` ergeben, ohne es
`WAISE`.

Eine naive Erweiterung des Defaults auf jedes `slice-*.md` hätte eine
Fehlerrichtung, die schlimmer ist als die behobene: das Beobachtungs-Register
führt jeden Beleg als `evidence/slice-<Kennung>.md`. In diesem Repo sind das 270
von 549 Dateien `slice-*.md` unter `docs/plan/planning/`. Jeder trüge den Namen
eines Slice und deckte, was er nennt — die Erweiterung würde Waisen verdecken,
wo sie bisher Slices übersah.

## Entscheidung

### 1. Der Default kennt zwei Fälle

Ziffern-Fall zuerst: passt der Basisname auf `^slice-\d+-.*\.md$`, bleibt die
Kennung `slice-NNN` wie bisher. Sonst gilt der Namens-Fall: passt der Basisname
auf `^slice-[a-z0-9]+(?:-[a-z0-9]+)*\.md$`, ist die Kennung der Basisname ohne
`.md`. Ein Punkt im Stamm, Großschreibung und ein leerer Name sind keine Kennung.

### 2. Belege unter `observations` sind im Default keine Slices

Eine Datei, die unterhalb von `trace.slices.dir` in einem Verzeichnis
`observations` liegt, trägt im Default keine Kennung. Das Segment in
`trace.slices.dir` selbst zählt nicht.

### 3. Ein gesetztes `trace.slices.file-pattern` gilt unverändert

Wer ein Muster setzt, entscheidet die Menge selbst — ohne Namens-Fall und ohne
den `observations`-Ausschluss.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| Nur den Schlüssel auffindbar machen (Handbuch, Gerüst) | kein Verhalten ändert sich | der Default bleibt falsch für jedes Repo mit Namens-Slices; die Waise erscheint, bevor jemand das Handbuch liest |
| Eine Warnung für `slice-*.md`, die das Muster nicht trifft | macht den Fall sichtbar | neue Ausgabe-Klasse im `--trace`-Lauf, bricht die Byte-Identität des Ziffern-Falls; meldet die Ursache, statt sie zu beseitigen |
| Jedes `slice-*.md` ohne Ausschluss | einfachste Regel | Belege des Beobachtungs-Registers deckten Anforderungen; Waisen würden verdeckt |
| Gewählt: Ziffern-Fall zuerst, Namens-Fall, `observations` ausgenommen | Ziffern-Repos unverändert; Namens-Repos ohne Konfiguration korrekt | Notiz-Dateien mit `slice-`-Präfix außerhalb von `observations` zählen als Slice; Abhilfe ist ein gesetztes `file-pattern` |

**Fitness-Funktion:**

- Die Tests in `internal/hexagon/core/app` (`trace_slice_names_test.go`) fahren
  Namens-Slice, Ziffern-Fall, die Negativliste, Kennungs-Gleichheit und das
  gesetzte Muster. Der Namens-Test wird ohne den Namens-Fall und ohne den
  `observations`-Ausschluss je einzeln rot — jede Stelle gebrochen und gesehen.
- Die Akzeptanz am Binary hält den Ziffern-Fall byte-identisch und den
  Namens-Fall bei Exit 0 unter `--require-complete`.

## Konsequenzen

- Ein Repo mit Namens-Slices und ohne `trace`-Block bekommt eine andere RTM:
  Waisen werden zu `ok`. Das ist beabsichtigt, aber nicht byte-identisch zum
  Vorgänger; das Lastenheft fasst das Kriterium *Default byte-identisch* auf den
  Ziffern-Fall.
- **Grenze:** Ob eine `slice-*.md` ein Slice-Plan ist, entscheidet der Dateiname
  nicht. Wer Notizen mit `slice-`-Präfix im Planning-Baum führt, setzt
  `trace.slices.file-pattern`.
- **Grenze:** Die Review-Deckung (`reviews`), `--suggest-config` und
  `tools/archive-wave` tragen weiter ihre ziffern-only Form; der Befund betraf
  nur den Trace.

## Re-Evaluierungs-Trigger

- Ein Adopter meldet eine `slice-*.md`, die der Namens-Fall fälschlich als Slice
  zählt, und das gesetzte Muster ist keine zumutbare Abhilfe.
- Ein zweiter Ort neben `observations` trägt Dateien mit Slice-Namen ohne Slice
  zu sein.
- Die Review-Deckung oder `--suggest-config` bekommen einen Namens-Fall — dann
  ist die Erkennungsform zu vereinheitlichen.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-29 | Entwurf zu slice-279 nach dem Feldbefund (GitHub-Issue 11). Status Proposed, bis die Umsetzung geprüft ist. |
