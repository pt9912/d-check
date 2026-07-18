# Slice slice-071: Trace-Kreuzverweis-Konsistenz-Gate (Vorwärts-RTM ↔ Rück-Kanten)

**Status:** in-progress — **am 2026-07-18 wieder aufgenommen**
(`open`→`in-progress`, Blocker aufgelöst). Code, Spezifikation,
[ADR-0038](../../adr/0038-trace-cross-consistency.md) und vier Review-Runden plus
Closure-Review liegen vor; v0.44.0/v0.45.0 sind getaggt und veröffentlicht. Der
letzte offene DoD-Punkt — der Realdatenbeleg gegen grid-gym — ist **erbracht**
(d-check v0.48.1): der `17. Testarchitektur`-Abschnitt, der bis v0.48.0 an
`spec/architecture.md:913` an d-checks eigener Ignore-Direktive in einer
Tabellenzeile abbrach (Exit 2), läuft mit der Direktiven-Toleranz aus slice-074
([ADR-0040](../../adr/0040-kommentar-suffix-in-tabellenzeilen.md)) durch.

**Welle:** welle-60-trace-cross-consistency — **inhaltlich weiter Teil der Welle**,
aber nicht in Arbeit. Modul 5: `in-progress→open` ist der Übergang für „Blocker,
Priorität offen"; ein blockierter Slice gehört **sichtbar zurück, nicht still
weitergeschoben**. Wiederaufnahme, sobald slice-074 eine tragende Regel hat.

**Bezug:** neuer Lastenheft-Change-Request
[`DC-FA-XREF-001`](../../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
(Lastenheft 0.44.2 — 0.44.0 initial, 0.44.1/0.44.2 als
Review-getriebene Vakuitäts-Schärfungen), Mit-Änderung
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
und Gate über
[`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code);
begründende Entscheidung [ADR-0038](../../adr/0038-trace-cross-consistency.md)
(Proposed). Produkt-/Config-Delta; Spezifikation und ADR sind doc-first vor der
Implementierung fortgeschrieben.

**Autor:** pt9912. **Datum:** 2026-07-16.

---

## 1. Ziel

Zwei unabhängig gepflegte Sichten derselben Anforderung→Design-Relation driften
heute unbemerkt: eine Vorwärts-RTM-Tabelle (Anforderung → Design-Artefaktmenge)
und Rückwärts-`Bezug`-Kanten (jedes Artefakt nennt aufwärts seine Anforderungen).
Kein Modul vergleicht die konkreten Mengen. Der Slice liefert einen opt-in
Mengenabgleich als `trace`-Unterfähigkeit
([`DC-FA-XREF-001`](../../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)),
der je Anforderung die beiden Richtungsdifferenzen meldet — mit vollem
Reader-Reuse und byte-identischer RTM für Nicht-Konsumenten.

## 2. Vertrag aus dem Change Request

- **Zwei Sichten, ein Reader:** Vorwärts- und Rück-Sicht sind kuratierte
  Markdown-Tabellen, gelesen über den header-gebundenen Reader
  ([`DC-FA-REQ-001`](../../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen))
  + die range-aware Span-Semantik von
  [`DC-FA-COV-001`](../../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in).
  Nur invertieren + Set-Diff ist neu.
- **Rück-Kanten = Quelle der Wahrheit:** je Tabelle mit `Bezug`-Spalte
  (header-gebunden) ist die Artefakt-ID die **erste Spalte** (positionell, da deren
  Header je Tabelle variiert), per `design-pattern` extrahiert; die `Bezug`-Zelle
  liefert range-aware die Anforderungs-IDs (`req-pattern`).
- **Vorwärts:** header-gebundene Anforderungs-/Design-Spalte; Design-Artefakte per
  `design-pattern`, Anforderungen range-aware.
- **Mengen-Diff:** je Anforderung `R` werden `F(R) \ B(R)` und `B(R) \ F(R)` mit
  Richtungslabel gemeldet; Modus `equal` (beide gaten) oder `superset` (nur
  `B \ F`). `1:N` Normalfall.
- **Ableitungssprünge:** `exclude-req` (RE2) nimmt Mittelschicht-IDs aus dem
  Abgleich (benanntes Ventil, driftet selbst).
- **Gate-Auslöser:** advisory unter `--trace`; Exit-Änderung nur über das globale
  `--require-complete`
  ([`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)),
  kein block-lokaler Schalter.
