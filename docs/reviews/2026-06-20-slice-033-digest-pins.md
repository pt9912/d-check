# Review — slice-033 Implementierung (Image-Pins auf Digest)

## Kopf-Metadaten

- **Review-Art:** Code-Review (Working-Tree-Diff gegen Plan/ADR/Anforderungen/
  Hard Rules — **kein Verifier**: DoD-Abhaken und Gate-Lauf-Bestätigung sind
  nicht Gegenstand; Gates/Registry-Auflösungen werden NICHT als grün/korrekt
  angenommen).
- **Datum:** 2026-06-20
- **Gegenstand:** Working-Tree-Diff (`git diff HEAD`) der slice-033-Umsetzung
  plus die untrackte Datei `docs/plan/adr/0011-digest-pins-build-gate-images.md`
  (NEU) und der gestagete Rename
  `docs/plan/planning/next/…` → `docs/plan/planning/done/slice-033-dockerfile-digest-pins.md`.
  Inhalt: alle drei externen `FROM`-Zeilen im `Dockerfile` (`golang`,
  `golangci/golangci-lint`, `gcr.io/distroless/static-debian12`) sowie das
  Scanner-Image in `tools/semgrep.sh` werden von reinem Tag auf
  `image:tag@sha256:…` gehoben; `make versions` bekommt eine semgrep-Image-Zeile
  (Version+Digest per `grep -oP` aus `tools/semgrep.sh`); ADR-0011 (NEU)
  vereinheitlicht ADR-0002 §1 / ADR-0010 auf Digest; ADR-Index ergänzt;
  Slice-Plan nach `in-progress` verschoben und auf inline-Digest-Politik
  umgeschrieben.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/done/slice-033-dockerfile-digest-pins.md`; ADR
  `docs/plan/adr/0011-digest-pins-build-gate-images.md` (vereinheitlicht
  ADR-0002 §1, hebt ADR-0010 von Tag auf Digest); die nach Index immutablen
  Accepted-ADRs `docs/plan/adr/0002-distribution-ghcr-image.md` und
  `docs/plan/adr/0010-semgrep-hermetisches-gate.md`; Anforderung `DC-QA-02`
  (Determinismus — identische Eingabe ⇒ identische Befunde); Hard Rules
  `AGENTS.md` §3 (insb. §3.1 Docker/make-only, §3.5 ADR-Immutability, §3.6
  Gate-Lockerung nur per ADR); ADR-Index-Regel (Accepted = immutable). Kein
  Vertrags-`DC-FA-*` (Prozess-/Reproduzierbarkeits-Härtung). Vorläufer:
  `docs/reviews/2026-06-20-slice-032-semgrep-gate.md` (INFO-1 = Tag-statt-Digest,
  hier adressiert; HIGH-1 = stilles Grün ohne Regel-Nachweis, bereits in
  `tools/semgrep.sh` Zeilen 64–73 gefixt und hier unberührt). **Netzzugriff
  (Bash-Registry-Tools, WebFetch, WebSearch) war in dieser Umgebung gesperrt —
  die Digest↔Manifest-Typ-Prüfung konnte nicht registry-seitig erfolgen; das ist
  als Verifizierbarkeitspfad notiert, nicht als bestätigte Tatsache.** Tests/Gates
  wurden NICHT ausgeführt.

## Findings

### HIGH

Keine.

### MEDIUM

#### MEDIUM-1 — Manifest-Listen-Eigenschaft der vier Digests nicht belegt; ein plattformspezifischer Digest bräche den arm64-Build

- **Kategorie:** MEDIUM
- **Quelle:** `DC-QA-02` (Determinismus/Reproduzierbarkeit); ADR-0011 §3
  (Manifest-Listen-Digest verbindlich „damit der Build auf amd64/arm64
  funktioniert")
- **Pfad:** `Dockerfile:34`
- **Befund:** ADR-0011 §3 entscheidet, gepinnt werde der Manifest-Listen-Digest,
  nicht ein plattformspezifischer. Der Diff liefert die vier `@sha256:`-Werte
  hart verdrahtet, ohne Beleg (Manifest-`mediaType`, Beschaffungs-Kommando), dass
  sie Index-/Listen-Digests sind und nicht ein einzelner Plattform-Manifest-Digest.
  Ist auch nur einer plattformspezifisch (z. B. amd64-Manifest statt Index),
  scheitert ein `docker build` auf arm64 beim Pull mit „no matching manifest for
  linux/arm64" — die von ADR-0011 zugesicherte Multi-Arch-Reproduzierbarkeit
  bricht. Dieselbe Stelle gilt für `Dockerfile:55`, `Dockerfile:113` und
  `tools/semgrep.sh:23`.
- **Verifizierbar:** ja — `docker buildx imagetools inspect <ref>@<digest>` bzw.
  `skopeo inspect --raw docker://<ref>@<digest>` je Digest: ist `mediaType` ein
  `image.index`/`manifest.list`, ist es korrekt; ist es ein
  `image.manifest`/`manifest.v2` mit einzelner Plattform, bestätigt das den
  Befund. (Konnte hier mangels Netz/Registry-Zugriff nicht ausgeführt werden.)

