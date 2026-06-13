# Lastenheft — d-check

**Version:** 0.8.0

**Status:** Draft

**Autor:** pt9912, **Datum:** 2026-06-10.

---

## 1. Zweck und Geltungsbereich

`d-check` ist ein Kommandozeilen-Tool, das Markdown-Dokumentation eines
Repositories auf kaputte Referenzen prüft: lokale Links und
Bildreferenzen, Heading-Anker, Linkpflicht für Anforderungs- und
Entscheidungs-Kennungen sowie Referenzrichtungs-Regeln zwischen
Dokumentklassen (Referenzmatrix). Es konsolidiert dreizehn funktional
überlappende Einzeltools (`check_refs.py`, `docs-check.js`,
`verify-doc-refs.sh`), die heute als Kopien in den
Schwester-Repositories des Entwicklungs-Workspace gepflegt werden.

Das Tool wird als Docker-Image über GHCR verteilt und in CI-Pipelines
sowie lokal als einzelner Prüfschritt aufgerufen. Die Funktionalität
ist in Regelmodule gegliedert, die pro Repository über
Kommandozeilen-Optionen und eine Konfigurationsdatei aktiviert werden —
ein Image, ein Update-Pfad, repo-spezifisches Verhalten per Config
statt per Code-Kopie.

## 2. Stakeholder

| Stakeholder | Rolle | Erwartung |
|---|---|---|
| Repo-Maintainer (pt9912) | Auftraggeber | Ein gepflegtes Tool statt dreizehn driftender Kopien; Fixes wirken überall |
| CI-Pipelines der Repos | Konsument | Ein Docker-Step mit deterministischen Exit-Codes und stabiler Ausgabe |
| AI-Agenten (Harness-Sensorik) | Konsument | Verlässlicher, maschinenlesbarer Doku-Sensor (z. B. als `make`-Gate eingebunden) |
| Harness-Projekte nach Kurs-Methodik | Nutznießer | Referenzmatrix- und ID-Regeln der Spec-Stratifizierung maschinell prüfbar |

## 3. Funktionale Anforderungen

> **Schema-Konvention.** Funktionale Anforderungen verwenden von Beginn
> an Bereichskürzel: `DC-FA-<BEREICH>-<NNN>`. Bereiche: `CLI`
> (Aufruf/Ausgabe), `SCAN` (Datei-Auswahl), `LINK` (Link-Modul), `ANCH`
> (Anker-Modul), `ID` (ID-Linkpflicht-Modul), `MTX`
> (Referenzmatrix-Modul), `EXT` (externe Links), `CODE`
> (Inline-Code-Pfade), `SPAN` (Markdown-Span-Artefakte), `HOST`
> (Host-Pfad-Hygiene), `CONF` (Konfiguration), `DIST` (Distribution).

### DC-FA-CLI-001 — Aufruf und Scan-Wurzel

**Beschreibung:** Das Tool wird als `d-check [pfad]` aufgerufen. Ohne
Argument ist die Scan-Wurzel das aktuelle Arbeitsverzeichnis; mit
Argument das angegebene Verzeichnis. Die Scan-Wurzel gilt als
Repository-Wurzel für alle Pfadauflösungen.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo ohne kaputte Referenzen, when `d-check` in der Repo-Wurzel aufgerufen wird, then Exit-Code 0 und eine Zusammenfassung „N Datei(en) geprüft, 0 Befunde".
- **Boundary:** Given ein Verzeichnis ohne Markdown-Dateien, when `d-check` aufgerufen wird, then Exit-Code 0, Zusammenfassung „0 Datei(en) geprüft", kein Fehler.
- **Negative:** Given ein nicht existierendes Verzeichnis als Argument, when `d-check /gibt/es/nicht` aufgerufen wird, then Exit-Code 2 und eine Fehlermeldung auf stderr.

**Out-of-Scope:** Prüfung mehrerer Repos in einem Aufruf.

---

### DC-FA-CLI-002 — Regelmodul-Auswahl

