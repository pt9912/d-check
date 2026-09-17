# Review-Report: slice-226 — 2026-09-17 (R1)

**Review-Art:** Code-Review — geprüft wird die Implementierung gegen den
Slice-Plan, den eingehenden CR, `DC-FA-VCS-001`/`DC-FA-VCS-001.a`, `ADR-0086`
(und dessen zitierten Präzedenzfall `ADR-0072`) sowie `AGENTS.md` §3.

**Gegenstand:** Commit-Range `7a6ee6f9^..37103614` (vier Commits:
`7a6ee6f9` Beanspruchung, `2209c919` Fix + Tests, `4e7c3818` Spec-Nachtrag +
CR-Antwort, `37103614` DoD-Häkchen). **Hinweis zur Auftrags-Range:** die im
Auftrag genannte Range `bc6ee3f6..HEAD` enthält zusätzlich zwei
`slice-219`-Commits (`54a187cb`, `05f8e00c`, CO-001-Auflösung/Release
`v0.76.1`), die bereits unter
`docs/reviews/2026-09-17-slice-219-release-loest-co-001-review-r1.md`
laufen und **nicht** Gegenstand dieses Reports sind — geprüft wurde exakt der
`slice-226`-eigene Diff (11 Dateien, 751 Zeilen, `git diff --stat
7a6ee6f9^..HEAD`).

**Skill:** `.harness/skills/reviewer.md` @ v1.16.0 / 2026-09-07 ·
**Modell:** `claude-sonnet-5` · **Datum:** 2026-09-17

**Eingangs-Kontext:**

- Slice-Plan `docs/plan/planning/in-progress/slice-226-vcs-pack-alias-fremdes-praefix.md`
- Eingehender CR `docs/plan/cr/2026-09-17-cr-eingehend-ai-harness-init-pack-praefix.md`
- `DC-FA-VCS-001`/`DC-FA-VCS-001.a` (`spec/lastenheft.md`, `spec/spezifikation.md`)
- `ADR-0086` (Accepted; Schärft-Ziel `DC-FA-VCS-001.a` Schritt 2) und sein
  zitierter Präzedenzfall `ADR-0072` (R3, `gopkg.in/yaml.v3`-Kapsel-Erweiterung
  Adapter → Liste)
- `AGENTS.md` §3 (Hard Rules), §6 (Workflow), §3.8 (Modul-Scan-Grenze)

---

## Verifikationen dieses Laufs

Alle folgenden Läufe wurden **selbst ausgeführt** (Docker/`make`-only, kein
Host-Go; `git` als normaler Host-Tool für die Probe-Repos), nicht nur aus der
Commit-Botschaft übernommen:

- `make test`: grün, alle Pakete inkl. `internal/adapter/driven/git`.
- `make gates`: **zweimal** grün gefahren (zweiter Lauf mit
  `--no-cache-filter` auf `lint`/`test`, also nicht aus dem Cache
  durchgewunken) — `baseline-verify + workflow-pins + doc-check + lint +
  test + arch-check + coverage-gate + semgrep + gate-consistency +
  planning-check green`. `make coverage-gate`: **94.70 %** (Schwelle 93 %).
  `make arch-check` (über `ghcr.io/pt9912/a-check:v0.19.0`): `gesamt: 0
  Befund(e)` — die `os`/`io/fs`-Kapsel-Erweiterung bleibt grün. `make
  semgrep`: `0 findings`.
- **Bewusstes Brechen (Modul 11) für DoD (1):** `git worktree add` auf den
  Parent-Commit `bc6ee3f6` (vor der Beanspruchung, also ohne den Fix), die
  neue `git_test.go` unverändert hineinkopiert, `docker build --target test`
  gefahren. Ergebnis: `TestAllPathsPackUnterFremdemPraefix` **schlägt fehl**
  mit genau dem CR-Fehlerbild (`Range-Basis "…" nicht auflösbar: reference
  not found`), während alle anderen Pakete grün bleiben. Der Test
  diskriminiert also echt zwischen vorher/nachher — keine der beiden neuen
  Testfunktionen ist tautologisch grün. `TestAllPathsPackMitUnbrauchbaremPraefixBleibtFehlerhaft`
  war in diesem Lauf nicht die fehlschlagende Funktion — konsistent damit,
  dass sie als Gegenprobe ohnehin **vor und nach** dem Fix bestehen soll, was
  gegenprüft wurde: kompilierbar und grün gegen den Parent-Commit ebenso wie
  gegen `HEAD`.