### LOW

#### LOW-1 — `make versions` semgrep-Zeile hängt an GNU/PCRE-`grep -oP`; ohne `-P`-Support stiller leerer Pin-Beleg

- **Kategorie:** LOW
- **Quelle:** Maintainability (Reviewer-Skill LOW-Anker „latente Wartungsfalle,
  die erst bei künftigem Edit/anderem Host zündet"); `DC-QA-02` (Pin-Beleg)
- **Pfad:** `Makefile:89`
- **Befund:** Die neue Beleg-Zeile extrahiert Version und Digest per
  `grep -oP 'SEMGREP_VERSION:-\K[^}]+'` / `'SEMGREP_DIGEST:-\K[^}]+'`. `-P`
  (PCRE, inkl. `\K`) ist eine GNU-grep-Eigenheit; BSD-/busybox-grep kennt es
  nicht. Auf einem Host ohne PCRE-grep liefert der `$(…)`-Teilausdruck leer,
  während die umschließende `@echo`-Zeile mit Exit 0 durchläuft — Ausgabe
  `semgrep-image=semgrep/semgrep:@`, also ein **stiller leerer Pin-Beleg** statt
  lautem Fehler. Das übrige Makefile nutzt portables `grep -E` (`Makefile:88`);
  diese Zeile führt erstmals eine `-oP`-Abhängigkeit in den Beleg-Pfad ein.
- **Verifizierbar:** ja — `make versions` auf einem Host mit BSD-/busybox-grep
  (oder `grep` ohne `-P`) zeigt die leere `semgrep-image=…:@`-Zeile bei Exit 0.

#### LOW-2 — Tag und Digest gemeinsam gehalten ohne Konsistenz-Prüfung; bei Einzel-Hebung wird der Tag zum irreführenden Label

- **Kategorie:** LOW
- **Quelle:** Maintainability (Reviewer-Skill LOW-Anker „latente Wartungsfalle");
  `DC-QA-02`
- **Pfad:** `tools/semgrep.sh:23`
- **Befund:** Image-Referenzen tragen Tag UND Digest (`image:tag@sha256:…`).
  Docker löst ausschließlich über den Digest auf und ignoriert den Tag-Teil
  dabei; es gibt keine Stelle (kein Gate, kein Skript), die prüft, dass der Tag
  tatsächlich auf den danebenstehenden Digest zeigt. Wird bei einer Hebung nur
  einer der beiden Werte geändert (entgegen der ADR-0011-§4-Politik „Version UND
  Digest gemeinsam"), zieht Docker still die Digest-Variante, während der
  Klartext-Tag eine andere Version vorspiegelt — der Tag, dessen einziger Zweck
  laut ADR die Lesbarkeit ist, wird zur Fehlinformation. Gleiches Muster:
  `Dockerfile:34`, `Dockerfile:55`, `Dockerfile:113`.
- **Verifizierbar:** ja — ein Edit, der `SEMGREP_VERSION` auf einen anderen Wert
  setzt, aber `SEMGREP_DIGEST` unverändert lässt, läuft fehlerfrei durch (Docker
  zieht den alten Digest); Tag und tatsächlich gezogenes Image divergieren ohne
  Warnung.

### INFO

#### INFO-1 — `latest`-Tag-Drift (ADR-0002 §4 ↔ release.yml) bleibt offen, ist aber sauber an eine Folge-ADR delegiert

- **Kategorie:** INFO
- **Quelle:** ADR-0002 §4 („kein `latest`"); ADR-0011 §Konsequenzen
  (`latest`-Politik „bleibt unberührt … eigene Folge-ADR")
- **Pfad:** `docs/plan/adr/0011-digest-pins-build-gate-images.md:66`
- **Befund:** `release.yml` taggt und pusht weiterhin `:latest`
  (`.github/workflows/release.yml` Zeilen 87/129), was ADR-0002 §4 („kein
  `latest`") widerspricht. ADR-0011 grenzt diesen Punkt ausdrücklich aus dem
  Slice-Umfang aus und weist ihn einer eigenen Folge-ADR zu; der Diff fasst
  `release.yml`/`.github` nicht an (Diff-Stat: nur Dockerfile, Makefile,
  ADR-README, Slice-Plan, semgrep.sh). Die Abgrenzung ist explizit und leckt
  nicht in diesen Diff — dokumentationswürdig als bewusst offener Folgepunkt.
- **Verifizierbar:** ja — `git diff HEAD --stat` zeigt keine `.github`-Änderung;
  `release.yml:87`/`:129` belegen den fortbestehenden `latest`-Push.

## Negativbefunde (geprüft, ohne Befund)

- **Immutability ADR-0002/ADR-0010 (AGENTS.md §3.5, Index-Regel):** REFUTED, dass
  eine Accepted-ADR editiert wurde. `git diff HEAD -- docs/plan/adr/0002-…md` und
  `… 0010-…md` sind **leer**; beide Dateien erscheinen nicht in `git status`. Die
  Vereinheitlichung läuft ausschließlich über die NEUE ADR-0011, die in §2 die
  Bezüge nennt („Vereinheitlicht ADR-0002 §1 … hebt ADR-0010 vom Tag auf den
  Digest"); ADR-0002 behält Status Accepted/2026-06-10, ADR-0010
  Accepted/2026-06-19 unverändert. Kein Inhaltsüberschreiben.
- **`make versions` Digest-Beleg über FROM-grep (DC-QA-02):** Die
  `grep -E '^FROM ' | grep -v '^FROM deps' | sort -u`-Pipeline (`Makefile:88`)
  gibt genau die drei externen, jetzt digest-tragenden `FROM`-Zeilen aus; die
  fünf internen `FROM deps`-Stages werden korrekt gefiltert. Der `@sha256:`-Teil
  steht in jeder Zeile — der Digest-Beleg ist vollständig. (Der Tag-Teil zeigt
  literal `${GO_VERSION}`/`${GOLANGCI_LINT_VERSION}`, da `make` Dockerfile-ARGs
  nicht expandiert; die aufgelösten Versionen stehen separat auf `Makefile:86–87`
  — Lesbarkeitsdetail, kein Beleg-Verlust, und vorbestehendes Verhalten der
  Zeile, nicht durch diesen Slice eingeführt.)
- **`grep -oP` semgrep-Extraktion — Mehrfachtreffer/Quote-Syntax:** Auf einem
  GNU-grep-Host extrahiert `SEMGREP_VERSION:-\K[^}]+` → `1.167.0` und
  `SEMGREP_DIGEST:-\K[^}]+` → `sha256:06938c1f…` jeweils **genau einen** Treffer;
  `[^}]+` endet sauber am schließenden `}` des `${…}`, die Werte sind unquoted
  (kein Quote-Stripping nötig), und `head -1` deckt einen künftigen
  Zweittreffer ab. Korrektheit auf GNU-Host bestätigt; die Host-Portabilität ist
  LOW-1.
- **Build-Hermetik / `--network none` (DC-QA-03):** Der Digest betrifft
  ausschließlich die Image-Auflösung beim Pull (Setup/Netz, wie der Tag zuvor).
  Der semgrep-Scan-Container trägt unverändert `--network none`
  (`tools/semgrep.sh:54`); die digest-gepinnte Referenz steht in der
  `docker run`-Zeile (`:57`), das Pull-Setup liegt davor. Kein
  digest-gepinntes `FROM` führt Netzzugriff in eine netzlose Stage ein —
  ADR-0010 Punkt 3 / ADR-0011 §Konsequenzen bleiben erfüllt.
- **Doku-Kohärenz ADR-0011 ↔ Code:** ADR-0011 §1/§2 (alle `FROM` + semgrep,
  Tag+Digest inline) deckt sich mit `Dockerfile:34/55/113` und
  `tools/semgrep.sh` Zeilen 23/57; §5 (Beleg via `make versions`) deckt sich mit
  `Makefile:88–89`; §4 (Doppel-Edit-Politik) deckt sich mit den Kommentaren in
  `Dockerfile:22–27` und `tools/semgrep.sh` Zeilen 21–22. Index-Zeile
  (`docs/plan/adr/README.md:25`) trägt Status Accepted/2026-06-20 und die Bezüge
  ADR-0002/0006/0010/`DC-QA-02` — konsistent mit dem ADR-Kopf. Slice-Plan-DoD
  (`…/slice-033-…:46–59`) nennt inline-Digest, semgrep `${SEMGREP_DIGEST}`,
  neue ADR-0011 statt Edit, Manifest-Listen-Digest — kohärent zum Diff.
- **Scope-Disziplin (slice-033 = nur Digest-Pins):** `git diff HEAD --stat`
  listet ausschließlich Dockerfile, Makefile, ADR-README, Slice-Plan,
  semgrep.sh (+ untrackte ADR-0011, + gestageter Slice-Rename). Keine
  `release.yml`/`.github`-, keine `spec/`-, keine `CHANGELOG`-Änderung. Die
  `latest`-Drift ist sauber ausgeklammert (INFO-1).
- **Lifecycle / AGENTS.md §3.3 (git mv + Inhalt = zwei Commits):** Der
  Slice-Plan ist zugleich umbenannt (`next/`→`in-progress/`, von Git als
  R69%-Rename erkannt) **und** inhaltlich geändert. §3.3 verlangt für die
  *Closure* zwei getrennte Commits (Move, dann Inhalt); im Working-Tree-Diff ist
  das noch nicht beobachtbar (uncommitted). Hinweis an die Implementation/
  Verifikation: Commit-Aufteilung beim Abschluss prüfen — als Working-Tree-Diff
  kein Finding, nur Verifier-relevant.
- **Stilles Grün aus slice-032 (HIGH-1) unberührt:** Der „Ran N rules"-Nachweis
  (`tools/semgrep.sh` Zeilen 70–73) bleibt im Diff unangetastet; die Digest-Ergänzung
  ändert nur die Image-Referenz (`:23`/`:57`) und schwächt die Regel-Nachweis-
  Schranke nicht.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 1 | 2 | 1 |

## Verdikt

**Vor Merge zu klären** (MEDIUM-1). Die Umsetzung ist im Kern sauber: ADR-0002 und
ADR-0010 sind nachweislich byte-unverändert (Vereinheitlichung korrekt über die
neue ADR-0011 statt Edit — §3.5 eingehalten), die Build-Hermetik bleibt erhalten
(Digest betrifft nur den Pull, der Scan bleibt `--network none`), `make versions`
belegt die Digests vollständig über den `FROM`-grep, und die Scope-Abgrenzung zur
`latest`-Drift ist explizit und dicht (release.yml unangetastet). Es bleibt
**MEDIUM-1**: ADR-0011 §3 macht den Manifest-Listen-Digest verbindlich, doch der
Diff belegt nicht, dass die vier gepinnten Digests Index-Digests und nicht
plattformspezifisch sind — ein plattformgebundener Digest bräche den arm64-Build
und damit die zugesicherte Multi-Arch-Reproduzierbarkeit (`DC-QA-02`). Das ist per
`buildx imagetools inspect`/`skopeo inspect --raw` je Digest zu klären; mangels
Netz-/Registry-Zugriff konnte der Reviewer den Manifest-Typ nicht selbst
bestätigen. LOW-1 (PCRE-`grep -oP` → stiller leerer Beleg auf Nicht-GNU-Host) und
LOW-2 (Tag↔Digest ohne Konsistenz-Prüfung) sind latente Wartungsfallen; INFO-1
(`latest`-Drift) ist ein bewusst delegierter Folgepunkt. Die Gate-/Build-
Bestätigung obliegt der getrennten Verifikation (hier NICHT als grün angenommen;
Tests/Gates/Registry-Auflösungen wurden nicht ausgeführt).
