# Welle welle-89: Wellenlose Review-Archivierung nachgerüstet

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach in `docs/plan/planning/`; bei der Closure wandert sie nach `done/` und
bekommt ihre Ergebnisnotiz `welle-89-results.md` daneben.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-04.

---

## 1. Welle-Ziel

**`docs/reviews/` trägt ~65 Review-Reports, die keiner archivierten Welle
angehören** — Reports zu wellenlos geschlossenen Slices (slice-102 bis
slice-189 u. a.), für die `tools/archive-wave` bisher keinen Träger hat: das
Werkzeug verlangt zwingend eine `-welle`-Kennung.

**Der Baseline-Bump auf `v6.0.0` (welle-88) hat die Kanon-Lücke geschlossen,
die dieses Repo im CR
[`2026-09-03-cr-ai-harness-course-wellenlose-review-archivierung.md`](../../cr/2026-09-03-cr-ai-harness-course-wellenlose-review-archivierung.md)
gemeldet hatte:** `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht
nennt jetzt explizit eine Form für den wellenlosen Fall — die Slice-Closure
selbst archiviert, Schlüssel ist der Slice: `done/slice-<NNN>-archiv.zip`,
flach neben dem Stub (statt eines Wellen-Verzeichnisses). Das war zuvor eine
benannte, bewusst nicht lokal gelöste Lücke (Auftraggeber-Entscheid: auf die
Kanon-Antwort warten statt selbst bauen) — jetzt ist sie umsetzbar.

**Zwei Slices, nach demselben Muster wie welle-87 (Werkzeug bauen, dann auf
den Bestand anwenden):**

- **slice-196** erweitert `tools/archive-wave` um einen zweiten,
  wellenlosen Einzel-Slice-Modus (`-slice=slice-<NNN>` statt `-welle`):
  Volltext + Review-Reports eines einzelnen `done/`-Slice nach
  `done/slice-<NNN>-archiv.zip`, gekürzter Stub an seiner Stelle,
  Verweis-Nachzug wie im Wellen-Modus. An einem konstruierten Fixture
  bewiesen, nicht am echten Bestand.
- **slice-197** wendet den neuen Modus auf **alle** wellenlos geschlossenen
  `done/`-Slices an, deren Review-Reports noch in `docs/reviews/` liegen
  (Zahl bei Ausführung exakt erhoben, ~43 erwartet).

## 2. Trigger (Welle startet)

Kein Vorwellen-Trigger nötig — der Anlass ist die eigene, bereits gemeldete
Beobachtung samt jetzt vorliegender Kanon-Antwort. `in-progress/` ist leer,
das WIP-Limit ist frei.

## 3. Closure-Trigger (Welle schließt)

- slice-196, slice-197 beide in `done/`.
- `docs/reviews/` enthält keinen Report mehr zu einem wellenlosen,
  archivierbaren `done/`-Slice (Reports zu `done/welle-*/`-archivierten
  Slices bleiben unberührt — sie liegen bereits im jeweiligen `archiv.zip`).
- `make gates` und `make fullbuild` grün auf dem Endstand.
- `make archive-wave-test` grün (Fixture-Beweis des neuen Modus).
- Closure-Notiz `welle-89-results.md` geschrieben.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| slice-196 | `tools/archive-wave`: wellenloser Einzel-Slice-Archivierungs-Modus | Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht |
| slice-197 | Wellenlosen Review-Bestand archivieren (~43 Slices) | slice-196 |

## 5. Abhängigkeiten

- slice-197 wird blockiert von slice-196 (braucht den neuen Werkzeug-Modus).

## 6. Out-of-Scope für diese Welle

- Eine Anpassung des Wellen-Modus (`-welle`) — bleibt unverändert.
- Reviews zu Slices, die bereits über eine Welle archiviert sind (liegen
  schon in einem `archiv.zip`) — nur der wellenlose Rest-Bestand.
- Eine rückwirkende `**Welle:**`-Feld-Vergabe für Slices, die tatsächlich
  wellenlos sind — sie bleiben `ohne Welle`, das ist ihr korrekter Zustand
  (Baseline-Regelwerk „wellenlos heißt nicht wächterlos").

## 7. Closure-Notiz

<!-- wird erst nach Welle-Abschluss gefüllt -->

Ergebnis: `done/welle-89-results.md`
