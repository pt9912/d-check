# Slice slice-086: Regelwerk-Migration Etappe C — MR-Bereinigung + Datei-Migration

**Status:** Done (welle-67, Etappe C abgeschlossen 2026-08-02).

**Welle:** welle-67-baseline-v500-migration (dritte Umsetzungs-Etappe, nach
[slice-085](../done/slice-085-etappe-b-modul-delta.md)).

**Bezug:** Umsetzung von **Etappe C** des in
[slice-083](../done/slice-083-regelwerk-v500-migration-analyse.md) §2.7
abgenommenen Schnitts, **verbindlich gemacht** durch die
[slice-085](../done/slice-085-etappe-b-modul-delta.md)-Finding-Liste (§3.1
C-1…C-8 + die §3.4-Korrekturen). Berührt die **Konventions-Identität** des Repos.
**Kein Change Request**, **kein ADR**, **kein Release** — Harness-/Konventions-
Änderung.

**Autor:** pt9912. **Datum:** 2026-08-02.

---

## 1. Ziel

`harness/conventions.md` von der Inline-Sammeldatei auf die v5.0.0-**Default-Form**
heben: **Index + eine Datei je Adaption** (`harness/conventions/MR-<NNN>-<titel>.md`),
aufgelöste Einträge → `conventions/done/` (Zustand = Verzeichnis-Position), die
neuen Pflichtfelder je Eintrag, **Forks reklassifiziert** (nach dem in slice-085
§3.4 korrigierten Kriterium), entwertete Adaptionen gestrichen/verschmolzen, die
**173** `conventions.md#mr-…`-Links via **Voll-Slug-Anker-Kompatibilitäts-Block**
erhalten und den d-check-`matrix`-Scope gegen die **8×8-Matrix** als bewusste
Grenze deklariert (C-4). **Nur die Konventions-Identität** — die Form-Struktur
(die 11 D-Findings) ist Etappe D; die **Historie-Provenance-Enforcement (C-3)** ist
als Code-Feature in einen **Folge-Produkt-Slice** herausgeschnitten (§3).

## 2. Vorgehen (slice-083 §2.7 Etappe C, gefüttert durch die C-Findings)

1. **Verbindlich machen** — die §2.4-MR-Einschätzungen gegen die slice-085-Findings
   festziehen; besonders die **korrigierte** Fork-Frage (§3.4: die Guard-Härtungs-
   Adaptionen sind keine Forks; Fork = *pauschale* Nichtanwendbarkeit).
2. **In Dateien splitten** — die Inline-Adaptionen aus `conventions.md` je in eine
   `harness/conventions/MR-<NNN>-<titel>.md` überführen; im Index bleibt die
   Adoptions-Erklärung + eine Zeile je aktiver Adaption; aufgelöste per `git mv`
   nach `conventions/done/`. **Anker-Erhalt (C-5/Review-N-1):** `conventions.md`
   behält je **von eingefrorener Doku** (immutable ADRs, `done/`-Slices,
   `docs/reviews/`) referenziertem MR einen **Voll-Slug-`<a id>`-Anker** — auch für
   die aufgelösten — in einem eigenen **Anker-Kompatibilitäts-Block**, unabhängig
   vom aktiv/`done`-Schnitt. So lösen die 173 Links **ohne Retarget, ohne ADR-Edit**
   auf; der Block wird als eigener `MR` deklariert (migrationsspezifisch).
3. **Neue Pflichtfelder** (C-6) je Eintrag: `Ersetzt-Baseline-Regel` (genau **eine**
   v5.0.0-Regel, Anker-Link in den vendored Baum), `Status: Accepted`, für
   Ablösungen `Löst auf` + `Ausgelöst durch Baseline-Stand` bzw. `(schärft …)` im
   Titel.
4. **Forks reklassifizieren** (C-Schritt 4 / §3.4) — **Abnahme-Punkt (§3.1)**.
5. **Entfall/Verschmelzung** (C-1/C-2/C-8) — die entwerteten Adaptionen (eigene
   Spec-Straten-Schicht → Default; Spec-Decke → Default; Template-Freiheit/-Cache)
   streichen bzw. zu Provenienz abrüsten; Vendoring/Currency-Audit in den Default
   überführen.
