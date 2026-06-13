# Slice slice-010: Image-Integrationstests + Reproduzierbarkeits-Belege

**Status:** done.

**Welle:** welle-04-distribution-und-migration.

**Bezug:** [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(Akzeptanzkriterien-Teil: identisches Verhalten Container vs. nativ),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(byte-identische Ausgabe als Vergleichskriterium),
[ADR-0002](../../adr/0002-distribution-ghcr-image.md) (Image-Form);
Kurs-Modul 14 (Reproduzierbarkeits-Bindung, Image-Hash).

**Autor:** pt9912. **Datum:** 2026-06-11.

---

## 1. Ziel

Die DIST-Akzeptanzkriterien sind gegen das **lokal gebaute** Image
automatisiert (kein externer Abhängigkeits-Anteil), und die
Reproduzierbarkeits-Targets `versions`/`fullbuild`/`ci` existieren —
die „Nicht behauptet"-Listen in AGENTS/harness sind danach leer.

## 2. Definition of Done

- [x] Image-Integrationstest (`tools/image-test.sh`, Make-Target
  `image-test`): die drei Akzeptanzkriterien von
  [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
  automatisiert — (1) Happy: Repo mit kaputtem Link → Befund-Ausgabe
  und Exit-Code des Containers **byte-identisch** zur nativen
  Ausführung (natives Binary aus der compile-Stage, Docker-only-Regel
  bleibt gewahrt); (2) Boundary: read-only-Mount (`:ro`) →
  vollständige Prüfung ohne Schreibfehler; (3) Negative: fehlender
  `/repo`-Mount → Exit 2 mit Mount-Hinweis.
- [x] `make versions`: gibt alle Pins reproduzierbar aus
  (`GO_VERSION`, `GOLANGCI_LINT_VERSION`, Basis-Images, Image-Digest
  des Runtime-Builds) — Grundlage der Reproduzierbarkeits-Bindung.
- [x] `make fullbuild`: volle Closure (gates + image-test + bench),
  schließt mit dem Image-Hash (`sha256:…`) ab; der Hash ist die
  Reproduzierbarkeits-Bindung in der Sensors-Tabelle (Kurs-Modul 14).
- [x] `make ci`: CI-äquivalenter Lauf (gates + image-test) — das
  Target, das die Release-Pipeline (slice-011) aufruft.
- [x] Sensors-Tabelle ([`harness/README.md`](../../../../harness/README.md))
  und [`AGENTS.md`](../../../../AGENTS.md) §4 um die neuen Targets
  ergänzt; **„Nicht behauptet"-Listen leer** (`make gate-consistency`
  erzwingt die Konsistenz in beide Richtungen).
- [x] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md);
  Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/image-test.sh` | neu | DIST-Akzeptanzkriterien als Skript (Fixture, Vergleich nativ/Container) |
| [`Dockerfile`](../../../../Dockerfile) | update | ggf. Binary-Export für den Nativ-Vergleich (`-o`-Stage) |
| [`Makefile`](../../../../Makefile) | update | Targets `image-test`, `versions`, `fullbuild`, `ci` |
| [`harness/README.md`](../../../../harness/README.md), [`AGENTS.md`](../../../../AGENTS.md) | update | Sensors-/Gates-Tabellen, „Nicht behauptet" leeren |

## 4. Trigger

Sofort — welle-04 ist aktiv; alles ist offline beweisbar (keine
externen Abhängigkeiten).

## 5. Closure-Trigger

DoD vollständig + Commit(s) auf `main` + Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- „Nativ" in einem Docker-only-Repo: der Vergleichslauf nutzt das
  statische Binary aus der compile-/build-Stage in einem schlanken
  Container — kein Host-Go (AGENTS.md §3.1 bleibt unverletzt).
- `fullbuild` läuft Minuten (bench + image-test); es ist bewusst
  **nicht** Teil von `make gates` (inner loop bleibt schnell), sondern
  Wellen-/Release-Schlusspunkt.
- Byte-Identität setzt deterministische Ausgabe voraus
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus))
  — relative Pfade sind bereits Vertragsbestandteil.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Commit `c5f8a2a`. **Belege:** `make ci` grün (sechs
Gates + image-test); `make fullbuild` grün — image-test OK,
Benchmark-Median 526 ms, Abschluss mit Image-Hash
`sha256:5a3a6a91…` (Reproduzierbarkeits-Bindung).

- **Was hat funktioniert:** Die drei DIST-Akzeptanzkriterien ließen
  sich als ein Shell-Skript mit einem Fixture abbilden; der
  byte-identische Vergleich (stdout **und** stderr) ist dank der
  relativen Pfade aus [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) trivial erfüllbar gewesen.
- **Anders als geplant:** Das „native" Binary kommt per `docker cp`
  aus dem **Runtime-Image** statt aus der compile-Stage (Plan §3
  hatte eine mögliche Export-Stage vorgesehen) — so vergleicht der
  Test exakt das Artefakt, das ausgeliefert wird, und das Dockerfile
  blieb unverändert.
- **Steering-Loop-Lerneintrag:** `gate-consistency` forderte die vier
  neuen Targets unmittelbar nach ihrer Entstehung ein (Erstlauf rot,
  vier Meldungen) — das Meta-Gate aus slice-009 wirkt bereits als
  Promotion-Trigger-Wächter: ein Target kann nicht mehr entstehen,
  ohne dokumentiert zu werden. Zweite Beobachtung: eine
  Gate-Verkettung über `make`-Prerequisites (`ci: gates image-test`)
  braucht keinerlei neue Logik — die `.NOTPARALLEL`-Disziplin aus
  [`MR-005`](../../../../harness/conventions.md#mr-005--härtung-ggü-b-cad-inhaltsbasierter-gate-nachweis-sub-shell-prüfung)
  trägt unverändert.
- **Folge-Slices:** keine neuen; slice-011 (GHCR-Release) ist durch
  `make ci` entriegelt (Trigger erfüllt).

**Review R1 (gebündelt über slice-009+010, Agent-Review mit
getrenntem Kontext):** 9 Findings (2 MEDIUM, 4 LOW, 3 INFO), alle
behoben — Kern: die einzige Stilles-Grün-Falle war die hart
verdrahtete Median-Position im Benchmark (latent, hätte erst bei
RUNS-Änderung gezündet); dazu die fehlende 2-vCPU-Normierung
(Spec-Treue der Messmethode), eine Mehrfach-Target-Lücke im
Meta-Gate-Parser (jetzt mit Parser-Selbsttest), `fullbuild: ci bench`
statt Kettenduplikat, `versions` ohne Stage-Rauschen, drei
dokumentierte Annahmen (QA-03-Format fail-closed, amd64,
Fixture-Verbleib). Kein akuter Stilles-Grün-Pfad im ausgelieferten
Stand — die Review-Beute waren latente Wartungsfallen und
Spec-Treue, nicht aktive Bugs.

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Gate-/Tooling-Arbeit; siehe Kurs Modul 5
§Worked Mini-Example).
