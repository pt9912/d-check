# Verifikation — slice-296: Neuerungen von `v7.0.0` sichten, Entscheidungsvorlage ohne Übernahme

- **Rolle:** Verifier (Modul 11) — DoD und Plan gegen den Stand
- **Gegenstand:** `slice-296` §2, §3 (samt den Plan-Änderungen), §5, §6; Liefergegenstand
  `docs/plan/planning/baseline-v700-neuerungen.md` (die Vorlage); Closure-Trigger (b) und (f) von `wave-98`
- **Range:** `5a14fe90..582f10b3`; `main` liegt 66 Commits vor der lokalen Referenz `origin/main`. Uncommitted lagen die
  Reports der Runden 1 und 2 vor
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-11-27
- **Wegwerf:** Klone und Fixtures im Scratchpad, Images unter eigenem Namen; am Repo ist allein dieser Report neu.
  Die schweren Läufe liefen nacheinander (Last zwischen 6 und 18 auf 20 Kernen, 14 GiB verfügbar), fremde Container
  blieben unberührt

## Verdikt

Die Messungen der Vorlage halten, auch die bisher nicht nachgefahrenen: `make gates` im Klon mit `git worktree` und den
vier Zeilen aus 5b ist grün, die zwei `make test`-Aussagen in B1 und B2 treffen. Die Abgrenzung „keine Übernahme“ ist
eingehalten. **Zwei MEDIUM** sind vor der Closure zu entscheiden: Zeile A1 führt „befolgt“, obwohl jeder Gate-Build ein
nicht digest-gepinntes Image bezieht (V-1, betrifft den Slice), und die Vorlage fällt aus der erwarteten Menge des
Closure-Triggers (b) der Welle (V-2, betrifft die Welle).

## Messungen (selbst gefahren)

| # | Lauf | Ergebnis |
|---|---|---|
| M1 | `make gates` im Repo, samt den zwei uncommitteten Reviews | `EXIT=0`; `[gates] baseline-verify + … + planning-check green`; `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%`; `Ran 55 rules on 75 files: 0 findings`; Doku-Lauf `1340 Datei(en) geprüft, 0 Befund(e)`, `pinsync: 2 Pin(s) geprüft, 4 Quelle(n) gelesen` |
| M2 | Klon, `git worktree` unter `.claude/worktrees/w1`, **ohne** Ausschlüsse | d-check `4054 Befund(e)`; a-check Exit 1 |
| M3 | dasselbe, nur `scan.ignore` und `.gitignore` | `100 Befund(e)`, alle `citation-out-of-range` und alle aus dem Arbeitsbaum |
| M4 | dasselbe, **alle vier** Zeilen (3 Dateien) | d-check 0, a-check 0 |
| M5 | `make gates` im selben Klon (eigener Image-Name) | `EXIT=0`, Schlusszeile green; Coverage 94.90 %, semgrep 0 findings. **5b ist bestätigt** |
| M6 | Klon, nur die 4 Code-Zeilen der Defaults auf `wave-*`, `make test` | `EXIT=2`; **21** oberste `--- FAIL`-Tests, **2** Pakete (`configyaml`, `rules`). **B1 bestätigt** |
| M7 | Klon, die 8 Zeilen aus B2 auf `specification.md`, `make test` | `EXIT=0`. **B2 bestätigt** |
| M8 | `planning`-Fixtures, `waves.dir` und `mode: many` | Default-Globs mit `wave-99` und Register-Zeile ohne Notiz: 0; beide Globs: `wave-drift` und `wave-results-missing`; Namens-Kennung `wave-cache-warmup`: 0; Gegenprobe `wave-7`: `wave-drift` (2f, B1 bestätigt) |
| M9 | Spiegelbild zu B1 (in der Vorlage „nicht eigens gemessen“) | `wave-*`-Globs, `welle-99` ohne Zeiger und `welle-98` ohne Notiz: 0 Befunde. Die Vermutung trifft |
| M10 | `--suggest-config ai-harness`, beide Spezifikations-Namen | schlägt beide Male den alten Namen vor; bei `specification.md` Abwärtsverweis 0 Befunde, mit korrigiertem Namen `matrix-downward`. Spiegelbild zu B2 (alter Name, neuer Vorschlag): 0 Befunde, trifft |
| M11 | A8 | `printenv GOLANG_VERSION` am Digest `sha256:4013…`: `1.27.0`; `pinsync` über beide Basis-Digests an Kopien: `pinsync-mismatch` |
| M12 | `make versions` | nennt `FROM`-, semgrep-, a-check-, Builder- und Trivy-Digests; kein Frontend (V-1) |
| M13 | `make trace-check RANGE=2831e071..HEAD`, `make adr-check RANGE=738d2bca..HEAD`, `make review-coverage` | 8 bzw. 60 Commits aufgelöst, je 0 Befunde |

