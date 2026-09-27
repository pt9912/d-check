# Review-Report — slice-234 (`--suggest-config ai-harness` erkennt `RB`), R1

**Review-Art:** Code (gegen Slice-Plan, betroffene `DC-*`-Anforderung, Hard Rules).

**Gegenstand:** `slice-234`, Commit-Range `d2bcd698..5197346e`
(fünf Commits: `09bd9f5c`, `ed6e7d14`, `06fc8d9c`, `c4ed0e9e`, `5197346e`).

**Skill:** `.harness/skills/reviewer.md` @ `8dcc0452` (Version 1.16.0, 2026-09-07).

**Modell-ID:** `claude-sonnet-5`.

**Datum:** 2026-09-27.

**Eingangs-Kontext:**

- Slice-Plan `docs/plan/planning/in-progress/slice-234-suggest-config-kennungsreihe-rb.md`
  (§1–§8, vollständig gelesen).
- `DC-FA-CLI-006` (`spec/lastenheft.md`), Spec-Stelle `DC-FA-CLI-006.a`
  (`spec/spezifikation.md`).
- `AGENTS.md` §3 (§3.7 Kommentar-Klassen, §3.8 Modul-Scan-Grenze), §5
  (Grenzen-Regel, MR-025, Zitat-Geltungsbereich, Commit-Overclaim).
- Kein ADR für diesen Slice — bewusste Entscheidung, geprüft unten (Finding
  R1-M3).
- `docs/reviews/` durchsucht nach früheren Reviews zu `suggest.go`/
  `DC-FA-CLI-006`/`--suggest-config`/slice-037: **keine eigenständige
  Review-Datei gefunden** (slice-037 liegt archiviert unter
  `docs/plan/planning/done/welle-26/`; sein Review-Report — falls vorhanden —
  ist Teil des Wellen-Archivs und wird laut Kanon nicht erneut gelesen). Als
  Ersatz-Signal geprüft: die Test-Kommentare `TestCLI037_IDPrefix_*` nennen
  zwei historische R1-MEDIUM-Findings (Override-Reihenfolge, Flag übergeht
  Konflikt) — beide bereits gefixt und getestet, von diesem Diff nicht berührt.
- **Nicht** erhalten (laut Skill): die DoD-Abhakung selbst.

**Gates unabhängig nachgefahren** (Docker/`make`, kein Host-Go):

| Gate | Ergebnis |
|---|---|
| `make test` | grün — alle Pakete `ok`, inkl. `internal/adapter/driving/cli` |
| `make lint` | grün — `golangci-lint run ./...`: `0 issues.` |
| `make doc-check` | grün — `835 Datei(en) geprüft, 0 Befund(e)` |
| `make coverage-gate` | grün — `94.30%` (Schwelle `93%`) |
| `make gates` (voll, alle zehn Glieder) | grün — `baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` |

---

## Zentraler Prüfpunkt: Byte-Gleichheit vs. bedingte Aktivierung

**Ergebnis: die Kern-Logik ist korrekt implementiert und durch Tests belegt.**
Nachvollzogen am Code (`internal/hexagon/core/app/suggest.go`):

- `sawRB` wird ausschließlich innerhalb `deriveReqPrefix` gesetzt, und
  `deriveReqPrefix` wird ausschließlich aufgerufen, wenn
  `reqPrefix == "" && harness` — `harness` ist nur dann `true`, wenn die
  Quelle `ai-harness` (nicht `ai-harness-init`) unter den `sources` ist. Für
  reines `--id-prefix` oder reines `ai-harness-init` bleibt `sawRB` immer
  `false`, `harnessIDPatterns` erhält `includeRB=false`, und das erzeugte
  Muster ist bit-identisch zum Vor-Slice-Stand (`FA-[A-Z]+|QA`, ohne `|RB`).
  Das deckt sich exakt mit der in Commit `06fc8d9c` nachgezogenen „benannten
  Grenze" in `spec/spezifikation.md`.
