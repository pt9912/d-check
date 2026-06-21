# Verifikation R2 — slice-036 LOW-1-Fix (`--trace` RTM, `DC-FA-CLI-009`)

## Kopf-Metadaten

- **Review-Art:** Unabhängige Verifikation/Reviewer-Runde R2 — adversariale
  Gegenprüfung des LOW-1-Fixes aus R1
  ([`docs/reviews/2026-06-21-slice-036-rtm-trace.md`](./2026-06-21-slice-036-rtm-trace.md)),
  Regressionssuche durch den Fix sowie Stichproben-Gegenprüfung der
  R1-Kernaussagen (Ableitung, Determinismus, read-only). Frischer Kontext,
  alles selbst reproduziert. Gates real ausgeführt (read-only, Docker).
- **Datum:** 2026-06-21
- **Gegenstand:** Fix-Commit `75c870a` (LOW-1) auf der Implementierung
  `9ef67fa`; geänderte Funktion `traceTitle`
  (`internal/hexagon/core/app/trace.go:125`), neuer Test
  `TestCLI036_Trace_BacktickHeading`
  (`internal/adapter/driving/cli/cli_acceptance_test.go:1316`).
- **Reviewer:** Claude (Agent), Skill
  [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Beleg-Setup:** HEAD `75c870a`, Image `d-check:latest` aus HEAD neu gebaut
  (`make build` ⇒ identische Image-ID, also deckungsgleich mit dem committeten
  Stand). Für den Vorher/Nachher-Vergleich wurde aus `9ef67fa` über die
  Projekt-Dockerfile ein Pre-Fix-Image gebaut (Worktree, danach entfernt; das
  d-check-Repo blieb unverändert). Test-Repos liegen außerhalb des Repos unter
  `/tmp`. Realer Aufruf:
  `docker run --rm -v "<tmp>:/repo:ro" d-check:latest --trace [--json|--yaml]`.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/in-progress/slice-036-rtm-trace.md`; Anforderung
  `DC-FA-CLI-009` (RTM, Titel-AK), `DC-QA-02` (Determinismus), `DC-QA-03`
  (read-only/netzlos); R1-Report (LOW-1, INFO-1..3).

## Verifikations-Kernaussage

**LOW-1 ist bestätigt behoben** — aber der Fix führt eine **neue Regression
derselben Klasse** ein (LOW): Titel, die unmittelbar nach der Kennung mit
einem Inline-Code-Span beginnen, verlieren ihren führenden Backtick.

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

#### LOW-2 (neu) — `traceTitle` strippt den führenden Backtick eines Titel-Code-Spans (Fix-Regression)

- **Kategorie:** LOW
- **Quelle:** `DC-FA-CLI-009` (AK „eine Zeile … mit … Titel"); Maintainability
- **Pfad:** `internal/hexagon/core/app/trace.go:130`
- **Befund:** Der Fix nimmt den Backtick in die abschließende
  `strings.TrimLeft(rest, …)`-Zeichenklasse auf. Beginnt der Titel unmittelbar
  nach Kennung/Trenner mit einem Inline-Code-Span, wird dessen öffnender
  Backtick mit abgeräumt; der schließende Backtick bleibt als Waise stehen.
  Vorher/Nachher-Beleg über identisches Repo (Pre-Fix-Image `9ef67fa` vs.
  HEAD-Image `75c870a`), Heading und resultierender JSON-`title`:

  ```text
  Heading (bare-ID, Titel beginnt mit Code-Span):
    ### DC-FA-LEAD-001 — `code-at-start` rest of title
  pre-fix  title: `code-at-start` rest of title     (korrekt)
  post-fix title: code-at-start` rest of title      (führender Backtick weg,
                                                     schließender verwaist)

  Realistische Form 1 (Titel startet mit Modul-Code):
    ### DC-FA-MOD-001 — `links`-Modul prüft lokale Referenzen
  post-fix title: links`-Modul prüft lokale Referenzen

  Realistische Form 2 (Titel startet mit Flag-Code):
    ### DC-FA-MOD-002 — `--flag` aktiviert das Verhalten
  post-fix title: flag` aktiviert das Verhalten   (auch das -- mitgestrippt,
                                                   da - in der Trim-Klasse)
  ```

  Wirkungsgleich zum ursprünglichen LOW-1: unsauberer Titel im
  präfix-agnostischen Fremd-Repo-Einsatz. Keine Auswirkung auf d-checks eigene
  Daten — kein eigenes Anforderungs-Heading beginnt den Titel mit einem
  Code-Span direkt hinter dem Trenner (geprüft mit
  `grep -nE` über `spec/lastenheft.md`, keine Treffer); Code-Spans stehen dort
  stets *später* im Titel und bleiben unberührt (Beleg: das Heading
  `### DC-FA-LINK-001 — Lokale Link- und Bildreferenzen (Modul links)` mit
  echtem Code-Span ergibt post-fix den unveränderten Titel inklusive Span).
- **Verifizierbar:** ja —
  `docker run --rm -v "<repo, Anforderungs-Titel beginnt mit Code-Span>:/repo:ro"
  d-check:latest --trace --json` zeigt den führenden Backtick im `title`-Feld
  gestrippt und einen verwaisten schließenden Backtick.

### INFO

#### INFO-4 — Em-Dash ohne umgebende Leerzeichen verhindert die Anforderungs-Erkennung (vorbestehend, nicht durch den Fix)

- **Kategorie:** INFO
- **Quelle:** Spezifikation §DC-FA-CLI-009.a Pkt. 1 (führende Heading-Kennung)
- **Pfad:** `internal/hexagon/core/app/trace.go:97` (`strings.Fields`)
- **Befund:** `traceRequirements` ermittelt die Kennung aus dem ersten Feld
  nach `strings.Fields`. Klebt ein Trenner ohne Leerzeichen an der Kennung
  (`### DC-FA-NOSP-001—NoSpaceEmDash`), bleibt das erste Feld
  `DC-FA-NOSP-001—NoSpaceEmDash` und scheitert an `isFullReqID` (exaktes
  Voll-Match) — die Anforderung wird **gar nicht** erkannt. Beleg: in einem
  Repo mit genau diesem Heading plus einem normalen `### DC-FA-NOSP-002 — Normal`
  liefert `--trace --json` nur das normale (`total:1`). Berührt den Fix nicht
  (greift vor `traceTitle`); rein vorbestehend. d-checks eigene Headings setzen
  den Em-Dash stets mit Leerzeichen, daher ohne praktische Wirkung.
  Dokumentationswürdige Annahme.

## Negativbefunde (geprüft, ohne Befund)

- **LOW-1 behoben (Backtick-Kennung im Heading):** Vorher/Nachher über
  identisches Repo (Heading-Shape ⇒ JSON-`title`):

  ```text
  ### `DC-FA-BTK-001` — LOW-1 Backtick-Kennung
    pre-fix:  Kennung + Backticks im Titel erhalten (unsauber)
    post-fix: LOW-1 Backtick-Kennung   (sauber)

  ### `DC-FA-BTK-007`            (nur Kennung, kein Titel)
    pre-fix:  `DC-FA-BTK-007`     (Kennung + Backticks als Titel)
    post-fix: ""                  (leer, kein Absturz, code=0)
  ```

  **Bestätigt behoben.**
- **Bare-ID-Headings nicht gebrochen:** Em-Dash (`DC-FA-BARE-002`) und
  Boundary „bare-ID ohne Titel" (`DC-FA-BARE-008`) liefern pre- und post-fix
  denselben Titel (`Bare Em-Dash` bzw. leer). Im breiten Lauf sind
  Em-Dash/Hyphen/Doppelpunkt bei bare-ID alle sauber extrahiert.
- **Gemischte Trenner mit Backtick-Kennung:**

  ```text
  ### `DC-FA-BTK-005` - …    (Hyphen)       ⇒ title: Backtick Hyphen
  ### `DC-FA-BTK-006`: …     (Doppelpunkt)  ⇒ title: Backtick Doppelpunkt
  ```

  Sauber.
- **Eingebetteter Code-Span MITTEN im Titel (Regressions-Abgrenzung):**

  ```text
  ### DC-FA-EMB-009 — Titel mit `embedded` darin
    pre-fix == post-fix: Titel mit `embedded` darin  (Backticks erhalten)
  ```

  Die Regression (LOW-2) trifft ausschließlich den *führenden* Backtick direkt
  hinter Kennung/Trenner, nicht spätere Spans.
- **Determinismus (DC-QA-02):** 5× Markdown über das Edge-Repo
  (Backtick-/Misch-Fälle) byte-identisch (`sha256` konstant); 3× JSON und 3×
  YAML über das Ableitungs-Repo je byte-identisch. Kein Map-Leak.
- **read-only / netzlos (DC-QA-03):** `--network none` + `:ro`-Mount läuft
  (Edge-Repo, `total:12`). Datei-Snapshot (`sha256sum` aller Dateien) vor und
  nach 3 Läufen unverändert; keine `.d-check.yml`/`.tmp` erzeugt.
- **Ableitungs-Kern (Stichprobe):** Anforderungen vs. Nicht-Anforderungen
  korrekt getrennt — erkannt: `DC-FA-FOO-001`, `DC-FA-FOO-002`, `DC-QA-01`;
  ausgeschlossen: `DC-FA-CLI-009.a` (`.a`-Variante),
  `## Anforderung zu DC-FA-MID-001` (mid-heading), `### ADR-0099`,
  `### MR-006`, `### slice-100` (Nicht-Anforderungs-Präfixe).
- **Waisen-Logik (Waise = kein Slice; ADR-only = Waise):** `DC-FA-FOO-002`
  (nur ADR) ⇒ `orphan:true`; `DC-QA-01` (nur Slice) ⇒ `orphan:false`;
  `DC-FA-FOO-001` (ADR + Slice) ⇒ `orphan:false`. `orphans:1` korrekt.
- **Referenz-Sammlung + README-Skip:** `DC-FA-FOO-001` wird auch in der
  `README.md` im ADR-Verzeichnis referenziert, doch README erscheint in keiner
  Owner-Zelle (kein Datei-Shape-Treffer); Owner sind ausschließlich
  `ADR-0099` und `slice-100` (aus den Dateinamen abgeleitet).
- **INFO-1 (R1) reproduziert:** `DC-FA-GHOST-777` (in ADR referenziert, nicht
  im Lastenheft definiert) erscheint nicht in der RTM — dangling reference
  still verworfen, wie in R1 als bewusste Annahme dokumentiert.
- **Neuer Test guardt den Fix:** `TestCLI036_Trace_BacktickHeading` erwartet
  den sauberen Titel aus einem Backtick-Kennungs-Heading; pre-fix produziert
  die Kennung-behaftete Variante, also failt der Test ohne den Fix — er
  sichert LOW-1 echt ab.
- **Regression Gesamt:** `make test` grün (alle Pakete `ok`, inkl. `cli`-Paket
  mit dem neuen Test). `make doc-check` grün (Dogfooding über die Repo-Wurzel,
  0 Befunde) — der Fix bricht das Dogfooding nicht.
- **Repo unberührt:** keine Schreibzugriffe in den Quell-/Spec-/Doc-Baum;
  HEAD bleibt `75c870a`, Worktree nach dem Vergleich entfernt
  (`git status` sauber).

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 1 | 1 |

(LOW-2 neu durch den Fix; INFO-4 neu, vorbestehend. Das ursprüngliche LOW-1
ist behoben und entfällt; INFO-1..3 aus R1 bleiben unberührt gültig.)

## Verdikt

**Closure-fähig (kein Blocker).** LOW-1 ist mit Vorher/Nachher-Beleg über
dasselbe Repo **bestätigt behoben** (Backtick-Kennung im Heading liefert jetzt
den sauberen Titel; bare-ID-Headings, gemischte Trenner und eingebettete
Code-Spans mitten im Titel bleiben intakt; ID-only-Heading ⇒ leerer Titel ohne
Absturz). Determinismus (DC-QA-02), read-only/netzlos (DC-QA-03) und der
Ableitungs-Kern (Anforderungs-Erkennung, Waisen-Semantik, Referenz-Sammlung
mit README-Skip) sind unabhängig reproduziert und tragfähig; `make test` und
`make doc-check` sind grün. Der Fix führt **eine neue Regression derselben
Klasse** ein (LOW-2): ein Titel, der direkt hinter der Kennung mit einem
Inline-Code-Span beginnt, verliert seinen führenden Backtick — wirkungsgleich
zum ursprünglichen LOW-1, ohne Effekt auf d-checks eigene Daten, relevant nur
im präfix-agnostischen Fremd-Repo-Einsatz. Da LOW-1 (das ältere Defizit
derselben Klasse) bereits als nicht-blockierend galt, blockiert auch LOW-2 die
Freigabe nicht; INFO-4 (Em-Dash ohne Leerzeichen) ist vorbestehend und ohne
praktische Wirkung. Empfehlung an die Übergabe: LOW-2 als bekannten
Robustheits-Rest führen.
