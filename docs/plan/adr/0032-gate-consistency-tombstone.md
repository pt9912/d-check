# ADR-0032 — `gate-consistency.sh` Voll-Tombstone: die `DC-QA-03`-Restprüfung wandert in einen Go-Test

**Status:** Accepted
**Datum:** 2026-07-05
**Autor:** pt9912
**Bezug:** [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(die von der Restprüfung bewachte Netzlos-Gate-Integrität);
[ADR-0028](0028-planning-lifecycle-modul.md) (Präzedenz: der `planning`-Tombstone
entfernte sein Skript **vollständig**); Abschluss des `tools/*.sh`-Audits.
**Revidiert:** die `Scope = Kern`-Konsequenz von
[ADR-0031](0031-targets-deklarations-konsistenz-modul.md) — dort wurde
entschieden, dass `tools/gate-consistency.sh` für die Prüfung-3-Restmenge
**nicht vollständig** entfernt wird; diese ADR hebt genau das auf (das Modul und
die Identitäts-Ausweitung von [ADR-0031](0031-targets-deklarations-konsistenz-modul.md)
bleiben gültig).

## Kontext

[ADR-0031](0031-targets-deklarations-konsistenz-modul.md) mechanisierte den
cross-repo-driftenden **Kern** von `tools/gate-consistency.sh` (Doku ↔ Makefile,
beide Richtungen) als das opt-in Modul `targets` und **behielt bewusst** die
dritte Prüfung im reduzierten Skript: die repo-spezifische
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Selbstkonsistenz
der `.d-check.yml`-Modulliste (das netzlose `doc-check`-Gate führt alle netzlosen
Doku-Module und **nicht** `external`). Begründung damals: kein Kopie-Drift-Kern,
also nicht Modul-würdig.

Zwei Punkte machen den **vollen** Tombstone jetzt fällig:

1. **Das Modul ist bewährt.** `targets` ist seit **v0.38.0** released, dogfoodet
   `make gate-consistency` und ist paritäts-verifiziert (Mutations-Beleg gegen
   das Skript). Der Kern-Grund für das Skript ist weg.
2. **Es ist das letzte `tools/*.sh`.** Der `tools/*.sh`-Audit hat alle anderen
   Gate-Skripte mechanisiert (`adr-immutable`→`vcs`, `trace`→`commits`,
   `planning`→`planning`, `completeness`→Produkt-Flag). Ein ~40-Zeilen-Shell-Rest
   für eine reine Datei-Assertion ist die letzte vermeidbare **Harness-Lügen-
   Oberfläche** (Skript ≠ getesteter Go-Code).

Die Restprüfung ist **keine** cross-repo-Drift (sie prüft d-checks *eigene*
`.d-check.yml`), taugt also nicht als verteilbares Modul — aber sie taugt als
**Go-Test**: sie ist genau eine Assertion über das Decode-Ergebnis der
Live-`.d-check.yml`.

## Entscheidung

`tools/gate-consistency.sh` wird **vollständig entfernt**. Die
Modullisten-Restprüfung wandert in einen **Go-Test** im `configyaml`-Paket: er
liest die Repo-`.d-check.yml` über einen Relativ-Pfad (etabliertes Muster aus
`internal/hexagon/core/app/diagnose_test.go`), dekodiert sie mit
`configyaml.Decode` und prüft, dass die `modules`-Liste alle netzlosen
Doku-Module (`links`/`anchors`/`ids`/`matrix`/`codepaths`) führt und **nicht**
`external`. `make gate-consistency` fährt danach **nur** den `targets`-Modul-
Dogfood; das Target bleibt (Name und Selbstbezug unverändert), verliert nur die
`bash`-Zeile.

## Konsequenzen

- **Null `tools/*.sh`.** Der `tools/*.sh`-Audit ist damit restlos abgeschlossen —
  jedes verbleibende Gate ist getesteter Go-Code, ein Produkt-Flag oder ein
  Modul-Dogfood.
- **Die Restprüfung wird im `make test`-Lauf geprüft** (typisierter Decode statt
  `grep`), nicht mehr als Shell-Hook in `make gate-consistency`. Robuster gegen
  YAML-Format-Änderungen (der Skript-`grep` nahm Flow-Style an; ein Wechsel zur
  Listenform hätte ihn still falsch-grün gemacht — der Go-Test dekodiert echt).
  Der Assertions-**Umfang** wird dabei an die aktuelle Netzlos-Messmethode
  gekoppelt (alle netzlosen Default-`modules` außer `external`/`vcs`) statt die
  alte 5-Modul-Skript-Teilmenge fest zu verdrahten — sonst bliebe er grün, wenn
  `spans`/`hostpaths`/`versions` aus der `modules`-Liste fielen (R1-F-6).
- **Trade-off:** die Restprüfung ist keine eigene `make`-Gate-Zeile mehr, sondern
  ein Unit-Test — minimal weniger „gate-sichtbar", aber verlässlicher und ohne
  Shell.
- **`codepaths.ignore-refs`-Tombstone (Pflicht, [ADR-0025](0025-codepaths-ignore-refs.md)).**
  ~30 Inline-Code-Verweise `tools/gate-consistency.sh` (CHANGELOG, Spec, ADRs,
  Roadmap, `done/`-Slices) werden nach dem `git rm` sonst `codepath-missing` —
  `tools/gate-consistency.sh` kommt in `.d-check.yml` `codepaths.ignore-refs`
  (wie die fünf Vorgänger-Skripte). Die **zwei Markdown-Links** aufs Skript
  (`harness/README.md` §Sensors, eine `done/`-Slice) fängt das Register **nicht**
  (die `links`-Achse ist per [ADR-0025](0025-codepaths-ignore-refs.md) bewusste
  Rest-Falle) → sie werden **editiert** (Code-Span bzw. Entfernung, inkl. eines
  `done/`-Inhalts-Edits); ebenso der live Go-Paritäts-Kommentar in
  `internal/hexagon/core/rules/targets.go`.
- **Netzlos-Bindung wandert mit.** `make gate-consistency` bindet danach nur
  noch das Modul `targets`; die Netzlos-Modullisten-Integrität prüft der Go-Test
  unter `make test`. Die **Bindungsspalte** in `harness/README.md` §Sensors und
  die `AGENTS.md` §4-Beschreibung werden entsprechend umgeschrieben (nicht nur die
  Prosa) — sonst behauptete die Sensors-Zeile eine Netzlos-Integritäts-Durchsetzung
  durch `gate-consistency`, die dort nicht mehr stattfindet. Die Roadmap-Notiz
  „`gate-consistency` bewusst nicht d-check-fähig" (seit der `targets`-Einführung
  ohnehin falsch) entfällt.
