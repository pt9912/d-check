# Slice slice-052: Modul `immutable` — Immutabilitäts-Pin gegen Core-Drift

**Status:** in-progress (welle-41-immutable).

**Welle:** welle-41-immutable (Trigger: Auftraggeber — „adr-immutable-check ist
nur ein Skript (Copy-Drift); können wir die Immutabilitäts-Prüfung wie die
übrigen Module verteilen?" → Hexagon-Analyse: git-Diff = nicht-hermetischer Port
analog `external`, Content-Pin = hermetische, im Arbeitsbaum entscheidbare Hälfte).

**Bezug:** Führt eine **neue IMM-Anforderung** im Lastenheft ein (Modul
`immutable`, Content-Pin/`core-drift`) plus einen **begleitenden ADR** (die
Mechanisierungs-Strategie: hermetischer Pin jetzt, VCS-Adapter vertagt). Mechanik
erbt die `pins`-Normalisierung ([ADR-0020](../../adr/0020-content-pin-fence-ausnahme.md));
Schwester-Gate ist das git-basierte `adr-check`
([ADR-0016](../../adr/0016-adr-immutable-gate.md)), das **unangetastet** bleibt.
Verteilung wie der Rest des Werkzeugs (gepinntes Image, kein kopiertes Skript —
[`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Linie).

**Autor:** pt9912. **Datum:** 2026-06-28.

---

## 1. Ziel

Die Immutabilität von `Accepted`-ADRs erzwingt heute `adr-immutable-check.sh`
([ADR-0016](../../adr/0016-adr-immutable-gate.md)) über einen **git-Diff**
(`core(BASE)` vs. `core(HEAD)`). Harte Garantie, aber als **Skript** nur per
Kopie in Schwester-Repos nutzbar (Copy-Drift). **Neu:** dieselbe Disziplin als
verteilbares, **hermetisches** Modul — eine Datei trägt
`<!-- immutable: sha256:<hex> -->`, d-check hasht ihren normalisierten **Core**
(Datei ohne Marker-Zeile + ohne `exclude-sections`) und meldet Abweichung als
`core-drift`. Kein git, rein im read-only gescannten Arbeitsbaum
([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
die Pin-Disziplin ist die des Image-Digests (pinnen → bei legitimer Änderung neu
segnen), auf Doku-Immutabilität angewandt.

## 2. Entscheidungen

- **Hermetisch statt git.** d-check ist ein read-only-Arbeitsbaum-Scanner (kein
  git im Image, `.git` im Mount nicht garantiert —
  [ADR-0008](../../adr/0008-reparatur-ableitbarkeit.md)). Die entscheidbare,
  hermetische Frage ist „Core == Pin?" (Hash-Vergleich), nicht „hat sich der Core
  über eine Commit-Range geändert?" (git-Diff). Letzteres bräuchte einen neuen
  nicht-hermetischen Port (ein `git`-Adapter neben `httpcheck`/`external`) und
  bleibt **Out-of-Scope** — eigene spätere Anforderung
  ([ADR-0023](../../adr/0023-immutable-core-pin.md)).
- **Zwei Backends koexistieren.** Die Pin-Garantie ist **schwächer** (neu-pinn-bar
  → Reviewer als Boden, wie `pins`/`versions`); die harte git-Garantie bleibt
  `adr-check` ([ADR-0016](../../adr/0016-adr-immutable-gate.md)), **unangetastet**.
  Schwester-Repos ohne git-Gate bekommen mit `immutable` die verteilbare Stufe;
  d-check selbst behält sein lokales Skript (Produzent — kein Copy-Drift).
- **Core = Datei ohne Marker-Zeile + `exclude-sections`.** Die Marker-Zeile fällt
  raus (sonst Selbstbezug); `exclude-sections` (für ADRs `[Geschichte]`) nimmt
  legitime Anhänge aus — dieselbe Section-Abgrenzung wie `matrix.exclude-sections`.
- **Normalisierung whitespace-/reflow-invariant**, SHA-256 — **dieselbe** wie
  `pins` ([ADR-0020](../../adr/0020-content-pin-fence-ausnahme.md)); roh inkl.
  Fenced-Code (Code-Beispiele im Core zählen).
- **Marker auf der vorverarbeiteten Zeile** (Fence-/Inline-Code-Vorkommen sind
  inert) — ein Syntax-Beispiel in der Doku wird **kein** Live-Pin.
- **opt-in pro Datei** (nur gepinnte Dateien zählen), strikt opt-in Modul
  (default-off byte-identisch); **diagnose-only** (kein `--repair`-Hunk —
  Neu-Pinnen ist menschlich).
- **Eigenes Modul `immutable`** (12.), nicht in `pins` verschmolzen: andere
  Bindung (Selbst-Pin der Datei statt Link-Ziel-Span), andere Semantik
  (Immutabilität statt Zitat-Frische).

## 3. Definition of Done

- [x] **Spec:** neue IMM-Anforderung
  [`DC-FA-IMM-001`](../../../../spec/lastenheft.md#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in)
  (Bereichskürzel `IMM` in §3, Versions-Bump 0.32.0 + §7-Historie, `immutable` in
  [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
  + Glossar) + begleitender [ADR-0023](../../adr/0023-immutable-core-pin.md) +
  spezifikation-`.a`-Algorithmus-Sektion + Grund-Code `core-drift` (§4) +
  Schema-Key `immutable.exclude-sections` (§2) + ADR-Index; doc-first vor Code.
- [ ] **Code:** Modul `immutable` (`rules/immutable.go`: Marker-Erkennung auf der
  vorverarbeiteten Zeile, erster Marker je Datei; Core = roh ohne Marker-Zeile +
  `exclude-sections`-Abschnitte; `pins`-Normalisierung + SHA-256; Vergleich →
  `core-drift`), Wiring in `run.go`, `model.ImmutableConfig` +
  `validModules()` + `configyaml`-Kompilierung (`exclude-sections`), Grund-Code in
  `diagnose.go`. Tests: Happy/Reflow-Boundary/ausgenommener-Abschnitt/Negative/
  kein-Marker/Modul-aus; Marker in Fence inert; erster-Marker-zählt.
- [ ] `make gates` grün; zwei unabhängige Reviews; Closure (Move nach `done/` +
  Roadmap-Flip, [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise))
  + Release v0.32.0.

## 4. Risiken / offene Punkte

- **Normalisierung ist ein Vertrag** (definiert „inhaltliche Core-Änderung") —
  geerbt von `pins`, dort spike-kalibriert; konservativ.
- **Garantie-Stufe ehrlich benennen:** der Pin ist neu-pinn-bar — wer den Core
  ändert *und* neu pinnt, kommt durch; der Reviewer ist der Boden (Doku im ADR).
  Die harte git-Garantie bleibt bewusst beim Schwester-Gate `adr-check`.
- **Status-Übergänge:** die Pin-Form sieht nur „Core unverändert/nicht". Wird die
  ausgeschlossene Status-Zeile **nicht** über `exclude-sections` abgedeckt (sie
  ist eine Zeile, kein Abschnitt), berührt ein Status-Wechsel den Core nur, wenn
  die Status-Zeile im Core liegt → bewusstes Neu-Pinnen. Für ADRs wird daher
  **nach Accept** gepinnt; ein späteres Supersede verlangt ein Neu-Pinnen
  (dokumentierte Grenze; die feinere „nur-Status-Zeile-strippen"-Semantik des
  Skripts ist eine mögliche Folge-CR, kein Teil von v1).
- **Kein Repo-Dogfooding der 21 Alt-ADRs:** ein Marker in eine `Accepted`-ADR zu
  setzen, löste `adr-check` selbst aus (Körper-Edit) — dieselbe Immutability-Falle
  wie slice-051. v1 liefert das Modul **un-dogfooded** aus (wie `pins`); Adoption
  (ADRs markieren / Skript ablösen) ist eine separate Entscheidung.

## 5. Trigger

Auftraggeber 2026-06-28: „Wir machen es regelkonform" — nach Design-Diskussion
(adr-check als Skript → Copy-Drift; Hexagon: git-Diff = nicht-hermetischer Port,
Content-Pin = hermetische Hälfte; „beide Wege anbieten" → Slice A jetzt, VCS-
Adapter als Slice B vertagt).

## 6. Sub-Area-Modus-Begründung

GF (Produkt-Code + Spec; „Doc führt, Code folgt"). Keine BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

_(folgt bei Closure)_
