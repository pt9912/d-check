# Review-Report: slice-133 — Baseline-Sensor eingesteckt — 2026-08-23

**Review-Art:** Code-Review (Diff gegen Slice-Plan, Hard Rules und die
zitierten `DC-*`/`MR-*`-Verträge, Modul 10 §Drei Review-Arten) · **Gegenstand:**
Commit `55b8815` — der Feature-Commit von slice-133 (`feat(harness):
Baseline-Sensor eingesteckt — verify in gates, freshness als Nachtlauf
(MR-011, DC-QA-03, slice-133)`). Die im Auftrag genannte Range
`f9a2fd0..55b8815` löst im Repo nicht auf; verwendet wurde der tatsächliche
Eltern-Commit `55b8815^..55b8815` (= `4b98bd6..55b8815`, ein reiner
Lifecycle-Move ohne Inhaltsänderung an den hier geprüften Dateien).

**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 (Stand `d87c4d1`) ·
**Modell-ID:** `claude-sonnet-5` · **Datum:** 2026-08-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-133-baseline-sensor-verdrahten.md`
  (§1–§9, insb. §2 Vorgehen, §3 NICHT, §4 DoD, §5 Risiken)
- `docs/plan/planning/welle-84-durchsetzung.md`
- `tools/harness/fetch-baseline-cache.sh` (vollständig, insb. `verify()`,
  `check_latest()`, Dispatch)
- `Makefile` (Targets `baseline-verify`/`baseline-freshness`, `gates`-Kette,
  `.PHONY`, `.NOTPARALLEL`)
- `.github/workflows/upstream-drift.yml` und `.github/workflows/ci.yml`
  (Vergleich), `.github/workflows/release.yml` (Pin-Vergleich)
- `AGENTS.md` §3.1, §3.9, §4, §5 (Botschafts-Regel/`BEO-009`)
- `harness/README.md` §Sensors
- `spec/lastenheft.md` → `DC-QA-03` (Volltext)
- `.d-check.yml` → Block `targets`
- `tools/harness/record-gates.sh`, `tools/harness/working-tree-hash.sh`,
  `tools/semgrep.sh`, `a-check.mk`, `.githooks/commit-msg`,
  `.githooks/pre-commit` (Gegenprobe der Host-Werkzeug-Aussage)
- `harness/conventions.md` §MR-011/§MR-021, Baseline
  `.harness/baseline/v5.11.0/regelwerk/modul-13-quality-gates.md`
  §Fitness Function

**Vom Reviewer selbst gefahren** (nur Lesekommandos, `git show`/`git diff`,
sowie **isolierte Reproduktionen** außerhalb des Repos): das Skript
`fetch-baseline-cache.sh` wurde unverändert in ein separates, temporäres
Git-Repo unter der Scratchpad-Verzeichnis kopiert (`/tmp/…/scratchpad/…`,
inzwischen entfernt) und dort mit `bash`/`make` gegen synthetische
Baseline-Fixtures gefahren, um Exit-Codes und Fehlerpolitik empirisch zu
verifizieren — **nicht** im geprüften Repo, keine `make`-Targets dort
ausgeführt, keine Datei im Repo geändert (`git status` vor/nach: clean).

**Verdikt: blockierend** — ein HIGH, vier MEDIUM, ein LOW, kein INFO.

---

## Findings

### F-1 — `baseline-verify`s Manifest-Deckung scannt nur zwei Unterbäume; eine Datei direkt unter `.harness/baseline/<tag>/` ist unsichtbar und das Gate bleibt grün

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §4 („Halluzinierte Gates sind die häufigste Form
  von Harness-Lüge") · Reviewer-Skill §HIGH-Anker „Stilles-Grün-Pfad in
  einem Gate oder Gate-Skript" · `AGENTS.md` §3.8 (ein Modul verspricht nur
  über das, was es scannt — die Folge kann still sein)
- **pfad:** `tools/harness/fetch-baseline-cache.sh:118` (`on_disk="$(find
  "${baseline}/regelwerk" "${baseline}/templates" -type f …)"`); Zusagen-Text
  in `AGENTS.md:248`, `harness/README.md:81`, `Makefile:159-163`
