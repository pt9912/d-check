# Slice slice-143: Der Platzhalter-Erkenner des Produkts sieht nur einen Abschnitt

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [slice-139](../done/slice-139-closure-ausgang-waechter.md); `tools/harness/closure-outcomes.sh`; [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md); [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in), [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in).

**Berührte Spec-Stellen:** `spec/lastenheft.md` — eine Anforderung wächst (Abschnitts-Skopus bzw. Muster-Verbot); Bump und Historie nach [`MR-032`](../../../../harness/conventions.md#mr-032).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

Das Produkt erkennt Vorlagen-Platzhalter fence- und inline-code-bewusst — aber
**nur im Abschnitt der Closure-Notiz**. Für den Rest des Slice steht seit
[slice-139](../done/slice-139-closure-ausgang-waechter.md) ein Zeichenketten-
Wächter als Bash-Skript daneben, mit zwei benannten Schwächen: seine
Platzhalter-Liste ist eine **Liste**, und er behandelt Fenced Code nicht.

Die saubere Form ist der **Abschnitts-Skopus im Produkt**: dieselbe Erkennung,
weiter gefasst — oder ein `forbid-match` in `structure` als Gegenstück zum
vorhandenen `headings-match`. Welche der beiden, entscheidet der Slice.

## 2. Vorgehen

1. **Die zwei Wege gegeneinander stellen**, mit ihren Folgen: den Skopus des
   `planning`-Platzhalter-Erkenners weiten, oder `structure` um ein
   Muster-Verbot ergänzen. Beides ist ein Produkt-Delta mit ADR.
2. **Am Bestand messen**, bevor scharfgeschaltet wird — beide Wege treffen mehr
   als den einen Abschnitt.
3. Das Bash-Skript **ersatzlos entfernen**, sobald der Skopus es abdeckt: sein
   Auflösungs-Trigger steht in seinem Kopf.
4. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Doppelführung.** Am Ende trägt **eine** Mechanik die Frage; die
  andere fällt weg.
- **Keine Ausweitung auf andere Dokumentklassen** in diesem Zug.

## 4. Definition of Done

- [x] Die Wegwahl ist begründet, mit der Folge des verworfenen Wegs — in
      [ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md)
      samt fünf verglichenen Alternativen. **Beide geplanten Wege waren falsch
      gerahmt** (§9).
- [x] Der Bestand ist gemessen: 461 Dateien, **0 Befunde** — vor und nach der
      Selektor-Weitung. Keine Fundstelle zu räumen, also auch keine
      auszuweisen.
- [x] Das Bash-Skript ist entfernt samt aller sieben Spiegel; sein
      Auflösungs-Trigger ist eingelöst — **auf einem anderen Weg, als sein
      eigener Text vorsah** (§9).
- [x] ADR geschrieben ([ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md),
      Index-Zeile); `make fullbuild` Exit 0 (fünf Closure-Glieder statt sechs).
      **Lastenheft-Bump entfällt und das ist der Befund, nicht das Versäumnis:**
      die benutzte Fähigkeit ist in
      [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
      bereits zugesagt, es gibt kein Produkt-Delta.
- [x] `make gates` Exit 0 (zehn Glieder); unabhängiger Review
      ([Report](../../../reviews/2026-08-26-slice-143-structure-abschnitts-skopus-review.md)),
      blockierend mit **zwei HIGH**, alle fünf Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Ein weiter gefasster Skopus meldet mehr.** Was heute in Slice-Abschnitten
  legitim in Winkelklammern steht,   wird dann zum Befund — die Messung entscheidet,
  ob der Weg tragbar ist. — **Ausgang:** *entfallen — die Messung sagt nein.*
  Weder mit dem engen Selektor (`^# Slice slice-`) noch mit dem weiten
  (`^# ` plus `sections: each`) meldet der Bestand etwas: 0 Befunde über 460
  bzw. 461 Dateien. Der Grund ist die Bereinigung, nicht Glück — was in
  Slice-Abschnitten legitim in Winkelklammern steht, steht in Inline-Code.
- **Zwei Mechaniken gleichzeitig sind schlechter als eine.** Bleibt das Skript
  neben dem Produkt stehen, ist die Doppelung dauerhaft statt übergangsweise. —
  **Ausgang:** *entfallen — das Skript ist weg.* Was bleibt, ist keine
  Doppelung, sondern ein **gemessenes** Komplement: `closure-note-placeholder`
  fängt jede whitespace-freie Winkelklammer-Form, aber nur im
  Closure-Abschnitt; die neue Regel vier benannte Formen über jeden
  H1-Abschnitt. Die Gegenprobe steht in §9.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung zeigt, dass der weitere Skopus den Bestand breit trifft.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25): [`BEO-011`](../observations.md) für die Bestandsmessung; [`BEO-015`](../observations.md), weil dieser Slice die dortige mechanische Form einlöst.

Slice-ID: slice-143. Betroffene IDs: [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in), [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in). Module: `planning`, `structure`.
Gates: `make doc-check`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Erweiterung einer bestehenden Modul-Fähigkeit.

## 9. Closure-Notiz (nach `done/`)

Geliefert: Die urteilsfreie Hälfte der Drei-Ausgänge-Regel wird von einer Regel
im Closure-Profil getragen; `tools/harness/closure-outcomes.sh` ist entfernt.
`make fullbuild` hat fünf Closure-Glieder statt sechs.

**Beide geplanten Wege waren falsch gerahmt, und beides zeigte erst das
Nachsehen.** Weg A — den Abschnitts-Skopus des `planning`-Erkenners weiten —
**kann** das Skript nicht ablösen: sein Muster verlangt whitespace-freie
Winkelklammern und deckt genau **eine** der vier Formen. Weg B braucht **nichts
Neues**: `structure` trägt `forbid-pattern` längst. Der Slice-Plan hielt es für
fehlend, und ich habe das aus dem Gedächtnis bestätigt statt nachzusehen.

**Der Auflösungs-Trigger des Skripts war in sich widersprüchlich.** Er versprach
*„sobald der Abschnitts-Skopus des Moduls den ganzen Slice umfasst, fällt dieses
Skript ersatzlos"* — zwei Zeilen darüber sagt derselbe Kopf, der Wächter decke
*„zwei Formen, die jenes Muster nicht kennt"*. Ein Trigger, der die eigene
Nachbarzeile nicht liest, ist eine Zusage ohne Deckung.

**Der teuerste Befund widerlegt meine eigene Commit-Botschaft.** Sie behauptete,
die Zusage werde *„nicht schwächer, sondern stärker"*. Auf einer gemessenen
Achse ist sie schwächer: das Skript paarte Backticks **je Zeile**, das Produkt
**absatzweise**. Ein offener Ausgang zwischen zwei Backticks desselben Absatzes
ist unsichtbar — Skript Exit 1, Regel Exit 0. Das ist
[`BEO-009`](../observations.md) Richtung (b): die Proben stimmten, der Schluss
reichte weiter.

**Die Richtigstellung ändert die Entscheidung nicht.** Die absatzweise Paarung
ist die **korrektere** Markdown-Lesart; der Vorteil des Skripts stammte aus
einer falscheren Lexik, nicht aus einem Entwurf. Gegengeprüft, ob ein anderes
Gate die Lücke schließt: `spans` kennt nur `span-unclosed`, also **ungerade**
Parität — der Fall hat gerade Parität und läuft durch. Keine Deckung von dort,
und darum steht die Differenz in
[ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md) statt in
einer Fußnote. Dass es überhaupt eine ADR gibt, ist der Befund des Reviews:
[`AGENTS.md` §3.6](../../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)
gilt einer Prüfregel, deren Reichweite sinkt, und das war unbelegt.

**Zwei Zusagen waren zu breit formuliert.** *„Über die ganze Slice-Datei"* stand
an zwei Stellen, während vier Bereiche außen vor bleiben — und der Vorbehalt
*„Inline-Code wird übersprungen"*, den die abgelöste `AGENTS.md`-Zeile noch
trug, war beim Umschreiben ersatzlos entfallen. Beide Stellen zählen die vier
jetzt einzeln auf: drei gewollt, der vierte als Preis.

**Eine stille Kappung ist zu.** Der erste Selektor traf die Titelzeile; eine
zweite H1 hätte die Spanne still beendet, und `sections: one` fängt das nicht,
weil es nur **Selektor-Treffer** zählt. Jetzt `^# ` mit `sections: each` —
gegengeprobt an einem Platzhalter hinter `# Anhang`, der jetzt gefunden wird.

**Was die Regel besser kann als das Skript**, über die Lexik hinaus: sie kann
nicht still leerlaufen. Greift der Selektor in einer Datei nicht, meldet sie
`section-missing`. Das Skript brauchte dafür eine eigene Fail-Closed-Klausel.

**Gemessen, nicht behauptet:** 461 Dateien, 0 Befunde. Vier Formen einzeln
gebrochen, jede rot mit gelesener Fundstelle, jede nach Rückbau grün.
Negativkontrolle: dieselbe Form in Inline-Code bleibt grün. Beides am
`make`-Target gefahren, nicht nur am Container-Aufruf.

**Register:** [`BEO-016`](../observations.md) angelegt — Prosa verschwindet in
einer absatzweiten Inline-Code-Spanne, Zähler 1. Die Klasse ist breiter als
dieser Wächter: sie trifft **jedes** prosa-lesende Modul.
