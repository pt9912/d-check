# Verifikation — slice-249: `hostpaths` erkennt Home-relative Pfade, Ziel-Ventil `exempt-targets`

- **Rolle:** Verifier (Modul 11) — Frage: „Bauen wir es richtig?" (DoD + Spec + Plan)
- **Gegenstand:** `slice-249` · Range `aa854193..0d3d10e3` — `21fde701` (Implementierung), `82d7e046` (R1-Report), `0d3d10e3` (R1-Einarbeitung)
- **Datum:** 2026-10-05 · **Modell-ID:** claude-opus-5-5
- **Eingangs-Kontext:** DoD-Bestätigung des Implementers samt Sensor-Belegen (Commit-Botschaften `21fde701`/`0d3d10e3`); Slice-Plan; `DC-FA-HOST-001` (Lastenheft 0.94.1); Spezifikation §DC-FA-HOST-001.a, `SPEC-005`-Schema-Zeilen, `SPEC-027`; `ADR-0098` (Proposed); R1-Report `2026-10-05-slice-249-hostpaths-tilde-r1.md`
- **Vorbehalt:** §7 (Closure-Notiz) leer, §6-Risiken offen, `ADR-0098` `Proposed`, README/Handbuch/CHANGELOG unberührt — laut Auftrag und Plan §1 (Release-Prep, `AGENTS.md` §5 Regel 17) der Stand **vor** der Closure. DoD 7 ist Closure-Schuld, keine DoD-Verletzung.
- **Arbeitsbaum:** alle Mutationen per `git checkout --` zurückgenommen; das Vorher-Image wurde aus `git archive 21fde701^` außerhalb des Repos gebaut und danach entfernt. `git status` nach dem Lauf: nur dieser Report.

---

## Eigene Sensor-Läufe

| Lauf | Ergebnis |
|---|---|
| `make gates` auf `0d3d10e3` | **Exit 0**, Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `coverage-gate: OK — Coverage 94.60% erfüllt Schwelle 93%`; semgrep 0 Findings |
| `make trace-check RANGE=aa854193..0d3d10e3` | 3 Commits, 0 Befunde |
| `make adr-check RANGE=aa854193..0d3d10e3` | 0 Befunde (`ADR-0098` ist `Proposed`, Körper-Änderung zulässig) |
| `make review-coverage` | 0 Befunde |

Die Implementer-Belege (Coverage 94,60 %, Gates grün) sind damit nachgefahren, nicht übernommen.

## Bewusstes Brechen (Modul 11)

Jede Zeile: genau eine Änderung am Arbeitsbaum, `make test`, Ursache gelesen, zurückgenommen.

| # | Mutation | Roter Test | Gelesene Ursache |
|---|---|---|---|
| B1 | `hostpaths.go` aus `21fde701^` (vor dem Slice) | `TestHostpathsTilde`, `TestHostpathsExemptTargets` | Tilde: nur der **abgeschnittene** Präfix-Treffer ohne Tilde, die übrigen Tilde-Pfade gar nicht (genau das Vorher-Verhalten aus Plan §1); Ventil wirkungslos (der ausgenommene Unix-Fund bleibt, der Tilde-Fund fehlt). `TestHostpathsTildeOhneUnixVerlust` bleibt grün — richtig, er behauptet „wie vor der Erweiterung" |
| B2 | `hostpaths.go` aus `21fde701` (vor der R1-Einarbeitung) | `TestHostpathsTildeOhneUnixVerlust` | `Targets = []` statt der beiden Präfix-Pfade (durchgestrichen, am Wort klebend) — exakt der R1-H-1-Verlust |
| B3 | Treffer-in-Treffer-Ausschluss aus (`hostpaths.go:56`, Bedingung immer wahr) | `TestHostpathsTilde` | Doppelbefund: abgeschnittener Präfix-Treffer **zusätzlich** zum vollen Tilde-Pfad |
| B4 | Punkt-Regel aus (`hostpaths.go:28`, `.` aus der ersten Zeichenklasse entfernt) | `TestHostpathsTilde` | das Punkt-Segment (Werkzeug-Konvention) wird gemeldet |
| B5 | Windows-Funde ventil-fähig (`hostpaths.go:62`, `false` → `true`) | `TestHostpathsExemptTargets` | der Laufwerks-Fund verschwindet unter dem Glob — Zusage „Windows/UNC fest" hängt am Test |
| B6 | `~`-Prüfung der Präfixe aus (`configyaml.go:1833`) | `TestDecode_HostpathsTildeUndExemptTargets` | `hostpaths.prefixes "~": err = <nil>` |
| B7 | Glob-Präfix-Prüfung `/`/`~` aus (`configyaml.go:1841`) | `TestDecode_HostpathsTildeUndExemptTargets` | `hostpaths.exempt-targets Library/**: err = <nil>` |
| B8 | `validateSegmentGlobs`-Aufruf aus (`configyaml.go:1837`) | `TestDecode_HostpathsTildeUndExemptTargets` | `err = <nil>` für das Glob mit ungültigem Segment `[x` (Tilde-Präfix) |

