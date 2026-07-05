# Slice slice-064: `gate-consistency.sh` Voll-Tombstone — Prüfung 3 in Go-Test

**Status:** next (Backlog, welle-53-gate-consistency-tombstone). Das
Doc-first-Fundament ([ADR-0032](../../adr/0032-gate-consistency-tombstone.md)) ist
geschrieben; die Implementierung startet **frisch** (sobald der Implementer
beginnt ⇒ Move nach `in-progress/`, Roadmap-Flip, [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

**Welle:** welle-53-gate-consistency-tombstone (Kandidat).

**Bezug:** [ADR-0032](../../adr/0032-gate-consistency-tombstone.md)
(Voll-Tombstone-Entscheidung; revidiert die `Scope = Kern`-Konsequenz von
[ADR-0031](../../adr/0031-targets-deklarations-konsistenz-modul.md)) +
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(die von der Restprüfung bewachte Netzlos-Gate-Integrität). **Kein**
Lastenheft-CR (keine neue Anforderung), **kein** Release (interne Gate-Mechanik).

**Autor:** pt9912. **Datum:** 2026-07-05.

---

## 1. Ziel

`tools/gate-consistency.sh` ist nach slice-063 auf die
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Modullisten-Restprüfung
reduziert — das **letzte** `tools/*.sh`. Diese Prüfung ist eine reine Assertion
über das Decode-Ergebnis der Live-`.d-check.yml` und gehört daher in einen
**Go-Test**, nicht in ein Shell-Skript. Ziel: das Skript **vollständig** entfernen
(voller Tombstone wie beim `planning`-Modul); `make gate-consistency` fährt nur
noch den `targets`-Modul-Dogfood.

## 2. Entscheidungen (aus [ADR-0032](../../adr/0032-gate-consistency-tombstone.md))

- **Prüfung 3 → Go-Test** im `configyaml`-Paket: liest die Repo-`.d-check.yml`
  über einen Relativ-Pfad (Muster aus
  `internal/hexagon/core/app/diagnose_test.go`), dekodiert via `configyaml.Decode`
  und prüft `modules` ⊇ {`links`, `anchors`, `ids`, `matrix`, `codepaths`} sowie
  `external` ∉ `modules`. Fail-closed (fehlende/undekodierbare Datei ⇒ Test rot).
- **`tools/gate-consistency.sh` entfernen;** das Makefile-Target `gate-consistency`
  verliert die `@bash …`-Zeile und fährt nur `--enable targets $(FOCUS_DISABLE)`.
  Target-Name und Selbstbezug (AGENTS §4 / §Sensors) bleiben — nur der
  Rest-Skript-Teil der Beschreibung entfällt.
- **Doku-Nachzug:** AGENTS §4, `harness/README.md` §Sensors, die Prosa in
  `.d-check.yml` und `config_template.go` (Skript-Verweise) sowie die
  Roadmap-§Nächste-Wellen-Zeile („`gate-consistency` bewusst nicht d-check-fähig",
  seit slice-063 falsch).

## 3. Definition of Done

- [ ] Go-Test für die Netzlos-Modulliste (Happy + fail-closed: ein Modul fehlt
  bzw. `external` gesetzt ⇒ rot); `tools/gate-consistency.sh` entfernt;
  Makefile-Target reduziert; `grep -r gate-consistency.sh` findet nur noch
  Historie (`done/`, `CHANGELOG`, ADR-Kontext).
- [ ] Doku-Nachzug (AGENTS §4, §Sensors, `.d-check.yml` / `config_template.go`,
  Roadmap-Zeile); [ADR-0032](../../adr/0032-gate-consistency-tombstone.md) auf
  Accepted + ADR-Index.
- [ ] `make gates` / `make ci` grün; **ein unabhängiger Review**; Closure-Move +
  Body + **Lerneintrag** (Modul 5). Kein Release.

## 4. Trigger

Freigabe durch den Auftraggeber — dieser Slice ist der bewusst als **Folge-Slice**
zurückgestellte Voll-Tombstone aus der slice-063-Closure-Diskussion (Nutzer-Frage
2026-07-05: „Können wir `tools/gate-consistency.sh` schon ersetzen?").
