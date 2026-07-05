# Review-Report — slice-064 (`gate-consistency.sh` Voll-Tombstone) — Plan-/Doc-first-Review R1

## Kopf-Metadaten

- **Datum:** 2026-07-05
- **Gegenstand:** Slice-Plan `docs/plan/planning/next/slice-064-gate-consistency-tombstone.md`
  (Voll-Tombstone von `tools/gate-consistency.sh`: die `DC-QA-03`-Modullisten-
  Restprüfung wandert in einen Go-Test im `configyaml`-Paket, das Skript wird
  entfernt, `make gate-consistency` fährt nur noch den `targets`-Modul-Dogfood)
  samt begleitender [ADR-0032](../plan/adr/0032-gate-consistency-tombstone.md)
  (Proposed) und ADR-Index-Zeile.
- **Reviewer-Rolle:** unabhängig/adversarial; nicht Mit-Autor des Plans.
- **Prüfmethode:** rein statisch — `grep`, `ls`, Lesen des Ist-Baums. **Keine**
  `make`-/`go`-/`docker`-Läufe (Verifier-getrennt). Der Plan liegt in `next/`;
  Implementierung hat **nicht** begonnen (kein `internal/`-Test, keine
  `.d-check.yml`-/Makefile-Änderung im Baum) — geprüft wird die Vollständigkeit
  und Konsistenz des Plans gegen den realen Repo-Stand.
- **Eingangs-Kontext:** Slice-Plan; ADR-0032; ADR-0031; ADR-0028; ADR-0025
  (`codepaths.ignore-refs`); `.d-check.yml`; `Makefile` (`gate-consistency`,
  `test`, `gates`); `harness/README.md` §Sensors; `AGENTS.md` §4;
  `spec/lastenheft.md` §DC-QA-03; `tools/gate-consistency.sh`;
  `internal/hexagon/core/app/diagnose_test.go` (Muster-Vorbild);
  `internal/adapter/driven/configyaml/`; Reviewer-Skill v1.2.0.
- **Baum-Scope (bestätigt):** `docs/plan/planning/next/slice-064-…md`,
  `docs/plan/adr/0032-…md` und die ADR-Index-Zeile sind neu; `tools/gate-consistency.sh`
  existiert noch; `slice-063` liegt in `done/`, `v0.38.0` ist in `CHANGELOG.md`
  und `slice-063` (done) als released markiert. Sauberer Doc-first-Schnitt.

## Findings

### F-1 — HIGH — Der `codepaths.ignore-refs`-Tombstone für `tools/gate-consistency.sh` fehlt im gesamten Plan → `make doc-check`/`make gates` läuft nach der Entfernung rot (`codepath-missing`)

- **kategorie:** HIGH
- **quelle:** Harness-Lüge (stilles-Grün-Erwartung im Gate-Pfad) / `DC-FA-CODE-001`
  (`codepaths`) / [ADR-0025](../plan/adr/0025-codepaths-ignore-refs.md) / DoD-Konformität
- **pfad:** Slice `docs/plan/planning/next/slice-064-gate-consistency-tombstone.md:38`–`:57`
  (§2 Entscheidungen + §3 DoD) und ADR-0032 `docs/plan/adr/0032-gate-consistency-tombstone.md:56`–`:74`
  (Konsequenzen); Gegen-Autorität `.d-check.yml:84`–`:90` (`codepaths.ignore-refs`
  mit fünf Einträgen)
