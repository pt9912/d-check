# Re-Review-Report (R2): slice-083 (Regelwerk-Migration v1.4.0 → v5.0.0) — 2026-08-01

**Review-Art:** Bestätigende Re-Review (R2) — Verifikation der R1-Nachbesserung
gegen die kanonischen Quellen, unabhängiger Frischkontext. Prüft, ob die R1-Befunde
quellen-gegengeprüft geheilt sind und ob die Fixes NEUE Probleme einführten.

**Gegenstand:** `docs/plan/planning/open/slice-083-regelwerk-v500-migration-analyse.md` (uncommittet).

**Vorlauf:** R1-Report `docs/reviews/2026-08-01-slice-083-v500-migration-review.md`
(2 MEDIUM F-1/F-2, 3 LOW F-3/F-4/F-5, 2 INFO F-6/F-7).

**Skill:** `.harness/skills/reviewer.md` v1.2.0 · **Modell:** claude-opus-4-8 · **Datum:** 2026-08-01

**Eingangs-Kontext (geprüfte Verträge/Quellen):**

- Aktuelle slice-083 (ganz).
- Anker-Engine: `internal/hexagon/core/rules/anchors.go` (DC-FA-ANCH-001).
- Doku-Gate: `.d-check.yml` (scan-Skopus, `modules`, `ids`-Ziele).
- d-check-Ist: `harness/conventions.md`, `AGENTS.md`, `harness/README.md`,
  die zwölf Accepted-immutablen ADRs unter `docs/plan/adr/`.
- Baseline v5.0.0: Schwester-Repo `ai-harness-course`, Tag `v5.0.0`
  (`lab/templates/AGENTS.template.md`, `lab/regelwerk/grundlagen-source-precedence.md`,
  `lab/regelwerk/grundlagen-*`). Das v5.0.0-Scratchpad-Bundle war nicht mehr präsent;
  gegengeprüft wurde durabel gegen den `v5.0.0`-Tag des Schwester-Repos.

**Methode:** jede R1-Heilung gegen die Quelle nachgemessen (Repo-Greps mit
Zeilenzählung, Anker-Engine am Code gelesen, git-Tag-Diffs). REFUTE nur mit Zitat.

---

## Verdikt je R1-Befund

### F-1 (MEDIUM, R1) — retired/entfallene Quellzeiger ohne Etappe-Schritt → **GEHEILT**

Der Fix sitzt in Etappe A, Schritt 4 (Slice-Zeilen 349–356): die **entfallenen**
Quellzeiger werden ausdrücklich **umgeschrieben statt retargetet**. Gegengeprüft:

- Die zwei benannten Fundstellen existieren wie beschrieben:
  `harness/conventions.md:37` trägt den Zeiger `kurs/de/grundlagen/konventionen.md`
  (Quelldatei entfällt in v5.0.0) und `:39` die `agents-regelwerk.md`-Raw-URL
  (v2.0.0 retired); `harness/README.md:42` führt dieselbe `agents-regelwerk.md`-URL
  in der §Guides-Tabelle als „adoptiertes Betriebsregelwerk".
