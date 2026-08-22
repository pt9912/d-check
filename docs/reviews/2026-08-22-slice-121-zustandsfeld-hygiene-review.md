# Review-Report: slice-121 — Zustandsfeld-Hygiene: dreizehn Felder geheilt, ein Anker geschärft, eine Vorrangregel

**Datum:** 2026-08-22 · **Review-Art:** Doku-/Konventions-Review (geprüft gegen Slice-Plan slice-121 §1/§2/§3/§4/§5/§7, Hard Rules AGENTS §3.5/§3.7/§5, Reviewer-Skill 1.7.0 (HIGH-Anker *Zustandsfeld trägt Chronik*), Kanon `.harness/baseline/v5.9.0/regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt (*Dieselbe Regel für Zustandsfelder*) und `.harness/baseline/v5.9.0/regelwerk/modul-05-planning-harness.md` §Lifecycle als State Machine, Ziel-Formen `.harness/baseline/v5.9.0/templates/docs/plan/planning/slice.template.md` und `.harness/baseline/v5.9.0/templates/docs/plan/adr/NNNN-titel.template.md`, Beobachtungs-Register `BEO-002`/`BEO-009`) mit eigenen Gegenproben am gebauten Image
**Gegenstand:** Commit `1c26b92` (Range `de6b705..1c26b92`) — dreizehn `**Status:**`-Kopffelder in `docs/plan/planning/done/slice-*.md` auf `done` gesetzt, der `<a id>`-Anker der Vorgänger-Version aus `version.md` entfernt, Vorrangregel in `AGENTS.md` §3.7 und dieselbe Nicht-Melde-Ausnahme in `.harness/skills/reviewer.md`; wellenlos, **vor** der Closure, kein Release, kein Produkt-Code
**Skill:** `.harness/skills/reviewer.md` @ 1.7.0 (Stand `1c26b92`; siehe F-5) · **Modell-ID:** `claude-opus-5[1m]`
**Eingangs-Kontext:** Slice-Plan `docs/plan/planning/in-progress/slice-121-zustandsfeld-hygiene.md` (§1 Ziel mit den gemessenen Zahlen, §2 Vorgehen, §3 NICHT-Liste, §4 DoD, §5 Risiken, §7 Vorgelagert); Herkunft der drei Befunde in `docs/reviews/2026-08-22-slice-119-kopf-zustandsfelder-review.md` (F-5) und `docs/reviews/2026-08-22-slice-120-register-und-drift-log-review.md` (F-6, F-7); geschlossene Ergebnisnotiz `docs/plan/planning/done/welle-81-results.md` §Folge-Slices; Beobachtungs-Register `docs/plan/planning/observations.md`; Konfiguration `.d-check.yml` (Blöcke `vcs`, `versions`, `structure`) und `.d-check.closure.yml`; Regelcode `internal/hexagon/core/rules/vcs.go`; Spiegel `harness/README.md`, `harness/conventions.md`, `docs/plan/planning/README.md`, `docs/user/releasing.md`. Nicht erhalten: die DoD-Abhakung (Verifikations-Rolle). Kein `make`-Target im echten Repo — alle Proben liefen als Image-Lauf (`d-check:latest`, `--network none`, read-only) gegen eine `.git`-freie Baum-Kopie außerhalb des Repos; Exit je Lauf explizit gelesen und in eine Datei umgeleitet (Arbeitsregel `BEO-007`).

## Findings

### F-1 · HIGH

- **kategorie:** HIGH
- **quelle:** AGENTS §3.7 („**Für Zustandsfelder gibt es keine Bestandsgrenze:** der vorhandene Bestand wird mit dem v5.9.0-Bump umgestellt, nicht grandfathered") / Reviewer-Skill 1.7.0, HIGH-Anker *Zustandsfeld trägt Chronik* / Kanon `.harness/baseline/v5.9.0/regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt, *Dieselbe Regel für Zustandsfelder* / Quell-Befund `docs/reviews/2026-08-22-slice-119-kopf-zustandsfelder-review.md` F-5 / `BEO-009`
- **pfad:** `docs/plan/planning/in-progress/slice-121-zustandsfeld-hygiene.md:84-86` (DoD-Punkt 1); Bestand exemplarisch `docs/plan/planning/done/slice-074-kommentar-suffix-tabellenzeilen.md:3`, `docs/plan/planning/done/slice-076-markdown-lexik-commonmark.md:3`, `docs/plan/planning/done/slice-077-stiller-tabellen-uebersprung.md:3`, `docs/plan/planning/done/slice-079-zitat-verifikation.md:3`, `docs/plan/planning/done/slice-036-rtm-trace.md:3`
- **befund:** Die DoD sagt „keines trägt Chronik im Feld (gemessen: vorher elf plus zwei, nachher null)". Absatzweise nachgezählt trägt der `done/`-Bestand nach dem Commit 90 `**Status:**`-Kopffelder; 13 lesen jetzt exakt `done`, 27 tragen Zustand oder Zustand plus Wellen-Beleg — und **50** erzählen weiter, *wie* der Zustand entstand (Kriterium: Datum, `abgeschlossen`, `Closure`, `Review`, `Release`, `Lifecycle`). Die zwei als „Chronik" geheilten Felder (`slice-071`, `slice-078`) sind genau und ausschließlich die beiden, deren Zeile die wörtliche Zeichenfolge `abgeschlossen am` enthält; `slice-074`, `slice-076` und `slice-077` tragen die praktisch gleichlautende Form `done — **welle-60, abgeschlossen 2026-07-18** (Review R… ` — ohne `am` — und bleiben stehen. Die drei Beispiele, die der Quell-Befund (slice-119-Review F-5) wörtlich als Chronik zitiert (`done (Closure 2026-06-21; Review R1+R2, kein ADR nötig)`, `done — **welle-60, abgeschlossen 2026-07-18** (Review R4 …)`, `done (welle-62), **abgeschlossen 2026-07-18 (v0.50.0)**`), sind **alle drei** unangetastet. Damit trennt kein benanntes Kriterium die zwei geheilten von den 50 verbliebenen: §3 begründet die Ausnahme mit „Sie widersprechen nichts", was auf `slice-071`/`slice-078` genauso zutraf. Versagen: Der DoD-Punkt wird bei der Closure gegen eine Zahl abgehakt, die der Bestand widerlegt, und ein späterer Leser des Slice hält die Fläche für geräumt.
- **verifizierbar:** nein — kein Gate liest Slice-Kopf-Felder (Probe 4, `--enable planning`, Exit 0 / 0 Befunde bei 50 chronikführenden Feldern; Probe 5, Closure-Profil, Exit 0 / 403 Dateien). Belegt per absatzweisem Zensus über alle 119 `done/`-Slices.
- **klasse:** heil-auswahl-folgt-dem-suchstring-nicht-der-klasse

### F-2 · HIGH

- **kategorie:** HIGH (Kontext-Eskalation: die Aussage beschreibt einen Gate-Pfad)
- **quelle:** AGENTS §3.5 („erlaubt bleiben `## Geschichte`-Anhänge + der `**Status:**`-Übergang") / `harness/README.md` §Traceability rules / ADR-0016 und ADR-0024 / `.d-check.yml` Block `vcs` (`status-line`, `head-allow`) / `internal/hexagon/core/rules/vcs.go` (`vcsCore`)
- **pfad:** `AGENTS.md:182-185`
- **befund:** Die neue Vorrangregel sagt: „das `**Status:**`-Feld einer ADR gehört zu ihrem Kern — ist die ADR `Accepted`, ist es immutabel, und §3.5 sticht". Beide Hälften der Begründung sind gegen den Mechanismus falsch. `AGENTS.md:147` sagt drei Abschnitte darüber das Gegenteil („erlaubt bleiben … der `**Status:**`-Übergang"), `harness/README.md:150-153` spiegelt denselben Satz, `.d-check.yml:310-311` deklariert `status-line: '^\*\*Status:\*\*'` und `head-allow: '^\*\*Status:\*\* (Accepted|Superseded by ADR-[0-9]{4})'`, und `internal/hexagon/core/rules/vcs.go` streicht in `vcsCore` genau diese Kopfzeile **aus** dem Core, bevor er BASE gegen HEAD vergleicht. Die Status-**Zeile** ist also ausdrücklich nicht Kern; Kern sind erst ihre Fortsetzungszeilen. Versagen: Ein Agent, der eine `Accepted`-ADR per `Supersedes ADR-NNNN` ablösen soll — der von §3.5 vorgesehene Korrekturweg —, liest §3.7 und verweigert den Übergang `Accepted` → `Superseded by ADR-NNNN`, den `make adr-check` erlaubt; umgekehrt meldet ein Reviewer den erlaubten Übergang als Regelbruch. Die Regel, die eine Kollision auflösen soll, erzeugt eine zweite.
- **verifizierbar:** ja — `make adr-check` (`RANGE=`/`STAGED=`) über einen Commit, der ausschließlich die erste `**Status:**`-Zeile einer `Accepted`-ADR umschreibt, bleibt grün; der Beleg für „immutabel" existiert nur für die Fortsetzungszeilen.
- **klasse:** vorrangregel-behauptet-eine-immutabilitaet-die-das-gate-nicht-kennt

### F-3 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** AGENTS §3.7 gegen Slice-Plan §2 Schritt 4 („für **neue** ADRs gilt §3.7 ab dem ersten Schreiben") / Reviewer-Skill 1.7.0, MEDIUM-Anker *Adressierungs-Form eines Neuzugangs* (dort ist derselbe Fall ausformuliert)
- **pfad:** `AGENTS.md:185`
- **befund:** Der Satz lautet „Für **neue** ADRs gilt diese Regel ab dem ersten Schreiben" und steht unmittelbar hinter der Vorrangregel. Das nächstliegende Bezugswort von „diese Regel" ist die eben gesetzte Vorrangregel (§3.5 sticht) — gelesen ergibt der Satz, dass **neue** ADRs schon beim ersten Schreiben unter die Ausnahme fallen, also das genaue Gegenteil der Absicht aus Slice §2 Schritt 4 und der Commit-Botschaft. Die gemeinte Regel (§3.7) wird nicht benannt, obwohl beide Regeln im selben Satzgefüge stehen. Der Skill löst dieselbe Zweideutigkeit an anderer Stelle ausdrücklich auf („die Form gilt beim Schreiben, nicht beim Status-Übergang"); hier fehlt der Zusatz. Versagen: Eine neu geschriebene ADR bekommt eine Chronik ins Status-Feld, weil das Briefing so gelesen wird, dass die Ausnahme von Anfang an gilt — und nach dem `Accepted`-Übergang ist die Chronik nicht mehr räumbar.
- **verifizierbar:** nein — Formulierungsfrage, kein Gate. Belegt per Vergleich mit `docs/plan/planning/in-progress/slice-121-zustandsfeld-hygiene.md:65-68`.
- **klasse:** pronominal-referenz-kehrt-die-vorrangregel-um

### F-4 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** Reviewer-Skill 1.7.0, HIGH-Anker *Zustandsfeld trägt Chronik* (neue Ausnahme) gegen AGENTS §3.7 (letzter Satz der Vorrangregel) / Reviewer-Skill, MEDIUM-Anker *Adressierungs-Form eines Neuzugangs*
- **pfad:** `.harness/skills/reviewer.md:42-46`
- **befund:** Briefing und Skill decken sich in dem, **was** sie ausnehmen (das `**Status:**`-Feld einer `Accepted`-ADR), aber nicht in dem, was die Ausnahme begrenzt: der Skill trägt die Neu-ADR-Hälfte gar nicht. Ein Reviewer, der nur den Skill liest — und der Skill ist für die Kategorien-Anker die führende Quelle —, nimmt das Status-Feld **jeder** `Accepted`-ADR aus, auch einer, die im gerade geprüften Slice `Accepted` wurde. Versagen: Ein Slice schreibt eine neue ADR mit Chronik im Status-Feld und flippt sie im selben Slice auf `Accepted`; der Review meldet nichts, weil die Ausnahme am Zustand `Accepted` hängt und nicht am Zeitpunkt des Schreibens, und §3.5 friert die Chronik danach ein. Genau diesen Weg schließt der Skill für `Schärft:`/`Bezug:`-Felder eine Kategorie tiefer ausdrücklich („auch dann, wenn die ADR im selben Slice `Accepted` wird") — beim neuen Anker fehlt der Vorbehalt.
- **verifizierbar:** nein — Skill-Text, kein Gate. Belegt per Textvergleich `AGENTS.md:182-185` gegen `.harness/skills/reviewer.md:42-46`.
- **klasse:** skill-ausnahme-ohne-neuzugangs-vorbehalt

### F-5 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** Reviewer-Skill §Eingangs-Kontext („Pflicht — sonst nicht reproduzierbar") und §Ablage (Kopf-Metadaten: „**Skill** (`reviewer.md` @ Version/Commit)") / gelebte Praxis der drei Vorgänger-Änderungen
- **pfad:** `.harness/skills/reviewer.md:3`, `.harness/skills/reviewer.md:45`
- **befund:** Der Commit ändert den Text des HIGH-Ankers, lässt aber `**Version:** 1.7.0 · **Datum:** 2026-08-22` stehen; die geänderte Passage sagt weiterhin „neuer HIGH-Eintrag seit 1.7.0". Damit tragen zwei verschiedene Regeltexte dieselbe Versionsnummer: `68bafad` (Einführung, ohne ADR-Ausnahme) und `1c26b92` (mit Ausnahme). Die drei vorangegangenen inhaltlichen Änderungen haben jeweils gebumpt (1.4.0 → 1.5.0 → 1.6.0 → 1.7.0). Die fünf zuletzt abgelegten Reports pinnen den Skill ausschließlich über die Version, nie über den Commit. Versagen: Wer die Reports zu slice-118/119/120 („@ 1.7.0") später gegen den Skill hält, kann nicht entscheiden, ob die ADR-Ausnahme beim jeweiligen Lauf in Kraft war — der Report ist als Lauf-Beleg genau in dem Punkt nicht mehr reproduzierbar, den der Skill für Pflicht erklärt.
- **verifizierbar:** nein — kein Gate liest Skill-Versionen. Belegt per `git show 68bafad:.harness/skills/reviewer.md` gegen `HEAD`, Zeile 3 in beiden `**Version:** 1.7.0`.
- **klasse:** skill-inhalt-geaendert-version-eingefroren

### F-6 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** Slice §5 Risiko 1 („Ein `done/`-Slice ist ein Lauf-Beleg") / `MR-013` (Bijektion Wellen-Zeiger ⟺ Dateien) / Maintainability
- **pfad:** `docs/plan/planning/done/slice-070-trace-tabellenquellen-nullmengen-guard.md:3` (geheiltes Feld), `:5-7` (`**Welle:**`)
- **befund:** Das geheilte Feld las `in-progress (welle-59-trace-tabellenquellen)`. Anders als in den zwölf übrigen geheilten Dateien fängt der `**Welle:**`-Block diese Angabe nicht auf: er lautet „aktiv; Vorgänger `slice-069` ist abgeschlossen" und nennt keinen Wellen-Namen; auch §7 nennt keinen. Ein repo-weites `grep` nach `welle-59` liefert nach dem Commit **null** Treffer — die Wellen-Zuordnung von `slice-070` existiert nur noch in der git-Historie. Für `slice-071` und `slice-078` gilt das nicht: dort steht jede entfernte Angabe (Datum, Lifecycle-Kette, ADR-Status, Release-Tags, Review-Verdikte, der `architecture.md`-Beleg) vollständig in der Closure-Notiz §7. Versagen: Wer die Wellen-Zuordnung eines geschlossenen Slice aus dem Baum rekonstruiert — etwa beim nächsten Wellen-Register-Abgleich —, findet `slice-070` keiner Welle zugeordnet, ohne dass eine Datei das behauptet.
- **verifizierbar:** nein — kein Gate misst Wellen-Zuordnungen geschlossener Slices (`wave-drift` misst nur `in-progress`). Belegt per repo-weitem `grep` nach `welle-59` vor und nach dem Commit.
- **klasse:** einziger-fundort-mit-dem-geheilten-feld-entfallen

### F-7 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** AGENTS §3.7 §Bestandsgrenze gegen den Bestand / AGENTS §5 („Alt-Slices in `done/` behalten ihr historisches Feld") / `BEO-002` bzw. `MR-025` (Spiegel vor dem Editieren)
- **pfad:** `AGENTS.md:191-193`, `AGENTS.md:278-280`
- **befund:** §3.7 sagt unverändert: „**Für Zustandsfelder gibt es keine Bestandsgrenze:** der vorhandene Bestand wird mit dem v5.9.0-Bump umgestellt, nicht grandfathered." Der Slice lässt per Auftraggeber-Entscheid 50 chronikführende Felder stehen (F-1), und weder §3.7 noch §5 noch ein Carveout unter `docs/plan/carveouts/` hält diese Ausnahme fest; §5 spricht nur über die *Existenz* des Feldes, nicht über seinen Inhalt. Die Spiegel-Liste des Slice (§2 Schritt 5) nennt `AGENTS.md` §5 ausdrücklich, die Datei wurde in demselben Commit angefasst — der Nachzug unterblieb trotzdem. Versagen: Der nächste Review, der den Skill-Anker regelkonform auf den `done/`-Bestand anwendet, muss 50 Findings schreiben, weil das Briefing die Ausnahme nirgends trägt; oder ein Implementer räumt sie in gutem Glauben und fasst 50 eingefrorene Lauf-Belege an.
- **verifizierbar:** nein — kein Gate hält Briefing-Aussagen gegen den Bestand. Belegt per Zensus (F-1) gegen `AGENTS.md:191-193`.
- **klasse:** briefing-behauptet-einen-umgestellten-bestand-den-es-nicht-gibt

### F-8 · LOW

- **kategorie:** LOW
- **quelle:** AGENTS §3.7 (Zeitform-Test: Indikativ über das, was ist) / `BEO-002` (Ränder bleiben stehen) / Kanon §Was ein Kommentar trägt
- **pfad:** `docs/plan/planning/done/slice-070-trace-tabellenquellen-nullmengen-guard.md:5`, `docs/plan/planning/done/slice-071-trace-cross-consistency-gate.md:5-8`, `docs/plan/planning/done/slice-071-trace-cross-consistency-gate.md:19`
- **befund:** In zwei der dreizehn angefassten Dateien steht zwei Zeilen unter dem geheilten Feld ein zweites Kopffeld mit derselben falschen Gegenwarts-Aussage: `slice-070` trägt `**Welle:** aktiv; …`, `slice-071` trägt `**Welle:** … **inhaltlich weiter Teil der Welle**, aber nicht in Arbeit … Wiederaufnahme, sobald slice-074 eine tragende Regel hat` — `slice-074` liegt seit 2026-07-18 in `done/`. Im selben Kopf sagt `slice-071` außerdem `**Bezug:** … ADR-0038 (Proposed)`, während ADR-0038 `Accepted` ist. Der Slice hat die Klasse „Kopffeld widerspricht dem Ist" in genau diesen Dateien geheilt und die Nachbarzeile stehen lassen — der von `BEO-002` benannte Blindfleck „in genau der Datei, die ohnehin bearbeitet wurde".
- **verifizierbar:** nein — kein Gate liest `**Welle:**`- oder `**Bezug:**`-Kopffelder. Belegt per Lesung der beiden Dateiköpfe im HEAD-Stand.
- **klasse:** rand-zustandsfeld-in-der-gerade-editierten-datei

### F-9 · LOW

- **kategorie:** LOW
- **quelle:** `BEO-009` (Botschaft/Plan behauptet eine Änderung, die so nicht stattfand) / Maintainability
- **pfad:** `docs/plan/planning/in-progress/slice-121-zustandsfeld-hygiene.md:36`, `:52`, `:56-57`, `:74`, `:111`
- **befund:** Die Umklassifikation von „zwölf" auf „elf plus zwei" ist nur teilweise nachgezogen. Stehen geblieben sind: §1 „Die zwölf werden auf den wahren Zustand gesetzt" (es sind dreizehn), §2 Schritt 1 „die zwölf widersprüchlichen Felder namentlich auflisten", §2 Schritt 2 „(ein Fall sagt `done — abgeschlossen am …`)" (es waren zwei), §3 „Sie zu entfernen hieße, 78 eingefrorene Lauf-Belege anzufassen" — im selben Aufzählungspunkt, dessen Kopf im Commit auf 77 korrigiert wurde —, und §6 „das Heilen der zwölf Felder". Fünf Zahlen im Plan widersprechen den drei korrigierten.
- **verifizierbar:** nein — kein Gate liest Zahlwörter. Belegt per `grep` im HEAD-Stand der Slice-Datei.
- **klasse:** teilweise-nachgezogene-zahl-nach-umklassifikation

### F-10 · LOW

- **kategorie:** LOW
- **quelle:** AGENTS §1 („Diese Datei trägt **Hard Rules und Pointer** … Bei Konflikt … gilt die kanonische Quelle") / Kanon `.harness/baseline/v5.9.0/regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt
- **pfad:** `AGENTS.md:185-186`, `.harness/skills/reviewer.md:44-46`
- **befund:** Beide neuen Passagen enden mit dem Zeiger auf den Baseline-Abschnitt *Dieselbe Regel für Zustandsfelder* („Kanon:" bzw. die Klammer im Skill), und beide Male steht der Zeiger jetzt hinter der Vorrangregel. Der zitierte Abschnitt zählt `Stand`/`Status` einer Register-Zeile, die Stand-Zelle eines Roadmap-Fadens und die Status-Spalte eines Meilensteins auf; über ADR-Kopffelder, über Immutabilität und über einen Vorrang gegenüber einer ADR-Regel sagt er nichts. Der Rang-Zeiger deckt damit mehr, als die Quelle trägt, und die Vorrangregel selbst ist eine repo-lokale Entscheidung ohne eigenen Anker — kein `MR-*`, keine ADR, kein Carveout.
- **verifizierbar:** nein — Anker-Semantik, kein Gate (der Link löst auf, Probe 1 grün). Belegt per Lesung des zitierten Baseline-Abschnitts.
- **klasse:** rang-zeiger-deckt-mehr-als-die-quelle-traegt

### F-11 · INFO

- **kategorie:** INFO
- **quelle:** AGENTS §3.7 (Aufzählung der Zustandsfeld-Orte)
- **pfad:** `AGENTS.md:173-175`
- **befund:** §3.7 nennt als Zustandsfeld „etwa eine `Stand`- oder `Status`-Zelle in Roadmap, Beobachtungs-Register oder Meilenstein-Tabelle". Ein Kopffeld einer ADR oder eines Slice-Plans ist keines der drei. Dieser Slice wendet die Regel gleichwohl auf beide an — dreizehn Slice-Kopffelder werden geheilt, und die neue Vorrangregel setzt voraus, dass §3.7 sonst bis in den ADR-Kopf reichte. Entweder ist die Aufzählung wirklich offen — dann weitet die Vorrangregel §3.7 stillschweigend auf alle Dokument-Kopffelder aus, was in der Regel selbst nicht steht —, oder die aufgelöste Kollision bestand nach dem Wortlaut nie. Dokumentiert als undokumentierte Annahme; die Auflösung ist eine Redaktionsentscheidung, kein Fehler.
- **verifizierbar:** nein — Auslegungsfrage.
- **klasse:** regel-wird-weiter-angewandt-als-ihre-eigene-aufzaehlung-reicht

## Negativbefunde

1. **Die Widerspruchs-Klasse ist vollständig geräumt.** `grep` über alle 90 `**Status:**`-Kopffelder in `docs/plan/planning/done/` nach `open`, `next`, `in-progress`, `In Arbeit`, `geplant`, `offen`, `WIP`, `laufend`, `blockiert`, `Proposed`: null Treffer. Die elf geheilten Dateien sind exakt die elf, die der Quell-Befund (slice-119-Review F-5) namentlich nennt — `slice-025`, `slice-026`, `slice-027`, `slice-028`, `slice-083` (`open`), `slice-056`, `slice-057`, `slice-068`, `slice-070` (`in-progress`), `slice-084`, `slice-085` (`In Arbeit`). Kein Kandidat dieser Klasse übersehen.
2. **Keine Fläche außerhalb `done/`.** `docs/plan/planning/open/` und `docs/plan/planning/next/` sind leer; `docs/plan/planning/in-progress/` trägt nur `roadmap.md` und den laufenden Slice, und der führt korrekt den `**Lifecycle:**`-Hinweis statt eines Status-Feldes (Ziel-Form `slice.template.md`).
3. **Kein Planungs-Dokument jenseits der Slices trägt die Klasse.** Alle 37 Nicht-Slice-Dateien in `done/` geprüft (22 Ergebnisnotizen `welle-*-results.md`, 15 flache Wellendokumente): kein `**Status:**`- und kein `**Stand:**`-Kopffeld; die Wellendokumente deklarieren ausdrücklich „kein Status-Feld". Auch `docs/plan/planning/observations.md` und `docs/plan/planning/in-progress/roadmap.md` tragen nach welle-81 keine Kopf-Zustandszeile mehr.
4. **Kein Substanzverlust in `slice-071` und `slice-078`.** Jede Angabe der entfernten Absätze steht vollständig in der Closure-Notiz §7 derselben Datei: Abschlussdatum, Lifecycle-Kette, Blocker-Auflösung, ADR-Status, Release-Tags v0.44.0/v0.45.0 bzw. v0.49.0, Review-Verdikte — einschließlich des `17. Testarchitektur`-Belegs samt `spec/architecture.md`-Zeilenangabe und des `ACCEPT-WITH-NITS`-Verdikts. Auch der Lastenheft-/ADR-Kettennachweis von `slice-078` steht dort.
5. **Kein Substanzverlust in `slice-068` und `slice-083`.** Der entfernte Text von `slice-068` war vorwärtsgewandte Plan-Prosa („Implementierung + Review + Release folgen") plus ein inzwischen überholter ADR-Status (`Proposed`); Welle, ADR und Ergebnis stehen im `**Welle:**`-/`**Bezug:**`-Block und in §7. Der entfernte Text von `slice-083` („reine Ist-Messung … die Etappen gehören **vor** die Umsetzung abgenommen") steht inhaltlich doppelt: im `**Bezug:**`-Block („kein Lastenheft-/Spezifikations-Bump, keine neue Kennung, keine ADR, kein Release"), im Hinweis-Blockzitat darunter und in §7 („Reine Analyse — kein Code/Spec/Harness-Delta").
6. **In elf der dreizehn Dateien überlebt der Wellen-Name.** `slice-056`, `slice-057`, `slice-068`, `slice-084`, `slice-085` und die vier `slice-025`…`slice-028` tragen ihre Welle im `**Welle:**`-Block; `slice-083` ist ausdrücklich keiner Welle zugeordnet, `slice-071`/`slice-078` nennen sie in §7. Einzige Ausnahme ist `slice-070` (F-6).
7. **`version.md` trägt genau einen Anker.** `grep` nach `a id=`: ein Treffer, `version.md:35`, auf `v0.62.0` — der aktuellen Version. Die Kopfregel („**Nur die aktuelle Version** trägt einen expliziten HTML-Anker") ist damit erfüllt; der Regeltext selbst ist unverändert.
8. **Jeder lebende `version.md#…`-Verweis löst auf.** Repo-weites `grep`: der einzige feste Markdown-Pin auf eine Nummer ist `docs/user/benutzerhandbuch.md:3` auf `version.md#v0.62.0` — die Zeile, die den Anker behält. Alle übrigen Verweise zeigen auf `version.md#aktuell` (Heading-Slug von `## Aktuell`, unberührt) oder stehen in Inline-Code bzw. in eingefrorenen Reports (`docs/reviews/2026-07-01-slice-056-commits-doc-r1.md:157`, `docs/reviews/2026-07-17-slice-071-closure-independent.md:201`, `docs/reviews/2026-07-17-slice-075-implementation-r1.md:47` — alle drei innerhalb einfacher Backticks, kein Link). Auf den entfernten Anker der Vorgänger-Version zeigte kein Dokument; kein bestehender Verweis wird durch den Fix rot.
9. **Der Detektor ist scharf — belegt in beide Richtungen.** Probe 2 (HEAD-Kopie, `docs/user/benutzerhandbuch.md:3` auf den Vorgänger-Anker umgebogen): **Exit 1**, ein Befund `anchor-missing`. Probe 3 (dieselbe Kopie, `version.md` auf den Stand vor dem Commit zurückgesetzt, also mit zwei Ankern): **Exit 0**, 0 Befunde. Der Commit stellt genau die Meldung her, die der Registerkopf als seinen Zweck nennt.
10. **Gate-Stand des Baumes.** Probe 1 (Default-Profil): **Exit 0**, 437 Dateien, 0 Befunde — die Angabe der Commit-Botschaft ist reproduziert. Probe 4 (`--enable planning`): **Exit 0**. Probe 5 (`--config .d-check.closure.yml --enable planning --enable structure`): **Exit 0**, 403 Dateien. Die dreizehn geänderten `done/`-Slices verletzen weder eine Closure-Notiz-Bedingung noch eine Abschnitts-Invariante.
11. **Das Anfassen eingefrorener `done/`-Dokumente ist gedeckt, nicht überdehnt.** Der Kanon `modul-05-planning-harness.md` §Lifecycle als State Machine wendet sich gegen das Feld als **zweite Quelle** desselben Zustands, nicht gegen das Editieren von `done/`; keine Baseline- und keine Repo-Regel erklärt `done/`-Dateien für immutabel (das tut §3.5 nur für `Accepted`-ADRs), und §3.7 versagt Zustandsfeldern die Bestandsgrenze ausdrücklich. Der Kanon nennt „die Slice-Datei in `done/`" zudem selbst als legitimen Ort der Schließung — die Chronik gehört dort in die Closure-Notiz, nicht ins Zustandsfeld. Die Begründung des Slice (falsche Gegenwarts-Aussage, nicht Alter) trägt.
12. **Spiegel geprüft, kein weiterer Nachzug fällig.** `docs/user/releasing.md:29-30` trägt die Anker-Wanderungs-Regel bereits korrekt und musste nicht mit; `docs/plan/planning/README.md` §Lifecycle macht über Kopffelder keine Aussage; `harness/conventions.md:61` sagt „Der **Zustand ist die Verzeichnis-Position**, kein Status-Feld" und bleibt gültig; `.harness/skills/closure-note-reviewer.md` berührt Zustandsfelder nicht; die Ziel-Formen `slice.template.md` (`**Lifecycle:**` statt Status) und `NNNN-titel.template.md` (`**Status:** Proposed | Accepted | Deprecated | Superseded by ADR-NNNN`) sind unverändert konsistent. Einziger Spiegel mit Konflikt ist `harness/README.md:150-153` — als Bestandteil von F-2 gemeldet, nicht doppelt.
13. **Commit-Form.** Ein einzelner Commit ohne `git mv`; §3.3 (Move plus Inhalt = zwei Commits) ist nicht einschlägig. Die Botschaft nennt `slice-121` und erfüllt damit die Traceability-Form aus §5; sie deklariert die Umklassifikation von „zwölf" auf „elf plus zwei" offen und benennt die Gegenprobe samt Ergebnis — der `BEO-009`-Klasse wird hier aktiv begegnet (die Rest-Zahlen im Plan sind F-9, nicht die Botschaft).
14. **`version.md`-Tabelle bleibt monotonic.** Die geänderte Zeile behält Spalte 1 (`v0.61.0` in Inline-Code) und ihre Position; die `structure`-Invariante `table-order: desc` auf `## Verlauf` ist erfüllt (Probe 1 grün, `structure` ist im Default-Profil aktiv).
15. **Nur eine ADR ist betroffen.** Über alle 57 ADRs geprüft: `docs/plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md:3` ist das einzige `**Status:**`-Feld mit Chronik; alle übrigen tragen exakt `Accepted` oder `Proposed`. Die Zahl „genau eine ADR" aus §1 stimmt, und der Slice hat sie erwartungsgemäß nicht angefasst (§3 NICHT-Liste eingehalten).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 5 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** heil-auswahl-folgt-dem-suchstring-nicht-der-klasse · vorrangregel-behauptet-eine-immutabilitaet-die-das-gate-nicht-kennt · pronominal-referenz-kehrt-die-vorrangregel-um · skill-ausnahme-ohne-neuzugangs-vorbehalt · skill-inhalt-geaendert-version-eingefroren · einziger-fundort-mit-dem-geheilten-feld-entfallen · briefing-behauptet-einen-umgestellten-bestand-den-es-nicht-gibt · rand-zustandsfeld-in-der-gerade-editierten-datei · teilweise-nachgezogene-zahl-nach-umklassifikation · rang-zeiger-deckt-mehr-als-die-quelle-traegt · regel-wird-weiter-angewandt-als-ihre-eigene-aufzaehlung-reicht

## Ausgänge zu den Slice-§5-Risiken

- **Risiko 1 (ein `done/`-Slice ist ein Lauf-Beleg).** Der Grundsatz ist gedeckt (Negativbefund 11), die Ausführung nur teilweise: ein Beleg ist ersatzlos entfallen (F-6), und die Grenze, die §3 zieht, trennt die geheilten nicht von den verbliebenen Feldern (F-1).
- **Risiko 2 (der Anker-Fix schärft einen Detektor — lebt im Baum ein solcher Pin?).** Adressiert und belegt: kein lebender Verweis auf die Vorgänger-Version (Negativbefund 8), Schärfe in beide Richtungen nachgewiesen (Negativbefund 9). Ohne Vorbehalt erfüllt.
- **Risiko 3 (die Vorrangregel darf kein Schlupfloch werden).** Nicht erfüllt: die Regel begründet sich auf einer Immutabilität, die weder §3.5 noch das Gate kennen (F-2), ihr Neu-ADR-Vorbehalt ist im Briefing zweideutig (F-3) und im Skill gar nicht vorhanden (F-4) — der tote Winkel für neue ADRs, den §5 ausschließen wollte, ist offen.

## Verdikt

**Merge-blockierend:** ja — zwei HIGH und fünf MEDIUM, sämtlich auf Flächen, die dieser Slice selbst angefasst hat. F-2 ist der schwerste: die Regel, die eine Kollision zweier Hard Rules auflösen soll, beruft sich für ihre Begründung auf eine Aussage, die §3.5 im selben Dokument, der Spiegel in `harness/README.md` und die Gate-Konfiguration übereinstimmend ausschließen — sie setzt an die Stelle einer offenen Frage eine falsche Antwort über einen Gate-Pfad. F-1 ist der zweite: das angewandte Heil-Kriterium ist ein Suchstring, nicht die Klasse, und die DoD sagt „nachher null", wo 50 Felder weiter Chronik tragen. F-7 hängt an F-1 — solange der Bestand bewusst teilgeräumt bleibt, muss die Ausnahme im Briefing stehen, sonst schreibt sie der nächste Review als 50 Findings.

**Nicht blockierend, obwohl MEDIUM:** keines. F-5 und F-6 sind je für sich klein, liegen aber beide innerhalb des Slice-Scopes und sind mit einer Zeile behebbar.

**Unstrittig erfüllt:** der Anker-Teil (§1 Punkt 2) — genau ein Anker, Kopfregel eingehalten, kein Verweis gebrochen, Detektor-Schärfe in beide Richtungen am Image belegt — und die vollständige Räumung der Widerspruchs-Klasse (elf von elf, keiner übersehen, keiner außerhalb `done/`).

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler des Beobachtungs-Registers. Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill-Stand, dieses Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD- und Spec-Konformität prüft der Verifier separat.
