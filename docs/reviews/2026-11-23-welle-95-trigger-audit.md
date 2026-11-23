# Trigger-Audit welle-95 — Zweige ADR und Hard Rule (Architect)

**Rolle:** Architect (Modul 8, §Rollen-Sequenz für eine Welle, Closure-Schritt 2: Planner → Architect → Planner).
**Gegenstand:** [welle-95](../plan/planning/done/welle-95/welle-95-linelength-modul.md) mit
[slice-289](../plan/planning/done/welle-95/slice-289-linelength-anforderung-und-adr.md) und
[slice-290](../plan/planning/done/welle-95/slice-290-linelength-modulkern.md).
**Eingang:** Audit-Vorlage des Planners; Verifikations-Beleg der Welle
([`2026-11-23-welle-95-verifikation.md`](2026-11-23-welle-95-verifikation.md): `make fullbuild` grün,
`blackbox-probe` byte-identisch gegen `v0.91.0`, `image-test` Phase 9).
**Nicht Gegenstand:** die Zweige Carveout und bootstrap-aware Gate (prüft der Planner selbst).
**Datum:** 2026-11-23.
**Geändert wurde nichts:** keine ADR, keine Hard Rule; dies ist nur das Verdikt.

## Gesamturteil

Kein Re-Evaluierungs-Trigger und kein Auflösungs-Trigger ist eingetreten. Es ist **keine Folge-ADR**
nötig, und **keine Hard-Rule-Zeile** wird entfernt. Die beiden Beobachtungen aus slice-290 berühren
keine Entscheidung, die eine ADR trägt. Eine davon (`regeltext-gegen-spiegel-test`) widerspricht aber
einem Regeltext, der niedriger rangiert als das Gate, und gehört als kleiner Nachzug an den Regeltext
(Punkt 1b) — ohne ADR.

Bestandsaufnahme der Gate-Lage über die Welle (`ca85c4c2..HEAD`): `.d-check.yml`, `Makefile`,
`.golangci.yml` und `tools/coverage-gate.sh` sind unverändert
(`git diff --stat ca85c4c2..HEAD -- .d-check.yml Makefile .golangci.yml tools/coverage-gate.sh` ⇒ leer),
`internal/hexagon/core/rules/markdown.go` ebenso.

## 1. ADR-0114 (`linelength`)

| Trigger | eingetreten? | Beleg |
|---|---|---|
| (a) die geteilte Lexik ändert ihre Fence- oder Tabellen-Antwort | **nein** | `markdown.go` ist über die Welle unverändert. Die Kopplungs-Tests aus slice-290 fahren `linelength` und die bisherigen Konsumenten durch dieselbe Eingabe; ab jetzt wird der Trigger rot sichtbar, statt still einzutreten |
| (b) ein Konsument muss eine Zeile ausnehmen, die weder Code noch Tabelle noch unumbrechbar ist | **nein** | kein Konsument; das eigene Profil schaltet das Modul nicht ein (Welle-Out-of-Scope) |
| (c) die Grobheit von (c) unter Variante H wird als Mangel gelesen, oder umgekehrt der Umbruch vor einem langen Link als Last | **nein** | die Reviews R1/R2 zu slice-289 betrafen die **Begründung** von (c), nicht die Regel; ADR-0114 ist danach berichtigt. Code-Reviews und Verifikation zu slice-290 melden zu (c) keinen Mangel. H ist die vom Auftraggeber am 2026-11-23 neu bestätigte Fassung, keine nachträgliche Lesart |
| (d) breite Zeichen, Codepunkte tragen nicht | **nein** | kein Konsument mit solchem Text |
| (e) eine andere Dateiart wird verlangt | **nein** | kein Antrag |
| (f) eine weitere Lint-, Stil- oder Rechtschreib-Prüfung wird gewünscht | **nein** | die Ausnahme in Lastenheft §5 deckt nur die Zeilenlänge; ein weiterer Antrag liegt nicht vor |

**Die zwei Beobachtungen aus slice-290 — kein Trigger für ADR-0114:**