## Inventar nachgezählt

Form: **Datei im Bundle** (`regelwerk/`, `templates/`) ohne `SHA256SUMS`; alt aus
`git archive 738d2bca`, neu aus dem vendorten Baum; „inhaltlich“ ohne `<!-- Quelle:`-Zeilen, Versionsstempel
normalisiert; Diff-Zeilen sind die `<`/`>`-Zeilen von `diff`.

| Größe | Vorlage | nachgezählt |
|---|---|---|
| Dateien alt / neu / gleichnamig | 54 / 55 / 47 | gleich |
| byte-verschieden / roh | 38 / 862 | gleich |
| inhaltlich / normalisiert | 31 (19 + 12) / 804 | gleich |
| umbenannt / neu / nur Stempel | 7 / 1 / 7 | gleich; die sieben Stempel-Dateien je 2 rohe Zeilen |
| Diff-Zahl je Zeile und Reihenfolge | absteigend | alle 39 gleich, auch die umbenannten gegen den alten Namen |

Die Dateinummern 1 bis 39 bilden die 39 Dateien eineindeutig ab, 40 ist die Sammelzeile. Alle 68 Inventar-, 20 A- und
3 B-Zeilen haben eine leere Entscheidungs-Spalte.

**Zählkommando der Welle, Trigger (f):** Er verweist auf „das Kommando aus `slice-296` §3“, dort stehen jetzt zwei.
`grep -c '^| [0-9]'` ergibt **68** (mit Unterzeilen) und passt nicht zu „39 Dateien“;
`grep -o '^| [0-9][0-9]*' <Vorlage> | sort -u | wc -l` ergibt **40** (39 plus Sammelzeile) und passt. Die Welle sollte
das zweite Kommando und die erwartete 40 nennen (V-3).

## DoD je Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| Inventar | **bestätigt** | Zählung oben. Die Diffs der als „nur Namen“ geführten Dateien selbst gelesen: kein weiterer inhaltlicher Unterschied. Die Zeilen 22, 31, 32 tragen zwei Änderungen (Runde 1, F-10) |
| Klasse und Empfehlung je Zeile | **teilweise** | Spalten, Klassen, Wortlaute (44 `d-check:cite`, M1 grün) stehen. Abweichungen: A1 (V-1), „Anker-Link“ (V-4), zwei Bedingungen (V-5) |
| Zwei benannte Fragen | **bestätigt**, A1 ausgenommen | A: jeder Bestandteil aus §3 Schritt 4 hat eine Zeile (A1 bis A3 Zusage, A5 Docker, A6/A7 zweiter Träger, A8/A9 zwei Prüfungen, A10 bis A13 Soll, Prüfform, Inventar, Kontext, A14 Probe), die sechs Abschnitte der Vorlage sind A15 bis A20. B: Defaults, Beispiel-Pfade, Konsumenten-Verhalten mit Reichweite und beiden Preisen; 10 Zeilen in 4 Dateien plus 2, und 8 Zeilen in 4 Dateien, nachgezählt; keinen weiteren Rest in Handbuch, `--print-config` und `--suggest-config` gefunden. „Lücke (wahlfrei)“ weicht von den drei Antworten des Plans ab (Runde 2, R2-2) |
| `make gates` grün | **bestätigt** | M1 |
| Review durchgeführt | **bestätigt** | zwei Reports unter `docs/reviews/`, noch uncommitted; F-1 bis F-7 geschlossen, höchstens zwei Runden eingehalten. Die Kontext-Trennung ist aus den Dateien nicht prüfbar |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | **offen** | Closure steht aus; §7 ist leer. `make verify-closure-notes` greift erst nach dem Move |

