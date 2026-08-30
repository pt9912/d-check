# Review slice-181 — Das Handbuch ist aus dem Werkzeug heraus auffindbar


**Review-Art:** Code/Design (gegen Slice-Plan, Lastenheft, Hard Rules)
**Gegenstand:** `84f7259~1..50d461e` (5 Commits); `ff21044` (eingehender CR) ausdrücklich **nicht** Gegenstand
**Skill:** `.harness/skills/reviewer.md` @ v1.13.0 · **Modell:** claude-opus-5[1m] · **Datum:** 2026-08-30
**Eingangs-Kontext:** `AGENTS.md` §2/§3.1/3.4/3.6/3.7/3.8/§4/§5, `harness/conventions.md` (MR-013, MR-053, MR-054), `.harness/baseline/v5.12.0/regelwerk/` (`grundlagen-harness-dateien` §Was ein Kommentar trägt, `modul-03-spec`, `modul-04-adrs`, `modul-05`, `modul-06`), `DC-FA-CLI-001`, `DC-FA-CLI-010`, `DC-FA-VER-001`, Beobachtungs-Register (`BEO-009`, `BEO-012`, `BEO-020`, `BEO-023`), `docs/user/releasing.md`

### Eigener Lauf

| Lauf | Ausgabe |
|---|---|
| `make doc-check` | `613 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `make test` | alle Pakete `ok` · Exit 0 |
| `make lint` | `0 issues.` · Exit 0 |
| `make arch-check` | `gesamt: 0 Befund(e)` · Exit 0 |
| `make planning-check` | `613 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `make trace-check RANGE=84f7259~1..50d461e` | `613 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `make adr-check RANGE=84f7259~1..50d461e` | `613 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `docker run d-check:latest --help` / `--print-mk` | je genau **eine** URL-Zeile; help auf stderr (stdout 0 Byte), Exit 0 |
| 6 Mutations-/Gegenproben in Kopien unter `scratchpad/rev181` | s. u. |

`make gates` und `make record-gates` bewusst **nicht** gefahren (Nachweis-Datei).

---

## Urteil

**BLOCKIERT** — 1 HIGH · 5 MEDIUM · 2 LOW · 1 INFO.

Der Bau selbst funktioniert und ist gegen das echte Binary belegt; die drei behaupteten Mutationsproben habe ich reproduziert und sie stimmen exakt. Blockierend ist, **wo die Begründung gelandet ist und was sie behauptet**: ein Konjunktiv über die verworfene Alternative in zwei neuen Kommentaren, eine über ihren Geltungsbereich hinaus zitierte Anforderung an fünf Stellen (darunter das Lastenheft), und die Spezifikation, die die beiden geänderten Ausgaben weiterhin abschließend ohne den neuen Zeiger aufzählt.

---

## Findings

### H-1 · Konjunktiv über die verworfene Alternative in zwei neuen Kommentaren

