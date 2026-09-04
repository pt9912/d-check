# Review — slice-156 / Commits `9c36964`, `04f7f2e`

**Review-Art:** Code + Doku · **Gegenstand:** `slice-156-hook-antwortform.md`,
`9c36964` (Wächter-Fix + zwei Kanäle + `MR-044`), `04f7f2e` (`AGENTS.md`
aufgeräumt) · **Skill:** `reviewer.md` @ 1.10.0 · **Modell:** claude-opus-5[1m] ·
**Datum:** 2026-08-27

**Selbst gemessen** (Auszug):

| Messung | Ergebnis |
|---|---|
| `set -e; false && cmd` (Test **und** Funktion) | läuft weiter, rc=0 — die AND-Listen sind errexit-exempt |
| `$(</dev/stdin)` je stdin-Form | Pipe ✓, Datei ✓, here-string ✓, `/dev/null` ✓; **stdin geschlossen ✗** (rc=1, **keine Ausgabe**) |
| `IFS= read -r -d '' input \|\| true` | md5-identisch bei Pipe mit/ohne Schluss-Newline, Datei-Redirect, mehrzeilig, 1 MiB; leer → len 0 |
| echter Guard, Rand-Eingaben | Pipe-Block rc=**2**, mehrzeilig rc=2, leer rc=2, `/dev/null` rc=2, **stdin geschlossen rc=2**, Tiefe >3 rc=2, leeres `PATH` rc=2 |
| End-to-End | `echo probe-a; rustup --version` → geblockt, `permissionDecisionReason` kam **wörtlich** an; stderr leer ⇒ das stdout-JSON wurde trotz `exit 2` geparst |
| `make guard-probe` | 33 Proben, 0 Fehlschläge, Exit 0 |
| Gegenprobe (`exit 2` entfernt) | **25 von 33 → `halb`**, exakt reproduziert |
| Gegenprobe (Guard stirbt vor Ausgabe) | 33 × `crash` — das Verdikt ist erreichbar |
| `git show 04f7f2e^:AGENTS.md \| wc -l` | **468** (nicht 452) |

## Findings

### F-1 · MEDIUM · `make guard-probe` meldet Fehlschläge und endet mit 0