- **Fail-closed / byte-identisch:** ungültiges Regex, fehlende Spalte, ID-Header
  nicht genau einmal, unbekannter `mode` ⇒ Exit 2; ebenso ein **vakuumer Abgleich**
  (beide Sichten kantenleer, oder Rück-Sicht kantenleer unter `superset`, oder
  `exclude-req` verschluckt alles) — nachgezogen als Lastenheft 0.44.1/0.44.2 nach
  den Review-Runden R1–R4; ohne `trace.cross-consistency`-Block ist die RTM
  byte-identisch
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)), read-only
  ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

## 3. Definition of Done

- [x] **Lastenheft-CR:**
  [`DC-FA-XREF-001`](../../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  mit Happy/Boundary/Negative, Out-of-Scope, Bereich `XREF`, Version 0.44.0 und
  Historie; additive `--trace`-Erweiterung
  ([`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
  ohne RTM-Änderung.
- [x] **Spezifikation:** Algorithmus [`DC-FA-XREF-001.a`](../../../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency) (Vorwärts-/Rück-Extraktion,
  Inversion, Mengen-Diff, `exclude-req`, Fehlerpräzedenz) + Config-Schema-Zeilen
  (`trace.cross-consistency.*`) + Historie.
- [x] **ADR + Index:** [ADR-0038](../../adr/0038-trace-cross-consistency.md)
  (Platzierung, Reader-Reuse, Generator-Sequenzierung), Status Proposed, im Index.
- [x] **Modell/Config:** `trace.cross-consistency` abbilden, vollständig
  validieren, in `--print-config` sichtbar machen.
- [x] **Extraktion:** Rück-Kanten (erste-Spalte-ID + header-gebundene `Bezug`) und
  Vorwärts-Sicht (header-gebunden) deterministisch/read-only; range-aware Reuse.
- [x] **Set-Diff:** Inversion + `F\B`/`B\F` je `R`, Modi `equal`/`superset`,
  `exclude-req`-Ventil; deterministisch sortierte Befunde.
- [x] **Gate-Bindung:** advisory unter `--trace`, Gatung über globales
  `--require-complete`.
- [x] **Tests:** konsistentes 1:N grün, beide Richtungsdifferenzen, range-aware,
  superset, Ableitungssprung-Ausschluss, fail-closed-Config, byte-identischer
  Default.
- [x] **Realdatenbeleg grid-gym:** erbracht mit d-check **v0.48.1** gegen die realen
  grid-gym-Quellen (§27.1-Vorwärts-Sicht gegen die `Bezug`-Rück-Kanten der
  Architektur). Der reale §27.1-↔-`Bezug`-Drift wird geflaggt —
  `GG-ARCH-005` (§27.1 nennt `COMP-CORE`/`COMP-DOMAIN`, Rück-Kanten
  `COMP-SCHED`/`P-005`/`P-009`, Schnittmenge null) und `GG-SIM-009`, beide Richtungen
  gelabelt; die per `exclude-req` ausgeschlossenen Familien nicht (Ventil greift,
  0 Befunde); das konsistente 1:N läuft grün. Die dokumentierte Messung
  (161 Differenzen) reproduziert. Der bis dahin blockierende `17. Testarchitektur`-
  Abschnitt läuft durch — bis v0.48.0 brach er an `architecture.md:913` mit „913 hat
  4 statt 3 Zellen" ab; die Direktiven-Toleranz aus slice-074 ist der entblockende
  Mechanismus (Vorher/Nachher belegt). Beleg:
  [`2026-07-18-slice-071-realdatenbeleg-grid-gym`](../../../reviews/2026-07-18-slice-071-realdatenbeleg-grid-gym.md).
- [x] **Nutzerdoku:** Handbuch 1.30 (§4.12 Aufgabe mit **am Image verifizierter**
  Ausgabe + E2E-Verankerung im Beispiel-Harness, §5 Config), CHANGELOG 0.44.0,
  `operations.md` (`--trace` + die zweite `--require-complete`-Gate-Ursache).
- [x] **Release:** v0.44.0 (Kern) + v0.45.0 (Scope-Trennung `forward.req-pattern`)
  getaggt, auf GHCR gepusht und als GitHub-Release veröffentlicht (2026-07-17); der
  Digest-Backfill erfolgte zum jeweiligen Release, der Handbuch-§4-Pin ist seither
  weitergezogen. Der Realdatenbeleg selbst ist Verifikation (kein Code) und braucht
  kein eigenes Release.
- [x] **Qualität:** vier unabhängige, kontext-getrennte Review-Runden (R1 REJECT →
  R4 ohne Code-Vorbehalt) plus unabhängiger Closure-Review
  ([`2026-07-17-slice-071-closure-independent.md`](../../../reviews/2026-07-17-slice-071-closure-independent.md),
  ACCEPT-WITH-NITS, Tag freigegeben; die Nits sind eingearbeitet); `make gates` und
  `make ci` grün.

## 4. Risiken / offene Designpunkte

- **Positioneller `first`-Modus:** die Backward-ID-Spalte hat heterogene Header
  (`Kennung`/`Port-ID`/`Tabu-ID`/`Komponente`), aber `Bezug` ist einheitlich → die
  Artefakt-ID kommt aus der ersten Spalte (via `design-pattern`), nicht per
  Header-Name. Ein neuer Extraktions-Sonderweg neben der Header-Bindung.
- **Namensraum-Kongruenz:** `forward.design-pattern` und die Backward-Erste-Spalte
  müssen denselben Namensraum liefern (`GG-AR-*`), sonst ist der Set-Diff
  bedeutungslos — als Vorbedingung zu spezifizieren.
- **`exclude-req`-Ventil:** kuratierte Kante, die mit der Schicht-Struktur synchron
  bleiben muss und selbst driften kann (wie `matrix.exclude-sections`).
- **Konsumenten-Vorarbeit:** die Vorwärts-Sicht (§27.1) muss auf konkrete IDs
  restrukturiert werden (Wildcards/Prosa raus), bevor das Gate grün wird — nicht
  d-checks Aufgabe, aber Voraussetzung des Realdatenbelegs.
- **Prosa-`Bezug`:** vereinzelte Prosa-`Bezug:`-Zeilen ohne Tabelle bleiben v1-
  Out-of-Scope; der Konsument tabellarisiert oder akzeptiert Nicht-Gatung.
- **Generator-Sequenz:** dieser Slice liefert **nur** das Gate; der Generator
  (Vorwärts-RTM aus Rück-Kanten erzeugen) ist eine eigene spätere CR
  ([ADR-0038](../../adr/0038-trace-cross-consistency.md) Entscheidung 7).

## 5. Trigger

Auftraggeber-Befund grid-gym (Trigger 088, ADR 0080 §4.4 iii): die §27.1-Zeile
einer Architektur-Anforderung nannte `{GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN}`,
während die `Bezug`-Rück-Kanten derselben Anforderung
`{GG-AR-P-005, GG-AR-P-009, GG-AR-COMP-SCHED}` kamen — Schnittmenge null, von
keinem Gate bemerkt; `--require-complete` ist mangels Fähigkeit bewusst nicht
verdrahtet.

## 6. Sub-Area-Modus-Begründung

GF für die additive `trace.cross-consistency`-Config und den Abgleich-Pfad: Der
neue Vertrag führt, Implementierung folgt. Bestehende `trace`-Pfade sind
Kompatibilitätsbaseline und bleiben durch byte-identische Tests geschützt.

## 7. Closure-Notiz (nach `done/`)

_Ausstehend — wird bei Abschluss mit Commit-Hash, Review-Verdikt und Lerneintrag
gefüllt._
