# Welle welle-67-baseline-v500-migration: Baseline-Regelwerk-Migration v1.4.0 → v5.0.0

**Lifecycle:** Diese Datei entstand — nachgezogen, weil die Welle bereits lief — bei
der Dokumentation der laufenden Welle und liegt **flach** unter
`docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/` (neben ihre
`welle-67-results.md`). Der Zustand ist die **Verzeichnis-Position** — kein
Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (interne Harness-/Spec-Konformität, kein
Release).

**Verantwortlich:** pt9912. **Datum:** 2026-08-02.

---

## 1. Welle-Ziel

Das adoptierte Baseline-Regelwerk ist von der gepinnten Version `v1.4.0` auf `v5.0.0`
gehoben, **vollständig** — Vendoring, Modul-Delta, Konventionsspeicher und
Form-Konformität. Leitsatz (Auftraggeber): **der Baseline-Default sticht jede
repo-lokale Adaption.** Beobachtbar erfüllt, wenn alle Etappen-Slices in `done/` liegen
und `make gates` + `make adr-check` grün sind — bei erhaltener 188-Link-Integrität und
**ohne** inhaltliche Berührung einer `Accepted`-ADR.

## 2. Trigger (Welle startet)

- [slice-083](done/slice-083-regelwerk-v500-migration-analyse.md) (v5.0.0-Migrations-
  Analyse, §2.7-Etappen-Schnitt) **abgenommen** — der Start-Trigger liegt **vor** der
  Welle, ist kein Ergebnis dieser Welle.

## 3. Closure-Trigger (Welle schließt)

- Alle Slices dieser Welle (084–091) liegen in `done/`.
- `make gates` grün **und** `make adr-check` grün (repo-weit — steht in keiner
  einzelnen Slice-DoD; das ist das *Mehr* gegenüber den Slice-DoDs).
- **Trigger-Audit** durchlaufen (`modul-06` Closure-Schritt 2): Carveout · bootstrap-
  aware Gate · ADR-Re-Evaluierungs-Trigger — je bestätigt/aufgelöst/permanent (aktuell
  **0 aktive Carveouts**, latent).
- Closure-Notiz `done/welle-67-results.md` geschrieben.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [slice-084](done/slice-084-etappe-a-vendoring.md) | Etappe A — Vendoring (self-contained v5.0.0-Bundle, Pin/Pointer/Tombstone) | Migrations-Schnitt §2.7 A |
| [slice-085](done/slice-085-etappe-b-modul-delta.md) | Etappe B — Modul-Delta (18 Findings, 8→C / 11→D) | Migrations-Schnitt §2.7 B |
| [slice-086](done/slice-086-etappe-c-mr-bereinigung.md) | Etappe C — Konventionsspeicher (Index + Datei je MR) | Migrations-Schnitt §2.7 C |
| [slice-087](done/slice-087-spec-historie-referenzrichtung.md) | C-3-Nachzug — Spec-§7-Referenzrichtung entkoppelt | Finding C-3 |
| [slice-088](in-progress/slice-088-etappe-d1-doc-form.md) | Etappe D — Planning-Layer-Form (Roadmap · Welle · Register) | Findings D-1/D-2/D-3/D-4/D-9 |
| slice-089 | Etappe D — Doc-Form (AGENTS.md · ADR-Trigger) | Findings D-8/D-11 |
| slice-090 | Etappe D — Review-Infrastruktur | Findings D-6/D-7/D-10 |
| slice-091 | Etappe D — Slice-`Status:`-Feld entfernen | Finding D-5 |

## 5. Abhängigkeiten

- **Serielle Etappen-Ordnung:** A → B → C → C-3-Nachzug → D; D ist als Mini-Welle
  intern seriell (088 → 089/090/091, in beliebiger Reihenfolge nach 088).
- **Wird blockiert von:** nichts Externem (keine Fremd-Welle, kein Release-Gate).
- **Blockiert:** nichts — welle-67 ist der letzte offene Migrations-Bogen.

## 6. Out-of-Scope für diese Welle

- **Kein Produkt-Release, kein neues Produkt-Feature.** Reine Harness-/Spec-/
  Konventions-Konformität (die einzige Ausnahme wäre gewesen ein C-3-Code-Feature —
  entfiel, weil die Heading-Namen bereits selektiv trennen; siehe slice-087).
- **Kein retroaktives Wellen-*Plandokument* für frühere Wellen** (welle-01…welle-66):
  dieses Wellendokument ist der Neu-Anfang der Lifecycle-Adoption. Für die im Roadmap-
  Closure-Log geführten Alt-Wellen (welle-60…66) werden in slice-088 **minimale,
  retroaktiv markierte** `done/welle-NN-results.md` nachgezogen, damit das
  `## Abgeschlossene Wellen`-Log auflöst (Nutzer-Entscheid) — kein volles Plandokument.
- **Kein Umbau des `planning`-Moduls** (der Ruhe-Marker bleibt als deklarierte
  Adaption; ein Produkt-Code-Umbau wäre ein eigener Change).

## 7. Closure-Notiz

_Ausstehend (wird bei der Welle-Closure gefüllt; die Pointer auf `welle-67-results.md`
und das Beobachtungs-Register werden dann so geschrieben, wie sie vom Ruheort `done/`
auflösen — Ruheort-Regel)._
