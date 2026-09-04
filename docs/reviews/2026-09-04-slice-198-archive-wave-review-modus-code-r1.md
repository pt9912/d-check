# Review-Report: slice-198 — `tools/archive-wave` Modus für eigenständige Review-Archivierung — 2026-09-04

**Review-Art:** Unabhängiger Code-Review (Modul 8, Reviewer-Rolle; kein geteilter Kontext mit der Implementierung).
**Gegenstand:** Commit `d820833` ("feat(planning): slice-198 -- archive-wave Modus fuer eigenstaendige Review-Archivierung (slice-198, welle-90)"), Diff gegen `461da52`. Geänderte Dateien: `AGENTS.md`, `Makefile`, `docs/plan/planning/in-progress/slice-198-archive-wave-review-modus.md`, `tools/archive-wave/{main.go,archive.go,collect.go,stub.go,Makefile}`, `tools/archive-wave/{review_mode_test.go,slice_mode_test.go}`.
**Skill/Modell/Datum:** Claude Sonnet 5, Reviewer-Rolle, 2026-09-04.
**Eingangs-Kontext:** `AGENTS.md` vollständig gelesen; kein Zugriff auf den Implementierungs-Kontext dieser Session — ausschließlich Diff, Bestand und eigene Läufe.

## Unabhängig nachgefahrene Läufe (echte Ausgabe, kein behaupteter Exit-Code)

- `make archive-wave-test` → `ok  archive-wave  0.019s` (grün).
- `make gates` → `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` (zehn Gates, Coverage 94,70 % ≥ 93 %, Semgrep 0 Findings).
- `make fullbuild` → `[fullbuild] green — image-hash sha256:1c9a353425ff0751c689f9f25d1c129efc338a1692a51a0aeb0ea939ba1aa85e`.
- `git status --porcelain` vor und nach allen Läufen leer — kein Lauf hat den echten Bestand verändert.
- Dry-Run (`make archive-wave REVIEW=<datei> `, kein `APPLY`) gegen alle elf echten Reviews unter `docs/reviews/` einzeln nachgefahren — alle elf liefern eine saubere Vorschau ohne Schreibzugriff.
- Zusätzlich (über den Slice-Anspruch hinaus, um die `ExtractFullHeading`-Design-Entscheidung *und* den `-apply`-Pfad gegen den vollen Bestand zu verifizieren, nicht nur gegen die drei unit-getesteten Formen): alle elf echten Reviews in eine Scratch-Kopie des Repos dupliziert und dort `-apply` gegen jeden einzeln gefahren. Alle elf archivieren korrekt (Zip + Stub in `docs/reviews/archiv/`, Original gelöscht); jeder Stub trägt die **volle**, unverstümmelte Original-Überschrift. Das reale Repo blieb dabei unberührt (Kopie in `/tmp/…/scratchpad/review-apply-test`, nicht committet).

## Findings

### F-1 · HIGH · `AGENTS.md` §3.7 (fünf Kommentar-Klassen; Herkunfts-Prosa/Slice-Nummer explizit verboten, Bestandsgrenze deckt nur Alt-Kommentare) · `tools/archive-wave/review_mode_test.go:94-98`

**Befund:** Der Doc-Kommentar über `TestExtractFullHeading_RealeUeberschriftenformen` zitiert die eigene Slice- und Plan-Absatznummer als Beleg:

```go
// TestExtractFullHeading_RealeUeberschriftenformen belegt slice-198 §2
// Punkt 4/§4: ExtractFullHeading darf das fuehrende Wort einer
// uneinheitlichen Review-Ueberschrift nicht verschlucken -- anders als
// ExtractTitle, das fuer das Slice-/Welle-Praefixschema gebaut ist.
// Gegen zwei der elf realen Formen geprueft (Fixture-Kopie).
```

„belegt slice-198 §2 Punkt 4/§4" ist Herkunfts-Prosa mit Slice-Nummer und Plan-Absatzverweis — keine der fünf erlaubten Klassen (Zusage · Kopplung · Abgrenzung · Rang-Zeiger · Grenze). §3.7 nennt „keine Slice-Nummern" ausdrücklich, und die Bestandsgrenze gilt nur für **vor** der Schärfung geschriebene Zeilen: „Neuzugänge fallen überall unter den Anker." Dieser Kommentar ist ein Neuzugang dieses Diffs.

