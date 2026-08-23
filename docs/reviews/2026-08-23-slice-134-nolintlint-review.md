# Review-Report: slice-134 — nolintlint im Profil — 2026-08-23

**Review-Art:** Code-Review (Konfigurations-/Doku-Diff gegen Kanon und
Slice-Plan, Modul 10 §Drei Review-Arten) · **Gegenstand:** Commit `8abe10b`
(Diff `HEAD~1..HEAD`) — der Feature-Commit von slice-134
(`feat(harness): nolintlint im Profil — die Lücke wird verengt, nicht
geschlossen (slice-134)`). Der Commit ändert `.golangci.yml`, `AGENTS.md`,
`docs/plan/planning/in-progress/slice-134-nolintlint.md` und
`harness/README.md`.

**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 · **Modell-ID:**
`claude-sonnet-5` · **Datum:** 2026-08-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-134-nolintlint.md` (§1–§9,
  vor und nach dem Commit)
- Wellendokument `docs/plan/planning/welle-84-durchsetzung.md` (§1, §3, §4, §6)
- Zensus `docs/plan/planning/done/slice-132-hard-rule-zensus.md` §9 (Tabellen-
  zeile §3.2), §5 (Risikosektion)
- `.golangci.yml` vollständig (vor und nach dem Commit)
- `AGENTS.md` §3.2, §3.6, §4 (vollständig gelesen: §1–§6)
- `harness/README.md` §Sensors (vollständige Tabelle, alle Zeilen)
- `docs/plan/adr/0006-lint-profil-solid.md` vollständig
- `Makefile` (Target `lint`), `Dockerfile` (Stage `lint`)
- `docs/plan/planning/done/slice-005-lint-profil-solid.md` (Ursprungs-Slice
  von ADR-0006, zum Abgleich der Linter-Zahl)
- **Vorheriger Review am selben Modul:**
  `docs/reviews/2026-08-23-slice-132-hard-rule-zensus-review.md` — dessen F-4
  (slice-134 versprach „gedeckt" statt „teilgedeckt") ist bereits vor diesem
  Commit korrigiert (Commit `d835701`); dessen Negativbefund zur
  §3.2-Kernbehauptung galt dem Config-Stand **ohne** `nolintlint` und deckt
  die hier geprüfte `require-explanation`-Frage nicht — neu für diesen Lauf.
- Nicht erhalten: die DoD-Abhakung (Verifikations-Rolle, getrennter Kontext)

**Vom Reviewer selbst gefahren** (nur Lesekommandos, kein `make`, keine
Dateiänderung): `git show`/`git log` auf `8abe10b`; `git grep` nach
`//nolint`/`//lint:ignore` über `*.go` (getrackt, inkl. Tests); `grep` über
alle `.md`/`.yml`-Dateien nach der Linter-Zahl (`23`/`24 Linter`, `SOLID-nah`)
zur Spiegel-Prüfung; Zeilenzählung der `linters.enable`-Liste in
`.golangci.yml`; Lektüre des `exclusions`-Blocks auf `nolintlint`-Ausnahmen.
**Kein** `make`-Target gelaufen (Auftrags-Verbot) — Aussagen zu Linter-
Verhalten sind Konfigurations-Plausibilität, keine Ausführungsbeobachtung.

**Verdikt: nicht blockierend** — zwei MEDIUM, ein INFO. Die Kernbehauptungen
des Commits (Schema-Platzierung, Nullbestand, Zahl 23→24, ADR-0006-Reichweite
für das *Ausgangsprofil*, Auflösungs-Trigger-Form) halten der Prüfung stand;
die Findings betreffen Genauigkeits-/Zurechnungslücken in den begleitenden
Erklärtexten, nicht die Konfiguration selbst.

---

## Findings

### F-1 — Das als „wohlgeformt" bezeichnete Zensus-Beispiel widerspricht der im selben Commit scharfgeschalteten `require-explanation`