6. **Nummern-Kollision** (§2.4) — **Abnahme-Punkt (§3.1)**.
7. **Prosa entpinnen + Referenzrichtung** — die interne `v1.4.0`-Prosa (slice-084
   LOW-2a: `<tag>=v1.4.0` in den Alt-Einträgen) auf `v5.0.0` ziehen; die
   Referenzrichtung gegen die **8×8-Matrix** stellen und den d-check-`matrix`-Scope
   (Spec-Decke + markierte ADR→Slice-Kante) als **bewusste Grenze** in der
   Referenzrichtungs-Adaption deklarieren (C-4; die Planungs-Kanten ADR→{Carveout,
   Welle,Roadmap}/Slice→Roadmap bleiben unbewacht, dokumentiert). **Die
   Historie-Provenance-Enforcement (C-3) ist herausgeschnitten** (§3 Abnahme-Punkt
   3 → Folge-Produkt-Slice).
8. **Gate** — `make gates` + `make adr-check` grün; die neuen `conventions/`-Dateien
   erfüllen die `ids`-Linkpflicht; **keine** `Accepted`-ADR inhaltlich berührt (die
   173 Links via Anker-Block, die eingefrorene ADR-Alt-Pfad-Fundstelle via den
   Tombstone aus Etappe A); unabhängiger Frischkontext-Review.

## 3. Abnahme-Punkte (Maintainer-Entscheid vor/während der Umsetzung)

1. **Fork-Klassifikation** (Vorgehen 4). Nach dem in slice-085 §3.4 am Quelltext
   geschärften Kriterium (*pauschale* Nichtanwendbarkeit der Baseline) ist **keine**
   d-check-Adaption ein echter Fork; die Guard-Härtungs-Adaptionen benennen
   `Ersetzt-Baseline-Regel` auf `grundlagen-durchsetzungsschicht`. **Entscheid:**
   (a) alle als **Adaptionen** mit `Ersetzt-Baseline-Regel` führen (empfohlen,
   quell-belegt), oder (b) die Rest-Kandidaten als Forks behalten.
   → **Entschieden 2026-08-02: (a)** — alle als Adaptionen, **keine Forks**; jeder
   Eintrag benennt eine `Ersetzt-Baseline-Regel`.
