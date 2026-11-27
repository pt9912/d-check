# Review slice-296 — Entscheidungsvorlage zu den Neuerungen von `v7.0.0`, Runde 1

- **Review-Art:** Doku (gegen Plan, Architect-Verdikt, Hard Rules); der Gegenstand ist eine Entscheidungsvorlage
- **Gegenstand:** `slice-296`, Diff `5a14fe90..433447b1` (Plan-Änderung `a6c71fec`, neue Datei `433447b1`:
  `docs/plan/planning/baseline-v700-neuerungen.md`)
- **Skill:** `reviewer.md` @ 1.20.0
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-11-27
- **Eingangs-Kontext:** `slice-296` (§1, §2 nur als Prüfmaßstab der Klassen und Belege, §3 samt Plan-Änderungen, §4, §6),
  `wave-98` (§1, §3 (f), §5, §6), Architect-Verdikt zu `wave-98` (Frage 11), `AGENTS.md` §3 und §5,
  Report zu `slice-295` (Runde 1), Baum `v7.0.0` gegen den alten Baum aus `git archive 738d2bca`
- **Zusätzlich gefahren** (Wegwerf-Bäume im Scratchpad, am Repo nichts geändert): Bundle-Zählung und Diff je Datei,
  zusätzlich mit herausnormalisierten Namens-Tokens; echter `git worktree` unter `.claude/worktrees/` gegen d-check
  (`d-check:latest`) und a-check (Digest aus `a-check.mk`); `planning` mit namensbasierten Wellen; `--suggest-config`
  mit beiden Spezifikations-Namen; der vorgeschlagene `pinsync`-Pin an Kopien beider Dockerfiles;
  `printenv GOLANG_VERSION` am Digest aus `tools/archive-wave/Dockerfile`.
- **Nicht gefahren:** die zwei Wegwerf-Läufe mit `make test` (B1 „21 Tests in 2 Paketen rot", B2 „grün"), `make gates`,
  `semgrep` mit Arbeitsbaum; ob der Digest des Wurzel-`Dockerfile` zur Version seines Tags passt.

## Findings

Keine HIGH. Sieben MEDIUM, zwei LOW, drei INFO.

### F-1 — MEDIUM — 5b: kleinste Antwort und Preis stimmen nicht; mit ihr bleiben zwei Gates rot

