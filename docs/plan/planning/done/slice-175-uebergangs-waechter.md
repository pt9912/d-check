# Slice slice-175: Die Closure-Vorbedingungen hängen jetzt am Übergang

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-86](welle-86-closure-uebergang-durchsetzen.md) — der
Closure-Trigger der Welle verlangt vier Slices in `done/`; dieser ist der
vierte und letzte, und er **setzt die drei anderen voraus** (welle-86 §5): er
bindet die Prüfungen, die 172–174 liefern, an den Übergang selbst.

**Bezug:** welle-86 §1 (die vier Vorbedingungen, davon vier ursprünglich
offen), §4 „Der Träger von slice-175 ist der git-Hook, nicht der
Werkzeug-Hook" (die Design-Entscheidung, hier umgesetzt), [`MR-042`](../../../../harness/conventions.md#mr-042)
(Werkzeug- vs. Repo-Reichweite der beiden Hook-Familien).

**Berührte Spec-Stellen:** — (Verdrahtung bestehender Fähigkeiten an einen
neuen Bindepunkt; keine neue `DC-FA`-Anforderung, kein neuer Grund-Code).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-03.

---

## 1. Ziel

Ein gestagter Rename/Add einer Slice-Datei nach `docs/plan/planning/done/`
löst automatisch `make verify-closure-notes` aus — im lokalen `pre-commit`-Hook
**und** in der PR-/Push-CI über die Commit-Range. Die Vorbedingungen für
`done/` (DoD-Häkchen, Closure-Notiz, Risiko-Ausgänge) hängen damit am
**Übergang**, nicht mehr nur an einer gelegentlichen `make fullbuild`-Prüfung.

## 2. Vorgehen

1. **Kein neues Modul, keine neue Prüf-Logik.** `verify-closure-notes`
   existiert (Module `planning`/`structure`/`spans`,
   [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md)/[ADR-0077](../../adr/0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md));
   die Lücke ist ausschließlich der **Bindepunkt**.
2. **`.githooks/pre-commit` erweitern:** nach dem bestehenden `doc-check` eine
   Erkennung, ob der gestagte Diff einen Rename/Add nach
   `docs/plan/planning/done/slice-<NNN>*.md` trägt (nicht rekursiv — ein
   archivierter Stub eine Ebene tiefer, `done/<welle-id>/…`, zählt nicht,
   siehe §3). Trifft das zu: `make verify-closure-notes` fahren; `set -e`
   trägt den Abbruch.
3. **`.github/workflows/ci.yml` um dieselbe Erkennung ergänzen**, aber über
   die Commit-**Range** statt `--staged` — `--no-verify` umgeht nur den
   lokalen Hook, nicht die CI (dieselbe Asymmetrie wie bei `adr-check`).
4. **Zwei reale Proben statt einer Behauptung:** ein konstruierter
   Test-Kandidat mit offenem DoD-Haken/fehlendem Abschnitt/dünner
   Closure-Notiz wird als echter `git commit`-Versuch abgewiesen; derselbe
   Kandidat, korrigiert, wird angenommen. Beide Läufe **und ihre Ausgabe**
   stehen in §9, nicht nur ihr Ergebnis. Der Test-Commit der Positiv-Probe
   wird danach per `git reset --soft` wieder entfernt — er gehört nicht in
   die echte Historie.
