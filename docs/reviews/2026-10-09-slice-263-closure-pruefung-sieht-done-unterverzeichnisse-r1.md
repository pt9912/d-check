# Review R1 — slice-263: Closure-Prüfung über Unterverzeichnisse und Inhalts-Ausnahme (Produkt)

- **Review-Art:** Code. Geprüft wird der Produkt-Diff gegen den Slice-Plan (`slice-263`, §1 Ziel und
  Abgrenzung, §3, §6, §8 samt Spiegel-Liste), die Entscheidungen (`ADR-0048`, `MR-025`) und die
  Hard Rules `AGENTS.md` §3.7/§3.8 sowie §5 Regel 3/13/15. Die DoD-Abhakung gehört nicht zu diesem
  Review.
- **Gegenstand:** `slice-263` · Produkt-Commit `487ee4de` (Kontext: `1bce45f0..e98c5fd2`, Plan und
  Übergang).
- **Skill:** `reviewer.md` @ 1.19.0 (`980d6314`)
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** Slice-Plan `slice-263`; `DC-FA-PLAN-001` und `DC-FA-STRUCT-001` im Lastenheft;
  Spezifikation §`DC-FA-PLAN-001.a` C1–C3, §`DC-FA-STRUCT-001.a` Schritte 1/2, §2-Schema, §4-Zeile
  `SPEC-039`, §8; vorherige Befunde am Modul: Review R1/R2 von `slice-260` (F-1 ist der Auslöser
  dieses Slice). Im Code gelesen: `internal/hexagon/core/rules/{planning,structure,scan}.go`,
  `internal/adapter/driven/fs/fs.go`, `internal/adapter/driven/configyaml/configyaml.go`
  (`applyClosure`, `structureBedingungsFehler`), `internal/hexagon/core/model/config.go`
  (`ClosureConfig`, `StructureRule`, `Identity`), `internal/hexagon/core/app/diagnose.go`,
  `internal/adapter/driving/cli/config_template.go`, die neuen Tests.
- **Proben** (Wegwerf-Kopie im Scratchpad, `make test IMAGE=dcheck-review-probe`; Arbeitsbaum
  unverändert, Probe-Image danach entfernt):
  1. Test mit `closure.dir: docs/fehlt/` (Schrägstrich am Ende) und einem Dateisystem, dessen `List`
     auf dem Verzeichnis scheitert, Erwartung = Meldung des Stands vor `487ee4de`. **Rot:**
     `Message:Closure-Verzeichnis docs/fehlt fehlt oder ist unlesbar (fail-closed)` bei
     `File:docs/fehlt/ Target:docs/fehlt/`.
  2. Test `recursive: true` mit `done/node_modules/slice-002-b.md` (dünn) und `done/a/b/slice-003-c.md`
     (dünn): genau ein Befund auf `a/b/…` — grün, `SKIP_DIRS` und zweistufiger Abstieg arbeiten.
  3. Mutation `if false && isSkipDir(e.Name)` in `closureCandidates`, Probe-Tests entfernt, ganze
     Suite: `ok github.com/pt9912/d-check/internal/hexagon/core/rules` — **kein** Test wird rot.

## Findings

### F-1 — MEDIUM — Fehlerpfad ohne die Schlüssel nicht byte-identisch

