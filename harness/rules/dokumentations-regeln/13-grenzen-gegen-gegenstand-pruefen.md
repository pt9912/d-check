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

**Eine Grenze, die an einer zweiten Stelle wiedergegeben wird**, nimmt ihre Bedingung mit — in ein
Kriterium, eine Historie-Zeile, die Spezifikation, eine ADR, einen Kommentar oder einen Testfall. Für
eine Begründung gilt dasselbe: Sie steht nur bei den Fällen, die sie trägt. Vor dem Handoff wird jede
Stelle, die dieselbe Aussage trägt, über ihren **Inhalt** gesucht, nicht über ihren Wortlaut; eine Stelle,
die die Grenze ohne ihre Bedingung trägt, sagt für den ausgenommenen Fall das Gegenteil zu. Im Slice-Plan
trägt das [Regel 18](18-plan-festlegung-zitiert-ihre-quelle.md) (zitieren statt nacherzählen), vor einem
Fix nach dem Review die Suche über die Aussage im Workflow-Skelett; dieser Absatz gilt beim ersten
Schreiben und überall sonst. Urteil, kein `grep`. *(Hard Rule aus dem Steering Loop,
[`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../../../docs/plan/planning/observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md),
seit welle-97; Auflösungs-Trigger: permanent.)*
