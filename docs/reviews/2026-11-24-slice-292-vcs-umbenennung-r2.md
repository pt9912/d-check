# Review-Report: slice-292 — Kern und git-Adapter, `vcs` folgt der Umbenennung (Code R2) — 2026-11-24

**Review-Art:** Code. Geprüft wird der Fix-Commit gegen die Befunde von R1, Slice-Plan, ADR und Hard Rules,
also Maintainability (`v6.18.1` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten). Gegen die DoD
prüft danach der Verifier.

**Gegenstand:** `slice-292`, ein Commit: `95aef23b` (nur lokal auf `main`). Er berührt drei Dateien,
21 Zeilen neu (`git show --stat 95aef23b`):

- `internal/adapter/driven/git/git_renames_test.go`: `TestRenamesCommitMengeOhneVorBase`
- `internal/hexagon/core/rules/vcs_follow_renames_test.go`: ein Fall in `TestVCSFollowRenamesRandfaelle`
- `docs/user/abdeckung-tests.md`: eine abgeleitete Zeile

**Skill:** `.harness/skills/reviewer.md` @ 1.20.0 (`d1c37c14`).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-24.

**Eingangs-Kontext:**

- Report [R1](2026-11-24-slice-292-vcs-umbenennung-r1.md), besonders M-1 und M-2 samt ihren Probe-Fixtures
  und den Mutationen M19 und M5.