- **quelle:** `DC-FA-PLAN-001` (Spezifikation §2-Zeile `planning.closure.recursive`: „Aus ⇒ nur das
  Verzeichnis selbst, Befundsatz byte-identisch"); Skill-Frage 2
- **pfad:** `internal/hexagon/core/rules/planning.go` · „`"Closure-Verzeichnis "+err.Error()+" fehlt oder ist unlesbar (fail-closed)"`"
  zusammen mit `closureCandidates` · „`cur := path.Join(dir, sub)`" / „`return nil, errors.New(cur)`"
- **befund:** Die Meldung nennt jetzt `path.Join(dir, "")`, also den **bereinigten** Pfad, während
  `file`/`target` weiter das rohe `dir` tragen. `applyClosure` normalisiert `dir` nicht (nur `/`-Präfix
  und `..` sind verboten). Ein Repo mit `dir: docs/done/` oder `./docs/done`, dessen Verzeichnis fehlt
  oder unlesbar ist, bekommt ohne einen der neuen Schlüssel eine andere Meldung als vor dem Commit
  (Probe 1) — und Befund-Datei und Meldung nennen zwei verschiedene Schreibweisen desselben Pfads.
- **verifizierbar:** ja — Probe 1 als Test; `make blackbox-probe` sieht es nicht, weil keine Fixture
  ein unbereinigtes, fehlendes `closure.dir` führt.
- **klasse:** `byte-identitaet-nur-auf-gemessenen-formen`

### F-2 — MEDIUM — Die neuen Schlüssel stehen nur in der Spezifikation, das Lastenheft schließt sie aus

- **quelle:** `AGENTS.md` §5 Regel 3 (neue/geänderte `DC-*`-Anforderungen nur im Lastenheft);
  Source Precedence (`harness/README.md` §Source precedence — die Spezifikation darf das Lastenheft
  schärfen, nicht erweitern)
- **pfad:** `spec/lastenheft.md` · „Geprüft wird jede Datei
  in `planning.closure.dir`, deren Basisname `planning.closure.glob` matcht" und „Geprüft wird
  **ausschließlich** `planning.closure.dir`"; `spec/lastenheft.md` (`DC-FA-STRUCT-001`) ·
  „Abgezogen wird `exempt-paths` der Regel."; Slice-Plan §3 (keine Lastenheft-Zeile)
- **befund:** Das Lastenheft nennt die Kandidatenmengen beider Module abschließend: jede
  basisnamen-passende Datei in `closure.dir` wird geprüft, bei `structure` wird nur `exempt-paths`
  abgezogen; auch die Exit-2-Listen am Config-Rand sind dort aufgezählt. `skip-pattern` nimmt eine
  basisnamen-passende Datei heraus, `recursive` prüft Dateien außerhalb von `closure.dir` selbst —
  beides steht nur in der Spezifikation (Rang 2), die damit dem Rang-1-Wortlaut widerspricht. Jeder
  frühere Schlüssel dieser Module (`closure.glob`, `placeholder`, `exempt-section-pattern`,
  `exempt-expect-count`, `tasks-ignore-pattern`, `open-tasks-require-marker-section`) steht im
  Lastenheft. Wer gegen das Lastenheft abnimmt, liest „jede Datei" und findet ein Ventil, das es nicht
  zusagt.
- **verifizierbar:** nein — `grep -c skip-pattern spec/lastenheft.md` ergibt 0; das Urteil über den
  Widerspruch ist keines, das ein Gate fällt.
- **klasse:** `vertrag-nur-im-niedrigeren-stratum`

### F-3 — MEDIUM — Beispiel in `--print-config` trifft den Stub der Baseline-Form nicht

- **quelle:** `AGENTS.md` §5 Regel 13 (Grenze/Beispiel gegen den Gegenstand prüfen); `DC-FA-CLI-005`
  (die Vorlage ist öffentliche Ausgabe)
- **pfad:** `internal/adapter/driving/cli/config_template.go` ·
  „`skip-pattern: '^> \*\*ARCHIVIERT'  # RE2 gegen den rohen INHALT`" (beide Blöcke, `closure` und
  `structure`)
- **befund:** RE2 bindet `^` ohne `(?m)` an den Anfang des **ganzen** Inhalts. Ein Stub nach
  Baseline-Form beginnt mit der Überschrift, der Marker steht auf Zeile 3 (gemessen an
  `done/wellenlos/slice-112-…md`; die Tests dieses Commits nutzen deshalb `(?m)^> \*\*ARCHIVIERT\*\*`).
  Wer das Beispiel übernimmt, nimmt keinen einzigen Stub aus — `planning` meldet dann jeden Stub als
  `closure-note-missing`, `structure` jeden als `section-missing`. Die Spezifikation sagt zur
  Anker-Semantik des Musters nichts.
- **verifizierbar:** ja — ein Test mit dem Vorlagen-Muster gegen einen Stub der Baseline-Form.
- **klasse:** `beispiel-nicht-am-gegenstand-geprueft`

### F-4 — MEDIUM — Die stille Richtung der Inhalts-Ausnahme ist unbenannt

- **quelle:** Skill-Frage 18 (Grenzen-Liste ohne ihre größte Lücke); `AGENTS.md` §3.8
- **pfad:** `spec/spezifikation.md` · „die Kandidaten ab, deren **rohen Inhalt** es trifft" und
  Slice-Plan §6 · „**Ein Stub ohne Marker**"