- **befund:** Sowohl `sha256sum -c` als auch die Manifest-Deckung
  beschränken sich auf die zwei benannten Unterbäume
  `${baseline}/regelwerk` und `${baseline}/templates`. Eine Datei, die
  direkt unter `.harness/baseline/<tag>/` (Geschwister von `regelwerk/` und
  `templates/`, z. B. neben `SHA256SUMS`) eingelegt wird, taucht in keinem
  `find`-Aufruf auf — weder im Manifest noch im Platten-Zähler — und
  `make baseline-verify` bleibt Exit 0. Empirisch reproduziert (isolierte
  Kopie): eine Datei `<tag>/rogue-outside.txt` neben `regelwerk/`/`templates/`
  eingelegt liefert „verify ok (2 Dateien, vollständig)", Exit 0 — obwohl der
  Bestand unter `.harness/baseline/<tag>/` nachweislich verändert wurde. Der
  Zusagen-Text in `AGENTS.md`/`harness/README.md`/`Makefile` spricht vom
  „committeten"/„vendorten Baseline-Bestand" bzw. „ist der committete Bestand
  unversehrt?" — nicht eingeschränkt auf die zwei Unterbäume — und ist damit
  breiter als das tatsächlich Geprüfte. Die DoD-geforderte zweite
  Brech-Richtung („Datei zusätzlich eingelegt") wurde nur **innerhalb** der
  beiden Bäume gemessen (Commit-Botschaft: „Manifest (51 Zeilen) != Dateien
  auf Platte (52)"), nicht an dieser dritten Form.
- **verifizierbar:** ja — `echo x > .harness/baseline/<pin>/rogue.txt &&
  make baseline-verify` (Exit-Code danach wieder auf 0 prüfen, Datei danach
  löschen).
- **klasse:** manifest-deckung-scan-luecke-ausserhalb-baum

### F-2 — `DC-QA-03` wird für ein bloßes Host-Skript ohne jede Container-/Netz-Sandbox zitiert, obwohl seine Messmethode einen `docker run --network none`-Lauf beschreibt

- **kategorie:** MEDIUM
- **quelle:** `spec/lastenheft.md#dc-qa-03` (Anforderung spricht von „Das
  Tool"; Messmethode: „Integrationstest mit read-only-Mount und
  netzwerkloser Umgebung (`docker run --network none`)")
- **pfad:** `Makefile:164-166` (Kommentar „NICHT in gates (DC-QA-03: make
  gates bleibt netzlos)"); `harness/README.md:79` (Quelle-Spalte „DC-QA-03
  (netzlos bleibt gates, nicht dieses Target)"); Slice-Plan §2 Schritt 1
- **befund:** `DC-QA-03`s Anforderungstext bindet an „das Tool" (den
  d-check-Prozess) und seine Messmethode ist wörtlich ein containerisierter
  `--network none`-Lauf. `baseline-verify`/`baseline-freshness` sind reine
  Host-Bash-Aufrufe ohne jede Container- oder Netzwerk-Namespace-Isolation —
  ihre Netzlosigkeit (bzw. Netz-Nutzung) ist eine Eigenschaft des
  Skript-Codes, nicht eine gemessene Sandbox-Eigenschaft. Das bricht mit dem
  bisherigen Repo-Präzedenzfall: `record-gates.sh`/`working-tree-hash.sh`
  (ebenfalls reine Host-Skripte, seit vor diesem Slice in `gates`) tragen in
  keiner der beiden Doku-Tabellen eine `DC-QA-03`-Zitation (`Quelle`-Spalte
  dort: „—"). Dieser Commit ist der erste, der `DC-QA-03` einem Host-Skript
  ohne jede Container-Messung zuschreibt — eine Reichweite, die die zitierte
  Anforderung selbst nicht deckt.
- **verifizierbar:** teilweise — kein Gate prüft Zitations-Reichweite; der
  Textvergleich (`DC-QA-03`-Wortlaut gegen die zitierende Zeile) ist
  Text-gegen-Text nachvollziehbar.
- **klasse:** requirement-zitat-ausserhalb-messmethode

### F-3 — Die Doku behauptet, `make baseline-freshness` selbst liefere Exit 3/Exit 4; GNU Make normalisiert jeden fehlschlagenden Recipe-Exit auf 2

- **kategorie:** MEDIUM
- **quelle:** AGENTS.md §5/Botschafts-Regel (Aussage muss der gemessenen
  Menge entsprechen)
- **pfad:** `.github/workflows/upstream-drift.yml:11-16` (Kopf-Kommentar
  „WAS ER MELDET (`make baseline-freshness`, Exit = schlimmster Teil): …
  Exit 3 … Exit 4 …"); `harness/README.md:79` („Exit 3" / „Exit 4" im
  `make baseline-freshness`-Zeilentext)
- **befund:** Empirisch reproduziert (isolierte Kopie, `make -f <mk> foo`
  mit `exit 1`/`exit 3`/`exit 4` im Recipe): GNU Make meldet bei jedem
  fehlschlagenden Recipe „Fehler N" im Log, der **eigene** Exit-Code des
  `make`-Prozesses ist in allen drei Fällen **2**, nie der ursprüngliche
  Wert. `make baseline-freshness` selbst kann also nie mit Exit 3 oder Exit
  4 enden — nur das direkt aufgerufene Skript
  (`fetch-baseline-cache.sh --check-latest`) tut das. Der rote Job in der
  Actions-Übersicht ist davon unberührt (jeder Nicht-Null-Exit lässt den
  Schritt fehlschlagen), aber die spezifische Exit-Code-Zuordnung in den
  beiden genannten Doku-Stellen ist für das benannte Kommando falsch.
- **verifizierbar:** ja — Recipe mit `exit 3`/`exit 4` in einem
  Test-Makefile, `make <target>; echo $?`.
- **klasse:** exit-code-attribution-durch-make-schicht-verschluckt

### F-4 — Commit-Botschaft schreibt `record-gates.sh`/`working-tree-hash.sh` die Aufrufe `sed` und `grep` zu; beide Skripte rufen keines von beiden

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §3.1 (Host-Werkzeug-Liste) · AGENTS.md §5
  (Botschaft darf nicht mehr behaupten als die Arbeit trägt)
- **pfad:** `tools/harness/record-gates.sh` (ruft nur `bash`, `mkdir`,
  delegiert an `working-tree-hash.sh`); `tools/harness/working-tree-hash.sh`
  (ruft `git`, `sort`, `sha256sum`, `awk`, `readlink` — **kein** `sed`,
  **kein** `grep`); Commit-Botschaft, Absatz „MITKORRIGIERT" (Behauptung:
  „record-gates.sh und working-tree-hash.sh … rufen sha256sum, awk, sed,
  grep und sort"); `AGENTS.md:86-88`
- **befund:** Gegen den tatsächlichen Skript-Text geprüft, rufen die beiden
  namentlich genannten Skripte weder `sed` noch `grep` auf — nur `sha256sum`,
  `awk`, `sort` (plus `git`/`readlink`). `grep` wird zwar an anderer Stelle
  im `gates`-Baum gebraucht (`tools/semgrep.sh:70`,
  `fetch-baseline-cache.sh`s `verify()`-Pfad), aber `sed` wird von **keinem**
  Host-Skript aufgerufen, das über `make gates` hängt — die einzigen
  `sed`-Aufrufe im Repo liegen in `fetch-baseline-cache.sh`s
  `check_latest()` (Ziel `baseline-freshness`, bewusst **nicht** in `gates`),
  in `tools/bench-fixture.sh` (Ziel `bench`, in `fullbuild`, nicht `gates`)
  und im `Makefile`-Ziel `versions` (nicht in `gates`). Die neue
  `AGENTS.md`-§3.1-Liste (`sha256sum, find, awk, sed, grep, sort`) nennt
  damit ein Werkzeug (`sed`), das kein an `make gates` hängendes Skript
  braucht — bei einer Korrektur, deren eigene Begründung ausdrücklich mit
  „Das war schon vor diesem Slice falsch" antritt.
- **verifizierbar:** ja — Volltext der drei genannten Skripte gegen die
  Commit-Aussage; `grep -rn 'sed' tools/ Makefile` gegen die Ziel-Zugehörigkeit
  zu `gates`.
- **klasse:** werkzeug-zuschreibung-ohne-fundstelle

### F-5 — `curl`-Aufrufe in `check_latest()` haben kein Timeout; ein hängender (nicht fehlschlagender) Netzpfad kann den fail-open-Nachtlauf in einen roten Job statt eines SKIP verwandeln

- **kategorie:** MEDIUM
- **quelle:** Slice-Plan §3 („Kein Fail-closed im Netz-Teil … ein Sensor,
  der bei Netzstörung rot wird, wird abgeschaltet")
- **pfad:** `tools/harness/fetch-baseline-cache.sh:138`,
  `tools/harness/fetch-baseline-cache.sh:159` (`curl -fsSL` ohne
  `--max-time`/`--connect-timeout`); `.github/workflows/upstream-drift.yml:36`
  (`timeout-minutes: 10`)
- **befund:** Beide `curl`-Aufrufe im Netz-Pfad tragen keine
  Zeitbegrenzung. Ein klar **fehlschlagender** Request (Verbindung
  abgelehnt, DNS-Fehler, HTTP-Fehlercode) wird von `curl -f` sofort
  quittiert und korrekt als `SKIP` behandelt (empirisch reproduziert: `curl`
  durch ein Fake-Binary mit sofortigem `exit 6` ersetzt → beide Teile
  `SKIP`, Exit 0). Ein **stockender** TCP-Pfad (Handshake hängt, Response
  bleibt aus, ohne dass die Verbindung terminiert) lässt `curl` dagegen ohne
  eigenes Timeout unbegrenzt warten; das Skript selbst käme dann nie in den
  SKIP-Zweig, sondern würde erst vom Job-`timeout-minutes: 10` extern
  beendet — ein anderes Ergebnis (roter/abgebrochener Job) als das im
  Slice-Plan §3 zugesagte „Ausfall ⇒ SKIP je Teil".
- **verifizierbar:** ja — `curl` gegen einen Endpunkt, der TCP annimmt, aber
  nie antwortet (z. B. `nc -l` ohne Response), zeigt das Fehlen des SKIP-Pfads.
- **klasse:** fail-open-lücke-fehlendes-netz-timeout

### F-6 — `AGENTS.md` §3.9 behauptet weiterhin „beide Workflows"; der Commit legt mit `upstream-drift.yml` den dritten an

- **kategorie:** LOW
- **quelle:** `AGENTS.md:228` selbst (unverändert von diesem Commit)
- **pfad:** `AGENTS.md:226-228` („Jeder `uses:`-Eintrag … Das gilt für
  **beide** Workflows gleich"); neu: `.github/workflows/upstream-drift.yml`
- **befund:** Vor diesem Commit gab es zwei Workflow-Dateien (`ci.yml`,
  `release.yml`), „beide" war akkurat (im Vorgänger-Review zu slice-131
  ausdrücklich als Negativbefund festgehalten: „es gibt keine dritte
  Workflow-Datei"). Dieser Commit legt `upstream-drift.yml` als dritte
  Workflow-Datei an — die Pin-Substanz stimmt bei allen drei (SHA-Pin mit
  Tag-Kommentar, gegengeprüft), aber der Prosa-Text „beide Workflows" in
  §3.9 wurde nicht auf drei nachgezogen.
- **verifizierbar:** ja — `ls .github/workflows/` gegen den Wortlaut in
  `AGENTS.md:228`.
- **klasse:** doku-prosa-zaehlwort-veraltet

## Negativbefunde

- geprüft, ohne Befund: Netzlosigkeit des `--verify`-Dispatch-Pfads —
  der gemeinsame Vorspann vor der `case`-Weiche (Tag-Ermittlung via
  lokalem `grep` auf `harness/conventions.md`, Funktionsdefinitionen) enthält
  keinen `curl`/`wget`/Netzaufruf; `verify()` selbst nutzt ausschließlich
  `sha256sum`, `find`, `grep -c` lokal
- geprüft, ohne Befund: Fail-closed-Mechanik von `--verify` **via
  `make baseline-verify`** — empirisch reproduziert: geänderte Datei
  innerhalb `regelwerk/`/`templates/` ⇒ Exit 2, zusätzlich eingelegte Datei
  innerhalb der beiden Bäume ⇒ Exit 2, Rückbau ⇒ Exit 0 — deckungsgleich mit
  der Commit-Botschaft, sofern über `make` (nicht das nackte Skript, das
  ohne die make-Schicht Exit 1 liefert) gemessen
- geprüft, ohne Befund: Dispatch-Wrapper für `--check-latest`
  (`set +e; check_latest; rc=$?; set -e; exit "$rc"`) — der Exit-Code wird
  korrekt nicht verschluckt, kein `|| true` am Dispatch selbst
- geprüft, ohne Befund: Prioritätslogik `rc` in `check_latest()`
  (4>3>0) — gelesen, für alle Kombinationen von `currency`/`authenticity`
  korrekt, kein Bug
- geprüft, ohne Befund: leeres `SHA256SUMS` — `sha256sum -c` verweigert
  selbst die Ausführung („keine korrekt formatierte Prüfsummenzeile
  gefunden"), Exit 1 unter `make` → Exit 2; kein stilles Grün
- geprüft, ohne Befund: gelöschte vendorte Datei — `sha256sum -c` meldet
  „FEHLSCHLAG öffnen oder lesen", korrekt rot
- geprüft, ohne Befund: nicht existierendes Tag-Verzeichnis — `[ -f "$sums"
  ]` schlägt vor jeder inhaltlichen Prüfung fehl, kein stilles Grün
- geprüft, ohne Befund: `Makefile` — `baseline-verify` ist tatsächlich das
  erste Element von `gates:` (Zeile 176); `.PHONY` trägt beide neuen Namen;
  `.NOTPARALLEL` (unverändert) hält `record-gates` weiterhin als letzten
  Schritt
- geprüft, ohne Befund: `AGENTS.md` §4 ↔ `harness/README.md` §Sensors —
  beide neuen Zeilenpaare (`baseline-verify`/`baseline-freshness`) sind
  inhaltlich deckungsgleich (Netzlos/fail-closed vs. Netz/fail-open,
  Gate-Zugehörigkeit); keine `BEO-010`-Drift zwischen den beiden Tabellen
- geprüft, ohne Befund: `.d-check.yml` → `targets:`-Block — unverändert,
  deckt weiterhin beide Doku-Tabellen (`AGENTS.md`, `harness/README.md`) ab;
  kein Nachzug nötig
- geprüft, ohne Befund: `upstream-drift.yml` — `permissions: {}` auf
  Top-Level, `contents: read` auf Job-Ebene, Trigger nur `schedule`/
  `workflow_dispatch` (kein `push`/`pull_request`, keine Doppelläufe mit
  `ci.yml`), `actions/checkout`-Pin identisch zu `ci.yml`/`release.yml`
  (SHA `de0fac2e…`, Tag-Kommentar `v6.0.2`) — §3.9 in der Substanz erfüllt
- geprüft, ohne Befund: `make baseline-freshness` hat im Makefile keine
  Docker-Vorbedingung (kein `build`-Prerequisite) — läuft auf
  `ubuntu-latest` ohne Image-Bau; `curl`/`unzip`/`sha256sum`/`find`/`git`
  sind Standard-Ausstattung gehosteter Ubuntu-Runner
- geprüft, ohne Befund: `--check-latest` bei vollständigem, sofortigem
  Netzausfall (curl schlägt sofort fehl) — empirisch reproduziert, beide
  Teile `SKIP`, Exit 0 (Kontrast zu F-5, das den hängenden statt den
  fehlschlagenden Fall betrifft)
- geprüft, ohne Befund: `modul-13-quality-gates.md` §Fitness Function —
  der Commit-Verweis „beide Richtungen gemessen (modul-13 §Fitness
  Function)" erweitert die dortige Ein-Beispiel-Formel um eine zweite
  Brech-Richtung; das ist eine plausible Anwendung der Methode auf ein Gate
  mit zwei Fehlerklassen, keine dem Modul widersprechende Zitation

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 4 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** `manifest-deckung-scan-luecke-ausserhalb-baum`
· `requirement-zitat-ausserhalb-messmethode` ·
`exit-code-attribution-durch-make-schicht-verschluckt` ·
`werkzeug-zuschreibung-ohne-fundstelle` ·
`fail-open-lücke-fehlendes-netz-timeout` · `doku-prosa-zaehlwort-veraltet`

## Verdikt

**Merge-blockierend:** ja — ein HIGH, vier MEDIUM offen. Die Kern-Prämisse
des Slice (Sensor eingesteckt statt neu gebaut, zwei Targets mit zwei
Fehlerpolitiken, Doku nachgezogen) trägt und ist an mehreren Stellen
nachweislich korrekt umgesetzt (Negativbefunde). F-1 ist der schwerste
Befund: die neu gegatete Integritätsprüfung hat einen stillen Grün-Pfad
genau dort, wo die Commit-Botschaft selbst die Latte legt („ohne die
Manifest-Deckung wäre 'prüft die Integrität' überdehnt") — die Deckung ist
auch mit ihr noch überdehnt, nur eine Ebene höher (Baum-Geschwister statt
Baum-Inhalt). F-2/F-3/F-4 sind Reichweiten- bzw. Zuschreibungs-Lücken in der
Doku- und Commit-Evidenz (BEO-011/BEO-009-Nachbarschaft), die die
Gate-Funktion selbst nicht brechen, aber die Zusagen überziehen, gegen die
sie zitiert werden. F-5 ist eine unadressierte Lücke in der fail-open-Zusage
für einen realistischen (nicht nur hypothetischen) Netzzustand. F-6 ist ein
reiner Doku-Nachzug.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser
Report selbst ist ein Lauf-Beleg — DoD-/Spec-Konformität prüft der Verifier
separat.
