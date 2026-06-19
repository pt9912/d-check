# Slice slice-030: `--suggest-config ai-harness`-Vorlage

**Status:** done (abgeschlossen 2026-06-19).

**Welle:** welle-19-suggest-ai-harness (Trigger: Change Request 0.18.0
akzeptiert; baut auf slice-020 (`--suggest-config`) auf).

**Bezug:**
[`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(Hauptanforderung — reservierte Quelle `ai-harness`),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(Ergebnis dekodiert über den eigenen Parser),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(deterministisch),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(read-only).

**Autor:** pt9912. **Datum:** 2026-06-19.

---

## 1. Ziel

`d-check --suggest-config <token>` gibt ein an die adoptierte
ai-harness-course-Konvention (Baseline `v1.3.0`) angelehntes,
kommentiertes `.d-check.yml` auf stdout aus: die **strukturellen**
Konventionen, die reine Quellen-Ableitung nicht erzeugen kann (kanonische
`ids`-Muster, `matrix`-Klassen samt Referenzrichtung, Standard-Modulset,
Scan-Scope). **Zwei Modi** (Henne-Ei — nicht aus der Repo-Existenz
ableitbar, daher explizit gewählt):

- **`ai-harness-init`** (Voll-Kanon): alle Blöcke aktiv, ohne
  Existenzprüfung — Zielbild fürs leere/frische Repo; läuft, sobald die
  Struktur existiert.
- **`ai-harness`** (repo-bewusst): nur existierende Pfade/Targets aktiv,
  fehlende auskommentiert mit Hinweis — läuft sofort gegen ein bestehendes
  Repo.

Read-only, advisory, deterministisch; mit echten Quellen kombinierbar.

## 2. Definition of Done

- [x] **CLI/Core:** `ai-harness` und `ai-harness-init` werden in
  `SuggestConfig` als **reservierte Schlüsselwörter** erkannt (nicht als
  Pfad → **kein** Exit 2 „Quelle existiert nicht"); kombinierbar mit
  echten Quellen.
- [x] **Vorlage (zwei Modi):** kanonische Blöcke aus d-checks
  `.d-check.yml` — `ids` (`ADR-\d{4}`→`docs/plan/adr/`,
  `MR-\d{3}`→`harness/conventions.md`,
  `DC-(FA-[A-Z]+|QA)-\d+`→`spec/lastenheft.md`,
  `slice-\d{3}`→`docs/plan/planning/`; `link-policy: always`,
  `exempt-paths: [CHANGELOG.md, "docs/reviews/**"]`), `matrix`
  (Klassen spec-straten/adr/slice, Regeln `spec-straten`→adr/slice =
  false, `status.forbidden`, `exclude-sections`), `modules`
  `[links, anchors, ids, matrix, codepaths]`, `scan.roots` +
  `ids.scope`. **`ai-harness-init`:** alle Blöcke aktiv (Voll-Kanon, keine
  Existenzprüfung). **`ai-harness`:** repo-bewusst — fehlendes Target/Pfad
  → **auskommentiert mit Hinweis** statt weglassen. Carveout (kein Target)
  in beiden auskommentiert. Kommentar-Header mit Baseline-Pin + Modus.
- [x] **Spezifikation:**
  [`DC-FA-CLI-006.a`](../../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag)
  um die beiden Modi erweitert (Schritte + Vorlagen-Inhalt +
  Baseline-Pin `v1.3.0`).
- [x] **Tests:** die CR-Akzeptanzkriterien — `ai-harness` (Happy: Blöcke
  vorhanden/Parser akzeptiert; Boundary: fehlendes `docs/plan/adr/` →
  ADR-Block auskommentiert, nur existierende roots; Abgrenzung: ≠ fehlende
  Quelle, Exit 0) und `ai-harness-init` (Voll-Kanon: alle Blöcke aktiv
  auch ohne Targets); Determinismus 10× ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus));
  read-only.
- [x] **Doku:** `docs/user/operations.md` (Options-Tabelle) +
  Benutzerhandbuch (§ zu `--suggest-config`) um den Modus ergänzt.
- [x] `make gates` grün; unabhängiges Review R1; Closure-Notiz.

## 3. Plan (vor Code)

| Datei | Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | [`DC-FA-CLI-006.a`](../../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag) um `ai-harness`-Modus (Vorlage, Hybrid-Regel, Baseline-Pin) |
| `internal/hexagon/core/suggest.go` | update | reserviertes Schlüsselwort, repo-bewusste Vorlagen-Erzeugung; Wiederverwendung `resolveConfigPath`/`fsys.Kind` |
| `internal/adapter/driving/cli/cli_acceptance_test.go` | update | drei AK + Determinismus + read-only |
| `docs/user/operations.md`, `docs/user/benutzerhandbuch.md` | update | Modus dokumentieren |

Lastenheft ist mit Change Request 0.18.0 gesetzt.

## 4. Trigger

Change Request 0.18.0 akzeptiert (Bezug oben).

## 5. Closure-Trigger

DoD vollständig inkl. grüner Gates und Review R1.

## 6. Risiken und offene Punkte

- **Kanon-Drift:** die Vorlage spiegelt d-checks `.d-check.yml`; weicht
  diese ab, driftet der Modus. Gegenmittel: Pin auf Baseline `v1.3.0` im
  Header + Spezifikation als Quelle der Vorlage; nicht aus der Laufzeit-
  `.d-check.yml` ableiten (sonst zirkulär).
- **Determinismus** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  feste Block-/Pfad-Reihenfolge, keine Map-Iteration; Existenz-Prüfung
  rein aus dem gescannten Baum (kein git/Netz).
- **Parser-Treue** ([`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)):
  der aktive (nicht auskommentierte) Teil muss strikt dekodieren — Test
  führt das Ergebnis durch `configyaml.Decode`.
