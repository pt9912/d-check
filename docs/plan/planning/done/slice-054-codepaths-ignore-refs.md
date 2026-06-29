# Slice slice-054: `codepaths.ignore-refs` — die Frozen-Doc-Refactoring-Falle auflösen und `adr-immutable-check.sh` entfernen

**Status:** done (welle-43-ignore-refs).

**Welle:** welle-43-ignore-refs (Trigger: Auftraggeber — „Wir brauchen eine bessere
Lösung, dass solch ein Problem — etwas wird refaktoriert/gelöscht — nicht immer
wieder auftaucht." Konkreter Auslöser: das in slice-053 abgelöste, aber **behaltene**
`tools/adr-immutable-check.sh` ([`slice-053`](../done/slice-053-vcs-modul.md),
[ADR-0024](../../adr/0024-vcs-immutable-gate.md) entschied „pfad-stabil behalten") —
weil die **immutable** [ADR-0016](../../adr/0016-adr-immutable-gate.md) es in
Inline-Code referenziert und `codepaths` das als Existenz-Pflicht erzwingt.)

**Bezug:** Erweitert die bestehende Anforderung
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(Modul `codepaths`) um `ignore-refs` plus eine begleitende **neue ADR**
(`Supersedes` die „Skript-behalten"-Teilentscheidung von
[ADR-0024](../../adr/0024-vcs-immutable-gate.md); ihr Rest — VCS-Port,
Modul `vcs` — bleibt gültig). Damit wird das in slice-053 nur als „pfad-stabiler
Fallback" behaltene `tools/adr-immutable-check.sh` **vollständig entfernt**.

**Autor:** pt9912. **Datum:** 2026-06-29.

---

## 1. Ziel

**Die Wurzel:** Immutable/historische Doku (Accepted-ADRs, `done/`-Slices) trägt
*eingefrorene* Inline-Code-Verweise (z. B. `tools/foo.sh`). <!-- d-check:ignore (illustrativer Beispielpfad) --> `codepaths` prüft sie als
**lebende** Zeiger (Pfad muss *jetzt* existieren). Wird der Code refaktoriert/gelöscht,
dangelt der eingefrorene Verweis → `codepath-missing` → **nicht fixbar** (die ADR ist
immutable). Das tritt bei **jedem** solchen Refactoring auf — zuletzt blockierte es in
slice-053 die Entfernung von `adr-immutable-check.sh` (deshalb dort als Workaround
„behalten").

**Die Lösung:** ein Config-Schlüssel `codepaths.ignore-refs` — ein **Tombstone-Register**
bewusst entfernter Artefakte: Inline-Code-Verweise, die einem `ignore-refs`-Glob
entsprechen, lösen kein `codepath-missing` aus (egal wo, auch in immutabler Doku).
Bewusster Akt **mit Gate**: wer löscht und `ignore-refs` vergisst, fällt weiter auf
`codepath-missing` — nichts dangelt still (wie Versions-/Digest-Pins). Damit wird das
Refactoring/Löschen wiederkehrend sauber, **ohne** Edits an immutabler Doku und **ohne**
ganze Doc-Klassen aus dem Check zu nehmen. Erster Anwendungsfall: `adr-immutable-check.sh`
entfernen.

## 2. Entscheidungen

- **`codepaths.ignore-refs` (Glob-Liste, `matchGlob` wie `exempt-paths`/`scan.ignore`).**
  Ein Inline-Code-Pfad, der einem Eintrag entspricht, wird **nicht** existenz-geprüft.
  Default leer → ohne `ignore-refs` byte-identisch
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- **Per-Pfad, nicht per-Datei-Klasse.** Bewusst **kein** pauschales Ausnehmen von
  `docs/plan/adr/**`/`done/**` (das verlöre die `codepaths`-Abdeckung der *übrigen*
  Verweise dieser Dateien). `ignore-refs` nimmt nur den **einen** entfernten Pfad aus;
  alles andere bleibt geprüft.
- **Bewusster Akt mit Gate.** `ignore-refs` ist ein Opt-in-Register: vergisst man beim
  Löschen den Eintrag, meldet `codepath-missing` weiter — der Gate erzwingt die
  Tombstone-Deklaration, statt still durchzulassen.
- **Kein Doc-Marker.** Ein `<!-- d-check:ignore -->` in der ADR-Zeile wäre eine
  Core-Änderung an immutabler Doku (→ `adr-check` FAIL) und kann Accepted-ADRs nicht
  nachträglich erreichen; der Default („Verweise sind live") wäre für historische
  Records falsch herum. `ignore-refs` lebt in der Config → kein Doku-Edit, retrofittet
  bestehende ADRs. (Der vorhandene `d-check:ignore`-Zeilen-Marker bleibt für historische
  Verweise in *lebenden* Dateien.)
- **`adr-immutable-check.sh` wird entfernt.** `git rm` + `tools/adr-immutable-check.sh`
  ins `ignore-refs` (die immutablen [ADR-0016](../../adr/0016-adr-immutable-gate.md)/[ADR-0024](../../adr/0024-vcs-immutable-gate.md)-Inline-Referenzen
  werden so zu deklarierten Historien-Verweisen). Die neue ADR nimmt die
  „Skript-behalten"-Teilentscheidung jener ADR zurück (`Supersedes`).
- **Config-Surface mitziehen** (Lehre aus slice-053): `--print-config` und das
  Benutzerhandbuch dokumentieren `ignore-refs`; `--suggest-config ai-harness` bei Bedarf.
- **Name/Form** als Auftraggeber-bestätigter Default: `codepaths.ignore-refs`, Glob.
  (Alternativen `removed-refs`/`tombstones` verworfen zugunsten der Auftraggeber-Nennung
  „ignore-code-ref".)

## 3. Definition of Done

### 3a. Artefakte (Doc-first → Code → Config-Surface → Entfernung)

- [ ] **Doc-first:**
  [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
  um `ignore-refs` erweitert (Beschreibung + Akzeptanzkriterien + Out-of-Scope),
  Versions-Bump **0.34.0** + §7-Historie; die `.a`-Algorithmus-Sektion um den
  `ignore-refs`-Schritt + Schema-Key (§2) ergänzt; **neue ADR** (`Supersedes` die
  „Skript-behalten"-Teilentscheidung von [ADR-0024](../../adr/0024-vcs-immutable-gate.md),
  Bezug [ADR-0016](../../adr/0016-adr-immutable-gate.md)) + ADR-Index.
- [ ] **Code:** `model.CodepathsConfig.IgnoreRefs`; `configyaml` parst `ignore-refs`
  (Glob-Validierung wie `exempt-paths`); die `codepaths`-Regel überspringt die
  Existenz-Prüfung für matchende Pfade (`anchor-missing` für solche Pfade entfällt
  mit). Tests: Happy (ignored Pfad → kein Befund), Negative (nicht-ignored fehlender
  Pfad → weiter `codepath-missing`), Glob, Default-leer byte-identisch.
- [ ] **Config-Surface (nicht vergessen):** `--print-config`
  ([config_template.go](../../../../internal/adapter/driving/cli/config_template.go),
  `codepaths`-Block + `ignore-refs`), Benutzerhandbuch (§6/§Weitere Module), und —
  falls einschlägig — der `--suggest-config ai-harness`-`codepaths`-Block.
- [ ] **Entfernung:** `git rm tools/adr-immutable-check.sh` + Eintrag in
  `.d-check.yml` `codepaths.ignore-refs`; Kommentare in `Makefile` (`adr-check`),
  [`harness/README.md`](../../../../harness/README.md) §Sensors und
  [`AGENTS.md`](../../../../AGENTS.md) §4 von „pfad-stabiler Fallback" entschlacken.

### 3b. Verifikation (Korrektheit — „richtig gebaut")

- [ ] `make ci` grün (doc-check, lint, test, arch-check, Coverage ≥ 93 %, semgrep,
  gate-consistency, planning-check, image-test) + `make completeness-check` (keine
  Requirements-Waisen).
- [ ] Akzeptanztests des neuen `ignore-refs`-Verhaltens spezifikations-konform
  (Happy/Negative/Glob/Default-leer); `configyaml`-Validierung der `ignore-refs`-Globs.
- [ ] Zwei **unabhängige** Reviews (R1 Doc, R2 Code/Config) — alle HIGH/MEDIUM behoben,
  Reports unter `docs/reviews/`.

### 3c. Validierung (Wirksamkeit — „das Richtige gebaut")

- [ ] **Die Falle ist am realen Fall aufgelöst:** `tools/adr-immutable-check.sh` ist
  **entfernt** und `make doc-check` bleibt **grün** — die eingefrorenen
  [ADR-0016](../../adr/0016-adr-immutable-gate.md)/[ADR-0024](../../adr/0024-vcs-immutable-gate.md)-Inline-Referenzen brechen *nicht* mehr (Beweis, dass `ignore-refs`
  immutable Doku retrofittet, ohne sie zu editieren).
- [ ] **Mechanismus generisch demonstriert:** ein entfernter Pfad ohne `ignore-refs`-
  Eintrag feuert `codepath-missing` (Gate erzwingt die Tombstone-Deklaration); mit
  Eintrag ist er still — kein klassenweites Ausnehmen nötig, übrige Verweise der
  Datei bleiben geprüft.
- [ ] **Auffindbarkeit (slice-053-Lehre):** `--print-config` und das Benutzerhandbuch
  zeigen `ignore-refs` — ein Adopter findet den Mechanismus ohne Quelltext-Studium.
- [ ] **Keine Regression der Verteilung/Hermetik:** kein neues Modul, keine neue
  Dependency; Default-Lauf byte-identisch.

### 3d. Closure

- [ ] Move nach `done/` + Roadmap-Flip
  ([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise));
  die neue ADR → Accepted + die [ADR-0024](../../adr/0024-vcs-immutable-gate.md)-Teil-Supersede-Annotation (Geschichte + Index,
  **erst bei Closure** — wie slice-053-R1-F-4; der Index trägt die Annotation **nicht**,
  solange die neue ADR `Proposed` ist, R1-MEDIUM-1) + ein `## Geschichte`-Append an die
  immutable [ADR-0016](../../adr/0016-adr-immutable-gate.md) (Skript entfernt, Verweis auf
  die neue ADR — erlaubter Geschichte-Anhang, ihr Core bleibt unangetastet; R1-INFO-3);
  Release **v0.34.0**.

## 4. Risiken / offene Punkte

- **`ignore-refs` als stilles Ventil missbrauchbar:** ein zu breiter Glob könnte echte
  `codepath-missing` verstecken. Mitigation: per-Pfad-Disziplin (nicht klassenweit), der
  Glob ist eng zu halten; der Eintrag ist ein bewusster, im Diff sichtbarer Akt.
- **Wächst über die Zeit:** das Register sammelt entfernte Pfade. Bewusst — es ist ein
  ehrliches Tombstone-Verzeichnis; bounded auf tatsächlich entfernte, noch zitierte
  Artefakte.
- **Die [ADR-0024](../../adr/0024-vcs-immutable-gate.md)-Teilentscheidung umkehren:** braucht regelkonform eine neue ADR
  (`Supersedes`); der VCS-Port/Modul-Teil jener ADR bleibt unberührt.
- **Globaler Scope:** `ignore-refs` greift überall (auch in lebenden Dateien). Akzeptabel —
  ein bewusst entfernter Pfad ist überall „nicht da"; lebende Verweise räumt man ohnehin auf.
- **`ignore-refs` ist `codepaths`-lokal — die `links`-Achse bleibt (R1-INFO-1):** ein
  entfernter Pfad, der als **Markdown-Link** (statt Inline-Code) in immutabler Doku steht,
  feuert weiter `target-missing` des Moduls `links`; `ignore-refs` deckt das nicht, `links`
  hat kein referenz-weites Pendant. Aktuell rein latent (kein immutables Dokument verlinkt
  das Skript in Link-Form). Folge: ein eigenes `links`-Tombstone-Ventil wäre eine
  Folge-Anforderung (in der neuen ADR als Re-Evaluierungs-Trigger verortet).
- **Rückwirkende Review-Record-Bearbeitung (R1-INFO-2):** beim Entfernen brach ein
  Markdown-Link aufs Skript im historischen Review `2026-06-28-slice-052-immutable-r2.md`
  (`links`/`target-missing` — nicht von `ignore-refs` gedeckt; `docs/reviews/**` ist nur
  `codepaths`-exempt). Er wurde zu Inline-Code entlinkt: **Inhalt unverändert, nur die
  Link-Form nachgezogen**, bewusst hier dokumentiert statt still.

## 5. Trigger

Auftraggeber 2026-06-29: nach der slice-053-Frage „Brauchen wir
`tools/adr-immutable-check.sh` noch?" und der Feststellung, dass das Behalten nur ein
Workaround gegen die `codepaths`-/Immutabilitäts-Falle war — „Wir brauchen eine bessere
Lösung, dass solch ein Problem … nicht immer wieder auftaucht." Idee `ignore-code-ref`
im `.d-check.yml` vom Auftraggeber.

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code + Spec; „Doc führt, Code folgt"). `ignore-refs` ist eine additive
`codepaths`-Erweiterung; keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

**Geliefert:** Modul `codepaths` um `codepaths.ignore-refs` erweitert — eine Glob-Liste
nimmt **aufgelöste Ziel-Pfade** referenz-weit (datei-/zeilen-unabhängig) von der
Existenz-/Escape-/Anker-Prüfung aus (Skip vor allen drei Grund-Codes, Spezifikation
[§DC-FA-CODE-001.a](../../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) Schritt 5). Tombstone-Register bewusst entfernter Artefakte; löst die
Frozen-Doc-Refactoring-Falle ohne Edit an immutabler Doku und ohne klassenweites
Ausnehmen. Default leer → byte-identisch. **Falle am realen Fall bewiesen:**
`tools/adr-immutable-check.sh` per `git rm` entfernt, `make doc-check` bleibt grün — die
eingefrorenen [ADR-0016](../../adr/0016-adr-immutable-gate.md)/[ADR-0024](../../adr/0024-vcs-immutable-gate.md)-Inline-Referenzen
sind als Tombstones deklariert.

**Doc-first:** [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(Lastenheft 0.34.0), Spezifikation [§DC-FA-CODE-001.a](../../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code) + §2-Schema;
[ADR-0025](../../adr/0025-codepaths-ignore-refs.md) (Accepted, `Supersedes` die
„Skript-behalten"-Teilentscheidung von [ADR-0024](../../adr/0024-vcs-immutable-gate.md) —
deren VCS-Kern bleibt); Geschichte-Annotation an [ADR-0024](../../adr/0024-vcs-immutable-gate.md)/[ADR-0016](../../adr/0016-adr-immutable-gate.md).

**Verifikation:** `make ci` grün (doc-check 155/0, lint, test, arch-check, Coverage
93,40 %, semgrep 0, gate-consistency, planning-check, image-test) + `completeness-check`
0 Waisen. Zwei unabhängige Reviews (R1 doc 0H/1M/1L/3I, R2 code 0H/1M/1I) — alle Befunde
behoben; der Escape-/Anker-Lock-Test wurde per Mutationstest verifiziert (Guard hinter
`escaped` → rot). Reports unter `docs/reviews/2026-06-29-slice-054-ignore-refs-doc-r1.md`
und `docs/reviews/2026-06-29-slice-054-ignore-refs-code-r2.md`.

**Validierung:** Mechanismus per-Pfad statt klassenweit (Test);
`--print-config`/`--suggest-config`/Benutzerhandbuch zeigen `ignore-refs`
(Auffindbarkeit); kein neues Modul, keine neue Dependency.

**Bewusst offen (Re-Eval):** `ignore-refs` ist `codepaths`-lokal — die `links`-Achse (ein
**Markdown-Link** statt Inline-Code auf einen entfernten Pfad) bleibt eine Rest-Falle, in
[ADR-0025](../../adr/0025-codepaths-ignore-refs.md) als Re-Evaluierungs-Trigger verortet.
Beim Entfernen wurde ein solcher Markdown-Link im historischen Review slice-052-R2 zu
Inline-Code entlinkt (Inhalt unverändert).

**Release:** **v0.34.0** auf GHCR (Pipeline-Run 28381441115 grün), Digest-Pin
`ghcr.io/pt9912/d-check@sha256:1cf0837fd62daa077be2705a2fa77d791a11b5f1d07ea65bb9c8c00c2116d64a`.
