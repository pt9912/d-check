# Review-Report — slice-056 (Modul `commits`), R1 Doc-first / Cross-Surface

## Kopf-Metadaten

- **Gegenstand:** slice-056 — neues opt-in Modul `commits` (Commit-Message-Traceability,
  mechanisiert `tools/trace-check.sh` über den VCS-Port). Fokus dieses Reviews:
  **Doc-first-Straten + Cross-Surface-Konsistenz** (NICHT der Go-Code — separater Reviewer).
- **Diff-Basis:** `git diff HEAD` + neue Dateien (`docs/plan/adr/0027-…`,
  `docs/plan/planning/in-progress/slice-056-…`, `rules/commits*.go`, `cli_commits_test.go`),
  Working Tree gegen `HEAD` = `9a40b95`.
- **Reviewer-Rolle:** unabhängig, adversarial. Baseline: `.harness/skills/reviewer.md`
  (v1.2.0), `AGENTS.md` §3 (Hard Rules), `harness/conventions.md` (MR-006 Referenzrichtung,
  MR-014 Doc-Struktur, §Anforderungs-Anlege-Prozess).
- **Kern-Artefakte geprüft:** `spec/lastenheft.md`, `spec/spezifikation.md`,
  `spec/architecture.md`, `docs/plan/adr/0027-commits-traceability-modul.md`,
  `docs/plan/adr/README.md`, `docs/plan/planning/in-progress/{slice-056,roadmap}.md`,
  `.d-check.yml`, `.githooks/commit-msg`, `Makefile`, `AGENTS.md`, `harness/README.md`,
  `CHANGELOG.md`, `version.md`, `README.md`, `docs/user/benutzerhandbuch.md`,
  Config-Surface (`config_template.go`, `suggest.go`, `print_mk.go`, `config.go`,
  `finding.go`).
- **Verifikations-Mittel:** `git diff`, Dateilektüre, gezielte `grep`-Proben. DoD-Abhakung
  bewusst **nicht** geprüft (nicht Reviewer-Rolle).

---

## Findings

### F-1 — MEDIUM — Lastenheft `DC-FA-CLI-010`-Körper + Akzeptanzkriterien nennen weiter „sieben" Targets und lassen `doc-commits` aus

- **kategorie:** MEDIUM
- **quelle:** `DC-FA-CLI-010` (Lastenheft, Source-Precedence Rang 1),
  `harness/conventions.md` §Anforderungs-Anlege-Prozess (geänderte `DC-*` → aktuelle
  Akzeptanzkriterien)
- **pfad:** `spec/lastenheft.md:377-378` (Beschreibung), `:402` (Happy-Path-AC),
  `:403` (Boundary-AC), `:406` (Out-of-Scope)
