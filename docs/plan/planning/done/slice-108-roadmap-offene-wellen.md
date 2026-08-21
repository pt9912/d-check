# Slice slice-108: Roadmap auf §Offene Wellen — Etappe C-1 des Baseline-Bumps

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** welle-78-baseline-v560-migration.

**Bezug:** Audit-Finding C-1 aus
[slice-107](../done/slice-107-baseline-v560-delta-audit.md) §9
(Stufe v5.5.0); betrifft
[`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform)
und die `planning`-Selbstkonfiguration
([`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
— nur Config, kein Produkt-Code).

**Autor:** pt9912. **Verantwortlich:** pt9912. **Datum:** 2026-08-21.

---

## 1. Ziel

Die Roadmap folgt der v5.6.0-Baseline-Form: **§Offene Wellen** statt
§Aktuelle Welle — derivativ (der Zustand sind die flachen Welle-Dateien; woran
gearbeitet wird, sagt das `Welle:`-Feld der Slices in `in-progress/`),
Ruhe-Marker **„Nichts in Arbeit"** bei leerem `in-progress/`, Wellen-Closure
Schritt 5 ohne Beförderung. Baseline-Default sticht die repo-lokale Adaption
(Auftraggeber-Linie seit der v5.0.0-Migration).

## 2. Vorgehen

1. Roadmap: §Aktuelle Welle → §Offene Wellen (Zeiger auf das flache
   Wellendokument statt Struktur-Felder; Marker-Text neu), Sektions-Regeln-
   Zitate auf die v5.6.0-Formulierungen; Drift-Log-Eintrag.
   **Spiegel-Liste der alten Form** (Review-Befund F-3 zum Audit — vor dem
   Editieren per grep verifizieren, nicht abschreiben): `AGENTS.md` §3.3
   (Lifecycle-Move-Ausnahme nennt „§Aktuelle Welle"/„Keine aktive Welle")
   und §4-planning-check-Zeile, `harness/README.md` §Sensors-planning-check-
   Zeile, der `Makefile`-Kommentar am planning-check-Target,
   [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)
   (Wortlaut der gebündelten Verweise) — plus alles, was der grep noch findet.
2. `.d-check.yml` (+ `.d-check.closure.yml` prüfen): `planning.heading` auf
   `## Offene Wellen`, `marker` auf `Nichts in Arbeit` — das Produkt deckt
   beides per Config; **kein** Produkt-Code. Vor dem Umstellen messen (ein
   Lauf je Profil).
3. [`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform)
   entscheiden: auflösen (Baseline-Stand v5.6.0 trägt Marker-Semantik selbst)
   oder fortschreiben — die Datei-Historie dokumentiert den Übergang.
4. Konsumenten-Blick: die Drift-Kopplung der dritten `planning`-Fähigkeit
   (Anspruch ⇔ genau ein flaches Wellendokument) bleibt im
   Ein-Wellen-Betrieb wahr; die Grenze für künftigen Mehr-Wellen-Betrieb wird
   im Config-Kommentar benannt, nicht gelöst.
5. **Benannte Grenze — Produkt-Default bleibt:** der `planning.heading`-Default
   des Produkts ist weiterhin `## Aktuelle Welle`; d-check lebt nach C-1
   dauerhaft auf Nicht-Default-Config. Eine Default-Änderung wäre ein
   Konsumenten-Breaking-Change und ist bewusst **nicht** Teil dieses Slice —
   sie braucht einen eigenen Entscheid (CR/ADR), falls die Baseline-Form sich
   bei den Konsumenten durchsetzt (Review-Befund F-10).

## 3. Definition of Done

- [x] Roadmap in v5.6.0-Form, `make planning-check`/`make gates` grün (beide
      Profile), Marker-Wächter belegt (Negativ-Probe: Marker bei belegtem
      `in-progress/` ⇒ rot).
- [x] [`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform)-Entscheid vollzogen und im Index nachvollziehbar.
- [x] Kein Produkt-Code berührt; unabhängiger Review.

## 4. Trigger

**Start** (`open` → `in-progress`):
[slice-107](../done/slice-107-baseline-v560-delta-audit.md) in `done/`
**und** WIP-Slot frei.

## 5. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Planning-/Harness-Doku + Selbstkonfiguration, GF.
- **Beobachtungen sichten** (Stand 2026-08-21: BEO-006 offen, 1×): BEO-006
  einschlägig als Arbeitsregel — die Roadmap-Umstellung erzeugt Moves; vor
  pfad-selektiven Commits den Index prüfen.

## 6. Sub-Area-Modus-Begründung

**GF (Repo-Default)** — Form-Angleichung an die adoptierte Baseline.

## 7. Closure-Notiz (nach `done/`)

**Geliefert:** die Roadmap läuft in der v5.6.0-Form §Offene Wellen (derivativ,
Ruhe-Marker „Nichts in Arbeit" mit Wächter), beide Prüf-Profile per
`planning.heading`/`marker` umgestellt,
[`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform)
aufgelöst — sein Gegenstand ist Baseline-Default geworden. Die Negativ-Probe
ist doppelt belegt (eigener Lauf und Reviewer-Reproduktion mit
SHA-verifiziertem Restore).

**Zwei Dinge trugen über die Config hinaus.** Erstens: der Marker wird als
**Substring der Sektion** gematcht — der Sektions-Regeln-Text darf ihn also
nie literal zitieren; das Baseline-Template selbst würde als Kommentar-Zitat
matchen. Zweitens fand der Review, dass der `planning-drift`-Klartext die
Sektion **verdrahtet** nannte statt aus der Config — d-check ist der erste
Nicht-Default-Konsument seines eigenen Moduls und hätte sich selbst eine
Sektion gemeldet, die es nicht mehr gibt. Ein Klartext, der Konfigurierbares
benennt, gehört an die Config gebunden — dieselbe Bindungs-Lehre wie bei den
Modul-Enumerationen.

**Benannte Grenze statt gelöstem Problem:** der baseline-legitime Zustand
„Welle offen, `in-progress/` leer" bleibt `wave-drift`-rot, und der
Produkt-Default (`## Aktuelle Welle`) bleibt unverändert — beides sind eigene
Entscheide (Konsumenten-Fläche), keine Nebeneffekte dieses Slice. Aus der
Definition of Done bleibt nichts offen; der eine berührte Produkt-String ist
ein Befunds-Klartext (nicht stabilitätszugesagt), im Review-Auflagen-Commit
deklariert.
