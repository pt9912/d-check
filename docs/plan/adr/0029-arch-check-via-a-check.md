# ADR-0029 — arch-check via Schwester-Tool a-check: das Gate konsumiert das digest-gepinnte a-check-Image statt `tools/arch-check.sh`

**Status:** Proposed
**Datum:** 2026-07-03
**Autor:** pt9912
**Bezug:** **Supersedes die Fitness-Function-Mechanik** von
[ADR-0005](0005-modul-layout-hexagon-ordner.md) (`tools/arch-check.sh`,
Dockerfile-Stage) und [ADR-0012](0012-kern-paketschnitt-model-rules-app.md)
(R6 im selben Skript) — die **Import-Regeln R1–R6 selbst bleiben unverändert
gültig** (Teil-Supersede, wie
[ADR-0026](0026-completeness-in-product-gate.md)/[ADR-0027](0027-commits-traceability-modul.md));
Mechanik-Familie [ADR-0024](0024-vcs-immutable-gate.md)/[ADR-0026](0026-completeness-in-product-gate.md)/[ADR-0027](0027-commits-traceability-modul.md)/[ADR-0028](0028-planning-lifecycle-modul.md)
(dieselbe „Skript → verteiltes, gepinntes Tool"-Linie — hier erstmals durch das
**Schwester-Tool** statt durch ein d-check-Modul); Gate-Image-Vorbild
[ADR-0010](0010-semgrep-hermetisches-gate.md) (externes, gepinntes Analyse-Image,
netzloser Lauf); Pin-Politik [ADR-0011](0011-digest-pins-build-gate-images.md)
(`@sha256:`-Digest für extern bezogene Build-/Gate-Images); Verteilungs-Kern
[`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)
(„verteilen statt kopieren"); die immutablen Skript-Referenzen werden ein
[`codepaths.ignore-refs`](../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)-Tombstone
([ADR-0025](0025-codepaths-ignore-refs.md)); keine Gate-Lockerung im Sinne von
[`AGENTS.md` §3.6](../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)
— genau deshalb diese ADR.
**Schärft:** keine neue Spec-Stelle — Prozess-/Mechanik-ADR (wie
[ADR-0026](0026-completeness-in-product-gate.md)); verbindlich für die Verdrahtung
des `make arch-check`-Gates. Die
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Bindung
des Gates bleibt.

## Kontext

`make arch-check` erzwingt die Import-Regeln des Hexagon-Schnitts: R1–R5 aus
[ADR-0005](0005-modul-layout-hexagon-ordner.md) (Kern-Reinheit, Tech-Kapselung
`net/http`/`yaml`/`os`, keine lateralen Adapter-Kanten) plus R6 aus
[ADR-0012](0012-kern-paketschnitt-model-rules-app.md) (Kern-Paket-Richtung
`model` ← `rules` ← `app`). Mechanik heute: `tools/arch-check.sh` übersetzt die
Regeln in `go list`-Prüfungen und läuft als Dockerfile-Stage — ein
handgeschriebenes, Go-spezifisches Shell-Skript.

Der Auftraggeber-`tools/*.sh`-Audit (2026-06-29) hat die vier d-check-fähigen
Gate-Skripte als d-check-Feature mechanisiert
([ADR-0024](0024-vcs-immutable-gate.md)/[ADR-0026](0026-completeness-in-product-gate.md)/[ADR-0027](0027-commits-traceability-modul.md)/[ADR-0028](0028-planning-lifecycle-modul.md));
`arch-check` wurde dabei bewusst **nicht** d-check zugeschlagen (Go-Import-Analyse
ist außerhalb des Doku-Checker-Produkt-Scopes), sondern dem **Schwester-Projekt
[a-check](https://github.com/pt9912/a-check)** zugewiesen (Roadmap §Nächste
Wellen). Dieses Ziel ist inzwischen einlösbar: a-check ist released (GHCR-Image,
distroless, aktuell v0.6.0), prüft Hexagon-Architekturen sprach-übergreifend
(sieben Sprach-Backends, darunter Go), konfiguriert über eine deklarative
`.a-check.yml` (`layers`-Globs, `edges`-Allowlist, `tech`-Kapselung,
`composition_root`) und verteilt sich wie d-check als digest-pinnbares Image plus
generiertes Makefile-Fragment (`a-check.mk` aus `a-check --print-mk`; der Lauf ist
`--network none` + read-only-Mount). Der strukturelle Anlass ist derselbe wie bei
[`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding):
vier Schwester-Repos pflegen heute divergente `arch-check.sh`-Varianten
(Copy-Drift-Klasse); a-check konsolidiert genau diese Familie. Umgekehrt
konsumiert a-check bereits d-check über ein eingebundenes `d-check.mk` — die
Schwester-Beziehung wird mit dieser Entscheidung symmetrisch.

## Entscheidung

`make arch-check` konsumiert das **digest-gepinnte a-check-Image**
(`@sha256:`-Pin gemäß [ADR-0011](0011-digest-pins-build-gate-images.md)) über ein
include-bares **`a-check.mk`** (aus `a-check --print-mk`, an die Repo-Politik
angepasst — analog dem `d-check.mk`, das a-check konsumiert) plus eine
repo-eigene **`.a-check.yml`**, die die Import-Regeln R1–R6 deklarativ ausdrückt
(`layers` für Kern/Ports/Adapter **und** die drei Kern-Pakete, `edges` für die
erlaubten Richtungen, `tech` für die `net/http`-/`yaml`-/`os`-Kapselung,
`composition_root` für CLI/`cmd`).

- **Policy unverändert, nur die Mechanik wechselt.** R1–R5
  ([ADR-0005](0005-modul-layout-hexagon-ordner.md)) und R6
  ([ADR-0012](0012-kern-paketschnitt-model-rules-app.md)) bleiben die
  verbindliche Regel-Menge; der Bindepunkt bleibt Produkt-Gate in
  `make gates`/`ci`.
- **Paritäts-Beleg vor Umstellung (Pflicht).** Je Regel R1–R6 eine adversariale
  Mutations-Probe (verbotener Import injiziert ⇒ `make arch-check` rot) — wie
  die Paritäts-Belege der Skript-Ablösungs-Präzedenzfälle. Verbleibende
  Präzisions-Deltas (text-heuristische Extraktion statt `go list`) werden
  ehrlich dokumentiert; ist ein Delta nicht per Config schließbar, ist das ein
  Change Request an das a-check-Lastenheft im Schwester-Repo — **keine** stille
  Lockerung hier.
- **Skript und Stage entfernt.** `tools/arch-check.sh` wird atomar per `git rm`
  entfernt, die Dockerfile-Stage `arch-check` (samt `--no-cache-filter`) entfällt;
  die immutablen Inline-Referenzen
  ([ADR-0005](0005-modul-layout-hexagon-ordner.md) §Fitness Function,
  [ADR-0012](0012-kern-paketschnitt-model-rules-app.md)) werden der fünfte
  [`codepaths.ignore-refs`](../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)-Tombstone
  ([ADR-0025](0025-codepaths-ignore-refs.md)).
- **Pin-Pflege wie semgrep.** Der a-check-Digest-Pin erscheint in
  `make versions`; eine Pin-Hebung ist ein bewusster Commit
  ([ADR-0011](0011-digest-pins-build-gate-images.md)-Politik).
- **Kein d-check-Produkt-Code.** Das d-check-Image bleibt byte-identisch —
  kein Release (wie [ADR-0026](0026-completeness-in-product-gate.md)).

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **a-check-Konsum via `a-check.mk` + `.a-check.yml` (gewählt)** | löst die Copy-Drift-Klasse der vier Schwester-Skripte ökosystemweit; Regeln deklarativ (Config) statt Bash; digest-gepinnt reproduzierbar, netzlos/read-only ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-konform); symmetrische Schwester-Beziehung | externe Abhängigkeit vom a-check-Release-Stand; text-heuristische Extraktion statt `go list`-Compiler-Wahrheit (Paritäts-Beleg nötig); Erst-Konfiguration der Schichten |
| Skript behalten | `go list`-Präzision; keine neue Abhängigkeit | Copy-Drift-Klasse bleibt in vier Repos; der Roadmap-Zeiger „→ a-check" bliebe uneingelöst |
| Regeln in d-check mechanisieren | ein Tool weniger im Stack | Go-Import-Analyse liegt außerhalb des d-check-Produkt-Scopes (Doku-Checker); im Audit explizit als nicht-d-check-fähig klassifiziert |
| Fertiges Go-Tool (go-arch-lint/depguard) | etabliert, präzise | Go-only — konsolidiert die C++/Rust/Kotlin-Schwestern nicht; neue Toolchain im Build statt gepinntes Schwester-Image |

## Konsequenzen

- `make arch-check` zieht ein externes Gate-Image (wie `make semgrep`,
  [ADR-0010](0010-semgrep-hermetisches-gate.md)): der Pull am Digest-Pin ist
  Setup (Netz, wie Image-Pull), der Prüf-Lauf selbst ist netzlos/read-only.
- **Präzisions-Trade:** `go list` (Compiler-Sicht, Build-Tags, Transitive)
  weicht der text-heuristischen Import-Extraktion (dokumentierte a-check-Grenze,
  dort `AC-QA-02`). Die Mutations-Proben je R1–R6 machen die reale Abdeckung
  sichtbar; jedes Rest-Delta wird beim Umsetzungs-Slice benannt statt
  verschwiegen. Sonderfälle, die die Proben klären müssen: die
  `net/url`-Erlaubnis in R1 (reiner Parser ohne I/O) und die R4-Dreifach-Zone
  (`os` in fs-Adapter, CLI und `cmd` — CLI/`cmd` als `composition_root`).
- Die Dockerfile verliert die `arch-check`-Stage; `make arch-check` wird vom
  `docker build --target` zum `docker run` des gepinnten Images.
- `gate-consistency` bleibt bindend: die Sensors-Tabelle
  ([`harness/README.md`](../../../harness/README.md)) und
  [`AGENTS.md`](../../../AGENTS.md) §4 ziehen auf die neue Mechanik um (Target-Name
  `arch-check` bleibt).
- Kein Selbstbezug/Henne-Ei (anders als die Dogfood-Module): a-check ist ein
  externes, bereits released Werkzeug.

## Fitness Function

- `make arch-check` läuft **rot** bei jeweils einer injizierten Verletzung der
  Regeln R1–R6 (Mutations-Proben, Beleg im Umsetzungs-Slice), **grün** auf dem
  sauberen Baum.
- Der Prüf-Lauf ist netzlos (`--network none`) + read-only-Mount
  ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Bindung
  gehalten); gleicher Baum ⇒ gleicher Befund
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)-analog,
  a-check-seitig als Determinismus-Zusage verankert).
- `make versions` weist den a-check-Digest-Pin aus.
- `make gate-consistency` grün (Doku ↔ Makefile in beiden Richtungen).

## Re-Evaluierungs-Trigger

- Ein dokumentiertes Paritäts-Rest-Delta wird durch ein a-check-Release
  schließbar → Pin-Hebung (bewusster Commit).
- Das Modul-Layout ändert sich (neue Pakete/Adapter) → `.a-check.yml`-Anpassung
  — jetzt Config, nicht Skript-Code; die Regel-Quelle bleiben
  [ADR-0005](0005-modul-layout-hexagon-ordner.md)/[ADR-0012](0012-kern-paketschnitt-model-rules-app.md)
  (bzw. deren Nachfolger).
- a-check-Regression oder ausbleibende Pflege → Rückkehr zur Skript-Mechanik
  aus der git-Historie ist möglich (Re-Eval, neue ADR).

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-07-03 | Entwurf (slice-058, welle-47): `make arch-check` konsumiert das digest-gepinnte a-check-Image via `a-check.mk` + `.a-check.yml` statt `tools/arch-check.sh` (Dockerfile-Stage); Regeln R1–R6 bleiben Policy, Paritäts-Beleg per Mutations-Proben Pflicht, Skript + Stage entfernt + fünfter Tombstone. Löst den Roadmap-Zeiger „arch-check → Schwester-Projekt a-check" ein. Kein d-check-Produkt-Code, kein Release. Status Proposed. |