- **befund:** Der Plan entfernt `tools/gate-consistency.sh` per `git rm`, nennt aber
  **an keiner Stelle** (§2 Entscheidungen, §3 DoD, ADR-0032 Konsequenzen) das
  `codepaths.ignore-refs`-Register. Genau dieses Register ist der etablierte,
  gate-erzwungene Mechanismus für die Entfernung eines zitierten Gate-Skripts: alle
  **fünf** zuvor abgelösten Skripte (`adr-immutable-check.sh`, `completeness-check.sh`,
  `trace-check.sh`, `planning-consistency.sh`, `arch-check.sh`) stehen in
  `.d-check.yml:90`, weil das `codepaths`-Modul jeden Inline-Code-Pfad mit einem
  konfigurierten Root-Präfix (`tools/…`) auf Existenz prüft. `tools/gate-consistency.sh`
  wird in Inline-Code über den gesamten **nicht** von `codepaths.exempt-paths`
  (`docs/reviews/**`) ausgenommenen Baum zitiert: u. a. `CHANGELOG.md:28`,
  `spec/lastenheft.md:1473`, `spec/spezifikation.md:1208`/`:1235`,
  `docs/plan/adr/0031-…md:20`/`:24`/`:125`, `docs/plan/adr/0032-…md:12`/`:20`/`:46`
  (ADR-0032 wird laut DoD **Accepted** → immutabel), `docs/plan/planning/in-progress/roadmap.md:17`/`:24`/`:324`,
  `docs/plan/planning/done/slice-063-…md:27`/`:195`, `…/slice-040-…md:49`,
  `…/slice-009-…md:33`, `…/slice-058-…md:50` — plus die Slice-Datei selbst
  (`:23`/`:38`/`:50`/`:63`), die bei Closure nach `done/` wandert (weiterhin
  gescannt). Beleg, dass dies real rot wird und nur das Register grün hält: die
  bereits entfernten Skripte werden **heute** in Inline-Code zitiert und sind grün —
  `CHANGELOG.md:90` (`tools/planning-consistency.sh`), `CHANGELOG.md:124`
  (`tools/trace-check.sh`), `done/slice-056-…md:7`, ADR-0027 (`:12`/`:105`/`:138`) —
  ausschließlich, weil sie in `.d-check.yml:90` stehen; ADR-0027 nennt das Muster
  kanonisch „Skript entfernt, **Tombstone statt Edit**". **Failure:** ohne den
  Eintrag `"tools/gate-consistency.sh"` in `codepaths.ignore-refs` meldet der
  nächste `make doc-check`/`make gates`-Lauf ~20 `codepath-missing`-Befunde und
  Exit 1 — der DoD-Punkt „`make gates` / `make ci` grün" (`:56`) ist nicht
  erreichbar, und die DoD-Aussage „`grep -r gate-consistency.sh` findet nur noch
  Historie (`done/`, `CHANGELOG`, ADR-Kontext)" (`:51`) verkennt gerade diese
  Referenzen: `done/`, `CHANGELOG` und die ADRs sind **nicht** aus dem
  `codepaths`-Scan ausgenommen (nur `docs/reviews/**` ist es), also keine harmlose
  „Historie".
- **verifizierbar:** ja — `make doc-check` (bzw. `make gates`) nach `git rm tools/gate-consistency.sh`
  **ohne** ignore-refs-Eintrag liefert die `codepath-missing`-Befunde; mit dem
  Eintrag schweigen sie. Der Widerspruch ist heute schon per `grep` an den fünf
  bestehenden Tombstone-Einträgen (`.d-check.yml:90`) gegen die grüne CI belegbar.

### F-2 — MEDIUM — Markdown-**Links** auf das Skript werden `target-missing`, und `codepaths.ignore-refs` deckt sie nicht (die `links`-Achse ist per ADR-0025 bewusste Rest-Falle) — inklusive eines Edits an einer `done/`-Slice

- **kategorie:** MEDIUM
- **quelle:** `DC-FA-LINK-001` (`links`) / [ADR-0025](../plan/adr/0025-codepaths-ignore-refs.md)
  (`links`-Achse als Rest-Falle) / DoD-Konformität
- **pfad:** `docs/plan/planning/done/slice-013-codepaths-modul.md:68`
  (Markdown-Link mit Zielpfad `tools/gate-consistency.sh`, vier `../`-Segmente),
  `harness/README.md:58` (Markdown-Link mit Zielpfad `tools/gate-consistency.sh`);
  Plan §2/§3 (`:38`–`:57`) nennt weder Link-Edit noch die Rest-Fallen-Grenze
