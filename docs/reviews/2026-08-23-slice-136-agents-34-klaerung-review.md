# Review-Report: slice-136 — §3.4-Klärung — 2026-08-23

**Review-Art:** Code-Review (Doku-/Konventions-Diff gegen Kanon und Slice-Plan,
Modul 10 §Drei Review-Arten) · **Gegenstand:** Commit `830a02f` (`HEAD`, Diff
`HEAD~1..HEAD`) — der Feature-Commit von slice-136 (`feat(harness): §3.4
gemessen — keine Doppelung, keine Verschärfung, aber eine falsche Ausweisung
(slice-136)`). Der Commit ändert ausschließlich `AGENTS.md` §3.4 (1 Datei,
8 Einfügungen, 3 Löschungen laut `git show --stat`); die tragende Arbeit liegt
in der Commit-Botschaft.

**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 (`9a7654a`) · **Modell-ID:**
`claude-sonnet-5` · **Datum:** 2026-08-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-136-agents-34-klaerung.md`
  (§1–§9, insbesondere §4 DoD und §5 Risiken)
- `AGENTS.md` §3.4 (Diff) und §3.1–§3.9 (Überblick)
- Baseline-Kanon `modul-09-implementierung.md` §AGENTS.md-Regeln, §Hard Rules
  (repo-spezifisch), §Regeln gegen typische Fehlannahmen
- Baseline-Kanon `modul-03-spec.md` §Ziel-Form: Architektur-Sicht
- Baseline-Kanon `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)
- Baseline-Kanon `modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin)
- Vendorte Vorlage `.harness/baseline/v5.11.0/templates/spec/architecture.template.md`
  (vollständig gelesen) und `.harness/baseline/v5.11.0/templates/AGENTS.template.md`
  §3.4 (selbst aufgefunden, nicht im Auftrag benannt, aber direkt einschlägig)
- `spec/architecture.md`, `spec/spezifikation.md` (Volltext-Zählungen)
- Zensus `docs/plan/planning/done/slice-132-hard-rule-zensus.md` §9
  (Zensus-Tabelle, insbesondere die §3.4-Zeile)
- `.d-check.yml` (Module `matrix`, `ids`, `codepaths`)
- `docs/plan/planning/observations.md` (`BEO-011`, `BEO-012`)
- Vorheriger Review-Report am selben Modul:
  `docs/reviews/2026-08-23-slice-132-hard-rule-zensus-review.md` (F-1: „Das
  GEDECKT-Verdikt für §3.4 deckt nur eine von zwei Teilaussagen" — Kalibrierungs-Anker)

**Nicht erhalten:** die DoD-Abhakung (Verifikations-Rolle, getrennter Kontext;
alle vier DoD-Punkte in slice-136.md stehen noch auf `[ ]`, §9 Closure-Notiz
ist noch leer — konsistent mit dem Commit-Grenzen-Muster feat→closure-body).

**Vom Reviewer selbst gefahren** (nur Lesekommandos, kein `make`, keine
Dateiänderung): `git show`/`git diff` auf `HEAD~1..HEAD`; `grep` über
`.harness/baseline/v5.11.0/` und das ganze Repo nach `Modul-Pfad`; eigene
Volltext-Zählung von Pfad-artigen Inline-Code-Spans in
`architecture.template.md`, `spec/architecture.md` und `spec/spezifikation.md`
gegen mehrere Suchmuster (breiter Backtick-Slash-Regex, `internal/`/`cmd/`-Muster,
`tools/*.sh`-Muster); Abgleich der `matrix`-/`ids`-/`codepaths`-Konfiguration in
`.d-check.yml` gegen die im Regeltext genannten fünf Abwärts-Kategorien
(ADRs, Wellen, Slices, Commit-Hashes, Closure-Daten).

**Verdikt: blockierend** — ein HIGH, vier MEDIUM.

---

## Findings

### F-1 — Die „gedeckt"-Hälfte der Abwärts-Sperre wird auf eine Probe von fünf Kategorien verallgemeinert

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §3.4 Satz „Kein Spec-Stratum … referenziert ADRs,
  Wellen, Slices, Commit-Hashes oder Closure-Daten" plus „**Zwei Aussagen,
  eine davon gedeckt:** die **Abwärts-Sperre** hält `make doc-check`" ·
  `.d-check.yml:184–205` (Modul `matrix`) · `.d-check.yml:119–149` (Modul
  `ids`) · Zensus-Zeile `docs/plan/planning/done/slice-132-hard-rule-zensus.md:145`
  · `AGENTS.md` §5 Botschafts-Regel („ihr Schluss reicht nicht weiter als die
  gemessene Menge") · `BEO-011`/`BEO-009` · Kalibrierungs-Präzedenz: F-1 in
  `docs/reviews/2026-08-23-slice-132-hard-rule-zensus-review.md` — exakt
  dasselbe Muster eine Ebene höher (dort: „gedeckt" deckte nur eine von zwei
  Teilaussagen der Regel; wurde als HIGH gewertet und über die Zensus-Revision
  auf „teilgedeckt" korrigiert).
- **pfad:** `AGENTS.md:166–168` (unverändert von diesem Commit übernommen)
- **befund:** Die als „gedeckt" geführte Abwärts-Sperre nennt fünf
  Referenz-Kategorien (ADRs, Wellen, Slices, Commit-Hashes, Closure-Daten),
  belegt aber laut Zensus nur **eine** mit einem konstruierten Verstoß
  (`slice-999` in `spec/spezifikation.md` ⇒ `matrix-forbidden`). Die
  `matrix`-Konfiguration (`.d-check.yml:184–205`) kennt genau drei Klassen
  (`spec-straten`, `adr`, `slice`) und genau ein Token-Muster
  (`slice: token: 'slice-\d{3}'`); für „Wellen" existiert **keine** Klasse
  überhaupt — ein Link oder Token, das auf `docs/plan/planning/welle-*.md`
  zeigt, matcht keine `rules`-Zeile und bleibt ungeprüft. Für „Commit-Hashes"
  und „Closure-Daten" existiert weder in `matrix` noch in `ids`
  (`.d-check.yml:119–149`, Muster nur für `ADR-\d{4}`, `MR-\d{3}`, `DC-*`,
  `SPEC-\d{3}`, `ARC-\d{3}`) irgendein Muster. Selbst „ADRs" ist in `matrix`
  nur über **Markdown-Links** erfasst (die `adr`-Klasse trägt kein `token:`,
  anders als `slice`); ein bloßer Fließtext-Verweis „ADR-0005" ohne Link wäre
  für `matrix` unsichtbar (er würde nur indirekt über `ids`s
  `link-policy: always` auffallen, und auch das nur, wenn er nicht per
  `d-check:ignore` ausgenommen ist). Drei der fünf benannten Kategorien
  (Wellen, Commit-Hashes, Closure-Daten) haben damit **keinen** Sensor,
  obwohl der Satz „eine davon gedeckt" (im Gegensatz zur jetzt sauber
  aufgeteilten Sprachfreiheits-Hälfte) keine weitere Unterteilung trägt und
  uneingeschränkt liest. Der Commit korrigiert genau diese Klasse von Fehler
  für die zweite Hälfte des Abschnitts (Rollen/Modul-Pfad), lässt die
  identische Lücke in der ersten Hälfte aber unangetastet stehen — obwohl
  sein eigener DoD-Anspruch („§3.4 trägt keine Aussage mehr, die eine
  Unmöglichkeit behauptet, ohne sie geprüft zu haben") wortgleich dafür
  gälte.
- **verifizierbar:** ja — ein konstruierter Verstoß mit einer Welle-Referenz
  (Link auf eine `welle-*.md`-Datei) oder einem Commit-Hash im Prosa-Körper
  von `spec/spezifikation.md` würde `make doc-check` **nicht** rot färben;
  nicht selbst ausgeführt (Auftrag untersagt `make`-Läufe).
- **klasse:** gedeckt-verdikt-deckt-nur-eine-von-fuenf-kategorien

### F-2 — Die „NULL Code-Pfad-Token"-Zählung übersieht den Bedienhinweis der Vorlage

- **kategorie:** MEDIUM
- **quelle:** `.harness/baseline/v5.11.0/templates/spec/architecture.template.md:1–8`
  · Slice-Plan §Abnahme-Punkte („Zähle selbst. … auch in Bedienhinweisen, auch
  in `<…>`-Platzhaltern") · `BEO-011` Ausprägung (a), Exklusivitäts-Aussage
- **pfad:** Commit-Botschaft, Absatz „VERSCHÄRFUNG? NEIN" (Zeile
  „architecture.template.md führt NULL Code-Pfad-Token")
- **befund:** Der Template-Hinweis-Block der Vorlage (Zeilen 1–8, per
  eigener Anweisung „lösche diesen Block" vor der Nutzung zu entfernen)
  enthält den Satz „Sie ist **sprach- und meilensteinfrei** (siehe
  [Baseline-Regelwerk Modul 4](../../regelwerk/modul-04-adrs.md))" — ein
  relativer Pfad, der wörtlich auf ein „Modul" zeigt. Unter der im Auftrag
  genannten Lesart „Modul-Pfad = Pfad zu einem Modul" ist das die
  buchstäblichste mögliche Instanz eines Modul-Pfads in der gesamten Datei.
  Der Commit zählt „NULL", ohne diesen Fund zu nennen oder den Ausschluss zu
  begründen (z. B. „zählt nicht, weil der Block vor Gebrauch gelöscht wird").
  Die Auslassung ist für die Kern-Aussage („die Vorlage praktiziert Null")
  nicht folgenlos: Eine uneingeschränkte Exklusivitäts-Aussage über eine
  ganze Datei ist genau die Form, die laut `BEO-011` wiederholt gekippt ist.
- **verifizierbar:** ja — Volltext-Lesung der ersten acht Zeilen der Vorlage.
- **klasse:** null-zaehlung-uebersieht-bedienhinweis

### F-3 — Der Verschärfungs-Schluss stützt sich nur auf Vorlagen-Praxis, nicht auf die direktere Text-Evidenz, die ihm widerspricht

- **kategorie:** MEDIUM
- **quelle:** `modul-03-spec.md:126–127` („referenziert Modul-Pfade, aber
  keine Wellen, Slices, Commit-Hashes oder Closure-Daten") ·
  `.harness/baseline/v5.11.0/templates/AGENTS.template.md:128–138` (§3.4,
  wörtlich: „`spec/architecture.md` referenziert Modul-Pfade, aber **keine**
  Wellen, Slices, Commit-Hashes oder Closure-Daten") · Slice-Plan §5
  („Eine Vorlage ist kein Regeltext. Was `architecture.template.md`
  praktiziert, ist starkes Indiz und keine Definition.") · `BEO-012`
- **pfad:** Commit-Botschaft, Absatz „VERSCHÄRFUNG? NEIN" und „GRENZE DIESES
  BELEGS"
- **befund:** Der Commit zieht seinen Verschärfungs-Beleg ausschließlich aus
  der **Praxis** der Vorlagen-Datei (Token-Zählung in
  `architecture.template.md`) und nennt die Grenze dieses Belegs selbst
  („eine Vorlage ist kein Regeltext"). Er prüft dabei nicht die direktere
  **Text**-Evidenz: `AGENTS.template.md` §3.4 — das eigentliche
  AGENTS.md-Analogon der Baseline, nicht bloß eine Sicht-Vorlage — wiederholt
  denselben Satz wörtlich als ausgefülltes Rule-Beispiel, nicht als
  Platzhalter-Kommentar. Beide Quellen (`modul-03`s Bullet und
  `AGENTS.template.md`s Beispieltext) behaupten explizit, dass die Sicht
  „Modul-Pfade" referenziert — als Aussage über erlaubtes Verhalten, nicht
  nur als beobachtete Form. Diese Quelle ist um eine Stufe direkter als die
  Vorlagen-Praxis, die der Commit stattdessen zählt, und sie zeigt in die
  Gegenrichtung des gezogenen Schlusses. Der Commit widerlegt sie nicht,
  erwähnt sie nicht und stützt „kein MR-Eintrag fällig" allein auf das
  schwächere Indiz. Das ist die Form, vor der `BEO-012` warnt: der Text der
  herangezogenen Quelle (Vorlage) stimmt, ihre Reichweite als Beleg für die
  Regel-Bedeutung wird überdehnt, während eine belastbarere, gegenläufige
  Quelle ungenannt bleibt. `grundlagen-referenz-richtung.md` trägt an keiner
  Stelle eine Auflösung des Begriffs „Modul-Pfad" (geprüft, kein Treffer) —
  die Mehrdeutigkeit bleibt damit im Kanon offen, unabhängig vom
  Commit-Ergebnis.
- **verifizierbar:** nein — semantische Auslegungsfrage, kein Gate.
- **klasse:** vorlagen-praxis-statt-direkterer-textbeleg

### F-4 — Die „fünf"-Messung in der Spezifikation ist nicht reproduzierbar und mischt zwei verschiedene Pfad-Kategorien

- **kategorie:** MEDIUM
- **quelle:** Slice-Plan §2 Schritt 4 („heute gemessen: `spec/architecture.md`
  null, `spec/spezifikation.md` fünf") · `AGENTS.md:162` eigener Gebrauch von
  „Modul-Pfade (Import-Regeln)" im Sinn von ADR-0005 · `spec/architecture.md:29,91`
  (dieselbe Verwendung) · `.d-check.yml:226–240` (Modul `codepaths`,
  `ignore-refs`-Tombstone-Liste) · `BEO-011`
- **pfad:** Commit-Botschaft, Zeile „Gemessen: null in der Sicht, fünf in der
  Spezifikation (dort legitim)"
- **befund:** Der Commit nennt keine Suchmethode. Ein eigener Nachvollzug mit
  einem breiten Backtick-Slash-Muster über `spec/spezifikation.md` liefert
  weit mehr als fünf Treffer (u. a. `docs/plan/adr/`, `docs/plan/planning/`,
  `spec/lastenheft.md`, `gopkg.in/yaml.v3`), von denen alle unter
  `codepaths.roots` (`docs, spec, tools, harness, internal, cmd`) fallen und
  damit ebenso „von `codepaths` gefunden" würden. Genau fünf Treffer ergibt
  nur ein enger auf `tools/*.sh` beschränktes Muster
  (`spec/spezifikation.md:2277,2305,2842,2845,2846` — drei ×
  `tools/gate-consistency.sh`, je eine × `tools/planning-consistency.sh` und
  `tools/trace-check.sh`); diese fünf sind laut `.d-check.yml:234–240`
  bereits als Tombstone-Referenzen auf **abgelöste** Skripte deklariert und
  über `codepaths.ignore-refs` von der Existenzprüfung ausgenommen. Ein Muster
  für „Modul-Pfad" im Sinn, den `AGENTS.md` selbst an anderer Stelle benutzt
  (§3.4 Satz 3: „Modul-Pfade, Import-Regeln" = die ADR-0005-Hexagon-Paketpfade,
  z. B. `internal/model`, `internal/rules`, `internal/app`, `cmd/…`), liefert
  in `spec/spezifikation.md` dagegen **null** Treffer — denselben Wert wie in
  der Sicht. Die Größe „fünf" hängt damit vollständig vom gewählten
  Suchmuster ab, und unter der Lesart, die `AGENTS.md` selbst für „Modul-Pfad"
  benutzt, lautet der Vergleich nicht 0-zu-5, sondern 0-zu-0 — was die neu
  eingeführte Formulierung „kein heutiges Modul trägt es" zwar unverändert
  stützt, die Beispiel-Zahl im Commit aber unbelegt lässt.
- **verifizierbar:** ja — Musterabhängigkeit durch Nachzählen mit
  unterschiedlichen Regex direkt reproduzierbar.
- **klasse:** suchmuster-abhaengige-kennzahl

### F-5 — Der letzte Satz von §3.4 trägt drei unausgewiesene „wo lebt was"-Aussagen

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md:162–164` („Die sprachkonkrete Übersetzung
  (Modul-Pfade, Import-Regeln) und die Begründungen leben in den ADRs, deren
  `Schärft:`-Feld aufwärts zeigt; die zeitliche Schicht lebt in
  `docs/plan/planning/`.") · `modul-13-quality-gates.md:51–53` („Eine neue
  Hard Rule trägt ab ihrer Einführung einen Auflösungs-Trigger oder die
  Kennzeichnung *permanent*") · slice-127-Lehre („ein Abschnitt ist keine
  Regel") · Slice-Plan §4 DoD Punkt 3
- **pfad:** `AGENTS.md:162–164`
- **befund:** Dieser Satz enthält drei eigenständige Behauptungen — (1) die
  sprachkonkrete Übersetzung lebt in den ADRs, (2) deren `Schärft:`-Feld
  zeigt aufwärts, (3) die zeitliche Schicht lebt in
  `docs/plan/planning/`. Keine davon fällt unter die beiden neu eingeführten,
  sauber ausgewiesenen Kategorien („Rollen statt Technologie" — Urteil,
  permanent; „Modul-Pfad" — detektierbar, auflösender Trigger). Für keine der
  drei Aussagen nennt der Abschnitt einen Gate-Bezug, ein Urteil-mit-Trigger
  oder eine explizite Out-of-Scope-Notiz — sie stehen unkommentiert neben den
  jetzt fein aufgeteilten Aussagen, obwohl der Rest von `AGENTS.md` §3
  durchgängig jede Aussage entweder mit einem Gate-Verweis oder einem
  `(Auflösungs-Trigger: …)`-Feld versieht (vgl. §3.2, §3.6, §3.9). Der Commit
  behauptet mit seiner Überschrift („§3.4 gemessen") und seinem DoD-Anspruch
  eine vollständige Klärung des Abschnitts; dieser Satz fällt durch das
  Raster, das die beiden anderen Aussagen jetzt sauber trägt.
- **verifizierbar:** nein — Vollständigkeits-Urteil, kein Gate.
- **klasse:** unausgewiesener-restsatz-im-abschnitt

## Negativbefunde

- geprüft, ohne Befund: **Doppelungs-Frage** (Auftrags-Punkt 1) — das
  `modul-09`-Zitat („Jede Hard Rule liegt in zwei Quadranten … Beides ist
  Pflicht") ist korrekt und vollständig wiedergegeben
  (`modul-09-implementierung.md:170–175,207`) und trägt die Schlussfolgerung
  „keine Doppelung". `modul-09` nennt sogar noch direkter, ungenutzt vom
  Commit: „Nicht jede Hard Rule in AGENTS.md ist repo-spezifisch: ‚Architektur
  ist sprach- und meilensteinfrei' folgt aus dem Sicht-Stratum … ein Repo
  verkörpert sie in AGENTS.md, es entscheidet sie nicht. Solche Regeln
  gehören in die Datei" (`modul-09-implementierung.md:70–74`) — eine noch
  stärkere Stütze desselben Ergebnisses.
- geprüft, ohne Befund: **Diff-Scope** — `git show --stat` bestätigt
  ausschließlich `AGENTS.md` (1 Datei, 8 Einfügungen, 3 Löschungen); keine
  Code-, Test- oder Config-Änderung im selben Commit.
- geprüft, ohne Befund: **Kommentar-Form** (§3.7) im neuen Text — der
  eingefügte `AGENTS.md`-Fließtext trägt keine Slice-Nummer, kein
  Mess-Label und keine Review-Historie im Dateikörper selbst; Herkunft bleibt
  in der Commit-Botschaft.
- geprüft, im Kern korrekt (Präzisierung siehe F-4): **codepaths-Aussage**
  („findet solche Token, verbietet sie aber nicht") — `.d-check.yml:226–240`
  bestätigt, dass `codepaths` Pfad-Existenz prüft, aber keine Regel führt, die
  ein role-statt-pfad-Vokabular in `spec/architecture.md` erzwingt; kein
  anderes Modul (`structure`, `spans`, `versions`) hat dafür Scope.
- geprüft, ohne auflösenden Treffer: **`grundlagen-referenz-richtung.md`** auf
  eine Stelle, die den Begriff „Modul-Pfad" definiert oder die Mehrdeutigkeit
  auflöst — keine gefunden (trägt zu F-3 bei).
- nicht nachvollzogen (Auftrag untersagt eigene `make`-Läufe): die Gate-Aussage
  „`make gates` Exit 0 (zehn Glieder, 477 Dateien, 0 Befunde)" in der
  Commit-Botschaft — weder bestätigt noch widerlegt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 4 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** gedeckt-verdikt-deckt-nur-eine-von-fuenf-kategorien
· null-zaehlung-uebersieht-bedienhinweis · vorlagen-praxis-statt-direkterer-textbeleg
· suchmuster-abhaengige-kennzahl · unausgewiesener-restsatz-im-abschnitt

## Verdikt

**Merge-blockierend:** ja — ein HIGH (F-1) und vier MEDIUM. F-1 wiederholt,
eine Ebene tiefer, exakt das Muster, das der vorherige Review am selben
Abschnitt bereits als HIGH gewertet hat (`docs/reviews/2026-08-23-slice-132-hard-rule-zensus-review.md`
F-1): eine „gedeckt"-Behauptung, die nur eine von mehreren in der Regel
genannten Kategorien tatsächlich mit einem Gate belegt. Die drei
Doppelungs-/Verschärfungs-Kernaussagen des Commits (Punkte 1 und 2 des
Auftrags) halten der Prüfung stand — F-2 und F-3 zeigen aber, dass die
Beleglage für „keine Verschärfung" dünner ist, als die Commit-Botschaft
suggeriert (übersehener Pfad im Bedienhinweis, ungenutzte direktere
Text-Evidenz, die dagegen spricht). F-4 zeigt, dass die konkrete Zahl „fünf"
nicht aus dem im Text verwendeten Begriffsverständnis von „Modul-Pfad" folgt.
F-5 ist eine Lücke in der behaupteten Vollständigkeit der Neu-Ausweisung.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den Beobachtungs-Zähler.
Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation (DoD-/
Spec-Konformität prüft der Verifier separat).