5. **Doku-Nachzug:** `AGENTS.md` §4 (`make hooks`-Zeile) und
   `harness/README.md` (Sensors-Tabelle, dieselbe Zeile).

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Stop-Hook.** Der Welle-Plan benennt ihn ausdrücklich als Option
  („slice-175 entscheidet, ob er beide Hälften nimmt") — ein `Stop`-Hook ist
  **werkzeug-lokal** ([`MR-042`](../../../../harness/conventions.md#mr-042)),
  ein Verstoß bliebe für jedes andere Werkzeug und für CI unsichtbar. Der
  git-Hook plus CI trägt die **Repo-Invariante** bereits vollständig; ein
  Stop-Hook wäre nur eine billigere, zusätzliche Bequemlichkeitsschicht
  **oberhalb** dieser Invariante, keine Voraussetzung für sie. Bleibt ein
  offener, unabhängiger Kandidat.
- **Kein pfadgebundener Rules-Kanal** (`.claude/rules/` mit `paths`) —
  dieselbe Werkzeug-Grenze, dazu bereits als [`BEO-024`](../observations.md)
  geführt; ändert nichts an der hier gebauten Repo-Invariante.
- **Rekursiver Scan von `done/`-Unterverzeichnissen.** Ein archivierter
  Slice-Stub (`done/<welle-id>/slice-*.md`) trägt keine DoD mehr;
  [`MR-059`](../../../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)s
  Ein-Commit-Form ersetzt dort ohnehin den Volltext, keine neue Zusage
  entsteht.
- **Aufnahme des Moduls `reviews` in diese Bindung.**
  [ADR-0081](../../adr/0081-reviews-modul.md) hat die Aufnahme in
  `gates`/erzwungene Läufe bewusst vertagt (eigener
  Re-Evaluierungs-Trigger); dieser Slice zieht das nicht vor.
- **Die Qualität eines Reviews oder die Rollen-Trennung selbst** — welle-86
  §6 schließt beides ausdrücklich aus dieser Welle aus.
- **Das Nachrüsten des Bestands.** Kein `done/`-Slice vor diesem Commit wird
  rückwirkend geprüft; der Wächter gilt für künftige Übergänge.

## 4. Definition of Done

- [x] `.githooks/pre-commit` löst `make verify-closure-notes` bei einem
      gestagten Rename/Add nach `docs/plan/planning/done/slice-*.md` aus,
      nicht rekursiv.
- [x] `.github/workflows/ci.yml` trägt dieselbe Bindung über die
      Commit-Range im bestehenden Traceability-/ADR-Immutable-Schritt.
- [x] **Negativ-Probe real gefahren:** ein konstruierter Test-Kandidat mit
      offenem DoD-Haken, fehlendem `## 5.`-Abschnitt und dünner
      Closure-Notiz wird von `git commit` abgewiesen (Ausgabe in §9).
- [x] **Positiv-Probe real gefahren:** derselbe Kandidat, korrigiert, wird
      angenommen (Ausgabe in §9); der Test-Commit ist per `git reset --soft`
      wieder entfernt, kein Test-Artefakt in der echten Historie.
- [x] Doku-Update: `AGENTS.md` §4 und `harness/README.md` Sensors-Tabelle
      nennen die neue Bindung.
- [x] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register
      fortgeschrieben (oder Fehlanzeige begründet); jedes Risiko aus §5 mit
      Ausgang; die drei Paarungen laufen bei der Welle-Closure.

## 5. Abnahme-Punkte / Risiken

- **Ein Teil-Stage (nur die Datei, nicht z. B. begleitende Doku-Änderungen)
  könnte den Hook gegen einen inkonsistenten Arbeitsbaum laufen lassen** —
  derselbe bekannte Git-Hook-Rand wie bei `adr-check STAGED=1`, nicht neu für
  diesen Slice. — **Ausgang:** *entfallen* — `verify-closure-notes` liest nur
  die Slice-Datei selbst (Closure-Notiz, DoD, Risiken); sie hat keine
  Teil-Stage-Abhängigkeit zu anderen Dateien.
- **Der CI-Range-Zweig baut ein zusätzliches Image**, weil
  `verify-closure-notes` von `build` abhängt — Laufzeit-Kosten für **jeden**
  Push, nicht nur für Closure-Commits. — **Ausgang:** *entfallen* — der Zweig
  fährt nur, wenn die Erkennung selbst (ein `git diff`, keine Docker-Aktion)
  einen Treffer meldet; ohne Treffer bleibt der Zusatzaufwand aus, und
  `adr-check` im selben Schritt baut das Image ohnehin schon.
- **Ein zukünftiger Konsument dieses Repos (Fork, Schwester-Repo) übernimmt
  `.githooks/pre-commit` unverändert, aber nicht `docs/plan/planning/`s
  Konvention** — der Wächter liefe dann gegen eine leere oder andersartige
  Menge. — **Ausgang:** *weiter offen* → kein neuer `BEO`-Eintrag (kein
  gemessener Fall, reine Möglichkeit); genannt als Grenze, nicht als Risiko
  mit Handlungsbedarf.
- **Kein Stop-Hook** (siehe §3) lässt einen Verstoß bis zum Commit-Versuch
  laufen, nicht schon beim Schreiben. — **Ausgang:** *entfallen* — bewusste
  Entscheidung dieses Slice (§3), keine Lücke: der git-Hook plus CI trägt die
  Repo-Invariante bereits vollständig.

## 6. Trigger

**Start** (`open` → `in-progress`): direkt beansprucht — WIP-Limit frei.

**Rückführungen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls sich zeigt,
  dass Hook **und** CI **und** ein Stop-Hook gemeinsam nötig sind, um die
  Welle zu schließen — trat nicht ein (§3 begründet, warum der Stop-Hook
  keine Voraussetzung ist).
- `in-progress` → `open` (blockiert): keine Bedingung erkannt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `.githooks/` (der Repo-Hook) und
  `.github/workflows/` (die CI) — beide fallen unter den Default `*` =
  **Greenfield** ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration pro Sub-Area). Die Regel, die diesen Schritt
  vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:219-220 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-09-03, höchste
  Kennung `BEO-027`): [`BEO-024`](../observations.md) (Zähler 1) — betrifft
  denselben Hooks-/Kanal-Kontext (ein pfadgebundener Rules-Kanal hängt an der
  Arbeitsweise, nicht am Inhalt), aber nicht die hier gebaute Repo-Invariante
  (§3 grenzt den Rules-Kanal ausdrücklich aus). Kein weiterer Treffer. Die
  Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:225-225 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): `upstream-drift.yml`
  meldet weiterhin **ROT** (Lauf 2026-09-03T05:23:56Z, unverändert seit
  slice-173: zwei planmäßige `VERALTET`-Meldungen, `go`/`semgrep`),
  `image-scan.yml` grün (2026-09-03T08:05:17Z). Keiner der beiden Funde
  berührt diesen Slice. **Dieser Block trägt bewusst keine `cite`-Direktive**
  — sein Ziel ist eine Repo-Adaption, kein Baseline-Abschnitt
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-175. Betroffene IDs: keine neue `DC-FA`-Anforderung. Module:
`planning`, `structure`, `spans` (unverändert, nur neu gebunden). Gates:
`make test`, `make verify-closure-notes`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Es entsteht keine neue Fähigkeit und kein
Fremdsystem; die Konventions-Dichte ist hoch (dasselbe `STAGED=`/`RANGE=`-Muster
trägt bereits `adr-check`, hier eins zu eins wiederverwendet).

