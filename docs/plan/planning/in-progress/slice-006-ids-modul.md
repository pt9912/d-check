# Slice slice-006: Modul `ids` — Linkpflicht für Kennungen

**Status:** in-progress.

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

- [ ] Akzeptanzkriterien von `DC-FA-ID-001` als Tests: ID als Link →
  kein Befund; ID in Inline-Code → linkpflichtfrei; nackte ID →
  `id-unlinked` (Grund-Code gemäß Spezifikation §4).
- [ ] Muster-Präzedenz getestet (Deklarationsreihenfolge, erstes
  Match gewinnt — Spezifikation §DC-FA-ID-001.a).
- [ ] „Verlinkt" = Vorkommen liegt im Linktext eines Markdown-Links
  (Link-Text-Spannen aus der Extraktion).
- [ ] Config-Constraint durchgesetzt: nicht existierendes
  `ids.patterns[].target` → Exit 2 (mit Test).
- [ ] `ids` in `isImplemented` aufgenommen (Interim-Hinweis entfällt).
- [ ] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md)
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

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (spec-first; siehe Kurs Modul 5 §Worked
Mini-Example).
