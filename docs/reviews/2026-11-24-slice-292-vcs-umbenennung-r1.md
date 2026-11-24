# Review-Report: slice-292 — Kern und git-Adapter, `vcs` folgt der Umbenennung (Code R1) — 2026-11-24

**Review-Art:** Code. Geprüft wird der Diff gegen Slice-Plan, ADR, Spezifikation und Hard Rules, also
Maintainability (`v6.18.1` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten). Gegen die DoD
prüft danach der Verifier.

**Gegenstand:** `slice-292`, Commit-Range `6aca94b0..5999fb2b`. Beide Commits liegen nur lokal auf `main`:

- `bea51566`: Stufe 1. Port-Methode `Renames`, Adapter, Kern, Konfigurations-Rand.
- `5999fb2b`: Stufe 2. Randfälle, Index-Angleichung, Handbuch, Sensor-Datei, Spezifikation Schritt 2.

16 Dateien, 2228 Zeilen neu (`git diff 6aca94b0..5999fb2b --stat`). Nicht berührt sind
[ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md), `.d-check.yml`, `AGENTS.md`, `harness/README.md`,
die READMEs und `CHANGELOG.md`; das deckt sich mit der Dateitabelle in §3 des Plans.

**Skill:** `.harness/skills/reviewer.md` @ 1.20.0 (`d1c37c14`).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-24.

**Eingangs-Kontext:**

- Plan `slice-292` §1, §3 (Dateitabelle, Testfälle) und §6 (Risiken).
- [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md), Entscheidung 1–14.
- Architect-Verdikt ([`2026-11-24-slice-292-architect-verdikt.md`](2026-11-24-slice-292-architect-verdikt.md)),
  F-4 bis F-6.
