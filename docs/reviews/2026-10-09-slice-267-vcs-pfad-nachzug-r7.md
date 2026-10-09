# Review R7: slice-267, `vcs` lässt einen reinen Pfad-Nachzug durch

- **Review-Art:** Code. Geprüft wird `0876b9d0..c9748c61`: `f89914c7` (Plan-Änderung nach R6),
  `ecb8443e` (`fix(vcs)`, Normierung auflösender Ziele statt Leerung) und `c9748c61`
  (`ADR-0103`-Geschichte). Geprüft gegen den Slice-Plan `slice-267` samt Plan-Änderung, die
  Findings aus R1 bis R6 (`2026-10-09-slice-267-vcs-pfad-nachzug-r1.md` bis `-r6.md`), `ADR-0103`,
  `DC-FA-VCS-001`, `MR-025`, `MR-066` und die Hard Rules `AGENTS.md` §3.5, §3.7, §3.8 sowie §5
  Regel 13 und 15. Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-267` · `0876b9d0..c9748c61`
- **Skill:** `reviewer.md` @ 1.19.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** Lastenheft 0.102.5 (`DC-FA-VCS-001`, Absatz „Pfad-Nachzug (opt-in)",
  Kriterium „Boundary (Pfad-Nachzug)", Historie 0.102.5), Spezifikation `DC-FA-VCS-001.a`
  Schritt 4, §2-Zeile `vcs.ignore-link-targets` und Historie, `harness/sensors/adr-check.md`
  Grenze 4, `ADR-0103` (Geschichte). Im Code: `vcs.go` (`CheckVCS`, `vcsModified`, `vcsCore`),
  `vcs_link_targets.go` (`normalizedLinkTargetLines`, `linkTargetResolver`, `pathTree`,
  `replaceRanges`, `linkTargetCuts`, `containsLink`, `opaqueLines`), zum Vergleich `paths.go`
  (`ResolveTarget` des Moduls `links`), `model/config.go`, `cli/config_template.go`, dazu
  `vcs_pfad_nachzug_test.go`. Vorherige Findings am Modul: R1 bis R6 zu `slice-267`.
- **Proben:** Das Image ist per `make build` vom Stand `c9748c61` gebaut (`d-check:latest`). Alle
  Proben laufen in einem Wegwerf-Klon des Repos im Scratchpad: Änderung, ein Commit, Lauf
  `--enable vcs` mit `FOCUS_DISABLE` aus dem `Makefile` und `--range HEAD~1..HEAD` gegen die
  `.d-check.yml` des Repos (`ignore-link-targets: true`), danach `git reset --hard origin/main`.
  Die Mutationsprobe läuft per `make test IMAGE=d-check-r7mut` im Klon. Klon und Mutations-Image
  sind danach entfernt.

## Proben am gebauten Image

| # | Änderung im Klon | Ergebnis |
|---|---|---|
| P1 | `git mv docs/user/releasing.md docs/user/maintainer/releasing.md`, die vier Links in `ADR-0014` nachgezogen | 0 Befunde, Exit 0 |
| P2 | wie P1, ein Link zeigt stattdessen auf `docs/user/operations.md` | `core-drift-vcs` `ADR-0014`, Exit 1 |
| P3 | ohne Umzug: ein Link auf `../../../docs/user/maintainer/fehlt.md` | `core-drift-vcs`, Exit 1 |
| P4 | ohne Umzug: ein Link in `ADR-0014` von `../../../docs/user/releasing.md` auf `releasing.md` (fehlt in `docs/plan/adr/`) | **0 Befunde, Exit 0** |
| P4c | ohne Umzug: `ADR-0024` Zeile 20, Ziel `../../../spec/spezifikation.md#dc-fa-vcs-001a--…` auf `spezifikation.md#dc-fa-vcs-001a--…` | **0 Befunde, Exit 0**; der Voll-Lauf (Modul `links`) meldet `target-missing`, Exit 1 |
| P5 | alle 13 Ziele `../../../harness/README.md` in den ADRs auf `README.md` (den ADR-Index) | 0 Befunde, Exit 0 (benannte Grenze, siehe I-2) |
| P6 | zweiter Umzug: `git mv harness/conventions.md docs/harness/conventions.md`, 46 Links in 21 ADRs nachgezogen (`../../harness/conventions.md#mr-…`) | 0 Befunde, Exit 0 |
| P6c | wie P6, dazu ein Wort im Körper von `ADR-0024` geändert | `core-drift-vcs` `ADR-0024`, Exit 1 |
| P7 | P1 gestagt, Lauf mit `--staged` (Weg des `pre-commit`-Hooks) | 0 Befunde, Exit 0 |
| P8 | ohne Umzug: `releasing.md` → `Releasing.md` im Ziel | `core-drift-vcs`, Exit 1 |
| P9 | ohne Umzug: `releasing.md` → `releasing%2Emd` im Ziel | `core-drift-vcs`, Exit 1 |
| P10 | ohne Umzug: Ziel `../../../docs/plan/../user/./releasing.md` (dieselbe Datei) | 0 Befunde, Exit 0 |
| Mut | `vcs.go`: `linkTargetResolver(path, baseTree)` durch `linkTargetResolver(path, headTree)` ersetzt | `make test` grün, Paket `core/rules` `ok` |