Jede DoD-Zusage mit Testbezug ist damit an einen Test gebunden, der ohne sie aus dem behaupteten Grund rot läuft.

## Black-Box gegen das gebaute Image

`make build` (`0d3d10e3`) und ein Vorher-Image (`git archive 21fde701^`), jeweils
`docker run --rm --network none -v <fixture>:/repo:ro … --enable hostpaths --disable links --disable anchors`.
Fixture-Zeilen (Fence, weil das Modul Inline-Code prüft):

```text
unix/a.md:  Pfad /home/alice/x.md und `/mnt/data/y`.  |  Link [a](/home/bob/z)
            URL https://example.org/home/x  |  C:\Users\a\x und \\srv\share\q
            Fence mit /home/fence/x  |  rel docs/home/x
tilde/b.md: ~/.claude/settings.json | ~/src/projekt/x.md, | ~/home/alice/r.md
            ~5 % und ~/ allein | ~alice/x | https://example.org/~/x
            `~/code/z` | ~~/home/alice/strike~~ | wort~/home/alice/klebt
            Fence mit ~/src/fence
```

| Prüfung (Lastenheft-Akzeptanzkriterium) | Vorher | Nachher | Urteil |
|---|---|---|---|
| Happy (relativ, Repo-Wurzel-relativ) / Boundary Fence / URL | still | still | erfüllt |
| Negative (Präfix-Pfad in Prosa und Inline-Code) | 5 Befunde, Exit 1 | identisch | erfüllt |
| Ohne `exempt-targets`: Befundsatz Unix-Korpus (`--json`) | — | `cmp` **byte-identisch** zum Vorher-Image | erfüllt (DC-QA-02) |
| `exempt-targets: []` vs. Schlüssel fehlt (`--json`) | — | byte-identisch | erfüllt |
| Home-relativ (Negative): voller Pfad samt Tilde, genau einmal | Präfix-Fall abgeschnitten ohne Tilde, Rest still | Projekt-, Präfix- und Inline-Code-Tilde je einmal, voll | erfüllt |
| Home-relativ (Boundary): Punkt-Segment, nackte Tilde, Benutzername, „ungefähr", URL, Fence | still | still | erfüllt |
| Home-relativ (kein Verlust): durchgestrichen, am Wort klebend | 2 Befunde | **dieselben 2**, gleiche Targets | erfüllt |
| Ziel-Ventil: drei Globs (Glob-Block unten) | — | genau diese drei Funde entfallen; der Tilde-Präfix-Fund hinterlässt **keinen** abgeschnittenen Unix-Rest; Windows/UNC bleiben | erfüllt |
| Config-Rand: Präfix `~` bzw. `~alice`; Glob leer, ungültig (`[x`), ohne `/`/`~` (`home/**`, `C:*`) | — | jeweils **Exit 2** mit eigener Meldung | erfüllt |
| `--doctor`-Klartext / `--print-config`-Kopf | „absolut" | „absoluter oder Home-relativer Pfad" bzw. „absolute und Home-relative Pfade", Gerüst nennt `exempt-targets` | erfüllt (R1 M-1) |

Ventil-Globs der Zeile „Ziel-Ventil" (Fence, weil das Modul Inline-Code prüft):

```yaml
hostpaths:
  exempt-targets: ["~/src/**", "/home/bob/**", "~/home/**"]
```

## DoD-Prüfung

### DoD 1 — Tilde-Erkennung (Lastenheft, Spezifikation, Kern-Regel, Tests, Rot-Beleg) — **BESTÄTIGT**

