# Slice slice-285: Tabellenzählung — Schätzung aus der Tabelle des Plans, Prüfung 2 und Hinweis „Marke ausstehend“

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-93

**Bezug:** die Anforderung des Moduls `metrics` aus [`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) (die Hälfte „Schätzung aus der
Tabelle“),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) (gleiche Tabelle ⇒ gleiche Zählung),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (lesend, hermetisch). Die ADR aus
[`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) ist die aktive Entscheidungsgrundlage (Entscheid 9).

**Berührte Spec-Stellen:** die in [`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) angelegten Stellen (Anforderung,
Schema `metrics.*` für das Tabellenformat, Eingabefehler, Hinweise) — gelesen, nicht
verändert, außer dem Fassungsstand der Historie.
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-30.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Modul `metrics` liest aus jedem Plan mit Block die eine
Markdown-Tabelle (konfigurierbarer Kopf, Schlüsselspalte, Klassenspalte), zählt die
Zeilen je Klasse und bildet die Schätzung aus Klassen-Faktoren und Block-Feldern.
Damit prüft es „Schätzung gleich Tabelle“ (nur solange die Messung 0 ist) und gibt
bei einem Block ohne Schätzung und Marke den Hinweis „Marke ausstehend“ mit der
berechneten Schätzung und Marke. Tabellenfehler sind Eingabefehler (Exit 2).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Block, Marke gegen Formel, leerer Grund, Messung** — Bestand von
  [`slice-283`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) und [`slice-284`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle), davor `done/`: Prüfung 2 braucht den Zustand
  „Messung 0 / ≥ 1“ aus der Messung; dieser Slice nimmt ihn als gegeben.
- **Die Prosa-Form-Erkennung** (ältere Zeile statt Block) — nach der Antwort des Issue-Autors
  kein Zustand des Moduls: sie kommt nur in Plänen vor, die das Verzeichnis „in Arbeit“
  nicht enthält; ein Plan ohne Block gibt nur den Hinweis „kein Block“ mit dem Plan-Namen.
- **Andere Tabellen-Auswertungen als die Zählung je Klasse** (Summen aus Zellen,
  mehrere Tabellen) — der Anlass verlangt Zeilen je Klasse; mehr wäre eine neue
  Anforderung.
- **Release-Prep (Handbuch-Kopf, CHANGELOG, Digest-Pins)** — ein anderer Vorgang
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Config-Schlüssel und Zustand:** `metrics.estimate-inputs` (Feld → Faktor, Faktor ohne Feld ⇒ Exit 2) und `metrics.table.*` kommen in den Config-Rand; der Zustand *Marke ausstehend* verlangt alle Felder aus `estimate-inputs` und nennt im Hinweis die berechnete Schätzung und Marke.
- [ ] **Tabellen-Leser und Konfiguration:** genau eine Tabelle pro Plan außerhalb
      umzäunter Blöcke; Kopf, Schlüsselspalte, Klassenspalte und Klassen samt Faktor sind
      Konfiguration; Trennzeile, Zellen (`\|`, fehlender Schlussstrich, CRLF),
      Tabellenende und Zählung je Klasse nach der Anforderung. Eingabefehler (Exit 2,
      Klartext): Tabelle fehlt · mehrdeutig · Trennzeile fehlt · Zeile mit falscher
      Zellenzahl · `Nr.` leer oder doppelt · Klasse unbekannt · Tabelle ohne Zeilen.
      Config-Rand (leere Klassenliste, Klasse ohne Faktor, Schlüsselspalte nicht im Kopf)
      → Exit 2. Negativliste (§3) getestet.
- [ ] **Schätzung und Prüfung 2:** `S = Σ Klassenfaktor · Zeilen + Σ Feldfaktor · Feld`
      (ganzzahlig); solange die Messung 0 ist, meldet ein Unterschied zu `schaetzung` im
      Block den Befund „Schätzung gleich Tabelle“ verfehlt, bei Messung ≥ 1 nur einen
      Hinweis. Die belegten Zahlen der Anforderung (Slices A, B, C und das durchgerechnete
      Beispiel) sind Testfälle.
- [ ] **Hinweis „Marke ausstehend“:** ein Block mit Pfad und Eingaben, aber ohne Basis,
      Schätzung und Marke, gibt die berechnete Schätzung und Marke als Hinweis (kein
      Befund, Exit unverändert).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update nur, wo ein Test es hält; Handbuch-Kopf und `CHANGELOG.md` **nicht**
- [ ] Grund-Code `metrics-estimate-mismatch`: Zeile in der Tabelle von `spec/spezifikation.md` §4, Eintrag in `AllReasons()` und `--doctor`-Klartext (`TestAllReasonsDeckungGegenSpezifikationGrundCodes`).
      (Release-Prep, §1).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `metrics.go` im Regel-Paket | update | Tabelle finden und zählen, Schätzung bilden, Prüfung 2 und Hinweis |
| `metrics_table.go` im Regel-Paket | neu (falls `metrics.go` zu lang wird) | Tabellen-Leser getrennt von der Prüfung |
| `internal/hexagon/core/model/config.go` | update | Tabellenformat und Klassenfaktoren in `MetricsConfig` |
| `internal/adapter/driven/configyaml/configyaml.go` | update | `metrics.tabelle.*` und Klassen lesen und validieren |
| `internal/adapter/driving/cli/config_template.go` | update | Gerüst führt das Tabellenformat (Spiegel-Test) |
| `metrics_test.go` im Regel-Paket · `configyaml_test.go` | update | Happy/Boundary/Negative nach der Anforderung; Negativliste unten |
| `spec/spezifikation.md` | update | Zeile `metrics-estimate-mismatch` in der Tabelle von §4, Tabellenende bei Anfangsstrich, Historie |
| `internal/hexagon/core/app/diagnose.go` | update | Grund-Code in `AllReasons()` samt Klartext |
| `tools/image-test.sh` | update | Fall (7): Schätzung aus der Tabelle am gebauten Image (Befund bei Messung 0, Hinweis danach); Fall (6) bekommt das Tabellenformat in der Konfiguration |
| `docs/user/abdeckung-e2e.md` | update | Abdeckung per `make abdeckung` |

**Negativliste der Tabellen-Erkennung (Schritt 19, vor dem Code-Commit zu testen):**
Tabelle mit ähnlichem Kopf, aber anderer Zellenzahl · zweite Tabelle mit passendem Kopf
(mehrdeutig) · Tabelle im umzäunten Block (Backticks und Tilden, Einrückung 0 bis 3
Leerzeichen, Zaun mit Backtick im Rest der Zeile) · Zeile mit `\|` in einer Zelle ·
`|` im Code-Span (kein Fehler: der Span ist eine Zelle, wie bei den übrigen Tabellen-Lesern des Produkts) · fehlender Schlussstrich ·
Leerraum hinter dem Schlussstrich · Tabelle endet an Leerzeile, an Text, an `>`-Zeile, an
Zaun · späteres `|`-Stück nach dem Ende gehört nicht dazu · Trennzeile mit `:` an einem
und beiden Enden, mit anderer Zellenzahl · `Nr.` leer und doppelt · Klasse leer,
`nebenlaeufig`, `Nebenläufig` · Tabelle ohne Datenzeile · CRLF · Plan ohne Block
(Tabelle wird nicht gelesen). Dazu der Lauf über den eigenen Bestand: mit aktiviertem
Modul und ohne Block **kein** Befund, gegen das Vorher-Image byte-gleich.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-284`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) liegt in `done/` (Messung und der
Zustand „Messung 0 / ≥ 1“ stehen), und welle-93 ist offen. Der Übergang landet auf dem
Hauptzweig, vor der Arbeit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Tabellen-Leser und Prüfung
  tragen zusammen mehr als eine Review-Sitzung — dann wandert der Leser in einen eigenen
  Slice davor, ausdrücklich nicht still weiter ([`MR-066`](../../../../harness/conventions.md#mr-066)).
- `in-progress` → `open` (blockiert — Carveout?): die Tabellenformate weiterer Pläne des
  Anlass-Repos weichen vom beschriebenen ab (Rückfrage im Issue).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Review-Report liegt vor, und am gebauten Image
meldet ein Fixture mit abweichender Schätzung bei Messung 0 den Befund, bei Messung ≥ 1
den Hinweis; ein Block ohne Schätzung und Marke gibt den Hinweis mit den berechneten
Zahlen. Der Lauf ohne `metrics`-Konfiguration ist gegen das Vorher-Image byte-gleich.
Dazu die Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Tabellen-Erkennung liest einen Zaun oder eine Zeile anders als der
  Markdown-Leser des Repos; eine Tabelle verschwindet still oder wird mehrdeutig —
  **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)
- Das Tabellenformat ist repo-spezifisch; ein zweiter Konsument bräuchte ein anderes
  Format, das die Konfiguration nicht ausdrückt — **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)
- Die Paritäts-Probe vergleicht Exit-Code und Inhalt der Befunde und Hinweise, nicht den
  Wortlaut; eine Abweichung im Inhalt zeigt sich erst beim Eigentümer — **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)

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

Wird bei der Closure gefüllt.

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
Regelmodul, Konfigurations-Rand, Gerüst; keine Sub-Area mit eigener Modus-Deklaration
(`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen, Sub-Area
`*`. Treffer:
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(7×, verkörpert in Schritt 19 des Workflow-Skeletts — die Negativliste in §3 ist seine
Anwendung; die Tabellen-Erkennung ist eine neue Erkennungsform) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(3×, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — die Tabelle liegt im selben
Plan, den das Modul liest; die Grenze steht in der Anforderung). Der Zähler folgt aus
den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-10-30: `upstream-drift` grün (jüngster Lauf 2026-10-30T06:41:15Z), `image-scan`
grün (2026-10-30T10:04:41Z). Nichts zu lesen, nichts ausstehend.

### Sub-Area: `*` (Default)

Alle berührten Sub-Areas GF.