## Findings

### M-1 — Ein Nachzug auf ein fehlendes Ziel geht durch, wenn der neue Rohtext dem normierten alten Ziel gleicht; zwei Grenz-Texte versprechen mehr Drift, als der Code liefert

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 10 und 18 · `DC-FA-VCS-001` Kriterium „Boundary (Pfad-Nachzug)" · `AGENTS.md`
  §5 Regel 13
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`return path.Base(p) + anchor, true`"
  und „`if norm, ok := resolve(r[c[0]:c[1]]); ok {`"; `spec/lastenheft.md` · „zeigt das neue Ziel
  auf eine Datei mit anderem Namen oder auf keine existierende"; `spec/spezifikation.md` · „Das
  Werkzeug prüft, dass das neue Ziel existiert und denselben Namen trägt" und (§2-Zeile) „ein
  Nachzug auf eine andere Datei und jede Änderung am Linktext"; `harness/sensors/adr-check.md` ·
  „Es prüft damit, dass das neue Ziel existiert"; `ADR-0103` Geschichte · „verengt sich damit auf
  Dateien gleichen Namens"
- **Befund:** Ein auflösendes Ziel wird zu `Dateiname#Anker`. Ein Ziel, das nicht auflöst, bleibt
  roh. Schreibt HEAD das Ziel als bloßen Dateinamen samt Anker, der im Verzeichnis der ADR nicht
  existiert, sind beide Fassungen byte-gleich. Der Nachzug auf ein fehlendes Ziel ist dann keine
  Drift (P4, P4c). Lastenheft-Kriterium, Spezifikation und Sensor-Datei sagen das Gegenteil
  („existiert"). Die `ADR-0103`-Geschichte sagt, die Grenze verenge sich auf Dateien gleichen
  Namens. Die §2-Zeile sagt außerdem, ein Nachzug auf „eine andere Datei" bleibe Drift. Das stimmt
  nur für einen anderen Namen, P5 geht durch. Schritt 4 derselben Spezifikation benennt das
  richtig. Der Test „neues ziel existiert nicht" nimmt ein Ziel mit anderem Dateinamen
  (`x/fehlt.md`). Er prüft damit die Namensbedingung, nicht die Existenz.
- **Failure-Szenario:** Ein Commit „vereinfacht" in einer `Accepted`-ADR den Verweis
  `../../../spec/spezifikation.md#…` zu `spezifikation.md#…`. `make adr-check` und der
  `pre-commit`-Hook bleiben grün (P4c). In diesem Repo fängt `make doc-check` den toten Link
  (`target-missing`). Ein Konsument, der `vcs` ohne `links` fährt, bekommt das Grün, das ihm das
  Lastenheft-Kriterium ausdrücklich nicht verspricht.
- **Warum nicht HIGH:** Sichtbarer Inhalt ändert sich nicht. Die Änderung trifft nur den
  unsichtbaren Zielslot, und im Repo meldet das Partner-Gate `links` sie.
- **Verifizierbar:** ja. P4 und P4c. Ein Kern-Test mit dem bloßen Dateinamen als neuem Ziel und
  erwartetem Befund läuft heute rot.
- **Klasse:** `grenze-gegen-gegenstand`

### M-2 — Dass jeder Stand gegen seinen eigenen Baum auflöst, hat keinen Test

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 13 · Spezifikation `DC-FA-VCS-001.a` Schritt 4 („BASE gegen den Pfad-Baum
  von BASE, HEAD gegen den von HEAD")
- **Pfad:** `internal/hexagon/core/rules/vcs.go` · „`linkTargetResolver(path, baseTree)`";
  `internal/hexagon/core/rules/vcs_pfad_nachzug_test.go` · „`for _, ref := range []string{"BASE",
  "HEAD"} {`"
- **Befund:** Der Test legt `nachzugBaum()` in **beide** Stände. Kein Fall hat eine Datei, die nur
  in BASE oder nur in HEAD liegt, also keinen echten Umzug. Wird der BASE-Baum durch den HEAD-Baum
  ersetzt (Mut), bleibt `make test` grün. Dasselbe gilt für jede andere Vertauschung oder
  Vereinigung der beiden Bäume. Die Mutationsliste in `ecb8443e` nennt die „Baum-Prüfung". Rot
  wird sie über „verzeichnis eines nicht aufloesenden ziels geaendert", also über die Existenz an
  sich, nicht über die Zuordnung zum Stand.
- **Failure-Szenario:** Eine Vereinfachung bildet einen gemeinsamen Baum aus BASE und HEAD. Die
  Suite bleibt grün. Danach geht ein Nachzug auf eine Datei durch, die in HEAD gelöscht ist (sie
  löst im gemeinsamen Baum noch auf). Mit dem HEAD-Baum für beide Seiten (Mut) wird stattdessen
  jeder echte Umzug laut zu Drift, also gerade der Anlassfall P1.
- **Verifizierbar:** ja. Mutation oben, `make test` grün.
- **Klasse:** `filterzweig-ohne-negativtest`

### L-1 — Drei Spiegel beschreiben noch die Leerung

- **Kategorie:** LOW
- **Quelle:** `MR-025` (Spiegel einer Semantik-Änderung) · `AGENTS.md` §3.7 („Ein Kommentar
  beschreibt, was da ist")
- **Pfad:** `internal/hexagon/core/model/config.go` · „`// Links auf eine leere Form: ein reiner
  Pfad-Nachzug ist keine Core-Drift,`"; `internal/hexagon/core/rules/vcs_link_targets.go` ·
  „`// opaqueLines liefert die 1-basierten Zeilen, deren Links nicht geleert`" und „`bleibt
  ungeleert`"; `internal/adapter/driving/cli/config_template.go` · „`reiner Pfad-Nachzug (nur
  Link-Ziele geändert) ist keine Drift`"
- **Befund:** Der Feldkommentar im Modell sagt, das Ziel werde auf eine leere Form gebracht. Der
  Kommentar zu `opaqueLines` spricht von „nicht geleert". Die `--print-config`-Vorlage sagt, eine
  Änderung nur an Link-Zielen sei keine Drift. Ein Nachzug auf einen anderen Namen oder auf ein
  fehlendes Ziel ist seit `ecb8443e` aber Drift. Der Plan zählt Modell und Vorlage selbst als
  Spiegel des Schlüssels auf.
- **Verifizierbar:** nein (Lesen). Kein Gate prüft Kommentar-Wahrheit.
- **Klasse:** `spiegel-nicht-nachgezogen`

### L-2 — R6 L-1 ist nur im Code benannt, nicht in Spezifikation und Sensor-Datei

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §5 Regel 13
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`Link, den ein zeilenlokal falsch
  gepaarter Code-Span verdeckt.`"; `spec/spezifikation.md` · „Fail-safe bleibt ein Nachzug Drift,
  wo die Erkennung oder ein Filter den Link nicht durchlässt:"; `harness/sensors/adr-check.md` ·
  „Fail-safe bleibt Drift: Ziel auf"
- **Befund:** Der `GRENZE`-Kommentar nennt den falsch gepaarten zeilenlokalen Code-Span. Die
  Fail-safe-Aufzählungen in Spezifikation und Sensor-Datei nennen ihn nicht. Er trifft aber den
  Bestand (`ADR-0028` Zeile 85, R6 B2a): Ein echter Nachzug dort meldet Drift, ohne dass die
  Grenze ihn erklärt. Die Richtung ist fail-safe.
- **Verifizierbar:** ja. R6 B2a.
- **Klasse:** `grenze-nur-eine-richtung`

### I-1 — Der Resolver liest Ziele anders als `ResolveTarget` des Moduls `links`

- **Kategorie:** INFO
- **Quelle:** Prüffrage 11 (unter MEDIUM-Schwelle: kein stiller Pfad)
- **Pfad:** `internal/hexagon/core/rules/vcs_link_targets.go` · „`if t == "" ||
  strings.HasPrefix(t, "/") ||`"; `internal/hexagon/core/rules/paths.go` · „`decoded, err :=
  url.PathUnescape(target)`"
- **Befund:** `links` dekodiert Prozent-Kodierung und liest `/…` als wurzel-relativ. Der Resolver
  von `vcs` tut beides nicht, solche Ziele bleiben roh. Ein echter Umzug einer so verlinkten Datei
  wird damit Drift (vgl. P9). Die Richtung ist fail-safe. Gezählt: 0 prozent-kodierte und 0
  absolute Ziele in den ADRs (Zählform: `](/` und `](…%XX…)` per `grep`). Die Spezifikation
  benennt das absolute Ziel, das kodierte nicht.
- **Verifizierbar:** ja. P9 als Form. Ein echter Umzug mit kodiertem Link ist nicht gefahren.
- **Klasse:** `zwei-module-eine-eingabe`

### I-2 — Größe der benannten Restgrenze „gleicher Name, andere Datei", gemessen

- **Kategorie:** INFO
- **Quelle:** `DC-FA-VCS-001.a` Schritt 4 („nicht, dass es dieselbe Datei ist")
- **Pfad:** `spec/spezifikation.md` · „nicht, dass es dieselbe Datei ist"
- **Befund:** Form des Gezählten: Ziel-Token `](…` bis `)`, `#` oder Leerraum in `Accepted`-ADRs
  (Geschichte eingeschlossen, ohne `://`), deren letzter Pfadteil unter `git ls-files` mehrfach
  vorkommt. Treffer: 111 von 1289 Tokens, verteilt auf `conventions.md` (46),
  `AGENTS.md` (28), `README.md` (18), `.d-check.yml` (11), `observation.md` (5) und drei
  `modul-*.md` (5). Angesehen: `conventions.md` und `AGENTS.md` haben je einen Namensvetter
  unter `.claude/rules/`, `observation.md` 64. Für diese Ziele kann ein Umlenken auf eine andere
  Datei gleichen Namens still bleiben (P5). Die Grenze ist benannt, und sichtbarer Text ändert
  sich dabei nicht. Ob diese Reichweite trägt, ist eine Frage an `ADR-0103`, kein Code-Befund.
- **Verifizierbar:** ja. P5.
- **Klasse:** `restgrenze-gemessen`

## Status der R6-Findings

- **R6 M-1:** strukturell aufgelöst. Inhaltstext in Klammern löst nicht als Pfad auf und bleibt roh.
  Die Formen (verschachtelter Link, `|` ohne Leerraum, Escapes, Titel einer Referenz-Definition)
  haben Kern-Tests mit erwartetem Befund. Am Image nicht einzeln nachgefahren.
- **R6 M-2:** aufgelöst und benannt. Hinter einer Listenmarke geht nur noch ein Nachzug eines
  auflösenden Ziels mit gleichem Namen durch, kein Inhaltswort. Spezifikation, Sensor-Datei und
  Kommentar nennen das.
- **R6 M-3:** aufgelöst. Die Tests „fence-schliesser anderes zeichen" und „fence-schliesser zu
  kurz" sind vorhanden.
- **R6 L-1:** teilweise. Im Code benannt, in Spezifikation und Sensor-Datei nicht (L-2).
- **R6 L-2:** aufgelöst. Alle drei Stellen nennen die CRLF-Grenze jetzt nur für die
  Referenz-Definition.
- **R6 I-1:** beantwortet durch den Auftraggeber-Entscheid. Die Klasse
  `muster-leert-mehr-als-gegenstand` ist für Inhaltstext geschlossen. Die neue Achse
  (Auflösung) hat eigene Ränder (M-1, I-2).

## Umgehungs-Achsen der Auflösung (Bewertung)

- **Verzeichnis, `.`, `..`:** lösen auf und werden auf den Verzeichnisnamen normiert. Ein Umlenken
  bleibt nur bei gleichem Namen still (I-2). Auf Wurzelebene löst `.` nicht auf (fail-safe).
  Realistisch nur als Nachzug.
- **Gleiche Datei, andere Schreibweise** (P10): 0 Befunde. Das ist richtig, es bleibt dieselbe
  Datei.
- **Groß-/Kleinschreibung** (P8), **URL-Kodierung** (P9): lösen nicht auf, Drift. Fail-safe.
- **Anker:** bleibt wörtlich im Vergleich. Eine Anker-Änderung ist Drift (Test „anker geaendert").
- **Inhaltswort, das zufällig ein Dateiname ist:** Normiert wird nur der Zielslot, und die beiden
  Fassungen sind nur gleich, wenn der Dateiname gleich bleibt. Sichtbarer Text kann sich darüber
  nicht ändern. Kein realistischer Pfad.
- **Code hinter einer Listenmarke:** sichtbarer Code geht nur durch, wenn ein auflösender Pfad auf
  einen gleichnamigen anderen Pfad wechselt. Benannt, im Bestand keine Vorkommen (R6-Zählung).
- **Bloßer Dateiname als fehlendes Ziel:** der einzige Restpfad gegen einen ausdrücklichen
  Vertragstext (M-1).

## Negativbefunde (geprüft, ohne Befund)

- **Anlass am Bestand:** P1 (vier Links, `ADR-0014`) und der zweite Umzug P6 (46 Links, 21 ADRs)
  ergeben 0 Befunde. P7 zeigt das auch für den gestagten Weg. Gegenprobe P6c: ein Wort mehr meldet.
  Kein Befund.
- **Gegenrichtung:** anderes Dokument (P2) und fehlendes Ziel mit anderem Namen (P3) melden. Kein
  Befund über M-1 hinaus.
- **`ADR-0103` gegen `273bebea`:** `git diff 273bebea..HEAD` an der Datei hat nur zwei
  `+`-Zeilen, beide in der Geschichte-Tabelle. `make adr-check RANGE=273bebea..HEAD` ergibt 0
  Befunde (13 Commits). §3.5 ist eingehalten. Zum Inhalt der neuen Zeile siehe M-1.
- **Plan-Änderung (`f89914c7`) und `MR-066`:** Die Notiz steht vor dem Code-Commit. Sie nennt den
  Grund und die Ersatz-Form der Prüfung (jede Runde misst am Image beide Richtungen), wie `MR-066`
  es verlangt. Diese Runde vollzieht sie. Die Rückführung aus §4 ist begründet ausgeschlossen. Kein
  Befund.
- **Lastenheft 0.102.5 und Spezifikations-Historie:** Versionskopf und neueste Zeile oben passen,
  die Anker lösen auf. Schritt 4 beschreibt Normierung, Stände, Filter und den Listen-Code
  entsprechend dem Code. Kein Befund über M-1 und L-2 hinaus.
- **`harness/sensors/adr-check.md`:** Vertrag und Grenze 4 sind nachgezogen, die CRLF-Grenze ist
  eingeengt. Kein Befund über M-1 und L-2 hinaus.
- **Kommentare §3.7 (Klasse):** `normalizedLinkTargetLines`, `linkTargetResolver`, `pathTree`,
  `replaceRanges`, `containsLink`, `linkTargetCuts` und der Test-Kommentar zu `nachzugBaum` tragen
  Zusage oder `GRENZE`. Keine Review-Historie, keine Befund- oder Slice-Nummern. Zur Wahrheit
  siehe L-1.
- **`replaceRanges`:** sortiert, überspringt Überlappungen, und die Ersetzung ist rein
  positionsbasiert auf der Rohzeile. Zeilen, deren vorverarbeitete Länge abweicht, werden
  übersprungen (fail-safe). Kein Befund.
- **Determinismus (`DC-QA-02`), Hexagon (`ADR-0005`), Netz (`DC-QA-03`), Suppressions:** Die
  Menge wird nur per Schlüssel nachgeschlagen, nie iteriert ausgegeben. Neu sind nur `path`,
  `regexp`, `sort` und `strings`. Kein Netz, kein `//nolint`. Kein Befund.
- **Opt-in:** Ohne Schlüssel bleiben die Bäume `nil`, der Resolver ist `nil`, der Zweig ist
  unverändert. Jeder Test fährt beide Modi. Kein Befund.
- **§3.8 (Eingaben außerhalb der Scan-Menge):** Der Pfad-Baum kommt aus `AllPaths`. Das ist
  dieselbe fail-closed Quelle, aus der die geschützte Klasse kommt, im gestagten Modus also der
  Index. Kein zusätzlicher Leseweg. Kein Befund.
- **Botschaft `ecb8443e` (§5 Regel 15):** Der Bestandsbeleg ist auf `ADR-0014` begrenzt und so
  benannt. „Strukturell geschlossen" bezieht sich auf Inhaltstext und hält. Die Probe „Ziel
  existiert nicht" ist korrekt gemessen, aber nur in einer Form (siehe M-1). Kein eigener Befund.
- **Traceability:** Alle drei Betreffs tragen `slice-267`, dazu `DC-FA-VCS-001`, `ADR-0103` oder
  `MR-066`. Kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 (M-1, M-2) |
| LOW | 2 (L-1, L-2) |
| INFO | 2 (I-1, I-2) |

## Verdikt

**Kein HIGH. Zwei MEDIUM sind vor Closure zu klären.** Die Normierung auflösender Ziele schließt
die Klasse aus R4 bis R6 für Inhaltstext. Der Anlassfall und ein zweiter realistischer Umzug
(`harness/conventions.md`, 46 Links) ergeben am gebauten Image 0 Befunde, auch gestagt. Ein
anderes Dokument und ein fehlendes Ziel mit anderem Namen melden. Offen bleibt ein schmaler
Restpfad gegen den Vertragstext: ein fehlendes Ziel, das als bloßer Dateiname geschrieben ist
(M-1). Im Repo meldet ihn `links`, `vcs` allein nicht. Außerdem hat die Zuordnung „jeder Stand
gegen seinen Baum" keinen Test, der eine Vertauschung fängt (M-2). Die zwei LOW betreffen
Spiegel- und Grenztexte.
