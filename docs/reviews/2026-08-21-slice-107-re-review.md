# Review-Report — slice-107 Stufen-Audit, bestätigende Re-Review (Heilung `c5c4b9d`)

- **Review-Art:** bestätigende Re-Review (Plan-/Audit-Review) — geprüft wird,
  ob die Heilung `c5c4b9d` die elf Befunde des Erst-Reviews schließt und ob
  sie neue Defekte einträgt.
- **Gegenstand:** Commit `c5c4b9d` (Audit
  `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md` §9,
  `docs/plan/planning/open/slice-108-roadmap-offene-wellen.md`,
  `docs/plan/planning/open/slice-109-v560-konventions-nachzuege.md`,
  Roadmap-Driftzeile) gegen den Erst-Report
  `docs/reviews/2026-08-21-slice-107-audit-review.md` (F-1…F-11).
- **Skill:** `.harness/skills/reviewer.md` @ 1.4.0.
- **Modell-ID:** `claude-fable-5`.
- **Datum:** 2026-08-21.
- **Eingangs-Kontext:** vendorter Baum `.harness/baseline/v5.6.0/regelwerk/`
  (26 Dateien); Stufen-Diffs und Tag-Notizen im Kurs-Repo (read-only:
  `git diff vX..vY -- lab/regelwerk`, `git tag -n99 v5.4.0`);
  `harness/conventions.md`, `AGENTS.md`, `harness/README.md`, `Makefile`,
  `.d-check.yml`, `spec/spezifikation.md` §DC-FA-ID-001.a,
  `spec/lastenheft.md`; eigener Gate-Lauf `make doc-check` (401/0, grün —
  deckt die Commit-Behauptung). Keine Host-Toolchain benutzt.

---

## Befund-Schließung (je Erst-Finding)