Das ist **derselbe Befund-Typ**, der in diesem Tool bereits einmal HIGH war: der Review-Report zu slice-196 (`done/welle-89/archiv.zip`, `docs/reviews/2026-09-04-slice-196-archive-wave-slice-modus-code-r1.md`) blockierte den Merge mit F-1 über exakt diese Form, unter anderem an einer wortgleich aufgebauten Stelle: „`slice_mode_test.go:115`: „belegt **slice-196** §4: ein Slice mit …" — dieselbe Form." Der hier vorliegende Kommentar wiederholt das Muster „belegt sliceNNN §N: …" fast wortgleich, obwohl der Vorgänger-Review in derselben Werkzeug-Codebasis (`tools/archive-wave`) genau davor gewarnt hat und die Empfehlung aussprach, die Herkunft gehöre „in Commit-Botschaft/Closure-Notiz, nicht in den Code-Kommentar."

Zusätzlich ist die Zahlenangabe im Kommentar selbst falsch: „Gegen zwei der elf realen Formen geprueft" — die Testtabelle enthält tatsächlich **drei** Fälle (Zeilen 100–107: „Review-Report: Change Request…", „Review — Lastenheft-CR…", „Review Release-Prep v0.70.0…"), was auch der Slice-DoD (§4) und die Closure-Notiz (§9) korrekt mit „drei" beziffern. Der Code-Kommentar selbst ist damit nicht nur klassenwidrig, sondern auch sachlich falsch — ein Kommentar, der beschreiben soll „was da ist" (Baseline-Merksatz), trifft es nicht.

**Verifizierbar:** ja — `grep -n "slice-198" tools/archive-wave/*.go` zeigt genau diese eine Stelle; `sed -n '100,107p' tools/archive-wave/review_mode_test.go` zeigt drei Tabellenzeilen gegen den im Kommentar behaupteten „zwei"; der Vorgänger-Befund liegt wortgleich referenzierbar in `docs/plan/planning/done/welle-89/archiv.zip` unter `docs/reviews/2026-09-04-slice-196-archive-wave-slice-modus-code-r1.md` (F-1).
**Klasse:** `kommentar-herkunfts-prosa-slice-nummer`

### F-2 · LOW · Kommentar-Aktualität · `tools/archive-wave/main.go:52` (Doc-Kommentar über `validateModeFlags`)

**Befund:** Der Kommentar direkt über der Funktion wurde beim Erweitern auf drei Flags nicht mitgezogen:

```go
// validateModeFlags erzwingt, dass genau eines von -welle/-slice gesetzt
// ist -- weder beides (mehrdeutig, welcher Modus?) noch keines (kein Ziel).
```

Die Funktion prüft seit diesem Diff drei Flags (`welle, sliceID, reviewFile`), nicht mehr zwei; die Fehlermeldung im Funktionskörper selbst wurde korrekt auf „genau eines von -welle, -slice oder -review" aktualisiert, der Doc-Kommentar eine Zeile darüber nicht. Ein Kommentar, der beschreibt, was **nicht mehr** da ist, ist funktional harmlos hier (die Fehlermeldung im Code trägt die korrekte Aussage), aber eine spätere Erweiterung um einen vierten Modus würde denselben Kommentar ein weiteres Mal stehen lassen.

**Verifizierbar:** ja — `sed -n '52,53p' tools/archive-wave/main.go` zeigt den Zwei-Flags-Wortlaut neben der Drei-Flags-Signatur eine Zeile darunter.
**Klasse:** `kommentar-veraltet-nach-erweiterung`

### F-3 · LOW · Testabdeckung · `tools/archive-wave/main.go:245-273` (`runReview`, Dry-Run-Zweig)

**Befund:** Der Dry-Run-Zweig von `runReview` kehrt zurück, **bevor** `ApplyReview` (und damit `ExtractFullHeading`) je aufgerufen wird — die Titel-Extraktion läuft ausschließlich im `-apply`-Pfad. Damit deckt der im Commit behauptete „Dry-Run smoke-getestet gegen alle elf echten Reviews" (Commit-Botschaft, Closure-Notiz §9) **nicht** die Design-Entscheidung, die dieser Slice eigentlich absichern soll (`ExtractFullHeading` vs. `ExtractTitle`) — dieser Pfad wird im Bestand nur an einer einzigen realen Überschrift (`TestRunReview_Apply`) plus drei Fixture-Kopien (`TestExtractFullHeading_RealeUeberschriftenformen`) unit-getestet, nicht an allen elf. Die Aussage selbst ist wörtlich korrekt (Dry-Run wurde smoke-getestet), aber ein Leser könnte daraus lesen, dass damit auch die Überschriften-Extraktion für den vollen Bestand geprüft wurde — das stimmt nicht ohne den zusätzlichen `-apply`-Lauf, den dieser Review nachgeholt hat (siehe oben, alle elf bestätigt korrekt). Kein Korrektheits-Fehler — die Funktion verhält sich, wie oben gegen alle elf real verifiziert, richtig — sondern eine Lücke zwischen behaupteter und tatsächlich durch das Werkzeug selbst (nicht durch diesen Review) abgedeckter Prüfmenge.
**Verifizierbar:** ja — `runReview`s `if !apply { … return nil }`-Zweig in `main.go` endet vor dem `ApplyReview`-Aufruf; `grep -n "ExtractFullHeading\|ApplyReview" tools/archive-wave/main.go` zeigt beide nur im `apply`-Zweig.
**Klasse:** `smoke-test-deckt-beworbenen-pfad-nicht`

