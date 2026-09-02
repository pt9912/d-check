# Review-Report — slice-187 (Handbuch-Zeiger nennt beide URL-Formen)

- **Review-Art:** Code (gegen Plan, `DC-FA-CLI-001`/`.a`, `DC-FA-CLI-010`/`.a`, Hard Rules)
- **Gegenstand:** `e25f8b0..HEAD` — Basis-Commit `e25f8b0` (Beanspruchung `open/`→`in-progress/`, Vorprüfungen aufgefrischt), geprüfter Commit `1ae05ca` (Implementierung). HEAD `1ae05ca`.
- **Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) v1.13.0
- **Modell-ID:** claude-sonnet-5
- **Datum:** 2026-09-02
- **Eingangs-Kontext:** Slice-Plan [`slice-187`](../plan/planning/in-progress/slice-187-handbuch-url-beide-formen.md); [`DC-FA-CLI-001`](../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)/[`DC-FA-CLI-001.a`](../../spec/spezifikation.md#dc-fa-cli-001a--ablauf-eines-prüflaufs), [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)/[`DC-FA-CLI-010.a`](../../spec/spezifikation.md#dc-fa-cli-010a--makefile-fragment); [`AGENTS.md`](../../AGENTS.md) §3.3 (Beanspruchungs-Ausnahme), §3.4, §3.7, §5; Baseline-Regelwerk [`modul-06-roadmap.md`](../../.harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md) §Das Beobachtungs-Register, [`modul-05-planning-harness.md`](../../.harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md) §Lifecycle als State Machine; Beobachtungs-Register [`BEO-002`](../plan/planning/observations.md), [`BEO-012`](../plan/planning/observations.md); vorherige Findings am gleichen Modul: [`2026-08-31-slice-185-code-r1.md`](2026-08-31-slice-185-code-r1.md), [`2026-08-31-slice-186-vcs-rename-code-r1.md`](2026-08-31-slice-186-vcs-rename-code-r1.md), [`2026-08-30-slice-181-handbuch-link-review.md`](2026-08-30-slice-181-handbuch-link-review.md).

Repo unverändert (`git status --short` leer, HEAD `1ae05ca`).

---

## Eigener Lauf (Ausgabe, nicht behauptet)

| Lauf | Ausgabe |
|---|---|
| `make test` | Docker-Build + `go test ./...` grün, alle Pakete `ok` |
| `make gates` | `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; semgrep `Ran 55 rules on 56 files: 0 findings`; `d-check: 650 Datei(en) geprüft, 0 Befund(e)` (×2, `targets`/`planning`-Läufe); Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` — deckt sich mit der Commit-Botschaft (zehn Gates, 650 Dateien, 0 Befunde) |
| `curl -sI https://raw.githubusercontent.com/pt9912/d-check/refs/heads/main/docs/user/benutzerhandbuch.md` | `HTTP/2 200`, `content-type: text/plain`, identer `etag` wie die Kurzform ohne `refs/heads/` — die neue Konstante ist eine echte, auflösende GitHub-Raw-URL |
| `git show e25f8b0^:.../observations.md` vs. `e25f8b0` vs. `HEAD` (BEO-002/BEO-012-Zeilen) | BEO-012 unverändert bei 12 (Sichtung, keine neue Instanz); BEO-002 springt in **`e25f8b0`** (Beanspruchungs-Commit, nicht im geprüften Diff) von 7 auf 8 samt vollständiger „Achte Instanz"-Erzählung im Perfekt |
| `git show e25f8b0:internal/adapter/driving/cli/cli.go \| grep handbuchURLRaw` | kein Treffer — die Konstante existiert zum Zeitpunkt des Beanspruchungs-Commits noch nicht |
| Vergleich mit Geschwister-Beanspruchungs-Commits `39258a4` (slice-185), `8117758` (slice-186), `dde0916` (slice-183), `fe00f17` (slice-189) | Alle vier lassen `docs/plan/planning/observations.md` **unverändert** und alle DoD-Häkchen **`[ ]`** — der Vorprüfungs-Schritt liest den Registerstand nur in den Slice-Plan hinein, schreibt aber nicht ins Register selbst. `e25f8b0` weicht davon ab |

---

## Findings

### F-1 · HIGH · Baseline-Regelwerk `modul-06-roadmap.md` §Das Beobachtungs-Register („Eingetragen wird bei der Slice-Closure") / `AGENTS.md` §3.3 Beanspruchungs-Ausnahme · `docs/plan/planning/observations.md` (BEO-002-Zeile) + `docs/plan/planning/in-progress/slice-187-handbuch-url-beide-formen.md:76-95`, beide geschrieben in Commit `e25f8b0`

**Befund:** Der Beanspruchungs-Commit `e25f8b0` (vor dem hier geprüften Diff, aber Gegenstand von Frage 5 der Aufgabe) schreibt bereits den vollständigen Registereintrag „Achte Instanz, slice-187" in `observations.md` — Zähler 7→8, im Perfekt erzählt („die Prüfung ergab: nein", „Geprüft und **tatsächlich** um die zweite Form ergänzt: die Test-Literale und die `--print-mk`-Illustration im Handbuch — beide reproduzieren den Bau wörtlich") — und hakt im selben Commit die DoD-Punkte 1–3 der Slice-Datei bereits ab (`[x]`), inklusive der Detail-Behauptung „Umgesetzt: `handbuchURLRaw`-Konstante neben `handbuchURL` in `cli.go`". Zu diesem Zeitpunkt existiert `handbuchURLRaw` im Repo nicht (`git show e25f8b0:internal/adapter/driving/cli/cli.go` enthält die Konstante nicht) — der Code, den die Erzählung als „tatsächlich ergänzt" beschreibt, entsteht erst zwei Minuten später im nächsten Commit (`1ae05ca`).

Das widerspricht doppelt der eigenen kanonischen Form: (a) Modul 6 legt fest *„Eingetragen wird bei der **Slice-Closure**"* — nicht bei der Beanspruchung; die vier Geschwister-Beanspruchungs-Commits dieser Session (`39258a4`/slice-185, `8117758`/slice-186, `dde0916`/slice-183, `fe00f17`/slice-189) belegen die korrekte Form: alle lassen `observations.md` bei der Beanspruchung unverändert (der Sichtungs-Schritt liest den Registerstand nur in den Slice-Plan, schreibt nicht ins Register), und die beiden mit DoD-Abschnitt (`slice-185`, `slice-186`) lassen alle Häkchen `[ ]`. (b) Ein DoD-Haken ist laut `AGENTS.md` §5 „eine Selbstauskunft" — ein Haken, der ein Ereignis als geschehen behauptet, bevor der Commit existiert, der es umsetzt, macht den Zustand der Beanspruchungs-Commit-Version des Repos nachweislich falsch: Wer `e25f8b0` auscheckt, liest eine Zusage über Code, der dort nicht liegt.

Der Schaden ist auf den ersten Blick harmlos, weil `1ae05ca` die Behauptung zwei Minuten später einlöst — aber die Reihenfolge selbst ist der Befund: Der Beanspruchungs-Commit sollte laut Modul 5 „vor der Arbeit" landen, trägt hier aber bereits deren Ergebnis. Sollte die Closure von slice-187 `docs/plan/planning/observations.md` bei Schritt 3 (Register fortschreiben) unverändert lassen, weil der Eintrag ja „schon da" ist, hat der eigentliche Closure-Schreibschritt nie stattgefunden — er wurde durch die Beanspruchung vorweggenommen.

**Verifizierbar:** ja — `git show e25f8b0^:docs/plan/planning/observations.md` vs. `git show e25f8b0:...` zeigt den Sprung 7→8 im Beanspruchungs-Commit; `git show e25f8b0:internal/adapter/driving/cli/cli.go` zeigt die fehlende Konstante; der Vergleich mit den vier genannten Geschwister-Commits ist mit `git show <hash> -- docs/plan/planning/observations.md` reproduzierbar (leerer Diff bei allen vieren).
**Klasse:** `register-und-dod-vorab-abgehakt`

### F-2 · MEDIUM · `BEO-012`-Klasse (Quelle über ihren Geltungsbereich hinaus zitiert) · Commit-Botschaft `1ae05ca`, Absatz 3

**Befund:** Die Commit-Botschaft schreibt: *„Diese Unterscheidung (Illustration vs. Verweis) ist die eigentliche Antwort auf die in **MR-021**/BEO-002 offene Kopplungsfrage."* [`MR-021`](../../harness/conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden) hat einen expliziten `Geltungsbereich`: *„alle Markdown-Links auf `.harness/baseline/<tag>/…` in der Live-Doku"* — die Pin-Bindung von in-repo-Verweisen auf das vendorte Regelwerk. Das hat mit der hier verhandelten Frage (ob eine zweite, unabhängige Zeichenkette in `packaging/dockerhub/overview.md` mitgezogen werden muss) nichts zu tun; keine der zwölf bisherigen `MR-021`-Nennungen im Repo (`grep -rln MR-021 docs/plan/planning/done/`) behandelt String-Duplikate. `BEO-002` selbst trägt die Frage korrekt und ist als **verkörpert in [`MR-025`](../../harness/conventions.md#mr-025)** ausgewiesen — vermutlich war `MR-025` gemeint, nicht `MR-021`. Der Fehlschluss selbst ist damit ein Beispiel der Klasse, die er in der Commit-Botschaft neben sich zitiert.
**Verifizierbar:** ja — `harness/conventions/MR-021-vendored-verweise-pin-gebunden.md` `Geltungsbereich`-Feld liegt lesbar vor.
**Klasse:** `quelle-ueber-geltungsbereich-zitiert`

### I-1 · INFO · `AGENTS.md` §3.7 (Herkunfts-Feld-Enumeration enger als Bestandspraxis) · `internal/adapter/driving/cli/cli.go:118,137`

**Befund:** Der neue Kommentar trägt `seit slice-187` und `(BEO-002)` als Herkunfts-Anker. `AGENTS.md` §3.7 zählt für Code/Konfiguration/Skripte explizit nur `DC-*`, `ADR-*`, `MR-*`, `seit welle-<NN>` auf — weder Slice-Nummern noch `BEO-*` stehen in dieser Liste, und derselbe Absatz nennt „keine Slice-Nummern" ausdrücklich als verbotenen Inhalt. In der Praxis ist beides jedoch seit Langem Bestand (`seit slice-073`/`-057`/`-059`/`-099` in `markdown.go`/`configyaml.go`/`trace.go`; `(BEO-003)`/`(BEO-023)` u. a. in `structure.go` und mehreren Testdateien) — slice-187 setzt eine etablierte, wiederholt reviewte Konvention fort, keine neue. Kein Finding gegen diesen Diff; die Diskrepanz zwischen Wortlaut und Bestand ist selbst meldenswert, aber unabhängig von slice-187 entstanden.
**Verifizierbar:** ja — `grep -rn "seit slice-\|BEO-[0-9]" internal/` zeigt den Bestand vor diesem Slice.
**Klasse:** `herkunfts-anker-enumeration-enger-als-bestand`

---

## Negativbefunde (geprüft, ohne Befund)

- **Prüffrage 1 (Silent-Grün).** `make gates` unabhängig gefahren: zehn Gates grün, 650 Dateien/0 Befunde — deckt sich mit der Commit-Botschaft, kein stiller Pfad gefunden.
- **Prüffrage 2 (Kern-Modul-Korrektheit).** Keine Kern-Module (`internal/hexagon/core/rules`) sind in diesem Diff berührt — reine CLI-Adapter-/Spec-/Doku-Änderung.
- **Prüffrage 3 (Hexagon-Import-Richtung).** `arch-check` läuft grün als Teil von `make gates`; keine neue Import-Kante (nur ein neues String-Konstantenpaar im bestehenden `cli`-Paket).
- **Prüffrage 5 (Netzzugriff außerhalb `external`).** Kein neuer Netzzugriff im Produktcode; die eigene `curl`-Probe gegen die neue URL war Teil der Review-Verifikation, nicht des Produkts.
- **Punkt 1 der Aufgabe — URL-Korrektheit.** `https://raw.githubusercontent.com/pt9912/d-check/refs/heads/main/docs/user/benutzerhandbuch.md` löst real auf (`HTTP/2 200`, identer `etag` zur Kurzform ohne `refs/heads/`) — beide Schreibweisen sind gültige, von GitHub unterstützte Raw-URL-Formen.
- **Punkt 2 der Aufgabe — Testabdeckung `TestHandbuchURL_TraegtKeineVersion`.** Die Fatalf-Prüfung `len(zeilen) != 2` fängt jede fehlende Zeile bereits vor der Fallunterscheidung ab; die Switch-Fallunterscheidung plus die abschließende `!sahBlob || !sahRaw`-Prüfung fängt auch den Fall zweier gleichartiger Zeilen (z. B. zwei `raw`-Zeilen statt einer von jeder Form) ab. Keine Lücke gefunden, die einen Test grün ließe, während eine Form fehlt oder dupliziert ist.
- **Punkt 3 der Aufgabe — Entscheidung gegen `overview.md`.** `packaging/dockerhub/overview.md` wird laut `.github/workflows/hub-description.yml` ausschließlich per `sed s/__VERSION__/…/` befüllt und dann roh an `dockerhub-description`-Action übergeben — keine Generierung aus `--print-mk`/`--help`-Output. Die Begründung „handkuratierter Link, keine Output-Illustration" trägt.
- **Punkt 4 der Aufgabe — Spec-Konsistenz.** `spec/lastenheft.md` (0.82.0) und `spec/spezifikation.md` (Historie 2026-09-01) beschreiben Reihenfolge (blob zuerst, raw danach), Beschriftung und Versionslosigkeit deckungsgleich mit dem Code (`writeUsage`, `mkTemplate`).
- **§3.4 (Architektur sprachfrei, Spec referenziert nie abwärts).** Weder `spec/lastenheft.md` noch `spec/spezifikation.md` nennen in den neuen Zeilen Slice-/ADR-/Commit-Bezüge; `spec/architecture.md` ist in diesem Diff nicht berührt.
- **`CHANGELOG.md`.** Nicht im Diff — konform zur Release-Prep-Regel (`AGENTS.md` §5).
- **Fünfter Spiegel / vollständige Spiegel-Menge.** `grep -rln` über den Baum nach der `blob`/`raw`-URL findet neben den vier bearbeiteten Stellen nur historische Review-Reports (`docs/reviews/2026-08-3*`) — korrekt unverändert, da Lauf-Belege, keine lebenden Spiegel.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 1 | F-1 |
| MEDIUM | 1 | F-2 |
| LOW | 0 | — |
| INFO | 1 | I-1 |

**Finding-Klassen dieses Laufs:** `register-und-dod-vorab-abgehakt` · `quelle-ueber-geltungsbereich-zitiert` · `herkunfts-anker-enumeration-enger-als-bestand`

## Verdikt

**Merge-blockierend: ja (F-1).**

Der eigentliche Code-Diff (`1ae05ca`) ist sauber: die neue Raw-URL löst real auf, die Testabdeckung hat keine Lücke, die Entscheidung gegen eine zweite Form in `packaging/dockerhub/overview.md` ist geprüft statt angenommen und trägt, Lastenheft und Spezifikation ziehen präzise nach, `make gates`/`make test` sind unabhängig grün nachgefahren. Kein Hexagon-, Suppressions- oder Netz-Verstoß; §3.4 unberührt.

**F-1 ist der schwerwiegende Befund**, liegt aber technisch im Beanspruchungs-Commit `e25f8b0`, nicht im als „Implementierung" bezeichneten `1ae05ca` — die Aufgabe verlangt jedoch ausdrücklich die Prüfung genau dieses Inhalts (Frage 5). Drei Geschwister-Beanspruchungs-Commits derselben Session zeigen, dass die korrekte Form bekannt und geübt ist (`observations.md` bleibt bei der Beanspruchung unangetastet); `e25f8b0` weicht davon ab, indem es einen vollständigen, im Perfekt erzählten Registereintrag samt abgehakter DoD-Punkte schreibt, bevor der beschriebene Code existiert. Empfehlung an die Implementer-Rolle: vor der Closure den BEO-002-Eintrag und die DoD-Reihenfolge gegen den tatsächlichen Commit-Verlauf prüfen und im Closure-Schritt (nicht rückwirkend in `e25f8b0`) bestätigen, dass der Schreib-Zeitpunkt korrekt nachgezogen wird — eine Umschreibung von `e25f8b0` selbst verbietet sich, weil er bereits `HEAD`-Vorgänger ist (keine `--amend`-Politik). F-2 ist vor Merge leicht behebbar (Korrektur der Commit-Botschaft ist nicht mehr möglich, aber die Closure-Notiz sollte die korrekte Kennung `MR-025` statt `MR-021` führen, falls die Aussage dort wiederholt wird).
