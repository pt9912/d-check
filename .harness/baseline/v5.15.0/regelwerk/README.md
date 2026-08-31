# Regelwerk — der Kurs als Betriebsregelwerk, nach Modulen

**Stand:** Kurs-Welle 106 · 2026-08-31.

**Geltungsbereich:** getestet und gelebt mit **einem schreibenden Menschen plus
Agenten** je Repository; die Mehr-Schreiber-Fassung ist entworfen, nicht
belegt.

Die 17 Module (0–16) **und die Grundlagen-Abschnitte** des Kurses als
**Betriebsregelwerk für Code-Agenten** — didaktik-freier Extrakt (Regeln, Konventionen, Abläufe in
Quellformulierung; weggelassen ist die Didaktik-Schicht, nicht verdichtet der
Inhalt — umgehängt werden beim Split nur relative Links und Selbstverweise:
„in diesem Kurs" wird „in diesem Regelwerk", „Kurs" als Name der *Quelle*
bleibt). Pro Abschnitt eine Datei, damit ein Agent einen einzelnen Abschnitt
laden kann, ohne das ganze Regelwerk im Kontext zu halten.

> **Was dieses Verzeichnis ist.** Das **kanonische Regelwerk-Artefakt** (login-frei
> ausgeliefert als `lab-regelwerk.zip`, self-navigierbar). Es trägt keine eigene
> Normativität: maßgeblich für den *Inhalt* bleibt der Kurs unter
> [`/kurs/de/`](https://github.com/pt9912/ai-harness-course/blob/v5.15.0/kurs/de/README.md) — die Module der Phasen 01–05 plus die
> normativen Grundlagen-Abschnitte.
>
> **Was dieses Verzeichnis NICHT ist.** Eine eigene Quelle der Wahrheit. Wer hier
> eine Regel ändert, ohne die Kurs-Quelle zu ändern, erzeugt genau die Drift, die
> das Regelwerk selbst verbietet.

**Vendored gelesen?** Dann liegt dieses Verzeichnis als `regelwerk/` unter
`.harness/baseline/<tag>/` eines Adopter-Repos — entpackt aus dem
self-contained Baseline-Bundle
(<https://github.com/pt9912/ai-harness-course/releases/latest/download/lab-regelwerk.zip>),
das `regelwerk/` und [`templates/`](../templates/README.md) **parallel** trägt.
Deshalb lösen die `../templates/…`-Verweise der Abschnitte („Ziel-Form: X")
netzlos auf. Der Einstieg ist dort nicht diese Datei, sondern **`AGENTS.md` des
Adopter-Repos** ([Vorlage](../templates/AGENTS.template.md)): Es briefingt die
verkörperte Form und verweist hierher als **nachschlagbare Vertiefung** — pro
Entscheidung **nur den benötigten Abschnitt** laden (die Liste unten ist der
Index), nicht das ganze Regelwerk im Kontext halten. Vorgehen beim Bootstrap:
[Modul 2 — Harness-Bootstrap](modul-02-harness-bootstrap.md) (liegt hier im
Bundle, netzlos).

**Links.** Im Repo relativ (lokal navigierbar, vom Doku-Gate validiert). Relativ
bleiben auch im ausgelieferten `lab-regelwerk.zip` die Modul-Querverweise
*innerhalb* dieses Verzeichnisses (`--keep-within-src`) **und** die
`../templates/…`-Verweise (`--keep-within=lab/templates`) — beide Verzeichnisse
reisen im Bundle parallel mit, es ist self-navigierbar. Nur Verweise auf den
Kurs, der *nicht* mitreist, werden beim Release auf absolute, auf den Tag
gepinnte GitHub-URLs umgeschrieben (`tools/rewrite-doc-links.py`).

## Abschnitte

### Grundlagen

- [Kernbegriffe und Trennschärfen](grundlagen-begriffe.md)
- [Source Precedence und Spec-Stratifizierung](grundlagen-source-precedence.md)
- [Referenz-Richtung (SDP)](grundlagen-referenz-richtung.md)
- [Die Harness-Dateien und ihre Form](grundlagen-harness-dateien.md)
- [Harness-Bootstrap](grundlagen-bootstrap.md)
- [Traceability-Constraint](grundlagen-traceability.md)
- [Klassifikation und Steering Loop](grundlagen-klassifikation.md)
- [Durchsetzungsschicht](grundlagen-durchsetzungsschicht.md)

### Einführung

- [Modul 0 — Einführung](modul-00-einfuehrung.md)

### Phase 01 — Spec und Architektur

- [Modul 1 — Der Entwicklungszyklus](modul-01-entwicklungszyklus.md)
- [Modul 2 — Harness-Bootstrap](modul-02-harness-bootstrap.md)
- [Modul 3 — Die Spec: Lastenheft, Spezifikation, Architektur](modul-03-spec.md)
- [Modul 4 — ADRs](modul-04-adrs.md)

### Phase 02 — Planung

- [Modul 5 — Planning Harness](modul-05-planning-harness.md)
- [Modul 6 — Roadmap Engineering](modul-06-roadmap.md)
- [Modul 7 — Carveout Management](modul-07-carveouts.md)

### Phase 03 — Agenten

- [Modul 8 — Agentenrollen](modul-08-agentenrollen.md)
- [Modul 9 — Implementierung durch KI-Agenten](modul-09-implementierung.md)

### Phase 04 — Qualität

- [Modul 10 — Review Harness](modul-10-review-harness.md)
- [Modul 11 — Verification Harness](modul-11-verification.md)
- [Modul 12 — Replay und Evaluierung](modul-12-replay-evaluierung.md)
- [Modul 13 — Quality Gates](modul-13-quality-gates.md)

### Phase 05 — Betrieb

- [Modul 14 — Docker Harness](modul-14-docker-harness.md)
- [Modul 15 — Observability](modul-15-observability.md)
- [Modul 16 — Produktiver Betrieb](modul-16-produktiver-betrieb.md)

## Lizenz

Wie der übrige Kurs: Texte unter CC BY 4.0, Code-Artefakte unter MIT. Details in
[`LICENSE.md`](https://github.com/pt9912/ai-harness-course/blob/v5.15.0/LICENSE.md).
