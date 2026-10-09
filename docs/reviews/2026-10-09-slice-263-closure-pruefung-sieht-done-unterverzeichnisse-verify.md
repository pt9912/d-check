# Verifikation — slice-263: Closure-Prüfung über Unterverzeichnisse und Stub-Ausnahme (Produkt)

- **Rolle:** Verifier (Modul 11). Die Frage ist: „Bauen wir es richtig?" Maßstab sind DoD, Spec und Plan.
- **Gegenstand:** `slice-263` · Range `e98c5fd2..e87dc084`. Code-tragend sind `487ee4de` (Feature),
  `5c3febc1` (R1-Einarbeitung) und `e87dc084` (R2-Einarbeitung, Spec-Wortlaut und Test-Kommentare).
  Der Claim-Commit `e98c5fd2` ist ein reiner Rename (100 %) plus Roadmap-Ruhe-Marker.
- **Datum:** 2026-10-09 · **Modell-ID:** claude-opus-5-5
- **Eingangs-Kontext:** Slice-Plan (§1–§3 mit zwei Plan-Änderungen nach R1/R2, §6, §8), Reports R1 und
  R2, die Commit-Botschaften als Sensor-Belege des Implementers, `spec/lastenheft.md` 0.99.1
  (`DC-FA-PLAN-001`, `DC-FA-STRUCT-001`), `spec/spezifikation.md` (C1/C2, `structure` Schritte 1/2,
  §2-Schema, `SPEC-039`, Historie).
- **Vorbehalt:** §7 ist leer, die DoD-Haken stehen offen. Das ist der Stand vor der Closure. DoD 5 ist
  deshalb Closure-Schuld und keine DoD-Verletzung.
- **Proben:** Alle Code-Proben liefen in Wegwerf-Kopien per `git archive HEAD` im Scratchpad, mit eigenem
  Image-Tag (`IMAGE=dcheck-mut-<n>`). Im Repo selbst liefen nur `make gates`, `make blackbox-probe` und
  `make build`. Danach habe ich das Mutations-Image entfernt. Während der Läufe hat eine parallele
  Sitzung `564077eb` committet (CR sf-connector, `slice-265`/`slice-266` in `open/`); verifiziert ist
  `e87dc084`. `git status` zeigt nur diesen Bericht.

---

## Sensor-Läufe (selbst gefahren)