| Erst-Finding | Status | Beleg |
|---|---|---|
| F-1 (MEDIUM, fehlende Widerspruchs-Ausgang-Zeile) | **geschlossen** — mit Rest RF-1 | Zeile steht in der v5.3.1/v5.4.0-Tabelle; die drei Ausgänge (behalten + Widerspruch benennen / übernehmen = Rückbau / „will, kann noch nicht" = Carveout) decken Tag-Notiz v5.4.0 NEU 2 exakt; der Ausgang ist auf den Struktur-ID-Fall angewandt |
| F-2 (MEDIUM, falscher Stichproben-Kandidat) | **geschlossen** | `git diff v5.0.0..v5.6.0 --stat -- lab/regelwerk/modul-14-docker-harness.md` ist **leer**; von 26 Regelwerk-Dateien sind 20 geändert, die delta-freie Menge ist exakt die genannte (`grundlagen-klassifikation`, `modul-00/01/14/15/16`); slice-109 C-6 zieht nach |
| F-3 (MEDIUM, Spiegel-Liste unvollständig) | **geschlossen** | alle fünf Fundstellen verifiziert: `AGENTS.md:99-100` (§3.3), `AGENTS.md:144` (§4), `harness/README.md:60`, `Makefile:81`, `harness/conventions/MR-013-lifecycle-move-buendelung.md:15,20`; eigener grep über beide Phrasen findet **keinen** weiteren lebenden Träger außerhalb der slice-108-Schritte (Roadmap → Schritt 1, `.d-check.yml:247` → Schritt 2, MR-024-Datei + `harness/conventions.md:102`-Indexzeile → Schritt 3; ADR-0028 immutable, Produkt-Default-Doku per Schritt 5 ausgenommen, `done/`/`docs/reviews/`/`CHANGELOG.md` eingefroren). „plus alles, was der grep noch findet" ist als Auftrag hinreichend präzise: beide Nadeln stehen wörtlich zitiert in der Liste, grep-vor-Edit ist ausdrücklich angeordnet, und die beiden Ausschluss-Grenzen (Produkt-Default, eingefrorene Artefakte) trägt der Slice selbst (Schritt 5) bzw. die stehende Immutabilitäts-Konvention |
| F-4 (MEDIUM, überdehnter Rückfallweg) | **geschlossen** | die MR-027-Linie (Abweichung → deklarierte Adaption, nicht still in MR-000) ist konsistent durch `slice-107:118` (v5.2.0-Zeile), `:143` (Widerspruchs-Ausgang „in diesem Audit angewandt"), `:183` (Etappe-C-Schnitt) und slice-109 C-2(b) gezogen — dort mit Begründung **und** Auflösungs-Trigger; MR-027 ist die nächste dichte Nummer nach MR-026. Rest: eine falsche Zahl in der Begründung → RF-2 |
| F-5 (MEDIUM, fehlende Vergabe-/Dissens-Zeilen) | **geschlossen** — mit Resten RF-3/RF-5 | beide Zeilen tragen die tragenden Regeln des v5.4.0..v5.5.0-Diffs: „Welle zählt repo-weit dicht" + „lokal ableitbar liest Verzeichnis **und** offene Welle-Dateien (§4-Zeilen ohne Datei)" wörtlich in `grundlagen-source-precedence.md`; die modul-10-Dissens-Regel („Skill schärfen, nie die mildere Lesart") wörtlich im modul-10-Diff |
| F-6 (LOW, Emissions-Regel ohne Zeile) | **geschlossen** | die Regel steht wörtlich im v5.2.0..v5.3.0-Diff von `modul-09-implementierung.md` („Emittierte Skripte, Makefiles und Konfigurationen tragen keine Slice- oder Befund-Nummer des Erzeuger-Repos"); Zeile korrekt in der v5.3.0-Tabelle, Befund referenziert die Gegenprobe des Erst-Reviews |
| F-7 (LOW, zwei v5.2.0-Kandidaten ohne Zeile) | **geschlossen** | beide Zeilen decken die v5.1.0..v5.2.0-Diffs: `grundlagen-bootstrap` „Beide Modi regeln dieselbe Achse … zuständig ist der Freshness-Audit" (Migration ≠ BF — als konform-gelebt richtig), `grundlagen-traceability` „die `ARC-*` der Sicht zählt nicht" (als geänderte Regel benannt, n. a. korrekt begründet: keine `ARC-*` im Repo) |
| F-8 (LOW, Zahl ohne Zählung im Drift-Log) | **geschlossen** | `roadmap.md:98` führt keine Belegzählung mehr („Konform-Belege ohne Handlung je Zeile (u. a. …)"); „sechs Stufen" und „C-2…C-6" sind konsistente, nachzählbare Größen |
| F-9 (LOW, konform ohne Stichprobe) | **geschlossen** — mit Rest RF-4 | die v5.3.0-Zeile trägt die gemessene Rest-Klasse samt C-3-Auftrag; slice-109 C-3 trägt das grep-Mandat samt Fundstellen. Verifiziert: `cli.go:331,438`, `planning.go:116`, `configyaml.go:1721,1731`, `structure_tableorder.go:71,178`, `markdown.go:218` — alle real |
| F-10 (INFO, Produkt-Default-Frage nirgends notiert) | **geschlossen** | slice-108 Schritt 5 benennt die Grenze (Default bleibt `## Aktuelle Welle`, Änderung = Konsumenten-Breaking-Change, eigener CR/ADR-Entscheid) |
| F-11 (INFO, Schutz-Aussage netzlos unbelegt) | **geschlossen** | `slice-107:163` behauptet nur noch das netzlos Belegbare (gelebter Direkt-Commit-Träger; Schutzregime ausdrücklich „netzlos nicht verifizierbar") |

## Neue Findings dieses Laufs

### RF-1

- **kategorie:** LOW
- **quelle:** slice-107 §1 („je Stufe gegen die Tag-Notizen")
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:143`
- **befund:** Die Regel-Zelle der Widerspruchs-Ausgang-Zeile mischt zwei
  Baseline-Passagen zweier Stufen: „ein Fund geht den Weg jeder Diskrepanz —
  übernehmen im nächsten Slice … mehrere Funde treffen die MR-000-Aussage"
  ist wörtlich der **v5.2.0**-Bestands-Stichproben-Absatz von `modul-02`
  (dort bereits von der v5.2.0-Zeile des Audits gedeckt); die
  **v5.4.0**-Regel-Ergänzung (Tag-Notiz NEU 2) ist allein die Wahl am
  Ausgang 5 der `MR`-Prüfung (behalten-mit-benanntem-Widerspruch /
  übernehmen als Rückbau / Carveout bei „will, kann noch nicht"). Die drei
  Ausgänge der Zeile und die Anwendung auf den Struktur-ID-Fall sind
  korrekt — nur die Stufen-Zuordnung der Rahmensätze nicht; ein Leser, der
  die Zeile gegen den v5.4.0-Diff nachliest, findet zwei ihrer Klauseln dort
  nicht.
- **verifizierbar:** nein (Lese-Abgleich `git diff v5.3.1..v5.4.0` vs. `git diff v5.1.0..v5.2.0 -- lab/regelwerk/modul-02-harness-bootstrap.md`)
- **klasse:** stufen-attribution-vermischt

### RF-2

- **kategorie:** LOW
- **quelle:** Maintainability (ehrlich zählen — exakt die in `c5c4b9d` selbst geheilte F-8-Klasse)
- **pfad:** `docs/plan/planning/open/slice-109-v560-konventions-nachzuege.md:26-27`
- **befund:** Die MR-027-Begründung behauptet „42 Anforderungen mit je
  eigener `.a`-Sektion". Nachgezählt: das Lastenheft führt **44**
  `DC-FA-*`-Anforderungen (plus 4 `DC-QA-*`), die Spezifikation §1 führt
  `.a`-Sektionen für **36** davon — acht FA-Anforderungen haben keine
  eigene `.a`-Sektion (`DC-FA-CLI-003/004`, `DC-FA-CONF-001`,
  `DC-FA-DIST-001`, `DC-FA-LINK-002`, `DC-FA-MTX-002/003`,
  `DC-FA-SCAN-001`). Weder die Zahl noch die Je-Aussage stimmt; unkorrigiert
  wandert beides bei C-2 in die Begründung der `MR-027`-Datei.
- **verifizierbar:** nein (Zählung per grep über `spec/lastenheft.md`-Headings und `### DC-`-Sektionen der Spezifikation; kein Gate-Gegenstand)
- **klasse:** zahl-ohne-zählung

### RF-3

- **kategorie:** LOW
- **quelle:** Maintainability (Retro-Belege ehrlich — slice-095-Lehre 15/19)
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:155`
- **befund:** Der Beleg der Vergabe-Zeile — „in dieser Welle **belegbar**:
  welle-78 §4 vergab slice-108/109 vor Datei-Existenz" — ist aus der Historie
  nicht belegbar: die §4-Zeilen und beide Slice-Dateien entstanden im
  **selben** Commit `a0e6ecc` (bei Eröffnung `0b9431f` und B-Start `37e1fd4`
  führt §4 nur den Platzhalter „Etappe-C-Slices folgen aus dem B-Befund").
  Es gab in keinem committeten Stand eine §4-Zeile ohne Datei; „belegbar"
  überzieht einen Arbeits-Reihenfolge-Anspruch, den kein Artefakt trägt.
  (Der Befund stammt aus dem Erst-Report F-5 und wurde ungeprüft als Beleg
  übernommen; die konform-Antwort der Zeile selbst bleibt haltbar.)
- **verifizierbar:** nein (`git log --follow` über das Wellendokument)
- **klasse:** beleg-nicht-nachvollziehbar

### RF-4

- **kategorie:** LOW
- **quelle:** Maintainability (Fundstellen ehrlich benennen)
- **pfad:** `docs/plan/planning/open/slice-109-v560-konventions-nachzuege.md:40`; `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:131`
- **befund:** Zwei Ungenauigkeiten in den F-9-Fundstellen-Nennungen:
  (a) slice-109 C-3 nennt `finding.go` als Träger nackter
  Review-Finding-Tokens — `internal/hexagon/core/model/finding.go` enthält
  **keinen** (grep über `R?-F-?`/`Review`-Muster leer); die übrigen fünf
  genannten Dateien tragen reale Fundstellen. (b) Die Audit-Zeile nennt
  `structure_tableorder.go`/`markdown.go` als „die in **dieser** Welle
  selbst geschriebenen" — beide Kommentare (`structure_tableorder.go:71,178`,
  `markdown.go:218`) stammen aus welle-77 (slice-105/ADR-0057), nicht
  welle-78; slice-109 sagt an derselben Stelle korrekt „welle-77/78".
  Folgenlos für die Räumung (das grep-Mandat ist der Auftrag, nicht die
  Liste), aber die Liste behauptet eine Fundstelle, die es nicht gibt.
- **verifizierbar:** nein (grep über die genannten Dateien)
- **klasse:** fundstelle-behauptung-falsch

### RF-5

- **kategorie:** INFO
- **quelle:** slice-107 §4 DoD („je Regel eine Zeile")
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:155`
- **befund:** Vom v5.5.0-§Vergabe-Block bleiben zwei im Erst-Report F-5
  ausdrücklich genannte Teilstücke ohne eigene Zeile: die
  `MR`-Hybrid-Kollisionsklasse und die Abgrenzung Nummern-Ableitung vs.
  Beanspruchen. Beide sind im Ein-Schreiber-Betrieb ohne Failure-Szenario
  (die Hybrid-Kollision braucht zwei gleichzeitige Schreiber; die
  Abgrenzung ist eine Nicht-Widerspruchs-Klärung, deren Beanspruchen-Hälfte
  die TA-7-Zeile deckt) — darum kein LOW, aber die Sammelzeile
  „§Vergabe-Ergänzungen" heißt vollständiger, als sie ist.
- **verifizierbar:** nein
- **klasse:** aufzählung-rest-ohne-zeile

---

## Negativbefunde (geprüft, ohne Befund)

- **F-1-Substanz:** die drei Ausgänge der neuen Zeile decken Tag-Notiz
  v5.4.0 NEU 2 vollständig (Wahl behalten/übernehmen; Carveout-Abgrenzung);
  die Anwendung auf den Struktur-ID-Fall ist unter **beiden** einschlägigen
  Baseline-Passagen tragfähig (v5.2.0: „`MR`-Einträge, die die Abweichungen
  deklarieren"; v5.4.0: deklarierte Abweichung statt stiller Nicht-Befolgung).
- **F-2-Nachrechnung:** `git diff v5.0.0..v5.6.0 --stat` über den
  Regelwerk-Baum: 20 von 26 Dateien geändert; die sechs delta-freien sind
  exakt die im Audit genannten. `modul-07-carveouts.md` ist mit +13/−2 in
  v5.2.0 korrekt als gebrannt aussortiert.
- **F-3-Gegen-grep:** beide Phrasen über den gesamten Arbeitsbaum gegrept;
  jeder Treffer ist entweder eine der fünf gelisteten Fundstellen, von
  slice-108 Schritt 1–3 gedeckt (`roadmap.md`, `.d-check.yml:247`,
  MR-024-Datei samt `conventions.md:102`-Indexzeile), eingefroren
  (`done/`, `docs/reviews/`, `CHANGELOG.md`, ADR-0028) oder
  Produkt-Default-Doku (Spec, Handbuch, READMEs, `config.go`,
  `config_template.go`, Tests — bleibt per Schritt 5 unangetastet). Kein
  lebender Träger rottet still.
- **F-4-Marker:** alle vier MR-027-Zeilen tragen `d-check:ignore`; das
  Ventil ist per Spezifikation (§DC-FA-ID-001.a, Ventil 4) ganzzeilig und
  wirkt auf `ids` (`MR-\d{3}`, `link-policy: always`, Target
  `harness/conventions/` — MR-027-Datei existiert planmäßig noch nicht).
  Die In-Zellen-Platzierung bricht den Tabellen-Reader nicht (die
  Direktiven-Toleranz betrifft nur überzählige Zellen hinter der letzten
  Pipe) und entspricht der etablierten Praxis der MR-026-Marker aus
  slice-106/welle-78 §4. Kein weiteres verdecktes Suppress-Risiko auf den
  vier Zeilen gefunden.
- **F-5-Quellen:** beide Zeilen gegen den Kurs-Diff v5.4.0..v5.5.0
  (`grundlagen-source-precedence.md` +27-Zeilen-Block, `modul-10` Dissens +
  Auflösungs-Trigger-Disziplin) wörtlich gegengelesen — die Paraphrasen
  tragen.
- **F-9-Stichprobe:** alle fünf verifizierbaren Fundstellen-Nennungen real
  (einzige Ausnahme `finding.go` → RF-4); das grep-Mandat („alle
  Fundstellen … ziehen oder streichen; Verhalten unverändert") ist als
  Auftrag vollständig.
- **Etappe-C-Konsistenz nach Heilung:** slice-108/109 tragen keine
  Behauptung, die das geheilte Audit nicht deckt; Trigger-Kette
  (107 → 108 → 109) unverändert schlüssig; welle-78 §4 und Driftzeile
  konsistent.
- **Gate-Lauf:** eigener `make doc-check` grün — „401 Datei(en) geprüft,
  0 Befund(e)" — identisch zur Commit-Behauptung; die Heilung erzeugt
  keinen Gate-Befund.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 4 | RF-1, RF-2, RF-3, RF-4 |
| INFO | 1 | RF-5 |

**Finding-Klassen dieses Laufs:** stufen-attribution-vermischt ·
zahl-ohne-zählung · beleg-nicht-nachvollziehbar ·
fundstelle-behauptung-falsch · aufzählung-rest-ohne-zeile.

## Verdikt

**APPROVE mit Auflagen (LOW, nicht blockierend).** Alle fünf blockierenden
MEDIUMs des Erst-Reviews sind in der Sache geschlossen — die nachgetragenen
Zeilen tragen die richtigen Regeln mit gegen Kurs-Diff und Tag-Notizen
haltbaren Antworten, der Stichproben-Kandidat ist verifiziert delta-frei,
die Spiegel-Liste ist vollständig und als grep-Auftrag präzise genug, die
MR-027-Einordnung ist konsistent durchgezogen, und die LOW-/INFO-Nachzüge
(F-6–F-11) sitzen. Etappe C kann starten. Auflagen — im Zuge der ohnehin
anstehenden Slices, ohne erneuten Review-Lauf: RF-2 **vor** der
MR-027-Formulierung in slice-109 korrigieren (sonst wandert die falsche
Zählung in einen dauerhaften Konventions-Eintrag), RF-4(a) beim
C-3-grep miterledigen; RF-1/RF-3/RF-5 sind Doku-Präzisierungen im
Audit-Dokument nach Ermessen.
