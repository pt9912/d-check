# Review — slice-157 / Commit `25f0290`

**Review-Art:** Code (gegen Slice-Plan §3–§5, `AGENTS.md` §3, MR-042/044/046/047,
`harness/README.md` §Sensors) · **Gegenstand:** `25f0290` —
`internal/adapter/driven/configyaml/enforcement_layer_test.go` (neu),
`AGENTS.md` §4-Zeile, `harness/README.md` §Sensors-Zeile · **Skill:**
`reviewer.md` @ 1.10.0 · **Modell:** claude-opus-5[1m] · **Datum:** 2026-08-27

**Gefahrene Messungen** (drei `make test`-Läufe, Exit aus dem Log):

| Lauf | Zustand | Exit | gelesene Meldung |
|---|---|---|---|
| Baseline | unverändert | 0 | `ok … /configyaml 0.016s` |
| Mutation A | Stop-Hook → `bash .claude/hooks/gibt-es-nicht.sh` | **0** | `ok …` — kein Befund |
| Mutation B | Stop-Hook → `bash "$CLAUDE_PROJECT_DIR"/.claude/hooks/gibt-es-nicht.sh` | 2 | `verdrahteter Hook ".claude/hooks/gibt-es-nicht.sh" existiert nicht` |

A und B unterscheiden sich **nur in der Schreibweise**, nicht im Zustand.

## Findings

### F-1 · HIGH · `hookPaths` ist blind für die einfachste Verdrahtungs-Form

**Pfad:** `enforcement_layer_test.go:94` (Marker), `:98–104`, `:148–163`,
`:88–91` (Grenz-Kommentar).

Der Marker ist `"/.claude/hooks/"` **mit führendem Schrägstrich**. Ein
Hook-Kommando in literaler Relativform trägt ihn nicht, fällt aus `hookPaths`
heraus und wird nie geprüft — Mutation A ist grün, obwohl die verdrahtete Datei
fehlt. Der Riegel `len(paths)==0` greift nicht, weil er erst feuert, wenn
**alle** Einträge den Marker verlieren. Der Test hält genau den Fall nicht, für
den er laut seinem eigenen Kommentar gebaut ist.

Verschärfend: der Grenz-Kommentar nennt als blinde Klasse „ein Hook, der seinen
Pfad **zusammensetzt oder aus einer Variablen holt**" — **umgekehrt** zum
gemessenen Verhalten. Wer ihn liest, schließt, dass literale Pfade sicher sind.

Zwei weitere Ausprägungen (aus dem Code gelesen): `strings.Index` findet nur das
**erste** Vorkommen je Kommando; und ein legitimer Nutzer-Hook
`~/.claude/hooks/x.sh` trägt den Marker, wird gegen `repoRoot()` aufgelöst und
färbt das Gate **falsch rot**.

**Warum es durchkam:** `TestDurchsetzung_GuardsMeldenDieMutation` fährt
`hookPaths` nur gegen `settings{}` — der einzige Guard im Diff ohne
Mutations-Abdeckung. Das ist die `BEO-018`-Klasse in Reinform.

### F-2 · MEDIUM · `assertDenyDecktGuard` misst Namens-Nennung, nicht Deckung

**Pfad:** `:119–135` (`:122–126`). Die Normalisierung schneidet am ersten
Leerzeichen **oder** `)`. Damit zählen als Deckung für `pip`:
`Bash(pip install *)` (sperrt nur einen Teilbaum), `Bash(pip)` (nur der
argumentlose Aufruf), und ein formloser Eintrag `pip` ohne `Bash(`-Hülle. Die
Zusage in `AGENTS.md` §4 liest sich als Klassen-Deckung; gemessen wird „es
existiert ein Eintrag, dessen erstes Token dem Namen gleicht".
*(Ein Lauf mit verengter Regel wurde vom Auto-Mode-Klassifizierer abgelehnt; der
Befund steht auf der deterministischen Zeichenketten-Logik.)*

### F-3 · MEDIUM · Die Versions-Suffix-Grenze nennt die falsche Schicht

**Pfad:** `:67–69`. Der Kommentar sagt, der Fall werde „über die Sperrliste
separat abgedeckt". Die Sperrliste führt `Bash(python *)`/`Bash(python3 *)` —
`python3.12 x.py` ist davon **nicht** gedeckt, und `MR-047` führt genau diese
Zeile in seiner Grenz-Tabelle. Gedeckt ist der Fall vom **Wächter**. Die
Zuordnung ist vertauscht und erklärt die einzige echte Lücke für geschlossen.

### F-4 · MEDIUM · §5 Abnahme-Punkt 1 ist vorausgesetzt, nicht entschieden

