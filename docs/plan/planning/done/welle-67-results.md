# Welle 67 — Baseline-Regelwerk-Migration v1.4.0 → v5.0.0 — Closure-Notiz

**Welle:** welle-67-baseline-v500-migration
**Abschluss:** 2026-08-03
**Verantwortlich:** pt9912

## Was wurde geliefert?

Das adoptierte Baseline-Regelwerk ist **vollständig** von `v1.4.0` auf `v5.0.0` gehoben —
über acht Slices (084–091), je mit unabhängigem Frischkontext-Review vor Closure:

- **Etappe A** (slice-084): self-contained v5.0.0-Bundle committet vendored (beide Bäume
  `{regelwerk,templates}`), Pin auf `v5.0.0`, das Materialisierungs-Skript aufs Bundle
  gehoben, drei eingefrorene Verweise via geteiltes `ignore-refs` getombstoned.
- **Etappe B** (slice-085): Modul-Delta gelesen — 18 Findings (8 → Etappe C, 11 → D);
  Kern: die Historie-Provenance-Ausnahme ist revoziert.
- **Etappe C** (slice-086): Konventionsspeicher auf **Index + Datei je MR** template-konform
  (8 aktiv / 15 aufgelöst); 188 Links über 12 immutable ADRs erhalten, keine ADR berührt.
- **C-3-Nachzug** (slice-087): Spec-§7-Referenzrichtung **entkoppelt**, `matrix.exclude-sections`
  auf `[Geschichte]` verengt; eine neue Accepted-ADR (Supersede-Verfeinerung) hält den Schnitt.
- **Etappe D** (slice-088–091): **Planning-Layer** (dieses Wellendokument + das
  Beobachtungs-Register + Roadmap auf sechs Baseline-Abschnitte), **AGENTS.md**-Currency +
  ADR-Re-Evaluierungs-Trigger-Konvention + Slice-Vorprüfungen, **reviewer.md**-Currency +
  Report-Kopffelder + Finding-Feld `klasse`, und das Slice-`Status:`-Feld go-forward
  abgeschafft (Lifecycle = Verzeichnis).

## Was hat funktioniert?

- Der in slice-083 **abgenommene Etappen-Schnitt** (A→B→C→D) trug durch; die Mini-Wellen-
  Aufteilung der Etappe D hielt gemischte (Doc/Prozess/Adoptions-)Naturen sauber getrennt.
- Der **GUARD** (netzloser Selbst-Scan vor **jedem** Commit) fing die id-unlinked-/
  codepath-/hostpath-Fallen früh, nicht erst in der CI.
- **Adversariale Frischkontext-Reviews** je Slice fingen echte Befunde (u.a. ein MEDIUM
  AGENTS-§3.3-Spiegel, mehrere hängende Prosa-Reste), die kein Gate sieht.

## Was ging anders als geplant?

- **C-3 war KEIN Code-Feature.** Die slice-086-Messung war zu grob (ganze Liste statt
  namens-selektiv); die Heading-Namen `Geschichte`/`7. Historie` trennen bereits sauber →
  chirurgischer Konfig-Schnitt (slice-087). Folge: der geplante C-3-Folge-Produkt-Slice entfiel.
- **Etappe D neu geordnet** (Nutzer-Hinweis: welle-67 lief ohne ihr Wellendokument) →
  Planning-Layer zuerst statt Doc-Form.
- **D-7** (Closure-Note-Reviewer + Gate) als **Folge-Produkt-Slice** herausgeschnitten
  (Code + ADR); **D-5** als **template-forward** umgesetzt (kein Retrofit der 90 done-Slices)
  — beide Nutzer-Entscheide.

## Steering-Loop-Einträge

- **— keine — bei 3×.** Die wiederkehrenden Lehren dieser Welle sind je in der Slice-§9-
  Closure-Notiz und (wo tragend) in `AGENTS.md`/Gates/Konventionen verkörpert; keine
  Beobachtung hat im Register die 3×-Schwelle erreicht (das Register entstand erst in
  dieser Welle und startete leer).

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`../observations.md`](../observations.md) (`— keine —`). Eine offene
Beobachtung ist notiert (Baseline-**interne** 5-vs-6-Finding-Feld-Drift, slice-090-§9) —
sie ist **upstream** und für den nächsten Baseline-Bump-Drift-Audit vermerkt, nicht
d-check-verkörperbar.

## Folge-Slices

- **— keine mit Datei —.** Ein Kandidat steht in der Roadmap §Nächste Wellen: der
  **Closure-Note-Reviewer + `verify-closure-notes`-Gate** (D-7, Code + eigene ADR) — noch
  **ohne** Slice-Datei, daher kein benannter Folge-Slice.

## Verifikation

- `make gates` grün (doc-check + lint + test + arch-check + coverage-gate + semgrep +
  gate-consistency + planning-check).
- `make adr-check` grün — **keine** `Accepted`-ADR über die Welle inhaltlich berührt (die
  eine neue ADR aus dem C-3-Nachzug ist neu Accepted, kein Edit einer bestehenden).
- Alle acht Slices (084–091) liegen in `done/`; der Baseline-Pin ist `v5.0.0`, beide Bäume
  `{regelwerk,templates}` sind vendored (`.harness/baseline/v5.0.0/`).