- **1a — [`BEO-ALL/yaml-dekoder-kuerzt-zahlen-still`](../plan/planning/observations/BEO-ALL/yaml-dekoder-kuerzt-zahlen-still/observation.md)** (1×, offen).
  ADR-0114 Entscheid 4 schließt die Lücke für `linelength` (nur YAML-Ganzzahl, nur `true`/`false`); die
  Beobachtung betrifft die **Nachbar**-Module (`file[].max-lines`, `cell-max-chars`, Wahrheitswerte anderer
  Blöcke). Sie widerspricht keiner Entscheidung: Weder ADR-0088 noch ADR-0069 hat die Dekodierung
  entschieden, und eine still gekürzte Kommazahl verletzt die vollständige Validierung aus
  [`DC-FA-CONF-001`](../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei). Ihr Weg ist deshalb ein
  **Defekt-Fix** am Konfigurations-Rand, kein Folge-ADR. Bis dahin trägt sie das Register; kein eigener
  Träger aus diesem Audit.
- **1b — [`BEO-ALL/regeltext-gegen-spiegel-test`](../plan/planning/observations/BEO-ALL/regeltext-gegen-spiegel-test/observation.md)** (1×, offen).
  Der Widerspruch besteht zwischen [Regel 17](../../harness/rules/dokumentations-regeln/17-changelog-release-prep.md)
  („Dieselbe Grenze gilt den beiden `README*.md` … Ein fehlender Eintrag dort im Feature-Commit ist deshalb
  kein Rückstand“) und dem Spiegel-Test `registry_mirror_test.go`, der die Modul-Bullets im Feature-Commit
  verlangt. ADR-0114 steht auf der Seite des Gates (Folgepflicht: das Modul wird zusammen mit der
  Modul-Aufzählung registriert) — **kein Trigger**. Nach Source Precedence wird die niedriger rangierte Stelle
  angepasst ([`AGENTS.md`](../../AGENTS.md) §1), und das ist der Regeltext, nicht das Gate. **Der Gegenstand
  ist älter als die Welle:** Dieselbe Ausnahme musste für `metrics` (Architect-Verdikt welle-93, F-1) und für
  `pinsync` (slice-289 §1: „bei `pinsync` gemessen“) im Plan neu begründet werden. Das Register führt bisher
  nur den Beleg aus slice-290. **Vorschlag an den Planner** (Closure-Schritt 3a): die beiden früheren Fälle
  als Belege prüfen. Erreicht der Eintrag damit die Schwelle, ist die Verkörperung ein Satz im Volltext von
  Regel 17: „Ausgenommen sind Modul-Aufzählungen, die ein Spiegel-Test gegen die Modul-Registrierung hält;
  sie kommen mit dem Feature-Commit, der das Modul registriert.“ Das ist eine Text-Angleichung an ein
  bestehendes Gate, keine Gate-Änderung und kein ADR-Anlass (§3.6 betrifft sie nicht).

## 2. Berührte und zitierte ADRs

| ADR | Trigger | eingetreten? | Beleg |
|---|---|---|---|
| [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) | ein Fall der Klasse tritt im Bestand ein | **nein** | `linelength` übernimmt die geteilte Antwort (`proseLineSet`, `tableRowLine`) und ist über Kopplungs-Tests daran festgemacht — die Form aus Entscheidung 4, kein Fall der Klasse |
| ADR-0054 | `spans` meldet `fence-unclosed` in einer Datei, die `vcs.paths` trifft | **nein** | `spans` ist im eigenen Profil aktiv; `make doc-check` an `HEAD` meldet 0 Befunde |
| ADR-0054 | eine **vierte** Stelle beantwortet eine Lexik-Frage selbst | **nein** | die Welle hat Stellen **zusammengelegt**, nicht hinzugefügt: `pinsync` liest seine Zeilen jetzt über die geteilte `textLines` (`internal/hexagon/core/rules/codepaths.go`), seine eigene `pinsyncLines` ist entfernt. **Akzeptiertes Negativ:** `metrics_table.go` schneidet `\r` weiter selbst ab — Bestand vor der Welle, mit derselben Antwort; kein eigener Träger, er fällt beim nächsten Eingriff in `metrics` mit |
| [ADR-0113](../plan/adr/0113-pinsync-modul.md) | (a) mehrzeiliger Wert · (b) Normalisierung · (c) Quellen über Symlink | **nein** | kein Konsument mit solchem Bedarf. Der Wechsel auf `textLines` ist eine reine Verschiebung: der entfernte und der geteilte Rumpf sind Zeile für Zeile gleich (`git diff ca85c4c2..HEAD -- internal/hexagon/core/rules/pinsync.go`); Entscheidung 2 (`\n`, ein `\r` am Zeilenende gehört nicht zur Zeile) gilt unverändert |
| [ADR-0088](../plan/adr/0088-file-modul-groessengrenzen.md) | (1) ein Folge-Slice setzt eine Schwelle für eine Datei · (2) `structure` liest rohen Text | **nein** | die Welle setzt keine Schwelle; die vorhandenen `file`-Regeln im eigenen Profil stammen aus slice-239 und sind über die Welle unverändert. `linelength` liest rohen Text, ist aber ein eigenes Modul — `structure` hat keine neue Fähigkeit |
| [ADR-0089](../plan/adr/0089-structure-max-lines-zwoelfte-bedingung.md) | (1) Schwelle für eine konkrete Datei · (2) wiederholte Verwechslung `structure[].max-lines`/`file[].max-lines` | **nein** | keine Schwelle gesetzt; [`BEO-ALL/max-lines-namenskollision-structure-file`](../plan/planning/observations/BEO-ALL/max-lines-namenskollision-structure-file/observation.md) hat keinen neuen Beleg (nur `slice-237`). ADR-0114 meidet einen dritten nahen Namen (`max-chars`) und grenzt die fünf Längen-Achsen in einer Tabelle ab |

