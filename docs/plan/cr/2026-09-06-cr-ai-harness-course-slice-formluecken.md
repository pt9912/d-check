# Change Request an `ai-harness-course` — Der Slice-Plan hat zwei Formlücken

**Absender:** d-check (Adopter) · **Datum:** 2026-09-06
**Richtung:** ausgehend ([`MR-035`](../../../harness/conventions.md#mr-035))
**Ziel:** `lab/templates/docs/plan/planning/slice.template.md`, `modul-09-implementierung.md`
**Baseline-Stand:** `v6.3.1`
**Stand:** **beantwortet, beide Bitten angenommen** — siehe
[Antwort](2026-09-06-antwort-ai-harness-course-slice-formluecken.md).

---

## Bitte 1 — Der Slice-Plan braucht einen Out-of-Scope-Abschnitt

`modul-06-roadmap.md` §Wellen-Closure-Prozedur schreibt über die Eröffnung:

> „Out-of-Scope gehört dazu — dieselbe Disziplin wie im Lastenheft (Modul 3)
> und **im Slice-Plan (Modul 9)**; was nicht ausdrücklich ausgeschlossen ist,
> dehnt die Welle, bis der Closure-Trigger unerreichbar wird."

Der Satz nimmt eine Disziplin im Slice-Plan als gegeben an. Gemessen am Bundle
`v6.3.1`:

| Quelle | Treffer für Out-of-Scope |
|---|---|
| `templates/docs/plan/planning/slice.template.md` | **0** |
| `regelwerk/modul-09-implementierung.md` | **0** |

Weder das Modul noch die Vorlage kennt den Abschnitt. Die Begründung, die
`modul-06` für die Welle gibt, gilt für den Slice unverändert: Was nicht
ausdrücklich ausgeschlossen ist, dehnt ihn — und anders als bei der Welle gibt
es beim Slice keinen zweiten Träger, der es auffängt.

**Vorschlag:** Ein Abschnitt zwischen *Vorgehen/Plan* und *Definition of Done*,
der nennt, was der Slice **nicht** tut, mit Begründung je Punkt.

## Bitte 2 — Die beiden `Vorgelagert`-Blöcke stehen im Abschnitt, der entfallen darf

`slice.template.md` §8 *Sub-Area-Modus-Begründung* sagt in einem Absatz beides:

> „Bei reinem Refactor ohne neue Sub-Area-Berührung **entfällt er ganz**. Die
> beiden *Vorgelagert*-Blöcke **entfallen nie**."

Die Blöcke stehen aber **innerhalb** von §8. `modul-05` begründet ihre
Unabhängigkeit ausdrücklich — sie „hängen weder am Modus noch am Slice-Typ und
stehen deshalb in **jedem** Slice-Plan". Genau diese Unabhängigkeit ist
strukturell nicht abgebildet: Ein Abschnitt, dessen Titel die Modus-Begründung
nennt und der ganz entfallen kann, beherbergt zwei Pflichtblöcke, die mit dem
Modus nichts zu tun haben. Heute überbrückt ein Satz, was die Gliederung
widerlegt.

**Vorschlag:** Die beiden Blöcke als **eigener Abschnitt** vor der
Modus-Begründung. Der erklärende Satz („entfallen nie") wird dadurch
entbehrlich, weil die Struktur die Aussage selbst trägt.

## Was wir dafür an Erfahrung anbieten

d-check führt beide Abschnitte als Haus-Form und hat sie an geschlossenen
Slices gelebt. Der Out-of-Scope-Abschnitt hat dabei nachweisbar getragen — er
ist die Stelle, an der ein wachsender Slice benennt, was er abgibt, statt es
stillschweigend mitzunehmen. Formulierungen stellen wir gern bereit.

## Abgrenzung

- **Keine Bitte um die Umbenennungen**, die d-checks Haus-Form sonst führt
  (*Vorgehen* statt *Plan (vor Code)*, *Abnahme-Punkte / Risiken* statt
  *Risiken und offene Punkte*, Closure-Notiz am Ende statt in der Mitte). Das
  ist Geschmack, nicht Lücke, und gehört nicht in einen CR.
- **Keine Änderung an den Pflichtfeldern** des Kopfes oder an der DoD.
- **Verworfene Alternative:** die Abschnitte dauerhaft als lokalen Fork zu
  führen. Das lässt jeden weiteren Adopter dieselbe Lücke einzeln entdecken —
  und die erste Bitte betrifft eine Aussage, die der Kanon über sich selbst
  macht, also nicht sinnvoll lokal auflösbar.