- A4 benennt beide Träger-Dateien korrekt („in `conventions.md` §Adoptierte Quellen
  und `harness/README.md` §Guides … umgeschrieben bzw. entfernt, **nicht** auf tote
  v5.0.0-Ziele umgehängt"). Die falsche Operation aus R1 (A4s pauschales „auf den
  neuen Pfad ziehen") ist damit für die retireten/entfallenen Assets ausgenommen.

Die Heilung deckt genau die von R1 benannten Fundstellen ab; der „umschreiben
statt retargeten"-Schritt ist vollständig. Kein Restbefund.

### F-2 (MEDIUM, grenzt an HIGH, R1) — `conventions.md`-Split bricht `MR`-Anker-Links → **TEILWEISE geheilt / MEDIUM bleibt offen**

Der Fix hat drei Teilaussagen; zwei sind belegt geheilt, die dritte (der Kern) ist es nicht.

**(b) Zählung — GEHEILT.** Die R1-Zahl (52/16) ist auf **173 Links / 57 Dateien /
12 immutable ADRs** korrigiert und **exakt** (nicht nur größenordnungs-):
`grep -roE 'conventions\.md#mr-[0-9]{3}'` über das Repo liefert **173** Vorkommen in
**57** Dateien; davon sind **12** Accepted-ADRs (0019, 0021, 0022, 0023, 0024, 0026,
0027, 0028, 0029, 0030, 0031, 0046 — je `**Status:** Accepted`). Die 173 zählt
korrekt den Gesamtbestand (inkl. `done/`-Slices, Reviews und der slice-083 selbst);
R1s 52/16 war die editierbar-relevante Teilmenge. Beide sind konsistent.

**(a) Anker-Semantik — im Prinzip tragfähig, aber im Slice unterspezifiziert.**
`anchors.go` bildet die gültige Anker-Menge als `AnchorSet` = Heading-Slugs ∪
`htmlAnchors` (Zeilen 155–161); `htmlAnchors` liest `id`-Attribute an **beliebigen**
Elementen **wörtlich** in die Menge (`set[v] = true`, Zeilen 120–135). Ein Link löst
auf, wenn `slugs[a.frag]` **exakt** trifft (Zeile 219 — Map-Lookup, **kein**
Präfix-/Teilstring-Abgleich). Ein `<a id="…">` im Index kann einen Link also decken —
**aber nur, wenn der `id`-Wert das Fragment wörtlich reproduziert.** Gemessen: **alle
173** Links tragen den **vollen Heading-Slug** (z. B.
`conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding`); **null** Links
tragen die Kurzform `#mr-NNN`. Der im Slice notierte Anker `<a id="mr-NNN">`
(Zeilen 210, 393–394) ist damit als konkrete Anweisung **falsch verkürzt**: ein
literaler `<a id="mr-007">` deckt **keinen einzigen** der 173 Links; der Anker müsste
`<a id="mr-007--auflösung-von-mr-003-doc-check-als-dogfooding">` (voller Slug je
Titel) lauten. Für sich genommen LOW; es fällt aber mit (c) zusammen.

**(c) C8-Zusage „keine Accepted-ADR berührt / 173 Links gültig" — NICHT erreichbar
(MEDIUM, siehe Neu-Befund N-1).** Der Fix trägt den Anker „je Index-Zeile" (C2,
Zeile 393), und der Index führt laut Slice **nur eine Zeile je *aktiver* Adaption**
(Zeilen 199, 390); aufgelöste Einträge wandern per `git mv` nach `conventions/done/`
(Zeile 390), als `entfällt`/`verschmilzt` markierte werden in C5 **gestrichen**
(Zeilen 402–405). Damit verlässt der Anker der aufgelösten/gestrichenen Einträge
`conventions.md`. Genau auf solche Einträge zeigen die immutablen ADRs: von den
zwölf verlinken **zehn** (alle außer 0021/0022) auf `mr-003` (aufgelöst→`done/`),
`mr-007` (historisch→`done/`), `mr-017` (`entfällt`), `mr-019`/`mr-022`
(`verschmilzt`) — alles Einträge, die den Index nach C verlassen. Diese ADRs sind
nach ADR-0016/ADR-0024 nicht editierbar. Der Index-Anker-Mechanismus deckt also die
**aktiven** Einträge, nicht aber die Mehrheit der Links; C8s Zusage bleibt so
unerreichbar. → Re-eröffnet als **N-1 (MEDIUM)**.

### F-3 (LOW, R1) — Bruch-6 „acht" → „sechs" → **GEHEILT**

Slice-Zeile 138 lautet jetzt „in **sechs neue** `grundlagen-*`-Dateien aufgesplittet
(grundlagen damit 3 → 8)". Quellentreu (v5.0.0 `lab/regelwerk/` trägt acht
`grundlagen-*`, davon sechs aus dem `konventionen.md`-Split). Zahl-Inkonsistenz behoben.