| Lauf | Ergebnis |
|---|---|
| `make gates` (HEAD `e87dc084`) | **Exit 0.** Alle zehn Glieder grün; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`, semgrep 0 Findings, `planning-check` 0 Befunde |
| `make blackbox-probe REF=e98c5fd2` | **Exit 0**: „byte-identisch über 20 Vergleiche (Vorher e98c5fd2, stdout/stderr/Exit getrennt, Kanarienlauf vorher und nachher)" |
| Zusatz: Closure-Profil, Vorher- gegen Nachher-Image | `--config .d-check.closure.yml --enable planning --enable structure --enable spans --enable reviews`, Standard/`--json`/`--doctor`: alle drei Ausgaben per `cmp` byte-gleich, Exit 0 beide |
| `make build` (danach) | Exit 0. Das stellt `d-check:latest` wieder her, das die Probe neu gebaut hatte |

Ich habe den Closure-Lauf ergänzt, weil die Fixtures der Probe den Closure-Pfad nicht tragen
(`ids`, `links`, `sauber`, `targets`, `repo` mit dem Default-Profil). Erst dieser Lauf zeigt die
Byte-Identität für genau die Fähigkeit, die der Slice ändert.

## DoD-Prüfung

### DoD 1 — `planning.closure.recursive` und `planning.closure.skip-pattern` — **BESTÄTIGT**

- **Verfeinerung in der Spezifikation:** C2 nennt die folgenden Punkte:
  - den Abstieg ohne `SKIP_DIRS`;
  - den Basisnamen-Filter;
  - „Verzeichnis ist Abstieg, nicht Kandidat";
  - den Abzug nach rohem Inhalt;
  - die unlesbare Datei als Kandidatin;
  - beide Grenzen (Symlink, zitierter Marker);
  - die Nullmenge über die **gesamte** Menge.

  Die §2-Zeilen und `SPEC-039` sind nachgezogen. Am Code (`closureCandidates`, `closureSkip`,
  `closureSkipNote`) stimmen alle Aussagen.
- **Konfig-Validierung, Exit 2:** Im Code steht `applyClosure` → `regexp.Compile`. Am Binary
  geprüft (Scratch-Repo, `skip-pattern: '^(['`):
  `d-check: error: bad1.yml: planning.closure.skip-pattern "^([" ist kein gültiges Regex: …`, **exit=2**.
  Die Meldung nennt den Schlüssel, wie es das Akzeptanzkriterium verlangt.
- **Tests, rot ohne die Änderung:** siehe §Bewusstes Brechen (M1, M2, M4–M8).

### DoD 2 — `structure[].skip-pattern` — **BESTÄTIGT**

- Spec: Schritt 1 (Exit-2-Liste) und Schritt 2 (Abzug, unlesbare Datei, Grenze, Nullmenge nach beiden
  Abzügen). Die §2-Zeilen `structure[].files` und `structure[].skip-pattern` sind nachgezogen.
  Lastenheft: Exit-2-Aufzählung und zwei neue Kriterien.
- Exit 2 am Binary: `structure[0]: skip-pattern "^([" ist kein gültiges Regex: …`, **exit=2**.
- Tests rot ohne die Änderung: M3, M7, M10.

### DoD 3 — Ausgabe unverändert, Spiegel, Gates — **BESTÄTIGT** (mit INFO V-3)

- Die Byte-Identität ohne die Schlüssel ist doppelt belegt: durch `blackbox-probe` und durch den
  Closure-Profil-Vergleich (siehe oben).
- **`--print-config`:** `docker run d-check:latest --print-config` enthält `recursive: false` und
  `skip-pattern: '(?m)^> \*\*ARCHIVIERT'` im `closure`-Block (Zeilen 151/152). Im
  `structure`-Block steht die `skip-pattern`-Zeile (Zeile 214).
- **Die Vorlage parst und trifft echte Stubs.** Ich habe den `planning`-Block der Ausgabe wörtlich
  einkommentiert. Die Abweichungen von der Vorlage:
  - `recursive: true`;
  - `glob: 'slice-*.md'`;
  - `heading`/`marker` dieses Repos.

  Dazu kam eine `structure`-Regel mit der wörtlichen `skip-pattern`-Zeile der Vorlage. Gefahren habe ich
  das gegen eine `git archive`-Kopie des Repos:
  - mit beiden Mustern: **0 Befunde**, Exit 0;
  - ohne die Muster: **458 Befunde** auf 229 Dateien. Alle 229 tragen `^> **ARCHIVIERT`
    (`grep -L` leer) und liegen unter `done/wellenlos/` und `done/welle-*/`;
  - Gegenprobe gegen die Lesart „skip-pattern verschluckt alles": In einem Volltext unter
    `done/wellenlos/` (`slice-249-…`) habe ich die Closure-Notiz verdünnt und die DoD-Überschrift
    entfernt. Mit beiden Mustern meldet der Lauf genau **2 Befunde** auf diese Datei
    (`closure-note-thin`, `section-missing`). Die Volltexte werden also geprüft, die Stubs nicht.
- Die Bestandszahlen des Plans habe ich nachgezählt: 98 Slice-Dateien unter `done/wellenlos/`, 80 mit
  Marker, 18 ohne. Unter `done/welle-*/` trägt jede Slice-Datei den Marker.
- **Übrige Spiegel:** `--suggest-config` ist laut Plan §8 kein Spiegel. Das Benutzerhandbuch fehlt
  (`grep skip-pattern\|recursive` leer). Das entspricht Plan §8 und `AGENTS.md` §5 Regel 17
  (Release-Prep), ist also keine Verletzung. Zur `--doctor`-Hinweiszeile siehe V-3.
- `make gates` ist grün (siehe oben).

### DoD 4 — Review und Verifikation — **TEILWEISE** (siehe V-1)

R1 und R2 liegen als Reports vor, diese Verifikation ist die zweite Hälfte. Das letzte Review-Verdikt
lautet allerdings **„Nicht freigegeben, Nachzug klein"** (R2). Der Nachzug `e87dc084` ist von keiner
Review-Runde gesehen worden.

### DoD 5 — Closure-Notiz, Register, Risiko-Ausgänge, Paarungen — **OFFEN (Closure-Schuld)**

Erwartbar vor der Closure. Mitzunehmen sind die folgenden Punkte:
- R1 F-6 (Botschaft behauptet zu viel) ist laut R2 an die Closure-Notiz übergeben.
- Plan §8 kündigt für `reviews.done-dir` einen Folge-Slice an, der bei der Closure geschnitten wird.
  Ob `open/slice-265-reviews-zusage-kennung-leerlauf.md` (`564077eb`, parallele Sitzung)
  dieser Slice ist, habe ich nicht geprüft.
- Beide Risiken in §6 stehen auf `*(offen)*` und brauchen je einen der drei Ausgänge.

## Akzeptanzkriterien → Test

| Kriterium (Lastenheft 0.99.1) | Test | Rot ohne Fix |
|---|---|---|
| PLAN Happy Path (Unterverzeichnisse; ohne Schlüssel ungelesen) | `TestClosureRecursive_SiehtUnterverzeichnis` | M1 |
| PLAN Boundary (Stub aus, Volltext geprüft) | `TestClosureSkipPattern_NimmtStubAus` | M1, M2 |
| PLAN Boundary (alles ausgenommen ⇒ Nullmenge mit Nennung) | `TestClosureSkipPattern_LeereMengeFailClosed` | M2, M8 |
| PLAN Boundary (unlesbare Datei bleibt Kandidatin) | `TestClosureSkipPattern_UnlesbareDateiBleibtKandidat` | M6 |
| PLAN fail-closed (unlesbares Unterverzeichnis, Pfad in der Meldung) | `TestClosureRecursive_UnlesbaresUnterverzeichnisFailClosed` | M1 |
| PLAN fail-closed (`skip-pattern` ⇒ Exit 2, Schlüssel genannt) | `TestDecode_ClosureFehler` „skip-pattern RE2" | M7, prüft aber nur `err != nil` (V-2) |
| STRUCT Boundary (Stub aus, Volltext geprüft) | `TestStructureSkipPattern_NimmtStubAus` | M3 |
| STRUCT Boundary (alles ausgenommen ⇒ `section-missing` auf dem Glob) | `TestStructureSkipPattern_LeereMengeFailClosed` | M3 |
| STRUCT Boundary (unlesbare Datei bleibt Kandidatin) | `TestStructureSkipPattern_UnlesbareDateiBleibtKandidat` | M10 |
| STRUCT ohne Schlüssel byte-identisch | `TestStructureSkipPattern_OhneSchluesselMeldungUnveraendert` + blackbox-probe | — (Erhaltungs-Test) |
| STRUCT fail-closed (Exit 2) | `TestDecode_StructureFehler` „skip-pattern RE2" | M7 |
| Plan R1 F-1 (ungereinigtes `dir` byte-identisch) | `TestClosureDir_UngereinigtMeldungUnveraendert` | M5 |
| Plan R1 F-5 (`SKIP_DIRS`) | `TestClosureRecursive_SkipDirsBleibenUnbetreten` | M4 |

Jedes neue Akzeptanzkriterium ist durch einen Test gedeckt.

## Bewusstes Brechen

Jede Mutation lief allein in einer eigenen `git archive`-Kopie mit `make test IMAGE=dcheck-mut-<n>`. Die
unveränderte Kontrollkopie lief **grün** (Exit 0).

| # | Mutation | Rote Tests (alle anderen grün) | Grund in der Fehlermeldung |
|---|---|---|---|
| M1 | kein Abstieg (`if false && recursive && …KindDir`) | `SiehtUnterverzeichnis`, `VerzeichnisNameOhneSchalterIstKandidat`, `UnlesbaresUnterverzeichnisFailClosed`, `ClosureSkipPattern_NimmtStubAus` | „mit recursive: genau closure-note-thin auf …/wellenlos/slice-002-b.md erwartet, got []"; „unlesbares Unterverzeichnis ⇒ … got []"; das Verzeichnis `slice-900-dir.md` meldet sich als Kandidat |
| M2 | kein Abzug im `planning` (`closureSkip` gibt ungefiltert zurück) | `ClosureSkipPattern_NimmtStubAus`, `…_LeereMengeFailClosed` | der Stub `wellenlos/slice-002-b.md` meldet; statt der Nullmenge erscheint ein Befund auf `slice-001-a.md` |
| M3 | kein Abzug in `structure` (`structureSkipped` immer `false`) | `StructureSkipPattern_NimmtStubAus`, `…_LeereMengeFailClosed` | `section-missing` auf dem Stub `slice-001-a.md` |
| M4 | `SKIP_DIRS` betreten (`if false && isSkipDir`) | `ClosureRecursive_SkipDirsBleibenUnbetreten` | „node_modules/ darf nicht betreten werden, got … node_modules/slice-002-b.md … closure-note-thin" |
| M5 | ungereinigtes `dir` bereinigt (`cur := path.Clean(dir)`) | `ClosureDir_UngereinigtMeldungUnveraendert` | „Meldung verändert: want "Closure-Verzeichnis docs/fehlt/ …"" |
| M6 | unlesbare Datei wird verschluckt (`err != nil \|\| re.Match`) | `ClosureSkipPattern_UnlesbareDateiBleibtKandidat` | „unlesbar ⇒ Kandidat, closure-note-missing auf …/slice-002-b.md, got []" |
| M7 | Validierung beider `skip-pattern` entfernt | `TestDecode_StructureFehler`, `TestDecode_ClosureFehler` | „skip-pattern RE2: ungültige …-Config akzeptiert" |
| M8 | Nullmengen-Meldung ohne Nennung des Musters (`closureSkipNote` → "") | `ClosureSkipPattern_LeereMengeFailClosed` | „… mit Nennung der Ausnahme" |
| M10 | unlesbare Datei in `structure` verschluckt (`return true`) | `StructureSkipPattern_UnlesbareDateiBleibtKandidat` | „Befund auf done/wellenlos/slice-001-a.md", stattdessen Leerlauf-Befund auf dem Glob |

Jeder Test wird aus dem behaupteten Grund rot, und es wird kein fremder Test rot.

## Plan-Konformität und Abgrenzung (§1, §3)

- Die geänderten Dateien entsprechen genau §3: Kern-Regeln `planning`/`structure`, Konfig-Modell,
  YAML-Adapter, Konfig-Vorlage, Lastenheft und Spezifikation. Die Plan-Änderungen nach R1 und R2 sind
  jeweils **vor** dem Code-Commit committet (`f4606192` vor `5c3febc1`, `cca5006d` vor `e87dc084`).
- Die Abgrenzung ist eingehalten. Nicht angefasst wurden:
  - `.d-check.yml` und `.d-check.closure.yml` (gehören zu slice-264);
  - das Modul `reviews`;
  - ADR-0048/0051.
- Die Kopf-Felder `Verantwortlich:`/`Autor:` stehen auf `pt9912`. Jeder Commit trägt `slice-263`.
- Hard Rules:
  - kein `//nolint`;
  - Kommentare ohne Slice-, Befund- oder Chronik-Bezug (die Treffer auf `slice-` in den Tests sind
    Fixture-Dateinamen);
  - Spec-Straten ohne ADR-, Slice- oder Wellen-Token (`doc-check` grün).