**Beschreibung:** Die Prüf-Funktionalität ist in benannte Regelmodule
gegliedert: `links`, `anchors`, `ids`, `matrix`, `external`,
`codepaths`, `spans`, `hostpaths`. Ohne Konfiguration sind `links` und
`anchors` aktiv. Module werden über
Kommandozeilen-Optionen (`--enable <modul>`, `--disable <modul>`)
und über die Konfigurationsdatei ([`DC-FA-CONF-001`](#dc-fa-conf-001--konfigurationsdatei))
aktiviert; Kommandozeilen-Optionen haben Vorrang vor der Konfiguration.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo mit einem kaputten Anker, aber korrekten Links, when `d-check --disable anchors` aufgerufen wird, then Exit-Code 0 (das Anker-Modul läuft nicht).
- **Boundary:** Given eine Konfigurationsdatei, die `ids` aktiviert, when `d-check --disable ids` aufgerufen wird, then läuft `ids` nicht (Kommandozeile schlägt Konfiguration).
- **Negative:** Given ein unbekannter Modulname, when `d-check --enable foo` aufgerufen wird, then Exit-Code 2 und eine Meldung mit der Liste gültiger Modulnamen.

**Out-of-Scope:** Dynamisches Nachladen externer Module (Plugins).

---

### DC-FA-CLI-003 — Exit-Codes

**Beschreibung:** Das Tool liefert deterministische Exit-Codes:
`0` = Prüfung gelaufen, keine Befunde; `1` = Prüfung gelaufen,
mindestens ein Befund; `2` = Nutzungs- oder Umgebungsfehler (ungültige
Option, ungültige Konfiguration, Scan-Wurzel nicht lesbar) — die
Prüfung hat dann keine verlässliche Aussage geliefert.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo ohne kaputte Referenzen, when `d-check` läuft, then Exit-Code 0.
- **Boundary:** Given genau ein Befund, when `d-check` läuft, then Exit-Code 1 (nicht 2), und der Befund wird ausgegeben.
- **Negative:** Given eine ungültige Option, when `d-check --foo` aufgerufen wird, then Exit-Code 2 und keine Prüf-Ausgabe.

**Out-of-Scope:** Differenzierte Exit-Codes pro Befund-Kategorie.

---

### DC-FA-CLI-004 — Ausgabeformate

**Beschreibung:** Befunde werden zeilenweise im Format
`<pfad>:<zeile>\t<ziel>\t<grund>` auf stdout ausgegeben (ein Befund pro
Zeile, Pfade relativ zur Scan-Wurzel); Zusammenfassung und Diagnose
gehen auf stderr. Mit `--json` erfolgt die gesamte Ausgabe auf stdout
als ein maschinenlesbares JSON-Dokument mit mindestens den Feldern
`findings` (Liste mit `file`, `line`, `target`, `rule`, `reason`),
`summary` (`filesChecked`, `findingCount`) und `exitCode`; stdout
enthält dann keine unstrukturierten Textzeilen.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo mit zwei kaputten Links, when `d-check` läuft, then genau zwei Befund-Zeilen auf stdout, je mit Datei, Zeilennummer, Ziel und Grund.
- **Boundary:** Given dieselben Befunde, when `d-check --json` läuft, then ist stdout als JSON parsbar, `summary.findingCount` ist 2 und stdout enthält keine Nicht-JSON-Zeilen.
- **Negative:** Given die unbekannte Option `--format` (kein Teil des CLI), when `d-check --json --format xml` aufgerufen wird, then Exit-Code 2 (ungültige Nutzung, vgl. [`DC-FA-CLI-003`](#dc-fa-cli-003--exit-codes)).

**Out-of-Scope:** Weitere Formate (SARIF, JUnit-XML) in dieser Version.

---

### DC-FA-SCAN-001 — Datei-Auswahl und Ignorier-Regeln

**Beschreibung:** Geprüft werden Markdown-Dateien (`*.md`) unterhalb
konfigurierter Scan-Wurzeln. Default ohne Konfiguration: `docs/`,
`spec/` (jeweils rekursiv, sofern vorhanden) sowie alle `*.md` direkt
in der Repository-Wurzel. Build- und Fremdverzeichnisse (`.git/`,
`node_modules/`, `build/`, `target/`, `.venv/`, `__pycache__/` u. ä.)
werden immer übersprungen. Zusätzliche Ignorier-Muster sind
konfigurierbar. Alle explizit in der Konfiguration deklarierten
Scan-Wurzeln müssen existieren; nur die Default-Wurzeln sind optional.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein kaputter Link in `docs/a.md`, when `d-check` ohne Konfiguration läuft, then wird der Befund gemeldet. <!-- d-check:ignore (AK-Beispiel) -->
- **Boundary:** Given ein kaputter Link in `node_modules/x/README.md`, when `d-check` läuft, then wird kein Befund gemeldet.
- **Negative:** Given eine Konfiguration, die explizit die Scan-Wurzel `handbuch/` deklariert, die nicht existiert, when `d-check` läuft, then Exit-Code 2 mit Hinweis auf die fehlende Wurzel (Default-Wurzeln dagegen sind optional).

**Out-of-Scope:** Prüfung anderer Markup-Formate (reStructuredText, AsciiDoc, HTML).

---

### DC-FA-LINK-001 — Lokale Link- und Bildreferenzen (Modul `links`)

**Beschreibung:** Markdown-Links `[text](ziel)` und Bildreferenzen
`![alt](ziel)` mit relativem Ziel werden geprüft: das Ziel muss nach
Pfadauflösung existieren und innerhalb der Repository-Wurzel liegen.
Die Pfadauflösung dekodiert Prozent-Escapes vollständig (RFC 3986,
z. B. `%20`); die Dekodierung erfolgt **vor** der Repo-Escape-Prüfung.
Ziele mit externem Schema (`http:`, `https:`, `mailto:`, `ftp:`)
werden von diesem Modul ignoriert; reine Anker-Ziele (`#...`) behandelt
das Modul `anchors`. Fragment-Teile gemischter Ziele (`ziel.md#anker`)
ignoriert dieses Modul ebenfalls — Anker-Validierung ist ausschließlich
Aufgabe des Moduls `anchors` (ist `anchors` deaktiviert, bleiben
fehlerhafte Anker unbemängelt). Vorkommen in Fenced-Code-Blöcken
(``` / ~~~) und in Inline-Code werden **von diesem Modul** nicht
geprüft (explizite Pfade in Inline-Code prüft das opt-in-Modul
`codepaths`,
[`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in));
HTML-Kommentare werden nicht gesondert behandelt (Links darin gelten
als Fließtext).

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Link auf eine existierende Datei im Repo, when das Modul `links` läuft, then kein Befund.
- **Boundary:** Given ein Link auf eine nicht existierende Datei innerhalb eines Fenced-Code-Blocks, when das Modul läuft, then kein Befund.
- **Negative:** Given ein Link `../../etc/passwd`, dessen Ziel die Repository-Wurzel verlässt, when das Modul läuft, then ein Befund mit Grund „verlässt Repository" — auch wenn der Zielpfad existiert. <!-- d-check:ignore (AK-Beispiel: Angriffs-Pfad) -->
- **Negative:** Given ein Link auf eine nicht existierende Datei, when das Modul läuft, then ein Befund mit Datei, Zeile, Ziel und Grund, Exit-Code 1.

**Out-of-Scope:** Reference-Style-Links (`[text][ref]` mit separater `[ref]: ziel`-Definition); semantische Prüfung, ob das verlinkte Dokument inhaltlich passt.

---

### DC-FA-LINK-002 — Symlink-Ablehnung

**Beschreibung:** Linkziele, die ein symbolischer Link sind oder deren
Pfad eine Symlink-Komponente enthält, erzeugen einen Befund —
unabhängig davon, wohin der Symlink zeigt (Defense-in-depth gegen
Verweise, die effektiv aus dem Repo herausführen). Die Symlink-Prüfung
hat Vorrang vor der Repo-Escape-Prüfung aus
[`DC-FA-LINK-001`](#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links);
pro Linkziel entsteht genau ein Befund.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Linkziel, das eine reguläre Datei ist, when das Modul `links` läuft, then kein Befund.
- **Boundary:** Given ein Symlink, der auf eine Datei *innerhalb* des Repos zeigt, when das Modul läuft, then trotzdem ein Befund mit Grund „Symlink".
- **Negative:** Given ein Symlink, der auf einen Pfad außerhalb des Repos zeigt, when das Modul läuft, then ein Befund; der Inhalt des Ziels wird nicht gelesen.

**Out-of-Scope:** Eine Konfigurations-Option, die Symlinks erlaubt.

---

### DC-FA-ANCH-001 — Heading-Anker-Validierung (Modul `anchors`)

**Beschreibung:** Links mit Fragment (`ziel.md#anker` sowie `#anker`
innerhalb derselben Datei) werden gegen die tatsächlichen Headings der
Zieldatei geprüft. Die Anker-Bildung folgt dem GitHub-Slug-Verfahren
(Kleinschreibung, Sonderzeichen-Entfernung, Leerzeichen → `-`,
Duplikat-Suffixe `-1`, `-2`, …). Headings aller Ebenen (`#`–`######`)
zählen; bei gleichlautenden Headings erhält das erste Vorkommen in
Dokumentreihenfolge den Basis-Slug, weitere die Suffixe `-1`, `-2`, ….
Existiert die Zieldatei nicht, wird die Anker-Prüfung für diesen Link
übersprungen — der Befund kommt vom Modul `links`
([`DC-FA-LINK-001`](#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)).

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Link `a.md#zweck-und-geltungsbereich` und ein Heading „Zweck und Geltungsbereich" in `a.md`, when das Modul läuft, then kein Befund.
- **Boundary:** Given zwei gleichlautende Headings „Beispiel" in der Zieldatei, when ein Link auf `#beispiel-1` zeigt, then kein Befund (Duplikat-Suffix korrekt aufgelöst).
- **Boundary:** Given ein Link `fehlt.md#x` auf eine nicht existierende Datei, when die Module `links` und `anchors` aktiv sind, then genau ein Befund — aus `links`; das Modul `anchors` schweigt.
- **Negative:** Given ein Link auf `a.md#gibt-es-nicht`, when das Modul läuft, then ein Befund mit Grund „Anker nicht gefunden".

**Out-of-Scope:** Anker in HTML-Tags (`<a name=…>`, `id=…`) als Linkziele.

---

### DC-FA-ID-001 — Linkpflicht für Kennungen (Modul `ids`)

**Beschreibung:** Die Konfiguration deklariert Kennungs-Muster (Regex,
z. B. `ADR-\d{4}`, `LH(-[A-Z0-9]+)+-\d{3}[A-Z]?`) und je Muster das
Ziel-Dokument bzw. Ziel-Verzeichnis, in dem die Kennung definiert ist.
Nackte Vorkommen einer Kennung im Fließtext — außerhalb von Inline-Code
und Fenced-Code-Blöcken — müssen als Markdown-Link auf ihre Definition
ausgeführt sein. Überlappen mehrere Muster für dasselbe Vorkommen,
gilt die Deklarationsreihenfolge in der Konfiguration: das erste
passende Muster gewinnt.

**Link-Politik (`link-policy`, je Muster konfigurierbar).** Pro Muster
ist wählbar, wie streng die Linkpflicht greift:

- `prose` (**Default**): wie oben — nur nackte Vorkommen im Fließtext
  sind linkpflichtig; Inline-Code-Vorkommen sind frei.
- `always`: auch Vorkommen **innerhalb von Inline-Code** müssen im
  Linktext eines Markdown-Links stehen (`` [`ADR-0042`](ziel) `` ist
  erfüllt, `` `ADR-0042` `` allein nicht). Fenced-Code-Blöcke,
  Heading-Zeilen und das `target` des Musters bleiben frei.

Der Default ist `prose`, damit bestehende Konfigurationen
byte-identisch bleiben und kein Repo ungefragt rote Läufe bekommt
(opt-in fürs Gating). Die Entdeckung ungenügender Verlinkung ist davon
unabhängig eine Bringschuld des Werkzeug-Betreibers (fleet-weiter
`always`-Lauf, [`DC-QA-04`](#dc-qa-04--migrationsabdeckung-der-alt-tools)-Muster);
„opt-in" heißt nicht „unsichtbar". Zwei Ventile halten `always`
treffsicher:

- `exempt-paths` (Glob-Liste je Muster): Dateien, in denen die strenge
  Regel nicht gilt (literal-schwere Artefakte wie Changelogs oder
  Review-Reports).
- Der Zeilen-Marker `d-check:ignore` (HTML-Kommentar, Begründung in
  Klammern empfohlen) nimmt eine Zeile von der `ids`-Prüfung aus — für
  bewusst illustrative Beispiel-IDs (gleiche Begründung wie bei
  [`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in);
  der Marker wirkt ab dieser Anforderung auf `codepaths` **und** `ids`).

**Akzeptanzkriterien:**

- **Happy Path:** Given das Muster `ADR-\d{4}` und ein Vorkommen `[ADR-0042](docs/plan/adr/0042-beispiel.md)`, when das Modul `ids` läuft, then kein Befund.
- **Boundary:** Given ein Vorkommen `` `ADR-0042` `` in Inline-Code und `link-policy: prose` (Default), when das Modul läuft, then kein Befund (Code-Vorkommen sind linkpflichtfrei).
- **Negative:** Given ein nacktes `ADR-0042` im Fließtext, when das Modul läuft, then ein Befund mit Grund „Kennung ohne Link".
- **`always` Happy Path:** Given `link-policy: always` und ein Vorkommen `` [`ADR-0042`](docs/plan/adr/0042-beispiel.md) ``, when das Modul läuft, then kein Befund (Code-Span im Linktext zählt als verlinkt).
- **`always` Negative:** Given `link-policy: always` und ein Vorkommen `` `ADR-0042` `` ohne Link (außerhalb `exempt-paths` und ohne `d-check:ignore`), when das Modul läuft, then ein Befund `id-unlinked`.
- **`always` Boundary (Ventile):** Given `link-policy: always`, ein `` `ADR-0042` `` in einer `exempt-paths`-Datei und ein zweites `` `ADR-0099` `` auf einer Zeile mit `d-check:ignore`, when das Modul läuft, then kein Befund für beide.

**Out-of-Scope:** Automatisches Ermitteln der Muster aus dem Repo-Inhalt; Prüfung, ob die verlinkte Definition inhaltlich zur Kennung passt; ein `link-policy`-Default abweichend von `prose` (bewusst opt-in).

---

### DC-FA-MTX-001 — Referenzmatrix zwischen Dokumentklassen (Modul `matrix`)

**Beschreibung:** Die Konfiguration deklariert Dokumentklassen über
Pfad-Muster (z. B. Contract-Spec `spec/lastenheft.md`, View-Spec
`spec/architecture.md`, ADR `docs/plan/adr/*.md`, Slice
`docs/plan/planning/**/slice-*.md`) und je Klassen-Paar, ob Referenzen
erlaubt sind. Zusätzlich sind Status-Bedingungen deklarierbar:
Referenzen auf Dokumente, deren Status-Feld (`## Status`-Heading oder
`**Status:**`-Zeile) einen verbotenen Wert hat (z. B. `superseded`,
`deprecated`), erzeugen einen Befund; fehlt das Status-Feld, gilt das
Dokument als aktiv. Damit sind die
Referenzrichtungs-Regeln der Spec-Stratifizierung (Lastenheft darf
nicht abwärts auf ADR/Slice verweisen; Slices nur auf aktive ADRs)
maschinell prüfbar.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Slice mit Link auf ein ADR mit Status `Accepted`, when das Modul `matrix` läuft, then kein Befund.
- **Boundary:** Given ein ADR mit Status `Superseded by ADR-0099`, when ein Slice darauf verlinkt, then ein Befund mit Grund „Referenz auf inaktives ADR".
- **Negative:** Given ein Link aus `spec/lastenheft.md` auf eine ADR-Datei, when das Modul läuft, then ein Befund mit Grund „verbotene Abwärtsreferenz" und Angabe beider Dokumentklassen.

**Out-of-Scope:** Semantische Unterscheidung von Verweis-Zwecken (z. B. Verifikations-Zeiger vs. Entscheidungsgrundlage) — das bleibt Review-Aufgabe; Provenance-/Historie-Sektionen können per Konfiguration von der Prüfung ausgenommen werden, eine automatische Erkennung solcher Sektionen ist nicht gefordert.

---

### DC-FA-EXT-001 — Externe Links (Modul `external`, opt-in)

**Beschreibung:** Nur bei explizit aktiviertem Modul werden `http:`-
und `https:`-Links per HTTP-Request (HEAD, Fallback GET) auf
Erreichbarkeit geprüft, mit konfigurierbarem Timeout (Default 10 s,
einstellbar über die Konfigurationsdatei) und begrenzter Parallelität.
Redirects (HTTP 301/302/303/307/308) werden bis zu fünf Stationen
gefolgt; längere Ketten erzeugen einen Befund. Das Modul ist nie Teil
der Default-Module —
Netzwerkzugriffe erfolgen ausschließlich opt-in.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Link auf eine erreichbare URL (HTTP-Status < 400), when das Modul `external` aktiviert ist, then kein Befund.
- **Boundary:** Given das Modul ist *nicht* aktiviert, when `d-check` in einer Umgebung ohne Netzwerk läuft, then erfolgt kein einziger Netzwerkzugriff und externe Links erzeugen keine Befunde.
- **Negative:** Given ein Link, der HTTP 404 liefert oder ins Timeout läuft, when das Modul aktiviert ist, then ein Befund mit Grund (Statuscode bzw. „Timeout").

**Out-of-Scope:** Link-Rotting-Archivierung, Redirect-Ketten-Bewertung, Authentifizierung gegen geschützte URLs.

---

### DC-FA-CODE-001 — Explizite Pfade in Inline-Code (Modul `codepaths`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `codepaths` werden
Inline-Code-Spans geprüft, deren Inhalt ein expliziter relativer Pfad
ist: beginnend mit `./` oder `../` (immer) oder mit einem der
konfigurierten Wurzel-Präfixe (z. B. `docs/`, `tools/` — relativ zur
Repository-Wurzel). Vor der Prüfung wird der Wert normalisiert:
umschließende Anführungszeichen und schließende Satzzeichen entfallen
(`` `../foo.md,` `` prüft `../foo.md`). Das Ziel muss nach <!-- d-check:ignore (Beispiel im Anforderungstext) -->
Pfadauflösung existieren und innerhalb der Repository-Wurzel liegen;
Fragment-Teile (`#…`) werden abgetrennt, und bei Markdown-Zielen wird
der Anker zusätzlich gegen die Headings der Zieldatei geprüft —
gleiches Slug-Verfahren wie
[`DC-FA-ANCH-001`](#dc-fa-anch-001--heading-anker-validierung-modul-anchors).
Werte mit Whitespace,
Platzhalter-/Glob-Zeichen oder Ellipsen gelten nicht als Pfad
(konservative Erkennung); Vorkommen in Fenced-Code-Blöcken werden
nicht geprüft. Ein HTML-Kommentar `d-check:ignore` auf derselben
Zeile (Begründung in Klammern empfohlen) nimmt die Zeile von **genau
dieser** Prüfung aus — bewusst nicht existierende Beispiel-Pfade
(etwa Angriffs-Beispiele in Lehrtexten) sind kein Fehler, sondern
eine zu dokumentierende Absicht. Für alle anderen Module existiert
kein Opt-out-Marker: deterministische Befunde werden behoben, nicht
unterdrückt.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Inline-Code-Span `` `docs/plan/adr/` `` auf ein existierendes Verzeichnis und das konfigurierte Präfix `docs/`, when das Modul `codepaths` läuft, then kein Befund.
- **Boundary:** Given ein Inline-Code-Span mit nicht existierendem Pfad und ein Kommentar `d-check:ignore` auf derselben Zeile, when das Modul läuft, then kein Befund — und der Marker hat keinerlei Wirkung auf Befunde anderer Module derselben Zeile.
- **Negative:** Given ein Inline-Code-Span `` `../fehlt.md` ``, dessen Ziel nicht existiert (oder die Repository-Wurzel verlässt), when das Modul läuft, then ein Befund mit Datei, Zeile, Ziel und Grund, Exit-Code 1.

**Out-of-Scope:** Pfad-Erkennung im Fließtext ohne Inline-Code; Pfade in Fenced-Code-Blöcken; Opt-out-Marker für andere Module; semantische Prüfung, ob der referenzierte Pfad inhaltlich passt.

---

### DC-FA-SPAN-001 — Markdown-Span-Artefakte (Modul `spans`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `spans` werden zwei
deterministische Artefakt-Klassen gemeldet, bei denen Markdown auf
GitHub nachweislich anders rendert, als der Quelltext nahelegt:

1. **Ungeschlossene Code-Spans** (`span-unclosed`): eine öffnende
   Backtick-Folge, die unmittelbar von Nicht-Whitespace gefolgt wird
   und im Absatz keine gleich lange schließende Folge findet
   (Absatzgrenzen: Leerzeile und Fence — dieselbe Absatz-Semantik wie
   die Vorverarbeitung der übrigen Module). Solche Opener kippen die
   Backtick-Parität des restlichen Absatzes; nachfolgende Links und
   Kennungen rendern zerrissen. Alleinstehende Backtick-Folgen
   (beidseitig Whitespace oder Zeilenrand) gelten als beabsichtigt
   literal und sind kein Befund (konservative Erkennung).
2. **Verschachtelte Link-Artefakte** (`span-nested-link`): Link-Syntax
   im Linktext eines weiteren Links (schließende Ziel-Klammer
   unmittelbar gefolgt von einer weiteren Linktext-Schließung mit
   Ziel-Öffnung) außerhalb von Fences und Code-Spans — für Leser
   sieht das wie ein Link aus, wird aber nur teilweise als Link
   geparst. Eine **Bildreferenz als Linktext** (Badge-Muster) ist
   legales Markdown und kein Treffer.

Der Befund nennt Datei, Zeile (Opener- bzw. Muster-Zeile), die
betroffene Backtick-Folge bzw. das Muster und den Grund. Es gilt die
allgemeine Opt-out-Regel unverändert: deterministische Befunde werden
behoben, nicht unterdrückt — der Zeilen-Marker `d-check:ignore`
bleibt auf das Modul `codepaths` beschränkt
([`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)).

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Absatz mit ausschließlich balancierten Code-Spans — auch solchen, die einen Zeilenumbruch enthalten —, when das Modul `spans` läuft, then kein Befund.
- **Boundary:** Given eine alleinstehende literale Backtick-Folge (beidseitig Whitespace), when das Modul läuft, then kein Befund.
- **Negative:** Given ein Listenpunkt, dessen öffnende Backtick-Folge unmittelbar vor Nicht-Whitespace steht und im Absatz ungeschlossen bleibt, when das Modul läuft, then ein Befund `span-unclosed` mit Datei, Zeile und Grund, Exit-Code 1.

**Out-of-Scope:** Emphasis-/Bold-Artefakte (`*`, `_`); Syntaxfehler in Fenced-Code-Blöcken; Reference-Style-Definitionen; automatische Korrektur; ein Opt-out-Marker für dieses Modul.

---

### DC-FA-HOST-001 — Host-lokale absolute Pfade (Modul `hostpaths`, opt-in)

**Beschreibung:** Bei explizit aktiviertem Modul `hostpaths` werden
absolute Pfade gemeldet, die ein Maschinen-Layout statt einer
Repo-Struktur beschreiben — sie funktionieren nur auf einem konkreten
Host und leaken dessen Verzeichnis-Aufbau. Erkannt werden, jeweils
als erstes Segment eines absoluten Pfads:

1. **Unix-Host-Präfixe** (Default-Liste, als Wurzel-Verzeichnisnamen
   deklariert: Development, home, Users, Volumes, mnt, media;
   per `hostpaths.prefixes` ersetzbar — tmp gehört bewusst nicht
   dazu: ein POSIX-Standard-Laufzeitort, dessen Erwähnung legitime
   Betriebs-Doku ist),
2. **Windows-Laufwerkspfade** (Laufwerksbuchstabe, Doppelpunkt,
   Backslash) und **UNC-Pfade** (doppelter Backslash plus
   Servername) — beide immer, nicht konfigurierbar.

Geprüft werden Prosa-Zeilen **einschließlich Inline-Code** (dort
leben solche Pfade typischerweise); Fenced-Code-Blöcke sind
ausgenommen — Beispiel- und Lehrinhalte mit bewussten Host-Pfaden
gehören in Fences. Vorbedingung ist eine Wortgrenze (kein
unmittelbar vorangehendes URL-, Pfad- oder Wortzeichen); schließende
Satzzeichen werden vom gemeldeten Pfad abgetrennt. Grund-Code
`hostpath-forbidden` mit Datei, Zeile und gefundenem Pfad. Es gibt
**keinen Opt-out-Marker** (der Zeilen-Marker `d-check:ignore` bleibt
auf das Modul `codepaths` beschränkt,
[`DC-FA-CODE-001`](#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)) —
die Auswege sind Fences für beabsichtigte Beispiele und die
Präfixliste für repo-spezifische Sonderfälle.

**Akzeptanzkriterien:**

- **Happy Path:** Given Dokumentation, deren Pfad-Angaben relativ oder Repo-Wurzel-absolut sind, when das Modul `hostpaths` läuft, then kein Befund.
- **Boundary:** Given ein host-lokaler absoluter Pfad innerhalb eines Fenced-Code-Blocks, when das Modul läuft, then kein Befund.
- **Negative:** Given eine Prosa-Zeile oder ein Inline-Code-Span mit einem Pfad unterhalb eines deklarierten Host-Präfixes, when das Modul läuft, then ein Befund `hostpath-forbidden` mit Datei, Zeile, Pfad und Grund, Exit-Code 1.

**Out-of-Scope:** Erkennung relativer Pfade (Aufgabe von `links`/`codepaths`); URL-Pfade hinter Schemata; Pfade in Fenced-Code-Blöcken; ein Opt-out-Marker für dieses Modul; automatische Umschreibung.

---

### DC-FA-CONF-001 — Konfigurationsdatei

**Beschreibung:** Eine optionale Datei `.d-check.yml` in der
Repository-Wurzel deklariert: Scan-Wurzeln, Ignorier-Muster, aktive
Module, Kennungs-Muster (für `ids`), Dokumentklassen und
Referenzregeln (für `matrix`) sowie Modul-Parameter (z. B. Timeout für
`external`). Ohne Konfigurationsdatei gelten die in den jeweiligen
Anforderungen definierten Defaults; eine vorhandene Datei wird beim
Start vollständig validiert — Syntax **und** Inhalt (z. B. ungültige
Regexe, unbekannte Modulnamen). Jeder Konfigurationsfehler führt zu
Exit-Code 2; es findet keine Prüfung mit stillschweigenden Defaults
statt.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine `.d-check.yml` mit einem Ignorier-Muster `docs/archive/**`, when `d-check` läuft, then erzeugen kaputte Links unter `docs/archive/` keine Befunde. <!-- d-check:ignore (AK-Beispiel) -->
- **Boundary:** Given keine `.d-check.yml`, when `d-check` läuft, then laufen genau die Default-Module (`links`, `anchors`) auf den Default-Scan-Wurzeln.
- **Negative:** Given eine syntaktisch ungültige `.d-check.yml`, when `d-check` läuft, then Exit-Code 2 mit Fehlermeldung inkl. Zeilenangabe; es wird keine Prüfung mit stillschweigenden Defaults durchgeführt.

**Out-of-Scope:** Konfigurations-Vererbung über mehrere Verzeichnisebenen; Migrations-Assistent von Alt-Tool-Konfigurationen.

---

### DC-FA-CONF-002 — Modul-lokaler Scan-Scope

**Beschreibung:** Jedes Regelmodul akzeptiert in der
Konfigurationsdatei optional einen Schlüssel `<modul>.scope` mit den
Unterschlüsseln `roots` und `ignore`. Ist er gesetzt, **ersetzt** er
für genau dieses Modul den globalen Scan-Scope
(`scan.roots`/`scan.ignore`) — bewusst keine Schnittmengen- oder
Vererbungs-Semantik (einfachste ehrliche Form). Module ohne `scope`
erben den globalen Scan unverändert; bestehende Konfigurationen
bleiben unverändert gültig (additiv, abwärtskompatibel). Für
`scope.roots` und `scope.ignore` gelten dieselben Regeln und
Constraints wie für die globalen Schlüssel: deklarierte Wurzeln
müssen existieren und innerhalb der Repository-Wurzel liegen (jede
Verletzung Exit 2,
[`DC-FA-CONF-001`](#dc-fa-conf-001--konfigurationsdatei)),
Ignorier-Muster prunen den Verzeichnis-Abstieg, die `SKIP_DIRS`
gelten immer, und eine explizit leere `roots`-Liste prüft nichts
(Leere-Menge-Semantik wie `scan.roots`).

**Akzeptanzkriterien:**

- **Happy Path:** Given global `scan.roots: ["."]` und ein `ids`-`scope` mit `roots: ["spec"]`, when `d-check` läuft, then stammen `id-unlinked`-Befunde ausschließlich aus `spec/`, während `links`/`anchors` weiterhin den globalen Scope prüfen.
- **Boundary:** Given ein Modul ohne `scope`-Schlüssel, when `d-check` läuft, then verhält es sich byte-identisch zum Verhalten ohne diese Anforderung.
- **Boundary:** Given ein Modul-`scope` mit explizit leerer `roots`-Liste, when `d-check` läuft, then prüft dieses Modul nichts — bewusste Setzung, kein Fehler.
- **Negative:** Given ein `scope.roots`-Eintrag, der nicht existiert oder die Repository-Wurzel verlässt, when `d-check` startet, then Exit-Code 2 mit Fehlermeldung; es findet keine Prüfung mit stillschweigenden Defaults statt.

**Out-of-Scope:** Scope je Muster (`ids.patterns[].roots` — erst bei zweitem Bedarfsträger); Schnittmengen-/Unions-Semantik zwischen globalem und Modul-Scope; alternative Konfigurationsdateien, Mehrfach-Configs und Vererbung (Out-of-Scope-Linie von [`DC-FA-CONF-001`](#dc-fa-conf-001--konfigurationsdatei)).

---

### DC-FA-DIST-001 — Docker-Image

**Beschreibung:** Das Tool wird als Container-Image
`ghcr.io/pt9912/d-check` mit Semver-Tags veröffentlicht. Aufruf:
`docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:<tag>`
prüft das nach `/repo` gemountete Repository. Das Image läuft als
Non-root-Prozess; ein read-only-Mount genügt, da das Tool das geprüfte
Repository nie beschreibt. Default-Befehl des Images ist die Prüfung
von `/repo`; CLI-Optionen werden als Container-Argumente angehängt.
Ergebnis und Exit-Code sind identisch zur nativen Ausführung.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Repo mit einem kaputten Link, when der Container mit gemountetem Repo läuft, then identische Befund-Ausgabe und Exit-Code 1 wie beim nativen Aufruf.
- **Boundary:** Given ein read-only gemountetes Repo (`:ro`), when der Container läuft, then keine Schreibfehler und vollständige Prüfung.
- **Negative:** Given kein Mount auf `/repo` (Verzeichnis leer oder fehlt), when der Container läuft, then Exit-Code 2 mit Hinweis auf den erwarteten Mount.

**Out-of-Scope:** Distributionswege jenseits GHCR (Homebrew, Paketmanager); deren Bewertung erfolgt später per ADR.

## 4. Nichtfunktionale Anforderungen

### DC-QA-01 — Performance

- **Anforderung:** Ein vollständiger Lauf der Default-Module über ein Repository mit 1.000 Markdown-Dateien (insgesamt ≤ 20 MB) terminiert in < 5 s auf Standard-Entwickler-Hardware (2 vCPU); das Modul `external` ist von dieser Schranke ausgenommen.
- **Messmethode:** Benchmark-Lauf gegen ein generiertes Fixture-Repo, Definition in `spec/spezifikation.md`.

### DC-QA-02 — Determinismus

- **Anforderung:** Identische Eingabe (Repo-Stand + Konfiguration + Optionen) liefert byte-identische Ausgabe; Befunde sind stabil sortiert (Pfad, dann Zeilennummer).
- **Messmethode:** Test führt denselben Lauf 10× aus und vergleicht Hashes der Ausgabe.

### DC-QA-03 — Seiteneffektfreiheit und Netzwerk-Sparsamkeit

- **Anforderung:** Das Tool schreibt nie in das geprüfte Repository und öffnet außer im explizit aktivierten Modul `external` keine Netzwerkverbindungen.
- **Messmethode:** Integrationstest mit read-only-Mount und netzwerkloser Umgebung (`docker run --network none`), alle Module außer `external` aktiv.

### DC-QA-04 — Migrationsabdeckung der Alt-Tools

- **Anforderung:** Jedes der dreizehn Quell-Tool-Vorkommen ist durch `d-check` mit passender Konfiguration ersetzbar: Auf dem jeweiligen Repo-Stand meldet `d-check` mindestens dieselben echten Befunde wie das Alt-Tool und erzeugt keine False-Positives, die eine bislang grüne CI brechen.
- **Messmethode:** Pilot-Migration in mindestens drei Repos — je ein Vertreter der Shell-Familie (`verify-doc-refs.sh`), der Python-Familie (`check_refs.py`, inkl. u-boot-Vollausbau) und der JS-Familie (`docs-check.js`) — mit Vergleichslauf Alt-Tool vs. `d-check`.

## 5. Globale Out-of-Scope-Punkte

- Architektur-/Boundary-Checks (C++-Include-Regeln, Go-Package-Graph, Rust-Crate-Grenzen) — bleiben repo-lokale, sprachspezifische Tools.
- Build-Reproduzierbarkeits-Prüfung (Funktionalität von `repro-check.sh`).
- Mathematik-/Formel-Validierung (MathJax-Rendering-Checks aus `euler-fourier-hilbert`) — Kandidat für eine spätere Version, nicht Teil von 0.x.
- Automatisches Reparieren kaputter Referenzen (Auto-Fix) — das Tool berichtet nur.
- Prüfung von Nicht-Markdown-Formaten (reStructuredText, AsciiDoc, HTML, Jupyter-Notebooks).
- Rechtschreib-, Stil- oder Markdown-Lint-Prüfungen (dafür existieren etablierte Tools).

## 6. Glossar

| Begriff | Bedeutung im Lastenheft |
|---|---|
| Befund | Eine einzelne festgestellte Regelverletzung mit Datei, Zeile, Ziel und Grund. |
| Regelmodul | Benannte, einzeln aktivierbare Prüf-Einheit (`links`, `anchors`, `ids`, `matrix`, `external`, `codepaths`, `spans`, `hostpaths`). |
| Scan-Wurzel | Verzeichnis, unterhalb dessen Markdown-Dateien gesucht werden; zugleich Bezugspunkt der Pfadauflösung. |
| Anker | Fragment-Teil eines Links (`#…`), das auf ein Heading der Zieldatei zeigt (GitHub-Slug-Verfahren). |
| Repo-Escape | Linkziel, dessen aufgelöster Pfad außerhalb der Repository-Wurzel liegt. |
| Kennung | Textuelle ID nach deklariertem Muster (z. B. `ADR-0042`), für die Linkpflicht gelten kann. |
| Dokumentklasse | Über Pfad-Muster definierte Gruppe von Dokumenten (z. B. Contract-Spec, ADR, Slice) als Knoten der Referenzmatrix. |
| Referenzmatrix | Deklaration, welche Dokumentklasse auf welche verweisen darf, inkl. Status-Bedingungen. |
| Aktives ADR | ADR, dessen Status-Feld keinen verbotenen Wert (`superseded`, `deprecated`) trägt. |
| Quell-Tools | Die dreizehn konsolidierten Alt-Tool-Vorkommen in den Schwester-Repositories des Entwicklungs-Workspace: zwölf aus drei Familien (Shell: `verify-doc-refs.sh`, Python: `check_refs.py`, JavaScript: `docs-check.js`) plus eine eigenständige Python-Linie (`check_markdown_links.py`, Inventur-Nachtrag 2026-06-12). |

## 7. Historie

| Version | Datum | Änderung | Verweis |
|---|---|---|---|
| 0.8.0 | 2026-06-13 | Change Request (Auftraggeber): `DC-FA-ID-001` um konfigurierbare `link-policy: prose\|always` (je Muster) erweitert — `always` macht auch Inline-Code-Vorkommen linkpflichtig, Default `prose` (opt-in, abwärtskompatibel). Zwei Ventile: `exempt-paths` (Glob-Liste je Muster) und der Zeilen-Marker `d-check:ignore`, dessen Geltungsbereich von `codepaths`-only auf `ids` erweitert wird (illustrative Beispiel-IDs). Anlass: der `ids`-Sensor maß bislang nicht das Ziel „gut verlinkt" — ein Code-Span konnte stillschweigend einen fehlenden Link verbergen (Ausgangsbefund `DC-QA-03` in slice-017). Kalibrierung über die drei `ids`-Repos (d-check 155, u-boot 9, b-trace 2) bestimmte die Ventil-Form | slice-018 |
| 0.1.0 | 2026-06-10 | Initiale Fassung (Konsolidierung von 12 Quell-Tools, Modul-Schnitt, Docker-Distribution) | — |
| 0.2.0 | 2026-06-10 | Review-Runde R1: Modul-Schnitt `links`/`anchors` präzisiert (Fragment-Zuständigkeit, fehlende Zieldatei), Slug-Duplikat-Reihenfolge, Symlink-Vorrang, RFC-3986-Dekodierung vor Escape-Prüfung, Redirect-Regel `external`, Muster-Präzedenz `ids`, Status-Default `matrix`, Scan-Wurzel- und Config-Vollvalidierung, Out-of-Scope Reference-Style-Links, Image-Default-Befehl | — |
| 0.2.1 | 2026-06-10 | Redaktionell: Beispiel-Kennungen in DC-FA-ID-001/DC-FA-MTX-001/Glossar auf fiktive Nummern (`ADR-0042`, `ADR-0099`) umgestellt — Kollision mit real entstandenen/zukünftigen eigenen ADRs vermeiden; keine inhaltliche Änderung | — |
| 0.2.2 | 2026-06-10 | Redaktionell: absolute Workspace-Pfade entfernt („Schwester-Repositories des Entwicklungs-Workspace" statt konkreter Pfade); keine inhaltliche Änderung | — |
| 0.2.3 | 2026-06-11 | Redaktionell: „(folgt)" in der `DC-QA-01`-Messmethode entfernt — die Benchmark-Definition existiert in der Spezifikation; keine inhaltliche Änderung | — |
| 0.7.2 | 2026-06-12 | Review zum Change Request `DC-FA-HOST-001` (Kalibrierungs-Befund slice-016): tmp aus der Default-Präfixliste gestrichen — der Kalibrierungslauf zeigte überwiegend legitime Laufzeit-Doku (Log-/Output-Pfade) statt Maschinen-Layout-Leaks; wer tmp prüfen will, konfiguriert es | slice-016 |
| 0.7.1 | 2026-06-12 | Review zum Change Request `DC-FA-SPAN-001` (Kalibrierungs-Befund slice-015): Bildreferenzen als Linktext (`[![…](…)](…)` — Badge-Muster, z. B. Shields in vendorten Paket-READMEs) sind legales Markdown und kein `span-nested-link`-Treffer | slice-015 |
| 0.7.0 | 2026-06-12 | Change Request (Erst-Bedarfsträger grid-gym): neue Anforderung `DC-FA-CONF-002` — optionaler modul-lokaler Scan-Scope `<modul>.scope` (`roots`/`ignore`), ersetzt für genau dieses Modul den globalen Scope; additiv und abwärtskompatibel, Constraints spiegeln `scan.*`. Anlass (gemessen, v0.2.0): `ids`-Aktivierung in grid-gym liefert global 2776 Befunde (Masse: Retro-Verlinkung historischer done-Planning-Docs = Umschreiben des Audit-Trails) vs. 312 echte, fixbare Befunde im kuratierten Scope `spec/` + `docs/user/` — der Linkpflicht-Nutzen ist heute nur um den Preis des globalen Sweeps oder des Verzichts auf breite `links`/`anchors`-Abdeckung zu haben | slice-017 |
| 0.6.0 | 2026-06-12 | Change Request (Auftraggeber): neue Anforderung `DC-FA-HOST-001` — Modul `hostpaths` meldet host-lokale absolute Pfade in Prosa und Inline-Code (opt-in; Unix-Präfixliste konfigurierbar, Windows-/UNC-Muster fest; Fences ausgenommen, kein Opt-out-Marker). Anlass: der bess-ems-Rest-Sensor generalisiert (dort als eigenes Tool gebaut) plus die 8 Host-Pfad-Links aus dem d-migrate-Vergleichslauf und die eigene 0.2.2-Hygiene-Korrektur — dieselbe Leak-Klasse dreimal unabhängig. Bereich `HOST` deklariert; Modul-Listen in `DC-FA-CLI-002` und Glossar ergänzt | slice-016 |
| 0.5.0 | 2026-06-12 | Change Request (Auftraggeber): neue Anforderung `DC-FA-SPAN-001` — Modul `spans` meldet ungeschlossene Code-Spans und verschachtelte Link-Artefakte (opt-in, konservative Erkennung: alleinstehende Backticks bleiben befundfrei; kein Opt-out-Marker). Anlass: die `DC-QA-04`-Vergleichsläufe aus slice-012/slice-014 — ~100 u-boot-Artefakte sowie Fälle in grid-gym und bess-ems wurden nur indirekt als `id-unlinked`-Folgefehler sichtbar; ein direkter Sensor meldet die Ursache statt der Symptome. Bereich `SPAN` in der Schema-Konvention deklariert; Modul-Listen in `DC-FA-CLI-002` und Glossar ergänzt | slice-015 |
| 0.3.0 | 2026-06-11 | Change Request (Auftraggeber): neue Anforderung `DC-FA-CODE-001` — Modul `codepaths` prüft explizite Pfade in Inline-Code (opt-in, konservative Erkennung, Zeilen-Opt-out `d-check:ignore` nur für dieses Modul). Anlass: `DC-QA-04`-Vergleichslauf gegen die JS-Familie (`docs-check.js`) zeigte die Prüfklasse als Konsolidierungs-Lücke. Bereich `CODE` in der Schema-Konvention deklariert; Modul-Listen in `DC-FA-CLI-002` und Glossar ergänzt | slice-013 |
| 0.4.0 | 2026-06-12 | Change Request (Auftraggeber): Inventur-Nachtrag zu `DC-QA-04` — dreizehntes Alt-Tool-Vorkommen entdeckt (`check_markdown_links.py` in bess-ems, eigenständige Python-Linie, entstanden 2026-05-24 und damit vor der Inventur vom 2026-06-10 übersehen); Anforderungs-Text, Einleitung, Stakeholder-Tabelle und Glossar von zwölf auf dreizehn fortgeschrieben. Messmethode (drei Familien-Piloten) unverändert | slice-014 |
| 0.3.1 | 2026-06-11 | Review R1 zum Change Request: `DC-FA-CODE-001` präzisiert — Wert-Normalisierung (Anführungszeichen, schließende Satzzeichen) und Anker-Prüfung bei Markdown-Zielen (`DC-QA-04`-Parität zur JS-Familie); `DC-FA-LINK-001`-Inline-Code-Aussage auf Modul-Bezug eingegrenzt | slice-013 |
