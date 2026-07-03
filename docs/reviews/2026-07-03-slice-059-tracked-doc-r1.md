# Review-Report — slice-059 (Modul `tracked`) — Doku-Kette & Spec-Treue (R1)

**Datum:** 2026-07-03
**Reviewer-Rolle:** unabhängig/adversarial, Fokus **Doc-Straten + Cross-Surface-Konsistenz
+ Doku↔Code-Ehrlichkeit** (NICHT Go-Code-Qualität — separater zweiter Reviewer).
**Gegenstand:** neues opt-in Modul `tracked` (16.) — Doku-Kette
[Lastenheft §DC-FA-TRK-001](../../spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in)
(+ Schema-Konvention `TRK`, CLI-002/-010, Glossar, §7 0.37.0),
[Spezifikation §DC-FA-TRK-001.a](../../spec/spezifikation.md#dc-fa-trk-001a--getrackt-status-auflösbarer-referenz-ziele-tracked)
(+ §2 `tracked.exempt-targets`, §4 `target-untracked`, CLI-010.a, §7),
[ADR-0030](../plan/adr/0030-tracked-referenz-ziele.md) (Proposed) + ADR-Index,
[slice-059](../plan/planning/in-progress/slice-059-tracked-modul.md) + Roadmap (welle-48),
[Benutzerhandbuch](../user/benutzerhandbuch.md) §4.13/§5/§6,
[CHANGELOG](../../CHANGELOG.md) `[Unreleased]` — jeweils gegen die Implementierung
(feat-Commit `d39dd3f`).
**Baseline:** `.harness/skills/reviewer.md` (Kategorien/Schema/Negativbefund-Pflicht),
`AGENTS.md` §3.
**NICHT geprüft:** Code-Qualität/Testarchitektur (R2), DoD-Abhakung (Verifikation),
ids-/matrix-Konformität (gate-geprüft).

**Verifikations-Läufe (lokales Image `d-check:latest`, Build-Stand des feat-Commits;
`--network none`, read-only; git-Fixture-Repos im Scratch):**

- `--print-mk | grep -cE '^doc-[a-z-]+:'` = **10**; `doc-tracked`-Recipe: `--enable
  tracked` + fokussierte `--disable`-Liste (inkl. `--disable links`), **ohne** Range.
- **Probe 1** (Fixture: `docs/a.md` verlinkt `sub/u.md`; Ziel existiert, untracked):
  `--enable tracked --disable links --disable anchors` ⇒ `docs/a.md:3  sub/u.md
  target-untracked`, Exit 1 — das Modul feuert **ohne** aktives `links`.
- **Probe 2** (Ventil): `exempt-targets: ["sub/u.md"]` (die Schreibweise aus der
  Befund-Spalte) ⇒ Befund **bleibt**; `exempt-targets: ["docs/sub/u.md"]`
  (aufgelöster Pfad) ⇒ 0 Befunde.
- **Probe 3** (Verzeichnis-Ziel): Link auf ein existierendes, komplett gitignoriertes
  **Verzeichnis** ⇒ mit `--enable tracked` **0 Befunde**, Exit 0; frischer Klon
  desselben Repos ⇒ `target-missing`, Exit 1.
- **Probe 4** (`--doctor --enable tracked`): Befundzeile zeigt den **nackten Code**
  `target-untracked` (kein Klartext).
- **Probe 5** (Config-Rand): `exempt-targets: [""]` ⇒ `error: .d-check.yml:
  tracked.exempt-targets enthält ein leeres Glob`, Exit 2 (config-zeitig).

---

## Findings (MEDIUM)

### MEDIUM-1 — §DC-FA-TRK-001.a beschreibt eine Mechanik („Post-Pass über die von `links` aufgelösten Ziele"), die es nicht gibt

- **kategorie:** MEDIUM
- **quelle:** DC-FA-TRK-001 (Spec-Treue der Messmethode)
- **pfad:** `spec/spezifikation.md:1130` (+ `:1143`), `docs/plan/adr/0030-tracked-referenz-ziele.md:41`,
  `docs/plan/planning/in-progress/slice-059-tracked-modul.md:44` (+ DoD `:77`),
  Roadmap `docs/plan/planning/in-progress/roadmap.md:14`
- **befund:** Die `.a`-Sektion eröffnet mit „Opt-in-Post-Pass" und Schritt 3 nennt als
  Kandidaten „alle von `links` erfolgreich aufgelösten … Ziele"; ADR-0030 („prüft als
  Post-Pass") und der Slice tragen dieselbe Formulierung. Der Code prüft jedoch **je
  Quell-Datei im Scan** (`run.go` `checkFile`, `rules/tracked.go`) mit einmal je Lauf
  geladener Index-Menge und **eigener** Extraktion/Auflösung (`ExtractLinks`/`localTarget`)
  — unabhängig davon, ob das Modul `links` aktiv ist. „Post-Pass" ist im selben Dokument
  und im Code der Fachbegriff für die **nach** dem Datei-Scan laufenden Module (`run.go`
  Kommentar „Post-Pässe": vcs/commits/planning; `runPostPasses` enthält `tracked` nicht).
  Die Sektion ist zudem in sich gespalten: Schritt 3 verlangt selbst den Modul-Scope
  „für die Quell-Dateien wie bei jedem Modul" — das passt zur realen per-Datei-Prüfung,
  nicht zu einem Pass über `links`-Ergebnisse.
- **Failure-Szenario:** Ein Verifizierer, der die `.a`-Sektion wörtlich nimmt, sagt für
  das in §DC-FA-CLI-010.a spezifizierte `doc-tracked`-Target (das `links` per
  Fokus-`--disable` abschaltet) einen **No-op** voraus — real liefert genau dieser
  Aufruf einen Befund (Probe 1). Die bindende Algorithmus-Beschreibung und der
  spezifizierte Verteilungs-Aufruf widersprechen einander; das Target funktioniert nur,
  **weil** die `.a`-Beschreibung nicht stimmt.
- **verifizierbar:** ja (Probe 1: `--enable tracked --disable links` ⇒ 1×
  `target-untracked`, Exit 1).

### MEDIUM-2 — Befund-Feld `target`: Spec sagt „aufgelöster Zielpfad", Code liefert die rohe Link-Schreibweise

- **kategorie:** MEDIUM
- **quelle:** DC-FA-TRK-001 (Negative-AK) / §DC-FA-TRK-001.a Schritt 5
- **pfad:** `spec/spezifikation.md:1157`, `spec/lastenheft.md:1378` vs.
  `internal/hexagon/core/rules/tracked.go:40`
- **befund:** Schritt 5 legt fest „`target` = der **aufgelöste** Zielpfad"; das
  Lastenheft-Negative-AK verlangt „(Datei, Zeile, **aufgelöstes Ziel**)". Der Code setzt
  `Target: ref.Target` — die **rohe** Schreibweise aus dem Dokument (inkl. relativer
  `../`-/Unterpfad-Form und ggf. Fragment). Probe 1 zeigt für das in `docs/sub/u.md`
  aufgelöste Ziel die Befund-Spalte `sub/u.md`. Die Unit-Tests fixieren nur Wurzel-Fälle,
  in denen roh == aufgelöst gilt — die Abweichung ist unverriegelt. Verschärfend:
  `tracked.exempt-targets` matcht laut §2 (und real, Probe 2) den **aufgelösten** Pfad —
  wer die Befund-Schreibweise ins Ventil kopiert, schreibt ein wirkungsloses Glob. Genau
  diese Roh-vs-aufgelöst-Falle dokumentiert das Handbuch für `codepaths.ignore-refs`
  explizit; der neue `tracked`-Absatz (§5) trägt den Hinweis nicht.
- **Failure-Szenario:** (a) Ein aus der Spec abgeleiteter Akzeptanztest
  (`target == "docs/sub/u.md"`) schlägt gegen die Implementierung fehl. (b) Ein Nutzer
  überträgt `sub/u.md` aus der Befund-Spalte nach `exempt-targets` — der Befund bleibt
  bestehen (Probe 2a), ohne dass irgendeine Doku-Schicht die Diskrepanz erklärt.
- **verifizierbar:** ja (Proben 1 + 2a/2b).

### MEDIUM-3 — Verzeichnis-Ziele sind still kein Kandidat — in keiner Doku-Schicht gesagt, Versprechens-Loch offen

- **kategorie:** MEDIUM
- **quelle:** DC-FA-TRK-001 (Beschreibung/Out-of-Scope) / §DC-FA-TRK-001.a Schritt 3
- **pfad:** `internal/hexagon/core/rules/tracked.go:28` (Skip `kind != KindFile`) vs.
  `spec/lastenheft.md:1331`–`1339` + `:1383`–`1391` (Out-of-Scope),
  `spec/spezifikation.md:1143`, `docs/user/benutzerhandbuch.md:813`
- **befund:** Der Code überspringt Ziele, die auf ein **Verzeichnis** auflösen
  (nur der Code-Kommentar und die Commit-Message benennen das). Das Lastenheft verspricht
  die Prüfung für „**jedes** auflösbare, existierende repo-interne Link-/Bild-Ziel" und
  führt Verzeichnis-Ziele **nicht** in der Out-of-Scope-Liste (die Submodule, gescannte
  Dateien selbst etc. explizit ausnimmt); die `.a`-Sektion sagt nur „Datei-Ziele", ohne
  die Ausschluss-Regel und ihre Konsequenz zu benennen; Handbuch/ADR/CHANGELOG schweigen.
  Das Loch ist real: `links` akzeptiert existierende Verzeichnis-Ziele, und ein Link auf
  ein existierendes, komplett untracktes/gitignoriertes Verzeichnis bleibt mit aktivem
  `tracked` grün, ist auf dem frischen Klon aber `target-missing` (Probe 3) — exakt die
  Umgebungs-Drift, deren Fang die Anforderung als Zweck deklariert.
- **Failure-Szenario:** Ein Adopter verlinkt einen lokal generierten Ordner (z. B.
  `build/`-Report-Verzeichnis), aktiviert `tracked` als Schutz vor der Frischer-Klon-Falle
  und wähnt sich gedeckt — der nächste Checkout ist rot, obwohl das Modul genau dieses
  Szenario zu fangen verspricht. Aus keiner Doku-Schicht ist ableitbar, dass hier eine
  bewusste Grenze liegt (weder als Verhalten noch als Out-of-Scope/Re-Eval-Trigger).
- **verifizierbar:** ja (Probe 3a/3b: Erzeuger Exit 0 mit `tracked`, Klon
  `target-missing`).

### MEDIUM-4 — DC-FA-CLI-010: Happy-AK, Boundary-AK und Out-of-Scope stehen auf dem Neun-Stand — Out-of-Scope widerspricht der eigenen Beschreibung

- **kategorie:** MEDIUM
- **quelle:** DC-FA-CLI-010 (Akzeptanzkriterien/Out-of-Scope)
- **pfad:** `spec/lastenheft.md:415` (Happy), `:416` (Boundary), `:419` (Out-of-Scope)
- **befund:** Die Beschreibung ist auf „**zehn** Targets" gezogen und enumeriert
  `doc-tracked` (Zeilen 379, 398–402). Alle drei nachgelagerten Anker sind stale: das
  Happy-AK endet bei `doc-planning` + `doc-help` ohne `doc-tracked`; das Boundary-AK
  enumeriert die Aufruf-Modi nur bis `doc-planning` (kein `doc-tracked → --enable
  tracked + Fokus-… ohne Range`); das Out-of-Scope schließt „weitere Targets jenseits
  der gelisteten **neun**" aus und listet neun ohne `doc-tracked` — nach dieser Klausel
  wäre das real ausgelieferte zehnte Target ausdrücklich out-of-scope. Der Vorgänger
  slice-057 ließ nur das Boundary-AK aus (R1-LOW); hier fehlt das neue Target in **allen
  drei** Ankern, und einer widerspricht aktiv der Beschreibung und dem Fragment
  (`--print-mk` liefert real 10, siehe Verifikations-Läufe).
- **Failure-Szenario:** Ein Verifizierer, der die AK als Checkliste nimmt, prüft
  `doc-tracked` nicht (Happy/Boundary) bzw. muss es nach der Out-of-Scope-Klausel als
  vertragswidriges Extra-Target werten — dieselbe Spec-Sektion sagt beides zugleich.
- **verifizierbar:** teils (`--print-mk`-Zählung ja; die AK-Prosa selbst prüft kein Gate).

### MEDIUM-5 — Neuer §4-Grund-Code ohne `--doctor`-Klartext — §DC-FA-CLI-007.a-Zusage („für jeden Grund-Code aus §4") greift nicht

- **kategorie:** MEDIUM
- **quelle:** DC-FA-CLI-007 (§DC-FA-CLI-007.a Schritt 3) / Harness-Ehrlichkeit
- **pfad:** `spec/spezifikation.md:225` (Zusage) + `:1414` (neue §4-Zeile) vs.
  `internal/hexagon/core/app/diagnose.go:67`–`99`
- **befund:** §CLI-007.a Schritt 3 behauptet: Klartext-Mapping „für jeden Grund-Code aus
  §4 genau ein Eintrag, abgesichert durch eine Vollständigkeits-Prüfung gegen die
  Reason-Konstanten". `AllReasons()`/`reasonTexts()` enden jedoch bei
  `hostpath-forbidden`; der neue Code `target-untracked` (und die Bestandsklasse
  `diagram-id-undefined`/`version-stale`/`link-stale`/`core-drift`/`core-drift-vcs`/
  `commit-untraceable`/`planning-drift` — Drift seit v0.25) fehlt in beiden. `--doctor`
  zeigt für den neuen Befund den nackten Code statt eines Klartexts (Probe 4); die
  „Vollständigkeits-Prüfung" (diagnose_test) vergleicht Mapping ↔ `AllReasons` und kann
  die Lücke prinzipbedingt nicht fangen, weil die handgepflegte `AllReasons`-Liste selbst
  nicht gegen §4 wächst. slice-059 verlängert die Reihe um den 8. Code; der
  Bestandsanteil ist nicht slice-verursacht, die Spec-Zusage ist aber für die neue
  §4-Zeile erneut falsch.
- **Failure-Szenario:** Ein Nutzer folgt Handbuch §4.9 („erklärende Diagnose") für einen
  `target-untracked`-Befund und erhält als „Klartext" den Code selbst; ein Verifizierer,
  der §CLI-007.a liest, hält das per Spec für unmöglich („abgesichert durch
  Vollständigkeits-Prüfung").
- **verifizierbar:** ja (Probe 4).

---

## Findings (LOW)

### LOW-1 — Handbuch behauptet die 0.37-Surface für das v0.36.0-Kommando (transient bis Release-Prep)

- **kategorie:** LOW
- **quelle:** Maintainability (Doku-Currency; bekannter Release-Prep-Blind-Spot)
- **pfad:** `docs/user/benutzerhandbuch.md:615`/`:627` (Pin v0.36.0) vs. `:621`–`624`
  („zehn … `doc-tracked`"); `:954` ff. (§11 koppelt Handbuch 1.18 an v0.36.0, keine
  1.19/0.37-Zeile)
- **befund:** §4.13 zeigt das Kommando mit dem `v0.36.0`-Pin und behauptet als dessen
  Ergebnis „zehn `##`-annotierte Targets (… `doc-tracked` …)" — das v0.36.0-Image liefert
  neun und kein `doc-tracked`. Ebenso dokumentieren §5/§6 das `tracked`-Modul unter der
  §11-Kopplung „Handbuch 1.18 / v0.36.0". Der Zustand heilt beim Release-Prep-Pin-Bump
  mechanisch (der `versions`-Gate fängt die `ghcr`-Pins beim `version.md`-Flip; die
  §11-Historien-Zeile bleibt der bekannte manuelle Schritt) — bis dahin macht das
  Handbuch auf `main` eine falsche Aussage über ein released Artefakt.
- **Failure-Szenario:** Ein Leser zieht heute per §4.13-Kommando das Fragment und findet
  weder zehn Targets noch `doc-tracked`.
- **verifizierbar:** ja (v0.36.0-Image ausführen); der Pin-Anteil zusätzlich durch den
  `versions`-Gate nach dem `version.md`-Flip.

### LOW-2 — DC-FA-CLI-006: AK- und Out-of-Scope-Enumeration der situativen Module enden bei `vcs`

- **kategorie:** LOW
- **quelle:** DC-FA-CLI-006 (AK „ai-harness Auffindbarkeit", Out-of-Scope)
- **pfad:** `spec/lastenheft.md:263` + `:265` vs. `internal/hexagon/core/app/suggest.go:392`
- **befund:** Beide Enumerationen nennen „(`external`, `spans`, `hostpaths`, `diagrams`,
  `versions`, `pins`, `immutable`, `vcs`)"; die reale `--suggest-config`-Ausgabe nennt
  seit diesem Commit zusätzlich `commits`, `planning`, `tracked`. Bestands-Drift (fehlte
  schon für `commits`/`planning`, slice-056/057); slice-059 verlängert sie um `tracked`,
  ohne die bindende AK-Liste nachzuziehen.
- **Failure-Szenario:** Ein Verifizierer, der den Kommentar-Inhalt exakt gegen die
  AK-Liste prüft, wertet drei Modul-Nennungen als nicht spezifiziert.
- **verifizierbar:** ja (`--suggest-config ai-harness` ausführen, Kommentar vergleichen).

### LOW-3 — Config-zeitige Glob-Validierung von `tracked.exempt-targets` (leer/ungültig ⇒ Exit 2) in keiner Schicht dokumentiert

- **kategorie:** LOW
- **quelle:** DC-FA-TRK-001 / DC-FA-CONF-001 (Konsistenz zur `planning.slice-glob`-Doku)
- **pfad:** `spec/spezifikation.md:1365` (§2-Zeile ohne Validierungs-Klausel) vs.
  `internal/adapter/driven/configyaml/configyaml.go:371`–`388`
- **befund:** `applyTracked` weist leere und `path.Match`-ungültige Globs config-zeitig
  mit Exit 2 ab (Probe 5; laut Commit-Message mutations-verriegelt). Keine Doku-Schicht
  sagt das: die §2-Zeile schreibt nur „Glob (wie `scan.ignore`)" — und `scan.ignore` ist
  gerade **nicht** config-zeitig validiert, die Analogie suggeriert also das Gegenteil.
  Der direkte Vorgänger dokumentiert denselben Guard explizit in seiner §2-Zeile
  (`planning.slice-glob`: „muss ein gültiges `path.Match`-Glob sein (sonst Exit 2 …)",
  Zeile 1364). Fail-closed-Richtung (strenger als dokumentiert, kein stilles Grün),
  daher LOW.
- **Failure-Szenario:** Ein aus der Spec ableitender Verifizierer kann den
  mutations-verriegelten Guard nicht herleiten; ein Nutzer mit leerem Glob-Eintrag
  bekommt einen von keiner Doku vorhergesagten Exit 2.
- **verifizierbar:** ja (Probe 5).

---

## Findings (INFO)

### INFO-1 — §4.13-Beispiel-Exzerpt elidiert vier der zehn Targets ohne Elisions-Marker

- **kategorie:** INFO
- **quelle:** Maintainability (Bestand, wächst je Modul weiter)
- **pfad:** `docs/user/benutzerhandbuch.md:626`–`651`
- **befund:** Das `text`-Exzerpt zeigt drei Recipes und benennt in den „…"-Kommentaren
  nur `doc-trace`/`doc-complete`/`doc-help`; `doc-immutable`/`doc-commits`/
  `doc-planning`/`doc-tracked` kommen im Exzerpt gar nicht vor. Bestands-Stil (nicht
  slice-059-verursacht), driftet aber mit jeder Modul-Welle weiter vom „zehn"-Anspruch
  der umgebenden Prosa weg.
- **verifizierbar:** nein (illustratives Exzerpt, kein Gate).

---

## Negativbefunde (geprüft, ohne Befund)

- **Lastenheft-Gerüst:** Schema-Konvention führt `TRK` (Z. 50); `DC-FA-CLI-002`-Modulliste,
  Glossar-Regelmodul-Zeile und §7-Historie (0.37.0, Version-Header 0.37.0) führen
  `tracked` konsistent; die 0.37.0-Zeile deckt sich inhaltlich mit ADR/Slice/CHANGELOG
  (Index-Wahrheit, dritte Port-Nutzung ohne Range, kein Doppelbefund, Ventil,
  fail-closed, 9→10).
- **Spezifikation §2/§4/§7 + CLI-010.a:** §2-Zeile `tracked.exempt-targets`
  (referenz-weit, aufgelöste Ziel-Pfade, ohne Eintrag byte-identisch), §4-Zeile
  `target-untracked`, §7-Historien-Zeile vorhanden und untereinander konsistent;
  §CLI-010.a („Zehn `.PHONY`-Targets", `doc-tracked`-Bullet ohne `--range`/`--staged`,
  abgeleitete `--disable`-Liste) deckt sich exakt mit `print_mk.go` (FÜNF fmt-Verben)
  und der realen Fragment-Ausgabe (10 Targets, Recipe siehe Verifikations-Läufe).
- **Gestagt = getrackt:** in allen Straten gleichlautend (Lastenheft-Boundary-AK,
  `.a` Schritt 2, ADR, Handbuch, CHANGELOG); der Adapter liest den go-git-Index
  (`Storer.Index()` — gestagte, nie committete Dateien enthalten). Deckungsgleich.
- **Escaped-Ziele:** alle Straten begrenzen auf „repo-interne" Ziele; der Code
  überspringt escapte Auflösungen (bleiben `repo-escape` bei `links`) — vom
  Unit-Test (`../raus.md`) fixiert. Kein Widerspruch.
- **Kein Doppelbefund:** nicht existierende Ziele bleiben `target-missing`
  (`links`), `tracked` meldet nur existierende untrackte Datei-Ziele — Code
  (`Kind`-Check), AK und alle Prosa-Schichten deckungsgleich; Probe 3b zeigt das
  `links`-Verhalten unangetastet.
- **Fail-closed `.git`:** Lastenheft/`.a`/ADR/Handbuch sagen übereinstimmend „aktiv ohne
  lesbares `.git` ⇒ Exit 2, kein stilles Grün"; `RunWithVCS` guardet `vcs == nil` und
  `TrackedPaths`-Fehler als error → Exit 2; die CLI öffnet den Adapter für `tracked`
  ohne Range-Pflicht (`resolveVCS`), `vcs`/`commits` verlangen die Range weiter. Das
  slice-Risiko „Konsumenten ohne git-Repo" ist im Handbuch adressiert (Voraussetzung
  explizit benannt).
- **Config-Surface-Currency (bis auf LOW-2/LOW-3):** `--print-config`-Template
  (Verfügbar-Liste + eigener `tracked`-Block mit opt-in-/fail-closed-Hinweis),
  `--suggest-config`-Kommentar, `validModules()` (16 Einträge), Handbuch §5
  (Überblicks-YAML + eigener Absatz + `doc-tracked`) und §6-Tabellenzeile führen
  `tracked` konsistent; „16. Regelmodul" stimmt über alle Surfaces.
- **ADR-0030 / ADR-Index:** Status Proposed beidseitig; `Schärft:` zeigt aufwärts auf
  die `.a`-Sektion, `Bezug:` aufs Lastenheft; Alternativen-Tabelle und
  Re-Evaluierungs-Trigger decken die Lastenheft-Out-of-Scope-Punkte (Inline-Code-Ziele,
  Gegenrichtung, Submodule); Geschichte-Eintrag ist Entstehungs-Provenance, keine
  getarnte Entscheidungsgrundlage. Kein unmarkierter Abwärts-Token im Körper.
- **Roadmap ↔ Slice:** welle-48 als aktive Welle benannt, slice-059 in `in-progress/`,
  Status-Zeile und Bezugs-Block des Slice konsistent (Lastenheft 0.37.0, CLI-010 9→10,
  ADR-0030 Proposed, Release v0.37.0 erwartet).
- **CHANGELOG `[Unreleased]`:** Aussagen decken sich mit dem Code; keine
  Phantom-Behauptung — das Handbuch ist diesmal real nachgezogen (§4.13/§5/§6), im
  Unterschied zum slice-057-R1-Befund (dort MEDIUM). Abschnitt korrekt als
  `[Unreleased]` statt vorgezogener Versions-Überschrift.

---

## Kategorie-Summary

| Kategorie | Anzahl |
| --- | --- |
| HIGH | 0 |
| MEDIUM | 5 |
| LOW | 3 |
| INFO | 1 |

## Verdikt

**NACHBESSERN.** Die Kette ist über die üblichen Currency-Flächen (Schema-Konvention,
Modul-Listen, Glossar, Historien, print-config/suggest, Handbuch §5/§6, CHANGELOG)
sauber und in den Kern-Zusagen (Index-Wahrheit, kein Doppelbefund, fail-closed,
ohne Range) code-deckungsgleich. Blockierend sind die fünf MEDIUMs: die bindende
`.a`-Sektion beschreibt eine andere Mechanik als gebaut und kollidiert wörtlich mit dem
eigenen `doc-tracked`-Target (MEDIUM-1); das spezifizierte Befund-Feld `target` weicht
beobachtbar vom Code ab und untergräbt die Ventil-Bedienung (MEDIUM-2 — hier muss
entschieden werden, ob Spec oder Code sich bewegt); der stille Verzeichnis-Ausschluss
lässt das zentrale Versprechen der Anforderung undokumentiert löchrig (MEDIUM-3);
DC-FA-CLI-010 widerspricht sich selbst zwischen Beschreibung und AK/Out-of-Scope
(MEDIUM-4); und der neue Grund-Code fällt in eine falsche §CLI-007.a-Vollständigkeits-
Zusage (MEDIUM-5, Bestandsklasse — der Alt-Anteil ist als eigener Fix/CR eskalierbar,
nicht slice-059 anzulasten). Die LOWs sind Currency-Nachzieher (LOW-1 heilt beim
Release-Prep-Pin-Bump mechanisch, LOW-2/LOW-3 sind Enumerations-/Dokumentations-Drift
mit klarer Vorgänger-Präzedenz).
