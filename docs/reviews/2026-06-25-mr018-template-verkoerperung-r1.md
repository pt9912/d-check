# Review-Report — MR-018 (d-check verkörpert als Producer-/Self-Hoster keine Templates), R1

## Kopf-Metadaten

- **Datum:** 2026-06-25
- **Reviewer:** unabhängiger adversarialer Reviewer (d-check)
- **Gegenstand:** unkommittierter `git diff` — `AGENTS.md` §1 +
  `harness/conventions.md` (neuer Block MR-018 + Schärfung der §Adoptierte-
  Lese-Form-Bullet). Kein Code, kein Skript, keine Gate-Config-Änderung.
- **Kontext:** Neu gefasste MR-018. Die frühere Fassung (behauptetes
  „Kanon-Loch", „0-Treffer"-Verifikation) wurde entfernt; geprüft wird die
  finale Producer-/Self-Hoster-Brücken-Fassung.
- **Pflichtlektüre gelesen:** `.harness/skills/reviewer.md` (Rolle/Kategorien/
  Schema/Negativbefund-Pflicht), `AGENTS.md` §3 (Hard Rules),
  `harness/conventions.md` MR-014 · MR-017 · neuer MR-018 + §Adoptierte.
- **Quellen-Verifikation (gelesen):** Live-Kurs `lab/templates/README.md`;
  gepinnte v1.4.0-Baseline `…/templates/README.md`; Schwester-Repo
  `bedrock-eu-guard` (`git ls-files '*.template.md'`); d-check selbst
  (`git ls-files '*.template.md'` = leer); Slugify-Engine
  `internal/hexagon/core/rules/anchors.go`; `.d-check.yml`.

## Findings

### F-1 — Prosa: „Kurs stützt diese Producer-Lesart bereits live" überschreibt für eine Atemzuglänge die Skelett-Grenze

- **kategorie:** LOW
- **quelle:** Maintainability (Doku-Drift ggü. Live-Quelle
  `lab/templates/README.md` §Self-Hosting-/Producer-Fall)
