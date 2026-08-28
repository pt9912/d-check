# ADR-0068: Lokale Workflow-Referenzen tragen keinen SHA — und statt des Pins die Existenz-Prüfung

**Status:** Accepted

**Datum:** 2026-08-28

**Autor:** pt9912

**Bezug:** [`AGENTS.md`](../../../AGENTS.md) §3.9 (die Regel, die hier präzisiert
wird), §3.6 (die Verfahrenspflicht, aus der diese ADR entsteht),
[ADR-0065](0065-spiegel-gleichheit-ist-der-config-digest.md) (der Anlass: der
Hub-Sync als aufrufbarer Workflow)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`AGENTS.md`](../../../AGENTS.md) §3.9 verlangt: *„Jeder `uses:`-Eintrag in
`.github/workflows/` nennt einen **vollen Commit-SHA** mit Tag-Kommentar
dahinter."* Die Begründung ist Supply-Chain-Härtung — ein Tag lässt sich
umhängen, ein SHA nicht.

Beim Bau des auslösbaren Hub-Sync
([ADR-0065](0065-spiegel-gleichheit-ist-der-config-digest.md)) entstand die
erste **lokale** Workflow-Referenz dieses Repos:

```yaml
uses: ./.github/workflows/hub-description.yml
```

`make workflow-pins` schlug an. **Für diese Klasse ist die Regel nicht zu
streng, sondern unerfüllbar:** eine `./`-Referenz *kann* keinen SHA tragen. Die
Alternative wäre, auf `workflow_call` zu verzichten und die Sync-Mechanik in
zwei Workflows zu duplizieren.

**Die Frage, ob das eine Lockerung ist, ist strittig — und deshalb steht sie
hier.** §3.6 verlangt für eine gelockerte Prüfregel eine ADR. Ob das
*Präzisieren* einer unerfüllbaren Regel dasselbe ist wie ihr *Senken*, ist ein
Urteil; eine ADR ist der Ort, an dem ein solches Urteil überprüfbar wird.

## Entscheidung

Wir wählen **die Ausnahme für den `./`-Präfix, mit einer anderen Prüfung an der
Stelle des Pins**.

1. **Eine `./`-Referenz braucht keinen Pin, weil sie stärker gebunden ist als
   einer.** Sie löst auf **denselben Commit** auf wie der aufrufende Workflow —
   sie kann per Konstruktion nicht driften. Der Zweck von §3.9 (keine
   beweglichen Referenzen) ist ohne Pin erfüllt. Das ist der Grund, warum wir
   dies nicht als Senkung führen.
2. **Statt des Pin-Checks die Frage, die hier Sinn ergibt: existiert das
   Ziel?** Ein vertippter lokaler Verweis trägt keinen Pin und fiele sonst erst
   zur Laufzeit auf. Der Wächter meldet ihn als `uses-local-missing`. **Die
   Ausnahme lässt damit keinen Eintrag ungeprüft** — sie tauscht eine Prüfung
   gegen eine andere, statt eine wegzunehmen.
3. **Die Ausnahme gilt ausschließlich dem `./`-Präfix.** Eine Referenz in ein
   fremdes Repository (`owner/repo/.github/workflows/x.yml@ref`) ist fremder
   Code und fällt unter §3.9 wie jede Action.
4. **Sie wird ausgewiesen, nicht verschwiegen.** Die Erfolgsmeldung nennt die
   Zahl: *„8 `uses:`-Einträge geprüft, davon 1 lokale Referenz(en) ohne
   Pin-Pflicht"*. Eine stille Ausnahme wäre genau der Grün-Pfad, den dieses
   Repo sonst jagt.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, §3.9 wörtlich lassen | die Hard Rule bleibt unangetastet | dann ist `workflow_call` in diesem Repo unmöglich, und die Sync-Mechanik lebt in **zwei** Dateien — Drift zwischen zwei Kopien ist die Klasse, die dieses Repo als `BEO-002`/`BEO-010` führt |
| B — den Wächter eine Ausnahme-**Liste** je Datei führen lassen | maximal explizit | eine Liste driftet und ist datei-, nicht klassen-gebunden; sie überlebt ihren Grund (`BEO-013`) |
| C — jeden `uses:`-Eintrag **ohne `@`** ausnehmen | eine Zeile weniger | zu weit: ein **fremder** Verweis ohne Version rutschte durch — genau der bewegliche Fall, den §3.9 verbietet |
| D — Ausnahme ohne Ersatz-Prüfung (nur `continue`) | am einfachsten | ein vertippter lokaler Verweis liefe still durch; die Ausnahme wäre ein Loch statt eines Tauschs |
| **E — `./`-Präfix ausnehmen, dafür Existenz prüfen, Zahl ausweisen (gewählt)** | kein Eintrag bleibt ungeprüft; die Ausnahme ist eng, sichtbar und begründet | §3.9 liest sich nicht mehr als ein Satz; wer nur die erste Zeile liest, hält jede Referenz für pin-pflichtig |

## Konsequenzen

- **Positiv:** `workflow_call` ist möglich, und die Mechanik der Hub-Darstellung
  lebt an **einem** Ort statt in zwei Kopien.
- **Positiv:** der Wächter sagt über lokale Referenzen etwas Wahres statt gar
  nichts. Gegengeprobt in beide Richtungen: ein vertippter Pfad ⇒
  `uses-local-missing`, der richtige ⇒ grün.
- **Negativ:** §3.9 trägt jetzt eine Ausnahme und ist damit länger als sein
  Kernsatz. Wer die Regel überfliegt, liest sie falsch — deshalb steht die
  Ausnahme **im selben Absatz**, nicht in einer Fußnote.
- **Negativ, benannt:** die Existenz-Prüfung sagt, dass die Datei **da** ist —
  nicht, dass sie einen passenden `workflow_call`-Eingang hat. Ein Aufruf mit
  falschen `inputs` fällt weiterhin erst zur Laufzeit auf.
- **Folgepflicht:** wächst die Zahl lokaler Referenzen über eine Handvoll, ist
  zu prüfen, ob der Wächter auch die `workflow_call`-Signatur vergleichen soll.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `workflow-pins.sh` | jeder **fremde** `uses:`-Eintrag trägt vollen SHA + Tag-Kommentar; jede **lokale** `./`-Referenz zeigt auf eine existierende Datei (`uses-local-missing`); die Zahl der lokalen Referenzen steht in der Erfolgsmeldung | `make workflow-pins` (in `gates`) |

## Re-Evaluierungs-Trigger

Führt GitHub eine Form ein, mit der eine lokale Workflow-Referenz **doch** an
einen Commit gebunden werden kann, entfällt die Begründung dieser Ausnahme — dann
gilt §3.9 wieder wörtlich.

Zweiter Trigger: rutscht jemals ein **fremder** Verweis über den `./`-Zweig
durch, ist die Erkennung zu weit gefasst und die Ausnahme neu zu schneiden.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-08-28 | Accepted — der Auftraggeber hat die ADR ausdrücklich verlangt, nachdem ich die Ausnahme zunächst nur in §3.9 dokumentiert hatte; die Frage „Präzisierung oder Lockerung" gehört in ein Entscheidungsprotokoll, nicht in einen Regelabsatz | [ADR-0065](0065-spiegel-gleichheit-ist-der-config-digest.md) |