- Lastenheft
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  in 0.114.3 (30 Kriterien) und Spezifikation
  §[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  Schritt 2 und 4.
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.2, §3.7, §3.8.
- Keine früheren Findings am Modul in diesem Slice; die Reports zu `slice-291` betreffen Anforderung und ADR.

**Wie gelesen und gefahren wurde.** Der Diff vollständig, dazu `git.go` im ganzen, `vcs.go` und
`vcs_link_targets.go`. `make test`, `make lint` und `make arch-check` im Stand `5999fb2b`: grün, `0 issues`,
`0 Befund(e)`. 22 Mutationen im Arbeitsbaum, je eine, danach mit `git checkout --` zurückgenommen (Tabelle
unten); dazu drei temporäre Probe-Tests, nach dem Lauf gelöscht. Am Ende sind `git status --short` und
`git diff` leer, nichts wurde committet.

---

## Je Schwerpunkt ein Satz

| # | Schwerpunkt | Urteil |
|---|---|---|
| 1 | Korrektheit des Adapters | Die Rename-Erkennung der Bibliothek wird nicht benutzt (`DiffTree`, `DetectRenames` und `Changes` kommen im neuen Code nicht vor); Kette und Mehrdeutigkeit sind Mengenfunktionen je Commit über Eltern-Zustände und unabhängig von jeder Aufzählungsreihenfolge; ein fehlender Fehlerpfad ist nicht zu finden, die drei Commit-Lesestellen sind redundant, nicht lückenhaft; zwei Lesepfade und eine Stelle der Commit-Menge tragen keinen eigenen Rot-Beleg (M-1, L-1, L-2) |
| 2 | Kern | Die Port-Methode wird nur mit dem Schalter gerufen (M4 rot), die drei Kombinationen lösen wie Entscheidung 10 auf, die Invariante bricht den Lauf (M13 rot); die doppelte Marke selbst ist ungetestet (M-2) |
| 3 | Index-Angleichung | Mit Plan, Spezifikation und Welle-Trigger konsistent; der Vorzustand ist gemessen Exit 2 (M6), der Nachzustand `core-drift-vcs` in beiden Modi, und der `--range`-Zweig ist durch M20 rot gedeckt |
| 4 | Tests | Jedes der 30 Kriterien trägt einen Test; von 22 Mutationen fielen 16 rot und 6 überlebten (ein erster Versuch war ein Übersetzungsfehler und wurde wiederholt); zwei Überlebende zeigen echte Lücken (M-1, M-2), drei sind Randfälle (L-1, L-2), einer ist äquivalent (I-1); die Fixtures sind stabil |
| 5 | Doku-Spiegel | Sieben Grenz-Aussagen am Code geprüft, alle stimmen; die Laufzeit-Aussage lässt aus, was `Renames` zusätzlich walkt (I-2), und die Gitlink-Aussage steht nur im Schlüssel-Absatz (I-4) |
| 6 | Kommentare und Hard Rules | Keine Review-Historie, keine Slice-Nummern, keine Mess-Labels, kein `//nolint`; Rang-Zeiger statt Herkunfts-Prosa; kein Befund |
| 7 | Abgrenzung | Die 16 Dateien der Stat-Liste stehen alle in §3; die sechs ausgenommenen Dateien sind unberührt; kein Befund |

---

## Findings

Die Spalte „Behandlung“ folgt der Vereinbarung zu höchstens zwei Review-Runden: HIGH und MEDIUM werden
eingearbeitet, LOW und INFO gehen ins Register.

| Nr. | Kategorie | Kurz | Behandlung |
|---|---|---|---|
| M-1 | MEDIUM | Commits vor BASE sind aus der Range ausgeschlossen, und kein Test zeigt es | einarbeiten |
| M-2 | MEDIUM | Die doppelte Marke des teilnehmenden Ziels ändert kein Testergebnis | einarbeiten |
| L-1 | LOW | Der Baum einer Merge-Seite und die Art eines gestagten Eintrags tragen keinen eigenen Rot-Beleg | ins Register |
| L-2 | LOW | „Neuanlage lässt die Linie wegfallen“ ist durch keinen Test gefahren | ins Register |
| I-1 | INFO | Der Teilbaum-Skip bei gleicher Kennung ist ungetestet und im Ergebnis äquivalent | ins Register |
| I-2 | INFO | Die Laufzeit-Aussage nennt nur die Länge der Range | ins Register |
| I-3 | INFO | Der Test „AusByteIdentisch“ belegt Gleichheit der zwei Konfigurationen, nicht der Vorfassung | ins Register |
| I-4 | INFO | Eine Aussage ohne den Schlüssel steht im Absatz des Schlüssels; `entryUnreadable` führt die Art-Liste ein zweites Mal | ins Register |

### M-1 — Die Menge „von HEAD erreichbar, von BASE nicht“ hat keinen Test; die Mutation, die BASE-Vorfahren mitwandert, bleibt grün

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 13;
  [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 2 und 13;
  [`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  Schritt 2 „Commit-Menge“; Plan §2, zweiter Punkt (Rot ohne die Regel für jeden Test)
- **Pfad:** `internal/adapter/driven/git/git.go` · „if w.excluded[h] || w.commits[h] != nil {“ (Z. 853)
- **Befund:** Nimmt man `w.excluded[h]` aus der Bedingung, wandert `rangeWalk.visit` die gesamte Historie
  bis zur Wurzel, und `make test` bleibt grün (Mutation M19). Kein Fixture hat ein Ereignis auf einem Pfad
  in einem Commit vor BASE, der einer Quelle in BASE gehört; die Fixtures „Merge: eine Seite zweigt vor
  BASE ab“ und das CLI-Fixture haben keine.
- **Failure-Szenario:** Ein Pfad `a.md`, der lange vor BASE gelöscht und vor BASE neu angelegt wurde, liegt
  in BASE. Wandert der Walk Commits vor BASE mit, trägt `a.md` die Linie „gelöscht“, die nie verschwindet
  (Entscheidung 11), und die Umbenennung `a.md` → `b.md` in der Range wird nicht erkannt: der Nachzug bleibt
  `core-drift-vcs`. Probe-Fixture (temporär, nicht im Repo): `r` mit `a.md` und `n.md`; `r2` löscht `a.md`;
  `base` legt `a.md` mit anderem Inhalt neu an; `head` benennt `a.md` nach `b.md` um. Der Code liefert
  `{a.md: b.md}`; mit M19 liefert er `{}`. Wer den Walk später umbaut, etwa auf eine gemeinsame
  Ausschluss-Menge mit `CommitMessages`, merkt es an keinem Test.
- **Verifizierbar:** ja. Ein Fixture mit Löschen und Neuanlegen vor BASE, Mutation M19 dagegen fahren.
- **Klasse:** ausschluss-menge-ohne-rot-beleg

### M-2 — Die doppelte Marke ändert in keinem Test ein Ergebnis; auf einer Wurzel-Datei als Ende fiele ein Umbiegen durch

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 13;
  [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 10 („die Marke zweimal“);
  [`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  Schritt 4 („doppelter Marke“, „die beiden Formen tragen verschiedene Marken und gleichen sich nie“)
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` ·
  „const renamedTargetMark = normalizedTargetMark + normalizedTargetMark“ (Z. 116)
- **Befund:** Setzt man die Konstante auf eine einfache Marke, bleibt `make test` grün (Mutation M5). Die
  Fälle in `vcs_follow_renames_test.go` tragen alle Enden in einem Verzeichnis; dort unterscheidet sich der
  Endpfad (`d3/a.md`) schon vom Dateinamen (`a.md`), die Marke trägt nichts. Gleich werden die beiden Formen
  nur, wenn das Ende im Wurzelverzeichnis liegt: Endpfad und Dateiname sind dann dieselbe Zeichenkette.
- **Failure-Szenario:** Beide Schlüssel, `a.md` wird nach `x.md` umbenannt (Wurzel), und ein Link in einer
  immutablen Datei wird von `d1/x.md` (nimmt nicht teil, BASE und HEAD) auf `x.md` umgebogen. Mit der
  doppelten Marke: `core-drift-vcs`, wie Entscheidung 10 und das Kriterium „Boundary (beide Schlüssel)“ es
  zusagen. Mit einfacher Marke gleichen beide Seiten `Marke + x.md`, und das Umbiegen auf ein anderes
  Dokument geht durch. Probe-Test (temporär): Befunde `{follow: 1, beide: 1, ignore: 0}`; mit M5 liefert
  der Lauf mit beiden Schlüsseln 0.
- **Verifizierbar:** ja. Ein Randfall in `TestVCSFollowRenamesRandfaelle` mit dem Ende `x.md` im Wurzelverzeichnis.
- **Klasse:** marke-ohne-rot-beleg

### L-1 — Zwei Lesepfade tragen keinen eigenen Rot-Beleg: der Baum einer Merge-Seite ohne Kind-Schritt und die Art eines gestagten Eintrags

- **Kategorie:** LOW
- **Quelle:** [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 13 („den Baum jedes
  Schritts“) und §Fitness Function; Prüffrage 13
- **Pfad:** `internal/adapter/driven/git/git.go` · „return stepEvents{}, fmt.Errorf("baum von commit %s nicht lesbar: %w", c.Hash, err)“
  (Z. 881–883) und „right[e.Name] = fileEntry{hash: e.Hash, symlink: e.Mode == filemode.Symlink}“ (Z. 993)
- **Befund:** Zwei Mutationen bleiben grün. M8 lässt `c.Tree()` den Fehler schlucken: der Test „Baum eines
  Schritts“ fällt trotzdem rot, weil derselbe Baum im Folgeschritt als Eltern-Baum gelesen wird (`pt`, M7
  fällt dort rot). Allein steht der Lesepfad `c.Tree()` bei einem Commit, dessen einziges Kind ein Merge
  ist; das Fixture dafür fehlt. M22 nimmt die Art (Symlink) aus dem Eintrag im Index: nur die
  Baum-Fassung (M9) ist durch „Symlink gegen Datei mit gleichem Blob“ gedeckt.
- **Verifizierbar:** ja. M8 gegen ein Fixture, dessen unlesbarer Baum einem Seiten-Commit mit Merge als
  Kind gehört; M22 gegen ein gestagtes Fixture mit Datei gegen Symlink gleichen Inhalts.
- **Klasse:** lesepfad-maskiert-durch-zweiten-leser

### L-2 — „Eine Neuanlage lässt die Linie wegfallen“ ist nur im Kommentar begründet

- **Kategorie:** LOW
- **Quelle:** [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 11 („Eine Linie, auf der der
  Pfad nicht liegt … ist unberührt“); Prüffrage 13
- **Pfad:** `internal/adapter/driven/git/git.go` · „if ev.created[p] {“ (Z. 637, in `lineSet.land`)
- **Befund:** Die Mutation, die den Zweig abschaltet, bleibt grün (M16). Sie ändert das Ergebnis dort, wo
  eine Seitenlinie, die vor BASE abzweigt, einen Pfad anlegt, den BASE schon trägt, und ihn dann umbenennt.
  Probe (temporär): `p0` mit `n.md`; `base` legt `a.md` an; `c1` (Haupt-Linie) ändert nur anderes; `s1` (von
  `p0`) legt `a.md` mit anderem Inhalt an, `s2` benennt `a.md` nach `q.md` um; `head` ist der Merge von `c1`
  und `s2` und trägt `a.md` und `q.md`. Der Code liefert `{}`; mit M16 liefert er `{a.md: q.md}`, obwohl
  `a.md` in HEAD weiter besteht. Der Code ist hier richtig, die Zeile hängt an keinem Test.
- **Verifizierbar:** ja, mit dem Probe-Fixture.
- **Klasse:** linien-randfall-ohne-rot-beleg

### I-1 — Der Teilbaum-Skip bei gleicher Kennung ist ungetestet; sein Fehlen änderte kein Ergebnis

- **Kategorie:** INFO
- **Quelle:** [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 13 („muss nicht gelesen
  werden“)
- **Pfad:** `internal/adapter/driven/git/git.go` · „if p.left.Hash == p.right.Hash {“ (Z. 796)
- **Befund:** Die Zeile steht im Code und liest einen Teilbaum gleicher Kennung nicht (Antwort auf die Frage
  des Auftrags). M15 nimmt sie heraus: `make test` bleibt grün. Der Skip ist im Ergebnis äquivalent (zwei
  gleiche Bäume haben keine Ereignisse) und nur eine Optimierung; ein Test bräuchte ein unlesbares Objekt,
  das BASE und HEAD nicht lesen und das kein Schritt als neu anlegt, das gibt es nur in einem Merge-Commit,
  der selbst kein Schritt ist. Die Spezifikation formuliert „muss nicht“, also eine Erlaubnis.
- **Verifizierbar:** nein
- **Klasse:** optimierung-ohne-beobachtbare-wirkung

### I-2 — Die Laufzeit-Aussage nennt die Länge der Range; `Renames` läuft zusätzlich über die Historie vor BASE

- **Kategorie:** INFO
- **Quelle:** Prüffrage 18; [`AGENTS.md`](../../AGENTS.md) §5 Regel 13; Plan §6, Risiko 3
- **Pfad:** `docs/user/benutzerhandbuch.md` · „die Laufzeit wächst mit der Länge“ (Z. 2223);
  `internal/adapter/driven/git/git.go` · „isAnc, err := a.isAncestor(*baseHash, *headHash)“ (Z. 920) und
  „excluded, err := a.ancestors(*baseHash)“ (Z. 927)
- **Befund:** `isAncestor` ist eine Tiefensuche, die zuerst den letzten Eltern-Commit weiterverfolgt und bei
  Merge-Commits bis zur Wurzel läuft, bevor sie BASE auf dem anderen Zweig findet; bei einer fremden Basis
  liest sie die ganze Historie von HEAD. Dazu kommt `ancestors(base)` ein zweites Mal neben dem in
  `AllPaths`. Das ist dieselbe Größenordnung wie die Enumeration von `AllPaths` ohne den Schlüssel, also kein
  neuer Aufwand-Typ, aber Handbuch und Sensor-Datei nennen nur die Range.
- **Messung (nur zur Einordnung, nicht die Closure-Messung des Plans):** das gebaute Image, Eigenlauf in
  diesem Repo (3046 Commits, 3 Merges), `--enable vcs`, `ignore-link-targets` an, ohne und mit
  `follow-renames`:

  | Range | ohne | mit |
  |---|---|---|
  | `HEAD~100..HEAD` | 2,4 s | 3,2 s |
  | `HEAD~300..HEAD` | 4,8 s | 6,1 s |
  | `HEAD~1000..HEAD` | 15,9 s | 15,5 s |
  | zwei Merge-Ranges (je 2 Commits) | 3,3 s / 1,5 s | 3,4 s / 1,7 s |

  Höchstens das 1,35-Fache. Die Bedingung des Plans (100 Commits mit einem Merge, dessen eine Seite vor
  BASE abzweigt) erfüllt keine dieser Ranges; die Historie hat nur drei Merges.
- **Verifizierbar:** ja, mit dem Messlauf der Closure-Notiz.
- **Klasse:** laufzeit-aussage-nennt-nur-die-range

### I-3 — „AusByteIdentisch“ belegt, dass fehlender Schlüssel und `false` gleich laufen; die Gleichheit zur Vorfassung trägt die Welle-Probe

- **Kategorie:** INFO
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 15; Kriterium „Boundary (ohne `vcs.follow-renames`
  byte-identisch)“ von
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
- **Pfad:** `internal/adapter/driving/cli/cli_vcs_renames_test.go` ·
  „TestVCS_FollowRenamesAusByteIdentisch“ (Z. 161)
- **Befund:** Beide verglichenen Läufe laufen im neuen Code. Dass ohne Schlüssel weder ein Commit zwischen
  BASE und HEAD gelesen wird (`TestVCS_FollowRenamesFailClosed`, Exit 0) noch `Renames` gerufen wird
  (`TestVCSFollowRenames`, Zähler), trägt die zweite Hälfte des Kriteriums. Die erste Hälfte, gleiche
  Ausgabe wie vor der Erweiterung, trägt `make blackbox-probe` der Welle. Der Test-Name sagt das nicht.
- **Verifizierbar:** nein
- **Klasse:** test-name-sagt-mehr-als-der-test

### I-4 — Zwei Kleinigkeiten gleicher Klasse: eine Aussage ohne Schlüssel im Schlüssel-Absatz, und die Art-Liste steht zweimal in `git.go`

- **Kategorie:** INFO
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 13; Prüffrage 21
- **Pfad:** `docs/user/benutzerhandbuch.md` · „ersetzt ein Gitlink eine immutable Datei an ihrem Pfad, meldet
  das Modul“ (Z. 2226–2227); `internal/adapter/driven/git/git.go` ·
  „return e.Mode == filemode.Regular || e.Mode == filemode.Executable ||“ (Z. 325, in `entryUnreadable`)
- **Befund:** Der Gitlink-Satz gilt mit `--staged` auch ohne den Schlüssel; das Handbuch führt ihn nur im
  Absatz „Was der Schlüssel liest“ (die Sensor-Datei sagt „mit und ohne `vcs.follow-renames`“). Und
  `isFileMode` ist jetzt die eine Antwort auf „ist das ein Datei-Eintrag“, `entryUnreadable` führt die
  drei Arten weiter als eigene Bedingung.
- **Verifizierbar:** nein
- **Klasse:** aussage-am-falschen-ort

---

## Mutationen

„rot“ heißt: `make test` fällt, und der benannte Test nennt die Ursache. Alle Mutationen sind
zurückgenommen.

| M | Mutation | Ergebnis | Fällt in |
|---|---|---|---|
| M1 | Bedingung „Ende nicht in BASE“ entfernt | rot | `TestRenamesMehrdeutigkeit` (Tausch, freigemachter Name, Rotation) |
| M2 | zwei neue Pfade mit gleichem Blob nicht mehrdeutig | rot | `TestRenamesMehrdeutigkeit`, `TestRenamesStaged` |
| M3 | Merge liest nur das erste Elternteil | rot | `TestRenamesMehrdeutigkeit` (zwei Enden, zwei Quellen), `TestRenamesMerge` |
| M4 | `Renames` auch mit `ignore-link-targets` allein gerufen | rot | `TestVCSFollowRenames`, `TestVCSFollowRenamesRandfaelle` |
| M5 | doppelte Marke einfach | **grün** | M-2 |
| M6 | Index-Angleichung zurück (`AllPaths` nimmt jeden Index-Eintrag) | rot | `TestAllPathsStagedGitlinkIstKeinePfad`, `TestVCS_ImmutableDateiDurchGitlinkErsetzt/--staged`; **gemessen: Exit 2** statt 1 |
| M7 | Baum des Eltern-Commits nicht fail-closed | rot | `TestRenamesUnlesbaresObjekt/Baum_des_Eltern-Commits_vor_BASE` |
| M8 | Baum des Commits nicht fail-closed | **grün** | L-1 |
| M9 | Art (Symlink) zählt nicht (Baum) | rot | `TestRenamesArtDesEintrags/Symlink_gegen_Datei_mit_gleichem_Blob` |
| M10 | „gelöscht ohne Ende“ nicht vermerkt | rot | `TestRenamesMerge` (Merge-Seite löscht und legt neu an) |
| M11 | Konfigurations-Rand nicht verdrahtet | rot | `TestDecode_VCSFollowRenames`, `TestVCS_FollowRenames` |
| M12 | Basis kein Vorfahr nicht geprüft (erster Versuch: Übersetzungsfehler, mit `&& false` wiederholt) | rot | `TestRenamesMerge/Basis_kein_Vorfahr_von_HEAD`, `TestVCS_FollowRenamesBasisKeinVorfahr` |
| M13 | Invariante „Quelle zugleich Ende“ nicht geprüft | rot | `TestVCSFollowRenamesPortMethode` |
| M14 | Auflöser baut den Baum nur mit `ignore-link-targets` | rot | `TestVCSFollowRenames`, `TestVCS_FollowRenames` |
| M15 | Teilbaum gleicher Kennung wird gelesen | grün | I-1 (äquivalent) |
| M16 | Neuanlage lässt die Linie nicht wegfallen | **grün** | L-2 |
| M17 | zwei Quellen auf ein Ende zulässig | rot | `TestRenamesMehrdeutigkeit/zwei_Quellen_auf_ein_Ende_(Merge)` |
| M18 | Ende nicht gleichgesetzt | rot | `TestVCSFollowRenames`, `TestVCS_FollowRenames` |
| M19 | ausgeschlossene Commits werden mitgewandert | **grün** | M-1 |
| M20 | die Baum-Enumeration nimmt einen Gitlink | rot | `TestAllPathsGitlinkWirdUebersprungen`, `TestVCS_ImmutableDateiDurchGitlinkErsetzt/--range` |
| M21 | das Ausführbar-Bit zählt | rot | `TestRenamesArtDesEintrags/nur_das_Ausführbar-Bit_ändert_sich` |
| M22 | gestagt: Art (Symlink) zählt nicht | **grün** | L-1 |

Die Probe-Tests zu M-1, L-2 und M-2 liefen gegen den unveränderten Code grün und gegen M19, M16 und M5
rot, mit den oben genannten Zuordnungen. Die Dateien sind gelöscht.

---

## Negativbefunde

Je ein betrachteter Bereich; „geprüft, ohne Befund“.

- **Rename-Erkennung der Bibliothek.** Im neuen Code kommen weder `DiffTree`, `DetectRenames` noch
  `Changes` vor; der Vergleich ist `diffFiles`/`eventsOf` über Tree-Einträge. Geprüft, ohne Befund.
- **Reihenfolge-Unabhängigkeit.** Der Zustand eines Commits ist eine Funktion der Zustände seiner Eltern;
  `mergeStates` vereint mengenweise (kommutativ); `visit` liefert Eltern vor Kindern; `entryNames` sortiert;
  das Ergebnis `finalRenames` ist eine Abbildung. Die drei Salze der Tabelle ändern die Kennungen, die
  Sortierung nach Kennung wird nirgends gelesen. Geprüft, ohne Befund.
- **Zustandsmengen.** `unchanged`, `deleted` und `at` werden nur über Kopien verändert (`with`, `union`,
  Wertempfänger); `advance` gibt den Eltern-Zustand unverändert weiter, wenn der Schritt kein Ereignis hat,
  und nichts schreibt danach. Geprüft, ohne Befund.
- **Merge und fremde Basis.** Merge-Commit ist kein Schritt, trägt nichts bei; Wurzel-Commit liefert `nil`;
  Basis kein Vorfahr liefert eine leere Abbildung ohne Fehler (M12 rot). Geprüft, ohne Befund.
- **fail-closed.** Gelesen wird: Commit (`visit`, Eltern-Commit in `step`, `ancestors`), Baum
  (`Tree()` je Seite), Unterbaum (`diffSubtrees`, `collectSide`, `walkEntries`), Index, BASE- und
  HEAD-Menge (`pathSet`). Jeder Lesefehler wird zurückgegeben, keiner geschluckt. Die drei Commit-Lesestellen
  sind redundant: ein Commit der Range liest `visit` und, als Eltern-Commit, `step`; ein Eltern-Commit vor
  BASE liest `ancestors` und danach `step`. Das ist Absicherung, keine Lücke. Geprüft, ohne Befund.
- **Gleichheit.** Blob-Kennung und Art, das Ausführbar-Bit zählt nicht (M21 rot); Symlink gegen Datei ist
  ungleich (M9 rot); Verzeichnis und Gitlink nehmen nicht teil (`TestRenamesArtDesEintrags`). Geprüft,
  ohne Befund.
- **Kern, drei Kombinationen.** Nur `follow-renames`: teilnehmend → doppelte Marke und Endpfad, sonst roh.
  Beide: teilnehmend wie zuvor, jedes andere wie `ignore-link-targets` allein. Nur `ignore-link-targets`:
  `followed` findet bei leerer Abbildung nichts, die Form bleibt Marke und Dateiname, und die Port-Methode
  wird nicht gerufen. Geprüft, ohne Befund (außer M-2).
- **Index-Angleichung.** Konsistent mit Plan §3, dem Architect-Verdikt F-4 und F-5, dem Spezifikations-Satz
  und dem Welle-Trigger, der die Abweichung benennt. `TrackedPaths` (Modul `tracked`) nimmt weiter jeden
  Index-Eintrag, weil sein Gegenstand der Getrackt-Status ist (Prüffrage 11: anderer Gegenstand, benannter
  Grund). Geprüft, ohne Befund.
- **Konfigurations-Rand.** `FollowRenames` ist ein `bool` wie `IgnoreLinkTargets`; `TestDecode_VCSFollowRenames`
  fährt beide Schlüssel nebeneinander (falscher Typ, kein Wert, `null`). Geprüft, ohne Befund.
- **Doku-Spiegel, sieben Grenz-Aussagen am Code (Regel 13).** (1) „ohne den Schlüssel liest das Modul keinen
  Commit zwischen BASE und HEAD“: `TestVCS_FollowRenamesFailClosed`, Exit 0 ohne, 2 mit. (2) „ein Gitlink
  und ein Verzeichnis nehmen nicht teil“: `TestRenamesArtDesEintrags`, M20. (3) Basis kein Vorfahr „ohne
  Exit 2“: `TestVCS_FollowRenamesBasisKeinVorfahr`, M12. (4) „mit `--staged` ist der einzige Schritt HEAD
  gegen den Index“: `stagedRenames`. (5) „ein unlesbares Objekt ist Exit 2, auch wenn BASE und HEAD lesbar
  sind“: `TestRenamesUnlesbaresObjekt` mit der Vorbedingung `AllPaths`. (6) „von einer Datei liest die
  Erkennung nur die Kennung“: im Adapter wird nie ein Blob geöffnet. (7) Sensor-Datei: Gitlink durch
  `TestVCS_ImmutableDateiDurchGitlinkErsetzt` in allen vier Kombinationen. Stimmen alle.
- **Hexagon-Richtung, Netz, Suppression.** `make arch-check`: 0 Befunde; `git.go` importiert nur den Port;
  kein Netz; kein `//nolint` im Diff. Geprüft, ohne Befund.
- **Kommentare.** Die neuen Kommentare tragen Zusage, Kopplung, Abgrenzung oder Rang-Zeiger; keine
  Review-Historie, keine Slice-Nummern, keine Mess-Labels (Stichwortsuche im Diff). Geprüft, ohne Befund.
- **Fixtures.** Die Adapter-Fixtures legen Objekte direkt im Speicher an, kein Arbeitsbaum, kein Index,
  also kein racily-clean (`BEO-ALL/racily-clean-git-fixture`). Die CLI-Fixtures schreiben keine Datei
  gleich lang neu: jede Neufassung von `adr-x.md` ändert den Link-Pfad und damit die Länge; `mvAt` legt einen
  neuen Pfad an. Geprüft, ohne Befund.
- **Abgrenzung.** Die 16 Dateien aus der Stat-Liste stehen alle in §3; `ADR-0115`, `.d-check.yml`,
  `AGENTS.md`, `harness/README.md`, READMEs und `CHANGELOG.md` sind unberührt. Der Handbuch-Kopf ist nicht
  berührt. Geprüft, ohne Befund.
- **Lint.** `make lint`: `0 issues`. `gofmt -l` meldet auch Dateien außerhalb des Diffs; das Profil führt
  `gofmt` nicht, also kein Konventions-Anker für ein Finding.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 2 | M-1, M-2 |
| LOW | 2 | L-1, L-2 |
| INFO | 4 | I-1, I-2, I-3, I-4 |

## Verdikt

**Der Produktcode ist korrekt, kein Finding betrifft sein Verhalten.** Beide MEDIUM-Findings sind Test-Lücken
bei einer zugesagten Regel: eine Mutation der Regel bleibt grün, und ein konkretes Fixture zeigt, dass sie
einen anderen Ausgang hätte. Sie blockieren nach der Vereinbarung die Verifikation und gehören eingearbeitet
(je ein Testfall, kein Produktcode). Mit diesen beiden ist der Stand reif für die Verifikation; LOW und INFO
gehen ins Register. Der Stand läuft in eine zweite und letzte Runde R2, in der die zwei Testfälle und die
Mutationen M19 und M5 gegen sie geprüft werden.
