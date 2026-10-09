# Verifikation — slice-266: Die Doku der Module ist netzlos aus dem Image lesbar

**Rolle:** Verifier (Baseline-Regelwerk `modul-11-verification.md`): DoD und Plan gegen den tatsächlichen Stand, nicht Diff gegen Plan
**Gegenstand:** slice-266 auf `0cf40674` (Feat `ce2a060d`, R1-Korrekturen `d04d5c8a` + `dbae92e5`, R2-Auflagen `0cf40674`); DoD-Punkte 1 bis 3
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan §1–§3, §6, §8;
[eingehender CR von `sf-connector`](../plan/cr/2026-10-09-cr-eingehend-sf-connector-reviews-zusage.md), Punkt 4;
[`DC-FA-CLI-013`](../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen)
samt [`DC-FA-CLI-013.a`](../../spec/spezifikation.md#dc-fa-cli-013a--handbuch-und-spezifikation-ausgeben---manual);
[ADR-0107](../plan/adr/0107-mitgelieferte-dokumente-im-modul-root.md);
Reviews R1/R2 zu slice-266 (nur als Hinweis, nicht als Beleg übernommen).

Alle Läufe in einem frischen Klon unter dem Session-Scratchpad, nur über `make`/Docker.
Am Arbeitsbaum ist nichts geändert außer dieser Datei.

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| V1 | `git clone` (HEAD `0cf40674`), `make gates IMAGE=d-check-v266` | Exit 0, `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `ok …/internal/adapter/driving/cli`, `ok …/internal/hexagon/core/app`; a-check `gesamt: 0 Befund(e)`; `coverage-gate: OK — Coverage 94.80% erfüllt Schwelle 93%`; `manualComboError`, `runManual`, `ManualSections`, `manualHeadingText` je 100,0 %, `ManualTitle` 83,3 %; doc-check `0 Befund(e)` |
| V2 | `docker run --rm --network none d-check-v266:latest --manual planning` (kein Mount) | Exit 0, stderr leer, Kopfzeilen `==> docs/user/benutzerhandbuch.md:936`, `:1149`, `:1907`, `==> spec/spezifikation.md:2139` — genau die vier echten Überschriften mit „planning" (`grep -n -i '^#.*planning'`). Abschnitt `spec:2139` per `cmp` byte-gleich zu den Quellzeilen 2139 ff.; er endet vor `### DC-FA-STRUCT-001.a` |
| V3 | Black-box-Kriterien gegen das Image (Tabelle unten) | alle wie zugesagt |
| V4 | Bewusstes Brechen im Klon, je `make test` / `make arch-check` (Tabelle unten) | siehe dort |
| V5 | `IMAGE_REF=ghcr.io/pt9912/d-check:v0.85.0 bash tools/image-test.sh` (Image per `docker pull`) | Exit 1; Phasen (1)–(4) grün, dann `image-test: FAIL — --manual: Exit 2 im Container ohne Netz und ohne Mount, want 0`; Ursache direkt: v0.85.0 meldet `flag provided but not defined: -manual`. Gegenprobe gegen `d-check-v266:latest`: `(5) --manual — netzlos im Container, nativ == Container, Exit 0`, `OK` |

Die Testbilder wurden danach entfernt.

### Black-box-Kriterien (V3, Image `d-check-v266:latest`, `--network none`, ohne Mount)

| Kriterium (`DC-FA-CLI-013`) | Aufruf | Ergebnis |
|---|---|---|
| Happy Path, Groß-/Kleinschreibung | `--manual "  PLANNING  "` (auch `--read-only`) | erste Kopfzeile `benutzerhandbuch.md:936`, wie mit `planning` |
| Ganzes Dokument (Titel) | `--manual "Benutzerhandbuch: d-check"` | eine Kopfzeile `:1`, Rest per `cmp` byte-gleich zum Handbuch ohne Leerzeilen am Ende |
| Verschachtelt | `--manual Spezifikation` (H1-Titel nennt den Begriff, ebenso `### DC-FA-CLI-013.a — Handbuch und Spezifikation ausgeben`) | genau **eine** Kopfzeile `spec/spezifikation.md:1`, der Unterabschnitt erscheint nicht ein zweites Mal; Inhalt byte-gleich zur Datei |
| Code-Block kein Treffer | `--manual Repository-Struktur` (nur `### F-1 — Repository-Struktur` im Fence, Handbuch:671); `--manual opt-in-Module` (nur `# Weitere opt-in-Module…` im Fence, Spez:202) | beide Exit 2, stdout leer |
| Code-Block kein Treffer, echte Überschrift daneben | `--manual Kreuzverweis-Konsistenz` (`## Kreuzverweis-Konsistenz` im Fence, Handbuch:836; echte Überschrift Spez:473) | nur `==> spec/spezifikation.md:473` |
| Code-Block keine Grenze | `--manual "4.15 "` (Abschnitt Handbuch:793 umschließt den Fence mit `## Kreuzverweis-Konsistenz`); `--manual DC-FA-CLI-006.a` (umschließt Spez:202–204) | Abschnitt läuft über die Fence-Zeile hinaus bis vor die nächste echte Überschrift |
| Kein Treffer | `--manual Repository-Struktur` | Exit 2, stderr `keine Überschrift nennt "Repository-Struktur" — der Titel eines Dokuments liefert es ganz: "Benutzerhandbuch: d-check", "Spezifikation — d-check"` |
| Leerer Begriff | `--manual ""`, `--manual "   "` | je Exit 2, stdout leer, `--manual braucht einen Begriff` |
| Kombination | `--manual planning` + je `--json`, `--yaml`, `--doctor`, `--repair`, `--repair-broad`, `--trace`, `--require-complete`, `--print-config`, `--print-mk`, `--suggest-config=ai-harness`, `--commit-msg=-`, `--range=HEAD~1..HEAD`, `--staged` | alle 13 Exit 2, stdout leer. `--require-complete` allein meldet `--require-complete erfordert --trace` (frühere Prüfung), mit `--trace` dazu die `--manual`-Meldung |

### Bewusstes Brechen (V4)

| # | Mutation im Klon | Sensor | Ergebnis |
|---|---|---|---|
| B1 | Handbuch Zeile 936 um ` VERIFY-DRIFT` ergänzt, nicht neu gebaut | `make test` | **grün** (`ok …/cli`). Erwartet — siehe F-1 |
| B2 | `//go:embed docs/user/benutzerhandbuch.md` → `//go:embed README.md` (Konstante `HandbuchPfad` unverändert) | `make test` | rot, `TestManual_EingebetteteDokumenteGleichDerQuelle: docs/user/benutzerhandbuch.md: eingebettete Fassung weicht von der Quelldatei ab` (dazu `TestManual_GibtAbschnitteAus`, `TestManual_Nutzungsfehler`) — richtiger Grund |
| B3 | `internal/hexagon/core/app/zz_verify.go` importiert `github.com/pt9912/d-check` | `make test` / `make arch-check` | test rot: `TestManual_RootPaketNurAusDerCompositionRoot: …/app/zz_verify.go importiert das Paket im Modul-Root`; arch-check **grün** (`gesamt: 0 Befund(e)`) — bestätigt die Grenze aus ADR-0107: die eingehende Kante hält nur der Test |
| B4 | `manual.go` importiert `internal/hexagon/core/app` und trägt eine Funktion | `make arch-check` / `make test` | arch-check rot: `manual.go:10: core-impurity: Kern importiert github.com/pt9912/d-check/internal/hexagon/core/app` — die Schicht `docs` greift für ausgehende Kanten; test rot: `TestManual_RootPaketTraegtNurDieEinbettung: … importiert …/core/app` und `… trägt eine Funktion` |
| B5 | `if h.Line < until` → `if false && …` (verschachtelte Treffer nicht mehr übersprungen) | `make test` | rot, `TestManualSections_TrefferUndVerschachtelung` |
| B6 | Treffer gegen die rohe Überschriftenzeile statt gegen den Text ohne `#`-Folge | `make test` | **grün** — siehe F-2 |

## DoD-Punkt 1 — netzlos aus dem Image, zugesagt, durch einen Test gehalten

**Erfüllt.**

- Zusage: Lastenheft 0.105.0, `DC-FA-CLI-013` mit acht Kriterien, Historie-Zeile; Spezifikation
  `DC-FA-CLI-013.a` Schritte 1–6 samt Grenze, Historie-Zeile. Hilfe nennt `--manual` (`DC-FA-CLI-001.a`
  nachgezogen).
- Verhalten: V2/V3 — alle Kriterien gegen das gebaute Image ohne Netz und ohne Mount bestätigt,
  auch hinter dem Pfad-Argument des `ENTRYPOINT ["/d-check", "/repo"]`.
- Test: `make test` hält das Verhalten (`TestManual_GibtAbschnitteAus` repo-frei über einen
  nicht existierenden Pfad, `TestManual_Nutzungsfehler`, die vier `app`-Tests); `tools/image-test.sh`
  Phase (5) hält „netzlos im Container". V5: Phase (5) wird gegen ein Image ohne `--manual` aus dem
  richtigen Grund rot. **Grenze:** `image-test` hängt an `ci`/`fullbuild`, nicht an `gates`; der
  Netzlos-Beleg im Container liegt damit im CI-Lauf, nicht im `make gates`-Nachweis.

## DoD-Punkt 2 — Test hält die mitgelieferte Doku gegen ihre Quelle; `make gates` grün

**Erfüllt, mit einer Präzisierung dessen, was der Test beweist (F-1).**

- `make gates` grün auf `0cf40674` (V1, selbst gefahren).
- `TestManual_EingebetteteDokumenteGleichDerQuelle` wird aus dem richtigen Grund rot, wenn die
  Einbettung auf eine andere Datei zeigt als die Pfad-Konstante (B2).
- Die Byte-Gleichheit zum Build-Stand selbst liefert `go:embed`: Test-Binary und Quelldatei
  stammen aus demselben Baum, eine „veraltete Einbettung" kann in `make test` nicht entstehen (B1).

## DoD-Punkt 3 — Folge-ADR (Proposed) zu ADR-0005; Schicht in `.a-check.yml`

**Erfüllt.**

- [ADR-0107](../plan/adr/0107-mitgelieferte-dokumente-im-modul-root.md): `**Status:** Proposed`,
  Bezug auf ADR-0005 (erweitert das Layout, ersetzt es nicht — kein `Supersedes`, sachgerecht), drei
  Alternativen plus Wahl, Konsequenzen mit Grenze, Fitness Function, Re-Evaluierungs-Trigger,
  Geschichte. Index-Zeile in `docs/plan/adr/README.md` vorhanden (`Proposed`, 2026-10-09).
- `.a-check.yml`: Schicht `docs: { globs: ["manual.go"], role: domain }` mit Kommentar (Zusage,
  Grenze, `DC-FA-CLI-013`, `ADR-0107`). B4 zeigt, dass die Schicht ausgehende Kanten wirklich prüft;
  B3, dass die eingehende Kante a-check entgeht und der Test sie hält — so wie die ADR es benennt.
- Plan-Änderung vor dem Code: `69e52c11` (Plan §3, zwei Zeilen) liegt vor `dbae92e5`.

## Findings

### F-1 (LOW) — Der Drift-Test hält die Pfad-Kopplung, nicht die Aktualität; die ADR sagt mehr

`TestManual_EingebetteteDokumenteGleichDerQuelle` vergleicht `dcheck.Handbuch` (zur Übersetzungszeit
eingebettet) mit `os.ReadFile(HandbuchPfad)` aus demselben Baum. Innerhalb eines `go test`-Laufs sind
beide per Konstruktion gleich; eine geänderte Handbuch-Zeile lässt ihn grün (B1). Rot wird er nur, wenn
`//go:embed`-Ziel und Pfad-Konstante auseinanderlaufen (B2) — also genau dann, wenn die Kopfzeile
`==> <pfad>:<zeile>` lügen würde. Das ist ein echter, sinnvoller Schutz; die Drift-Freiheit trägt aber
`go:embed`. R1 hat das als F-6 (INFO) benannt. Was bleibt: ADR-0107 §Fitness Function sagt, der Test
„hält die Einbettung gegen die Quelldateien" — gemessen hält er die **Zuordnung** von Pfad-Konstante und
Einbettung. Vor `Accepted` lässt sich der Satz noch präzisieren (Dokumentations-Regel 15).

### F-2 (LOW) — Spezifikations-Schritt 3 „ohne die führende `#`-Folge" ist ungetestet

Mutation B6 (Vergleich gegen die rohe Zeile `## …` statt gegen den Überschriftentext) überlebt
`make test`. Beobachtbar wäre sie an einem Begriff wie `#` oder `##`, der dann jede Überschrift träfe.
Kein Lastenheft-Kriterium nennt den Fall ausdrücklich; die Spezifikation legt ihn fest. Ein Testfall
(etwa `ManualSections(docs, "#")` ohne Treffer) schlösse die Lücke.

### F-3 (INFO) — Nicht-Modus-Optionen werden still ignoriert

`--manual planning --config x.yml` liefert Exit 0 und ignoriert `--config`; ebenso `--enable`/`--disable`.
Das widerspricht keiner Zusage — die Kombinationsliste nennt nur Modus-Optionen, und der Modus ist
repo-frei wie `--print-config`. Benannt, damit es als Entscheidung gelesen wird, nicht als Lücke.

### F-4 (INFO) — „Drift"-Kriterium und Ausgabeform

Das Lastenheft-Kriterium „die ausgegebenen Dokumente … byte-gleich" steht neben Schritt 5 der
Spezifikation (Kopfzeile, ohne Leerzeilen am Ende). V3 bestätigt: der Titel-Aufruf liefert das Dokument
byte-gleich bis auf Kopfzeile und abschließende Leerzeilen. Gemeint ist die Fassung im Binary; der
Wortlaut „ausgegeben" ist enger, als die Ausgabe es einlöst. Kein Handlungsbedarf vor der Closure.

## Negativbefunde

- Kein Host-Go/-Python/-Perl; alle Messungen über `make`/Docker.
- Kombinationsprüfung deckt alle 13 im Lastenheft genannten Optionen (`--repair-broad` über `repair`).
- `runManual` liest weder Konfiguration noch Scan-Wurzel (Kurzschluss in `earlyGenerators` vor
  `openRoot`); im Container ohne Mount Exit 0.
- Abschnitt-Inhalt byte-gleich zur Quelle (V2, V3); Fence-Zeilen sind weder Treffer noch Grenze.
- Abdeckungs-Dateien stimmen mit ihrer Ableitung (in `make test` gehalten, V1).

## Verdikt

**DoD-Punkte 1–3 erfüllt.** Das Verhalten ist gegen das gebaute Image black-box bestätigt, Image-Test
Phase (5) und die Kanten-/Einbettungs-Tests werden aus dem richtigen Grund rot. Zwei LOW-Befunde, beide
vor `Accepted` von ADR-0107 bzw. vor der Closure billig: F-1 den Fitness-Function-Satz auf das
präzisieren, was der Test hält; F-2 einen Testfall für Schritt 3 ergänzen. Kein Befund blockiert die
Closure.