- **quelle:** Prüffragen 8 und 17, `AGENTS.md` §5 Regel 13
- **pfad:** Vorlage, Zeile 5b · „je 1 Zeile; das Closure-Profil hat keinen `scan`-Block und wäre eine dritte Stelle"
- **befund:** Gemessen mit echtem `git worktree` unter `.claude/worktrees/w1`, `.gitignore`-Zeile und `scan.ignore`-Eintrag:
  d-check meldet weiter **100** `citation-out-of-range`, alle aus dem Arbeitsbaum (`citations.scope` ersetzt den
  globalen Scan; erst ein Eintrag auch dort ergibt 0), und a-check endet mit Exit 1 („Schicht in layers deklarieren oder
  Datei in exclude aufnehmen"), also `make arch-check` in `make gates`. Das Closure-Profil ist dagegen **keine** dritte
  Stelle: ohne `scan.roots` gelten `docs` und `spec`, es prüft mit und ohne Kopie 1191 Dateien.
- **Failure-Szenario:** Der Auftraggeber gibt „zwei Zeilen" frei; der Folge-Slice liefert genau das, und der erste
  Arbeitsbaum färbt `make gates` rot. „d-check meldet einen Befund" stammt aus einem Baum, der die Größe nicht trägt
  (volle Kopie ohne Ausschluss: 4054 Befunde).
- **verifizierbar:** ja — Arbeitsbaum anlegen, `make doc-check` und `make arch-check`
- **klasse:** preis-an-zu-kleinem-gegenstand-gemessen

### F-2 — MEDIUM — 7a und 7b geben die Regeln ohne ihre Bedingung wieder

- **quelle:** Prüffrage 9, `AGENTS.md` §5 Regel 13 und 16, `slice-296` §6 Risiken 3 und 5
- **pfad:** Vorlage, Zeile 7a · „Prüffrage *Berührt der Schnitt eine Vertrauensgrenze?*, fehlender Ausgang ist MEDIUM";
  Zeile 7b · „Was eine aktive ADR an der Grenze akzeptiert und nennt, meldet der Reviewer nicht"
- **befund:** `v7.0.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten stellt die Frage ausdrücklich „für
  **Plan- und Design-Review**"; die Zeile und ihr Wortlaut 7a (beginnt eine Zeile unter der Bedingung) lassen das weg.
  §Finding-Kategorien sagt „meldet der Reviewer in keiner Review-Art **als Sicherheits-Anti-Pattern**" und „Jeder
  Sicherheits-Verstoß bleibt HIGH"; beides fehlt.
- **Failure-Szenario:** Der Auftraggeber liest 7a als Frage in jedem Review (auch am Code-Diff, den dieses Repo
  überwiegend prüft) und 7b als pauschales Schweigen zu allem, was eine ADR nennt. Er entscheidet über zwei andere
  Regeln als die der Baseline.
- **verifizierbar:** ja — Lesen der zwei Absätze
- **klasse:** grenz-aussage-ohne-ihre-bedingung

### F-3 — MEDIUM — 7b/17b: die Ausnahme hat im Repo-Skill keine Regel, von der sie ausnimmt

- **quelle:** `slice-296` §2 (kleinste Antwort mit Preis), `AGENTS.md` §5 Regel 13
- **pfad:** Vorlage, Zeile 7b · „**Kleinste:** Ausnahme in [`reviewer.md`], rund 4 Zeilen"
- **befund:** `.harness/skills/reviewer.md` führt keine Klasse „Sicherheits-Anti-Pattern"
  (`grep -i 'sicherheits-anti\|injection' .harness/skills/reviewer.md` ist leer); die Vorlage des Bundles führt sie als
  HIGH, der Repo-Skill nicht. Die Vorlage sagt das nicht.
- **Failure-Szenario:** Freigegeben wird eine Ausnahme ohne Gegenstand; die eigentliche Frage (bekommt der Skill den
  HIGH-Anker, zu dem die Ausnahme gehört?) wird dem Auftraggeber nicht gestellt.
- **verifizierbar:** ja — das `grep` oben
- **klasse:** antwort-setzt-regel-voraus-die-fehlt

### F-4 — MEDIUM — Neuerung fehlt: die Vorlagen-`.d-check.yml` benennt eine Grenze des Produkts

- **quelle:** Prüffrage 18, `DC-FA-PLAN-001`, `slice-296` §1 (Reste der Namensregel für Konsumenten)
- **pfad:** Vorlage, Zeilen 2d und B1 · „mit beiden Globs `wave-drift` und `wave-results-missing`"
- **befund:** Der neue `planning`-Block in `v7.0.0` · `templates/.d-check.yml` trägt eine „GRENZE (d-check v0.91.0, auf
  v0.92.0 nachgemessen)": Eine Wellen-Kennung ohne Ziffern (`wave-cache-warmup`) liest `waves` nicht, Datei ohne Zeiger
  und Zeiger ohne Datei bleiben still. Nachgemessen mit **beiden** Globs und `mode: many`: für `wave-ohne-zeiger.md` und
  einen Zeiger `wave-cache-warmup` ohne Datei kein `wave-drift`; die nummerierte Gegenprobe wird gemeldet. Keine Zeile
  der Vorlage (2d „bereits befolgt", B1 bis B3) nennt das; das Handbuch nennt die Ziffern-Bedingung auch nicht
  (`grep -i ziffer` im Wellen-Abschnitt leer, die Spezifikation führt sie in den Schritten zu `planning.waves`).
- **Failure-Szenario:** Der Auftraggeber nimmt den Handbuch-Satz aus B1 und hält den Konsumenten, der dem Regelwerk
  folgt, damit für versorgt; ein Konsument mit Namens-Kennungen bleibt trotz beider Globs ungeprüft. Es ist die einzige
  Stelle des Bundle-Diffs, die ich ohne Zeile gefunden habe.
- **verifizierbar:** ja — Wegwerf-Baum mit namensbasierter Welle, `make planning-check`-Konfiguration
- **klasse:** grenzen-liste-ohne-groesste-luecke

### F-5 — MEDIUM — die Check-Zeilen heißen „nicht anwendbar", obwohl das Kapitel das Docker-Repo ausdrücklich meint

- **quelle:** Prüffragen 9 und 20, `AGENTS.md` §5 Regel 13
- **pfad:** Vorlage, Zeile A6 · „wer die Zusage geprüft haben will, braucht ihn"; Zeilen 1b, 3, A6, A9 bis A20, 18a,
  18c, 28
- **befund:** Keine Pflicht ist richtig gelesen, „nicht anwendbar" nicht: `v7.0.0` ·
  `regelwerk/modul-14-docker-harness.md` nennt den Check im Docker-only-Repo („Auch im Docker-only-Repo prüft er
  etwas, das der Digest nicht sagt"; Ausführungsort „nur Host (Docker-only …)"). Er ist anwendbar und wahlfrei. Zweitens
  zählt §Zwei Prüfungen auch den Stellen-Abgleich als Prüfung derselben Zusage, den das Repo mit `pinsync` fährt; der
  Check wird nur gebraucht, um **installierte** Werkzeuge gegen den Pin zu halten. 1c („neu und anwendbar") verweist für
  seine zweite Hälfte auf das „nicht anwendbare" A9.
- **Failure-Szenario:** A6 wird abgelehnt in dem Glauben, das Kapitel betreffe dieses Repo nicht. Dabei ist der Fund aus
  A8 genau die Klasse, die das Kapitel allein dem Check zuschreibt („dass darin die Version steckt, die der Pin
  behauptet"); A6 sagt das nicht. (Am Wurzel-`Dockerfile` meldet `make go-base-digest`, Netz und fail-open, die
  verwandte Abweichung Tag gegen Digest.)
- **verifizierbar:** ja — Lesen des Kapitels
- **klasse:** klasse-widerspricht-dem-wortlaut

### F-6 — MEDIUM — „Bedingung: … (Wortlaut X)": das zitierte Wort ist keine Bedingung

- **quelle:** Prüffrage 9, `slice-296` §2 („Grund an der Bedingung der Regel") und §6 Risiko 3
- **pfad:** Vorlage, Zeile A10 · „Bedingung: es gibt ein Gate, das Ist gegen Soll vergleicht (Wortlaut A10)"
- **befund:** Wortlaut A10 („Das Gate vergleicht Ist gegen Soll."), A12 („Der Check prüft, was dort läuft …"), A13
  („Das Laufzeit-Gate ist ein Skript des Repos"), A15 für Zeile 3 („Muster, kein Skript"), 28 („hat eine maschinelle
  Vorstufe"), 18c („ergänzt den Wächter") und A6 für Zeile 1b beschreiben den Check; keiner sagt „gilt nur, wenn". 18a
  verweist auf A6. Eine echte Bedingung steht in `v7.0.0` · `templates/AGENTS.template.md` („führt das Repo einen
  Tool-/Versions-Check", Wortlaut 10a) und dort nur für den Inventar-Satz.
- **Failure-Szenario:** Die Zitate sehen wie Belege aus, `make doc-check` bestätigt sie wörtlich, und der Auftraggeber
  hält die Nicht-Anwendbarkeit für vom Kapitel verfügt, wo sie eine Lesart der Vorlage ist.
- **verifizierbar:** nein — Urteil; `citations` prüft nur Wortgleichheit
- **klasse:** citation-stretched-beyond-scope

### F-7 — MEDIUM — A5 „befolgt": `harness/image-hash.txt` ist kein Beispiel

- **quelle:** Prüffrage 20, `ADR-0002` §Fitness Function, `slice-296` §6 Risiko 1
- **pfad:** Vorlage, Zeile A5 · „wie das Kapitel sie als Beispiel nennt, führt das Repo nicht"
- **befund:** `v7.0.0` · `regelwerk/modul-14-docker-harness.md` §Multi-Stage-Build führt unter „verlangt:" das
  „einzeilige Beleg-Artefakt `harness/image-hash.txt`, referenziert in `harness/README.md`". `ADR-0002` (Accepted) nennt
  dieselbe Datei als Beleg-Manifest; sie fehlt. Die Regel ist alt, nicht neu in `v7.0.0`; die Antwort „befolgt" trägt sie
  trotzdem nicht. A2 trägt nur, weil die Digests im Commit stehen; der Image-Hash geht allein nach stdout.
- **Failure-Szenario:** Der Auftraggeber liest den Docker-Träger als erfüllt; eine verlangte und in einer eigenen ADR
  angekündigte Datei kommt nie auf den Tisch.
- **verifizierbar:** ja — `ls harness/`, Lesen von `ADR-0002`
- **klasse:** begruendung-trifft-gegenstand-nicht

### F-8 — LOW — A1 nennt `make tidy` nicht

- **quelle:** `AGENTS.md` §5 Regel 13
- **pfad:** `Makefile`, Target `tidy` · „-v "$(CURDIR)":/src -w /src golang:$(GO_VERSION)"
- **befund:** Der Lauf nimmt das Go-Image nur per Tag, ohne Digest, und schreibt `go.mod`/`go.sum`. A1 nennt als
  Ausnahmen der Container-Seite nur A8 und (Host) A7.
- **verifizierbar:** ja — `grep -n 'golang:' Makefile`
- **klasse:** grenzen-liste-ohne-groesste-luecke

### F-9 — LOW — 5a kürzt das Zitat aus `implement-slice.md` um seine Bedingung

- **quelle:** `AGENTS.md` §5 Regel 13
- **pfad:** Vorlage, Zeile 5a · „schließt heute ein paralleles `in-progress/` aus"
- **befund:** Die Stelle lautet „WIP-Limit = 1 **pro Implementer** (Modul 5): kein paralleles `in-progress/`"; sie
  schließt keinen zweiten Implementer-Lauf aus. Die Empfehlung („nichts tun") ändert sich dadurch nicht.
- **verifizierbar:** ja
- **klasse:** grenz-aussage-ohne-ihre-bedingung

### F-10 — INFO — Zeilen 22, 31 und 32 tragen zwei Neuerungen

Je eine inhaltliche Änderung und der entschiedene Namens-Nachzug in einer Zeile mit einer Klasse (`slice-296` §6
Risiko 2: „keine Zeile mehr als eine Neuerung"). Auf die Entscheidung wirkt es nicht, weil der Namensteil nichts fragt.
**klasse:** zeile-buendelt-neuerungen

### F-11 — INFO — Lesbarkeit ohne Sitzungsgeschichte

„Übernommen ist hiermit nichts" (Kopf) steht neben „0 Dateien; die Regel wirkt über den Pin" (4b, 5c, 13a): Dass ein
leeres Entscheidungsfeld dort die Regel gelten lässt, muss man sich aus §So lesen erschließen. „Closure-Profil",
„Break-Test", „Waise" und „Wellen-Präfix" sind nicht erklärt. Die Zellen von A6, A8 und B1 sind Absätze. Zu 7a: Mit 4a
und 7a zusammen wird jede Bedrohung ohne Ausgang zum Befund; die als „größer" geführte Durchsicht in 4a folgt dann.
**klasse:** vorlage-setzt-sitzungswissen-voraus

### F-12 — INFO — A8: der neue Pin bewacht die Versions-Zeilen von `tools/archive-wave/` nicht

Der Fund selbst trägt (siehe unten). Die genannte Grenze nennt Tag gegen Digest; nicht genannt ist, dass
`ARG GO_VERSION` und `GO_VERSION ?=` dort auch nach der Antwort ohne Quelle im Pin `GO_VERSION` stehen.
**klasse:** grenzen-liste-ohne-groesste-luecke

## Belege der Zeilen „bereits befolgt" und „befolgt"

| Zeile | Trägt? | Geöffnet |
|---|---|---|
| 2d | ja (die Grenze des Blocks fehlt, F-4) | `.d-check.yml` `planning.waves`: beide Globs `wave-*`, `mode: many`, Überschrift und Marker wie in der Vorlage |
| 2e | ja | `.d-check.yml` `pinsync`, zwei Pins; Lauf meldet „2 Pin(s) geprüft, 4 Quelle(n)" |
| 9b | ja | `wave-98` §4 verlinkt mit Lifecycle-Verzeichnis; `tools/harness/slice-mv.sh` zieht Verweise nach |
| 10a | ja | `AGENTS.md` §3.1 nennt `git`, `make`, `bash`, Docker und die POSIX-Werkzeuge als Klasse |
| A1 | ja, mit ungenannter Ausnahme (F-8) | `FROM`-Zeilen im `Dockerfile`, Digests in `tools/semgrep.sh`, `tools/image-scan.sh`, `a-check.mk`, `ADR-0011` |
| A2 | ja | `make versions`-Rezept; `fullbuild` gibt den Image-Hash aus |
| A5 | **nein**, nicht vollständig (F-7) | Stages, Digests, `go.sum` ja; Beleg-Artefakt nein |

## Der Fund in A8 am Gegenstand

Er trägt in allen Teilen: `make archive-wave-test` und `make archive-wave` reichen `GO_VERSION=1.27.4` durch; die
`FROM`-Zeile in `tools/archive-wave/Dockerfile` hält den Digest `sha256:4013…`, an dem `printenv GOLANG_VERSION`
**1.27.0** ausgibt; `make go-base-digest` liest nur das Wurzel-`Dockerfile`; ein `pinsync`-Pin über beide Basis-Digests
meldet an Kopien heute `pinsync-mismatch`.

## Antworten auf die vier Urteilsfragen

- **(a) Zweiter Träger Pflicht?** Nein. „Die Umgebungs-Zusage lässt sich auf zwei Wegen erfüllen", „Die Wahl trifft das
  Repo", und die `AGENTS`-Vorlage setzt den Check bedingt. Pflicht ist die Zusage als Anforderung mit einem Sensor als
  Messmethode (A3). „Wo Docker festlegt, ist der Check die Prüfung der Zusage" beschreibt seine Rolle und verlangt ihn
  nicht. Anwendbar ist er trotzdem; richtig wäre „anwendbar, wahlfrei, Empfehlung nicht übernehmen" (F-5).
- **(b) „Nicht anwendbar" ohne Kapitel-Bedingung:** die Zeilen aus F-6. 2b zitiert eine echte Bedingung; die Zeilen mit
  „entschieden" brauchen nach `slice-296` §3 Schritt 3 keine; 1d, A4, 6, 36, 39 und 40 verlangen keinen Akt.
- **(c) Bündel:** F-10. 1b fächert in eigene Zeilen auf und ist unschädlich; die elf Zeilen unter A6 sind vom Plan
  verlangt (je Bestandteil eine).
- **(d) In einer Sitzung entscheidbar?** Ja. Die Entscheidungsfläche ist die Tabelle mit sechs Gegenständen, und ich
  habe alle Zeilen in einer Sitzung geprüft; die Rückführung aus `slice-296` §4 greift nicht. Vor der Vorlage an den
  Auftraggeber stimmen aber bei „Arbeitsbäume" (F-1), „Bedrohungen" (F-2, F-3), „Umgebungs-Zusage" (F-5) und „Reste der
  Namensregel" (F-4) die Grundlagen nicht. Die Behebung braucht keine neuen Zeilen außer der zu F-4.

## Negativbefunde

- **Zählung:** geprüft, ohne Befund — 54/55, 47, 38, 862, 31 (19 + 12), 804, 7 umbenannt, 1 neu, 7 nur Stempel;
  Dateinummern 40, Unterzeilen-Zählkommando liefert die genannte Form.
- **Vollständigkeit des Bundle-Diffs:** `modul-14`, `.d-check.yml`, `modul-04`, `modul-05`, `modul-10`, `modul-03`,
  `modul-13`, `modul-12`, `modul-08`, die `AGENTS`-, Sicht-, Slice-, Reviewer-, ADR- und Lastenheft-Vorlage und der
  Vorlagen-Index gelesen; ohne Befund außer F-4. Die als „nur Namen" geführten Dateien tragen nach Herausnormalisieren
  der Namens-Tokens keinen weiteren inhaltlichen Unterschied (Rest je eine Beispiel-Kennung).
- **Frage B, Reichweiten:** geprüft, ohne Befund — 10 Zeilen in 4 Dateien plus 2 (`welle-<n>`), 8 Beispiel-Zeilen in
  4 Dateien, 0 Treffer `archive.zip`; `--suggest-config ai-harness` schlägt neben `spec/specification.md` den deutschen
  Namen vor, damit 0 Befunde, mit dem richtigen Namen `matrix-downward`. Nicht geprüft: die zwei `make test`-Aussagen.
- **4a:** geprüft, ohne Befund — `git grep -i vertrauensgrenze -- spec` leer, die Sicht führt `ARC-001` bis `ARC-007`.
- **A3, A7, 1d, 2b, 2c, 12a:** geprüft, ohne Befund — vier `DC-QA-*`, keine sagt die Zusage zu; Workflows auf
  `ubuntu-latest`; kein lebender Link auf den alten Anker; kein `d-check.mk` eingebunden, `exempt-targets` leer.
- **Wortlaute:** geprüft — ein Lauf über `git archive HEAD` meldet 0 Befunde, die Zitate sind wortgleich. Ihre Reichweite
  ist Gegenstand von F-2 und F-6.
- **Plan-Änderung `a6c71fec`:** geprüft, ohne Befund — vor dem Commit der Vorlage, nennt Stand und Anlass.
- **Keine Übernahme eingeschleust:** geprüft, ohne Befund — der Diff berührt nur die Vorlage und den Plan.
- **Hard Rules** (`AGENTS.md` §3): geprüft, ohne Befund — kein Code, keine Konfiguration, keine ADR, keine Spec im Diff.

## Zu den Risiken aus `slice-296` §6 (Eingabe für die Closure)

Risiko 1: ein Beleg trägt nicht vollständig (F-7). Risiko 2: drei Zeilen bündeln (F-10). Risiko 3: Bedingungen fehlen
oder sind nicht zitiert (F-2, F-6, F-9). Risiko 4: am geprüften Diff nicht eingetreten. Risiko 5: Zitate wortgleich,
Reichweite überdehnt (F-2, F-6).

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 7 (F-1 bis F-7) | 2 (F-8, F-9) | 3 (F-10 bis F-12) |

## Verdikt

**Nicht abnahmefähig in dieser Fassung** — sieben MEDIUM blockieren. Eingearbeitet werden nach der Review-Vorgabe der
Welle nur die MEDIUM; höchstens eine weitere Runde. Die Empfehlungen selbst habe ich in keiner Zeile widerlegt; falsch
oder verkürzt sind Preis, Bedingung, Klasse oder Beleg, auf denen sie stehen.
