# Review-Report: slice-219 — 2026-09-17 (R1)

**Review-Art:** Release-/Carveout-Auflösungs-Review — kein Code-Slice.
Geprüft wird der gesamte Commit-Bereich seiner Beanspruchung: Behauptungen im
Slice-Plan/Carveout gegen die Realität (gezogenes Image, GitHub-Release),
`CO-001`s interne Konsistenz, `carveouts/README.md`, die beiden Guide-Docs
`harness/sensors/adr-check.md`/`trace-check.md` gegen den tatsächlichen
Code-Stand, sowie CHANGELOG/version.md/README/Handbuch auf `v0.76.1`.

**Gegenstand:** Commit-Range `64206d53~1..HEAD` (sechs Commits:
`64206d53` Beanspruchung, `8c97cff3` Vorprüfungen/DoD-Präzisierung,
`39870bc9` Release-Prep, `528b55c4` CO-001 inhaltlich aufgelöst,
`3e0f621d` reiner `git mv`, `bc6ee3f6` Link-Tiefen + Sensor-Docs
nachgezogen), **plus** der unstaged/uncommittete Arbeitsbaum-Stand von
`docs/plan/planning/in-progress/slice-219-release-loest-co-001.md` (DoD-Haken,
Risiko-Ausgänge, Closure-Notiz §7) und der neuen, unversionierten Beobachtung
`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/`.

**Skill:** `.harness/skills/reviewer.md` @ v1.16.0 / 2026-09-07 ·
**Modell:** `claude-sonnet-5` · **Datum:** 2026-09-17

**Eingangs-Kontext:**

- Slice-Plan `docs/plan/planning/in-progress/slice-219-release-loest-co-001.md`
  (Arbeitsbaum-Stand)
- `CO-001` (`docs/plan/carveouts/done/CO-001-vcs-range-stiller-skip.md`)
- Vorgänger-Slice `docs/plan/planning/done/slice-220-vcs-pfadmenge-statt-diff.md`
  samt seiner zwei Review-Reports (R1/R2)
- `docs/user/releasing.md`
- `harness/sensors/adr-check.md`, `harness/sensors/trace-check.md`
- `internal/hexagon/core/rules/vcs.go`, `internal/adapter/driven/git/git.go`
  (aktueller HEAD-Stand)
- `docs/plan/carveouts/README.md`
- `DC-FA-VCS-001`/`DC-FA-VCS-001.a`, `DC-FA-DIST-001`/`DC-FA-DIST-002`
  (`spec/lastenheft.md`, `spec/spezifikation.md`)
- `AGENTS.md` §3 (Hard Rules), §5 (CHANGELOG/README/Handbuch-Disziplin), §6
  (Workflow)
- `ADR-0014` (`:latest`-Tag), `ADR-0064`/`ADR-0065` (Docker-Hub-Spiegel,
  fail-closed, Config-Digest-Gleichheit)

---

## Verifikationen dieses Laufs

Selbst ausgeführt, unabhängig vom Slice-Plan-Text:

- **GHCR-Digest gezogen:** `docker pull
  ghcr.io/pt9912/d-check@sha256:1470ecdcaa686a5ef4513dee9b0ae522586f54b87d568b06fc6b5b2741b633b3`
  — erfolgreich, Digest bestätigt sich selbst.
- **OCI-Label geprüft:** `org.opencontainers.image.version` des gezogenen
  Images = `0.76.1` — deckt sich mit dem Tag.
- **`gh release view v0.76.1`:** Release existiert, `prerelease: false`,
  Notizen nennen denselben Digest
  (`sha256:1470ecdcaa686a5ef4513dee9b0ae522586f54b87d568b06fc6b5b2741b633b3`)
  wortgleich zum Slice-Plan/CO-001-Beleg.
- **Docker-Hub-Spiegel-Gleichheit selbst nachvollzogen** (nicht nur die
  Pipeline-Aussage übernommen): `docker manifest inspect` gegen
  `ghcr.io/pt9912/d-check:v0.76.1` und `pt9912/d-check:v0.76.1` liefert
  denselben **Config-Digest** (`sha256:9213c6aa4200…`) — die Prüfgröße, die
  ADR-0065 selbst verwendet (Layer-Digests unterscheiden sich zwischen den
  Registries, das ist erwartet und irrelevant für den Vertrag).
