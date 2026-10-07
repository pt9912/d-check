# Eingehender Change Request — `targets.authority` nimmt eine Liste an

**Absender:** Adopter `ai-harness-init` · **Eingegangen:** 2026-10-06
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(Modul `targets`, Schlüssel `targets.authority`)
**Stand:** **angenommen und umgesetzt am 2026-10-06** — alle sieben
Akzeptanzkriterien erfüllt, beide offenen Fragen beantwortet (unten); Träger
`slice-251`, Lastenheft 0.96.2.

**Ablage-Hinweis.** Ein **eingehender** CR ist die dritte Klasse neben
[`MR-035`](../../../harness/conventions.md#mr-035) (ausgehend) und
[`MR-036`](../../../harness/conventions.md#mr-036) (Antworten darauf). Die
Datei liegt hier aus demselben Grund wie ihre Vorgänger in diesem
Verzeichnis: Was gebeten wurde und wie entschieden wird, soll im Repo
stehen und nicht nur im Vorgang. Im Wortlaut sind Kennungen verlinkt
(Linkpflicht des eigenen Gates), sonst unverändert.

**Abhängigkeit des Absenders:** Der CR setzt eine Regelwerks-Änderung im
Kurs voraus, die der Absender parallel an `ai-harness-course` richtet
(werkzeug-eigener Teil des Gate-Index neben `harness/README.md`, gemessen am
Stand `v6.13.0`). Jener CR richtet sich nicht an dieses Repo und liegt
deshalb nicht hier; seine für diese Entscheidung relevante Bedingung — kein
Target steht in beiden Teilen — ist unten unter *Entscheidung* aufgenommen.

---

## Wortlaut

> **CR an d-check: targets.authority nimmt eine Liste an**
>
> Absender: ai-harness-init (Adopter von d-check). Gegenstand:
> [`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in), Modul targets, Schlüssel targets.authority. Gemessen am Stand: v0.81.0.
> Abhängigkeit: Er setzt eine Regelwerks-Änderung im Kurs voraus (CR oben).
> Die d-check-Seite ist aber eigenständig sinnvoll und rückwärtskompatibel.
>
> ### Anlass
>
> Ein gebootstrapptes Ziel trägt zwei Teile seines Gate-Index.
> harness/README.md (Adopter) dokumentiert die eigenen Targets. Eine
> werkzeug-eigene, bei jedem Bootstrap neu geschriebene Datei dokumentiert
> die Targets der Werkzeug-Fragmente. Der Deklarations-Sensor soll
> Vollständigkeit gegen beide zusammen messen.
>
> ### Ist-Verhalten (v0.81.0)
>
> - Schema: targets.authority ist ein einzelner String (Authority string,
>   configyaml.go).
> - gate-undocumented: misst nur gegen diese eine Datei
>   (undocumentedFindings(fsys, cfg.Authority, …), rules/targets.go).
> - doc-tables: ist eine Liste, wird aber nur für gate-phantom gelesen.
>
> Folge: Ein Target, das nur in der zweiten Datei steht, meldet
> gate-undocumented. Der einzige Ausweg ist exempt-targets mit allen
> Werkzeug-Targets, auch echten Gates.
>
> ### Bitte
>
> targets.authority nimmt neben einem String auch eine Liste von Pfaden an.
> gate-undocumented misst gegen die Vereinigung der dort dokumentierten
> Targets.
>
> ### Akzeptanzkriterien
>
> 1. String bleibt gültig: authority: harness/README.md verhält sich
>    byte-identisch wie heute.
> 2. Liste: Bei authority: [harness/README.md, harness/targets.md] gilt ein
>    Target als dokumentiert, wenn es in mindestens einer Datei steht.
> 3. Fehlende Datei: Ein Listeneintrag, der nicht existiert, bricht
>    fail-closed mit Exit 2 ab, wie heute bei einer fehlenden
>    Autoritäts-Datei.
> 4. Doppelt dokumentiert: Ein Target in zwei Autoritäts-Dateien ergibt
>    keinen Befund. Ob es ein Hinweis wert ist, entscheidet d-check, siehe
>    Frage 2.
> 5. Pfad-Regel: Jeder Eintrag ist relativ, kein /, kein .., wie bei
>    makefiles.
> 6. Rotes Gegenbeispiel: Ein Target in keiner der Dateien meldet
>    gate-undocumented mit Fundstelle im Makefile-Fragment.
> 7. Glob: Ob authority auch Glob-Muster annimmt (analog makefiles),
>    entscheidet d-check. Für den Anlass genügen wörtliche Pfade.
>
> ### Abgrenzung
>
> - Nicht Gegenstand: doc-tables, exempt-targets und die Regel-Extraktion
>   aus Makefiles.
>
> ### Offene Fragen an d-check
>
> 1. Soll --print-config ein Listen-Beispiel tragen?
> 2. Soll ein in zwei Autoritäts-Dateien doppelt dokumentiertes Target ein
>    Befund sein (Drift-Signal), oder nicht?

## Entscheidung

**Angenommen.** `targets.authority` ist ein Pfad oder eine Liste wörtlicher
Pfade; `gate-undocumented` misst gegen die Vereinigung der dort
dokumentierten Targets. Die Abgrenzung des CR ist übernommen (`doc-tables`,
`exempt-targets` und die Regel-Extraktion unverändert). Begründung in
begleitender [ADR-0100](../adr/0100-targets-authority-liste.md).

**Zu den Akzeptanzkriterien:**

- **1 (String):** byte-identisch, auch mit `--json` und `--doctor`; das
  gilt ebenso für `null`/`~`, einen YAML-Alias und einen wörtlichen Pfad mit
  Glob-Zeichen. Eine einelementige Liste verhält sich wie der String.
- **3 (fehlende Datei):** Exit 2 mit ihrem Namen, an jeder Position der Liste.
- **4 (doppelt dokumentiert):** kein Befund (siehe Frage 2).
- **5 (Pfad-Regel):** gilt je Eintrag.
- **7 (Glob):** **nein** — Einträge sind wörtliche Pfade, auch mit
  Glob-Zeichen. Die Autoritäts-Menge ist klein und fest; ein Muster, das als
  Datei nicht existiert, scheitert laut mit Exit 2.

**Antworten auf die offenen Fragen:**

1. **`--print-config`: ja.** Das Gerüst nennt die Listenform mit dem
   Beispiel `[harness/README.md, harness/targets.md]`.
2. **Doppelnennung: kein Befund — vorerst.** Die Regel „kein Target steht in
   beiden Teilen" steht im parallelen CR des Absenders an den Kurs, nicht in
   der gepinnten Baseline. Sie jetzt zu prüfen, hieße, der Baseline-
   Entscheidung vorzugreifen. **Benannter Folgeschritt:** Nimmt die
   Kurs-Baseline die Regel an, bekommt `targets` eine opt-in-Prüfung der
   Disjunktheit mit eigenem Grund-Code (Re-Evaluierungs-Trigger der ADR).

**Vom CR nicht benannte Festlegungen, die der Absender kennen sollte:**

- Eine leere Liste lässt Richtung 2 entfallen, wie `doc-tables: []`
  Richtung 1. Ein leerer, Null- oder Nicht-Skalar-**Eintrag** einer Liste ist
  dagegen ein Konfigurationsfehler (Exit 2, mit Zeile) — er fällt nicht still
  weg.
- Dieselbe Datei in zwei Schreibweisen (mit und ohne führendes `./`) zählt
  einmal; die Meldung bleibt dann im Singular.
- Mit mehreren Dateien lautet die Meldung „… ohne Deklaration in einer der
  Autoritäts-Dokus `<a>, <b>`"; der `--doctor`-Klartext des Grund-Codes
  bleibt unverändert.

**Verfügbar ab:** d-check `v0.82.0`.

**Nachtrag 2026-10-07:** Der benannte Folgeschritt ist eingelöst. Die
Baseline führt die Disjunktheit seit `v6.16.0`; `targets` prüft sie opt-in
über `targets.authority-disjoint: true` mit dem Grund-Code
`gate-declared-twice` ([ADR-0101](../adr/0101-targets-authority-disjunkt.md),
d-check `v0.83.0`). Ohne den Schalter bleibt die Doppelnennung still wie oben
beschrieben.