## Befunde

### V-1 — MEDIUM — Der R2-Nachzug ist von keiner Review-Runde gesehen

- **pfad:** `e87dc084` (`spec/lastenheft.md`, `spec/spezifikation.md`, zwei Test-Dateien)
- **befund:** R2 schließt mit „Nicht freigegeben, Nachzug klein". Eine Runde, die den Nachzug freigibt,
  gibt es nicht. Inhaltlich habe ich die R2-Punkte gegengelesen:
  - F-1: `skip-pattern` steht in der Exit-2-Liste von `DC-FA-STRUCT-001`.
  - F-2: „Gesamtmenge" steht in `SPEC-039`, in C2 und im fail-closed-Absatz des Lastenhefts. Der Code
    aggregiert erst und zählt dann, ein leeres Unterverzeichnis ergibt keinen Befund.
  - F-3: Die Test-Kommentare sind ohne Chronik-Wort.
  - F-5: Zwei neue Kriterien und die Historie der Spezifikation sind ergänzt.

  Ich finde keinen Rest. DoD 4 verlangt aber ein *durchgeführtes* Review, und das freigebende Urteil
  fehlt formal.
- **vorschlag:** Entweder eine kurze R3 über `e87dc084`, oder die Closure-Notiz hält ausdrücklich fest,
  dass der Wortlaut-Nachzug ohne erneute Review-Runde schließt und dass die Verifikation ihn
  gegengelesen hat. Das entscheiden Planner oder Auftraggeber, nicht der Verifier.

