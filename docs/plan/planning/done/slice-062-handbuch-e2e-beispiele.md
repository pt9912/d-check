# Slice slice-062: Handbuch-E2E-Beispiel-Verankerung (Kommando/Ausgabe ↔ echtes Binary)

**Status:** done (welle-51-handbuch-e2e-verankerung, Closure 2026-07-05).

**Welle:** welle-51-handbuch-e2e-verankerung (Folge von welle-50;
ausgegliedert aus slice-061 per Nutzer-Entscheid 2026-07-04 „A jetzt, B als
slice-062", in einer Sitzung aus dem Backlog gezogen und abgeschlossen).

**Bezug:** Verifikations-Mechanik gegen bestehende Verträge
([`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)
Exit-Codes;
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
Image-Verhalten). **Kein Change Request** (kein neuer Vertrag), **kein ADR**
(E2E-Test-Erweiterung im bestehenden Schnitt). Schwester-Slice zu
[`slice-061`](welle-50/slice-061-doc-config-beispiel-verifikation.md)
(dort Dimension A: Config-Fragmente ↔ Parse; hier Dimension B: Kommando-/
Ausgabe-Beispiele ↔ Verhalten).

**Autor:** pt9912. **Datum:** 2026-07-04.

---

## 1. Ziel

Das Benutzerhandbuch führt **Kommando-Aufrufe mit dokumentierter Ausgabe/
Exit-Code** (` ```bash ` + ` ```text `) — sauberes Repo ⇒ Exit 0 +
„0 Befund(e)"; kaputter Link ⇒ Befund-Zeile `Datei:Zeile  Ziel  Grund` +
Exit 1; `--doctor` ⇒ gruppierte Diagnose. Diese Verhaltensbehauptungen werden
nie gegen das echte Binary geprüft; driftet das CLI-Verhalten, bleibt das
Handbuch still falsch. **Neu:** repräsentative Kommando-Beispiele werden als
E2E-Fälle gegen das reale Binary/Image über Fixtures verankert (Anregung des
Auftraggebers: „von E2E-Tests könnte ein Handbuch profitieren, da ja im
Handbuch Beispiele aufgeführt sind").

## 2. Entscheidungen

- **Vorhandene E2E-Infrastruktur nutzen** statt neuen Runner bauen:
  `tools/image-test.sh` ([`make image-test`](../../../../harness/README.md#sensors-feedback-gates))
  fährt schon nativ == Container mit Fixtures + Exit-Code-/Ausgabe-Prüfung;
  `cli_acceptance_test.go` fährt CLI-E2E gegen `MemFS`/Temp-Repos.
- **Nur Beispiele mit prüfbarer Verhaltensbehauptung** verankern: je ein
  Fixture, das die Prämisse herstellt (sauberes Repo / kaputter Link /
  `--doctor`-Fall), der dokumentierte Aufruf (die **Flags**, nicht der
  wörtliche `docker run …@sha256`-Pull), Assertion auf Exit-Code + Ausgabe-
  **Form** (Befund-Zeilen-Schema, „N Befund(e)"-Zeile, Diagnose-Kopf).
- **Auf Form prüfen, nicht auf wörtliche Zeilen:** Regex/Schema statt
  Datei-Zahlen/Versions-Pins — sonst bricht der Test bei jeder harmlosen
  Doku-/Release-Änderung (Wartungsfalle statt Wert).
- **Nicht-replaybare Beispiele markieren** (externer Zustand, konkrete
  Digests, Netz) — gleicher Opt-out-Marker-Ansatz wie slice-061, mit Grund;
  kein stiller Ausschluss (welche Beispiele E2E-verankert sind, welche nicht
  und warum).
- **Determinismus/netzlos:** lokal gebautes Image (wie `image-test`), kein
  Pull.

## 3. Definition of Done

- [x] repräsentative Handbuch-Kommando-Beispiele mit Verhaltensbehauptung als
  E2E-Fälle (Fixture + Flags + Exit-Code-/Ausgabe-Form-Assertion) — sieben
  verankert in `handbook_examples_test.go` (eigene Datei im Paket `cli_test`,
  nutzt `run`/`write`/`traceRepo` aus `cli_acceptance_test.go`); die
  Container-/`--network none`-Claims bleiben in `tools/image-test.sh` verankert.
- [x] Beispiel-Auswahl dokumentiert (E2E-verankert vs. begründet ausgenommen —
  im Datei-Kopf-Kommentar des Harness).
- [x] Fail-closed-Guards mutations-verifiziert; adversariale Probe
  (dokumentierte Ausgabe künstlich verletzt ⇒ rot) — s. §7.
- [x] `make gates`/`make ci` grün; unabhängiges Review R1 (ACCEPT);
  Closure-Move + Body. **Kein Produkt-Code, kein Release** (Test-Infra).

## 4. Risiken / offene Punkte

- **Ausgabe-Matching-Stabilität** (s. §2): Form statt Wortlaut — der
  Kern-Entscheid, damit der Test kein Wartungsklotz wird.
- **Beispiel-Auswahl:** nicht jedes der ~15 Kommando-Beispiele lohnt einen
  E2E-Fall; die mit einer eindeutigen, stabilen Verhaltensbehauptung zuerst.
- **Abgrenzung zu bestehenden Tests:** `cli_acceptance_test.go`/`image-test`
  decken viele Verhalten schon ab — Dimension B ergänzt gezielt die
  **im Handbuch gezeigten** Fälle, ohne Bestehendes zu duplizieren.

## 5. Trigger

Ausgegliedert aus [`slice-061`](welle-50/slice-061-doc-config-beispiel-verifikation.md)
per Nutzer-Entscheid 2026-07-04 (Stufung „A jetzt / B als slice-062"); die
E2E-Anregung selbst kam vom Auftraggeber in derselben Sitzung. Aufnahme in eine
Welle, wenn slice-061 (Dimension A) abgeschlossen ist.

## 6. Sub-Area-Modus-Begründung

GF (E2E-Test-Erweiterung im bestehenden Schnitt). Kein neuer Adapter, keine
BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung.** Dimension B als Go-E2E-Harness
`internal/adapter/driving/cli/handbook_examples_test.go` (Paket `cli_test`,
eigene Datei statt Anbau an das ~1600-Zeilen-`cli_acceptance_test.go`; nutzt
dessen `run`/`write`/`traceRepo`). Er liest das Benutzerhandbuch über einen
repo-relativen Pfad (Präzedenz slice-060/061) und koppelt die dokumentierten
Kommando-/Ausgabe-Beispiele an das **echte Binary** (`cli.Run`):

- **`TestHandbook_AnchoredExamplesReplay`** — sieben verankerte Beispiele
  (§3/§4.1 sauberes Repo + kaputter Link, §4.9 `--doctor` + `--doctor --json`,
  §4.11 `--json` + `--yaml`, §4.12 `--trace`): je Beispiel eine
  Prämissen-Fixture, die **dokumentierten** Flags (per Assertion auch im
  `bash`-Block belegt), Prüfung auf Exit-Code + Ausgabe-**Form**. Form-Anker
  müssen in **beidem** stehen — dem dokumentierten Ausgabeblock UND der echten
  Ausgabe (Drift auf jeder Seite ⇒ rot); für `json`/`yaml` zusätzlich
  strukturelle Schlüssel-Kopplung (dokumentiert ⊆ real). Geprüft wird die
  **Form** (Befund-Zeilen-Schema, „N Befund(e)"-Zeile, Diagnose-Kopf, RTM-Kopf,
  Schlüssel-Mengen), nicht Datei-Zahlen/Pfade — keine Wartungsfalle
  (Slice-§4-Kern-Risiko).
- **`TestHandbook_OutputBlocksClassified`** — fail-closed: jeder `text`/`json`-
  **und** jeder per slice-061-`not-config`-Marker als CLI-Ausgabe deklarierte
  `yaml`-Block ist entweder an ein Beispiel verankert oder trägt den
  `not-replayable`-Marker. Ein neuer, unklassifizierter Ausgabeblock ⇒ rot
  (kein stiller Ausschluss).
- **Marker & Ausgenommen mit Grund** (im Datei-Kopf dokumentiert): genau ein
  Nicht-replaybarer Block — das abgekürzte `--print-mk`-Fragment (Elision) —
  trägt `<!-- d-check-test:not-replayable: … -->` (doc-first). §4.5–4.8
  (prosaische Verhaltensaussagen ohne Ausgabeblock) sind von den
  Modul-Akzeptanztests gedeckt; die Container-Mount-/`--network none`-/
  nativ==Container-Claims von `tools/image-test.sh`; die Digest-/`docker pull`-
  Beispiele sind nicht replaybar (Netz/Registry).

**Belege.**

- `make gates`/`make ci` **grün** (doc-check 184/0, lint, test, arch-check,
  Coverage 93,80 %, semgrep 0, gate-consistency, planning-check; image-test).
- **Mutations-Belege:** (1) adversariale Probe — ein Form-Anker im
  `--doctor`-Doku-Block verfälscht („ohne Markdown-Link" → „ohne MD-Link")
  ⇒ rot an der Block-Zeile; (2) Silent-Exclusion-Guard — den
  `not-replayable`-Marker entfernt ⇒ der `--print-mk`-Block unklassifiziert
  ⇒ rot; beide reversibel per Edit wiederhergestellt (nie `git checkout`).
  Extraktor-Fail-closed (unbalancierter Fence) und die Schlüssel-Rename-
  Erkennung sind per Unit-Test (`TestExtractFencedBlocks`,
  `TestHandbookVerifiers`) permanent verriegelt.
- **Unabhängiges Review R1**
  ([Report](../../../reviews/2026-07-05-slice-062-handbuch-e2e-harness-r1.md)):
  **ACCEPT**, 0 HIGH/0 MEDIUM/1 LOW/1 INFO — beide eingearbeitet. F-1 (LOW):
  die Fail-closed-Sweep überging `yaml` ganz, sodass eine künftige
  yaml-CLI-Ausgabe (z. B. das in §4.9 zugesagte `--doctor --yaml`) zwischen
  slice-061 (per `not-config` ausgenommen) und slice-062 durchgerutscht wäre —
  geschlossen über die **Brücke** (ein `not-config`-yaml IST CLI-Ausgabe ⇒
  gehört in die Sweep). F-2 (INFO): die Aussage „Drift in beide Richtungen"
  gilt streng nur für die formTokens; die strukturelle
  `dokumentiert ⊆ real`-Kopplung (bewusst, sonst False-Positive bei
  abgekürzten Doku-Blöcken) präzisiert.
- **Kein Change Request, kein ADR, kein Release** (Test-/Harness-Infra + ein
  HTML-Kommentar-Marker; Produkt-Image byte-identisch, `make image-test` grün).

**Lerneintrag.** (1) Der Schwester-Bezug zu slice-060/061 hält: alle drei
koppeln **Doku an Code** über eine Mengen-/Format-Verriegelung mit
Fail-closed-Guards — hier Handbuch-Ausgabe-Beispiel ↔ echtes Binary (statt
Config-Fragment ↔ Validator bzw. Grund-Code-Liste ↔ Spec-§4). (2) Zwei
Harnesse an derselben Doku brauchen eine **explizite Brücke**: der
slice-061-`not-config`-Marker wird zum Übergabepunkt, sonst entsteht ein
Klassen-Loch genau zwischen den Zuständigkeiten (R1-F-1). (3) Ausgabe-Matching
**auf Form, nicht Wortlaut** ist der Kern-Entscheid gegen die Wartungsfalle —
Zahlen/Pfade normalisiert, invariante Template-Wörter/Schlüssel gekoppelt.
(4) Ein fremd-generierter Review-Report ist selbst Live-Doku: die als
Inline-Code notierten Fence-Typen (drei Backticks plus Sprachname) brachen die
`spans`-Backtick-Parität pro Absatz (doc-check rot) — Reports dürfen
Fence-Typen nur backtick-neutral nennen. Steering-Loop: die
Handbuch-Verhaltensbehauptungen sind jetzt als Sensor verankert, nicht nur
einmal von Hand geprüft. Commit-Kette: doc-first `2b72d23` → feat `abc1045` →
R1 `eb1bd75` → closure-move `c0bfa7e` → closure-body.
