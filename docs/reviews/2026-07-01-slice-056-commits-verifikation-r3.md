# Verifikation R3 — slice-056 Modul `commits`: Fix-Verifikation (adversarial)

## Kopf-Metadaten

- **Datum:** 2026-07-01
- **Rolle:** unabhängiger Verifikations-Reviewer (R3). **Kein** dritter Voll-Review —
  ausschließlich Verifikation, dass die aus R1/R2 eingearbeiteten Fixes korrekt +
  vollständig sind und nichts Neues eingeschleppt wurde.
- **Gegenstand:** slice-056 — opt-in Modul `commits` (Commit-Message-Traceability,
  löst `tools/trace-check.sh` über den VCS-Port ab).
- **Vor-Reviews:** `2026-07-01-slice-056-commits-doc-r1.md` (F-1 MEDIUM, F-2/F-3 LOW),
  `2026-07-01-slice-056-commits-code-r2.md` (MEDIUM-1 Negativtest-Lücke).
- **Baseline:** `.harness/skills/reviewer.md` (v1.2.0). Working Tree gegen `HEAD` = `9a40b95`.
- **Methode:** Dateilektüre + gezielte `grep`; doc-check/Adversarial-Proben gegen das
  gebaute `d-check:latest` (read-only, `--network none`); **echte Mutation** für R2-MEDIUM
  (`make test` = Test-Image-Rebuild `d-check:test`; Runtime-`d-check:latest` unberührt).
  Kein neues Runtime-Image gebaut. DoD-Abhakung bewusst **nicht** geprüft (nicht Rolle).
- **Tree-Integrität:** Arbeitsbaum am Ende byte-identisch zum Vorgefundenen
  (git status deckungsgleich; `configyaml.go` sha256 `3d6f29…d82f` wiederhergestellt);
  einziges neues Artefakt = dieser Report.

---

## Fix-Verifikation (pro Befund: verifiziert ja/nein + Beleg)

### F-1 (R1, MEDIUM) — Lastenheft `DC-FA-CLI-010`-Körper „sieben" → „acht" inkl. `doc-commits`

**VERIFIZIERT — bestätigt.**

- **Beschreibung** (`spec/lastenheft.md:377`): „sowie **acht** `##`-annotierte Targets" —
  die Aufzählung nennt jetzt `doc-immutable` (`:384`) **und** `doc-commits` (`:389`) mit
  Anker-Link auf `DC-FA-COMMITS-001`.
- **Happy-Path-AC** (`:406`): enumeriert exakt acht Targets — `doc-check`, `doc-trace`,
  `doc-complete`, `doc-doctor`, `doc-repair`, `doc-immutable`, **`doc-commits`**, `doc-help`.
- **Boundary-AC** (`:407`): führt `…/doc-immutable/doc-commits` und deren Modus
  (`doc-commits` → `--enable commits` + Fokus-`--disable` + `--range $(RANGE)`).
- **Out-of-Scope** (`:410`): „jenseits der gelisteten **acht**
  (`…/doc-immutable/doc-commits/doc-help`)" — `doc-commits` ist damit **nicht** mehr als
  out-of-scope klassifiziert (der genaue R1-F-1-Widerspruch ist behoben).
- **Konsistenz-Anker:**
  - (a) `spec/spezifikation.md:335` „**Acht** `.PHONY`-Targets", `doc-immutable` (`:349`),
    `doc-commits` (`:357`).
  - (b) real ausgeliefertes `--print-mk` (`docker run … d-check:latest --print-mk`):
    **8** `doc-*`-Targets — `doc-check, doc-commits, doc-complete, doc-doctor, doc-help,
    doc-immutable, doc-repair, doc-trace` (`grep -cE '^doc-[a-z]+:'` = 8).
  - (c) `config_template.go`: `commits` in Verfügbar-/opt-in-Liste (R1-Negativbefund,
    unverändert grün — doc-check bestätigt).
  - (d) Handbuch §6/§4.13 „acht" (siehe F-2).