- Die Ambiguitäts-Konsequenz (ein zweites Präfix, das nur über eine `-RB-`
  Überschrift auftritt, kippt einen vorher stillen Ein-Präfix-Erfolg in einen
  Exit-2-Fehler) **liegt außerhalb** des Geltungsbereichs der
  Byte-Gleichheits-Zusage — die ist in Lastenheft/Spezifikation explizit auf
  „ohne `-RB-`-Überschrift im Repo" verengt. Diese Konsequenz ist zusätzlich
  **dreifach benannt**: Slice-Plan §6 (Risiko, Ausgang „bei Closure zu
  vergeben"), Lastenheft-AK „Randbedingungs-Reihe (RB) Mehrdeutigkeit", und
  die 0.91.0-Historie-Zeile. Kein Verschweigen.
- Regex-Grenzfall `AC-RBX-01` von Hand durchgerechnet: `reqShape`
  (`^([A-Z][A-Z0-9]*)-(?:FA-[A-Z]+|QA|RB)-\d+[A-Za-z]?$`) kann die
  Capture-Gruppe nur bis zum ersten `-` fassen (kein `-` im Zeichensatz), das
  ergibt Präfix `AC`, Rest `RBX-01` — keine der drei Alternativen
  (`FA-[A-Z]+`/`QA`/`RB`) matcht `RBX-01` (nach `RB` verlangt die Regex
  sofort `-`, tatsächlich folgt `X`). `AC-RBX-01` matcht `reqShape` also
  **gar nicht** — bestätigt am Regex, nicht nur am Kommentar. Die
  `sawRB`-Markierung (`strings.HasPrefix(tok[len(m[1])+1:], "RB-")`) wird
  dadurch nie mit einem falschen Positiv erreicht; ein Off-by-one/Panik-Risiko
  besteht nicht, da `tok` bereits vollständig gegen `reqShape` gematcht ist.
- `reqIDFull` (`trace.go`) und das statische Gerüst
  (`internal/hexagon/core/app/config_template.go`) sind laut Diff
  **unverändert** (`git diff d2bcd698 5197346e -- trace.go config_template.go`
  liefert keine Ausgabe) — deckt sich mit der Abgrenzung in Slice-Plan §1/§3.

Zwei Findings unten (R1-M1, R1-M2) betreffen **nicht** diese Kern-Logik,
sondern zwei begleitende Vertragsstellen, die dem sonst sauberen Kern nicht
gerecht werden.

---

## Findings

| # | Kategorie | Quelle | Pfad | Befund | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| R1-M1 | MEDIUM | `DC-FA-CLI-006.a` | `spec/spezifikation.md` (Abschnitt „Kanonische Vorlage", Zeile mit `regex: '<PREFIX>-(FA-[A-Z]+\|QA\|RB)-\d+'`) | Die als „Spiegel der Repo-Konvention; `ai-harness-init` gibt sie vollständig aktiv aus" deklarierte Vorlage zeigt `RB` **unbedingt** in der Alternation — im Widerspruch zu (a) der unmittelbar darüber stehenden Prosa („ein Repo, das … den Voll-Kanon `ai-harness-init` nutzt, liest `spec/lastenheft.md` dafür nicht — `RB` bleibt dort außen vor"), (b) dem eigenen YAML-Kommentar zwei Zeilen darüber („`\|RB` nur, wenn derselbe Ableitungs-Durchlauf eine `-RB-`-Überschrift sah") und (c) dem tatsächlichen Code (`harness`-Flag ist bei reinem `ai-harness-init` immer `false`, `sawRB` bleibt `false`, das reale Muster enthält kein `RB`). Ein Leser (namentlich der CR-Konsument `ai-harness-course`, der diese Vorlage referenziert), der die Vorlage wörtlich übernimmt, bekäme `RB` unbedingt in jedes frische Repo eingebacken — genau das Verhalten, das Commit `06fc8d9c` bewusst ausgeschlossen hat. | ja — `d-check --suggest-config ai-harness-init` gegen ein leeres Repo ausführen und die `ids`-Anforderungs-Zeile gegen den Fenced-Block in `spec/spezifikation.md` diffen | spec-example-drift |
| R1-M2 | MEDIUM | `DC-FA-CLI-006` (Lastenheft-AK „Randbedingungs-Reihe (RB) Happy") | `internal/adapter/driving/cli/cli_acceptance_test.go:~1300` (`TestCLI234_RB_Happy`) | Das Lastenheft-Akzeptanzkriterium und die Slice-Plan-DoD (§2: „fehlende `-RB-07` wird gemeldet — Rot-Beleg gegen den alten Stand") versprechen zusätzlich zur Muster-Erweiterung, dass ein anschließender Lauf gegen das vorgeschlagene Gerüst eine nicht aufgelöste `<PREFIX>-RB-07` als `id-unlinked` meldet. `TestCLI234_RB_Happy` prüft ausschließlich, dass das erzeugte Regex-Objekt `AC-RB-01` matcht und `FA`/`QA` nicht verliert — es gibt keinen Test, der das vorgeschlagene Muster tatsächlich in einen `ids`-Lauf einspeist und den `id-unlinked`-Befund für die fehlende `-RB-07` beobachtet. Die zweite Hälfte des öffentlichen Vertrags ist unbelegt. | ja — Suche nach einem Test, der `--suggest-config`-Ausgabe in eine `.d-check.yml` überführt und gegen ein Lastenheft mit unaufgelöster `-RB-07` laufen lässt; keiner existiert | ak-half-untested |
| R1-M3 | MEDIUM | Governance/Modul 4 (ADR) | `spec/lastenheft.md` Historie 0.91.0; `spec/spezifikation.md` Historie-Zeile | Die dokumentierte „Keine begleitende ADR"-Begründung adressiert ausschließlich die RTM-Abgrenzungsfrage (Fachfrage, keine Architekturentscheidung). Sie äußert sich **nicht** zur tatsächlich nicht-trivialen Entscheidung dieses Slice — unbedingte vs. bedingte Aufnahme von `RB` ins Muster, um die Byte-Gleichheit zu erhalten —, obwohl genau diese Entscheidung im ersten Anlauf (Commit `ed6e7d14`) falsch getroffen und in Commit `06fc8d9c` korrigiert werden musste. Zwei directe Präzedenzfälle vergleichbaren Zuschnitts (`DC-FA-MTX-003` Instanz-Identitäts-Ausnahme, `DC-FA-STRUCT-001` `max-lines`) haben für ähnlich schmale Muster-Erweiterungen jeweils eine begleitende ADR bekommen. Die „Verglichene Alternativen"-Funktion einer ADR (warum nicht unbedingt, warum nicht per Flag) ist hier nur aus der Commit-Botschaft von `06fc8d9c` rekonstruierbar, nicht aus einem stehenden Artefakt. | teilweise — Urteilsfrage, kein Gate; verifizierbar durch Vergleich der Präzedenzfälle | missing-adr-for-corrected-design |
| R1-L1 | LOW | Slice-Plan §2 (Prozess) | `docs/plan/planning/in-progress/slice-234-suggest-config-kennungsreihe-rb.md` §2/§7 | Die DoD verlangt, die ADR-Notwendigkeits-Entscheidung „vor dem ersten Edit" zu treffen und „in §7 [zu] benennen". §7 (Closure-Notiz) ist zum Review-Zeitpunkt noch leer (`*(Bei der Closure zu füllen.)*`); die tatsächliche Begründung landete stattdessen in der Lastenheft-Historie und der Commit-Botschaft von `06fc8d9c`. Inhaltlich ist die Begründung vorhanden, nur nicht am selbst zugesagten Ort. | ja — Diff von §7 zum Review-Zeitpunkt | dod-self-commitment-location-drift |
| R1-I1 | INFO | Repo-Konvention (nicht formalisiert, `MR-035` deckt nur ausgehende CRs) | `docs/plan/cr/` | Zwölf vorangehende **eingehende** CRs (u. a. derselbe Tag, derselbe Konsument: `2026-09-27-cr-eingehend-ai-harness-course-structure-max-lines.md`) liegen als datierte Datei unter `docs/plan/cr/`. Der in diesem Slice referenzierte CR („Change Request des Konsumenten `ai-harness-course` … 2026-09-27") hat kein Gegenstück dort — nur Prosa im Slice-Kopf und in der Lastenheft-Historie. Kein Hard-Rule-Verstoß (der Kanon erklärt einen eingehenden CR ausdrücklich zu „bewusst keinem Harness-Konstrukt", `MR-035` bindet nur ausgehende CRs), aber eine Abweichung von einer bislang ausnahmslosen Bestandspraxis. | ja — `ls docs/plan/cr/` zeigt die Lücke | incoming-cr-not-filed |

## Negativbefunde (geprüft, ohne Befund)

- **Gates:** `make test`, `make lint`, `make doc-check`, `make coverage-gate`
  und der volle `make gates`-Lauf (zehn Glieder) unabhängig nachgefahren —
  alle grün, echte Ausgabe oben zitiert.
- **Hexagon-Import-Richtung (ADR-0005):** `suggest.go` importiert nur
  `core/rules`, `core/model`, `port/driven` — keine neue Abhängigkeit, keine
  Verletzung.
- **Gate-Suppression / Schwellen-Senkung ohne ADR (§3.2/§3.6):** keine
  `//nolint`, keine Schwellen-Änderung im Diff.
- **Netzzugriff außerhalb `external`:** keiner — reiner Lese-Pfad auf
  Filesystem-Abstraktion.
- **Kommentar-Klassen (§3.7):** die neuen/geänderten Kommentare an
  `reqShape`, `deriveReqPrefix`, `harnessIDPatterns` tragen Zusage- und
  Grenz-Aussagen über das erzeugte Muster (warum bedingt statt unbedingt) —
  keine Review-Historie, keine Deliberation über Verworfenes, keine
  Slice-Nummer im Code-Kommentar (die Kennung `DC-FA-CLI-006` ist ein
  auflösbares Rang-Zeiger-Feld).
- **Zustandsfelder:** keine Roadmap-/Register-Zustandszeile im Diff außer der
  bereits etablierten Lifecycle-Move-Mechanik (`Nichts in Arbeit` entfernt,
  reiner `git mv`-Begleiteffekt).
- **Commit-Zerlegung:** Commit `5197346e` (Code) ist rein `internal/`
  (`cli_acceptance_test.go`, `suggest.go`); Commit `06fc8d9c` (Korrektur) ist
  rein `spec/`, fasst keinen Code an; die Lifecycle-Move-Commits
  (`09bd9f5c`, `d2bcd698`, `c4ed0e9e`) sind reine Planning-Doku-Commits.
- **Spec-Straten-Referenzrichtung (§3.4/MR-006):** weder Lastenheft noch
  Spezifikation nennen in diesem Diff eine ADR, einen Slice oder einen
  Commit-Hash im Körper — `matrix`-Modul bestätigt das über `make doc-check`.
- **RTM/`reqIDFull` (Out-of-Scope):** `trace.go` und `config_template.go` im
  Diff unverändert — Abgrenzung eingehalten.
- **Byte-Gleichheits-Kern (`sawRB`-Gating):** wie oben ausgeführt, am Code
  nachvollzogen und durch `TestCLI234_RB_AbwesendByteGleich` (Byte-exakter
  Vergleich gegen die alte Formel `AC-(FA-[A-Z]+|QA)-\d+`) belegt — kein
  Finding.
- **Regex-Grenzfall `AC-RBX-01`:** von Hand durchgerechnet, matcht `reqShape`
  nicht — kein falsches Positiv in der `sawRB`-Erkennung.
- **Test-Authentizität:** `reqPattern()` dekodiert die reale Stdout-Ausgabe
  über den echten `configyaml`-Decoder und liefert das reale `regexp.Regexp`
  — keine Tautologie, alle drei neuen Tests laufen durch den vollen
  CLI-Pfad.
- **Referenz-Richtung/Provenance-Marker (§Reviewer-Anker MEDIUM):** kein
  `<!-- d-check:status-provenance -->`-Marker im Diff, keine neuen
  Abwärts-Token in Spec-Straten.
- **`DC-QA-04` Alt-Tool-Migrationsabdeckung:** nicht berührt — kein
  Migrations-Modul im Diff.
- **Modul-Scan-Grenze (§3.8):** `deriveReqPrefix` liest `spec/lastenheft.md`
  außerhalb der `.d-check.yml`-Scan-Achse (es ist kein Audit-Modul, sondern
  ein CLI-Generator) — die Grenze ist in `spec/spezifikation.md` explizit als
  „benannte Grenze" benannt (welche Modi lesen, welche nicht), nicht
  stillschweigend.

## Kategorie-Summary

- HIGH: 0
- MEDIUM: 3 (R1-M1, R1-M2, R1-M3)
- LOW: 1 (R1-L1)
- INFO: 1 (R1-I1)

## Verdikt

**Merge-blockierend: ja, wegen R1-M1 und R1-M2.**

Begründung: R1-M1 ist eine in sich widersprüchliche Spezifikationsstelle
innerhalb derselben `DC-FA-CLI-006.a`-Passage (Prosa sagt „RB bleibt bei
`ai-harness-init` außen vor", die direkt danebenstehende „kanonische"
Beispiel-Ausgabe zeigt `RB` unbedingt) — genau die Art Drift, die
`--suggest-config`s Konsumenten (hier der CR-Absender selbst) in die Irre
führt, wenn sie die Vorlage wörtlich übernehmen. R1-M2 lässt die stärkere
Hälfte des neuen öffentlichen Vertrags (End-zu-Ende-`id-unlinked`-Meldung)
ungetestet, obwohl DoD und Lastenheft explizit einen Rot-Beleg dafür
zusagen. Beide sind mit überschaubarem Aufwand zu schließen (R1-M1: die
Fenced-Vorlage entweder ohne `RB` zeigen oder den Bedingungssatz auch dort
abbilden; R1-M2: ein Test, der das vorgeschlagene Muster real gegen ein
Lastenheft mit unaufgelöster `-RB-07` laufen lässt) und sollten vor Closure
behoben werden. R1-M3 ist eine Governance-Ermessensfrage (ob eine ADR fällig
gewesen wäre) und blockiert für sich allein nicht zwingend, sollte aber vor
Closure explizit entschieden (nicht nur implizit belassen) werden. R1-L1 und
R1-I1 sind nicht merge-blockierend.
