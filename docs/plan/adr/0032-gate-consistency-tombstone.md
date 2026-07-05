# ADR-0032 — `gate-consistency.sh` Voll-Tombstone: die `DC-QA-03`-Restprüfung wandert in einen Go-Test

**Status:** Proposed
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
- **Trade-off:** die Restprüfung ist keine eigene `make`-Gate-Zeile mehr, sondern
  ein Unit-Test — minimal weniger „gate-sichtbar", aber verlässlicher und ohne
  Shell.
- **Doku-Nachzug:** die `make gate-consistency`-Beschreibung (AGENTS §4,
  `harness/README.md` §Sensors) und die Skript-Verweise in `.d-check.yml` /
  `config_template.go` verlieren ihren „Rest-Skript"-Teil; die Roadmap-Notiz
  „`gate-consistency` bewusst nicht d-check-fähig" (seit der `targets`-Einführung
  ohnehin falsch) entfällt.
- **Reversibel:** rein interne Gate-Mechanik, kein nutzersichtbares Verhalten,
  **kein Release** (wie die `arch-check`- und `completeness`-Rückbauten).
