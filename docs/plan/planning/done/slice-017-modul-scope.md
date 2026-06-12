# Slice slice-017: Modul-lokaler Scan-Scope (`<modul>.scope`)

**Status:** done.

**Welle:** welle-07-modul-scope (per Roadmap-Fortschreibung; Start
bei Priorisierung durch den Auftraggeber — der Erst-Bedarfsträger
grid-gym wartet konkret).

**Bezug:** [`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope)
(Change Request 0.7.0, Erst-Bedarfsträger grid-gym),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(Vollvalidierung, Constraint-Spiegelung),
[`DC-FA-SCAN-001`](../../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln)
(Scan-Regeln, `SKIP_DIRS`, Pruning),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Determinismus über Modul-Scopes hinweg).

**Autor:** pt9912 (CR-Einreichung grid-gym). **Datum:** 2026-06-12.

---

## 1. Ziel

Regelmodule können ihren Scan-Scope lokal ersetzen
(`<modul>.scope.roots`/`.ignore`) — gemessen am Erst-Bedarf: grid-gym
aktiviert `ids` kuratiert auf `spec/` + `docs/user/` (312 echte
Befunde, davon 190 Traceability-Gewinne in der Architektur-Spec)
statt global (2776 Befunde, deren Masse Retro-Verlinkung
historischer done-Planning-Docs wäre — Umschreiben des
Audit-Trails ohne Qualitätsgewinn), während `links`/`anchors`
weiterhin die volle Repo-Wurzel prüfen.

## 2. Definition of Done

- [x] Spezifikation fortgeschrieben: §`DC-FA-CONF-002.a`
  (effektiver Scope pro Modul: eigener Discover-Lauf mit
  Modul-`roots`/`ignore` — der Modul-Scope kann Dateien umfassen,
  die der globale Scan nicht enthält; Konstraint-Spiegelung von
  `scan.*` inkl. Existenz, Repo-Escape-Verbot, Abstiegs-Pruning,
  `SKIP_DIRS`, Leere-Menge-Semantik), Schema-Tabelle
  §`.d-check.yml` um `<modul>.scope.roots`/`.ignore`; **keine**
  neuen Grund-Codes.
- [x] Implementierung: Datei-Ermittlung pro Modul-Scope
  (deterministische Reihenfolge je Scope; Datei-Inhalte werden
  weiterhin genau einmal gelesen — Scope-Discovery getrennt vom
  Inhalts-Cache); die vier Akzeptanzkriterien aus
  [`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope)
  als Tests.
- [x] **Abwärtskompatibilitäts-Beleg:** Configs ohne `scope`
  verhalten sich byte-identisch (Regressions-Test über die
  bestehende Test-Suite hinaus: Eigenlauf und mindestens ein
  migriertes Schwester-Repo vor/nach, identische Ausgabe-Hashes —
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)-Methodik).
- [x] **Konsumenten-Abnahme grid-gym** (dortige Overlay-Config):
  `ids.scope.roots: ["spec", "docs/user"]` liefert exakt die
  kuratierten Befunde, `links`/`anchors` unverändert über den
  globalen Scope; Ergebnis als Vergleichszahlen dokumentiert. Der
  anschließende gestaffelte Fix-Sweep (erst `docs/user/`, dann der
  Architektur-Traceability-Sweep) ist grid-gym-Arbeit, nicht
  Slice-Bestandteil.
- [x] Dogfooding: Selbstkonfiguration bleibt unverändert lauffähig
  (kein eigener `scope`-Bedarf); `gate-consistency`-gebundene
  `DC-QA-03`-Modulliste unberührt.
- [x] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md);
  Closure-Notiz mit Steering-Loop-Lerneintrag. Release-Hinweis:
  fließt in das nächste Minor-Release (mit welle-06 gebündelt als
  v0.3.0, sonst eigenes v0.x) — Digest in den Release-Notes für den
  grid-gym-Pin-Bump.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | §`DC-FA-CONF-002.a`, Schema-Tabelle |
| Config-Decoding + Validierung | update | `<modul>.scope` strikt validiert (Exit 2) |
| Datei-Ermittlung im Hexagon-Kern | update | Discover pro Modul-Scope, geteilte Helfer |
| Tests (4 AKs, Abwärtskompatibilität) | neu | Beleg-Pflicht |
| [`CHANGELOG.md`](../../../../CHANGELOG.md) | update | nutzersichtbarer Config-Schlüssel |