**Bewusstes Brechen:** Der Slice liefert keinen Test. Wo die Vorlage ihre Aussage auf eine Messung stützt, ist das Rot
oben mit Ursache gelesen (M2, M3, M6, M8, M10, M11).

**Plan-Konformität:** Beide Implementer-Plan-Änderungen liegen vor ihrem Commit (`a6c71fec` vor `433447b1`,
`b620261c` vor `582f10b3`). Die Evidenz-Datei im Register (§3, Zeile 2 der Tabelle) steht mit der Closure aus.

**Abgrenzung „keine Übernahme“:** `git log --name-status 5a14fe90..582f10b3` nennt allein die Vorlage (`A`, `M`) und den
Slice-Plan (`M`). `open/` und `next/` sind leer, es gibt keinen neuen MR-Eintrag. Die zwei Move-Commits davor ziehen
Link-Pfade nach (darunter in `MR-081`) und den Ruhe-Marker der Roadmap.

## Stichprobe der Zeilen, die kein Review namentlich geprüft hat

Geöffnet gegen den Diff der Bundle-Datei bzw. den Repo-Beleg: 1a, 1c, 2a, 4c, 5c, 5d, 6, 7c, 8, 9a, 10b bis 10d, 11,
13a, 13b, 14, 15a, 15b, 16, 17a, 17b, 19 bis 21, 23 bis 26, 27a, 27b, 29, 30, 33 bis 40, A4, A11, A14, A16 bis A20.
Ohne Befund außer den unten genannten. 7c trägt: Der Eingang in `reviewer.md` nennt Sicht, ADRs der Schnittstelle und
offene Slices nicht. 10b trägt: `AGENTS.md` §3.1 führt den alten Wortlaut.

Zitate: an jeder der zitierten Stellen den umgebenden Absatz gelesen, nicht nur die zitierte Zeile. Abweichungen allein
in V-5.

## Findings

### V-1 — MEDIUM — A1 „befolgt“: jeder Build bezieht ein nur per Tag benanntes Image

- **pfad:** `Dockerfile` Z. 1 und `tools/archive-wave/Dockerfile` Z. 1 · `# syntax=docker/dockerfile:1.7`
- **befund:** Das BuildKit-Frontend, das das Dockerfile liest, wird ohne Digest bezogen. Der Log aus M1 zeigt je Build
  `[auth] docker/dockerfile:pull token for registry-1.docker.io` und
  `resolve image config for docker-image://docker.io/docker/dockerfile:1.7`; welcher Digest läuft, entscheidet die
  Registry. `ADR-0011` §Entscheidung 1 sagt „Jedes extern bezogene Image wird per `@sha256:`-Digest gepinnt“ (die
  Aufzählung dahinter nennt `FROM`-Zeilen und Scanner, nicht das Frontend). `make versions` nennt es nicht. A1 führt als
  Ausnahmen allein A8 und A7, die Kurzfassung zu Frage A „zwei Lücken“. Daneben `make tidy` mit `golang:` nur per Tag
  (Runde 1, F-8, nicht eingearbeitet).
- **wirkung:** Der Auftraggeber liest die Container-Seite als festgelegt und entscheidet „nichts tun“. Von den Reviews
  ist die Stelle nicht benannt. Dieselbe Klasse wie F-7.
- **kleinste antwort:** A1 nennt beide Ausnahmen (oder führt „Lücke“, Preis: je ein Digest an den zwei Zeilen), die
  Kurzfassung zählt mit.
- **nicht geprüft:** ob sich der Tag `1.7` je auf andere Bytes bewegt und ein Ergebnis verändert hat.
- **verifizierbar:** ja

### V-2 — MEDIUM (Welle, nicht DoD) — die Vorlage steht außerhalb der erwarteten Menge von Trigger (b)

