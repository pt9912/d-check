# Review-Report — slice-186 (`vcs` meldet den reinen Rename auch über `--range`)

- **Review-Art:** Code (gegen Plan, ADR-0024/ADR-0016, `DC-FA-VCS-001`/`.a`, Hard Rules)
- **Gegenstand:** `8117758~2..HEAD` — Commits `f01b4d9` (Plan) · `8117758` (Lifecycle-Move `open/`→`in-progress/`) · `84c0ca8` (Fix + Test + Kommentar-Korrekturen + Spezifikation + Dokumente) · `4e9cbe8` (Test: Rename aus der immutablen Klasse heraus). HEAD `4e9cbe8`.
- **Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) v1.13.0
- **Modell-ID:** claude-sonnet-5
- **Datum:** 2026-08-31
- **Eingangs-Kontext:** Slice-Plan [`slice-186`](../plan/planning/in-progress/slice-186-vcs-rename-im-range.md); eingehender Befund [`2026-08-31-befund-ai-harness-course-vcs-rename.md`](../plan/cr/2026-08-31-befund-ai-harness-course-vcs-rename.md); [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in) und die Verfeinerung [`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs); [ADR-0024](../plan/adr/0024-vcs-immutable-gate.md), [ADR-0016](../plan/adr/0016-adr-immutable-gate.md); [`AGENTS.md`](../../AGENTS.md) §3 (bes. §3.5/§3.7), §4/§5; vorherige Findings am selben Modul: [`2026-06-29-slice-053-vcs-code-r2.md`](2026-06-29-slice-053-vcs-code-r2.md), [`2026-06-29-slice-053-vcs-doc-first-r1.md`](2026-06-29-slice-053-vcs-doc-first-r1.md), [`2026-06-28-slice-052-immutable-r1.md`](2026-06-28-slice-052-immutable-r1.md), [`2026-06-28-slice-052-immutable-r2.md`](2026-06-28-slice-052-immutable-r2.md); Beobachtungs-Register [`BEO-023`](../plan/planning/observations.md), [`BEO-024`](../plan/planning/observations.md), [`BEO-012`](../plan/planning/observations.md).

Repo unverändert (`git status --short` leer, HEAD `4e9cbe8`); alle Proben liefen entweder gegen das gebaute `d-check:latest`-Image oder in einer git-Repo-Kopie unter der Scratchpad-Wurzel, außerhalb des Repos.

---

## Eigener Lauf (Ausgabe, nicht behauptet)

| Lauf | Ausgabe |
|---|---|
| `make gates` | Build + Tests grün, `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`, semgrep `0 findings`, `d-check: 638 Datei(en) geprüft, 0 Befund(e)` (×2, `doc-immutable`/`doc-planning`-artige Läufe), Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` |
| Bestandsmessung, unabhängig nachgefahren (`--range 361ed39..HEAD`, 61 Commits statt der behaupteten 60) über `docs/plan/adr/[0-9]*.md` | `d-check: 638 Datei(en) geprüft, 0 Befund(e)` — reproduziert die Commit-Botschaft |
| Fünf End-zu-Ende-Proben gegen `d-check:latest` in einer Scratch-Repo-Kopie (eigenes `.d-check.yml` mit `vcs.paths`, unabhängig vom Repo-Bestand aufgebaut): reines Löschen · reiner Rename · Rename **mit** starker Umformulierung · Rename **aus** der Klasse heraus · `--staged`-reiner-Rename | alle fünf melden `core-drift-vcs` auf dem **alten** Pfad, Exit 1 — keine der fünf ist stillgeblieben |

Die dritte Zeile deckt mehr als die Commit-Botschaft selbst belegt (die nennt nur „die drei Kontrollen bewegen sich nicht" + den einen Fix-Fall) und bestätigt zusätzlich den vierten Fall aus `4e9cbe8` end-to-end statt nur über die core-`fakeVCS`-Unit.

---

## Findings

### F-1 · MEDIUM · `AGENTS.md` §3.7 (Herkunfts-Prosa/Review-Historie) · `.githooks/pre-commit:12-16`

**Befund:** Der neue Kommentar-Absatz trägt einen vollständigen, im Präteritum
indikativisch erzählten Satz über einen **behobenen** Vorfall: *„und dort lief
die Zusage aus DC-FA-VCS-001 einmal auseinander — ein reiner Rename meldete nur
ueber --staged."* Das ist nicht die sanktionierte Konjunktiv-Form „was müsste
passieren, damit die Zusage bricht" (Baseline
[`grundlagen-harness-dateien.md` §Was ein Kommentar trägt](../../.harness/baseline/v5.12.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte),
Zeitform-Test), sondern eine Chronik dessen, was einmal (und jetzt nicht mehr)
der Fall war — genau die Baseline-Beispielform für Herkunfts-Prosa
(*„die frühere Zusage wurde ersetzt"*). Der tragende Satz davor (*„Gleiches Gate
heisst aber nicht von selbst gleiche Antwort"*) und danach (*„Die
Uebereinstimmung ist eine gehaltene Eigenschaft, keine, die aus dem gemeinsamen
Gate folgt"*) sind für sich genommen vollständige Abgrenzung/Grenze-Aussagen;
der Vorfall-Satz dazwischen liefert dem Leser, der den Hook **gleich ändert**,
keine zusätzliche Handlungsgrundlage — er begründet nicht die Regel, sondern
erzählt ihre Entstehung.
**Verifizierbar:** nein — kein Gate prüft Kommentar-Klassen (Urteil, kein
`grep`, wie der Skill selbst festhält).
**Klasse:** `kommentar-herkunfts-prosa`

### F-2 · MEDIUM · `AGENTS.md` §5 (Botschaft verallgemeinert über die Messung hinaus, `BEO-009`) · Commit `84c0ca8` Body + `internal/hexagon/port/driven/vcs.go:1-8`

**Befund:** Commit `84c0ca8` behauptet Vollständigkeit: *„Die zwei falschen
Zusagen sind weg"* (DoD-Wortlaut identisch), und zählt namentlich
`internal/adapter/driven/git/git.go` und `.githooks/pre-commit` auf. Eine
dritte Stelle trug **denselben** Fehlschluss und wurde vom Audit nicht
gefunden: der Port-Kommentar in `internal/hexagon/port/driven/vcs.go:2-5`
sagt seit slice-053 unverändert *„Renames werden pfadbasiert als Löschung des
alten Pfads + Hinzufügung des neuen erkannt … (kein
Inhalts-Ähnlichkeits-Matching …)"* — eine Aussage, die für den `--range`-Pfad
vor diesem Fix ebenso falsch war wie die beiden korrigierten (go-gits
Default-Diff **maß** Inhalts-Ähnlichkeit). Die Datei ist in diesem Diff nicht
angefasst (`git log -- internal/hexagon/port/driven/vcs.go` zeigt zuletzt
`d39dd3f`, lange vor `slice-186`); der Kommentar ist heute **zufällig** wieder
wahr, weil der Fix ihn einholt — nicht, weil er geprüft wurde. Ein einfacher
`grep -rn "Inhalts-Ähnlichkeit\|pfadbasiert" --include=*.go` hätte ihn
gefunden. Die „zwei"-Zählung ist damit eine Untermessung der eigenen
Bestandsprüfung, nicht nur eine Wortwahl — der Skill klassifiziert genau diese
Richtung als MEDIUM (*„suche die N+1-te Form"*).
**Verifizierbar:** ja — der `grep`-Befund oben ist reproduzierbar; kein Gate
fängt das.
**Klasse:** `botschaft-vollstaendigkeit-ungeprueft`

### F-3 · INFO · `AGENTS.md` §3.7 (Herkunfts-Prosa, kurzes Fragment) · `internal/adapter/driven/git/git.go:96`

**Befund:** Der neue `diffTrees`-Kommentar ist über weite Strecken sauberer
Zusage-Klasse (er nennt exakt, was garantiert ist und was mit eingeschalteter
Rename-Erkennung bräche — die sanktionierte Form aus der Baseline-Tabelle).
Der Parenthese-Zusatz *„(DC-FA-VCS-001, der Fall trat real ein)"* hängt hinter
die eine erlaubte Herkunfts-ID ein zweites, narratives Element — dass der Fall
„real eintrat" ist eine Chronik-Aussage, kein auflösbares Feld. Anders als F-1
ist das ein knappes, in eine sonst tragfähige Zusage eingebettetes Fragment
(vier Wörter) — dieselbe Gestalt wie ein bereits präzedierter Fall in
[`2026-08-30-slice-179-strukture-teilmenge-review.md`](2026-08-30-slice-179-strukture-teilmenge-review.md)
(dort I-1, ebenfalls INFO), der dort ebenso eine kurze Zuschreibung neben einer
sonst vollständigen Abgrenzung war.
**Verifizierbar:** nein.
**Klasse:** `kommentar-herkunfts-prosa`

### F-4 · LOW · Maintainability (Grammatik) · `docs/plan/cr/2026-08-31-cr-ai-harness-course-observations-relational.md:45-46`

**Befund:** *„Er steht hier stehen, weil dieses Dokument den CR **wie
empfangen** hält"* — doppeltes Verb, der Satz ist nicht wohlgeformt (vermutlich
„Er bleibt hier stehen" oder „Er steht hier weiterhin" gemeint). Rein
sprachlich, kein Gate betroffen (`make doc-check` prüft Links/Anker/Struktur,
keine Grammatik).
**Verifizierbar:** nein.
**Klasse:** `doku-tippfehler`

### F-5 · LOW · `DC-FA-VCS-001.a` Schritt 2 / Maintainability (Testabdeckung) · `internal/adapter/driven/git/git_test.go`

**Befund:** Von den drei im Plan (§2 Schritt 3) und in der Commit-Botschaft
genannten Kontrollfällen ist nur einer — `--staged`, unverändert seit slice-053
über `TestStaged`/`TestRangeAndFileAt` (reines Löschen) — durch **bestehende**
automatisierte Tests am echten Adapter gedeckt. Der dritte Kontrollfall
(„Rename **mit** Umformulierung, `--range`") ist ausschließlich über die
manuelle End-zu-Ende-Probe der Commit-Botschaft belegt; kein neuer Test in
`git_test.go` fährt ihn. Risiko ist niedrig, weil `DiffTreeOptions{DetectRenames:
false}` keinen eigenen Ähnlichkeits-Schwellenwert-Zweig hat, den nur dieser Fall
träfe (dieselbe Abschaltung deckt beide Fälle mechanisch) — aber ein künftiger
Wechsel auf einen similarity-basierten Ansatz mit eigener Schwelle würde nur
gegen den bereits vorhandenen Byte-identisch-Test (`TestRangePureRenameYieldsDelete`)
laufen, nicht gegen den Umformulierungs-Fall.
**Verifizierbar:** ja — Testdatei-Diff zeigt nur einen neuen Testfall
(`TestRangePureRenameYieldsDelete`), keinen für den Umformulierungs-Fall.
**Klasse:** `kontrollfall-nur-manuell-belegt`

---

## Negativbefunde (geprüft, ohne Befund)

- **Prüffrage 1 (Silent-Grün) — der eigentliche Gegenstand, mehrfach unabhängig
  gemessen.** Fünf End-zu-Ende-Läufe gegen `d-check:latest` in einer
  eigenständig aufgebauten Scratch-Repo-Kopie (reines Löschen, reiner Rename,
  Rename mit Umformulierung, Rename aus der Klasse heraus, `--staged`-Rename) —
  alle fünf melden `core-drift-vcs` auf dem alten Pfad, keine bleibt bei 0
  Befunden stehen. Die Bestandsmessung über 61 Commits (`361ed39..HEAD`, eine
  echte Range statt einer künstlichen) liefert dieselben 638 Dateien / 0 Befunde
  wie die Commit-Botschaft (dort 60 Commits — Off-by-one gegenüber meiner
  Zählung, ohne Belang für die Aussage: der Bereich enthält so oder so keinen
  ADR-Rename).
- **Prüffrage 2 (Kern-Modul-Korrektheit).** `CheckVCS`/`vcsDeleted`/`vcsModified`
  (`internal/hexagon/core/rules/vcs.go`) sind in diesem Diff **unverändert** —
  der Fix liegt ausschließlich in der Adapter-Übersetzung (`diffTrees`); die
  Kern-Logik, die Delete/Modified/Added interpretiert, war bereits korrekt und
  bleibt es. Der neue Core-Test (`TestVCSRenameOutOfClass`) zeigt zusätzlich
  End-zu-Ende (eigene Probe oben, vierter Fall), dass ein Rename aus der
  geschützten Klasse heraus genau einen Befund auf dem **alten** Pfad ergibt.
- **Prüffrage 3 (Hexagon-Import-Richtung, ADR-0005).** `make gates` schließt
  `arch-check` grün ein; der einzige go-git-Import bleibt in
  `internal/adapter/driven/git` (unverändert gegenüber slice-053, keine neue
  Import-Kante in diesem Diff).
- **Prüffrage 4 (Suppression ohne ADR).** `git diff 8117758~2..HEAD | grep
  nolint` liefert nichts; `make lint` läuft grün als Teil von `make gates`.
- **Prüffrage 5 (Netzzugriff außerhalb `external`).** Der Fix ändert nur die
  lokale Baum-Diff-Option (`DiffTreeOptions.DetectRenames`); kein neuer
  Netz-Pfad, `DC-QA-03` bleibt unberührt.
- **Prüffrage 10 (Messmethode vs. Spec-Stelle).** `spec/spezifikation.md`
  §DC-FA-VCS-001.a Schritt 2 benennt den Mechanismus jetzt (`DetectRenames:
  false`) und stimmt mit dem Code überein — durch die eigenen End-zu-Ende-Läufe
  oben unabhängig verifiziert, nicht nur gelesen.
- **`AGENTS.md` §3.5 — die explizit verlangte Prüfung.** Die Datei ist in diesem
  Diff nicht angefasst. Das ist **keine** Lücke: §3.5s Wortlaut („Maschinell
  erzwungen über `make adr-check` … pre-commit-Hook + PR-/Push-CI") macht keine
  modus- oder rename-spezifische Zusage, die eine Ergänzung bräuchte — es war
  **vor** diesem Slice bereits unwahr in der Praxis (der `--range`-Pfad
  enforced nicht, was der Satz versprach), und der Fix stellt die Wahrheit des
  unveränderten Wortlauts wieder her, statt eine neue Grenze zu benötigen. Die
  eigenen End-zu-Ende-Proben (reines Löschen, reiner Rename, Umformulierung,
  Klassen-Austritt, `--staged`) finden keinen verbleibenden Fall, in dem
  `make adr-check` über einen der beiden Modi eine Umbenennung/Löschung einer
  `Accepted`-ADR still ließe — es bleibt also nichts zu benennen. Damit steht
  auch fest: die DoD-Formulierung „oder trägt die verbleibende Grenze benannt"
  ist mit „keine Änderung" korrekt eingelöst, nicht übersehen.
- **MR-013 (Lifecycle-Move-Commit).** `8117758` ist ein reiner `git mv`
  (Rename-Score 100 %, 0 Zeilenänderungen an der Slice-Datei) plus der
  gekoppelten Roadmap-Änderung (Ruhe-Marker „Nichts in Arbeit" entfernt); die
  Inhaltsänderungen (Kommentar-Korrekturen, Fix, Test, Spezifikation) liegen
  in den beiden Folge-Commits — konform zur Ausnahme in `AGENTS.md` §3.3.
  WIP-Limit gewahrt: `in-progress/` trug vor dem Move keinen Slice.
- **`CHANGELOG.md`.** Nicht im Diff — konform zur expliziten DoD-Vorgabe „Kein
  CHANGELOG-Eintrag im Feature-Commit" und zur repo-weiten Release-Prep-Regel.
- **`docs/plan/planning/observations.md`.** Nicht im Diff — konform: Der
  Register-Eintrag entsteht bei der Closure, der Slice ist noch
  `in-progress/`. Die vom Plan zitierten Zähler (`BEO-023`=7, `BEO-024`=1,
  `BEO-012`=11) stimmen mit dem aktuellen Registerstand überein — keine
  vorgezogene oder verspätete Fortschreibung.
- **ADR-0024 Fitness Function.** *„`make adr-check` läuft rot bei einem
  Körper-Edit / Status-Rückfall an einer `Accepted`-ADR (bzw. Löschung/
  Umbenennung), grün sonst"* — die ADR selbst schränkt keinen Eingabe-Modus
  ein; der Fix schließt exakt die Lücke, die diese seit slice-053 stehende
  Zusage bisher nicht einlöste. ADR-0024 ist in diesem Diff nicht angefasst,
  was korrekt ist (§3.5-Immutabilität; kein inhaltlicher Widerspruch
  entstanden, der einen Folge-ADR bräuchte).
- **Referenz-Richtung / Beleg-Anker.** Die beiden `d-check:cite`-Direktiven in
  §7 des Slice-Plans zeigen auf `modul-05-planning-harness.md:213-214` bzw.
  `:219` und tragen dort das wörtliche Zitat — stimmt gegen die Baseline-Datei
  (durch `make doc-check` `citations`-Modul ohnehin scharf geprüft, hier
  zusätzlich gegen den Baseline-Text gegengelesen).
- **Prüffrage 9 (Quelle über Geltungsbereich hinaus zitiert).** Der Plan zitiert
  `DC-FA-VCS-001` korrekt als die Anforderung ohne Modus-Einschränkung (Wortlaut
  gegen `spec/lastenheft.md:2026-2034` gegengelesen — „gelöschte oder umbenannte
  immutable Datei → Befund `core-drift-vcs`" trägt keine Modus-Qualifikation);
  `ADR-0024`s Fitness Function wird korrekt als bereits bestehende, jetzt erst
  eingelöste Zusage gelesen, nicht als neue Entscheidung ausgegeben.
- **Test-Hygiene der neuen Fälle.** Beide neuen Tests
  (`TestRangePureRenameYieldsDelete`, `TestVCSRenameOutOfClass`) nutzen
  ausschließlich bestehende Helfer (`repoAt`/`put`/`snapshot`/`statusOf` bzw.
  `fakeVCS`/`adr`/`refs`) ohne Duplikation; beide Kommentare benennen die
  Grenze, die der jeweilige Test **nicht** deckt (keine
  Inhalts-Ähnlichkeits-Messung bzw. Klassen-Filter-Fehlrichtung) — das ist die
  sanktionierte Grenze-Klasse, kein Finding.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 2 | F-1, F-2 |
| LOW | 2 | F-4, F-5 |
| INFO | 1 | F-3 |

---

## Verdikt

**Bedingt mergebar — zwei MEDIUM (F-1, F-2) vor Merge zu klären, kein HIGH.**

Der eigentliche Gegenstand — der still-grüne `--range`-Pfad — ist geschlossen
und über den Bestand der Commit-Behauptungen hinaus **unabhängig** verifiziert:
fünf End-zu-Ende-Proben (Löschen, reiner Rename, Rename mit Umformulierung,
Rename aus der Klasse heraus, `--staged`-Rename) melden alle `core-drift-vcs`
auf dem korrekten (alten) Pfad, keine bleibt still. Die Kern-Modul-Logik ist
unverändert und war bereits korrekt; der Fix liegt sauber in der
Adapter-Übersetzung, ohne Hexagon-Grenzverletzung, ohne neue Suppression, ohne
Netz-Ausweitung. `AGENTS.md` §3.5 wurde zu Recht nicht angefasst — die
Prüfung ergibt, dass der unveränderte Wortlaut nach dem Fix wieder vollständig
trägt, keine benannte Grenze fehlt.

Die zwei MEDIUM sind beide aus derselben Wurzel wie der behobene Bug — ein
Kommentar/eine Botschaft behauptet mehr Gewissheit, als sie trägt —, nur eine
Ebene über dem Code: **F-1** trägt eine Chronik-Aussage statt einer
Zusage/Grenze im `.githooks/pre-commit`-Kommentar (Herkunfts-Prosa nach
`AGENTS.md` §3.7). **F-2** ist die interessantere: der Slice korrigiert „die
zwei falschen Zusagen" mit einer Vollständigkeits-Behauptung, die selbst nicht
gemessen wurde — eine dritte, wortgleich falsche Stelle
(`internal/hexagon/port/driven/vcs.go`) blieb unentdeckt und ist heute nur
durch den Code-Fix zufällig wieder wahr, nicht durch Prüfung. Beide sind mit
einer Zeile bzw. einem `grep`-Fund behebbar und blockieren nicht die
Korrektheit des Fixes selbst. F-3 (INFO) und F-4/F-5 (LOW) sind
nicht-blockierend.