- Lastenheft `**Version:** 0.94.1`; Historie-Zeilen 0.94.0 und 0.94.1 (`spec/lastenheft.md:3992-3993`), vierspaltig — `MR-032`-Form. `DC-FA-HOST-001` Punkt 2, Ventil-Absatz, vier neue Akzeptanzkriterien plus „kein Verlust", Out-of-Scope um Benutzername-Tilde/Variablen/Windows-Ventil erweitert.
- Spezifikation §DC-FA-HOST-001.a Schritt 2 (Home-relativ, Treffer-in-Treffer-Regel), Schritt 4 (Ventil), bekannte Grenze; `SPEC-027`-Zeile erweitert; Historie-Zeile vorhanden.
- Kern-Regel `internal/hexagon/core/rules/hostpaths.go:28` (`tildeRE`), `:56` (Ausschluss); Regex deckt die Spezifikation wörtlich (erstes Zeichen weder Punkt noch Schrägstrich noch Endzeichen; Vorbedingung wie Unix plus `~`).
- Tests Happy/Negative/Boundary in `TestHostpathsTilde` inklusive aller im DoD genannten Ränder und des §6-Falls „ungefähr"; Rot-Beleg B1 zeigt den Negative-Fall **aus dem richtigen Grund** (abgeschnittener Präfix-Treffer statt voller Form).

### DoD 2 — Ziel-Ventil `hostpaths.exempt-targets` — **BESTÄTIGT**

- Schema-Zeile `spec/spezifikation.md:3308`; Modell-Feld `ExemptTargets`; Config-Rand `configyaml.go:1837-1844` (segmentweise Glob-Validierung + Präfix `/`/`~`, beides per B7/B8 gebunden); Kern-Prüfung über `exemptTarget` nur für Unix/Tilde (B5); Gerüst-Zeile im `--print-config` (Image-Ausgabe gesehen).
- Byte-Identität ohne Schlüssel: Black-Box `cmp` gegen das Vorher-Image auf einem Unix-Korpus und `[]` gegen „fehlt". Die bestehenden Tests `TestHostpathsModul`/`TestHostpathsPrefixesKonfigurierbar` sind im Diff unverändert.

### DoD 3 — Präfix mit `~` ⇒ Exit 2 mit Hinweis — **BESTÄTIGT**

`configyaml.go:1833-1835`; Meldung nennt die eigene Tilde-Erkennung und das Ventil; Exit 2 im Image für `~` und `~alice`; B6 rot.

### DoD 4 — ADR + Index — **BESTÄTIGT**

`docs/plan/adr/0098-hostpaths-home-relativ-und-ziel-ventil.md` trägt alle drei geforderten Entscheidungen (Punkt-Regel, Ventil statt Zeilen-Marker, Schärfung ohne Byte-Identität für das Muster), Alternativen-Tabelle, Fitness-Function-Tabelle (alle vier Tests existieren und sind per B1–B8 gebunden), Re-Evaluierungs-Trigger. Index-Zeile `docs/plan/adr/README.md:108`. Status `Proposed` — beabsichtigt (Übergang mit der Closure).

### DoD 5 — `make gates` grün — **BESTÄTIGT** (eigener Lauf, s. o.)

### DoD 6 — Review mit Report — **BESTÄTIGT**

`docs/reviews/2026-10-05-slice-249-hostpaths-tilde-r1.md`, anderer Kontext. Verdikt R1 war „Blockiert"; die Einarbeitung ist unten befund-weise nachgeprüft. Eine zweite Review-Runde verlangt weder DoD noch Reviewer-Skill.

### DoD 7 — Closure-Notiz, Register, Risiko-Ausgänge, Paarungen — **OFFEN (erwartet)**

