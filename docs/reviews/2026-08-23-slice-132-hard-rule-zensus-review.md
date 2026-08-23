# Review-Report: slice-132 — Hard-Rule-Zensus — 2026-08-23

**Review-Art:** Code-Review (Doku-/Konventions-Diff gegen Kanon und Slice-Plan,
Modul 10 §Drei Review-Arten) · **Gegenstand:** Commit `15a7830` (`HEAD`,
Diff `HEAD~1..HEAD`) — der Zensus-Commit von slice-132 (`feat(harness): Zensus
über zehn Hard Rules — fünf Ausweisungen, zwei Slices geschnitten (slice-132,
slice-134, slice-135)`). Der Commit ändert `AGENTS.md`, legt
`docs/plan/planning/open/slice-134-nolintlint.md` und
`docs/plan/planning/open/slice-135-uses-pin-sensor.md` neu an und aktualisiert
`docs/plan/planning/welle-84-durchsetzung.md` §4.

**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 (`9a7654a`) · **Modell-ID:**
`claude-sonnet-5` · **Datum:** 2026-08-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-132-hard-rule-zensus.md` (§1–§9)
- Wellendokument `docs/plan/planning/welle-84-durchsetzung.md` (§1, §3, §4)
- `AGENTS.md` §3 (vollständig, §3.1–§3.9), §5, §6
- `docs/plan/planning/open/slice-134-nolintlint.md`,
  `docs/plan/planning/open/slice-135-uses-pin-sensor.md` (beide neu in diesem
  Commit)
- `.golangci.yml`, `.d-check.yml`, `Makefile`,
  `.claude/hooks/pretooluse-command-guard.sh`
- Baseline-Kanon `.harness/baseline/v5.11.0/regelwerk/modul-09-implementierung.md`
  §AGENTS.md-Regeln, `modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin) und
  §Fitness Function aus einem ADR-Satz
- `docs/plan/planning/observations.md` (`BEO-004`, `BEO-007`, `BEO-009`,
  `BEO-010`, `BEO-011`, `BEO-012`)
- Modul-Implementierung zur Verifikation der genannten Befund-Codes:
  `internal/hexagon/core/rules/matrix.go`, `vcs.go`, `planning.go`, `commits.go`
- Nicht erhalten: die DoD-Abhakung (Verifikations-Rolle, getrennter Kontext)

**Vom Reviewer selbst gefahren** (nur Lesekommandos, kein `make`, keine
Dateiänderung): `git show`/`git diff` auf `HEAD~1..HEAD`, `grep` über
`.d-check.yml`, `.github/workflows/*.yml`, `internal/hexagon/core/rules/*.go`
nach den vier im Commit genannten Befund-Codes (`matrix-forbidden`,
`core-drift-vcs`, `planning-drift`, `commit-untraceable`) und nach verdeckten
Gates für die fünf als *einseitig* ausgewiesenen Regeln.

**Verdikt: blockierend** — drei HIGH, zwei MEDIUM, kein LOW/INFO-Substanzfund
über die unten benannten hinaus.

---

## Findings

### F-1 — Das GEDECKT-Verdikt für §3.4 deckt nur eine von zwei Teilaussagen der Regel

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §3.4 (Zeile 153–161) — die Regel trägt zwei
  Teilaussagen: (a) `spec/architecture.md` ist sprach-/meilensteinfrei, (b)
  kein Spec-Stratum referenziert ADRs/Wellen/Slices/Commit-Hashes/Closure-Daten
  abwärts · `.d-check.yml` Modul `matrix` (Zeile 179–222) · slice-132 §5
  („„Gate genannt" ist nicht „Regel getragen". Ein Gate kann einen Teil der
  Regel prüfen und der Zeile trotzdem ein *gedeckt* verschaffen.")
- **pfad:** Commit-Botschaft, Absatz „GEDECKT (2)" (Zeile 6–9 der
  Commit-Message); `.d-check.yml:184–222`