- **`repackToPack`-Helfer geprüft (nicht vertraut):** Die Behauptung des DoD
  (1) — „gemessen an einem Repo, dessen Objekte nur noch im umbenannten Pack
  liegen (lose Kopien entfernt)“ — wird durch den Helfer tatsächlich
  eingelöst: er ruft `RepackObjects` und entfernt danach jedes lose
  Objektverzeichnis (`objects/<2-hex>/`) außer `pack`/`info`. Ohne diesen
  zweiten Schritt wäre der Test tautologisch grün gewesen (Objekt über die
  lose Kopie lesbar, Pack-Alias irrelevant) — genau die in der Aufgabe
  benannte Tautologie-Gefahr. Sie tritt hier **nicht** ein.
- **Eigene End-to-End-Probe, unabhängig von der mitgelieferten Test-Suite**
  (echter `git`-Binary + gebautes `d-check:latest`-Image, `--enable vcs`,
  Config mit `vcs.paths`/`immutable-when` auf ein Muster wie
  `docs/plan/adr/[0-9]*.md`):
  - Repo mit zwei Commits, `git repack -a -d`, alle `pack-*` auf `loose-*`
    umbenannt (exakter CR-Repro-Schritt): meldet den echten Core-Drift-Befund
    (`core-drift-vcs`, Exit 1) — **nicht** mehr Exit 2. Vor dem Fix (Parent
    `bc6ee3f6`, separat gebaut) wäre das derselbe stille/abbrechende Zustand,
    den der CR beschreibt.
  - Dasselbe Repo, Pack+Idx+Rev stattdessen auf `xpack-garbage.{pack,idx,rev}`
    umbenannt (kein gültiges Hash-Suffix): `error: Range-Basis "…" nicht
    auflösbar: reference not found`, **Exit 2** — der fail-closed-Abbruch aus
    der Nicht-Bitte des CR bleibt exakt erhalten.
- **go-git-Quelltext gelesen** (v5.19.2, aus dem `deps`-Layer via
  `docker run … cat`, nicht per Host-Go/Netz-Vermutung): `DotGit.objectPacks()`
  ruft `ReadDir("objects/pack")` genau **einmal** pro `ObjectStorage`-Instanz
  (`requireIndex()` cached über `s.index != nil`); ein `Adapter` hält genau
  eine `ObjectStorage` über die gesamte CLI-Laufzeit. `packAliases()` scannt
  das Verzeichnis also nicht pro Commit/Objekt neu, sondern höchstens einmal
  pro Prozesslauf zusätzlich zu einem `Open()`-Aufruf pro referenziertem Pack
  — kein DC-QA-01-Risiko bei der heute üblichen Pack-Anzahl. Multi-Pack-Index
  (`.midx`) und Reverse-Index (`.rev`) werden von dieser go-git-Version
  **nirgends** gelesen (`objectPacks()` betrachtet ausschließlich
  `pack-*.pack`-Namen) — die bewusste Nicht-Behandlung von `.rev`/`.midx` in
  `packAliasFS` ist damit nicht nur harmlos, sondern zutreffend.

## Findings

### R1-F-1 (MEDIUM) — Drei Doku-Stellen beschreiben weiterhin die vor `slice-226` geltende Pack-Sichtbarkeitsgrenze als aktuell