§7 leer, beide §6-Risiken `(offen)` — Closure-Schuld, kein Befund. Für die Closure vorgemerkt: R1 nennt zwei Register-Klassen (`BEO-ALL/semantic-change-body-only-edges-stale` für M-1, `BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet` für H-1/L-2); Risiko 2 („ungefähr") ist durch `TestHostpathsTilde` getragen, Risiko 1 hängt an der Release-Prep (CHANGELOG/Handbuch).

## R1-Einarbeitung

| R1 | Einarbeitung | Nachgeprüft |
|---|---|---|
| H-1 (HIGH) | pauschale Tilde-Sperre im Unix-Muster zurückgenommen, Treffer-in-Treffer-Ausschluss; Test `TestHostpathsTildeOhneUnixVerlust`; Lastenheft-Kriterium „kein Verlust"; Spezifikation und `ADR-0098` Entscheidung 2 nachgezogen | **erledigt** — B2 rot gegen `21fde701`, Black-Box-Targets identisch zum Vorher-Image |
| M-1 (MEDIUM) | `--doctor`-Klartext (`diagnose.go:134`) und Gerüst-Kopf | **erledigt** — beide im Image-Output gesehen |
| L-1 (LOW) | Glob ohne `/`/`~` ⇒ Exit 2, Test-Fall `Library/**` | **erledigt** — B7 rot, Image Exit 2 |
| L-2 (LOW) | Hervorhebungs-Ränder als bekannte Grenze in der Spezifikation | **erledigt** — Spezifikation §DC-FA-HOST-001.a §Bekannte Grenze; der neue Test pinnt das Endzeichen-Verhalten |
| I-1 (INFO) | Lastenheft-Vorbedingung um „vor einem Home-relativen Pfad auch keine Tilde" präzisiert | **teilweise** — der Tilde-Teil stimmt; die pauschale Nennung „URL-, Pfad- oder Wortzeichen" gilt für die Windows-/UNC-Muster weiter nicht (Bestand vor dem Slice, R1 hat das selbst so eingeordnet) — siehe V-I-3 |
| I-2 (INFO) | Release-Prep | offen, planmäßig |

## Abgrenzung (Plan §1)

Kein Zeilen-Marker, keine Benutzername-Tilde, keine Variablen-Formen, kein Windows-/UNC-Ventil (ein Windows-Glob ist am Config-Rand sogar Exit 2), kein README/Handbuch/CHANGELOG im Diff. `diagnose.go` steht nicht in der Plan-§3-Tabelle, ist aber ein `MR-025`-Spiegel (R1 M-1) — keine Ausweitung des Gegenstands. Keine `//nolint` im Diff.

## Neue Befunde

### V-I-1 — Review-Befund-Kennung in der Lastenheft-Historie

- **Kategorie:** INFO
- **Pfad:** `spec/lastenheft.md:3992`
- **Kern:** Die Historie-Zeile 0.94.1 nennt „R1-H-1, HIGH" — eine Kennung eines Review-Reports aus `docs/reviews/`. Die fünf Kategorien von `AGENTS.md` §3.4 nennen Review-Kennungen nicht, und die Historie ist aus der Matrix ausgenommen (`MR-0098`); die zeitliche Schicht gehört aber nach `docs/plan/`. Die Vorläufer-Zeile 0.93.2 derselben Form („Nachzug nach unabhängigem Review") kommt ohne Kennung aus. Kein DoD-Verstoß; Vorschlag: die Kennung streichen, die Beschreibung trägt den Inhalt.

### V-I-2 — `ADR-0098` §Bezug nennt nur Lastenheft 0.94.0

- **Kategorie:** INFO
- **Pfad:** `docs/plan/adr/0098-hostpaths-home-relativ-und-ziel-ventil.md:10`
- **Kern:** Entscheidung 2 und das Ventil-Präfix-Kriterium der ADR entsprechen dem Stand 0.94.1, der Bezug sagt „erweitert, Lastenheft 0.94.0". Solange die ADR `Proposed` ist, lässt sich das vor dem `Accepted`-Übergang ohne Geschichte-Anhang korrigieren; danach nicht mehr.

### V-I-3 — Lastenheft-Vorbedingung pauschaler als die Windows-/UNC-Muster

- **Kategorie:** INFO
- **Pfad:** `spec/lastenheft.md:1733-1735`
- **Kern:** Rest von R1 I-1, Bestand vor dem Slice: Die Spezifikation unterscheidet korrekt („Vorbedingung hier: kein Wortzeichen davor"), das Lastenheft nicht. Kein Akzeptanzkriterium hängt daran.


### V-I-4 — Ein Ventil-Glob in Prosa oder Inline-Code ist selbst ein Befund

- **Kategorie:** INFO
- **Pfad:** `docs/reviews/2026-10-05-slice-249-hostpaths-tilde-verify.md` (dieser Report, erster Entwurf)
- **Kern:** Gemessen beim eigenen `make doc-check`: Globs wie die drei aus dem Ventil-Block oben werden in Inline-Code als `hostpath-forbidden` gemeldet (Tilde- und Präfix-Globs erfüllen das Muster; `*` ist kein Endzeichen). Folgerichtig nach Spezifikation — Inline-Code wird geprüft —, aber jedes Dokument, das ein `hostpaths.exempt-targets`-Beispiel nennt, muss es in einen Fence setzen. Für die Release-Prep (Benutzerhandbuch §5/§6, CHANGELOG) relevant; das `--print-config`-Gerüst ist Go-Quelltext und nicht betroffen.
## Verdikt

**DoD 1–6 bestätigt, DoD 7 erwartet offen.** Alle testgestützten Zusagen sind per Bewusstem Brechen gebunden (B1–B8), die Akzeptanzkriterien des Lastenhefts sind am gebauten Image erfüllt, die Byte-Identität ohne Schlüssel ist per `cmp` gegen das Vorher-Image belegt, die R1-Einarbeitung trägt. Keine HIGH-/MEDIUM-/LOW-Befunde; vier INFO, davon V-I-2 vor dem `Accepted`-Übergang sinnvoll. Freigabe für die Closure.