- **befund:** Das Muster läuft gegen den rohen Inhalt einschließlich Fenced-Code und Zitaten. Ein
  **Volltext**, der die Stub-Form zeigt (etwa ein Slice über das Archivieren, der die Vorlage mit
  `> **ARCHIVIERT** — Volltext:` in einem Codeblock zitiert), fällt still aus der Prüfung; die
  Nullmengen-Regel fängt nur den Fall, dass **alles** ausgenommen ist. Benannt sind die laute
  Richtung (§6: Stub ohne Marker wird als Volltext geprüft, also rot) und die unlesbare Datei, nicht
  aber diese stille. Heute trägt keine Datei unter `docs/plan/planning/` außerhalb der Stubs eine
  Zeile, die mit dem Marker beginnt — die Lücke ist latent, nicht eingetreten.
- **verifizierbar:** ja — ein Kandidat mit vollständiger, dünner Notiz und dem Marker in einem
  Fenced-Block bleibt unter `skip-pattern` befundfrei.
- **klasse:** `ausnahme-verschluckt-volltext-still`

### F-5 — MEDIUM — Negativtest für `SKIP_DIRS` unter `recursive` fehlt

- **quelle:** Skill-Frage 13; Spezifikation C2 · „steigt das Listing in jedes Unterverzeichnis ab
  (die `SKIP_DIRS` ausgenommen)"
- **pfad:** `internal/hexagon/core/rules/planning.go` · „`if isSkipDir(e.Name) {`";
  `internal/hexagon/core/rules/planning_closure_rekursiv_test.go`
- **befund:** Die Ausnahme ist Teil des neuen Vertrags, aber kein Test hält sie: mit
  ausgeschaltetem `isSkipDir` bleibt die ganze Suite grün (Probe 3). Eine Änderung, die in
  `node_modules/`, `build/` oder `vendor/` unter `closure.dir` absteigt, liefe unbemerkt durch — oder
  umgekehrt eine, die ein legitimes Unterverzeichnis namens `build` still ausließe.
- **verifizierbar:** ja — Probe 2 als Test, gegengeprüft mit der Mutation aus Probe 3.
- **klasse:** `neuer-vertrag-ohne-negativtest`

### F-6 — MEDIUM — Commit-Botschaft behauptet mehr, als gemessen wurde

- **quelle:** `AGENTS.md` §5 Regel 15; Skill-Frage 8 (`BEO-ALL/commit-message-overclaims-work`)
- **pfad:** Commit `487ee4de` · „ohne Abstieg bzw. ohne Abzug werden alle neuen Tests aus dem
  richtigen Grund rot" und „Drei opt-in-Schluessel, ohne die der Befundsatz byte-identisch bleibt"
- **befund:** Drei der neuen Tests bleiben ohne Abzug grün, weil sie die Gegenrichtung absichern:
  `TestClosureSkipPattern_UnlesbareDateiBleibtKandidat`,
  `TestStructureSkipPattern_UnlesbareDateiBleibtKandidat`,
  `TestStructureSkipPattern_OhneSchluesselMeldungUnveraendert`; dasselbe gilt für die erste Hälfte von
  `TestClosureRecursive_VerzeichnisNameOhneSchalterUnveraendert`. Als Wächter sind diese Tests
  richtig, die Botschaft zählt sie aber mit. Die Byte-Identität ist über 20 Vergleiche der
  `blackbox-probe` gemessen und wird als allgemeine Eigenschaft behauptet; F-1 ist die N+1-te Form.
  Die Commit-Botschaft ist eingefroren; die Korrektur gehört in die Closure-Notiz.
- **verifizierbar:** ja — die genannten Tests gegen den Stand vor dem Commit (bzw. mit ausgeschaltetem
  Abzug) fahren.
- **klasse:** `botschaft-ueberdehnt-messung`

### F-7 — MEDIUM — Symlink-Unterverzeichnis wird unter `recursive` still übergangen

- **quelle:** `AGENTS.md` §3.8; Skill-Frage 15
- **pfad:** `internal/hexagon/core/rules/planning.go` · „`if recursive && e.Kind == driven.KindDir {`";
  Spezifikation C2 · „steigt das Listing in jedes Unterverzeichnis ab"
