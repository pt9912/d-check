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

Drei Frischkontext-Leser (grundlagen · Prozess-Module 02/05/06/07 · Spec/Review
03-spec/04-adrs/10/11/13) haben die 25 v5.0.0-Dateien gegen den d-check-Ist
gegengelesen; jeder Treffer ist am Quelltext belegt (Reports als Rückgabe im
Slice-Bogen). **18 Findings**, dedupliziert und nach Ziel-Etappe geordnet; zwei
Funde **korrigieren die abgenommene Analyse** (§3.4). Die niedrig-priorisierten
Module (00/01/08/09/12/14/15/16) sind spot-gecheckt: keine neue Pflicht über die
Liste hinaus.

### 3.1 → Etappe C (Adaptions-/Konventionsspeicher-Bereinigung)

| # | v5.0.0-Quelle (belegt) | Regel-Delta | d-check betroffen |
|---|---|---|---|
| C-1 | `modul-03-spec` / `grundlagen-source-precedence`: „alle drei Straten obligatorisch" | Drei Spec-Straten sind **Default**, nicht mehr Adaption | [`MR-001`](../../../../harness/conventions.md#mr-001--source-precedence-mit-eigener-spezifikations-schicht) — Prämisse „Kurs-Default = zwei Ränge" entwertet → auf Provenienz abrüsten |
| C-2 | `grundlagen-referenz-richtung` §Spec-Decke | „kein Spec-Stratum verweist abwärts" ist Baseline-Default | [`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs) → Adaption zu Provenienz; der `matrix`-Mechanismus bleibt |
| C-3 | `grundlagen-referenz-richtung` §Prüfung: Spec-Check „ganzes Dokument, **ohne** ausgenommene Sektion" | die Historie-Ausnahme über Spec-Straten ist **verboten** (§3.3) | `.d-check.yml` `matrix.exclude-sections` (für `spec-straten` zu weit) + [`MR-012`](../../../../harness/conventions.md#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011)-Begründung („legitime Provenance") |
| C-4 | `grundlagen-referenz-richtung` 8×8-Matrix | d-checks Kodierung deckt Spec-Decke + ADR→Slice, ist aber **zu eng**: ADR→{Carveout,Welle,Roadmap}, Slice→Roadmap unbewacht | `.d-check.yml` `matrix` → Scope erweitern **oder** Grenze deklarieren |
| C-5 | `grundlagen-harness-dateien` §Konventionsspeicher | `conventions.md` = **Index + Datei je MR**, Zustand = Verzeichnis-Position | der ganze §Adaptions-Block (24 Inline-Einträge) → Datei-Split (Etappe-C-Kern) |
| C-6 | `grundlagen-harness-dateien` Pflichtgliederung | neue Pflichtfelder je Eintrag: **`Ersetzt-Baseline-Regel`**, `Löst auf`/`Ausgelöst durch Baseline-Stand`, `(schärft …)` im Titel | jeder `### MR-…`-Eintrag → Felder ergänzen (Status-Feld siehe §3.4) |
| C-7 | `grundlagen-durchsetzungsschicht` §Guard-Härtung: „Grenz-Zeile mitziehen" | jeder Wächter-MR muss nennen, was der Guard **nicht** kann (`python -c`, `env`, Wrapper) | [`MR-005`](../../../../harness/conventions.md#mr-005--härtung-ggü-b-cad-inhaltsbasierter-gate-nachweis-sub-shell-prüfung) → Grenz-Zeile ergänzen |
| C-8 | `grundlagen-harness-dateien` §Template-Schichtung (kein Producer-Ausnahme-Abschnitt) | die Baseline-Brücke von „keine Templates verkörpert" entfällt | [`MR-018`](../../../../harness/conventions.md#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates) + [`MR-020`](../../../../harness/conventions.md#mr-020--baseline-template-propagation-per-drift-audit-template-frei-bestätigt) → Entfall/Reklassifikation, co-located Vorlagen |

### 3.2 → Etappe D (Form-Konformität / Artefakte)