- **befund:** Der Zensus belegt „§3.4 gedeckt" ausschließlich mit der Probe
  „`slice-999` in `spec/spezifikation.md` ⇒ `matrix-forbidden`" — das prüft nur
  Teilaussage (b). Für Teilaussage (a) (Sprach-/Modul-Pfad-Freiheit von
  `spec/architecture.md`) existiert **keine** Prüfung: kein Modul in
  `.d-check.yml` liest den Inhalt von `spec/architecture.md` auf
  Sprach-/Technologie-Begriffe (`matrix` prüft nur Referenzrichtung zwischen
  Klassen, nicht Wortschatz innerhalb einer Datei); kein anderes Modul
  (`structure`, `spans`, `versions`) hat einen Scope, der das leisten würde.
  Anders als bei §3.5 trägt `AGENTS.md` §3.4 zudem **keinen** neuen Text, der
  die Deckung (auch nur teilweise) dokumentiert — die einzige Spur der
  „gedeckt"-Behauptung ist die Commit-Botschaft selbst.
- **verifizierbar:** ja — `.d-check.yml`-Modulliste und `matrix`-Konfiguration
  enthalten kein Muster für Sprach-/Technologie-Begriffe in
  `spec/architecture.md`; ein konstruierter Verstoß gegen Teilaussage (a)
  (z. B. ein Go-Import-Pfad in `spec/architecture.md`) würde kein Gate rot
  färben. Nicht selbst ausgeführt (Verbot, `make`-Targets zu laufen).
- **klasse:** gedeckt-verdikt-deckt-nur-halbe-regel

### F-2 — Die Zeilen-Arithmetik des Zensus geht nicht auf: „zehn" vs. tatsächlich elf Verdikt-Einträge, „EINSEITIG (5)" vs. sechs gelistete Regeln

