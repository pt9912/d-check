# Slice slice-111: `planning.waves.mode` — Kennungs-Bijektion als opt-in (`one`|`many`)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** welle-79-zwei-haelften-ein-waechter (zugeordnet bei der Eröffnung).

**Bezug:**
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
§Wellen-Invariante (Aussage 1/2, `wave-drift`),
[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
(Proposed — wird per `## Geschichte` fortgeschrieben, nicht ersetzt),
[`MR-025`](../../../../harness/conventions.md#mr-025) (Semantik-Fläche ⇒
Spiegel-Liste vor dem Editieren); **formaler CR des Konsumenten
ai-harness-course vom 2026-08-21** („planning.waves: Bijektion statt
Singleton" — Wortlaut beim Auftraggeber; die Landung ist per CR-Definition
der Lastenheft-Commit, keine eigene Datei); Baseline-Regelwerk
`modul-06-roadmap.md` §Offene Wellen in der v5.7.0-Fassung (vendored durch
slice-110 — bindende Reihenfolge).

**Verantwortlich:** pt9912. **Datum:** 2026-08-21.

---

## 1. Ziel

Die erste Aussage der Wellen-Invariante bekommt einen **opt-in
Kardinalitäts-Modus** `planning.waves.mode: one | many` (Default `one`,
Befundsatz ohne den Schlüssel byte-identisch —
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
Unter `many` vergleicht `wave-drift` **Kennungs-Mengen statt Marker gegen
Anzahl**: die im `planning.heading`-Block genannte Kennungs-Menge gegen die
Menge der flachen Wellendokumente, beide Richtungen, jede Kardinalität
einschließlich null; der Ruhe-Marker geht **nicht** ein — er ist Gegenstand
von `planning-drift` (die Marker-Hälfte hält per team-sim-Messung
s04c/s04d des Konsumenten in beiden Richtungen). Das ist genau das „eigene
Prädikat" der Listen-Hälfte aus der Baseline v5.7.0, und es macht die drei
unter `one` roten, baseline-legitimen Zustände unter `many` grün
(Eröffnungs-Fenster · wellenloser Slice · Mehr-Wellen-Betrieb, gemessen
s04a/s04b). Konsumentensichtbar additiv ⇒ **Release Minor v0.62.0**;
anschließend stellt dieses Repo (es fährt das Offene-Wellen-Modell selbst)
seine Prüf-Profile auf `many` um.

## 2. Vorgehen

1. **CR-Commit (nur das Lastenheft, vor jedem weiteren Edit):** Version
   0.61.1 → 0.62.0, §Wellen-Invariante — Zeile 1/2 der Aussage-Tabelle um
   beide Modi, neuer Prosa-Absatz „Zwei Kardinalitäts-Modelle, ein Prädikat
   (`planning.waves.mode`)" nach „Vier Reparaturen, vier Grund-Codes", fünf
   neue Akzeptanzkriterien (Happy-Path `many` · Negative beide Richtungen
   mit Kennung als Befund-Ziel · Marker-Orthogonalität samt
   `one`-Abgrenzung · Modus-Default byte-identisch · fail-closed
   unbekannter/explizit leerer Modus, Exit 2), §7-Historie-Zeile in der
   Langform (Anlass, Messung, Nicht-Ziele, Schnitt-Begründung). Wortlaut
   aus CR §4, redaktionell eingepasst.
2. **[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
   fortschreiben** (`## Geschichte`, Status bleibt Proposed —
   kein Immutabilitäts-Konflikt): der Mode-Entscheid samt Begründung, die
   Kennungs-Erkennung (dasselbe `waveID`-Verfahren wie in den Registern —
   literales Glob-Präfix + Ziffernfolge, zeilenweise über die
   **Prosa-Zeilen** des Blocks, Fence-Inhalte zählen nicht, Mehrfachnennung
   zählt einmal, layout-agnostisch für Tabellen- wie Listen-Form), der
   Grund-Code-Entscheid (`wave-drift` bleibt — die Reparatur ist dieselbe;
   beide Richtungen über das **`target`** unterscheidbar: Kennung statt
   `waves.dir`, wie es die Register-Aussagen praktizieren; Dedup-Probe über
   das Befund-Tupel), fail-closed-Rand (unbekannter Modus inkl. explizit
   leerem String ⇒ Exit 2, Zeiger-Disziplin wie die übrigen
   `waves`-Schlüssel).
3. **[`MR-025`](../../../../harness/conventions.md#mr-025)-Spiegel-Liste
   vor dem Editieren** — der CR benennt die Flächen: Lastenheft
   §Wellen-Invariante (Schritt 1), Spezifikation §Algorithmus W3 +
   §2-Config-Schema + §4-Grund-Code-Zeile, `planning_waves.go` (Prädikat),
   `planning.go` (`planningActiveStatus` liefert neben dem Bool die
   Kennungs-Liste des Blocks), Config-Rand (`configyaml` `rawWaves`/
   `applyWaves`, `model`-Config samt Effective-Zugriff), Benutzerhandbuch
   (§5-Config-Referenz + planning-Abschnitt), `print-config`-Template,
   CHANGELOG, Grenz-Kommentare der eigenen Prüf-Profile, AGENTS §3.3.
4. **Code + Tests:** Modus-Parsing fail-closed; Bijektion beide Richtungen
   mit Kennung als `target`; die vier team-sim-Zustände als eigene
   Testfälle nachgebaut (s04b wird unter `many` grün und bleibt unter `one`
   rot — die Abgrenzung ist ein Boundary-Test, kein Fehler), Fence-Probe,
   Mehrfachnennung, Kardinalität null (leerer Block + kein Dokument grün,
   Marker-Anwesenheit ohne Einfluss), Default-Probe byte-identisch.
5. **Eigene Profile + Grenz-Kommentare + AGENTS:** nach Release und
   Digest-Backfill (der gepinnte Prüfer muss den Schlüssel kennen — vorher
   wäre `mode:` Exit 2) `.d-check.yml`/`.d-check.closure.yml` auf
   `mode: many`; die Grenz-Kommentare („Welle offen, `in-progress/` leer
   ist `wave-drift`-rot") und den AGENTS-§3.3-Atomicity-Hinweis auf die
   neue Lage entspannen — die Eröffnungs-Atomicity ist dann Konvention,
   nicht Gate-Zwang. Positiv-/Negativ-Probe am eigenen Bestand.
6. **Release-Prep + Release:** Benutzerhandbuch-Versionsschritt
   (planning-Abschnitt + Config-Tabelle um `waves.mode`; **inkl. des
   vorgemerkten §6-Schliffs**: planning-Zeile auf „im `heading`-Block
   (Default `## Aktuelle Welle`)"), CHANGELOG, Handbuch-§11-Zeile
   chronologisch **unter** die letzte, Release **v0.62.0** → CI → Tag →
   GHCR → Digest-Backfill.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Änderung an `planning-drift`** (die Marker-Hälfte ist korrekt und
  gemessen), **kein zweiter Grund-Code**, **keine Default-Änderung** —
  CR §6 wörtlich.
- **W4/W5 unberührt** (Vorschau- und Register-Aussagen), ebenso Slice- und
  Closure-Fähigkeit.
- **Keine Festlegung**, ob der Aktiv-Block seine Kennungen als Tabelle oder
  Liste führt — die Erkennung ist layout-agnostisch.
- **Kein Mehr-Wellen-Betrieb dieses Repos** — `many` macht ihn prüfbar; ob
  d-check je mehr als eine Welle öffnet, ist ein eigener Roadmap-Entscheid.

## 4. Definition of Done

- [ ] CR-Commit (Lastenheft allein: 0.62.0 + §4-Text + Historie-Zeile)
      liegt **vor** der Implementierung in der Historie.
- [ ] [ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
      per `## Geschichte` fortgeschrieben (Proposed, kein neues
      ADR-Dokument).
- [ ] [`MR-025`](../../../../harness/conventions.md#mr-025)-Spiegel-Liste
      vor dem ersten Edit dokumentiert und am Ende abgehakt.
- [ ] Default-Probe: ohne `mode`-Schlüssel Befundsatz byte-identisch
      ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus));
      fail-closed-Probe: unbekannter und explizit leerer Modus ⇒ Exit 2 mit
      Schlüssel-Nennung.
- [ ] Testfälle: beide Bijektions-Richtungen rot mit Kennung als `target`,
      die vier team-sim-Zustände nachgebaut, Fence-/Mehrfach-/Null-Proben.
- [ ] `make fullbuild` grün (Exit explizit); Release v0.62.0 auf GHCR,
      Digest-Backfill committet; **danach** eigene Profile auf `many`,
      Gate grün am eigenen Bestand.
- [ ] Unabhängiger Review vor der Closure.

## 5. Abnahme-Punkte / Risiken

- **Die Sektions-Prosa wird Teil der Messfläche:** unter `many` zählt jede
  Wellen-Kennung in den Prosa-Zeilen des Blocks als Zeiger — die
  Paraphrase-Disziplin des Blocks erweitert sich vom Marker-Wortlaut auf
  Wellen-Kennungen (die Sektionsregel darf keine `welle-<n>` nennen). Vor
  der Profil-Umstellung den eigenen Block prüfen.
- **`target`-Wechsel nur unter `many`:** der `one`-Pfad behält Klartexte
  und `target` unverändert — sonst wäre der Default nicht byte-identisch.
  Die Default-Probe ist der Wächter dieses Versprechens.
- **Zwei Wahrheiten während der Umbau-Phase:** zwischen Release und
  Digest-Backfill prüft das gepinnte Alt-Image weiter `one` — die eigene
  Profil-Umstellung ist deshalb ausdrücklich **nach** dem Backfill
  sequenziert (Schritt 5), nicht mit dem Code-Commit.

## 6. Trigger

**Start** (`open` → `in-progress`): slice-110 in `done/` (bindend — die
ADR-Fortschreibung zitiert die **vendorte** v5.7.0-Formulierung, nicht das
Kurs-Repo).

**Rückführungen:** `in-progress` → `next`, falls die Umsetzung zeigt, dass
die Kennungs-Erkennung eine über `mode` hinausgehende Config-Fläche braucht
(dann erst CR-Rücksprache mit dem Konsumenten, dann Umbau).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Kern (`internal/hexagon`, GF —
  Repo-Default) samt Config-Rand (`configyaml`), dazu Spec- und
  Harness-Doku (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-21:
  BEO-006/BEO-007/BEO-008 offen bei je 2, alle mit gelebtem
  Gegenmittel): die BEO-006/007-Arbeitsregeln gelten in diesem Slice
  ausdrücklich (`git status` vor pfad-selektiven Commits; Gate-Exits
  explizit statt Pipe); BEO-002 wirkt verkörpert als
  [`MR-025`](../../../../harness/conventions.md#mr-025) — der Slice
  ändert eine zugesagte Semantik-Fläche, die Spiegel-Liste steht in §2
  Schritt 3.

Slice-ID: slice-111. Betroffene IDs:
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus);
[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
(Fortschreibung). Module: `planning` (Kern `rules/`), Config-Rand
(`configyaml`, `model`), Spec (Lastenheft + Spezifikation), Handbuch,
eigene Prüf-Profile, AGENTS. Gates: `make gates`, `make fullbuild`, GUARD;
Release-Pipeline.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — additive Erweiterung einer eigenen,
spezifizierten Fähigkeit auf formalen Konsumenten-CR; kein Legacy-Import.

## 9. Closure-Notiz (nach `done/`)

*Wird bei der Closure geschrieben (Struktur nach `closure.heading-pattern`).*