- **befund:** Zwei Referenzen auf `tools/gate-consistency.sh` sind **Markdown-Links**,
  keine Inline-Code-Spans: die live `harness/README.md:58` und — kritischer —
  `done/slice-013-codepaths-modul.md:68` (eine als „Historie" eingestufte
  `done/`-Slice). Ein Link auf die gelöschte Datei ist ein `links`-Befund
  (`target-missing`), **nicht** `codepaths` — und `codepaths.ignore-refs` wirkt nur
  im `codepaths`-Modul; ADR-0025 hält die `links`-Achse ausdrücklich als „bewusste
  Rest-Falle" offen (kein Link-Tombstone-Ventil). Solche Links müssen daher
  **editiert** werden (Präzedenz: die verify-doc-refs.sh-Entfernung stellte ihren
  Geltungsbereich-Link in `MR-003` auf einen Code-Span um). Der Plan adressiert
  das nicht: der `harness/README.md:58`-Link fällt zwar plausibel unter den
  „Doku-Nachzug … §Sensors" (§2, `:42`), aber nur, wenn der Implementierer den
  **Link** (nicht bloß die Prosa) tilgt; der `done/slice-013:68`-Link ist gar nicht
  erwähnt und verlangt einen Inhalts-Edit an einem `done/`-Dokument. **Failure:**
  bleibt einer der beiden Links stehen, meldet `make doc-check` `target-missing`
  und Exit 1 — derselbe rote Gate-Ausgang wie F-1, aber über ein anderes Modul,
  das der ignore-refs-Eintrag aus F-1 **nicht** abfängt.
- **verifizierbar:** ja — `make doc-check` nach der Entfernung ohne Link-Edit
  liefert `target-missing` auf `done/slice-013:68` (und `harness/README.md:58`,
  falls der Link dort überlebt).

### F-3 — MEDIUM — Die `DC-QA-03`-Bindung wandert weg von `gate-consistency`, aber der Plan spezifiziert nur „Rest-Skript-Teil entfällt", nicht die Umschreibung der Sensors-/§4-Bindungsspalte

- **kategorie:** MEDIUM
- **quelle:** Harness-Ehrlichkeit (Gate behauptet eine Bindung, die es nicht mehr
  erfüllt) / `DC-QA-03` (DC-Bindung)
