# ADR-0023 — Immutabilität als Content-Pin: das Modul `immutable` prüft den Core hermetisch im Arbeitsbaum

**Status:** Proposed
**Datum:** 2026-06-28
**Autor:** pt9912
**Bezug:** [`DC-FA-IMM-001`](../../../spec/lastenheft.md#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in)
(opt-in Modul `immutable`); Mechanik-Präzedenz [ADR-0020](0020-content-pin-fence-ausnahme.md)
(Content-Pin/Normalisierung `pins`); Schwester-Gate [ADR-0016](0016-adr-immutable-gate.md)
(git-basiertes `adr-check`); Hermetik-Grenze [ADR-0008](0008-reparatur-ableitbarkeit.md)
(git-Historie außerhalb des read-only-Baums).
**Schärft:** die neue `immutable`-Algorithmus-Sektion
[`spec/spezifikation.md` §DC-FA-IMM-001.a](../../../spec/spezifikation.md#dc-fa-imm-001a--immutabilitäts-pin-gegen-core-drift-immutable)
(Core-Definition, Normalisierung, Grund-Code `core-drift`).

## Kontext

Die Immutabilität von `Accepted`-ADRs erzwingt heute das Shell-Skript
`adr-immutable-check.sh` ([ADR-0016](0016-adr-immutable-gate.md)): es vergleicht
über eine git-Range den **Core** jeder geänderten Accepted-ADR
(`core(BASE)` vs. `core(HEAD)`, Core = Datei ohne `## Geschichte` und ohne die
Status-Zeile). Das ist eine **harte** Garantie — fängt jede Körper-Änderung am
Integrationspunkt —, aber als **Skript** trägt es das Copy-Drift-Problem der
ganzen Werkzeug-Familie: nutzbar in Schwester-Repos nur per **Kopie**, die
auseinanderläuft (dieselbe Klasse, die
[`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)
für `doc-check` mit „verteilen statt kopieren" gelöst hat). Der Auftraggeber
fragt: lässt sich die Immutabilitäts-Prüfung wie die übrigen Module über das
gepinnte Image **verteilen**?

Die Spannung: d-check ist ein **read-only-Arbeitsbaum-Scanner** — hermetisch,
`--network none`, **kein git** ([ADR-0008](0008-reparatur-ableitbarkeit.md):
das distroless-Image enthält kein git, `.git` ist im read-only-Mount nicht
garantiert; eine git-historienbasierte Prüfung „bräuchte einen neuen
nicht-hermetischen Port, analog `external`"). Der git-Diff von `adr-check` passt
**nicht** in dieses Paradigma. Die entscheidbare, hermetische Frage lautet
stattdessen — analog [ADR-0020](0020-content-pin-fence-ausnahme.md): **ist der
Core einer Datei seit einem gesetzten Pin unverändert?** Ein Hash-Vergleich im
Arbeitsbaum, ohne git.

## Entscheidung

Ein opt-in Modul `immutable` (Default aus) mechanisiert die Immutabilität als
**Content-Pin**: eine Datei trägt `<!-- immutable: sha256:<hex> -->`; d-check
hasht ihren **whitespace-normalisierten Core** (Datei **ohne** die Marker-Zeile
und **ohne** die per `immutable.exclude-sections` benannten Abschnitte — für
ADRs `Geschichte`) und vergleicht mit dem Pin. Abweichung → `core-drift`.
Normalisierung und SHA-256 sind dieselben wie bei
[`pins`](0020-content-pin-fence-ausnahme.md). **Diagnose-only**, **opt-in pro
Datei** (nur gepinnte Dateien zählen), default-off byte-identisch. Verteilung
im gepinnten Image, Konsum über `doc-check` — **kein kopiertes Skript**.

**Zwei Backends, bewusst koexistent.** Die Pin-Form ist die hermetische,
verteilbare Hälfte; ihre Garantie ist **schwächer** als der git-Diff (der Pin
ist neu-pinn-bar — wer den Core ändert *und* neu pinnt, kommt durch; der
Reviewer ist der Boden, genau wie bei `pins`/`versions`). Umgekehrt ist der Pin
an **einer** Stelle **strenger** als das Skript: die `**Status:**`-Zeile liegt
stets im Core und ist per `exclude-sections` nicht ausnehmbar (das Ventil greift
auf Abschnitte, nicht auf Einzelzeilen) — ein vom Skript erlaubter
Supersede-Übergang berührt damit den Core und verlangt ein bewusstes Neu-Pinnen;
die feinere „nur die Kopf-Status-Zeile strippen"-Semantik bleibt dem späteren
VCS-Adapter (oder einer Folge-CR) vorbehalten. Die **harte**
git-Garantie bleibt das bestehende `adr-check`
([ADR-0016](0016-adr-immutable-gate.md)) — es bleibt **unangetastet** und
bewacht d-checks eigene 21 Accepted-ADRs weiter über die Range (der Produzent
hat das Skript lokal, kein Copy-Drift). Schwester-Repos, die kein git-Gate
pflegen wollen, bekommen mit `immutable` die verteilbare Stufe.

**Die git-Stufe ist nicht verworfen, nur vertagt.** Eine künftige opt-in-Stufe
(`core(BASE)` vs. `core(HEAD)`) wäre — wie [ADR-0008](0008-reparatur-ableitbarkeit.md)
es vorzeichnet — ein **neuer nicht-hermetischer Driven-Port** (ein `git`-Adapter
neben `httpcheck`/`external`) als **eigene Anforderung/ADR**. Diese ADR
entscheidet bewusst **nur** die hermetische Pin-Hälfte; der VCS-Adapter ist
Out-of-Scope.

**Eigenes Modul** (nicht in `pins` verschmolzen): andere Bindung (Selbst-Pin der
Datei statt eines Link-Ziel-Spans), andere Semantik (Immutabilität statt
Zitat-Frische), je einzeln opt-in und testbar.

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **Content-Pin `immutable` (gewählt)** | hermetisch/offline/deterministisch; verteilbar im Image (löst Copy-Drift); reuse `pins`-Normalisierung; opt-in pro Datei | Garantie schwächer (neu-pinn-bar) → Reviewer als Boden |
| git-Diff ins Modul ziehen (jetzt) | volle Range-Garantie | bräche den hermetischen read-only-Scan-Vertrag ([ADR-0008](0008-reparatur-ableitbarkeit.md): kein git im Image, `.git` nicht garantiert); nur als nicht-hermetischer Port — eigene Anforderung |
| Skript belassen, nicht mechanisieren | volle Garantie, schon da | Copy-Drift in Schwester-Repos bleibt ([`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Problem ungelöst) |
| `pins` erweitern (Selbst-Link mit `dpin`) | kein neues Modul | Core-Ausnahme (`Geschichte`/Marker) nicht ausdrückbar; Selbst-Link ist ein Hack |
| Default-on / Pflicht-Pin | maximale Abdeckung | nicht jede Datei ist eingefroren → Lärm; bräche Abwärtskompatibilität |

## Fitness Function

- Ohne aktives `immutable` ist der Befundsatz byte-identisch (opt-in-Selbsttest).
- Read-only: das Werkzeug schreibt nichts; Pinnen/Neu-Pinnen ist menschlich.
- Determinismus: gleicher Baum + gleicher Pin ⇒ gleicher Befundsatz.
- Reflow-Boundary-Test (nur-Whitespace-Änderung am Core → kein `core-drift`).
- Ausgenommener-Abschnitt-Test (Anhang unter `exclude-sections` → kein
  `core-drift`); Core ohne die Marker-Zeile (kein Selbstbezug).
- Negative: inhaltliche Core-Änderung außerhalb der Ausnahme → `core-drift`.
- Kein `--repair`-Hunk für `core-drift`.
- Kein git: der Lauf bleibt im read-only-Arbeitsbaum
  ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-06-28 | Entwurf nach Auftraggeber-Diskussion (adr-check ist nur ein Skript → Copy-Drift; Hexagon-Analyse: git-Diff = nicht-hermetischer Port analog `external`, Pin-Form = hermetische Hälfte). Zwei-Backend-Entscheid: `immutable`-Pin jetzt, VCS-Adapter vertagt. Begleitet slice-052. Status Proposed. |