## 9. Closure-Notiz (nach `done/`)

- **Was hat funktioniert:** Das Muster lag tatsächlich vollständig vor, wie
  welle-86 §4 vorhersagte — `.githooks/pre-commit` musste nur um eine
  Erkennung plus einen bedingten `make`-Aufruf wachsen, keine neue Logik.
  Beide Proben liefen **real** als `git commit`-Versuche, nicht simuliert:
  die Negativ-Probe zeigte `d-check: 425 Datei(en) geprüft, 3 Befund(e)`
  (`section-missing`, `section-tasks-open`, `closure-note-thin`) und wurde
  mit `make: *** [Makefile:348: verify-closure-notes] Fehler 1` abgewiesen
  — kein neuer Commit entstand (`git log` blieb auf dem Vorgänger-Commit
  stehen). Die Positiv-Probe (derselbe Kandidat mit gesetztem Haken, einem
  `## 5.`-Abschnitt und vier Sätzen in der Closure-Notiz) wurde als Commit
  `07afe62` angenommen, sofort per `git reset --soft HEAD~1` wieder entfernt.
- **Was ging anders als geplant:** Die erste Fassung der Positiv-Probe
  scheiterte noch einmal an `closure-note-thin` (3 statt 4
  Satzende-Zeichen verlangt) — derselbe Sensor, der die Negativ-Probe
  zu Recht abwies, maß auch die eigene Korrektur zunächst zu knapp. Ein
  vierter Satz behob es. Nebenbefund: `git commit`s Exit-Code in einer
  Pipe (`| tail -30`) ist der von `tail`, nicht von `git` — für die Proben
  ohne Bedeutung, weil `git log`/`git status` das eigentliche Ergebnis
  direkt zeigen, aber notiert, damit niemand `$?` danach vertraut.
- **Der unabhängige Review fand ein LOW, vor der Closure behoben:**
  `awk '{print $NF}'` spaltet auf **jedem** Whitespace, nicht nur auf dem
  Tab, den `git diff --name-status` als Feldtrenner benutzt — ein
  Dateiname mit Leerzeichen hätte die Erkennung **still** (fail-open, nicht
  fail-closed) verfehlt. Behoben mit `awk -F'\t'`, in beiden Dateien; beide
  Regressions-Fälle (slice-173-Move erkannt, welle-87-Archiv-Stub nicht
  erkannt) bleiben nach dem Fix unverändert richtig. **Zusätzlich vom
  Review verifiziert, ohne Änderungsbedarf:** ein pfadgebundener `git diff`
  auf `docs/plan/planning/done/` zeigt einen `in-progress→done`-Move **nie**
  als `R` (nur als `A`) — `--diff-filter=AR` fängt ihn trotzdem über die
  `A`-Hälfte; und `set -e` propagiert eine fehlschlagende Zeile **innerhalb**
  eines `if …; then …; fi`-Blocks korrekt, nur die Bedingung selbst ist
  ausgenommen.
- **Steering-Loop-Eintrag:** keiner verkörpert — alle drei Beobachtungen
  dieser Session (zwei eigene, eine vom Review gefunden) sind Einzelfälle,
  keine dritte Instanz einer bereits geführten Beobachtung.
- **Beobachtungs-Register (`../observations.md`):** keine Beobachtung
  angefallen, die eine neue Kennung oder einen Zähler-Schritt rechtfertigt.
- **Folge-Slices:** keiner für den Stop-Hook — er bleibt ein benannter,
  unabhängiger Kandidat (§3/§5), kein verbindlicher Nachfolger.
- **Risiken aus §5:** siehe §5, je Zeile ein Ausgang.
- **Drei Paarungen:** nicht hier geprüft — Repo mit Wellen-Betrieb, prüft die
  Closure von welle-86, die mit diesem Slice ihre vier Vorbedingungen
  vollständig liefert.
