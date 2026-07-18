# Slice slice-079: Zitat-Verifikation — `datei:zeile` als gemessenes Property

**Status:** in-progress (welle-62), **aktiviert 2026-07-18** — beide Vorfragen
entschieden.
**Adopter-Rückfrage empirisch beantwortet:** die `datei:zeile`-Zitate im Adopter-Repo
`ai-harness-init` stehen **33/33 in Inline-Code**, null in nackter Prosa (gemessen);
der Adopter hat seinen eigenen Sensor-Slice deshalb selbst blockiert. **⇒
`codepaths`-Erweiterung**, kein Prosa-Scanning. **Vorfrage 1 (Zuschnitt) entschieden
(Auftraggeber):** alle drei Stufen jetzt, **Form (c)** — Stufe 1/2 als Erweiterung von
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
Stufe 3 (`verbatim`) als **eigenes Modul** (setzt auf `codepaths`' Detektion auf, wie
`anchors` auf `links`). **Vorfrage 2** war bereits aufgelöst (gemeinsames
Kürzel-Kriterium mit
[slice-078](../done/slice-078-ignore-refs-quell-skopus.md)).

**Bezug:** **Change Request** eines Adopters (`ai-harness-init`), eingereicht
2026-07-17. Berührt
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(Modul `codepaths`) — dessen Algorithmus
[`DC-FA-CODE-001.a`](../../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)
Schritt 3 die Zeilennummer heute **ausdrücklich verwirft**. **Noch kein ADR und
keine Lastenheft-Änderung** — beides folgt, wenn §4 entschieden ist. IDs werden
nicht ad hoc vergeben ([`AGENTS.md`](../../../../AGENTS.md) §5).

**Autor:** pt9912 (CR: Adopter `ai-harness-init`). **Datum:** 2026-07-17.

---

## 1. Ziel

d-check verifiziert **Struktur**: Links lösen auf, Anker existieren, Kennungen
sind klickbar, die Referenz-Richtung stimmt. Es verifiziert **keine Behauptung**.
Eine `datei:zeile`-Angabe kann falsch sein, und jeder Gate-Lauf bleibt grün.

Belege des Adopters aus einem realen Planungs-Zug (2026-07-16/17, alle
nachgemessen):

| Behauptung im Dokument | Wirklichkeit | Gate |
|---|---|---|
| „54 Dateien" im Baseline-ZIP | 42 | grün |
| `modul-02:173-176` als Zitat-Beleg | Zitat endet bei 175 | grün |
| `modul-05:79-80` als Zitat-Beleg | beginnt bei 78 | grün |
| „~200 auf 120 Zeilen geschrumpft" | 219 auf 120 | grün |
| „193 Zeilen, fast alle Tag-Bumps" | 150 von 193 (77 %) | grün |
| „zwei von drei Linsen" | eine von drei | grün |

Die letzte Zeile ist die aussagekräftigste: sie stand **im Review-Report über
genau diese Fehlerklasse**, geschrieben nachdem die Klasse fünfmal aufgefallen
war. Sorgfalt hat dreimal nicht gereicht — Vorsession, drei Reviewer-Läufe, Autor
des Reports. **Das ist das Argument für einen Sensor statt für mehr Disziplin.**

**Warum jetzt.** Das Regelwerk schreibt Adoptern vor, die Baseline **committet
vendored** zu führen. Damit entsteht in jedem konformen Repo ein Korpus von
Zitaten auf einen in-tree, versionierten Fremdbaum. Diese Zitate sind heute
korrekt; **beim nächsten Tag-Bump zeigen sie still ins Leere.** Der Adopter misst
neun Stunden zwischen zwei Tags und verschobene Zeilen in 35 von 42 Dateien.
Zugleich ist das Vendoring die Bedingung, die den Check **erst möglich** macht:
in-tree ist die Referenz mechanisch prüfbar, ein gitignorierter Fetch-Cache war es
nicht. Die Vorgabe schafft die Lücke und den Hebel gleichzeitig.

## 2. Entscheidungen / Regel

**Keine.** Der eingereichte Vorschlag steht; die Erdung gegen unsere Artefakte
(§2.2) verschiebt aber seinen Zuschnitt, und die Vorfragen aus §4 sind offen.