| # | v5.0.0-Quelle (belegt) | Regel-Delta | d-check betroffen |
|---|---|---|---|
| D-1 | `modul-06` §Roadmap-Struktur | Roadmap = **fünf** Abschnitte | `roadmap.md` führt drei (Meilensteine + Abgeschlossene Wellen fehlen; Closure-Bestand steckt in der `Vorgänger`-Prosa) |
| D-2 | `modul-06` §Wellen-Closure (5 Schritte) + `modul-05` | Welle-Plandatei + Closure-Notiz (`done/welle-NN-results.md`) + Wave-Self-Close-Commit | d-check führt Wellen nur in der Roadmap; `welle.template.md`/`welle-results.template.md` vendored ungenutzt |
| D-3 | `modul-06` §Beobachtungs-Register + `grundlagen-traceability` §Anker-Paarung | `observations.md` + `BEO-NNN` + Pflichtfeld `liegt in <Zielort>` in Welle-/Slice-§7 | fehlt komplett (`observations.template.md` vendored ungenutzt) |
| D-4 | `modul-05` §Zwei Schritte vor der Modus-Begründung | Slice-§ trägt „Vorgelagert — Sub-Area prüfen / offene Beobachtungen sichten" | d-check-Slice-§6 ohne die zwei Vorprüfungen (hängt an D-3) |
| D-5 | `modul-05` §Lifecycle: „Zustand = Verzeichnis, **kein** Status-Feld" | Slice- **und** MR-`Status:`-Feld ist eine zweite Zustandsquelle | Slice-Header-Konvention ([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise) / [`MR-014`](../../../../harness/conventions.md#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template)) + jeder MR-Eintrag |
| D-6 | `modul-10` §Report-Gerüst + `review-report.template` | Kopffelder **Review-Art / Skill-Version / Modell-ID** + Finding-`klasse` + Finding-Klassen-Zeile | `.harness/skills/reviewer.md` §Output-Schema/§Ablage; Reports uneinheitlich (erst ab 2026-08-01 geführt) |
| D-7 | `modul-10`/`modul-11` §closure-note-reviewer | Schwester-Skill **über** einem Closure-Note-Struktur-Gate | `.harness/skills/` (nur `reviewer.md`); `AGENTS.md` §4 kennt kein `verify-closure-notes` → an ein Gate gekoppelt |
| D-8 | `modul-04` §Re-Evaluierungs-Trigger | jede ADR trägt Trigger (oder „permanent"), im Welle-Closure-Trigger-Audit geprüft | d-check-ADR-Kopf ohne Trigger-Feld (kein Template-Zwang → gering) |
| D-9 | `modul-07` §Carveout-Audit-Slice | pro Welle-Closure ein Audit-Slice **vor** `done/` | fehlt (an D-2 gekoppelt; heute latent — 0 aktive Carveouts) |
| D-10 | `grundlagen-referenz-richtung` (Baseline-Kopf-Currency) | der Skill-Kopf zitiert noch die retirete `grundlagen-konventionen.md` + „Kurs-Welle 18 §8" | `.harness/skills/reviewer.md` Z. 4–6 (Z. 39 ist bereits korrekt — Kopf hinkt nach; **vom Auftraggeber vorab bemerkt**) |

### 3.3 Der schärfste Fund — Historie-Provenance revoziert (C-3, DOPPELT belegt)

Reader A (grundlagen, F-8) und Reader C (Spec-Module, B-5) fanden **unabhängig**:
`modul-03-spec` / `grundlagen-referenz-richtung` verbieten ADR-/Slice-Verweise in
den Spec-Straten **in jedem Abschnitt, auch der Historie** — der Spec-Check läuft
whole-document ohne Ausnahme-Sektion („gäbe es dort eine ausgenommene Sektion,
wäre sie genau die Stelle, an der die Verweise landen"). d-check nimmt genau die
Historie aus (`.d-check.yml` `matrix.exclude-sections`, begründet in C-3 als
„legitime Provenance"). **Real nicht-konform:** `spec/spezifikation.md` §7
(ADR-Links + `slice-NNN`-Spalte) und `spec/lastenheft.md` §7 (`slice-NNN`-Spalte).
Klassischer stiller Regel-Wechsel (slice-083 §4-Sinn). Handlung: **C** (die
Ausnahme für die `spec-straten`-Klasse zurücknehmen, die Planungs-Klassen behalten
sie) **+ D** (die Verweise aus den Spec-§7 entfernen — **berührt die kanonische
Spec, ausdrücklicher Abnahme-Punkt für Etappe C/D**).

### 3.4 Korrekturen an der abgenommenen Analyse (slice-083 §2.3/§2.4)

- **Fork-Klassifikation (Reader A F-7/F-10):** §2.4 markiert u. a.
  [`MR-004`](../../../../harness/conventions.md#mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild)
  und C-7 als baseline-ungeregelte **Forks**. Das ist unter v5.0.0 falsch:
  `grundlagen-durchsetzungsschicht` **kanonisiert** genau die Guard-Härtung
  (Landung als MR, Grenz-Zeile) und den inhaltsbasierten Gate-Nachweis
  (Content-Hash = Design-Eigenschaft 2) — beide können `Ersetzt-Baseline-Regel`
  auf diese Datei benennen = **Adaptionen, keine Forks**. Zudem ist die
  Quell-Fork-Definition enger (**pauschale** Nichtanwendbarkeit der Baseline) —
  **keine** d-check-Adaption erfüllt das. Etappe C muss das Fork-Kriterium am
  Quelltext schärfen **bevor** es reklassifiziert; die §2.4-Fork-Liste schrumpft
  (Rest-Kandidaten allenfalls
  [`MR-015`](../../../../harness/conventions.md#mr-015--auflösung-der-mr-012-pointer-drift-agentsmd-routet-spiegelt-nicht-mehr) /
  [`MR-021`](../../../../harness/conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden),
  und selbst die finden einen Anker in `grundlagen-harness-dateien`).
- **`Status: Accepted`-Feld (Reader A F-4):** slice-083 §2.3 nennt es als neues
  Pflichtfeld. Die **grundlagen** stützen das **nicht** („Zustand = Verzeichnis,
  kein Status-Feld"). Es kann nur aus dem MR-**Template** stammen — in Etappe C
  aus der vendorten Vorlage zu belegen, nicht als grundlagen-Pflicht zu behandeln.

### 3.5 Bestätigt konform / kein Delta (Auswahl der Negativbefunde)

- **Change-Request-Fußabdruck** (`modul-03`): d-checks §Anforderungs-Anlege-Prozess
  führt CR bereits als reinen Doc-Vorgang (kein Schema/Datei/Gate) — konform.
- **ADR-Immutabilität / MADR-Form** (`modul-04`): via `vcs`-Modul / `adr-check`
  konform (nur der Re-Eval-Trigger offen, D-8).
- **Traceability-Kern** „Commit trägt ≥1 ID" (`grundlagen-traceability`): durch das
  `commits`-Modul deckungsgleich (nur die Herkunfts-Anker-Schicht offen, D-3).
- **Halluzinierte-Gates-Hard-Rule** (`modul-13`): d-check namentlich als konform
  genannt; **`modul-15` §Doku-Konsistenz-Drift** beschreibt genau, was d-checks
  `targets`-Modul tut — d-check **realisiert** das Baseline-Konzept (INFO, keine
  Handlung).
- Lifecycle-Mechanik, Modus-pro-Sub-Area, Source-Precedence-Rangliste,
  Reviewer-vs-Verifier-Trennung: geprüft, konform.

## 4. Flotten-Stand

u-boot + a-check stehen auf `v3.5.2` (`harness/conventions.md` §Baseline),
ai-harness-init ist nicht lokal prüfbar. d-check ist mit `v5.0.0` **vor** der
Flotte. Folge für den slice-083-§2.1-Kostensenker: a-checks v3.5.2-Analyse trägt
nur die `v1.4.0`→`v3.5.2`-Hälfte; die **v4.0.0/v5.0.0-Deltas** (3-Straten-Default,
8×8-Matrix, `conventions`-Index-Form, neue MR-Pflichtfelder, Historie-Provenance-
Revocation, `observations.md`/Wellen-Closure) sind **d-check-first** — kein
übertragbarer Flotten-Präzedenzfall.

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
