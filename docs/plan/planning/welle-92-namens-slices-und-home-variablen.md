# Welle welle-92: Namens-Slices im Trace und Home-Variablen in hostpaths

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-10-29.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Zwei im Feld gemessene Erkennungslücken sind geschlossen, ohne dass ein bisher
grüner Lauf ungefragt rot wird:

- Der Trace (`--trace`, [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
  ordnet einen Slice mit Namens-Kennung (`slice-<name>.md`) ohne Konfiguration
  einer Anforderung zu; ein Repo mit Namens-Slices endet bei
  `--trace --require-complete` nicht mehr wegen unsichtbarer Slices mit Exit 1
  (GitHub-Issue 11).
- Das Modul `hostpaths` ([`DC-FA-HOST-001`](../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in))
  erkennt `$HOME/…` und `${HOME}/…` auf Wunsch als Home-Verweis, auch in
  Inline-Code; ohne den neuen Schlüssel bleibt der Befundsatz byte-identisch
  (GitHub-Issue 12).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Die GitHub-Issues 11 und 12 liegen offen vor, und der Auftraggeber hat ihre
  Bearbeitung als Welle angefordert (2026-10-29).
- Für Issue 11 ist der Weg entschieden: Default erweitern, der Ziffern-Fall
  bleibt unverändert (Auftraggeber-Entscheid 2026-10-29).
- Die Lifecycle-Verzeichnisse `open/` bis `in-progress/` tragen keinen Slice
  außer den zwei dieser Welle.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle sind `done/`.
- `make fullbuild` grün.
- **Das Mehr gegenüber den Slice-DoDs** (repo-weit, in keiner einzelnen DoD):
  `make blackbox-probe` gegen das Vorher-Image `v0.87.1` zeigt **dieselbe**
  Ausgabe (stdout, stderr und Exit getrennt) für zwei Eingaben, die der
  Eigenlauf allein nicht abdeckt — dieses Repo setzt `trace.slices.file-pattern`
  selbst, ein Lauf gegen es prüft den neuen Default also nicht: eine Kopie
  des Planning-Baums **ohne** `trace`-Block (Ziffern-Slices, Beobachtungs-Belege
  unter `observations/`) und ein Fixture mit `hostpaths` ohne den neuen
  Schlüssel. Dazu liefern beide Fixtures aus den Issues (Namens-Slice,
  `$HOME`-Pfad) am gebauten Image das erwartete Ergebnis. Jeder Slice belegt
  seine Hälfte; dass keine die andere verschoben hat, belegt erst der
  gemeinsame Lauf.
- Das Trigger-Audit (seit slice-241) ist über die Welle-Closure vollzogen und
  in `welle-92-results.md` belegt.
- Beide Issues tragen eine Antwort mit Verweis auf das Ergebnis; Schließen und
  Release-Prep sind Folge-Handlungen des Auftraggebers (siehe §6).

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-279`](open/slice-279-trace-erkennt-namens-slices.md) | Der Trace erkennt Slices mit Namens-Kennung ohne Konfiguration | [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix) |
| [`slice-280`](open/slice-280-hostpaths-home-variablen.md) | `hostpaths` erkennt `$HOME` und `${HOME}` als Home-Verweis (opt-in) | [`DC-FA-HOST-001`](../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine Welle.
- Wird blockiert von: keiner.
- **Reihenfolge der Slices:** erst [`slice-279`](open/slice-279-trace-erkennt-namens-slices.md),
  dann [`slice-280`](open/slice-280-hostpaths-home-variablen.md). Beide
  schreiben eine Fassung ins Lastenheft (Versions-Zahl und Historie-Zeile in
  [`spec/lastenheft.md`](../../../spec/lastenheft.md)); das WIP-Limit 1 hält
  sie ohnehin nacheinander, die Reihenfolge verhindert den Versions-Konflikt.
  Fachlich sind sie unabhängig.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Release, Tag und Release-Prep** (CHANGELOG, READMEs, Handbuch-Kopf,
  Digest-Pins): eigener Vorgang nach der Welle ([`AGENTS.md`](../../../AGENTS.md)
  §5 Zeile 17) — die Welle liefert Code, Spec und Handbuch-Inhalt, keine
  Distribution.
- **Die Erkennung von Namens-Kennungen außerhalb des Trace**: die
  Review-Deckung (`reviews`), die Konfigurations-Vorschläge
  (`--suggest-config`) und `tools/archive-wave` tragen dieselbe ziffern-only
  Form. Der Gegenstand ist ein anderer Vorgang (andere Module, andere
  Anforderungen), und der Befund des Feldes betraf nur den Trace. Was der
  Slice über die Spiegel feststellt, steht dort im Plan; Folgearbeit wäre ein
  eigener Slice.
- **Weitere Variablen-Formen** (`$USERPROFILE`, `%USERPROFILE%`, `$XDG_*`):
  nicht gemeldet, nicht gemessen; der Schlüssel aus
  [`slice-280`](open/slice-280-hostpaths-home-variablen.md) ist so geschnitten,
  dass eine weitere Variable ein Eintrag in seiner Liste wäre, aber keine
  Anforderung dieser Welle.
- **Ein Zeilen-Marker für `hostpaths`**: bleibt ausgeschlossen
  ([`DC-FA-HOST-001`](../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in));
  die Issues verlangen ihn nicht.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: bei der Closure — Zeiger auf `welle-92-results.md`, Geschwister im Ruheort `done/`.
Zähler: bei der Closure — Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort.
