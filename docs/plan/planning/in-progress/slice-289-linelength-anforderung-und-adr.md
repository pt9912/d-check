# Slice slice-289: Anforderung, Spezifikation und ADR für das Modul `linelength`; Entscheide festschreiben

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-95-linelength-modul`](../welle-95-linelength-modul.md).

**Bezug:** [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in) — die Anforderung des Moduls
`linelength`, die dieser Slice ins Lastenheft schreibt (Bereichs-Kürzel `LLEN`),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(Scope), [`DC-FA-SCAN-001`](../../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei),
[`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope);
Abgrenzung gegen [`DC-FA-FILE-001`](../../../../spec/lastenheft.md#dc-fa-file-001--zeilen--und-byte-obergrenzen-einer-ganzen-datei-modul-file-opt-in).

**Berührte Spec-Stellen:** `lastenheft.md` §3 (Bereichs-Kürzel-Liste und die neue
Anforderung), `spezifikation.md` §1 (neues Verfahren), §2 (Schema `linelength.*`),
§4 (Grund-Code — die Tabellenzeile kommt mit dem Code); `architecture.md` bleibt
unberührt (die Rolle „Regel-Modul“ besteht). Der Verweis zeigt **aufwärts**: Die Spec
nennt diesen Slice nie.

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-11-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Anforderung des Moduls `linelength` steht im Lastenheft, ihr Verfahren,
ihr Schema und ihr Grund-Code in der Spezifikation und die Entscheidung in einer
`Accepted`-ADR — und die Entscheide unten sind **vor dem ersten Code** vom
Auftraggeber bestätigt und festgeschrieben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Konfigurations-Rand, Tests des Moduls** — gehören
  [`slice-290`](../open/slice-290-linelength-modulkern.md); dieser Slice schreibt, was der Code
  einlösen soll, damit die Reviews dort gegen eine Anforderung prüfen.
- **Die Zeile in der Grund-Code-Tabelle (§4) und die Modul-Aufzählung in
  [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) und
  im Glossar** — kommen mit dem Code in
  [`slice-290`](../open/slice-290-linelength-modulkern.md): die Spiegel-Tests halten beide gegen
  das registrierte Modul und wären vorher rot (bei `pinsync` gemessen).
- **Andere Dateiarten** — vom Auftraggeber ausgeschlossen (2026-11-22); die Anforderung
  nennt es in ihrem Out-of-Scope, damit es nicht über `file` oder einen Glob-Schalter
  still hineinwächst.
- **CHANGELOG, READMEs, Handbuch** — Release-Prep
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17), kein Feature-Commit.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Die Anforderung steht im Lastenheft (Bereichs-Kürzel in der Liste ergänzt,
      Version gehoben, Historie-Zeile nach
      [`MR-032`](../../../../harness/conventions.md#mr-032)) mit Akzeptanzkriterien je
      Entscheid (Zähleinheit, jede Ausnahme samt ihren Rändern, Grenzwert genau auf der
      Schwelle, leere Menge, Zeile über der Grenze mit `d-check:ignore` ⇒ Befund), den Rändern
      mit Exit 2, den benannten Grenzen der geteilten Lexik und ihrem Out-of-Scope.
- [ ] Die Spezifikation führt das Verfahren (Leseweise, Zählregel, Ausnahmen, Reihenfolge
      der Befunde, Bezugsmenge), das Schema `linelength.*` und den Grund-Code im Verfahren;
      jede Aussage über die Grenzen (welche Eingaben das Modul liest, [`AGENTS.md`](../../../../AGENTS.md)
      §3.8) ist am Gegenstand gezählt oder gemessen, nicht aus der Beschreibung übernommen.
- [ ] [ADR-0114](../../adr/0114-linelength-modul.md) liegt `Proposed` vor, ist im
      [ADR-Index](../../adr/README.md) geführt, trägt `## Re-Evaluierungs-Trigger`, und die
      Anforderungs-Kennung und das `Schärft:`-Feld sind verlinkt (Architect). `Accepted` wird sie
      mit der Closure von [`slice-290`](../open/slice-290-linelength-modulkern.md) (Verdikt F-8).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit** Wellen von der nächsten Welle-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | neue Anforderung, Kürzel in der Liste, Version und Historie |
| `spec/spezifikation.md` | update | Verfahren, Schema `linelength.*`, Grund-Code im Verfahren |
| [ADR-0114](../../adr/0114-linelength-modul.md) | update — **schreibt der Architect** | liegt `Proposed` vor (vom Architect geschrieben, samt Index-Zeile); nachzuziehen sind die Verlinkung der Anforderungs-Kennung und das `Schärft:`-Feld, sobald die Anforderung im Lastenheft steht. Der Implementer liest die ADR als Constraint und meldet eine nötige Änderung als Befund an den Architect |
| Plan von [`slice-290`](../open/slice-290-linelength-modulkern.md) | update — Planner | die Testfälle, die aus den bestätigten Entscheiden und dem Architect-Verdikt folgen, und die zwei mit dem Code kommenden Spiegel-Stellen |

**Plan-Änderung nach dem Architect-Verdikt (2026-11-22).** Das
[Architect-Verdikt](../../../reviews/2026-11-22-welle-95-architect-verdikt.md) (F-1 bis F-8)
ist eingearbeitet: Entscheid 1, 5 und 7 sind am Gegenstand korrigiert bzw. präzisiert, Entscheid 8
ist gegen [`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope)
neu gefasst, Entscheid 9 ist neu, und die ADR bleibt `Proposed` bis zur Closure von
[`slice-290`](../open/slice-290-linelength-modulkern.md) (F-8). Die ADR schreibt der Architect,
nicht dieser Slice (Modul 8 §Rollen-Regeln); die Zeile „ADR-Index“ entfällt deshalb aus der Tabelle.
Die Reihenfolge war nicht regelkonform: der Slice lag in `in-progress/`, bevor das Verdikt vorlag
(dort als einmaliges akzeptiertes Negativ geführt; kein Code war entstanden).

**Festzuschreibende Entscheide** — vom Auftraggeber am 2026-11-22 bestätigt: 1–7 mit den
Empfehlungen (1, 5 und 7 danach vom Architect am Gegenstand geschärft, ohne die Wahl zu ändern),
8 in der neuen Fassung und 9 neu, beide nach dem Architect-Verdikt. Maßgeblich für die Begründung
ist [ADR-0114](../../adr/0114-linelength-modul.md).

1. **Eigenes Modul statt Bedingung.** Modul `linelength`, Bereichs-Kürzel `LLEN`
   (`LINK` ist belegt). Verworfen: eine dritte Schwelle im Modul `file` — es liest jede
   Dateiart roh über eigene Globs, die Markdown-Ausnahmen (Fenced-Code, Tabellen)
   brächen seine Zusage „unabhängig von der Dateiart“; eine Bedingung in `structure` —
   es wählt seine Dateien über eigene Globs und kennt kein `scope`, jede Regel braucht einen
   Abschnitts-Selektor, und es misst den bereinigten Text (Verdikt F-6; der frühere Grund
   „meldet auf der Überschriftszeile“ stimmte am Gegenstand nicht, `cell-max-chars` meldet auf
   der Zeile selbst).
2. **Eingabemenge:** die Markdown-Dateien der Scan-Menge (`scan.roots`/`scan.ignore`,
   [`DC-FA-SCAN-001`](../../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln))
   mit modul-lokalem `scope`
   ([`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope))
   und `exempt-paths` wie die anderen Prosa-Module. Das Modul liest keine Datei außerhalb
   dieser Menge — die §3.8-Frage hat damit die Antwort „keine“, und die Spezifikation
   sagt es.
3. **Zähleinheit:** Unicode-Codepunkte der rohen Zeile ohne Zeilenende (`\n`, ein `\r`
   davor zählt nicht) — dieselbe Einheit wie `cell-max-chars`
   (`utf8.RuneCountInString`). Ein Tab zählt als ein Zeichen. Verworfen: Bytes (ein
   Umlaut zählte doppelt) und Anzeigebreite (braucht eine Breiten-Tabelle, East-Asian
   Wide, Emoji — nicht deterministisch über Unicode-Versionen hinweg). Ungültiges UTF-8
   zählt je Byte als ein Zeichen (Verhalten von `RuneCountInString`), die Anforderung
   nennt es.
4. **Schwelle und Schlüsselname:** `linelength.max-chars` (Pflicht bei aktivem Modul,
   ≥ 1, sonst Exit 2); geprüft wird `Länge > max-chars`, eine Zeile genau auf der
   Grenze ist befundfrei. Bewusst **nicht** `max-lines`/`max-length`: der Name
   `max-lines` trägt schon zwei Bedeutungen
   ([`BEO-ALL/max-lines-namenskollision-structure-file`](../observations/BEO-ALL/max-lines-namenskollision-structure-file/observation.md)),
   eine dritte nahe Form vergrößerte die Verwechslung.
5. **Ausnahmen, je einzeln schaltbar, Default „ausgenommen“:** (a) Zeilen in
   Fenced-Code-Blöcken samt ihrer Zäune (`exempt-code-blocks`), (b) Zeilen außerhalb von
   Fenced-Code, die mit `|` beginnen (`exempt-tables`; die Lexik verlangt keine Trennzeile und
   trimmt nicht — Verdikt F-4), (c) Zeilen, deren Teil ab der Grenze keinen Leerraum enthält —
   Leerraum sind Space und Tab, geprüft werden die Codepunkte an den Positionen `max-chars + 1`
   bis Zeilenende; eine lange URL, ein langer Link, ein langer Inline-Code-Pfad lässt sich
   nicht umbrechen (`exempt-unbreakable`; Verdikt F-5). Überschriften zählen mit. Die Erkennung von
   Fenced-Code und Tabellenzeile nimmt das Modul aus der geteilten Lexik
   (`proseLineSet`, `tableRowLine`), kein eigenes Muster daneben
   ([`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md));
   die Grenzen dieser Lexik (Umschalten bei jeder Zaun-Zeile, eingerückter Code-Block und Zaun im
   Zitat werden gemessen, offener Zaun nimmt den Rest der Datei aus — gemeldet von `spans` als
   `fence-unclosed`; eingerückte, zitierte und GFM-Zeilen ohne führende Pipe sind keine
   Tabellenzeilen) nennt die Anforderung (Verdikt F-3, F-4).
   Front-Matter ist keine eigene Ausnahme (d-check kennt ihn heute nicht als Form).
6. **Befund:** je überlanger Zeile ein Befund, Ort = Datei und Zeile, Meldung nennt
   gemessene Länge und Grenze; Grund-Code `line-chars-exceeded` (Empfehlung, analog
   `file-lines-exceeded`/`file-bytes-exceeded`). Reihenfolge nach Pfad, dann Zeile
   ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
7. **Bezugsmenge:** die Vollständigkeits-Zeile nennt geprüfte Dateien und geprüfte
   Zeilen **nach** Abzug der Ausnahmen, damit ein grünes Ergebnis nicht mehr behauptet,
   als es misst (Baseline-Regelwerk `modul-13-quality-gates.md` §Hard Rule, *Ein Gate
   ohne seine Grenze*). Geschärft (Verdikt F-7): Zeilen zählen wie bei `file`
   (Zeilenumbrüche plus eine unvollständige Schlusszeile); „geprüft“ ist die Menge nach `scope`,
   `exempt-paths`, (a) und (b); die über (c) ausgenommenen Zeilen stehen als eigene Zahl
   daneben; Ausgabe auf stderr bzw. in `summary.notes`.
8. **Leere Menge ist kein Fehler** (neu gefasst nach Verdikt F-1, bestätigt 2026-11-22): trifft
   die Eingabemenge keine Markdown-Datei, meldet das Modul nichts, und die Bezugsmenge nennt
   `0`; eine explizit leere `linelength.scope.roots` prüft nichts, wie
   [`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope) es
   jedem scannenden Modul zusagt. Die frühere Fassung (Exit 2, „fail-closed wie
   `structure`/`file`“) widersprach dieser Anforderung, und die Präzedenz trug nicht: dort ist die
   Nullmenge ein Befund mit Exit 1 eines nicht scannenden Moduls.
9. **Kein Zeilen-Marker** (neu nach Verdikt F-2, bestätigt 2026-11-22): `linelength` honoriert
   `d-check:ignore` nicht; die Ventile sind (a)–(c), `scope` und `exempt-paths`. Eine Zeile über
   der Grenze mit Marker ist ein Befund (der Marker zählt mit).

**Gegenprobe vor dem Schreiben (Schritt 19 des Workflow-Skeletts):** die Zählregel läuft
am eigenen Bestand (`awk '{ if (length($0) > 120) n++ } END { print n }'` über die
gescannten `.md`, in einer UTF-8-Locale; Zahl mit und ohne Fenced-Code steht im
Plan-Nachtrag) — das belegt die Out-of-Scope-Aussage „eigenes Profil nicht scharf“ mit
einer Zahl statt einer Vermutung. Die Negativfälle (Zeile genau auf der Grenze,
Umlaut/Emoji, `\r\n`, Tab, nicht geschlossener Fence bis Dateiende, Tabellenzeile ohne
Trennzeile, URL mit Leerzeichen davor) stehen als Testfälle für
[`slice-290`](../open/slice-290-linelength-modulkern.md).

**Plan-Nachtrag — Gegenprobe gelaufen (2026-11-22, Implementer).** Gemessen auf dem Stand `f3dc94da` (vor diesem Nachtrag) über die
**globale** Scan-Menge des eigenen Profils (`scan.roots: ["."]`, die drei `scan.ignore`-Muster,
die feste Skip-Liste, nur reguläre `*.md`-Dateien, keine Symlinks) — die Menge, die `linelength`
ohne eigenen `scope` erbt. Grenze 120, Zähleinheit Codepunkte der rohen Zeile ohne `\n`, ein `\r`
vor dem `\n` zählt nicht; Fenced-Code nach dem Umschalt-Automaten der geteilten Lexik.

| Größe | Zahl |
|---|---|
| Dateien · Zeilen | 1237 · 86275 |
| Zeilen über 120, mit Fenced-Code | 6108 |
| Zeilen über 120, ohne Fenced-Code | 6082 |
| davon Tabellenzeilen (Ausnahme b) | 3130 |
| davon ohne Leerraum ab Position 121 (Ausnahme c) | 828 |
| Befunde mit allen drei Ausnahmen an (Default) | 2124 |
| zum Vergleich: über 120 **Bytes** statt Codepunkten, mit Fenced-Code | 6233 |

Damit ist „das eigene Profil schaltet das Modul nicht ein“ belegt: bei 120 stünden 2124 Befunde
an. Kommando (aus der Repo-Wurzel, eine einzige `awk`-Instanz; das Ende des Programms gibt die
Zahlen aus):

```sh
find . \( -type d \( -name .git -o -name node_modules -o -name build -o -name target \
    -o -name dist -o -name vendor -o -name .venv -o -name __pycache__ -o -name .idea \
    -o -name .vscode -o -name .gradle \) -prune \) \
  -o \( -path ./.harness/baseline -o -path ./.harness/cache \
    -o -path ./tools/blackbox-probe/fixtures \) -prune \
  -o -type f -name '*.md' -print0 \
| LC_ALL=C.UTF-8 xargs -0 awk -v G=120 '
  FNR == 1 { infence = 0; files++ }
  { raw = $0; sub(/\r$/, "", raw); lines++
    t = raw; sub(/^[ \t]+/, "", t); toggle = 0
    if (substr(t, 1, 3) == "~~~") toggle = 1
    else if (substr(t, 1, 3) == "\140\140\140") {
      rest = t; sub(/^\140+/, "", rest); if (index(rest, "\140") == 0) toggle = 1 }
    code = (toggle || infence); if (toggle) infence = !infence
    if (length(raw) > G) { all++
      if (!code) { prose++
        if (substr(raw, 1, 1) == "|") table++
        else if (substr(raw, G + 1) ~ /[ \t]/) hit++
        else unbreak++ } } }
  END { print files, lines, all, prose, table, unbreak, hit }'
```

Die Bytes-Zeile ist dasselbe Kommando unter `LC_ALL=C`. Gegenprobe der Menge: `make doc-check`
meldet 1243 Dateien — die Vereinigung aller Modul-Mengen; die 6 Dateien mehr sind die
Fixture-Dateien unter `tools/blackbox-probe/fixtures/`, die `citations.scope` nicht ausnimmt
(`find tools/blackbox-probe/fixtures -type f -name '*.md' | wc -l` ⇒ 6).

**Am Gegenstand gelesen, für die Grenz-Aussagen der Spezifikation** (Kommando neben der Aussage):

- Der Fence-Automat teilt an `\n` und lässt ein `\r` in der Zeile:
  `grep -n 'strings.Split(string(content), "\\n")' internal/hexagon/core/rules/markdown.go` ⇒ Z. 91.
- Leerraum vor dem Zaun sind Space und Tab, beliebig viele:
  `grep -n 'TrimLeft(raw, " \\t")' internal/hexagon/core/rules/markdown.go` ⇒ Z. 53 (`TrimFenceIndent`).
- Tabellenzeile = in der Prosa-Menge und die **rohe** Zeile beginnt mit `|`:
  `grep -n 'HasPrefix(lines\[i\], "|")' internal/hexagon/core/rules/markdown.go` ⇒ Z. 193.
- Die Scan-Menge nimmt nur Dateien mit der Endung `.md` (Groß-/Kleinschreibung zählt) und läuft
  nicht durch Symlinks: `grep -n 'HasSuffix(e.Name, ".md")\|case driven.Kind' internal/hexagon/core/rules/scan.go`
  ⇒ Z. 94, 112, 119, 120.
- Eine unlesbare Datei der Scan-Menge beendet den ganzen Lauf mit Exit 2:
  `grep -n 'nicht lesbar: ' internal/hexagon/core/rules/run.go` ⇒ Z. 232.
- Den Schlüssel `<modul>.scope` nehmen heute 16 Module an; `linelength` kommt als siebzehntes dazu:
  `sed -n '/^func applyScopes/,/^	}$/p' internal/adapter/driven/configyaml/configyaml.go | grep -c '{"'` ⇒ 16.
- Den Zeilen-Marker honorieren vier Module (`codepaths`, `ids`, `versions`, `diagrams`):
  `grep -n 'ignoreMarker nimmt' -A1 internal/hexagon/core/rules/codepaths.go`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Welle `welle-95-linelength-modul` ist offen und das
WIP-Limit ist frei (`in-progress/` trägt nur die Roadmap).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Entscheide 1–8 brauchen
  mehr als eine Review-Sitzung, oder die Ausnahme (c) verlangt eine eigene Link-Lexik,
  die die geteilte nicht hergibt.
- `in-progress` → `open` (blockiert): der Auftraggeber entscheidet einen der Punkte
  1, 3, 5, 8 oder 9 nicht; ohne Entscheid beginnt kein Code.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Lastenheft, Spezifikation und ADR liegen vor und `make gates` ist grün),
Review- und Verifikations-Report liegen unter `docs/reviews/`, Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Ausnahme (c) „kein Leerzeichen ab der Grenze“ ist zu grob: eine Zeile mit langem
  Link **und** umbrechbarem Text davor wird ausgenommen, obwohl sie kürzbar wäre —
  **Ausgang:** *entfallen*, wenn die Anforderung diese Grenze der Ausnahme ausdrücklich
  nennt und die ADR die strenge Variante (nur ein Wort, das allein länger als die Grenze ist,
  wird ausgenommen) verglichen hat — die früher hier genannte Alternative „nur der Überschuss
  nach dem letzten Leerzeichen zählt“ ist bis auf eine Randposition dieselbe Regel (Verdikt F-5);
  *weiter offen* ins Register, wenn Review oder Verifikation die Grenze als Mangel lesen.
- Ein vergessener Schluss-Fence macht für die geteilte Lexik den Rest der Datei zu Code
  (die Fence-Hälfte von
  [`BEO-ALL/prosa-verschwindet-in-inline-code-spanne`](../observations/BEO-ALL/prosa-verschwindet-in-inline-code-spanne/observation.md),
  2×); nimmt das Modul Fenced-Code über dieselbe Lexik aus, fällt jede überlange Zeile
  danach still aus der Prüfung — ein dritter Beleg wäre die Schwelle. Die
  Inline-Code-Hälfte trifft das Modul nicht: es misst die rohe Zeile —
  **Ausgang:** *entfallen*, wenn die Anforderung das Verhalten beim offenen Fence festlegt,
  als Grenze benennt und für die Meldung auf `spans` (`fence-unclosed`) verweist — ein zweiter
  Wächter in `linelength` wäre eine zweite Antwort (Verdikt F-3; der Test dazu ist Lieferung von
  [`slice-290`](../open/slice-290-linelength-modulkern.md)); *weiter offen* — ein Beleg in
  `BEO-ALL/prosa-verschwindet-in-inline-code-spanne`, der den Zähler auf 3× hebt —, wenn
  sie es offen lässt.
- Der Name `linelength` liest sich neben `file.max-lines` und `structure.max-lines` als
  vierte Längen-Achse und wird verwechselt — **Ausgang:** *entfallen*, wenn Anforderung und
  ADR die Längen-Achsen gegeneinander abgrenzen und das Review dazu keinen Befund meldet;
  *weiter offen* — ein Beleg in `BEO-ALL/max-lines-namenskollision-structure-file` —, wenn
  es einen meldet.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist der Default `*` (`ALL`) —
Lastenheft, Spezifikation, ADR; keine Sub-Area mit eigener Modus-Deklaration
(`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand
2026-11-22), Sub-Area `*`. Treffer:
[`BEO-ALL/max-lines-namenskollision-structure-file`](../observations/BEO-ALL/max-lines-namenskollision-structure-file/observation.md)
(offen, 1× — Entscheid 4 meidet einen dritten nahen Schlüsselnamen; Risiko 3 in §6),
[`BEO-ALL/prosa-verschwindet-in-inline-code-spanne`](../observations/BEO-ALL/prosa-verschwindet-in-inline-code-spanne/observation.md)
(offen, 2× — Risiko 2 in §6; ein dritter Beleg wäre die Schwelle),
[`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md)
(offen, 2× — Entscheid 5 verlangt die geteilte Lexik),
[`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md)
(offen, 1× — jede Randbedingung der Anforderung braucht in
[`slice-290`](../open/slice-290-linelength-modulkern.md) einen Test) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — Entscheid 2 beantwortet die
Frage). Der Zähler folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-22: **beide rot**. `upstream-drift` (jüngster Lauf
2026-11-22T07:40:16Z): lokal nachgeprüft fünf veraltete Pins — Go 1.27.3 → 1.27.4,
Semgrep 1.183.0 → 1.186.0, a-check 0.23.3 → 0.24.0, `docker/login-action` 4.7.0 → 4.8.0,
Digest von `golang:1.27.3`. `image-scan` (2026-11-22T11:29:01Z): behebbare
CRITICAL/HIGH im publizierten Image, die Kennungen stehen im Job-Summary des Laufs (das
Log war über `gh` nicht abrufbar). Beides ist planmäßige Fremd-Drift, berührt dieses
Modul nicht und gehört als eigene Pin-Hebung vor das nächste Release
(Welle-Out-of-Scope).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