### F-4 · INFO · `docs/plan/planning/in-progress/slice-198-archive-wave-review-modus.md` §2 Punkt 3

**Befund:** Der ursprüngliche Plan-Text in §2 („Vorgehen") sagt noch „Stub-Titel via `ExtractTitle`" — überholt durch die tatsächliche Umsetzung (`ExtractFullHeading`, korrekt in §4 als „Präzisiert" nachgetragen und in §9 korrekt beschrieben). Kein Verstoß — §4/§9 sind die maßgeblichen, aktuellen Abschnitte, und das Auseinanderfallen von Plan (§2) und Ergebnis ist hier sogar die dokumentierte Pointe des Slice (§5, Risiko 1: das Risiko trat vor der Implementierung ein und wurde aufgefangen) — aber ein Leser, der nur §2 liest, bekäme den falschen Eindruck.
**Verifizierbar:** ja — §2 Punkt 3 vs. §4/§9 derselben Datei.
**Klasse:** `plan-vs-ergebnis-abschnitt-nicht-nachgezogen`

## Positiv geprüft (kein Finding)

- **`FindReview`** (collect.go): exakter Dateiname-Lookup unter `docs/reviews/`, `os.Stat`-Fehler und `IsDir()` beide auf denselben Fehlertext gemappt — korrekt für den Zweck (kein Verzeichnis als Review akzeptiert).
- **`ApplyReview`** (archive.go): baut das Zip **vor** dem Löschen aus dem Originalinhalt (`buildZip` vor `os.Remove`), schreibt den Stub mit demselben Basisnamen in `ReviewArchiveDir`, löscht nur die Originaldatei, liefert einen korrekten `Move` für `RewriteRepo`. Durch eigenen `-apply`-Lauf gegen alle elf realen Reviews verifiziert (kein Bestandslauf im echten Repo, siehe oben).
- **`runReview`** (main.go): `SliceIDFromPath(filename)` korrekt auf dem **Dateinamen** (nicht dem vollen Pfad) aufgerufen — konsistent mit der Funktionssignatur, die selbst `filepath.Base` anwendet.
- **`ExtractFullHeading` vs. `ExtractTitle`**: Behauptung im Commit/Plan verifiziert. Gegen die elf echten Überschriften geprüft: `h1RE` (ExtractTitle) verschluckt bei „# Review-Report: Change Request…" das Wort „Review-Report:" (Gruppe 1 beginnt erst bei „Change Request…") und bei „# Review — Lastenheft-CR…" das Wort „Review" (Gruppe 1 beginnt bei „— Lastenheft-CR…") — die Design-Begründung trägt. `fullH1RE` liefert für alle elf Formen die vollständige, unveränderte Überschrift.
- **Testqualität (Mutationscheck von Hand):** Ein vergessenes `os.Remove` in `ApplyReview` würde `TestRunReview_Apply` (Zeile 47-49, `os.Stat`-Existenzprüfung am alten Pfad) rot machen. Eine Regression von `ExtractFullHeading` auf `ExtractTitle`-Verhalten würde `TestExtractFullHeading_RealeUeberschriftenformen` an allen drei Fällen rot machen. `TestRunReview_RejectsSliceReview` konstruiert eine real existierende, mit gültiger Überschrift versehene Datei (`buildReviewFixture`) — ohne die `SliceIDFromPath`-Prüfung in `runReview` würde der Aufruf mit `apply=true` erfolgreich durchlaufen (verifiziert: die Fixture-Datei trägt „# Review slice-137\n", `ExtractFullHeading` liefert dafür einen nicht-leeren Titel), der Test schlägt also tatsächlich auf die geprüfte Bedingung an, nicht auf einen Nebeneffekt.
- **`validateModeFlags`/`TestValidateModeFlags`:** Erzwingt korrekt „genau eins von drei" (Zähl-Schleife statt Boolean-XOR, notwendig für drei statt zwei Flags). Die Testtabelle deckt alle 2³ = 8 Kombinationen der drei Flags vollständig ab (0, je-1-von-3, je-2-von-3, alle-3) — keine Lücke.
- **`main()`-Switch:** `case *welle != "" / case *sliceID != "" / default` — bei durch `validateModeFlags` bereits erzwungenem „genau eins gesetzt" korrekt und eindeutig; `runWelle`/`runSlice` im Diff inhaltlich unverändert (nur Signatur-Aufrufer angepasst).
- **Makefiles:** Sowohl `tools/archive-wave/Makefile`s `run`-Target als auch das Root-`Makefile`s `archive-wave`-Target verlangen über eine POSIX-kompatible Zähl-Schleife (`n=$$((n+1))`) korrekt genau eines von `WELLE`/`SLICE`/`REVIEW`; nachgefahren gegen „kein Argument" (Exit 2, korrekte Meldung), „zwei Argumente" (Exit 2), und den Grenzfall „eine Variable explizit leer gesetzt vs. gar nicht gesetzt" (`WELLE= SLICE=slice-198 REVIEW=`) — in Make sind beide Formen ununterscheidbar und werden beide korrekt als „nicht gesetzt" behandelt (kein Fehlverhalten, da konsistent).
- **§3.1 (Docker/make-only):** Alle eigenen Nachvollzugs-Läufe dieses Reviews liefen über `make archive-wave`/`make archive-wave-test`/`make gates`/`make fullbuild` bzw. `docker run` auf dem bereits von `make` gebauten Image — kein roher Host-`go`. Die Implementierung selbst benutzt laut Commit-Botschaft dieselben Targets.
- **`AGENTS.md`-Zeile `make archive-wave`:** Beschreibt den neuen `REVIEW`-Modus korrekt (Pfad, Ablehnung von `slice-<NNN>`-Dateinamen, Stub-Unterschied zu Slice-/Wellen-Reviews) und korrigiert dabei zutreffend die veraltete Passage zum `SLICE`-Modus (alter flacher `done/slice-<NNN>-archiv.zip`-Pfad → aktuelles `docs/plan/planning/done/wellenlos/`, per Code-Konstante `WellenlosArchiveDir` bestätigt). Der neu eingefügte `d-check:ignore`-Marker auf `docs/reviews/archiv/` ist notwendig und korrekt begründet: ohne ihn meldet `codepaths` `codepath-missing`, da das Verzeichnis vor der ersten Anwendung (slice-199) nicht existiert — durch testweises Entfernen des Markers selbst nachgefahren und wiederhergestellt.
- **DoD/§5/§9 der Slice-Datei:** Checkbox-Stand ehrlich (Review/Verifikation/Closure-Freigabe korrekt unchecked, da noch ausstehend). §5-Risiken alle mit einem der drei zulässigen Ausgänge versehen; „entfallen"-Begründungen sind jeweils nachvollziehbar (WIP-Limit 1, kein Bestandsbezug in diesem Slice). §9-Closure-Notiz beschreibt den tatsächlichen Diff akkurat (Funktionsliste stimmt mit dem Diff überein).

## Verdikt

**Merge-blocking: JA — wegen F-1.**

Die Implementierung selbst ist korrekt: alle Gates grün, eigene Läufe bestätigen `ApplyReview`/`FindReview`/`runReview` und die `ExtractFullHeading`-Design-Entscheidung gegen den vollständigen Bestand aller elf realen Reviews (nicht nur die drei unit-getesteten Formen), keine Bestandsmutation durch diesen Review verursacht. F-1 ist jedoch derselbe Befund-Typ, der im unmittelbaren Vorgänger-Slice (slice-196, dasselbe Werkzeug, dieselbe Welle-Abstammung) bereits HIGH und blockierend war — inklusive derselben Kommentar-Form („belegt sliceNNN §N: …") und trotz der dort ausgesprochenen Empfehlung, Herkunftsbezüge in Commit-Botschaft/Closure-Notiz statt in Code-Kommentare zu legen. Empfehlung: den Kommentar über `TestExtractFullHeading_RealeUeberschriftenformen` auf seine tatsächliche Aussage kürzen (was die Funktion testet, ohne Slice-/Absatz-Bezug) und dabei die Zahlenangabe auf „drei" korrigieren; F-2/F-3/F-4 sind nicht blockierend, aber vor Merge günstig mitzunehmen, da sie an derselben Stelle bzw. Datei hängen.