## 3. Hard Rules in `AGENTS.md` §3

| Regel | Auflösungs-Trigger | eingetreten? |
|---|---|---|
| §3.1 Docker/make-only | keiner benannt (Bestand, kein Nachrüsten) | — |
| §3.2 Suppression-Verbot | permanent | nein |
| §3.3 `git mv` + Inhalt = zwei Commits | permanent | nein |
| §3.4 Architektur sprachfrei, Straten nie abwärts | Closure-Datum: keiner; Rollen-statt-Technologie: permanent | nein — `architecture.md` ist unberührt |
| §3.5 ADRs nach `Accepted` immutable | durchgesetzt (`make adr-check`) | — ADR-0114 wurde mit `6b56bd57` `Accepted`, `make adr-check` grün |
| §3.6 Gates nicht ohne ADR lockern | permanent | nein — keine Schwelle gesenkt: die Gate-Konfiguration ist über die Welle unverändert (oben), `coverage-gate` 94,90 % gegen die Schwelle aus `SPEC-089`; `linelength` ist ein neues opt-in-Modul, nicht in `gates` aufgenommen |
| §3.7 Kommentare tragen eine der fünf Klassen | keiner benannt | — |
| §3.8 Ein Modul verspricht nur über seine Scan-Menge | permanent | nein — die Welle hat die Regel **angewandt**, nicht aufgelöst: `DC-FA-LLEN-001` beantwortet die Frage („außer der Konfiguration keine Datei außerhalb der Scan-Menge“), und ADR-0114 Entscheid 2 trägt sie. Eine Anwendung ist kein Trigger einer permanenten Regel |
| §3.9 Action-Referenzen SHA-gepinnt | durchgesetzt (`make workflow-pins`) | — keine Workflow-Änderung in der Welle |

## Vorschläge an den Planner

- **Keine Folge-ADR.** Kein Trigger ist eingetreten.
- **Regel 17** (Punkt 1b): Belege `metrics`/`pinsync` prüfen; bei Erreichen der Schwelle die
  Ausnahme-Zeile im Volltext verkörpern. Sonst bleibt der Eintrag offen im Register.
- **`yaml-dekoder-kuerzt-zahlen-still`** (Punkt 1a): bleibt im Register; der Weg ist ein Defekt-Fix nach
  `DC-FA-CONF-001`, nicht eine Entscheidung.

## Nicht geprüft

- Die Zweige Carveout und bootstrap-aware Gate (Planner).
- Ob die Fälle `metrics` und `pinsync` im Sinn des Registers als Belege zählen — das ist die Zählung des
  Planners, nicht dieses Audits.
- ADRs außerhalb der fünf genannten; die Welle zitiert weitere (ADR-0063, ADR-0069) nur als Präzedenz, ohne
  ihren Gegenstand zu berühren.
