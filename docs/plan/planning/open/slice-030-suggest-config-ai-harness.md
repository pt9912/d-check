# Slice slice-030: `--suggest-config ai-harness`-Vorlage

**Status:** open (geplant).

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

- [ ] **CLI/Core:** `ai-harness` und `ai-harness-init` werden in
  `SuggestConfig` als **reservierte Schlüsselwörter** erkannt (nicht als
  Pfad → **kein** Exit 2 „Quelle existiert nicht"); kombinierbar mit
  echten Quellen.
- [ ] **Vorlage (zwei Modi):** kanonische Blöcke aus d-checks
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
- [ ] **Spezifikation:**
  [`DC-FA-CLI-006.a`](../../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag)
  um die beiden Modi erweitert (Schritte + Vorlagen-Inhalt +
  Baseline-Pin `v1.3.0`).
- [ ] **Tests:** die CR-Akzeptanzkriterien — `ai-harness` (Happy: Blöcke
  vorhanden/Parser akzeptiert; Boundary: fehlendes `docs/plan/adr/` →
  ADR-Block auskommentiert, nur existierende roots; Abgrenzung: ≠ fehlende
  Quelle, Exit 0) und `ai-harness-init` (Voll-Kanon: alle Blöcke aktiv
  auch ohne Targets); Determinismus 10× ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus));
  read-only.
- [ ] **Doku:** `docs/user/operations.md` (Options-Tabelle) +
  Benutzerhandbuch (§ zu `--suggest-config`) um den Modus ergänzt.
- [ ] `make gates` grün; unabhängiges Review R1; Closure-Notiz.

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

*(wird bei Closure gefüllt — Umsetzung, Belege, Lerneintrag, Review-Runde.)*

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Spec-/Code-/Doku-Arbeit; Greenfield-Default).
