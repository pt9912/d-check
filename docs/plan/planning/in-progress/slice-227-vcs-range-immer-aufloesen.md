# Slice slice-227: Eine angegebene Range wird immer aufgelöst, auch ohne Klassen-Config

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
[`DC-FA-COMMITS-001`](../../../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in).

**Berührte Spec-Stellen:** —
(kein Spec-Nachtrag geplant: die Zusage „eine angegebene Range wird
aufgelöst, sonst Exit 2" steht bereits im VCS-Port-Vertrag
([`internal/hexagon/port/driven/vcs.go`](../../../../internal/hexagon/port/driven/vcs.go)) und in `AllPaths`/`CommitMessages`
selbst — dieser Slice schließt eine Lücke in der **Aufrufreihenfolge**
zweier Regelmodule, keine neue Zusage).

**Verantwortlich:** pt9912 (Implementer-Rolle).

**Autor:** pt9912. **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar.

**Anlass ist [GitHub Issue #4](https://github.com/pt9912/d-check/issues/4)
Punkt 1** (Adopter-Meldung, 2026-09-17): Ist `vcs` per `--enable vcs`
eingeschaltet, fehlt aber der `vcs:`-Block in der Konfiguration (oder die
ganze `.d-check.yml`), endet `--range <base>..<head>` mit `0 Befund(e)`,
Exit 0 — **auch wenn die Range gar nicht auflösbar ist**
(`deadbeef..cafebabe`). Reproduziert und bestätigt gegen den aktuellen
Quellstand.

**Root Cause, gelesen im Quelltext.** `CheckVCS`
([`internal/hexagon/core/rules/vcs.go:30-32`](../../../../internal/hexagon/core/rules/vcs.go))
und `CheckCommits`
([`internal/hexagon/core/rules/commits.go:23-26`](../../../../internal/hexagon/core/rules/commits.go))
tragen **denselben** Frühausstieg:

```go
if len(cfg.Paths) == 0 || vcs == nil {   // CheckVCS
        return nil, nil // inert
}
```

Fehlt die Klassen-Config (`vcs.paths` bzw. `commits.id-patterns`), kehrt die
Funktion zurück, **bevor** sie `vcs.AllPaths`/`vcs.CommitMessages` aufruft —
genau der Aufruf, der die Range fail-closed auflöst
([`internal/hexagon/port/driven/vcs.go:22-42`](../../../../internal/hexagon/port/driven/vcs.go)).
Der CLI-seitige Fail-Closed-Pfad
([`internal/adapter/driving/cli/cli.go:697`](../../../../internal/adapter/driving/cli/cli.go),
`resolveVCS`/`vcsRefs`) fängt nur eine **fehlende** `--range`/`--staged`-Angabe
ab — eine **syntaktisch** gültige, aber inhaltlich unauflösbare Range
(`deadbeef..cafebabe`) erreicht `resolveVCS` unbeanstandet und wird erst in
`CheckVCS`/`CheckCommits` selbst geprüft, wo der Frühausstieg sie überspringt.
**Betroffen sind beide Module** — `commits` trägt denselben Fehler, das
Issue nennt nur `vcs`.

**Ziel:** Ist `vcs` bzw. `commits` aktiv und eine Range angegeben
(`--range`/`--staged`), wird sie **immer** über den VCS-Port aufgelöst — auch
wenn die Klassen-Config leer ist. Eine nicht auflösbare Range bricht dann
mit Exit 2 ab, unabhängig davon, ob es eine geschützte Datei bzw. ein
ID-Muster gibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ob „aktiv ohne jede Klassen-Config" selbst ein Konfigurationsfehler
  wird** (Issue #4s erste Erwartungs-Hälfte: „`--enable vcs` ohne
  `vcs:`-Block ist ein Konfigurationsfehler mit Exit 2"). Beide Module
  tragen den dokumentierten, bewussten Vertrag „Opt-in (leere Paths ⇒
  inert)" ([`vcs.go:15`](../../../../internal/hexagon/core/rules/vcs.go),
  [`commits.go`](../../../../internal/hexagon/core/rules/commits.go)) — ihn
  umzukehren ist eine **Verhaltensänderung** der Opt-in-Semantik (nicht nur
  ein Bugfix) mit möglicher Breaking-Wirkung für Konsumenten, die das Modul
  bewusst ohne Klassen-Config in einer geteilten Pipeline aktivieren, bevor
  sie Paths/Muster nachziehen. Mit der Range-Auflösung aus diesem Slice ist
  der harte Teil der Meldung („eine falsche Range fällt nie auf") ohnehin
  gelöst; ob die zweite Hälfte zusätzlich gewollt ist, ist ein eigener
  Entscheid — *Bestand bleibt bewusst stehen*, bis er getroffen ist.
- **Kein pauschales Audit aller übrigen opt-in-Module auf dasselbe Muster.**
  `vcs`/`commits` sind die beiden Post-Pässe, die zusätzlich eine
  CLI-Range-Angabe tragen (`--range`/`--staged`) — bei den übrigen
  Datei-scannenden opt-in-Modulen gibt es keine vergleichbare zweite,
  unabhängig gültige Eingabe, die ein Frühausstieg verschlucken könnte. Ein
  Audit auf strukturell ähnliche Fälle wäre *ein anderer Vorgang*.
- **Keine Änderung an der git-Objektauflösung selbst** (`packAliasFS`,
  go-git-Nutzung) — dieser Bug liegt ausschließlich in der
  Aufrufreihenfolge der beiden Regelmodule, nicht im VCS-Port.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** `--enable vcs --range <unauflösbar>` ohne (oder mit leerem)
      `vcs:`-Block bricht mit Exit 2 ab, in derselben Fehlermeldungsform wie
      mit konfigurierten `vcs.paths` — Regressionstest gegen die
      Issue-#4-Reproduktion.
- [ ] **(2)** Dieselbe Zusage für `commits`: `--enable commits --range
      <unauflösbar>` ohne (oder mit leerem) `commits:`-Block bricht mit
      Exit 2 ab.
- [ ] **(3)** Eine **auflösbare** Range mit leerer Klassen-Config bleibt
      weiterhin befundfrei (Exit 0) — Regressionstest für die unveränderte
      Rückwärtskompatibilität des dokumentierten Opt-in-Vertrags.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/core/rules/vcs.go` (`CheckVCS`) | update | Frühausstieg zerlegen: erst `vcs == nil` prüfen und `vcs.AllPaths` aufrufen (fail-closed), danach `len(cfg.Paths) == 0` prüfen und ggf. ohne Befunde zurückkehren |
| `internal/hexagon/core/rules/commits.go` (`CheckCommits`) | update | dieselbe Zerlegung für `vcs.CommitMessages` / `cfg.IDPatterns` |
| `internal/hexagon/core/rules/vcs_test.go` bzw. neue Testdatei | neu/update | Happy: auflösbare Range + leere Config ⇒ 0 Befunde (DoD 3); Negative: unauflösbare Range + leere Config ⇒ Fehler (DoD 1) |
| `internal/hexagon/core/rules/commits_test.go` bzw. neue Testdatei | neu/update | dieselben zwei Fälle für `commits` (DoD 2) |

## 4. Trigger

**Start** (`next` → `in-progress`): Freigabe zur Umsetzung (Issue #4 ist
öffentlich gemeldet, keine gesonderte Freigabe über den CR-Kanal nötig).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls sich beim
  Umsetzen zeigt, dass `vcs` und `commits` nicht symmetrisch reparierbar
  sind (z. B. eine der beiden Port-Methoden hat einen abweichenden
  Fehler-Vertrag) und die beiden Hälften getrennt geplant werden müssen.
- `in-progress` → `open` (blockiert): falls die Team-/Auftraggeber-Linie zu
  Ausschlusspunkt 1 (die Config-Fehler-Hälfte) doch in diesem Slice verlangt
  wird — dann wächst der Umfang über die Abgrenzung hinaus und braucht einen
  neuen Zuschnitt.

## 5. Closure-Trigger

DoD (1) bis (3) sowie `make gates`, Review, Closure-Notiz, Register,
Risiko-Ausgänge, Paarungen — alle abgehakt.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der zusätzliche `AllPaths`/`CommitMessages`-Aufruf bei leerer
  Klassen-Config kostet einen git-Lese-Vorgang, der bisher übersprungen
  wurde.** Für den in [`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance)
  gemessenen Rahmen (übliche Range-Größen) vermutlich vernachlässigbar,
  aber nicht gemessen. — **Ausgang:** \<offen\>
- **Der ausgeschlossene Punkt (Ausschlusspunkt 1) bleibt eine benannte
  Lücke**, falls der Auftraggeber ihn doch für nötig hält — dann ist er ein
  Folge-Slice, keine Nacharbeit an diesem. — **Ausgang:** \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende.

Dieses Repo führt **drei** Prüfungen — die zwei kanonischen und, als
Adaption, den Nachtlauf-Stand ([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `internal/hexagon/core/rules/` (die beiden Regelmodule) —
trägt keine eigene Modus-Deklaration in `harness/conventions.md` und fällt
damit unter den Default `*`. Greenfield, wie der Rest des Produkts. Der
CLI-Layer (`internal/adapter/driving/cli/`) wird nur lesend zitiert
(Root-Cause-Beleg), nicht geändert — keine zweite Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, 45 Verzeichnisse). **Keine** der
Einträge trifft `internal/hexagon/core/rules/` als Sub-Area.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-17 gelesen: `image-scan.yml` grün
(2026-09-17T08:41:27Z). `upstream-drift.yml` **rot** (2026-09-17T05:38:51Z)
— laut eigener Meldung eine **planmäßige** Fremd-Release-Benachrichtigung
([`MR-051`](../../../../harness/conventions.md#mr-051)), keine unerwartete;
betrifft gepinnte Fremd-Bestände, nicht diesen Slice.

**Modus-Begründung:** die einzige berührte Sub-Area ist GF — kein
Begründungsblock nötig.
