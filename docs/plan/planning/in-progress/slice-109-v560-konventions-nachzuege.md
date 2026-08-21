# Slice slice-109: v5.6.0-Konventions-Nachzüge — Etappe C-2…C-6 des Baseline-Bumps

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** welle-78-baseline-v560-migration.

**Bezug:** Audit-Findings C-2…C-6 aus
[slice-107](../done/slice-107-baseline-v560-delta-audit.md) §9. Reine
Harness-/Konventions-Doku; kein `DC-*`-Produktvertrag berührt.

**Autor:** pt9912. **Verantwortlich:** pt9912. **Datum:** 2026-08-21.

---

## 1. Ziel — fünf mechanische Angleichungen

1. **C-2 — ID-Schema-Deklaration in zwei Teilen** (Review-Befund F-4 zum
   Audit): (a) die **Vergabe** (dichte Nummern, ein Schreiber — kein
   Bereichssegment) gehört in die
   [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)-ID-Schema-Aussage;
   (b) der **Verzicht auf Struktur-IDs** ist eine echte Abweichung
   (`modul-03` schreibt `SPEC-*` für die Sektionstypen der Spezifikation
   §2–§6 vor) und wird als eigene Adaption **MR-027** deklariert <!-- d-check:ignore -->
   (dichte `.a`-Verfeinerungen + `§`-Anker; Begründung: von den 44
   FA-Anforderungen tragen 36 eine eigene Verfeinerungs-Sektion (37
   Sektionen inkl. einer `.b`) — sie decken die Adressierbarkeit der
   technischen Verfeinerungen, der Rest adressiert per `§`-Anker; ein
   SPEC-Retrofit
   über die gewachsene Spezifikation hätte keinen Konsumenten;
   Auflösungs-Trigger: eine ADR kann ihr `Schärft:`-Ziel nicht mehr eindeutig
   per `§`-Anker adressieren).
2. **C-3 — Kommentar-Regel-Träger:** AGENTS-Hard-Rule als Pointer auf
   [Baseline §Was ein Kommentar trägt](../../../../.harness/baseline/v5.6.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte)
   + HIGH-Eintrag „Kommentar trägt keine der fünf Klassen" in `reviewer.md`
   (Version 1.4.0 → 1.5.0; neuer HIGH-Eintrag mit Auflösungs-Trigger oder
   „permanent"). Dazu going-forward: Slice-Köpfe tragen `Verantwortlich:`
   (AGENTS §5-Halbsatz; kein Retrofit). **Und die gemessene Rest-Klasse
   räumen** (Review-Befund F-9): nackte Review-Finding-Tokens als
   Herkunfts-Feld — per grep alle Fundstellen (u. a. `cli.go`, `planning.go`,
   `configyaml.go` sowie die in welle-77 geschriebenen in
   `structure_tableorder.go`/`markdown.go`) auf eine der drei
   zugelassenen Formen ziehen oder streichen; Verhalten unverändert.
3. **C-4 — Kennungs-Anker im MR-Index:** je Index-Zeile zusätzlich
   `<a id="mr-<NNN>">`; die Voll-Slug-Anker bleiben als Migrations-Schuld
   stehen (Baseline: beim Formwechsel trägt die Zeile den alten Slug
   **zusätzlich**). Neue Verweise nutzen die Kennungs-Form.
4. **C-5 — Leseordnung** in `harness/README.md`: drei bis fünf geordnete
   Zeiger (Menschen-Hälfte des Einstiegs; kein Closure-Prüfpunkt).
5. **C-6 — Bestands-Stichprobe** dieses Bumps: ein delta-freier Abschnitt
   (**`modul-14-docker-harness.md`** — per Kurs-Diff verifiziert delta-frei;
   der erste Kandidat `modul-07` änderte sich in v5.2.0, Review-Befund F-2)
   je Regel gegen die eigene Verkörperung geprüft;
   Ergebnis hier notiert (Frage je Regel: im Artefakt oder als deklarierte
   Abweichung? Zweimal nein = nie übernommen ⇒ Weg jeder Diskrepanz).

## 2. Definition of Done

- [ ] Alle fünf Punkte umgesetzt bzw. (C-6) ausgeführt und notiert;
      `make gates` grün; unabhängiger Review; kein Release (Harness).

## 3. Trigger

**Start** (`open` → `in-progress`):
[slice-108](../done/slice-108-roadmap-offene-wellen.md) in `done/` (Reihenfolge:
erst die Struktur-Frage, dann die Nachzüge — C-4/C-5 schreiben in Dateien,
die C-1 nicht berührt, aber der Index-Anker-Punkt soll auf dem End-Layout
landen) **und** WIP-Slot frei.

## 4. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Doku, GF.
- **Beobachtungen sichten** (Stand 2026-08-21: BEO-006 offen, 1×): BEO-006
  als Arbeitsregel (Index vor pfad-selektiven Commits prüfen);
  [`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
  gilt für C-2 (die ID-Schema-Aussage spiegelt sich in [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage) und AGENTS §5).

## 5. Sub-Area-Modus-Begründung

**GF (Repo-Default)** — Konventions-Nachzug an die adoptierte Baseline.

## 6. Closure-Notiz (nach `done/`)

_Ausstehend._
