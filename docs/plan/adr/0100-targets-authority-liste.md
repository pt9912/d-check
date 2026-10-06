# ADR-0100: `targets.authority` nimmt eine Liste an; Doppelnennung ist kein Befund

**Status:** Accepted

**Datum:** 2026-10-06

**Autor:** pt9912

**Bezug:** [`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(erweitert, Lastenheft 0.96.0, nach Review präzisiert in 0.96.1); Anlass ist der eingehende
[CR von `ai-harness-init`](../cr/2026-10-06-cr-eingehend-ai-harness-init-targets-authority-liste.md);
Vorgängerin am selben Modul [ADR-0099](0099-targets-makefiles-glob.md)
(`authority` bleibt wörtlich — hier bestätigt); Schnitt-Kriterium aus
[ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md);
slice-251 <!-- d-check:status-provenance -->.

**Schärft:**
[`spec/spezifikation.md` §DC-FA-TGT-001.a](../../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
(Schritt 5), die Schema-Zeile `targets.authority` unter
[`SPEC-005`](../../../spec/spezifikation.md#spec-005--d-checkyml) und die
Grund-Code-Zeile
[`SPEC-061`](../../../spec/spezifikation.md#4-grund--und-fehler-codes).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`gate-undocumented` misst gegen genau eine Autoritäts-Datei. Ein Adopter, dessen
Gate-Index aus zwei Teilen besteht — einem eigenen und einem, den ein
Bootstrap-Werkzeug bei jedem Lauf neu schreibt —, kann die Werkzeug-Targets
heute nur über `exempt-targets` stilllegen, echte Gates eingeschlossen. Ein
stilles Gate in einer Ausnahmeliste ist genau das, was das Modul verhindern
soll.

Der Absender hat parallel einen CR an die Kurs-Baseline gestellt, der einen
werkzeug-eigenen Teil des Gate-Index zulassen soll, mit der Bedingung, dass
kein Target in beiden Teilen steht. Die gepinnte Baseline (`v6.13.0`) kennt
diese Regel nicht; dort gibt es genau eine Autoritäts-Doku.

## Entscheidung

1. **`targets.authority` ist ein Pfad oder eine Liste wörtlicher Pfade.**
   `gate-undocumented` misst gegen die Vereinigung der dort dokumentierten
   Targets. Ein String ist die einelementige Liste; mit einer Datei bleibt der
   Befundsatz samt Meldungstext byte-identisch.
2. **Eine Doppelnennung ist kein Befund.** Die Disjunktheit mehrerer Teile ist
   eine Regel, die nur in einem offenen CR an die Baseline steht. Sie jetzt im
   Produkt zu verankern, hieße, einer Baseline-Entscheidung vorzugreifen und
   Konsumenten rot zu machen, deren Baseline sie nicht verlangt. Nimmt der Kurs
   die Regel an, ist ein opt-in-Schlüssel mit eigenem Grund-Code der
   Folgeschritt (Re-Evaluierungs-Trigger unten).
3. **Kein Glob für `authority`, und keine Ablehnung von Glob-Zeichen** — die
   Entscheidung der Vorgängerin bleibt: Einträge sind wörtliche Pfade, auch
   mit Glob-Zeichen. Ein Muster, das als Datei nicht existiert, scheitert
   laut zur Laufzeit (Exit 2). Eine Ablehnung beim Laden hätte einen
   bisher gültigen wörtlichen Pfad mit eckiger Klammer gebrochen.
4. **Leer heißt entfällt — aber nur auf oberster Ebene.** Die String-Form
   dekodiert wie zuvor (ein leerer oder Null-Pfad lässt Richtung 2 entfallen,
   ein Alias wird aufgelöst); eine leere Liste entfällt wie `doc-tables: []`.
   Ein leeres, Null- oder Nicht-Skalar-**Element** einer Liste ist dagegen ein
   Konfigurationsfehler — geprüft je YAML-Knoten, weil der Decoder
   Null-Elemente beim Dekodieren in eine String-Liste still verwirft.
5. **Der `--doctor`-Klartext bleibt im Singular** („in der Autoritäts-Doku").
   Er beschreibt den Grund-Code, nicht den Befund; die Befund-Meldung nennt
   bei mehreren Dateien alle. Eine Änderung verschöbe den Klartext für jeden
   Konsumenten, auch für die unveränderte String-Form.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Disjunktheit sofort als Befund** | fängt die vom Kurs-CR verbotene Doppelnennung | verankert eine Regel, die die Baseline nicht hat; macht Konsumenten rot ohne Grundlage |
| **Disjunktheit sofort als opt-in** | kein Default-Bruch | Umfang ohne Anlass vor der Baseline-Entscheidung; Grund-Code auf Vorrat |
| **Werkzeug-Targets in `exempt-targets`** | keine Änderung | stilllegt echte Gates; vom CR als Anlass benannt |
| **Liste, Vereinigung, Doppelnennung still** (gewählt) | rückwärtskompatibel; schließt die Lücke ohne Baseline-Vorgriff | die Disjunktheit bleibt bis zur Baseline-Entscheidung ungeprüft (benannt) |

## Konsequenzen

- Ein Adopter mit zweiteiligem Gate-Index trägt beide Dateien ein, statt
  Werkzeug-Gates auszunehmen.
- Die Meldung von `gate-undocumented` nennt bei mehreren Dateien alle.
- Schema, Grund-Code-Zeile und `--print-config`-Gerüst nennen die Listenform.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test `TestCheckTargetsAuthorityListe` | Vereinigung, Doppelnennung still, Befund im Fragment mit allen Dateien | `make test` |
| Go-Test `TestCheckTargetsAuthorityEinzelnWortlaut` | eine Datei: Meldungstext unverändert | `make test` |
| Go-Test `TestCheckTargetsAuthorityListeFehlend` | fehlender Listeneintrag ⇒ Fehler (Exit 2) | `make test` |
| Go-Test `TestDecode_TargetsAuthority` | String/Null/Alias/Liste/leer wie zuvor, Pfad mit Glob-Zeichen wörtlich; leeres/Null-/Nicht-Skalar-Element, Abbildung, Wurzel-Flucht ⇒ Konfigurationsfehler | `make test` |

## Re-Evaluierungs-Trigger

Die Kurs-Baseline nimmt die Regel an, dass kein Target in mehreren Teilen des
Gate-Index steht — dann Folgeschritt: opt-in-Prüfung der Disjunktheit mit
eigenem Grund-Code. Ebenso: ein Adopter, dessen Autoritäts-Menge wächst, macht
einen Glob nötig.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-06 | Angelegt als `Proposed`; `Accepted` erst mit der Closure des Vorgangs, nach Review und Verifikation |
| 2026-10-06 | Nach R1 (HIGH, MEDIUM): Entscheidung 3 auf „keine Ablehnung von Glob-Zeichen" umgestellt, Entscheidung 4 auf Knoten-Prüfung der Listenelemente (Null-Elemente verschwanden still), Entscheidung 5 (`--doctor`-Klartext) ergänzt. Noch `Proposed`, Körper daher geändert statt angehängt |
| 2026-10-06 | Nach Verifikation: `Nicht-Pfad-` zu `Nicht-Skalar-Element` präzisiert (V5). `Accepted` mit der Closure des Vorgangs |
