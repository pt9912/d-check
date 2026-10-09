# Review R2 — slice-266: Die Doku der Module ist netzlos aus dem Image lesbar

**Review-Art:** Code (Diff gegen Plan, ADRs, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-266, Commits `69e52c11` (Plan-Änderung), `d04d5c8a` (Löschung `manual_test.go`, Botschaft falsch), `dbae92e5` (Korrekturen zu R1)
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** R1-Report zu slice-266 (F-1 bis F-7); Slice-Plan slice-266 nach der
Plan-Änderung (§2, §3); [ADR-0005](../plan/adr/0005-modul-layout-hexagon-ordner.md),
[ADR-0012](../plan/adr/0012-kern-paketschnitt-model-rules-app.md),
[ADR-0029](../plan/adr/0029-arch-check-via-a-check.md),
[ADR-0104](../plan/adr/0104-test-nachweise-entlasten-in-der-rtm.md),
[ADR-0107](../plan/adr/0107-mitgelieferte-dokumente-im-modul-root.md) (Proposed);
[`DC-FA-CLI-013`](../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen),
`DC-FA-CLI-001.a`; [`MR-025`](../../harness/conventions.md#mr-025); `AGENTS.md` §3.5, §3.7, §3.8, §5 Regel 15, §6 Schritt 4.

## Messungen (Kommando und Ergebnis)

Wegwerf-Klon auf `dbae92e5` unter dem Scratchpad.

| # | Fall | Ergebnis |
|---|---|---|
| A | `make arch-check` auf HEAD | grün, der Hinweis „in keiner Schicht" aus R1 ist weg |
| M1 | `app/manual.go` importiert das Root-Paket; `make arch-check` | grün — a-check sieht die Kante weiterhin nicht (Angabe aus ADR-0107 bestätigt) |
| M1' | M1 plus `driven/report/report.go` importiert das Root-Paket; `make test` | rot: `TestManual_RootPaketNurAusDerCompositionRoot` nennt beide Dateien — richtiger Grund |
| M2 | `manual.go` importiert `strings` und `internal/hexagon/core/rules` | rot: `core-impurity: Kern importiert …/core/rules` |
| M3 | `manual.go` importiert nur zusätzlich `strings` | grün |
| B | `ls internal internal/adapter/driving` | außer `hexagon`, `adapter/driven` und `adapter/driving/cli` gibt es nichts — der Walk des Kanten-Tests deckt alle Pakete außerhalb der Composition Root (`cmd/` ist Composition Root) |

## Abgleich mit R1

| R1 | Stand |
|---|---|
| F-1 Schicht und Entscheidung | **erledigt.** Schicht `docs` in `.a-check.yml`, ADR-0107 Proposed samt Index-Zeile, Plan-Änderung (`69e52c11`) vor dem Code. Die Kante, die a-check nicht sieht, hält ein Test, der aus dem richtigen Grund rot wird (M1'). Neue Restpunkte: F-1 und F-2 unten |
| F-2 Drift-Test außerhalb der Ableitung | **erledigt.** Der Test liegt in `internal/adapter/driving/cli/manual_root_test.go` und steht in `abdeckung-tests.md` |
| F-3 Image-Test-Spiegel | **teilweise.** Skriptkopf, Schlusszeile, Vertrag der Sensor-Datei und Bindung (`DC-FA-CLI-013`) sind nachgezogen. Rest siehe F-4 |
| F-4 Hilfe | **erledigt.** Die Hilfe-Zeile und `DC-FA-CLI-001.a` samt Historie sind nachgezogen |
| F-5 eigene Erkennung | **erledigt.** `rules.HeadingText` statt eigener Ableitung; `app.ManualTitle` über `FindSectionHeads` mit Fence-Behandlung und Test |
| F-6, F-7 (INFO) | unverändert, ohne Handlungsbedarf |

## Findings

### F-1 — ADR-0107 sagt „keine Logik, keine Importe außer `embed`" zu, und nichts hält das

- **kategorie:** MEDIUM
- **quelle:** ADR-0107 (Entscheidung 1, Konsequenzen „Grenze"); Skill-Prüffrage 18; `AGENTS.md` §3.8
- **pfad:** `docs/plan/adr/0107-mitgelieferte-dokumente-im-modul-root.md` · „trägt sonst nichts: keine Logik, keine Importe außer"
- **befund:** Die Rolle `domain` fängt Tech-Muster und Importe interner Pakete (M2), aber keinen Standardbibliotheks-Import und keine Logik (M3 grün). Das Root-Paket liegt außerdem außerhalb von `internal/...` und damit außerhalb der Messbasis von `make coverage-gate`. Die Grenze der ADR nennt nur die Kante Kern → Root. Folge: Wandert Logik ins Root-Paket, meldet weder `make arch-check` noch das Coverage-Gate etwas, obwohl die ADR genau das ausschließt.
- **verifizierbar:** ja — Mutation M3 plus `make arch-check`.
- **klasse:** grenzen-liste-ohne-groesste-luecke

### F-2 — Eine Begründung in ADR-0107 trifft am Repo nicht zu

- **kategorie:** MEDIUM
- **quelle:** Skill-Prüffrage 20; `AGENTS.md` §3.1, §3.5
- **pfad:** `docs/plan/adr/0107-mitgelieferte-dokumente-im-modul-root.md` · „ein Build-Schritt, den `make test` außerhalb von Docker nicht kennt"
- **befund:** In diesem Repo gibt es kein `make test` außerhalb von Docker (`AGENTS.md` §3.1). Ein Kopierschritt stünde in der Test-Stage des Dockerfiles und wäre dort bekannt. Das tragende Contra der Option ist die zweite Fassung im Repo bzw. die Drift, und die bleibt bestehen. An der Entscheidung ändert sich also nichts, die Begründung ist trotzdem falsch. Wird die ADR `Accepted`, ist sie unveränderlich (§3.5), und eine Re-Evaluierung liest dann einen Grund, den es nicht gibt.
- **verifizierbar:** nein.
- **klasse:** begruendung-trifft-gegenstand-nicht

### F-3 — Die Kommentarzeile in `.a-check.yml` schreibt a-check eine Zusage zu, die der Test trägt

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §3.7 (Grenze); ADR-0107 Entscheidung 4
- **pfad:** `.a-check.yml` · „Ohne Kante — nur die Composition Root liest es."
- **befund:** Am Ort der a-check-Konfiguration liest sich der Satz als Zusage dieses Gates. a-check hält aber nur die ausgehende Hälfte (M2); die eingehende hält `TestManual_RootPaketNurAusDerCompositionRoot` (M1, M1'). Der Kommentar nennt diese Grenze nicht und zeigt nicht auf den Test.
- **verifizierbar:** ja — M1.
- **klasse:** zusage-am-falschen-ort

### F-4 — Bindungs-Spiegel des Image-Tests unvollständig

- **kategorie:** LOW
- **quelle:** [`MR-025`](../../harness/conventions.md#mr-025)
- **pfad:** `harness/sensors/image-test.md` · Abschnitt „Bindung"; `harness/README.md` · Zeile `` [`make image-test`] ``
- **befund:** Der Vertrag der Sensor-Datei nennt jetzt `DC-FA-CLI-007`, `-008` und `-013`. Die Bindung nennt davon nur `-013`, und die Gate-Index-Zeile in `harness/README.md` nennt weiterhin nur `DC-FA-DIST-001`/`DC-QA-02`. Drei Stellen beschreiben dasselbe Skript verschieden.
- **verifizierbar:** nein.
- **klasse:** spiegel-nicht-nachgezogen

### F-5 — `d04d5c8a` trägt eine Botschaft, die seinen Inhalt nicht beschreibt

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §5 Regel 15
- **pfad:** Commit `d04d5c8a` · „Plan-Aenderung nach Review R1 — Schicht und ADR fuer das Root-Paket"
- **befund:** Der Commit löscht nur `manual_test.go`. Seine Botschaft ist wortgleich mit der von `69e52c11` und behauptet eine Plan-Änderung. `dbae92e5` stellt das in seiner Botschaft richtig, damit ist die Kette lesbar. `git log -- manual_test.go` zeigt die Löschung trotzdem unter der falschen Botschaft. Kein Gate fängt das, `trace-check` prüft nur die Kennung.
- **verifizierbar:** nein.
- **klasse:** botschaft-beschreibt-inhalt-nicht

### F-6 — Rolle `domain` für `docs`: in der Wirkung passend, nirgends begründet

- **kategorie:** INFO
- **quelle:** ADR-0107 Entscheidung 3; ADR-0005/ADR-0012 (Begriff „Kern")
- **pfad:** `.a-check.yml` · `docs:     { globs: ["manual.go"], role: domain }`
- **befund:** `domain` ohne Kanten ist die strengste verfügbare Einstufung: Tech-Muster sind kategorisch verboten, und jeder Import eines geschichteten Pakets wird gemeldet (M2). Das passt zu „trägt sonst nichts". a-check nennt einen Verstoß aber „Kern importiert …". In ADR-0005 und ADR-0012 ist der Kern `internal/hexagon/core`, und ADR-0107 Entscheidung 2 grenzt das Root-Paket ausdrücklich vom Kern ab. Warum gerade `domain` gewählt ist, sagt die ADR nicht.
- **verifizierbar:** ja — M2.
- **klasse:** undokumentierte-annahme

### F-7 — ADR-0005 zeigt nicht auf ihre Erweiterung

- **kategorie:** INFO
- **quelle:** ADR-0005 (Geschichte); `AGENTS.md` §3.5
- **pfad:** `docs/plan/adr/0005-modul-layout-hexagon-ordner.md` · Abschnitt „Geschichte"
- **befund:** Erweiterung ohne `Supersedes` ist hier richtig (siehe Negativbefunde). Wer ADR-0005 liest, findet ADR-0107 aber nur über den Index. Für die Teil-Ablösung durch ADR-0029 trägt die Geschichte von ADR-0005 eine Zeile; ein solcher Anhang ist nach §3.5 erlaubt.
- **verifizierbar:** nein.
- **klasse:** rueckverweis-fehlt

## Antworten auf die zwei Zusatzfragen

- **ADR-0107 als Erweiterung ohne `Supersedes`:** richtig. Keine Aussage von ADR-0005 wird falsch: Die Layout-Tabelle wird um einen Ort außerhalb von `internal/` ergänzt, und die Import-Regeln 1–5 gelten unverändert. `Supersedes` ist für Korrekturen da (§3.5). Hier wird nichts korrigiert, sondern ein Fall entschieden, den ADR-0005 nicht kannte. Offen bleibt nur der Rückverweis (F-7).
- **Rolle `domain` ohne Kanten:** In der Wirkung passend, weil sie die strengste Einstufung ist (M2). Sie hält aber „nur `embed`" nicht (F-1), und weder ADR noch Kommentar begründen die Wahl. Die Meldungssprache „Kern" kollidiert mit dem Kern-Begriff der ADRs (F-6).

## Negativbefunde

- **Plan-Änderung vor Code (§6 Schritt 4):** geprüft. `69e52c11` (§2-DoD-Punkt, zwei §3-Zeilen) liegt vor `dbae92e5`. Die neue öffentliche Funktion `rules.HeadingText` steht nicht als eigene §3-Zeile im Plan; sie ist eine dünne Hülle um die vorhandene Erkennung und berührt keine Ausnahme aus §1. Ohne Befund.
- **Kanten-Test:** geprüft. Der Walk deckt jedes Paket außerhalb der Composition Root (B), vergleicht den Importpfad exakt (Alias- und Blank-Import werden mit erfasst) und wird aus dem richtigen Grund rot (M1').
- **Drift-Test nach dem Umzug:** geprüft. Der relative Pfad zur Wurzel stimmt mit der Tiefe des Pakets überein, und der Test steht in der Ableitung.
- **`rules.HeadingText` / `app.ManualTitle`:** geprüft. Eine Erkennung, ein Test für den Fence-Fall. Kommentare tragen Zusage und Kopplung (§3.7).
- **Spezifikation `DC-FA-CLI-001.a`:** geprüft. Text und Hilfe-Ausgabe stimmen in der Reihenfolge überein (URLs, dann `--manual`), die Historie-Zeile ist vorhanden.
- **ADR-0107 Form:** geprüft. Re-Evaluierungs-Trigger vorhanden, Index nachgezogen. Der Provenance-Marker zeigt auf die Herkunft und begründet nichts. `Schärft: —` ist vertretbar, weil die Entscheidung kein Rollen-Element der Sicht berührt.
- **Hexagon, Netz, Suppressions:** keine neuen Produkt-Imports außer `app` → `rules` (erlaubte Kante), kein `//nolint`, kein Netz. Testdateien mit `os` und `go/parser` sind von a-check ausgenommen (`exclude: **/*_test.go`).

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 2 | 3 | 2 |

Wiederkehrende Klassen: Grenze auf der Scan-Achse (§3.8, jetzt F-1: der Ausschnitt, den Gate und Coverage am Root-Paket sehen) und Prüffrage 20 (F-2).

## Verdikt

**Nachbessern vor `Accepted` von ADR-0107.** Die R1-Befunde sind in der Sache erledigt. F-1 und F-2 betreffen den Text der ADR, nicht das Verhalten, und lassen sich nur ändern, solange sie `Proposed` ist. F-3 bis F-5 sind kleine Spiegel- bzw. Botschafts-Punkte.
