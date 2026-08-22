# Welle welle-82-config-flaechen: Drei Konfigurations-Flächen, jede eine Kerbe zu schmal

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-82-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug; die Welle schließt mit einem
**Minor-Release**, weil alle drei Erweiterungen konsumentensichtbar sind.

**Verantwortlich:** pt9912. **Datum:** 2026-08-22.

---

## 1. Welle-Ziel

Drei Module tragen eine Konfigurations-Fläche, die **eine Kerbe zu schmal**
ist. Das ist keine Vermutung — jede der drei hat in den letzten beiden Wellen
einen realen Schaden angerichtet, und jeder ist dokumentiert:

| Modul | Heutige Enge | Realer Schaden |
|---|---|---|
| [`versions`](../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) | genau **ein** `pin-pattern` gegen **eine** `current-from`-Quelle | die 3×-Form der Beobachtung BEO-008 (Baseline-Tag in URLs und Prosa gegen den Pin) ist damit **nicht baubar** — die Klasse ist dreimal eingetreten und bleibt gate-blind |
| [`structure`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) | kein Schlüssel „**jede** Überschrift des Abschnitts matcht dieses Muster" | die Ersatz-Konstruktion (ausgeschriebene Präfix-Negation, weil RE2 keinen Lookahead kennt) hatte ein **stilles Falsch-Negativ**: eine eingerückte Sektion entkam der Kennungs-Pflicht |
| [`diagrams`](../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in) | **kein** Datei-Ventil, **kein** Zeilen-Marker, **keine** `§2`-Schema-Zeilen | ein Beispiel-Diagramm mit erfundener Kennung in einem Report hätte über den `pre-commit`-Hook jeden Commit blockiert; das eigene Profil musste auf `spec/` gescopt werden, um das zu umgehen |

**Alle drei Erweiterungen sind additiv und opt-in:** ohne den neuen Schlüssel
ist der Befundsatz byte-identisch
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)). Keine
bestehende Konfiguration ändert ihr Verhalten — das ist die Zusage, an der
jeder der drei Slices gemessen wird.

## 2. Trigger (Welle startet)

Auftraggeber-Freigabe 2026-08-22 (alle drei Change Requests bestätigt);
[slice-121](done/slice-121-zustandsfeld-hygiene.md) ist geschlossen, die
Hygiene-Hälfte liegt. Die drei Engstellen sind in den Closure-Notizen der
welle-80 und welle-81 sowie im Beobachtungs-Register belegt.

## 3. Closure-Trigger (Welle schließt)

- Alle vier Slices in `done/`; `make fullbuild` grün (Exit explizit).
- Je Erweiterung ein **Default-Beweis**: ohne den neuen Schlüssel ist der
  Befundsatz byte-identisch gegen das gepinnte Vorgänger-Image — grün wie rot.
- Je Erweiterung eine **konstruierte Gegenprobe**, die ohne sie stumm bliebe.
- **Release `v0.63.0`** auf GHCR mit Digest-Pin, Doku-Currency nachgezogen,
  Digest-Backfill committet.
- Die 3×-Form von BEO-008 ist danach **baubar** — ob sie gebaut wird, ist ein
  eigener Entscheid.
- Ergebnisnotiz `welle-82-results.md` mit Register-Lese-Schritt.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [slice-122](open/slice-122-versions-musterliste.md) | `versions`: mehrere Muster-Quellen-Paare statt eines — macht die Beobachtungs-3×-Form baubar | [`DC-FA-VER-001`](../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) |
| [slice-123](open/slice-123-structure-heading-muster.md) | `structure`: ein Schlüssel „jede Überschrift des Abschnitts matcht dieses Muster" statt ausgeschriebener Negation | [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) |
| [slice-124](open/slice-124-diagrams-ventile.md) | `diagrams`: Datei-Ventil, Zeilen-Marker und die fehlenden §2-Schema-Zeilen — Ventil-Parität zu den übrigen Modulen | [`DC-FA-DIAG-001`](../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in) |
| [slice-125](open/slice-125-release-v0630.md) | Release-Prep über alle drei Erweiterungen und Release `v0.63.0` (Tag, GHCR, Digest-Backfill) | [`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image) |

**Eine ADR für die Welle**, nicht drei: die drei Entscheidungen teilen eine
Begründung (eine Fläche additiv weiten statt eine Ersatz-Konstruktion zu
pflegen) und stehen als eigene Entscheidungen darin — dieselbe Form, in der
die Wellen-Invarianten-ADR mehrere Entscheidungen trägt.

## 5. Abhängigkeiten

- Blockiert: nichts Geplantes.
- Wird blockiert von: nichts. Reihenfolge innerhalb: 122, 123 und 124 sind
  voneinander unabhängig (drei Module, drei Anforderungen); 125 setzt alle drei
  voraus.

## 6. Out-of-Scope für diese Welle

- **Die 3×-Form von BEO-008 wird nicht gebaut.** Diese Welle macht sie
  *baubar*; ob das eigene Profil sie dann fährt, ist ein eigener Entscheid mit
  eigener Messung.
- **Kein Rückbau der Ersatz-Konstruktionen im eigenen Profil**, außer wo der
  neue Schlüssel sie unmittelbar ablöst (die Präfix-Negation in slice-123) —
  das `diagrams`-Scoping bleibt, bis eine Messung zeigt, dass das Ventil
  reicht.
- **Keine Default-Änderung.** Jede Erweiterung ist opt-in; ohne den Schlüssel
  bleibt der Befundsatz byte-identisch.
- **Keine neuen Anforderungs-Kennungen:** alle drei erweitern eine bestehende
  Anforderung (Einzelmodul-Frage ⇒ bestehende Anforderung ändern, das
  etablierte Schnitt-Kriterium).
