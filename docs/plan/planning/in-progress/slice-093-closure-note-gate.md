# Slice slice-093: Closure-Note-Reviewer-Skill + `verify-closure-notes`-Gate (D-7)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** welle-68-planning-roadmap-harness (zweiter Slice, nach slice-092).

**Bezug:** Etappe-D-Finding **D-7** (aus slice-090 als Folge-Produkt-Slice
herausgeschnitten): der dedizierte **Closure-Note-Qualitäts-Nachlauf** fehlt in d-check
komplett. Die Baseline vendored beide Ziel-Formen: `review-report.template.md` und
`closure-note-reviewer.template.md` (der *inferentielle* Skill **über** einem
*strukturellen* `verify-closure-notes`-Gate). **Produkt-/Harness-Feature ⇒ CR + eigene
ADR + Spec + Code + Tests + Release.**

**Autor:** pt9912. **Datum:** 2026-08-03.

---

## 1. Ziel

Den Closure-Note-Qualitäts-Nachlauf **mechanisieren**: (a) ein **strukturelles Gate**
(`make verify-closure-notes`) prüft die Closure-Notizen der `done/`-Slices maschinell
(Heading vorhanden · Mindest-Satzzahl außerhalb Code-Blöcken · Floskel-Liste); (b) der
**inferentielle** `closure-note-reviewer.md`-Skill deckt darüber, was Struktur allein
nicht fängt (Inhalt vs. Floskel — semantisch, Reviewer-Sache).

## 2. Vorgehen

1. **Design-Entscheid (Abnahme-Punkt §3).** Wie wird das Gate gebaut — als **neues
   Go-Regelmodul** (d-check-Identität: `tools/*.sh` → verteilbare Go-Module, wie
   `planning`/`commits`/`vcs`) vs. Erweiterung des `planning`-Moduls; die genaue
   Prüf-Semantik; die Bereichs-/`DC-FA-*`-Kennung; ADR-Nummer.
2. **CR + Spec.** Change Request in `spec/lastenheft.md` (neue `DC-FA-*`-Anforderung) +
   `spec/spezifikation.md`-`.a`-Algorithmus + §4-Grund-Code(s).
3. **ADR.** Entscheidungs-Record (neues Gate; Modul-Schnitt; Struktur-vs-Inferenz-Grenze).
4. **Implementierung.** Go-Modul(-Fähigkeit) + `Makefile`-Target `verify-closure-notes` +
   der `.harness/skills/closure-note-reviewer.md`-Skill (aus `closure-note-reviewer.template.md`
   an den Haus-Stil angepasst); Tests (Golden/Mutation, fail-closed).
5. **Release.** Voller Release-Flow (Push + Tag + GHCR), Digest-Backfill.
6. **Gate.** `make gates` + das neue `make verify-closure-notes` grün (repo-weit über den
   `done/`-Bestand — das *Mehr* der welle-68-Closure); unabhängiger Frischkontext-Review.

## 3. Abnahme-Punkte

1. **Gate-Bauform + Modul-Schnitt.** (a) Neues Go-Regelmodul; (b) Erweiterung des
   `planning`-Moduls; (c) `check_closure_notes.py`-Skript wie im Template.
   → **Entschieden 2026-08-03: (b)** — das **`planning`-Modul** um eine
   **Closure-Note-Struktur-Fähigkeit** erweitern (Nutzer-Entscheid): schärft
   [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
   + eigene ADR; **kein** neues Modul. **Offen (im Slice zu entwerfen):** die exakte
   Struktur-Semantik (Heading-Pflicht · Mindest-Satzzahl außerhalb Code-Blöcken ·
   Floskel-Liste), der Grund-Code, das Config-Feld, ob `make verify-closure-notes` ein
   eigenes Target ist, und ob der `closure-note-reviewer.md`-Skill an das Struktur-Gate
   **oder** an den allgemeinen `reviewer.md` koppelt.

## 4. Definition of Done

- [ ] Design-Entscheid (Abnahme-Punkt 1) festgehalten; CR + `DC-FA-*` + ADR angelegt.
- [ ] `make verify-closure-notes`-Gate implementiert (Struktur) + `closure-note-reviewer.md`-
  Skill (Inferenz); Tests fail-closed.
- [ ] Release (Tag + GHCR + Digest-Backfill).
- [ ] `make gates` + `make verify-closure-notes` grün; unabhängiger Frischkontext-Review.

## 5. Risiken / offene Punkte

- **Produkt-Code + Release** — grösster Slice der welle-68; der volle Release-Flow ist Teil.
- **Struktur-vs-Inferenz-Grenze** sauber ziehen (das Gate prüft Struktur, der Skill
  Inhalt) — sonst Doppelbefund oder Lücke.
- **Dogfood:** d-check muss sein eigenes `verify-closure-notes` über die eigenen
  `done/`-Closure-Notizen grün bekommen (ggf. Nachbesserung bestehender Notizen).

## 6. Trigger

Freigabe des aus slice-090 herausgeschnittenen D-7 als welle-68-Slice.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** berührt mehrere Sub-Areas — Produkt-Code (`internal/`,
  `planning`-Modul), Spec (`spec/`) und Harness/Skills (`.harness/skills/`). **Alle unter
  dem Repo-Default GF** (`harness/conventions.md` §Modus: `*` = Greenfield). **Kein
  Brownfield:** der Code entsteht **spec-first**, wird nicht rückdokumentiert.
- **Offene Beobachtungen sichten:** `observations.md` = `— keine —`; nichts zu berücksichtigen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — „Spec führt, Code folgt: wir versprechen X, dann
liefern wir X". slice-093 ist ein **Produkt-Feature** (neues Gate: Go-Code + Spec +
Release), aber der **Modus** ist GF wie bei jedem d-check-Slice: der Change Request +
die Schärfung der `planning`-Anforderung + die ADR schreiben die **Zusage**, die
Go-Implementierung **liefert** sie. Der Unterschied zu den Doc-only-Migrations-Slices ist der **Scope**
(Produkt-Code + Release statt reiner Doc/Harness), **nicht** der Modus — Brownfield wäre
nur die Inventur bestehenden undokumentierten Codes, was hier nicht vorliegt.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
