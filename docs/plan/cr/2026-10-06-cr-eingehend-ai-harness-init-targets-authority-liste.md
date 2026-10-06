# Eingehender Change Request — `targets.authority` nimmt eine Liste an

**Absender:** Adopter `ai-harness-init` · **Eingegangen:** 2026-10-06
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(Modul `targets`, Schlüssel `targets.authority`)
**Stand:** **eingegangen**; Entscheidung und Umsetzung trägt `slice-251`.

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