**Pfad:** Slice `:76-82`; Test `:168–183`. §5 verlangt eine Entscheidung, „ob
ein **Repo**-Gate über eine Werkzeug-Datei wachen darf". Weder Commit noch Test
beantworten sie; der Test-Kopf begründet nur den **Ort** per Analogie zu
`gate_consistency_test.go`, dessen Gegenstand aber eine **Produkt**-Konfiguration
ist. Daneben berührt der dritte Test §3 Ausschluss 3 („Kein Gate auf den
Wächter selbst"), zu dem der Commit schweigt: er liest den Wächter und macht
`make gates` rot, wenn dessen `BLOCKED`-Zuweisung fehlt oder umbenannt wird.
Da `make test` aus dem **Arbeitsbaum** baut, färbt schon eine lokale,
uncommittete Änderung an einer Werkzeug-Datei das Repo-Gate rot.

*(Die Bewertung-vs-Kopplung-Unterscheidung des Commits **trägt** — sie
beantwortet aber Ausschluss 2; berührt ist Ausschluss 3.)*

### F-5 · HIGH · Herkunfts-Felder im Kommentar verletzen §3.7

**Pfad:** `:11`, `:167`, `:188`. `:11` trägt zwei Felder und eine Slice-Nummer;
`:188` einen `BEO-`Verweis (nicht im Feld-Schema); `:167` Herkunfts-Prosa mit
Mess-Label. Die Bestandsgrenze greift nicht — Neuzugänge fallen unter den Anker.

### F-6 · LOW · Die Ort-Begründung trifft den Bestand nicht

Die vorhandenen repo-lesenden Tests des Pakets lesen ausnahmslos **YAML, das
dieser Adapter dekodiert**. Die neue Datei liest JSON und ein bash-Skript und
importiert `configyaml` nicht. `arch-check` nimmt `**/*_test.go` aus — ein
Refactoring des Pakets nähme die Prüfung mit, und kein Gate sähe es.

### F-7 · LOW · `matcher`, `type` und `ask` werden dekodiert und nie gelesen

`matcher` von `"Bash"` auf einen anderen Wert gedreht ⇒ der Wächter feuert auf
keinen Bash-Aufruf mehr, alle drei Tests bleiben grün. Die `ask`-Hälfte von
`MR-047` hat gar keine Kopplung.

### F-8 · LOW · `blockedFromGuard`: die erste Fundstelle gewinnt, still

Eine zweite Zuweisung oder eine Fortschreibung wird still unter-gelesen; ein
Kommentar mit Beispiel-Liste in derselben Schreibweise verdeckt die echte.

### F-9 · INFO · Eingaben, die der Test liest bzw. nicht liest (§3.8)

Neben `.claude/settings.json` liegt ein ungetracktes `.claude/settings.local.json`
(27 KB, Top-Level `permissions`/`enabledPlugins`, heute ohne `deny`-Block), dazu
die Nutzer-Ebene außerhalb des Repos. Die Sensors-Zeile spricht von „diese
Dateien". Zweite Achse: gelesen wird der Wächter (`MR-042`-Geltungsbereich),
zitiert nur `MR-047`.

## Negativbefunde

- **Ursachen-Aussage** belegt in drei unabhängigen Stücken: `1fd1d6a` änderte
  `$(cat)` → `$(</dev/stdin)`; Letzteres scheitert gemessen, sobald stdin keine
  wieder-öffenbare Datei ist; die Vorfassung trug dieselbe veraltete Antwortform
  und wurde durchgesetzt.
- **Der neue Leseweg verkürzt nichts** — md5-identisch für Pipe mit/ohne
  Schluss-Newline, Datei-Redirect, mehrzeilig, 1 MiB. Einzige Verkürzung: NUL.
- **Zwei Kanäle, jeder Block-Pfad** — vier Ablehnungs-Pfade, alle über
  `emit_block`; `exit 0` steht genau einmal (Pass-Fall). rc=2 gemessen für Pipe,
  mehrzeiliges JSON, Tiefe >3, leere Eingabe, `/dev/null`, geschlossenes stdin,
  leeres `PATH`. **Kein Pfad gefunden, der ablehnt und mit 0 endet.**
- **`set -e`-Fallen** — `false && cmd` läuft unter `set -e` weiter (gemessen);
  vor dem ersten `emit_block` stehen nur zwei Parameter-Expansionen und das mit
  `|| true` abgesicherte `read`.
- **Verdikt `halb` korrekt konstruiert**; `crash` bleibt erreichbar (mutierter
  Guard: 33 × `crash`).
- **§3.2**, **arch-check**, **testpackage**, **Netzlosigkeit**, **§3.3
  Commit-Zerlegung**, **planning-check**: ohne Befund.

## Kategorie-Summary

HIGH 2 · MEDIUM 3 · LOW 3 · INFO 1

## Urteil

**Schließbar nach Nacharbeit.** F-1 ist blockierend: der zweiten von drei
Zusagen fehlt genau der Fall, für den sie gebaut wurde, die dazu benannte Grenze
beschreibt die umgekehrte Klasse, und der Guard ist der einzige im Diff ohne
Mutations-Abdeckung — gemessen, nicht vermutet. F-5 ist billig zu beheben, steht
aber als HIGH-Anker im Rollen-Skill. F-4 entscheidet, ob der Slice überhaupt in
dieser Form schließt.