- **DoD (2) unabhängig gemessen, nicht übernommen:** vier eigene Probe-Repos
  gebaut (git init, ADR-Datei mit `**Status:** Accepted`, echte
  Core-Änderung), je eine der vier Ausprägungen aus `CO-001`s Tabelle
  nachgestellt (unlesbares BASE-Blob per entferntem Loose-Objekt; unlesbarer
  BASE-Tree mit Pendant auf der HEAD-Seite; unlesbarer BASE-Tree ohne
  Pendant, Verzeichnis komplett gelöscht; unlesbarer HEAD-Tree), plus ein
  fünftes Kontroll-Repo mit einer **echten**, nicht verdeckten
  Core-Änderung. Alle fünf gegen das **gezogene** `v0.76.1`-Image gefahren
  (`docker run --rm --network none -v <probe>:/repo:ro <digest> --enable vcs
  --disable links --range <base>..<head>`):

  | Probe | Ergebnis |
  |---|---|
  | unlesbares BASE-Blob | `nicht lesbares Objekt zu "docs/plan/adr/0001-test.md" in "<base>": file not found`, **Exit 2** |
  | unlesbarer BASE-Tree, mit Pendant | `Range-Basis "<base>" nicht auflösbar: … nicht lesbarer Unterbaum "docs/plan/adr": object not found`, **Exit 2** |
  | unlesbarer BASE-Tree, ohne Pendant | dieselbe Fehlerform, **Exit 2** |
  | unlesbarer HEAD-Tree | `Range-Spitze "<head>" nicht auflösbar: … nicht lesbarer Unterbaum "docs/plan/adr": object not found`, **Exit 2** |
  | Kontrolle (echte Änderung, nichts verdeckt) | `1 Befund(e)`, `core-drift-vcs`, **Exit 1** — exakt die von `CO-001` behauptete kanonische Form |

  Alle vier Defekt-Ausprägungen brechen unabhängig bestätigt mit Exit 2 ab;
  die Kontrollprobe zeigt, dass das Modul bei einer echten, unverdeckten
  Verletzung weiterhin korrekt `Exit 1`/`1 Befund` meldet statt
  überzureagieren. Deckt sich mit `CO-001`s Auflösungs-Trigger-Tabelle und
  mit dem Slice-Plan-DoD (2)-Text.
- **`make gates`:** grün — `baseline-verify + workflow-pins + doc-check +
  lint + test + arch-check + coverage-gate + semgrep + gate-consistency +
  planning-check`; `coverage-gate` meldet 94.70 % (Schwelle 93 %); die
  beiden abschließenden `doc-check`-Läufe (Module `targets`/`planning`)
  melden je `802 Datei(en) geprüft, 0 Befund(e)`. Zehn Gates, wie in den
  Commit-Botschaften behauptet.
- **`make ci`/`image-test` nicht erneut lokal gefahren** — stattdessen das
  tatsächlich publizierte Artefakt direkt geprüft (Digest-Pull, OCI-Label,
  Docker-Hub-Config-Digest-Vergleich, vier funktionale Probe-Läufe). Das ist
  eine stärkere Prüfung des DoD-(1)-Anspruchs „am gezogenen Image
  verifiziert" als ein erneuter lokaler Image-Build. Negativbefund unten.

## Findings

Keine.

## Negativbefunde (Pflicht)

- **Faktentreue der Kern-Behauptungen** (Punkt 1/2 des Auftrags): geprüft —
  Digest, OCI-Label, GitHub-Release und alle vier `CO-001`-Ausprägungen sind
  unabhängig bestätigt, nicht nur aus dem Slice-Text übernommen. Kein
  Befund.
- **`CO-001`s interne Verweise** (Punkt 3): `Folge-Slice:`-Feld nennt
  `slice-220` (Klassen-Fix) und danach `slice-219` (Release) mit korrekter
  Begründung, warum nicht `slice-218`; `Geltungsbereich` wurde von „aktiv,
  zwei ungleiche Hälften" korrekt auf „aufgelöst" umgeschrieben, samt
  ehrlicher Selbstkorrektur zweier früherer Fassungen; `Auflösungs-Trigger`
  nennt den tatsächlichen Digest und die tatsächlichen Fehlermeldungs-Werte
  (durch meine eigenen Probes bestätigt, nicht nur zitiert); `Geschichte`-
  Tabelle trägt Anlage- und Auflösungszeile mit Verweis auf `slice-220`. Die
  Link-Tiefen wurden nach dem `git mv` (`3e0f621d`) korrekt um eine Ebene
  vertieft (`bc6ee3f6`) — `make gates`/`doc-check` bestätigt das mit 0
  Befunden über den gesamten Baum. Kein Befund.
- **`carveouts/README.md`** (Punkt 4): „keine aktiven Carveouts mehr, zwei
  aufgelöst" korrekt; Tabelle „Aktive Carveouts" leer, „Aufgelöste
  Carveouts" führt `CO-001` (Zeiger auf `done/`) und `CO-002` mit
  korrektem `aufgelöst durch`-Feld. Kein Befund.