- **Kategorie:** MEDIUM
- **Quelle:** `AGENTS.md` §6 Schritt 7 (*„Doku/Indizes aktualisieren, falls
  ein öffentlicher Vertrag berührt"*)
- **Pfad:** `harness/sensors/adr-check.md:33-53` (§Grenze, Punkt 3),
  `harness/sensors/trace-check.md:20-37` (§Grenze, Punkt 3),
  `docs/user/benutzerhandbuch.md:1745-1792` (Modul `vcs`) und
  `docs/user/benutzerhandbuch.md:1808-1823` (Modul `commits`)
- **Befund:** Alle vier Stellen behaupten unverändert, ein Pack unter einem
  anderen Namen als `pack-<Hash>.{idx,pack}` sei „unsichtbar“ und breche den
  Lauf mit Exit 2 ab — `adr-check.md`/`trace-check.md` mit einer konkreten
  Sechs-Proben-Tabelle, in der `loose-<eigener Hash>` ⇒ Exit 2 gelistet ist;
  das Benutzerhandbuch nennt denselben Fall mit demselben Beispiel und dem
  Workaround `git repack -A -d`. Nach diesem Slice ist genau das für einen
  Pack mit gültigem Hash-Suffix **und** passender `.idx`-Datei nicht mehr
  wahr — dieser Fall wird jetzt aufgelöst und meldet den echten Befund statt
  Exit 2. `spec/spezifikation.md` wurde korrekt nachgezogen (DoD-Punkt
  erfüllt), diese vier Stellen — zwei Sensor-Dokumente und zwei
  Handbuch-Abschnitte, alle mit öffentlichem Charakter (Sensors-Doku wird von
  jedem Agentenlauf als Kanon gelesen, das Handbuch ist die
  Nutzer-Referenz) — nicht.
- **Verifizierbar:** ja — Textvergleich der vier Stellen gegen den in
  `spec/spezifikation.md:1801-1813` beschriebenen (und empirisch bestätigten,
  s. o.) neuen Mechanismus.
- **Klasse:** `guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`
  (dieselbe Registerklasse wie bei der `CO-001`-Auflösung in `slice-219` —
  `docs/plan/planning/observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/`,
  Stand dort aktuell 1× belegt. Dieser Befund wäre, sofern bei Closure
  eingetragen, die **zweite** Evidenz derselben Beobachtung — noch nicht die
  3×-Schwelle, aber ein zweiter Treffer derselben Klasse innerhalb von zwei
  aufeinanderfolgenden Slices).

### R1-F-2 (MEDIUM) — `packAliases()` löst eine Namenskollision nicht deterministisch auf

- **Kategorie:** MEDIUM
- **Quelle:** `DC-QA-02` (Determinismus: identische Eingabe ⇒ identische
  Ausgabe)
- **Pfad:** `internal/adapter/driven/git/packalias.go:47-76`
  (`packAliases()`)
- **Befund:** `packAliases()` iteriert mit `for name := range onDisk` — Gos
  Map-Iterationsreihenfolge ist bewusst randomisiert. Liegen zwei
  alias-fähige, nicht-kanonisch benannte Dateien mit **identischem**
  Hash-Suffix im selben `objects/pack`-Verzeichnis (z. B. zwei Kopien
  desselben Packs unter verschiedenen Fremdnamen, oder — als
  Angriffs-/Fehlerfall außerhalb des CR-Anlasses, wie im Slice-Plan §6 selbst
  benannt — zwei inhaltlich verschiedene Dateien mit zufällig identischem
  Namens-Suffix), gewinnt bei `aliases[canonicalPack] = name` /
  `aliases[canonicalIdx] = idxName` diejenige, die die Karte zuletzt sieht —
  nicht deterministisch über Prozessläufe hinweg. Schwerer wiegt: `ReadDir`
  und `Open` rufen `packAliases()` **unabhängig voneinander** neu auf (kein
  gemeinsamer Cache); bei einer solchen Kollision kann der beim `ReadDir`
  gewählte Gewinner ein anderer sein als der beim nachfolgenden
  `Open("pack-<Hash>.idx")` bzw. `Open("pack-<Hash>.pack")` — Idx und
  Pack-Bytes könnten dann aus zwei verschiedenen realen Dateien stammen.
  Praktisch dürfte dies meist über `loadIdxFile`s
  Checksummen-Vergleich (`idxf.PackfileChecksum != h`) fail-closed
  auffliegen, solange der Namens-Hash nicht zufällig mit dem echten
  Idx-Checksum übereinstimmt — das ist aber eine Folge, keine Garantie, und
  nirgends als Risiko benannt oder getestet (weder Slice-Plan §6 noch
  `git_test.go` konstruieren den Kollisionsfall).
- **Verifizierbar:** ja — zwei Dateien `alpha-<H>.{pack,idx}` und
  `beta-<H>.{pack,idx}` mit identischem `<H>`, aber unterschiedlichem Inhalt,
  in `objects/pack` anlegen und `packAliases()` über mehrere Prozessläufe
  beobachten (oder mit `GODEBUG`/eigenem Testtreiber die Karten-Iteration
  fixieren).
- **Klasse:** `pack-alias-map-iteration-nondeterminism`

## Negativbefunde

- **Trailing-Slash-Behandlung in `ReadDir`** (`packalias.go:96-98`): geprüft
  — `strings.TrimSuffix(p, "/") != packDir` normalisiert `"objects/pack/"`
  und `"objects/pack"` gleich, bevor verglichen wird. Kein Befund.
- **`.midx`/`.rev`-Behandlung**: geprüft gegen den go-git-v5.19.2-Quelltext
  (s. o., Abschnitt Verifikationen) — diese go-git-Version liest weder
  Multi-Pack-Index noch Reverse-Index-Dateien; die Nicht-Behandlung in
  `packAliasFS` entspricht dem tatsächlichen Konsumenten-Verhalten, ist keine
  übersehene Lücke.
- **Performance/DC-QA-01**: geprüft — `packAliases()` scannt `objects/pack`
  höchstens einmal je `Open()`-Aufruf pro referenziertem Pack (go-gits
  `requireIndex()`-Cache greift pro `ObjectStorage`-Instanz, die über die
  gesamte CLI-Laufzeit eines `Adapter` besteht), nicht pro Commit oder
  Objekt in einer Range. Bei der heute üblichen Pack-Anzahl (einstellig)
  kein beobachtbares Risiko für den `DC-FA-VCS-001.a`/`DC-QA-01`-Pfad. Kein
  Befund — siehe aber R1-F-2 zur (unabhängigen) Kollisions-Frage.
- **`os`/`io/fs`-Kapsel-Erweiterung auf Scope-Creep geprüft**
  (`.a-check.yml`, `internal/adapter/driven/git/packalias.go`): `os`
  erscheint in `packalias.go` ausschließlich als `os.FileInfo`-Typreferenz
  (Interface-Signatur, Embedding in `aliasFileInfo`) — kein `os.Open`,
  `os.ReadFile` oder anderer Host-I/O-Aufruf; `io/fs` wird in diesem Diff gar
  nicht importiert (die Kapsel-Erweiterung dafür ist vorsorglich, symmetrisch
  zu `os`, aber ungenutzt). `make arch-check` bestätigt `0 Befund(e)`. Kein
  Befund.
- **Präzedenzfall-Behauptung von `ADR-0086` gegen `ADR-0072` geprüft**: ADR-
  0086 zitiert ADR-0072 als Präzedenz für „Kapsel-Erweiterung Adapter → Liste“
  (R4, `os`/`io/fs`). `ADR-0072` §Entscheidung Punkt 8 belegt tatsächlich
  genau diese Bewegung, nur für R3/`gopkg.in/yaml.v3` („Die yaml-Allowlist …
  wächst dafür von zwei auf drei Adapter“). Die Präzedenz-Behauptung trägt —
  sie ist nicht bloß behauptet, sondern am zitierten Text nachvollziehbar.
  Kein Befund.
- **Zwei DoD-Testbehauptungen gegen „bewusstes Brechen“ (Modul 11) geprüft**:
  s. o., Abschnitt Verifikationen — beide neuen Tests diskriminieren echt,
  keiner ist tautologisch grün. Kein Befund.
- **Abweichung „Plan vor Code“ (Slice-Plan §1) geprüft**: die Abweichung ist
  benannt (nicht verschwiegen), mit nachvollziehbarer Begründung (go-gits
  private Pack-Discovery-Mechanik war ohne Lektüre nicht seriös planbar) und
  einem konkreten Lerneintrag-Verweis auf §7. Am Ergebnis (ADR mit
  Alternativen-Vergleich, zwei diskriminierende Tests, konsistente
  Spec-/CR-Nachträge) ist keine Spur einer durch diese Reihenfolge
  verursachten Qualitätslücke sichtbar — die Prämisse „Recherche war die
  Plan-Arbeit“ hält der Überprüfung stand. Kein Befund.
- **CR-Format und Nachtrag-Konsistenz geprüft**: Kopf-Felder, `## Antwort`-
  Sektion und `<!-- d-check:status-provenance -->`-Marker der eingehenden
  CR-Datei entsprechen exakt dem durch die beiden Vorgänger-CRs
  (2026-09-06/07, „eingehend“-Klasse) etablierten Muster; der publizierte
  Digest `sha256:1470ecdc…` für `v0.76.1` im CR-Wortlaut stimmt mit dem
  lokal vorhandenen Image überein (`docker inspect`). Kein Befund.
- **`docs/plan/adr/README.md`-Indexeintrag für `ADR-0086`** geprüft: vorhanden,
  Spaltenform konsistent mit Nachbareinträgen. Kein Befund.
- **CHANGELOG.md/README/Handbuch-Kopf** geprüft: in den Feature-Commits
  korrekt unangetastet (Release-Prep-Grenze aus `AGENTS.md` §5 eingehalten)
  — unabhängig davon bleibt die in R1-F-1 gemeldete **bestehende** Prosa im
  Handbuch stehen, was ein anderer Befund ist als eine verfrühte neue
  Changelog-Zeile.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:**
`guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen` ·
`pack-alias-map-iteration-nondeterminism` ·
`stale-kommentar-technische-begruendung`

## LOW-Findings (nice-to-fix)

- **R1-F-3** — `internal/adapter/driven/git/packalias.go:130-136`
  (`splitPackPath`-Doc-Kommentar): *„ohne path/filepath, weil DotGit seine
  eigenen Pfade unabhängig vom Host-Trennzeichen immer mit '/'
  zusammensetzt“* — geprüft gegen den go-billy-v5.9.1-Quelltext:
  `ChrootHelper.Join`/`ChrootOS.Join` delegieren an `filepath.Join`, also
  **host-separator-abhängig**, nicht separatorunabhängig. Der Grund, warum
  `splitPackPath`s fest verdrahtetes `"/"` trotzdem sicher ist, ist nicht
  „DotGit ist trennzeichenunabhängig“, sondern dass dieses Produkt
  ausschließlich unter Linux läuft (Docker-only-Distribution, `AGENTS.md`
  §3.1/ADR-0002) — dort ist `filepath.Join` ohnehin `"/"`. Die im Kommentar
  gegebene technische Begründung ist damit unzutreffend, auch wenn die
  Code-Zeile selbst für den tatsächlichen Laufzeit-Kontext korrekt bleibt.
  Ein künftiger Leser, der daraus eine allgemeine Eigenschaft von
  go-git/go-billy ableitet, wird fehlgeleitet. Quelle: `AGENTS.md` §3.7
  (Kommentar trägt eine der fünf Klassen — hier eine falsche
  Abgrenzungs-Begründung). Verifizierbar: ja, Quelltext-Lesung
  `go-billy/v5@v5.9.1/helper/chroot/chroot.go:298-300` und
  `osfs/os_chroot.go:99-101`. Klasse: `stale-kommentar-technische-begruendung`.

## Verdikt

**Blockiert.** R1-F-1 (MEDIUM) betrifft öffentlich gelesene Dokumentation
(zwei Sensor-Guides, die jeder Agentenlauf als Kanon konsultiert, plus zwei
Abschnitte der Nutzer-Handbuchs) und sollte vor der nächsten Closure-Runde
auf den seit diesem Slice geltenden Mechanismus nachgezogen werden — dieselbe
Fehlerklasse, die `slice-219`s Closure bereits einmal für die
`CO-001`-Auflösung gemeldet hat
(`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`).
R1-F-2 (MEDIUM) ist ein am Code nachvollziehbarer, aber eng umrissener
Determinismus-Rand (setzt eine Namenskollision zwischen zwei alias-fähigen
Dateien voraus, die außerhalb bekannter `git`-Tooling-Praxis liegt) — vor der
Closure sollte er zumindest als Risiko in Slice-Plan §6 nachgetragen und
idealerweise mit einer deterministischen Tie-Break-Regel (z. B. lexikografisch
kleinster Name gewinnt, oder Abbruch bei Mehrdeutigkeit statt stillem
Auswählen) geschlossen werden. R1-F-3 (LOW) hält die Closure für sich
genommen nicht auf. Die eigentliche Fix-Logik, ihre beiden DoD-Tests, die
Architektur-Entscheidung (`ADR-0086`) und ihr Präzedenzfall-Bezug halten der
Überprüfung — einschließlich dreier unabhängiger empirischer Gegenproben mit
echtem `git`/gebautem Image — stand.
