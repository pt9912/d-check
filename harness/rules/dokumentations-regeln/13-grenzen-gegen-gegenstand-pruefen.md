# Wer eine Grenze aufschreibt, prüft sie gegen den Gegenstand, nicht gegen seine Beschreibung

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

Ein Abschnitt, der aufzählt, **was ein grüner Lauf nicht abdeckt**, liest
sich durch seine Form als Menge und ist immer eine Auswahl. Vor dem
Handoff deshalb zweierlei: **den Vertrags-Teil desselben Artefakts
durchgehen und jede Zusage einmal umdrehen** — was folgt daraus für das
Grün? —, und **wo der Gegenstand Code oder Konfiguration ist, gegen diese
prüfen statt gegen die Prosa darüber**. In den belegten Fällen stand die
fehlende Grenze fast immer bereits im Vertrags-Teil oder in der
Konfiguration; in einem stand sie **nur im Code**, und der Vertrags-Text
daneben sagte das Gegenteil des Verhaltens — **dafür ist die zweite Hälfte
der Regel da**, denn eine nur aus der Prosa abgeleitete Grenze beschrieb
dort einen Mechanismus, den es nicht gibt.

**Nächste Verwandte:** [`AGENTS.md` §3.8](../../../AGENTS.md#38-ein-modul-verspricht-nur-über-das-was-es-scannt)
verlangt dieselbe Umkehrung für ein **Modul** und seine Scan-Menge; dieser
Absatz verlangt sie für **jede** aufgeschriebene Grenze. **Drei Grenzen:**
Die Regel gilt dem **Autor vor der Übergabe** und ersetzt den fremden
Leser nicht — in den belegten Fällen fand die Lücke **jemand anderes als
der Autor**; [`AGENTS.md` §6](../../../AGENTS.md#6-minimal-agent-workflow)
richtet davon den Review ein, nicht jeden fremden Leser. Ihre erste Hälfte
setzt einen **korrekten** Vertrags-Text voraus — wo er lügt, fängt nur die
zweite. Und belegt ist sie an **Sensor-Beschreibungen**, nicht an
Grenzen-Listen überhaupt. Urteil, kein `grep`; der Reviewer-Skill trägt den
Anker dazu. *(Hard Rule aus dem Steering Loop,
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../../../docs/plan/planning/observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md),
seit slice-213; Auflösungs-Trigger: permanent.)*
