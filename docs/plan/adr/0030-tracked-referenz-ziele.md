# ADR-0030 — Getrackt-Status von Referenz-Zielen: das Modul `tracked` prüft auflösbare Ziele gegen den git-Index (VCS-Port, ohne Range)

**Status:** Proposed
**Datum:** 2026-07-03
**Autor:** pt9912
**Bezug:** [`DC-FA-TRK-001`](../../../spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)
(opt-in Modul `tracked`); VCS-Port-Fundament [ADR-0024](0024-vcs-immutable-gate.md)
(reine-Go-git, ohne Binary/Netz — dritte Nutzung nach `vcs` und `commits`);
Port-Erweiterungs-Präzedenz [ADR-0027](0027-commits-traceability-modul.md)
(`commits` erweiterte den Port um Message-Lesen, `tracked` erweitert um die
Index-Abfrage); Kein-Doppelbefund-Prinzip [ADR-0020](0020-content-pin-fence-ausnahme.md)
(nur auflösbare Ziele — der strukturelle Befund bleibt beim Struktur-Modul);
Ventil-Vorbild [ADR-0025](0025-codepaths-ignore-refs.md) (referenz-weites
Ziel-Glob); die heute gelebten Workarounds
[`MR-017`](../../../harness/conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)
(gitignore + `scan.ignore`-Doppel) und
[`MR-019`](../../../harness/conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)
(Vendoring für Präsenz auf jedem Checkout) sind Einzelfall-Disziplin — dieses
Modul mechanisiert die Klasse.
**Schärft:** die neue Algorithmus-Sektion
[`spec/spezifikation.md` §DC-FA-TRK-001.a](../../../spec/spezifikation.md#dc-fa-trk-001a--getrackt-status-auflösbarer-referenz-ziele-tracked)
(Index-Menge, Kandidaten, Ventil, Grund-Code `target-untracked`).

## Kontext

d-check ist im Default **git-agnostisch**: die Wahrheit ist der Arbeitsbaum,
`.gitignore` wird nie gelesen (kommt in Code und Spezifikation nicht vor).
Referenziert ein Dokument eine Datei, die nur lokal existiert — untracked oder
gitignoriert —, ist der Lauf beim Erzeuger grün, auf **jedem frischen Klon**
aber rot (`target-missing`): Umgebungs-Drift zwischen Arbeitsbäumen, kein
Determinismus-Bruch je Baum (Auftraggeber-Frage 2026-07-03, per Fixture-Demo
belegt). Die heutigen Ventile fangen das spät oder nur je Einzelfall: die CI
des nächsten Checkouts (Feedback erst nach dem Push), das gitignore+`scan.ignore`-Doppel
und das Vendoring (beide setzen voraus, dass jemand die Falle schon kennt).
Seit [ADR-0024](0024-vcs-immutable-gate.md) existiert ein VCS-Port, der `.git`
read-only in reinem Go liest — die Infrastruktur, um den Getrackt-Status am
**Entstehungsort** zu prüfen, ist vorhanden.

## Entscheidung

Ein opt-in Modul **`tracked`** (16. Modul, Default aus) prüft als Post-Pass die
**Datei-Ebene** der von `links` aufgelösten, **existierenden** repo-internen
Link-/Bild-Ziele gegen den **git-Index**; ein untracktes Ziel ⇒ Grund-Code
`target-untracked`. **Diagnose-only.**

- **Index statt `.gitignore`-Interpretation.** d-check baut keinen zweiten
  ignore-Regel-Interpreter; der Index ist die einzige Wahrheit. Eine frisch
  per `git add` gestagte Datei gilt als getrackt — neue Doku wird mit dem
  Staging grün, der Inner-Loop bleibt arbeitsfähig.
- **VCS-Port, dritte Nutzung — ohne Range.** `tracked` liest ausschließlich
  den Index-Stand (kein `--range`/`--staged`); der Port wird um die
  Index-Abfrage erweitert (wie [ADR-0027](0027-commits-traceability-modul.md)
  ihn um Message-Lesen erweiterte). Eingabe erweitert, aber lokal, lesend,
  deterministisch —
  [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)/[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  in der `vcs`-Lesart.
- **Kein Doppelbefund.** Nicht existierende Ziele bleiben `target-missing`
  (`links`); `tracked` prüft nur, was die Link-Auflösung fand
  ([ADR-0020](0020-content-pin-fence-ausnahme.md)-Prinzip).
- **Ventil `tracked.exempt-targets`** (Glob über den aufgelösten Zielpfad,
  referenz-weit wie [ADR-0025](0025-codepaths-ignore-refs.md)): absichtlich
  untrackte Ziele (lokal generierte Artefakte) bleiben deklariert erlaubt.
- **Fail-closed:** aktives `tracked` ohne lesbares `.git` ⇒ Exit 2 — kein
  stilles Grün (wie `vcs`/`commits`).
- **Verteilung:** `--print-mk` trägt ein `doc-tracked`-Target
  ([`DC-FA-CLI-010`](../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
  9→10); `--print-config`/`--suggest-config` führen `tracked` als opt-in-Modul.
- **Kein neuer Gate-Bindepunkt in d-check selbst:** das eigene Repo verlinkt
  keine untrackten Ziele (Beleg-Lauf beim Umsetzungs-Slice); ein dauerhaftes
  Make-Target entsteht erst bei Bedarf — Konsumenten binden `doc-tracked`.

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **Opt-in Modul `tracked` über den VCS-Port (gewählt)** | fängt die Drift am Entstehungsort; Index = eine Wahrheit, kein Regel-Nachbau; Port vorhanden (dritte Nutzung); Default bleibt hermetisch | Eingabe-Scope erweitert (nur opt-in); ein weiteres Modul |
| `links` git-aware machen | kein neues Modul | bricht die Hermetik des Default-Moduls — jedes Repo ohne `.git` (Doku-Ordner, Tarball) bekäme Exit 2 oder stilles Teilverhalten |
| `.gitignore` parsen (hermetisch) | kein `.git` nötig | zweiter Regel-Interpreter (Syntax-Nachbau, Nested-ignores, Negationen) — genau die Fehler-Klasse, die der Index-Blick vermeidet; untracked-aber-nicht-ignoriert bliebe unsichtbar |
| Nichts tun (CI als Netz) | kein Aufwand | Feedback erst nach Push im nächsten Checkout; die Falle bleibt Einzelfall-Wissen ([`MR-017`](../../../harness/conventions.md#mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen)/[`MR-019`](../../../harness/conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017) mussten sie je einzeln lernen) |
| Auch gescannte Dateien selbst prüfen | deckt die Gegenrichtung | flaggt jede WIP-Datei vor dem ersten `git add` — Inner-Loop-feindlich; bewusst Out-of-Scope (Re-Eval) |

## Konsequenzen

- **Produkt-Code ändert sich** (Port-Erweiterung + Modul) → Minor-Release
  nötig (wie [ADR-0024](0024-vcs-immutable-gate.md)/[ADR-0027](0027-commits-traceability-modul.md)).
- Die `.d-check.yml`-Basis der fokussierten Gates wächst **nicht** (`tracked`
  ist opt-in, nicht in `modules:`); die aus `ValidModules` **abgeleiteten**
  Fokus-`--disable`-Listen nehmen das neue Modul automatisch auf — die
  handgepflegte `FOCUS_DISABLE`-Variable im eigenen Makefile ist bei der
  Umsetzung zu prüfen (Kommentar-Regel dort: neues Default-Modul ⇒ nachziehen;
  `tracked` ist keins).
- Config-Surface-Currency: `--print-config`, `--suggest-config`
  (opt-in-Hinweis-Kommentar), Benutzerhandbuch, `--print-mk`.
- Die [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Modulliste
  des `gate-consistency`-Selbsttests ist unberührt (opt-in, nicht in
  `modules:`) — wie bei `vcs`/`commits`/`planning`.

## Fitness Function

- Akzeptanztests entlang der
  [`DC-FA-TRK-001`](../../../spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)-Kriterien
  (Happy/Index-Wahrheit/Modul-aus/Negative/Kein-Doppelbefund/Ventil/fail-closed)
  gegen git-Fixture-Repos in `make test`.
- Ohne aktives `tracked` byte-identischer Befundsatz (opt-in-Selbsttest).
- Beleg-Lauf gegen das eigene Repo (`--enable tracked`): grün; adversariale
  Probe mit temporär untracktem, verlinktem Ziel: rot mit `target-untracked`.
- `make gates` unverändert grün (kein Default-Verhalten berührt).

## Re-Evaluierungs-Trigger

- Praxis-Evidenz für untrackte **Inline-Code**-Ziele → Ausdehnung auf
  [`DC-FA-CODE-001`](../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)-Ziele
  (Folge-CR).
- Bedarf, die **Gegenrichtung** (gescannte untrackte Dateien) zu prüfen →
  eigener Schalter mit WIP-tauglicher Semantik (z. B. nur gitignorierte, nicht
  bloß ungestagte Dateien) — heute bewusst Out-of-Scope.
- Submodule/verschachtelte Arbeitsbäume als reale Konsumenten-Fälle →
  Port-Semantik erweitern.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-07-03 | Entwurf (slice-059, welle-48): opt-in Modul `tracked` prüft auflösbare, existierende Link-/Bild-Ziele gegen den git-Index (`target-untracked`); dritte VCS-Port-Nutzung (ohne Range), Index statt `.gitignore`-Interpretation, kein Doppelbefund, Ventil `tracked.exempt-targets`, fail-closed ohne `.git`. Anlass: Auftraggeber-Frage „Was passiert, wenn ein Dokument ein gitignoriertes Dokument referenziert?" + Fixture-Demo (Erzeuger grün, frischer Klon rot). Lastenheft 0.37.0 ([`DC-FA-TRK-001`](../../../spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in) + [`DC-FA-CLI-010`](../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) 9→10). Status Proposed. |