- **kategorie:** MEDIUM
- **quelle:** `.golangci.yml:103–106` (`nolintlint`-Settings, `require-
  explanation: true`) · nolintlint-Semantik („require an explanation of
  nonzero length after each nolint directive")
- **pfad:** `docs/plan/planning/in-progress/slice-134-nolintlint.md:38–44`
  (§2 Punkt 3, in diesem Commit umgeschrieben)
- **befund:** Der Plan-Text nennt die Konstruktion
  `` `//nolint:unused,gochecknoglobals` `` (ohne jeden Begründungs-
  Kommentar) „wohlgeformt" und behauptet, sie „wird von `nolintlint`
  **nicht** gemeldet". Bei `require-explanation: true` fehlt genau diesem
  String die Begründung — nolintlint würde ihn als „missing explanation"
  meckern, nicht durchlassen. Wohlgeformt im Sinne aller drei scharfen
  Schalter ist erst die Form, die die Commit-Botschaft selbst als Probe
  zeigt: `//nolint:unused,gochecknoglobals // Probe A` (mit angehängtem
  Begründungs-Kommentar). Der Plan-Text, den dieser Commit gerade zu diesem
  Zweck umformuliert hat, bleibt hinter der eigenen Konfiguration zurück —
  wer die DoD mit exakt der im Plan gezeigten Zeichenkette nachvollzieht,
  bekommt keinen grünen, sondern einen roten Befund.
- **verifizierbar:** ja — `.golangci.yml:103–106` gegen
  `slice-134-nolintlint.md:39` gelesen; kein `make lint`-Lauf nötig, die
  nolintlint-Settings-Semantik ist eindeutig. Ein `make lint`-Lauf mit exakt
  der Plan-Zeichenkette (ohne Kommentar) würde den Unterschied endgültig
  bestätigen.
- **klasse:** beispiel-widerspricht-eigener-konfiguration

### F-2 — `nolintlint` wird im selben Atemzug ADR-0006 zugeschrieben, den der Commit selbst als „über das adoptierte Profil hinaus" ausweist

- **kategorie:** MEDIUM
- **quelle:** `docs/plan/adr/0006-lint-profil-solid.md` (Entscheidung/Re-
  Evaluierungs-Trigger — deckt die SOLID-nahe/u-boot-Paritäts-Linter, nicht
  Suppression-Governance) · `AGENTS.md` §4 (Zeile 276: `make lint` bereits an
  „(§3.2)" gebunden) · `Makefile:65` (Kommentar „AGENTS.md §3.2") ·
  Beobachtungsklasse BEO-012 („Quelle über ihren Geltungsbereich hinaus
  zitiert")
- **pfad:** `.golangci.yml:1–3` und `:9–14` (Kopf-Kommentar, beide Absätze
  in diesem Commit erweitert); `harness/README.md:71` (Spalte „Bindung" der
  Sensors-Zeile `make lint`, in diesem Commit nur im Vertrag-Text, nicht in
  der Bindung-Spalte geändert)
- **befund:** Der Kopf-Kommentar von `.golangci.yml` zählt `nolintlint`
  in „5 Default- + 24 SOLID-nahe Linter nach dem Vorbild u-boot —
  **beschlossen in ADR-0006**" mit — sechs Zeilen weiter, im selben
  Kommentarblock, steht aber ausdrücklich: „`nolintlint` … über das
  adoptierte Profil **hinaus** aufgenommen". Beide Sätze widersprechen sich:
  entweder ist `nolintlint` Teil dessen, was ADR-0006 beschlossen hat, oder
  es liegt außerhalb — der Kommentar behauptet nacheinander beides. In der
  YAML-Struktur selbst steht `nolintlint` zudem kommentarlos unter dem
  bestehenden Abschnitt „# SOLID-nahe Linter", ohne eigene Kennzeichnung als
  andersartig motiviert. Dieselbe Zurechnung fehlt in
  `harness/README.md:71`: die Spalte „Bindung" der `make lint`-Zeile nennt
  weiterhin ausschließlich `ADR-0006`, obwohl `AGENTS.md` §4 denselben
  Sensor bereits an „(§3.2)" bindet, der Makefile-Kommentar dasselbe tut, und
  andere Sensors-Zeilen mit vergleichbarer Doppel-Bindung (`coverage-gate`,
  `adr-check`, `planning-check`) konsequent **beide** Quellen (ADR **und**
  `AGENTS.md`-Paragraph) nennen. Der Commit vermeidet im Fließtext sorgfältig
  genau diesen Fehler für ADR-0006s eigene „24 weitere Linter"-Formulierung
  („Beinahe hätte ich daraus eine Drift gemeldet, die es nie gab") — und
  begeht ihn dann selbst an der neu geschriebenen Stelle.
- **verifizierbar:** ja — Textvergleich der beiden Kommentar-Absätze in
  `.golangci.yml`; Spaltenvergleich der Sensors-Tabelle in
  `harness/README.md` gegen `AGENTS.md` §4 und `Makefile:65`.
- **klasse:** quelle-ueber-geltungsbereich-hinaus-zitiert (BEO-012)

## Negativbefunde

- geprüft, ohne Befund: **Schema-Platzierung.** `nolintlint:` steht unter
  `linters.settings` auf derselben Einrücktiefe wie `maintidx`/`nestif`
  (`.golangci.yml:94–108`) — korrekt fürs v2-Schema, keine tote Einstellung.
- geprüft, ohne Befund: **Bestand.** `git grep -n "//nolint"` und
  `git grep -n "//lint:ignore"` über alle getrackten `*.go`-Dateien
  (inklusive `_test.go`) liefern **null** Treffer — die im Commit behauptete
  Nullmessung ist bestätigt.
- geprüft, ohne Befund: **Linter-Zahl.** Die `linters.enable`-Liste trägt
  nach dem Commit 5 Default- + 24 „SOLID-nahe" Einträge (Handzählung,
  `.golangci.yml:20–52`) — die im Commit behauptete Hebung 23→24 stimmt.
- geprüft, ohne Befund: **ADR-0006-Reichweite für das Ausgangsprofil.**
  ADR-0006 §Kontext beschreibt „24 weitere Linter" ausdrücklich als
  u-boot-Eigenschaft („Das Ökosystem-Vorbild u-boot fährt … ein SOLID-nahes
  Profil"), nicht als d-check-Zahl; `CHANGELOG.md:1790` dokumentiert
  d-checks eigenen Stand nach slice-005 korrekt als „23 Linter" (u-boot
  minus `depguard`). Die commit-eigene Herleitung „minus depguard ergab das
  immer 23" ist für den *Vor-Commit-Stand* zutreffend — unabhängig von
  F-2, die den *Nach-Commit*-Zurechnungsfehler betrifft.
- geprüft, ohne Befund: **Exclusions-Block.** Keine der fünf Regeln unter
  `linters.exclusions.rules` (`.golangci.yml:164–200`) referenziert
  `nolintlint` — insbesondere die `_test.go`-Ausnahmen (cyclop/gocognit/
  gocyclo/nestif/funlen, noctx/unparam, revive unused-parameter/-receiver)
  lassen `nolintlint` unangetastet; er gilt unverändert repo-weit
  einschließlich Testdateien.
- geprüft, ohne Befund: **Andere Suppressionswege.** `//lint:ignore` (die
  staticcheck-native Direktive) liegt außerhalb dessen, was `nolintlint`
  prüft (es kennt nur `//nolint`-Direktiven) — der Commit behauptet an
  keiner Stelle, `nolintlint` decke sie ab; die Bestandsmessung schließt sie
  separat mit ein (0 Treffer), und slice-134 §3 weist einen zweiten Wächter
  dafür ausdrücklich als Out-of-Scope aus. `.golangci.yml`-Pfad-Ausnahmen
  sind der von §3.2 selbst sanktionierte, zentrale Ausnahmeweg — keine
  Umgehung, die der Commit fälschlich als geschlossen ausgäbe.
- geprüft, ohne Befund: **§3.6-Anwendbarkeit.** §3.6 verlangt eine ADR für
  jede „Schwellen-**Senkung**"; die Aufnahme von `nolintlint` ist eine
  Verschärfung (23→24 aktive Linter), keine Senkung, und ADR-0006s eigene
  Re-Evaluierungs-Trigger („u-boot hebt sein Profil substanziell" /
  „systematische False Positives") greifen nicht. Die commit-eigene Annahme
  „Verschärfung ⇒ keine ADR nötig" hält am Wortlaut von §3.6 — unabhängig
  von der in F-2 gemeldeten Zurechnungsungenauigkeit, ob die Verschärfung
  *ADR-0006* oder *§3.2* zuzuordnen ist.
- geprüft, ohne Befund: **§3.2-Text nach dem Commit.** Die neue Fassung
  („meldet eine Direktive ohne benannten Linter, ohne Begründung oder ohne
  Wirkung … prüft die Form … nicht ihre Berechtigung") deckt sich exakt mit
  den drei konfigurierten Schaltern (`require-specific`, `require-
  explanation`, `allow-unused: false`) — weder Über- noch Unterclaim
  gegenüber der tatsächlichen `nolintlint`-Konfiguration.
- geprüft, ohne Befund: **Auflösungs-Trigger-Form.** „permanent" ist die in
  `AGENTS.md` §3.3/§3.4/§3.6/§3.8 etablierte Form für Regeln, deren
  verbleibende Lücke ein Urteil statt ein Gate ist; die Begründung („die
  Berechtigungsfrage ist ein Urteil, ein zweiter Wächter wäre der von §6
  ausgeschlossene Heuristik-Wächter") ist wörtlich durch
  `welle-84-durchsetzung.md` §6 gedeckt.
- geprüft, ohne Befund: **CHANGELOG.md/ADR-0006 unberührt.** Beide
  Vorgänger-Slices derselben Welle (slice-132, slice-133) ließen
  `CHANGELOG.md` ebenfalls unberührt (Durchsetzungs-Welle ohne
  Lastenheft-Delta/Release, `welle-84-durchsetzung.md` §6) — die
  Nicht-Änderung ist konsistent mit dem Wellen-Zuschnitt, keine Auslassung.

### INFO — Ein dritter, vorbestehender Spiegel der Linter-Zahl bleibt inkonsistent (nicht durch diesen Commit verursacht)

- **kategorie:** INFO
- **quelle:** `docs/plan/planning/done/slice-005-lint-profil-solid.md` §1
  (Ziel-Absatz) vs. `CHANGELOG.md:1789–1790` (beide aus demselben
  Ursprungs-Slice)
- **pfad:** `docs/plan/planning/done/slice-005-lint-profil-solid.md:18–19`
- **befund:** slice-005s eigener Ziel-Satz nennt „5 Default- + 24 SOLID-nahe
  Linter" als erreichten Zustand, während derselbe Slice im CHANGELOG-Eintrag
  „23 Linter" dokumentiert (korrekt, da ohne `depguard`) — ein
  Selbstwiderspruch, der seit slice-005 besteht und mit diesem Commit weder
  entstanden noch berührt wurde. Der Commit hebt zwei Spiegel („`.golangci.
  yml`-Kopf und Sensors-Tabelle") und begründet explizit, warum zwei weitere
  (CHANGELOG, ADR-0006) unberührt bleiben — dieser dritte, vorbestehende
  Spiegel taucht in keiner der beiden Listen auf. Da er `done/` und
  historisch ist, kein Merge-Hindernis; nur zur Vollständigkeit der
  „wo steht die Zahl noch"-Fläche notiert.
- **verifizierbar:** ja — Textvergleich der beiden Zeilen.
- **klasse:** vorbestehender-spiegel-nicht-in-audit-liste

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** `beispiel-widerspricht-eigener-
konfiguration` · `quelle-ueber-geltungsbereich-hinaus-zitiert` (BEO-012) ·
`vorbestehender-spiegel-nicht-in-audit-liste`

## Verdikt

**Nicht merge-blockierend.** Die fünf zu prüfenden Kernbehauptungen des
Commits (Schema-Platzierung, Nullbestand, Zahl 23→24, ADR-0006-Reichweite
für den *Ausgangsstand*, §3.2 teilgedeckt/permanent) sind einzeln nachgemessen
und tragen. Beide MEDIUM-Findings sind Genauigkeitslücken in begleitenden
Erklärtexten desselben Commits, nicht in der Konfiguration: F-1 zeigt, dass
das im Slice-Plan gezeigte „wohlgeformt"-Beispiel bei genauer Lektüre der
selbst gesetzten `require-explanation`-Schwelle nicht besteht — wer die DoD
mit der wörtlich gezeigten Zeichenkette nachvollzieht, misst das Gegenteil
der Behauptung. F-2 zeigt, dass der Commit an zwei Stellen (`.golangci.yml`-
Kopf, `harness/README.md`-Bindung-Spalte) exakt die Zurechnungsfalle
begeht, die er im Fließtext für ADR-0006s ursprüngliche „24"-Formulierung
selbst noch vermeidet — `nolintlint` ist eine §3.2-Maßnahme, keine
ADR-0006-Entscheidung, aber beide Spiegel zählen sie ADR-0006 zu. Beide
Findings sind vor dem nächsten Anfassen dieser Dateien zu klären, blockieren
aber nicht die Aufnahme von `nolintlint` selbst, die konfigurativ korrekt und
wirksam ist.

**Übergabe:** Findings gehen an den Implementer (Rückkante Review → Plan/
Config bei Textdefekt); dieser Report ist ein Lauf-Beleg — DoD-/
Spec-Konformität prüft der Verifier separat.
