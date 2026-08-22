# Welle welle-81-zustandsfelder: Zustandsfelder tragen Zustand, keine Chronik (Baseline v5.9.0)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-81-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Adoptions-Pflege der Baseline plus
die Konsequenzen, die sie am eigenen Bestand fällig macht).

**Verantwortlich:** pt9912. **Datum:** 2026-08-22.

---

## 1. Welle-Ziel

Die Baseline von `v5.7.0` auf **`v5.9.0`** heben (zwei Stufen, Kurs-Tags vom
2026-08-22) und die eine Regel ziehen, die sie bringt. Gemessen statt
hergeleitet: von 33 geänderten Bundle-Dateien tragen **26 nur den
Quell-URL-Stempel**; der Inhalt ändert sich in vier Regelwerks- und sieben
Template-Dateien, und alle sagen dasselbe.

**Die Regel:** Ein Feld, das einen *Zustand* trägt — die `Stand`-Zelle einer
Register-Zeile, die `Status`-Spalte eines Meilensteins, die Kopfzeile eines
lebenden Registers —, ist ein Zustands-Artefakt wie ein Kommentar, nur im
Rumpf. Es nennt **den Zustand und den Beleg als auflösbaren Anker**, nicht die
Chronik, wie der Zustand zustande kam. Was sonst in der Zelle stand, hat
eigene Orte: Behauptung und vorgeschlagene Handlung beim Vorhaben selbst, die
**Schließung im Closure-Log**, die **Umplanung im Drift-Log** — und was keines
davon ist, hält `git`. Ein Drift-Log, das Schließungen protokolliert, ist ein
zweites Closure-Log, und zwei Logs driften.

**Vier Treffer im eigenen Bestand, alle gemessen:**

| Fläche | Ist-Zustand | Ziel-Form |
|---|---|---|
| Kopfzeilen lebender Register | drei Dateien tragen `**Status:** Aktiv. **Letzte Änderung:** …` | ersatzlos weg — der Zustand eines Registers ist sein Inhalt, sein Änderungsdatum hält `git` |
| Kopf des Technik-Stratums | dieselbe Zeile in der Spezifikation | weg — die Historie trägt das Datum, zwei Felder für eines driften |
| Kopf der Sicht | dieselbe Zeile in der Architektur | **bleibt** — dort ist sie der bewusste Frische-Marker (die Sicht hat keine Historie) |
| `Stand`-Zellen des Beobachtungs-Registers | acht Zellen, 169 bis **3 011** Zeichen, überwiegend Chronik | Zustand + Beleg als Anker |
| Drift-Log der Roadmap | 71 Zeilen, die obersten sind **Schließungen** („welle-79 geschlossen", „slice-111 abgeschlossen") | nur Umplanungen; Schließungen stehen im Closure-Log |

## 2. Trigger (Welle startet)

Auftraggeber-Anstoß 2026-08-22 („danach können wir auf das neue Regelwerk
v5.9.0 migrieren"); beide Kurs-Tags sind verifiziert und das Currency-Audit
meldet sie (`--check-latest`, Exit 3), der gepinnte Stand ist upstream
unverändert (kein Content-Drift). welle-80 ist geschlossen.

## 3. Closure-Trigger (Welle schließt)

- Alle vier Slices in `done/`; `make fullbuild` grün (Exit explizit).
- Der Pin steht auf `v5.9.0`, der alte Baum ist entfernt, `--verify` grün und
  `--check-latest` meldet keinen neueren Release mehr.
- Kein lebendes Register und kein Technik-Stratum trägt eine
  Kopf-Zustandszeile; die Sicht trägt ihre weiter, und die Begründung steht
  dort.
- Keine `Stand`-Zelle und keine Drift-Log-Zeile erzählt eine Chronik —
  gemessen an der eigenen Datei, nicht behauptet.
- Ergebnisnotiz `welle-81-results.md` mit Register-Lese-Schritt.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [slice-117](in-progress/slice-117-baseline-v590-bump.md) | Pin-Hebung auf `v5.9.0`: Bundle vendored, alter Baum entfernt, Nachfolge-Adaption, Drei-Klassen-Zensus der pin-gebundenen Verweise | §Baseline, [`MR-028`](../../../harness/conventions.md#mr-028) |
| [slice-118](open/slice-118-zustandsfeld-regel.md) | Die Regel verkörpern, bevor sie angewandt wird: Briefing-Hard-Rule auf Zustandsfelder erweitert, Reviewer-Skill mit dem HIGH-Anker | [`MR-000`](../../../harness/conventions.md#mr-000--baseline-aussage), Reviewer-Skill |
| [slice-119](open/slice-119-kopf-zustandsfelder.md) | Die drei Kopf-Zustandszeilen entfernen; die Sicht behält ihre samt Begründung | Spezifikation, Roadmap, Beobachtungs-Register |
| [slice-120](open/slice-120-register-und-drift-log.md) | `Stand`-Zellen auf Zustand + Beleg, Drift-Log auf Umplanungen zurückschneiden, Meilenstein-Status-Form | Beobachtungs-Register, Roadmap |

## 5. Abhängigkeiten

- Blockiert: nichts Geplantes.
- Wird blockiert von: nichts. Reihenfolge innerhalb: 117 zuerst (die Regel
  muss vendored sein, bevor sie zitiert wird), dann 118 (Regel verkörpern),
  dann 119 und 120 (anwenden); 119 und 120 sind voneinander unabhängig.

## 6. Out-of-Scope für diese Welle

- **Kein Produkt-Code, kein Release.** Die Regel ist eine Doku-/Prozess-Regel;
  kein Modul und kein Grund-Code ändert sich.
- **Kein neuer Sensor.** Die Baseline sagt ausdrücklich, dass kein Gate das
  fängt — Träger sind Briefing und Reviewer-Skill. Ob ein `structure`-Muster
  einen Teil davon prüfen könnte, ist eine eigene Frage.
- **Der Vertrag bleibt unberührt:** das Lastenheft trägt Version und Status
  weiter (die Vorlage führt beide unverändert) — die Regel zielt auf lebende
  Register und das Technik-Stratum.
- **Kein Rückbau der Historie:** was aus einer `Stand`-Zelle verschwindet, ist
  nicht verloren — es steht in `git`, in den Closure-Notizen und in den
  Review-Reports. Nichts davon wird nachgetragen oder umgeschrieben.
