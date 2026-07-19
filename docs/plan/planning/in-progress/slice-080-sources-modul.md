# Slice slice-080: Modul `sources` — Upstream-Content-Drift externer Quellen

**Status:** In Arbeit (doc-first abgeschlossen; Code offen)
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
  `pins`/dpin-Hash-Ergonomie-Fix bleibt **separat** ([`slice-072`](../open/slice-072-handbuch-aufgabenorientierung.md)) — dieser Slice schneidet nicht quer.
- Begründung/Alternativen in [ADR-0046](../../adr/0046-sources-upstream-content-drift.md).

## 3. Definition of Done

- [x] **Spec/Doc (doc-first):** [`DC-FA-SRC-001`](../../../../spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz) (Lastenheft 0.49.0) + Bereich `SRC` §3 + [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Erweiterung + §7 + `sources` in der Regelmodul-Auswahl + Glossar; [ADR-0046](../../adr/0046-sources-upstream-content-drift.md) (`Proposed`); Spezifikation [`DC-FA-SRC-001.a`](../../../../spec/spezifikation.md#dc-fa-src-001a--upstream-content-drift-externer-quellen-sources) + §2-Schema (`sources[]`).
- [ ] **Code:** `rules/sources.go` (`CheckSources`, Netz-Post-Pass wie `external`) — Marker-Parser + Config-`sources[]`; `model.validModules()` + `SourcesConfig`; `configyaml` (raw/`applySources`/`scopeOfSources`); CLI `httpChecker()` deckt `sources`; `run.go` sammelt Pins + ruft nach dem Scan; `archive/zip`-Manifest-Hash; voller Ist-Hash im Befund; fail-closed. Grund-Codes `source-drift`/`source-unreachable` (§4) + `AllReasons()`/`reasonTexts()` (Lockstep). Netzlos-Test `forbiddenInNetless()` += `sources` (+ Regressionsfall).
- [ ] **Tests:** die sechs Akzeptanzkriterien (Happy · Archiv-Determinismus · Modul-aus/netzlos · Negative-Drift · unreachable ≠ Drift · fail-closed) als Go-Tests gegen einen `httptest`-Server (kein echtes Netz im Unit-Test); Guards mutations-verifiziert (u. a. Manifest-Sortierung entfernt ⇒ Reorder-Test kippt).
- [ ] **Config-Surface:** `--print-config`-Template (`sources`-Block + „einzige Netz-Tür"→zwei), `--print-mk` (automatisch über `ValidModules()`); Handbuch §6-Zeile + Config-Beispiel + §11 (Handbuch-Version 1.40) + `operations.md` + README EN/DE (die „einzige Netzwerk-Tür"-Stellen revidieren); `CHANGELOG.md`.
- [ ] **Belege:** `make ci`/`make fullbuild` **grün**; Realdatenbeleg (`--enable sources` gegen eine echte gepinnte Quelle: unverändert grün, ein-Byte-Drift → `source-drift`); zwei unabhängige Reviews (R1 doc + R2 code), alle Befunde eingearbeitet; release-prep; Closure Move+Roadmap-Flip; [ADR-0046](../../adr/0046-sources-upstream-content-drift.md) → `Accepted`; Release **v0.51.0** + Digest-Backfill.

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

*(bei Closure gefüllt: Umsetzung, Belege, Lerneintrag.)*