- Plan `slice-292` §2 (Rot ohne die Regel für jeden Test) und §3 (Dateitabelle).
- [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md), Entscheidung 2, 10 und 11.
- Lastenheft
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  und Spezifikation
  §[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  Schritt 2 und 4.
- [`AGENTS.md`](../../AGENTS.md) §3.7.
- Vorherige Findings am Modul: die acht aus R1. Diese Runde ändert an L-1, L-2 und I-1 bis I-4 nichts; sie
  gehen nach der Vereinbarung ins Register.

**Wie gelesen und gefahren wurde.** `git show 95aef23b` vollständig; die Fixture-Hilfen (`graph`, `newRepo`,
`runRenameCases`) und die Tabelle samt Schleife von `TestVCSFollowRenamesRandfaelle`. Zwei Mutationen im
Arbeitsbaum, je eine, mit `make test` gefahren und danach mit `git checkout --` zurückgenommen. Am Ende sind
`git status --short` und `git diff` leer; nichts wurde committet.

---

## Je Prüffrage ein Satz

| # | Prüffrage | Urteil |
|---|---|---|
| 1 | M-1 geschlossen? | Ja: das Fixture fährt die Menge „von HEAD erreichbar, von BASE nicht“ und ist rot unter M19, und nur dort |
| 2 | M-2 geschlossen? | Ja: der Fall ist unter M5 rot, mit der Meldung des Falles selbst, und grün gegen den unveränderten Code |
| 3 | Botschaft gegen Diff | Beide Aussagen sind gemessen; keine Verallgemeinerung |
| 4 | Fixtures stabil | Adapter: Objekte im Speicher ohne Arbeitsbaum und Index; Kern: Fake ohne Dateisystem; kein racily-clean |
| 5 | Abgrenzung und Kommentar | Nur die drei Dateien; der Kommentar trägt Zusage und Abgrenzung |

---

## Findings

Keine. Es gibt kein neues Finding, weder HIGH noch MEDIUM noch LOW noch INFO.

| Nr. | Kategorie | Kurz | Behandlung |
|---|---|---|---|
| — | — | kein Finding | — |

## Zu M-1: `TestRenamesCommitMengeOhneVorBase`

Das Fixture ist das Probe-Fixture aus R1, Zeile für Zeile: `r` mit `a.md` (Inhalt A) und `n.md`; `r2` löscht
`a.md`; `base` legt `a.md` mit anderem Inhalt (B) neu an; `head` benennt `a.md` nach `b.md` um. Es fährt die
Commit-Menge also wirklich: `r` und `r2` liegen vor BASE, sind von HEAD erreichbar und von BASE ebenfalls, und
der Pfad trägt in ihnen ein Ereignis (Löschen), das zu einer Quelle in BASE gehört. Der Test fährt über
`runRenameCases` in drei Kennungs-Salzen; die erwartete Zuordnung ist `{a.md: b.md}`.

**Mutation M19 (selbst gesetzt).** In `git.go` Zeile 853 aus `if w.excluded[h] || w.commits[h] != nil {` den
Teil `w.excluded[h] ||` entfernt, `make test`:

- rot; als einziger Test fällt `TestRenamesCommitMengeOhneVorBase/Pfad_vor_BASE_gelöscht_und_neu_angelegt`
  mit `Salz "": Zuordnung = map[], want map[a.md:b.md]`.
- Die Ursache ist die richtige: die Umbenennung geht verloren, weil die Linie „gelöscht“ aus der Historie vor
  BASE mitwandert. Die Mutation ist zurückgenommen (`git checkout -- internal/adapter/driven/git/git.go`).

Damit ist der Rot-Beleg für die Regel „Commits vor BASE tragen nichts bei“ vorhanden, mit dem Failure-Szenario
aus R1. **M-1 ist geschlossen.**

## Zu M-2: der Fall „beide Schlüssel: Ende im Wurzelverzeichnis …“

Der Fall: Link `d1/x.md` → `x.md`, BASE-Baum `{a.md, d1/x.md}`, HEAD-Baum `{x.md, d1/x.md}`, Zuordnung
`{a.md: x.md}`, erwartet `counts{follow 1, both 1, ignore 0}`. Das ist das Probe-Fixture aus R1 mit denselben
Zählern. `d1/x.md` nimmt nicht teil (liegt in BASE und HEAD), `x.md` ist das Ende einer teilnehmenden
Umbenennung: Endpfad und Dateiname sind dieselbe Zeichenkette, nur die doppelte Marke trennt die beiden Formen.

**Mutation M5 (selbst gesetzt).** In `vcs_link_targets.go` Zeile 116 `const renamedTargetMark =
normalizedTargetMark + normalizedTargetMark` auf `= normalizedTargetMark`, `make test`:

- rot; als einziger Test fällt `TestVCSFollowRenamesRandfaelle/beide_Schlüssel:_Ende_im_Wurzelverzeichnis,_Umbiegen_auf_eine_andere_Datei_gleichen_Namens`
  mit `beide schlüssel: Befunde = 0, want 1 ([])`.
- Der Fall ist nicht aus einem anderen Grund rot: die Meldung nennt den Modus „beide schlüssel“ und die Zahl
  der Befunde, die ausbleibt, weil beide Seiten zu `Marke + x.md` werden und das Umbiegen auf ein anderes
  Dokument durchgeht; die Modi „follow-renames allein“ (1) und „ignore-link-targets allein“ (0) sind davor
  und danach unverändert und tragen die Mutation nicht. Die Mutation ist zurückgenommen
  (`git checkout -- internal/hexagon/core/rules/vcs_link_targets.go`).

**M-2 ist geschlossen.**

## Mutationen

| M | Mutation | Ergebnis | Fällt in |
|---|---|---|---|
| M19 | `w.excluded[h] \|\|` aus der Bedingung der Commit-Menge entfernt | rot (vorher grün) | `TestRenamesCommitMengeOhneVorBase`, sonst nichts |
| M5 | doppelte Marke einfach | rot (vorher grün) | `TestVCSFollowRenamesRandfaelle`, der neue Fall, sonst nichts |

Beide sind zurückgenommen; `git status --short` und `git diff` sind leer.

---

## Negativbefunde

- **Botschaft gegen Diff (Prüffrage 8).** Die Botschaft behauptet drei Dinge: (a) ein Pfad, der vor BASE
  gelöscht und neu angelegt wurde, behält seine Umbenennung, und mit mitwandernden Commits geht sie
  verloren; (b) mit beiden Schlüsseln und einem Ende im Wurzelverzeichnis ist das Umbiegen ein Befund, mit
  einfacher Marke fiele es durch; (c) nur Tests und die abgeleitete Datei, kein Produktcode. (a) und (b)
  sind in dieser Runde genau so gemessen (M19, M5); (c) stimmt mit der Stat-Liste überein. Die Botschaft
  behauptet nicht, dass die Menge oder die Marke damit „vollständig“ gedeckt seien. Geprüft, ohne Befund.
- **Fixtures (`BEO-ALL/racily-clean-git-fixture`).** Der Adapter-Fall legt Commits über `graph`/`commitIn`
  direkt als Objekte an; es gibt weder Arbeitsbaum noch Index, also keinen Stat-Vergleich, der gleich lang
  neu geschriebene Dateien übersehen könnte. Die drei Salze ändern die Kennungen, und der Fall läuft in allen
  dreien grün. Der Kern-Fall nutzt `vcsTreeFake` und einen Fake ohne Dateisystem. Geprüft, ohne Befund.
- **Abgrenzung.** Der Commit berührt genau die drei genannten Dateien; kein Produktcode, keine ADR, kein
  `AGENTS.md`, kein Handbuch. Geprüft, ohne Befund.
- **Abgeleitete Datei.** Die eine neue Zeile in `docs/user/abdeckung-tests.md` steht in der Sortierung der
  Nachbarzeilen (`TestRenamesArtDesEintrags`, `TestRenamesCommitMengeOhneVorBase`, `TestRenamesFehlerFormen`);
  `make test` hält die Datei gegen ihre Ableitung und war in beiden Mutationsläufen daran nicht rot.
  Geprüft, ohne Befund.
- **Kommentar (`AGENTS.md` §3.7, Prüffrage 6).** Der neue Kommentar trägt eine Zusage („Die Commit-Menge ist
  ‚von HEAD erreichbar, von BASE nicht‘“) und eine Abgrenzung (ein Commit vor BASE trägt nichts bei), dazu
  ein auflösbares Herkunftsfeld (`DC-FA-VCS-001`, `ADR-0115 Entscheidung 2`). Keine Review-Historie, kein
  Befund-Marker, keine Slice-Nummer, kein Mess-Label. Der Konjunktiv („läse die Range diese Commits mit,
  trüge …“) benennt die Fehlerform, die der Test hält, und ist damit Grenze, nicht Deliberation über
  Verworfenes. Im Kern-Fall steht der Grund im Fallnamen; kein Kommentar. Geprüft, ohne Befund.
- **Hexagon, Netz, Suppression.** Keine neue Import-Zeile außer den schon vorhandenen Test-Hilfen; kein
  `//nolint`. Geprüft, ohne Befund.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 0 | — |
| INFO | 0 | — |

Die aus R1 offenen LOW und INFO (L-1, L-2, I-1 bis I-4) bleiben unverändert und gehen ins Register.

## Verdikt

**M-1 und M-2 sind geschlossen.** Beide neuen Testfälle fahren die zugesagte Regel und werden unter der
jeweiligen Mutation aus dem richtigen Grund rot, und nur dort. Die Botschaft behauptet nicht mehr als die
Messung, die Fixtures sind stabil, der Kommentar hält §3.7, der Commit bleibt in seinen drei Dateien. Der
Code ist nach diesem Urteil **reif für die Verifikation**; es war die letzte Review-Runde.