### V-2 — LOW — Der Exit-2-Test prüft nicht, dass die Meldung den Schlüssel nennt

- **pfad:** `configyaml_test.go` · `TestDecode_ClosureFehler`/`TestDecode_StructureFehler`, Fall
  „skip-pattern RE2"
- **befund:** Das Lastenheft-Kriterium verlangt „eine Meldung, die den Schlüssel nennt". Die Tests prüfen
  nur `err != nil`. Eine Mutation, die den Schlüsselnamen aus der Meldung nimmt, bliebe grün. Am Binary
  ist das Kriterium heute erfüllt (beide Meldungen oben nennen `skip-pattern`). Die Lücke ist eine
  Lücke im Test, nicht im Verhalten, und sie gilt für die bestehenden Fälle dieser Tabellentests
  genauso.
- **vorschlag:** `strings.Contains(err.Error(), "skip-pattern")` für die zwei neuen Fälle, oder die
  Restlücke benennen.

### V-3 — INFO — Die `--doctor`-Hinweiszeile von `closure-note-missing` kennt die neuen Ursachen nicht

- **pfad:** `internal/hexagon/core/app/diagnose.go:142`, „…oder Closure-Verzeichnis fehlt bzw. ist leer —
  fail-closed"
- **befund:** Das unlesbare Unterverzeichnis und die durch `skip-pattern` geleerte Menge fehlen hier,
  `SPEC-039` nennt beide. Die Zeile wird ohne die Schlüssel ausgegeben. Sie zu ändern bräche die
  Byte-Identität von `--doctor`, die DoD 3 verlangt. Dass sie stehen bleibt, ist deshalb vertretbar,
  steht aber nirgends.
