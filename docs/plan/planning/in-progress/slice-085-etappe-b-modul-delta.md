# Slice slice-085: Regelwerk-Migration Etappe B — Modul-Delta lesen

**Status:** In Arbeit (welle-67).

**Welle:** welle-67-baseline-v500-migration (zweite Umsetzungs-Etappe, nach
[slice-084](../done/slice-084-etappe-a-vendoring.md)).

**Bezug:** Umsetzung von **Etappe B** des in
[slice-083](../done/slice-083-regelwerk-v500-migration-analyse.md) §2.7
abgenommenen Migrations-Schnitts. Reine **Lese-/Analyse-Etappe** — kein
Code/Config-Delta. **Kein Change Request**, **kein ADR**, **kein Release**.

**Autor:** pt9912. **Datum:** 2026-08-02.

---

## 1. Ziel

Das seit Etappe A **vendorte** v5.0.0-Regelwerk (8 `grundlagen-*` + 17 Module,
netzlos lesbar) gegen den d-check-Ist gegenlesen und die Regel-Deltas als
**Finding-Liste** sammeln. Diese Liste macht Etappe C (Adaptions-/Konventions-
speicher-Bereinigung) und Etappe D (Form-Konformität) verbindlich: sie sagt je
Treffer, **welche** v5.0.0-Regel d-check **wie** trifft und **wohin** die Handlung
gehört. Die in slice-083 §2.3 vorab identifizierten „Zugänge" werden dabei gegen
die **Quelle** bestätigt (oder korrigiert) und um am Text gefundene Deltas ergänzt.

## 2. Prozedur (aus slice-083 §2.7 Etappe B)

1. **Gegenlesen** — jede `grundlagen-*`-Datei und jedes Modul gegen `v5.0.0`,
   Priorität nach §2.1: substanziell umgeschrieben zuerst (grundlagen-Split;
   Module 2/5/6/7/10/11/13; die umbenannten `modul-03-spec`/`modul-04-adrs`);
   `modul-00`/`modul-09`/`grundlagen-durchsetzungsschicht` zuletzt.
2. **Finding-Schema** — je Treffer ein Eintrag: *{Quelle (Modul/§) · Regel-Delta ·
   betroffene d-check-Adaption/-Artefakt · Handlung → C oder D}*.
3. **Flotten-Stand** — u-boot / a-check / ai-harness-init auf `v5.0.0`? Bestimmt,
   ob a-checks Analyse noch überträgt (slice-083 §2.1-Kostensenker).
4. **Frischkontext-Review Pflicht** (slice-083 §4) — der Bump ist ein Re-Adopt.
   Ergebnis: die Finding-Liste, die C und D speist.

## 3. Findings (Modul-Delta)

_Wird beim Gegenlesen gefüllt (Schema §2.2). Bis dahin: der §2.3-Vorab-Befund aus
slice-083 ist der Ausgangspunkt, gegen die vendorte Quelle zu bestätigen._

## 4. Flotten-Stand

_Wird erhoben._

## 5. Definition of Done

- [ ] Alle 8 `grundlagen-*` + 17 Module gegen `v5.0.0` gegengelesen (Priorität §2.1).
- [ ] Finding-Liste im Schema vollständig; die slice-083-§2.3-Zugänge gegen die
  Quelle bestätigt/korrigiert und um Text-Deltas ergänzt; je Finding die Handlung
  (C oder D) zugeordnet.
- [ ] Flotten-Stand erhoben (überträgt a-checks Analyse noch?).
- [ ] `make gates` grün; unabhängiger Frischkontext-Review.

## 6. Risiken / offene Punkte

- **Urteilslast.** „Neue/geänderte Baseline-Regel trifft d-check" ist eine
  Ist-Abgleich-Frage; die Finding-Liste muss den d-check-Ist (Adaptionen,
  Artefakte, Config) korrekt spiegeln, sonst speist sie C/D falsch.
- **Abgrenzung zu C/D.** Etappe B **entscheidet** nichts und **ändert** nichts —
  sie sammelt nur. Jede „Handlung" ist ein Zeiger auf C oder D, keine Umsetzung.

## 7. Trigger

Abschluss von slice-084 (Etappe A): ab hier ist die v5.0.0-Quelle netzlos
vendored und belegbar lesbar.

## 8. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc/Prozess führt. Reine Analyse gegen die vendorte Baseline —
kein Brownfield-Spec-Bezug.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
