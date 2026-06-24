# ADR-0020 — Content-Pin: das Modul `pins` hasht den rohen, normalisierten Ziel-Span

**Status:** Accepted
**Datum:** 2026-06-24
**Autor:** pt9912
**Bezug:** [`DC-FA-PIN-001`](../../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in)
(opt-in Modul `pins`); Mechanik-Präzedenz [ADR-0018](0018-diagram-fence-ausnahme.md)
(Fence-Öffnung `diagrams`), [ADR-0019](0019-versions-pin-fence-ausnahme.md)
(Fence-Lesen `versions`), [ADR-0010](0010-semgrep-hermetisches-gate.md)
(distroless/netzlos).
**Schärft:** [`spec/spezifikation.md` §DC-FA-LINK-001.a](../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
(Fence-Opazität der Vorverarbeitung) und die neue `pins`-Algorithmus-Sektion (`.a`).

## Kontext

d-check fängt strukturellen Drift (`target-missing`/`anchor-missing`), aber nicht
**inhaltlichen**: ein Link löst weiter auf, während der Ziel-Abschnitt sich seit
dem Verlinken geändert hat — die zitierende Aussage daneben veraltet still
("stale citation"). Spike 2026-06-24: 86/116 zitierte Sections wurden nach Anlage
geändert (real). Eine *semantische* Konsistenzprüfung scheidet aus (unentscheidbar,
bricht [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)/[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Die entscheidbare Frage ist: **ist der Ziel-Span seit dem gesetzten Pin
unverändert?** — ein Hash-Vergleich. Roh-Byte-Hashing wäre laut (Reflow); der
Spike zeigt 0/87 kosmetische Trips unter Whitespace-Normalisierung. Und: Drift in
**Code-Beispielen** ist ein Zielfall → der Span muss **roh inkl. Fences** gehasht
werden, gegen die Fence-Opazität der Vorverarbeitung.

## Entscheidung

Ein opt-in Modul `pins` (Default aus). Ein Link kann einen Content-Pin tragen
(`<!-- dpin: sha256:<hex> -->`, an den unmittelbar vorausgehenden Link derselben
Zeile gebunden); d-check hasht den **rohen, whitespace-normalisierten** Ziel-Span
(ganze Datei oder Heading-Section) und vergleicht. Mismatch → `link-stale`.
**Diagnose-only** (kein Auto-Re-Pin — das wäre stilles Segnen; Re-Pinnen bleibt
menschlich, `--bless` spätere CR). Nur auflösbare Links; der strukturelle Befund
bleibt `links`/`anchors`. Die Fence-Öffnung ist die bewusste Ausnahme — wie
[ADR-0018](0018-diagram-fence-ausnahme.md)/[ADR-0019](0019-versions-pin-fence-ausnahme.md),
hier über den **ganzen Ziel-Span** (kein Parser, nur Hash). Verteilung im
gepinnten Image, Konsum über `doc-check` — kein kopiertes Skript.

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **Content-Pin (gewählt)** | entscheidbar/offline/deterministisch; reuse Section-Logik; opt-in pro Link = Zitat-Intent | Pin-Pflege manuell |
| semantischer/LLM-Check | „echte" Zitat-Konsistenz | unentscheidbar, bricht den Determinismus-/Netzlos-Kernvertrag |
| Byte-Hash ohne Normalisierung | maximal einfach | laut bei Reflow (Spike: Rauschen, das die Normalisierung vermeidet) |
| Manifest-Datei statt inline-Pin | weniger Prosa-Clutter | Pin und Inhalt driften in getrennten Diffs, kein lokaler Kontext |
| Default-on / Pflicht-Pin | maximale Abdeckung | Pins sind Zitat-Intent, nicht jeder Link ist ein Zitat → Lärm |

## Fitness Function

- Ohne aktives `pins` ist der Befundsatz byte-identisch (opt-in-Selbsttest).
- Read-only: das Werkzeug schreibt nichts; Pin-Pflege ist menschlich.
- Determinismus: gleicher Baum + gleiche Pins ⇒ gleicher Befundsatz.
- Reflow-Boundary-Test (nur-Whitespace-Änderung am Span → kein `link-stale`).
- Marker-Bindung deterministisch (Ambiguitäts-Tests); nicht auflösbares Ziel →
  kein `pins`-Befund (auch pins-only), kein Doppelbefund.
- Kein `--repair`-Hunk für `link-stale`.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-06-24 | Entwurf nach Spike (Idee 2; 86/116 Sections driften, 0/87 kosmetisch) + Plan-Review R1 (Marker-Bindung / Ziel-weg / Modul-Scope geschärft). Status Proposed. |
| 2026-06-24 | Angenommen mit der slice-049-Closure: Plan-Review R1→R3 + unabhängiges Impl-Review (alle behoben), Modul `pins` implementiert + getestet, `make gates` grün. Status Accepted. |