- **befund:** Diese Slice **erweitert** `DC-FA-CLI-010` von 7 auf 8 `--print-mk`-Targets —
  so deklariert es die Slice selbst in `spec/lastenheft.md:1384` (§7-Historie:
  „`DC-FA-CLI-010`-Erweiterung (7→8 Targets): `--print-mk` trägt ein `doc-commits`-Target"),
  ebenso `spec/spezifikation.md:335` („**Acht** `.PHONY`-Targets" + `doc-commits`-Bullet)
  und der tatsächlich ausgelieferte Fragment-Generator (`print_mk.go` emittiert 8 `doc-*`-
  Targets inkl. `doc-commits`). Der **Anforderungs-Körper** von `DC-FA-CLI-010` wurde jedoch
  nicht mitgezogen: die Beschreibung nennt „sowie **sieben** `##`-annotierte Targets" und
  zählt sie ohne `doc-commits` auf; die **Happy-Path-AC** (`:402`) enumeriert exakt dieselben
  sieben (`doc-check`, `doc-trace`, `doc-complete`, `doc-doctor`, `doc-repair`,
  `doc-immutable`, `doc-help`); die Out-of-Scope-Zeile (`:406`) erklärt „weitere Targets
  **jenseits der gelisteten sieben** … Konsumenten komponieren weitere `gates` selbst" —
  klassifiziert `doc-commits` damit **aktiv als out-of-scope/konsumenten-komponiert**,
  im direkten Widerspruch zur Spezifikation und zum ausgelieferten Fragment. Der Rang-1-
  Vertrag widerspricht damit seiner eigenen §7-Änderungshistorie.
- **Failure-Szenario:** Ein Verifizierer/Umsetzer schreibt/prüft den `--print-mk`-
  Akzeptanztest strikt gegen die Happy-Path-AC von `DC-FA-CLI-010` — die exakt sieben
  Targets aufzählt und `doc-commits` (das Headline-Distributions-Target dieser Slice)
  **nicht** fordert. Ein Reviewer, der dem bindenden Lastenheft folgt, liest aus der
  Out-of-Scope-Zeile, `doc-commits` gehöre nicht zum `DC-FA-CLI-010`-Vertrag. Die verteilte
  Commit-Traceability („verteilen statt kopieren", `MR-007`-Kern dieser Slice) fehlt damit
  im bindenden Vertrag, obwohl sie im Produkt existiert.
- **verifizierbar:** ja — `docker run --rm ghcr.io/pt9912/d-check:v0.35.0 --print-mk | grep -c '^doc-.*:'`
  ⇒ 8 (inkl. `doc-commits`), gegen das „sieben"/7-Enumerationen in `spec/lastenheft.md:377,402,406`;
  zugleich Selbst-Widerspruch Lastenheft-Körper ⟷ Lastenheft-§7-Historie (`:1384`).

### F-2 — LOW — Benutzerhandbuch `--print-mk`-How-to nennt „sechs" Targets, lässt `doc-immutable` UND `doc-commits` aus

- **kategorie:** LOW (Doku-Drift Prosa-Liste; Reviewer-Skill §LOW „veraltete Beispiele")
- **quelle:** `DC-FA-CLI-010`, Cross-Surface-Ehrlichkeit (docs/user, Rang 6)
- **pfad:** `docs/user/benutzerhandbuch.md:622-624`
- **befund:** Der Fließtext im Abschnitt „Makefile-Fragment ausgeben" beschreibt
  „`TRACE_FLAGS` und **sechs** `##`-annotierte Targets (`doc-check`, `doc-trace`,
  `doc-complete`, `doc-doctor`, `doc-repair`, `doc-help`)". Der Generator liefert real
  **acht** Targets; die Aufzählung lässt `doc-immutable` (bereits seit slice-053 fehlend)
  **und** das neue `doc-commits` aus. Slice-056 hat die unmittelbar angrenzende Zeile
  (`DCHECK_IMAGE …:v0.34.0`→`:v0.35.0`, `:625`) im selben Fenced-Block angefasst und
  direkt darunter (§5, `:781-786`) eine `doc-commits`-How-to ergänzt, die Prosa-Zählung
  hier aber nicht mitgezogen.
- **Failure-Szenario:** Ein Konsument liest die `--print-mk`-How-to, entnimmt der
  vollständig wirkenden Aufzählung „sechs Targets (…, `doc-help`)", dass `doc-commits`
  nicht Teil des Fragments ist, und verdrahtet die verteilte Commit-Traceability nicht —
  obwohl `d-check.mk` das Target enthält (auffindbar nur über `make doc-help` / Lektüre
  des generierten Fragments).
- **verifizierbar:** ja — `--print-mk`-Ausgabe (8 `doc-*`-Targets) gegen die „sechs"-Prosa
  in `docs/user/benutzerhandbuch.md:622`.

### F-3 — LOW — Benutzerhandbuch „Versionen und Tags" nennt weiter `:v0.34.0` als empfohlene feste Version (Pin-Sweep-Lücke, vom `versions`-Gate nicht gefangen)

- **kategorie:** LOW (latente Pin-Drift / veraltetes Beispiel)
- **quelle:** `DC-FA-VER-001` (Versions-Pin-Konsistenz) — Blind-Spot; Pin-Sweep-Disziplin
- **pfad:** `docs/user/benutzerhandbuch.md:74`
- **befund:** Im Abschnitt „### Versionen und Tags" illustriert die Aufzählungszeile
  `:v0.34.0` die feste-Version-Variante („empfohlen für reproduzierbare Läufe"). Der
  `docker pull`-Befehl acht Zeilen darüber (`:66`) wurde auf `:v0.35.0` gesweept, ebenso
  alle übrigen `ghcr.io/pt9912/d-check:v…`-Pins im Handbuch. Dieses Beispiel blieb auf
  der abgelösten `v0.34.0`. Das `versions`-Gate fängt es **nicht**, weil es ein bloßes
  `:v0.34.0`-Tag ohne `ghcr.io/pt9912/d-check`-Präfix ist (darum bleibt `doc-check` grün) —
  eine echte Sweep-Lücke, kein Gate-Fang.
- **Failure-Szenario:** Ein Leser sieht im selben Abschnitt oben `pull …:v0.35.0` und
  wenige Zeilen darunter `:v0.34.0` als „empfohlene feste Version" und pinnt auf die
  superseded-Version im Glauben, sie sei aktuell.
- **verifizierbar:** ja — `grep -n ':v0.34.0' docs/user/benutzerhandbuch.md` (Zeile 74
  vs. `:v0.35.0`-Pull in Zeile 66); `versions`-Gate grün trotz Abweichung belegt den
  Blind-Spot.

---

## Negativbefunde (geprüft, ohne Befund)

- **MR-006 Referenzrichtung (Spec-Straten):** Kein Spec-Stratum verweist abwärts auf
  ADR-0027 oder Slices. `DC-FA-COMMITS-001` (Lastenheft) referenziert nur `DC-*`-IDs
  (`DC-FA-VCS-001`, `DC-QA-02/03`, `DC-FA-CONF-002`, `DC-FA-CLI-008/009`, `DC-FA-ID-001`,
  `DC-FA-MTX-001`) — **kein** ADR-/Slice-Link im Anforderungs-Körper; „in begleitender ADR
  festgehalten" ist die etablierte marker-lose Phrasierung (kein `ADR-NNNN`-Token).
  `DC-FA-COMMITS-001.a` (Spezifikation) referenziert nur `DC-*`/§-Sektionen.
  `architecture.md` benennt nur „opt-in Module `vcs`, `commits`" (keine ADR/Slice).
  Slice-Token stehen ausschließlich in den §7-Historie-Verweis-Spalten (Provenance, legitim,
  deckungsgleich mit dem `matrix`-`exclude-sections`-Muster).
- **ADR-0027 Referenzrichtung/Marker:** Der ADR-Körper trägt keine undeklarierten
  Abwärts-Slice-Token; `slice-056` erscheint nur in der `## Geschichte` (Zeile 198),
  `slice-*` in Zeile 31 ist ein **Zitat** der Traceability-Regel (kein `slice-\d+`-Token).
  ADR→ADR-Verweise sind für ADRs zulässig.
- **Anforderungs-Anlege-Prozess (`DC-FA-COMMITS-001`):** vier Akzeptanzkriterien
  (Happy / Boundary „Modul-aus/git-frei" / Negative / fail-closed) + explizite
  Out-of-Scope-Liste; Versions-Bump 0.35.0 + §7-Historie vorhanden; Bereichskürzel
  `COMMITS` in §3, Modul in `DC-FA-CLI-002` + Glossar. ACs deckungsgleich mit
  `DC-FA-COMMITS-001.a` (Exit 2 fail-closed, Exit 1 Befund, Exit 0 Happy).
- **Konsistenz commit-untraceable / Schema / Finding-Felder / Modi / Merge-Revert:**
  Grund-Code `commit-untraceable` identisch in Lastenheft, Spezifikation §4, `finding.go`
  (`ReasonCommitUntraceable`), `.d-check.yml`-Kommentar, Handbuch. Schema-Keys
  `commits.id-patterns`/`commits.exempt-pattern` identisch in Lastenheft/Spezifikation §2/
  `config.go` (`IDPatterns`/`ExemptPattern`)/`.d-check.yml`/`config_template.go`.
  Finding-Felder `file`=`target`=Kurz-SHA bzw. `pending` konsistent Lastenheft↔Spezifikation.
  Modi `--range`/`--commit-msg` konsistent. Betreff-Ausnahme `^(Merge |Revert )` identisch
  über alle Surfaces. Die dokumentierte Asymmetrie (leere `id-patterns`: Range-Modus inert
  vs. Message-Modus Exit 2) ist im §2-Schema-Kommentar und in `DC-FA-COMMITS-001.a` Schritt 1
  explizit spezifiziert (kein verdeckter Silent-Green).
- **Tombstone-Korrektheit:** `tools/trace-check.sh` ist in `codepaths.ignore-refs`
  (`.d-check.yml`, referenz-weit ⇒ deckt **alle** Inline-Code-Vorkommen: immutable ADR-0013,
  ADR-0027, `roadmap.md`, `slice-039`/`slice-056`, Spec-Historie, `docs/reviews/*`). **Keine**
  Markdown-**Links** auf `trace-check.sh` (grep leer) ⇒ keine `target-missing`-Gefahr über
  das `links`-Modul. Der `.githooks/commit-msg`-Kommentar-Verweis ist eine `.sh`-Datei
  (nicht markdown-gescannt). Logik hält, nicht nur das grüne Ergebnis.
- **ADR-Immutabilität:** ADR-0013 **unangetastet** (nicht im `git status`); kein
  `Accepted`-ADR-Körper inhaltlich geändert (nur der ADR-**Index** `README.md`, kein ADR-Body).
  Die ADR-0013-Teil-Supersede-Annotation ist korrekt auf Closure vertagt (Slice-DoD
  `:88-92`, `:110-112`) — jetzt **nicht** angefasst. ADR-0026-Index-Korrektur
  `Proposed`→`Accepted` faktisch korrekt (ADR-0026-Datei `Status:` = `Accepted`).
- **Cross-Surface-Vollständigkeit (`commits`):** Modul geführt in `validModules()`,
  `CommitsConfig` (`config.go`), `ReasonCommitUntraceable` (`finding.go`),
  `config_template.go` (Verfügbar-Liste + opt-in-Block), `suggest.go` (situative opt-in-Liste),
  `print_mk.go` (`doc-commits`-Target, `disableAllExcept("commits")`), Handbuch §5-How-to +
  §6-Modultabelle + §Änderungsverlauf, `DC-FA-CLI-002`, Glossar, `.d-check.yml`. Das Modul
  selbst ist **lückenlos** verkörpert; die einzige Rest-Lücke ist die `doc-commits`-**Target**-
  Enumeration (F-1/F-2), nicht das Modul.
- **`FOCUS_DISABLE` ↔ `.d-check.yml modules`:** `modules: [links, anchors, ids, matrix,
  codepaths, spans, hostpaths, versions]` == `FOCUS_DISABLE`-Liste (8/8 gespiegelt) ⇒
  `make trace-check` läuft fokussiert nur `commits`, kein Über-Feuern/Silent-Green.
  `commits` **nicht** in der Default-`modules`-Liste ⇒ Default-`doc-check` byte-identisch.
- **version.md / Pin-Sweep:** `<a id="v0.35.0">`-Anker gewandert (v0.34.0 verliert ihn);
  **keine** dangling `#v0.34.0`-Links in Live-Doku (grep leer); README-`docker run` und
  Handbuch-Voll-Pins auf `:v0.35.0`; Handbuch-Header `1.17` / `[v0.35.0](../../version.md#v0.35.0)`
  verweist auf den existierenden Anker; `## Aktuell` → `#v0.35.0`. (Ausnahme: das bloße
  Tag-Beispiel `:v0.34.0`, F-3.)
- **Modulname-Entscheidung:** `commits`-statt-`trace`-Begründung (Kollision mit `--trace`/RTM
  `DC-FA-CLI-009`) sauber dokumentiert in ADR-0027 §Kontext (`:51-56`) + §Verglichene
  Alternativen (`:125`) + Slice §2 (`:43-48`). Gate-Target bleibt `make trace-check`.
- **`--commit-msg`-Modus-Konsistenz:** in `cli.go` verdrahtet (`--commit-msg`-Flag,
  `runCommitMsg`, Kurzschluss-/Ausschluss-Logik) — die Doku-referenzierten Flag-Namen
  existieren (Feinprüfung der Go-Mechanik ist Zweit-Reviewer).

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 |
| LOW | 2 | F-2, F-3 |
| INFO | 0 | — |

---

## Verdikt

**NACHBESSERN.** Der Doc-first-Kern der Slice ist stark: `DC-FA-COMMITS-001` ist
prozesskonform angelegt (3+1 ACs + Out-of-Scope, Bump, Historie), MR-006 gewahrt, der
Tombstone deckt alle Referenzen, keine `Accepted`-ADR angetastet, das Modul `commits` ist
über alle Surfaces vollständig verkörpert, und die Pin-/Anker-Migration ist sauber.

Blockierend ist **F-1 (MEDIUM):** die Slice erweitert `DC-FA-CLI-010` laut eigener
§7-Historie, Spezifikation und ausgeliefertem Fragment auf **acht** Targets, lässt aber den
bindenden `DC-FA-CLI-010`-Anforderungs-Körper samt Happy-Path-AC und Out-of-Scope auf
**sieben** stehen — der Rang-1-Vertrag widerspricht sich selbst und klassifiziert
`doc-commits` sogar als out-of-scope. Das ist eine Spec-Treue-Lücke auf genau dem Vertrag,
den diese Slice anfasst, und per Reviewer-Skill blockieren MEDIUM-Befunde typischerweise den
Merge. **F-2/F-3 (LOW)** sind Doku-Drift im Benutzerhandbuch (Target-Zählung; ein
vom `versions`-Gate nicht gefangenes `:v0.34.0`-Beispiel) und sollten im selben Zug
mitgezogen werden.
