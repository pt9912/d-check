# Review-Report: slice-120 — `Stand`-Zellen und Drift-Log auf Zustand und Beleg

**Datum:** 2026-08-22 · **Review-Art:** Doku-/Konventions-Review (geprüft gegen Slice-Plan slice-120 §1/§2 inkl. Schritt 4a/§3 NICHT-Liste/§5 Risiken, Wellendokument welle-81 §1 Treffer-Tabelle/§3 Closure-Trigger/§6 Out-of-Scope, Hard Rules AGENTS §3.5/§3.6/§3.7, `MR-000` (Baseline-Aussage), `MR-025` (Spiegel vor dem Editieren), `MR-029` (v5.9.0-Pin), Kanon `.harness/baseline/v5.9.0/regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt (*Dieselbe Regel für Zustandsfelder*, *Die Kopfzeile eines lebenden Registers ist derselbe Fall*) und `.harness/baseline/v5.9.0/regelwerk/modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte / §Das Beobachtungs-Register, Ziel-Formen `.harness/baseline/v5.9.0/templates/docs/plan/planning/roadmap.template.md` und `.harness/baseline/v5.9.0/templates/docs/plan/planning/observations.template.md`) mit eigenen Gegenproben am gebauten Image
**Gegenstand:** Commit `8c6e8db` (Range `86188c4..8c6e8db`) — acht `Stand`-Zellen in `docs/plan/planning/observations.md` neu gefasst, Drift-Log in `docs/plan/planning/in-progress/roadmap.md` von 69 auf 10 Datenzeilen zurückgeschnitten samt neuer Sektions-Prosa, Meilenstein-Status-Form ergänzt, `**Stand:**`-Zeile in `harness/conventions.md` §Baseline neu gefasst; **vor** der Closure, kein Release, kein Produkt-Code (drei Markdown-Dateien, 28 Einfügungen / 87 Löschungen)
**Skill:** `.harness/skills/reviewer.md` @ 1.7.0 · **Modell-ID:** `claude-opus-5[1m]`
**Eingangs-Kontext:** Slice-Plan `docs/plan/planning/in-progress/slice-120-register-und-drift-log.md`; Wellendokument `docs/plan/planning/welle-81-zustandsfelder.md`; geschlossene Vorgänger `docs/plan/planning/done/slice-117-baseline-v590-bump.md`, `docs/plan/planning/done/slice-118-zustandsfeld-regel.md`, `docs/plan/planning/done/slice-119-kopf-zustandsfelder.md` samt ihren Reports in `docs/reviews/`; Vorfassungen `86188c4:docs/plan/planning/observations.md` und `86188c4:docs/plan/planning/in-progress/roadmap.md`; `harness/conventions.md` §Baseline und die sieben Einträge der Pin-Serie in `harness/conventions/` und `harness/conventions/done/`; `.d-check.yml` (§structure-Chronologie-Regeln), `.d-check.closure.yml`, `Makefile`. Nicht erhalten: die DoD-Abhakung (Verifikations-Rolle). Kein `make`-Target im echten Repo — alle Proben liefen als Image-Lauf (`d-check:latest`, netzlos, read-only) gegen `.git`-freie Baum-Kopien außerhalb des Repos; Exit je Lauf in eine Datei umgeleitet und separat gelesen (Arbeitsregel `BEO-007`).

## Findings