- **Carveout-Kennungen** ohne festes Definitions-Target (bislang
  ungenutzt) → nur als auskommentierter Hinweis, kein aktives Muster.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Vertrag
[`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(Lastenheft 0.18.1, zwei Modi) + Spezifikation
[`DC-FA-CLI-006.a`](../../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag).
Reservierte Quellen `ai-harness` (repo-bewusst) und `ai-harness-init`
(Voll-Kanon); `renderHarness` mit `repoAware`-Schalter über gemeinsamem
Kanon (`harnessIDPatterns`/`harnessClasses`, Spiegel von `.d-check.yml`).
CLI unverändert (das Schlüsselwort wird im Core erkannt). `make gates` grün
(Coverage 93,70 %).

**Belege:**

- 5 Akzeptanztests: `ai-harness` Happy/Boundary/Abgrenzung,
  `ai-harness-init` Voll-Kanon, Determinismus 10×
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- Beide Modi **am Image verifiziert**: `ai-harness-init` auf leerem Repo →
  Voll-Kanon aktiv; `ai-harness` auf d-check → repo-bewusst (gespiegelte
  `.d-check.yml`).
- read-only ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit))
  per `os.Stat` in beiden Happy-Tests belegt; rein lesender Core-Pfad.

**Lerneintrag:** Henne-Ei — ein einzelner auto-hybrider Modus kommentierte
im leeren Repo *alles* aus (nutzlos als Bootstrap), und ein naives „immer
voll-aktiv" stirbt beim Lauf an `ensureIDTargetsExist` (Exit 2, fehlendes
ids-Target). Auflösung: **zwei explizite, nutzergewählte Modi** statt
Auto-Erkennung. Der Einwand kam vom Auftraggeber *vor* der Closure — der
Review-/Abnahme-Schritt vor Release hat den Designfehler gefangen.

**Review R1** (unabhängiger Reviewer-Subagent, eigener Kontext ohne
DoD-Wissen,
[Report](../../../reviews/2026-06-19-slice-030-suggest-config-ai-harness.md)):
HIGH 0 / MEDIUM 0 / LOW 0 / INFO 2 — freigegeben; beide INFO im selben Stand
geschlossen (Doppel-Token-Vorrang in der Spezifikation dokumentiert;
read-only-Test für den `ai-harness-init`-Pfad ergänzt). Der erste Entwurf
(Einzel-Hybrid) wurde nach dem Henne-Ei-Einwand verworfen und das
Zwei-Modi-Design neu reviewt.

**Welle:** welle-19-suggest-ai-harness ist damit vollständig (slice-030).

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Spec-/Code-/Doku-Arbeit; Greenfield-Default).
