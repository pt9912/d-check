# ADR-0019 — Versions-Pin-Fence-Ausnahme: das Modul `versions` liest Pins auch in Fences

**Status:** Proposed
**Datum:** 2026-06-24
**Autor:** pt9912
**Bezug:** [`DC-FA-VER-001`](../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
(opt-in Modul `versions`), Mechanik-Präzedenz der opt-in-Module
[ADR-0018](0018-diagram-fence-ausnahme.md) (Fence-Öffnung für `diagrams`) und
[ADR-0010](0010-semgrep-hermetisches-gate.md) (distroless/netzlos).
**Schärft:** [`spec/spezifikation.md` §DC-FA-LINK-001.a](../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
(Fence-Opazität der Vorverarbeitung) und die neue `versions`-Algorithmus-Sektion
(`.a`) der Spezifikation.

## Kontext

Versionsnummern werden über die Doku hart gepinnt — vor allem als
`ghcr.io/…/d-check:vX.Y.Z` in Kommando-Beispielen (Spike 2026-06-24: ~18× im
Benutzerhandbuch + README). Beim Release veralten sie still: ein vergessener Pin
bleibt unbemerkt. Strukturell fängt `d-check` das heute nicht — die
Markdown-Vorverarbeitung entfernt Fenced-Code für **alle** Module
([`spec/spezifikation.md` §DC-FA-LINK-001.a](../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)),
und die Pins leben genau dort.

Zwei naheliegende Wege wurden verworfen. Ein kopiertes Meta-Gate (ein
`version-consistency`-Shell-Skript) würde über die Repo-Familie per Datei-Kopie
driften — exakt die
[`MR-003`](../../../harness/conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh)→[`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Lektion
(vendorter Sensor → durch Dogfooding im gepinnten Werkzeug ersetzt); ein neues
kopiertes Skript wäre ein Rückschritt gegen den Daseinszweck von `d-check`. Der
Git-Tag als Wahrheitsquelle scheidet aus: er liegt außerhalb des read-only
gemounteten Baums und bräche
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus).

Die `d-check`-förmige, entscheidbare Frage ist: **stimmen alle Versions-Pins mit
der deklarierten aktuellen Version überein?** — ein Gleichheits-Vergleich, kein
semantisches Urteil. „Alle Pins == aktuelle Version" ist das Problem jedes Repos
mit versioniertem Artefakt, also eine generische Regel, kein Repo-Spezifikum.

## Entscheidung

Ein opt-in Modul `versions` (Default aus). Es liest Versions-Pins (Muster
`versions.pin-pattern`) **auch innerhalb von Fenced-Code** und vergleicht jeden
Treffer mit der aus `versions.current-from` (Default `version.md#aktuell`)
gelesenen aktuellen Version; Abweichung → Befund `version-stale`. Die
Fence-Öffnung ist eine bewusste, auf das `pin-pattern` gescopte Ausnahme —
analog zu [ADR-0018](0018-diagram-fence-ausnahme.md), das Fences für `diagrams`
öffnet, hier aber **breiter** (Muster-Scan über alle Fences statt nur gelistete
Diagramm-Fences): deshalb diese eigene ADR. Kein Sprach-Parser, nur ein
Muster-Scan. Verteilt wird die Regel wie der Rest des Werkzeugs — im gepinnten
Image, konfiguriert über `.d-check.yml`, Make-Glue über das `--print-mk`-Fragment,
nicht als kopiertes Skript.

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **Modul `versions` (gewählt)** | im gepinnten Werkzeug verteilt, kein Copy-Drift; generisch für jedes Repo mit versioniertem Artefakt | Fence-Öffnung breiter als `diagrams` (Muster-Scan über alle Fences) |
| kopiertes `version-consistency`-Shell-Skript | schnell, kein Release nötig | Copy-Drift über die Repo-Familie (s. Kontext); Rückschritt gegen den `d-check`-Zweck |
| Pins nur außerhalb von Fences zulassen | keine Fence-Öffnung nötig | unrealistisch — Kommando-Beispiele müssen in Fences stehen |
| Git-Tag als Wahrheitsquelle | „echte" Release-Wahrheit | außerhalb des read-only-Baums, bräche den Determinismus-Kernvertrag |

## Fitness Function

- Ohne `versions`-Block ist der Befundsatz byte-identisch (opt-in-Selbsttest, wie
  `diagrams`/`hostpaths`).
- Read-only: das Werkzeug schreibt nichts; das Register pflegt der Mensch.
- Determinismus: gleicher Baum + gleiche Konfiguration ⇒ gleicher Befundsatz.
- Der Fence-Scan greift ausschließlich für das konfigurierte `pin-pattern` — alle
  anderen Module behandeln Fences unverändert opak.
- Negativ-Test: ein stale Pin **innerhalb eines Fenced-Blocks** wird gemeldet;
  historische Pins (`exempt-paths`/`d-check:ignore`) nicht.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-06-24 | Entwurf (Idee 1; Spike-Befund ~18 Fenced-Pins; kopiertes Meta-Gate wegen Copy-Drift verworfen). Status Proposed. |