### 2.1 Der eingereichte Vorschlag

Neues Modul `citations` in drei unabhängig aktivierbaren Stufen:

- **Stufe 1 `exists`** — `<pfad>:<zeile>` bzw. `<pfad>:<von>-<bis>` muss auf eine
  existierende Datei mit mindestens `<bis>` Zeilen zeigen. Fängt Zitat-Fäule nach
  Tag-Bump. Ambiguität: keine.
- **Stufe 2 `in-range`** — `<von> ≤ <bis> ≤ Zeilenzahl`. Fängt Tippfehler und
  invertierte Bereiche. Ambiguität: keine.
- **Stufe 3 `verbatim`** (opt-in, Markup-Pflicht) — ein ausgezeichneter Zitatblock
  wird **zeichengenau** gegen die referenzierte Spanne geprüft. Fängt
  Off-by-one-Bereiche, driftende Zitate, paraphrasierte „wortgleich"-Behauptungen.

Stufe 3 ist der eigentliche Gewinn: sie macht „wortgleich" von einer **Zusage** zu
einem **gemessenen Property**. Der Adopter hat genau diese Zusage schon einmal
gebrochen — ein als „wortgleich" deklarierter Regelwerks-Cache war eine selbst
geschriebene Kurzfassung. Kein Gate fand es; ein Mensch fand es.

