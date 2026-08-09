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
   + eigene ADR; **kein** neues Modul.

2. **Struktur-Semantik, Bindepunkt, Schwelle.** → **Entschieden 2026-08-09**
   (Nutzer-Entscheide, festgehalten in
   [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md)):

   - **Bindepunkt:** baseline-treu **getrennt** — die Closure-Prüfung hängt an der
     Closure, nicht am inneren Loop. Da die Konfiguration konventionell aus **einer**
     Datei kommt, verlangt das eine Herkunfts-Umschaltung: die **neue** Anforderung
     [`DC-FA-CLI-012`](../../../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)
     (`--config <datei>`) erlaubt zwei disjunkte Prüf-Profile im selben Repo.
   - **Struktur-Semantik:** `planning.closure.dir` (opt-in **innerhalb** des opt-in
     Moduls) · `heading-pattern` (RE2, Default deckt alle gemessenen Bestands-Formen)
     · `min-sentences` außerhalb Fenced-Code · `boilerplate` (literale Teilstrings,
     case-insensitiv).
   - **Grund-Codes:** **drei** statt einem — fehlender Abschnitt, zu dünn, Floskel
     (drei verschiedene Reparaturen; Vorbild ist die Zwei-Richtungs-Meldung des
     `targets`-Moduls).
   - **Floskel-Liste:** per Default **leer** — sie entscheidet über rot/grün und wäre
     sprach-gebunden falsch; d-check deklariert seine eigenen Phrasen in der
     Selbstkonfiguration.
   - **Schwelle:** `min-sentences` **4** (statt der Baseline-2). Belegt am eigenen
     Bestand: 92/92 abgeschlossene Slices tragen einen Abschnitt, das Minimum an
     Satzende-Zeichen außerhalb Code-Blöcken ist **5** ⇒ **kein Retrofit**; der
     Platzhalter eines offenen Slice zählt **1** ⇒ das Gate hat Zähne.
   - **Skill-Kopplung:** eigener `closure-note-reviewer.md` als Schwester des
     allgemeinen `reviewer.md` (Baseline-Default), der ausdrücklich **nicht**
     doppelt meldet, was die Struktur schon abdeckt.

   **Offen (Implementierungs-Entscheid, kein Vertragsdelta):** wie das
   `verify-closure-notes`-Target die zweite Profil-Datei benennt und wo sie liegt.

## 4. Definition of Done