- **Repo-weite Enumerations-Suche:** `grep "sieben\|sechs"` in `spec/` + `docs/user/` = **leer**.
  Repo-weit außerhalb der Review-Reports treten `sieben`/`sechs` nur in `.harness/baseline/**`
  (Fremd-Baseline), `docs/plan/planning/done/**` (eingefrorene Historie, u. a. slice-047
  „sechs Targets" der Vorstufe), `CHANGELOG.md:651` (6-Modul-Ära) und `adr/0024:160`
  („sieben Klassen des Skript-Selbsttests" — nicht print-mk) auf. **Keine** dieser Stellen
  ist die `DC-FA-CLI-010`-Target-Enumeration — Fix vollständig.

### F-2 (R1, LOW) — Handbuch §4.13-How-to „sechs" → „acht" inkl. `doc-immutable` + `doc-commits`

**VERIFIZIERT — bestätigt.** `docs/user/benutzerhandbuch.md:620-622`: „`TRACE_FLAGS` und
**acht** `##`-annotierte Targets (`doc-check`, `doc-trace`, `doc-complete`, `doc-doctor`,
`doc-repair`, `doc-immutable`, `doc-commits`, `doc-help`)". Sowohl das zuvor fehlende
`doc-immutable` als auch das neue `doc-commits` sind gelistet. §6-Modultabelle + §5-How-to
(`:781-786`, `doc-commits RANGE=…`) unverändert konsistent.

### F-3 (R1, LOW) — stale `:v0.34.0` in Live-Prosa

**VERIFIZIERT — bestätigt.** `benutzerhandbuch.md:74` zeigt jetzt `:v0.35.0` („eine feste
Version, empfohlen"). `grep "v0.34.0" docs/user/benutzerhandbuch.md README.md` liefert nur
**`benutzerhandbuch.md:923`** — eine Verlaufstabellen-Zeile (Handbuch-Version 1.16 ↔ Software
v0.34.0), historisch legitim, kein `:v0.34.0`-Pin. README.md: kein Treffer. Header
(`:3`) verweist `[v0.35.0](../../version.md#v0.35.0)` — Ziel-Anker existiert
(`version.md:35` `<a id="v0.35.0"></a>`; `## Aktuell` → `v0.35.0`).

### R2-MEDIUM-1 — Negativtest für die drei `applyCommits`-Ablehnungs-Guards (MUTATIONS-VERIFIKATION)

**VERIFIZIERT — bestätigt durch echte Mutation.** Der Test bindet den Silent-Grün-Guard
tatsächlich.

- **Guards vorhanden** (`configyaml.go:287-296`): leerer Eintrag (`:288`), nicht
  kompilierbares Regex (`:292`), **Leerstring-Match** `re.MatchString("")` (`:294`).
- **Neue Tests** (`configyaml_test.go:153` `TestDecode_CommitsFehler`, `:177`
  `TestDecode_CommitsHappy`): der Fehler-Test prüft alle drei Guards; für den
  Leerstring-Guard fordert er über `.*`/`X*`/`(ADR)?` explizit eine Fehlermeldung mit
  Substring „Leerstring".
- **Mutation (a):** Guard-Block `if re.MatchString("") { … }` aus `applyCommits` entfernt.
- **Mutation (b/c):** `make test` → **`--- FAIL: TestDecode_CommitsFehler`** mit
  `configyaml_test.go:162: commits.id-patterns ".*": err = <nil> (Leerstring-Ablehnung
  erwartet)`; Paket `configyaml` FAIL, `make` Exit 1. Der Test fällt also real → der
  Negativtest ist **wirksam** (kein toter Test).
- **Mutation (d):** `configyaml.go` byte-exakt aus Backup wiederhergestellt —
  sha256 `3d6f29178f0b25829023f9feff5e838da98d15d11198478ce7efdca794f6d82f` == Ausgangswert;
  Guard wieder bei `:294`; `git diff --stat` zeigt die intakten slice-056-Änderungen
  (61 insertions, 7 deletions).
- **Mutation (e):** `make test` wieder grün (`ok … configyaml`).

Kein HIGH. Die R2-MEDIUM-Regressionslücke ist geschlossen und der Schutz mutations-belegt.

---

## Kein neuer Schaden durch die Fixes

- **doc-check grün:** `docker run … d-check:latest` (read-only, `--network none`) ⇒
  „163 Datei(en) geprüft, 0 Befund(e)", Exit 0. Der **neue** Anker-Link in
  `DC-FA-CLI-010` (`doc-commits` → `#dc-fa-commits-001--…-modul-commits-opt-in`) löst auf —
  Ziel-Heading `spec/lastenheft.md:1192`; ein dangling Anker würde das `anchors`-Modul melden.
- **Review-Reports doc-check-sauber:** `docs/reviews/` steht **nicht** in `scan.ignore`
  (`.d-check.yml:15` = nur `.harness/baseline/**`, `.harness/cache/**`); die Reports sind mithin
  gescannt (bei matrix/ids nur token-exempt, `:23`) — das grüne doc-check über 163 Dateien
  belegt: keine `span-unclosed`/Anker-Defekte in R1/R2 (oder diesem R3).
- **Keine neue Inkonsistenz:** die Fixes berührten nur Zähl-/Aufzählungs-Prosa +
  Pin/Anker + einen Testblock; alle Konsistenz-Anker (F-1 (a)-(d)) ziehen jetzt gleich.

---

## Finaler adversarialer End-to-End-Check (kein Silent-Grün-Pfad)

Alle Proben gegen `d-check:latest`, read-only, `--network none`:

- **(a) `--commit-msg -`, kennungslose Message** („chore: adjust some wording without any
  tracer"): Ausgabe `pending:1 pending commit-untraceable` +
  `d-check: commit-untraceable — Commit-Message ohne DC-/ADR-/MR-/slice-ID`, **Exit 1**.
- **(b) `--commit-msg -`, ID-tragende Message** („feat(commits): … (slice-056, ADR-0027)"):
  keine Ausgabe, **Exit 0** (legitim still).
- **(c) Range-Probe** gegen ein Wegwerf-git-Repo (scratchpad, eigenes `.d-check.yml` mit
  `commits.id-patterns`; Basis + kennungsloser refactor-Commit + `ADR-0027`-Commit),
  Aufruf `--enable commits <FOCUS_DISABLE> --range BASE..HEAD`: **genau 1 Befund** —
  `5a883ec:1 5a883ec commit-untraceable` (der kennungslose Commit), Basis ausgeschlossen,
  der `ADR-0027`-Commit sauber, **Exit 1**.
- **Positiv-Kontrolle** (saubere Teil-Range, nur der `ADR-0027`-Commit): „0 Befund(e)",
  **Exit 0** — bestätigt, dass Exit 0 ein echter Pass ist, kein verschluckter Fehler.

**Bestätigt: kein Silent-Grün-Pfad.** Der einzige stille Pass ist der legitim kennungs-
tragende Fall; kennungslose Messages/Commits erzeugen Befund + non-zero in beiden Modi.

---

## Negativbefunde (geprüft, ohne Befund)

- **F-1/F-2/F-3-Regression:** keine verbliebene „sieben"/„sechs"-Enumeration in bindenden
  oder Nutzer-Surfaces; alle acht Targets deckungsgleich über Lastenheft/Spezifikation/
  print-mk/Handbuch.
- **Version/Anker:** `version.md#v0.35.0` existiert; keine `:v0.34.0`-Live-Pins.
- **Tree-Integrität:** git status nach allen Proben deckungsgleich zum Ausgangszustand
  (staged Delete `tools/trace-check.sh`, 24 modifizierte, 7 untracked); `configyaml.go`
  sha256 wiederhergestellt; Wegwerf-Repo + Backup gelöscht; `make test` grün.
- **Runtime-Image:** `d-check:latest` unberührt (nur `d-check:test` durch `make test`).

---

## Kategorie-Summary (nur durch die Fixes neu eingeschleppt)

| Kategorie | Anzahl | IDs |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 0 | — |
| INFO | 0 | — |

---

## Verdikt

**VERIFIED.** Alle vier eingearbeiteten Fixes sind korrekt + vollständig belegt:
F-1 (Lastenheft-Körper „acht" inkl. `doc-commits`, konsistent zu Spezifikation/print-mk/
config_template/Handbuch), F-2/F-3 (Handbuch-How-to „acht", kein stale `:v0.34.0`, Anker
existent). Die R2-MEDIUM-Testlücke ist geschlossen und der Leerstring-Silent-Grün-Guard
**mutations-belegt** wirksam getestet (Guard-Entfernung ⇒ `TestDecode_CommitsFehler` fällt).
Keine neue Inkonsistenz eingeschleppt (doc-check grün, Anker löst auf, Reports sauber);
der adversariale End-to-End-Check zeigt in beiden Modi keinen Silent-Grün-Pfad. Der
Arbeitsbaum ist byte-identisch zum Vorgefundenen hinterlassen (nur dieser Report neu).