- **pfad:** Vorlage Z. 3, 47 und in Zeile 7b · `v6.18.1`; `wave-98` §3 (b); `MR-081`
- **befund:** Das Kommando aus (b) nennt am Stand zehn Dateien, darunter die Vorlage. Der Bump-Eintrag führt als
  erwartete Menge sich selbst, die `.d-check.yml`, fünf ADRs und drei Pläne der Welle, die mit dem Umzug nach `done/`
  herausfallen. Die Vorlage nennt er nicht, und sie liegt flach in `planning/`, wandert also nicht.
- **wirkung:** (b) verlangt „genau die Menge“; bei der Welle-Closure ist das so nicht erfüllt.
- **entscheidung:** beim Planner (Menge benannt fortschreiben, Nennung umformulieren oder die Vorlage anders ablegen).
- **verifizierbar:** ja — `git grep -l 'v6\.18\.1'` mit den Ausschlüssen aus (b)

### V-3 — LOW (Welle) — das Zählkommando von Trigger (f) ist mehrdeutig

Siehe §Inventar nachgezählt. **verifizierbar:** ja

### V-4 — LOW — DoD verlangt einen Anker-Link, 40 Zeilen verlinken nur die Datei

- **befund:** DoD „Die Fundstelle ist ein Anker-Link in den vendorten Baum“ (ebenso §3 §Spalten). Von 91 Zeilen haben
  48 keinen `#`-Anker: 8 zeigen auf die `.d-check.yml`-Vorlage (dort nicht möglich), **40** auf Markdown-Dateien, etwa
  8, 10a bis 10d, 14, 17a bis 17c, 20, 21, 23, 33, A15 bis A20.
- **wirkung:** Auffindbarkeit, nicht Entscheidung. Vor dem DoD-Haken nachziehen oder als Abweichung notieren.
- **verifizierbar:** ja

### V-5 — LOW — zwei Zeilen geben ihre Quelle ohne deren Bedingung wieder

- **28:** „mit einem Check bekommt die erste Drift-Frage eine maschinelle Vorstufe“. Der Absatz in `modul-12` setzt
  fort: „für die Werkzeuge im Image nur, sofern er im Image läuft“. In einem Docker-only-Repo ist das die tragende
  Bedingung. Von den Reviews nicht benannt.
- **5a:** `implement-slice.md` „schließt heute ein paralleles `in-progress/` aus“; die Quelle sagt „WIP-Limit = 1 **pro
  Implementer**“ (Runde 1, F-9, nach der Review-Vorgabe nicht eingearbeitet).
- **wirkung:** keine auf die Empfehlung; beide zählen für die Risiken 3 und 5.
- **verifizierbar:** ja

### V-6 — INFO — 4a bis 4d nennen die Reichweiten-Grenze des Kapitels nicht

`modul-04` §Bedrohung, Grenze, Ausgang: „Das Modell des Harness selbst — Prompt-Injection über Repo-Inhalte, die Rechte
des Agenten, der Container — … dieser Abschnitt behandelt es nicht.“ Wer „Bedrohungen und Vertrauensgrenzen“ freigibt,
erfährt aus der Vorlage nicht, dass der Agenten-Betrieb nicht darunter fällt.

### V-7 — INFO — 1d nennt eine der umformulierten Stellen in `modul-14` nicht

§Reproduzierbarkeits-Regeln bindet „Host-Lauf ist keine valide Gate-Ausführung“ neu an die fehlende Umgebungs-Zusage
und ergänzt einen Satz zum Tool-Pin mit Check. Sachlich decken 1b und A7 das ab.

## Risiken aus §6 — welchen Ausgang die Belege tragen (Entscheidung: Planner)

