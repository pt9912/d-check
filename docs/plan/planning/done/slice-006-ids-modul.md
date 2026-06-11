# Slice slice-006: Modul `ids` — Linkpflicht für Kennungen

**Status:** done.

**Welle:** welle-03-regelmodule.

**Bezug:** [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(Constraint `ids.patterns[].target` muss existieren),
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl);
[ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md).

**Autor:** pt9912. **Datum:** 2026-06-10.

---

## 1. Ziel

Das Regelmodul `ids` ist implementiert: nackte Kennungen im Fließtext
(konfigurierte Regex-Muster) müssen Markdown-Links auf ihre
Definition sein.

## 2. Definition of Done

- [x] Akzeptanzkriterien von `DC-FA-ID-001` als Tests: ID als Link →
  kein Befund; ID in Inline-Code → linkpflichtfrei; nackte ID →
  `id-unlinked` (Grund-Code gemäß Spezifikation §4).
- [x] Muster-Präzedenz getestet (Deklarationsreihenfolge, erstes
  Match gewinnt — Spezifikation §`DC-FA-ID-001.a`).
- [x] „Verlinkt" = Vorkommen liegt im Linktext eines Markdown-Links
  (Link-Text-Spannen aus der Extraktion).
- [x] Config-Constraint durchgesetzt: nicht existierendes
  `ids.patterns[].target` → Exit 2 (mit Test).
- [x] `ids` in `isImplemented` aufgenommen (Interim-Hinweis entfällt).
- [x] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md)
  aktualisiert; Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/core/ids.go` (+ Tests) | neu | Prüf-Logik des Moduls |
| `internal/hexagon/core/markdown.go` | update | Link-Text-Spannen pro Zeile exportieren (für „in Linktext"-Erkennung) |
| `internal/hexagon/core/run.go`, `config.go` | update | Modul-Registrierung; Target-Existenz beim Lauf-Start |
| Akzeptanztests (cli_test) | update | Black-Box-Abdeckung |

## 4. Trigger

Sofort — welle-03 ist aktiv.

## 5. Closure-Trigger

DoD vollständig + Commit(s) auf `main` + Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Überlappende Treffer im selben Textbereich (Muster-Präzedenz pro
  Vorkommen) — Spezifikation ist eindeutig, Implementierung braucht
  saubere Span-Arithmetik.
- Die Selbstkonfiguration (eigene `DC-*`/`MR-*`/`ADR-*`-Muster)
  erfolgt bewusst erst in slice-007 zusammen mit `matrix` — nackte
  IDs in immutablen ADR-/MR-Texten brauchen die dort geplante
  Ausnahme-Betrachtung.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Commit `6add093` (Modul `ids`, Link-Spannen-Export,
Target-Constraint, Spec-Präzisierung, Tests).

- **Was hat funktioniert:** Die Span-Arithmetik aus dem Risiko-Block
  (überlappende Treffer) ließ sich klein lösen — pro Zeile eine
  `claimed`-Intervallliste, erstes Muster gewinnt; beide
  Deklarationsreihenfolgen sind getestet. Die Link-Text-Spannen kamen
  als Refactoring von `parseLinkAt` ohne Verhaltensänderung für
  `links`/`anchors` heraus (alle Bestandstests blieben grün).
- **Anders als geplant:** Eine Spec-Lücke wurde sichtbar:
  §`DC-FA-ID-001.a` hätte wörtlich genommen auch Vorkommen in der
  Ziel-Klammer von Links und in Bildreferenzen als `id-unlinked`
  geflaggt — in Repos, deren Definitions-Dateinamen die Kennung
  tragen (`ADR-0042-beispiel.md`), würde damit jedes korrekt
  verlinkte Vorkommen einen False-Positive aus der Ziel-Klammer
  erzeugen. Die Spezifikation wurde fortgeschrieben (Ziel-Klammern
  und Bildreferenzen sind kein Fließtext); das Lastenheft („im
  Fließtext") blieb unberührt.
- **Steering-Loop-Lerneintrag:** (a) Spec-Lücken-Klasse *inferential
  feedforward* — die Algorithmus-Operationalisierung („sonst Befund")
  hatte den Negativ-Raum (was ist *kein* Fließtext?) nicht
  durchdekliniert; bei den kommenden Modulen `matrix`/`external` den
  „sonst"-Zweig der Spezifikation vor Implementierung explizit gegen
  Beispiel-Korpora prüfen. (b) Der Interim-Mechanismus
  (`isImplemented`/`SkippedModules`) hängt jetzt nur noch an
  `matrix`/`external` — nach slice-008 ist er toter Code und gehört
  entfernt (Hinweis in slice-008 einplanen).
- **Folge-Slices:** keine neuen; die Selbstkonfiguration (eigene
  `DC-*`/`MR-*`/`ADR-*`-Muster in `.d-check.yml`) folgt wie geplant
  in slice-007 (Ausnahme-Betrachtung für immutable ADR-/MR-Texte).

**Review R1 (nach Closure, Agent-Review mit getrenntem Kontext):**
6 Findings (2 MEDIUM, 3 LOW, 1 INFO), alle nachgeschärft in einem
Folge-Commit: Repo-Escape-Verbot für `ids.patterns[].target` *und*
`scan.roots` (gemeinsamer Helfer `resolveConfigPath` — die
Scan-Wurzel-Lücke war eine Bestandslücke derselben Klasse aus
slice-003, vom Review-Fix mit abgedeckt); Leerstring-matchende
ids-Regexe als Konfigurationsfehler; Inline-Code-Stripping
positionserhaltend (Leerzeichen statt Entfernen — keine
Phantom-Kennungen); Schleifen-Dopplung über gemeinsamen
Link-Iterator beseitigt; zeilenbasierte Extraktion als normative
Grenze in der Spezifikation dokumentiert (Konsequenz für `ids` als
Risiko-Punkt in slice-007 eingetragen). Steering-Loop-Beleg: das
Review mit getrenntem Kontext fand eine Fehlerklasse (Repo-Escape in
Config-Pfaden), die Implementierungs- und Selbstprüfungs-Kontext
beide übersehen hatten — Bestätigung der Rollentrennung.

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (spec-first; siehe Kurs Modul 5 §Worked
Mini-Example).