- **kategorie:** HIGH
- **quelle:** slice-132 DoD-Punkt 1 („Je Hard Rule … genau eine Zeile; keine
  Regel ohne Zeile, keine Zeile ohne Regel") · `AGENTS.md` §5 Botschafts-Regel
  (Zeile 349–356, „eine genannte Probe … ihr Schluss reicht nicht weiter als
  die gemessene Menge") · `BEO-011` (Vollständigkeits-Aussagen, achtmal in
  welle-82 gekippt)
- **pfad:** Commit-Botschaft (Titel + Absatz „EINSEITIG (5)");
  `docs/plan/planning/welle-84-durchsetzung.md:26–30`
- **befund:** Der Commit löst §3.7 bewusst in **zwei** Regeln auf („Auflösung
  der §3.7-Zeile in ZWEI Regeln folgt der slice-127-Lehre: ein Abschnitt ist
  keine Regel"). Rechnet man diese Auflösung durch, ergeben sich GEDECKT (2) +
  TEILGEDECKT (3) + EINSEITIG (§3.2, §3.6, §3.7a, §3.7b, §3.8, §3.9 = 6) = **11**
  Verdikt-Zeilen — nicht die im Titel behaupteten „zehn Hard Rules". Der
  Absatz-Kopf „EINSEITIG (**5**)" listet im selben Satz **sechs** §-Kennungen;
  „5" passt nur, wenn man die fünf **Absätze** zählt, die in `AGENTS.md`
  hinterlegt sind (§3.2, §3.6, §3.7-kombiniert, §3.8, §3.9 — §3.7a/§3.7b teilen
  sich einen Absatz), während die Liste daneben die sechs **Regeln** aufzählt.
  Beide Einheiten (Absätze vs. Regeln) werden vermischt, ohne dass der Commit
  das benennt. `welle-84-durchsetzung.md` §1 (unverändert von diesem Commit)
  trägt weiterhin „§3 trägt neun Hard Rules, dazu die Botschafts-Regel in §5"
  — die Zehn-Zählung, die genau die Sektions-statt-Regel-Logik fortschreibt,
  die dieser Commit an §3.7 explizit korrigiert.
- **verifizierbar:** ja — Volltext-Zählung der Commit-Botschaft (sechs
  §-Kennungen unter der Überschrift „(5)") und Abgleich mit
  `welle-84-durchsetzung.md:26–30`.
- **klasse:** vollstaendigkeits-arithmetik-inkonsistent

### F-3 — §5 wird nur mit ihrer einen „Botschafts-Regel" zitiert; die übrigen Hard Rules des Abschnitts bleiben ungenannt und unverdiktet

- **kategorie:** HIGH
- **quelle:** `AGENTS.md` §5 „Dokumentations-Regeln" (Zeile 299–357, mindestens
  zwölf eigenständige Imperative: Kennungs-Vergabe nur beim Spec-/ADR-Schreiben,
  `DC-*`-Anlage nur im Lastenheft, ADR-Index-Pflicht, Re-Evaluierungs-Trigger-
  Pflicht neuer ADRs, `**Verantwortlich:**`-Feld-Pflicht, `**Berührte Spec-
  Stellen:**`-Form, Status-Feld-Verbot für neue Slices, Vorprüfungs-Pflicht vor
  der Sub-Area-Modus-Begründung, CHANGELOG-Pflege) · `BEO-012` („eine Quelle
  wird über ihren Geltungsbereich hinaus zitiert … ein Eintrag wird nach seinem
  Titel zitiert statt nach seinem Feld")
- **pfad:** slice-132.md Zeile 23 („dazu die Botschafts-Regel in §5"),
  `welle-84-durchsetzung.md:26,29` (dieselbe Formel, von diesem Commit nicht
  korrigiert)
- **befund:** Der Zensus-Auftrag lautet wörtlich „je Hard Rule in §3 … und der
  Botschafts-Regel in §5" — §5 wird damit von vornherein auf **eine einzige**
  ihrer Regeln verengt, ohne dass der Commit begründet, warum die übrigen
  Regeln des Abschnitts (der als Ganzes „Dokumentations-Regeln" heißt, nicht
  „Botschafts-Regel") aus dem Zensus fallen. Für keine der übrigen §5-Regeln
  gibt es eine Zeile, einen Verweis auf eine frühere Zensus-Antwort oder eine
  explizite Out-of-Scope-Notiz — sie werden weder als *gedeckt* noch
  *teilgedeckt* noch *einseitig* geführt, sie fehlen ganz. Das ist exakt das
  Muster, gegen das `BEO-012` warnt: eine Quelle (§5) wird nach ihrem im
  Zensus-Titel genannten Etikett zitiert, nicht nach ihrem tatsächlichen
  Geltungsfeld gelesen.
- **verifizierbar:** ja — Volltext-Lesung von `AGENTS.md` §5 gegen die
  Zensus-Zeilenliste in Commit-Botschaft und `slice-132.md`.
- **klasse:** geltungsbereich-per-titel-statt-inhalt

### F-4 — slice-134 verspricht im Vorgehen einen vollen Deckungswechsel, den die eigene Risikosektion als unhaltbar benennt

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §5 Botschafts-Regel („eine genannte Probe … ihr
  Schluss reicht nicht weiter als die gemessene Menge") · Reviewer-Skill
  MEDIUM-Anker „Botschaft verallgemeinert über die Messung hinaus"
- **pfad:** `docs/plan/planning/open/slice-134-nolintlint.md:41–44` vs.
  `:73–76`
- **befund:** §2 Schritt 4 sagt voraus: „`AGENTS.md` §3.2 nachziehen: der
  Auflösungs-Trigger ist eingelöst, die Regel wechselt von *einseitig* auf
  *gedeckt*". §5 (Abnahme-Punkte/Risiken) desselben Dokuments widerspricht dem
  bereits jetzt: „`nolintlint` prüft die Form der Direktive, nicht ihre
  Berechtigung. Eine begründete, spezifische `//nolint` besteht ihn — die Regel
  verbietet sie trotzdem. Die Deckung wird also **teilweise** sein, und das
  gehört benannt statt aufgerundet." Das Vorgehen bucht das volle *gedeckt*
  bereits als Ergebnis, während die eigene Risikosektion im selben Slice-Plan
  festhält, dass genau das eine Überdehnung wäre.
- **verifizierbar:** ja — Textvergleich der beiden Abschnitte in derselben
  Datei.
- **klasse:** plan-verspricht-mehr-als-eigenes-risiko-zulaesst

### F-5 — §3.1 wird *teilgedeckt* genannt, obwohl die tragende Mechanik im selben Satz als „kein Repo-Gate" ausgewiesen ist

- **kategorie:** MEDIUM
- **quelle:** slice-132 §2 Punkt 2 (Teilgedeckt-Definition: „ein Gate prüft
  einen Teil der Regel") · `AGENTS.md` §4 Kopfsatz („Nur hier gelistete
  Targets existieren im Makefile")
- **pfad:** `AGENTS.md:96–101`
- **befund:** Die neue §3.1-Passage benennt den Tool-Call-Wächter selbst als
  „werkzeug-lokal, **kein Repo-Gate**: kein `make`-Target und keine CI ruft
  ihn". Bei §3.3 und der §5-Botschaftsregel bedeutet *teilgedeckt*, dass ein
  echtes Gate (`planning-check`, `trace-check`) einen **Teil des Regel-Inhalts**
  abdeckt und ein anderer Teil ungedeckt bleibt — bei §3.1 deckt dieselbe
  Wortmarke stattdessen einen Fall, in dem **kein** Gate beteiligt ist, sondern
  ein Mechanismus außerhalb der in §4 definierten Gate-Landschaft. Die
  Leitfrage des Zensus lautet „welcher Gate-**Lauf** trägt sie" — nach dieser
  Definition wäre die konsequente Antwort für §3.1 „keiner, aber ein
  Nicht-Gate-Mechanismus deckt einen Teil" oder schlicht *einseitig* (wie
  §3.8/§3.9, die ebenfalls ohne jedes Gate auskommen). Der Bedeutungswechsel
  von *teilgedeckt* zwischen den drei so verdikteten Zeilen wird im Commit
  nicht benannt.
- **verifizierbar:** teilweise — die Inkonsistenz ist Text-gegen-Text
  verifizierbar; ob §3.1 richtig kategorisiert ist, ist ein Urteil, kein Gate.
- **klasse:** teilgedeckt-verdikt-uneinheitlich-begruendet

## Negativbefunde

- geprüft, ohne Befund: §3.5 GEDECKT — die Probe (Kern-Satz-Edit in ADR-0005,
  gestagt ⇒ `core-drift-vcs`) trifft die Regel vollständig; beide von §3.5
  genannten Ausnahmen (`## Geschichte`-Anhänge, `**Status:**`-Übergang) sind in
  `.d-check.yml` (`vcs.exclude-sections`, `vcs.head-allow`) exakt gespiegelt —
  kein unentdeckter Teil-Gap wie bei §3.4
- geprüft, ohne Befund: der §3.2-Kernbefund selbst — `.golangci.yml` aktiviert
  `unused` und `gochecknoglobals`, führt aber `nolintlint` nicht in der
  `enable`-Liste; die im Commit behauptete dreistufige Messkette (rot ohne
  Direktive → rot mit nur `//nolint:unused` [gochecknoglobals feuert weiter] →
  Exit 0 mit beiden) ist mit dieser Konfiguration plausibel und in sich
  konsistent
- geprüft, ohne Befund: die fünf neu/vorbestehend in `AGENTS.md` §3
  ausgewiesenen „Kein Gate prüft"-Passagen (§3.1, §3.2, §3.3, §3.6, §3.7)
  tragen weder Slice-Verweise noch Review-Historie/Forensik — konform zu §3.7
  Kommentar-/Zustandsfeld-Klassen und der Slice-Verweis-Verbot in `AGENTS.md`
  selbst
- geprüft, ohne Befund: die vier im Commit genannten Befund-Codes
  (`matrix-forbidden`, `core-drift-vcs`, `planning-drift`,
  `commit-untraceable`) existieren alle mit der behaupteten Bedeutung im Code
  (`rules/matrix.go:14`, `rules/vcs.go` via `model.ReasonCoreDriftVCS`,
  `rules/planning.go:363`, `rules/commits.go:61`) — keine Harness-Lüge
  (halluzinierter Gate-/Befund-Name)
- geprüft, ohne Befund: für §3.6, §3.8, §3.9 (einseitig) findet eine
  repo-weite Suche (Go-Tests, `.d-check.yml`-Module, Workflows, Hooks) keinen
  übersehenen tragenden Lauf — insbesondere kein `uses:`-Pin-Sensor in
  `.github/workflows/*.yml` (bestätigt §3.9) und kein Test, der die
  Scan-Menge-Zusage eines Moduls gegen seine gelesenen-aber-nicht-gescannten
  Eingaben hält (bestätigt §3.8)
- geprüft, ohne Befund: `slice-135-uses-pin-sensor.md` — Zusage sauber auf
  „Form, nicht Gültigkeit" geschnitten (§2 Punkt 2), Risiken-Sektion benennt
  die Überdehnungsgefahr („„Die Workflows sind gehärtet" wäre die Aussage, die
  der Sensor nicht trägt") korrekt vorab, kein Widerspruch zwischen Vorgehen
  und Risiken wie bei slice-134
- geprüft, ohne Befund: §3.3 Ausnahme „MR-/Wellen-Lifecycle-Move" verlangt bei
  sinkendem Rename-Score eine explizite `git mv`-Deklaration in der
  Commit-Botschaft; diese Teilfacette ist im `commits`-Modul (`.d-check.yml`
  `id-patterns`) nicht geprüft — deckt sich mit, verschärft aber nicht die im
  Commit bereits benannte Lücke „Commit-Zerlegung sieht kein Gate" (kein
  eigenständiges Finding, da bereits von der bestehenden Ausweisung
  eingeschlossen)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 3 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** `gedeckt-verdikt-deckt-nur-halbe-regel` ·
`vollstaendigkeits-arithmetik-inkonsistent` ·
`geltungsbereich-per-titel-statt-inhalt` ·
`plan-verspricht-mehr-als-eigenes-risiko-zulaesst` ·
`teilgedeckt-verdikt-uneinheitlich-begruendet`

## Verdikt

**Merge-blockierend:** ja — drei HIGH offen. Kein Finding widerlegt, dass die
fünf im Regeltext ausgewiesenen einseitigen Zeilen (§3.2, §3.6, §3.7a, §3.7b,
§3.8, §3.9 — sechs Regeln, s. F-2) tatsächlich ungedeckt sind, und der
§3.2-Kernbefund (wirkender Umgehungspfad) hält der Prüfung stand (Negativbefund
oben). Die drei HIGH-Findings treffen die Kern-Prämisse des Slice direkt: F-1
zeigt, dass eines der beiden GEDECKT-Verdikte selbst in die Falle läuft, die
slice-132 §5 als Risiko benennt („Gate genannt ≠ Regel getragen"); F-2 zeigt,
dass die Bijektion „keine Regel ohne Zeile, keine Zeile ohne Regel" (DoD-Punkt
1) nach der eigenen §3.7-Auflösung nicht mehr aufgeht; F-3 zeigt, dass der
Zensus eine ganze Sektion (§5) auf ihr Titel-Etikett verengt, statt ihren
Inhalt zu zählen — dieselbe Klasse (`BEO-012`), die slice-132 §7 selbst als zu
prüfendes Risiko nennt. F-4/F-5 sind Konsistenz-Lücken in den geschnittenen
Folge-Slices bzw. in der Verdikt-Definition, die vor Bau/Closure zu klären
sind.

**Übergabe:** Findings gehen an den Implementer (Rückkante Review → Plan bei
Plan-Defekt); die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und
von dort in den Zähler. Dieser Report selbst ist ein Lauf-Beleg — DoD-/
Spec-Konformität prüft der Verifier separat.