- [x] Design-Entscheid (Abnahme-Punkte 1 + 2) festgehalten; CR + `DC-FA-*` + ADR angelegt
  (Lastenheft 0.50.0: [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  geschärft + [`DC-FA-CLI-012`](../../../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)
  neu; Spezifikation: C1–C5 in
  §[`DC-FA-PLAN-001.a`](../../../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  + §[`DC-FA-CLI-012.a`](../../../../spec/spezifikation.md#dc-fa-cli-012a--konfigurations-pfad---config)
  + §2-Schema;
  [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md) `Proposed`).
- [x] `make verify-closure-notes`-Gate implementiert (Struktur) + `closure-note-reviewer.md`-
  Skill (Inferenz); Tests fail-closed. Gate über ein **eigenes** Profil
  (`--config`, Closure-Bindepunkt in `fullbuild`, nicht in `gates`);
  Mutations-Gegenprobe am echten Bestand: alle drei Grund-Codes feuern, ohne
  Mutation 299 Dateien / 0 Befunde (kein Retrofit nötig).
- [x] Release (Tag + GHCR + Digest-Backfill). **v0.52.0** veröffentlicht
  (Pipeline-Lauf grün, Digest `sha256:412a6fd3…662c`); das gepinnte Release-Image
  gegen dieses Repo gegengeprüft.
- [x] `make gates` + `make verify-closure-notes` grün; unabhängiger Frischkontext-Review.
  **Zwei** Reviews mit getrennten Linsen (Vertrag + Code), beide blockierend —
  1 HIGH / 10 MEDIUM / 6 LOW / 3 INFO, alle eingearbeitet.

## 5. Risiken / offene Punkte

- **Produkt-Code + Release** — grösster Slice der welle-68; der volle Release-Flow ist Teil.
- **Struktur-vs-Inferenz-Grenze** sauber ziehen (das Gate prüft Struktur, der Skill
  Inhalt) — sonst Doppelbefund oder Lücke.
- **Dogfood:** d-check muss sein eigenes `verify-closure-notes` über den eigenen
  Bestand grün bekommen. **Vorab gemessen (2026-08-09), Risiko entschärft:** 92/92
  abgeschlossene Slices tragen einen Closure-Notiz-Abschnitt, das Minimum an
  Satzende-Zeichen außerhalb Code-Blöcken ist 5 — bei Schwelle 4 ist **keine**
  Nachbesserung bestehender Notizen nötig. Verbleibendes Restrisiko: die
  Heading-Form variiert (Nummer + Suffix), das Default-Muster muss alle Varianten
  decken; die Messung deckt genau das ab.
- **Zweite Profil-Datei = zweite Pflegestelle:** die Netzlos-Modullisten-Integrität
  muss beide Konfigurationen abdecken, sonst entsteht eine ungeprüfte zweite Tür.

## 6. Trigger

Freigabe des aus slice-090 herausgeschnittenen D-7 als welle-68-Slice.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** berührt mehrere Sub-Areas — Produkt-Code (`internal/`,
  `planning`-Modul), Spec (`spec/`) und Harness/Skills (`.harness/skills/`). **Alle unter
  dem Repo-Default GF** (`harness/conventions.md` §Modus: `*` = Greenfield). **Kein
  Brownfield:** der Code entsteht **spec-first**, wird nicht rückdokumentiert.
- **Offene Beobachtungen sichten:** `observations.md` = `— keine —`; nichts zu berücksichtigen.
  *Nachtrag 2026-08-09:* während dieses Slice ist **BEO-001** entstanden (Datei-Register
  ↔ Autoritäts-Tabelle driften unbemerkt, gefunden am fehlenden ADR-Index-Eintrag). Sie ist
  im [Register](../observations.md) eingetragen und **bewusst nicht** in diesen Scope
  gezogen — WIP-Limit, und slice-093 trägt bereits Produkt-Code + Release.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — „Spec führt, Code folgt: wir versprechen X, dann
liefern wir X". slice-093 ist ein **Produkt-Feature** (neues Gate: Go-Code + Spec +
Release), aber der **Modus** ist GF wie bei jedem d-check-Slice: der Change Request +
die Schärfung der `planning`-Anforderung + die ADR schreiben die **Zusage**, die
Go-Implementierung **liefert** sie. Der Unterschied zu den Doc-only-Migrations-Slices ist der **Scope**
(Produkt-Code + Release statt reiner Doc/Harness), **nicht** der Modus — Brownfield wäre
nur die Inventur bestehenden undokumentierten Codes, was hier nicht vorliegt.

## 9. Closure-Notiz (nach `done/`)

Umgesetzt und als **v0.52.0** veröffentlicht: das Modul `planning` trägt eine
zweite Fähigkeit (Closure-Note-Struktur, drei Grund-Codes, opt-in über
`closure.dir`), dazu die neue Option `--config` und der inferentielle
Schwester-Skill. Der Bestand blieb ohne Retrofit grün — die Vorab-Messung (92/92
Notizen, Minimum 5 Satzende-Zeichen) hat die Schwelle 4 getragen, statt sie zu
raten.

- **Lernsignal — eigene Lexik statt geteilter Parser war der teuerste Fehler.**
  Der Code-Review fand einen stillen Grün-Pfad **im Gate selbst**: die
  Abschnitts-Grenze zählte nackte `#`, weil ich sie neu geschrieben habe, statt
  den vorhandenen `parseATXHeading` zu nutzen. Eine Zeile wie `#1 war ein Thema`
  galt damit als H1 und schnitt die Notiz ab — eine Floskel dahinter blieb
  unsichtbar. Die Regel daraus ist allgemein: **wer Markdown-Lexik im zweiten
  Anlauf implementiert, baut eine Abweichung**, und in einem Gate ist jede
  Abweichung ein Kandidat für stilles Grün.
- **Architektur-Beobachtung — eine Config-Lokalität hat die CLI-Oberfläche
  erzwungen.** Der Modul-Schnitt (Fähigkeit *im* `planning`-Modul statt neues
  Modul) war richtig, hatte aber eine nicht offensichtliche Folge: weil die
  Konfiguration konventionell aus **einer** Datei kommt, läuft alles, was im
  Modul wohnt, dort mit, wo das Modul läuft — also in `gates`. Der gewünschte
  Closure-Bindepunkt war damit nur über eine **neue öffentliche
  CLI-Anforderung** erreichbar. Merke: die Frage „wo wohnt die Fähigkeit?" und
  die Frage „wann läuft sie?" sind über die Config-Lokalität gekoppelt, nicht
  unabhängig.
- **Fail-closed ist eine Zusage, die man zweimal treffen muss.** Drei der
  Review-Befunde waren derselbe Fehler an verschiedenen Stellen: ein leerer
  `--config`-Wert fiel still zurück, die Wurzel-Grenze war nur lexikalisch (ein
  Verzeichnis-Symlink entkam), und ein kandidatenfreies Closure-Verzeichnis
  meldete Erfolg. Jedes Mal war der *Vertrag* streng und die *Kante* offen.
- **Zwei Fixes waren von keinem Test gehalten** — ihr Rückbau blieb grün. Ein
  Test, der bei zurückgebauter Logik grün bleibt, belegt nichts; dasselbe galt
  für den Determinismus-Test, der nur die Sortierung des Test-Doubles maß.

**Review:** zwei unabhängige Frischkontext-Reviews mit getrennten Linsen
([Vertrag](../../../reviews/2026-08-09-slice-093-vertrag-review.md),
[Code](../../../reviews/2026-08-09-slice-093-code-review.md)), beide
**blockierend** — 1 HIGH / 10 MEDIUM / 6 LOW / 3 INFO, alle eingearbeitet. Eine
Reviewer-Aussage wurde begründet zurückgewiesen statt übernommen (nicht alle
Config-Meldungen tragen den Datei-Präfix; die ohne nennen gar keine Datei).
`make ci` grün, Digest `sha256:412a6fd3…662c`, das gepinnte Release-Image gegen
dieses Repo gegengeprüft.

**Anschluss:** [BEO-001](../observations.md) im Beobachtungs-Register — Datei-
Register driften unbemerkt gegen ihre Autoritäts-Tabelle (gefunden am fehlenden
ADR-Index-Eintrag während dieses Slice, bewusst nicht hier hineingezogen).
Vorschlag dort: ein kleines opt-in-Modul für die Richtung „Artefakt ⇒
registriert". Danach steht die welle-68-Closure an.

**Selbsttest:** diese Notiz ist die erste, die das in diesem Slice gebaute Gate
prüft.
