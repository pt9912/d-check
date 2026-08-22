# Slice slice-083: Regelwerk-Migration v1.4.0 → v5.0.0 — Delta-Analyse und Etappen-Schnitt

**Status:** done

**Welle:** keiner Welle zugeordnet — die [Roadmap](../in-progress/roadmap.md)
steht in Ruhe (welle-66 abgeschlossen); die Einplanung ist Teil der Abnahme (§5).

**Bezug:** Baseline-Deklaration
[`harness/conventions.md` §Baseline](../../../../harness/conventions.md#baseline)
und [§Adoptierte Konventions-Quellen](../../../../harness/conventions.md#adoptierte-konventions-quellen),
Briefing [`AGENTS.md`](../../../../AGENTS.md) §1, Harness-Einstieg
[`harness/README.md`](../../../../harness/README.md), Doku-Gate
[`.d-check.yml`](../../../../.d-check.yml), Materialisierungs-Skript
[`tools/harness/fetch-baseline-cache.sh`](../../../../tools/harness/fetch-baseline-cache.sh).
Kein Vertrag der Produkt-Achse berührt: **kein** Lastenheft-/Spezifikations-Bump,
**keine** neue Anforderungs-Kennung, **keine** ADR, **kein** Release — dies ist
eine Harness-/Konventions-Änderung. **Kein Change Request** (der Vertrag bleibt
unangetastet).

**Autor:** pt9912. **Datum:** 2026-07-25.

> **Hinweis.** In dieser Datei werden **keine** neuen Kennungen vergeben und
> **keine** Artefakte geändert. Sie misst den Ist-Zustand und schlägt einen
> Schnitt vor; jede Umsetzung ist ein eigener Slice.

---

## 1. Ziel

Den Sprung der adoptierten Baseline von `v1.4.0` auf `v5.0.0` **vollständig**
vorbereiten — mit der Lese-Leistung, die einen Re-Adopt von einem Auto-Update
unterscheidet: Was hat sich geändert, und welche Adaption dieses Repos wird davon
entwertet?

> **Neubasierung 2026-08-01.** Diese Analyse zielte ursprünglich auf `v3.5.2`
> (Stand 2026-07-25). In der Woche danach kamen **zwei weitere Majors** dazu —
> `v4.0.0` (Asset-Umbenennung: Module 03/04, drei Straten Pflicht, Matrix 8×8,
> Adaptions-Block verlässt den Lesepfad) und `v5.0.0` (`grundlagen-konventionen.md`
> **entfällt ersatzlos**, grundlagen 3 → 8 aufgesplittet). Ziel ist jetzt `v5.0.0`;
> die v3.5.2-Befunde bleiben als **Untergrenze** gültig, die v4/v5-Brüche liegen
> darüber. Wo unten „gegen `v5.0.0`" steht, ist der geprüfte Zielstand gemeint.

| | |
|---|---|
| d-check gepinnt auf | **`v1.4.0`** (Adoption 2026-06-10; auf den Release-Tag gepinnt mit [`MR-011`](../../../../harness/conventions.md#mr-011--baseline-auf-release-tag-gepinnt), gehoben mit [`MR-012`](../../../../harness/conventions.md#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011) und [`MR-016`](../../../../harness/conventions.md#mr-016--baseline-pin-hebung-zweiter-nachtrag-zu-mr-011)) |
| aktuelles Kurs-Release | **`v5.0.0`** (Kurs-Welle 64, 2026-08-01) |
| dazwischen | **≥ 14 Releases**, **vier** Major-Sprünge (v1 → v2 → v3 → v4 → v5); die zwei jüngsten (`v4.0.0` Welle 62, `v5.0.0` Welle 64) liegen über der ursprünglichen v3.5.2-Analyse |
| Flotte | Erhebung 2026-07-25 (ggü. `v5.0.0` nachzuprüfen — Etappe B): u-boot/a-check waren auf `v3.5.2` (a-check Etappe A durch), ai-harness-init `v3.5.1`; ob die Flotte inzwischen v4/v5 nachzog, ist offen. **d-check bleibt Schlusslicht** (gepinnt `v1.4.0`) |

**Auftrag (Maintainer-Vorgabe 2026-07-25, neubasiert 2026-08-01):** vollständige
Migration nach `v5.0.0` — **auch dort, wo ein `MR-*` heute etwas anderes sagt**.
Wo eine Adaption nur deshalb existiert, weil der alte Default anders war,
**gewinnt der `v5.0.0`-Default**. Erst analysieren.

Bis die Prüfung Adaption-für-Adaption erfolgt ist (Etappe C, §2.6), bleiben die
deklarierten `MR-*` **in Kraft**: mehrere sind maschinell gegatet (die
`ids`-Muster in [`.d-check.yml`](../../../../.d-check.yml) erzwingen die
Linkpflicht auf sie), und ein pauschaler Sofort-Vorrang stellte zwei einander
ausschließende Regeln nebeneinander.

## 2. Entscheidungen / Regel

### 2.1 Umfang des Sprungs (gemessen)

Gelesen wurde das **kanonische Artefakt** — das `v5.0.0`-Release-Bundle
(`lab-regelwerk.zip`, netzlos in den Scratchpad geladen; trägt `regelwerk/` **und**
`templates/` parallel), nicht der Git-Baum des Kurses.

- **Datei-Bestand (der teuerste Teil):** grundlagen **3 → 8 Dateien** —
  `grundlagen-konventionen.md` (unter `v1.4.0` ~29 % des Regelwerks) ist
  **aufgesplittet und entfallen**, an ihre Stelle treten `grundlagen-begriffe`,
  `-bootstrap`, `-harness-dateien`, `-referenz-richtung`, `-source-precedence`,
  `-traceability`. Zwei Module **umbenannt** (`modul-03-lastenheft` →
  `modul-03-spec`, `modul-04-architektur-adrs` → `modul-04-adrs`). **Alle 17
  Module** und der Index sind inhaltlich verschieden (Datei-für-Datei-Diff der
  beiden Bundles).
- **Umfang:** Gesamt **4030 → 3851 Zeilen** (netto −179; weiter **verdichtet**,
  nicht additiv — die Annahme „unser Stand bleibt gültig" trägt erst recht
  nicht). Der Löwenanteil der Verdichtung bleibt **Didaktik-Abbau**
  (`Worked Example …` → `Ziel-Form: …` mit Vorlagen-Verweis), nicht Regel-Abbau;
  das entlastet die Migration, ersetzt aber die Modul-Lektüre nicht (§2.3).
- **Zwei Majors über der v3.5.2-Analyse:** `v4.0.0` (Welle 62) — Asset-Umbenennung
  Module 03/04, drei Spec-Straten Pflicht, Referenzmatrix 8×8, Adaptions-Block aus
  dem Pflicht-Lesepfad heraus; `v5.0.0` (Welle 64) — `konventionen.md`-Wegweiser
  ersatzlos entfernt (der Split selbst war `v4.1.0`, additiv, mit Weiterleitung).
  Beide sind nach der Asset-/Layout-Policy des Kurses MAJOR.

**Kostensenker — jetzt nur eingeschränkt:** Das Regelwerk war zwischen `v1.3.0`
und `v1.4.0` **byte-identisch**, weshalb a-checks Sprung `v1.3.0 → v3.5.2` und
d-checks `v1.4.0 → v3.5.2` **derselbe** Regelwerk-Diff waren (a-checks Analyse 1:1
übertragbar). Mit dem Ziel `v5.0.0` gilt das nur noch für die v3.5.2-Untergrenze:
a-check stand (Erhebung 2026-07-25) selbst auf `v3.5.2` und hat den v4/v5-Sprung
also **nicht** notwendig hinter sich; ob seine Analyse für `v5.0.0` trägt, ist
erst nach Prüfung des Flotten-Stands (Etappe B) sicher.

### 2.2 Die sechs Brüche

1. **`agents-regelwerk.md` ist retired** (v2.0.0, BREAKING — der Modul-Split ist
   jetzt kanonisch). Zwei Live-Fundstellen führen es weiterhin als adoptierte
   Quelle: die Guides-Zeile in [`harness/README.md`](../../../../harness/README.md)
   und §Adoptierte Konventions-Quellen in
   [`harness/conventions.md`](../../../../harness/conventions.md).
2. **Das Bundle ist self-contained** (v3.0.0): ein Asset trägt `regelwerk/` und
   `templates/` als Geschwister; `lab-templates.zip` **existiert nicht mehr**.
   Die `../templates/…`-Verweise der Module („Ziel-Form: X") lösen netzlos nur
   auf, wenn beide Bäume parallel vendored sind.
3. **Das Materialisierungs-Skript bricht dreifach** auf dem neuen Layout: das
   Entpacken in den Regelwerk-Pfad erzeugt eine Doppel-Verschachtelung, das
   Manifest-Glob greift dann ins Leere; der zweite Entpack-Lauf lädt ein Asset,
   das es nicht mehr gibt (Abbruch); und die Content-Drift-Prüfung des
   Currency-Audits meldete aus demselben Glob-Grund einen **Fehl-Alarm**.
   Referenz-Lösung liegt vor: u-boots Fassung desselben Skripts löst das Layout
   tolerant auf (Regelwerk am Modul-0-Abschnitt erkennen, Templates als
   Geschwister fordern), vendored beide Bäume, bildet das Manifest über den
   tatsächlichen Dateibestand und zieht vor dem Manifest eine
   **Under-Copy-Barriere** (Quell-Dateizahl = vendorte Dateizahl).
4. **Pin-gebundene Verweise** ([`MR-021`](../../../../harness/conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden)):
   Entfernt der Bump das alte Tag-Verzeichnis — laut derselben Adaption
   Pflicht-Teil —, laufen **sechs Live-Fundstellen** (Briefing, Harness-Einstieg,
   Konventionen, [Planning-Index](../README.md), [Roadmap](../in-progress/roadmap.md)
   und der Reviewer-Skill `.harness/skills/reviewer.md`, Z. 6 + 39 — von der
   v3.5.2-Analyse übersehen, s. Bruch 6) **und drei historische** ins Leere:
   [ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md) (Status
   `Accepted`, also immutabel nach
   [ADR-0016](../../adr/0016-adr-immutable-gate.md) /
   [ADR-0024](../../adr/0024-vcs-immutable-gate.md)),
   [slice-080](../done/slice-080-sources-modul.md) und
   [slice-081](../done/slice-081-pins-hash-ergonomie.md).
5. **Module 03/04 umbenannt** (`v4.0.0`, MAJOR): `modul-03-lastenheft.md` →
   `modul-03-spec.md`, `modul-04-architektur-adrs.md` → `modul-04-adrs.md`. Bricht
   das vendorte Layout ein zweites Mal. **Keine** Live-d-check-Fundstelle nennt
   die Modul-Dateinamen (geprüft) — nur der vendorte `v1.4.0`-Baum trägt sie, den
   Etappe A ohnehin ersetzt; das tolerante Entpacken (Bruch 3) muss die neuen
   Namen aber kennen.
6. **`grundlagen-konventionen.md` entfällt ersatzlos** (`v5.0.0`, MAJOR): die eine
   große Konventions-Datei ist in **sechs neue** `grundlagen-*`-Dateien aufgesplittet (grundlagen damit 3 → 8)
   (`v4.1.0`, additiv, mit Wegweiser), der Wegweiser dann entfernt (`v5.0.0`). Das
   ist die **neue**, in #4 nachgetragene Live-Fundstelle: der Reviewer-Skill
   verlinkt `grundlagen-konventionen.md#referenz-richtung-sdp-…` — der Inhalt
   wandert nach `grundlagen-referenz-richtung.md`, der Link muss **retargetet**
   werden. Die v3.5.2-Analyse hatte diesen Bruch nicht (dort existierte die Datei
   noch). Dieselbe Anker-Verschiebung trifft die immutable
   [ADR-0022](../../adr/0022-matrix-token-richtung-provenance-marker.md) (§2.2 #4).

### 2.3 Was `v5.0.0` inhaltlich verlangt (die echten Zugänge)

- **Roadmap-Struktur, fünf Abschnitte** (Modul 6): *Aktuelle Welle · Nächste
  Wellen · Meilensteine · Abgeschlossene Wellen · Historische
  Trigger-Verschiebungen*. d-checks Roadmap führt **drei** davon; das
  Closure-Log und die Meilenstein-Ebene fehlen, der geschlossene Bestand steckt
  in der Prosa-Kette unter §Aktuelle Welle.
- **Wellen-Closure-Prozedur, fünf Schritte** (Modul 6): Trigger prüfen ·
  Carveout-Audit · Welle nach `done/` schließen (Closure-Notiz *und* `git mv`
  der Welle-Plan-Datei) · ein Wave-Self-Close-Commit · Roadmap fortschreiben.
  d-check führt Wellen heute ausschließlich in der Roadmap; es gibt weder
  Welle-Plandateien noch Wellen-Closure-Notizen.
- **Freshness-Audit als deklarierte Routine** (Modul 2, seit v3.2.0
  prozeduralisiert): beobachtbarer Auslöser statt Kalenderpflicht, Netz-Operation
  **außerhalb** der Gates, und die **Release-Liste** prüfen statt des Assets.
  d-check hat den Sensor, aber die Routine (Kadenz, Zuständigkeit, Nicht-Ziele)
  lebt nur als Prosa in
  [`MR-022`](../../../../harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019)
  statt als eigener Konventions-Abschnitt.
- **Change-Request-Fußabdruck** (neu in v3.5.2): „Change Request" ist bewusst
  **kein** Harness-Konstrukt — kein eigenes ID-Schema, keine Datei, kein Gate;
  ein angenommener CR hinterlässt nur einen Version-Bump des Lastenhefts, eine
  Zeile in dessen Historie und die geänderten Anforderungen. Hard Rule: weder ADR
  noch Slice dürfen Lastenheft-Kennungen je ändern. Gegen d-checks
  Anforderungs-Anlege-Prozess zu prüfen (Erst-Eindruck: erfüllt).
- **Fehlende Artefakte:** die Vorlage für einen `closure-note-reviewer`-Skill
  (semantische Schicht über einem Struktur-Gate für Closure-Notizen) hat kein
  Gegenstück im Repo; die Review-Report-Vorlage verlangt Kopf-Felder
  (**Review-Art**, **Skill-Version**, **Modell-ID**), die d-checks Reports heute
  nicht führen; zudem fehlt das **Beobachtungs-Register** (`observations.md`,
  v5.0.0-Standard-Artefakt mit eigener Vorlage) ganz.
- **Drei Spec-Straten Pflicht** (`v4.0.0`, Modul 3): Lastenheft, Spezifikation und
  Architektur sind jetzt Baseline-**Default**, nicht mehr optional. d-check führt
  sie bereits (deklariert als [`MR-001`](../../../../harness/conventions.md#mr-001--source-precedence-mit-eigener-spezifikations-schicht)) —
  die Adaption wird damit zum Default und entfällt (§2.4).
- **Referenzmatrix 8×8** (`v4.0.0`): die erlaubten Referenzrichtungen zwischen den
  Dokumentklassen sind auf eine 8×8-Matrix erweitert. d-checks `matrix`-Config und
  [`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)
  sind gegen die neue Matrix zu stellen (Etappe B) — Erst-Eindruck deckungsgleich.
- **Adaptions-Block außerhalb des Pflicht-Lesepfads** (`v4.0.0`): der
  `MR-*`-Adaptions-Block gehört nicht mehr in den obligatorischen Lesepfad. d-check
  führt ihn heute mitten in `harness/conventions.md`; die Neuform trennt
  „einmal lesen" von „bei Bedarf nachschlagen" — eine Struktur-Angleichung für
  Etappe D.
- **grundlagen-Navigation je Abschnitt** (`v5.0.0`): der Split (8 Dateien) macht
  „nur den benötigten Abschnitt laden" datei-granular. Jeder d-check-Zeiger auf
  grundlagen-Inhalt (Referenz-Richtung, Source-Precedence, Traceability, Begriffe)
  muss auf die **richtige** Split-Datei zeigen, nicht mehr auf die eine große.
- **`conventions.md`: Index statt Inline-Block, eine Datei je Adaption** (`v4.0.0`
  „Adaptions-Block verlässt den Lesepfad"; kanonisch: Regelwerk
  `grundlagen-harness-dateien.md` §Konventionsspeicher). Der **Default** ist jetzt
  die **Verzeichnis-Form**: `conventions.md` trägt nur noch den **Index**
  (Adoptions-Erklärung + eine Tabellenzeile je Adaption); die Einträge leben je
  Adaption in **einer eigenen Datei** (`harness/conventions/MR-<NNN>-<titel>.md`),
  ein aufgelöster Eintrag wandert per `git mv` nach `conventions/done/` — **Zustand = Verzeichnis-Position,
  kein Status-Feld** (dieselbe Lifecycle-Form wie Slices). Grund: `conventions.md`
  liest **jeder** Lauf; inline wächst der Pflicht-Kontext mit jeder — auch jeder
  *aufgelösten* — Adaption, und ein aufgelöster Eintrag liest sich wie ein
  geltender. **d-check trifft das hart:** die 23 Einträge (die Adoptions-Erklärung
  + 22 Adaptionen/Forks) liegen heute inline als Prosa in **einer** `conventions.md`.
  **Anker-Kaskade (Review-F-2/N-1):** der Split verlagert die `### MR-…`-Heading-Anker
  aus `conventions.md` — **173** `conventions.md#mr-…`-Links in **57** Dateien würden
  brechen; **12** liegen in Accepted-**immutablen** ADRs (nicht editierbar), und
  **zehn** dieser ADR-Links zeigen auf **nicht-aktive** MRs (aufgelöst/entfällt/
  verschmilzt), deren Anker aus dem Index verschwinden. Die Links sind **Voll-Slug**
  (`#mr-NNN--voller-titel-slug`, nicht kurz `#mr-NNN`). **Gegenmittel** (Etappe C):
  `conventions.md` behält je **von immutabler/eingefrorener Doku referenziertem** MR
  einen **Voll-Slug-`<a id>`-Anker** — **auch für die aufgelösten** — in einem eigenen
  **Anker-Kompatibilitäts-Block** (unabhängig vom aktiv-/`done/`-Schnitt), sodass alle
  173 Links **ohne Retarget und ohne ADR-Edit** auflösen. Das ist eine d-check-
  **migrationsspezifische** Maßnahme (ein frisches v5.0.0-Repo hat keine Alt-Links) und
  wird als eigener `MR`-Eintrag deklariert.
  **Nebeneffekt:** erst die Einzeldatei-Form macht
  die Append-only-Disziplin *pinbar* — d-checks `immutable`/`vcs`-Module können dann
  je Eintrag gegen Core-Drift wachen, was eine Sammeldatei nicht kann.
- **Neue Pflichtfelder je Eintrag** (Template `conventions/MR-NNN-titel.template.md`):
  vor allem **`Ersetzt-Baseline-Regel`** — **genau eine** Baseline-Regel, als Link
  mit Abschnitts-Anker in die vendored v5.0.0-Fassung; **wer keine benannte Regel
  ersetzt, ist ein Fork, keine Adaption** (kanonisch enger:
  `grundlagen-source-precedence.md`). Dazu `Status: Accepted` und — wenn ein
  Eintrag einen früheren ablöst — `Löst auf` + `Ausgelöst durch Baseline-Stand`;
  *schärft* er ihn nur (der alte gilt weiter, strenger), steht `(schärft …)` im
  Titel. d-checks Einträge führen heute **keines** dieser Felder; die Migration
  ergänzt sie je Eintrag und **reklassifiziert die Forks** (§2.4).
- **AGENTS.md gegen `AGENTS.template.md` angleichen** (`v5.0.0`; kanonisch
  `modul-09-implementierung.md` §AGENTS.md-Regeln). **Struktur deckungsgleich**
  (§1–§6, §3.1–§3.6), aber der **§1-Inhalt driftet:** der Kanon-Zeiger auf Modul 9
  §AGENTS.md-Regeln fehlt, die Baseline-URL und das `{regelwerk,templates}`-Layout
  stehen auf `v1.4.0`, und §1 trägt v1.4.0-/`MR`-gebundenen Text (die MR-Verweise
  migrieren mit). §3.1 (d-checks „Docker/make-only" statt Template-„Docker-only")
  und §3.4 (d-checks „Spec-Straten nie abwärts") sind Zusätze, die an Adaptionen
  hängen (§2.4). **Etappe D.**

### 2.4 Der `MR-*`-Bestand unter „Default gewinnt"

**Erst-Einschätzung** aus der Lektüre der neuen Vorlagen — verbindlich erst nach
der Prüfung in Etappe C:

| Adaption | Erste Einschätzung gegen `v5.0.0` |
|---|---|
| [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage) Baseline-Aussage | bleibt — die Vorlage kennt den Abschnitt weiterhin; Inhalt abgleichen |
| [`MR-001`](../../../../harness/conventions.md#mr-001--source-precedence-mit-eigener-spezifikations-schicht) eigene Spezifikations-Schicht | **entfällt** — `v4.0.0` macht die drei Spec-Straten (Lastenheft/Spez/Architektur) zum Pflicht-Default (§2.3) |
| [`MR-002`](../../../../harness/conventions.md#mr-002--id-schema-mit-bereichskürzeln-ab-initialer-fassung) ID-Schema mit Bereichskürzeln | **Kandidat zum Entfall** — Default-Schema prüfen |
| [`MR-003`](../../../../harness/conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh) vendorter Bootstrap-Sensor (aufgelöst) | **Nummern-Kollision**: die Vorlage belegt dieselbe Nummer mit der Vendoring-Adaption |
| [`MR-004`](../../../../harness/conventions.md#mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild) Gate-Nachweis-Mechanik | bleibt — repo-eigene Hook-Mechanik, von der Baseline nicht geregelt |
| [`MR-005`](../../../../harness/conventions.md#mr-005--härtung-ggü-b-cad-inhaltsbasierter-gate-nachweis-sub-shell-prüfung) inhaltsbasierter Gate-Nachweis | bleibt — wie oben |
| [`MR-006`](../../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs) Referenzrichtung | prüfen gegen die 8×8-Matrix (`v4.0.0`) + `grundlagen-referenz-richtung.md` (`v5.0.0`, eigene Datei); Erst-Eindruck deckungsgleich |
| [`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding) Auflösung der Vorgänger-Adaption | historisch — Provenienz |
| [`MR-008`](../../../../harness/conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage) ID-Schema-Deklaration | **Kandidat zum Entfall** — die Vorlage fordert die Deklaration selbst |
| [`MR-009`](../../../../harness/conventions.md#mr-009--source-precedence-ohne-docsuser-rang) / [`MR-010`](../../../../harness/conventions.md#mr-010--auflösung-von-mr-009-docsuser-rang-eingefügt) `docs/user`-Rang | historisch (die spätere löst die frühere auf) |
| [`MR-011`](../../../../harness/conventions.md#mr-011--baseline-auf-release-tag-gepinnt) / [`MR-012`](../../../../harness/conventions.md#mr-012--baseline-pin-hebung-nachtrag-zu-mr-011) / [`MR-016`](../../../../harness/conventions.md#mr-016--baseline-pin-hebung-zweiter-nachtrag-zu-mr-011) Pin + Hebungen | Historie; die Bump-Prozedur zieht diese Migration nach |
| [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise) Lifecycle-Move-Commit | **prüfen** — die Baseline führt „`git mv` + Inhaltsänderung = zwei Commits" als Hard Rule; d-checks Bündelungs-Ausnahme ist gate-getrieben und bleibt vermutlich |
| [`MR-014`](../../../../harness/conventions.md#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template) Slice-/ADR-Haus-Stil | **entfällt** — sagt wörtlich „Haus-Stil **ggü. Baseline-Template**"; Folgen in §2.5 |
| [`MR-015`](../../../../harness/conventions.md#mr-015--auflösung-der-mr-012-pointer-drift-agentsmd-routet-spiegelt-nicht-mehr) AGENTS.md routet | bleibt — repo-eigene Pointer-Disziplin |
| [`MR-017`](../../../../harness/conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen) Cache aus dem Selbst-Scan | **entfällt** — der Cache-Zweig verschwindet mit dem retireten Asset |
| [`MR-018`](../../../../harness/conventions.md#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates) keine Templates verkörpert | **entfällt vollständig** — Folgen in §2.5 |
| [`MR-019`](../../../../harness/conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017) Regelwerk committet vendored | **verschmilzt** in den neuen Default (die Vorlage führt das Vendoring selbst) |
| [`MR-020`](../../../../harness/conventions.md#mr-020--baseline-template-propagation-per-drift-audit-template-frei-bestätigt) Template-Propagation, template-frei bestätigt | **entfällt** mit der Template-Freiheits-Adaption darüber |
| [`MR-021`](../../../../harness/conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden) pin-gebundene Verweise | bleibt — repo-eigene Gate-Kopplung; **Prosa entpinnen** (Version steht einmal im Pin) |
| [`MR-022`](../../../../harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019) Currency-Audit-Modus | **verschmilzt** in den Konventions-Abschnitt §Freshness-Audit (§2.3) |

**Zwei neue v5.0.0-Achsen je Eintrag** (Etappe C, aus §2.3), quer zur Tabelle:
1. **`Ersetzt-Baseline-Regel` → Adaption oder Fork.** Jeder Eintrag muss **genau
   eine** Baseline-Regel benennen, die er ersetzt (anker-verlinkt in den vendored
   Baum). Die **repo-eigenen, von der Baseline nicht geregelten** Einträge (die
   Gate-Nachweis-Mechanik + ihre Härtung, die AGENTS.md-Pointer-Disziplin — oben als
   „bleibt" markiert) ersetzen **keine** Baseline-Regel und sind damit **Forks**,
   keine Adaptionen. Wo Forks dokumentiert werden (Fork-Register vs. repo-lokaler
   Regel-Ort), ist ein Etappe-C-/Abnahme-Detail.
2. **Supersede-/Schärfungs-Ketten.** Die oben als „historisch"/„Nachtrag" geführten
   Paare werden zu `Löst auf` + `Ausgelöst durch Baseline-Stand`-Feldern (Ablösung,
   Nachfolger nach `conventions/done/`) bzw. `(schärft …)` im Titel (der alte gilt
   weiter), statt Prosa-„Nachtrag zu …".

### 2.5 Die zwei Adaptionen, die am weitesten reichen

- **[`MR-018`](../../../../harness/conventions.md#mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates) fällt vollständig**, nicht nur im Cache-Pfad. Ihre
  Brücken-Begründung (der Kurs nenne d-check im Self-Hosting-/Producer-Fall und
  nehme es von der `harness.mk`-Adoption aus) ist **erloschen**: `harness.mk`
  wurde in v2.0.0 retired, und die Templates-Übersicht von `v5.0.0` führt keinen
  Self-Hosting-Abschnitt mehr. Ihr eigener Auflösungs-Trigger sah genau diese
  Prüfung beim nächsten Bump vor. **Folge:** d-check verkörpert die fünf
  wiederkehrenden Skelette künftig co-located (ADR, Slice, Welle, Carveout,
  Review-Report). **Gate-Folge:** [`.d-check.yml`](../../../../.d-check.yml)
  braucht den Suffix-Ausschluss für Vorlagendateien, sonst färben deren
  Platzhalter die Kennungs- und Referenz-Prüfung rot. `matrix`, `vcs` und
  `trace` sind **nicht** betroffen — deren Globs greifen auf Vorlagen-Dateinamen
  nicht.
- **[`MR-014`](../../../../harness/conventions.md#mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template) fällt ebenfalls.** Gemessene Struktur-Differenz: die
  Slice-Vorlage führt einen eigenen Abschnitt *Plan (vor Code)*, einen eigenen
  *Closure-Trigger* und die Sub-Area-Begründung als §8; die ADR-Vorlage führt
  Alternativen als Option-Prosa, die Fitness Function als Tabelle und die
  Geschichte dreispaltig. **Der Bestand bleibt unangetastet** — 46 ADRs (nach
  `Accepted` immutabel) und 82 abgeschlossene Slices sind Historie; die Vorlagen
  sind Kopiervorlagen für **neue** Instanzen. Keine Rückmigration; das ist keine
  Ausnahme von der Vorgabe, sondern ihre korrekte Reichweite.

### 2.6 Entwarnungen (geprüft, kein Handlungsbedarf)

- **Review-Harness weitgehend stabil:** die Finding-Kategorien, das Output-Schema
  und die Negativbefund-Pflicht standen in `v3.5.2` wie gehabt (gegen `v5.0.0` in
  Etappe B zu bestätigen — mehrere Wellen berührten den Review-Harness, u. a. die
  Review-Report-Kopf-Felder aus §2.3 und ein `closure-note-reviewer`-Skill). *Neu und schärfer:*
  die repo-konkrete Liste schwerer Befunde muss **mindestens zwei
  repo-spezifische Regeln** nennen; die Skill-Datei ist daraufhin zu prüfen.
- **Lifecycle-Verzeichnisse** (`open/` → `next/` → `in-progress/` → `done/`) und
  der **Ort der Roadmap** entsprechen der Baseline-Vorlage (`v5.0.0`); `done/` archiviert
  dort zusätzlich abgeschlossene Nicht-Slice-Records (Welle-Closure, aufgelöste
  Carveouts) — genau der Ort, den die Wellen-Closure-Prozedur braucht.
- **Der Content-Drift-Teil des Currency-Audits ist Baseline-konform**, nicht
  Baseline-abweichend: Modul 2 nennt Release-Listen-Prüfung und
  Asset-Integrität ausdrücklich nebeneinander („ersetzt die
  Release-Listen-Prüfung nicht"). Er bleibt erhalten; nur das Verfahren des
  Currency-Teils zieht auf die Release-Liste um.

### 2.7 Migrationsschritte — der Etappen-Schnitt im Detail

Überblick zuerst, dann jeder Schritt einzeln. A ist vollständig spezifiziert (das
Layout ist gemessen); B produziert die Findings, die C und D verbindlich machen —
dort ist die **Prozedur** präzise, der **Inhalt** je Fund. Jede Etappe ist ein
eigener reviewbarer Commit-Bogen und schließt mit grünem `make gates`.

| Etappe | Inhalt | Warum getrennt |
|---|---|---|
| **A — Vendoring** | Skript auf das neue Bundle-Layout heben (tolerantes Entpacken, beide Bäume, Manifest über den tatsächlichen Bestand, Under-Copy-Barriere, Currency-Teil auf die Release-Liste); `.harness/baseline/<tag>/` neu materialisieren; Alt-Stand entfernen; Pin und Pointer nachziehen; die drei historischen Verweise (§2.2) über ein gescoptes Referenz-Ventil abfangen **und** die neue Live-Fundstelle Reviewer-Skill retargeten (`konventionen.md`→`referenz-richtung.md`, §2.2 Bruch 6); Vendoring-Adaption unter der nächsten freien `MR`-Nummer deklarieren | mechanisch, sofort prüfbar, **Voraussetzung** für alles Weitere: danach arbeitet jeder Schritt netzlos gegen die vendorte Quelle statt gegen eine URL |
| **B — Modul-Delta lesen** | die substanziell umgeschriebenen Module gegenlesen (Priorität: 2, 5, 6, 7, 10, 11, 13) und die Treffer als Findings sammeln | Lesearbeit; das Ergebnis bestimmt C und D |
| **C — `MR-*`-Bereinigung + Datei-Migration** | die 23 Adaptionen in Einzeldateien splitten (Index + `conventions/` + `done/`, **Anker-Erhalt** für die 173 Links), die neuen Pflichtfelder ergänzen, **Forks** reklassifizieren, entwertete streichen, Nummern-Kollision auflösen, Prosa entpinnen | berührt die Konventions-Identität des Repos — eigener, reviewbarer Schritt |
| **D — Form-Konformität** | Roadmap auf fünf Abschnitte, Wellen-Closure-Prozedur, **AGENTS.md gegen Template angleichen**, fehlende Artefakte + `observations.md` anlegen (§2.3), Vorlagen co-located (§2.5), Asset-Integrität dogfooden | breit, aber mechanisch; nach C, weil C die Form mitbestimmt |

**Reihenfolge-Argument:** A zuerst — genau die Eigenschaft, die das Vendoring
begründet (netzloser, reproduzierbarer Nachschlag), ist die Voraussetzung dafür,
dass B–D belegbar gegen die Quelle arbeiten.

#### Etappe A — Vendoring (der netzlose Boden). *Vollständig spezifiziert.*

1. **Skript heben** — `tools/harness/fetch-baseline-cache.sh` auf das
   v5.0.0-Bundle-Layout (Referenz: u-boots Fassung desselben Skripts): tolerantes
   Entpacken (`regelwerk/` am Modul-0-Abschnitt erkennen, `templates/` als
   Geschwister im **selben** Asset fordern — das Bundle trägt beide parallel);
   Manifest über den **tatsächlichen** Dateibestand statt festem Glob;
   **Under-Copy-Barriere** (Quell-Dateizahl == vendorte Dateizahl vor dem
   Manifest); der `--check-latest`-Content-Teil auf die Release-Liste statt
   Asset-Glob (behebt den v3.5.2-Fehl-Alarm, §2.2 Bruch 3).
2. **Materialisieren** — `fetch-baseline-cache.sh v5.0.0`: `.harness/baseline/v5.0.0/`
   mit **beiden** Bäumen (`regelwerk/` = 8 grundlagen + 17 Module + Index,
   `templates/`), `SHA256SUMS` über den vendorten Bestand, `--verify` grün.
3. **Alt-Stand entfernen** — `.harness/baseline/v1.4.0/` löschen (der Wegfall des
   alten Tag-Verzeichnisses ist Adaptions-Pflicht, §2.2 Bruch 4).
4. **Pin + Pointer** — `§Baseline`-Stand in `harness/conventions.md` `v1.4.0` →
   `v5.0.0`; die **sechs** Live-Fundstellen (§2.2 #4) auf den neuen Pfad ziehen
   (Briefing, Harness-Einstieg, Konventionen, Planning-Index, Roadmap,
   Reviewer-Skill). **Aber (Review-F-1):** die **entfallenen** Quellzeiger — der
   §Adoptierte-Quellen-Eintrag zur (v2.0.0 retireten) `agents-regelwerk.md` und der
   Zeiger auf die (v5.0.0 entfallene) Kurs-Konventionsdatei — lassen sich **nicht**
   retargeten; sie werden in `conventions.md` §Adoptierte Quellen und
   `harness/README.md` §Guides auf das vendorte Modul-Bundle **umgeschrieben** bzw.
   entfernt, **nicht** auf tote v5.0.0-Ziele umgehängt.
5. **Reviewer-Skill retargeten** — `…/grundlagen-konventionen.md#referenz-richtung-sdp-…`
   → `…/v5.0.0/regelwerk/grundlagen-referenz-richtung.md#…`; den Ziel-Anker im
   Split-File verifizieren (§2.2 Bruch 6).
6. **Historische Verweise abfangen** — die drei eingefrorenen Fundstellen
   (§2.2 #4: eine immutable ADR und zwei `done/`-Slices) auf den entfallenen Alt-Pfad
   über ein gescoptes `ignore-refs`/Tombstone (Präzedenz: d-checks
   Tombstone-Register für entfernte Artefakte).
7. **Vendoring-Adaption deklarieren** — unter der **nächsten freien** `MR`-Nummer
   (die konkrete Nummer erst hier vergeben — kein Vorgriff auf einen noch nicht
   existierenden Anker); die Nummern-Kollision mit der Baseline-Vendoring-Adaption
   (§2.4) löst Etappe C.
8. **Gate** — `make gates` grün. Ab hier arbeiten B–D netzlos gegen die vendorte
   v5.0.0-Quelle.

#### Etappe B — Modul-Delta lesen (die Findings). *Prozedur präzise, Inhalt je Fund.*

1. **Gegenlesen** — jede der 8 `grundlagen-*`-Dateien und jedes der 17 Module gegen
   `v5.0.0`, Priorität nach §2.1 (substanziell umgeschrieben zuerst: der
   grundlagen-Split, Modul 2/5/6/7/10/11/13, die umbenannten `modul-03-spec`/
   `modul-04-adrs`; Modul 0/9/Durchsetzungsschicht zuletzt).
2. **Finding-Schema** — je Treffer ein Eintrag: *{Quelle (Modul/§), Regel-Delta,
   betroffene d-check-Adaption/-Artefakt, Handlung → C oder D}*.
3. **Flotten-Stand** — u-boot/a-check/ai-harness-init auf `v5.0.0`? Das bestimmt,
   ob a-checks Analyse noch überträgt (§2.1-Kostensenker).
4. **Frischkontext-Review Pflicht** (§4) — der Bump ist ein Re-Adopt, kein Patch.
   Ergebnis: die Finding-Liste, die C und D speist.

#### Etappe C — `MR-*`-Bereinigung **und** Datei-Migration (die Konventions-Identität).

1. **Verbindlich machen** — die 23 §2.4-Einschätzungen gegen die B-Findings festziehen.
2. **In Dateien splitten** (§2.3) — die inline-Adaptionen aus `conventions.md` in
   Einzeldateien `harness/conventions/MR-<NNN>-<titel>.md` überführen (**eine je
   Adaption**); in `conventions.md`
   bleibt der **Index** (Adoptions-Erklärung + eine Zeile je aktiver Adaption);
   aufgelöste Einträge per `git mv` nach `conventions/done/` (Zustand =
   Verzeichnis-Position). **Anker-Erhalt (Review-F-2/N-1):** `conventions.md` behält
   je referenziertem MR einen **Voll-Slug**-`<a id="mr-NNN--voller-titel-slug">`-Anker
   — für **alle** von immutabler/eingefrorener Doku verlinkten MRs, **auch die
   aufgelösten/entfallenen** (zehn der zwölf immutablen ADRs zeigen auf nicht-aktive
   MRs; **plus** die nur von eingefrorenen **Review-Reports** referenzierten wie
   `mr-015`/`mr-016` — das **Inventar** ist repo-weit über **alle** eingefrorene Doku:
   ADRs, `done/`-Slices **und** `docs/reviews/`, das `anchors` ebenfalls prüft, nicht
   nur die ADRs), als eigener **Anker-Kompatibilitäts-Block** **unabhängig** vom aktiv-/`done/`-
   Schnitt der Index-Zeilen — damit Schritt 5 (Streichen/Verschieben) die Anker nicht
   mitnimmt. So lösen die **173** `#mr-…`-Voll-Slug-Links **ohne Retarget, ohne
   ADR-Edit** auf; der Block wird als eigener `MR` deklariert (migrationsspezifisch).
3. **Neue Pflichtfelder** je Eintrag (§2.3): `Ersetzt-Baseline-Regel` (genau **eine**
   v5.0.0-Regel, Anker-Link in den vendored Baum), `Status: Accepted`, und für
   Ablösungen `Löst auf` + `Ausgelöst durch Baseline-Stand` bzw. `(schärft …)` im Titel.
4. **Forks reklassifizieren** (§2.4) — die repo-eigenen, baseline-ungeregelten
   Einträge (Gate-Nachweis-Mechanik + Härtung, Pointer-Disziplin) ersetzen keine
   Baseline-Regel → **Fork**, keine Adaption; ihr Ort ist Abnahme-Punkt.
5. **Entfall/Verschmelzung** — die in §2.4 als **entfällt** markierten Adaptionen
   streichen (eigene Spec-Straten-Schicht, jetzt Default; Template-Freiheit;
   Cache-Ausnahme) und die als **verschmilzt** markierten (Vendoring, Currency-Audit)
   in den Default überführen.
6. **Nummern-Kollision** (§2.4) — die Baseline belegt dieselbe `MR`-Nummer mit ihrer
   Vendoring-Adaption wie d-checks historische, aufgelöste Vorgänger-Adaption →
   Entscheid (Abnahme-Punkt): eigene Nummern behalten (Provenienz) oder angleichen.
7. **Prosa entpinnen** (§2.2 Bruch 4) + **Referenzrichtung** (§2.3/§2.4) gegen die
   **8×8-Matrix**, `matrix`-Config ggf. nachziehen.
8. **Gate** — `make gates` + `make adr-check` grün; die neuen `conventions/`-Dateien
   erfüllen die `ids`-Linkpflicht (jeder `MR`-Verweis verlinkt) und sind — sobald
   `Accepted` — je Datei gegen Core-Drift **pinbar**. Keine `Accepted`-ADR wird
   inhaltlich berührt: die 173 `#mr-…`-Voll-Slug-Links bleiben via den
   **Anker-Kompatibilitäts-Block** gültig (Schritt 2 — auch für die aufgelösten,
   von immutabler Doku referenzierten MRs), die eingefrorene ADR-Fundstelle auf den
   Alt-Regelwerk-Pfad via **Tombstone** aus A6 — die zwei Mechaniken, ohne die C8s
   Zusage nicht erreichbar wäre.

#### Etappe D — Form-Konformität (die Struktur).

1. **Roadmap fünf Abschnitte** (§2.3) — die fehlenden zwei (*Meilensteine*,
   *Abgeschlossene Wellen*/Closure-Log) ergänzen; den geschlossenen Bestand aus der
   Prosa-Kette befüllen **oder** ab nächster Welle führen (Abnahme-Punkt, §2.8).
2. **Wellen-Closure-Prozedur** (fünf Schritte, §2.3) verankern — Welle-Plandateien,
   Wellen-Closure-Notizen, der Wave-Self-Close-Commit.
3. **AGENTS.md angleichen** (§2.3) — gegen `AGENTS.template.md`: den §1-Inhalt auf
   `v5.0.0` (Kanon-Zeiger Modul 9 §AGENTS.md-Regeln, `{regelwerk,templates}`-URL +
   Layout, die mit-migrierten `MR`-Verweise); die d-check-Zusätze §3.1/§3.4 als
   eigene Einträge führen (§2.4). Die Struktur (§1–§6) bleibt deckungsgleich.
4. **Fehlende Artefakte** anlegen — `closure-note-reviewer`-Skill-Vorlage;
   Review-Report-Kopf-Felder (*Review-Art*, *Skill-Version*, *Modell-ID*); das
   **Beobachtungs-Register** `observations.md` (v5.0.0-Vorlage).
5. **Vorlagen co-located** (Template-Freiheits-Folge, §2.5) — die fünf wiederkehrenden Skelette
   (ADR, Slice, Welle, Carveout, Review-Report) als `*.template.md`; `.d-check.yml`
   um den **Suffix-Ausschluss** für Vorlagen-Platzhalter ergänzen, sonst färbt deren
   Platzhalter-Kennung `ids`/`links` rot.
6. **Adaptions-Block** aus dem Pflicht-Lesepfad ziehen (§2.3, `v4.0.0`).
7. **Asset-Integrität** per eigenem Modul dogfooden (der Bundle-Wächter, Kurs-Welle 53).
8. **Freshness-Audit** als eigener Konventions-Abschnitt (Currency-Audit-Verschmelzung, §2.3).
9. **Gate** — `make gates` grün.

**Granularität der Umsetzung** (Abnahme-Punkt): A–D können **ein Slice je Etappe**
sein (die ursprüngliche Zerlegung) oder — wenn dies „die Migration" wird — als eine
Welle mit vier reviewbaren Etappen-Commits laufen. A ist in beiden Fällen der erste,
gate-belegte Schritt; B ist der Fixpunkt, an dem der unabhängige Review sitzt.

### 2.8 Was diese Analyse **nicht** geklärt hat

- **Modul-für-Modul-Delta gegen `v5.0.0`:** vermessen ist der Datei-Bestand beider
  Bundles (§2.1); inhaltlich gegen `v5.0.0` gegengelesen ist noch **nichts**. Die
  frühere v3.5.2-Lektüre (Grundlagen-Konventionen, Modul 2/6/10) ist durch die zwei
  Majors + den grundlagen-Split teils überholt (`konventionen.md` existiert nicht
  mehr). Das volle Gegenlesen aller 17 Module (inkl. der umbenannten
  `modul-03-spec`/`modul-04-adrs`) und der 8 grundlagen-Dateien ist **Etappe B**.
- **Ob die Entfall-Kandidaten in §2.4 wirklich entfallen** — das ist Lesearbeit
  am neuen Default, keine Messung.
- **Der Umgang mit dem geschlossenen Wellen-Bestand:** ob das Closure-Log aus der
  bestehenden Prosa-Kette befüllt wird (Zeiger auf Slice-Closures) oder erst ab
  der nächsten Welle geführt wird, entscheidet Etappe D.
- **Ob die Migration Gate-Änderungen erzwingt** über den Suffix-Ausschluss
  (§2.5) hinaus.
- **Die Fork-vs-Adaption-Zuordnung je Eintrag** und der genaue Ort der Forks
  (Fork-Register vs. repo-lokaler Regel-Ort): benannt (§2.4), aber die konkrete
  `Ersetzt-Baseline-Regel`-Zeile je Adaption ist Lesearbeit gegen die vendored
  v5.0.0-Fassung — Etappe C.

## 3. Definition of Done

- [x] Ist-Stand, Pin-Stellen und Flotten-Vergleich erhoben (§1).
- [x] Sprung-Umfang am kanonischen Artefakt gemessen, Rohmessung **und**
  normalisierte Messung ausgewiesen (§2.1).
- [x] Die sechs Brüche benannt und je an einer konkreten Fundstelle belegt (§2.2).
- [x] Die inhaltlichen Zugänge von `v5.0.0` gegen den Ist-Zustand gestellt (§2.3),
  inkl. der zwei Majors `v4.0.0`/`v5.0.0` über der ursprünglichen v3.5.2-Analyse.
- [x] Der vollständige `MR-*`-Bestand mit Erst-Einschätzung tabelliert (§2.4),
  die zwei weitreichenden Adaptionen einzeln begründet (§2.5).
- [x] Entwarnungen geprüft und benannt statt stillschweigend angenommen (§2.6).
- [x] Offene Lücken der Analyse ausgewiesen (§2.8) statt als Vollständigkeit
  ausgegeben.
- [x] **Jeder Migrationsschritt** je Etappe A–D präzise beschrieben (§2.7 Detail):
  A vollständig spezifiziert, B–D als präzise Prozedur mit Inhalt je Fund.
- [x] **Abnahme:** Etappen-Schnitt A–D (§2.7) und der Umgang mit den drei
  historischen Verweisen (§2.2) durch den Maintainer — **erteilt 2026-08-01** (nach R1/R2/R3).

## 4. Risiken / offene Punkte

- **Stiller Regel-Wechsel.** Eine Regeländerung in `v5.0.0` kann eine bestehende
  Adaption entwerten, ohne dass ein Gate ausschlägt. Genau dagegen steht die
  Gegenprobe in Etappe C — nicht der Skript-Lauf in Etappe A.
- **Der Bump ist keine Wartung.** Über 10 Releases mit zwei Major-Sprüngen ist
  dies ein Re-Adopt, kein Patch-Bump. Ein unabhängiger Frischkontext-Review ist
  für die inhaltlichen Etappen (B–D) **Pflicht**, nicht optional.
- **Historische Verweise vs. Alt-Stand-Entfernung.** Die Pin-Adaption verlangt
  das Entfernen des alten Tag-Verzeichnisses; drei Verweise darauf liegen in
  immutabler bzw. abgeschlossener Doku (§2.2). Das Referenz-Ventil ist der
  präzedenzgestützte Weg (das Repo führt bereits ein Tombstone-Register für
  entfernte Artefakte), aber es ist ein bewusster Entscheid, kein Automatismus.
- **Reihenfolge-Abhängigkeit.** Wird D vor C gezogen, entstehen neue Artefakte in
  einer Form, die C danach wieder anfasst.
- **Fremd-Repo-Vorlagen sind Anregung, kein Kanon.** u-boots Skript und a-checks
  Etappen-Schnitt sind gelebte Präzedenz in der Flotte, aber weder kanonische
  Quelle noch Baseline; bei Konflikt gilt die Baseline.

## 5. Trigger

Maintainer-Frage am 2026-07-25 („welche Regelwerksversion verwenden wir?") und
die anschließende Vorgabe: vollständige Migration, auch gegen bestehende `MR-*`.
Der Trigger ist damit gefeuert; die Einplanung der Etappen ist eine
Priorisierungs-Entscheidung, kein weiterer Trigger. **Neubasierung 2026-08-01:**
der Upstream sprang inzwischen auf `v5.0.0` (zwei weitere Majors); Maintainer-
Vorgabe „Analyse auf v5.0.0 neu basieren, dann Abnahme" — der Zielstand wandert
von `v3.5.2` auf `v5.0.0`, der Auftrag bleibt.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc führt. Berührt sind die Sub-Areas *Harness/Konventionen*
und *Harness-Tooling* — beide greenfield gewachsen, ohne Brownfield-Spec. Diese
Analyse ändert keine Artefakte; die Etappen A–D sind vor der Umsetzung je für
sich zu begründen.

## 7. Closure-Notiz (nach `done/`)

**Abgenommen 2026-08-01** (Maintainer). Reine Analyse — kein Code/Spec/Harness-Delta,
kein CR/ADR/Release.

**Review-Kette.** Drei unabhängige Frischkontext-Reviews: R1 (nicht abnahmereif, 2
MEDIUM) → geheilt; R2 (bestätigt F-1/LOW/INFO/AGENTS, fand aber N-1 = Fehler im
F-2-Fix) → geheilt (Voll-Slug-Anker-Kompatibilitäts-Block für **alle** von
eingefrorener Doku referenzierten MRs); R3 (eng auf N-1, am Anker-Engine-Code
belegt) = **abnahmereif**, ein LOW (F-R3-1, Inventar repo-weit) eingearbeitet. Drei
Review-Reports unter `docs/reviews/`.

**Abgenommener Schnitt.** Vier Etappen A (Vendoring) → B (Modul-Delta lesen) → C
(MR-Bereinigung + Datei-Migration) → D (Form-Konformität), **je ein Slice**;
Umsetzung ab welle-67 (Etappe A = slice-084).

**Entscheid historische Verweise.** Die drei eingefrorenen Regelwerk-Pfad-Verweise
(§2.2 #4: eine immutable ADR + zwei `done/`-Slices) auf den entfallenden
v1.4.0-Alt-Pfad → gescoptes `ignore-refs`/Tombstone (Etappe A6). Die 173
`conventions.md#mr-…`-Anker-Links (57 Dateien, 12 immutable ADRs, inkl.
Review-Report-referenzierte) → Voll-Slug-Anker-Kompatibilitäts-Block (Etappe C),
**kein ADR-Edit**.

**Commit-Kette.** `3677681` (Analyse) · `404a576` (R1) · `42bd3fe` (N-1/F-R3-1-Fix)
· `c94c44d` (R2/R3) · Abnahme-Move + dieser Body-Commit.

**Lerneintrag.** Der Wert des Frischkontext-Zwangs (§4): R2 fand einen echten Fehler
*im Fix* von R1, R3 prüfte den Fix von R2 — jede Runde härtete den kritischen
`conventions.md`-Anker-Mechanismus, den ein Ein-Pass-Review nicht erreicht hätte.