- **`harness/sensors/adr-check.md`/`trace-check.md` gegen den Code**
  (Punkt 5): beide Dateien beschreiben nach `bc6ee3f6` exakt das aktuelle
  Verhalten von `internal/hexagon/core/rules/vcs.go` und
  `internal/adapter/driven/git/git.go` — vier Ausprägungen, alle auf
  denselben `Exit 2`-Codepfad zusammengeführt (bestätigt durch meine eigenen
  Probes oben), `CO-001` korrekt als „inzwischen aufgelöst" verlinkt
  (`carveouts/done/…`). Die genannten Testnamen
  (`TestFileAtUnlesbaresObjekt`, `TestFileAtEintragOhneBlob`) existieren
  tatsächlich in `internal/adapter/driven/git/git_test.go`; die vorher
  genannten, inzwischen entfernten Namen
  (`TestVCSAddedMeldetUnlesbareBasis`) sind korrekt nicht mehr referenziert.
  Kein Befund.
- **CHANGELOG/version.md/README/Handbuch auf `v0.76.1`** (Punkt 6):
  `version.md` §Aktuell + Verlauf-Zeile korrekt, `<a id>`-Anker korrekt von
  `v0.76.0` auf `v0.76.1` verschoben (die `v0.76.0`-Zeile verliert ihn, wie
  `releasing.md` §Release-Prep verlangt). Beide READMEs (`docker run …
  :v0.76.1`) und alle ~26 Pin-Stellen im Handbuch (`grep` bestätigt: kein
  `v0.76.0`-Pin außerhalb der bewusst historischen Erwähnungen) sind
  gehoben; Handbuch-Kopf (`1.73`/`v0.76.1`), §11-Verlaufszeile
  (chronologisch unten angehängt, nicht oben) und der `vcs`-Modul-Abschnitt
  (die Ausprägungs-Tabelle korrekt auf vier statt zwei Zeilen erweitert)
  sind konsistent aktualisiert. `CHANGELOG.md` trägt den neuen
  `[0.76.1]`-Abschnitt korrekt **unter** seiner Versionsnummer (kein
  `[Unreleased]`), im Release-Prep-Commit (nicht im Feature-Commit) — nach
  `AGENTS.md` §5. Verbleibende `v0.76.0`-Nennungen (`.d-check.yml`-Kommentar,
  beide Sensor-Docs, Handbuch-Warnkasten, Handbuch-§11-Historie) sind
  durchweg bewusste historische Referenzen („bis `v0.76.0` galt …"), keine
  vergessenen Pins. **Patch- statt Minor-Bump korrekt begründet**: reiner
  Bugfix am internen Mechanismus, keine neue Config-/CLI-Fläche, kein
  Lastenheft-Bump, `spec/spezifikation.md`-Historie bestätigt „kein
  Lastenheft-Bump — dieselbe Zusage, ein anderer Mechanismus". Kein Befund.
- **`make gates`** (Punkt 7): selbst gefahren, grün, zehn Gates — deckt sich
  mit allen sechs Commit-Botschaften der Range. Kein Befund.
- **Botschaft/Beleg verallgemeinert über die gemessene Menge hinaus**
  (Anker 8): der Slice-Plan-DoD (2) und `CO-001`s Auflösungs-Trigger
  behaupten „alle vier Ausprägungen … gemessen" — durch meine eigenen,
  unabhängig konstruierten Probes bestätigt, keine überdehnte
  Schlussfolgerung („damit ist jede Diff-Abhängigkeit im Repo behoben" o. ä.)
  wird an keiner Stelle gezogen; im Gegenteil benennt `harness/sensors/
  adr-check.md` explizit die verbleibende Grenze (Gitlink-Ersatz einer
  immutablen Datei bleibt unbemerkt) sowie, dass nur drei Module (`vcs`,
  `commits`, `tracked`) untersucht wurden. Kein Befund.
- **Zählung „41 Verzeichnisse" im Beobachtungs-Register** (Anker 17,
  vorgelagerte Sichtung §8): gegen den Commit-Stand `8c97cff3` nachgezählt
  (`git ls-tree`) — exakt 41 `observation.md`-Dateien zum Zeitpunkt der
  Vorprüfung. Dass jetzt (nach der später in `bc6ee3f6` neu angelegten
  Beobachtung) 42 Verzeichnisse existieren, ist kein Widerspruch, sondern
  der erwartete zeitliche Ablauf — die Sichtung lief vor der Neuanlage. Kein
  Befund.
- **`d-check:cite`-Direktiven in slice-219 §8** (Anker 9, MR-054): beide
  Zitate (`modul-05-planning-harness.md:363-364` und `:369-369`) wortgleich
  gegen die Baseline-Datei nachgeprüft (`grep -n`/`sed -n`) — treffen exakt.
  Kein Befund.
- **Grenzen-Liste ohne größte Lücke** (Anker 18): `CO-001`s eigener Abschnitt
  „Rest-Risiko-Umfang" benennt explizit die verbleibenden Bedingungen (a)/(b)
  für den Regelfall und (c) für die dritte Ausprägung — keine verschwiegene
  Kategorie gefunden. `harness/sensors/adr-check.md` benennt zusätzlich den
  Gitlink-Ersatz-Fall als bewusst bestehende, unveränderte Grenze. Kein
  Befund.
- **Neue Beobachtung `guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-
  stehen`**: `observation.md`/`state.md`/`evidence/slice-219.md` folgen der
  Ziel-Form (unveränderliche Identität, veränderlicher Stand `offen — 1×`,
  eine Evidence-Datei je Vorgang); Sub-Area `*` korrekt (kein spezifischeres
  Modus-deklariertes Verzeichnis berührt); Kennungs-Pfad `BEO-ALL/…` korrekt
  aus der Modus-Deklaration abgeleitet. Der Rückverweis in
  `evidence/slice-219.md` auf `docs/plan/planning/done/slice-219-…md` (statt
  `in-progress/`) ist eine bewusste Vorwegnahme des Zielpfads nach dieser
  Slice-Closure — als Inline-Code, nicht als Markdown-Link, löst er keinen
  Link-Sensor aus und wird mit dem `git mv` dieser Slice wahr. Kein Befund.
- **`Verantwortlich:`/`Autor:`-Felder, `Berührte Spec-Stellen:`,
  Sub-Area-§8-Block**: Feldform korrekt (`§Ziel-Form: Slice`); zwei
  Sub-Areas benannt (`*`, `docs/plan/carveouts/`), beide korrekt als GF ohne
  eigene Modus-Deklaration eingeordnet; Nachtlauf-Stand (`MR-053`) gelesen
  und mit Datum/Zeitstempel belegt. Kein Befund.
- **DoD-(1)-Anspruch „am gezogenen Image verifiziert (Digest, OCI-Label,
  Smoke)"**: Digest und OCI-Label von mir selbst bestätigt (siehe oben); den
  in `releasing.md` genannten Smoke-Test (`make run`) und den vollen
  `make ci`/`image-test`-Lauf habe ich **nicht** erneut lokal gefahren —
  stattdessen das reale publizierte Artefakt über vier funktionale Probes
  direkt getestet, was den Anspruch aus Sicht dieses Reviews eher über- als
  unterdeckt. Benannt, kein Blocker: `make ci` lief bereits in der
  Release-Pipeline vor dem Tag (durch das erfolgreiche, nicht-draft
  GitHub-Release indirekt belegt).

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** keine.

## Verdikt

**Freigegeben.** Alle sieben im Auftrag benannten Prüfpunkte sind unabhängig
verifiziert, nicht nur aus dem Slice-Plan-Text übernommen: der GHCR-Digest
`sha256:1470ecdcaa686a5ef4513dee9b0ae522586f54b87d568b06fc6b5b2741b633b3`
existiert, trägt das OCI-Label `0.76.1` und ist im nicht-draft
GitHub-Release `v0.76.1` genannt; der Docker-Hub-Spiegel trägt denselben
Config-Digest; alle vier in `CO-001` behaupteten Ausprägungen des
`vcs`-Defekts brechen gegen vier selbst gebauten Probe-Repos und das
tatsächlich gezogene Image mit Exit 2 ab, während eine Kontrollprobe mit
echter Verletzung weiterhin korrekt mit Exit 1/1 Befund meldet; `CO-001`s
Verweise, `carveouts/README.md`, die beiden Guide-Sensor-Docs (gegen den
aktuellen `vcs.go`/`git.go`-Code) und die Release-Doku (CHANGELOG,
version.md, beide READMEs, Handbuch) sind alle konsistent und korrekt auf
`v0.76.1` gehoben. `make gates` bestätigt zehn grüne Gates. Aus Review-Sicht
steht der Closure von slice-219 nichts entgegen; das Abhaken des
DoD-Punkts „Unabhängiger Review durchgeführt" und der Verweis auf diesen
Report sind Sache der anschließenden Verifikation/Closure, nicht dieses
Reviews.
