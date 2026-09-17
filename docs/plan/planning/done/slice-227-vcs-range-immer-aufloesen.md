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

- [x] **(1)** `--enable vcs --range <unauflösbar>` ohne (oder mit leerem)
      `vcs:`-Block bricht mit Exit 2 ab, in derselben Fehlermeldungsform wie
      mit konfigurierten `vcs.paths` — Regressionstest gegen die
      Issue-#4-Reproduktion.
- [x] **(2)** Dieselbe Zusage für `commits`: `--enable commits --range
      <unauflösbar>` ohne (oder mit leerem) `commits:`-Block bricht mit
      Exit 2 ab.
- [x] **(3)** Eine **auflösbare** Range mit leerer Klassen-Config bleibt
      weiterhin befundfrei (Exit 0) — Regressionstest für die unveränderte
      Rückwärtskompatibilität des dokumentierten Opt-in-Vertrags.
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

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
  aber nicht gemessen. — **Ausgang:** entfallen — der Aufruf war bereits
  vorher unconditional fällig, sobald die Klassen-Config **nicht** leer
  war (der Regelfall); nur der leere Fall ist neu, und der ist per
  Definition der seltenere. Kein eigenes Performance-Risiko.
- **Der ausgeschlossene Punkt (Ausschlusspunkt 1) bleibt eine benannte
  Lücke**, falls der Auftraggeber ihn doch für nötig hält — dann ist er ein
  Folge-Slice, keine Nacharbeit an diesem. — **Ausgang:** entfallen — bei
  näherer Betrachtung kein Risiko dieses Slice, sondern eine bereits in §1
  getroffene, bewusste Abgrenzungs-Entscheidung; sie bleibt dort (in dem
  dann archivierten Slice-Plan) auffindbar, falls sie später aufgegriffen
  wird. Kein Steering-Loop-Muster, das eine eigene Register-Zeile trüge.

## 7. Closure-Notiz

**Geliefert.** `CheckVCS` und `CheckCommits`
(`internal/hexagon/core/rules/vcs.go`/`commits.go`) lösen eine angegebene
Range jetzt immer über den VCS-Port auf, **bevor** die Klassen-Config
(`vcs.paths`/`commits.id-patterns`) geprüft wird — eine syntaktisch
gültige, aber unauflösbare Range bricht seither mit Exit 2 ab, unabhängig
davon, ob es etwas zu prüfen gibt. Eine auflösbare Range mit leerer
Klassen-Config bleibt weiterhin befundfrei (Exit 0) — die dokumentierte
Opt-in-Trägheit ist unverändert, nur nicht mehr an eine ungeprüfte Range
gekoppelt.

**Was funktionierte.** Die Root-Cause-Analyse traf im ersten Anlauf zu
(derselbe Frühausstieg in zwei Funktionen) und wurde durch bewusstes
Brechen (Modul 11) unabhängig zweimal bestätigt — einmal durch mich selbst,
einmal durch den Review. Die End-to-End-Reproduktion gegen das gebaute
Image lieferte in beiden Läufen identische Ergebnisse.

**Was anders lief.** Der unabhängige Review (R1) fand zwei Lücken, die die
eigene Prüfung vor dem Review nicht sah:

- **R1-F-1 (HIGH):** Zwei neue Testkommentare trugen „(GitHub Issue #4
  Punkt 1)" als Herkunfts-Prosa statt einer der fünf nach `AGENTS.md` §3.7
  zulässigen Kommentarklassen. Behoben durch Entfernen der Klammer-Referenz.
  **Lerneintrag:** Eine externe Vorgangs-Kennung (Issue-Nummer) ist
  dieselbe verbotene Klasse wie eine Slice- oder Befund-Nummer im
  Kommentar — sie fällt aber leichter durch die eigene Aufmerksamkeit,
  weil sie wie ein legitimer, extern nachprüfbarer Beleg aussieht.
- **R1-F-2 (MEDIUM):** Die Reihenfolge-Änderung in `CheckCommits` schließt
  die Klasse vollständig (sie ist nicht pfad-, sondern funktionsspezifisch)
  — aber die Root-Cause-Analyse, das DoD und die ursprünglichen zwei Tests
  benannten nur `--range`. `--enable commits --staged` mit leerer
  `commits:`-Config teilte denselben Fehler (still Exit 0 vor dem Fix) und
  wechselt durch denselben Fix korrekt auf fail-closed Exit 2 — unbenannt
  und ungetestet, bis der Review es empirisch fand. Geschlossen durch einen
  dritten Regressionstest.

**Register.** Zwei neue Einträge:
[`BEO-ALL/kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen`](../observations/BEO-ALL/kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen/observation.md)
(R1-F-1) und
[`BEO-ALL/fix-aendert-unbenannten-zweiten-pfad-mit`](../observations/BEO-ALL/fix-aendert-unbenannten-zweiten-pfad-mit/observation.md)
(R1-F-2), beide 1×.

**Risiken aus §6:** beide aufgelöst — der Performance-Punkt entfällt
(der teure Fall war ohnehin schon der Regelfall), der Abgrenzungs-Punkt
entfällt (er war nie ein Risiko dieses Slice, sondern eine bereits in §1
getroffene Entscheidung).

**Drei Paarungen** (Repo ohne Wellen-Betrieb, hier geprüft): Anker — kein
`liegt in`-Feld in dieser Notiz, keine Paarung fällig. Folge-Slice — keiner
genannt. Register — beide zitierten `BEO-ALL/…`-Pfade existieren mit
nicht-leerem `evidence/` (oben angelegt).

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