### F-1 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** Kanon `.harness/baseline/v5.9.0/regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt (*Dieselbe Regel für Zustandsfelder*: „nennt den Zustand und den Beleg — als **auflösbaren** Anker") / Slice-Plan slice-120 §2 Schritt 4a
- **pfad:** `harness/conventions.md:31-32`
- **befund:** Die neue Stand-Zeile ersetzt die Ketten-Nacherzählung durch eine Anweisung, wo die Kette zu finden ist: „jeder Eintrag nennt seinen Vorgänger im Feld `Löst auf:`". Das Feld existiert im ganzen Repo **einmal** — in `harness/conventions/MR-029-baseline-v590.md:67`, wo es die slice-117-Review-Auflage eingetragen hat (`6c42d82`). Von den sieben Einträgen der Pin-Serie tragen es sechs nicht: `MR-023-baseline-v500.md`, `MR-026-baseline-v560.md` und `MR-028-baseline-v570.md` tragen **weder** `Löst auf:` **noch** `Aufgelöst durch:`, und `MR-011`, `MR-012`, `MR-016` tragen ausschließlich `Aufgelöst durch:` — die Gegenrichtung. Wer der genannten Anweisung folgt, kommt von `MR-029` genau einen Hop weit nach `MR-028` und steht dann vor einem Eintrag ohne Vorgänger-Feld; die Serie zurück bis `MR-011` ist über das genannte Feld nicht erreichbar. Der zweite Halbsatz („die aufgelösten liegen in `harness/conventions/done/`") führt zusätzlich an `MR-023` vorbei, das als Serien-Glied im aktiven Verzeichnis steht. Die Kette **ist** rekonstruierbar — jeder Eintrag nennt seinen Vorgänger im `Adaption:`-Absatz und im Titel („N-ter Nachtrag zu MR-011") —, aber das ist nicht der Weg, den das Zustandsfeld nennt. Damit hat der Commit die Kette gegen einen Anker eingetauscht, der sie nicht auflöst; dieselbe Ungenauigkeit steht bereits im Slice-Plan §2 Schritt 4a, der zusätzlich `Ausgelöst durch Baseline-Stand:` nennt — ebenfalls ein Feld, das nur `MR-029` führt.
- **verifizierbar:** nein — kein Gate liest Zustandsfeld-Semantik oder Feld-Vollständigkeit in `MR-*`-Einträgen (Gegenprobe 1 bleibt bei Exit 0 / 0 Befunden). Belegt per `grep -rl "Löst auf" harness/conventions/` (ein Treffer) und `grep -rl "Aufgelöst durch" harness/conventions/` (16 Treffer, darunter weder `MR-026` noch `MR-028`).
- **klasse:** zustandsfeld-nennt-einen-aufloesungsweg-den-der-bestand-nicht-fuehrt

### F-2 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** Kanon `.harness/baseline/v5.9.0/regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register („Der Beleg ist formgebunden — drei Prüfungen ohne Urteil: **Form** · **Anzahl** (so viele wie der Zähler) · **Lage**") / Slice-Plan slice-120 §5 Risiko 1
- **pfad:** `docs/plan/planning/observations.md:12` und `docs/plan/planning/observations.md:16`
- **befund:** Zwei Zeilen des Registers weichen von der Beleg-Anzahl-Form ab: `BEO-009` trägt Zähler 2 mit **einem** Beleg (`slice-111`), `BEO-004` trägt Zähler 3 mit **einem** Beleg (`slice-101`). Beide Vorfassungen deklarierten die Abweichung ausdrücklich und begründeten sie mit der Form selbst — „ein Beleg-Slice (die Beleg-Form verlangt Slice-Kennungen; die Klasse ist dichter als der Zähler)" bzw. „die Beleg-Form verlangt Slice-Kennungen, darum **eine** Kennung bei Zähler 3". Das ist keine Chronik, sondern eine benannte Grenze der eigenen Form; mit der Neufassung ist sie aus beiden Zellen verschwunden, und die Formulierung steht danach nirgends mehr im Dokumentbestand (`grep` nach „dichter als der Zähler", „Beleg-Form verlangt", „ein Beleg-Slice" außerhalb des vendorten Baums und der Review-Reports: null Treffer). Das Register trägt jetzt in zwei von acht Zeilen einen unerklärten Widerspruch zu der Form, die es im eigenen Kopf zitiert („so viele wie der Zähler"). Versagen: Wer den Sichtungs- oder Lese-Schritt fährt, kann nicht unterscheiden, ob die Zeile eine deklarierte Abweichung oder einen Fehler trägt — und die im Kanon vorgezeichnete Mechanisierung der Anzahl-Prüfung meldete beide Zeilen rot, ohne dass das Register ihr etwas entgegenzuhalten hätte. Der Erklärungstext ist je einen Hop entfernt (`docs/plan/planning/done/slice-111-wave-drift-zwei-haelften.md:210-214` bzw. `docs/plan/planning/done/welle-70-results.md:75-77`), aber genau diese Auflösbarkeit behauptet die Zelle für den **Zustand**, nicht für die Formabweichung.
- **verifizierbar:** nein — Gegenprobe 6 setzt den `BEO-009`-Zähler auf 7 bei unverändert einem Beleg: volles Profil **und** `planning`-Profil bleiben bei Exit 0 / 0 Befunden. Die Anzahl-Prüfung ist in diesem Repo nicht gebaut (kein Grund-Code dafür in `internal/hexagon/core/model/finding.go`).
- **klasse:** deklarierte-abweichung-von-der-eigenen-beleg-form-mitentfernt

### F-3 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** Kanon `.harness/baseline/v5.9.0/regelwerk/modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte, Bullet *Meilensteine* („Ein erreichter Meilenstein bleibt in der Tabelle: `Status` sagt *erreicht* mit Datum und Beleg") / Slice-Plan slice-120 §2 Schritt 4
- **pfad:** `docs/plan/planning/in-progress/roadmap.md:55` gegenüber `docs/plan/planning/in-progress/roadmap.md:57-59`
- **befund:** Der Commit setzt unter die Meilenstein-Tabelle den Satz „Erreichte Meilensteine bleiben **hier** stehen; die Status-Zelle erzählt nicht, wie es dazu kam." Die Tabelle darüber trägt unverändert die Platzhalter-Zeile `— keine offenen —`. Beides zusammen ist falsch: das Repo hat drei Meilensteine erreicht (`M1: Spec-Fundament steht` / welle-01, `M2: Dogfooding` / welle-02, `M3: erstes GHCR-Release + Pilot-Migration` / welle-04, letzterer mit Status „**erreicht** (2026-06-12: v0.1.0/v0.2.0 auf GHCR, drei Repos migriert)"), und alle drei wurden am 2026-08-02 mit `b84b2cd` aus der Tabelle entfernt. Der Abschnitt behauptet damit eine Regel, die er selbst verletzt, und der Platzhalter-Wortlaut *keine **offenen*** kodiert obendrein die Gegenannahme — er liest sich als gefilterte Sicht, in der erreichte Einträge legitim fehlen. Versagen: Ein Audit fragt „welche extern beobachtbaren Zustände hat dieses Repo erreicht?", bekommt von der Roadmap „keine" und findet die Antwort nur in `git`; und der nächste erreichte Meilenstein wird nach demselben Muster aus der Tabelle genommen, sobald er nicht mehr offen ist. Der Slice hat die Fläche ausdrücklich angefasst (§2 Schritt 4) und den Ist-Zustand als „sie führt derzeit keinen offenen Eintrag" gelesen, statt zu prüfen, ob erreichte fehlen.
- **verifizierbar:** nein — `structure` führt auf `## Meilensteine` keine Regel (die vier Roadmap-nahen Regeln in `.d-check.yml` liegen auf `## Historische Trigger-Verschiebungen` und `## Abgeschlossene Wellen`); Gegenprobe 1 bleibt bei 0 Befunden. Belegt per `git log -p --follow` auf `docs/plan/planning/in-progress/roadmap.md` (drei `erreicht`-Zeilen, entfernt in `b84b2cd`).
- **klasse:** neue-regel-prosa-widerspricht-der-tabelle-ueber-der-sie-steht

