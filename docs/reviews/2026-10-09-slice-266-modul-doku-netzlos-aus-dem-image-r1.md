# Review R1 — slice-266: Die Doku der Module ist netzlos aus dem Image lesbar

**Review-Art:** Code (Diff gegen Plan, ADRs, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-266, Commit `ce2a060d` (`git show HEAD`)
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-266 (§1 Ziel, Abgrenzung und Auftraggeber-Entscheid
zur Form, §3 Plan, §8 Spiegel-Liste); eingehender CR
`docs/plan/cr/2026-10-09-cr-eingehend-sf-connector-reviews-zusage.md`, Punkt 4;
[`DC-FA-CLI-013`](../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen)
samt Verfeinerung `DC-FA-CLI-013.a`;
[ADR-0002](../plan/adr/0002-distribution-ghcr-image.md),
[ADR-0005](../plan/adr/0005-modul-layout-hexagon-ordner.md),
[ADR-0012](../plan/adr/0012-kern-paketschnitt-model-rules-app.md),
[ADR-0029](../plan/adr/0029-arch-check-via-a-check.md),
[ADR-0104](../plan/adr/0104-test-nachweise-entlasten-in-der-rtm.md);
[`MR-025`](../../harness/conventions.md#mr-025); `AGENTS.md` §3.1, §3.2, §3.6, §3.7, §3.8, §5 Regel 17.
Vorherige Findings am selben Modul: keine Reports, die das Modul-Root, `go:embed`
oder eine schichtlose Datei behandeln.

## Messungen (Kommando und Ergebnis)

Image `d-check:latest` (gebaut aus dem Stand von HEAD), Mutationen im Wegwerf-Klon
unter dem Scratchpad, nicht im Arbeitsbaum.

| # | Fall | Ergebnis |
|---|---|---|
| A | `make arch-check` auf HEAD | grün, mit Hinweis „1 gescannte Datei(en) liegen in keiner Schicht und bleiben ungeprüft: manual.go" |
| B | `make arch-check` auf HEAD~1 | grün, **ohne** diesen Hinweis |
| M1 | Klon: `manual.go` importiert `net/http` und `os` | rot: zwei `tech-leak` — die Tech-Regeln gelten auch im Modul-Root |
| M2 | Klon: `manual.go` importiert `internal/adapter/driven/fs`; `internal/hexagon/core/app/manual.go` importiert das Root-Paket | **ein** Befund: `wrong-direction: (ohne Schicht) -> adapters`; der Import Kern → Root bleibt **ohne** Befund |
| M3 | Klon: `//go:embed` des Handbuchs zeigt auf `README.md`; `make test` | rot: `TestEingebetteteDokumenteGleichDerQuelle` („weicht von der Quelldatei ab") samt beiden CLI-Tests — richtiger Grund; das Root-Paket läuft in `make test` mit |
| C | `--manual "Spezifikation — d-check"` bzw. `"Benutzerhandbuch: d-check"`, Kopfzeile entfernt, `cmp` gegen die Quelldatei | je eine Kopfzeile `:1`, beide byte-gleich |
| D | `--manual zzzq` | Exit 2, stdout leer, Hinweis nennt beide Titel |
| E | `--manual ""`, `--manual "   "`, `--manual` ohne Wert | Exit 2, stdout leer |
| F | `--manual planning` mit je `--repair-broad`, `--suggest-config x`, `--staged`, `--range a..b`, `--commit-msg -`, `--print-mk`, `--require-complete` | jeweils Exit 2, stdout leer |
| G | `--manual planning` mit `--enable links`, `--disable links`, `--config nix.yml`, `--id-prefix AC`; zweimal `--manual` | Exit 0 — Werte still ignoriert bzw. der letzte gewinnt; `--print-config` verhält sich mit denselben Optionen gleich |
| H | `--manual planning`: Ausgabe | vier Abschnitte (Handbuch 936, 1149, 1907; Spezifikation 2138), eine Leerzeile dazwischen, Ende ohne Leerzeile |
| I | `IMAGE_REF=d-check:slice232-fix bash tools/image-test.sh` (Image ohne `--manual`) | Phasen 1–4 grün, dann `FAIL — --manual: Exit 2 im Container ohne Netz und ohne Mount, want 0` — richtiger Grund |
| J | dasselbe gegen `d-check:latest` | Phase 5 grün, Gesamtlauf OK |

## Findings

### F-1 — Ein Paket außerhalb des Layouts von ADR-0005 und außerhalb jeder a-check-Schicht

- **kategorie:** MEDIUM
- **quelle:** ADR-0005 (Entscheidungs-Tabelle Modul-Layout); `AGENTS.md` §3.8; Skill-Prüffragen 3 und 15
- **pfad:** `manual.go` · „Das Paket liegt im Modul-Root, weil go:embed nur Dateien unter dem eigenen Verzeichnis erreicht."; `.a-check.yml` · `layers:`
- **befund:** Die Tabelle von ADR-0005 legt das Modul-Layout abschließend fest und kennt kein Paket im Modul-Root. Weder eine ADR noch `.a-check.yml` gibt dem neuen Paket eine Rolle. Der Gate-Lauf bleibt grün und meldet selbst, dass `manual.go` keiner Schicht angehört und bei den Kanten ungeprüft bleibt (A gegen B). Die Tech-Regeln greifen dort weiter (M1). Die Kantenregeln greifen nur in einer Richtung: Ein Import aus dem Kern in das Root-Paket bleibt ohne Befund (M2). Ein späterer Kern-Code kann also von einem schichtlosen Paket abhängen, und `make arch-check` meldet es nicht.
- **verifizierbar:** ja — `make arch-check` (Hinweiszeile) und Mutation M2.
- **klasse:** paket-ausserhalb-der-schichten

### F-2 — Der Drift-Test liegt außerhalb der Abdeckungs-Ableitung

- **kategorie:** MEDIUM
- **quelle:** ADR-0104 (Entscheidung 2); `AGENTS.md` §3.8; Skill-Prüffragen 15 und 18
- **pfad:** `manual_test.go` · „Die eingebetteten Dokumente sind byte-gleich zu ihren Quelldateien"; `internal/adapter/driven/configyaml/abdeckung_test.go` · `for _, dir := range []string{"internal", "cmd"}`
- **befund:** Der Test nennt `DC-FA-CLI-013` im Doc-Kommentar. Die Ableitung durchläuft aber nur `internal/` und `cmd/`, deshalb fehlt die Zeile in `docs/user/abdeckung-tests.md`. ADR-0104 und `harness/sensors/test.md` beschreiben die Regel als „eine Kennung im Doc-Kommentar einer Testfunktion" und nennen diese Einschränkung nicht. Nur die abgeleitete Datei selbst nennt sie, und erst mit diesem Commit gibt es einen Test außerhalb der beiden Wurzeln. Folge: Die RTM zeigt `DC-FA-CLI-013` mit zwei CLI-Tests belegt, und keiner davon prüft das Kriterium „Drift". Der einzige Test, der es prüft, ist dort unsichtbar. Hätte eine Anforderung ihren einzigen Test im Modul-Root, bliebe sie eine Waise, obwohl ihr Test eine Kennung deklariert.
- **verifizierbar:** ja — `grep manual_test docs/user/abdeckung-tests.md` liefert nichts.
- **klasse:** abdeckung-scan-menge-ohne-modul-root

### F-3 — Image-Test: Kopf, Schlusszeile und Sensor-Datei nennen nur DC-FA-DIST-001

- **kategorie:** LOW
- **quelle:** [`MR-025`](../../harness/conventions.md#mr-025) (Spiegel)
- **pfad:** `tools/image-test.sh` · „image-test: OK — DC-FA-DIST-001-Akzeptanzkriterien erfüllt"; `harness/sensors/image-test.md` · „Die Akzeptanzkriterien von"
- **befund:** Phase 5 prüft `DC-FA-CLI-013`, aber der Kopf des Skripts (Zeile 2), die Schlusszeile sowie Vertrag und Bindung der Sensor-Datei nennen weiterhin nur `DC-FA-DIST-001`, `DC-QA-02` und ADR-0102. Für Phase 4 (`DC-FA-CLI-007`/`008`) bestand dieselbe Lücke schon vorher; neu ist, dass sie mit diesem Commit eine zweite Anforderung betrifft.
- **verifizierbar:** nein — kein Gate liest die Sensor-Datei gegen die Phasen.
- **klasse:** spiegel-nicht-nachgezogen

### F-4 — Die Hilfe-Ausgabe verweist für das Handbuch nur auf das Netz

- **kategorie:** LOW
- **quelle:** [`MR-025`](../../harness/conventions.md#mr-025); Plan §8 nennt die Hilfe-Ausgabe als Spiegel
- **pfad:** `internal/adapter/driving/cli/cli.go` · „Benutzerhandbuch (aufgabenorientiert, deutsch):"
- **befund:** `--manual` steht in der Flag-Liste. Der Block der Hilfe, der sagt, *wo* das Handbuch liegt, nennt aber nur die beiden GitHub-URLs. Wer ohne Netz die Hilfe liest, findet an der Stelle, die ihm das Handbuch zeigt, keinen Hinweis auf den netzlosen Weg.
- **verifizierbar:** nein.
- **klasse:** spiegel-nicht-nachgezogen

### F-5 — Eigene Überschriften-Erkennung neben der vorhandenen

- **kategorie:** LOW
- **quelle:** Skill-Prüffrage 21 (seit slice-267)
- **pfad:** `internal/hexagon/core/app/manual.go` · „manualHeadingText ist der Text einer ATX-Überschrift ohne die führende"; `internal/adapter/driving/cli/manual.go` · `strings.HasPrefix(l, "# ")`
- **befund:** `parseATXHeading` liefert den Überschriften-Text bereits. Weil `FindSectionHeads` dem Prädikat nur die Rohzeile übergibt, leitet `manualHeadingText` den Text ein zweites Mal ab. Heute ist das Ergebnis gleichwertig (gelesen am Code, nicht gemessen). `titel()` erkennt die Titelzeile mit einem dritten Muster, das Fenced-Code nicht beachtet. Ein konkretes Versagen tritt bei den beiden heutigen Dokumenten nicht auf (C), deshalb LOW.
- **verifizierbar:** nein.
- **klasse:** eigene-erkennung-neben-produkt

### F-6 — Der Drift-Test hält die Kopplung, die Byte-Gleichheit liefert `go:embed`

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `manual_test.go` · `TestEingebetteteDokumenteGleichDerQuelle`
- **befund:** Einbettung und `os.ReadFile` lesen denselben Baum. Der Test kann deshalb nur rot werden, wenn Pfad-Konstante und `//go:embed`-Direktive auseinanderlaufen (M3). Genau das prüft er, und er ist in diesem Sinn richtig. Die Byte-Gleichheit der *Ausgabe* ergibt sich aus der Konstruktion; gemessen ist sie hier nur von Hand (C). Wer den Test als Drift-Sensor zwischen Image und Repo liest, liest ihn weiter, als er reicht.
- **verifizierbar:** ja — M3.
- **klasse:** test-haelt-kopplung-nicht-zusage

### F-7 — Nicht-Modus-Optionen neben `--manual` werden still angenommen

- **kategorie:** INFO
- **quelle:** `DC-FA-CLI-013.a` Schritt 1
- **pfad:** `internal/adapter/driving/cli/manual.go` · „--manual ist ein eigener Modus"
- **befund:** `--enable`, `--disable`, `--config`, `--id-prefix` und ein zweites `--manual` führen zu Exit 0 (G). Die Spezifikation zählt nur Modus-Optionen auf, und `--print-config` verhält sich gleich. Das ist also konsistent und kein Vertragsbruch, aber eine undokumentierte Annahme.
- **verifizierbar:** ja — Fall G.
- **klasse:** stille-nebenoption

## Negativbefunde

- **Lastenheft ↔ Spezifikation ↔ Code:** geprüft. Treffer nur in Überschriften (Code: Prädikat auf ATX-Zeilen; D, Test `KeinTreffer`), Vergleich in Kleinbuchstaben auf beiden Seiten, Begriff getrimmt, verschachtelte Treffer einmal (`until`-Sprung), Fences weder Treffer noch Grenze (`FindSectionHeads`/`SectionEnd`, dieselbe Lexik wie `structure`), Kopfzeile `==> <pfad>:<zeile>` 1-basiert, Leerzeilen am Abschnittsende entfernt, eine Leerzeile zwischen Abschnitten (H), der Titel liefert das ganze Dokument byte-gleich (C), kein Treffer, leerer Begriff und jede Modus-Kombination ergeben Exit 2 ohne stdout (D–F). Die Modus-Liste in Lastenheft und Spezifikation (13 Optionen) deckt sich mit `manualComboError`; `--repair-broad` setzt `o.repair` mit (`repair: *repairOut || *repairBroadOut`). Ohne Befund.
- **Repo-frei / DC-QA-03:** geprüft. `earlyGenerators` läuft vor `openRoot`; die Scan-Wurzel existiert nicht, im Container gibt es weder Netz noch Mount (Test, I/J). Kein Netzzugriff, keine Schreiboperation.
- **Image-Test Phase 5:** geprüft. Gegen ein Image ohne `--manual` wird sie aus dem richtigen Grund rot (I), gegen HEAD grün (J). Der Abgleich nativ gegen Container ist byte-genau (`cmp -s`).
- **Drift-Test in `make test`:** geprüft. `go test ./...` schließt das Root-Paket ein; M3 wird rot.
- **CLI-Tests über das Produkt:** geprüft. Beide getaggten Tests laufen über `cli.Run`, also den Produkt-Eingang, nicht über die interne Funktion. Die Kriterien „ganzes Dokument", „verschachtelt" und „Code-Block" prüfen die ungetaggten Tests in `app` (für die RTM unsichtbar, beanspruchen dort aber auch nichts). Ob die Tests die Kriterien vollständig belegen, ist Sache der Verifikation.
- **ADR-0012:** geprüft. `app` → `rules` ist eine erlaubte Kante; `ManualSections` ist rein und ohne I/O.
- **ADR-0002:** geprüft. Das Image bleibt das eine Binary, die Dokumente reisen eingebettet mit. Es kommt keine Datei hinzu, und der Build-Kontext hat keine `.dockerignore`, die die Quellen ausschlösse.
- **§3.2 / §3.6:** kein `//nolint`, keine gesenkte Schwelle. `coverage-gate` misst weiter über `internal/...`, das Root-Paket enthält keine Logik.
- **§3.7 Kommentare:** `manual.go` (Zusage plus Grenze zum Ort des Pakets), `cli/manual.go`, `app/manual.go`, die Testkommentare und die Phase-Kommentare in `image-test.sh` tragen Zusage, Abgrenzung oder ein Kennungsfeld. Keine Review-Historie, keine Slice-Nummer. Ohne Befund.
- **Abgrenzung des Plans:** keine Mitnahme — weder die `reviews`-Regel (slice-265) noch ein neuer Handbuchtext. README und Handbuch sind unberührt (`AGENTS.md` §5 Regel 17), und keine Aussage dort ist durch den Diff falsch geworden (geprüft per `grep` nach repo-frei/netzlos-Aussagen).
- **Spec-Historien:** Lastenheft 0.105.0 und die Historie-Zeile der Spezifikation sind vorhanden. Keine Abwärts-Referenz in den Spec-Straten.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 2 | 3 | 2 |

Wiederkehrende Klasse: Modul-Grenze auf der Scan-Achse (`AGENTS.md` §3.8) — zweimal in diesem Lauf (F-1 a-check-Schichten, F-2 Abdeckungs-Ableitung). Beide Male hat das neue Modul-Root eine Scan-Menge verlassen, die bisher stillschweigend „alle Go-Pakete" bedeutete.

## Verdikt

**Nachbessern vor Closure.** F-1 braucht eine Entscheidung über den Ort und die Rolle des Root-Pakets: eine Deklaration in `.a-check.yml` und/oder eine ADR zu ADR-0005, gefällt vom Architect und nicht hier. F-2 braucht eine Entscheidung über die Scan-Menge der Ableitung oder deren benannte Grenze. Das Produktverhalten selbst entspricht in allen gemessenen Fällen Lastenheft und Spezifikation.