2. **Nummern-Kollision** (Vorgehen 6, §2.4). Eine `MR`-Nummer ist in der Baseline
   anders belegt als in d-checks aufgelöster Vorgänger-Adaption. **Entscheid:**
   eigene Nummern behalten (Provenienz) oder an die Baseline angleichen.
   → **Entschieden 2026-08-02: eigene Nummern behalten (Provenienz)** —
   [`MR-003`](../../../../harness/conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh)
   bleibt der historische Bootstrap-Sensor, das Vendoring bleibt
   [`MR-023`](../../../../harness/conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
   (die Nummern tragen die 173 Links; die Baseline-Nummer ist ein
   Vorlagen-Beispiel, keine Regel).
3. **Historie-Provenance-Spec-Bereinigung** (Vorgehen 7 / C-3). Die ADR-/Slice-
   Verweise aus den Spec-§7-Historien entfernen **berührt die kanonische Spec**.
   **Entscheid:** jetzt bereinigen (Baseline-Default sticht) — und ob der interne
   Slice-Bezug nach v5.0.0-Verfahren in die Slice-Closure-Notiz wandert.
   → **Revidiert 2026-08-02 (nach Messung): in einen Folge-Produkt-Slice
   herausgeschnitten.** Die Messung (`exclude-sections` testweise entfernt) zeigte:
   `matrix.exclude-sections` ist **global**, nicht pro Klasse — global entfernen
   bräche **65 Befunde in `Accepted`-immutablen ADRs** (0022–0046, deren
   `Geschichte` Slices als marker-lose Provenance nennt). C-3 braucht ein
   **d-check-Code-Feature** (per-Klasse/-Pfad exclude-sections) **+** die gekoppelte
   Spec-§7-Bereinigung → eigener Produkt-Slice (CR + ADR + Slice am `matrix`-Modul),
   **nicht** diese Doc-Migrations-Etappe C. Die anderen zwei Entscheide bleiben.

## 4. Definition of Done

- [x] `conventions.md` = Index (Adoptions-**Pointer** + Zeile je aktiver Adaption);
  je Adaption eine Datei `harness/conventions/MR-<NNN>-<titel>.md`, aufgelöste unter
  `conventions/done/`. **Form-Delta ggü. Plan:** die Voll-Slug-Anker sind in die
  Index-Zeilen **gefaltet** (kein Sonderblock, §8); vollständiger Template-Abgleich
  gegen `conventions.template.md` (Purpose · Baseline · Adoptierte-Pointer ·
  Adaptions-Block: Inline-Baseline-Aussage + `### Aktive`/`### Aufgelöste` ·
  Zusatzklassen · Modus).
- [x] Neue Pflichtfelder je Eintrag (`Ersetzt-Baseline-Regel`/`Aufgelöst durch`,
  `Status: Accepted`); Forks reklassifiziert → **keine** (Abnahme-Punkt 1a);
  entwertetes `§Anforderungs-Anlege-Prozess`-Duplikat **gelöscht** (Baseline-Deckung
  `modul-03-spec`, §8); Nummern-Kollision: eigene Nummern behalten (Abnahme-Punkt 2).
- [x] Referenzrichtung gegen die 8×8-Matrix gestellt; d-check-`matrix`-Scope als
  bewusste Grenze in der Referenzrichtungs-Adaption deklariert (C-4). (C-3-Enforcement
  herausgeschnitten →
  Folge-Produkt-Slice, Abnahme-Punkt 3.)
- [x] Die `conventions.md#mr-…`-Links lösen auf (gefaltete Index-Anker; **188**
  statt der geplanten 173 gemessen); **keine** `Accepted`-ADR inhaltlich berührt
  (verifiziert über den vollen slice-086-Bereich).
- [x] `make gates` + `make adr-check` grün; umfangreicher Nutzer-Live-Template-
  Abgleich **+** unabhängiger Frischkontext-Review (abnahmereif, 4 LOW eingearbeitet).

## 5. Risiken / offene Punkte

- **Größter, identitäts-berührender Schnitt.** 24 Einträge, 173 Links, 12 immutable
  ADRs — der Anker-Kompatibilitäts-Block (Vorgehen 2) ist der load-bearing Kern;
  ein Fehler dort bricht `anchors` repo-weit.
- **Abnahme-Punkte sind Voraussetzung** für Vorgehen 4/6/7 — ohne Entscheid keine
  Umsetzung dieser Schritte.
- **Kanonische Spec** (Vorgehen 7): der einzige Schritt, der `spec/`-Dateien ändert
  — vorsichtig, review-pflichtig.

## 6. Trigger

Abschluss von slice-085 (Etappe B): die Finding-Liste macht C verbindlich.

## 7. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc/Prozess führt. Berührt *Harness/Konventionen* (der
Konventionsspeicher) und punktuell die *Spec*-Sub-Area (Vorgehen 7) — greenfield,
ohne Brownfield-Spec.

## 8. Closure-Notiz (nach `done/`)

Umgesetzt: `harness/conventions.md` von der Inline-Sammeldatei (895 Zeilen) auf die
v5.0.0-**Default-Form** (Index + Datei je Adaption) gehoben — **8 aktive** Einträge
unter `harness/conventions/`, **15 aufgelöste** unter `harness/conventions/done/`,
jeder mit den neuen Pflichtfeldern (`Ersetzt-Baseline-Regel` bzw. `Aufgelöst durch`,
`Status: Accepted`). Der Index folgt jetzt vollständig `conventions.template.md`:
Purpose · Baseline · **Adoptierte als Pointer** (keine Wiederholung) · Adaptions-Block
(Inline-Baseline-Aussage + `### Aktive`/`### Aufgelöste`-Tabellen) · Zusatzklassen ·
Modus.