- **pfad:** `harness/README.md:58` (Sensors-Zeile `make gate-consistency`,
  Bindungsspalte „`DC-QA-03` (DC-Bindung)") und `:53` (`make test` = `DC-QA-02`);
  `AGENTS.md:140`; Plan §2 (`:42`–`:45`) / ADR-0032 Konsequenz „Doku-Nachzug"
  (`:68`–`:72`)
- **befund:** Nach der Verlagerung führt `make gate-consistency` **nur** noch den
  `targets`-Dogfood (Bindung `DC-FA-TGT-001`); die `DC-QA-03`-Modullisten-Integrität
  wird künftig im Go-Test unter `make test` geprüft (ADR-0032 `:61`). Die
  Sensors-Tabelle bindet die `gate-consistency`-Zeile heute aber genau an
  `DC-QA-03` (`harness/README.md:58`, rechte Spalte), und `make test` trägt dort
  nur `DC-QA-02` (`:53`). Der Plan/ADR beschreibt den Doku-Nachzug als „verlieren
  ihren ‚Rest-Skript'-Teil" — das trifft die **Beschreibung**, nicht die
  **Bindungsspalte**. **Failure:** wird nur die Prosa gekürzt, behauptet die
  Sensors-Zeile weiter eine `DC-QA-03`-Durchsetzung durch `gate-consistency`, die
  dort nicht mehr stattfindet, während der neue Bindepunkt (`make test`) sie nicht
  ausweist — genau die „Gate behauptet, was es nicht tut"-Drift, gegen die das
  Meta-Gate-Programm antritt. Dieselbe Umschreibung betrifft die `AGENTS.md:140`-Zelle
  („+ `DC-QA-03`-Modulliste (Rest-Skript)").
- **verifizierbar:** teilweise — kein Sensor prüft die Bindungsspalten-Korrektheit;
  per `grep`/Lesen der Sensors-Tabelle gegen den reduzierten Target-Recipe
  belegbar. (Dass kein Gate dies fängt, ist Teil des Befunds.)

### F-4 — MEDIUM — ADR-0031 steht noch auf `Proposed` (Datei + Index), obwohl slice-063 in `done/` und `v0.38.0` released ist; ADR-0032 „revidiert" eine Konsequenz dieser noch nicht angenommenen ADR

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §3.5 (ADR-Lifecycle `Proposed → Accepted`, Immutabilität)
  / Konsistenz Plan-Prämisse ↔ Repo-Stand
- **pfad:** `docs/plan/adr/0031-targets-deklarations-konsistenz-modul.md:3`
  (`**Status:** Proposed`), `docs/plan/adr/README.md:46` (Index-Zeile ADR-0031
  „Proposed"); vs. `docs/plan/planning/done/slice-063-targets-modul.md:3`
  („Status: done …, Release **v0.38.0**") und `CHANGELOG.md:7` (`## [0.38.0]`);
  ADR-0032 Bezug/Revidiert (`docs/plan/adr/0032-gate-consistency-tombstone.md:10`–`:15`),
  Prämisse (`:30`–`:33`)
- **befund:** ADR-0032 stützt sich auf zwei Behauptungen: (1) „`targets` ist seit
  **v0.38.0** released … bewährt" und (2) es „revidiert die `Scope = Kern`-Konsequenz
  von ADR-0031". (1) stimmt (`v0.38.0` ist in `CHANGELOG.md`, slice-063 in `done/`).
  Aber ADR-0031 trägt in Datei **und** Index unverändert `Proposed` — obwohl seine
  Slice abgeschlossen und das Modul released ist. Die Präzedenz (ADR-0027 Geschichte:
  „Angenommen mit der slice-056-Closure … Status Accepted") zeigt, dass die
  `Accepted`-Setzung ein Closure-Schritt der umsetzenden Slice ist; für slice-063
  scheint er zu fehlen. Damit „revidiert" eine `Proposed`-ADR (0032) eine Konsequenz
  einer anderen, ebenfalls `Proposed` gebliebenen ADR (0031) — die §3.5-Ceremony
  („neue ADR statt Überschreiben") setzt eine **Accepted** Vor-ADR voraus; eine
  `Proposed`-ADR wäre direkt fortschreibbar. **Failure:** wird slice-064 auf ADR-0031
  als „entschieden/immutabel" aufgesetzt, während sie real `Proposed` ist, entsteht
  eine ungeklärte Reihenfolge — entweder ist die ADR-0031-`Accepted`-Umschreibung
  ein offener slice-063-Closure-Defekt (dann vor slice-064 zu schließen), oder die
  „Revidiert"-Relation von ADR-0032 ist verfrüht formuliert. Der Plan macht die
  Abhängigkeit nicht sichtbar.
- **verifizierbar:** ja — `grep '**Status:**' docs/plan/adr/0031-*.md` und die
  Index-Zeile zeigen `Proposed` gegen den `done/`+`v0.38.0`-Stand.

### F-5 — LOW — DoD-`grep`-Aufzählung unvollständig: `internal/hexagon/core/rules/targets.go:20` hält einen live Skript-Verweis, der nicht `done/`/`CHANGELOG`/ADR ist

- **kategorie:** LOW
- **quelle:** Maintainability / Doku-Drift (stale Verweis in Produkt-Code) /
  DoD-Genauigkeit
- **pfad:** `docs/plan/planning/next/slice-064-gate-consistency-tombstone.md:51`
  (DoD: „findet nur noch Historie (`done/`, `CHANGELOG`, ADR-Kontext)") vs.
  `internal/hexagon/core/rules/targets.go:20` (Paritäts-Kommentar
  „`tools/gate-consistency.sh` (`^[a-zA-Z]…`)")
- **befund:** Die DoD behauptet, nach der Entfernung fänden `grep`-Treffer „nur
  noch Historie". Real trägt `targets.go:20` einen Kommentar, der die
  Regel-Heuristik-Parität gegen `tools/gate-consistency.sh` beschreibt — weder
  `done/`, noch `CHANGELOG`, noch ADR-Kontext, sondern **live Produkt-Code**. Nach
  der Entfernung verweist er auf eine nicht mehr existierende Datei. Kein
  Gate-Bruch (`codepaths` scannt Markdown, keine Go-Kommentare), aber die
  DoD-Aufzählung ist faktisch falsch und der Plan terminiert die Kommentar-Pflege
  nicht. (Verstärkt F-1: das DoD-Modell „verbleibende Verweise sind harmlose
  Historie" ist auf zwei Achsen unzutreffend — hier der live Go-Kommentar, dort die
  gate-geprüften Inline-Refs.)
- **verifizierbar:** ja — `grep -rn gate-consistency.sh internal/` zeigt den
  Treffer außerhalb der drei DoD-genannten Historien-Klassen.

### F-6 — INFO — Der spezifizierte Go-Test ist enger als die `DC-QA-03`-Messmethode: er prüft nur `external ∉ modules` und eine hart kodierte 5-Modul-Teilmenge, nicht „alle netzlosen Doku-Module"

- **kategorie:** INFO
- **quelle:** `DC-QA-03` (Messmethode) / `DC-QA-04` (Fidelity zum abgelösten Skript)
- **pfad:** Slice `:33`–`:37` und ADR-0032 `:48`–`:52` (Prüf-Set
  „`modules ⊇ {links, anchors, ids, matrix, codepaths}` sowie `external ∉ modules`")
  vs. `spec/lastenheft.md:1570` (Messmethode: „alle Module **außer `external` und
  `vcs`** aktiv") und `.d-check.yml:16` (`modules: [links, anchors, ids, matrix,
  codepaths, spans, hostpaths, versions]`, acht netzlose Doku-Module)
- **befund:** Der geplante Go-Test bildet das abgelöste Skript **paritätstreu** ab
  (`tools/gate-consistency.sh:24`–`:34` prüft ebenfalls nur die fünf Module und
  `external`). Er ist damit aber enger als (a) die aktuelle `DC-QA-03`-Messmethode,
  die `external` **und `vcs`** ausschließt, und (b) die reale Default-`modules`-Liste,
  die zusätzlich `spans`/`hostpaths`/`versions` als netzlose Doku-Module führt. §1
  des Plans (`:24`) und ADR-0032 (`:51`) formulieren „alle netzlosen Doku-Module",
  die konkrete Assertion pinnt jedoch fünf namentlich. **Failure (latent):** würde
  `spans`/`hostpaths`/`versions` aus der Default-Liste fallen oder `vcs` hinzugefügt,
  bliebe der Go-Test grün, obwohl die netzlose Beweis-Aussage geschwächt wäre. Auf
  dem aktuellen Baum ohne Wirkung; der Übergang zum typisierten Test wäre der
  natürliche Moment, die Assertion an die heutige `DC-QA-03`-Messmethode zu
  koppeln, statt die Skript-Teilmenge zu konservieren. Bewusst als INFO — Parität
  zum Skript ist eine legitime Wahl, aber die Prosa/Assertion-Differenz ist
  unbenannt.
- **verifizierbar:** teilweise — nur auf einem manipulierten `modules`-Stand
  sichtbar (Modul entfernt / `vcs` gesetzt).

## Negativbefunde (geprüft, ohne Befund)

- **Muster-Übertrag `diagnose_test.go` → `configyaml`-Paket:** geprüft — das
  Vorbild `internal/hexagon/core/app/diagnose_test.go:71` liest über
  `filepath.Join("..","..","..","..", …)` (vier Segmente tief). Das Zielpaket
  liegt in `internal/adapter/driven/configyaml/` — **ebenfalls vier Segmente tief**;
  der Relativ-Pfad-Tiefengrad überträgt sich deckungsgleich. Kein Befund.
- **Fail-closed-Design des Go-Tests:** „fehlende/undekodierbare `.d-check.yml` ⇒
  Test rot" (Slice `:37`, `:49`) deckt sich mit dem Fail-closed-Muster des
  Vorbilds (`diagnose_test.go:71`–`:107`, `t.Fatal` bei unlesbar/leer/mehrdeutig).
  Kein Befund.
- **Parität des Prüf-Kerns zum Skript:** das Prüf-Set (`{links, anchors, ids,
  matrix, codepaths}` + `external ∉`) entspricht `tools/gate-consistency.sh:25`–`:34`
  zeichensatz-genau; die Verlagerung ändert die Logik nicht (Breite s. F-6). Kein
  Befund.
- **Reduziertes Makefile-Target:** `gate-consistency: build` +
  `$(DCHECK_RUN) --enable targets $(FOCUS_DISABLE)` bleibt nach Entfall der
  `@bash …`-Zeile lauffähig; `build`-Prereq weiterhin nötig (Image-Lauf).
  `.PHONY`-Liste unverändert korrekt (Target existiert weiter). Kein Befund.
- **`targets`-Dogfood-Selbstbezug nach Entfernung:** das `gate-consistency`-Target
  bleibt bestehen → kein selbst-induziertes `gate-phantom`/`gate-undocumented` im
  `targets`-Lauf (`.d-check.yml:150`–`:154`, `authority: AGENTS.md`). Kein Befund.
- **`DC-QA-03`-Determinismus des Go-Tests:** der Test liest die live `.d-check.yml`
  read-only im Docker-`make test`-Kontext (Build-Kontext = Repo-Wurzel, wie
  `diagnose_test.go` `spec/spezifikation.md`); kein Schreibzugriff, kein Netz.
  Kein Befund.
- **Kein Lastenheft-CR / kein Release:** korrekt — keine neue/geänderte
  `DC-*`-Anforderung (die Verlagerung ist interne Gate-Mechanik), nur
  `internal/`-**Test**-Zuwachs, kein Produkt-Verhalten; analog `arch-check`-/
  `completeness`-Rückbau (ADR-0032 `:73`–`:74`). Kein Befund.
- **ADR-Nummer / Bezug-Richtung:** ADR-0032 = nächste freie Nummer (nach 0031);
  „Revidiert"/„Bezug" zeigen auf 0031/0028/`DC-QA-03` (aufwärts bzw. auf frühere
  ADRs), kein Abwärts-Verweis eines Spec-Stratums. Index-Zeile (`README.md:47`)
  vorhanden. Kein Befund (der Status-Flip ist F-4, nicht hier).
- **Referenz-Richtung (SDP) / Marker-Ehrlichkeit:** der Slice-Plan trägt keine
  getarnte ADR→Slice-Entscheidungsgrundlage; Slice-/ADR-Tokens erscheinen als
  Provenance. Kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 3 |
| LOW | 1 |
| INFO | 1 |

## Verdikt

**NACHBESSERN (blockierend).**

Begründung: Der Grundgedanke ist tragfähig und präzedenzgestützt (Voll-Tombstone
wie `planning`; Go-Test statt Shell für eine reine Decode-Assertion; das
`configyaml`-Paket ist tiefengleich zum genannten Vorbild; fail-closed sauber
übernommen; kein CR/kein Release korrekt). Blockierend ist **F-1 (HIGH):** der
Plan entfernt das letzte zitierte Gate-Skript, ohne den `codepaths.ignore-refs`-
Tombstone einzuplanen, der bei allen fünf Vorgänger-Entfernungen zwingend war —
`make doc-check`/`make gates` läuft danach mit ~20 `codepath-missing`-Befunden rot,
und die DoD-Aussage „`make gates` grün" plus „nur noch Historie" ist so nicht
erreichbar (die Fix-Größe ist eine Zeile, die Lücke ist die fehlende Planung).
**F-2 (MEDIUM)** ergänzt die zweite, ignore-refs-**un**abhängige Bruchachse
(`links`/`target-missing`, inkl. `done/`-Edit). **F-3/F-4 (MEDIUM)** betreffen
Harness-Ehrlichkeit (Sensors-`DC-QA-03`-Bindung wandert, wird aber nicht
umgeschrieben) und die Plan-Prämisse (ADR-0031 real noch `Proposed`). **F-5 (LOW)**
und **F-6 (INFO)** sind nachrangig, aber F-5 belegt die falsche DoD-Annahme
zusätzlich. Alle Befunde sind reine Plan-/Doku-Ergänzungen ohne Konzeptbruch.
