# Slice slice-026: Reparatur-Patch (`--repair`)

**Status:** open (geplant).

**Welle:** welle-15-doctor-repair (Trigger: slice-025 done — baut auf
dem dort eingeführten Fix-Kandidaten-Modell auf).

**Bezug:**
[`DC-FA-CLI-008`](../../../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch)
(Hauptanforderung),
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
(Fix-Kandidaten-Mechanik wird wiederverwendet),
[`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)
(Exit-Codes der Modi),
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
(Patch ersetzt Default-stdout; nicht mit `--json` kombinierbar),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(byte-identischer Patch je Stufe),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(read-only; Patch nur auf stdout, kein In-place).

**Autor:** pt9912. **Datum:** 2026-06-18.

---

## 1. Ziel

`d-check --repair` gibt einen **unified diff auf stdout** aus, der
ableitbare Befunde behebt (`git apply`-kompatibel); das Werkzeug schreibt
selbst **nichts**. **Zwei Stufen** per Schalter: **konservativ** (Default,
nur eindeutige Fixes — v1 v. a. `id-unlinked` → Definitions-Link) und
**breit** (opt-in, Best-Guess wie `target-missing` → nächstliegende
Überschrift/Datei). Die Best-Guess-Markierung erscheint auf **stderr**,
damit der Patch auf stdout `git apply`-rein bleibt. Baut auf dem
Fix-Kandidaten-Modell aus slice-025 auf.

## 2. Definition of Done

- [ ] **Spezifikation** zu [`DC-FA-CLI-008`](../../../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch) (Spec-Abschnitt `…008.a` entsteht in dieser Slice): unified-diff-Rendering
  (`git apply`-Format), Stufen-Schalter (Name + Werte), Zuordnung
  Grund-Code → Stufe (konservativ vs. breit), stderr-Kanal der
  review-pflichtig-Markierung, Determinismus (gleicher Input + Stufe →
  byte-identisch), saubere Anwendbarkeit gegen den unveränderten Baum.
- [ ] **Implementierung:** Flag `--repair` + Stufen-Schalter;
  Patch-Renderer, der die slice-025-Fix-Kandidaten in Hunks übersetzt;
  Konservativ-Filter (nur eindeutige Kandidaten).
- [ ] **Round-Trip-Beleg** (Happy): Patch → `git apply` → erneuter
  `d-check`-Lauf meldet den `id-unlinked`-Befund nicht mehr.
- [ ] **Stufen-Verhalten** (Boundary): konservativ leerer Patch bei nur
  best-guess-fähigen Befunden (`target-missing`); breite Stufe markierte
  Hunks, Kennzeichnung auf stderr; Patch bleibt `git apply --check`-rein.
- [ ] **Negativtest** ([`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)):
  ungültige Stufen-Wahl bzw. `--repair --json` → exit 2.
- [ ] **Read-only-Beleg** ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit))
  und **Determinismus-Beleg** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- [ ] **Doku** unter `docs/user/`; `make gates` grün; Closure-Notiz.

## 3. Plan (vor Code)

| Datei | Art | Begründung |
|---|---|---|
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | update | Spec-Abschnitt zu [`DC-FA-CLI-008`](../../../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch): Diff-Rendering, Stufen, stderr-Markierung, Determinismus |
| `internal/…` (CLI + Patch-Renderer) | update/neu | Flag, Stufen-Schalter, Kandidaten→Hunk, Konservativ-Filter (Wiederverwendung der slice-025-Funktion) |
| `docs/user/operations.md` | update | Option `--repair` + Stufen dokumentieren |

Das Lastenheft ([`DC-FA-CLI-008`](../../../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch)) ist mit Change Request 0.16.0 bereits
gesetzt — kein weiterer Vertrags-Change (außer einer Schärfung, falls die
Spezifikation eine aufdeckt).

## 4. Trigger

slice-025 done (Fix-Kandidaten-Modell vorhanden und stabil).

## 5. Closure-Trigger

DoD vollständig inkl. Round-Trip-Beleg, `git apply --check`-Reinheit,
Determinismus-/read-only-Beleg und grüner Gates.

## 6. Risiken und offene Punkte

- **`git apply`-Reinheit:** Marker dürfen **nicht** in den Patch (→
  stderr), sonst bricht `git apply`. Mit `git apply --check` auf dem
  emittierten Patch absichern (deckt den Review-Punkt INFO-4 ab).
- **Best-Guess-Treffsicherheit:** die breite Stufe rät (z. B.
  `target-missing` → falsche Überschrift) — deshalb review-pflichtig;
  nie als „korrekt" verkaufen (Vertrag, kein Orakel — analog slice-020).
- **Konservativ-Disziplin:** die Default-Stufe darf nur wirklich
  eindeutige Fixes emittieren, sonst entsteht ein stiller Fehl-Patch ohne
  Review — derselbe „grün ≠ richtig"-Pfad, den der Reviewer-Skill
  bewacht.
- **Determinismus** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  stabile Hunk-Reihenfolge und Kontextzeilen — gleicher Input + Stufe →
  byte-identischer Patch.

## 7. Closure-Notiz (nach `done/`)

*(wird bei Closure gefüllt — Umsetzung, Belege, Lerneintrag, Review-Runde.)*

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Spec-/Code-/Doku-Arbeit; Greenfield-Default).