**Der CR benennt seine Grenzen selbst, und das ist seine Stärke:** freie Zahlen
mit externer Grundwahrheit („42 Dateien im ZIP") und Prosa-Quantoren („fast alle",
„zwei von drei") bleiben ungeprüft. Von seinen sechs Belegzeilen fängt der
Vorschlag **zwei**. Er schließt die Lücke nicht — er schneidet den
mechanisierbaren Teil heraus und benennt den Rest als Review-Territorium.

### 2.2 Erdung gegen d-checks eigene Artefakte (Antwort auf §7 des CR)

Der CR stellt vier Fragen, die er selbst nicht beantworten kann. Gemessen
2026-07-17:

**Vorbesteht die Fähigkeit? — Ja, größtenteils. Der CR wäre für Stufe 1/2 eine
Dopplung.** `codepaths` scannt Inline-Code-Spans auf explizite Pfade, **erkennt
die `datei:zeile`-Konvention bereits** (`trimLineSuffix` in
`internal/hexagon/core/rules/codepaths.go`) und **verwirft die Zeilennummer
bewusst** — normativ in
[`DC-FA-CODE-001.a`](../../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)
Schritt 3 („ein Zeilen-Suffix `:NNN` abtrennen"). Es hat zudem `exempt-paths` mit
**exakt** der Begründung, die der CR für `citations.exempt-paths` gibt:
Review-Reports zitieren naturgemäß `Datei:Zeile`. Stufe 1/2 sind damit keine neue
Fähigkeit, sondern **das Weglassen eines `trim`** plus Zeilenzahl-Prüfung. Ein
eigenes Modul duplizierte Detektion, `roots` und Ventil.

**Aber: `codepaths` sieht nur Inline-Code, der CR spricht von Fließtext.** Ein
`modul-02:173-176` ohne Backticks ist für d-check heute unsichtbar — und das ist
**Absicht**: der Inline-Code-Marker *ist* das Absichtssignal des Autors
([`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in): „Explizite Pfade in **Inline-Code**"). Prosa-Scanning wäre eine
andere, deutlich fehleranfälligere Fläche. **Das ist die schärfste Rückfrage an
den Adopter:** stehen seine Zitate in Inline-Code? Dann ist der CR eine
`codepaths`-Erweiterung. Stehen sie in nackter Prosa, ist er ein Bruch mit dem
Absichtssignal — und braucht dafür eine eigene Begründung.

**Lastenheft-Bezug?** Es gibt **keine** Dach-Anforderung „Referenzen sind
verifiziert"; jedes Modul trägt seine eigene `DC-FA-*`. Siehe §4.

**ADR-Pflicht?** **Ja.** Jedes Modul hat eine ADR
([ADR-0020](../../adr/0020-content-pin-fence-ausnahme.md) `pins`,
[ADR-0024](../../adr/0024-vcs-immutable-gate.md) `vcs`,
[ADR-0030](../../adr/0030-tracked-referenz-ziele.md) `tracked`). d-check hält das
strenger als die adoptierte Lehre.

**Zitat-Syntax `d-check:cite`?** **Zurückstellen.** d-checks bestehende
Direktiven-Konvention (`d-check:ignore`) blendet in einer Tabellenzeile den
eigenen Tabellen-Reader; der Defekt ist **ausgeliefert und nach fünf Anläufen
ungelöst** ([slice-074](../done/slice-074-kommentar-suffix-tabellenzeilen.md)). Eine
**zweite** Direktive einzuführen, während die Platzierungsregeln der ersten offen
sind, vergrößert eine ungelöste Klasse. Stufe 3 hängt daran — Stufe 1/2 nicht.

**Dogfood-Reichweite bei uns — nahe null:**

| | Anzahl |
|---|---|
| `datei:zeile`-Zitate in Inline-Code (`docs`+`spec`+`harness`) | 837 |
| davon in `docs/reviews/` — bereits via `exempt-paths` ausgenommen (Zeitdokumente) | 830 |
| **verbleibend prüfbar** | **7** |

Kein Gegenargument, aber eine **Warnung**: wir können das Gate an uns selbst kaum
erproben. Genau diese Dogfood-Lücke hat 2026-07-17 dreimal zugeschlagen
(slice-073 §4, slice-074 §4, [slice-076](../done/slice-076-markdown-lexik-commonmark.md) §4).
Der Realdatenbeleg gegen das Adopter-Repo ist daher **nicht optional**.

## 3. Definition of Done

- [x] **CR erfasst**, §7 gegen unsere Artefakte beantwortet (§2.2).
- [x] **Rückfrage an den Adopter beantwortet:** Inline-Code (gemessen: 33/33
  `datei:zeile`-Zitate in `ai-harness-init` backticked, null Prosa) ⇒
  `codepaths`-Erweiterung.
- [x] **Vorfragen entschieden** (§4) — Zuschnitt: alle drei Stufen jetzt, Form (c)
  (Stufe 1/2 `codepaths`-Erweiterung, Stufe 3 eigenes Modul); Verortung bereits
  aufgelöst (Kürzel-Kriterium mit slice-078).
- [x] **Lastenheft-CR:** das Modul `codepaths` erweitert (opt-in
  `codepaths.check-lines`, 2 neue AKs;
  [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in))
  **+** neue Anforderung
  [`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in)
  (18. Modul `citations`, Bereich `CITE` in §3, 3 AKs); `citations` in Modul-Liste +
  Glossar, Version 0.47.0→0.48.0, Historie.
- [x] **ADR:** [ADR-0045](../../adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)
  (Proposed) — Erweiterung vs. eigenes Modul begründet (Existenz/Bereich → `codepaths`,
  Inhalt → eigenes Modul wie `anchors`↔`links`), `d-check:cite`-Direktive,
  zeichengenau, fail-closed, nur ausgezeichnete Blöcke.
- [x] **Spezifikation-`.a`:**
  [`DC-FA-CODE-001.a`](../../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)
  Schritt 3/6 (Zeilen-Check) + neue
  [`DC-FA-CITE-001.a`](../../../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
  (Direktive → Zitatblock → zeichengenauer Vergleich) + §2-Schema
  (`codepaths.check-lines`). Grund-Codes (§4) folgen mit der Implementierung (Lockstep).
- [ ] **Tests:** Zitat auf zu kurze Datei ⇒ Befund · invertierter Bereich ⇒ Befund
  · `exempt-paths` greift · **Negativ: ein korrektes Zitat bleibt grün nach einem
  Tag-Bump, der die Datei nicht anfasst** (kein Fehlalarm durch Nachbar-Drift).
- [ ] **Realdatenbeleg gegen das Adopter-Repo** — **nicht optional** (§2.2:
  Dogfood-Reichweite 7).
- [ ] **Qualität:** unabhängiger, kontext-getrennter Review **vor** dem Release;
  `make gates`/`make ci` grün.

## 4. Risiken / offene Punkte

- **Vorfrage 1 — Zuschnitt: ENTSCHIEDEN (Auftraggeber, 2026-07-18): Form (c), alle
  drei Stufen jetzt.** Die Adopter-Rückfrage ist empirisch beantwortet (33/33 Zitate
  in Inline-Code) ⇒ `codepaths`-Erweiterung, kein Prosa-Scanning. **Stufe 1/2** sind
  eine Erweiterung von
  [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
  (die Zeile wird heute erkannt und verworfen), **kein neues Modul**. **Stufe 3**
  (`verbatim`) ist genuin neu (Inhalts-Verifikation statt Existenz) und wird als
  **eigenes Modul** gebaut, das `codepaths`' Detektion **aufgreift** (kein
  Re-Detect — wie `anchors` auf `links` aufsetzt), mit eigener Direktive
  `d-check:cite` (durch slice-074 entblockt). — _Verworfen: (a) Stufe 3 später
  (Auftraggeber will alle drei jetzt); (b) eigenes Modul, das die Detektion
  **dupliziert** (§2.2)._
- **Vorfrage 2 — Lastenheft-Verortung: AUFGELÖST (Auftraggeber, 2026-07-18)** über ein
  **gemeinsames Kürzel-Kriterium** mit
  [slice-078](../done/slice-078-ignore-refs-quell-skopus.md):
  **querschnittlich (mehrere Module teilen die Fähigkeit) → neues Bereichskürzel;
  Einzelmodul-Erweiterung → bestehende Anforderung ändern**
  ([`MR-002`](../../../../harness/conventions.md#mr-002--id-schema-mit-bereichskürzeln-ab-initialer-fassung)).
  Anwendung hier: Stufe 1/2 (`datei:zeile`-Verifikation) sind eine
  Einzelmodul-Erweiterung von `codepaths` (die Zeile wird dort schon erkannt und
  verworfen, §2.2) → **Änderung von
  [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)**,
  **kein** neues Kürzel; Stufe 3 (`verbatim`) bleibt separat (eigene Direktive,
  zurückgestellt). So driften 078 (querschnittlich → neues Kürzel) und 079
  (Einzelmodul → bestehende Anforderung) **nicht** — ein Prinzip, zwei
  situationsgerechte Ergebnisse.
- **Zweite Direktive bei ungelöster erster.** Siehe §2.2. Stufe 3 sollte warten,
  bis slice-074 eine tragende Regel hat.
- **Dogfood-Lücke.** Reichweite bei uns: 7 Referenzen. Wir können den Sensor an
  uns selbst nicht erproben.
- **Der Vorschlag fängt zwei von sechs Belegzeilen** — das steht so im CR, und es
  ist seine Stärke, nicht seine Schwäche. Ein überversprochenes Gate wäre
  schlimmer als keines. Die verworfene Alternative `claims` (Provenienz-Pflicht
  für Zahlen) gehört ausdrücklich **nicht** hierher: höheres
  False-Positive-Risiko, unscharfe Abgrenzung, und sie sollte auf
  Betriebserfahrung mit dieser Stufe aufsetzen statt darauf zu wetten.
- **Quell-Artefakt beschädigt.** Der eingereichte CR-Text lag doppelt vor (§2–§7
  zweimal, teils mitten im Satz abgeschnitten). Inhaltlich rekonstruiert; vor der
  Annahme gegen das Original des Adopters abgleichen.

## 5. Trigger

Adopter-CR `ai-harness-init` (2026-07-17), Beleg-Artefakt: dessen
Plan-Review-Report vom selben Tag. Anlass ist strukturell, nicht rückblickend: die
Vorgabe, die Baseline committet vendored zu führen, erzeugt in jedem konformen
Repo einen Korpus von `datei:zeile`-Zitaten auf einen versionierten Fremdbaum, der
beim nächsten Tag-Bump **still** verfault.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Die Spezifikation führt, der Code folgt. Hier verschärft — das
**Lastenheft** führt: der CR ändert oder ergänzt einen Vertrag, und dieser Slice
hält an §4 an, statt Code zu schreiben.

## 7. Closure-Notiz (nach `done/`)

_Ausstehend._
