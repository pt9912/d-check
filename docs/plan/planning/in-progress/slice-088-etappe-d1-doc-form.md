# Slice slice-088: Regelwerk-Migration Etappe D-1 — Doc-Form-Konformität (Roadmap/AGENTS/ADR-Trigger)

**Status:** In Arbeit (welle-67).

**Welle:** welle-67-baseline-v500-migration (Etappe D, erster von vier
„Mini-Welle"-Slices, nach
[slice-087](../done/slice-087-spec-historie-referenzrichtung.md)).

**Bezug:** Erster Slice der **Etappe D (Form-Konformität)** aus
[slice-083](../done/slice-083-regelwerk-v500-migration-analyse.md) §2.7 /
[slice-085](../done/slice-085-etappe-b-modul-delta.md) §3.2 (Findings D-1…D-11).
Etappe D ist als **Mini-Welle** in vier thematische Slices geschnitten
(Nutzer-Entscheid 2026-08-02): **slice-088 Doc-Form** (D-1/D-8/D-11), slice-089
Review-Infrastruktur (D-6/D-7/D-10), slice-090 Wellen-Lifecycle +
Beobachtungs-Register (D-2/D-3/D-4/D-9), slice-091 Slice-`Status:`-Feld (D-5). Dieser
Slice trägt die rein **mechanische Doc-Form** — plus die eine Adaptions-Entscheidung
D-1 (Ruhe-Marker ↔ Template). **Kein Release.**

**Autor:** pt9912. **Datum:** 2026-08-02.

---

## 1. Ziel

Drei bestehende Doc-Artefakte an die v5.0.0-Baseline-Form angleichen: die **Roadmap**
(D-1: drei fehlende Abschnitte), **AGENTS.md** (D-11: §1-Drift + stale Cache-Prosa) und
die **ADR-Re-Evaluierungs-Trigger-Konvention** (D-8: konsequent). **Keine** neuen
Prozesse/Artefakte (die kommen in slice-090), kein Release.

## 2. Vorgehen

1. **D-1 Roadmap-Abschnitte.** `roadmap.md` um die drei fehlenden Baseline-Abschnitte
   ergänzen: `## Meilensteine`, `## Abhängigkeitsgraph` (mermaid) und
   `## Abgeschlossene Wellen` (Tabelle); den Closure-Bestand aus der
   `## Aktuelle Welle`-„Vorgänger"-Prosa in `## Abgeschlossene Wellen` überführen.
   Reihenfolge/Regeln je `modul-06-roadmap.md`. **Abnahme-Punkt (§3):** die
   `## Aktuelle Welle`-Form (Ruhe-Marker vs. Template-Felder).
2. **D-11 AGENTS.md.** §1 gegen `AGENTS.template.md` angleichen: Modul-9-Kanon-Zeiger,
   `{regelwerk,templates}`-Layout, und die retirete Templates-Cache-Prosa
   (lab-templates.zip / .harness/cache) entfernen/umschreiben — das in slice-084 als
   Deferral vermerkte Nachziehen einlösen.
3. **D-8 ADR-Re-Evaluierungs-Trigger.** Die Konvention „jede ADR trägt die Sektion
   `## Re-Evaluierungs-Trigger` (oder „permanent")" bestätigen/dokumentieren
   (Haus-Stil, bereits 27/46; ältere ohne Sektion sind immutable/grandfathered). Der
   **Welle-Closure-Trigger-Audit** koppelt an D-9 → **slice-090** (hier nur die
   Konvention, kein Audit-Mechanismus).
4. **Gate.** `make gates` (inkl. planning-check über die neue Roadmap-Struktur) +
   `make adr-check` grün; unabhängiger Frischkontext-Review.

## 3. Abnahme-Punkte

1. **`## Aktuelle Welle`-Form (D-1) — Ruhe-Marker ↔ Template.** Das Baseline-Template
   führt Struktur-Felder (**Welle-ID / Start / Geplantes Ende / Closure-Trigger**) und
   nennt **immer** eine laufende Welle. d-checks `planning`-Modul
   ([`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in))
   **erzwingt** dagegen den „Keine aktive Welle"-Ruhe-Marker in genau diesem Abschnitt
   (Roadmap-↔-in-progress-Invariante). **Entscheid:** (a) den Ruhe-Marker als
   **deklarierte Adaption** (neuer Eintrag im Konventionsspeicher) behalten — die
   Mechanisierung ist d-check-Identität — und die Aktive-Welle-Prosa optional um die
   Template-Felder ergänzen; oder (b) auf die reine Template-Feld-Form heben und das
   `planning`-Modul anpassen (größer, berührt Produkt-Code + Spec).

## 4. Definition of Done

- [ ] `roadmap.md` führt die sechs Baseline-`##`-Abschnitte; der Closure-Bestand steht
  in `## Abgeschlossene Wellen` (nicht mehr in der „Vorgänger"-Prosa).
- [ ] `AGENTS.md` §1 template-konform (Modul-9-Zeiger, `{regelwerk,templates}`-Layout,
  keine stale Templates-Cache-Prosa).
- [ ] ADR-Re-Evaluierungs-Trigger-Konvention dokumentiert (Audit → slice-090).
- [ ] `## Aktuelle Welle`-Form-Entscheid (Abnahme-Punkt 1) umgesetzt.
- [ ] `make gates` + `make adr-check` grün; unabhängiger Frischkontext-Review.

## 5. Risiken / offene Punkte

- **planning-check-Kopplung:** die Roadmap-Umbauten dürfen den Ruhe-Marker- bzw.
  Heading-Guard nicht brechen (das `planning`-Modul, `## Aktuelle Welle`-Heading exakt).
- **AGENTS.md ist Source-Precedence-Anker** — Änderungen an §1 vorsichtig, review-pflichtig.
- **Abhängigkeitsgraph (mermaid):** die Wellen-Historie ist weitgehend linear — der
  Graph ist eher formal; Nutzen vs. Pflege im Review prüfen.

## 6. Trigger

Abschluss von [slice-087](../done/slice-087-spec-historie-referenzrichtung.md)
(C-3-Nachzug); Etappe D als Mini-Welle geschnitten (Nutzer-Entscheid).

## 7. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc/Prozess führt. Berührt die *Harness/Prozess*-Doku (Roadmap,
`AGENTS.md`, ADR-Konvention); greenfield-Form-Angleich an die adoptierte Baseline,
ohne Brownfield-Spec.

## 8. Closure-Notiz (nach `done/`)

_Ausstehend._