- **quelle:** `AGENTS.md` §3.7, Baseline `grundlagen-harness-dateien.md` §Was ein Kommentar trägt (Zeitform-Test)
- **pfad:** `internal/adapter/driving/cli/cli.go:122-124` · `internal/adapter/driving/cli/cli_handbuch_link_test.go:54-58`
- **zugesagt:** Kanon, wörtlich: *„Unzulässig ist er [der Konjunktiv] über die **verworfene Alternative**: Die ist entschieden, und die Entscheidung steht in der ADR. Die Probe dafür ist die Zeitrichtung — zeigt der Konjunktiv nach vorn oder zurück?"* Und: *„Ein Kommentar schreibt an jemanden, der den Code gleich ändert — nicht an jemanden, der die Entscheidung noch einmal treffen will. Der zweite Leser hat die ADR."*
- **gemessen:** beide Kommentare tragen denselben rückwärts zeigenden Konjunktiv — `cli.go:123`: *„eine versionierte URL **waere** eine Release-Prep-Flaeche, die kein Gate deckt"*; `cli_handbuch_link_test.go:54-55`: *„Eine versionierte URL **waere** eine Release-Prep-Flaeche…"*. Beide begründen die **verworfene** Option, nicht die vorhandene Zeile. Der jeweilige Rest (*„Der Preis ist benannt: wer ein altes Image faehrt, liest ein neueres Handbuch"* / *„er faellt, sobald jemand `blob/v0.68.0` schreibt"*) ist eine saubere Grenze bzw. Zusage — der Befund gilt nur der Konjunktiv-Klausel dazwischen.
- **warum es zählt:** die Deliberation hat in diesem Slice einen legitimen Ort — §2 des Slice-Plans — und ist von dort in den Produktionscode und den Test kopiert worden, wo sie mit jeder Zeile weiterreist und bei jeder künftigen Änderung erneut gegen den Zeitform-Test läuft. Der Skill führt „Deliberation über Verworfenes im Kommentar" ausdrücklich als HIGH.
- **verifizierbar:** nein — kein Gate; `make lint` meldet `0 issues` (gemessen).
- **klasse:** `kommentar-konjunktiv-rueckwaerts`
- **billigster Fix:** die Konjunktiv-Klausel in beiden Kommentaren streichen; die Grenze („der Preis") und die Zusage stehen lassen.

### M-1 · `DC-FA-VER-001` wird über seinen Geltungsbereich hinaus zitiert — an fünf Stellen

- **quelle:** `BEO-012`, `AGENTS.md` §5 (*„Eine zitierte Quelle trägt nur, was in ihrem Geltungsbereich steht … die **direkteste** Quelle wählen"*)
- **pfad:** `spec/lastenheft.md:98` · `spec/lastenheft.md:3343` · `internal/adapter/driving/cli/cli.go:124` · `internal/adapter/driving/cli/cli_handbuch_link_test.go:55-56` · Commit-Botschaft `ca43027`
- **zugesagt:** *„[`DC-FA-VER-001`] hält ausschließlich `ghcr`-präfixierte Pins gegen `version.md`, nackte Tags in Prosa driften still."*
- **gemessen:** `DC-FA-VER-001` sagt das nicht. Der Anforderungstext (`spec/lastenheft.md:1728-1731`) führt das ghcr-Muster als **Beispiel** (*„z. B. `ghcr.io/…/d-check:v\\d+\\.\\d+\\.\\d+`"*) eines frei konfigurierbaren `versions.pin-pattern`; §Mehrere Muster-Quellen-Paare (`:1759-1766`) nennt als Anwendungsfall ausdrücklich *„ein fremder Baseline-Tag gegen den Pin im Konventionsspeicher"* — also genau Tags in URLs und Prosa. Das eigene Beobachtungs-Register sagt es noch deutlicher: `BEO-008` Stand — *„Die benannte mechanische Form ist seit slice-122 baubar: `versions.patterns` trägt mehrere Muster-Quellen-Paare, ein zweiter Abgleich (Baseline-Tag in URLs und Prosa gegen den §Baseline-Pin) ist damit konfigurierbar. Das eigene Profil fährt ihn **nicht**."* Wahr ist die Aussage über **dieses Repo-Profil**: `.d-check.yml:397-402` führt die Kurzform mit genau einem `pin-pattern`. Die direkteste Quelle dafür ist `docs/user/releasing.md:35-37` — *„Der `versions`-Gate prüft nur `ghcr`-**präfixierte** Pins gegen `version.md#aktuell`, nicht Fließtext oder nackte Tags"* —, und sie spricht korrekt vom **Gate**, nicht von der Anforderung. Der Slice-Plan §2 formuliert es ebenfalls korrekt („der `versions`-Gate"); die Überdehnung entsteht erst beim Übertragen ins Lastenheft, in den Code und in den Test.
- **Nebenbefund derselben Messung:** für den *tatsächlichen* Fundort der URL — eine Go-Konstante — trägt das Argument ohnehin nicht. `d-check` scannt ausschließlich `*.md` (`spec/lastenheft.md:924`); **keine** `versions`-Konfiguration könnte `cli.go` je erreichen. Der belastbare Grund ist die Scan-Menge (§3.8), nicht der Musterzuschnitt.
- **warum es zählt:** die Aussage steht jetzt im **ranghöchsten** Stratum als Eigenschaft einer *anderen* Anforderung. Wer später einen `versions.patterns`-Zweitabgleich erwägt (den `BEO-008` als baubar führt), liest im Lastenheft, das gehe nicht.
- **verifizierbar:** nein — kein Gate; Beleg ist der Anforderungstext selbst.
- **klasse:** `zitat-ueberdehnt-anforderung-statt-profil`
- **billigster Fix:** an allen fünf Stellen von der Anforderung auf das Profil umstellen („dieses Repo konfiguriert `versions` nur auf `ghcr`-präfixierte Pins").

### M-2 · Die Spezifikation zählt beide geänderten Ausgaben weiterhin abschließend ohne den Zeiger auf

- **quelle:** `AGENTS.md` §2 (Source Precedence: *„die niedriger rangierte wird angepasst"*), §6 Schritt 7
- **pfad:** `spec/spezifikation.md:48-56` · `spec/spezifikation.md:683-684`
- **zugesagt:** Lastenheft 0.77.0 (Rang 1) fordert den Handbuch-Zeiger in beiden Ausgaben.
- **gemessen:** die Spezifikation (Rang 2) wurde im ganzen Diff nicht angefasst und widerspricht jetzt. `:49-56`: *„Die Usage-Ausgabe führt (**in dieser Reihenfolge**) eine Kurzbeschreibung, die Synopsis …, eine Zeile zum Pfad-Argument …, die Flag-Liste … **und** einen Konfigurations-Hinweis, der auf `DC-FA-CLI-005` … und `DC-FA-CLI-006` … verweist"* — eine geschlossene, geordnete Fünfer-Liste; die reale Ausgabe hat sechs Elemente. `:683-684`: *„1. Kommentar-Kopf: Einbindung via `include`, Hinweis zum Digest-Pin via `DCHECK_IMAGE` bzw. bequemer `DCHECK_DIGEST`."* — die Handbuch-Zeile fehlt. Die Historie-Tabelle der Spezifikation (`:3104 ff.`) trägt keine Zeile für 2026-08-30.
- **warum es zählt:** die Spezifikation ist die Stelle, an der ein Implementer die Ausgabe-Reihenfolge nachschlägt; wer sie beim nächsten Umbau befolgt, entfernt den Zeiger regelkonform wieder — genau das Versagen, gegen das §2 des Slice-Plans die beiden Akzeptanzkriterien geschrieben hat. `**Berührte Spec-Stellen:**` im Slice-Kopf nennt die Spezifikation nicht, §3 *Ausdrücklich NICHT* auch nicht — die Auslassung ist weder gedeckt noch deklariert.
- **verifizierbar:** nein — `doc-check` prüft Links und Anker, nicht die Deckung zweier Prosa-Aufzählungen (`613/0` gemessen).
- **klasse:** `stratum-nachzug-fehlt`
- **billigster Fix:** beide Aufzählungen um den Zeiger ergänzen, Historie-Zeile in der Spezifikation.

### M-3 · Der `writeUsage`-Doc-Kommentar hängt jetzt an `const handbuchURL`

- **quelle:** `AGENTS.md` §3.7 (*„Ein Kommentar beschreibt, was da ist"*), Rang-Zeiger-Klasse
- **pfad:** `internal/adapter/driving/cli/cli.go:116-127`
- **zugesagt:** der Block beginnt mit *„writeUsage gibt die Hilfe aus (`DC-FA-CLI-001.a`): Kurzbeschreibung, Synopsis mit dem Pfad-Argument, Flag-Liste …"*
- **gemessen:** die Konstante ist **zwischen** Kommentar und Funktion eingefügt (`:126`), ohne Leerzeile davor. Nach Go-Doc-Regeln ist die zusammenhängende Kommentargruppe damit die Doku von `handbuchURL`; `writeUsage` (`:127`) hat keinen Doc-Kommentar mehr. Der Rang-Zeiger `DC-FA-CLI-001.a` — das einzige auflösbare Herkunfts-Feld der Funktion — sitzt jetzt an der Konstante.
- **warum es zählt:** wer `handbuchURL` liest, bekommt eine Beschreibung der Hilfe-Ausgabe; wer `writeUsage` ändert, findet weder Zusage noch Rang-Zeiger. `revive` ist im Profil aktiv, greift aber nur bei exportierten Deklarationen — `make lint` meldet `0 issues` (gemessen), das Muster ist unbewacht.
- **verifizierbar:** nein.
- **klasse:** `doc-kommentar-am-falschen-element`
- **billigster Fix:** Konstante mit eigenem Kommentar vor den `writeUsage`-Block ziehen (oder eine Leerzeile setzen und den Funktionskommentar wieder direkt an `writeUsage`).

### M-4 · Die „im Kopf"-Zusage von `TestCLI010` hat einen stillen Ausweich-Pfad — und keine Probe

- **quelle:** `BEO-023` (Zähler **4**, Schwelle überschritten), `AGENTS.md` §5 (Proben-Disziplin)
- **pfad:** `internal/adapter/driving/cli/cli_handbuch_link_test.go:41-51`
- **zugesagt:** Kommentar `:41-44`: *„Im Kopf, nicht irgendwo: die Zeile steht VOR der ersten Variablen-ZUWEISUNG, sonst stuende sie zwischen Targets statt in der Einleitung."* Slice-DoD: *„jeder wird von genau der Mutation rot, gegen die er steht."*
- **gemessen:** `:45-48` — `kopf := stdout; if i := strings.Index(stdout, "\
DCHECK_IMAGE ?="); i > 0 { kopf = stdout[:i] }`. Fehlt der Anker, bleibt `kopf` die **ganze** Ausgabe und die Zusage kollabiert kommentarlos auf die schwächere Prüfung eine Zeile darüber. **Probe C** (Kopie, `DCHECK_IMAGE` → `DCHECK_IMG`, Handbuch-Zeile als **allerletzte** Zeile hinter allen zwölf Targets, fokussierter Lauf `go test -v -run Handbuch`): `--- PASS: TestCLI010_PrintMKNenntHandbuch`, Exit 0. Die drei Mutationen des Autors erreichen die Zeilen 45-51 nie: sie entfernen die URL ganz, also feuert schon `t.Fatalf` auf `:39`. Erst meine **Probe B** (Zeile unter `DCHECK_IMAGE ?=` verschoben, Anker intakt) bringt die Kopf-Assertion zum Sprechen — `cli_handbuch_link_test.go:50: Handbuch-URL steht nicht im Kopf des Fragments`. Die Assertion beißt heute also; belegt war das nicht.
- **Abschwächung, gemessen:** `cli_acceptance_test.go:2539` prüft `"DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v"` — eine schlichte Umbenennung fiele dort auf. Erreichbar bleibt der Pfad, wenn `DCHECK_IMAGE ?=` je die erste Fragment-Zeile würde.
- **warum es zählt:** genau die Klasse, die `BEO-023` bei Zähler 4 führt — *„Ein Wächter, der nie fangen konnte, liest sich wie einer, der fängt"*; hier: ein Wächter, dessen benannte Hälfte still schwächer werden kann, ohne dass der Kommentar diese Grenze nennt.
- **verifizierbar:** ja — der obige fokussierte Lauf.
- **klasse:** `assertion-faellt-still-auf-schwaechere-form-zurueck`
- **billigster Fix:** den Ausweich-Pfad in ein `t.Fatalf` drehen oder die Grenze im Kommentar benennen; die fehlende Mutation (Zeile unter den Anker) nachziehen.

### M-5 · Die Lösungs-Begründung steht im Akzeptanzkriterium des Lastenhefts

- **quelle:** Baseline `modul-04-adrs.md` §Regeln gegen typische Fehlannahmen — *„ADRs begründen die **Lösung**. Anforderungen begründet die Spec. Wer ADRs zur Spec macht, kann später keine Architektur ohne Lastenheft-Änderung wechseln."* · `AGENTS.md` §3.4 (*„die Begründungen leben in den ADRs"*) · `modul-03-spec.md` §Ziel-Form: Akzeptanzkriterium (*„Das Lastenheft sagt, **was** zugesagt ist"*)
- **pfad:** `spec/lastenheft.md:98` (Kriterium *Hilfe*), Folgeverweis in `:516`
- **zugesagt:** Slice §3: *„die eine nicht-offensichtliche Wahl — `main` statt Version — steht mit ihrer Begründung im Lastenheft, wo ein Leser sie sucht."*
- **gemessen:** das Kriterium trägt nach dem prüfbaren `then` zwei reine Begründungssätze (*„eine versionierte Form wäre eine Release-Prep-Fläche, die kein Gate deckt …"*, *„Der Preis ist benannt: …"*). Es ist das **einzige** von **74** Given/When/Then-Kriterien des Lastenhefts, das eine Begründung dieser Art führt (gemessen über `grep` auf `weil |Der Preis|Begründung` in den Kriterien-Zeilen: 1 Treffer, und das ist `:98`). Die Zusage selbst („die URL zeigt auf den Hauptzweig") ist dort richtig aufgehoben und getestet — der Befund gilt nur den beiden Begründungssätzen.
- **warum es zählt:** genau der Effekt, den `modul-04` benennt: eine spätere Umstellung auf eine versionierte Form (oder auf ein `versions.patterns`-Paar) verlangt jetzt eine Änderung am **vertraglich abnahmebindenden** Stratum statt einer Folge-ADR. Und sie transportiert M-1 an die rangoberste Stelle.
- **Nicht** Gegenstand dieses Befunds: dass keine ADR geschrieben wurde. §3.6 ist nicht ausgelöst (nichts gelockert), und die ADR-Pflicht ist hier nicht automatisch — 17 von 58 jüngeren Lastenheft-Historie-Zeilen tragen *„Begründung in begleitender ADR"*, 41 nicht (gemessen). Der Slice-Plan §2 ist ein legitimer Ort für die Begründung; das Lastenheft ist der zusätzliche, falsche.
- **verifizierbar:** nein.
- **klasse:** `loesungsbegruendung-im-lastenheft`
- **billigster Fix:** die beiden Begründungssätze aus `:98` streichen (die Historie-Zeile `:3343` darf sie tragen), die prüfbare Zusage stehen lassen.

### L-1 · „elf `##`-annotierte Targets" — es sind zwölf

- **quelle:** Doku-Drift (Maintainability); `docs/user/releasing.md:88-91` (Enumerationen ohne Gate)
- **pfad:** `docs/user/benutzerhandbuch.md:836-839` · gleichlautend `spec/spezifikation.md:692`
- **gemessen:** `docker run d-check:latest --print-mk | grep -c '^[a-z-]*: ##'` = **12**; es fehlt `doc-structure` in beiden Aufzählungen. `DC-FA-CLI-010` Happy Path (`spec/lastenheft.md:516`) führt `doc-structure` korrekt — Rang 1 und Rang 2/6 widersprechen sich.
- **Bestandslage:** vorbestehend seit `e93d6a9` (2026-08-15, slice-099), **nicht** von diesem Slice eingeführt. Gemeldet, weil der Diff genau diesen Satz im Handbuch angefasst hat (`50d461e` hängt hinter der schließenden Klammer an).
- **verifizierbar:** ja — der obige `--print-mk`-Zähler.
- **klasse:** `enumeration-hinter-dem-modulsatz`

### L-2 · Das Handbuch verkauft den Preis als Vorteil

- **quelle:** `AGENTS.md` §3.8 (Grenze dort benennen, wo sie gilt)
- **pfad:** `docs/user/benutzerhandbuch.md:842`
- **zugesagt:** Lastenheft `:98`/`:3343` nennt dieselbe Tatsache als **Preis**: *„wer ein älteres Image fährt, liest ein neueres Handbuch."*
- **gemessen:** das Handbuch formuliert sie ausschließlich als Nutzen: *„so kann er nicht veralten, und Sie lesen den aktuellen Stand, auch wenn Sie ein älteres Image pinnen."* Die Kehrseite — der Adopter mit altem Pin liest über Fähigkeiten, die sein Bild nicht hat — fehlt an genau der Stelle, an der dieser Leser steht. Abgeschwächt dadurch, dass der Handbuch-Kopf die dokumentierte Software-Version stempelt (Release-Prep, `9d22a44`).
- **verifizierbar:** nein.
- **klasse:** `grenze-nur-auf-der-produktseite-benannt`

### I-1 · Die URL existiert jetzt zweimal ohne Kopplung

- **pfad:** `internal/adapter/driving/cli/cli.go:126` · `packaging/dockerhub/overview.md:62`
- **gemessen:** beide tragen zeichengleich `https://github.com/pt9912/d-check/blob/main/docs/user/benutzerhandbuch.md`. Kein Sensor koppelt sie: `d-check` scannt nur `*.md` (`spec/lastenheft.md:924`), erreicht die Go-Konstante also nie; `external` ist weder in `.d-check.yml` konfiguriert noch von einem Target oder Workflow aktiviert (gemessen: kein `external:`-Block, keine `--enable external`-Stelle). Eine Umbenennung des Repos oder eine Verschiebung der Datei bricht beide still. Slice §5 Risiko 1 nennt die Klasse, nicht die Verdopplung.

---

## Negativbefunde (geprüft, ohne Befund)

- **Vorzustand „null URLs" (Slice §1, Lastenheft-Historie).** Statisch über das ganze `cli`-Paket bei `84f7259~1` gemessen: `://` kommt ausschließlich in `config_template.go:239/241` vor — das ist die `--print-config`-Ausgabe. Weder `writeUsage` noch `mkTemplate` trugen eine URL. Die Aussage gilt damit für die **ganzen** beiden Ausgaben, nicht nur den Trailer; die Menge ist in §7 (`BEO-020`, Zähler 4 — Registerstand bestätigt) ausdrücklich deklariert. ✔
- **Nachzustand gegen das echte Binary.** `--help` → stderr (stdout 0 Byte), Exit 0, genau eine URL-Zeile; `--print-mk` → stdout Zeile 12, im Kopfkommentar **vor** `DCHECK_IMAGE ?=`, genau eine URL-Zeile. Beide Akzeptanzkriterien treffen zu. ✔
- **Die drei behaupteten Mutationsproben, nachgefahren** (Kopien aus `git archive 50d461e`, `make test IMAGE=…`): Handbuch-Zeilen aus `--help` entfernt → **2 rot** (`TestCLI001_HilfeNenntHandbuch`, `TestHandbuchURL_TraegtKeineVersion`); aus dem `mkTemplate` entfernt → **2 rot** (`TestCLI010_…`, `TestHandbuchURL_…`); URL auf `blob/v0.68.0` → **3 rot**. Exakt wie in `ca43027` behauptet. ✔
- **`TestHandbuchURL_TraegtKeineVersion` bei einseitiger Version** (Probe A: nur `--help` bekommt ein Literal mit `blob/v0.68.0`, das Fragment bleibt auf `main`): rot, und die Meldung benennt die betroffene Ausgabe — `--help: Handbuch-URL zeigt nicht auf main` und `--help: Handbuch-URL traegt eine Version`. Der `t.Errorf`-Loop über beide Ausgaben trägt. ✔
- **Reichweite / §3.8.** Kein Gate dieses Repos löst die URL je auf: `external` nicht konfiguriert und nirgends aktiviert; `d-check` scannt nur `*.md`. Die Risiko-Formulierung in Slice §5 und in der Lastenheft-Historie ist zutreffend. ✔
- **„Dieselbe Form wie die Docker-Hub-Overview" (Slice §2).** Zeichengleich mit `packaging/dockerhub/overview.md:62`. ✔
- **„README.md führt das Handbuch als *(German)*" (Slice §2 Punkt 3).** `README.md:259` — `is the [user handbook](docs/user/benutzerhandbuch.md) (German)`. ✔
- **„zwei benannte Stellen im Handbuch" (Slice §2 Punkt 2).** `docs/user/releasing.md:32-37` nennt genau zwei bare Tags ohne `ghcr`-Präfix (§Versionen und Tags, Docker-Hub-Pull in §Docker-Image). ✔
- **Release-Prep-Kopplung (`50d461e`).** Handbuch- und Software-Version bleiben zu Recht unberührt: `releasing.md` §4 führt den Kopfstempel als Release-Prep-Aufgabe, und `9d22a44` hat ihn dort gesetzt. ✔
- **Kein `CHANGELOG.md`-Eintrag ist Repo-Praxis, kein Defekt.** Gemessen: die Feat-Commits von slice-177 (`e498646`-Kette), slice-179 (`d009e17`) und slice-180 (`6835944`) fassen `CHANGELOG.md` nicht an; geschrieben wird er im Release-Prep-Commit (`9d22a44`, `63534fc`). ✔
- **§3.6 nicht ausgelöst.** Keine Schwelle gesenkt, kein Gate gelockert; das Vergleichs-ADR ADR-0068 existiert, *weil* es eine Hard Rule präzisiert. Die Feststellung von Slice §3 trägt insoweit. ✔
- **§3.4 Abwärts-Sperre.** Kein ADR-, Slice-, Wellen- oder Commit-Token im neuen Lastenheft-Text; `doc-check` (Modul `matrix`) grün. ✔
- **Lifecycle-Move `3458610`.** Reiner Rename (`{open => in-progress}`, 0 Zeilen Inhalt) plus der gekoppelte Roadmap-Flip (Ruhe-Marker entfernt, Zeiger-Liste unberührt) — genau die MR-013/§3.3-Ausnahme; `make planning-check` grün. ✔
- **`d-check:cite`-Direktiven im Slice-Plan (MR-054).** `modul-05:213-214` ist tatsächlich *„Sub-Area-Wahl prüfen"*, `:219` tatsächlich *„Offene Beobachtungen sichten"* — beide zeigen auf die **vorschreibende** Zeile; `citations` grün. Der dritte Block trägt zu Recht keine (MR-053-Ziel ist repo-eigen). MR-013/MR-053/MR-054 sind alle innerhalb ihres Geltungsbereichs zitiert (Felder gelesen). ✔
- **§3.1 Docker/make-only.** Kein Host-Toolchain-Aufruf im Diff; die neue Template-Zeile führt kein `%` und lässt die zugesagten sieben `fmt`-Verben (`print_mk.go:26-28`) unberührt. ✔
- **Hexagon-Richtung.** Kein neuer Import; `handbuchURL` wird paketintern konsumiert. `make arch-check` `0 Befund(e)`. ✔
- **Kommentar-Klassen im Übrigen.** `cli_handbuch_link_test.go:8-11` (warum das Literal statt des Imports) ist eine saubere Abgrenzung; `:14-16` und `:30-32` tragen Zusage plus Rang-Zeiger, je **ein** auflösbares Herkunfts-Feld. Kein Slice-Nummer-, Review- oder Mess-Label in Code oder Test. ✔
- **Zustandsfelder im Slice-Plan.** Kein `**Status:**`, `**Lifecycle:**` vorhanden, `**Verantwortlich:**` bei der Beanspruchung gesetzt (`361ed39`), `**Berührte Spec-Stellen:**` mit Kennungen statt §-Ankern. ✔

---

## Repo-unverändert-Beleg

```
$ git status --porcelain | wc -l
0
$ git rev-parse HEAD
ff21044adbeab7ee20c07291fae72f94b67b99f7
```

Alle Proben liefen in Kopien unter `…/scratchpad/rev181/` (`git archive 50d461e`), das Verzeichnis ist entfernt; die sechs Proben-Images (`dcheck-rev181*`) sind gelöscht (`docker images | grep rev181` leer). `make gates`/`make record-gates` nicht gefahren, `.harness/state/` unberührt.