**Pfad:** `tools/harness/guard-probe.sh:110`. `$fails` wird nie in den Exit
überführt. Gemessen: „`== Fehlschläge: 25`", **Exit 0**; mit einem Guard, der vor
jeder Ausgabe stirbt: „33", **Exit 0**. Genau die Bauform, gegen die dieser
Slice argumentiert. *(Bestand seit `617fbfa`, aber dieser Slice macht die Proben
zur einzigen Instanz zwischen einem halb kaputten Wächter und „grün".)*

### F-2 · MEDIUM · Die Proben können die Ausfall-Klasse dieses Slice nicht fangen

**Pfad:** `guard-probe.sh:19`, `:91`. Beide Verdikt-Wege bespielen stdin über
eine **Pipe**. Gemessen: `$(</dev/stdin)` funktioniert über eine Pipe
einwandfrei. Der Ausfall wäre unter `make guard-probe` **grün geblieben** — an
dem Tag, an dem der Wächter jeden Befehl durchließ, hätten die Proben 0
Fehlschläge gemeldet. Damit ist §5 Abnahme-Punkt 1 strukturell **offen**.

### F-3 · MEDIUM · Neue Regel in einer Commit-Botschaft, im Widerspruch zu einer aktiven Adaption

**Pfad:** `04f7f2e` gegen `MR-013-…md:22` und `AGENTS.md:154-156`. Die Botschaft
begründet eine Entfernung mit *„In AGENTS.md und harness/README.md stehen nur
Normatives und ADRs"*. `MR-013` schreibt vor, dass der Lifecycle-Move-Commit
*„alle Pfad-Verweise auf den Slice (Roadmap, `AGENTS.md` §4,
`harness/README.md` §Sensors)"* mitzieht. Gemessen: beide Dateien tragen heute
**null** Slice-Verweise. Entweder ist die MR-013-Klausel toter Text, oder die
Begründung ist falsch; aufgelöst wird der Widerspruch in keiner Richtung, und
die neue Regel steht nur in einer Commit-Botschaft.

### F-4 · MEDIUM · Der Vorfall hat keinen Eintrag im Beobachtungs-Register

Das Register führt BEO-001…017; keine Zeile trägt die Klasse *„eine Härtung am
Leseweg kippt die Fehlerpolitik, und die Proben sehen es nicht, weil sie den
Kanal anders bespielen als der Ernstfall"*. `BEO-007` und `BEO-017` liegen
daneben, decken es aber nicht. *(Zeitlich: Register-Einträge entstehen hier bei
der Closure — als Closure-Pflicht zu führen, nicht als Versäumnis.)*

### F-5 · LOW · `harness/README.md` schreibt `MR-044` eine Tabelle zu, die `MR-042` trägt

### F-6 · LOW · Zahl in der Botschaft `04f7f2e` stimmt nicht

„452 → 424 Zeilen"; gemessen **468 → 424** (Diff-Stat 26+/70− stützt 468).

### F-7 · LOW · Eine entfernte §3.4-Aussage steht nirgends sonst

Gedeckt sind Token-Breite, die vier benannten Grenzen, `modul-pfad` als
Token-Ziel, das Rohzeilen-Lesen, die nachgebaute Abwärts-Kante. **Nicht**
gedeckt: *„ein bares `ADR-<NNNN>` im Fließtext ist es **nicht** — die Klasse
trägt kein `token`"* — die `adr`-Klasse trägt in der Config keinen Kommentar.
Ebenfalls ersatzlos: „Großschreibung fällt heraus", „ein Pfad über einen
Zeilenumbruch".

### F-8 · LOW · Mess-Gütesiegel im Skript-Kommentar, vier Minuten nach seiner Entfernung aus `AGENTS.md`

`.claude/hooks/pretooluse-command-guard.sh:51` — *„Beide zugleich sind
**gemessen** konfliktfrei"*.

### F-9 · LOW · Kopf-Kommentar der Proben ist eine Fassung hinterher, Verdikt-Logik doppelt

Der Kopf sagt „das **dritte** Verdikt `crash`", während derselbe Commit das
vierte einführt; der Verdikt-Block steht **wörtlich zweimal** im Skript.

### F-10 · INFO · `read -d ''` verkürzt an einem NUL-Byte

Gemessen: 9-Byte-Eingabe mit NUL an Position 5 wird als 4 Bytes gelesen. NUL
**im** `command`-Wert → abgeschnittenes JSON → Exit 3 → **Block** (richtig); NUL
**zwischen** zwei vollständigen Objekten → das zweite fällt weg → **pass**.
JSON verbietet rohes NUL; die Eingabe ist maschinell erzeugt. Der Kommentar
begründet ausführlich, *warum nicht* `cat` und *warum nicht* `$(</dev/stdin)`,
nennt aber die eigene Grenze des gewählten Wegs nicht.
*(Nebenbei: der `read`-Weg kostet 9 ms Grundlast, 376 ms/MiB über eine Pipe.)*

## Negativbefunde

- **Ursachen-Aussage belegt, nicht erzählt** — in drei unabhängigen Stücken.
  Nicht reproduzierbar war nur die konkrete stdin-Form des realen Hooks (siehe
  F-2).
- **Der neue Leseweg verkürzt nichts** außer am NUL-Byte.
- **Zwei Kanäle, jeder Block-Pfad** — vier Ablehnungs-Pfade, alle über
  `emit_block`; `exit 0` steht genau einmal. **Kein Pfad gefunden, der ablehnt
  und mit 0 endet.**
- **`set -e`-Fallen** — vor dem ersten möglichen `emit_block` stehen nur zwei
  Parameter-Expansionen und das mit `|| true` abgesicherte `read`. Ein harter
  Lesefehler lässt `input` unset; `set -u` tötet die awk-Subshell, rc≠0 → Block.
- **Verdikt `halb` korrekt konstruiert**; die Umsortierung (`grep` vor `rc`) war
  nötig und ist richtig.
- **Zusagen gegen Bestand** — `AGENTS.md` §3.1 stimmt; kein `guard-probe`-Aufruf
  in CI oder Hooks; „zehn Glieder" stimmt. Von den drei Mess-Behauptungen in
  `MR-044` eine unabhängig reproduziert; „25 von 33" exakt nachgestellt.
- **Slice-Kopf und Struktur**, **DoD 1–3**, **§3.3**, **`planning-check`**:
  ohne Befund.

## Kategorie-Summary

HIGH 0 · MEDIUM 4 · LOW 5 · INFO 1

## Urteil

**Schließbar nach Nacharbeit.** Der Kern trägt: Ursachenaussage belegt, neuer
Leseweg verkürzt nichts Erreichbares, alle vier Ablehnungs-Pfade enden
nicht-null, beide Kanäle wirken nachweislich zusammen. Was fehlt, liegt eine
Ebene darüber: die Proben, die diese Zusage halten sollen, haben **keinen**
eigenen Exit-Kanal (F-1) und bespielen stdin anders als der Ernstfall (F-2) —
zusammen genau die Konstellation, die den Ausfall einen Tag lang unsichtbar
gemacht hat.
