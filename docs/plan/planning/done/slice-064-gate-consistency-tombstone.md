# Slice slice-064: `gate-consistency.sh` Voll-Tombstone — Prüfung 3 in Go-Test

**Status:** done (welle-53-gate-consistency-tombstone). Lifecycle abgeschlossen
(`next`→`in-progress`→`done`, Roadmap-Flip §Aktuelle Welle,
[`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise));
[ADR-0032](../../adr/0032-gate-consistency-tombstone.md) auf **Accepted**
(ADR-Annotation bei Closure). Ergebnis + Belege in §7.

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
  und prüft, dass `modules` **alle netzlosen Default-Doku-Module** führt und
  `external`/`vcs` ∉ `modules`. **Umfang an die Netzlos-Messmethode gekoppelt**
  (alle Default-`modules` außer `external`/`vcs`), **nicht** die alte 5-Modul-
  Skript-Teilmenge fest verdrahtet — sonst bliebe er grün, wenn
  `spans`/`hostpaths`/`versions` aus der Liste fielen (R1-F-6). Fail-closed
  (fehlende/undekodierbare Datei ⇒ Test rot).
- **`tools/gate-consistency.sh` entfernen (`git rm`);** das Makefile-Target
  `gate-consistency` verliert die `@bash …`-Zeile und fährt nur
  `--enable targets $(FOCUS_DISABLE)`. Target-Name und Dogfood-Selbstbezug bleiben.
- **`codepaths.ignore-refs`-Tombstone (R1-F-1, PFLICHT).** ~30 Inline-Code-Verweise
  `tools/gate-consistency.sh` (CHANGELOG, Spec, ADRs, Roadmap, `done/`-Slices)
  werden nach dem `git rm` sonst `codepath-missing` (Exit 1, `make gates` rot).
  `tools/gate-consistency.sh` kommt in `.d-check.yml` `codepaths.ignore-refs` (wie
  die fünf Vorgänger-Skripte, [ADR-0025](../../adr/0025-codepaths-ignore-refs.md)).
- **Markdown-Links editieren (R1-F-2).** Die zwei **Links** aufs Skript —
  `harness/README.md` §Sensors und **`done/slice-013-codepaths-modul.md`** (ein
  `done/`-Inhalts-Edit) — fängt `ignore-refs` **nicht** (`links`-Achse ist per
  [ADR-0025](../../adr/0025-codepaths-ignore-refs.md) bewusste Rest-Falle); sie
  werden auf Code-Span umgestellt bzw. entfernt (sonst `target-missing`).
- **Doku-/Bindungs-Nachzug (R1-F-3/F-5):** in `harness/README.md` §Sensors und
  `AGENTS.md` §4 nicht nur die Prosa kürzen, sondern die **Netzlos-Bindungs-
  spalte umschreiben** — `gate-consistency` bindet nur noch das Modul `targets`,
  `make test` bekommt die Netzlos-Integrität dazu (sonst behauptet die Sensors-
  Zeile eine Bindung, die dort nicht mehr stattfindet). Ferner: die Skript-Prosa in
  `.d-check.yml` / `config_template.go`, der **live Go-Kommentar**
  `internal/hexagon/core/rules/targets.go:20` (Skript-Paritäts-Verweis) und die
  Roadmap-Notiz („`gate-consistency` bewusst nicht d-check-fähig" — bereits
  korrigiert).
- **Vorbedingung (R1-F-4, erledigt):** [ADR-0031](../../adr/0031-targets-deklarations-konsistenz-modul.md)
  ist auf **Accepted** gesetzt (slice-063-Closure-Backfill) — erst damit ist die
  „Revidiert"-Relation von [ADR-0032](../../adr/0032-gate-consistency-tombstone.md)
  gegen eine Accepted-Vor-ADR sauber.

## 3. Definition of Done

- [ ] Go-Test für die Netzlos-Modulliste (Happy + fail-closed: ein Modul fehlt
  bzw. `external`/`vcs` gesetzt ⇒ rot; Umfang an die Netzlos-Messmethode
  gekoppelt, R1-F-6); `tools/gate-consistency.sh` per `git rm` entfernt;
  Makefile-Target reduziert.
- [ ] **Bruch-Achsen geschlossen (R1-F-1/F-2/F-5):** `codepaths.ignore-refs` +=
  `tools/gate-consistency.sh`; die zwei Markdown-**Links** (`harness/README.md`
  §Sensors, `done/slice-013`) editiert; der Go-Kommentar `targets.go:20`
  nachgezogen. **Kontrolle:** `make doc-check` findet **null** `codepath-missing`
  / `target-missing` — nicht „`grep` findet nur Historie" (falsche Annahme, R1-F-5:
  `done/`/`CHANGELOG`/ADR sind gate-gescannt, `targets.go` ist Produkt-Code).
- [ ] **Doku-/Bindungs-Nachzug (R1-F-3):** die Netzlos-**Bindungsspalte** in
  `harness/README.md` §Sensors + `AGENTS.md` §4 umgeschrieben (`gate-consistency`
  → nur noch Modul `targets`; `make test` → + Netzlos-Integrität); `.d-check.yml` /
  `config_template.go`-Prosa; [ADR-0032](../../adr/0032-gate-consistency-tombstone.md)
  auf Accepted + ADR-Index. (Der R1-F-4-Backfill — die Vor-ADR auf Accepted — ist
  bereits erledigt.)
- [ ] `make gates` / `make ci` grün; **ein unabhängiger Impl-Review**;
  Closure-Move + Body + **Lerneintrag** (Modul 5). **Kein Release.**

## 4. Trigger

Freigabe durch den Auftraggeber — dieser Slice ist der bewusst als **Folge-Slice**
zurückgestellte Voll-Tombstone aus der slice-063-Closure-Diskussion (Nutzer-Frage
2026-07-05: „Können wir `tools/gate-consistency.sh` schon ersetzen?").

## 5. Review-Nachtrag (Plan-/Doc-first-Review R1)

Unabhängiges Plan-Review (1 HIGH / 3 MEDIUM / 1 LOW / 1 INFO, Verdikt
**NACHBESSERN**) — alle Befunde in der ADR und §2/§3 eingearbeitet, **bevor** die
Umsetzung beginnt:

- **F-1 (HIGH, `codepaths.ignore-refs` fehlte):** der Tombstone-Eintrag ist jetzt
  §2/§3-DoD-Pflicht — ohne ihn ~30 `codepath-missing` + rotes Gate.
- **F-2 (MED, Markdown-Links):** die zwei Link-Referenzen (`harness/README.md`
  §Sensors, `done/slice-013`) sind als eigener Edit benannt (`links`-Achse, kein
  `ignore-refs`-Ventil, [ADR-0025](../../adr/0025-codepaths-ignore-refs.md)).
- **F-3 (MED, Netzlos-Bindung):** die Sensors-/§4-**Bindungsspalte** wird
  umgeschrieben (nicht nur Prosa) — sonst behauptet das Gate eine Netzlos-
  Integritäts-Bindung, die es nicht mehr erfüllt.
- **F-4 (MED, Vor-ADR war `Proposed`):** slice-063-Closure-Backfill — die Vor-ADR
  (Modul `targets`) auf **Accepted** gesetzt (Datei + Index + Geschichte); erst
  damit ist die „Revidiert"-Relation dieses Slice gegen eine Accepted-Vor-ADR
  sauber.
- **F-5 (LOW, DoD-`grep`-Annahme):** „findet nur noch Historie" war falsch
  (`done/`/`CHANGELOG`/ADR sind gate-gescannt, `targets.go:20` ist Produkt-Code)
  → DoD auf „`make doc-check` null Befunde" umgestellt.
- **F-6 (INFO, Test-Umfang):** die Assertion koppelt an die Netzlos-Messmethode
  (alle Default-`modules` außer `external`/`vcs`) statt an die 5-Modul-Skript-
  Teilmenge (§2 und die ADR-Konsequenz).

## 6. Review-Nachtrag (Impl-R1)

Unabhängiger Impl-Review (19 Tool-Uses, alle Artefakte real gelesen) — **Verdikt
NACHBESSERN → behoben** (Kern korrekt/vollständig, ein Bindungs-MEDIUM; Report
[`2026-07-06-slice-064-gate-consistency-tombstone`](../../../reviews/2026-07-06-slice-064-gate-consistency-tombstone.md)):

- **MEDIUM (behoben):** in `harness/README.md` §Sensors war nur die Vertrags-,
  nicht die **Bindungs**-Spalte migriert — die `make gate-consistency`-Zeile band
  noch die Netzlos-Modullisten-Integrität, obwohl der Vertrag sie als „voll
  abgelöst → `make test`" ausweist (Row-interner Widerspruch), und `make test`
  trug sie in der Bindung nicht. Genau die von
  [ADR-0032](../../adr/0032-gate-consistency-tombstone.md)/DoD-F-3 verlangte „nicht
  nur Prosa"-Migration. **Behoben:** `gate-consistency`-Bindung → das
  `targets`-Modul-Requirement ([ADR-0031](../../adr/0031-targets-deklarations-konsistenz-modul.md),
  parallel zu `planning-check`), `make test`-Bindung um die Netzlos-Integrität
  ergänzt.
- **Verifiziert sauber:** Paritäts-Äquivalenz (Go-Test strikt stärker als das
  Skript — fängt auch den Abfall von `spans`/`hostpaths`/`versions`, R1-F-6),
  Fail-closed (Read-/Decode-Fehler → rot, Relativpfad korrekt), Guard-Isolation
  (synthetische Listen treffen nur die Assertion, drei Mutationen sterben),
  Tombstone-Vollständigkeit (0 verbliebene Skript-Links, kein dangling
  Produkt-Code-Verweis), Makefile-Rückbau (`targets`-Dogfood bleibt grün).

## 7. Closure-Notiz (nach done)

**Umgesetzt:** `tools/gate-consistency.sh` ist **vollständig entfernt** (`git rm`)
— das **letzte** `tools/*.sh`; der `tools/*.sh`-Audit ist restlos abgeschlossen
(adr-immutable→`vcs`, completeness→Flag, trace→`commits`, planning→`planning`,
arch→a-check, gate-consistency→`targets`+Go-Test). Die Prüfung-3-Restmenge (die
Netzlos-Modulliste der `.d-check.yml`, Vertrag
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit))
ist ein getippter `configyaml.Decode`-Go-Test
(`internal/adapter/driven/configyaml/gate_consistency_test.go`): Live-Prüfung
(fail-closed bei fehlender/undekodierbarer Datei) + Guard-Mutationstest, an die
Netzlos-Messmethode gekoppelt (alle 8 netzlosen Doku-Module präsent,
`external`/`vcs` abwesend) statt der alten 5-Modul-Skript-Teilmenge (R1-F-6).
`make gate-consistency` fährt nur noch `--enable targets` (Modul-Dogfood).

**Belege:** `make gates` + `make ci` grün (gate-consistency = nur `targets`,
`planning-check` grün, neuer `configyaml`-Test grün, image-test byte-identisch).
**Kein Release** (interne Gate-Mechanik, Image byte-identisch — GHCR bleibt
v0.39.0), **kein Lastenheft-CR** (keine neue Anforderung). Impl-Review R1 ACCEPT
nach Nachbesserung (§6). Commit-Kette: feat+lifecycle-in `ea857a6` · closure-move
· closure-body (kein digest-backfill — kein Release).

**Lehre:** (i) beim Bindungs-Nachzug in Sensor-Tabellen ist die **Bindungs-Spalte**
(nicht nur der Vertrags-Text) zu migrieren — sonst überlebt die Harness-Lüge, die
der Slice tilgt, und **kein** Gate fängt es (`targets` prüft nur `make X`↔Regel);
(ii) `git rm` eines gate-gescannten Skripts braucht `codepaths.ignore-refs` **vor**
dem Entfernen (Inline-Verweise) UND separate Link-Edits (die `links`-Achse fängt
`ignore-refs` nicht); (iii) awk-`mv` aus `mktemp` überträgt 600-Rechte → im
Container (nonroot, read-only) „permission denied", `chmod 644` nachziehen.