**Design-Evolution ggü. dem Plan (Vorgehen 2), getrieben durch den Nutzer-Live-Review:**

- **Anker gefaltet statt Sonderblock.** Der geplante separate
  „Anker-Kompatibilitäts-Block" entfiel; die Voll-Slug-`<a id>`-Anker liegen **in den
  Index-Zeilen** (aktiv wie aufgelöst). `anchors` prüft nur `anchor-missing` (kein
  Duplikat-Check), also tragen die Zeilen-Anker die Links ohne Sonderabschnitt —
  schlanker und ohne redundante Struktur.
- **§Anforderungs-Anlege-Prozess gelöscht statt migriert.** Der Block duplizierte die
  Baseline (`modul-03-spec`: Akzeptanzkriterien-Trias, CR-Verfahren, „ADR schärft
  Spec, nicht Lastenheft") — nach dem „keine Wiederholung"-Prinzip **entfernt**; die
  vier eingefrorenen Verweise darauf auf [`AGENTS.md`
  §5](../../../../AGENTS.md#5-dokumentations-regeln) retargetet (Etappe-A-Muster:
  eingefrorene Doku via Retarget/Tombstone, nicht editieren).
- **`ids`-target auf `harness/conventions/` gezogen** — die MR-Nennungen in den nun
  ausgelagerten Dateien erfüllen die Linkpflicht „im target".

**Zahlen-Korrektur:** der Plan nannte **173** `conventions.md#mr-…`-Links; gemessen
sind es **188** (12 immutable ADRs) — alle erhalten, **keine `Accepted`-ADR**
inhaltlich berührt (über den vollen slice-086-Bereich verifiziert).

**Abnahme-Punkte:** (1a) alle Adaptionen, **keine Forks**; (2) eigene Nummern behalten
(Provenienz); (3) **revidiert** — die Historie-Provenance-Spec-Bereinigung (C-3) ist
in einen **Folge-Produkt-Slice** herausgeschnitten: die Messung zeigte
`matrix.exclude-sections` **global** (nicht pro Klasse), globales Entfernen bräche 65
Befunde in immutablen ADRs → braucht ein d-check-**Code-Feature** (per-Klasse/-Pfad-
`exclude-sections`) **+** die gekoppelte Spec-§7-Bereinigung (CR + ADR + Slice am
`matrix`-Modul). Der C-4-Scope (Spec-Decke + markierte ADR→Slice-Kante, nicht die
volle 8×8-Matrix) ist in der Referenzrichtungs-Adaption deklariert.

**Reviews:** ein **umfangreicher Nutzer-Live-Template-Abgleich** fing mehrere
Template-Abweichungen des ersten Wurfs (Purpose-Wortlaut nicht übernommen · §Adoptierte
forensisch statt Pointer · Anker-Block redundant · §Anforderungs-Duplikat · eine
Fehlklassifikation aktiv/aufgelöst) — alle geheilt; danach **unabhängiger
Frischkontext-Review** (`docs/reviews/2026-08-02-slice-086-etappe-c-review.md`):
abnahmereif, HIGH 0 / MEDIUM 0 / LOW 4 / INFO 2, die 4 LOW eingearbeitet. `make gates`
+ `make adr-check` grün.

**Lehre:** `conventions.md` gegen `conventions.template.md` **abgleichen** ist
Pflichtteil einer Baseline-Hebung — nicht nur der MR-Datei-Split. Die Template-Form
(Pointer statt Wiederholung, `###`-Subsektionen, keine Baseline-Duplikate) ist Prüfstoff
für sich.

**Anschluss:** **Etappe D** (Form-Konformität: Roadmap-5-Abschnitte,
Welle-Closure-Artefakte, `observations.md`, `AGENTS.md`↔Template,
Review-Report-Kopffelder, reviewer.md-Currency) offen; der C-3-Folge-Produkt-Slice
separat.