| Risiko | Belege tragen | Grund |
|---|---|---|
| 1 „bereits befolgt“ ohne Beleg | **eingetreten** (Evidenz-Datei) | F-7 (A5, behoben) und V-1 (A1, offen). Die übrigen Belege (2d, 2e, 9b, 10a, A2) tragen |
| 2 Empfehlung weitet über den Anlass | nach dem **Wortlaut** des Kriteriums **eingetreten**, in der Sache nicht beobachtet | „keine Zeile mehr als eine Neuerung“ ist nicht gehalten (22, 31, 32; 7b stellt eine zweite Frage, R2-4). Jede Empfehlung *übernehmen* nennt die kleinste Antwort mit Preis; eine Weitung habe ich nicht gefunden. Das Register-Auftreten (dann 2×) hängt an dieser Lesart |
| 3 Grenze ohne ihre Bedingung | **eingetreten** | F-2, F-6 (behoben), V-5 (offen). Das Kriterium „jede *nicht anwendbar*-Zeile zitiert die Bedingung“ ist zudem wörtlich nicht erfüllbar für die Zeilen „entschieden“ und „kein Akt“, die §3 Schritt 3 anders begründen lässt |
| 4 Sichtung schleust eine Übernahme ein | **entfallen** | siehe §Abgrenzung |
| 5 Wortlaut gestreckt | **eingetreten** | F-6 (behoben); in V-5 weicht die Wiedergabe vom Feld ab |

## Nicht geprüft

- Die Zeilen- und Datei-Schätzungen der Preise (nicht gebaut), darunter A5 und R2-1.
- Das Closure-Profil im Klon mit Arbeitsbaum (5b „braucht nichts“): von beiden Reviews gemessen, von mir nicht;
  `make gates` enthält es nicht.
- `make fullbuild` und die übrigen Aussagen des Welle-Triggers außer (b) und (f); CI.
- Die Lesbarkeit der Vorlage für den Auftraggeber — Frage des Validators.

## Nachprüfung zu V-1 — `582f10b3..60b4ea1f`

Geprüft sind allein `cd4fafbc` (Plan-Notiz) und `60b4ea1f` (Vorlage). Der Stand oben (Verdikt, DoD, Risiken) gilt für
`582f10b3`; was sich daran ändert, steht hier.

**Verdikt: V-1 ist geschlossen.** Keine neuen HIGH oder MEDIUM; zwei LOW, zwei INFO.

### A1 am Gegenstand

| Prüfpunkt | Urteil | Beleg |
|---|---|---|
| Klasse „Lücke“ | trägt | Das Kapitel verlangt die Umgebung gepinnt (Z. 15–16 des Absatzes zur Umgebungs-Zusage); zwei Bezüge sind es nicht |
| „Nur per Tag … zwei Bezüge“ | vollständig | `git grep -n '^# syntax='` 2 Zeilen; Dockerfiles, Makefiles, `a-check.mk`, `tools/` und Workflows nach weiteren Image-Bezügen durchsucht: Builder- und binfmt-Image in `release.yml` mit Digest, Scanner mit Digest, die Tag-Refs in `tools/image-scan.sh` sind Prüfgegenstand, nicht Umgebung. Docker und BuildKit des Hosts deckt A7 |
| Beleg „`make versions` nennt das Frontend nicht“ | trägt | 0 Treffer, auch nach dem Pin |
| Kleinste Antwort und Preis | trägt, gemessen (unten) | 3 Dateien, 3 Zeilen, kein neues Gate; „Diese Digests hielte kein `*-digest`-Target aktuell“ stimmt. Siehe N-2 |
| Empfehlung „übernehmen“ | folgt aus Klasse und Preis | Entscheidungs-Feld leer |
| A2, A5, Kurzfassung, Übersicht | stimmig | Frage A führt A1, A3, A5, A7, A8 als „Lücke“, davon im Container A1, A5, A8: „drei“ stimmt; die Übersichtszeile nennt A1 mit derselben Empfehlung |
| Zählung | unverändert | 40 Dateinummern, 68 Inventar-, 20 A-, 3 B-Zeilen; 46 `d-check:cite` (zwei neu), `make gates` grün |
| Plan-Notiz | vor der Korrektur committet | die genannten Stellen decken sich mit dem Diff; kein `Dockerfile`, kein `Makefile` geändert |

### Gemessen: Digest an den `syntax`-Zeilen

Wegwerf-Klon bei `60b4ea1f`, `# syntax=docker/dockerfile:1.7@sha256:a57df69d…` in beiden Dockerfiles. Das ist der Digest,
auf den der Tag `1.7` am Messtag zeigt (`docker buildx imagetools inspect`, ein Multi-Plattform-Index). Zusätzlich der
Go-Digest des Wurzel-`Dockerfile` an der Image-Zeile von `make tidy`.