- **pfad:** `harness/conventions.md:551`
- **befund:** Der Einstiegssatz der Begründung sagt pauschal „Der **Kurs
  stützt diese Producer-Lesart bereits live**". Der zitierte Live-Abschnitt
  §Self-Hosting-/Producer-Fall stützt die Producer-Lesart aber ausschließlich
  für die `harness.mk`-Konsumenten-Integration und sagt im selben Absatz
  ausdrücklich, dass die **Dokument-Skelette adoptiert werden** („adoptiert
  werden die Dokument-Skelette plus `Makefile` und `.d-check.yml`, nicht
  `harness.mk`"). Der Pauschalsatz greift damit für einen Satz weiter als die
  Quelle, bevor die Folgesätze ihn selbst korrekt einengen.
- **Mildernd (kein REFUTE, aber relevant):** Dieselbe Begründung korrigiert
  sich zwei Sätze später explizit: „nimmt es von der `harness.mk`-Adoption
  aus" und „zieht die Producer-Logik zugleich auf die wiederkehrenden
  **Dokument**-Skelette weiter — den Schritt, den der Kurs auch live (noch)
  nicht ausspricht." Das überdehnt also nicht im Saldo; der Eröffnungssatz ist
  nur lokal unscharf.
- **verifizierbar:** nein — kein Gate misst Prosa-Treue gegen eine externe,
  ungepinnte Quelle; reiner Reviewer-Befund (Klasse §Ein-vs-wiederkehrende,
  reviewer.md: „prüft Existenz, nicht ob Vorhandenes korrekt beschrieben ist").

### F-2 — AGENTS.md nennt „(ADR, Slice)", MR-018 nennt fünf Skelette

- **kategorie:** INFO
- **quelle:** Maintainability (Pointer-/Kanon-Konsistenz, AGENTS.md §1 →
  MR-018)
- **pfad:** `AGENTS.md:32`
- **befund:** Die operative AGENTS.md-Kurzform schreibt „wiederkehrende
  Artefakte (ADR, Slice) entstehen nativ im Haus-Stil", während der kanonische
  MR-018-Block fünf wiederkehrende Skelette führt (ADR, Slice, Welle,
  Carveout, Review-Report). Das ist ein illustrativer Teilausschnitt mit
  Pointer auf MR-018 (kanonisch vollständig), keine widersprüchliche Aussage —
  d-check autoriert Review-Reports und Carveouts nachweislich nativ
  (`docs/reviews/*`, `docs/plan/carveouts/`), ohne co-located `*.template.md`.
  Bei künftigem Edit könnte die verkürzte Klammer als „nur diese zwei" fehl-
  gelesen werden.
- **verifizierbar:** nein — Klassen-Konsistenz-Befund, kein Gate-Bindepunkt
  (Rollen-Trennung operativ↔kanonisch ist bewusst, MR-018 ist die Quelle).

## Negativbefunde (geprüft, ohne blockierenden Befund)

- **Faktencheck Live-Kurs §Self-Hosting-/Producer-Fall:** geprüft —
  existiert in `lab/templates/README.md`, nennt „das Tool-Repo selbst
  (d-check), das seinen Doku-Gate via `make doc-check` direkt dogfooded" und
  nimmt es von der `harness.mk`-Adoption aus. MR-018s Zitat ist wortgetreu.
  Kein Befund.
- **Faktencheck Live-Kurs §Ein-vs-wiederkehrende:** geprüft — sagt für die
  wiederkehrenden Skelette (ADR, Slice, Welle, Carveout, Review-Report)
  wörtlich „als `.template.md` **co-located** im Repo behalten; jede neue
  Instanz wird daneben kopiert"; keine Producer-Ausnahme. MR-018s
  Charakterisierung („d-check weicht bewusst von dieser Regel ab") ist damit
  die ehrliche Lesart — die Baseline trägt für die **Dokument**-Skelette keine
  Producer-Carve-out, also ist „Abweichung" korrekt und nicht überzogen
  (adversariale Frage 1 widerlegt). Kein Befund.
- **Faktencheck gepinnte v1.4.0-Baseline (`…/templates/README.md`):** geprüft
  — §Self-Hosting-/Producer-Fall fehlt dort tatsächlich (die §Gate-Baseline
  trägt keinen solchen Unterabschnitt); §Ein-vs-wiederkehrende ist bereits
  vorhanden. MR-018s Aussagen „post-v1.4.0" und „v1.4.0 trägt diesen Abschnitt
  noch nicht" sind belegt. Kein Befund.
- **Faktencheck Schwester-Repo `bedrock-eu-guard`:** geprüft —
  `git ls-files '*.template.md'` liefert genau die fünf wiederkehrenden
  Skelette co-located (`adr/NNNN-titel`, `slice`, `welle`, `carveout`,
  `review-report`). MR-018s „verkörpert die fünf wiederkehrenden Skelette
  co-located" ist exakt belegt. Kein Befund.
- **Faktencheck d-check selbst:** geprüft — `git ls-files '*.template.md'`
  und ein cache-ausgenommener `find` liefern **keine** co-located Templates;
  Review-Reports/Carveouts entstehen nativ. Der zentrale Tatsachenkern von
  MR-018 hält. Kein Befund.
- **Adversarial 2 (Selbstwiderspruch Producer-Argument):** geprüft — d-check
  autoriert ADRs/Slices/Reviews wiederkehrend, tut dies aber nachweislich
  „nativ im Haus-Stil" aus dem Bestand (Verweis MR-014, der den gelebten
  Haus-Stil ggü. Baseline-Template deklariert). Das Producer-Argument („nativ
  statt aus Skelett") ist intern konsistent mit dem realen Repo-Zustand und
  mit MR-014. Kein Befund.
- **Adversarial — Auflösungs-Trigger-Konsistenz:** geprüft — der Trigger
  („Re-Evaluation beim nächsten Baseline-Bump; aufgelöst, falls der Kurs die
  Producer-Lesart auch für die wiederkehrenden Dokument-Skelette übernimmt")
  deckt sich mit dem in der Begründung benannten Bridge-Charakter (der Kurs
  spricht den Dokument-Skelett-Schritt live noch nicht aus). Brücke ↔ Trigger
  sind konsistent. Kein Befund.
- **Append-only (MR-014/MR-017):** geprüft — der Diff fügt MR-018 als neuen
  Block an und ändert die §Adoptierte-Lese-Form-Bullet; die Block-Körper von
  MR-014 (Z. 366–401) und MR-017 (Z. 493–525) bleiben unangetastet, beide
  werden nur per Link referenziert. Kein inhaltliches Umschreiben. Kein Befund.
- **ID-Hygiene / nacktes `MR-003`:** geprüft — der Diff führt **kein** neues
  `MR-003`-Vorkommen ein; bestehende nackte/Link-`MR-003`-Token
  (`conventions.md` Z. 59/184/189) sind vorbestehend und nicht Teil dieser
  Änderung. Kein neuer Befund durch den Diff.
- **Anker-Auflösung (MR-018 + Cross-Links):** geprüft — die Slugify-Engine
  (`anchors.go`: ToLower, Markup-Drop, em-dash/`/` verworfen, Spaces→`-`)
  erzeugt für die Überschrift exakt
  `mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates`
  (inkl. Doppel-Bindestrich um den em-dash); beide Link-Anker (AGENTS.md→
  conventions.md und der In-File-Selbstlink) treffen diesen Slug zeichengenau.
  Der `anchors`-Gate würde grün laufen. Kein Befund.
- **Rollen-Trennung (kanonisch↔operativ):** geprüft — conventions.md trägt
  die kanonische Adaptions-Deklaration (MR-018-Block), AGENTS.md §1 trägt die
  operative Kurzregel mit Pointer auf MR-018; keine kanonische Setzung
  dupliziert in AGENTS.md. Sauber (abgesehen vom INFO-Teilausschnitt F-2).
  Kein blockierender Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 (F-1) |
| INFO | 1 (F-2) |

## Verdikt

**Nicht blockierend — mergebar.** Der Tatsachenkern von MR-018 ist gegen alle
vier Quellen (Live-Kurs, gepinnte v1.4.0, `bedrock-eu-guard`, d-check selbst)
belegt; die frühere „Kanon-Loch"-Überdehnung ist entfernt. Die finale Fassung
deklariert ihren Brücken-/Eigenanteil transparent („den Schritt, den der Kurs
live noch nicht ausspricht"), Anker lösen auf, Append-only ist gewahrt, die
Rollen-Trennung ist sauber. Kein HIGH/MEDIUM. F-1 (LOW) ist ein lokal
unscharfer Eröffnungssatz, der sich im selben Bullet selbst einengt; F-2
(INFO) ein verkürzter, gepointerter Teilausschnitt in der operativen
Kurzform. Beide vor oder ohne Merge optional zu glätten, keiner blockiert.