## 4. Trigger

Change Request 0.7.0 im Lastenheft (erfüllt 2026-06-12) **und**
Priorisierung durch den Auftraggeber. Unabhängig von
[slice-015](../done/slice-015-spans-modul.md)/[slice-016](../in-progress/slice-016-hostpaths-modul.md)
implementierbar (orthogonal: Scope vs. neue Module); bei Bündelung
in ein gemeinsames v0.3.0 profitieren künftige Module sofort vom
Scope-Schlüssel.

## 5. Closure-Trigger

DoD vollständig (insbesondere Abwärtskompatibilitäts-Beleg und
dokumentierte grid-gym-Abnahme) + Closure-Notiz.

## 6. Risiken und offene Punkte

- **Performance/Architektur:** Heute gibt es einen globalen
  Discover-Lauf; pro-Modul-Scopes brauchen mehrere. Die
  [`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance)-Schranke
  gilt weiter — Scope-Discovery ist billig (Verzeichnis-Walk),
  Datei-Inhalte dürfen nicht mehrfach gelesen werden (Inhalts-Cache
  über Scopes hinweg).
- **Determinismus-Falle:** Befund-Sortierung ist global (Pfad,
  Zeile) — Befunde verschiedener Module mit verschiedenen Scopes
  müssen stabil zusammengeführt werden (bestehende Sortier-Regel
  reicht, Test belegt es).
- Per-Pattern-Scope ist bewusst out-of-scope (erst bei zweitem
  Bedarfsträger) — im Lastenheft dokumentiert, nicht stillschweigend
  vorgebaut.

## 7. Closure-Notiz (nach `done/`)

**Abgeschlossen:** 2026-06-12, am Tag des Change Requests.

- **Umsetzung** (933cf14): Adapter-seitig `scope` an allen sechs
  Modul-Sektionen (`links`/`anchors` neu als reine scope-Träger),
  `roots`-Pflicht strikt; Kern-seitig Discover pro Modul-Scope mit
  Vereinigungsmenge und Einmal-Lese-Garantie; `modules: []` behält
  die globale Datei-Zählung. gocognit-Grenzen per Refactoring
  eingehalten (`runState`/`checkFile`, `discoverInto`,
  `applyScopes`) — keine Suppressions. Coverage 95,9 %.
- **Abwärtskompatibilität belegt:** v0.2.1-Release-Image vs. neuer
  Stand auf zwei Korpora (Eigenlauf, m-trace) — Ausgabe
  byte-identisch (Hash-Vergleich).
- **Konsumenten-Abnahme grid-gym** (Overlay nach CR-Vorlage,
  Muster-Näherung an deren Kennungs-Schema, lokales Image):
  global 2378 `id-unlinked` (CR maß 2776 mit deren exaktem
  Muster-Satz auf älterem Stand — gleiche Größenordnung, gleiche
  Diagnose), kuratiert (`ids.scope.roots: [spec, docs/user]`)
  **311** Befunde (CR-Prognose: 312), **0 außerhalb** des Scopes;
  `links`/`anchors` unverändert global (205 Dateien, 0 Befunde).
  Der produktive Umstieg dort folgt nach dem Release per
  Digest-Bump (CR §Abnahme).

### Lerneintrag (Steering Loop)

Der erste extern eingereichte Change Request hat die
CR-Maschinerie validiert: Weil der Bedarfsträger Motivation als
**Messung** lieferte (2776 vs. 312) und AKs im Haus-Stil mitbrachte,
war der Weg CR → Lastenheft 0.7.0 → Spec → Implementierung →
Abnahme an einem Tag gangbar. Geschärfte Regel für künftige CRs:
die Messung des Bedarfsträgers ist das Abnahme-Soll — die
Konsumenten-Abnahme gehört als eigener DoD-Punkt in jeden
CR-Slice (hier erstmals so geschnitten, Wiederverwendung
empfohlen).

## 8. Sub-Area-Modus-Begründung

Dieses Repo: GF (Spec führt, Code folgt — Change Request 0.7.0 vor
Implementierung). grid-gym: Bedarfsträger und Abnahme-Instanz,
Fix-Sweeps dort folgen deren Prozess.
