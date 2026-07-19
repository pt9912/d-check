# Slice slice-080: Modul `sources` — Upstream-Content-Drift externer Quellen

**Status:** Done (Release **v0.51.0**, 2026-07-19)
**Welle:** welle-63-sources (Trigger: WIP-Slot frei nach welle-62-Abschluss, v0.50.0)
**Bezug:** [`DC-FA-SRC-001`](../../../../spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz), [ADR-0046](../../adr/0046-sources-upstream-content-drift.md), [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (Amendment); Vorläufer als Harness-Tooling [`MR-022`](../../../../harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019)
**Autor:** pt9912
**Datum:** 2026-07-19

---

## 1. Ziel

Das 19. Regelmodul **`sources`** (opt-in, **Netz**) erkennt **Upstream-Content-Drift**:
eine auf einen `sha256` gepinnte **externe** Quelle wird geholt, gehasht und gegen
den Pin verglichen — driftet der Inhalt, `source-drift`; ist die Quelle nicht
erreichbar, `source-unreachable`. Es ist das Netz-/Inhalts-Gegenstück zu `pins`
(in-repo Content-Hash) und `external` (Netz-Erreichbarkeit) und **produktisiert** das
Kurs-Beispiel `check_regelwerk_drift.py` als reusables Modul für jeden Adopter.
d-checks eigene Baseline deckt der Bash-Helfer `fetch-baseline-cache.sh --check-latest`
([`MR-022`](../../../../harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019))
schon ab; dieses Modul macht es allgemein.

## 2. Entscheidungen

- **Zwei Deklarations-Flächen (beide):** Marker `<!-- source-pin: [zip] sha256:<hex> -->`
  am externen `http(s)`-Link (dpin-Stil, per-Referenz) **und** Config-Block `sources: [{url, sha256, unpack}]`.
- **Zwei Quelltypen (beide):** Einzeldatei (Roh-Byte-`sha256`) und Archiv
  (`unpack: zip` → `sha256` eines **pfad-sortierten**, byte-genau definierten
  **Content-Manifests**, nicht der Zip-Roh-Bytes → reihenfolge-invariant;
  konzeptionell wie `SHA256SUMS`, eigenständig kanonisiert).
- **Archiv-Keyword explizit `zip`** (parallel zu `unpack: zip`; kein Ableiten aus der `.zip`-Endung).
- **Robustheit:** größenbegrenzter Fetch (Body ≤ 64 MiB, Entpack ≤ 256 MiB/≤ 10 000 Einträge, Zip-Bomben-Schutz), Redirects wie `external` (bis fünf), `sha256` case-insensitiv; Nicht-Zip-2xx → `source-unreachable`.
- **`source-drift` emittiert den vollen Ist-`sha256`** (Re-Pin-Vorlage) — schließt die
  dpin-Ergonomie-Sackgasse für dieses Modul nativ.
- **`source-unreachable` getrennt von `source-drift`** (unerreichbar ≠ gedriftet).
- **Amendment der Netz-Sparsamkeit:** Netz jetzt in `external` **und** `sources` (beide opt-in, nie
  Default); der getippte Netzlos-Modullisten-Test führt `sources` als zweite Netz-Ausnahme.
- **Abgrenzung:** Currency/„neuerer Tag" bleibt im Bash-Helfer `--check-latest`; der
  `pins`/dpin-Hash-Ergonomie-Fix bleibt **separat** ([`slice-072`](../in-progress/slice-072-handbuch-aufgabenorientierung.md)) — dieser Slice schneidet nicht quer.
- Begründung/Alternativen in [ADR-0046](../../adr/0046-sources-upstream-content-drift.md).

## 3. Definition of Done

- [x] **Spec/Doc (doc-first):** [`DC-FA-SRC-001`](../../../../spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz) (Lastenheft 0.49.0) + Bereich `SRC` §3 + [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Erweiterung + §7 + `sources` in der Regelmodul-Auswahl + Glossar; [ADR-0046](../../adr/0046-sources-upstream-content-drift.md) (`Accepted`); Spezifikation [`DC-FA-SRC-001.a`](../../../../spec/spezifikation.md#dc-fa-src-001a--upstream-content-drift-externer-quellen-sources) + §2-Schema (`sources[]`).
- [x] **Code:** `rules/sources.go` (`CheckSources`, Netz-Post-Pass wie `external`) — Marker-Parser + Config-`sources[]`; `model.validModules()` + `SourcesConfig`; `configyaml` (raw/`applySources`; **kein** `sources.scope` — bare Liste, dok. Scope-Ausnahme); CLI `httpChecker()` deckt `sources`; `run.go` sammelt Pins + ruft nach dem Scan; `archive/zip`-Manifest-Hash; voller Ist-Hash im Befund; fail-closed. Grund-Codes `source-drift`/`source-unreachable` (§4) + `AllReasons()`/`reasonTexts()` (Lockstep). Netzlos-Test `forbiddenInNetless()` += `sources` (+ Regressionsfall).
- [x] **Tests:** die sechs Akzeptanzkriterien (Happy · Archiv-Determinismus · Modul-aus/netzlos · Negative-Drift · unreachable ≠ Drift · fail-closed) als Go-Tests gegen einen `httptest`-Server (kein echtes Netz im Unit-Test); Guards mutations-verifiziert (u. a. Manifest-Sortierung entfernt ⇒ Reorder-Test kippt).
- [x] **Config-Surface:** `--print-config`-Template (`sources`-Block + „einzige Netz-Tür"→zwei), `--print-mk` (automatisch über `ValidModules()`); Handbuch §6-Zeile + Config-Beispiel + §11 (Handbuch-Version 1.40) + `operations.md` + README EN/DE (die „einzige Netzwerk-Tür"-Stellen revidieren); `CHANGELOG.md`.
- [x] **Belege:** `make ci`/`make fullbuild` **grün**; Realdatenbeleg (`--enable sources` gegen eine echte gepinnte Quelle: unverändert grün, ein-Byte-Drift → `source-drift`); zwei unabhängige Reviews (R1 doc + R2 code), alle Befunde eingearbeitet; release-prep; Closure Move+Roadmap-Flip; [ADR-0046](../../adr/0046-sources-upstream-content-drift.md) → `Accepted`; Release **v0.51.0** + Digest-Backfill.

## 4. Risiken / offene Punkte

- **Die Netz-Sparsamkeit ist ein QA-Kernvertrag** — die zweite Netz-Tür wird bewusst amendiert
  (opt-in, nie Default); der Netzlos-Test (`forbiddenInNetless`) und die
  Messmethode müssen `sources` als Netz-Ausnahme kohärent führen.
- **Netz-Fläche im Go-Kern:** Download + `archive/zip`-Entpacken (Standardbibliothek,
  keine neue Dependency); Unit-Tests laufen gegen `httptest`, nie gegen echte URLs
  (Determinismus [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- **Doku-Currency-Falle:** viele „einzige Netzwerk-Tür"-Behauptungen (Lastenheft,
  `operations.md`, `config_template.go`, README EN/DE, Benutzerhandbuch) — der
  `versions`-Gate fängt sie **nicht**, manuell nachziehen.
- **Archiv-Hash-Vertrag:** Content-Manifest (pfad-sortiert, Pfad-normalisiert),
  nicht Zip-Roh-Bytes — der Reorder-Determinismus-Test ist der Beleg;
  Verzeichnis-Einträge des Zip ignorieren.
- **R1-doc-Review** ([Report](../../../reviews/2026-07-19-slice-080-sources-doc-first-r1.md), gegen `b36dc58`): **BLOCK** auf F-1/F-2 (Manifest-Kanonisierung widersprüchlich/unterbestimmt) → in [`DC-FA-SRC-001.a`](../../../../spec/spezifikation.md#dc-fa-src-001a--upstream-content-drift-externer-quellen-sources) byte-genau nachgezogen (pfad-sortiert, Pfad-Normalisierung, keine `sha256sum`-Byte-Parität); MEDIUMs F-3…F-6 (Case-Norm, Nicht-Zip → `source-unreachable`, Größen-/Zip-Bomben-Limit, Boundary-AKs) + LOW/INFO F-7…F-10 (Config-Pin-`line`, Keyword `archive`→`zip`, inerter Marker dokumentiert, Redirect-Politik) eingearbeitet. **Gegenprobe: ACCEPT-WITH-NITS** (BLOCK aufgehoben, F-1/F-2 verifiziert geschlossen); F-11 (Keyword-`archive`-Rest in Spec/ADR-Beispiel), F-12 (Multi-Link-Bindung), F-13 (Manifest-Tie-Break bei Duplikat-Pfad) nachgezogen.

## 5. Trigger

WIP-Slot frei nach welle-62-Abschluss (slice-079 `done`, v0.50.0). Die §4-Vorfragen
sind vom Nutzer entschieden: Pin-Deklaration **beides** (Marker + Config), Quelltypen
**Einzeldatei + Archiv**. Anlass: Nutzer-Frage „Drift gegen Upstream in d-check
einbauen" (2026-07-19).

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code + Spec; „Doc führt, Code folgt"). Keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** 19. Regelmodul `sources` (opt-in, **Netz**): pinnt eine externe
`http(s)`-Quelle auf einen `sha256` (Marker `<!-- source-pin: [zip] sha256:… -->`
am Link **oder** Config-Block `sources:`), holt sie als Netz-Post-Pass (geteilter
HTTP-Client wie `external`, Redirects ≤ 5, Body ≤ 64 MiB), hasht und vergleicht
case-insensitiv → `source-drift` (Meldung mit **vollem** Ist-Hash) bzw.
`source-unreachable`. Archiv (`unpack: zip`): byte-genaues, pfad-sortiertes
Content-Manifest (Basisname-frei, Verzeichnis-Einträge raus, Limits
256 MiB / 10 000). **Zweite Netz-Tür** — amendiert
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(`forbiddenInNetless` führt `sources`). Doc-first:
[`DC-FA-SRC-001`](../../../../spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz)
(Lastenheft 0.49.0) + [ADR-0046](../../adr/0046-sources-upstream-content-drift.md)
(`Accepted`) + Spezifikation `.a`/§2/§4 gingen dem Code voraus.

**Belege.**
- `make ci` **grün** (doc-check 261/0, lint, test, arch-check via a-check, Coverage
  **94,2 %**, semgrep 0, gate-consistency, planning-check; image-test nativ ==
  Container); `make completeness-check` **grün** (46 Anforderungen / **0 Waisen**).
- **AK-Tests** gegen `httptest` (Happy · Archiv-Reorder-Invarianz **mutations-echt
  gegen ein unabhängiges Manifest-Orakel** · Modul-aus/netzlos · Drift mit vollem
  Hash · unreachable ≠ Drift · fail-closed · Marker-Bindung/inert · Case ·
  verschachteltes Zip · Golden-Hash unabhängig via `sha256sum`).
- **Realdatenbeleg** ([Report](../../../reviews/2026-07-19-slice-080-sources-realdatenbeleg.md))
  gegen das echte `lab-regelwerk.zip` des Kurses (Archiv-Pfad, mit Netz): Dummy-Pin
  → `source-drift` mit vollem Ist-Hash, korrekter Pin → grün, GROSS-Pin → grün, ein
  Byte anders → `source-drift`.
- **Zwei unabhängige Reviews**
  ([R1-doc](../../../reviews/2026-07-19-slice-080-sources-doc-first-r1.md) /
  [R2-code](../../../reviews/2026-07-19-slice-080-sources-code-r2.md)): R1-doc
  **BLOCK** auf den Manifest-Kern (widersprüchliche Sortierung, undefinierte
  Pfad-Form) → byte-genau nachgezogen → Gegenprobe **ACCEPT-WITH-NITS**; R2-code
  **ACCEPT-WITH-NITS** (Marker-Hash-64-Validierung, Limit-Tests, Golden-Anker) —
  **alle Befunde eingearbeitet**.
- Release **v0.51.0** auf GHCR (Pipeline-Run 29679753205 grün), Digest-Pin
  `ghcr.io/pt9912/d-check@sha256:9197fcf0b6dd029637ba80088fff1f7a858287a0a0af13517f360d1437fc1d98`.

**Steering-Loop.** Klassifikation gemäß
[`grundlagen-klassifikation.md` §Steering Loop](../../../../.harness/baseline/v1.4.0/regelwerk/grundlagen-klassifikation.md#steering-loop):
Feedforward (doc-first-Vertrag + ADR) → Bau → Feedback (Doppel-Review +
Realdatenbeleg + Gates) → Release.

**Lerneintrag.**
1. **Der Doppel-Review trug messbar:** R1-doc fing den Manifest-Widerspruch
   (zeilen- vs. pfad-sortiert; „byte-gleich zu unsortiertem `sha256sum`") **vor**
   dem Code, R2-code die Marker-Hash-Divergenz (die Config erzwang 64 Hex, der
   Marker nicht → Falsch-`source-drift`) **nach** dem Code — beides wären sonst
   Bugs geworden. Design-Review-first *und* Code-Review-after zahlen sich getrennt
   aus.
2. **Ein Marker-Keyword-Rename ist ein Multi-Datei-Abgleich** (Lastenheft +
   normative Spec + ADR-Entscheidung + ADR-Beispiel + Slice): der halbe Rename
   (`archive`→`zip`, F-11) blieb in Spec-Schritt-1 + ADR-Beispiel stehen und wäre
   ein Falsch-Befund gewesen — die **Gegenprobe** desselben Reviewers fing genau das.
3. **Ein zweites Netz-Modul amendiert den Netz-Sparsamkeits-Kernvertrag:** der getippte
   Netzlos-Test (`forbiddenInNetless`) führt die Netz-Ausnahmen explizit; die
   §4-Grund-Codes kommen mit dem Code (AllReasons-↔-§4-Lockstep), nicht im doc-first.