### F-4 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** Kanon `.harness/baseline/v5.9.0/regelwerk/modul-06-roadmap.md` §Roadmap-Struktur, Bullet *Historische Trigger-Verschiebungen* („Das Drift-Log sagt, *was umgeplant* wurde … **und sonst nichts**") / `BEO-009` (Botschaft behauptet eine Änderung, die so nicht stattfand)
- **pfad:** `docs/plan/planning/in-progress/roadmap.md:121`, `docs/plan/planning/in-progress/roadmap.md:123`, `docs/plan/planning/in-progress/roadmap.md:124`
- **befund:** Die Commit-Botschaft sagt „drei Wellen-Eröffnungen sind auf ihren Umplanungs-Kern getrimmt". Gemessen betrifft der Trim ausschließlich Spalte 2 (*Was wurde geändert?*), die von 252, 345 bzw. 378 auf je 56 Zeichen fällt; Spalte 3 (*Warum?*) ist in allen drei Zeilen **byte-identisch** mit der Vorfassung und trägt 838, 1 063 bzw. 1 076 Zeichen — rund 95 Prozent der jeweiligen Zeile. Die welle-79-Zeile bleibt damit mit 1 753 Zeichen die längste der ganzen Tabelle, die welle-77-Zeile die zweitlängste, und ihr Inhalt ist zu großen Teilen keine Umplanungs-Begründung: das Eröffnungs-Sichtungs-Protokoll der Welle („**Beim Eröffnungs-Sichten:** das Register führt keine unverkörperte Beobachtung — BEO-006/BEO-007 stehen bei 2 mit gelebtem Gegenmittel"), eine Review-Auflagen-Nummer („Zählung präzisiert auf Review-Auflage F-2"), eine Release-Ankündigung („Release Minor v0.62.0") und Delta-Audit-Ergebnisse. Der Kanon weist genau diese Inhalte anderen Orten zu (Welle-Datei, Review-Report, Closure-Log). Versagen: Der Closure-Trigger der Welle („Keine `Stand`-Zelle und keine Drift-Log-Zeile erzählt eine Chronik — gemessen an der eigenen Datei, nicht behauptet") wird gegen eine Botschaft gelesen, die einen Trim behauptet, den rund 95 Prozent des Zeilentextes nicht erfahren haben; und die Zeile hält einen Register-Stand fest (`BEO-007` bei 2), den das Register zwei Dateien weiter mit 3 und *verkörpert* beantwortet.
- **verifizierbar:** nein — kein Gate liest Drift-Log-Semantik; die Chronologie-Regel prüft nur Spalte 1. Belegt per Feld-Vergleich der drei Zeilen gegen `86188c4:docs/plan/planning/in-progress/roadmap.md` (Spalte 3 identisch, Spalte 2 gekürzt).
- **klasse:** trim-erreicht-nur-eine-spalte-die-botschaft-behauptet-die-zeile

### F-5 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** AGENTS §3.7 („**Für Zustandsfelder gibt es keine Bestandsgrenze:** der vorhandene Bestand wird mit dem v5.9.0-Bump umgestellt, nicht grandfathered") / Kanon `.harness/baseline/v5.9.0/regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt, *Dieselbe Regel für Zustandsfelder* / Wellendokument welle-81 §1 (Treffer-Tabelle, „Fünf Treffer im eigenen Bestand, alle gemessen")
- **pfad:** `docs/plan/planning/done/slice-025-doctor.md:3` (exemplarisch; gleiche Bauform in 89 weiteren Dateien)
- **befund:** Nach diesem Commit erzählt eine **sechste** Fläche weiter Chronik in einem Zustandsfeld, und sie ist im Zensus der Welle nicht enthalten: 90 der 118 Slice-Dateien in `docs/plan/planning/done/` tragen im Kopf ein `**Status:**`-Feld — die Form, die der aktuelle Slice-Kopf ausdrücklich abgeschafft hat („Der Zustand dieses Slice ist das **Verzeichnis** …, bewegt per `git mv` — kein Status-Feld"). Viele davon erzählen, *wie* der Zustand entstand: `**Status:** done (Closure 2026-06-21; Review R1+R2, kein ADR nötig).`, `**Status:** done — **welle-60, abgeschlossen 2026-07-18** (Review R4 …)`, `**Status:** done (welle-62), **abgeschlossen 2026-07-18 (v0.50.0)**`. Elf davon behaupten überdies einen Zustand, der dem Verzeichnis **widerspricht**: `slice-025`, `slice-026`, `slice-027`, `slice-028` und `slice-083` sagen `open`, `slice-056`, `slice-057`, `slice-068`, `slice-070` sagen `in-progress`, `slice-084` und `slice-085` sagen `In Arbeit` — alle elf liegen in `done/`. Versagen: Zwei Zustandsaussagen derselben Datei sagen Gegensätzliches, und das Feld ist die auffälligere von beiden; das ist genau der Grund, aus dem `slice-091` es abgeschafft hat. Die Fläche steht außerhalb des Closure-Triggers der Welle (der nennt `Stand`-Zelle und Drift-Log-Zeile) und außerhalb des Slice-Scopes — sie gehört damit in den Register-Lese-Schritt der Welle-Closure, nicht in diesen Commit; gemeldet ist, dass der Zensus „fünf Treffer, alle gemessen" sie nicht kennt, obwohl §3.7 Zustandsfeldern die Bestandsgrenze ausdrücklich versagt.
- **verifizierbar:** nein — kein Gate liest Slice-Kopf-Felder gegen die Verzeichnis-Lage (Gegenprobe 2, das `planning`-Profil, bleibt bei Exit 0 / 0 Befunden, obwohl elf Dateien in `done/` `open` bzw. `in-progress` behaupten). Belegt per Auszählung über `docs/plan/planning/done/`.
- **klasse:** sechste-zustandsfeld-flaeche-ausserhalb-des-wellen-zensus

### F-6 · LOW

- **kategorie:** LOW
- **quelle:** Maintainability — die im Kopf von `version.md` deklarierte Invariante („**Nur die aktuelle Version** trägt einen expliziten HTML-Anker … Beim Release **wandert** der Anker zur neuen aktuellen Version; die bisherige Zeile verliert ihn")
- **pfad:** `version.md:35-36`
- **befund:** Zwei Zeilen tragen den Anker: `v0.62.0` und `v0.61.0`. Die Release-Prep zu `v0.62.0` (`f06fe8a`) hat den neuen Anker gesetzt, den alten aus `2b5d7f4` aber stehen lassen. Damit ist der Zweck des Registers für genau diese Version außer Kraft: ein fester Pin auf `version.md#v0.61.0` löst weiter auf, statt als `anchor-missing` zu melden, und der Kopf des Registers sagt ausdrücklich, dass dieses Melden der Zweck ist. Bestandslage außerhalb des Slice-Diffs; gefunden bei der von der Prüf-Aufgabe verlangten Sichtung von `version.md` auf Zustandsfelder.
- **verifizierbar:** ja — ein `anchors`-Lauf gegen ein Dokument, das `version.md#v0.61.0` referenziert, bleibt grün, obwohl `v0.61.0` nicht mehr die aktuelle Version ist; heute referenziert kein Dokument diesen Anker, deshalb schweigt jeder Gate-Lauf.
- **klasse:** wandernder-anker-nicht-mitgewandert

### F-7 · INFO

- **kategorie:** INFO
- **quelle:** AGENTS §3.5 (ADRs sind nach `Accepted` immutable) gegen AGENTS §3.7 (Zustandsfelder, keine Bestandsgrenze)
- **pfad:** `docs/plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md:3`
- **befund:** Das einzige `**Status:**`-Feld im ADR-Bestand, das mehr trägt als seinen Zustand, erzählt eine Chronik: „Accepted — **wieder aufgenommen und umgesetzt 2026-07-18, sicher auf ADR-0043.** Die 2026-07-17 ausgesetzte Fassung scheiterte fünfmal, weil …". Nach §3.7 wäre das ein Zustandsfeld ohne Bestandsgrenze, nach §3.5 ist die Zeile unantastbar (erlaubt bleiben nur `## Geschichte`-Anhänge und der Status-Übergang). Die beiden Hard Rules treffen sich hier ohne Vorrangregel; dokumentiert, weil die Welle die Frage sonst offen zurücklässt, wenn der nächste Durchgang den ADR-Bestand mit demselben Maß liest.
- **verifizierbar:** ja — `make adr-check` (Modul `vcs`) meldete eine inhaltliche Änderung dieser Zeile als `core-drift-vcs`; das ist der Beleg, dass sie nicht geräumt werden kann.
- **klasse:** zustandsfeld-regel-trifft-immutabilitaets-regel

### F-8 · INFO

- **kategorie:** INFO
- **quelle:** Kanon `.harness/baseline/v5.9.0/regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt („Ein zweites Log … ist eine Kopie, und Kopien driften")
- **pfad:** `docs/plan/planning/observations.md:11`
- **befund:** Die neue `Stand`-Zelle von `BEO-010` eröffnet mit „Gate-blind: `gate-consistency` prüft Target-**Namen**, nicht Modul-**Mengen**." — wörtlich der Schlusssatz der Beobachtungs-Spalte derselben Zeile. Als einzige der acht Zellen wiederholt sie ihre Nachbarzelle statt den Zustand zu nennen; wird die Beobachtung später präzisiert, ist die Kopie zwei Spalten weiter mitzuziehen. Keine der übrigen sieben Zellen tut das.
- **verifizierbar:** nein — kein Gate vergleicht Zellen einer Tabellenzeile. Belegt per Textvergleich innerhalb `docs/plan/planning/observations.md:11`.
- **klasse:** stand-zelle-wiederholt-ihre-beobachtungs-spalte

## Negativbefunde

1. **Zähler-, Belege- und Sub-Area-Spalten unverändert.** Alle acht Datenzeilen mechanisch feldweise gegen `86188c4` verglichen (Spalten 1 bis `NF-2`): byte-identisch. Die NICHT-Liste §3 („keine Zählerstände ändern sich") ist eingehalten.
2. **Sektion *Gestrichene Einträge* und Register-Kopf unberührt.** `BEO-001` und `BEO-005` sind zeilengleich mit der Vorfassung; der Kanon erlaubt dort die Begründung des Zustands, und beide tragen genau die (mechanisierte Instanz, benannter Rest, vorgezeichnete Form) — kein Handlungsbedarf.
3. **`BEO-007`: Gegenmittel erhalten und geschärft.** Die Verkörperungsstelle (`pre-commit` läuft den vollen `doc-check`) steht, und die Arbeitsregel ist ausdrücklicher als vorher („Gate-Läufe in eine Datei umleiten, den Exit **bindend** verwenden").
4. **`BEO-006`: Gegenmittel und 3×-Form vollständig.** Beide Halbregeln (`git status` vor pfad-selektiven Commits; Moves/Löschungen erst unmittelbar vor ihrem Commit stagen) und die Heuristik-Grenze der mechanischen Form sind erhalten.
5. **`BEO-002`: kein Substanzverlust.** Die entfallene Spiegel-Liste und der Blindfleck „in genau der Datei, die ohnehin bearbeitet wurde" stehen vollständig in `harness/conventions/MR-025-spiegel-vor-dem-editieren.md` — dem Eintrag, auf den die neue Zelle verlinkt; auch der Ableiter „`grep` nach dem alten Wortlaut" steht dort.
6. **`BEO-003`: Verkörperungs-Ort auflösbar, Zusatz belegt.** Der Name des Kopplungs-Tests entfällt aus der Zelle, steht aber unter dem genannten Anker (`docs/plan/planning/done/welle-74-results.md:91`, ebenso `docs/plan/planning/done/slice-103-geteilte-lexik-raender.md:166`); der Test existiert (`internal/hexagon/core/rules/lexikon_kopplung_test.go:215`). Der neue Zusatz „konfigurationsseitig" ist kein Neuzugang, sondern die bereits belegte vierte Achse aus `docs/plan/planning/done/welle-80-results.md:72` und `docs/plan/planning/done/slice-114-spec-vergabe-spezifikation.md:174`.
7. **`BEO-004`: Begründung der Verkörperungs-Form auflösbar.** Warum die kanonische Form (Frage im Slice-Template) hier nicht gewählt wurde, steht unter dem genannten welle-73-Anker: `docs/plan/planning/done/welle-73-results.md:106-108`, samt `MR-018`-Zeiger.
8. **`BEO-010`: der zitierte Kopplungs-Kommentar steht noch.** Die Vorfassung zitierte einen `Makefile`-Kommentar; er ist unverändert vorhanden (`Makefile:207`), die Aussage geht mit dem Kürzen nicht verloren.
9. **`BEO-008`: beide Richtungen und die CR-Grenze erhalten.** Vergessene Hebung **und** Über-Hebung, der Drei-Klassen-Zensus als Gegenmittel und die Begründung, warum die mechanische Form heute nicht konfigurierbar ist, stehen weiter; entfallen ist einzig der Klammerzusatz „(Schema geprüft)", der eine bereits erbrachte Probe markierte.
10. **Drift-Log-Deckung mechanisch nachgezählt.** Alle **20** im alten Log genannten Wellen-Kennungen haben eine Zeile im Closure-Log (das 21 Wellen führt), und alle **42** genannten Slice-Kennungen haben eine Datei in `docs/plan/planning/done/`. Die Deckungs-Behauptung des Commits hält.
11. **Die zehn verbliebenen Zeilen tragen echte Umplanungen.** Sieben sind byte-identisch übernommen (Etappe-C-Schnitt, `slice-104` als Change-Request-Einschub, Etappe-D-Neuordnung, Wiederaufnahme `slice-071`, Rückstellung `slice-074`, WIP-Limit-Wiederherstellung, `slice-012`-Trigger-Verschiebung), drei sind die im Kopf getrimmten Wellen-Eröffnungen mit „neu geschnitten" — je ein im Kanon aufgezählter Umplanungs-Typ. Keine Schließung, kein erreichter Meilenstein ist stehen geblieben.
12. **Die Umplanungs-Verdachtsfälle unter den 59 entfernten Zeilen einzeln nachgeprüft — keine ist die einzige Spur.** (a) „welle-70 eröffnet … der stille Grün-Pfad wird **vorgezogen**: bindende Vorbedingung der structure-Umsetzung" — die Vorrang-Beziehung steht in `docs/plan/planning/done/welle-69-results.md:104` (§Folge-Slices) und `docs/plan/planning/done/welle-70-fence-lexik.md:35` (§2 Trigger). (b) „welle-67 Etappe D eröffnet als Mini-Welle … in vier thematische Slices geschnitten (Nutzer-Entscheid)" — der Schnitt steht in `docs/plan/planning/done/welle-67-baseline-v500-migration.md:49` (§4) und `docs/plan/planning/done/welle-67-results.md:21`; die verbliebene Zeile „Etappe D neu geordnet" nennt die Endzuordnung aller vier Slices vollständig. (c) Der Reihenfolge-Entscheid „erst BEO-005, dann Migration" aus der getrimmten welle-77-Zeile steht unverändert in der verbliebenen welle-78-Zeile („Reihenfolge-Entscheid 2026-08-21: nach der Chronologie-Welle"). (d) Die welle-69-Closure-Zeile enthielt eine Trigger-Korrektur (`slice-099` fälschlich als Wellen-Mitglied) — sie steht in der Ergebnis-Notiz der Welle, dem vom Kanon vorgesehenen Ort.
13. **Die drei getrimmten Zeilen sind inhaltlich wahr.** „slice-105/106/110 neu geschnitten beim Öffnen von welle-77/78/79" deckt sich mit der jeweiligen Vorfassung („neu geschnitten und `open`→`in-progress`"); die Formulierung nimmt nichts hinzu.
14. **Chronologie-Regeln unberührt und scharf.** `.d-check.yml` ist nicht im Diff; die vier für die Roadmap relevanten Regeln stehen unverändert (Drift-Log `table-order: desc` auf Spalte 1, Closure-Log `table-order: desc` mit `table-column: 2`). Die verbliebene Drift-Tabelle ist monoton absteigend (2026-08-21 ×4 → 2026-08-10 → 2026-08-02 → 2026-07-18 → 2026-07-17 ×2 → 2026-06-11). Gegenproben 4 und 5 zeigen, dass beide Regeln nach dem Schnitt weiter feuern — sie sind nicht still übersprungen worden.
15. **Kopf-Zustandszeilen (slice-119) bleiben entfernt.** Weder `docs/plan/planning/in-progress/roadmap.md` noch `docs/plan/planning/observations.md` noch `spec/spezifikation.md` trägt eine `Status: Aktiv`-Kopfzeile; `spec/architecture.md:3` behält ihre samt Begründung — dieser Commit hat daran nichts verschoben.
16. **Kein weiteres Kopf-Datum ohne benannten Trigger.** Die `**Stand:**`-Zeile in `docs/user/benutzerhandbuch.md:4` wird nachweislich mit jedem Handbuch-Versions-Bump gesetzt (fünf aufeinanderfolgende Release-Prep-Commits in ihrer Zeilen-Historie), die Kopf-Daten der Skills unter `.harness/skills/` mit deren Versions-Bump — beide fallen unter die Trigger-Ausnahme des Kanons.
17. **MR-Bestand und ADR-Bestand sonst sauber.** Alle `**Status:**`-Felder in `harness/conventions/` und `harness/conventions/done/` tragen ausschließlich `Accepted`; im ADR-Bestand trägt genau eines mehr (F-7).
18. **`version.md` und `CHANGELOG.md` tragen kein Zustandsfeld mit Chronik.** Das Release-Register führt Versions-Koordinaten (Version, Datum, Tag) ohne Stand-Spalte; der `CHANGELOG`-Kopf führt Unreleased-Einträge, kein Statusfeld. (Der Anker-Befund F-6 ist eine andere Klasse.)
19. **Slice-§5-Risiken 2 und 3 sind adressiert und belegt.** Risiko 2 (einzige Spur) durch die Deckungs-Zählung (Negativbefund 10 und 12), Risiko 3 (Chronologie) durch Gegenproben 1, 4 und 5. Risiko 1 (Substanzverlust) ist überwiegend abgedeckt; die Ausnahme ist F-2.
20. **Scope-Treue.** Der Diff fasst drei Markdown-Dateien an — kein Produkt-Code, keine Prüf-Config, kein Gate, kein Release-Artefakt; §3 („Kein Produkt-Code") ist eingehalten. Die `86188c4`-Vorstufe ist ein reiner Lifecycle-Move nach `MR-013` und außerhalb dieses Prüfgegenstands unauffällig.

## Gegenproben (Exit explizit)

Alle Läufe netzlos und read-only gegen eine `.git`-freie Baum-Kopie außerhalb des Repos, Ausgabe je Lauf in eine Datei umgeleitet und der Exit separat gelesen (`BEO-007`).

| # | Baum | Lauf | Exit | Ausgabe |
|---|---|---|---|---|
| 1 | `8c6e8db` | Sollform, volles Profil (`doc-check`-Äquivalent) | 0 | 434 Datei(en), 0 Befund(e) |
| 2 | `8c6e8db` | `--enable planning` mit abgewählten Datei-Modulen (`planning-check`-Äquivalent) | 0 | 434 Datei(en), 0 Befund(e) |
| 3 | `8c6e8db` | `--config .d-check.closure.yml --enable planning --enable structure` (`verify-closure-notes`-Äquivalent) | 0 | 400 Datei(en), 0 Befund(e) |
| 4 | `8c6e8db` **plus** zwei getauschte Drift-Log-Zeilen | volles Profil | 1 | 2 Befund(e), `section-unordered` auf `## Historische Trigger-Verschiebungen` |
| 5 | `8c6e8db` **plus** getauschte erste/letzte Closure-Log-Zeile | volles Profil | 1 | 2 Befund(e), `section-unordered` auf `## Abgeschlossene Wellen` |
| 6 | `8c6e8db` **plus** `BEO-009`-Zähler auf 7 (ein Beleg unverändert) | volles Profil **und** `planning`-Profil | 0 / 0 | je 434 Datei(en), 0 Befund(e) |
| 7 | Arbeitsstand inklusive dieses Reports | volles Profil | 0 | 435 Datei(en), 0 Befund(e) |

Lesart: Proben 4 und 5 belegen, dass die beiden Chronologie-Regeln den Schnitt überlebt haben und auf den verkürzten Tabellen weiter scharf sind — die Spalten-Lage hat sich nicht verschoben. Probe 6 belegt die Gate-Blindheit hinter F-2: die Beleg-**Anzahl** ist in diesem Repo nicht mechanisiert, ein Zähler ohne passende Belegzahl bleibt still. Proben 1 bis 3 bestätigen die Probe der Commit-Botschaft (434 Dateien / 0 Befunde) am eigenen Lauf.

## Welle-Abschluss (welle-81 §3)

- **Alle vier Slices in `done/`** — offen: `slice-120` liegt noch in `in-progress/`; das ist der Closure-Move selbst.
- **`make fullbuild` grün (Exit explizit)** — Verifikations-Rolle, nicht geprüft.
- **Pin auf `v5.9.0`, alter Baum entfernt** — erfüllt: `.harness/baseline/` führt genau ein Verzeichnis (`v5.9.0`).
- **Kein lebendes Register und kein Technik-Stratum mit Kopf-Zustandszeile, die Sicht behält ihre** — erfüllt (Negativbefund 15).
- **Keine `Stand`-Zelle und keine Drift-Log-Zeile erzählt eine Chronik** — für die acht `Stand`-Zellen erfüllt; für das Drift-Log mit dem Vorbehalt aus F-4 (drei Zeilen tragen in Spalte 3 unverändert Eröffnungs-Protokoll, Review-Auflagen-Nummer und Release-Ankündigung). F-1 betrifft denselben Trigger auf dem fünften Treffer: die Chronik ist dort weg, aber der Ersatz-Anker löst nicht auf.
- **Ergebnisnotiz `welle-81-results.md` mit Register-Lese-Schritt** — steht aus. Für den Lese-Schritt liegen `BEO-008` (Schwelle erreicht, 3×, unverkörpert) und die in F-5 gemeldete sechste Fläche an.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 5 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** zustandsfeld-nennt-einen-aufloesungsweg-den-der-bestand-nicht-fuehrt · deklarierte-abweichung-von-der-eigenen-beleg-form-mitentfernt · neue-regel-prosa-widerspricht-der-tabelle-ueber-der-sie-steht · trim-erreicht-nur-eine-spalte-die-botschaft-behauptet-die-zeile · sechste-zustandsfeld-flaeche-ausserhalb-des-wellen-zensus · wandernder-anker-nicht-mitgewandert · zustandsfeld-regel-trifft-immutabilitaets-regel · stand-zelle-wiederholt-ihre-beobachtungs-spalte

## Verdikt

**Merge-blockierend:** ja — F-1 bis F-4 blockieren nach der Regel des Skills, und alle vier liegen auf Flächen, die dieser Slice selbst angefasst hat. F-1 ist der schwerste: die Kette wurde gegen einen Anker eingetauscht, der sie nicht auflöst, und derselbe Fehlgriff steht bereits im Slice-Plan §2 Schritt 4a. F-2 ist der einzige gefundene Substanzverlust im Sinne von §5 Risiko 1 — alles andere, was aus den acht Zellen verschwunden ist, war Chronik oder ist unter dem genannten Anker in einem Hop auflösbar (Negativbefunde 5 bis 9). F-3 und F-4 sind je eine Aussage, die der Commit neu setzt und die der Bestand daneben widerlegt.

**Nicht blockierend, obwohl MEDIUM:** F-5. Die Fläche liegt außerhalb des Slice-Ziels (§1: die beiden lebenden Register plus der fünfte Treffer) und außerhalb des Closure-Triggers der Welle; sie gehört in den Register-Lese-Schritt der `welle-81`-Closure — als Beobachtung oder als Folge-Slice. Gemeldet ist sie hier, weil AGENTS §3.7 Zustandsfeldern die Bestandsgrenze ausdrücklich versagt und der Wellen-Zensus „fünf Treffer, alle gemessen" damit unvollständig ist.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD- und Spec-Konformität prüft der Verifier separat.