- **Reversibel:** rein interne Gate-Mechanik, kein nutzersichtbares Verhalten,
  **kein Release** (wie die `arch-check`- und `completeness`-Rückbauten).

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-07-06 | Angenommen mit der slice-064-Closure: `tools/gate-consistency.sh` per `git rm` entfernt (**letztes** `tools/*.sh`, Audit restlos abgeschlossen); die Netzlos-Modullisten-Restprüfung als getippter `configyaml.Decode`-Go-Test (Live + Guards + fail-closed, an die Netzlos-Messmethode gekoppelt — alle 8 netzlosen Doku-Module präsent, `external`/`vcs` abwesend, R1-F-6); `make gate-consistency` fährt nur noch `--enable targets`; `codepaths.ignore-refs` += Skript-Pfad (R1-F-1) + zwei Markdown-Link-Edits (R1-F-2); Bindungsspalten (`harness/README.md` §Sensors, `AGENTS.md` §4) auf `make test` bzw. [ADR-0031](0031-targets-deklarations-konsistenz-modul.md) migriert (R1-F-3, im Impl-Review als MEDIUM nachgezogen). Impl-Review R1 NACHBESSERN → behoben, dann ACCEPT. `make ci` grün, kein Lastenheft-CR, kein Release (Image byte-identisch). Status Accepted. |