- **vorschlag:** in der Closure-Notiz benennen, oder mit dem nächsten Release ändern, das die Ausgabe
  ohnehin verschiebt.

### V-4 — INFO — `structureSkipped` kompiliert das Muster je Datei

- **pfad:** `structure.go` · `regexp.MustCompile(r.SkipPattern)` in `structureSkipped`
- **befund:** `planning` kompiliert einmal pro Lauf, `structure` einmal je Kandidat. Dazu kommt, dass
  jede Datei einmal für den Abzug und einmal für die Prüfung gelesen wird. Die Validierung am Rand macht
  `MustCompile` panikfrei. Kein Verstoß gegen eine Zusage, nur eine Asymmetrie zwischen den beiden
  Modulen.

### V-5 — INFO — Lifecycle `open → in-progress` ohne `next`

- **pfad:** `e98c5fd2` (`rename open => in-progress`)
- **befund:** Die Baseline-State-Machine kennt `open → in-progress` nicht. `Verantwortlich:` war vorher
  gesetzt (`cabfde4b`), und `planning-check` ist grün. Der Befund gehört in die Planner-Sicht. Er ist
  hier nur notiert, weil er im Commit-Verlauf sichtbar ist.

## Negativbefunde (geprüft, ohne Befund)

- Ohne `recursive` bleibt ein Verzeichnis mit passendem Namen Kandidat und meldet sich unlesbar, wie
  zuvor. Mit `recursive` wird es betreten (M1 belegt beide Richtungen).
- Die Sortierung über relative Pfade ist stabil (`sort.Strings` nach dem Abzug, DC-QA-02). Die
  Byte-Identität der Proben bestätigt das.
- Für einen Volltext, der den Marker zitiert, und für den Symlink auf ein Unterverzeichnis stehen die
  Grenzen in Spec, Kommentar und §6. Die Kommentar-Aussage zur Lstat-Sicht stimmt mit der
  `KindDir`-Abfrage überein.
- R1 F-10 (ein unlesbares Unterverzeichnis verdeckt die Restmenge) ist von der Spec gedeckt
  (fail-closed). Die Meldung nennt den Pfad.

## Kategorie-Summary

| Schwere | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 (V-1) |
| LOW | 1 (V-2) |
| INFO | 3 (V-3, V-4, V-5) |

## Verdikt

**DoD 1–3 bestätigt.** Das gilt für das Verhalten am Binary, die Tests und das bewusste Brechen.
Jeder getestete Fix wird ohne sich selbst aus dem richtigen Grund rot. Die Byte-Identität ohne die
Schlüssel ist belegt, auch auf dem Closure-Pfad, und `make gates` ist grün.

**Offen sind DoD 4 und DoD 5.** DoD 4 ist formal offen, weil kein freigebendes Review-Urteil über den
R2-Nachzug vorliegt (V-1). DoD 5 ist Closure-Schuld: Notiz, Register, zwei Risiko-Ausgänge, R1 F-6 und
der Folge-Slice für `reviews.done-dir`.

Es gibt keine DoD-Verletzung im Produkt.
