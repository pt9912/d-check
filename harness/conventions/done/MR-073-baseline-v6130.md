# MR-073 — Baseline-Pin-Hebung auf `v6.13.0` (fünfzehnter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(Pin-Fortschreibung im Bundle-Layout des
  Vorgängers; keine inhaltliche Abweichung)*
- **Datum:** 2026-09-29
- **Geltungsbereich:** §Baseline, pin-gebundene Verweise,
  `.harness/baseline/v6.13.0/`
- **Adaption:** Der vendorte Bestand steht auf
  [`v6.13.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.13.0).
  Drei Zwischen-Tags (`v6.10.0`, `v6.11.0`, `v6.12.0`) liegen zwischen dem
  alten Pin und diesem — ein Bump zielt auf den jeweils aktuellen Tag, nicht
  auf jeden dazwischen (Präzedenz: `v5.7.0`→`v5.9.0`,
  `v6.6.0`→`v6.9.0`, u. a.).

  **Der Delta, gemessen und nicht geschätzt** (`diff -r -I '<!-- Quelle:'`
  gegen den Vorzustand, beide Bäume noch vorhanden): 55 Dateien vorher wie
  nachher, keine neu, keine entfallen. **28** tragen ein inhaltliches Delta
  (19 Regelwerk-Dateien, 8 Templates, `SHA256SUMS`). Größte Bewegungen:
  `modul-06-roadmap.md` mit 67 Diff-Zeilen, `templates/AGENTS.template.md`
  mit 54, `templates/.d-check.yml` mit 49, `modul-04-adrs.md` mit 49.

  **Vier inhaltliche Schlagzeilen, keine davon übernommen** (Abgrenzung 1,
  slice-240): Erstens trägt `modul-06-roadmap.md` ein **Trigger-Audit der
  Welle** — vier Artefaktklassen (Carveout, Hard Rule, ADR,
  Bootstrap-aware Gate) werden bei der Welle-Closure auf ihren
  Auflösungs-Trigger geprüft. Zweitens hält `modul-04-adrs.md` fest, dass
  **eine Gate-Erweiterung nicht automatisch ein ADR-Anlass** ist — die
  Aufnahme eines unabhängig lauffähigen Wächters in `make gates` genügt
  einem Verweis auf seine tragende ADR. Drittens verschärft
  `templates/AGENTS.template.md` das Schnitt-Prinzip für
  Dokumentations-Regeln (Kurzzeile im Briefing, Volltext in
  `harness/rules/<name>.md`). Viertens dokumentiert `templates/.d-check.yml`
  die Linkpflicht für ADR-Kennungen (`ids` + `link-policy` statt nackter
  `token`-Klasse in `matrix`). Ob und wie d-check davon etwas adoptiert, ist
  Sache eigener Folge-Slices.

  **Spiegel-Klassen, gemessen am echten Vorzustand** (HEAD vor dem Swap):
  75 Dateien außerhalb des vendorten Baums nannten `v6.9.0` — **161**
  Vorkommen. **36** blieb der erste Ersetzungs-Wurf ganz schuldig (Pfad-
  Verweise, `d-check:cite`-Direktiven, die Skills unter `.harness/skills/`,
  `spec/architecture.md`); **10** trugen Mischfundstellen — lebende Verweise
  neben frozen Vergangenheits-Aussagen — und wurden zeilenweise getrennt:
  `AGENTS.md`, `harness/README.md`, `harness/conventions.md`, `roadmap.md`,
  `observations/README.md`, `MR-021`, `MR-049`, `MR-053` (je Release-/Zip-URL
  oder bare Nennung „Baseline-Form `v6.9.0` §N" retargetet, Prosa bleibt),
  `MR-056` (Baseline-Link des Zitats retargetet, der Wortlaut-Vermerk bleibt)
  und `spec/spezifikation.md` (§Rolle-Pointer retargetet, Historie-Zeilen
  bleiben) — die letzten beiden erst im Review-R1-Report als Mischfälle
  erkannt, nicht in der ersten Liste;
  **27** blieben ganz frozen: 13 `done/`-Slices (inkl. Stub slice-224), 2
  Review-Reports, 3 aufgelöste MR-Dateien, CO-002, der eingehende CR, 1
  Evidence-Datei, [ADR-0085](../../docs/plan/adr/0085-bedingte-pflicht-marke-open-tasks.md)
  (Provenanz-Prosa + Link, siehe unten —
  [`docs/plan/adr/0085-bedingte-pflicht-marke-open-tasks.md`](../../docs/plan/adr/0085-bedingte-pflicht-marke-open-tasks.md)), 1
  Spec-Historie-Zeilen-Träger (`spec/lastenheft.md`), 3 Code-Kommentar-Träger,
  `.d-check.closure.yml`
  (2 Provenanz-Kommentare). Dazu der
  slice-240-Plan selbst (Gegenstands-Nennung) und die MR-072-Datei (Zug nach
  `conventions/done/`, byte-stabil — ihr Geltungsbereich nennt den
  `v6.9.0`-Baum als Vergangenheits-Aussage korrekt). Die 8 Baseline-Symlinks
  unter `.claude/rules/` (gemessen; der Plan sagte 7) lösen gegen `v6.13.0`.

  **`ignore-refs` wächst diesmal** — anders als beim Vorgänger, und das ist
  gemessen, nicht angenommen: **vier** eingefrorene Artefakte tragen den
  entfernten `v6.9.0`-Baum als Markdown-**Link** ([ADR-0085](../../docs/plan/adr/0085-bedingte-pflicht-marke-open-tasks.md)
  in seiner Geltungsbegründung, die drei aufgelösten MR-Dateien
  MR-059/061/062 je in ihrer Begründung) und melden nach dem Entfernen
  `target-missing` — vier neue Einträge in `.d-check.yml`
  ([`MR-069`](../conventions.md#mr-069--das-ignore-refs-ventil-ist-eine-deklarierte-gate-senkung-und-es-wächst-mit-jedem-bump)).

  **`d-check:cite`-Neu-Ankern, kein Zitat-Delta** — die Abwesenheits-Behauptung
  der ersten Fassung dieses Eintrags war falsch (Review-R1-F-1, gemessen am
  Diff statt behauptet), und die erste Korrektur-Zahl (15/6/9) trug nicht
  diff-genau (Verifier V-1): gemessen **14** Direktiven, **7** davon neu
  geankert — reine Zeilenverschiebung bei unverändertem Wortlaut (`MR-043`
  `grundlagen-durchsetzungsschicht.md` 100-102→116-118; `MR-005` ebenda
  50→66; `MR-031` `modul-09-implementierung.md` 196→203; `MR-049`
  `modul-05-planning-harness.md` 170-171→180-181 und 160→170; die beiden
  Spannen des slice-240-Plans 363-364→373-374 und 369→379) — und **7**
  tag-only (Spanne identisch, nur der Baum-Pfad wechselte). Kein Zitat-Delta
  ([`MR-039`](../conventions.md#mr-039)): die wörtlichen Zitate in lebenden
  Dokumenten (MR-056 zu `modul-05` §Lifecycle) sind im `v6.13.0`-Wortlaut
  unverändert — der Satz „…die Bedingung dafür, dass die Datei überhaupt
  nach `done/` darf…" steht zeilenidentisch (Zeile 36, im Delta nicht
  berührt).
- **Begründung:** Der Nachtlauf meldete vier neuere Tags (`make
  baseline-freshness`, Exit 3, 2026-09-29) bei inhaltlich unverändertem
  Vorgänger — die Hebung ist eine reine Fortschreibung; die inhaltlichen
  Deltas (oben) warten auf eigene Folge-Slices.
- **Auflösungs-Trigger:** die nächste Pin-Hebung.