| Lauf | Ergebnis |
|---|---|
| `make gates` (eigener Image-Name) | `EXIT=0`, Schlusszeile green; Coverage 94.90 %, semgrep 0, d-check 0. Der Log nennt den Digest-Bezug in 13 Zeilen und den Bezug nur per Tag in keiner |
| `make build` | `EXIT=0` |
| `make archive-wave-test` und `build` in `tools/archive-wave/` | `EXIT=0`, `ok  archive-wave` |
| Gegenprobe: Digest aus 64 Nullen, `make archive-wave-test` | `EXIT=2`, `docker.io/docker/dockerfile:1.7@sha256:0000…: not found` — der Digest wird ausgewertet, nicht überlesen |

**Ein Build läuft damit.** Nicht gemessen: ein Lauf mit leerem Layer-Cache (das Frontend wurde aufgelöst und las das
Dockerfile, die `RUN`-Schichten kamen teils aus dem Cache), der Multi-Plattform-Build des Release-Pfads
(`image-publish`, `image-test-arm64`) und `make tidy` selbst (lädt Module aus dem Netz).

### Neue Findings

**N-1 — LOW — `ADR-0011` ist ohne ihr Geltungs-Feld gelesen.** A1: „`ADR-0011` sagt es für jedes extern bezogene Image
…, ihre Aufzählung dahinter nennt das Frontend nicht.“ Der Satz aus §Entscheidung 1 ist wortgleich. Das Kopf-Feld
`Schärft:` der ADR sagt aber „verbindlich für jede `FROM`-Zeile im `Dockerfile`, das Scanner-Image in
`tools/semgrep.sh` und den Pin-Beleg in `make versions`“; Frontend und `make tidy` fallen nicht darunter
(`AGENTS.md` §5 Regel 16). Der Auftraggeber kann einen Verstoß gegen eine angenommene ADR lesen, wo die ADR nichts bindet.
Klasse und Empfehlung stehen auf dem Kapitel und bleiben. **Mein V-1 oben zitiert denselben Satz ohne dieses Feld** und
ist insoweit zu berichtigen.

**N-2 — LOW — der Preis nennt nicht, was der Digest an `make tidy` neu anlegt.** `golang:$(GO_VERSION)@sha256:…` ist eine
weitere Stelle des Go-Digests, und der Digest gewinnt über den Tag. Hebt eine Go-Hebung `GO_VERSION` und `Dockerfile`,
aber nicht diese Zeile, läuft `tidy` still mit dem alten Go: dieselbe Klasse wie der Fund in A8. Die Vorlage sagt, kein
`*-digest`-Target halte ihn aktuell, nicht aber, dass auch kein `pinsync`-Pin ihn gegen das `Dockerfile` hält.

**N-3 — INFO — Stütze für A1, an diesem Host.** Der lokale Tag `golang:1.27.4` zeigt auf `sha256:205c7636…`, das
`Dockerfile` und die Registry auf `sha256:5f1b98fa…`. `docker run` zieht bei vorhandenem Tag nicht nach: `make tidy`
liefe hier heute mit einem anderen Image als der Build.

**N-4 — INFO — zwei Nebenstellen.** Zeile 1a folgt weiter „A3, A7“, obwohl A1 jetzt eine eigene Empfehlung trägt (A1
steht in der Übersicht, der Verlust ist gering). A1 sagt, das Kapitel nenne beim Digest das Base-Image; dasselbe Kapitel
verlangt für die Rezept-Form „jede Eingabe ist digest-gepinnt“, was die Zeile stärker stützte.

### Stand der übrigen Punkte

V-2 bis V-7 sind von den zwei Commits nicht berührt. Zu den Risiken: 1 bleibt **eingetreten** (beide Auftreten, F-7 und
V-1, sind behoben); 5 bleibt **eingetreten**, mit N-1 als weiterem Auftreten. DoD „Klasse und Empfehlung“: die
Einschränkung zu A1 entfällt, V-4 und V-5 bleiben.