- **befund:** `fs.List` klassifiziert per `Lstat`; ein Symlink auf ein Verzeichnis ist `KindSymlink`,
  wird nicht betreten und ist, weil sein Name den Basisnamen-Filter nicht trifft, auch kein Kandidat.
  Seine Slices fallen still aus der Prüfung. Zyklen sind dadurch ausgeschlossen, und das Verhalten ist
  vertretbar, aber die Spezifikation sagt „jedes Unterverzeichnis" und nennt die Grenze nicht. Bei
  `structure` steht sie ausdrücklich („keine Symlinks", Schritt 2).
- **verifizierbar:** ja — Test über den echten Dateisystem-Adapter mit einem symbolisch verlinkten
  Unterverzeichnis.
- **klasse:** `grenze-der-neuen-achse-unbenannt`

### F-8 — LOW — Grund-Code-Zeile `SPEC-039` nicht nachgezogen

- **quelle:** `MR-025` (Spiegel vor dem Editieren)
- **pfad:** `spec/spezifikation.md` · „Kandidat im `planning.closure.dir` (Filter:
  `planning.closure.glob`, sonst `planning.slice-glob`)"
- **befund:** Die §4-Zeile beschreibt die Auslöser von `closure-note-missing` weiter ohne
  Unterverzeichnis-Kandidaten, ohne unlesbares Unterverzeichnis und ohne die Leere nach `skip-pattern`.
  Bei der letzten Änderung der Kandidatenmenge (`closure.glob`, Historie 2026-08-10) wurde genau diese
  Zeile mitgezogen. Die Spiegel-Liste in §8 des Plans führt sie nicht. Ebenso nennt die §2-Zeile
  `structure[].files` beim Leerlauf weiter nur `exempt-paths`.
- **verifizierbar:** nein (Doku-Drift, kein Gate).
- **klasse:** `spiegel-unvollstaendig`

### F-9 — LOW — Kommentar erzählt den Vorzustand

- **quelle:** `AGENTS.md` §3.7 („Ein Kommentar beschreibt, was da ist")
- **pfad:** `internal/hexagon/core/rules/planning.go` · „Ohne recursive wird ein Verzeichnis behandelt
  wie bisher"
- **befund:** Der Kommentar trägt eine Zusage (ein filter-treffendes Verzeichnis ist ohne `recursive`
  Kandidat und meldet sich als unlesbar), formuliert sie aber relativ zu einem Vorzustand. „Bisher"
  hat für einen späteren Leser keinen Bezugspunkt. Eine Klasse ist da; LOW, nicht HIGH.
- **verifizierbar:** nein.
- **klasse:** `kommentar-mit-chronik-wort`

### F-10 — INFO — Ein unlesbares Unterverzeichnis verdeckt alle übrigen Befunde

- **quelle:** Maintainability
- **pfad:** `internal/hexagon/core/rules/planning.go` · „`return nil, err`" in `closureCandidates`
- **befund:** Scheitert `List` auf einem Unterverzeichnis, bricht die ganze Kandidatenwahl ab und es
  bleibt ein einziger `closure-note-missing` auf `dir`; die übrigen Notizen werden in diesem Lauf
  nicht gemessen. Das ist fail-closed und von der Spezifikation gedeckt, für den Leser einer roten
  Ausgabe aber nicht offensichtlich.
- **verifizierbar:** ja (`TestClosureRecursive_UnlesbaresUnterverzeichnisFailClosed` zeigt genau einen
  Befund).
- **klasse:** `fail-closed-verdeckt-restmenge`

## Negativbefunde

- **Hexagon-Richtung (ADR-0005):** `rules` importiert neu nur `errors` aus der Standardbibliothek;
  kein Adapter-Import. Ohne Befund.
- **Netz, Suppression:** kein Netzzugriff, kein `//nolint`, keine gesenkte Schwelle. Ohne Befund.
- **Kandidatenmenge ohne Schlüssel, Normalpfad:** ohne `recursive` ist `rel = path.Join("", e.Name)
  = e.Name`, der Filter läuft wie zuvor über jeden Eintrag, auch über ein Verzeichnis und einen
  Symlink mit passendem Namen; ohne `skip-pattern` gibt `closureSkip` die Menge unverändert zurück,
  `closureSkipNote` liefert `""`; die `structure`-Leerlauf-Meldung bleibt ohne Muster wörtlich gleich
  (Test vorhanden). Abweichung nur im Fehlerpfad, siehe F-1.
- **`out := names[:0]`:** `names` ist eine lokale, frisch aufgebaute Slice aus `closureCandidates`;
  es gibt keinen zweiten Halter, das Überschreiben ist folgenlos. Ohne Befund.
- **Zyklen, Tiefe:** Symlinks werden nicht betreten (F-7), echte Verzeichnisse können keinen Zyklus
  bilden; zwei Ebenen Abstieg in Probe 2 grün. Ohne Befund.
- **Fehlerformen des Lesewegs (Frage 19):** unlesbare Datei bleibt Kandidatin (beide Module,
  getestet, mit falscher Implementierung rot); unlesbares Unterverzeichnis meldet fail-closed
  (getestet); alles ausgenommen ⇒ Nullmengen-Befund mit Nennung des Musters (beide Module,
  getestet); leere Unterverzeichnisse tragen nichts bei, die Nullmengen-Regel greift danach;
  ungültiges Muster ⇒ Exit 2 am Config-Rand (beide Module, getestet), im Kern zusätzlich
  abgesichert (`planning`) bzw. per `MustCompile` wie jeder andere RE2-Schlüssel des Moduls
  (`structure`). Ohne Befund außer F-5.
- **Konsistenz zwischen den Modulen (Frage 11):** beide behandeln „unlesbar" gleich (Kandidat,
  fail-closed) und gleichen den Leerlauf-Befund gleich an. Symlink-Dateien: `structure` schließt sie
  aus, `planning` nimmt sie als Kandidaten — Bestand, nicht durch diesen Diff entstanden.
- **Regel-Identität:** `skip-pattern` geht wie `exempt-paths` nicht in `Identity()` ein; zwei Regeln,
  die sich nur darin unterscheiden, kollidieren (Exit 2). Gleiches Verhalten wie beim Geschwister
  `exempt-paths`, nicht wie bei `exempt-section-pattern`; mangels Grandfathering-Paar-Fall ohne Befund.
- **Weitere Spiegel:** `--doctor`-Text zu `closure-note-missing` ist allgemein genug („fehlt bzw. ist
  leer"); `--suggest-config` erzeugt keinen `closure`-Block mit Unterverzeichnissen und keine
  `structure`-Regel; YAML-Durchreichung beider Module getestet. Handbuch gehört in die Release-Prep.
  Ohne Befund außer F-8.
- **Spezifikation am Code (Regel 13, Schritt 18):** „ein Verzeichnis ist dann Abstieg, nicht
  Kandidat", „Ohne den Schlüssel … ist Kandidat und meldet sich als unlesbar", „unlesbare Datei bleibt
  Kandidatin und meldet `section-missing` auf sich", „die Meldung nennt dann beide Abzüge", Exit-2-
  Listen um `skip-pattern` ergänzt — am Code nachgelesen, stimmen. Abweichungen: F-1, F-3, F-7.
- **Kommentare (§3.7):** übrige neue Kommentare in `model/config.go`, `configyaml.go`,
  `planning.go`, `structure.go` tragen Zusage, Kopplung oder Grenze; keine Slice-Nummern, keine
  Befund-Marker. Ohne Befund außer F-9.
- **Abgrenzung des Plans:** `reviews.done-dir` ist nicht angefasst, keine Repo-Konfiguration geändert
  (`slice-264`), kein `**`-Glob-Umbau. Ohne Befund.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 7 | F-1 bis F-7 |
| LOW | 2 | F-8, F-9 |
| INFO | 1 | F-10 |

Wiederkehrende Finding-Klasse für die Closure-Notiz: `beispiel-nicht-am-gegenstand-geprueft` /
`byte-identitaet-nur-auf-gemessenen-formen` — beide sind die Regel-13- bzw. Regel-15-Klasse
(`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`, `BEO-ALL/commit-message-overclaims-work`).

## Verdikt

**Nicht freigegeben.** Kein HIGH; die sieben MEDIUM blockieren typischerweise. Am schwersten wiegen
F-2 (der Vertrag auf Rang 1 widerspricht dem neuen Verhalten — das ist eine Plan-Lücke, §3 führt das
Lastenheft nicht; ob das Lastenheft nachgezogen wird oder die Schlüssel als Verfeinerung gelten, ist
eine Entscheidung vor dem Code, kein Nachtrag) und F-3 (das öffentliche Beispiel funktioniert am
Gegenstand nicht). F-1 und F-5 sind kleine Code-/Test-Nachzüge, F-4 und F-7 Grenz-Sätze in der
Spezifikation, F-6 gehört in die Closure-Notiz. Kein Befund begründet eine Rückführung des Slice.
