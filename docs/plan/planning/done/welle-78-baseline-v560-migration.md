# Welle welle-78-baseline-v560-migration: Baseline-Regelwerk v5.0.0 → v5.6.0

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-78-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Adoptions-Pflege der Baseline).

**Verantwortlich:** pt9912. **Datum:** 2026-08-21.

---

## 1. Welle-Ziel

Die adoptierte Baseline vollständig von `v5.0.0` auf **`v5.6.0`** heben
(Kurs-Tag 2026-08-16; sechs additive Stufen, 20 Regelwerks-Dateien,
+902/−152 Zeilen) — nach der welle-67-Präzedenz in Etappen: **A** Vendoring,
Pin und Verweis-Hebung ([slice-106](slice-106-baseline-v560-vendoring.md)),
**B** Stufen-Audit je Regel gegen die Tag-Notizen
([slice-107](slice-107-baseline-v560-delta-audit.md)), **C** der
Konformitäts-Nachzug — dessen Slices werden aus dem B-Befund geschnitten
(Roadmap-Drift-Log dokumentiert den Zuwachs). Der größte B-Gegenstand ist die
**Team-Fähigkeit** (v5.5.0): ob sie ein Ein-Operator-Repo bindet oder eine
deklarierte Adaption verlangt, entscheidet der Audit, nicht die Annahme.

## 2. Trigger (Welle startet)

Freigabe des Auftraggebers (angekündigt 2026-08-16, Reihenfolge-Entscheid
2026-08-21: nach der Chronologie-Welle), WIP-Slot frei (welle-77 geschlossen,
`in-progress/` trägt nur die Roadmap). Der Kurs-Tag v5.6.0 samt
`lab-regelwerk.zip`-Release-Asset liegt vor (verifiziert 2026-08-21).

## 3. Closure-Trigger (Welle schließt)

- Alle Etappen-Slices in `done/` — [slice-106](slice-106-baseline-v560-vendoring.md),
  [slice-107](slice-107-baseline-v560-delta-audit.md) und die aus dem
  B-Befund geschnittenen Etappe-C-Slices (§4 wird beim B-Abschluss
  nachgeführt).
- Pin `v5.6.0` vendored, `--verify` offline grün **und** `--check-latest`
  ohne Currency-/Content-Drift-Befund.
- Kein **lebender** Verweis nennt mehr `baseline/v5.0.0`; die eingefrorenen
  sind quell-skopiert getombstoned (`make doc-check` belegt beides).
- Der Stufen-Audit trägt **je Regel** eine Antwort (konform · angepasst ·
  nicht anwendbar mit Begründung) — kein pauschales Konform.
- `make fullbuild` grün. **Release nur, falls Etappe C Produkt-Verhalten
  ändert** (dann Minor) — die Hebung selbst ist Harness, nicht Produkt
  (Präzedenz welle-67: kein Release).

## 4. Slices in dieser Welle

| Slice | Rolle |
|---|---|
| [slice-106](slice-106-baseline-v560-vendoring.md) | Etappe A: Bundle v5.6.0 vendored, Pin + MR-026, Verweis-Hebung, Tombstones | <!-- d-check:ignore -->
| [slice-107](slice-107-baseline-v560-delta-audit.md) | Etappe B: Stufen-Audit v5.1.0–v5.6.0, je Regel eine Antwort; schneidet Etappe C |
| [slice-108](slice-108-roadmap-offene-wellen.md) | Etappe C-1: Roadmap auf §Offene Wellen (Form + `planning`-Config + [`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform)-Entscheid) |
| [slice-109](slice-109-v560-konventions-nachzuege.md) | Etappe C-2…C-6: ID-Schema-Deklaration, Kommentar-Regel-Träger, Kennungs-Anker, Leseordnung, Bestands-Stichprobe |

*Nachgeführt beim B-Abschluss (Drift-Log-Eintrag), wie bei der Eröffnung
angekündigt.*

## 5. Abhängigkeiten

- Kurs-Repo-Tag `v5.6.0` mit Release-Asset `lab-regelwerk.zip` (liegt vor).
- Das Materialisierungs-Skript
  (`tools/harness/fetch-baseline-cache.sh`, [`MR-023`](../../../../harness/conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)-Layout)
  nimmt ein explizites Tag-Argument — der erste Vendor-Lauf braucht den
  neuen Pin nicht.
- Etappe B liest den **vendorten** Baum (netzlos), nicht das Kurs-Repo —
  darum A vor B (bindende Reihenfolge).

## 6. Out-of-Scope für diese Welle

- **Kein pauschales Übernehmen** neuer Kurs-Artefakte ohne Audit-Antwort —
  jede Stufe wird gelesen, bevor etwas entsteht.
- **Kein Retrofit** eingefrorener Artefakte (immutable ADRs, `done/`-Slices,
  Review-Reports) auf neue Formen — dieselbe template-forward-Disziplin wie
  bei welle-67 (D-5).
- **Keine Produkt-Feature-Arbeit** außerhalb dessen, was der Audit als
  Etappe C verlangt.

## 7. Closure-Notiz

Geschlossen am 2026-08-21. Alle Closure-Trigger erfüllt: die vier
Etappen-Slices liegen in `done/`, der Pin v5.6.0 ist vendored und beidseitig
auditiert (`--verify` 51 Dateien, `--check-latest` Currency + Content OK),
kein lebender Verweis nennt den entfernten v5.0.0-Baum, der Stufen-Audit
trägt je Regel eine Antwort, und der Release-Punkt ist ehrlich entschieden
(kein Release — Klartexte, kein Vertrag; Lesart in der Ergebnisnotiz).
Die Migration hat zwei eigene Beobachtungs-Klassen erzeugt (BEO-006 auf 2,
BEO-007 neu bei 2) und die Baseline-Disziplin an sich selbst erprobt: der
Widerspruchs-Ausgang entschied den Struktur-ID-Fall, die Spiegel-Regel fand
die Lücken der eigenen Umbauten. Was wirkte und was anders lief:
[welle-78-results.md](welle-78-results.md).