### F-4 (LOW, R1) — DoD „vier Brüche" → „sechs" → **GEHEILT**

DoD-Zeile 471 lautet jetzt „Die **sechs** Brüche benannt"; deckungsgleich mit der
§2.2-Überschrift (Zeile 98) und dem sechs-elementigen Bestand. Kein Rest-Drift.

### F-5 (LOW, R1) — `observations.md` fehlte in „fehlende Artefakte" → **GEHEILT**

Das Beobachtungs-Register ist an drei Stellen ergänzt: §2.3 „Fehlende Artefakte"
(Zeile 176, „zudem fehlt das **Beobachtungs-Register** (`observations.md`,
v5.0.0-Standard-Artefakt)"), Etappe-D-Übersicht (Zeile 327) und Etappe D Schritt 4
(Zeile 431, „das **Beobachtungs-Register** `observations.md` (v5.0.0-Vorlage)").

### F-6 (INFO, R1) — Fork-Grenze: weite Template- vs. engere source-precedence-Fassung → **GEHEILT (aufgeführt)**

Der Slice benennt die Divergenz jetzt: §2.3 (Zeile 218) „wer keine benannte Regel
ersetzt, ist ein Fork, keine Adaption (**kanonisch enger:
`grundlagen-source-precedence.md`**)" und deklariert die verbindliche
Fork-Zuordnung als Abnahme-/Etappe-C-Detail (§2.4 Zeilen 262–268, §2.8 Zeilen
461–464). Gegen die Quelle bestätigt: `ai-harness-course` `v5.0.0`
`lab/regelwerk/grundlagen-source-precedence.md` Zeilen 79–88 fasst Fork enger („eine
`MR-<NNN>`, die die Baseline **pauschal für nicht anwendbar erklärt**, statt eine
benannte Regel zu ersetzen") und MR-000 als „weder Fork noch Adaption". R1s Forderung
(die zwei Fassungen beim Abnahme-Entscheid benennen) ist erfüllt. Rest-Nuance (kein
Finding): der Slice hängt die *weite* Formulierung sprachlich an die *enge* Quelle,
statt den Unterschied auszuformulieren — surfacet die Divergenz aber ausreichend.

### F-7 (INFO, R1) — „23 Adaptionen" zählt `MR-000` mit → **GEHEILT (an der Definitionsstelle)**

§2.3 (Zeilen 205–206) disaggregiert jetzt: „die 23 Einträge (die Adoptions-Erklärung
+ 22 Adaptionen/Forks)". An der begrifflichen Stelle korrekt. Rest (kein Finding):
die Kurzform „die 23 Adaptionen" steht noch in der Etappen-Tabelle (Zeile 326) und
C1 (Zeile 386); die Definition in §2.3/§2.4 stellt das richtig.

---

## Neue Findings (durch die Nachbesserung eingeführt/verbliebene Lücke)

### N-1 (MEDIUM) — Der F-2-Anker-Fix deckt nur aktive Einträge; die Links der aufgelösten/gestrichenen Einträge (u. a. in zehn immutablen ADRs) brechen weiterhin

- **kategorie:** MEDIUM
- **quelle:** §2.3 (Index-Anker) / §2.7 Etappe C (Schritte 2 + 5 + 8) / DC-FA-ANCH-001 / ADR-0016 + ADR-0024 (Immutabilität)
- **pfad:** `docs/plan/planning/open/slice-083-...md:390` (C2 „je aktiver Adaption"),
  `:393-394` (Anker-Erhalt für 173), `:402-405` (C5 streicht `entfällt`/`verschmilzt`),
  `:414-416` (C8-Zusage); Belege `internal/hexagon/core/rules/anchors.go:219`
  (Exakt-Match), `.d-check.yml:13-16` (scan-Skopus + `anchors` aktiv),
  `docs/plan/adr/0026-...md:16`, `0030-...md:15+17`, `0046-...md:20+50+88` (immutable
  Links auf gestrichene/aufgelöste `MR`)
- **befund:** C2 verspricht, „jede Index-Zeile" trage einen Anker, damit **alle 173**
  Links auflösen; der Index führt aber laut C2 nur „eine Zeile je **aktiver**
  Adaption", und C5 **streicht** die `entfällt`-/`verschmilzt`-Einträge, während
  aufgelöste per `git mv` nach `conventions/done/` gehen. Die Anker der nicht-aktiven
  Einträge verlassen damit `conventions.md`. Von den 173 Links zeigen **mindestens
  111** auf solche Einträge (u. a. `mr-007` mit 50, `mr-003` mit 13, `mr-022` mit 9,
  `mr-019` mit 8, `mr-017` mit 7, `mr-001` mit 5); **zehn der zwölf immutablen ADRs**
  (alle außer 0021/0022, die nur auf das aktive `mr-006` zeigen) verlinken auf
  `mr-003`/`mr-007`/`mr-017`/`mr-019`/`mr-022`. Nach C melden `anchors` für diese
  Links `anchor-missing`; die ADRs sind nach ADR-0016/ADR-0024 nicht editierbar, das
  A6-Tombstone deckt laut Slice nur die drei v1.4.0-**Pfad**-Verweise, nicht die
  `#mr-NNN`-Anker. C8s Zusage „keine Accepted-ADR berührt / 173 Links via Index-Anker
  gültig" (Zeilen 414–416) ist so nicht erreichbar. Verschärfend: der Anker müsste
  ohnehin den **vollen** Titel-Slug tragen (`<a id="mr-007--auflösung-…">`), nicht die
  im Slice notierte Kurzform `<a id="mr-NNN">` — alle 173 Links sind voll-slug, null
  kurz. Die Schritt-Reihenfolge C2 (Anker-Erhalt für alle 173) vs. C5 (Streichung)
  ist damit in sich widersprüchlich.
- **verifizierbar:** ja — nach Ausführung von Etappe C meldet `make gates`
  (`anchors`-Modul) `anchor-missing` für die auf nicht-aktive Einträge zeigenden
  Links (≥111, davon in zehn immutablen ADRs); ohne benannte Tombstone-/Anker-Erhalt-
  Kampagne **auch** für die nicht-aktiven Einträge ist C8 nicht grün.

---

## Negativbefunde (je Prüfachse — geprüft)

**Achse F-1 (entfallene Quellzeiger).** A4s „umschreiben statt retargeten" trifft die
zwei realen Fundstellen (`conventions.md:37/:39`, `README.md:42`) exakt; kein
Restbefund. Geprüft, ohne Befund.

**Achse F-2b (Zählung).** 173/57/12 am Repo-Grep nachgemessen und **exakt** bestätigt
(alle zwölf ADR-Ziele `Status: Accepted`). Geprüft, ohne Befund.

**Achse F-2a (Anker-Engine).** `htmlAnchors` honoriert `<a id>` wörtlich, Exakt-Match
in `CheckAnchors` (anchors.go:120–135, 219): der Mechanismus **kann** tragen — die
Unterspezifikation (Kurz- vs. Voll-Slug) und die Nicht-Abdeckung nicht-aktiver
Einträge sind in **N-1** gebündelt.

**Achse F-2c / C8 (Erreichbarkeit).** Der scan-Skopus (`.d-check.yml:13-16`,
`roots: ["."]`, nur `.harness/{baseline,cache}/**` ignoriert; `anchors` in `modules`)
bestätigt, dass die immutablen ADRs im Prüf-Skopus liegen — der Bruch ist gate-real.
**Befund N-1.**

**Achse F-3/F-4/F-5 (LOW).** Alle drei Textstellen nachgezogen und quellentreu
(Zeilen 138, 471, 176/327/431). Geprüft, ohne Befund.

**Achse 4 — AGENTS.md gegen `AGENTS.template.md`.** Gegen `ai-harness-course`
`v5.0.0` `lab/templates/AGENTS.template.md` bestätigt: **Struktur deckungsgleich**
(Template und `AGENTS.md` tragen je §1–§6 und §3.1–§3.6 mit identischer Nummerierung).
**§1-Drift bestätigt:** das Template führt in §1 den Kanon-Zeiger „Baseline-Regelwerk
`modul-09-implementierung.md` §Ziel-Form: AGENTS.md" und das
`{regelwerk,templates}`-Layout unter `.harness/baseline/<tag>/`; `AGENTS.md` §1 trägt
weder den Modul-9-Zeiger noch dieses Layout. **§3.1 bestätigt:** Template „Docker-only"
vs. `AGENTS.md:64` „Docker/make-only". **§3.4 bestätigt:** Template „Architektur ist
sprach- und meilensteinfrei" vs. `AGENTS.md:108` „Architektur sprach-/meilensteinfrei;
**Spec-Straten nie abwärts**" (angehängter Zusatz). Rest-Nuance (kein Finding): der
Slice benennt die Template-Sektion als „Modul 9 §**AGENTS.md-Regeln**"; kanonisch
heißt sie „§Ziel-Form: AGENTS.md" — reine Namens-Ungenauigkeit, ohne Migrations-Wirkung.
Geprüft, ohne blockierenden Befund.

**Achse F-6/F-7 (INFO).** Beide an der Quelle gegengeprüft und ausreichend gesurfacet
(source-precedence Zeilen 79–88; §2.3 Zeilen 205–206/218). Geprüft, ohne Befund.

**Achse 5 — neue Probleme / Etappe-C-Kohärenz.** Ein neues Problem gefunden: die C2/C5-
Reihenfolge widerspricht sich am Anker-Erhalt (→ **N-1**). Die übrige Reihenfolge
(splitten → Felder → Forks → Entfall → Nummern-Kollision → entpinnen) ist konsistent;
die F-1-, F-3/4/5-, F-6/7- und AGENTS.md-Ergänzungen führen **keine** weiteren
Widersprüche oder Doppelungen ein. Geprüft; ein Befund (N-1).

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | N-1 (Re-Eröffnung von R1-F-2c) |
| LOW | 0 | — |
| INFO | 0 | — |

R1-Heilungsstand: F-1 GEHEILT · F-2 TEILWEISE (b/a-Prinzip geheilt, c offen → N-1) ·
F-3 GEHEILT · F-4 GEHEILT · F-5 GEHEILT · F-6 GEHEILT · F-7 GEHEILT.

---

## Gesamt-Verdikt

**Noch nicht abnahmereif** — ein offener MEDIUM (N-1). Sechs der sieben R1-Befunde
sind quellen-gegengeprüft geheilt; die Nachbesserung ist substanziell und führte —
außer N-1 — keine neuen Widersprüche ein. Blockierend bleibt allein, dass der zur
Heilung von F-2 ergänzte **Index-Anker-Mechanismus** die von F-2 adressierte Kern-Not
(die Links der immutablen ADRs ohne ADR-Edit erhalten) **nicht** trägt: er deckt nur
die im Index verbleibenden **aktiven** Einträge, während die immutablen ADRs
mehrheitlich auf aufgelöste/gestrichene Einträge (`mr-003/007/017/019/022`) zeigen,
deren Anker den Index nach Etappe C verlassen. C8s Zusage „keine Accepted-ADR berührt /
173 Links gültig" ist damit weiterhin unerreichbar. N-1 ist — wie R1-F-2 — durch
Ergänzung der Etappe-C-Beschreibung heilbar (Anker-Erhalt/Tombstone auch für die
verlinkten nicht-aktiven Einträge, Anker als **voller** Titel-Slug), ohne den Schnitt
umzukehren. Der Abnahme-Entscheid sollte N-1 explizit adressieren.
